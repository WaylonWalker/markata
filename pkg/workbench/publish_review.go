package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

// publishReview is deliberately read-only. HEAD means "last committed", not
// "deployed"; no endpoint may call this a live/published-site revision.
type publishReview struct {
	Path           string `json:"path"`
	Head           string `json:"head"`
	Baseline       string `json:"baseline"`
	Current        string `json:"current"`
	BaselineExists bool   `json:"baseline_exists"`
	Changed        bool   `json:"changed"`
	Version        string `json:"version"`
	BaselineKind   string `json:"baseline_kind"`
	CanPublish     bool   `json:"can_publish"`
}

// gitOutput runs bounded, non-interactive reads. It is never used to execute
// an editor-supplied Git command or mutate repository state.
func gitOutput(ctx context.Context, dir string, args ...string) (string, error) {
	// #nosec G204 G702 -- args are fixed Git subcommands assembled by this package; no shell is used.
	command := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	command.Env = append(command.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0")
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	if len(output) > 2<<20 {
		return "", fmt.Errorf("git response is too large")
	}
	return string(output), nil
}

func (h *Handler) servePublishReview(w http.ResponseWriter, r *http.Request) {
	// serveWorkspace already enforces loopback-only access.
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		fail(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	path := r.URL.Query().Get("path")
	if path == "" || !filepath.IsLocal(filepath.FromSlash(path)) {
		fail(w, http.StatusBadRequest, "choose a file from the workspace")
		return
	}
	markdown := strings.HasSuffix(strings.ToLower(path), ".md")
	var current string
	var err error
	if markdown {
		current, err = h.workspace.Read(path)
		if err != nil {
			fail(w, http.StatusNotFound, err.Error())
			return
		}
	}
	absoluteWorkspace, err := filepath.Abs(h.root)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "invalid workspace root")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	root, err := gitOutput(ctx, absoluteWorkspace, "rev-parse", "--show-toplevel")
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "this workspace is not in a Git repository")
		return
	}
	root = strings.TrimSpace(root)
	resolvedPath, pathErr := filepath.EvalSymlinks(filepath.Join(absoluteWorkspace, filepath.FromSlash(path)))
	if pathErr != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot resolve source path")
		return
	}
	relative, err := filepath.Rel(root, resolvedPath)
	if err != nil || !filepath.IsLocal(relative) {
		fail(w, http.StatusUnprocessableEntity, "file is outside Git repository")
		return
	}
	if !markdown {
		current, err = readChangedWorkspaceFile(ctx, root, resolvedPath, relative)
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
	}
	head, err := gitOutput(ctx, absoluteWorkspace, "rev-parse", "--verify", "HEAD")
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "the repository has no commit to compare against")
		return
	}
	// The colon form anchors the path in this exact revision; it does not
	// resolve arbitrary refs or allow a user to supply Git revision syntax.
	head = strings.TrimSpace(head)
	baseline, err := gitOutput(ctx, absoluteWorkspace, "show", head+":"+filepath.ToSlash(relative))
	exists := err == nil
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			fail(w, http.StatusUnprocessableEntity, "cannot read committed revision")
			return
		}
		// Missing tracked files are a valid "new post" case. A missing blob
		// must not be mistaken for an existing empty file.
		baseline = ""
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(publishReview{
		Path: path, Head: head, Baseline: baseline, Current: current,
		BaselineExists: exists, Changed: !exists || baseline != current,
		Version:      editorworkspace.SourceVersion(current),
		BaselineKind: "last-committed", CanPublish: false,
	}); err != nil {
		return
	}
}

// readChangedWorkspaceFile returns the current text of a non-Markdown file only
// when Git reports it as changed, so this read-only endpoint cannot browse the
// repository.
func readChangedWorkspaceFile(ctx context.Context, root, resolvedPath, relative string) (string, error) {
	status, err := gitOutput(ctx, root, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--", filepath.ToSlash(relative))
	if err != nil || strings.TrimSpace(status) == "" {
		return "", errors.New("file has no Git changes to compare")
	}
	file, err := os.Open(resolvedPath)
	if err != nil {
		return "", errors.New("cannot read current file")
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxSource+1))
	if err != nil {
		return "", errors.New("cannot read current file")
	}
	if len(data) > maxSource {
		return "", errors.New("file exceeds 200,000 bytes")
	}
	if bytes.IndexByte(data, 0) >= 0 {
		return "", errors.New("binary files cannot be compared")
	}
	return string(data), nil
}
