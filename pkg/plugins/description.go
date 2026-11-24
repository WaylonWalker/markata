// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// DescriptionPlugin auto-generates descriptions for posts that don't have one.
// It extracts the first paragraph of content, strips markdown formatting,
// and truncates to a reasonable length for meta descriptions.
type DescriptionPlugin struct {
	// maxLength is the maximum length for generated descriptions (default: 160)
	maxLength int
}

// NewDescriptionPlugin creates a new DescriptionPlugin with default settings.
func NewDescriptionPlugin() *DescriptionPlugin {
	return &DescriptionPlugin{
		maxLength: 160,
	}
}

// Name returns the unique name of the plugin.
func (p *DescriptionPlugin) Name() string {
	return "description"
}

// Configure reads configuration options for the plugin.
func (p *DescriptionPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		if maxLength, ok := parseIntFromInterface(config.Extra["description_max_length"]); ok && maxLength > 0 {
			p.maxLength = maxLength
		}
	}
	return nil
}

// Priority returns the plugin priority for the given stage.
// Description should run early in transform to have content available for other plugins.
func (p *DescriptionPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityEarly
	}
	return lifecycle.PriorityDefault
}

// Transform generates descriptions for posts that don't have one,
// and strips wikilinks from all descriptions (including user-provided ones).
func (p *DescriptionPlugin) Transform(m *lifecycle.Manager) error {
	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip {
			return false
		}
		// Skip private posts - the encryption plugin will set a generic
		// description after encrypting. Generating a description from
		// private content would be wasteful and could leak if encryption
		// fails to run for any reason.
		if post.Private {
			return false
		}
		if post.Description != nil && *post.Description != "" {
			return true
		}
		return post.Content != ""
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		if post.Description != nil && *post.Description != "" {
			cleaned := p.stripMarkdown(p.stripWikilinks(*post.Description))
			post.Description = &cleaned
			return nil
		}

		description := p.generateDescription(post.Content)
		if description != "" {
			post.Description = &description
		}

		return nil
	})
}

// Regex patterns for stripping markdown
var (
	// Match markdown links [text](url) or [text][ref]
	markdownLinkRegex = regexp.MustCompile(`\[([^\]]+)\]\([^)]+\)|\[([^\]]+)\]\[[^\]]*\]`)

	// Match markdown images ![alt](url)
	markdownImageRegex = regexp.MustCompile(`!\[([^\]]*)\]\([^)]+\)`)

	// Match inline code `code`
	inlineCodeRegex = regexp.MustCompile("`[^`]+`")

	// Match bold/italic **text**, *text*, __text__, _text_
	emphasisRegex = regexp.MustCompile(`(\*\*|__)[^*_]+(\*\*|__)|(\*|_)[^*_]+(\*|_)`)

	// Match headers # Header
	headerRegex = regexp.MustCompile(`(?m)^#{1,6}\s+`)

	// Match HTML tags
	htmlTagRegex = regexp.MustCompile(`<[^>]+>`)

	// Match block quotes > text
	blockquoteRegex = regexp.MustCompile(`(?m)^>\s*`)

	// Match unordered list markers - item, * item, + item
	unorderedListRegex = regexp.MustCompile(`(?m)^\s*[-*+]\s+`)

	// Match ordered list markers 1. item
	orderedListRegex = regexp.MustCompile(`(?m)^\s*\d+\.\s+`)

	// Match horizontal rules ---, ***, ___
	hrRegex = regexp.MustCompile(`(?m)^[-*_]{3,}\s*$`)

	// Match code blocks ```code``` or ~~~code~~~
	codeBlockRegex = regexp.MustCompile("(?s)```.*?```|~~~.*?~~~")

	// Match reference-style link definitions [ref]: url
	linkDefRegex = regexp.MustCompile(`(?m)^\[[^\]]+\]:\s+.+$`)

	// Match wikilinks with optional spaces: [[slug]], [[ slug ]], [[slug|text]], [[ slug | text ]]
	// This is more flexible than wikilinkRegex in wikilinks.go to handle user-typed variations
	flexibleWikilinkRegex = regexp.MustCompile(`\[\[\s*([^\]|]+?)\s*(?:\|\s*([^\]]+?)\s*)?\]\]`)

	// Match multiple whitespace
	multiSpaceRegex = regexp.MustCompile(`\s+`)

	// Match inline attributes {.class #id key="value"}
	inlineAttributesRegex = regexp.MustCompile(`\{[.#][^}]+\}`)

	// Match Jinja/template expressions {{ ... }} and {% ... %} and {# ... #}
	// These may appear in content with jinja: true before the jinja_md plugin resolves them
	jinjaExprRegex = regexp.MustCompile(`\{\{.*?\}\}|\{%.*?%\}|\{#.*?#\}`)
)

// generateDescription creates a description from markdown content.
func (p *DescriptionPlugin) generateDescription(content string) string {
	// Get the first paragraph
	firstParagraph := p.extractFirstParagraph(content)
	if firstParagraph == "" {
		return ""
	}

	// Strip markdown formatting
	text := p.stripMarkdown(firstParagraph)
	if text == "" {
		return ""
	}

	// Truncate to max length
	text = p.truncate(text, p.maxLength)

	return text
}

// extractFirstParagraph extracts the first paragraph of content.
// Skips frontmatter, headers, code blocks, and empty lines.
func (p *DescriptionPlugin) extractFirstParagraph(content string) string {
	// Remove code blocks first to avoid parsing their content
	content = codeBlockRegex.ReplaceAllString(content, "")

	// Split into lines
	lines := strings.Split(content, "\n")

	var paragraph strings.Builder
	inParagraph := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			if inParagraph {
				// End of paragraph
				break
			}
			continue
		}

		// Skip headers
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Skip horizontal rules
		if hrRegex.MatchString(trimmed) {
			continue
		}

		// Skip link definitions
		if linkDefRegex.MatchString(trimmed) {
			continue
		}

		// Skip images on their own line
		if strings.HasPrefix(trimmed, "![") && strings.HasSuffix(trimmed, ")") {
			continue
		}

		// Found content - add to paragraph
		inParagraph = true
		if paragraph.Len() > 0 {
			paragraph.WriteString(" ")
		}
		paragraph.WriteString(trimmed)
	}

	return paragraph.String()
}

// stripMarkdown removes markdown formatting from text.
func (p *DescriptionPlugin) stripMarkdown(text string) string {
	// Remove Jinja/template expressions (before other processing)
	text = jinjaExprRegex.ReplaceAllString(text, "")
	// Remove images (before links to handle nested patterns)
	text = markdownImageRegex.ReplaceAllString(text, "")

	// Replace links with their text
	text = markdownLinkRegex.ReplaceAllString(text, "$1$2")

	// Replace wikilinks with their display text (or slug if no display text)
	// Uses flexibleWikilinkRegex to handle both [[slug]] and [[ slug ]] formats
	text = flexibleWikilinkRegex.ReplaceAllStringFunc(text, func(match string) string {
		groups := flexibleWikilinkRegex.FindStringSubmatch(match)
		if len(groups) >= 3 && groups[2] != "" {
			return strings.TrimSpace(groups[2]) // Use display text
		}
		if len(groups) >= 2 {
			return strings.TrimSpace(groups[1]) // Use slug
		}
		return ""
	})

	// Remove inline attributes
	text = inlineAttributesRegex.ReplaceAllString(text, "")

	// Remove inline code
	text = inlineCodeRegex.ReplaceAllString(text, "")

	// Remove emphasis markers but keep text
	text = emphasisRegex.ReplaceAllString(text, "$2$4")

	// Remove headers
	text = headerRegex.ReplaceAllString(text, "")

	// Remove HTML tags
	text = htmlTagRegex.ReplaceAllString(text, "")

	// Remove blockquote markers
	text = blockquoteRegex.ReplaceAllString(text, "")

	// Remove list markers
	text = unorderedListRegex.ReplaceAllString(text, "")
	text = orderedListRegex.ReplaceAllString(text, "")

	// Normalize whitespace
	text = multiSpaceRegex.ReplaceAllString(text, " ")
	text = strings.TrimSpace(text)

	return text
}

// stripWikilinks removes wikilink syntax from text, keeping the display text or slug.
// This handles both [[slug]] and [[ slug ]] formats (with optional spaces).
func (p *DescriptionPlugin) stripWikilinks(text string) string {
	return flexibleWikilinkRegex.ReplaceAllStringFunc(text, func(match string) string {
		groups := flexibleWikilinkRegex.FindStringSubmatch(match)
		if len(groups) >= 3 && groups[2] != "" {
			return strings.TrimSpace(groups[2]) // Use display text
		}
		if len(groups) >= 2 {
			return strings.TrimSpace(groups[1]) // Use slug
		}
		return ""
	})
}

// truncate shortens text to maxLength characters, breaking at word boundaries.
// Appends ellipsis if truncated.
func (p *DescriptionPlugin) truncate(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}

	// Find the last space before maxLength - 3 (for "...")
	truncateAt := maxLength - 3
	if truncateAt < 0 {
		truncateAt = 0
	}

	// Find word boundary
	lastSpace := -1
	for i, r := range text {
		if i > truncateAt {
			break
		}
		if unicode.IsSpace(r) {
			lastSpace = i
		}
	}

	if lastSpace > 0 {
		return strings.TrimSpace(text[:lastSpace]) + "..."
	}

	// No space found, hard truncate
	return text[:truncateAt] + "..."
}

// SetMaxLength sets the maximum length for generated descriptions.
func (p *DescriptionPlugin) SetMaxLength(length int) {
	if length > 0 {
		p.maxLength = length
	}
}

// Ensure DescriptionPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*DescriptionPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*DescriptionPlugin)(nil)
	_ lifecycle.TransformPlugin = (*DescriptionPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*DescriptionPlugin)(nil)
)
