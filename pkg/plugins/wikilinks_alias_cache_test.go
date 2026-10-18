package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestWikilinksPlugin_CacheRestoredAliasRepresentations(t *testing.T) {
	for _, aliases := range []interface{}{[]string{"Old Name"}, []interface{}{"Old Name"}} {
		manager := lifecycle.NewManager()
		target := &models.Post{
			Slug: "person", Href: "/person/",
			Extra: map[string]interface{}{"aliases": aliases},
		}
		source := &models.Post{
			Slug: "notes", Href: "/notes/", Content: "See [[ Old Name ]]",
			Extra: make(map[string]interface{}),
		}
		manager.SetPosts([]*models.Post{target, source})
		if err := NewWikilinksPlugin().Transform(manager); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(source.Content, `href="/person/"`) || strings.Contains(source.Content, "[[") {
			t.Fatalf("alias representation %T lost link: %s", aliases, source.Content)
		}
	}
}
