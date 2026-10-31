package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSitemapPlugin_Name(t *testing.T) {
	plugin := NewSitemapPlugin()
	if plugin.Name() != "sitemap" {
		t.Errorf("Name() = %q, want %q", plugin.Name(), "sitemap")
	}
}

func TestSitemapPlugin_Priority(t *testing.T) {
	plugin := NewSitemapPlugin()

	if got := plugin.Priority(lifecycle.StageWrite); got != lifecycle.PriorityLate {
		t.Errorf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityLate)
	}

	if got := plugin.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestSitemapPlugin_ExcludesPrivatePosts(t *testing.T) {
	plugin := NewSitemapPlugin()
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{
			Slug:      "public-post",
			Href:      "/public-post/",
			Published: true,
			Draft:     false,
			Skip:      false,
			Private:   false,
			Date:      &date,
		},
		{
			Slug:      "private-post",
			Href:      "/private-post/",
			Published: true,
			Draft:     false,
			Skip:      false,
			Private:   true,
			Date:      &date,
		},
		{
			Slug:      "another-public",
			Href:      "/another-public/",
			Published: true,
			Draft:     false,
			Skip:      false,
			Private:   false,
			Date:      &date,
		},
	})

	siteURL := "https://example.com"
	sitemap := plugin.buildSitemap(m, siteURL)

	// Count post URLs (excluding home page)
	postURLCount := 0
	for _, url := range sitemap.URLs {
		if url.Loc == siteURL+"/" {
			continue // Skip home page
		}
		postURLCount++
		// Verify private post URL is not present
		if url.Loc == siteURL+"/private-post/" {
			t.Error("private post should not appear in sitemap")
		}
	}

	// Should have 2 public posts (no feeds configured)
	if postURLCount != 2 {
		t.Errorf("expected 2 post URLs in sitemap, got %d", postURLCount)
	}
}

func TestSitemapPlugin_ExcludesDraftAndSkipPosts(t *testing.T) {
	plugin := NewSitemapPlugin()
	m := lifecycle.NewManager()

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m.SetPosts([]*models.Post{
		{
			Slug:      "published",
			Href:      "/published/",
			Published: true,
			Draft:     false,
			Skip:      false,
			Private:   false,
			Date:      &date,
		},
		{
			Slug:      "draft",
			Href:      "/draft/",
			Published: true,
			Draft:     true,
			Skip:      false,
			Private:   false,
			Date:      &date,
		},
		{
			Slug:      "skipped",
			Href:      "/skipped/",
			Published: true,
			Draft:     false,
			Skip:      true,
			Private:   false,
			Date:      &date,
		},
		{
			Slug:      "unpublished",
			Href:      "/unpublished/",
			Published: false,
			Draft:     false,
			Skip:      false,
			Private:   false,
			Date:      &date,
		},
	})

	siteURL := "https://example.com"
	sitemap := plugin.buildSitemap(m, siteURL)

	// Count post URLs (excluding home page)
	postURLCount := 0
	for _, url := range sitemap.URLs {
		if url.Loc == siteURL+"/" {
			continue
		}
		postURLCount++
	}

	// Should only have 1 published, non-draft, non-skip, non-private post
	if postURLCount != 1 {
		t.Errorf("expected 1 post URL in sitemap, got %d", postURLCount)
	}
}

func TestSitemapPlugin_Write_GeneratesSitemapIndex(t *testing.T) {
	plugin := NewSitemapPlugin()
	m := lifecycle.NewManager()
	config := m.Config()
	config.OutputDir = t.TempDir()
	config.Extra = map[string]interface{}{
		"url":           "https://example.com",
		"feed_defaults": models.NewFeedDefaults(),
		"feeds_page":    models.NewFeedsPageConfig(),
	}

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{{
		Slug:      "published",
		Href:      "/published/",
		Published: true,
		Date:      &date,
	}})
	feedConfigs := []models.FeedConfig{{
		Slug:    "blog",
		Formats: models.FeedFormats{HTML: true, Sitemap: true},
	}}
	defaults := models.NewFeedDefaults()
	for i := 0; i < defaults.ItemsPerPage+1; i++ {
		feedConfigs = append(feedConfigs, models.FeedConfig{
			Slug:    fmt.Sprintf("generated-%d", i),
			Formats: models.FeedFormats{HTML: true},
		})
	}
	m.Cache().Set("feed_configs", feedConfigs)

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	indexContent, err := os.ReadFile(filepath.Join(config.OutputDir, "sitemap.xml"))
	if err != nil {
		t.Fatalf("ReadFile(sitemap.xml) error = %v", err)
	}
	if !strings.Contains(string(indexContent), "<sitemapindex") {
		t.Fatalf("root sitemap should be a sitemap index")
	}
	if !strings.Contains(string(indexContent), "https://example.com/blog/sitemap.xml") {
		t.Fatalf("root sitemap index should reference feed sitemap")
	}
	if !strings.Contains(string(indexContent), "<lastmod>2024-01-15</lastmod>") {
		t.Fatalf("root sitemap index should include sitemap lastmod")
	}

	pagesContent, err := os.ReadFile(filepath.Join(config.OutputDir, "sitemap-pages.xml"))
	if err != nil {
		t.Fatalf("ReadFile(sitemap-pages.xml) error = %v", err)
	}
	if !strings.Contains(string(pagesContent), "https://example.com/feeds/") {
		t.Fatalf("pages sitemap should include feeds listing page")
	}
	if !strings.Contains(string(pagesContent), "https://example.com/feeds/generated/") {
		t.Fatalf("pages sitemap should include generated feeds page")
	}
}

// Ensure SitemapPlugin implements the required interfaces.
func TestSitemapPlugin_ImplementsInterfaces(_ *testing.T) {
	var _ lifecycle.Plugin = (*SitemapPlugin)(nil)
	var _ lifecycle.WritePlugin = (*SitemapPlugin)(nil)
	var _ lifecycle.PriorityPlugin = (*SitemapPlugin)(nil)
}

func TestPagesSitemapExcludesNoindexAndDuplicateLandingPages(t *testing.T) {
	m := lifecycle.NewManager()
	m.Config().Extra = map[string]interface{}{
		"feeds_page": models.FeedsPageConfig{SlugPrefix: "feeds", Robots: "noindex,follow"},
	}
	m.SetPosts([]*models.Post{
		{Slug: "blog", Href: "/blog/", Published: true},
		{Slug: "noindex", Href: "/noindex/", Published: true, Extra: map[string]interface{}{"robots": "noindex, follow"}},
		{Slug: "public", Href: "/public/", Published: true},
	})
	m.Cache().Set("feed_configs", []models.FeedConfig{
		{Slug: "blog", Formats: models.FeedFormats{HTML: true, Sitemap: true}},
		{Slug: "draft", Robots: "noindex,follow", Formats: models.FeedFormats{HTML: true, Sitemap: true}},
		{Slug: "all", Formats: models.FeedFormats{HTML: true, Sitemap: false}},
		{Slug: "public-feed", Formats: models.FeedFormats{HTML: true, Sitemap: true}},
	})
	got := NewSitemapPlugin().buildPagesSitemap(m, "https://example.com")
	seen := make(map[string]int)
	for _, entry := range got.URLs {
		seen[entry.Loc]++
	}
	for _, forbidden := range []string{"/noindex/", "/draft/", "/all/", "/feeds/"} {
		if seen["https://example.com"+forbidden] != 0 {
			t.Errorf("included excluded landing: %s", forbidden)
		}
	}
	for _, included := range []string{"/blog/", "/public/", "/public-feed/"} {
		if seen["https://example.com"+included] != 1 {
			t.Errorf("wanted one %s, got %d", included, seen["https://example.com"+included])
		}
	}
}

func TestSitemapNoindexDirectiveParsing(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  bool
	}{
		{"noindex,follow", true},
		{"NOINDEX", true},
		{"index,follow", false},
		{"noimageindex", false},
		{"follow noindex", true},
		{"", false},
	} {
		if got := sitemapNoindex(tc.input); got != tc.want {
			t.Errorf("sitemapNoindex(%q)=%v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestSitemapIndexOmitsNoindexAndDuplicateFeedSitemaps(t *testing.T) {
	m := lifecycle.NewManager()
	m.Config().Extra = map[string]interface{}{"feeds_page": models.NewFeedsPageConfig()}
	m.Cache().Set("feed_configs", []models.FeedConfig{
		{Slug: "blog", Formats: models.FeedFormats{HTML: true, Sitemap: true}},
		{Slug: "draft", Robots: "NOINDEX,follow", Formats: models.FeedFormats{HTML: true, Sitemap: true}},
		{Slug: "blog", Formats: models.FeedFormats{HTML: true, Sitemap: true}},
		{Slug: "disabled", Formats: models.FeedFormats{HTML: true}},
		{Slug: "private", IncludePrivate: true, Formats: models.FeedFormats{HTML: true, Sitemap: true}},
		{Slug: "sitemap-only", Formats: models.FeedFormats{Sitemap: true}},
	})

	index := NewSitemapPlugin().buildSitemapIndex(m, "https://example.com")
	want := map[string]int{
		"https://example.com/sitemap-pages.xml":        1,
		"https://example.com/blog/sitemap.xml":         1,
		"https://example.com/sitemap-only/sitemap.xml": 1,
	}
	got := make(map[string]int)
	for _, entry := range index.Sitemaps {
		got[entry.Loc]++
	}
	if len(got) != len(want) {
		t.Fatalf("sitemap index entries = %v, want %v", got, want)
	}
	for loc, count := range want {
		if got[loc] != count {
			t.Errorf("sitemap index %q count = %d, want %d", loc, got[loc], count)
		}
	}

	// Exercise the production writer, not only its in-memory index builder.
	m.Config().OutputDir = t.TempDir()
	m.Config().Extra["url"] = "https://example.com"
	if err := NewSitemapPlugin().Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(m.Config().OutputDir, "sitemap.xml"))
	if err != nil {
		t.Fatalf("ReadFile(sitemap.xml) error = %v", err)
	}
	for loc, count := range want {
		if actual := strings.Count(string(data), "<loc>"+loc+"</loc>"); actual != count {
			t.Errorf("emitted sitemap.xml %q count = %d, want %d: %s", loc, actual, count, data)
		}
	}
	for _, excluded := range []string{"draft", "disabled", "private"} {
		loc := "https://example.com/" + excluded + "/sitemap.xml"
		if strings.Contains(string(data), "<loc>"+loc+"</loc>") {
			t.Errorf("emitted sitemap.xml advertises excluded feed %q: %s", loc, data)
		}
	}
}

func TestSitemapIndexLastModExcludesNonEmittedPosts(t *testing.T) {
	plugin := NewSitemapPlugin()
	m := lifecycle.NewManager()
	m.Config().Extra = map[string]interface{}{"feeds_page": models.NewFeedsPageConfig()}
	publishedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	futureDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	public := &models.Post{Slug: "public", Href: "/public/", Published: true, Date: &publishedDate}
	undated := &models.Post{Slug: "undated", Href: "/undated/", Published: true}
	draft := &models.Post{Slug: "draft", Href: "/draft/", Published: true, Draft: true, Date: &futureDate}
	scheduled := &models.Post{Slug: "scheduled", Href: "/scheduled/", Date: &futureDate}
	private := &models.Post{Slug: "private", Href: "/private/", Published: true, Private: true, Date: &futureDate}
	skipped := &models.Post{Slug: "skipped", Href: "/skipped/", Published: true, Skip: true, Date: &futureDate}
	noindex := &models.Post{
		Slug: "noindex", Href: "/noindex/", Published: true, Date: &futureDate,
		Extra: map[string]interface{}{"robots": "noindex,follow"},
	}
	m.SetPosts([]*models.Post{nil, public, undated, draft, scheduled, private, skipped, noindex})
	m.Cache().Set("feed_configs", []models.FeedConfig{
		{
			Slug: "blog", Formats: models.FeedFormats{HTML: true, Sitemap: true},
			Posts: []*models.Post{nil, public, undated, draft, scheduled, private, skipped},
		},
		{
			Slug: "empty", Formats: models.FeedFormats{HTML: true, Sitemap: true},
			Posts: []*models.Post{draft, scheduled, private, skipped},
		},
		{
			Slug: "disabled", Formats: models.FeedFormats{HTML: true},
			Posts: []*models.Post{public},
		},
	})

	pages := plugin.buildPagesSitemap(m, "https://example.com")
	for _, entry := range pages.URLs {
		if entry.LastMod == "2026-12-31" {
			t.Errorf("excluded future post affected pages sitemap: %+v", entry)
		}
	}
	index := plugin.buildSitemapIndex(m, "https://example.com")
	want := map[string]string{
		"https://example.com/sitemap-pages.xml": "2024-01-15",
		"https://example.com/blog/sitemap.xml":  "2024-01-15",
		"https://example.com/empty/sitemap.xml": "",
	}
	if len(index.Sitemaps) != len(want) {
		t.Fatalf("index has %d entries, want %d: %+v", len(index.Sitemaps), len(want), index.Sitemaps)
	}
	for _, entry := range index.Sitemaps {
		expected, ok := want[entry.Loc]
		if !ok {
			t.Errorf("unexpected sitemap entry %q", entry.Loc)
			continue
		}
		if entry.LastMod != expected {
			t.Errorf("%s lastmod = %q, want %q", entry.Loc, entry.LastMod, expected)
		}
	}
}

func TestSitemapIndexLastModIncludesGeneratedLandingDate(t *testing.T) {
	plugin := NewSitemapPlugin()
	m := lifecycle.NewManager()
	m.Config().Extra = map[string]interface{}{"feeds_page": models.NewFeedsPageConfig()}
	m.SetPosts([]*models.Post{{Slug: "undated", Href: "/undated/", Published: true}})
	feedDate := time.Date(2025, 7, 4, 0, 0, 0, 0, time.UTC)
	m.Cache().Set("feed_configs", []models.FeedConfig{{
		Slug:    "blog",
		Formats: models.FeedFormats{HTML: true, Sitemap: true},
		Posts:   []*models.Post{{Slug: "feed-only", Href: "/feed-only/", Published: true, Date: &feedDate}},
	}})

	index := plugin.buildSitemapIndex(m, "https://example.com")
	if len(index.Sitemaps) != 2 {
		t.Fatalf("index entries = %+v, want pages and blog", index.Sitemaps)
	}
	for _, entry := range index.Sitemaps {
		if entry.LastMod != "2025-07-04" {
			t.Errorf("%s lastmod = %q, want generated landing date", entry.Loc, entry.LastMod)
		}
	}
}

func TestPublishFeedSitemapEmittedXMLSkipsExcludedPosts(t *testing.T) {
	outputDir := t.TempDir()
	publishedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	futureDate := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	fc := &models.FeedConfig{
		Slug: "blog",
		Posts: []*models.Post{
			nil,
			{Slug: "public", Href: "/public/", Published: true, Date: &publishedDate},
			{Slug: "draft", Href: "/draft/", Published: true, Draft: true, Date: &futureDate},
			{Slug: "private", Href: "/private/", Published: true, Private: true, Date: &futureDate},
			{Slug: "noindex", Href: "/noindex/", Published: true, Extra: map[string]interface{}{"robots": "noindex,follow"}},
		},
	}
	config := &lifecycle.Config{Extra: map[string]interface{}{"url": "https://example.com"}}
	if err := NewPublishFeedsPlugin().publishSitemap(fc, config, outputDir); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(outputDir, "sitemap.xml"))
	if err != nil {
		t.Fatal(err)
	}
	xml := string(data)
	if !strings.Contains(xml, "<loc>https://example.com/public/</loc>") {
		t.Errorf("published post missing from sitemap: %s", xml)
	}
	for _, excluded := range []string{"/draft/", "/private/", "/noindex/"} {
		if strings.Contains(xml, "<loc>https://example.com"+excluded+"</loc>") {
			t.Errorf("excluded post %s present in sitemap: %s", excluded, xml)
		}
	}
}
