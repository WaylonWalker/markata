package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
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
	used, stats, err := scanHTMLFilesForSelectorsWithStats(outputDir, cacheDir, concurrency, verbose)
	return used, stats.CacheHit, err
}

type htmlSelectorScanStats struct {
	Files       int
	FilesParsed int
	FilesReused int
	BytesHashed int64
	CacheHit    bool
	CacheMode   string
}

const (
	htmlSelectorCacheModeIdentity        = "identity"
	htmlSelectorCacheModeIdentityPartial = "identity-partial"
	htmlSelectorCacheModeContent         = "content"
	htmlSelectorCacheModePartial         = "partial"
	htmlSelectorCacheModeScan            = "scan"
	htmlSelectorFileCacheName            = "html-selectors-per-file.json"
)

func scanHTMLFilesForSelectorsWithStats(outputDir, cacheDir string, concurrency int, verbose bool) (*csspurge.UsedSelectors, htmlSelectorScanStats, error) {
	htmlFiles, err := findHTMLFiles(outputDir)
	if err != nil {
		return nil, htmlSelectorScanStats{}, fmt.Errorf("failed to find HTML files: %w", err)
	}

	if len(htmlFiles) == 0 {
		if verbose {
			cssPurgeLog.Printf("No HTML files found, skipping")
		}
		return nil, htmlSelectorScanStats{}, nil
	}

	if verbose {
		cssPurgeLog.Printf("Found %d HTML files to analyze", len(htmlFiles))
	}

	cachePath := filepath.Join(cacheDir, "html-selectors.json")
	cached, cacheExists := loadCSSPurgeHTMLCache(cachePath)
	cacheValid := cacheExists && validCSSPurgeHTMLCache(cached)
	var identities map[string]fileidentity.Identity
	identitiesComplete := false
	if cacheValid && fileidentity.Available {
		identities, identitiesComplete = currentHTMLFileIdentities(htmlFiles, outputDir)
		if identitiesComplete && cssPurgeHTMLIdentityCacheMatches(cached, identities, time.Now()) {
			stats := htmlSelectorScanStats{
				Files:     len(htmlFiles),
				CacheHit:  true,
				CacheMode: htmlSelectorCacheModeIdentity,
			}
			cssPurgeLog.Printf("HTML selector cache hit: files=%d bytes_hashed=0 mode=identity", stats.Files)
			return &cached.Selectors, stats, nil
		}
	}

	fileCachePath := filepath.Join(cacheDir, htmlSelectorFileCacheName)
	fileCache, fileCacheExists := loadCSSPurgeHTMLFileCache(fileCachePath)
	fileCacheValid := fileCacheExists && validCSSPurgeHTMLFileCache(fileCache)
	if cacheValid && fileCacheValid && fileCache.Fingerprint == cached.Fingerprint && identitiesComplete {
		used, stats, handled, err := scanHTMLFilesByIdentity(htmlFiles, outputDir, cachePath, fileCachePath, concurrency, cached, fileCache, identities)
		if handled || err != nil {
			return used, stats, err
		}
	}

	fingerprint := fingerprintHTMLFiles(htmlFiles, outputDir, concurrency)
	stats := htmlSelectorScanStats{Files: len(htmlFiles), BytesHashed: fingerprint.scannedBytes}
	if !fingerprint.ok {
		used, complete := scanHTMLFilesConcurrently(htmlFiles, concurrency)
		stats.FilesParsed = len(htmlFiles)
		stats.CacheMode = htmlSelectorCacheModeScan
		cssPurgeLog.Printf("HTML selector scan: files=%d bytes_hashed=%d parsed=%d reused=0 cache_hit=false mode=%s", stats.Files, stats.BytesHashed, stats.FilesParsed, stats.CacheMode)
		if !complete {
			return nil, stats, fmt.Errorf("HTML selector scan incomplete; aborting CSS purge")
		}
		return used, stats, nil
	}
	if cacheValid && cached.Fingerprint == fingerprint.fingerprint {
		cached.Observed = time.Now()
		cached.HTMLBytes = fingerprint.scannedBytes
		cached.Files = fingerprint.identities
		saveCSSPurgeHTMLCache(cachePath, cached)
		stats.CacheHit = true
		stats.CacheMode = htmlSelectorCacheModeContent
		cssPurgeLog.Printf("HTML selector cache hit: files=%d bytes_hashed=%d mode=content", len(htmlFiles), fingerprint.scannedBytes)
		return &cached.Selectors, stats, nil
	}

	used := csspurge.NewUsedSelectors()
	fileSelectors := make(map[string]cssPurgeHTMLFileSelectors, len(fingerprint.fileHashes))
	changedFiles := make([]string, 0, len(fingerprint.fileHashes))
	for relative, contentHash := range fingerprint.fileHashes {
		if fileCacheValid {
			if previous, ok := fileCache.Files[relative]; ok && previous.ContentHash == contentHash {
				fileSelectors[relative] = previous
				used.Merge(previous.Selectors.toUsedSelectors())
				stats.FilesReused++
				continue
			}
		}
		changedFiles = append(changedFiles, filepath.Join(outputDir, filepath.FromSlash(relative)))
	}
	sort.Strings(changedFiles)
	parsedSelectors, complete, hashesMatch := scanHTMLFilesWithExpectedHashesConcurrently(changedFiles, outputDir, fingerprint.fileHashes, concurrency)
	stats.FilesParsed = len(parsedSelectors)
	if err := validateCSSPurgeHTMLScan(complete, hashesMatch); err != nil {
		return nil, stats, err
	}
	for path, selectors := range parsedSelectors {
		relative, err := filepath.Rel(outputDir, path)
		if err != nil || !filepath.IsLocal(relative) {
			complete = false
			continue
		}
		relative = filepath.ToSlash(relative)
		fileSelectors[relative] = cssPurgeHTMLFileSelectors{
			ContentHash: fingerprint.fileHashes[relative],
			Selectors:   cssPurgeCachedSelectorsFromUsed(selectors),
		}
		fileSelectors[relative] = sealCSSPurgeHTMLFileSelectors(fileSelectors[relative])
		used.Merge(selectors)
	}
	if complete && hashesMatch {
		fileCache := cssPurgeHTMLFileCache{
			Schema:      cssPurgeHTMLFileCacheSchema,
			Fingerprint: fingerprint.fingerprint,
			Files:       fileSelectors,
		}
		saveCSSPurgeHTMLFileCache(fileCachePath, fileCache)
		saveCSSPurgeHTMLCache(cachePath, cssPurgeHTMLCache{
			Schema:       cssPurgeHTMLCacheSchema,
			Fingerprint:  fingerprint.fingerprint,
			SelectorHash: cssPurgeAggregateSelectorHash(used),
			Observed:     time.Now(),
			HTMLBytes:    fingerprint.scannedBytes,
			Files:        fingerprint.identities,
			Selectors:    *used,
		})
	}
	if stats.FilesReused > 0 {
		stats.CacheMode = htmlSelectorCacheModePartial
	} else {
		stats.CacheMode = htmlSelectorCacheModeScan
	}
	cssPurgeLog.Printf("HTML selector scan: files=%d bytes_hashed=%d parsed=%d reused=%d cache_hit=false mode=%s", stats.Files, stats.BytesHashed, stats.FilesParsed, stats.FilesReused, stats.CacheMode)

	if verbose {
		cssPurgeLog.Printf("Found %d classes, %d IDs, %d elements, %d attributes",
			len(used.Classes), len(used.IDs), len(used.Elements), len(used.Attributes))
	}

	return used, stats, nil
}

// scanHTMLFilesByIdentity reuses per-file selector data for files that still
// have the identity captured during the prior exact-content scan. Only files
// whose identity is new or uncertain need an exact content hash.
func scanHTMLFilesByIdentity(
	files []string,
	outputDir, cachePath, fileCachePath string,
	concurrency int,
	cached cssPurgeHTMLCache,
	fileCache cssPurgeHTMLFileCache,
	currentIdentities map[string]fileidentity.Identity,
) (*csspurge.UsedSelectors, htmlSelectorScanStats, bool, error) {
	stats := htmlSelectorScanStats{Files: len(files)}
	if now := time.Now(); cached.Observed.IsZero() || now.Before(cached.Observed) {
		// A missing or future observation cannot authorize identity reuse. The
		// caller will use the exact-content fallback for the full file set.
		return nil, stats, false, nil
	}

	used, fileSelectors, changedFiles, reused, ok := reuseCSSPurgeSelectorsByIdentity(files, outputDir, cached, fileCache, currentIdentities)
	if !ok {
		return nil, stats, false, nil
	}
	stats.FilesReused = reused
	fingerprint, toParse, reused, ok := fingerprintChangedCSSPurgeFiles(changedFiles, outputDir, concurrency, fileCache, used, fileSelectors)
	if !ok {
		return nil, stats, false, nil
	}
	stats.BytesHashed = fingerprint.scannedBytes
	stats.FilesReused += reused

	parsedSelectors, complete, hashesMatch := scanHTMLFilesWithExpectedHashesConcurrently(toParse, outputDir, fingerprint.fileHashes, concurrency)
	stats.FilesParsed = len(parsedSelectors)
	if err := validateCSSPurgeHTMLScan(complete, hashesMatch); err != nil {
		return nil, stats, true, err
	}
	if err := mergeParsedCSSPurgeSelectors(parsedSelectors, outputDir, fingerprint.fileHashes, fileSelectors, used); err != nil {
		return nil, stats, true, err
	}
	if len(fileSelectors) != len(files) {
		return nil, stats, false, nil
	}

	nextIdentities, ok := cssPurgeIdentitiesAfterFingerprint(outputDir, currentIdentities, changedFiles, fingerprint.identities)
	if !ok {
		return nil, stats, false, nil
	}
	contentFingerprint := cssPurgeHTMLFingerprintFromCache(fileSelectors)
	if contentFingerprint == "" {
		return nil, stats, false, nil
	}
	fileCache = cssPurgeHTMLFileCache{
		Schema:      cssPurgeHTMLFileCacheSchema,
		Fingerprint: contentFingerprint,
		Files:       fileSelectors,
	}
	saveCSSPurgeHTMLFileCache(fileCachePath, fileCache)

	var htmlBytes int64
	for _, identity := range currentIdentities {
		htmlBytes += identity.Size
	}
	saveCSSPurgeHTMLCache(cachePath, cssPurgeHTMLCache{
		Schema:       cssPurgeHTMLCacheSchema,
		Fingerprint:  contentFingerprint,
		SelectorHash: cssPurgeAggregateSelectorHash(used),
		Observed:     time.Now(),
		HTMLBytes:    htmlBytes,
		Files:        nextIdentities,
		Selectors:    *used,
	})

	stats.CacheHit = contentFingerprint == cached.Fingerprint
	stats.CacheMode = htmlSelectorCacheModeIdentityPartial
	cssPurgeLog.Printf("HTML selector scan: files=%d bytes_hashed=%d parsed=%d reused=%d cache_hit=%t mode=%s", stats.Files, stats.BytesHashed, stats.FilesParsed, stats.FilesReused, stats.CacheHit, stats.CacheMode)
	return used, stats, true, nil
}

func reuseCSSPurgeSelectorsByIdentity(
	files []string,
	outputDir string,
	cached cssPurgeHTMLCache,
	fileCache cssPurgeHTMLFileCache,
	currentIdentities map[string]fileidentity.Identity,
) (used *csspurge.UsedSelectors, fileSelectors map[string]cssPurgeHTMLFileSelectors, changedFiles []string, reused int, ok bool) {
	used = csspurge.NewUsedSelectors()
	fileSelectors = make(map[string]cssPurgeHTMLFileSelectors, len(files))
	changedFiles = make([]string, 0, len(files))
	for _, path := range files {
		relative, err := filepath.Rel(outputDir, path)
		if err != nil || !filepath.IsLocal(relative) {
			return nil, nil, nil, reused, false
		}
		relative = filepath.ToSlash(relative)
		identity := currentIdentities[relative]
		previousIdentity, hasPreviousIdentity := cached.Files[relative]
		previousSelectors, hasPreviousSelectors := fileCache.Files[relative]
		if hasPreviousIdentity && hasPreviousSelectors && identity == previousIdentity && identity.StableBefore(cached.Observed) {
			fileSelectors[relative] = previousSelectors
			used.Merge(previousSelectors.Selectors.toUsedSelectors())
			reused++
			continue
		}
		changedFiles = append(changedFiles, path)
	}
	sort.Strings(changedFiles)
	return used, fileSelectors, changedFiles, reused, true
}

func fingerprintChangedCSSPurgeFiles(
	changedFiles []string,
	outputDir string,
	concurrency int,
	fileCache cssPurgeHTMLFileCache,
	used *csspurge.UsedSelectors,
	fileSelectors map[string]cssPurgeHTMLFileSelectors,
) (fingerprint htmlFingerprintSummary, toParse []string, reused int, ok bool) {
	fingerprint = htmlFingerprintSummary{fileHashes: make(map[string]string), identities: make(map[string]fileidentity.Identity), ok: true}
	if len(changedFiles) > 0 {
		fingerprint = fingerprintHTMLFiles(changedFiles, outputDir, concurrency)
		if !fingerprint.ok {
			return fingerprint, nil, 0, false
		}
	}

	toParse = make([]string, 0, len(changedFiles))
	for _, path := range changedFiles {
		relative, err := filepath.Rel(outputDir, path)
		if err != nil || !filepath.IsLocal(relative) {
			return fingerprint, nil, reused, false
		}
		relative = filepath.ToSlash(relative)
		contentHash := fingerprint.fileHashes[relative]
		if previous, ok := fileCache.Files[relative]; ok && previous.ContentHash == contentHash {
			fileSelectors[relative] = previous
			used.Merge(previous.Selectors.toUsedSelectors())
			reused++
			continue
		}
		toParse = append(toParse, path)
	}
	return fingerprint, toParse, reused, true
}

func mergeParsedCSSPurgeSelectors(
	parsed map[string]*csspurge.UsedSelectors,
	outputDir string,
	expectedHashes map[string]string,
	fileSelectors map[string]cssPurgeHTMLFileSelectors,
	used *csspurge.UsedSelectors,
) error {
	for path, selectors := range parsed {
		relative, err := filepath.Rel(outputDir, path)
		if err != nil || !filepath.IsLocal(relative) {
			return fmt.Errorf("invalid HTML path %s", path)
		}
		relative = filepath.ToSlash(relative)
		fileSelectors[relative] = sealCSSPurgeHTMLFileSelectors(cssPurgeHTMLFileSelectors{
			ContentHash: expectedHashes[relative],
			Selectors:   cssPurgeCachedSelectorsFromUsed(selectors),
		})
		used.Merge(selectors)
	}
	return nil
}

func cssPurgeIdentitiesAfterFingerprint(
	outputDir string,
	currentIdentities map[string]fileidentity.Identity,
	changedFiles []string,
	fingerprintedIdentities map[string]fileidentity.Identity,
) (map[string]fileidentity.Identity, bool) {
	identities := make(map[string]fileidentity.Identity, len(currentIdentities))
	for relative, identity := range currentIdentities {
		identities[relative] = identity
	}
	for _, path := range changedFiles {
		relative, err := filepath.Rel(outputDir, path)
		if err != nil || !filepath.IsLocal(relative) {
			return nil, false
		}
		delete(identities, filepath.ToSlash(relative))
	}
	for relative, identity := range fingerprintedIdentities {
		identities[relative] = identity
	}
	if len(identities) == 0 {
		return nil, true
	}
	return identities, true
}

const cssPurgeHTMLCacheSchema = 4
const cssPurgeHTMLFileCacheSchema = 1

type cssPurgeHTMLFileSelectors struct {
	ContentHash  string                  `json:"content_hash"`
	SelectorHash string                  `json:"selector_hash"`
	Selectors    cssPurgeCachedSelectors `json:"selectors"`
}

type cssPurgeHTMLFileCache struct {
	Schema      int                                  `json:"schema"`
	Fingerprint string                               `json:"fingerprint"`
	Files       map[string]cssPurgeHTMLFileSelectors `json:"files"`
}

type cssPurgeCachedSelectors struct {
	Classes    []string `json:"classes"`
	IDs        []string `json:"ids"`
	Elements   []string `json:"elements"`
	Attributes []string `json:"attributes"`
}

func cssPurgeCachedSelectorsFromUsed(used *csspurge.UsedSelectors) cssPurgeCachedSelectors {
	return cssPurgeCachedSelectors{
		Classes:    sortedCSSPurgeSelectorNames(used.Classes),
		IDs:        sortedCSSPurgeSelectorNames(used.IDs),
		Elements:   sortedCSSPurgeSelectorNames(used.Elements),
		Attributes: sortedCSSPurgeSelectorNames(used.Attributes),
	}
}

func sealCSSPurgeHTMLFileSelectors(entry cssPurgeHTMLFileSelectors) cssPurgeHTMLFileSelectors {
	data, err := json.Marshal(entry.Selectors)
	if err != nil {
		return entry
	}
	entry.SelectorHash = cssPurgeBytesHash(data)
	return entry
}

func sortedCSSPurgeSelectorNames(values map[string]bool) []string {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (cached cssPurgeCachedSelectors) toUsedSelectors() *csspurge.UsedSelectors {
	used := csspurge.NewUsedSelectors()
	for _, name := range cached.Classes {
		used.Classes[name] = true
	}
	for _, name := range cached.IDs {
		used.IDs[name] = true
	}
	for _, name := range cached.Elements {
		used.Elements[name] = true
	}
	for _, name := range cached.Attributes {
		used.Attributes[name] = true
	}
	return used
}

type cssPurgeHTMLCache struct {
	Schema       int                              `json:"schema"`
	Fingerprint  string                           `json:"fingerprint"`
	SelectorHash string                           `json:"selector_hash"`
	Observed     time.Time                        `json:"observed,omitempty"`
	HTMLBytes    int64                            `json:"html_bytes,omitempty"`
	Files        map[string]fileidentity.Identity `json:"files,omitempty"`
	Selectors    csspurge.UsedSelectors           `json:"selectors"`
}

func validCSSPurgeHTMLCache(cached cssPurgeHTMLCache) bool {
	return cached.Schema == cssPurgeHTMLCacheSchema && cached.Fingerprint != "" && cached.HTMLBytes >= 0 &&
		cached.Selectors.Classes != nil && cached.Selectors.IDs != nil &&
		cached.Selectors.Elements != nil && cached.Selectors.Attributes != nil &&
		cached.SelectorHash != "" && cached.SelectorHash == cssPurgeAggregateSelectorHash(&cached.Selectors)
}

func cssPurgeAggregateSelectorHash(used *csspurge.UsedSelectors) string {
	data, err := json.Marshal(used)
	if err != nil {
		return ""
	}
	return cssPurgeBytesHash(data)
}

func validateCSSPurgeHTMLScan(complete, hashesMatch bool) error {
	if !complete {
		return fmt.Errorf("HTML selector scan incomplete; aborting CSS purge")
	}
	if !hashesMatch {
		return fmt.Errorf("HTML changed during selector scan; aborting CSS purge")
	}
	return nil
}

func loadCSSPurgeHTMLFileCache(path string) (cssPurgeHTMLFileCache, bool) {
	var entry cssPurgeHTMLFileCache
	b, err := os.ReadFile(path)
	if err != nil || json.Unmarshal(b, &entry) != nil {
		return entry, false
	}
	return entry, true
}

func validCSSPurgeHTMLFileCache(cached cssPurgeHTMLFileCache) bool {
	if cached.Schema != cssPurgeHTMLFileCacheSchema || cached.Fingerprint == "" || len(cached.Files) == 0 {
		return false
	}
	for path := range cached.Files {
		entry := cached.Files[path]
		if !filepath.IsLocal(filepath.FromSlash(path)) || len(entry.ContentHash) != sha256.Size*2 ||
			entry.Selectors.Classes == nil || entry.Selectors.IDs == nil ||
			entry.Selectors.Elements == nil || entry.Selectors.Attributes == nil {
			return false
		}
		if _, err := hex.DecodeString(entry.ContentHash); err != nil {
			return false
		}
		sealed := sealCSSPurgeHTMLFileSelectors(entry)
		if entry.SelectorHash == "" || sealed.SelectorHash != entry.SelectorHash {
			return false
		}
	}
	return cssPurgeHTMLFingerprintFromCache(cached.Files) == cached.Fingerprint
}

func cssPurgeHTMLFingerprintFromCache(files map[string]cssPurgeHTMLFileSelectors) string {
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	input := sha256.New()
	for _, path := range paths {
		hash, err := hex.DecodeString(files[path].ContentHash)
		if err != nil || len(hash) != sha256.Size {
			return ""
		}
		_, _ = input.Write([]byte(path))
		_, _ = input.Write([]byte{0})
		_, _ = input.Write(hash)
	}
	return hex.EncodeToString(input.Sum(nil))
}

func currentHTMLFileIdentities(files []string, outputDir string) (map[string]fileidentity.Identity, bool) {
	if !fileidentity.Available {
		return nil, false
	}
	identities := make(map[string]fileidentity.Identity, len(files))
	for _, path := range files {
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return nil, false
		}
		identity := fileidentity.Read(info)
		if identity.Inode == 0 {
			return nil, false
		}
		relative, err := filepath.Rel(outputDir, path)
		if err != nil || !filepath.IsLocal(relative) {
			return nil, false
		}
		identities[filepath.ToSlash(relative)] = identity
	}
	return identities, len(identities) == len(files)
}

func cssPurgeHTMLIdentityCacheMatches(cached cssPurgeHTMLCache, current map[string]fileidentity.Identity, now time.Time) bool {
	if !fileidentity.Available || !validCSSPurgeHTMLCache(cached) || cached.Observed.IsZero() ||
		now.Before(cached.Observed) || len(cached.Files) == 0 || len(cached.Files) != len(current) {
		return false
	}
	for path, identity := range current {
		if previous, ok := cached.Files[path]; !ok || previous != identity || !identity.StableBefore(cached.Observed) {
			return false
		}
	}
	return true
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
	entry.SelectorHash = cssPurgeAggregateSelectorHash(&entry.Selectors)
	saveCSSPurgeCacheArtifact(path, ".html-selectors-*.tmp", "HTML selector", entry)
}

func saveCSSPurgeHTMLFileCache(path string, entry cssPurgeHTMLFileCache) {
	saveCSSPurgeCacheArtifact(path, ".html-selector-files-*.tmp", "per-file HTML selector", entry)
}

func saveCSSPurgeCacheArtifact(path, tempPattern, label string, entry any) {
	b, err := json.Marshal(entry)
	if err != nil || os.MkdirAll(filepath.Dir(path), 0o755) != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), tempPattern)
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
		cssPurgeLog.Printf("could not update %s cache: %v", label, err)
	}
}

type htmlFingerprintResult struct {
	path     string
	sum      [sha256.Size]byte
	size     int64
	identity fileidentity.Identity
	err      error
}

type htmlFingerprintSummary struct {
	fingerprint  string
	scannedBytes int64
	identities   map[string]fileidentity.Identity
	fileHashes   map[string]string
	ok           bool
}

func fingerprintHTMLFiles(files []string, outputDir string, concurrency int) htmlFingerprintSummary {
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
				result := fingerprintHTMLFile(path)
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
			return htmlFingerprintSummary{scannedBytes: totalBytes}
		}
		totalBytes += result.size
		items = append(items, result)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].path < items[j].path })
	input := sha256.New()
	fileHashes := make(map[string]string, len(items))
	var identities map[string]fileidentity.Identity
	if fileidentity.Available {
		identities = make(map[string]fileidentity.Identity, len(items))
	}
	for _, item := range items {
		nativeRel, err := filepath.Rel(outputDir, item.path)
		if err != nil {
			nativeRel = item.path
		}
		rel := filepath.ToSlash(nativeRel)
		fileHashes[rel] = hex.EncodeToString(item.sum[:])
		_, _ = input.Write([]byte(rel))
		_, _ = input.Write([]byte{0})
		_, _ = input.Write(item.sum[:])
		if identities != nil {
			if item.identity.Inode == 0 || !filepath.IsLocal(nativeRel) {
				identities = nil
			} else if identities != nil {
				identities[filepath.ToSlash(rel)] = item.identity
			}
		}
	}
	if identities != nil && len(identities) != len(items) {
		identities = nil
	}
	return htmlFingerprintSummary{
		fingerprint:  hex.EncodeToString(input.Sum(nil)),
		scannedBytes: totalBytes,
		identities:   identities,
		fileHashes:   fileHashes,
		ok:           true,
	}
}

type htmlFingerprintReader struct{ io.Reader }

var cssPurgeFingerprintBufferPool = sync.Pool{
	New: func() any {
		buffer := make([]byte, 64<<10)
		return &buffer
	},
}

func fingerprintHTMLFile(path string) htmlFingerprintResult {
	result := htmlFingerprintResult{path: path}
	file, err := os.Open(path)
	if err != nil {
		result.err = err
		return result
	}
	var before os.FileInfo
	if fileidentity.Available {
		before, err = file.Stat()
		if err != nil {
			_ = file.Close()
			result.err = err
			return result
		}
	}
	hasher := sha256.New()
	buffer, ok := cssPurgeFingerprintBufferPool.Get().(*[]byte)
	if !ok || buffer == nil {
		fresh := make([]byte, 64<<10)
		buffer = &fresh
	}
	result.size, err = io.CopyBuffer(hasher, htmlFingerprintReader{Reader: file}, *buffer)
	cssPurgeFingerprintBufferPool.Put(buffer)
	var after os.FileInfo
	var afterErr error
	if fileidentity.Available {
		after, afterErr = file.Stat()
	}
	closeErr := file.Close()
	if err != nil {
		result.err = err
		return result
	}
	if afterErr != nil {
		result.err = afterErr
		return result
	}
	if closeErr != nil {
		result.err = closeErr
		return result
	}
	hasher.Sum(result.sum[:0])
	if !fileidentity.Available || !before.Mode().IsRegular() || !after.Mode().IsRegular() {
		return result
	}
	beforeIdentity := fileidentity.Read(before)
	afterIdentity := fileidentity.Read(after)
	pathInfo, err := os.Lstat(path)
	if err != nil || !pathInfo.Mode().IsRegular() {
		return result
	}
	pathIdentity := fileidentity.Read(pathInfo)
	if beforeIdentity.Inode != 0 && beforeIdentity == afterIdentity && afterIdentity == pathIdentity {
		result.identity = afterIdentity
	}
	return result
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
	perFile, complete, _ := scanHTMLFilesWithSelectors(files, concurrency, "", nil)
	combined := csspurge.NewUsedSelectors()
	for _, used := range perFile {
		combined.Merge(used)
	}
	return combined, complete
}

func scanHTMLFilesWithExpectedHashesConcurrently(files []string, outputDir string, expectedHashes map[string]string, concurrency int) (perFile map[string]*csspurge.UsedSelectors, complete, hashesMatch bool) {
	return scanHTMLFilesWithSelectors(files, concurrency, outputDir, expectedHashes)
}

func scanHTMLFilesWithSelectors(files []string, concurrency int, outputDir string, expectedHashes map[string]string) (perFile map[string]*csspurge.UsedSelectors, complete, hashesMatch bool) {
	if len(files) == 0 {
		return map[string]*csspurge.UsedSelectors{}, true, true
	}
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > len(files) {
		concurrency = len(files)
	}

	// Use channels for work distribution
	jobs := make(chan string, len(files))
	type result struct {
		path        string
		used        *csspurge.UsedSelectors
		err         error
		hashMatches bool
	}
	results := make(chan result, len(files))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				used := csspurge.NewUsedSelectors()
				if expectedHashes != nil {
					content, err := os.ReadFile(path)
					if err != nil {
						results <- result{path: path, err: fmt.Errorf("reading %s: %w", path, err)}
						continue
					}
					if err := csspurge.ScanHTMLContent(string(content), used); err != nil {
						results <- result{path: path, err: fmt.Errorf("scanning %s: %w", path, err)}
						continue
					}
					relative, err := filepath.Rel(outputDir, path)
					if err != nil || !filepath.IsLocal(relative) {
						results <- result{path: path, err: fmt.Errorf("invalid HTML path %s", path)}
						continue
					}
					expected := expectedHashes[filepath.ToSlash(relative)]
					actual := cssPurgeBytesHash(content)
					results <- result{path: path, used: used, hashMatches: actual == expected}
					continue
				}
				if err := csspurge.ScanHTML(path, used); err != nil {
					results <- result{path: path, err: fmt.Errorf("scanning %s: %w", path, err)}
					continue
				}
				results <- result{path: path, used: used, hashMatches: true}
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

	complete, hashesMatch = true, true
	perFile = make(map[string]*csspurge.UsedSelectors, len(files))
	for item := range results {
		if item.err != nil {
			// Log but don't fail - partial analysis is still useful.
			cssPurgeLog.Warnf("%v", item.err)
			complete = false
			continue
		}
		if !item.hashMatches {
			hashesMatch = false
		}
		perFile[item.path] = item.used
	}

	return perFile, complete, hashesMatch
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
