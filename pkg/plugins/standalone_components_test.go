package plugins

import (
	"strings"
	"testing"
)

func TestStandaloneIcons_ProseAndCode(t *testing.T) {
	r, err := RenderStandalone("## :lucide-globe: Hello\n\n:lucide/heart: :lucide-no-such-icon:\n\n`:lucide-globe:`\n\n```md\n:lucide-globe:\n```", StandaloneOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(r.HTML, `data-icon="lucide/globe"`) != 1 || !strings.Contains(r.HTML, `data-icon="lucide/heart"`) {
		t.Fatal("icons not expanded in prose")
	}
	if !strings.Contains(r.HTML, "<code>:lucide-globe:</code>") || !strings.Contains(r.HTML, ":lucide-no-such-icon:") {
		t.Fatal("code or unknown icon changed")
	}
	if strings.Contains(r.HTML, "<!-- raw HTML omitted -->") {
		t.Fatal("generated SVG must work with authored HTML disabled")
	}
}

func TestStandaloneAppearance_CatalogAndAssets(t *testing.T) {
	r, err := BuildStandaloneAppearance()
	if err != nil {
		t.Fatal(err)
	}
	if r.SchemaVersion != 1 || !strings.Contains(r.PickerHTML, `data-picker-tab="font"`) {
		t.Fatal("missing controls")
	}
	for _, name := range append(r.Styles, r.Scripts...) {
		if len(r.Assets[name]) == 0 {
			t.Errorf("missing %s", name)
		}
	}
	for _, item := range []struct{ name, token string }{{"css/palette.css", "--palette-manifest"}, {"css/aesthetic.css", "--aesthetic-manifest"}, {"css/fonts.css", "--fontpack-manifest"}, {"js/standalone-appearance.js", "window.__markataSeasonal ="}} {
		if !strings.Contains(string(r.Assets[item.name]), item.token) {
			t.Errorf("missing %s", item.token)
		}
	}
	found := false
	for name := range r.Assets {
		if strings.HasPrefix(name, "assets/fonts/") {
			found = true
		}
	}
	if !found {
		t.Fatal("font binaries missing")
	}
}
