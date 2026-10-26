package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestValidateConfig_Nil(t *testing.T) {
	errs := ValidateConfig(nil)
	if len(errs) == 0 {
		t.Error("ValidateConfig(nil) should return errors")
	}
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	config := &models.Config{
		URL:         "https://example.com",
		Concurrency: 4,
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
		Feeds: []models.FeedConfig{
			{
				Slug: "blog",
				Formats: models.FeedFormats{
					HTML: true,
				},
			},
		},
	}

	errs := ValidateConfig(config)
	if HasErrors(errs) {
		t.Errorf("ValidateConfig() unexpected errors: %v", errs)
	}
}

func TestValidateConfig_InvalidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"no scheme", "example.com"},
		{"invalid scheme", "ftp://example.com"},
		{"no host", "https://"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &models.Config{
				URL: tt.url,
				GlobConfig: models.GlobConfig{
					Patterns: []string{"**/*.md"},
				},
			}

			errs := ValidateConfig(config)
			if !HasErrors(errs) {
				t.Errorf("ValidateConfig() should return error for URL %q", tt.url)
			}
		})
	}
}

func TestValidateConfig_ValidURLs(t *testing.T) {
	tests := []string{
		"https://example.com",
		"http://localhost:8000",
		"https://sub.domain.example.com",
		"https://example.com/path/to/page",
	}

	for _, url := range tests {
		t.Run(url, func(t *testing.T) {
			config := &models.Config{
				URL: url,
				GlobConfig: models.GlobConfig{
					Patterns: []string{"**/*.md"},
				},
			}

			errs := ValidateConfig(config)
			for _, err := range errs {
				var ve ValidationError
				if errors.As(err, &ve) && ve.Field == "url" && !ve.IsWarn {
					t.Errorf("ValidateConfig() should not error for URL %q: %v", url, err)
				}
			}
		})
	}
}

func TestValidateConfig_NegativeConcurrency(t *testing.T) {
	config := &models.Config{
		Concurrency: -1,
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
	}

	errs := ValidateConfig(config)
	if !HasErrors(errs) {
		t.Error("ValidateConfig() should error for negative concurrency")
	}

	// Find the concurrency error
	found := false
	for _, err := range errs {
		var ve ValidationError
		if errors.As(err, &ve) && ve.Field == "concurrency" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected concurrency validation error")
	}
}

func TestValidateConfig_ZeroConcurrency(t *testing.T) {
	config := &models.Config{
		Concurrency: 0, // 0 means auto-detect, should be valid
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
	}

	errs := ValidateConfig(config)
	for _, err := range errs {
		var ve ValidationError
		if errors.As(err, &ve) && ve.Field == "concurrency" && !ve.IsWarn {
			t.Errorf("ValidateConfig() should not error for concurrency=0: %v", err)
		}
	}
}

func TestValidateConfig_EmptyGlobPatterns(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{},
		},
	}

	errs := ValidateConfig(config)
	if !HasWarnings(errs) {
		t.Error("ValidateConfig() should warn for empty glob patterns")
	}

	// Should be a warning, not an error
	actualErrors, warnings := SplitErrorsAndWarnings(errs)
	if len(actualErrors) > 0 {
		t.Errorf("Empty glob patterns should be warning, not error: %v", actualErrors)
	}
	if len(warnings) == 0 {
		t.Error("Expected warning for empty glob patterns")
	}
}

func TestValidateConfig_FeedEmptySlugAllowed(t *testing.T) {
	// Empty slug is now allowed - it represents the home page feed (index.html)
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
		Feeds: []models.FeedConfig{
			{
				Slug: "", // Empty slug is valid for home page
				Formats: models.FeedFormats{
					HTML: true,
				},
			},
		},
	}

	errs := ValidateConfig(config)
	actualErrors, _ := SplitErrorsAndWarnings(errs)
	if len(actualErrors) > 0 {
		t.Errorf("ValidateConfig() should allow empty slug for home page feed, got errors: %v", actualErrors)
	}
}

func TestValidateConfig_FeedNegativeItemsPerPage(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
		Feeds: []models.FeedConfig{
			{
				Slug:         "blog",
				ItemsPerPage: -5,
				Formats: models.FeedFormats{
					HTML: true,
				},
			},
		},
	}

	errs := ValidateConfig(config)
	if !HasErrors(errs) {
		t.Error("ValidateConfig() should error for negative items_per_page")
	}
}

func TestValidateConfig_FeedNegativeLimitOffset(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
		Feeds: []models.FeedConfig{
			{
				Slug:   "blog",
				Limit:  -1,
				Offset: -2,
				Formats: models.FeedFormats{
					HTML: true,
				},
			},
		},
	}

	errs := ValidateConfig(config)
	if !HasErrors(errs) {
		t.Error("ValidateConfig() should error for negative limit/offset")
	}
}

func TestValidateConfig_FeedNoFormats(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
		Feeds: []models.FeedConfig{
			{
				Slug:    "blog",
				Formats: models.FeedFormats{}, // All formats disabled
			},
		},
	}

	errs := ValidateConfig(config)
	if !HasWarnings(errs) {
		t.Error("ValidateConfig() should warn for feed with no formats")
	}
}

func TestValidateConfig_FeedDefaultsNegativeValues(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
		FeedDefaults: models.FeedDefaults{
			ItemsPerPage:    -1,
			OrphanThreshold: -1,
		},
	}

	errs := ValidateConfig(config)
	if !HasErrors(errs) {
		t.Error("ValidateConfig() should error for negative feed default values")
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Field:   "url",
		Message: "must include scheme",
		IsWarn:  false,
	}

	expected := "config error: url: must include scheme"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}

	warn := ValidationError{
		Field:   "glob.patterns",
		Message: "empty patterns",
		IsWarn:  true,
	}

	expectedWarn := "config warning: glob.patterns: empty patterns"
	if warn.Error() != expectedWarn {
		t.Errorf("Error() = %q, want %q", warn.Error(), expectedWarn)
	}
}

func TestHasErrors(t *testing.T) {
	tests := []struct {
		name string
		errs []error
		want bool
	}{
		{
			name: "no errors",
			errs: nil,
			want: false,
		},
		{
			name: "only warnings",
			errs: []error{
				ValidationError{Field: "f", Message: "m", IsWarn: true},
			},
			want: false,
		},
		{
			name: "has errors",
			errs: []error{
				ValidationError{Field: "f", Message: "m", IsWarn: false},
			},
			want: true,
		},
		{
			name: "mixed",
			errs: []error{
				ValidationError{Field: "f1", Message: "m1", IsWarn: true},
				ValidationError{Field: "f2", Message: "m2", IsWarn: false},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasErrors(tt.errs); got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasWarnings(t *testing.T) {
	tests := []struct {
		name string
		errs []error
		want bool
	}{
		{
			name: "no errors",
			errs: nil,
			want: false,
		},
		{
			name: "only errors",
			errs: []error{
				ValidationError{Field: "f", Message: "m", IsWarn: false},
			},
			want: false,
		},
		{
			name: "has warnings",
			errs: []error{
				ValidationError{Field: "f", Message: "m", IsWarn: true},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasWarnings(tt.errs); got != tt.want {
				t.Errorf("HasWarnings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitErrorsAndWarnings(t *testing.T) {
	errs := []error{
		ValidationError{Field: "f1", Message: "error", IsWarn: false},
		ValidationError{Field: "f2", Message: "warning", IsWarn: true},
		ValidationError{Field: "f3", Message: "error2", IsWarn: false},
		ValidationError{Field: "f4", Message: "warning2", IsWarn: true},
	}

	errsOut, warnings := SplitErrorsAndWarnings(errs)

	if len(errsOut) != 2 {
		t.Errorf("len(errsOut) = %d, want 2", len(errsOut))
	}
	if len(warnings) != 2 {
		t.Errorf("len(warnings) = %d, want 2", len(warnings))
	}
}

func TestSortErrors(t *testing.T) {
	errs := []error{
		ValidationError{Field: "f1", Message: "warning1", IsWarn: true},
		ValidationError{Field: "f2", Message: "error1", IsWarn: false},
		ValidationError{Field: "f3", Message: "warning2", IsWarn: true},
		ValidationError{Field: "f4", Message: "error2", IsWarn: false},
	}

	sortErrors(errs)

	// First two should be errors
	for i := 0; i < 2; i++ {
		if isWarning(errs[i]) {
			t.Errorf("errs[%d] should be an error, got warning", i)
		}
	}
	// Last two should be warnings
	for i := 2; i < 4; i++ {
		if !isWarning(errs[i]) {
			t.Errorf("errs[%d] should be a warning, got error", i)
		}
	}
}

func TestValidateAndWarn(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{}, // Will generate warning
		},
		Concurrency: -1, // Will generate error
	}

	var warnings []error
	warnFunc := func(err error) {
		warnings = append(warnings, err)
	}

	errs := ValidateAndWarn(config, warnFunc)

	// Should return errors
	if len(errs) == 0 {
		t.Error("ValidateAndWarn() should return errors")
	}

	// Warnings should be passed to warnFunc
	if len(warnings) == 0 {
		t.Error("ValidateAndWarn() should pass warnings to warnFunc")
	}

	// Errors should not include warnings
	for _, err := range errs {
		if isWarning(err) {
			t.Error("ValidateAndWarn() should not return warnings")
		}
	}
}

func TestValidateAndWarn_NilWarnFunc(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{}, // Will generate warning
		},
	}

	// Should not panic with nil warnFunc
	errs := ValidateAndWarn(config, nil)
	// Should return empty (warnings filtered out, no errors)
	if len(errs) != 0 {
		t.Errorf("ValidateAndWarn() returned errors: %v", errs)
	}
}

func TestHasAnyFormat(t *testing.T) {
	tests := []struct {
		name    string
		formats models.FeedFormats
		want    bool
	}{
		{
			name:    "all false",
			formats: models.FeedFormats{},
			want:    false,
		},
		{
			name:    "HTML true",
			formats: models.FeedFormats{HTML: true},
			want:    true,
		},
		{
			name:    "RSS true",
			formats: models.FeedFormats{RSS: true},
			want:    true,
		},
		{
			name:    "multiple true",
			formats: models.FeedFormats{HTML: true, RSS: true, Atom: true},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasAnyFormat(tt.formats); got != tt.want {
				t.Errorf("hasAnyFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"https://example.com", false},
		{"http://localhost:8000", false},
		{"example.com", true},       // no scheme
		{"ftp://example.com", true}, // invalid scheme
		{"https://", true},          // no host
		{"://example.com", true},    // missing scheme name
		{"https://example.com/path", false},
		{"https://sub.domain.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfig_MultipleFeedErrors(t *testing.T) {
	config := &models.Config{
		GlobConfig: models.GlobConfig{
			Patterns: []string{"**/*.md"},
		},
		Feeds: []models.FeedConfig{
			{Slug: "", Formats: models.FeedFormats{HTML: true}},                      // Empty slug is now valid
			{Slug: "valid", Formats: models.FeedFormats{}},                           // No formats (warning)
			{Slug: "bad", ItemsPerPage: -1, Formats: models.FeedFormats{HTML: true}}, // Negative items
		},
	}

	errs := ValidateConfig(config)
	actualErrors, warnings := SplitErrorsAndWarnings(errs)

	// Only 1 error now: negative items_per_page
	if len(actualErrors) < 1 {
		t.Errorf("Expected at least 1 error, got %d", len(actualErrors))
	}
	// 2 warnings: no formats on "valid" feed
	if len(warnings) < 1 {
		t.Errorf("Expected at least 1 warning, got %d", len(warnings))
	}
}

func TestValidateConfigWithPositions(t *testing.T) {
	tomlContent := `[markata-go]
title = "My Site"
url = "example.com"
concurrency = -1

[markata-go.glob]
patterns = []
`
	tracker := NewPositionTracker([]byte(tomlContent), "markata-go.toml")

	config := &models.Config{
		URL:         "example.com",
		Concurrency: -1,
		GlobConfig: models.GlobConfig{
			Patterns: []string{},
		},
	}

	configErrors := ValidateConfigWithPositions(config, tracker)

	if !configErrors.HasErrors() {
		t.Error("Expected errors for invalid config")
	}

	if !configErrors.HasWarnings() {
		t.Error("Expected warnings for empty patterns")
	}

	// Check that we have specific errors
	foundURL := false
	foundConcurrency := false
	foundPatterns := false

	for _, err := range configErrors.Errors {
		switch err.Field {
		case "url":
			foundURL = true
			if err.Line == 0 {
				t.Error("URL error should have line number")
			}
			if err.Fix == "" {
				t.Error("URL error should have fix suggestion")
			}
		case "concurrency":
			foundConcurrency = true
			if err.Line == 0 {
				t.Error("Concurrency error should have line number")
			}
		case "glob.patterns":
			foundPatterns = true
			if !err.IsWarn {
				t.Error("Empty patterns should be a warning")
			}
		}
	}

	if !foundURL {
		t.Error("Expected URL validation error")
	}
	if !foundConcurrency {
		t.Error("Expected concurrency validation error")
	}
	if !foundPatterns {
		t.Error("Expected patterns validation warning")
	}
}

func TestValidateConfigWithPositions_NilTracker(t *testing.T) {
	config := &models.Config{
		URL: "example.com",
	}

	configErrors := ValidateConfigWithPositions(config, nil)

	if !configErrors.HasErrors() {
		t.Error("Expected errors even without tracker")
	}

	// Errors should still be present, just without position info
	for _, err := range configErrors.Errors {
		if err.Field == "url" {
			if err.Line != 0 {
				t.Error("Without tracker, Line should be 0")
			}
			return
		}
	}
	t.Error("URL error not found")
}

func TestValidateConfigWithPositions_NilConfig(t *testing.T) {
	configErrors := ValidateConfigWithPositions(nil, nil)

	if len(configErrors.Errors) != 1 {
		t.Errorf("Expected 1 error for nil config, got %d", len(configErrors.Errors))
	}

	if configErrors.Errors[0].Field != "config" {
		t.Errorf("Expected 'config' field, got %q", configErrors.Errors[0].Field)
	}
}

func TestValidateConfig_CustomProjectPalette(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "palettes"), 0o755); err != nil {
		t.Fatal(err)
	}
	palette := "[palette]\nname = \"Custom\"\nvariant = \"dark\"\n\n[palette.colors]\nbg = \"#101112\"\ntext = \"#e0e1e2\"\naccent = \"#778899\"\n\n[palette.semantic]\nbg-primary = \"bg\"\ntext-primary = \"text\"\naccent = \"accent\"\nlink = \"accent\"\n"
	if err := os.WriteFile(filepath.Join(dir, "palettes", "custom-dark.toml"), []byte(palette), 0o600); err != nil {
		t.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	config := &models.Config{URL: "https://example.com", Concurrency: 1}
	config.Theme.Palette = "custom-dark"
	if errs := ValidateConfig(config); HasErrors(errs) {
		t.Errorf("custom project palette rejected: %v", errs)
	}

	config.Theme.Palette = "definitely-missing"
	if errs := ValidateConfig(config); !HasErrors(errs) {
		t.Error("expected unknown palette error")
	}
}

func TestValidateConfig_ThemeFontpackAcceptsCatalogPacks(t *testing.T) {
	tests := []struct {
		name     string
		fontpack string
		wantErr  bool
	}{
		{"contract fontpack", "brush", false},
		{"catalog-only fontpack", "typewriter", false},
		{"catalog alias", "reader", false},
		{"unknown fontpack", "definitely-missing", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &models.Config{URL: "https://example.com", Concurrency: 1}
			config.Theme.Fontpack = tt.fontpack
			errs := ValidateConfig(config)
			if HasErrors(errs) != tt.wantErr {
				t.Fatalf("HasErrors = %v, want %v: %v", HasErrors(errs), tt.wantErr, errs)
			}
		})
	}
}
