package themes

import (
	"regexp"
	"strings"
	"testing"
)

// A stylesheet that is late, blocked, or purged must not let the keyboard
// shortcuts dialog create a large scrollable region after the site footer.
func TestShortcutsModalHiddenWithoutCSS(t *testing.T) {
	tmpl, err := ReadTemplate("partials/shortcuts-modal.html")
	if err != nil {
		t.Fatal(err)
	}

	// The browser's native hidden attribute works without JavaScript or CSS.
	if !regexp.MustCompile(`<div id="shortcuts-modal"[^>]*\bhidden\b`).Match(tmpl) {
		t.Fatal("closed shortcuts dialog must have native hidden attribute")
	}

	css, err := ReadStatic("css/components.css")
	if err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`(?s)\.shortcuts-modal\[hidden\]\s*\{\s*display:\s*none;`).Match(css) {
		t.Fatal("author CSS must respect hidden even though the dialog uses display:flex")
	}

	js, err := ReadStatic("js/search-shortcuts.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"modal.removeAttribute('hidden');",
		"modal.setAttribute('hidden', '');",
		"modal.removeAttribute('inert');",
		"modal.setAttribute('inert', '');",
	} {
		if !strings.Contains(string(js), expected) {
			t.Errorf("shortcuts dialog must manage both hidden and inert: missing %q", expected)
		}
	}
}
