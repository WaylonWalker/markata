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
