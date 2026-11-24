package plugins

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// =============================================================================
// SeriesPlugin Name and Priority
// =============================================================================

func TestSeriesPlugin_Name(t *testing.T) {
	plugin := NewSeriesPlugin()
	if plugin.Name() != "series" {
		t.Errorf("Name() = %q, want %q", plugin.Name(), "series")
	}
}

func TestSeriesPlugin_Priority(t *testing.T) {
	plugin := NewSeriesPlugin()

	if got := plugin.Priority(lifecycle.StageCollect); got != lifecycle.PriorityEarly {
		t.Errorf("Priority(StageCollect) = %d, want %d (PriorityEarly)", got, lifecycle.PriorityEarly)
	}

	if got := plugin.Priority(lifecycle.StageWrite); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageWrite) = %d, want %d (PriorityDefault)", got, lifecycle.PriorityDefault)
	}
}

// =============================================================================
// SeriesPlugin Collect Tests
// =============================================================================

func TestSeriesPlugin_NoPosts(t *testing.T) {
	m := lifecycle.NewManager()

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// No feeds should be created
	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 0 {
		t.Errorf("expected 0 feed configs, got %d", len(feedConfigs))
	}
}

func TestSeriesPlugin_NoSeriesPosts(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2"), Date: &date},
	})

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 0 {
		t.Errorf("expected 0 feed configs, got %d", len(feedConfigs))
	}
}

func TestSeriesPlugin_BasicSeries(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	post1 := &models.Post{
		Path: "part1.md", Slug: "part1", Title: strPtr("Part 1"), Date: &date1,
		Published: true,
		Extra:     map[string]interface{}{"series": "building-rest-api"},
	}
	post2 := &models.Post{
		Path: "part2.md", Slug: "part2", Title: strPtr("Part 2"), Date: &date2,
		Published: true,
		Extra:     map[string]interface{}{"series": "building-rest-api"},
	}
	post3 := &models.Post{
		Path: "part3.md", Slug: "part3", Title: strPtr("Part 3"), Date: &date3,
		Published: true,
		Extra:     map[string]interface{}{"series": "building-rest-api"},
	}

	m.SetPosts([]*models.Post{post3, post1, post2}) // out of order

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	fc := feedConfigs[0]

	// Check slug
	if fc.Slug != "series/building-rest-api" {
		t.Errorf("feed slug = %q, want %q", fc.Slug, "series/building-rest-api")
	}

	// Check type
	if fc.Type != models.FeedTypeSeries {
		t.Errorf("feed type = %q, want %q", fc.Type, models.FeedTypeSeries)
	}

	// Check title (title-cased from series name)
	if fc.Title != "Building Rest Api" {
		t.Errorf("feed title = %q, want %q", fc.Title, "Building Rest Api")
	}

	// Check sidebar
	if fc.Sidebar == nil || !*fc.Sidebar {
		t.Error("feed sidebar should be true")
	}

	// Check post order (should be date ascending: post1, post2, post3)
	if len(fc.Posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(fc.Posts))
	}
	if fc.Posts[0].Slug != "part1" {
		t.Errorf("first post = %q, want %q", fc.Posts[0].Slug, "part1")
	}
	if fc.Posts[1].Slug != "part2" {
		t.Errorf("second post = %q, want %q", fc.Posts[1].Slug, "part2")
	}
	if fc.Posts[2].Slug != "part3" {
		t.Errorf("third post = %q, want %q", fc.Posts[2].Slug, "part3")
	}

	// Check prev/next navigation
	if post1.Prev != nil {
		t.Errorf("post1.Prev should be nil, got %v", post1.Prev)
	}
	if post1.Next != post2 {
		t.Error("post1.Next should be post2")
	}
	if post2.Prev != post1 {
		t.Error("post2.Prev should be post1")
	}
	if post2.Next != post3 {
		t.Error("post2.Next should be post3")
	}
	if post3.Prev != post2 {
		t.Error("post3.Prev should be post2")
	}
	if post3.Next != nil {
		t.Errorf("post3.Next should be nil, got %v", post3.Next)
	}

	// Check PrevNextFeed
	if post1.PrevNextFeed != "series/building-rest-api" {
		t.Errorf("post1.PrevNextFeed = %q, want %q", post1.PrevNextFeed, "series/building-rest-api")
	}

	// Check PrevNextContext
	if post1.PrevNextContext == nil {
		t.Fatal("post1.PrevNextContext should not be nil")
	}
	if post1.PrevNextContext.Position != 1 {
		t.Errorf("post1 position = %d, want 1", post1.PrevNextContext.Position)
	}
	if post1.PrevNextContext.Total != 3 {
		t.Errorf("post1 total = %d, want 3", post1.PrevNextContext.Total)
	}
	if post2.PrevNextContext.Position != 2 {
		t.Errorf("post2 position = %d, want 2", post2.PrevNextContext.Position)
	}
	if post3.PrevNextContext.Position != 3 {
		t.Errorf("post3 position = %d, want 3", post3.PrevNextContext.Position)
	}

	// Check series metadata in Extra
	if post1.Extra["series_slug"] != "series/building-rest-api" {
		t.Errorf("post1.Extra[series_slug] = %v, want %q", post1.Extra["series_slug"], "series/building-rest-api")
	}
	if post1.Extra["series_total"] != 3 {
		t.Errorf("post1.Extra[series_total] = %v, want 3", post1.Extra["series_total"])
	}
}

func TestSeriesPlugin_ExplicitOrder(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	// Posts with explicit order, reversed from date order
	post1 := &models.Post{
		Path: "setup.md", Slug: "setup", Title: strPtr("Setup"), Date: &date3,
		Published: true,
		Extra:     map[string]interface{}{"series": "tutorial", "series_order": 1},
	}
	post2 := &models.Post{
		Path: "basics.md", Slug: "basics", Title: strPtr("Basics"), Date: &date1,
		Published: true,
		Extra:     map[string]interface{}{"series": "tutorial", "series_order": 2},
	}
	post3 := &models.Post{
		Path: "advanced.md", Slug: "advanced", Title: strPtr("Advanced"), Date: &date2,
		Published: true,
		Extra:     map[string]interface{}{"series": "tutorial", "series_order": 3},
	}

	m.SetPosts([]*models.Post{post3, post2, post1}) // random order

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	fc := feedConfigs[0]

	// Should be ordered by series_order, not date
	if len(fc.Posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(fc.Posts))
	}
	if fc.Posts[0].Slug != "setup" {
		t.Errorf("first post = %q, want %q", fc.Posts[0].Slug, "setup")
	}
	if fc.Posts[1].Slug != "basics" {
		t.Errorf("second post = %q, want %q", fc.Posts[1].Slug, "basics")
	}
	if fc.Posts[2].Slug != "advanced" {
		t.Errorf("third post = %q, want %q", fc.Posts[2].Slug, "advanced")
	}
}

func TestSeriesPlugin_MixedOrder(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	// Some posts have explicit order, others don't
	post1 := &models.Post{
		Path: "intro.md", Slug: "intro", Title: strPtr("Intro"), Date: &date3,
		Published: true,
		Extra:     map[string]interface{}{"series": "guide", "series_order": 1},
	}
	post2 := &models.Post{
		Path: "setup.md", Slug: "setup", Title: strPtr("Setup"), Date: &date1,
		Published: true,
		Extra:     map[string]interface{}{"series": "guide", "series_order": 2},
	}
	// No series_order - should come after ordered posts, sorted by date
	post3 := &models.Post{
		Path: "extra.md", Slug: "extra", Title: strPtr("Extra"), Date: &date2,
		Published: true,
		Extra:     map[string]interface{}{"series": "guide"},
	}

	m.SetPosts([]*models.Post{post3, post1, post2})

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	// Order: post1 (order=1), post2 (order=2), post3 (no order, by date)
	fc := feedConfigs[0]
	if len(fc.Posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(fc.Posts))
	}
	if fc.Posts[0].Slug != "intro" {
		t.Errorf("first post = %q, want %q", fc.Posts[0].Slug, "intro")
	}
	if fc.Posts[1].Slug != "setup" {
		t.Errorf("second post = %q, want %q", fc.Posts[1].Slug, "setup")
	}
	if fc.Posts[2].Slug != "extra" {
		t.Errorf("third post = %q, want %q", fc.Posts[2].Slug, "extra")
	}
}

func TestSeriesPlugin_MultipleSeries(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "go1.md", Slug: "go1", Title: strPtr("Go Part 1"), Date: &date1, Published: true,
			Extra: map[string]interface{}{"series": "learn-go"}},
		{Path: "go2.md", Slug: "go2", Title: strPtr("Go Part 2"), Date: &date2, Published: true,
			Extra: map[string]interface{}{"series": "learn-go"}},
		{Path: "py1.md", Slug: "py1", Title: strPtr("Python Part 1"), Date: &date2, Published: true,
			Extra: map[string]interface{}{"series": "learn-python"}},
		{Path: "py2.md", Slug: "py2", Title: strPtr("Python Part 2"), Date: &date1, Published: true,
			Extra: map[string]interface{}{"series": "learn-python"}},
	})

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 2 {
		t.Fatalf("expected 2 feed configs, got %d", len(feedConfigs))
	}

	// Build map for lookup
	configBySlug := make(map[string]models.FeedConfig)
	for _, fc := range feedConfigs {
		configBySlug[fc.Slug] = fc
	}

	goFeed, ok := configBySlug["series/learn-go"]
	if !ok {
		t.Fatal("series/learn-go feed not found")
	}
	if len(goFeed.Posts) != 2 {
		t.Fatalf("learn-go: expected 2 posts, got %d", len(goFeed.Posts))
	}
	// Date ascending: go1 (Jan) before go2 (Feb)
	if goFeed.Posts[0].Slug != "go1" {
		t.Errorf("learn-go first post = %q, want %q", goFeed.Posts[0].Slug, "go1")
	}

	pyFeed, ok := configBySlug["series/learn-python"]
	if !ok {
		t.Fatal("series/learn-python feed not found")
	}
	if len(pyFeed.Posts) != 2 {
		t.Fatalf("learn-python: expected 2 posts, got %d", len(pyFeed.Posts))
	}
	// Date ascending: py2 (Jan) before py1 (Feb)
	if pyFeed.Posts[0].Slug != "py2" {
		t.Errorf("learn-python first post = %q, want %q", pyFeed.Posts[0].Slug, "py2")
	}
}

func TestSeriesPlugin_SinglePost(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	post := &models.Post{
		Path: "solo.md", Slug: "solo", Title: strPtr("Solo Post"), Date: &date,
		Published: true,
		Extra:     map[string]interface{}{"series": "one-part"},
	}
	m.SetPosts([]*models.Post{post})

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	// Single post: no prev/next
	if post.Prev != nil {
		t.Error("single post should have nil Prev")
	}
	if post.Next != nil {
		t.Error("single post should have nil Next")
	}
	if post.PrevNextContext == nil {
		t.Fatal("PrevNextContext should not be nil")
	}
	if post.PrevNextContext.Position != 1 {
		t.Errorf("position = %d, want 1", post.PrevNextContext.Position)
	}
	if post.PrevNextContext.Total != 1 {
		t.Errorf("total = %d, want 1", post.PrevNextContext.Total)
	}
}

// =============================================================================
// Configuration Tests
// =============================================================================

func TestSeriesPlugin_ConfigOverride(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{
		{Path: "p1.md", Slug: "p1", Title: strPtr("P1"), Date: &date,
			Published: true,
			Extra:     map[string]interface{}{"series": "my-series"}},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"series": map[string]interface{}{
			"overrides": map[string]interface{}{
				"my-series": map[string]interface{}{
					"title":       "Custom Title",
					"description": "Custom description",
				},
			},
		},
	}
	m.SetConfig(config)

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	fc := feedConfigs[0]
	if fc.Title != "Custom Title" {
		t.Errorf("title = %q, want %q", fc.Title, "Custom Title")
	}
	if fc.Description != "Custom description" {
		t.Errorf("description = %q, want %q", fc.Description, "Custom description")
	}
}

func TestSeriesPlugin_CustomSlugPrefix(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{
		{Path: "p1.md", Slug: "p1", Title: strPtr("P1"), Date: &date,
			Published: true,
			Extra:     map[string]interface{}{"series": "my-series"}},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"series": map[string]interface{}{
			"slug_prefix": "courses",
		},
	}
	m.SetConfig(config)

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	if feedConfigs[0].Slug != "courses/my-series" {
		t.Errorf("slug = %q, want %q", feedConfigs[0].Slug, "courses/my-series")
	}
}

func TestSeriesPlugin_DefaultsConfig(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{
		{Path: "p1.md", Slug: "p1", Title: strPtr("P1"), Date: &date,
			Published: true,
			Extra:     map[string]interface{}{"series": "test"}},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"series": map[string]interface{}{
			"defaults": map[string]interface{}{
				"items_per_page": 5,
				"sidebar":        false,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	fc := feedConfigs[0]
	if fc.ItemsPerPage != 5 {
		t.Errorf("items_per_page = %d, want 5", fc.ItemsPerPage)
	}
	if fc.Sidebar != nil && *fc.Sidebar {
		t.Error("sidebar should be false")
	}
}

func TestSeriesPlugin_PreservesExistingFeeds(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{
		{Path: "p1.md", Slug: "p1", Title: strPtr("P1"), Date: &date,
			Published: true,
			Extra:     map[string]interface{}{"series": "test"}},
		{Path: "p2.md", Slug: "p2", Title: strPtr("P2"), Date: &date, Tags: []string{"go"}},
	})

	// Pre-existing feed config
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{Slug: "blog", Title: "Blog", Filter: `"go" in tags`},
		},
	}
	m.SetConfig(config)

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 2 {
		t.Fatalf("expected 2 feed configs (1 existing + 1 series), got %d", len(feedConfigs))
	}

	// First should be the existing feed
	if feedConfigs[0].Slug != "blog" {
		t.Errorf("first feed slug = %q, want %q", feedConfigs[0].Slug, "blog")
	}

	// Second should be the series feed
	if feedConfigs[1].Slug != "series/test" {
		t.Errorf("second feed slug = %q, want %q", feedConfigs[1].Slug, "series/test")
	}
}

// =============================================================================
// Series Normalization Tests
// =============================================================================

func TestSeriesPlugin_NameNormalization(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// Different posts use different casings of the same series name
	m.SetPosts([]*models.Post{
		{Path: "p1.md", Slug: "p1", Title: strPtr("P1"), Date: &date,
			Published: true,
			Extra:     map[string]interface{}{"series": "My Series"}},
		{Path: "p2.md", Slug: "p2", Title: strPtr("P2"), Date: &date,
			Published: true,
			Extra:     map[string]interface{}{"series": "my-series"}},
	})

	plugin := NewSeriesPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	// Both should be grouped into the same series (slugified)
	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config (both names slugify the same), got %d", len(feedConfigs))
	}
	if feedConfigs[0].Slug != "series/my-series" {
		t.Errorf("slug = %q, want %q", feedConfigs[0].Slug, "series/my-series")
	}
	if len(feedConfigs[0].Posts) != 2 {
		t.Errorf("expected 2 posts, got %d", len(feedConfigs[0].Posts))
	}
}

// =============================================================================
// Sorting Helper Tests
// =============================================================================

func TestGetSeriesOrder(t *testing.T) {
	tests := []struct {
		name   string
		extra  map[string]interface{}
		want   int
		wantOk bool
	}{
		{"nil extra", nil, 0, false},
		{"no series_order", map[string]interface{}{"series": "test"}, 0, false},
		{"int value", map[string]interface{}{"series_order": 3}, 3, true},
		{"float64 value", map[string]interface{}{"series_order": 2.0}, 2, true},
		{"int64 value", map[string]interface{}{"series_order": int64(5)}, 5, true},
		{"string value", map[string]interface{}{"series_order": "not a number"}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Extra: tt.extra}
			got, ok := getSeriesOrder(post)
			if ok != tt.wantOk {
				t.Errorf("getSeriesOrder() ok = %v, want %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("getSeriesOrder() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTieBreakByDateThenPath(t *testing.T) {
	date1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		a    *models.Post
		b    *models.Post
		want bool
	}{
		{
			"earlier date first",
			&models.Post{Path: "b.md", Date: &date1},
			&models.Post{Path: "a.md", Date: &date2},
			true,
		},
		{
			"later date second",
			&models.Post{Path: "a.md", Date: &date2},
			&models.Post{Path: "b.md", Date: &date1},
			false,
		},
		{
			"same date - path tiebreak",
			&models.Post{Path: "a.md", Date: &date1},
			&models.Post{Path: "b.md", Date: &date1},
			true,
		},
		{
			"no dates - path tiebreak",
			&models.Post{Path: "a.md"},
			&models.Post{Path: "b.md"},
			true,
		},
		{
			"first has date, second nil",
			&models.Post{Path: "b.md", Date: &date1},
			&models.Post{Path: "a.md"},
			true,
		},
		{
			"first nil, second has date",
			&models.Post{Path: "a.md"},
			&models.Post{Path: "b.md", Date: &date1},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tieBreakByDateThenPath(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("tieBreakByDateThenPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// =============================================================================
// Configuration Parsing Tests
// =============================================================================

func TestParseSeriesConfig_Defaults(t *testing.T) {
	config := lifecycle.NewConfig()
	cfg := parseSeriesConfig(config)

	if cfg.SlugPrefix != "series" {
		t.Errorf("SlugPrefix = %q, want %q", cfg.SlugPrefix, "series")
	}
	if !cfg.AutoSidebar {
		t.Error("AutoSidebar should be true by default")
	}
	if cfg.Defaults.ItemsPerPage != 0 {
		t.Errorf("Defaults.ItemsPerPage = %d, want 0", cfg.Defaults.ItemsPerPage)
	}
	if !cfg.Defaults.Sidebar {
		t.Error("Defaults.Sidebar should be true by default")
	}
	if len(cfg.Overrides) != 0 {
		t.Errorf("Overrides should be empty, got %d", len(cfg.Overrides))
	}
}

func TestParseSeriesConfig_Full(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra["series"] = map[string]interface{}{
		"slug_prefix":  "courses",
		"auto_sidebar": false,
		"defaults": map[string]interface{}{
			"items_per_page": 20,
			"sidebar":        false,
			"formats": map[string]interface{}{
				"html": true,
				"rss":  true,
			},
		},
		"overrides": map[string]interface{}{
			"golang-101": map[string]interface{}{
				"title":          "Golang 101",
				"description":    "Learn Go from scratch",
				"items_per_page": 10,
			},
		},
	}

	cfg := parseSeriesConfig(config)

	if cfg.SlugPrefix != "courses" {
		t.Errorf("SlugPrefix = %q, want %q", cfg.SlugPrefix, "courses")
	}
	if cfg.AutoSidebar {
		t.Error("AutoSidebar should be false")
	}
	if cfg.Defaults.ItemsPerPage != 20 {
		t.Errorf("Defaults.ItemsPerPage = %d, want 20", cfg.Defaults.ItemsPerPage)
	}
	if cfg.Defaults.Sidebar {
		t.Error("Defaults.Sidebar should be false")
	}
	if cfg.Defaults.Formats == nil {
		t.Fatal("Defaults.Formats should not be nil")
	}
	if !cfg.Defaults.Formats.HTML {
		t.Error("Defaults.Formats.HTML should be true")
	}
	if !cfg.Defaults.Formats.RSS {
		t.Error("Defaults.Formats.RSS should be true")
	}

	override, ok := cfg.Overrides["golang-101"]
	if !ok {
		t.Fatal("override for golang-101 not found")
	}
	if override.Title != "Golang 101" {
		t.Errorf("override title = %q, want %q", override.Title, "Golang 101")
	}
	if override.Description != "Learn Go from scratch" {
		t.Errorf("override description = %q, want %q", override.Description, "Learn Go from scratch")
	}
	if override.ItemsPerPage == nil {
		t.Fatal("override ItemsPerPage should not be nil")
	}
	if *override.ItemsPerPage != 10 {
		t.Errorf("override ItemsPerPage = %d, want 10", *override.ItemsPerPage)
	}
}

func TestParseSeriesConfig_NilExtra(t *testing.T) {
	config := &lifecycle.Config{}
	cfg := parseSeriesConfig(config)

	// Should return defaults without panic
	if cfg.SlugPrefix != "series" {
		t.Errorf("SlugPrefix = %q, want %q", cfg.SlugPrefix, "series")
	}
}

func TestParseSeriesConfig_InvalidType(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra["series"] = "not a map"

	cfg := parseSeriesConfig(config)
	// Should return defaults
	if cfg.SlugPrefix != "series" {
		t.Errorf("SlugPrefix = %q, want %q", cfg.SlugPrefix, "series")
	}
}

// =============================================================================
// FeedFormats Parsing Tests
// =============================================================================

func TestParseFeedFormatsFromMap(t *testing.T) {
	tests := []struct {
		name string
		raw  interface{}
		want *models.FeedFormats
	}{
		{
			"nil input",
			nil,
			nil,
		},
		{
			"non-map input",
			"invalid",
			nil,
		},
		{
			"all enabled",
			map[string]interface{}{
				"html":        true,
				"simple_html": true,
				"rss":         true,
				"atom":        true,
				"json":        true,
				"markdown":    true,
				"text":        true,
				"sitemap":     true,
			},
			&models.FeedFormats{
				HTML: true, SimpleHTML: true, RSS: true, Atom: true,
				JSON: true, Markdown: true, Text: true, Sitemap: true,
			},
		},
		{
			"partial formats",
			map[string]interface{}{
				"html": true,
				"rss":  true,
			},
			&models.FeedFormats{HTML: true, RSS: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseFeedFormatsFromMap(tt.raw)
			if tt.want == nil {
				if got != nil {
					t.Errorf("expected nil, got %+v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected non-nil, got nil")
			}
			if *got != *tt.want {
				t.Errorf("got %+v, want %+v", *got, *tt.want)
			}
		})
	}
}

// =============================================================================
// Interface Compliance
// =============================================================================

func TestSeriesPlugin_ImplementsInterfaces(_ *testing.T) {
	var _ lifecycle.Plugin = (*SeriesPlugin)(nil)
	var _ lifecycle.CollectPlugin = (*SeriesPlugin)(nil)
	var _ lifecycle.PriorityPlugin = (*SeriesPlugin)(nil)
}
