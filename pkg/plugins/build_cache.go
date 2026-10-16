// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata-go/pkg/buildcache"
	"github.com/WaylonWalker/markata-go/pkg/config"
	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/logging"
	"github.com/WaylonWalker/markata-go/pkg/themes"
)

var buildCacheLog = logging.Component("build_cache")

// BuildCachePlugin manages incremental build caching.
// It loads the build cache at the start and saves it at the end.
// Other plugins can use the cache to skip unchanged posts.
type BuildCachePlugin struct {
	cache   *buildcache.Cache
	enabled bool
}

// NewBuildCachePlugin creates a new BuildCachePlugin.
func NewBuildCachePlugin() *BuildCachePlugin {
	return &BuildCachePlugin{
		enabled: true,
	}
}

// Name returns the unique name of the plugin.
func (p *BuildCachePlugin) Name() string {
	return "build_cache"
}

// Priority returns high priority so it runs early in configure and late in cleanup.
// Load invalidation runs after posts are available but before transforms consume cached
// derivatives. Transform runs late to collect dependencies after links are processed.
func (p *BuildCachePlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageConfigure:
		return lifecycle.PriorityEarly - 100 // Very early in configure
	case lifecycle.StageLoad:
		return lifecycle.PriorityLate + 100 // After all posts and synthetic posts are loaded
	case lifecycle.StageTransform:
		return lifecycle.PriorityLate + 100 // Very late in transform (after wikilinks/embeds)
	case lifecycle.StageCleanup:
		return lifecycle.PriorityLate + 100 // Very late in cleanup
	default:
		return lifecycle.PriorityDefault
	}
}

// Configure loads the build cache and checks for global invalidation.
func (p *BuildCachePlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Check if caching is disabled
	if !p.isEnabled(config) {
		p.enabled = false
		return nil
	}

	// Load existing cache. Default to content-dir-local .markata so incremental
	// build state survives when output is redirected to a temp/workspace path.
	cacheDir := filepath.Join(config.ContentDir, ".markata")
	if config.Extra != nil {
		if dir, ok := config.Extra["cache_dir"].(string); ok && dir != "" {
			cacheDir = dir
		}
	}

	cache, err := buildcache.Load(cacheDir)
	if err != nil {
		// Non-fatal: start with empty cache
		cache = buildcache.New(cacheDir)
	}
	p.cache = cache

	// Store cache in manager for other plugins to access
	m.Cache().Set("build_cache", cache)

	// Compute and check config hash
	// Use explicit config path(s) when available, otherwise fall back to local defaults
	configFiles := []string{"markata-go.toml", "markata-go.yaml", "markata-go.json"}
	if config.Extra != nil {
		if paths, ok := config.Extra["config_paths"].([]string); ok && len(paths) > 0 {
			configFiles = paths
		} else if path, ok := config.Extra["config_path"].(string); ok && path != "" {
			configFiles = []string{path}
		}
	}
	configHash := buildcache.ContentHash(configHashInput(config, configFiles))
	if configHash != "" && cache.SetConfigHash(configHash) {
		// Config changed - cache was invalidated
		buildCacheLog.Phase("configure").Printf("Config changed, full rebuild required")
	}

	// Compute and check templates hash. The fingerprint is useful for avoiding
	// metadata work, but it is not sufficient to invalidate rendered pages:
	// template contents can change while preserving size and modtime (for
	// example after a checkout or archive extraction). Always compare the
	// content hash so cached full-page HTML cannot outlive its template.
	templatesDir := PluginNameTemplates // default "templates"
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}
	if fingerprint, fingerprintErr := buildcache.HashDirectoryState(templatesDir, []string{".html", ".txt", ".md"}); fingerprintErr == nil && fingerprint != "" {
		cache.SetTemplatesFingerprint(fingerprint)
	}
	projectTemplatesHash, err := buildcache.HashDirectory(templatesDir, []string{".html", ".txt", ".md"})
	if err != nil {
		return fmt.Errorf("hash project templates: %w", err)
	}
	embeddedTemplatesHash, err := buildcache.HashFS(themes.DefaultTemplates(), []string{".html", ".txt", ".md"})
	if err != nil {
		return fmt.Errorf("hash embedded default templates: %w", err)
	}
	templatesHash := buildcache.ContentHash("project:" + projectTemplatesHash + "\nembedded:" + embeddedTemplatesHash)
	if cache.SetTemplatesHash(templatesHash) {
		// Project or embedded templates changed - cache was invalidated.
		buildCacheLog.Phase("configure").Printf("Templates changed, full rebuild required")
	}

	// Reset stats for this build
	cache.ResetStats()

	if lifecycle.IsServeFullRebuild(m) {
		lifecycle.SetServeChangedPaths(m, nil)
		lifecycle.SetServeAffectedPaths(m, nil)
		return nil
	}

	return p.configureIncrementalServe(m, cache)
}

// Load discovers changed post inputs and expands their persisted dependency
// closure before transform and render plugins decide whether to restore caches.
func (p *BuildCachePlugin) Load(m *lifecycle.Manager) error {
	if !p.enabled || p.cache == nil {
		return nil
	}

	posts := m.Posts()
	batch := make([]struct{ Path, InputHash, Template string }, 0, len(posts))
	identitiesByPath := make(map[string][]string, len(posts))
	for _, post := range posts {
		// Empty slugs (the homepage) cannot participate in slug-based dependency
		// expansion, but they must still be marked affected when their input
		// changes, otherwise incremental/fast builds never render them.
		if post.Skip || post.Path == "" || post.InputHash == "" {
			continue
		}
		batch = append(batch, struct{ Path, InputHash, Template string }{
			Path:      post.Path,
			InputHash: post.InputHash,
			Template:  post.Template,
		})
		identitiesByPath[post.Path] = postDependencyIdentities(post)
	}

	changedPaths := p.cache.ShouldRebuildBatch(batch)
	affected := lifecycle.GetServeAffectedPaths(m)
	if affected == nil {
		affected = make(map[string]bool, len(changedPaths))
	}
	for path := range changedPaths {
		// The previous slug must be invalidated before transform refreshes the
		// graph's path-to-slug entry, otherwise inbound links keep old routes.
		if oldSlug := p.cache.Graph.SlugForPath(path); oldSlug != "" {
			p.cache.MarkSlugChanged(oldSlug)
		}
		for _, identity := range identitiesByPath[path] {
			p.cache.MarkSlugChanged(identity)
		}
		affected[path] = true
	}

	changedSlugs := p.cache.GetChangedSlugs()
	dependentPaths := p.cache.GetAffectedPosts(changedSlugs)
	p.cache.MarkAffectedDependents(changedSlugs)
	for _, path := range dependentPaths {
		affected[path] = true
	}
	if len(affected) > 0 {
		lifecycle.SetServeAffectedPaths(m, affected)
	}

	if len(changedPaths) > 0 || len(dependentPaths) > 0 {
		buildCacheLog.Phase("load").Printf(
			"Invalidated %d changed posts and %d dependent posts",
			len(changedPaths),
			len(dependentPaths),
		)
	}

	return nil
}

func (p *BuildCachePlugin) isEnabled(config *lifecycle.Config) bool {
	if config == nil || config.Extra == nil {
		return true
	}
	cacheConfig, ok := config.Extra["build_cache"].(map[string]interface{})
	if !ok {
		return true
	}
	enabled, ok := cacheConfig["enabled"].(bool)
	if !ok {
		return true
	}
	return enabled
}

func configFilesHash(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	normalized := make([]string, 0, len(paths))
	seen := make(map[string]bool, len(paths))
	add := func(path string) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			absPath = filepath.Clean(path)
		}
		if !seen[absPath] {
			seen[absPath] = true
			normalized = append(normalized, absPath)
		}
	}
	for _, path := range paths {
		if path == "" {
			continue
		}
		add(path)
		// Included config files (include = [...]) affect output too, so an
		// edit to any of them must invalidate the cache.
		if included, err := config.DiscoverIncludedConfigPaths(path); err == nil {
			for _, inc := range included {
				add(inc)
			}
		}
	}

	sort.Strings(normalized)

	hashes := make([]string, 0, len(normalized))
	for _, path := range normalized {
		if hash, err := buildcache.HashFile(path); err == nil && hash != "" {
			hashes = append(hashes, path+":"+hash)
		}
	}

	return strings.Join(hashes, "\n")
}

func configHashInput(config *lifecycle.Config, paths []string) string {
	components := make([]string, 0, 3)
	if pathHash := configFilesHash(paths); pathHash != "" {
		components = append(components, pathHash)
	}
	if config != nil {
		components = append(components, config.ContentDir, strings.Join(config.GlobPatterns, "\x00"))
		// Unsaved settings previewed by `markata-go serve` change the config
		// without changing any file.
		if overlay, ok := config.Extra["config_overlay"].(string); ok && overlay != "" {
			components = append(components, "overlay:"+overlay)
		}
	}
	return strings.Join(components, "\n")
}

func (p *BuildCachePlugin) configureIncrementalServe(m *lifecycle.Manager, cache *buildcache.Cache) error {
	changedPaths := lifecycle.GetServeChangedPaths(m)
	if len(changedPaths) == 0 {
		lifecycle.SetServeAffectedPaths(m, nil)
		return nil
	}

	cache.MarkChangedPaths(changedPaths)
	removedPaths := lifecycle.GetServeRemovedPaths(m)
	for _, path := range removedPaths {
		cache.Graph.RemoveSource(path)
	}

	changedSlugs := cache.GetChangedSlugs()
	cache.MarkAffectedDependents(changedSlugs)
	changedSlugs = cache.GetChangedSlugs()
	affectedPaths := cache.GetAffectedPosts(changedSlugs)

	affected := make(map[string]bool, len(changedPaths)+len(affectedPaths))
	for _, path := range changedPaths {
		affected[path] = true
	}
	for _, path := range affectedPaths {
		affected[path] = true
	}
	for _, path := range removedPaths {
		affected[path] = true
	}
	lifecycle.SetServeAffectedPaths(m, affected)
	return nil
}

// Cleanup saves the build cache to disk.
func (p *BuildCachePlugin) Cleanup(m *lifecycle.Manager) error {
	if !p.enabled || p.cache == nil {
		return nil
	}

	if extra := m.Config().Extra; extra != nil {
		if async, ok := extra["cache_cleanup_async"].(bool); ok && async {
			go func() {
				if err := p.cleanupCache(m); err != nil {
					buildCacheLog.Phase("cleanup").Errorf("async cleanup failed: %v", err)
				}
			}()
			return nil
		}
	}

	return p.cleanupCache(m)
}

func (p *BuildCachePlugin) cleanupCache(m *lifecycle.Manager) error {
	if p.cache == nil {
		return nil
	}

	// Remove stale entries (posts that no longer exist)
	if config := m.Config(); config.Extra != nil {
		if fast, ok := config.Extra["fast_mode"].(bool); ok && fast {
			return p.saveCache(m)
		}
	}
	posts := m.Posts()
	currentPaths := make(map[string]bool, len(posts))
	for _, post := range posts {
		currentPaths[post.Path] = true
	}
	removed := p.cache.RemoveStale(currentPaths)
	if removed > 0 {
		buildCacheLog.Phase("cleanup").Printf("Removed %d stale cache entries", removed)
	}
	removedMermaid, err := p.cache.CleanupMermaidSVG()
	if err != nil {
		return err
	}
	if removedMermaid > 0 {
		buildCacheLog.Phase("cleanup").Printf("Removed %d stale Mermaid SVG cache entries", removedMermaid)
	}

	// Save cache
	return p.saveCache(m)
}

func (p *BuildCachePlugin) saveCache(m *lifecycle.Manager) error {
	if p.cache == nil {
		return nil
	}
	if err := p.cache.Save(); err != nil {
		buildCacheLog.Phase("cleanup").Errorf("failed to save build cache: %v", err)
	}

	// Log stats
	skipped, rebuilt := p.cache.Stats()
	if skipped > 0 || rebuilt > 0 {
		buildCacheLog.Phase("cleanup").Printf("Incremental build: %d skipped, %d rebuilt", skipped, rebuilt)
		m.Cache().Set("build_cache_skipped", skipped)
		m.Cache().Set("build_cache_rebuilt", rebuilt)
	}

	// Log dependency graph stats
	graphSize := p.cache.GraphSize()
	if graphSize > 0 {
		buildCacheLog.Phase("cleanup").Printf("Dependency graph: %d posts with dependencies tracked", graphSize)
	}

	return nil
}

// Transform collects dependencies from posts after wikilinks/embeds have processed them.
// This runs late in the transform stage to ensure all dependencies have been collected.
func (p *BuildCachePlugin) Transform(m *lifecycle.Manager) error {
	if !p.enabled || p.cache == nil {
		return nil
	}

	posts := m.Posts()
	depsRecorded := 0

	for _, post := range posts {
		if post.Skip {
			continue
		}
		dependencies := append([]string(nil), post.Dependencies...)
		dependencies = append(dependencies, unresolvedLogicalDependencies(post.Content)...)
		if len(dependencies) > 0 {
			p.cache.SetDependencies(post.Path, post.Slug, dependencies)
			depsRecorded++
		}
	}

	if depsRecorded > 0 {
		buildCacheLog.Phase("transform").Printf("Recorded dependencies for %d posts", depsRecorded)
	}

	return nil
}

// Cache returns the build cache instance.
func (p *BuildCachePlugin) Cache() *buildcache.Cache {
	return p.cache
}

// Enabled returns whether the build cache is enabled.
func (p *BuildCachePlugin) Enabled() bool {
	return p.enabled && p.cache != nil
}

// GetBuildCache retrieves the build cache from the manager's cache.
// Returns nil if not found or caching is disabled.
func GetBuildCache(m *lifecycle.Manager) *buildcache.Cache {
	if cached, ok := m.Cache().Get("build_cache"); ok {
		if bc, ok := cached.(*buildcache.Cache); ok {
			return bc
		}
	}
	return nil
}

// Ensure BuildCachePlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*BuildCachePlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*BuildCachePlugin)(nil)
	_ lifecycle.LoadPlugin      = (*BuildCachePlugin)(nil)
	_ lifecycle.TransformPlugin = (*BuildCachePlugin)(nil)
	_ lifecycle.CleanupPlugin   = (*BuildCachePlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*BuildCachePlugin)(nil)
)
