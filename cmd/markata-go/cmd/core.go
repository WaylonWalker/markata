package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/servecontrol"
)

// createManager creates and configures a lifecycle manager with all plugins.
func createManager(cfgPath string) (*lifecycle.Manager, error) {
	return createManagerWithPlugins(cfgPath, plugins.DefaultPlugins)
}

// createManagerWithPlugins loads configuration and registers the plugin set
// returned by pluginSet.
func createManagerWithPlugins(cfgPath string, pluginSet func() []lifecycle.Plugin) (*lifecycle.Manager, error) {
	cfg, configPathUsed, configPaths, err := loadManagerConfig(cfgPath)
	if err != nil {
		return nil, err
	}

	// Apply output directory override from CLI flag
	if outputDir != "" {
		cfg.OutputDir = outputDir
	}

	baseDir := resolveConfigBaseDir(configPathUsed)
	contentDir := "."
	if baseDir != "" {
		contentDir = baseDir
	}
	setSourcePathBaseDir(contentDir)
	cfg.OutputDir = resolveConfigRelativePath(baseDir, cfg.OutputDir)
	cfg.AssetsDir = resolveConfigRelativePath(baseDir, cfg.AssetsDir)
	cfg.TemplatesDir = resolveConfigRelativePath(baseDir, cfg.TemplatesDir)

	// Validate config
	validationErrs := config.ValidateConfig(cfg)
	actualErrors, warnings := config.SplitErrorsAndWarnings(validationErrs)

	// Print warnings
	for _, w := range warnings {
		if runtime, jobID, _ := currentServeControl(); runtime != nil {
			runtime.AddDiagnostic(servecontrol.Diagnostic{
				Code: "config.validation_warning", Severity: "warning", Message: w.Error(),
				JobID: jobID, SuggestedFix: "Review this setting in markata-go.toml, then rebuild.",
			})
		}
		if verbose || isLicenseWarning(w) {
			warnf("%v", w)
		}
	}

	// Return errors
	if len(actualErrors) > 0 {
		return nil, fmt.Errorf("config validation failed: %w", actualErrors[0])
	}

	// Create manager
	m := lifecycle.NewManager()

	// Convert models.Config to lifecycle.Config
	lcConfig := &lifecycle.Config{
		ContentDir:   contentDir,
		OutputDir:    cfg.OutputDir,
		GlobPatterns: cfg.GlobConfig.Patterns,
		Extra:        make(map[string]interface{}),
	}

	// Copy config values to Extra for plugins to access
	lcConfig.Extra["url"] = cfg.URL
	lcConfig.Extra["title"] = cfg.Title
	lcConfig.Extra["description"] = cfg.Description
	lcConfig.Extra["author"] = cfg.Author
	lcConfig.Extra["language"] = cfg.Language
	lcConfig.Extra["author_url"] = cfg.AuthorURL
	lcConfig.Extra["managing_editor"] = cfg.ManagingEditor
	lcConfig.Extra["webmaster"] = cfg.WebMaster
	lcConfig.Extra["copyright"] = cfg.Copyright
	lcConfig.Extra["templates_dir"] = cfg.TemplatesDir
	lcConfig.Extra["assets_dir"] = cfg.AssetsDir
	lcConfig.Extra["template_presets"] = cfg.TemplatePresets
	lcConfig.Extra["default_templates"] = cfg.DefaultTemplates
	lcConfig.Extra["theme_calendar"] = cfg.ThemeCalendar
	lcConfig.Extra["error_pages"] = cfg.ErrorPages
	lcConfig.Extra["resource_hints"] = cfg.ResourceHints
	lcConfig.Extra["markdown"] = cfg.MarkdownConfig
	lcConfig.Extra["feeds"] = cfg.Feeds
	lcConfig.Extra["feed_defaults"] = cfg.FeedDefaults
	lcConfig.Extra["use_gitignore"] = cfg.GlobConfig.UseGitignore
	lcConfig.Extra["nav"] = cfg.Nav
	lcConfig.Extra["footer"] = cfg.Footer
	lcConfig.Extra["post_formats"] = cfg.PostFormats
	lcConfig.Extra["templates"] = cfg.Templates
	lcConfig.Extra["websub"] = cfg.WebSub
	lcConfig.Extra["well_known"] = cfg.WellKnown
	lcConfig.Extra["seo"] = cfg.SEO
	lcConfig.Extra["search"] = cfg.Search
	lcConfig.Extra["components"] = cfg.Components
	lcConfig.Extra["header"] = cfg.Header
	lcConfig.Extra["head"] = cfg.Head
	lcConfig.Extra["toc"] = cfg.Toc
	lcConfig.Extra["sidebar"] = cfg.Sidebar
	fontpack := cfg.Fontpack
	if cfg.Theme.Fontpack != "" {
		fontpack = cfg.Theme.Fontpack
	}
	lcConfig.Extra["fontpack"] = fontpack
	lcConfig.Extra["fontpacks_file"] = resolveConfigRelativePath(baseDir, cfg.FontpacksFile)

	// Pass theme configuration to plugins
	lcConfig.Extra["theme"] = cfg.Theme

	// Pass layout configuration for automatic layout selection
	lcConfig.Extra["layout"] = &cfg.Layout

	// Pass blogroll configuration
	lcConfig.Extra["blogroll"] = cfg.Blogroll

	// Pass mentions configuration
	lcConfig.Extra["mentions"] = cfg.Mentions

	// Pass assets configuration for self-hosting CDN assets
	lcConfig.Extra["assets"] = cfg.Assets

	// Pass sidebar, toc, and header configurations
	lcConfig.Extra["sidebar"] = cfg.Sidebar
	lcConfig.Extra["toc"] = cfg.Toc
	lcConfig.Extra["header"] = cfg.Header

	// Pass tags configuration
	lcConfig.Extra["tags"] = cfg.Tags
	lcConfig.Extra["feeds_page"] = cfg.FeedsPage

	// Pass view transitions configuration
	lcConfig.Extra["view_transitions"] = cfg.ViewTransitions

	// Pass garden configuration
	lcConfig.Extra["garden"] = cfg.Garden

	// Pass image library configuration
	lcConfig.Extra["images"] = cfg.Images

	// Pass search configuration with verbose flag override from CLI
	searchConfig := cfg.Search
	if verbose {
		// CLI --verbose flag overrides config setting
		v := true
		searchConfig.Pagefind.Verbose = &v
	}
	lcConfig.Extra["search"] = searchConfig

	// Copy arbitrary plugin configs from cfg.Extra (e.g., image_zoom, wikilinks)
	if cfg.Extra != nil {
		for key, value := range cfg.Extra {
			lcConfig.Extra[key] = value
		}
	}
	// Build identity is supplied by the binary, not by site configuration.
	lcConfig.Extra["markata_version"] = Version
	lcConfig.Extra["markata_commit"] = Commit
	// Canonical nested theme values remain authoritative after arbitrary plugin
	// extras are copied. Legacy root keys are compatibility inputs only.
	lcConfig.Extra["fontpack"] = fontpack
	lcConfig.Extra["theme"] = cfg.Theme
	// Keep resolved path-bearing settings authoritative after copying arbitrary
	// extras. Otherwise a raw "templates" value can replace the config-relative
	// absolute path and make the engine load the default theme tree instead of
	// the site's templates.
	lcConfig.Extra["templates_dir"] = cfg.TemplatesDir
	lcConfig.Extra["assets_dir"] = cfg.AssetsDir
	// Keep path-bearing plugin configuration aligned with the resolved build
	// paths even when the raw config also contains the option in Extra.
	lcConfig.Extra["fontpacks_file"] = resolveConfigRelativePath(baseDir, cfg.FontpacksFile)

	// Previews build into their own cache so the site's warm cache survives
	// previewing, resetting, and stopping serve.
	if _, previewing := lcConfig.Extra[configOverlayExtraKey]; previewing {
		cacheDir := filepath.Join(contentDir, ".markata")
		if dir, ok := lcConfig.Extra["cache_dir"].(string); ok && dir != "" {
			cacheDir = dir
		}
		lcConfig.Extra["cache_dir"] = filepath.Join(cacheDir, servePreviewCacheDir)
	}

	// Store full models.Config for components that need direct access (e.g., 404 page handler)
	lcConfig.Extra["models_config"] = cfg
	if configPathUsed != "" {
		lcConfig.Extra["config_path"] = configPathUsed
	}
	if len(configPaths) > 0 {
		lcConfig.Extra["config_paths"] = configPaths
	}

	m.SetConfig(lcConfig)

	// Set concurrency if specified
	if cfg.Concurrency > 0 {
		m.SetConcurrency(cfg.Concurrency)
	}

	registerConfiguredPlugins(m, cfg, pluginSet())

	return m, nil
}

// registerConfiguredPlugins always loads the standard set and appends any
// registered optional hooks named in the site configuration.
func registerConfiguredPlugins(m *lifecycle.Manager, cfg *models.Config, defaults []lifecycle.Plugin) {
	registered := make(map[string]struct{}, len(defaults))
	for _, plugin := range defaults {
		m.RegisterPlugin(plugin)
		registered[plugin.Name()] = struct{}{}
	}
	for _, name := range cfg.Hooks {
		if name == sourceDefault || !cfg.IsHookEnabled(name) {
			continue
		}
		if _, ok := registered[name]; ok {
			continue
		}
		plugin, ok := plugins.PluginByName(name)
		if !ok {
			continue
		}
		m.RegisterPlugin(plugin)
		registered[name] = struct{}{}
	}
}

// createSinglePageManager configures the normal renderer for one Markdown
// source and publishes it at output/index.html. Site-level output such as
// feeds, listings, sitemaps, and search indexes is not generated.
func createSinglePageManager(cfgPath, sourcePath string) (*lifecycle.Manager, error) {
	m, err := createManagerWithPlugins(cfgPath, plugins.SinglePagePlugins)
	if err != nil {
		return nil, err
	}

	contentRoot, err := filepath.Abs(m.Config().ContentDir)
	if err != nil {
		return nil, fmt.Errorf("resolve content directory: %w", err)
	}
	sourceAbs, err := filepath.Abs(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("resolve Markdown file %q: %w", sourcePath, err)
	}
	relativePath, err := filepath.Rel(contentRoot, sourceAbs)
	if err != nil {
		return nil, fmt.Errorf("resolve Markdown file %q relative to content directory: %w", sourcePath, err)
	}
	if relativePath == ".." || strings.HasPrefix(relativePath, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("markdown file %q is outside content directory %q", sourcePath, m.Config().ContentDir)
	}
	info, err := os.Stat(sourceAbs)
	if err != nil {
		return nil, fmt.Errorf("markdown file %q: %w", sourcePath, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("markdown file %q is a directory", sourcePath)
	}
	extension := strings.ToLower(filepath.Ext(sourceAbs))
	if extension != ".md" && extension != ".markdown" {
		return nil, fmt.Errorf("markdown file %q must end in .md or .markdown", sourcePath)
	}

	m.Config().GlobPatterns = []string{relativePath}
	m.Config().Extra["feeds"] = []models.FeedConfig{}
	m.Config().Extra["subscription_feeds_disabled"] = true
	m.Config().Extra["single_page"] = true
	// The single-page plugin set has no search index, so hide the search UI.
	search, ok := m.Config().Extra["search"].(models.SearchConfig)
	if !ok {
		search = models.NewSearchConfig()
	}
	searchDisabled := false
	search.Enabled = &searchDisabled
	m.Config().Extra["search"] = search

	return m, nil
}

// configOverlayExtraKey carries the serve settings preview fingerprint to the
// build cache (see plugins.configHashInput).
const configOverlayExtraKey = "config_overlay"

// servePreviewCacheDir is the build cache subdirectory used while unsaved
// settings are previewed, so previews never overwrite the site's cache.
const servePreviewCacheDir = "serve-preview"

// configPreview holds unsaved settings previewed by `markata-go serve`:
// overlay sets keys and remove resets keys to their defaults.
type configPreview struct {
	overlay map[string]any
	remove  [][]string
}

func (p configPreview) empty() bool {
	return len(p.overlay) == 0 && len(p.remove) == 0
}

func loadManagerConfig(cfgPath string) (cfg *models.Config, configPathUsed string, configPaths []string, err error) {
	return loadManagerConfigWith(cfgPath, servePreviewConfig())
}

// loadManagerConfigWith loads the config files and applies preview (unsaved
// settings previewed by `markata-go serve`) on top, below env overrides.
func loadManagerConfigWith(cfgPath string, preview configPreview) (cfg *models.Config, configPathUsed string, configPaths []string, err error) {
	configPathUsed = cfgPath

	if len(mergeConfigFiles) > 0 || !preview.empty() {
		basePath := cfgPath
		if basePath == "" {
			discovered, discoverErr := config.Discover()
			if discoverErr == nil {
				basePath = discovered
			}
		}
		configPathUsed = basePath

		options := config.LoadOptions{Overlay: preview.overlay, Remove: preview.remove}
		cfg, err = config.LoadWithMergeOptions(options, basePath, mergeConfigFiles...)
		if err != nil {
			return nil, "", nil, fmt.Errorf("loading merged config: %w", err)
		}
		if !preview.empty() {
			// The build cache hashes config files; the preview changes the
			// config without touching them, so it must be part of the hash.
			fingerprint, err := json.Marshal(map[string]any{"set": preview.overlay, "reset": preview.remove})
			if err != nil {
				return nil, "", nil, fmt.Errorf("encoding settings preview: %w", err)
			}
			if cfg.Extra == nil {
				cfg.Extra = map[string]any{}
			}
			cfg.Extra[configOverlayExtraKey] = string(fingerprint)
		}
		if basePath != "" {
			configPaths = append(configPaths, basePath)
		}
		for _, path := range mergeConfigFiles {
			if path != "" {
				configPaths = append(configPaths, path)
			}
		}
		return cfg, configPathUsed, configPaths, nil
	}

	cfg, err = config.Load(cfgPath)
	if err != nil {
		return nil, "", nil, fmt.Errorf("loading config: %w", err)
	}
	if configPathUsed == "" {
		if discovered, discoverErr := config.Discover(); discoverErr == nil {
			configPathUsed = discovered
		}
	}
	if configPathUsed != "" {
		configPaths = append(configPaths, configPathUsed)
	}

	return cfg, configPathUsed, configPaths, nil
}

func resolveConfigBaseDir(configPath string) string {
	if configPath == "" {
		return ""
	}
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return filepath.Dir(configPath)
	}
	return filepath.Dir(absPath)
}

func resolveConfigRelativePath(baseDir, path string) string {
	if path == "" || baseDir == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(baseDir, path)
}

func licenseWarningMessage(cfg *models.Config) string {
	if cfg == nil || !cfg.NeedsLicenseWarning() {
		return ""
	}
	return fmt.Sprintf("License not configured. Set license = %q (recommended) or license = false.", models.DefaultLicenseKey)
}

func isLicenseWarning(err error) bool {
	var vErr config.ValidationError
	if !errors.As(err, &vErr) {
		return false
	}
	return vErr.IsWarn && vErr.Field == "license"
}

// applyFastMode sets the fast_mode flag in the manager's config Extra map,
// which tells plugins to skip expensive non-essential work (minification, CSS purging).
func applyFastMode(m *lifecycle.Manager) {
	if m.Config().Extra == nil {
		m.Config().Extra = make(map[string]any)
	}
	m.Config().Extra["fast_mode"] = true
	m.Config().Extra["blogroll_disabled"] = true
	m.Config().Extra["mentions_disabled"] = true
	m.Config().Extra["feeds_incremental"] = true
	verbosef("Fast mode: skipping minification, CSS purging, tailwind rebuilds, pagefind indexing, blogroll, and mentions")
}

// BuildResult holds the result of a build operation.
type BuildResult struct {
	Executor       lifecycle.BuildExecutor
	PostsProcessed int
	FeedsGenerated int
	FilesWritten   int
	Warnings       []string
	Duration       float64
	Benchmark      buildstats.Summary
	Content        diagnostics.ContentLedgerSnapshot

	// BlogrollStatus holds blogroll feature status
	BlogrollStatus BlogrollStatus
}

// BlogrollStatus holds information about the blogroll feature.
type BlogrollStatus struct {
	// Configured indicates if blogroll section exists in config
	Configured bool
	// Enabled indicates if blogroll is enabled
	Enabled bool
	// FeedsConfigured is the number of feeds configured
	FeedsConfigured int
	// FeedsFetched is the number of feeds successfully fetched
	FeedsFetched int
}

// runBuild executes a full build and returns the result.
func runBuild(m *lifecycle.Manager) (result *BuildResult, err error) {
	return runBuildObserved(m, nil)
}

// runBuildObserved reports lifecycle stage transitions to serve clients. The
// ordinary build command passes no observer and keeps its existing behavior.
func runBuildObserved(m *lifecycle.Manager, observe func(lifecycle.Stage, bool, error)) (result *BuildResult, err error) {
	if dagBuildEnabled() {
		return runDAGBuildObserved(m, observe)
	}
	if err := m.SetBuildExecutor(lifecycle.BuildExecutorLegacy); err != nil {
		return nil, fmt.Errorf("select legacy build executor: %w", err)
	}

	profile := buildstats.Start()
	defer func() {
		summary := profile.Stop()
		if result != nil {
			result.Benchmark = summary
		}
	}()

	// Run all lifecycle stages with verbose output if enabled
	stages := []lifecycle.Stage{
		lifecycle.StageConfigure,
		lifecycle.StageValidate,
		lifecycle.StageGlob,
		lifecycle.StageLoad,
		lifecycle.StageTransform,
		lifecycle.StageRender,
		lifecycle.StageCollect,
		lifecycle.StageWrite,
		lifecycle.StageCleanup,
	}

	for _, stage := range stages {
		stageStart := time.Now()
		if observe != nil {
			observe(stage, true, nil)
		}
		buildstats.SetActiveStage(string(stage))
		verbosef("  [%s] running...", stage)
		if err := m.RunTo(stage); err != nil {
			if observe != nil {
				observe(stage, false, err)
			}
			buildstats.SetActiveStage("")
			return nil, fmt.Errorf("stage %s: %w", stage, err)
		}
		if observe != nil {
			observe(stage, false, nil)
		}
		buildstats.SetActiveStage("")
		stageElapsed := time.Since(stageStart)
		buildstats.RecordStage(string(stage), stageElapsed)
		if verbose {
			verbosef("  [%s] done in %s", stage, stageElapsed.Truncate(100*time.Microsecond))
			switch stage {
			case lifecycle.StageGlob:
				verbosef("  [%s] discovered %d files", stage, len(m.Files()))
			case lifecycle.StageLoad:
				verbosef("  [%s] loaded %d posts", stage, len(m.Posts()))
			case lifecycle.StageCollect:
				verbosef("  [%s] collected %d feeds", stage, len(m.Feeds()))
			case lifecycle.StageConfigure, lifecycle.StageValidate, lifecycle.StageTransform,
				lifecycle.StageRender, lifecycle.StageWrite, lifecycle.StageCleanup:
				// No extra logging for these stages
			}
		}
	}

	// Collect results
	result = &BuildResult{
		Executor:       m.BuildExecutor(),
		PostsProcessed: len(m.Posts()),
		FeedsGenerated: len(m.Feeds()),
		Content:        m.ContentDiagnostics(),
	}

	// Collect blogroll status
	result.BlogrollStatus = getBlogrollStatus(m)

	// Collect warnings
	for _, w := range m.Warnings() {
		result.Warnings = append(result.Warnings, w.Error())
	}

	return result, nil
}

// getBlogrollStatus extracts blogroll feature status from the manager.
func getBlogrollStatus(m *lifecycle.Manager) BlogrollStatus {
	status := BlogrollStatus{}

	cfg := m.Config()
	if cfg == nil || cfg.Extra == nil {
		return status
	}

	blogrollVal, ok := cfg.Extra["blogroll"]
	if !ok {
		return status
	}

	blogrollConfig, ok := blogrollVal.(models.BlogrollConfig)
	if !ok {
		return status
	}

	status.Configured = true
	status.Enabled = blogrollConfig.Enabled
	status.FeedsConfigured = len(blogrollConfig.Feeds)

	// Get fetched feeds count from cache
	if feedsVal, ok := m.Cache().Get("blogroll_feeds"); ok {
		if feeds, ok := feedsVal.([]*models.ExternalFeed); ok {
			status.FeedsFetched = len(feeds)
		}
	}

	return status
}
