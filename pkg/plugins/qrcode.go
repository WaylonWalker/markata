// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	qrcode "github.com/skip2/go-qrcode"
)

// QR code format constants.
const (
	qrFormatSVG = "svg"
	qrFormatPNG = "png"
)

// QRCodePlugin generates QR code images for each post's URL.
// It runs at the write stage to generate QR code files.
type QRCodePlugin struct {
	config    models.QRCodeConfig
	outputDir string
	baseURL   string
}

// NewQRCodePlugin creates a new QRCodePlugin with default settings.
func NewQRCodePlugin() *QRCodePlugin {
	return &QRCodePlugin{
		config: models.NewQRCodeConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *QRCodePlugin) Name() string {
	return "qrcode"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs late in write stage to ensure all post data is finalized.
func (p *QRCodePlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "qrcode" key.
func (p *QRCodePlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	p.outputDir = config.OutputDir

	if config.Extra == nil {
		return nil
	}

	// Get base URL from Extra (where the full models.Config is stored)
	if url, ok := config.Extra["url"].(string); ok {
		p.baseURL = url
	}

	// Check for qrcode config in Extra
	pluginConfig, ok := config.Extra["qrcode"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := pluginConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = enabled
		}
		if format, ok := cfgMap["format"].(string); ok && (format == qrFormatSVG || format == qrFormatPNG) {
			p.config.Format = format
		}
		if size, ok := parseIntFromInterface(cfgMap["size"]); ok && size > 0 {
			p.config.Size = size
		}
		if outputDir, ok := cfgMap["output_dir"].(string); ok && outputDir != "" {
			p.config.OutputDir = outputDir
		}
		if errorCorrection, ok := cfgMap["error_correction"].(string); ok {
			p.config.ErrorCorrection = errorCorrection
		}
		if foreground, ok := cfgMap["foreground"].(string); ok && foreground != "" {
			p.config.Foreground = foreground
		}
		if background, ok := cfgMap["background"].(string); ok && background != "" {
			p.config.Background = background
		}
	}

	return nil
}

// Write generates QR code images for all posts.
func (p *QRCodePlugin) Write(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	// Create QR code output directory
	qrDir := filepath.Join(p.outputDir, p.config.OutputDir)
	if err := os.MkdirAll(qrDir, 0o755); err != nil {
		return fmt.Errorf("failed to create QR code directory: %w", err)
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// processPost generates a QR code for a single post.
func (p *QRCodePlugin) processPost(post *models.Post) error {
	// Skip posts marked as skip
	if post.Skip {
		return nil
	}

	// Build the absolute URL for the post
	absoluteURL := p.buildAbsoluteURL(post)
	if absoluteURL == "" {
		return nil
	}

	// Generate QR code filename
	filename := p.buildFilename(post)
	qrPath := filepath.Join(p.outputDir, p.config.OutputDir, filename)

	// Get error correction level
	level := p.getErrorCorrectionLevel()

	// Generate the QR code
	var err error
	if p.config.Format == "png" {
		err = qrcode.WriteFile(absoluteURL, level, p.config.Size, qrPath)
	} else {
		// SVG format
		err = p.writeSVG(absoluteURL, level, qrPath)
	}

	if err != nil {
		return fmt.Errorf("failed to generate QR code for %s: %w", post.Slug, err)
	}

	// Add qrcode_url to post's Extra map
	qrURL := "/" + p.config.OutputDir + "/" + filename
	post.Set("qrcode_url", qrURL)

	return nil
}

// buildAbsoluteURL constructs the absolute URL for a post.
func (p *QRCodePlugin) buildAbsoluteURL(post *models.Post) string {
	href := post.Href
	if href == "" && post.Slug != "" {
		href = "/" + post.Slug + "/"
	}
	if href == "" {
		return ""
	}

	// Combine with base URL
	if p.baseURL != "" {
		baseURL := strings.TrimSuffix(p.baseURL, "/")
		return baseURL + href
	}

	// If no base URL, just return the relative path
	return href
}

// buildFilename generates the QR code filename for a post.
func (p *QRCodePlugin) buildFilename(post *models.Post) string {
	slug := post.Slug
	if slug == "" {
		// Fallback to path-based name
		slug = strings.TrimSuffix(filepath.Base(post.Path), filepath.Ext(post.Path))
	}
	return slug + "." + p.config.Format
}

// getErrorCorrectionLevel returns the QR code error correction level.
func (p *QRCodePlugin) getErrorCorrectionLevel() qrcode.RecoveryLevel {
	switch strings.ToUpper(p.config.ErrorCorrection) {
	case "L":
		return qrcode.Low
	case "Q":
		return qrcode.High
	case "H":
		return qrcode.Highest
	default: // "M" or anything else
		return qrcode.Medium
	}
}

// writeSVG generates an SVG QR code file.
func (p *QRCodePlugin) writeSVG(content string, level qrcode.RecoveryLevel, path string) error {
	// Generate QR code data
	qr, err := qrcode.New(content, level)
	if err != nil {
		return err
	}

	// Get the QR code as a 2D bitmap
	bitmap := qr.Bitmap()
	size := len(bitmap)

	// Calculate module size for the target dimensions
	moduleSize := p.config.Size / size
	if moduleSize < 1 {
		moduleSize = 1
	}
	actualSize := moduleSize * size

	// Build SVG
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" version="1.1" viewBox="0 0 %d %d" width="%d" height="%d">
`, actualSize, actualSize, actualSize, actualSize))

	// Background
	sb.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="%s"/>
`, actualSize, actualSize, p.config.Background))

	// Draw modules
	sb.WriteString(fmt.Sprintf(`<g fill="%s">
`, p.config.Foreground))

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if bitmap[y][x] {
				sb.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="%d" height="%d"/>
`, x*moduleSize, y*moduleSize, moduleSize, moduleSize))
			}
		}
	}

	sb.WriteString(`</g>
</svg>`)

	// Write to file (0o644 is appropriate for public web content)
	return os.WriteFile(path, []byte(sb.String()), 0o644) //nolint:gosec // QR codes are public web content
}

// SetConfig sets the plugin configuration directly.
// This is useful for testing or programmatic configuration.
func (p *QRCodePlugin) SetConfig(config models.QRCodeConfig) {
	p.config = config
}

// Config returns the current plugin configuration.
func (p *QRCodePlugin) Config() models.QRCodeConfig {
	return p.config
}

// SetOutputDir sets the output directory for testing.
func (p *QRCodePlugin) SetOutputDir(dir string) {
	p.outputDir = dir
}

// SetBaseURL sets the base URL for testing.
func (p *QRCodePlugin) SetBaseURL(url string) {
	p.baseURL = url
}

// Ensure QRCodePlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*QRCodePlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*QRCodePlugin)(nil)
	_ lifecycle.WritePlugin     = (*QRCodePlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*QRCodePlugin)(nil)
)
