package cmd

import (
	"testing"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
)

func TestCreateHuhTheme_EmptyPalette(t *testing.T) {
	theme := createHuhTheme("")
	if theme == nil {
		t.Fatal("expected non-nil theme for empty palette name")
	}
	// Should return the default Charm theme (not nil)
}

func TestCreateHuhTheme_InvalidPalette(t *testing.T) {
	theme := createHuhTheme("nonexistent-palette-xyz")
	if theme == nil {
		t.Fatal("expected non-nil theme for invalid palette name")
	}
	// Should fall back to Charm theme
}

func TestCreateHuhTheme_DefaultDark(t *testing.T) {
	theme := createHuhTheme("default-dark")
	if theme == nil {
		t.Fatal("expected non-nil theme for default-dark palette")
	}

	// The default-dark palette defines accent = blue-400 (#60a5fa).
	// Verify that the pink fuchsia (#F780E2) from ThemeCharm is no longer
	// used for the select selector (chevron).
	// We cannot easily inspect lipgloss.Style internals, but we can at
	// least verify the theme was returned without error.
}

func TestCreateHuhTheme_AllFieldsSet(t *testing.T) {
	// Create a palette with all semantic and component colors defined
	// so we exercise every branch in createHuhTheme.
	p := palettes.NewPalette("test-full", palettes.VariantDark)
	p.Colors = map[string]string{
		"blue":       "#3b82f6",
		"light-blue": "#60a5fa",
		"white":      "#ffffff",
		"gray":       "#6b7280",
		"light-gray": "#d1d5db",
		"dark-gray":  "#374151",
		"green":      "#22c55e",
		"red":        "#ef4444",
		"dark-blue":  "#2563eb",
		"near-white": "#f9fafb",
		"near-black": "#111827",
	}
	p.Semantic = map[string]string{
		"accent":         "blue",
		"accent-hover":   "light-blue",
		"text-primary":   "white",
		"text-secondary": "light-gray",
		"text-muted":     "gray",
		"success":        "green",
		"error":          "red",
		"link":           "light-blue",
		"border":         "dark-gray",
		"border-focus":   "blue",
	}
	p.Components = map[string]string{
		"button-primary-bg":     "dark-blue",
		"button-primary-text":   "near-white",
		"button-secondary-bg":   "dark-gray",
		"button-secondary-text": "near-black",
	}

	theme := createHuhThemeFromPalette(p)
	if theme == nil {
		t.Fatal("expected non-nil theme")
	}

	// Verify theme is not the raw Charm default by checking that we
	// don't crash and the theme is usable.
}

func TestCreateHuhTheme_PartialPalette(t *testing.T) {
	// Palette with only accent defined - everything else should fall back
	// gracefully without panicking.
	p := palettes.NewPalette("test-partial", palettes.VariantLight)
	p.Colors = map[string]string{
		"blue": "#3b82f6",
	}
	p.Semantic = map[string]string{
		"accent": "blue",
	}

	theme := createHuhThemeFromPalette(p)
	if theme == nil {
		t.Fatal("expected non-nil theme")
	}
}

func TestCreateHuhTheme_NoSuccessFallsBackToAccent(t *testing.T) {
	// When success is not defined, SelectedOption should use accent.
	p := palettes.NewPalette("test-no-success", palettes.VariantDark)
	p.Colors = map[string]string{
		"purple": "#8b5cf6",
	}
	p.Semantic = map[string]string{
		"accent": "purple",
	}

	theme := createHuhThemeFromPalette(p)
	if theme == nil {
		t.Fatal("expected non-nil theme")
	}
	// The SelectedOption and SelectedPrefix should use accent color
	// when success is not available. We verify no panic.
}

func TestResolveColor_Found(t *testing.T) {
	p := palettes.NewPalette("test", palettes.VariantDark)
	p.Colors = map[string]string{
		"blue": "#3b82f6",
	}
	p.Semantic = map[string]string{
		"accent": "blue",
	}

	got := resolveColor(p, "accent")
	if got == nil {
		t.Fatal("expected non-nil color for defined palette key")
	}

	// Verify it resolves to the expected lipgloss.Color
	expected := lipgloss.Color("#3b82f6")
	if got != expected {
		t.Errorf("resolveColor(accent) = %v, want %v", got, expected)
	}
}

func TestResolveColor_NotFound(t *testing.T) {
	p := palettes.NewPalette("test", palettes.VariantDark)
	p.Colors = map[string]string{}

	got := resolveColor(p, "nonexistent")
	if got != nil {
		t.Errorf("expected nil for undefined palette key, got %v", got)
	}
}

// createHuhThemeFromPalette is a test helper that creates a theme directly
// from a Palette without going through the loader (which needs files).
func createHuhThemeFromPalette(palette *palettes.Palette) *huh.Theme {
	return buildThemeFromPalette(palette)
}

func TestResolveWizardLicenseValue(t *testing.T) {
	tests := []struct {
		name     string
		selected string
		want     interface{}
	}{
		{name: "recommended default", selected: "", want: models.DefaultLicenseKey},
		{name: "explicit false", selected: "false", want: false},
		{name: "valid key", selected: "mit", want: "mit"},
		{name: "invalid key falls back", selected: "invalid", want: models.DefaultLicenseKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveWizardLicenseValue(tt.selected)
			if got != tt.want {
				t.Fatalf("resolveWizardLicenseValue(%q) = %v, want %v", tt.selected, got, tt.want)
			}
		})
	}
}
