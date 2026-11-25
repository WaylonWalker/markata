package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

func TestPublishFeeds_ConcurrentHTMLRenderingSharesPosts(t *testing.T) {
	outputDir := t.TempDir()
	m := lifecycle.NewManager()
	m.Config().OutputDir = outputDir
	m.Config().Extra = map[string]interface{}{
		"title": "Concurrent feeds",
		"url":   "https://example.com",
	}
	m.SetConcurrency(16)

	date := time.Date(2026, time.November, 21, 12, 0, 0, 0, time.UTC)
	posts := make([]*models.Post, 64)
	for i := range posts {
		title := fmt.Sprintf("Post %d", i)
		posts[i] = &models.Post{
			Path:      fmt.Sprintf("posts/%03d.md", i),
			Slug:      fmt.Sprintf("post-%03d", i),
			Href:      fmt.Sprintf("/post-%03d/", i),
			Title:     &title,
			Date:      &date,
			Template:  "photo",
			Published: true,
			Extra: map[string]interface{}{
				"image": "/example.png",
				"link_metadata": map[string]interface{}{
					"site_name": "Example",
					"image":     "/example.png",
				},
			},
		}
	}
	m.SetPosts(posts)

	feeds := make([]models.FeedConfig, 32)
	for i := range feeds {
		feeds[i] = models.FeedConfig{
			Slug:         fmt.Sprintf("feed-%02d", i),
			Title:        fmt.Sprintf("Feed %d", i),
			Presentation: presentationEditorial,
			Formats:      models.FeedFormats{HTML: true},
			ItemsPerPage: len(posts),
			Posts:        posts,
		}
	}
	m.Cache().Set("feed_configs", feeds)

	if err := NewPublishFeedsPlugin().Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	page, err := os.ReadFile(filepath.Join(outputDir, "feed-00", "index.html"))
	if err != nil {
		t.Fatalf("reading rendered feed page: %v", err)
	}
	if len(page) == 0 {
		t.Fatal("rendered feed page is empty")
	}
}

func TestBuildEditorialFeedSections_DoesNotMutateCachedPostMap(t *testing.T) {
	date := time.Date(2026, time.November, 21, 12, 0, 0, 0, time.UTC)
	post := &models.Post{
		Slug:     "photo",
		Template: "photo",
		Date:     &date,
		Extra: map[string]interface{}{
			"image": "/example.png",
		},
	}
	cachedPostMap := templates.GetPostMap(post)
	if _, exists := cachedPostMap["editorial_photo_date"]; exists {
		t.Fatal("cached post map unexpectedly has an editorial-only field before grouping")
	}

	sections := buildEditorialFeedSections([]*models.Post{post})
	if len(sections) != 1 || len(sections[0].Items) != 1 {
		t.Fatalf("editorial sections = %+v, want one photo item", sections)
	}
	if _, exists := cachedPostMap["editorial_photo_date"]; exists {
		t.Fatal("editorial grouping mutated the shared cached post map")
	}
	if got := sections[0].Items[0].Posts[0]["editorial_photo_date"]; got == nil {
		t.Fatal("editorial photo date missing from the feed-local post map")
	}
}
