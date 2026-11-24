package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestPrevNextPlugin_Name(t *testing.T) {
	p := NewPrevNextPlugin()
	if got := p.Name(); got != "prevnext" {
		t.Errorf("Name() = %q, want %q", got, "prevnext")
	}
}

func TestPrevNextPlugin_Priority(t *testing.T) {
	p := NewPrevNextPlugin()

	// Should run late in collect stage (after feeds)
	if got := p.Priority(lifecycle.StageCollect); got != lifecycle.PriorityLate {
		t.Errorf("Priority(StageCollect) = %d, want %d", got, lifecycle.PriorityLate)
	}

	// Default for other stages
	if got := p.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestPrevNextPlugin_Collect_Disabled(t *testing.T) {
	m := lifecycle.NewManager()

	// Create posts
	posts := []*models.Post{
		{Slug: "post-1"},
		{Slug: "post-2"},
	}
	m.SetPosts(posts)

	// Disable plugin
	m.Config().Extra = map[string]interface{}{
		"prevnext": models.PrevNextConfig{
			Enabled: false,
		},
	}

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Errorf("Collect() error = %v", err)
	}

	// Posts should not have prev/next set
	for _, post := range m.Posts() {
		if post.Prev != nil || post.Next != nil {
			t.Errorf("Post %s has prev/next set despite plugin being disabled", post.Slug)
		}
	}
}

func TestPrevNextPlugin_Collect_FirstFeedStrategy(t *testing.T) {
	m := lifecycle.NewManager()

	// Create posts
	post1 := &models.Post{Slug: "post-1"}
	post2 := &models.Post{Slug: "post-2"}
	post3 := &models.Post{Slug: "post-3"}
	posts := []*models.Post{post1, post2, post3}
	m.SetPosts(posts)

	// Create a feed containing the posts in order
	feeds := []*lifecycle.Feed{
		{
			Name:  "blog",
			Title: "Blog",
			Posts: []*models.Post{post1, post2, post3},
		},
	}
	m.SetFeeds(feeds)

	// Default config (first_feed strategy)
	m.Config().Extra = map[string]interface{}{}

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Verify post1
	if post1.Prev != nil {
		t.Errorf("post1.Prev = %v, want nil", post1.Prev)
	}
	if post1.Next != post2 {
		t.Errorf("post1.Next = %v, want post2", post1.Next)
	}
	if post1.PrevNextFeed != "blog" {
		t.Errorf("post1.PrevNextFeed = %q, want %q", post1.PrevNextFeed, "blog")
	}
	if post1.PrevNextContext == nil {
		t.Fatal("post1.PrevNextContext is nil")
	}
	if post1.PrevNextContext.Position != 1 {
		t.Errorf("post1.PrevNextContext.Position = %d, want 1", post1.PrevNextContext.Position)
	}
	if post1.PrevNextContext.Total != 3 {
		t.Errorf("post1.PrevNextContext.Total = %d, want 3", post1.PrevNextContext.Total)
	}

	// Verify post2
	if post2.Prev != post1 {
		t.Errorf("post2.Prev = %v, want post1", post2.Prev)
	}
	if post2.Next != post3 {
		t.Errorf("post2.Next = %v, want post3", post2.Next)
	}
	if post2.PrevNextContext.Position != 2 {
		t.Errorf("post2.PrevNextContext.Position = %d, want 2", post2.PrevNextContext.Position)
	}

	// Verify post3
	if post3.Prev != post2 {
		t.Errorf("post3.Prev = %v, want post2", post3.Prev)
	}
	if post3.Next != nil {
		t.Errorf("post3.Next = %v, want nil", post3.Next)
	}
	if post3.PrevNextContext.Position != 3 {
		t.Errorf("post3.PrevNextContext.Position = %d, want 3", post3.PrevNextContext.Position)
	}
}

func TestPrevNextPlugin_Collect_ExplicitFeedStrategy(t *testing.T) {
	m := lifecycle.NewManager()

	// Create posts
	post1 := &models.Post{Slug: "post-1"}
	post2 := &models.Post{Slug: "post-2"}
	post3 := &models.Post{Slug: "post-3"}
	posts := []*models.Post{post1, post2, post3}
	m.SetPosts(posts)

	// Create two feeds with different post orders
	feeds := []*lifecycle.Feed{
		{
			Name:  "blog",
			Title: "Blog",
			Posts: []*models.Post{post1, post2, post3},
		},
		{
			Name:  "featured",
			Title: "Featured",
			Posts: []*models.Post{post3, post1}, // Different order, post2 missing
		},
	}
	m.SetFeeds(feeds)

	// Use explicit_feed strategy with "featured" feed
	m.Config().Extra = map[string]interface{}{
		"prevnext": models.PrevNextConfig{
			Enabled:     true,
			Strategy:    models.StrategyExplicitFeed,
			DefaultFeed: "featured",
		},
	}

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// post1 is second in featured feed
	if post1.Prev != post3 {
		t.Errorf("post1.Prev = %v, want post3", post1.Prev)
	}
	if post1.Next != nil {
		t.Errorf("post1.Next = %v, want nil", post1.Next)
	}
	if post1.PrevNextFeed != "featured" {
		t.Errorf("post1.PrevNextFeed = %q, want %q", post1.PrevNextFeed, "featured")
	}

	// post2 is not in featured feed, should have no prev/next
	if post2.Prev != nil || post2.Next != nil {
		t.Errorf("post2 should have no prev/next (not in featured feed)")
	}

	// post3 is first in featured feed
	if post3.Prev != nil {
		t.Errorf("post3.Prev = %v, want nil", post3.Prev)
	}
	if post3.Next != post1 {
		t.Errorf("post3.Next = %v, want post1", post3.Next)
	}
}

func TestPrevNextPlugin_Collect_SeriesStrategy(t *testing.T) {
	m := lifecycle.NewManager()

	// Create posts - post2 has a series frontmatter
	post1 := &models.Post{Slug: "post-1"}
	post2 := &models.Post{
		Slug:  "post-2",
		Extra: map[string]interface{}{"series": "tutorial"},
	}
	post3 := &models.Post{Slug: "post-3"}
	posts := []*models.Post{post1, post2, post3}
	m.SetPosts(posts)

	// Create feeds - tutorial feed has different order
	feeds := []*lifecycle.Feed{
		{
			Name:  "blog",
			Title: "Blog",
			Posts: []*models.Post{post1, post2, post3},
		},
		{
			Name:  "series/tutorial",
			Title: "Tutorial Series",
			Posts: []*models.Post{post3, post2}, // post2 is last in tutorial
		},
	}
	m.SetFeeds(feeds)

	// Use series strategy
	m.Config().Extra = map[string]interface{}{
		"prevnext": models.PrevNextConfig{
			Enabled:  true,
			Strategy: models.StrategySeries,
		},
	}

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// post1 uses first_feed (blog) since no series
	if post1.PrevNextFeed != "blog" {
		t.Errorf("post1.PrevNextFeed = %q, want %q (first_feed fallback)", post1.PrevNextFeed, "blog")
	}

	// post2 uses tutorial series
	if post2.PrevNextFeed != "series/tutorial" {
		t.Errorf("post2.PrevNextFeed = %q, want %q (from series)", post2.PrevNextFeed, "series/tutorial")
	}
	if post2.Prev != post3 {
		t.Errorf("post2.Prev = %v, want post3 (in tutorial)", post2.Prev)
	}
	if post2.Next != nil {
		t.Errorf("post2.Next = %v, want nil (last in tutorial)", post2.Next)
	}

	// post3 uses first_feed (blog) since no series
	if post3.PrevNextFeed != "blog" {
		t.Errorf("post3.PrevNextFeed = %q, want %q (first_feed fallback)", post3.PrevNextFeed, "blog")
	}
}

func TestPrevNextPlugin_Collect_FrontmatterStrategy(t *testing.T) {
	m := lifecycle.NewManager()

	// Create posts - post2 has prevnext_feed frontmatter
	post1 := &models.Post{Slug: "post-1"}
	post2 := &models.Post{
		Slug:  "post-2",
		Extra: map[string]interface{}{"prevnext_feed": "featured"},
	}
	post3 := &models.Post{Slug: "post-3"}
	posts := []*models.Post{post1, post2, post3}
	m.SetPosts(posts)

	// Create feeds
	feeds := []*lifecycle.Feed{
		{
			Name:  "blog",
			Title: "Blog",
			Posts: []*models.Post{post1, post2, post3},
		},
		{
			Name:  "featured",
			Title: "Featured",
			Posts: []*models.Post{post2}, // Only post2 in featured
		},
	}
	m.SetFeeds(feeds)

	// Use frontmatter strategy
	m.Config().Extra = map[string]interface{}{
		"prevnext": models.PrevNextConfig{
			Enabled:  true,
			Strategy: models.StrategyFrontmatter,
		},
	}

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// post1 uses first_feed (blog) since no prevnext_feed
	if post1.PrevNextFeed != "blog" {
		t.Errorf("post1.PrevNextFeed = %q, want %q", post1.PrevNextFeed, "blog")
	}

	// post2 uses featured feed from frontmatter
	if post2.PrevNextFeed != "featured" {
		t.Errorf("post2.PrevNextFeed = %q, want %q", post2.PrevNextFeed, "featured")
	}
	// post2 is alone in featured, so no prev/next
	if post2.Prev != nil || post2.Next != nil {
		t.Errorf("post2 should have no prev/next (alone in featured)")
	}
}

func TestPrevNextPlugin_Collect_NoFeeds(t *testing.T) {
	m := lifecycle.NewManager()

	// Create posts
	posts := []*models.Post{
		{Slug: "post-1"},
		{Slug: "post-2"},
	}
	m.SetPosts(posts)

	// No feeds
	m.SetFeeds([]*lifecycle.Feed{})

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Posts should have no prev/next
	for _, post := range m.Posts() {
		if post.Prev != nil || post.Next != nil {
			t.Errorf("Post %s has prev/next despite no feeds", post.Slug)
		}
	}
}

func TestPrevNextPlugin_Collect_PostNotInAnyFeed(t *testing.T) {
	m := lifecycle.NewManager()

	// Create posts
	post1 := &models.Post{Slug: "post-1"}
	post2 := &models.Post{Slug: "post-2"}
	orphan := &models.Post{Slug: "orphan"}
	posts := []*models.Post{post1, post2, orphan}
	m.SetPosts(posts)

	// Feed only contains post1 and post2
	feeds := []*lifecycle.Feed{
		{
			Name:  "blog",
			Title: "Blog",
			Posts: []*models.Post{post1, post2},
		},
	}
	m.SetFeeds(feeds)

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Orphan should have no prev/next
	if orphan.Prev != nil || orphan.Next != nil {
		t.Errorf("orphan should have no prev/next")
	}
	if orphan.PrevNextContext != nil {
		t.Errorf("orphan.PrevNextContext should be nil")
	}
}

func TestPrevNextPlugin_Collect_SinglePost(t *testing.T) {
	m := lifecycle.NewManager()

	// Single post
	post := &models.Post{Slug: "only-post"}
	m.SetPosts([]*models.Post{post})

	feeds := []*lifecycle.Feed{
		{
			Name:  "blog",
			Title: "Blog",
			Posts: []*models.Post{post},
		},
	}
	m.SetFeeds(feeds)

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Single post has no prev/next
	if post.Prev != nil {
		t.Errorf("post.Prev = %v, want nil", post.Prev)
	}
	if post.Next != nil {
		t.Errorf("post.Next = %v, want nil", post.Next)
	}

	// But it should have context
	if post.PrevNextContext == nil {
		t.Fatal("post.PrevNextContext is nil")
	}
	if post.PrevNextContext.Position != 1 {
		t.Errorf("post.PrevNextContext.Position = %d, want 1", post.PrevNextContext.Position)
	}
	if post.PrevNextContext.Total != 1 {
		t.Errorf("post.PrevNextContext.Total = %d, want 1", post.PrevNextContext.Total)
	}
	if !post.PrevNextContext.IsFirst() {
		t.Error("post.PrevNextContext.IsFirst() should be true")
	}
	if !post.PrevNextContext.IsLast() {
		t.Error("post.PrevNextContext.IsLast() should be true")
	}
}

func TestPrevNextPlugin_Collect_ConfigFromMap(t *testing.T) {
	m := lifecycle.NewManager()

	post1 := &models.Post{Slug: "post-1"}
	post2 := &models.Post{Slug: "post-2"}
	m.SetPosts([]*models.Post{post1, post2})

	feeds := []*lifecycle.Feed{
		{Name: "blog", Title: "Blog", Posts: []*models.Post{post1, post2}},
		{Name: "alt", Title: "Alt", Posts: []*models.Post{post2, post1}},
	}
	m.SetFeeds(feeds)

	// Config as map (as it might come from TOML parsing)
	m.Config().Extra = map[string]interface{}{
		"prevnext": map[string]interface{}{
			"enabled":      true,
			"strategy":     "explicit_feed",
			"default_feed": "alt",
		},
	}

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Should use alt feed
	if post1.PrevNextFeed != "alt" {
		t.Errorf("post1.PrevNextFeed = %q, want %q", post1.PrevNextFeed, "alt")
	}
	// In alt feed, post1 is second (after post2)
	if post1.Prev != post2 {
		t.Errorf("post1.Prev = %v, want post2", post1.Prev)
	}
}

func TestPrevNextPlugin_Collect_MultipleFeeds(t *testing.T) {
	m := lifecycle.NewManager()

	// Posts that appear in multiple feeds
	post1 := &models.Post{Slug: "post-1"}
	post2 := &models.Post{Slug: "post-2"}
	post3 := &models.Post{Slug: "post-3"}
	m.SetPosts([]*models.Post{post1, post2, post3})

	// post2 appears in both feeds but in different positions
	feeds := []*lifecycle.Feed{
		{
			Name:  "blog",
			Title: "Blog",
			Posts: []*models.Post{post1, post2, post3}, // post2 in middle
		},
		{
			Name:  "featured",
			Title: "Featured",
			Posts: []*models.Post{post2, post3}, // post2 is first
		},
	}
	m.SetFeeds(feeds)

	// Default (first_feed) strategy
	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// post2 should use first feed it appears in (blog)
	if post2.PrevNextFeed != "blog" {
		t.Errorf("post2.PrevNextFeed = %q, want %q (first feed)", post2.PrevNextFeed, "blog")
	}
	if post2.Prev != post1 {
		t.Errorf("post2.Prev = %v, want post1 (from blog)", post2.Prev)
	}
}

func TestPrevNextPlugin_Collect_SeriesFallback(t *testing.T) {
	m := lifecycle.NewManager()

	// Post with invalid series (feed doesn't exist)
	post := &models.Post{
		Slug:  "post-1",
		Extra: map[string]interface{}{"series": "nonexistent"},
	}
	m.SetPosts([]*models.Post{post})

	feeds := []*lifecycle.Feed{
		{Name: "blog", Title: "Blog", Posts: []*models.Post{post}},
	}
	m.SetFeeds(feeds)

	m.Config().Extra = map[string]interface{}{
		"prevnext": models.PrevNextConfig{
			Enabled:  true,
			Strategy: models.StrategySeries,
		},
	}

	p := NewPrevNextPlugin()
	if err := p.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Should fall back to first_feed (blog)
	if post.PrevNextFeed != "blog" {
		t.Errorf("post.PrevNextFeed = %q, want %q (fallback)", post.PrevNextFeed, "blog")
	}
}

func TestPrevNextContext_Helpers(t *testing.T) {
	prev := &models.Post{Slug: "prev"}
	next := &models.Post{Slug: "next"}

	tests := []struct {
		name    string
		ctx     *models.PrevNextContext
		hasPrev bool
		hasNext bool
		isFirst bool
		isLast  bool
	}{
		{
			name:    "middle post",
			ctx:     &models.PrevNextContext{Position: 2, Total: 3, Prev: prev, Next: next},
			hasPrev: true,
			hasNext: true,
			isFirst: false,
			isLast:  false,
		},
		{
			name:    "first post",
			ctx:     &models.PrevNextContext{Position: 1, Total: 3, Prev: nil, Next: next},
			hasPrev: false,
			hasNext: true,
			isFirst: true,
			isLast:  false,
		},
		{
			name:    "last post",
			ctx:     &models.PrevNextContext{Position: 3, Total: 3, Prev: prev, Next: nil},
			hasPrev: true,
			hasNext: false,
			isFirst: false,
			isLast:  true,
		},
		{
			name:    "only post",
			ctx:     &models.PrevNextContext{Position: 1, Total: 1, Prev: nil, Next: nil},
			hasPrev: false,
			hasNext: false,
			isFirst: true,
			isLast:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ctx.HasPrev(); got != tt.hasPrev {
				t.Errorf("HasPrev() = %v, want %v", got, tt.hasPrev)
			}
			if got := tt.ctx.HasNext(); got != tt.hasNext {
				t.Errorf("HasNext() = %v, want %v", got, tt.hasNext)
			}
			if got := tt.ctx.IsFirst(); got != tt.isFirst {
				t.Errorf("IsFirst() = %v, want %v", got, tt.isFirst)
			}
			if got := tt.ctx.IsLast(); got != tt.isLast {
				t.Errorf("IsLast() = %v, want %v", got, tt.isLast)
			}
		})
	}
}

func TestPrevNextStrategy_IsValid(t *testing.T) {
	tests := []struct {
		strategy models.PrevNextStrategy
		valid    bool
	}{
		{models.StrategyFirstFeed, true},
		{models.StrategyExplicitFeed, true},
		{models.StrategySeries, true},
		{models.StrategyFrontmatter, true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.strategy), func(t *testing.T) {
			if got := tt.strategy.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestPrevNextConfig_ApplyDefaults(t *testing.T) {
	cfg := models.PrevNextConfig{}
	cfg.ApplyDefaults()

	if cfg.Strategy != models.StrategyFirstFeed {
		t.Errorf("Strategy = %q, want %q", cfg.Strategy, models.StrategyFirstFeed)
	}
}

// Verify interface implementations
func TestPrevNextPlugin_Interfaces(_ *testing.T) {
	var _ lifecycle.Plugin = (*PrevNextPlugin)(nil)
	var _ lifecycle.CollectPlugin = (*PrevNextPlugin)(nil)
	var _ lifecycle.PriorityPlugin = (*PrevNextPlugin)(nil)
}
