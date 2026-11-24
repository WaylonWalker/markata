package plugins

import (
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// =============================================================================
// FeedsPlugin Tests
// =============================================================================

func TestFeedsPlugin_Name(t *testing.T) {
	plugin := NewFeedsPlugin()
	if plugin.Name() != "feeds" {
		t.Errorf("Name() = %q, want %q", plugin.Name(), "feeds")
	}
}

func TestFeedsPlugin_BasicFeed(t *testing.T) {
	// Create a manager with posts
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Date: &date1},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2"), Date: &date2},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Post 3"), Date: &date3},
	})

	// Configure a basic feed
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:  "all",
				Title: "All Posts",
			},
		},
	}
	m.SetConfig(config)

	// Run the plugin
	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// Verify feed was created
	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	feed := feeds[0]
	if feed.Name != "all" {
		t.Errorf("feed.Name = %q, want %q", feed.Name, "all")
	}
	if feed.Title != "All Posts" {
		t.Errorf("feed.Title = %q, want %q", feed.Title, "All Posts")
	}
	if len(feed.Posts) != 3 {
		t.Errorf("feed should have 3 posts, got %d", len(feed.Posts))
	}

	// Verify default sorting (by date, reversed - newest first)
	if feed.Posts[0].Slug != "post1" {
		t.Errorf("first post should be post1 (newest), got %q", feed.Posts[0].Slug)
	}
	if feed.Posts[2].Slug != "post3" {
		t.Errorf("last post should be post3 (oldest), got %q", feed.Posts[2].Slug)
	}
}

func TestFeedsPlugin_FilteredFeed(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Python Tutorial"), Tags: []string{"python", "tutorial"}, Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Go Basics"), Tags: []string{"go", "tutorial"}, Date: &date},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Python Advanced"), Tags: []string{"python", "advanced"}, Date: &date},
	})

	// Configure a filtered feed for Python posts
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:   "python",
				Title:  "Python Posts",
				Filter: `'python' in tags`,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	// Should only have Python posts
	if len(feeds[0].Posts) != 2 {
		t.Errorf("expected 2 Python posts, got %d", len(feeds[0].Posts))
	}

	for _, post := range feeds[0].Posts {
		hasPython := false
		for _, tag := range post.Tags {
			if tag == "python" {
				hasPython = true
				break
			}
		}
		if !hasPython {
			t.Errorf("post %q should have python tag", post.Slug)
		}
	}
}

func TestFeedsPlugin_SortedFeed(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post B"), Date: &date1},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post A"), Date: &date2},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Post C"), Date: &date3},
	})

	// Configure feed sorted by title (ascending)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:    "by-title",
				Title:   "Posts by Title",
				Sort:    "title",
				Reverse: false,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	// Verify alphabetical order by title
	posts := feeds[0].Posts
	if *posts[0].Title != "Post A" {
		t.Errorf("first post should be 'Post A', got %q", *posts[0].Title)
	}
	if *posts[1].Title != "Post B" {
		t.Errorf("second post should be 'Post B', got %q", *posts[1].Title)
	}
	if *posts[2].Title != "Post C" {
		t.Errorf("third post should be 'Post C', got %q", *posts[2].Title)
	}
}

func TestFeedsPlugin_ReverseSorting(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post B"), Date: &date1},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post A"), Date: &date2},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Post C"), Date: &date3},
	})

	// Configure feed sorted by title (reversed - descending)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:    "by-title-desc",
				Title:   "Posts by Title Descending",
				Sort:    "title",
				Reverse: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	posts := feeds[0].Posts

	// Verify reverse alphabetical order by title
	if *posts[0].Title != "Post C" {
		t.Errorf("first post should be 'Post C', got %q", *posts[0].Title)
	}
	if *posts[2].Title != "Post A" {
		t.Errorf("last post should be 'Post A', got %q", *posts[2].Title)
	}
}

func TestFeedsPlugin_LimitOffset(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	date4 := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Date: &date1},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2"), Date: &date2},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Post 3"), Date: &date3},
		{Path: "post4.md", Slug: "post4", Title: strPtr("Post 4"), Date: &date4},
	})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:   "limited",
				Title:  "Limited",
				Limit:  2,
				Offset: 1,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	posts := feeds[0].Posts
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts after limit/offset, got %d", len(posts))
	}
	if posts[0].Slug != "post2" {
		t.Errorf("first post should be post2 after offset, got %q", posts[0].Slug)
	}
	if posts[1].Slug != "post3" {
		t.Errorf("second post should be post3 after limit, got %q", posts[1].Slug)
	}
}

func TestFeedsPlugin_RecordsDistinctFeedSelectionReasons(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	posts := []*models.Post{
		{Path: "included.md", Slug: "included", Tags: []string{"keep"}, Published: true, Date: &date},
		{Path: "filtered.md", Slug: "filtered", Tags: []string{"other"}, Published: true, Date: &date},
		{Path: "private.md", Slug: "private", Tags: []string{"keep"}, Published: true, Private: true, Date: &date},
		{Path: "draft.md", Slug: "draft", Tags: []string{"keep"}, Published: true, Draft: true, Date: &date},
		{Path: "skipped.md", Slug: "skipped", Tags: []string{"keep"}, Published: true, Skip: true, Date: &date},
		{Path: "unpublished.md", Slug: "unpublished", Tags: []string{"keep"}, Published: false, Date: &date},
	}

	m := lifecycle.NewManager()
	m.SetPosts(posts)
	m.SetFiles([]string{"included.md", "filtered.md", "private.md", "draft.md", "skipped.md", "unpublished.md"})
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{{
			Slug:   "diagnostics",
			Filter: "published == true and 'keep' in tags",
		}},
	}
	m.SetConfig(config)

	if err := NewFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	assertFeedReason := func(path, reason string, included bool) {
		t.Helper()
		entry := findContentDisposition(m.ContentDiagnostics(), path)
		if entry == nil {
			t.Fatalf("missing content entry for %s", path)
		}
		var feed *diagnostics.ContentFeedDisposition
		for index := range entry.Feeds {
			if entry.Feeds[index].Feed == "diagnostics" {
				feed = &entry.Feeds[index]
				break
			}
		}
		if feed == nil {
			t.Fatalf("missing diagnostics feed entry for %s: %+v", path, entry.Feeds)
		}
		if feed.Included != included {
			t.Errorf("%s included = %t, want %t (%v)", path, feed.Included, included, feed.Reasons)
		}
		if reason != "" && !containsReason(feed.Reasons, reason) {
			t.Errorf("%s feed reasons = %v, want %q", path, feed.Reasons, reason)
		}
	}

	assertFeedReason("included.md", "", true)
	assertFeedReason("filtered.md", diagnostics.ReasonContentFiltered, false)
	assertFeedReason("private.md", diagnostics.ReasonContentPrivate, false)
	assertFeedReason("draft.md", diagnostics.ReasonContentDraft, false)
	assertFeedReason("skipped.md", diagnostics.ReasonContentSkip, false)
	assertFeedReason("unpublished.md", diagnostics.ReasonContentPublishedFalse, false)

	filtered := findContentDisposition(m.ContentDiagnostics(), "filtered.md")
	if filtered == nil || !containsReason(filtered.Feeds[0].Reasons, diagnostics.ReasonContentFiltered) {
		t.Fatalf("filtered feed disposition = %+v", filtered)
	}
	private := findContentDisposition(m.ContentDiagnostics(), "private.md")
	if private == nil || containsReason(private.Feeds[0].Reasons, diagnostics.ReasonContentFiltered) {
		t.Fatalf("private feed disposition was mislabeled as filtered: %+v", private)
	}
	unpublished := findContentDisposition(m.ContentDiagnostics(), "unpublished.md")
	if unpublished == nil || containsReason(unpublished.Feeds[0].Reasons, diagnostics.ReasonContentFiltered) {
		t.Fatalf("unpublished feed disposition was mislabeled as filtered: %+v", unpublished)
	}
}

func TestFeedsPlugin_RecordsOffsetAndLimitReasonsAfterFilteringAndSorting(t *testing.T) {
	dates := []time.Time{
		time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC),
	}
	posts := []*models.Post{
		{Path: "first.md", Slug: "first", Tags: []string{"keep"}, Published: true, Date: &dates[0]},
		{Path: "second.md", Slug: "second", Tags: []string{"keep"}, Published: true, Date: &dates[1]},
		{Path: "third.md", Slug: "third", Tags: []string{"keep"}, Published: true, Date: &dates[2]},
		{Path: "fourth.md", Slug: "fourth", Tags: []string{"keep"}, Published: true, Date: &dates[3]},
	}

	m := lifecycle.NewManager()
	m.SetPosts(posts)
	m.SetFiles([]string{"first.md", "second.md", "third.md", "fourth.md"})
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{{
			Slug:   "window",
			Filter: "'keep' in tags",
			Offset: 1,
			Limit:  2,
		}},
	}
	m.SetConfig(config)

	if err := NewFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	feeds := m.Feeds()
	if len(feeds) != 1 || len(feeds[0].Posts) != 2 || feeds[0].Posts[0].Slug != "second" || feeds[0].Posts[1].Slug != "third" {
		t.Fatalf("window feed posts = %+v", feeds)
	}

	for path, reason := range map[string]string{
		"first.md":  diagnostics.ReasonFeedOffset,
		"fourth.md": diagnostics.ReasonFeedLimit,
	} {
		entry := findContentDisposition(m.ContentDiagnostics(), path)
		if entry == nil || len(entry.Feeds) != 1 {
			t.Fatalf("%s feed disposition = %+v", path, entry)
		}
		feed := entry.Feeds[0]
		if feed.Included || !containsReason(feed.Reasons, reason) {
			t.Errorf("%s feed disposition = %+v, want excluded for %s", path, feed, reason)
		}
		if containsReason(feed.Reasons, diagnostics.ReasonContentFiltered) {
			t.Errorf("%s feed disposition = %+v, must not be labeled filtered", path, feed)
		}
	}
	for _, path := range []string{"second.md", "third.md"} {
		entry := findContentDisposition(m.ContentDiagnostics(), path)
		if entry == nil || len(entry.Feeds) != 1 || !entry.Feeds[0].Included {
			t.Errorf("%s feed disposition = %+v, want included", path, entry)
		}
	}
}

func TestFeedsPlugin_IncrementalRebuildsWithoutSerializedFeedPosts(t *testing.T) {
	post := &models.Post{Path: "post.md", Slug: "post", Published: true, Content: "body"}
	m := lifecycle.NewManager()
	m.SetPosts([]*models.Post{post})
	m.SetFiles([]string{"post.md"})
	m.Config().Extra = map[string]interface{}{
		"feeds_incremental": true,
		"feeds":             []models.FeedConfig{{Slug: "archive"}},
	}
	m.Cache().Set("build_cache", buildcache.New(t.TempDir()))

	if err := NewFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	feeds := m.Feeds()
	if len(feeds) != 1 || len(feeds[0].Posts) != 1 || feeds[0].Posts[0].Path != post.Path {
		t.Fatalf("incremental feeds = %+v, want runtime feed membership restored", feeds)
	}
	entry := findContentDisposition(m.ContentDiagnostics(), post.Path)
	if entry == nil || len(entry.Feeds) != 1 || !entry.Feeds[0].Included {
		t.Fatalf("incremental feed diagnostics = %+v, want included post", entry)
	}
}

func TestFeedsPlugin_Pagination(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create 15 posts
	var posts []*models.Post
	for i := 1; i <= 15; i++ {
		posts = append(posts, &models.Post{
			Path:  "post" + string(rune('0'+i)) + ".md",
			Slug:  "post-" + string(rune('0'+i)),
			Title: strPtr("Post " + string(rune('0'+i))),
			Date:  &date,
		})
	}
	m.SetPosts(posts)

	// Configure feed with 5 items per page
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:         "paginated",
				Title:        "Paginated Posts",
				ItemsPerPage: 5,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	// Check that feed_configs are stored in cache with pagination info
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

	// 15 posts / 5 per page = 3 pages
	if len(fc.Pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(fc.Pages))
	}

	// Verify page sizes
	if len(fc.Pages) >= 1 && len(fc.Pages[0].Posts) != 5 {
		t.Errorf("page 1 should have 5 posts, got %d", len(fc.Pages[0].Posts))
	}
	if len(fc.Pages) >= 2 && len(fc.Pages[1].Posts) != 5 {
		t.Errorf("page 2 should have 5 posts, got %d", len(fc.Pages[1].Posts))
	}
	if len(fc.Pages) >= 3 && len(fc.Pages[2].Posts) != 5 {
		t.Errorf("page 3 should have 5 posts, got %d", len(fc.Pages[2].Posts))
	}
}

func TestFeedsPlugin_OrphanThreshold(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create 12 posts: with 5 per page, that's 2 pages of 5, plus 2 orphans
	var posts []*models.Post
	for i := 1; i <= 12; i++ {
		posts = append(posts, &models.Post{
			Path:  "post.md",
			Slug:  "post",
			Title: strPtr("Post"),
			Date:  &date,
		})
	}
	m.SetPosts(posts)

	// Configure feed with orphan threshold of 3 (remaining 2 items should merge into previous page)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:            "orphan-test",
				Title:           "Orphan Test",
				ItemsPerPage:    5,
				OrphanThreshold: 3,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	cached, _ := m.Cache().Get("feed_configs")
	feedConfigs, _ := cached.([]models.FeedConfig) //nolint:errcheck // test code
	fc := feedConfigs[0]

	// With 12 posts, 5 per page, and orphan threshold 3:
	// Normal: 5 + 5 + 2 = 3 pages, but 2 < 3 (threshold)
	// So: 5 + 7 = 2 pages (last 2 merged into page 2)
	if len(fc.Pages) != 2 {
		t.Errorf("expected 2 pages (with orphan merge), got %d", len(fc.Pages))
	}

	if len(fc.Pages) >= 2 {
		// Second page should have 7 posts (5 + 2 orphans)
		if len(fc.Pages[1].Posts) != 7 {
			t.Errorf("page 2 should have 7 posts (with orphans), got %d", len(fc.Pages[1].Posts))
		}
	}
}

func TestFeedsPlugin_FeedDefaults(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	var posts []*models.Post
	for i := 1; i <= 25; i++ {
		posts = append(posts, &models.Post{
			Path:  "post.md",
			Slug:  "post",
			Title: strPtr("Post"),
			Date:  &date,
		})
	}
	m.SetPosts(posts)

	// Configure defaults with 10 items per page
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feed_defaults": models.FeedDefaults{
			ItemsPerPage:    10,
			OrphanThreshold: 5,
			Formats: models.FeedFormats{
				HTML: true,
				RSS:  true,
			},
		},
		"feeds": []models.FeedConfig{
			{
				Slug:  "with-defaults",
				Title: "Feed Using Defaults",
				// ItemsPerPage not set - should inherit from defaults
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	cached, _ := m.Cache().Get("feed_configs")
	feedConfigs, _ := cached.([]models.FeedConfig) //nolint:errcheck // test code
	fc := feedConfigs[0]

	// Should have inherited ItemsPerPage = 10
	if fc.ItemsPerPage != 10 {
		t.Errorf("ItemsPerPage should be 10 (from defaults), got %d", fc.ItemsPerPage)
	}

	// 25 posts / 10 per page, with orphan threshold 5
	// Normal: 10 + 10 + 5 = 3 pages, but 5 >= 5 (threshold) so no merge
	if len(fc.Pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(fc.Pages))
	}
}

func TestFeedsPlugin_MultipleFeedsConfig(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Python Post"), Tags: []string{"python"}, Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Go Post"), Tags: []string{"go"}, Date: &date},
		{Path: "post3.md", Slug: "post3", Title: strPtr("Tutorial"), Tags: []string{"tutorial"}, Date: &date},
	})

	// Configure multiple feeds
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:   "all",
				Title:  "All Posts",
				Filter: "",
			},
			{
				Slug:   "python",
				Title:  "Python Posts",
				Filter: `'python' in tags`,
			},
			{
				Slug:   "go",
				Title:  "Go Posts",
				Filter: `'go' in tags`,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 3 {
		t.Fatalf("expected 3 feeds, got %d", len(feeds))
	}

	// Check each feed has correct posts
	feedByName := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedByName[f.Name] = f
	}

	if allFeed, ok := feedByName["all"]; ok {
		if len(allFeed.Posts) != 3 {
			t.Errorf("'all' feed should have 3 posts, got %d", len(allFeed.Posts))
		}
	} else {
		t.Error("'all' feed not found")
	}

	if pyFeed, ok := feedByName["python"]; ok {
		if len(pyFeed.Posts) != 1 {
			t.Errorf("'python' feed should have 1 post, got %d", len(pyFeed.Posts))
		}
	} else {
		t.Error("'python' feed not found")
	}

	if goFeed, ok := feedByName["go"]; ok {
		if len(goFeed.Posts) != 1 {
			t.Errorf("'go' feed should have 1 post, got %d", len(goFeed.Posts))
		}
	} else {
		t.Error("'go' feed not found")
	}
}

func TestFeedsPlugin_EmptyFilter(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Date: &date},
		{Path: "post2.md", Slug: "post2", Title: strPtr("Post 2"), Date: &date},
	})

	// Configure feed with empty filter (should get all posts)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:   "all",
				Title:  "All Posts",
				Filter: "",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds[0].Posts) != 2 {
		t.Errorf("empty filter should return all posts, got %d", len(feeds[0].Posts))
	}
}

func TestFeedsPlugin_InvalidFilter(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Date: &date},
	})

	// Configure feed with invalid filter (unclosed parenthesis)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:   "invalid",
				Title:  "Invalid Filter",
				Filter: "( unclosed paren",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err == nil {
		t.Error("expected error for invalid filter, got nil")
	}
}

func TestFeedsPlugin_NoPosts(t *testing.T) {
	m := lifecycle.NewManager()
	m.SetPosts([]*models.Post{})

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:  "empty",
				Title: "Empty Feed",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}
	if len(feeds[0].Posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(feeds[0].Posts))
	}

	// Check that pagination creates empty pages slice
	cached, _ := m.Cache().Get("feed_configs")
	feedConfigs, _ := cached.([]models.FeedConfig) //nolint:errcheck // test code
	if len(feedConfigs[0].Pages) != 0 {
		t.Errorf("expected 0 pages for empty feed, got %d", len(feedConfigs[0].Pages))
	}
}

func TestFeedsPlugin_NoFeedsConfigured(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "post1.md", Slug: "post1", Title: strPtr("Post 1"), Date: &date},
	})

	// No feeds in config
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 0 {
		t.Errorf("expected 0 feeds when none configured, got %d", len(feeds))
	}
}

func TestFeedsPlugin_SortByExtraField(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{
			Path:  "post1.md",
			Slug:  "post1",
			Title: strPtr("Post 1"),
			Date:  &date,
			Extra: map[string]interface{}{"priority": 3},
		},
		{
			Path:  "post2.md",
			Slug:  "post2",
			Title: strPtr("Post 2"),
			Date:  &date,
			Extra: map[string]interface{}{"priority": 1},
		},
		{
			Path:  "post3.md",
			Slug:  "post3",
			Title: strPtr("Post 3"),
			Date:  &date,
			Extra: map[string]interface{}{"priority": 2},
		},
	})

	// Sort by custom "priority" field
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:    "by-priority",
				Title:   "By Priority",
				Sort:    "priority",
				Reverse: false,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	posts := feeds[0].Posts

	// Should be sorted by priority ascending: post2 (1), post3 (2), post1 (3)
	if posts[0].Slug != "post2" {
		t.Errorf("first post should be post2 (priority 1), got %q", posts[0].Slug)
	}
	if posts[1].Slug != "post3" {
		t.Errorf("second post should be post3 (priority 2), got %q", posts[1].Slug)
	}
	if posts[2].Slug != "post1" {
		t.Errorf("third post should be post1 (priority 3), got %q", posts[2].Slug)
	}
}

func TestFeedsPlugin_PageNavigation(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create 10 posts for 2 pages
	var posts []*models.Post
	for i := 1; i <= 10; i++ {
		posts = append(posts, &models.Post{
			Path:  "post.md",
			Slug:  "post",
			Title: strPtr("Post"),
			Date:  &date,
		})
	}
	m.SetPosts(posts)

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:         "nav-test",
				Title:        "Navigation Test",
				ItemsPerPage: 5,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	cached, _ := m.Cache().Get("feed_configs")
	feedConfigs, _ := cached.([]models.FeedConfig) //nolint:errcheck // test code
	fc := feedConfigs[0]

	if len(fc.Pages) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(fc.Pages))
	}

	// Page 1: HasPrev=false, HasNext=true
	page1 := fc.Pages[0]
	if page1.HasPrev {
		t.Error("page 1 should not have previous page")
	}
	if !page1.HasNext {
		t.Error("page 1 should have next page")
	}
	if page1.NextURL != "/nav-test/page/2/" {
		t.Errorf("page 1 NextURL = %q, want %q", page1.NextURL, "/nav-test/page/2/")
	}

	// Page 2: HasPrev=true, HasNext=false
	page2 := fc.Pages[1]
	if !page2.HasPrev {
		t.Error("page 2 should have previous page")
	}
	if page2.HasNext {
		t.Error("page 2 should not have next page")
	}
	if page2.PrevURL != "/nav-test/" {
		t.Errorf("page 2 PrevURL = %q, want %q", page2.PrevURL, "/nav-test/")
	}
}

// =============================================================================
// Helper Function Tests
// =============================================================================

func TestFilterPosts(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	posts := []*models.Post{
		{Slug: "post1", Tags: []string{"python"}, Date: &date, Template: "photo"},
		{Slug: "post2", Tags: []string{"go"}, Date: &date, Template: "video"},
		{Slug: "post3", Tags: []string{"python", "go"}, Date: &date, Template: "shot"},
	}

	tests := []struct {
		name      string
		filter    string
		wantCount int
		wantSlugs []string
		wantErr   bool
	}{
		{
			name:      "empty filter returns all",
			filter:    "",
			wantCount: 3,
		},
		{
			name:      "filter by tag",
			filter:    `'python' in tags`,
			wantCount: 2,
			wantSlugs: []string{"post1", "post3"},
		},
		{
			name:      "filter by templateKey alias",
			filter:    `templateKey == "shots"`,
			wantCount: 2,
			wantSlugs: []string{"post1", "post3"},
		},
		{
			name:    "invalid filter",
			filter:  "not valid @#$",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filterPosts(posts, tt.filter, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(result) != tt.wantCount {
				t.Errorf("got %d posts, want %d", len(result), tt.wantCount)
			}
		})
	}
}

func TestFilterPosts_IncludePrivate(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	posts := []*models.Post{
		{Slug: "public1", Tags: []string{"python"}, Date: &date, Private: false, ArticleHTML: "<p>Public 1</p>"},
		{Slug: "private1", Tags: []string{"python"}, Date: &date, Private: true, ArticleHTML: "<p>Private 1</p>"},
		{Slug: "public2", Tags: []string{"go"}, Date: &date, Private: false, ArticleHTML: "<p>Public 2</p>"},
		{Slug: "private2", Tags: []string{"go"}, Date: &date, Private: true, ArticleHTML: "<p>Private 2</p>"},
	}

	tests := []struct {
		name           string
		filter         string
		includePrivate bool
		wantCount      int
		wantSlugs      []string
	}{
		{
			name:           "exclude private by default",
			filter:         "",
			includePrivate: false,
			wantCount:      2,
			wantSlugs:      []string{"public1", "public2"},
		},
		{
			name:           "include private when flag is true",
			filter:         "",
			includePrivate: true,
			wantCount:      4,
			wantSlugs:      []string{"public1", "private1", "public2", "private2"},
		},
		{
			name:           "exclude private with filter",
			filter:         `'python' in tags`,
			includePrivate: false,
			wantCount:      1,
			wantSlugs:      []string{"public1"},
		},
		{
			name:           "include private with filter",
			filter:         `'python' in tags`,
			includePrivate: true,
			wantCount:      2,
			wantSlugs:      []string{"public1", "private1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filterPosts(posts, tt.filter, tt.includePrivate)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != tt.wantCount {
				t.Errorf("got %d posts, want %d", len(result), tt.wantCount)
			}
			// Verify the slugs
			resultSlugs := make([]string, len(result))
			for i, p := range result {
				resultSlugs[i] = p.Slug
			}
			for _, wantSlug := range tt.wantSlugs {
				found := false
				for _, gotSlug := range resultSlugs {
					if gotSlug == wantSlug {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected slug %q not found in result %v", wantSlug, resultSlugs)
				}
			}
		})
	}
}

func TestFeedsPlugin_IncludePrivate(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{Path: "public1.md", Slug: "public1", Title: strPtr("Public 1"), Date: &date, Private: false},
		{Path: "private1.md", Slug: "private1", Title: strPtr("Private 1"), Date: &date, Private: true},
		{Path: "public2.md", Slug: "public2", Title: strPtr("Public 2"), Date: &date, Private: false},
	})

	// Configure two feeds: one excluding private (default), one including
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:           "public-only",
				Title:          "Public Only",
				IncludePrivate: false, // default behavior
			},
			{
				Slug:           "with-private",
				Title:          "With Private",
				IncludePrivate: true,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 2 {
		t.Fatalf("expected 2 feeds, got %d", len(feeds))
	}

	// Build map for easier access
	feedByName := make(map[string]*lifecycle.Feed)
	for _, f := range feeds {
		feedByName[f.Name] = f
	}

	// Public-only feed should have 2 posts
	if publicFeed, ok := feedByName["public-only"]; ok {
		if len(publicFeed.Posts) != 2 {
			t.Errorf("'public-only' feed should have 2 posts, got %d", len(publicFeed.Posts))
		}
		// Verify no private posts
		for _, post := range publicFeed.Posts {
			if post.Private {
				t.Errorf("'public-only' feed should not contain private post %q", post.Slug)
			}
		}
	} else {
		t.Error("'public-only' feed not found")
	}

	// With-private feed should have 3 posts
	if privateFeed, ok := feedByName["with-private"]; ok {
		if len(privateFeed.Posts) != 3 {
			t.Errorf("'with-private' feed should have 3 posts, got %d", len(privateFeed.Posts))
		}
		// Verify private post is included
		hasPrivate := false
		for _, post := range privateFeed.Posts {
			if post.Private {
				hasPrivate = true
				break
			}
		}
		if !hasPrivate {
			t.Error("'with-private' feed should contain at least one private post")
		}
	} else {
		t.Error("'with-private' feed not found")
	}
}

func TestFeedsPlugin_PrivateAliasIncludesPrivate(t *testing.T) {
	m := lifecycle.NewManager()
	m.SetPosts([]*models.Post{
		{Path: "public.md", Slug: "public", Href: "/public/", Published: true},
		{Path: "private.md", Slug: "private", Href: "/private/", Published: true, Private: true},
	})
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{{Slug: "private", Private: true}},
	}
	m.SetConfig(config)

	if err := NewFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error: %v", err)
	}
	feeds := m.Feeds()
	if len(feeds) != 1 || !feeds[0].IncludePrivate {
		t.Fatalf("private feed alias was not normalized: %#v", feeds)
	}
	if len(feeds[0].Posts) != 2 {
		t.Fatalf("private feed posts = %d, want 2", len(feeds[0].Posts))
	}
}

func TestSortPosts(t *testing.T) {
	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2024, 1, 5, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		posts     []*models.Post
		field     string
		reverse   bool
		wantOrder []string
	}{
		{
			name: "sort by date ascending",
			posts: []*models.Post{
				{Slug: "post1", Date: &date1},
				{Slug: "post2", Date: &date2},
				{Slug: "post3", Date: &date3},
			},
			field:     "date",
			reverse:   false,
			wantOrder: []string{"post3", "post2", "post1"},
		},
		{
			name: "sort by date descending",
			posts: []*models.Post{
				{Slug: "post1", Date: &date1},
				{Slug: "post2", Date: &date2},
				{Slug: "post3", Date: &date3},
			},
			field:     "date",
			reverse:   true,
			wantOrder: []string{"post1", "post2", "post3"},
		},
		{
			name: "sort by title ascending",
			posts: []*models.Post{
				{Slug: "post1", Title: strPtr("C Title")},
				{Slug: "post2", Title: strPtr("A Title")},
				{Slug: "post3", Title: strPtr("B Title")},
			},
			field:     "title",
			reverse:   false,
			wantOrder: []string{"post2", "post3", "post1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortPosts(tt.posts, tt.field, tt.reverse)

			for i, wantSlug := range tt.wantOrder {
				if tt.posts[i].Slug != wantSlug {
					t.Errorf("position %d: got %q, want %q", i, tt.posts[i].Slug, wantSlug)
				}
			}
		})
	}
}

func TestGetFieldValue(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	post := &models.Post{
		Slug:  "test-post",
		Title: strPtr("Test Title"),
		Date:  &date,
		Extra: map[string]interface{}{
			"custom": "value",
			"number": 42,
		},
	}

	tests := []struct {
		name  string
		field string
		want  interface{}
	}{
		{"slug field", "slug", "test-post"},
		{"title field (pointer)", "title", "Test Title"},
		{"date field (pointer)", "date", date},
		{"extra custom field", "custom", "value"},
		{"extra number field", "number", 42},
		{"non-existent field", "nonexistent", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFieldValue(post, tt.field)
			if got != tt.want {
				t.Errorf("getFieldValue(%q) = %v, want %v", tt.field, got, tt.want)
			}
		})
	}
}

func TestCompareFieldValues(t *testing.T) {
	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		a    interface{}
		b    interface{}
		want int // -1, 0, or 1
	}{
		{"both nil", nil, nil, 0},
		{"a nil", nil, "x", -1},
		{"b nil", "x", nil, 1},
		{"equal strings", "abc", "abc", 0},
		{"strings a < b", "abc", "xyz", -1},
		{"strings a > b", "xyz", "abc", 1},
		{"dates equal", date1, date1, 0},
		{"dates a < b", date2, date1, -1},
		{"dates a > b", date1, date2, 1},
		{"integers equal", 5, 5, 0},
		{"integers a < b", 1, 10, -1},
		{"integers a > b", 10, 1, 1},
		{"int64 comparison", int64(5), int64(10), -1},
		{"float64 as int", float64(3), float64(7), -1},
		{"mixed int types", int(5), int64(10), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareFieldValues(tt.a, tt.b)
			// Normalize to -1, 0, 1 for comparison
			if got < 0 {
				got = -1
			} else if got > 0 {
				got = 1
			}
			if got != tt.want {
				t.Errorf("compareFieldValues(%v, %v) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// =============================================================================
// Interface Compliance
// =============================================================================

func TestFeedsPlugin_ImplementsInterfaces(_ *testing.T) {
	var _ lifecycle.Plugin = (*FeedsPlugin)(nil)
	var _ lifecycle.CollectPlugin = (*FeedsPlugin)(nil)
}

// =============================================================================
// Guide Series Tests
// =============================================================================

func TestFeedsPlugin_GuideTypeFeed(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create posts with guide_order in Extra
	m.SetPosts([]*models.Post{
		{
			Path:  "guide3.md",
			Slug:  "guide3",
			Title: strPtr("Guide Part 3"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 3},
		},
		{
			Path:  "guide1.md",
			Slug:  "guide1",
			Title: strPtr("Guide Part 1"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 1},
		},
		{
			Path:  "guide2.md",
			Slug:  "guide2",
			Title: strPtr("Guide Part 2"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 2},
		},
	})

	// Configure a guide-type feed
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:  "getting-started",
				Title: "Getting Started Guide",
				Type:  models.FeedTypeGuide,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	// Verify posts are sorted by guide_order (ascending)
	posts := feeds[0].Posts
	if len(posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(posts))
	}

	if posts[0].Slug != "guide1" {
		t.Errorf("first post should be guide1 (order 1), got %q", posts[0].Slug)
	}
	if posts[1].Slug != "guide2" {
		t.Errorf("second post should be guide2 (order 2), got %q", posts[1].Slug)
	}
	if posts[2].Slug != "guide3" {
		t.Errorf("third post should be guide3 (order 3), got %q", posts[2].Slug)
	}
}

func TestFeedsPlugin_GuideNavigation(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Create posts with guide_order
	m.SetPosts([]*models.Post{
		{
			Path:  "guide1.md",
			Slug:  "guide1",
			Title: strPtr("Guide Part 1"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 1},
		},
		{
			Path:  "guide2.md",
			Slug:  "guide2",
			Title: strPtr("Guide Part 2"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 2},
		},
		{
			Path:  "guide3.md",
			Slug:  "guide3",
			Title: strPtr("Guide Part 3"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 3},
		},
	})

	// Configure a guide-type feed
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:  "tutorial",
				Title: "Tutorial Series",
				Type:  models.FeedTypeGuide,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	posts := feeds[0].Posts

	// First post: no prev, has next
	if posts[0].Prev != nil {
		t.Error("first post should have no Prev")
	}
	if posts[0].Next == nil || posts[0].Next.Slug != "guide2" {
		t.Errorf("first post.Next should be guide2, got %v", posts[0].Next)
	}
	if posts[0].PrevNextFeed != "tutorial" {
		t.Errorf("first post.PrevNextFeed should be 'tutorial', got %q", posts[0].PrevNextFeed)
	}

	// Middle post: has prev and next
	if posts[1].Prev == nil || posts[1].Prev.Slug != "guide1" {
		t.Errorf("middle post.Prev should be guide1, got %v", posts[1].Prev)
	}
	if posts[1].Next == nil || posts[1].Next.Slug != "guide3" {
		t.Errorf("middle post.Next should be guide3, got %v", posts[1].Next)
	}

	// Last post: has prev, no next
	if posts[2].Prev == nil || posts[2].Prev.Slug != "guide2" {
		t.Errorf("last post.Prev should be guide2, got %v", posts[2].Prev)
	}
	if posts[2].Next != nil {
		t.Error("last post should have no Next")
	}
}

func TestFeedsPlugin_SeriesTypeFeed(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{
			Path:  "part1.md",
			Slug:  "part1",
			Title: strPtr("Series Part 1"),
			Date:  &date,
		},
		{
			Path:  "part2.md",
			Slug:  "part2",
			Title: strPtr("Series Part 2"),
			Date:  &date,
		},
	})

	// Configure a series-type feed
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:  "my-series",
				Title: "My Series",
				Type:  models.FeedTypeSeries,
				Sort:  "title",
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	posts := feeds[0].Posts

	// Verify navigation is set up for series type too
	if posts[0].Next == nil || posts[0].Next.Slug != "part2" {
		t.Error("series posts should have prev/next navigation")
	}
	if posts[0].PrevNextFeed != "my-series" {
		t.Errorf("PrevNextFeed should be 'my-series', got %q", posts[0].PrevNextFeed)
	}
}

func TestFeedsPlugin_BlogTypeDefaultBehavior(t *testing.T) {
	m := lifecycle.NewManager()

	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{
			Path:  "post1.md",
			Slug:  "post1",
			Title: strPtr("Older Post"),
			Date:  &date2,
		},
		{
			Path:  "post2.md",
			Slug:  "post2",
			Title: strPtr("Newer Post"),
			Date:  &date1,
		},
	})

	// Configure a blog-type feed (should sort by date descending by default)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:  "blog",
				Title: "Blog",
				Type:  models.FeedTypeBlog,
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	posts := feeds[0].Posts

	// Blog type should NOT set up guide navigation
	// Newer post should be first (reverse date sort)
	if posts[0].Slug != "post2" {
		t.Errorf("blog feed should sort by date descending, expected post2 first, got %q", posts[0].Slug)
	}
}

func TestFeedsPlugin_GuideWithExplicitSort(t *testing.T) {
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{
			Path:  "guide-c.md",
			Slug:  "guide-c",
			Title: strPtr("Guide C"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 3},
		},
		{
			Path:  "guide-a.md",
			Slug:  "guide-a",
			Title: strPtr("Guide A"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 1},
		},
		{
			Path:  "guide-b.md",
			Slug:  "guide-b",
			Title: strPtr("Guide B"),
			Date:  &date,
			Extra: map[string]interface{}{"guide_order": 2},
		},
	})

	// Configure guide feed with explicit sort override to title
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{
			{
				Slug:  "guides",
				Title: "Guides",
				Type:  models.FeedTypeGuide,
				Sort:  "title", // Override default guide_order sorting
			},
		},
	}
	m.SetConfig(config)

	plugin := NewFeedsPlugin()
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error: %v", err)
	}

	feeds := m.Feeds()
	posts := feeds[0].Posts

	// Should be sorted by title (ascending) since we explicitly set Sort
	if posts[0].Slug != "guide-a" {
		t.Errorf("expected guide-a first (by title), got %q", posts[0].Slug)
	}
	if posts[1].Slug != "guide-b" {
		t.Errorf("expected guide-b second (by title), got %q", posts[1].Slug)
	}
	if posts[2].Slug != "guide-c" {
		t.Errorf("expected guide-c third (by title), got %q", posts[2].Slug)
	}
}

// =============================================================================
// Test Helpers
// =============================================================================

// strPtr returns a pointer to a string value.
func strPtr(s string) *string {
	return &s
}
