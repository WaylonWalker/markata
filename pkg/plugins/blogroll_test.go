package plugins

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

// =============================================================================
// extractFirstImageFromHTML Tests
// =============================================================================

func TestExtractFirstImageFromHTML_BasicImage(t *testing.T) {
	html := `<p>Some text</p><img src="https://example.com/image.jpg" alt="test">`
	result := extractFirstImageFromHTML(html)
	if result != "https://example.com/image.jpg" {
		t.Errorf("extractFirstImageFromHTML() = %q, want %q", result, "https://example.com/image.jpg")
	}
}

func TestExtractFirstImageFromHTML_SingleQuotes(t *testing.T) {
	html := `<img src='https://example.com/image.png'>`
	result := extractFirstImageFromHTML(html)
	if result != "https://example.com/image.png" {
		t.Errorf("extractFirstImageFromHTML() = %q, want %q", result, "https://example.com/image.png")
	}
}

func TestExtractFirstImageFromHTML_NoImage(t *testing.T) {
	html := `<p>No images here</p>`
	result := extractFirstImageFromHTML(html)
	if result != "" {
		t.Errorf("extractFirstImageFromHTML() = %q, want empty string", result)
	}
}

func TestExtractFirstImageFromHTML_HTMLEntities(t *testing.T) {
	// This tests the bug fix for Atom feeds that encode content as HTML entities
	html := `&lt;p&gt;Some text&lt;/p&gt;&lt;img src="https://example.com/atom-image.jpg" alt="test"&gt;`
	result := extractFirstImageFromHTML(html)
	if result != "https://example.com/atom-image.jpg" {
		t.Errorf("extractFirstImageFromHTML() with HTML entities = %q, want %q", result, "https://example.com/atom-image.jpg")
	}
}

func TestExtractFirstImageFromHTML_NestedHTMLEntities(t *testing.T) {
	// Test deeply encoded content (edge case)
	html := `&lt;img src=&quot;https://example.com/deeply-encoded.jpg&quot;&gt;`
	result := extractFirstImageFromHTML(html)
	if result != "https://example.com/deeply-encoded.jpg" {
		t.Errorf("extractFirstImageFromHTML() with nested entities = %q, want %q", result, "https://example.com/deeply-encoded.jpg")
	}
}

func TestExtractFirstImageFromHTML_MultipleImages(t *testing.T) {
	// Should return the first image only
	html := `<img src="first.jpg"><img src="second.jpg">`
	result := extractFirstImageFromHTML(html)
	if result != "first.jpg" {
		t.Errorf("extractFirstImageFromHTML() = %q, want %q", result, "first.jpg")
	}
}

// =============================================================================
// generateFallbackImageURL Tests
// =============================================================================

func TestGenerateFallbackImageURL_Basic(t *testing.T) {
	template := "https://shots.example.com/shot/?url={url}&width=1200"
	entryURL := "https://blog.example.com/my-post"
	result := generateFallbackImageURL(template, entryURL)
	expected := "https://shots.example.com/shot/?url=https%3A%2F%2Fblog.example.com%2Fmy-post&width=1200"
	if result != expected {
		t.Errorf("generateFallbackImageURL() = %q, want %q", result, expected)
	}
}

func TestGenerateFallbackImageURL_WithSpecialChars(t *testing.T) {
	template := "https://screenshot.service/{url}"
	entryURL := "https://example.com/post?foo=bar&baz=qux"
	result := generateFallbackImageURL(template, entryURL)
	expected := "https://screenshot.service/https%3A%2F%2Fexample.com%2Fpost%3Ffoo%3Dbar%26baz%3Dqux"
	if result != expected {
		t.Errorf("generateFallbackImageURL() = %q, want %q", result, expected)
	}
}

func TestGenerateFallbackImageURL_NoPlaceholder(t *testing.T) {
	// If template has no {url} placeholder, return as-is
	template := "https://default-image.com/fallback.png"
	entryURL := "https://blog.example.com/post"
	result := generateFallbackImageURL(template, entryURL)
	if result != template {
		t.Errorf("generateFallbackImageURL() = %q, want %q", result, template)
	}
}

func TestGenerateFallbackImageURL_EmptyEntryURL(t *testing.T) {
	template := "https://shots.example.com/?url={url}"
	entryURL := ""
	result := generateFallbackImageURL(template, entryURL)
	expected := "https://shots.example.com/?url="
	if result != expected {
		t.Errorf("generateFallbackImageURL() = %q, want %q", result, expected)
	}
}

func TestGenerateFallbackImageURL_URLWithUnicode(t *testing.T) {
	template := "https://shots.example.com/?url={url}"
	entryURL := "https://example.com/post/日本語"
	result := generateFallbackImageURL(template, entryURL)
	// Unicode characters should be percent-encoded
	expected := "https://shots.example.com/?url=https%3A%2F%2Fexample.com%2Fpost%2F%E6%97%A5%E6%9C%AC%E8%AA%9E"
	if result != expected {
		t.Errorf("generateFallbackImageURL() = %q, want %q", result, expected)
	}
}

func TestBlogrollPlugin_LoadFromCache_UsesLastFetchedTimestamp(t *testing.T) {
	tmpDir := t.TempDir()
	p := NewBlogrollPlugin()
	now := time.Now()
	feed := &models.ExternalFeed{
		FeedURL:     "https://example.com/feed.xml",
		Title:       "Example",
		LastFetched: &now,
	}
	p.saveToCache(feed, tmpDir)

	cachePath := filepath.Join(tmpDir, p.cacheKey(feed.FeedURL)+".json")
	oldTime := now.Add(-48 * time.Hour)
	if err := os.Chtimes(cachePath, oldTime, oldTime); err != nil {
		t.Fatalf("Chtimes() error = %v", err)
	}

	cached := p.loadFromCache(feed.FeedURL, tmpDir, 24*time.Hour)
	if cached == nil {
		t.Fatal("loadFromCache() = nil, want cached feed")
	}
	if cached.Title != "Example" {
		t.Fatalf("cached.Title = %q, want Example", cached.Title)
	}
}

func TestBlogrollPlugin_AggregateCache_RoundTripAndHashCheck(t *testing.T) {
	tmpDir := t.TempDir()
	p := NewBlogrollPlugin()

	feeds := []*models.ExternalFeed{{
		FeedURL: "https://example.com/feed.xml",
		Title:   "Example",
	}}
	entries := []*models.ExternalEntry{{
		FeedURL: "https://example.com/feed.xml",
		ID:      "entry-1",
		Title:   "First Entry",
	}}

	p.saveAggregateCache(tmpDir, "hash-a", feeds, entries)

	cached := p.loadAggregateCache(tmpDir, 24*time.Hour, "hash-a")
	if cached == nil {
		t.Fatal("loadAggregateCache() = nil, want cached snapshot")
	}
	if len(cached.Feeds) != 1 || cached.Feeds[0].Title != "Example" {
		t.Fatalf("cached.Feeds = %#v, want one Example feed", cached.Feeds)
	}
	if len(cached.Entries) != 1 || cached.Entries[0].Title != "First Entry" {
		t.Fatalf("cached.Entries = %#v, want one First Entry", cached.Entries)
	}

	if got := p.loadAggregateCache(tmpDir, 24*time.Hour, "hash-b"); got != nil {
		t.Fatalf("loadAggregateCache() with mismatched hash = %#v, want nil", got)
	}
}

func TestBlogrollPlugin_LoadAggregateCacheAny_IgnoresAgeButChecksHash(t *testing.T) {
	tmpDir := t.TempDir()
	p := NewBlogrollPlugin()

	feeds := []*models.ExternalFeed{{
		FeedURL: "https://example.com/feed.xml",
		Title:   "Example",
	}}
	entries := []*models.ExternalEntry{{
		FeedURL: "https://example.com/feed.xml",
		ID:      "entry-1",
		Title:   "First Entry",
	}}

	p.saveAggregateCache(tmpDir, "hash-a", feeds, entries)
	cachePath := filepath.Join(tmpDir, aggregateCacheFile)
	oldTime := time.Now().Add(-7 * 24 * time.Hour)
	if err := os.Chtimes(cachePath, oldTime, oldTime); err != nil {
		t.Fatalf("Chtimes() error = %v", err)
	}

	cached := p.loadAggregateCacheAny(tmpDir, "hash-a")
	if cached == nil {
		t.Fatal("loadAggregateCacheAny() = nil, want cached snapshot")
	}
	if got := p.loadAggregateCacheAny(tmpDir, "hash-b"); got != nil {
		t.Fatalf("loadAggregateCacheAny() with mismatched hash = %#v, want nil", got)
	}
}

func TestBlogrollPlugin_FetchFeed_UsesStaleCacheWhenRefreshDisabled(t *testing.T) {
	tmpDir := t.TempDir()
	p := NewBlogrollPlugin()
	oldTime := time.Now().Add(-48 * time.Hour)
	feed := &models.ExternalFeed{
		FeedURL:     "https://example.com/feed.xml",
		Title:       "Cached",
		LastFetched: &oldTime,
	}
	p.saveToCache(feed, tmpDir)

	config := models.ExternalFeedConfig{URL: "https://example.com/feed.xml", Title: "Configured"}
	got, summary := p.fetchFeed(config, tmpDir, time.Hour, 1, 10, blogrollFetchOptions{refreshOnBuild: false})

	if got == nil {
		t.Fatal("fetchFeed() = nil, want cached feed")
	}
	if got.Title != "Configured" {
		t.Fatalf("fetchFeed().Title = %q, want config override", got.Title)
	}
	if summary.stale != 1 {
		t.Fatalf("summary.stale = %d, want 1", summary.stale)
	}
	if summary.refreshed != 0 {
		t.Fatalf("summary.refreshed = %d, want 0", summary.refreshed)
	}
}

func TestBlogrollPlugin_FetchFeed_ReturnsPlaceholderWhenRefreshDisabledAndCacheMissing(t *testing.T) {
	p := NewBlogrollPlugin()
	config := models.ExternalFeedConfig{URL: "https://example.com/feed.xml", Title: "Configured"}

	got, summary := p.fetchFeed(config, t.TempDir(), time.Hour, 1, 10, blogrollFetchOptions{refreshOnBuild: false})

	if got == nil {
		t.Fatal("fetchFeed() = nil, want placeholder feed")
	}
	if got.Error == "" {
		t.Fatal("fetchFeed().Error = empty, want missing cache message")
	}
	if summary.failed != 1 {
		t.Fatalf("summary.failed = %d, want 1", summary.failed)
	}
}

func TestMergeCachedFeed_AppliesConfigOverrides(t *testing.T) {
	cached := &models.ExternalFeed{
		Config: models.ExternalFeedConfig{
			URL:      "https://example.com/feed.xml",
			Handle:   "old-handle",
			Aliases:  []string{"old"},
			SiteURL:  "http://example.com",
			ImageURL: "http://example.com/old.png",
		},
		Title:       "Cached Title",
		Description: "Cached description",
		SiteURL:     "http://example.com",
		ImageURL:    "http://example.com/old.png",
		Category:    "Old",
		Tags:        []string{"cached"},
	}
	config := models.ExternalFeedConfig{
		URL:         "https://example.com/feed.xml",
		Title:       "Config Title",
		Description: "Config description",
		Category:    "Config",
		Tags:        []string{"fresh"},
		SiteURL:     "https://example.com",
		ImageURL:    "https://example.com/new.png",
		Handle:      "new-handle",
		Aliases:     []string{"new"},
	}

	merged := mergeCachedFeed(cached, config)

	if merged.Config.Handle != "new-handle" {
		t.Fatalf("merged.Config.Handle = %q, want %q", merged.Config.Handle, "new-handle")
	}
	if merged.Title != "Config Title" {
		t.Fatalf("merged.Title = %q, want %q", merged.Title, "Config Title")
	}
	if merged.Description != "Config description" {
		t.Fatalf("merged.Description = %q, want %q", merged.Description, "Config description")
	}
	if merged.Category != "Config" {
		t.Fatalf("merged.Category = %q, want %q", merged.Category, "Config")
	}
	if merged.SiteURL != "https://example.com" {
		t.Fatalf("merged.SiteURL = %q, want %q", merged.SiteURL, "https://example.com")
	}
	if merged.ImageURL != "https://example.com/new.png" {
		t.Fatalf("merged.ImageURL = %q, want %q", merged.ImageURL, "https://example.com/new.png")
	}
	if len(merged.Tags) != 1 || merged.Tags[0] != "fresh" {
		t.Fatalf("merged.Tags = %#v, want %#v", merged.Tags, []string{"fresh"})
	}
}
