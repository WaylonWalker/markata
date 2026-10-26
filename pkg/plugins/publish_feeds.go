// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
	"github.com/WaylonWalker/markata/pkg/themes"
)

var publishFeedsLog = logging.Component("publish_feeds").Phase("write")

// sidebarNavChunkSize bounds the bytes fetched and DOM work when browsing long sidebars.
const sidebarNavChunkSize = 64

// PublishFeedsPlugin writes feeds to multiple output formats during the write stage.
// It also registers synthetic posts in the Configure stage so they can be resolved by wikilinks.
type PublishFeedsPlugin struct {
	navPreviews map[string]map[string]interface{}
	// engineCache caches template engines to avoid re-parsing templates for each feed
	engineMu    sync.RWMutex
	engineCache map[string]*templates.Engine

	// postFeedHashCache memoizes per-post feed semantic hashes for the current build.
	postFeedHashMu    sync.RWMutex
	postFeedHashCache map[string]string
}

// NewPublishFeedsPlugin creates a new PublishFeedsPlugin.
func NewPublishFeedsPlugin() *PublishFeedsPlugin {
	return &PublishFeedsPlugin{
		engineCache:       make(map[string]*templates.Engine),
		postFeedHashCache: make(map[string]string),
	}
}

// Name returns the unique name of the plugin.
func (p *PublishFeedsPlugin) Name() string {
	return "publish_feeds"
}

// getOrCreateEngine returns a cached template engine, or creates one if not cached.
func (p *PublishFeedsPlugin) getOrCreateEngine(templatesDir, themeName string) (*templates.Engine, error) {
	cacheKey := templatesDir + ":" + themeName

	// Fast path: check cache with read lock
	p.engineMu.RLock()
	if engine, ok := p.engineCache[cacheKey]; ok {
		p.engineMu.RUnlock()
		return engine, nil
	}
	p.engineMu.RUnlock()

	// Slow path: create engine with write lock
	p.engineMu.Lock()
	defer p.engineMu.Unlock()

	// Double-check after acquiring write lock
	if engine, ok := p.engineCache[cacheKey]; ok {
		return engine, nil
	}

	engine, err := templates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return nil, err
	}

	p.engineCache[cacheKey] = engine
	return engine, nil
}

// Configure registers synthetic posts for feed pages so they can be resolved by wikilinks.
// These posts are marked with Skip: true so they don't interfere with normal rendering.
func (p *PublishFeedsPlugin) Configure(m *lifecycle.Manager) error {
	// Get feed configs from cache (set by FeedsPlugin)
	var feedConfigs []models.FeedConfig
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			feedConfigs = fcs
		}
	}

	if len(feedConfigs) == 0 {
		return nil
	}

	// Helper to create string pointer
	strPtr := func(s string) *string { return &s }

	// Register synthetic post for each feed
	for i := range feedConfigs {
		fc := &feedConfigs[i]

		// Determine title
		title := fc.Title
		if title == "" {
			title = fc.Slug
		}

		// Create synthetic post
		feedPost := &models.Post{
			Slug:        fc.Slug,
			Title:       strPtr(title),
			Description: strPtr(fc.Description),
			Href:        "/" + fc.Slug + "/",
			Published:   true,
			Skip:        true,
		}
		m.AddPost(feedPost)
	}

	return nil
}

// Write generates and writes feed files in all configured formats.
// Feed generation is parallelized for better performance with many feeds.
// Uses incremental build cache to skip feeds with unchanged content.
func (p *PublishFeedsPlugin) Write(m *lifecycle.Manager) error {
	p.resetPostFeedHashCache()
	p.navPreviews = nil
	if previews, ok := m.Cache().Get("nav_previews"); ok {
		if resolved, ok := previews.(map[string]map[string]interface{}); ok {
			p.navPreviews = resolved
		}
	}

	config := m.Config()
	feedConfigs := getCachedFeedConfigs(m)
	if err := validateShortsRoutes(feedConfigs, m.Posts()); err != nil {
		return err
	}

	// Reconcile routes even if the last Shorts-enabled feed was removed.
	// Never remove stale public output until all current feeds published.
	if len(feedConfigs) != 0 {
		var err error
		if shouldPublishFeedsAsync(m) {
			// Serve workers must finish before stale routes are removed.
			err = p.publishFeedsAsync(m, feedConfigs)
		} else {
			err = p.publishFeeds(m, config, feedConfigs)
		}
		if err != nil {
			return err
		}
	}
	return reconcileShortsRoutes(config.OutputDir, feedConfigs, m.Posts())
}

func getCachedFeedConfigs(m *lifecycle.Manager) []models.FeedConfig {
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			return fcs
		}
	}

	return nil
}

func shouldPublishFeedsAsync(m *lifecycle.Manager) bool {
	if !lifecycle.IsServeFastMode(m) {
		return false
	}

	extra := m.Config().Extra
	if extra == nil {
		return false
	}

	async, ok := extra["feeds_async"].(bool)
	return ok && async
}

func (p *PublishFeedsPlugin) publishFeeds(m *lifecycle.Manager, config *lifecycle.Config, feedConfigs []models.FeedConfig) error {
	outputDir := config.OutputDir
	// Copy XSL stylesheets to output directory for styled RSS/Atom feeds
	if err := p.copyXSLStylesheets(config, outputDir); err != nil {
		return fmt.Errorf("copying XSL stylesheets: %w", err)
	}

	// Get build cache for incremental builds
	buildCache := GetBuildCache(m)
	affectedPaths := lifecycle.GetServeAffectedPaths(m)
	// Track skipped feeds
	var skippedCount int
	var rebuiltCount int

	// Process feeds concurrently with a worker pool
	// Use lifecycle concurrency (auto-detected from CPU cores, capped at 16)
	maxConcurrency := m.Concurrency()
	if maxConcurrency < 8 {
		maxConcurrency = 8
	}
	numFeeds := len(feedConfigs)

	// For small numbers of feeds, just process sequentially
	if numFeeds <= 2 {
		for i := range feedConfigs {
			fc := &feedConfigs[i]
			skip, hash := p.shouldSkipFeedWithConfigAndChanges(fc, buildCache, outputDir, config, affectedPaths)
			if skip {
				skippedCount++
				recordFeedOutput(m, fc)
				continue
			}
			if err := p.publishFeed(fc, config, outputDir); err != nil {
				recordFeedWriteError(m, fc)
				return fmt.Errorf("publishing feed %q: %w", fc.Slug, err)
			}
			recordFeedOutput(m, fc)
			p.cacheFeedHash(fc, buildCache, hash)
			rebuiltCount++
		}
		return nil
	}

	// Use a semaphore to limit concurrency
	semaphore := make(chan struct{}, maxConcurrency)
	errChan := make(chan error, numFeeds)
	var wg sync.WaitGroup
	var countMu sync.Mutex
	for i := range feedConfigs {
		fc := &feedConfigs[i]
		skip, hash := p.shouldSkipFeedWithConfigAndChanges(fc, buildCache, outputDir, config, affectedPaths)
		if skip {
			skippedCount++
			recordFeedOutput(m, fc)
			continue
		}
		wg.Add(1)
		go func(fc *models.FeedConfig, hash string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			if err := p.publishFeed(fc, config, outputDir); err != nil {
				recordFeedWriteError(m, fc)
				errChan <- fmt.Errorf("publishing feed %q: %w", fc.Slug, err)
				return
			}
			recordFeedOutput(m, fc)
			p.cacheFeedHash(fc, buildCache, hash)
			countMu.Lock()
			rebuiltCount++
			countMu.Unlock()
		}(fc, hash)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Log incremental stats if any feeds were skipped
	if skippedCount > 0 {
		publishFeedsLog.Printf("Incremental: %d feeds skipped, %d rebuilt", skippedCount, rebuiltCount)
	}

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PublishFeedsPlugin) publishFeedsAsync(m *lifecycle.Manager, feedConfigs []models.FeedConfig) error {
	config := m.Config()
	outputDir := config.OutputDir
	if len(feedConfigs) == 0 {
		return nil
	}

	if err := p.copyXSLStylesheets(config, outputDir); err != nil {
		return fmt.Errorf("copying XSL stylesheets: %w", err)
	}

	buildCache := GetBuildCache(m)
	affectedPaths := lifecycle.GetServeAffectedPaths(m)
	var skippedCount int
	var rebuiltCount int
	var countMu sync.Mutex

	maxConcurrency := m.Concurrency()
	if maxConcurrency < 8 {
		maxConcurrency = 8
	}
	numFeeds := len(feedConfigs)
	if numFeeds <= 2 {
		for i := range feedConfigs {
			fc := &feedConfigs[i]
			skip, hash := p.shouldSkipFeedWithConfigAndChanges(fc, buildCache, outputDir, config, affectedPaths)
			if skip {
				skippedCount++
				recordFeedOutput(m, fc)
				continue
			}
			if err := p.publishFeed(fc, config, outputDir); err != nil {
				recordFeedWriteError(m, fc)
				return fmt.Errorf("publishing feed %q: %w", fc.Slug, err)
			}
			recordFeedOutput(m, fc)
			rebuiltCount++
			p.cacheFeedHash(fc, buildCache, hash)
		}
	} else {
		sem := make(chan struct{}, maxConcurrency)
		errChan := make(chan error, numFeeds)
		var wg sync.WaitGroup

		for i := range feedConfigs {
			fc := &feedConfigs[i]
			wg.Add(1)
			go func(fc *models.FeedConfig) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				skip, hash := p.shouldSkipFeedWithConfigAndChanges(fc, buildCache, outputDir, config, affectedPaths)
				if skip {
					countMu.Lock()
					skippedCount++
					countMu.Unlock()
					recordFeedOutput(m, fc)
					return
				}
				if err := p.publishFeed(fc, config, outputDir); err != nil {
					recordFeedWriteError(m, fc)
					errChan <- fmt.Errorf("publishing feed %q: %w", fc.Slug, err)
					return
				}
				recordFeedOutput(m, fc)
				countMu.Lock()
				rebuiltCount++
				countMu.Unlock()
				p.cacheFeedHash(fc, buildCache, hash)
			}(fc)
		}

		wg.Wait()
		close(errChan)
		for err := range errChan {
			if err != nil {
				return err
			}
		}
	}

	if skippedCount > 0 || rebuiltCount > 0 {
		publishFeedsLog.Printf("Async: %d feeds skipped, %d rebuilt", skippedCount, rebuiltCount)
	}

	return nil
}

// computeFeedHash computes a hash of the feed's content and configuration.
// It includes post slugs (in feed order to detect sort changes), and all
// configuration fields that affect feed output.
func (p *PublishFeedsPlugin) computeFeedHash(fc *models.FeedConfig) string {
	return p.computeFeedHashWithConfigAndCache(fc, nil, nil)
}

func (p *PublishFeedsPlugin) computeFeedHashWithConfigAndCache(fc *models.FeedConfig, config *lifecycle.Config, cache *buildcache.Cache) string {
	h := sha256.New()
	syndication := getSyndicationConfig(config)
	writeStringField := func(value string) {
		fmt.Fprintf(h, "%d:", len(value))
		h.Write([]byte(value))
		h.Write([]byte{0})
	}
	writeBoolField := func(value bool) {
		fmt.Fprintf(h, "%t", value)
		h.Write([]byte{0})
	}
	writeIntField := func(value int) {
		fmt.Fprintf(h, "%d", value)
		h.Write([]byte{0})
	}

	// Hash post slugs in their current (sorted) order, not alphabetically.
	// This ensures sort order changes are detected.
	for _, post := range fc.Posts {
		writeStringField(post.Slug)
		writeStringField(p.getPostFeedItemHash(post, cache))
		if fc.HasView(models.FeedViewShorts) && post != nil {
			for _, key := range []string{"video", "cover", "cover_image", "og_image", "image_alt", "alt", "caption", "poster", "poster_image", "video_poster", "video_thumbnail", "thumbnail", "thumb"} {
				if value, ok := post.Extra[key].(string); ok {
					writeStringField(value)
				}
			}
		}
	}

	for i := range fc.Pages {
		page := &fc.Pages[i]
		writeIntField(page.Number)
		writeBoolField(page.HasPrev)
		writeBoolField(page.HasNext)
		writeStringField(page.PrevURL)
		writeStringField(page.NextURL)
		writeIntField(page.TotalPages)
		writeIntField(page.TotalItems)
		writeIntField(page.ItemsPerPage)
		for _, url := range page.PageURLs {
			writeStringField(url)
		}
	}

	// Hash all feed config fields that affect output
	writeStringField(fc.Slug)
	writeStringField(fc.Title)
	writeStringField(fc.Description)
	writeStringField(string(fc.Type))
	writeStringField(fc.Filter)
	writeStringField(fc.Sort)
	writeBoolField(fc.Reverse)
	writeBoolField(fc.Primary)
	writeIntField(fc.ItemsPerPage)
	writeIntField(fc.OrphanThreshold)
	writeStringField(string(fc.PaginationType))
	writeBoolField(fc.IncludesPrivate())
	writeBoolField(fc.ArchiveDisabled)
	writeIntField(syndication.MaxItems)
	writeBoolField(syndication.IncludeContent)
	writeBoolField(syndication.SiteArchiveDisabled)
	writeBoolField(syndication.FeedArchivesDisabled)

	// Hash sidebar config
	writeBoolField(fc.Sidebar != nil && *fc.Sidebar)
	writeStringField(fc.SidebarTitle)
	writeIntField(fc.SidebarOrder)
	writeStringField(fc.SidebarGroupBy)
	writeStringField(fc.Presentation)

	// Hash format flags
	writeBoolField(fc.Formats.HTML)
	writeBoolField(fc.Formats.SimpleHTML)
	writeBoolField(fc.HasView(models.FeedViewDefault))
	writeBoolField(fc.HasView(models.FeedViewSimple))
	writeBoolField(fc.HasView(models.FeedViewCalendar))
	writeBoolField(fc.HasView(models.FeedViewShorts))
	writeStringField(fc.ShortsURL())
	writeBoolField(fc.Formats.RSS)
	writeBoolField(fc.Formats.Atom)
	writeBoolField(fc.Formats.JSON)
	writeBoolField(fc.Formats.Markdown)
	writeBoolField(fc.Formats.Text)
	writeBoolField(fc.Formats.Sitemap)

	// Hash template overrides
	writeStringField(fc.Templates.HTML)
	writeStringField(fc.Templates.SimpleHTML)
	writeStringField(fc.Templates.RSS)
	writeStringField(fc.Templates.Atom)
	writeStringField(fc.Templates.JSON)
	writeStringField(fc.Templates.Card)
	writeStringField(fc.Templates.Sitemap)

	return hex.EncodeToString(h.Sum(nil))
}

func (p *PublishFeedsPlugin) resetPostFeedHashCache() {
	p.postFeedHashMu.Lock()
	defer p.postFeedHashMu.Unlock()
	p.postFeedHashCache = make(map[string]string)
}

func (p *PublishFeedsPlugin) getPostFeedItemHash(post *models.Post, cache *buildcache.Cache) string {
	if post == nil {
		return ""
	}
	// The load-stage semantic hash predates encryption and cannot detect a
	// changed encrypted wrapper or key rotation. Private feed entries must use
	// the current safe projection after encryption instead of that cached hash.
	if post.Private {
		return computePrivateFeedItemHash(post)
	}
	if post.Path == "" || cache == nil {
		return computePostFeedItemHash(post)
	}

	p.postFeedHashMu.RLock()
	if cachedHash, ok := p.postFeedHashCache[post.Path]; ok {
		p.postFeedHashMu.RUnlock()
		return cachedHash
	}
	p.postFeedHashMu.RUnlock()

	feedHash := ""
	feedHash, _, _ = cache.GetPostSemanticHashes(post.Path)
	if feedHash == "" {
		return computePostFeedItemHash(post)
	}

	p.postFeedHashMu.Lock()
	p.postFeedHashCache[post.Path] = feedHash
	p.postFeedHashMu.Unlock()

	return feedHash
}

func computePrivateFeedItemHash(post *models.Post) string {
	safe := safeFeedPost(post)
	if safe == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(computePostFeedItemHash(safe))
	b.WriteByte('\x00')
	b.WriteString(safe.ArticleHTML)
	b.WriteByte('\x00')
	b.WriteString("private")
	return buildcache.ContentHash(b.String())
}

// shouldSkipFeed checks if a feed can be skipped (incremental build).
// Returns (skip bool, hash string) - hash is returned so callers can reuse it
// for caching without recomputing.
func (p *PublishFeedsPlugin) shouldSkipFeed(fc *models.FeedConfig, cache interface{}, outputDir string) (skip bool, hash string) {
	return p.shouldSkipFeedWithConfig(fc, cache, outputDir, nil)
}

func (p *PublishFeedsPlugin) shouldSkipFeedWithConfig(fc *models.FeedConfig, cache interface{}, outputDir string, config *lifecycle.Config) (skip bool, hash string) {
	return p.shouldSkipFeedWithConfigAndChanges(fc, cache, outputDir, config, nil)
}

func (p *PublishFeedsPlugin) shouldSkipFeedWithConfigAndChanges(
	fc *models.FeedConfig,
	cache interface{},
	outputDir string,
	config *lifecycle.Config,
	affectedPaths map[string]bool,
) (skip bool, hash string) {
	if cache == nil {
		return false, p.computeFeedHashWithConfigAndCache(fc, config, nil)
	}

	bc, ok := cache.(*buildcache.Cache)
	if !ok || bc == nil {
		return false, p.computeFeedHashWithConfigAndCache(fc, config, nil)
	}

	// Always compute hash since we return it for caching.
	currentHash := p.computeFeedHashWithConfigAndCache(fc, config, bc)
	for _, post := range fc.Posts {
		if post != nil && post.Path != "" && affectedPaths[post.Path] {
			return false, currentHash
		}
	}

	// Check if feed hash changed (membership, order, page structure, content)
	cachedHash := bc.GetFeedHash(fc.Slug)
	if cachedHash != currentHash {
		return false, currentHash // Need to rebuild
	}

	for _, path := range p.expectedFeedOutputPaths(fc, outputDir, config) {
		if _, err := os.Stat(path); err != nil {
			return false, currentHash
		}
	}

	// A prior build may have cached a classic/fallback HTML page despite an
	// editorial opt-in. Existence and a matching feed hash do not prove the
	// expected template actually made it into the published release.
	if resolveFeedPresentation(fc) == presentationEditorial && fc.PaginationType != models.PaginationJS {
		indexPath := filepath.Join(p.determineFeedDir(outputDir, fc.Slug), "index.html")
		if !editorialHTMLPublished(indexPath) {
			publishFeedsLog.Infof("feed %q cached index lacks editorial markup; rebuilding", fc.Slug)
			return false, currentHash
		}
	}

	return true, currentHash // Can skip
}

// cacheFeedHash stores the feed hash in the build cache.
func (p *PublishFeedsPlugin) cacheFeedHash(fc *models.FeedConfig, cache interface{}, hash string) {
	if cache == nil {
		return
	}
	bc, ok := cache.(*buildcache.Cache)
	if !ok || bc == nil {
		return
	}
	bc.SetFeedHash(fc.Slug, hash)
}

//nolint:gocyclo // output path expansion is a direct format-to-path mapping
func (p *PublishFeedsPlugin) expectedFeedOutputPaths(fc *models.FeedConfig, outputDir string, config *lifecycle.Config) []string {
	feedDir := p.determineFeedDir(outputDir, fc.Slug)
	syndication := getSyndicationConfig(config)
	paths := make([]string, 0, len(fc.Pages)*2+8)
	htmlPages := fc.Pages
	if fc.Formats.HTML || simpleHTMLViewEnabled(fc) {
		htmlFeed := feedConfigWithRenderablePosts(fc)
		ensureDefaultCollectionPage(htmlFeed)
		htmlPages = htmlFeed.Pages
	}
	add := func(path string) {
		if path != "" {
			paths = append(paths, path)
		}
	}

	if fc.Formats.HTML && len(htmlPages) > 0 {
		add(filepath.Join(feedDir, "index.html"))
		if lastPage := htmlPages[len(htmlPages)-1].Number; lastPage > 1 {
			// The feed hash already covers page count and pagination structure, so
			// checking the first and last HTML pages is enough to catch missing
			// outputs without repeating every paginated page stat on warm builds.
			add(filepath.Join(feedDir, "page", fmt.Sprintf("%d", lastPage), "index.html"))
		}
	}

	if fc.HasView(models.FeedViewShorts) {
		shortsDir := filepath.Join(outputDir, strings.Trim(fc.ShortsURL(), "/"))
		add(filepath.Join(shortsDir, "index.html"))
		add(filepath.Join(shortsDir, "data", "index.json"))
		if count := shortsPostCount(fc.Posts); count > 0 {
			add(filepath.Join(shortsDir, "data", fmt.Sprintf("%04d.json", (count-1)/shortsChunkSize)))
		}
	}

	if simpleHTMLViewEnabled(fc) && len(htmlPages) > 0 {
		add(filepath.Join(feedDir, "simple", "index.html"))
		if lastPage := htmlPages[len(htmlPages)-1].Number; lastPage > 1 {
			add(filepath.Join(feedDir, "simple", "page", fmt.Sprintf("%d", lastPage), "index.html"))
		}
	}

	if fc.Formats.RSS {
		add(filepath.Join(feedDir, "rss.xml"))
		add(filepath.Join(outputDir, "rss.xsl"))
	}
	if fc.Formats.Atom {
		add(filepath.Join(feedDir, "atom.xml"))
		add(filepath.Join(outputDir, "atom.xsl"))
	}
	if fc.Formats.JSON {
		add(filepath.Join(feedDir, "feed.json"))
		if fc.Slug != "" {
			add(filepath.Join(outputDir, fc.Slug+".json", "index.html"))
		}
	}
	if fc.Formats.Markdown {
		if fc.Slug == "" {
			add(filepath.Join(outputDir, "index.md"))
		} else {
			add(filepath.Join(outputDir, fc.Slug+".md"))
			add(filepath.Join(outputDir, fc.Slug, "index.md", "index.html"))
		}
	}
	if fc.Formats.Text {
		if fc.Slug == "" {
			add(filepath.Join(outputDir, "index.txt"))
		} else {
			add(filepath.Join(outputDir, fc.Slug+".txt"))
			add(filepath.Join(outputDir, fc.Slug, "index.txt", "index.html"))
		}
	}
	if fc.Formats.Sitemap {
		add(filepath.Join(feedDir, "sitemap.xml"))
	}
	if fc.Formats.HTML && fc.HasView(models.FeedViewCalendar) {
		add(filepath.Join(feedDir, "calendar.json"))
	}
	if sidebarChunksEnabled(fc) {
		count := len(publicSidebarChunkPosts(fc))
		if count > 0 {
			add(filepath.Join(feedDir, "_sidebar", "0000.json"))
			add(filepath.Join(feedDir, "_sidebar", fmt.Sprintf("%04d.json", (count-1)/sidebarNavChunkSize)))
		}
	}

	archiveDir := feedArchiveDir(feedDir)
	if fc.Formats.RSS && shouldGenerateFeedArchive(fc, syndication) {
		add(filepath.Join(archiveDir, "rss.xml"))
	}
	if fc.Formats.Atom && shouldGenerateFeedArchive(fc, syndication) {
		add(filepath.Join(archiveDir, "atom.xml"))
	}
	if fc.Formats.JSON && shouldGenerateFeedArchive(fc, syndication) {
		add(filepath.Join(archiveDir, "feed.json"))
	}

	return paths
}

// feedFormatPublisher defines how to publish a specific feed format.
type feedFormatPublisher struct {
	name       string // Format name for error messages
	enabled    bool   // Whether this format is enabled
	publish    func() error
	ext        string // File extension for redirect (empty if no redirect needed)
	targetFile string // Target file name for redirect
}

type feedRenderContext struct {
	totalPosts        int
	latestPostDate    string
	sparklinePoints   string
	sparklineData     []SparklinePoint
	sparklineValues   string
	sparklineTitle    string
	sparklineSummary  string
	sparklineStart    string
	sparklineStartISO string
	sparklineEnd      string
}

// publishFeed publishes a single feed in all configured formats.
func (p *PublishFeedsPlugin) publishFeed(fc *models.FeedConfig, config *lifecycle.Config, outputDir string) error {
	feedDir := p.determineFeedDir(outputDir, fc.Slug)
	syndication := getSyndicationConfig(config)
	pagePosts, outputPosts := splitFeedRenderablePosts(fc.Posts, fc.IncludesPrivate())
	needsHTML := fc.Formats.HTML || simpleHTMLViewEnabled(fc) || fc.HasView(models.FeedViewShorts)
	needsOutputPosts := fc.Formats.RSS || fc.Formats.Atom || fc.Formats.JSON || fc.Formats.Markdown || fc.Formats.Text || fc.Formats.Sitemap

	var modelsConfig *models.Config
	var htmlFC *models.FeedConfig
	var syndicationFC *models.FeedConfig
	var renderCtx *feedRenderContext

	if needsHTML {
		modelsConfig = ToModelsConfig(config)
		htmlFC = cloneFeedConfigWithPosts(fc, pagePosts)
		baseURL := "/" + htmlFC.Slug
		if htmlFC.Slug == "" {
			baseURL = "/"
		}
		htmlFC.Paginate(baseURL)
		ensureDefaultCollectionPage(htmlFC)
		renderCtx = buildFeedRenderContext(htmlFC)
	}

	if needsOutputPosts {
		syndicationFC = cloneFeedConfigWithPosts(fc, outputPosts)
	}

	if err := os.MkdirAll(feedDir, 0o755); err != nil {
		return fmt.Errorf("creating feed directory: %w", err)
	}
	if err := cleanupDisabledFeedViewOutputs(fc, outputDir, feedDir); err != nil {
		return err
	}

	// Define all format publishers with their configurations
	publishers := []feedFormatPublisher{
		{name: "HTML", enabled: fc.Formats.HTML, publish: func() error { return p.publishHTMLPages(htmlFC, config, modelsConfig, feedDir, renderCtx) }},
		{name: "Calendar data", enabled: fc.Formats.HTML && fc.HasView(models.FeedViewCalendar), publish: func() error { return p.publishCalendarJSON(htmlFC, feedDir) }},
		{name: "Sidebar chunks", enabled: sidebarChunksEnabled(fc), publish: func() error { return p.publishSidebarChunks(fc, feedDir) }},
		{name: "SimpleHTML", enabled: simpleHTMLViewEnabled(fc), publish: func() error { return p.publishSimpleHTMLPages(htmlFC, config, modelsConfig, feedDir, renderCtx) }},
		{name: "Shorts", enabled: fc.HasView(models.FeedViewShorts), publish: func() error { return p.publishShortsPages(fc, config, outputDir) }},
		{name: "RSS", enabled: fc.Formats.RSS, publish: func() error { return p.publishRSS(syndicationFC, config, feedDir, false) }},
		{name: "Atom", enabled: fc.Formats.Atom, publish: func() error { return p.publishAtom(syndicationFC, config, feedDir, false) }},
		{name: "JSON", enabled: fc.Formats.JSON, publish: func() error { return p.publishJSON(syndicationFC, config, feedDir, false) }, ext: "json", targetFile: "feed.json"},
		{name: "Markdown", enabled: fc.Formats.Markdown, publish: func() error { return p.publishMarkdown(syndicationFC, fc.Slug, outputDir) }, ext: "md", targetFile: ""},
		{name: "Text", enabled: fc.Formats.Text, publish: func() error { return p.publishText(syndicationFC, fc.Slug, outputDir) }, ext: "txt", targetFile: ""},
		{name: "Sitemap", enabled: fc.Formats.Sitemap, publish: func() error { return p.publishSitemap(syndicationFC, config, feedDir) }},
	}

	for _, pub := range publishers {
		if err := p.publishFormat(pub, fc.Slug, outputDir); err != nil {
			return err
		}
	}

	if shouldGenerateFeedArchive(fc, syndication) {
		archiveDir := feedArchiveDir(feedDir)
		if err := os.MkdirAll(archiveDir, 0o755); err != nil {
			return fmt.Errorf("creating archive feed directory: %w", err)
		}

		archivePublishers := []feedFormatPublisher{
			{name: "Archive RSS", enabled: fc.Formats.RSS, publish: func() error { return p.publishRSS(syndicationFC, config, archiveDir, true) }},
			{name: "Archive Atom", enabled: fc.Formats.Atom, publish: func() error { return p.publishAtom(syndicationFC, config, archiveDir, true) }},
			{name: "Archive JSON", enabled: fc.Formats.JSON, publish: func() error { return p.publishJSON(syndicationFC, config, archiveDir, true) }},
		}

		for _, pub := range archivePublishers {
			if err := pub.publish(); err != nil {
				return fmt.Errorf("publishing %s for %q: %w", pub.name, fc.Slug, err)
			}
		}
	}

	return nil
}

// cleanupDisabledFeedViewOutputs prunes only outputs owned by this feed.
func cleanupDisabledFeedViewOutputs(fc *models.FeedConfig, outputDir, feedDir string) error {
	if !fc.HasView(models.FeedViewShorts) {
		if err := cleanupDisabledShorts(outputDir, fc.ShortsURL(), fc.Slug); err != nil {
			return err
		}
	}
	if !simpleHTMLViewEnabled(fc) {
		if err := os.RemoveAll(filepath.Join(feedDir, "simple")); err != nil {
			return fmt.Errorf("removing disabled simple feed output: %w", err)
		}
	}
	if !fc.Formats.HTML || !fc.HasView(models.FeedViewCalendar) {
		if err := os.Remove(filepath.Join(feedDir, "calendar.json")); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing disabled calendar data: %w", err)
		}
	}
	if !sidebarChunksEnabled(fc) {
		if err := os.RemoveAll(filepath.Join(feedDir, "_sidebar")); err != nil {
			return fmt.Errorf("removing disabled sidebar chunks: %w", err)
		}
	}
	return nil
}

func recordFeedOutput(m *lifecycle.Manager, fc *models.FeedConfig) {
	if m == nil || fc == nil || m.ContentLedger() == nil {
		return
	}

	pagePosts, outputPosts := splitFeedRenderablePosts(fc.Posts, fc.IncludesPrivate())
	emittedPaths := feedEmittedPaths(fc, pagePosts, outputPosts)

	for _, post := range fc.Posts {
		if post == nil || post.Path == "" {
			continue
		}
		reasons := feedEligibilityReasons(post, fc.IncludesPrivate())
		if _, emitted := emittedPaths[post.Path]; emitted {
			m.ContentLedger().RecordFeed(post.Path, fc.Slug, len(reasons) == 0, reasons...)
			m.ContentLedger().MarkEmitted(post.Path)
			continue
		}

		if len(reasons) == 0 {
			reasons = append(reasons, diagnostics.ReasonContentNoOutput)
		}
		m.ContentLedger().RecordFeed(post.Path, fc.Slug, false, reasons...)
	}
}

func feedEmittedPaths(fc *models.FeedConfig, pagePosts, outputPosts []*models.Post) map[string]struct{} {
	emittedPaths := make(map[string]struct{}, len(pagePosts)+len(outputPosts))
	addPosts := func(posts []*models.Post) {
		for _, post := range posts {
			if post != nil && post.Path != "" {
				emittedPaths[post.Path] = struct{}{}
			}
		}
	}

	if fc.Formats.HTML || simpleHTMLViewEnabled(fc) || fc.HasView(models.FeedViewShorts) {
		addPosts(pagePosts)
	}
	if fc.Formats.RSS || fc.Formats.Atom || fc.Formats.JSON || fc.Formats.Markdown || fc.Formats.Text {
		addPosts(outputPosts)
	}
	if fc.Formats.Sitemap {
		for _, post := range outputPosts {
			if post == nil || !post.Published || post.Private {
				continue
			}
			if post.Path != "" {
				emittedPaths[post.Path] = struct{}{}
			}
		}
	}

	return emittedPaths
}

func recordFeedWriteError(m *lifecycle.Manager, fc *models.FeedConfig) {
	if m == nil || fc == nil || m.ContentLedger() == nil {
		return
	}
	for _, post := range fc.Posts {
		if post == nil || post.Path == "" {
			continue
		}
		m.ContentLedger().RecordError(post.Path, diagnostics.ReasonContentWriteError, "feed output could not be written")
	}
}

func feedConfigWithRenderablePosts(fc *models.FeedConfig) *models.FeedConfig {
	clone := cloneFeedConfigWithPosts(fc, filterFeedPagePosts(fc.Posts, fc.IncludesPrivate()))
	baseURL := "/" + clone.Slug
	if clone.Slug == "" {
		baseURL = "/"
	}
	clone.Paginate(baseURL)
	return clone
}

// ensureDefaultCollectionPage gives a configless site addressable home and
// archive pages before it has published content. Other empty feeds retain
// their existing no-page behavior.
func ensureDefaultCollectionPage(fc *models.FeedConfig) {
	if fc == nil || (fc.Slug != "" && fc.Slug != defaultArchivePrefix) || len(fc.Pages) != 0 {
		return
	}
	fc.Pages = []models.FeedPage{{
		Number:         1,
		Posts:          []*models.Post{},
		TotalPages:     1,
		ItemsPerPage:   fc.ItemsPerPage,
		PaginationType: fc.PaginationType,
	}}
}

func feedConfigWithOutputPosts(fc *models.FeedConfig) *models.FeedConfig {
	clone := cloneFeedConfigWithPosts(fc, filterFeedOutputPosts(fc.Posts, fc.IncludesPrivate()))
	baseURL := "/" + clone.Slug
	if clone.Slug == "" {
		baseURL = "/"
	}
	clone.Paginate(baseURL)
	return clone
}

// determineFeedDir returns the output directory for a feed based on its slug.
func (p *PublishFeedsPlugin) determineFeedDir(outputDir, slug string) string {
	if slug == "" {
		return outputDir
	}
	return filepath.Join(outputDir, slug)
}

// prepareWriteFile removes any existing directory at a file path.
// This handles the case where a previous build created a directory where a file should be.
func (p *PublishFeedsPlugin) prepareWriteFile(path string) error {
	// Check if path exists as a directory
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		// Remove the directory so we can write a file
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("removing directory at file path %s: %w", path, err)
		}
	}
	return nil
}

// safeWriteFile writes content to a file, removing any existing directory at that path.
// If the file already exists with identical content, it is left untouched.
func (p *PublishFeedsPlugin) safeWriteFile(path string, content []byte) error {
	if err := p.prepareWriteFile(path); err != nil {
		return err
	}

	if existing, err := os.ReadFile(path); err == nil && bytes.Equal(existing, content) {
		return nil
	}

	return p.writeFile(path, content)
}

// writeFile writes content directly without a read-before-write check.
// Use this only when the caller already knows the output is dirty.
func (p *PublishFeedsPlugin) writeFile(path string, content []byte) error {
	if err := p.prepareWriteFile(path); err != nil {
		return err
	}

	//nolint:gosec // G306: Output files need 0644 for web serving
	return os.WriteFile(path, content, 0o644)
}

// publishFormat publishes a single format if enabled and handles redirects.
func (p *PublishFeedsPlugin) publishFormat(pub feedFormatPublisher, slug, outputDir string) error {
	if !pub.enabled {
		return nil
	}

	if err := pub.publish(); err != nil {
		return fmt.Errorf("publishing %s: %w", pub.name, err)
	}

	// Write redirect for non-root feeds with redirect configuration
	if slug != "" && pub.ext != "" {
		if pub.targetFile == "" {
			// Reversed redirect: /slug/index.ext -> /slug.ext (for Markdown/Text)
			if err := p.writeReversedFeedRedirect(slug, pub.ext, outputDir); err != nil {
				return fmt.Errorf("writing %s redirect: %w", pub.name, err)
			}
		} else {
			// Forward redirect: /slug.ext -> /slug/targetFile (for JSON)
			if err := p.writeFeedFormatRedirect(slug, pub.ext, pub.targetFile, outputDir); err != nil {
				return fmt.Errorf("writing %s redirect: %w", pub.name, err)
			}
		}
	}

	return nil
}

// publishHTMLPages publishes HTML pages for a paginated feed.
func (p *PublishFeedsPlugin) publishHTMLPages(fc *models.FeedConfig, config *lifecycle.Config, modelsConfig *models.Config, feedDir string, renderCtx *feedRenderContext) error {
	if err := p.cleanupPaginatedFeedDirs(feedDir, "", fc.Pages); err != nil {
		return fmt.Errorf("cleaning html pagination dirs: %w", err)
	}

	for i := range fc.Pages {
		page := &fc.Pages[i]
		// Determine output path
		var pagePath string
		if page.Number == 1 {
			pagePath = filepath.Join(feedDir, "index.html")
		} else {
			pageDir := filepath.Join(feedDir, "page", fmt.Sprintf("%d", page.Number))
			if err := os.MkdirAll(pageDir, 0o755); err != nil {
				return fmt.Errorf("creating page directory: %w", err)
			}
			pagePath = filepath.Join(pageDir, "index.html")
		}

		// Generate HTML content
		html, err := p.generateFeedPageHTML(fc, page, config, modelsConfig, renderCtx)
		if err != nil {
			return fmt.Errorf("generating page %d: %w", page.Number, err)
		}

		if err := p.writeFile(pagePath, []byte(html)); err != nil {
			return fmt.Errorf("writing page %d: %w", page.Number, err)
		}
	}

	return nil
}

// publishSimpleHTMLPages publishes the simple (compact list) HTML pages for a feed.
// Output is written to feedDir/simple/ with pagination at feedDir/simple/page/N/.
func (p *PublishFeedsPlugin) publishSimpleHTMLPages(fc *models.FeedConfig, config *lifecycle.Config, modelsConfig *models.Config, feedDir string, renderCtx *feedRenderContext) error {
	simpleDir := filepath.Join(feedDir, "simple")
	if err := p.cleanupPaginatedFeedDirs(feedDir, "simple", fc.Pages); err != nil {
		return fmt.Errorf("cleaning simple pagination dirs: %w", err)
	}

	// Compute the base URL prefix for simple feed pagination links.
	// For a feed with slug "blog", this is "/blog/simple".
	// For the root feed (slug ""), this is "/simple".
	simpleBaseURL := "/" + fc.Slug + "/simple"
	if fc.Slug == "" {
		simpleBaseURL = "/simple"
	}

	for i := range fc.Pages {
		srcPage := &fc.Pages[i]

		// Create an adjusted copy of the page with URLs pointing to /simple/ paths
		adjustedPage := p.adjustPageURLsForSimple(srcPage, simpleBaseURL)

		// Determine output path
		var pagePath string
		if adjustedPage.Number == 1 {
			if err := os.MkdirAll(simpleDir, 0o755); err != nil {
				return fmt.Errorf("creating simple feed directory: %w", err)
			}
			pagePath = filepath.Join(simpleDir, "index.html")
		} else {
			pageDir := filepath.Join(simpleDir, "page", fmt.Sprintf("%d", adjustedPage.Number))
			if err := os.MkdirAll(pageDir, 0o755); err != nil {
				return fmt.Errorf("creating simple feed page directory: %w", err)
			}
			pagePath = filepath.Join(pageDir, "index.html")
		}

		// Generate HTML content using simple-feed.html template
		htmlContent, err := p.generateSimpleFeedPageHTML(fc, &adjustedPage, config, modelsConfig, renderCtx)
		if err != nil {
			return fmt.Errorf("generating simple page %d: %w", adjustedPage.Number, err)
		}

		if err := p.writeFile(pagePath, []byte(htmlContent)); err != nil {
			return fmt.Errorf("writing simple page %d: %w", adjustedPage.Number, err)
		}
	}

	return nil
}

func (p *PublishFeedsPlugin) cleanupPaginatedFeedDirs(feedDir, formatSubdir string, pages []models.FeedPage) error {
	pageRoot := filepath.Join(feedDir, formatSubdir, "page")
	entries, err := os.ReadDir(pageRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	validPages := make(map[int]struct{}, len(pages))
	for i := range pages {
		if pages[i].Number > 1 {
			validPages[pages[i].Number] = struct{}{}
		}
	}

	if len(validPages) == 0 {
		return os.RemoveAll(pageRoot)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pageNum, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if _, ok := validPages[pageNum]; ok {
			continue
		}
		if err := os.RemoveAll(filepath.Join(pageRoot, entry.Name())); err != nil {
			return err
		}
	}

	return nil
}

// adjustPageURLsForSimple creates a copy of FeedPage with URLs adjusted for the simple feed subdirectory.
func (p *PublishFeedsPlugin) adjustPageURLsForSimple(page *models.FeedPage, simpleBaseURL string) models.FeedPage {
	adjusted := *page

	// Adjust prev/next URLs
	if adjusted.HasPrev {
		if adjusted.Number == 2 {
			adjusted.PrevURL = simpleBaseURL + "/"
		} else {
			adjusted.PrevURL = simpleBaseURL + "/page/" + fmt.Sprintf("%d", adjusted.Number-1) + "/"
		}
	}
	if adjusted.HasNext {
		adjusted.NextURL = simpleBaseURL + "/page/" + fmt.Sprintf("%d", adjusted.Number+1) + "/"
	}

	// Adjust page URLs for numbered navigation
	adjustedURLs := make([]string, len(adjusted.PageURLs))
	for i := range adjusted.PageURLs {
		if i == 0 {
			adjustedURLs[i] = simpleBaseURL + "/"
		} else {
			adjustedURLs[i] = simpleBaseURL + "/page/" + fmt.Sprintf("%d", i+1) + "/"
		}
	}
	adjusted.PageURLs = adjustedURLs

	return adjusted
}

// generateSimpleFeedPageHTML generates HTML for a simple feed page using the simple-feed.html template.
func (p *PublishFeedsPlugin) generateSimpleFeedPageHTML(fc *models.FeedConfig, page *models.FeedPage, config *lifecycle.Config, modelsConfig *models.Config, renderCtx *feedRenderContext) (string, error) {
	// Get templates directory from config
	templatesDir := PluginNameTemplates
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}

	// Get theme name from config (default to "default")
	themeName := ThemeDefault
	if extra := config.Extra; extra != nil {
		if theme, ok := extra["theme"].(models.ThemeConfig); ok {
			if theme.Name != "" {
				themeName = theme.Name
			}
		}
		if theme, ok := extra["theme"].(map[string]interface{}); ok {
			if name, ok := theme["name"].(string); ok && name != "" {
				themeName = name
			}
		}
		if name, ok := extra["theme"].(string); ok && name != "" {
			themeName = name
		}
	}

	// Determine which template to use
	templateName := fc.Templates.SimpleHTML
	if templateName == "" {
		templateName = "simple-feed.html"
	}

	// Try to use pongo2 template engine (cached)
	engine, err := p.getOrCreateEngine(templatesDir, themeName)
	if err == nil && engine.TemplateExists(templateName) {
		if modelsConfig == nil {
			modelsConfig = ToModelsConfig(config)
		}
		ctx := templates.NewFeedContext(fc, page, modelsConfig)
		ctx.Set("nav_previews", p.navPreviews)
		p.addFeedStatsContext(&ctx, fc, renderCtx)
		setEncryptedFeedContext(&ctx, page)

		// Feed pages always need cards CSS
		ctx.Set("needs_cards_css", true)

		htmlContent, err := engine.Render(templateName, ctx)
		if err != nil {
			publishFeedsLog.Warnf("template rendering failed for %s: %v (falling back to built-in template)", templateName, err)
		} else {
			return htmlContent, nil
		}
	}

	// Fallback: use the standard HTML fallback (better than nothing)
	return p.generateFeedPageHTMLFallback(fc, page, config)
}

// generateFeedPageHTML generates HTML for a feed page.
//
//nolint:gocyclo // Template selection, pagination and fail-closed presentation validation share the render path.
func (p *PublishFeedsPlugin) generateFeedPageHTML(fc *models.FeedConfig, page *models.FeedPage, config *lifecycle.Config, modelsConfig *models.Config, renderCtx *feedRenderContext) (string, error) {
	presentation := resolveFeedPresentation(fc)
	// JS pagination requires direct children and cannot render nested weeks.
	if page.PaginationType == models.PaginationJS && presentation == presentationEditorial {
		presentation = presentationClassic
	}
	templateName := fc.Templates.HTML
	if templateName == "" {
		templateName = "feed.html"
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

	// Try to use pongo2 template engine with the configured HTML template (cached)
	engine, err := p.getOrCreateEngine(templatesDir, themeName)
	if err == nil && engine.TemplateExists(templateName) {
		// Use shared config conversion to ensure all fields are available
		// (search, head, components, theme, etc. - required by base.html)
		if modelsConfig == nil {
			modelsConfig = ToModelsConfig(config)
		}

		// Create feed context
		ctx := templates.NewFeedContext(fc, page, modelsConfig)
		ctx.Set("nav_previews", p.navPreviews)
		ctx.Set("feed_robots", fc.Robots)
		ctx.Set("feed_presentation", presentation)
		if presentation == presentationEditorial {
			ctx.Set("feed_sections", buildEditorialFeedSections(page.Posts))
		}
		p.addFeedStatsContext(&ctx, fc, renderCtx)

		// Feed pages always need cards CSS
		ctx.Set("needs_cards_css", true)

		// If any post on this page has encrypted content, load decryption JS/CSS
		setEncryptedFeedContext(&ctx, page)

		// Render with pongo2 template. An explicitly selected editorial feed
		// must not silently degrade to the legacy HTML after a template error.
		html, renderErr := engine.Render(templateName, ctx)
		if renderErr != nil {
			if presentation == presentationEditorial {
				return "", fmt.Errorf("editorial feed %q page %d template %q failed: %w", fc.Slug, page.Number, templateName, renderErr)
			}
			publishFeedsLog.Warnf("template rendering failed for %s: %v (falling back to built-in template)", templateName, renderErr)
		} else {
			if presentation == presentationEditorial && !editorialHTMLMarkupValid(html) {
				return "", fmt.Errorf("editorial feed %q page %d template %q rendered without editorial markers; check theme/template overrides", fc.Slug, page.Number, templateName)
			}
			if page.Number == 1 && presentation == presentationEditorial {
				publishFeedsLog.Infof("feed %q rendered with editorial presentation and %d weekly sections", fc.Slug, len(buildEditorialFeedSections(page.Posts)))
			}
			return html, nil
		}
	}

	if presentation == presentationEditorial {
		if err != nil {
			return "", fmt.Errorf("editorial feed %q template engine unavailable: %w", fc.Slug, err)
		}
		return "", fmt.Errorf("editorial feed %q template %q not found; cannot use classic fallback", fc.Slug, templateName)
	}

	// Legacy feeds continue to support the built-in Go template fallback.
	return p.generateFeedPageHTMLFallback(fc, page, config)
}

// setEncryptedFeedContext enables the assets required to unlock encrypted feed
// entries. Keep this shared by the normal and simple feed page renderers.
func setEncryptedFeedContext(ctx *templates.Context, page *models.FeedPage) {
	for _, post := range page.Posts {
		if post != nil {
			if v, ok := post.Extra["has_encrypted_content"].(bool); ok && v {
				ctx.Set("has_encrypted_content", true)
				return
			}
		}
	}
}

func buildFeedRenderContext(fc *models.FeedConfig) *feedRenderContext {
	if fc == nil {
		return nil
	}
	totalPosts, latestPostDate, _ := feedStats(fc.Posts, fc.IncludesPrivate())
	window := computeSparklineWindow(fc.Posts, fc.IncludesPrivate())
	sparklineData := buildFeedSparklineData(fc.Posts, window, fc.IncludesPrivate())
	return &feedRenderContext{
		totalPosts:        totalPosts,
		latestPostDate:    latestPostDate,
		sparklinePoints:   buildFeedSparkline(fc.Posts, window, fc.IncludesPrivate()),
		sparklineData:     sparklineData,
		sparklineValues:   encodeFeedSparklineValues(sparklineData),
		sparklineTitle:    buildFeedSparklineTitle(fc.Posts, window, fc.IncludesPrivate()),
		sparklineSummary:  buildFeedSparklineSummary(fc.Posts, window, fc.IncludesPrivate()),
		sparklineStart:    buildFeedSparklineStart(window),
		sparklineStartISO: buildFeedSparklineStartISO(window),
		sparklineEnd:      buildFeedSparklineEnd(window),
	}
}

func (p *PublishFeedsPlugin) addFeedStatsContext(ctx *templates.Context, fc *models.FeedConfig, renderCtx *feedRenderContext) {
	if ctx == nil || fc == nil || renderCtx == nil {
		return
	}
	ctx.Set("feed_stats_total_posts", renderCtx.totalPosts)
	ctx.Set("feed_stats_latest_post", renderCtx.latestPostDate)
	ctx.Set("feed_sparkline_points", renderCtx.sparklinePoints)
	ctx.Set("feed_sparkline_data", renderCtx.sparklineData)
	ctx.Set("feed_sparkline_values", renderCtx.sparklineValues)
	ctx.Set("feed_sparkline_title", renderCtx.sparklineTitle)
	ctx.Set("feed_sparkline_summary", renderCtx.sparklineSummary)
	ctx.Set("feed_sparkline_start", renderCtx.sparklineStart)
	ctx.Set("feed_sparkline_start_iso", renderCtx.sparklineStartISO)
	ctx.Set("feed_sparkline_end", renderCtx.sparklineEnd)
}

// generateFeedPageHTMLFallback generates HTML using a built-in Go template.
func (p *PublishFeedsPlugin) generateFeedPageHTMLFallback(fc *models.FeedConfig, page *models.FeedPage, config *lifecycle.Config) (string, error) {
	siteURL := getSiteURL(config)
	siteTitle := getSiteTitle(config)

	title := fc.Title
	if title == "" {
		title = siteTitle
	}

	// Simple default template with CSS links
	tmplStr := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/css/variables.css">
    <link rel="stylesheet" href="/css/palette.css">
    <link rel="stylesheet" href="/css/main.css">
    <link rel="stylesheet" href="/css/admonitions.css">
    <link rel="stylesheet" href="/css/code.css">
    <link rel="alternate" type="application/rss+xml" title="RSS Feed" href="rss.xml">
    <link rel="alternate" type="application/atom+xml" title="Atom Feed" href="atom.xml">
    <link rel="alternate" type="application/feed+json" title="JSON Feed" href="feed.json">
</head>
<body>
    <header>
        <nav>
            <a href="/">{{.SiteTitle}}</a>
            <a href="/blog/">Blog</a>
        </nav>
    </header>
    <main>
        <h1>{{.Title}}</h1>
        {{if .Description}}<p class="description">{{.Description}}</p>{{end}}
        <div class="post-list">
        {{range .Posts}}
            <article class="post-card">
                <a href="{{.Href}}">
                    <h2>{{if .Title}}{{.Title}}{{else}}{{.Slug}}{{end}}</h2>
                </a>
				{{if .Date}}<time datetime="{{.Date.Format "2006-01-02"}}">{{humanDate .Date}}</time>{{end}}
                {{if .Description}}<p>{{.Description}}</p>{{end}}
            </article>
        {{end}}
        </div>
    </main>
    <nav class="pagination">
        {{if .HasPrev}}<a href="{{.PrevURL}}" rel="prev">&laquo; Newer</a>{{end}}
        <span>Page {{.Number}}</span>
        {{if .HasNext}}<a href="{{.NextURL}}" rel="next">Older &raquo;</a>{{end}}
    </nav>
    <footer>
        <p>&copy; {{.SiteTitle}}</p>
    </footer>
</body>
</html>`

	tmpl, err := template.New("feed").Funcs(template.FuncMap{"humanDate": templates.FormatHumanDate}).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	// Prepare template data
	data := struct {
		Title       string
		Description string
		Posts       []*models.Post
		HasPrev     bool
		HasNext     bool
		PrevURL     string
		NextURL     string
		Number      int
		SiteURL     string
		SiteTitle   string
	}{
		Title:       title,
		Description: fc.Description,
		Posts:       page.Posts,
		HasPrev:     page.HasPrev,
		HasNext:     page.HasNext,
		PrevURL:     page.PrevURL,
		NextURL:     page.NextURL,
		Number:      page.Number,
		SiteURL:     siteURL,
		SiteTitle:   siteTitle,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

// publishRSS generates and writes an RSS feed.
func (p *PublishFeedsPlugin) publishRSS(fc *models.FeedConfig, config *lifecycle.Config, feedDir string, archive bool) error {
	rssFeed := p.syndicationFeedConfig(fc, config, archive)
	rss, err := GenerateRSSFromFeedConfig(rssFeed, config)
	if err != nil {
		return err
	}

	rssPath := filepath.Join(feedDir, "rss.xml")
	return p.writeFile(rssPath, []byte(rss))
}

// publishAtom generates and writes an Atom feed.
func (p *PublishFeedsPlugin) publishAtom(fc *models.FeedConfig, config *lifecycle.Config, feedDir string, archive bool) error {
	atomFeed := p.syndicationFeedConfig(fc, config, archive)
	atom, err := GenerateAtomFromFeedConfig(atomFeed, config)
	if err != nil {
		return err
	}

	atomPath := filepath.Join(feedDir, "atom.xml")
	return p.writeFile(atomPath, []byte(atom))
}

// publishJSON generates and writes a JSON feed.
func (p *PublishFeedsPlugin) publishJSON(fc *models.FeedConfig, config *lifecycle.Config, feedDir string, archive bool) error {
	jsonFC := p.syndicationFeedConfig(fc, config, archive)
	jsonFeed, err := GenerateJSONFeedFromFeedConfig(jsonFC, config)
	if err != nil {
		return err
	}

	jsonPath := filepath.Join(feedDir, "feed.json")
	return p.writeFile(jsonPath, []byte(jsonFeed))
}

type calendarFeedItem struct {
	Date        string `json:"date"`
	Href        string `json:"href"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Image       string `json:"image,omitempty"`
}

// sidebarChunksEnabled excludes private feeds: static JSON must never reveal private
// titles or URLs. All public HTML feeds get chunks, including one-feed sidebars.
func sidebarChunksEnabled(fc *models.FeedConfig) bool {
	return fc != nil && fc.Formats.HTML && !fc.IncludesPrivate()
}

func publicSidebarChunkPosts(fc *models.FeedConfig) []*models.Post {
	if !sidebarChunksEnabled(fc) {
		return nil
	}
	pagePosts, _ := splitFeedRenderablePosts(fc.Posts, false)
	return pagePosts
}

type sidebarNavChunk struct {
	Start int               `json:"start"`
	Posts []sidebarPostJSON `json:"posts"`
}

// publishSidebarChunks writes bounded, cacheable static navigation data. It
// never includes article bodies or unpublished/private posts. Delete obsolete
// tail chunks when a feed shrinks, and rebuild them with the feed's content hash.
func (p *PublishFeedsPlugin) publishSidebarChunks(fc *models.FeedConfig, feedDir string) error {
	if !sidebarChunksEnabled(fc) {
		return nil
	}
	dir := filepath.Join(feedDir, "_sidebar")
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("cleaning old sidebar chunks: %w", err)
	}
	posts := publicSidebarChunkPosts(fc)
	if len(posts) == 0 {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating sidebar chunk directory: %w", err)
	}
	for start := 0; start < len(posts); start += sidebarNavChunkSize {
		end := start + sidebarNavChunkSize
		if end > len(posts) {
			end = len(posts)
		}
		chunk := sidebarNavChunk{Start: start, Posts: make([]sidebarPostJSON, 0, end-start)}
		for _, post := range posts[start:end] {
			chunk.Posts = append(chunk.Posts, sidebarPostFromKey(sidebarKey(post, fc.Slug), false))
		}
		data, err := json.Marshal(chunk)
		if err != nil {
			return fmt.Errorf("encoding sidebar chunk: %w", err)
		}
		name := fmt.Sprintf("%04d.json", start/sidebarNavChunkSize)
		if err := p.writeFile(filepath.Join(dir, name), data); err != nil {
			return fmt.Errorf("writing sidebar chunk %s: %w", name, err)
		}
	}
	return nil
}

// publishCalendarJSON writes only the metadata needed to build a feed calendar.
// Keeping it out of the HTML avoids one hidden DOM row per post on initial load.
func (p *PublishFeedsPlugin) publishCalendarJSON(fc *models.FeedConfig, feedDir string) error {
	items := make([]calendarFeedItem, 0, len(fc.Posts))
	for _, post := range fc.Posts {
		if post == nil || post.Date == nil || post.Href == "" {
			continue
		}
		title := post.PlainTitle()
		if title == "" {
			title = post.Slug
		}
		item := calendarFeedItem{
			Date:  post.Date.Format("2006-01-02"),
			Href:  post.Href,
			Title: title,
			Image: calendarImageForPost(post),
		}
		if post.Description != nil {
			item.Description = *post.Description
		}
		items = append(items, item)
	}

	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("encoding calendar metadata: %w", err)
	}
	return p.writeFile(filepath.Join(feedDir, "calendar.json"), data)
}

func calendarImageForPost(post *models.Post) string {
	if post == nil {
		return ""
	}
	media := postExtraString(post, "image")
	if media == "" {
		media = postExtraString(post, "cover")
	}
	if media == "" {
		media = postExtraString(post, "video")
	}
	if templates.IsVideoURL(media) {
		media = templates.PosterURLFromMap(templates.GetPostMap(post), media)
	}
	return media
}

func postExtraString(post *models.Post, key string) string {
	if post == nil || post.Extra == nil {
		return ""
	}
	switch value := post.Extra[key].(type) {
	case string:
		return strings.TrimSpace(value)
	case *string:
		if value != nil {
			return strings.TrimSpace(*value)
		}
	}
	return ""
}

func (p *PublishFeedsPlugin) syndicationFeedConfig(fc *models.FeedConfig, config *lifecycle.Config, archive bool) *models.FeedConfig {
	if archive {
		clone := cloneFeedConfigWithPosts(fc, fc.Posts)
		clone.Slug = path.Join(fc.Slug, defaultArchivePrefix)
		clone.Title = archiveTitleSuffix(clone.Title)
		return clone
	}

	if isArchiveFeed(fc) {
		return cloneFeedConfigWithPosts(fc, fc.Posts)
	}

	syndication := getSyndicationConfig(config)
	return cloneFeedConfigWithPosts(fc, limitFeedPosts(fc.Posts, syndication.MaxItems))
}

// publishMarkdown generates and writes a Markdown feed listing.
// For non-root feeds, content is written to /slug.md (canonical URL).
// For root feeds (slug=""), content is written to /index.md.
// HTML entities in titles and descriptions are decoded to readable characters.
func (p *PublishFeedsPlugin) publishMarkdown(fc *models.FeedConfig, slug, outputDir string) error {
	var sb strings.Builder

	// Title (decode HTML entities for plain text output)
	title := fc.Title
	if title == "" {
		title = "Posts"
	}
	sb.WriteString("# " + html.UnescapeString(title) + "\n\n")

	// Description (decode HTML entities)
	if fc.Description != "" {
		sb.WriteString(html.UnescapeString(fc.Description) + "\n\n")
	}

	// Posts list
	for _, post := range fc.Posts {
		postTitle := post.PlainTitle()
		if postTitle == "" {
			postTitle = post.Slug
		}

		sb.WriteString("- [" + postTitle + "](" + post.Href + ")")

		if post.Date != nil {
			sb.WriteString(" - " + post.Date.Format("2006-01-02"))
		}

		sb.WriteString("\n")
	}

	// Determine output path: /slug.md for non-root, /index.md for root
	var mdPath string
	if slug == "" {
		mdPath = filepath.Join(outputDir, "index.md")
	} else {
		mdPath = filepath.Join(outputDir, slug+".md")
	}
	return p.writeFile(mdPath, []byte(sb.String()))
}

// publishText generates and writes a plain text feed listing.
// For non-root feeds, content is written to /slug.txt (canonical URL).
// For root feeds (slug=""), content is written to /index.txt.
// HTML entities in titles and descriptions are decoded to readable characters.
func (p *PublishFeedsPlugin) publishText(fc *models.FeedConfig, slug, outputDir string) error {
	var sb strings.Builder

	// Title (decode HTML entities for plain text output)
	title := fc.Title
	if title == "" {
		title = "Posts"
	}
	title = html.UnescapeString(title)
	sb.WriteString(title + "\n")
	sb.WriteString(strings.Repeat("=", len(title)) + "\n\n")

	// Description (decode HTML entities)
	if fc.Description != "" {
		sb.WriteString(html.UnescapeString(fc.Description) + "\n\n")
	}

	// Posts list
	for _, post := range fc.Posts {
		postTitle := post.PlainTitle()
		if postTitle == "" {
			postTitle = post.Slug
		}

		if post.Date != nil {
			sb.WriteString(post.Date.Format("2006-01-02") + " - ")
		}

		sb.WriteString(postTitle + "\n")
		sb.WriteString("  " + post.Href + "\n\n")
	}

	// Determine output path: /slug.txt for non-root, /index.txt for root
	var txtPath string
	if slug == "" {
		txtPath = filepath.Join(outputDir, "index.txt")
	} else {
		txtPath = filepath.Join(outputDir, slug+".txt")
	}
	return p.writeFile(txtPath, []byte(sb.String()))
}

// publishSitemap generates and writes a sitemap XML file for feed posts.
func (p *PublishFeedsPlugin) publishSitemap(fc *models.FeedConfig, config *lifecycle.Config, feedDir string) error {
	// Get site URL
	siteURL := getSiteURL(config)
	if siteURL == "" {
		siteURL = DefaultSiteURL
	}

	// Build sitemap for this feed's posts
	sitemap := &URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]SitemapURL, 0, len(fc.Posts)),
	}

	// Add all posts in this feed
	for _, post := range fc.Posts {
		if !post.Published || post.Draft || post.Skip || post.Private {
			continue
		}

		url := SitemapURL{
			Loc: siteURL + post.Href,
		}

		// Use post.date for lastmod
		if post.Date != nil {
			url.LastMod = post.Date.Format("2006-01-02")
		}

		// Support frontmatter fields: changefreq and priority
		// Default values from spec
		changefreq := "weekly"
		priority := "0.5"

		if post.Extra != nil {
			if cf, ok := post.Extra["changefreq"].(string); ok && cf != "" {
				changefreq = cf
			}
			if p, ok := post.Extra["priority"].(string); ok && p != "" {
				priority = p
			}
			// Also support float64 for priority from TOML/JSON parsing
			if p, ok := post.Extra["priority"].(float64); ok {
				priority = fmt.Sprintf("%.1f", p)
			}
		}

		url.ChangeFreq = changefreq
		url.Priority = priority

		sitemap.URLs = append(sitemap.URLs, url)
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(sitemap, "", "    ")
	if err != nil {
		return fmt.Errorf("marshaling sitemap: %w", err)
	}

	// Add XML declaration
	xmlContent := xml.Header + string(output)

	// Write sitemap.xml
	sitemapPath := filepath.Join(feedDir, "sitemap.xml")
	return p.writeFile(sitemapPath, []byte(xmlContent))
}

// writeFeedFormatRedirect writes a redirect from /slug.ext to /slug/targetFile.
// This creates a file at slug.ext/index.html, which allows the URL /slug.ext
// (without trailing slash) to serve the HTML redirect on most static hosts.
//
// For example, requesting /archive.json will serve the redirect HTML that
// points to /archive/feed.json where the actual JSON content lives.
//
// Note: Web servers serve slug.ext/index.html when /slug.ext is requested,
// without adding a trailing slash redirect (unlike directory-only approaches).
func (p *PublishFeedsPlugin) writeFeedFormatRedirect(slug, ext, targetFile, outputDir string) error {
	// Create redirect HTML that points to the actual file
	targetURL := fmt.Sprintf("/%s/%s", slug, targetFile)
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

	// Create directory at /slug.ext/ (e.g., /archive.md/)
	redirectDir := filepath.Join(outputDir, slug+"."+ext)
	if err := os.MkdirAll(redirectDir, 0o755); err != nil {
		return fmt.Errorf("creating redirect directory %s: %w", redirectDir, err)
	}

	// Write index.html inside the directory
	// This allows /slug.ext to be served without trailing slash on most static hosts
	outputPath := filepath.Join(redirectDir, "index.html")
	if err := p.writeFile(outputPath, []byte(redirectHTML)); err != nil {
		return fmt.Errorf("writing redirect %s: %w", outputPath, err)
	}

	return nil
}

// writeReversedFeedRedirect writes a redirect from /slug/index.ext to /slug.ext.
// This is the "reversed" direction from writeFeedFormatRedirect - content is at the
// short URL (/slug.ext) and the redirect points there from the long URL (/slug/index.ext).
//
// Creates a file at /slug/index.ext/index.html, which allows the URL /slug/index.ext
// (without trailing slash) to serve the HTML redirect on most static hosts.
//
// For example, requesting /archive/index.md will serve the redirect HTML that
// points to /archive.md where the actual markdown content lives.
func (p *PublishFeedsPlugin) writeReversedFeedRedirect(slug, ext, outputDir string) error {
	// Create redirect HTML that points to the canonical short URL
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

	// Create directory at /slug/index.ext/ (e.g., /archive/index.md/)
	redirectDir := filepath.Join(outputDir, slug, "index."+ext)
	if err := os.MkdirAll(redirectDir, 0o755); err != nil {
		return fmt.Errorf("creating redirect directory %s: %w", redirectDir, err)
	}

	// Write index.html inside the directory
	// This allows /slug/index.ext to be served without trailing slash on most static hosts
	outputPath := filepath.Join(redirectDir, "index.html")
	if err := p.writeFile(outputPath, []byte(redirectHTML)); err != nil {
		return fmt.Errorf("writing redirect %s: %w", outputPath, err)
	}

	return nil
}

// copyXSLStylesheets copies XSL stylesheets to the output directory for styled RSS/Atom feeds.
// It searches for XSL files in the following order:
// 1. User's templates directory (if configured)
// 2. Embedded default theme templates (fallback)
//
// Stylesheets containing the markata:theme-head marker region get the site's
// theme head injected (see themeFeedXSL).
func (p *PublishFeedsPlugin) copyXSLStylesheets(config *lifecycle.Config, outputDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Get templates directory from config
	templatesDir := PluginNameTemplates
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}

	// List of XSL files to copy
	xslFiles := []string{"rss.xsl", "atom.xsl"}

	for _, xslFile := range xslFiles {
		var content []byte
		var err error

		// First, try to read from user's templates directory
		srcPath := filepath.Join(templatesDir, xslFile)
		if _, statErr := os.Stat(srcPath); statErr == nil {
			// User has their own XSL file, use it
			content, err = os.ReadFile(srcPath)
			if err != nil {
				return fmt.Errorf("reading XSL file %s: %w", srcPath, err)
			}
		} else if os.IsNotExist(statErr) {
			// No user file, try embedded templates as fallback
			content, err = themes.ReadTemplate(xslFile)
			if err != nil {
				// XSL file doesn't exist in embedded templates either, skip it
				continue
			}
		} else {
			return fmt.Errorf("checking XSL file %s: %w", srcPath, statErr)
		}

		content = p.themeFeedXSL(content, config, templatesDir)

		// Write to output directory
		dstPath := filepath.Join(outputDir, xslFile)
		if err := p.safeWriteFile(dstPath, content); err != nil {
			return fmt.Errorf("writing XSL file %s: %w", dstPath, err)
		}
	}

	return nil
}

// Ensure PublishFeedsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*PublishFeedsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*PublishFeedsPlugin)(nil)
	_ lifecycle.WritePlugin     = (*PublishFeedsPlugin)(nil)
)
