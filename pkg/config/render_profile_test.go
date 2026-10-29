package config

import (
	"strings"
	"testing"
)

func TestRenderProfile_ConfigFormatsAndSidebar(t *testing.T) {
	for _, tc := range []struct{
		name, data string
		format Format
	}{
		{"toml", "[markata-go.markdown]\nrender_profile = \"offline-interactive\"\n", FormatTOML},
		{"yaml", "markata-go:\n  markdown:\n    render_profile: offline-interactive\n", FormatYAML},
		{"json", `{"markata-go":{"markdown":{"render_profile":"offline-interactive"}}}`, FormatJSON},
	} {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := LoadFromString(tc.data, tc.format)
			if err != nil {
				t.Fatal(err)
			}
			if cfg.MarkdownConfig.RenderProfile != "offline-interactive" {
				t.Fatalf("config lost render profile: %q", cfg.MarkdownConfig.RenderProfile)
			}
			field, ok := LookupSettingIn(cfg, "markdown.render_profile")
			if !ok || !field.Editable || !field.Closed {
				t.Fatalf("serve sidebar cannot edit profile: %+v", field)
			}
			if !strings.Contains(strings.Join(field.Options, ","), "offline-static") {
				t.Fatalf("incomplete profile options: %v", field.Options)
			}
		})
	}
}

func TestRenderProfile_ConfigMergeAndValidation(t *testing.T) {
	base, err := LoadFromString("[markata-go.markdown]\nrender_profile = \"offline-interactive\"\n", FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	other, err := LoadFromString("[markata-go]\ntitle = \"Untouched\"\n", FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	merged := mergeMarkdownConfig(base.MarkdownConfig, other.MarkdownConfig)
	if merged.RenderProfile != "offline-interactive" {
		t.Fatalf("unrelated override clobbered profile: %q", merged.RenderProfile)
	}
	if f, ok := LookupSettingIn(base, "markdown.render_profile"); !ok {
		t.Fatal("missing config setting")
	} else if _, err := CoerceSettingValue(f, "not-a-profile"); err == nil {
		t.Fatal("accepted unrecognized profile in settings")
	}
}
