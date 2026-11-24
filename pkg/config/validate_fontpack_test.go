package config

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestValidateRenderingTheme_CatalogFontpacks(t *testing.T) {
	tests := []struct {
		name     string
		fontpack string
		file     string
		wantErr  bool
	}{
		{"catalog pack", "typewriter", "", false},
		{"catalog alias", "notebook", "", false},
		{"contract pack", "brush", "", false},
		{"unknown pack", "not-a-pack", "", true},
		{"custom catalog", "my-pack", "fonts/catalog.yaml", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &models.Config{}
			cfg.Theme.Fontpack = tt.fontpack
			cfg.FontpacksFile = tt.file
			gotErr := false
			for _, err := range validateRenderingTheme(cfg) {
				if strings.Contains(err.Error(), "theme.fontpack") {
					gotErr = true
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("fontpack %q: gotErr=%v, want %v", tt.fontpack, gotErr, tt.wantErr)
			}
		})
	}
}

func TestValidateRenderingTheme_LegacyTopLevelFontpack(t *testing.T) {
	tests := []struct {
		name     string
		fontpack string
		wantErr  bool
	}{
		{"known legacy pack", "typewriter", false},
		{"unknown legacy pack", "not-a-pack", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &models.Config{Fontpack: tt.fontpack}
			cfg.Theme.Fontpack = tt.fontpack
			gotErr := false
			for _, err := range validateRenderingTheme(cfg) {
				if strings.Contains(err.Error(), "fontpack") {
					gotErr = true
				}
			}
			if gotErr != tt.wantErr {
				t.Errorf("legacy fontpack %q: gotErr=%v, want %v", tt.fontpack, gotErr, tt.wantErr)
			}
		})
	}
}
