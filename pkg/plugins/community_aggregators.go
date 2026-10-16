package plugins

import (
	"html"
	"net/url"
	"regexp"
	"strings"
)

var (
	lemmyPostIDRegex    = regexp.MustCompile(`^\d+package plugins

import (
	"html"
	"net/url"
	"regexp"
	"strings"
)

var (
	lemmyPostIDRegex    = regexp.MustCompile()
	redditPostIDRegex   = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	redditFeedLinkRegex = regexp.MustCompile(`(?is)<a\b[^>]*href=["']([^"']+)["'][^>]*>\s*\[link\]\s*</a>`)
)

func isLemmyDiscussionURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != schemeHTTP && parsed.Scheme != schemeHTTPS) || parsed.Hostname() == "" {
		return false
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	return len(parts) == 2 && parts[0] == "post" && lemmyPostIDRegex.MatchString(parts[1]) //nolint:goconst // URL path segment, unrelated to garden node type
}

func isRedditDiscussionURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != schemeHTTP && parsed.Scheme != schemeHTTPS) {
		return false
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	host = strings.TrimPrefix(host, "old.")
	if host != "reddit.com" {
		return false
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) >= 4 && parts[0] == "r" && parts[2] == "comments" {
		return redditPostIDRegex.MatchString(parts[3])
	}
	return len(parts) >= 2 && parts[0] == "comments" && redditPostIDRegex.MatchString(parts[1])
}

// Reddit's temporary Atom feeds link to the discussion page. For external
// link posts the entry HTML commonly contains an explicitly labeled [link]
// anchor. Prefer only this anchor, not the first arbitrary URL in the post.
// This is best-effort; self, media and unrecognized posts stay on Reddit.
func redditOutboundArticleURL(discussionURL, feedContent string) string {
	if !isRedditDiscussionURL(discussionURL) {
		return ""
	}
	decoded := html.UnescapeString(feedContent)
	match := redditFeedLinkRegex.FindStringSubmatch(decoded)
	if len(match) < 2 {
		return ""
	}
	candidate := strings.TrimSpace(html.UnescapeString(match[1]))
	if !isHTTPArticleURL(candidate) || isRedditDiscussionURL(candidate) {
		return ""
	}
	parsed, err := url.Parse(candidate)
	if err != nil {
		return ""
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	if host == "reddit.com" || strings.HasSuffix(host, ".reddit.com") {
		return ""
	}
	return candidate
}
