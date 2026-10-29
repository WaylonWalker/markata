package themes

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

// Both the legacy and embedded templates must be safe before CSS and JS load.
// The dialog is present on every generated page but should not expose focusable
// controls until the shortcuts handler explicitly opens it.
func TestShortcutsDialogAccessibilityInBothTemplates(t *testing.T) {
	for _, path := range []string{
		"../../templates/partials/shortcuts-modal.html",
		"default/templates/partials/shortcuts-modal.html",
	} {
		t.Run(path, func(t *testing.T) {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			markup := string(data)
			for _, attribute := range []string{"hidden", "inert"} {
				pattern := regexp.MustCompile(`<div id="shortcuts-modal"[^>]*\s+` + attribute + `(?:\s|>)`)
				if !pattern.MatchString(markup) {
					t.Errorf("closed shortcuts dialog must have native %s attribute", attribute)
				}
			}
			if !strings.Contains(markup, `aria-labelledby="shortcuts-modal-title"`) {
				t.Error("shortcuts dialog must have an accessible title")
			}

			tables := strings.Count(markup, `<table class="shortcuts-table">`)
			if tables < 4 {
				t.Fatalf("expected at least four shortcut tables, got %d", tables)
			}
			for _, element := range []string{
				`<caption class="visually-hidden">`,
				`<th scope="col">Keys</th>`,
				`<th scope="col">Action</th>`,
			} {
				if got := strings.Count(markup, element); got != tables {
					t.Errorf("expected %d occurrences of %q, got %d", tables, element, got)
				}
			}
		})
	}
}
