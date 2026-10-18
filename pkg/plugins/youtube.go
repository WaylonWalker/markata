// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// YouTubePlugin converts YouTube URLs into embedded iframes.
// It runs at the render stage (after markdown conversion) and handles
// both plain URLs and autolinked URLs (from the Linkify extension).
type YouTubePlugin struct {
	config models.YouTubeConfig
}

// NewYouTubePlugin creates a new YouTubePlugin with default settings.
func NewYouTubePlugin() *YouTubePlugin {
	return &YouTubePlugin{
		config: models.NewYouTubeConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *YouTubePlugin) Name() string {
	return "youtube"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs after render_markdown (which has default priority 0).
func (p *YouTubePlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "youtube" key.
func (p *YouTubePlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Check for youtube config in Extra
	pluginConfig, ok := config.Extra["youtube"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := pluginConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = enabled
		}
		if privacyEnhanced, ok := cfgMap["privacy_enhanced"].(bool); ok {
			p.config.PrivacyEnhanced = privacyEnhanced
		}
		if containerClass, ok := cfgMap["container_class"].(string); ok && containerClass != "" {
			p.config.ContainerClass = containerClass
		}
		if lazyLoad, ok := cfgMap["lazy_load"].(bool); ok {
			p.config.LazyLoad = lazyLoad
		}
	}

	return nil
}

// Render processes YouTube URLs in the rendered HTML.
func (p *YouTubePlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, "youtube.com") || strings.Contains(post.ArticleHTML, "youtu.be")
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// youtubeURLRegex matches YouTube URLs in various formats.
// Matches both plain URLs in <p> tags AND autolinked URLs (from Linkify extension).
// Pattern 1: <p>URL</p> (plain URL, no Linkify)
// Pattern 2: <p><a href="URL">URL</a></p> (autolinked by Linkify extension)
var youtubeURLRegex = regexp.MustCompile(
	`<p>\s*(?:<a href="[^"]*">)?(https?://(?:www\.|m\.)?(?:youtube\.com/watch\?v=|youtu\.be/)([a-zA-Z0-9_-]{11})(?:[^<\s]*)?)\s*(?:</a>)?\s*</p>`,
)

// youtubeInCodeBlockRegex detects if a YouTube URL is inside a code block.
var youtubeInCodeBlockRegex = regexp.MustCompile(
	`<code[^>]*>.*?https?://(?:www\.|m\.)?(?:youtube\.com/watch\?v=|youtu\.be/)[a-zA-Z0-9_-]{11}.*?</code>`,
)

// processPost processes a single post's HTML for YouTube URLs.
func (p *YouTubePlugin) processPost(post *models.Post) error {
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	if !strings.Contains(post.ArticleHTML, "youtube.com") && !strings.Contains(post.ArticleHTML, "youtu.be") {
		return nil
	}

	result := youtubeURLRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		if youtubeInCodeBlockRegex.MatchString(post.ArticleHTML) {
			matchIdx := strings.Index(post.ArticleHTML, match)
			if matchIdx >= 0 {
				beforeMatch := post.ArticleHTML[:matchIdx]
				openCodes := strings.Count(beforeMatch, "<code")
				closeCodes := strings.Count(beforeMatch, "</code>")
				if openCodes > closeCodes {
					return match
				}
			}
		}

		submatches := youtubeURLRegex.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		fullURL := submatches[1]
		videoID := submatches[2]

		// Parse timestamp from URL
		startTime := p.parseTimestamp(fullURL)

		return p.buildEmbed(videoID, startTime)
	})

	post.ArticleHTML = result
	return nil
}

// parseTimestamp extracts and converts timestamp from YouTube URL.
// Supports formats: ?t=123, ?t=1h2m3s, ?t=2m30s, ?start=123
// Returns the start time in seconds, or 0 if no timestamp found.
func (p *YouTubePlugin) parseTimestamp(rawURL string) int {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return 0
	}

	query := parsedURL.Query()

	// Check for ?start= parameter first
	if start := query.Get("start"); start != "" {
		if seconds, err := strconv.Atoi(start); err == nil {
			return seconds
		}
	}

	// Check for ?t= parameter
	t := query.Get("t")
	if t == "" {
		return 0
	}

	// Try simple integer (seconds)
	if seconds, err := strconv.Atoi(t); err == nil {
		return seconds
	}

	// Parse duration format (1h2m3s)
	var total int

	// Parse hours
	if idx := strings.Index(t, "h"); idx != -1 {
		if hours, err := strconv.Atoi(t[:idx]); err == nil {
			total += hours * 3600
		}
		t = t[idx+1:]
	}

	// Parse minutes
	if idx := strings.Index(t, "m"); idx != -1 {
		if minutes, err := strconv.Atoi(t[:idx]); err == nil {
			total += minutes * 60
		}
		t = t[idx+1:]
	}

	// Parse seconds
	if idx := strings.Index(t, "s"); idx != -1 {
		if seconds, err := strconv.Atoi(t[:idx]); err == nil {
			total += seconds
		}
	}

	return total
}

// buildEmbed creates the HTML for a YouTube embed.
func (p *YouTubePlugin) buildEmbed(videoID string, startTime int) string {
	domain := "www.youtube.com"
	if p.config.PrivacyEnhanced {
		domain = "www.youtube-nocookie.com"
	}

	// Build embed URL with optional start time
	embedURL := fmt.Sprintf("https://%s/embed/%s", domain, html.EscapeString(videoID))
	if startTime > 0 {
		embedURL += fmt.Sprintf("?start=%d", startTime)
	}

	var sb strings.Builder
	sb.WriteString(`<div class="`)
	sb.WriteString(html.EscapeString(p.config.ContainerClass))
	sb.WriteString("\">\n")
	sb.WriteString(`  <iframe`)
	sb.WriteString("\n")
	sb.WriteString(`    src="`)
	sb.WriteString(embedURL)
	sb.WriteString("\"\n")
	sb.WriteString(`    title="YouTube video player"`)
	sb.WriteString("\n")
	sb.WriteString(`    frameborder="0"`)
	sb.WriteString("\n")
	sb.WriteString(`    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"`)
	sb.WriteString("\n")
	sb.WriteString(`    allowfullscreen`)

	if p.config.LazyLoad {
		sb.WriteString("\n")
		sb.WriteString(`    loading="lazy"`)
	}

	sb.WriteString(">\n")
	sb.WriteString(`  </iframe>`)
	sb.WriteString("\n")
	sb.WriteString(`</div>`)

	return sb.String()
}

// SetConfig sets the plugin configuration directly.
func (p *YouTubePlugin) SetConfig(config models.YouTubeConfig) {
	p.config = config
}

// Config returns the current plugin configuration.
func (p *YouTubePlugin) Config() models.YouTubeConfig {
	return p.config
}

var (
	_ lifecycle.Plugin          = (*YouTubePlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*YouTubePlugin)(nil)
	_ lifecycle.RenderPlugin    = (*YouTubePlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*YouTubePlugin)(nil)
)
