package plugins

import (
	"context"
	"encoding/json"
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
		"lineColor":           colors.Accent,
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
