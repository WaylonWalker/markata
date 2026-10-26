package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// =============================================================================
// AutoFeedsPlugin Tests
// =============================================================================

func TestAutoFeedsPlugin_Name(t *testing.T) {
	plugin := NewAutoFeedsPlugin()
	if plugin.Name() != "auto_feeds" {
		t.Errorf("Name() = %q, want %q", plugin.Name(), "auto_feeds")
	}
}

func TestAutoFeedsPlugin_TagFeeds(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Python Tutorial"), Tags: []string{"python", "tutorial"}, Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Go Basics"), Tags: []string{"go", "tutorial"}, Date: &date},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Python Advanced"), Tags: []string{"python", "advanced"}, Date: &date},
	})

	// Configure auto feeds for tags
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "tags",
				Formats: models.FeedFormats{
					HTML: true,
					RSS:  true,
				},
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Should have 4 tag feeds: advanced, go, python, tutorial (alphabetical)
	if len(feeds) != 4 {
		t.Fatalf("expected 4 tag feeds, got %d", len(feeds))
	}

	// Create a map for easy lookup
	feedMap := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedMap[f.Name] = f
	}

	// Check that expected tag feeds exist with correct titles
	expectedFeeds := map[string]string{
		"tags/python":   "Posts tagged: python",
		"tags/go":       "Posts tagged: go",
		"tags/tutorial": "Posts tagged: tutorial",
		"tags/advanced": "Posts tagged: advanced",
	}

	for slug, expectedTitle := range expectedFeeds {
		feed, ok := feedMap[slug]
		if !ok {
			t.Errorf("expected feed %q not found", slug)
			continue
		}
		if feed.Title != expectedTitle {
			t.Errorf("feed %q title = %q, want %q", slug, feed.Title, expectedTitle)
		}
	}
}

func TestAutoFeedsPlugin_CustomSlugPrefix(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Tags: []string{"python"}, Date: &date},
	})

	// Configure with custom slug prefix
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "topics",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	// Should use custom prefix
	if feeds[0].Name != "topics/python" {
		t.Errorf("feed slug = %q, want %q", feeds[0].Name, "topics/python")
	}
}

func TestAutoFeedsPlugin_TagFeedsRobotsFromRawConfig(t *testing.T) {
	m := lifecycle.NewManager()
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{{Path: "post1.md", Slug: "post1", Title: strPtr("Python Tutorial"), Tags: []string{"python"}, Date: &date}})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": map[string]interface{}{
			"tags": map[string]interface{}{
				"enabled":     true,
				"slug_prefix": "tags",
				"robots":      "noindex,follow",
				"formats": map[string]interface{}{
					"html": true,
					"rss":  true,
				},
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("expected feed_configs cache entry")
	}
	configs, ok := cached.([]models.FeedConfig)
	if !ok || len(configs) != 1 {
		t.Fatalf("expected 1 feed config, got %#v", cached)
	}
	if configs[0].Robots != "noindex,follow" {
		t.Fatalf("feed robots = %q, want %q", configs[0].Robots, "noindex,follow")
	}
}

func TestAutoFeedsPlugin_TagFeeds_CollapsesSlugCollisions(t *testing.T) {
	m := lifecycle.NewManager()
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Tags: []string{"open source", "3d printing"}, Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2"), Tags: []string{"open-source", "3d-printing"}, Date: &date},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{Enabled: true, SlugPrefix: "tags"},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 2 {
		t.Fatalf("expected 2 collapsed tag feeds, got %d", len(feeds))
	}

	feedMap := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedMap[f.Name] = f
	}

	feed, ok := feedMap["tags/open-source"]
	if !ok {
		t.Fatalf("expected tags/open-source feed")
	}
	if feed.Title != "Posts tagged: open source" {
		t.Fatalf("feed title = %q, want %q", feed.Title, "Posts tagged: open source")
	}

	allConfigs, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs missing from cache")
	}
	configs, ok := allConfigs.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs has wrong type: %T", allConfigs)
	}
	var matched *models.FeedConfig
	for i := range configs {
		if configs[i].Slug == "tags/open-source" {
			matched = &configs[i]
			break
		}
	}
	if matched == nil {
		t.Fatal("tags/open-source config not found")
	}
	wantFilter := `"open source" in tags or "open-source" in tags`
	if matched.Filter != wantFilter {
		t.Fatalf("filter = %q, want %q", matched.Filter, wantFilter)
	}
}

func TestAutoFeedsPlugin_GenerateTagFeedsForChanged_UsesCanonicalGroup(t *testing.T) {
	plugin := NewAutoFeedsPlugin()
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	posts := []*models.Post{
		{Path: "post1.md", Slug: "post1", Tags: []string{"open source"}, Date: &date},
		{Path: "post2.md", Slug: "post2", Tags: []string{"open-source"}, Date: &date},
	}

	feeds := plugin.generateTagFeedsForChanged(posts, AutoFeedTypeConfig{SlugPrefix: "tags"}, map[string]bool{"post2": true}, map[string]bool{})
	if len(feeds) != 1 {
		t.Fatalf("expected 1 changed feed, got %d", len(feeds))
	}
	if feeds[0].Title != "Posts tagged: open source" {
		t.Fatalf("title = %q, want %q", feeds[0].Title, "Posts tagged: open source")
	}
	wantFilter := `"open source" in tags or "open-source" in tags`
	if feeds[0].Filter != wantFilter {
		t.Fatalf("filter = %q, want %q", feeds[0].Filter, wantFilter)
	}
}

func TestAutoFeedsPlugin_CategoryFeeds(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{
			Path:    "post1.md",
			Slug:    "post1",
			Title:   strPtr("Tech Post 1"),
			Content: "tech post 1",
			Date:    &date,
			Extra:   map[string]interface{}{"category": "Technology"},
		},
		{
			Path:    "post2.md",
			Slug:    "post2",
			Title:   strPtr("Life Post"),
			Content: "life post",
			Date:    &date,
			Extra:   map[string]interface{}{"category": "Lifestyle"},
		},
		{
			Path:    "post3.md",
			Slug:    "post3",
			Title:   strPtr("Tech Post 2"),
			Content: "tech post 2",
			Date:    &date,
			Extra:   map[string]interface{}{"category": "Technology"},
		},
	})

	// Configure auto feeds for categories
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Categories: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "categories",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Should have 2 category feeds: Lifestyle, Technology
	if len(feeds) != 2 {
		t.Fatalf("expected 2 category feeds, got %d", len(feeds))
	}

	feedMap := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedMap[f.Name] = f
	}

	// Check technology feed
	if techFeed, ok := feedMap["categories/technology"]; ok {
		if len(techFeed.Posts) != 2 {
			t.Errorf("technology category should have 2 posts, got %d", len(techFeed.Posts))
		}
		if techFeed.Title != "Category: Technology" {
			t.Errorf("technology feed title = %q, want %q", techFeed.Title, "Category: Technology")
		}
	} else {
		t.Error("categories/technology feed not found")
	}

	// Check lifestyle feed
	if lifeFeed, ok := feedMap["categories/lifestyle"]; ok {
		if len(lifeFeed.Posts) != 1 {
			t.Errorf("lifestyle category should have 1 post, got %d", len(lifeFeed.Posts))
		}
	} else {
		t.Error("categories/lifestyle feed not found")
	}
}

func TestAutoFeedsPlugin_YearlyArchives(t *testing.T) {
	m := lifecycle.NewManager()

	date2024 := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	date2023a := time.Date(2023, 3, 10, 0, 0, 0, 0, time.UTC)
	date2023b := time.Date(2023, 8, 20, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 2024"), Date: &date2024},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2023 A"), Date: &date2023a},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Post 2023 B"), Date: &date2023b},
	})

	// Configure yearly archive feeds
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Archives: AutoArchiveConfig{
				Enabled:      true,
				SlugPrefix:   "archive",
				YearlyFeeds:  true,
				MonthlyFeeds: false,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Should have 2 yearly archive feeds: 2024, 2023
	if len(feeds) != 2 {
		t.Fatalf("expected 2 yearly archive feeds, got %d", len(feeds))
	}

	feedMap := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedMap[f.Name] = f
	}

	// Check 2024 archive exists with correct title
	if feed2024, ok := feedMap["archive/2024"]; ok {
		if feed2024.Title != "Archive: 2024" {
			t.Errorf("2024 archive title = %q, want %q", feed2024.Title, "Archive: 2024")
		}
	} else {
		t.Error("archive/2024 feed not found")
	}

	// Check 2023 archive exists
	if _, ok := feedMap["archive/2023"]; !ok {
		t.Error("archive/2023 feed not found")
	}
}

func TestAutoFeedsPlugin_MonthlyArchives(t *testing.T) {
	m := lifecycle.NewManager()

	dateJan := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	dateMarA := time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)
	dateMarB := time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("January Post"), Date: &dateJan},
		{Path: "post2.md", Slug: "post2", Title: strPtr("March Post A"), Date: &dateMarA},
		{Path: "post3.md", Slug: "post3", Title: strPtr("March Post B"), Date: &dateMarB},
	})

	// Configure monthly archive feeds
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Archives: AutoArchiveConfig{
				Enabled:      true,
				SlugPrefix:   "archive",
				YearlyFeeds:  false,
				MonthlyFeeds: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Should have 2 monthly archive feeds: 2024/01, 2024/03
	if len(feeds) != 2 {
		t.Fatalf("expected 2 monthly archive feeds, got %d", len(feeds))
	}

	feedMap := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedMap[f.Name] = f
	}

	// Check January archive exists with correct title
	if feedJan, ok := feedMap["archive/2024/01"]; ok {
		if !strings.Contains(feedJan.Title, "January") {
			t.Errorf("January archive title should contain 'January', got %q", feedJan.Title)
		}
	} else {
		t.Error("archive/2024/01 feed not found")
	}

	// Check March archive exists
	if _, ok := feedMap["archive/2024/03"]; !ok {
		t.Error("archive/2024/03 feed not found")
	}
}

func TestAutoFeedsPlugin_CombinedArchives(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post"), Date: &date},
	})

	// Configure both yearly and monthly archive feeds
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Archives: AutoArchiveConfig{
				Enabled:      true,
				SlugPrefix:   "archive",
				YearlyFeeds:  true,
				MonthlyFeeds: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Should have 2 feeds: yearly 2024 + monthly 2024/01
	if len(feeds) != 2 {
		t.Fatalf("expected 2 archive feeds, got %d", len(feeds))
	}

	hasYearly := false
	hasMonthly := false
	for _, f := range feeds {
		if f.Name == "archive/2024" {
			hasYearly = true
		}
		if f.Name == "archive/2024/01" {
			hasMonthly = true
		}
	}

	if !hasYearly {
		t.Error("expected yearly archive feed")
	}
	if !hasMonthly {
		t.Error("expected monthly archive feed")
	}
}

func TestAutoFeedsPlugin_NoAutoFeedsConfigured(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post"), Tags: []string{"test"}, Date: &date},
	})

	// No auto_feeds config (tag feeds enabled by default)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	// With tag feeds enabled by default, we should get 1 feed for the "test" tag
	if len(feeds) != 1 {
		t.Errorf("expected 1 feed (tag feed) with default config, got %d", len(feeds))
	}
}

func TestAutoFeedsPlugin_AllDisabled(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post"), Tags: []string{"test"}, Date: &date},
	})

	// All auto feeds explicitly disabled
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags:       AutoFeedTypeConfig{Enabled: false},
			Categories: AutoFeedTypeConfig{Enabled: false},
			Archives:   AutoArchiveConfig{Enabled: false},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds when all disabled, got %d", len(feeds))
	}
}

func TestAutoFeedsPlugin_InheritsFeedDefaults(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create enough posts to trigger pagination
	var posts []*models.Post
	for i := 0; i < 15; i++ {
		posts = append(posts, &models.Post{
			Path:    "post.md",
			Slug:    "post",
			Title:   strPtr("Post"),
			Tags:    []string{"python"},
			Content: "post",
			Date:    &date,
		})
	}
	m.SetPosts(posts)

	// Configure feed defaults and tag auto feeds
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feed_defaults": models.FeedDefaults{
			ItemsPerPage:    5,
			OrphanThreshold: 2,
		},
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// Check that feed configs inherit defaults
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}

	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs has wrong type: %T", cached)
	}

	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	fc := feedConfigs[0]

	// Should have inherited defaults
	if fc.ItemsPerPage != 5 {
		t.Errorf("ItemsPerPage should be 5 (from defaults), got %d", fc.ItemsPerPage)
	}

	// 15 posts / 5 per page = 3 pages
	if len(fc.Pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(fc.Pages))
	}
}

func TestAutoFeedsPlugin_PostsSortedByDateDescending(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "oldest", Title: strPtr("Oldest"), Tags: []string{"python"}, Content: "oldest", Date: &date3},
		{Path: "post2.md", Slug: "newest", Title: strPtr("Newest"), Tags: []string{"python"}, Content: "newest", Date: &date1},
		{Path: "post3.md", Slug: "middle", Title: strPtr("Middle"), Tags: []string{"python"}, Content: "middle", Date: &date2},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	posts := feeds[0].Posts

	// Should be sorted by date descending (newest first)
	if posts[0].Slug != "newest" {
		t.Errorf("first post should be 'newest', got %q", posts[0].Slug)
	}
	if posts[1].Slug != "middle" {
		t.Errorf("second post should be 'middle', got %q", posts[1].Slug)
	}
	if posts[2].Slug != "oldest" {
		t.Errorf("third post should be 'oldest', got %q", posts[2].Slug)
	}
}

func TestAutoFeedsPlugin_SlugifyTag(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Tags with special characters that need slugification
	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post"), Tags: []string{"C++ Programming", "Machine Learning"}, Date: &date},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Check that slugs are properly formatted
	for _, f := range feeds {
		// Slugs should be lowercase with hyphens
		if strings.ContainsAny(f.Name, " +") {
			t.Errorf("feed slug should not contain spaces or +: %q", f.Name)
		}
		if f.Name != strings.ToLower(f.Name) {
			// Check slug part only (after prefix)
			parts := strings.SplitN(f.Name, "/", 2)
			if len(parts) == 2 && parts[1] != strings.ToLower(parts[1]) {
				t.Errorf("feed slug should be lowercase: %q", f.Name)
			}
		}
	}
}

func TestAutoFeedsPlugin_NoTagsNoPosts(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Posts without tags
	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Tags: []string{}, Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2"), Date: &date}, // nil tags
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// No tags = no tag feeds
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds when no tags, got %d", len(feeds))
	}
}

func TestAutoFeedsPlugin_NoCategoriesNoPosts(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Posts without category
	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Date: &date},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Categories: AutoFeedTypeConfig{
				Enabled: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// No categories = no category feeds
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds when no categories, got %d", len(feeds))
	}
}

func TestAutoFeedsPlugin_NoDatesNoArchives(t *testing.T) {
	m := lifecycle.NewManager()

	// Posts without dates
	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1")},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2")},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Archives: AutoArchiveConfig{
				Enabled:     true,
				YearlyFeeds: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// No dates = no archive feeds
	if len(feeds) != 0 {
		t.Errorf("expected 0 archive feeds when no dates, got %d", len(feeds))
	}
}

func TestAutoFeedsPlugin_AppendsToExistingFeeds(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post"), Tags: []string{"python"}, Date: &date},
	})

	// Pre-populate existing feeds (as if FeedsPlugin ran first)
	existingFeed := &lifecycle.Feed{
		Name:  "existing",
		Title: "Existing Feed",
	}
	m.SetFeeds([]*lifecycle.Feed{existingFeed})

	// Also set existing feed configs in cache
	m.Cache().Set("feed_configs", []models.FeedConfig{
		{Slug: "existing", Title: "Existing Feed"},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Should have both existing and auto-generated feeds
	if len(feeds) != 2 {
		t.Fatalf("expected 2 feeds (1 existing + 1 auto), got %d", len(feeds))
	}

	hasExisting := false
	hasPython := false
	for _, f := range feeds {
		if f.Name == "existing" {
			hasExisting = true
		}
		if f.Name == "tags/python" {
			hasPython = true
		}
	}

	if !hasExisting {
		t.Error("existing feed should be preserved")
	}
	if !hasPython {
		t.Error("auto-generated python tag feed should exist")
	}
}

func TestAutoFeedsPlugin_TagFeedFilterExpression(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Python Tutorial"), Tags: []string{"python", "tutorial"}, Content: "python tutorial", Date: &date, Published: true},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Go Basics"), Tags: []string{"go", "tutorial"}, Content: "go basics", Date: &date, Published: true},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Python Advanced"), Tags: []string{"python", "advanced"}, Content: "python advanced", Date: &date, Published: true},
	})

	// Configure auto feeds for tags
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "tags",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()

	// Create a map for easy lookup
	feedMap := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedMap[f.Name] = f
	}

	// Check python feed has correct posts (only posts with "python" tag)
	pythonFeed, ok := feedMap["tags/python"]
	if !ok {
		t.Fatal("tags/python feed not found")
	}

	// Python feed should have exactly 2 posts (post1 and post3)
	if len(pythonFeed.Posts) != 2 {
		t.Errorf("python feed should have 2 posts, got %d", len(pythonFeed.Posts))
	}

	// Verify the posts are the correct ones
	slugs := make(map[string]bool)
	for _, p := range pythonFeed.Posts {
		slugs[p.Slug] = true
	}
	if !slugs["post1"] || !slugs["post3"] {
		t.Errorf("python feed should contain post1 and post3, got %v", slugs)
	}
	if slugs["post2"] {
		t.Errorf("python feed should NOT contain post2 (go post)")
	}

	// Check go feed has correct posts (only posts with "go" tag)
	goFeed, ok := feedMap["tags/go"]
	if !ok {
		t.Fatal("tags/go feed not found")
	}

	// Go feed should have exactly 1 post (post2)
	if len(goFeed.Posts) != 1 {
		t.Errorf("go feed should have 1 post, got %d", len(goFeed.Posts))
	}
	if goFeed.Posts[0].Slug != "post2" {
		t.Errorf("go feed should contain post2, got %s", goFeed.Posts[0].Slug)
	}
}

// =============================================================================
// Private Tag Feed Tests
// =============================================================================

func TestAutoFeedsPlugin_PrivateTagIncludesPrivatePost(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("A"), Tags: []string{"python"}, Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("B"), Tags: []string{"gratitude"}, Date: &date, Private: true},
		{Path: "post3.md", Slug: "post3", Title: strPtr("C"), Tags: []string{"tutorial"}, Date: &date},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "tags",
			},
		},
		"models_config": &models.Config{
			Encryption: models.EncryptionConfig{
				Enabled:     true,
				PrivateTags: map[string]string{"gratitude": "default"},
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// Extract feed configs from cache to check IncludePrivate and posts
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}

	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs has wrong type: %T", cached)
	}

	// Build map of feed configs by slug
	configMap := make(map[string]models.FeedConfig)
	for _, fc := range feedConfigs {
		configMap[fc.Slug] = fc
	}

	// Private tag feed should include private posts (with encrypted content).
	gratitudeFeed, ok := configMap["tags/gratitude"]
	if !ok {
		t.Fatal("tags/gratitude feed config not found")
	}
	if !gratitudeFeed.IncludePrivate {
		t.Error("tags/gratitude feed should have IncludePrivate=true")
	}
	if len(gratitudeFeed.Posts) != 1 {
		t.Errorf("tags/gratitude feed should have 1 post, got %d", len(gratitudeFeed.Posts))
	}

	// Python feed should NOT have IncludePrivate
	pythonFeed, ok := configMap["tags/python"]
	if !ok {
		t.Fatal("tags/python feed config not found")
	}
	if pythonFeed.IncludePrivate {
		t.Error("tags/python feed should have IncludePrivate=false, got true")
	}
	if len(pythonFeed.Posts) != 1 {
		t.Errorf("tags/python feed should have 1 post, got %d", len(pythonFeed.Posts))
	}

	// Tutorial feed should NOT have IncludePrivate
	tutorialFeed, ok := configMap["tags/tutorial"]
	if !ok {
		t.Fatal("tags/tutorial feed config not found")
	}
	if tutorialFeed.IncludePrivate {
		t.Error("tags/tutorial feed should have IncludePrivate=false, got true")
	}
	if len(tutorialFeed.Posts) != 1 {
		t.Errorf("tags/tutorial feed should have 1 post, got %d", len(tutorialFeed.Posts))
	}
}

func TestAutoFeedsPlugin_PrivateTagCaseInsensitiveIncludesPrivatePost(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Entry"), Tags: []string{"Gratitude"}, Date: &date, Private: true},
	})

	// Private tag in config is lowercase, tag on post is capitalized
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "tags",
			},
		},
		"models_config": &models.Config{
			Encryption: models.EncryptionConfig{
				Enabled:     true,
				PrivateTags: map[string]string{"gratitude": "default"},
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}

	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs has wrong type: %T", cached)
	}

	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	if !feedConfigs[0].IncludePrivate {
		t.Error("feed for tag 'Gratitude' should have IncludePrivate=true")
	}
	if len(feedConfigs[0].Posts) != 1 {
		t.Errorf("feed should have 1 post, got %d", len(feedConfigs[0].Posts))
	}
}

func TestAutoFeedsPlugin_IncrementalStillKeepsUnchangedAutoFeeds(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Tutorial"), Tags: []string{"tutorial"}, Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Gratitude"), Tags: []string{"gratitude"}, Date: &date, Private: true},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds_incremental": true,
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "tags",
				Formats: models.FeedFormats{
					HTML: true,
				},
			},
		},
		"models_config": &models.Config{
			Encryption: models.EncryptionConfig{
				PrivateTags: map[string]string{"gratitude": "default"},
			},
		},
	}
	m.SetConfig(config)

	cache := buildcache.New(t.TempDir())
	cache.MarkFeedSlugChanged("post2")
	m.Cache().Set("build_cache", cache)

	plugin := NewAutoFeedsPlugin()
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feedsValue, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}
	feedConfigs, ok := feedsValue.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs has wrong type: %T", feedsValue)
	}

	feedMap := make(map[string]models.FeedConfig, len(feedConfigs))
	for _, fc := range feedConfigs {
		feedMap[fc.Slug] = fc
	}

	if _, ok := feedMap["tags/tutorial"]; !ok {
		t.Fatal("tags/tutorial feed should remain present during incremental auto-feed rebuilds")
	}
	gratitudeFeed, ok := feedMap["tags/gratitude"]
	if !ok {
		t.Fatal("tags/gratitude feed should remain present during incremental auto-feed rebuilds")
	}
	if !gratitudeFeed.IncludePrivate {
		t.Fatal("tags/gratitude should keep IncludePrivate=true during incremental auto-feed rebuilds")
	}
}

func TestAutoFeedsPlugin_NoEncryptionConfig(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post"), Tags: []string{"python"}, Date: &date},
	})

	// No models_config at all — should not crash, IncludePrivate should be false
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"auto_feeds": AutoFeedsConfig{
			Tags: AutoFeedTypeConfig{
				Enabled:    true,
				SlugPrefix: "tags",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewAutoFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}

	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs has wrong type: %T", cached)
	}

	if len(feedConfigs) != 1 {
		t.Fatalf("expected 1 feed config, got %d", len(feedConfigs))
	}

	if feedConfigs[0].IncludePrivate {
		t.Error("IncludePrivate should be false when no encryption config exists")
	}
}

// =============================================================================
// Slugify Tests
// =============================================================================

func TestSlugify(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Python", "python"},
		{"Machine Learning", "machine-learning"},
		{"C++", "c"},
		{"Web Dev!", "web-dev"},
		{"  spaces  ", "spaces"},
		{"multiple---hyphens", "multiple-hyphens"},
		{"UPPERCASE", "uppercase"},
		{"Mix3d C4se", "mix3d-c4se"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := slugify(tt.input)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// =============================================================================
// Interface Compliance
// =============================================================================

func TestAutoFeedsPlugin_ImplementsInterfaces(_ *testing.T) {
	var _ lifecycle.Plugin = (*AutoFeedsPlugin)(nil)
	var _ lifecycle.CollectPlugin = (*AutoFeedsPlugin)(nil)
}
