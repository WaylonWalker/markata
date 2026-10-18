package config

import (
	"encoding/json"
	"fmt"

	"github.com/WaylonWalker/markata/pkg/models"
)

// CanonicalJSON serializes a config without legacy flat theme aliases.
// Consumers can use the result as a portable presentation config.
func CanonicalJSON(config *models.Config) ([]byte, error) {
	raw, err := rawWrapperFromConfig(config)
	if err != nil {
		return nil, err
	}
	markata, ok := raw["markata-go"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("markata-go config section is missing")
	}
	theme, ok := markata["theme"].(map[string]any)
	if !ok {
		theme = map[string]any{}
	}
	for _, key := range []string{"aesthetic", "palette", "palette_light", "palette_dark", "fontpack", "texture", "texture_strength", "texture_scale", "texture_scope", "heading_texture", "heading_texture_strength", "heading_texture_scale", "motif", "motif_glyph", "motif_size", "motif_gap", "motif_row_offset", "motif_wobble", "motif_scatter", "motif_layer", "motif_color", "motif_color_distance"} {
		delete(markata, key)
	}
	// Rebuild the canonical dials explicitly. JSON struct tags use omitempty
	// elsewhere for historical config compatibility, but zero is a meaningful
	// contract value for every normalized color_mix/geometry dial.
	canonical := func(value any) map[string]any {
		result, ok := value.(map[string]any)
		if !ok {
			result = nil
		}
		if result == nil {
			result = map[string]any{}
		}
		return result
	}
	texture := canonical(theme["texture"])
	heading := canonical(theme["heading_texture"])
	motif := canonical(theme["motif"])
	texture["color_mix"] = config.Theme.Texture.ColorMix
	texture["scale"] = config.Theme.Texture.Scale
	texture["scope"] = config.Theme.Texture.Scope
	heading["color_mix"] = config.Theme.HeadingTexture.ColorMix
	heading["scale"] = config.Theme.HeadingTexture.Scale
	motif["row_offset"] = config.Theme.Motif.RowOffset
	motif["wobble"] = config.Theme.Motif.Wobble
	motif["scatter"] = config.Theme.Motif.Scatter
	motif["color_mix"] = config.Theme.Motif.ColorMix
	theme["texture"], theme["heading_texture"], theme["motif"] = texture, heading, motif
	markata["theme"] = theme
	return json.MarshalIndent(raw, "", "  ")
}
