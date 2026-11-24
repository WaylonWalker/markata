package templates

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestCalendarFeedTemplateRendersCompleteFeed(t *testing.T) {
	engine, err := NewEngineWithTheme("", "default")
	if err != nil {
		t.Fatalf("NewEngineWithTheme() error = %v", err)
	}

	firstTitle := "Leap day, first"
	secondTitle := "Leap day, second"
	olderTitle := "Older post"
	undatedTitle := "Undated note"
	leapDay := time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC)
	olderDay := time.Date(2023, time.December, 31, 12, 0, 0, 0, time.UTC)

	feed := &models.FeedConfig{
		Slug:  "archive",
		Title: "Archive",
		Posts: []*models.Post{
			{Slug: "leap-one", Href: "/leap-one/", Title: &firstTitle, Date: &leapDay, Published: true},
			{Slug: "leap-two", Href: "/leap-two/", Title: &secondTitle, Date: &leapDay, Published: true},
			{Slug: "older", Href: "/older/", Title: &olderTitle, Date: &olderDay, Published: true},
			{Slug: "undated", Href: "/undated/", Title: &undatedTitle, Published: true},
		},
	}

	// Deliberately expose only one post on the current pagination page. The
	// calendar template must use feed.posts so the visual archive is complete.
	page := &models.FeedPage{
		Number:     1,
		Posts:      feed.Posts[:1],
		TotalPages: 4,
		TotalItems: len(feed.Posts),
	}
	config := &models.Config{Title: "Test Site", URL: "https://example.com"}
	ctx := NewFeedContext(feed, page, config)
	ctx.Set("feed_stats_total_posts", len(feed.Posts))
	ctx.Set("feed_stats_latest_post", "February 29, 2024")

	html, err := engine.Render("calendar-feed.html", ctx)
	if err != nil {
		t.Fatalf("Render(calendar-feed.html) error = %v", err)
	}

	if got := strings.Count(html, `data-date="2024-02-29"`); got != 2 {
		t.Fatalf("leap-day entries = %d, want 2", got)
	}
	if !strings.Contains(html, `data-date="2023-12-31"`) || !strings.Contains(html, olderTitle) {
		t.Fatalf("calendar source omitted a post outside the current pagination page")
	}
	if !strings.Contains(html, undatedTitle) || !strings.Contains(html, "Undated") {
		t.Fatalf("undated post was omitted from the list fallback")
	}
	if !strings.Contains(html, "css/calendar-feed") {
		t.Fatalf("calendar stylesheet is not linked")
	}
	if !strings.Contains(html, "js/calendar-feed") {
		t.Fatalf("calendar script is not linked")
	}
	if !strings.Contains(html, `data-calendar-mode="list"`) || !strings.Contains(html, `data-calendar-mode="calendar"`) {
		t.Fatalf("calendar/list view controls are missing")
	}
}
