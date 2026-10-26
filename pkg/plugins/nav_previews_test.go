package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

func TestNavPreviews_ResolveAndRender(t *testing.T) {
	visible := true
	description := "A focused article"
	public := &models.Post{
		Href: "/article/", Published: true, Description: &description,
		Tags: []string{"go", "notes"}, Date: datePtr(time.January),
		Extra: map[string]interface{}{"word_count": 800, "reading_time": 4},
	}
	second := &models.Post{
		Href: "/another/", Published: true, Date: datePtr(time.February),
		Extra: map[string]interface{}{"word_count": 200, "reading_time": 1},
	}
	private := &models.Post{Href: "/secret/", Published: true, Private: true, Extra: map[string]interface{}{"word_count": 9000}}
	config := models.NewConfig()
	config.Components.Nav.Enabled = &visible
	config.Components.Nav.Items = []models.NavItem{
		{Label: "Journal", URL: "/journal/"},
		{Label: "Article", URL: "/article/"},
		{Label: "Secret", URL: "/secret/"},
		{Label: "Outside", URL: "https://example.com/journal/", External: true},
	}
	feeds := []models.FeedConfig{{Slug: "journal", Description: "Occasional dispatches", Posts: []*models.Post{public, second, private}}}
	previews := buildNavPreviews(config, []*models.Post{public, second, private}, feeds, models.BlogrollConfig{}, RandomPostConfig{})
	if len(previews) != 2 {
		t.Fatalf("got %d previews, want 2", len(previews))
	}
	feed := previews["/journal/"]
	if feed["count"] != 2 || feed["sparkline"] == "" || feed["sparkline_title"] == "" {
		t.Fatalf("unexpected feed preview: %#v", feed)
	}
	for _, volatile := range []string{"words", "minutes", "words_display", "reading_display"} {
		if _, ok := feed[volatile]; ok {
			t.Fatalf("feed preview retained body-derived %q metadata: %#v", volatile, feed)
		}
	}
	if previews["/article/"]["description"] != description {
		t.Fatalf("missing article description: %#v", previews["/article/"])
	}
	if _, ok := previews["/article/"]["words"]; ok {
		t.Fatalf("article preview retained volatile word count: %#v", previews["/article/"])
	}

	engine, err := templates.NewEngine("../../templates")
	if err != nil {
		t.Fatal(err)
	}
	ctx := templates.NewContext(nil, "", config)
	ctx.Set("nav_previews", previews)
	html, err := engine.Render("components/nav.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, fragment := range []string{"Occasional dispatches", "Publication rhythm", "A focused article", "target=\"_blank\""} {
		if !strings.Contains(html, fragment) {
			t.Errorf("rendered nav missing %q", fragment)
		}
	}
	if strings.Contains(html, "1,000</strong> words") || strings.Contains(html, "9000") || strings.Contains(html, "Secret</span>") {
		t.Errorf("body-derived/private metadata appeared in nav: %s", html)
	}
}

func TestNavPreviews_GeneratedRoutesAndFooter(t *testing.T) {
	config := models.NewConfig()
	enabled := true
	config.Components.Nav.Enabled = &enabled
	config.Nav = []models.NavItem{
		{Label: "Reader", URL: "/reading/"},
		{Label: "Random", URL: "/surprise/"},
	}
	config.Components.Footer.Links = []models.NavItem{{Label: "Blogroll", URL: "/following/"}}
	posts := []*models.Post{
		{Href: "/one/", Published: true},
		{Href: "/two/", Published: true, Tags: []string{"exclude"}},
		{Href: "/secret/", Published: true, Private: true},
	}
	blogroll := models.BlogrollConfig{
		Enabled: true, BlogrollSlug: "following", ReaderSlug: "reading",
		Feeds: []models.ExternalFeedConfig{{URL: "https://one.example/feed"}, {URL: "https://two.example/feed"}},
	}
	random := RandomPostConfig{Enabled: true, Path: "surprise", ExcludeTags: []string{"exclude"}}
	previews := buildNavPreviews(config, posts, nil, blogroll, random)
	if previews["/reading/"]["kind"] != "reader" || previews["/reading/"]["count"] != 2 {
		t.Fatalf("reader preview = %#v", previews["/reading/"])
	}
	if previews["/following/"]["kind"] != "followed" || previews["/following/"]["count"] != 2 {
		t.Fatalf("blogroll preview = %#v", previews["/following/"])
	}
	if previews["/surprise/"]["kind"] != "surprise" || previews["/surprise/"]["count"] != 1 {
		t.Fatalf("random preview = %#v", previews["/surprise/"])
	}
	engine, err := templates.NewEngine("../../templates")
	if err != nil {
		t.Fatal(err)
	}
	ctx := templates.NewContext(nil, "", config)
	ctx.Set("nav_previews", previews)
	html, err := engine.Render("components/footer.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Followed sites") || !strings.Contains(html, "2</strong> sources") {
		t.Fatalf("footer link lacks blogroll preview: %s", html)
	}
	navHTML, err := engine.Render("components/nav.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(navHTML, "RSS reader") || !strings.Contains(navHTML, "followed RSS feeds") {
		t.Fatalf("reader link lacks RSS preview: %s", navHTML)
	}
}

func datePtr(month time.Month) *time.Time {
	date := time.Date(2026, month, 1, 0, 0, 0, 0, time.UTC)
	return &date
}
