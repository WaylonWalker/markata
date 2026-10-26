package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// Asset loader code should only be emitted on pages with Lite YouTube markup,
// even when another page enables the feature globally.
func TestLiteYouTubeAssetsPageScopedRendering(t *testing.T) {
	m := lifecycle.NewManager()
	m.Config().Extra["templates_dir"] = "/nonexistent"
	m.Config().Extra["lite_youtube_enabled"] = true
	title := "Video"
	plainTitle := "Plain article"

	video := &models.Post{
		Path:        "video.md",
		Slug:        "video",
		Title:       &title,
		Template:    "post.html",
		ArticleHTML: `<lite-youtube videoid="dQw4w9WgXcQ"></lite-youtube>`,
		Extra:       map[string]interface{}{"needs_lite_youtube": true},
	}
	plain := &models.Post{
		Path:        "plain.md",
		Slug:        "plain",
		Title:       &plainTitle,
		Template:    "post.html",
		ArticleHTML: "<p>Only prose.</p>",
	}
	// A generated/handwritten HTML embed should work without the flag.
	raw := &models.Post{
		Path:        "raw.md",
		Slug:        "raw",
		Title:       &title,
		Template:    "post.html",
		ArticleHTML: `<lite-youtube videoid="dQw4w9WgXcQ"></lite-youtube>`,
	}
	m.SetPosts([]*models.Post{video, plain, raw})
	plugin := NewTemplatesPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Render(m); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		post    *models.Post
		enabled bool
	}{
		{post: video, enabled: true},
		{post: plain, enabled: false},
		{post: raw, enabled: true},
	} {
		hasLoader := strings.Contains(tc.post.HTML, "lite-yt-embed.js") &&
			strings.Contains(tc.post.HTML, "lite-yt-embed.css")
		if hasLoader != tc.enabled {
			t.Errorf("post %s loader = %v, want %v", tc.post.Slug, hasLoader, tc.enabled)
		}
		hasTransitionHints := strings.Contains(tc.post.HTML, `data-markata-lite-youtube-css="`) &&
			strings.Contains(tc.post.HTML, `data-markata-lite-youtube-js="`)
		if hasTransitionHints != tc.enabled {
			t.Errorf("post %s transition asset hints = %v, want %v", tc.post.Slug, hasTransitionHints, tc.enabled)
		}
	}
}

// Hand-authored Lite YouTube markup must load even when the transform phase
// hasn't set the site-wide capability flag (e.g. content made after Transform).
func TestLiteYouTubeAssetsWithoutGlobalFlag(t *testing.T) {
	m := lifecycle.NewManager()
	m.Config().Extra["templates_dir"] = "/nonexistent"
	m.Config().Extra["lite_youtube_enabled"] = false
	title := "Video"
	video := &models.Post{
		Path:        "raw.md",
		Slug:        "raw",
		Title:       &title,
		Template:    "post.html",
		ArticleHTML: `<lite-youtube videoid="dQw4w9WgXcQ"></lite-youtube>`,
	}
	m.SetPosts([]*models.Post{video})
	plugin := NewTemplatesPlugin()
	if err := plugin.Configure(m); err != nil {
		t.Fatal(err)
	}
	if err := plugin.Render(m); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(video.HTML, "lite-yt-embed.js") || !strings.Contains(video.HTML, "lite-yt-embed.css") {
		t.Fatal("hand-authored Lite YouTube must load even without the global feature flag")
	}
}
