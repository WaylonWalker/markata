package plugins

import (
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

var absoluteHTTPURLRegex = regexp.MustCompile(`http://[^\s"'<>)]*`)

// normalizeExternalURL upgrades safe absolute http URLs to https.
//
// This is intentionally conservative for local/private hosts, where forcing
// https would often be wrong. For public hosts, preferring https avoids mixed
// content when third-party metadata still advertises legacy http asset URLs.
func normalizeExternalURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.HasPrefix(raw, "//") {
		return "https:" + raw
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if parsed.Scheme != schemeHTTP || parsed.Hostname() == "" {
		return raw
	}
	if shouldKeepHTTP(parsed.Hostname()) {
		return raw
	}

	parsed.Scheme = schemeHTTPS
	return parsed.String()
}

func shouldKeepHTTP(host string) bool {
	host = strings.ToLower(strings.TrimSpace(host))
	if host == "" || host == "localhost" || strings.HasSuffix(host, ".local") {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

func normalizeExternalFeed(feed *models.ExternalFeed) {
	if feed == nil {
		return
	}
	feed.FeedURL = normalizeExternalURL(feed.FeedURL)
	feed.SiteURL = normalizeExternalURL(feed.SiteURL)
	feed.ImageURL = normalizeExternalURL(feed.ImageURL)
	feed.AvatarURL = normalizeExternalURL(feed.AvatarURL)
	for _, entry := range feed.Entries {
		normalizeExternalEntry(entry)
	}
}

func normalizeExternalEntry(entry *models.ExternalEntry) {
	if entry == nil {
		return
	}
	entry.FeedURL = normalizeExternalURL(entry.FeedURL)
	entry.URL = normalizeExternalURL(entry.URL)
	entry.ImageURL = normalizeExternalURL(entry.ImageURL)
}

func normalizeMentionMetadata(metadata *models.MentionMetadata) {
	if metadata == nil {
		return
	}
	metadata.URL = normalizeExternalURL(metadata.URL)
	metadata.Avatar = normalizeExternalURL(metadata.Avatar)
}

func normalizeOGMetadata(metadata *OGMetadata) {
	if metadata == nil {
		return
	}
	metadata.Image = normalizeExternalURL(metadata.Image)
	metadata.AuthorURL = normalizeExternalURL(metadata.AuthorURL)
}

func normalizeHTMLFragmentURLs(fragment string) string {
	if fragment == "" {
		return ""
	}
	return absoluteHTTPURLRegex.ReplaceAllStringFunc(fragment, normalizeExternalURL)
}

func normalizeReceivedWebMention(mention *ReceivedWebMention) {
	if mention == nil {
		return
	}
	mention.URL = normalizeExternalURL(mention.URL)
	mention.Source = normalizeExternalURL(mention.Source)
	mention.Target = normalizeExternalURL(mention.Target)
	mention.OriginalURL = normalizeExternalURL(mention.OriginalURL)
	mention.Author.URL = normalizeExternalURL(mention.Author.URL)
	mention.Author.Photo = normalizeExternalURL(mention.Author.Photo)
	mention.Content.Text = normalizeHTMLFragmentURLs(mention.Content.Text)
	mention.Content.HTML = normalizeHTMLFragmentURLs(mention.Content.HTML)
}
