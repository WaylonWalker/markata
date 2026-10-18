package terminalpage

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"regexp"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"golang.org/x/net/html"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

var (
	spacePattern = regexp.MustCompile(`\s+`)
	ansiPattern  = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

const summaryTag = "summary"

const (
	DoubleRule = "━"
	SingleRule = "─"
)

type Options struct {
	ANSI        bool
	Palette     string
	ChromaStyle string
}

type Renderer struct {
	options Options
	theme   theme
}

type theme struct {
	headline  ansiStyle
	muted     ansiStyle
	link      ansiStyle
	quote     ansiStyle
	code      ansiStyle
	border    ansiStyle
	success   ansiStyle
	warning   ansiStyle
	danger    ansiStyle
	info      ansiStyle
	strong    ansiStyle
	emph      ansiStyle
	underline ansiStyle
	strike    ansiStyle
}

type ansiStyle struct {
	prefix string
	suffix string
}

func New(options Options) *Renderer {
	return &Renderer{
		options: options,
		theme:   buildTheme(options),
	}
}

func RenderHTML(src string, options Options) string {
	return New(options).Render(src)
}

func StripANSI(src string) string {
	return ansiPattern.ReplaceAllString(src, "")
}

func (r *Renderer) Render(src string) string {
	trimmed := strings.TrimSpace(src)
	if trimmed == "" {
		return ""
	}

	doc, err := html.Parse(strings.NewReader("<html><body><div>" + trimmed + "</div></body></html>"))
	if err != nil {
		return strings.TrimSpace(stdhtml.UnescapeString(trimmed))
	}
	root := findTerminalRoot(doc)
	if root == nil {
		return strings.TrimSpace(stdhtml.UnescapeString(trimmed))
	}

	blocks := []string{}
	for node := root.FirstChild; node != nil; node = node.NextSibling {
		for _, block := range r.renderBlocks(node) {
			block = strings.TrimSpace(block)
			if block != "" {
				blocks = append(blocks, block)
			}
		}
	}

	return strings.TrimSpace(strings.Join(blocks, "\n\n"))
}

//nolint:gocyclo // HTML block dispatch is clearer as a single switch.
func (r *Renderer) renderBlocks(node *html.Node) []string {
	if node == nil {
		return nil
	}

	if node.Type == html.TextNode {
		text := normalizeText(node.Data)
		if text == "" {
			return nil
		}
		return []string{text}
	}

	if node.Type != html.ElementNode {
		return r.renderChildBlocks(node)
	}

	switch node.Data {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return []string{r.renderHeading(node)}
	case "p":
		text := strings.TrimSpace(r.renderInlineChildren(node))
		if text == "" {
			return nil
		}
		return []string{text}
	case "blockquote":
		return []string{r.renderBlockquote(node)}
	case "pre":
		return []string{r.renderCodeBlock(node)}
	case "ul":
		return []string{r.renderList(node, false)}
	case "ol":
		return []string{r.renderList(node, true)}
	case "table":
		return []string{r.renderTable(node)}
	case "hr":
		line := strings.Repeat("-", 8)
		return []string{r.theme.border.wrap(line)}
	case "details":
		if hasClass(node, "admonition") {
			return []string{r.renderAdmonition(node)}
		}
		return []string{r.renderDetails(node)}
	case "div", "section", "article", "header", "footer", "aside", "main":
		if hasClass(node, "admonition") {
			return []string{r.renderAdmonition(node)}
		}
		if hasClass(node, "chroma") {
			return []string{r.renderCodeBlock(node)}
		}
		return r.renderChildBlocks(node)
	case "figure":
		return r.renderChildBlocks(node)
	case "br":
		return []string{""}
	default:
		text := strings.TrimSpace(r.renderInlineChildren(node))
		if text != "" {
			return []string{text}
		}
		return r.renderChildBlocks(node)
	}
}

func (r *Renderer) renderChildBlocks(node *html.Node) []string {
	blocks := []string{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		blocks = append(blocks, r.renderBlocks(child)...)
	}
	return blocks
}

func (r *Renderer) renderInlineChildren(node *html.Node) string {
	var parts []string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		part := r.renderInlineNode(child)
		if part != "" {
			parts = append(parts, part)
		}
	}
	return strings.TrimSpace(collapseInlineWhitespace(strings.Join(parts, "")))
}

//nolint:gocyclo // Inline HTML handling is easier to reason about as one switch.
func (r *Renderer) renderInlineNode(node *html.Node) string {
	if node == nil {
		return ""
	}

	switch node.Type {
	case html.TextNode:
		return normalizeText(node.Data)
	case html.ElementNode:
		content := r.renderInlineChildren(node)
		switch node.Data {
		case "strong", "b":
			return r.theme.strong.wrap(content)
		case "em", "i":
			return r.theme.emph.wrap(content)
		case "u":
			return r.theme.underline.wrap(content)
		case "s", "del":
			return r.theme.strike.wrap(content)
		case "code":
			if node.Parent != nil && node.Parent.Data == "pre" {
				return ""
			}
			return r.theme.code.wrap(content)
		case "a":
			href := getAttr(node, "href")
			if shouldSkipAnchorLink(node, href) {
				return ""
			}
			if href == "" {
				return content
			}
			if content == "" {
				content = href
			}
			if sameLinkText(content, href) {
				return r.theme.link.wrap(content)
			}
			return r.theme.link.wrap(content) + " <" + r.theme.muted.wrap(href) + ">"
		case "img":
			return renderImageReference(node)
		case "video":
			return renderMediaReference(node, "Video")
		case "audio":
			return renderMediaReference(node, "Audio")
		case "source":
			return ""
		case "br":
			return "\n"
		default:
			return content
		}
	default:
		return ""
	}
}

func (r *Renderer) renderHeading(node *html.Node) string {
	text := strings.TrimSpace(r.renderInlineChildren(node))
	if text == "" {
		return ""
	}

	level := 1
	if parsedLevel, err := strconv.Atoi(strings.TrimPrefix(node.Data, "h")); err == nil {
		level = parsedLevel
	}
	plain := StripANSI(text)

	switch level {
	case 1:
		return r.theme.headline.wrap(text) + "\n" + r.theme.border.wrap(strings.Repeat(DoubleRule, len([]rune(plain))))
	case 2:
		return r.theme.headline.wrap(text) + "\n" + r.theme.border.wrap(strings.Repeat(SingleRule, len([]rune(plain))))
	default:
		prefix := strings.Repeat("#", level) + " "
		return r.theme.headline.wrap(prefix + text)
	}
}

func (r *Renderer) renderBlockquote(node *html.Node) string {
	content := joinBlocks(r.renderChildBlocks(node))
	if content == "" {
		content = strings.TrimSpace(r.renderInlineChildren(node))
	}
	return prefixLines(content, r.theme.quote.wrap("│ "))
}

func (r *Renderer) renderList(node *html.Node, ordered bool) string {
	items := []string{}
	index := 1
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode || child.Data != "li" {
			continue
		}
		marker := "- "
		if ordered {
			marker = fmt.Sprintf("%d. ", index)
		}
		items = append(items, r.renderListItem(child, marker))
		index++
	}
	return strings.Join(items, "\n")
}

func (r *Renderer) renderListItem(node *html.Node, marker string) string {
	parts := []string{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parts = append(parts, r.renderBlocks(child)...)
	}
	content := joinBlocks(parts)
	if content == "" {
		content = strings.TrimSpace(r.renderInlineChildren(node))
	}
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " ")
	}
	if len(lines) == 0 {
		return marker
	}
	result := marker + lines[0]
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" {
			result += "\n"
			continue
		}
		result += "\n" + strings.Repeat(" ", len(marker)) + line
	}
	return result
}

func (r *Renderer) renderCodeBlock(node *html.Node) string {
	codeNode := findDescendant(node, "code")
	lang := detectLanguage(node)
	if codeNode != nil && lang == "" {
		lang = detectLanguage(codeNode)
	}

	code := extractText(node, true)
	code = strings.Trim(code, "\n")
	if code == "" {
		return ""
	}

	if !r.options.ANSI {
		fence := "```"
		if lang != "" {
			fence += lang
		}
		return fence + "\n" + code + "\n```"
	}

	highlighted := r.highlightCode(code, lang)
	label := r.theme.muted.wrap("[code]")
	if lang != "" {
		label = r.theme.muted.wrap("[" + lang + "]")
	}
	return label + "\n" + prefixLines(highlighted, "  ")
}

//nolint:gocyclo // Table extraction and formatting stay together for readability.
func (r *Renderer) renderTable(node *html.Node) string {
	rows := [][]string{}
	headerRows := 0
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}
		if child.Data == "thead" {
			headerRows = countRows(child)
		}
		if child.Data == "tbody" || child.Data == "thead" || child.Data == "tfoot" {
			for tr := child.FirstChild; tr != nil; tr = tr.NextSibling {
				if tr.Type == html.ElementNode && tr.Data == "tr" {
					rows = append(rows, collectCells(r, tr))
				}
			}
		}
		if child.Data == "tr" {
			rows = append(rows, collectCells(r, child))
		}
	}
	if len(rows) == 0 {
		return ""
	}

	widths := make([]int, 0, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			for len(widths) <= i {
				widths = append(widths, 0)
			}
			if w := len([]rune(StripANSI(cell))); w > widths[i] {
				widths[i] = w
			}
		}
	}

	separator := buildTableSeparator(widths)
	lines := []string{separator}
	for i, row := range rows {
		cells := make([]string, len(widths))
		for j := range widths {
			value := ""
			if j < len(row) {
				value = row[j]
			}
			if i < headerRows {
				value = r.theme.strong.wrap(value)
			}
			cells[j] = padRightANSI(value, widths[j])
		}
		lines = append(lines, "| "+strings.Join(cells, " | ")+" |")
		if i == headerRows-1 || (i == 0 && headerRows == 0) {
			lines = append(lines, separator)
		}
	}
	lines = append(lines, separator)
	return strings.Join(lines, "\n")
}

func (r *Renderer) renderAdmonition(node *html.Node) string {
	kind := admonitionKind(node)
	title := strings.TrimSpace(findAdmonitionTitle(r, node))
	content := strings.TrimSpace(findAdmonitionBody(r, node))
	if content == "" {
		content = strings.TrimSpace(r.renderInlineChildren(node))
	}

	labelText := strings.ToUpper(kind)
	if labelText == "" {
		labelText = "NOTE"
	}
	labelStyle := r.admonitionStyle(kind)
	header := labelStyle.wrap(labelText)
	if title != "" && !strings.EqualFold(title, labelText) {
		header += " " + r.theme.strong.wrap(title)
	}
	if node.Data == "details" && getAttr(node, "open") == "" {
		header += " " + r.theme.muted.wrap("(collapsed by default)")
	}
	if content == "" {
		return header
	}
	return header + "\n" + prefixLines(content, labelStyle.wrap("│ "))
}

func (r *Renderer) renderDetails(node *html.Node) string {
	summary := "Details"
	if summaryNode := findDescendant(node, summaryTag); summaryNode != nil {
		summary = strings.TrimSpace(r.renderInlineChildren(summaryNode))
	}
	content := []string{r.theme.strong.wrap(summary)}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == summaryTag {
			continue
		}
		content = append(content, r.renderBlocks(child)...)
	}
	return joinBlocks(content)
}

func (r *Renderer) highlightCode(code, lang string) string {
	lexer := lexers.Get(lang)
	if lexer == nil {
		//nolint:misspell // Chroma exposes Analyse with British spelling.
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		return code
	}

	styleName := r.options.ChromaStyle
	if styleName == "" {
		styleName = palettes.DefaultChromaThemeDark
	}
	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	if err != nil {
		return code
	}

	var buf bytes.Buffer
	if err := formatters.TTY16m.Format(&buf, style, iterator); err != nil {
		return code
	}
	return strings.TrimRight(buf.String(), "\n")
}

func (r *Renderer) admonitionStyle(kind string) ansiStyle {
	switch kind {
	case "success", "tip", "hint":
		return r.theme.success
	case "warning", "caution", "important", "attention":
		return r.theme.warning
	case "danger", "error", "bug":
		return r.theme.danger
	case "info", "abstract", "example", "seealso", "reminder":
		return r.theme.info
	default:
		return r.theme.border
	}
}

func buildTheme(options Options) theme {
	base := theme{
		headline:  makeStyle(options.ANSI, "#7c3aed", true, false, false, false),
		muted:     makeStyle(options.ANSI, "#6b7280", false, false, false, false),
		link:      makeStyle(options.ANSI, "#2563eb", false, false, true, false),
		quote:     makeStyle(options.ANSI, "#6b7280", false, false, false, false),
		code:      makeStyle(options.ANSI, "#dc2626", false, false, false, false),
		border:    makeStyle(options.ANSI, "#94a3b8", false, false, false, false),
		success:   makeStyle(options.ANSI, "#16a34a", true, false, false, false),
		warning:   makeStyle(options.ANSI, "#d97706", true, false, false, false),
		danger:    makeStyle(options.ANSI, "#dc2626", true, false, false, false),
		info:      makeStyle(options.ANSI, "#0f766e", true, false, false, false),
		strong:    makeStyle(options.ANSI, "", true, false, false, false),
		emph:      makeStyle(options.ANSI, "", false, true, false, false),
		underline: makeStyle(options.ANSI, "", false, false, true, false),
		strike:    makeStyle(options.ANSI, "", false, false, false, true),
	}

	if options.Palette == "" {
		return base
	}

	loader := palettes.NewLoader()
	palette, err := loader.Load(options.Palette)
	if err != nil || palette == nil {
		return base
	}

	if hex := palette.Resolve("accent"); hex != "" {
		base.headline = makeStyle(options.ANSI, hex, true, false, false, false)
	}
	if hex := palette.Resolve("text-muted"); hex != "" {
		base.muted = makeStyle(options.ANSI, hex, false, false, false, false)
		base.quote = base.muted
	}
	if hex := palette.Resolve("link"); hex != "" {
		base.link = makeStyle(options.ANSI, hex, false, false, true, false)
	}
	if hex := palette.Resolve("success"); hex != "" {
		base.success = makeStyle(options.ANSI, hex, true, false, false, false)
	}
	if hex := palette.Resolve("warning"); hex != "" {
		base.warning = makeStyle(options.ANSI, hex, true, false, false, false)
	}
	if hex := palette.Resolve("error"); hex != "" {
		base.danger = makeStyle(options.ANSI, hex, true, false, false, false)
	}
	if hex := palette.Resolve("info"); hex != "" {
		base.info = makeStyle(options.ANSI, hex, true, false, false, false)
	}
	if hex := palette.Resolve("border"); hex != "" {
		base.border = makeStyle(options.ANSI, hex, false, false, false, false)
	}
	if hex := palette.Resolve("accent"); hex != "" {
		base.code = makeStyle(options.ANSI, hex, false, false, false, false)
	}

	return base
}

func makeStyle(enabled bool, fg string, bold, italic, underline, strike bool) ansiStyle {
	if !enabled {
		return ansiStyle{}
	}
	parts := []string{}
	if bold {
		parts = append(parts, "\x1b[1m")
	}
	if italic {
		parts = append(parts, "\x1b[3m")
	}
	if underline {
		parts = append(parts, "\x1b[4m")
	}
	if strike {
		parts = append(parts, "\x1b[9m")
	}
	if fg != "" {
		parts = append(parts, colorSequence(fg, false))
	}
	if len(parts) == 0 {
		return ansiStyle{}
	}
	return ansiStyle{prefix: strings.Join(parts, ""), suffix: "\x1b[0m"}
}

func colorSequence(hex string, background bool) string {
	parsed, err := palettes.ParseHexColor(hex)
	if err != nil {
		return ""
	}
	mode := 38
	if background {
		mode = 48
	}
	return fmt.Sprintf("\x1b[%d;2;%d;%d;%dm", mode, parsed.R, parsed.G, parsed.B)
}

func (s ansiStyle) wrap(text string) string {
	if s.prefix == "" || text == "" {
		return text
	}
	return s.prefix + text + s.suffix
}

func normalizeText(text string) string {
	text = stdhtml.UnescapeString(text)
	if strings.TrimSpace(text) == "" {
		if strings.Contains(text, "\n") {
			return "\n"
		}
		return " "
	}
	return collapseInlineWhitespace(text)
}

func collapseInlineWhitespace(text string) string {
	text = strings.ReplaceAll(text, "\u00a0", " ")
	return spacePattern.ReplaceAllString(text, " ")
}

func sameLinkText(text, href string) bool {
	return strings.TrimSuffix(text, "/") == strings.TrimSuffix(href, "/")
}

func shouldSkipAnchorLink(node *html.Node, href string) bool {
	if href == "" {
		return false
	}
	if hasClass(node, "anchor") || hasClass(node, "heading-anchor") {
		return true
	}
	return strings.HasPrefix(href, "#") && isHeadingTag(node.Parent)
}

func isHeadingTag(node *html.Node) bool {
	if node == nil || node.Type != html.ElementNode {
		return false
	}
	switch node.Data {
	case "h1", "h2", "h3", "h4", "h5", "h6":
		return true
	default:
		return false
	}
}

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func renderImageReference(node *html.Node) string {
	label := firstNonEmpty(getAttr(node, "alt"), getAttr(node, "title"), "Image")
	src := firstNonEmpty(getAttr(node, "src"), getAttr(node, "data-src"))
	if src == "" {
		return "Image: " + label
	}
	if label == "Image" {
		return "Image: <" + src + ">"
	}
	return "Image: " + label + " <" + src + ">"
}

func renderMediaReference(node *html.Node, fallback string) string {
	label := firstNonEmpty(getAttr(node, "aria-label"), getAttr(node, "title"), fallback)
	src := firstNonEmpty(getAttr(node, "src"), getAttr(node, "data-src"), mediaSourceFromChildren(node))
	if src == "" {
		return fallback + ": " + label
	}
	if label == fallback {
		return fallback + ": <" + src + ">"
	}
	return fallback + ": " + label + " <" + src + ">"
}

func mediaSourceFromChildren(node *html.Node) string {
	if node == nil {
		return ""
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "source" {
			if src := firstNonEmpty(getAttr(child, "src"), getAttr(child, "data-src")); src != "" {
				return src
			}
		}
		if src := mediaSourceFromChildren(child); src != "" {
			return src
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func hasClass(node *html.Node, class string) bool {
	classes := strings.Fields(getAttr(node, "class"))
	for _, candidate := range classes {
		if candidate == class {
			return true
		}
	}
	return false
}

func joinBlocks(blocks []string) string {
	trimmed := []string{}
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block != "" {
			trimmed = append(trimmed, block)
		}
	}
	return strings.Join(trimmed, "\n\n")
}

func prefixLines(text, prefix string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line == "" {
			lines[i] = strings.TrimRight(prefix, " ")
			continue
		}
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
}

func extractText(node *html.Node, preserveWhitespace bool) string {
	if node == nil {
		return ""
	}
	if node.Type == html.TextNode {
		if preserveWhitespace {
			return stdhtml.UnescapeString(node.Data)
		}
		return normalizeText(node.Data)
	}
	var builder strings.Builder
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		builder.WriteString(extractText(child, preserveWhitespace))
	}
	return builder.String()
}

func detectLanguage(node *html.Node) string {
	class := getAttr(node, "class")
	for _, item := range strings.Fields(class) {
		if strings.HasPrefix(item, "language-") {
			return strings.TrimPrefix(item, "language-")
		}
	}
	if dataLang := getAttr(node, "data-language"); dataLang != "" {
		return dataLang
	}
	return ""
}

func findDescendant(node *html.Node, tag string) *html.Node {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == tag {
			return child
		}
		if found := findDescendant(child, tag); found != nil {
			return found
		}
	}
	return nil
}

func findTerminalRoot(node *html.Node) *html.Node {
	if node == nil {
		return nil
	}
	if node.Type == html.ElementNode && node.Data == "div" && node.Parent != nil && node.Parent.Data == "body" {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findTerminalRoot(child); found != nil {
			return found
		}
	}
	return nil
}

func countRows(node *html.Node) int {
	count := 0
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "tr" {
			count++
		}
	}
	return count
}

func collectCells(r *Renderer, tr *html.Node) []string {
	row := []string{}
	for cell := tr.FirstChild; cell != nil; cell = cell.NextSibling {
		if cell.Type != html.ElementNode || (cell.Data != "td" && cell.Data != "th") {
			continue
		}
		value := strings.TrimSpace(r.renderInlineChildren(cell))
		if value == "" {
			value = strings.TrimSpace(joinBlocks(r.renderChildBlocks(cell)))
		}
		row = append(row, value)
	}
	return row
}

func buildTableSeparator(widths []int) string {
	parts := make([]string, len(widths))
	for i, width := range widths {
		parts[i] = strings.Repeat("-", width+2)
	}
	return "+" + strings.Join(parts, "+") + "+"
}

func padRightANSI(text string, width int) string {
	pad := width - len([]rune(StripANSI(text)))
	if pad <= 0 {
		return text
	}
	return text + strings.Repeat(" ", pad)
}

func admonitionKind(node *html.Node) string {
	for _, class := range strings.Fields(getAttr(node, "class")) {
		if class != "admonition" {
			return class
		}
	}
	return "note"
}

func findAdmonitionTitle(r *Renderer, node *html.Node) string {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type != html.ElementNode {
			continue
		}
		if child.Data == summaryTag || hasClass(child, "admonition-title") {
			return r.renderInlineChildren(child)
		}
	}
	return ""
}

func findAdmonitionBody(r *Renderer, node *html.Node) string {
	parts := []string{}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && (child.Data == summaryTag || hasClass(child, "admonition-title")) {
			continue
		}
		parts = append(parts, r.renderBlocks(child)...)
	}
	return joinBlocks(parts)
}
