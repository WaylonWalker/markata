package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestEditorialHTMLContractAndCachedFallback(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "index.html")
	classic := `<div class="posts posts-list">old fallback</div>`
	editorial := `<div class="feed feed--editorial" data-feed-presentation="editorial"><div class="editorial-feed posts-list" data-editorial-version="4"></div></div>`
	priorEditorial := `<div class="feed feed--editorial" data-feed-presentation="editorial"><div class="editorial-feed posts-list" data-editorial-version="3"></div></div>`
	if editorialHTMLMarkupValid(priorEditorial) {
		t.Fatal("prior editorial markup must be rebuilt for the new layout")
	}
	if editorialHTMLPublished(file) {
		t.Fatal("missing output cannot count as editorial")
	}
	if err := os.WriteFile(file, []byte(classic), 0o600); err != nil {
		t.Fatal(err)
	}
	if editorialHTMLPublished(file) {
		t.Fatal("classic fallback must invalidate editorial cache")
	}
	if !editorialHTMLMarkupValid(editorial) {
		t.Fatal("valid editorial markup was rejected")
	}
	if err := os.WriteFile(file, []byte(editorial), 0o600); err != nil {
		t.Fatal(err)
	}
	if !editorialHTMLPublished(file) {
		t.Fatal("editorial output must be reusable")
	}
}

func TestEditorialPageRendersMarkupNotFallback(t *testing.T) {
	cfg := lifecycle.NewConfig()
	cfg.Extra = map[string]interface{}{"title": "Test site", "url": "https://example.com"}
	title := "An article"
	date := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	post := &models.Post{
		Slug: "blog/example", Href: "/blog/example/", Title: &title,
		Published: true, Template: "post", Date: &date,
		ArticleHTML: "<p>Example.</p>",
	}
	fc := &models.FeedConfig{
		Slug: "archive", Title: "Archive", Presentation: "editorial",
		Posts: []*models.Post{post},
	}
	page := &models.FeedPage{Number: 1, Posts: fc.Posts, TotalPages: 1, TotalItems: 1}
	html, err := NewPublishFeedsPlugin().generateFeedPageHTML(fc, page, cfg, nil, buildFeedRenderContext(fc))
	if err != nil {
		t.Fatalf("editorial page failed: %v", err)
	}
	if !editorialHTMLMarkupValid(html) || !strings.Contains(html, "data-editorial-relative-week") ||
		!strings.Contains(html, "data-week-start=\"2026-10-12\"") {
		t.Fatal("editorial page rendered without weekly markup")
	}
}

func TestEditorialMissingTemplateIsError(t *testing.T) {
	cfg := lifecycle.NewConfig()
	cfg.Extra = map[string]interface{}{"title": "Test site"}
	fc := &models.FeedConfig{
		Slug: "archive", Presentation: "editorial",
		Templates: models.FeedTemplates{HTML: "missing-editorial-template-2026.html"},
	}
	page := &models.FeedPage{Number: 1, TotalPages: 1}
	_, err := NewPublishFeedsPlugin().generateFeedPageHTML(fc, page, cfg, nil, buildFeedRenderContext(fc))
	if err == nil || !strings.Contains(err.Error(), "cannot use classic fallback") {
		t.Fatalf("expected actionable error, got %v", err)
	}
}

func TestCachedEditorialFeedRepublishesClassicHTML(t *testing.T) {
	plugin := NewPublishFeedsPlugin()
	outputDir := t.TempDir()
	config := lifecycle.NewConfig()
	config.OutputDir = outputDir
	config.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Test Site",
	}
	title := "A post"
	date := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	post := &models.Post{
		Path: "pages/post.md", Slug: "post", Href: "/post/",
		Title: &title, Published: true, Template: "post", Date: &date,
		ArticleHTML: "<p>hello</p>",
	}
	feed := &models.FeedConfig{
		Slug: "archive", Title: "Archive", Presentation: "editorial",
		ItemsPerPage: 10, Formats: models.FeedFormats{HTML: true},
		Posts: []*models.Post{post},
	}
	feed.Pages = []models.FeedPage{{
		Number: 1, Posts: feed.Posts, TotalPages: 1,
		TotalItems: 1, ItemsPerPage: 10, PageURLs: []string{"/archive/"},
	}}
	if err := plugin.publishFeed(feed, config, outputDir); err != nil {
		t.Fatalf("editorial publication: %v", err)
	}
	index := filepath.Join(outputDir, "archive", "index.html")
	if !editorialHTMLPublished(index) {
		t.Fatal("published archive lacks editorial markup")
	}
	cache := buildcache.New(t.TempDir())
	hash := plugin.computeFeedHash(feed)
	cache.SetFeedHash(feed.Slug, hash)
	skip, _ := plugin.shouldSkipFeed(feed, cache, outputDir)
	if !skip {
		t.Fatal("valid editorial archive should be reused")
	}
	if err := os.WriteFile(index, []byte("<div class=\"feed--editorial\" data-feed-presentation=\"editorial\"><div class=\"editorial-feed posts-list\" data-editorial-version=\"2\">stale editorial layout</div></div>"), 0o600); err != nil {
		t.Fatal(err)
	}
	skip, _ = plugin.shouldSkipFeed(feed, cache, outputDir)
	if skip {
		t.Fatal("cached old editorial layout must be rebuilt despite matching feed hash")
	}
	if err := os.WriteFile(index, []byte("<div class=\"posts-list\">classic fallback</div>"), 0o600); err != nil {
		t.Fatal(err)
	}
	skip, _ = plugin.shouldSkipFeed(feed, cache, outputDir)
	if skip {
		t.Fatal("cached classic fallback must be rebuilt despite matching feed hash")
	}
}
