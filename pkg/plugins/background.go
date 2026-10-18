// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// BackgroundPlugin adds configurable multi-layered background decorations to pages.
// It runs at the configure stage to validate configuration.
type BackgroundPlugin struct {
	config models.BackgroundConfig
}

// NewBackgroundPlugin creates a new BackgroundPlugin with default settings.
func NewBackgroundPlugin() *BackgroundPlugin {
	return &BackgroundPlugin{
		config: models.NewBackgroundConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *BackgroundPlugin) Name() string {
	return "background"
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "theme.background" key.
func (p *BackgroundPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Check for theme config first
	themeConfig, ok := config.Extra["theme"]
	if !ok {
		return nil
	}

	// Handle parsed ThemeConfig struct
	if tc, ok := themeConfig.(models.ThemeConfig); ok {
		p.config = tc.Background
		return p.validate()
	}

	// Handle map configuration (raw TOML)
	themeMap, ok := themeConfig.(map[string]interface{})
	if !ok {
		return nil
	}

	// Check for background config within theme
	bgConfig, ok := themeMap["background"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := bgConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = &enabled
		}

		if css, ok := cfgMap["css"].(string); ok {
			p.config.CSS = css
		}

		// Parse scripts array
		if scripts, ok := cfgMap["scripts"].([]interface{}); ok {
			p.config.Scripts = make([]string, 0, len(scripts))
			for _, s := range scripts {
				if str, ok := s.(string); ok {
					p.config.Scripts = append(p.config.Scripts, str)
				}
			}
		}

		// Parse backgrounds array
		if backgrounds, ok := cfgMap["backgrounds"].([]interface{}); ok {
			p.config.Backgrounds = make([]models.BackgroundElement, 0, len(backgrounds))
			for _, bg := range backgrounds {
				if bgMap, ok := bg.(map[string]interface{}); ok {
					element := models.BackgroundElement{}
					if html, ok := bgMap["html"].(string); ok {
						element.HTML = html
					}
					if zIndex, ok := parseIntFromInterface(bgMap["z_index"]); ok {
						element.ZIndex = zIndex
					}
					p.config.Backgrounds = append(p.config.Backgrounds, element)
				}
			}
		}
	}

	return p.validate()
}

// validate checks that the configuration is valid.
func (p *BackgroundPlugin) validate() error {
	if !p.config.IsEnabled() {
		return nil
	}

	// Validate that HTML in background elements doesn't contain dangerous content
	for i, bg := range p.config.Backgrounds {
		if bg.HTML == "" {
			return fmt.Errorf("background element %d has empty HTML", i)
		}
		// Basic validation - warn if script tags are embedded directly
		if strings.Contains(strings.ToLower(bg.HTML), "<script") {
			return fmt.Errorf("background element %d contains <script> tag; use the scripts array instead", i)
		}
	}

	return nil
}

// IsEnabled returns whether the background plugin is enabled.
func (p *BackgroundPlugin) IsEnabled() bool {
	return p.config.IsEnabled()
}

// Config returns the current background configuration.
func (p *BackgroundPlugin) Config() models.BackgroundConfig {
	return p.config
}

// SetConfig sets the background configuration directly.
// This is useful for testing or programmatic configuration.
func (p *BackgroundPlugin) SetConfig(config models.BackgroundConfig) {
	p.config = config
}

// GenerateBackgroundHTML generates the HTML for background elements.
// This is called by templates to inject background decorations.
func (p *BackgroundPlugin) GenerateBackgroundHTML() string {
	if !p.config.IsEnabled() || len(p.config.Backgrounds) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("<!-- Background Decorations -->\n")

	for _, bg := range p.config.Backgrounds {
		if bg.ZIndex != 0 {
			sb.WriteString(fmt.Sprintf(`<div class="background-layer" style="z-index: %d; position: fixed; inset: 0; pointer-events: none;">`, bg.ZIndex))
		} else {
			sb.WriteString(`<div class="background-layer" style="position: fixed; inset: 0; pointer-events: none; z-index: -1;">`)
		}
		sb.WriteString("\n  ")
		sb.WriteString(bg.HTML)
		sb.WriteString("\n</div>\n")
	}

	return sb.String()
}

// GenerateBackgroundCSS generates the CSS for background elements.
// This is called by templates to inject custom CSS.
func (p *BackgroundPlugin) GenerateBackgroundCSS() string {
	if !p.config.IsEnabled() || p.config.CSS == "" {
		return ""
	}

	return fmt.Sprintf("<style>\n%s\n</style>", p.config.CSS)
}

// GenerateBackgroundScripts generates the script tags for background elements.
// This is called by templates to inject scripts.
func (p *BackgroundPlugin) GenerateBackgroundScripts() string {
	if !p.config.IsEnabled() || len(p.config.Scripts) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("<!-- Background Scripts -->\n")

	for _, script := range p.config.Scripts {
		sb.WriteString(fmt.Sprintf(`<script src=%q></script>`, script))
		sb.WriteString("\n")
	}

	return sb.String()
}

// Ensure BackgroundPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*BackgroundPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*BackgroundPlugin)(nil)
)
