package plugins

import (
	"encoding/json"
	"log"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
)

// mermaidRenderer is an interface for rendering Mermaid diagrams in different modes
type mermaidRenderer interface {
	// render processes a mermaid diagram code and returns the rendered output
	// For client mode: returns the diagram code wrapped in mermaid markers
	// For cli/chromium modes: returns the SVG output
	render(diagramCode string) (string, error)
	// close cleans up any resources (e.g., browser connections)
	close() error
}

// mermaidPaletteColors holds resolved palette colors for Mermaid theming.
// These are resolved at build time from the site's palette and passed to
// mermaid.initialize() as themeVariables for consistent diagram styling.
type mermaidPaletteColors struct {
	Background string // bg-primary -> --color-background
	TextColor  string // text-primary -> --color-text
	Accent     string // accent -> --color-border fallback
	Surface    string // bg-surface -> --color-surface
}

// resolvePaletteColors loads and resolves palette colors for mermaid theming.
// Returns nil if no palette is configured or colors cannot be resolved.
func resolvePaletteColors(extra map[string]interface{}) *mermaidPaletteColors {
	if extra == nil {
		return nil
	}

	// Extract palette name from config.Extra["theme"]
	// The theme config may be a map[string]interface{} (raw TOML) or
	// models.ThemeConfig (if already parsed by another plugin).
	var paletteName string
	switch theme := extra["theme"].(type) {
	case map[string]interface{}:
		if p, ok := theme["palette"].(string); ok {
			paletteName = p
		}
	case models.ThemeConfig:
		paletteName = theme.Palette
	case *models.ThemeConfig:
		if theme != nil {
			paletteName = theme.Palette
		}
	}

	if paletteName == "" {
		return nil
	}

	loader := palettes.NewLoader()
	palette, err := loader.Load(paletteName)
	if err != nil {
		log.Printf("[mermaid] could not load palette %q: %v", paletteName, err)
		return nil
	}

	colors := &mermaidPaletteColors{
		Background: palette.Resolve("bg-primary"),
		TextColor:  palette.Resolve("text-primary"),
		Accent:     palette.Resolve("accent"),
		Surface:    palette.Resolve("bg-surface"),
	}

	// Only return if we got at least the accent color
	if colors.Accent == "" {
		log.Printf("[mermaid] palette %q has no accent color, using default theme", paletteName)
		return nil
	}

	log.Printf("[mermaid] resolved palette %q colors for theming (accent=%s, bg=%s, text=%s)",
		paletteName, colors.Accent, colors.Background, colors.TextColor)
	return colors
}

// Default fallback colors for mermaid theme variables when palette values are missing.
const (
	defaultMermaidBg      = "#ffffff"
	defaultMermaidText    = "#1f2937"
	defaultMermaidSurface = "#f9fafb"
)

// themeVariablesJSON returns the mermaid themeVariables object as a JSON string.
func (c *mermaidPaletteColors) themeVariablesJSON() string {
	bg := c.Background
	if bg == "" {
		bg = defaultMermaidBg
	}
	text := c.TextColor
	if text == "" {
		text = defaultMermaidText
	}
	surface := c.Surface
	if surface == "" {
		surface = defaultMermaidSurface
	}

	vars := map[string]interface{}{
		"background":          bg,
		"primaryColor":        surface,
		"primaryTextColor":    text,
		"primaryBorderColor":  c.Accent,
		"lineColor":           text,
		"textColor":           text,
		"nodeBkg":             surface,
		"nodeBorder":          c.Accent,
		"nodeTextColor":       text,
		"fontSize":            "16px",
		"clusterBkg":          bg,
		"clusterBorder":       c.Accent,
		"clusterTextColor":    text,
		"titleColor":          text,
		"edgeLabelBackground": surface,
	}

	b, err := json.Marshal(vars)
	if err != nil {
		log.Printf("[mermaid] failed to marshal theme variables: %v", err)
		return "{}"
	}
	return string(b)
}

// clientRenderer renders diagrams client-side via JavaScript
type clientRenderer struct {
	config models.MermaidConfig
}

func (r *clientRenderer) render(diagramCode string) (string, error) {
	// Client-side rendering: just return the diagram code, the script will render it
	return diagramCode, nil
}

func (r *clientRenderer) close() error {
	return nil
}

const mermaidJSVersion = "10.9.5"
const nativeMermaidVersion = "go-mermaid/v0.1.6"
