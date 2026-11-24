package plugins

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestPaletteCSSPlugin_CorePreservesDefaultsAndDefersOtherPalettes(t *testing.T) {
	// Exercise only shipped palettes, independent of the developer's collection.
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("APPDATA", t.TempDir())
	includeAll := true
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprint(enabled), func(t *testing.T) {
			m := lifecycle.NewManager()
			m.SetConfig(&lifecycle.Config{OutputDir: t.TempDir(), Extra: map[string]interface{}{
				"models_config": &models.Config{Theme: models.ThemeConfig{Palette: "ayu-dark", Aesthetic: "precision"}},
				"theme": models.ThemeConfig{
					Palette: "ayu-dark", PaletteLight: "ayu-light", PaletteDark: "ayu-dark",
					Switcher:  models.ThemeSwitcherConfig{Enabled: &enabled, IncludeAll: &includeAll},
					Variables: map[string]string{"--test-override": "42px"},
				},
			}})
			p := NewPaletteCSSPlugin()
			if err := p.Configure(m); err != nil {
				t.Fatal(err)
			}
			if err := p.Write(m); err != nil {
				t.Fatal(err)
			}
			corePath := filepath.Join(m.Config().OutputDir, "css", "palette-core.css")
			core, err := os.ReadFile(corePath)
			if !enabled {
				if !os.IsNotExist(err) || m.Config().Extra["palette_core_enabled"] == true {
					t.Fatal("switcher disabled must retain the single stylesheet")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			full, err := os.ReadFile(filepath.Join(m.Config().OutputDir, "css", "palette.css"))
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{`[data-palette="ayu-dark"]`, `[data-palette="ayu-light"]`, "--palette-manifest:", "--test-override: 42px", `--aesthetic: "precision"`, "--color-border-strong:"} {
				if !strings.Contains(string(core), want) {
					t.Errorf("core missing %s", want)
				}
			}
			if strings.Contains(string(core), `nord-dark`) || !strings.Contains(string(full), `nord-dark`) {
				t.Error("optional palette colors and metadata must remain in the full stylesheet only")
			}
			for name, css := range map[string][]byte{"core": core, "full": full} {
				_, value, ok := strings.Cut(string(css), "--palette-manifest: '")
				if !ok {
					t.Fatalf("%s missing manifest", name)
				}
				value, _, _ = strings.Cut(value, "';")
				var entries []PaletteManifestEntry
				if err := json.Unmarshal([]byte(strings.ReplaceAll(value, "\\'", "'")), &entries); err != nil {
					t.Fatal(err)
				}
				if name == "core" && len(entries) != 2 {
					t.Fatalf("core must contain the configured pair, got %d palettes", len(entries))
				}
				if name == "full" && len(entries) <= 2 {
					t.Fatal("full stylesheet must retain the complete palette catalog")
				}
			}
			if len(core)*2 >= len(full) {
				t.Errorf("core %d bytes should be less than half of full %d", len(core), len(full))
			}
			hash := fmt.Sprintf("%x", sha256.Sum256(core))[:8]
			if hash != m.GetAssetHash("css/palette-core.css") {
				t.Fatal("registered core hash differs from the emitted bytes")
			}
			hashed, err := os.ReadFile(filepath.Join(m.Config().OutputDir, "css", "palette-core."+hash+".css"))
			if err != nil || !bytes.Equal(hashed, core) {
				t.Fatal("hashed core artifact must match the unversioned stylesheet")
			}
		})
	}
}

func TestPaletteCSSPlugin_Name(t *testing.T) {
	p := NewPaletteCSSPlugin()
	if got := p.Name(); got != "palette_css" {
		t.Errorf("Name() = %q, want %q", got, "palette_css")
	}
}

func TestPaletteCSSPlugin_GetPaletteConfig_WithThemeConfig(t *testing.T) {
	p := NewPaletteCSSPlugin()

	// Test with models.ThemeConfig (as stored by core.go)
	extra := map[string]interface{}{
		"theme": models.ThemeConfig{
			Name:         "default",
			Palette:      "catppuccin-mocha",
			PaletteLight: "catppuccin-latte",
			PaletteDark:  "catppuccin-mocha",
		},
	}

	palette, paletteLight, paletteDark, _ := p.getPaletteConfig(extra)

	if palette != "catppuccin-mocha" {
		t.Errorf("palette = %q, want %q", palette, "catppuccin-mocha")
	}
	if paletteLight != "catppuccin-latte" {
		t.Errorf("paletteLight = %q, want %q", paletteLight, "catppuccin-latte")
	}
	if paletteDark != "catppuccin-mocha" {
		t.Errorf("paletteDark = %q, want %q", paletteDark, "catppuccin-mocha")
	}
}

func TestPaletteCSSPlugin_GetPaletteConfig_WithMapInterface(t *testing.T) {
	p := NewPaletteCSSPlugin()

	// Test with map[string]interface{} (as from raw TOML)
	extra := map[string]interface{}{
		"theme": map[string]interface{}{
			"name":          "default",
			"palette":       "nord",
			"palette_light": "nord-light",
			"palette_dark":  "nord",
		},
	}

	palette, paletteLight, paletteDark, _ := p.getPaletteConfig(extra)

	if palette != "nord" {
		t.Errorf("palette = %q, want %q", palette, "nord")
	}
	if paletteLight != "nord-light" {
		t.Errorf("paletteLight = %q, want %q", paletteLight, "nord-light")
	}
	if paletteDark != "nord" {
		t.Errorf("paletteDark = %q, want %q", paletteDark, "nord")
	}
}

func TestPaletteCSSPlugin_Write_GeneratesCSS(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name:    "default",
				Palette: "catppuccin-mocha",
			},
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Check the file was created
	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)

	// Verify header comment
	if !strings.Contains(css, "CSS Custom Properties - Generated by markata-go") {
		t.Error("expected header comment in CSS")
	}

	// Verify it mentions catppuccin
	if !strings.Contains(css, "catppuccin") {
		t.Error("expected catppuccin in CSS")
	}

	// Verify it has color variables
	if !strings.Contains(css, "--color-primary") {
		t.Error("expected --color-primary in CSS")
	}
	if !strings.Contains(css, "--color-code-keyword") {
		t.Error("expected syntax highlighting code variables in CSS")
	}
	// Heading contours must track palette switches AND author-defined highlight
	// overrides without the old computed-style/inline-style updater.
	for _, alias := range []string{
		"--heading-mark-bg: var(--color-highlight);",
		"--heading-mark-fg: var(--color-highlight-text);",
	} {
		if !strings.Contains(css, alias) {
			t.Errorf("expected CSS-driven highlight alias %q", alias)
		}
	}
}

func TestPaletteCSSPlugin_Write_DefaultsToDarkWithoutSystemPreference(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name:         "default",
				Palette:      "catppuccin-latte",
				PaletteDark:  "catppuccin-mocha",
				PaletteLight: "catppuccin-latte",
			},
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)
	if !strings.Contains(css, ":root:not([data-theme=\"light\"])") {
		t.Error("expected dark mode to be the default fallback selector")
	}
	if !strings.Contains(css, "[data-theme=\"light\"]") {
		t.Error("expected explicit light mode selector")
	}
	if strings.Contains(css, "prefers-color-scheme") {
		t.Error("did not expect prefers-color-scheme media queries")
	}
}

func TestPaletteCSSPlugin_Write_MultiPaletteDefaultsToDarkWithoutSystemPreference(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": map[string]interface{}{
				"palette":       "catppuccin-latte",
				"palette_light": "catppuccin-latte",
				"palette_dark":  "catppuccin-mocha",
				"switcher": map[string]interface{}{
					"enabled": true,
				},
			},
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)
	if !strings.Contains(css, ":root:not([data-theme=\"light\"]):not([data-palette])") {
		t.Error("expected multi-palette default dark fallback selector")
	}
	if !strings.Contains(css, "[data-theme=\"light\"]:not([data-palette])") {
		t.Error("expected explicit light mode selector in multi-palette CSS")
	}
	if strings.Contains(css, "prefers-color-scheme") {
		t.Error("did not expect prefers-color-scheme media queries in multi-palette CSS")
	}
}

func TestPaletteCSSPlugin_Write_UsesLightFallbackModeWhenConfigured(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name:         "default",
				Palette:      "catppuccin-latte",
				PaletteDark:  "catppuccin-mocha",
				PaletteLight: "catppuccin-latte",
				FallbackMode: "light",
			},
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)
	if !strings.Contains(css, ":root:not([data-theme=\"dark\"])") {
		t.Error("expected light mode to be the default fallback selector")
	}
	if !strings.Contains(css, "[data-theme=\"dark\"]") {
		t.Error("expected explicit dark mode selector when fallback_mode=light")
	}
	if strings.Contains(css, "prefers-color-scheme") {
		t.Error("did not expect prefers-color-scheme media queries when fallback_mode=light")
	}
}

func TestPaletteCSSPlugin_Write_MultiPaletteUsesLightFallbackModeWhenConfigured(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": map[string]interface{}{
				"palette":       "catppuccin-latte",
				"palette_light": "catppuccin-latte",
				"palette_dark":  "catppuccin-mocha",
				"fallback_mode": "light",
				"switcher": map[string]interface{}{
					"enabled": true,
				},
			},
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)
	if !strings.Contains(css, ":root:not([data-theme=\"dark\"]):not([data-palette])") {
		t.Error("expected multi-palette light fallback selector when fallback_mode=light")
	}
	if !strings.Contains(css, "[data-theme=\"dark\"]:not([data-palette])") {
		t.Error("expected explicit dark mode selector in multi-palette CSS when fallback_mode=light")
	}
	if strings.Contains(css, "prefers-color-scheme") {
		t.Error("did not expect prefers-color-scheme media queries in multi-palette CSS when fallback_mode=light")
	}
}

func TestPaletteCSSPlugin_Write_AppliesThemeVariables(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name:    "default",
				Palette: "catppuccin-mocha",
				Variables: map[string]string{
					"--text-base":     "1.125rem",
					"--content-width": "70ch",
				},
			},
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)
	if !strings.Contains(css, "--text-base: 1.125rem;") {
		t.Error("expected theme variable override for --text-base")
	}
	if !strings.Contains(css, "--content-width: 70ch;") {
		t.Error("expected theme variable override for --content-width")
	}
}

func TestPaletteCSSPlugin_Write_AppliesThemeVariablesFromModelsConfig(t *testing.T) {
	tmpDir := t.TempDir()

	modelsConfig := &models.Config{}
	modelsConfig.Theme = models.ThemeConfig{
		Name:    "default",
		Palette: "catppuccin-mocha",
		Variables: map[string]string{
			"--text-base":     "1.5rem",
			"--content-width": "72ch",
		},
	}

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"models_config": modelsConfig,
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)
	if !strings.Contains(css, "--text-base: 1.5rem;") {
		t.Error("expected theme variable override from models_config for --text-base")
	}
	if !strings.Contains(css, "--content-width: 72ch;") {
		t.Error("expected theme variable override from models_config for --content-width")
	}
}

func TestPaletteCSSPlugin_Write_NoPaletteSkips(t *testing.T) {
	tmpDir := t.TempDir()

	p := NewPaletteCSSPlugin()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name: "default",
				// No palette set
			},
		},
	})

	if err := p.Write(m); err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// File should NOT be created when no palette is configured
	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	if _, err := os.Stat(cssPath); err == nil {
		t.Error("expected palette.css to NOT be created when no palette is configured")
	}
}

func TestPaletteCSSPlugin_Write_ContractOwnsActiveVariables(t *testing.T) {
	tmpDir := t.TempDir()
	modelsConfig := &models.Config{}
	modelsConfig.Theme = models.ThemeConfig{Palette: "nord", Aesthetic: "brutal"}
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{OutputDir: tmpDir, Extra: map[string]interface{}{"models_config": modelsConfig}})
	if err := NewPaletteCSSPlugin().Write(m); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(tmpDir, "css", "palette.css"))
	if err != nil {
		t.Fatal(err)
	}
	css := string(content)
	for _, want := range []string{"--bg: #242933", "--panel: #2e3440", "--ink: #d8dee9", "--accent: #88c0d0", "--color-background: #242933", "--color-surface: #2e3440", "--color-text: #d8dee9", "--color-primary: #88c0d0", "--aesthetic: \"brutal\"", "--theme-radius: 0px"} {
		if !strings.Contains(css, want) {
			t.Errorf("CSS missing %q", want)
		}
	}
}

func TestPaletteCSSPlugin_Write_ResolvesContractAliasBeforeLegacyLoader(t *testing.T) {
	tmpDir := t.TempDir()
	modelsConfig := &models.Config{}
	modelsConfig.Theme = models.ThemeConfig{Palette: "catppuccin"}
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{OutputDir: tmpDir, Extra: map[string]interface{}{"models_config": modelsConfig}})
	if err := NewPaletteCSSPlugin().Write(m); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(tmpDir, "css", "palette.css"))
	if err != nil {
		t.Fatal(err)
	}
	css := string(content)
	for _, want := range []string{"--bg: #1e1e2e", "--panel: #313244", "--ink: #cdd6f4", "--accent: #cba6f7"} {
		if !strings.Contains(css, want) {
			t.Errorf("CSS missing Catppuccin role %q", want)
		}
	}
}

// TestPaletteCSSPlugin_ServeMode_ConfigChange simulates the serve mode scenario
// where palette config changes should be detected and new CSS should be generated.
func TestPaletteCSSPlugin_ServeMode_ConfigChange(t *testing.T) {
	tmpDir := t.TempDir()

	// First build with catppuccin-mocha
	p1 := NewPaletteCSSPlugin()
	m1 := lifecycle.NewManager()
	m1.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name:    "default",
				Palette: "catppuccin-mocha",
			},
		},
	})

	if err := p1.Write(m1); err != nil {
		t.Fatalf("First Write error: %v", err)
	}

	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content1, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after first build: %v", err)
	}

	css1 := string(content1)
	if !strings.Contains(css1, "catppuccin-mocha") {
		t.Errorf("First build should contain catppuccin-mocha, got:\n%s", css1[:min(500, len(css1))])
	}

	// Second build with dracula - simulates user changing config
	p2 := NewPaletteCSSPlugin()
	m2 := lifecycle.NewManager()
	m2.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name:    "default",
				Palette: "dracula",
			},
		},
	})

	if err := p2.Write(m2); err != nil {
		t.Fatalf("Second Write error: %v", err)
	}

	content2, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after second build: %v", err)
	}

	css2 := string(content2)
	if !strings.Contains(css2, "dracula") {
		t.Errorf("Second build should contain dracula, got:\n%s", css2[:min(500, len(css2))])
	}

	// Verify the CSS actually changed
	if css1 == css2 {
		t.Error("CSS should have changed between builds with different palettes")
	}
}

func TestPaletteCSSPlugin_Priority(t *testing.T) {
	p := NewPaletteCSSPlugin()

	// Should have default priority for write stage (runs after static_assets)
	if got := p.Priority(lifecycle.StageWrite); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

// TestPaletteCSSPlugin_FullBuildFlow simulates the full build flow with
// static_assets running before palette_css to verify palette CSS overwrites defaults.
func TestPaletteCSSPlugin_FullBuildFlow(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup: Create a manager with both plugins
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name:    "default",
				Palette: "catppuccin-mocha",
			},
		},
	})

	// Step 1: Run static_assets (writes default palette.css)
	staticPlugin := NewStaticAssetsPlugin()
	if err := staticPlugin.Write(m); err != nil {
		t.Fatalf("StaticAssetsPlugin.Write error: %v", err)
	}

	// Step 2: Run palette_css (should overwrite with catppuccin)
	palettePlugin := NewPaletteCSSPlugin()
	m.SetAssetHash("css/palette.css", "deadbeef")
	if err := palettePlugin.Write(m); err != nil {
		t.Fatalf("PaletteCSSPlugin.Write error: %v", err)
	}

	// Verify: CSS should have catppuccin colors, not defaults
	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css: %v", err)
	}

	css := string(content)

	// Should have catppuccin-mocha in the header
	if !strings.Contains(css, "catppuccin-mocha") {
		t.Error("expected catppuccin-mocha in CSS header")
	}

	// Should have "Generated by markata-go" header (from palette_css, not default)
	if !strings.Contains(css, "CSS Custom Properties - Generated by markata-go") {
		t.Error("expected Generated by markata-go in CSS header")
	}

	hashedPath := filepath.Join(tmpDir, "css", "palette.deadbeef.css")
	if _, err := os.Stat(hashedPath); err != nil {
		t.Fatalf("expected hashed palette.css to exist: %v", err)
	}
}

// TestPaletteCSSPlugin_FullBuildFlow_NoPalette tests that when no palette is configured,
// the default CSS from static_assets is preserved.
func TestPaletteCSSPlugin_FullBuildFlow_NoPalette(t *testing.T) {
	tmpDir := t.TempDir()

	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{
		OutputDir: tmpDir,
		Extra: map[string]interface{}{
			"theme": models.ThemeConfig{
				Name: "default",
				// No Palette set
			},
		},
	})

	// Step 1: Run static_assets
	staticPlugin := NewStaticAssetsPlugin()
	if err := staticPlugin.Write(m); err != nil {
		t.Fatalf("StaticAssetsPlugin.Write error: %v", err)
	}

	// Check if default CSS exists
	cssPath := filepath.Join(tmpDir, "css", "palette.css")
	_, err := os.Stat(cssPath)
	defaultExists := err == nil

	// Step 2: Run palette_css (should skip since no palette)
	palettePlugin := NewPaletteCSSPlugin()
	if err := palettePlugin.Write(m); err != nil {
		t.Fatalf("PaletteCSSPlugin.Write error: %v", err)
	}

	// If default existed, verify it's still there and unmodified (or modified by palette_css is fine too)
	if defaultExists {
		content, err := os.ReadFile(cssPath)
		if err != nil {
			t.Fatalf("failed to read palette.css: %v", err)
		}
		// Just verify it's not empty
		if len(content) == 0 {
			t.Error("palette.css should not be empty")
		}
	}
}
