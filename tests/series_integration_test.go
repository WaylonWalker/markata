package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

func TestIntegration_FullBuild_GeneratesSeriesFeedsFromFrontmatter(t *testing.T) {
	site := newTestSite(t)
	site.addPost("posts/part-1.md", `---
title: Part 1
slug: part-1
published: true
date: 2024-01-01
series: "go tutorial"
---

Part 1`)
	site.addPost("posts/part-2.md", `---
title: Part 2
slug: part-2
published: true
date: 2024-01-02
series: "go tutorial"
---

Part 2`)
	site.addConfig("" +
		"[markata-go]\n" +
		"url = \"https://example.com\"\n" +
		"title = \"Test Site\"\n" +
		"output_dir = \"" + filepath.ToSlash(site.outputDir) + "\"\n\n" +
		"[markata-go.glob]\n" +
		"patterns = [\"posts/**/*.md\"]\n")

	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir(site.dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", site.dir, err)
	}
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			t.Fatalf("failed to restore cwd: %v", err)
		}
	}()

	cfg, err := config.Load(filepath.Join(site.dir, "markata-go.toml"))
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	m := lifecycle.NewManager()
	lcConfig := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    cfg.OutputDir,
		GlobPatterns: cfg.GlobConfig.Patterns,
		Extra:        make(map[string]interface{}),
	}
	lcConfig.Extra["url"] = cfg.URL
	lcConfig.Extra["title"] = cfg.Title
	lcConfig.Extra["feeds"] = cfg.Feeds
	lcConfig.Extra["feed_defaults"] = cfg.FeedDefaults
	lcConfig.Extra["theme"] = cfg.Theme
	lcConfig.Extra["models_config"] = cfg
	for key, value := range cfg.Extra {
		lcConfig.Extra[key] = value
	}
	m.SetConfig(lcConfig)
	m.RegisterPlugins(plugins.DefaultPlugins()...)

	if err := m.Run(); err != nil {
		t.Fatalf("full build failed: %v", err)
	}

	if !site.fileExists("series/go-tutorial/index.html") {
		t.Fatal("expected full build to generate series/go-tutorial/index.html")
	}
	if !site.fileExists("series/go-tutorial/rss.xml") {
		t.Fatal("expected full build to generate series/go-tutorial/rss.xml")
	}

	seriesFeed := findFeedConfigBySlug(m, "series/go-tutorial")
	if seriesFeed == nil {
		t.Fatal("expected generated series feed config in full build")
	}
	if seriesFeed.Type != models.FeedTypeSeries {
		t.Fatalf("series feed type = %q, want %q", seriesFeed.Type, models.FeedTypeSeries)
	}
	if len(seriesFeed.Posts) != 2 {
		t.Fatalf("series feed posts = %d, want 2", len(seriesFeed.Posts))
	}

	part1 := findPostBySlug(m.Posts(), "part-1")
	part2 := findPostBySlug(m.Posts(), "part-2")
	if part1 == nil || part2 == nil {
		t.Fatal("expected built posts to be available")
	}
	if part1.PrevNextFeed != "series/go-tutorial" {
		t.Fatalf("part1 PrevNextFeed = %q, want %q", part1.PrevNextFeed, "series/go-tutorial")
	}
	if part1.Next == nil || part1.Next.Slug != "part-2" {
		t.Fatalf("part1 next post = %#v, want part-2", part1.Next)
	}
	if part2.Prev == nil || part2.Prev.Slug != "part-1" {
		t.Fatalf("part2 prev post = %#v, want part-1", part2.Prev)
	}
}

func findFeedConfigBySlug(m *lifecycle.Manager, slug string) *models.FeedConfig {
	if m == nil {
		return nil
	}
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		return nil
	}
	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		return nil
	}
	for i := range feedConfigs {
		if feedConfigs[i].Slug == slug {
			return &feedConfigs[i]
		}
	}
	return nil
}

func findPostBySlug(posts []*models.Post, slug string) *models.Post {
	for _, post := range posts {
		if post.Slug == slug {
			return post
		}
	}
	return nil
}
