// Package buildcache provides incremental build caching for markata-go.
//
// The cache tracks input hashes (content + frontmatter + template) for each post
// and allows skipping rebuild when inputs haven't changed. It also tracks global
// hashes (templates, config) to invalidate the entire cache when needed.
//
// # Cache File Structure
//
// The cache is stored in .markata/build-cache.json with the structure:
//
//	{
//	  "version": 1,
//	  "config_hash": "abc123",
//	  "templates_hash": "def456",
//	  "posts": {
//	    "path/to/post.md": {
//	      "input_hash": "xyz789",
//	      "output_path": "output/post/index.html",
//	      "template": "post.html"
//	    }
//	  },
//	  "graph": {
//	    "dependencies": {"path/to/post.md": ["linked-slug"]},
//	    "path_to_slug": {"path/to/post.md": "post-slug"}
//	  }
//	}
package buildcache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

// CacheVersion is incremented when the cache format changes.
const CacheVersion = 3

// DefaultCacheDir is the directory for cache files.
const DefaultCacheDir = ".markata"

// CacheFileName is the name of the build cache file.
const CacheFileName = "build-cache.json"

// Cache manages incremental build state.
type Cache struct {
	mu sync.RWMutex

	// Version of the cache format
	Version int `json:"version"`

	// ConfigHash is the hash of the config file contents
	ConfigHash string `json:"config_hash"`

	// TemplatesHash is the combined hash of all template files
	TemplatesHash string `json:"templates_hash"`
	// NavPreviewHash invalidates pages when shared navigation preview data changes.
	NavPreviewHash string `json:"nav_preview_hash,omitempty"`
	// FontpackHash invalidates rendered HTML when shared font CSS or preloads change.
	FontpackHash string `json:"fontpack_hash,omitempty"`

	// TemplatesFingerprint is a cheap filesystem fingerprint used to skip
	// recomputing TemplatesHash when the template tree is unchanged.
	TemplatesFingerprint string `json:"templates_fingerprint,omitempty"`

	// AssetsHash is the combined hash of all static assets (JS/CSS)
	AssetsHash string `json:"assets_hash,omitempty"`

	// Posts maps source path to cached post metadata
	Posts map[string]*PostCache `json:"posts"`

	// Feeds maps feed slug to cached feed metadata
	Feeds map[string]*FeedCache `json:"feeds,omitempty"`

	// TagsListingHash caches the tags listing output state
	TagsListingHash string `json:"tags_listing_hash,omitempty"`

	// GardenHash caches the garden view output state
	GardenHash string `json:"garden_hash,omitempty"`

	// FeedsListingHash caches the /feeds listing output state.
	FeedsListingHash string `json:"feeds_listing_hash,omitempty"`

	// ImageLibraryHash caches the generated image inventory and library output state.
	ImageLibraryHash string `json:"image_library_hash,omitempty"`

	// ImageLibraryContentRoot identifies the site root for the recorded image
	// library files. It scopes ownership when a cache is reused by another site.
	ImageLibraryContentRoot string `json:"image_library_content_root,omitempty"`

	// ImageLibraryOutputRoot identifies the output root for the recorded image
	// library files. It scopes cleanup when a cache is reused with another output.
	ImageLibraryOutputRoot string `json:"image_library_output_root,omitempty"`

	// ImageLibraryOutputHashes records image-library output paths relative to
	// ImageLibraryOutputRoot and their last generated content hashes. Content
	// hashes prevent cleanup from deleting a file that a user later changed.
	ImageLibraryOutputHashes map[string]string `json:"image_library_output_hashes,omitempty"`

	// ImageLibraryMedia stores content fingerprints for local image-library
	// media. The map key is a normalized source path. Size, ModTime, and an
	// optional ChangeTime are cheap invalidation state; ContentHash is reused
	// while that state matches.
	ImageLibraryMedia map[string]ImageLibraryMediaFingerprint `json:"image_library_media,omitempty"`

	// TailwindManifestHash caches the generated Tailwind token manifest state.
	TailwindManifestHash string `json:"tailwind_manifest_hash,omitempty"`

	// PagefindCorpusHash caches the searchable HTML corpus state.
	PagefindCorpusHash string `json:"pagefind_corpus_hash,omitempty"`

	// Graph tracks dependencies between posts for transitive invalidation
	Graph *DependencyGraph `json:"graph,omitempty"`

	// path to the cache file (not serialized)
	path string

	// stalePosts preserves output ownership while a global cache invalidation is
	// being processed. The entries are intentionally not serialized: a
	// successful build replaces them with current entries before the cache is
	// saved.
	stalePosts map[string]*PostCache

	// dirty tracks whether cache needs saving
	dirty bool

	// skippedCount tracks how many posts were skipped this build (lock-free)
	skippedCount atomic.Int64

	// rebuiltCount tracks how many posts were rebuilt this build (lock-free)
	rebuiltCount atomic.Int64

	// changedSlugs tracks slugs that changed this build (for dependency invalidation)
	changedSlugs map[string]bool

	// changedFeedSlugs tracks slugs with feed-relevant changes this build
	changedFeedSlugs map[string]bool

	// tagsDirty marks whether tag outputs should rebuild this build
	tagsDirty bool

	// gardenDirty marks whether garden outputs should rebuild this build
	gardenDirty bool

	// GlobFiles caches the list of discovered files from glob stage
	GlobFiles []string `json:"glob_files,omitempty"`

	// GlobPatternHash detects when glob patterns change
	GlobPatternHash string `json:"glob_pattern_hash,omitempty"`

	// In-memory caches to avoid repeated per-post disk reads during hot builds.
	// These are populated lazily on first access (backfilled from disk on cache
	// hit) and updated by CacheFullHTML/CacheArticleHTML/CachePostData on writes.
	// Using sync.Map for lock-free concurrent reads from worker pools.

	// fullHTMLMemory maps FullHTMLPath -> HTML string
	fullHTMLMemory sync.Map
	// articleHTMLMemory maps ArticleHTMLPath -> HTML string
	articleHTMLMemory sync.Map
	// postDataMemory maps post cache file path -> *CachedPostData
	postDataMemory sync.Map
	// encryptedHTMLMemory maps EncryptedHTMLPath -> HTML string
	encryptedHTMLMemory sync.Map
	// mermaidSVGMemory maps MermaidSVGPath -> rendered SVG string
	mermaidSVGMemory sync.Map
	// mermaidSVGUsed tracks Mermaid SVG cache file paths referenced in the current build.
	mermaidSVGUsed sync.Map
}

// PostCache stores cached metadata for a single post.
type PostCache struct {
	// InputHash is the hash of content + frontmatter + template
	InputHash string `json:"input_hash"`

	// OutputPath is the primary output file path
	OutputPath string `json:"output_path"`

	// Template is the template name used for rendering
	Template string `json:"template"`

	// OutputHash is the hash of the rendered output (optional, for verification)
	OutputHash string `json:"output_hash,omitempty"`

	// ContentHash is the hash of just the markdown content (for render caching)
	ContentHash string `json:"content_hash,omitempty"`

	// ArticleHTMLPath is the path to the cached rendered HTML file
	ArticleHTMLPath string `json:"article_html_path,omitempty"`

	// FullHTMLPath is the path to the cached full page HTML file
	FullHTMLPath string `json:"full_html_path,omitempty"`

	// HeadingHighlightRevision records the last successfully cached page revision.
	HeadingHighlightRevision int `json:"heading_highlight_revision,omitempty"`

	// ModTime is the file modification time (Unix nanoseconds)
	ModTime int64 `json:"mod_time,omitempty"`

	// Slug is the post's slug for dependency tracking
	Slug string `json:"slug,omitempty"`

	// LinkHrefsHash is a hash of the post's article HTML used for link extraction caching
	LinkHrefsHash string `json:"link_hrefs_hash,omitempty"`

	// LinkHrefs caches extracted href values for the post
	LinkHrefs []string `json:"link_hrefs,omitempty"`

	// FeedMembershipHash is the hash of the sorted slugs of feed co-members.
	// When feed membership changes (posts added/removed from a tag), this hash
	// changes and the post is rebuilt with the updated sidebar.
	FeedMembershipHash string `json:"feed_membership_hash,omitempty"`
	// LocalPreviewHash tracks preview data used by links and tags on this page.
	LocalPreviewHash string `json:"local_preview_hash,omitempty"`

	// LinkAvatarsHash is a hash of the ArticleHTML input used for link_avatars caching.
	// When ArticleHTML changes, this hash changes and the post is re-processed.
	LinkAvatarsHash string `json:"link_avatars_hash,omitempty"`

	// LinkAvatarsHTML is the cached ArticleHTML output after link_avatars processing.
	LinkAvatarsHTML string `json:"link_avatars_html,omitempty"`

	// EmbedsHash is a hash of the post Content input used for embeds transform caching.
	// When Content changes, this hash changes and the post is re-processed.
	EmbedsHash string `json:"embeds_hash,omitempty"`

	// EmbedsContent is the cached Content output after embeds transform processing.
	EmbedsContent string `json:"embeds_content,omitempty"`

	// GlossaryHash is a hash of the ArticleHTML input + glossary terms state used for caching.
	GlossaryHash string `json:"glossary_hash,omitempty"`

	// GlossaryHTML is the cached ArticleHTML output after glossary processing.
	GlossaryHTML string `json:"glossary_html,omitempty"`

	// EncryptedHash is a hash of inputs for encrypted HTML caching.
	EncryptedHash string `json:"encrypted_hash,omitempty"`

	// EncryptedHTMLPath is the path to cached encrypted HTML wrapper.
	EncryptedHTMLPath string `json:"encrypted_html_path,omitempty"`

	// FeedItemHash is a hash of fields that affect feed output for this post.
	FeedItemHash string `json:"feed_item_hash,omitempty"`

	// TagIndexHash is a hash of fields that affect tag listings for this post.
	TagIndexHash string `json:"tag_index_hash,omitempty"`

	// GardenHash is a hash of fields that affect garden view output for this post.
	GardenHash string `json:"garden_hash,omitempty"`

	// TailwindHTMLHash is a hash of the final HTML used for Tailwind token extraction.
	TailwindHTMLHash string `json:"tailwind_html_hash,omitempty"`

	// TailwindTokens stores the extracted Tailwind-relevant tokens for the post.
	TailwindTokens string `json:"tailwind_tokens,omitempty"`

	// MentionsHash is a hash of the Content input + mention resolution state.
	// When either the raw content or resolved mention metadata changes, this hash changes.
	MentionsHash string `json:"mentions_hash,omitempty"`

	// MentionsContent is the cached Content output after mentions transform processing.
	MentionsContent string `json:"mentions_content,omitempty"`
}

// ImageLibraryMediaFingerprint stores the cheap filesystem state and content
// hash for one local image-library media file. ChangeTime is an operating
// system change-time signal when the filesystem exposes one.
type ImageLibraryMediaFingerprint struct {
	Size            int64  `json:"size"`
	ModTime         int64  `json:"mod_time"`
	ChangeTime      int64  `json:"change_time,omitempty"`
	ChangeTimeKnown bool   `json:"change_time_known,omitempty"`
	ContentHash     string `json:"content_hash"`
}

// FeedCache stores cached metadata for a single feed.
type FeedCache struct {
	// Hash is the hash of the feed's content (post slugs, config)
	Hash string `json:"hash"`
}

// New creates a new empty cache.
func New(cacheDir string) *Cache {
	if cacheDir == "" {
		cacheDir = DefaultCacheDir
	}
	return &Cache{
		Version:           CacheVersion,
		Posts:             make(map[string]*PostCache),
		Feeds:             make(map[string]*FeedCache),
		ImageLibraryMedia: make(map[string]ImageLibraryMediaFingerprint),
		Graph:             NewDependencyGraph(),
		path:              filepath.Join(cacheDir, CacheFileName),
		stalePosts:        make(map[string]*PostCache),
		changedSlugs:      make(map[string]bool),
		changedFeedSlugs:  make(map[string]bool),
	}
}

// Load reads the cache from disk. Returns a new empty cache if file doesn't exist.
func Load(cacheDir string) (*Cache, error) {
	cache := New(cacheDir)

	data, err := os.ReadFile(cache.path)
	if err != nil {
		if os.IsNotExist(err) {
			return cache, nil // Empty cache is fine
		}
		return nil, fmt.Errorf("reading cache file: %w", err)
	}

	if err := json.Unmarshal(data, cache); err != nil {
		// Corrupt cache, start fresh
		return New(cacheDir), nil
	}

	// Version mismatch, start fresh
	if cache.Version != CacheVersion {
		return New(cacheDir), nil
	}

	// Ensure Graph is initialized (might be nil from old cache versions)
	if cache.Graph == nil {
		cache.Graph = NewDependencyGraph()
	} else {
		// Rebuild the reverse index (Dependents) from persisted Dependencies
		cache.Graph.RebuildReverse()
	}

	// Ensure Feeds is initialized
	if cache.Feeds == nil {
		cache.Feeds = make(map[string]*FeedCache)
	}
	if cache.ImageLibraryMedia == nil {
		cache.ImageLibraryMedia = make(map[string]ImageLibraryMediaFingerprint)
	}

	// Ensure changedSlugs is initialized
	cache.changedSlugs = make(map[string]bool)
	cache.changedFeedSlugs = make(map[string]bool)

	return cache, nil
}

// preservePostOwnershipLocked keeps the previous post metadata available to
// cleanup code after a global invalidation clears Posts. The caller must hold
// c.mu.
func (c *Cache) preservePostOwnershipLocked() {
	if len(c.Posts) == 0 {
		return
	}
	if c.stalePosts == nil {
		c.stalePosts = make(map[string]*PostCache, len(c.Posts))
	}
	for path, post := range c.Posts {
		if post != nil {
			c.stalePosts[path] = post
		}
	}
}

// Save writes the cache to disk.
func (c *Cache) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.dirty {
		return nil
	}

	// Ensure directory exists
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating cache directory: %w", err)
	}

	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("marshaling cache: %w", err)
	}

	//nolint:gosec // G306: Cache file needs 0644 for sharing across builds
	if err := os.WriteFile(c.path, data, 0o644); err != nil {
		return fmt.Errorf("writing cache file: %w", err)
	}

	c.dirty = false
	return nil
}

// SetConfigHash updates the config hash and invalidates if changed.
func (c *Cache) SetConfigHash(hash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.ConfigHash == hash {
		return false // No change
	}

	// Config changed - invalidate all posts
	c.ConfigHash = hash
	c.preservePostOwnershipLocked()
	c.Posts = make(map[string]*PostCache)
	c.Feeds = make(map[string]*FeedCache)
	c.FeedsListingHash = ""
	c.ImageLibraryHash = ""
	c.TailwindManifestHash = ""
	c.PagefindCorpusHash = ""
	c.dirty = true
	return true
}

// SetTemplatesHash updates the templates hash and invalidates if changed.
func (c *Cache) SetTemplatesHash(hash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.TemplatesHash == hash {
		return false // No change
	}

	// Templates changed - invalidate all posts
	c.TemplatesHash = hash
	c.preservePostOwnershipLocked()
	c.Posts = make(map[string]*PostCache)
	c.Feeds = make(map[string]*FeedCache)
	c.ImageLibraryHash = ""
	c.TailwindManifestHash = ""
	c.PagefindCorpusHash = ""
	c.dirty = true
	return true
}

// SetNavPreviewHash invalidates rendered pages when shared nav metadata changes.
func (c *Cache) SetNavPreviewHash(hash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.NavPreviewHash == hash {
		return false
	}
	c.NavPreviewHash = hash
	c.preservePostOwnershipLocked()
	c.Posts = make(map[string]*PostCache)
	c.Feeds = make(map[string]*FeedCache)
	c.dirty = true
	return true
}

// SetFontpackHash invalidates cached rendered pages when shared font CSS or
// preload URLs change. The fontpack plugin runs after content loading, so this
// separate context hash avoids changing source input hashes and warm load work.
func (c *Cache) SetFontpackHash(hash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.FontpackHash == hash {
		return false
	}
	c.FontpackHash = hash
	c.preservePostOwnershipLocked()
	c.Posts = make(map[string]*PostCache)
	c.Feeds = make(map[string]*FeedCache)
	c.FeedsListingHash = ""
	c.dirty = true
	return true
}

// GetTemplatesHash returns the cached template content hash.
func (c *Cache) GetTemplatesHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TemplatesHash
}

// GetTemplatesFingerprint returns the cached template tree fingerprint.
func (c *Cache) GetTemplatesFingerprint() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TemplatesFingerprint
}

// SetTemplatesFingerprint updates the template tree fingerprint and marks the
// cache dirty only when the fingerprint actually changes.
func (c *Cache) SetTemplatesFingerprint(hash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.TemplatesFingerprint == hash {
		return false
	}

	c.TemplatesFingerprint = hash
	c.dirty = true
	return true
}

// SetAssetsHash updates the assets hash and invalidates if changed.
// This ensures pages are rebuilt when JS/CSS files change so they reference new hashed filenames.
func (c *Cache) SetAssetsHash(hash string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.AssetsHash == hash {
		return false // No change
	}

	// Assets changed - invalidate all posts
	c.AssetsHash = hash
	c.preservePostOwnershipLocked()
	c.Posts = make(map[string]*PostCache)
	c.Feeds = make(map[string]*FeedCache)
	c.ImageLibraryHash = ""
	c.TailwindManifestHash = ""
	c.PagefindCorpusHash = ""
	c.dirty = true
	return true
}

// ShouldRebuild checks if a post needs rebuilding based on input hash.
// Returns true if the post should be rebuilt (hash mismatch or not in cache).
// Also returns true if any post this one depends on has changed.
func (c *Cache) ShouldRebuild(sourcePath, inputHash, template string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.shouldRebuildLocked(sourcePath, inputHash, template)
}

// shouldRebuildLocked shares metadata gates with ReasonForRebuild and batch checks.
// Caller must hold at least c.mu.RLock().
func (c *Cache) shouldRebuildLocked(sourcePath, inputHash, template string) bool {
	return c.rebuildReasonLocked(sourcePath, inputHash, template) != RebuildReasonNone
}

func (c *Cache) rebuildReasonLocked(sourcePath, inputHash, template string) RebuildReason {
	cached, ok := c.Posts[sourcePath]
	if !ok || cached == nil {
		return RebuildReasonMissingEntry
	}

	if cached.InputHash != inputHash {
		return RebuildReasonInputChanged
	}
	if cached.Template != template {
		return RebuildReasonTemplateChanged
	}

	return RebuildReasonNone
}

// ShouldRebuildBatch checks multiple posts against the cache in a single lock acquisition.
// Returns a map of sourcePath -> true for posts that need rebuilding.
// This is more efficient than calling ShouldRebuild in a loop.
func (c *Cache) ShouldRebuildBatch(posts []struct {
	Path, InputHash, Template string
}) map[string]bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]bool, len(posts)/10) // expect ~10% changed
	for _, p := range posts {
		if c.shouldRebuildLocked(p.Path, p.InputHash, p.Template) {
			result[p.Path] = true
		}
	}
	return result
}

// ShouldRebuildWithSlug checks if a post needs rebuilding based on input hash
// and also checks if any of the posts it depends on have changed this build.
// The slug is used for dependency tracking.
func (c *Cache) ShouldRebuildWithSlug(sourcePath, _, inputHash, template string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return true // Not in cache
	}

	// Check if hash matches and template is the same
	if cached.InputHash != inputHash || cached.Template != template {
		return true // Changed
	}

	// Check if any dependency has changed this build
	if deps := c.Graph.GetDependencies(sourcePath); len(deps) > 0 {
		for _, dep := range deps {
			if c.changedSlugs[dep] {
				return true // A dependency changed
			}
		}
	}

	return false
}

// MarkRebuilt records that a post was rebuilt with the given hash.
// The slug is used for dependency invalidation tracking.
func (c *Cache) MarkRebuilt(sourcePath, inputHash, outputPath, template string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.stalePosts, sourcePath)

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.InputHash = inputHash
		cached.OutputPath = outputPath
		cached.Template = template
	} else {
		c.Posts[sourcePath] = &PostCache{
			InputHash:  inputHash,
			OutputPath: outputPath,
			Template:   template,
		}
	}
	c.dirty = true
	c.rebuiltCount.Add(1)
}

// MarkRebuiltWithSlug records that a post was rebuilt with the given hash.
// Also records that this slug changed, for dependency invalidation.
func (c *Cache) MarkRebuiltWithSlug(sourcePath, slug, inputHash, outputPath, template string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.stalePosts, sourcePath)

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.InputHash = inputHash
		cached.OutputPath = outputPath
		cached.Template = template
		cached.Slug = slug
	} else {
		c.Posts[sourcePath] = &PostCache{
			InputHash:  inputHash,
			OutputPath: outputPath,
			Template:   template,
			Slug:       slug,
		}
	}
	c.dirty = true
	c.rebuiltCount.Add(1)

	// Track that this slug changed (for dependency invalidation)
	if slug != "" {
		c.changedSlugs[slug] = true
	}
}

// MarkInputProcessed records the input state after a post's enabled outputs
// were processed successfully, without counting an HTML write. Posts without
// HTML output can still be dependency targets, so their inputs must not appear
// changed on every build. The method clears obsolete HTML ownership and reports
// whether the cached state changed.
func (c *Cache) MarkInputProcessed(sourcePath, slug, inputHash, template string) bool {
	if sourcePath == "" || inputHash == "" {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		if cached.InputHash == inputHash && cached.Template == template && cached.OutputPath == "" && (slug == "" || cached.Slug == slug) {
			return false
		}
		delete(c.stalePosts, sourcePath)
		cached.InputHash = inputHash
		cached.OutputPath = ""
		cached.Template = template
		if slug != "" {
			cached.Slug = slug
		}
	} else {
		delete(c.stalePosts, sourcePath)
		c.Posts[sourcePath] = &PostCache{
			InputHash: inputHash,
			Slug:      slug,
			Template:  template,
		}
	}
	c.dirty = true
	return true
}

// UpdatePostSemanticHashes updates cached per-post hashes for feed/tags/garden.
// Returns which hashes changed compared to cache.
func (c *Cache) UpdatePostSemanticHashes(sourcePath, feedHash, tagHash, gardenHash string) (feedChanged, tagChanged, gardenChanged bool) {
	if sourcePath == "" {
		return false, false, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	cacheEntry, ok := c.Posts[sourcePath]
	if !ok {
		cacheEntry = &PostCache{}
		c.Posts[sourcePath] = cacheEntry
	}

	if cacheEntry.FeedItemHash != feedHash {
		feedChanged = true
		cacheEntry.FeedItemHash = feedHash
	}
	if cacheEntry.TagIndexHash != tagHash {
		tagChanged = true
		cacheEntry.TagIndexHash = tagHash
	}
	if cacheEntry.GardenHash != gardenHash {
		gardenChanged = true
		cacheEntry.GardenHash = gardenHash
	}
	if tagChanged {
		c.tagsDirty = true
	}
	if gardenChanged {
		c.gardenDirty = true
	}

	if feedChanged || tagChanged || gardenChanged {
		c.dirty = true
	}

	return feedChanged, tagChanged, gardenChanged
}

// GetPostSemanticHashes returns the cached feed/tag/garden hashes for a post.
func (c *Cache) GetPostSemanticHashes(sourcePath string) (feedHash, tagHash, gardenHash string) {
	feedHash, tagHash, gardenHash, _ = c.GetPostSemanticHashBaseline(sourcePath)
	return feedHash, tagHash, gardenHash
}

// GetPostSemanticHashBaseline returns prior hashes and whether they can be
// compared canonically. Empty/partial hashes are valid for new entries without
// an input identity, but missing derived metadata on an established entry is
// unknown, not a semantic change. All fields are read under the same lock.
func (c *Cache) GetPostSemanticHashBaseline(sourcePath string) (feedHash, tagHash, gardenHash string, available bool) {
	if sourcePath == "" {
		return "", "", "", true
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if entry := c.Posts[sourcePath]; entry != nil {
		available = entry.InputHash == "" ||
			(entry.FeedItemHash != "" && entry.TagIndexHash != "" && entry.GardenHash != "")
		return entry.FeedItemHash, entry.TagIndexHash, entry.GardenHash, available
	}
	return "", "", "", true
}

// TagsDirty reports whether any tag-relevant fields changed this build.
func (c *Cache) TagsDirty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tagsDirty
}

// GardenDirty reports whether any garden-relevant fields changed this build.
func (c *Cache) GardenDirty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.gardenDirty
}

// SetTagsListingHash stores the tags listing hash in the cache.
func (c *Cache) SetTagsListingHash(hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.TagsListingHash == hash {
		return
	}
	c.TagsListingHash = hash
	c.dirty = true
}

// GetTagsListingHash returns the cached tags listing hash.
func (c *Cache) GetTagsListingHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TagsListingHash
}

// SetGardenHash stores the garden view hash in the cache.
func (c *Cache) SetGardenHash(hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.GardenHash == hash {
		return
	}
	c.GardenHash = hash
	c.dirty = true
}

// GetGardenHash returns the cached garden view hash.
func (c *Cache) GetGardenHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.GardenHash
}

// SetFeedsListingHash stores the feeds listing hash in the cache.
func (c *Cache) SetFeedsListingHash(hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.FeedsListingHash == hash {
		return
	}
	c.FeedsListingHash = hash
	c.dirty = true
}

// GetFeedsListingHash returns the cached feeds listing hash.
func (c *Cache) GetFeedsListingHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.FeedsListingHash
}

// SetImageLibraryHash stores the image-library output hash in the cache.
func (c *Cache) SetImageLibraryHash(hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ImageLibraryHash == hash {
		return
	}
	c.ImageLibraryHash = hash
	c.dirty = true
}

// SetImageLibraryOutputs records image-library output hashes relative to the
// supplied output root and scoped to the supplied site root.
func (c *Cache) SetImageLibraryOutputs(contentRoot, outputRoot string, outputHashes map[string]string) {
	siteRoot, err := filepath.Abs(contentRoot)
	if err != nil {
		return
	}
	root, err := filepath.Abs(outputRoot)
	if err != nil {
		return
	}
	siteRoot = filepath.Clean(siteRoot)
	root = filepath.Clean(root)
	normalized := make(map[string]string, len(outputHashes))
	for path, hash := range outputHashes {
		absolute, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		relative, err := filepath.Rel(root, filepath.Clean(absolute))
		if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) || hash == "" {
			continue
		}
		normalized[filepath.Clean(relative)] = hash
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ImageLibraryContentRoot == siteRoot && c.ImageLibraryOutputRoot == root && equalStringMaps(c.ImageLibraryOutputHashes, normalized) {
		return
	}
	c.ImageLibraryContentRoot = siteRoot
	c.ImageLibraryOutputRoot = root
	c.ImageLibraryOutputHashes = normalized
	c.dirty = true
}

func equalStringMaps(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}

// GetImageLibraryOutputs returns the recorded site root, output root, and
// relative path hashes.
func (c *Cache) GetImageLibraryOutputs() (contentRoot, outputRoot string, outputs map[string]string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	outputs = make(map[string]string, len(c.ImageLibraryOutputHashes))
	for path, hash := range c.ImageLibraryOutputHashes {
		outputs[path] = hash
	}
	return c.ImageLibraryContentRoot, c.ImageLibraryOutputRoot, outputs
}

// OwnsImageLibraryOutput reports whether the cache records the exact path under
// the supplied site and output roots. It does not verify the current file
// contents.
func (c *Cache) OwnsImageLibraryOutput(contentRoot, outputRoot, path string) bool {
	siteRoot, err := filepath.Abs(contentRoot)
	if err != nil {
		return false
	}
	root, err := filepath.Abs(outputRoot)
	if err != nil {
		return false
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(absolute))
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if filepath.Clean(c.ImageLibraryContentRoot) != filepath.Clean(siteRoot) || filepath.Clean(c.ImageLibraryOutputRoot) != filepath.Clean(root) {
		return false
	}
	_, ok := c.ImageLibraryOutputHashes[filepath.Clean(relative)]
	return ok
}

// ImageLibraryOutputHash returns the expected hash for a recorded exact path.
func (c *Cache) ImageLibraryOutputHash(contentRoot, outputRoot, path string) (string, bool) {
	siteRoot, err := filepath.Abs(contentRoot)
	if err != nil {
		return "", false
	}
	root, err := filepath.Abs(outputRoot)
	if err != nil {
		return "", false
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(absolute))
	if err != nil || relative == "." || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return "", false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()
	if filepath.Clean(c.ImageLibraryContentRoot) != filepath.Clean(siteRoot) || filepath.Clean(c.ImageLibraryOutputRoot) != filepath.Clean(root) {
		return "", false
	}
	hash, ok := c.ImageLibraryOutputHashes[filepath.Clean(relative)]
	return hash, ok
}

// ClearImageLibraryOutputs forgets the image-library output ownership records.
func (c *Cache) ClearImageLibraryOutputs() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ImageLibraryContentRoot == "" && c.ImageLibraryOutputRoot == "" && len(c.ImageLibraryOutputHashes) == 0 {
		return
	}
	c.ImageLibraryContentRoot = ""
	c.ImageLibraryOutputRoot = ""
	c.ImageLibraryOutputHashes = nil
	c.dirty = true
}

// GetImageLibraryMedia returns a copy of the cached local-media fingerprints.
func (c *Cache) GetImageLibraryMedia() map[string]ImageLibraryMediaFingerprint {
	c.mu.RLock()
	defer c.mu.RUnlock()

	media := make(map[string]ImageLibraryMediaFingerprint, len(c.ImageLibraryMedia))
	for path, fingerprint := range c.ImageLibraryMedia {
		media[path] = fingerprint
	}
	return media
}

// SetImageLibraryMedia replaces the cached local-media fingerprints. It
// reports whether the persisted cache state changed.
func (c *Cache) SetImageLibraryMedia(media map[string]ImageLibraryMediaFingerprint) bool {
	normalized := make(map[string]ImageLibraryMediaFingerprint, len(media))
	for path, fingerprint := range media {
		if path == "" || fingerprint.ContentHash == "" {
			continue
		}
		normalized[filepath.ToSlash(filepath.Clean(path))] = fingerprint
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if equalImageLibraryMedia(c.ImageLibraryMedia, normalized) {
		return false
	}
	c.ImageLibraryMedia = normalized
	c.dirty = true
	return true
}

func equalImageLibraryMedia(left, right map[string]ImageLibraryMediaFingerprint) bool {
	if len(left) != len(right) {
		return false
	}
	for path, fingerprint := range left {
		if right[path] != fingerprint {
			return false
		}
	}
	return true
}

// GetImageLibraryHash returns the cached image-library output hash.
func (c *Cache) GetImageLibraryHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ImageLibraryHash
}

// MarkSkipped records that a post was skipped (already up to date).
// Uses atomic increment -- no mutex needed.
func (c *Cache) MarkSkipped() {
	c.skippedCount.Add(1)
}

// IsFileUnchanged checks if a file's ModTime matches the cached value.
// Returns true if the file has not changed since last build.
// Returns false if file is not in cache or ModTime differs.
func (c *Cache) IsFileUnchanged(sourcePath string, modTime int64) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return false
	}
	return cached.ModTime == modTime && cached.ModTime != 0
}

// GetCachedPost returns cached post metadata, including ownership metadata
// preserved during a global invalidation. Returns nil if the source is unknown.
func (c *Cache) GetCachedPost(sourcePath string) *PostCache {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if post := c.Posts[sourcePath]; post != nil {
		return post
	}
	return c.stalePosts[sourcePath]
}

// UpdateModTime updates the ModTime for a post.
func (c *Cache) UpdateModTime(sourcePath string, modTime int64, slug string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.ModTime = modTime
		cached.Slug = slug
	} else {
		c.Posts[sourcePath] = &PostCache{
			ModTime: modTime,
			Slug:    slug,
		}
	}
	c.dirty = true
}

// GetCachedLinkHrefs returns cached hrefs for a post if the article hash matches.
// Returns nil if no matching cache exists.
func (c *Cache) GetCachedLinkHrefs(sourcePath, articleHash string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return nil
	}
	if cached.LinkHrefsHash == "" || cached.LinkHrefsHash != articleHash {
		return nil
	}
	if len(cached.LinkHrefs) == 0 {
		return []string{}
	}

	hrefs := make([]string, len(cached.LinkHrefs))
	copy(hrefs, cached.LinkHrefs)
	return hrefs
}

// CacheLinkHrefs stores extracted hrefs for a post keyed by article hash.
func (c *Cache) CacheLinkHrefs(sourcePath, articleHash string, hrefs []string) {
	if sourcePath == "" || articleHash == "" {
		return
	}

	copyHrefs := make([]string, len(hrefs))
	copy(copyHrefs, hrefs)

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.LinkHrefsHash = articleHash
		cached.LinkHrefs = copyHrefs
	} else {
		c.Posts[sourcePath] = &PostCache{
			LinkHrefsHash: articleHash,
			LinkHrefs:     copyHrefs,
		}
	}
	c.dirty = true
}

// GetCachedLinkAvatarsHTML returns cached link_avatars output if the article hash matches.
// Returns the cached HTML and true if valid, empty string and false otherwise.
func (c *Cache) GetCachedLinkAvatarsHTML(sourcePath, articleHash string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return "", false
	}
	if cached.LinkAvatarsHash == "" || cached.LinkAvatarsHash != articleHash {
		return "", false
	}
	return cached.LinkAvatarsHTML, true
}

// CacheLinkAvatarsHTML stores the link_avatars processed HTML keyed by article hash.
func (c *Cache) CacheLinkAvatarsHTML(sourcePath, articleHash, html string) {
	if sourcePath == "" || articleHash == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.LinkAvatarsHash = articleHash
		cached.LinkAvatarsHTML = html
	} else {
		c.Posts[sourcePath] = &PostCache{
			LinkAvatarsHash: articleHash,
			LinkAvatarsHTML: html,
		}
	}
	c.dirty = true
}

// GetCachedEmbedsContent returns cached embeds transform output if the content hash matches.
// Returns the cached content and true if valid, empty string and false otherwise.
func (c *Cache) GetCachedEmbedsContent(sourcePath, contentHash string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return "", false
	}
	if cached.EmbedsHash == "" || cached.EmbedsHash != contentHash {
		return "", false
	}
	return cached.EmbedsContent, true
}

// CacheEmbedsContent stores the embeds transform output keyed by content hash.
func (c *Cache) CacheEmbedsContent(sourcePath, contentHash, content string) {
	if sourcePath == "" || contentHash == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.EmbedsHash = contentHash
		cached.EmbedsContent = content
	} else {
		c.Posts[sourcePath] = &PostCache{
			EmbedsHash:    contentHash,
			EmbedsContent: content,
		}
	}
	c.dirty = true
}

// GetCachedGlossaryHTML returns cached glossary output if the combined hash matches.
// Returns the cached HTML and true if valid, empty string and false otherwise.
func (c *Cache) GetCachedGlossaryHTML(sourcePath, combinedHash string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return "", false
	}
	if cached.GlossaryHash == "" || cached.GlossaryHash != combinedHash {
		return "", false
	}
	return cached.GlossaryHTML, true
}

// CacheGlossaryHTML stores the glossary processed HTML keyed by combined hash.
func (c *Cache) CacheGlossaryHTML(sourcePath, combinedHash, html string) {
	if sourcePath == "" || combinedHash == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.GlossaryHash = combinedHash
		cached.GlossaryHTML = html
	} else {
		c.Posts[sourcePath] = &PostCache{
			GlossaryHash: combinedHash,
			GlossaryHTML: html,
		}
	}
	c.dirty = true
}

// GetCachedMentionsContent returns cached mentions transform output if the hash matches.
// Returns the cached content and true if valid, empty string and false otherwise.
func (c *Cache) GetCachedMentionsContent(sourcePath, contentHash string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return "", false
	}
	if cached.MentionsHash == "" || cached.MentionsHash != contentHash {
		return "", false
	}
	return cached.MentionsContent, true
}

// CacheMentionsContent stores the mentions transform output keyed by content hash.
func (c *Cache) CacheMentionsContent(sourcePath, contentHash, content string) {
	if sourcePath == "" || contentHash == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.MentionsHash = contentHash
		cached.MentionsContent = content
	} else {
		c.Posts[sourcePath] = &PostCache{
			MentionsHash:    contentHash,
			MentionsContent: content,
		}
	}
	c.dirty = true
}

// GetCachedTailwindTokens returns cached Tailwind tokens if the HTML hash matches.
func (c *Cache) GetCachedTailwindTokens(sourcePath, htmlHash string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok {
		return "", false
	}
	if cached.TailwindHTMLHash == "" || cached.TailwindHTMLHash != htmlHash {
		return "", false
	}
	return cached.TailwindTokens, true
}

// GetStoredTailwindTokens returns cached Tailwind tokens without validating the
// HTML hash. This is safe for unchanged posts whose final HTML is known to be
// identical to the prior build even if the current build has not eagerly
// restored the cached full-page HTML into memory.
func (c *Cache) GetStoredTailwindTokens(sourcePath string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, ok := c.Posts[sourcePath]
	if !ok || cached.TailwindTokens == "" {
		return "", false
	}
	return cached.TailwindTokens, true
}

// CacheTailwindTokens stores extracted Tailwind tokens keyed by final HTML hash.
func (c *Cache) CacheTailwindTokens(sourcePath, htmlHash, tokens string) {
	if sourcePath == "" || htmlHash == "" {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.TailwindHTMLHash = htmlHash
		cached.TailwindTokens = tokens
	} else {
		c.Posts[sourcePath] = &PostCache{
			TailwindHTMLHash: htmlHash,
			TailwindTokens:   tokens,
		}
	}
	c.dirty = true
}

// GetTailwindManifestHash returns the cached Tailwind manifest hash.
func (c *Cache) GetTailwindManifestHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.TailwindManifestHash
}

// SetTailwindManifestHash stores the Tailwind manifest hash.
func (c *Cache) SetTailwindManifestHash(hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.TailwindManifestHash == hash {
		return
	}
	c.TailwindManifestHash = hash
	c.dirty = true
}

// GetPagefindCorpusHash returns the cached Pagefind corpus hash.
func (c *Cache) GetPagefindCorpusHash() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.PagefindCorpusHash
}

// SetPagefindCorpusHash stores the Pagefind corpus hash.
func (c *Cache) SetPagefindCorpusHash(hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.PagefindCorpusHash == hash {
		return
	}
	c.PagefindCorpusHash = hash
	c.dirty = true
}

// SetFeedMembershipHash stores the feed membership hash for a post.
func (c *Cache) SetFeedMembershipHash(sourcePath, hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cached, ok := c.Posts[sourcePath]; ok {
		cached.FeedMembershipHash = hash
	}
	c.dirty = true
}

// GetFeedMembershipHash returns the cached feed membership hash for a post.
func (c *Cache) GetFeedMembershipHash(sourcePath string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if cached, ok := c.Posts[sourcePath]; ok {
		return cached.FeedMembershipHash
	}
	return ""
}

// GetHeadingHighlightRevision returns zero when no live post entry exists.
func (c *Cache) GetHeadingHighlightRevision(sourcePath string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if cached := c.Posts[sourcePath]; cached != nil {
		return cached.HeadingHighlightRevision
	}
	return 0
}

// SetHeadingHighlightRevision updates only existing live metadata, not stale
// ownership entries retained during invalidation.
func (c *Cache) SetHeadingHighlightRevision(sourcePath string, revision int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cached := c.Posts[sourcePath]; cached != nil && cached.HeadingHighlightRevision != revision {
		cached.HeadingHighlightRevision = revision
		c.dirty = true
	}
}

// SetLocalPreviewHash stores the preview fingerprint for a rendered page.
func (c *Cache) SetLocalPreviewHash(sourcePath, hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cached, ok := c.Posts[sourcePath]; ok {
		cached.LocalPreviewHash = hash
		c.dirty = true
	}
}

// GetLocalPreviewHash returns the preview fingerprint for a rendered page.
func (c *Cache) GetLocalPreviewHash(sourcePath string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if cached, ok := c.Posts[sourcePath]; ok {
		return cached.LocalPreviewHash
	}
	return ""
}

// ComputeFeedMembershipHash computes a hash of the sorted slugs of feed co-members.
// Returns empty string if the slug list is empty.
func ComputeFeedMembershipHash(slugs []string) string {
	if len(slugs) == 0 {
		return ""
	}
	sorted := make([]string, len(slugs))
	copy(sorted, slugs)
	sort.Strings(sorted)
	return HashContent(strings.Join(sorted, "\x00"))
}

// Stats returns build statistics.
func (c *Cache) Stats() (skipped, rebuilt int) {
	return int(c.skippedCount.Load()), int(c.rebuiltCount.Load())
}

// CacheStats holds build statistics.
type CacheStats struct {
	Skipped int
	Rebuilt int
}

// GetStats returns build statistics as a struct.
func (c *Cache) GetStats() CacheStats {
	return CacheStats{
		Skipped: int(c.skippedCount.Load()),
		Rebuilt: int(c.rebuiltCount.Load()),
	}
}

// ResetStats resets the build statistics for a new build.
func (c *Cache) ResetStats() {
	c.skippedCount.Store(0)
	c.rebuiltCount.Store(0)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changedSlugs = make(map[string]bool)
	c.changedFeedSlugs = make(map[string]bool)
	c.tagsDirty = false
	c.gardenDirty = false
}

// RemoveStale removes cache entries for posts that no longer exist.
// Returns the number of entries removed.
func (c *Cache) RemoveStale(currentPaths map[string]bool) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	removed := 0
	for path := range c.Posts {
		if !currentPaths[path] {
			delete(c.Posts, path)
			c.Graph.RemoveSource(path)
			removed++
			c.dirty = true
		}
	}
	for path := range c.stalePosts {
		if !currentPaths[path] {
			c.Graph.RemoveSource(path)
			removed++
		}
	}
	if len(c.stalePosts) > 0 {
		c.stalePosts = make(map[string]*PostCache)
		c.dirty = true
	}
	return removed
}

// HashContent computes a SHA256 hash of the given content.
func HashContent(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])
}

// HashFile computes a SHA256 hash of a file's contents.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashDirectory computes a combined hash of all files in a directory.
// Files are processed in sorted order for deterministic output.
func HashDirectory(dir string, extensions []string) (string, error) {
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if len(extMap) == 0 || extMap[ext] {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil // Directory doesn't exist, empty hash
		}
		return "", err
	}

	// Sort for deterministic ordering
	sort.Strings(files)

	h := sha256.New()
	for _, path := range files {
		// Include relative path in hash (so renames are detected)
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path // fallback to absolute path
		}
		h.Write([]byte(relPath))

		// Include file contents
		content, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		h.Write(content)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashFS computes a deterministic content hash for matching files in an fs.FS.
// File paths are included so adding, removing, or renaming a template changes
// the hash as well as editing its contents.
func HashFS(fsys fs.FS, extensions []string) (string, error) {
	if fsys == nil {
		return "", nil
	}
	extMap := make(map[string]bool, len(extensions))
	for _, ext := range extensions {
		extMap[ext] = true
	}

	var files []string
	if err := fs.WalkDir(fsys, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		if len(extMap) == 0 || extMap[filepath.Ext(path)] {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return "", err
	}
	sort.Strings(files)

	h := sha256.New()
	for _, path := range files {
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return "", err
		}
		h.Write([]byte(path))
		h.Write(content)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashDirectoryState computes a cheap fingerprint for a directory tree.
// It includes each matching file's relative path, size, and modtime so callers
// can skip a full content hash when the tree is unchanged.
func HashDirectoryState(dir string, extensions []string) (string, error) {
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}

	type fileState struct {
		path    string
		size    int64
		modTime int64
	}

	states := make([]fileState, 0, 64)
	err := filepath.WalkDir(dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if len(extMap) > 0 && !extMap[ext] {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path
		}

		states = append(states, fileState{
			path:    relPath,
			size:    info.Size(),
			modTime: info.ModTime().UnixNano(),
		})
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	sort.Slice(states, func(i, j int) bool { return states[i].path < states[j].path })

	h := sha256.New()
	for _, state := range states {
		h.Write([]byte(state.path))
		h.Write([]byte{0})
		h.Write([]byte(strconv.FormatInt(state.size, 10)))
		h.Write([]byte{0})
		h.Write([]byte(strconv.FormatInt(state.modTime, 10)))
		h.Write([]byte{0})
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashAssetMap computes a combined hash from a map of asset paths to content hashes.
// This is used to detect when any JS/CSS file changes so we can invalidate cached HTML.
func HashAssetMap(assetHashes map[string]string) string {
	if len(assetHashes) == 0 {
		return ""
	}

	// Get sorted list of paths for deterministic ordering
	paths := make([]string, 0, len(assetHashes))
	for path := range assetHashes {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Combine path+hash pairs
	h := sha256.New()
	for _, path := range paths {
		h.Write([]byte(path))
		h.Write([]byte(assetHashes[path]))
	}

	return hex.EncodeToString(h.Sum(nil))
}

// ComputePostInputHash computes the input hash for a post.
// This combines: content + frontmatter fields that affect output + template name.
func ComputePostInputHash(content, frontmatter, template string) string {
	h := sha256.New()
	h.Write([]byte(content))
	h.Write([]byte("\x00")) // Separator
	h.Write([]byte(frontmatter))
	h.Write([]byte("\x00"))
	h.Write([]byte(template))
	return hex.EncodeToString(h.Sum(nil))
}

// SetDependencies records what targets a source post links to.
// This delegates to the underlying DependencyGraph.
func (c *Cache) SetDependencies(sourcePath, sourceSlug string, targets []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Graph.SetDependencies(sourcePath, sourceSlug, targets) {
		c.dirty = true
	}
}

// GetAffectedPosts returns all posts that need rebuilding when the given
// slugs change. This performs transitive closure via the dependency graph.
func (c *Cache) GetAffectedPosts(changedSlugs []string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Graph.GetAffectedPosts(changedSlugs)
}

// GetChangedSlugs returns the slugs that changed during this build.
func (c *Cache) GetChangedSlugs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]string, 0, len(c.changedSlugs))
	for slug := range c.changedSlugs {
		result = append(result, slug)
	}
	sort.Strings(result)
	return result
}

// GetChangedFeedSlugs returns slugs that changed in feed-relevant ways this build.
func (c *Cache) GetChangedFeedSlugs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.changedFeedSlugs) == 0 {
		return nil
	}
	result := make([]string, 0, len(c.changedFeedSlugs))
	for slug := range c.changedFeedSlugs {
		result = append(result, slug)
	}
	sort.Strings(result)
	return result
}

// ClearChangedFeedSlugs resets the feed change tracking for the current build.
func (c *Cache) ClearChangedFeedSlugs() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changedFeedSlugs = make(map[string]bool)
}

// MarkSlugChanged records that a slug changed this build.
// Used for dependency invalidation.
func (c *Cache) MarkSlugChanged(slug string) {
	if slug == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changedSlugs[slug] = true
}

// MarkFeedSlugChanged records that a slug changed in feed-relevant ways.
func (c *Cache) MarkFeedSlugChanged(slug string) {
	if slug == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changedFeedSlugs[slug] = true
}

// SetPostSlug records the slug for a source path in the dependency graph.
func (c *Cache) SetPostSlug(sourcePath, slug string) {
	if sourcePath == "" || slug == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Graph.PathToSlug[sourcePath] = slug
}

// MarkChangedPaths records slugs for the provided source paths if known in the cache.
// This is used to seed dependency invalidation from filesystem change events.
func (c *Cache) MarkChangedPaths(paths []string) {
	if len(paths) == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, path := range paths {
		if path == "" {
			continue
		}
		if slug := c.Graph.PathToSlug[path]; slug != "" {
			c.changedSlugs[slug] = true
		}
	}
}

// MarkMissingPaths records slugs for cached posts that are absent from the
// current glob result. The stale entries remain until cleanup so output
// cleanup and dependent invalidation can use their cached ownership data.
func (c *Cache) MarkMissingPaths(currentPaths []string) []string {
	if c == nil {
		return nil
	}

	current := make(map[string]struct{}, len(currentPaths))
	for _, path := range currentPaths {
		if path != "" {
			current[path] = struct{}{}
		}
	}

	type missingPost struct {
		path string
		slug string
	}

	c.mu.RLock()
	missing := make([]missingPost, 0)
	seen := make(map[string]struct{}, len(c.Posts)+len(c.stalePosts))
	for path, post := range c.Posts {
		if _, ok := current[path]; ok {
			continue
		}
		slug := ""
		if post != nil {
			slug = post.Slug
		}
		missing = append(missing, missingPost{path: path, slug: slug})
		seen[path] = struct{}{}
	}
	for path, post := range c.stalePosts {
		if _, ok := current[path]; ok {
			continue
		}
		if _, ok := seen[path]; ok {
			continue
		}
		slug := ""
		if post != nil {
			slug = post.Slug
		}
		missing = append(missing, missingPost{path: path, slug: slug})
	}
	c.mu.RUnlock()

	for i := range missing {
		if missing[i].slug == "" && c.Graph != nil {
			missing[i].slug = c.Graph.SlugForPath(missing[i].path)
		}
		if missing[i].slug != "" {
			c.MarkSlugChanged(missing[i].slug)
		}
	}

	paths := make([]string, 0, len(missing))
	for _, post := range missing {
		paths = append(paths, post.path)
	}
	sort.Strings(paths)
	return paths
}

// MarkAffectedDependents records all dependents (transitive) of changed slugs as changed.
func (c *Cache) MarkAffectedDependents(changedSlugs []string) {
	if len(changedSlugs) == 0 {
		return
	}

	c.mu.RLock()
	affectedPaths := c.Graph.GetAffectedPosts(changedSlugs)
	c.mu.RUnlock()
	if len(affectedPaths) == 0 {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, path := range affectedPaths {
		if slug := c.Graph.PathToSlug[path]; slug != "" {
			c.changedSlugs[slug] = true
		}
	}
}

// GetGlobCache returns the cached glob file list and pattern hash.
func (c *Cache) GetGlobCache() (files []string, patternHash string) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.GlobFiles, c.GlobPatternHash
}

// SetGlobCache stores the glob file list and pattern hash.
func (c *Cache) SetGlobCache(files []string, patternHash string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.GlobFiles = files
	c.GlobPatternHash = patternHash
	c.dirty = true
}

// GraphSize returns the number of posts with dependencies tracked.
func (c *Cache) GraphSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Graph.Size()
}

// HTMLCacheDir is the subdirectory for cached rendered HTML files.
const HTMLCacheDir = "html-cache"

// EncryptedHTMLCacheDir is the subdirectory for cached encrypted HTML files.
const EncryptedHTMLCacheDir = "encrypted-html-cache"

// MermaidSVGCacheDir is the subdirectory for cached prerendered Mermaid SVGs.
const MermaidSVGCacheDir = "mermaid-svg-cache"

// GetCachedArticleHTML returns the cached rendered HTML for a post if available.
// Returns empty string if not cached or cache is stale.
func (c *Cache) GetCachedArticleHTML(sourcePath, contentHash string) string {
	c.mu.RLock()
	cached, ok := c.Posts[sourcePath]
	c.mu.RUnlock()

	if !ok || cached.ContentHash != contentHash || cached.ArticleHTMLPath == "" {
		return ""
	}

	// Check in-memory cache first (populated by preloadCaches)
	if val, ok := c.articleHTMLMemory.Load(cached.ArticleHTMLPath); ok {
		if html, ok := val.(string); ok {
			return html
		}
	}

	// Fallback to disk read
	html, err := readCacheTextFile(cached.ArticleHTMLPath)
	if err != nil {
		return ""
	}

	// Store in memory for future calls
	c.articleHTMLMemory.Store(cached.ArticleHTMLPath, html)
	return html
}

// CacheArticleHTML stores rendered HTML for a post.
func (c *Cache) CacheArticleHTML(sourcePath, contentHash, articleHTML string) error {
	htmlPath, err := c.cacheHTMLFile(HTMLCacheDir, contentHash, articleHTML, &c.articleHTMLMemory)
	if err != nil {
		return err
	}

	// Update cache metadata
	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.ContentHash = contentHash
		cached.ArticleHTMLPath = htmlPath
	} else {
		c.Posts[sourcePath] = &PostCache{
			ContentHash:     contentHash,
			ArticleHTMLPath: htmlPath,
		}
	}
	c.dirty = true

	return nil
}

// GetCachedMermaidSVG returns cached prerendered Mermaid SVG output if the hash matches.
func (c *Cache) GetCachedMermaidSVG(sourcePath, renderHash string) string {
	if sourcePath == "" || renderHash == "" {
		return ""
	}
	svgPath := c.mermaidSVGPath(renderHash)
	c.markMermaidSVGUsed(svgPath)

	if val, ok := c.mermaidSVGMemory.Load(svgPath); ok {
		if svg, ok := val.(string); ok {
			return svg
		}
	}

	data, err := os.ReadFile(svgPath)
	if err != nil {
		return ""
	}

	svg := string(data)
	c.mermaidSVGMemory.Store(svgPath, svg)
	return svg
}

// CacheMermaidSVG stores prerendered Mermaid SVG output keyed by render hash.
func (c *Cache) CacheMermaidSVG(sourcePath, renderHash, svg string) error {
	if sourcePath == "" || renderHash == "" {
		return nil
	}
	c.markMermaidSVGUsed(c.mermaidSVGPath(renderHash))
	_, err := c.cacheHTMLFile(MermaidSVGCacheDir, renderHash, svg, &c.mermaidSVGMemory)
	return err
}

// CleanupMermaidSVG removes cached Mermaid SVG files not used in the current build.
func (c *Cache) CleanupMermaidSVG() (int, error) {
	cacheDir := filepath.Join(filepath.Dir(c.path), MermaidSVGCacheDir)
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("reading mermaid svg cache dir: %w", err)
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(cacheDir, entry.Name())
		if _, ok := c.mermaidSVGUsed.Load(path); ok {
			continue
		}
		if err := os.Remove(path); err != nil {
			return removed, fmt.Errorf("removing stale mermaid svg cache %s: %w", path, err)
		}
		c.mermaidSVGMemory.Delete(path)
		removed++
	}

	return removed, nil
}

// ResetMermaidSVGUsage clears per-build Mermaid SVG usage tracking.
func (c *Cache) ResetMermaidSVGUsage() {
	c.mermaidSVGUsed = sync.Map{}
}

func (c *Cache) mermaidSVGPath(renderHash string) string {
	cacheDir := filepath.Dir(c.path)
	hashPrefix := renderHash
	if len(hashPrefix) > 16 {
		hashPrefix = hashPrefix[:16]
	}
	return filepath.Join(cacheDir, MermaidSVGCacheDir, hashPrefix+".html")
}

func (c *Cache) markMermaidSVGUsed(path string) {
	if path == "" {
		return
	}
	c.mermaidSVGUsed.Store(path, struct{}{})
}

// ContentHash computes a hash of just the markdown content.
func ContentHash(content string) string {
	h := sha256.New()
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}

// GetFeedHash returns the cached hash for a feed, or empty string if not cached.
func (c *Cache) GetFeedHash(slug string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cached, ok := c.Feeds[slug]; ok {
		return cached.Hash
	}
	return ""
}

// SetFeedHash stores the hash for a feed in the cache.
func (c *Cache) SetFeedHash(slug, hash string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Feeds[slug] = &FeedCache{Hash: hash}
	c.dirty = true
}

// FullHTMLCacheDir is the subdirectory for cached full page HTML files.
const FullHTMLCacheDir = "fullhtml-cache"

// GetCachedFullHTML returns the cached full page HTML for a post if available.
func (c *Cache) GetCachedFullHTML(sourcePath string) string {
	c.mu.RLock()
	cached, ok := c.Posts[sourcePath]
	var htmlPath string
	if ok {
		htmlPath = cached.FullHTMLPath
	}
	c.mu.RUnlock()

	if htmlPath == "" {
		return ""
	}

	// Check in-memory cache first (populated by preloadCaches)
	if val, ok := c.fullHTMLMemory.Load(htmlPath); ok {
		if html, ok := val.(string); ok {
			return html
		}
	}

	// Fallback to disk read
	html, err := readCacheTextFile(htmlPath)
	if err != nil {
		return ""
	}

	// Concurrent restorations of a shared path must retain one backing string.
	if value, loaded := c.fullHTMLMemory.LoadOrStore(htmlPath, html); loaded {
		if cachedHTML, ok := value.(string); ok {
			return cachedHTML
		}
	}
	return html
}

// InvalidateFullHTML clears a post's full-page cache reference and persists the
// removal on Save. Other metadata and potentially shared cache files remain.
func (c *Cache) InvalidateFullHTML(sourcePath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if cached, ok := c.Posts[sourcePath]; ok && cached.FullHTMLPath != "" {
		cached.FullHTMLPath = ""
		c.dirty = true
	}
}

// CacheFullHTML stores the full page HTML for a post.
func (c *Cache) CacheFullHTML(sourcePath, fullHTML string) error {
	cacheDir := filepath.Dir(c.path)
	htmlCacheDir := filepath.Join(cacheDir, FullHTMLCacheDir)

	if err := os.MkdirAll(htmlCacheDir, 0o755); err != nil {
		return fmt.Errorf("creating fullhtml cache dir: %w", err)
	}

	// Use input hash as filename
	c.mu.RLock()
	cached, ok := c.Posts[sourcePath]
	var hashPrefix string
	if ok && cached.InputHash != "" {
		hashPrefix = cached.InputHash
		if len(hashPrefix) > 16 {
			hashPrefix = hashPrefix[:16]
		}
	} else {
		// Generate hash from path
		h := sha256.New()
		h.Write([]byte(sourcePath))
		hashPrefix = hex.EncodeToString(h.Sum(nil))[:16]
	}
	c.mu.RUnlock()

	htmlPath := filepath.Join(htmlCacheDir, hashPrefix+".html")

	//nolint:gosec // G306: cache files need 0644 for reading by other processes
	if err := os.WriteFile(htmlPath, []byte(fullHTML), 0o644); err != nil {
		return fmt.Errorf("writing fullhtml cache: %w", err)
	}

	// Store in memory cache for future reads within this build
	c.fullHTMLMemory.Store(htmlPath, fullHTML)

	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.FullHTMLPath = htmlPath
	} else {
		c.Posts[sourcePath] = &PostCache{
			FullHTMLPath: htmlPath,
		}
	}
	c.dirty = true

	return nil
}

// GetCachedEncryptedHTML returns cached encrypted HTML if hash matches.
func (c *Cache) GetCachedEncryptedHTML(sourcePath, encryptedHash string) string {
	c.mu.RLock()
	cached, ok := c.Posts[sourcePath]
	c.mu.RUnlock()
	if !ok || cached.EncryptedHash != encryptedHash || cached.EncryptedHTMLPath == "" {
		return ""
	}

	if val, ok := c.encryptedHTMLMemory.Load(cached.EncryptedHTMLPath); ok {
		if html, ok := val.(string); ok {
			return html
		}
	}

	data, err := os.ReadFile(cached.EncryptedHTMLPath)
	if err != nil {
		return ""
	}

	html := string(data)
	c.encryptedHTMLMemory.Store(cached.EncryptedHTMLPath, html)
	return html
}

// CacheEncryptedHTML stores encrypted HTML for a post.
func (c *Cache) CacheEncryptedHTML(sourcePath, encryptedHash, html string) error {
	htmlPath, err := c.cacheHTMLFile(EncryptedHTMLCacheDir, encryptedHash, html, &c.encryptedHTMLMemory)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	if cached, ok := c.Posts[sourcePath]; ok {
		cached.EncryptedHash = encryptedHash
		cached.EncryptedHTMLPath = htmlPath
	} else {
		c.Posts[sourcePath] = &PostCache{
			EncryptedHash:     encryptedHash,
			EncryptedHTMLPath: htmlPath,
		}
	}
	c.dirty = true
	return nil
}

func (c *Cache) cacheHTMLFile(cacheDirName, hash, html string, memory *sync.Map) (string, error) {
	cacheDir := filepath.Dir(c.path)
	cacheSubDir := filepath.Join(cacheDir, cacheDirName)
	if err := os.MkdirAll(cacheSubDir, 0o755); err != nil {
		return "", fmt.Errorf("creating html cache dir: %w", err)
	}

	hashPrefix := hash
	if len(hashPrefix) > 16 {
		hashPrefix = hashPrefix[:16]
	}
	htmlPath := filepath.Join(cacheSubDir, hashPrefix+".html")

	//nolint:gosec // G306: cache files need 0644 for reading by other processes
	if err := os.WriteFile(htmlPath, []byte(html), 0o644); err != nil {
		return "", fmt.Errorf("writing html cache: %w", err)
	}

	if memory != nil {
		memory.Store(htmlPath, html)
	}

	return htmlPath, nil
}

// PostCacheDir is the subdirectory for cached parsed post JSON files.
const PostCacheDir = "post-cache"

// CachedPostData holds the serializable parts of a Post for caching.
// This excludes rendered HTML which is cached separately.
type CachedPostData struct {
	Path               string              `json:"path"`
	Content            string              `json:"content"`
	Slug               string              `json:"slug"`
	Href               string              `json:"href"`
	Title              *string             `json:"title,omitempty"`
	Date               *time.Time          `json:"date,omitempty"`
	Modified           *time.Time          `json:"modified,omitempty"`
	Published          bool                `json:"published"`
	Draft              bool                `json:"draft"`
	Private            bool                `json:"private"`
	PrivateOverride    *bool               `json:"private_override,omitempty"`
	Skip               bool                `json:"skip"`
	SecretKey          string              `json:"secret_key,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	Description        *string             `json:"description,omitempty"`
	Template           string              `json:"template"`
	Templates          map[string]string   `json:"templates,omitempty"`
	FrontmatterPresent bool                `json:"frontmatter_present"`
	FrontmatterValid   bool                `json:"frontmatter_valid"`
	FrontmatterIssues  []diagnostics.Issue `json:"frontmatter_issues,omitempty"`
	RawFrontmatter     string              `json:"raw_frontmatter"`
	InputHash          string              `json:"input_hash"`
	Authors            []string            `json:"authors,omitempty"`
	Author             *string             `json:"author,omitempty"`
	Extra              map[string]any      `json:"extra,omitempty"`
}

// GetCachedPostData returns cached post data if ModTime matches.
// Returns nil if post is not cached or file has changed.
func (c *Cache) GetCachedPostData(sourcePath string, modTime int64) *CachedPostData {
	c.mu.RLock()
	cached, ok := c.Posts[sourcePath]
	c.mu.RUnlock()

	if !ok || cached.ModTime != modTime || cached.ModTime == 0 {
		return nil
	}

	// Compute the post cache file path
	cacheDir := filepath.Dir(c.path)
	postCacheDir := filepath.Join(cacheDir, PostCacheDir)
	h := sha256.Sum256([]byte(sourcePath))
	hashPrefix := hex.EncodeToString(h[:])[:16]
	postPath := filepath.Join(postCacheDir, hashPrefix+".json")

	// Check in-memory cache first (populated by preloadCaches)
	if val, ok := c.postDataMemory.Load(postPath); ok {
		if pd, ok := val.(*CachedPostData); ok {
			return pd
		}
	}

	// Fallback to disk read
	data, err := os.ReadFile(postPath)
	if err != nil {
		return nil
	}

	var postData CachedPostData
	if err := json.Unmarshal(data, &postData); err != nil {
		return nil
	}

	// Store in memory for future calls
	c.postDataMemory.Store(postPath, &postData)
	return &postData
}

// GetCachedPostDataLatest returns cached post data without checking ModTime.
// Use only when external change detection is trusted.
func (c *Cache) GetCachedPostDataLatest(sourcePath string) *CachedPostData {
	if sourcePath == "" {
		return nil
	}

	// Compute the post cache file path
	cacheDir := filepath.Dir(c.path)
	postCacheDir := filepath.Join(cacheDir, PostCacheDir)
	h := sha256.Sum256([]byte(sourcePath))
	hashPrefix := hex.EncodeToString(h[:])[:16]
	postPath := filepath.Join(postCacheDir, hashPrefix+".json")

	// Check in-memory cache first
	if val, ok := c.postDataMemory.Load(postPath); ok {
		if pd, ok := val.(*CachedPostData); ok {
			return pd
		}
	}

	data, err := os.ReadFile(postPath)
	if err != nil {
		return nil
	}

	var postData CachedPostData
	if err := json.Unmarshal(data, &postData); err != nil {
		return nil
	}

	c.postDataMemory.Store(postPath, &postData)
	return &postData
}

// CachePostData stores parsed post data to disk.
func (c *Cache) CachePostData(sourcePath string, modTime int64, postData *CachedPostData) error {
	cacheDir := filepath.Dir(c.path)
	postCacheDir := filepath.Join(cacheDir, PostCacheDir)

	if err := os.MkdirAll(postCacheDir, 0o755); err != nil {
		return fmt.Errorf("creating post cache dir: %w", err)
	}

	h := sha256.Sum256([]byte(sourcePath))
	hashPrefix := hex.EncodeToString(h[:])[:16]
	postPath := filepath.Join(postCacheDir, hashPrefix+".json")

	data, err := json.Marshal(postData)
	if err != nil {
		return fmt.Errorf("marshaling post data: %w", err)
	}

	//nolint:gosec // G306: cache files need 0644 for reading
	if err := os.WriteFile(postPath, data, 0o644); err != nil {
		return fmt.Errorf("writing post cache: %w", err)
	}

	// Store in memory cache for future reads within this build
	c.postDataMemory.Store(postPath, postData)

	// Update ModTime in cache
	c.mu.Lock()
	defer c.mu.Unlock()

	if cached, ok := c.Posts[sourcePath]; ok {
		cached.ModTime = modTime
		cached.Slug = postData.Slug
	} else {
		c.Posts[sourcePath] = &PostCache{
			ModTime: modTime,
			Slug:    postData.Slug,
		}
	}
	c.dirty = true

	return nil
}
