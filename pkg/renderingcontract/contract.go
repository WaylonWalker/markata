// Package renderingcontract contains the versioned, language-neutral presentation contract.
package renderingcontract

import (
	"embed"
	"encoding/json"
	"fmt"
	"math"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

const (
	textureScopeQuiet = "quiet"
)

//go:embed contract-v1.json
var contractFS embed.FS

type Palette struct {
	ID      string            `json:"id"`
	Family  string            `json:"family"`
	Variant string            `json:"variant"`
	Name    string            `json:"name"`
	Roles   map[string]string `json:"roles"`
}
type Dial struct {
	Kind      string  `json:"kind"`
	ColorMix  float64 `json:"color_mix"`
	Scale     float64 `json:"scale"`
	Scope     string  `json:"scope,omitempty"`
	Glyph     string  `json:"glyph,omitempty"`
	Size      string  `json:"size,omitempty"`
	Gap       string  `json:"gap,omitempty"`
	RowOffset float64 `json:"row_offset,omitempty"`
	Wobble    float64 `json:"wobble,omitempty"`
	Scatter   float64 `json:"scatter,omitempty"`
	Layer     string  `json:"layer,omitempty"`
	Color     string  `json:"color,omitempty"`
	URL       string  `json:"url,omitempty"`
}

// MotifState is the public normalized motif shape. It intentionally excludes
// the old internal scale artifact; motif scale is not a v1 authoring key.
type MotifState struct {
	Kind      string  `json:"kind"`
	ColorMix  float64 `json:"color_mix"`
	Glyph     string  `json:"glyph,omitempty"`
	Size      string  `json:"size,omitempty"`
	Gap       string  `json:"gap,omitempty"`
	RowOffset float64 `json:"row_offset"`
	Wobble    float64 `json:"wobble"`
	Scatter   float64 `json:"scatter"`
	Layer     string  `json:"layer,omitempty"`
	Color     string  `json:"color,omitempty"`
	URL       string  `json:"url,omitempty"`
}
type Defaults struct {
	Palette, Aesthetic, Fontpack string
	Texture                      Dial `json:"texture"`
	HeadingTexture               Dial `json:"heading_texture"`
	Motif                        Dial `json:"motif"`
}
type Contract struct {
	ContractVersion           int                           `json:"contract_version"`
	Defaults                  Defaults                      `json:"defaults"`
	Enums                     map[string][]string           `json:"enums"`
	Bounds                    map[string][2]float64         `json:"bounds"`
	Aesthetics                map[string]map[string]any     `json:"aesthetics"`
	IncompletePaletteFamilies []string                      `json:"incomplete_palette_families"`
	Textures                  map[string]map[string]float64 `json:"textures"`
	Fontpacks                 map[string]map[string]string  `json:"fontpacks"`
	ConfigKeys                map[string][]string           `json:"config_keys"`
	Palettes                  []Palette                     `json:"palettes"`
	Aliases                   map[string]string             `json:"aliases"`
	MotifURLSchemes           []string                      `json:"motif_url_schemes"`
	Presentation              map[string]any                `json:"presentation"`
}

// FinalRenderPalette is the semantic palette emitted by Markata's final CSS
// projection. It is deliberately not the raw contract palette: the projection
// applies the same AA-preserving adjustments as palette_css.go.
func FinalRenderPalette(p Palette) map[string]string {
	background := p.Roles["background"]
	surface := p.Roles["surface"]
	ink := finalAccessibleColor(p.Roles["ink"], background, 4.5)
	ink = finalAccessibleColor(ink, surface, 4.5)
	// The accent is used as normal link text in the active render, so it does
	// not receive a large-text exemption. Component-only boundaries are audited
	// separately at 3:1 below.
	accent := finalAccessibleColor(p.Roles["accent"], background, 4.5)
	return map[string]string{
		"background": background,
		"surface":    surface,
		"text":       ink,
		"link":       accent,
		"focus":      accent,
		"border":     ink,
	}
}

// FinalRenderColor applies the final CSS projection's contrast adjustment.
func FinalRenderColor(foreground, background string, required float64) string {
	return finalAccessibleColor(foreground, background, required)
}

// FinalRenderContrast checks the values a Markata render actually emits.
// Normal text stays at 4.5:1; controls and focus indicators use 3:1.
func FinalRenderContrast(p Palette) []palettes.ContrastCheck {
	roles := FinalRenderPalette(p)
	specs := []struct {
		fg, bg   string
		required float64
		level    string
	}{
		{"text", "background", 4.5, "AA"},
		{"text", "surface", 4.5, "AA"},
		{"link", "background", 4.5, "AA"},
		{"focus", "background", 3, "UI"},
		{"border", "surface", 3, "UI"},
		{"link", "background", 4.5, "AA highlight"},
	}
	results := make([]palettes.ContrastCheck, 0, len(specs))
	for _, spec := range specs {
		ratio, err := palettes.ContrastRatioFromHex(roles[spec.fg], roles[spec.bg])
		results = append(results, palettes.ContrastCheck{Foreground: spec.fg, Background: spec.bg, ForegroundHex: roles[spec.fg], BackgroundHex: roles[spec.bg], Ratio: ratio, Required: spec.required, Level: spec.level, Passed: err == nil && ratio >= spec.required})
	}
	return results
}

func finalAccessibleColor(foreground, background string, required float64) string {
	fg, fgErr := palettes.ParseHexColor(foreground)
	bg, bgErr := palettes.ParseHexColor(background)
	if fgErr != nil || bgErr != nil || palettes.ContrastRatio(fg, bg) >= required {
		return foreground
	}
	if adjusted, ok := fg.AdjustForContrast(bg, required); ok {
		if palettes.ContrastRatio(adjusted, bg) >= required {
			return adjusted.Hex()
		}
	}
	// Keep the final projection a real pass even for colors that cannot be
	// adjusted by the legacy HSL helper due to 8-bit rounding.
	black, err := palettes.ParseHexColor("#000000")
	if err != nil {
		return foreground
	}
	white, err := palettes.ParseHexColor("#ffffff")
	if err != nil {
		return foreground
	}
	if palettes.ContrastRatio(black, bg) >= required {
		return black.Hex()
	}
	if palettes.ContrastRatio(white, bg) >= required {
		return white.Hex()
	}
	return foreground
}

// NormalizeTheme migrates legacy flat settings into the canonical nested shape.
// It intentionally accepts map data so TOML, YAML, JSON, and browser consumers
// can share the same precedence rule.
//
//nolint:gocyclo // Legacy-key migration must preserve each field's precedence and warning.
func NormalizeTheme(raw map[string]any) (normalized map[string]any, warnings []string) {
	out := map[string]any{}
	warnings = []string{}
	nested, nestedOK := raw["theme"].(map[string]any)
	if !nestedOK {
		nested = nil
	}
	for key, value := range nested {
		out[key] = value
	}
	resolve := func(key, legacy string, fallback any) {
		if value, ok := nested[key]; ok {
			out[key] = value
			if old, exists := raw[legacy]; exists && fmt.Sprint(value) != fmt.Sprint(old) {
				warnings = append(warnings, legacy+" conflicts with theme."+key+"; canonical value wins")
			}
			return
		}
		if value, ok := raw[legacy]; ok {
			out[key] = value
			return
		}
		out[key] = fallback
	}
	resolve("contract_version", "contract_version", 1)
	resolve("palette", "palette", "ayu-dark")
	resolve("aesthetic", "aesthetic", "minimal")
	// Match the source rendering contract and the core zero-download default.
	resolve("fontpack", "fontpack", "system")
	if value, ok := raw["texture"].(string); ok {
		if _, exists := nested["texture"]; !exists {
			out["texture"] = map[string]any{"kind": value}
		} else if texture, ok := nested["texture"].(map[string]any); ok {
			if canonical, exists := texture["kind"]; exists && fmt.Sprint(canonical) != value {
				warnings = append(warnings, "texture conflicts with theme.texture.kind; canonical value wins")
			}
		}
	}
	canonicalTextureKind := false
	canonicalTextureScope, hasCanonicalTextureScope := "", false
	if texture, ok := nested["texture"].(map[string]any); ok {
		_, canonicalTextureKind = texture["kind"]
		canonicalTextureScope, hasCanonicalTextureScope = texture["scope"].(string)
	}
	if value, ok := raw["texture_scope"]; ok {
		texture, textureOK := out["texture"].(map[string]any)
		if !textureOK {
			texture = nil
		}
		if texture == nil {
			texture = map[string]any{}
			out["texture"] = texture
		}
		if scope, ok := value.(string); ok && scope == "headings" {
			if hasCanonicalTextureScope && canonicalTextureScope != textureScopeQuiet {
				warnings = append(warnings, "texture_scope=headings conflicts with theme.texture.scope; canonical value wins")
			}
			// The legacy setting selected the heading-only projection of the
			// texture. Keep that meaning instead of emitting an invalid scope.
			if _, exists := out["heading_texture"]; !exists {
				out["heading_texture"] = map[string]any{}
			}
			heading, headingOK := out["heading_texture"].(map[string]any)
			if !headingOK {
				heading = map[string]any{}
				out["heading_texture"] = heading
			}
			if _, exists := heading["kind"]; !exists {
				if kind, ok := texture["kind"]; ok {
					heading["kind"] = kind
				}
			}
			if !hasCanonicalTextureScope {
				texture["scope"] = textureScopeQuiet
			}
			// The old setting disabled the surface texture everywhere. Keep that
			// visible effect in the dedicated heading-texture projection.
			if canonicalTextureKind {
				warnings = append(warnings, "texture_scope=headings conflicts with theme.texture.kind; canonical value wins")
			} else {
				texture["kind"] = "none"
			}
			warnings = append(warnings, "texture_scope=headings migrated to theme.heading_texture.kind; headings-only behavior preserved")
		} else if canonical, exists := texture["scope"]; exists {
			if fmt.Sprint(canonical) != fmt.Sprint(value) {
				warnings = append(warnings, "texture_scope conflicts with theme.texture.scope; canonical value wins")
			}
		} else {
			texture["scope"] = value
		}
	}
	resolveMix := func(table, field, legacy string, fallback float64) {
		values, valuesOK := out[table].(map[string]any)
		if !valuesOK {
			values = nil
		}
		if values == nil {
			values = map[string]any{}
			out[table] = values
		}
		if canonical, ok := values[field]; ok {
			if old, exists := raw[legacy]; exists && NormalizeMix(canonical, fallback) != NormalizeMix(old, fallback) {
				warnings = append(warnings, legacy+" conflicts with theme."+table+"."+field+"; canonical value wins")
			}
			values[field] = NormalizeMix(canonical, fallback)
			return
		}
		if old, ok := raw[legacy]; ok {
			values[field] = NormalizeMix(old, fallback)
			return
		}
		values[field] = fallback
	}
	resolveMix("texture", "color_mix", "texture_strength", .35)
	resolveScale := func(table, field, legacy string, fallback float64) {
		values, valuesOK := out[table].(map[string]any)
		if !valuesOK {
			values = nil
		}
		if values == nil {
			values = map[string]any{}
			out[table] = values
		}
		value, ok := values[field]
		canonical := ok
		if !ok {
			value, ok = raw[legacy]
		}
		if !ok {
			values[field] = fallback
			return
		}
		n := NormalizeScale(value, fallback)
		if canonical {
			if old, exists := raw[legacy]; exists && n != NormalizeScale(old, fallback) {
				warnings = append(warnings, legacy+" conflicts with theme."+table+"."+field+"; canonical value wins")
			}
		}
		values[field] = n
	}
	resolveScale("texture", "scale", "texture_scale", 1)
	resolveMix("heading_texture", "color_mix", "heading_texture_strength", .45)
	resolveScale("heading_texture", "scale", "heading_texture_scale", 1)
	resolveMix("motif", "color_mix", "motif_color_distance", .01)
	resolveMix("motif", "row_offset", "motif_row_offset", .24)
	resolveMix("motif", "wobble", "motif_wobble", .18)
	resolveMix("motif", "scatter", "motif_scatter", 0)
	resolveLegacy := func(table, key, legacy string, fallback any) {
		values, valuesOK := out[table].(map[string]any)
		if !valuesOK {
			values = nil
		}
		if values == nil {
			values = map[string]any{}
			out[table] = values
		}
		if value, ok := values[key]; ok {
			if old, exists := raw[legacy]; exists && fmt.Sprint(value) != fmt.Sprint(old) {
				warnings = append(warnings, legacy+" conflicts with theme."+table+"."+key+"; canonical value wins")
			}
			return
		}
		if value, ok := raw[legacy]; ok {
			values[key] = value
		} else {
			values[key] = fallback
		}
	}
	resolveLegacy("heading_texture", "kind", "heading_texture", "inherit")
	resolveLegacy("motif", "kind", "motif", "off")
	resolveLegacy("motif", "glyph", "motif_glyph", "W")
	resolveLegacy("motif", "size", "motif_size", "78px")
	resolveLegacy("motif", "gap", "motif_gap", "10px")
	resolveLegacy("motif", "layer", "motif_layer", "sandwich")
	resolveLegacy("motif", "color", "motif_color", "ink")
	resolveLegacy("motif", "url", "motif_url", "")
	delete(out, "texture_strength")
	delete(out, "heading_texture_strength")
	delete(out, "motif_color_distance")
	return out, warnings
}

func Load() (Contract, error) {
	var c Contract
	data, err := contractFS.ReadFile("contract-v1.json")
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(data, &c)
	return c, err
}
func PaletteIDs() []string {
	c, err := Load()
	if err != nil {
		return nil
	}
	out := make([]string, len(c.Palettes))
	for i := range c.Palettes {
		out[i] = c.Palettes[i].ID
	}
	return out
}

// NormalizeMix accepts canonical numbers and legacy percentage strings.
func NormalizeMix(value any, fallback float64) float64 {
	return normalizeBounded(value, fallback, 0, 1)
}

func normalizeBounded(value any, fallback, minimum, maximum float64) float64 {
	var n float64
	switch v := value.(type) {
	case float64:
		n = v
	case int:
		n = float64(v)
	case int8:
		n = float64(v)
	case int16:
		n = float64(v)
	case int32:
		n = float64(v)
	case int64:
		n = float64(v)
	case uint:
		n = float64(v)
	case uint8:
		n = float64(v)
	case uint16:
		n = float64(v)
	case uint32:
		n = float64(v)
	case uint64:
		n = float64(v)
	case string:
		if _, err := fmt.Sscanf(v, "%f", &n); err != nil {
			return fallback
		}
		if v != "" && v[len(v)-1] == '%' {
			n /= 100
		}
	default:
		return fallback
	}
	return math.Round(math.Max(minimum, math.Min(maximum, n))*1000) / 1000
}

// NormalizeScale accepts canonical unitless scale values and legacy percentages.
func NormalizeScale(value any, fallback float64) float64 {
	return normalizeBounded(value, fallback, .25, 3)
}
