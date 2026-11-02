package plugins

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

// Code fence info-string extras.
//
// Authors can annotate fenced code blocks with a title and highlighted lines
// using either the goldmark attribute syntax or the shorthand common in other
// static site generators:
//
//	```python title="app.py" {2,4-6}
//	```python {title="app.py" hl_lines=[2,"4-6"]}
//	```python {2,4-6}
//
// codeFenceTransformer normalises the shorthand into node attributes that the
// highlighting extension understands (hl_lines, title). codeFenceWrapper then
// emits a wrapper with data-lang/data-title so the theme can render a header
// with a language badge, filename and copy button.

var (
	// key="value" or key='value' or key=value (no spaces) pairs in an info string.
	codeFenceKVPattern = regexp.MustCompile(`([A-Za-z_][\w-]*)=("([^"]*)"|'([^']*)'|([^\s{}]+))`)
	// Bare {1,3-5} highlight shorthand.
	codeFenceRangePattern = regexp.MustCompile(`^\{\s*(\d+(?:\s*-\s*\d+)?(?:\s*,\s*\d+(?:\s*-\s*\d+)?)*)\s*\}$`)
)

// pluginFenceLanguages are fence languages consumed by other plugins
// (csv_fence, chartjs, mermaid) that match on goldmark's plain
// <pre><code class="language-x"> markup. Chroma has lexers for some of them
// (csv), so they are routed around the highlighter entirely.
var pluginFenceLanguages = map[string]bool{
	"csv":     true,
	"chartjs": true,
	"mermaid": true,
}

// KindPluginFence is the AST node kind for fences handed to other plugins.
var KindPluginFence = ast.NewNodeKind("PluginFence")

// pluginFence is a fenced code block rendered as plain <pre><code> so that a
// language plugin can transform it after markdown rendering.
type pluginFence struct {
	ast.BaseBlock
	lang  []byte
	lines *text.Segments
}

// Kind implements ast.Node.
func (n *pluginFence) Kind() ast.NodeKind { return KindPluginFence }

// Dump implements ast.Node.
func (n *pluginFence) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, map[string]string{"Lang": string(n.lang)}, nil)
}

// IsRaw implements ast.Node; the block body is not parsed as inline content.
func (n *pluginFence) IsRaw() bool { return true }

// pluginFenceRenderer writes pluginFence nodes as <pre><code class="language-x">.
type pluginFenceRenderer struct{}

// RegisterFuncs implements renderer.NodeRenderer.
func (r *pluginFenceRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPluginFence, r.render)
}

//nolint:errcheck // Goldmark owns the writer lifecycle and handles write errors at the top level.
func (r *pluginFenceRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	n, ok := node.(*pluginFence)
	if !ok {
		return ast.WalkStop, fmt.Errorf("unexpected node type %T for plugin fence", node)
	}
	if !entering {
		_, _ = w.WriteString("</code></pre>\n")
		return ast.WalkContinue, nil
	}
	_, _ = w.WriteString(`<pre><code class="language-`)
	_, _ = w.Write(util.EscapeHTML(n.lang))
	_, _ = w.WriteString(`">`)
	for i := 0; i < n.lines.Len(); i++ {
		line := n.lines.At(i)
		_, _ = w.Write(util.EscapeHTML(line.Value(source)))
	}
	return ast.WalkContinue, nil
}

// codeFenceTransformer copies title/hl_lines shorthand from a fence info
// string onto the FencedCodeBlock node attributes, and swaps fences whose
// language belongs to another plugin for pluginFence nodes.
type codeFenceTransformer struct{}

// Transform implements parser.ASTTransformer.
func (t *codeFenceTransformer) Transform(node *ast.Document, reader text.Reader, _ parser.Context) {
	source := reader.Source()
	_ = ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) { //nolint:errcheck // callback always returns a nil error
		if !entering {
			return ast.WalkContinue, nil
		}
		fence, ok := n.(*ast.FencedCodeBlock)
		if !ok || fence.Info == nil {
			return ast.WalkContinue, nil
		}
		info := strings.TrimSpace(string(fence.Info.Segment.Value(source)))
		if lang := strings.ToLower(string(fence.Language(source))); pluginFenceLanguages[lang] {
			replacement := &pluginFence{lang: []byte(lang), lines: fence.Lines()}
			if parent := fence.Parent(); parent != nil {
				parent.ReplaceChild(parent, fence, replacement)
			}
			return ast.WalkSkipChildren, nil
		}
		applyCodeFenceInfo(fence, info)
		return ast.WalkContinue, nil
	})
}

// applyCodeFenceInfo parses the extras after the language token and sets
// attributes on the node. Attributes already present (from goldmark's own
// {key=value} parsing) are preserved.
//
//nolint:gocyclo // The supported fence syntaxes are independent parsing cases.
func applyCodeFenceInfo(fence *ast.FencedCodeBlock, info string) {
	if info == "" {
		return
	}
	rest := info
	if i := strings.IndexAny(info, " \t{"); i >= 0 {
		rest = strings.TrimSpace(info[i:])
	} else {
		return
	}
	if rest == "" {
		return
	}

	// Bare {1,3-5} shorthand → hl_lines
	if m := codeFenceRangePattern.FindStringSubmatch(rest); m != nil {
		if _, exists := fence.Attribute([]byte("hl_lines")); !exists {
			fence.SetAttribute([]byte("hl_lines"), parseHighlightRanges(m[1]))
		}
		return
	}

	// Brace-wrapped attribute block: let goldmark parse it if the core parser
	// did not already (it only does so when the block is the whole info tail).
	if strings.HasPrefix(rest, "{") {
		if attrs, ok := parser.ParseAttributes(text.NewReader([]byte(rest))); ok {
			for _, attr := range attrs {
				if _, exists := fence.Attribute(attr.Name); !exists {
					fence.SetAttribute(attr.Name, attr.Value)
				}
			}
			// A trailing bare range after the block, e.g. {title="x"} {1,3}
			if idx := strings.LastIndex(rest, "}"); idx >= 0 && idx < len(rest)-1 {
				tail := strings.TrimSpace(rest[idx+1:])
				if m := codeFenceRangePattern.FindStringSubmatch(tail); m != nil {
					if _, exists := fence.Attribute([]byte("hl_lines")); !exists {
						fence.SetAttribute([]byte("hl_lines"), parseHighlightRanges(m[1]))
					}
				}
			}
			return
		}
	}

	// Docusaurus/MkDocs-style key="value" pairs, optionally followed by {1,3-5}
	for _, m := range codeFenceKVPattern.FindAllStringSubmatch(rest, -1) {
		key := m[1]
		value := m[3]
		if value == "" {
			value = m[4]
		}
		if value == "" {
			value = m[5]
		}
		if _, exists := fence.Attribute([]byte(key)); exists {
			continue
		}
		switch key {
		case "hl_lines", "highlight":
			fence.SetAttribute([]byte("hl_lines"), parseHighlightRanges(strings.Trim(value, "{}")))
		default:
			fence.SetAttribute([]byte(key), []byte(value))
		}
	}
	if idx := strings.Index(rest, "{"); idx >= 0 {
		if m := codeFenceRangePattern.FindStringSubmatch(strings.TrimSpace(rest[idx:])); m != nil {
			if _, exists := fence.Attribute([]byte("hl_lines")); !exists {
				fence.SetAttribute([]byte("hl_lines"), parseHighlightRanges(m[1]))
			}
		}
	}
}

// parseHighlightRanges turns "1, 3-5" into the []interface{} shape the
// highlighting extension expects ([]uint8 range strings).
func parseHighlightRanges(spec string) []interface{} {
	parts := strings.Split(spec, ",")
	out := make([]interface{}, 0, len(parts))
	for _, part := range parts {
		part = strings.ReplaceAll(strings.TrimSpace(part), " ", "")
		if part == "" {
			continue
		}
		if _, err := strconv.Atoi(strings.Split(part, "-")[0]); err != nil {
			continue
		}
		out = append(out, []uint8(part))
	}
	return out
}

// codeFenceDisplayLang maps lexer names to friendlier badge labels.
var codeFenceDisplayLang = map[string]string{
	"js": "javascript", "ts": "typescript", "py": "python", "sh": "shell",
	"yml": "yaml", "md": "markdown", "rb": "ruby", "rs": "rust", "kt": "kotlin",
	"console": "shell", "shell-session": "shell", "plaintext": "", "text": "", "txt": "",
}

// writeCodeFenceHeader emits the static portion of the code chrome (title and
// language badge). The copy button is added client-side by reading.js.
//
//nolint:errcheck // Goldmark owns the writer lifecycle and handles write errors at the top level.
func writeCodeFenceHeader(w util.BufWriter, lang, title string) {
	display := strings.ToLower(lang)
	if mapped, ok := codeFenceDisplayLang[display]; ok {
		display = mapped
	}
	if title == "" && (display == "shell" || display == "bash") {
		title = "Terminal"
	}
	if display == "" && title == "" {
		return
	}
	_, _ = w.WriteString(`<div class="code-block__header"><div class="code-block__meta">`)
	if title != "" {
		_, _ = w.WriteString(`<span class="code-block__title">`)
		_, _ = w.Write(util.EscapeHTML([]byte(title)))
		_, _ = w.WriteString(`</span>`)
	}
	if display != "" {
		_, _ = w.WriteString(`<span class="code-block__lang">`)
		_, _ = w.Write(util.EscapeHTML([]byte(display)))
		_, _ = w.WriteString(`</span>`)
	}
	_, _ = w.WriteString(`</div></div>`)
}

// codeFenceWrapper renders the wrapper around highlighted code blocks with
// data attributes consumed by the theme's code chrome.
//
//nolint:errcheck // Goldmark owns the writer lifecycle and handles write errors at the top level.
func codeFenceWrapper(w util.BufWriter, c highlighting.CodeBlockContext, entering bool) {
	lang := ""
	if l, ok := c.Language(); ok {
		lang = string(l)
	}
	title := ""
	if attrs := c.Attributes(); attrs != nil {
		if v, ok := attrs.Get([]byte("title")); ok {
			switch tv := v.(type) {
			case []uint8:
				title = string(tv)
			case string:
				title = tv
			}
		}
	}
	display := strings.ToLower(lang)
	if mapped, ok := codeFenceDisplayLang[display]; ok {
		display = mapped
	}
	if title == "" && (display == "shell" || display == "bash") {
		title = "Terminal"
	}

	if !entering {
		if !c.Highlighted() {
			_, _ = w.WriteString("</code></pre>")
		}
		_, _ = w.WriteString("</div>\n")
		return
	}

	_, _ = w.WriteString(`<div class="code-block"`)
	if lang != "" {
		_, _ = w.WriteString(` data-lang="`)
		_, _ = w.Write(util.EscapeHTML([]byte(lang)))
		_, _ = w.WriteString(`"`)
	}
	if title != "" {
		_, _ = w.WriteString(` data-title="`)
		_, _ = w.Write(util.EscapeHTML([]byte(title)))
		_, _ = w.WriteString(`"`)
	}
	_, _ = w.WriteString(">")
	writeCodeFenceHeader(w, lang, title)
	if !c.Highlighted() {
		_, _ = w.WriteString("<pre><code")
		if len(lang) > 0 {
			_, _ = w.WriteString(` class="language-`)
			_, _ = w.Write(util.EscapeHTML([]byte(lang)))
			_, _ = w.WriteString(`"`)
		}
		_, _ = w.WriteString(">")
	}
}
