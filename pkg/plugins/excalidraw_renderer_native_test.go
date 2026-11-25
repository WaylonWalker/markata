//go:build !js

package plugins

import (
	"bytes"
	"context"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WaylonWalker/markata/pkg/visualizations"
)

func TestExcalidrawRendererServesUpstreamFontURL(t *testing.T) {
	assets, err := visualizations.ExcalidrawAssetFS()
	if err != nil {
		t.Fatal(err)
	}
	fontPath := "excalidraw-fonts/Excalifont/Excalifont-Regular-a88b72a24fb54c9f94e3b5fdaa7481c9.woff2"
	want, err := fs.ReadFile(assets, fontPath)
	if err != nil {
		t.Fatal(err)
	}
	handler := excalidrawFontAssetAlias(http.FileServer(http.FS(assets)))

	response := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(context.Background(), http.MethodGet,
		"http://localhost/excalidraw-fonts/fonts/Excalifont/Excalifont-Regular-a88b72a24fb54c9f94e3b5fdaa7481c9.woff2", http.NoBody)
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("upstream font URL returned %d: %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Content-Type") != "font/woff2" {
		t.Errorf("font content type = %q, want font/woff2", response.Header().Get("Content-Type"))
	}
	if !bytes.Equal(response.Body.Bytes(), want) {
		t.Fatal("upstream font URL did not return the bundled font bytes")
	}
}
