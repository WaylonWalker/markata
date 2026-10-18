package templates

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestFeedViewOptOutsSuppressCalendarAssetsAndSimpleLink(t *testing.T) {
	engine, err := NewEngineWithTheme("", "default")
	if err != nil {
		t.Fatal(err)
	}
	title := "Post"
	date := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	feed := &models.FeedConfig{
		Slug: "blog", Title: "Blog", Views: []string{models.FeedViewDefault},
		Formats: models.FeedFormats{HTML: true, SimpleHTML: true},
		Posts:   []*models.Post{{Slug: "post", Href: "/post/", Title: &title, Date: &date, Published: true}},
	}
	page := &models.FeedPage{Number: 1, Posts: feed.Posts, TotalPages: 1, TotalItems: 1}
	ctx := NewFeedContext(feed, page, &models.Config{Title: "Test", URL: "https://example.com"})
	html, err := engine.Render("feed.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"css/calendar-feed", "js/calendar-feed", "data-calendar-feed", "?view=calendar", "/blog/simple/"} {
		if strings.Contains(html, forbidden) {
			t.Fatalf("disabled view leaked %q", forbidden)
		}
	}
}

func TestSimpleFeedHidesDisabledCalendarPeer(t *testing.T) {
	engine, err := NewEngineWithTheme("", "default")
	if err != nil {
		t.Fatal(err)
	}
	feed := &models.FeedConfig{Slug: "blog", Views: []string{models.FeedViewDefault, models.FeedViewSimple}}
	page := &models.FeedPage{Number: 1}
	ctx := NewFeedContext(feed, page, &models.Config{Title: "Test"})
	html, err := engine.Render("simple-feed.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "?view=calendar") {
		t.Fatal("calendar link rendered after opt-out")
	}
}
