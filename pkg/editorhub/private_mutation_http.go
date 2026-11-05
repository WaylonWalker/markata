package editorhub

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/editorsync"
)

// privateMutationHTTP is intentionally unmounted. Only a hosted application
// with a proven trusted ForwardAuth path can safely register this handler.
// No browser-provided subject or builder operator capability is accepted.
type privateMutationHTTP struct {
	registry     *Registry
	authenticate Authenticator
	engine       *editorsync.Engine
	origin       string
}

// NewPrivateMutationHandler creates one private-draft write route. It must
// never be mounted beside the loopback workbench or Builder Admin routes.
// allowedOrigin is the exact externally visible HTTPS editor origin; it is
// not learned from browser-controlled Host or forwarded headers.
func NewPrivateMutationHandler(registry *Registry, authenticate Authenticator, engine *editorsync.Engine, allowedOrigin string) (http.Handler, error) {
	parsed, err := url.Parse(allowedOrigin)
	if err != nil || parsed == nil || parsed.Scheme != "https" || parsed.Host == "" ||
		parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New("private mutation handler needs an explicit HTTPS editor origin")
	}
	if registry == nil || authenticate == nil || engine == nil {
		return nil, errors.New("private mutation handler needs a registry, trusted authenticator and SQLite writer")
	}
	handler := &privateMutationHTTP{
		registry: registry, authenticate: authenticate, engine: engine, origin: allowedOrigin,
	}
	mux := http.NewServeMux()
	mux.Handle("POST /api/v1/sites/{site}/mutations", http.HandlerFunc(handler.apply))
	mux.Handle("GET /api/v1/sites/{site}/conflicts/{mutation}", http.HandlerFunc(handler.getConflict))
	return mux, nil
}

func (h *privateMutationHTTP) apply(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Vary", "Origin, Cookie, Authorization")
	if !h.validWriteOrigin(r) {
		http.Error(w, "cross-site write forbidden", http.StatusForbidden)
		return
	}
	subject, err := h.authenticate(r)
	if err != nil || !validMutationSubject(subject) {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	site := r.PathValue("site")
	if !h.registry.Allows(subject, site, PermissionWrite) {
		http.Error(w, "write forbidden", http.StatusForbidden)
		return
	}
	mutation, err := decodePrivateMutation(w, r, site, subject)
	if err != nil {
		var sizeError *http.MaxBytesError
		if errors.As(err, &sizeError) {
			http.Error(w, "mutation too large", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "invalid private mutation", http.StatusBadRequest)
		}
		return
	}
	receipt, err := h.engine.Apply(r.Context(), subject, mutation)
	if err != nil {
		h.writeMutationError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// The internal per-site change sequence may contain other authors' private
	// edits. Never expose that sequence in a browser receipt.
	if err := json.NewEncoder(w).Encode(struct {
		MutationID string `json:"mutation_id"`
		Revision   int64  `json:"revision"`
		Digest     string `json:"digest"`
	}{
		MutationID: receipt.MutationID,
		Revision:   receipt.Revision,
		Digest:     receipt.PayloadDigest,
	}); err != nil {
		return // Response may already be committed.
	}
}

// getConflict returns a durable private rejection copy only to its original
// authenticated author; site write access alone is not sufficient.
func (h *privateMutationHTTP) getConflict(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Vary", "Origin, Cookie, Authorization")
	if r.Header.Get("Sec-Fetch-Site") == "cross-site" ||
		r.Header.Get("Sec-Fetch-Site") == "cross-origin" {
		http.Error(w, "cross-site request forbidden", http.StatusForbidden)
		return
	}
	if origins := r.Header.Values("Origin"); len(origins) > 1 ||
		(len(origins) == 1 && origins[0] != h.origin) {
		http.Error(w, "wrong editor origin", http.StatusForbidden)
		return
	}
	subject, err := h.authenticate(r)
	if err != nil || !validMutationSubject(subject) {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	site := r.PathValue("site")
	if !h.registry.Allows(subject, site, PermissionWrite) {
		http.Error(w, "write forbidden", http.StatusForbidden)
		return
	}
	record, err := h.engine.ReadPrivateConflict(r.Context(), subject, site, r.PathValue("mutation"))
	if err != nil {
		if errors.Is(err, editorsync.ErrConflictNotFound) {
			// Unknown and another author's conflict must be indistinguishable.
			http.Error(w, "conflict not found", http.StatusNotFound)
		} else if errors.Is(err, editorsync.ErrWriteForbidden) {
			http.Error(w, "write forbidden", http.StatusForbidden)
		} else {
			http.Error(w, "private conflict temporarily unavailable", http.StatusServiceUnavailable)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(struct {
		MutationID        string `json:"mutation_id"`
		DocumentID        string `json:"document_id"`
		StaleBaseRevision int64  `json:"stale_base_revision"`
		CurrentRevision   int64  `json:"current_revision"`
		Digest            string `json:"digest"`
		SourceBase64      string `json:"source_base64"`
		Deleted           bool   `json:"deleted"`
	}{
		MutationID: record.MutationID,
		DocumentID: record.DocumentID,
		StaleBaseRevision: record.StaleBaseRevision,
		CurrentRevision: record.CurrentRevision,
		Digest: record.PayloadDigest,
		SourceBase64: base64.StdEncoding.EncodeToString(record.Source),
		Deleted: record.Deleted,
	}); err != nil {
		return // Response may already be committed.
	}
}

func (h *privateMutationHTTP) validWriteOrigin(r *http.Request) bool {
	if r.Header.Get("Sec-Fetch-Site") == "cross-site" ||
		r.Header.Get("Sec-Fetch-Site") == "cross-origin" {
		return false
	}
	origins := r.Header.Values("Origin")
	if len(origins) != 1 || origins[0] != h.origin {
		return false
	}
	intents := r.Header.Values("X-Markata-Sync-Intent")
	return len(intents) == 1 && intents[0] == "write"
}

func validMutationSubject(subject string) bool {
	return subject != "" && len(subject) <= 256 && utf8.ValidString(subject) &&
		strings.TrimSpace(subject) == subject && !strings.ContainsAny(subject, "\r\n\x00")
}

type privateMutationRequest struct {
	DocumentID   string `json:"document_id"`
	MutationID   string `json:"mutation_id"`
	BaseRevision uint64 `json:"base_revision"`
	Source       string `json:"source"`
	Deleted      bool   `json:"deleted"`
}

func decodePrivateMutation(w http.ResponseWriter, r *http.Request, site, subject string) (editorsync.Mutation, error) {
	contentType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || contentType != "application/json" {
		return editorsync.Mutation{}, errors.New("invalid mutation content type")
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20+4096)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var body privateMutationRequest
	if err := decoder.Decode(&body); err != nil {
		return editorsync.Mutation{}, err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return editorsync.Mutation{}, errors.New("invalid trailing mutation JSON")
	}
	mutation := editorsync.Mutation{
		SiteID: site, SubjectID: subject,
		DocumentID: body.DocumentID, MutationID: body.MutationID,
		BaseRevision: body.BaseRevision, Visibility: editorsync.Private,
		Source: []byte(body.Source), Deleted: body.Deleted,
	}
	if err := mutation.Validate(); err != nil {
		return editorsync.Mutation{}, err
	}
	return mutation, nil
}

func (*privateMutationHTTP) writeMutationError(w http.ResponseWriter, err error) {
	if errors.Is(err, editorsync.ErrWriteForbidden) {
		http.Error(w, "write forbidden", http.StatusForbidden)
		return
	}
	if errors.Is(err, editorsync.ErrMutationIDReuse) {
		http.Error(w, "mutation ID reused", http.StatusConflict)
		return
	}
	var conflict *editorsync.ConflictError
	if errors.As(err, &conflict) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusConflict)
		if err := json.NewEncoder(w).Encode(struct {
			Code            string `json:"code"`
			CurrentRevision int64  `json:"current_revision"`
			Preserved       bool   `json:"preserved"`
		}{"revision_conflict", conflict.CurrentRevision, conflict.Preserved}); err != nil {
			return // Response may already be committed.
		}
		return
	}
	// Database and transaction errors are not reflected to the browser.
	http.Error(w, "private mutation could not be committed", http.StatusServiceUnavailable)
}
