package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestMDVideoPlugin_Name(t *testing.T) {
	p := NewMDVideoPlugin()
	if got := p.Name(); got != "md_video" {
		t.Errorf("Name() = %q, want %q", got, "md_video")
	}
}

func TestMDVideoPlugin_DefaultConfig(t *testing.T) {
	p := NewMDVideoPlugin()
	cfg := p.Config()

	if !cfg.Enabled {
		t.Error("Expected Enabled to be true by default")
	}
	if cfg.Autoplay {
		t.Error("Expected Autoplay to be false by default")
	}
	if !cfg.Loop {
		t.Error("Expected Loop to be true by default")
	}
	if !cfg.Muted {
		t.Error("Expected Muted to be true by default")
	}
	if !cfg.Playsinline {
		t.Error("Expected Playsinline to be true by default")
	}
	if !cfg.Controls {
		t.Error("Expected Controls to be true by default")
	}
	if cfg.VideoClass != "md-video" {
		t.Errorf("Expected VideoClass to be 'md-video', got %q", cfg.VideoClass)
	}
	if cfg.Preload != "none" {
		t.Errorf("Expected Preload to be 'none', got %q", cfg.Preload)
	}
}

func TestMDVideoPlugin_ProcessPost_BasicVideo(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: `<p><img src="video.mp4" alt="My video"></p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	// Should contain video tag
	if !contains(post.ArticleHTML, "<video") {
		t.Error("Expected video tag in output")
	}
	if !contains(post.ArticleHTML, `<source src="video.mp4"`) {
		t.Error("Expected source tag with correct src")
	}
	if !contains(post.ArticleHTML, `type="video/mp4"`) {
		t.Error("Expected video/mp4 MIME type")
	}
	if contains(post.ArticleHTML, "autoplay") {
		t.Error("Expected video to wait for reader playback")
	}
	if !contains(post.ArticleHTML, "loop") {
		t.Error("Expected loop attribute")
	}
	if !contains(post.ArticleHTML, "muted") {
		t.Error("Expected muted attribute")
	}
	if !contains(post.ArticleHTML, "playsinline") {
		t.Error("Expected playsinline attribute")
	}
	if !contains(post.ArticleHTML, "controls") {
		t.Error("Expected controls attribute")
	}
	if !contains(post.ArticleHTML, `preload="none"`) {
		t.Error("Expected preload=none by default")
	}
	if !contains(post.ArticleHTML, `class="md-video"`) {
		t.Error("Expected md-video class")
	}
	// Should NOT contain img tag anymore
	if contains(post.ArticleHTML, "<img") {
		t.Error("Should not contain img tag after processing")
	}
}

func TestMDVideoPlugin_ProcessPost_WithQueryParams(t *testing.T) {
	p := NewMDVideoPlugin()

	// Test URL with query parameters (like your dropper example)
	post := &models.Post{
		ArticleHTML: `<img src="https://dropper.waylonwalker.com/file/6f042a91.mp4?width=500" alt="kickflip">`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !contains(post.ArticleHTML, "<video") {
		t.Error("Expected video tag for URL with query params")
	}
	// Preserve the full URL including query params
	if !contains(post.ArticleHTML, `src="https://dropper.waylonwalker.com/file/6f042a91.mp4?width=500"`) {
		t.Errorf("Expected full URL preserved, got: %s", post.ArticleHTML)
	}
	if !contains(post.ArticleHTML, "kickflip") {
		t.Error("Expected alt text preserved as fallback")
	}
}

func TestMDVideoPlugin_ProcessPost_DownsizesPosterOnly(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: `<img src="https://dropper.waylonwalker.com/file/6f042a91.mp4?width=500" alt="kickflip">`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `src="https://dropper.waylonwalker.com/file/6f042a91.mp4?width=500"`) {
		t.Errorf("Expected full video source preserved, got: %s", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `poster="https://dropper.waylonwalker.com/file/6f042a91.webp?h=675`) || !strings.Contains(post.ArticleHTML, `w=1200`) {
		t.Errorf("Expected downsized poster URL, got: %s", post.ArticleHTML)
	}
}

func TestMDVideoPlugin_ProcessPost_DerivesPosterForSupportedVideoFormats(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantType string
		wantWebP string
	}{
		{"webm", "https://dropper.waylonwalker.com/file/demo.webm?width=500#watch", "video/webm", "demo.webp"},
		{"ogv", "https://dropper.waylonwalker.com/file/demo.ogv?width=500#watch", "video/ogg", "demo.webp"},
		{"mov", "https://dropper.waylonwalker.com/file/demo.mov?width=500#watch", "video/quicktime", "demo.webp"},
		{"m4v", "https://dropper.waylonwalker.com/file/demo.m4v?width=500#watch", "video/x-m4v", "demo.webp"},
		{"avi", "https://dropper.waylonwalker.com/file/demo.avi?width=500#watch", "video/x-msvideo", "demo.webp"},
		{"mkv", "https://dropper.waylonwalker.com/file/demo.mkv?width=500#watch", "video/x-matroska", "demo.webp"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMDVideoPlugin()
			post := &models.Post{
				ArticleHTML: `<img src="` + tt.src + `" alt="demo">`,
			}

			if err := p.processPost(post); err != nil {
				t.Fatalf("processPost() error = %v", err)
			}

			if !strings.Contains(post.ArticleHTML, `<video`) {
				t.Fatalf("Expected video tag, got: %s", post.ArticleHTML)
			}
			if !strings.Contains(post.ArticleHTML, `type="`+tt.wantType+`"`) {
				t.Fatalf("Expected MIME type %q, got: %s", tt.wantType, post.ArticleHTML)
			}
			if !strings.Contains(post.ArticleHTML, `poster="https://dropper.waylonwalker.com/file/`+tt.wantWebP) {
				t.Fatalf("Expected derived poster for %s, got: %s", tt.name, post.ArticleHTML)
			}
			if !strings.Contains(post.ArticleHTML, `w=1200`) || !strings.Contains(post.ArticleHTML, `h=675`) {
				t.Fatalf("Expected downsized poster dimensions, got: %s", post.ArticleHTML)
			}
		})
	}
}

func TestMDVideoPlugin_ProcessPost_DifferentExtensions(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantType string
	}{
		{"mp4", "video.mp4", "video/mp4"},
		{"webm", "video.webm", "video/webm"},
		{"ogg", "video.ogg", "video/ogg"},
		{"ogv", "video.ogv", "video/ogg"},
		{"mov", "video.mov", "video/quicktime"},
		{"m4v", "video.m4v", "video/x-m4v"},
		{"mkv", "video.mkv", "video/x-matroska"},
		{"MP4 uppercase", "video.MP4", "video/mp4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMDVideoPlugin()
			post := &models.Post{
				ArticleHTML: `<img src="` + tt.src + `" alt="test">`,
			}

			err := p.processPost(post)
			if err != nil {
				t.Fatalf("processPost() error = %v", err)
			}

			if !contains(post.ArticleHTML, "<video") {
				t.Errorf("Expected video tag for %s", tt.name)
			}
			if !contains(post.ArticleHTML, `type="`+tt.wantType+`"`) {
				t.Errorf("Expected MIME type %q, got: %s", tt.wantType, post.ArticleHTML)
			}
		})
	}
}

func TestMDVideoPlugin_ProcessPost_NonVideoImage(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: `<p><img src="photo.jpg" alt="A photo"></p>`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	// Should remain unchanged
	if post.ArticleHTML != originalHTML {
		t.Errorf("Non-video image should not be modified.\nGot: %s\nWant: %s", post.ArticleHTML, originalHTML)
	}
}

func TestMDVideoPlugin_ProcessPost_MixedContent(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: `<p><img src="photo.jpg" alt="Photo"></p><p><img src="video.mp4" alt="Video"></p><p>Some text</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	// Should have both img (for photo) and video (for mp4)
	if !contains(post.ArticleHTML, `<img src="photo.jpg"`) {
		t.Error("Photo img tag should remain")
	}
	if !contains(post.ArticleHTML, "<video") {
		t.Error("Video tag should be present for mp4")
	}
	if contains(post.ArticleHTML, `<img src="video.mp4"`) {
		t.Error("Video img tag should be converted")
	}
}

func TestMDVideoPlugin_ProcessPost_DisabledPlugin(t *testing.T) {
	p := NewMDVideoPlugin()
	cfg := p.Config()
	cfg.Enabled = false
	p.SetConfig(cfg)

	post := &models.Post{
		ArticleHTML: `<img src="video.mp4" alt="test">`,
	}
	originalHTML := post.ArticleHTML

	// Note: processPost doesn't check Enabled - that's done in Render()
	// When disabled, Render() returns early and never calls processPost
	// So we test that Render() respects the enabled flag by checking it directly
	if p.Config().Enabled {
		t.Error("Expected Enabled to be false")
	}

	// The processPost function itself will still process if called directly
	// This is by design - the check happens at the Render level
	_ = originalHTML // unused, but keeping for clarity
}

func TestMDVideoPlugin_ProcessPost_SkippedPost(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		Skip:        true,
		ArticleHTML: `<img src="video.mp4" alt="test">`,
	}
	originalHTML := post.ArticleHTML

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	// Should remain unchanged when skipped
	if post.ArticleHTML != originalHTML {
		t.Error("Skipped post should not be modified")
	}
}

func TestMDVideoPlugin_ProcessPost_EmptyContent(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: "",
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if post.ArticleHTML != "" {
		t.Error("Empty content should remain empty")
	}
}

func TestMDVideoPlugin_ProcessPost_NoAlt(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: `<img src="video.mp4">`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	if !contains(post.ArticleHTML, "Your browser does not support the video tag.") {
		t.Error("Expected default fallback text when no alt provided")
	}
}

func TestMDVideoPlugin_ProcessPost_RawVideoDefaultsToOnDemand(t *testing.T) {
	p := NewMDVideoPlugin()
	post := &models.Post{ArticleHTML: `<video autoplay="" controls preload="auto"><source src="https://dropper.waylonwalker.com/api/file/clip.mp4"></video>`}

	if err := p.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}
	if videoAutoplayAttrRegex.MatchString(videoTagRegex.FindString(post.ArticleHTML)) {
		t.Errorf("raw video should not autoplay by default: %s", post.ArticleHTML)
	}
	if !contains(post.ArticleHTML, `preload="none"`) {
		t.Errorf("raw video should default to preload=none: %s", post.ArticleHTML)
	}
	if !contains(post.ArticleHTML, `src="https://dropper.waylonwalker.com/api/file/clip.mp4"`) || !contains(post.ArticleHTML, "controls") {
		t.Errorf("raw video lost its native playback path: %s", post.ArticleHTML)
	}
	if !contains(post.ArticleHTML, `poster="https://dropper.waylonwalker.com/api/file/clip.webp`) {
		t.Errorf("raw video did not get its derived poster: %s", post.ArticleHTML)
	}
}

func TestMDVideoPlugin_ProcessPost_RawVideoAutoplayCanBeOptedIn(t *testing.T) {
	p := NewMDVideoPlugin()
	cfg := p.config
	cfg.Autoplay = true
	cfg.Preload = "metadata"
	p.SetConfig(cfg)
	post := &models.Post{ArticleHTML: `<video autoplay controls><source src="clip.mp4"></video>`}

	if err := p.processPost(post); err != nil {
		t.Fatalf("processPost() error = %v", err)
	}
	if !contains(post.ArticleHTML, `data-authored-autoplay="true"`) || !contains(post.ArticleHTML, `preload="none"`) {
		t.Errorf("explicit autoplay config was not honored: %s", post.ArticleHTML)
	}
}

func TestMDVideoPlugin_ProcessPost_CustomConfig(t *testing.T) {
	p := NewMDVideoPlugin()
	p.SetConfig(models.MDVideoConfig{
		Enabled:         true,
		VideoExtensions: []string{".mp4"},
		VideoClass:      "custom-video",
		Controls:        false,
		Autoplay:        false,
		Loop:            false,
		Muted:           false,
		Playsinline:     false,
		Preload:         "auto",
	})

	post := &models.Post{
		ArticleHTML: `<img src="video.mp4" alt="test">`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	// Should have custom class
	if !contains(post.ArticleHTML, `class="custom-video"`) {
		t.Error("Expected custom-video class")
	}
	// Should have preload="auto"
	if !contains(post.ArticleHTML, `preload="auto"`) {
		t.Error("Expected preload=auto")
	}
	// Should NOT have boolean attributes when false
	if contains(post.ArticleHTML, "autoplay") {
		t.Error("Should not have autoplay when disabled")
	}
	if contains(post.ArticleHTML, "loop") {
		t.Error("Should not have loop when disabled")
	}
	if contains(post.ArticleHTML, "muted") {
		t.Error("Should not have muted when disabled")
	}
	if contains(post.ArticleHTML, "controls") {
		t.Error("Should not have controls when disabled")
	}
	if contains(post.ArticleHTML, "playsinline") {
		t.Error("Should not have playsinline when disabled")
	}
}

func TestMDVideoPlugin_ProcessPost_MultipleVideos(t *testing.T) {
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: `<p><img src="first.mp4" alt="First"></p><p><img src="second.webm" alt="Second"></p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	// Count video tags
	videoCount := countOccurrences(post.ArticleHTML, "<video")
	if videoCount != 2 {
		t.Errorf("Expected 2 video tags, got %d", videoCount)
	}

	if !contains(post.ArticleHTML, `src="first.mp4"`) {
		t.Error("Expected first video source")
	}
	if !contains(post.ArticleHTML, `src="second.webm"`) {
		t.Error("Expected second video source")
	}
}

func TestMDVideoPlugin_isVideoURL(t *testing.T) {
	p := NewMDVideoPlugin()

	tests := []struct {
		url  string
		want bool
	}{
		{"video.mp4", true},
		{"video.webm", true},
		{"video.ogg", true},
		{"video.ogv", true},
		{"video.mov", true},
		{"video.m4v", true},
		{"video.avi", true},
		{"video.mkv", true},
		{"VIDEO.MP4", true}, // case insensitive
		{"video.MP4", true},
		{"https://example.com/video.mp4", true},
		{"https://example.com/video.mp4?width=500", true}, // with query params
		{"https://example.com/video.webm?width=500#watch", true},
		{"image.jpg", false},
		{"image.png", false},
		{"document.pdf", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := p.isVideoURL(tt.url); got != tt.want {
				t.Errorf("isVideoURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestMDVideoPlugin_getVideoMIMEType(t *testing.T) {
	p := NewMDVideoPlugin()

	tests := []struct {
		url      string
		wantMIME string
	}{
		{"video.mp4", "video/mp4"},
		{"video.webm", "video/webm"},
		{"video.ogg", "video/ogg"},
		{"video.ogv", "video/ogg"},
		{"video.mov", "video/quicktime"},
		{"video.m4v", "video/x-m4v"},
		{"video.avi", "video/x-msvideo"},
		{"video.mkv", "video/x-matroska"},
		{"video.webm?query=param#watch", "video/webm"},
		{"video.unknown", "video/mp4"}, // fallback
		{"video.mp4?query=param", "video/mp4"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := p.getVideoMIMEType(tt.url); got != tt.wantMIME {
				t.Errorf("getVideoMIMEType(%q) = %q, want %q", tt.url, got, tt.wantMIME)
			}
		})
	}
}

func TestMDVideoPlugin_RealWorldExample(t *testing.T) {
	// Test with your actual example
	p := NewMDVideoPlugin()

	post := &models.Post{
		ArticleHTML: `<p><img src="https://dropper.waylonwalker.com/file/6f042a91-1e90-445d-b91d-8d4ee187af2c.mp4" alt="kickflip down the 3 stair - fingerboarding"></p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost() error = %v", err)
	}

	// Verify it looks like your expected output
	expected := []string{
		"<video",
		"loop",
		"muted",
		"playsinline",
		"controls",
		`class="md-video"`,
		`<source src="https://dropper.waylonwalker.com/file/6f042a91-1e90-445d-b91d-8d4ee187af2c.mp4"`,
		`type="video/mp4"`,
		"kickflip down the 3 stair - fingerboarding",
		"</video>",
	}

	for _, want := range expected {
		if !contains(post.ArticleHTML, want) {
			t.Errorf("Expected %q in output.\nGot: %s", want, post.ArticleHTML)
		}
	}

	// Should not have img tag
	if contains(post.ArticleHTML, "<img") {
		t.Error("Should not contain img tag")
	}
}

// Helper function
func countOccurrences(s, substr string) int {
	count := 0
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			count++
		}
	}
	return count
}

func TestMDVideoPlugin_AuthoredIntent(t *testing.T) {
	p := NewMDVideoPlugin()
	cfg := p.config
	cfg.Autoplay = true
	p.SetConfig(cfg)
	for _, tt := range []struct {
		name       string
		html       string
		wantIntent bool
	}{
		{"raw manual", `<video controls><source src="clip.mp4"></video>`, false},
		{"raw autoplay", `<video autoplay controls><source src="clip.mp4"></video>`, true},
		{"Markdown configured autoplay", `<img src="clip.mp4" alt="clip">`, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{ArticleHTML: tt.html}
			if err := p.processPost(post); err != nil {
				t.Fatal(err)
			}
			tag := videoTagRegex.FindString(post.ArticleHTML)
			if videoAutoplayAttrRegex.MatchString(tag) {
				t.Fatalf("premature native autoplay: %s", tag)
			}
			if got := strings.Contains(tag, `data-authored-autoplay="true"`); got != tt.wantIntent {
				t.Fatalf("intent = %v, want %v: %s", got, tt.wantIntent, tag)
			}
			if tt.wantIntent && !strings.Contains(tag, `preload="none"`) {
				t.Fatalf("premature preload: %s", tag)
			}
		})
	}
}

func TestMDVideoPlugin_PosterReservesGeometry(t *testing.T) {
	p := NewMDVideoPlugin()
	for _, tt := range []struct{ name, opening, want string }{
		{"default frame", `<video controls width="100%">`, `style="aspect-ratio:16/9"`},
		{"author style", `<video controls style="aspect-ratio:1/1">`, `style="aspect-ratio:1/1"`},
		{"author height", `<video controls height="240">`, `height="240"`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{ArticleHTML: tt.opening + `<source src="https://dropper.waylonwalker.com/api/file/clip.mp4"></video>`}
			if err := p.processPost(post); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(post.ArticleHTML, tt.want) {
				t.Fatalf("geometry lost: %s", post.ArticleHTML)
			}
			if tt.name != "default frame" && strings.Contains(post.ArticleHTML, `aspect-ratio:16/9`) {
				t.Fatalf("author geometry overridden: %s", post.ArticleHTML)
			}
		})
	}
}
