package cmd

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestBlogrollMentionResolver(t *testing.T) {
	cfg := models.BlogrollConfig{
		Enabled: true,
		Feeds: []models.ExternalFeedConfig{{
			URL:     "https://example.com/feed.xml",
			SiteURL: "https://example.com",
			Handle:  "Poteto",
			Aliases: []string{"Lauren"},
		}, {
			URL: "https://simonwillison.net/atom/everything/",
		}},
	}
	r := newBlogrollMentionResolver(cfg)
	if r == nil {
		t.Fatal("expected resolver")
	}
	for _, handle := range []string{"poteto", "POTETO", "lauren", "example.com", "simonwillison", "simonwillison.net"} {
		if !r.ResolveHandle(handle) {
			t.Errorf("expected %q to resolve", handle)
		}
	}
	if r.ResolveHandle("not-in-blogroll") {
		t.Fatal("unknown handle resolved")
	}
	cfg.Enabled = false
	if newBlogrollMentionResolver(cfg) != nil {
		t.Fatal("disabled blogroll should not enable reference diagnostics")
	}
}
