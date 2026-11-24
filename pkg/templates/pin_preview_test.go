package templates

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"

	"github.com/flosch/pongo2/v6"
)

func TestPinPreview_MediaAndCommentary(t *testing.T) {
	tests := []struct {
		name  string
		data  map[string]interface{}
		image string
		note  string
	}{
		{"explicit image wins", map[string]interface{}{
			"image": "/chosen.webp", "cover": "/cover.webp", "article_html": `<img src="/body.webp"><p>My note.</p>`,
		}, "/chosen.webp", "My note."},
		{"cover fallback", map[string]interface{}{"image": "  ", "cover": "/cover.webp", "cover_image": "/other.webp"}, "/cover.webp", ""},
		{"cover image fallback", map[string]interface{}{"cover_image": "/cover.webp", "og_image": "/og.webp"}, "/cover.webp", ""},
		{"og image fallback", map[string]interface{}{"og_image": "/og.webp"}, "/og.webp", ""},
		{"embed media without metadata commentary", map[string]interface{}{
			"article_html": `<img src="/body.webp"><div class="embed-card"><img src="https://cdn.example/cover?a=1&amp;b=2"><p>Destination description</p></div><p>My <em>own</em> note &amp; reasons.</p>`,
		}, "https://cdn.example/cover?a=1&b=2", "My own note & reasons."},
		{"youtube thumbnail", map[string]interface{}{
			"article_html": `<div class="embed-card"><lite-youtube videoid="9fcfDF8SBnc"></lite-youtube></div><p>Useful talk.</p>`,
		}, "https://i.ytimg.com/vi/9fcfDF8SBnc/hqdefault.jpg", "Useful talk."},
		{"invalid youtube id", map[string]interface{}{"article_html": `<lite-youtube videoid="../invalid/"></lite-youtube>`}, "", ""},
		{"body image and list notes", map[string]interface{}{
			"article_html": `<img class="link-avatar" src="/icon.ico"><img src="/photo.webp"><ul><li>First point</li><li><p>Second point</p></li></ul>`,
		}, "/photo.webp", "First point\n\nSecond point"},
		{"skip code and standalone urls", map[string]interface{}{
			"article_html": `<pre><p>Not commentary</p></pre><p>!https://example.com</p><p>https://example.com</p><p>Actual commentary.</p>`,
		}, "", "Actual commentary."},
		{"description fallback", map[string]interface{}{
			"description": "A useful bookmark.", "article_html": `<div class="embed-card"><p>Destination</p></div>`,
		}, "", "A useful bookmark."},
		{"url description omitted", map[string]interface{}{"description": "!https://example.com"}, "", ""},
		{"no data", nil, "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filterPinPreview(pongo2.AsValue(tt.data), nil)
			if err != nil {
				t.Fatal(err)
			}
			got, ok := result.Interface().(map[string]interface{})
			if !ok {
				t.Fatalf("preview type = %T, want map", result.Interface())
			}
			if got["image"] != tt.image || got["commentary"] != tt.note {
				t.Fatalf("preview = %#v, want image %q and note %q", got, tt.image, tt.note)
			}
		})
	}
}

func TestPinsTemplate_RendersPreviewsWithoutJavaScript(t *testing.T) {
	for _, directory := range []string{"", "../../templates"} {
		t.Run(directory, func(t *testing.T) {
			engine, err := NewEngineWithTheme(directory, "default")
			if err != nil {
				t.Fatal(err)
			}
			title := "Saved page"
			post := &models.Post{
				Slug: "saved-page", Href: "/saved-page/", Published: true, Title: &title,
				Extra:       map[string]interface{}{"link": "https://example.com/source"},
				ArticleHTML: `<div class="embed-card"><img src="/embed.webp"><p>Destination metadata</p></div><p>Use &lt;script&gt; as text.</p>`,
			}
			posts := []*models.Post{post}
			feed := &models.FeedConfig{Slug: "pins", Title: "Pins", Posts: posts}
			page := &models.FeedPage{Number: 1, Posts: posts, TotalPages: 1, TotalItems: 1}
			ctx := NewFeedContext(feed, page, &models.Config{URL: "https://example.com"})
			output, err := engine.Render("pins.html", ctx)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{`src="/embed.webp`, `<details class="pin-commentary">`, `Use &lt;script&gt; as text.`, `href="/saved-page/"`, `href="https://example.com/source"`, "Saved links"} {
				if !strings.Contains(output, want) {
					t.Errorf("rendered Pins missing %q", want)
				}
			}
			for _, unwanted := range []string{"Destination metadata", "Worth keeping", "pin-preview-source", "js/pins.js", "Use <script> as text."} {
				if strings.Contains(output, unwanted) {
					t.Errorf("rendered Pins contains %q", unwanted)
				}
			}
		})
	}
}
