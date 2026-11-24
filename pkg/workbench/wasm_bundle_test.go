package workbench

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestHandler_OptionalBundledWasmPreview(t *testing.T) {
	h := newHandler(t)
	index := httptest.NewRecorder()
	h.ServeHTTP(index, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody))
	if index.Code != http.StatusOK {
		t.Fatalf("GET / returned %d", index.Code)
	}
	if !strings.Contains(index.Body.String(), "MARKATA_WASM_VERSION=") {
		t.Fatal("editor does not expose WASM capability handshake")
	}
	wantFingerprint := wasmPreviewFingerprint(h.wasmVersion, h.siteConfig)
	if !strings.Contains(index.Body.String(), `MARKATA_PREVIEW_FINGERPRINT="`+wantFingerprint+`";`) {
		t.Fatalf("editor preview fingerprint does not match active renderer config: %q", wantFingerprint)
	}

	worker := httptest.NewRecorder()
	h.ServeHTTP(worker, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/editor/wasm/worker.js", http.NoBody))
	if worker.Code != http.StatusOK || !strings.Contains(worker.Body.String(), "markata.wasm") {
		t.Fatalf("bundled worker not served: %d", worker.Code)
	}

	_, bundled := web.ReadFile("web/wasm/markata.wasm")
	response := httptest.NewRecorder()
	h.ServeHTTP(response, httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/editor/wasm/markata.wasm", http.NoBody))
	if bundled == nil {
		if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "application/wasm" {
			t.Fatalf("bundled WASM returned %d, MIME %q", response.Code, response.Header().Get("Content-Type"))
		}
		if response.Header().Get("X-Markata-Wasm-Version") != h.wasmVersion {
			t.Fatalf("WASM version response = %q, want %q", response.Header().Get("X-Markata-Wasm-Version"), h.wasmVersion)
		}
		if strings.Contains(index.Body.String(), "MARKATA_WASM_VERSION=\"\";") {
			t.Fatal("embedded WASM was not advertised")
		}
	} else if response.Code != http.StatusNotFound || !strings.Contains(index.Body.String(), "MARKATA_WASM_VERSION=\"\";") {
		t.Fatalf("native-only build must safely disable WASM; status=%d", response.Code)
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
			request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "https://edit.waylonwalker.com"+asset, http.NoBody)
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
			assetName := strings.TrimPrefix(asset, "/editor/")
			if _, err := web.ReadFile("web/" + assetName); err != nil {
				wantStatus = http.StatusNotFound
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

func TestHandler_WasmPrecompressedResponse(t *testing.T) {
	raw, err := web.ReadFile("web/wasm/markata.wasm")
	if err != nil {
		t.Skip("native-only build; no bundled WASM")
	}
	compressed, err := web.ReadFile("web/wasm/markata.wasm.gz")
	if err != nil {
		t.Skip("build did not include the optional precompressed WASM")
	}
	h := newHandler(t)
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/editor/wasm/markata.wasm", http.NoBody)
	request.Header.Set("Accept-Encoding", "gzip, deflate")
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
	if response.Code != http.StatusOK || response.Header().Get("Content-Encoding") != "gzip" ||
		response.Header().Get("Vary") != "Accept-Encoding" ||
		response.Header().Get("Content-Type") != "application/wasm" ||
		response.Header().Get("X-Markata-Wasm-Version") != wasmPreviewFingerprint(h.wasmVersion, nil) {
		t.Fatalf("WASM gzip status=%d encoding=%q vary=%q type=%q",
			response.Code, response.Header().Get("Content-Encoding"),
			response.Header().Get("Vary"), response.Header().Get("Content-Type"))
	}
	if !bytes.Equal(response.Body.Bytes(), compressed) {
		t.Fatal("compressed WASM response does not match build artifact")
	}
	reader, err := gzip.NewReader(bytes.NewReader(response.Body.Bytes()))
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := io.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	if err := reader.Close(); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(decoded, raw) {
		t.Fatal("gzip response did not decode to the release-matched WASM runtime")
	}
	request.Header.Set("Accept-Encoding", "identity")
	plain := httptest.NewRecorder()
	h.ServeHTTP(plain, request)
	if plain.Code != http.StatusOK || plain.Header().Get("Content-Encoding") != "" ||
		!bytes.Equal(plain.Body.Bytes(), raw) {
		t.Fatalf("WASM identity fallback returned status=%d encoding=%q",
			plain.Code, plain.Header().Get("Content-Encoding"))
	}
}

func TestWasmPreviewFingerprintRequiresDefaultSiteRenderer(t *testing.T) {
	if got := wasmPreviewFingerprint("abc123", nil); got != "abc123" {
		t.Fatalf("unconfigured renderer fingerprint = %q, want bundled version", got)
	}
	if got := wasmPreviewFingerprint("abc123", &models.Config{}); got != "" {
		t.Fatalf("configured renderer fingerprint = %q, want WASM disabled", got)
	}
	if got := wasmPreviewFingerprint("", nil); got != "" {
		t.Fatalf("native-only renderer fingerprint = %q, want empty", got)
	}
}
