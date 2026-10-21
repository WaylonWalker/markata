// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// KindContainer is the AST node kind for container elements.
var KindContainer = ast.NewNodeKind("Container")

// Container is an AST node representing a container block (::: class).
type Container struct {
	ast.BaseBlock
	// Classes holds the CSS classes for this container
	Classes []string
	// ContainerID holds the optional ID for this container
	ContainerID string
	// ExtraAttrs holds additional attributes
	ExtraAttrs map[string]string
	// Depth holds the number of leading colons used to open the container
	Depth int
}

// Kind returns the kind of this node.
func (n *Container) Kind() ast.NodeKind {
	return KindContainer
}

// Dump dumps the node for debugging.
func (n *Container) Dump(source []byte, level int) {
	attrs := map[string]string{
		"Classes": strings.Join(n.Classes, " "),
	}
	if n.ContainerID != "" {
		attrs["ID"] = n.ContainerID
	}
	ast.DumpHelper(n, source, level, attrs, nil)
}

// NewContainer creates a new Container node.
func NewContainer(classes []string, id string, attrs map[string]string, depth int) *Container {
	return &Container{
		Classes:     classes,
		ContainerID: id,
		ExtraAttrs:  attrs,
		Depth:       depth,
	}
}

// containerRegex matches container opening syntax:
// ::: class1 class2 {#id .extra-class key=value}
// Group 1: classes (space-separated)
// Group 2: optional attribute block
var containerRegex = regexp.MustCompile(`^(:{3,})\s*(\S.*?)?\s*(?:\{([^}]*)\})?\s*$`)

var containerAttrRegex = regexp.MustCompile(`(?:^|\s)(#[^\s]+|\.[^\s]+|[A-Za-z_:][-A-Za-z0-9_:.]*=(?:"[^"]*"|'[^']*'|[^\s]+))`)

// containerParser parses ::: container syntax.
type containerParser struct{}

// newContainerParser creates a new container parser.
func newContainerParser() parser.BlockParser {
	return &containerParser{}
}

// Trigger returns the trigger bytes for this parser.
func (p *containerParser) Trigger() []byte {
	return []byte{':'}
}

// Open is called when a new block starts.
func (p *containerParser) Open(_ ast.Node, reader text.Reader, _ parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()
	lineStr := string(line)
	trimmedLine := strings.TrimLeft(lineStr, " \t")

	// Check for opening :::
	if !strings.HasPrefix(trimmedLine, ":::") {
		return nil, parser.NoChildren
	}

	// Check for closing ::: (or more) with optional whitespace
	trimmed := strings.TrimSpace(trimmedLine)
	if isContainerClose(trimmed) {
		return nil, parser.NoChildren
	}

	// Parse the opening line
	matches := containerRegex.FindStringSubmatch(trimmedLine)
	if matches == nil {
		return nil, parser.NoChildren
	}

	marker := matches[1]
	depth := len(marker)

	// Extract classes
	var classes []string
	if matches[2] != "" {
		// Split by whitespace
		classes = append(classes, strings.Fields(matches[2])...)
	}

	// Extract attributes from {#id .class key=value} block
	var id string
	attrs := make(map[string]string)

	if matches[3] != "" {
		for _, attrMatch := range containerAttrRegex.FindAllStringSubmatch(matches[3], -1) {
			if len(attrMatch) != 2 {
				continue
			}
			part := strings.TrimSpace(attrMatch[1])
			if strings.HasPrefix(part, "#") {
				id = part[1:]
			} else if strings.HasPrefix(part, ".") {
				classes = append(classes, part[1:])
			} else if idx := strings.Index(part, "="); idx > 0 {
				key := part[:idx]
				value := strings.Trim(part[idx+1:], `"'`)
				attrs[key] = value
			}
		}
	}

	// Create the container node
	node := NewContainer(classes, id, attrs, depth)

	// Advance past the opening line
	reader.Advance(segment.Len() - 1)
	// Find the newline
	for {
		line, _ := reader.PeekLine()
		if len(line) == 0 {
			break
		}
		if line[0] == '\n' {
			reader.Advance(1)
			break
		}
		reader.Advance(1)
	}

	return node, parser.HasChildren
}

// Continue is called for each subsequent line.
func (p *containerParser) Continue(node ast.Node, reader text.Reader, _ parser.Context) parser.State {
	line, segment := reader.PeekLine()
	lineStr := strings.TrimSpace(string(line))
	trimmed := strings.TrimLeft(lineStr, " \t")

	if isContainerClose(trimmed) {
		depth := len(trimmed)
		if container, ok := node.(*Container); ok {
			if depth == container.Depth {
				reader.Advance(segment.Len())
				return parser.Close
			}
			if depth < container.Depth {
				return parser.Close
			}
		}
	}

	return parser.Continue | parser.HasChildren
}

// Close is called when the block ends.
func (p *containerParser) Close(_ ast.Node, _ text.Reader, _ parser.Context) {
	// Nothing special to do
}

func isContainerClose(line string) bool {
	if len(line) < 3 {
		return false
	}
	for _, r := range line {
		if r != ':' {
			return false
		}
	}
	return true
}

// CanInterruptParagraph returns true if this parser can interrupt a paragraph.
func (p *containerParser) CanInterruptParagraph() bool {
	return true
}

// CanAcceptIndentedLine returns true if this parser can accept an indented line.
func (p *containerParser) CanAcceptIndentedLine() bool {
	return true
}

// containerHTMLRenderer renders Container nodes to HTML.
type containerHTMLRenderer struct {
	html.Config
}

// newContainerHTMLRenderer creates a new container HTML renderer.
func newContainerHTMLRenderer() renderer.NodeRenderer {
	return &containerHTMLRenderer{
		Config: html.NewConfig(),
	}
}

// RegisterFuncs registers the render functions.
func (r *containerHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindContainer, r.renderContainer)
}

// renderContainer renders a Container node to HTML.
//
//nolint:errcheck // WriteString errors are handled at a higher level in goldmark
func (r *containerHTMLRenderer) renderContainer(w util.BufWriter, _ []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	container, ok := n.(*Container)
	if !ok {
		return ast.WalkContinue, nil
	}

	// Source/preview disclosures are native HTML. Content inside the
	// container stays ordinary fenced Markdown, with no JS to initialize.
	// The source panel is marked as UI so rich post copying can omit it.
	for _, class := range container.Classes {
		if class == "md-source" {
			if entering {
				label := strings.TrimSpace(container.ExtraAttrs["data-label"])
				if label == "" {
					label = "Show Markdown source"
				}
				_, _ = w.WriteString(`<details class="md-source" data-markata-ui="source-preview"><summary>`)
				_, _ = w.WriteString(escapeHTMLAttr(label))
				_, _ = w.WriteString("</summary>\n")
			} else {
				_, _ = w.WriteString("</details>\n")
			}
			return ast.WalkContinue, nil
		}
	}

	if entering {
		_, _ = w.WriteString("<div")

		// Add ID if present
		if container.ContainerID != "" {
			_, _ = w.WriteString(` id="`)
			_, _ = w.WriteString(escapeHTMLAttr(container.ContainerID))
			_, _ = w.WriteString(`"`)
		}

		// Add classes if present
		if len(container.Classes) > 0 {
			_, _ = w.WriteString(` class="`)
			for i, class := range container.Classes {
				if i > 0 {
					_, _ = w.WriteString(" ")
				}
				_, _ = w.WriteString(escapeHTMLAttr(class))
			}
			_, _ = w.WriteString(`"`)
		}

		// Add other attributes
		for key, value := range container.ExtraAttrs {
			_, _ = w.WriteString(` `)
			_, _ = w.WriteString(escapeHTMLAttr(key))
			_, _ = w.WriteString(`="`)
			_, _ = w.WriteString(escapeHTMLAttr(value))
			_, _ = w.WriteString(`"`)
		}

		_, _ = w.WriteString(">\n")
		r.renderGauge(w, container)
	} else {
		_, _ = w.WriteString("</div>\n")
	}

	return ast.WalkContinue, nil
}

// renderGauge turns an opt-in Markdown container into a native, accessible
// HTML meter/progress element. A visible Markdown caption remains ordinary
// container content, and missing or invalid numeric data leaves that caption
// intact without inserting a misleading bar. No client JavaScript is needed.
//
//nolint:errcheck // Writes to the Goldmark renderer are handled at the caller level.
func (r *containerHTMLRenderer) renderGauge(w util.BufWriter, container *Container) {
	tag := ""
	for _, class := range container.Classes {
		switch class {
		case "md-meter":
			tag = "meter"
		case "md-progress":
			tag = "progress"
		}
	}
	if tag == "" {
		return
	}

	valueText, hasValue := container.ExtraAttrs["data-value"]
	maxText, hasMax := container.ExtraAttrs["data-max"]
	label := strings.TrimSpace(container.ExtraAttrs["data-label"])
	if !hasValue || !hasMax || label == "" {
		return
	}
	value, valueErr := strconv.ParseFloat(valueText, 64)
	maximum, maxErr := strconv.ParseFloat(maxText, 64)
	if valueErr != nil || maxErr != nil ||
		math.IsNaN(value) || math.IsNaN(maximum) ||
		math.IsInf(value, 0) || math.IsInf(maximum, 0) ||
		maximum <= 0 || value < 0 || value > maximum {
		return
	}

	// Emit canonical numbers only; label text remains safely escaped.
	_, _ = w.WriteString("<" + tag + ` value="` + strconv.FormatFloat(value, 'f', -1, 64))
	_, _ = w.WriteString(`" max="` + strconv.FormatFloat(maximum, 'f', -1, 64))
	_, _ = w.WriteString(`" aria-label="` + escapeHTMLAttr(label) + `">`)
	_, _ = w.WriteString(strconv.FormatFloat(value, 'f', -1, 64) + " / " + strconv.FormatFloat(maximum, 'f', -1, 64))
	_, _ = w.WriteString("</" + tag + ">\n")
}

// escapeHTMLAttr escapes a string for use in an HTML attribute.
func escapeHTMLAttr(s string) string {
	var b strings.Builder
	for _, c := range s {
		switch c {
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '&':
			b.WriteString("&amp;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&#39;")
		default:
			b.WriteRune(c)
		}
	}
	return b.String()
}

// ContainerExtension is a goldmark extension for container syntax.
type ContainerExtension struct{}

// Extend adds the container parser and renderer to goldmark.
func (e *ContainerExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			// Lower priority than admonitions so !!! takes precedence
			util.Prioritized(newContainerParser(), 800),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(newContainerHTMLRenderer(), 500),
		),
	)
}
