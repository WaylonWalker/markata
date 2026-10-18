package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestOneLineLinkPlugin_Name(t *testing.T) {
	p := NewOneLineLinkPlugin()
	if got := p.Name(); got != "one_line_link" {
		t.Errorf("Name() = %q, want %q", got, "one_line_link")
	}
}

func TestOneLineLinkPlugin_ProcessPost_NoURLs(t *testing.T) {
	p := NewOneLineLinkPlugin()

	post := &models.Post{
		ArticleHTML: "<p>Hello world with no URLs</p>",
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged
	if post.ArticleHTML != originalHTML {
		t.Error("HTML was modified when no URLs present")
	}
}

func TestOneLineLinkPlugin_ProcessPost_InlineURL(t *testing.T) {
	p := NewOneLineLinkPlugin()

	// URL is inline with other text, should NOT be converted
	post := &models.Post{
		ArticleHTML: "<p>Check out https://example.com for more info</p>",
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged (inline URLs are not expanded)
	if post.ArticleHTML != originalHTML {
		t.Error("Inline URL was incorrectly converted")
	}
}

func TestOneLineLinkPlugin_ProcessPost_StandaloneURL(t *testing.T) {
	p := NewOneLineLinkPlugin()

	// URL alone in paragraph should be converted
	post := &models.Post{
		ArticleHTML: `<p>Check this out:</p>
<p>https://example.com/awesome-article</p>
<p>More content here</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should contain link card
	if !strings.Contains(post.ArticleHTML, `class="link-card"`) {
		t.Error("Expected link-card class in output")
	}
	if !strings.Contains(post.ArticleHTML, `href="https://example.com/awesome-article"`) {
		t.Error("Expected href to original URL")
	}
	if !strings.Contains(post.ArticleHTML, "example.com") {
		t.Error("Expected domain in link card")
	}
	if !strings.Contains(post.ArticleHTML, `aria-label="Link on example.com"`) {
		t.Error("Expected accessible aria-label on link-domain anchor")
	}
	// Should preserve surrounding content
	if !strings.Contains(post.ArticleHTML, "Check this out:") {
		t.Error("Lost content before URL")
	}
	if !strings.Contains(post.ArticleHTML, "More content here") {
		t.Error("Lost content after URL")
	}
}

func TestOneLineLinkPlugin_ProcessPost_MultipleURLs(t *testing.T) {
	p := NewOneLineLinkPlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://example.com</p>
<p>https://other.org/page</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have two link cards
	count := strings.Count(post.ArticleHTML, `class="link-card"`)
	if count != 2 {
		t.Errorf("Expected 2 link cards, got %d", count)
	}
}

func TestOneLineLinkPlugin_ProcessPost_SkipPost(t *testing.T) {
	p := NewOneLineLinkPlugin()

	post := &models.Post{
		Skip:        true,
		ArticleHTML: `<p>https://example.com</p>`,
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

func TestOneLineLinkPlugin_ProcessPost_EmptyHTML(t *testing.T) {
	p := NewOneLineLinkPlugin()

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

func TestOneLineLinkPlugin_ProcessPost_ExcludePattern(t *testing.T) {
	p := NewOneLineLinkPlugin()
	p.SetConfig(models.OneLineLinkConfig{
		Enabled:         true,
		CardClass:       "link-card",
		ExcludePatterns: []string{`^https://twitter\.com`, `^https://x\.com`},
	})

	// Twitter URL should not be converted
	post := &models.Post{
		ArticleHTML: `<p>https://twitter.com/user/status/123</p>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// HTML should be unchanged (excluded pattern)
	if post.ArticleHTML != originalHTML {
		t.Error("Excluded URL was converted")
	}
}

func TestOneLineLinkPlugin_ProcessPost_CustomCardClass(t *testing.T) {
	p := NewOneLineLinkPlugin()
	p.SetConfig(models.OneLineLinkConfig{
		Enabled:       true,
		CardClass:     "custom-preview",
		FallbackTitle: "Visit",
	})

	post := &models.Post{
		ArticleHTML: `<p>https://example.com</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="custom-preview"`) {
		t.Error("Expected custom card class")
	}
	if !strings.Contains(post.ArticleHTML, "Visit") {
		t.Error("Expected custom fallback title")
	}
}

func TestOneLineLinkPlugin_ProcessPost_HTTPScheme(t *testing.T) {
	p := NewOneLineLinkPlugin()

	// HTTP (not HTTPS) should also work
	post := &models.Post{
		ArticleHTML: `<p>http://example.com</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `href="http://example.com"`) {
		t.Error("Expected HTTP URL to be converted")
	}
}

func TestOneLineLinkPlugin_ProcessPost_URLWithQueryParams(t *testing.T) {
	p := NewOneLineLinkPlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://example.com/page?foo=bar&amp;baz=qux</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should handle query params
	if !strings.Contains(post.ArticleHTML, `class="link-card"`) {
		t.Error("Expected link card for URL with query params")
	}
}

func TestOneLineLinkPlugin_ProcessPost_URLWithPath(t *testing.T) {
	p := NewOneLineLinkPlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://example.com/path/to/page</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `href="https://example.com/path/to/page"`) {
		t.Error("Expected URL path to be preserved")
	}
}

func TestOneLineLinkPlugin_ProcessPost_StripWWW(t *testing.T) {
	p := NewOneLineLinkPlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://www.example.com</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Domain display should strip www.
	if !strings.Contains(post.ArticleHTML, `>example.com<`) {
		t.Error("Expected www. to be stripped from display domain")
	}
}

func TestOneLineLinkPlugin_ProcessPost_WhitespaceAroundURL(t *testing.T) {
	p := NewOneLineLinkPlugin()

	// URL with whitespace in paragraph should still work
	post := &models.Post{
		ArticleHTML: `<p>  https://example.com  </p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="link-card"`) {
		t.Error("Expected URL with whitespace to be converted")
	}
}

func TestOneLineLinkPlugin_ProcessPost_Disabled(t *testing.T) {
	p := NewOneLineLinkPlugin()
	p.SetConfig(models.OneLineLinkConfig{
		Enabled: false,
	})

	// Verify config is disabled
	if p.Config().Enabled {
		t.Error("Expected plugin to be disabled")
	}
}
