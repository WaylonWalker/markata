package workbench

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/WaylonWalker/markata/pkg/prose"
)

// serveProseDictionaryAdd persists a known spelling term only at the
// server-selected workspace root. The browser cannot supply a filesystem
// destination or arbitrary preference keys through this endpoint.
func (h *Handler) serveProseDictionaryAdd(w http.ResponseWriter, r *http.Request) {
	if !h.sameOriginWorkspaceMutation(r) {
		fail(w, http.StatusForbidden, "same-origin workspace writes are required")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		fail(w, http.StatusUnsupportedMediaType, "send Content-Type: application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var request struct {
		Word string `json:"word"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		fail(w, http.StatusBadRequest, "invalid dictionary request")
		return
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		fail(w, http.StatusBadRequest, "expected one dictionary request")
		return
	}
	if _, err := prose.CanonicalWord(request.Word); err != nil {
		fail(w, http.StatusBadRequest, "invalid dictionary word")
		return
	}
	h.workspaceWriteMu.Lock()
	defer h.workspaceWriteMu.Unlock()
	if _, err := prose.AddWordAtRoot(h.root, request.Word); err != nil {
		fail(w, http.StatusUnprocessableEntity, "could not save dictionary word")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(map[string]bool{"ok": true}); err != nil {
		return
	}
}
