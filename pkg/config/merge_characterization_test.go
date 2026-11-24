package config

import (
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

// TestMergeConfigs_CharacterizesTopLevelFields records the complete result of
// merging every top-level field handled directly by MergeConfigs. The fields
// that MergeConfigs intentionally leaves alone are populated on base as well,
// so a refactor cannot accidentally broaden its contract.
func TestMergeConfigs_CharacterizesTopLevelFields(t *testing.T) {
	base := &models.Config{
		Fontpack:         "base-fontpack",
		FontpacksFile:    "base-fontpacks.toml",
		OutputDir:        "base-output",
		URL:              "https://base.example",
		Title:            "Base title",
		Description:      "Base description",
		Author:           "Base author",
		Language:         "en",
		AuthorURL:        "https://base.example/about",
		ManagingEditor:   "base-editor@example",
		WebMaster:        "base-webmaster@example",
		Copyright:        "Base copyright",
		License:          models.LicenseValue{Raw: "cc-by-4.0"},
		AssetsDir:        "base-static",
		TemplatesDir:     "base-templates",
		Hooks:            []string{"base-hook"},
		DisabledHooks:    []string{"base-disabled"},
		Nav:              []models.NavItem{{Label: "Base", URL: "/base"}},
		Concurrency:      2,
		Feeds:            []models.FeedConfig{{Slug: "base-feed"}},
		Extra:            map[string]any{"base": "value"},
		TemplatePresets:  map[string]models.TemplatePreset{"base": {}},
		DefaultTemplates: map[string]string{"html": "base.html"},
	}
	override := &models.Config{
		Fontpack:       "override-fontpack",
		FontpacksFile:  "override-fontpacks.toml",
		OutputDir:      "override-output",
		URL:            "https://override.example",
		Title:          "Override title",
		Description:    "Override description",
		Author:         "Override author",
		Language:       "fr",
		AuthorURL:      "https://override.example/about",
		ManagingEditor: "override-editor@example",
		WebMaster:      "override-webmaster@example",
		Copyright:      "Override copyright",
		License:        models.LicenseValue{Raw: false},
		AssetsDir:      "override-static",
		TemplatesDir:   "override-templates",
		Hooks:          []string{"override-hook"},
		DisabledHooks:  []string{"override-disabled"},
		Nav:            []models.NavItem{{Label: "Override", URL: "/override"}},
		Concurrency:    8,
		Feeds:          []models.FeedConfig{{Slug: "override-feed"}},
		Extra:          map[string]any{"override": "value"},
		// These fields are intentionally non-zero to characterize that the
		// current top-level merge leaves them at their base values.
		ThemeCalendar: models.ThemeCalendarConfig{},
	}

	baseBefore := *base
	overrideBefore := *override
	got := MergeConfigs(base, override)

	want := *base
	want.OutputDir = override.OutputDir
	want.URL = override.URL
	want.Title = override.Title
	want.Description = override.Description
	want.Author = override.Author
	want.Language = override.Language
	want.AuthorURL = override.AuthorURL
	want.ManagingEditor = override.ManagingEditor
	want.WebMaster = override.WebMaster
	want.Copyright = override.Copyright
	want.License = override.License
	want.AssetsDir = override.AssetsDir
	want.TemplatesDir = override.TemplatesDir
	want.Hooks = override.Hooks
	want.DisabledHooks = override.DisabledHooks
	want.Nav = override.Nav
	want.Concurrency = override.Concurrency
	want.Feeds = override.Feeds
	want.Extra = map[string]any{"base": "value", "override": "value"}

	if !reflect.DeepEqual(*got, want) {
		t.Fatalf("merged config differs from characterized result\n got: %#v\nwant: %#v", *got, want)
	}
	if !reflect.DeepEqual(*base, baseBefore) {
		t.Fatal("MergeConfigs mutated base")
	}
	if !reflect.DeepEqual(*override, overrideBefore) {
		t.Fatal("MergeConfigs mutated override")
	}
}

func TestMergeConfigs_CharacterizesNilIdentity(t *testing.T) {
	base := &models.Config{OutputDir: "base"}
	override := &models.Config{OutputDir: "override"}
	if got := MergeConfigs(nil, override); got != override {
		t.Fatal("nil base should return override identity")
	}
	if got := MergeConfigs(base, nil); got != base {
		t.Fatal("nil override should return base identity")
	}
}
