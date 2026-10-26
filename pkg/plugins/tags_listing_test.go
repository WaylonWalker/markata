package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/models"
)

func strPtrTags(s string) *string { return &s }

func TestComputeTagsListingHash_UsesCachedSemanticHashes(t *testing.T) {
	post := models.NewPost("posts/test.md")
	post.Slug = "test"
	post.Href = "/test/"
	post.Title = strPtrTags("Test")
	post.Published = true
	post.Tags = []string{"go", "perf"}

	tagsConfig := models.NewTagsConfig()
	cache := buildcache.New("")
	cache.UpdatePostSemanticHashes(post.Path, "feed-hash", computePostTagIndexHash(post), "garden-hash")

	want := computeTagsListingHash([]*models.Post{post}, &tagsConfig, nil)
	got := computeTagsListingHash([]*models.Post{post}, &tagsConfig, cache)
	if got != want {
		t.Fatalf("computeTagsListingHash with cache = %q, want %q", got, want)
	}
}

func TestComputeTagsListingHash_ChangesWhenConfigChanges(t *testing.T) {
	post := models.NewPost("posts/test.md")
	post.Slug = "test"
	post.Href = "/test/"
	post.Title = strPtrTags("Test")
	post.Published = true
	post.Tags = []string{"go", "perf"}

	base := models.NewTagsConfig()
	baseHash := computeTagsListingHash([]*models.Post{post}, &base, nil)

	modified := models.NewTagsConfig()
	modified.Blacklist = []string{"perf"}
	modifiedHash := computeTagsListingHash([]*models.Post{post}, &modified, nil)

	if baseHash == modifiedHash {
		t.Fatal("computeTagsListingHash did not change after blacklist update")
	}
}
