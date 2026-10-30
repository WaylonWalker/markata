package plugins

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/aesthetic"
	"github.com/WaylonWalker/markata/pkg/fontpacks"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/templates"
	"github.com/WaylonWalker/markata/pkg/themes"
)

// StandaloneAppearance is a reusable appearance component and its ordered assets.
// Assets contains bytes, encoded as base64 by encoding/json.
type StandaloneAppearance struct {
	Defaults      map[string]string `json:"defaults,omitempty"`
	SchemaVersion int               `json:"schema_version"`
	PickerHTML    string            `json:"picker_html"`
	Assets        map[string][]byte `json:"assets"`
	Styles        []string          `json:"styles"`
	Scripts       []string          `json:"scripts"`
}

// BuildStandaloneAppearance uses the same embedded catalogs and template as sites.
// It does not discover a site, write files, or download assets.
func BuildStandaloneAppearance() (*StandaloneAppearance, error) {
	return BuildStandaloneAppearanceForSite(nil, "")
}

type standaloneSiteStyle struct {
	theme                                 models.ThemeConfig
	fontpack, mode, light, dark, selected string
}

func standaloneSiteTheme(cfg *models.Config) standaloneSiteStyle {
	var theme models.ThemeConfig
	var fontpack, mode, light, dark, selected string
	theme = models.ThemeConfig{Palette: templateTypeDefault, PaletteLight: "default-light", PaletteDark: "default-dark", Aesthetic: "minimal", Fontpack: "system", Switcher: models.NewThemeSwitcherConfig()}
	if cfg != nil {
		theme = cfg.Theme
	}
	if theme.Palette == "" {
		theme.Palette = templateTypeDefault
	}
	if theme.Aesthetic == "" {
		theme.Aesthetic = "minimal"
	}
	fontpack = theme.Fontpack
	if fontpack == "" && cfg != nil {
		fontpack = cfg.Fontpack
	}
	if fontpack == "" {
		fontpack = "system"
	}
	mode = theme.FallbackMode
	if mode != "light" {
		mode = "dark"
	}
	light, dark = palettes.GetEffectivePalettes(theme.Palette, theme.PaletteLight, theme.PaletteDark)
	selected = dark
	if mode == "light" {
		selected = light
	}
	return standaloneSiteStyle{theme, fontpack, mode, light, dark, selected}
}

// BuildStandaloneAppearanceForSite applies site appearance settings and local catalogs.
func BuildStandaloneAppearanceForSite(cfg *models.Config, root string) (*StandaloneAppearance, error) {
	style := standaloneSiteTheme(cfg)
	theme, fontpack, mode := style.theme, style.fontpack, style.mode
	light, dark, selected := style.light, style.dark, style.selected
	extra := map[string]interface{}{"theme": theme}
	paletteLoader := palettes.NewLoaderWithPaths(nil)
	aestheticLoader := aesthetic.NewLoaderWithPaths(nil)
	if root != "" {
		paletteLoader.AddPath(filepath.Join(root, "palettes"))
		aestheticLoader.AddPath(filepath.Join(root, "aesthetics"))
	}
	r := &StandaloneAppearance{Defaults: map[string]string{"data-theme": mode, "data-palette": normalizePaletteName(selected), "data-aesthetic": theme.Aesthetic, "data-fontpack": fontpack, "data-text-size": theme.EffectiveTextSize()}, SchemaVersion: 1, Assets: map[string][]byte{},
		Styles:  []string{"css/variables.css", "css/main.css", "css/components.css", "css/palette.css", "css/aesthetic.css", "css/fonts.css", "css/text-size.css", "css/palette-switcher.css", "css/standalone-appearance.css"},
		Scripts: []string{"js/standalone-appearance.js", "js/text-size.js", "js/adaptive-loading.js", "js/palette-switcher.js"}}
	for _, name := range append(append([]string{}, r.Styles...), r.Scripts...) {
		if strings.Contains(name, "standalone-appearance") || name == "css/palette.css" || name == "css/fonts.css" {
			continue
		}
		data, err := themes.ReadStatic(name)
		if err != nil {
			return nil, err
		}
		r.Assets[name] = data
	}
	paletteCSS := NewPaletteCSSPlugin().generateMultiPaletteCSS(paletteLoader, extra, theme.Palette, light, dark, theme.Variables, mode)
	r.Assets["css/palette.css"] = []byte(paletteCSS)
	a := NewAestheticCSSPlugin()
	r.Assets["css/aesthetic.css"] = append([]byte(a.generateMultiAestheticCSS(aestheticLoader, extra, theme.Aesthetic)+aestheticSurfaceCSS), r.Assets["css/aesthetic.css"]...)
	source, err := fontpacks.BuiltinSource()
	if cfg != nil && cfg.FontpacksFile != "" {
		filename := cfg.FontpacksFile
		if !filepath.IsAbs(filename) {
			filename = filepath.Join(root, filename)
		}
		source, err = fontpacks.LoadSource(filename)
	}
	if err != nil {
		return nil, err
	}
	names := append([]string{"system"}, fontpacks.SortedKeys(source.Catalog.FontPacks)...)
	fonts, err := source.Catalog.ResolveManyFSWithOptions(names, source.FS, source.Root, "", fontpacks.ResolveOptions{})
	if err != nil {
		return nil, err
	}
	r.Assets["css/fonts.css"] = []byte(fonts.CSS + fontpackManifestCSS(source.Catalog, fontpack, fonts.Packs))
	for i := range fonts.Assets {
		asset := &fonts.Assets[i]
		data, err := fs.ReadFile(source.FS, path.Join(source.Root, asset.Source, asset.File))
		if err != nil {
			return nil, fmt.Errorf("font %s: %w", asset.File, err)
		}
		r.Assets[strings.TrimPrefix(asset.URL, "/")] = data
	}
	engine := templates.NewEmbeddedEngine()
	r.PickerHTML, err = engine.RenderToString("partials/palette-switcher.html", map[string]interface{}{"config": map[string]interface{}{"theme": map[string]interface{}{"switcher": map[string]interface{}{"enabled": true, "mode_toggle": true}, "show_text_size_control": true}, "header": map[string]interface{}{"show_theme_toggle": true}}})
	if err != nil {
		return nil, err
	}
	for _, name := range []string{"css/standalone-appearance.css", "js/standalone-appearance.js"} {
		data, err := themes.ReadStatic(name)
		if err != nil {
			return nil, err
		}
		r.Assets[name] = data
	}
	defaults, err := json.Marshal(map[string]string{"light": normalizePaletteName(light), "dark": normalizePaletteName(dark)})
	if err != nil {
		return nil, err
	}
	script := string(r.Assets["js/standalone-appearance.js"])
	script = strings.Replace(script, "window.__markataThemeFallbackMode = 'dark';", "window.__markataThemeFallbackMode = '"+mode+"';", 1)
	script = strings.Replace(script, "window.__markataThemeDefaults = {light:'default-light',dark:'default-dark'};", "window.__markataThemeDefaults = "+string(defaults)+";", 1)
	r.Assets["js/standalone-appearance.js"] = []byte(script)
	if err := r.addSiteCSS(theme.CustomCSS, root); err != nil {
		return nil, err
	}
	if schedule, ok := extra["palette_seasonal_json"].(string); ok {
		r.Assets["js/standalone-appearance.js"] = append([]byte("window.__markataSeasonal = "+schedule+";\n"), r.Assets["js/standalone-appearance.js"]...)
	}
	return r, nil
}

func (r *StandaloneAppearance) addSiteCSS(filename, root string) error {
	if filename == "" {
		return nil
	}
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(root, filename)
	}
	css, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read configured custom CSS: %w", err)
	}
	r.Assets["css/site-custom.css"] = css
	r.Styles = append(r.Styles, "css/site-custom.css")
	return nil
}
