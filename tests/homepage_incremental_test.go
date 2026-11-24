package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

func runHomepageIncrementalBuild(t *testing.T, contentDir, outputDir string) {
	t.Helper()
	m := lifecycle.NewManager()
	modelsConfig := models.NewConfig()
	modelsConfig.OutputDir = outputDir
	modelsConfig.Images = models.NewImagesConfig()
	m.SetConfig(&lifecycle.Config{
		ContentDir:   contentDir,
		OutputDir:    outputDir,
		GlobPatterns: []string{"pages/**/*.md"},
		Extra: map[string]interface{}{
			"models_config": modelsConfig,
			"title":         "Homepage regression fixture",
			"url":           "https://example.com",
		},
	})
	publish := plugins.NewPublishHTMLPlugin()
	images := plugins.NewImageLibraryPlugin()
	for _, plugin := range []lifecycle.Plugin{
		plugins.NewBuildCachePlugin(),
		plugins.NewGlobPlugin(),
		plugins.NewLoadPlugin(),
		plugins.NewRenderMarkdownPlugin(),
		plugins.NewTemplatesPlugin(),
	} {
		m.RegisterPlugin(plugin)
	}
	if err := m.RunTo(lifecycle.StageCollect); err != nil {
		t.Fatalf("run content stages: %v", err)
	}
	if err := publish.Write(m); err != nil {
		t.Fatalf("publish HTML: %v", err)
	}
	if err := images.Write(m); err != nil {
		t.Fatalf("publish image library: %v", err)
	}
	if err := m.RunTo(lifecycle.StageCleanup); err != nil {
		t.Fatalf("run cleanup: %v", err)
	}
}

func TestIncrementalBuildWithImageLibraryPreservesHomepage(t *testing.T) {
	root := t.TempDir()
	contentDir := filepath.Join(root, "content")
	outputDir := filepath.Join(root, "output")
	if err := os.MkdirAll(filepath.Join(contentDir, "pages", "post"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "pages", "post", "index.md"), []byte(`---
title: Homepage
slug: ""
published: true
---

![Homepage image](https://example.com/home.webp)
`), 0o600); err != nil {
		t.Fatal(err)
	}
	firstPost := filepath.Join(contentDir, "pages", "first.md")
	firstPostContent := `---
title: First post
published: true
---

First post.
`
	if err := os.WriteFile(firstPost, []byte(firstPostContent), 0o600); err != nil {
		t.Fatal(err)
	}

	runHomepageIncrementalBuild(t, contentDir, outputDir)
	if _, err := os.Stat(filepath.Join(outputDir, "index.html")); err != nil {
		t.Fatalf("first build did not write homepage: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "images", "index.html")); err != nil {
		t.Fatalf("first build did not write image library: %v", err)
	}

	if err := os.WriteFile(firstPost, []byte(firstPostContent+"\nChanged."), 0o600); err != nil {
		t.Fatal(err)
	}
	runHomepageIncrementalBuild(t, contentDir, outputDir)
	if _, err := os.Stat(filepath.Join(outputDir, "index.html")); err != nil {
		t.Fatalf("incremental build removed homepage: %v", err)
	}
}
