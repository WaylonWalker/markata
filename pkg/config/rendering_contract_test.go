package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLoadRenderingContract_FontpackWarnings(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		fontpack string
		warnings []string
	}{
		{name: "defaults", data: "[markata-go]\n", fontpack: "system"},
		{name: "canonical", data: "[markata-go.theme]\nfontpack = \"typewriter\"\n", fontpack: "typewriter"},
		{name: "legacy", data: "[markata-go]\nfontpack = \"typewriter\"\n", fontpack: "typewriter"},
		{
			name: "matching", fontpack: "typewriter",
			data: "[markata-go]\nfontpack = \"typewriter\"\n[markata-go.theme]\nfontpack = \"typewriter\"\n",
		},
		{
			name: "conflicting", fontpack: "typewriter",
			data:     "[markata-go]\nfontpack = \"system\"\n[markata-go.theme]\nfontpack = \"typewriter\"\n",
			warnings: []string{"fontpack conflicts with theme.fontpack; canonical value wins"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "markata-go.toml")
			if err := os.WriteFile(path, []byte(tt.data), 0o600); err != nil {
				t.Fatal(err)
			}
			loaders := map[string]func() (*models.Config, error){
				"file": func() (*models.Config, error) { return loadResolvedConfig(path) },
				"merge": func() (*models.Config, error) {
					return LoadWithMergeOptions(LoadOptions{DisableDotEnv: true, DisableEnvOverrides: true}, path)
				},
				"string": func() (*models.Config, error) { return LoadFromString(tt.data, FormatTOML) },
				"single": func() (*models.Config, error) { return LoadSingleConfig(path) },
			}
			for name, load := range loaders {
				t.Run(name, func(t *testing.T) {
					cfg, err := load()
					if err != nil {
						t.Fatal(err)
					}
					if cfg.Theme.Fontpack != tt.fontpack {
						t.Errorf("fontpack = %q, want %q", cfg.Theme.Fontpack, tt.fontpack)
					}
					var warnings []string
					if warningValue, exists := cfg.Extra["theme_migration_warnings"]; exists {
						var ok bool
						warnings, ok = warningValue.([]string)
						if !ok {
							t.Fatalf("warnings type = %T, want []string", warningValue)
						}
					}
					if !reflect.DeepEqual(warnings, tt.warnings) {
						t.Errorf("warnings = %v, want %v", warnings, tt.warnings)
					}
				})
			}
		})
	}
}

func TestLoadRenderingContract_DefaultsWithoutWarnings(t *testing.T) {
	cfg, err := LoadWithDefaults()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Theme.Fontpack != "system" {
		t.Errorf("fontpack = %q, want system", cfg.Theme.Fontpack)
	}
	if warnings, ok := cfg.Extra["theme_migration_warnings"]; ok {
		t.Errorf("unexpected default migration warnings: %v", warnings)
	}
}

func TestLoadRenderingContract_LegacyMigratesWithCanonicalPrecedence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "markata-go.toml")
	contents := `[markata-go]
texture_strength = "10%"
heading_texture_strength = "20%"
motif_color_distance = "1%"

[markata-go.theme]
contract_version = 1
palette = "ayu-dark"
fontpack = "brush-poster"

[markata-go.theme.texture]
kind = "screenprint"
color_mix = 0.9
scale = 1.0
scope = "all"
`
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := loadResolvedConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.Theme.Texture.ColorMix != .9 {
		t.Fatalf("color_mix = %v", config.Theme.Texture.ColorMix)
	}
	if config.Theme.Fontpack != "brush-poster" {
		t.Fatalf("fontpack = %q", config.Theme.Fontpack)
	}
	warnings, ok := config.Extra["theme_migration_warnings"].([]string)
	if !ok || len(warnings) == 0 {
		t.Fatal("expected migration warning")
	}
}

func TestLoadRenderingContract_LegacyOnlyValuesSurviveDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "markata-go.toml")
	if err := os.WriteFile(path, []byte("[markata-go]\ntexture_strength = \"10%\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	config, err := loadResolvedConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.Theme.Texture.ColorMix != .1 {
		t.Fatalf("legacy color_mix = %v, want .1", config.Theme.Texture.ColorMix)
	}
}
