// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
)

var pagefindLog = logging.Component("pagefind").Phase("cleanup")

// defaultBundleDir is the default directory name for Pagefind search index files.
const defaultBundleDir = "_pagefind"

// PagefindPlugin runs Pagefind to generate a search index after all HTML files are written.
// Pagefind is a static site search tool that creates an optimized search index
// that can be queried entirely client-side.
//
// This plugin runs in the Cleanup stage with PriorityLast to ensure all HTML
// files have been written before indexing begins.
//
// For more information about Pagefind, see: https://pagefind.app/
type PagefindPlugin struct{}

// NewPagefindPlugin creates a new PagefindPlugin.
func NewPagefindPlugin() *PagefindPlugin {
	return &PagefindPlugin{}
}

// Name returns the unique name of the plugin.
func (p *PagefindPlugin) Name() string {
	return pagefindBinaryName
}

// Cleanup runs Pagefind to index the generated site.
// This runs after all HTML files have been written in the Write stage.
// On incremental builds with few changes, Pagefind is skipped for speed.
func (p *PagefindPlugin) Cleanup(m *lifecycle.Manager) error {
	config := m.Config()
	searchConfig := getSearchConfig(config)

	// Skip if search is disabled
	if !searchConfig.IsEnabled() {
		return nil
	}

	// Skip when bleve is the active search backend.
	if searchConfig.SearchBackend() == "bleve" {
		return nil
	}

	// Skip in fast mode - search indexing is unnecessary during development
	if fast, ok := config.Extra["fast_mode"].(bool); ok && fast {
		return nil
	}

	cache := GetBuildCache(m)
	corpusHash := p.computeCorpusHash(m, searchConfig)
	if cache != nil && corpusHash != "" && cache.GetPagefindCorpusHash() == corpusHash && p.indexExists(config, searchConfig) {
		pagefindLog.Printf("Skipping search index update (search corpus unchanged)")
		return nil
	}

	verbose := searchConfig.Pagefind.IsVerbose()

	// Try to find or install Pagefind
	pagefindPath, err := p.findOrInstallPagefind(searchConfig, verbose)
	if err != nil {
		// Log warning but don't fail the build
		pagefindLog.Warnf("%v", err)
		pagefindLog.Printf("The site will work fine, just without search functionality")
		return nil
	}

	if pagefindPath == "" {
		// No Pagefind available - already warned
		return nil
	}

	if err := p.runPagefind(pagefindPath, config, searchConfig, verbose); err != nil {
		return err
	}

	if cache != nil && corpusHash != "" {
		cache.SetPagefindCorpusHash(corpusHash)
		if err := cache.Save(); err != nil {
			pagefindLog.Warnf("failed to persist pagefind corpus hash: %v", err)
		}
	}

	return nil
}

// findOrInstallPagefind locates or automatically installs the Pagefind binary.
// It first checks the system PATH, then attempts auto-install if enabled.
func (p *PagefindPlugin) findOrInstallPagefind(searchConfig models.SearchConfig, verbose bool) (string, error) {
	// First, check if pagefind is in PATH
	pagefindPath, err := exec.LookPath("pagefind")
	if err == nil {
		return pagefindPath, nil
	}

	// Check if auto-install is enabled
	if !searchConfig.Pagefind.IsAutoInstallEnabled() {
		pagefindLog.Warnf("pagefind not found in PATH, skipping search index generation")
		pagefindLog.Printf("Install with: npm install -g pagefind OR cargo install pagefind")
		pagefindLog.Printf("Or enable auto_install in config: [search.pagefind] auto_install = true")
		return "", nil
	}

	// Check if we have a cached version for the requested version
	installer := NewPagefindInstallerWithConfig(PagefindInstallerConfig{
		Version:  searchConfig.Pagefind.Version,
		CacheDir: searchConfig.Pagefind.CacheDir,
	})
	installer.Verbose = verbose

	// Attempt to install
	if verbose {
		pagefindLog.Printf("Pagefind not found in PATH, attempting auto-install...")
	}
	installedPath, err := installer.Install()
	if err != nil {
		return "", fmt.Errorf("auto-install failed: %w", err)
	}

	return installedPath, nil
}

// runPagefind executes the Pagefind CLI to generate the search index.
func (p *PagefindPlugin) runPagefind(pagefindPath string, config *lifecycle.Config, searchConfig models.SearchConfig, verbose bool) error {
	outputDir := config.OutputDir

	// Verify output directory exists
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return fmt.Errorf("output directory does not exist: %s", outputDir)
	}

	// Build command arguments
	args := []string{
		"--site", outputDir,
	}

	// Configure output subdirectory (where search index files go)
	// Note: Using --output-subdir instead of deprecated --bundle-dir (deprecated in Pagefind 1.0)
	bundleDir := searchConfig.Pagefind.BundleDir
	if bundleDir == "" {
		bundleDir = defaultBundleDir
	}
	args = append(args, "--output-subdir", bundleDir)

	// Configure root selector if specified
	if searchConfig.Pagefind.RootSelector != "" {
		args = append(args, "--root-selector", searchConfig.Pagefind.RootSelector)
	}

	// Configure exclude selectors
	for _, selector := range searchConfig.Pagefind.ExcludeSelectors {
		args = append(args, "--exclude-selectors", selector)
	}

	// Run pagefind
	if verbose {
		pagefindLog.Printf("Generating search index in %s/%s", outputDir, bundleDir)
	}
	cmd := exec.Command(pagefindPath, args...)
	cmd.Dir = "." // Run from project root
	start := time.Now()
	pagefindLog.Printf("Running subprocess: %s %s", pagefindPath, strings.Join(args, " "))

	var stdoutLog, stderrLog *subprocessLineWriter
	var stderrTail *subprocessTailWriter
	cmd.Stdout = io.Discard
	if verbose {
		stdoutLog = newSubprocessLineWriter(pagefindLog)
		stderrLog = newSubprocessLineWriter(pagefindLog.Level("error"))
		cmd.Stdout, cmd.Stderr = stdoutLog, stderrLog
	} else {
		stderrTail = newSubprocessTailWriter(64 * 1024)
		cmd.Stderr = stderrTail
	}
	if err := cmd.Run(); err != nil {
		if stdoutLog != nil {
			stdoutLog.Flush()
			stderrLog.Flush()
		} else if output := stderrTail.String(); output != "" {
			logSubprocessOutput(pagefindLog.Level("error"), output)
		}
		return fmt.Errorf("pagefind indexing failed: %w", err)
	}
	if stdoutLog != nil {
		stdoutLog.Flush()
		stderrLog.Flush()
	}
	pagefindLog.Printf("Subprocess completed in %v", time.Since(start))

	// Verify the index was created
	indexPath := filepath.Join(outputDir, bundleDir, "pagefind.js")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return fmt.Errorf("pagefind did not create expected index file: %s", indexPath)
	}

	pagefindLog.Printf("Search index generated successfully")
	return nil
}

// Priority returns the plugin priority for the cleanup stage.
// Pagefind runs last in cleanup to ensure all HTML files are written first.
func (p *PagefindPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCleanup {
		return lifecycle.PriorityLast
	}
	return lifecycle.PriorityDefault
}

// getSearchConfig extracts SearchConfig from lifecycle.Config.Extra.
func getSearchConfig(config *lifecycle.Config) models.SearchConfig {
	if config.Extra != nil {
		if sc, ok := config.Extra["search"].(models.SearchConfig); ok {
			return sc
		}
	}
	return models.NewSearchConfig()
}

func (p *PagefindPlugin) computeCorpusHash(m *lifecycle.Manager, searchConfig models.SearchConfig) string {
	posts := m.Posts()
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Path < posts[j].Path
	})

	var builder strings.Builder
	builder.WriteString(searchConfig.Position)
	builder.WriteString("\n")
	builder.WriteString(searchConfig.Placeholder)
	builder.WriteString("\n")
	builder.WriteString(searchConfig.Pagefind.BundleDir)
	builder.WriteString("\n")
	builder.WriteString(searchConfig.Pagefind.RootSelector)
	for _, selector := range searchConfig.Pagefind.ExcludeSelectors {
		builder.WriteString("\n")
		builder.WriteString(selector)
	}

	for _, post := range posts {
		if post == nil || post.Skip || post.Draft || !post.Published || post.Private || isSearchExcluded(post) {
			continue
		}
		builder.WriteString("\n--post:")
		builder.WriteString(post.Path)
		builder.WriteString("\n--href:")
		builder.WriteString(post.Href)
		builder.WriteString("\n--slug:")
		builder.WriteString(post.Slug)
		builder.WriteString("\n--title:")
		builder.WriteString(postTitle(post))
		builder.WriteString("\n--description:")
		builder.WriteString(postDescription(post))
		builder.WriteString("\n--article-hash:")
		builder.WriteString(buildcache.ContentHash(post.ArticleHTML))
		for _, tag := range post.Tags {
			builder.WriteString("\n--tag:")
			builder.WriteString(tag)
		}
	}

	return buildcache.ContentHash(builder.String())
}

func (p *PagefindPlugin) indexExists(config *lifecycle.Config, searchConfig models.SearchConfig) bool {
	bundleDir := searchConfig.Pagefind.BundleDir
	if bundleDir == "" {
		bundleDir = defaultBundleDir
	}
	indexPath := filepath.Join(config.OutputDir, bundleDir, "pagefind.js")
	_, err := os.Stat(indexPath)
	return err == nil
}

func isSearchExcluded(post *models.Post) bool {
	if post == nil || post.Extra == nil {
		return false
	}
	if value, ok := post.Extra["search_exclude"].(bool); ok {
		return value
	}
	return false
}

func postTitle(post *models.Post) string {
	if post == nil || post.PlainTitle() == "" {
		return ""
	}
	return post.PlainTitle()
}

func postDescription(post *models.Post) string {
	if post == nil || post.Description == nil {
		return ""
	}
	return *post.Description
}

// Ensure PagefindPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*PagefindPlugin)(nil)
	_ lifecycle.CleanupPlugin  = (*PagefindPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*PagefindPlugin)(nil)
)
