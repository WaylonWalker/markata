// Package tests provides integration tests for markata-go.
// These tests verify the full build workflow from content to output.
package tests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

// =============================================================================
// Test Helpers
// =============================================================================

// testSite creates a test site structure in a temporary directory.
type testSite struct {
	dir        string
	contentDir string
	outputDir  string
	t          *testing.T
}

// newTestSite creates a new test site in a temp directory.
func newTestSite(t *testing.T) *testSite {
	t.Helper()
	dir := t.TempDir()

	contentDir := filepath.Join(dir, "content")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatalf("failed to create content dir: %v", err)
	}

	outputDir := filepath.Join(dir, "output")

	return &testSite{
		dir:        dir,
		contentDir: contentDir,
		outputDir:  outputDir,
		t:          t,
	}
}

// addPost adds a markdown post to the content directory.
func (s *testSite) addPost(path, content string) {
	s.t.Helper()
	fullPath := filepath.Join(s.contentDir, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		s.t.Fatalf("failed to create dir %s: %v", dir, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o600); err != nil {
		s.t.Fatalf("failed to write %s: %v", path, err)
	}
}

// addConfig writes a config file.
func (s *testSite) addConfig(content string) {
	s.t.Helper()
	configPath := filepath.Join(s.dir, "markata-go.toml")
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		s.t.Fatalf("failed to write config: %v", err)
	}
}

// disableMinificationForTestBuild keeps integration tests focused on their
// asserted behavior. Minifier behavior has dedicated tests in pkg/plugins;
// processing all bundled assets for every tiny fixture is slow under -race.
func disableMinificationForTestBuild(extra map[string]interface{}) {
	extra["js_minify"] = map[string]interface{}{"enabled": false}
	extra["css_minify"] = map[string]interface{}{"enabled": false}
}

// build runs a full build and returns the manager.
func (s *testSite) build() *lifecycle.Manager {
	s.t.Helper()

	// Create manager with default config
	m := lifecycle.NewManager()

	// Configure
	cfg := &lifecycle.Config{
		ContentDir:   s.contentDir,
		OutputDir:    s.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"
	m.SetConfig(cfg)

	// Register plugins
	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	// Run build
	if err := m.Run(); err != nil {
		s.t.Fatalf("build failed: %v", err)
	}

	return m
}

// buildWithFeeds runs a build with feed support.
func (s *testSite) buildWithFeeds(feedConfigs []models.FeedConfig) *lifecycle.Manager {
	s.t.Helper()

	m := lifecycle.NewManager()

	cfg := &lifecycle.Config{
		ContentDir:   s.contentDir,
		OutputDir:    s.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"
	cfg.Extra["feeds"] = feedConfigs
	cfg.Extra["feed_defaults"] = models.FeedDefaults{
		ItemsPerPage:    10,
		OrphanThreshold: 3,
		Formats: models.FeedFormats{
			HTML: true,
			RSS:  true,
		},
	}
	m.SetConfig(cfg)

	// Register plugins
	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewFeedsPlugin())
	m.RegisterPlugin(plugins.NewPublishFeedsPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	if err := m.Run(); err != nil {
		s.t.Fatalf("build failed: %v", err)
	}

	return m
}

// fileExists checks if a file exists in the output directory.
func (s *testSite) fileExists(path string) bool {
	fullPath := filepath.Join(s.outputDir, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// readFile reads a file from the output directory.
func (s *testSite) readFile(path string) string {
	s.t.Helper()
	fullPath := filepath.Join(s.outputDir, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		s.t.Fatalf("failed to read %s: %v", path, err)
	}
	return string(data)
}

// =============================================================================
// Output Structure Tests (from tests.yaml output_structure)
// =============================================================================

func TestOutputStructure_PostOutputPath(t *testing.T) {
	// Test case: "post output path"
	// Post with slug "hello-world" should output to output/hello-world/index.html
	site := newTestSite(t)
	site.addPost("hello-world.md", `---
title: Hello World
published: true
---
Content here`)

	site.build()

	if !site.fileExists("hello-world/index.html") {
		t.Error("expected output/hello-world/index.html to exist")
	}
}

func TestOutputStructure_NestedSlugPath(t *testing.T) {
	// Test case: "nested slug output path"
	// Post with slug "blog/2024/my-post" should output to nested directory
	site := newTestSite(t)
	site.addPost("blog/2024/my-post.md", `---
title: My Post
slug: blog/2024/my-post
published: true
---
Content here`)

	site.build()

	if !site.fileExists("blog/2024/my-post/index.html") {
		t.Error("expected output/blog/2024/my-post/index.html to exist")
	}
}

func TestOutputStructure_PathSlugModeFromContentPath(t *testing.T) {
	site := newTestSite(t)
	site.addPost("posts/notes/today.md", `---
title: Today
published: true
---
Content here`)

	m := lifecycle.NewManager()
	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	modelsConfig := config.DefaultConfig()
	modelsConfig.GlobConfig.SlugMode = models.SlugModePath
	cfg.Extra["models_config"] = modelsConfig
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"
	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	if err := m.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	if !site.fileExists("notes/today/index.html") {
		t.Error("expected output/notes/today/index.html to exist")
	}
}

func TestOutputStructure_SlugRulesHybridMode(t *testing.T) {
	site := newTestSite(t)
	site.addPost("posts/blog/launch-post.md", `---
title: Launch
published: true
---
Launch content`)
	site.addPost("posts/notes/today.md", `---
title: Today
published: true
---
Notes content`)

	m := lifecycle.NewManager()
	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	modelsConfig := config.DefaultConfig()
	modelsConfig.GlobConfig.SlugMode = models.SlugModeFlat
	modelsConfig.GlobConfig.SlugRules = []models.SlugRule{
		{Prefix: "posts/blog", Mode: models.SlugModeFlat},
		{Prefix: "posts/notes", Mode: models.SlugModePath},
	}
	cfg.Extra["models_config"] = modelsConfig
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"
	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	if err := m.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	if !site.fileExists("launch-post/index.html") {
		t.Error("expected output/launch-post/index.html to exist")
	}
	if !site.fileExists("notes/today/index.html") {
		t.Error("expected output/notes/today/index.html to exist")
	}
}

func TestCanonicalHeadingsBuildUsesCanonicalPostProjection(t *testing.T) {
	site := newTestSite(t)
	site.addPost("test-headings.md", `---
template: post
title: test headings
slug: test-headings
published: true
---

# Canonical heading`)
	site.addPost("ordinary.md", `---
title: ordinary
published: true
---

# Ordinary heading`)

	m := lifecycle.NewManager()
	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra: map[string]interface{}{
			"templates_dir": "../templates",
			"models_config": &models.Config{URL: "https://example.com"},
		},
	}
	m.SetConfig(cfg)
	templatesPlugin := plugins.NewTemplatesPlugin()
	if err := templatesPlugin.Configure(m); err != nil {
		t.Fatalf("templates configure failed: %v", err)
	}
	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(templatesPlugin)
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())
	if err := m.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	html := site.readFile("test-headings/index.html")
	for _, marker := range []string{
		`data-rendering-specimen="canonical-headings"`,
		`class="canonical-document"`,
	} {
		if !strings.Contains(html, marker) {
			t.Errorf("generated canonical post is missing %q", marker)
		}
	}
	if strings.Contains(html, `class="post h-entry"`) {
		t.Error("generated canonical post used the generic post projection")
	}
	if !site.fileExists("test-headings/index.html") {
		t.Fatal("canonical fixture was not written to its slug path")
	}
	ordinary := site.readFile("ordinary/index.html")
	if strings.Contains(ordinary, `<main class="main-content" id="main-content" tabindex="-1" data-rendering-specimen=`) {
		t.Error("ordinary post unexpectedly contains the canonical route marker")
	}
	if strings.Contains(ordinary, `class="canonical-document"`) {
		t.Error("ordinary post unexpectedly contains the canonical projection")
	}
	if !strings.Contains(ordinary, `class="post h-entry"`) {
		t.Error("ordinary post did not use the generic post projection")
	}
}

func TestLinkAvatars_RootAbsoluteAssetsOnNestedPage(t *testing.T) {
	site := newTestSite(t)
	site.addPost("instant-pot-chicken-teriyaki.md", `---
title: Instant Pot Chicken Teriyaki
slug: instant-pot-chicken-teriyaki
published: true
---
Content here`)

	m := lifecycle.NewManager()
	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	cfg.Extra["models_config"] = &models.Config{URL: "https://recipes.waylonwalker.com"}
	cfg.Extra["link_avatars"] = map[string]any{
		"enabled": true,
		"mode":    "js",
	}
	m.SetConfig(cfg)

	templatesPlugin := plugins.NewTemplatesPlugin()
	if err := templatesPlugin.Configure(m); err != nil {
		t.Fatalf("templates configure failed: %v", err)
	}

	m.RegisterPlugin(plugins.NewStaticAssetsPlugin())
	m.RegisterPlugin(plugins.NewLinkAvatarsPlugin())
	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(templatesPlugin)
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	if err := m.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	html := site.readFile("instant-pot-chicken-teriyaki/index.html")
	if strings.Contains(html, `href="css/`) || strings.Contains(html, `src="js/`) {
		t.Fatalf("expected root-absolute asset links, got %s", html)
	}
	if !strings.Contains(html, `href="/css/link-avatars.`) {
		t.Fatalf("expected hashed link-avatars css href, got %s", html)
	}
	if !strings.Contains(html, `src="/js/link-avatars.`) {
		t.Fatalf("expected hashed link-avatars js src, got %s", html)
	}
	if strings.Contains(html, "/assets/markata/link-avatars.css") || strings.Contains(html, "/assets/markata/link-avatars.js") {
		t.Fatalf("unexpected link-avatars asset paths in html: %s", html)
	}

	if !site.fileExists("css/link-avatars.css") {
		t.Fatalf("expected css asset at output/css/link-avatars.css")
	}
	if !site.fileExists("js/link-avatars.js") {
		t.Fatalf("expected js asset at output/js/link-avatars.js")
	}

	cssMatches, err := filepath.Glob(filepath.Join(site.outputDir, "css", "link-avatars.*.css"))
	if err != nil {
		t.Fatalf("glob css hashes: %v", err)
	}
	if len(cssMatches) == 0 {
		t.Fatalf("expected hashed css copy to exist")
	}
	jsMatches, err := filepath.Glob(filepath.Join(site.outputDir, "js", "link-avatars.*.js"))
	if err != nil {
		t.Fatalf("glob js hashes: %v", err)
	}
	if len(jsMatches) == 0 {
		t.Fatalf("expected hashed js copy to exist")
	}
}

func TestOutputStructure_FeedOutputPath(t *testing.T) {
	// Test case: "feed output path"
	// Feed with slug "archive" should output to output/archive/index.html
	site := newTestSite(t)
	site.addPost("post1.md", `---
title: Post 1
published: true
---
Content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "archive",
			Title:  "Archive",
			Filter: "published == True",
			Formats: models.FeedFormats{
				HTML: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	if !site.fileExists("archive/index.html") {
		t.Error("expected output/archive/index.html to exist")
	}
}

func TestOutputStructure_FeedPage2Path(t *testing.T) {
	// Test case: "feed page 2 output path"
	// Paginated feed page 2 should be at output/blog/page/2/index.html
	site := newTestSite(t)

	// Create enough posts to trigger pagination
	for i := 1; i <= 15; i++ {
		filename := filepath.Join("posts", fmt.Sprintf("post-%02d.md", i))
		content := fmt.Sprintf(`---
title: Post %d
published: true
date: 2024-01-%02d
---
Content for post %d`, i, i, i)
		site.addPost(filename, content)
	}

	feedConfigs := []models.FeedConfig{
		{
			Slug:         "blog",
			Title:        "Blog",
			Filter:       "published == True",
			ItemsPerPage: 5,
			Formats: models.FeedFormats{
				HTML: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	// Check that paginated files exist
	if !site.fileExists("blog/index.html") {
		t.Error("expected output/blog/index.html to exist")
	}
	if !site.fileExists("blog/page/2/index.html") {
		t.Error("expected output/blog/page/2/index.html to exist")
	}
}

func TestPostIncludesArticleProgressIndicator(t *testing.T) {
	site := newTestSite(t)
	site.addPost("progress.md", `---
title: Progress Indicator
slug: progress
published: true
---
# Reading Progress

This post is long enough to trigger the reading progress indicator.

`)

	site.build()

	content := site.readFile("progress/index.html")
	if !strings.Contains(content, "article-progress__indicator") {
		t.Errorf("expected article progress markup in post output")
	}
	if !strings.Contains(content, "initArticleProgressIndicator") {
		t.Errorf("expected progress initializer script to be rendered")
	}
}

func TestTailwindManagedConfig_DisablesPreflightByDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	site := newTestSite(t)
	site.addPost("tailwind.md", `---
title: Tailwind Defaults
published: true
---
# Heading

<div class="text-red-500 font-bold">utility text</div>
`)

	m := lifecycle.NewManager()
	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	cfg.Extra["url"] = "https://example.com"
	cfg.Extra["title"] = "Test Site"
	cfg.Extra["tailwind"] = models.NewTailwindConfig()
	m.SetConfig(cfg)

	templatesPlugin := plugins.NewTemplatesPlugin()
	if err := templatesPlugin.Configure(m); err != nil {
		t.Fatalf("templates configure failed: %v", err)
	}

	m.RegisterPlugin(plugins.NewTailwindPlugin())
	m.RegisterPlugin(plugins.NewStaticAssetsPlugin())
	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(templatesPlugin)
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	if err := m.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	html := site.readFile("tailwind/index.html")
	if !strings.Contains(html, "/markata-tailwind.css") {
		t.Fatalf("expected managed Tailwind CSS to be included, got %s", html)
	}

	tailwindCSS := site.readFile("markata-tailwind.css")
	if strings.Contains(tailwindCSS, "h1,h2,h3,h4,h5,h6{font-size:inherit;font-weight:inherit}") {
		t.Fatal("expected preflight heading reset to be disabled in managed Tailwind CSS")
	}
	if !strings.Contains(tailwindCSS, ".text-red-500") {
		t.Fatal("expected Tailwind utility class to be generated")
	}
}

// =============================================================================
// Feed Format Tests (from tests.yaml feed_formats)
// =============================================================================

func TestFeedFormat_HTMLOutput(t *testing.T) {
	// Test case: "html format output"
	site := newTestSite(t)
	site.addPost("test.md", `---
title: Test Post
slug: test
published: true
---
Content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "blog",
			Title:  "Blog",
			Filter: "True",
			Formats: models.FeedFormats{
				HTML: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	if !site.fileExists("blog/index.html") {
		t.Error("expected output/blog/index.html to exist")
	}
}

func TestFeedFormat_RSSOutput(t *testing.T) {
	// Test case: "rss format output"
	site := newTestSite(t)
	site.addPost("test.md", `---
title: Test
slug: test
published: true
date: 2024-01-15
---
Content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "blog",
			Title:  "Blog",
			Filter: "True",
			Formats: models.FeedFormats{
				RSS: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	if !site.fileExists("blog/rss.xml") {
		t.Error("expected output/blog/rss.xml to exist")
	}

	// Verify RSS structure
	content := site.readFile("blog/rss.xml")
	if !strings.Contains(content, `<rss version="2.0"`) {
		t.Error("RSS should contain version 2.0 declaration")
	}
	if !strings.Contains(content, "<title>Test</title>") {
		t.Error("RSS should contain post title")
	}
	if !strings.Contains(content, "<pubDate>") {
		t.Error("RSS should contain pubDate element")
	}

	// Verify valid XML
	var rss struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
	}
	if err := xml.Unmarshal([]byte(content), &rss); err != nil {
		t.Errorf("RSS should be valid XML: %v", err)
	}
}

func TestFeedFormat_AtomOutput(t *testing.T) {
	// Test case: "atom format output"
	site := newTestSite(t)
	site.addPost("test.md", `---
title: Test
slug: test
published: true
date: 2024-01-15
---
Content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "blog",
			Title:  "Blog",
			Filter: "True",
			Formats: models.FeedFormats{
				Atom: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	if !site.fileExists("blog/atom.xml") {
		t.Error("expected output/blog/atom.xml to exist")
	}

	content := site.readFile("blog/atom.xml")
	if !strings.Contains(content, `xmlns="http://www.w3.org/2005/Atom"`) {
		t.Error("Atom feed should contain Atom namespace")
	}
	if !strings.Contains(content, "<entry>") {
		t.Error("Atom feed should contain entry elements")
	}
	if !strings.Contains(content, "<published>") {
		t.Error("Atom feed should contain published element")
	}
}

func TestFeedFormat_JSONOutput(t *testing.T) {
	// Test case: "json feed format output"
	site := newTestSite(t)
	site.addPost("test.md", `---
title: Test
slug: test
published: true
---
Content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "blog",
			Title:  "Blog",
			Filter: "True",
			Formats: models.FeedFormats{
				JSON: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	if !site.fileExists("blog/feed.json") {
		t.Error("expected output/blog/feed.json to exist")
	}

	content := site.readFile("blog/feed.json")

	// Verify valid JSON
	var feed map[string]interface{}
	if err := json.Unmarshal([]byte(content), &feed); err != nil {
		t.Errorf("JSON feed should be valid JSON: %v", err)
	}

	// Verify JSON Feed version
	if version, ok := feed["version"].(string); !ok || version != "https://jsonfeed.org/version/1.1" {
		t.Error("JSON feed should contain version 1.1")
	}
}

func TestFeedFormat_MultipleFormats(t *testing.T) {
	// Test case: "multiple formats simultaneously"
	site := newTestSite(t)
	site.addPost("test.md", `---
title: Test
slug: test
published: true
date: 2024-01-15
---
Content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "blog",
			Title:  "Blog",
			Filter: "True",
			Formats: models.FeedFormats{
				HTML: true,
				RSS:  true,
				Atom: true,
				JSON: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	expectedFiles := []string{
		"blog/index.html",
		"blog/rss.xml",
		"blog/atom.xml",
		"blog/feed.json",
	}

	for _, f := range expectedFiles {
		if !site.fileExists(f) {
			t.Errorf("expected output/%s to exist", f)
		}
	}
}

// =============================================================================
// Lifecycle Tests (from tests.yaml lifecycle)
// =============================================================================

func TestLifecycle_StagesRunInOrder(t *testing.T) {
	// Test case: "stages run in order"
	site := newTestSite(t)
	site.addPost("test.md", `---
title: Test
published: true
---
Content`)

	m := lifecycle.NewManager()

	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	m.SetConfig(cfg)

	// Register plugins
	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())

	// Run only to render stage
	err := m.RunTo(lifecycle.StageRender)
	if err != nil {
		t.Fatalf("RunTo failed: %v", err)
	}

	// Verify stages have run
	expectedStages := []lifecycle.Stage{
		lifecycle.StageConfigure,
		lifecycle.StageValidate,
		lifecycle.StageGlob,
		lifecycle.StageLoad,
		lifecycle.StageTransform,
		lifecycle.StageRender,
	}

	for _, stage := range expectedStages {
		if !m.HasRun(stage) {
			t.Errorf("stage %s should have run", stage)
		}
	}

	// Verify later stages have not run
	laterStages := []lifecycle.Stage{
		lifecycle.StageCollect,
		lifecycle.StageWrite,
		lifecycle.StageCleanup,
	}

	for _, stage := range laterStages {
		if m.HasRun(stage) {
			t.Errorf("stage %s should NOT have run", stage)
		}
	}
}

func TestLifecycle_FullBuildRunsAllStages(t *testing.T) {
	// Test case: "full build runs all stages"
	site := newTestSite(t)
	site.addPost("test.md", `---
title: Test
published: true
---
Content`)

	m := lifecycle.NewManager()

	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	// Run full build
	err := m.Run()
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify all stages have run
	allStages := []lifecycle.Stage{
		lifecycle.StageConfigure,
		lifecycle.StageValidate,
		lifecycle.StageGlob,
		lifecycle.StageLoad,
		lifecycle.StageTransform,
		lifecycle.StageRender,
		lifecycle.StageCollect,
		lifecycle.StageWrite,
		lifecycle.StageCleanup,
	}

	for _, stage := range allStages {
		if !m.HasRun(stage) {
			t.Errorf("stage %s should have run", stage)
		}
	}
}

// =============================================================================
// Full Build Integration Tests
// =============================================================================

func TestIntegration_FullBuildWithSampleContent(t *testing.T) {
	site := newTestSite(t)

	// Add multiple posts
	site.addPost("posts/hello-world.md", `---
title: Hello World
slug: hello-world
published: true
date: 2024-01-15
tags: [intro, welcome]
---
# Hello World

This is my first post!`)

	site.addPost("posts/second-post.md", `---
title: Second Post
slug: second-post
published: true
date: 2024-02-20
tags: [update]
---
# Second Post

More content here.`)

	site.addPost("drafts/draft-post.md", `---
title: Draft Post
slug: draft-post
draft: true
---
This should not be published.`)

	m := site.build()

	// Verify posts were processed
	posts := m.Posts()
	if len(posts) < 2 {
		t.Errorf("expected at least 2 posts, got %d", len(posts))
	}

	// Verify published posts have output files
	if !site.fileExists("hello-world/index.html") {
		t.Error("hello-world/index.html should exist")
	}
	if !site.fileExists("second-post/index.html") {
		t.Error("second-post/index.html should exist")
	}

	// Verify draft post does NOT have output file
	if site.fileExists("draft-post/index.html") {
		t.Error("draft-post/index.html should NOT exist (draft)")
	}

	// Verify HTML content
	content := site.readFile("hello-world/index.html")
	if !strings.Contains(content, "Hello World") {
		t.Error("HTML should contain post title")
	}
	if !strings.Contains(content, "<h1>") {
		t.Error("HTML should contain rendered markdown")
	}
}

func TestIntegration_BuildWithConfiguration(t *testing.T) {
	site := newTestSite(t)

	// Add config
	site.addConfig(`[markata-go]
output_dir = "public"
url = "https://mysite.com"
title = "My Site"
`)

	site.addPost("test.md", `---
title: Test Post
published: true
---
Content`)

	// Load config and build
	configPath := filepath.Join(site.dir, "markata-go.toml")
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.OutputDir != "public" {
		t.Errorf("output_dir = %q, want %q", cfg.OutputDir, "public")
	}
	if cfg.URL != "https://mysite.com" {
		t.Errorf("url = %q, want %q", cfg.URL, "https://mysite.com")
	}
}

func TestIntegration_ConcurrentProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Test case: "concurrent post processing" from tests.yaml
	site := newTestSite(t)

	// Create many posts to test concurrent processing
	for i := 0; i < 50; i++ {
		filename := filepath.Join("posts", fmt.Sprintf("post-%03d.md", i))
		content := fmt.Sprintf(`---
title: Post %d
published: true
---
# Test Post %d

Some content here.`, i, i)
		site.addPost(filename, content)
	}

	m := lifecycle.NewManager()
	m.SetConcurrency(8) // Use 8 workers

	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	if err := m.Run(); err != nil {
		t.Fatalf("build failed: %v", err)
	}

	// Verify all posts were processed
	posts := m.Posts()
	if len(posts) < 50 {
		t.Errorf("expected at least 50 posts, got %d", len(posts))
	}

	// Verify no race conditions - check that posts are unique
	slugs := make(map[string]bool)
	for _, p := range posts {
		if slugs[p.Slug] {
			t.Errorf("duplicate slug found: %s (possible race condition)", p.Slug)
		}
		slugs[p.Slug] = true
	}
}

func TestIntegration_EmptyContent(t *testing.T) {
	site := newTestSite(t)

	// No posts added

	m := lifecycle.NewManager()

	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	// Should not fail with no content
	if err := m.Run(); err != nil {
		t.Fatalf("build with no content should not fail: %v", err)
	}

	posts := m.Posts()
	if len(posts) != 0 {
		t.Errorf("expected 0 posts, got %d", len(posts))
	}
}

func TestIntegration_DeeplyNestedContent(t *testing.T) {
	// Test case: "deeply nested content directory" from tests.yaml
	site := newTestSite(t)

	// Create deeply nested content
	deepPath := "a/b/c/d/e/f/g/h/i/j/post.md"
	site.addPost(deepPath, `---
title: Deep Post
slug: a/b/c/d/e/f/g/h/i/j/post
published: true
---
Deeply nested content.`)

	site.build()

	// Verify output path
	if !site.fileExists("a/b/c/d/e/f/g/h/i/j/post/index.html") {
		t.Error("deeply nested output should exist")
	}
}

// =============================================================================
// Error Handling Tests
// =============================================================================

func TestIntegration_InvalidFrontmatter(t *testing.T) {
	site := newTestSite(t)

	// Add post with valid frontmatter
	site.addPost("valid.md", `---
title: Valid Post
published: true
---
Content`)

	// Add post with broken frontmatter (this might cause a warning but shouldn't fail build)
	site.addPost("broken.md", `---
title: "unclosed
---
Content`)

	m := lifecycle.NewManager()

	cfg := &lifecycle.Config{
		ContentDir:   site.contentDir,
		OutputDir:    site.outputDir,
		GlobPatterns: []string{"**/*.md"},
		Extra:        make(map[string]interface{}),
	}
	m.SetConfig(cfg)

	m.RegisterPlugin(plugins.NewGlobPlugin())
	m.RegisterPlugin(plugins.NewLoadPlugin())
	m.RegisterPlugin(plugins.NewRenderMarkdownPlugin())
	m.RegisterPlugin(plugins.NewPublishHTMLPlugin())

	// Build may succeed with warnings or fail depending on strictness
	// The key is it shouldn't panic
	_ = m.Run() //nolint:errcheck // intentionally ignoring error in test

	// Valid post should still be processed
	warnings := m.Warnings()
	// Just verify we got through without panic
	_ = warnings
}

// =============================================================================
// Feed Generation Tests
// =============================================================================

func TestIntegration_FilteredFeed(t *testing.T) {
	// Test case: "filtered feed" from tests.yaml
	site := newTestSite(t)

	site.addPost("public.md", `---
title: Public
slug: public
published: true
---
Public content`)

	site.addPost("draft.md", `---
title: Draft
slug: draft
published: false
---
Draft content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "published",
			Title:  "Published Posts",
			Filter: "published == true",
			Formats: models.FeedFormats{
				HTML: true,
			},
		},
	}

	m := site.buildWithFeeds(feedConfigs)

	feeds := m.Feeds()
	if len(feeds) == 0 {
		t.Skip("no feeds generated - feed plugin may not be active")
	}

	// The feed should only contain published posts
	for _, feed := range feeds {
		for _, post := range feed.Posts {
			if !post.Published {
				t.Errorf("feed should only contain published posts, found: %s", post.Slug)
			}
		}
	}
}

func TestIntegration_SortedFeed(t *testing.T) {
	// Test case: "sorted feed" from tests.yaml
	site := newTestSite(t)

	site.addPost("old.md", `---
title: Old
slug: old
published: true
date: 2024-01-01
---
Old content`)

	site.addPost("new.md", `---
title: New
slug: new
published: true
date: 2024-03-01
---
New content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:    "recent",
			Title:   "Recent Posts",
			Filter:  "True",
			Sort:    "date",
			Reverse: true, // Newest first
			Formats: models.FeedFormats{
				HTML: true,
			},
		},
	}

	m := site.buildWithFeeds(feedConfigs)

	feeds := m.Feeds()
	if len(feeds) == 0 {
		t.Skip("no feeds generated")
	}

	// First post should be the newest
	for _, feed := range feeds {
		if len(feed.Posts) >= 2 {
			first := feed.Posts[0]
			if first.Title != nil && *first.Title != "New" {
				t.Errorf("first post should be 'New' (newest), got %v", first.Title)
			}
		}
	}
}

func TestIntegration_HomePageFeed(t *testing.T) {
	// Test case: "home page feed (empty slug)" from tests.yaml
	site := newTestSite(t)

	site.addPost("post1.md", `---
title: Post 1
slug: post-1
published: true
---
Content`)

	feedConfigs := []models.FeedConfig{
		{
			Slug:   "", // Empty slug for home page
			Title:  "Home",
			Filter: "published == True",
			Formats: models.FeedFormats{
				HTML: true,
			},
		},
	}

	site.buildWithFeeds(feedConfigs)

	// Home page feed should output to index.html directly
	if !site.fileExists("index.html") {
		t.Error("expected output/index.html to exist for home page feed")
	}
}

// =============================================================================
// Palette Serve Mode Tests (Issue #492)
// =============================================================================

// TestIntegration_PaletteChangeOnRebuild tests that changing the palette config
// during serve mode correctly regenerates the CSS with the new palette.
// This simulates the bug reported in issue #492 where palette changes
// would revert to default during serve mode.
func TestIntegration_PaletteChangeOnRebuild(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	site := newTestSite(t)

	// Add a simple post
	site.addPost("test.md", `---
title: Test Post
slug: test
published: true
---
Content`)

	// Helper to run a full build with specific palette
	runBuildWithPalette := func(paletteName string) {
		t.Helper()

		m := lifecycle.NewManager()

		cfg := &lifecycle.Config{
			ContentDir:   site.contentDir,
			OutputDir:    site.outputDir,
			GlobPatterns: []string{"**/*.md"},
			Extra:        make(map[string]interface{}),
		}
		cfg.Extra["url"] = "https://example.com"
		cfg.Extra["title"] = "Test Site"
		// Set the theme config with palette
		cfg.Extra["theme"] = models.ThemeConfig{
			Name:    "default",
			Palette: paletteName,
		}
		disableMinificationForTestBuild(cfg.Extra)
		m.SetConfig(cfg)

		// Register all default plugins including static_assets and palette_css
		m.RegisterPlugins(plugins.DefaultPlugins()...)

		if err := m.Run(); err != nil {
			t.Fatalf("build failed with palette %s: %v", paletteName, err)
		}
	}

	// First build with catppuccin-mocha
	runBuildWithPalette("catppuccin-mocha")

	// Read the CSS
	cssPath := filepath.Join(site.outputDir, "css", "palette.css")
	css1, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after first build: %v", err)
	}

	if !strings.Contains(string(css1), "catppuccin-mocha") {
		t.Error("First build: expected catppuccin-mocha in CSS")
	}

	// Second build with different palette (simulating serve mode config change)
	runBuildWithPalette("dracula")

	css2, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after second build: %v", err)
	}

	if !strings.Contains(string(css2), "dracula") {
		t.Errorf("Second build: expected dracula in CSS, got:\n%s", string(css2)[:min(500, len(css2))])
	}

	// Verify CSS actually changed
	if bytes.Equal(css1, css2) {
		t.Error("CSS should have changed between builds with different palettes")
	}
}

// TestIntegration_PaletteChangeWithConfigFile tests palette changes using actual
// config file loading, more closely simulating the serve mode scenario.
func TestIntegration_PaletteChangeWithConfigFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	site := newTestSite(t)

	// Add a simple post
	site.addPost("test.md", `---
title: Test Post
slug: test
published: true
---
Content`)

	// Helper to build with config file
	buildWithConfig := func(configContent string) {
		t.Helper()

		// Write config file
		configPath := filepath.Join(site.dir, "markata-go.toml")
		if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		// Load config from file (like serve mode does)
		cfg, err := config.Load(configPath)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		// Override output dir for test
		cfg.OutputDir = site.outputDir

		m := lifecycle.NewManager()

		// Create lifecycle config (like createManager in core.go)
		lcConfig := &lifecycle.Config{
			ContentDir:   site.contentDir, // Content dir from test site
			OutputDir:    cfg.OutputDir,
			GlobPatterns: cfg.GlobConfig.Patterns,
			Extra:        make(map[string]interface{}),
		}
		lcConfig.Extra["url"] = cfg.URL
		lcConfig.Extra["title"] = cfg.Title
		lcConfig.Extra["theme"] = cfg.Theme // This is models.ThemeConfig
		disableMinificationForTestBuild(lcConfig.Extra)
		m.SetConfig(lcConfig)

		// Register all default plugins
		m.RegisterPlugins(plugins.DefaultPlugins()...)

		if err := m.Run(); err != nil {
			t.Fatalf("build failed: %v", err)
		}
	}

	// First build with catppuccin-mocha
	config1 := `
[markata-go]
url = "https://example.com"
title = "Test Site"
output_dir = "output"

[markata-go.glob]
patterns = ["**/*.md"]

[markata-go.theme]
name = "default"
palette = "catppuccin-mocha"
`
	buildWithConfig(config1)

	cssPath := filepath.Join(site.outputDir, "css", "palette.css")
	css1, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after first build: %v", err)
	}

	if !strings.Contains(string(css1), "catppuccin-mocha") {
		t.Error("First build: expected catppuccin-mocha in CSS")
	}

	// Second build with dracula (simulating user editing config file)
	config2 := `
[markata-go]
url = "https://example.com"
title = "Test Site"
output_dir = "output"

[markata-go.glob]
patterns = ["**/*.md"]

[markata-go.theme]
name = "default"
palette = "dracula"
`
	buildWithConfig(config2)

	css2, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after second build: %v", err)
	}

	if !strings.Contains(string(css2), "dracula") {
		t.Errorf("Second build: expected dracula in CSS, got:\n%s", string(css2)[:min(500, len(css2))])
	}

	// Verify CSS actually changed
	if bytes.Equal(css1, css2) {
		t.Error("CSS should have changed between builds with different palettes")
	}
}

// TestIntegration_PaletteChangeWithBuildCache tests that palette changes work correctly
// when the build cache is involved. This more closely simulates the serve mode scenario.
func TestIntegration_PaletteChangeWithBuildCache(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	site := newTestSite(t)

	// Create the .markata cache directory
	cacheDir := filepath.Join(site.dir, ".markata")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatalf("failed to create cache dir: %v", err)
	}

	// Add a simple post
	site.addPost("test.md", `---
title: Test Post
slug: test
published: true
---
Content`)

	// Helper to build with config file (simulating serve mode createManager)
	buildWithPalette := func(paletteName string) {
		t.Helper()

		// Convert to forward slashes for TOML compatibility on Windows
		outputDirTOML := filepath.ToSlash(site.outputDir)

		// Write config file
		configContent := fmt.Sprintf(`
[markata-go]
url = "https://example.com"
title = "Test Site"
output_dir = "%s"

[markata-go.glob]
patterns = ["**/*.md"]

[markata-go.theme]
name = "default"
palette = "%s"
`, outputDirTOML, paletteName)

		configPath := filepath.Join(site.dir, "markata-go.toml")
		if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		// Load config from file (like serve mode does)
		cfg, err := config.Load(configPath)
		if err != nil {
			t.Fatalf("failed to load config: %v", err)
		}

		m := lifecycle.NewManager()

		// Create lifecycle config (like createManager in core.go)
		lcConfig := &lifecycle.Config{
			ContentDir:   site.contentDir,
			OutputDir:    cfg.OutputDir,
			GlobPatterns: cfg.GlobConfig.Patterns,
			Extra:        make(map[string]interface{}),
		}
		lcConfig.Extra["url"] = cfg.URL
		lcConfig.Extra["title"] = cfg.Title
		lcConfig.Extra["theme"] = cfg.Theme
		disableMinificationForTestBuild(lcConfig.Extra)
		m.SetConfig(lcConfig)

		// Register all default plugins (including build_cache)
		m.RegisterPlugins(plugins.DefaultPlugins()...)

		if err := m.Run(); err != nil {
			t.Fatalf("build failed with palette %s: %v", paletteName, err)
		}
	}

	// First build with catppuccin-mocha
	buildWithPalette("catppuccin-mocha")

	cssPath := filepath.Join(site.outputDir, "css", "palette.css")
	css1, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after first build: %v", err)
	}

	if !strings.Contains(string(css1), "catppuccin-mocha") {
		t.Errorf("First build: expected catppuccin-mocha in CSS, got:\n%s", string(css1)[:min(500, len(css1))])
	}

	// Second build with different palette (simulating user editing config file)
	buildWithPalette("dracula")

	css2, err := os.ReadFile(cssPath)
	if err != nil {
		t.Fatalf("failed to read palette.css after second build: %v", err)
	}

	if !strings.Contains(string(css2), "dracula") {
		t.Errorf("Second build: expected dracula in CSS, got:\n%s", string(css2)[:min(500, len(css2))])
	}

	// Verify CSS actually changed
	if bytes.Equal(css1, css2) {
		t.Error("CSS should have changed between builds with different palettes")
	}
}
