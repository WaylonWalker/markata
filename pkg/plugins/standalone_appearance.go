package plugins

import (
	"fmt"
	"io/fs"
	"path"
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
	SchemaVersion int               `json:"schema_version"`
	PickerHTML    string            `json:"picker_html"`
	Assets        map[string][]byte `json:"assets"`
	Styles        []string          `json:"styles"`
	Scripts       []string          `json:"scripts"`
}

// BuildStandaloneAppearance uses the same embedded catalogs and template as sites.
// It does not discover a site, write files, or download assets.
func BuildStandaloneAppearance() (*StandaloneAppearance, error) {
	theme := models.ThemeConfig{Palette: "default", Aesthetic: "minimal", Fontpack: "system", Switcher: models.NewThemeSwitcherConfig()}
	extra := map[string]interface{}{"theme": theme}
	r := &StandaloneAppearance{SchemaVersion: 1, Assets: map[string][]byte{},
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
	paletteCSS := NewPaletteCSSPlugin().generateMultiPaletteCSS(palettes.NewLoaderWithPaths(nil), extra, "default", "default-light", "default-dark", nil, "dark")
	r.Assets["css/palette.css"] = []byte(paletteCSS)
	a := NewAestheticCSSPlugin()
	r.Assets["css/aesthetic.css"] = append([]byte(a.generateMultiAestheticCSS(aesthetic.NewLoaderWithPaths(nil), extra, "minimal")+aestheticSurfaceCSS), r.Assets["css/aesthetic.css"]...)
	source, err := fontpacks.BuiltinSource()
	if err != nil {
		return nil, err
	}
	names := append([]string{"system"}, fontpacks.SortedKeys(source.Catalog.FontPacks)...)
	fonts, err := source.Catalog.ResolveManyFSWithOptions(names, source.FS, source.Root, "", fontpacks.ResolveOptions{})
	if err != nil {
		return nil, err
	}
	r.Assets["css/fonts.css"] = []byte(fonts.CSS + fontpackManifestCSS(source.Catalog, "system", fonts.Packs))
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
	if schedule, ok := extra["palette_seasonal_json"].(string); ok {
		r.Assets["js/standalone-appearance.js"] = append([]byte("window.__markataSeasonal = "+schedule+";\n"), r.Assets["js/standalone-appearance.js"]...)
	}
	return r, nil
}
