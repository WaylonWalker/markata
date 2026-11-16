package editorhub

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/editorsync"
)

const privateReadCrossSite = "cross-site"

// PrivateReadHTTP is deliberately unmounted. No browser can access it unless
// an application explicitly registers it behind a verified trusted ingress.
type PrivateReadHTTP struct {
	registry     *Registry
	authenticate Authenticator
	reader       *editorsync.PrivateReader
	origin       string
}

func NewPrivateReadHandler(registry *Registry, authenticate Authenticator, reader *editorsync.PrivateReader, allowedOrigin string) (http.Handler, error) {
	parsed, err := url.Parse(allowedOrigin)
	if err != nil || parsed == nil || parsed.Scheme != "https" || parsed.Host == "" ||
		parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" ||
		registry == nil || authenticate == nil || reader == nil {
		return nil, errors.New("private read handler requires verified auth, site registry, private reader and exact HTTPS origin")
	}
	h := &PrivateReadHTTP{registry: registry, authenticate: authenticate, reader: reader, origin: allowedOrigin}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/sites/{site}/snapshots", h.snapshot)
	mux.HandleFunc("GET /api/v1/sites/{site}/documents", h.documents)
	mux.HandleFunc("GET /api/v1/sites/{site}/changes", h.changes)
	mux.HandleFunc("GET /api/v1/sites/{site}/documents/{document}", h.source)
	return mux, nil
}

func (h *PrivateReadHTTP) authorize(w http.ResponseWriter, r *http.Request) (string, bool) {
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Vary", "Cookie, Authorization, Origin")
	if r.Header.Get("Sec-Fetch-Site") == privateReadCrossSite || r.Header.Get("Sec-Fetch-Site") == "cross-origin" {
		http.Error(w, "cross-site access forbidden", http.StatusForbidden)
		return "", false
	}
	subject, err := h.authenticate(r)
	if err != nil || !validReadSubject(subject) {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return "", false
	}
	if !h.registry.Allows(subject, r.PathValue("site"), PermissionRead) {
		http.Error(w, "read forbidden", http.StatusForbidden)
		return "", false
	}
	return subject, true
}

func validReadSubject(subject string) bool {
	return subject != "" && len(subject) <= 256 && strings.TrimSpace(subject) == subject &&
		!strings.ContainsAny(subject, "\x00\r\n") && utf8.ValidString(subject)
}

func privateReadPageLimit(r *http.Request, allowCursor bool) (limit int, cursor string, resultErr error) {
	q := r.URL.Query()
	for name, values := range q {
		if name != "limit" && (name != "cursor" || !allowCursor) || len(values) != 1 {
			return 0, "", errors.New("invalid read query")
		}
	}
	limit = 50
	if raw := q.Get("limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > 100 {
			return 0, "", errors.New("invalid page size")
		}
		limit = n
	}
	if allowCursor && q.Get("cursor") == "" {
		return 0, "", errors.New("snapshot cursor is required")
	}
	return limit, q.Get("cursor"), nil
}

// Stable JSON shapes avoid exposing Go implementation details or per-site sequence.
type privateHTTPManifestItem struct {
	DocumentID string `json:"document_id"`
	Revision   int64  `json:"revision"`
	Hash       string `json:"hash"`
}
type privateHTTPSnapshotPage struct {
	Items         []privateHTTPManifestItem `json:"items"`
	NextCursor    string                    `json:"next_cursor,omitempty"`
	ChangesCursor string                    `json:"changes_cursor,omitempty"`
}

func snapshotHTTPPage(page editorsync.PrivateSnapshotPage) privateHTTPSnapshotPage {
	out := privateHTTPSnapshotPage{Items: make([]privateHTTPManifestItem, 0, len(page.Items)),
		NextCursor: page.NextCursor, ChangesCursor: page.ChangesCursor}
	for _, item := range page.Items {
		out.Items = append(out.Items, privateHTTPManifestItem{DocumentID: item.DocumentID, Revision: item.Revision, Hash: item.Hash})
	}
	return out
}

type privateHTTPChange struct {
	DocumentID string `json:"document_id"`
	Revision   int64  `json:"revision"`
	Deleted    bool   `json:"deleted"`
}
type privateHTTPChangesPage struct {
	Items      []privateHTTPChange `json:"items"`
	NextCursor string              `json:"next_cursor"`
}

func changesHTTPPage(page editorsync.PrivateChangesPage) privateHTTPChangesPage {
	out := privateHTTPChangesPage{Items: make([]privateHTTPChange, 0, len(page.Items)), NextCursor: page.NextCursor}
	for _, item := range page.Items {
		out.Items = append(out.Items, privateHTTPChange{DocumentID: item.DocumentID, Revision: item.Revision, Deleted: item.Deleted})
	}
	return out
}

func privateReadJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		return // The response may already have been committed.
	}
}

func privateReadError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, editorsync.ErrReadForbidden):
		http.Error(w, "read forbidden", http.StatusForbidden)
	case errors.Is(err, editorsync.ErrPrivateVersionNotFound):
		http.Error(w, "document not found", http.StatusNotFound)
	case errors.Is(err, editorsync.ErrExpiredSyncCursor):
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusGone)
		if err := json.NewEncoder(w).Encode(map[string]string{"code": "checkpoint_expired"}); err != nil {
			return // The response may already have been committed.
		}
	case errors.Is(err, editorsync.ErrInvalidSyncCursor):
		http.Error(w, "invalid cursor", http.StatusBadRequest)
	default:
		http.Error(w, "private read temporarily unavailable", http.StatusServiceUnavailable)
	}
}

func (h *PrivateReadHTTP) snapshot(w http.ResponseWriter, r *http.Request) {
	subject, ok := h.authorize(w, r)
	if !ok {
		return
	}
	if origins := r.Header.Values("Origin"); len(origins) != 1 || origins[0] != h.origin ||
		r.Header.Get("Sec-Fetch-Site") == privateReadCrossSite || r.ContentLength != 0 {
		http.Error(w, "invalid snapshot request", http.StatusForbidden)
		return
	}
	limit, _, err := privateReadPageLimit(r, false)
	if err != nil {
		http.Error(w, "invalid page size", http.StatusBadRequest)
		return
	}
	page, err := h.reader.Snapshot(r.Context(), subject, r.PathValue("site"), "", limit)
	if err != nil {
		privateReadError(w, err)
		return
	}
	privateReadJSON(w, snapshotHTTPPage(page))
}

func (h *PrivateReadHTTP) documents(w http.ResponseWriter, r *http.Request) {
	subject, ok := h.authorize(w, r)
	if !ok {
		return
	}
	limit, cursor, err := privateReadPageLimit(r, true)
	if err != nil {
		http.Error(w, "invalid snapshot paging query", http.StatusBadRequest)
		return
	}
	page, err := h.reader.Snapshot(r.Context(), subject, r.PathValue("site"), cursor, limit)
	if err != nil {
		privateReadError(w, err)
		return
	}
	privateReadJSON(w, snapshotHTTPPage(page))
}

func (h *PrivateReadHTTP) changes(w http.ResponseWriter, r *http.Request) {
	subject, ok := h.authorize(w, r)
	if !ok {
		return
	}
	limit, cursor, err := privateReadPageLimit(r, true)
	if err != nil {
		http.Error(w, "invalid change feed query", http.StatusBadRequest)
		return
	}
	page, err := h.reader.Changes(r.Context(), subject, r.PathValue("site"), cursor, limit)
	if err != nil {
		privateReadError(w, err)
		return
	}
	privateReadJSON(w, changesHTTPPage(page))
}

func (h *PrivateReadHTTP) source(w http.ResponseWriter, r *http.Request) {
	subject, ok := h.authorize(w, r)
	if !ok {
		return
	}
	q := r.URL.Query()
	if len(q) != 1 || len(q["revision"]) != 1 {
		http.Error(w, "revision required", http.StatusBadRequest)
		return
	}
	revision, err := strconv.ParseInt(q.Get("revision"), 10, 64)
	if err != nil || revision < 1 {
		http.Error(w, "invalid revision", http.StatusBadRequest)
		return
	}
	raw, err := h.reader.SourceAtRevision(r.Context(), subject, r.PathValue("site"), r.PathValue("document"), revision)
	if err != nil {
		privateReadError(w, err)
		return
	}
	// Exact bytes, including non-UTF-8 source, survive JSON transport.
	// No filenames, filesystem paths, or other subjects' metadata are emitted.
	privateReadJSON(w, map[string]any{"revision": revision, "source_base64": base64.StdEncoding.EncodeToString(raw)})
}
