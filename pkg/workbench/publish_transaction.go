package workbench

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
	"gopkg.in/yaml.v3"
)

type commitRequest struct {
	Path    string `json:"path"`
	Head    string `json:"head"`
	Version string `json:"version"`
}
type commitResult struct {
	Status  string `json:"status"`
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
	Message string `json:"message"`
}
type pushRequest struct {
	Commit string `json:"commit"`
	Branch string `json:"branch"`
}
type pushResult struct {
	Status  string `json:"status"`
	Commit  string `json:"commit"`
	Message string `json:"message"`
}

var gitHashPattern = regexp.MustCompile(`^[0-9a-f]{40,64}$`)
var branchPattern = regexp.MustCompile(`^refs/heads/[a-zA-Z0-9][a-zA-Z0-9/_-]*$`)

// commitGit runs a Git operation without interactive prompts or shell evaluation.
func commitGit(ctx context.Context, dir string, index string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_OPTIONAL_LOCKS=0")
	if index != "" {
		cmd.Env = append(cmd.Env, "GIT_INDEX_FILE="+index)
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s failed: %w", args[0], err)
	}
	if len(output) > 2<<20 {
		return "", fmt.Errorf("git response exceeds limit")
	}
	return strings.TrimSpace(string(output)), nil
}

func readPublishInput(w http.ResponseWriter, r *http.Request, dst any) bool {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		fail(w, http.StatusMethodNotAllowed, "POST required")
		return false
	}
	if !sameOriginWorkspaceMutation(r) {
		fail(w, http.StatusForbidden, "same-origin confirmation required")
		return false
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		fail(w, http.StatusBadRequest, "invalid publish request")
		return false
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		fail(w, http.StatusBadRequest, "unexpected additional request data")
		return false
	}
	return true
}

func validatePublicationMetadata(source string) error {
	// Do not silently rewrite YAML or publish a page still explicitly marked a draft.
	// Recognize only unambiguous top-level fields; reject malformed YAML and duplicates.
	if !strings.HasPrefix(source, "---
") && !strings.HasPrefix(source, "---
") {
		return fmt.Errorf("add YAML frontmatter and set published: true before committing")
	}
	end := strings.Index(source[4:], "
---")
	if end < 0 {
		return fmt.Errorf("invalid frontmatter boundary")
	}
	body := source[4 : 4+end+1]
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(body), &doc); err != nil {
		return fmt.Errorf("invalid YAML frontmatter: %w", err)
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return fmt.Errorf("frontmatter must be a YAML mapping")
	}
	var published, draft *bool
	keys := map[string]bool{}
	mapping := doc.Content[0]
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		key, value := mapping.Content[i], mapping.Content[i+1]
		if keys[key.Value] {
			return fmt.Errorf("duplicate frontmatter key %q", key.Value)
		}
		keys[key.Value] = true
		if key.Value != "published" && key.Value != "draft" {
			continue
		}
		var flag bool
		if value.Kind != yaml.ScalarNode || value.Tag != "!!bool" || value.Decode(&flag) != nil {
			return fmt.Errorf("%s must be a YAML boolean", key.Value)
		}
		if key.Value == "published" {
			published = &flag
		} else {
			draft = &flag
		}
	}
	if published == nil || !*published || (draft != nil && *draft) {
		return fmt.Errorf("publish requires published: true and draft not true; review frontmatter first")
	}
	return nil
}

func (h *Handler) publishPath(path string) (string, string, error) {
	// Reuse workspace restrictions to reject paths outside the configured content set.
	if !filepath.IsLocal(filepath.FromSlash(path)) || !strings.HasSuffix(strings.ToLower(path), ".md") {
		return "", "", fmt.Errorf("invalid Markdown path")
	}
	if _, err := h.workspace.Read(path); err != nil {
		return "", "", err
	}
	workspace, err := filepath.Abs(h.root)
	if err != nil {
		return "", "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	root, err := commitGit(ctx, workspace, "", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", fmt.Errorf("workspace requires a Git checkout")
	}
	absolute, err := filepath.EvalSymlinks(filepath.Join(workspace, filepath.FromSlash(path)))
	if err != nil {
		return "", "", err
	}
	relative, err := filepath.Rel(root, absolute)
	if err != nil || !filepath.IsLocal(relative) {
		return "", "", fmt.Errorf("file outside Git repository")
	}
	return root, filepath.ToSlash(relative), nil
}

func (h *Handler) servePublishCommit(w http.ResponseWriter, r *http.Request) {
	var req commitRequest
	if !readPublishInput(w, r, &req) {
		return
	}
	if !gitHashPattern.MatchString(req.Head) || !gitHashPattern.MatchString(req.Version) {
		fail(w, 400, "expected revision hashes are required")
		return
	}
	root, relative, err := h.publishPath(req.Path)
	if err != nil {
		fail(w, 422, err.Error())
		return
	}
	source, err := h.workspace.Read(req.Path)
	if err != nil {
		fail(w, 422, "cannot load current source")
		return
	}
	if editorworkspace.SourceVersion(source) != req.Version {
		fail(w, http.StatusConflict, "file changed since review; reopen review")
		return
	}
	if err := validatePublicationMetadata(source); err != nil {
		fail(w, 422, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	head, err := commitGit(ctx, root, "", "rev-parse", "HEAD")
	if err != nil || head != req.Head {
		fail(w, http.StatusConflict, "repository changed since review")
		return
	}
	ref, err := commitGit(ctx, root, "", "symbolic-ref", "-q", "HEAD")
	if err != nil || !branchPattern.MatchString(ref) {
		fail(w, 422, "publishing requires a named local branch")
		return
	}
	// The active index can hold unrelated staged changes; never include those
	// in a published commit. A staged version of this target is ambiguous.
	staged, err := commitGit(ctx, root, "", "diff", "--cached", "--name-only", "--", relative)
	if err != nil || staged != "" {
		fail(w, http.StatusConflict, "this file has staged changes; resolve them before publishing")
		return
	}
	temp, err := os.MkdirTemp("", "markata-publish-index-")
	if err != nil {
		fail(w, 500, "cannot create temporary Git index")
		return
	}
	defer os.RemoveAll(temp)
	index := filepath.Join(temp, "index")
	if _, err = commitGit(ctx, root, index, "read-tree", req.Head); err != nil {
		fail(w, 422, "cannot prepare isolated Git tree")
		return
	}
	if _, err = commitGit(ctx, root, index, "add", "--", relative); err != nil {
		fail(w, 422, "cannot add selected file to isolated Git tree")
		return
	}
	tree, err := commitGit(ctx, root, index, "write-tree")
	if err != nil {
		fail(w, 422, "cannot prepare publish tree")
		return
	}
	currentTree, err := commitGit(ctx, root, "", "rev-parse", "HEAD^{tree}")
	if err != nil {
		fail(w, 422, "cannot read committed tree")
		return
	}
	if tree == currentTree {
		fail(w, http.StatusConflict, "file has no new changes to commit")
		return
	}
	message := "Publish: " + relative
	commit, err := commitGit(ctx, root, "", "commit-tree", tree, "-p", req.Head, "-m", message)
	if err != nil {
		fail(w, 422, "commit failed; configure your Git user name and email")
		return
	}
	// Compare-and-swap is essential: no force updates and no stale commits.
	if _, err = commitGit(ctx, root, "", "update-ref", ref, commit, req.Head); err != nil {
		fail(w, http.StatusConflict, "repository changed while committing; reopen review")
		return
	}
	// Reconcile only the committed path in the *real* index. This preserves
	// unrelated staged files and ensures the editor doesn't appear dirty after
	// successfully committing its selected document.
	if _, err := commitGit(ctx, root, "", "reset", "-q", "HEAD", "--", relative); err != nil {
		fail(w, 500, "commit created but index reconciliation failed; inspect Git status")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(commitResult{
		Status: "committed-local", Commit: commit, Branch: ref,
		Message: "Committed locally. Nothing has been pushed or deployed.",
	})
}

func (h *Handler) servePublishPush(w http.ResponseWriter, r *http.Request) {
	var req pushRequest
	if !readPublishInput(w, r, &req) {
		return
	}
	if !gitHashPattern.MatchString(req.Commit) || !branchPattern.MatchString(req.Branch) {
		fail(w, 400, "expected commit and branch are required")
		return
	}
	root, _, err := h.publishPath(r.URL.Query().Get("path"))
	if err != nil {
		fail(w, 422, err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 35*time.Second)
	defer cancel()
	head, err := commitGit(ctx, root, "", "rev-parse", "HEAD")
	if err != nil || head != req.Commit {
		fail(w, http.StatusConflict, "branch changed; review Git status before pushing")
		return
	}
	ref, err := commitGit(ctx, root, "", "symbolic-ref", "-q", "HEAD")
	if err != nil || ref != req.Branch {
		fail(w, http.StatusConflict, "active branch changed")
		return
	}
	if _, err = commitGit(ctx, root, "", "remote", "get-url", "--push", "origin"); err != nil {
		fail(w, 422, "origin push remote is not configured")
		return
	}
	// Non-forced push; remote rejects non-fast-forwards. This may kick off the
	// user's external build, but it does not prove the site is deployed.
	if _, err = commitGit(ctx, root, "", "push", "--porcelain", "origin", ref+":"+ref); err != nil {
		fail(w, 422, "remote push failed; local commit is retained")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(pushResult{
		Status: "pushed-awaiting-deployment", Commit: req.Commit,
		Message: "Pushed to remote. Deployment has not been verified.",
	})
}
