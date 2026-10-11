package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/models"
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
		Path: "video.md", Slug: "video", Title: &title,
		Template: "post.html",
		ArticleHTML: `<lite-youtube videoid="dQw4w9WgXcQ"></lite-youtube>`,
		Extra: map[string]interface{}{"needs_lite_youtube": true},
	}
	plain := &models.Post{
		Path: "plain.md", Slug: "plain", Title: &plainTitle,
		Template: "post.html",
		ArticleHTML: "<p>Only prose.</p>",
	}
	// A generated/handwritten HTML embed should work without the flag.
	raw := &models.Post{
		Path: "raw.md", Slug: "raw", Title: &title,
		Template: "post.html",
		ArticleHTML: `<lite-youtube videoid="dQw4w9WgXcQ"></lite-youtube>`,
	}
	m.SetPosts([]*models.Post{video, plain, raw})
	if err := NewTemplatesPlugin().Render(m); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		post    *models.Post
		enabled bool
	}{
		{video, true}, {plain, false}, {raw, true},
	} {
		hasLoader := strings.Contains(tc.post.HTML, "lite-yt-embed.js") &&
			strings.Contains(tc.post.HTML, "lite-yt-embed.css")
		if hasLoader != tc.enabled {
			t.Errorf("post %s loader = %v, want %v", tc.post.Slug, hasLoader, tc.enabled)
		}
	}
}
