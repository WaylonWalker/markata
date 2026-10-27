package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/csspurge"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
)

var cssPurgeLog = logging.Component("css_purge").Phase("cleanup")

// CSSPurgePlugin removes unused CSS rules from stylesheets based on
// which selectors are actually used in the generated HTML files.
//
// This plugin runs in the Cleanup stage after all HTML has been written,
// allowing it to scan the final output and optimize CSS accordingly.
//
// The plugin supports preserving dynamically-added classes (like those
// from JavaScript frameworks) via glob patterns in the configuration.
type CSSPurgePlugin struct{}

// NewCSSPurgePlugin creates a new CSSPurgePlugin.
func NewCSSPurgePlugin() *CSSPurgePlugin {
	return &CSSPurgePlugin{}
}

// Name returns the unique name of the plugin.
func (p *CSSPurgePlugin) Name() string {
	return "css_purge"
}

// Cleanup scans HTML files and removes unused CSS rules.
// This runs after all content has been written to the output directory.
// Skipped in fast mode (--fast flag) for faster development builds.
func (p *CSSPurgePlugin) Cleanup(m *lifecycle.Manager) error {
	config := m.Config()

	// Skip in fast mode
	if fast, ok := config.Extra["fast_mode"].(bool); ok && fast {
		return nil
	}

	purgeConfig := getCSSPurgeConfig(config)

	// Skip if disabled
	if !purgeConfig.Enabled {
		return nil
	}

	outputDir := config.OutputDir
	verbose := purgeConfig.Verbose

	if verbose {
		cssPurgeLog.Printf("Analyzing CSS usage in %s", outputDir)
	}

	// Step 1: Find and scan HTML files
	cacheDir := filepath.Join(config.ContentDir, ".markata-cache", "css-purge")
	used, _, err := scanHTMLFilesForSelectors(outputDir, cacheDir, m.Concurrency(), verbose)
	if err != nil {
		return err
	}
	if used == nil {
		return nil // No HTML files found
	}

	// Step 2: Find CSS files
	cssFiles, err := findCSSFiles(outputDir)
	if err != nil {
		return fmt.Errorf("failed to find CSS files: %w", err)
	}

	if len(cssFiles) == 0 {
		if verbose {
			cssPurgeLog.Printf("No CSS files found, skipping")
		}
		return nil
	}

	// Step 3: Build purge options
	opts := buildPurgeOptions(purgeConfig, verbose)

	// Step 4: Process CSS files
	stats := processCSSFiles(cssFiles, outputDir, cacheDir, used, opts, purgeConfig, verbose)

	// Step 5: Report summary
	reportPurgeSummary(stats, purgeConfig, verbose)

	return nil
}

// purgeProcessingStats holds statistics from CSS processing.
type purgeProcessingStats struct {
	totalOriginal  int
	totalPurged    int
	filesProcessed int
	filesSkipped   int
	cacheHits      int
}

// scanHTMLFilesForSelectors finds and scans HTML files for used selectors.
func scanHTMLFilesForSelectors(outputDir, cacheDir string, concurrency int, verbose bool) (*csspurge.UsedSelectors, bool, error) {
	htmlFiles, err := findHTMLFiles(outputDir)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find HTML files: %w", err)
	}

	if len(htmlFiles) == 0 {
		if verbose {
			cssPurgeLog.Printf("No HTML files found, skipping")
		}
		return nil, false, nil
	}

	if verbose {
		cssPurgeLog.Printf("Found %d HTML files to analyze", len(htmlFiles))
	}

	fingerprint, scannedBytes, fingerprintOK := fingerprintHTMLFiles(htmlFiles, outputDir, concurrency)
	cachePath := filepath.Join(cacheDir, "html-selectors.json")
	if fingerprintOK {
		if cached, ok := loadCSSPurgeHTMLCache(cachePath); ok && cached.Schema == cssPurgeHTMLCacheSchema && cached.Fingerprint == fingerprint && cached.Selectors.Classes != nil && cached.Selectors.IDs != nil && cached.Selectors.Elements != nil && cached.Selectors.Attributes != nil {
			cssPurgeLog.Printf("HTML selector cache hit: files=%d bytes=%d", len(htmlFiles), scannedBytes)
			return &cached.Selectors, true, nil
		}
	}

	used, complete := scanHTMLFilesConcurrently(htmlFiles, concurrency)
	if fingerprintOK && complete {
		saveCSSPurgeHTMLCache(cachePath, cssPurgeHTMLCache{Schema: cssPurgeHTMLCacheSchema, Fingerprint: fingerprint, Selectors: *used})
	}
	cssPurgeLog.Printf("HTML selector scan: files=%d bytes=%d cache_hit=false", len(htmlFiles), scannedBytes)

	if verbose {
		cssPurgeLog.Printf("Found %d classes, %d IDs, %d elements, %d attributes",
			len(used.Classes), len(used.IDs), len(used.Elements), len(used.Attributes))
	}

	return used, false, nil
}

const cssPurgeHTMLCacheSchema = 1

type cssPurgeHTMLCache struct {
	Schema      int                    `json:"schema"`
	Fingerprint string                 `json:"fingerprint"`
	Selectors   csspurge.UsedSelectors `json:"selectors"`
}

func loadCSSPurgeHTMLCache(path string) (cssPurgeHTMLCache, bool) {
	var entry cssPurgeHTMLCache
	b, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(b, &entry) != nil {
		return entry, false
	}
	return entry, true
}

func saveCSSPurgeHTMLCache(path string, entry cssPurgeHTMLCache) {
	b, err := json.Marshal(entry)
	if err != nil || os.MkdirAll(filepath.Dir(path), 0o755) != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".html-selectors-*.tmp")
	if err != nil {
		return
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err = tmp.Write(b); err != nil {
		_ = tmp.Close()
		return
	}
	if err = tmp.Close(); err != nil {
		return
	}
	if err := replaceDiagnosticsArtifact(name, path); err != nil {
		cssPurgeLog.Printf("could not update HTML selector cache: %v", err)
	}
}

type htmlFingerprintResult struct {
	path string
	sum  [sha256.Size]byte
	size int64
	err  error
}

func fingerprintHTMLFiles(files []string, outputDir string, concurrency int) (fingerprint string, scannedBytes int64, ok bool) {
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > len(files) {
		concurrency = len(files)
	}
	jobs := make(chan string, len(files))
	results := make(chan htmlFingerprintResult, len(files))
	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				content, err := os.ReadFile(path)
				result := htmlFingerprintResult{path: path, err: err}
				if err == nil {
					result.sum = sha256.Sum256(content)
					result.size = int64(len(content))
				}
				results <- result
			}
		}()
	}
	for _, path := range files {
		jobs <- path
	}
	close(jobs)
	wg.Wait()
	close(results)

	items := make([]htmlFingerprintResult, 0, len(files))
	var totalBytes int64
	for result := range results {
		if result.err != nil {
			cssPurgeLog.Warnf("fingerprinting %s: %v", result.path, result.err)
			return "", totalBytes, false
		}
		totalBytes += result.size
		items = append(items, result)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].path < items[j].path })
	input := sha256.New()
	for _, item := range items {
		rel, err := filepath.Rel(outputDir, item.path)
		if err != nil {
			rel = item.path
		}
		_, _ = input.Write([]byte(filepath.ToSlash(rel)))
		_, _ = input.Write([]byte{0})
		_, _ = input.Write(item.sum[:])
	}
	return hex.EncodeToString(input.Sum(nil)), totalBytes, true
}

// buildPurgeOptions creates PurgeOptions from configuration.
func buildPurgeOptions(purgeConfig models.CSSPurgeConfig, verbose bool) csspurge.PurgeOptions {
	preserve := purgeConfig.Preserve
	if len(preserve) == 0 {
		preserve = csspurge.DefaultPreservePatterns()
	}

	preserveAttrs := purgeConfig.PreserveAttributes
	if len(preserveAttrs) == 0 {
		preserveAttrs = csspurge.DefaultPreserveAttributes()
	}

	return csspurge.PurgeOptions{
		Preserve:           preserve,
		PreserveAttributes: preserveAttrs,
		Verbose:            verbose,
	}
}

// processCSSFiles processes each CSS file concurrently and returns statistics.
func processCSSFiles(cssFiles []string, outputDir, cacheDir string, used *csspurge.UsedSelectors, opts csspurge.PurgeOptions, purgeConfig models.CSSPurgeConfig, verbose bool) purgeProcessingStats {
	filteredFiles := make([]string, 0, len(cssFiles))
	skippedCount := 0

	for _, cssFile := range cssFiles {
		relPath, err := filepath.Rel(outputDir, cssFile)
		if err != nil {
			relPath = cssFile
		}

		if shouldSkipCSSFile(relPath, purgeConfig.SkipFiles) {
			if verbose {
				cssPurgeLog.Printf("Skipping %s (matches skip pattern)", relPath)
			}
			skippedCount++
			continue
		}

		filteredFiles = append(filteredFiles, cssFile)
	}

	if len(filteredFiles) == 0 {
		return purgeProcessingStats{filesSkipped: skippedCount}
	}

	return processCSSFilesConcurrently(filteredFiles, outputDir, cacheDir, used, opts, verbose, skippedCount)
}

// cssFileResult holds the result of processing a single CSS file.
type cssFileResult struct {
	cssFile    string
	relPath    string
	origSize   int
	purgedSize int
	rules      int
	removed    int
	cacheHit   bool
}

type cssPurgeCacheEntry struct {
	SelectorHash string `json:"selector_hash"`
	OptionsHash  string `json:"options_hash"`
	InputHash    string `json:"input_hash"`
	OutputHash   string `json:"output_hash"`
	OutputCSS    string `json:"output_css"`
	OriginalSize int    `json:"original_size"`
	PurgedSize   int    `json:"purged_size"`
	Rules        int    `json:"rules"`
	Removed      int    `json:"removed"`
}

func cssPurgeHash(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		// An unhashable configuration must never produce a reusable cache key.
		return ""
	}
	return cssPurgeBytesHash(b)
}

func cssPurgeBytesHash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func cssPurgeCachePath(cacheDir, outputDir, cssFile string) string {
	rel, err := filepath.Rel(outputDir, cssFile)
	if err != nil {
		rel = filepath.Base(cssFile)
	}
	sum := sha256.Sum256([]byte(filepath.Clean(rel)))
	return filepath.Join(cacheDir, hex.EncodeToString(sum[:])+".json")
}

func loadCSSPurgeCache(path string) (cssPurgeCacheEntry, bool) {
	var entry cssPurgeCacheEntry
	b, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(b, &entry) != nil {
		return entry, false
	}
	// Reject partially written or tampered entries. Legacy entries without
	// source bytes also miss safely and are regenerated on the next purge.
	if entry.InputHash == "" || entry.OutputHash != cssPurgeBytesHash([]byte(entry.OutputCSS)) ||
		entry.OriginalSize < 0 || entry.PurgedSize < 0 || entry.Rules < 0 || entry.Removed < 0 {
		return entry, false
	}
	return entry, true
}

func saveCSSPurgeCache(path string, entry cssPurgeCacheEntry) {
	b, err := json.Marshal(entry)
	if err != nil {
		return
	}
	if os.MkdirAll(filepath.Dir(path), 0o755) != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".css-purge-*.tmp")
	if err != nil {
		return
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err = tmp.Write(b); err != nil {
		_ = tmp.Close()
		return
	}
	if err = tmp.Close(); err != nil {
		return
	}
	if err := replaceDiagnosticsArtifact(name, path); err != nil {
		cssPurgeLog.Printf("could not update purge cache %s: %v", filepath.Base(path), err)
	}
}

// processSingleCSSPurgeFile either restores validated output or purges the
// current input. Rebuilt originals and already-purged files both hit the cache.
func processSingleCSSPurgeFile(cssFile, outputDir, cacheDir string, used *csspurge.UsedSelectors, opts csspurge.PurgeOptions, selectorsHash, optionsHash string) (cssFileResult, error) {
	relPath, err := filepath.Rel(outputDir, cssFile)
	if err != nil {
		relPath = cssFile
	}
	result := cssFileResult{cssFile: cssFile, relPath: relPath}
	content, err := os.ReadFile(cssFile)
	if err != nil {
		return result, fmt.Errorf("reading %s: %w", relPath, err)
	}

	cachePath := cssPurgeCachePath(cacheDir, outputDir, cssFile)
	inputHash := cssPurgeBytesHash(content)
	cacheable := selectorsHash != "" && optionsHash != ""
	if cacheable {
		cached, ok := loadCSSPurgeCache(cachePath)
		if ok && cached.SelectorHash == selectorsHash && cached.OptionsHash == optionsHash &&
			(cached.InputHash == inputHash || cached.OutputHash == inputHash) {
			if cached.OutputHash != inputHash {
				//nolint:gosec // G306: generated CSS must be web-readable
				if err := os.WriteFile(cssFile, []byte(cached.OutputCSS), 0o644); err != nil {
					return result, fmt.Errorf("restoring cached CSS %s: %w", relPath, err)
				}
			}
			result.origSize, result.purgedSize = cached.OriginalSize, cached.PurgedSize
			result.rules, result.removed, result.cacheHit = cached.Rules, cached.Removed, true
			return result, nil
		}
	}

	purged, purgeStats := csspurge.PurgeCSS(string(content), used, opts)
	result.origSize = purgeStats.OriginalSize
	result.purgedSize = purgeStats.PurgedSize
	result.rules = purgeStats.TotalRules
	result.removed = purgeStats.RemovedRules

	// Avoid overwriting the original bytes if nothing was actually removed.
	outputCSS := string(content)
	if purgeStats.RemovedRules > 0 {
		outputCSS = purged
		//nolint:gosec // G306: generated CSS must be web-readable
		if err := os.WriteFile(cssFile, []byte(purged), 0o644); err != nil {
			return result, fmt.Errorf("writing %s: %w", relPath, err)
		}
	}
	if cacheable {
		saveCSSPurgeCache(cachePath, cssPurgeCacheEntry{
			SelectorHash: selectorsHash,
			OptionsHash:  optionsHash,
			InputHash:    inputHash,
			OutputHash:   cssPurgeBytesHash([]byte(outputCSS)),
			OutputCSS:    outputCSS,
			OriginalSize: purgeStats.OriginalSize,
			PurgedSize:   purgeStats.PurgedSize,
			Rules:        purgeStats.TotalRules,
			Removed:      purgeStats.RemovedRules,
		})
	}
	return result, nil
}

// processCSSFilesConcurrently processes CSS files using a worker pool.
func processCSSFilesConcurrently(cssFiles []string, outputDir, cacheDir string, used *csspurge.UsedSelectors, opts csspurge.PurgeOptions, verbose bool, skippedCount int) purgeProcessingStats {
	concurrency := len(cssFiles)
	if concurrency > 16 {
		concurrency = 16
	}

	if verbose {
		cssPurgeLog.Infof("Processing %d CSS files with %d workers", len(cssFiles), concurrency)
	}

	jobs := make(chan string, len(cssFiles))
	results := make(chan cssFileResult, len(cssFiles))
	errors := make(chan error, len(cssFiles))

	selectorsHash, optionsHash := cssPurgeHash(used), cssPurgeHash(opts)

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for cssFile := range jobs {
				result, err := processSingleCSSPurgeFile(cssFile, outputDir, cacheDir, used, opts, selectorsHash, optionsHash)
				if err != nil {
					errors <- err
					continue
				}
				results <- result
			}
		}()
	}

	for _, cssFile := range cssFiles {
		jobs <- cssFile
	}
	close(jobs)

	wg.Wait()
	close(results)
	close(errors)

	for err := range errors {
		cssPurgeLog.Warnf("%v", err)
	}

	var stats purgeProcessingStats
	stats.filesSkipped = skippedCount

	for result := range results {
		stats.totalOriginal += result.origSize
		stats.totalPurged += result.purgedSize
		stats.filesProcessed++
		if result.cacheHit {
			stats.cacheHits++
		}

		if verbose && result.removed > 0 {
			savings := float64(result.origSize-result.purgedSize) / float64(result.origSize) * 100
			cssPurgeLog.Infof("%s: removed %d/%d rules (%.1f%% reduction, %d -> %d bytes)",
				result.relPath, result.removed, result.rules, savings, result.origSize, result.purgedSize)
		} else if verbose {
			cssPurgeLog.Infof("%s: all %d rules are used", result.relPath, result.rules)
		}
	}

	return stats
}

// reportPurgeSummary reports the purging summary.
func reportPurgeSummary(stats purgeProcessingStats, purgeConfig models.CSSPurgeConfig, verbose bool) {
	if stats.filesProcessed > 0 {
		savings := float64(stats.totalOriginal-stats.totalPurged) / float64(stats.totalOriginal) * 100
		cssPurgeLog.Printf("Processed %d CSS files (%d cache hits): %d -> %d bytes (%.1f%% reduction)",
			stats.filesProcessed, stats.cacheHits, stats.totalOriginal, stats.totalPurged, savings)

		if purgeConfig.WarningThreshold > 0 && int(savings) > purgeConfig.WarningThreshold {
			cssPurgeLog.Warnf("Removed %.1f%% of CSS (threshold: %d%%). This might indicate overly aggressive purging. Consider adding patterns to 'preserve' config.",
				savings, purgeConfig.WarningThreshold)
		}
	}

	if stats.filesSkipped > 0 && verbose {
		cssPurgeLog.Printf("Skipped %d CSS files", stats.filesSkipped)
	}
}

// Priority returns the plugin priority for the cleanup stage.
// CSS purge runs before Pagefind (PriorityDefault) since Pagefind indexes
// content but doesn't need CSS, and we want smaller CSS files deployed.
func (p *CSSPurgePlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCleanup {
		return lifecycle.PriorityDefault - 10 // Before Pagefind
	}
	return lifecycle.PriorityDefault
}

// findHTMLFiles recursively finds all HTML files in a directory.
func findHTMLFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".html") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// findCSSFiles recursively finds all CSS files in a directory.
func findCSSFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".css") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// scanHTMLFilesConcurrently scans HTML files using a worker pool.
func scanHTMLFilesConcurrently(files []string, concurrency int) (*csspurge.UsedSelectors, bool) {
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > len(files) {
		concurrency = len(files)
	}

	// Use channels for work distribution
	jobs := make(chan string, len(files))
	results := make(chan *csspurge.UsedSelectors, len(files))
	errors := make(chan error, len(files))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				used := csspurge.NewUsedSelectors()
				if err := csspurge.ScanHTML(path, used); err != nil {
					errors <- fmt.Errorf("scanning %s: %w", path, err)
					continue
				}
				results <- used
			}
		}()
	}

	// Send jobs
	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
	close(results)
	close(errors)

	// Check for errors
	complete := true
	for err := range errors {
		// Log but don't fail - partial analysis is still useful
		cssPurgeLog.Warnf("%v", err)
		complete = false
	}

	// Merge results
	combined := csspurge.NewUsedSelectors()
	for used := range results {
		combined.Merge(used)
	}

	return combined, complete
}

// shouldSkipCSSFile checks if a CSS file matches any skip pattern.
func shouldSkipCSSFile(relPath string, patterns []string) bool {
	if isVersionedCSSFile(relPath) {
		return true
	}
	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}
		// Also check just the filename
		matched, err = filepath.Match(pattern, filepath.Base(relPath))
		if err == nil && matched {
			return true
		}
	}
	return false
}

// getCSSPurgeConfig extracts CSSPurgeConfig from config.Extra.
func getCSSPurgeConfig(config *lifecycle.Config) models.CSSPurgeConfig {
	if config.Extra == nil {
		return models.NewCSSPurgeConfig()
	}

	// Try direct type assertion
	if pc, ok := config.Extra["css_purge"].(models.CSSPurgeConfig); ok {
		return pc
	}

	// Try to parse from map if stored as map[string]interface{}
	rawConfig, ok := config.Extra["css_purge"].(map[string]interface{})
	if !ok {
		return models.NewCSSPurgeConfig()
	}

	return parseCSSPurgeConfigFromMap(rawConfig)
}

// parseCSSPurgeConfigFromMap parses CSSPurgeConfig from a raw map.
func parseCSSPurgeConfigFromMap(rawConfig map[string]interface{}) models.CSSPurgeConfig {
	result := models.NewCSSPurgeConfig()

	if enabled, ok := rawConfig["enabled"].(bool); ok {
		result.Enabled = enabled
	}
	if verbose, ok := rawConfig["verbose"].(bool); ok {
		result.Verbose = verbose
	}
	if preserve, ok := rawConfig["preserve"]; ok {
		result.Preserve = parseStringSlice(preserve)
	}
	if skip, ok := rawConfig["skip_files"]; ok {
		result.SkipFiles = parseStringSlice(skip)
	}
	if preserveAttrs, ok := rawConfig["preserve_attributes"]; ok {
		result.PreserveAttributes = parseStringSlice(preserveAttrs)
	}
	if threshold, ok := parseIntFromInterface(rawConfig["warning_threshold"]); ok {
		result.WarningThreshold = threshold
	}

	return result
}

// parseStringSlice extracts a string slice from an interface value.
func parseStringSlice(value interface{}) []string {
	switch values := value.(type) {
	case []interface{}:
		result := make([]string, 0, len(values))
		for _, v := range values {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return values
	default:
		return nil
	}
}

// Ensure CSSPurgePlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*CSSPurgePlugin)(nil)
	_ lifecycle.CleanupPlugin  = (*CSSPurgePlugin)(nil)
	_ lifecycle.PriorityPlugin = (*CSSPurgePlugin)(nil)
)
