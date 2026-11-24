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
