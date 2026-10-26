package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSimpleHTMLViewEnabledRequiresFormatAndView(t *testing.T) {
	fc := &models.FeedConfig{Formats: models.FeedFormats{SimpleHTML: true}, Views: []string{models.FeedViewDefault}}
	if simpleHTMLViewEnabled(fc) {
		t.Fatal("simple output enabled after view opt-out")
	}
	fc.Views = []string{models.FeedViewDefault, models.FeedViewSimple}
	if !simpleHTMLViewEnabled(fc) {
		t.Fatal("simple output should be enabled")
	}
}

func TestSidebarVariantsHideDisabledSimpleView(t *testing.T) {
	fc := &models.FeedConfig{Slug: "blog", Views: []string{models.FeedViewDefault}, Formats: models.FeedFormats{SimpleHTML: true}}
	variants := buildSidebarVariants(fc, models.SyndicationConfig{}, models.PostFormatsConfig{})
	for _, variant := range variants {
		if variant.Key == "simple" {
			t.Fatal("sidebar exposed disabled simple view")
		}
	}
}

func TestPublishFeedViewChangeRemovesStaleSimpleOutput(t *testing.T) {
	output := t.TempDir()
	cfg := lifecycle.NewConfig()
	cfg.OutputDir = output
	plugin := NewPublishFeedsPlugin()
	feed := &models.FeedConfig{
		Slug: "blog", Title: "Blog",
		Posts:   []*models.Post{{Slug: "post", Href: "/post/", Published: true, ArticleHTML: "<p>Post</p>"}},
		Views:   []string{models.FeedViewDefault, models.FeedViewSimple},
		Formats: models.FeedFormats{HTML: true, SimpleHTML: true},
	}
	if err := plugin.publishFeed(feed, cfg, output); err != nil {
		t.Fatal(err)
	}
	simple := filepath.Join(output, "blog", "simple", "index.html")
	if _, err := os.Stat(simple); err != nil {
		t.Fatalf("simple output missing: %v", err)
	}
	hash := plugin.computeFeedHash(feed)
	feed.Views = []string{models.FeedViewDefault}
	if plugin.computeFeedHash(feed) == hash {
		t.Fatal("view change did not invalidate feed hash")
	}
	if err := plugin.publishFeed(feed, cfg, output); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(simple); !os.IsNotExist(err) {
		t.Fatalf("disabled simple output remains: %v", err)
	}
	if _, err := os.Stat(filepath.Join(output, "blog", "index.html")); err != nil {
		t.Fatalf("primary output missing: %v", err)
	}
}
