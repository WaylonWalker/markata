package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

func TestLocalLinkPreviews_ResolveAndAnnotate(t *testing.T) {
	title := "Start <here>"
	public := &models.Post{Href: "/start/", Title: &title, Published: true, Tags: []string{"Go"}, Extra: map[string]interface{}{"word_count": 850, "reading_time": 4}}
	private := &models.Post{Href: "/secret/", Published: true, Private: true}
	feed := models.FeedConfig{Slug: "blog", Title: "The blog", Posts: []*models.Post{public, private}}
	previews := encodeLocalPreviews(buildLocalPreviews([]*models.Post{public, private}, []models.FeedConfig{feed}))
	article := `<a href="http://localhost:8000/blog/">Blog</a> <a href="/start/" class="wikilink">Start</a> <a href="/secret/">Secret</a> <a href="https://other.example/blog/">Other</a>`
	got := annotateLocalPreviewLinks(article, "/entry/", "https://example.com", previews)
	if strings.Count(got, "data-local-preview=") != 2 || !strings.Contains(got, `Start \u003chere\u003e`) || !strings.Contains(got, `&#34;count&#34;:1`) {
		t.Fatalf("wrong local annotations: %s", got)
	}
	if strings.Contains(got, "secret&quot;") || strings.Contains(got, "other.example/blog/&quot;") {
		t.Fatalf("private or external preview leaked: %s", got)
	}
	if path := localPreviewPath("../blog/#recent", "/entry/deep/", "https://example.com"); path != "/entry/blog/" {
		t.Fatalf("relative path = %q", path)
	}
	if path := localPreviewPath("https://other.example/blog/", "/entry/", "https://example.com"); path != "" {
		t.Fatalf("external path = %q", path)
	}
}

func TestLocalLinkPreviews_AutoTagsAndHash(t *testing.T) {
	title := "A post"
	post := &models.Post{Href: "/entry/", Published: true, Title: &title, Tags: []string{"Go"}, Date: datePtr(time.January), ArticleHTML: `<a href="/tags/go/">Go</a>`}
	private := &models.Post{Href: "/secret/", Published: true, Private: true, Tags: []string{"Go"}}
	config := &lifecycle.Config{Extra: map[string]interface{}{}}
	feeds := previewFeeds(config, []*models.Post{post, private}, nil)
	if len(feeds) != 1 || feeds[0].Slug != "tags/go" || len(feeds[0].Posts) != 1 {
		t.Fatalf("tag feeds = %#v", feeds)
	}
	encoded := encodeLocalPreviews(buildLocalPreviews([]*models.Post{post, private}, feeds))
	if localPreviewForTag("Go", encoded) == "" {
		t.Fatal("tag preview missing")
	}
	before := localPreviewHash(post, encoded, "https://example.com")
	second := &models.Post{Href: "/second/", Published: true, Tags: []string{"Go"}, Date: datePtr(time.February)}
	updated := previewFeeds(config, []*models.Post{post, second, private}, nil)
	encoded = encodeLocalPreviews(buildLocalPreviews([]*models.Post{post, second, private}, updated))
	if before == localPreviewHash(post, encoded, "https://example.com") {
		t.Fatal("tag feed change did not invalidate referencing post")
	}
}

func TestLocalLinkPreviews_TagAttributeEscaped(t *testing.T) {
	engine, err := templates.NewEngine("../../templates")
	if err != nil {
		t.Fatal(err)
	}
	ctx := templates.NewContext(nil, "", models.NewConfig())
	ctx.Set("tag_preview_json", map[string]string{"Go": `{"title":"A \"quote\""}`})
	ctx.Set("tag", "Go")
	got, err := engine.RenderString(`<a data-local-preview="{{ tag_preview_json[tag] }}">Tag</a>`, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, `&quot;title&quot;`) {
		t.Fatalf("preview JSON is not attribute-escaped: %s", got)
	}
}
