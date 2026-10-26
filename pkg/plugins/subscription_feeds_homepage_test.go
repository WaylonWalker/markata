package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func cachedRootFeed(t *testing.T, m *lifecycle.Manager) *models.FeedConfig {
	t.Helper()

	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}
	feeds, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatalf("feed_configs type = %T, want []models.FeedConfig", cached)
	}
	root := GetFeedBySlug("", feeds)
	if root == nil {
		t.Fatal("root feed not found")
	}
	return root
}

func TestSubscriptionFeedsPlugin_Collect_AuthoredHomepageKeepsRootSyndicationWithoutFeedHTML(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title": "Test Site",
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)
	m.SetPosts([]*models.Post{
		{Path: "index.md", Slug: "", Published: true},
		{Path: "post.md", Slug: "post", Published: true},
	})

	if err := NewSubscriptionFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	root := cachedRootFeed(t, m)
	if root.Formats.HTML {
		t.Error("implicit root feed HTML should be disabled when an authored homepage owns /index.html")
	}
	if !root.Formats.RSS || !root.Formats.Atom {
		t.Fatalf("root syndication should remain enabled: formats = %#v", root.Formats)
	}
}

func TestSubscriptionFeedsPlugin_Collect_DraftHomepageDoesNotSuppressDefaultHomepage(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title": "Test Site",
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)
	m.SetPosts([]*models.Post{
		{Path: "index.md", Slug: "", Draft: true},
	})

	if err := NewSubscriptionFeedsPlugin().Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	root := cachedRootFeed(t, m)
	if !root.Formats.HTML {
		t.Error("draft homepage should not suppress the generated default homepage")
	}
}

func TestSubscriptionFeedsPlugin_Collect_ReevaluatesImplicitRootAcrossRebuilds(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title": "Test Site",
	}
	m := lifecycle.NewManager()
	m.SetConfig(config)

	m.SetPosts([]*models.Post{{Path: "post.md", Slug: "post", Published: true}})
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("first Collect() error = %v", err)
	}
	if root := cachedRootFeed(t, m); !root.Formats.HTML {
		t.Fatal("initial implicit root should render generated homepage HTML")
	}

	m.SetPosts([]*models.Post{
		{Path: "index.md", Slug: "", Published: true},
		{Path: "post.md", Slug: "post", Published: true},
	})
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() after adding homepage error = %v", err)
	}
	root := cachedRootFeed(t, m)
	if root.Formats.HTML {
		t.Error("implicit root HTML should turn off when an authored homepage appears")
	}
	if !root.Formats.RSS || !root.Formats.Atom {
		t.Fatalf("root syndication should remain enabled after homepage appears: %#v", root.Formats)
	}

	m.SetPosts([]*models.Post{{Path: "post.md", Slug: "post", Published: true}})
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() after removing homepage error = %v", err)
	}
	if root := cachedRootFeed(t, m); !root.Formats.HTML {
		t.Error("implicit root HTML should turn back on when authored homepage disappears")
	}
}

func TestSubscriptionFeedsPlugin_Collect_ConfiguredRootRemainsAuthoritative(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title": "Test Site",
	}
	m := lifecycle.NewManager()
	m.SetConfig(config)

	// First collect creates an implicit root and records its runtime provenance.
	if err := plugin.Collect(m); err != nil {
		t.Fatalf("initial Collect() error = %v", err)
	}

	explicit := models.FeedConfig{
		Slug:  "",
		Title: "Configured Home",
		Formats: models.FeedFormats{
			HTML: true,
			RSS:  true,
		},
	}
	modelsConfig := &models.Config{Feeds: []models.FeedConfig{explicit}}
	config.Extra["models_config"] = modelsConfig
	config.Extra["feeds"] = []models.FeedConfig{explicit}
	m.SetPosts([]*models.Post{{Path: "index.md", Slug: "", Published: true}})

	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() with configured root error = %v", err)
	}
	root := cachedRootFeed(t, m)
	if root.Title != explicit.Title || !root.Formats.HTML {
		t.Fatalf("configured root was changed: %#v", root)
	}
}
