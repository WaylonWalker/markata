package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// newTestPostIndex creates a PostIndex from a map of href -> post for testing.
func newTestPostIndex(posts map[string]*models.Post) *lifecycle.PostIndex {
	idx := &lifecycle.PostIndex{
		BySlug:      make(map[string]*models.Post),
		BySlugified: make(map[string]*models.Post),
		ByHref:      make(map[string]*models.Post),
		ByPath:      make(map[string]*models.Post),
	}
	for href, post := range posts {
		idx.ByHref[href] = post
		if post.Slug != "" {
			idx.BySlug[strings.ToLower(post.Slug)] = post
			idx.BySlugified[models.Slugify(post.Slug)] = post
		}
	}
	return idx
}

func TestWikilinkHoverPlugin_Name(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	if got := p.Name(); got != "wikilink_hover" {
		t.Errorf("Name() = %q, want %q", got, "wikilink_hover")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_NoWikilinks(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.postIdx = newTestPostIndex(map[string]*models.Post{})

	post := &models.Post{
		ArticleHTML: "<p>Hello world with no wikilinks</p>",
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged
	if post.ArticleHTML != originalHTML {
		t.Error("HTML was modified when no wikilinks present")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_BasicWikilink(t *testing.T) {
	p := NewWikilinkHoverPlugin()

	// Set up target post
	description := "This is a great article about testing"
	targetPost := &models.Post{
		Slug:        "my-article",
		Href:        "/my-article/",
		Description: &description,
	}
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/my-article/": targetPost,
	})

	post := &models.Post{
		ArticleHTML: `<p>Check out <a href="/my-article/" class="wikilink">My Article</a> for more info.</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have data-preview attribute
	if !strings.Contains(post.ArticleHTML, `data-preview="`) {
		t.Error("Expected data-preview attribute")
	}
	if !strings.Contains(post.ArticleHTML, "This is a great article") {
		t.Error("Expected description in preview")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_WithImage(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.config.IncludeImage = true

	// Set up target post with image
	description := "Article description"
	targetPost := &models.Post{
		Slug:        "my-article",
		Href:        "/my-article/",
		Description: &description,
		Extra: map[string]interface{}{
			"image": "/images/featured.jpg",
		},
	}
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/my-article/": targetPost,
	})

	post := &models.Post{
		ArticleHTML: `<a href="/my-article/" class="wikilink">My Article</a>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have data-preview-image attribute
	if !strings.Contains(post.ArticleHTML, `data-preview-image="/images/featured.jpg"`) {
		t.Error("Expected data-preview-image attribute")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_NoImage(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.config.IncludeImage = true

	// Set up target post without image
	description := "Article description"
	targetPost := &models.Post{
		Slug:        "my-article",
		Href:        "/my-article/",
		Description: &description,
	}
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/my-article/": targetPost,
	})

	post := &models.Post{
		ArticleHTML: `<a href="/my-article/" class="wikilink">My Article</a>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should NOT have data-preview-image attribute
	if strings.Contains(post.ArticleHTML, `data-preview-image`) {
		t.Error("Should not have data-preview-image when no image exists")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_ImageDisabled(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.config.IncludeImage = false

	// Set up target post with image
	description := "Article description"
	targetPost := &models.Post{
		Slug:        "my-article",
		Href:        "/my-article/",
		Description: &description,
		Extra: map[string]interface{}{
			"image": "/images/featured.jpg",
		},
	}
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/my-article/": targetPost,
	})

	post := &models.Post{
		ArticleHTML: `<a href="/my-article/" class="wikilink">My Article</a>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should NOT have data-preview-image when disabled
	if strings.Contains(post.ArticleHTML, `data-preview-image`) {
		t.Error("Should not have data-preview-image when disabled")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_ScreenshotService(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.config.ScreenshotService = "https://screenshot.example.com/capture?url="

	description := "Article description"
	targetPost := &models.Post{
		Slug:        "my-article",
		Href:        "/my-article/",
		Description: &description,
	}
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/my-article/": targetPost,
	})

	post := &models.Post{
		ArticleHTML: `<a href="/my-article/" class="wikilink">My Article</a>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have data-preview-screenshot attribute
	if !strings.Contains(post.ArticleHTML, `data-preview-screenshot="https://screenshot.example.com/capture?url=/my-article/"`) {
		t.Error("Expected data-preview-screenshot attribute")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_SkipPost(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.postIdx = newTestPostIndex(map[string]*models.Post{})

	post := &models.Post{
		Skip:        true,
		ArticleHTML: `<a href="/test/" class="wikilink">Test</a>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged for skipped posts
	if post.ArticleHTML != originalHTML {
		t.Error("Skip post HTML was modified")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_EmptyHTML(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.postIdx = newTestPostIndex(map[string]*models.Post{})

	post := &models.Post{
		ArticleHTML: "",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if post.ArticleHTML != "" {
		t.Error("Empty HTML was modified")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_MultipleWikilinks(t *testing.T) {
	p := NewWikilinkHoverPlugin()

	desc1 := "First article description"
	desc2 := "Second article description"
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/first/":  {Slug: "first", Href: "/first/", Description: &desc1},
		"/second/": {Slug: "second", Href: "/second/", Description: &desc2},
	})

	post := &models.Post{
		ArticleHTML: `<p>See <a href="/first/" class="wikilink">First</a> and <a href="/second/" class="wikilink">Second</a>.</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have two data-preview attributes
	count := strings.Count(post.ArticleHTML, `data-preview="`)
	if count != 2 {
		t.Errorf("Expected 2 data-preview attributes, got %d", count)
	}
}

func TestWikilinkHoverPlugin_ProcessPost_UnknownTarget(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.postIdx = newTestPostIndex(map[string]*models.Post{})

	post := &models.Post{
		ArticleHTML: `<a href="/unknown/" class="wikilink">Unknown</a>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged for unknown targets
	if post.ArticleHTML != originalHTML {
		t.Error("Wikilink to unknown target was modified")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_TruncateLongDescription(t *testing.T) {
	p := NewWikilinkHoverPlugin()
	p.config.PreviewLength = 50

	longDesc := "This is a very long description that should be truncated to fit within the preview length limit"
	targetPost := &models.Post{
		Slug:        "long",
		Href:        "/long/",
		Description: &longDesc,
	}
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/long/": targetPost,
	})

	post := &models.Post{
		ArticleHTML: `<a href="/long/" class="wikilink">Long Post</a>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should contain truncated text with ellipsis
	if !strings.Contains(post.ArticleHTML, "...") {
		t.Error("Expected truncated preview with ellipsis")
	}
	// Should not contain the full description
	if strings.Contains(post.ArticleHTML, "preview length limit") {
		t.Error("Description should be truncated")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_FallbackToContent(t *testing.T) {
	p := NewWikilinkHoverPlugin()

	// Post with no description but has content
	targetPost := &models.Post{
		Slug:    "content-only",
		Href:    "/content-only/",
		Content: "This is the raw markdown content of the post.",
	}
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/content-only/": targetPost,
	})

	post := &models.Post{
		ArticleHTML: `<a href="/content-only/" class="wikilink">Content Post</a>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should use content as preview
	if !strings.Contains(post.ArticleHTML, "raw markdown content") {
		t.Error("Expected content to be used as preview fallback")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_RegularLink(t *testing.T) {
	p := NewWikilinkHoverPlugin()

	desc := "Article description"
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/my-article/": {Slug: "my-article", Href: "/my-article/", Description: &desc},
	})

	// Regular link (not wikilink) should not be modified
	post := &models.Post{
		ArticleHTML: `<a href="/my-article/">Regular Link</a>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged for non-wikilink anchors
	if post.ArticleHTML != originalHTML {
		t.Error("Regular link was modified")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_PreservesExistingAttributes(t *testing.T) {
	p := NewWikilinkHoverPlugin()

	desc := "Article description"
	p.postIdx = newTestPostIndex(map[string]*models.Post{
		"/my-article/": {Slug: "my-article", Href: "/my-article/", Description: &desc},
	})

	// Wikilink with existing attributes
	post := &models.Post{
		ArticleHTML: `<a href="/my-article/" class="wikilink" title="Existing Title">My Article</a>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should preserve existing attributes
	if !strings.Contains(post.ArticleHTML, `title="Existing Title"`) {
		t.Error("Existing title attribute was lost")
	}
	if !strings.Contains(post.ArticleHTML, `class="wikilink"`) {
		t.Error("Existing class attribute was lost")
	}
	// And should add data-preview
	if !strings.Contains(post.ArticleHTML, `data-preview="`) {
		t.Error("Expected data-preview to be added")
	}
}

func TestWikilinkHoverPlugin_ProcessPost_ImageFieldVariants(t *testing.T) {
	testCases := []struct {
		name       string
		imageField string
	}{
		{"image", "image"},
		{"featured_image", "featured_image"},
		{"cover_image", "cover_image"},
		{"og_image", "og_image"},
		{"thumbnail", "thumbnail"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewWikilinkHoverPlugin()
			p.config.IncludeImage = true

			desc := "Description"
			targetPost := &models.Post{
				Slug:        "test",
				Href:        "/test/",
				Description: &desc,
				Extra: map[string]interface{}{
					tc.imageField: "/images/test.jpg",
				},
			}
			p.postIdx = newTestPostIndex(map[string]*models.Post{
				"/test/": targetPost,
			})

			post := &models.Post{
				ArticleHTML: `<a href="/test/" class="wikilink">Test</a>`,
			}

			err := p.processPost(post)
			if err != nil {
				t.Errorf("processPost() error = %v", err)
			}

			if !strings.Contains(post.ArticleHTML, `data-preview-image="/images/test.jpg"`) {
				t.Errorf("Expected %s field to be used for preview image", tc.imageField)
			}
		})
	}
}

func TestTruncatePreviewText(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "short text unchanged",
			input:  "Hello world",
			maxLen: 50,
			want:   "Hello world",
		},
		{
			name:   "truncate at word boundary",
			input:  "This is a longer text that needs truncation",
			maxLen: 20,
			want:   "This is a longer...",
		},
		{
			name:   "collapse whitespace",
			input:  "Text  with   multiple    spaces",
			maxLen: 100,
			want:   "Text with multiple spaces",
		},
		{
			name:   "collapse newlines",
			input:  "Text\nwith\nnewlines",
			maxLen: 100,
			want:   "Text with newlines",
		},
		{
			name:   "trim whitespace",
			input:  "  trimmed  ",
			maxLen: 100,
			want:   "trimmed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncatePreviewText(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncatePreviewText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStripHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "remove tags",
			input: "<p>Hello <strong>world</strong></p>",
			want:  "Hello world",
		},
		{
			name:  "decode entities",
			input: "Tom &amp; Jerry &lt;3",
			want:  "Tom & Jerry <3",
		},
		{
			name:  "plain text unchanged",
			input: "Plain text",
			want:  "Plain text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripHTML(tt.input)
			if got != tt.want {
				t.Errorf("stripHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestWikilinkHoverPlugin_Config(t *testing.T) {
	p := NewWikilinkHoverPlugin()

	// Default config
	cfg := p.Config()
	if !cfg.Enabled {
		t.Error("Expected Enabled to be true by default")
	}
	if cfg.PreviewLength != 200 {
		t.Errorf("Expected PreviewLength to be 200, got %d", cfg.PreviewLength)
	}
	if !cfg.IncludeImage {
		t.Error("Expected IncludeImage to be true by default")
	}
	if cfg.ScreenshotService != "" {
		t.Error("Expected ScreenshotService to be empty by default")
	}

	// Set custom config
	customCfg := models.WikilinkHoverConfig{
		Enabled:           false,
		PreviewLength:     100,
		IncludeImage:      false,
		ScreenshotService: "https://example.com/",
	}
	p.SetConfig(customCfg)

	cfg = p.Config()
	if cfg.Enabled {
		t.Error("Expected Enabled to be false")
	}
	if cfg.PreviewLength != 100 {
		t.Errorf("Expected PreviewLength to be 100, got %d", cfg.PreviewLength)
	}
}
