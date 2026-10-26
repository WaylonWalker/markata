package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestRedirectsPlugin_HTMLFallbackDefaultsEnabled(t *testing.T) {
	p := NewRedirectsPlugin()
	if !p.HTMLFallbackEnabled() {
		t.Fatal("HTML fallback should be enabled by default")
	}
}

func TestRedirectsPlugin_Write_GeneratesNginxConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	redirectsFile := filepath.Join(tmpDir, "_redirects")

	redirectsContent := `/old /new
/go/docs https://docs.example.com/start
`
	if err := os.WriteFile(redirectsFile, []byte(redirectsContent), 0o600); err != nil {
		t.Fatalf("write redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir

	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: redirectsFile})
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, nginxRedirectsFilename))
	if err != nil {
		t.Fatalf("read generated nginx config: %v", err)
	}

	got := string(content)
	for _, want := range []string{
		`location = "/old" {`,
		`return 301 "/new";`,
		`location = "/go/docs" {`,
		`return 301 "https://docs.example.com/start";`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("redirects.conf missing %q\n%s", want, got)
		}
	}
}

func TestRedirectsPlugin_Write_EmptySourceClearsNativeRules(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	if err := os.WriteFile(redirectsFile, []byte("/old /new\n"), 0o600); err != nil {
		t.Fatalf("write redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: redirectsFile})
	if err := p.Write(m); err != nil {
		t.Fatalf("initial Write() error = %v", err)
	}

	if err := os.WriteFile(redirectsFile, nil, 0o600); err != nil {
		t.Fatalf("empty redirects file: %v", err)
	}
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() after clearing source error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, nginxRedirectsFilename))
	if err != nil {
		t.Fatalf("read refreshed nginx config: %v", err)
	}
	if strings.Contains(string(content), "location =") || strings.Contains(string(content), "return 301") {
		t.Fatalf("stale native redirect remained after source was emptied:\n%s", content)
	}
}

func TestRedirectsPlugin_Write_MissingSourceRemovesGeneratedNativeConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	if err := os.WriteFile(redirectsFile, []byte("/old /new\n"), 0o600); err != nil {
		t.Fatalf("write redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: redirectsFile})
	if err := p.Write(m); err != nil {
		t.Fatalf("initial Write() error = %v", err)
	}

	if err := os.Remove(redirectsFile); err != nil {
		t.Fatalf("remove redirects source: %v", err)
	}
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() after source removal error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, nginxRedirectsFilename)); !os.IsNotExist(err) {
		t.Fatalf("generated redirects.conf should be removed; stat error = %v", err)
	}
}

func TestRedirectsPlugin_Write_MissingSourcePreservesUserNativeConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	userConfig := "# user managed\nlocation = /keep { return 302 /still-here; }\n"
	outputPath := filepath.Join(outputDir, nginxRedirectsFilename)
	if err := os.WriteFile(outputPath, []byte(userConfig), 0o600); err != nil {
		t.Fatalf("write user config: %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: filepath.Join(tmpDir, "missing", "_redirects")})
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read user config: %v", err)
	}
	if string(content) != userConfig {
		t.Fatalf("user-managed redirects.conf was modified:\n%s", content)
	}
}

func TestRedirectsPlugin_Write_RefusesToOverwriteUserNativeConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	redirectsFile := filepath.Join(tmpDir, "_redirects")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("create output dir: %v", err)
	}
	if err := os.WriteFile(redirectsFile, []byte("/old /new\n"), 0o600); err != nil {
		t.Fatalf("write redirects file: %v", err)
	}

	userConfig := "# user managed\nlocation = /keep { return 302 /still-here; }\n"
	outputPath := filepath.Join(outputDir, nginxRedirectsFilename)
	if err := os.WriteFile(outputPath, []byte(userConfig), 0o600); err != nil {
		t.Fatalf("write user config: %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	p := NewRedirectsPlugin()
	p.SetConfig(RedirectsConfig{RedirectsFile: redirectsFile})
	err := p.Write(m)
	if err == nil {
		t.Fatal("expected Write() to refuse user-managed redirects.conf")
	}
	if !strings.Contains(err.Error(), "refusing to overwrite user-managed") {
		t.Fatalf("unexpected error: %v", err)
	}

	content, readErr := os.ReadFile(outputPath)
	if readErr != nil {
		t.Fatalf("read user config: %v", readErr)
	}
	if string(content) != userConfig {
		t.Fatalf("user-managed redirects.conf was modified:\n%s", content)
	}
}

func TestRedirectsPlugin_Write_HTMLFallbackCanBeDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "output")
	redirectsFile := filepath.Join(tmpDir, "_redirects")

	if err := os.WriteFile(redirectsFile, []byte("/old /new\n"), 0o600); err != nil {
		t.Fatalf("write redirects file: %v", err)
	}

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = outputDir
	cfg.Extra = map[string]interface{}{
		"redirects": map[string]interface{}{
			"redirects_file": redirectsFile,
			"html_fallback":  false,
		},
	}

	p := NewRedirectsPlugin()
	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure() error = %v", err)
	}
	if p.HTMLFallbackEnabled() {
		t.Fatal("HTML fallback should be disabled by config")
	}
	if err := p.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputDir, nginxRedirectsFilename)); err != nil {
		t.Fatalf("expected redirects.conf: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "old", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("HTML fallback should not be generated; stat error = %v", err)
	}
}

func TestQuoteNginxString(t *testing.T) {
	got := quoteNginxString(`/old"$host`)
	want := `"/old\"\$host"`
	if got != want {
		t.Fatalf("quoteNginxString() = %q, want %q", got, want)
	}
}
