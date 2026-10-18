// Package blogroll provides functionality for managing blogroll metadata.
package blogroll

import (
	"context"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
)

// Pre-compiled regex patterns for HTML metadata extraction.
var (
	// parseHTMLMetadata patterns
	htmlTitleRe        = regexp.MustCompile(`(?i)<title[^>]*>([^<]*)</title>`)
	htmlMetaRe         = regexp.MustCompile(`(?i)<meta\s+[^>]*(?:property|name)=["']([^"']+)["'][^>]*content=["']([^"']*)["'][^>]*>`)
	htmlMetaRe2        = regexp.MustCompile(`(?i)<meta\s+[^>]*content=["']([^"']*)["'][^>]*(?:property|name)=["']([^"']+)["'][^>]*>`)
	htmlIconRe         = regexp.MustCompile(`(?i)<link\s+[^>]*rel=["'](?:icon|shortcut icon|apple-touch-icon)[^"']*["'][^>]*href=["']([^"']+)["'][^>]*>`)
	htmlIconRe2        = regexp.MustCompile(`(?i)<link\s+[^>]*href=["']([^"']+)["'][^>]*rel=["'](?:icon|shortcut icon|apple-touch-icon)[^"']*["'][^>]*>`)
	rssCategoryRe      = regexp.MustCompile(`(?is)<category[^>]*>(.*?)</category>`)
	atomCategoryTermRe = regexp.MustCompile(`(?is)<category\b[^>]*term=["']([^"']+)["'][^>]*/?>`)
	atomCategoryTextRe = regexp.MustCompile(`(?is)<category\b[^>]*>(.*?)</category>`)
	atomAltLinkRe      = regexp.MustCompile(`(?is)<link\b[^>]*rel=["']alternate["'][^>]*href=["']([^"']+)["'][^>]*/?>`)
	atomAltLinkRe2     = regexp.MustCompile(`(?is)<link\b[^>]*href=["']([^"']+)["'][^>]*rel=["']alternate["'][^>]*/?>`)
	quotedTagRe        = regexp.MustCompile(`"([^"]+)"|(\S+)`)
)

// xmlTagCache caches compiled regex patterns for XML tag extraction.
// The tag set is finite (~13 known tags), so this is bounded.
var xmlTagCache sync.Map // map[string][]*regexp.Regexp

// getXMLTagPatterns returns pre-compiled regex patterns for the given XML tag name.
func getXMLTagPatterns(tag string) []*regexp.Regexp {
	if cached, ok := xmlTagCache.Load(tag); ok {
		return cached.([]*regexp.Regexp) //nolint:errcheck // type is guaranteed by xmlTagCache.Store
	}
	patterns := []*regexp.Regexp{
		regexp.MustCompile(fmt.Sprintf("<%s[^>]*>([^<]*)</%s>", tag, tag)),
		regexp.MustCompile(fmt.Sprintf("<%s[^>]*>(.*?)</%s>", tag, tag)),
	}
	xmlTagCache.Store(tag, patterns)
	return patterns
}

// Updater handles fetching and extracting metadata from external sites.
type Updater struct {
	client  *http.Client
	timeout time.Duration
}

// NewUpdater creates a new Updater with the given timeout.
func NewUpdater(timeout time.Duration) *Updater {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := &http.Client{
		Timeout: timeout,
	}
	buildstats.InstrumentHTTPClient(client)

	return &Updater{
		client:  client,
		timeout: timeout,
	}
}

// Metadata represents extracted metadata from a website.
type Metadata struct {
	// From OpenGraph
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
	SiteURL     string   `json:"site_url,omitempty"`
	Tags        []string `json:"tags,omitempty"`

	// Author information
	Author string `json:"author,omitempty"`

	// Avatar (person/site representative image)
	// Distinct from ImageURL which may be an article/page image
	AvatarURL    string       `json:"avatar_url,omitempty"`
	AvatarSource AvatarSource `json:"avatar_source,omitempty"`

	// From feed
	FeedTitle       string     `json:"feed_title,omitempty"`
	FeedDescription string     `json:"feed_description,omitempty"`
	FeedAuthor      string     `json:"feed_author,omitempty"`
	FeedImageURL    string     `json:"feed_image_url,omitempty"`
	FeedTags        []string   `json:"feed_tags,omitempty"`
	LastUpdated     *time.Time `json:"last_updated,omitempty"`

	// Source tracking
	Source string `json:"source,omitempty"` // "opengraph", "meta", "feed"
}

// UpdateResult contains the result of updating a single feed's metadata.
type UpdateResult struct {
	FeedURL     string    `json:"feed_url"`
	Handle      string    `json:"handle,omitempty"`
	OldMetadata *Metadata `json:"old_metadata,omitempty"`
	NewMetadata *Metadata `json:"new_metadata,omitempty"`
	Updated     bool      `json:"updated"`
	Error       string    `json:"error,omitempty"`
}

// FetchMetadata fetches metadata from a site URL or feed URL.
// It tries multiple sources in order: OpenGraph, HTML meta tags, feed metadata.
// Additionally, it attempts avatar discovery via h-card, WebFinger, and .well-known/avatar.
func (u *Updater) FetchMetadata(ctx context.Context, feedURL string) (*Metadata, error) {
	return u.FetchMetadataWithResource(ctx, feedURL, "")
}

// FetchMetadataWithResource fetches metadata and attempts avatar discovery.
// The resource parameter is used for WebFinger lookups (e.g., "acct:user@example.com").
func (u *Updater) FetchMetadataWithResource(ctx context.Context, feedURL, resource string) (*Metadata, error) {
	metadata := &Metadata{}

	// First, try to extract site URL from feed URL
	defaultSiteURL, err := extractSiteURL(feedURL)
	if err != nil {
		return nil, fmt.Errorf("invalid feed URL: %w", err)
	}

	// Fetch feed metadata
	feedMetadata, feedErr := u.fetchFeedMetadata(ctx, feedURL)
	siteURL := defaultSiteURL
	if feedErr == nil && feedMetadata != nil && feedMetadata.SiteURL != "" {
		siteURL = feedMetadata.SiteURL
	}

	// Fetch site metadata (OpenGraph + HTML meta)
	siteMetadata, siteErr := u.fetchSiteMetadata(ctx, siteURL)
	if siteErr == nil && siteMetadata != nil {
		mergeMetadata(metadata, siteMetadata)
	}

	if feedErr == nil && feedMetadata != nil {
		// Feed metadata fills in gaps but doesn't overwrite
		mergeMetadataWithoutOverwrite(metadata, feedMetadata)
	}

	// If we couldn't get any metadata, return an error
	if siteErr != nil && feedErr != nil {
		return nil, fmt.Errorf("failed to fetch metadata from %s: %w", feedURL, siteErr)
	}

	// Set the site URL if not already set
	if metadata.SiteURL == "" {
		metadata.SiteURL = siteURL
	}

	// Attempt avatar discovery (best-effort, doesn't fail the overall operation)
	if avatarResult, _ := u.DiscoverAvatar(ctx, siteURL, resource); avatarResult != nil { //nolint:errcheck // intentionally ignoring error for best-effort operation
		metadata.AvatarURL = avatarResult.URL
		metadata.AvatarSource = avatarResult.Source
	}

	return metadata, nil
}

// fetchURL is a helper that fetches a URL and returns the body as bytes.
func (u *Updater) fetchURL(ctx context.Context, targetURL, accept string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (Blogroll Updater)")
	req.Header.Set("Accept", accept)

	resp, err := u.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Read body with a limit (1MB)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return body, nil
}

// fetchSiteMetadata fetches OpenGraph and HTML meta tags from a site.
func (u *Updater) fetchSiteMetadata(ctx context.Context, siteURL string) (*Metadata, error) {
	body, err := u.fetchURL(ctx, siteURL, "text/html,application/xhtml+xml")
	if err != nil {
		return nil, err
	}
	return parseHTMLMetadata(string(body), siteURL), nil
}

// fetchFeedMetadata fetches metadata from an RSS/Atom feed.
func (u *Updater) fetchFeedMetadata(ctx context.Context, feedURL string) (*Metadata, error) {
	body, err := u.fetchURL(ctx, feedURL, "application/rss+xml, application/atom+xml, application/xml, text/xml")
	if err != nil {
		return nil, err
	}
	return parseFeedMetadata(string(body)), nil
}

// parseHTMLMetadata extracts OpenGraph and meta tags from HTML using regex.
func parseHTMLMetadata(htmlContent, baseURL string) *Metadata {
	metadata := &Metadata{}

	// Extract <title> tag
	if matches := htmlTitleRe.FindStringSubmatch(htmlContent); len(matches) > 1 {
		metadata.Title = decodeMetadataText(matches[1])
	}

	// Extract meta tags with pattern: <meta property="..." content="..."> or <meta name="..." content="...">
	processMetaMatch := func(prop, content string) {
		processOpenGraphMeta(prop, content, metadata, baseURL)
		processHTMLMeta(prop, content, metadata, baseURL)
	}

	for _, matches := range htmlMetaRe.FindAllStringSubmatch(htmlContent, -1) {
		if len(matches) > 2 {
			processMetaMatch(matches[1], matches[2])
		}
	}

	for _, matches := range htmlMetaRe2.FindAllStringSubmatch(htmlContent, -1) {
		if len(matches) > 2 {
			processMetaMatch(matches[2], matches[1])
		}
	}

	// Extract link rel="icon" for favicon
	if metadata.ImageURL == "" {
		if matches := htmlIconRe.FindStringSubmatch(htmlContent); len(matches) > 1 {
			metadata.ImageURL = resolveURL(matches[1], baseURL)
		} else if matches := htmlIconRe2.FindStringSubmatch(htmlContent); len(matches) > 1 {
			metadata.ImageURL = resolveURL(matches[1], baseURL)
		}
	}

	metadata.Title = normalizeSiteMetadataTitle(metadata.Title, baseURL)

	if metadata.Title != "" || metadata.Description != "" || metadata.ImageURL != "" {
		metadata.Source = "opengraph"
	}

	return metadata
}

// processOpenGraphMeta handles OpenGraph protocol meta tags.
func processOpenGraphMeta(property, content string, metadata *Metadata, baseURL string) {
	content = decodeMetadataText(content)
	property = strings.ToLower(strings.TrimSpace(property))

	switch property {
	case "og:title":
		if metadata.Title == "" {
			metadata.Title = content
		}
	case "og:description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "og:image":
		if metadata.ImageURL == "" {
			metadata.ImageURL = resolveURL(content, baseURL)
		}
	case "og:url":
		if metadata.SiteURL == "" {
			metadata.SiteURL = content
		}
	case "og:site_name":
		// Can be used as fallback title
		if metadata.Title == "" {
			metadata.Title = content
		}
	}
}

// processHTMLMeta handles standard HTML meta tags.
func processHTMLMeta(name, content string, metadata *Metadata, baseURL string) {
	name = strings.ToLower(strings.TrimSpace(name))
	content = decodeMetadataText(content)

	switch name {
	case "description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "keywords":
		metadata.Tags = uniqueTags(append(metadata.Tags, splitMetadataTags(content)...))
	case "author":
		if metadata.Author == "" {
			metadata.Author = content
		}
	case "article:tag":
		metadata.Tags = uniqueTags(append(metadata.Tags, content))
	case "twitter:title":
		if metadata.Title == "" {
			metadata.Title = content
		}
	case "twitter:description":
		if metadata.Description == "" {
			metadata.Description = content
		}
	case "twitter:image":
		if metadata.ImageURL == "" {
			metadata.ImageURL = resolveURL(content, baseURL)
		}
	}
}

// parseFeedMetadata extracts metadata from RSS/Atom feed content.
func parseFeedMetadata(feedContent string) *Metadata {
	metadata := &Metadata{
		Source: "feed",
	}

	// Detect feed type and parse accordingly
	if strings.Contains(feedContent, "<feed") && strings.Contains(feedContent, "xmlns=\"http://www.w3.org/2005/Atom\"") {
		parseAtomFeedMetadata(feedContent, metadata)
	} else {
		parseRSSFeedMetadata(feedContent, metadata)
	}

	return metadata
}

// parseAtomFeedMetadata parses Atom feed metadata.
func parseAtomFeedMetadata(content string, metadata *Metadata) {
	// Extract title
	if title := extractXMLTag(content, "title"); title != "" {
		metadata.FeedTitle = stripCDATA(title)
	}

	// Extract subtitle (description)
	if subtitle := extractXMLTag(content, "subtitle"); subtitle != "" {
		metadata.FeedDescription = stripCDATA(subtitle)
	}

	// Extract author
	if authorName := extractNestedXMLTag(content, "author", "name"); authorName != "" {
		metadata.FeedAuthor = stripCDATA(authorName)
	}
	if authorURI := extractNestedXMLTag(content, "author", "uri"); authorURI != "" {
		metadata.SiteURL = stripCDATA(authorURI)
	}
	if metadata.SiteURL == "" {
		metadata.SiteURL = extractAtomAlternateLink(content)
	}

	// Extract icon/logo
	if icon := extractXMLTag(content, "icon"); icon != "" {
		metadata.FeedImageURL = stripCDATA(icon)
	} else if logo := extractXMLTag(content, "logo"); logo != "" {
		metadata.FeedImageURL = stripCDATA(logo)
	}

	// Extract updated timestamp
	if updated := extractXMLTag(content, "updated"); updated != "" {
		if t, err := time.Parse(time.RFC3339, stripCDATA(updated)); err == nil {
			metadata.LastUpdated = &t
		}
	}

	metadata.FeedTags = extractAtomFeedTags(content)
}

// parseRSSFeedMetadata parses RSS feed metadata.
func parseRSSFeedMetadata(content string, metadata *Metadata) {
	// Find the channel section
	channelStart := strings.Index(content, "<channel")
	channelEnd := strings.Index(content, "</channel>")
	if channelStart == -1 || channelEnd == -1 {
		return
	}
	channel := content[channelStart:channelEnd]

	// Extract title (avoid item titles)
	if title := extractChannelTag(channel, "title"); title != "" {
		metadata.FeedTitle = stripCDATA(title)
	}

	// Extract description
	if desc := extractChannelTag(channel, "description"); desc != "" {
		metadata.FeedDescription = stripCDATA(desc)
	}

	if link := extractChannelTag(channel, "link"); link != "" {
		metadata.SiteURL = stripCDATA(link)
	}

	// Extract managing editor or webmaster as author
	if author := extractChannelTag(channel, "managingEditor"); author != "" {
		metadata.FeedAuthor = stripCDATA(extractEmailName(author))
	} else if author := extractChannelTag(channel, "webMaster"); author != "" {
		metadata.FeedAuthor = stripCDATA(extractEmailName(author))
	} else if author := extractChannelTag(channel, "author"); author != "" {
		metadata.FeedAuthor = stripCDATA(author)
	}

	// Extract image (try multiple sources)
	if imgURL := extractNestedXMLTag(channel, "image", "url"); imgURL != "" {
		metadata.FeedImageURL = stripCDATA(imgURL)
	} else if icon := extractXMLTag(channel, "icon"); icon != "" {
		// Fallback to Atom-style icon tag (some feeds mix formats)
		metadata.FeedImageURL = stripCDATA(icon)
	} else if logo := extractXMLTag(channel, "logo"); logo != "" {
		// Fallback to Atom-style logo tag
		metadata.FeedImageURL = stripCDATA(logo)
	}

	// Extract lastBuildDate or pubDate
	if lastBuild := extractChannelTag(channel, "lastBuildDate"); lastBuild != "" {
		if t, err := parseRSSDate(stripCDATA(lastBuild)); err == nil {
			metadata.LastUpdated = &t
		}
	} else if pubDate := extractChannelTag(channel, "pubDate"); pubDate != "" {
		if t, err := parseRSSDate(stripCDATA(pubDate)); err == nil {
			metadata.LastUpdated = &t
		}
	}

	metadata.FeedTags = extractRSSFeedTags(channel)
}

func extractRSSFeedTags(channel string) []string {
	itemStart := strings.Index(channel, "<item")
	searchContent := channel
	if itemStart > 0 {
		searchContent = channel[:itemStart]
	}

	tags := make([]string, 0)
	for _, match := range rssCategoryRe.FindAllStringSubmatch(searchContent, -1) {
		if len(match) > 1 {
			tags = append(tags, stripCDATA(match[1]))
		}
	}

	return uniqueTags(tags)
}

func extractAtomFeedTags(content string) []string {
	entryStart := strings.Index(content, "<entry")
	searchContent := content
	if entryStart > 0 {
		searchContent = content[:entryStart]
	}

	tags := make([]string, 0)
	for _, match := range atomCategoryTermRe.FindAllStringSubmatch(searchContent, -1) {
		if len(match) > 1 {
			tags = append(tags, stripCDATA(match[1]))
		}
	}
	for _, match := range atomCategoryTextRe.FindAllStringSubmatch(searchContent, -1) {
		if len(match) > 1 {
			tags = append(tags, stripCDATA(match[1]))
		}
	}

	return uniqueTags(tags)
}

func extractAtomAlternateLink(content string) string {
	for _, match := range atomAltLinkRe.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			return strings.TrimSpace(match[1])
		}
	}
	for _, match := range atomAltLinkRe2.FindAllStringSubmatch(content, -1) {
		if len(match) > 1 {
			return strings.TrimSpace(match[1])
		}
	}
	return ""
}

func splitMetadataTags(content string) []string {
	content = decodeMetadataText(content)
	if content == "" {
		return nil
	}

	if strings.Contains(content, ",") {
		parts := strings.Split(content, ",")
		tags := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.Trim(strings.TrimSpace(part), `"`)
			if trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
		return tags
	}

	tags := make([]string, 0)
	for _, match := range quotedTagRe.FindAllStringSubmatch(content, -1) {
		switch {
		case len(match) > 1 && strings.TrimSpace(match[1]) != "":
			tags = append(tags, strings.TrimSpace(match[1]))
		case len(match) > 2 && strings.TrimSpace(match[2]) != "":
			tags = append(tags, strings.TrimSpace(match[2]))
		}
	}

	return tags
}

func uniqueTags(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(tags))
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}

		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, trimmed)
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// extractChannelTag extracts a tag value from the channel section, avoiding item content.
func extractChannelTag(channel, tag string) string {
	// Find first occurrence before any <item> tag
	itemStart := strings.Index(channel, "<item")
	searchContent := channel
	if itemStart > 0 {
		searchContent = channel[:itemStart]
	}
	return extractXMLTag(searchContent, tag)
}

// extractXMLTag extracts the content of an XML tag.
// Patterns are cached per tag name via xmlTagCache.
func extractXMLTag(content, tag string) string {
	for _, re := range getXMLTagPatterns(tag) {
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

// extractNestedXMLTag extracts a nested tag value.
// Uses cached patterns via getXMLTagPatterns.
func extractNestedXMLTag(content, parent, child string) string {
	// Find parent tag using cached pattern (second pattern: greedy inner match)
	patterns := getXMLTagPatterns(parent)
	if len(patterns) < 2 {
		return ""
	}
	matches := patterns[1].FindStringSubmatch(content)
	if len(matches) > 1 {
		return extractXMLTag(matches[1], child)
	}
	return ""
}

// stripCDATA removes CDATA wrapper and trims whitespace.
func stripCDATA(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "<![CDATA[") && strings.HasSuffix(s, "]]>") {
		s = s[9 : len(s)-3]
	}
	return decodeMetadataText(s)
}

func decodeMetadataText(s string) string {
	return strings.TrimSpace(html.UnescapeString(s))
}

func normalizeSiteMetadataTitle(title, baseURL string) string {
	title = decodeMetadataText(title)
	if title == "" {
		return ""
	}

	u, err := url.Parse(baseURL)
	if err == nil {
		host := strings.ToLower(strings.TrimPrefix(u.Host, "www."))
		if host == "youtube.com" {
			title = strings.TrimSpace(strings.TrimSuffix(title, " - YouTube"))
		}
	}

	return title
}

// extractEmailName extracts a name from an email format like "email (Name)" or "Name <email>".
func extractEmailName(s string) string {
	// Format: email (Name)
	if idx := strings.Index(s, "("); idx > 0 {
		end := strings.Index(s, ")")
		if end > idx {
			return strings.TrimSpace(s[idx+1 : end])
		}
	}
	// Format: Name <email>
	if idx := strings.Index(s, "<"); idx > 0 {
		return strings.TrimSpace(s[:idx])
	}
	return s
}

// parseRSSDate parses RSS date formats.
func parseRSSDate(s string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"02 Jan 2006 15:04:05 -0700",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}

// extractSiteURL extracts the base site URL from a feed URL.
func extractSiteURL(feedURL string) (string, error) {
	u, err := url.Parse(feedURL)
	if err != nil {
		return "", err
	}

	// Return just the scheme and host
	return fmt.Sprintf("%s://%s", u.Scheme, u.Host), nil
}

// resolveURL resolves a potentially relative URL against a base URL.
func resolveURL(href, baseURL string) string {
	if href == "" {
		return ""
	}

	// Already absolute
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	// Protocol-relative
	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}

	// Relative URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	ref, err := url.Parse(href)
	if err != nil {
		return href
	}

	return base.ResolveReference(ref).String()
}

// mergeMetadata merges source metadata into target, overwriting non-empty values.
func mergeMetadata(target, source *Metadata) {
	if source.Title != "" {
		target.Title = source.Title
	}
	if source.Description != "" {
		target.Description = source.Description
	}
	if source.ImageURL != "" {
		target.ImageURL = source.ImageURL
	}
	if source.SiteURL != "" {
		target.SiteURL = source.SiteURL
	}
	if len(source.Tags) > 0 {
		target.Tags = append([]string{}, source.Tags...)
	}
	if source.Author != "" {
		target.Author = source.Author
	}
	if source.AvatarURL != "" {
		target.AvatarURL = source.AvatarURL
		target.AvatarSource = source.AvatarSource
	}
	if source.FeedTitle != "" {
		target.FeedTitle = source.FeedTitle
	}
	if source.FeedDescription != "" {
		target.FeedDescription = source.FeedDescription
	}
	if source.FeedAuthor != "" {
		target.FeedAuthor = source.FeedAuthor
	}
	if source.FeedImageURL != "" {
		target.FeedImageURL = source.FeedImageURL
	}
	if len(source.FeedTags) > 0 {
		target.FeedTags = append([]string{}, source.FeedTags...)
	}
	if source.LastUpdated != nil {
		target.LastUpdated = source.LastUpdated
	}
	if source.Source != "" {
		target.Source = source.Source
	}
}

// mergeMetadataWithoutOverwrite merges source into target only for empty fields.
func mergeMetadataWithoutOverwrite(target, source *Metadata) {
	if target.Title == "" && source.FeedTitle != "" {
		target.Title = source.FeedTitle
	}
	if target.Description == "" && source.FeedDescription != "" {
		target.Description = source.FeedDescription
	}
	if target.ImageURL == "" && source.FeedImageURL != "" {
		target.ImageURL = source.FeedImageURL
	}
	if target.Author == "" && source.FeedAuthor != "" {
		target.Author = source.FeedAuthor
	}
	if len(target.Tags) == 0 && len(source.FeedTags) > 0 {
		target.Tags = append([]string{}, source.FeedTags...)
	}
	if target.AvatarURL == "" && source.AvatarURL != "" {
		target.AvatarURL = source.AvatarURL
		target.AvatarSource = source.AvatarSource
	}
	if target.LastUpdated == nil && source.LastUpdated != nil {
		target.LastUpdated = source.LastUpdated
	}
	// Keep feed-specific fields
	if source.FeedTitle != "" {
		target.FeedTitle = source.FeedTitle
	}
	if source.FeedDescription != "" {
		target.FeedDescription = source.FeedDescription
	}
	if source.FeedAuthor != "" {
		target.FeedAuthor = source.FeedAuthor
	}
	if source.FeedImageURL != "" {
		target.FeedImageURL = source.FeedImageURL
	}
	if len(source.FeedTags) > 0 {
		target.FeedTags = append([]string{}, source.FeedTags...)
	}
}
