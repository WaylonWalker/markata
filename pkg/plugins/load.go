package plugins

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/encryption"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// Pre-compiled regex patterns for date normalization.
// These are compiled once at package init instead of per-call.
var (
	// timeFixRegex fixes malformed time components like "8:011:00" -> "08:11:00"
	timeFixRegex = regexp.MustCompile(`(\d{1,2}):0*(\d{1,2}):0*(\d{1,2})`)

	// singleDigitHourRegex normalizes single-digit hours in time component
	singleDigitHourRegex = regexp.MustCompile(`([ T])(\d):(\d{2})`)

	// startSingleDigitHourRegex matches time at start of string
	startSingleDigitHourRegex = regexp.MustCompile(`^\d:\d{2}`)
)

const sourceEncryptedPostKey = "_source_encrypted"

// LoadPlugin parses markdown files into Post objects.
type LoadPlugin struct {
	slugMode          string
	slugRules         []models.SlugRule
	encryptionEnabled bool
	encryptionDefault string
	encryptionTags    map[string]string
	encryptionKeys    map[string]string
}

// NewLoadPlugin creates a new LoadPlugin.
func NewLoadPlugin() *LoadPlugin {
	return &LoadPlugin{
		slugMode:       models.SlugModeFlat,
		encryptionTags: make(map[string]string),
		encryptionKeys: make(map[string]string),
	}
}

// Name returns the plugin identifier.
func (p *LoadPlugin) Name() string {
	return "load"
}

func (p *LoadPlugin) Configure(m *lifecycle.Manager) error {
	p.slugMode, p.slugRules = configuredSlugSettings(m.Config())
	p.configureSourceEncryption(m.Config())
	return nil
}

func (p *LoadPlugin) configureSourceEncryption(cfg *lifecycle.Config) {
	p.encryptionEnabled = false
	p.encryptionDefault = ""
	p.encryptionTags = make(map[string]string)
	p.encryptionKeys = make(map[string]string)

	if modelsConfig, ok := getModelsConfig(cfg); ok && modelsConfig != nil {
		p.encryptionEnabled = modelsConfig.Encryption.Enabled
		p.encryptionDefault = modelsConfig.Encryption.DefaultKey
		for tag, key := range modelsConfig.Encryption.PrivateTags {
			p.encryptionTags[strings.ToLower(tag)] = key
		}
	}

	if cfg != nil && cfg.Extra != nil {
		if enabled, ok := cfg.Extra["encryption_enabled"].(bool); ok {
			p.encryptionEnabled = enabled
		}
		if defaultKey, ok := cfg.Extra["encryption_default_key"].(string); ok {
			p.encryptionDefault = defaultKey
		}
	}

	if !p.encryptionEnabled {
		return
	}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, EncryptionEnvPrefix) {
			continue
		}
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 || parts[1] == "" {
			continue
		}
		keyName := strings.ToLower(strings.TrimPrefix(parts[0], EncryptionEnvPrefix))
		p.encryptionKeys[keyName] = parts[1]
	}
}

// Load reads and parses all discovered files into Post objects.
// Files are loaded in parallel using a worker pool for improved I/O performance.
// Uses ModTime-based caching to skip re-parsing unchanged files.
func (p *LoadPlugin) Load(m *lifecycle.Manager) error {
	resetSemanticHashBaseline(m)
	files := m.Files()
	config := m.Config()
	baseDir := config.ContentDir
	if baseDir == "" {
		baseDir = "."
	}

	// A normal build rescans the content globs, so the cache is the source of
	// ownership information for files that disappeared since the last build.
	// Keep those entries alive until the write stage can remove their outputs.
	if cache := GetBuildCache(m); cache != nil {
		missing := cache.MarkMissingPaths(files)
		if len(missing) > 0 {
			removed := append(lifecycle.GetServeRemovedPaths(m), missing...)
			lifecycle.SetServeRemovedPaths(m, removed)

			// Normal builds discover removals after the cache configure hook has
			// run. Compute dependent paths here so transform and write plugins
			// do not restore stale output for posts that referenced a deletion.
			changedSlugs := cache.GetChangedSlugs()
			affectedPaths := cache.GetAffectedPosts(changedSlugs)
			cache.MarkAffectedDependents(changedSlugs)
			affected := lifecycle.GetServeAffectedPaths(m)
			if affected == nil {
				affected = make(map[string]bool)
			}
			for _, path := range affectedPaths {
				affected[path] = true
			}
			lifecycle.SetServeAffectedPaths(m, affected)
		}
	}

	if cachedPosts := lifecycle.GetServeCachedPosts(m); len(cachedPosts) > 0 {
		if err := p.loadFromCachedPosts(m, files, baseDir, cachedPosts); err != nil {
			return err
		}
		recordDuplicateSlugReasons(m)
		return nil
	}

	// Get build cache for ModTime-based skipping
	cache := GetBuildCache(m)
	if lifecycle.IsServeFullRebuild(m) || cache == nil {
		if err := p.loadAllFiles(m, files, baseDir, cache, false); err != nil {
			return err
		}
		recordDuplicateSlugReasons(m)
		return nil
	}

	affected := lifecycle.GetServeAffectedPaths(m)
	if len(affected) == 0 {
		if err := p.loadAllFiles(m, files, baseDir, cache, true); err != nil {
			return err
		}
		recordDuplicateSlugReasons(m)
		return nil
	}

	restored, _, err := p.restoreFromCacheOrLoadChanged(m, files, baseDir, cache, affected)
	m.SetPosts(restored)
	recordDuplicateSlugReasons(m)
	if err != nil {
		return err
	}
	return nil
}

func (p *LoadPlugin) loadFromCachedPosts(
	m *lifecycle.Manager,
	files []string,
	baseDir string,
	cachedPosts map[string]*models.Post,
) error {
	cache := GetBuildCache(m)
	posts := make([]*models.Post, 0, len(files))
	affected := lifecycle.GetServeAffectedPaths(m)
	useFastCache := lifecycle.IsServeIncremental(m) && len(affected) > 0
	var firstOperationalErr error
	for _, file := range files {
		post, err := p.loadCachedPost(m, file, baseDir, cachedPosts, cache, useFastCache, affected)
		if err != nil {
			if !isRecoverableContentError(err) {
				recordLoadError(m, file, err)
				if firstOperationalErr == nil {
					firstOperationalErr = err
				}
			}
			continue
		}
		if post != nil {
			posts = append(posts, post)
		}
	}
	m.SetPosts(posts)
	return firstOperationalErr
}

func (p *LoadPlugin) loadCachedPost(
	m *lifecycle.Manager,
	file string,
	baseDir string,
	cachedPosts map[string]*models.Post,
	cache *buildcache.Cache,
	useFastCache bool,
	affected map[string]bool,
) (*models.Post, error) {
	if existing, ok := cachedPosts[file]; ok && !useFastCache {
		modTime, err := resolveFileModTime(m, file, baseDir)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s: %w", file, err)
		}
		if cache != nil {
			if cachedData := cache.GetCachedPostData(file, modTime); cachedData != nil {
				post := p.restorePostFromCache(cachedData, cache)
				recordLoadedPost(m, post)
				return post, nil
			}
		}
		recordLoadedPost(m, existing)
		return existing, nil
	}

	if useFastCache {
		if len(affected) > 0 && affected[file] {
			return p.loadFile(m, file, baseDir, cache)
		}
		if cache != nil {
			if cachedData := cache.GetCachedPostDataLatest(file); cachedData != nil {
				post := p.restorePostFromCache(cachedData, cache)
				recordLoadedPost(m, post)
				return post, nil
			}
		}
	}

	return p.loadFile(m, file, baseDir, cache)
}

func (p *LoadPlugin) loadAllFiles(m *lifecycle.Manager, files []string, baseDir string, cache *buildcache.Cache, useCache bool) error {
	if cache == nil {
		return p.loadSequential(m, files, baseDir, cache)
	}
	if !useCache {
		posts, err := p.loadFilesConcurrent(m, files, func(file string) (*models.Post, error) {
			return p.loadFile(m, file, baseDir, cache)
		})
		for _, post := range posts {
			m.AddPost(post)
		}
		return err
	}

	affected := lifecycle.GetServeAffectedPaths(m)
	useFastCache := lifecycle.IsServeIncremental(m) && len(affected) > 0
	posts, err := p.loadFilesConcurrent(m, files, func(file string) (*models.Post, error) {
		if useFastCache {
			if len(affected) > 0 && affected[file] {
				return p.loadFile(m, file, baseDir, cache)
			}
			if cachedData := cache.GetCachedPostDataLatest(file); cachedData != nil {
				post := p.restorePostFromCache(cachedData, cache)
				recordLoadedPost(m, post)
				return post, nil
			}
		}

		modTime, err := resolveFileModTime(m, file, baseDir)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s: %w", file, err)
		}
		if cachedData := cache.GetCachedPostData(file, modTime); cachedData != nil {
			post := p.restorePostFromCache(cachedData, cache)
			recordLoadedPost(m, post)
			return post, nil
		}
		return p.loadFile(m, file, baseDir, cache)
	})
	for _, post := range posts {
		m.AddPost(post)
	}
	return err
}

func (p *LoadPlugin) loadFilesConcurrent(m *lifecycle.Manager, files []string, fn func(string) (*models.Post, error)) ([]*models.Post, error) {
	if len(files) == 0 {
		return nil, nil
	}

	type result struct {
		index int
		post  *models.Post
		err   error
	}

	workers := m.Concurrency()
	if workers > len(files) {
		workers = len(files)
	}
	jobs := make(chan int, len(files))
	results := make(chan result, len(files))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				post, err := fn(files[idx])
				results <- result{index: idx, post: post, err: err}
			}
		}()
	}

	for i := range files {
		jobs <- i
	}
	close(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()

	ordered := make([]*models.Post, len(files))
	orderedErrors := make([]error, len(files))
	var firstOperationalErr error
	for res := range results {
		if res.err != nil {
			orderedErrors[res.index] = res.err
			continue
		}
		ordered[res.index] = res.post
	}
	for index, err := range orderedErrors {
		if err == nil || isRecoverableContentError(err) {
			continue
		}
		recordLoadError(m, files[index], err)
		if firstOperationalErr == nil {
			firstOperationalErr = err
		}
	}

	posts := make([]*models.Post, 0, len(ordered))
	for _, post := range ordered {
		if post != nil {
			posts = append(posts, post)
		}
	}
	return posts, firstOperationalErr
}

func (p *LoadPlugin) restoreFromCacheOrLoadChanged(
	m *lifecycle.Manager,
	files []string,
	baseDir string,
	cache *buildcache.Cache,
	affected map[string]bool,
) (posts, changedPosts []*models.Post, err error) {
	if cache == nil {
		return nil, nil, fmt.Errorf("build cache unavailable")
	}

	posts = make([]*models.Post, 0, len(files))
	var firstOperationalErr error

	for _, file := range files {
		if affected[file] {
			post, err := p.loadFile(m, file, baseDir, cache)
			if err != nil {
				if !isRecoverableContentError(err) {
					recordLoadError(m, file, err)
					if firstOperationalErr == nil {
						firstOperationalErr = err
					}
				}
				continue
			}
			posts = append(posts, post)
			changedPosts = append(changedPosts, post)
			continue
		}

		modTime, err := resolveFileModTime(m, file, baseDir)
		if err != nil {
			statErr := fmt.Errorf("failed to stat %s: %w", file, err)
			recordLoadError(m, file, statErr)
			if firstOperationalErr == nil {
				firstOperationalErr = statErr
			}
			continue
		}
		cachedData := cache.GetCachedPostData(file, modTime)
		if cachedData == nil {
			post, err := p.loadFile(m, file, baseDir, cache)
			if err != nil {
				if !isRecoverableContentError(err) {
					recordLoadError(m, file, err)
					if firstOperationalErr == nil {
						firstOperationalErr = err
					}
				}
				continue
			}
			posts = append(posts, post)
			continue
		}
		post := p.restorePostFromCache(cachedData, cache)
		recordLoadedPost(m, post)
		posts = append(posts, post)
	}

	return posts, changedPosts, firstOperationalErr
}

// loadFile loads a single file, using cache if ModTime is unchanged.
func (p *LoadPlugin) loadFile(m *lifecycle.Manager, file, baseDir string, cache *buildcache.Cache) (*models.Post, error) {
	fullPath := fullContentPath(file, baseDir)

	modTime, err := resolveFileModTime(m, file, baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to stat %s: %w", file, err)
	}

	// Read file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", file, err)
	}
	frontmatterAnalysis := recordSourceDiagnostics(m, file, string(content))
	sourceEncrypted := encryption.IsSourceEncrypted(sourceBodyForContent(string(content)))

	// Check cache for unchanged plaintext files. Source-encrypted files are
	// parsed every time so decrypted Markdown is never restored from disk cache.
	if cache != nil && !sourceEncrypted {
		if cachedData := cache.GetCachedPostData(file, modTime); cachedData != nil {
			post := p.restorePostFromCache(cachedData, cache)
			recordLoadedPost(m, post)
			return post, nil
		}
	}

	// Parse the file
	post, err := p.parseFile(file, string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", file, err)
	}

	// Capture before either parsed-cache storage or semantic updates can
	// mutate the prior build's entry. Empty hashes are a valid new-post baseline.
	captureSemanticHashBaseline(m, cache, file)

	// Cache the parsed post. Source-encrypted files intentionally skip parsed
	// post and article caches because those would contain decrypted content.
	if cache != nil && !isSourceEncryptedPost(post) {
		postData := p.postToCachedData(post, frontmatterAnalysis.Issues)
		//nolint:errcheck // caching is best-effort
		cache.CachePostData(file, modTime, postData)
		contentHash := buildcache.ContentHash(post.Content)
		//nolint:errcheck // caching is best-effort
		cache.CacheArticleHTML(file, contentHash, post.ArticleHTML)
		cache.SetPostSlug(file, post.Slug)
	}

	if cache != nil {
		feedHash := computePostFeedItemHash(post)
		tagHash := computePostTagIndexHash(post)
		gardenHash := computePostGardenHash(post)
		feedChanged, _, _ := cache.UpdatePostSemanticHashes(file, feedHash, tagHash, gardenHash)
		if feedChanged {
			cache.MarkFeedSlugChanged(post.Slug)
		}
	}
	recordLoadedPost(m, post)

	return post, nil
}

func recordSourceDiagnostics(m *lifecycle.Manager, path, content string) diagnostics.FrontmatterAnalysis {
	analysis := diagnostics.AnalyzeFrontmatter(path, content)
	if m == nil || m.ContentLedger() == nil {
		return analysis
	}

	ledger := m.ContentLedger()
	ledger.MarkLoaded(path)
	ledger.MarkFrontmatter(path, analysis.Inspection.HasFrontmatter, analysis.Valid)
	for _, issue := range analysis.Issues {
		ledger.AddIssue(issue)
	}
	return analysis
}

func recordLoadedPost(m *lifecycle.Manager, post *models.Post) {
	if m == nil || post == nil || post.Path == "" || m.ContentLedger() == nil {
		return
	}

	ledger := m.ContentLedger()
	ledger.MarkLoaded(post.Path)
	for _, issue := range cachedFrontmatterDiagnostics(post) {
		ledger.AddIssue(issue)
	}
	present, valid := cachedFrontmatterState(post)
	ledger.MarkFrontmatter(post.Path, present, valid)
	eligible := post.Published && !post.Draft && !post.Skip
	ledger.MarkPost(post.Path, eligible)
	if !post.Published {
		ledger.AddReason(post.Path, diagnostics.ReasonContentPublishedFalse)
	}
	if post.Draft {
		ledger.AddReason(post.Path, diagnostics.ReasonContentDraft)
	}
	if post.Skip {
		ledger.AddReason(post.Path, diagnostics.ReasonContentSkip)
	}
	if post.Private {
		ledger.AddReason(post.Path, diagnostics.ReasonContentPrivate)
	}
}

func cachedFrontmatterDiagnostics(post *models.Post) []diagnostics.Issue {
	if post == nil {
		return nil
	}
	issues, ok := post.Get("_frontmatter_issues").([]diagnostics.Issue)
	if !ok {
		return nil
	}
	return issues
}

func cachedFrontmatterState(post *models.Post) (present, valid bool) {
	if post == nil {
		return false, false
	}
	if value, ok := post.Get("_frontmatter_present").(bool); ok {
		present = value
	} else {
		present = post.RawFrontmatter != ""
	}
	if value, ok := post.Get("_frontmatter_valid").(bool); ok {
		valid = value
	} else {
		valid = present
	}
	return present, valid
}

func recordLoadError(m *lifecycle.Manager, path string, err error) {
	if m == nil || m.ContentLedger() == nil {
		return
	}
	if isRecoverableContentError(err) {
		return
	}
	ledger := m.ContentLedger()
	if ledger.HasIssueCode(path, diagnostics.ReasonContentLoadError) {
		return
	}
	// Keep source diagnostics safe for machine-readable output. In particular,
	// do not copy wrapped errors that can contain environment-derived values.
	ledger.RecordError(path, diagnostics.ReasonContentLoadError, "content could not be loaded")
}

func isRecoverableContentError(err error) bool {
	return errors.Is(err, ErrRecoverableContent)
}

func recordDuplicateSlugReasons(m *lifecycle.Manager) {
	if m == nil || m.ContentLedger() == nil {
		return
	}

	bySlug := make(map[string][]string)
	for _, post := range m.Posts() {
		if post == nil || post.Path == "" || post.Skip || post.Draft {
			continue
		}
		slug := strings.ToLower(post.Slug)
		bySlug[slug] = append(bySlug[slug], post.Path)
	}
	for _, paths := range bySlug {
		if len(paths) < 2 {
			continue
		}
		for _, path := range paths {
			m.ContentLedger().AddReason(path, diagnostics.ReasonContentDuplicateSlug)
		}
	}
}

func fullContentPath(file, baseDir string) string {
	if filepath.IsAbs(file) {
		return file
	}

	return filepath.Join(baseDir, file)
}

func resolveFileModTime(m *lifecycle.Manager, file, baseDir string) (int64, error) {
	if m != nil {
		if cached, ok := m.Cache().Get(cacheKeyGlobFileModTimes); ok {
			if modTimes, ok := cached.(map[string]int64); ok {
				if modTime, ok := modTimes[file]; ok && modTime != 0 {
					return modTime, nil
				}
			}
		}
	}

	stat, err := os.Stat(fullContentPath(file, baseDir))
	if err != nil {
		return 0, err
	}

	return stat.ModTime().UnixNano(), nil
}

// restorePostFromCache creates a Post from cached data.
func (p *LoadPlugin) restorePostFromCache(data *buildcache.CachedPostData, cache *buildcache.Cache) *models.Post {
	post := models.NewPost(data.Path)
	post.Content = data.Content
	post.Slug = data.Slug
	post.Href = data.Href
	post.Title = data.Title
	post.Date = data.Date
	post.Modified = data.Modified
	post.Published = data.Published
	post.Draft = data.Draft
	post.Private = data.Private
	post.PrivateOverride = data.PrivateOverride
	post.Skip = data.Skip
	post.Tags = data.Tags
	post.Description = data.Description
	post.Template = data.Template
	post.Templates = data.Templates
	post.RawFrontmatter = data.RawFrontmatter
	present := data.FrontmatterPresent
	valid := data.FrontmatterValid
	// Older cache entries did not store explicit frontmatter state. Infer the
	// state for them from the retained raw frontmatter when possible.
	if !present && data.RawFrontmatter != "" {
		present = true
		valid = true
	}
	post.InputHash = data.InputHash
	post.Authors = data.Authors
	post.Author = data.Author
	post.SecretKey = data.SecretKey // pragma: allowlist secret
	if data.Extra != nil {
		for k, v := range data.Extra {
			post.Set(k, v)
		}
	}
	// Set internal state after restoring Extra so older or hand-authored cache
	// entries cannot override the typed frontmatter state above.
	post.Set("_frontmatter_present", present)
	post.Set("_frontmatter_valid", valid)
	post.Set("_frontmatter_issues", append([]diagnostics.Issue{}, data.FrontmatterIssues...))

	if cache != nil {
		contentHash := buildcache.ContentHash(data.Content)
		if cachedHTML := cache.GetCachedArticleHTML(data.Path, contentHash); cachedHTML != "" {
			post.ArticleHTML = cachedHTML
		}
	}
	return post
}

func computePostFeedItemHash(post *models.Post) string {
	if post == nil {
		return ""
	}
	var b strings.Builder
	if post.Title != nil {
		b.WriteString(*post.Title)
	}
	b.WriteByte('\x00')
	b.WriteString(post.Slug)
	b.WriteByte('\x00')
	b.WriteString(post.Href)
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Private))
	b.WriteByte('\x00')
	if post.Date != nil {
		b.WriteString(post.Date.UTC().Format(time.RFC3339Nano))
	}
	b.WriteByte('\x00')
	if post.Description != nil {
		b.WriteString(*post.Description)
	}
	b.WriteByte('\x00')
	b.WriteString(post.Content)
	b.WriteByte('\x00')
	if len(post.Tags) > 0 {
		tags := append([]string(nil), post.Tags...)
		sort.Strings(tags)
		for _, tag := range tags {
			b.WriteString(tag)
			b.WriteByte('\x01')
		}
	}
	b.WriteByte('\x00')
	if post.Extra != nil {
		if img, ok := post.Extra["image"].(string); ok {
			b.WriteString(img)
		}
	}
	return buildcache.ContentHash(b.String())
}

func computePostTagIndexHash(post *models.Post) string {
	if post == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(post.Slug)
	b.WriteByte('\x00')
	if post.Title != nil {
		b.WriteString(*post.Title)
	}
	b.WriteByte('\x00')
	b.WriteString(post.Href)
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Published))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Draft))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Private))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Skip))
	b.WriteByte('\x00')
	if len(post.Tags) > 0 {
		tags := append([]string(nil), post.Tags...)
		sort.Strings(tags)
		for _, tag := range tags {
			b.WriteString(tag)
			b.WriteByte('\x01')
		}
	}
	return buildcache.ContentHash(b.String())
}

func computePostGardenHash(post *models.Post) string {
	if post == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(post.Slug)
	b.WriteByte('\x00')
	if post.Title != nil {
		b.WriteString(*post.Title)
	}
	b.WriteByte('\x00')
	b.WriteString(post.Href)
	b.WriteByte('\x00')
	if post.Date != nil {
		b.WriteString(post.Date.UTC().Format(time.RFC3339Nano))
	}
	b.WriteByte('\x00')
	if post.Description != nil {
		b.WriteString(*post.Description)
	}
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Published))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Draft))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Private))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", post.Skip))
	b.WriteByte('\x00')
	if len(post.Tags) > 0 {
		tags := append([]string(nil), post.Tags...)
		sort.Strings(tags)
		for _, tag := range tags {
			b.WriteString(tag)
			b.WriteByte('\x01')
		}
	}
	b.WriteByte('\x00')
	if len(post.Outlinks) > 0 {
		links := make([]string, 0, len(post.Outlinks))
		for _, link := range post.Outlinks {
			links = append(links, link.RawTarget)
		}
		sort.Strings(links)
		for _, href := range links {
			b.WriteString(href)
			b.WriteByte('\x01')
		}
	}
	return buildcache.ContentHash(b.String())
}

// postToCachedData converts a Post to cacheable data.
func (p *LoadPlugin) postToCachedData(post *models.Post, frontmatterIssues []diagnostics.Issue) *buildcache.CachedPostData {
	return &buildcache.CachedPostData{
		Path:               post.Path,
		Content:            post.Content,
		Slug:               post.Slug,
		Href:               post.Href,
		Title:              post.Title,
		Date:               post.Date,
		Modified:           post.Modified,
		Published:          post.Published,
		Draft:              post.Draft,
		Private:            post.Private,
		PrivateOverride:    post.PrivateOverride,
		Skip:               post.Skip,
		Tags:               post.Tags,
		Description:        post.Description,
		Template:           post.Template,
		Templates:          post.Templates,
		RawFrontmatter:     post.RawFrontmatter,
		FrontmatterPresent: post.Has("_frontmatter_present") && post.Get("_frontmatter_present") == true,
		FrontmatterValid:   post.Has("_frontmatter_valid") && post.Get("_frontmatter_valid") == true,
		FrontmatterIssues:  append([]diagnostics.Issue{}, frontmatterIssues...),
		InputHash:          post.InputHash,
		Authors:            post.Authors,
		Author:             post.Author,
		SecretKey:          post.SecretKey, // pragma: allowlist secret
		Extra:              post.Extra,
	}
}

// loadSequential loads files one at a time (used for small file counts).
func (p *LoadPlugin) loadSequential(m *lifecycle.Manager, files []string, baseDir string, cache *buildcache.Cache) error {
	var firstOperationalErr error
	for _, file := range files {
		post, err := p.loadFile(m, file, baseDir, cache)
		if err != nil {
			if !isRecoverableContentError(err) {
				recordLoadError(m, file, err)
				if firstOperationalErr == nil {
					firstOperationalErr = err
				}
			}
			continue
		}
		if post != nil {
			m.AddPost(post)
		}
	}

	return firstOperationalErr
}

// ParsePostFromContent parses a markdown file content into a Post.
// This is a lightweight helper for cache-backed workflows that need
// frontmatter parsing and slug generation without running the full lifecycle.
func ParsePostFromContent(path, content string) (*models.Post, error) {
	return ParsePostFromContentWithConfig(path, content, nil)
}

func ParsePostFromContentWithConfig(path, content string, cfg *models.Config) (*models.Post, error) {
	loader := NewLoadPlugin()
	if cfg != nil {
		loader.slugMode = models.NormalizeSlugMode(cfg.GlobConfig.SlugMode)
		loader.slugRules = append([]models.SlugRule{}, cfg.GlobConfig.SlugRules...)
		loader.encryptionEnabled = cfg.Encryption.Enabled
		loader.encryptionDefault = cfg.Encryption.DefaultKey
		loader.encryptionTags = make(map[string]string, len(cfg.Encryption.PrivateTags))
		for tag, key := range cfg.Encryption.PrivateTags {
			loader.encryptionTags[strings.ToLower(tag)] = key
		}
		if loader.encryptionEnabled {
			for _, env := range os.Environ() {
				if !strings.HasPrefix(env, EncryptionEnvPrefix) {
					continue
				}
				parts := strings.SplitN(env, "=", 2)
				if len(parts) == 2 && parts[1] != "" {
					loader.encryptionKeys[strings.ToLower(strings.TrimPrefix(parts[0], EncryptionEnvPrefix))] = parts[1]
				}
			}
		}
	}
	post, err := loader.parseFile(path, content)
	if err != nil {
		return post, err
	}

	// Standalone callers do not run the transform stage. Complete the same
	// title derivation here even when parseFile already derived a filename slug.
	if post.Title != nil {
		m := lifecycle.NewManager()
		if cfg != nil {
			m.Config().Extra = cfg.Extra
		}
		renderer := NewRenderMarkdownPlugin()
		if err := renderer.Configure(m); err != nil {
			return nil, err
		}
		m.AddPost(post)
		if err := NewInlineTitlesPlugin().Transform(m); err != nil {
			return nil, err
		}
	}

	// Only generate a title-based fallback when parseFile did not already
	// produce a slug and the user did not explicitly provide one.
	if post.Slug == "" && !post.Has("_slug_explicit") {
		post.GenerateSlugWithMode(models.SlugModeForPath(path, loader.slugMode, loader.slugRules))
		post.GenerateHref()
	}
	return post, nil
}

// parseFile parses a markdown file's content into a Post object.
func (p *LoadPlugin) parseFile(path, content string) (*models.Post, error) {
	inspection := diagnostics.InspectFrontmatter(content)

	// Parse frontmatter and get raw frontmatter for hashing
	metadata, body, rawFrontmatter, err := ParseFrontmatterWithRaw(content)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRecoverableContent, err)
	}

	// Create post with defaults
	post := models.NewPost(path)
	post.Content = body
	post.RawFrontmatter = rawFrontmatter
	post.Set("_frontmatter_present", inspection.HasFrontmatter)
	post.Set("_frontmatter_valid", inspection.HasFrontmatter)

	// Apply metadata to post
	if err := p.applyMetadata(post, metadata); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRecoverableContent, err)
	}

	decryptedBody, err := p.decryptSourceBody(post, body)
	if err != nil {
		return nil, err
	}
	if decryptedBody != body {
		body = decryptedBody
		post.Content = decryptedBody
	}

	// Generate slug if not explicitly set in frontmatter
	// If slug was explicitly set (even to empty string), respect it
	if !post.Has("_slug_explicit") && post.Slug == "" {
		candidate := models.NewPost(path)
		candidate.GenerateSlugWithMode(models.SlugModeForPath(path, p.slugMode, p.slugRules))
		if candidate.Slug != "" {
			post.Slug = candidate.Slug
		}
	}

	// Generate href from slug
	post.GenerateHref()

	// Compute input hash (content + frontmatter + template)
	// Template may be resolved later, so we use what we have now
	post.InputHash = buildcache.ComputePostInputHash(body, rawFrontmatter, post.Template)

	return post, nil
}

func sourceBodyForContent(content string) string {
	_, body, err := ExtractFrontmatter(content)
	if err != nil {
		return content
	}
	return body
}

func isSourceEncryptedPost(post *models.Post) bool {
	value, ok := post.Get(sourceEncryptedPostKey).(bool)
	return ok && value
}

func (p *LoadPlugin) decryptSourceBody(post *models.Post, body string) (string, error) {
	envelope, encrypted, err := encryption.ParseSourceEnvelope(body)
	if err != nil {
		return "", err
	}
	if !encrypted {
		return body, nil
	}

	post.Set(sourceEncryptedPostKey, true)
	if !p.encryptionEnabled {
		return body, nil
	}

	p.applySourcePrivateTags(post)
	keyName := strings.TrimSpace(envelope.KeyName)
	if keyName != "" {
		post.SecretKey = keyName // pragma: allowlist secret
		post.Private = true
	} else {
		keyName = p.resolveSourceKeyName(post)
	}
	if keyName == "" {
		return "", fmt.Errorf("encrypted source body for %s has no key; set secret_key or encryption.default_key", post.Path)
	}

	password, ok := p.encryptionKeys[strings.ToLower(keyName)]
	if !ok || password == "" {
		return "", fmt.Errorf("encrypted source body for %s requires %s%s", post.Path, EncryptionEnvPrefix, strings.ToUpper(keyName))
	}

	decrypted, _, err := encryption.DecryptSourceMarkdown(body, password)
	if err != nil {
		return "", fmt.Errorf("decrypt encrypted source body for %s with key %q: %w", post.Path, keyName, err)
	}

	post.SecretKey = keyName // pragma: allowlist secret
	post.Private = true
	return decrypted, nil
}

func (p *LoadPlugin) resolveSourceKeyName(post *models.Post) string {
	if post.SecretKey != "" {
		return post.SecretKey
	}
	if p.encryptionDefault != "" {
		return p.encryptionDefault
	}
	return ""
}

func (p *LoadPlugin) applySourcePrivateTags(post *models.Post) {
	if len(p.encryptionTags) == 0 || post == nil || post.Skip || post.Draft {
		return
	}
	if post.IsExplicitlyPublic() {
		return
	}
	for _, tag := range post.Tags {
		if keyName, ok := p.encryptionTags[strings.ToLower(tag)]; ok {
			post.Private = true
			if post.SecretKey == "" {
				post.SecretKey = keyName // pragma: allowlist secret
			}
			return
		}
	}
	if post.Template == "" {
		return
	}
	if keyName, ok := p.encryptionTags[strings.ToLower(post.Template)]; ok {
		post.Private = true
		if post.SecretKey == "" {
			post.SecretKey = keyName // pragma: allowlist secret
		}
	}
}

func configuredSlugSettings(cfg *lifecycle.Config) (string, []models.SlugRule) {
	if cfg == nil {
		return models.SlugModeFlat, nil
	}
	if modelsConfig, ok := cfg.Extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		return models.NormalizeSlugMode(modelsConfig.GlobConfig.SlugMode), append([]models.SlugRule{}, modelsConfig.GlobConfig.SlugRules...)
	}
	if raw, ok := cfg.Extra["slug_mode"].(string); ok {
		return models.NormalizeSlugMode(raw), nil
	}
	return models.SlugModeFlat, nil
}

func configuredSlugMode(cfg *lifecycle.Config, path string) string {
	mode, rules := configuredSlugSettings(cfg)
	return models.SlugModeForPath(path, mode, rules)
}

// Date field aliases for publication date (first match wins).
// Precedence follows common SSG conventions (publishdate first).
var dateFieldAliases = []string{"publishdate", "date", "pubdate"}

// Modified field aliases for last modified date (first match wins).
// Precedence follows common SSG conventions (lastmod first).
var modifiedFieldAliases = []string{"lastmod", "modified", "updated", "updated_at", "last_modified"}

// applyMetadata applies parsed frontmatter metadata to a Post.
func (p *LoadPlugin) applyMetadata(post *models.Post, metadata map[string]interface{}) error {
	// Known fields to extract
	knownFields := map[string]bool{
		"title":                 true,
		"date":                  true,
		"published":             true,
		"draft":                 true,
		"private":               true,
		"skip":                  true,
		"tags":                  true,
		"description":           true,
		"template":              true,
		"templates":             true,
		"slug":                  true,
		"secret_key":            true,
		"private_key":           true,
		"encryption_key":        true,
		"author":                true,
		"authors":               true,
		"by":                    true,
		"writer":                true,
		"modified":              true,
		"lastmod":               true,
		"updated":               true,
		"last_modified":         true,
		"updated_at":            true,
		"publishdate":           true,
		"pubdate":               true,
		"_title_explicit":       true,
		"_description_explicit": true,
	}

	// Title
	if title := GetString(metadata, "title"); title != "" {
		post.Title = &title
		post.Set("_title_explicit", true)
	}

	if err := applyDateFields(post, metadata); err != nil {
		return err
	}

	// Published
	post.Published = GetBool(metadata, "published", false)

	// Draft
	post.Draft = GetBool(metadata, "draft", false)

	// Private
	if err := applyPrivateMetadata(post, metadata); err != nil {
		return err
	}

	// Skip
	post.Skip = GetBool(metadata, "skip", false)

	// Tags
	if tags := GetStringSlice(metadata, "tags"); tags != nil {
		post.Tags = tags
	}

	// Description
	if desc, ok := metadata["description"]; ok {
		if value, ok := desc.(string); ok {
			cleaned := strings.TrimSpace(value)
			post.Description = &cleaned
			post.Set("_description_explicit", true)
		}
	}

	// Template - support both 'template' and 'templateKey' for Python markata compatibility
	if template := GetString(metadata, "template"); template != "" {
		post.Template = template
	} else if template := GetString(metadata, "templateKey"); template != "" {
		post.Template = template
	}

	// Templates - per-format template overrides
	if templatesVal, ok := metadata["templates"]; ok {
		post.Templates = parseTemplatesMap(templatesVal)
	}

	// Slug - support custom slugs including explicit empty string for homepage
	if slugVal, exists := metadata["slug"]; exists {
		slug := normalizeCustomSlug(GetString(metadata, "slug"))
		post.Slug = slug
		// Mark that slug was explicitly set (prevents auto-generation)
		post.Set("_slug_explicit", true)
		_ = slugVal // Exists check used, value handled via GetString
	}

	// SecretKey - for encrypted posts  // pragma: allowlist secret
	// Supports aliases: secret_key, private_key, encryption_key (first non-empty wins)
	if secretKey := GetString(metadata, "secret_key"); secretKey != "" {
		post.SecretKey = secretKey // pragma: allowlist secret
	} else if privateKey := GetString(metadata, "private_key"); privateKey != "" {
		post.SecretKey = privateKey // pragma: allowlist secret
	} else if encryptionKey := GetString(metadata, "encryption_key"); encryptionKey != "" {
		post.SecretKey = encryptionKey // pragma: allowlist secret
	}

	// Author fields - support both legacy 'author' and new 'authors' arrays
	// Also support aliases: 'by' and 'writer' map to 'author'
	// SetAuthors handles: strings, []string, and mixed []interface{} with {id, role} maps
	if authorsVal, ok := metadata["authors"]; ok {
		post.SetAuthors(authorsVal)
	} else if author := GetString(metadata, "author"); author != "" {
		// Legacy single author field
		post.Author = &author
	} else if by := GetString(metadata, "by"); by != "" {
		// Alias: 'by' -> 'author'
		post.Author = &by
	} else if writer := GetString(metadata, "writer"); writer != "" {
		// Alias: 'writer' -> 'author'
		post.Author = &writer
	}

	// Store unknown fields in Extra
	for key, value := range metadata {
		if !knownFields[key] {
			post.Set(key, value)
		}
	}

	return nil
}

func applyPrivateMetadata(post *models.Post, metadata map[string]interface{}) error {
	post.Private = false
	post.PrivateOverride = nil
	rawPrivate, ok := metadata["private"]
	if !ok {
		return nil
	}
	private, ok := rawPrivate.(bool)
	if !ok {
		return fmt.Errorf("frontmatter field private must be a boolean")
	}
	post.Private = private
	post.PrivateOverride = &private
	return nil
}

func applyDateFields(post *models.Post, metadata map[string]interface{}) error {
	if dateVal := getFirstDateValue(metadata, dateFieldAliases); dateVal != nil {
		date, err := parseDate(dateVal)
		if err != nil {
			return fmt.Errorf("invalid date: %w", err)
		}
		post.Date = &date
	}

	if modVal := getFirstDateValue(metadata, modifiedFieldAliases); modVal != nil {
		modified, err := parseDate(modVal)
		if err != nil {
			return fmt.Errorf("invalid modified date: %w", err)
		}
		post.Modified = &modified
	}

	return nil
}

// parseDate attempts to parse a date value from various formats.
func parseDate(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		return parseDateString(v)
	default:
		return time.Time{}, fmt.Errorf("unsupported date type: %T", value)
	}
}

// parseDateString parses a date string using common formats.
func parseDateString(s string) (time.Time, error) {
	// Normalize the date string first
	s = normalizeDateString(s)

	// Common date formats to try (most specific first)
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02 15:04",
		"2006/01/02",
		"01/02/2006 15:04:05",
		"01/02/2006 15:04",
		"01/02/2006",
		"02-01-2006",
		"January 2, 2006",
		"Jan 2, 2006",
		"2 January 2006",
		"2 Jan 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// normalizeDateString normalizes a date string to handle common variations.
// It handles:
// - Single-digit hours/minutes (e.g., "1:00:00" -> "01:00:00")
// - Malformed time components (e.g., "8:011:00" -> "8:11:00")
// - Extra whitespace
func normalizeDateString(s string) string {
	// Trim whitespace
	s = strings.TrimSpace(s)

	// Early return if no time component
	if !strings.Contains(s, ":") {
		return s
	}

	// Fix malformed time components like "8:011:00" -> "08:11:00"
	// This regex finds time components with potentially extra leading zeros
	s = timeFixRegex.ReplaceAllStringFunc(s, func(match string) string {
		parts := timeFixRegex.FindStringSubmatch(match)
		if len(parts) == 4 {
			h, err := strconv.Atoi(parts[1])
			if err != nil {
				return match
			}
			m, err := strconv.Atoi(parts[2])
			if err != nil {
				return match
			}
			sec, err := strconv.Atoi(parts[3])
			if err != nil {
				return match
			}
			return fmt.Sprintf("%02d:%02d:%02d", h, m, sec)
		}
		return match
	})

	// Normalize single-digit hours in time component
	// Match patterns like " 1:00" or "T1:00" and pad the hour
	s = singleDigitHourRegex.ReplaceAllString(s, "${1}0${2}:${3}")

	// Handle time at start of string or after date with space
	if startSingleDigitHourRegex.MatchString(s) {
		s = "0" + s
	}

	return s
}

// normalizeCustomSlug normalizes a slug from frontmatter.
// It handles:
//   - "/" or "" -> "" (homepage)
//   - "/docs/page" -> "docs/page" (strip leading slash)
//   - "docs/page/" -> "docs/page" (strip trailing slash)
//   - Preserves internal structure for nested paths
func normalizeCustomSlug(slug string) string {
	// Trim whitespace
	slug = strings.TrimSpace(slug)

	// "/" means homepage (empty slug)
	if slug == "/" {
		return ""
	}

	// Strip leading and trailing slashes
	slug = strings.Trim(slug, "/")

	return slug
}

// parseTemplatesMap parses the templates field from frontmatter.
// It handles both map[string]interface{} and map[string]string formats.
func parseTemplatesMap(val interface{}) map[string]string {
	result := make(map[string]string)

	switch v := val.(type) {
	case map[string]interface{}:
		for key, value := range v {
			if str, ok := value.(string); ok {
				result[key] = str
			}
		}
	case map[string]string:
		for key, value := range v {
			result[key] = value
		}
	case map[interface{}]interface{}:
		// Handle YAML's default map type
		for key, value := range v {
			if keyStr, ok := key.(string); ok {
				if valStr, ok := value.(string); ok {
					result[keyStr] = valStr
				}
			}
		}
	}

	return result
}

// getFirstDateValue returns the first non-nil value from metadata using a list of alias keys.
// Returns nil if no value is found for any of the keys.
func getFirstDateValue(metadata map[string]interface{}, aliases []string) interface{} {
	for _, key := range aliases {
		if val, ok := metadata[key]; ok && val != nil {
			return val
		}
	}
	return nil
}
