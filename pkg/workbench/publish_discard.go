package workbench

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

// Local workbench only. Hosted draft discard must use the durable sync journal.
const discardVersionKey = "version"

type discardRequest struct {
	Path    string `json:"path"`
	Head    string `json:"head"`
	Version string `json:"version"`
}

func (h *Handler) servePublishDiscard(w http.ResponseWriter, r *http.Request) {
	var req discardRequest
	if !readPublishInput(w, r, &req) {
		return
	}
	// Protect the source/version check and write from other requests to this editor.
	// External filesystem editors still require a separate conflict protocol.
	h.workspaceWriteMu.Lock()
	defer h.workspaceWriteMu.Unlock()
	if !gitHashPattern.MatchString(req.Head) || !gitHashPattern.MatchString(req.Version) {
		fail(w, http.StatusBadRequest, "review revision and Git HEAD required")
		return
	}
	root, relative, err := h.publishPath(req.Path)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "file is not a valid workspace Markdown document")
		return
	}
	source, err := h.workspace.Read(req.Path)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot read current draft")
		return
	}
	if editorworkspace.SourceVersion(source) != req.Version {
		fail(w, http.StatusConflict, "draft changed after review; reopen Review changes")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	head, err := gitOutput(ctx, root, "rev-parse", "--verify", "HEAD")
	if err != nil || strings.TrimSpace(head) != req.Head {
		fail(w, http.StatusConflict, "Git baseline changed after review")
		return
	}
	staged, err := gitOutput(ctx, root, "diff", "--cached", "--name-only", "--", relative)
	if err != nil || strings.TrimSpace(staged) != "" {
		fail(w, http.StatusConflict, "this file has staged changes; discard refused")
		return
	}
	// Resolve only a blob in the reviewed HEAD. Untracked files need a separate
	// guarded deletion path rather than an unchecked filesystem removal.
	baseline, err := gitOutput(ctx, root, "show", req.Head+":"+relative)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "new file: discard deletion is not supported yet")
		return
	}
	if source == baseline {
		fail(w, http.StatusConflict, "no changes remain to discard")
		return
	}
	// Re-read immediately before the workspace write. This is the local
	// workbench's optimistic check, not a hosted multi-device transaction.
	again, err := h.workspace.Read(req.Path)
	if err != nil || editorworkspace.SourceVersion(again) != req.Version {
		fail(w, http.StatusConflict, "draft changed during discard")
		return
	}
	if err := h.workspace.Write(req.Path, baseline); err != nil {
		fail(w, http.StatusUnprocessableEntity, "could not restore committed source")
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":          "restored-local",
		discardVersionKey: editorworkspace.SourceVersion(baseline),
	}); err != nil {
		return
	}
}
