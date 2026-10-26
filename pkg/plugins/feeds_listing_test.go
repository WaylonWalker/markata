package plugins

import (
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestFeedsListingPlugin_Write(t *testing.T) {
	plugin := NewFeedsListingPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.OutputDir = t.TempDir()
	feedsPage := models.NewFeedsPageConfig()
	feedsPage.Robots = "noindex,follow"
	defaults := models.NewFeedDefaults()
	config.Extra = map[string]interface{}{
		"title":         "Test Site",
		"description":   "A test site",
		"url":           "https://example.com",
		"feeds_page":    feedsPage,
		"feed_defaults": defaults,
	}

	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	title := "Post"
	m.Cache().Set("feed_configs", []models.FeedConfig{
		{
			Slug:        "blog",
			Title:       "Blog",
			Description: "All posts",
			Formats:     models.FeedFormats{HTML: true, RSS: true, Atom: true, JSON: true, Sitemap: true},
			Posts: []*models.Post{{
				Slug:      "post",
				Href:      "/post/",
				Title:     &title,
				Published: true,
				Date:      &now,
			}},
		},
		{
			Slug:           "private-feed",
			Title:          "Private",
			IncludePrivate: true,
			Formats:        models.FeedFormats{RSS: true},
		},
	})

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(config.OutputDir, "feeds", "index.html"))
	if err != nil {
		t.Fatalf("ReadFile(feeds page) error = %v", err)
	}
	body := string(content)
	if !strings.Contains(body, "Blog") {
		t.Fatalf("feeds page should contain public feed title")
	}
	if !strings.Contains(body, `<meta name="robots" content="noindex,follow">`) {
		t.Fatalf("feeds page should include robots noindex meta")
	}
	if !strings.Contains(body, "/blog/archive/rss.xml") {
		t.Fatalf("feeds page should link to archive rss variant")
	}
	if strings.Contains(body, "Private") {
		t.Fatalf("feeds page should not contain private feeds by default")
	}
}

func TestFeedsListingPlugin_Write_IncludesOptedInPrivateFeed(t *testing.T) {
	plugin := NewFeedsListingPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.OutputDir = t.TempDir()
	feedsPage := models.NewFeedsPageConfig()
	feedsPage.ShowPrivateFeeds = []string{"private-feed"}
	defaults := models.NewFeedDefaults()
	config.Extra = map[string]interface{}{
		"title":         "Test Site",
		"description":   "A test site",
		"url":           "https://example.com",
		"feeds_page":    feedsPage,
		"feed_defaults": defaults,
	}

	m.Cache().Set("feed_configs", []models.FeedConfig{{
		Slug:           "private-feed",
		Title:          "Private",
		IncludePrivate: true,
		Formats:        models.FeedFormats{HTML: true, RSS: true},
	}})

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(config.OutputDir, "feeds", "index.html"))
	if err != nil {
		t.Fatalf("ReadFile(feeds page) error = %v", err)
	}
	body := string(content)
	if !strings.Contains(body, "Private") {
		t.Fatalf("feeds page should contain opted-in private feed")
	}
}

func TestFeedsListingPlugin_Write_TruncatesGeneratedFeedsOnMainPage(t *testing.T) {
	plugin := NewFeedsListingPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.OutputDir = t.TempDir()
	feedsPage := models.NewFeedsPageConfig()
	defaults := models.NewFeedDefaults()
	config.Extra = map[string]interface{}{
		"title":         "Test Site",
		"description":   "A test site",
		"url":           "https://example.com",
		"feeds_page":    feedsPage,
		"feed_defaults": defaults,
	}

	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	feedConfigs := make([]models.FeedConfig, 0, defaults.ItemsPerPage+16)
	for i := 0; i < defaults.ItemsPerPage+16; i++ {
		title := "Generated " + strconv.Itoa(i)
		slug := "generated-" + strconv.Itoa(i)
		feedConfigs = append(feedConfigs, models.FeedConfig{
			Slug:        slug,
			Title:       title,
			Description: "Generated feed",
			Formats:     models.FeedFormats{HTML: true, RSS: true},
			Posts: []*models.Post{{
				Slug:      slug + "-post",
				Href:      "/" + slug + "/",
				Title:     &title,
				Published: true,
				Date:      &now,
			}},
		})
	}
	m.Cache().Set("feed_configs", feedConfigs)

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	mainBodyBytes, err := os.ReadFile(filepath.Join(config.OutputDir, "feeds", "index.html"))
	if err != nil {
		t.Fatalf("ReadFile(main feeds page) error = %v", err)
	}
	mainBody := string(mainBodyBytes)
	if !strings.Contains(mainBody, `aria-label="Generated Feeds pages"`) {
		t.Fatalf("main feeds page should include generated feed pagination")
	}
	if !strings.Contains(mainBody, `href="/feeds/generated/page/2/"`) {
		t.Fatalf("main feeds page should link to generated feeds page 2")
	}
	if got := strings.Count(mainBody, `class="feed-row"`); got != defaults.ItemsPerPage {
		t.Fatalf("main feeds page should render %d preview rows, got %d", defaults.ItemsPerPage, got)
	}

	generatedBodyBytes, err := os.ReadFile(filepath.Join(config.OutputDir, "feeds", "generated", "index.html"))
	if err != nil {
		t.Fatalf("ReadFile(generated feeds page) error = %v", err)
	}
	generatedBody := string(generatedBodyBytes)
	if !strings.Contains(generatedBody, `aria-label="Pagination"`) {
		t.Fatalf("generated feeds page should include pagination controls")
	}
	if strings.Contains(generatedBody, "Generated 25") {
		t.Fatalf("first generated feeds page should be paginated")
	}
	generatedPage2BodyBytes, err := os.ReadFile(filepath.Join(config.OutputDir, "feeds", "generated", "page", "2", "index.html"))
	if err != nil {
		t.Fatalf("ReadFile(generated feeds page 2) error = %v", err)
	}
	if !strings.Contains(string(generatedPage2BodyBytes), "Generated 25") {
		t.Fatalf("generated feeds page 2 should include remaining generated feeds")
	}
}

func TestFeedsListingPlugin_Write_SkipsWhenCachedAndOutputsExist(t *testing.T) {
	plugin := NewFeedsListingPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.OutputDir = t.TempDir()
	feedsPage := models.NewFeedsPageConfig()
	defaults := models.NewFeedDefaults()
	config.Extra = map[string]interface{}{
		"title":         "Test Site",
		"description":   "A test site",
		"url":           "https://example.com",
		"feeds_page":    feedsPage,
		"feed_defaults": defaults,
	}
	m.Cache().Set("build_cache", buildcache.New(t.TempDir()))

	now := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	title := "Post"
	m.Cache().Set("feed_configs", []models.FeedConfig{{
		Slug:        "blog",
		Title:       "Blog",
		Description: "All posts",
		Formats:     models.FeedFormats{HTML: true, RSS: true},
		Posts: []*models.Post{{
			Slug:      "post",
			Href:      "/post/",
			Title:     &title,
			Published: true,
			Date:      &now,
		}},
	}})

	if err := plugin.Write(m); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}

	outputPath := filepath.Join(config.OutputDir, "feeds", "index.html")
	//nolint:gosec // test fixture simulates world-readable output file
	if err := os.WriteFile(outputPath, []byte("sentinel"), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", outputPath, err)
	}

	if err := plugin.Write(m); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}

	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) error = %v", outputPath, err)
	}
	if string(content) != "sentinel" {
		t.Fatalf("expected cached write skip to preserve sentinel content, got %q", string(content))
	}
}

func TestMonthlyPostBuckets_UsesSharedWindow(t *testing.T) {
	window := sparklineWindow{
		Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
	}
	title := "Post"
	posts := []*models.Post{
		{Title: &title, Published: true, Date: testDate(1, 15)},
		{Title: &title, Published: true, Date: testDate(4, 2)},
	}
	buckets, _ := monthlyPostBuckets(posts, window, false)
	if len(buckets) != 4 {
		t.Fatalf("expected 4 monthly buckets, got %d", len(buckets))
	}
	if buckets[0] != 1 || buckets[1] != 0 || buckets[2] != 0 || buckets[3] != 1 {
		t.Fatalf("unexpected bucket distribution: %#v", buckets)
	}
}

func TestBuildFeedSparklineData_ComputesHitboxOffsets(t *testing.T) {
	window := sparklineWindow{
		Start: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
	}
	title := "Post"
	posts := []*models.Post{
		{Title: &title, Published: true, Date: testDate(1, 15)},
		{Title: &title, Published: true, Date: testDate(3, 15)},
	}
	data := buildFeedSparklineData(posts, window, false)
	if len(data) != 3 {
		t.Fatalf("expected 3 sparkline points, got %d", len(data))
	}
	wantHitX := []float64{-2.6, 45.4, 93.4}
	for i, point := range data {
		if math.Abs(point.HitX-wantHitX[i]) > 0.000001 {
			t.Fatalf("point %d hitbox x = %v, want %v", i, point.HitX, wantHitX[i])
		}
	}
}

func testDate(month time.Month, day int) *time.Time {
	t := time.Date(2024, month, day, 0, 0, 0, 0, time.UTC)
	return &t
}
