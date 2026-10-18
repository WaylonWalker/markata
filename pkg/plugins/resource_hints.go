// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/resourcehints"
)

// ResourceHintsPlugin generates and injects resource hints (preconnect, dns-prefetch, etc.)
// into HTML files for improved network performance.
//
// The plugin runs during the Write stage with late priority (after HTML files are written)
// to scan generated HTML for external resources and inject appropriate resource hints.
type ResourceHintsPlugin struct {
	// enabled controls whether the plugin processes files
	enabled bool

	// autoDetect enables automatic detection of external domains
	autoDetect bool

	// config holds the resource hints configuration
	config *models.ResourceHintsConfig

	// detector detects external domains in HTML content
	detector *resourcehints.Detector

	// generator generates hint HTML tags
	generator *resourcehints.Generator
}

// NewResourceHintsPlugin creates a new ResourceHintsPlugin with default settings.
func NewResourceHintsPlugin() *ResourceHintsPlugin {
	return &ResourceHintsPlugin{
		enabled:    true,
		autoDetect: true,
		detector:   resourcehints.NewDetector(),
		generator:  resourcehints.NewGenerator(),
	}
}

// Name returns the unique name of the plugin.
func (p *ResourceHintsPlugin) Name() string {
	return "resource_hints"
}

// Priority returns the plugin's priority for the write stage.
// Returns a late priority to run after HTML files are written.
func (p *ResourceHintsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		// Run late in write stage, after publish_html
		return lifecycle.PriorityLate + 100
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options from config.ResourceHints.
func (p *ResourceHintsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Get ResourceHints config from lifecycle config
	rhConfig := getResourceHintsConfig(config)
	p.config = &rhConfig

	// Apply configuration
	p.enabled = rhConfig.IsEnabled()
	p.autoDetect = rhConfig.IsAutoDetectEnabled()

	// Set excluded domains on detector
	if len(rhConfig.ExcludeDomains) > 0 {
		p.detector.SetExcludeDomains(rhConfig.ExcludeDomains)
	}

	return nil
}

// Write scans HTML files for external resources and injects resource hints.
// Each page gets hints only for domains detected on that specific page.
func (p *ResourceHintsPlugin) Write(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}

	config := m.Config()
	outputDir := config.OutputDir

	// Process each HTML file individually for page-specific hints
	return filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-HTML files
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".html") {
			return nil
		}

		// Read HTML content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		htmlContent := string(content)

		// Skip files that don't have a <head> tag
		if !strings.Contains(htmlContent, "<head") {
			return nil
		}

		// Skip files that already have resource hints
		if strings.Contains(htmlContent, "<!-- Auto-generated resource hints -->") {
			return nil
		}

		// Detect external domains for THIS page only
		var detectedDomains []resourcehints.DetectedDomain
		if p.autoDetect {
			detectedDomains = p.detector.DetectExternalDomains(htmlContent)
		}

		// Generate hint tags for this page
		hintTags := p.generator.GenerateFromConfig(p.config, detectedDomains)
		if hintTags == "" {
			return nil // No hints to inject for this page
		}

		// Wrap in comment
		hintBlock := resourcehints.GenerateComment(hintTags)

		// Inject hints after <head> or after <meta charset>
		modifiedContent := p.injectHints(htmlContent, hintBlock)
		if modifiedContent == htmlContent {
			return nil // No changes made
		}

		// Write modified content back
		//nolint:gosec // G306: HTML files need 0644 for web serving
		return os.WriteFile(path, []byte(modifiedContent), 0o644)
	})
}

// headOpenRegex matches the opening <head> tag.
var headOpenRegex = regexp.MustCompile(`(?i)(<head[^>]*>)`)

// charsetMetaRegex matches a charset meta tag.
var charsetMetaRegex = regexp.MustCompile(`(?i)(<meta[^>]*charset[^>]*>)`)

// injectHints injects resource hint tags into HTML content.
// Prioritizes injection after charset meta tag, falls back to after <head>.
func (p *ResourceHintsPlugin) injectHints(htmlContent, hintBlock string) string {
	// Try to inject after charset meta tag (preferred position)
	if charsetMetaRegex.MatchString(htmlContent) {
		return charsetMetaRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
			return match + "\n\n" + hintBlock
		})
	}

	// Fall back to injecting after <head>
	if headOpenRegex.MatchString(htmlContent) {
		return headOpenRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
			return match + "\n" + hintBlock
		})
	}

	// No suitable injection point found
	return htmlContent
}

// getResourceHintsConfig extracts ResourceHintsConfig from lifecycle.Config.
func getResourceHintsConfig(config *lifecycle.Config) models.ResourceHintsConfig {
	// Check if we have a ResourceHints field in Extra
	if config.Extra != nil {
		if rh, ok := config.Extra["resource_hints"].(models.ResourceHintsConfig); ok {
			return rh
		}
	}
	return models.NewResourceHintsConfig()
}

// Ensure ResourceHintsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*ResourceHintsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ResourceHintsPlugin)(nil)
	_ lifecycle.WritePlugin     = (*ResourceHintsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*ResourceHintsPlugin)(nil)
)
