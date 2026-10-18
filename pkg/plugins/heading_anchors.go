// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// HeadingAnchorsPlugin adds anchor links to headings in rendered HTML.
// It processes article_html during the render stage, after render_markdown.
type HeadingAnchorsPlugin struct {
	// enabled controls whether the plugin processes headings
	enabled bool

	// minLevel is the minimum heading level to process (1-6, default: 2)
	minLevel int

	// maxLevel is the maximum heading level to process (1-6, default: 4)
	maxLevel int

	// position is where to insert the anchor link ("start" or "end", default: "end")
	position string

	// symbol is the link text for the anchor (default: "#")
	symbol string

	// class is the CSS class for the anchor link (default: "heading-anchor")
	class string
	wear  *HeadingWearPlugin
}

// HeadingWearPlugin applies heading texture glyph wrappers independently of
// heading anchor configuration.
type HeadingWearPlugin struct {
	enabled bool
}

// NewHeadingWearPlugin creates a disabled heading-wear pass.
func NewHeadingWearPlugin() *HeadingWearPlugin { return &HeadingWearPlugin{} }

func (p *HeadingWearPlugin) Name() string { return "heading_wear" }

func (p *HeadingWearPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

func (p *HeadingWearPlugin) Configure(m *lifecycle.Manager) error {
	p.enabled = false
	if configured, ok := m.Config().Extra["models_config"].(*models.Config); ok {
		p.enabled = configured.Theme.HeadingTexture.Kind != renderingTextureNone && configured.Theme.HeadingTexture.ColorMix > 0
	}
	return nil
}

func (p *HeadingWearPlugin) SetEnabled(enabled bool) { p.enabled = enabled }

func (p *HeadingWearPlugin) Render(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.ArticleHTML != "" && strings.Contains(post.ArticleHTML, "<h")
	})
	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

func (p *HeadingWearPlugin) processPost(post *models.Post) error {
	post.ArticleHTML = headingTagRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		submatches := headingTagRegex.FindStringSubmatch(match)
		if len(submatches) < 4 {
			return match
		}
		return fmt.Sprintf("<h%s%s>%s</h%s>", submatches[1], submatches[2], wrapHeadingGlyphText(submatches[3]), submatches[1])
	})
	return nil
}

// NewHeadingAnchorsPlugin creates a new HeadingAnchorsPlugin with default settings.
func NewHeadingAnchorsPlugin() *HeadingAnchorsPlugin {
	return &HeadingAnchorsPlugin{
		enabled:  true,
		minLevel: 2,
		maxLevel: 4,
		position: PositionEnd,
		symbol:   "#",
		class:    "heading-anchor",
		wear:     NewHeadingWearPlugin(),
	}
}

// Name returns the unique name of the plugin.
func (p *HeadingAnchorsPlugin) Name() string {
	return "heading_anchors"
}

// Priority returns the plugin's priority for the render stage.
// Returns a positive priority to run after render_markdown.
func (p *HeadingAnchorsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after render_markdown
	}
	return lifecycle.PriorityDefault
}

// Configure reads configuration options for the plugin from config.Extra.
// Configuration is expected in config.Extra["heading_anchors"] as a map.
func (p *HeadingAnchorsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if p.wear == nil {
		p.wear = NewHeadingWearPlugin()
	}
	if err := p.wear.Configure(m); err != nil {
		return err
	}
	if config.Extra == nil {
		return nil
	}

	// Look for heading_anchors configuration
	haConfig, ok := config.Extra["heading_anchors"]
	if !ok {
		return nil
	}

	// Handle map[string]interface{} configuration
	if cfgMap, ok := haConfig.(map[string]interface{}); ok {
		if enabled, ok := cfgMap["enabled"].(bool); ok {
			p.enabled = enabled
		}
		if minLevel, ok := parseIntFromInterface(cfgMap["min_level"]); ok && minLevel >= 1 && minLevel <= 6 {
			p.minLevel = minLevel
		}
		if maxLevel, ok := parseIntFromInterface(cfgMap["max_level"]); ok && maxLevel >= 1 && maxLevel <= 6 {
			p.maxLevel = maxLevel
		}
		if position, ok := cfgMap["position"].(string); ok && (position == PositionStart || position == PositionEnd) {
			p.position = position
		}
		if symbol, ok := cfgMap["symbol"].(string); ok {
			p.symbol = symbol
		}
		if class, ok := cfgMap["class"].(string); ok {
			p.class = class
		}
	}

	return nil
}

// headingTagRegex matches opening heading tags with optional attributes.
// Captures: Group 1 = level (1-6), Group 2 = attributes (including id), Group 3 = content, Group 4 = closing tag
var headingTagRegex = regexp.MustCompile(`(?i)<h([1-6])([^>]*)>(.*?)</h[1-6]>`)

// idAttrRegex extracts the id attribute value from tag attributes.
var idAttrRegex = regexp.MustCompile(`(?i)\bid=["']([^"']+)["']`)

// Render processes article_html and adds anchor links to headings.
// Posts with Skip=true or empty ArticleHTML are skipped.
func (p *HeadingAnchorsPlugin) Render(m *lifecycle.Manager) error {
	// Run wear before inserting generated anchors, and independently of the
	// anchor enabled and level-range settings.
	if p.wear != nil {
		if err := p.wear.Render(m); err != nil {
			return err
		}
	}
	if !p.enabled {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return strings.Contains(post.ArticleHTML, "<h")
	})

	return m.ProcessPostsSliceConcurrently(posts, p.processPost)
}

// processPost adds anchor links to headings in a single post's ArticleHTML.
func (p *HeadingAnchorsPlugin) processPost(post *models.Post) error {
	if post.Skip || post.ArticleHTML == "" {
		return nil
	}

	// Track IDs to handle duplicates
	idCounts := make(map[string]int)

	post.ArticleHTML = headingTagRegex.ReplaceAllStringFunc(post.ArticleHTML, func(match string) string {
		return p.processHeading(match, idCounts)
	})

	return nil
}

// processHeading processes a single heading match and adds an anchor link.
func (p *HeadingAnchorsPlugin) processHeading(match string, idCounts map[string]int) string {
	submatches := headingTagRegex.FindStringSubmatch(match)
	if len(submatches) < 4 {
		return match
	}

	levelStr := submatches[1]
	attrs := submatches[2]
	content := submatches[3]

	// Skip headings that already have an anchor link (e.g., from goldmark anchor extension)
	if strings.Contains(content, `class="anchor"`) {
		return match
	}

	// Skip feed card titles embedded in article content via render_feed();
	// they are already links and must not gain ids or anchor glyphs.
	if strings.Contains(attrs, "card-title") {
		return match
	}

	// Parse level
	level := int(levelStr[0] - '0')
	if level < p.minLevel || level > p.maxLevel {
		return match
	}

	// Extract or generate ID
	id := p.extractID(attrs)
	if id == "" {
		id = p.generateID(content, idCounts)
		// Add id attribute if not present
		if attrs == "" {
			attrs = fmt.Sprintf(` id=%q`, id)
		} else {
			attrs = fmt.Sprintf(` id=%q%s`, id, attrs)
		}
	} else {
		// Track existing ID for duplicates
		idCounts[id]++
	}

	// Create anchor link
	anchor := p.createAnchor(id)

	// Insert anchor at configured position
	var newContent string
	if p.position == PositionStart {
		newContent = anchor + content
	} else {
		newContent = content + anchor
	}

	return fmt.Sprintf("<h%s%s>%s</h%s>", levelStr, attrs, newContent, levelStr)
}

func wrapHeadingGlyphText(content string) string {
	// Parse the fragment as HTML. Wrapping serialized strings is unsafe: an
	// inline element can contain another element, quoted attributes can contain
	// angle brackets, and the heading anchor must never be worn.
	root, err := html.ParseFragment(strings.NewReader(content), &html.Node{Type: html.ElementNode, Data: "span", DataAtom: atom.Span})
	if err != nil {
		return content
	}
	var wearText func(*html.Node, bool)
	wearText = func(node *html.Node, inAnchor bool) {
		anchor := inAnchor || (node.Type == html.ElementNode && node.Data == "a" && isGeneratedHeadingAnchor(node))
		if node.Type == html.TextNode && !anchor && strings.TrimSpace(node.Data) != "" {
			parent := node.Parent
			if parent != nil && parent.Type == html.ElementNode && parent.Data == "span" && hasClass(parent, "heading-wear-glyph") {
				return
			}
			wrapped := &html.Node{Type: html.ElementNode, Data: "span", Attr: []html.Attribute{{Key: "class", Val: "heading-wear-glyph"}}}
			wrapped.AppendChild(&html.Node{Type: html.TextNode, Data: node.Data})
			parent.InsertBefore(wrapped, node)
			parent.RemoveChild(node)
			return
		}
		for child := node.FirstChild; child != nil; {
			next := child.NextSibling
			wearText(child, anchor)
			child = next
		}
	}
	container := &html.Node{Type: html.ElementNode, Data: "div"}
	for _, node := range root {
		container.AppendChild(node)
	}
	wearText(container, false)
	var out bytes.Buffer
	for child := container.FirstChild; child != nil; child = child.NextSibling {
		if err := html.Render(&out, child); err != nil {
			return ""
		}
	}
	return out.String()
}

func hasClass(node *html.Node, class string) bool {
	for _, attr := range node.Attr {
		if attr.Key == htmlAttrClass {
			for _, value := range strings.Fields(attr.Val) {
				if value == class {
					return true
				}
			}
		}
	}
	return false
}

func isGeneratedHeadingAnchor(node *html.Node) bool {
	return hasClass(node, "anchor") || hasClass(node, "heading-anchor")
}

// extractID extracts the id attribute value from heading attributes.
func (p *HeadingAnchorsPlugin) extractID(attrs string) string {
	matches := idAttrRegex.FindStringSubmatch(attrs)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// generateID creates a URL-safe ID from heading content.
// Handles duplicate IDs by appending numbers.
func (p *HeadingAnchorsPlugin) generateID(content string, idCounts map[string]int) string {
	// Strip HTML tags from content
	text := stripHTMLTags(content)

	// Use the shared Slugify function for consistent slug generation
	id := models.Slugify(text)

	// Handle empty ID
	if id == "" {
		id = "heading"
	}

	// Handle duplicates
	baseID := id
	count := idCounts[baseID]
	idCounts[baseID] = count + 1

	if count > 0 {
		id = fmt.Sprintf("%s-%d", baseID, count)
	}

	return id
}

// stripHTMLRegex matches HTML tags for removal.
var stripHTMLRegex = regexp.MustCompile(`<[^>]*>`)

// stripHTMLTags removes HTML tags from content.
func stripHTMLTags(content string) string {
	return stripHTMLRegex.ReplaceAllString(content, "")
}

// createAnchor creates an anchor link element for the given ID.
func (p *HeadingAnchorsPlugin) createAnchor(id string) string {
	id = stdhtml.EscapeString(id)
	class := stdhtml.EscapeString(p.class)
	symbol := stdhtml.EscapeString(p.symbol)
	// Add leading space for "end" position, trailing space for "start"
	if p.position == PositionStart {
		return fmt.Sprintf(`<a href="#%s" class=%q>%s</a> `, id, class, symbol)
	}
	return fmt.Sprintf(` <a href="#%s" class=%q>%s</a>`, id, class, symbol)
}

// SetEnabled enables or disables the plugin.
func (p *HeadingAnchorsPlugin) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// SetLevelRange sets the minimum and maximum heading levels to process.
func (p *HeadingAnchorsPlugin) SetLevelRange(minLevel, maxLevel int) {
	if minLevel >= 1 && minLevel <= 6 {
		p.minLevel = minLevel
	}
	if maxLevel >= 1 && maxLevel <= 6 && maxLevel >= minLevel {
		p.maxLevel = maxLevel
	}
}

// SetPosition sets the anchor position ("start" or "end").
func (p *HeadingAnchorsPlugin) SetPosition(position string) {
	if position == PositionStart || position == PositionEnd {
		p.position = position
	}
}

// SetSymbol sets the anchor link symbol.
func (p *HeadingAnchorsPlugin) SetSymbol(symbol string) {
	p.symbol = symbol
}

// SetClass sets the anchor link CSS class.
func (p *HeadingAnchorsPlugin) SetClass(class string) {
	p.class = class
}

// Ensure HeadingAnchorsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*HeadingAnchorsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*HeadingAnchorsPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*HeadingAnchorsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*HeadingAnchorsPlugin)(nil)
	_ lifecycle.Plugin          = (*HeadingWearPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*HeadingWearPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*HeadingWearPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*HeadingWearPlugin)(nil)
)
