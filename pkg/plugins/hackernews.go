package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	neturl "net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
)

var hackerNewsItemAPIBaseURL = "https://hacker-news.firebaseio.com/v0/item/%s.json"

var hackerNewsHTTPClient = instrumentedHackerNewsHTTPClient()

func instrumentedHackerNewsHTTPClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	buildstats.InstrumentHTTPClient(client)
	return client
}

var hackerNewsURLCache sync.Map

var (
	hackerNewsHTMLTitleRegex = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)
	hackerNewsOGTitleRegex   = regexp.MustCompile(`(?is)<meta[^>]+property=['"]og:title['"][^>]+content=['"]([^'"]+)['"]|<meta[^>]+content=['"]([^'"]+)['"][^>]+property=['"]og:title['"]`)
	hackerNewsOGDescRegex    = regexp.MustCompile(`(?is)<meta[^>]+property=['"]og:description['"][^>]+content=['"]([^'"]+)['"]|<meta[^>]+content=['"]([^'"]+)['"][^>]+property=['"]og:description['"]`)
	hackerNewsOGImageRegex   = regexp.MustCompile(`(?is)<meta[^>]+property=['"]og:image['"][^>]+content=['"]([^'"]+)['"]|<meta[^>]+content=['"]([^'"]+)['"][^>]+property=['"]og:image['"]`)
	hackerNewsDescRegex      = regexp.MustCompile(`(?is)<meta[^>]+name=['"]description['"][^>]+content=['"]([^'"]+)['"]|<meta[^>]+content=['"]([^'"]+)['"][^>]+name=['"]description['"]`)
)

type hackerNewsItem struct {
	URL string `json:"url"`
}

type hackerNewsMetadata struct {
	Title       string
	Description string
	Image       string
}

// normalizeHackerNewsURL resolves Hacker News discussion URLs to their outbound article URLs.
// If resolution fails or the item has no outbound URL, the original URL is returned.
func normalizeHackerNewsURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	if cached, ok := hackerNewsURLCache.Load(rawURL); ok {
		if resolved, ok := cached.(string); ok && resolved != "" {
			return resolved
		}
		return rawURL
	}

	if !isHackerNewsDiscussionURL(rawURL) {
		return rawURL
	}

	resolvedURL, ok := fetchHackerNewsArticleURL(rawURL)
	if !ok {
		return rawURL
	}
	if resolvedURL == "" {
		hackerNewsURLCache.Store(rawURL, "")
		return rawURL
	}

	hackerNewsURLCache.Store(rawURL, resolvedURL)
	return resolvedURL
}

func resolveHackerNewsReference(rawURL string) (resolvedURL, originalURL string) {
	if !isHackerNewsDiscussionURL(rawURL) {
		return rawURL, ""
	}

	resolvedURL = normalizeHackerNewsURL(rawURL)
	if resolvedURL == "" {
		return rawURL, ""
	}
	if resolvedURL != rawURL {
		return resolvedURL, rawURL
	}
	return rawURL, ""
}

func isHackerNewsDiscussionURL(rawURL string) bool {
	parsed, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}

	host := strings.ToLower(parsed.Hostname())
	host = strings.TrimPrefix(host, "www.")
	if host != "news.ycombinator.com" {
		return false
	}

	if !strings.EqualFold(parsed.Path, "/item") {
		return false
	}

	return strings.TrimSpace(parsed.Query().Get("id")) != ""
}

func fetchHackerNewsArticleURL(rawURL string) (string, bool) {
	parsed, err := neturl.Parse(rawURL)
	if err != nil {
		return "", false
	}

	id := strings.TrimSpace(parsed.Query().Get("id"))
	if id == "" {
		return "", false
	}

	apiURL := fmt.Sprintf(hackerNewsItemAPIBaseURL, id)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, http.NoBody)
	if err != nil {
		return "", false
	}

	resp, err := hackerNewsHTTPClient.Do(req)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false
	}

	var item hackerNewsItem
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return "", false
	}

	return strings.TrimSpace(item.URL), true
}

func fetchHackerNewsArticleMetadata(articleURL string) *hackerNewsMetadata {
	if articleURL == "" {
		return nil
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, articleURL, http.NoBody)
	if err != nil {
		return nil
	}
	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")

	resp, err := hackerNewsHTTPClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	if err != nil {
		return nil
	}

	htmlContent := string(body)
	metadata := &hackerNewsMetadata{
		Title:       extractHackerNewsOGTitle(htmlContent),
		Description: extractHackerNewsOGDescription(htmlContent),
		Image:       extractHackerNewsOGImage(htmlContent),
	}
	if metadata.Title == "" {
		metadata.Title = extractHackerNewsHTMLTitle(htmlContent)
	}
	if metadata.Description == "" {
		metadata.Description = extractHackerNewsDescription(htmlContent)
	}

	return metadata
}

func enrichHackerNewsFeedEntry(rawURL string) (resolvedURL, originalURL string, metadata *hackerNewsMetadata) {
	resolvedURL, originalURL = resolveHackerNewsReference(rawURL)
	if originalURL == "" {
		return resolvedURL, originalURL, nil
	}

	metadata = fetchHackerNewsArticleMetadata(resolvedURL)
	return resolvedURL, originalURL, metadata
}

func extractHackerNewsOGTitle(htmlContent string) string {
	return extractHackerNewsMetaMatch(htmlContent, hackerNewsOGTitleRegex)
}

func extractHackerNewsOGDescription(htmlContent string) string {
	return extractHackerNewsMetaMatch(htmlContent, hackerNewsOGDescRegex)
}

func extractHackerNewsOGImage(htmlContent string) string {
	return extractHackerNewsMetaMatch(htmlContent, hackerNewsOGImageRegex)
}

func extractHackerNewsDescription(htmlContent string) string {
	return extractHackerNewsMetaMatch(htmlContent, hackerNewsDescRegex)
}

func extractHackerNewsMetaMatch(htmlContent string, pattern *regexp.Regexp) string {
	if match := pattern.FindStringSubmatch(htmlContent); len(match) > 1 {
		for _, candidate := range match[1:] {
			if candidate != "" {
				return html.UnescapeString(candidate)
			}
		}
	}
	return ""
}

func extractHackerNewsHTMLTitle(htmlContent string) string {
	if match := hackerNewsHTMLTitleRegex.FindStringSubmatch(htmlContent); len(match) > 1 {
		return strings.TrimSpace(html.UnescapeString(match[1]))
	}
	return ""
}
