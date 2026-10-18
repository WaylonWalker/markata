package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/fontpacks"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/themes"
)

var paletteManifestPattern = regexp.MustCompile(`--palette-manifest: '([^\n]*)';`)

func writePickerPaletteCSS(t *testing.T, theme models.ThemeConfig) string {
	t.Helper()
	tmpDir := t.TempDir()
	modelsConfig := &models.Config{}
	modelsConfig.Theme = theme
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{OutputDir: tmpDir, Extra: map[string]interface{}{"models_config": modelsConfig}})
	if err := NewPaletteCSSPlugin().Write(m); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(tmpDir, "css", "palette.css"))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestThemePicker_EnabledByDefault(t *testing.T) {
	var cfg models.ThemeSwitcherConfig
	if !cfg.IsEnabled() {
		t.Fatal("unset switcher config should be enabled by default")
	}
	defaults := models.NewThemeSwitcherConfig()
	if !defaults.IsEnabled() {
		t.Fatal("NewThemeSwitcherConfig should be enabled")
	}
	if !NewPaletteCSSPlugin().isSwitcherEnabled(map[string]interface{}{}) {
		t.Fatal("palette CSS should include every palette when the switcher is unset")
	}
}

func TestThemePicker_DisabledWritesSinglePalette(t *testing.T) {
	off := false
	css := writePickerPaletteCSS(t, models.ThemeConfig{Palette: "nord", Switcher: models.ThemeSwitcherConfig{Enabled: &off}})
	if strings.Contains(css, "--palette-manifest") {
		t.Fatal("disabled switcher should not emit a palette manifest")
	}
}

func TestThemePicker_ManifestIsCompleteAndConsistent(t *testing.T) {
	css := writePickerPaletteCSS(t, models.ThemeConfig{Palette: "catppuccin-mocha"})
	match := paletteManifestPattern.FindStringSubmatch(css)
	if match == nil {
		t.Fatal("palette manifest not found")
	}
	var manifest []PaletteManifestEntry
	if err := json.Unmarshal([]byte(strings.ReplaceAll(match[1], `\'`, "'")), &manifest); err != nil {
		t.Fatalf("manifest is not valid JSON: %v", err)
	}
	if len(manifest) < 20 {
		t.Fatalf("manifest has %d entries, want the full theme catalog", len(manifest))
	}

	byName := make(map[string]PaletteManifestEntry, len(manifest))
	for _, entry := range manifest {
		if _, dup := byName[entry.Name]; dup {
			t.Errorf("duplicate manifest entry %q", entry.Name)
		}
		byName[entry.Name] = entry
		if entry.Variant != "light" && entry.Variant != "dark" {
			t.Errorf("%s: variant %q, want light or dark", entry.Name, entry.Variant)
		}
		if strings.Contains(entry.Name, "-light-dark") || strings.Contains(entry.Name, "-dark-light") {
			t.Errorf("derived-of-derived palette %q", entry.Name)
		}
		if strings.HasSuffix(entry.DisplayName, " Light") || strings.HasSuffix(entry.DisplayName, " Dark") {
			t.Errorf("%s: label %q should drop the mode suffix", entry.Name, entry.DisplayName)
		}
		if !strings.Contains(css, `[data-palette="`+entry.Name+`"]`) {
			t.Errorf("%s: no scoped CSS block for manifest entry", entry.Name)
		}
	}
	for _, entry := range manifest {
		if entry.Counterpart == "" {
			continue
		}
		other, ok := byName[entry.Counterpart]
		if !ok {
			t.Errorf("%s: counterpart %q missing from manifest", entry.Name, entry.Counterpart)
			continue
		}
		if other.Variant == entry.Variant {
			t.Errorf("%s: counterpart %q has the same variant", entry.Name, entry.Counterpart)
		}
	}
	for _, name := range []string{"catppuccin-mocha", "catppuccin-latte"} {
		if _, ok := byName[name]; !ok {
			t.Errorf("site default %q missing from manifest", name)
		}
	}
	if got := byName["gruvbox-dark"].Counterpart; got != "gruvbox-light" {
		t.Errorf("gruvbox-dark counterpart = %q, want gruvbox-light", got)
	}
}

func TestThemePicker_ContractFallbackYieldsToPickedPalette(t *testing.T) {
	css := writePickerPaletteCSS(t, models.ThemeConfig{Palette: "catppuccin-mocha"})
	if strings.Contains(css, ":root:not([data-theme=\"light\"]), [data-theme=\"dark\"] {") {
		t.Error("contract dark fallback must not outrank [data-palette] blocks")
	}
	if !strings.Contains(css, `:where(:root:not([data-theme="light"]), [data-theme="dark"]) {`) {
		t.Error("contract dark fallback should use zero specificity")
	}
}

func TestPaletteFamilyLabel(t *testing.T) {
	tests := map[string]string{
		"Ayu Dark":         "Ayu",
		"Gruvbox Light":    "Gruvbox",
		"Catppuccin Mocha": "Catppuccin Mocha",
		"Dark":             "Dark",
	}
	for in, want := range tests {
		if got := paletteFamilyLabel(in); got != want {
			t.Errorf("paletteFamilyLabel(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPaletteDiscover_Deterministic(t *testing.T) {
	first, err := palettes.NewLoader().Discover()
	if err != nil {
		t.Fatal(err)
	}
	second, err := palettes.NewLoader().Discover()
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != len(second) {
		t.Fatalf("Discover() returned %d then %d palettes", len(first), len(second))
	}
	for i := range first {
		if first[i].Name != second[i].Name {
			t.Fatalf("Discover() order differs at %d: %q vs %q", i, first[i].Name, second[i].Name)
		}
	}
}

func TestLiveMotifCSS(t *testing.T) {
	theme := models.NewThemeConfig()
	theme.Palette = "catppuccin-mocha"
	theme.Motif.Kind = renderingMotifBlockW
	css := liveMotifCSS(theme, `url("/assets/motif-block-w-v1.svg")`, 0.01)
	for _, want := range []string{
		`html[data-palette]:not([data-palette="catppuccin-mocha"]) body::after`,
		`mask-image: url("/assets/motif-block-w-v1.svg")`,
		"color-mix(in srgb, var(--color-text) 1.00%, var(--color-background))",
	} {
		if !strings.Contains(css, want) {
			t.Errorf("live motif CSS missing %q\n%s", want, css)
		}
	}
	theme.Motif.Kind = renderingMotifOff
	if got := liveMotifCSS(theme, `url("/x.svg")`, 0.01); got != "" {
		t.Errorf("motif off should emit nothing, got %q", got)
	}
	theme.Motif.Kind = renderingMotifBlockW
	if got := liveMotifCSS(theme, renderingTextureNone, 0.01); got != "" {
		t.Errorf("no motif image should emit nothing, got %q", got)
	}
}

func TestAestheticManifest_UsesNormalizedIDs(t *testing.T) {
	if got := aestheticID("Minimal"); got != "minimal" {
		t.Errorf("aestheticID(Minimal) = %q", got)
	}
	if got := aestheticID("Soft Glow_X"); got != "soft-glow-x" {
		t.Errorf("aestheticID = %q", got)
	}
}

func TestAestheticCSSPlugin_DefaultMotifOffKeepsTexture(t *testing.T) {
	theme := models.NewThemeConfig()
	if theme.Motif.Kind != renderingMotifOff {
		t.Fatalf("default motif kind = %q, want off", theme.Motif.Kind)
	}
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
	bundle, err := compileMotifBundle(config)
	if err != nil {
		t.Fatalf("compileMotifBundle() error = %v", err)
	}
	if _, ok := bundle.Assets["assets/surface-screenprint-v1.svg"]; !ok {
		t.Fatalf("texture asset missing with motif off: %v", bundle.Assets)
	}
	if _, ok := bundle.Assets["assets/motif-block-w-v1.svg"]; ok {
		t.Fatal("motif asset emitted while motif is off")
	}
	css := NewAestheticCSSPlugin().generatePresentationCSS(config)
	if !strings.Contains(css, "--theme-motif-image: none") {
		t.Fatalf("motif image should be none when off:\n%s", css)
	}
}

func TestFontpackManifestCSS_ListsPacksWithFamilies(t *testing.T) {
	source, err := fontpacks.BuiltinSource()
	if err != nil {
		t.Fatalf("BuiltinSource() error = %v", err)
	}
	css := fontpackManifestCSS(source.Catalog, "brush", source.Catalog.FontPacks)
	for _, want := range []string{
		`--fontpack-default:"brush"`,
		`{"name":"brush","displayName":"Brush","heading":"Knewave","body":"Space Grotesk","code":"DM Mono","headingFont":"\'Knewave\', cursive"`,
		`"name":"system"`,
	} {
		if !strings.Contains(css, want) {
			t.Errorf("fontpack manifest missing %q\n%s", want, css)
		}
	}
	if strings.Contains(css, `\"`) {
		t.Error("manifest must not contain escaped double quotes")
	}
	if strings.Count(css, `"name":`) != len(source.Catalog.FontPacks) {
		t.Errorf("manifest entries = %d, want %d", strings.Count(css, `"name":`), len(source.Catalog.FontPacks))
	}
}

func TestAestheticSurfaceCSS_DistinctPerAesthetic(t *testing.T) {
	// The generated stylesheet is aestheticSurfaceCSS followed by the theme's
	// aesthetic.css, which owns the per-aesthetic radius and surface tokens.
	themeCSS, err := themes.ReadStatic("css/aesthetic.css")
	if err != nil {
		t.Fatalf("read aesthetic.css: %v", err)
	}
	combined := aestheticSurfaceCSS + strings.ReplaceAll(string(themeCSS), "\r\n", "\n")
	for _, name := range []string{"balanced", "elevated", "precision", "brutal", "minimal"} {
		if !strings.Contains(combined, `[data-aesthetic="`+name+`"] {`+"\n  --radius-sm:") {
			t.Errorf("aesthetic %q has no radius tokens", name)
		}
	}
	for _, name := range []string{"precision", "elevated", "minimal"} {
		if !strings.Contains(combined, `html[data-aesthetic][data-aesthetic="`+name+`"]`) {
			t.Errorf("aesthetic %q has no signature surface rules", name)
		}
	}
	if !strings.Contains(aestheticSurfaceCSS, "--brutal-ink:") {
		t.Error("brutal ink token missing from surface CSS")
	}
	if strings.Contains(combined, `html[data-aesthetic][data-aesthetic="balanced"]`) {
		t.Error("balanced is the default look and must not get signature rules")
	}
}
