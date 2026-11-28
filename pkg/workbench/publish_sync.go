package workbench

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

const publishSyncStatus = "synced"

type publishSyncRequest struct {
	Path    string `json:"path"`
	Head    string `json:"head"`
	Version string `json:"version"`
}

type publishSyncResult struct {
	Status       string `json:"status"`
	Head         string `json:"head"`
	RemoteHead   string `json:"remote_head,omitempty"`
	Message      string `json:"message,omitempty"`
	RemoteSource string `json:"remote_source,omitempty"`
	RemoteExists bool   `json:"remote_exists"`
}

// servePublishSync makes one guarded attempt to align the editor checkout with
// origin. It only fast-forwards or rebases one reviewed, single-file commit.
//
//nolint:gocyclo // The explicit Git safety gates make the sync transaction auditable.
func (h *Handler) servePublishSync(w http.ResponseWriter, r *http.Request) {
	var req publishSyncRequest
	if !h.readPublishInput(w, r, &req) {
		return
	}
	if !gitHashPattern.MatchString(req.Head) || !gitHashPattern.MatchString(req.Version) {
		fail(w, http.StatusBadRequest, "expected revision and source hashes are required")
		return
	}
	root, relative, err := h.publishPath(req.Path)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "select a publishable site Markdown post")
		return
	}
	source, err := h.workspace.Read(req.Path)
	if err != nil || editorworkspace.SourceVersion(source) != req.Version {
		fail(w, http.StatusConflict, "saved post changed; review the current version before syncing")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	head, err := commitGit(ctx, root, "", "rev-parse", "HEAD")
	if err != nil || head != req.Head {
		fail(w, http.StatusConflict, "repository changed; refresh the review before syncing")
		return
	}
	branch, err := commitGit(ctx, root, "", "symbolic-ref", "-q", "HEAD")
	if err != nil || !branchPattern.MatchString(branch) {
		fail(w, http.StatusConflict, "publishing requires a named local branch")
		return
	}
	if err := rejectActiveGitOperation(ctx, root); err != nil {
		fail(w, http.StatusConflict, "finish the active Git operation before syncing")
		return
	}
	if _, err := commitGit(ctx, root, "", "remote", "get-url", "--push", "origin"); err != nil {
		fail(w, http.StatusUnprocessableEntity, "a Git origin remote is required for publishing")
		return
	}

	remoteHead, err := h.fetchPublishRemote(ctx, root, branch, head)
	if err != nil {
		fail(w, http.StatusServiceUnavailable, "cannot fetch origin; no local changes were changed")
		return
	}
	staged, err := publishGitOutput(ctx, root, "diff", "--cached", "--name-only", "-z")
	if err != nil {
		fail(w, http.StatusServiceUnavailable, "cannot inspect staged Git changes")
		return
	}
	if len(staged) != 0 {
		h.writePublishSyncReview(ctx, w, root, relative, head, remoteHead,
			"Staged changes are present. Your saved files are safe; review them before syncing.")
		return
	}
	if remoteHead == head {
		writePublishSyncResult(w, http.StatusOK, publishSyncResult{Status: publishSyncStatus, Head: head, RemoteHead: remoteHead})
		return
	}

	base, err := commitGit(ctx, root, "", "merge-base", head, remoteHead)
	if err != nil || !gitHashPattern.MatchString(base) {
		h.writePublishSyncReview(ctx, w, root, relative, head, remoteHead,
			"Origin has a different Git history. Your saved version is safe; compare it with the origin copy before choosing what to keep.")
		return
	}

	if base == head {
		if changed, checkErr := publishGitPathChanged(ctx, root, head, remoteHead, relative); checkErr != nil {
			fail(w, http.StatusServiceUnavailable, "cannot compare the post with origin")
			return
		} else if changed {
			h.writePublishSyncReview(ctx, w, root, relative, head, remoteHead,
				"Origin changed this post. Your saved version is safe; compare both copies before choosing what to keep.")
			return
		}
		incoming, err := publishChangedPaths(ctx, root, head, remoteHead)
		if err != nil {
			fail(w, http.StatusServiceUnavailable, "cannot inspect files changed on origin")
			return
		}
		working, err := publishWorkingPaths(ctx, root, head)
		if err != nil {
			fail(w, http.StatusServiceUnavailable, "cannot inspect saved workspace changes")
			return
		}
		if overlap := intersectPaths(incoming, working); len(overlap) != 0 {
			h.writePublishSyncReview(ctx, w, root, relative, head, remoteHead,
				"Origin and this workspace both changed at least one file. Your saved files are safe; review the overlap before syncing.")
			return
		}
		if current, err := commitGit(ctx, root, "", "rev-parse", "HEAD"); err != nil || current != head {
			fail(w, http.StatusConflict, "repository changed during sync; refresh the review")
			return
		}
		if _, err := commitGit(ctx, root, "", "merge", "--ff-only", "--no-edit", remoteHead); err != nil {
			h.writePublishSyncReview(ctx, w, root, relative, head, remoteHead,
				"Git could not safely fast-forward this workspace. Your saved version is safe; compare it with origin.")
			return
		}
		newHead, err := commitGit(ctx, root, "", "rev-parse", "HEAD")
		if err != nil || newHead != remoteHead {
			fail(w, http.StatusServiceUnavailable, "sync result could not be confirmed; inspect the local Git state")
			return
		}
		writePublishSyncResult(w, http.StatusOK, publishSyncResult{Status: publishSyncStatus, Head: newHead, RemoteHead: remoteHead})
		return
	}

	if base == remoteHead {
		h.writePublishSyncReview(ctx, w, root, relative, head, remoteHead,
			"This workspace has unpublished commits. Review them before publishing this post.")
		return
	}
	if err := h.rebaseReviewedPublishCommit(ctx, root, relative, head, base, remoteHead); err != nil {
		h.writePublishSyncReview(ctx, w, root, relative, head, remoteHead, err.Error())
		return
	}
	newHead, err := commitGit(ctx, root, "", "rev-parse", "HEAD")
	if err != nil || !gitHashPattern.MatchString(newHead) {
		fail(w, http.StatusServiceUnavailable, "sync result could not be confirmed; inspect the local Git state")
		return
	}
	writePublishSyncResult(w, http.StatusOK, publishSyncResult{Status: publishSyncStatus, Head: newHead, RemoteHead: remoteHead})
}

func (h *Handler) fetchPublishRemote(ctx context.Context, root, branch, head string) (string, error) {
	var remote string
	var err error
	if h.gitPublisher != nil {
		remote, err = h.gitPublisher.FetchRemote(ctx, branch, head)
	} else {
		ref := gitPublisherSyncRef(branch)
		_, err = commitGit(ctx, root, "", "fetch", "--no-tags", "origin", "+"+branch+":"+ref)
		if err == nil {
			remote, err = commitGit(ctx, root, "", "rev-parse", "--verify", ref)
		}
	}
	if err != nil || !gitHashPattern.MatchString(remote) {
		return "", errors.New("remote fetch failed")
	}
	return remote, nil
}

func (h *Handler) rebaseReviewedPublishCommit(ctx context.Context, root, relative, head, base, remoteHead string) error {
	count, err := commitGit(ctx, root, "", "rev-list", "--count", base+".."+head)
	if err != nil || count != "1" {
		return errors.New("only this post's single unpublished commit can be synced automatically; review other local commits before publishing")
	}
	changed, err := publishChangedPaths(ctx, root, base, head)
	if err != nil || len(changed) != 1 || !changed[relative] {
		return errors.New("the unpublished commit contains other files; review local commits before publishing")
	}
	if changedOnOrigin, err := publishGitPathChanged(ctx, root, base, remoteHead, relative); err != nil {
		return errors.New("cannot compare this post with origin")
	} else if changedOnOrigin {
		return errors.New("origin changed this post; compare both copies before choosing what to keep")
	}
	remoteChanges, err := publishChangedPaths(ctx, root, base, remoteHead)
	if err != nil {
		return errors.New("cannot inspect files changed on origin")
	}
	working, err := publishWorkingPaths(ctx, root, head)
	if err != nil {
		return errors.New("cannot inspect saved workspace changes")
	}
	if working[relative] {
		return errors.New("this post has additional unsaved edits; review it before syncing the committed copy")
	}
	if overlap := intersectPaths(remoteChanges, working); len(overlap) != 0 {
		return errors.New("origin and this workspace changed at least one file; review the overlap before syncing")
	}
	if current, err := commitGit(ctx, root, "", "rev-parse", "HEAD"); err != nil || current != head {
		return errors.New("repository changed during sync; refresh the review")
	}
	if _, err := commitGit(ctx, root, "", "rebase", "--autostash", "--onto", remoteHead, base); err != nil {
		_, abortErr := commitGit(context.WithoutCancel(ctx), root, "", "rebase", "--abort")
		if abortErr != nil {
			return errors.New("git could not complete the safe rebase; inspect the repository before trying again")
		}
		return errors.New("git could not safely reapply this post; compare it with origin before trying again")
	}
	if current, err := commitGit(ctx, root, "", "rev-parse", "HEAD"); err != nil || current == head {
		return errors.New("git did not create an updated commit")
	}
	return nil
}

func (h *Handler) writePublishSyncReview(ctx context.Context, w http.ResponseWriter, root, relative, head, remoteHead, message string) {
	result := publishSyncResult{
		Status: "review-required", Head: head, RemoteHead: remoteHead,
		Message: message,
	}
	remoteSource, exists, err := publishSourceAtRevision(ctx, root, remoteHead, relative)
	switch {
	case err == nil && exists:
		result.RemoteSource = remoteSource
		result.RemoteExists = true
	case err == nil:
		result.RemoteExists = false
	default:
		result.Message += " The origin copy could not be loaded here; keep the local file and review origin in Git."
	}
	writePublishSyncResult(w, http.StatusConflict, result)
}

func writePublishSyncResult(w http.ResponseWriter, status int, result publishSyncResult) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		return
	}
}

func rejectActiveGitOperation(ctx context.Context, root string) error {
	for _, marker := range []string{"MERGE_HEAD", "CHERRY_PICK_HEAD", "REVERT_HEAD", "rebase-merge", "rebase-apply", "BISECT_LOG"} {
		path, err := commitGit(ctx, root, "", "rev-parse", "--git-path", marker)
		if err != nil {
			return err
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(root, path)
		}
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("git operation %s is active", marker)
		} else if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	return nil
}

func publishGitOutput(ctx context.Context, root string, args ...string) ([]byte, error) {
	// #nosec G204 G702 -- arguments are fixed Git commands assembled from validated revisions and paths.
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if len(output) > 2<<20 {
		return nil, errors.New("git response exceeds limit")
	}
	return output, nil
}

func publishGitPathChanged(ctx context.Context, root, left, right, path string) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", root, "diff", "--quiet", "--no-ext-diff", left, right, "--", filepath.ToSlash(path))
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0")
	err := cmd.Run()
	if err == nil {
		return false, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == 1 {
		return true, nil
	}
	return false, err
}

func publishChangedPaths(ctx context.Context, root, left, right string) (map[string]bool, error) {
	output, err := publishGitOutput(ctx, root, "diff", "--name-only", "-z", left, right)
	if err != nil {
		return nil, err
	}
	return gitNULPaths(output), nil
}

func publishWorkingPaths(ctx context.Context, root, head string) (map[string]bool, error) {
	tracked, err := publishGitOutput(ctx, root, "diff", "--name-only", "-z", head)
	if err != nil {
		return nil, err
	}
	untracked, err := publishGitOutput(ctx, root, "ls-files", "--others", "--exclude-standard", "-z")
	if err != nil {
		return nil, err
	}
	paths := gitNULPaths(tracked)
	for path := range gitNULPaths(untracked) {
		paths[path] = true
	}
	return paths, nil
}

func gitNULPaths(output []byte) map[string]bool {
	paths := make(map[string]bool)
	for _, path := range strings.Split(string(output), "\x00") {
		if path != "" {
			paths[filepath.ToSlash(path)] = true
		}
	}
	return paths
}

func intersectPaths(left, right map[string]bool) []string {
	var overlap []string
	for path := range left {
		if right[path] {
			overlap = append(overlap, path)
		}
	}
	return overlap
}

func publishSourceAtRevision(ctx context.Context, root, revision, path string) (source string, exists bool, err error) {
	output, err := publishGitOutput(ctx, root, "ls-tree", "-z", "--full-tree", revision, "--", filepath.ToSlash(path))
	if err != nil {
		return "", false, err
	}
	if len(output) == 0 {
		return "", false, nil
	}
	entry := strings.SplitN(strings.TrimSuffix(string(output), "\x00"), "\t", 2)
	if len(entry) != 2 || entry[1] != filepath.ToSlash(path) {
		return "", false, errors.New("cannot identify origin file")
	}
	fields := strings.Fields(entry[0])
	if len(fields) != 3 || fields[1] != "blob" || !gitHashPattern.MatchString(fields[2]) || fields[0] != "100644" && fields[0] != "100755" {
		return "", false, errors.New("origin file is not a regular Markdown document")
	}
	content, err := publishGitOutput(ctx, root, "cat-file", "blob", fields[2])
	if err != nil {
		return "", false, err
	}
	if len(content) > maxSource || strings.ContainsRune(string(content), '\x00') {
		return "", false, errors.New("origin file exceeds the review limit")
	}
	return string(content), true, nil
}
