package plugins

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/bmatcuk/doublestar/v4"
)

const cacheKeyGlobFileModTimes = "glob.file_mod_times"
const cacheKeyGlobFileInfo = "glob.file_info"

// GlobFileInfo describes filesystem metadata collected during a glob scan.
type GlobFileInfo struct {
	ModTime int64
	Size    int64
}

// GlobFileInfoMap returns metadata collected during the most recent glob
// scan so later stages do not need to stat every matched file again.
func GlobFileInfoMap(m *lifecycle.Manager) map[string]GlobFileInfo {
	if m == nil {
		return nil
	}
	value, ok := m.Cache().Get(cacheKeyGlobFileInfo)
	if !ok {
		return nil
	}
	modTimes, ok := value.(map[string]GlobFileInfo)
	if !ok {
		return nil
	}
	result := make(map[string]GlobFileInfo, len(modTimes))
	for path, info := range modTimes {
		result[path] = info
	}
	return result
}

// GlobPlugin discovers content files using glob patterns.
type GlobPlugin struct {
	// patterns are the glob patterns to match files against.
	// Supports ** for recursive matching (doublestar patterns).
	patterns []string

	// useGitignore determines whether to parse and respect .gitignore.
	useGitignore bool

	// gitignorePatterns holds parsed gitignore patterns.
	gitignorePatterns []string
	ignoreRules       []ignoreRule
}

type ignoreRule struct {
	pattern  string
	filename string
	meta     bool
	pathRule bool
}

// NewGlobPlugin creates a new GlobPlugin with default settings.
func NewGlobPlugin() *GlobPlugin {
	return &GlobPlugin{
		patterns:     []string{"**/*.md"},
		useGitignore: true,
	}
}

// Name returns the plugin identifier.
func (p *GlobPlugin) Name() string {
	return "glob"
}

// Configure reads configuration from the manager and initializes the plugin.
func (p *GlobPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Get glob patterns from config
	if len(config.GlobPatterns) > 0 {
		p.patterns = config.GlobPatterns
	}

	// Check for useGitignore setting in Extra config
	if extra := config.Extra; extra != nil {
		if useGitignore, ok := extra["use_gitignore"].(bool); ok {
			p.useGitignore = useGitignore
		}
	}

	// Parse .gitignore if enabled
	if p.useGitignore {
		if err := p.loadGitignore(config.ContentDir); err != nil {
			// Don't fail if .gitignore doesn't exist
			if !os.IsNotExist(err) {
				return err
			}
		}
		p.prepareIgnoreRules()
	}

	return nil
}

func (p *GlobPlugin) prepareIgnoreRules() {
	p.ignoreRules = make([]ignoreRule, 0, len(p.gitignorePatterns))
	for _, pattern := range p.gitignorePatterns {
		pattern = filepath.ToSlash(strings.TrimSuffix(pattern, "/"))
		p.ignoreRules = append(p.ignoreRules, ignoreRule{
			pattern:  pattern,
			filename: filepath.Base(pattern),
			meta:     strings.ContainsAny(pattern, "*?[{"),
			pathRule: strings.Contains(pattern, "/"),
		})
	}
}

// loadGitignore reads and parses .gitignore patterns.
func (p *GlobPlugin) loadGitignore(baseDir string) error {
	gitignorePath := filepath.Join(baseDir, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return err
	}
	defer file.Close()

	p.gitignorePatterns = make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		p.gitignorePatterns = append(p.gitignorePatterns, line)
	}

	return scanner.Err()
}

// isIgnored checks if a path matches any gitignore pattern.
func (p *GlobPlugin) isIgnored(path string) bool {
	if !p.useGitignore || len(p.ignoreRules) == 0 {
		return false
	}

	// Normalize path separators
	normalizedPath := filepath.ToSlash(path)
	filename := filepath.Base(normalizedPath)

	for _, rule := range p.ignoreRules {
		pattern := rule.pattern
		// Handle negation patterns (patterns starting with !)
		if strings.HasPrefix(pattern, "!") {
			continue // Skip negation for now in ignore check
		}

		// Most .gitignore entries are literal directory or file names. Avoid
		// running doublestar's pattern validator for those entries on every
		// matched file. This is particularly important for large sites with
		// generated output and dependency directories.
		if !rule.meta {
			if normalizedPath == pattern || strings.HasPrefix(normalizedPath, pattern+"/") ||
				filename == rule.filename || strings.HasSuffix(normalizedPath, "/"+pattern) ||
				strings.HasSuffix(normalizedPath, "/"+pattern+"/") {
				return true
			}
			continue
		}

		// Try different matching strategies

		// 1. Direct match with the pattern
		if doublestar.MatchUnvalidated(pattern, normalizedPath) {
			return true
		}

		// 2. Pattern as prefix (for directory patterns)
		if strings.HasPrefix(normalizedPath, pattern+"/") {
			return true
		}

		// 3. Patterns containing a path separator cannot match a basename
		// (except **/ rules, which intentionally match at any depth).
		if !rule.pathRule || strings.HasPrefix(pattern, "**/") {
			if doublestar.MatchUnvalidated(pattern, filename) {
				return true
			}
		}

		// 4. Try with **/ prefix for patterns that should match anywhere
		if !strings.HasPrefix(pattern, "**/") && !strings.HasPrefix(pattern, "/") {
			if doublestar.MatchUnvalidated("**/"+pattern, normalizedPath) {
				return true
			}
		}
	}

	return false
}

// Glob discovers content files matching the configured patterns.
// Uses cached file list when patterns haven't changed.
func (p *GlobPlugin) Glob(m *lifecycle.Manager) error {
	config := m.Config()
	baseDir := config.ContentDir
	if baseDir == "" {
		baseDir = "."
	}

	// Convert to absolute path for consistent matching
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return err
	}

	// Check for cached file list
	cache := GetBuildCache(m)
	patternHash := buildcache.HashContent(strings.Join(p.patterns, "\n"))

	if lifecycle.IsServeFullRebuild(m) {
		files, modTimes, err := p.scanFiles(absBaseDir)
		if err != nil {
			return &contentDiscoveryError{err: err}
		}
		if cache != nil && len(files) > 0 {
			cache.SetGlobCache(files, patternHash)
		}
		setGlobMetadata(m, modTimes)
		m.SetFiles(files)
		return nil
	}

	if cache != nil {
		cachedFiles, cachedHash := cache.GetGlobCache()
		if cachedHash == patternHash && len(cachedFiles) > 0 && shouldReuseCachedGlobFiles(m) {
			m.SetFiles(cachedFiles)
			return nil
		}
	}

	// Full scan
	files, modTimes, err := p.scanFiles(absBaseDir)
	if err != nil {
		return &contentDiscoveryError{err: err}
	}

	// Cache for next build
	if cache != nil && len(files) > 0 {
		cache.SetGlobCache(files, patternHash)
	}

	setGlobMetadata(m, modTimes)
	m.SetFiles(files)
	return nil
}

// contentDiscoveryError preserves the underlying filesystem error for
// callers while keeping absolute discovery paths out of CLI error output.
type contentDiscoveryError struct {
	err error
}

func (e *contentDiscoveryError) Error() string {
	return "content discovery failed"
}

func (e *contentDiscoveryError) Unwrap() error {
	return e.err
}

func setGlobMetadata(m *lifecycle.Manager, info map[string]GlobFileInfo) {
	modTimes := make(map[string]int64, len(info))
	for path, fileInfo := range info {
		modTimes[path] = fileInfo.ModTime
	}
	m.Cache().Set(cacheKeyGlobFileModTimes, modTimes)
	m.Cache().Set(cacheKeyGlobFileInfo, info)
}

func shouldReuseCachedGlobFiles(m *lifecycle.Manager) bool {
	if !lifecycle.IsServeIncremental(m) || lifecycle.IsServeGlobDirty(m) {
		return false
	}

	return len(lifecycle.GetServeChangedPaths(m)) > 0 ||
		len(lifecycle.GetServeRemovedPaths(m)) > 0 ||
		len(lifecycle.GetServeAffectedPaths(m)) > 0
}

// scanFiles performs full glob scan and records file modtimes for later stages.
func (p *GlobPlugin) scanFiles(absBaseDir string) ([]string, map[string]GlobFileInfo, error) {
	for _, pattern := range p.patterns {
		if filepath.IsAbs(pattern) {
			return p.scanFilesWithGlob(absBaseDir)
		}
	}

	fileSet := make(map[string]struct{})
	modTimes := make(map[string]GlobFileInfo)

	err := filepath.WalkDir(absBaseDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == absBaseDir {
			return nil
		}

		relPath, err := filepath.Rel(absBaseDir, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if p.isIgnored(relPath) {
				return filepath.SkipDir
			}
			return nil
		}
		// Apply the content glob before the ignore rules. The old glob walker
		// only presented matching files to isIgnored; preserving that order is
		// important because .gitignore files can contain hundreds of rules.
		if !p.matchesAnyPattern(relPath) || p.isIgnored(relPath) {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}
		fileSet[relPath] = struct{}{}
		modTimes[relPath] = GlobFileInfo{ModTime: info.ModTime().UnixNano(), Size: info.Size()}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	files := make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	sort.Strings(files)
	return files, modTimes, nil
}

func (p *GlobPlugin) scanFilesWithGlob(absBaseDir string) ([]string, map[string]GlobFileInfo, error) {
	fileSet := make(map[string]struct{})
	modTimes := make(map[string]GlobFileInfo)

	for _, pattern := range p.patterns {
		fullPattern := pattern
		if !filepath.IsAbs(pattern) {
			fullPattern = filepath.Join(absBaseDir, pattern)
		}

		matches, err := doublestar.FilepathGlob(fullPattern)
		if err != nil {
			return nil, nil, err
		}

		for _, match := range matches {
			relPath, err := filepath.Rel(absBaseDir, match)
			if err != nil {
				relPath = match
			}
			if p.isIgnored(relPath) {
				continue
			}

			info, err := os.Stat(match)
			if err != nil {
				return nil, nil, err
			}
			if info.IsDir() {
				continue
			}
			fileSet[relPath] = struct{}{}
			modTimes[relPath] = GlobFileInfo{ModTime: info.ModTime().UnixNano(), Size: info.Size()}
		}
	}

	return sortedGlobFiles(fileSet), modTimes, nil
}

func (p *GlobPlugin) matchesAnyPattern(path string) bool {
	path = filepath.ToSlash(path)
	for _, pattern := range p.patterns {
		if doublestar.MatchUnvalidated(filepath.ToSlash(pattern), path) {
			return true
		}
	}
	return false
}

func sortedGlobFiles(fileSet map[string]struct{}) []string {
	files := make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	sort.Strings(files)
	return files
}

// SetPatterns sets the glob patterns to use for file discovery.
func (p *GlobPlugin) SetPatterns(patterns []string) {
	p.patterns = patterns
}

// SetUseGitignore enables or disables gitignore support.
func (p *GlobPlugin) SetUseGitignore(use bool) {
	p.useGitignore = use
}
