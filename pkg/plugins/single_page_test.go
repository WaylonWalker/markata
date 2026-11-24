package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSinglePagePlugins_ExcludesSiteGenerators(t *testing.T) {
	defaultNames := map[string]bool{}
	for _, p := range DefaultPlugins() {
		defaultNames[p.Name()] = true
	}
	for name := range singlePageExcludedPlugins {
		if !defaultNames[name] {
			t.Errorf("excluded plugin %q is not a default plugin name", name)
		}
	}

	got := map[string]bool{}
	for _, p := range SinglePagePlugins() {
		got[p.Name()] = true
	}
	for name := range singlePageExcludedPlugins {
		if got[name] {
			t.Errorf("SinglePagePlugins() includes excluded plugin %q", name)
		}
	}
	for _, name := range []string{"load", "render_markdown", "templates", "publish_html", "static_assets", "palette_css", "single_page_root"} {
		if !got[name] {
			t.Errorf("SinglePagePlugins() missing %q", name)
		}
	}
}

func TestSinglePageRootPlugin_Load(t *testing.T) {
	m := lifecycle.NewManager()
	post := &models.Post{Path: "pages/solo.md", Slug: "custom-slug", Href: "/custom-slug/"}
	m.SetPosts([]*models.Post{post})

	if err := NewSinglePageRootPlugin().Load(m); err != nil {
		t.Fatal(err)
	}
	if post.Slug != "" || post.Href != "/" || !post.Has("_slug_explicit") {
		t.Errorf("post slug=%q href=%q explicit=%v, want root", post.Slug, post.Href, post.Has("_slug_explicit"))
	}
}
