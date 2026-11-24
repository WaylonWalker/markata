package workbench

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_OptionalBundledWasmPreview(t *testing.T) {
	h := newHandler(t)
	index := httptest.NewRecorder()
	h.ServeHTTP(index, httptest.NewRequest(http.MethodGet, "/", http.NoBody))
	if index.Code != http.StatusOK {
		t.Fatalf("GET / returned %d", index.Code)
	}
	if !strings.Contains(index.Body.String(), "MARKATA_WASM_VERSION=") {
		t.Fatal("editor does not expose WASM capability handshake")
	}

	worker := httptest.NewRecorder()
	h.ServeHTTP(worker, httptest.NewRequest(http.MethodGet, "/editor/wasm/worker.js", http.NoBody))
	if worker.Code != http.StatusOK || !strings.Contains(worker.Body.String(), "markata.wasm") {
		t.Fatalf("bundled worker not served: %d", worker.Code)
	}

	_, bundled := web.ReadFile("web/wasm/markata.wasm")
	response := httptest.NewRecorder()
	h.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/editor/wasm/markata.wasm", http.NoBody))
	if bundled == nil {
		if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "application/wasm" {
			t.Fatalf("bundled WASM returned %d, MIME %q", response.Code, response.Header().Get("Content-Type"))
		}
		if strings.Contains(index.Body.String(), "MARKATA_WASM_VERSION=\"\";") {
			t.Fatal("embedded WASM was not advertised")
		}
	} else {
		if response.Code != http.StatusNotFound || !strings.Contains(index.Body.String(), "MARKATA_WASM_VERSION=\"\";") {
			t.Fatalf("native-only build must safely disable WASM; status=%d", response.Code)
		}
	}
}

func TestHandler_HostedWasmAssetsKeepAuthBoundary(t *testing.T) {
	h := newHandler(t)
	if err := h.EnableHostedAuth(HostedAuthConfig{
		PublicOrigin:      "https://edit.waylonwalker.com",
		TrustedProxyCIDRs: []string{"10.42.0.0/24"},
		RequiredScope:     "markata-editor:admin",
	}); err != nil {
		t.Fatal(err)
	}
	for _, asset := range []string{
		"/editor/wasm/worker.js",
		"/editor/wasm/wasm_exec.js",
		"/editor/wasm/markata.wasm",
	} {
		t.Run(asset, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "https://edit.waylonwalker.com"+asset, http.NoBody)
			request.RemoteAddr = "10.42.0.12:8443"
			denied := httptest.NewRecorder()
			h.ServeHTTP(denied, request)
			if denied.Code != http.StatusUnauthorized || denied.Header().Get("Cache-Control") != "private, no-store" {
				t.Fatalf("unauthorized WASM asset: status=%d cache=%q",
					denied.Code, denied.Header().Get("Cache-Control"))
			}
			request.Header.Set("X-Hlab-User-Id", "usr_123")
			request.Header.Set("X-Hlab-Scopes", "markata-editor:admin")
			allowed := httptest.NewRecorder()
			h.ServeHTTP(allowed, request)
			wantStatus := http.StatusOK
			if asset != "/editor/wasm/worker.js" {
				assetName := strings.TrimPrefix(asset, "/editor/")
				if _, err := web.ReadFile("web/" + assetName); err != nil {
					wantStatus = http.StatusNotFound
				}
			}
			if allowed.Code != wantStatus || allowed.Header().Get("Cache-Control") != "private, no-store" ||
				allowed.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Fatalf("authorized WASM asset: status=%d want=%d cache=%q cors=%q",
					allowed.Code, wantStatus, allowed.Header().Get("Cache-Control"),
					allowed.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}
