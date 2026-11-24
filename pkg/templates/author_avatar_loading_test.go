package templates

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestBylineAvatars_VisibleEagerCollapsedLazy(t *testing.T) {
	authorRole, editorRole := "author", "copy editor"
	writerAvatar, editorAvatar, secondAvatar := "/writer.png", "/editor.png", "/second.png"

	for _, directory := range []string{"", "../../templates"} {
		t.Run(directory, func(t *testing.T) {
			engine, err := NewEngine(directory)
			if err != nil {
				t.Fatal(err)
			}

			render := func(editors int) string {
				t.Helper()
				authors := []models.Author{
					{ID: "writer", Name: "Writer", Role: &authorRole, Avatar: &writerAvatar},
					{ID: "editor", Name: "Editor", Role: &editorRole, Avatar: &editorAvatar},
				}
				if editors == 2 {
					authors = append(authors, models.Author{ID: "second", Name: "Second Editor", Role: &editorRole, Avatar: &secondAvatar})
				}
				post := &models.Post{AuthorObjects: authors}
				got, err := engine.Render("components/post_byline.html", NewContext(post, "", &models.Config{}))
				if err != nil {
					t.Fatalf("render byline: %v", err)
				}
				return got
			}

			one := render(1)
			if strings.Count(one, `loading="eager"`) != 2 || strings.Contains(one, `loading="lazy"`) {
				t.Fatalf("visible writer and single editor avatars should be eager: %s", one)
			}

			multiple := render(2)
			// One visible writer avatar and two visible editor-stack thumbnails.
			if strings.Count(multiple, `loading="eager"`) != 3 {
				t.Fatalf("visible writer and editor-stack avatars should be eager: %s", multiple)
			}
			// The two full-size editor avatars are hidden until <details> opens.
			if strings.Count(multiple, `loading="lazy"`) != 2 {
				t.Fatalf("collapsed editor-list avatars should remain lazy: %s", multiple)
			}
		})
	}
}
