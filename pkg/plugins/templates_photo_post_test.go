package plugins

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestTemplatesPlugin_RenderPhotoPostAliases(t *testing.T) {
	const image = "https://dropper.example.com/file/skate-shelf.webp"
	const caption = "<p>Skate shelf in action.</p>"
	for _, tt := range []struct {
		name     string
		template string
		body     string
		image    string
		private  bool
		want     bool
	}{
		{name: "editor photo", template: "photo", body: caption, image: image, want: true},
		{name: "legacy shot", template: "shot", body: caption, image: image, want: true},
		{name: "shots alias", template: "shots", body: caption, image: image, want: true},
		{name: "image alias", template: "image", body: caption, image: image, want: true},
		{name: "gallery retains its own layout", template: "gallery", body: caption, image: image},
		{name: "already inline", template: "shot", body: "<img src=\"" + image + "\">" + caption, image: image},
		{name: "no media", template: "photo", body: caption},
		{name: "private", template: "photo", body: caption, image: image, private: true},
		{name: "ordinary article", template: "post.html", body: caption, image: image},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTemplatesPlugin()
			m := lifecycle.NewManager()
			m.Config().Extra["templates_dir"] = filepath.Join(t.TempDir(), "missing")
			if err := p.Configure(m); err != nil {
				t.Fatal(err)
			}
			title := "Skate Shelf"
			post := &models.Post{
				Title:       &title,
				Slug:        "skate-shelf",
				Href:        "/skate-shelf/",
				Template:    tt.template,
				ArticleHTML: tt.body,
				Published:   true,
				Private:     tt.private,
				Extra:       map[string]interface{}{"image": tt.image},
			}
			m.AddPost(post)
			if err := p.Render(m); err != nil {
				t.Fatal(err)
			}
			const marker = "class=\"u-photo\""
			if got := strings.Contains(post.HTML, marker); got != tt.want {
				t.Errorf("primary photo present = %v, want %v", got, tt.want)
			}
			if !strings.Contains(post.HTML, caption) {
				t.Error("caption missing")
			}
			if tt.want {
				if !strings.Contains(post.HTML, image) {
					t.Error("photo source missing")
				}
				if strings.Index(post.HTML, marker) > strings.Index(post.HTML, caption) {
					t.Error("photo must precede caption")
				}
			}
		})
	}
}
