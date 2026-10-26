// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"html"
	"net/url"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// OneLineLinkPlugin expands URLs that appear alone in a paragraph into rich preview cards.
// It runs at the render stage (after markdown conversion).
type OneLineLinkPlugin struct {
	config          models.OneLineLinkConfig
	excludePatterns []*regexp.Regexp
}

// NewOneLineLinkPlugin creates a new OneLineLinkPlugin with default settings.
func NewOneLineLinkPlugin() *OneLineLinkPlugin {
	return &OneLineLinkPlugin{
		config: models.NewOneLineLinkConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *OneLineLinkPlugin) Name() string {
	return "one_line_link"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs after render_markdown (which has default priority 0).
func (p *OneLineLinkPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "one_line_link" key.
func (p *OneLineLinkPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Check for one_line_link config in Extra
	pluginConfig, ok := config.Extra["one_line_link"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := pluginConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = enabled
		}
		if cardClass, ok := cfgMap["card_class"].(string); ok && cardClass != "" {
			p.config.CardClass = cardClass
		}
		if fetchMetadata, ok := cfgMap["fetch_metadata"].(bool); ok {
			p.config.FetchMetadata = fetchMetadata
		}
		if fallbackTitle, ok := cfgMap["fallback_title"].(string); ok && fallbackTitle != "" {
			p.config.FallbackTitle = fallbackTitle
		}
		if timeout, ok := parseIntFromInterface(cfgMap["timeout"]); ok && timeout > 0 {
			p.config.Timeout = timeout
		}
		if patterns, ok := cfgMap["exclude_patterns"].([]interface{}); ok {
			for _, pattern := range patterns {
				if str, ok := pattern.(string); ok {
					p.config.ExcludePatterns = append(p.config.ExcludePatterns, str)
				}
			}
		}
	}

	// Compile exclude patterns
	p.excludePatterns = make([]*regexp.Regexp, 0, len(p.config.ExcludePatterns))
	for _, pattern := range p.config.ExcludePatterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			p.excludePatterns = append(p.excludePatterns, re)
		}
	}

	return nil
}

// Render processes paragraphs containing only URLs in the rendered HTML.
func (p *OneLineLinkPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, "http://") || strings.Contains(post.ArticleHTML, "https://")
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// oneLineLinkRegex matches paragraphs containing only a URL.
// This matches <p>https://...</p> or <p>http://...</p> patterns.
var oneLineLinkRegex = regexp.MustCompile(
	`<p>\s*(https?://[^\s<>"]+)\s*</p>`,
)

// processPost processes a single post's HTML for standalone URL paragraphs.
func (p *OneLineLinkPlugin) processPost(post *models.Post) error {
	// Skip posts marked as skip or with no HTML content
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	// Quick check for any potential matches
	if !strings.Contains(post.ArticleHTML, "http://") && !strings.Contains(post.ArticleHTML, "https://") {
		return nil
	}

	// Replace standalone URL paragraphs with link cards
	result := oneLineLinkRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		// Extract the URL
		submatches := oneLineLinkRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		rawURL := strings.TrimSpace(submatches[1])

		// Check if URL should be excluded
		if p.isExcluded(rawURL) {
			return match
		}

		// Parse URL to extract domain
		parsedURL, err := url.Parse(rawURL)
		if err != nil {
			return match
		}

		// Build the link card
		return p.buildLinkCard(rawURL, parsedURL)
	})

	post.ArticleHTML = result
	return nil
}

// isExcluded checks if a URL matches any exclude pattern.
func (p *OneLineLinkPlugin) isExcluded(urlStr string) bool {
	for _, re := range p.excludePatterns {
		if re.MatchString(urlStr) {
			return true
		}
	}
	return false
}

// buildLinkCard creates the HTML for a link preview card.
func (p *OneLineLinkPlugin) buildLinkCard(rawURL string, parsedURL *url.URL) string {
	domain := parsedURL.Host
	// Remove www. prefix for cleaner display
	domain = strings.TrimPrefix(domain, "www.")

	title := p.config.FallbackTitle
	description := ""

	// Note: Metadata fetching is disabled by default for performance.
	// When enabled, it would fetch the page and extract og:title, og:description, og:image.
	// This is left as a future enhancement to avoid blocking builds on network requests.

	// Build the card HTML
	var sb strings.Builder
	sb.WriteString(`<a href="`)
	sb.WriteString(html.EscapeString(rawURL))
	sb.WriteString(`" class="`)
	sb.WriteString(html.EscapeString(p.config.CardClass))
	sb.WriteString(`" target="_blank" rel="noopener noreferrer" aria-label="`)
	sb.WriteString(html.EscapeString(title))
	sb.WriteString(` on `)
	sb.WriteString(html.EscapeString(domain))
	sb.WriteString(`">`)
	sb.WriteString("\n")

	sb.WriteString(`  <div class="`)
	sb.WriteString(html.EscapeString(p.config.CardClass))
	sb.WriteString(`-content">`)
	sb.WriteString("\n")

	sb.WriteString(`    <div class="`)
	sb.WriteString(html.EscapeString(p.config.CardClass))
	sb.WriteString(`-title">`)
	sb.WriteString(html.EscapeString(title))
	sb.WriteString(`</div>`)
	sb.WriteString("\n")

	if description != "" {
		sb.WriteString(`    <div class="`)
		sb.WriteString(html.EscapeString(p.config.CardClass))
		sb.WriteString(`-description">`)
		sb.WriteString(html.EscapeString(description))
		sb.WriteString(`</div>`)
		sb.WriteString("\n")
	}

	sb.WriteString(`    <div class="`)
	sb.WriteString(html.EscapeString(p.config.CardClass))
	sb.WriteString(`-url">`)
	sb.WriteString(html.EscapeString(domain))
	sb.WriteString(`</div>`)
	sb.WriteString("\n")

	sb.WriteString(`  </div>`)
	sb.WriteString("\n")
	sb.WriteString(`</a>`)

	return sb.String()
}

// SetConfig sets the plugin configuration directly.
// This is useful for testing or programmatic configuration.
func (p *OneLineLinkPlugin) SetConfig(config models.OneLineLinkConfig) {
	p.config = config
	// Recompile exclude patterns
	p.excludePatterns = make([]*regexp.Regexp, 0, len(config.ExcludePatterns))
	for _, pattern := range config.ExcludePatterns {
		re, err := regexp.Compile(pattern)
		if err == nil {
			p.excludePatterns = append(p.excludePatterns, re)
		}
	}
}

// Config returns the current plugin configuration.
func (p *OneLineLinkPlugin) Config() models.OneLineLinkConfig {
	return p.config
}

// Ensure OneLineLinkPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*OneLineLinkPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*OneLineLinkPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*OneLineLinkPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*OneLineLinkPlugin)(nil)
)
