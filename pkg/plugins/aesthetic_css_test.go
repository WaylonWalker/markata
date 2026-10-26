package plugins

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"
	"github.com/WaylonWalker/markata/pkg/renderingrecipe"
	"github.com/WaylonWalker/markata/pkg/themes"
)

func TestAestheticCSSPlugin_VariablesAreCanonical(t *testing.T) {
	m := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.OutputDir = t.TempDir()
	theme := models.NewThemeConfig()
	theme.Aesthetic = "precision"
	theme.Variables = map[string]string{
		"--text-base":     "1.2rem",
		"--content-width": "clamp(56ch, 60vw, 72ch)",
		"--invalid;name":  "ignored",
	}
	config.Extra = map[string]interface{}{
		"theme":         theme,
		"models_config": &models.Config{Theme: theme},
	}
	m.SetConfig(config)
	plugin := NewAestheticCSSPlugin()
	var previous []byte
	for i := 0; i < 20; i++ {
		if err := plugin.Configure(m); err != nil {
			t.Fatal(err)
		}
		if err := plugin.Write(m); err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(filepath.Join(config.OutputDir, "css", "aesthetic.css"))
		if err != nil {
			t.Fatal(err)
		}
		if i > 0 && !bytes.Equal(previous, got) {
			t.Fatal("identical variables produced different CSS")
		}
		wantHash := fmt.Sprintf("%x", sha256.Sum256(got))[:8]
		if hash := m.GetAssetHash("css/aesthetic.css"); hash != wantHash {
			t.Fatalf("configured hash %q does not match generated CSS %q", hash, wantHash)
		}
		css := string(got)
		if strings.Index(css, "--content-width: clamp") > strings.Index(css, "--text-base: 1.2rem") {
			t.Fatal("theme variables are not sorted by property name")
		}
		if strings.Contains(css, "--invalid;name") {
			t.Fatal("invalid CSS variable name was emitted")
		}
		previous = got
	}
}

func TestAestheticCSSPlugin_BundlesConsumptionRules(t *testing.T) {
	m := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.OutputDir = t.TempDir()
	config.Extra = map[string]interface{}{"theme": models.ThemeConfig{Aesthetic: "brutal"}}
	m.SetConfig(config)
	plugin := NewAestheticCSSPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Write(m); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join(config.OutputDir, "css", "aesthetic.css"))
	if err != nil {
		t.Fatal(err)
	}
	surfaces, err := themes.ReadStatic("css/aesthetic.css")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(got), string(surfaces)) {
		t.Fatal("generated aesthetic stylesheet must end with default theme's surface rules")
	}
	hash := m.GetAssetHash("css/aesthetic.css")
	hashed, err := os.ReadFile(filepath.Join(config.OutputDir, "css", "aesthetic."+hash+".css"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(hashed, got) {
		t.Fatal("hashed asset differs from generated aesthetic CSS")
	}
}

func TestAestheticCSSPlugin_Configure(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	m := lifecycle.NewManager()

	// Create config
	cfg := lifecycle.NewConfig()

	// Setup extra
	b := true
	themeCfg := models.NewThemeConfig()
	themeCfg.Aesthetic = "brutal"
	themeCfg.Switcher.Enabled = &b

	cfg.Extra = map[string]interface{}{
		"theme": themeCfg,
	}
	m.SetConfig(cfg)

	err := plugin.Configure(m)
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}

	hash := m.GetAssetHash("css/aesthetic.css")
	if hash == "" {
		t.Error("Expected hash to be set for css/aesthetic.css")
	}
}

func TestAestheticCSSPlugin_PresentationUsesContractTextureSemantics(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Texture.Kind = "screenprint"
	theme.Texture.ColorMix = 0.35
	theme.Texture.Scope = "all"
	theme.HeadingTexture.Kind = "inherit"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, `--theme-texture-opacity: 1.000;`) {
		t.Errorf("compiled texture paint should remain opaque at nonzero mix: %s", css)
	}
	if !strings.Contains(css, `--theme-heading-texture-kind: "screenprint";`) {
		t.Errorf("inherit heading texture was not resolved: %s", css)
	}
	if strings.Contains(css, `opacity: var(--theme-texture-color-mix)`) {
		t.Error("surface opacity must not use public color_mix directly")
	}
}

func TestCompileMotifBundle_NonCanonicalFontpacksUseLegacyProjection(t *testing.T) {
	for _, fontpack := range []string{"field-notebook", "editorial"} {
		t.Run(fontpack, func(t *testing.T) {
			config := lifecycle.NewConfig()
			theme := models.NewThemeConfig()
			theme.Fontpack = fontpack
			config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
			bundle, err := compileMotifBundle(config)
			if err != nil {
				t.Fatalf("legacy fontpack must not fail recipe compilation: %v", err)
			}
			if len(bundle.Assets) != 0 {
				t.Fatalf("legacy fontpack received canonical assets: %v", bundle.Assets)
			}
		})
	}
}

func TestCompileMotifBundle_CanonicalBrushStillCompiles(t *testing.T) {
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Fontpack = "brush"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
	bundle, err := compileMotifBundle(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Assets) == 0 {
		t.Fatal("canonical brush did not compile assets")
	}
}

func TestAestheticCSSPlugin_PresentationPropagatesCanonicalContentWidth(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Variables = map[string]string{"--content-width": "68ch"}
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, "--content-width: 68ch;") {
		t.Fatalf("canonical content width was not emitted in the active presentation CSS: %s", css)
	}
}

func TestAestheticCSSPlugin_PresentationEmitsMotifLengthsAsCSSLengths(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Motif.Size = "71px"
	theme.Motif.Gap = "11px"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, "--theme-motif-size: 71px;") || !strings.Contains(css, "--theme-motif-gap: 11px;") {
		t.Fatalf("motif lengths must be emitted as unquoted CSS tokens: %s", css)
	}
	if strings.Contains(css, `--theme-motif-size: "71px"`) || strings.Contains(css, `--theme-motif-gap: "11px"`) {
		t.Fatal("motif lengths must not be quoted")
	}
}

func TestAestheticCSSPlugin_CanonicalMotifUsesCompiledAsset(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Motif.Kind = "block-w"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, `--theme-motif-image: url("/assets/motif-block-w-v1.svg")`) || strings.Contains(css, "data:image/svg+xml;base64,") {
		t.Fatalf("canonical motif must attach compiler asset without data URI: %s", css)
	}
	bundle, err := compileMotifBundle(config)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(string(bundle.Assets["assets/motif-block-w-v1.svg"]), `data-index="`) != 160 {
		t.Fatal("compiled canonical motif must contain 160 marks")
	}
}

func TestAestheticCSSPlugin_EmptyTextureKeepsSupportedHeadingAndMotif(t *testing.T) {
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Texture.Kind = ""
	theme.HeadingTexture.Kind = "splatter"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	bundle, err := compileMotifBundle(config)
	if err != nil {
		t.Fatalf("empty texture should be a no-op: %v", err)
	}
	for _, pass := range bundle.Manifest.Passes {
		if pass.ID == "surface-texture" {
			t.Fatal("empty texture must not emit a surface texture pass")
		}
	}
	if !containsPass(bundle.Manifest.Passes, "heading-wear") || !containsPass(bundle.Manifest.Passes, "motif-under") || !containsPass(bundle.Manifest.Passes, "motif-over") {
		t.Fatalf("heading and motif passes should still compile: %+v", bundle.Manifest.Passes)
	}
}

func TestAestheticCSSPlugin_InvalidTextureStillRejected(t *testing.T) {
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Texture.Kind = "unsupported"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	if _, err := compileMotifBundle(config); err == nil {
		t.Fatal("unsupported texture kind must be rejected")
	}
}

func containsPass(passes []renderingrecipe.Pass, id string) bool {
	for i := range passes {
		pass := &passes[i]
		if pass.ID == id {
			return true
		}
	}
	return false
}

func TestAestheticCSSPlugin_PresentationRendersOverMotifPass(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Motif.Layer = "over"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, "body::after { content: ''; pointer-events: none;") {
		t.Fatalf("over motif pass must create a non-interactive pseudo-element: %s", css)
	}
	if !strings.Contains(css, "background-image: var(--theme-motif-over-image);") {
		t.Fatalf("over motif pass must paint its configured image: %s", css)
	}
}

func TestAestheticCSSPlugin_PresentationUsesNormalizedSeparationAndLayers(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Texture.ColorMix = 2
	theme.HeadingTexture.ColorMix = 0.5
	theme.Motif.ColorMix = 0.5
	theme.Motif.Layer = "under"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, `--theme-texture-color-mix: 1.000;`) || !strings.Contains(css, `--theme-heading-texture-color-mix: 0.500;`) {
		t.Fatalf("color mixes were not normalized at endpoints: %s", css)
	}
	if !strings.Contains(css, `--theme-motif-color-mix: 0.500;`) || !strings.Contains(css, `--theme-motif-z: -2;`) {
		t.Fatalf("intermediate motif mix or under layer missing: %s", css)
	}
	if !strings.Contains(css, "mask-image: var(--theme-heading-texture-mask)") || !strings.Contains(css, "/assets/heading-splatter-v1.svg") {
		t.Fatal("heading wear must use a stable SVG mask")
	}
	if strings.Contains(css, "opacity: var(--theme-motif-color-mix)") {
		t.Fatal("motif color_mix must not be applied twice through opacity")
	}
}

func TestAestheticCSSPlugin_HeadingColorMixChangesResolvedWearPaint(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	colors := make([]string, 0, 3)
	for _, mix := range []float64{0, 0.29, 1} {
		config := lifecycle.NewConfig()
		theme := models.NewThemeConfig()
		theme.HeadingTexture.ColorMix = mix
		config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
		css := plugin.generatePresentationCSS(config)
		marker := "--theme-heading-texture-color: "
		start := strings.Index(css, marker)
		if start < 0 {
			t.Fatalf("missing resolved heading paint for mix %v", mix)
		}
		value := css[start+len(marker):]
		value = value[:strings.IndexByte(value, ';')]
		colors = append(colors, value)
	}
	if colors[0] == colors[1] || colors[1] == colors[2] || colors[0] == colors[2] {
		t.Fatalf("heading endpoint paints must progress, got %v", colors)
	}
}

func TestAestheticCSSPlugin_MotifLayersHaveStableZSemantics(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	for layer, want := range map[string]string{"under": "-2", "sandwich": "0", "over": "1"} {
		config := lifecycle.NewConfig()
		theme := models.NewThemeConfig()
		theme.Motif.Layer = layer
		config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
		css := plugin.generatePresentationCSS(config)
		if !strings.Contains(css, `--theme-motif-z: `+want+`;`) {
			t.Errorf("layer %q did not emit z-index %s", layer, want)
		}
		if layer == "sandwich" && !strings.Contains(css, "body { background-image: var(--theme-motif-under-image)") {
			t.Errorf("sandwich layer must have an under motif pass: %s", css)
		}
	}
}

func TestAestheticCSSPlugin_PresentationKeepsTextureAndMotifPassesIndependent(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	for _, layer := range []string{"under", "sandwich"} {
		config := lifecycle.NewConfig()
		theme := models.NewThemeConfig()
		theme.Texture.Kind = "screenprint"
		theme.Texture.ColorMix = 0
		theme.Motif.Layer = layer
		theme.Motif.Kind = "block-w"
		config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
		css := plugin.generatePresentationCSS(config)
		if !strings.Contains(css, `--theme-texture-opacity: 0.000;`) {
			t.Errorf("%s texture pass should be transparent at mix zero: %s", layer, css)
		}
		if !strings.Contains(css, "body { background-image: var(--theme-motif-under-image)") || !strings.Contains(css, `--theme-motif-under-image: url("/assets/motif-block-w-v1.svg")`) {
			t.Errorf("%s motif pass must remain independent of texture opacity: %s", layer, css)
		}
	}
}

func TestAestheticCSSPlugin_PresentationMotifMixEndpoints(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	for mix, want := range map[float64]string{
		0: `color-mix(in srgb, var(--color-background, #fff) 100.0%`,
		1: `color-mix(in srgb, var(--color-background, #fff) 0.0%`,
	} {
		config := lifecycle.NewConfig()
		theme := models.NewThemeConfig()
		theme.Texture.Kind = "screenprint"
		theme.Texture.ColorMix = 0.5
		theme.Motif.ColorMix = mix
		theme.Motif.Layer = "over"
		config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
		css := plugin.generatePresentationCSS(config)
		_ = want // Endpoint paint is now resolved to an opaque concrete SVG color.
		if !strings.Contains(css, `--theme-texture-opacity: 1.000;`) {
			t.Errorf("texture paint should remain opaque at nonzero mix: %s", css)
		}
		if !strings.Contains(css, `background-size: var(--theme-motif-field-size, calc(16 * (var(--theme-motif-size) + var(--theme-motif-gap)))) auto`) {
			t.Errorf("motif size and opacity must be independent of color mix: %s", css)
		}
	}
}

func TestAestheticCSSPlugin_LetterMotifDataURLIsEncoded(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Motif.Kind = "letter"
	theme.Motif.Glyph = "W"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, `data:image/svg+xml;base64,`) {
		t.Fatalf("letter motif SVG must be encoded: %s", css)
	}
}

func TestAestheticCSSPlugin_PresentationValidatesMotifURLSchemes(t *testing.T) {
	contract, err := renderingcontract.Load()
	if err != nil {
		t.Fatalf("load rendering contract: %v", err)
	}
	plugin := NewAestheticCSSPlugin()
	for _, scheme := range contract.MotifURLSchemes {
		config := lifecycle.NewConfig()
		theme := models.NewThemeConfig()
		theme.Motif.Kind = "letter"
		theme.Motif.URL = scheme + "example"
		config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
		css := plugin.generatePresentationCSS(config)
		if !strings.Contains(css, `--theme-motif-mask: none;`) || !strings.Contains(css, `--theme-motif-image: url("data:image/svg+xml;base64,`) {
			t.Errorf("contract scheme %q was not normalized as a recolorable mask: %s", scheme, css)
		}
	}

	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Motif.URL = "javascript:alert(1)"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}
	if css := plugin.generatePresentationCSS(config); strings.Contains(css, "javascript:alert") {
		t.Fatal("unsupported motif URL scheme must fall back to the generated field")
	}
}

func TestAestheticCSSPlugin_PresentationQuietScopeProtectsReadingSurfaces(t *testing.T) {
	plugin := NewAestheticCSSPlugin()
	config := lifecycle.NewConfig()
	theme := models.NewThemeConfig()
	theme.Texture.Kind = "screenprint"
	theme.Texture.Scope = "quiet"
	theme.HeadingTexture.Kind = "inherit"
	config.Extra = map[string]interface{}{"models_config": &models.Config{Theme: theme}}

	css := plugin.generatePresentationCSS(config)
	if !strings.Contains(css, `--theme-texture-kind: "screenprint";`) {
		t.Error("quiet scope must preserve the configured framing texture")
	}
	if !strings.Contains(css, `--theme-heading-texture-kind: "screenprint";`) {
		t.Error("quiet scope must preserve the heading texture kind")
	}
	if !strings.Contains(css, `[data-theme-texture-scope=quiet] article`) {
		t.Error("quiet scope must protect article reading surfaces")
	}
}
