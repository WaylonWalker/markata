package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func testLobstersArticleServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write([]byte(`<html><head>
<meta property="og:title" content="Original Article Title">
<meta property="og:description" content="Original article summary">
<meta property="og:image" content="https://example.com/preview.jpg">
</head></html>`))
		if err != nil {
			t.Errorf("write article: %v", err)
		}
	}))
}

func TestLobstersRSSArticleAndDiscussion(t *testing.T) {
	article := testLobstersArticleServer(t)
	defer article.Close()

	discussion := "https://lobste.rs/s/abc123/example_story"
	// The official Lobsters RSS feed uses an article <link> and discussion <guid>.
	feedXML := []byte(fmt.Sprintf(`<rss version="2.0"><channel>
		<title>Lobsters</title>
		<item>
			<title>Lobsters headline</title>
			<link>%s</link>
			<guid isPermaLink="false">%s</guid>
			<pubDate>Thu, 15 Oct 2026 13:00:00 GMT</pubDate>
		</item>
	</channel></rss>`, article.URL, discussion))
	_, entries, err := parseRSS2Feed(feedXML)
	if err != nil {
		t.Fatalf("parseRSS2Feed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one entry, got %d", len(entries))
	}
	entry := entries[0]
	if entry.URL != article.URL || entry.OriginalURL != discussion {
		t.Fatalf("wrong article/discussion links: %+v", entry)
	}
	if entry.Title != "Original Article Title" || entry.Description != "Original article summary" ||
		entry.ImageURL != "https://example.com/preview.jpg" {
		t.Fatalf("article metadata not enriched: %+v", entry)
	}

	entry.FeedURL = "https://lobste.rs/rss"
	entry.FeedTitle = "Lobsters"
	plugin := NewBlogrollPlugin()
	entryMap := plugin.readerEntryMap(entry, nil, models.BlogrollConfig{})
	if got := entryMap["discussion_source"]; got != "Lobsters" {
		t.Fatalf("discussion_source = %v, want Lobsters", got)
	}
	groups := []map[string]interface{}{{
		"count":   1,
		"entries": []map[string]interface{}{entryMap},
	}}
	markup := plugin.renderReaderTimeline(groups, false)
	if !strings.Contains(markup, "Found on Lobsters:") ||
		!strings.Contains(markup, `href="https://lobste.rs/s/abc123/example_story"`) ||
		!strings.Contains(markup, `href="`+article.URL+`"`) {
		t.Fatalf("reader partial missing correct discussion/article links: %s", markup)
	}
}

func TestLobstersRSSCommentsField(t *testing.T) {
	article := testLobstersArticleServer(t)
	defer article.Close()

	// Some syndicated RSS feeds supply the discussion as <comments>.
	feedXML := []byte(fmt.Sprintf(`<rss><channel>
		<item><title>Story</title><link>%s</link>
			<guid>opaque-story-id</guid>
			<comments>https://lobste.rs/s/def456/comments</comments>
		</item>
	</channel></rss>`, article.URL))
	_, entries, err := parseRSS2Feed(feedXML)
	if err != nil {
		t.Fatalf("parseRSS2Feed: %v", err)
	}
	if len(entries) != 1 || entries[0].OriginalURL != "https://lobste.rs/s/def456/comments" {
		t.Fatalf("comments permalink not preserved: %+v", entries)
	}
}

func TestLobstersAtomEntryIDPreservesDiscussion(t *testing.T) {
	article := testLobstersArticleServer(t)
	defer article.Close()

	feedXML := []byte(fmt.Sprintf(`<feed xmlns="http://www.w3.org/2005/Atom">
		<title>Lobsters</title><entry>
			<id>https://lobste.rs/s/ghj789/story</id>
			<title>Feed headline</title>
			<link rel="alternate" href="%s"/>
		</entry></feed>`, article.URL))
	_, entries, err := parseAtomFeed(feedXML)
	if err != nil {
		t.Fatalf("parseAtomFeed: %v", err)
	}
	if len(entries) != 1 || entries[0].URL != article.URL ||
		entries[0].OriginalURL != "https://lobste.rs/s/ghj789/story" {
		t.Fatalf("atom outbound/discussion links incorrect: %+v", entries)
	}
}

func TestLobstersDiscussionOnlyFeedResolvesAPI(t *testing.T) {
	article := testLobstersArticleServer(t)
	defer article.Close()
	hits := 0
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		if r.URL.Path != "/s/xyz987.json" {
			t.Errorf("unexpected Lobsters API path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := fmt.Fprintf(w, `{"url":%q}`, article.URL); err != nil {
			t.Errorf("write article json: %v", err)
		}
	}))
	defer api.Close()

	oldBase, oldClient := lobstersStoryAPIBaseURL, lobstersHTTPClient
	lobstersStoryAPIBaseURL = api.URL + "/s/%s.json"
	lobstersHTTPClient = api.Client()
	defer func() {
		lobstersStoryAPIBaseURL, lobstersHTTPClient = oldBase, oldClient
	}()

	discussion := "https://lobste.rs/s/xyz987/article_slug"
	feedXML := []byte(`<rss><channel><item>
		<title>Submitted story</title><link>https://lobste.rs/s/xyz987/article_slug</link>
	</item></channel></rss>`)
	_, entries, err := parseRSS2Feed(feedXML)
	if err != nil {
		t.Fatalf("parseRSS2Feed: %v", err)
	}
	if len(entries) != 1 || entries[0].URL != article.URL ||
		entries[0].OriginalURL != discussion || entries[0].Title != "Original Article Title" {
		t.Fatalf("discussion-only feed not enriched: %+v", entries)
	}
	if hits != 1 {
		t.Fatalf("API hits = %d, want 1", hits)
	}
	if got := normalizeLobstersURL(discussion); got != article.URL || hits != 1 {
		t.Fatalf("cached URL = %q, API hits = %d", got, hits)
	}
}

func TestLobstersTextPostAndFailureStayOnDiscussion(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/s/tex123.json" {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte(`{"url":null}`)); err != nil {
				t.Errorf("write text story json: %v", err)
			}
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer api.Close()

	oldBase, oldClient := lobstersStoryAPIBaseURL, lobstersHTTPClient
	lobstersStoryAPIBaseURL = api.URL + "/s/%s.json"
	lobstersHTTPClient = api.Client()
	defer func() {
		lobstersStoryAPIBaseURL, lobstersHTTPClient = oldBase, oldClient
	}()
	for _, discussion := range []string{"https://lobste.rs/s/tex123/text_post", "https://lobste.rs/s/fai123/unavailable"} {
		resolved, original := resolveAggregatorFeedReference(discussion, "")
		if resolved != discussion || original != "" {
			t.Fatalf("text/error story (%s): resolved=%q original=%q", discussion, resolved, original)
		}
	}
}

func TestAggregatorDiscussionURLValidationAndLabels(t *testing.T) {
	testCases := []struct {
		url     string
		isStory bool
		label   string
	}{
		{"https://lobste.rs/s/abc123", true, "Lobsters"},
		{"https://www.lobste.rs/s/abc123/a-slug?foo=1", true, "Lobsters"},
		{"https://lobste.rs/t/python", false, ""},
		{"https://lobste.rs.evil.example/s/abc123", false, ""},
		{"javascript://lobste.rs/s/abc123", false, ""},
		{"https://lobste.rs/s/../../private", false, ""},
		{"https://news.ycombinator.com/item?id=123", false, "HN"},
	}
	for _, tc := range testCases {
		if got := isLobstersDiscussionURL(tc.url); got != tc.isStory {
			t.Errorf("isLobstersDiscussionURL(%q) = %v, want %v", tc.url, got, tc.isStory)
		}
		if got := aggregatorDiscussionLabel(tc.url); got != tc.label {
			t.Errorf("aggregatorDiscussionLabel(%q) = %q, want %q", tc.url, got, tc.label)
		}
	}
}

func TestHackerNewsDiscussionRemainsInReaderPartial(t *testing.T) {
	plugin := NewBlogrollPlugin()
	entry := &models.ExternalEntry{
		FeedURL:     "https://hnrss.org/frontpage",
		FeedTitle:   "Hacker News",
		Title:       "Article",
		URL:         "https://example.com/article",
		OriginalURL: "https://news.ycombinator.com/item?id=456",
	}
	entryMap := plugin.readerEntryMap(entry, nil, models.BlogrollConfig{})
	if got := entryMap["discussion_source"]; got != "HN" {
		t.Fatalf("discussion_source = %v, want HN", got)
	}
	markup := plugin.renderReaderTimeline([]map[string]interface{}{{
		"count":   1,
		"entries": []map[string]interface{}{entryMap},
	}}, false)
	if !strings.Contains(markup, "Found on HN:") ||
		!strings.Contains(markup, "https://news.ycombinator.com/item?id=456") {
		t.Fatalf("HN discussion missing from partial: %s", markup)
	}
}
