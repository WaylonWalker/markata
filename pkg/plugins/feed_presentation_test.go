package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

func presentationTestPost(day int, template string) *models.Post {
	date := time.Date(2026, 10, day, 12, 0, 0, 0, time.UTC)
	return &models.Post{Date: &date, Template: template}
}

func TestEditorialGroupsWeeklyShotsAndPreservesEveryPost(t *testing.T) {
	posts := []*models.Post{
		presentationTestPost(13, "ping"),
		presentationTestPost(12, "photo"),
		presentationTestPost(12, "post"),
		presentationTestPost(14, "photo"),
		presentationTestPost(16, "photo"),
		presentationTestPost(5, "post"),
	}
	sections := buildEditorialFeedSections(posts)
	if len(sections) != 2 {
		t.Fatalf("weekly sections = %d; want 2", len(sections))
	}
	total, clusters := 0, 0
	for _, section := range sections {
		for _, item := range section.Items {
			total += len(item.Posts)
			if item.Kind == embedOptionShots {
				clusters++
				if len(item.Posts) != 3 {
					t.Fatalf("shot cluster size = %d; want 3", len(item.Posts))
				}
				if !item.Mosaic || item.DateLabel != section.Label {
					t.Fatalf("weekly photo collection lost its week label or mosaic: %+v", item)
				}
				for j, want := range []string{"October 12, 2026", "October 14, 2026", "October 16, 2026"} {
					if got := item.Posts[j]["editorial_photo_date"]; got != want {
						t.Errorf("photo %d date = %v, want %s", j, got, want)
					}
				}
			}
		}
	}
	if total != len(posts) || clusters != 1 {
		t.Fatalf("posts %d/%d, shot clusters %d", total, len(posts), clusters)
	}
	if sections[0].Items[0].Importance != "featured" {
		t.Fatalf("expected article ahead of ordinary posts: %+v", sections[0].Items)
	}
	if sections[1].Start != "2026-10-05" || sections[1].RelativeLabel == "" {
		t.Fatalf("section missing week start or relative label: %+v", sections[1])
	}
	if got := sections[1].Label; got != "Oct 5 – Oct 11, 2026" {
		t.Fatalf("previous week label = %q", got)
	}
}

func TestAdaptiveFeedClassificationAndLegacyDefault(t *testing.T) {
	fc := &models.FeedConfig{Presentation: "auto"}
	for n := 0; n < 12; n++ {
		fc.Posts = append(fc.Posts, presentationTestPost(12, "photo"))
	}
	if got := resolveFeedPresentation(fc); got != "gallery" {
		t.Fatalf("photo feed auto = %s; want gallery", got)
	}
	fc.Templates.HTML = "feed-photo-grid.html"
	if got := resolveFeedPresentation(fc); got != "classic" {
		t.Fatalf("custom template must win: %s", got)
	}
	fc.Templates.HTML = ""
	fc.Presentation = ""
	if got := resolveFeedPresentation(fc); got != "classic" {
		t.Fatalf("default must remain classic: %s", got)
	}
	fc.Presentation = "editorial"
	if got := resolveFeedPresentation(fc); got != "editorial" {
		t.Fatalf("explicit selection ignored: %s", got)
	}
}

func TestSingleShotKeepsGallerySizingWithinItsWeek(t *testing.T) {
	sections := buildEditorialFeedSections([]*models.Post{
		presentationTestPost(12, "photo"), presentationTestPost(5, "photo"),
	})
	if len(sections) != 2 {
		t.Fatalf("weekly sections = %d, want 2", len(sections))
	}
	for _, section := range sections {
		if len(section.Items) != 1 || section.Items[0].Kind != embedOptionShots ||
			len(section.Items[0].Posts) != 1 || section.Items[0].Mosaic {
			t.Fatalf("single photo should remain a one-cell gallery, not a full-width card: %+v", section)
		}
	}
}

func TestEditorialTemplateRendersEveryShotThroughCardRouter(t *testing.T) {
	engine, err := templates.NewEngine("../../templates")
	if err != nil {
		t.Fatal(err)
	}
	date := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	otherDate := time.Date(2026, 10, 14, 12, 0, 0, 0, time.UTC)
	a := "Bat Cap"
	b := "Ghost Cap"
	posts := []*models.Post{
		{Slug: "bat-cap", Href: "/shots/bat-cap/", Template: "photo", Date: &date, Title: &a,
			Extra: map[string]interface{}{"image": "https://dropper.waylonwalker.com/bat.webp"}},
		{Slug: "ghost-cap", Href: "/shots/ghost-cap/", Template: "photo", Date: &otherDate, Title: &b,
			Extra: map[string]interface{}{"image": "https://dropper.waylonwalker.com/ghost.webp"}},
	}
	ctx := map[string]interface{}{
		"feed_sections": buildEditorialFeedSections(posts),
		"config": map[string]interface{}{
			"components": map[string]interface{}{
				"card_router": map[string]interface{}{
					"mappings": map[string]interface{}{"photo": "photo"},
				},
			},
		},
	}
	html, err := engine.RenderToString("partials/feed-editorial.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`class="editorial-feed posts-list"`,
		`class="editorial-shot-grid"`,
		`href="/shots/bat-cap/"`,
		`href="/shots/ghost-cap/"`,
		`datetime="2026-10-12"`,
		`datetime="2026-10-14"`,
		`October 12, 2026`,
		`October 14, 2026`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in rendered editorial gallery: %s", want, html)
		}
	}
}

func TestEditorialLeadPhotoEssayAndNotebook(t *testing.T) {
	featured := presentationTestPost(12, "post")
	featured.Extra = map[string]interface{}{"featured": true}
	article := presentationTestPost(13, "post")
	link := presentationTestPost(14, "link")
	thought := presentationTestPost(14, "thought")
	thought.Path = "pages/thoughts/note.md"
	posts := []*models.Post{
		thought, link,
		presentationTestPost(16, "photo"),
		presentationTestPost(16, "photo"),
		presentationTestPost(16, "photo"),
		article, featured,
	}
	sections := buildEditorialFeedSections(posts)
	if len(sections) != 1 {
		t.Fatalf("sections = %d, want 1", len(sections))
	}
	leadCount, photoCount, quietCount := 0, 0, 0
	for i, item := range sections[0].Items {
		if item.Lead {
			leadCount++
			if i != 0 || item.TypeLabel != "Writing" {
				t.Errorf("expected leading featured writing, got item %d: %+v", i, item)
			}
		}
		if item.Kind == embedOptionShots {
			photoCount++
			if !item.Mosaic || item.DateLabel != "Oct 12 – Oct 18, 2026" || len(item.Posts) != 3 {
				t.Errorf("wrong photo essay: %+v", item)
			}
		}
		if item.Importance == presentationQuiet {
			quietCount++
			if item.TypeLabel != "Note" {
				t.Errorf("thought misclassified: %+v", item)
			}
		}
	}
	if leadCount != 1 || photoCount != 1 || quietCount != 1 {
		t.Fatalf("weekly lead count %d, photo essays %d, quiet thoughts %d", leadCount, photoCount, quietCount)
	}
	if len(sections[0].Items) != 5 {
		t.Fatalf("weekly items = %d, want 5", len(sections[0].Items))
	}
}

func TestEditorialNoArtificialNoteLead(t *testing.T) {
	sections := buildEditorialFeedSections([]*models.Post{presentationTestPost(12, "note")})
	if len(sections) != 1 || len(sections[0].Items) != 1 || sections[0].Items[0].Lead {
		t.Fatalf("note-only week should not invent lead story: %+v", sections)
	}
}

func TestEditorialRelativeWeekLabels(t *testing.T) {
	now := time.Date(2026, 10, 18, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct{ start, want string }{
		{"2026-10-12", "This week"},
		{"2026-10-05", "Last week"},
		{"2026-09-28", "2 weeks ago"},
		{"2026-09-14", "4 weeks ago"},
		{"2026-08-17", "2 months ago"},
		{"2025-10-13", "1 year ago"},
		{"2024-10-14", "2 years ago"},
		{"2026-10-19", "Next week"},
		{"2026-10-26", "In 2 weeks"},
		{"undated", "Undated"},
	} {
		t.Run(tc.start, func(t *testing.T) {
			if got := editorialRelativeWeekLabel(tc.start, now); got != tc.want {
				t.Fatalf("relative label for %q = %q, want %q", tc.start, got, tc.want)
			}
		})
	}
	// Calendar weeks, not elapsed seven-day periods, define the boundary.
	monday := now.AddDate(0, 0, 1)
	if got := editorialRelativeWeekLabel("2026-10-12", monday); got != "Last week" {
		t.Fatalf("Monday boundary = %q, want Last week", got)
	}
}
