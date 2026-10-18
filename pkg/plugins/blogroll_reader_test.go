package plugins

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestBlogrollPlugin_ReaderDayGroups(t *testing.T) {
	plugin := NewBlogrollPlugin()
	base := time.Date(2026, 4, 17, 9, 0, 0, 0, time.UTC)
	config := models.BlogrollConfig{FallbackImageService: "https://shots.example.com/?url={url}"}

	entries := []*models.ExternalEntry{
		{
			FeedURL:   "https://example.com/feed.xml",
			FeedTitle: "Example",
			Title:     "Morning post",
			Published: &base,
		},
		{
			FeedURL:   "https://example.com/feed.xml",
			FeedTitle: "Example",
			Title:     "Later post",
			Published: func() *time.Time { t := base.Add(2 * time.Hour); return &t }(),
		},
		{
			FeedURL:   "https://other.example/rss",
			FeedTitle: "Other",
			Title:     "Yesterday post",
			Published: func() *time.Time { t := base.Add(-24 * time.Hour); return &t }(),
		},
	}
	feeds := []*models.ExternalFeed{{
		FeedURL:   "https://example.com/feed.xml",
		Title:     "Example",
		SiteURL:   "https://example.com",
		AvatarURL: "https://example.com/avatar.png",
	}}

	groups := plugin.readerDayGroups(entries, buildFeedIndex(feeds), config)
	if got, want := len(groups), 2; got != want {
		t.Fatalf("len(groups) = %d, want %d", got, want)
	}

	first := groups[0]
	if got, want := first["date_key"], "2026-04-17"; got != want {
		t.Fatalf("first date_key = %v, want %v", got, want)
	}
	if got, want := first["count"], 2; got != want {
		t.Fatalf("first count = %v, want %v", got, want)
	}

	groupEntries, ok := first["entries"].([]map[string]interface{})
	if !ok {
		t.Fatalf("first entries has type %T, want []map[string]interface{}", first["entries"])
	}
	if got, want := len(groupEntries), 2; got != want {
		t.Fatalf("len(first entries) = %d, want %d", got, want)
	}

	iconURL, ok := groupEntries[0]["source_icon_url"].(string)
	if !ok || iconURL == "" || !strings.Contains(iconURL, "example.com") {
		t.Fatalf("source_icon_url = %q, want favicon URL for example.com", iconURL)
	}

	if got, want := groupEntries[0]["published_label"], "Apr 17, 2026"; got != want {
		t.Fatalf("published_label = %v, want %v", got, want)
	}

	if got, want := groupEntries[0]["preview_kind"], "source-image"; got != want {
		t.Fatalf("preview_kind = %v, want %v", got, want)
	}

	markup := plugin.renderReaderTimeline(groups, false)
	if !strings.Contains(markup, `class="reader-stream posts-list"`) {
		t.Fatalf("timeline markup missing wrapper: %s", markup)
	}
	if !strings.Contains(markup, `class="reader-day"`) || !strings.Contains(markup, `class="reader-entry-source-icon"`) {
		t.Fatalf("timeline markup missing expected reader classes: %s", markup)
	}
}

func TestBuildReaderVariants(t *testing.T) {
	config := models.BlogrollConfig{ReaderSlug: "reader"}
	feeds := []*models.ExternalFeed{
		{FeedURL: "https://example.com/written.xml", Config: models.ExternalFeedConfig{Type: models.ReaderFeedTypeWritten}, Type: models.ReaderFeedTypeWritten},
		{FeedURL: "https://example.com/video.xml", Config: models.ExternalFeedConfig{Type: models.ReaderFeedTypeVideo}, Type: models.ReaderFeedTypeVideo},
	}
	entries := []*models.ExternalEntry{
		{FeedURL: "https://example.com/written.xml", Title: "Post"},
		{FeedURL: "https://example.com/video.xml", Title: "Video"},
	}
	variants := buildReaderVariants(feeds, entries, config, buildFeedIndex(feeds), t.TempDir())
	if len(variants) != 3 {
		t.Fatalf("len(variants) = %d, want 3", len(variants))
	}
	if variants[0].Key != "all" || variants[1].Key != "written" || variants[2].Key != "video" {
		t.Fatalf("variant keys = %v, want [all written video]", []string{variants[0].Key, variants[1].Key, variants[2].Key})
	}
	if variants[2].BaseURL != "/reader/video" {
		t.Fatalf("video BaseURL = %q, want %q", variants[2].BaseURL, "/reader/video")
	}
}

func TestRenderReaderTypeLinks(t *testing.T) {
	markup := renderReaderTypeLinks([]map[string]interface{}{
		{"title": "Reader", "url": "/reader/", "active": true, "entry_count": 10},
		{"title": "Reader: Video", "url": "/reader/video/", "active": false, "entry_count": 3},
	})
	if !strings.Contains(markup, `Reader (10)`) {
		t.Fatalf("markup missing all-reader count: %s", markup)
	}
	if !strings.Contains(markup, `href="/reader/video/"`) {
		t.Fatalf("markup missing video link: %s", markup)
	}
	if !strings.Contains(markup, `class="active"`) {
		t.Fatalf("markup missing active reader type: %s", markup)
	}
}

func TestReaderPreviewForEntry_Hierarchy(t *testing.T) {
	tests := []struct {
		name            string
		entry           *models.ExternalEntry
		sourceImageURL  string
		fallbackService string
		wantURL         string
		wantKind        string
	}{
		{
			name:            "article image wins",
			entry:           &models.ExternalEntry{URL: "https://example.com/post", ImageURL: "https://cdn.example.com/post.jpg"},
			sourceImageURL:  "https://cdn.example.com/source.png",
			fallbackService: "https://shots.example.com/?url={url}",
			wantURL:         "https://cdn.example.com/post.jpg",
			wantKind:        "article-image",
		},
		{
			name:            "source image beats screenshot",
			entry:           &models.ExternalEntry{URL: "https://example.com/post"},
			sourceImageURL:  "https://cdn.example.com/source.png",
			fallbackService: "https://shots.example.com/?url={url}",
			wantURL:         "https://cdn.example.com/source.png",
			wantKind:        "source-image",
		},
		{
			name:            "screenshot before source tile",
			entry:           &models.ExternalEntry{URL: "https://example.com/post?id=1"},
			fallbackService: "https://shots.example.com/?url={url}",
			wantURL:         "https://shots.example.com/?url=https%3A%2F%2Fexample.com%2Fpost%3Fid%3D1",
			wantKind:        "screenshot",
		},
		{
			name:     "source tile last",
			entry:    &models.ExternalEntry{URL: "https://example.com/post"},
			wantKind: "source-tile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotKind := readerPreviewForEntry(tt.entry, tt.sourceImageURL, tt.fallbackService)
			if gotURL != tt.wantURL || gotKind != tt.wantKind {
				t.Fatalf("readerPreviewForEntry() = (%q, %q), want (%q, %q)", gotURL, gotKind, tt.wantURL, tt.wantKind)
			}
		})
	}
}

func TestNormalizeHackerNewsURL_ResolvesArticleURL(t *testing.T) {
	articleServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Article Title</title>
	<meta property="og:title" content="Article Title">
	<meta property="og:description" content="Article Description">
	<meta property="og:image" content="https://example.com/article.jpg">
</head>
<body></body>
</html>`)); err != nil {
			t.Errorf("write article html: %v", err)
		}
	}))
	defer articleServer.Close()
	articleURL := articleServer.URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v0/item/123.json" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"url":"` + articleURL + `"}`)); err != nil {
			t.Errorf("write HN json: %v", err)
		}
	}))
	defer server.Close()

	oldBaseURL := hackerNewsItemAPIBaseURL
	oldClient := hackerNewsHTTPClient
	hackerNewsItemAPIBaseURL = server.URL + "/v0/item/%s.json"
	hackerNewsHTTPClient = server.Client()
	defer func() {
		hackerNewsItemAPIBaseURL = oldBaseURL
		hackerNewsHTTPClient = oldClient
	}()

	rawURL := "https://news.ycombinator.com/item?id=123"
	if got := normalizeHackerNewsURL(rawURL); got != articleURL {
		t.Fatalf("normalizeHackerNewsURL() = %q, want %q", got, articleURL)
	}

	feedXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Hacker News</title>
  <entry>
    <id>123</id>
    <title>Story title</title>
    <published>2026-04-17T12:17:00Z</published>
    <link rel="alternate" href="` + rawURL + `"/>
  </entry>
</feed>`)

	_, entries, err := parseAtomFeed(feedXML)
	if err != nil {
		t.Fatalf("parseAtomFeed() error = %v", err)
	}
	if got, want := len(entries), 1; got != want {
		t.Fatalf("len(entries) = %d, want %d", got, want)
	}
	if got, want := entries[0].URL, articleURL; got != want {
		t.Fatalf("entry.URL = %q, want %q", got, want)
	}
	if got, want := entries[0].Title, "Article Title"; got != want {
		t.Fatalf("entry.Title = %q, want %q", got, want)
	}
	if got, want := entries[0].Description, "Article Description"; got != want {
		t.Fatalf("entry.Description = %q, want %q", got, want)
	}
	if got, want := entries[0].ImageURL, "https://example.com/article.jpg"; got != want {
		t.Fatalf("entry.ImageURL = %q, want %q", got, want)
	}
	if got, want := entries[0].OriginalURL, rawURL; got != want {
		t.Fatalf("entry.OriginalURL = %q, want %q", got, want)
	}
}

func TestParseAtomFeed_YouTubeEntriesUseVideoThumbnail(t *testing.T) {
	feedXML := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:yt="http://www.youtube.com/xml/schemas/2015">
  <title>YouTube Channel</title>
  <link rel="alternate" href="https://www.youtube.com/channel/abc"/>
  <entry>
    <id>yt:video:Oq5e_8zvick</id>
    <title>It's all fake</title>
    <published>2026-04-17T12:17:00Z</published>
    <updated>2026-04-17T18:26:19Z</updated>
    <link rel="alternate" href="https://www.youtube.com/watch?v=Oq5e_8zvick"/>
    <author><name>The PrimeTime</name></author>
  </entry>
</feed>`)

	_, entries, err := parseAtomFeed(feedXML)
	if err != nil {
		t.Fatalf("parseAtomFeed() error = %v", err)
	}
	if got, want := len(entries), 1; got != want {
		t.Fatalf("len(entries) = %d, want %d", got, want)
	}
	if got, want := entries[0].ImageURL, "https://img.youtube.com/vi/Oq5e_8zvick/hqdefault.jpg"; got != want {
		t.Fatalf("entry.ImageURL = %q, want %q", got, want)
	}
}

func TestReaderCSS_DuckDBOverflowRegressionRules(t *testing.T) {
	cssPath := filepath.Join("..", "themes", "default", "static", "css", "components.css")
	content, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("read css: %v", err)
	}
	css := string(content)

	checks := []string{
		".reader-day-entries > *",
		"min-width: 0;",
		".reader-entry-meta-row",
		"grid-template-columns: minmax(0, 1fr) auto;",
		".reader-entry-source-link",
		"grid-template-columns: auto minmax(0, 1fr);",
		".reader-entry-source",
		"overflow-wrap: anywhere;",
		"word-break: break-word;",
	}

	for _, check := range checks {
		if !strings.Contains(css, check) {
			t.Fatalf("components.css missing regression rule fragment %q", check)
		}
	}
}
