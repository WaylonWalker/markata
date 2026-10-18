package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSlugConflictsPlugin_Name(t *testing.T) {
	p := NewSlugConflictsPlugin()
	if got := p.Name(); got != "slug_conflicts" {
		t.Errorf("Name() = %q, want %q", got, "slug_conflicts")
	}
}

func TestSlugConflictsPlugin_Priority(t *testing.T) {
	p := NewSlugConflictsPlugin()

	// Should have highest priority (first) during collect stage
	if got := p.Priority(lifecycle.StageCollect); got != lifecycle.PriorityFirst {
		t.Errorf("Priority(StageCollect) = %d, want %d", got, lifecycle.PriorityFirst)
	}

	// Should have default priority for other stages
	if got := p.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestSlugConflictsPlugin_NoConflicts(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add posts with unique slugs
	post1 := models.NewPost("posts/first.md")
	post1.Slug = "first"
	post1.Published = true

	post2 := models.NewPost("posts/second.md")
	post2.Slug = "second"
	post2.Published = true

	m.AddPost(post1)
	m.AddPost(post2)

	// No feeds configured
	err := p.Collect(m)
	if err != nil {
		t.Errorf("Collect() returned error for non-conflicting posts: %v", err)
	}

	if len(p.Conflicts()) != 0 {
		t.Errorf("Conflicts() = %d, want 0", len(p.Conflicts()))
	}
}

func TestSlugConflictsPlugin_PostPostConflict(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add two posts with the same slug
	post1 := models.NewPost("posts/first.md")
	post1.Slug = "same-slug"
	post1.Published = true

	post2 := models.NewPost("posts/second.md")
	post2.Slug = "same-slug"
	post2.Published = true

	m.AddPost(post1)
	m.AddPost(post2)

	err := p.Collect(m)
	if err == nil {
		t.Fatal("Collect() should return error for conflicting post slugs")
	}

	if len(p.Conflicts()) != 1 {
		t.Fatalf("Conflicts() = %d, want 1", len(p.Conflicts()))
	}

	conflict := p.Conflicts()[0]
	if conflict.Slug != "same-slug" {
		t.Errorf("conflict.Slug = %q, want %q", conflict.Slug, "same-slug")
	}
	if conflict.ConflictType != "post-post" {
		t.Errorf("conflict.ConflictType = %q, want %q", conflict.ConflictType, "post-post")
	}
	if len(conflict.Sources) != 2 {
		t.Errorf("len(conflict.Sources) = %d, want 2", len(conflict.Sources))
	}

	// Check error message contains useful information
	errMsg := err.Error()
	if !strings.Contains(errMsg, "same-slug") {
		t.Errorf("error message should contain slug: %s", errMsg)
	}
	if !strings.Contains(errMsg, "post-post") {
		t.Errorf("error message should contain conflict type: %s", errMsg)
	}
}

func TestSlugConflictsPlugin_TreatsSlugCaseInsensitively(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	first := models.NewPost("posts/first.md")
	first.Slug = "Same-Slug"
	first.Published = true
	second := models.NewPost("posts/second.md")
	second.Slug = "same-slug"
	second.Published = true
	m.AddPost(first)
	m.AddPost(second)

	if err := p.Collect(m); err == nil {
		t.Fatal("Collect() error = nil, want case-insensitive slug conflict")
	}
	if len(p.Conflicts()) != 1 {
		t.Fatalf("Conflicts() = %d, want 1", len(p.Conflicts()))
	}
}

func TestSlugConflictsPlugin_PostFeedConflict(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add a post with slug "blog"
	post := models.NewPost("posts/blog-post.md")
	post.Slug = "blog"
	post.Published = true
	m.AddPost(post)

	// Configure a feed with the same slug
	feedConfigs := []models.FeedConfig{
		{
			Slug:  "blog",
			Title: "Blog Feed",
			Formats: models.FeedFormats{
				HTML: true,
				RSS:  true,
			},
		},
	}
	m.Cache().Set("feed_configs", feedConfigs)

	err := p.Collect(m)
	if err == nil {
		t.Fatal("Collect() should return error for post-feed slug conflict")
	}

	// Should have 1 post-feed conflict
	conflicts := p.Conflicts()
	var postFeedConflicts []SlugConflict
	for _, c := range conflicts {
		if c.ConflictType == "post-feed" {
			postFeedConflicts = append(postFeedConflicts, c)
		}
	}

	if len(postFeedConflicts) != 1 {
		t.Fatalf("post-feed conflicts = %d, want 1", len(postFeedConflicts))
	}

	conflict := postFeedConflicts[0]
	if conflict.Slug != "blog" {
		t.Errorf("conflict.Slug = %q, want %q", conflict.Slug, "blog")
	}

	// Check error message
	errMsg := err.Error()
	if !strings.Contains(errMsg, "post-feed") {
		t.Errorf("error message should contain conflict type: %s", errMsg)
	}
}

func TestSlugConflictsPlugin_HomepageConflict(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add a post with empty slug (homepage)
	post := models.NewPost("posts/blog/index.md")
	post.Slug = "" // Empty slug = homepage
	post.Published = true
	m.AddPost(post)

	// Configure homepage feed (empty slug)
	feedConfigs := []models.FeedConfig{
		{
			Slug:  "", // Empty slug = homepage
			Title: "Home",
			Formats: models.FeedFormats{
				HTML: true,
			},
		},
	}
	m.Cache().Set("feed_configs", feedConfigs)

	err := p.Collect(m)
	if err == nil {
		t.Fatal("Collect() should return error for homepage slug conflict")
	}

	// Check that empty slug is properly displayed
	errMsg := err.Error()
	if !strings.Contains(errMsg, "homepage") || !strings.Contains(errMsg, "empty slug") {
		t.Errorf("error message should indicate homepage conflict: %s", errMsg)
	}
}

func TestSlugConflictsPlugin_SkipsNonHTMLFeeds(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add a post with slug "data"
	post := models.NewPost("posts/data.md")
	post.Slug = "data"
	post.Published = true
	m.AddPost(post)

	// Configure a feed with the same slug but only RSS (no HTML)
	feedConfigs := []models.FeedConfig{
		{
			Slug:  "data",
			Title: "Data Feed",
			Formats: models.FeedFormats{
				HTML: false,
				RSS:  true,
				JSON: true,
			},
		},
	}
	m.Cache().Set("feed_configs", feedConfigs)

	err := p.Collect(m)
	if err != nil {
		t.Errorf("Collect() should not error when feed doesn't generate HTML: %v", err)
	}
}

func TestSlugConflictsPlugin_SkipsDraftPosts(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add a published post
	post1 := models.NewPost("posts/first.md")
	post1.Slug = "same-slug"
	post1.Published = true

	// Add a draft post with same slug (should be ignored)
	post2 := models.NewPost("posts/second.md")
	post2.Slug = "same-slug"
	post2.Draft = true

	m.AddPost(post1)
	m.AddPost(post2)

	err := p.Collect(m)
	if err != nil {
		t.Errorf("Collect() should not error when only one non-draft post has slug: %v", err)
	}
}

func TestSlugConflictsPlugin_SkipsSkippedPosts(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add a published post
	post1 := models.NewPost("posts/first.md")
	post1.Slug = "same-slug"
	post1.Published = true

	// Add a skipped post with same slug (should be ignored)
	post2 := models.NewPost("posts/second.md")
	post2.Slug = "same-slug"
	post2.Skip = true

	m.AddPost(post1)
	m.AddPost(post2)

	err := p.Collect(m)
	if err != nil {
		t.Errorf("Collect() should not error when only one non-skipped post has slug: %v", err)
	}
}

func TestSlugConflictsPlugin_Disabled(t *testing.T) {
	p := NewSlugConflictsPlugin()
	p.SetEnabled(false)

	m := lifecycle.NewManager()

	// Add conflicting posts
	post1 := models.NewPost("posts/first.md")
	post1.Slug = "same-slug"
	post1.Published = true

	post2 := models.NewPost("posts/second.md")
	post2.Slug = "same-slug"
	post2.Published = true

	m.AddPost(post1)
	m.AddPost(post2)

	err := p.Collect(m)
	if err != nil {
		t.Errorf("Collect() should not error when plugin is disabled: %v", err)
	}
}

func TestSlugConflictsPlugin_MultipleConflicts(t *testing.T) {
	p := NewSlugConflictsPlugin()
	m := lifecycle.NewManager()

	// Add posts with multiple slug conflicts
	post1 := models.NewPost("posts/a.md")
	post1.Slug = "conflict-a"
	post1.Published = true

	post2 := models.NewPost("posts/b.md")
	post2.Slug = "conflict-a"
	post2.Published = true

	post3 := models.NewPost("posts/c.md")
	post3.Slug = "conflict-b"
	post3.Published = true

	post4 := models.NewPost("posts/d.md")
	post4.Slug = "conflict-b"
	post4.Published = true

	m.AddPost(post1)
	m.AddPost(post2)
	m.AddPost(post3)
	m.AddPost(post4)

	err := p.Collect(m)
	if err == nil {
		t.Fatal("Collect() should return error for multiple conflicts")
	}

	if len(p.Conflicts()) != 2 {
		t.Errorf("Conflicts() = %d, want 2", len(p.Conflicts()))
	}

	// Verify both conflicts are detected
	slugs := make(map[string]bool)
	for _, c := range p.Conflicts() {
		slugs[c.Slug] = true
	}

	if !slugs["conflict-a"] || !slugs["conflict-b"] {
		t.Error("expected both conflict-a and conflict-b to be detected")
	}
}
