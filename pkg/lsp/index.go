package lsp

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"gopkg.in/yaml.v3"
)

// Index maintains an index of all markdown posts in the workspace.
type Index struct {
	logger *log.Logger

	// posts maps slug to PostInfo for quick lookup
	posts map[string]*PostInfo
	mu    sync.RWMutex

	// uriToSlug maps file URI to slug for reverse lookup
	uriToSlug map[string]string

	// mentions maps handle/alias to MentionInfo for quick lookup
	mentions  map[string]*MentionInfo
	mentionMu sync.RWMutex
}

// PostInfo contains indexed information about a post.
type PostInfo struct {
	// URI is the file URI (file:///path/to/file.md)
	URI string

	// Path is the file system path
	Path string

	// Slug is the URL-safe identifier
	Slug string

	// Title is the post title (from frontmatter or filename)
	Title string

	// Description is the post description/excerpt
	Description string

	// Metadata is the parsed frontmatter for filter evaluation
	Metadata map[string]interface{}

	// Aliases are alternative slugs that resolve to this post
	Aliases []string

	// Wikilinks contains all wikilinks found in the post
	Wikilinks []WikilinkInfo
}

// MentionInfo contains indexed information about a mention.
// Mentions can be external (from blogroll) or internal (from posts matching a filter).
type MentionInfo struct {
	// Handle is the primary handle (e.g., "daverupert")
	Handle string

	// Aliases are alternative handles that resolve to this mention
	Aliases []string

	// Title is the display name (e.g., "Dave Rupert")
	Title string

	// Description is a short description
	Description string

	// --- External (blogroll) fields ---

	// SiteURL is the website URL (external mentions only)
	SiteURL string

	// FeedURL is the RSS/Atom feed URL (external mentions only)
	FeedURL string

	// --- Internal (from_posts) fields ---

	// IsInternal indicates this mention comes from a post, not blogroll
	IsInternal bool

	// Slug is the post slug (internal mentions only)
	Slug string

	// Path is the file path (internal mentions only)
	Path string
}

// WikilinkInfo contains information about a wikilink in a post.
type WikilinkInfo struct {
	// Target is the slug being linked to
	Target string

	// DisplayText is the optional display text
	DisplayText string

	// Line is the 0-based line number
	Line int

	// StartChar is the 0-based character position of [[
	StartChar int

	// EndChar is the 0-based character position after ]]
	EndChar int
}

// NewIndex creates a new post index.
func NewIndex(logger *log.Logger) *Index {
	return &Index{
		logger:    logger,
		posts:     make(map[string]*PostInfo),
		uriToSlug: make(map[string]string),
		mentions:  make(map[string]*MentionInfo),
	}
}

// Build indexes all markdown files in the workspace.
func (idx *Index) Build(rootPath string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Clear existing index
	idx.posts = make(map[string]*PostInfo)
	idx.uriToSlug = make(map[string]string)

	// Index blogroll mentions from config
	idx.indexBlogrollMentions(rootPath)

	// Walk the directory tree
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil // Skip files we can't access
		}

		// Skip directories and non-markdown files
		if info.IsDir() {
			// Skip common non-content directories
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == "output" || name == ".markata" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}

		// Index this file
		if indexErr := idx.indexFile(path); indexErr != nil {
			idx.logger.Printf("Failed to index %s: %v", path, indexErr)
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Index mentions from posts matching from_posts config
	idx.indexFromPostsMentions(rootPath)

	return nil
}

// indexFile indexes a single markdown file.
func (idx *Index) indexFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return idx.indexContent(path, string(content))
}

// indexContent indexes content for a given path.
func (idx *Index) indexContent(path, content string) error {
	// Parse frontmatter
	metadata, body, err := plugins.ParseFrontmatter(content)
	if err != nil {
		// Still index the file, just without metadata
		metadata = make(map[string]interface{})
		body = content
	}

	// Generate slug
	slug := generateSlug(path, metadata)

	// Get title
	title := plugins.GetString(metadata, "title")
	if title == "" {
		// Derive from filename
		base := filepath.Base(path)
		title = strings.TrimSuffix(base, filepath.Ext(base))
		title = strings.ReplaceAll(title, "-", " ")
		title = strings.ReplaceAll(title, "_", " ")
	}

	// Get description
	description := plugins.GetString(metadata, "description")
	if description == "" {
		// Extract first paragraph as description
		description = extractExcerpt(body, 200)
	}

	// Find wikilinks
	wikilinks := findWikilinks(body)

	// Extract aliases from frontmatter
	aliases := extractAliases(metadata)

	// Create post info
	uri := pathToURI(path)
	info := &PostInfo{
		URI:         uri,
		Path:        path,
		Slug:        slug,
		Title:       title,
		Description: description,
		Metadata:    metadata,
		Aliases:     aliases,
		Wikilinks:   wikilinks,
	}

	// Store in index by slug
	idx.posts[slug] = info
	idx.uriToSlug[uri] = slug

	// Store alias entries (only if not already taken by a slug)
	for _, alias := range aliases {
		normalizedAlias := strings.ToLower(alias)
		if _, exists := idx.posts[normalizedAlias]; !exists {
			idx.posts[normalizedAlias] = info
		}
	}

	return nil
}

// Update updates the index for a single file.
func (idx *Index) Update(uri, content string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	path := uriToPath(uri)

	// Remove old entry if exists (including aliases)
	if oldSlug, ok := idx.uriToSlug[uri]; ok {
		if oldInfo, exists := idx.posts[oldSlug]; exists {
			// Remove alias entries that point to this post
			for _, alias := range oldInfo.Aliases {
				normalizedAlias := strings.ToLower(alias)
				if existingInfo, aliasExists := idx.posts[normalizedAlias]; aliasExists && existingInfo == oldInfo {
					delete(idx.posts, normalizedAlias)
				}
			}
		}
		delete(idx.posts, oldSlug)
		delete(idx.uriToSlug, uri)
	}

	return idx.indexContent(path, content)
}

// Remove removes a file from the index.
func (idx *Index) Remove(uri string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if slug, ok := idx.uriToSlug[uri]; ok {
		if info, exists := idx.posts[slug]; exists {
			// Remove alias entries that point to this post
			for _, alias := range info.Aliases {
				normalizedAlias := strings.ToLower(alias)
				if existingInfo, aliasExists := idx.posts[normalizedAlias]; aliasExists && existingInfo == info {
					delete(idx.posts, normalizedAlias)
				}
			}
		}
		delete(idx.posts, slug)
		delete(idx.uriToSlug, uri)
	}
}

// GetBySlug returns post info for a slug.
func (idx *Index) GetBySlug(slug string) *PostInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	// Try exact match first
	if info, ok := idx.posts[slug]; ok {
		return info
	}

	// Try case-insensitive match
	normalizedSlug := normalizeSlug(slug)
	for postSlug, info := range idx.posts {
		if normalizeSlug(postSlug) == normalizedSlug {
			return info
		}
	}

	return nil
}

// GetByURI returns post info for a URI.
func (idx *Index) GetByURI(uri string) *PostInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if slug, ok := idx.uriToSlug[uri]; ok {
		return idx.posts[slug]
	}
	return nil
}

// AllPosts returns all indexed posts (unique by slug, excluding alias entries).
func (idx *Index) AllPosts() []*PostInfo {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	// Use a map to deduplicate (aliases point to same PostInfo)
	seen := make(map[*PostInfo]bool)
	posts := make([]*PostInfo, 0, len(idx.posts))
	for _, info := range idx.posts {
		if !seen[info] {
			seen[info] = true
			posts = append(posts, info)
		}
	}
	return posts
}

// Match type constants for PostSearchResult.
const (
	MatchTypeSlug  = "slug"
	MatchTypeAlias = "alias"
	MatchTypeTitle = "title"
)

// PostSearchResult contains a post and how it was matched.
type PostSearchResult struct {
	Post       *PostInfo
	MatchedBy  string // "slug", "alias", or "title"
	MatchedKey string // The actual slug/alias that matched
}

// SearchPosts returns posts matching a prefix.
func (idx *Index) SearchPosts(prefix string) []*PostInfo {
	results := idx.SearchPostsWithMatch(prefix)
	posts := make([]*PostInfo, 0, len(results))
	for _, r := range results {
		posts = append(posts, r.Post)
	}
	return posts
}

// SearchPostsWithMatch returns posts matching a prefix with match information.
// Results are deduplicated - each post appears once with the best match type
// (slug matches take precedence over alias matches, which take precedence over title matches).
func (idx *Index) SearchPostsWithMatch(prefix string) []PostSearchResult {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	prefix = strings.ToLower(prefix)

	// Track best match for each post
	bestMatch := make(map[*PostInfo]*PostSearchResult)

	for key, info := range idx.posts {
		keyLower := strings.ToLower(key)

		// Check if this key matches the prefix
		if strings.HasPrefix(keyLower, prefix) {
			// Determine match type
			var matchType string
			if key == info.Slug {
				matchType = MatchTypeSlug
			} else {
				matchType = MatchTypeAlias
			}

			// Check if this is a better match than existing
			existing := bestMatch[info]
			if existing == nil || matchTypePriority(matchType) < matchTypePriority(existing.MatchedBy) {
				bestMatch[info] = &PostSearchResult{
					Post:       info,
					MatchedBy:  matchType,
					MatchedKey: key,
				}
			}
		} else if strings.Contains(strings.ToLower(info.Title), prefix) {
			// Title match - only add if no better match exists
			if bestMatch[info] == nil {
				bestMatch[info] = &PostSearchResult{
					Post:       info,
					MatchedBy:  MatchTypeTitle,
					MatchedKey: info.Slug,
				}
			}
		}
	}

	// Convert map to slice
	results := make([]PostSearchResult, 0, len(bestMatch))
	for _, r := range bestMatch {
		results = append(results, *r)
	}

	return results
}

// matchTypePriority returns priority for match types (lower = better).
func matchTypePriority(matchType string) int {
	switch matchType {
	case MatchTypeSlug:
		return 0
	case MatchTypeAlias:
		return 1
	case MatchTypeTitle:
		return 2
	default:
		return 3
	}
}

// AllMentions returns all indexed mentions (unique by handle).
func (idx *Index) AllMentions() []*MentionInfo {
	idx.mentionMu.RLock()
	defer idx.mentionMu.RUnlock()

	// Deduplicate by handle (aliases point to same MentionInfo)
	seen := make(map[string]bool)
	mentions := make([]*MentionInfo, 0)
	for _, info := range idx.mentions {
		if !seen[info.Handle] {
			seen[info.Handle] = true
			mentions = append(mentions, info)
		}
	}
	return mentions
}

// SearchMentions returns mentions matching a prefix.
func (idx *Index) SearchMentions(prefix string) []*MentionInfo {
	idx.mentionMu.RLock()
	defer idx.mentionMu.RUnlock()

	prefix = strings.ToLower(prefix)

	// Deduplicate results by handle
	seen := make(map[string]bool)
	var results []*MentionInfo

	for key, info := range idx.mentions {
		if seen[info.Handle] {
			continue
		}
		// Match against handle, aliases, or title
		if strings.HasPrefix(strings.ToLower(key), prefix) ||
			strings.Contains(strings.ToLower(info.Title), prefix) {
			seen[info.Handle] = true
			results = append(results, info)
		}
	}

	return results
}

// GetByHandle returns a mention by handle or alias (case-insensitive).
func (idx *Index) GetByHandle(handle string) *MentionInfo {
	idx.mentionMu.RLock()
	defer idx.mentionMu.RUnlock()

	handle = strings.ToLower(handle)
	return idx.mentions[handle]
}

// indexBlogrollMentions reads the markata config and indexes blogroll feeds as mentions.
func (idx *Index) indexBlogrollMentions(rootPath string) {
	idx.mentionMu.Lock()
	defer idx.mentionMu.Unlock()

	idx.mentions = make(map[string]*MentionInfo)

	// Try to find and parse config file
	configPaths := []string{
		filepath.Join(rootPath, "markata-go.toml"),
		filepath.Join(rootPath, "markata.toml"),
		filepath.Join(rootPath, "markata-go.yaml"),
		filepath.Join(rootPath, "markata.yaml"),
	}

	for _, configPath := range configPaths {
		content, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		var feeds []blogrollFeedConfig
		if strings.HasSuffix(configPath, ".toml") {
			feeds = idx.parseBlogrollFromTOML(content)
		} else {
			feeds = idx.parseBlogrollFromYAML(content)
		}

		for _, feed := range feeds {
			idx.indexFeedAsMention(&feed)
		}

		if len(feeds) > 0 {
			idx.logger.Printf("Indexed %d blogroll mentions from %s", len(idx.mentions), configPath)
		}
		return // Only use first config found
	}
}

// blogrollFeedConfig is a minimal struct for parsing feed configs.
type blogrollFeedConfig struct {
	URL         string   `toml:"url" yaml:"url"`
	Title       string   `toml:"title" yaml:"title"`
	Description string   `toml:"description" yaml:"description"`
	SiteURL     string   `toml:"site_url" yaml:"site_url"`
	Handle      string   `toml:"handle" yaml:"handle"`
	Aliases     []string `toml:"aliases" yaml:"aliases"`
	Active      *bool    `toml:"active" yaml:"active"`
}

// parseBlogrollFromTOML extracts feed configs from TOML content.
func (idx *Index) parseBlogrollFromTOML(content []byte) []blogrollFeedConfig {
	// Try markata-go.blogroll first (new format)
	var configNew struct {
		MarktaGo struct {
			Blogroll struct {
				Enabled bool                 `toml:"enabled"`
				Feeds   []blogrollFeedConfig `toml:"feeds"`
			} `toml:"blogroll"`
		} `toml:"markata-go"`
	}

	if err := toml.Unmarshal(content, &configNew); err == nil {
		if configNew.MarktaGo.Blogroll.Enabled && len(configNew.MarktaGo.Blogroll.Feeds) > 0 {
			return configNew.MarktaGo.Blogroll.Feeds
		}
	}

	// Fall back to blogroll (old format)
	var configOld struct {
		Blogroll struct {
			Enabled bool                 `toml:"enabled"`
			Feeds   []blogrollFeedConfig `toml:"feeds"`
		} `toml:"blogroll"`
	}

	if err := toml.Unmarshal(content, &configOld); err != nil {
		idx.logger.Printf("Failed to parse TOML config: %v", err)
		return nil
	}

	if !configOld.Blogroll.Enabled {
		return nil
	}

	return configOld.Blogroll.Feeds
}

// parseBlogrollFromYAML extracts feed configs from YAML content.
func (idx *Index) parseBlogrollFromYAML(content []byte) []blogrollFeedConfig {
	var config struct {
		Blogroll struct {
			Enabled bool                 `yaml:"enabled"`
			Feeds   []blogrollFeedConfig `yaml:"feeds"`
		} `yaml:"blogroll"`
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		idx.logger.Printf("Failed to parse YAML config: %v", err)
		return nil
	}

	if !config.Blogroll.Enabled {
		return nil
	}

	return config.Blogroll.Feeds
}

// indexFeedAsMention creates a MentionInfo from a feed config.
func (idx *Index) indexFeedAsMention(feed *blogrollFeedConfig) {
	// Skip inactive feeds
	if feed.Active != nil && !*feed.Active {
		return
	}

	// Determine handle
	handle := feed.Handle
	if handle == "" {
		handle = extractHandleFromFeedURL(feed.SiteURL)
		if handle == "" {
			handle = extractHandleFromFeedURL(feed.URL)
		}
	}
	if handle == "" {
		return
	}

	handle = strings.ToLower(handle)

	// Create mention info
	info := &MentionInfo{
		Handle:      handle,
		Aliases:     feed.Aliases,
		Title:       feed.Title,
		SiteURL:     feed.SiteURL,
		FeedURL:     feed.URL,
		Description: feed.Description,
	}

	// Index by handle
	if _, exists := idx.mentions[handle]; !exists {
		idx.mentions[handle] = info
	}

	// Index by aliases
	for _, alias := range feed.Aliases {
		alias = strings.ToLower(alias)
		if alias != "" {
			if _, exists := idx.mentions[alias]; !exists {
				idx.mentions[alias] = info
			}
		}
	}

	// Auto-register domain alias
	if feed.SiteURL != "" {
		domain := extractDomainAlias(feed.SiteURL)
		if domain != "" && domain != handle {
			if _, exists := idx.mentions[domain]; !exists {
				idx.mentions[domain] = info
			}
		}
	}
}

// indexFromPostsMentions reads the mentions.from_posts config and indexes
// matching posts as mentions for LSP completion and hover.
func (idx *Index) indexFromPostsMentions(rootPath string) {
	// Try to find and parse config file
	configPaths := []string{
		filepath.Join(rootPath, "markata-go.toml"),
		filepath.Join(rootPath, "markata.toml"),
		filepath.Join(rootPath, "markata-go.yaml"),
		filepath.Join(rootPath, "markata.yaml"),
	}

	for _, configPath := range configPaths {
		content, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		var sources []mentionPostSource
		if strings.HasSuffix(configPath, ".toml") {
			sources = idx.parseMentionsFromPostsTOML(content)
		} else {
			sources = idx.parseMentionsFromPostsYAML(content)
		}

		if len(sources) == 0 {
			// Fall back to default: template == 'contact' posts
			sources = []mentionPostSource{
				{
					Filter:      "template == 'contact'",
					HandleField: "handle",
				},
			}
		}

		// Get all unique posts (filter out alias entries)
		uniquePosts := idx.getUniquePosts()

		for _, source := range sources {
			if source.Filter == "" {
				continue
			}

			// Parse the filter expression
			f, err := filter.Parse(source.Filter)
			if err != nil {
				idx.logger.Printf("mentions from_posts: error parsing filter %q: %v", source.Filter, err)
				continue
			}

			// Evaluate filter against each post
			matchCount := 0
			for _, postInfo := range uniquePosts {
				// Convert PostInfo to a minimal models.Post for filter evaluation
				post := idx.postInfoToPost(postInfo)

				match, err := f.Match(post)
				if err != nil {
					continue
				}

				if match {
					idx.indexPostAsMention(postInfo, source)
					matchCount++
				}
			}

			if matchCount > 0 {
				idx.logger.Printf("Indexed %d mentions from posts matching filter %q", matchCount, source.Filter)
			}
		}

		return // Only use first config found
	}
}

// mentionPostSource mirrors models.MentionPostSource for config parsing.
type mentionPostSource struct {
	Filter       string `toml:"filter" yaml:"filter"`
	HandleField  string `toml:"handle_field" yaml:"handle_field"`
	AliasesField string `toml:"aliases_field" yaml:"aliases_field"`
	AvatarField  string `toml:"avatar_field" yaml:"avatar_field"`
}

// parseMentionsFromPostsTOML extracts from_posts config from TOML content.
func (idx *Index) parseMentionsFromPostsTOML(content []byte) []mentionPostSource {
	// Try markata-go.mentions first (new format)
	var configNew struct {
		MarktaGo struct {
			Mentions struct {
				FromPosts []mentionPostSource `toml:"from_posts"`
			} `toml:"mentions"`
		} `toml:"markata-go"`
	}

	if err := toml.Unmarshal(content, &configNew); err == nil {
		if len(configNew.MarktaGo.Mentions.FromPosts) > 0 {
			return configNew.MarktaGo.Mentions.FromPosts
		}
	}

	// Fall back to mentions (old format)
	var configOld struct {
		Mentions struct {
			FromPosts []mentionPostSource `toml:"from_posts"`
		} `toml:"mentions"`
	}

	if err := toml.Unmarshal(content, &configOld); err != nil {
		return nil
	}

	return configOld.Mentions.FromPosts
}

// parseMentionsFromPostsYAML extracts from_posts config from YAML content.
func (idx *Index) parseMentionsFromPostsYAML(content []byte) []mentionPostSource {
	var config struct {
		Mentions struct {
			FromPosts []mentionPostSource `yaml:"from_posts"`
		} `yaml:"mentions"`
	}

	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil
	}

	return config.Mentions.FromPosts
}

// getUniquePosts returns all unique posts (excluding alias entries).
func (idx *Index) getUniquePosts() []*PostInfo {
	seen := make(map[string]bool)
	var posts []*PostInfo

	for _, post := range idx.posts {
		if !seen[post.Path] {
			seen[post.Path] = true
			posts = append(posts, post)
		}
	}

	return posts
}

// postInfoToPost converts a PostInfo to a minimal models.Post for filter evaluation.
func (idx *Index) postInfoToPost(info *PostInfo) *models.Post {
	post := &models.Post{
		Path: info.Path,
		Slug: info.Slug,
		Href: "/" + info.Slug + "/",
	}

	// Set title if present
	if info.Title != "" {
		post.Title = &info.Title
	}

	// Set description if present
	if info.Description != "" {
		post.Description = &info.Description
	}

	// Extract tags from metadata
	if info.Metadata != nil {
		if tags, ok := info.Metadata["tags"]; ok {
			post.Tags = extractTagsFromMetadata(tags)
		}

		// Copy other metadata to Extra for filter access
		post.Extra = make(map[string]interface{})
		for k, v := range info.Metadata {
			post.Extra[k] = v
		}
	}

	return post
}

// extractTagsFromMetadata converts various tag formats to []string.
func extractTagsFromMetadata(tags interface{}) []string {
	switch t := tags.(type) {
	case []string:
		return t
	case []interface{}:
		result := make([]string, 0, len(t))
		for _, v := range t {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case string:
		// Single tag as string
		return []string{t}
	default:
		return nil
	}
}

// indexPostAsMention indexes a post as a mention.
func (idx *Index) indexPostAsMention(postInfo *PostInfo, source mentionPostSource) {
	idx.mentionMu.Lock()
	defer idx.mentionMu.Unlock()

	// Determine handle from configured field or fall back to slug
	handle := idx.getHandleFromPostInfo(postInfo, source.HandleField)
	if handle == "" {
		return
	}

	handle = strings.ToLower(handle)

	// Create mention info for internal post
	info := &MentionInfo{
		Handle:      handle,
		Title:       postInfo.Title,
		Description: postInfo.Description,
		IsInternal:  true,
		Slug:        postInfo.Slug,
		Path:        postInfo.Path,
	}

	// Index by handle (first entry wins)
	if _, exists := idx.mentions[handle]; !exists {
		idx.mentions[handle] = info
	}

	// Index aliases if configured
	if source.AliasesField != "" && postInfo.Metadata != nil {
		if aliasesRaw, ok := postInfo.Metadata[source.AliasesField]; ok {
			aliases := extractTagsFromMetadata(aliasesRaw) // Reuse tag extraction logic
			for _, alias := range aliases {
				alias = strings.ToLower(alias)
				if alias != "" && alias != handle {
					if _, exists := idx.mentions[alias]; !exists {
						idx.mentions[alias] = info
					}
				}
			}
			info.Aliases = aliases
		}
	}
}

// getHandleFromPostInfo extracts a handle from a post's metadata.
func (idx *Index) getHandleFromPostInfo(postInfo *PostInfo, fieldName string) string {
	// If no field specified, use slug
	if fieldName == "" {
		return postInfo.Slug
	}

	// Check metadata for the field
	if postInfo.Metadata == nil {
		return postInfo.Slug
	}

	// Special case: if field is "slug", return slug directly
	if fieldName == "slug" {
		return postInfo.Slug
	}

	// Look up the field in metadata
	if val, ok := postInfo.Metadata[fieldName]; ok {
		if str, ok := val.(string); ok && str != "" {
			return str
		}
	}

	// Fall back to slug
	return postInfo.Slug
}

// extractHandleFromFeedURL extracts a handle from a URL's domain.
func extractHandleFromFeedURL(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	// Simple extraction - find the domain part
	urlStr = strings.TrimPrefix(urlStr, "https://")
	urlStr = strings.TrimPrefix(urlStr, "http://")
	urlStr = strings.TrimPrefix(urlStr, "www.")
	urlStr = strings.TrimPrefix(urlStr, "blog.")

	// Get first part before /
	if idx := strings.Index(urlStr, "/"); idx > 0 {
		urlStr = urlStr[:idx]
	}

	// Get first part before .
	parts := strings.Split(urlStr, ".")
	if len(parts) > 0 {
		return strings.ToLower(parts[0])
	}

	return ""
}

// extractDomainAlias extracts the full domain as an alias.
func extractDomainAlias(urlStr string) string {
	if urlStr == "" {
		return ""
	}

	urlStr = strings.TrimPrefix(urlStr, "https://")
	urlStr = strings.TrimPrefix(urlStr, "http://")

	// Get domain before /
	if idx := strings.Index(urlStr, "/"); idx > 0 {
		urlStr = urlStr[:idx]
	}

	return strings.ToLower(urlStr)
}

// generateSlug generates a slug from path and metadata.
func generateSlug(path string, metadata map[string]interface{}) string {
	// Check for slug in metadata
	if slug := plugins.GetString(metadata, "slug"); slug != "" {
		return slug
	}

	// Generate from filename
	base := filepath.Base(path)

	// Handle index.md special case
	if strings.EqualFold(base, "index.md") {
		dir := filepath.Dir(path)
		dir = filepath.Clean(dir)
		if dir == "." {
			return "" // Root index.md becomes homepage
		}
		slug := strings.ToLower(dir)
		slug = strings.ReplaceAll(slug, string(filepath.Separator), "/")
		slug = strings.TrimPrefix(slug, "./")
		return slug
	}

	// Use filename without known extension (consistent with models.Post)
	basename := models.StripKnownExtension(base)
	return models.Slugify(basename)
}

// extractAliases extracts aliases from frontmatter metadata.
// Aliases can be specified as a list of strings in the "aliases" field.
func extractAliases(metadata map[string]interface{}) []string {
	aliasesRaw, ok := metadata["aliases"]
	if !ok {
		return nil
	}

	var aliases []string

	switch v := aliasesRaw.(type) {
	case []interface{}:
		for _, alias := range v {
			if aliasStr, ok := alias.(string); ok && aliasStr != "" {
				aliases = append(aliases, aliasStr)
			}
		}
	case []string:
		for _, alias := range v {
			if alias != "" {
				aliases = append(aliases, alias)
			}
		}
	}

	return aliases
}

// wikilinkRegex matches [[slug]] and [[slug|display text]] patterns.
var wikilinkRegex = regexp.MustCompile(`\[\[([^\]|]+)(?:\|([^\]]+))?\]\]`)

// findWikilinks finds all wikilinks in content.
func findWikilinks(content string) []WikilinkInfo {
	var results []WikilinkInfo
	lines := strings.Split(content, "\n")

	for lineNum, line := range lines {
		matches := wikilinkRegex.FindAllStringSubmatchIndex(line, -1)
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}

			// match[0:2] is the full match position
			// match[2:4] is group 1 (slug)
			// match[4:6] is group 2 (display text) if present
			fullMatch := line[match[0]:match[1]]
			groups := wikilinkRegex.FindStringSubmatch(fullMatch)

			target := strings.TrimSpace(groups[1])
			displayText := ""
			if len(groups) > 2 && groups[2] != "" {
				displayText = strings.TrimSpace(groups[2])
			}

			results = append(results, WikilinkInfo{
				Target:      target,
				DisplayText: displayText,
				Line:        lineNum,
				StartChar:   match[0],
				EndChar:     match[1],
			})
		}
	}

	return results
}

// inlineAttributesRegex matches inline attributes {.class #id key="value"}
var inlineAttributesRegex = regexp.MustCompile(`\{[.#][^}]+\}`)

// extractExcerpt extracts the first paragraph as an excerpt.
func extractExcerpt(content string, maxLen int) string {
	// Remove leading whitespace and headers
	content = strings.TrimSpace(content)

	// Skip leading headers
	lines := strings.Split(content, "\n")
	textLines := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			if len(textLines) > 0 {
				break // End of first paragraph
			}
			continue
		}
		textLines = append(textLines, line)
	}

	if len(textLines) == 0 {
		return ""
	}

	excerpt := strings.Join(textLines, " ")

	// Strip inline attributes like {.text-gray-500}
	excerpt = inlineAttributesRegex.ReplaceAllString(excerpt, "")
	excerpt = strings.TrimSpace(excerpt)

	if len(excerpt) > maxLen {
		excerpt = excerpt[:maxLen-3] + "..."
	}

	return excerpt
}

// normalizeSlug normalizes a slug for comparison.
func normalizeSlug(slug string) string {
	return models.Slugify(slug)
}

// pathToURI converts a file system path to a file URI.
func pathToURI(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	// Normalize separators for URI
	absPath = filepath.ToSlash(absPath)
	return "file://" + absPath
}

// uriToPath converts a file URI to a file system path.
func uriToPath(uri string) string {
	path := strings.TrimPrefix(uri, "file://")
	// Handle Windows paths
	if len(path) > 2 && path[0] == '/' && path[2] == ':' {
		path = path[1:]
	}
	return filepath.FromSlash(path)
}
