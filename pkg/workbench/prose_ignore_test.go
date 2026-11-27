package workbench

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/prose"
)

func TestProseIgnoreRoutePersistsProjectRule(t *testing.T) {
	root := t.TempDir()
	h := &Handler{root: root}
	req := httptest.NewRequest(http.MethodPost, "http://localhost/workspace/prose/ignore-rule",
		strings.NewReader(`{"rule":"prose/cliche"}`))
	req.RemoteAddr = "127.0.0.1:42000"
	req.Header.Set("Origin", "http://localhost")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	h.ServeHTTP(response, req)
	if response.Code != http.StatusOK {
		t.Fatalf("ignore rule request returned %d: %s", response.Code, response.Body.String())
	}
	prefs, err := prose.LoadPreferences(filepath.Join(root, "post.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !prefs.Options().IgnoredRules["prose/cliche"] {
		t.Fatalf("browser preference not shared with CLI/LSP: %+v", prefs)
	}
}

func TestProseIgnoreRouteRejectsCrossOriginAndUnknownRules(t *testing.T) {
	for _, tc := range []struct {
		name   string
		origin string
		body   string
	}{
		{name: "cross-origin", origin: "https://example.org", body: `{"rule":"prose/cliche"}`},
		{name: "unknown rule", origin: "http://localhost", body: `{"rule":"example/unknown"}`},
		{name: "arbitrary path", origin: "http://localhost", body: `{"rule":"prose/cliche","path":"../../other"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			h := &Handler{root: root}
			req := httptest.NewRequest(http.MethodPost, "http://localhost/workspace/prose/ignore-rule",
				strings.NewReader(tc.body))
			req.RemoteAddr = "127.0.0.1:42001"
			req.Header.Set("Origin", tc.origin)
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			h.ServeHTTP(response, req)
			if response.Code < 400 {
				t.Fatalf("expected rejection, got HTTP %d", response.Code)
			}
			if _, err := os.Stat(filepath.Join(root, prose.PreferencesFile)); !os.IsNotExist(err) {
				t.Fatalf("unexpected prose preference file after rejected request: %v", err)
			}
		})
	}
}
