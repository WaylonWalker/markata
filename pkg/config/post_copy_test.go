package config

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestPostCopyConfig_ParsingAndMerge(t *testing.T) {
	parsers := []struct {
		name  string
		parse func([]byte) (*models.Config, error)
		data  string
	}{
		{"TOML", ParseTOML, "[markata-go.components.post_copy]\nenabled = false\n"},
		{"YAML", ParseYAML, "markata-go:\n  components:\n    post_copy:\n      enabled: false\n"},
		{"JSON", ParseJSON, `{"markata-go":{"components":{"post_copy":{"enabled":false}}}}`},
	}
	for _, tc := range parsers {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := tc.parse([]byte(tc.data))
			if err != nil {
				t.Fatal(err)
			}
			if parsed.Components.PostCopy.IsEnabled() {
				t.Fatal("explicit false must disable the copy menu")
			}
			merged := mergeComponentsConfig(models.NewComponentsConfig(), parsed.Components)
			if merged.PostCopy.IsEnabled() {
				t.Fatal("merged config lost explicit false")
			}
		})
	}
	if !models.NewComponentsConfig().PostCopy.IsEnabled() {
		t.Fatal("post copy menu must remain on by default")
	}
	field, ok := LookupSetting("components.post_copy.enabled")
	if !ok || field.Kind != SettingBool || !field.Editable {
		t.Fatalf("post copy setting is not editable: %+v, found %v", field, ok)
	}
	for _, setting := range Settings(models.NewConfig()) {
		if setting.Key == field.Key {
			if setting.Value != true {
				t.Fatalf("post copy default = %v, want true", setting.Value)
			}
			return
		}
	}
	t.Fatal("post copy setting missing from settings list")
}
