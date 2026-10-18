package lsp

import (
	"strings"

	"github.com/WaylonWalker/markata/pkg/renderingcontract"
)

const themePaletteDetail = "palette"

// ThemeCompletionItems exposes canonical theme values to configuration-aware clients.
// The values come from the versioned rendering contract, not a second enum list.
func ThemeCompletionItems(prefix string) []CompletionItem {
	c, err := renderingcontract.Load()
	if err != nil {
		return nil
	}
	values := renderingcontract.PaletteIDs()
	items := make([]CompletionItem, 0, len(values))
	for _, value := range values {
		if prefix != "" && !strings.HasPrefix(value, prefix) {
			continue
		}
		detail := themePaletteDetail
		for _, palette := range c.Palettes {
			if palette.ID == value {
				detail = palette.Family + "/" + palette.Variant
				break
			}
		}
		items = append(items, CompletionItem{Label: value, Kind: CompletionItemKindEnumMember, InsertText: value, Detail: detail})
	}
	return items
}

//nolint:gocyclo // Completion context parsing covers the supported config syntaxes.
func getThemeConfigCompletions(content string, lineNumber int, current string) []CompletionItem {
	lines := strings.Split(content, "\n")
	if lineNumber >= len(lines) {
		return nil
	}
	section := ""
	for i := 0; i <= lineNumber; i++ {
		trimmed := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			section = strings.Trim(trimmed, "[]")
		}
	}
	if section == "" {
		for i := lineNumber; i >= 0; i-- {
			trimmed := strings.TrimSpace(lines[i])
			switch {
			case strings.HasPrefix(trimmed, "motif:") || strings.HasPrefix(trimmed, `"motif"`):
				section = "markata-go.theme.motif"
			case strings.HasPrefix(trimmed, "heading_texture:") || strings.HasPrefix(trimmed, `"heading_texture"`):
				section = "markata-go.theme.heading_texture"
			case strings.HasPrefix(trimmed, "texture:") || strings.HasPrefix(trimmed, `"texture"`):
				section = "markata-go.theme.texture"
			case strings.HasPrefix(trimmed, "theme:") || strings.HasPrefix(trimmed, `"theme"`):
				section = "markata-go.theme"
			}
			if section != "" {
				break
			}
		}
	}
	if !strings.HasPrefix(section, "markata-go.theme") {
		if section != "markata-go" {
			return nil
		}
	}
	keyGroup := "theme"
	if strings.HasSuffix(section, ".variables") {
		common := []string{"--content-width", "--content-max-width", "--font-size"}
		prefix := strings.Trim(strings.TrimSpace(current), "= \"")
		items := make([]CompletionItem, 0, len(common)+1)
		for _, key := range common {
			if prefix == "" || strings.HasPrefix(key, prefix) {
				items = append(items, CompletionItem{Label: key, Kind: CompletionItemKindProperty, InsertText: key, Detail: "CSS custom property; CSS string"})
			}
		}
		items = append(items, CompletionItem{Label: "--foo", Kind: CompletionItemKindProperty, InsertText: "--foo", Detail: "arbitrary CSS custom property allowed; CSS string value"})
		return items
	}
	if strings.HasSuffix(section, ".texture") {
		keyGroup = "texture"
	}
	if strings.HasSuffix(section, ".heading_texture") {
		keyGroup = "heading_texture"
	}
	if strings.HasSuffix(section, ".motif") {
		keyGroup = "motif"
	}
	trimmed := strings.TrimSpace(current)
	if strings.Contains(trimmed, "=") || strings.Contains(trimmed, ":") {
		delimiter := "="
		if !strings.Contains(trimmed, "=") {
			delimiter = ":"
		}
		parts := strings.SplitN(trimmed, delimiter, 2)
		field := strings.Trim(strings.TrimSpace(parts[0]), "\"")
		prefix := strings.Trim(strings.TrimSpace(parts[1]), "\" ,")
		kindGroup := map[string]string{"texture": "textures", "heading_texture": "heading_textures", "motif": "motifs"}[keyGroup]
		if kindGroup == "" {
			kindGroup = "textures"
		}
		groups := map[string]string{"palette": "palette", "aesthetic": "aesthetics", "fontpack": "fontpacks", "kind": kindGroup, "scope": "scopes", "layer": "motif_layers", "color": "motif_colors"}
		if group, ok := groups[field]; ok {
			if group == "palette" {
				return ThemeCompletionItems(prefix)
			}
			c, err := renderingcontract.Load()
			if err != nil {
				return nil
			}
			items := make([]CompletionItem, 0)
			for _, value := range c.Enums[group] {
				if group == "fontpacks" && value == "system" {
					continue
				}
				if strings.HasPrefix(value, prefix) {
					items = append(items, CompletionItem{Label: value, Kind: CompletionItemKindEnumMember, InsertText: value, Detail: group})
				}
			}
			return items
		}
	}
	c, err := renderingcontract.Load()
	if err != nil {
		return nil
	}
	keys := append([]string{}, c.ConfigKeys[keyGroup]...)
	if section == "markata-go" {
		for alias, replacement := range c.Aliases {
			if strings.HasPrefix(replacement, "theme.") {
				keys = append(keys, alias)
			}
		}
	}
	prefix := strings.TrimSpace(strings.TrimSuffix(trimmed, "="))
	items := make([]CompletionItem, 0, len(keys))
	for _, key := range keys {
		if prefix == "" || strings.HasPrefix(key, prefix) {
			item := CompletionItem{Label: key, Kind: CompletionItemKindProperty, InsertText: key}
			if replacement, ok := c.Aliases[key]; ok {
				item.Detail = "deprecated; use " + replacement
			}
			items = append(items, item)
		}
	}
	return items
}
