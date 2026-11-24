package aesthetic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/themes"
)

// Built-in presets feed the generated stylesheet; the embedded default CSS is
// also served by the theme and must not silently drift from those presets.
func TestBuiltinAestheticTokensMatchThemeCSS(t *testing.T) {
	css, err := themes.ReadStatic("css/aesthetic.css")
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"balanced", "brutal", "precision", "elevated", "minimal"} {
		t.Run(name, func(t *testing.T) {
			a, err := LoadBuiltin(name)
			if err != nil {
				t.Fatal(err)
			}
			block := regexp.MustCompile(`(?s)\[data-aesthetic="` + name + `"\]\s*\{([^}]*)\}`).FindSubmatch(css)
			if block == nil {
				t.Fatal("missing scoped tokens")
			}
			tokens := make(map[string]string)
			for _, line := range strings.Split(string(block[1]), ";") {
				key, val, ok := strings.Cut(strings.TrimSpace(line), ":")
				if ok {
					tokens[strings.TrimSpace(key)] = strings.TrimSpace(val)
				}
			}
			for category, entries := range map[string]map[string]string{
				"radius":     a.Tokens.Radius,
				"border":     a.Tokens.Border,
				"shadow":     a.Tokens.Shadow,
				"typography": a.Tokens.Typography,
			} {
				for key, want := range entries {
					// Not every preset token is consumed by the default theme.
					cssKey := "--" + tokenToCSSName(key, map[string]string{"radius": "radius", "border": "border", "shadow": "shadow", "typography": ""}[category])
					if got, ok := tokens[cssKey]; ok && got != want {
						t.Errorf("%s: CSS %q, preset %q", cssKey, got, want)
					}
				}
			}
			got, err := strconv.ParseFloat(tokens["--spacing-scale"], 64)
			if err != nil || got != a.GetSpacingScale() {
				t.Errorf("spacing-scale: CSS %q, preset %s", tokens["--spacing-scale"], fmt.Sprint(a.GetSpacingScale()))
			}
		})
	}
}

func TestGenerateCSS_Brutal(t *testing.T) {
	// Create a brutalist aesthetic
	brutal := &Aesthetic{
		Name:        "Brutal",
		Description: "Brutalist design aesthetic with bold typography and sharp edges",
		Typography: map[string]string{
			"font-body":       "system-ui, sans-serif",
			"font-heading":    "Impact, Haettenschweiler, sans-serif",
			"font-mono":       "Consolas, monospace",
			"text-base":       "1rem",
			"text-lg":         "1.25rem",
			"text-xl":         "1.5rem",
			"text-2xl":        "2rem",
			"text-3xl":        "2.5rem",
			"leading-tight":   "1.25",
			"leading-normal":  "1.5",
			"leading-relaxed": "1.75",
		},
		Spacing: map[string]string{
			"space-1":  "0.25rem",
			"space-2":  "0.5rem",
			"space-3":  "0.75rem",
			"space-4":  "1rem",
			"space-6":  "1.5rem",
			"space-8":  "2rem",
			"space-12": "3rem",
			"space-16": "4rem",
		},
		Borders: map[string]string{
			"radius":       "0",
			"radius-lg":    "0",
			"border-width": "3px",
		},
		Shadows: map[string]string{
			"shadow-sm": "4px 4px 0 #000",
			"shadow-md": "6px 6px 0 #000",
			"shadow-lg": "8px 8px 0 #000",
		},
	}

	css := brutal.GenerateCSS()

	// Should contain :root selector
	if !strings.Contains(css, ":root") {
		t.Error("CSS should contain :root selector")
	}

	// Should contain header comment
	if !strings.Contains(css, "Brutal") {
		t.Error("CSS should contain palette name in header")
	}

	// Should contain typography variables
	if !strings.Contains(css, "--font-body:") {
		t.Error("CSS should contain --font-body")
	}
	if !strings.Contains(css, "--font-heading:") {
		t.Error("CSS should contain --font-heading")
	}
	if !strings.Contains(css, "Impact") {
		t.Error("CSS should contain Impact font for brutal heading")
	}

	// Should contain spacing variables
	if !strings.Contains(css, "--space-1:") {
		t.Error("CSS should contain --space-1")
	}
	if !strings.Contains(css, "--space-4:") {
		t.Error("CSS should contain --space-4")
	}

	// Should contain border variables - brutal has 0 radius
	if !strings.Contains(css, "--radius:") {
		t.Error("CSS should contain --radius")
	}
	// Check that radius is 0 for brutal
	if !strings.Contains(css, "--radius: 0") && !strings.Contains(css, "--radius:0") {
		t.Error("Brutal aesthetic should have --radius: 0")
	}

	// Should contain shadow variables - brutal has hard shadows
	if !strings.Contains(css, "--shadow-sm:") {
		t.Error("CSS should contain --shadow-sm")
	}
	if !strings.Contains(css, "4px 4px 0 #000") {
		t.Error("Brutal aesthetic should have hard shadow")
	}

	// Should not contain color variables (those come from palette)
	if strings.Contains(css, "--color-text:") {
		t.Error("Aesthetic CSS should not contain color variables")
	}
}

func TestGenerateCSS_AllAesthetics(t *testing.T) {
	// Get all built-in aesthetics
	names := BuiltinNames()
	if len(names) == 0 {
		t.Skip("No built-in aesthetics available")
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			a, err := LoadBuiltin(name)
			if err != nil {
				t.Fatalf("LoadBuiltin(%q) error = %v", name, err)
			}

			css := a.GenerateCSS()

			// Every aesthetic should produce valid CSS
			if !strings.Contains(css, ":root") {
				t.Error("CSS should contain :root selector")
			}

			// CSS should have opening and closing braces
			if !strings.Contains(css, "{") || !strings.Contains(css, "}") {
				t.Error("CSS should be valid with braces")
			}

			// Should contain some CSS variables
			if !strings.Contains(css, "--") {
				t.Error("CSS should contain CSS custom properties")
			}

			// Should not be empty
			if len(css) < 50 {
				t.Errorf("CSS seems too short (%d bytes), may be invalid", len(css))
			}
		})
	}
}

func TestGenerateCSS_Format(t *testing.T) {
	a := &Aesthetic{
		Name: "Test",
		Typography: map[string]string{
			"font-body": "system-ui",
		},
		Spacing: map[string]string{
			"space-1": "0.25rem",
		},
		Borders: map[string]string{
			"radius": "4px",
		},
	}

	t.Run("default format", func(t *testing.T) {
		css := a.GenerateCSS()

		// Should have nice formatting with newlines and indentation
		lines := strings.Split(css, "\n")
		if len(lines) < 3 {
			t.Error("Default CSS format should have multiple lines")
		}

		// Should have indented properties
		hasIndent := false
		for _, line := range lines {
			if strings.HasPrefix(line, "  ") {
				hasIndent = true
				break
			}
		}
		if !hasIndent {
			t.Error("Default CSS format should have indented properties")
		}
	})

	t.Run("minified format", func(t *testing.T) {
		css := a.GenerateCSSMinified()

		// Should be compact
		lines := strings.Split(strings.TrimSpace(css), "\n")
		if len(lines) > 2 {
			t.Errorf("Minified CSS should have few lines, got %d", len(lines))
		}
	})

	t.Run("with format options", func(t *testing.T) {
		format := CSSFormat{
			IncludeTypography: true,
			IncludeSpacing:    true,
			IncludeBorders:    false,
			IncludeShadows:    false,
		}

		css := a.GenerateCSSWithFormat(format)

		if !strings.Contains(css, "--font-body") {
			t.Error("CSS should include typography when enabled")
		}
		if !strings.Contains(css, "--space-1") {
			t.Error("CSS should include spacing when enabled")
		}
		if strings.Contains(css, "--radius") {
			t.Error("CSS should not include borders when disabled")
		}
	})
}

func TestGenerateCSS_CustomPrefix(t *testing.T) {
	a := &Aesthetic{
		Name: "Test",
		Typography: map[string]string{
			"font-body": "system-ui",
		},
	}

	format := CSSFormat{
		Prefix:            "theme",
		IncludeTypography: true,
	}

	css := a.GenerateCSSWithFormat(format)

	if !strings.Contains(css, "--theme-font-body") {
		t.Error("CSS should use custom prefix")
	}
}

func TestGenerateCSS_Sections(t *testing.T) {
	a := &Aesthetic{
		Name: "Complete",
		Typography: map[string]string{
			"font-body": "Arial",
		},
		Spacing: map[string]string{
			"space-1": "4px",
		},
		Borders: map[string]string{
			"radius": "8px",
		},
		Shadows: map[string]string{
			"shadow-sm": "0 1px 2px rgba(0,0,0,0.1)",
		},
	}

	css := a.GenerateCSS()

	// Should have section comments
	if !strings.Contains(css, "Typography") {
		t.Error("CSS should have Typography section comment")
	}
	if !strings.Contains(css, "Spacing") {
		t.Error("CSS should have Spacing section comment")
	}
	if !strings.Contains(css, "Border") {
		t.Error("CSS should have Borders section comment")
	}
	if !strings.Contains(css, "Shadow") {
		t.Error("CSS should have Shadows section comment")
	}
}

func TestGenerateCSS_EmptyAesthetic(t *testing.T) {
	a := &Aesthetic{
		Name: "Empty",
	}

	css := a.GenerateCSS()

	// Should still produce valid CSS structure
	if !strings.Contains(css, ":root") {
		t.Error("Empty aesthetic should still have :root")
	}
	if !strings.Contains(css, "{") || !strings.Contains(css, "}") {
		t.Error("Empty aesthetic should produce valid CSS structure")
	}
}

func TestCombineWithPalette(t *testing.T) {
	// Test that aesthetic CSS can be combined with palette CSS
	a := &Aesthetic{
		Name: "Test",
		Typography: map[string]string{
			"font-body": "system-ui",
		},
		Borders: map[string]string{
			"radius": "4px",
		},
	}

	aestheticCSS := a.GenerateCSS()

	// Simulate palette CSS
	paletteCSS := `:root {
  --palette-red: #ff0000;
  --color-primary: var(--palette-red);
}`

	// Combined CSS should have no conflicts
	combinedCSS := paletteCSS + "\n" + aestheticCSS

	// Both sections should be present
	if !strings.Contains(combinedCSS, "--palette-red") {
		t.Error("Combined CSS should contain palette variables")
	}
	if !strings.Contains(combinedCSS, "--font-body") {
		t.Error("Combined CSS should contain aesthetic typography")
	}
	if !strings.Contains(combinedCSS, "--radius") {
		t.Error("Combined CSS should contain aesthetic borders")
	}

	// Count :root occurrences (should be 2 - one from each)
	rootCount := strings.Count(combinedCSS, ":root")
	if rootCount != 2 {
		t.Errorf("Combined CSS has %d :root blocks, want 2", rootCount)
	}
}

func TestDefaultCSSFormat(t *testing.T) {
	format := DefaultCSSFormat()

	if !format.IncludeTypography {
		t.Error("Default format should include typography")
	}
	if !format.IncludeSpacing {
		t.Error("Default format should include spacing")
	}
	if !format.IncludeBorders {
		t.Error("Default format should include borders")
	}
	if !format.IncludeShadows {
		t.Error("Default format should include shadows")
	}
	if format.Minify {
		t.Error("Default format should not be minified")
	}
	if format.Prefix != "" {
		t.Errorf("Default format prefix = %q, want empty", format.Prefix)
	}
}
