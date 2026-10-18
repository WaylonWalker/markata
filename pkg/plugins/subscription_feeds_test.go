package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestNewSubscriptionFeedsPlugin(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()
	if plugin == nil {
		t.Fatal("NewSubscriptionFeedsPlugin() returned nil")
	}
	if plugin.Name() != "subscription_feeds" {
		t.Errorf("Name() = %q, want %q", plugin.Name(), "subscription_feeds")
	}
}

func TestSubscriptionFeedsPlugin_Priority(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()

	// Should run early in Collect stage
	if got := plugin.Priority(lifecycle.StageCollect); got != lifecycle.PriorityEarly {
		t.Errorf("Priority(StageCollect) = %d, want %d", got, lifecycle.PriorityEarly)
	}

	// Default priority for other stages
	if got := plugin.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestSubscriptionFeedsPlugin_Collect_InjectsFeedConfigs(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title":       "Test Site",
		"description": "A test site",
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	// Run Collect
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Check that feed_configs were added to cache
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}

	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatal("feed_configs is not []models.FeedConfig")
	}

	// Should have root and archive feeds
	if len(feedConfigs) < 2 {
		t.Errorf("Expected at least 2 feed configs, got %d", len(feedConfigs))
	}

	// Check for root feed (slug="")
	foundRoot := false
	foundArchive := false
	for _, fc := range feedConfigs {
		if fc.Slug == "" {
			foundRoot = true
			// Verify root feed settings
			if fc.Title != "Test Site Feed" {
				t.Errorf("Root feed title = %q, want %q", fc.Title, "Test Site Feed")
			}
			if !fc.Formats.HTML {
				t.Error("Root feed should have HTML=true")
			}
			if !fc.Formats.RSS {
				t.Error("Root feed should have RSS=true")
			}
			if !fc.Formats.Atom {
				t.Error("Root feed should have Atom=true")
			}
		}
		if fc.Slug == "archive" {
			foundArchive = true
			// Verify archive feed settings
			if fc.Title != "Test Site Archive Feed" {
				t.Errorf("Archive feed title = %q, want %q", fc.Title, "Test Site Archive Feed")
			}
			if fc.Formats.HTML {
				t.Error("Archive feed should have HTML=false")
			}
			if !fc.Formats.RSS {
				t.Error("Archive feed should have RSS=true")
			}
			if !fc.Formats.Atom {
				t.Error("Archive feed should have Atom=true")
			}
		}
	}

	if !foundRoot {
		t.Error("Root subscription feed (slug='') not found")
	}
	if !foundArchive {
		t.Error("Archive subscription feed (slug='archive') not found")
	}
}

func TestSubscriptionFeedsPlugin_Collect_PreservesConfiguredArchive(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()
	archive := models.FeedConfig{
		Slug:  "archive",
		Title: "Configured Archive",
		Formats: models.FeedFormats{
			HTML: true,
		},
	}
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"feeds": []models.FeedConfig{archive},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	feedConfigs := getFeedConfigs(m.Config())
	if len(feedConfigs) != 3 {
		t.Fatalf("feed configs = %#v, want configured archive, implicit root, and implicit pins", feedConfigs)
	}
	seen := make(map[string]bool, len(feedConfigs))
	for _, fc := range feedConfigs {
		if seen[fc.Slug] {
			t.Errorf("duplicate feed %q", fc.Slug)
		}
		seen[fc.Slug] = true
		switch fc.Slug {
		case "":
			if !fc.Formats.HTML {
				t.Error("implicit root feed should render the homepage")
			}
		case "archive":
			if fc.Title != archive.Title || !fc.Formats.HTML {
				t.Errorf("archive = %#v, want configured archive %#v", fc, archive)
			}
		case "pins":
			if fc.Filter != "published == true and (link or url)" || !fc.Formats.HTML || fc.Templates.HTML != "pins.html" {
				t.Errorf("pins = %#v, want implicit HTML link board", fc)
			}
		default:
			t.Errorf("unexpected feed %q", fc.Slug)
		}
	}
}

func TestSubscriptionFeedsPlugin_Collect_DoesNotDuplicateExisting(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title": "Test Site",
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	// Pre-populate cache with existing root feed
	existingFeeds := []models.FeedConfig{
		{
			Slug:  "",
			Title: "Custom Root Feed",
			Formats: models.FeedFormats{
				HTML: true,
				RSS:  true,
			},
		},
	}
	m.Cache().Set("feed_configs", existingFeeds)

	// Run Collect
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Check feed configs
	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("Expected feed_configs in cache")
	}
	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatal("Expected feed_configs to be []models.FeedConfig")
	}

	// Count root feeds
	rootCount := 0
	for _, fc := range feedConfigs {
		if fc.Slug == "" {
			rootCount++
			// Should preserve the existing custom root feed
			if fc.Title != "Custom Root Feed" {
				t.Errorf("Expected custom root feed to be preserved, got title %q", fc.Title)
			}
		}
	}

	if rootCount != 1 {
		t.Errorf("Expected exactly 1 root feed, got %d", rootCount)
	}
}

func TestSubscriptionFeedsPlugin_Collect_CanDisableSiteArchive(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()
	defaults := models.NewFeedDefaults()
	defaults.Syndication.SiteArchiveDisabled = true

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title":         "Test Site",
		"feed_defaults": defaults,
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	if err := plugin.Collect(m); err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	cached, ok := m.Cache().Get("feed_configs")
	if !ok {
		t.Fatal("feed_configs not found in cache")
	}
	feedConfigs, ok := cached.([]models.FeedConfig)
	if !ok {
		t.Fatal("feed_configs is not []models.FeedConfig")
	}
	for _, fc := range feedConfigs {
		if fc.Slug == "archive" {
			t.Fatal("archive subscription feed should not be injected when disabled")
		}
	}
}

func TestSubscriptionFeedsPlugin_Collect_Disabled(t *testing.T) {
	plugin := NewSubscriptionFeedsPlugin()

	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"title":                       "Test Site",
		"subscription_feeds_disabled": true,
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	// Run Collect
	err := plugin.Collect(m)
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	// Check that no feed_configs were added
	_, ok := m.Cache().Get("feed_configs")
	if ok {
		t.Error("Expected no feed_configs when subscription_feeds_disabled=true")
	}
}

func TestGetDiscoveryFeed_WithSidebarFeed(t *testing.T) {
	sidebarFeed := &models.FeedConfig{
		Slug:  "tags/python",
		Title: "Python Posts",
		Formats: models.FeedFormats{
			RSS:  true,
			Atom: true,
			JSON: false,
		},
	}

	allFeeds := []models.FeedConfig{
		{Slug: "", Title: "Site Feed"},
	}

	post := &models.Post{Slug: "test-post"}

	discovery := GetDiscoveryFeed(post, sidebarFeed, allFeeds)

	if discovery.Slug != "tags/python" {
		t.Errorf("Discovery slug = %q, want %q", discovery.Slug, "tags/python")
	}
	if discovery.RSSURL != "/tags/python/rss.xml" {
		t.Errorf("RSS URL = %q, want %q", discovery.RSSURL, "/tags/python/rss.xml")
	}
	if discovery.AtomURL != "/tags/python/atom.xml" {
		t.Errorf("Atom URL = %q, want %q", discovery.AtomURL, "/tags/python/atom.xml")
	}
	if discovery.HasRSS != true {
		t.Error("Expected HasRSS=true")
	}
	if discovery.HasAtom != true {
		t.Error("Expected HasAtom=true")
	}
	if discovery.HasJSON != false {
		t.Error("Expected HasJSON=false")
	}
}

func TestGetDiscoveryFeed_WithoutSidebarFeed(t *testing.T) {
	allFeeds := []models.FeedConfig{
		{
			Slug:  "",
			Title: "Site Feed",
			Formats: models.FeedFormats{
				RSS:  true,
				Atom: true,
			},
		},
		{
			Slug:  "archive",
			Title: "Archive",
		},
	}

	post := &models.Post{Slug: "test-post"}

	discovery := GetDiscoveryFeed(post, nil, allFeeds)

	if discovery.Slug != "" {
		t.Errorf("Discovery slug = %q, want empty string", discovery.Slug)
	}
	if discovery.Title != "Site Feed" {
		t.Errorf("Title = %q, want %q", discovery.Title, "Site Feed")
	}
	if discovery.RSSURL != "/rss.xml" {
		t.Errorf("RSS URL = %q, want %q", discovery.RSSURL, "/rss.xml")
	}
	if discovery.AtomURL != "/atom.xml" {
		t.Errorf("Atom URL = %q, want %q", discovery.AtomURL, "/atom.xml")
	}
}

func TestGetDiscoveryFeed_Fallback(t *testing.T) {
	// No feeds configured at all
	allFeeds := []models.FeedConfig{}

	post := &models.Post{Slug: "test-post"}

	discovery := GetDiscoveryFeed(post, nil, allFeeds)

	// Should return default fallback
	if discovery.RSSURL != "/rss.xml" {
		t.Errorf("RSS URL = %q, want %q", discovery.RSSURL, "/rss.xml")
	}
	if discovery.AtomURL != "/atom.xml" {
		t.Errorf("Atom URL = %q, want %q", discovery.AtomURL, "/atom.xml")
	}
}

func TestDiscoveryFeedToMap(t *testing.T) {
	df := &DiscoveryFeed{
		Slug:    "tags/go",
		Title:   "Go Posts",
		RSSURL:  "/tags/go/rss.xml",
		AtomURL: "/tags/go/atom.xml",
		JSONURL: "",
		HasRSS:  true,
		HasAtom: true,
		HasJSON: false,
	}

	m := DiscoveryFeedToMap(df)

	if m["slug"] != "tags/go" {
		t.Errorf("map[slug] = %v, want %q", m["slug"], "tags/go")
	}
	if m["title"] != "Go Posts" {
		t.Errorf("map[title] = %v, want %q", m["title"], "Go Posts")
	}
	if m["rss_url"] != "/tags/go/rss.xml" {
		t.Errorf("map[rss_url] = %v, want %q", m["rss_url"], "/tags/go/rss.xml")
	}
	if m["has_rss"] != true {
		t.Errorf("map[has_rss] = %v, want true", m["has_rss"])
	}
	if m["has_json"] != false {
		t.Errorf("map[has_json] = %v, want false", m["has_json"])
	}
}

func TestDiscoveryFeedToMap_Nil(t *testing.T) {
	m := DiscoveryFeedToMap(nil)
	if m != nil {
		t.Errorf("DiscoveryFeedToMap(nil) = %v, want nil", m)
	}
}

func TestGetFeedBySlug(t *testing.T) {
	feeds := []models.FeedConfig{
		{Slug: "", Title: "Root"},
		{Slug: "tags/python", Title: "Python"},
		{Slug: "archive", Title: "Archive"},
	}

	// Found cases
	if fc := GetFeedBySlug("tags/python", feeds); fc == nil || fc.Title != "Python" {
		t.Error("Expected to find tags/python feed")
	}
	if fc := GetFeedBySlug("", feeds); fc == nil || fc.Title != "Root" {
		t.Error("Expected to find root feed")
	}

	// Not found case
	if fc := GetFeedBySlug("nonexistent", feeds); fc != nil {
		t.Error("Expected nil for nonexistent feed")
	}
}

func TestFindPostSidebarFeed_TagMatch(t *testing.T) {
	enabled := true
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"components": models.ComponentsConfig{
				FeedSidebar: models.FeedSidebarConfig{
					Enabled: &enabled,
					Feeds:   []string{"tags/daily-note"},
				},
			},
		},
	}

	feeds := []models.FeedConfig{
		{Slug: "tags/daily-note", Title: "Daily Notes"},
	}

	// Post with matching tag
	post := &models.Post{
		Slug: "my-note",
		Tags: []string{"daily-note", "golang"},
	}

	fc := FindPostSidebarFeed(post, config, feeds)
	if fc == nil {
		t.Fatal("Expected to find sidebar feed for post with matching tag")
	}
	if fc.Slug != "tags/daily-note" {
		t.Errorf("Feed slug = %q, want %q", fc.Slug, "tags/daily-note")
	}
}

func TestFindPostSidebarFeed_NoMatch(t *testing.T) {
	enabled := true
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"components": models.ComponentsConfig{
				FeedSidebar: models.FeedSidebarConfig{
					Enabled: &enabled,
					Feeds:   []string{"tags/daily-note"},
				},
			},
		},
	}

	feeds := []models.FeedConfig{
		{Slug: "tags/daily-note", Title: "Daily Notes"},
	}

	// Post without matching tag
	post := &models.Post{
		Slug: "my-post",
		Tags: []string{"golang", "programming"},
	}

	fc := FindPostSidebarFeed(post, config, feeds)
	if fc != nil {
		t.Errorf("Expected nil for post without matching tag, got %+v", fc)
	}
}

func TestFindPostSidebarFeed_Disabled(t *testing.T) {
	disabled := false
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"components": models.ComponentsConfig{
				FeedSidebar: models.FeedSidebarConfig{
					Enabled: &disabled,
					Feeds:   []string{"tags/daily-note"},
				},
			},
		},
	}

	feeds := []models.FeedConfig{
		{Slug: "tags/daily-note", Title: "Daily Notes"},
	}

	post := &models.Post{
		Slug: "my-note",
		Tags: []string{"daily-note"},
	}

	fc := FindPostSidebarFeed(post, config, feeds)
	if fc != nil {
		t.Error("Expected nil when feed sidebar is disabled")
	}
}
