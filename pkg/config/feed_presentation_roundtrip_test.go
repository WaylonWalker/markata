package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

// Exercise actual format converters: model-only unit tests never found this bug.
func TestFeedPresentationParsesAcrossFormats(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte) (*models.Config, error)
		data  string
	}{
		{
			name:  "toml",
			parse: ParseTOML,
			data:  "[markata-go.feed_defaults]\npresentation = \"timeline\"\n[[markata-go.feeds]]\nslug = \"archive\"\npresentation = \"editorial\"\npagination_type = \"htmx\"\n[[markata-go.feeds]]\nslug = \"ping\"\n",
		},
		{
			name:  "yaml",
			parse: ParseYAML,
			data:  "markata-go:\n  feed_defaults:\n    presentation: timeline\n  feeds:\n    - slug: archive\n      presentation: editorial\n      pagination_type: htmx\n    - slug: ping\n",
		},
		{
			name:  "json",
			parse: ParseJSON,
			data:  "{\"markata-go\":{\"feed_defaults\":{\"presentation\":\"timeline\"},\"feeds\":[{\"slug\":\"archive\",\"presentation\":\"editorial\",\"pagination_type\":\"htmx\"},{\"slug\":\"ping\"}]}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := tt.parse([]byte(tt.data))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if cfg.FeedDefaults.Presentation != "timeline" {
				t.Fatalf("default presentation = %q", cfg.FeedDefaults.Presentation)
			}
			if len(cfg.Feeds) != 2 {
				t.Fatalf("feeds = %d, want 2", len(cfg.Feeds))
			}
			archive := cfg.Feeds[0]
			archive.ApplyDefaults(cfg.FeedDefaults)
			if archive.Presentation != "editorial" || archive.PaginationType != models.PaginationHTMX {
				t.Fatalf("archive config lost presentation: %+v", archive)
			}
			ping := cfg.Feeds[1]
			ping.ApplyDefaults(cfg.FeedDefaults)
			if ping.Presentation != "timeline" {
				t.Fatalf("default presentation not inherited: %+v", ping)
			}
		})
	}
}

func TestIncludedTOMLFeedPresentationSurvivesMergedLoader(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "markata-go.toml")
	feeds := filepath.Join(dir, "feeds.toml")
	if err := os.WriteFile(base, []byte("[markata-go]\ninclude = [\"feeds.toml\"]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(feeds, []byte("[[markata-go.feeds]]\nslug = \"archive\"\npresentation = \"editorial\"\npagination_type = \"htmx\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadWithMerge(base)
	if err != nil {
		t.Fatalf("load includes: %v", err)
	}
	if len(cfg.Feeds) != 1 {
		t.Fatalf("feeds = %d, want 1", len(cfg.Feeds))
	}
	archive := cfg.Feeds[0]
	archive.ApplyDefaults(cfg.FeedDefaults)
	if archive.Presentation != "editorial" || archive.PaginationType != models.PaginationHTMX {
		t.Fatalf("merged config lost editorial: %+v", archive)
	}
}

func TestMergeFeedDefaultsPreservesPresentation(t *testing.T) {
	got := mergeFeedDefaults(models.NewFeedDefaults(), models.FeedDefaults{Presentation: "reading"})
	if got.Presentation != "reading" {
		t.Fatalf("merged presentation = %q", got.Presentation)
	}
	got = mergeFeedDefaults(got, models.FeedDefaults{})
	if got.Presentation != "reading" {
		t.Fatalf("empty override erased presentation: %q", got.Presentation)
	}
}
