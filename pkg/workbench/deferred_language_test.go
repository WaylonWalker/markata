package workbench

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeferredLanguageBecomesAvailableWithoutBlockingEditor(t *testing.T) {
	service := &deferredLanguage{}
	req := httptest.NewRequest(http.MethodGet, "/lsp", http.NoBody)
	first := httptest.NewRecorder()
	service.ServeHTTP(first, req)
	if first.Code != http.StatusServiceUnavailable || first.Header().Get("Retry-After") == "" {
		t.Fatalf("starting language service: %d %q", first.Code, first.Body.String())
	}
	service.ready(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	next := httptest.NewRecorder()
	service.ServeHTTP(next, req)
	if next.Code != http.StatusNoContent {
		t.Fatalf("ready language service not available: %d", next.Code)
	}
}
