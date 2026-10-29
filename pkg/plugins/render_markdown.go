// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bytes"
	stdhtml "html"
	"regexp"
	"strings"
	"sync"
	"unicode"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	emojiast "github.com/yuin/goldmark-emoji/ast"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	figure "github.com/mangoumbrella/goldmark-figure"
	"go.abhg.dev/goldmark/anchor"
)

// markdownBufferPool is a sync.Pool for reusing bytes.Buffer instances
// during markdown rendering. This significantly reduces allocations and
// GC pressure when processing many posts.
var markdownBufferPool = sync.Pool{
	New: func() interface{} {
		// Pre-allocate 32KB buffer - typical blog post is 5-20KB HTML
		buf := bytes.NewBuffer(make([]byte, 0, 32*1024))
		return buf
	},
}

var figureBlockquoteCaptionParagraphRegex = regexp.MustCompile(`(?s)^<p>.*?</p>$`)

// attributionMarker is an HTML comment inserted before attribution lines
// in the markdown pre-processing step, and detected in the post-processing step.
const attributionMarker = "<!--markata-attribution-->"

const disabledHighlightTheme = "monokailight"

// RenderMarkdownPlugin converts markdown content to HTML using goldmark.
type RenderMarkdownPlugin struct {
	md       goldmark.Markdown
	inlineMD goldmark.Markdown
	cache    *buildcache.Cache // build cache for HTML caching
}

// CacheKeyMarkdownRenderer is the manager cache key for the markdown render
// function.  Other plugins (e.g. feed helpers) can retrieve and call this to
// render markdown on-demand when ArticleHTML has not yet been populated.
const CacheKeyMarkdownRenderer = "markdown.renderer"

// headingHighlightRevision migrates only pages whose canonical articles contain
// marked H1/H2 headings. Certification belongs to the final page cache writer.
const headingHighlightRevision = 1

// canonicalArticlePathsKey tracks which articles were actually restored or
// rendered this pass, so an omitted incremental post cannot be certified.
const canonicalArticlePathsKey = "markdown.canonical_article_paths"

// MarkdownRenderFunc is the signature of the on-demand markdown renderer
// stored under CacheKeyMarkdownRenderer.
type MarkdownRenderFunc func(content string) (string, error)

// InlineResult contains the two representations needed when inline Markdown
// is used outside an article body. HTML is for trusted display contexts; Text
// is for metadata, feeds, search, and other plain-text consumers.
type InlineResult struct {
	HTML string
	Text string
}

// InlineRenderFunc renders an inline Markdown fragment with the site's
// configured Markdown extensions and safe HTML policy.
type InlineRenderFunc func(source string) (InlineResult, error)

// CacheKeyInlineRenderer is the manager cache key for the shared inline
// renderer used by title processing and other plugins.
const CacheKeyInlineRenderer = "markdown.inline_renderer"

// NewRenderMarkdownPlugin creates a new RenderMarkdownPlugin with goldmark configured.
// The goldmark instance is configured with:
// - GFM extensions (tables, strikethrough, autolinks, task lists)
// - Syntax highlighting using chroma
// - HTML rendering with unsafe mode (allow raw HTML)
// - Auto heading IDs
// - Custom admonition support
//
// The initial configuration uses a default theme. The Configure() method will
// reconfigure the markdown renderer with the appropriate theme based on the
// site's palette configuration.
func NewRenderMarkdownPlugin() *RenderMarkdownPlugin {
	return &RenderMarkdownPlugin{
		md:       createMarkdownRenderer(palettes.DefaultChromaThemeDark, false, DefaultMarkdownExtensionConfig(), true),
		inlineMD: createMarkdownRenderer(palettes.DefaultChromaThemeDark, false, DefaultMarkdownExtensionConfig(), false),
	}
}

// MarkdownExtensionConfig holds configuration for optional markdown extensions.
type MarkdownExtensionConfig struct {
	TypographerEnabled       bool
	DefinitionListEnabled    bool
	FootnoteEnabled          bool
	CJKEnabled               bool
	FigureEnabled            bool
	AnchorEnabled            bool
	TypographerSubstitutions map[extension.TypographicPunctuation]string
}

// DefaultMarkdownExtensionConfig returns the default configuration with all extensions enabled.
func DefaultMarkdownExtensionConfig() MarkdownExtensionConfig {
	return MarkdownExtensionConfig{
		TypographerEnabled:       true,
		DefinitionListEnabled:    true,
		FootnoteEnabled:          true,
		CJKEnabled:               true,
		FigureEnabled:            true,
		AnchorEnabled:            false,
		TypographerSubstitutions: nil, // nil means use goldmark defaults
	}
}

// createMarkdownRenderer creates a goldmark instance with the specified highlighting options.
func createMarkdownRenderer(chromaTheme string, lineNumbers bool, extConfig MarkdownExtensionConfig, unsafeOptions ...bool) goldmark.Markdown {
	unsafeHTML := true
	if len(unsafeOptions) > 0 {
		unsafeHTML = unsafeOptions[0]
	}
	// Use CSS classes instead of inline styles for syntax highlighting.
	// This enables theme customization via external CSS files.
	formatOptions := []chromahtml.Option{
		chromahtml.WithClasses(true),
		chromahtml.WithAllClasses(true),
	}

	if lineNumbers {
		formatOptions = append(formatOptions, chromahtml.WithLineNumbers(true))
	}

	highlightOpts := []highlighting.Option{
		highlighting.WithStyle(chromaTheme),
		highlighting.WithFormatOptions(formatOptions...),
	}

	extensions := []goldmark.Extender{
		// GFM extensions
		extension.GFM,
		extension.Table,
		extension.Strikethrough,
		extension.Linkify,
		extension.TaskList,
		// Syntax highlighting with chroma
		highlighting.NewHighlighting(highlightOpts...),
		// Custom admonition extension
		&AdmonitionExtension{},
		// Mark extension for ==highlighted text==
		&MarkExtension{},
		&ScriptExtension{},
		// Keys extension for ++Ctrl+Alt+Del++
		&KeysExtension{},
		// Container extension for ::: class
		&ContainerExtension{},
		// Emoji extension for :smile: syntax
		emoji.Emoji,
	}

	// Add CJK extension (Chinese/Japanese/Korean line break support)
	if extConfig.CJKEnabled {
		extensions = append(extensions, extension.NewCJK())
	}

	// Add Figure extension for <figure> elements from images with captions
	if extConfig.FigureEnabled {
		extensions = append(extensions, figure.Figure)
	}

	// Add Anchor extension for heading permalinks
	if extConfig.AnchorEnabled {
		extensions = append(extensions, &anchor.Extender{})
	}

	// Add Typographer extension (smart quotes, dashes, ellipses)
	if extConfig.TypographerEnabled {
		if extConfig.TypographerSubstitutions != nil {
			extensions = append(extensions, extension.NewTypographer(
				extension.WithTypographicSubstitutions(extConfig.TypographerSubstitutions),
			))
		} else {
			extensions = append(extensions, extension.Typographer)
		}
	}

	// Add DefinitionList extension (PHP Markdown Extra style definition lists)
	if extConfig.DefinitionListEnabled {
		extensions = append(extensions, extension.DefinitionList)
	}

	// Add Footnote extension (PHP Markdown Extra style footnotes)
	if extConfig.FootnoteEnabled {
		extensions = append(extensions, extension.Footnote)
	}

	options := []goldmark.Option{
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			// Enable attribute syntax for {.class}, {#id}, {key=value} on block elements
			parser.WithAttribute(),
			// Enable inline attribute syntax for inline and block elements
			parser.WithASTTransformers(
				util.Prioritized(&AttributeTransformer{}, -100),
				util.Prioritized(&HeadingIDTransformer{}, -90),
			),
		),
	}
	if unsafeHTML {
		options = append(options, goldmark.WithRendererOptions(html.WithUnsafe()))
	}
	return goldmark.New(options...)
}

// Name returns the unique name of the plugin.
func (p *RenderMarkdownPlugin) Name() string {
	return "render_markdown"
}

// Configure initializes the plugin during the configure stage.
// It reads the highlight theme from configuration, falling back to the theme
// derived from the site's color palette if not explicitly set.
//
// Configuration priority:
// 1. markdown.highlight.theme in config (explicit override)
// 2. Theme derived from theme.palette (automatic matching)
// 3. Default theme based on palette variant (github/github-dark)
//
// Markdown extension configuration:
// 1. markdown.extensions.typographer - Enable smart quotes (default: true)
// 2. markdown.extensions.definition_list - Enable definition lists (default: true)
// 3. markdown.extensions.footnote - Enable footnotes (default: true)
// 4. markdown.extensions.cjk - Enable CJK line breaks (default: true)
// 5. markdown.extensions.figure - Enable figure from images with captions (default: true)
// 6. markdown.extensions.anchor - Enable heading permalinks (default: true)
func (p *RenderMarkdownPlugin) Configure(m *lifecycle.Manager) error {
	chromaTheme, lineNumbers := p.resolveHighlightConfig(m.Config().Extra)
	extConfig := p.resolveExtensionConfig(m.Config().Extra)

	// Reconfigure the markdown renderer with the resolved theme and extensions
	p.md = createMarkdownRenderer(chromaTheme, lineNumbers, extConfig, true)
	p.inlineMD = createMarkdownRenderer(chromaTheme, false, extConfig, false)

	// Register the render function so other plugins (e.g. feed helpers during
	// jinja_md transform) can render markdown on-demand when ArticleHTML has
	// not yet been populated by the Render stage.
	m.Cache().Set(CacheKeyMarkdownRenderer, MarkdownRenderFunc(p.doRender))
	m.Cache().Set(CacheKeyInlineRenderer, InlineRenderFunc(p.renderInline))

	return nil
}

// renderInline renders a Markdown fragment without a paragraph wrapper. It
// parses through the same configured Goldmark instance as body content while
// deliberately keeping raw HTML disabled for this metadata-adjacent path.
func (p *RenderMarkdownPlugin) renderInline(source string) (InlineResult, error) {
	// A leading text token keeps tag-looking input in an inline paragraph
	// rather than letting Goldmark classify a leading HTML block. Raw HTML is
	// still disabled on this renderer, so the source cannot become executable.
	bytesSource := []byte("x " + source)
	document := p.inlineMD.Parser().Parse(text.NewReader(bytesSource))
	paragraph := document.FirstChild()
	if paragraph == nil {
		return InlineResult{}, nil
	}

	var htmlBuffer bytes.Buffer
	var plain strings.Builder
	for child := paragraph.FirstChild(); child != nil; child = child.NextSibling() {
		if err := p.inlineMD.Renderer().Render(&htmlBuffer, bytesSource, child); err != nil {
			return InlineResult{}, err
		}
		appendInlineText(&plain, child, bytesSource)
	}
	richHTML := strings.TrimPrefix(htmlBuffer.String(), "x ")
	plainText := strings.TrimPrefix(plain.String(), "x ")
	// Goldmark's typographer extension may represent smart punctuation as
	// HTML character references in the AST. TitleText is plain text, so decode
	// those references before exposing it to templates and metadata writers.
	plainText = stdhtml.UnescapeString(plainText)
	return InlineResult{HTML: richHTML, Text: strings.Join(strings.Fields(plainText), " ")}, nil
}

func appendInlineText(out *strings.Builder, node ast.Node, source []byte) {
	if err := ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Text:
			out.Write(n.Segment.Value(source))
		case *ast.String:
			out.Write(n.Value)
		case *ast.AutoLink:
			// AutoLink is a leaf node, so its visible label is not visited
			// as a child Text node. Use Goldmark's supported label API so
			// URL and email autolinks retain their human-visible text.
			out.Write(n.Label(source))
		case *emojiast.Emoji:
			if n.Value != nil {
				out.WriteString(string(n.Value.Unicode))
			}
		}
		return ast.WalkContinue, nil
	}); err != nil {
		return
	}
}

// resolveHighlightConfig extracts highlight configuration from the config.Extra map.
// Returns the Chroma theme name and whether line numbers should be shown.
func (p *RenderMarkdownPlugin) resolveHighlightConfig(extra map[string]interface{}) (string, bool) {
	var chromaTheme string
	var lineNumbers bool
	var paletteVariant = palettes.VariantDark

	if highlight, ok := highlightConfigFromExtra(extra); ok {
		if !highlight.IsEnabled() {
			// Use a neutral theme while preserving the existing renderer path.
			return disabledHighlightTheme, false
		}
		chromaTheme = highlight.Theme
		lineNumbers = highlight.LineNumbers
	}

	// If no explicit theme, derive from palette
	if chromaTheme == "" {
		paletteName := p.getPaletteName(extra)
		if paletteName != "" {
			// Try to get theme from palette mapping
			chromaTheme = palettes.ChromaTheme(paletteName)

			// Determine palette variant for fallback
			paletteVariant = p.getPaletteVariant(paletteName)
		}
	}

	// Final fallback based on palette variant
	if chromaTheme == "" {
		chromaTheme = palettes.ChromaThemeForVariant(paletteVariant)
	}

	return chromaTheme, lineNumbers
}

// highlightConfigFromExtra accepts the canonical typed config and the legacy
// map representation used by direct plugin callers.
func highlightConfigFromExtra(extra map[string]interface{}) (models.HighlightConfig, bool) {
	if extra == nil {
		return models.HighlightConfig{}, false
	}

	switch markdown := extra["markdown"].(type) {
	case models.MarkdownConfig:
		return markdown.Highlight, true
	case *models.MarkdownConfig:
		if markdown != nil {
			return markdown.Highlight, true
		}
	case map[string]interface{}:
		highlightValue, ok := markdown["highlight"]
		if !ok {
			return models.HighlightConfig{}, false
		}
		highlightMap, ok := highlightValue.(map[string]interface{})
		if !ok {
			return models.HighlightConfig{}, false
		}

		highlight := models.HighlightConfig{}
		if enabled, ok := highlightMap["enabled"].(bool); ok {
			highlight.Enabled = &enabled
		}
		if theme, ok := highlightMap["theme"].(string); ok {
			highlight.Theme = theme
		}
		if lineNumbers, ok := highlightMap["line_numbers"].(bool); ok {
			highlight.LineNumbers = lineNumbers
		}
		return highlight, true
	}

	return models.HighlightConfig{}, false
}

// getPaletteName extracts the palette name from config.Extra.
func (p *RenderMarkdownPlugin) getPaletteName(extra map[string]interface{}) string {
	if theme, ok := extra["theme"].(map[string]interface{}); ok {
		if palette, ok := theme["palette"].(string); ok && palette != "" {
			return palette
		}
	}
	return ""
}

// getPaletteVariant determines the variant (light/dark) of a palette by name.
// Uses naming conventions to infer variant when palette data isn't loaded.
func (p *RenderMarkdownPlugin) getPaletteVariant(paletteName string) palettes.Variant {
	// Check common light palette name patterns
	lightPatterns := []string{
		"-light", "-latte", "-dawn", "-day", "-lotus",
	}
	for _, pattern := range lightPatterns {
		if strings.Contains(paletteName, pattern) {
			return palettes.VariantLight
		}
	}
	return palettes.VariantDark
}

// resolveExtensionConfig extracts markdown extension configuration from config.Extra map.
// Returns a MarkdownExtensionConfig with settings for typographer, definition lists, and footnotes.
// All extensions are enabled by default unless explicitly disabled in config.
func (p *RenderMarkdownPlugin) resolveExtensionConfig(extra map[string]interface{}) MarkdownExtensionConfig {
	config := DefaultMarkdownExtensionConfig()

	// Try to get extension config from markdown.extensions
	if markdown, ok := extra["markdown"].(map[string]interface{}); ok {
		if extensions, ok := markdown["extensions"].(map[string]interface{}); ok {
			// Typographer (smart quotes)
			if enabled, ok := extensions["typographer"].(bool); ok {
				config.TypographerEnabled = enabled
			}

			// Definition lists
			if enabled, ok := extensions["definition_list"].(bool); ok {
				config.DefinitionListEnabled = enabled
			}

			// Footnotes
			if enabled, ok := extensions["footnote"].(bool); ok {
				config.FootnoteEnabled = enabled
			}

			// CJK (Chinese/Japanese/Korean line breaks)
			if enabled, ok := extensions["cjk"].(bool); ok {
				config.CJKEnabled = enabled
			}

			// Figure (images with captions to <figure>)
			if enabled, ok := extensions["figure"].(bool); ok {
				config.FigureEnabled = enabled
			}

			// Anchor (heading permalinks)
			if enabled, ok := extensions["anchor"].(bool); ok {
				config.AnchorEnabled = enabled
			}
		}
	}

	return config
}

// Render converts markdown content to HTML for all posts.
// Posts with Skip=true are skipped.
// The rendered HTML is stored in post.ArticleHTML.
// Uses build cache to skip re-rendering unchanged content.
//
// Uses two-phase processing for incremental optimization:
// Phase 1: Quick single-threaded pass to restore cached HTML (no worker overhead)
// Phase 2: Concurrent processing only for posts that need rendering
func (p *RenderMarkdownPlugin) Render(m *lifecycle.Manager) error {
	// Get build cache for HTML caching
	p.cache = GetBuildCache(m)
	canonicalPaths := make(map[string]bool)
	m.Cache().Set(canonicalArticlePathsKey, canonicalPaths)

	// Phase 1: Pre-filter posts and restore cached HTML
	// This avoids worker pool overhead for the ~98% of posts that are cached
	postsNeedingRender := make([]*models.Post, 0, len(m.Posts()))
	var migrationPosts []*models.Post
	for _, post := range m.Posts() {
		// Skip posts marked as skip
		if post.Skip {
			continue
		}

		// Skip posts with no content
		if post.Content == "" {
			post.ArticleHTML = ""
			canonicalPaths[post.Path] = true
			if post.HTML != "" {
				m.ContentLedger().MarkRendered(post.Path)
			}
			continue
		}

		// Try to get cached HTML if content hasn't changed
		if p.cache != nil && !isSourceEncryptedPost(post) {
			contentHash := buildcache.ContentHash(post.Content)
			if cachedHTML := p.cache.GetCachedArticleHTML(post.Path, contentHash); cachedHTML != "" {
				post.ArticleHTML = cachedHTML
				canonicalPaths[post.Path] = true
				// Detect CSS requirements from cached HTML
				p.detectCSSRequirements(post)
				m.ContentLedger().MarkRendered(post.Path)
				if p.needsHeadingHighlightMigration(post) {
					migrationPosts = append(migrationPosts, post)
				}
				continue // Already handled, no concurrent processing needed
			}
		}

		migrationNeeded, err := p.classifyIncrementalHeadingMigration(m, post)
		if err != nil {
			return err
		}
		if migrationNeeded {
			migrationPosts = append(migrationPosts, post)
		}
		postsNeedingRender = append(postsNeedingRender, post)
	}
	p.markHeadingHighlightMigration(m, migrationPosts)

	if lifecycle.IsServeIncremental(m) {
		if affected := lifecycle.GetServeAffectedPaths(m); len(affected) > 0 {
			filtered := postsNeedingRender[:0]
			for _, post := range postsNeedingRender {
				if affected[post.Path] {
					filtered = append(filtered, post)
				}
			}
			postsNeedingRender = filtered
		}
	}

	// Phase 2: Process only posts that need rendering concurrently
	// Each worker owns its post; map mutation and dependency expansion happen
	// only after the workers finish.
	err := m.ProcessPostsSliceConcurrently(postsNeedingRender, func(post *models.Post) error {
		if err := p.renderPost(post); err != nil {
			m.ContentLedger().RecordError(post.Path, diagnostics.ReasonContentRenderError, "Markdown rendering failed")
			return err
		}
		if post.ArticleHTML != "" {
			m.ContentLedger().MarkRendered(post.Path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	migrationPosts = migrationPosts[:0]
	for _, post := range postsNeedingRender {
		canonicalPaths[post.Path] = true
		if p.needsHeadingHighlightMigration(post) {
			migrationPosts = append(migrationPosts, post)
		}
	}
	p.markHeadingHighlightMigration(m, migrationPosts)
	return nil
}

// classifyIncrementalHeadingMigration probes eligible legacy cache misses
// without caching unrelated omitted articles or probing encrypted sources.
func (p *RenderMarkdownPlugin) classifyIncrementalHeadingMigration(m *lifecycle.Manager, post *models.Post) (bool, error) {
	if !lifecycle.IsServeIncremental(m) || !p.headingHighlightMigrationEligible(post) {
		return false, nil
	}
	article, err := p.renderArticleHTML(post.Content)
	if err != nil {
		m.ContentLedger().RecordError(post.Path, diagnostics.ReasonContentRenderError, "Markdown migration classification failed")
		return false, err
	}
	return hasHeadingMarkHighlights(article), nil
}

func (p *RenderMarkdownPlugin) headingHighlightMigrationEligible(post *models.Post) bool {
	return p.cache != nil && !post.Skip && post.Path != "" && !isSourceEncryptedPost(post) &&
		p.cache.GetHeadingHighlightRevision(post.Path) < headingHighlightRevision
}

func (p *RenderMarkdownPlugin) needsHeadingHighlightMigration(post *models.Post) bool {
	return p.headingHighlightMigrationEligible(post) && hasHeadingMarkHighlights(post.ArticleHTML)
}

func (p *RenderMarkdownPlugin) markHeadingHighlightMigration(m *lifecycle.Manager, posts []*models.Post) {
	if len(posts) == 0 {
		return
	}
	affected := lifecycle.GetServeAffectedPaths(m)
	if affected == nil {
		affected = make(map[string]bool)
	}
	for _, post := range posts {
		affected[post.Path] = true
		for _, identity := range postDependencyIdentities(post) {
			p.cache.MarkSlugChanged(identity)
		}
		p.cache.MarkFeedSlugChanged(post.Slug)
	}
	changed := p.cache.GetChangedSlugs()
	for _, path := range p.cache.GetAffectedPosts(changed) {
		affected[path] = true
	}
	p.cache.MarkAffectedDependents(changed)
	lifecycle.SetServeAffectedPaths(m, affected)
}

// renderPost renders a single post's markdown content to HTML.
// Uses a buffer pool to reduce allocations and GC pressure.
// Note: Cache lookup is now done in the filter phase, so posts reaching here need rendering.
func (p *RenderMarkdownPlugin) renderPost(post *models.Post) error {
	// Skip posts marked as skip (defensive check, filter should catch this)
	if post.Skip {
		return nil
	}

	// Skip posts with no content (defensive check, filter should catch this)
	if post.Content == "" {
		post.ArticleHTML = ""
		return nil
	}

	// Render the markdown
	renderedHTML, err := p.renderArticleHTML(post.Content)
	if err != nil {
		return err
	}
	post.ArticleHTML = renderedHTML

	// Cache the result for future incremental builds
	if p.cache != nil && !isSourceEncryptedPost(post) {
		contentHash := buildcache.ContentHash(post.Content)
		//nolint:errcheck // caching is best-effort, failures are non-fatal
		p.cache.CacheArticleHTML(post.Path, contentHash, renderedHTML)
	}

	// Detect CSS requirements from rendered HTML
	p.detectCSSRequirements(post)
	return nil
}

// renderArticleHTML produces the canonical article-cache boundary without
// mutating a post or cache. Migration probes and fresh renders use the same
// Markdown conversion and post-processing pipeline.
func (p *RenderMarkdownPlugin) renderArticleHTML(content string) (string, error) {
	renderedHTML, err := p.doRender(content)
	if err != nil {
		return "", err
	}
	renderedHTML = mergeFigureBlockquoteCaptions(renderedHTML)
	renderedHTML = mergeBlockquoteAttributions(renderedHTML)
	return wrapHeadingMarkHighlights(renderedHTML), nil
}

// doRender performs the actual markdown to HTML conversion.
func (p *RenderMarkdownPlugin) doRender(content string) (string, error) {
	// Pre-process: prepare blockquote attributions before goldmark rendering
	content = prepareMarkdownAttribution(content)

	// Get buffer from pool
	buf, ok := markdownBufferPool.Get().(*bytes.Buffer)
	if !ok {
		buf = new(bytes.Buffer)
	}
	buf.Reset()
	defer markdownBufferPool.Put(buf)

	// Convert markdown to HTML
	if err := p.md.Convert([]byte(content), buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func mergeFigureBlockquoteCaptions(renderedHTML string) string {
	var b strings.Builder
	for pos := 0; pos < len(renderedHTML); {
		figureStart := strings.Index(renderedHTML[pos:], "<figure")
		if figureStart < 0 {
			b.WriteString(renderedHTML[pos:])
			break
		}
		figureStart += pos
		b.WriteString(renderedHTML[pos:figureStart])

		figureEnd := strings.Index(renderedHTML[figureStart:], "</figure>")
		if figureEnd < 0 {
			b.WriteString(renderedHTML[figureStart:])
			break
		}
		figureEnd += figureStart + len("</figure>")

		figureHTML := renderedHTML[figureStart:figureEnd]
		if strings.Contains(figureHTML, "<figcaption") || insideWebAwesomeComparisonContainer(renderedHTML, figureStart) {
			b.WriteString(figureHTML)
			pos = figureEnd
			continue
		}

		afterFigure := figureEnd
		for afterFigure < len(renderedHTML) {
			c := renderedHTML[afterFigure]
			if c != ' ' && c != '\n' && c != '\t' && c != '\r' {
				break
			}
			afterFigure++
		}

		if !strings.HasPrefix(renderedHTML[afterFigure:], "<blockquote>") {
			b.WriteString(figureHTML)
			pos = figureEnd
			continue
		}

		blockquoteEnd := strings.Index(renderedHTML[afterFigure:], "</blockquote>")
		if blockquoteEnd < 0 {
			b.WriteString(figureHTML)
			pos = figureEnd
			continue
		}
		blockquoteEnd += afterFigure + len("</blockquote>")

		blockquoteInner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(renderedHTML[afterFigure:blockquoteEnd], "<blockquote>"), "</blockquote>"))
		if !figureBlockquoteCaptionParagraphRegex.MatchString(blockquoteInner) {
			b.WriteString(figureHTML)
			pos = figureEnd
			continue
		}

		b.WriteString(strings.Replace(figureHTML, "</figure>", "<figcaption>"+blockquoteInner+"</figcaption></figure>", 1))
		pos = blockquoteEnd
	}
	return b.String()
}

func insideWebAwesomeComparisonContainer(renderedHTML string, pos int) bool {
	openIndex := strings.LastIndex(renderedHTML[:pos], "<div")
	closeIndex := strings.LastIndex(renderedHTML[:pos], "</div>")
	if openIndex < 0 || closeIndex > openIndex {
		return false
	}

	openTagEnd := strings.Index(renderedHTML[openIndex:], ">")
	if openTagEnd < 0 {
		return false
	}

	openTag := renderedHTML[openIndex : openIndex+openTagEnd+1]
	return strings.Contains(openTag, `class="wa-comparison"`) ||
		(strings.Contains(openTag, `class="`) && strings.Contains(openTag, "webawesome") && strings.Contains(openTag, "comparison"))
}

// prepareMarkdownAttribution detects blockquote lines followed by attribution
// lines (no blank line between) and inserts a blank line + marker to prevent
// goldmark's lazy continuation from merging them into the blockquote.
func prepareMarkdownAttribution(markdown string) string {
	lines := strings.Split(markdown, "\n")
	result := make([]string, 0, len(lines))

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			result = append(result, line)
			continue
		}

		if strings.HasPrefix(strings.TrimLeft(line, " \t"), ">") {
			result = append(result, line)
			continue
		}

		// Non-blockquote line; check if it follows a blockquote with no blank line
		if isImmediatelyAfterBlockquote(lines, i) && isAttributionText(trimmed) {
			result = append(result, "", attributionMarker, line)
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// isImmediatelyAfterBlockquote checks whether lines[i] is directly after a
// blockquote line without any intervening blank line.
func isImmediatelyAfterBlockquote(lines []string, i int) bool {
	for j := i - 1; j >= 0; j-- {
		prev := strings.TrimSpace(lines[j])
		if prev == "" {
			return false
		}
		if strings.HasPrefix(strings.TrimLeft(lines[j], " \t"), ">") {
			return true
		}
	}
	return false
}

// isAttributionText checks whether text looks like a blockquote attribution.
// This runs after the Transform stage, so mentions/wikilinks may already be
// replaced with <a> tags in the content.
func isAttributionText(text string) bool {
	if text == "" {
		return false
	}

	// Already-processed <a> tag (from mentions or wikilinks transforms)
	if strings.HasPrefix(text, "<a ") && strings.Contains(text, "</a>") {
		return true
	}

	// @mention — @username (pre-transform state)
	if strings.HasPrefix(text, "@") && len(text) > 1 {
		return true
	}

	// Wikilink — [[Page Name]] (pre-transform state)
	if strings.Contains(text, "[[") && strings.Contains(text, "]]") {
		return true
	}

	// Markdown link — [text](url) at start or preceded by plain text
	if strings.Contains(text, "](") {
		return true
	}

	// Dash prefix — attribution or — attribution
	if strings.HasPrefix(text, "—") || strings.HasPrefix(text, "--") {
		return true
	}

	// Short plain text (≤ 60 chars) — likely a name or short attribution
	if len(text) <= 60 && !strings.HasPrefix(text, "http") {
		// Ensure it's mostly unicode letters and spacing to avoid matching code
		letterCount := 0
		for _, r := range text {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) || strings.ContainsRune(".,'-@_", r) {
				letterCount++
			}
		}
		return letterCount*2 >= len([]rune(text)) // at least half the chars are "word-like"
	}

	return false
}

// attributionFooter holds the rendered footer HTML and optional cite URL.
type attributionFooter struct {
	footerHTML string
	sourceURL  string
}

// attributionLinkInfo holds parsed information about a link in attribution HTML.
type attributionLinkInfo struct {
	fullTag    string // the complete <a...>...</a> tag
	href       string // the href URL
	text       string // the inner text content
	isMention  bool   // has mention class (from mentions plugin)
	isWikilink bool   // has wikilink class (from wikilinks plugin)
}

// mergeBlockquoteAttributions merges blockquotes followed by attribution
// markers and paragraphs, inserting the attribution as a <footer> inside the
// <blockquote>. Source links are wrapped in <cite> and blockquote[cite] is
// set from the source URL.
func mergeBlockquoteAttributions(renderedHTML string) string {
	marker := attributionMarker
	result := renderedHTML

	for {
		markerPos := strings.Index(result, marker)
		if markerPos < 0 {
			break
		}

		bqClose := strings.LastIndex(result[:markerPos], "</blockquote>")
		if bqClose < 0 {
			result = strings.Replace(result, marker, "", 1)
			continue
		}

		bqOpen := strings.LastIndex(result[:bqClose], "<blockquote>")
		if bqOpen < 0 {
			result = strings.Replace(result, marker, "", 1)
			continue
		}
		bqEnd := bqClose + len("</blockquote>")

		between := strings.TrimSpace(result[bqEnd:markerPos])
		if between != "" {
			result = strings.Replace(result, marker, "", 1)
			continue
		}

		afterMarker := markerPos + len(marker)
		afterMarker = skipHTMLWhitespace(result, afterMarker)

		if !strings.HasPrefix(result[afterMarker:], "<p>") {
			result = strings.Replace(result, marker, "", 1)
			continue
		}

		pTagEnd := strings.Index(result[afterMarker:], "</p>")
		if pTagEnd < 0 {
			result = strings.Replace(result, marker, "", 1)
			continue
		}
		pTagEnd += afterMarker + len("</p>")

		// Extract attribution content (between <p> and </p>)
		attrContent := result[afterMarker+3 : pTagEnd-4]

		bqHTML := result[bqOpen:bqEnd]
		footer := buildAttributionFooter(attrContent)

		// Add cite attribute to blockquote if source URL found
		if footer.sourceURL != "" {
			escapedURL := strings.ReplaceAll(footer.sourceURL, `"`, "&quot;")
			bqHTML = strings.Replace(bqHTML, "<blockquote>", `<blockquote cite="`+escapedURL+`">`, 1)
		}

		replacement := strings.Replace(bqHTML, "</blockquote>", "  "+footer.footerHTML+"\n</blockquote>", 1)
		result = result[:bqOpen] + replacement + result[pTagEnd:]
	}

	return result
}

func skipHTMLWhitespace(s string, pos int) int {
	for pos < len(s) {
		c := s[pos]
		if c != ' ' && c != '\n' && c != '\t' && c != '\r' {
			break
		}
		pos++
	}
	return pos
}

// findAttributionLinks extracts all <a> tags from attribution paragraph HTML.
func findAttributionLinks(html string) []attributionLinkInfo {
	var links []attributionLinkInfo
	pos := 0
	for pos < len(html) {
		aStart := strings.Index(html[pos:], "<a ")
		if aStart < 0 {
			break
		}
		aStart += pos
		aEnd := strings.Index(html[aStart:], "</a>")
		if aEnd < 0 {
			break
		}
		aEnd += aStart + len("</a>")

		fullTag := html[aStart:aEnd]
		href := extractAttrValue(fullTag, "href")
		text := extractLinkText(fullTag)
		isMention := strings.Contains(fullTag, "mention")
		isWiki := strings.Contains(fullTag, "wikilink")

		links = append(links, attributionLinkInfo{
			fullTag:    fullTag,
			href:       href,
			text:       text,
			isMention:  isMention,
			isWikilink: isWiki,
		})
		pos = aEnd
	}
	return links
}

// extractAttrValue extracts the value of an HTML attribute from a tag string.
func extractAttrValue(tag, attr string) string {
	pattern := " " + attr + "=\""
	idx := strings.Index(tag, pattern)
	if idx < 0 {
		return ""
	}
	start := idx + len(pattern)
	end := strings.Index(tag[start:], "\"")
	if end < 0 {
		return ""
	}
	return tag[start : start+end]
}

// extractLinkText extracts the visible text content from an <a> tag.
func extractLinkText(tag string) string {
	// Find the first > that closes the opening <a> tag
	openBracket := strings.Index(tag, ">")
	if openBracket < 0 {
		return ""
	}
	// tag[openBracket:] starts after opening > (e.g. ">text</a>...")
	// Find </a> relative to that position
	closeTag := strings.Index(tag[openBracket:], "</a>")
	if closeTag < 0 {
		return ""
	}
	// Text is between the opening > and the start of </a>
	return tag[openBracket+1 : openBracket+closeTag]
}

// isSourceLinkText checks if link text indicates a source (not person) link.
func isSourceLinkText(text string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	return lower == "via" || lower == "source" || lower == "read more" || lower == gardenEdgeTypeLink || lower == "permalink"
}

// textBeforeTag returns content before the first occurrence of tag in html.
func textBeforeTag(html, tag string) string {
	idx := strings.Index(html, tag)
	if idx < 0 {
		return html
	}
	return html[:idx]
}

// buildAttributionFooter parses attribution paragraph HTML and returns the
// footer markup plus an optional source URL for blockquote[cite].
//
// Classification rules:
//   - 0 links: all content is the person name
//   - 1 link that is a mention: person
//   - 1 link with source text (via, source, etc.): source only
//   - 1 bare-URL autolink: source only
//   - 1 link with text before it: text=person, link=source
//   - 1 link standalone (person name in markdown link): person
//   - 2+ links: first=person, last=source, interleaving text preserved
func buildAttributionFooter(attrContent string) attributionFooter {
	links := findAttributionLinks(attrContent)

	var personHTML string
	var sourceURL string
	var sourceLink string

	switch {
	case len(links) == 0:
		personHTML = attrContent

	case len(links) == 1:
		l := links[0]
		switch {
		case l.isMention:
			personHTML = l.fullTag
		case isSourceLinkText(l.text):
			sourceURL = l.href
			sourceLink = l.fullTag
		case l.href == l.text || strings.TrimPrefix(l.href, "https://") == strings.TrimSpace(l.text) ||
			strings.TrimPrefix(l.href, "http://") == strings.TrimSpace(l.text):
			sourceURL = l.href
			sourceLink = l.fullTag
		default:
			preText := strings.TrimSpace(textBeforeTag(attrContent, l.fullTag))
			if preText != "" {
				personHTML = preText
				sourceURL = l.href
				sourceLink = l.fullTag
			} else {
				personHTML = l.fullTag
			}
		}

	default:
		// Multiple links: everything before the last link is person
		// (preserving interleaving text), last link is source.
		last := links[len(links)-1]
		lastPos := strings.Index(attrContent, last.fullTag)
		if lastPos >= 0 {
			personHTML = attrContent[:lastPos]
		} else {
			personHTML = attrContent
		}
		sourceURL = last.href
		sourceLink = last.fullTag
	}

	var footer strings.Builder
	footer.WriteString("<footer>")
	// Only prepend em dash if the attribution doesn't already start with one
	if !strings.HasPrefix(strings.TrimSpace(personHTML), "\u2014") &&
		!strings.HasPrefix(strings.TrimSpace(personHTML), "&mdash;") {
		footer.WriteString("\u2014 ")
	}
	footer.WriteString(strings.TrimRight(personHTML, " "))
	if sourceLink != "" {
		footer.WriteString(" <cite>")
		footer.WriteString(sourceLink)
		footer.WriteString("</cite>")
	}
	footer.WriteString("</footer>")

	return attributionFooter{
		footerHTML: footer.String(),
		sourceURL:  sourceURL,
	}
}

// detectCSSRequirements scans the rendered HTML and sets flags in post.Extra
// to indicate which CSS files are needed. This follows the ImageZoom pattern
// for conditional asset loading.
func (p *RenderMarkdownPlugin) detectCSSRequirements(post *models.Post) {
	if post.ArticleHTML == "" {
		return
	}

	// Initialize Extra map if needed
	if post.Extra == nil {
		post.Extra = make(map[string]interface{})
	}

	// Detect admonitions - check for admonition class
	if strings.Contains(post.ArticleHTML, `class="admonition`) {
		post.Extra["needs_admonitions_css"] = true
	}

	// Detect code blocks - check for chroma syntax highlighting classes or code/pre tags
	if strings.Contains(post.ArticleHTML, `class="chroma"`) ||
		strings.Contains(post.ArticleHTML, `class="highlight"`) ||
		strings.Contains(post.ArticleHTML, "<pre><code") ||
		strings.Contains(post.ArticleHTML, `<code class="language-`) ||
		strings.Contains(post.ArticleHTML, `data-needs-code-css="true"`) {
		post.Extra["needs_code_css"] = true
	}

	// Detect images that will get GLightbox treatment (image_zoom plugin
	// runs after this, but we can detect existing glightbox markers or
	// images that will be processed).
	if strings.Contains(post.ArticleHTML, `class="glightbox"`) ||
		strings.Contains(post.ArticleHTML, `data-glightbox`) ||
		strings.Contains(post.ArticleHTML, `{data-zoomable}`) ||
		strings.Contains(post.ArticleHTML, `{.zoomable}`) {
		post.Extra["needs_image_zoom"] = true
	}

	// Detect Mermaid diagrams - check for mermaid class or pre tags
	if strings.Contains(post.ArticleHTML, `class="mermaid"`) ||
		strings.Contains(post.ArticleHTML, `<pre class="mermaid"`) ||
		strings.Contains(post.ArticleHTML, `<div class="mermaid"`) {
		post.Extra["has_mermaid"] = true
	}

	// Detect figure elements - check for figure tag (goldmark-figure extension)
	if strings.Contains(post.ArticleHTML, "<figure") {
		post.Extra["has_figure"] = true
	}

	// Detect anchor links - check for anchor class (goldmark-anchor extension)
	if strings.Contains(post.ArticleHTML, `class="anchor"`) {
		post.Extra["has_anchor_links"] = true
	}

	// Detect roads - `::: road STYLE` containers or elements with data-road.
	if strings.Contains(post.ArticleHTML, "road") && roadMarkupPattern.MatchString(post.ArticleHTML) {
		post.Extra["needs_roads"] = true
	}
}

// roadMarkupPattern matches a "road" class token or a data-road attribute.
var roadMarkupPattern = regexp.MustCompile(`class="(?:[^"]*\s)?road(?:\s[^"]*)?"|\sdata-road[\s=>]`)

// Ensure RenderMarkdownPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*RenderMarkdownPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*RenderMarkdownPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*RenderMarkdownPlugin)(nil)
)
