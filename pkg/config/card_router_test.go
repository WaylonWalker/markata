package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

// Regression for #1549: all format-specific converters and the resolved
// config loader must preserve custom card router mappings.
func TestCardRouterMappingsAcrossConfigFormats(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		ext    string
		data   string
		parse  func([]byte) (*models.Config, error)
	}{
		{
			name:   "toml",
			format: FormatTOML,
			ext:    ".toml",
			parse:  ParseTOML,
			data: `[markata-go]
title = "Test"
url = "https://example.com"

[markata-go.components.card_router.mappings]
note = "article"
custom = "custom-card"
`,
		},
		{
			name:   "yaml",
			format: FormatYAML,
			ext:    ".yaml",
			parse:  ParseYAML,
			data: `markata-go:
  title: Test
  url: https://example.com
  components:
    card_router:
      mappings:
        note: article
        custom: custom-card
`,
		},
		{
			name:   "json",
			format: FormatJSON,
			ext:    ".json",
			parse:  ParseJSON,
			data:   `{"markata-go":{"title":"Test","url":"https://example.com","components":{"card_router":{"mappings":{"note":"article","custom":"custom-card"}}}}}`,
		},
	}

	assertMappings := func(t *testing.T, config *models.Config) {
		t.Helper()
		for key, want := range map[string]string{"note": "article", "custom": "custom-card"} {
			if got := config.Components.CardRouter.Mappings[key]; got != want {
				t.Errorf("CardRouter.Mappings[%q] = %q, want %q", key, got, want)
			}
			if got := config.Components.CardRouter.MergedMappings()[key]; got != want {
				t.Errorf("CardRouter.MergedMappings()[%q] = %q, want %q", key, got, want)
			}
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := tt.parse([]byte(tt.data))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			assertMappings(t, parsed)

			fromString, err := LoadFromString(tt.data, tt.format)
			if err != nil {
				t.Fatalf("LoadFromString: %v", err)
			}
			assertMappings(t, fromString)

			path := filepath.Join(t.TempDir(), "markata-go"+tt.ext)
			//nolint:gosec // Temporary test fixture.
			if err := os.WriteFile(path, []byte(tt.data), 0o644); err != nil {
				t.Fatal(err)
			}
			loaded, err := Load(path)
			if err != nil {
				t.Fatalf("Load: %v", err)
			}
			assertMappings(t, loaded)
		})
	}
}
