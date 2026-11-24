package templates

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestDedicatedCalendarCarriesHoverPreviewMetadata(t *testing.T) {
	engine, err := NewEngineWithTheme("", "default")
	if err != nil {
		t.Fatalf("NewEngineWithTheme() error = %v", err)
	}

	title := "Archive post"
	description := "This summary should remain available to the dedicated calendar preview."
	date := time.Date(2026, time.September, 28, 12, 0, 0, 0, time.UTC)
	feed := &models.FeedConfig{
		Title: "Archive",
		Posts: []*models.Post{{
			Slug:        "archive",
			Href:        "/archive/",
			Title:       &title,
			Description: &description,
			Date:        &date,
			Published:   true,
			Extra:       map[string]interface{}{"image": "/images/archive.webp"},
		}},
	}
	page := &models.FeedPage{Number: 1, Posts: feed.Posts, TotalPages: 1, TotalItems: 1}
	ctx := NewFeedContext(feed, page, &models.Config{Title: "Test Site", URL: "https://example.com"})
	ctx.Set("feed_stats_total_posts", 1)

	html, err := engine.Render("calendar-feed.html", ctx)
	if err != nil {
		t.Fatalf("Render(calendar-feed.html) error = %v", err)
	}
	for _, want := range []string{
		`data-calendar-list data-calendar-primary`,
		`aria-label="Posts by date"`,
		`data-description="This summary should remain available to the dedicated calendar preview."`,
		`data-image=`,
		`/images/archive.webp`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("dedicated calendar missing preview parity hook %q", want)
		}
	}
}
