package plugins

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/aesthetic"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"
	"github.com/WaylonWalker/markata/pkg/renderingrecipe"
	"github.com/WaylonWalker/markata/pkg/templates"
	"github.com/WaylonWalker/markata/pkg/themes"
)

// AestheticCSSPlugin generates CSS variables from the configured aesthetic.
// It runs during the Write stage and creates/overwrites css/aesthetic.css
// with the aesthetic's CSS custom properties.
type AestheticCSSPlugin struct{}

const (
	renderingMotifOff           = "off"
	renderingMotifBlockW        = "block-w"
	renderingLayerUnder         = "under"
	renderingLayerOver          = "over"
	renderingLayerSandwich      = "sandwich"
	renderingFontpackBrush      = "brush"
	renderingTextureNone        = "none"
	renderingTextureScreenprint = "screenprint"
	renderingHeadingSplatter    = "splatter"
	renderingScopeAll           = "all"
	renderingHeadingInherit     = "inherit"
	renderingMotifLetter        = "letter"
	renderingMotifAccent        = "accent"
	renderingMotifMuted         = "muted"
	renderingMotifShadow        = "shadow"
)

// NewAestheticCSSPlugin creates a new AestheticCSSPlugin.
func NewAestheticCSSPlugin() *AestheticCSSPlugin {
	return &AestheticCSSPlugin{}
}

// Name returns the unique name of the plugin.
func (p *AestheticCSSPlugin) Name() string {
	return "aesthetic_css"
}

// Configure generates CSS from the configured aesthetic and registers its hash
func (p *AestheticCSSPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	aestheticName := p.getAestheticConfig(config.Extra)
	if aestheticName == "" {
		return nil
	}

	loader := aesthetic.NewLoader()
	switcherEnabled := p.isSwitcherEnabled(config.Extra)

	var css string
	if switcherEnabled {
		css = p.generateMultiAestheticCSS(loader, config.Extra, aestheticName)
	} else {
		css = p.generateSingleAestheticCSS(loader, aestheticName)
	}
	css += p.generatePresentationCSS(config) + aestheticSurfaceCSS
	surfaceCSS, err := themes.ReadStatic("css/aesthetic.css")
	if err != nil {
		return fmt.Errorf("reading default aesthetic surfaces: %w", err)
	}
	css += "\n" + string(surfaceCSS)
	if bundle, err := compileMotifBundle(config); err != nil {
		return err
	} else if len(bundle.Assets) != 0 {
		hashes := make(map[string]string)
		for path, data := range bundle.Assets {
			hashes[path] = renderingrecipe.AssetHash(data)
			m.SetAssetHash(path, hashes[path])
		}
		templates.SetAssetHashes(hashes)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(css)))[:8]

	m.SetAssetHash("css/aesthetic.css", hash)
	templates.SetAssetHashes(map[string]string{"css/aesthetic.css": hash})

	log.Printf("[aesthetic_css] Registered hash %s for aesthetic.css", hash)

	return nil
}

// Write generates CSS from the configured aesthetic and writes it to the output directory.
func (p *AestheticCSSPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir
	if hashedAssetReadyInFastMode(m, "css/aesthetic.css") {
		return nil
	}

	aestheticName := p.getAestheticConfig(config.Extra)
	if aestheticName == "" {
		return nil
	}

	log.Printf("[aesthetic_css] Generating CSS for aesthetic: %s", aestheticName)

	switcherEnabled := p.isSwitcherEnabled(config.Extra)
	loader := aesthetic.NewLoader()

	var css string
	if switcherEnabled {
		css = p.generateMultiAestheticCSS(loader, config.Extra, aestheticName)
	} else {
		css = p.generateSingleAestheticCSS(loader, aestheticName)
	}
	css += p.generatePresentationCSS(config) + aestheticSurfaceCSS
	surfaceCSS, err := themes.ReadStatic("css/aesthetic.css")
	if err != nil {
		return fmt.Errorf("reading default aesthetic surfaces: %w", err)
	}
	css += "\n" + string(surfaceCSS)

	cssDir := filepath.Join(outputDir, "css")
	cssPath := filepath.Join(cssDir, "aesthetic.css")
	if bundle, err := compileMotifBundle(config); err != nil {
		return err
	} else if len(bundle.Assets) != 0 {
		for path, data := range bundle.Assets {
			assetPath := filepath.Join(outputDir, filepath.FromSlash(path))
			if err := os.MkdirAll(filepath.Dir(assetPath), 0o755); err != nil {
				return fmt.Errorf("creating compiled asset directory: %w", err)
			}
			//nolint:gosec // Generated static assets must be readable by the web server.
			if err := os.WriteFile(assetPath, data, 0o644); err != nil {
				return fmt.Errorf("writing compiled asset: %w", err)
			}
		}
	}
	if existing, err := os.ReadFile(cssPath); err == nil {
		if bytes.Equal(existing, []byte(css)) {
			return nil
		}
	}

	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		return fmt.Errorf("creating css directory: %w", err)
	}
	//nolint:gosec // Generated CSS must be readable by the web server.
	if err := os.WriteFile(cssPath, []byte(css), 0o644); err != nil {
		return fmt.Errorf("writing aesthetic CSS: %w", err)
	}

	if hash := m.GetAssetHash("css/aesthetic.css"); hash != "" {
		base := strings.TrimSuffix(filepath.Base(cssPath), filepath.Ext(cssPath))
		hashedPath := filepath.Join(cssDir, fmt.Sprintf("%s.%s.css", base, hash))
		//nolint:gosec // Generated CSS must be readable by the web server.
		if err := os.WriteFile(hashedPath, []byte(css), 0o644); err != nil {
			return fmt.Errorf("writing hashed aesthetic CSS: %w", err)
		}
	}

	return nil
}

func (p *AestheticCSSPlugin) generatePresentationCSS(config *lifecycle.Config) string {
	body := p.generatePresentationCSSBody(config)
	if config != nil {
		if configured, ok := config.Extra["models_config"].(*models.Config); ok && configured.Theme.Motif.URL != "" {
			contract, err := renderingcontract.Load()
			if err != nil {
				return body
			}
			if !validMotifURL(contract, configured.Theme.Motif.URL) || strings.ContainsAny(configured.Theme.Motif.URL, "\"'()\n\r\t") {
				body = strings.ReplaceAll(body, configured.Theme.Motif.URL, "")
			}
		}
	}
	body = unquoteCSSLengthProperty(body, "--theme-motif-size")
	body = unquoteCSSLengthProperty(body, "--theme-motif-gap")
	return body + `
@layer tokens {
  /* Motif color_mix selects the color. Motif paint stays opaque. */
	  body { background-size: var(--theme-motif-field-size, calc(var(--theme-motif-size) + var(--theme-motif-gap))) auto; }
	  body::after { background-size: var(--theme-motif-field-size, calc(var(--theme-motif-size) + var(--theme-motif-gap))) auto; opacity: 1; }
   body::after { content: ''; pointer-events: none; position: fixed; inset: 0; z-index: 0; background-color: transparent; background-image: var(--theme-motif-over-image); background-repeat: repeat; }
   body > * { position: relative; z-index: 2; }
   /* Fixed overlays (reading progress bar) must stay above the page wrapper. */
   body > #view-transition-progress { z-index: 1000; }
   [data-rendering-specimen="canonical-headings"] { position: relative; z-index: 4; }
  body::before { background-size: calc(180px * var(--theme-texture-scale)); }
	  body::after { transform: none; }
  h1, h2, h3, h4, h5, h6 { mask-image: none !important; -webkit-mask-image: none !important; background: none !important; color: inherit !important; -webkit-text-fill-color: currentColor !important; }
  /* Wear removes ink through a mask. It must not replace the semantic paint
     inherited from mark, links, emphasis, or strong. */
   .heading-wear-glyph { display: inline; position: relative; z-index: 5; color: inherit; background: none; -webkit-text-fill-color: currentColor; mask-image: var(--theme-heading-texture-mask) !important; -webkit-mask-image: var(--theme-heading-texture-mask) !important; mask-size: calc(180px * var(--theme-heading-texture-scale)); -webkit-mask-size: calc(180px * var(--theme-heading-texture-scale)); mask-repeat: repeat; -webkit-mask-repeat: repeat; }
  .heading-wear-glyph .heading-anchor, .heading-anchor { mask: none !important; background: none !important; color: inherit !important; -webkit-text-fill-color: currentColor !important; }
}
` + `
@layer tokens {
	  body, body::after { background-size: var(--theme-motif-field-size, calc(16 * (var(--theme-motif-size) + var(--theme-motif-gap)))) auto; background-position: 0 0; }
}
`
}

func unquoteCSSLengthProperty(css, property string) string {
	prefix := property + `: "`
	start := strings.Index(css, prefix)
	if start < 0 {
		return css
	}
	valueStart := start + len(prefix)
	end := strings.Index(css[valueStart:], `";`)
	if end < 0 {
		return css
	}
	end += valueStart
	return css[:start] + property + ": " + css[valueStart:end] + css[end+1:]
}

//nolint:gocyclo // CSS projection keeps independent legacy and canonical paths together.
func (p *AestheticCSSPlugin) generatePresentationCSSBody(config *lifecycle.Config) string {
	if config == nil {
		return ""
	}
	theme := models.NewThemeConfig()
	if configured, ok := config.Extra["models_config"].(*models.Config); ok {
		theme = configured.Theme
	}
	contract, err := renderingcontract.Load()
	if err != nil {
		return ""
	}
	textureKind := theme.Texture.Kind
	headingTextureKind := theme.HeadingTexture.Kind
	if headingTextureKind == renderingHeadingInherit {
		headingTextureKind = textureKind
	}
	// quiet keeps the texture on framing surfaces. The CSS projection below
	// excludes the reading surface instead of silently discarding the texture.
	textureMix := normalizedColorMix(theme.Texture.ColorMix)
	headingMix := normalizedColorMix(theme.HeadingTexture.ColorMix)
	motifMix := normalizedColorMix(theme.Motif.ColorMix)
	// Heading color_mix controls glyph wear. At zero the heading remains solid
	// text (the texture is background-equivalent); increasing the dial enables
	// the same glyph-scoped mask used by the browser consumers.
	headingColor := "var(--color-text, #222)"
	textureImage := renderingTextureNone
	headingImage := renderingTextureNone
	if bundle, err := compileMotifBundle(config); err == nil {
		for i := range bundle.Manifest.Passes {
			pass := &bundle.Manifest.Passes[i]
			switch pass.ID {
			case "surface-texture":
				textureImage = `url("/` + pass.Asset + `")`
			case "heading-wear":
				headingImage = `url("/` + pass.Asset + `")`
				// The recipe owns the semantic projection of color_mix. Do not
				// discard it at the CSS boundary and replace it with a fixed ink
				// color: that makes all non-zero endpoint builds paint identically.
				if pass.Paint != "" {
					headingColor = pass.Paint
				}
				// Keep the public dial as semantic color separation while using a
				// compiler-owned coverage projection for the alpha mask. This makes
				// the static endpoint builds materially distinct without changing
				// glyph geometry or using span opacity.
				if headingMix <= 0 {
					headingImage = renderingTextureNone
				} else {
					coverage := math.Max(0.25, 1-(headingMix*.75))
					headingImage = fmt.Sprintf("linear-gradient(rgba(0,0,0,%.3f),rgba(0,0,0,%.3f)), %s", coverage, coverage, headingImage)
				}
			}
		}
	} else {
		log.Printf("[aesthetic_css] compiled texture/heading assets unavailable: %v", err)
	}
	motifImage := renderingTextureNone
	motifMask := renderingTextureNone
	motifPaint := renderingTextureNone
	if theme.Motif.Kind != renderingMotifOff {
		motifColor := resolveMotifPaint(contract, theme, motifMix)
		motifPaint = motifColor
		if theme.Motif.Kind != renderingMotifBlockW {
			motifImage = motifImageFor(theme.Motif.Kind, motifColor, theme.Motif.Glyph)
		}
	}
	if bundle, err := compileMotifBundle(config); err == nil && len(bundle.Assets) != 0 && theme.Motif.Kind == renderingMotifBlockW {
		// Canonical consumers attach the compiler-owned field. Keep the local CSS
		// layer and pseudo-element integration, but do not regenerate geometry.
		if _, exists := bundle.Assets["assets/motif-block-w-v1.svg"]; exists {
			motifImage = `url("/assets/motif-block-w-v1.svg")`
		}
	}
	customMotifURL := validMotifURL(contract, theme.Motif.URL) && !strings.ContainsAny(theme.Motif.URL, "\"'()\n\r\t")
	// The canonical W is portable contract artwork, not a network dependency.
	// Keep the URL as authoring metadata but render the vendored path locally.
	if customMotifURL && theme.Motif.Kind == renderingMotifBlockW {
		// block-w is compiler-owned. Custom artwork is not silently substituted
		// into the canonical field because that would change its digest.
		customMotifURL = false
	}
	if customMotifURL && theme.Motif.URL == "https://waylonwalker.com/w.svg" {
		customMotifURL = false
	}
	if customMotifURL {
		// Custom artwork is not supported by the compiled block-w asset. Reject it
		// explicitly instead of falling back to a second geometry generator.
		log.Printf("[aesthetic_css] unsupported motif custom URL: block-w uses the canonical compiled asset")
	} else if theme.Motif.URL != "" {
		log.Printf("[aesthetic_css] invalid-or-unsafe-custom-url: motif URL ignored")
	}
	textureOpacity := 1.0
	if textureImage == renderingTextureNone || textureMix == 0 {
		textureOpacity = 0
	}
	motifZ := motifLayerZ(theme.Motif.Layer)
	aestheticCSS := ""
	canonicalHighlightRadius := "0px"
	if contract, err := renderingcontract.Load(); err == nil {
		if tokens, ok := contract.Aesthetics[theme.Aesthetic]; ok {
			aestheticCSS = fmt.Sprintf("--radius: %v; --shadow: %v; --space: %v;", tokens["radius"], tokens["shadow"], tokens["spacing"])
		}
		if presentation, ok := contract.Presentation["canonical_document"].(map[string]any); ok {
			if radius, ok := presentation["highlight_radius"].(float64); ok {
				canonicalHighlightRadius = fmt.Sprintf("%.2fpx", radius)
			}
		}
	}
	underImage := renderingTextureNone
	if theme.Motif.Layer == renderingLayerUnder || theme.Motif.Layer == renderingLayerSandwich {
		underImage = motifImage
	}
	overImage := renderingTextureNone
	if theme.Motif.Layer == renderingLayerOver || theme.Motif.Layer == renderingLayerSandwich {
		overImage = motifImage
	}
	variableNames := make([]string, 0, len(theme.Variables))
	for name := range theme.Variables {
		variableNames = append(variableNames, name)
	}
	sort.Strings(variableNames)
	for _, name := range variableNames {
		value := theme.Variables[name]
		if strings.HasPrefix(name, "--") && !strings.ContainsAny(name+value, "{};\n\r") {
			aestheticCSS += fmt.Sprintf(" %s: %s;", name, value)
		}
	}
	aestheticCSS += fmt.Sprintf(" --canonical-highlight-radius: %s;", canonicalHighlightRadius)
	liveMotif := liveMotifCSS(theme, motifImage, motifMix)
	return fmt.Sprintf("\n@layer tokens {\n  :root {\n    %s\n    --theme-contract-version: %d;\n    --theme-texture-kind: %q; --theme-texture-color-mix: %.3f; --theme-texture-scale: %.3f; --theme-texture-scope: %q; --theme-texture-opacity: %.3f; --theme-texture-image: %s;\n    --theme-heading-texture-kind: %q; --theme-heading-texture-color-mix: %.3f; --theme-heading-texture-scale: %.3f; --theme-heading-texture-color: %s; --theme-heading-texture-mask: %s;\n    --theme-motif-kind: %q; --theme-motif-glyph: %q; --theme-motif-color-mix: %.3f; --theme-motif-layer: %q; --theme-motif-image: %s; --theme-motif-under-image: %s; --theme-motif-over-image: %s; --theme-motif-mask: %s; --theme-motif-paint: %s; --theme-motif-size: %q; --theme-motif-gap: %q; --theme-motif-row-offset: %.3f; --theme-motif-wobble: %.3f; --theme-motif-scatter: %.3f; --theme-motif-color: %q; --theme-motif-url: %q; --theme-motif-z: %s;\n  }\n  body { background-image: var(--theme-motif-under-image); background-size: var(--theme-motif-field-size); background-position: 0 0; }\n  body::before { content: ''; pointer-events: none; position: fixed; inset: 0; z-index: -1; background-image: var(--theme-texture-image); background-size: calc(180px * var(--theme-texture-scale)); background-position: 0 0; opacity: var(--theme-texture-opacity); }\n  [data-theme-texture-scope=quiet] main, [data-theme-texture-scope=quiet] article, [data-theme-texture-scope=quiet] [data-reading-surface], [data-theme-texture-scope=quiet] .reading-surface { background-color: var(--color-background, #fff); }\n  h1, h2, h3, h4, h5, h6 { color: transparent; background-color: var(--color-text, #222); background-image: var(--theme-heading-texture-color); background-clip: text; -webkit-background-clip: text; -webkit-text-fill-color: transparent; mask-image: var(--theme-heading-texture-mask); -webkit-mask-image: var(--theme-heading-texture-mask); mask-size: calc(180px * var(--theme-heading-texture-scale)); -webkit-mask-size: calc(180px * var(--theme-heading-texture-scale)); mask-repeat: repeat; -webkit-mask-repeat: repeat; }\n  body::after { content: ''; pointer-events: none; position: fixed; inset: 0; z-index: var(--theme-motif-z); background-color: var(--theme-motif-paint); background-image: var(--theme-motif-over-image); mask-image: var(--theme-motif-mask); -webkit-mask-image: var(--theme-motif-mask); background-size: var(--theme-motif-field-size); mask-size: var(--theme-motif-field-size); background-position: 0 0; mask-position: 0 0; background-repeat: repeat; mask-repeat: repeat; opacity: 1; transform: none; }\n}\n", aestheticCSS, theme.ContractVersion, textureKind, textureMix, theme.Texture.Scale, theme.Texture.Scope, textureOpacity, textureImage, headingTextureKind, headingMix, theme.HeadingTexture.Scale, headingColor, headingImage, theme.Motif.Kind, theme.Motif.Glyph, motifMix, theme.Motif.Layer, motifImage, underImage, overImage, motifMask, motifPaint, theme.Motif.Size, theme.Motif.Gap, theme.Motif.RowOffset, theme.Motif.Wobble, theme.Motif.Scatter, theme.Motif.Color, theme.Motif.URL, motifZ) + liveMotif
}

// liveMotifCSS lets the motif follow a palette chosen in the theme picker.
// Compiled motif artwork bakes the configured palette's paint into the image,
// so when the page runs any other palette the artwork is reused as an alpha
// mask and painted from the live palette variables instead.
func liveMotifCSS(theme models.ThemeConfig, motifImage string, mix float64) string {
	if theme.Motif.Kind == renderingMotifOff || motifImage == "" || motifImage == renderingTextureNone {
		return ""
	}
	target := "var(--color-text)"
	weight := mix
	switch theme.Motif.Color {
	case renderingMotifAccent:
		target = "var(--color-link, var(--color-primary))"
	case renderingMotifMuted:
		weight = mix * .55
	case renderingMotifShadow:
		weight = mix * .28
	}
	z := "var(--theme-motif-z)"
	if theme.Motif.Layer == renderingLayerUnder {
		z = "0"
	}
	selector := "html[data-palette]"
	if theme.Palette != "" && !strings.ContainsAny(theme.Palette, "\"\\]\n\r") {
		selector += fmt.Sprintf(`:not([data-palette=%q])`, theme.Palette)
	}
	size := "var(--theme-motif-field-size, calc(16 * (var(--theme-motif-size) + var(--theme-motif-gap)))) auto"
	return fmt.Sprintf(`
@layer tokens {
  %[1]s body { background-image: none; }
  %[1]s body::after { z-index: %[2]s; background-image: none; background-color: color-mix(in srgb, %[3]s %.2[4]f%%, var(--color-background)); mask-image: %[5]s; -webkit-mask-image: %[5]s; mask-size: %[6]s; -webkit-mask-size: %[6]s; mask-repeat: repeat; -webkit-mask-repeat: repeat; mask-position: 0 0; -webkit-mask-position: 0 0; }
}
`, selector, z, target, weight*100, motifImage, size)
}

// compiledMotifBundle is the single bridge from the model configuration to
// the shared compiler. Unsupported themes keep the legacy local projection;
// they must not silently attach an asset compiled from different semantics.
func compileMotifBundle(config *lifecycle.Config) (renderingrecipe.Bundle, error) {
	if config == nil || config.Extra == nil {
		return renderingrecipe.Bundle{}, nil
	}
	configured, ok := config.Extra["models_config"].(*models.Config)
	if !ok {
		return renderingrecipe.Bundle{}, nil
	}
	theme := configured.Theme
	// Recipe v1 is the canonical brush projection. Existing bundled fontpacks
	// keep their legacy CSS projection until their recipes are compiled.
	if theme.Fontpack != "" && theme.Fontpack != renderingFontpackBrush {
		return renderingrecipe.Bundle{}, nil
	}
	if theme.Fontpack == "" {
		theme.Fontpack = renderingFontpackBrush
	}
	// Site-local palettes (./palettes) are not part of the rendering contract,
	// so the compiled motif bundle is skipped; the legacy CSS projection is used.
	if !contractHasPalette(theme.Palette) {
		return renderingrecipe.Bundle{}, nil
	}
	// An omitted texture is the canonical no-op. Keep the other supported
	// projections compilable: the heading wear and motif passes are independent
	// of the surface texture pass.
	if strings.TrimSpace(theme.Texture.Kind) == "" {
		theme.Texture.Kind = renderingTextureNone
	}
	if strings.TrimSpace(theme.Texture.Scope) == "" {
		theme.Texture.Scope = renderingScopeAll
	}
	if theme.HeadingTexture.Kind == renderingHeadingInherit {
		theme.HeadingTexture.Kind = renderingHeadingSplatter
	}
	// A disabled motif still needs the texture and heading passes. Compile
	// them with the canonical motif geometry and drop the unused motif asset.
	motifOff := theme.Motif.Kind == renderingMotifOff
	if motifOff {
		theme.Motif.Kind = renderingMotifBlockW
		theme.Motif.URL = ""
	}
	if (theme.Texture.Kind != renderingTextureNone && theme.Texture.Kind != renderingTextureScreenprint) || theme.HeadingTexture.Kind != renderingHeadingSplatter || theme.Motif.Kind != renderingMotifBlockW {
		return renderingrecipe.Bundle{}, fmt.Errorf("unsupported rendering recipe: texture=%q heading_texture=%q motif=%q", theme.Texture.Kind, theme.HeadingTexture.Kind, theme.Motif.Kind)
	}
	if theme.Motif.URL != "" && theme.Motif.URL != "https://waylonwalker.com/w.svg" {
		return renderingrecipe.Bundle{}, fmt.Errorf("motif block-w custom URL is unsupported by the compiler")
	}
	bundleTheme := renderingrecipe.Theme{
		ContractVersion: theme.ContractVersion,
		Palette:         theme.Palette,
		Aesthetic:       theme.Aesthetic,
		Fontpack:        theme.Fontpack,
		Texture:         renderingcontract.Dial{Kind: theme.Texture.Kind, ColorMix: theme.Texture.ColorMix, Scale: theme.Texture.Scale, Scope: theme.Texture.Scope},
		HeadingTexture:  renderingcontract.Dial{Kind: theme.HeadingTexture.Kind, ColorMix: theme.HeadingTexture.ColorMix, Scale: theme.HeadingTexture.Scale},
		Motif:           renderingcontract.MotifState{Kind: theme.Motif.Kind, ColorMix: theme.Motif.ColorMix, Glyph: theme.Motif.Glyph, Size: theme.Motif.Size, Gap: theme.Motif.Gap, RowOffset: theme.Motif.RowOffset, Wobble: theme.Motif.Wobble, Scatter: theme.Motif.Scatter, Layer: theme.Motif.Layer, Color: theme.Motif.Color, URL: theme.Motif.URL},
		Variables:       theme.Variables,
	}
	bundle, err := renderingrecipe.Compile(bundleTheme)
	if err != nil {
		return renderingrecipe.Bundle{}, fmt.Errorf("compile motif bundle: %w", err)
	}
	if motifOff {
		delete(bundle.Assets, "assets/motif-block-w-v1.svg")
	}
	return bundle, nil
}

// hashedAssetReadyInFastMode reports whether a fast (incremental) build can
// skip regenerating a hashed CSS asset because the file from a previous build
// is already present in the output directory.
func hashedAssetReadyInFastMode(m *lifecycle.Manager, asset string) bool {
	config := m.Config()
	if config == nil || config.Extra == nil {
		return false
	}
	fast, ok := config.Extra["fast_mode"].(bool)
	if !ok || !fast {
		return false
	}
	hash := m.GetAssetHash(asset)
	if hash == "" {
		return false
	}
	base := strings.TrimSuffix(asset, filepath.Ext(asset))
	hashed := filepath.Join(config.OutputDir, filepath.FromSlash(fmt.Sprintf("%s.%s%s", base, hash, filepath.Ext(asset))))
	_, err := os.Stat(hashed)
	return err == nil
}

func contractHasPalette(id string) bool {
	if id == "" {
		return true
	}
	contract, err := renderingcontract.Load()
	if err != nil {
		return false
	}
	if alias, ok := contract.Aliases[id]; ok {
		id = alias
	}
	for _, palette := range contract.Palettes {
		if palette.ID == id {
			return true
		}
	}
	return false
}

func resolveMotifPaint(contract renderingcontract.Contract, theme models.ThemeConfig, mix float64) string {
	for _, palette := range contract.Palettes {
		if palette.ID != theme.Palette && contract.Aliases[theme.Palette] != palette.ID {
			continue
		}
		final := renderingcontract.FinalRenderPalette(palette)
		background := final["background"]
		target := final["text"]
		switch theme.Motif.Color {
		case "accent":
			target = final["link"]
		case "muted":
			target = mixHex(background, final["text"], .55)
		case "shadow":
			target = mixHex(background, final["text"], .28)
		}
		return mixHex(background, target, mix)
	}
	return sharedColorMix(theme.Motif.Color, mix)
}

func mixHex(background, foreground string, amount float64) string {
	parse := func(value string) ([3]int, bool) {
		value = strings.TrimPrefix(strings.TrimSpace(value), "#")
		if len(value) != 6 {
			return [3]int{}, false
		}
		var result [3]int
		for index := range result {
			var channel int
			if _, err := fmt.Sscanf(value[index*2:index*2+2], "%02x", &channel); err != nil {
				return result, false
			}
			result[index] = channel
		}
		return result, true
	}
	bg, ok := parse(background)
	if !ok {
		return foreground
	}
	fg, ok := parse(foreground)
	if !ok {
		return foreground
	}
	if amount < 0 {
		amount = 0
	}
	if amount > 1 {
		amount = 1
	}
	return fmt.Sprintf("#%02x%02x%02x", int(float64(bg[0])+(float64(fg[0])-float64(bg[0]))*amount+.5), int(float64(bg[1])+(float64(fg[1])-float64(bg[1]))*amount+.5), int(float64(bg[2])+(float64(fg[2])-float64(bg[2]))*amount+.5))
}

func normalizedColorMix(value float64) float64 { return math.Max(0, math.Min(1, value)) }

func validMotifURL(contract renderingcontract.Contract, value string) bool {
	if value == "" {
		return false
	}
	for _, scheme := range contract.MotifURLSchemes {
		if strings.HasPrefix(strings.ToLower(value), strings.ToLower(scheme)) {
			return true
		}
	}
	return false
}

func sharedColorMix(role string, mix float64) string {
	color := "var(--color-text, #222)"
	switch role {
	case "accent":
		color = "var(--color-primary, #369)"
	case "muted":
		color = "color-mix(in srgb, var(--color-text, #222) 55%, var(--color-background, #fff))"
	case "shadow":
		color = "color-mix(in srgb, var(--color-text, #222) 28%, var(--color-background, #fff))"
	}
	return fmt.Sprintf("color-mix(in srgb, var(--color-background, #fff) %.1f%%, %s %.1f%%)", (1-mix)*100, color, mix*100)
}

func motifImageFor(kind, color, glyph string) string {
	if kind == renderingMotifLetter {
		if glyph == "" {
			glyph = "W"
		}
		glyph = strings.ReplaceAll(strings.ReplaceAll(glyph, "&", "&amp;"), "<", "&lt;")
		svg := fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' width='100' height='100'><text x='50' y='70' text-anchor='middle' font-size='72' fill='%s'>%s</text></svg>", color, glyph)
		return `url("data:image/svg+xml;base64,` + base64.StdEncoding.EncodeToString([]byte(svg)) + `")`
	}
	return fmt.Sprintf("linear-gradient(135deg, transparent 0 42%%, %s 42%% 58%%, transparent 58%% 100%%)", color)
}

func motifLayerZ(layer string) string {
	switch layer {
	case "under":
		return "-2"
	case "sandwich":
		return "0"
	default:
		return "1"
	}
}

func (p *AestheticCSSPlugin) isSwitcherEnabled(extra map[string]interface{}) bool {
	return themeSwitcherEnabled(extra)
}

func (p *AestheticCSSPlugin) getSwitcherConfig(extra map[string]interface{}) models.ThemeSwitcherConfig {
	if extra == nil {
		return models.NewThemeSwitcherConfig()
	}
	if themeConfig, ok := extra["theme"].(models.ThemeConfig); ok {
		return themeConfig.Switcher
	}
	return models.NewThemeSwitcherConfig()
}

func (p *AestheticCSSPlugin) generateSingleAestheticCSS(loader *aesthetic.Loader, aestheticName string) string {
	a, err := loader.Load(aestheticName)
	if err != nil {
		return ""
	}

	var buf bytes.Buffer
	buf.WriteString("@layer reset, tokens, base, components, utilities, overrides;\n\n")
	buf.WriteString("@layer tokens {\n")
	buf.WriteString(fmt.Sprintf("/* Aesthetic: %s */\n", a.Name))

	// Just use the root level CSS block directly, but strip out the :root { ... } to format it ourselves
	// or we can just append a.GenerateCSS() directly. Wait, a.GenerateCSS() returns :root { ... }
	cssContent := a.GenerateCSS()
	buf.WriteString(cssContent)

	buf.WriteString("}\n")
	return buf.String()
}

type AestheticManifestEntry struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

func (p *AestheticCSSPlugin) generateMultiAestheticCSS(loader *aesthetic.Loader, extra map[string]interface{}, defaultAestheticName string) string {
	var buf bytes.Buffer

	buf.WriteString("@layer reset, tokens, base, components, utilities, overrides;\n\n")
	buf.WriteString("@layer tokens {\n")

	allAesthetics, err := loader.Discover()
	if err != nil {
		return p.generateSingleAestheticCSS(loader, defaultAestheticName)
	}

	// Filter based on switcher
	// ... we will filter just by name for now, using the palette include/exclude logic
	switcherConfig := p.getSwitcherConfig(extra)
	filteredAesthetics := p.filterAesthetics(allAesthetics, switcherConfig)

	manifest := make([]AestheticManifestEntry, 0, len(filteredAesthetics))
	for _, info := range filteredAesthetics {
		manifest = append(manifest, AestheticManifestEntry{
			Name:        aestheticID(info.Name),
			DisplayName: info.Name,
		})
	}

	sort.Slice(manifest, func(i, j int) bool {
		return manifest[i].Name < manifest[j].Name
	})

	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		manifestJSON = []byte("[]")
	}
	escapedManifest := strings.ReplaceAll(string(manifestJSON), "'", "\\'")

	buf.WriteString("/* Global aesthetic configuration */\n")
	buf.WriteString(":root {\n")
	buf.WriteString(fmt.Sprintf("  --aesthetic-default: %q;\n", defaultAestheticName))
	buf.WriteString(fmt.Sprintf("  --aesthetic-manifest: '%s';\n", escapedManifest))
	buf.WriteString("}\n\n")

	for _, info := range filteredAesthetics {
		a, err := loader.Load(info.Name)
		if err != nil {
			continue
		}

		buf.WriteString(fmt.Sprintf("/* Aesthetic: %s */\n", info.Name))
		// Scope by the normalized id so it matches the configured aesthetic
		// name rendered into data-aesthetic (e.g. "minimal", not "Minimal").
		buf.WriteString(fmt.Sprintf("[data-aesthetic=%q] {\n", aestheticID(info.Name)))

		// Extract variables directly to inject in [data-aesthetic] scope
		p.writeAestheticVariablesIndented(&buf, a, "  ")

		buf.WriteString("}\n\n")
	}

	if defaultAestheticName != "" {
		a, err := loader.Load(defaultAestheticName)
		if err == nil {
			buf.WriteString(fmt.Sprintf("/* Default aesthetic - %s */\n", defaultAestheticName))
			buf.WriteString(":root:not([data-aesthetic]) {\n")
			p.writeAestheticVariablesIndented(&buf, a, "  ")
			buf.WriteString("}\n\n")
		}
	}

	buf.WriteString("}\n")
	return buf.String()
}

// aestheticID normalizes an aesthetic name the same way the aesthetic loader does.
func aestheticID(name string) string {
	return strings.NewReplacer(" ", "-", "_", "-").Replace(strings.ToLower(name))
}

func (p *AestheticCSSPlugin) writeAestheticVariablesIndented(buf *bytes.Buffer, a *aesthetic.Aesthetic, indent string) {
	cssLines := strings.Split(a.GenerateCSS(), "\n")
	for _, line := range cssLines {
		if strings.HasPrefix(line, "--") {
			buf.WriteString(indent + line + "\n")
		} else if strings.HasPrefix(strings.TrimSpace(line), "--") {
			buf.WriteString(indent + strings.TrimSpace(line) + "\n")
		}
	}
}

func (p *AestheticCSSPlugin) filterAesthetics(all []aesthetic.Info, switcherConfig models.ThemeSwitcherConfig) []aesthetic.Info {
	if switcherConfig.IsIncludeAll() {
		excludeSet := make(map[string]bool)
		for _, name := range switcherConfig.Exclude {
			excludeSet[strings.ToLower(name)] = true
		}

		var result []aesthetic.Info
		for _, info := range all {
			lowerName := strings.ToLower(info.Name)
			if !excludeSet[lowerName] {
				result = append(result, info)
			}
		}
		return result
	}

	includeSet := make(map[string]bool)
	for _, name := range switcherConfig.Include {
		includeSet[strings.ToLower(name)] = true
	}

	var result []aesthetic.Info
	for _, info := range all {
		lowerName := strings.ToLower(info.Name)
		if includeSet[lowerName] {
			result = append(result, info)
		}
	}
	return result
}

func (p *AestheticCSSPlugin) getAestheticConfig(extra map[string]interface{}) string {
	if extra == nil {
		return ""
	}

	if modelsConfig, ok := extra["models_config"].(*models.Config); ok {
		if modelsConfig.Theme.Aesthetic != "" {
			return modelsConfig.Theme.Aesthetic
		}
	}

	if themeConfig, ok := extra["theme"].(models.ThemeConfig); ok {
		return themeConfig.Aesthetic
	}

	theme, ok := extra["theme"].(map[string]interface{})
	if !ok {
		return ""
	}

	if asth, ok := theme["aesthetic"].(string); ok && asth != "" {
		return asth
	}

	return ""
}

func (p *AestheticCSSPlugin) Priority(_ lifecycle.Stage) int {
	return lifecycle.PriorityDefault
}

var (
	_ lifecycle.Plugin          = (*AestheticCSSPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*AestheticCSSPlugin)(nil)
	_ lifecycle.WritePlugin     = (*AestheticCSSPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*AestheticCSSPlugin)(nil)
)
