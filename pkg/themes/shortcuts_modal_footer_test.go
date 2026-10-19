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


func TestFooterShortcutsTriggerIsCompactAndClickable(t *testing.T) {
	footer, err := ReadTemplate("components/footer.html")
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{
		`type="button"`,
		`class="footer-shortcuts-trigger"`,
		`data-shortcuts-trigger`,
		`aria-label="Show keyboard shortcuts"`,
		`aria-controls="shortcuts-modal"`,
		`<kbd aria-hidden="true">?</kbd>`,
	} {
		if !strings.Contains(string(footer), needle) {
			t.Errorf("footer missing compact keyboard shortcut trigger: %s", needle)
		}
	}
	if strings.Contains(string(footer), `shortcuts-modal-body`) {
		t.Error("expanded shortcuts list must not live in footer")
	}

	js, err := ReadStatic("js/search-shortcuts.js")
	if err != nil {
		t.Fatal(err)
	}
	for _, needle := range []string{
		`document.addEventListener('click', function(e)`,
		`e.target.closest('[data-shortcuts-trigger]')`,
		`showShortcutsModal();`,
	} {
		if !strings.Contains(string(js), needle) {
			t.Errorf("shortcut button must survive SPA footer replacement: missing %s", needle)
		}
	}

	css, err := ReadStatic("css/components.css")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(css), ".footer-shortcuts-trigger {") {
		t.Fatal("footer shortcut button missing minimal styling")
	}
}
