package templates

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestPostWideWidthOptIn(t *testing.T) {
	for _, directory := range []string{"", "../../templates"} {
		for _, width := range []string{"", "wide", "unknown"} {
			t.Run(directory+"/"+width, func(t *testing.T) {
				engine, err := NewEngine(directory)
				if err != nil {
					t.Fatal(err)
				}
				title := "Story"
				post := &models.Post{
					Title: &title, Slug: "story", Href: "/story/",
					Extra: map[string]interface{}{"content_width": width},
				}
				rendered, err := engine.Render("post.html", NewContext(post, "<p>Opening.</p>", &models.Config{}))
				if err != nil {
					t.Fatal(err)
				}
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(rendered))
				if err != nil {
					t.Fatal(err)
				}
				article := doc.Find("article.post")
				if article.Length() != 1 || article.HasClass("post--wide") != (width == "wide") {
					t.Fatalf("unexpected article classes for width %q: %s", width, rendered)
				}
				_, autoFit := article.Find("h1").Attr("data-title-size")
				if autoFit == (width == "wide") {
					t.Fatalf("wide titles must balance instead of auto-fitting to one line: width %q", width)
				}
			})
		}
	}
}
