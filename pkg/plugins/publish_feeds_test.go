package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// TestPublishFeedsPlugin_Name tests the plugin name.
func TestPublishFeedsPlugin_Name(t *testing.T) {
	p := NewPublishFeedsPlugin()
	if got := p.Name(); got != "publish_feeds" {
		t.Errorf("Name() = %q, want %q", got, "publish_feeds")
	}
}

func TestPublishCalendarJSONWritesDatedPostMetadata(t *testing.T) {
	outputDir := t.TempDir()
	title := "Calendar entry"
	description := "Preview text"
	date := time.Date(2026, time.October, 13, 23, 15, 0, 0, time.UTC)
	undatedTitle := "Undated"
	feed := &models.FeedConfig{Posts: []*models.Post{
		{
			Slug:        "calendar-entry",
			Href:        "/calendar-entry/",
			Title:       &title,
			Description: &description,
			Date:        &date,
			Extra:       map[string]interface{}{"image": "/images/calendar-entry.webp"},
		},
		{Slug: "undated", Href: "/undated/", Title: &undatedTitle},
		{Slug: "missing-href", Title: &title, Date: &date},
	}}

	if err := NewPublishFeedsPlugin().publishCalendarJSON(feed, outputDir); err != nil {
		t.Fatalf("publishCalendarJSON() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(outputDir, "calendar.json"))
	if err != nil {
		t.Fatalf("ReadFile(calendar.json) error = %v", err)
	}
	var got []calendarFeedItem
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal(calendar.json) error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("calendar items = %d, want only the dated post with an href", len(got))
	}
	want := calendarFeedItem{
		Date:        "2026-10-13",
		Href:        "/calendar-entry/",
		Title:       title,
		Description: description,
		Image:       "/images/calendar-entry.webp",
	}
	if got[0] != want {
		t.Fatalf("calendar item = %+v, want %+v", got[0], want)
	}
}

func TestPublishFeedsPlugin_FastFeedPublishingWaitsForErrors(t *testing.T) {
	outputDir := t.TempDir()
	blockedFeedPath := filepath.Join(outputDir, "blocked")
	if err := os.WriteFile(blockedFeedPath, []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	m.Config().Extra = map[string]interface{}{
		"fast_mode":   true,
		"feeds_async": true,
	}
	post := &models.Post{Path: "post.md", Slug: "post", Published: true, Content: "body"}
	m.SetFiles([]string{"post.md"})
	m.Cache().Set("feed_configs", []models.FeedConfig{{
		Slug:    "blocked",
		Formats: models.FeedFormats{RSS: true},
		Posts:   []*models.Post{post},
	}})

	if err := NewPublishFeedsPlugin().Write(m); err == nil {
		t.Fatal("Write() error = nil, want feed publication error")
	}
	entry := findContentDisposition(m.ContentDiagnostics(), post.Path)
	if entry == nil || !containsDiagnosticCode(entry.Diagnostics, diagnostics.ReasonContentWriteError) {
		t.Fatalf("feed write diagnostics = %+v, want content.write_error", entry)
	}
}

func TestPublishFeedsPlugin_FastFeedPublishingCompletesBeforeReturn(t *testing.T) {
	outputDir := t.TempDir()
	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	m.Config().Extra = map[string]interface{}{
		"fast_mode":   true,
		"feeds_async": true,
		"title":       "Test Site",
		"url":         "https://example.com",
	}
	post := &models.Post{Path: "post.md", Slug: "post", Published: true, Content: "body"}
	m.SetFiles([]string{"post.md"})
	cache := buildcache.New(t.TempDir())
	m.Cache().Set("build_cache", cache)
	m.Cache().Set("feed_configs", []models.FeedConfig{{
		Slug:    "recent",
		Formats: models.FeedFormats{RSS: true},
		Posts:   []*models.Post{post},
	}})

	if err := NewPublishFeedsPlugin().Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "recent", "rss.xml")); err != nil {
		t.Fatalf("feed output missing after Write() returned: %v", err)
	}
	if cache.GetFeedHash("recent") == "" {
		t.Fatal("feed hash missing after Write() returned")
	}
	entry := findContentDisposition(m.ContentDiagnostics(), post.Path)
	if entry == nil || !entry.Emitted || len(entry.Feeds) != 1 || !entry.Feeds[0].Included {
		t.Fatalf("feed output diagnostics = %+v, want emitted and included", entry)
	}
}

func TestRecordFeedOutput_OnlyMarksEnabledFormatOutputs(t *testing.T) {
	title := "Title only"
	posts := []*models.Post{
		{Path: "title-only.md", Slug: "title-only", Title: &title, Published: true},
		{Path: "article.md", Slug: "article", Published: true, Content: "article body"},
	}
	m := lifecycle.NewManager()
	m.SetFiles([]string{"title-only.md", "article.md"})

	recordFeedOutput(m, &models.FeedConfig{
		Slug:    "rss-only",
		Formats: models.FeedFormats{RSS: true},
		Posts:   posts,
	})

	snapshot := m.ContentDiagnostics()
	titleOnly := findContentDisposition(snapshot, "title-only.md")
	article := findContentDisposition(snapshot, "article.md")
	if titleOnly == nil || titleOnly.Emitted {
		t.Fatalf("title-only disposition = %+v, want not emitted", titleOnly)
	}
	if article == nil || !article.Emitted {
		t.Fatalf("article disposition = %+v, want emitted", article)
	}
	if len(titleOnly.Feeds) != 1 || titleOnly.Feeds[0].Included || !containsReason(titleOnly.Feeds[0].Reasons, diagnostics.ReasonContentNoOutput) {
		t.Fatalf("title-only feed disposition = %+v", titleOnly.Feeds)
	}
}

func TestRecordFeedOutput_UnpublishedShadowRemainsExcludedFromFeed(t *testing.T) {
	m := lifecycle.NewManager()
	m.SetFiles([]string{"unpublished.md"})
	post := &models.Post{
		Path:        "unpublished.md",
		Slug:        "unpublished",
		Published:   false,
		ArticleHTML: "<p>shadow</p>",
	}

	recordFeedOutput(m, &models.FeedConfig{
		Slug:    "public",
		Formats: models.FeedFormats{HTML: true},
		Posts:   []*models.Post{post},
	})

	entry := findContentDisposition(m.ContentDiagnostics(), post.Path)
	if entry == nil || !entry.Emitted || entry.Disposition != diagnostics.DispositionShadow {
		t.Fatalf("unpublished shadow disposition = %+v, want emitted shadow", entry)
	}
	if len(entry.Feeds) != 1 || entry.Feeds[0].Included || !containsReason(entry.Feeds[0].Reasons, diagnostics.ReasonContentPublishedFalse) {
		t.Fatalf("unpublished feed disposition = %+v, want excluded for published_false", entry.Feeds)
	}
}

func TestGenerateFeedPageHTML_UsesConfiguredHTMLTemplate(t *testing.T) {
	t.Parallel()

	p := NewPublishFeedsPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Example Site",
	}

	title := "Grid Shot"
	description := "Overlay copy"
	now := time.Date(2026, time.January, 2, 15, 4, 5, 0, time.UTC)

	fc := &models.FeedConfig{
		Slug:        "shots",
		Title:       "Shots",
		Description: "Photo and video posts",
		Templates: models.FeedTemplates{
			HTML: "feed-photo-grid.html",
		},
		Posts: []*models.Post{{
			Slug:        "shots/grid-shot",
			Href:        "/shots/grid-shot/",
			Title:       &title,
			Description: &description,
			Date:        &now,
			Published:   true,
			Extra: map[string]interface{}{
				"image": "https://example.com/grid-shot.webp",
			},
		}},
	}
	page := &models.FeedPage{Posts: fc.Posts, TotalPages: 1}

	html, err := p.generateFeedPageHTML(fc, page, config, nil, buildFeedRenderContext(fc))
	if err != nil {
		t.Fatalf("generateFeedPageHTML() error = %v", err)
	}

	if !strings.Contains(html, `feed--photo-grid`) {
		t.Fatalf("expected configured template output to contain photo grid feed class, got %q", html)
	}
	if !strings.Contains(html, `shot-card`) {
		t.Fatalf("expected configured template output to contain shot-card markup, got %q", html)
	}
}

func TestGenerateFeedPageHTML_PinsSourcePrefersLinkThenURL(t *testing.T) {
	t.Parallel()

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"url": "https://example.com", "title": "Example Site"}
	title := "Saved source"
	posts := []*models.Post{
		{Slug: "url-only", Href: "/url-only/", Title: &title, Published: true,
			Extra: map[string]interface{}{"url": "https://source.example/url-only"}},
		{Slug: "both", Href: "/both/", Title: &title, Published: true,
			Extra: map[string]interface{}{"link": "https://source.example/preferred", "url": "https://source.example/fallback"}},
	}
	feed := &models.FeedConfig{
		Slug: "pins", Title: "Pins", Templates: models.FeedTemplates{HTML: "pins.html"}, Posts: posts,
	}
	page := &models.FeedPage{Posts: posts, TotalPages: 1}

	html, err := NewPublishFeedsPlugin().generateFeedPageHTML(feed, page, config, nil, buildFeedRenderContext(feed))
	if err != nil {
		t.Fatalf("generateFeedPageHTML() error = %v", err)
	}
	for _, source := range []string{"https://source.example/url-only", "https://source.example/preferred"} {
		if !strings.Contains(html, `class="pin-source" href="`+source+`"`) {
			t.Errorf("Pins source link %q missing", source)
		}
	}
	if strings.Contains(html, `class="pin-source" href="https://source.example/fallback"`) {
		t.Error("Pins source used url when link was present")
	}
}

func TestGenerateSimpleFeedPageHTML_LoadsDecryptionAssets(t *testing.T) {
	t.Parallel()

	p := NewPublishFeedsPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Example Site",
	}

	fc := &models.FeedConfig{
		Slug:  "private",
		Title: "Private",
		Posts: []*models.Post{{
			Slug:  "private/entry",
			Href:  "/private/entry/",
			Extra: map[string]interface{}{"has_encrypted_content": true},
		}},
	}
	page := &models.FeedPage{Posts: fc.Posts, TotalPages: 1}

	html, err := p.generateSimpleFeedPageHTML(fc, page, config, nil, buildFeedRenderContext(fc))
	if err != nil {
		t.Fatalf("generateSimpleFeedPageHTML() error = %v", err)
	}
	if !containsRenderedAsset(html, "js/decryption.js") {
		t.Fatalf("expected simple feed output to load decryption.js, got %q", html)
	}
}

func TestCleanupPaginatedFeedDirs_RemovesStalePageDirectories(t *testing.T) {
	t.Parallel()

	plugin := NewPublishFeedsPlugin()
	feedDir := t.TempDir()

	stalePaths := []string{
		filepath.Join(feedDir, "page", "3"),
		filepath.Join(feedDir, "page", "4"),
		filepath.Join(feedDir, "simple", "page", "3"),
		filepath.Join(feedDir, "simple", "page", "4"),
	}
	for _, stalePath := range stalePaths {
		if err := os.MkdirAll(stalePath, 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", stalePath, err)
		}
	}

	pages := []models.FeedPage{{Number: 1}}

	if err := plugin.cleanupPaginatedFeedDirs(feedDir, "", pages); err != nil {
		t.Fatalf("cleanupPaginatedFeedDirs(html) error = %v", err)
	}
	if err := plugin.cleanupPaginatedFeedDirs(feedDir, "simple", pages); err != nil {
		t.Fatalf("cleanupPaginatedFeedDirs(simple) error = %v", err)
	}

	for _, removedPath := range []string{
		filepath.Join(feedDir, "page"),
		filepath.Join(feedDir, "simple", "page"),
	} {
		if _, err := os.Stat(removedPath); !os.IsNotExist(err) {
			t.Fatalf("expected %q to be removed, stat err = %v", removedPath, err)
		}
	}
}

// TestPublishFeedsPlugin_FormatRedirectsCreateDirectories tests that format redirects work correctly.
// - For md/txt: Content at /slug.ext, reversed redirect from /slug/index.ext/index.html -> /slug.ext
// - For json: Content at /slug/feed.json, forward redirect from /slug.json/index.html -> /slug/feed.json
// Fixes: https://github.com/WaylonWalker/markata/issues/248
func TestPublishFeedsPlugin_FormatRedirectsCreateDirectories(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishFeedsPlugin()

	// Create manager with feed configs in cache
	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = tempDir
	cfg.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Test Site",
	}

	// Create a feed config with all format types enabled
	feedConfigs := []models.FeedConfig{
		{
			Slug:        "archive",
			Title:       "Archive",
			Description: "All posts",
			Posts:       []*models.Post{},
			Pages:       []models.FeedPage{},
			Formats: models.FeedFormats{
				HTML:     true,
				RSS:      true,
				Atom:     true,
				JSON:     true,
				Markdown: true,
				Text:     true,
			},
		},
	}
	m.Cache().Set("feed_configs", feedConfigs)

	// Run write
	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Test that md/txt content files exist at canonical short URLs
	t.Run("markdown content at short URL", func(t *testing.T) {
		contentPath := filepath.Join(tempDir, "archive.md")
		info, err := os.Stat(contentPath)
		if err != nil {
			t.Fatalf("markdown content file %s not found: %v", contentPath, err)
		}
		if info.IsDir() {
			t.Errorf("expected %s to be a file, but it's a directory", contentPath)
		}
	})

	t.Run("text content at short URL", func(t *testing.T) {
		contentPath := filepath.Join(tempDir, "archive.txt")
		info, err := os.Stat(contentPath)
		if err != nil {
			t.Fatalf("text content file %s not found: %v", contentPath, err)
		}
		if info.IsDir() {
			t.Errorf("expected %s to be a file, but it's a directory", contentPath)
		}
	})

	// Test that md/txt reversed redirects exist
	t.Run("markdown reversed redirect", func(t *testing.T) {
		// Redirect from /archive/index.md -> /archive.md
		redirectPath := filepath.Join(tempDir, "archive", "index.md", "index.html")
		content, err := os.ReadFile(redirectPath)
		if err != nil {
			t.Fatalf("markdown redirect %s not found: %v", redirectPath, err)
		}
		contentStr := string(content)
		if !strings.Contains(contentStr, "<!DOCTYPE html>") {
			t.Error("redirect file should contain DOCTYPE")
		}
		if !strings.Contains(contentStr, "meta http-equiv=\"refresh\"") {
			t.Error("redirect file should contain meta refresh")
		}
		if !strings.Contains(contentStr, "/archive.md") {
			t.Errorf("redirect should point to /archive.md, got: %s", contentStr)
		}
	})

	t.Run("text reversed redirect", func(t *testing.T) {
		// Redirect from /archive/index.txt -> /archive.txt
		redirectPath := filepath.Join(tempDir, "archive", "index.txt", "index.html")
		content, err := os.ReadFile(redirectPath)
		if err != nil {
			t.Fatalf("text redirect %s not found: %v", redirectPath, err)
		}
		contentStr := string(content)
		if !strings.Contains(contentStr, "<!DOCTYPE html>") {
			t.Error("redirect file should contain DOCTYPE")
		}
		if !strings.Contains(contentStr, "meta http-equiv=\"refresh\"") {
			t.Error("redirect file should contain meta refresh")
		}
		if !strings.Contains(contentStr, "/archive.txt") {
			t.Errorf("redirect should point to /archive.txt, got: %s", contentStr)
		}
	})

	// Test that JSON still uses forward redirect (unchanged behavior)
	t.Run("json forward redirect", func(t *testing.T) {
		// JSON content at /archive/feed.json
		jsonContentPath := filepath.Join(tempDir, "archive", "feed.json")
		if _, err := os.Stat(jsonContentPath); err != nil {
			t.Errorf("JSON content file %s not found: %v", jsonContentPath, err)
		}

		// Redirect from /archive.json/index.html -> /archive/feed.json
		redirectDir := filepath.Join(tempDir, "archive.json")
		info, err := os.Stat(redirectDir)
		if err != nil {
			t.Fatalf("json redirect path %s not found: %v", redirectDir, err)
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory, but it's a file", redirectDir)
		}

		indexPath := filepath.Join(redirectDir, "index.html")
		content, err := os.ReadFile(indexPath)
		if err != nil {
			t.Fatalf("index.html not found in %s: %v", redirectDir, err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "<!DOCTYPE html>") {
			t.Error("redirect file should contain DOCTYPE")
		}
		if !strings.Contains(contentStr, "meta http-equiv=\"refresh\"") {
			t.Error("redirect file should contain meta refresh")
		}
		if !strings.Contains(contentStr, "/archive/feed.json") {
			t.Errorf("redirect file should point to /archive/feed.json, got: %s", contentStr)
		}
	})
}

// TestPublishFeedsPlugin_RootFeedNoRedirects tests that root feeds (empty slug) don't create redirects.
// Root feeds write to the output root, so there's no slug to redirect from.
func TestPublishFeedsPlugin_RootFeedNoRedirects(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishFeedsPlugin()

	// Create manager with root feed config
	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = tempDir
	cfg.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Test Site",
	}

	// Create a root feed config (empty slug) with formats enabled
	feedConfigs := []models.FeedConfig{
		{
			Slug:        "", // Root feed
			Title:       "Home",
			Description: "Home page",
			Posts:       []*models.Post{},
			Pages:       []models.FeedPage{},
			Formats: models.FeedFormats{
				Markdown: true,
				Text:     true,
				JSON:     true,
			},
		},
	}
	m.Cache().Set("feed_configs", feedConfigs)

	// Run write
	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify that no redirect directories were created at root level
	// (They would be named ".md", ".txt", ".json" which doesn't make sense)
	redirectPaths := []string{
		filepath.Join(tempDir, ".md"),
		filepath.Join(tempDir, ".txt"),
		filepath.Join(tempDir, ".json"),
	}

	for _, path := range redirectPaths {
		if _, err := os.Stat(path); err == nil {
			t.Errorf("root feed should not create redirect at %s", path)
		}
	}

	// Verify the actual format files were created at root
	expectedFiles := []string{
		filepath.Join(tempDir, "index.md"),
		filepath.Join(tempDir, "index.txt"),
		filepath.Join(tempDir, "feed.json"),
	}

	for _, path := range expectedFiles {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("expected file %s not found: %v", path, err)
		}
	}
}

func TestPublishFeedsPlugin_ExcludesPostsWithoutRenderableOutput(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishFeedsPlugin()

	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = tempDir
	cfg.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Test Site",
	}

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	posts := []*models.Post{
		{Path: "rendered.md", Slug: "rendered", Href: "/rendered/", Title: strPtr("Rendered"), Content: "rendered", ArticleHTML: "<p>Rendered</p>", Date: &date},
		{Path: "empty.md", Slug: "empty", Href: "/empty/", Title: strPtr("Empty"), Date: &date},
		{Path: "skipped.md", Slug: "skipped", Href: "/skipped/", Title: strPtr("Skipped"), Content: "skipped", ArticleHTML: "<p>Skipped</p>", Date: &date, Skip: true},
	}
	feedConfig := models.FeedConfig{
		Slug:        "blog",
		Title:       "Blog",
		Description: "Posts",
		Posts:       posts,
		Formats: models.FeedFormats{
			RSS:      true,
			Atom:     true,
			JSON:     true,
			Markdown: true,
			Text:     true,
		},
	}
	feedConfig.Paginate("/blog")
	m.Cache().Set("feed_configs", []models.FeedConfig{feedConfig})

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	checks := []string{
		filepath.Join(tempDir, "blog", "rss.xml"),
		filepath.Join(tempDir, "blog", "atom.xml"),
		filepath.Join(tempDir, "blog", "feed.json"),
		filepath.Join(tempDir, "blog.md"),
		filepath.Join(tempDir, "blog.txt"),
	}
	for _, path := range checks {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("ReadFile(%q) error = %v", path, err)
		}
		text := string(content)
		if !strings.Contains(text, "rendered") {
			t.Fatalf("%s should contain rendered post", path)
		}
		if strings.Contains(text, "empty") {
			t.Fatalf("%s should not contain empty post", path)
		}
		if strings.Contains(text, "skipped") {
			t.Fatalf("%s should not contain skipped post", path)
		}
	}
}

// TestPublishFeedsPlugin_writeFeedFormatRedirect tests the redirect file generation directly.
func TestPublishFeedsPlugin_writeFeedFormatRedirect(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishFeedsPlugin()

	tests := []struct {
		name       string
		slug       string
		ext        string
		targetFile string
		wantURL    string
	}{
		{
			name:       "json redirect",
			slug:       "posts",
			ext:        "json",
			targetFile: "feed.json",
			wantURL:    "/posts/feed.json",
		},
		{
			name:       "nested slug redirect",
			slug:       "tags/python",
			ext:        "json",
			targetFile: "feed.json",
			wantURL:    "/tags/python/feed.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh subdirectory for each test
			testDir := filepath.Join(tempDir, tt.name)
			if err := os.MkdirAll(testDir, 0o755); err != nil {
				t.Fatalf("failed to create test dir: %v", err)
			}

			// Write redirect
			if err := plugin.writeFeedFormatRedirect(tt.slug, tt.ext, tt.targetFile, testDir); err != nil {
				t.Fatalf("writeFeedFormatRedirect() error = %v", err)
			}

			// Verify redirect directory was created
			redirectDir := filepath.Join(testDir, tt.slug+"."+tt.ext)
			info, err := os.Stat(redirectDir)
			if err != nil {
				t.Fatalf("redirect directory not found: %v", err)
			}
			if !info.IsDir() {
				t.Error("expected redirect path to be a directory")
			}

			// Verify index.html exists
			indexPath := filepath.Join(redirectDir, "index.html")
			content, err := os.ReadFile(indexPath)
			if err != nil {
				t.Fatalf("failed to read index.html: %v", err)
			}

			contentStr := string(content)

			// Verify HTML structure
			if !strings.Contains(contentStr, "<!DOCTYPE html>") {
				t.Error("missing DOCTYPE")
			}
			if !strings.Contains(contentStr, "meta http-equiv=\"refresh\"") {
				t.Error("missing meta refresh")
			}
			if !strings.Contains(contentStr, tt.wantURL) {
				t.Errorf("redirect should point to %s, got: %s", tt.wantURL, contentStr)
			}
			if !strings.Contains(contentStr, `rel="canonical"`) {
				t.Error("missing canonical link")
			}
		})
	}
}

// TestComputeFeedHash_SortChangeProducesDifferentHash verifies that changing
// the sort or reverse config fields produces a different feed hash, preventing
// stale cache hits.
func TestComputeFeedHash_SortChangeProducesDifferentHash(t *testing.T) {
	p := NewPublishFeedsPlugin()

	baseFeed := &models.FeedConfig{
		Slug:         "archive",
		Title:        "Archive",
		Description:  "All posts",
		Sort:         "date",
		Reverse:      false,
		ItemsPerPage: 10,
		Filter:       `published == true`,
		Posts: []*models.Post{
			{Slug: "post-a"},
			{Slug: "post-b"},
		},
	}

	tests := []struct {
		name   string
		modify func(fc *models.FeedConfig)
	}{
		{
			name: "sort field change",
			modify: func(fc *models.FeedConfig) {
				fc.Sort = "title"
			},
		},
		{
			name: "reverse change",
			modify: func(fc *models.FeedConfig) {
				fc.Reverse = true
			},
		},
		{
			name: "filter change",
			modify: func(fc *models.FeedConfig) {
				fc.Filter = `tags contains "go"`
			},
		},
		{
			name: "description change",
			modify: func(fc *models.FeedConfig) {
				fc.Description = "Updated description"
			},
		},
		{
			name: "items per page change",
			modify: func(fc *models.FeedConfig) {
				fc.ItemsPerPage = 20
			},
		},
		{
			name: "format change",
			modify: func(fc *models.FeedConfig) {
				fc.Formats.RSS = true
			},
		},
		{
			name: "template change",
			modify: func(fc *models.FeedConfig) {
				fc.Templates.HTML = "custom.html"
			},
		},
		{
			name: "post order change",
			modify: func(fc *models.FeedConfig) {
				fc.Posts = []*models.Post{
					{Slug: "post-b"},
					{Slug: "post-a"},
				}
			},
		},
		{
			name: "include private change",
			modify: func(fc *models.FeedConfig) {
				fc.IncludePrivate = true
			},
		},
	}

	baseHash := p.computeFeedHash(baseFeed)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy and modify it
			modified := *baseFeed
			modified.Formats = baseFeed.Formats
			modified.Templates = baseFeed.Templates
			// Deep copy posts slice
			modified.Posts = make([]*models.Post, len(baseFeed.Posts))
			for i, post := range baseFeed.Posts {
				p := *post
				modified.Posts[i] = &p
			}
			tt.modify(&modified)

			modifiedHash := p.computeFeedHash(&modified)
			if modifiedHash == baseHash {
				t.Errorf("hash should change when %s, but both are %s", tt.name, baseHash)
			}
		})
	}
}

// TestComputeFeedHash_IdenticalConfigsProduceSameHash verifies determinism.
func TestComputeFeedHash_IdenticalConfigsProduceSameHash(t *testing.T) {
	p := NewPublishFeedsPlugin()

	fc := &models.FeedConfig{
		Slug:         "archive",
		Title:        "Archive",
		Description:  "All posts",
		Sort:         "date",
		Reverse:      true,
		ItemsPerPage: 10,
		Filter:       `published == true`,
		Posts: []*models.Post{
			{Slug: "post-a"},
			{Slug: "post-b"},
		},
		Formats: models.FeedFormats{
			HTML: true,
			RSS:  true,
		},
	}

	hash1 := p.computeFeedHash(fc)
	hash2 := p.computeFeedHash(fc)

	if hash1 != hash2 {
		t.Errorf("identical configs should produce same hash: %s != %s", hash1, hash2)
	}
}

func TestComputeFeedHash_PostContentChangeProducesDifferentHash(t *testing.T) {
	p := NewPublishFeedsPlugin()
	title := "Post"
	description := "Desc"
	feed := &models.FeedConfig{
		Slug:         "archive",
		Title:        "Archive",
		ItemsPerPage: 10,
		Formats:      models.FeedFormats{HTML: true},
		Posts: []*models.Post{{
			Path:        "pages/post.md",
			Slug:        "post",
			Href:        "/post/",
			Title:       &title,
			Description: &description,
			Published:   true,
			Content:     "hello",
			ArticleHTML: "<p>hello</p>",
		}},
		Pages: []models.FeedPage{{Number: 1, Posts: []*models.Post{{Slug: "post"}}}},
	}

	hash1 := p.computeFeedHash(feed)
	feed.Posts[0].Content = "goodbye"
	hash2 := p.computeFeedHash(feed)

	if hash1 == hash2 {
		t.Fatalf("hash should change when post feed content changes: %s", hash1)
	}
}

func TestComputeFeedHash_PrivateEncryptedWrapperChangeProducesDifferentHash(t *testing.T) {
	p := NewPublishFeedsPlugin()
	feed := &models.FeedConfig{
		Slug:           "private",
		IncludePrivate: true,
		ItemsPerPage:   10,
		Formats:        models.FeedFormats{HTML: true, RSS: true, JSON: true},
		Posts: []*models.Post{{
			Path:        "private.md",
			Slug:        "private",
			Href:        "/private/",
			Private:     true,
			Published:   true,
			ArticleHTML: `<div class="encrypted-content" data-encrypted="ciphertext-one" data-key-name="private">locked</div>`,
		}},
	}

	first := p.computeFeedHash(feed)
	feed.Posts[0].ArticleHTML = `<div class="encrypted-content" data-encrypted="ciphertext-two" data-key-name="private">locked</div>`
	second := p.computeFeedHash(feed)
	if first == second {
		t.Fatalf("private feed hash should change when encrypted content changes: %s", first)
	}
}

func TestGetPostFeedItemHash_PrivateIgnoresStaleSemanticCache(t *testing.T) {
	p := NewPublishFeedsPlugin()
	post := &models.Post{
		Path:        "private.md",
		Slug:        "private",
		Href:        "/private/",
		Private:     true,
		ArticleHTML: `<div class="encrypted-content" data-encrypted="current" data-key-name="private">locked</div>`,
	}
	cache := buildcache.New(t.TempDir())
	cache.UpdatePostSemanticHashes(post.Path, "stale-load-hash", "", "")

	got := p.getPostFeedItemHash(post, cache)
	if got == "stale-load-hash" {
		t.Fatal("private feed hash should not use the stale load-stage semantic hash")
	}
}

func TestFeedConfigWithRenderablePosts_KeepsTitleOnlyPosts(t *testing.T) {
	title := "Gratitude"
	feed := &models.FeedConfig{
		Slug:           "tags/gratitude",
		IncludePrivate: true,
		Posts: []*models.Post{
			{Slug: "entry-1", Title: &title, Published: true, Private: true},
			{Slug: "entry-2", Title: &title, Published: true, Private: true, Content: "has body"},
		},
	}

	renderable := feedConfigWithRenderablePosts(feed)
	if len(renderable.Posts) != 2 {
		t.Fatalf("expected title-only posts to remain in feed pages, got %d posts", len(renderable.Posts))
	}
	if len(renderable.Pages) != 1 || len(renderable.Pages[0].Posts) != 2 {
		t.Fatalf("expected pagination to include both posts, got %#v", renderable.Pages)
	}
}

func TestEnsureDefaultCollectionPage_EmptyRootGetsHomepage(t *testing.T) {
	feed := &models.FeedConfig{
		ItemsPerPage:   10,
		PaginationType: models.PaginationManual,
	}

	ensureDefaultCollectionPage(feed)

	if len(feed.Pages) != 1 {
		t.Fatalf("pages = %#v, want one homepage", feed.Pages)
	}
	page := feed.Pages[0]
	if page.Number != 1 || page.TotalPages != 1 || len(page.Posts) != 0 {
		t.Fatalf("homepage page = %#v, want an empty first page", page)
	}
}

func TestEnsureDefaultCollectionPage_EmptyArchiveGetsPage(t *testing.T) {
	feed := &models.FeedConfig{Slug: "archive"}

	ensureDefaultCollectionPage(feed)

	if len(feed.Pages) != 1 || feed.Pages[0].Number != 1 {
		t.Fatalf("pages = %#v, want one archive page", feed.Pages)
	}
}

func TestEnsureDefaultCollectionPage_EmptyNonDefaultFeedStaysWithoutPages(t *testing.T) {
	feed := &models.FeedConfig{Slug: "tags/go"}

	ensureDefaultCollectionPage(feed)

	if len(feed.Pages) != 0 {
		t.Fatalf("pages = %#v, want no pages", feed.Pages)
	}
}

func TestComputeFeedHash_RenderDecorationsDoNotChangeHash(t *testing.T) {
	p := NewPublishFeedsPlugin()
	title := "Post"
	feed := &models.FeedConfig{
		Slug:         "archive",
		Title:        "Archive",
		ItemsPerPage: 10,
		Formats:      models.FeedFormats{HTML: true},
		Posts: []*models.Post{{
			Path:        "pages/post.md",
			Slug:        "post",
			Href:        "/post/",
			Title:       &title,
			Published:   true,
			Content:     "hello",
			ArticleHTML: "<p>hello</p>",
		}},
	}

	hash1 := p.computeFeedHash(feed)
	feed.Posts[0].ArticleHTML = `<p class="has-avatar">hello</p>`
	hash2 := p.computeFeedHash(feed)

	if hash1 != hash2 {
		t.Fatalf("hash should ignore downstream render decorations: %s != %s", hash1, hash2)
	}
}

func TestPublishFeedsPlugin_DoesNotSkipFeedContainingAffectedDependent(t *testing.T) {
	plugin := NewPublishFeedsPlugin()
	outputDir := t.TempDir()
	feed := &models.FeedConfig{
		Slug:  "archive",
		Posts: []*models.Post{{Path: "content/source.md", Slug: "source", Content: "same source"}},
	}
	cache := buildcache.New(t.TempDir())
	hash := plugin.computeFeedHashWithConfigAndCache(feed, nil, cache)
	cache.SetFeedHash(feed.Slug, hash)
	if skip, _ := plugin.shouldSkipFeedWithConfigAndChanges(feed, cache, outputDir, nil, nil); !skip {
		t.Fatal("unchanged feed was not eligible for skipping")
	}
	cache.MarkSlugChanged("source")
	if skip, _ := plugin.shouldSkipFeedWithConfigAndChanges(feed, cache, outputDir, nil, nil); !skip {
		t.Fatal("changed slug without an affected path unnecessarily rebuilt the feed")
	}

	skip, gotHash := plugin.shouldSkipFeedWithConfigAndChanges(
		feed,
		cache,
		outputDir,
		nil,
		map[string]bool{"content/source.md": true},
	)
	if skip {
		t.Fatal("feed containing a changed dependent was skipped")
	}
	if gotHash != hash {
		t.Fatalf("feed hash = %q, want unchanged semantic hash %q", gotHash, hash)
	}
}

func TestPublishFeedsPlugin_ShouldSkipFeedWhenOutputsExist(t *testing.T) {
	plugin := NewPublishFeedsPlugin()
	outputDir := t.TempDir()
	title := "Post"

	config := lifecycle.NewConfig()
	config.OutputDir = outputDir
	config.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Test Site",
	}

	feed := &models.FeedConfig{
		Slug:         "archive",
		Title:        "Archive",
		Description:  "All posts",
		ItemsPerPage: 10,
		Formats: models.FeedFormats{
			HTML:     true,
			JSON:     true,
			Markdown: true,
		},
		Posts: []*models.Post{{
			Path:        "pages/post.md",
			Slug:        "post",
			Href:        "/post/",
			Title:       &title,
			Published:   true,
			ArticleHTML: "<p>hello</p>",
		}},
	}
	feed.Pages = []models.FeedPage{{
		Number:       1,
		Posts:        feed.Posts,
		TotalPages:   1,
		TotalItems:   1,
		ItemsPerPage: 10,
		PageURLs:     []string{"/archive/"},
	}}

	if err := plugin.publishFeed(feed, config, outputDir); err != nil {
		t.Fatalf("publishFeed() error = %v", err)
	}

	cache := buildcache.New(t.TempDir())
	hash := plugin.computeFeedHash(feed)
	cache.SetFeedHash(feed.Slug, hash)

	skip, gotHash := plugin.shouldSkipFeed(feed, cache, outputDir)
	if !skip {
		t.Fatalf("shouldSkipFeed() = false, want true")
	}
	if gotHash != hash {
		t.Fatalf("shouldSkipFeed() hash = %q, want %q", gotHash, hash)
	}

	if err := os.Remove(filepath.Join(outputDir, "archive", "feed.json")); err != nil {
		t.Fatalf("Remove(feed.json) error = %v", err)
	}
	if skip, _ := plugin.shouldSkipFeed(feed, cache, outputDir); skip {
		t.Fatalf("shouldSkipFeed() = true after deleting output, want false")
	}
}

func TestPublishFeedsPlugin_GeneratesArchiveVariants(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishFeedsPlugin()
	m := lifecycle.NewManager()
	cfg := m.Config()
	cfg.OutputDir = tempDir
	defaults := models.NewFeedDefaults()
	defaults.Syndication.MaxItems = 1
	cfg.Extra = map[string]interface{}{
		"url":           "https://example.com",
		"title":         "Test Site",
		"feed_defaults": defaults,
	}

	firstTitle := "First"
	secondTitle := "Second"
	firstDate := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
	secondDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	feedConfigs := []models.FeedConfig{{
		Slug:        "blog",
		Title:       "Blog",
		Description: "All posts",
		Formats:     models.FeedFormats{RSS: true, Atom: true, JSON: true},
		Posts: []*models.Post{
			{Slug: "first", Href: "/first/", Title: &firstTitle, Published: true, Date: &firstDate, ArticleHTML: "<p>first</p>"},
			{Slug: "second", Href: "/second/", Title: &secondTitle, Published: true, Date: &secondDate, ArticleHTML: "<p>second</p>"},
		},
	}}
	m.Cache().Set("feed_configs", feedConfigs)

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	primaryRSS, err := os.ReadFile(filepath.Join(tempDir, "blog", "rss.xml"))
	if err != nil {
		t.Fatalf("ReadFile(primary rss) error = %v", err)
	}
	archiveRSS, err := os.ReadFile(filepath.Join(tempDir, "blog", "archive", "rss.xml"))
	if err != nil {
		t.Fatalf("ReadFile(archive rss) error = %v", err)
	}

	if strings.Contains(string(primaryRSS), "/second/") {
		t.Fatalf("primary rss should be truncated to max_items=1")
	}
	if !strings.Contains(string(archiveRSS), "/second/") {
		t.Fatalf("archive rss should include older entries")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "blog", "archive", "atom.xml")); err != nil {
		t.Fatalf("archive atom not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tempDir, "blog", "archive", "feed.json")); err != nil {
		t.Fatalf("archive json not written: %v", err)
	}
}

func TestPublishFeedsPlugin_ShouldSkipFeed_UsesRenderableHTMLPagination(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishFeedsPlugin()
	config := lifecycle.NewManager().Config()
	config.OutputDir = tempDir
	config.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Test Site",
	}

	titleA := "Published"
	titleB := "Draft"
	feed := &models.FeedConfig{
		Slug:         "published",
		Title:        "Published",
		Description:  "Published posts",
		ItemsPerPage: 1,
		Formats: models.FeedFormats{
			HTML: true,
		},
		Posts: []*models.Post{
			{
				Path:        "pages/published.md",
				Slug:        "published",
				Href:        "/published/",
				Title:       &titleA,
				Published:   true,
				ArticleHTML: "<p>published</p>",
			},
			{
				Path:        "pages/draft.md",
				Slug:        "draft",
				Href:        "/draft/",
				Title:       &titleB,
				Draft:       true,
				Published:   false,
				ArticleHTML: "<p>draft</p>",
			},
		},
	}
	feed.Pages = []models.FeedPage{
		{Number: 1, Posts: []*models.Post{feed.Posts[0]}, TotalPages: 2, TotalItems: 2, ItemsPerPage: 1, PageURLs: []string{"/published/", "/published/page/2/"}, HasNext: true, NextURL: "/published/page/2/"},
		{Number: 2, Posts: []*models.Post{feed.Posts[1]}, TotalPages: 2, TotalItems: 2, ItemsPerPage: 1, PageURLs: []string{"/published/", "/published/page/2/"}, HasPrev: true, PrevURL: "/published/"},
	}

	if err := plugin.publishFeed(feed, config, tempDir); err != nil {
		t.Fatalf("publishFeed() error = %v", err)
	}

	cache := buildcache.New(t.TempDir())
	hash := plugin.computeFeedHash(feed)
	cache.SetFeedHash(feed.Slug, hash)

	skip, gotHash := plugin.shouldSkipFeed(feed, cache, tempDir)
	if !skip {
		t.Fatalf("shouldSkipFeed() = false, want true")
	}
	if gotHash != hash {
		t.Fatalf("shouldSkipFeed() hash = %q, want %q", gotHash, hash)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "published", "page", "2", "index.html")); !os.IsNotExist(err) {
		t.Fatalf("expected renderable pagination to omit page 2, stat err = %v", err)
	}
}
