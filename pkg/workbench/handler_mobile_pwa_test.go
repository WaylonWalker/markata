package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMobilePWAAssetsAndShareTarget(t *testing.T) {
	h := newHandler(t)
	get := func(path string) *httptest.ResponseRecorder {
		t.Helper()
		response := httptest.NewRecorder()
		h.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://localhost"+path, http.NoBody))
		if response.Code != http.StatusOK {
			t.Fatalf("GET %s: %d %s", path, response.Code, response.Body.String())
		}
		return response
	}
	page := get("/")
	if !strings.Contains(page.Body.String(), "rel=\"manifest\" href=\"/editor/manifest.webmanifest") ||
		!strings.Contains(page.Body.String(), "/editor/mobile-share.js") {
		t.Fatal("editor page did not load the mobile PWA manifest and share inbox")
	}
	manifest := get("/editor/manifest.webmanifest")
	if got := manifest.Header().Get("Content-Type"); got != "application/manifest+json" {
		t.Fatalf("manifest type = %q", got)
	}
	var config map[string]any
	if err := json.Unmarshal(manifest.Body.Bytes(), &config); err != nil {
		t.Fatal(err)
	}
	target, ok := config["share_target"].(map[string]any)
	if !ok || config["display"] != "standalone" || target["action"] != "/share-target" ||
		target["method"] != "POST" || target["enctype"] != "multipart/form-data" {
		t.Fatalf("invalid mobile share manifest: %+v", config)
	}
	worker := get("/editor/share-sw.js")
	if worker.Header().Get("Service-Worker-Allowed") != "/" {
		t.Fatal("share worker must be allowed to handle root-scope navigation")
	}
	if strings.Contains(worker.Body.String(), "caches.open") {
		t.Fatal("share worker must not cache private workspace data")
	}
	for _, path := range []string{"/editor/icon-192.png", "/editor/icon-512.png"} {
		icon := get(path)
		if icon.Header().Get("Content-Type") != "image/png" ||
			!bytes.HasPrefix(icon.Body.Bytes(), []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
			t.Fatalf("invalid icon %s", path)
		}
	}
	// Only the installed service worker should handle OS shares.
	// Never expose a public server-side POST route bypassing workspace ACL.
	post := httptest.NewRecorder()
	h.ServeHTTP(post, httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost/share-target", strings.NewReader("untrusted")))
	if post.Code != http.StatusMethodNotAllowed {
		t.Fatalf("unhandled public share POST = %d, want 405", post.Code)
	}
}
