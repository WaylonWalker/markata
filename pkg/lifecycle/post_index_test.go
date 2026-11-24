package lifecycle

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestManagerPostIndex_SetPostsInvalidatesIndex(t *testing.T) {
	m := NewManager()
	oldPost := &models.Post{Slug: "old-post", Href: "/old-post/", Path: "old.md"}
	newPost := &models.Post{Slug: "new-post", Href: "/new-post/", Path: "new.md"}

	m.SetPosts([]*models.Post{oldPost})
	idx1 := m.PostIndex()
	m.SetPosts([]*models.Post{newPost})
	idx2 := m.PostIndex()

	if idx2 == idx1 {
		t.Fatal("PostIndex() after SetPosts() returned the stale index")
	}
	if got := idx2.LookupBySlug("old-post"); got != nil {
		t.Fatalf("old post remained in index after SetPosts(): %v", got)
	}
	if got := idx2.LookupBySlug("new-post"); got != newPost {
		t.Fatalf("new post missing after SetPosts(): got %v, want new post", got)
	}
}

func TestManagerPostIndex_AddPostInvalidatesIndex(t *testing.T) {
	m := NewManager()
	oldPost := &models.Post{Slug: "old-post", Href: "/old-post/", Path: "old.md"}
	newPost := &models.Post{Slug: "new-post", Href: "/new-post/", Path: "new.md"}

	m.SetPosts([]*models.Post{oldPost})
	idx1 := m.PostIndex()
	m.AddPost(newPost)
	idx2 := m.PostIndex()

	if idx2 == idx1 {
		t.Fatal("PostIndex() after AddPost() returned the stale index")
	}
	if got := idx2.LookupBySlug("old-post"); got != oldPost {
		t.Fatalf("old post missing after AddPost(): got %v, want old post", got)
	}
	if got := idx2.LookupBySlug("new-post"); got != newPost {
		t.Fatalf("new post missing after AddPost(): got %v, want new post", got)
	}
}

func TestPostIndex_Lookups(t *testing.T) {
	post := &models.Post{
		Slug: "Hello, World!",
		Href: "/hello-world/",
		Path: "posts/hello.md",
	}
	aliasPost := &models.Post{
		Slug: "alias-owner",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"Old Name"},
		},
	}

	m := NewManager()
	m.SetPosts([]*models.Post{post, aliasPost})
	idx := m.PostIndex()

	for _, test := range []struct {
		name string
		key  string
		want *models.Post
	}{
		{name: "case insensitive slug", key: "HELLO, WORLD!", want: post},
		{name: "slugified slug", key: "hello-world", want: post},
		{name: "case insensitive alias", key: "OLD NAME", want: aliasPost},
		{name: "slugified alias", key: "old-name", want: aliasPost},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := idx.LookupBySlug(test.key); got != test.want {
				t.Fatalf("LookupBySlug(%q) = %v, want %v", test.key, got, test.want)
			}
		})
	}

	if got := idx.ByHref[post.Href]; got != post {
		t.Fatalf("ByHref[%q] = %v, want post", post.Href, got)
	}
	if got := idx.ByPath[post.Path]; got != post {
		t.Fatalf("ByPath[%q] = %v, want post", post.Path, got)
	}
}

func TestPostIndex_DuplicateKeys(t *testing.T) {
	first := &models.Post{
		Slug: "duplicate",
		Href: "/duplicate/",
		Path: "duplicate.md",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"shared alias"},
		},
	}
	last := &models.Post{
		Slug: "duplicate",
		Href: "/duplicate/",
		Path: "duplicate.md",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"shared alias"},
		},
	}

	m := NewManager()
	m.SetPosts([]*models.Post{first, last})
	idx := m.PostIndex()

	if got := idx.LookupBySlug("duplicate"); got != last {
		t.Fatalf("duplicate slug resolved to %v, want last post", got)
	}
	if got := idx.ByHref["/duplicate/"]; got != last {
		t.Fatalf("duplicate href resolved to %v, want last post", got)
	}
	if got := idx.ByPath["duplicate.md"]; got != last {
		t.Fatalf("duplicate path resolved to %v, want last post", got)
	}
	if got := idx.LookupBySlug("shared alias"); got != first {
		t.Fatalf("duplicate alias resolved to %v, want first alias owner", got)
	}
}

func TestPostIndex_RefreshAfterDirectPostMutation(t *testing.T) {
	post := &models.Post{
		Slug: "old-slug",
		Href: "/old-slug/",
		Path: "old.md",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"old alias"},
		},
	}
	m := NewManager()
	m.SetPosts([]*models.Post{post})
	idx := m.PostIndex()

	post.Slug = "new-slug"
	post.Href = "/new-slug/"
	post.Path = "new.md"
	post.Extra["aliases"] = []interface{}{"new alias"}
	idx.Refresh(m)

	if got := idx.LookupBySlug("old-slug"); got != nil {
		t.Fatalf("old slug remained after Refresh(): %v", got)
	}
	if got := idx.LookupBySlug("new-slug"); got != post {
		t.Fatalf("new slug missing after Refresh(): got %v, want post", got)
	}
	if got := idx.LookupBySlug("old alias"); got != nil {
		t.Fatalf("old alias remained after Refresh(): %v", got)
	}
	if got := idx.LookupBySlug("new alias"); got != post {
		t.Fatalf("new alias missing after Refresh(): got %v, want post", got)
	}
	if _, ok := idx.ByHref["/old-slug/"]; ok {
		t.Fatal("old href remained after Refresh()")
	}
	if got := idx.ByHref["/new-slug/"]; got != post {
		t.Fatalf("new href missing after Refresh(): got %v, want post", got)
	}
	if _, ok := idx.ByPath["old.md"]; ok {
		t.Fatal("old path remained after Refresh()")
	}
	if got := idx.ByPath["new.md"]; got != post {
		t.Fatalf("new path missing after Refresh(): got %v, want post", got)
	}
}

func TestPostIndex_NilReceiverLookup(t *testing.T) {
	var idx *PostIndex
	if got := idx.LookupBySlug("missing"); got != nil {
		t.Fatalf("nil PostIndex LookupBySlug() = %v, want nil", got)
	}
}

func TestPostIndex_AliasSliceRepresentations(t *testing.T) {
	for _, test := range []struct {
		name    string
		aliases interface{}
	}{
		{name: "typed", aliases: []string{"Old Name"}},
		{name: "decoded", aliases: []interface{}{"Old Name"}},
		{name: "decoded mixed", aliases: []interface{}{42, "Old Name", nil}},
	} {
		t.Run(test.name, func(t *testing.T) {
			post := &models.Post{Slug: "owner", Extra: map[string]interface{}{"aliases": test.aliases}}
			m := NewManager()
			m.SetPosts([]*models.Post{post})
			idx := m.PostIndex()
			for _, refresh := range []bool{false, true} {
				if refresh {
					idx.Refresh(m)
				}
				for _, alias := range []string{"OLD NAME", "old-name"} {
					if got := idx.LookupBySlug(alias); got != post {
						t.Fatalf("LookupBySlug(%q), refresh=%v: got %v, want owner", alias, refresh, got)
					}
				}
			}
		})
	}
}

func TestPostIndex_RefreshAcrossAliasRepresentationChanges(t *testing.T) {
	post := &models.Post{Slug: "owner", Extra: map[string]interface{}{"aliases": []string{"First Name"}}}
	m := NewManager()
	m.SetPosts([]*models.Post{post})
	idx := m.PostIndex()
	for _, test := range []struct {
		aliases interface{}
		name    string
	}{
		{aliases: []interface{}{"Second Name"}, name: "second-name"},
		{aliases: []string{"Third Name"}, name: "third-name"},
	} {
		post.Extra["aliases"] = test.aliases
		idx.Refresh(m)
		for _, old := range []string{"first-name", "second-name"} {
			if old != test.name && idx.LookupBySlug(old) != nil {
				t.Fatalf("stale alias %q survived refresh", old)
			}
		}
		if got := idx.LookupBySlug(test.name); got != post {
			t.Fatalf("representation change lost %q: %v", test.name, got)
		}
	}
}
