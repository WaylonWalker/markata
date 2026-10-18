// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

// MentionSource represents the detected source platform of a webmention.
type MentionSource int

const (
	// SourceWeb is a standard web webmention
	SourceWeb MentionSource = iota
	// SourceBridgyBluesky is a Bluesky mention via Bridgy
	SourceBridgyBluesky
	// SourceBridgyTwitter is a Twitter mention via Bridgy
	SourceBridgyTwitter
	// SourceBridgyMastodon is a Mastodon mention via Bridgy
	SourceBridgyMastodon
	// SourceBridgyGitHub is a GitHub mention via Bridgy
	SourceBridgyGitHub
	// SourceBridgyFlickr is a Flickr mention via Bridgy
	SourceBridgyFlickr
	// SourceCustomBridge is a mention from a custom bridge
	SourceCustomBridge
)

// String returns the string representation of MentionSource.
func (s MentionSource) String() string {
	switch s {
	case SourceBridgyBluesky:
		return "bluesky"
	case SourceBridgyTwitter:
		return "twitter"
	case SourceBridgyMastodon:
		return "mastodon"
	case SourceBridgyGitHub:
		return "github"
	case SourceBridgyFlickr:
		return "flickr"
	case SourceCustomBridge:
		return "custom"
	default:
		return "web"
	}
}

// ReceivedWebMention represents an incoming webmention with bridging metadata.
type ReceivedWebMention struct {
	// Standard webmention.io fields
	URL        string `json:"url"`
	Source     string `json:"wm-source"`
	Target     string `json:"wm-target"`
	Published  string `json:"published,omitempty"`
	WMProperty string `json:"wm-property"`

	// Author information
	Author MentionAuthor `json:"author,omitempty"`

	// Content information
	Content MentionContent `json:"content,omitempty"`

	// Bridging metadata (enriched by detection)
	SourceSite  MentionSource `json:"source_site,omitempty"`
	Platform    string        `json:"platform,omitempty"`
	Handle      string        `json:"handle,omitempty"`
	OriginalURL string        `json:"original_url,omitempty"`
}

// MentionAuthor represents the author of a webmention.
type MentionAuthor struct {
	Name  string `json:"name,omitempty"`
	Photo string `json:"photo,omitempty"`
	URL   string `json:"url,omitempty"`
}

// MentionContent represents the content of a webmention.
type MentionContent struct {
	Text string `json:"text,omitempty"`
	HTML string `json:"html,omitempty"`
}

// InteractionType returns the interaction type based on wm-property.
func (m *ReceivedWebMention) InteractionType() string {
	switch m.WMProperty {
	case "like-of":
		return "like"
	case "repost-of":
		return "repost"
	case "in-reply-to":
		return "reply"
	case "bookmark-of":
		return "bookmark"
	case "mention-of":
		return "mention"
	default:
		return "mention"
	}
}

// BridgingDetector detects and enriches webmentions with bridging metadata.
type BridgingDetector struct {
	config models.BridgesConfig
}

// NewBridgingDetector creates a new BridgingDetector.
func NewBridgingDetector(config models.BridgesConfig) *BridgingDetector {
	return &BridgingDetector{config: config}
}

// DetectSource detects the source platform from URL patterns.
func (d *BridgingDetector) DetectSource(sourceURL, sourceContent string) MentionSource {
	// Check Bridgy Fed patterns first
	if d.config.BridgyFediverse {
		if source := d.detectBridgySource(sourceURL); source != SourceWeb {
			return source
		}
	}

	// Check content-based detection as fallback
	return d.detectFromContent(sourceURL, sourceContent)
}

// detectBridgySource detects Bridgy Fed source from URL patterns.
func (d *BridgingDetector) detectBridgySource(sourceURL string) MentionSource {
	lowerURL := strings.ToLower(sourceURL)

	// Bridgy Fed patterns: brid.gy/publish/<platform> or brid.gy/<platform>/
	if strings.Contains(lowerURL, "brid.gy") {
		if strings.Contains(lowerURL, "/bluesky") || strings.Contains(lowerURL, "bsky") {
			if d.config.Bluesky {
				return SourceBridgyBluesky
			}
		}
		if strings.Contains(lowerURL, "/twitter") || strings.Contains(lowerURL, "/x.com") {
			if d.config.Twitter {
				return SourceBridgyTwitter
			}
		}
		if strings.Contains(lowerURL, "/mastodon") || strings.Contains(lowerURL, "/fediverse") {
			if d.config.Mastodon {
				return SourceBridgyMastodon
			}
		}
		if strings.Contains(lowerURL, "/github") {
			if d.config.GitHub {
				return SourceBridgyGitHub
			}
		}
		if strings.Contains(lowerURL, "/flickr") {
			if d.config.Flickr {
				return SourceBridgyFlickr
			}
		}
	}

	// Also check for direct platform URLs that Bridgy might convert
	if strings.Contains(lowerURL, "bsky.app") || strings.Contains(lowerURL, "bsky.social") {
		if d.config.Bluesky {
			return SourceBridgyBluesky
		}
	}

	return SourceWeb
}

// detectFromContent detects platform from URL domain or content patterns.
func (d *BridgingDetector) detectFromContent(sourceURL, content string) MentionSource {
	lowerURL := strings.ToLower(sourceURL)
	lowerContent := strings.ToLower(content)

	// Check domain patterns
	if strings.Contains(lowerURL, "bsky.app") || strings.Contains(lowerURL, "bsky.social") {
		if d.config.Bluesky {
			return SourceBridgyBluesky
		}
	}
	if strings.Contains(lowerURL, "twitter.com") || strings.Contains(lowerURL, "x.com") {
		if d.config.Twitter {
			return SourceBridgyTwitter
		}
	}
	if d.isMastodonURL(lowerURL) {
		if d.config.Mastodon {
			return SourceBridgyMastodon
		}
	}
	if strings.Contains(lowerURL, "github.com") {
		if d.config.GitHub {
			return SourceBridgyGitHub
		}
	}
	if strings.Contains(lowerURL, "flickr.com") {
		if d.config.Flickr {
			return SourceBridgyFlickr
		}
	}

	// Check content patterns as last resort
	if strings.Contains(lowerContent, "bsky.app") || strings.Contains(lowerContent, "@bsky") {
		if d.config.Bluesky {
			return SourceBridgyBluesky
		}
	}

	return SourceWeb
}

// Common Mastodon instance domains
var mastodonDomains = []string{
	"mastodon.social",
	"mastodon.online",
	"mstdn.social",
	"fosstodon.org",
	"hachyderm.io",
	"infosec.exchange",
	"tech.lgbt",
	"social.coop",
	"aus.social",
}

// isMastodonURL checks if a URL is from a known Mastodon instance.
func (d *BridgingDetector) isMastodonURL(lowerURL string) bool {
	for _, domain := range mastodonDomains {
		if strings.Contains(lowerURL, domain) {
			return true
		}
	}
	// Also check for ActivityPub indicators in URL patterns
	if strings.Contains(lowerURL, "/@") || strings.Contains(lowerURL, "/users/") {
		return true
	}
	return false
}

// EnrichMention adds platform-specific metadata to a webmention.
func (d *BridgingDetector) EnrichMention(mention *ReceivedWebMention) {
	source := d.DetectSource(mention.Source, mention.Content.Text)
	mention.SourceSite = source
	mention.Platform = source.String()

	// Extract platform-specific handle
	switch source {
	case SourceBridgyBluesky:
		mention.Handle = extractBlueskyHandle(mention.Author.URL, mention.Source)
	case SourceBridgyTwitter:
		mention.Handle = extractTwitterHandle(mention.Author.URL, mention.Source)
	case SourceBridgyMastodon:
		mention.Handle = extractMastodonHandle(mention.Author.URL, mention.Source)
	case SourceBridgyGitHub:
		mention.Handle = extractGitHubHandle(mention.Author.URL, mention.Source)
	case SourceBridgyFlickr:
		mention.Handle = extractFlickrHandle(mention.Author.URL, mention.Source)
	case SourceWeb, SourceCustomBridge:
		// No handle extraction for generic web mentions or custom bridges
	}

	// Set original URL if different from source
	if mention.Platform != "web" && mention.URL != mention.Source {
		mention.OriginalURL = mention.Source
	}
}

// ShouldAccept checks if a mention should be accepted based on filters.
func (d *BridgingDetector) ShouldAccept(mention *ReceivedWebMention) bool {
	filters := d.config.Filters

	// Check platform filter
	if len(filters.Platforms) > 0 {
		found := false
		for _, p := range filters.Platforms {
			if strings.EqualFold(p, mention.Platform) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check interaction type filter
	if len(filters.InteractionTypes) > 0 {
		found := false
		interactionType := mention.InteractionType()
		for _, t := range filters.InteractionTypes {
			if strings.EqualFold(t, interactionType) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check content length
	if filters.MinContentLength > 0 {
		if len(mention.Content.Text) < filters.MinContentLength {
			return false
		}
	}

	// Check blocked domains
	if len(filters.BlockedDomains) > 0 {
		parsed, err := url.Parse(mention.Source)
		if err == nil {
			host := strings.ToLower(parsed.Host)
			for _, blocked := range filters.BlockedDomains {
				if strings.Contains(host, strings.ToLower(blocked)) {
					return false
				}
			}
		}
	}

	return true
}

// Handle extraction helpers

// blueskyHandleRegex matches Bluesky handles in URLs
var blueskyHandleRegex = regexp.MustCompile(`bsky\.(?:app|social)/profile/([a-zA-Z0-9._-]+(?:\.bsky\.social)?)`)

func extractBlueskyHandle(authorURL, sourceURL string) string {
	// Try author URL first
	if authorURL != "" {
		if match := blueskyHandleRegex.FindStringSubmatch(authorURL); len(match) > 1 {
			return "@" + match[1]
		}
	}
	// Fallback to source URL
	if sourceURL != "" {
		if match := blueskyHandleRegex.FindStringSubmatch(sourceURL); len(match) > 1 {
			return "@" + match[1]
		}
	}
	return ""
}

// twitterHandleRegex matches Twitter handles in URLs
var twitterHandleRegex = regexp.MustCompile(`(?:twitter\.com|x\.com)/(@?[a-zA-Z0-9_]+)`)

func extractTwitterHandle(authorURL, sourceURL string) string {
	// Try author URL first
	if authorURL != "" {
		if match := twitterHandleRegex.FindStringSubmatch(authorURL); len(match) > 1 {
			handle := match[1]
			if !strings.HasPrefix(handle, "@") {
				handle = "@" + handle
			}
			return handle
		}
	}
	// Fallback to source URL
	if match := twitterHandleRegex.FindStringSubmatch(sourceURL); len(match) > 1 {
		handle := match[1]
		if !strings.HasPrefix(handle, "@") {
			handle = "@" + handle
		}
		return handle
	}
	return ""
}

// mastodonHandleRegex matches Mastodon handles in URLs
var mastodonHandleRegex = regexp.MustCompile(`(?:/@|/users/)([a-zA-Z0-9_]+)(?:@([a-zA-Z0-9.-]+))?`)

func extractMastodonHandle(authorURL, sourceURL string) string {
	// Try to extract from URL pattern
	urlToCheck := authorURL
	if urlToCheck == "" {
		urlToCheck = sourceURL
	}

	if match := mastodonHandleRegex.FindStringSubmatch(urlToCheck); len(match) > 1 {
		handle := "@" + match[1]
		// Extract domain from URL if not in handle
		if len(match) > 2 && match[2] != "" {
			handle += "@" + match[2]
		} else {
			// Try to get domain from URL
			if parsed, err := url.Parse(urlToCheck); err == nil {
				handle += "@" + parsed.Host
			}
		}
		return handle
	}
	return ""
}

// githubHandleRegex matches GitHub handles in URLs
var githubHandleRegex = regexp.MustCompile(`github\.com/([a-zA-Z0-9_-]+)`)

func extractGitHubHandle(authorURL, sourceURL string) string {
	// Try author URL first
	if authorURL != "" {
		if match := githubHandleRegex.FindStringSubmatch(authorURL); len(match) > 1 {
			return "@" + match[1]
		}
	}
	// Fallback to source URL
	if match := githubHandleRegex.FindStringSubmatch(sourceURL); len(match) > 1 {
		return "@" + match[1]
	}
	return ""
}

// flickrHandleRegex matches Flickr handles in URLs
var flickrHandleRegex = regexp.MustCompile(`flickr\.com/(?:photos|people)/([a-zA-Z0-9@_-]+)`)

func extractFlickrHandle(authorURL, sourceURL string) string {
	// Try author URL first
	if authorURL != "" {
		if match := flickrHandleRegex.FindStringSubmatch(authorURL); len(match) > 1 {
			return "@" + match[1]
		}
	}
	// Fallback to source URL
	if match := flickrHandleRegex.FindStringSubmatch(sourceURL); len(match) > 1 {
		return "@" + match[1]
	}
	return ""
}

// PlatformColors returns CSS color values for each platform.
func PlatformColors() map[string]string {
	return map[string]string{
		"bluesky":  "#0085ff",
		"twitter":  "#1da1f2",
		"mastodon": "#6364ff",
		"github":   "#333333",
		"flickr":   "#ff0084",
		"web":      "#718096",
	}
}

// PlatformEmoji returns an emoji for each platform.
func PlatformEmoji() map[string]string {
	return map[string]string{
		"bluesky":  "\U0001F98B", // butterfly
		"twitter":  "\U0001F426", // bird
		"mastodon": "\U0001F418", // elephant
		"github":   "\U0001F419", // octopus
		"flickr":   "\U0001F4F7", // camera
		"web":      "\U0001F310", // globe
	}
}

// PlatformName returns a display name for each platform.
func PlatformName() map[string]string {
	return map[string]string{
		"bluesky":  "Bluesky",
		"twitter":  "Twitter",
		"mastodon": "Mastodon",
		"github":   "GitHub",
		"flickr":   "Flickr",
		"web":      "Web",
	}
}
