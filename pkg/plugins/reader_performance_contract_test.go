package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Keep both reader templates aligned: production sites may override the theme
// template, while new sites use the packaged theme by default.
func TestReaderImageDecodeHints(t *testing.T) {
	templates := []string{
		filepath.Join("..", "..", "templates", "reader.html"),
		filepath.Join("..", "themes", "default", "templates", "reader.html"),
	}
	for _, path := range templates {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			source := string(data)
			if count := strings.Count(source, `decoding="async"`); count != 3 {
				t.Errorf("reader needs async decoding on the source icon, preview, and tile icon: got %d", count)
			}
			if count := strings.Count(source, `loading="lazy"`); count != 3 {
				t.Errorf("reader must preserve lazy loading for all three images: got %d", count)
			}
			if !strings.Contains(source, `width="960" height="540" loading="lazy" decoding="async"`) {
				t.Error("reader preview should preserve width, height, and lazy loading")
			}
		})
	}
}

func TestReaderOffscreenCardsKeepFirstPaint(t *testing.T) {
	path := filepath.Join("..", "themes", "default", "static", "css", "components.css")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	css := string(data)
	for _, hint := range []string{
		`@supports (content-visibility: auto)`,
		`.reader-day-entries > .reader-entry:nth-child(n+5)`,
		`content-visibility: auto;`,
		`contain-intrinsic-size: auto 24rem;`,
	} {
		if !strings.Contains(css, hint) {
			t.Errorf("missing reader offscreen rendering hint %q", hint)
		}
	}
}
