// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"html"
	"log"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// HashtagTagsPlugin transforms #tag syntax into HTML links with tag statistics.
// The hashtag references tags but does not modify the post's actual tags.
type HashtagTagsPlugin struct {
	cssClass string
}

// NewHashtagTagsPlugin creates a new HashtagTagsPlugin.
func NewHashtagTagsPlugin() *HashtagTagsPlugin {
	return &HashtagTagsPlugin{
		cssClass: "hashtag-tag",
	}
}

// Name returns the unique name of the plugin.
func (p *HashtagTagsPlugin) Name() string {
	return "hashtag_tags"
}

// Configure reads configuration options for the plugin.
func (p *HashtagTagsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Check for custom CSS class in config
	if config.Extra != nil {
		if cssClass, ok := config.Extra["hashtag_tags_css_class"].(string); ok && cssClass != "" {
			p.cssClass = cssClass
		}
	}
	return nil
}

// Priority returns the plugin's priority for a given stage.
// Runs late in transform stage, after markdown rendering but before other text transforms.
func (p *HashtagTagsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Transform processes #tag references in all post content.
func (p *HashtagTagsPlugin) Transform(m *lifecycle.Manager) error {
	// Fast incremental mode: skip when no feed-relevant changes
	if config := m.Config(); config.Extra != nil {
		if incremental, ok := config.Extra["feeds_incremental"].(bool); ok && incremental {
			if cache := GetBuildCache(m); cache != nil {
				if len(cache.GetChangedFeedSlugs()) == 0 {
					return nil
				}
			}
		}
	}

	// Build tag statistics map from all published posts
	tagStats := p.buildTagStats(m)

	if len(tagStats) == 0 {
		log.Printf("[hashtag_tags] No tags found, skipping")
		return nil
	}

	log.Printf("[hashtag_tags] Found %d tags", len(tagStats))

	// Cache tag stats for other plugins
	m.Cache().Set("hashtag_tag_stats", tagStats)

	// Filter posts to process
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.Content != ""
	})

	if lifecycle.IsServeIncremental(m) {
		if affected := lifecycle.GetServeAffectedPaths(m); len(affected) > 0 {
			filtered := posts[:0]
			for _, post := range posts {
				if affected[post.Path] {
					filtered = append(filtered, post)
				}
			}
			posts = filtered
		}
	}

	log.Printf("[hashtag_tags] Processing %d posts", len(posts))

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		content := p.processHashtagsInContent(post.Content, tagStats)
		post.Content = content
		return nil
	})
}

// TagStats holds statistics about a tag.
type TagStats struct {
	Tag             string // The tag name
	Slug            string // URL-safe slug
	Count           int    // Number of posts with this tag
	ReadingTime     int    // Total reading time (minutes)
	ReadingTimeText string // Formatted reading time text
}

// buildTagStats collects tag statistics from all published posts.
func (p *HashtagTagsPlugin) buildTagStats(m *lifecycle.Manager) map[string]*TagStats {
	config := m.Config()
	tagsConfig := getTagsConfig(config)

	tagStats := make(map[string]*TagStats)

	// Iterate through all posts
	for _, post := range m.Posts() {
		// Skip draft/unpublished/private/skip posts
		if post.Draft || !post.Published || post.Private || post.Skip {
			continue
		}

		// Get reading time for this post
		readingTime := 0
		if rt, ok := post.Extra["reading_time"].(int); ok {
			readingTime = rt
		}

		// Process each tag in the post
		for _, tag := range post.Tags {
			// Skip blacklisted and private tags
			if tagsConfig.IsBlacklisted(tag) || tagsConfig.IsPrivate(tag) {
				continue
			}

			if stat, exists := tagStats[tag]; exists {
				stat.Count++
				stat.ReadingTime += readingTime
			} else {
				slug := models.Slugify(tag)
				tagStats[tag] = &TagStats{
					Tag:         tag,
					Slug:        slug,
					Count:       1,
					ReadingTime: readingTime,
				}
			}
		}
	}

	// Format reading time text for each tag
	for _, stat := range tagStats {
		stat.ReadingTimeText = formatReadingTime(stat.ReadingTime)
	}

	return tagStats
}

// formatReadingTime converts minutes into human-readable text.
func formatReadingTime(minutes int) string {
	if minutes == 0 {
		return "less than a minute"
	}
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}

// hashtagRegex matches #tag patterns that are not inside code blocks.
// Pattern: matches #word or #word-with-dashes, surrounded by non-word characters or at string boundaries.
var hashtagRegex = regexp.MustCompile(`(?:^|[^\w#])#([a-zA-Z0-9](?:[a-zA-Z0-9_-]*[a-zA-Z0-9])?)`)

// hashtagCodeBlockRegex matches fenced code blocks (``` or ~~~).
var hashtagCodeBlockRegex = regexp.MustCompile("(?:```|~~~)[\\s\\S]*?(?:```|~~~)")

// processHashtagsInContent processes hashtag references in post content.
// Skips processing inside code blocks.
func (p *HashtagTagsPlugin) processHashtagsInContent(content string, tagStats map[string]*TagStats) string {
	// Split content by fenced code blocks to avoid transforming hashtags inside them
	codeBlocks := hashtagCodeBlockRegex.FindAllStringIndex(content, -1)

	if len(codeBlocks) == 0 {
		return p.processHashtagsInText(content, tagStats)
	}

	// Process content in segments, skipping code blocks
	var result strings.Builder
	lastEnd := 0

	for _, block := range codeBlocks {
		start, end := block[0], block[1]

		// Process text before this code block
		if start > lastEnd {
			processed := p.processHashtagsInText(content[lastEnd:start], tagStats)
			result.WriteString(processed)
		}

		// Keep code block unchanged
		result.WriteString(content[start:end])
		lastEnd = end
	}

	// Process any remaining text after the last code block
	if lastEnd < len(content) {
		processed := p.processHashtagsInText(content[lastEnd:], tagStats)
		result.WriteString(processed)
	}

	return result.String()
}

// processHashtagsInText processes hashtags in a text segment (not inside code blocks).
func (p *HashtagTagsPlugin) processHashtagsInText(text string, tagStats map[string]*TagStats) string {
	return hashtagRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the tag from the match
		groups := hashtagRegex.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}

		tag := groups[1]
		stat, found := tagStats[tag]
		if !found {
			return match
		}

		// Determine the prefix from the match (space, newline, etc.)
		hashPos := strings.Index(match, "#")
		prefix := ""
		if hashPos > 0 {
			prefix = match[:hashPos]
		}

		// Build data attributes
		dataAttrs := fmt.Sprintf(` data-tag=%q data-count=%d data-reading-time=%d data-reading-time-text=%q`,
			html.EscapeString(tag),
			stat.Count,
			stat.ReadingTime,
			html.EscapeString(stat.ReadingTimeText),
		)

		// Build the HTML link
		link := fmt.Sprintf(`<a href="/tags/%s/" class=%q%s>#%s</a>`,
			html.EscapeString(stat.Slug),
			p.cssClass,
			dataAttrs,
			html.EscapeString(tag),
		)

		return prefix + link
	})
}

// Ensure HashtagTagsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*HashtagTagsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*HashtagTagsPlugin)(nil)
	_ lifecycle.TransformPlugin = (*HashtagTagsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*HashtagTagsPlugin)(nil)
)
