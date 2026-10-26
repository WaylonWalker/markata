package listcache

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

func TestLoadChangedPosts_ParsesAllChangedFiles(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	changed := make(map[string]bool, 12)
	for i := 0; i < 12; i++ {
		path := filepath.Join("post-", string(rune('a'+i)), "index.md")
		fullPath := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) error = %v", fullPath, err)
		}
		if err := os.WriteFile(fullPath, []byte("---\ntitle: Test\n---\n\nContent"), 0o600); err != nil {
			t.Fatalf("WriteFile(%q) error = %v", fullPath, err)
		}
		changed[path] = true
	}

	posts, err := loadChangedPosts(dir, changed, nil, 4)
	if err != nil {
		t.Fatalf("loadChangedPosts() error = %v", err)
	}
	if len(posts) != len(changed) {
		t.Fatalf("loadChangedPosts() returned %d posts, want %d", len(posts), len(changed))
	}
	for _, post := range posts {
		if !changed[post.Path] {
			t.Errorf("unexpected post path %q", post.Path)
		}
	}
}

func TestCachedPost_PostStatsRoundTrip(t *testing.T) {
	t.Parallel()

	want := &plugins.PostStats{
		WordCount:       123,
		CharCount:       456,
		ReadingTime:     7,
		ReadingTimeText: "7 min read",
		CodeLines:       8,
		CodeBlocks:      2,
	}
	post := models.NewPost("stats.md")
	post.Set("stats", want)

	restored := roundTripCachedPost(t, post)
	got, ok := restored.Get("stats").(*plugins.PostStats)
	if !ok {
		t.Fatalf("restored stats type = %T, want *plugins.PostStats", restored.Get("stats"))
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("restored stats = %+v, want %+v", got, want)
	}
}

func TestCachedPost_PostStatsRoundTripPreservesAggregateAuthority(t *testing.T) {
	t.Parallel()

	content := strings.TrimSpace(strings.Repeat("prose ", 101)) +
		"\n\n```go\n" + strings.TrimSpace(strings.Repeat("code ", 100)) + "\n```"
	transformManager := lifecycle.NewManager()
	transformManager.SetConcurrency(1)
	transformManager.Config().Extra = map[string]interface{}{
		"stats": map[string]interface{}{
			"words_per_minute":      100,
			"include_code_in_count": true,
		},
	}
	post := models.NewPost("divergent.md")
	post.Content = content
	transformManager.SetPosts([]*models.Post{post})

	statsPlugin := plugins.NewStatsPlugin()
	if err := statsPlugin.Configure(transformManager); err != nil {
		t.Fatalf("configure stats: %v", err)
	}
	if err := statsPlugin.Transform(transformManager); err != nil {
		t.Fatalf("transform stats: %v", err)
	}
	readingTimePlugin := plugins.NewReadingTimePlugin()
	if err := readingTimePlugin.Configure(transformManager); err != nil {
		t.Fatalf("configure reading time: %v", err)
	}
	if err := readingTimePlugin.Transform(transformManager); err != nil {
		t.Fatalf("transform reading time: %v", err)
	}

	statsBefore, ok := post.Get("stats").(*plugins.PostStats)
	if !ok {
		t.Fatalf("stats before round trip = %T, want *plugins.PostStats", post.Get("stats"))
	}
	publicWordCount, ok := post.Get("word_count").(int)
	if !ok {
		t.Fatalf("public word_count = %T, want int", post.Get("word_count"))
	}
	publicReadingTime, ok := post.Get("reading_time").(int)
	if !ok {
		t.Fatalf("public reading_time = %T, want int", post.Get("reading_time"))
	}
	if statsBefore.WordCount == publicWordCount || statsBefore.ReadingTime == publicReadingTime {
		t.Fatalf(
			"test setup did not diverge: Stats = (%d, %d), public ReadingTime = (%d, %d)",
			statsBefore.WordCount,
			statsBefore.ReadingTime,
			publicWordCount,
			publicReadingTime,
		)
	}

	restored := roundTripCachedPost(t, post)
	aggregateManager := lifecycle.NewManager()
	aggregateManager.SetPosts([]*models.Post{restored})
	aggregateManager.SetFeeds([]*lifecycle.Feed{{Name: "posts", Posts: []*models.Post{restored}}})
	aggregator := plugins.NewStatsPlugin()
	if err := aggregator.Collect(aggregateManager); err != nil {
		t.Fatalf("collect stats: %v", err)
	}

	feedStats := plugins.GetFeedStats(aggregateManager, "posts")
	if feedStats == nil {
		t.Fatal("feed stats were not stored")
	}
	if feedStats.TotalWords != statsBefore.WordCount || feedStats.TotalReadingTime != statsBefore.ReadingTime {
		t.Fatalf(
			"feed aggregate = (%d, %d), want Stats values (%d, %d), not public values (%d, %d)",
			feedStats.TotalWords,
			feedStats.TotalReadingTime,
			statsBefore.WordCount,
			statsBefore.ReadingTime,
			publicWordCount,
			publicReadingTime,
		)
	}

	siteStats := plugins.GetSiteStats(aggregateManager)
	if siteStats == nil {
		t.Fatal("site stats were not stored")
	}
	if siteStats.TotalWords != statsBefore.WordCount || siteStats.TotalReadingTime != statsBefore.ReadingTime {
		t.Fatalf(
			"site aggregate = (%d, %d), want Stats values (%d, %d), not public values (%d, %d)",
			siteStats.TotalWords,
			siteStats.TotalReadingTime,
			statsBefore.WordCount,
			statsBefore.ReadingTime,
			publicWordCount,
			publicReadingTime,
		)
	}
}

func TestCachedPost_PostStatsRoundTripPreservesDisabledCodeMetrics(t *testing.T) {
	t.Parallel()

	manager := lifecycle.NewManager()
	manager.SetConcurrency(1)
	manager.Config().Extra = map[string]interface{}{
		"stats": map[string]interface{}{
			"track_code_blocks": false,
		},
	}
	post := models.NewPost("no-code-metrics.md")
	post.Content = "ordinary prose\n\n```go\nalpha\nbeta\ngamma\n```"
	manager.SetPosts([]*models.Post{post})

	statsPlugin := plugins.NewStatsPlugin()
	if err := statsPlugin.Configure(manager); err != nil {
		t.Fatalf("configure stats: %v", err)
	}
	if err := statsPlugin.Transform(manager); err != nil {
		t.Fatalf("transform stats: %v", err)
	}

	before, ok := post.Get("stats").(*plugins.PostStats)
	if !ok {
		t.Fatalf("stats before round trip = %T, want *plugins.PostStats", post.Get("stats"))
	}
	if before.CodeBlocks != 0 || before.CodeLines != 0 {
		t.Fatalf("stats before round trip code metrics = (%d, %d), want (0, 0)", before.CodeBlocks, before.CodeLines)
	}

	restored := roundTripCachedPost(t, post)
	after, ok := restored.Get("stats").(*plugins.PostStats)
	if !ok {
		t.Fatalf("stats after round trip = %T, want *plugins.PostStats", restored.Get("stats"))
	}
	if after.CodeBlocks != 0 || after.CodeLines != 0 {
		t.Fatalf("stats after round trip code metrics = (%d, %d), want (0, 0)", after.CodeBlocks, after.CodeLines)
	}
}

func TestCachedPost_WithoutStatsRoundTripsNormally(t *testing.T) {
	t.Parallel()

	post := models.NewPost("without-stats.md")
	post.Set("custom", map[string]interface{}{"value": "kept"})

	restored := roundTripCachedPost(t, post)
	if restored.Get("stats") != nil {
		t.Fatal("restored post unexpectedly contains stats")
	}
	custom, ok := restored.Get("custom").(map[string]interface{})
	if !ok {
		t.Fatalf("restored custom value type = %T, want map[string]interface{}", restored.Get("custom"))
	}
	if custom["value"] != "kept" {
		t.Fatalf("restored custom value = %#v, want kept", custom["value"])
	}
}

func TestCachedPost_MalformedStatsDoesNotPanic(t *testing.T) {
	t.Parallel()

	cache := newCache("test", ".", nil)
	cache.Posts["malformed.md"] = CachedPost{
		Path: "malformed.md",
		Extra: map[string]interface{}{
			"stats": map[string]interface{}{
				"word_count": "not a number",
			},
		},
	}
	path := filepath.Join(t.TempDir(), CacheFileName)
	if err := saveCache(path, cache); err != nil {
		t.Fatalf("saveCache() error = %v", err)
	}
	loaded, err := loadCache(path)
	if err != nil {
		t.Fatalf("loadCache() error = %v", err)
	}
	restored := cachedPostToModel(loaded.Posts["malformed.md"])
	if _, ok := restored.Get("stats").(*plugins.PostStats); ok {
		t.Fatal("malformed stats unexpectedly restored as *plugins.PostStats")
	}
}

func roundTripCachedPost(t *testing.T, post *models.Post) *models.Post {
	t.Helper()

	cache := newCache("test", ".", nil)
	cache.Posts[post.Path] = modelToCachedPost(post)
	path := filepath.Join(t.TempDir(), CacheFileName)
	if err := saveCache(path, cache); err != nil {
		t.Fatalf("saveCache() error = %v", err)
	}
	loaded, err := loadCache(path)
	if err != nil {
		t.Fatalf("loadCache() error = %v", err)
	}
	cached, ok := loaded.Posts[post.Path]
	if !ok {
		t.Fatalf("post %q missing from restored cache", post.Path)
	}
	return cachedPostToModel(cached)
}
