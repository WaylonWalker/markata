package plugins

import (
	"fmt"
	"sync"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestTemplatesPlugin_SidebarProjectionParity(t *testing.T) {
	plugin := NewTemplatesPlugin()
	for _, href := range []string{
		"/post/", "/post/?z=1&feed=old#section", "/caf\u00e9/",
		"https://example.com/post?a=two+words", "/bad%ZZ",
	} {
		post := &models.Post{Slug: "post", Href: href, TitleText: "Plain title"}
		for _, feed := range []string{"", "tags/go", "a b&c", "caf\u00e9"} {
			for _, active := range []bool{true, false, true} {
				if got, want := plugin.postToSidebarJSON(post, active, feed), postToSidebarJSON(post, active, feed); got != want {
					t.Fatalf("href=%q feed=%q active=%v: got %#v, want %#v", href, feed, active, got, want)
				}
			}
		}
	}
}

func TestTemplatesPlugin_SidebarProjectionIdentityAndReset(t *testing.T) {
	plugin := NewTemplatesPlugin()
	post := &models.Post{Slug: "old", Href: "/old/"}
	plugin.postToSidebarJSON(post, true, "feed")
	post.Slug = "new"
	post.Href = "/new/?x=1"
	post.TitleText = "New title"
	got := plugin.postToSidebarJSON(post, false, "other")
	if want := postToSidebarJSON(post, false, "other"); got != want {
		t.Fatalf("changed source returned stale projection: got %#v, want %#v", got, want)
	}
	plugin.resetSidebarPosts()
	if len(plugin.sidebarPosts) != 0 {
		t.Fatal("new render retained old sidebar projections")
	}
	if got := plugin.postToSidebarJSON(post, true, "other"); !got.Active {
		t.Fatal("reset lost page-specific active state")
	}
}

func TestTemplatesPlugin_SidebarProjectionConcurrentActiveState(t *testing.T) {
	plugin := NewTemplatesPlugin()
	post := &models.Post{Slug: "post", Href: "/post/", TitleText: "Title"}
	var workers sync.WaitGroup
	for i := range 64 {
		workers.Add(1)
		go func(active bool) {
			defer workers.Done()
			for range 100 {
				if got, want := plugin.postToSidebarJSON(post, active, "feed"), postToSidebarJSON(post, active, "feed"); got != want {
					t.Errorf("shared active state: got %#v, want %#v", got, want)
					return
				}
			}
		}(i%2 == 0)
	}
	workers.Wait()
}

func TestTemplatesPlugin_SidebarProjectionEachIdentityField(t *testing.T) {
	for _, field := range []string{"slug", "title", "href", "feed"} {
		t.Run(field, func(t *testing.T) {
			plugin := NewTemplatesPlugin()
			post := &models.Post{Slug: "post", Href: "/post/", TitleText: "Title"}
			feed := "feed"
			plugin.postToSidebarJSON(post, false, feed)
			switch field {
			case "slug":
				post.Slug = "new"
			case "title":
				post.TitleText = "Changed title"
			case "href":
				post.Href = "/changed/"
			case "feed":
				feed = "changed-feed"
			}
			if got, want := plugin.postToSidebarJSON(post, true, feed), postToSidebarJSON(post, true, feed); got != want {
				t.Fatalf("changed %s retained old projection: got %#v, want %#v", field, got, want)
			}
		})
	}
}

func BenchmarkSidebarPostProjection(b *testing.B) {
	post := &models.Post{Slug: "post", Href: "/post/?existing=one#section", TitleText: "Title"}
	for _, cached := range []bool{false, true} {
		b.Run(fmt.Sprintf("cached=%v", cached), func(b *testing.B) {
			plugin := NewTemplatesPlugin()
			plugin.postToSidebarJSON(post, false, "tags/go")
			b.ReportAllocs()
			b.ResetTimer()
			for i := range b.N {
				if cached {
					plugin.postToSidebarJSON(post, i%2 == 0, "tags/go")
				} else {
					postToSidebarJSON(post, i%2 == 0, "tags/go")
				}
			}
		})
	}
}
