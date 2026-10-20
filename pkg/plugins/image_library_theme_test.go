package plugins

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestImageLibraryThemeContract guards the shared, themed shell separately
// from the server-rendered pagination/search behavior tested elsewhere.
func TestImageLibraryThemeContract(t *testing.T) {
	rootPath := filepath.Join("..", "..", "templates", "images.html")
	embeddedPath := filepath.Join("..", "themes", "default", "templates", "images.html")
	root, err := os.ReadFile(rootPath)
	if err != nil {
		t.Fatal(err)
	}
	embedded, err := os.ReadFile(embeddedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(root, embedded) {
		t.Fatal("root and embedded image-library templates must remain identical")
	}
	html := string(embedded)
	for _, contract := range []string{
		"{% extends \"base.html\" %}",
		"css/image-library.css",
		"data-image-library",
		"data-image-grid",
		"data-image-card",
		"data-image-filter",
		"data-image-sort",
		"data-image-pagination",
		"data-copy-value",
	} {
		if !strings.Contains(html, contract) {
			t.Errorf("image-library template missing %q", contract)
		}
	}
	if strings.Contains(html, "image-library__eyebrow") {
		t.Error("image-library template reintroduced a separate asset-index eyebrow")
	}

	stylesPath := filepath.Join("..", "themes", "default", "static", "css", "image-library.css")
	styles, err := os.ReadFile(stylesPath)
	if err != nil {
		t.Fatal(err)
	}
	css := string(styles)
	// The library must honor user-selected palettes, fontpacks, text sizes and
	// aesthetics, rather than inventing its own fixed colors, sizes and radii.
	for _, token := range []string{
		"--color-background",
		"--color-surface",
		"--color-text",
		"--color-text-muted",
		"--color-primary",
		"--color-link",
		"--color-border",
		"--font-body",
		"--font-heading",
		"--font-display",
		"--text-sm",
		"--space-4",
		"--spacing-scale",
		"--radius-sm",
		"--radius-lg",
		"--border-width-normal",
		"--shadow-sm",
		"--grid-max-width",
		".image-card[hidden]",
		".image-library__pagination[hidden]",
		"content-visibility: auto",
		"prefers-reduced-motion: reduce",
	} {
		if !strings.Contains(css, token) {
			t.Errorf("image-library styles missing shared site token/contract %q", token)
		}
	}
	if strings.Contains(css, "translateY(") {
		t.Error("image cards should not lift and shift on hover")
	}
}
