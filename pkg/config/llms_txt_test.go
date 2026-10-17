package config

import (
	"testing"

	"github.com/WaylonWalker/markata-go/pkg/models"
)

func TestParseLLMsTXTPluginOptionsAcrossFormats(t *testing.T) {
	cases := []struct {
		name  string
		parse func([]byte) (*models.Config, error)
		input string
	}{
		{
			name:  "TOML",
			parse: ParseTOML,
			input: "[markata-go]\ntitle = \"Test\"\n[markata-go.llms_txt]\nenabled = false\nmax_pages = 7\n",
		},
		{
			name:  "YAML",
			parse: ParseYAML,
			input: "markata-go:\n  title: Test\n  llms_txt:\n    enabled: false\n    max_pages: 7\n",
		},
		{
			name:  "JSON",
			parse: ParseJSON,
			input: "{\"markata-go\":{\"title\":\"Test\",\"llms_txt\":{\"enabled\":false,\"max_pages\":7}}}",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := tc.parse([]byte(tc.input))
			if err != nil {
				t.Fatal(err)
			}
			raw, ok := cfg.Extra["llms_txt"].(map[string]any)
			if !ok {
				t.Fatalf("missing llms_txt options: %#v", cfg.Extra)
			}
			if enabled, ok := raw["enabled"].(bool); !ok || enabled {
				t.Fatalf("expected enabled=false, got %#v", raw["enabled"])
			}
			switch n := raw["max_pages"].(type) {
			case int:
				if n != 7 {
					t.Fatalf("max_pages = %d", n)
				}
			case int64:
				if n != 7 {
					t.Fatalf("max_pages = %d", n)
				}
			case float64:
				if n != 7 {
					t.Fatalf("max_pages = %v", n)
				}
			default:
				t.Fatalf("unexpected max_pages type %T: %v", n, n)
			}
		})
	}
}
