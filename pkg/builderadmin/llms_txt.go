package builderadmin

import (
	_ "embed"
	"io"
	"net/http"
)

// builderAdminLLMsTXT is served directly by the admin service rather than being
// produced by the site's publishing pipeline.
//
//go:embed llms.txt
var builderAdminLLMsTXT string

// handleLLMsTXT serves the same read-only guide to browser operators and
// bearer-authenticated agents. Authentication is applied by registerRoutes.
func (s *Service) handleLLMsTXT(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "private, no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if r.Method == http.MethodHead {
		return
	}
	_, _ = io.WriteString(w, builderAdminLLMsTXT)
}
