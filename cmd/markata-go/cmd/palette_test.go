package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

func TestPaletteCloneSpecificPalette(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := t.TempDir()

	// Create a mock palettes directory with a test palette
	palettesDir := filepath.Join(tempDir, "palettes")
	if err := os.MkdirAll(palettesDir, 0o755); err != nil {
		t.Fatalf("failed to create palettes directory: %v", err)
	}

	// Create a test palette file
	testPaletteContent := `[palette]
name = "test-palette"
variant = "dark"
description = "A test palette"

[palette.colors]
text = "#ffffff"
background = "#000000"
primary = "#7c3aed"

[palette.semantic]
text-primary = "text"
bg-primary = "background"
`
	testPalettePath := filepath.Join(palettesDir, "test-palette.toml")
	if err := os.WriteFile(testPalettePath, []byte(testPaletteContent), 0o600); err != nil {
		t.Fatalf("failed to write test palette: %v", err)
	}

	// Create a loader with the test palettes directory
	loader := palettes.NewLoaderWithPaths([]string{palettesDir})

	// Load the test palette
	p, err := loader.Load("test-palette")
	if err != nil {
		t.Fatalf("failed to load test palette: %v", err)
	}

	// Verify the palette was loaded correctly
	if p.Name != "test-palette" {
		t.Errorf("expected name 'test-palette', got %q", p.Name)
	}

	if p.Variant != palettes.VariantDark {
		t.Errorf("expected variant 'dark', got %q", p.Variant)
	}

	// Test cloning
	cloned := p.Clone()
	cloned.Name = "cloned-palette"
	cloned.Description = "Cloned from test-palette"

	if cloned.Name != "cloned-palette" {
		t.Errorf("expected cloned name 'cloned-palette', got %q", cloned.Name)
	}

	// Verify colors are preserved
	if cloned.Colors["text"] != "#ffffff" {
		t.Errorf("expected text color '#ffffff', got %q", cloned.Colors["text"])
	}

	if cloned.Colors["background"] != "#000000" {
		t.Errorf("expected background color '#000000', got %q", cloned.Colors["background"])
	}

	// Verify semantic colors are preserved
	if cloned.Semantic["text-primary"] != "text" {
		t.Errorf("expected semantic text-primary 'text', got %q", cloned.Semantic["text-primary"])
	}
}

func TestPaletteCloneNonExistentPalette(t *testing.T) {
	// Create a temporary directory with no palettes
	tempDir := t.TempDir()
	palettesDir := filepath.Join(tempDir, "palettes")
	if err := os.MkdirAll(palettesDir, 0o755); err != nil {
		t.Fatalf("failed to create palettes directory: %v", err)
	}

	// Create a loader with the empty palettes directory (no built-in fallback)
	loader := palettes.NewLoaderWithPaths([]string{palettesDir})

	// Try to load a non-existent palette
	_, err := loader.Load("non-existent-palette")
	if err == nil {
		t.Error("expected error when loading non-existent palette, got nil")
	}
}

func TestGeneratePaletteTOML(t *testing.T) {
	// Create a test palette
	p := palettes.NewPalette("my-test-theme", palettes.VariantDark)
	p.Description = "A test theme"
	p.Colors["text"] = "#e0e0e0"
	p.Colors["background"] = "#1e1e1e"
	p.Semantic["text-primary"] = "text"

	// Generate TOML
	tomlContent := generatePaletteTOML(p)

	// Verify key content is present
	if !contains(tomlContent, `name = "my-test-theme"`) {
		t.Error("TOML should contain palette name")
	}

	if !contains(tomlContent, `variant = "dark"`) {
		t.Error("TOML should contain variant")
	}

	if !contains(tomlContent, `description = "A test theme"`) {
		t.Error("TOML should contain description")
	}

	if !contains(tomlContent, `text = "#e0e0e0"`) {
		t.Error("TOML should contain text color")
	}

	if !contains(tomlContent, `background = "#1e1e1e"`) {
		t.Error("TOML should contain background color")
	}
}

func TestFormatPalettePreview(t *testing.T) {
	info := palettes.PaletteInfo{
		Name:        "test-palette",
		Variant:     palettes.VariantDark,
		Author:      "Test Author",
		Description: "A test description",
		Source:      "built-in",
	}

	preview := formatPalettePreview(info)

	if !contains(preview, "Name: test-palette") {
		t.Error("preview should contain palette name")
	}

	if !contains(preview, "Variant: dark") {
		t.Error("preview should contain variant")
	}

	if !contains(preview, "Author: Test Author") {
		t.Error("preview should contain author")
	}

	if !contains(preview, "Description: A test description") {
		t.Error("preview should contain description")
	}

	if !contains(preview, "Source: built-in") {
		t.Error("preview should contain source")
	}
}

func TestNormalizeFileName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"My Theme", "my-theme"},
		{"Catppuccin Mocha", "catppuccin-mocha"},
		{"test_theme", "test-theme"},
		{"already-normalized", "already-normalized"},
		{"UPPERCASE", "uppercase"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeFileName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeFileName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// contains checks if a string contains a substring.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
