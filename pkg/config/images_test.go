package config

import (
	"errors"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLoadFromString_ImagesConfigAcrossFormats(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		data   string
	}{
		{"toml", FormatTOML, "[markata-go.images]\nenabled = false\npath = \"media\"\ntemplate = \"media.html\"\nexport_json = false\ninclude_unreferenced = false\n"},
		{"yaml", FormatYAML, "markata-go:\n  images:\n    enabled: false\n    path: media\n    template: media.html\n    export_json: false\n    include_unreferenced: false\n"},
		{"json", FormatJSON, `{"markata-go":{"images":{"enabled":false,"path":"media","template":"media.html","export_json":false,"include_unreferenced":false}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadFromString(tt.data, tt.format)
			if err != nil {
				t.Fatalf("LoadFromString() error = %v", err)
			}
			if config.Images.IsEnabled() || config.Images.ShouldExportJSON() || config.Images.ShouldIncludeUnreferenced() || config.Images.Path != "media" || config.Images.Template != "media.html" {
				t.Fatalf("Images = %#v", config.Images)
			}
		})
	}
}

func TestImagesConfig_IncludeUnreferencedDefaultsFalse(t *testing.T) {
	var zero models.ImagesConfig
	if zero.ShouldIncludeUnreferenced() {
		t.Fatal("zero-valued ImagesConfig includes unreferenced media")
	}

	var nilConfig *models.ImagesConfig
	if nilConfig.ShouldIncludeUnreferenced() {
		t.Fatal("nil ImagesConfig includes unreferenced media")
	}

	defaults := DefaultConfig()
	if defaults.Images.ShouldIncludeUnreferenced() {
		t.Fatal("DefaultConfig includes unreferenced media")
	}
}

func TestMergeConfigs_ImagesPreservesExplicitFalse(t *testing.T) {
	base := models.NewConfig()
	base.Images.Path = "images"
	falseValue := false
	override := &models.Config{Images: models.ImagesConfig{
		Enabled:             &falseValue,
		ExportJSON:          &falseValue,
		IncludeUnreferenced: &falseValue,
	}}

	merged := MergeConfigs(base, override)
	if merged.Images.IsEnabled() || merged.Images.ShouldExportJSON() || merged.Images.ShouldIncludeUnreferenced() {
		t.Fatalf("merged image config = %#v", merged.Images)
	}
	if merged.Images.Path != "images" || merged.Images.Template != "images.html" {
		t.Fatalf("merged defaults were lost = %#v", merged.Images)
	}
}

func TestApplyEnvOverrides_Images(t *testing.T) {
	t.Setenv("MARKATA_GO_IMAGES_ENABLED", "false")
	t.Setenv("MARKATA_GO_IMAGES_PATH", "media")
	t.Setenv("MARKATA_GO_IMAGES_TEMPLATE", "media.html")
	t.Setenv("MARKATA_GO_IMAGES_EXPORT_JSON", "false")
	t.Setenv("MARKATA_GO_IMAGES_INCLUDE_UNREFERENCED", "false")

	config := DefaultConfig()
	if err := ApplyEnvOverrides(config); err != nil {
		t.Fatalf("ApplyEnvOverrides() error = %v", err)
	}
	if config.Images.IsEnabled() || config.Images.ShouldExportJSON() || config.Images.ShouldIncludeUnreferenced() || config.Images.Path != "media" || config.Images.Template != "media.html" {
		t.Fatalf("Images = %#v", config.Images)
	}
}

func TestValidateConfig_RejectsImagePathTraversal(t *testing.T) {
	config := models.NewConfig()
	config.Images.Path = "../outside"
	errs := ValidateConfig(config)
	for _, err := range errs {
		var validation ValidationError
		if errors.As(err, &validation) && validation.Field == "images.path" {
			return
		}
	}
	t.Fatalf("ValidateConfig() errors = %v, want images.path error", errs)
}
