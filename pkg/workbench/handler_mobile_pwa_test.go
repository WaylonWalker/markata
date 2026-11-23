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
	shortcuts, ok := config["shortcuts"].([]any)
	if !ok || len(shortcuts) != 4 {
		t.Fatalf("mobile manifest must offer four shortcuts: %+v", config["shortcuts"])
	}
	wantShortcuts := map[string]bool{
		"/?compose=note":  false,
		"/?compose=photo": false,
		"/?compose=link":  false,
		"/?share_inbox=1": false,
	}
	for _, item := range shortcuts {
		entry, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("invalid mobile shortcut: %+v", item)
		}
		url, ok := entry["url"].(string)
		if !ok {
			t.Fatalf("invalid mobile shortcut URL: %+v", entry)
		}
		if _, exists := wantShortcuts[url]; !exists {
			t.Fatalf("unexpected shortcut %q", url)
		}
		wantShortcuts[url] = true
	}
	for url, seen := range wantShortcuts {
		if !seen {
			t.Fatalf("missing shortcut %q", url)
		}
	}
	if !strings.Contains(page.Body.String(), "create-camera-file") ||
		!strings.Contains(page.Body.String(), "/editor/share-target-utils.js") {
		t.Fatal("editor must include direct camera capture and shared file validation")
	}
	worker := get("/editor/share-sw.js")
	if worker.Header().Get("Service-Worker-Allowed") != "/" {
		t.Fatal("share worker must be allowed to handle root-scope navigation")
	}
	if !strings.Contains(worker.Body.String(), "MARKATA_ACTIVATE_UPDATE") ||
		strings.Contains(worker.Body.String(), "addEventListener('install', event => event.waitUntil(self.skipWaiting())") {
		t.Fatal("service worker updates must be deliberately activated")
	}
	offline := get("/editor/offline-capture.html")
	if !strings.Contains(offline.Body.String(), "markata-share-inbox") ||
		!strings.Contains(offline.Body.String(), "Offline capture") {
		t.Fatal("offline fallback must be a self-contained device-local capture UI")
	}
	if !strings.Contains(worker.Body.String(), "cache.add(OFFLINE_PAGE)") ||
		!strings.Contains(worker.Body.String(), "const OFFLINE_PAGE = '/editor/offline-capture.html'") {
		t.Fatal("share worker must precache only the public offline capture fallback")
	}
	if strings.Contains(worker.Body.String(), "cache.put(") ||
		strings.Contains(worker.Body.String(), "cache.addAll(") {
		t.Fatal("share worker must never cache private editor or workspace responses")
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
