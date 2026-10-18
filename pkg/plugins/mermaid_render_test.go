package plugins

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

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
