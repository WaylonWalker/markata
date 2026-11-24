// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"html"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// WikilinkHoverPlugin adds hover preview data attributes to wikilinks.
// It runs at the render stage (after wikilinks have been converted to HTML).
// This plugin finds <a class="wikilink"> tags and adds:
// - data-preview: truncated description/content for hover tooltips
// - data-preview-image: featured image URL if available
// - data-preview-screenshot: screenshot URL if service is configured
type WikilinkHoverPlugin struct {
	config  models.WikilinkHoverConfig
	postIdx *lifecycle.PostIndex
}

// NewWikilinkHoverPlugin creates a new WikilinkHoverPlugin with default settings.
func NewWikilinkHoverPlugin() *WikilinkHoverPlugin {
	return &WikilinkHoverPlugin{
		config: models.NewWikilinkHoverConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *WikilinkHoverPlugin) Name() string {
	return "wikilink_hover"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs late in render stage (after wikilinks have been converted).
func (p *WikilinkHoverPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate + 10 // Run after wikilinks transform
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "wikilink_hover" key.
func (p *WikilinkHoverPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Check for wikilink_hover config in Extra
	pluginConfig, ok := config.Extra["wikilink_hover"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := pluginConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = enabled
		}
		if previewLength, ok := parseIntFromInterface(cfgMap["preview_length"]); ok && previewLength > 0 {
			p.config.PreviewLength = previewLength
		}
		if includeImage, ok := cfgMap["include_image"].(bool); ok {
			p.config.IncludeImage = includeImage
		}
		if screenshotService, ok := cfgMap["screenshot_service"].(string); ok {
			p.config.ScreenshotService = screenshotService
		}
	}

	return nil
}

// Render processes wikilinks in all post HTML to add hover data attributes.
func (p *WikilinkHoverPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	// Use the shared PostIndex from the lifecycle manager
	p.postIdx = m.PostIndex()

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, `class="wikilink"`)
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// wikilinkAnchorRegex matches wikilink anchor tags created by the wikilinks plugin.
// Captures: full tag attributes, href value
var wikilinkAnchorRegex = regexp.MustCompile(
	`<a\s+([^>]*class="[^"]*wikilink[^"]*"[^>]*)>([^<]*)</a>`,
)

// wikilinkHrefRegex extracts href attribute value from tag attributes.
var wikilinkHrefRegex = regexp.MustCompile(`href="([^"]*)"`)

// processPost processes a single post's HTML to add hover data to wikilinks.
func (p *WikilinkHoverPlugin) processPost(post *models.Post) error {
	// Skip posts marked as skip or with no HTML content
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	// Quick check for wikilinks
	if !strings.Contains(post.ArticleHTML, `class="wikilink"`) {
		return nil
	}

	// Replace wikilink anchors with enhanced versions
	result := wikilinkAnchorRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		return p.enhanceWikilink(match)
	})

	post.ArticleHTML = result
	return nil
}

// enhanceWikilink adds data attributes to a wikilink anchor tag.
func (p *WikilinkHoverPlugin) enhanceWikilink(match string) string {
	// Extract href from the match
	hrefMatches := wikilinkHrefRegex.FindStringSubmatch(match)
	if len(hrefMatches) < 2 {
		return match
	}

	href := hrefMatches[1]

	// Look up the target post using the shared PostIndex
	targetPost := p.postIdx.ByHref[href]
	if targetPost == nil {
		// Try without trailing slash
		targetPost = p.postIdx.ByHref[strings.TrimSuffix(href, "/")]
	}
	if targetPost == nil {
		// Try with trailing slash
		targetPost = p.postIdx.ByHref[href+"/"]
	}
	if targetPost == nil {
		return match
	}

	// Skip hover data for private posts to prevent metadata leaks
	if targetPost.Private {
		return match
	}

	// Build data attributes
	attrs := p.buildDataAttributes(targetPost)
	if attrs == "" {
		return match
	}

	// Insert data attributes before the closing >
	// Find the position of the first > in the <a ...> tag
	tagEnd := strings.Index(match, ">")
	if tagEnd == -1 {
		return match
	}

	// Insert attributes before the >
	return match[:tagEnd] + " " + attrs + match[tagEnd:]
}

// buildDataAttributes creates the data attribute string for hover previews.
func (p *WikilinkHoverPlugin) buildDataAttributes(post *models.Post) string {
	var attrs []string

	// Get preview text from description or content
	previewText := p.getPreviewText(post)
	if previewText != "" {
		attrs = append(attrs, `data-preview="`+html.EscapeString(previewText)+`"`)
	}

	// Add preview image if configured and available
	if p.config.IncludeImage {
		imageURL := p.getPostImage(post)
		if imageURL != "" {
			attrs = append(attrs, `data-preview-image="`+html.EscapeString(imageURL)+`"`)
		}
	}

	// Add screenshot URL if service is configured
	if p.config.ScreenshotService != "" && post.Href != "" {
		screenshotURL := p.config.ScreenshotService + post.Href
		attrs = append(attrs, `data-preview-screenshot="`+html.EscapeString(screenshotURL)+`"`)
	}

	return strings.Join(attrs, " ")
}

// getPreviewText extracts preview text from post description or content.
func (p *WikilinkHoverPlugin) getPreviewText(post *models.Post) string {
	// First try description
	if post.Description != nil && *post.Description != "" {
		return truncatePreviewText(*post.Description, p.config.PreviewLength)
	}

	// Fall back to content (strip HTML if any)
	if post.Content != "" {
		// Use ArticleHTML if available (already rendered), otherwise raw content
		text := post.Content
		if post.ArticleHTML != "" {
			text = stripHTML(post.ArticleHTML)
		}
		return truncatePreviewText(text, p.config.PreviewLength)
	}

	return ""
}

// getPostImage finds a featured image URL for the post.
func (p *WikilinkHoverPlugin) getPostImage(post *models.Post) string {
	if post.Extra == nil {
		return ""
	}

	// Check common image field names
	imageFields := []string{"image", "featured_image", "cover_image", "og_image", "thumbnail"}
	for _, field := range imageFields {
		if img, ok := post.Extra[field].(string); ok && img != "" {
			return img
		}
	}

	return ""
}

// truncatePreviewText truncates text to maxLen characters, adding ellipsis if needed.
// It tries to break at word boundaries.
func truncatePreviewText(text string, maxLen int) string {
	// Clean up whitespace
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	// Collapse multiple spaces
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	if utf8.RuneCountInString(text) <= maxLen {
		return text
	}

	// Truncate to maxLen runes
	runes := []rune(text)
	truncated := string(runes[:maxLen])

	// Try to break at last space
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > maxLen/2 {
		truncated = truncated[:lastSpace]
	}

	return strings.TrimSpace(truncated) + "..."
}

// stripHTML removes HTML tags from text.
func stripHTML(s string) string {
	// Simple regex-based HTML stripping
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	result := tagRegex.ReplaceAllString(s, "")

	// Decode common HTML entities
	result = strings.ReplaceAll(result, "&nbsp;", " ")
	result = strings.ReplaceAll(result, "&amp;", "&")
	result = strings.ReplaceAll(result, "&lt;", "<")
	result = strings.ReplaceAll(result, "&gt;", ">")
	result = strings.ReplaceAll(result, "&quot;", `"`)
	result = strings.ReplaceAll(result, "&#39;", "'")

	return result
}

// SetConfig sets the plugin configuration directly.
// This is useful for testing or programmatic configuration.
func (p *WikilinkHoverPlugin) SetConfig(config models.WikilinkHoverConfig) {
	p.config = config
}

// Config returns the current plugin configuration.
func (p *WikilinkHoverPlugin) Config() models.WikilinkHoverConfig {
	return p.config
}

// Ensure WikilinkHoverPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*WikilinkHoverPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*WikilinkHoverPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*WikilinkHoverPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*WikilinkHoverPlugin)(nil)
)
