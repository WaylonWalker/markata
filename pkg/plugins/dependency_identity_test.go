package plugins

import (
	"slices"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestUnresolvedLogicalDependencies(t *testing.T) {
	content := "Missing [[Future Target]] and ![[Another Target]].\n\n```\n[[ignored-target]]\n![[ignored-embed]]\n```\n"

	deps := unresolvedLogicalDependencies(content)
	for _, want := range []string{"future-target", "another-target"} {
		if !slices.Contains(deps, want) {
			t.Fatalf("dependencies = %v, want %q", deps, want)
		}
	}
	for _, unwanted := range []string{"ignored-target", "ignored-embed"} {
		if slices.Contains(deps, unwanted) {
			t.Fatalf("dependencies = %v, must not include fenced-code target %q", deps, unwanted)
		}
	}
}

func TestBuildCacheTransformPersistsUnresolvedTargets(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	plugin := &BuildCachePlugin{cache: cache, enabled: true}
	manager := lifecycle.NewManager()
	manager.SetPosts([]*models.Post{
		{
			Path:    "content/source.md",
			Slug:    "source",
			Content: "See [[Future Target]] and ![[Another Target]]",
		},
	})

	if err := plugin.Transform(manager); err != nil {
		t.Fatalf("Transform() error = %v", err)
	}

	deps := cache.Graph.GetDependencies("content/source.md")
	for _, want := range []string{"future-target", "another-target"} {
		if !slices.Contains(deps, want) {
			t.Fatalf("persisted dependencies = %v, want %q", deps, want)
		}
	}
}

func TestBuildCacheLoadInvalidatesUnresolvedAliasReferrers(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	cache.Posts["content/source-a.md"] = &buildcache.PostCache{InputHash: "source-a-hash", Template: "post.html"}
	cache.Posts["content/source-b.md"] = &buildcache.PostCache{InputHash: "source-b-hash", Template: "post.html"}
	cache.Graph.SetDependencies("content/source-a.md", "source-a", []string{"future-alias"})
	cache.Graph.SetDependencies("content/source-b.md", "source-b", []string{"future-alias"})

	plugin := &BuildCachePlugin{cache: cache, enabled: true}
	manager := lifecycle.NewManager()
	manager.SetPosts([]*models.Post{
		{Path: "content/source-a.md", Slug: "source-a", InputHash: "source-a-hash", Template: "post.html"},
		{Path: "content/source-b.md", Slug: "source-b", InputHash: "source-b-hash", Template: "post.html"},
		{
			Path:      "content/target.md",
			Slug:      "canonical-target",
			InputHash: "target-hash",
			Template:  "post.html",
			Extra: map[string]interface{}{
				"aliases": []interface{}{"Future Alias"},
			},
		},
	})

	if err := plugin.Load(manager); err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	affected := lifecycle.GetServeAffectedPaths(manager)
	for _, path := range []string{"content/target.md", "content/source-a.md", "content/source-b.md"} {
		if !affected[path] {
			t.Fatalf("affected = %v, want %q invalidated", affected, path)
		}
	}
}

func TestBuildCacheLoadResolvedNoOpDoesNotInvalidateReferrers(t *testing.T) {
	cache := buildcache.New(t.TempDir())
	cache.Posts["content/source.md"] = &buildcache.PostCache{InputHash: "source-hash", Template: "post.html"}
	cache.Posts["content/target.md"] = &buildcache.PostCache{InputHash: "target-hash", Template: "post.html"}
	cache.Graph.SetDependencies("content/source.md", "source", []string{"target"})

	plugin := &BuildCachePlugin{cache: cache, enabled: true}
	manager := lifecycle.NewManager()
	manager.SetPosts([]*models.Post{
		{Path: "content/source.md", Slug: "source", InputHash: "source-hash", Template: "post.html"},
		{Path: "content/target.md", Slug: "target", InputHash: "target-hash", Template: "post.html"},
	})

	if err := plugin.Load(manager); err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if affected := lifecycle.GetServeAffectedPaths(manager); len(affected) != 0 {
		t.Fatalf("no-op affected = %v, want none", affected)
	}
}

func TestPostDependencyIdentitiesIncludeSupportedAliases(t *testing.T) {
	post := &models.Post{
		Slug: "canonical-target",
		Extra: map[string]interface{}{
			"aliases": []interface{}{"Future Alias", "legacy/path"},
		},
	}
	identities := postDependencyIdentities(post)
	for _, want := range []string{"canonical-target", "future-alias", "legacy/path"} {
		if !slices.Contains(identities, want) {
			t.Fatalf("identities = %v, want %q", identities, want)
		}
	}
}
