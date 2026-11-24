package templates

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/WaylonWalker/markata/pkg/models"
)

// Readability-style article extractors should only receive the post, not site
// controls such as the share panel. Keep both copies of the default template
// aligned, and make sure the share controls remain usable on the page.
func TestPostShareOutsideArticleForReaders(t *testing.T) {
	enabled := true
	for _, directory := range []string{"", "../../templates"} {
		t.Run(directory, func(t *testing.T) {
			engine, err := NewEngine(directory)
			if err != nil {
				t.Fatal(err)
			}
			title := "Read this article"
			post := &models.Post{Title: &title, Slug: "read-this", Href: "/read-this/"}
			config := &models.Config{}
			config.Components.Share.Enabled = &enabled
			ctx := NewContext(post, "<p>Actual article body for reading aloud.</p>", config)
			ctx.Set("share_buttons", []map[string]interface{}{{
				"Action": "share", "Link": "https://example.com/share",
				"AriaLabel": "Share on Example", "Name": "Example",
				"Icon": "/example.svg",
			}})
			rendered, err := engine.Render("post.html", ctx)
			if err != nil {
				t.Fatal(err)
			}
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(rendered))
			if err != nil {
				t.Fatal(err)
			}
			article := doc.Find("article.post")
			if article.Length() != 1 || !strings.Contains(article.Text(), "Actual article body for reading aloud.") {
				t.Fatalf("post article missing its text: %s", rendered)
			}
			if article.Find(".share-panel").Length() != 0 || strings.Contains(strings.ToLower(article.Text()), "share this post") {
				t.Fatal("share controls leaked into the article reader boundary")
			}
			share := doc.Find("aside.post-share")
			if share.Length() != 1 || share.Find(".share-panel .share-button").Length() != 1 {
				t.Fatal("share controls disappeared after leaving the article")
			}
			if _, ok := share.Attr("data-pagefind-ignore"); !ok {
				t.Fatal("share controls must not be indexed as article text")
			}
			if share.Parent().ChildrenFiltered("article.post").Length() != 1 {
				t.Fatal("share panel must follow the article in the same post shell")
			}
		})
	}
}
