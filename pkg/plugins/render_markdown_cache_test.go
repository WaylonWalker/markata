package plugins

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestRenderMarkdownPlugin_CachedHeadingHighlightsMatchFresh(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "cache")
	cache := buildcache.New(dir)
	const content = "# ==First==\n\n## A ==**second**== heading\n\nProse ==ordinary==.\n"
	var expected string
	for iteration := 0; iteration < 4; iteration++ {
		m := lifecycle.NewManager()
		m.Cache().Set("build_cache", cache)
		post := models.NewPost("original.md")
		post.Content = content
		m.AddPost(post)
		if iteration > 0 {
			probe := models.NewPost("new.md")
			probe.Content = "Unrelated new content."
			m.AddPost(probe)
		}
		renderer := NewRenderMarkdownPlugin()
		if err := renderer.Configure(m); err != nil {
			t.Fatal(err)
		}
		if err := renderer.Render(m); err != nil {
			t.Fatal(err)
		}
		if iteration == 0 {
			expected = post.ArticleHTML
			if strings.Count(expected, `class="heading-highlight"`) != 2 {
				t.Fatalf("fresh heading wrappers incorrect: %s", expected)
			}
		} else if post.ArticleHTML != expected {
			t.Fatalf("cached article changed after adding unrelated content (iteration %d):\ngot: %s\nwant: %s",
				iteration, post.ArticleHTML, expected)
		}
		if cached := cache.GetCachedArticleHTML(post.Path, buildcache.ContentHash(content)); cached != expected {
			t.Fatalf("cache no longer contains the fresh render-stage result: %s", cached)
		}
		if err := cache.Save(); err != nil {
			t.Fatal(err)
		}
		var err error
		cache, err = buildcache.Load(dir)
		if err != nil {
			t.Fatal(err)
		}
	}
}
