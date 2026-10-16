package plugins

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/templates"
	"strings"
)

func presentationTestPost(day int, template string) *models.Post {
	date := time.Date(2026, 10, day, 12, 0, 0, 0, time.UTC)
	return &models.Post{Date: &date, Template: template}
}

func TestEditorialGroupsDailyShotsAndPreservesEveryPost(t *testing.T) {
	posts := []*models.Post{
		presentationTestPost(13, "ping"),
		presentationTestPost(12, "photo"),
		presentationTestPost(12, "post"),
		presentationTestPost(12, "photo"),
		presentationTestPost(12, "photo"),
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
			if item.Kind == "shots" {
				clusters++
				if len(item.Posts) != 3 {
					t.Fatalf("shot cluster size = %d; want 3", len(item.Posts))
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

func TestSingleShotNotMadeIntoOneItemGrid(t *testing.T) {
	sections := buildEditorialFeedSections([]*models.Post{
		presentationTestPost(12, "photo"), presentationTestPost(5, "photo"),
	})
	if len(sections) != 2 || sections[0].Items[0].Kind != "post" {
		t.Fatalf("unexpected clustering across weeks: %+v", sections)
	}
}


func TestEditorialTemplateRendersEveryShotThroughCardRouter(t *testing.T) {
	engine, err := templates.NewEngine("../../templates")
	if err != nil {
		t.Fatal(err)
	}
	date := time.Date(2026, 10, 12, 12, 0, 0, 0, time.UTC)
	a := "Bat Cap"
	b := "Ghost Cap"
	posts := []*models.Post{
		{Slug: "bat-cap", Href: "/shots/bat-cap/", Template: "photo", Date: &date, Title: &a,
			Extra: map[string]interface{}{"image": "https://dropper.waylonwalker.com/bat.webp"}},
		{Slug: "ghost-cap", Href: "/shots/ghost-cap/", Template: "photo", Date: &date, Title: &b,
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
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in rendered editorial gallery: %s", want, html)
		}
	}
}
