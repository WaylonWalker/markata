// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// TocPlugin extracts headings from markdown content and builds a
// hierarchical table of contents during the transform stage.
type TocPlugin struct {
	// minLevel is the minimum heading level to include (default: 2)
	minLevel int

	// maxLevel is the maximum heading level to include (default: 4)
	maxLevel int

	// globalTocConfig is the global TOC configuration
	globalTocConfig *models.TocConfig
}

// NewTocPlugin creates a new TocPlugin with default settings.
func NewTocPlugin() *TocPlugin {
	return &TocPlugin{
		minLevel: 2,
		maxLevel: 4,
	}
}

// Name returns the unique name of the plugin.
func (p *TocPlugin) Name() string {
	return "toc"
}

// Configure reads configuration options for the plugin.
func (p *TocPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		if minLevel, ok := parseIntFromInterface(config.Extra["toc_min_level"]); ok && minLevel >= 1 && minLevel <= 6 {
			p.minLevel = minLevel
		}
		if maxLevel, ok := parseIntFromInterface(config.Extra["toc_max_level"]); ok && maxLevel >= 1 && maxLevel <= 6 {
			p.maxLevel = maxLevel
		}

		// Get global TOC config from Extra
		if toc, ok := config.Extra["toc"].(models.TocConfig); ok {
			if toc.IsEnabled() || toc.IsAutoEnable() {
				p.globalTocConfig = &toc
			}
		}
	}

	return nil
}

// tocPlaceholderRegex matches [[toc]] markdown placeholder
var tocPlaceholderRegex = regexp.MustCompile(`(?i)\[\[toc\]\]`)

// Transform extracts headings from markdown and builds a TOC for each post.
func (p *TocPlugin) Transform(m *lifecycle.Manager) error {
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.Content != ""
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		// Check for [[toc]] placeholder
		hasPlaceholder := tocPlaceholderRegex.MatchString(post.Content)
		post.TocPlaceholder = hasPlaceholder

		// Extract TOC from content
		toc := p.extractTOC(post.Content)
		tocLinkCount := countTOCLinks(toc)

		// Count words in content
		wordCount := countWords(post.Content)

		// Determine if TOC should be enabled using priority logic
		shouldShow := post.IsTocEnabled(p.globalTocConfig, tocLinkCount, wordCount)

		if shouldShow && len(toc) > 0 {
			post.Set("toc", toc)
		}

		// If placeholder exists, replace it with a marker for later rendering
		if hasPlaceholder {
			// Replace [[toc]] with a marker that templates can detect
			post.Set("toc_placeholder_present", true)
		}

		return nil
	})
}

// countTOCLinks counts the total number of TOC entries (including nested)
func countTOCLinks(entries []*TocEntry) int {
	count := 0
	for _, entry := range entries {
		count += 1 + countTOCLinks(entry.Children)
	}
	return count
}

// countWords counts the number of words in text
func countWords(s string) int {
	words := 0
	inWord := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			inWord = false
		} else if !inWord {
			words++
			inWord = true
		}
	}
	return words
}

// TocEntry represents a single entry in the table of contents.
type TocEntry struct {
	// Level is the heading level (1-6)
	Level int `json:"level"`

	// Text is the heading text
	Text string `json:"text"`

	// ID is the anchor ID for the heading
	ID string `json:"id"`

	// Children contains nested headings
	Children []*TocEntry `json:"children,omitempty"`
}

// headingRegex matches ATX-style markdown headings (# Heading).
// Captures: Group 1 = hash marks, Group 2 = heading text
var headingRegex = regexp.MustCompile(`(?m)^(#{1,6})\s+(.+?)(?:\s*#*)?\s*$`)

// extractTOC extracts headings from markdown content and builds a hierarchical TOC.
func (p *TocPlugin) extractTOC(content string) []*TocEntry {
	matches := headingRegex.FindAllStringSubmatch(stripFencedCodeBlocks(content), -1)
	if len(matches) == 0 {
		return nil
	}

	// Extract flat list of headings
	headings := make([]*TocEntry, 0, len(matches))
	idCounts := make(map[string]int)

	for _, match := range matches {
		level := len(match[1])
		text := cleanHeadingText(match[2])

		// Explicit {#id} attributes win, matching the rendered heading.
		// They do not consume a slot in the collision counter, mirroring
		// HeadingIDTransformer which leaves explicit IDs untouched. Every
		// heading level consumes collision slots, as it does when rendering,
		// so this runs before the level filter.
		var id string
		if explicit := tocExplicitIDRegex.FindStringSubmatch(match[2]); explicit != nil {
			id = explicitHeadingID(explicit[1])
		} else {
			id = p.generateID(text, idCounts)
		}

		// Skip headings outside our level range
		if level < p.minLevel || level > p.maxLevel {
			continue
		}

		headings = append(headings, &TocEntry{
			Level:    level,
			Text:     text,
			ID:       id,
			Children: make([]*TocEntry, 0),
		})
	}

	if len(headings) == 0 {
		return nil
	}

	// Build hierarchical structure
	return p.buildHierarchy(headings)
}

var (
	// tocHTMLTagRegex matches HTML tags such as resolved wikilinks (<a class="wikilink">).
	tocHTMLTagRegex = regexp.MustCompile(`<[^>]+>`)
	// tocMarkdownLinkRegex matches [text](url) and ![alt](src), keeping the text.
	tocMarkdownLinkRegex = regexp.MustCompile(`!?\[([^\]]*)\]\([^)]*\)`)
	// tocWikilinkRegex matches unresolved [[target]] or [[target|alias]].
	tocWikilinkRegex = regexp.MustCompile(`\[\[\s*([^\]|]+?)\s*(?:\|\s*([^\]]+?)\s*)?\]\]`)
	// tocHeadingAttrRegex matches a trailing {#id .class} attribute block.
	tocHeadingAttrRegex = regexp.MustCompile(`\s*\{[^}]*\}\s*$`)
	// tocExplicitIDRegex captures the id from a trailing {#id .class} block.
	// Capture to whitespace, then validate ASCII punctuation separately. This
	// avoids an ambiguous punctuation range while preserving _-:.
	tocExplicitIDRegex = regexp.MustCompile(`\{[^}]*#([^\s}]+)[^}]*\}\s*$`)
	// tocInlineMarkupRegex matches emphasis, strike and mark delimiters.
	// Underscores are handled separately because intraword `_` is literal.
	tocInlineMarkupRegex = regexp.MustCompile("[*`~=]+")
	// tocWhitespaceRegex collapses runs of whitespace.
	tocWhitespaceRegex = regexp.MustCompile(`\s+`)
)

func explicitHeadingID(candidate string) string {
	for i, r := range candidate {
		if r < '!' || r > '~' || r == '_' || r == '-' || r == ':' || r == '.' {
			continue
		}
		if (r >= '!' && r <= '/') || (r >= ':' && r <= '@') || (r >= '[' && r <= '`') || (r >= '{' && r <= '~') {
			return candidate[:i]
		}
	}
	return candidate
}

// cleanHeadingText reduces a markdown heading line to its visible text so the
// TOC never shows raw HTML or markdown syntax. Earlier Transform plugins
// (wikilinks, jinja_md) may already have rewritten inline markup to HTML.
func cleanHeadingText(raw string) string {
	text := tocHeadingAttrRegex.ReplaceAllString(raw, "")

	// Code span contents are rendered literally (e.g. `update_meta`), so keep
	// them out of the markup-stripping passes below.
	text, codeSpans := protectCodeSpans(text)

	text = tocHTMLTagRegex.ReplaceAllString(text, "")
	text = tocMarkdownLinkRegex.ReplaceAllString(text, "$1")
	text = tocWikilinkRegex.ReplaceAllStringFunc(text, func(m string) string {
		parts := tocWikilinkRegex.FindStringSubmatch(m)
		if parts[2] != "" {
			return parts[2]
		}
		return parts[1]
	})
	text = tocInlineMarkupRegex.ReplaceAllString(text, "")
	text = stripEmphasisUnderscores(text)
	for i, code := range codeSpans {
		text = strings.Replace(text, codeSpanPlaceholder(i), code, 1)
	}
	text = html.UnescapeString(text)
	text = tocWhitespaceRegex.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

// protectCodeSpans swaps inline code spans for placeholders and returns
// their literal contents in order.
func protectCodeSpans(s string) (protected string, spans []string) {
	if !strings.Contains(s, "`") {
		return s, nil
	}
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] != '`' {
			b.WriteByte(s[i])
			i++
			continue
		}
		n := 0
		for i+n < len(s) && s[i+n] == '`' {
			n++
		}
		fence := s[i : i+n]
		closeAt := -1
		for j := i + n; j < len(s); {
			k := strings.Index(s[j:], fence)
			if k < 0 {
				break
			}
			k += j
			end := k + n
			if end < len(s) && s[end] == '`' {
				for end < len(s) && s[end] == '`' {
					end++
				}
				j = end
				continue
			}
			closeAt = k
			break
		}
		if closeAt < 0 {
			b.WriteString(fence)
			i += n
			continue
		}
		inner := s[i+n : closeAt]
		if len(inner) > 1 && strings.HasPrefix(inner, " ") && strings.HasSuffix(inner, " ") && strings.TrimSpace(inner) != "" {
			inner = inner[1 : len(inner)-1]
		}
		b.WriteString(codeSpanPlaceholder(len(spans)))
		spans = append(spans, inner)
		i = closeAt + n
	}
	return b.String(), spans
}

func codeSpanPlaceholder(i int) string {
	return "\x00" + strconv.Itoa(i) + "\x00"
}

// stripEmphasisUnderscores removes paired `_` emphasis delimiters while
// keeping intraword underscores (snake_case) and unpaired ones (_headers),
// approximating CommonMark's rules for heading text.
//
//nolint:gocyclo // The scanner coordinates delimiter-run pairing and context rules in one pass.
func stripEmphasisUnderscores(s string) string {
	if !strings.Contains(s, "_") {
		return s
	}
	runes := []rune(s)
	isWord := func(i int) bool {
		return i >= 0 && i < len(runes) && runes[i] != '_' && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i]))
	}
	isSpace := func(i int) bool {
		return i < 0 || i >= len(runes) || unicode.IsSpace(runes[i])
	}

	type delimRun struct{ start, end int }
	var runs []delimRun
	for i := 0; i < len(runes); i++ {
		if runes[i] != '_' {
			continue
		}
		start := i
		for i+1 < len(runes) && runes[i+1] == '_' {
			i++
		}
		if !(isWord(start-1) && isWord(i+1)) {
			runs = append(runs, delimRun{start, i + 1})
		}
	}

	drop := make([]bool, len(runes))
	for i := 0; i < len(runs); i++ {
		opener := runs[i]
		if isSpace(opener.end) {
			continue
		}
		for j := i + 1; j < len(runs); j++ {
			closer := runs[j]
			if isSpace(closer.start - 1) {
				continue
			}
			for k := opener.start; k < opener.end; k++ {
				drop[k] = true
			}
			for k := closer.start; k < closer.end; k++ {
				drop[k] = true
			}
			i = j
			break
		}
	}

	var b strings.Builder
	b.Grow(len(s))
	for i, r := range runes {
		if !drop[i] {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// stripFencedCodeBlocks blanks out fenced code blocks so `#` lines inside
// examples are not mistaken for headings.
func stripFencedCodeBlocks(content string) string {
	if !strings.Contains(content, "```") && !strings.Contains(content, "~~~") {
		return content
	}
	lines := strings.Split(content, "\n")
	var fence string
	for i, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		if len(line)-len(trimmed) > 3 {
			if fence != "" {
				lines[i] = ""
			}
			continue
		}
		if fence == "" {
			if marker := fenceMarker(trimmed); marker != "" {
				if marker[0] == '`' && strings.Contains(trimmed[len(marker):], "`") {
					continue
				}
				fence = marker
				lines[i] = ""
			}
			continue
		}
		lines[i] = ""
		if marker := fenceMarker(trimmed); marker != "" && marker[0] == fence[0] &&
			len(marker) >= len(fence) && strings.TrimSpace(trimmed[len(marker):]) == "" {
			fence = ""
		}
	}
	return strings.Join(lines, "\n")
}

func fenceMarker(line string) string {
	if line == "" || (line[0] != '`' && line[0] != '~') {
		return ""
	}
	n := 0
	for n < len(line) && line[n] == line[0] {
		n++
	}
	if n < 3 {
		return ""
	}
	return line[:n]
}

// generateID creates a URL-safe ID from heading text.
// Handles duplicate IDs by appending numbers.
func (p *TocPlugin) generateID(text string, idCounts map[string]int) string {
	return generateHeadingID(text, idCounts)
}

// buildHierarchy converts a flat list of headings into a nested structure.
func (p *TocPlugin) buildHierarchy(headings []*TocEntry) []*TocEntry {
	if len(headings) == 0 {
		return nil
	}

	// Find the minimum level to use as root level
	minLevel := 6
	for _, h := range headings {
		if h.Level < minLevel {
			minLevel = h.Level
		}
	}

	roots := make([]*TocEntry, 0, len(headings))
	stack := make([]*TocEntry, 0, len(headings))

	for _, heading := range headings {
		// Adjust level relative to minimum
		entry := &TocEntry{
			Level:    heading.Level,
			Text:     heading.Text,
			ID:       heading.ID,
			Children: make([]*TocEntry, 0),
		}

		// Pop stack until we find a parent at a lower level
		for len(stack) > 0 && stack[len(stack)-1].Level >= entry.Level {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			// This is a root-level heading
			roots = append(roots, entry)
		} else {
			// Add as child of the top of stack
			parent := stack[len(stack)-1]
			parent.Children = append(parent.Children, entry)
		}

		// Push current heading onto stack
		stack = append(stack, entry)
	}

	return roots
}

// SetLevelRange sets the minimum and maximum heading levels to include.
func (p *TocPlugin) SetLevelRange(minLevel, maxLevel int) {
	if minLevel >= 1 && minLevel <= 6 {
		p.minLevel = minLevel
	}
	if maxLevel >= 1 && maxLevel <= 6 && maxLevel >= minLevel {
		p.maxLevel = maxLevel
	}
}

// Ensure TocPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*TocPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*TocPlugin)(nil)
	_ lifecycle.TransformPlugin = (*TocPlugin)(nil)
)
