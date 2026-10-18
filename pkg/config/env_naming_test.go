package config

import (
	"strings"
	"testing"
)

func TestMarkataEnvCanonicalOverride(t *testing.T) {
	t.Setenv("MARKATA_TITLE", "Canonical title")
	config := DefaultConfig()
	if err := ApplyEnvOverrides(config); err != nil {
		t.Fatal(err)
	}
	if config.Title != "Canonical title" {
		t.Fatalf("title = %q, want canonical title", config.Title)
	}
}

func TestMarkataEnvLegacyStillWorks(t *testing.T) {
	t.Setenv("MARKATA_GO_TITLE", "Legacy title")
	config := DefaultConfig()
	if err := ApplyEnvOverrides(config); err != nil {
		t.Fatal(err)
	}
	if config.Title != "Legacy title" {
		t.Fatalf("title = %q, want legacy title", config.Title)
	}
}

func TestMarkataEnvCanonicalWinsAndWarnsWithoutValues(t *testing.T) {
	t.Setenv("MARKATA_GO_TITLE", "legacy-secret")
	t.Setenv("MARKATA_TITLE", "canonical-secret")
	t.Setenv("MARKATA_GO_OUTPUT_DIR", "legacy-output")
	t.Setenv("MARKATA_OUTPUT_DIR", "canonical-output")

	config := DefaultConfig()
	if err := ApplyEnvOverrides(config); err != nil {
		t.Fatal(err)
	}
	if config.Title != "canonical-secret" || config.OutputDir != "canonical-output" {
		t.Fatalf("canonical values not selected: title=%q output=%q", config.Title, config.OutputDir)
	}
	warnings, ok := config.Extra["name_migration_warnings"].([]string)
	if !ok || len(warnings) != 2 {
		t.Fatalf("name_migration_warnings = %#v, want two", config.Extra["name_migration_warnings"])
	}
	if warnings[0] != "MARKATA_OUTPUT_DIR overrides MARKATA_GO_OUTPUT_DIR" ||
		warnings[1] != "MARKATA_TITLE overrides MARKATA_GO_TITLE" {
		t.Fatalf("warnings = %#v, want sorted key names", warnings)
	}
	for _, warning := range warnings {
		if strings.Contains(warning, "secret") {
			t.Fatalf("migration warning exposed a config value: %q", warning)
		}
	}
}

func TestMarkataEnvMatchingAliasesNoWarning(t *testing.T) {
	t.Setenv("MARKATA_GO_TITLE", "Same title")
	t.Setenv("MARKATA_TITLE", "Same title")
	config := DefaultConfig()
	if err := ApplyEnvOverrides(config); err != nil {
		t.Fatal(err)
	}
	if config.Title != "Same title" {
		t.Fatalf("title = %q, want Same title", config.Title)
	}
	if _, ok := config.Extra["name_migration_warnings"]; ok {
		t.Fatalf("unexpected warning: %#v", config.Extra["name_migration_warnings"])
	}
}

func TestGetEnvValuePrefersMarkata(t *testing.T) {
	t.Setenv("MARKATA_GO_MIGRATION_TEST", "legacy")
	t.Setenv("MARKATA_MIGRATION_TEST", "canonical")
	if got, ok := GetEnvValue("migration_test"); !ok || got != "canonical" {
		t.Fatalf("GetEnvValue(migration_test) = %q, %v", got, ok)
	}
	t.Setenv("MARKATA_MIGRATION_TEST", "")
	if got, ok := GetEnvValue("migration_test"); !ok || got != "" {
		t.Fatalf("explicit empty canonical value lost: %q, %v", got, ok)
	}
}

func TestGetEnvValueLegacyFallback(t *testing.T) {
	t.Setenv("MARKATA_GO_NAMING_LEGACY_TEST", "still supported")
	if got, ok := GetEnvValue("naming_legacy_test"); !ok || got != "still supported" {
		t.Fatalf("legacy fallback = %q, %v", got, ok)
	}
}
