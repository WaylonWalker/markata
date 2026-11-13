package editorhub

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

const privateSessionHTTPS = "https"

// PrivateSessionHTTP returns the current verified subject and authorized sites.
// It never returns cookies, tokens, proxy assertions or Git credentials.
// No public editor mounts this handler until the real ForwardAuth path is proven.
type PrivateSessionHTTP struct {
	registry     *Registry
	authenticate Authenticator
	origin       string
}

func NewPrivateSessionHandler(registry *Registry, authenticate Authenticator, allowedOrigin string) (http.Handler, error) {
	origin, err := url.Parse(allowedOrigin)
	if err != nil || origin == nil || origin.Scheme != privateSessionHTTPS || origin.Host == "" ||
		origin.User != nil || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" ||
		registry == nil || authenticate == nil {
		return nil, errors.New("private session needs a registry, trusted authenticator and exact HTTPS origin")
	}
	h := &PrivateSessionHTTP{registry: registry, authenticate: authenticate, origin: allowedOrigin}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/session", h.serve)
	return mux, nil
}

func (h *PrivateSessionHTTP) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
	w.Header().Set("Vary", "Origin, Cookie, Authorization")
	if r.Header.Get("Sec-Fetch-Site") == privateWriteCrossSite ||
		r.Header.Get("Sec-Fetch-Site") == privateCrossOrigin {
		http.Error(w, "cross-site session request forbidden", http.StatusForbidden)
		return
	}
	origins := r.Header.Values("Origin")
	if len(origins) > 1 || (len(origins) == 1 && origins[0] != h.origin) {
		http.Error(w, "wrong session origin", http.StatusForbidden)
		return
	}
	subject, err := h.authenticate(r)
	if err != nil || !validReadSubject(subject) {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(struct {
		SubjectID string        `json:"subject_id"`
		Sites     []SiteSummary `json:"sites"`
	}{
		SubjectID: subject,
		Sites:     h.registry.ListFor(subject),
	}); err != nil {
		return // Response may already be committed.
	}
}
