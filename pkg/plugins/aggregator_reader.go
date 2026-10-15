package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

var (
	lobstersStoryAPIBaseURL = "https://lobste.rs/s/%s.json"
	lobstersHTTPClient      = instrumentedHackerNewsHTTPClient()
	lobstersURLCache        sync.Map
	lobstersShortIDRegex    = regexp.MustCompile(`^[a-zA-Z0-9]{6}$`)
)

type lobstersStory struct {
	URL string `json:"url"`
}

// isLobstersDiscussionURL recognizes story permalinks, with or without a
// human-readable slug, rather than arbitrary URLs on the Lobsters domain.
func isLobstersDiscussionURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	host := strings.TrimPrefix(strings.ToLower(parsed.Hostname()), "www.")
	if host != "lobste.rs" {
		return false
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	return (len(parts) == 2 || len(parts) == 3) &&
		parts[0] == "s" && lobstersShortIDRegex.MatchString(parts[1])
}

func isHTTPArticleURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	return err == nil &&
		(parsed.Scheme == "http" || parsed.Scheme == "https") &&
		parsed.Hostname() != ""
}

// normalizeLobstersURL handles feeds which link to a Lobsters discussion
// instead of to the article. Lobsters' official RSS usually already links
// directly to the article, with the discussion permalink in the GUID.
func normalizeLobstersURL(rawURL string) string {
	if !isLobstersDiscussionURL(rawURL) {
		return rawURL
	}
	if cached, ok := lobstersURLCache.Load(rawURL); ok {
		return cached.(string)
	}

	parsed, _ := url.Parse(rawURL) // already validated
	shortID := strings.Split(strings.Trim(parsed.Path, "/"), "/")[1]
	apiURL := fmt.Sprintf(lobstersStoryAPIBaseURL, shortID)
	req, err := http.NewRequestWithContext(
		context.Background(), http.MethodGet, apiURL, http.NoBody,
	)
	if err != nil {
		return rawURL
	}
	resp, err := lobstersHTTPClient.Do(req)
	if err != nil {
		return rawURL
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return rawURL
	}
	var story lobstersStory
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1024*1024)).Decode(&story); err != nil {
		return rawURL
	}
	articleURL := strings.TrimSpace(story.URL)
	if !isHTTPArticleURL(articleURL) {
		// Text-only stories have no outbound URL.
		lobstersURLCache.Store(rawURL, rawURL)
		return rawURL
	}
	lobstersURLCache.Store(rawURL, articleURL)
	return articleURL
}

// resolveAggregatorFeedReference keeps the community discussion URL alongside
// the article, including feeds whose link already points at the article.
func resolveAggregatorFeedReference(link, discussionHint string) (resolved, original string) {
	resolved, original = resolveHackerNewsReference(link)
	if original != "" {
		return resolved, original
	}
	if isLobstersDiscussionURL(link) {
		resolved = normalizeLobstersURL(link)
		if resolved != link {
			return resolved, link
		}
		return resolved, ""
	}
	if isLobstersDiscussionURL(discussionHint) {
		if isHTTPArticleURL(link) && link != discussionHint {
			return link, discussionHint
		}
		if link == "" {
			return discussionHint, ""
		}
	}
	return resolved, ""
}

func enrichAggregatorFeedEntry(link, discussionHint string) (resolved, original string, metadata *hackerNewsMetadata) {
	resolved, original = resolveAggregatorFeedReference(link, discussionHint)
	if original == "" {
		return resolved, original, nil
	}
	metadata, _ = fetchHackerNewsArticleMetadata(resolved)
	return resolved, original, metadata
}

func aggregatorDiscussionLabel(rawURL string) string {
	switch {
	case isHackerNewsDiscussionURL(rawURL):
		return "HN"
	case isLobstersDiscussionURL(rawURL):
		return "Lobsters"
	default:
		return ""
	}
}
