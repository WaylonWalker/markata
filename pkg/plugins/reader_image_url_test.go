package plugins

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

func TestNormalizeReaderImageURL(t *testing.T) {
	tests := []struct {
		name, image, article, want string
	}{
		{"absolute image", "https://cdn.example.com/image.png", "https://example.com/post", "https://cdn.example.com/image.png"},
		{"relative image", "/images/post.jpg", "https://example.com/articles/post", "https://example.com/images/post.jpg"},
		{"relative to article directory", "photo.png", "https://example.com/articles/post", "https://example.com/articles/photo.png"},
		{"protocol relative", "//cdn.example.com/photo.jpg", "https://example.com/post", "https://cdn.example.com/photo.jpg"},
		{"unsafe scheme", "javascript:alert(1)", "https://example.com/post", ""},
		{"data scheme", "data:image/png;base64,abcd", "https://example.com/post", ""},
		{"relative without base", "/pic.png", "", ""},
		{"non-http article", "photo.png", "mailto:hello@example.com", ""},
		{"empty image", "", "https://example.com/post", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeReaderImageURL(tt.image, tt.article); got != tt.want {
				t.Errorf("normalizeReaderImageURL(%q, %q) = %q, want %q", tt.image, tt.article, got, tt.want)
			}
		})
	}
}

func TestReaderImageFallbackScriptIncluded(t *testing.T) {
	for _, path := range []string{
		filepath.Join("..", "..", "templates", "reader.html"),
		filepath.Join("..", "themes", "default", "templates", "reader.html"),
	} {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(content), "reader-image-fallback.js") {
			t.Errorf("%s missing image fallback script", path)
		}
	}
	path := filepath.Join("..", "themes", "default", "static", "js", "reader-image-fallback.js")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "addEventListener('error'") ||
		!strings.Contains(string(content), "reader-entry-source-tile") {
		t.Fatal("reader image fallback must capture failed previews and render source tiles")
	}
}

func TestReaderImageSrcset(t *testing.T) {
	const image = "https://storage.ghost.io/c/site/content/images/size/w1200/photo.png?token=a%2Fb&v=2"
	for _, tc := range []struct {
		name, image string
		widths      []int
	}{
		{"Ghost raster", image, []int{320, 480, 640, 768, 960, 1200}},
		{"existing width caps candidates", strings.ReplaceAll(image, "w1200", "w800"), []int{320, 480, 640, 768, 800}},
		{"escaped path preserved", strings.ReplaceAll(image, "photo.png", "a%2Fb%20photo.png"), []int{320, 480, 640, 768, 960, 1200}},
		{"unknown provider", strings.ReplaceAll(image, "storage.ghost.io", "example.com"), nil},
		{"lookalike host", strings.ReplaceAll(image, "storage.ghost.io", "storage.ghost.io.example.com"), nil},
		{"credentials", strings.ReplaceAll(image, "https://", "https://user@"), nil},
		{"HTTP unchanged", strings.ReplaceAll(image, "https://", "http://"), nil},
		{"animated GIF unchanged", strings.ReplaceAll(image, ".png", ".gif"), nil},
		{"vector unchanged", strings.ReplaceAll(image, ".png", ".svg"), nil},
		{"non-size path unchanged", strings.ReplaceAll(image, "/size/w1200", ""), nil},
		{"small original unchanged", strings.ReplaceAll(image, "w1200", "w320"), nil},
		{"invalid width", strings.ReplaceAll(image, "w1200", "w0"), nil},
		{"overflow width", strings.ReplaceAll(image, "w1200", "w99999999999999999999999999999999"), nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := readerImageSrcset(tc.image)
			if len(tc.widths) == 0 {
				if got != "" {
					t.Fatalf("unsupported preview gained resizing URLs: %q", got)
				}
				return
			}
			parts := strings.Split(got, ", ")
			if len(parts) != len(tc.widths) {
				t.Fatalf("candidates = %q", got)
			}
			for i, width := range tc.widths {
				originalWidth := tc.widths[len(tc.widths)-1]
				want := strings.Replace(tc.image, "/size/w"+strconv.Itoa(originalWidth)+"/", "/size/w"+strconv.Itoa(width)+"/", 1) + " " + strconv.Itoa(width) + "w"
				if parts[i] != want {
					t.Errorf("candidate = %q, want %q", parts[i], want)
				}
			}
		})
	}
}

func TestReaderResponsiveImageRendering(t *testing.T) {
	plugin := NewBlogrollPlugin()
	image := "https://storage.ghost.io/c/site/content/images/size/w1200/photo.png?x=1&y=2"
	entry := plugin.readerEntryMap(&models.ExternalEntry{Title: "Story", URL: "https://example.com/story", ImageURL: image}, nil, models.BlogrollConfig{})
	if entry["image_url"] != image || entry["image_srcset"] == "" {
		t.Fatalf("reader must retain original src and offer candidates: %v", entry)
	}
	engine, err := templates.NewEngine("")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range []string{"../../templates/reader.html", "../themes/default/templates/reader.html"} {
		data, readErr := os.ReadFile(file)
		if readErr != nil {
			t.Fatal(readErr)
		}
		found := false
		for _, line := range strings.Split(string(data), "\n") {
			if !strings.Contains(line, `class="reader-entry-image `) {
				continue
			}
			found = true
			markup, renderErr := engine.RenderString(line, templates.Context{Extra: map[string]interface{}{"entry": entry}})
			if renderErr != nil {
				t.Fatal(renderErr)
			}
			assertReaderResponsiveMarkup(t, markup)
		}
		if !found {
			t.Fatalf("preview image missing from %s", file)
		}
	}
	markup := plugin.renderReaderTimeline([]map[string]interface{}{{"entries": []map[string]interface{}{entry}}}, false)
	assertReaderResponsiveMarkup(t, markup)

	entry["image_srcset"] = ""
	markup = plugin.renderReaderTimeline([]map[string]interface{}{{"entries": []map[string]interface{}{entry}}}, false)
	if strings.Contains(markup, "srcset=") || strings.Contains(markup, "sizes=") {
		t.Fatal("unsupported image should keep its original src without responsive attributes")
	}
}

func assertReaderResponsiveMarkup(t *testing.T, markup string) {
	t.Helper()
	for _, want := range []string{`src="https://storage.ghost.io/c/site/content/images/size/w1200/photo.png?x=1&amp;y=2"`, `srcset="https://storage.ghost.io/c/site/content/images/size/w320/`, `sizes="100vw"`, `loading="lazy" decoding="async"`, `width="960" height="540"`} {
		if !strings.Contains(markup, want) {
			t.Errorf("responsive image markup missing %q: %s", want, markup)
		}
	}
}
