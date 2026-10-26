package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/terminalpage"
)

// createTestManager creates a minimal lifecycle manager for testing.
func createTestManager(t *testing.T, config *lifecycle.Config) *lifecycle.Manager {
	t.Helper()
	m := lifecycle.NewManager()
	// Set the config in the manager
	m.SetConfig(config)
	return m
}

func renderOGHTML(t *testing.T, post *models.Post) string {
	t.Helper()
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url":          "https://example.com",
			"title":        "Test Site",
			"post_formats": models.PostFormatsConfig{OG: true},
		},
	}
	m := createTestManager(t, config)
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}
	ogPath := filepath.Join(tempDir, post.Slug, "og", "index.html")
	content, err := os.ReadFile(ogPath)
	if err != nil {
		t.Fatalf("failed to read OG card: %v", err)
	}
	return string(content)
}

// TestPublishHTMLPlugin_ShadowPages tests that unpublished posts are rendered as shadow pages.
func TestPublishHTMLPlugin_ShadowPages(t *testing.T) {
	tests := []struct {
		name          string
		post          *models.Post
		wantRendered  bool
		wantInSitemap bool
	}{
		{
			name: "published post is rendered",
			post: &models.Post{
				Path:        "published.md",
				Slug:        "published-post",
				HTML:        "<p>Published content</p>",
				Published:   true,
				Draft:       false,
				Skip:        false,
				ArticleHTML: "<p>Published content</p>",
			},
			wantRendered:  true,
			wantInSitemap: true,
		},
		{
			name: "unpublished post is rendered as shadow page",
			post: &models.Post{
				Path:        "unpublished.md",
				Slug:        "shadow-post",
				HTML:        "<p>Shadow content</p>",
				Published:   false,
				Draft:       false,
				Skip:        false,
				ArticleHTML: "<p>Shadow content</p>",
			},
			wantRendered:  true,
			wantInSitemap: false,
		},
		{
			name: "draft post is not rendered",
			post: &models.Post{
				Path:        "draft.md",
				Slug:        "draft-post",
				HTML:        "<p>Draft content</p>",
				Published:   false,
				Draft:       true,
				Skip:        false,
				ArticleHTML: "<p>Draft content</p>",
			},
			wantRendered:  false,
			wantInSitemap: false,
		},
		{
			name: "skipped post is not rendered",
			post: &models.Post{
				Path:        "skipped.md",
				Slug:        "skipped-post",
				HTML:        "<p>Skipped content</p>",
				Published:   true,
				Draft:       false,
				Skip:        true,
				ArticleHTML: "<p>Skipped content</p>",
			},
			wantRendered:  false,
			wantInSitemap: false,
		},
		{
			name: "published draft is not rendered",
			post: &models.Post{
				Path:        "published-draft.md",
				Slug:        "published-draft",
				HTML:        "<p>Published draft content</p>",
				Published:   true,
				Draft:       true,
				Skip:        false,
				ArticleHTML: "<p>Published draft content</p>",
			},
			wantRendered:  false,
			wantInSitemap: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tempDir := t.TempDir()

			// Create plugin
			plugin := NewPublishHTMLPlugin()

			// Create config
			config := &lifecycle.Config{
				OutputDir: tempDir,
				Extra:     make(map[string]interface{}),
			}

			// Create manager for testing
			m := createTestManager(t, config)

			// Write post
			err := plugin.writePost(tt.post, config, nil, m)
			if err != nil {
				t.Fatalf("writePost() error = %v", err)
			}

			// Check if file was created
			outputPath := filepath.Join(tempDir, tt.post.Slug, "index.html")
			_, statErr := os.Stat(outputPath)
			wasRendered := statErr == nil

			if wasRendered != tt.wantRendered {
				t.Errorf("post rendered = %v, want %v", wasRendered, tt.wantRendered)
			}
		})
	}
}

func TestSafeOutputPath(t *testing.T) {
	root := t.TempDir()
	path, err := safeOutputPath(root, filepath.Join("nested", "slug"))
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(root, "nested", "slug"); path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}
	absolute := filepath.VolumeName(root) + string(filepath.Separator) + "tmp" + string(filepath.Separator) + "unsafe"
	for _, traversal := range []string{"..", filepath.Join("..", "shared"), absolute} {
		if _, err := safeOutputPath(root, traversal); err == nil {
			t.Fatalf("safeOutputPath accepted traversal path %q", traversal)
		}
	}
}

func TestPublishHTMLPlugin_WritesExplicitHomepageAtOutputRoot(t *testing.T) {
	outputDir := t.TempDir()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra:     make(map[string]interface{}),
	}
	post := &models.Post{
		Path:        "pages/post/index.md",
		Slug:        "",
		HTML:        "<html><body>Homepage</body></html>",
		ArticleHTML: "<p>Homepage</p>",
		Published:   true,
	}
	post.Set("_slug_explicit", true)

	if err := NewPublishHTMLPlugin().writePost(post, config, nil, createTestManager(t, config)); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	content, err := os.ReadFile(filepath.Join(outputDir, "index.html"))
	if err != nil {
		t.Fatalf("read homepage output: %v", err)
	}
	if string(content) != post.HTML {
		t.Fatalf("homepage output = %q, want %q", content, post.HTML)
	}
}

// TestPublishHTMLPlugin_OGCardCanonicalURL tests that OG cards include canonical URL and robots meta.
func TestPublishHTMLPlugin_OGCardCanonicalURL(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	// Create config with OG format enabled and a site URL
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url":          "https://example.com",
			"title":        "Test Site",
			"post_formats": models.PostFormatsConfig{OG: true},
		},
	}

	// Create test post
	title := "Test Post Title"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "test-post",
		Title:       &title,
		HTML:        "<html><body>Test content</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>Test content</p>",
	}

	// Create manager for testing
	m := createTestManager(t, config)

	// Write post (which includes OG format)
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	// Read OG card content
	ogPath := filepath.Join(tempDir, "test-post", "og", "index.html")
	content, err := os.ReadFile(ogPath)
	if err != nil {
		t.Fatalf("failed to read OG card: %v", err)
	}

	ogHTML := string(content)

	// Verify canonical URL is present
	expectedCanonical := `<link rel="canonical" href="https://example.com/test-post/">`
	if !strings.Contains(ogHTML, expectedCanonical) {
		t.Errorf("OG card should contain canonical URL.\nExpected: %s\nGot: %s", expectedCanonical, ogHTML)
	}

	// Verify robots meta is present
	expectedRobots := `<meta name="robots" content="noindex, nofollow">`
	if !strings.Contains(ogHTML, expectedRobots) {
		t.Errorf("OG card should contain robots meta.\nExpected: %s\nGot: %s", expectedRobots, ogHTML)
	}
}

func TestPublishHTMLPlugin_RecordsEmptyPostInputWithoutHidingStaleHTML(t *testing.T) {
	for _, test := range []struct {
		name             string
		createStaleHTML  bool
		wantNeedsRebuild bool
	}{
		{name: "no prior HTML", wantNeedsRebuild: false},
		{name: "stale HTML exists", createStaleHTML: true, wantNeedsRebuild: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			outputDir := t.TempDir()
			config := &lifecycle.Config{OutputDir: outputDir, Extra: map[string]interface{}{}}
			post := &models.Post{Path: "empty.md", Slug: "empty", InputHash: "empty-hash", Published: true}
			postDir := filepath.Join(outputDir, post.Slug)
			if test.createStaleHTML {
				if err := os.MkdirAll(postDir, 0o755); err != nil {
					t.Fatalf("create post directory: %v", err)
				}
				if err := os.WriteFile(filepath.Join(postDir, "index.html"), []byte("stale"), 0o600); err != nil {
					t.Fatalf("write stale HTML: %v", err)
				}
			}

			cache := buildcache.New(t.TempDir())
			manager := createTestManager(t, config)
			manager.Cache().Set("build_cache", cache)
			if err := NewPublishHTMLPlugin().writePost(post, config, nil, manager); err != nil {
				t.Fatalf("writePost() error = %v", err)
			}

			if got := cache.ShouldRebuild(post.Path, post.InputHash, post.Template); got != test.wantNeedsRebuild {
				t.Fatalf("ShouldRebuild() = %t, want %t", got, test.wantNeedsRebuild)
			}
		})
	}
}

func TestPublishHTMLPlugin_DoesNotCommitInputAfterPartialFormatFailure(t *testing.T) {
	outputDir := t.TempDir()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{Markdown: true},
		},
	}
	post := &models.Post{
		Path:        "post.md",
		Slug:        "post",
		InputHash:   "new-hash",
		Published:   true,
		ArticleHTML: "<p>new</p>",
	}
	cache := buildcache.New(t.TempDir())
	cache.MarkRebuilt(post.Path, "old-hash", filepath.Join(outputDir, post.Slug, "index.html"), post.Template)
	manager := createTestManager(t, config)
	manager.Cache().Set("build_cache", cache)

	if err := os.Mkdir(filepath.Join(outputDir, post.Slug+".md"), 0o755); err != nil {
		t.Fatalf("create Markdown output collision: %v", err)
	}
	if err := NewPublishHTMLPlugin().writePost(post, config, nil, manager); err == nil {
		t.Fatal("writePost() unexpectedly succeeded")
	}
	if !cache.ShouldRebuild(post.Path, post.InputHash, post.Template) {
		t.Fatal("partial format failure committed the new input hash")
	}
	if _, err := os.Stat(filepath.Join(outputDir, post.Slug, "index.html")); err != nil {
		t.Fatalf("HTML format was not written before the induced later failure: %v", err)
	}
}

func TestPublishHTMLPlugin_RecordsSlugForFormatOnlyPost(t *testing.T) {
	outputDir := t.TempDir()
	disabled := false
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{HTML: &disabled, Markdown: true},
		},
	}
	post := &models.Post{Path: "post.md", Slug: "post", InputHash: "hash", Published: true, Content: "body"}
	cache := buildcache.New(t.TempDir())
	manager := createTestManager(t, config)
	manager.Cache().Set("build_cache", cache)

	if err := NewPublishHTMLPlugin().writePost(post, config, nil, manager); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}
	if got := cache.Posts[post.Path].Slug; got != post.Slug {
		t.Fatalf("cached slug = %q, want %q", got, post.Slug)
	}
	if got := cache.Posts[post.Path].OutputPath; got != "" {
		t.Fatalf("cached HTML output path = %q, want empty for format-only post", got)
	}
}

func TestPublishHTMLPlugin_SkipsUnchangedOGOnlyPost(t *testing.T) {
	outputDir := t.TempDir()
	disabled := false
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{HTML: &disabled, OG: true},
		},
	}
	post := &models.Post{Path: "post.md", Slug: "post", InputHash: "hash", Published: true, Content: "body"}
	cache := buildcache.New(t.TempDir())
	manager := createTestManager(t, config)
	manager.Cache().Set("build_cache", cache)
	manager.SetPosts([]*models.Post{post})
	plugin := NewPublishHTMLPlugin()

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}
	cache.ResetStats()
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if skipped, rebuilt := cache.Stats(); skipped != 1 || rebuilt != 0 {
		t.Fatalf("stats = (%d skipped, %d rebuilt), want (1, 0)", skipped, rebuilt)
	}
}

func TestPublishHTMLPlugin_RecreatesMissingFormatOnlyRedirect(t *testing.T) {
	outputDir := t.TempDir()
	disabled := false
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{HTML: &disabled, Markdown: true},
		},
	}
	post := &models.Post{Path: "post.md", Slug: "post", InputHash: "hash", Published: true, Content: "body"}
	cache := buildcache.New(t.TempDir())
	manager := createTestManager(t, config)
	manager.Cache().Set("build_cache", cache)
	manager.SetPosts([]*models.Post{post})
	plugin := NewPublishHTMLPlugin()

	if err := plugin.Write(manager); err != nil {
		t.Fatalf("first Write() error = %v", err)
	}
	redirectPath := filepath.Join(outputDir, post.Slug, "index.html")
	formatRedirectPath := filepath.Join(outputDir, post.Slug, "index.md", "index.html")
	for _, path := range []string{redirectPath, formatRedirectPath} {
		if err := os.Remove(path); err != nil {
			t.Fatalf("remove redirect %s: %v", path, err)
		}
	}
	if err := plugin.Write(manager); err != nil {
		t.Fatalf("second Write() error = %v", err)
	}
	if _, err := os.Stat(redirectPath); err != nil {
		t.Fatalf("format-only redirect was not recreated: %v", err)
	}
	if _, err := os.Stat(formatRedirectPath); err != nil {
		t.Fatalf("format-specific redirect was not recreated: %v", err)
	}
}

func TestPublishHTMLPlugin_RecordsDraftInputAfterRemovingOutput(t *testing.T) {
	outputDir := t.TempDir()
	config := &lifecycle.Config{OutputDir: outputDir, Extra: map[string]interface{}{}}
	post := &models.Post{Path: "draft.md", Slug: "draft", InputHash: "new-hash", Template: "post.html", Draft: true}
	outputPath := filepath.Join(outputDir, post.Slug, "index.html")
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		t.Fatalf("create output directory: %v", err)
	}
	if err := os.WriteFile(outputPath, []byte("stale"), 0o600); err != nil {
		t.Fatalf("write stale output: %v", err)
	}
	cache := buildcache.New(t.TempDir())
	cache.MarkRebuiltWithSlug(post.Path, post.Slug, "old-hash", outputPath, post.Template)
	manager := createTestManager(t, config)
	manager.Cache().Set("build_cache", cache)
	manager.SetPosts([]*models.Post{post})

	if err := NewPublishHTMLPlugin().Write(manager); err != nil {
		t.Fatalf("Write() error = %v", err)
	}
	if cache.ShouldRebuild(post.Path, post.InputHash, post.Template) {
		t.Fatal("draft input was not recorded after stale output removal")
	}
	if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
		t.Fatalf("draft output still exists: %v", err)
	}
}

func TestPublishHTMLPlugin_OGCardIncludesFrontmatterImage(t *testing.T) {
	tests := []struct {
		name     string
		extra    map[string]interface{}
		imageURL string
	}{
		{
			name: "uses image field first",
			extra: map[string]interface{}{
				"image":       "https://cdn.example.com/posts/image.webp",
				"cover_image": "https://cdn.example.com/posts/cover.webp",
				"og_image":    "https://cdn.example.com/posts/og.webp",
			},
			imageURL: "https://cdn.example.com/posts/image.webp",
		},
		{
			name: "falls back to cover_image",
			extra: map[string]interface{}{
				"cover_image": "https://cdn.example.com/posts/cover.webp",
			},
			imageURL: "https://cdn.example.com/posts/cover.webp",
		},
		{
			name: "falls back to og_image",
			extra: map[string]interface{}{
				"og_image": "https://cdn.example.com/posts/og.webp",
			},
			imageURL: "https://cdn.example.com/posts/og.webp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			plugin := NewPublishHTMLPlugin()

			config := &lifecycle.Config{
				OutputDir: tempDir,
				Extra: map[string]interface{}{
					"url":          "https://example.com",
					"title":        "Test Site",
					"post_formats": models.PostFormatsConfig{OG: true},
				},
			}

			title := "Post With Frontmatter Image"
			post := &models.Post{
				Path:        "test.md",
				Slug:        "test-post",
				Title:       &title,
				HTML:        "<html><body>Test content</body></html>",
				Published:   true,
				Draft:       false,
				Skip:        false,
				ArticleHTML: "<p>Test content</p>",
				Extra:       tt.extra,
			}

			m := createTestManager(t, config)
			if err := plugin.writePost(post, config, nil, m); err != nil {
				t.Fatalf("writePost() error = %v", err)
			}

			ogPath := filepath.Join(tempDir, "test-post", "og", "index.html")
			content, err := os.ReadFile(ogPath)
			if err != nil {
				t.Fatalf("failed to read OG card: %v", err)
			}

			ogHTML := string(content)
			expectedImageTag := fmt.Sprintf("src=%q", tt.imageURL)
			if !strings.Contains(ogHTML, expectedImageTag) {
				t.Errorf("OG card should include expected image source.\nExpected: %s\nGot: %s", expectedImageTag, ogHTML)
			}
		})
	}
}

func TestPublishHTMLPlugin_OGCardWithoutFrontmatterImageOmitsImageTag(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url":          "https://example.com",
			"title":        "Test Site",
			"post_formats": models.PostFormatsConfig{OG: true},
		},
	}

	title := "Post Without Image"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "test-post",
		Title:       &title,
		HTML:        "<html><body>Test content</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>Test content</p>",
	}

	m := createTestManager(t, config)
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	ogPath := filepath.Join(tempDir, "test-post", "og", "index.html")
	content, err := os.ReadFile(ogPath)
	if err != nil {
		t.Fatalf("failed to read OG card: %v", err)
	}

	ogHTML := string(content)
	if strings.Contains(ogHTML, "class=\"og-image\"") {
		t.Errorf("OG card should not include image tag when no frontmatter image is set. Got: %s", ogHTML)
	}
}

func TestPublishHTMLPlugin_OGCardSupportsVideoFrontmatter(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url":          "https://example.com",
			"title":        "Test Site",
			"post_formats": models.PostFormatsConfig{OG: true},
		},
	}

	title := "Post With Video"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "test-post",
		Title:       &title,
		HTML:        "<html><body>Test content</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>Test content</p>",
		Extra: map[string]interface{}{
			"video": "https://cdn.example.com/posts/clip.mp4",
		},
	}

	m := createTestManager(t, config)
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	ogPath := filepath.Join(tempDir, "test-post", "og", "index.html")
	content, err := os.ReadFile(ogPath)
	if err != nil {
		t.Fatalf("failed to read OG card: %v", err)
	}

	ogHTML := string(content)
	if !strings.Contains(ogHTML, "<video") {
		t.Errorf("OG card should include video tag when video frontmatter is set. Got: %s", ogHTML)
	}
	if strings.Contains(ogHTML, "<img") {
		t.Errorf("OG card should not include image tag when video frontmatter is set. Got: %s", ogHTML)
	}
}

func TestPublishHTMLPlugin_OGCardVideoPosterAlias(t *testing.T) {
	title := "Video Poster Alias"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "poster-alias",
		Title:       &title,
		HTML:        "<html><body>Test content</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>Test content</p>",
		Extra: map[string]interface{}{
			"video":        "https://dropper.wayl.one/media/clip.mp4?token=abc#frag",
			"poster_image": "https://dropper.wayl.one/media/poster.png",
		},
	}
	ogHTML := renderOGHTML(t, post)
	if !strings.Contains(ogHTML, "poster=\"https://dropper.wayl.one/media/poster.png") || !strings.Contains(ogHTML, "w=360") || !strings.Contains(ogHTML, "h=320") {
		t.Errorf("OG card should include poster alias with size parameters. Got: %s", ogHTML)
	}
	if !strings.Contains(ogHTML, "type=\"video/mp4\"") {
		t.Errorf("OG card video source should include correct MIME type even with query/fragment. Got: %s", ogHTML)
	}
}

func TestPublishHTMLPlugin_OGCardVideoPosterFallback(t *testing.T) {
	title := "Video Poster Fallback"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "poster-fallback",
		Title:       &title,
		HTML:        "<html><body>Test content</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>Test content</p>",
		Extra: map[string]interface{}{
			"video": "https://dropper.wayl.one/media/clip.mp4",
		},
	}
	ogHTML := renderOGHTML(t, post)
	if !strings.Contains(ogHTML, "poster=\"https://dropper.wayl.one/media/clip.webp") || !strings.Contains(ogHTML, "w=360") || !strings.Contains(ogHTML, "h=320") {
		t.Errorf("OG card should derive .webp poster for trusted video. Got: %s", ogHTML)
	}
}

// TestPublishHTMLPlugin_ShadowPagesDocumentation tests the expected behavior is documented.
func TestPublishHTMLPlugin_ShadowPagesDocumentation(t *testing.T) {
	// This test documents the shadow pages behavior:
	// - published: true → rendered + in feeds + in sitemap
	// - published: false → rendered (shadow page) + NOT in feeds + NOT in sitemap
	// - draft: true → NOT rendered (regardless of published status)
	// - skip: true → NOT rendered (regardless of published status)

	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra:     make(map[string]interface{}),
	}

	// Shadow page scenario
	shadowPost := &models.Post{
		Path:        "shadow.md",
		Slug:        "shadow-page",
		HTML:        "<html><body>Shadow content accessible via direct URL</body></html>",
		Published:   false, // Not in feeds
		Draft:       false, // Not a draft, so it will be rendered
		Skip:        false,
		ArticleHTML: "<p>Shadow content accessible via direct URL</p>",
	}

	// Create manager for testing
	m := createTestManager(t, config)

	if err := plugin.writePost(shadowPost, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	// Verify shadow page was created
	outputPath := filepath.Join(tempDir, "shadow-page", "index.html")
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Shadow page should be rendered even though published=false")
	}

	// Read and verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read shadow page: %v", err)
	}

	if len(content) == 0 {
		t.Error("Shadow page content should not be empty")
	}
}

// TestPublishHTMLPlugin_FormatRedirectsCreateDirectories tests that .md and .txt formats
// use the directory-based structure with redirects at /slug.ext/index.html pointing
// to content at /slug/index.ext.
// This ensures web servers serve the redirect as HTML at /slug.ext (without trailing slash).
// Fixes: https://github.com/WaylonWalker/markata/issues/395
// Fixes: https://github.com/WaylonWalker/markata/issues/437
// Fixes: https://github.com/WaylonWalker/markata/issues/465
func TestPublishHTMLPlugin_FormatRedirectsCreateDirectories(t *testing.T) {
	// Test each format independently to avoid redirect file overwrites
	// (when both formats are enabled, they share the same slug directory)
	tests := []struct {
		name          string
		enabledFormat string
		ext           string
	}{
		{
			name:          "markdown uses directory-based redirect",
			enabledFormat: "markdown",
			ext:           "md",
		},
		{
			name:          "text uses directory-based redirect",
			enabledFormat: "text",
			ext:           "txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			plugin := NewPublishHTMLPlugin()

			// Create config with only one format enabled at a time
			// HTML is DISABLED to test redirect behavior
			// (when HTML is enabled, /slug/index.html has main content instead of redirect)
			htmlEnabled := false
			postFormats := models.PostFormatsConfig{
				HTML: &htmlEnabled,
			}
			if tt.enabledFormat == "markdown" {
				postFormats.Markdown = true
			} else if tt.enabledFormat == "text" {
				postFormats.Text = true
			}

			config := &lifecycle.Config{
				OutputDir: tempDir,
				Extra: map[string]interface{}{
					"post_formats": postFormats,
				},
			}

			// Create test post
			title := "Test Post"
			post := &models.Post{
				Path:        "test.md",
				Slug:        "test-post",
				Title:       &title,
				Content:     "Test content",
				HTML:        "<html><body>Test</body></html>",
				Published:   true,
				Draft:       false,
				Skip:        false,
				ArticleHTML: "<p>Test</p>",
			}

			// Create manager for testing
			m := createTestManager(t, config)

			// Write post (which includes format redirects)
			if err := plugin.writePost(post, config, nil, m); err != nil {
				t.Fatalf("writePost() error = %v", err)
			}

			// Content should be at /slug.ext
			contentPath := filepath.Join(tempDir, "test-post."+tt.ext)
			// Redirect should be at /slug/index.ext/index.html
			redirectPath := filepath.Join(tempDir, "test-post", "index."+tt.ext, "index.html")
			// When HTML is disabled, /slug/index.html should also be a redirect
			slugRedirectPath := filepath.Join(tempDir, "test-post", "index.html")

			// Check that primary content file exists at /slug.ext
			info, err := os.Stat(contentPath)
			if err != nil {
				t.Fatalf("primary content file %s not found: %v", contentPath, err)
			}
			if info.IsDir() {
				t.Errorf("expected %s to be a file, but it's a directory", contentPath)
			}

			// Read content and verify it's actual content (not redirect HTML)
			content, err := os.ReadFile(contentPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", contentPath, err)
			}
			contentStr := string(content)
			if strings.Contains(contentStr, "<!DOCTYPE html>") {
				t.Error("primary content file should NOT be HTML redirect")
			}
			if !strings.Contains(contentStr, "Test") {
				t.Error("primary content file should contain post content")
			}

			// Check that redirect file exists at /slug/index.ext/index.html (always HTML!)
			redirectInfo, err := os.Stat(redirectPath)
			if err != nil {
				t.Fatalf("redirect file %s not found: %v", redirectPath, err)
			}
			if redirectInfo.IsDir() {
				t.Errorf("expected %s to be a file, but it's a directory", redirectPath)
			}

			// Read redirect and verify it points to the content URL
			redirectContent, err := os.ReadFile(redirectPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", redirectPath, err)
			}
			redirectStr := string(redirectContent)
			if !strings.Contains(redirectStr, "<!DOCTYPE html>") {
				t.Error("redirect file should contain DOCTYPE")
			}
			if !strings.Contains(redirectStr, "meta http-equiv=\"refresh\"") {
				t.Error("redirect file should contain meta refresh")
			}
			// Redirect should point TO the content at /slug.ext
			expectedTarget := "/test-post." + tt.ext
			if !strings.Contains(redirectStr, expectedTarget) {
				t.Errorf("redirect file should point to %s, got: %s", expectedTarget, redirectStr)
			}

			// When HTML is disabled, /slug/index.html should also be a redirect
			slugRedirectContent, err := os.ReadFile(slugRedirectPath)
			if err != nil {
				t.Fatalf("slug redirect file %s not found: %v", slugRedirectPath, err)
			}
			slugRedirectStr := string(slugRedirectContent)
			if !strings.Contains(slugRedirectStr, "meta http-equiv=\"refresh\"") {
				t.Error("/slug/index.html should be a redirect when HTML is disabled")
			}
		})
	}
}

// TestPublishHTMLPlugin_StandardWebTxtFiles tests that standard web txt files
// (robots.txt, llms.txt, humans.txt) are accessible via redirects from their expected URLs.
// Content is at /slug/index.txt, redirect at /slug.txt/index.html points there.
// Fixes: https://github.com/WaylonWalker/markata/issues/395
// Fixes: https://github.com/WaylonWalker/markata/issues/465
func TestPublishHTMLPlugin_StandardWebTxtFiles(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	// Create config with Text format enabled
	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
			},
		},
	}

	// Create manager for testing
	m := createTestManager(t, config)

	// Standard web txt files that should be at root level
	standardFiles := []struct {
		slug    string
		title   string
		content string
	}{
		{"robots", "Robots", "User-agent: *\nAllow: /"},
		{"llms", "LLMs", "# LLMs.txt\n\nThis site allows AI training."},
		{"humans", "Humans", "/* TEAM */\nDeveloper: Test"},
	}

	for _, sf := range standardFiles {
		title := sf.title
		post := &models.Post{
			Path:        sf.slug + ".md",
			Slug:        sf.slug,
			Title:       &title,
			Content:     sf.content,
			HTML:        "<html><body>Test</body></html>",
			Published:   true,
			Draft:       false,
			Skip:        false,
			ArticleHTML: "<p>Test</p>",
		}

		if err := plugin.writePost(post, config, nil, m); err != nil {
			t.Fatalf("writePost() error for %s: %v", sf.slug, err)
		}
	}

	// Verify standard web files have content and redirects
	for _, sf := range standardFiles {
		t.Run(sf.slug+".txt structure", func(t *testing.T) {
			// Check content at /slug.txt (root level for special files)
			contentPath := filepath.Join(tempDir, sf.slug+".txt")
			content, err := os.ReadFile(contentPath)
			if err != nil {
				t.Fatalf("content file %s.txt not found: %v", sf.slug, err)
			}

			// Should contain actual content, not redirect HTML
			contentStr := string(content)
			if strings.Contains(contentStr, "<!DOCTYPE html>") {
				t.Errorf("%s.txt should be content, not HTML redirect", sf.slug)
			}

			// Check HTML redirect at /slug/index.txt/index.html
			redirectPath := filepath.Join(tempDir, sf.slug, "index.txt", "index.html")
			redirectContent, err := os.ReadFile(redirectPath)
			if err != nil {
				t.Fatalf("redirect file %s/index.txt/index.html not found: %v", sf.slug, err)
			}

			// Should be an HTML redirect pointing to the content
			redirectStr := string(redirectContent)
			if !strings.Contains(redirectStr, "http-equiv=\"refresh\"") {
				t.Errorf("%s/index.txt/index.html should be an HTML redirect", sf.slug)
			}
			expectedTarget := fmt.Sprintf("/%s.txt", sf.slug)
			if !strings.Contains(redirectStr, expectedTarget) {
				t.Errorf("%s/index.txt/index.html should redirect to %s, got: %s", sf.slug, expectedTarget, redirectStr)
			}
		})
	}
}

// TestPublishHTMLPlugin_PostFormatsDefaultEnabled tests that Markdown and Text
// formats are enabled by default.
// Fixes: https://github.com/WaylonWalker/markata/issues/395
func TestPublishHTMLPlugin_PostFormatsDefaultEnabled(t *testing.T) {
	defaults := models.NewPostFormatsConfig()

	if !defaults.IsHTMLEnabled() {
		t.Error("HTML should be enabled by default")
	}

	if !defaults.Markdown {
		t.Error("Markdown should be enabled by default")
	}

	if !defaults.Text {
		t.Error("Text should be enabled by default")
	}

	if !defaults.OG {
		t.Error("OG should be enabled by default")
	}
}

// TestPublishHTMLPlugin_TxtTemplateRendering tests that txt output uses templates
// when a template engine is available.
// Fixes: https://github.com/WaylonWalker/markata/issues/397
func TestPublishHTMLPlugin_TxtTemplateRendering(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	// Create config with Text format enabled
	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML: &htmlEnabled,
				Text: true,
			},
		},
	}

	// Create test post
	title := "My Test Title"
	description := "A test description"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "test-post",
		Title:       &title,
		Description: &description,
		Content:     "This is the post content.",
		HTML:        "<html><body>Test</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>This is the post content.</p>",
	}

	// Create manager for testing
	m := createTestManager(t, config)

	// Write post without template engine (should use fallback)
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	// Read txt content from /slug.txt
	txtPath := filepath.Join(tempDir, "test-post.txt")
	content, err := os.ReadFile(txtPath)
	if err != nil {
		t.Fatalf("failed to read txt file: %v", err)
	}

	contentStr := string(content)

	// Verify title is present with underline
	if !strings.Contains(contentStr, "My Test Title") {
		t.Error("txt file should contain the title")
	}
	if !strings.Contains(contentStr, strings.Repeat(terminalpage.DoubleRule, len("My Test Title"))) {
		t.Error("txt file should contain title underline")
	}

	// Verify description is present
	if !strings.Contains(contentStr, "A test description") {
		t.Error("txt file should contain the description")
	}

	// Verify content is present
	if !strings.Contains(contentStr, "This is the post content.") {
		t.Error("txt file should contain the post content")
	}
}

func TestPublishHTMLPlugin_TerminalOutputsIncludeFrontmatterMediaLinks(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML: &htmlEnabled,
				Text: true,
				ANSI: true,
			},
		},
	}

	title := "Media Post"
	post := &models.Post{
		Path:        "media.md",
		Slug:        "media-post",
		Title:       &title,
		Content:     "post body",
		ArticleHTML: "<p>post body</p>",
		Published:   true,
		Extra: map[string]interface{}{
			"image": "https://cdn.example.com/poster.webp",
			"video": "https://cdn.example.com/demo.mp4",
		},
	}

	m := createTestManager(t, config)
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	for _, path := range []string{
		filepath.Join(tempDir, "media-post.txt"),
		filepath.Join(tempDir, "media-post.ansi"),
	} {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("failed to read %s: %v", path, err)
		}
		text := terminalpage.StripANSI(string(content))
		if !strings.Contains(text, "Image: https://cdn.example.com/poster.webp") {
			t.Fatalf("expected image link in %s, got:\n%s", path, text)
		}
		if !strings.Contains(text, "Video: https://cdn.example.com/demo.mp4") {
			t.Fatalf("expected video link in %s, got:\n%s", path, text)
		}
	}
}

// TestPublishHTMLPlugin_FormatExtRedirectsWithHTMLEnabled tests that redirect pages
// are created at /slug.ext/index.html (e.g., /test.txt/index.html) even when HTML
// format is enabled. This allows users to navigate to /test.txt/ and be redirected
// to /test/index.txt.
// Fixes: https://github.com/WaylonWalker/markata/issues/465
func TestPublishHTMLPlugin_FormatExtRedirectsWithHTMLEnabled(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	// Create config with HTML AND other formats enabled
	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     true,
			},
		},
	}

	// Create test post
	title := "Test Post"
	post := &models.Post{
		Path:        "test.md",
		Slug:        "test",
		Title:       &title,
		Content:     "Test content",
		HTML:        "<html><body>Test HTML content</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>Test content</p>",
	}

	// Create manager for testing
	m := createTestManager(t, config)

	// Write post
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	// Expected file structure:
	// /test.txt      -> text content
	// /test.md       -> markdown content
	// /test/index.html    -> HTML content (NOT a redirect)
	// /test/index.txt/index.html -> redirect to /test.txt
	// /test/index.ansi/index.html -> redirect to /test.ansi
	// /test/index.md/index.html  -> redirect to /test.md

	// 1. Verify /test/index.html exists and contains HTML content (NOT a redirect)
	htmlPath := filepath.Join(tempDir, "test", "index.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("failed to read HTML file %s: %v", htmlPath, err)
	}
	htmlStr := string(htmlContent)
	if strings.Contains(htmlStr, "http-equiv=\"refresh\"") {
		t.Error("/test/index.html should contain HTML content, not a redirect")
	}
	if !strings.Contains(htmlStr, "Test HTML content") {
		t.Error("/test/index.html should contain the post's HTML content")
	}

	// 2. Verify /test.txt exists with text content
	txtContentPath := filepath.Join(tempDir, "test.txt")
	txtContent, err := os.ReadFile(txtContentPath)
	if err != nil {
		t.Fatalf("failed to read txt file %s: %v", txtContentPath, err)
	}
	if strings.Contains(string(txtContent), "<!DOCTYPE html>") {
		t.Error("/test.txt should contain text content, not HTML")
	}

	// 3. Verify /test.ansi exists with ANSI terminal content
	ansiContentPath := filepath.Join(tempDir, "test.ansi")
	ansiContent, err := os.ReadFile(ansiContentPath)
	if err != nil {
		t.Fatalf("failed to read ansi file %s: %v", ansiContentPath, err)
	}
	if strings.Contains(string(ansiContent), "<!DOCTYPE html>") {
		t.Error("/test.ansi should contain terminal content, not HTML")
	}

	// 4. Verify /test.md exists with markdown content
	mdContentPath := filepath.Join(tempDir, "test.md")
	mdContent, err := os.ReadFile(mdContentPath)
	if err != nil {
		t.Fatalf("failed to read md file %s: %v", mdContentPath, err)
	}
	if strings.Contains(string(mdContent), "<!DOCTYPE html>") {
		t.Error("/test.md should contain markdown content, not HTML")
	}

	// 5. Verify /test/index.txt/index.html exists and is a redirect to /test.txt
	txtIndexRedirectPath := filepath.Join(tempDir, "test", "index.txt", "index.html")
	txtIndexRedirectContent, err := os.ReadFile(txtIndexRedirectPath)
	if err != nil {
		t.Fatalf("failed to read txt index redirect %s: %v", txtIndexRedirectPath, err)
	}
	if !strings.Contains(string(txtIndexRedirectContent), "/test.txt") {
		t.Error("/test/index.txt/index.html should redirect to /test.txt")
	}

	// 6. Verify /test/index.ansi/index.html exists and is a redirect to /test.ansi
	ansiIndexRedirectPath := filepath.Join(tempDir, "test", "index.ansi", "index.html")
	ansiIndexRedirectContent, err := os.ReadFile(ansiIndexRedirectPath)
	if err != nil {
		t.Fatalf("failed to read ansi index redirect %s: %v", ansiIndexRedirectPath, err)
	}
	if !strings.Contains(string(ansiIndexRedirectContent), "/test.ansi") {
		t.Error("/test/index.ansi/index.html should redirect to /test.ansi")
	}

	// 7. Verify /test/index.md/index.html exists and is a redirect to /test.md
	mdIndexRedirectPath := filepath.Join(tempDir, "test", "index.md", "index.html")
	mdIndexRedirectContent, err := os.ReadFile(mdIndexRedirectPath)
	if err != nil {
		t.Fatalf("failed to read md index redirect %s: %v", mdIndexRedirectPath, err)
	}
	if !strings.Contains(string(mdIndexRedirectContent), "/test.md") {
		t.Error("/test/index.md/index.html should redirect to /test.md")
	}
}

// TestPublishHTMLPlugin_RawTxtForSpecialFiles tests that special files
// (robots.txt, llms.txt, etc.) can use a raw.txt template for content-only output.
// Fixes: https://github.com/WaylonWalker/markata/issues/397
func TestPublishHTMLPlugin_RawTxtForSpecialFiles(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	// Create config with Text format enabled
	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML: &htmlEnabled,
				Text: true,
			},
		},
	}

	// Create robots post - should just output raw content
	robotsContent := "User-agent: *\nAllow: /"
	robotsPost := &models.Post{
		Path:        "robots.md",
		Slug:        "robots",
		Content:     robotsContent,
		HTML:        "<html><body>Test</body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		ArticleHTML: "<p>Test</p>",
	}

	// Create manager for testing
	m := createTestManager(t, config)

	// Write post without template engine (should use fallback)
	if err := plugin.writePost(robotsPost, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	// Read robots content from /robots.txt (root level for special files)
	robotsPath := filepath.Join(tempDir, "robots.txt")
	content, err := os.ReadFile(robotsPath)
	if err != nil {
		t.Fatalf("failed to read robots.txt: %v", err)
	}

	contentStr := string(content)

	// Verify content is present (without title/description header since no title set)
	if !strings.Contains(contentStr, "User-agent: *") {
		t.Error("robots.txt should contain the user-agent directive")
	}
	if !strings.Contains(contentStr, "Allow: /") {
		t.Error("robots.txt should contain the allow directive")
	}

	// Verify HTML redirect exists at /robots/index.txt/index.html
	redirectPath := filepath.Join(tempDir, "robots", "index.txt", "index.html")
	redirectContent, err := os.ReadFile(redirectPath)
	if err != nil {
		t.Fatalf("failed to read robots/index.txt/index.html: %v", err)
	}

	expectedTarget := "/robots.txt"
	if !strings.Contains(string(redirectContent), expectedTarget) {
		t.Errorf("robots/index.txt/index.html should redirect to %s", expectedTarget)
	}
	if !strings.Contains(string(redirectContent), "http-equiv=\"refresh\"") {
		t.Error("robots/index.txt/index.html should be an HTML redirect")
	}
}

// TestPublishHTMLPlugin_PrivatePostsNoAlternateFormats tests that private posts
// only get HTML output (with encrypted content) and do NOT get .md, .txt, or OG
// output files, which would leak plaintext content.
func TestPublishHTMLPlugin_PrivatePostsNoAlternateFormats(t *testing.T) {
	tempDir := t.TempDir()

	htmlEnabled := true

	// Enable all output formats
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     true,
				OG:       true,
			},
		},
	}

	plugin := NewPublishHTMLPlugin()
	m := createTestManager(t, config)

	title := "My Private Post"

	// Private post with encrypted ArticleHTML (simulating what encryption plugin does)
	privatePost := &models.Post{
		Path:        "secret.md",
		Slug:        "secret-post",
		Title:       &title,
		Content:     "", // scrubbed by encryption plugin
		HTML:        "<html><body><div class=\"encrypted-content\" data-encrypted=\"abc123\">Encrypted</div></body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		Private:     true,
		ArticleHTML: "<div class=\"encrypted-content\" data-encrypted=\"abc123\">Encrypted</div>",
	}

	if err := plugin.writePost(privatePost, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	// HTML should be written (contains encrypted content)
	htmlPath := filepath.Join(tempDir, "secret-post", "index.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Error("Private post HTML should be written (contains encrypted content)")
	}

	// .md should NOT be written
	mdPath := filepath.Join(tempDir, "secret-post.md")
	if _, err := os.Stat(mdPath); !os.IsNotExist(err) {
		t.Error("Private post .md file should NOT be written (would leak plaintext)")
	}

	// .txt should NOT be written
	txtPath := filepath.Join(tempDir, "secret-post.txt")
	if _, err := os.Stat(txtPath); !os.IsNotExist(err) {
		t.Error("Private post .txt file should NOT be written (would leak plaintext)")
	}

	// .ansi should NOT be written
	ansiPath := filepath.Join(tempDir, "secret-post.ansi")
	if _, err := os.Stat(ansiPath); !os.IsNotExist(err) {
		t.Error("Private post .ansi file should NOT be written (would leak plaintext)")
	}

	// OG card should NOT be written
	ogPath := filepath.Join(tempDir, "secret-post", "og", "index.html")
	if _, err := os.Stat(ogPath); !os.IsNotExist(err) {
		t.Error("Private post OG card should NOT be written (would leak metadata)")
	}
}

func TestPublishHTMLPlugin_PrivateTransitionRemovesExistingAlternateFormats(t *testing.T) {
	tempDir := t.TempDir()
	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url": "https://example.com",
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     true,
				OG:       true,
			},
		},
	}
	plugin := NewPublishHTMLPlugin()
	m := createTestManager(t, config)

	post := &models.Post{
		Path:        "transition.md",
		Slug:        "transition",
		Content:     "public body",
		HTML:        "<p>public body</p>",
		ArticleHTML: "<p>public body</p>",
		Published:   true,
	}
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("public writePost() error = %v", err)
	}

	post.Private = true
	post.Content = ""
	post.HTML = `<div class="encrypted-content" data-encrypted="ciphertext">Encrypted</div>`
	post.ArticleHTML = post.HTML
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("private writePost() error = %v", err)
	}

	for _, path := range []string{
		filepath.Join(tempDir, "transition.md"),
		filepath.Join(tempDir, "transition.txt"),
		filepath.Join(tempDir, "transition.ansi"),
		filepath.Join(tempDir, "transition", "index.md"),
		filepath.Join(tempDir, "transition", "index.txt"),
		filepath.Join(tempDir, "transition", "index.ansi"),
		filepath.Join(tempDir, "transition", "og"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("private transition retained alternate output %s: %v", path, err)
		}
	}
	html, err := os.ReadFile(filepath.Join(tempDir, "transition", "index.html"))
	if err != nil {
		t.Fatalf("private HTML output was removed: %v", err)
	}
	if !strings.Contains(string(html), "encrypted-content") {
		t.Fatalf("private HTML output did not contain encrypted content: %s", html)
	}
}

func TestPublishHTMLPlugin_PrivateDraftOrSkipRemovesExistingOutputs(t *testing.T) {
	for _, tt := range []struct {
		name  string
		draft bool
		skip  bool
	}{
		{name: "draft", draft: true},
		{name: "skip", skip: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			htmlEnabled := true
			config := &lifecycle.Config{
				OutputDir: tempDir,
				Extra: map[string]interface{}{
					"url": "https://example.com",
					"post_formats": models.PostFormatsConfig{
						HTML:     &htmlEnabled,
						Markdown: true,
						Text:     true,
						ANSI:     true,
						OG:       true,
					},
				},
			}
			plugin := NewPublishHTMLPlugin()
			m := createTestManager(t, config)
			post := &models.Post{
				Path:        "private-transition.md",
				Slug:        "private-transition",
				Content:     "public body",
				HTML:        "<p>public body</p>",
				ArticleHTML: "<p>public body</p>",
				Published:   true,
			}
			if err := plugin.writePost(post, config, nil, m); err != nil {
				t.Fatalf("public writePost() error = %v", err)
			}

			post.Private = true
			post.Draft = tt.draft
			post.Skip = tt.skip
			m.SetPosts([]*models.Post{post})
			if err := plugin.Write(m); err != nil {
				t.Fatalf("private %s Write() error = %v", tt.name, err)
			}

			for _, path := range []string{
				filepath.Join(tempDir, "private-transition"),
				filepath.Join(tempDir, "private-transition.md"),
				filepath.Join(tempDir, "private-transition.txt"),
				filepath.Join(tempDir, "private-transition.ansi"),
			} {
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Errorf("private %s retained output %s: %v", tt.name, path, err)
				}
			}
		})
	}
}

func TestPublishHTMLPlugin_PrivateToDraftOrSkipRemovesExistingOutputs(t *testing.T) {
	for _, tt := range []struct {
		name  string
		draft bool
		skip  bool
	}{
		{name: "draft", draft: true},
		{name: "skip", skip: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			htmlEnabled := true
			config := &lifecycle.Config{
				OutputDir: tempDir,
				Extra: map[string]interface{}{
					"url": "https://example.com",
					"post_formats": models.PostFormatsConfig{
						HTML:     &htmlEnabled,
						Markdown: true,
						Text:     true,
						ANSI:     true,
						OG:       true,
					},
				},
			}
			plugin := NewPublishHTMLPlugin()
			m := createTestManager(t, config)
			post := &models.Post{
				Path:        "private-to-draft.md",
				Slug:        "private-to-draft",
				Content:     "public body",
				HTML:        "<p>public body</p>",
				ArticleHTML: "<p>public body</p>",
				Published:   true,
			}
			if err := plugin.writePost(post, config, nil, m); err != nil {
				t.Fatalf("public writePost() error = %v", err)
			}

			post.Private = true
			post.Content = ""
			post.HTML = `<div class="encrypted-content" data-encrypted="ciphertext">Encrypted</div>`
			post.ArticleHTML = post.HTML
			if err := plugin.writePost(post, config, nil, m); err != nil {
				t.Fatalf("private writePost() error = %v", err)
			}

			post.Private = false
			post.Draft = tt.draft
			post.Skip = tt.skip
			m.SetPosts([]*models.Post{post})
			if err := plugin.Write(m); err != nil {
				t.Fatalf("public %s Write() error = %v", tt.name, err)
			}

			for _, path := range []string{
				filepath.Join(tempDir, "private-to-draft"),
				filepath.Join(tempDir, "private-to-draft.md"),
				filepath.Join(tempDir, "private-to-draft.txt"),
				filepath.Join(tempDir, "private-to-draft.ansi"),
			} {
				if _, err := os.Stat(path); !os.IsNotExist(err) {
					t.Errorf("public %s retained output %s: %v", tt.name, path, err)
				}
			}
		})
	}
}

func TestPublishHTMLPlugin_PrivateSlugChangeRemovesOldOutputs(t *testing.T) {
	tempDir := t.TempDir()
	cache := buildcache.New(filepath.Join(tempDir, "cache"))
	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url": "https://example.com",
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     true,
				OG:       true,
			},
		},
	}
	plugin := NewPublishHTMLPlugin()
	m := createTestManager(t, config)
	m.Cache().Set("build_cache", cache)
	post := &models.Post{
		Path:        "private-slug-change.md",
		Slug:        "old-slug",
		Content:     "public body",
		HTML:        "<p>public body</p>",
		ArticleHTML: "<p>public body</p>",
		Published:   true,
	}
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("public writePost() error = %v", err)
	}
	cache.Posts[post.Path] = &buildcache.PostCache{
		Slug:       "old-slug",
		OutputPath: filepath.Join(tempDir, "old-slug", "index.html"),
	}

	post.Private = true
	post.Slug = "new-slug"
	post.Content = ""
	post.HTML = `<div class="encrypted-content" data-encrypted="ciphertext">Encrypted</div>`
	post.ArticleHTML = post.HTML
	m.SetPosts([]*models.Post{post})
	if err := plugin.Write(m); err != nil {
		t.Fatalf("private slug change Write() error = %v", err)
	}

	for _, path := range []string{
		filepath.Join(tempDir, "old-slug"),
		filepath.Join(tempDir, "old-slug.md"),
		filepath.Join(tempDir, "old-slug.txt"),
		filepath.Join(tempDir, "old-slug.ansi"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("private slug change retained old output %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(tempDir, "new-slug", "index.html")); err != nil {
		t.Fatalf("private slug change did not write new HTML: %v", err)
	}
}

func TestPublishHTMLPlugin_RemovedPostRemovesOwnedOutputsOutsideServe(t *testing.T) {
	tempDir := t.TempDir()
	cache := buildcache.New(filepath.Join(tempDir, "cache"))
	cache.Posts["posts/deleted.md"] = &buildcache.PostCache{
		Slug:       "deleted",
		OutputPath: filepath.Join(tempDir, "deleted", "index.html"),
	}

	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.NewPostFormatsConfig(),
		},
	}
	m := createTestManager(t, config)
	m.Cache().Set("build_cache", cache)
	lifecycle.SetServeRemovedPaths(m, []string{"posts/deleted.md"})

	owned := []string{
		"deleted/index.html",
		"deleted/index.md/index.html",
		"deleted/index.txt/index.html",
		"deleted/index.ansi/index.html",
		"deleted/og/index.html",
		"deleted.md",
		"deleted.txt",
		"deleted.ansi",
	}
	for _, relative := range owned {
		path := filepath.Join(tempDir, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("MkdirAll(%s) error = %v", relative, err)
		}
		if err := os.WriteFile(path, []byte("generated"), 0o600); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", relative, err)
		}
	}
	unrelated := filepath.Join(tempDir, "deleted", "plugin-output.json")
	if err := os.WriteFile(unrelated, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := NewPublishHTMLPlugin().Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	for _, relative := range owned {
		if _, err := os.Stat(filepath.Join(tempDir, filepath.FromSlash(relative))); !os.IsNotExist(err) {
			t.Errorf("removed post retained generated output %s: %v", relative, err)
		}
	}
	if data, err := os.ReadFile(unrelated); err != nil || string(data) != "keep" {
		t.Fatalf("unrelated output = %q, error = %v; want keep", data, err)
	}
}

func TestPublishHTMLPlugin_PrivateTraversalSlugFailsBeforeRemoval(t *testing.T) {
	tempDir := t.TempDir()
	victimPath := filepath.Join(tempDir, "..", "victim.md")
	//nolint:gosec // The sentinel stays inside the test's temporary hierarchy.
	if err := os.WriteFile(victimPath, []byte("sentinel"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(victimPath) })

	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{HTML: boolPtr(true)},
		},
	}
	m := createTestManager(t, config)
	m.SetPosts([]*models.Post{{Path: "private.md", Slug: "../victim", Private: true, Published: true, HTML: "encrypted"}})
	if err := NewPublishHTMLPlugin().Write(m); err == nil {
		t.Fatal("expected traversal slug error")
	}
	data, err := os.ReadFile(victimPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "sentinel" {
		t.Fatalf("traversal slug modified path outside output directory: %q", data)
	}
}

func TestPublishHTMLPlugin_PrivateToPublicSlugChangeRemovesOldOutputs(t *testing.T) {
	tempDir := t.TempDir()
	cache := buildcache.New(filepath.Join(tempDir, "cache"))
	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url": "https://example.com",
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     true,
				OG:       true,
			},
		},
	}
	plugin := NewPublishHTMLPlugin()
	m := createTestManager(t, config)
	m.Cache().Set("build_cache", cache)
	post := &models.Post{
		Path:        "private-to-public.md",
		Slug:        "private-old",
		Content:     "public body",
		HTML:        "<p>public body</p>",
		ArticleHTML: "<p>public body</p>",
		Published:   true,
	}
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("initial public writePost() error = %v", err)
	}
	post.Private = true
	post.Content = ""
	post.HTML = `<div class="encrypted-content" data-encrypted="ciphertext">Encrypted</div>`
	post.ArticleHTML = post.HTML
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("private writePost() error = %v", err)
	}
	cache.Posts[post.Path] = &buildcache.PostCache{
		Slug:       "private-old",
		OutputPath: filepath.Join(tempDir, "private-old", "index.html"),
	}

	post.Private = false
	post.Slug = "public-new"
	post.Content = "new public body"
	post.HTML = "<p>new public body</p>"
	post.ArticleHTML = post.HTML
	m.SetPosts([]*models.Post{post})
	if err := plugin.Write(m); err != nil {
		t.Fatalf("public slug change Write() error = %v", err)
	}

	for _, path := range []string{
		filepath.Join(tempDir, "private-old"),
		filepath.Join(tempDir, "private-old.md"),
		filepath.Join(tempDir, "private-old.txt"),
		filepath.Join(tempDir, "private-old.ansi"),
	} {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("public slug change retained old private output %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(tempDir, "public-new", "index.html")); err != nil {
		t.Fatalf("public slug change did not write new HTML: %v", err)
	}
}

// TestPublishHTMLPlugin_NonPrivatePostGetsAllFormats verifies that non-private posts
// still get all enabled output formats (regression test for the private post fix).
func TestPublishHTMLPlugin_NonPrivatePostGetsAllFormats(t *testing.T) {
	tempDir := t.TempDir()

	htmlEnabled := true

	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"url":   "https://example.com",
			"title": "Test Site",
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     true,
				OG:       false, // OG requires template engine, skip for simplicity
			},
		},
	}

	plugin := NewPublishHTMLPlugin()
	m := createTestManager(t, config)

	title := "Public Post"

	publicPost := &models.Post{
		Path:        "public.md",
		Slug:        "public-post",
		Title:       &title,
		Content:     "Hello world, this is public content.",
		HTML:        "<html><body><p>Hello world</p></body></html>",
		Published:   true,
		Draft:       false,
		Skip:        false,
		Private:     false,
		ArticleHTML: "<p>Hello world</p>",
	}

	if err := plugin.writePost(publicPost, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	// HTML should be written
	htmlPath := filepath.Join(tempDir, "public-post", "index.html")
	if _, err := os.Stat(htmlPath); os.IsNotExist(err) {
		t.Error("Public post HTML should be written")
	}

	// .md should be written
	mdPath := filepath.Join(tempDir, "public-post.md")
	if _, err := os.Stat(mdPath); os.IsNotExist(err) {
		t.Error("Public post .md file should be written")
	}

	// .txt should be written
	txtPath := filepath.Join(tempDir, "public-post.txt")
	if _, err := os.Stat(txtPath); os.IsNotExist(err) {
		t.Error("Public post .txt file should be written")
	}

	// .ansi should be written
	ansiPath := filepath.Join(tempDir, "public-post.ansi")
	if _, err := os.Stat(ansiPath); os.IsNotExist(err) {
		t.Error("Public post .ansi file should be written")
	}
}

func TestResolvePostFormats_MergesFrontmatterOverride(t *testing.T) {
	htmlEnabled := true
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     false,
				OG:       true,
			},
		},
	}

	post := &models.Post{
		Extra: map[string]interface{}{
			"post_formats": map[string]interface{}{
				"ansi": true,
				"og":   false,
			},
		},
	}

	resolved := resolvePostFormats(post, config)

	if !resolved.IsHTMLEnabled() {
		t.Fatal("expected html to inherit enabled site default")
	}
	if !resolved.Markdown {
		t.Fatal("expected markdown to inherit enabled site default")
	}
	if !resolved.Text {
		t.Fatal("expected text to inherit enabled site default")
	}
	if !resolved.ANSI {
		t.Fatal("expected ansi to be enabled by post override")
	}
	if resolved.OG {
		t.Fatal("expected og to be disabled by post override")
	}
}

func TestPublishHTMLPlugin_PostFrontmatterDisablesInheritedFormats(t *testing.T) {
	tempDir := t.TempDir()
	plugin := NewPublishHTMLPlugin()

	htmlEnabled := true
	config := &lifecycle.Config{
		OutputDir: tempDir,
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     true,
				ANSI:     false,
				OG:       false,
			},
		},
	}

	title := "Frontmatter Override"
	post := &models.Post{
		Path:        "override.md",
		Slug:        "override-post",
		Title:       &title,
		Content:     "body",
		HTML:        "<html><body>html</body></html>",
		Published:   true,
		ArticleHTML: "<p>body</p>",
		Extra: map[string]interface{}{
			"post_formats": map[string]interface{}{
				"markdown": false,
				"text":     false,
				"ansi":     true,
			},
		},
	}

	m := createTestManager(t, config)
	if err := plugin.writePost(post, config, nil, m); err != nil {
		t.Fatalf("writePost() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(tempDir, "override-post", "index.html")); os.IsNotExist(err) {
		t.Fatal("expected html output to be written")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "override-post.md")); !os.IsNotExist(err) {
		t.Fatal("expected markdown output to be suppressed by post override")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "override-post.txt")); !os.IsNotExist(err) {
		t.Fatal("expected text output to be suppressed by post override")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "override-post.ansi")); os.IsNotExist(err) {
		t.Fatal("expected ansi output to be enabled by post override")
	}
}
