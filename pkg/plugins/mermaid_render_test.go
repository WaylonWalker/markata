package plugins

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

func TestMermaidPaletteThemeVariablesUseThemeSurfaces(t *testing.T) {
	colors := &mermaidPaletteColors{
		Background: "#111111",
		TextColor:  "#eeeeee",
		Accent:     "#ffcc00",
		Surface:    "#222222",
	}

	var got map[string]string
	if err := json.Unmarshal([]byte(colors.themeVariablesJSON()), &got); err != nil {
		t.Fatalf("themeVariablesJSON() returned invalid JSON: %v", err)
	}

	for key, want := range map[string]string{
		"primaryColor":        colors.Surface,
		"nodeBkg":             colors.Surface,
		"edgeLabelBackground": colors.Surface,
		"clusterBkg":          colors.Background,
		"primaryTextColor":    colors.TextColor,
		"nodeTextColor":       colors.TextColor,
		"clusterTextColor":    colors.TextColor,
		"primaryBorderColor":  colors.Accent,
		"nodeBorder":          colors.Accent,
		"lineColor":           colors.TextColor,
		"clusterBorder":       colors.Accent,
	} {
		if got[key] != want {
			t.Errorf("theme variable %q = %q, want %q", key, got[key], want)
		}
	}
}

func TestLoadOrDownloadJSSource_UsesBundledMermaidDir(t *testing.T) {
	bundledDir := t.TempDir()
	want := "window.mermaid = {};"
	cacheFile := filepath.Join(bundledDir, "mermaid-v"+mermaidJSVersion+".min.js")
	if err := os.WriteFile(cacheFile, []byte(want), 0o600); err != nil {
		t.Fatalf("failed to write bundled MermaidJS source: %v", err)
	}

	t.Setenv(runtimeenv.EnvBundledMermaidDir, bundledDir)
	t.Setenv(runtimeenv.EnvOffline, "true")
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	got, err := loadOrDownloadJSSource(context.Background(), mermaidJSVersion)
	if err != nil {
		t.Fatalf("loadOrDownloadJSSource() error = %v", err)
	}
	if got != want {
		t.Fatalf("loadOrDownloadJSSource() = %q, want %q", got, want)
	}
}

func TestLoadOrDownloadJSSource_OfflineMissingFails(t *testing.T) {
	t.Setenv(runtimeenv.EnvBundledMermaidDir, "")
	t.Setenv(runtimeenv.EnvOffline, "true")
	t.Setenv("XDG_CACHE_HOME", t.TempDir())

	_, err := loadOrDownloadJSSource(context.Background(), mermaidJSVersion)
	if err == nil {
		t.Fatal("expected offline MermaidJS load to fail when cache is empty")
	}
	if !strings.Contains(err.Error(), "offline") {
		t.Fatalf("expected offline error, got %v", err)
	}
}

func TestDownloadMermaidJSSource(t *testing.T) {
	const version = "10.9.5"
	want := "window.mermaid = {};"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/npm/mermaid@"+version+"/dist/mermaid.min.js" {
			http.NotFound(w, r)
			return
		}
		_, _ = io.WriteString(w, want)
	}))
	defer server.Close()

	client := &http.Client{Transport: mermaidDownloadRoundTripper(func(request *http.Request) (*http.Response, error) {
		request = request.Clone(request.Context())
		request.URL.Scheme = "http"
		request.URL.Host = strings.TrimPrefix(server.URL, "http://")
		return http.DefaultTransport.RoundTrip(request)
	})}
	got, err := downloadMermaidJSSource(context.Background(), version, client)
	if err != nil {
		t.Fatalf("downloadMermaidJSSource() error = %v", err)
	}
	if got != want {
		t.Fatalf("downloadMermaidJSSource() = %q, want %q", got, want)
	}
}

func TestDownloadMermaidJSSource_RejectsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	defer server.Close()
	client := &http.Client{Transport: mermaidDownloadRoundTripper(func(request *http.Request) (*http.Response, error) {
		request = request.Clone(request.Context())
		request.URL.Scheme = "http"
		request.URL.Host = strings.TrimPrefix(server.URL, "http://")
		return http.DefaultTransport.RoundTrip(request)
	})}

	_, err := downloadMermaidJSSource(context.Background(), "10.9.5", client)
	if err == nil || !strings.Contains(err.Error(), "HTTP status 404") {
		t.Fatalf("downloadMermaidJSSource() error = %v, want HTTP status 404", err)
	}
}

type mermaidDownloadRoundTripper func(*http.Request) (*http.Response, error)

func (roundTripper mermaidDownloadRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTripper(request)
}
