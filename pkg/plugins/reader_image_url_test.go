package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
