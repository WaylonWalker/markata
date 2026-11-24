package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestPublishSidebarChunks_BoundedPublicAndPrunesOldChunks(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	posts := make([]*models.Post, 0, 142)
	for i := 0; i < 140; i++ {
		title := fmt.Sprintf("Post %d", i)
		posts = append(posts, &models.Post{Slug: fmt.Sprintf("post-%d", i), Href: fmt.Sprintf("/post-%d/", i), Title: &title, Published: true})
	}
	privateTitle := "private-should-never-leak"
	posts = append(posts,
		&models.Post{Slug: "private", Href: "/secret/", Title: &privateTitle, Published: true, Private: true},
		&models.Post{Slug: "skipped", Href: "/skipped/", Skip: true, Published: true},
	)
	fc := &models.FeedConfig{Slug: "blog", Formats: models.FeedFormats{HTML: true}, Posts: posts}
	p := NewPublishFeedsPlugin()
	if err := p.publishSidebarChunks(fc, dir); err != nil {
		t.Fatal(err)
	}
	for i, want := range []int{64, 64, 12} {
		data, err := os.ReadFile(filepath.Join(dir, "_sidebar", fmt.Sprintf("%04d.json", i)))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(data), "private-should-never-leak") || strings.Contains(string(data), "/secret/") || strings.Contains(string(data), "/skipped/") {
			t.Fatal("sidebar chunks leaked non-public posts")
		}
		var chunk sidebarNavChunk
		if err := json.Unmarshal(data, &chunk); err != nil {
			t.Fatal(err)
		}
		if chunk.Start != i*sidebarNavChunkSize || len(chunk.Posts) != want {
			t.Fatalf("chunk %d start=%d posts=%d, expected start=%d posts=%d", i, chunk.Start, len(chunk.Posts), i*sidebarNavChunkSize, want)
		}
	}
	fc.Posts = posts[:7]
	if err := p.publishSidebarChunks(fc, dir); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "_sidebar", "0002.json")); !os.IsNotExist(err) {
		t.Fatalf("old chunk was not pruned: %v", err)
	}
	fc.IncludePrivate = true
	if sidebarChunksEnabled(fc) {
		t.Fatal("private feeds must never expose static sidebar chunks")
	}
}

func TestSidebarChunkBaseAndSingleFeedHydration(t *testing.T) {
	t.Parallel()
	m := lifecycle.NewManager()
	enabled := true
	m.Config().Extra["components"] = models.ComponentsConfig{FeedSidebar: models.FeedSidebarConfig{Enabled: &enabled, MaxPosts: 1}}
	title := "First post"
	posts := []*models.Post{
		{Slug: "first", Href: "/first/", Title: &title, Published: true},
		{Slug: "second", Href: "/second/", Published: true},
		{Slug: "third", Href: "/third/", Published: true},
	}
	feed := models.FeedConfig{Slug: "blog", Primary: true, Formats: models.FeedFormats{HTML: true}, Posts: posts}
	m.Cache().Set("feed_configs", []models.FeedConfig{feed})
	p := NewTemplatesPlugin()
	got := p.buildSidebarFeedsJSON(posts[1], m.Config(), m, &feed)
	if got == "" {
		t.Fatal("single-feed sidebar must include lazy navigation metadata")
	}
	var data sidebarFeedsDataJSON
	if err := json.Unmarshal([]byte(got), &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Feeds) != 1 {
		t.Fatalf("expected one feed, got %d", len(data.Feeds))
	}
	entry := data.Feeds[0]
	if entry.ChunkBase != "/blog/_sidebar/" || entry.ChunkSize != sidebarNavChunkSize {
		t.Fatalf("missing static chunk metadata: %#v", entry)
	}
	if entry.TotalPosts != 3 || len(entry.Posts) != 1 || entry.WindowStart != 1 || entry.Posts[0].Slug != "second" {
		t.Fatalf("unexpected initial window: %#v", entry)
	}
	// A feed containing an item that would be filtered during publishing must
	// not advertise a sidecar with mismatched offsets.
	posts[0].Private = true
	if source := sidebarChunkBase(&feed, posts, m); source != "" {
		t.Fatalf("private post should make sidecar offsets unsafe: %q", source)
	}
}

func TestCompactAlternateSidebarPosts(t *testing.T) {
	t.Parallel()
	feeds := []sidebarFeedJSON{
		{Slug: "current", ChunkBase: "/current/_sidebar/", ChunkSize: 64, TotalPosts: 2, Posts: []sidebarPostJSON{{Slug: "first"}}},
		{Slug: "other", ChunkBase: "/other/_sidebar/", ChunkSize: 64, TotalPosts: 2, WindowStart: 1, Posts: []sidebarPostJSON{{Slug: "second"}}},
		{Slug: "synthetic", TotalPosts: 3, Posts: []sidebarPostJSON{{Slug: "third"}}},
		{Slug: "unsafe", ChunkBase: "/unsafe/_sidebar/", TotalPosts: 3, Posts: []sidebarPostJSON{{Slug: "fourth"}}},
	}
	compactAlternateSidebarPosts(feeds, 0)
	if len(feeds[0].Posts) != 1 || len(feeds[1].Posts) != 0 || len(feeds[2].Posts) != 1 || len(feeds[3].Posts) != 1 {
		t.Fatalf("unexpected alternate feed windows: %#v", feeds)
	}
	if feeds[1].WindowStart != 1 || feeds[1].ChunkBase != "/other/_sidebar/" {
		t.Fatalf("lost alternate hydration metadata: %#v", feeds[1])
	}
}

func TestSidebarChunkBaseRequiresMatchingPublicURLsAndTitles(t *testing.T) {
	t.Parallel()
	m := lifecycle.NewManager()
	originalTitle := "Published title"
	feedPosts := []*models.Post{
		{Slug: "one", Href: "/one/", Title: &originalTitle, Published: true},
		{Slug: "two", Href: "/two/", Published: true},
	}
	published := models.FeedConfig{
		Slug:    "blog",
		Formats: models.FeedFormats{HTML: true},
		Posts:   feedPosts,
	}
	m.Cache().Set("feed_configs", []models.FeedConfig{published})

	matchingTitle := "Published title"
	matching := []*models.Post{
		{Slug: "one", Href: "/one/", Title: &matchingTitle, Published: true},
		{Slug: "two", Href: "/two/", Published: true},
	}
	if got := sidebarChunkBase(&published, matching, m); got != "/blog/_sidebar/" {
		t.Fatalf("matching navigation must permit deferred sidecar: %q", got)
	}

	staleURL := []*models.Post{
		{Slug: "one", Href: "/old-one/", Title: &matchingTitle, Published: true},
		{Slug: "two", Href: "/two/", Published: true},
	}
	if got := sidebarChunkBase(&published, staleURL, m); got != "" {
		t.Fatalf("mismatched href must retain embedded fallback, got sidecar %q", got)
	}
	staleTitle := "Outdated title"
	staleTitles := []*models.Post{
		{Slug: "one", Href: "/one/", Title: &staleTitle, Published: true},
		{Slug: "two", Href: "/two/", Published: true},
	}
	if got := sidebarChunkBase(&published, staleTitles, m); got != "" {
		t.Fatalf("mismatched title must retain embedded fallback, got sidecar %q", got)
	}
	feedData := []sidebarFeedJSON{
		{Slug: "current", Posts: []sidebarPostJSON{{Slug: "visible"}}},
		{Slug: "old", ChunkBase: sidebarChunkBase(&published, staleURL, m), ChunkSize: sidebarNavChunkSize, TotalPosts: 2, Posts: []sidebarPostJSON{{Slug: "one", Href: "/old-one/"}}},
	}
	compactAlternateSidebarPosts(feedData, 0)
	if len(feedData[1].Posts) != 1 {
		t.Fatal("alternate window must not be omitted when its sidecar URL is unsafe")
	}
}
