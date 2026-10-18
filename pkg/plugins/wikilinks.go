// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"html"
	"regexp"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/yuin/goldmark/parser"
)

// WikilinksPlugin transforms [[slug]] and [[slug|text]] wikilink syntax
// into HTML anchor tags during the transform stage.
type WikilinksPlugin struct {
	// warnOnBroken controls whether to warn about broken links
	warnOnBroken bool
}

// NewWikilinksPlugin creates a new WikilinksPlugin.
func NewWikilinksPlugin() *WikilinksPlugin {
	return &WikilinksPlugin{
		warnOnBroken: true,
	}
}

// Name returns the unique name of the plugin.
func (p *WikilinksPlugin) Name() string {
	return "wikilinks"
}

// Configure reads configuration options for the plugin.
func (p *WikilinksPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		if warnOnBroken, ok := config.Extra["wikilinks_warn_broken"].(bool); ok {
			p.warnOnBroken = warnOnBroken
		}
	}
	return nil
}

// Transform processes wikilinks in all post content.
func (p *WikilinksPlugin) Transform(m *lifecycle.Manager) error {
	// Use the shared post index instead of building our own maps
	postIndex := m.PostIndex()
	postIndex.Refresh(m)

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

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		content, warnings, dependencies := p.processWikilinks(post.Content, postIndex)
		post.Content = content

		// Record dependencies for incremental build cache
		for _, dep := range dependencies {
			post.AddDependency(dep)
		}

		// Store warnings if any
		if len(warnings) > 0 && p.warnOnBroken {
			existingWarnings, _ := post.Extra["wikilink_warnings"].([]string) //nolint:errcheck // type assertion ok to fail
			post.Set("wikilink_warnings", append(existingWarnings, warnings...))
		}

		return nil
	})
}

// wikilinkRegex matches [[slug]] and [[slug|display text]] patterns.
// It captures:
// - Group 1: The slug
// - Group 2: Optional display text (after the pipe)
var wikilinkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\](\{[^}]+\})?`)

// wikilinksCodeBlockRegex matches fenced code blocks to avoid transforming wikilinks inside them.
var wikilinksCodeBlockRegex = regexp.MustCompile("(?s)(```[^`]*```|~~~[^~]*~~~)")

// processWikilinks replaces wikilink syntax with HTML anchor tags.
// Returns the processed content, any warnings about broken links, and a list of resolved slugs (dependencies).
// Wikilinks inside fenced code blocks are preserved and not transformed.
func (p *WikilinksPlugin) processWikilinks(content string, postIndex *lifecycle.PostIndex) (processed string, warnings, dependencies []string) {
	// Split content by fenced code blocks to avoid transforming wikilinks inside them
	// Match ``` or ~~~ fenced code blocks (with optional language identifier)

	// Find all code blocks and their positions
	codeBlocks := wikilinksCodeBlockRegex.FindAllStringIndex(content, -1)

	// If no code blocks, process the entire content
	if len(codeBlocks) == 0 {
		processed = p.processWikilinksInText(content, postIndex, &warnings, &dependencies)
		return processed, warnings, dependencies
	}

	// Process content in segments, skipping code blocks
	var result strings.Builder
	lastEnd := 0

	for _, block := range codeBlocks {
		start, end := block[0], block[1]

		// Process text before this code block
		if start > lastEnd {
			processed := p.processWikilinksInText(content[lastEnd:start], postIndex, &warnings, &dependencies)
			result.WriteString(processed)
		}

		// Keep code block unchanged
		result.WriteString(content[start:end])
		lastEnd = end
	}

	// Process any remaining text after the last code block
	if lastEnd < len(content) {
		processed := p.processWikilinksInText(content[lastEnd:], postIndex, &warnings, &dependencies)
		result.WriteString(processed)
	}

	return result.String(), warnings, dependencies
}

// processWikilinksInText processes wikilinks in a text segment (not inside code blocks).
// Uses the shared PostIndex for O(1) lookups.
// Records successfully resolved slugs in the dependencies slice.

func (p *WikilinksPlugin) processWikilinksInText(text string, postIndex *lifecycle.PostIndex, warnings, dependencies *[]string) string {
	return wikilinkRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract groups from the match
		groups := wikilinkRegex.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}

		slug := strings.TrimSpace(groups[1])
		displayText := ""
		if len(groups) >= 3 && groups[2] != "" {
			displayText = strings.TrimSpace(groups[2])
		}

		attrBlock := ""
		if len(groups) >= 4 && groups[3] != "" {
			attrBlock = groups[3]
		}

		return p.renderWikilink(match, slug, displayText, attrBlock, postIndex, warnings, dependencies)
	})
}

func (p *WikilinksPlugin) renderWikilink(match, slug, displayText, attrBlock string, postIndex *lifecycle.PostIndex, warnings, dependencies *[]string) string {
	// Extract fragment identifier if present (e.g., "my-post#section" -> slug="my-post", fragment="#section")
	var fragment string
	if idx := strings.Index(slug, "#"); idx >= 0 {
		fragment = slug[idx:] // includes the "#"
		slug = slug[:idx]
	}

	// Use the shared PostIndex for lookup
	targetPost := postIndex.LookupBySlug(slug)

	if targetPost == nil {
		// Target post not found - warn and keep original syntax
		*warnings = append(*warnings, fmt.Sprintf("broken wikilink: [[%s]]", slug))
		return match
	}

	// Record this as a dependency for incremental builds
	*dependencies = append(*dependencies, targetPost.Slug)

	// Determine the display text
	if displayText == "" {
		// Use post title if available, otherwise use slug
		if targetPost.Title != nil && *targetPost.Title != "" {
			displayText = *targetPost.Title
		} else {
			displayText = targetPost.Slug
		}
	}

	// Generate HTML anchor tag with wikilink class for styling
	href := targetPost.Href
	if href == "" {
		href = "/" + targetPost.Slug + "/"
	}
	// Append fragment identifier (e.g., "#section") after the trailing slash
	if fragment != "" {
		href += fragment
	}

	dataAttrs := p.renderWikilinkDataAttrs(targetPost)

	classAttr, extra := buildWikilinkAttributes(attrBlock)

	return fmt.Sprintf(`<a href=%q class=%q%s%s>%s</a>`,
		html.EscapeString(href),
		classAttr,
		dataAttrs,
		extra,
		html.EscapeString(displayText))
}

func (p *WikilinksPlugin) renderWikilinkDataAttrs(targetPost *models.Post) string {
	// Build data attributes for tooltip
	// Skip metadata for private posts to prevent content leaks
	dataAttrs := ""
	if !targetPost.Private {
		if targetPost.Title != nil && *targetPost.Title != "" {
			dataAttrs += fmt.Sprintf(` data-title=%q`, html.EscapeString(*targetPost.Title))
		}
		if targetPost.Description != nil && *targetPost.Description != "" {
			dataAttrs += fmt.Sprintf(` data-description=%q`, html.EscapeString(*targetPost.Description))
		}
		if targetPost.Date != nil {
			dataAttrs += fmt.Sprintf(` data-date=%q`, targetPost.Date.Format("2006-01-02"))
		}
	}

	return dataAttrs
}

func buildWikilinkAttributes(attrBlock string) (classAttr, extra string) {
	classAttr = "wikilink"
	idAttr := ""
	extraAttrs := map[string]string{}
	if attrBlock != "" {
		if attrs, ok := parseAttributeOnly(attrBlock); ok {
			classes, id, extra := splitAttributes(attrs)
			if len(classes) > 0 {
				classAttr = classAttr + " " + strings.Join(classes, " ")
			}
			if id != "" {
				idAttr = id
			}
			for key, val := range extra {
				extraAttrs[key] = val
			}
		}
	}

	attrParts := make([]string, 0, len(extraAttrs))
	for key := range extraAttrs {
		attrParts = append(attrParts, key)
	}
	sort.Strings(attrParts)

	var extraBuilder strings.Builder
	for _, key := range attrParts {
		extraBuilder.WriteString(fmt.Sprintf(` %s=%q`, html.EscapeString(key), html.EscapeString(extraAttrs[key])))
	}
	if idAttr != "" {
		extraBuilder.WriteString(fmt.Sprintf(` id=%q`, html.EscapeString(idAttr)))
	}

	return classAttr, extraBuilder.String()
}

func splitAttributes(attrs parser.Attributes) (classes []string, id string, extra map[string]string) {
	classes = []string{}
	id = ""
	extra = map[string]string{}
	for _, attr := range attrs {
		name := string(attr.Name)
		value := normalizeAttributeValue(attr.Value)
		if name == "class" {
			classes = append(classes, strings.Fields(value)...)
			continue
		}
		if name == "id" {
			id = value
			continue
		}
		if name == "" {
			continue
		}
		extra[name] = value
	}
	return classes, id, extra
}

// SetWarnOnBroken enables or disables warnings for broken links.
func (p *WikilinksPlugin) SetWarnOnBroken(warn bool) {
	p.warnOnBroken = warn
}

// Ensure WikilinksPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*WikilinksPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*WikilinksPlugin)(nil)
	_ lifecycle.TransformPlugin = (*WikilinksPlugin)(nil)
)
