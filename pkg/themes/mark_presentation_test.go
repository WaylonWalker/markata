package themes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The default theme is embedded for released binaries, while the repository's
// rendering-contract fixture uses its own templates/base.html. Both must avoid
// the old forced-reflow script, and both must keep CSS-driven mark colors.
func TestMarkPresentationDoesNotRequireDOMStyleUpdater(t *testing.T) {
	builtin, err := ReadTemplate("base.html")
	if err != nil {
		t.Fatalf("read embedded base template: %v", err)
	}
	fixture, err := os.ReadFile(filepath.Join("..", "..", "templates", "base.html"))
	if err != nil {
		t.Fatalf("read fixture base template: %v", err)
	}

	for name, template := range map[string][]byte{
		"embedded": builtin,
		"fixture":  fixture,
	} {
		t.Run(name, func(t *testing.T) {
			source := string(template)
			for _, legacy := range []string{
				"applyMarkPresentation",
				"scheduleMarkPresentation",
				"querySelectorAll('mark').forEach",
			} {
				if strings.Contains(source, legacy) {
					t.Errorf("legacy mark updater %q remains in template", legacy)
				}
			}
			if !strings.Contains(source, "id=\"markata-mark-presentation\"") ||
				!strings.Contains(source, "html mark { background-color: var(--color-highlight)") {
				t.Error("template must retain CSS-driven mark presentation")
			}
		})
	}

	// CSS custom properties update automatically on theme changes and DOM
	// insertion, without the old mutation observers or computed-style reads.
	for _, tc := range []struct {
		path     string
		required []string
	}{
		{"css/main.css", []string{"html[data-fontpack] body mark {", "background-color: var(--color-highlight) !important;", "color: var(--color-highlight-text) !important;"}},
		{"css/components.css", []string{"html mark,", "html[data-palette] mark", "background-color: var(--color-highlight) !important;"}},
		{"css/variables.css", []string{"--heading-mark-bg: var(--color-highlight, var(--color-primary));", "--heading-mark-fg: var(--color-highlight-text, var(--color-background));"}},
	} {
		t.Run(tc.path, func(t *testing.T) {
			content, err := ReadStatic(tc.path)
			if err != nil {
				t.Fatalf("read %s: %v", tc.path, err)
			}
			for _, rule := range tc.required {
				if !strings.Contains(string(content), rule) {
					t.Errorf("%s missing presentation contract %q", tc.path, rule)
				}
			}
		})
	}
}
