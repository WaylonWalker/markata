package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestYouTubePlugin_Name(t *testing.T) {
	p := NewYouTubePlugin()
	if got := p.Name(); got != "youtube" {
		t.Errorf("Name() = %q, want %q", got, "youtube")
	}
}

func TestYouTubePlugin_ProcessPost_NoURLs(t *testing.T) {
	p := NewYouTubePlugin()

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
		t.Error("HTML was modified when no YouTube URLs present")
	}
}

func TestYouTubePlugin_ProcessPost_StandardURL(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>Check this out:</p>
<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ</p>
<p>More content here</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should contain youtube embed
	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Error("Expected youtube-embed class in output")
	}
	if !strings.Contains(post.ArticleHTML, `src="https://www.youtube-nocookie.com/embed/dQw4w9WgXcQ"`) {
		t.Error("Expected privacy-enhanced embed URL")
	}
	if !strings.Contains(post.ArticleHTML, `allowfullscreen`) {
		t.Error("Expected allowfullscreen attribute")
	}
	// Should preserve surrounding content
	if !strings.Contains(post.ArticleHTML, "Check this out:") {
		t.Error("Lost content before URL")
	}
	if !strings.Contains(post.ArticleHTML, "More content here") {
		t.Error("Lost content after URL")
	}
}

func TestYouTubePlugin_ProcessPost_ShortURL(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://youtu.be/dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should contain youtube embed
	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Error("Expected youtube-embed class in output")
	}
	if !strings.Contains(post.ArticleHTML, `embed/dQw4w9WgXcQ`) {
		t.Error("Expected video ID in embed URL")
	}
}

func TestYouTubePlugin_ProcessPost_MobileURL(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://m.youtube.com/watch?v=dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Error("Expected youtube-embed class for mobile URL")
	}
}

func TestYouTubePlugin_ProcessPost_URLWithoutWWW(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://youtube.com/watch?v=dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Error("Expected youtube-embed class for URL without www")
	}
}

func TestYouTubePlugin_ProcessPost_URLWithExtraParams(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ&t=42</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should extract video ID and timestamp
	if !strings.Contains(post.ArticleHTML, `embed/dQw4w9WgXcQ`) {
		t.Error("Expected video ID in embed URL")
	}
	// Timestamp should be included as start parameter
	if !strings.Contains(post.ArticleHTML, `?start=42`) {
		t.Error("Expected start=42 parameter in embed URL")
	}
}

func TestYouTubePlugin_ProcessPost_InlineURL(t *testing.T) {
	p := NewYouTubePlugin()

	// URL is inline with other text, should NOT be converted
	post := &models.Post{
		ArticleHTML: "<p>Check out https://www.youtube.com/watch?v=dQw4w9WgXcQ for more info</p>",
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

func TestYouTubePlugin_ProcessPost_MultipleURLs(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ</p>
<p>https://youtu.be/xvFZjo5PgG0</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have two embeds
	count := strings.Count(post.ArticleHTML, `class="youtube-embed"`)
	if count != 2 {
		t.Errorf("Expected 2 youtube embeds, got %d", count)
	}
}

func TestYouTubePlugin_ProcessPost_SkipPost(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		Skip:        true,
		ArticleHTML: `<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ</p>`,
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

func TestYouTubePlugin_ProcessPost_EmptyHTML(t *testing.T) {
	p := NewYouTubePlugin()

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

func TestYouTubePlugin_ProcessPost_PrivacyDisabled(t *testing.T) {
	p := NewYouTubePlugin()
	p.SetConfig(models.YouTubeConfig{
		Enabled:         true,
		PrivacyEnhanced: false,
		ContainerClass:  "youtube-embed",
		LazyLoad:        true,
	})

	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should use regular youtube.com
	if !strings.Contains(post.ArticleHTML, `src="https://www.youtube.com/embed/`) {
		t.Error("Expected regular youtube.com embed URL when privacy disabled")
	}
	if strings.Contains(post.ArticleHTML, "youtube-nocookie") {
		t.Error("Should not use youtube-nocookie when privacy disabled")
	}
}

func TestYouTubePlugin_ProcessPost_CustomContainerClass(t *testing.T) {
	p := NewYouTubePlugin()
	p.SetConfig(models.YouTubeConfig{
		Enabled:         true,
		PrivacyEnhanced: true,
		ContainerClass:  "custom-video-embed",
		LazyLoad:        true,
	})

	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="custom-video-embed"`) {
		t.Error("Expected custom container class")
	}
}

func TestYouTubePlugin_ProcessPost_LazyLoadDisabled(t *testing.T) {
	p := NewYouTubePlugin()
	p.SetConfig(models.YouTubeConfig{
		Enabled:         true,
		PrivacyEnhanced: true,
		ContainerClass:  "youtube-embed",
		LazyLoad:        false,
	})

	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if strings.Contains(post.ArticleHTML, `loading="lazy"`) {
		t.Error("Should not have lazy loading when disabled")
	}
}

func TestYouTubePlugin_ProcessPost_LazyLoadEnabled(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `loading="lazy"`) {
		t.Error("Expected lazy loading by default")
	}
}

func TestYouTubePlugin_ProcessPost_Disabled(t *testing.T) {
	p := NewYouTubePlugin()
	p.SetConfig(models.YouTubeConfig{
		Enabled: false,
	})

	// Verify config is disabled
	if p.Config().Enabled {
		t.Error("Expected plugin to be disabled")
	}
}

func TestYouTubePlugin_ProcessPost_InvalidVideoID(t *testing.T) {
	p := NewYouTubePlugin()

	// Video ID must be exactly 11 characters
	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=shortID</p>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should not be converted (invalid video ID)
	if post.ArticleHTML != originalHTML {
		t.Error("Invalid video ID URL was incorrectly converted")
	}
}

func TestYouTubePlugin_ProcessPost_WhitespaceAroundURL(t *testing.T) {
	p := NewYouTubePlugin()

	// URL with whitespace in paragraph should still work
	post := &models.Post{
		ArticleHTML: `<p>  https://www.youtube.com/watch?v=dQw4w9WgXcQ  </p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Error("Expected URL with whitespace to be converted")
	}
}

func TestYouTubePlugin_ProcessPost_HTTPScheme(t *testing.T) {
	p := NewYouTubePlugin()

	// HTTP (not HTTPS) should also work
	post := &models.Post{
		ArticleHTML: `<p>http://www.youtube.com/watch?v=dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Error("Expected HTTP URL to be converted")
	}
}

func TestYouTubePlugin_ProcessPost_SpecialCharactersInVideoID(t *testing.T) {
	p := NewYouTubePlugin()

	// Video IDs can contain underscores and hyphens
	post := &models.Post{
		ArticleHTML: `<p>https://www.youtube.com/watch?v=a_B-c1D2e3F</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `embed/a_B-c1D2e3F`) {
		t.Error("Expected video ID with special characters")
	}
}

// =============================================================================
// Timestamp Parsing Tests
// =============================================================================

func TestYouTubePlugin_ParseTimestamp_Seconds(t *testing.T) {
	p := NewYouTubePlugin()

	tests := []struct {
		url      string
		expected int
	}{
		{"https://www.youtube.com/watch?v=abc12345678&t=42", 42},
		{"https://youtu.be/abc12345678?t=123", 123},
		{"https://www.youtube.com/watch?v=abc12345678&start=60", 60},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := p.parseTimestamp(tt.url)
			if got != tt.expected {
				t.Errorf("parseTimestamp(%q) = %d, want %d", tt.url, got, tt.expected)
			}
		})
	}
}

func TestYouTubePlugin_ParseTimestamp_Duration(t *testing.T) {
	p := NewYouTubePlugin()

	tests := []struct {
		url      string
		expected int
	}{
		{"https://youtu.be/abc12345678?t=1h2m3s", 3723}, // 1*3600 + 2*60 + 3
		{"https://youtu.be/abc12345678?t=2m30s", 150},   // 2*60 + 30
		{"https://youtu.be/abc12345678?t=1h", 3600},     // 1*3600
		{"https://youtu.be/abc12345678?t=5m", 300},      // 5*60
		{"https://youtu.be/abc12345678?t=45s", 45},      // 45
		{"https://youtu.be/abc12345678?t=1h30m", 5400},  // 1*3600 + 30*60
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got := p.parseTimestamp(tt.url)
			if got != tt.expected {
				t.Errorf("parseTimestamp(%q) = %d, want %d", tt.url, got, tt.expected)
			}
		})
	}
}

func TestYouTubePlugin_ParseTimestamp_NoTimestamp(t *testing.T) {
	p := NewYouTubePlugin()

	urls := []string{
		"https://www.youtube.com/watch?v=abc12345678",
		"https://youtu.be/abc12345678",
		"https://www.youtube.com/watch?v=abc12345678&list=PLxyz",
	}

	for _, u := range urls {
		t.Run(u, func(t *testing.T) {
			got := p.parseTimestamp(u)
			if got != 0 {
				t.Errorf("parseTimestamp(%q) = %d, want 0", u, got)
			}
		})
	}
}

func TestYouTubePlugin_ProcessPost_TimestampInEmbed(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://youtu.be/dQw4w9WgXcQ?t=1m30s</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have start=90 (1*60 + 30)
	if !strings.Contains(post.ArticleHTML, `?start=90`) {
		t.Error("Expected start=90 in embed URL for 1m30s timestamp")
	}
}

func TestYouTubePlugin_ProcessPost_NoTimestamp(t *testing.T) {
	p := NewYouTubePlugin()

	post := &models.Post{
		ArticleHTML: `<p>https://youtu.be/dQw4w9WgXcQ</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should NOT have any start parameter
	if strings.Contains(post.ArticleHTML, `?start=`) {
		t.Error("Should not have start parameter when no timestamp in URL")
	}
}

// =============================================================================
// Autolinked URL Tests (Linkify extension)
// =============================================================================

func TestYouTubePlugin_ProcessPost_AutolinkedURL(t *testing.T) {
	p := NewYouTubePlugin()

	// This is what goldmark produces when Linkify extension is enabled
	post := &models.Post{
		ArticleHTML: `<p><a href="https://www.youtube.com/watch?v=dQw4w9WgXcQ">https://www.youtube.com/watch?v=dQw4w9WgXcQ</a></p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should contain youtube embed
	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Errorf("Expected youtube-embed class in output, got: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `embed/dQw4w9WgXcQ`) {
		t.Error("Expected video ID in embed URL")
	}
}

func TestYouTubePlugin_ProcessPost_AutolinkedShortURL(t *testing.T) {
	p := NewYouTubePlugin()

	// youtu.be format with Linkify
	post := &models.Post{
		ArticleHTML: `<p><a href="https://youtu.be/dQw4w9WgXcQ">https://youtu.be/dQw4w9WgXcQ</a></p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Errorf("Expected youtube-embed class for autolinked youtu.be URL, got: %s", post.ArticleHTML)
	}
}

func TestYouTubePlugin_ProcessPost_AutolinkedURLWithTimestamp(t *testing.T) {
	p := NewYouTubePlugin()

	// Autolinked URL with timestamp
	post := &models.Post{
		ArticleHTML: `<p><a href="https://youtu.be/dQw4w9WgXcQ?t=1m30s">https://youtu.be/dQw4w9WgXcQ?t=1m30s</a></p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="youtube-embed"`) {
		t.Errorf("Expected youtube-embed class, got: %s", post.ArticleHTML)
	}
	// Should have start=90 (1*60 + 30)
	if !strings.Contains(post.ArticleHTML, `?start=90`) {
		t.Errorf("Expected start=90 in embed URL for 1m30s timestamp, got: %s", post.ArticleHTML)
	}
}

func TestYouTubePlugin_ProcessPost_MixedAutolinkedAndPlain(t *testing.T) {
	p := NewYouTubePlugin()

	// Mix of autolinked and plain URLs
	post := &models.Post{
		ArticleHTML: `<p><a href="https://youtu.be/dQw4w9WgXcQ">https://youtu.be/dQw4w9WgXcQ</a></p>
<p>https://youtu.be/xvFZjo5PgG0</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Should have two embeds
	count := strings.Count(post.ArticleHTML, `class="youtube-embed"`)
	if count != 2 {
		t.Errorf("Expected 2 youtube embeds, got %d. HTML: %s", count, post.ArticleHTML)
	}
}
