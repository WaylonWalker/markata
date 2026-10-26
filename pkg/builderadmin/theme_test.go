package builderadmin

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLoadUITheme_UsesConfiguredSitePalette(t *testing.T) {
	sourceDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(sourceDir, "palettes"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "markata-go.toml"), []byte("[markata-go.theme]\npalette = \"builder\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "palettes", "builder-dark.toml"), []byte("[palette]\nname = \"Builder test\"\nvariant = \"dark\"\n\n[palette.colors]\nbg = \"#101112\"\npanel = \"#202122\"\nsurface = \"#303132\"\nelevated = \"#404142\"\ntext = \"#e0e1e2\"\nmuted = \"#a0a1a2\"\naccent = \"#778899\"\nsuccess = \"#11aa22\"\nwarning = \"#ddee33\"\nerror = \"#ee4455\"\ninfo = \"#3366dd\"\n\n[palette.semantic]\nbg-primary = \"bg\"\nbg-secondary = \"panel\"\nbg-surface = \"surface\"\nbg-elevated = \"elevated\"\ntext-primary = \"text\"\ntext-muted = \"muted\"\naccent = \"accent\"\nlink = \"accent\"\nborder = \"elevated\"\nborder-focus = \"accent\"\nsuccess = \"success\"\nwarning = \"warning\"\nerror = \"error\"\ninfo = \"info\"\n\n[palette.components]\ncode-bg = \"bg\"\ncode-text = \"text\"\nbutton-primary-bg = \"accent\"\nbutton-primary-text = \"bg\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	theme := loadUITheme(Config{SourceDir: sourceDir})
	if theme.Background != "#101112" || theme.Panel != "#202122" || theme.Accent != "#778899" {
		t.Fatalf("theme = %#v, want configured palette colors", theme)
	}
	if !theme.IsDark {
		t.Fatal("theme IsDark = false, want true")
	}
}

func TestSelectedPaletteName_UsesFallbackMode(t *testing.T) {
	if got := selectedPaletteName(models.ThemeConfig{Palette: "everforest", PaletteLight: "everforest-light", FallbackMode: "light"}); got != "everforest-light" {
		t.Fatalf("selectedPaletteName(light) = %q, want everforest-light", got)
	}
	if got := selectedPaletteName(models.ThemeConfig{Palette: "everforest", PaletteDark: "everforest-dark", FallbackMode: "dark"}); got != "everforest-dark" {
		t.Fatalf("selectedPaletteName(dark) = %q, want everforest-dark", got)
	}
}

func TestLoadUITheme_ResolvesRelativeConfigPathFromSource(t *testing.T) {
	sourceDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(sourceDir, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "config", "site.toml"), []byte("[markata-go.theme]\npalette = \"everforest-dark\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	theme := loadUITheme(Config{SourceDir: sourceDir, ConfigPath: "config/site.toml"})
	if theme.Background != "#2d353b" {
		t.Fatalf("theme.Background = %q, want everforest-dark background", theme.Background)
	}
}

func TestHandleIndex_RendersConfiguredTheme(t *testing.T) {
	sourceDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(sourceDir, "markata-go.toml"), []byte("[markata-go.theme]\npalette = \"everforest-dark\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc, err := New(Config{SourceDir: sourceDir, SiteDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	recorder := httptest.NewRecorder()
	svc.handleIndex(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	for _, want := range []string{"--markata-background:#2d353b", "--markata-success:#a7c080", "--markata-button-background:#a7c080"} {
		if !strings.Contains(recorder.Body.String(), want) {
			t.Errorf("rendered index missing %q", want)
		}
	}
}

func TestBuildDetail_RendersConfiguredTheme(t *testing.T) {
	sourceDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(sourceDir, "markata-go.toml"), []byte("[markata-go.theme]\npalette = \"everforest-dark\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	svc, err := New(Config{SourceDir: sourceDir, SiteDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = svc.leaderLock.Close() })
	svc.leader = true
	svc.state.Builds = []BuildRecord{{ID: "build-1"}}
	recorder := httptest.NewRecorder()
	svc.handleBuildDetail(recorder, httptest.NewRequest(http.MethodGet, "/builds/build-1", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), "--markata-background:#2d353b") {
		t.Error("rendered build detail is missing the configured background color")
	}
}
