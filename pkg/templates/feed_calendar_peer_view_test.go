package templates

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestDefaultFeedDefersCalendarHistoryToJSON(t *testing.T) {
	engine, err := NewEngineWithTheme("", "default")
	if err != nil {
		t.Fatalf("NewEngineWithTheme() error = %v", err)
	}

	newTitle := "Newest post"
	oldTitle := "Older post outside page one"
	oldDescription := "A compact summary for the calendar hover preview."
	newDate := time.Date(2026, time.September, 28, 12, 0, 0, 0, time.UTC)
	oldDate := time.Date(2024, time.February, 29, 12, 0, 0, 0, time.UTC)
	feed := &models.FeedConfig{
		Slug:  "blog",
		Title: "Blog",
		Formats: models.FeedFormats{
			HTML:       true,
			SimpleHTML: true,
		},
		Posts: []*models.Post{
			{Slug: "new", Href: "/new/", Title: &newTitle, Date: &newDate, Published: true},
			{
				Slug:        "old",
				Href:        "/old/",
				Title:       &oldTitle,
				Description: &oldDescription,
				Date:        &oldDate,
				Published:   true,
				Extra:       map[string]interface{}{"image": "/images/old.webp"},
			},
		},
	}
	page := &models.FeedPage{
		Number:         1,
		Posts:          feed.Posts[:1],
		TotalPages:     2,
		TotalItems:     len(feed.Posts),
		ItemsPerPage:   1,
		PaginationType: models.PaginationManual,
	}
	ctx := NewFeedContext(feed, page, &models.Config{Title: "Test Site", URL: "https://example.com"})
	ctx.Set("feed_stats_total_posts", len(feed.Posts))
	ctx.Set("feed_stats_latest_post", "September 28, 2026")

	html, err := engine.Render("feed.html", ctx)
	if err != nil {
		t.Fatalf("Render(feed.html) error = %v", err)
	}

	for _, want := range []string{
		`data-calendar-feed`,
		`data-calendar-url-state`,
		`data-calendar-src="/blog/calendar.json"`,
		`data-calendar-primary`,
		`data-calendar-mode="list"`,
		`data-calendar-mode="calendar"`,
		`data-calendar-error`,
		"css/calendar-feed",
		"js/calendar-feed",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("default feed missing calendar peer-view hook %q", want)
		}
	}

	for _, excluded := range []string{oldTitle, oldDescription, "/images/old.webp", `data-calendar-post`} {
		if strings.Contains(html, excluded) {
			t.Fatalf("default feed embedded full-history calendar data %q", excluded)
		}
	}
	if !strings.Contains(html, `href="/blog/simple/"`) || !strings.Contains(html, ">Simple</a>") {
		t.Fatalf("default feed no longer exposes the Simple peer view")
	}
}

func TestSimpleFeedLinksToCalendarPeerView(t *testing.T) {
	engine, err := NewEngineWithTheme("", "default")
	if err != nil {
		t.Fatalf("NewEngineWithTheme() error = %v", err)
	}

	title := "Post"
	date := time.Date(2026, time.September, 28, 12, 0, 0, 0, time.UTC)
	feed := &models.FeedConfig{
		Slug:  "blog",
		Title: "Blog",
		Posts: []*models.Post{{Slug: "post", Href: "/post/", Title: &title, Date: &date, Published: true}},
	}
	page := &models.FeedPage{Number: 1, Posts: feed.Posts, TotalPages: 1, TotalItems: 1}
	ctx := NewFeedContext(feed, page, &models.Config{Title: "Test Site", URL: "https://example.com"})
	ctx.Set("feed_stats_total_posts", 1)

	html, err := engine.Render("simple-feed.html", ctx)
	if err != nil {
		t.Fatalf("Render(simple-feed.html) error = %v", err)
	}
	if !strings.Contains(html, `href="/blog/?view=calendar"`) || !strings.Contains(html, ">Calendar</a>") {
		t.Fatalf("simple feed does not link to the calendar peer view")
	}
}

func TestPhotoGridFeedExposesCalendarPeerView(t *testing.T) {
	engine, err := NewEngineWithTheme("", "default")
	if err != nil {
		t.Fatalf("NewEngineWithTheme() error = %v", err)
	}

	title := "Shot"
	date := time.Date(2026, time.September, 28, 12, 0, 0, 0, time.UTC)
	feed := &models.FeedConfig{
		Slug:  "shots",
		Title: "Shots",
		Posts: []*models.Post{{Slug: "shot", Href: "/shot/", Title: &title, Date: &date, Published: true}},
	}
	page := &models.FeedPage{Number: 1, Posts: feed.Posts, TotalPages: 1, TotalItems: 1}
	ctx := NewFeedContext(feed, page, &models.Config{Title: "Test Site", URL: "https://example.com"})
	ctx.Set("feed_stats_total_posts", 1)

	html, err := engine.Render("feed-photo-grid.html", ctx)
	if err != nil {
		t.Fatalf("Render(feed-photo-grid.html) error = %v", err)
	}
	if !strings.Contains(html, `data-calendar-mode="calendar"`) || !strings.Contains(html, ">Grid</button>") {
		t.Fatalf("photo-grid feed does not expose Grid and Calendar peer views")
	}
	if !strings.Contains(html, `data-calendar-src="/shots/calendar.json"`) || strings.Contains(html, `data-calendar-post`) {
		t.Fatalf("photo-grid feed should defer full-history calendar metadata to JSON")
	}
}
