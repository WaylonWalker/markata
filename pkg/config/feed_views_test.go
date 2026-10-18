package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestFeedViewsParseAcrossFormats(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte) (*models.Config, error)
		data  string
	}{
		{"toml", ParseTOML, `[markata-go.feed_defaults]
views = ["default", "calendar"]
[[markata-go.feeds]]
slug = "notes"
views = ["default", "simple"]
`},
		{"yaml", ParseYAML, "markata-go:\n  feed_defaults:\n    views: [default, calendar]\n  feeds:\n    - slug: notes\n      views: [default, simple]\n"},
		{"json", ParseJSON, `{"markata-go":{"feed_defaults":{"views":["default","calendar"]},"feeds":[{"slug":"notes","views":["default","simple"]}]}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := tt.parse([]byte(tt.data))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}
			if !reflect.DeepEqual(cfg.FeedDefaults.Views, []string{"default", "calendar"}) {
				t.Fatalf("defaults views = %#v", cfg.FeedDefaults.Views)
			}
			if len(cfg.Feeds) != 1 || !reflect.DeepEqual(cfg.Feeds[0].Views, []string{"default", "simple"}) {
				t.Fatalf("feed views = %#v", cfg.Feeds)
			}
		})
	}
}

func TestValidateFeedViews(t *testing.T) {
	errs := validateFeedViews("feeds[0].views", []string{"simple", "unknown"})
	if len(errs) != 2 {
		t.Fatalf("errors = %d, want 2", len(errs))
	}
	joined := errs[0].Error() + "\n" + errs[1].Error()
	if !strings.Contains(joined, "unknown feed view") || !strings.Contains(joined, `must include "default"`) {
		t.Fatalf("unexpected validation: %s", joined)
	}
}

func TestMergeFeedDefaultsPreservesExplicitEmptyViews(t *testing.T) {
	got := mergeFeedDefaults(models.NewFeedDefaults(), models.FeedDefaults{Views: []string{}})
	if got.Views == nil {
		t.Fatal("explicit empty views became inherited views")
	}
	if len(validateFeedViews("feed_defaults.views", got.Views)) == 0 {
		t.Fatal("merged empty view list escaped validation")
	}
}
