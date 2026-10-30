// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// MentionsPlugin transforms @handle syntax into HTML links
// by resolving handles against blogroll entries.
type MentionsPlugin struct {
	// cssClass is the CSS class applied to mention links
	cssClass string
	// metadataHTTPClient is an optional test/client override. A nil client
	// preserves the production timeout and HTTP instrumentation defaults.
	metadataHTTPClient *http.Client
}

const mentionsCacheVersion = "mentions:v2"

// NewMentionsPlugin creates a new MentionsPlugin.
func NewMentionsPlugin() *MentionsPlugin {
	return &MentionsPlugin{
		cssClass: "mention",
	}
}

// Name returns the unique name of the plugin.
func (p *MentionsPlugin) Name() string {
	return "mentions"
}

// Configure reads configuration options for the plugin.
func (p *MentionsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// First check the new Mentions config
	mentionsConfig := getMentionsConfig(config)
	if mentionsConfig.CSSClass != "" {
		p.cssClass = mentionsConfig.CSSClass
	}

	// Fall back to legacy Extra config for backwards compatibility
	if config.Extra != nil {
		if cssClass, ok := config.Extra["mentions_css_class"].(string); ok && cssClass != "" {
			p.cssClass = cssClass
		}
	}
	return nil
}

// Priority returns the plugin's priority for a given stage.
// Runs after blogroll has cached feed data in Collect stage.
func (p *MentionsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		// Run late in transform stage, after blogroll config is available
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Transform processes @mentions in all post content.
func (p *MentionsPlugin) Transform(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		if disabled, ok := config.Extra["mentions_disabled"].(bool); ok && disabled {
			return nil
		}
	}

	// Build handle resolution map from blogroll config
	handleMap := p.buildHandleMap(m)

	log.Printf("[mentions] Found %d handles in handleMap", len(handleMap))

	// Keep processing when the blogroll is empty: unknown ~attributions
	// still need distinct styling and are caught by the offline linter.

	// Only contacts referenced by loaded posts need network metadata. Keep
	// the full handle map for unknown mentions and author rendering.
	referenced := p.referencedHandles(m, handleMap)

	// A post with no relevant handles must not contact unrelated sites.
	domainMap := p.fetchAllMetadata(m, referenced)

	// Store in memory cache for other plugins
	m.Cache().Set("mentions_metadata", domainMap)

	// Attach metadata to mention entries
	p.attachMetadataToEntries(handleMap, domainMap)

	cache := GetBuildCache(m)
	cacheSignature := p.computeCacheSignature(handleMap)
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.Content != ""
	})

	log.Printf("[mentions] Processing %d posts", len(posts))

	var needProcessing []*models.Post
	if cache != nil {
		for _, post := range posts {
			contentHash := p.computePostHash(post, cacheSignature)
			if cachedContent, ok := cache.GetCachedMentionsContent(post.Path, contentHash); ok {
				post.Content = cachedContent
				continue
			}
			needProcessing = append(needProcessing, post)
		}
	} else {
		needProcessing = posts
	}

	if len(needProcessing) == 0 {
		return nil
	}

	return m.ProcessPostsSliceConcurrently(needProcessing, func(post *models.Post) error {
		contentHash := p.computePostHash(post, cacheSignature)
		content := p.processChatAdmonitionTitles(post, handleMap)
		content = p.processMentionsWithMetadata(content, handleMap)
		post.Content = content
		if cache != nil {
			cache.CacheMentionsContent(post.Path, contentHash, content)
		}
		return nil
	})
}

// referencedHandles selects contacts that the rendering passes can actually
// link, without fetching unrelated blogroll entries or fenced-code literals.
func (p *MentionsPlugin) referencedHandles(m *lifecycle.Manager, handleMap map[string]*mentionEntry) map[string]*mentionEntry {
	referenced := make(map[string]*mentionEntry)
	addHandle := func(handle string) {
		normalized := strings.ToLower(handle)
		if entry := handleMap[normalized]; entry != nil {
			referenced[normalized] = entry
			return
		}
		// Preserve exact domain handles, but retry without sentence punctuation.
		normalized = mentionTrailingPunctRegex.ReplaceAllString(normalized, "")
		if entry := handleMap[normalized]; entry != nil {
			referenced[normalized] = entry
		}
	}
	for _, post := range m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.Content != ""
	}) {
		// Match the rendering pipeline: fenced code is inert, while chat
		// titles are an exception to ordinary skipped admonition headers.
		content := mentionsCodeBlockRegex.ReplaceAllString(post.Content, "")
		for _, line := range strings.Split(content, "\n") {
			if admonitionHeaderRegex.MatchString(line) {
				groups := chatAdmonitionTitleRegex.FindStringSubmatch(line)
				if len(groups) >= 7 {
					if groups[3] != "" {
						addHandle(groups[3])
					} else {
						addHandle(groups[5])
					}
				}
				continue
			}
			for _, match := range mentionRegex.FindAllStringSubmatch(line, -1) {
				addHandle(match[3])
			}
		}
	}
	return referenced
}

func (p *MentionsPlugin) computePostHash(post *models.Post, cacheSignature string) string {
	authorName, authorAvatar := getPrimaryAuthorDisplay(post)
	var b strings.Builder
	b.WriteString(post.Content)
	b.WriteByte('\x00')
	b.WriteString(cacheSignature)
	b.WriteByte('\x00')
	b.WriteString(authorName)
	b.WriteByte('\x00')
	b.WriteString(authorAvatar)
	return buildcache.ContentHash(b.String())
}

func (p *MentionsPlugin) computeCacheSignature(handleMap map[string]*mentionEntry) string {
	type cacheMetadata struct {
		Name   string `json:"name,omitempty"`
		Bio    string `json:"bio,omitempty"`
		Avatar string `json:"avatar,omitempty"`
		URL    string `json:"url,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	type cacheEntry struct {
		Key      string         `json:"key"`
		Handle   string         `json:"handle"`
		SiteURL  string         `json:"site_url"`
		Title    string         `json:"title,omitempty"`
		Internal bool           `json:"internal,omitempty"`
		Metadata *cacheMetadata `json:"metadata,omitempty"`
	}

	keys := make([]string, 0, len(handleMap))
	for key := range handleMap {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	entries := make([]cacheEntry, 0, len(keys))
	for _, key := range keys {
		entry := handleMap[key]
		if entry == nil {
			continue
		}

		cachedEntry := cacheEntry{
			Key:      key,
			Handle:   entry.Handle,
			SiteURL:  entry.SiteURL,
			Title:    entry.Title,
			Internal: entry.Internal,
		}
		if entry.Metadata != nil {
			cachedEntry.Metadata = &cacheMetadata{
				Name:   entry.Metadata.Name,
				Bio:    entry.Metadata.Bio,
				Avatar: entry.Metadata.Avatar,
				URL:    entry.Metadata.URL,
				Error:  entry.Metadata.Error,
			}
		}
		entries = append(entries, cachedEntry)
	}

	payload := struct {
		Version  string       `json:"version"`
		CSSClass string       `json:"css_class"`
		Entries  []cacheEntry `json:"entries"`
	}{
		Version:  mentionsCacheVersion,
		CSSClass: p.cssClass,
		Entries:  entries,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return mentionsCacheVersion
	}

	return buildcache.ContentHash(string(data))
}

// mentionEntry holds resolved information for a handle.
type mentionEntry struct {
	Handle   string
	SiteURL  string
	Title    string
	Internal bool // true for from_posts entries, false for blogroll entries
	Metadata *models.MentionMetadata
}

// registerFeedConfig registers a feed config's handle and aliases in the map.
func (p *MentionsPlugin) registerFeedConfig(feedConfig *models.ExternalFeedConfig, handleMap map[string]*mentionEntry) {
	if !feedConfig.IsActive() {
		return
	}

	// Determine the site URL
	siteURL := feedConfig.SiteURL
	if siteURL == "" {
		siteURL = extractSiteURL(feedConfig.URL)
	}
	siteURL = normalizeExternalURL(siteURL)

	if siteURL == "" {
		return
	}

	// Determine the handle
	handle := feedConfig.Handle
	if handle == "" {
		handle = extractHandleFromURL(siteURL)
	}

	if handle == "" {
		return
	}

	handle = strings.ToLower(handle)

	// Create and store the entry
	entry := &mentionEntry{
		Handle:  handle,
		SiteURL: siteURL,
		Title:   feedConfig.Title,
	}

	if _, exists := handleMap[handle]; !exists {
		handleMap[handle] = entry
	}

	// Auto-register domain alias
	domain := extractDomainFromURL(siteURL)
	if domain != "" && domain != handle {
		if _, exists := handleMap[domain]; !exists {
			handleMap[domain] = entry
		}
	}

	// Register manual aliases
	for _, alias := range feedConfig.Aliases {
		p.registerAlias(alias, entry, handleMap)
	}
}

// registerAlias registers an alias in the handle map.
func (p *MentionsPlugin) registerAlias(alias string, entry *mentionEntry, handleMap map[string]*mentionEntry) {
	normalizedAlias := strings.ToLower(alias)
	if normalizedAlias == "" {
		return
	}
	if _, exists := handleMap[normalizedAlias]; exists {
		log.Printf("[mentions] Warning: duplicate alias %q (first entry wins)", normalizedAlias)
		return
	}
	handleMap[normalizedAlias] = entry
}

// registerAuthors registers authors from the site configuration as mentionable contacts.
// Authors with a URL are registered using their config key (ID) as the handle.
// Author metadata (name, bio, avatar) is used for hovercards.
func (p *MentionsPlugin) registerAuthors(config *lifecycle.Config, handleMap map[string]*mentionEntry) {
	modelsConfig, ok := getModelsConfig(config)
	if !ok {
		return
	}

	authorMap := modelsConfig.Authors.Authors
	if len(authorMap) == 0 {
		return
	}

	registered := 0
	for id := range authorMap {
		author := authorMap[id]
		// Authors need a URL to be linkable mentions
		if author.URL == nil || *author.URL == "" {
			continue
		}

		handle := strings.ToLower(id)

		// Build metadata from author fields
		metadata := &models.MentionMetadata{
			Name: author.Name,
			URL:  *author.URL,
		}
		if author.Bio != nil {
			metadata.Bio = *author.Bio
		}
		if author.Avatar != nil {
			metadata.Avatar = normalizeExternalURL(*author.Avatar)
		}

		entry := &mentionEntry{
			Handle:   handle,
			SiteURL:  *author.URL,
			Title:    author.Name,
			Internal: true,
			Metadata: metadata,
		}

		// Register the handle (first entry wins)
		if _, exists := handleMap[handle]; !exists {
			handleMap[handle] = entry
			registered++
		}
	}

	if registered > 0 {
		log.Printf("[mentions] Registered %d authors as mentionable contacts", registered)
	}
}

// registerCachedFeed registers a cached feed's handle in the map.
func (p *MentionsPlugin) registerCachedFeed(feed *models.ExternalFeed, handleMap map[string]*mentionEntry) {
	if feed.SiteURL == "" {
		return
	}
	normalizeExternalFeed(feed)

	handle := feed.Config.Handle
	if handle == "" {
		handle = extractHandleFromURL(feed.SiteURL)
	}

	if handle == "" {
		return
	}

	handle = strings.ToLower(handle)

	if _, exists := handleMap[handle]; !exists {
		handleMap[handle] = &mentionEntry{
			Handle:  handle,
			SiteURL: feed.SiteURL,
			Title:   feed.Title,
		}
	}
}

// buildHandleMap builds a map of handles to their site URLs from blogroll config,
// configured authors, and internal posts configured via from_posts.
// Resolution order:
// 1. Explicit handle from blogroll config
// 2. Auto-generated handle from domain
// 3. Authors from site configuration
// 4. Internal posts matching from_posts filters
func (p *MentionsPlugin) buildHandleMap(m *lifecycle.Manager) map[string]*mentionEntry {
	handleMap := make(map[string]*mentionEntry)

	config := m.Config()
	blogrollConfig := getBlogrollConfig(config)
	mentionsConfig := getMentionsConfig(config)

	log.Printf("[mentions] Blogroll config enabled: %v, feeds count: %d", blogrollConfig.Enabled, len(blogrollConfig.Feeds))

	// Register from blogroll if enabled
	if blogrollConfig.Enabled {
		// Register feed configs
		for i := range blogrollConfig.Feeds {
			p.registerFeedConfig(&blogrollConfig.Feeds[i], handleMap)
		}

		// Register cached feeds
		if cachedFeeds, ok := m.Cache().Get("blogroll_feeds"); ok {
			if feeds, ok := cachedFeeds.([]*models.ExternalFeed); ok {
				for _, feed := range feeds {
					p.registerCachedFeed(feed, handleMap)
				}
			}
		}
	}

	// Register authors from site configuration
	p.registerAuthors(config, handleMap)

	// Register from internal posts (from_posts sources)
	p.registerFromPosts(m, mentionsConfig, handleMap)

	return handleMap
}

// cacheKey generates a cache filename for a domain.
func (p *MentionsPlugin) cacheKey(domain string) string {
	// Normalize domain and replace dots with underscores for safe filenames
	return strings.ReplaceAll(strings.ToLower(domain), ".", "_")
}

// loadFromCache loads mention metadata from the file cache.
// Returns nil if cache doesn't exist, is expired, or has errors.
func (p *MentionsPlugin) loadFromCache(domain, cacheDir string, maxAge time.Duration) *models.MentionMetadata {
	cacheFile := filepath.Join(cacheDir, p.cacheKey(domain)+".json")

	if _, err := os.Stat(cacheFile); err != nil {
		return nil // Cache file doesn't exist
	}

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}

	var metadata models.MentionMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil
	}
	normalizeMentionMetadata(&metadata)

	// Additional validation
	if metadata.IsExpired(maxAge) {
		return nil
	}

	return &metadata
}

// saveToCache saves mention metadata to the file cache.
func (p *MentionsPlugin) saveToCache(metadata *models.MentionMetadata, cacheDir string) error {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return err
	}

	cacheFile := filepath.Join(cacheDir, p.cacheKey(metadata.Domain)+".json")
	normalizeMentionMetadata(metadata)

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0o600)
}

// fetchMetadata fetches metadata for a single domain, using cache if available.
func (p *MentionsPlugin) fetchMetadata(domain, cacheDir string, maxAge, timeout time.Duration) *models.MentionMetadata {
	// Try cache first
	if cached := p.loadFromCache(domain, cacheDir, maxAge); cached != nil {
		return cached
	}

	// Cache miss or expired, fetch from network
	url := "https://" + domain
	metadata := &models.MentionMetadata{
		Domain:      domain,
		URL:         url,
		LastFetched: time.Now(),
	}

	client := p.metadataHTTPClient
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	buildstats.InstrumentHTTPClient(client)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		metadata.Error = fmt.Sprintf("failed to create request: %v", err)
		//nolint:errcheck // Cache even errors to prevent repeated failed requests
		p.saveToCache(metadata, cacheDir)
		return metadata
	}

	// Set user agent to be respectful
	req.Header.Set("User-Agent", "markata-go/1.0 mentions-plugin")

	resp, err := client.Do(req)
	if err != nil {
		metadata.Error = fmt.Sprintf("HTTP request failed: %v", err)
		//nolint:errcheck // Cache even errors to prevent repeated failed requests
		p.saveToCache(metadata, cacheDir)
		return metadata
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metadata.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		//nolint:errcheck // Cache even errors to prevent repeated failed requests
		p.saveToCache(metadata, cacheDir)
		return metadata
	}

	// Parse HTML to extract metadata
	if err := p.extractMetadataFromHTML(resp, metadata); err != nil {
		metadata.Error = fmt.Sprintf("failed to parse HTML: %v", err)
		//nolint:errcheck // Cache even errors to prevent repeated failed requests
		p.saveToCache(metadata, cacheDir)
		return metadata
	}

	// Cache successful metadata
	if err := p.saveToCache(metadata, cacheDir); err != nil {
		log.Printf("mentions: failed to cache metadata for %s: %v", domain, err)
	}

	return metadata
}

// extractMetadataFromHTML parses HTML response to extract mention metadata.
func (p *MentionsPlugin) extractMetadataFromHTML(resp *http.Response, metadata *models.MentionMetadata) error {
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	htmlContent := string(body)

	// Extract title/name
	if name := p.extractMetaContent(htmlContent, []*regexp.Regexp{
		metaOGSiteNameRe,
		metaOGTitleRe,
		metaAuthorRe,
		metaTitleRe,
	}); name != "" {
		metadata.Name = name
	} else {
		metadata.Name = metadata.Domain // fallback to domain
	}

	// Extract description/bio
	if bio := p.extractMetaContent(htmlContent, []*regexp.Regexp{
		metaOGDescriptionRe,
		metaDescriptionRe,
	}); bio != "" {
		metadata.Bio = bio
	}

	// Extract avatar/image
	if avatar := p.extractMetaContent(htmlContent, []*regexp.Regexp{
		metaOGImageRe,
		metaIconRe,
		metaAppleTouchIconRe,
	}); avatar != "" {
		metadata.Avatar = p.resolveURL(avatar, metadata.URL)
	}
	normalizeMentionMetadata(metadata)

	return nil
}

// extractMetaContent extracts content from HTML using multiple pre-compiled regex patterns.
// Returns the first match from the patterns list.
func (p *MentionsPlugin) extractMetaContent(htmlContent string, patterns []*regexp.Regexp) string {
	for _, re := range patterns {
		matches := re.FindStringSubmatch(htmlContent)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}
	return ""
}

// resolveURL makes relative URLs absolute based on the base URL.
func (p *MentionsPlugin) resolveURL(imageURL, baseURL string) string {
	if imageURL == "" {
		return ""
	}

	// If already absolute, return as-is
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return normalizeExternalURL(imageURL)
	}

	// Make relative URLs absolute
	if strings.HasPrefix(imageURL, "/") {
		if parsed, err := url.Parse(baseURL); err == nil {
			return normalizeExternalURL(parsed.Scheme + "://" + parsed.Host + imageURL)
		}
	}

	// Return original if can't resolve
	return normalizeExternalURL(imageURL)
}

// fetchAllMetadata fetches metadata for all unique domains concurrently.
func (p *MentionsPlugin) fetchAllMetadata(m *lifecycle.Manager, handleMap map[string]*mentionEntry) map[string]*models.MentionMetadata {
	config := getMentionsConfig(m.Config())

	// Extract unique domains
	domains := p.extractUniqueDomains(handleMap)
	if len(domains) == 0 {
		return make(map[string]*models.MentionMetadata)
	}

	// Parse cache duration
	cacheDuration, err := time.ParseDuration(config.GetCacheDuration())
	if err != nil {
		cacheDuration = 24 * time.Hour // fallback to 24h
	}

	timeout := time.Duration(config.GetTimeout()) * time.Second
	cacheDir := config.GetCacheDir()

	// Concurrent fetching with semaphore
	semaphore := make(chan struct{}, config.GetConcurrentRequests())
	var wg sync.WaitGroup
	metadataMap := make(map[string]*models.MentionMetadata)
	var mu sync.RWMutex

	for _, domain := range domains {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			metadata := p.fetchMetadata(d, cacheDir, cacheDuration, timeout)
			mu.Lock()
			metadataMap[d] = metadata
			mu.Unlock()
		}(domain)
	}

	wg.Wait()
	return metadataMap
}

// extractUniqueDomains extracts unique domains from handleMap.
func (p *MentionsPlugin) extractUniqueDomains(handleMap map[string]*mentionEntry) []string {
	domainSet := make(map[string]bool)
	for _, entry := range handleMap {
		if entry.SiteURL != "" {
			if parsed, err := url.Parse(entry.SiteURL); err == nil {
				domain := strings.ToLower(parsed.Hostname())
				if domain != "" {
					domainSet[domain] = true
				}
			}
		}
	}

	domains := make([]string, 0, len(domainSet))
	for domain := range domainSet {
		domains = append(domains, domain)
	}
	return domains
}

// attachMetadataToEntries attaches fetched metadata to mention entries.
// Entries that already have metadata (e.g., from internal posts) are skipped.
func (p *MentionsPlugin) attachMetadataToEntries(handleMap map[string]*mentionEntry, domainMap map[string]*models.MentionMetadata) {
	for _, entry := range handleMap {
		// Skip entries that already have metadata (e.g., from internal posts)
		if entry.Metadata != nil {
			continue
		}
		if entry.SiteURL != "" {
			if parsed, err := url.Parse(entry.SiteURL); err == nil {
				domain := strings.ToLower(parsed.Hostname())
				if metadata, exists := domainMap[domain]; exists {
					entry.Metadata = metadata
				}
			}
		}
	}
}

// Pre-compiled regexes for HTML metadata extraction.
var (
	metaOGSiteNameRe     = regexp.MustCompile(`<meta[^>]+property=["']og:site_name["'][^>]+content=["']([^"']+)["']`)
	metaOGTitleRe        = regexp.MustCompile(`<meta[^>]+property=["']og:title["'][^>]+content=["']([^"']+)["']`)
	metaAuthorRe         = regexp.MustCompile(`<meta[^>]+name=["']author["'][^>]+content=["']([^"']+)["']`)
	metaTitleRe          = regexp.MustCompile(`<title>([^<]+)</title>`)
	metaOGDescriptionRe  = regexp.MustCompile(`<meta[^>]+property=["']og:description["'][^>]+content=["']([^"']+)["']`)
	metaDescriptionRe    = regexp.MustCompile(`<meta[^>]+name=["']description["'][^>]+content=["']([^"']+)["']`)
	metaOGImageRe        = regexp.MustCompile(`<meta[^>]+property=["']og:image["'][^>]+content=["']([^"']+)["']`)
	metaIconRe           = regexp.MustCompile(`<link[^>]+rel=["']icon["'][^>]+href=["']([^"']+)["']`)
	metaAppleTouchIconRe = regexp.MustCompile(`<link[^>]+rel=["']apple-touch-icon["'][^>]+href=["']([^"']+)["']`)
)

// mentionRegex matches @handle references and ~handle quote attributions.
// Handles can contain alphanumeric characters, underscores, hyphens, and dots.
// This supports both simple handles like @daverupert and domain-style handles
// like @simonwillison.net. Must start with a letter and not be preceded by
// another @ or word character.
var mentionRegex = regexp.MustCompile(`((?:^|[^@~\w])([@~])([a-zA-Z][a-zA-Z0-9_.-]*))([^a-zA-Z0-9_.-]|$)`)

// mentionTrailingPunctRegex strips trailing punctuation from handles.
// When the greedy handle match includes trailing dots, commas, etc. that are
// actually sentence punctuation (e.g., "@contact." at end of sentence),
// we strip them and retry the lookup. This preserves domain-style handles
// like @simonwillison.net when they exist in the handle map.
var mentionTrailingPunctRegex = regexp.MustCompile(`[.,;:!?]+$`)

// mentionsCodeBlockRegex matches fenced code blocks to avoid transforming mentions inside them.
var mentionsCodeBlockRegex = regexp.MustCompile("(?s)(```[^`]*```|~~~[^~]*~~~)")

// admonitionHeaderRegex matches admonition header lines (!!!, ???, ???+).
// These lines must be skipped during general mention processing because
// @handles in admonition titles are either:
// 1. Already enriched by processChatAdmonitionTitles (for chat/chat-reply), or
// 2. Part of the admonition title syntax that goldmark needs to parse intact.
// Transforming @handles on these lines into <a> tags would break the admonition parser.
var admonitionHeaderRegex = regexp.MustCompile(`(?m)^(?:\?{3}\+?|!!!)\s+\S`)

// processMentionsWithMetadata replaces @handle syntax with HTML anchor tags including metadata.
func (p *MentionsPlugin) processMentionsWithMetadata(content string, handleMap map[string]*mentionEntry) string {
	// Split content by fenced code blocks to avoid transforming mentions inside them
	codeBlocks := mentionsCodeBlockRegex.FindAllStringIndex(content, -1)

	if len(codeBlocks) == 0 {
		return p.processMentionsInText(content, handleMap)
	}

	// Process content in segments, skipping code blocks
	var result strings.Builder
	lastEnd := 0

	for _, block := range codeBlocks {
		start, end := block[0], block[1]

		// Process text before this code block
		if start > lastEnd {
			processed := p.processMentionsInText(content[lastEnd:start], handleMap)
			result.WriteString(processed)
		}

		// Keep code block unchanged
		result.WriteString(content[start:end])
		lastEnd = end
	}

	// Process any remaining text after the last code block
	if lastEnd < len(content) {
		processed := p.processMentionsInText(content[lastEnd:], handleMap)
		result.WriteString(processed)
	}

	return result.String()
}

// processMentionsInText processes @mentions in a text segment (not inside code blocks).
// Admonition header lines are skipped to avoid breaking admonition title parsing.
func (p *MentionsPlugin) processMentionsInText(text string, handleMap map[string]*mentionEntry) string {
	// Find admonition header lines to skip
	admonitionLines := admonitionHeaderRegex.FindAllStringIndex(text, -1)

	// If there are admonition lines, process text in segments
	if len(admonitionLines) > 0 {
		return p.processMentionsSkippingAdmonitions(text, handleMap, admonitionLines)
	}

	return p.replaceMentionsInText(text, handleMap)
}

// processMentionsSkippingAdmonitions processes mentions while skipping admonition header lines.
func (p *MentionsPlugin) processMentionsSkippingAdmonitions(text string, handleMap map[string]*mentionEntry, admonitionLines [][]int) string {
	// Expand each admonition match to cover the full line
	skipRanges := make([][]int, 0, len(admonitionLines))
	for _, loc := range admonitionLines {
		lineStart := loc[0]
		// Find end of line
		lineEnd := strings.Index(text[lineStart:], "\n")
		if lineEnd == -1 {
			lineEnd = len(text)
		} else {
			lineEnd = lineStart + lineEnd
		}
		skipRanges = append(skipRanges, []int{lineStart, lineEnd})
	}

	var result strings.Builder
	lastEnd := 0

	for _, r := range skipRanges {
		start, end := r[0], r[1]

		// Process text before this admonition line
		if start > lastEnd {
			processed := p.replaceMentionsInText(text[lastEnd:start], handleMap)
			result.WriteString(processed)
		}

		// Keep admonition line unchanged
		result.WriteString(text[start:end])
		lastEnd = end
	}

	// Process any remaining text
	if lastEnd < len(text) {
		processed := p.replaceMentionsInText(text[lastEnd:], handleMap)
		result.WriteString(processed)
	}

	return result.String()
}

// replaceMentionsInText replaces @mentions with HTML links in a text segment.
func (p *MentionsPlugin) replaceMentionsInText(text string, handleMap map[string]*mentionEntry) string {
	return mentionRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract the handle from the match
		// Groups: [0]=full match, [1]=prefix+marker+handle, [2]=marker,
		// [3]=handle, [4]=suffix.
		groups := mentionRegex.FindStringSubmatch(match)
		if len(groups) < 5 {
			return match
		}

		marker := groups[2]
		handle := strings.ToLower(groups[3])
		suffix := groups[4]
		prefix, _, _ := strings.Cut(match, marker)

		// Look up the handle
		entry, found := handleMap[handle]
		if !found {
			// Try stripping trailing punctuation from the handle.
			// The regex greedily captures dots/commas in handles (e.g., "@contact."
			// at end of sentence captures handle as "contact."). Strip trailing
			// punctuation and retry. This preserves domain-style handles like
			// @simonwillison.net because their exact match succeeds above.
			stripped := mentionTrailingPunctRegex.ReplaceAllString(handle, "")
			if stripped != handle && stripped != "" {
				entry, found = handleMap[stripped]
				if found {
					// Move the stripped punctuation to the suffix
					suffix = handle[len(stripped):] + suffix
					handle = stripped
				}
			}
			if !found {
				if marker == "~" {
					return prefix + fmt.Sprintf(`<span class=%q data-handle=%q>%s%s</span>`,
						html.EscapeString(p.cssClass+" mention--attribution mention--unresolved"),
						html.EscapeString("~"+handle), marker, html.EscapeString(handle)) + suffix
				}
				return match
			}
		}

		// Build data attributes if metadata is available
		dataAttrs := ""
		if entry.Metadata != nil && entry.Metadata.IsValid() {
			dataAttrs += fmt.Sprintf(` data-name=%q`, html.EscapeString(entry.Metadata.Name))
			if entry.Metadata.Bio != "" {
				dataAttrs += fmt.Sprintf(` data-bio=%q`, html.EscapeString(entry.Metadata.Bio))
			}
			if entry.Metadata.Avatar != "" {
				dataAttrs += fmt.Sprintf(` data-avatar=%q`, html.EscapeString(entry.Metadata.Avatar))
			}
			dataAttrs += fmt.Sprintf(` data-handle=%q`, html.EscapeString("@"+entry.Handle))
		}

		// Quote attributions get a distinct, themeable class and retain the
		// source marker in the visible link text.
		cssClass := p.cssClass
		if marker == "~" {
			cssClass += " mention--attribution"
		}
		link := fmt.Sprintf(`<a href=%q class=%q%s>%s%s</a>`,
			html.EscapeString(entry.SiteURL),
			html.EscapeString(cssClass),
			dataAttrs,
			marker,
			html.EscapeString(entry.Handle))

		return prefix + link + suffix
	})
}

// chatAdmonitionTitleRegex matches chat/chat-reply admonition lines whose title
// starts with @handle. Supports both quoted and unquoted forms:
//   - !!! chat "@alice"   (quoted)
//   - !!! chat @alice     (unquoted)
//
// Captures:
// Group 1: marker and type prefix (e.g., '!!! chat ')
// Group 2: opening quote ("") or empty string
// Group 3: handle (e.g., 'alice')
// Group 4: closing quote and anything after, or end-of-line remainder
var chatAdmonitionTitleRegex = regexp.MustCompile(`(?m)^((?:\?\?\?\+?|!!!)\s+(?:chat|chat-reply)\s+)(?:(")@([a-zA-Z][a-zA-Z0-9_.-]*)(".*)|@([a-zA-Z][a-zA-Z0-9_.-]*)(.*))$`)

// chatReplyAdmonitionRegex matches chat-reply lines without explicit titles.
// Example: !!! chat-reply
var chatReplyAdmonitionRegex = regexp.MustCompile(`(?m)^((?:\?{3}\+?|!!!)\s+chat-reply)\s*$`)

// processChatAdmonitionTitles finds chat/chat-reply admonition lines with @handle
// titles and replaces the title with enriched HTML containing the contact's avatar
// and a linked mention. Supports both quoted ("@handle") and unquoted (@handle)
// forms. The enriched title is always emitted unquoted so the admonition parser
// captures it via its unquoted-title group; wrapping in quotes would break because
// the HTML attributes contain internal double-quote characters.

func (p *MentionsPlugin) processChatAdmonitionTitles(post *models.Post, handleMap map[string]*mentionEntry) string {
	if post == nil {
		return ""
	}

	content := post.Content
	if content == "" {
		return content
	}

	authorName, authorAvatar := getPrimaryAuthorDisplay(post)
	content = chatAdmonitionTitleRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := chatAdmonitionTitleRegex.FindStringSubmatch(match)
		if len(groups) < 7 {
			return match
		}

		prefix := groups[1] // e.g., "!!! chat "

		// Determine quoted vs unquoted form:
		// Quoted:   groups[2]=`"`, groups[3]=handle, groups[4]=`"...`
		// Unquoted: groups[5]=handle, groups[6]=rest
		var handle, suffix string
		if groups[2] != "" {
			// Quoted form: !!! chat "@alice"
			handle = groups[3]
			// Drop the quotes; any trailing content after the closing quote is kept
			rest := strings.TrimPrefix(groups[4], `"`)
			suffix = strings.TrimSpace(rest)
		} else {
			// Unquoted form: !!! chat @alice
			handle = groups[5]
			suffix = strings.TrimSpace(groups[6])
		}

		normalizedHandle := strings.ToLower(handle)
		entry, found := handleMap[normalizedHandle]
		var enrichedTitle string
		switch {
		case found:
			// Build enriched title HTML with avatar and mention link.
			// Emit as unquoted title so the admonition parser captures it via (.*)
			enrichedTitle = p.buildChatEnrichedTitle(entry)
		case isChatReplyPrefix(prefix) && authorName != "":
			enrichedTitle = p.buildChatReplyTitle(authorName, authorAvatar)
		default:
			return match
		}

		if suffix != "" {
			return prefix + enrichedTitle + "\n" + suffix
		}
		return prefix + enrichedTitle + "\n"
	})

	if authorName == "" {
		return content
	}

	return chatReplyAdmonitionRegex.ReplaceAllStringFunc(content, func(match string) string {
		groups := chatReplyAdmonitionRegex.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}
		prefix := groups[1]
		enrichedTitle := p.buildChatReplyTitle(authorName, authorAvatar)
		return prefix + " " + enrichedTitle + "\n"
	})
}

func (p *MentionsPlugin) buildChatReplyTitle(name, avatar string) string {
	var titleHTML strings.Builder

	if avatar != "" {
		titleHTML.WriteString(fmt.Sprintf(`<img class="chat-contact-avatar" src=%q alt=%q />`,
			html.EscapeString(avatar),
			html.EscapeString(name)))
	}

	titleHTML.WriteString(fmt.Sprintf(`<span class="chat-contact-name">%s</span>`, html.EscapeString(name)))

	return fmt.Sprintf(`<span class="chat-contact">%s</span>`, titleHTML.String())
}

func isChatReplyPrefix(prefix string) bool {
	normalized := strings.TrimSpace(prefix)
	return strings.Contains(normalized, "chat-reply")
}

func getPrimaryAuthorDisplay(post *models.Post) (name, avatar string) {
	if post == nil {
		return "", ""
	}
	if len(post.AuthorObjects) > 0 {
		author := post.AuthorObjects[0]
		name := author.Name
		avatar := ""
		if author.Avatar != nil {
			avatar = *author.Avatar
		}
		return name, avatar
	}
	if post.Author != nil && *post.Author != "" {
		return *post.Author, ""
	}
	return "", ""
}

// buildChatEnrichedTitle builds the enriched HTML for a chat admonition title.
func (p *MentionsPlugin) buildChatEnrichedTitle(entry *mentionEntry) string {
	var titleHTML strings.Builder

	// Avatar image
	if entry.Metadata != nil && entry.Metadata.Avatar != "" {
		titleHTML.WriteString(fmt.Sprintf(`<img class="chat-contact-avatar" src=%q alt=%q />`,
			html.EscapeString(entry.Metadata.Avatar),
			html.EscapeString(entry.Metadata.Name)))
	}

	// Build data attributes for hovercard
	dataAttrs := ""
	if entry.Metadata != nil && entry.Metadata.IsValid() {
		dataAttrs += fmt.Sprintf(` data-name=%q`, html.EscapeString(entry.Metadata.Name))
		if entry.Metadata.Bio != "" {
			dataAttrs += fmt.Sprintf(` data-bio=%q`, html.EscapeString(entry.Metadata.Bio))
		}
		if entry.Metadata.Avatar != "" {
			dataAttrs += fmt.Sprintf(` data-avatar=%q`, html.EscapeString(entry.Metadata.Avatar))
		}
		dataAttrs += fmt.Sprintf(` data-handle=%q`, html.EscapeString("@"+entry.Handle))
	}

	// Mention link
	titleHTML.WriteString(fmt.Sprintf(` <a href=%q class=%q%s>@%s</a>`,
		html.EscapeString(entry.SiteURL),
		html.EscapeString(p.cssClass),
		dataAttrs,
		html.EscapeString(entry.Handle)))

	// Wrap in a span for styling
	return fmt.Sprintf(`<span class="chat-contact">%s</span>`, titleHTML.String())
}

// extractSiteURL extracts the base site URL from a feed URL.
func extractSiteURL(feedURL string) string {
	if feedURL == "" {
		return ""
	}

	parsed, err := url.Parse(feedURL)
	if err != nil {
		return ""
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}

	// Return scheme + host
	return fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
}

// extractHandleFromURL extracts a handle from a URL's domain.
// For example:
// - "https://daverupert.com" -> "daverupert"
// - "https://www.example.com" -> "example"
// - "https://blog.jane.dev" -> "jane"
func extractHandleFromURL(siteURL string) string {
	parsed, err := url.Parse(siteURL)
	if err != nil {
		return ""
	}

	host := parsed.Hostname()
	if host == "" {
		return ""
	}

	// Remove common prefixes
	host = strings.TrimPrefix(host, "www.")
	host = strings.TrimPrefix(host, "blog.")

	// Get the first part of the domain (before the TLD)
	parts := strings.Split(host, ".")
	if len(parts) == 0 {
		return ""
	}

	// Take the first part (subdomain or main domain name)
	handle := parts[0]

	// Clean up the handle - only allow alphanumeric, underscores, and hyphens
	cleanHandle := strings.Builder{}
	for _, r := range handle {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' || r == '-' {
			cleanHandle.WriteRune(r)
		}
	}

	return strings.ToLower(cleanHandle.String())
}

// extractDomainFromURL extracts the full domain from a URL.
// For example:
// - "https://daverupert.com" -> "daverupert.com"
// - "https://www.example.com" -> "www.example.com"
// - "https://blog.jane.dev" -> "blog.jane.dev"
func extractDomainFromURL(siteURL string) string {
	parsed, err := url.Parse(siteURL)
	if err != nil {
		return ""
	}

	host := parsed.Hostname()
	if host == "" {
		return ""
	}

	return strings.ToLower(host)
}

// registerFromPosts registers handles from internal posts matching from_posts sources.
func (p *MentionsPlugin) registerFromPosts(m *lifecycle.Manager, config models.MentionsConfig, handleMap map[string]*mentionEntry) {
	if len(config.FromPosts) == 0 {
		return
	}

	posts := m.Posts()

	for _, source := range config.FromPosts {
		if source.Filter == "" {
			continue
		}

		// Parse and apply the filter
		f, err := filter.Parse(source.Filter)
		if err != nil {
			log.Printf("mentions: error parsing from_posts filter %q: %v", source.Filter, err)
			continue
		}

		matchedPosts := f.MatchAll(posts)

		for _, post := range matchedPosts {
			p.registerPostAsHandle(post, source, handleMap)
		}
	}
}

// registerPostAsHandle registers a single post as a handle source.
func (p *MentionsPlugin) registerPostAsHandle(post *models.Post, source models.MentionPostSource, handleMap map[string]*mentionEntry) {
	// Get the handle from the specified field or fall back to slug
	handle := p.getHandleFromPost(post, source.HandleField)
	if handle == "" {
		return
	}

	handle = strings.ToLower(handle)

	// Get the title for the entry
	title := post.PlainTitle()

	// Build metadata from the post's frontmatter fields
	metadata := p.buildPostMetadata(post, source)

	// Create the entry with the post's Href as the URL
	entry := &mentionEntry{
		Handle:   handle,
		SiteURL:  post.Href,
		Title:    title,
		Internal: true,
		Metadata: metadata,
	}

	// Register the handle (first entry wins)
	if _, exists := handleMap[handle]; !exists {
		handleMap[handle] = entry
	}

	// Register aliases if configured
	if source.AliasesField != "" {
		aliases := p.getAliasesFromPost(post, source.AliasesField)
		for _, alias := range aliases {
			p.registerAlias(alias, entry, handleMap)
		}
	}
}

// getHandleFromPost extracts a handle from a post's frontmatter field or falls back to slug.
func (p *MentionsPlugin) getHandleFromPost(post *models.Post, fieldName string) string {
	// If no field specified, use slug
	if fieldName == "" {
		return post.Slug
	}

	// Try to get the field from Extra
	if post.Extra == nil {
		return post.Slug
	}

	value, ok := post.Extra[fieldName]
	if !ok {
		return post.Slug
	}

	// Handle string value
	if str, ok := value.(string); ok && str != "" {
		return str
	}

	// Fall back to slug
	return post.Slug
}

// getAliasesFromPost extracts aliases from a post's frontmatter field.
func (p *MentionsPlugin) getAliasesFromPost(post *models.Post, fieldName string) []string {
	if fieldName == "" || post.Extra == nil {
		return nil
	}

	value, ok := post.Extra[fieldName]
	if !ok {
		return nil
	}

	// Handle []string
	if aliases, ok := value.([]string); ok {
		return aliases
	}

	// Handle []interface{} (common from YAML/JSON parsing)
	if aliases, ok := value.([]interface{}); ok {
		result := make([]string, 0, len(aliases))
		for _, alias := range aliases {
			if str, ok := alias.(string); ok && str != "" {
				result = append(result, str)
			}
		}
		return result
	}

	return nil
}

// buildPostMetadata creates MentionMetadata from a post's frontmatter fields.
// This allows internal posts (contacts, characters, etc.) to provide hovercard
// data directly from frontmatter instead of fetching from HTTP.
func (p *MentionsPlugin) buildPostMetadata(post *models.Post, source models.MentionPostSource) *models.MentionMetadata {
	metadata := &models.MentionMetadata{
		URL:         post.Href,
		LastFetched: time.Now(),
	}

	// Name from post title
	metadata.Name = post.PlainTitle()
	if metadata.Name == "" {
		metadata.Name = post.Slug
	}

	// Bio from post description
	if post.Description != nil {
		metadata.Bio = *post.Description
	}

	// Avatar: try configured field first, then default fields
	metadata.Avatar = p.getAvatarFromPost(post, source.AvatarField)
	normalizeMentionMetadata(metadata)

	return metadata
}

// avatarDefaultFields are the frontmatter fields checked for avatar URLs
// when no explicit avatar_field is configured. Checked in order.
var avatarDefaultFields = []string{"avatar", "image", "icon"}

// getAvatarFromPost extracts an avatar URL from a post's frontmatter.
// If avatarField is specified, it checks only that field.
// Otherwise, it checks "avatar", "image", and "icon" fields in order.
func (p *MentionsPlugin) getAvatarFromPost(post *models.Post, avatarField string) string {
	if post.Extra == nil {
		return ""
	}

	// If a specific field is configured, use only that
	if avatarField != "" {
		return getStringField(post.Extra, avatarField)
	}

	// Check default fields in order
	for _, field := range avatarDefaultFields {
		if value := getStringField(post.Extra, field); value != "" {
			return value
		}
	}

	return ""
}

// getStringField extracts a string value from a map, returning "" if not found or not a string.
func getStringField(extra map[string]interface{}, key string) string {
	if extra == nil {
		return ""
	}
	value, ok := extra[key]
	if !ok {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

// getMentionsConfig retrieves mentions configuration from the manager config.
func getMentionsConfig(config *lifecycle.Config) models.MentionsConfig {
	if config.Extra == nil {
		return models.NewMentionsConfig()
	}

	// First check for typed MentionsConfig in Extra
	if mentions, ok := config.Extra["mentions"]; ok {
		if mc, ok := mentions.(models.MentionsConfig); ok {
			return mc
		}
	}

	return models.NewMentionsConfig()
}

// Ensure MentionsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*MentionsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*MentionsPlugin)(nil)
	_ lifecycle.TransformPlugin = (*MentionsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*MentionsPlugin)(nil)
)
