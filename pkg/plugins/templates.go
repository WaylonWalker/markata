package plugins

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/sidebar"
	"github.com/WaylonWalker/markata/pkg/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var templatesLog = logging.Component("templates").Phase("render")

// Output format constants.
const (
	formatHTML          = "html"
	formatTxt           = "txt"
	formatText          = "text"
	formatANSI          = "ansi"
	formatMarkdown      = "markdown"
	formatMD            = "md"
	formatOG            = "og"
	feedPriorityPrivate = "private"

	// defaultTemplate is the default template name for posts.
	defaultTemplate = "post.html"

	// fullHTMLRestoreConcurrency bounds simultaneous full-page cache reads.
	// Manager concurrency may further reduce the number of workers.
	// Keep parallel disk reads bounded independently of rendering workers.
	fullHTMLRestoreConcurrency = 4
)

// TemplatesPlugin wraps rendered markdown content in HTML templates.
// It operates during the render stage, after markdown has been converted to HTML.
type TemplatesPlugin struct {
	engine        *templates.Engine
	layoutConfig  *models.LayoutConfig
	config        *lifecycle.Config
	localPreviews map[string]string
	siteURL       string
	sidebar       *sidebar.Builder
	sidebarPostMu sync.RWMutex
	sidebarPosts  map[sidebarPostKey]sidebarPostJSON
}

// NewTemplatesPlugin creates a new templates plugin.
func NewTemplatesPlugin() *TemplatesPlugin {
	return &TemplatesPlugin{}
}

// Name returns the plugin name.
func (p *TemplatesPlugin) Name() string {
	return PluginNameTemplates
}

// Configure initializes the template engine from the config.
func (p *TemplatesPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if modelsConfig, ok := config.Extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		templates.SetTrustedMediaDomains(modelsConfig.Templates.Media.TrustedDomains)
	}

	// Get templates directory from config
	templatesDir := PluginNameTemplates
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}

	// Get theme name from config (default to "default")
	themeName := ThemeDefault
	if extra := config.Extra; extra != nil {
		// Check for typed ThemeConfig struct (set by core.go)
		if theme, ok := extra["theme"].(models.ThemeConfig); ok {
			if theme.Name != "" {
				themeName = theme.Name
			}
		}
		// Also check for map[string]interface{} (legacy/dynamic config)
		if theme, ok := extra["theme"].(map[string]interface{}); ok {
			if name, ok := theme["name"].(string); ok && name != "" {
				themeName = name
			}
		}
		// Also check for simple theme string
		if name, ok := extra["theme"].(string); ok && name != "" {
			themeName = name
		}
	}

	// Initialize template engine with theme support
	engine, err := templates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return fmt.Errorf("failed to initialize template engine: %w", err)
	}
	p.engine = engine

	// Store engine in cache for other plugins to use
	m.Cache().Set("templates.engine", engine)

	// Store this plugin in cache for other plugins to use per-format resolution
	m.Cache().Set("templates.plugin", p)

	// Get layout config if available
	switch lc := config.Extra["layout"].(type) {
	case *models.LayoutConfig:
		p.layoutConfig = lc
	case models.LayoutConfig:
		p.layoutConfig = &lc
	}

	// Store config reference for template preset resolution
	p.config = config

	return nil
}

// resolveTemplate determines the template to use for a post (HTML format).
// This is a convenience wrapper for resolveTemplateForFormat with "html" format.
// Priority: frontmatter per-format -> preset -> simple template -> layout config -> global default
func (p *TemplatesPlugin) resolveTemplate(post *models.Post) string {
	return p.resolveTemplateForFormat(post, "html")
}

// resolveTemplateForFormat determines the template to use for a post and output format.
// Resolution priority:
// 1. Frontmatter per-format override (templates.html, templates.txt, etc.)
// 2. Frontmatter template preset (template: blog → expand to preset)
// 3. Frontmatter simple template (template: post.html → use for current format)
// 4. Layout config (path/feed-based)
// 5. Global default for format (default_templates.html, etc.)
// 6. Hardcoded default (post.html, default.txt, etc.)
//
//nolint:gocyclo // Template precedence is intentionally explicit per format.
func (p *TemplatesPlugin) resolveTemplateForFormat(post *models.Post, format string) string {
	// The canonical specimen uses the post template preset. Keep its
	// projection explicit so a preset resolved from frontmatter cannot fall
	// through to a generic post layout.
	if format == formatHTML && isCanonicalRenderingFixture(post) {
		return defaultTemplate
	}

	// 1. Check per-format override in frontmatter
	if post.Templates != nil {
		if tmpl, ok := post.Templates[format]; ok && tmpl != "" {
			return tmpl
		}
	}

	// 2. Check if template is a preset name
	if post.Template != "" && p.config != nil {
		presets := getTemplatePresets(p.config)
		if preset, ok := presets[post.Template]; ok {
			tmpl := preset.TemplateForFormat(format)
			if tmpl != "" {
				return tmpl
			}
		}
	}

	// 3. Use template as explicit file (if has extension) and adapt for format
	if post.Template != "" && strings.Contains(post.Template, ".") {
		return adaptTemplateForFormat(post.Template, format)
	}

	// 4. Use template as-is if it doesn't have an extension
	// This might be a preset name that wasn't found, fall through to layout
	if post.Template != "" {
		// For HTML, use the template directly
		if format == formatHTML {
			return post.Template + ".html"
		}
		// For other formats, try to adapt it
		return adaptTemplateForFormat(post.Template+".html", format)
	}

	// 5. Use layout configuration to determine template
	if p.layoutConfig != nil {
		// Get feed slug for feed-based layout lookup
		feedSlug := post.PrevNextFeed
		if feedSlug == "" {
			if feed, ok := post.Extra["feed"].(string); ok {
				feedSlug = feed
			}
		}

		// Get post path for path-based layout lookup
		postPath := post.Href
		if postPath == "" {
			postPath = "/" + strings.TrimPrefix(post.Path, "/")
		}

		// Resolve layout based on path and feed
		layout := p.layoutConfig.ResolveLayout(postPath, feedSlug)
		if layout != "" {
			baseTemplate := models.LayoutToTemplate(layout)
			return adaptTemplateForFormat(baseTemplate, format)
		}
	}

	// 6. Check global default templates from config
	if p.config != nil {
		defaultTemplates := getDefaultTemplates(p.config)
		if tmpl, ok := defaultTemplates[format]; ok && tmpl != "" {
			return tmpl
		}
	}

	// 7. Fall back to hardcoded defaults per format
	return getHardcodedDefault(format)
}

// adaptTemplateForFormat adapts a template name for a specific output format.
// For example: post.html → post.txt, post.md, post-og.html
func adaptTemplateForFormat(template, format string) string {
	ext := filepath.Ext(template)
	base := strings.TrimSuffix(template, ext)

	switch format {
	case formatHTML:
		return template
	case formatTxt, formatText:
		return base + ".txt"
	case formatANSI:
		return base + ".ansi"
	case formatMarkdown, formatMD:
		return base + ".md"
	case formatOG:
		return base + "-og.html"
	default:
		return template
	}
}

// getHardcodedDefault returns the hardcoded default template for a format.
func getHardcodedDefault(format string) string {
	switch format {
	case formatHTML:
		return defaultTemplate
	case formatTxt, formatText:
		return "default.txt"
	case formatANSI:
		return "default.ansi"
	case formatMarkdown, formatMD:
		return "raw.txt"
	case formatOG:
		return "og-card.html"
	default:
		return defaultTemplate
	}
}

// getTemplatePresets extracts TemplatePresets from lifecycle.Config.Extra.
func getTemplatePresets(config *lifecycle.Config) map[string]models.TemplatePreset {
	if config.Extra == nil {
		return nil
	}
	if presets, ok := config.Extra["template_presets"].(map[string]models.TemplatePreset); ok {
		return presets
	}
	return nil
}

// getDefaultTemplates extracts DefaultTemplates from lifecycle.Config.Extra.
func getDefaultTemplates(config *lifecycle.Config) map[string]string {
	if config.Extra == nil {
		return nil
	}
	if defaults, ok := config.Extra["default_templates"].(map[string]string); ok {
		return defaults
	}
	return nil
}

// Render wraps markdown content in templates.
// This runs after markdown rendering, using post.ArticleHTML as the body.
// Skips posts that don't need rebuilding (incremental builds).
//
// Uses three-phase processing for incremental optimization:
// ensureFeedConfigsCached pre-computes feed configs and caches them if
// the Collect stage (which normally does this) hasn't run yet.
// This allows the feed sidebar auto-discovery to work during the Render stage.
func ensureFeedConfigsCached(config *lifecycle.Config, m *lifecycle.Manager) {
	if _, ok := m.Cache().Get("feed_configs"); ok {
		return // Already cached (e.g., Collect ran first)
	}

	feedConfigs := getFeedConfigs(config)
	if len(feedConfigs) == 0 {
		return
	}

	feedDefaults := getFeedDefaults(config)
	posts := m.Posts()
	fc := newFeedFilterCache(posts)

	for i := range feedConfigs {
		feedCfg := &feedConfigs[i]
		feedCfg.ApplyDefaults(feedDefaults)

		filteredPosts, err := fc.FilterPosts(feedCfg.Filter, feedCfg.IncludesPrivate())
		if err != nil {
			continue
		}

		// Sort posts (same logic as feeds.go)
		sortField := feedCfg.Sort
		reverse := feedCfg.Reverse
		if sortField == "" {
			sortField = "date"
			reverse = true
		}
		sorted := cloneFeedPosts(filteredPosts)
		sortPosts(sorted, sortField, reverse)
		feedCfg.Posts = sorted
	}

	m.Cache().Set("feed_configs", feedConfigs)
}

// Phase 1a: Quick single-threaded pass to classify posts (no disk I/O)
// Phase 1b: Restore cached full-page HTML for unchanged posts.
// Phase 2: Concurrent rendering only for posts that need it
func (p *TemplatesPlugin) Render(m *lifecycle.Manager) error {
	ledger := m.ContentLedger()
	ledger.SetTemplateCache(nil)
	if p.engine == nil {
		return fmt.Errorf("template engine not initialized")
	}
	p.resetSidebarPosts()

	// Get config for template context
	config := m.Config()

	// Ensure feed_configs are available in cache for sidebar auto-discovery.
	// The Collect stage (feeds) runs AFTER Render, so we pre-compute here.
	ensureFeedConfigsCached(config, m)
	var navFeeds []models.FeedConfig
	if feedConfigs, ok := m.Cache().Get("feed_configs"); ok {
		if feeds, ok := feedConfigs.([]models.FeedConfig); ok {
			navFeeds = feeds
		}
	}
	navFeeds = previewFeeds(config, m.Posts(), navFeeds)
	modelsConfig := ToModelsConfig(config)
	navPreviews := buildNavPreviews(modelsConfig, m.Posts(), navFeeds, getBlogrollConfig(config), parseRandomPostConfig(config))
	feedMap := make(map[string]*models.FeedConfig, len(navFeeds))
	for i := range navFeeds {
		feedMap[navFeeds[i].Slug] = &navFeeds[i]
	}
	p.sidebar = sidebar.NewBuilder(modelsConfig, feedMap, m.Posts())
	m.Cache().Set("nav_previews", navPreviews)
	p.localPreviews = encodeLocalPreviews(buildLocalPreviews(m.Posts(), navFeeds))
	p.siteURL = ToModelsConfig(config).URL

	// Get build cache to check if posts need rebuilding
	cache := GetBuildCache(m)
	var stats diagnostics.TemplateCacheStats
	if cache != nil {
		encoded, err := json.Marshal(navPreviews)
		if err != nil {
			return fmt.Errorf("encode navigation previews: %w", err)
		}
		stats.NavPreviewReset = cache.SetNavPreviewHash(buildcache.ContentHash(string(encoded)))
	}
	changedSlugs := getChangedSlugsMap(cache)
	affected := lifecycle.GetServeAffectedPaths(m)
	var canonicalPaths map[string]bool
	if value, ok := m.Cache().Get(canonicalArticlePathsKey); ok {
		var valid bool
		canonicalPaths, valid = value.(map[string]bool)
		if !valid {
			return fmt.Errorf("canonical article paths have invalid type %T", value)
		}
	}

	// Collect private paths for robots.txt template variable
	privatePaths := collectPrivatePaths(m.Posts())

	// Pre-compute feed membership hashes once (O(N)) instead of per-post (O(N^2))
	feedMembershipHashes := precomputeFeedMembershipHashes(config, m)

	// Phase 1a: Classify posts into "cacheable" vs "needs rendering" without disk I/O.
	t0 := time.Now()
	cacheablePosts := make([]*models.Post, 0, len(m.Posts()))
	var postsNeedingRender []*models.Post

	for _, post := range m.Posts() {
		// Skip posts marked to skip. Posts with an empty body still get a page
		// (title, metadata, feed membership) so feed links never dangle.
		if post.Skip {
			stats.Skipped++
			continue
		}

		// Check if we can use cached HTML (no disk I/O -- just map lookups)
		stats.Classified++
		reason := p.templateCacheReason(post, cache, changedSlugs, feedMembershipHashes, affected)
		if reason == templateCacheHit {
			stats.Cacheable++
			cacheablePosts = append(cacheablePosts, post)
		} else {
			recordTemplateCacheMiss(&stats.MissReasons, reason)
			postsNeedingRender = append(postsNeedingRender, post)
		}
	}
	t1 := time.Now()
	templatesLog.Printf("Phase 1a classify: %d cacheable, %d need render (took %v)", len(cacheablePosts), len(postsNeedingRender), t1.Sub(t0))

	// Phase 1b: Restore the full page before the write stage. The output tree may
	// have been removed while the persistent cache remains, so later stages must
	// receive the same HTML that a cold build would have rendered. If the full
	// page cache is unavailable (for example, from an older cache format), render
	// the post instead of treating it as successfully restored.
	restoreStart := time.Now()
	restored, unavailable, err := restoreCachedFullHTML(m, cacheablePosts, cache.GetCachedFullHTML, fullHTMLRestoreConcurrency)
	if err != nil {
		return fmt.Errorf("restore full-page cache: %w", err)
	}
	stats.Restored = restored
	stats.MissReasons.FullHTMLUnavailable = len(unavailable)
	postsNeedingRender = append(postsNeedingRender, unavailable...)
	restoreEnd := time.Now()
	templatesLog.Printf("Phase 1b batch restore: took %v, %d now need render", restoreEnd.Sub(restoreStart), len(postsNeedingRender))
	stats.RenderRequired = len(postsNeedingRender)

	// A missing article cache does not authorize rendering an unrelated page
	// outside the incremental selection. Full-page hits above remain usable,
	// but fresh pages require an article handled by Markdown in this pass.
	if lifecycle.IsServeIncremental(m) && canonicalPaths != nil {
		filtered := postsNeedingRender[:0]
		for _, post := range postsNeedingRender {
			if canonicalPaths[post.Path] {
				filtered = append(filtered, post)
			} else {
				stats.ServeDeferred++
			}
		}
		postsNeedingRender = filtered
	}

	// Phase 2: Process only posts that need rendering concurrently
	// The worker pool processes every selected post even after an error.
	renderStart := time.Now()
	var renderFailed atomic.Int64
	err = m.ProcessPostsSliceConcurrently(postsNeedingRender, func(post *models.Post) error {
		// Render the template
		html, err := p.renderPost(post, config, m, privatePaths)
		if err != nil {
			renderFailed.Add(1)
			return err
		}
		post.HTML = html

		p.cacheRenderedPage(post, cache, feedMembershipHashes, canonicalPaths)

		return nil
	})
	t3 := time.Now()
	templatesLog.Printf("Phase 2 render: took %v", t3.Sub(renderStart))
	stats.RenderFailed = int(renderFailed.Load())
	stats.RenderSucceeded = len(postsNeedingRender) - stats.RenderFailed
	ledger.SetTemplateCache(&stats)
	return err
}

func (p *TemplatesPlugin) cacheRenderedPage(post *models.Post, cache *buildcache.Cache, feedMembershipHashes map[string]string, canonicalPaths map[string]bool) {
	if cache == nil || post.InputHash == "" {
		return
	}
	cacheErr := cache.CacheFullHTML(post.Path, post.HTML)
	cache.SetLocalPreviewHash(post.Path, localPreviewHash(post, p.localPreviews, p.siteURL))
	if membershipHash := lookupFeedMembershipHash(post, feedMembershipHashes); membershipHash != "" {
		cache.SetFeedMembershipHash(post.Path, membershipHash)
	}
	if cacheErr == nil && post.Path != "" && !isSourceEncryptedPost(post) &&
		(canonicalPaths == nil || canonicalPaths[post.Path]) {
		cache.SetHeadingHighlightRevision(post.Path, headingHighlightRevision)
	}
}

// canUseCachedHTML checks if a post can use cached HTML without doing any disk I/O.
// This is the "decision" phase that determines cache eligibility.
func canUseCachedHTML(post *models.Post, cache *buildcache.Cache, changedSlugs map[string]bool, feedMembershipHashes map[string]string) bool {
	return cachedHTMLReason(post, cache, changedSlugs, feedMembershipHashes) == templateCacheHit
}

// templateCacheReason preserves the outer affected/local-preview gates. Neither
// it nor cachedHTMLReason evaluates gates masked by an earlier miss.
func (p *TemplatesPlugin) templateCacheReason(post *models.Post, cache *buildcache.Cache, changedSlugs map[string]bool, feedMembershipHashes map[string]string, affected map[string]bool) templateCacheReason {
	if affected[post.Path] {
		return templateCacheAffectedPath
	}
	if reason := cachedHTMLReason(post, cache, changedSlugs, feedMembershipHashes); reason != templateCacheHit {
		return reason
	}
	if cache.GetLocalPreviewHash(post.Path) != localPreviewHash(post, p.localPreviews, p.siteURL) {
		return templateCacheLocalPreviewChanged
	}
	return templateCacheHit
}

type templateCacheReason uint8

const (
	templateCacheHit templateCacheReason = iota
	templateCacheAffectedPath
	templateCacheUnavailable
	templateCacheInputHashMissing
	templateCacheEntryMissing
	templateCacheInputHashMismatch
	templateCacheTemplateMismatch
	templateCacheDependencyChanged
	templateCacheSlugChanged
	templateCacheFeedMembershipChanged
	templateCacheLocalPreviewChanged
)

func cachedHTMLReason(post *models.Post, cache *buildcache.Cache, changedSlugs map[string]bool, feedMembershipHashes map[string]string) templateCacheReason {
	if cache == nil {
		return templateCacheUnavailable
	}
	if post.InputHash == "" {
		return templateCacheInputHashMissing
	}

	// Check if post itself changed
	switch cache.ReasonForRebuild(post.Path, post.InputHash, post.Template) {
	case buildcache.RebuildReasonMissingEntry:
		return templateCacheEntryMissing
	case buildcache.RebuildReasonInputChanged:
		return templateCacheInputHashMismatch
	case buildcache.RebuildReasonTemplateChanged:
		return templateCacheTemplateMismatch
	case buildcache.RebuildReasonNone:
		// Continue with the existing post-local gates.
	}

	// Check if any dependency changed
	if len(changedSlugs) > 0 {
		for _, dep := range post.Dependencies {
			if changedSlugs[dep] {
				return templateCacheDependencyChanged
			}
		}
		// Check if this post's slug is in changedSlugs
		if changedSlugs[post.Slug] {
			return templateCacheSlugChanged
		}
	}

	// Check if feed membership changed (for sidebar invalidation)
	if currentHash := lookupFeedMembershipHash(post, feedMembershipHashes); currentHash != "" {
		cachedHash := cache.GetFeedMembershipHash(post.Path)
		if cachedHash != currentHash {
			return templateCacheFeedMembershipChanged
		}
	}

	return templateCacheHit
}

func recordTemplateCacheMiss(reasons *diagnostics.TemplateCacheMissReasons, reason templateCacheReason) {
	switch reason {
	case templateCacheAffectedPath:
		reasons.AffectedPath++
	case templateCacheUnavailable:
		reasons.CacheUnavailable++
	case templateCacheInputHashMissing:
		reasons.InputHashMissing++
	case templateCacheEntryMissing:
		reasons.EntryMissing++
	case templateCacheInputHashMismatch:
		reasons.InputHashMismatch++
	case templateCacheTemplateMismatch:
		reasons.TemplateMismatch++
	case templateCacheDependencyChanged:
		reasons.DependencyChanged++
	case templateCacheSlugChanged:
		reasons.SlugChanged++
	case templateCacheFeedMembershipChanged:
		reasons.FeedMembershipChanged++
	case templateCacheLocalPreviewChanged:
		reasons.LocalPreviewChanged++
	case templateCacheHit:
		// Hits have no miss reason.
	}
}

// getChangedSlugsMap returns a map of slugs that changed in this build.
func getChangedSlugsMap(cache *buildcache.Cache) map[string]bool {
	if cache == nil {
		return nil
	}
	changedSlugs := make(map[string]bool)
	for _, slug := range cache.GetChangedSlugs() {
		changedSlugs[slug] = true
	}
	return changedSlugs
}

// precomputeFeedMembershipHashes builds a map of tag -> membership hash in O(N) time.
// This replaces the per-post O(N) scan that previously made the overall complexity O(N^2).
func precomputeFeedMembershipHashes(config *lifecycle.Config, m *lifecycle.Manager) map[string]string {
	components, ok := config.Extra["components"].(models.ComponentsConfig)
	if !ok {
		return nil
	}
	if components.FeedSidebar.Enabled == nil || !*components.FeedSidebar.Enabled {
		return nil
	}
	feedSlugs := components.FeedSidebar.Feeds
	if len(feedSlugs) == 0 {
		return nil
	}

	// Collect the tag names we care about
	tagNames := make(map[string]bool, len(feedSlugs))
	for _, feedSlug := range feedSlugs {
		if strings.HasPrefix(feedSlug, "tags/") {
			tagNames[strings.TrimPrefix(feedSlug, "tags/")] = true
		}
	}
	if len(tagNames) == 0 {
		return nil
	}

	// Single pass over all posts: group member slugs by tag
	tagMembers := make(map[string][]string, len(tagNames))
	for _, post := range m.Posts() {
		if !post.Published || post.Draft || post.Skip {
			continue
		}
		for _, t := range post.Tags {
			if tagNames[t] {
				tagMembers[t] = append(tagMembers[t], post.Slug)
			}
		}
	}

	// Compute hash per tag
	result := make(map[string]string, len(tagMembers))
	for tag, slugs := range tagMembers {
		result[tag] = buildcache.ComputeFeedMembershipHash(slugs)
	}
	return result
}

// lookupFeedMembershipHash finds the feed membership hash for a post using the pre-computed map.
// Returns the hash of the first matching tag's membership, or empty string.
func lookupFeedMembershipHash(post *models.Post, feedMembershipHashes map[string]string) string {
	if len(feedMembershipHashes) == 0 {
		return ""
	}
	for _, tag := range post.Tags {
		if h, ok := feedMembershipHashes[tag]; ok {
			return h
		}
	}
	return ""
}

// canSkipPostCardStyles is intentionally conservative: only the stock reading
// template with no card/feed/component markers can safely skip cards.css.
// This avoids changing custom templates or a post embedding a feed widget.
func canSkipPostCardStyles(post *models.Post, templateName string) bool {
	if post == nil || templateName != defaultTemplate {
		return false
	}
	content := strings.ToLower(post.ArticleHTML + " " + post.Content)
	for _, marker := range []string{
		// Intentionally favor extra CSS over a styling regression. The file
		// also styles photo-card, simple-list and span/aspect utilities.
		"card", "feed", "simple-list", "simple-item", "col-span",
		"row-span", "aspect-auto", "render_slashes", "include_post",
	} {
		if strings.Contains(content, marker) {
			return false
		}
	}
	return true
}

// renderPost renders a single post using the appropriate template.
func (p *TemplatesPlugin) renderPost(post *models.Post, config *lifecycle.Config, m *lifecycle.Manager, privatePaths []string) (string, error) {
	// Determine which template to use
	templateName := p.resolveTemplate(post)

	// Check if template exists, fall back to post.html if not
	if !p.engine.TemplateExists(templateName) {
		templateName = defaultTemplate
		if !p.engine.TemplateExists(templateName) {
			return post.ArticleHTML, nil
		}
	}

	// Create template context
	modelsConfig := applyPostFormatsToConfig(ToModelsConfig(config), resolvePostFormats(post, config))
	articleHTML := annotateLocalPreviewLinks(post.ArticleHTML, post.Href, p.siteURL, p.localPreviews)
	ctx := templates.NewContext(post, articleHTML, modelsConfig)
	if p.sidebar != nil {
		items, title := p.sidebar.ResolveForPost(post)
		ctx = ctx.WithSidebar(items, title)
	}
	tagPreviews := make(map[string]string, len(post.Tags))
	for _, tag := range post.Tags {
		if preview := localPreviewForTag(tag, p.localPreviews); preview != "" {
			tagPreviews[tag] = preview
		}
	}
	ctx.Set("tag_preview_json", tagPreviews)
	ctx = ctx.WithCore(m).WithBlogLayout(templateName == models.LayoutToTemplate("blog"))
	ctx.Set("canonical_rendering_fixture", isCanonicalRenderingFixture(post))
	// Feed card CSS is not on a normal article's critical rendering path.
	// Unknown/custom templates and embedded card/feed markup keep the CSS.
	ctx.Set("skip_cards_css", canSkipPostCardStyles(post, templateName))
	ctx.Set("feed_posts", createFeedPostsFunc(m))
	ctx.Set("render_feed", createRenderFeedFunc(m))
	ctx.Set("render_slashes", createRenderSlashesFunc(m))
	ctx.Set("include_post", createIncludePostFunc(m))
	ctx.Set("private_paths", privatePaths)
	if previews, ok := m.Cache().Get("nav_previews"); ok {
		ctx.Set("nav_previews", previews)
	}
	if modelsConfig.Garden.IsExportJSON() {
		ctx.Set("graph_json", "/"+modelsConfig.Garden.GetPath()+"/graph.json")
	}

	// Share buttons at the end of posts
	if modelsConfig.Components.Share.IsEnabled() {
		shareButtons := models.BuildShareButtons(modelsConfig.Components.Share, modelsConfig.URL, modelsConfig.Title, post)
		if len(shareButtons) > 0 {
			ctx.Set("share_buttons", shareButtons)
		}
	}

	if modelsConfig.Components.PostCopy.IsEnabled() {
		postCopyPayloads := buildPostCopyPayloads(post, config, modelsConfig.URL)
		ctx.Set("post_copy_payloads", postCopyPayloads)
		ctx.Set("post_copy_payloads_json", postCopyPayloads.JSON())
	}

	// Inject feed sidebar posts if configured
	sidebarPosts, sidebarFeed := p.getFeedSidebarPosts(post, config, m)
	if sidebarPosts != nil {
		ctx.Set("sidebar_posts", sidebarPosts)
		if sidebarFeed != nil {
			ctx.Set("sidebar_feed", sidebarFeed)
			ctx.Set("sidebar_calendar_enabled", sidebarFeed.HasView(models.FeedViewCalendar))
		}

		// Calculate prev/next within the sidebar feed
		sidebarPrev, sidebarNext := p.getSidebarPrevNext(post, sidebarPosts)
		if sidebarPrev != nil {
			ctx.Set("sidebar_prev", sidebarPrev)
		}
		if sidebarNext != nil {
			ctx.Set("sidebar_next", sidebarNext)
		}

		// Build JSON for all candidate feeds (for client-side {/} cycling)
		feedsJSON := p.buildSidebarFeedsJSON(post, config, m, sidebarFeed)
		if feedsJSON != "" {
			ctx.Set("sidebar_feeds_json", feedsJSON)
		}
	}

	// Inject discovery feed for per-page feed discovery links
	// If post has a sidebar feed, use that; otherwise use site default
	// Single-file builds publish no feeds, so skip feed discovery links and
	// site-level fallbacks that would point at missing pages.
	if singlePage, ok := config.Extra["single_page"].(bool); ok && singlePage {
		ctx.Set("single_page", true)
	} else if discoveryFeed := p.getDiscoveryFeed(post, sidebarFeed, m); discoveryFeed != nil {
		ctx.Set("discovery_feed", DiscoveryFeedToMap(discoveryFeed))
	}

	// Render the template
	html, err := p.engine.Render(templateName, ctx)
	if err != nil {
		return "", fmt.Errorf("failed to render template %q for post %q: %w", templateName, post.Path, err)
	}

	return html, nil
}

// isCanonicalRenderingFixture identifies the one explicit post used by the
// rendering contract. The slug alone is not sufficient: keep the authored
// title and template as part of the identity so an ordinary post cannot
// inherit the canonical projection accidentally.
func isCanonicalRenderingFixture(post *models.Post) bool {
	return post != nil && post.Slug == "test-headings" && post.PlainTitle() == "test headings" && post.Template == gardenNodeTypePost
}

// getFeedSidebarPosts returns the posts for the feed sidebar if the post belongs to a configured feed.
// Resolution priority:
//  1. Series – post has a "series" key in Extra
//  2. Explicit frontmatter – post.PrevNextFeed or post.Extra["sidebar_feed"]
//  3. Tag-based feeds – from config feed_sidebar.feeds
//  4. Auto-discovery – find the best (smallest non-excluded) feed containing this post,
//     preferring primary feeds only as a tie-breaker for equally specific candidates.
func (p *TemplatesPlugin) getFeedSidebarPosts(post *models.Post, config *lifecycle.Config, m *lifecycle.Manager) ([]*models.Post, *models.FeedConfig) {
	maxPosts := feedSidebarMaxPosts(config)

	// 1. Series
	seriesPosts, seriesFeed := p.getSeriesSidebarPosts(post, config, m)
	if seriesPosts != nil {
		return trimSidebarPosts(seriesPosts, post.Slug, maxPosts), seriesFeed
	}

	// 2. Explicit frontmatter: sidebar_feed or PrevNextFeed
	if explicitSlug := p.getExplicitFeedSlug(post); explicitSlug != "" {
		posts, fc := p.feedFromCachedConfigs(explicitSlug, post, m)
		if posts != nil {
			return trimSidebarPosts(posts, post.Slug, maxPosts), fc
		}
	}

	// 3. Tag-based feeds from feed_sidebar config
	tagPosts, tagFeed := p.getTagFeedSidebarPosts(post, config, m)
	if tagPosts != nil {
		return trimSidebarPosts(tagPosts, post.Slug, maxPosts), tagFeed
	}

	// 4. Auto-discovery: find the best feed containing this post
	autoPosts, autoFeed := p.autoDiscoverFeed(post, config, m)
	if autoPosts != nil {
		return trimSidebarPosts(autoPosts, post.Slug, maxPosts), autoFeed
	}

	return nil, nil
}

func trimSidebarPosts(posts []*models.Post, currentSlug string, maxPosts int) []*models.Post {
	if maxPosts <= 0 || len(posts) <= maxPosts {
		return posts
	}

	currentIndex := -1
	for i, p := range posts {
		if p != nil && p.Slug == currentSlug {
			currentIndex = i
			break
		}
	}
	if currentIndex == -1 {
		trimmed := make([]*models.Post, maxPosts)
		copy(trimmed, posts[:maxPosts])
		return trimmed
	}

	half := maxPosts / 2
	start := currentIndex - half
	if start < 0 {
		start = 0
	}
	end := start + maxPosts
	if end > len(posts) {
		end = len(posts)
		start = end - maxPosts
		if start < 0 {
			start = 0
		}
	}

	trimmed := make([]*models.Post, end-start)
	copy(trimmed, posts[start:end])
	return trimmed
}

func feedSidebarMaxPosts(config *lifecycle.Config) int {
	if config == nil {
		return 0
	}
	components, ok := config.Extra["components"].(models.ComponentsConfig)
	if !ok {
		return 0
	}
	return components.FeedSidebar.MaxPosts
}

// getExplicitFeedSlug returns a feed slug from post frontmatter, if any.
// Checks post.Extra["sidebar_feed"] first, then post.PrevNextFeed.
func (p *TemplatesPlugin) getExplicitFeedSlug(post *models.Post) string {
	if post == nil {
		return ""
	}
	if post.Extra != nil {
		if sf, ok := post.Extra["sidebar_feed"].(string); ok && sf != "" {
			return sf
		}
	}
	if post.PrevNextFeed != "" {
		return post.PrevNextFeed
	}
	return ""
}

// feedFromCachedConfigs looks up a feed by slug in cached feed_configs and
// returns its posts if the current post is a member.
func (p *TemplatesPlugin) feedFromCachedConfigs(slug string, post *models.Post, m *lifecycle.Manager) ([]*models.Post, *models.FeedConfig) {
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return nil, nil
	}
	configs, ok := cached.([]models.FeedConfig)
	if !ok {
		return nil, nil
	}
	fc := GetFeedBySlug(slug, configs)
	if fc == nil || len(fc.Posts) == 0 {
		return nil, nil
	}
	// Verify the post is actually in the feed
	found := false
	for _, fp := range fc.Posts {
		if fp.Slug == post.Slug {
			found = true
			break
		}
	}
	if !found {
		return nil, nil
	}
	return fc.Posts, fc
}

// getTagFeedSidebarPosts is the original tag-based feed sidebar logic extracted
// from the old getFeedSidebarPosts.
func (p *TemplatesPlugin) getTagFeedSidebarPosts(post *models.Post, config *lifecycle.Config, m *lifecycle.Manager) ([]*models.Post, *models.FeedConfig) {
	// Get components config
	components, ok := config.Extra["components"].(models.ComponentsConfig)
	if !ok {
		return nil, nil
	}

	// Check if feed sidebar is enabled
	if components.FeedSidebar.Enabled == nil || !*components.FeedSidebar.Enabled {
		return nil, nil
	}

	// Get configured feed slugs (e.g., ["tags/daily-note"])
	feedSlugs := components.FeedSidebar.Feeds
	if len(feedSlugs) == 0 {
		return nil, nil
	}

	// Check if this post belongs to any of the configured feeds
	// For tag-based feeds (tags/xxx), check if post has the tag
	for _, feedSlug := range feedSlugs {
		if !strings.HasPrefix(feedSlug, "tags/") {
			continue // Only handle tag feeds for now
		}

		// Extract tag name from feed slug (e.g., "tags/daily-note" -> "daily-note")
		tagName := strings.TrimPrefix(feedSlug, "tags/")

		// Check if post has this tag
		hasTag := false
		for _, postTag := range post.Tags {
			if postTag == tagName {
				hasTag = true
				break
			}
		}

		if !hasTag {
			continue
		}

		// Post belongs to this feed - collect all posts with this tag
		allPosts := m.Posts()
		sidebarPosts := make([]*models.Post, 0)
		for _, feedPost := range allPosts {
			for _, t := range feedPost.Tags {
				if t == tagName && feedPost.Published && !feedPost.Draft && !feedPost.Skip {
					sidebarPosts = append(sidebarPosts, feedPost)
					break
				}
			}
		}

		// Sort by date (newest first)
		sortPostsByDate(sidebarPosts, true)

		// Create a feed config for template display
		feedConfig := &models.FeedConfig{
			Slug:  feedSlug,
			Title: fmt.Sprintf("Posts tagged: %s", tagName),
			Posts: sidebarPosts,
		}
		feedConfig.ApplyDefaults(ToModelsConfig(config).FeedDefaults)

		return sidebarPosts, feedConfig
	}

	return nil, nil
}

type autoFeedCandidate struct {
	fc    *models.FeedConfig
	count int
}

func selectBestAutoFeedCandidate(candidates []autoFeedCandidate) autoFeedCandidate {
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.count < best.count || (c.count == best.count && c.fc.Primary && !best.fc.Primary) {
			best = c
		}
	}
	return best
}

// autoDiscoverFeed finds the best feed containing this post from cached feed_configs.
// It prefers smaller/more specific feeds over large catch-all feeds.
// Feeds with slugs like "archive", "all", "subscription-*" are deprioritized.
func (p *TemplatesPlugin) autoDiscoverFeed(post *models.Post, config *lifecycle.Config, m *lifecycle.Manager) ([]*models.Post, *models.FeedConfig) {
	if post == nil || post.Slug == "" {
		return nil, nil
	}

	// Check if feed sidebar is enabled (must be enabled for auto-discovery)
	components, ok := config.Extra["components"].(models.ComponentsConfig)
	if !ok {
		return nil, nil
	}
	if components.FeedSidebar.Enabled == nil || !*components.FeedSidebar.Enabled {
		return nil, nil
	}

	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return nil, nil
	}
	configs, ok := cached.([]models.FeedConfig)
	if !ok {
		return nil, nil
	}

	// Slugs to skip in auto-discovery (catch-all or meta feeds)
	skipSlugs := map[string]bool{
		"archive": true, "all": true, "sitemap": true, "search": true,
	}
	// Prefixes to skip
	skipPrefixes := []string{"subscription-", "tags/"}

	var candidates []autoFeedCandidate

	for i := range configs {
		fc := &configs[i]
		if skipSlugs[fc.Slug] {
			continue
		}
		skip := false
		for _, prefix := range skipPrefixes {
			if strings.HasPrefix(fc.Slug, prefix) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		// Check if this post is in the feed
		for _, fp := range fc.Posts {
			if fp.Slug == post.Slug {
				candidates = append(candidates, autoFeedCandidate{fc: fc, count: len(fc.Posts)})
				break
			}
		}
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	// Pick the smallest feed (most specific). If two candidates are equally
	// specific, prefer the primary feed so configured navigation still wins
	// when the scope is otherwise identical.
	best := selectBestAutoFeedCandidate(candidates)

	return best.fc.Posts, best.fc
}

func (p *TemplatesPlugin) getSeriesSidebarPosts(post *models.Post, config *lifecycle.Config, m *lifecycle.Manager) ([]*models.Post, *models.FeedConfig) {
	if post == nil || config == nil || m == nil {
		return nil, nil
	}

	seriesName := getStringFromExtra(post.Extra, seriesKey)
	if seriesName == "" {
		return nil, nil
	}

	seriesCfg := parseSeriesConfig(config)
	if !seriesCfg.AutoSidebar {
		return nil, nil
	}

	seriesSlug := buildSeriesFeedSlug(seriesCfg.SlugPrefix, slugify(seriesName))

	allPosts := m.Posts()
	seriesPosts := make([]*models.Post, 0)
	for _, feedPost := range allPosts {
		postSeries := getStringFromExtra(feedPost.Extra, seriesKey)
		if postSeries == "" {
			continue
		}
		postSeriesSlug := buildSeriesFeedSlug(seriesCfg.SlugPrefix, slugify(postSeries))
		if postSeriesSlug == seriesSlug {
			seriesPosts = append(seriesPosts, feedPost)
		}
	}

	if len(seriesPosts) == 0 {
		return nil, nil
	}

	seriesSlugValue := slugify(seriesName)
	group := &seriesGroup{
		name:  seriesName,
		slug:  seriesSlugValue,
		posts: seriesPosts,
		cfg:   resolveSeriesOverride(seriesCfg, seriesName, seriesSlugValue),
	}

	sortSeriesPosts(group, false)
	publishedPosts := filterSeriesOutputPosts(group.posts)
	if len(publishedPosts) == 0 {
		return nil, nil
	}

	feedConfig := &models.FeedConfig{
		Slug:  seriesSlug,
		Title: seriesDisplayTitle(group.name, group.cfg),
		Posts: publishedPosts,
	}
	if group.cfg != nil && group.cfg.Description != "" {
		feedConfig.Description = group.cfg.Description
	}
	feedConfig.ApplyDefaults(ToModelsConfig(config).FeedDefaults)

	return publishedPosts, feedConfig
}

// getSidebarPrevNext finds the previous and next posts relative to the current post
// within the sidebar posts list. The sidebar posts are sorted by date (newest first),
// so "prev" is the newer post (earlier in the list) and "next" is the older post.
func (p *TemplatesPlugin) getSidebarPrevNext(currentPost *models.Post, sidebarPosts []*models.Post) (prev, next *models.Post) {
	if len(sidebarPosts) == 0 {
		return nil, nil
	}

	// Find current post's position in sidebar posts
	position := -1
	for i, post := range sidebarPosts {
		if post.Slug == currentPost.Slug {
			position = i
			break
		}
	}

	if position == -1 {
		return nil, nil
	}

	// Since posts are sorted newest first:
	// - prev (newer) is at position-1
	// - next (older) is at position+1
	if position > 0 {
		prev = sidebarPosts[position-1]
	}
	if position < len(sidebarPosts)-1 {
		next = sidebarPosts[position+1]
	}

	return prev, next
}

// sidebarFeedJSON is the JSON structure for a single feed candidate
// used by the client-side feed cycling feature.
type sidebarFeedJSON struct {
	Slug                string            `json:"slug"`
	Title               string            `json:"title"`
	Priority            string            `json:"priority"`
	Primary             bool              `json:"primary,omitempty"`
	ContainsCurrentPost bool              `json:"containsCurrentPost,omitempty"`
	Variants            []sidebarVariant  `json:"variants,omitempty"`
	Calendar            bool              `json:"calendar"`
	Posts               []sidebarPostJSON `json:"posts"`
	TotalPosts          int               `json:"totalPosts"`
	ChunkBase           string            `json:"chunkBase,omitempty"`
	ChunkSize           int               `json:"chunkSize,omitempty"`
	WindowStart         int               `json:"windowStart,omitempty"`
	Prev                *sidebarPostJSON  `json:"prev,omitempty"`
	Next                *sidebarPostJSON  `json:"next,omitempty"`
}

type sidebarVariant struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Href  string `json:"href"`
}

type sidebarPostJSON struct {
	Slug   string `json:"slug"`
	Title  string `json:"title"`
	Href   string `json:"href"`
	Active bool   `json:"active,omitempty"`
}

type sidebarPostKey struct {
	slug     string
	title    string
	href     string
	feedSlug string
}

// sidebarFeedsDataJSON is the top-level JSON structure embedded in the page.
type sidebarFeedsDataJSON struct {
	Feeds             []sidebarFeedJSON `json:"feeds"`
	RotationFeedSlugs []string          `json:"rotationFeedSlugs"`
	CurrentFeedIndex  int               `json:"currentFeedIndex"`
	CurrentPostSlug   string            `json:"currentPostSlug"`
}

// getAllCandidateFeeds collects all feeds that contain this post across all
// 4 priority levels. Unlike getFeedSidebarPosts which returns the first match,
// this collects ALL matches so the user can cycle between them with {/}.
func (p *TemplatesPlugin) getAllCandidateFeeds(
	post *models.Post, config *lifecycle.Config, m *lifecycle.Manager, postFormats models.PostFormatsConfig,
) []sidebarFeedJSON {
	if post == nil {
		return nil
	}

	syndication := getSyndicationConfig(config)
	maxPosts := feedSidebarMaxPosts(config)
	seen := make(map[string]bool)
	var feeds []sidebarFeedJSON

	addFeed := func(fc *models.FeedConfig, posts []*models.Post, priority string) {
		if fc == nil || fc.IncludesPrivate() || seen[fc.Slug] {
			return
		}
		// Exclude feeds explicitly opted out of the sidebar picker
		if fc.Sidebar != nil && !*fc.Sidebar {
			return
		}
		seen[fc.Slug] = true
		entry := p.buildSidebarFeedEntry(post, fc, posts, priority, syndication, postFormats, maxPosts)
		entry.ChunkBase = sidebarChunkBase(fc, posts, m)
		if entry.ChunkBase != "" {
			entry.ChunkSize = sidebarNavChunkSize
		}
		feeds = append(feeds, entry)
	}

	// 1. Series
	seriesPosts, seriesFeed := p.getSeriesSidebarPosts(post, config, m)
	if seriesPosts != nil {
		addFeed(seriesFeed, seriesPosts, "series")
	}

	// 2. Explicit frontmatter
	if explicitSlug := p.getExplicitFeedSlug(post); explicitSlug != "" {
		explicitPosts, explicitFc := p.feedFromCachedConfigs(explicitSlug, post, m)
		if explicitPosts != nil {
			addFeed(explicitFc, explicitPosts, "explicit")
		}
	}

	// 3. Tag-based feeds -- check ALL tag feeds, not just the first match
	p.collectTagFeeds(post, config, m, addFeed)

	// 4. Auto-discovery -- include matching feeds, sorted by size (smallest first)
	p.collectAutoDiscoveredFeeds(post, config, m, seen, addFeed)

	return feeds
}

// buildSidebarFeedEntry constructs a sidebarFeedJSON for a single feed,
// windowing the post list around the current post for large feeds.
func (p *TemplatesPlugin) buildSidebarFeedEntry(
	currentPost *models.Post, fc *models.FeedConfig,
	posts []*models.Post, priority string, syndication models.SyndicationConfig, postFormats models.PostFormatsConfig,
	maxPosts int,
) sidebarFeedJSON {
	const legacyMaxWindowPosts = 50

	prev, next := p.getSidebarPrevNext(currentPost, posts)

	// Find current post position for windowing
	currentPos := -1
	for i, fp := range posts {
		if fp.Slug == currentPost.Slug {
			currentPos = i
			break
		}
	}

	// A configured max_posts is the shared contract for both the static sidebar
	// and feed-cycling JSON. Preserve the historical member-feed window when the
	// setting is unset, while configured limits also bound feeds that do not
	// contain the current post.
	windowedPosts := posts
	if maxPosts > 0 {
		windowedPosts = trimSidebarPosts(posts, currentPost.Slug, maxPosts)
	} else if len(posts) > legacyMaxWindowPosts && currentPos >= 0 {
		half := legacyMaxWindowPosts / 2
		start := currentPos - half
		end := currentPos + half + 1
		if start < 0 {
			end -= start
			start = 0
		}
		if end > len(posts) {
			start -= end - len(posts)
			end = len(posts)
		}
		if start < 0 {
			start = 0
		}
		windowedPosts = posts[start:end]
	}

	feedTitle := fc.Title
	if feedTitle == "" {
		feedTitle = strings.ReplaceAll(fc.Slug, "-", " ")
		feedTitle = cases.Title(language.English).String(feedTitle)
	}

	feed := sidebarFeedJSON{
		Slug:                fc.Slug,
		Title:               feedTitle,
		Priority:            priority,
		Primary:             fc.Primary,
		ContainsCurrentPost: currentPos >= 0,
		Variants:            buildSidebarVariants(fc, syndication, postFormats),
		Calendar:            fc.HasView(models.FeedViewCalendar),
		TotalPosts:          len(posts),
		Posts:               make([]sidebarPostJSON, 0, len(windowedPosts)),
	}

	if len(windowedPosts) > 0 {
		for i, fp := range posts {
			if fp.Slug == windowedPosts[0].Slug {
				feed.WindowStart = i
				break
			}
		}
	}
	for _, fp := range windowedPosts {
		feed.Posts = append(feed.Posts, p.postToSidebarJSON(fp, fp.Slug == currentPost.Slug, fc.Slug))
	}

	if prev != nil {
		pj := p.postToSidebarJSON(prev, false, fc.Slug)
		feed.Prev = &pj
	}
	if next != nil {
		nj := p.postToSidebarJSON(next, false, fc.Slug)
		feed.Next = &nj
	}

	return feed
}

func buildSidebarVariants(fc *models.FeedConfig, syndication models.SyndicationConfig, postFormats models.PostFormatsConfig) []sidebarVariant {
	if fc == nil {
		return nil
	}

	baseHref := "/"
	if fc.Slug != "" {
		baseHref = "/" + strings.Trim(fc.Slug, "/") + "/"
	}

	variants := make([]sidebarVariant, 0, 8)
	add := func(key, label, href string) {
		variants = append(variants, sidebarVariant{Key: key, Label: label, Href: href})
	}
	canonicalVariantHref := func(ext string) string {
		if fc.Slug == "" {
			return "/index." + ext
		}
		return "/" + strings.Trim(fc.Slug, "/") + "." + ext
	}

	if fc.Formats.Markdown && postFormats.Markdown {
		add("md", "md", canonicalVariantHref("md"))
	}
	if fc.Formats.JSON {
		add("json", "json", baseHref+"feed.json")
	}
	if shouldGenerateFeedArchive(fc, syndication) && fc.Formats.RSS {
		add("archive-rss", "archive-rss", feedArchiveURL(fc.Slug, "rss.xml"))
	}
	if fc.Formats.RSS {
		add("rss", "rss", baseHref+"rss.xml")
	}
	if fc.Formats.Atom {
		add("atom", "atom", baseHref+"atom.xml")
	}
	if fc.Formats.HTML && postFormats.IsHTMLEnabled() {
		add("html", "html", baseHref)
	}
	if simpleHTMLViewEnabled(fc) {
		add("simple", "simple", baseHref+"simple/")
	}
	if fc.Formats.Text && postFormats.Text {
		add("txt", "txt", canonicalVariantHref("txt"))
	}

	return variants
}

func appendFeedParamToHref(href, feedSlug string) string {
	if feedSlug == "" {
		return href
	}

	parsed, err := url.Parse(href)
	if err != nil {
		separator := "?"
		if strings.Contains(href, "?") {
			separator = "&"
		}
		return href + separator + "feed=" + url.QueryEscape(feedSlug)
	}

	query := parsed.Query()
	query.Set("feed", feedSlug)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}

// postToSidebarJSON converts a Post to a sidebarPostJSON.
func postToSidebarJSON(fp *models.Post, active bool, feedSlug string) sidebarPostJSON {
	return sidebarPostFromKey(sidebarKey(fp, feedSlug), active)
}

func sidebarKey(fp *models.Post, feedSlug string) sidebarPostKey {
	title := fp.PlainTitle()
	if title == "" {
		title = fp.Slug
	}
	return sidebarPostKey{slug: fp.Slug, title: title, href: fp.Href, feedSlug: feedSlug}
}

func sidebarPostFromKey(key sidebarPostKey, active bool) sidebarPostJSON {
	return sidebarPostJSON{
		Slug:   key.slug,
		Title:  key.title,
		Href:   appendFeedParamToHref(key.href, key.feedSlug),
		Active: active,
	}
}

func (p *TemplatesPlugin) postToSidebarJSON(fp *models.Post, active bool, feedSlug string) sidebarPostJSON {
	key := sidebarKey(fp, feedSlug)
	p.sidebarPostMu.RLock()
	projection, ok := p.sidebarPosts[key]
	p.sidebarPostMu.RUnlock()
	if !ok {
		projection = sidebarPostFromKey(key, false)
		p.sidebarPostMu.Lock()
		if p.sidebarPosts == nil {
			p.sidebarPosts = make(map[sidebarPostKey]sidebarPostJSON)
		}
		p.sidebarPosts[key] = projection
		p.sidebarPostMu.Unlock()
	}
	projection.Active = active
	return projection
}

func (p *TemplatesPlugin) resetSidebarPosts() {
	p.sidebarPostMu.Lock()
	p.sidebarPosts = nil
	p.sidebarPostMu.Unlock()
}

// collectTagFeeds finds all tag-based feeds from the sidebar config that
// contain the given post and calls addFeed for each.
func (p *TemplatesPlugin) collectTagFeeds(
	post *models.Post, config *lifecycle.Config, m *lifecycle.Manager,
	addFeed func(*models.FeedConfig, []*models.Post, string),
) {
	components, ok := config.Extra["components"].(models.ComponentsConfig)
	if !ok {
		return
	}
	if components.FeedSidebar.Enabled == nil || !*components.FeedSidebar.Enabled {
		return
	}

	for _, feedSlug := range components.FeedSidebar.Feeds {
		if !strings.HasPrefix(feedSlug, "tags/") {
			continue
		}
		tagName := strings.TrimPrefix(feedSlug, "tags/")
		if !postHasTag(post, tagName) {
			continue
		}
		tagPosts := filterPostsByTag(m.Posts(), tagName)
		sortPostsByDate(tagPosts, true)
		tagFc := &models.FeedConfig{
			Slug:  feedSlug,
			Title: fmt.Sprintf("Posts tagged: %s", tagName),
			Posts: tagPosts,
		}
		addFeed(tagFc, tagPosts, "tag")
	}
}

// postHasTag returns true if the post has the given tag.
func postHasTag(post *models.Post, tag string) bool {
	for _, t := range post.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// filterPostsByTag returns all published, non-draft, non-skip posts with the given tag.
func filterPostsByTag(allPosts []*models.Post, tag string) []*models.Post {
	var result []*models.Post
	for _, fp := range allPosts {
		if !fp.Published || fp.Draft || fp.Skip {
			continue
		}
		for _, t := range fp.Tags {
			if t == tag {
				result = append(result, fp)
				break
			}
		}
	}
	return result
}

// collectAutoDiscoveredFeeds finds feeds from the cache that contain the given
// post, sorted by size (smallest first), skipping catch-all feeds. Results are
// capped at maxAutoFeeds.
func (p *TemplatesPlugin) collectAutoDiscoveredFeeds(
	post *models.Post, config *lifecycle.Config, m *lifecycle.Manager,
	seen map[string]bool,
	addFeed func(*models.FeedConfig, []*models.Post, string),
) {
	const maxAutoFeedPosts = 5000
	const maxAutoFeeds = 10

	components, ok := config.Extra["components"].(models.ComponentsConfig)
	if !ok || components.FeedSidebar.Enabled == nil || !*components.FeedSidebar.Enabled {
		return
	}

	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return
	}
	configs, ok := cached.([]models.FeedConfig)
	if !ok {
		return
	}

	skipPrefixes := []string{"subscription-", "tags/"}

	// Feeds can be excluded from the picker via include_private = true
	// in their feed config. No hardcoded skip list — all feeds that
	// contain the current post are candidates.

	type autoCandidate struct {
		fc    *models.FeedConfig
		count int
	}
	var candidates []autoCandidate

	for i := range configs {
		fc := &configs[i]
		if seen[fc.Slug] || fc.IncludesPrivate() || len(fc.Posts) > maxAutoFeedPosts {
			continue
		}
		if hasAnyPrefix(fc.Slug, skipPrefixes) {
			continue
		}
		for _, fp := range fc.Posts {
			if fp.Slug == post.Slug {
				candidates = append(candidates, autoCandidate{fc: fc, count: len(fc.Posts)})
				break
			}
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].count < candidates[j].count
	})

	if len(candidates) > maxAutoFeeds {
		candidates = candidates[:maxAutoFeeds]
	}

	for _, ac := range candidates {
		addFeed(ac.fc, ac.fc.Posts, "auto")
	}
}

// hasAnyPrefix returns true if s starts with any of the given prefixes.
func hasAnyPrefix(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// buildSidebarFeedsJSON builds the JSON string for all candidate feeds,
// for embedding in the page for client-side feed cycling.
func (p *TemplatesPlugin) buildSidebarFeedsJSON(
	post *models.Post, config *lifecycle.Config, m *lifecycle.Manager,
	primaryFeed *models.FeedConfig,
) string {
	postFormats := resolvePostFormats(post, config)
	feeds := p.getAllCandidateFeeds(post, config, m, postFormats)
	feeds = filterPublicSidebarFeeds(feeds)
	seen := makeSidebarFeedSet(feeds)
	feeds = p.appendMissingPrimarySidebarFeeds(feeds, seen, post, config, m, postFormats)

	if len(feeds) == 0 {
		return ""
	}
	// Keep one-feed metadata: lazy navigation works even without feed cycling.

	rotationFeedSlugs := p.buildRotationFeedSlugs(m, seen)

	// Find the index of the primary (currently displayed) feed
	currentIndex := 0
	if primaryFeed != nil {
		for i := range feeds {
			if feeds[i].Slug == primaryFeed.Slug {
				currentIndex = i
				break
			}
		}
	}

	// Keep the current feed's full server-rendered fallback. For alternate
	// feeds with a verified public sidecar, send only metadata until selected.
	// Feeds without a safe sidecar retain their embedded post window.
	compactAlternateSidebarPosts(feeds, currentIndex)

	data := sidebarFeedsDataJSON{
		Feeds:             feeds,
		RotationFeedSlugs: rotationFeedSlugs,
		CurrentFeedIndex:  currentIndex,
		CurrentPostSlug:   post.Slug,
	}

	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}

// compactAlternateSidebarPosts trims duplicate alternate feed lists from
// every article. The visible/current feed and unsafe/synthetic feeds keep
// their full initial window so no-JS navigation stays usable.
func compactAlternateSidebarPosts(feeds []sidebarFeedJSON, currentIndex int) {
	for i := range feeds {
		if i == currentIndex || feeds[i].ChunkBase == "" || feeds[i].ChunkSize <= 0 || feeds[i].TotalPosts == 0 {
			continue
		}
		feeds[i].Posts = []sidebarPostJSON{}
	}
}

// sidebarChunkBase only advertises a sidecar when its emitted public entries
// match the navigation feed exactly. Synthetic tag/series feeds can have a
// different order or may not be published at all; retain their HTML fallback.
func sidebarChunkBase(fc *models.FeedConfig, posts []*models.Post, m *lifecycle.Manager) string {
	if !sidebarChunksEnabled(fc) || m == nil {
		return ""
	}
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return ""
	}
	configs, ok := cached.([]models.FeedConfig)
	if !ok {
		return ""
	}
	for i := range configs {
		candidate := &configs[i]
		if candidate.Slug != fc.Slug || !sidebarChunksEnabled(candidate) {
			continue
		}
		published := publicSidebarChunkPosts(candidate)
		if len(published) != len(posts) || len(posts) == 0 {
			return ""
		}
		for j := range posts {
			if posts[j] == nil || published[j] == nil {
				return ""
			}
			// Most feeds share the same immutable Post pointers. Avoid
			// recomputing display titles for these common matches.
			if posts[j] == published[j] {
				continue
			}
			// The slug alone is insufficient: synthetic/sidebar copies can
			// differ from what the publisher emitted at the same position.
			// Never discard embedded navigation in favor of stale URLs or
			// titles from a chunk that no longer represents this feed.
			if sidebarKey(posts[j], fc.Slug) != sidebarKey(published[j], candidate.Slug) {
				return ""
			}
		}
		if fc.Slug == "" {
			return "/_sidebar/"
		}
		return "/" + strings.Trim(fc.Slug, "/") + "/_sidebar/"
	}
	return ""
}

func filterPublicSidebarFeeds(feeds []sidebarFeedJSON) []sidebarFeedJSON {
	publicFeeds := feeds[:0]
	for i := range feeds {
		if feeds[i].Priority == feedPriorityPrivate {
			continue
		}
		publicFeeds = append(publicFeeds, feeds[i])
	}
	return publicFeeds
}

func makeSidebarFeedSet(feeds []sidebarFeedJSON) map[string]bool {
	seen := make(map[string]bool, len(feeds))
	for i := range feeds {
		seen[feeds[i].Slug] = true
	}
	return seen
}

func (p *TemplatesPlugin) appendMissingPrimarySidebarFeeds(
	feeds []sidebarFeedJSON, seen map[string]bool, post *models.Post,
	config *lifecycle.Config, m *lifecycle.Manager, postFormats models.PostFormatsConfig,
) []sidebarFeedJSON {
	maxPosts := feedSidebarMaxPosts(config)
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return feeds
	}
	configs, ok := cached.([]models.FeedConfig)
	if !ok {
		return feeds
	}

	syndication := getSyndicationConfig(config)
	for i := range configs {
		fc := &configs[i]
		if seen[fc.Slug] || fc.IncludesPrivate() {
			continue
		}
		if fc.Sidebar != nil && !*fc.Sidebar {
			continue
		}
		// Include if explicitly sidebar=true OR if primary (and not sidebar=false)
		sidebarExplicit := fc.Sidebar != nil && *fc.Sidebar
		if !sidebarExplicit && !fc.Primary {
			continue
		}
		entry := p.buildSidebarFeedEntry(post, fc, fc.Posts, "primary", syndication, postFormats, maxPosts)
		entry.ChunkBase = sidebarChunkBase(fc, fc.Posts, m)
		if entry.ChunkBase != "" {
			entry.ChunkSize = sidebarNavChunkSize
		}
		feeds = append(feeds, entry)
		seen[fc.Slug] = true
	}

	return feeds
}

func (p *TemplatesPlugin) buildRotationFeedSlugs(m *lifecycle.Manager, seen map[string]bool) []string {
	rotationFeedSlugs := make([]string, 0, len(seen))
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return rotationFeedSlugs
	}
	configs, ok := cached.([]models.FeedConfig)
	if !ok {
		return rotationFeedSlugs
	}

	// First pass: collect primary or sidebar=true feeds that appear in the sidebar candidate set
	for i := range configs {
		fc := &configs[i]
		if fc.IncludesPrivate() || !seen[fc.Slug] {
			continue
		}
		sidebarExplicit := fc.Sidebar != nil && *fc.Sidebar
		if !sidebarExplicit && !fc.Primary {
			continue
		}
		rotationFeedSlugs = append(rotationFeedSlugs, fc.Slug)
	}

	// Fallback: if no primary feeds matched, include all sidebar candidate feeds
	// so users can still cycle through available feeds.
	// Use the seen set directly since it represents all candidate feeds for
	// this post, including auto-generated tag feeds not in feed_configs.
	if len(rotationFeedSlugs) == 0 {
		for slug := range seen {
			rotationFeedSlugs = append(rotationFeedSlugs, slug)
		}
		// Sort for deterministic order
		sort.Strings(rotationFeedSlugs)
	}

	return rotationFeedSlugs
}

// getDiscoveryFeed returns the discovery feed for a post.
// If the post has a sidebar feed, that feed is used for discovery.
// Otherwise, the site default feed (root subscription feed) is used.
func (p *TemplatesPlugin) getDiscoveryFeed(post *models.Post, sidebarFeed *models.FeedConfig, m *lifecycle.Manager) *DiscoveryFeed {
	// Get feed configs from cache
	var feedConfigs []models.FeedConfig
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			feedConfigs = fcs
		}
	}

	return GetDiscoveryFeed(post, sidebarFeed, feedConfigs)
}

// sortPostsByDate sorts posts by date.
// If reverse is true, sorts newest first.
func sortPostsByDate(posts []*models.Post, reverse bool) {
	sort.Slice(posts, func(i, j int) bool {
		// Handle nil dates
		if posts[i].Date == nil && posts[j].Date == nil {
			return false
		}
		if posts[i].Date == nil {
			return !reverse
		}
		if posts[j].Date == nil {
			return reverse
		}
		if reverse {
			return posts[i].Date.After(*posts[j].Date)
		}
		return posts[i].Date.Before(*posts[j].Date)
	})
}

// collectPrivatePaths returns a list of paths (hrefs) for all private posts.
// These paths are used in robots.txt templates to add Disallow directives.
// Includes all format variants (.txt, .ansi, .md, .og) and excludes the robots post itself.
func collectPrivatePaths(posts []*models.Post) []string {
	var paths []string
	for _, post := range posts {
		if post.Private && !post.Draft && !post.Skip {
			// Skip the robots post itself to avoid self-reference
			if post.Slug == "robots" {
				continue
			}
			// Add base href (e.g., /slug/)
			paths = append(paths, post.Href)
			// Add format variants
			// For regular posts: /slug.txt, /slug.ansi, /slug.md, /slug.og/
			if post.Slug != "" {
				paths = append(paths,
					"/"+post.Slug+".txt",
					"/"+post.Slug+".ansi",
					"/"+post.Slug+".md",
					"/"+post.Slug+".og/",
				)
			}
		}
	}
	return paths
}

// Priority returns the plugin priority for the given stage.
// Templates should run late in the render stage, after markdown rendering.
func (p *TemplatesPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate // Run after markdown rendering
	}
	return lifecycle.PriorityDefault
}

// Engine returns the template engine for use by other plugins.
func (p *TemplatesPlugin) Engine() *templates.Engine {
	return p.engine
}

// ToModelsConfig converts lifecycle.Config to models.Config for template context.
// This is exported for use by other plugins that need to build template contexts
// with the full config (e.g., publish_feeds, blogroll).
// Note: This is not cached because lifecycle.Config.Extra is mutable and may be
// modified by plugins throughout the build process (e.g., image_zoom sets glightbox_enabled).
func ToModelsConfig(config *lifecycle.Config) *models.Config {
	return toModelsConfigUncached(config)
}

// toModelsConfigUncached is the actual conversion logic, used by the cache.
func toModelsConfigUncached(config *lifecycle.Config) *models.Config {
	if config == nil {
		return nil
	}
	// Convert lifecycle.Config to models.Config
	modelsConfig := &models.Config{
		OutputDir:    config.OutputDir,
		Title:        getStringFromExtra(config.Extra, "title"),
		URL:          getStringFromExtra(config.Extra, "url"),
		Description:  getStringFromExtra(config.Extra, "description"),
		Author:       getStringFromExtra(config.Extra, "author"),
		TemplatesDir: getStringFromExtra(config.Extra, "templates_dir"),
		Templates:    models.NewTemplatesConfig(),
	}

	// Copy nav items if available
	if navItems, ok := config.Extra["nav"].([]models.NavItem); ok {
		modelsConfig.Nav = navItems
	}

	// Copy footer config if available
	if footer, ok := config.Extra["footer"].(models.FooterConfig); ok {
		modelsConfig.Footer = footer
	}

	// Copy layout config if available
	switch layoutVal := config.Extra["layout"].(type) {
	case *models.LayoutConfig:
		modelsConfig.Layout = *layoutVal
	case models.LayoutConfig:
		modelsConfig.Layout = layoutVal
	}

	// Copy sidebar config if available
	if sidebar, ok := config.Extra["sidebar"].(models.SidebarConfig); ok {
		modelsConfig.Sidebar = sidebar
	}

	// Copy toc config if available
	if toc, ok := config.Extra["toc"].(models.TocConfig); ok {
		modelsConfig.Toc = toc
	}

	// Copy header config if available
	if header, ok := config.Extra["header"].(models.HeaderLayoutConfig); ok {
		modelsConfig.Header = header
	}

	// Copy feeds page config if available
	modelsConfig.FeedsPage = feedsPageConfigFromExtra(config.Extra["feeds_page"])

	// Copy SEO config if available
	switch seoVal := config.Extra["seo"].(type) {
	case models.SEOConfig:
		modelsConfig.SEO = seoVal
	case map[string]interface{}:
		modelsConfig.SEO = models.SEOConfig{
			TwitterHandle: getStringFromMap(seoVal, "twitter_handle"),
			DefaultImage:  getStringFromMap(seoVal, "default_image"),
			LogoURL:       getStringFromMap(seoVal, "logo_url"),
			AuthorImage:   getStringFromMap(seoVal, "author_image"),
		}
	}

	// Copy Search config if available, otherwise use defaults
	// This ensures search is enabled by default with position "navbar"
	if search, ok := config.Extra["search"].(models.SearchConfig); ok {
		modelsConfig.Search = search
	} else {
		modelsConfig.Search = models.NewSearchConfig()
	}

	if viewTransitions, ok := config.Extra["view_transitions"].(models.ViewTransitionsConfig); ok {
		modelsConfig.ViewTransitions = viewTransitions
	} else {
		modelsConfig.ViewTransitions = models.NewViewTransitionsConfig()
	}

	// Copy remaining plugin configs
	copyPluginConfigs(config, modelsConfig)

	// Copy the entire Extra map so templates can access dynamic plugin config
	// (e.g., glightbox_enabled, glightbox_options set by image_zoom plugin)
	if config.Extra != nil {
		modelsConfig.Extra = make(map[string]any)
		for k, v := range config.Extra {
			modelsConfig.Extra[k] = v
		}
	}

	return modelsConfig
}

// copyPluginConfigs copies plugin-specific config sections from lifecycle.Config to models.Config.
func copyPluginConfigs(config *lifecycle.Config, modelsConfig *models.Config) {
	if templatesCfg, ok := config.Extra["templates"].(models.TemplatesConfig); ok {
		modelsConfig.Templates = templatesCfg
	} else if templatesCfgPtr, ok := config.Extra["templates"].(*models.TemplatesConfig); ok && templatesCfgPtr != nil {
		modelsConfig.Templates = *templatesCfgPtr
	}

	// Copy Components config if available
	if components, ok := config.Extra["components"].(models.ComponentsConfig); ok {
		modelsConfig.Components = components
	}

	// Copy PostFormats config if available
	if postFormats, ok := config.Extra["post_formats"].(models.PostFormatsConfig); ok {
		modelsConfig.PostFormats = postFormats
	}

	// Copy WebSub config if available
	if websub, ok := config.Extra["websub"].(models.WebSubConfig); ok {
		modelsConfig.WebSub = websub
	}

	// Copy Theme config if available
	if theme, ok := config.Extra["theme"].(models.ThemeConfig); ok {
		modelsConfig.Theme = theme
	}
	if theme, ok := config.Extra["theme"].(map[string]interface{}); ok {
		modelsConfig.Theme = themeFromMap(theme, modelsConfig.Theme)
	}
	// Copy Head config if available
	if head, ok := config.Extra["head"].(models.HeadConfig); ok {
		modelsConfig.Head = head
	}
	if headMap, ok := config.Extra["head"].(map[string]interface{}); ok {
		modelsConfig.Head = headFromMap(headMap, modelsConfig.Head)
	}

	// Copy Assets config if available
	if assetsConfig, ok := config.Extra["assets"].(models.AssetsConfig); ok {
		modelsConfig.Assets = assetsConfig
	}
	if assetsConfig, ok := config.Extra["assets"].(*models.AssetsConfig); ok && assetsConfig != nil {
		modelsConfig.Assets = *assetsConfig
	}

	// Copy Tags config if available
	if tags, ok := config.Extra["tags"].(models.TagsConfig); ok {
		modelsConfig.Tags = tags
	} else {
		modelsConfig.Tags = models.NewTagsConfig()
	}

	// Copy Garden config if available
	if garden, ok := config.Extra["garden"].(models.GardenConfig); ok {
		modelsConfig.Garden = garden
	} else {
		modelsConfig.Garden = models.NewGardenConfig()
	}

	// Copy image library config if available
	if images, ok := config.Extra["images"].(models.ImagesConfig); ok {
		modelsConfig.Images = images
	} else {
		modelsConfig.Images = models.NewImagesConfig()
	}
}

// getStringFromExtra safely gets a string value from the Extra map.
func getStringFromExtra(extra map[string]interface{}, key string) string {
	if extra == nil {
		return ""
	}
	if v, ok := extra[key].(string); ok {
		return v
	}
	return ""
}

func themeFromMap(theme map[string]interface{}, fallback models.ThemeConfig) models.ThemeConfig {
	result := fallback
	if theme == nil {
		return result
	}
	if name, ok := theme["name"].(string); ok && name != "" {
		result.Name = name
	}
	if palette, ok := theme["palette"].(string); ok && palette != "" {
		result.Palette = palette
	}
	if paletteLight, ok := theme["palette_light"].(string); ok && paletteLight != "" {
		result.PaletteLight = paletteLight
	}
	if paletteDark, ok := theme["palette_dark"].(string); ok && paletteDark != "" {
		result.PaletteDark = paletteDark
	}
	if customCSS, ok := theme["custom_css"].(string); ok {
		result.CustomCSS = customCSS
	}
	switch variables := theme["variables"].(type) {
	case map[string]string:
		result.Variables = variables
	case map[string]interface{}:
		result.Variables = make(map[string]string, len(variables))
		for key, value := range variables {
			if v, ok := value.(string); ok {
				result.Variables[key] = v
			}
		}
	}
	return result
}

func headFromMap(head map[string]interface{}, fallback models.HeadConfig) models.HeadConfig {
	result := fallback
	if head == nil {
		return result
	}
	if text, ok := head["text"].(string); ok {
		result.Text = text
	}
	result.Meta = metaTagsFromValue(head["meta"], result.Meta)
	result.Link = linkTagsFromValue(head["link"], result.Link)
	result.Script = scriptTagsFromValue(head["script"], result.Script)
	return result
}

func metaTagsFromValue(value interface{}, fallback []models.MetaTag) []models.MetaTag {
	switch meta := value.(type) {
	case []models.MetaTag:
		return meta
	case []interface{}:
		result := make([]models.MetaTag, 0, len(meta))
		for _, entry := range meta {
			m, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			tag := models.MetaTag{}
			if name, ok := m["name"].(string); ok {
				tag.Name = name
			}
			if property, ok := m["property"].(string); ok {
				tag.Property = property
			}
			if content, ok := m["content"].(string); ok {
				tag.Content = content
			}
			result = append(result, tag)
		}
		return result
	default:
		return fallback
	}
}

func linkTagsFromValue(value interface{}, fallback []models.LinkTag) []models.LinkTag {
	switch links := value.(type) {
	case []models.LinkTag:
		return links
	case []interface{}:
		result := make([]models.LinkTag, 0, len(links))
		for _, entry := range links {
			m, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			tag := models.LinkTag{}
			if rel, ok := m["rel"].(string); ok {
				tag.Rel = rel
			}
			if href, ok := m["href"].(string); ok {
				tag.Href = href
			}
			if crossorigin, ok := m["crossorigin"].(bool); ok {
				tag.Crossorigin = crossorigin
			}
			result = append(result, tag)
		}
		return result
	default:
		return fallback
	}
}

func scriptTagsFromValue(value interface{}, fallback []models.ScriptTag) []models.ScriptTag {
	switch scripts := value.(type) {
	case []models.ScriptTag:
		return scripts
	case []interface{}:
		result := make([]models.ScriptTag, 0, len(scripts))
		for _, entry := range scripts {
			m, ok := entry.(map[string]interface{})
			if !ok {
				continue
			}
			tag := models.ScriptTag{}
			if src, ok := m["src"].(string); ok {
				tag.Src = src
			}
			result = append(result, tag)
		}
		return result
	default:
		return fallback
	}
}

// getStringFromMap safely gets a string value from a map.
func getStringFromMap(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func feedsPageConfigFromExtra(raw interface{}) models.FeedsPageConfig {
	defaults := models.NewFeedsPageConfig()

	switch feedsPageVal := raw.(type) {
	case models.FeedsPageConfig:
		return feedsPageVal
	case map[string]interface{}:
		feedsPage := models.FeedsPageConfig{
			Enabled:     defaults.Enabled,
			Title:       getStringFromMap(feedsPageVal, "title"),
			Description: getStringFromMap(feedsPageVal, "description"),
			Template:    getStringFromMap(feedsPageVal, "template"),
			SlugPrefix:  getStringFromMap(feedsPageVal, "slug_prefix"),
			Robots:      getStringFromMap(feedsPageVal, "robots"),
		}
		if rawShowPrivate, ok := feedsPageVal["show_private_feeds"].([]interface{}); ok {
			feedsPage.ShowPrivateFeeds = make([]string, 0, len(rawShowPrivate))
			for _, value := range rawShowPrivate {
				if s, ok := value.(string); ok {
					feedsPage.ShowPrivateFeeds = append(feedsPage.ShowPrivateFeeds, s)
				}
			}
		}
		if rawShowPrivate, ok := feedsPageVal["show_private_feeds"].([]string); ok {
			feedsPage.ShowPrivateFeeds = append([]string{}, rawShowPrivate...)
		}
		if enabled, ok := feedsPageVal["enabled"].(bool); ok {
			feedsPage.Enabled = &enabled
		}
		if feedsPage.Title == "" {
			feedsPage.Title = defaults.Title
		}
		if feedsPage.Description == "" {
			feedsPage.Description = defaults.Description
		}
		if feedsPage.Template == "" {
			feedsPage.Template = defaults.Template
		}
		if feedsPage.SlugPrefix == "" {
			feedsPage.SlugPrefix = defaults.SlugPrefix
		}
		return feedsPage
	default:
		return defaults
	}
}
