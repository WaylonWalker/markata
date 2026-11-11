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
	// Resolve only a blob in the reviewed HEAD. A missing baseline may
	// represent a new draft, but never infer that it is safe to delete merely
	// from a failed git-show: independently verify exact untracked status.
	baseline, err := gitOutput(ctx, root, "show", req.Head+":"+relative)
	if err != nil {
		h.serveUntrackedDiscard(w, ctx, req, root, relative)
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

// serveUntrackedDiscard is the local-only deletion path. It never operates on
// files that Git considers tracked, staged, ignored, or otherwise ambiguous.
func (h *Handler) serveUntrackedDiscard(
	w http.ResponseWriter, ctx context.Context, req discardRequest, root, relative string,
) {
	status, statusErr := gitOutput(ctx, root, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--", relative)
	if statusErr != nil || status != "?? "+relative+"\x00" {
		fail(w, http.StatusConflict, "cannot discard: file is tracked, staged, ignored, or Git status is uncertain")
		return
	}
	// Never fall back to raw os.Remove for an unsupported workspace source.
	remover, ok := h.workspace.(interface{ Delete(string, string) error })
	if !ok {
		fail(w, http.StatusUnprocessableEntity, "guarded deletion is not available in this workspace")
		return
	}
	again, readErr := h.workspace.Read(req.Path)
	if readErr != nil || editorworkspace.SourceVersion(again) != req.Version {
		fail(w, http.StatusConflict, "draft changed during discard")
		return
	}
	if err := remover.Delete(req.Path, req.Version); err != nil {
		fail(w, http.StatusConflict, "could not remove new draft safely: "+err.Error())
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "removed-local"}); err != nil {
		return
	}
}
