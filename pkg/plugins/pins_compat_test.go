package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSubscriptionFeedsPlugin_Collect_AddsPinsFeed(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"title": "Pins Site"}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	if err := NewSubscriptionFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	for _, feed := range getFeedConfigs(m.Config()) {
		if feed.Slug != pinsFeedSlug {
			continue
		}
		if !feed.Formats.HTML || feed.Templates.HTML != "pins.html" {
			t.Fatalf("pins feed = %#v, want HTML pins template", feed)
		}
		if feed.Filter != "published == true and (link or url)" {
			t.Fatalf("pins filter = %q", feed.Filter)
		}
		if feed.ItemsPerPage != 100 || feed.PaginationType != models.PaginationManual {
			t.Fatalf("pins pagination = %d / %q, want 100 / manual", feed.ItemsPerPage, feed.PaginationType)
		}
		return
	}
	t.Fatal("implicit pins feed was not injected")
}

func TestSubscriptionFeedsPlugin_Collect_PinsFilterSelectsSavedLinks(t *testing.T) {
	config := lifecycle.NewConfig()
	m := lifecycle.NewManager()
	m.SetConfig(config)
	if err := NewSubscriptionFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	var pinsFilter *filter.Filter
	for _, feed := range getFeedConfigs(m.Config()) {
		if feed.Slug == pinsFeedSlug {
			var err error
			pinsFilter, err = filter.Parse(feed.Filter)
			if err != nil {
				t.Fatalf("Parse(%q) error = %v", feed.Filter, err)
			}
			break
		}
	}
	if pinsFilter == nil {
		t.Fatal("implicit pins feed was not injected")
	}

	for _, tt := range []struct {
		name      string
		published bool
		extra     map[string]interface{}
		want      bool
	}{
		{name: "link", published: true, extra: map[string]interface{}{"link": "https://example.com/link"}, want: true},
		{name: "url", published: true, extra: map[string]interface{}{"url": "https://example.com/url"}, want: true},
		{name: "both", published: true, extra: map[string]interface{}{"link": "https://example.com/link", "url": "https://example.com/url"}, want: true},
		{name: "empty link", published: true, extra: map[string]interface{}{"link": ""}},
		{name: "empty url", published: true, extra: map[string]interface{}{"url": ""}},
		{name: "no source", published: true},
		{name: "unpublished", extra: map[string]interface{}{"url": "https://example.com/url"}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pinsFilter.Match(&models.Post{Published: tt.published, Extra: tt.extra})
			if err != nil {
				t.Fatalf("Match() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("Match() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestSubscriptionFeedsPlugin_Collect_PreservesConfiguredPins(t *testing.T) {
	want := models.FeedConfig{
		Slug: "pins", Title: "Reading shelf", ItemsPerPage: 7,
		PaginationType: models.PaginationJS,
		Templates:      models.FeedTemplates{HTML: "custom-pins.html"},
	}
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"feeds": []models.FeedConfig{want}}
	m := lifecycle.NewManager()
	m.SetConfig(config)
	if err := NewSubscriptionFeedsPlugin().Collect(m); err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, feed := range getFeedConfigs(m.Config()) {
		if feed.Slug == "pins" {
			count++
			if feed.Title != want.Title || feed.ItemsPerPage != want.ItemsPerPage ||
				feed.PaginationType != want.PaginationType || feed.Templates != want.Templates {
				t.Fatalf("configured pins changed: %#v", feed)
			}
		}
	}
	if count != 1 {
		t.Fatalf("configured pins count = %d, want 1", count)
	}
}

func TestSubscriptionFeedsPlugin_Collect_PreservesExistingPinsPost(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{"title": "Existing Site"}

	m := lifecycle.NewManager()
	m.SetConfig(config)
	m.SetPosts([]*models.Post{{
		Path:      "pages/pins.md",
		Slug:      pinsFeedSlug,
		Published: true,
	}})

	if err := NewSubscriptionFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	for _, feed := range getFeedConfigs(m.Config()) {
		if feed.Slug == pinsFeedSlug {
			t.Fatalf("implicit pins feed was injected over existing /pins/ post: %#v", feed)
		}
	}
}

func TestPostOwnsFeedSlug_IgnoresNonOutputPosts(t *testing.T) {
	posts := []*models.Post{
		{Slug: pinsFeedSlug, Published: true, Draft: true},
		{Slug: "/pins/", Published: true, Skip: true},
		{Slug: pinsFeedSlug, Published: false},
		{Slug: "other", Published: true},
	}
	if postOwnsFeedSlug(posts, pinsFeedSlug) {
		t.Fatal("draft, skipped, or unpublished posts should not reserve the feed output path")
	}
}
