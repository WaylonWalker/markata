// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"html"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// MDVideoPlugin converts markdown image syntax for video files into HTML video elements.
// It runs at the render stage (late priority, after markdown conversion).
//
// Example input (markdown):
//
//	![Video description](video.mp4)
//
// Example output (HTML):
//
//	<video loop muted playsinline controls preload="none" class="md-video" style="margin-inline:auto">
//	  <source src="video.mp4" type="video/mp4">
//	  Your browser does not support the video tag.
//	</video>
const mdVideoPreloadNone = "none"

type MDVideoPlugin struct {
	config models.MDVideoConfig
}

// NewMDVideoPlugin creates a new MDVideoPlugin with default settings.
func NewMDVideoPlugin() *MDVideoPlugin {
	return &MDVideoPlugin{
		config: models.NewMDVideoConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *MDVideoPlugin) Name() string {
	return "md_video"
}

// Priority returns the plugin's priority for a given stage.
// This plugin runs after render_markdown (which has default priority 0).
func (p *MDVideoPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected under the "md_video" key.
func (p *MDVideoPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Check for md_video config in Extra
	mdVideoConfig, ok := config.Extra["md_video"]
	if !ok {
		return nil
	}

	// Handle map configuration
	if cfgMap, ok := mdVideoConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.config.Enabled = enabled
		}
		if videoClass, ok := cfgMap["video_class"].(string); ok && videoClass != "" {
			p.config.VideoClass = videoClass
		}
		if controls, ok := cfgMap["controls"].(bool); ok {
			p.config.Controls = controls
		}
		if autoplay, ok := cfgMap["autoplay"].(bool); ok {
			p.config.Autoplay = autoplay
		}
		if loop, ok := cfgMap["loop"].(bool); ok {
			p.config.Loop = loop
		}
		if muted, ok := cfgMap["muted"].(bool); ok {
			p.config.Muted = muted
		}
		if playsinline, ok := cfgMap["playsinline"].(bool); ok {
			p.config.Playsinline = playsinline
		}
		if preload, ok := cfgMap["preload"].(string); ok && preload != "" {
			p.config.Preload = preload
		}

		// Handle video_extensions as []interface{} or []string
		switch extensions := cfgMap["video_extensions"].(type) {
		case []interface{}:
			p.config.VideoExtensions = make([]string, 0, len(extensions))
			for _, ext := range extensions {
				if s, ok := ext.(string); ok {
					p.config.VideoExtensions = append(p.config.VideoExtensions, s)
				}
			}
		case []string:
			p.config.VideoExtensions = extensions
		}
	}

	return nil
}

// Render processes video image tags in the rendered HTML for all posts.
func (p *MDVideoPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, "<img") || strings.Contains(strings.ToLower(post.ArticleHTML), "<video")
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// imgTagRegex matches <img> tags and captures src and alt attributes.
// It handles both src="..." and alt="..." in either order.
var imgTagRegex = regexp.MustCompile(`<img\s+([^>]*)>`)

var videoTagRegex = regexp.MustCompile(`(?i)<video\b[^>]*>`)

var videoBlockRegex = regexp.MustCompile(`(?is)<video\b[^>]*>.*?</video\s*>`)

var videoAutoplayAttrRegex = regexp.MustCompile(`(?i)\s+autoplay(?:\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+))?`)

var videoPreloadAttrRegex = regexp.MustCompile(`(?i)\s+preload\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)`)

var videoGeometryAttrRegex = regexp.MustCompile(`(?i)\s+(?:height|style)\s*=`)

var videoPosterAttrRegex = regexp.MustCompile(`(?i)\s+poster\s*=`)

var videoSourceRegex = regexp.MustCompile(`(?i)<source\b[^>]*\bsrc\s*=\s*(?:"([^"]+)"|'([^']+)')[^>]*>`)

var videoSrcRegex = regexp.MustCompile(`(?i)\bsrc\s*=\s*(?:"([^"]+)"|'([^']+)')`)

// srcAttrRegex extracts the src attribute value.
var srcAttrRegex = regexp.MustCompile(`src="([^"]*)"`)

// altAttrRegex extracts the alt attribute value.
var altAttrRegex = regexp.MustCompile(`alt="([^"]*)"`)

// processPost processes a single post's HTML for video image tags.
func (p *MDVideoPlugin) processPost(post *models.Post) error {
	// Skip posts marked as skip or with no HTML content
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	// Normalize raw video tags as well as Markdown image syntax converted below.
	if strings.Contains(strings.ToLower(post.ArticleHTML), "<video") {
		post.ArticleHTML = videoBlockRegex.ReplaceAllStringFunc(post.ArticleHTML, func(block string) string {
			tag := videoTagRegex.FindString(block)
			if tag == "" {
				return block
			}
			normalized := p.normalizeVideoTag(post, tag, block)
			return strings.Replace(block, tag, normalized, 1)
		})
	}

	if !strings.Contains(post.ArticleHTML, "<img") {
		return nil
	}

	// Replace img tags that have video extensions
	result := imgTagRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		// Extract src attribute
		srcMatch := srcAttrRegex.FindStringSubmatch(match)
		if len(srcMatch) < 2 {
			return match // No src found, leave unchanged
		}
		src := srcMatch[1]

		// Check if src ends with a video extension
		if !p.isVideoURL(src) {
			return match // Not a video, leave unchanged
		}

		// Extract alt attribute (optional)
		alt := ""
		altMatch := altAttrRegex.FindStringSubmatch(match)
		if len(altMatch) >= 2 {
			alt = altMatch[1]
		}

		// Build the video tag
		return p.buildVideoTag(post, src, alt)
	})

	post.ArticleHTML = result
	return nil
}

func (p *MDVideoPlugin) normalizeVideoTag(post *models.Post, tag, block string) string {
	if videoAutoplayAttrRegex.MatchString(tag) {
		tag = videoAutoplayAttrRegex.ReplaceAllString(tag, "")
		tag = strings.TrimSuffix(tag, ">") + ` data-authored-autoplay="true">`
	}

	preload := p.config.Preload
	if strings.Contains(tag, `data-authored-autoplay="true"`) {
		preload = mdVideoPreloadNone
	}
	if preload != mdVideoPreloadNone && preload != "metadata" && preload != "auto" {
		preload = mdVideoPreloadNone
	}
	tag = videoPreloadAttrRegex.ReplaceAllString(tag, "")
	tag = strings.TrimSuffix(tag, ">") + ` preload="` + preload + `">`
	if !videoPosterAttrRegex.MatchString(tag) {
		source := videoSourceRegex.FindStringSubmatch(block)
		if len(source) == 0 {
			openTag := strings.SplitN(tag, ">", 2)[0]
			source = videoSrcRegex.FindStringSubmatch(openTag)
		}
		if len(source) > 0 {
			mediaURL := source[1]
			if mediaURL == "" && len(source) > 2 {
				mediaURL = source[2]
			}
			if poster := templates.PosterURLFromMap(templates.GetPostMap(post), html.UnescapeString(mediaURL)); poster != "" {
				poster = templates.WithSize(poster, 1200, 675)
				if !videoGeometryAttrRegex.MatchString(tag) {
					tag = strings.TrimSuffix(tag, ">") + ` style="aspect-ratio:16/9">`
				}
				tag = strings.TrimSuffix(tag, ">") + ` poster="` + html.EscapeString(poster) + `">`
			}
		}
	}
	return tag
}

// isVideoURL checks if a URL ends with a recognized video extension.
// It handles URLs with query parameters.
func (p *MDVideoPlugin) isVideoURL(url string) bool {
	urlLower := videoExtension(url)
	for _, ext := range p.config.VideoExtensions {
		if strings.EqualFold(urlLower, ext) {
			return true
		}
	}
	return false
}

// getVideoMIMEType returns the MIME type for a video URL based on extension.
func (p *MDVideoPlugin) getVideoMIMEType(url string) string {
	if mimeType := templates.VideoMIMEType(url); mimeType != "" {
		return mimeType
	}
	return "video/mp4" // Default fallback
}

func videoExtension(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return strings.ToLower(filepath.Ext(raw))
	}
	return strings.ToLower(filepath.Ext(u.Path))
}

// buildVideoTag constructs the HTML video element.
func (p *MDVideoPlugin) buildVideoTag(post *models.Post, src, alt string) string {
	var attrs []string

	// Add boolean attributes (order matters for consistency)
	if p.config.Autoplay {
		attrs = append(attrs, `data-authored-autoplay="true"`)
	}
	if p.config.Loop {
		attrs = append(attrs, "loop")
	}
	if p.config.Muted {
		attrs = append(attrs, "muted")
	}
	if p.config.Playsinline {
		attrs = append(attrs, "playsinline")
	}
	if p.config.Controls {
		attrs = append(attrs, "controls")
	}

	// Add preload attribute if not empty
	if p.config.Autoplay {
		attrs = append(attrs, `preload="none"`)
	} else if p.config.Preload != "" {
		attrs = append(attrs, `preload="`+p.config.Preload+`"`)
	}

	// Add class if configured
	if p.config.VideoClass != "" {
		attrs = append(attrs, `class="`+p.config.VideoClass+`"`)
	}

	// Markdown images are centered when they are standalone block media. Preserve
	// that layout when md_video replaces the image with a video, including when a
	// Dropper query parameter gives the video a narrower intrinsic width.
	style := "margin-inline:auto"

	if poster := templates.PosterURLFromMap(templates.GetPostMap(post), src); poster != "" {
		poster = templates.WithSize(poster, 1200, 675)
		style += ";aspect-ratio:16/9"
		attrs = append(attrs, `poster="`+poster+`"`)
	}

	attrs = append(attrs, `style="`+style+`"`)

	// Build the opening tag
	attrStr := ""
	if len(attrs) > 0 {
		attrStr = " " + strings.Join(attrs, " ")
	}

	// Get MIME type for the source
	mimeType := p.getVideoMIMEType(src)

	// Build the video element
	var sb strings.Builder
	sb.WriteString("<video")
	sb.WriteString(attrStr)
	sb.WriteString(">")
	sb.WriteString(`<source src="`)
	sb.WriteString(src)
	sb.WriteString(`" type="`)
	sb.WriteString(mimeType)
	sb.WriteString(`">`)

	// Add fallback text (use alt if provided, otherwise generic message)
	if alt != "" {
		sb.WriteString(alt)
	} else {
		sb.WriteString("Your browser does not support the video tag.")
	}

	sb.WriteString("</video>")
	return sb.String()
}

// SetConfig sets the md_video configuration directly.
// This is useful for testing or programmatic configuration.
func (p *MDVideoPlugin) SetConfig(config models.MDVideoConfig) {
	p.config = config
}

// Config returns the current md_video configuration.
func (p *MDVideoPlugin) Config() models.MDVideoConfig {
	return p.config
}

// Ensure MDVideoPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*MDVideoPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*MDVideoPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*MDVideoPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*MDVideoPlugin)(nil)
)
