package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/models"
)

func TestEditorialHTMLContractAndCachedFallback(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "index.html")
	classic := `<div class="posts posts-list">old fallback</div>`
	editorial := `<div class="feed feed--editorial" data-feed-presentation="editorial"><div class="editorial-feed posts-list"></div></div>`
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
	if !editorialHTMLMarkupValid(html) || !strings.Contains(html, "The week of") {
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
