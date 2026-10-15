// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata-go/pkg/buildcache"
	"github.com/WaylonWalker/markata-go/pkg/diagnostics"
	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/logging"
	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/palettes"
	"github.com/WaylonWalker/markata-go/pkg/templates"
)

var publishHTMLLog = logging.Component("publish_html").Phase("write")

// defaultTxtTemplate is the default template name for txt output.
const defaultTxtTemplate = "default.txt"

// rawTxtTemplate is the template name for raw txt output.
const rawTxtTemplate = "raw.txt"

// defaultAnsiTemplate is the default template name for ANSI output.
const defaultAnsiTemplate = "default.ansi"

// rawAnsiTemplate is the template name for raw ANSI output.
const rawAnsiTemplate = "raw.ansi"

// rawMdTemplate is the template name for raw markdown output.
const rawMdTemplate = "raw.md"

// specialFiles are slugs that should have their content at /slug.ext rather than /slug/index.ext.
// These are standard web files that are expected at specific root-level locations.
var specialFiles = []string{"robots", "llms", "humans", "security", "ads"}

// isSpecialFile returns true if the slug is a special file that should be served at root level.
func isSpecialFile(slug string) bool {
	for _, special := range specialFiles {
		if slug == special {
			return true
		}
	}
	return false
}

func getPostExtraString(post *models.Post, keys ...string) string {
	if post == nil || post.Extra == nil {
		return ""
	}

	for _, key := range keys {
		value, ok := post.Extra[key]
		if !ok {
			continue
		}

		strValue, ok := value.(string)
		if !ok {
			continue
		}

		trimmed := strings.TrimSpace(strValue)
		if trimmed != "" {
			return trimmed
		}
	}

	return ""
}

// PublishHTMLPlugin writes individual post HTML files during the write stage.
// It supports multiple output formats: HTML, Markdown source, and OG card HTML.
type PublishHTMLPlugin struct {
	templatesPlugin *TemplatesPlugin
}

// NewPublishHTMLPlugin creates a new PublishHTMLPlugin.
func NewPublishHTMLPlugin() *PublishHTMLPlugin {
	return &PublishHTMLPlugin{}
}

// Name returns the unique name of the plugin.
func (p *PublishHTMLPlugin) Name() string {
	return "publish_html"
}

// Write outputs each post to the configured output directory in enabled formats.
func (p *PublishHTMLPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Get template engine from cache (may be nil if templates plugin not used)
	var engine *templates.Engine
	if cached, ok := m.Cache().Get("templates.engine"); ok && cached != nil {
		if e, ok := cached.(*templates.Engine); ok {
			engine = e
		}
	}

	// Get templates plugin from cache for per-format template resolution
	if cached, ok := m.Cache().Get("templates.plugin"); ok {
		if tp, ok := cached.(*TemplatesPlugin); ok {
			p.templatesPlugin = tp
		}
	}

	// Pre-pass: Identify all posts that changed (hash mismatch) and mark their slugs.
	// This enables dependency-based invalidation in FilterPosts below.
	// Uses ShouldRebuildBatch to acquire the lock once instead of per-post.
	cache := GetBuildCache(m)
	affected := lifecycle.GetServeAffectedPaths(m)
	if err := p.removeStalePostOutputs(m.Posts(), config, cache); err != nil {
		return err
	}
	p.recordDraftInputs(m.Posts(), cache)
	if err := p.markChangedPosts(cache, m, config); err != nil {
		return err
	}

	// Single-pass filter: identify posts that need writing.
	// This replaces the previous two-pass approach (pre-filter + writeHTMLFormat check).
	postsNeedingWrite := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.Draft {
			return false
		}
		if affected[post.Path] {
			return true
		}

		// For incremental builds, check if we can skip this post
		if cache != nil && post.InputHash != "" {
			if !cache.ShouldRebuildWithSlug(post.Path, post.Slug, post.InputHash, post.Template) {
				if p.expectedPostOutputsExist(post, config) && !projectStaticShadowsPost(post, config) {
					cache.MarkSkipped()
					m.ContentLedger().MarkEmitted(post.Path)
					return false
				}
			}
		}

		return true
	})

	// Phase 2: Process only posts that need writing concurrently
	return m.ProcessPostsSliceConcurrently(postsNeedingWrite, func(post *models.Post) error {
		if err := p.writePost(post, config, engine, m); err != nil {
			m.ContentLedger().RecordError(post.Path, diagnostics.ReasonContentWriteError, "post output could not be written")
			return err
		}
		if p.postOutputExists(post, config) {
			m.ContentLedger().MarkEmitted(post.Path)
		}
		return nil
	})
}

// projectStaticShadowsPost prevents a warm build from letting the static
// copier replace an authored page that won the same path on a cold build.
// Ordinary cached posts remain skippable; only a conflicting page is rewritten.
func projectStaticShadowsPost(post *models.Post, config *lifecycle.Config) bool {
	if post == nil || config == nil {
		return false
	}
	staticDir := NewStaticAssetsPlugin().projectStaticDir(config)
	name := filepath.Join(staticDir, filepath.FromSlash(strings.Trim(post.Slug, "/")), "index.html")
	info, err := os.Stat(name)
	return err == nil && info.Mode().IsRegular()
}

func (p *PublishHTMLPlugin) recordDraftInputs(posts []*models.Post, cache *buildcache.Cache) {
	if cache == nil {
		return
	}
	for _, post := range posts {
		if post != nil && post.Draft && !post.Skip && post.Path != "" && post.InputHash != "" {
			cache.MarkInputProcessed(post.Path, post.Slug, post.InputHash, post.Template)
		}
	}
}

func (p *PublishHTMLPlugin) expectedPostOutputsExist(post *models.Post, config *lifecycle.Config) bool {
	formats := resolvePostFormats(post, config)
	expected := make([]string, 0, 5)
	if formats.IsHTMLEnabled() && (post.HTML != "" || post.ArticleHTML != "" || post.Content != "") {
		expected = append(expected, filepath.Join(config.OutputDir, post.Slug, "index.html"))
	}
	if !post.Private {
		hasReversedFormat := formats.Markdown || formats.Text || formats.ANSI
		if !formats.IsHTMLEnabled() && hasReversedFormat {
			expected = append(expected, filepath.Join(config.OutputDir, post.Slug, "index.html"))
		}
		if formats.Markdown {
			expected = append(expected,
				filepath.Join(config.OutputDir, post.Slug+".md"),
				filepath.Join(config.OutputDir, post.Slug, "index.md", "index.html"),
			)
		}
		if formats.Text {
			expected = append(expected,
				filepath.Join(config.OutputDir, post.Slug+".txt"),
				filepath.Join(config.OutputDir, post.Slug, "index.txt", "index.html"),
			)
		}
		if formats.ANSI {
			expected = append(expected,
				filepath.Join(config.OutputDir, post.Slug+".ansi"),
				filepath.Join(config.OutputDir, post.Slug, "index.ansi", "index.html"),
			)
		}
		if formats.OG {
			expected = append(expected, filepath.Join(config.OutputDir, post.Slug, "og", "index.html"))
		}
	}
	if len(expected) == 0 {
		return false
	}
	for _, path := range expected {
		if _, err := os.Stat(path); err != nil {
			return false
		}
	}
	return true
}

func (p *PublishHTMLPlugin) postOutputExists(post *models.Post, config *lifecycle.Config) bool {
	if post == nil || config == nil || post.Skip || post.Draft {
		return false
	}
	formats := resolvePostFormats(post, config)
	paths := make([]string, 0, 5)
	if formats.IsHTMLEnabled() && (post.HTML != "" || post.ArticleHTML != "") {
		paths = append(paths, filepath.Join(config.OutputDir, post.Slug, "index.html"))
	}
	if !post.Private {
		if formats.Markdown {
			paths = append(paths, filepath.Join(config.OutputDir, post.Slug+".md"))
		}
		if formats.Text {
			paths = append(paths, filepath.Join(config.OutputDir, post.Slug+".txt"))
		}
		if formats.ANSI {
			paths = append(paths, filepath.Join(config.OutputDir, post.Slug+".ansi"))
		}
		if formats.OG {
			paths = append(paths, filepath.Join(config.OutputDir, post.Slug, "og", "index.html"))
		}
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}
	return false
}

func (p *PublishHTMLPlugin) markChangedPosts(cache *buildcache.Cache, m *lifecycle.Manager, config *lifecycle.Config) error {
	if cache == nil || m == nil || config == nil {
		return nil
	}

	posts := m.Posts()
	removed := lifecycle.GetServeRemovedPaths(m)
	postFormats := getPostFormatsConfig(config)

	batch := make([]struct{ Path, InputHash, Template string }, 0, len(posts))
	slugByPath := make(map[string]string, len(posts))
	for _, post := range posts {
		if post.Skip || post.Draft || post.InputHash == "" {
			continue
		}
		batch = append(batch, struct{ Path, InputHash, Template string }{
			Path:      post.Path,
			InputHash: post.InputHash,
			Template:  post.Template,
		})
		slugByPath[post.Path] = post.Slug
	}
	changed := cache.ShouldRebuildBatch(batch)
	for path := range changed {
		slug := slugByPath[path]
		// An explicit slug rename invalidates links and backlinks pointing at
		// both the old and the new route. Invalidate before updating the graph.
		if oldSlug := cache.Graph.SlugForPath(path); oldSlug != "" && oldSlug != slug {
			cache.MarkSlugChanged(oldSlug)
		}
		cache.MarkSlugChanged(slug)
	}
	for _, removedPath := range removed {
		if slug := cache.Graph.SlugForPath(removedPath); slug != "" {
			cache.MarkSlugChanged(slug)
		}
		if err := p.removePostOutputs(removedPath, config, postFormats, cache); err != nil {
			return err
		}
	}

	return nil
}

// writePost writes a single post to its output location in all enabled formats.
// Shadow pages: Unpublished posts are still rendered but not included in feeds.
// This allows sharing draft content via direct URL while keeping it out of public listings.
//
//nolint:gocyclo // Format-specific output and privacy rules are kept in one transaction.
func (p *PublishHTMLPlugin) writePost(post *models.Post, config *lifecycle.Config, engine *templates.Engine, m *lifecycle.Manager) error {
	// Determine output path before early returns so stale output can be removed
	// for private drafts and skipped posts as well.
	if !post.Has("_slug_explicit") && post.Slug == "" {
		post.GenerateSlugWithMode(configuredSlugMode(config, post.Path))
	}

	// Skip posts marked as skip
	if post.Skip {
		if post.Private && post.Slug != "" {
			return removeAllPostOutputs(config.OutputDir, post.Slug)
		}
		return nil
	}

	// Skip drafts - these are truly private work-in-progress content
	if post.Draft {
		if post.Private && post.Slug != "" {
			return removeAllPostOutputs(config.OutputDir, post.Slug)
		}
		return nil
	}

	// Note: Unpublished posts (published: false) are still rendered as "shadow pages"
	// They won't appear in feeds (which filter by published == True) but can be
	// accessed via direct URL for review/sharing purposes.

	// Determine output path
	// Use slug to create: output_dir/slug/index.html
	outputDir := config.OutputDir
	postDir, err := safePostOutputPath(outputDir, post.Slug)
	if err != nil {
		return err
	}

	// Create post directory
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		return fmt.Errorf("creating post directory %s: %w", postDir, err)
	}

	// Resolve site defaults plus any per-post frontmatter override.
	postFormats := resolvePostFormats(post, config)

	// Get build cache for incremental builds
	cache := GetBuildCache(m)
	removeDisabledPostOutputs(config.OutputDir, post.Slug, postFormats)
	if post.Private {
		if err := removePrivatePostOutputs(config.OutputDir, post.Slug); err != nil {
			return err
		}
	}

	// Write HTML format (default)
	htmlOutputPath := ""
	if postFormats.IsHTMLEnabled() {
		htmlOutputPath, err = p.writeHTMLFormat(post, config, postDir)
		if err != nil {
			return err
		}
	}

	// Write Markdown format (raw source)
	// Uses reversed redirect: content at /slug.md, redirect at /slug/index.html
	// Skip redirect if HTML is enabled (index.html already has main content)
	// Skip for private posts to prevent plaintext content leaks
	if postFormats.Markdown && !post.Private {
		mdContent := p.buildFormatContent(post, config, m, "markdown")
		skipRedirect := postFormats.IsHTMLEnabled()
		if err := p.writeReversedFormatOutput(post.Slug, "md", mdContent, config.OutputDir, skipRedirect); err != nil {
			return err
		}
	}

	// Write Text format (plain text)
	// Uses reversed redirect: content at /slug.txt, redirect at /slug/index.html
	// Skip redirect if HTML is enabled (index.html already has main content)
	// Skip for private posts to prevent plaintext content leaks
	if postFormats.Text && !post.Private {
		// Use renderTextContent for txt format to leverage main's sophisticated template resolution
		txtContent := p.renderTextContent(post, config, engine)
		skipRedirect := postFormats.IsHTMLEnabled()
		if err := p.writeReversedFormatOutput(post.Slug, "txt", txtContent, config.OutputDir, skipRedirect); err != nil {
			return err
		}
	}

	// Write ANSI terminal format
	if postFormats.ANSI && !post.Private {
		ansiContent := p.renderANSIContent(post, config, engine)
		skipRedirect := postFormats.IsHTMLEnabled()
		if err := p.writeReversedFormatOutput(post.Slug, "ansi", ansiContent, config.OutputDir, skipRedirect); err != nil {
			return err
		}
	}

	// Write OG format (social card HTML)
	// Skip for private posts to prevent metadata leaks in OG cards
	if postFormats.OG && !post.Private {
		if err := p.writeOGFormat(post, config, postDir, engine); err != nil {
			return err
		}
	}

	if cache != nil && post.InputHash != "" {
		if htmlOutputPath != "" {
			cache.MarkRebuiltWithSlug(post.Path, post.Slug, post.InputHash, htmlOutputPath, post.Template)
		} else if p.canRecordProcessedInput(post, postFormats, postDir) {
			cache.MarkInputProcessed(post.Path, post.Slug, post.InputHash, post.Template)
		}
	}

	return nil
}

func (p *PublishHTMLPlugin) canRecordProcessedInput(post *models.Post, formats models.PostFormatsConfig, postDir string) bool {
	if !formats.IsHTMLEnabled() || post.HTML != "" || post.ArticleHTML != "" {
		return true
	}

	_, err := os.Stat(filepath.Join(postDir, "index.html"))
	return os.IsNotExist(err)
}

func (p *PublishHTMLPlugin) removeStalePostOutputs(posts []*models.Post, config *lifecycle.Config, cache *buildcache.Cache) error {
	if config == nil {
		return nil
	}

	for _, post := range posts {
		if post == nil {
			continue
		}
		if !post.Has("_slug_explicit") && post.Slug == "" {
			post.GenerateSlugWithMode(configuredSlugMode(config, post.Path))
		}

		if post.Slug != "" {
			if post.Path != "" && (post.Skip || post.Draft) {
				if err := removeAllPostOutputs(config.OutputDir, post.Slug); err != nil {
					return err
				}
			} else if post.Private {
				if err := removePrivatePostOutputs(config.OutputDir, post.Slug); err != nil {
					return err
				}
			}
		}

		if cache == nil {
			continue
		}
		cached := cache.GetCachedPost(post.Path)
		if cached == nil {
			continue
		}
		oldSlug := cachedOutputSlug(config.OutputDir, cached.OutputPath)
		if oldSlug == "" {
			oldSlug = cached.Slug
		}
		if oldSlug != "" && oldSlug != post.Slug {
			if err := removeAllPostOutputs(config.OutputDir, oldSlug); err != nil {
				return err
			}
		}
	}

	return nil
}

func cachedOutputSlug(outputDir, outputPath string) string {
	if outputDir == "" || outputPath == "" {
		return ""
	}
	root, err := filepath.Abs(outputDir)
	if err != nil {
		return ""
	}
	target, err := filepath.Abs(outputPath)
	if err != nil {
		return ""
	}
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.Base(relative) != "index.html" {
		return ""
	}
	slug := filepath.Dir(relative)
	// The homepage writes output_dir/index.html; its slug is "" not ".".
	if slug == "." {
		return ""
	}
	return filepath.ToSlash(slug)
}

func (p *PublishHTMLPlugin) removePostOutputs(sourcePath string, config *lifecycle.Config, postFormats models.PostFormatsConfig, cache *buildcache.Cache) error {
	if sourcePath == "" || config == nil {
		return nil
	}
	if cache == nil {
		return nil
	}
	postCache := cache.GetCachedPost(sourcePath)
	if postCache == nil {
		return nil
	}

	// OutputPath is the authoritative ownership record when available. Fall
	// back to Slug for older cache entries and entries without a rendered HTML
	// output.
	slug := cachedOutputSlug(config.OutputDir, postCache.OutputPath)
	if slug == "" {
		slug = postCache.Slug
	}
	if slug == "" {
		return nil
	}

	outputDir := config.OutputDir
	_ = postFormats
	return removeAllPostOutputs(outputDir, slug)
}

func removeAllPostOutputs(outputDir, slug string) error {
	// Remove only paths generated by PublishHTMLPlugin. A post directory can
	// also contain user-owned or plugin-owned files, so do not remove the whole
	// directory as a side effect of deleting one post.
	paths := []string{
		filepath.Join(slug, "index.html"),
		filepath.Join(slug, "index.md"),
		filepath.Join(slug, "index.txt"),
		filepath.Join(slug, "index.ansi"),
		filepath.Join(slug, "og"),
		slug + ".md",
		slug + ".txt",
		slug + ".ansi",
	}
	for _, relative := range paths {
		path, err := safeOutputPath(outputDir, relative)
		if err != nil {
			return err
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove generated output %s: %w", path, err)
		}
	}

	// Remove the post directory only when it is empty. This keeps unrelated
	// files in the same directory intact.
	postDir, err := safeOutputPath(outputDir, slug)
	if err != nil {
		return err
	}
	return removeEmptyPostDirectory(postDir)
}

func removeEmptyPostDirectory(path string) error {
	entries, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read post directory %s: %w", path, err)
	}
	if len(entries) > 0 {
		return nil
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove empty post directory %s: %w", path, err)
	}
	return nil
}

func safeOutputPath(outputDir, relative string) (string, error) {
	if filepath.IsAbs(relative) {
		return "", fmt.Errorf("refusing absolute output path: %q", relative)
	}
	root, err := filepath.Abs(outputDir)
	if err != nil {
		return "", fmt.Errorf("resolve output directory: %w", err)
	}
	target, err := filepath.Abs(filepath.Join(root, relative))
	if err != nil {
		return "", fmt.Errorf("resolve output path: %w", err)
	}
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("refusing to remove output path outside %s: %q", root, relative)
	}
	resolvedRoot, rootErr := filepath.EvalSymlinks(root)
	resolvedParent, parentErr := filepath.EvalSymlinks(filepath.Dir(target))
	if rootErr == nil && parentErr == nil {
		resolvedRel, relErr := filepath.Rel(resolvedRoot, resolvedParent)
		if relErr != nil || resolvedRel == ".." || strings.HasPrefix(resolvedRel, ".."+string(filepath.Separator)) {
			return "", fmt.Errorf("refusing symlinked output path outside %s: %q", root, relative)
		}
	}
	return target, nil
}

// safePostOutputPath resolves the directory used for a post's HTML output.
// An explicitly empty slug is the homepage and writes directly beneath the
// configured output root; all other slugs must resolve below that root.
func safePostOutputPath(outputDir, slug string) (string, error) {
	if slug == "" {
		return safeOutputRoot(outputDir)
	}
	return safeOutputPath(outputDir, slug)
}

// writeHTMLFormat writes the standard HTML output for a post.
// If incremental build caching is enabled, skips posts that haven't changed.
func (p *PublishHTMLPlugin) writeHTMLFormat(post *models.Post, config *lifecycle.Config, postDir string) (string, error) {
	// Determine HTML content to write
	var htmlContent string
	switch {
	case post.HTML != "":
		// Use pre-rendered HTML if available
		htmlContent = post.HTML
	case post.ArticleHTML != "":
		// Wrap ArticleHTML in a basic template
		htmlContent = p.wrapInTemplate(post, config)
	default:
		// No HTML content available
		return "", nil
	}

	// Write index.html
	outputPath := filepath.Join(postDir, "index.html")

	//nolint:gosec // G306: HTML output files need 0644 for web serving
	if err := os.WriteFile(outputPath, []byte(htmlContent), 0o644); err != nil {
		return "", fmt.Errorf("writing %s: %w", outputPath, err)
	}

	return outputPath, nil
}

// buildFormatContent builds content for a specific output format.
// It checks for per-format template override, and if found, renders using the template.
// Otherwise, falls back to the default content builders for that format.
func (p *PublishHTMLPlugin) buildFormatContent(post *models.Post, config *lifecycle.Config, m *lifecycle.Manager, format string) string {
	// Determine the template to use for this format
	templateName := p.resolveTemplateForFormat(post, format)

	// Check if it's a "raw" template (raw.txt means just output raw content)
	if templateName == rawTxtTemplate || templateName == rawMdTemplate {
		return post.Content
	}

	// Try to get the template engine from cache
	var engine *templates.Engine
	if cached, ok := m.Cache().Get("templates.engine"); ok {
		if e, ok := cached.(*templates.Engine); ok {
			engine = e
		}
	}

	// If we have an engine and the template exists, render it
	if engine != nil && templateName != "" && engine.TemplateExists(templateName) {
		// Create template context
		ctx := templates.NewContext(post, post.Content, applyPostFormatsToConfig(ToModelsConfig(config), resolvePostFormats(post, config)))
		ctx = ctx.WithCore(m)

		// Render the template
		result, err := engine.Render(templateName, ctx)
		if err == nil {
			return result
		}
		// On error, fall through to default content builder
	}

	// Fall back to default content builders
	switch format {
	case formatTxt, formatText:
		return p.buildTextContentFallback(post)
	case formatANSI:
		return p.buildANSIContentFallback(post, config)
	case formatMarkdown, formatMD:
		return buildMarkdownContent(post)
	default:
		return post.Content
	}
}

// resolveTemplateForFormat determines the template to use for a post and output format.
// Uses the TemplatesPlugin if available, otherwise falls back to defaults.
func (p *PublishHTMLPlugin) resolveTemplateForFormat(post *models.Post, format string) string {
	// If we have the templates plugin, use its resolution logic
	if p.templatesPlugin != nil {
		return p.templatesPlugin.resolveTemplateForFormat(post, format)
	}

	// Fallback: check per-format override in frontmatter
	if post.Templates != nil {
		if tmpl, ok := post.Templates[format]; ok && tmpl != "" {
			return tmpl
		}
	}

	// Return empty to use default content builder
	return ""
}

func (p *PublishHTMLPlugin) renderTextContent(post *models.Post, config *lifecycle.Config, engine *templates.Engine) string {
	return p.renderTerminalContent(post, config, engine, false)
}

func (p *PublishHTMLPlugin) renderANSIContent(post *models.Post, config *lifecycle.Config, engine *templates.Engine) string {
	return p.renderTerminalContent(post, config, engine, true)
}

func (p *PublishHTMLPlugin) renderTerminalContent(
	post *models.Post,
	config *lifecycle.Config,
	engine *templates.Engine,
	ansi bool,
) string {
	if isSpecialFile(post.Slug) {
		return post.Content
	}

	body := buildTerminalPage(post, config, ansi)
	if engine == nil {
		return body
	}

	templateName := p.resolveTerminalTemplate(post, engine, ansi)
	if templateName == rawTxtTemplate || templateName == rawAnsiTemplate {
		return post.Content
	}

	ctx := templates.NewContext(post, body, applyPostFormatsToConfig(ToModelsConfig(config), resolvePostFormats(post, config)))
	result, err := engine.Render(templateName, ctx)
	if err != nil {
		return body
	}
	return result
}

func (p *PublishHTMLPlugin) resolveTerminalTemplate(post *models.Post, engine *templates.Engine, ansi bool) string {
	format := formatTxt
	defaultTemplateName := defaultTxtTemplate
	rawTemplateName := rawTxtTemplate
	formatTemplateKey := formatTxt
	shorthandTemplateKey := "txt_template"
	if ansi {
		format = formatANSI
		defaultTemplateName = defaultAnsiTemplate
		rawTemplateName = rawAnsiTemplate
		formatTemplateKey = formatANSI
		shorthandTemplateKey = "ansi_template"
	}

	templateName := p.resolveTemplateForFormat(post, format)
	if templateName != "" && engine.TemplateExists(templateName) {
		return templateName
	}

	if post.Extra != nil {
		if templatesMap, ok := post.Extra["templates"].(map[string]interface{}); ok {
			if formatTemplate, ok := templatesMap[formatTemplateKey].(string); ok && formatTemplate != "" {
				if engine.TemplateExists(formatTemplate) {
					return formatTemplate
				}
			}
		}
		if formatTemplate, ok := post.Extra[shorthandTemplateKey].(string); ok && formatTemplate != "" {
			if engine.TemplateExists(formatTemplate) {
				return formatTemplate
			}
		}
	}

	if isSpecialFile(post.Slug) && engine.TemplateExists(rawTemplateName) {
		return rawTemplateName
	}
	if engine.TemplateExists(defaultTemplateName) {
		return defaultTemplateName
	}
	if engine.TemplateExists(rawTemplateName) {
		return rawTemplateName
	}
	return defaultTemplateName
}

func (p *PublishHTMLPlugin) buildTextContentFallback(post *models.Post) string {
	return buildTerminalPage(post, nil, false)
}

func (p *PublishHTMLPlugin) buildANSIContentFallback(post *models.Post, config *lifecycle.Config) string {
	return buildTerminalPage(post, config, true)
}

func resolveTerminalPalette(config *lifecycle.Config) (string, palettes.Variant) {
	if config == nil {
		return "", palettes.VariantDark
	}
	modelsConfig := ToModelsConfig(config)
	theme := modelsConfig.Theme
	if theme.Palette != "" {
		return theme.Palette, inferPaletteVariant(theme.Palette)
	}
	if theme.PaletteDark != "" {
		return theme.PaletteDark, palettes.VariantDark
	}
	if theme.PaletteLight != "" {
		return theme.PaletteLight, palettes.VariantLight
	}
	return "", palettes.VariantDark
}

func inferPaletteVariant(paletteName string) palettes.Variant {
	lightPatterns := []string{"-light", "-latte", "-dawn", "-day", "-lotus", "modus-operandi"}
	for _, pattern := range lightPatterns {
		if strings.Contains(paletteName, pattern) {
			return palettes.VariantLight
		}
	}
	return palettes.VariantDark
}

// writeReversedFormatOutput writes content for .txt and .md formats.
//
// For special files (robots, llms, humans, security, ads):
// - /slug.<ext> - actual content (e.g., /robots.txt)
// - /slug/index.<ext> - redirect to /slug.<ext> (e.g., /robots/index.txt -> /robots.txt)
// - /slug/index.html - redirect to /slug.<ext> if HTML is disabled
//
// For regular files:
// - /slug.<ext> - actual content (e.g., /test.txt)
// - /slug/index.<ext> - redirect to /slug.<ext> (e.g., /test/index.txt -> /test.txt)
// - /slug/index.html - redirect to /slug.<ext> if HTML is disabled
//
// If skipSlugRedirect is true, the /slug/index.html redirect is skipped
// (used when HTML format is also enabled, since /slug/index.html has the main HTML content).
//
// Fixes: https://github.com/WaylonWalker/markata-go/issues/465
func (p *PublishHTMLPlugin) writeReversedFormatOutput(slug, ext, content, outputDir string, skipSlugRedirect bool) error {
	// The homepage (empty slug) lives at the output root, so its alternate
	// formats are /index.<ext>. There is no /<slug>.<ext> location or redirect
	// directory that would not collide with that file.
	if slug == "" {
		contentPath := filepath.Join(outputDir, "index."+ext)
		if err := removeLegacyRootFormatRedirect(contentPath, ext); err != nil {
			return err
		}
		//nolint:gosec // G306: Output files need 0644 for web serving
		if err := os.WriteFile(contentPath, []byte(content), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", contentPath, err)
		}
		return nil
	}

	// Special files get content at root level (e.g., /robots.txt)
	if isSpecialFile(slug) {
		return p.writeSpecialFileOutput(slug, ext, content, outputDir, skipSlugRedirect)
	}

	// Regular files get content at /slug.<ext> (e.g., /test.txt)
	return p.writeRegularFormatOutput(slug, ext, content, outputDir, skipSlugRedirect)
}

// removeLegacyRootFormatRedirect migrates only the exact redirect directory
// emitted by older homepage publishers. Other directory contents are preserved.
func removeLegacyRootFormatRedirect(path, ext string) error {
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("inspect root format output %s: %w", path, err)
	}
	if !info.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("inspect root format directory %s: %w", path, err)
	}
	if len(entries) != 1 || entries[0].Name() != "index.html" || !entries[0].Type().IsRegular() {
		return nil
	}
	redirectPath := filepath.Join(path, "index.html")
	content, err := os.ReadFile(redirectPath)
	if err != nil {
		return fmt.Errorf("read legacy root redirect %s: %w", redirectPath, err)
	}
	target := "/." + ext
	expected := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta http-equiv="refresh" content="0; url=%s">
<link rel="canonical" href="%s">
<title>Redirecting...</title>
</head>
<body>
<p>Redirecting to <a href="%s">%s</a>...</p>
</body>
</html>`, target, target, target, target)
	if string(content) != expected {
		return nil
	}
	if err := os.Remove(redirectPath); err != nil {
		return fmt.Errorf("remove legacy root redirect %s: %w", redirectPath, err)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove legacy root format directory %s: %w", path, err)
	}
	return nil
}

// writeSpecialFileOutput writes output for special files like robots.txt, llms.txt, etc.
// Content is placed at /slug.<ext> with redirects from /slug/index.<ext>.
func (p *PublishHTMLPlugin) writeSpecialFileOutput(slug, ext, content, outputDir string, skipSlugRedirect bool) error {
	// Write actual content at /slug.<ext> (e.g., /robots.txt)
	contentPath := filepath.Join(outputDir, slug+"."+ext)
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(contentPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", contentPath, err)
	}

	// Ensure slug directory exists for redirects
	slugDir := filepath.Join(outputDir, slug)
	if err := os.MkdirAll(slugDir, 0o755); err != nil {
		return fmt.Errorf("creating slug directory %s: %w", slugDir, err)
	}

	// Create redirect HTML that points to /slug.<ext>
	targetURL := fmt.Sprintf("/%s.%s", slug, ext)
	redirectHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta http-equiv="refresh" content="0; url=%s">
<link rel="canonical" href="%s">
<title>Redirecting...</title>
</head>
<body>
<p>Redirecting to <a href="%s">%s</a>...</p>
</body>
</html>`, targetURL, targetURL, targetURL, targetURL)

	// Create HTML redirect at /slug/index.<ext>/index.html (e.g., /robots/index.txt/index.html)
	// This handles requests to /robots/index.txt by serving the redirect HTML
	extRedirectDir := filepath.Join(slugDir, "index."+ext)
	if err := os.MkdirAll(extRedirectDir, 0o755); err != nil {
		return fmt.Errorf("creating format redirect directory %s: %w", extRedirectDir, err)
	}

	extRedirectPath := filepath.Join(extRedirectDir, "index.html")
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(extRedirectPath, []byte(redirectHTML), 0o644); err != nil {
		return fmt.Errorf("writing format redirect %s: %w", extRedirectPath, err)
	}

	// Create HTML redirect at /slug/index.html if HTML is not enabled
	if !skipSlugRedirect {
		htmlRedirectPath := filepath.Join(slugDir, "index.html")
		//nolint:gosec // G306: Output files need 0644 for web serving
		if err := os.WriteFile(htmlRedirectPath, []byte(redirectHTML), 0o644); err != nil {
			return fmt.Errorf("writing HTML redirect %s: %w", htmlRedirectPath, err)
		}
	}

	return nil
}

// writeRegularFormatOutput writes output for regular (non-special) files.
// Content is placed at /slug.<ext> with redirects from /slug.<ext>/index.html.
func (p *PublishHTMLPlugin) writeRegularFormatOutput(slug, ext, content, outputDir string, skipSlugRedirect bool) error {
	// Write actual content at /slug.<ext> (e.g., /test.txt)
	contentPath := filepath.Join(outputDir, slug+"."+ext)
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(contentPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", contentPath, err)
	}

	// Create redirect HTML that points to /slug.<ext>
	targetURL := fmt.Sprintf("/%s.%s", slug, ext)
	redirectHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta http-equiv="refresh" content="0; url=%s">
<link rel="canonical" href="%s">
<title>Redirecting...</title>
</head>
<body>
<p>Redirecting to <a href="%s">%s</a>...</p>
</body>
</html>`, targetURL, targetURL, targetURL, targetURL)

	// Create redirect at /slug/index.html if HTML is not enabled
	// (when HTML is enabled, /slug/index.html already has the main HTML content)
	if !skipSlugRedirect {
		slugDir := filepath.Join(outputDir, slug)
		if err := os.MkdirAll(slugDir, 0o755); err != nil {
			return fmt.Errorf("creating slug directory %s: %w", slugDir, err)
		}

		redirectPath := filepath.Join(slugDir, "index.html")
		//nolint:gosec // G306: Output files need 0644 for web serving
		if err := os.WriteFile(redirectPath, []byte(redirectHTML), 0o644); err != nil {
			return fmt.Errorf("writing redirect %s: %w", redirectPath, err)
		}
	}

	// Create redirect at /slug/index.<ext> to /slug.<ext>
	slugExtRedirectDir := filepath.Join(outputDir, slug, "index."+ext)
	if err := os.MkdirAll(slugExtRedirectDir, 0o755); err != nil {
		return fmt.Errorf("creating slug index redirect directory %s: %w", slugExtRedirectDir, err)
	}

	slugExtRedirectPath := filepath.Join(slugExtRedirectDir, "index.html")
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(slugExtRedirectPath, []byte(redirectHTML), 0o644); err != nil {
		return fmt.Errorf("writing slug index redirect %s: %w", slugExtRedirectPath, err)
	}

	return nil
}

// writeOGFormat writes the OpenGraph card HTML for social image generation.
func (p *PublishHTMLPlugin) writeOGFormat(post *models.Post, config *lifecycle.Config, postDir string, engine *templates.Engine) error {
	// Create og subdirectory
	ogDir := filepath.Join(postDir, "og")
	if err := os.MkdirAll(ogDir, 0o755); err != nil {
		return fmt.Errorf("creating og directory %s: %w", ogDir, err)
	}

	// Generate OG HTML
	ogHTML := p.generateOGHTML(post, config, engine)

	// Write og/index.html
	outputPath := filepath.Join(ogDir, "index.html")
	//nolint:gosec // G306: HTML output files need 0644 for web serving
	if err := os.WriteFile(outputPath, []byte(ogHTML), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", outputPath, err)
	}

	return nil
}

// generateOGHTML generates OpenGraph card HTML optimized for 1200x630 screenshots.
// It first tries to use a theme template (og-card.html), falling back to a built-in template.
func (p *PublishHTMLPlugin) generateOGHTML(post *models.Post, config *lifecycle.Config, engine *templates.Engine) string {
	// Try theme template first if engine is available
	if engine != nil && engine.TemplateExists("og-card.html") {
		return p.renderOGWithThemeTemplate(post, config, engine)
	}

	// Fall back to built-in template
	return p.renderOGWithBuiltinTemplate(post, config)
}

// renderOGWithThemeTemplate renders the OG card using the theme's og-card.html template.
func (p *PublishHTMLPlugin) renderOGWithThemeTemplate(post *models.Post, config *lifecycle.Config, engine *templates.Engine) string {
	// Build context for pongo2 template
	ctx := templates.NewContext(post, "", applyPostFormatsToConfig(ToModelsConfig(config), resolvePostFormats(post, config)))

	result, err := engine.Render("og-card.html", ctx)
	if err != nil {
		// Log error and fall back to built-in template
		publishHTMLLog.Warnf("failed to render og-card.html template: %v, falling back to built-in", err)
		return p.renderOGWithBuiltinTemplate(post, config)
	}

	return result
}

// renderOGWithBuiltinTemplate renders the OG card using the built-in Go template.
func (p *PublishHTMLPlugin) renderOGWithBuiltinTemplate(post *models.Post, config *lifecycle.Config) string {
	siteTitle := getSiteTitle(config)
	siteURL := getSiteURL(config)

	title := post.PlainTitle()
	if title == "" {
		title = post.Slug
	}

	description := ""
	if post.Description != nil {
		description = *post.Description
	}

	dateStr := ""
	if post.Date != nil {
		dateStr = templates.FormatHumanDate(*post.Date)
	}

	imageURL := getPostExtraString(post, "image", "cover_image", "og_image")
	videoURL := getPostExtraString(post, "video")
	mediaURL := imageURL
	if mediaURL == "" {
		mediaURL = videoURL
	}
	isVideo := templates.IsVideoURL(mediaURL)
	mediaType := ""
	if isVideo {
		mediaType = templates.VideoMIMEType(mediaURL)
	}
	decoratedMediaURL := mediaURL
	if decoratedMediaURL != "" {
		decoratedMediaURL = templates.WithSize(decoratedMediaURL, 360, 320)
	}
	posterURL := ""
	if isVideo {
		posterURL = templates.PosterURLFromMap(templates.GetPostMap(post), mediaURL)
		if posterURL != "" {
			posterURL = templates.WithSize(posterURL, 360, 320)
		}
	}

	// Build canonical URL for the original post
	canonicalURL := siteURL + "/" + post.Slug + "/"

	// Built-in OG card template (1200x630 optimized for social images)
	ogTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=1200, height=630">
    <link rel="canonical" href="{{.CanonicalURL}}">
    <meta name="robots" content="noindex, nofollow">
    <title>{{.Title}} - OG Card</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        html, body {
            width: 1200px;
            height: 630px;
            overflow: hidden;
        }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 60px;
        }
        .og-card {
            background: white;
            border-radius: 20px;
            padding: 60px;
            width: 100%;
            height: 100%;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
        }
        .og-content {
            flex: 1;
            display: flex;
            flex-direction: column;
            justify-content: center;
            gap: 24px;
        }
        .og-layout {
            flex: 1;
            display: flex;
            align-items: center;
            gap: 36px;
        }
        .og-copy {
            flex: 1;
            min-width: 0;
        }
        .og-image-wrap {
            width: 360px;
            height: 320px;
            border-radius: 18px;
            overflow: hidden;
            border: 2px solid #e2e8f0;
            box-shadow: 0 14px 35px -16px rgba(0, 0, 0, 0.45);
            flex-shrink: 0;
        }
        .og-image,
        .og-video {
            width: 100%;
            height: 100%;
            object-fit: cover;
        }
        h1 {
            font-size: 56px;
            font-weight: 700;
            line-height: 1.2;
            color: #1a202c;
            margin-bottom: 24px;
            overflow: hidden;
            display: -webkit-box;
            -webkit-line-clamp: 3;
            -webkit-box-orient: vertical;
        }
        .description {
            font-size: 28px;
            color: #4a5568;
            line-height: 1.5;
            overflow: hidden;
            display: -webkit-box;
            -webkit-line-clamp: 2;
            -webkit-box-orient: vertical;
        }
        .og-footer {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding-top: 30px;
            border-top: 2px solid #e2e8f0;
        }
        .site-name {
            font-size: 24px;
            font-weight: 600;
            color: #667eea;
        }
        .date {
            font-size: 20px;
            color: #718096;
        }
        {{if .Tags}}
        .tags {
            display: flex;
            gap: 12px;
            margin-top: 20px;
            flex-wrap: wrap;
        }
        .tag {
            background: #edf2f7;
            color: #4a5568;
            padding: 8px 16px;
            border-radius: 20px;
            font-size: 18px;
        }
        {{end}}
    </style>
</head>
<body>
    <div class="og-card">
        <div class="og-content">
            <div class="og-layout">
                <div class="og-copy">
                    <h1>{{.Title}}</h1>
                    {{if .Description}}<p class="description">{{.Description}}</p>{{end}}
                    {{if .Tags}}
                    <div class="tags">
                        {{range .TagsDisplay}}
                        <span class="tag">{{.}}</span>
                        {{end}}
                    </div>
                    {{end}}
                </div>
                {{if .HasMedia}}
                <div class="og-image-wrap">
                    {{if .IsVideo}}
                    <video class="og-video" autoplay muted loop playsinline{{if .Poster}} poster="{{.Poster}}"{{end}}>
                        <source src="{{.MediaURL}}"{{if .MediaType}} type="{{.MediaType}}"{{end}}>
                    </video>
                    {{else}}
                    <img src="{{.MediaURL}}" alt="{{.Title}}" class="og-image">
                    {{end}}
                </div>
                {{end}}
            </div>
        </div>
        <div class="og-footer">
            <span class="site-name">{{.SiteTitle}}</span>
            {{if .DateStr}}<span class="date">{{.DateStr}}</span>{{end}}
        </div>
    </div>
</body>
</html>`

	tmpl, err := template.New("og").Parse(ogTemplate)
	if err != nil {
		// Return minimal HTML on template error
		return fmt.Sprintf("<html><body><h1>%s</h1></body></html>",
			template.HTMLEscapeString(title))
	}

	// Limit tags to first 3 for display
	tagsDisplay := post.Tags
	if len(tagsDisplay) > 3 {
		tagsDisplay = tagsDisplay[:3]
	}

	data := struct {
		Title        string
		Description  string
		DateStr      string
		Tags         []string
		TagsDisplay  []string
		HasMedia     bool
		IsVideo      bool
		MediaType    string
		MediaURL     string
		Poster       string
		SiteTitle    string
		CanonicalURL string
	}{
		Title:        title,
		Description:  description,
		DateStr:      dateStr,
		Tags:         post.Tags,
		TagsDisplay:  tagsDisplay,
		HasMedia:     mediaURL != "",
		IsVideo:      isVideo,
		MediaType:    mediaType,
		MediaURL:     decoratedMediaURL,
		Poster:       posterURL,
		SiteTitle:    siteTitle,
		CanonicalURL: canonicalURL,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// Return minimal HTML on execution error
		return fmt.Sprintf("<html><body><h1>%s</h1></body></html>",
			template.HTMLEscapeString(title))
	}

	return buf.String()
}

// wrapInTemplate wraps post content in a basic HTML template.
func (p *PublishHTMLPlugin) wrapInTemplate(post *models.Post, config *lifecycle.Config) string {
	siteURL := getSiteURL(config)
	siteTitle := getSiteTitle(config)

	title := post.PlainTitle()
	if title == "" {
		title = post.Slug
	}

	description := ""
	if post.Description != nil {
		description = *post.Description
	}

	dateStr := ""
	dateISO := ""
	if post.Date != nil {
		dateStr = templates.FormatHumanDate(*post.Date)
		dateISO = post.Date.Format("2006-01-02")
	}
	fontpack := "system"
	if value, ok := config.Extra["fontpack"].(string); ok && value != "" {
		fontpack = value
	}
	if value, ok := post.Extra["_resolved_fontpack"].(string); ok && value != "" {
		fontpack = value
	}

	// Simple default template
	tmplStr := `<!DOCTYPE html>
<html lang="en" data-fontpack="{{.Fontpack}}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} | {{.SiteTitle}}</title>
    {{if .Description}}<meta name="description" content="{{.Description}}">{{end}}
    <link rel="canonical" href="{{.SiteURL}}{{.Href}}">
    <link rel="stylesheet" href="/css/variables.css">
    <link rel="stylesheet" href="/css/palette.css">
    <link rel="stylesheet" href="/css/main.css">
    <link rel="stylesheet" href="/css/aesthetic.css">
    <link rel="stylesheet" href="/css/code.css">
    <link rel="stylesheet" href="/css/admonitions.css">
    <link rel="stylesheet" href="/css/fonts.css">
</head>
<body>
    <header>
        <nav>
            <a href="{{.SiteURL}}">{{.SiteTitle}}</a>
        </nav>
    </header>
    <div class="article-progress" role="presentation" aria-hidden="true"><span class="article-progress__indicator"></span></div>
    <main>
        <article class="post">
            <header>
                <h1>{{.Title}}</h1>
                {{if .DateStr}}<time datetime="{{.DateISO}}">{{.DateStr}}</time>{{end}}
            </header>
            <div class="content">
                {{.Content}}
            </div>
            {{if .Tags}}
            <footer>
                <ul class="tags">
                {{range .Tags}}
                    <li><a href="/tags/{{.}}/">{{.}}</a></li>
                {{end}}
                </ul>
            </footer>
            {{end}}
        </article>
    </main>
    <script>
      (function initArticleProgressIndicator() {
        var indicator = document.querySelector('.article-progress__indicator');
        var article = document.querySelector('article.post');
        if (!indicator || !article) {
          return;
        }

        var raf = window.requestAnimationFrame ? window.requestAnimationFrame.bind(window) : function(cb) { return window.setTimeout(cb, 16); };
        var scheduled = false;
        var lastValue = -1;

        function clamp(value, min, max) {
          if (value < min) return min;
          if (value > max) return max;
          return value;
        }

        function refresh() {
          scheduled = false;
          var rect = article.getBoundingClientRect();
          var articleTop = window.scrollY + rect.top;
          var articleHeight = rect.height;
          var viewportHeight = window.innerHeight;
          var scrollableDistance = articleHeight - viewportHeight;
          var ratio;
          if (scrollableDistance <= 0) {
            ratio = 1;
          } else {
            ratio = clamp((window.scrollY - articleTop) / scrollableDistance, 0, 1);
          }
          if (Math.abs(ratio - lastValue) < 0.001) {
            return;
          }
          lastValue = ratio;
          indicator.style.transform = 'scaleX(' + ratio + ')';
        }

        function schedule() {
          if (scheduled) {
            return;
          }
          scheduled = true;
          raf(refresh);
        }

        window.addEventListener('scroll', schedule, { passive: true });
        window.addEventListener('resize', schedule);
        schedule();
      })();
    </script>
    <footer>
        <p><a href="{{.SiteURL}}">{{.SiteTitle}}</a></p>
    </footer>
</body>
</html>`

	tmpl, err := template.New("post").Parse(tmplStr)
	if err != nil {
		// Return basic HTML on template error
		return fmt.Sprintf("<html><head><title>%s</title></head><body>%s</body></html>",
			template.HTMLEscapeString(title), post.ArticleHTML)
	}

	data := struct {
		Title       string
		Description string
		Content     template.HTML
		DateStr     string
		DateISO     string
		Tags        []string
		Href        string
		SiteURL     string
		SiteTitle   string
		Fontpack    string
	}{
		Title:       title,
		Description: description,
		Content:     template.HTML(post.ArticleHTML), //nolint:gosec // G203: ArticleHTML is sanitized markdown output
		DateStr:     dateStr,
		DateISO:     dateISO,
		Tags:        post.Tags,
		Href:        post.Href,
		SiteURL:     siteURL,
		SiteTitle:   siteTitle,
		Fontpack:    fontpack,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		// Return basic HTML on execution error
		return fmt.Sprintf("<html><head><title>%s</title></head><body>%s</body></html>",
			template.HTMLEscapeString(title), post.ArticleHTML)
	}

	return buf.String()
}

// Ensure PublishHTMLPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin      = (*PublishHTMLPlugin)(nil)
	_ lifecycle.WritePlugin = (*PublishHTMLPlugin)(nil)
)

// getPostFormatsConfig extracts PostFormatsConfig from lifecycle.Config.Extra.
func getPostFormatsConfig(config *lifecycle.Config) models.PostFormatsConfig {
	if config.Extra != nil {
		if pf, ok := config.Extra["post_formats"].(models.PostFormatsConfig); ok {
			return pf
		}
	}
	return models.NewPostFormatsConfig()
}
