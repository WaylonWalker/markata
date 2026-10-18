package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

func TestFeedHelpers_FeedPostsCallable(t *testing.T) {
	m := newFeedTestManager(t)
	callable := createFeedPostsFunc(m)
	posts, err := callable("blog", map[string]interface{}{"limit": 1})
	if err != nil {
		t.Fatalf("feed_posts returned error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if posts[0]["slug"] != "second" {
		t.Fatalf("unexpected post slug: %v", posts[0]["slug"])
	}
}

func TestFeedHelpers_RenderFeedCallable(t *testing.T) {
	m := newFeedTestManager(t)
	engine, err := templates.NewEngine("../../templates")
	if err != nil {
		t.Fatalf("templates.NewEngine: %v", err)
	}
	m.Cache().Set("templates.engine", engine)

	render := createRenderFeedFunc(m)
	value, err := render("blog", map[string]interface{}{"limit": 2})
	if err != nil {
		t.Fatalf("render_feed returned error: %v", err)
	}
	output := value.String()
	if !strings.Contains(output, "feed h-feed") {
		t.Fatalf("rendered snippet missing expected wrapper: %s", output)
	}
}

func TestFeedHelpers_RenderFeedWithoutCachedEngine(t *testing.T) {
	m := newFeedTestManager(t)
	config := m.Config()
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}
	config.Extra["theme"] = ThemeDefault

	render := createRenderFeedFunc(m)
	value, err := render("blog", map[string]interface{}{"limit": 1})
	if err != nil {
		t.Fatalf("render_feed returned error: %v", err)
	}
	output := value.String()
	if !strings.Contains(output, "feed h-feed") {
		t.Fatalf("rendered snippet missing expected wrapper: %s", output)
	}
}

func TestFeedHelpers_RenderFeedPhotoCard(t *testing.T) {
	title := "Wicket's Lab"
	description := "Wicket in his lab"
	date := time.Now()
	photo := &models.Post{
		Slug:        "lab",
		Href:        "/lab/",
		Published:   true,
		Date:        &date,
		Template:    "photo",
		Title:       &title,
		Description: &description,
		Extra: map[string]interface{}{
			"image": "https://dropper.waylonwalker.com/lab.webp?token=abc&w=600",
		},
	}
	feed := models.FeedConfig{
		Slug:    "photo-feed",
		Title:   "Photos",
		Filter:  "template == 'photo'",
		Sort:    "date",
		Reverse: true,
	}
	m := newFeedTestManagerWithCustom(t, []models.FeedConfig{feed}, []*models.Post{photo})
	engine, err := templates.NewEngine("../../templates")
	if err != nil {
		t.Fatalf("templates.NewEngine: %v", err)
	}
	m.Cache().Set("templates.engine", engine)

	render := createRenderFeedFunc(m)
	value, err := render("photo-feed", map[string]interface{}{"limit": 1})
	if err != nil {
		t.Fatalf("render_feed returned error: %v", err)
	}
	output := value.String()
	if !strings.Contains(output, "photo-figure") {
		t.Fatalf("expected photo figure class, got: %s", output)
	}
	if !strings.Contains(output, `aria-label="Wicket in his lab"`) {
		t.Fatalf("expected accessible name on media anchor, got: %s", output)
	}
	if !strings.Contains(output, "<figure") || !strings.Contains(output, "<figcaption") {
		t.Fatalf("expected figure markup for photo card, got: %s", output)
	}
	if !strings.Contains(output, description) {
		t.Fatalf("expected photo description in figcaption, got: %s", output)
	}
	if !strings.Contains(output, `srcset="https://dropper.waylonwalker.com/lab.webp?token=abc&amp;w=160 160w`) {
		t.Fatalf("expected responsive candidates for trusted remote image, got: %s", output)
	}
	if !strings.Contains(output, `data-href="https://dropper.waylonwalker.com/lab.webp?token=abc"`) {
		t.Fatalf("expected full-size URL for the lightbox, got: %s", output)
	}
}

func TestFeedHelpers_RenderPhotoFigureAddsAccessibleName(t *testing.T) {
	output := renderPhotoFigure(map[string]interface{}{
		"href":        "/lab/",
		"image":       "https://dropper.waylonwalker.com/lab.webp?token=abc&w=400&h=200",
		"description": "Wicket in his lab",
	})

	if !strings.Contains(output, `aria-label="Wicket in his lab"`) {
		t.Fatalf("expected aria-label on fallback media anchor, got: %s", output)
	}
	if !strings.Contains(output, `src="https://dropper.waylonwalker.com/lab.webp?token=abc&amp;w=1200"`) {
		t.Fatalf("expected width-only image variant with access token preserved, got: %s", output)
	}
	if strings.Contains(output, `height="675"`) || strings.Contains(output, `h=675`) {
		t.Fatalf("photo fallback must not force a 16:9 aspect ratio: %s", output)
	}
}

func newFeedTestManager(t *testing.T) *lifecycle.Manager {
	date := time.Now()
	older := date.Add(-time.Hour)
	posts := []*models.Post{
		{
			Slug:      "first",
			Href:      "/first/",
			Published: true,
			Date:      &older,
		},
		{
			Slug:      "second",
			Href:      "/second/",
			Published: true,
			Date:      &date,
		},
	}
	feeds := []models.FeedConfig{
		{
			Slug:    "blog",
			Title:   "Blog",
			Filter:  "published == true",
			Sort:    "date",
			Reverse: true,
		},
	}
	return newFeedTestManagerWithCustom(t, feeds, posts)
}

func newFeedTestManagerWithCustom(t *testing.T, feeds []models.FeedConfig, posts []*models.Post) *lifecycle.Manager {
	t.Helper()
	m := lifecycle.NewManager()
	config := m.Config()
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}
	config.Extra["feeds"] = feeds
	m.SetPosts(posts)
	return m
}

func TestNormalizeTemplateWhitespace_SingleHTMLBlock(t *testing.T) {
	in := "\n  <div class=\"feed\">\n\n    <article>\n\n\n      <p>one</p>\n    </article>\n\n</div>\n\n"
	got := normalizeTemplateWhitespace(in)
	if strings.Contains(got, "\n\n") {
		t.Fatalf("expected no blank lines so goldmark keeps one HTML block, got %q", got)
	}
	if strings.Contains(got, "\n ") || strings.HasPrefix(got, " ") {
		t.Fatalf("expected left-trimmed lines, got %q", got)
	}
	want := "<div class=\"feed\">\n<article>\n<p>one</p>\n</article>\n</div>"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
