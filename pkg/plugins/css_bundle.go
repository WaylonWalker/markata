// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// CSSBundlePlugin combines multiple CSS files into optimized bundles.
// It runs during the Write stage, after other CSS-generating plugins
// (palette_css, chroma_css) have completed. This reduces HTTP requests
// and improves page load performance.
//
// The plugin:
// - Discovers CSS files in the output directory
// - Matches files to bundle configurations
// - Respects CSS load order (files are concatenated in config order)
// - Writes bundled output with optional source comments
// - Stores bundle paths in the cache for template access
type CSSBundlePlugin struct {
	config  models.CSSBundleConfig
	bundles map[string]string // bundle name -> output path
	exclude map[string]bool   // excluded file patterns
}

// NewCSSBundlePlugin creates a new CSSBundlePlugin.
func NewCSSBundlePlugin() *CSSBundlePlugin {
	return &CSSBundlePlugin{
		bundles: make(map[string]string),
		exclude: make(map[string]bool),
	}
}

// Name returns the unique name of the plugin.
func (p *CSSBundlePlugin) Name() string {
	return "css_bundle"
}

// Configure reads the CSS bundle configuration from the manager's config.
func (p *CSSBundlePlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Try to get css_bundle config from Extra
	// It may be a models.CSSBundleConfig or a map[string]interface{} from TOML parsing
	switch v := config.Extra["css_bundle"].(type) {
	case models.CSSBundleConfig:
		p.config = v
	case map[string]interface{}:
		p.config = p.parseConfigFromMap(v)
	}

	// Build exclude map for fast lookups
	for _, pattern := range p.config.Exclude {
		p.exclude[pattern] = true
	}

	return nil
}

// parseConfigFromMap parses CSSBundleConfig from a raw map (TOML parsing result).
func (p *CSSBundlePlugin) parseConfigFromMap(m map[string]interface{}) models.CSSBundleConfig {
	cfg := models.NewCSSBundleConfig()

	if enabled, ok := m["enabled"].(bool); ok {
		cfg.Enabled = enabled
	}

	if minify, ok := m["minify"].(bool); ok {
		cfg.Minify = minify
	}

	if addComments, ok := m["add_source_comments"].(bool); ok {
		cfg.AddSourceComments = &addComments
	}

	// Parse page selection.
	switch pages := m["pages"].(type) {
	case []string:
		cfg.Pages = pages
	case []interface{}:
		for _, page := range pages {
			if path, ok := page.(string); ok {
				cfg.Pages = append(cfg.Pages, path)
			}
		}
	}

	// Parse exclude list
	if exclude, ok := m["exclude"].([]interface{}); ok {
		cfg.Exclude = make([]string, 0, len(exclude))
		for _, e := range exclude {
			if s, ok := e.(string); ok {
				cfg.Exclude = append(cfg.Exclude, s)
			}
		}
	}

	// Parse bundles
	switch bundles := m["bundles"].(type) {
	case []map[string]interface{}:
		cfg.Bundles = p.parseBundlesFromSlice(bundles)
	case []interface{}:
		cfg.Bundles = p.parseBundlesFromInterface(bundles)
	}

	return cfg
}

// parseBundlesFromSlice parses bundle configs from a slice of maps.
func (p *CSSBundlePlugin) parseBundlesFromSlice(bundles []map[string]interface{}) []models.BundleConfig {
	result := make([]models.BundleConfig, 0, len(bundles))
	for _, b := range bundles {
		bundle := p.parseSingleBundle(b)
		if bundle.Name != "" {
			result = append(result, bundle)
		}
	}
	return result
}

// parseBundlesFromInterface parses bundle configs from a slice of interface{}.
func (p *CSSBundlePlugin) parseBundlesFromInterface(bundles []interface{}) []models.BundleConfig {
	result := make([]models.BundleConfig, 0, len(bundles))
	for _, b := range bundles {
		if bMap, ok := b.(map[string]interface{}); ok {
			bundle := p.parseSingleBundle(bMap)
			if bundle.Name != "" {
				result = append(result, bundle)
			}
		}
	}
	return result
}

// parseSingleBundle parses a single bundle config from a map.
func (p *CSSBundlePlugin) parseSingleBundle(m map[string]interface{}) models.BundleConfig {
	bundle := models.BundleConfig{}

	if name, ok := m["name"].(string); ok {
		bundle.Name = name
	}

	if output, ok := m["output"].(string); ok {
		bundle.Output = output
	}

	// Parse sources
	switch sources := m["sources"].(type) {
	case []interface{}:
		bundle.Sources = make([]string, 0, len(sources))
		for _, s := range sources {
			if str, ok := s.(string); ok {
				bundle.Sources = append(bundle.Sources, str)
			}
		}
	case []string:
		bundle.Sources = sources
	}

	return bundle
}

// Write performs CSS bundling in the output directory.
func (p *CSSBundlePlugin) Write(m *lifecycle.Manager) error {
	// Skip if not enabled or no bundles configured
	if !p.config.Enabled || len(p.config.Bundles) == 0 {
		return nil
	}

	config := m.Config()
	outputDir := config.OutputDir

	log.Printf("[css_bundle] Starting CSS bundling with %d bundle(s)", len(p.config.Bundles))

	// Process each bundle configuration
	for _, bundleConfig := range p.config.Bundles {
		if err := p.processBundle(bundleConfig, outputDir); err != nil {
			return fmt.Errorf("processing bundle %q: %w", bundleConfig.Name, err)
		}
	}

	// Store bundle paths in cache for template access
	p.storeBundlePaths(m)

	log.Printf("[css_bundle] Completed CSS bundling: %d bundle(s) created", len(p.bundles))

	return nil
}

// processBundle creates a single CSS bundle from the configuration.
func (p *CSSBundlePlugin) processBundle(bundleConfig models.BundleConfig, outputDir string) error {
	if bundleConfig.Name == "" {
		return fmt.Errorf("bundle name is required")
	}
	if bundleConfig.Output == "" {
		return fmt.Errorf("bundle output path is required")
	}
	if len(bundleConfig.Sources) == 0 {
		return fmt.Errorf("bundle sources cannot be empty")
	}

	var buf bytes.Buffer

	// Write bundle header
	buf.WriteString(fmt.Sprintf("/* CSS Bundle: %s */\n", bundleConfig.Name))
	buf.WriteString("/* Generated by markata-go css_bundle plugin */\n")
	buf.WriteString(fmt.Sprintf("/* Sources: %d file(s) */\n\n", len(bundleConfig.Sources)))

	filesIncluded := 0
	totalBytes := 0

	// Process each source pattern
	for _, source := range bundleConfig.Sources {
		files, err := p.resolveSourcePattern(source, outputDir)
		if err != nil {
			log.Printf("[css_bundle] Warning: failed to resolve pattern %q: %v", source, err)
			continue
		}

		// Sort files for deterministic output
		sort.Strings(files)

		for _, file := range files {
			// Check if file should be excluded
			if p.isExcluded(file) {
				log.Printf("[css_bundle] Skipping excluded file: %s", file)
				continue
			}

			content, err := os.ReadFile(file)
			if err != nil {
				log.Printf("[css_bundle] Warning: failed to read %q: %v", file, err)
				continue
			}

			// Add source comment if enabled
			if p.config.IsAddSourceComments() {
				relPath, err := filepath.Rel(outputDir, file)
				if err != nil || relPath == "" {
					relPath = filepath.Base(file)
				}
				// Always use forward slashes in source comments (web convention)
				relPath = filepath.ToSlash(relPath)
				buf.WriteString(fmt.Sprintf("/* === Source: %s === */\n", relPath))
			}

			// Write the CSS content
			buf.Write(content)

			// Ensure newlines between files
			if len(content) > 0 && content[len(content)-1] != '\n' {
				buf.WriteByte('\n')
			}
			buf.WriteByte('\n')

			filesIncluded++
			totalBytes += len(content)
		}
	}

	if filesIncluded == 0 {
		log.Printf("[css_bundle] Warning: no files found for bundle %q", bundleConfig.Name)
		return nil
	}

	// Write the bundle file
	outputPath := filepath.Join(outputDir, bundleConfig.Output)
	outputDirPath := filepath.Dir(outputPath)

	if err := os.MkdirAll(outputDirPath, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	//nolint:gosec // G306: CSS bundle files need 0644 for web serving
	if err := os.WriteFile(outputPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("writing bundle: %w", err)
	}

	// Record the bundle
	p.bundles[bundleConfig.Name] = "/" + bundleConfig.Output

	log.Printf("[css_bundle] Created bundle %q: %d files, %d bytes -> %s",
		bundleConfig.Name, filesIncluded, totalBytes, bundleConfig.Output)

	return nil
}

// resolveSourcePattern resolves a source pattern to a list of files.
// Supports both direct file paths and glob patterns.
func (p *CSSBundlePlugin) resolveSourcePattern(pattern, outputDir string) ([]string, error) {
	// Remove leading / if present (absolute from output dir)
	pattern = strings.TrimPrefix(pattern, "/")

	fullPattern := filepath.Join(outputDir, pattern)

	// Check if it's a direct file path
	if !strings.ContainsAny(pattern, "*?[") {
		if _, err := os.Stat(fullPattern); err == nil {
			return []string{fullPattern}, nil
		}
		// File doesn't exist, return empty
		return nil, nil
	}

	// It's a glob pattern
	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return nil, fmt.Errorf("glob pattern error: %w", err)
	}

	return matches, nil
}

// isExcluded checks if a file should be excluded from bundling.
func (p *CSSBundlePlugin) isExcluded(filePath string) bool {
	filename := filepath.Base(filePath)

	// Check exact match
	if p.exclude[filename] {
		return true
	}

	// Check pattern match
	for pattern := range p.exclude {
		if strings.ContainsAny(pattern, "*?[") {
			matched, err := filepath.Match(pattern, filename)
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}

// storeBundlePaths stores bundle information in the cache for template access.
func (p *CSSBundlePlugin) storeBundlePaths(m *lifecycle.Manager) {
	cache := m.Cache()

	// Store individual bundle paths
	cache.Set("css_bundles", p.bundles)

	// Store as a list for template iteration
	bundleList := make([]map[string]string, 0, len(p.bundles))
	for name, path := range p.bundles {
		bundleList = append(bundleList, map[string]string{
			"name": name,
			"path": path,
		})
	}
	cache.Set("css_bundle_list", bundleList)

	// Store enabled flag
	cache.Set("css_bundling_enabled", p.config.Enabled)
}

// Priority returns the plugin priority for the write stage.
// Should run after other CSS-generating plugins (palette_css, chroma_css)
// so all CSS files are available for bundling.
func (p *CSSBundlePlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCleanup {
		return lifecycle.PriorityLast + 10 // After CSS purge and minification.
	}
	if stage == lifecycle.StageWrite {
		// Run late in Write stage, after CSS generators
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Ensure CSSBundlePlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*CSSBundlePlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*CSSBundlePlugin)(nil)
	_ lifecycle.WritePlugin     = (*CSSBundlePlugin)(nil)
	_ lifecycle.CleanupPlugin   = (*CSSBundlePlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*CSSBundlePlugin)(nil)
)
