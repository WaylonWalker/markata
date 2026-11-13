package workbench

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// publishReadiness gates one-click publication. A Git push sends a branch,
// not a file: only an exact remote/local HEAD match proves that publishing
// the next isolated single-page commit will not also publish older commits.
type publishReadiness struct {
	Ready  bool   `json:"ready"`
	Head   string `json:"head"`
	Branch string `json:"branch"`
}

func (h *Handler) servePublishReadiness(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		fail(w, http.StatusMethodNotAllowed, "GET required")
		return
	}
	root, _, err := h.publishPath(r.URL.Query().Get("path"))
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "select a publishable site Markdown post")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	head, err := commitGit(ctx, root, "", "rev-parse", "HEAD")
	if err != nil || !gitHashPattern.MatchString(head) {
		fail(w, http.StatusServiceUnavailable, "cannot read current local revision")
		return
	}
	branch, err := commitGit(ctx, root, "", "symbolic-ref", "-q", "HEAD")
	if err != nil || !branchPattern.MatchString(branch) {
		fail(w, http.StatusConflict, "select a named Git branch before publishing")
		return
	}
	if _, err = commitGit(ctx, root, "", "remote", "get-url", "--push", "origin"); err != nil {
		fail(w, http.StatusUnprocessableEntity, "a Git origin remote is required for publishing")
		return
	}
	remote, err := commitGit(ctx, root, "", "ls-remote", "--heads", "origin", branch)
	if err != nil {
		fail(w, http.StatusServiceUnavailable, "cannot verify remote Git state; nothing published")
		return
	}
	parts := strings.Fields(remote)
	if len(parts) != 2 || parts[1] != branch || !gitHashPattern.MatchString(parts[0]) {
		fail(w, http.StatusConflict, "remote branch is missing or ambiguous; nothing published")
		return
	}
	if parts[0] != head {
		fail(w, http.StatusConflict, "local branch differs from origin. Review and sync pending commits before publishing this post.")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(publishReadiness{Ready: true, Head: head, Branch: branch}); err != nil {
		return
	}
}
