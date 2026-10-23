// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/WaylonWalker/markata/pkg/blogroll"
	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// categoryUncategorized is the default category name for feeds without a category.
const categoryUncategorized = "Uncategorized"

// Default slug constants for blogroll pages.
const (
	defaultBlogrollSlug  = "blogroll"
	defaultReaderSlug    = "reader"
	blogrollCacheVersion = "community-meta-v4"
	aggregateCacheFile   = "blogroll-aggregate-v1.json"
)

var readerFeedTypeOrder = []models.ReaderFeedType{
	models.ReaderFeedTypeWritten,
	models.ReaderFeedTypeVideo,
	models.ReaderFeedTypePodcast,
}

type readerVariant struct {
	Key         string
	Title       string
	Description string
	BaseURL     string
	OutputDir   string
	FeedType    models.ReaderFeedType
	EntryCount  int
}

// Default directory constants.
const (
	defaultOutputDir  = "output"
	blogrollBundleDir = "blogroll"
	mediumImage       = "image"
)

// extractFirstImageRegex matches the first img src attribute in HTML.
var extractFirstImageRegex = regexp.MustCompile(`<img[^>]+src\s*=\s*["']([^"']+)["']`)

// extractFirstImageFromHTML extracts the first image URL from HTML content.
func extractFirstImageFromHTML(htmlContent string) string {
	// Decode HTML entities first
	decoded := html.UnescapeString(htmlContent)

	// Simple regex to find first img src attribute
	matches := extractFirstImageRegex.FindStringSubmatch(decoded)
	if len(matches) > 1 {
		return matches[1]
	}

	// Try to extract YouTube thumbnail if no img tag found
	if thumbnail := extractYouTubeThumbnail(decoded); thumbnail != "" {
		return thumbnail
	}

	return ""
}

// youtubeThumbnailRegex matches YouTube URLs and extracts the video ID.
var youtubeThumbnailRegex = regexp.MustCompile(
	`https?://(?:www\.)?(?:youtube\.com/(?:watch\?v=|shorts/)|youtu\.be/)([a-zA-Z0-9_-]{11})(?:[^\s<]*)`,
)

// youtubeVideoIDRegex extracts a video ID from YouTube Atom IDs like yt:video:abc123def45.
var youtubeVideoIDRegex = regexp.MustCompile(`(?:^|:)video:([a-zA-Z0-9_-]{11})$`)

// youtubeThumbnailURL returns a reliable YouTube thumbnail URL for a video ID.
func youtubeThumbnailURL(videoID string) string {
	if len(videoID) != 11 {
		return ""
	}
	return "https://img.youtube.com/vi/" + videoID + "/hqdefault.jpg"
}

// extractYouTubeThumbnail extracts a YouTube thumbnail URL from content.
// Returns a stable thumbnail URL if a YouTube URL is found.
func extractYouTubeThumbnail(content string) string {
	matches := youtubeThumbnailRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return youtubeThumbnailURL(matches[1])
	}
	return ""
}

func extractYouTubeThumbnailFromID(id string) string {
	matches := youtubeVideoIDRegex.FindStringSubmatch(strings.TrimSpace(id))
	if len(matches) > 1 {
		return youtubeThumbnailURL(matches[1])
	}
	return ""
}

// extractEntryImage extracts an image URL from feed entry metadata.
// Priority: media:thumbnail -> media:content -> enclosure -> YouTube -> HTML img
func extractEntryImage(item *rss2Item) string {
	// Try media:thumbnail first (most specific for images)
	for _, thumb := range item.MediaThumb {
		if thumb.URL != "" {
			return thumb.URL
		}
	}

	// Try media:content with medium=image
	for _, mc := range item.MediaContent {
		if mc.URL != "" && mc.Medium == mediumImage {
			return mc.URL
		}
	}

	// Try enclosure if it's an image
	if item.Enclosure.URL != "" && strings.HasPrefix(item.Enclosure.Type, "image/") {
		return item.Enclosure.URL
	}

	// Combine content and description for YouTube/HTML extraction
	combined := item.Content
	if combined == "" {
		combined = item.Description
	}

	// Try YouTube thumbnail from combined content
	if thumbnail := extractYouTubeThumbnail(combined); thumbnail != "" {
		return thumbnail
	}

	// Fall back to HTML img extraction
	return extractFirstImageFromHTML(combined)
}

// getAtomEnclosureURL extracts enclosure URL from Atom entry links.
func getAtomEnclosureURL(links []atomLink) string {
	for _, link := range links {
		if link.Rel == "enclosure" && strings.HasPrefix(link.Type, "image/") {
			return link.Href
		}
	}
	return ""
}

// extractAtomEntryImage extracts an image URL from Atom feed entry metadata.
func extractAtomEntryImage(entry *atomEntry) string {
	// Try media:thumbnail first
	for _, thumb := range entry.MediaThumb {
		if thumb.URL != "" {
			return thumb.URL
		}
	}

	// Try media:content with medium=image
	for _, mc := range entry.MediaContent {
		if mc.URL != "" && mc.Medium == mediumImage {
			return mc.URL
		}
	}

	// Try enclosure from links
	if enclosureURL := getAtomEnclosureURL(entry.Link); enclosureURL != "" {
		return enclosureURL
	}

	// Combine content and summary for YouTube/HTML extraction
	combined := entry.Content.Body
	if combined == "" {
		combined = entry.Summary
	}

	// Try YouTube thumbnail
	if thumbnail := extractYouTubeThumbnail(combined); thumbnail != "" {
		return thumbnail
	}

	// YouTube Atom feeds include stable entry IDs like yt:video:<id>.
	if thumbnail := extractYouTubeThumbnailFromID(entry.ID); thumbnail != "" {
		return thumbnail
	}

	for _, link := range entry.Link {
		if thumbnail := extractYouTubeThumbnail(link.Href); thumbnail != "" {
			return thumbnail
		}
	}

	// Fall back to HTML img extraction
	return extractFirstImageFromHTML(combined)
}

// blogrollParsedFeed represents a parsed feed response for blogroll plugin.
type blogrollParsedFeed struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Language     string     `json:"language"`
	SiteURL      string     `json:"site_url"`
	ImageURL     string     `json:"image_url"`
	AvatarURL    string     `json:"avatar_url"`
	AvatarSource string     `json:"avatar_source"`
	LastUpdated  *time.Time `json:"last_updated"`
}

// RSS 2.0 structures for feed parsing.
type rss2Feed struct {
	XMLName xml.Name    `xml:"rss"`
	Channel rss2Channel `xml:"channel"`
}

type rss2Channel struct {
	Title       string     `xml:"title"`
	Link        string     `xml:"link"`
	Description string     `xml:"description"`
	Language    string     `xml:"language"`
	Image       rss2Image  `xml:"image"`
	Items       []rss2Item `xml:"item"`
}

type rss2Image struct {
	URL string `xml:"url"`
}

type rss2Item struct {
	Title        string               `xml:"title"`
	Link         string               `xml:"link"`
	Description  string               `xml:"description"`
	Content      string               `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
	PubDate      string               `xml:"pubDate"`
	GUID         string               `xml:"guid"`
	Comments     string               `xml:"comments"`
	Author       string               `xml:"author"`
	Creator      string               `xml:"http://purl.org/dc/elements/1.1/ creator"`
	Categories   []string             `xml:"category"`
	Enclosure    rss2Enclosure        `xml:"enclosure"`
	MediaThumb   []rss2MediaThumbnail `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaContent []rss2MediaContent   `xml:"http://search.yahoo.com/mrss/ content"`
}

type rss2Enclosure struct {
	URL  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type rss2MediaThumbnail struct {
	URL string `xml:"url,attr"`
}

type rss2MediaContent struct {
	URL    string `xml:"url,attr"`
	Medium string `xml:"medium,attr"`
}

// Atom structures for feed parsing.
type atomFeed struct {
	XMLName  xml.Name    `xml:"feed"`
	Title    string      `xml:"title"`
	Subtitle string      `xml:"subtitle"`
	Link     []atomLink  `xml:"link"`
	Icon     string      `xml:"icon"`
	Logo     string      `xml:"logo"`
	Entries  []atomEntry `xml:"entry"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type atomEntry struct {
	Title        string               `xml:"title"`
	Link         []atomLink           `xml:"link"`
	ID           string               `xml:"id"`
	Updated      string               `xml:"updated"`
	Published    string               `xml:"published"`
	Summary      string               `xml:"summary"`
	Content      atomText             `xml:"content"`
	Author       atomAuthor           `xml:"author"`
	Category     []atomCat            `xml:"category"`
	MediaThumb   []atomMediaThumbnail `xml:"http://search.yahoo.com/mrss/ thumbnail"`
	MediaContent []atomMediaContent   `xml:"http://search.yahoo.com/mrss/ content"`
}

type atomText struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

type atomAuthor struct {
	Name string `xml:"name"`
}

type atomCat struct {
	Term string `xml:"term,attr"`
}

type atomMediaThumbnail struct {
	URL string `xml:"url,attr"`
}

type atomMediaContent struct {
	URL    string `xml:"url,attr"`
	Medium string `xml:"medium,attr"`
}

// parseBlogrollFeedResponse parses an HTTP response into a feed structure for blogroll plugin.
func parseBlogrollFeedResponse(resp *http.Response) (*blogrollParsedFeed, []*models.ExternalEntry, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read response: %w", err)
	}

	// Try RSS 2.0 first
	feed, entries, err := parseRSS2Feed(body)
	if err == nil {
		// Successfully parsed as RSS 2.0
		return feed, entries, nil
	}
	rssErr := err

	// Try Atom
	feed, entries, err = parseAtomFeed(body)
	if err == nil {
		// Successfully parsed as Atom
		return feed, entries, nil
	}

	// Both parsers failed - return combined error
	return nil, nil, fmt.Errorf("failed to parse feed: %w", errors.Join(
		fmt.Errorf("rss 2.0: %w", rssErr),
		fmt.Errorf("atom: %w", err),
	))
}

// parseRSS2Feed parses an RSS 2.0 feed.
func parseRSS2Feed(data []byte) (*blogrollParsedFeed, []*models.ExternalEntry, error) {
	var feed rss2Feed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, nil, err
	}

	if feed.XMLName.Local != "rss" {
		return nil, nil, fmt.Errorf("not an RSS 2.0 feed")
	}

	parsed := &blogrollParsedFeed{
		Title:       feed.Channel.Title,
		Description: feed.Channel.Description,
		Language:    feed.Channel.Language,
		SiteURL:     normalizeExternalURL(feed.Channel.Link),
		ImageURL:    normalizeExternalURL(feed.Channel.Image.URL),
	}

	entries := make([]*models.ExternalEntry, 0, len(feed.Channel.Items))
	for i := range feed.Channel.Items {
		item := &feed.Channel.Items[i]

		// Parse publish date
		pubDate := parseRSSDate(item.PubDate)
		var pubDatePtr *time.Time
		if !pubDate.IsZero() {
			pubDatePtr = &pubDate
		}

		// Determine content - prefer content:encoded over description
		content := item.Content
		if content == "" {
			content = item.Description
		}

		// Determine author
		author := item.Author
		if author == "" {
			author = item.Creator
		}

		// Determine ID
		id := item.GUID
		if id == "" {
			id = item.Link
		}

		// Extract image from various sources (priority order)
		imageURL := extractEntryImage(item)
		discussionURL := item.GUID
		if isLobstersDiscussionURL(item.Comments) || isLemmyDiscussionURL(item.Comments) {
			discussionURL = item.Comments
		}
		resolvedURL, originalURL, metadata := enrichAggregatorFeedEntry(item.Link, discussionURL, content)
		if metadata != nil {
			if metadata.Title != "" {
				item.Title = metadata.Title
			}
			if metadata.Description != "" {
				item.Description = metadata.Description
			}
			if metadata.Image != "" {
				imageURL = metadata.Image
			}
		}

		// Create entry
		entryURL := normalizeExternalURL(resolvedURL)
		entry := &models.ExternalEntry{
			ID:          id,
			Title:       item.Title,
			URL:         entryURL,
			OriginalURL: originalURL,
			Published:   pubDatePtr,
			Updated:     pubDatePtr,
			Author:      author,
			Content:     content,
			Description: stripBlogrollHTML(item.Description),
			ImageURL:    normalizeReaderImageURL(imageURL, resolvedURL),
			Categories:  item.Categories,
		}

		entries = append(entries, entry)
	}

	return parsed, entries, nil
}

// parseAtomFeed parses an Atom feed.
func parseAtomFeed(data []byte) (*blogrollParsedFeed, []*models.ExternalEntry, error) {
	var feed atomFeed
	if err := xml.Unmarshal(data, &feed); err != nil {
		return nil, nil, err
	}

	// Check if it's actually an Atom feed
	if feed.XMLName.Local != DefaultFeedPath {
		return nil, nil, fmt.Errorf("not an Atom feed")
	}

	// Get site URL from links
	siteURL := ""
	for _, link := range feed.Link {
		if link.Rel == "" || link.Rel == "alternate" {
			siteURL = link.Href
			break
		}
	}

	// Prefer logo over icon
	imageURL := feed.Logo
	if imageURL == "" {
		imageURL = feed.Icon
	}

	parsed := &blogrollParsedFeed{
		Title:       feed.Title,
		Description: feed.Subtitle,
		SiteURL:     normalizeExternalURL(siteURL),
		ImageURL:    normalizeExternalURL(imageURL),
	}

	entries := make([]*models.ExternalEntry, 0, len(feed.Entries))
	for i := range feed.Entries {
		entry := &feed.Entries[i]

		// Parse dates
		pubDate := parseAtomDate(entry.Published)
		updDate := parseAtomDate(entry.Updated)

		var pubDatePtr, updDatePtr *time.Time
		if !pubDate.IsZero() {
			pubDatePtr = &pubDate
		}
		if !updDate.IsZero() {
			updDatePtr = &updDate
		}

		// Get entry URL
		entryURL := ""
		for _, link := range entry.Link {
			if link.Rel == "" || link.Rel == "alternate" {
				entryURL = link.Href
				break
			}
		}

		// Determine content
		content := entry.Content.Body
		if content == "" {
			content = entry.Summary
		}

		// Get tags
		var tags []string
		for _, cat := range entry.Category {
			if cat.Term != "" {
				tags = append(tags, cat.Term)
			}
		}

		// Extract first image from content
		imageURL := extractAtomEntryImage(entry)
		resolvedURL, originalURL, metadata := enrichAggregatorFeedEntry(entryURL, entry.ID, content)
		if metadata != nil {
			if metadata.Title != "" {
				entry.Title = metadata.Title
			}
			if metadata.Description != "" {
				entry.Summary = metadata.Description
			}
			if metadata.Image != "" {
				imageURL = metadata.Image
			}
		}

		// Create entry
		extEntry := &models.ExternalEntry{
			ID:          entry.ID,
			Title:       entry.Title,
			URL:         normalizeExternalURL(resolvedURL),
			OriginalURL: originalURL,
			Published:   pubDatePtr,
			Updated:     updDatePtr,
			Author:      entry.Author.Name,
			Content:     content,
			Description: stripBlogrollHTML(entry.Summary),
			ImageURL:    normalizeReaderImageURL(imageURL, resolvedURL),
			Categories:  tags,
		}

		entries = append(entries, extEntry)
	}

	return parsed, entries, nil
}

// parseRSSDate parses an RSS date string.
func parseRSSDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	// RFC822, RFC822Z, RFC1123, RFC1123Z are common RSS date formats
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
		"2 Jan 2006 15:04:05 -0700",
		"2 Jan 2006 15:04:05 MST",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Time{}
}

// parseAtomDate parses an Atom date string (RFC3339).
func parseAtomDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Time{}
	}

	return t
}

// stripBlogrollHTML strips HTML tags from content.
func stripBlogrollHTML(s string) string {
	// Decode HTML entities
	decoded := html.UnescapeString(s)
	// Remove HTML tags
	stripped := blogrollHTMLTagRegex.ReplaceAllString(decoded, "")
	// Trim whitespace
	return strings.TrimSpace(stripped)
}

// blogrollHTMLTagRegex matches HTML tags for stripping.
var blogrollHTMLTagRegex = regexp.MustCompile(`<[^>]*>`)

// BlogrollPlugin fetches and processes external RSS/Atom feeds.
// It runs in the Configure stage to register synthetic posts for wikilink resolution,
// in the Collect stage to gather external feed entries,
// and in the Write stage to generate blogroll and reader pages.
type BlogrollPlugin struct {
	feeds     []*models.ExternalFeed
	entries   []*models.ExternalEntry
	mu        sync.RWMutex
	assetURLs map[string]string
	urlMentions map[string][]map[string]string
}

type blogrollFetchOptions struct {
	forceRefresh   bool
	refreshOnBuild bool
	onFeedComplete func(BlogrollCacheRefreshProgress)
}

type blogrollRefreshSummary struct {
	refreshed int
	stale     int
	failed    int
}

type blogrollAggregateCache struct {
	ConfigHash string                  `json:"config_hash"`
	Feeds      []*models.ExternalFeed  `json:"feeds"`
	Entries    []*models.ExternalEntry `json:"entries"`
}

// NewBlogrollPlugin creates a new BlogrollPlugin.
func NewBlogrollPlugin() *BlogrollPlugin {
	return &BlogrollPlugin{
		feeds:   make([]*models.ExternalFeed, 0),
		entries: make([]*models.ExternalEntry, 0),
	}
}

// Name returns the unique name of the plugin.
func (p *BlogrollPlugin) Name() string {
	return "blogroll"
}

// Priority returns the plugin's priority for a given stage.
func (p *BlogrollPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageConfigure:
		// Run early to register synthetic posts before wikilinks are processed
		return lifecycle.PriorityDefault
	case lifecycle.StageCollect:
		// Run after feeds plugin to not interfere
		return lifecycle.PriorityLate + 10
	case lifecycle.StageWrite:
		// Run after publish_feeds to write blogroll pages
		return lifecycle.PriorityLate + 20
	default:
		return lifecycle.PriorityDefault
	}
}

// Configure registers synthetic posts for blogroll and reader pages
// so they can be resolved by wikilinks during the Transform stage.
func (p *BlogrollPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		p.assetURLs = assetURLsFromConfig(config)
		if disabled, ok := config.Extra["blogroll_disabled"].(bool); ok && disabled {
			return nil
		}
	}
	blogrollConfig := getBlogrollConfig(config)

	if !blogrollConfig.Enabled {
		return nil
	}

	// Register synthetic posts for wikilink resolution
	p.registerSyntheticPosts(m, blogrollConfig)

	return nil
}

func (p *BlogrollPlugin) resolveAssetURL(name, fallback string) string {
	return resolveAssetURL(p.assetURLs, name, fallback)
}

// Collect fetches and parses configured external feeds.
func (p *BlogrollPlugin) Collect(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		if disabled, ok := config.Extra["blogroll_disabled"].(bool); ok && disabled {
			return nil
		}
	}
	blogrollConfig := getBlogrollConfig(config)

	if !blogrollConfig.Enabled {
		return nil
	}

	// Fetch all feeds concurrently
	feeds, entries, err := p.fetchFeeds(blogrollConfig)
	if err != nil {
		return fmt.Errorf("blogroll: %w", err)
	}

	p.mu.Lock()
	p.feeds = feeds
	p.entries = entries
	p.mu.Unlock()

	// Store in cache for templates
	m.Cache().Set("blogroll_feeds", feeds)
	m.Cache().Set("blogroll_entries", entries)
	m.Cache().Set("blogroll_categories", p.groupByCategory(feeds))

	return nil
}

// registerSyntheticPosts creates synthetic Post objects for blogroll and reader pages
// so they can be resolved by wikilinks. These posts are marked with Skip: true
// so they don't interfere with normal rendering.
func (p *BlogrollPlugin) registerSyntheticPosts(m *lifecycle.Manager, config models.BlogrollConfig) {
	// Get configured slugs with defaults
	blogrollSlug := config.BlogrollSlug
	if blogrollSlug == "" {
		blogrollSlug = defaultBlogrollSlug
	}
	readerSlug := config.ReaderSlug
	if readerSlug == "" {
		readerSlug = defaultReaderSlug
	}

	// Helper to create string pointer
	strPtr := func(s string) *string { return &s }

	// Register blogroll page
	blogrollPost := &models.Post{
		Slug:        blogrollSlug,
		Title:       strPtr("Blogroll"),
		Description: strPtr("Blogs and feeds I follow"),
		Href:        "/" + blogrollSlug + "/",
		Published:   true,
		Skip:        true,
	}
	m.AddPost(blogrollPost)

	// Register reader page
	readerPost := &models.Post{
		Slug:        readerSlug,
		Title:       strPtr("Reader"),
		Description: strPtr("Latest posts from blogs I follow"),
		Href:        "/" + readerSlug + "/",
		Published:   true,
		Skip:        true,
	}
	m.AddPost(readerPost)
}

// Write generates the blogroll and reader pages.
func (p *BlogrollPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		if disabled, ok := config.Extra["blogroll_disabled"].(bool); ok && disabled {
			return nil
		}
	}
	blogrollConfig := getBlogrollConfig(config)

	if !blogrollConfig.Enabled {
		return nil
	}

	p.mu.RLock()
	feeds := p.feeds
	entries := p.entries
	p.mu.RUnlock()

	if len(feeds) == 0 {
		return nil
	}

	outputDir := config.OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	// Opt-in provider indexes are fetched once, before any reader pages render.
	p.urlMentions = nil
	if config.Extra != nil {
		p.urlMentions = readerURLMentions(config.Extra["url_mentions_providers"])
	}

	// Generate blogroll page
	if err := p.writeBlogrollPage(m, outputDir, feeds, blogrollConfig); err != nil {
		return fmt.Errorf("blogroll page: %w", err)
	}

	// Generate reader page
	if err := p.writeReaderPage(m, outputDir, feeds, entries, blogrollConfig); err != nil {
		return fmt.Errorf("reader page: %w", err)
	}

	return nil
}

// fetchFeeds fetches all configured feeds concurrently.
func (p *BlogrollPlugin) fetchFeeds(config models.BlogrollConfig) ([]*models.ExternalFeed, []*models.ExternalEntry, error) {
	feeds, entries, _, err := p.fetchFeedsWithOptions(config, blogrollFetchOptions{refreshOnBuild: config.GetRefreshOnBuild()})
	return feeds, entries, err
}

// fetchFeedsWithOptions fetches all configured feeds concurrently with optional
// cache-bypass behavior for direct cache refresh commands.
func (p *BlogrollPlugin) fetchFeedsWithOptions(config models.BlogrollConfig, options blogrollFetchOptions) ([]*models.ExternalFeed, []*models.ExternalEntry, blogrollRefreshSummary, error) {
	var activeFeeds []models.ExternalFeedConfig
	for i := range config.Feeds {
		if config.Feeds[i].IsActive() {
			activeFeeds = append(activeFeeds, config.Feeds[i])
		}
	}

	if len(activeFeeds) == 0 {
		return nil, nil, blogrollRefreshSummary{}, nil
	}

	// Parse cache duration
	cacheDuration, err := time.ParseDuration(config.CacheDuration)
	if err != nil {
		cacheDuration = time.Hour
	}

	// Create cache directory
	if config.CacheDir != "" {
		if err := os.MkdirAll(config.CacheDir, 0o755); err != nil {
			return nil, nil, blogrollRefreshSummary{}, fmt.Errorf("create cache dir: %w", err)
		}
	}

	if config.CacheDir != "" && !options.forceRefresh {
		configHash := p.aggregateConfigHash(activeFeeds)
		if cached := p.loadAggregateCache(config.CacheDir, cacheDuration, configHash); cached != nil {
			return cached.Feeds, cached.Entries, blogrollRefreshSummary{}, nil
		}
		if !options.refreshOnBuild {
			if cached := p.loadAggregateCacheAny(config.CacheDir, configHash); cached != nil {
				return cached.Feeds, cached.Entries, blogrollRefreshSummary{stale: len(cached.Feeds)}, nil
			}
		}
	}

	// Fetch feeds concurrently
	concurrency := config.ConcurrentRequests
	if concurrency <= 0 {
		concurrency = 5
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	globalMaxEntries := config.MaxEntriesPerFeed
	if globalMaxEntries <= 0 {
		globalMaxEntries = 50
	}

	semaphore := make(chan struct{}, concurrency)
	resultsCh := make(chan *models.ExternalFeed, len(activeFeeds))
	summaryCh := make(chan blogrollRefreshSummary, len(activeFeeds))
	var wg sync.WaitGroup

	for i := range activeFeeds {
		wg.Add(1)
		go func(feedConfig models.ExternalFeedConfig) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Use per-feed max_entries if set, otherwise global default
			maxEntries := feedConfig.GetMaxEntries(globalMaxEntries)
			feed, summary := p.fetchFeed(feedConfig, config.CacheDir, cacheDuration, timeout, maxEntries, options)
			resultsCh <- feed
			summaryCh <- summary
		}(activeFeeds[i])
	}

	wg.Wait()
	close(resultsCh)
	close(summaryCh)

	// Collect results
	feeds := make([]*models.ExternalFeed, 0, len(activeFeeds))
	var allEntries []*models.ExternalEntry
	summary := blogrollRefreshSummary{}
	for feed := range resultsCh {
		feeds = append(feeds, feed)
		allEntries = append(allEntries, feed.Entries...)
	}
	for result := range summaryCh {
		summary.refreshed += result.refreshed
		summary.stale += result.stale
		summary.failed += result.failed
	}

	// Sort feeds by title
	sort.Slice(feeds, func(i, j int) bool {
		return strings.ToLower(feeds[i].Title) < strings.ToLower(feeds[j].Title)
	})

	// Sort entries deterministically: date desc, then feed URL, entry ID, title
	sort.SliceStable(allEntries, func(i, j int) bool {
		return compareEntries(allEntries[i], allEntries[j])
	})

	if config.CacheDir != "" {
		p.saveAggregateCache(config.CacheDir, p.aggregateConfigHash(activeFeeds), feeds, allEntries)
	}

	return feeds, allEntries, summary, nil
}

func (p *BlogrollPlugin) aggregateConfigHash(feeds []models.ExternalFeedConfig) string {
	data, err := json.Marshal(feeds)
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(append([]byte(blogrollCacheVersion+"|aggregate|"), data...))
	return hex.EncodeToString(hash[:])
}

func (p *BlogrollPlugin) loadAggregateCache(cacheDir string, maxAge time.Duration, configHash string) *blogrollAggregateCache {
	cached := p.loadAggregateCacheAny(cacheDir, configHash)
	if cached == nil {
		return nil
	}

	cacheFile := filepath.Join(cacheDir, aggregateCacheFile)
	info, err := os.Stat(cacheFile)
	if err != nil || time.Since(info.ModTime()) > maxAge {
		return nil
	}

	return cached
}

func (p *BlogrollPlugin) loadAggregateCacheAny(cacheDir, configHash string) *blogrollAggregateCache {
	cacheFile := filepath.Join(cacheDir, aggregateCacheFile)
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}
	var cached blogrollAggregateCache
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil
	}
	if cached.ConfigHash != configHash {
		return nil
	}
	for _, feed := range cached.Feeds {
		if feed != nil {
			normalizeExternalFeed(feed)
		}
	}
	return &cached
}

func (p *BlogrollPlugin) saveAggregateCache(cacheDir, configHash string, feeds []*models.ExternalFeed, entries []*models.ExternalEntry) {
	cacheFile := filepath.Join(cacheDir, aggregateCacheFile)
	payload := blogrollAggregateCache{
		ConfigHash: configHash,
		Feeds:      feeds,
		Entries:    entries,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	if err := os.WriteFile(cacheFile, data, 0o600); err != nil {
		return
	}
}

// initFeedFromConfig creates a new ExternalFeed initialized from configuration.
func initFeedFromConfig(config models.ExternalFeedConfig) *models.ExternalFeed {
	feed := &models.ExternalFeed{
		Config:   config,
		FeedURL:  config.URL,
		Category: config.Category,
		Type:     config.GetReaderType(),
		Tags:     config.Tags,
		Entries:  make([]*models.ExternalEntry, 0),
	}

	// Use config values if set
	if config.Title != "" {
		feed.Title = config.Title
	}
	if config.Description != "" {
		feed.Description = config.Description
	}
	if config.SiteURL != "" {
		feed.SiteURL = normalizeExternalURL(config.SiteURL)
	}
	if config.ImageURL != "" {
		feed.ImageURL = normalizeExternalURL(config.ImageURL)
	}
	normalizeExternalFeed(feed)

	return feed
}

// mergeCachedFeed merges config values into a cached feed.
func mergeCachedFeed(cached *models.ExternalFeed, config models.ExternalFeedConfig) *models.ExternalFeed {
	normalizeExternalFeed(cached)
	cached.Config = config
	cached.Type = config.GetReaderType()
	if config.Title != "" {
		cached.Title = config.Title
	}
	if config.Description != "" {
		cached.Description = config.Description
	}
	if config.Category != "" {
		cached.Category = config.Category
	}
	if config.SiteURL != "" {
		cached.SiteURL = normalizeExternalURL(config.SiteURL)
	}
	if config.ImageURL != "" {
		cached.ImageURL = normalizeExternalURL(config.ImageURL)
	}
	cached.Tags = config.Tags
	normalizeExternalFeed(cached)
	return cached
}

// updateFeedFromParsed updates feed metadata from parsed feed data.
func updateFeedFromParsed(feed *models.ExternalFeed, parsed *blogrollParsedFeed) {
	if feed.Title == "" {
		feed.Title = parsed.Title
	}
	if feed.Description == "" {
		feed.Description = parsed.Description
	}
	if feed.SiteURL == "" {
		feed.SiteURL = normalizeExternalURL(parsed.SiteURL)
	}
	if feed.ImageURL == "" {
		feed.ImageURL = normalizeExternalURL(parsed.ImageURL)
	}
	if feed.AvatarURL == "" && parsed.AvatarURL != "" {
		feed.AvatarURL = normalizeExternalURL(parsed.AvatarURL)
		feed.AvatarSource = parsed.AvatarSource
	}
	normalizeExternalFeed(feed)

	now := time.Now()
	feed.LastFetched = &now
	if parsed.LastUpdated != nil {
		feed.LastUpdated = parsed.LastUpdated
	}
}

// fetchFeed fetches a single feed with caching.
func (p *BlogrollPlugin) fetchFeed(config models.ExternalFeedConfig, cacheDir string, cacheDuration time.Duration, timeout, maxEntries int, options blogrollFetchOptions) (*models.ExternalFeed, blogrollRefreshSummary) {
	feed := initFeedFromConfig(config)
	var staleCached *models.ExternalFeed

	// Check cache
	if cacheDir != "" {
		if !options.forceRefresh {
			if cached := p.loadFromCache(config.URL, cacheDir, cacheDuration); cached != nil {
				cachedFeed := mergeCachedFeed(cached, config)
				p.emitFetchProgress(options, config, cachedFeed, "cached", "")
				return cachedFeed, blogrollRefreshSummary{}
			}
		}
		staleCached = p.loadFromCacheAny(config.URL, cacheDir)
		if staleCached != nil && !options.forceRefresh && !options.refreshOnBuild {
			cachedFeed := mergeCachedFeed(staleCached, config)
			p.emitFetchProgress(options, config, cachedFeed, "stale", "")
			return cachedFeed, blogrollRefreshSummary{stale: 1}
		}
	}

	if !options.forceRefresh && !options.refreshOnBuild {
		feed.Error = "blogroll refresh disabled and no cached feed available"
		p.emitFetchProgress(options, config, feed, "missing_cache", feed.Error)
		return feed, blogrollRefreshSummary{failed: 1}
	}

	// Fetch the feed
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", config.URL, http.NoBody)
	if err != nil {
		feed.Error = fmt.Sprintf("create request: %v", err)
		failedFeed, summary := p.handleFeedFetchFailure(feed, staleCached, config, cacheDir)
		p.emitFetchFailureProgress(options, config, failedFeed, summary)
		return failedFeed, summary
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (RSS Reader)")
	req.Header.Set("Accept", "application/rss+xml, application/atom+xml, application/xml, text/xml")

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	buildstats.InstrumentHTTPClient(client)

	resp, err := client.Do(req)
	if err != nil {
		feed.Error = fmt.Sprintf("fetch: %v", err)
		failedFeed, summary := p.handleFeedFetchFailure(feed, staleCached, config, cacheDir)
		p.emitFetchFailureProgress(options, config, failedFeed, summary)
		return failedFeed, summary
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		feed.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
		failedFeed, summary := p.handleFeedFetchFailure(feed, staleCached, config, cacheDir)
		p.emitFetchFailureProgress(options, config, failedFeed, summary)
		return failedFeed, summary
	}

	// Parse the feed using simple XML parsing
	parsedFeed, entries, err := parseBlogrollFeedResponse(resp)
	if err != nil {
		feed.Error = fmt.Sprintf("parse: %v", err)
		failedFeed, summary := p.handleFeedFetchFailure(feed, staleCached, config, cacheDir)
		p.emitFetchFailureProgress(options, config, failedFeed, summary)
		return failedFeed, summary
	}

	// Update feed with parsed values
	updateFeedFromParsed(feed, parsedFeed)

	// Attempt avatar discovery if no avatar set from config
	// Best-effort: failures don't affect the feed
	if feed.AvatarURL == "" && feed.SiteURL != "" {
		p.discoverFeedAvatar(ctx, feed, config, timeout)
	}

	// Limit entries
	if len(entries) > maxEntries {
		entries = entries[:maxEntries]
	}

	// Add feed info to entries
	for _, entry := range entries {
		entry.FeedURL = config.URL
		entry.FeedTitle = feed.Title
	}

	feed.Entries = entries
	feed.EntryCount = len(entries)

	// Save to cache
	if cacheDir != "" {
		p.saveToCache(feed, cacheDir)
	}

	p.emitFetchProgress(options, config, feed, "refreshed", "")
	return feed, blogrollRefreshSummary{refreshed: 1}
}

func (p *BlogrollPlugin) handleFeedFetchFailure(feed, staleCached *models.ExternalFeed, config models.ExternalFeedConfig, cacheDir string) (*models.ExternalFeed, blogrollRefreshSummary) {
	now := time.Now()
	if feed != nil {
		feed.LastFetched = &now
	}
	if staleCached != nil {
		return mergeCachedFeed(staleCached, config), blogrollRefreshSummary{stale: 1}
	}
	if cacheDir != "" && feed != nil {
		p.saveToCache(feed, cacheDir)
	}
	return feed, blogrollRefreshSummary{failed: 1}
}

func (p *BlogrollPlugin) emitFetchProgress(options blogrollFetchOptions, config models.ExternalFeedConfig, feed *models.ExternalFeed, status, errorMessage string) {
	if options.onFeedComplete == nil {
		return
	}
	title := config.Title
	if feed != nil && feed.Title != "" {
		title = feed.Title
	}
	if title == "" {
		title = config.URL
	}
	options.onFeedComplete(BlogrollCacheRefreshProgress{
		FeedURL:   config.URL,
		FeedTitle: title,
		Status:    status,
		Error:     errorMessage,
	})
}

func (p *BlogrollPlugin) emitFetchFailureProgress(options blogrollFetchOptions, config models.ExternalFeedConfig, feed *models.ExternalFeed, summary blogrollRefreshSummary) {
	status := "failed"
	if summary.stale > 0 {
		status = "stale"
	}
	errorMessage := ""
	if feed != nil {
		errorMessage = feed.Error
	}
	p.emitFetchProgress(options, config, feed, status, errorMessage)
}

func compareEntries(a, b *models.ExternalEntry) bool {
	// Date desc: prefer Published, then Updated
	ai := entryDate(a)
	bj := entryDate(b)
	if ai == nil && bj != nil {
		return false
	}
	if ai != nil && bj == nil {
		return true
	}
	if ai != nil && bj != nil {
		if ai.After(*bj) {
			return true
		}
		if bj.After(*ai) {
			return false
		}
	}

	// Tie-breakers for deterministic ordering
	if a.FeedURL != b.FeedURL {
		return a.FeedURL < b.FeedURL
	}
	if a.ID != b.ID {
		return a.ID < b.ID
	}
	return a.Title < b.Title
}

func entryDate(entry *models.ExternalEntry) *time.Time {
	if entry == nil {
		return nil
	}
	if entry.Published != nil {
		return entry.Published
	}
	if entry.Updated != nil {
		return entry.Updated
	}
	return nil
}

// loadFromCache loads a feed from cache if valid.
func (p *BlogrollPlugin) loadFromCache(url, cacheDir string, maxAge time.Duration) *models.ExternalFeed {
	feed := p.loadFromCacheAny(url, cacheDir)
	if feed == nil {
		return nil
	}

	fetchedAt := time.Time{}
	if feed.LastFetched != nil {
		fetchedAt = *feed.LastFetched
	}
	if fetchedAt.IsZero() {
		cacheFile := filepath.Join(cacheDir, p.cacheKey(url)+".json")
		if info, err := os.Stat(cacheFile); err == nil {
			fetchedAt = info.ModTime()
		}
	}
	if fetchedAt.IsZero() || time.Since(fetchedAt) > maxAge {
		return nil
	}

	return feed
}

func (p *BlogrollPlugin) loadFromCacheAny(url, cacheDir string) *models.ExternalFeed {
	cacheFile := filepath.Join(cacheDir, p.cacheKey(url)+".json")

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil
	}

	var feed models.ExternalFeed
	if err := json.Unmarshal(data, &feed); err != nil {
		return nil
	}
	normalizeExternalFeed(&feed)

	return &feed
}

// saveToCache saves a feed to cache.
// Cache writes are best-effort; errors are silently ignored.
func (p *BlogrollPlugin) saveToCache(feed *models.ExternalFeed, cacheDir string) {
	cacheFile := filepath.Join(cacheDir, p.cacheKey(feed.FeedURL)+".json")

	data, err := json.MarshalIndent(feed, "", "  ")
	if err != nil {
		return
	}

	//nolint:errcheck // Cache writes are best-effort
	os.WriteFile(cacheFile, data, 0o600)
}

// cacheKey generates a cache key from a URL.
func (p *BlogrollPlugin) cacheKey(url string) string {
	hash := sha256.Sum256([]byte(blogrollCacheVersion + "|" + url))
	return hex.EncodeToString(hash[:8])
}

// groupByCategory groups feeds by their category.
func (p *BlogrollPlugin) groupByCategory(feeds []*models.ExternalFeed) []*models.BlogrollCategory {
	categoryMap := make(map[string]*models.BlogrollCategory)

	for _, feed := range feeds {
		category := feed.Category
		if category == "" {
			category = categoryUncategorized
		}

		cat, ok := categoryMap[category]
		if !ok {
			cat = &models.BlogrollCategory{
				Name: category,
				Slug: blogrollSlugify(category),
			}
			categoryMap[category] = cat
		}
		cat.Feeds = append(cat.Feeds, feed)
	}

	// Convert to slice and sort
	categories := make([]*models.BlogrollCategory, 0, len(categoryMap))
	for _, cat := range categoryMap {
		categories = append(categories, cat)
	}

	sort.Slice(categories, func(i, j int) bool {
		// Put "Uncategorized" last
		if categories[i].Name == categoryUncategorized {
			return false
		}
		if categories[j].Name == categoryUncategorized {
			return true
		}
		return categories[i].Name < categories[j].Name
	})

	return categories
}

// writeBlogrollPage generates the blogroll page at the configured slug path.
func (p *BlogrollPlugin) writeBlogrollPage(m *lifecycle.Manager, outputDir string, feeds []*models.ExternalFeed, config models.BlogrollConfig) error {
	// Use configured slug or default to "blogroll"
	slug := config.BlogrollSlug
	if slug == "" {
		slug = defaultBlogrollSlug
	}

	// Create output directory
	blogrollDir := filepath.Join(outputDir, slug)
	if err := os.MkdirAll(blogrollDir, 0o755); err != nil {
		return err
	}

	// Group feeds by category
	categories := p.groupByCategory(feeds)

	// Get reader slug for cross-linking
	readerSlug := config.ReaderSlug
	if readerSlug == "" {
		readerSlug = defaultReaderSlug
	}

	// Build template context with config for theme inheritance
	ctx := map[string]interface{}{
		"title":        "Blogroll",
		"description":  "Blogs and feeds I follow",
		"feeds":        p.feedsToMaps(feeds),
		"categories":   p.categoriesToMaps(categories),
		"feed_count":   len(feeds),
		"config":       p.configToMap(m.Config()),
		"blogroll_url": "/" + slug + "/",
		"reader_url":   "/" + readerSlug + "/",
	}

	// Try to render with template engine
	content, err := p.renderTemplate(m, config.Templates.Blogroll, ctx)
	if err != nil {
		// Fall back to built-in template
		content = p.renderBlogrollFallback(feeds, categories, config)
	}

	// Write the file
	outputFile := filepath.Join(blogrollDir, "index.html")
	return os.WriteFile(outputFile, []byte(content), 0o644) //nolint:gosec // G306: Public-facing HTML needs 644 permissions
}

// writeReaderPage generates the paginated /reader pages.
// Pages are written concurrently for better I/O throughput.
func (p *BlogrollPlugin) writeReaderPage(m *lifecycle.Manager, outputDir string, feeds []*models.ExternalFeed, entries []*models.ExternalEntry, config models.BlogrollConfig) error {
	feedIndex := buildFeedIndex(feeds)
	variants := buildReaderVariants(feeds, entries, config, feedIndex, outputDir)
	for _, variant := range variants {
		variantEntries := filterEntriesByReaderVariant(entries, feedIndex, variant)
		if err := p.writeReaderVariantPages(m, feedIndex, config, variant, variants, variantEntries); err != nil {
			return err
		}
	}

	return nil
}

func buildReaderVariants(feeds []*models.ExternalFeed, entries []*models.ExternalEntry, config models.BlogrollConfig, feedIndex map[string]*models.ExternalFeed, outputDir string) []readerVariant {
	readerSlug := config.ReaderSlug
	if readerSlug == "" {
		readerSlug = defaultReaderSlug
	}

	variants := []readerVariant{{
		Key:         "all",
		Title:       "Reader",
		Description: "Latest posts from blogs I follow",
		BaseURL:     "/" + readerSlug,
		OutputDir:   filepath.Join(outputDir, readerSlug),
		EntryCount:  len(entries),
	}}

	present := make(map[models.ReaderFeedType]bool)
	for _, feed := range feeds {
		if feed == nil {
			continue
		}
		present[effectiveReaderFeedType(feed)] = true
	}

	for _, feedType := range readerFeedTypeOrder {
		if !present[feedType] {
			continue
		}
		variant := readerVariant{
			Key:         string(feedType),
			Title:       "Reader: " + readerFeedTypeLabel(feedType),
			Description: "Latest " + string(feedType) + " posts from blogs I follow",
			BaseURL:     "/" + readerSlug + "/" + string(feedType),
			OutputDir:   filepath.Join(outputDir, readerSlug, string(feedType)),
			FeedType:    feedType,
		}
		variant.EntryCount = len(filterEntriesByReaderVariant(entries, feedIndex, variant))
		variants = append(variants, variant)
	}

	return variants
}

func readerFeedTypeLabel(feedType models.ReaderFeedType) string {
	value := string(feedType)
	if value == "" {
		return "Reader"
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func effectiveReaderFeedType(feed *models.ExternalFeed) models.ReaderFeedType {
	if feed == nil {
		return models.ReaderFeedTypeWritten
	}
	if feed.Type != "" {
		return feed.Type
	}
	return feed.Config.GetReaderType()
}

func filterEntriesByReaderVariant(entries []*models.ExternalEntry, feedIndex map[string]*models.ExternalFeed, variant readerVariant) []*models.ExternalEntry {
	if variant.FeedType == "" {
		return entries
	}
	filtered := make([]*models.ExternalEntry, 0, len(entries))
	for _, entry := range entries {
		if entry == nil {
			continue
		}
		feed := feedIndex[entry.FeedURL]
		if effectiveReaderFeedType(feed) == variant.FeedType {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (p *BlogrollPlugin) writeReaderVariantPages(m *lifecycle.Manager, feedIndex map[string]*models.ExternalFeed, config models.BlogrollConfig, variant readerVariant, variants []readerVariant, entries []*models.ExternalEntry) error {
	if err := os.MkdirAll(variant.OutputDir, 0o755); err != nil {
		return err
	}

	pages := p.paginateEntries(entries, config, variant.BaseURL)
	readerTypes := readerVariantsToMaps(variants, variant.Key)
	if len(pages) == 0 {
		return p.writeReaderPageFile(m, variant.OutputDir, feedIndex, config, models.ReaderPage{
			Number:         1,
			Entries:        []*models.ExternalEntry{},
			TotalPages:     1,
			TotalItems:     0,
			ItemsPerPage:   config.ItemsPerPage,
			PaginationType: config.PaginationType,
			PageURLs:       []string{variant.BaseURL + "/"},
		}, true, variant, readerTypes)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(pages))
	for i := range pages {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			if err := p.writeReaderPageFile(m, variant.OutputDir, feedIndex, config, pages[idx], idx == 0, variant, readerTypes); err != nil {
				errCh <- err
			}
		}(i)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		return err
	}
	return nil
}

func readerVariantsToMaps(variants []readerVariant, activeKey string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(variants))
	for _, variant := range variants {
		result = append(result, map[string]interface{}{
			"key":         variant.Key,
			"title":       variant.Title,
			"description": variant.Description,
			"url":         variant.BaseURL + "/",
			"active":      variant.Key == activeKey,
			"entry_count": variant.EntryCount,
		})
	}
	return result
}

// writeReaderPageFile writes a single reader page (full and partial for HTMX).
func (p *BlogrollPlugin) writeReaderPageFile(m *lifecycle.Manager, readerDir string, feedIndex map[string]*models.ExternalFeed, config models.BlogrollConfig, page models.ReaderPage, isFirstPage bool, variant readerVariant, readerTypes []map[string]interface{}) error {
	// Get slugs for cross-linking
	blogrollSlug := config.BlogrollSlug
	if blogrollSlug == "" {
		blogrollSlug = defaultBlogrollSlug
	}

	// Build template context with config for theme inheritance
	updated := latestEntryDate(page.Entries)
	ctx := map[string]interface{}{
		"title":           variant.Title,
		"description":     variant.Description,
		"entries":         p.entriesToMaps(page.Entries, feedIndex, config),
		"day_groups":      p.readerDayGroups(page.Entries, feedIndex, config),
		"entry_count":     page.TotalItems,
		"config":          p.configToMap(m.Config()),
		"page":            p.readerPageToMap(page),
		"pagination_type": string(page.PaginationType),
		"blogroll_url":    "/" + blogrollSlug + "/",
		"reader_url":      variant.BaseURL + "/",
		"reader_type":     variant.Key,
		"reader_types":    readerTypes,
		"updated":         updated,
	}

	// Determine output path
	var outputFile string
	var partialDir string
	if isFirstPage {
		outputFile = filepath.Join(readerDir, "index.html")
		partialDir = filepath.Join(readerDir, "partial")
	} else {
		pageDir := filepath.Join(readerDir, "page", fmt.Sprintf("%d", page.Number))
		if err := os.MkdirAll(pageDir, 0o755); err != nil {
			return err
		}
		outputFile = filepath.Join(pageDir, "index.html")
		partialDir = filepath.Join(pageDir, "partial")
	}

	// Try to render with template engine
	content, err := p.renderTemplate(m, config.Templates.Reader, ctx)
	if err != nil {
		// Fall back to built-in template
		content = p.renderReaderFallback(page.Entries, feedIndex, page, config, variant, readerTypes)
	}

	// Write the full page
	if err := os.WriteFile(outputFile, []byte(content), 0o644); err != nil { //nolint:gosec // G306: Public-facing HTML needs 644 permissions
		return err
	}

	// For HTMX pagination, also generate partial pages
	if page.PaginationType == models.PaginationHTMX {
		if err := os.MkdirAll(partialDir, 0o755); err != nil {
			return err
		}

		partialContent := p.renderReaderPartial(page.Entries, feedIndex, page, config, variant, readerTypes)
		partialFile := filepath.Join(partialDir, "index.html")
		if err := os.WriteFile(partialFile, []byte(partialContent), 0o644); err != nil { //nolint:gosec // G306: Public-facing HTML needs 644 permissions
			return err
		}
	}

	return nil
}

// paginateEntries divides entries into pages based on ItemsPerPage and OrphanThreshold.
func (p *BlogrollPlugin) paginateEntries(entries []*models.ExternalEntry, config models.BlogrollConfig, baseURL string) []models.ReaderPage {
	if len(entries) == 0 {
		return []models.ReaderPage{}
	}

	itemsPerPage := config.ItemsPerPage
	if itemsPerPage <= 0 {
		itemsPerPage = 50
	}

	orphanThreshold := config.OrphanThreshold
	if orphanThreshold <= 0 {
		orphanThreshold = 3
	}

	paginationType := config.PaginationType
	if paginationType == "" {
		paginationType = models.PaginationManual
	}

	totalEntries := len(entries)
	var pages []models.ReaderPage

	for i := 0; i < totalEntries; i += itemsPerPage {
		end := i + itemsPerPage
		if end > totalEntries {
			end = totalEntries
		}

		// Check orphan threshold: if remaining items are below threshold,
		// add them to the current page instead of creating a new page
		remaining := totalEntries - end
		if remaining > 0 && remaining < orphanThreshold {
			end = totalEntries
		}

		pageNum := len(pages) + 1
		page := models.ReaderPage{
			Number:  pageNum,
			Entries: entries[i:end],
			HasPrev: pageNum > 1,
		}

		pages = append(pages, page)

		if end >= totalEntries {
			break
		}
	}

	totalPages := len(pages)

	// Generate page URLs for numbered navigation
	pageURLs := make([]string, totalPages)
	for i := 0; i < totalPages; i++ {
		if i == 0 {
			pageURLs[i] = baseURL + "/"
		} else {
			pageURLs[i] = baseURL + "/page/" + fmt.Sprintf("%d", i+1) + "/"
		}
	}

	// Set HasNext, URLs, and metadata for each page
	for i := range pages {
		pages[i].HasNext = i < totalPages-1
		pages[i].TotalPages = totalPages
		pages[i].TotalItems = totalEntries
		pages[i].ItemsPerPage = itemsPerPage
		pages[i].PageURLs = pageURLs
		pages[i].PaginationType = paginationType

		if pages[i].HasPrev {
			if i == 1 {
				pages[i].PrevURL = baseURL + "/"
			} else {
				pages[i].PrevURL = baseURL + "/page/" + fmt.Sprintf("%d", i) + "/"
			}
		}

		if pages[i].HasNext {
			pages[i].NextURL = baseURL + "/page/" + fmt.Sprintf("%d", i+2) + "/"
		}
	}

	return pages
}

// readerPageToMap converts a ReaderPage to a template-friendly map.
func (p *BlogrollPlugin) readerPageToMap(page models.ReaderPage) map[string]interface{} {
	return map[string]interface{}{
		"number":          page.Number,
		"has_prev":        page.HasPrev,
		"has_next":        page.HasNext,
		"prev_url":        page.PrevURL,
		"next_url":        page.NextURL,
		"total_pages":     page.TotalPages,
		"total_items":     page.TotalItems,
		"items_per_page":  page.ItemsPerPage,
		"page_urls":       page.PageURLs,
		"pagination_type": string(page.PaginationType),
	}
}

// readerDayGroups groups entries by day in reverse chronological order.
func (p *BlogrollPlugin) readerDayGroups(entries []*models.ExternalEntry, feedIndex map[string]*models.ExternalFeed, config models.BlogrollConfig) []map[string]interface{} {
	if len(entries) == 0 {
		return nil
	}

	groups := make([]map[string]interface{}, 0)
	for _, entry := range entries {
		if entry == nil {
			continue
		}
		entryMap := p.readerEntryMap(entry, feedIndex[entry.FeedURL], config)
		dateKey, dateISO, weekday, title, year := readerDateFields(entryDate(entry))

		if len(groups) == 0 || stringValue(groups[len(groups)-1], "date_key") != dateKey {
			groups = append(groups, map[string]interface{}{
				"date_key": dateKey,
				"date_iso": dateISO,
				"weekday":  weekday,
				"title":    title,
				"year":     year,
				"count":    0,
				"entries":  []map[string]interface{}{},
			})
		}

		last := groups[len(groups)-1]
		lastEntries := entriesValue(last)
		lastEntries = append(lastEntries, entryMap)
		last["entries"] = lastEntries
		last["count"] = groupCount(last) + 1
	}

	return groups
}

func readerDateFields(date *time.Time) (dateKey, dateISO, weekday, title, year string) {
	if date == nil {
		return "undated", "", "Undated", "No date", ""
	}

	dateKey = date.Format("2006-01-02")
	dateISO = date.Format(time.RFC3339)
	weekday = date.Format("Mon")
	title = date.Format("Jan 2")
	year = date.Format("2006")
	return dateKey, dateISO, weekday, title, year
}

func (p *BlogrollPlugin) renderReaderTimeline(groups []map[string]interface{}, wrap bool) string {
	var sb strings.Builder
	if wrap {
		sb.WriteString(`    <div class="reader-stream posts-list">
`)
	} else {
		sb.WriteString(`<div class="reader-stream posts-list">
`)
	}

	if len(groups) == 0 {
		sb.WriteString(`      <p class="reader-empty">No reader entries were found.</p>
`)
	} else {
		for _, group := range groups {
			sb.WriteString(`      <section class="reader-day">
        <aside class="reader-day-rail">
          <time`)
			if dateISO := stringValue(group, "date_iso"); dateISO != "" {
				sb.WriteString(` datetime="`)
				sb.WriteString(html.EscapeString(dateISO))
				sb.WriteString(`"`)
			}
			sb.WriteString(`>
            <span class="reader-day-weekday">`)
			sb.WriteString(html.EscapeString(stringValue(group, "weekday")))
			sb.WriteString(`</span>
            <span class="reader-day-title">`)
			sb.WriteString(html.EscapeString(stringValue(group, "title")))
			sb.WriteString(`</span>`)
			if year := stringValue(group, "year"); year != "" {
				sb.WriteString(`
            <span class="reader-day-year">`)
				sb.WriteString(html.EscapeString(year))
				sb.WriteString(`</span>`)
			}
			sb.WriteString(`
          </time>
          <span class="reader-day-count">`)
			count := groupCount(group)
			label := "stories"
			if count == 1 {
				label = "story"
			}
			sb.WriteString(fmt.Sprintf("%d %s", count, label))
			sb.WriteString(`</span>
        </aside>
        <div class="reader-day-entries">
`)

			entries := entriesValue(group)
			for _, entry := range entries {
				previewKind := stringValue(entry, "preview_kind")
				sb.WriteString(`          <article class="reader-entry`)
				if imageURL := stringValue(entry, "image_url"); imageURL != "" {
					sb.WriteString(` has-image`)
				}
				if previewKind != "" {
					sb.WriteString(` preview-`)
					sb.WriteString(html.EscapeString(previewKind))
				}
				sb.WriteString(`">
            <div class="reader-entry-meta-row">
              <a href="`)
				sb.WriteString(html.EscapeString(stringValue(entry, "feed_url")))
				sb.WriteString(`" target="_blank" rel="noopener noreferrer" class="reader-entry-source-link">
`)
				if iconURL := stringValue(entry, "source_icon_url"); iconURL != "" {
					sb.WriteString(`                <img src="`)
					sb.WriteString(html.EscapeString(iconURL))
					sb.WriteString(`" alt="" class="reader-entry-source-icon" width="16" height="16" loading="lazy" decoding="async">
`)
				} else {
					sb.WriteString(`                <span class="reader-entry-source-icon reader-entry-source-icon--fallback" aria-hidden="true"></span>
`)
				}
				sb.WriteString(`                <span class="reader-entry-source">`)
				sb.WriteString(html.EscapeString(stringValue(entry, "feed_title")))
				sb.WriteString(`</span>
              </a>`)

				if publishedDatetime := stringValue(entry, "published_datetime"); publishedDatetime != "" {
					sb.WriteString(`
              <time datetime="`)
					sb.WriteString(html.EscapeString(publishedDatetime))
					sb.WriteString(`">`)
					sb.WriteString(html.EscapeString(stringValue(entry, "published_label")))
					sb.WriteString(`</time>`)
				}

				sb.WriteString(`
            </div>
`)
				if originalURL := stringValue(entry, "original_url"); originalURL != "" {
					if discussionSource := stringValue(entry, "discussion_source"); discussionSource != "" {
						sb.WriteString(`
            <div class="reader-entry-found-on">Found on `)
						sb.WriteString(html.EscapeString(discussionSource))
						sb.WriteString(`: <a href="`)
						sb.WriteString(html.EscapeString(originalURL))
						sb.WriteString(`" target="_blank" rel="noopener noreferrer">discussion</a></div>
`)
					}
				}
				sb.WriteString(`            <h2 class="reader-entry-title"><a href="`)
				sb.WriteString(html.EscapeString(stringValue(entry, "url")))
				sb.WriteString(`" target="_blank" rel="noopener noreferrer">`)
				sb.WriteString(html.EscapeString(stringValue(entry, "title")))
				sb.WriteString(`<span class="visually-hidden">(opens in new tab)</span></a></h2>
`)

				if description := stringValue(entry, "description"); description != "" {
					sb.WriteString(`            <p class="reader-entry-description">`)
					sb.WriteString(html.EscapeString(blogrollTruncateString(blogrollStripHTML(description), 200)))
					sb.WriteString(`</p>
`)
				}

				if imageURL := stringValue(entry, "image_url"); imageURL != "" {
					sb.WriteString(`            <a href="`)
					sb.WriteString(html.EscapeString(stringValue(entry, "url")))
					sb.WriteString(`" target="_blank" rel="noopener noreferrer" class="reader-entry-image-link reader-entry-image-link--`)
					sb.WriteString(html.EscapeString(previewKind))
					sb.WriteString(`" aria-label="Read `)
					sb.WriteString(html.EscapeString(stringValue(entry, "title")))
					sb.WriteString(` (opens in new tab)">
              <img src="`)
					sb.WriteString(html.EscapeString(imageURL))
					if srcset := stringValue(entry, "image_srcset"); srcset != "" {
						sb.WriteString(`" srcset="`)
						sb.WriteString(html.EscapeString(srcset))
						sb.WriteString(`" sizes="100vw`)
					}
					sb.WriteString(`" alt="" class="reader-entry-image reader-entry-image--`)
					sb.WriteString(html.EscapeString(previewKind))
					sb.WriteString(`" width="960" height="540" loading="lazy" decoding="async">
            </a>
`)
				} else {
					sb.WriteString(`            <a href="`)
					sb.WriteString(html.EscapeString(stringValue(entry, "url")))
					sb.WriteString(`" target="_blank" rel="noopener noreferrer" aria-label="Read `)
					sb.WriteString(html.EscapeString(stringValue(entry, "title")))
					sb.WriteString(` (opens in new tab)" class="reader-entry-image-link reader-entry-image-link--source-tile">
              <span class="reader-entry-source-tile" aria-hidden="true">
                <span class="reader-entry-source-tile-top">`)
					if iconURL := stringValue(entry, "source_icon_url"); iconURL != "" {
						sb.WriteString(`<img src="`)
						sb.WriteString(html.EscapeString(iconURL))
						sb.WriteString(`" alt="" class="reader-entry-source-tile-icon" width="20" height="20" loading="lazy" decoding="async">`)
					} else {
						sb.WriteString(`<span class="reader-entry-source-tile-badge">`)
						sb.WriteString(html.EscapeString(stringValue(entry, "source_initial")))
						sb.WriteString(`</span>`)
					}
					sb.WriteString(`<span class="reader-entry-source-tile-label">`)
					sb.WriteString(html.EscapeString(stringValue(entry, "feed_title")))
					sb.WriteString(`</span>
                </span>
                <span class="reader-entry-source-tile-domain">`)
					sb.WriteString(html.EscapeString(stringValue(entry, "source_domain")))
					sb.WriteString(`</span>
              </span>
            </a>
`)
				}

				sb.WriteString(`          </article>
`)
			}

			sb.WriteString(`        </div>
      </section>
`)
		}
	}

	sb.WriteString(`</div>
`)
	return sb.String()
}

func stringValue(m map[string]interface{}, key string) string {
	if value, ok := m[key].(string); ok {
		return value
	}
	return ""
}

func groupCount(group map[string]interface{}) int {
	if count, ok := group["count"].(int); ok {
		return count
	}
	return 0
}

func entriesValue(group map[string]interface{}) []map[string]interface{} {
	if entries, ok := group["entries"].([]map[string]interface{}); ok {
		return entries
	}
	return nil
}

func (p *BlogrollPlugin) readerEntryMap(entry *models.ExternalEntry, feed *models.ExternalFeed, config models.BlogrollConfig) map[string]interface{} {
	if entry == nil {
		return map[string]interface{}{}
	}

	sourceTitle := entry.FeedTitle
	if sourceTitle == "" && feed != nil {
		sourceTitle = feed.Title
	}
	sourceSiteURL := ""
	sourceImageURL := ""
	if feed != nil {
		sourceSiteURL = feed.SiteURL
		if feed.AvatarURL != "" {
			sourceImageURL = feed.AvatarURL
		} else if feed.ImageURL != "" {
			sourceImageURL = feed.ImageURL
		}
	}
	previewURL, previewKind := readerPreviewForEntry(entry, sourceImageURL, config.FallbackImageService)
	sourceDomain := readerDomain(entry.FeedURL)
	if sourceDomain == "" {
		sourceDomain = readerDomain(sourceSiteURL)
	}
	feedType := models.ReaderFeedTypeWritten
	if feed != nil {
		feedType = effectiveReaderFeedType(feed)
	}

	result := map[string]interface{}{
		"feed_url":          entry.FeedURL,
		"feed_title":        sourceTitle,
		"feed_type":         string(feedType),
		"id":                entry.ID,
		"url":               entry.URL,
		"original_url":      entry.OriginalURL,
		"discussion_source": aggregatorDiscussionLabel(entry.OriginalURL),
		"title":             entry.Title,
		"description":       entry.Description,
		"content":           entry.Content,
		"author":            entry.Author,
		"updated":           entry.Updated,
		"categories":        entry.Categories,
		"image_url":         previewURL,
		"image_srcset":      readerImageSrcset(previewURL),
		"preview_url":       previewURL,
		"preview_kind":      previewKind,
		"source_image_url":  sourceImageURL,
		"source_site_url":   sourceSiteURL,
		"reading_time":      entry.ReadingTime,
	}

	if date := entryDate(entry); date != nil {
		result["published"] = date
		result["published_datetime"] = date.Format(time.RFC3339)
		result["published_label"] = templates.FormatHumanDate(*date)
		result["date_key"] = date.Format("2006-01-02")
		result["date_iso"] = date.Format(time.RFC3339)
	}

	if sourceDomain != "" {
		result["source_domain"] = sourceDomain
		result["source_icon_url"] = "https://www.google.com/s2/favicons?domain=" + neturl.QueryEscape(sourceDomain) + "&sz=32"
	}
	if matches := readerMentionsForURL(p.urlMentions, entry.URL); len(matches) > 0 {
		result["url_mentions"] = matches
	}
	result["source_initial"] = readerSourceInitial(sourceTitle, sourceDomain)

	return result
}

func buildFeedIndex(feeds []*models.ExternalFeed) map[string]*models.ExternalFeed {
	if len(feeds) == 0 {
		return nil
	}

	index := make(map[string]*models.ExternalFeed, len(feeds))
	for _, feed := range feeds {
		if feed == nil || feed.FeedURL == "" {
			continue
		}
		index[feed.FeedURL] = feed
	}
	return index
}

func readerPreviewForEntry(entry *models.ExternalEntry, sourceImageURL, fallbackImageService string) (previewURL, previewKind string) {
	if entry != nil && entry.ImageURL != "" {
		return normalizeExternalURL(entry.ImageURL), "article-image"
	}
	if sourceImageURL != "" {
		return normalizeExternalURL(sourceImageURL), "source-image"
	}
	if fallbackImageService != "" && entry != nil && entry.URL != "" {
		return generateFallbackImageURL(fallbackImageService, entry.URL), "screenshot"
	}
	return "", "source-tile"
}

func readerSourceInitial(sourceTitle, sourceDomain string) string {
	text := strings.TrimSpace(sourceTitle)
	if text == "" {
		text = strings.TrimSpace(sourceDomain)
	}
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return strings.ToUpper(string(r))
		}
	}
	return "?"
}

func readerDomain(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	parsed, err := neturl.Parse(rawURL)
	if (err != nil || parsed.Hostname() == "") && !strings.Contains(rawURL, "://") {
		parsed, err = neturl.Parse("https://" + rawURL)
	}
	if err == nil && parsed.Hostname() != "" {
		return parsed.Hostname()
	}

	if strings.Contains(rawURL, ".") {
		return strings.TrimPrefix(rawURL, "https://")
	}

	return ""
}

// renderTemplate attempts to render using the template engine.
func (p *BlogrollPlugin) renderTemplate(m *lifecycle.Manager, templateName string, ctx map[string]interface{}) (string, error) {
	if previews, ok := m.Cache().Get("nav_previews"); ok {
		ctx["nav_previews"] = previews
	}
	// Check if template engine is available
	// The templates plugin stores it as "templates.engine"
	engine, ok := m.Cache().Get("templates.engine")
	if !ok {
		return "", fmt.Errorf("template engine not available")
	}

	// Try to use the engine
	if eng, ok := engine.(interface {
		RenderToString(string, map[string]interface{}) (string, error)
	}); ok {
		return eng.RenderToString(templateName, ctx)
	}

	return "", fmt.Errorf("template engine does not support RenderToString")
}

// feedsToMaps converts feeds to template-friendly maps.
func (p *BlogrollPlugin) feedsToMaps(feeds []*models.ExternalFeed) []map[string]interface{} {
	result := make([]map[string]interface{}, len(feeds))
	for i, feed := range feeds {
		sourceURL, sourceLabel := blogrollSourceLink(feed.SiteURL, feed.FeedURL)
		result[i] = map[string]interface{}{
			"title":         feed.Title,
			"description":   feed.Description,
			"site_url":      feed.SiteURL,
			"feed_url":      feed.FeedURL,
			"source_url":    sourceURL,
			"source_label":  sourceLabel,
			"image_url":     feed.ImageURL,
			"avatar_url":    feed.AvatarURL,
			"avatar_source": feed.AvatarSource,
			"category":      feed.Category,
			"tags":          feed.Tags,
			"entry_count":   feed.EntryCount,
			"last_fetched":  feed.LastFetched,
			"last_updated":  feed.LastUpdated,
			"error":         feed.Error,
		}
	}
	return result
}

func blogrollSourceLink(siteURL, feedURL string) (sourceURL, sourceLabel string) {
	if rootURL, host := blogrollRootURL(siteURL); rootURL != "" && host != "" {
		return rootURL, blogrollDisplayDomain(host)
	}

	if rootURL, host := blogrollRootURL(feedURL); rootURL != "" && host != "" {
		return rootURL, blogrollDisplayDomain(host)
	}

	return "", ""
}

func blogrollRootURL(rawURL string) (rootURL, host string) {
	if rawURL == "" {
		return "", ""
	}

	parsed, err := neturl.Parse(rawURL)
	if (err != nil || parsed.Hostname() == "") && !strings.Contains(rawURL, "://") {
		parsed, err = neturl.Parse("https://" + rawURL)
	}
	if err != nil || parsed.Hostname() == "" {
		return "", ""
	}

	scheme := parsed.Scheme
	if scheme == "" {
		scheme = schemeHTTPS
	}

	return scheme + "://" + parsed.Host, parsed.Hostname()
}

func blogrollDisplayDomain(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
	host = strings.TrimPrefix(host, "www.")
	return host
}

// entriesToMaps converts entries to template-friendly maps.
func (p *BlogrollPlugin) entriesToMaps(entries []*models.ExternalEntry, feedIndex map[string]*models.ExternalFeed, config models.BlogrollConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, len(entries))
	for i, entry := range entries {
		if entry == nil {
			result[i] = map[string]interface{}{}
			continue
		}
		result[i] = p.readerEntryMap(entry, feedIndex[entry.FeedURL], config)
	}
	return result
}

func latestEntryDate(entries []*models.ExternalEntry) *time.Time {
	var latest *time.Time
	for _, entry := range entries {
		if entry == nil {
			continue
		}
		candidate := entryDate(entry)
		if candidate == nil {
			continue
		}
		if latest == nil || candidate.After(*latest) {
			t := *candidate
			latest = &t
		}
	}
	return latest
}

// categoriesToMaps converts categories to template-friendly maps.
func (p *BlogrollPlugin) categoriesToMaps(categories []*models.BlogrollCategory) []map[string]interface{} {
	result := make([]map[string]interface{}, len(categories))
	for i, cat := range categories {
		result[i] = map[string]interface{}{
			"name":  cat.Name,
			"slug":  cat.Slug,
			"feeds": p.feedsToMaps(cat.Feeds),
		}
	}
	return result
}

// configToMap converts config to a template-friendly map with essential fields.
// Note: This should eventually be refactored to use ToModelsConfig + configToMap
// from templates package for full consistency.
func (p *BlogrollPlugin) configToMap(config *lifecycle.Config) map[string]interface{} {
	if config == nil {
		return nil
	}

	result := map[string]interface{}{
		"output_dir": config.OutputDir,
		"lang":       "en", // Default language
	}

	if config.Extra != nil {
		p.extractStringFields(config.Extra, result)
		p.extractNavItems(config.Extra, result)
		p.extractSearchConfig(config.Extra, result)
		p.extractStructConfigs(config.Extra, result)
		p.extractLayoutConfig(config.Extra, result)
		p.extractSEOConfig(config.Extra, result)
		p.extractMiscConfigs(config.Extra, result)
	}

	return result
}

// extractStringFields extracts common string fields from config.Extra.
func (p *BlogrollPlugin) extractStringFields(extra, result map[string]interface{}) {
	stringFields := []string{"title", "description", "url", "author", "templates_dir"}
	for _, field := range stringFields {
		if value, ok := extra[field].(string); ok {
			result[field] = value
		}
	}
}

// extractNavItems converts nav items from config.Extra.
func (p *BlogrollPlugin) extractNavItems(extra, result map[string]interface{}) {
	navItems, ok := extra["nav"].([]models.NavItem)
	if !ok {
		return
	}
	navMaps := make([]map[string]interface{}, len(navItems))
	for i, nav := range navItems {
		navMaps[i] = map[string]interface{}{
			"label":    nav.Label,
			"url":      nav.URL,
			"external": nav.External,
		}
	}
	result["nav"] = navMaps
}

// extractSearchConfig converts search config from config.Extra.
// Sets defaults (enabled=true, position="navbar") when no config is provided.
func (p *BlogrollPlugin) extractSearchConfig(extra, result map[string]interface{}) {
	// Start with defaults
	search := models.NewSearchConfig()

	// Override with config from Extra if available
	if sc, ok := extra["search"].(models.SearchConfig); ok {
		search = sc
	}

	// Use the same normalized search context as normal pages. In particular,
	// backend and bleve.endpoint must survive into reader.html's base template,
	// or Bleve-only sites will request Pagefind assets that were never built.
	result["search"] = templates.SearchToMap(&search)
}

// extractStructConfigs extracts simple struct configs from config.Extra.
// All structs must be converted to maps for pongo2 template access.
func (p *BlogrollPlugin) extractStructConfigs(extra, result map[string]interface{}) {
	// Convert components config if available (fixes #316)
	// Must convert to map for pongo2 template dot-notation access
	if components, ok := extra["components"].(models.ComponentsConfig); ok {
		result["components"] = p.componentsToMap(&components)
	}
	if footer, ok := extra["footer"].(models.FooterConfig); ok {
		result["footer"] = p.footerConfigToMap(&footer)
	}
	if sidebar, ok := extra["sidebar"].(models.SidebarConfig); ok {
		result["sidebar"] = p.sidebarToMap(&sidebar)
	}
	if toc, ok := extra["toc"].(models.TocConfig); ok {
		result["toc"] = p.tocToMap(&toc)
	}
	if header, ok := extra["header"].(models.HeaderLayoutConfig); ok {
		result["header"] = p.headerToMap(&header)
	}
	if postFormats, ok := extra["post_formats"].(models.PostFormatsConfig); ok {
		result["post_formats"] = p.postFormatsToMap(&postFormats)
	}
	if head, ok := extra["head"].(models.HeadConfig); ok {
		result["head"] = p.headToMap(&head)
	}
}

// componentsToMap converts ComponentsConfig to a template-friendly map.
// This is required because pongo2 cannot access Go struct fields directly.
func (p *BlogrollPlugin) componentsToMap(c *models.ComponentsConfig) map[string]interface{} {
	if c == nil {
		return nil
	}

	// Convert nav component
	navEnabled := true
	if c.Nav.Enabled != nil {
		navEnabled = *c.Nav.Enabled
	}
	navItems := make([]map[string]interface{}, len(c.Nav.Items))
	for i, item := range c.Nav.Items {
		navItems[i] = map[string]interface{}{
			"label":    item.Label,
			"url":      item.URL,
			"external": item.External,
		}
	}
	navMap := map[string]interface{}{
		"enabled":  navEnabled,
		"position": c.Nav.Position,
		"style":    c.Nav.Style,
		"items":    navItems,
	}

	// Convert footer component
	footerEnabled := true
	if c.Footer.Enabled != nil {
		footerEnabled = *c.Footer.Enabled
	}
	showCopyright := true
	if c.Footer.ShowCopyright != nil {
		showCopyright = *c.Footer.ShowCopyright
	}
	footerLinks := make([]map[string]interface{}, len(c.Footer.Links))
	for i, link := range c.Footer.Links {
		footerLinks[i] = map[string]interface{}{
			"label":    link.Label,
			"url":      link.URL,
			"external": link.External,
		}
	}
	footerMap := map[string]interface{}{
		"enabled":        footerEnabled,
		"text":           c.Footer.Text,
		"show_copyright": showCopyright,
		"links":          footerLinks,
	}

	// Convert doc_sidebar component
	docSidebarEnabled := false
	if c.DocSidebar.Enabled != nil {
		docSidebarEnabled = *c.DocSidebar.Enabled
	}
	docSidebarMap := map[string]interface{}{
		"enabled":   docSidebarEnabled,
		"position":  c.DocSidebar.Position,
		"width":     c.DocSidebar.Width,
		"min_depth": c.DocSidebar.MinDepth,
		"max_depth": c.DocSidebar.MaxDepth,
	}

	// Convert feed_sidebar component
	feedSidebarEnabled := false
	if c.FeedSidebar.Enabled != nil {
		feedSidebarEnabled = *c.FeedSidebar.Enabled
	}
	feedSidebarMap := map[string]interface{}{
		"enabled":  feedSidebarEnabled,
		"position": c.FeedSidebar.Position,
		"width":    c.FeedSidebar.Width,
		"title":    c.FeedSidebar.Title,
		"feeds":    c.FeedSidebar.Feeds,
	}

	return map[string]interface{}{
		"nav":          navMap,
		"footer":       footerMap,
		"doc_sidebar":  docSidebarMap,
		"feed_sidebar": feedSidebarMap,
	}
}

// footerConfigToMap converts FooterConfig to a template-friendly map.
func (p *BlogrollPlugin) footerConfigToMap(f *models.FooterConfig) map[string]interface{} {
	if f == nil {
		return nil
	}
	showCopyright := true
	if f.ShowCopyright != nil {
		showCopyright = *f.ShowCopyright
	}
	return map[string]interface{}{
		"text":           f.Text,
		"show_copyright": showCopyright,
	}
}

// sidebarToMap converts SidebarConfig to a template-friendly map.
func (p *BlogrollPlugin) sidebarToMap(s *models.SidebarConfig) map[string]interface{} {
	if s == nil {
		return nil
	}
	result := map[string]interface{}{
		"position": s.Position,
		"width":    s.Width,
		"title":    s.Title,
	}
	if s.Enabled != nil {
		result["enabled"] = *s.Enabled
	} else {
		result["enabled"] = true
	}
	if s.Collapsible != nil {
		result["collapsible"] = *s.Collapsible
	}
	if s.DefaultOpen != nil {
		result["default_open"] = *s.DefaultOpen
	}
	return result
}

// tocToMap converts TocConfig to a template-friendly map.
func (p *BlogrollPlugin) tocToMap(t *models.TocConfig) map[string]interface{} {
	if t == nil {
		return nil
	}
	result := map[string]interface{}{
		"position":  t.Position,
		"width":     t.Width,
		"min_depth": t.MinDepth,
		"max_depth": t.MaxDepth,
		"title":     t.Title,
	}
	if t.Enabled != nil {
		result["enabled"] = *t.Enabled
	} else {
		result["enabled"] = true
	}
	if t.Collapsible != nil {
		result["collapsible"] = *t.Collapsible
	}
	if t.DefaultOpen != nil {
		result["default_open"] = *t.DefaultOpen
	}
	return result
}

// headerToMap converts HeaderLayoutConfig to a template-friendly map.
func (p *BlogrollPlugin) headerToMap(h *models.HeaderLayoutConfig) map[string]interface{} {
	if h == nil {
		return nil
	}
	result := map[string]interface{}{
		"style": h.Style,
	}
	if h.Sticky != nil {
		result["sticky"] = *h.Sticky
	} else {
		result["sticky"] = true
	}
	if h.ShowLogo != nil {
		result["show_logo"] = *h.ShowLogo
	} else {
		result["show_logo"] = true
	}
	if h.ShowTitle != nil {
		result["show_title"] = *h.ShowTitle
	} else {
		result["show_title"] = true
	}
	if h.ShowNav != nil {
		result["show_nav"] = *h.ShowNav
	} else {
		result["show_nav"] = true
	}
	if h.ShowSearch != nil {
		result["show_search"] = *h.ShowSearch
	} else {
		result["show_search"] = true
	}
	if h.ShowThemeToggle != nil {
		result["show_theme_toggle"] = *h.ShowThemeToggle
	} else {
		result["show_theme_toggle"] = true
	}
	return result
}

// postFormatsToMap converts PostFormatsConfig to a template-friendly map.
func (p *BlogrollPlugin) postFormatsToMap(pf *models.PostFormatsConfig) map[string]interface{} {
	if pf == nil {
		return nil
	}
	htmlEnabled := true
	if pf.HTML != nil {
		htmlEnabled = *pf.HTML
	}
	return map[string]interface{}{
		"html":     htmlEnabled,
		"markdown": pf.Markdown,
		"text":     pf.Text,
		"og":       pf.OG,
	}
}

// headToMap converts HeadConfig to a template-friendly map.
func (p *BlogrollPlugin) headToMap(h *models.HeadConfig) map[string]interface{} {
	if h == nil {
		return nil
	}

	// Convert meta tags
	metaTags := make([]map[string]interface{}, len(h.Meta))
	for i, meta := range h.Meta {
		metaTags[i] = map[string]interface{}{
			"name":     meta.Name,
			"property": meta.Property,
			"content":  meta.Content,
		}
	}

	// Convert link tags
	linkTags := make([]map[string]interface{}, len(h.Link))
	for i, link := range h.Link {
		linkTags[i] = map[string]interface{}{
			"rel":         link.Rel,
			"href":        link.Href,
			"crossorigin": link.Crossorigin,
		}
	}

	// Convert script tags
	scriptTags := make([]map[string]interface{}, len(h.Script))
	for i, script := range h.Script {
		scriptTags[i] = map[string]interface{}{
			"src": script.Src,
		}
	}

	// Convert alternate feeds
	alternateFeeds := make([]map[string]interface{}, len(h.AlternateFeeds))
	for i, feed := range h.AlternateFeeds {
		alternateFeeds[i] = map[string]interface{}{
			"type":      feed.Type,
			"title":     feed.Title,
			"href":      feed.Href,
			"mime_type": feed.GetMIMEType(),
		}
	}

	return map[string]interface{}{
		"text":            h.Text,
		"meta":            metaTags,
		"link":            linkTags,
		"script":          scriptTags,
		"alternate_feeds": alternateFeeds,
	}
}

// extractLayoutConfig converts layout config from config.Extra, handling both pointer and value types.
func (p *BlogrollPlugin) extractLayoutConfig(extra, result map[string]interface{}) {
	switch layoutVal := extra["layout"].(type) {
	case *models.LayoutConfig:
		result["layout"] = *layoutVal
	case models.LayoutConfig:
		result["layout"] = layoutVal
	}
}

// extractSEOConfig converts SEO config from config.Extra, handling both struct and map types.
func (p *BlogrollPlugin) extractSEOConfig(extra, result map[string]interface{}) {
	switch seoVal := extra["seo"].(type) {
	case models.SEOConfig:
		result["seo"] = seoVal
	case map[string]interface{}:
		result["seo"] = models.SEOConfig{
			TwitterHandle: p.getStringFromMap(seoVal, "twitter_handle"),
			DefaultImage:  p.getStringFromMap(seoVal, "default_image"),
			LogoURL:       p.getStringFromMap(seoVal, "logo_url"),
			AuthorImage:   p.getStringFromMap(seoVal, "author_image"),
		}
	}
}

// extractMiscConfigs extracts miscellaneous configs that don't need type conversion.
func (p *BlogrollPlugin) extractMiscConfigs(extra, result map[string]interface{}) {
	if webmention, ok := extra["webmention"]; ok {
		result["webmention"] = webmention
	}
	if feeds, ok := extra["feeds"]; ok {
		result["feeds"] = feeds
	}
	// Extract theme config for background-css partial and other theme features
	if theme, ok := extra["theme"].(models.ThemeConfig); ok {
		result["theme"] = p.themeToMap(theme)
	}
}

// themeToMap converts ThemeConfig to a template-friendly map.
// This delegates to the shared templates.ThemeToMap to ensure consistency
// across all pages that use base.html.
func (p *BlogrollPlugin) themeToMap(theme models.ThemeConfig) map[string]interface{} {
	return templates.ThemeToMap(&theme)
}

// getStringFromMap safely gets a string value from a map.
func (p *BlogrollPlugin) getStringFromMap(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// renderBlogrollFallback generates a basic blogroll page that uses theme CSS if available.
func (p *BlogrollPlugin) renderBlogrollFallback(feeds []*models.ExternalFeed, categories []*models.BlogrollCategory, config models.BlogrollConfig) string {
	var sb strings.Builder

	// Get configured slugs with defaults
	blogrollSlug := config.BlogrollSlug
	if blogrollSlug == "" {
		blogrollSlug = defaultBlogrollSlug
	}
	readerSlug := config.ReaderSlug
	if readerSlug == "" {
		readerSlug = defaultReaderSlug
	}

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Blogroll</title>
  <!-- Theme CSS - uses site's configured palette if available -->
  <link rel="stylesheet" href="/css/variables.css">
  <link rel="stylesheet" href="/css/palette.css">
  <link rel="stylesheet" href="/css/main.css">
  <link rel="stylesheet" href="/css/components.css">
  <style>
    /* Fallback styles if theme CSS is not available */
    :root {
      --color-background: #ffffff;
      --color-text: #1a1a1a;
      --color-text-muted: #666666;
      --color-border: #e0e0e0;
      --color-surface: #f8f9fa;
      --color-primary: #3b82f6;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --color-background: #1a1a1a;
        --color-text: #f0f0f0;
        --color-text-muted: #999999;
        --color-border: #333333;
        --color-surface: #2a2a2a;
        --color-primary: #60a5fa;
      }
    }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      max-width: 900px;
      margin: 0 auto;
      padding: 2rem;
      background: var(--color-background);
      color: var(--color-text);
      line-height: 1.6;
    }
    h1 { margin-bottom: 0.5rem; }
    .subtitle { color: var(--color-text-muted); margin-bottom: 2rem; }
    .category { margin-bottom: 2rem; }
    .category h2 { border-bottom: 1px solid var(--color-border); padding-bottom: 0.5rem; }
    .feed-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
      gap: 1rem;
    }
    .feed-card {
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: 8px;
      padding: 1rem;
      display: flex;
      gap: 0.75rem;
    }
    .feed-avatar {
      width: 48px;
      height: 48px;
      border-radius: 6px;
      object-fit: cover;
      flex-shrink: 0;
      background: var(--color-border);
    }
    .feed-avatar-placeholder {
      width: 48px;
      height: 48px;
      border-radius: 6px;
      flex-shrink: 0;
      background: var(--color-border);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.25rem;
      color: var(--color-text-muted);
    }
    .feed-content {
      flex: 1;
      min-width: 0;
    }
    .feed-card h3 {
      margin: 0 0 0.5rem 0;
      font-size: 1rem;
    }
    .feed-card h3 a {
      color: var(--color-text);
      text-decoration: none;
    }
    .feed-card h3 a:hover { text-decoration: underline; }
    .feed-card p {
      margin: 0;
      font-size: 0.875rem;
      color: var(--color-text-muted);
    }
    .feed-meta {
      margin-top: 0.5rem;
      font-size: 0.75rem;
      color: var(--color-text-muted);
    }
    .nav-links {
      margin-bottom: 1rem;
    }
    .nav-links a {
      color: var(--color-text);
      margin-right: 1rem;
    }
  </style>
</head>
<body>
  <nav class="blogroll-nav" style="justify-content: flex-start; padding: 1rem 0;">
    <a href="/">Home</a>
    <a href="/` + blogrollSlug + `/">Blogroll</a>
    <a href="/` + readerSlug + `/">Reader</a>
  </nav>
  <div class="blogroll-page">
    <header class="blogroll-header" style="text-align: left;">
      <h1>Blogroll</h1>
      <p class="blogroll-subtitle">`)
	sb.WriteString(fmt.Sprintf("%d blogs and feeds I follow", len(feeds)))
	sb.WriteString(`</p>
    </header>
`)

	for _, cat := range categories {
		sb.WriteString(fmt.Sprintf(`    <section class="blogroll-category">
      <h2>%s</h2>
      <div class="blogroll-grid">
`, html.EscapeString(cat.Name)))

		for _, feed := range cat.Feeds {
			sb.WriteString(`        <article class="blogroll-card feed-card">
`)
			// Feed avatar - prefer AvatarURL (person/site avatar) over ImageURL (may be article image)
			avatarURL := feed.AvatarURL
			if avatarURL == "" {
				avatarURL = feed.ImageURL
			}
			if avatarURL != "" {
				sb.WriteString(fmt.Sprintf(`          <img src=%q alt="" class="feed-avatar" loading="lazy">
`, avatarURL))
			} else {
				// Placeholder with first letter
				initial := "?"
				if feed.Title != "" {
					initial = strings.ToUpper(string([]rune(feed.Title)[0:1]))
				}
				sb.WriteString(fmt.Sprintf(`          <div class="feed-avatar-placeholder">%s</div>
`, initial))
			}
			sb.WriteString(`          <div class="feed-content">
            <h3 class="blogroll-card-title">`)
			if feed.SiteURL != "" {
				sb.WriteString(fmt.Sprintf(`<a href=%q target="_blank" rel="noopener">%s</a>`,
					feed.SiteURL,
					html.EscapeString(feed.Title)))
			} else {
				sb.WriteString(html.EscapeString(feed.Title))
			}
			sb.WriteString(`</h3>
`)
			if feed.Description != "" {
				sb.WriteString(fmt.Sprintf(`            <p class="blogroll-card-description">%s</p>
`, html.EscapeString(blogrollTruncateString(feed.Description, 150))))
			}
			sb.WriteString(fmt.Sprintf(`            <footer class="blogroll-card-meta feed-meta">
              <span class="blogroll-card-count">%d posts</span>
            </footer>
          </div>
        </article>
`, feed.EntryCount))
		}

		sb.WriteString(`      </div>
    </section>
`)
	}

	sb.WriteString(`  </div>
</body>
</html>`)

	return sb.String()
}

// renderReaderFallback generates a basic reader page that uses theme CSS if available.
func (p *BlogrollPlugin) renderReaderFallback(entries []*models.ExternalEntry, feedIndex map[string]*models.ExternalFeed, page models.ReaderPage, config models.BlogrollConfig, variant readerVariant, readerTypes []map[string]interface{}) string {
	var sb strings.Builder

	// Get configured slugs with defaults
	blogrollSlug := config.BlogrollSlug
	if blogrollSlug == "" {
		blogrollSlug = defaultBlogrollSlug
	}
	readerSlug := config.ReaderSlug
	if readerSlug == "" {
		readerSlug = defaultReaderSlug
	}

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>` + html.EscapeString(variant.Title) + `</title>
  <!-- Theme CSS - uses site's configured palette if available -->
  <link rel="stylesheet" href="/css/variables.css">
  <link rel="stylesheet" href="/css/palette.css">
  <link rel="stylesheet" href="/css/main.css">
  <link rel="stylesheet" href="/css/components.css">
  <style>
    /* Fallback styles if theme CSS is not available */
    :root {
      --color-background: #ffffff;
      --color-text: #1a1a1a;
      --color-text-muted: #666666;
      --color-border: #e0e0e0;
      --color-surface: #f8f9fa;
      --color-primary: #3b82f6;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --color-background: #1a1a1a;
        --color-text: #f0f0f0;
        --color-text-muted: #999999;
        --color-border: #333333;
        --color-surface: #2a2a2a;
        --color-primary: #60a5fa;
      }
    }
    .pagination {
      display: flex;
      justify-content: center;
      gap: 1rem;
      margin: 2rem 0;
      padding: 1rem 0;
    }
    .pagination a, .pagination span {
      padding: 0.5rem 1rem;
      border: 1px solid var(--color-border);
      border-radius: 4px;
      text-decoration: none;
      color: var(--color-text);
    }
    .pagination a:hover {
      background: var(--color-surface);
    }
    .pagination .current {
      background: var(--color-primary);
      color: white;
      border-color: var(--color-primary);
    }
    .pagination .disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      max-width: 900px;
      margin: 0 auto;
      padding: 2rem;
      background: var(--color-background);
      color: var(--color-text);
      line-height: 1.6;
    }
    h1 { margin-bottom: 0.5rem; }
    .subtitle { color: var(--color-text-muted); margin-bottom: 2rem; }
    .entry-list { list-style: none; padding: 0; margin: 0; }
    .entry {
      border-bottom: 1px solid var(--color-border);
      padding: 1.5rem 0;
      display: flex;
      gap: 1rem;
    }
    .entry:last-child { border-bottom: none; }
    .entry-image {
      width: 120px;
      height: 80px;
      border-radius: 6px;
      object-fit: cover;
      flex-shrink: 0;
      background: var(--color-border);
    }
    .entry-content {
      flex: 1;
      min-width: 0;
    }
    .entry h2 {
      margin: 0 0 0.5rem 0;
      font-size: 1.25rem;
    }
    .entry h2 a {
      color: var(--color-text);
      text-decoration: none;
    }
    .entry h2 a:hover { text-decoration: underline; }
    .entry-meta {
      font-size: 0.875rem;
      color: var(--color-text-muted);
      margin-bottom: 0.5rem;
    }
    .entry-meta a { color: var(--color-primary); }
    .entry-description {
      color: var(--color-text);
      font-size: 0.9375rem;
    }
    .nav-links {
      margin-bottom: 1rem;
    }
    .nav-links a {
      color: var(--color-text);
      margin-right: 1rem;
    }
    .reader-type-links {
      display: flex;
      flex-wrap: wrap;
      gap: 0.75rem;
      margin: 0 0 1.5rem;
    }
    .reader-type-links a, .reader-type-links span {
      border: 1px solid var(--color-border);
      border-radius: 999px;
      padding: 0.35rem 0.8rem;
      text-decoration: none;
      color: var(--color-text);
      background: var(--color-surface);
      font-size: 0.95rem;
    }
    .reader-type-links .active {
      background: var(--color-primary);
      border-color: var(--color-primary);
      color: white;
    }
    @media (max-width: 600px) {
      .entry {
        flex-direction: column;
      }
      .entry-image {
        width: 100%;
        height: 160px;
      }
    }
  </style>
`)
	// Add HTMX if using that pagination type
	if page.PaginationType == models.PaginationHTMX {
		htmxURL := p.resolveAssetURL("htmx", "https://unpkg.com/htmx.org@1.9.10")
		sb.WriteString(fmt.Sprintf("  <script src=%q></script>\n", htmxURL))
	}
	sb.WriteString(`</head>
<body>
  <nav class="reader-nav" style="justify-content: flex-start; padding: 1rem 0;">
    <a href="/">Home</a>
    <a href="/` + blogrollSlug + `/">Blogroll</a>
    <a href="/` + readerSlug + `/">Reader</a>
  </nav>
  <div class="reader-page">
    <header class="reader-header" style="text-align: left;">
      <h1>` + html.EscapeString(variant.Title) + `</h1>
      <p class="reader-subtitle">` + html.EscapeString(variant.Description) + `</p>
    </header>
`)
	sb.WriteString(renderReaderTypeLinks(readerTypes))

	// Add content container for HTMX
	if page.PaginationType == models.PaginationHTMX {
		sb.WriteString(`    <div id="reader-content">
`)
	}

	sb.WriteString(p.renderReaderTimeline(p.readerDayGroups(entries, feedIndex, config), true))

	// Render pagination
	sb.WriteString(p.renderPagination(page))

	if page.PaginationType == models.PaginationHTMX {
		sb.WriteString(`    </div>
`)
	}

	sb.WriteString(`  </div>
</body>
</html>`)

	return sb.String()
}

// renderReaderPartial generates just the content portion for HTMX updates.
func (p *BlogrollPlugin) renderReaderPartial(entries []*models.ExternalEntry, feedIndex map[string]*models.ExternalFeed, page models.ReaderPage, config models.BlogrollConfig, _ readerVariant, readerTypes []map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString(renderReaderTypeLinks(readerTypes))

	sb.WriteString(p.renderReaderTimeline(p.readerDayGroups(entries, feedIndex, config), false))

	// Render pagination
	sb.WriteString(p.renderPagination(page))

	return sb.String()
}

func renderReaderTypeLinks(readerTypes []map[string]interface{}) string {
	if len(readerTypes) <= 1 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("    <nav class=\"reader-type-links\" aria-label=\"Reader types\">\n")
	for _, item := range readerTypes {
		title := stringValue(item, "title")
		url := stringValue(item, "url")
		count := 0
		if v, ok := item["entry_count"].(int); ok {
			count = v
		}
		label := html.EscapeString(title)
		if count > 0 {
			label += fmt.Sprintf(" (%d)", count)
		}
		if active, ok := item["active"].(bool); ok && active {
			sb.WriteString("      <span class=\"active\">" + label + "</span>\n")
		} else {
			sb.WriteString("      <a href=\"" + html.EscapeString(url) + "\">" + label + "</a>\n")
		}
	}
	sb.WriteString("    </nav>\n")
	return sb.String()
}

// renderPagination generates pagination navigation HTML.
func (p *BlogrollPlugin) renderPagination(page models.ReaderPage) string {
	if page.TotalPages <= 1 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(`    <nav class="pagination" aria-label="Page navigation">
`)

	// Previous link
	if page.HasPrev {
		if page.PaginationType == models.PaginationHTMX {
			sb.WriteString(fmt.Sprintf(`      <a href="%s" hx-get="%spartial/" hx-target="#reader-content" hx-swap="innerHTML">&laquo; Previous</a>
`, page.PrevURL, page.PrevURL))
		} else {
			sb.WriteString(fmt.Sprintf(`      <a href="%s">&laquo; Previous</a>
`, page.PrevURL))
		}
	} else {
		sb.WriteString(`      <span class="disabled">&laquo; Previous</span>
`)
	}

	// Page numbers
	for i, url := range page.PageURLs {
		pageNum := i + 1
		if pageNum == page.Number {
			sb.WriteString(fmt.Sprintf(`      <span class="current">%d</span>
`, pageNum))
		} else {
			if page.PaginationType == models.PaginationHTMX {
				sb.WriteString(fmt.Sprintf(`      <a href="%s" hx-get="%spartial/" hx-target="#reader-content" hx-swap="innerHTML">%d</a>
`, url, url, pageNum))
			} else {
				sb.WriteString(fmt.Sprintf(`      <a href="%s">%d</a>
`, url, pageNum))
			}
		}
	}

	// Next link
	if page.HasNext {
		if page.PaginationType == models.PaginationHTMX {
			sb.WriteString(fmt.Sprintf(`      <a href="%s" hx-get="%spartial/" hx-target="#reader-content" hx-swap="innerHTML">Next &raquo;</a>
`, page.NextURL, page.NextURL))
		} else {
			sb.WriteString(fmt.Sprintf(`      <a href="%s">Next &raquo;</a>
`, page.NextURL))
		}
	} else {
		sb.WriteString(`      <span class="disabled">Next &raquo;</span>
`)
	}

	sb.WriteString(`    </nav>
`)

	return sb.String()
}

// getBlogrollConfig retrieves blogroll configuration from the manager config.
func getBlogrollConfig(config *lifecycle.Config) models.BlogrollConfig {
	if config.Extra == nil {
		return models.NewBlogrollConfig()
	}

	if blogroll, ok := config.Extra["blogroll"]; ok {
		if bc, ok := blogroll.(models.BlogrollConfig); ok {
			return bc
		}
	}

	return models.NewBlogrollConfig()
}

// blogrollSlugify converts a string to a URL-safe slug.
func blogrollSlugify(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove non-alphanumeric characters except hyphens
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}

	// Remove multiple consecutive hyphens
	slug := result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim leading/trailing hyphens
	return strings.Trim(slug, "-")
}

// blogrollTruncateString truncates a string to the specified length.
func blogrollTruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// blogrollStripHTML removes HTML tags from a string.
func blogrollStripHTML(s string) string {
	// Use pre-compiled regex for HTML stripping
	s = blogrollHTMLTagRegex.ReplaceAllString(s, "")

	// Decode common HTML entities
	s = html.UnescapeString(s)

	// Normalize whitespace
	s = strings.Join(strings.Fields(s), " ")

	return s
}

// generateFallbackImageURL generates a fallback image URL from a template.
// The template should contain {url} as a placeholder for the URL-encoded entry URL.
func generateFallbackImageURL(template, entryURL string) string {
	encodedURL := neturl.QueryEscape(entryURL)
	return strings.ReplaceAll(template, "{url}", encodedURL)
}

// discoverFeedAvatar attempts to discover an avatar for a feed using IndieWeb methods.
// It tries h-card u-photo, WebFinger rel=avatar, and .well-known/avatar.
// This is best-effort and failures are silently ignored.
func (p *BlogrollPlugin) discoverFeedAvatar(ctx context.Context, feed *models.ExternalFeed, config models.ExternalFeedConfig, timeout int) {
	// Create an updater for avatar discovery
	updater := blogroll.NewUpdater(time.Duration(timeout) * time.Second)

	// Determine the resource for WebFinger lookup
	// Use handle if available (e.g., "acct:user@example.com")
	resource := ""
	if config.Handle != "" {
		// If handle looks like an email or acct URI, use it directly
		if strings.Contains(config.Handle, "@") {
			if !strings.HasPrefix(config.Handle, "acct:") {
				resource = "acct:" + config.Handle
			} else {
				resource = config.Handle
			}
		}
	}

	// Attempt avatar discovery
	avatarResult, err := updater.DiscoverAvatar(ctx, feed.SiteURL, resource)
	if err != nil {
		// Silent failure - avatar discovery is best-effort
		return
	}

	if avatarResult != nil && avatarResult.URL != "" {
		feed.AvatarURL = avatarResult.URL
		feed.AvatarSource = string(avatarResult.Source)
	}
}

// Ensure BlogrollPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*BlogrollPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*BlogrollPlugin)(nil)
	_ lifecycle.CollectPlugin   = (*BlogrollPlugin)(nil)
	_ lifecycle.WritePlugin     = (*BlogrollPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*BlogrollPlugin)(nil)
)

// CI trigger
