// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

const (
	rootSubscriptionFeedType = "root"
	pinsFeedSlug             = "pins"
)

// SubscriptionFeedsPlugin creates built-in subscription feeds at root, /archive, and /pins.
// The root feed also renders the default homepage, while /archive retains the
// configured archive page and both locations expose RSS and Atom feeds.
//
// This plugin also computes the discovery feed for each post during rendering,
// enabling dynamic <link rel="alternate"> tags based on the post's sidebar feed.
type SubscriptionFeedsPlugin struct {
	implicitRoot bool
}

// NewSubscriptionFeedsPlugin creates a new SubscriptionFeedsPlugin.
func NewSubscriptionFeedsPlugin() *SubscriptionFeedsPlugin {
	return &SubscriptionFeedsPlugin{}
}

// Name returns the unique name of the plugin.
func (p *SubscriptionFeedsPlugin) Name() string {
	return "subscription_feeds"
}

// Priority returns the plugin's priority for a given stage.
func (p *SubscriptionFeedsPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageCollect:
		// Run early in Collect to inject subscription feeds before other feed processing
		return lifecycle.PriorityEarly
	default:
		return lifecycle.PriorityDefault
	}
}

// Collect injects built-in subscription feeds into the feed configs.
// The implicit root feed generates an HTML homepage as well as RSS and Atom.
// Explicit root, archive, or pins feed configuration remains authoritative.
func (p *SubscriptionFeedsPlugin) Collect(m *lifecycle.Manager) error {
	config := m.Config()
	syndication := getSyndicationConfig(config)

	// Check if subscription feeds are disabled
	if config.Extra != nil {
		if disabled, ok := config.Extra["subscription_feeds_disabled"].(bool); ok && disabled {
			return nil
		}
	}

	// Start from the resolved config so built-in feeds (including /archive/)
	// and feeds injected earlier in Collect (e.g. series) are present. Cached
	// configs may predate those injections, so they only add feeds whose slug
	// the resolved config does not already define.
	feedConfigs := getFeedConfigs(config)
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			known := make(map[string]bool, len(feedConfigs))
			for i := range feedConfigs {
				known[feedConfigs[i].Slug] = true
			}
			for i := range fcs {
				if !known[fcs[i].Slug] {
					known[fcs[i].Slug] = true
					feedConfigs = append(feedConfigs, fcs[i])
				}
			}
		}
	}

	// Check if built-in feeds already exist.
	rootFeedIndex := -1
	hasArchiveFeed := false
	for i := range feedConfigs {
		if feedConfigs[i].Slug == "" {
			rootFeedIndex = i
		}
		if feedConfigs[i].Slug == defaultArchivePrefix {
			hasArchiveFeed = true
		}
	}

	authoredHomepage := hasAuthoredHomepage(m.Posts())
	configuredRoot := hasConfiguredRootFeed(config)

	// A root feed we injected on a previous Collect call is runtime state, not
	// explicit site configuration. Re-evaluate whether it should own index.html
	// so serve/incremental rebuilds react when an authored homepage appears or
	// disappears. A configured root feed always remains authoritative.
	if rootFeedIndex >= 0 && p.implicitRoot && !configuredRoot {
		feedConfigs[rootFeedIndex].Formats.HTML = !authoredHomepage
	}
	if configuredRoot {
		p.implicitRoot = false
	}

	// Create root subscription feed (slug="") if not already defined. An authored
	// homepage keeps ownership of /index.html while the implicit feed still
	// provides RSS and Atom at the root.
	if rootFeedIndex < 0 {
		rootFeed := models.FeedConfig{
			Slug:        "",
			Title:       getSubscriptionFeedTitle(config, rootSubscriptionFeedType),
			Description: getSubscriptionFeedDescription(config, rootSubscriptionFeedType),
			Filter:      "published == true",
			Sort:        "date",
			Reverse:     true,
			Formats: models.FeedFormats{
				HTML: !authoredHomepage,
				RSS:  true,
				Atom: true,
				JSON: false,
			},
		}
		feedConfigs = append(feedConfigs, rootFeed)
		p.implicitRoot = true
	}

	// Create archive subscription feed (slug="archive") if not already defined
	if !hasArchiveFeed && !syndication.SiteArchiveDisabled {
		archiveFeed := models.FeedConfig{
			Slug:        defaultArchivePrefix,
			Title:       getSubscriptionFeedTitle(config, defaultArchivePrefix),
			Description: getSubscriptionFeedDescription(config, defaultArchivePrefix),
			Filter:      "published == true",
			Sort:        "date",
			Reverse:     true,
			Formats: models.FeedFormats{
				HTML: false, // Don't generate index.html at /archive
				RSS:  true,
				Atom: true,
				JSON: false,
			},
		}
		feedConfigs = append(feedConfigs, archiveFeed)
	}

	// Create a visual pins feed unless the site configured one explicitly or an
	// authored publishable post already owns /pins/. Link frontmatter remains
	// render-neutral outside this feed.
	if shouldAddImplicitPinsFeed(feedConfigs, m.Posts()) {
		pinsFeed := models.FeedConfig{
			Slug:           pinsFeedSlug,
			ItemsPerPage:   100,
			PaginationType: models.PaginationManual,
			Title:          "Pins",
			Description:    "A field notebook of saved links",
			Filter:         "published == true and (link or url)",
			Sort:           "date",
			Reverse:        true,
			Templates: models.FeedTemplates{
				HTML: "pins.html",
			},
			Formats: models.FeedFormats{
				HTML: true,
				RSS:  false,
				Atom: false,
				JSON: false,
			},
		}
		feedConfigs = append(feedConfigs, pinsFeed)
	}

	// FeedsPlugin reads the resolved config during collection. Keep it in sync
	// with the cache so the injected feeds are collected on this build.
	config.Extra["feeds"] = feedConfigs
	m.Cache().Set("feed_configs", feedConfigs)

	return nil
}

func hasAuthoredHomepage(posts []*models.Post) bool {
	for _, post := range posts {
		if post == nil || post.Skip || post.Draft {
			continue
		}
		if post.Slug == "" {
			return true
		}
	}
	return false
}

func hasConfiguredRootFeed(config *lifecycle.Config) bool {
	modelsConfig, ok := getModelsConfig(config)
	if !ok || modelsConfig == nil {
		return false
	}
	for i := range modelsConfig.Feeds {
		if modelsConfig.Feeds[i].Slug == "" {
			return true
		}
	}
	return false
}

func shouldAddImplicitPinsFeed(feedConfigs []models.FeedConfig, posts []*models.Post) bool {
	for i := range feedConfigs {
		if feedConfigs[i].Slug == pinsFeedSlug {
			return false
		}
	}
	return !postOwnsFeedSlug(posts, pinsFeedSlug)
}

func postOwnsFeedSlug(posts []*models.Post, slug string) bool {
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if slug == "" {
		return false
	}
	for _, post := range posts {
		if post == nil || post.Skip || post.Draft || !post.Published {
			continue
		}
		if strings.Trim(strings.TrimSpace(post.Slug), "/") == slug {
			return true
		}
	}
	return false
}

// getSubscriptionFeedTitle returns the title for a subscription feed.
func getSubscriptionFeedTitle(config *lifecycle.Config, feedType string) string {
	siteTitle := "Site"
	if config.Extra != nil {
		if title, ok := config.Extra["title"].(string); ok && title != "" {
			siteTitle = title
		}
	}

	switch feedType {
	case rootSubscriptionFeedType:
		return siteTitle + " Feed"
	case defaultArchivePrefix:
		return siteTitle + " Archive Feed"
	default:
		return siteTitle + " Feed"
	}
}

// getSubscriptionFeedDescription returns the description for a subscription feed.
func getSubscriptionFeedDescription(config *lifecycle.Config, feedType string) string {
	siteDescription := ""
	if config.Extra != nil {
		if desc, ok := config.Extra["description"].(string); ok {
			siteDescription = desc
		}
	}

	if siteDescription != "" {
		return siteDescription
	}

	switch feedType {
	case rootSubscriptionFeedType:
		return "All published posts"
	case defaultArchivePrefix:
		return "Archive of all published posts"
	default:
		return "Posts feed"
	}
}

// DiscoveryFeed represents feed discovery information for templates.
// This is injected into template context to enable per-page feed discovery.
type DiscoveryFeed struct {
	// Slug is the feed slug (e.g., "tags/python", "archive", "")
	Slug string

	// Title is the feed title for display
	Title string

	// RSSURL is the RSS feed URL (if RSS format is enabled)
	RSSURL string

	// AtomURL is the Atom feed URL (if Atom format is enabled)
	AtomURL string

	// JSONURL is the JSON feed URL (if JSON format is enabled)
	JSONURL string

	// HasRSS indicates whether RSS format is enabled
	HasRSS bool

	// HasAtom indicates whether Atom format is enabled
	HasAtom bool

	// HasJSON indicates whether JSON format is enabled
	HasJSON bool
}

// GetDiscoveryFeed returns the discovery feed for a post.
// If the post has a sidebar_feed, that feed is used for discovery.
// Otherwise, the site default feed (root subscription feed) is used.
//
// The post parameter is reserved for future use when discovery logic
// may need to inspect post metadata (e.g., explicit feed assignment).
// This function is called from templates.go renderPost to inject discovery_feed context.
func GetDiscoveryFeed(_ *models.Post, sidebarFeed *models.FeedConfig, allFeeds []models.FeedConfig) *DiscoveryFeed {
	// If post has a sidebar feed, use that for discovery
	if sidebarFeed != nil {
		return feedConfigToDiscoveryFeed(sidebarFeed)
	}

	// Otherwise, use the root subscription feed (slug="")
	for i := range allFeeds {
		if allFeeds[i].Slug == "" {
			return feedConfigToDiscoveryFeed(&allFeeds[i])
		}
	}

	// Fallback: return a default discovery feed pointing to root feeds
	return &DiscoveryFeed{
		Slug:    "",
		Title:   "Site Feed",
		RSSURL:  "/rss.xml",
		AtomURL: "/atom.xml",
		HasRSS:  true,
		HasAtom: true,
		HasJSON: false,
	}
}

// feedConfigToDiscoveryFeed converts a FeedConfig to a DiscoveryFeed.
func feedConfigToDiscoveryFeed(fc *models.FeedConfig) *DiscoveryFeed {
	df := &DiscoveryFeed{
		Slug:    fc.Slug,
		Title:   fc.Title,
		HasRSS:  fc.Formats.RSS,
		HasAtom: fc.Formats.Atom,
		HasJSON: fc.Formats.JSON,
	}

	// Generate feed URLs based on slug
	baseURL := ""
	if fc.Slug != "" {
		baseURL = "/" + fc.Slug
	}

	if df.HasRSS {
		df.RSSURL = baseURL + "/rss.xml"
	}
	if df.HasAtom {
		df.AtomURL = baseURL + "/atom.xml"
	}
	if df.HasJSON {
		df.JSONURL = baseURL + "/feed.json"
	}

	return df
}

// DiscoveryFeedToMap converts a DiscoveryFeed to a map for template context.
func DiscoveryFeedToMap(df *DiscoveryFeed) map[string]interface{} {
	if df == nil {
		return nil
	}
	return map[string]interface{}{
		"slug":     df.Slug,
		"title":    df.Title,
		"rss_url":  df.RSSURL,
		"atom_url": df.AtomURL,
		"json_url": df.JSONURL,
		"has_rss":  df.HasRSS,
		"has_atom": df.HasAtom,
		"has_json": df.HasJSON,
	}
}

// GetFeedBySlug finds a feed config by slug from the feed configs list.
func GetFeedBySlug(slug string, feedConfigs []models.FeedConfig) *models.FeedConfig {
	for i := range feedConfigs {
		if feedConfigs[i].Slug == slug {
			return &feedConfigs[i]
		}
	}
	return nil
}

// FindPostSidebarFeed finds the sidebar feed for a post.
// It checks if the post belongs to any feed configured for sidebar display.
// Returns the feed config and a boolean indicating if a match was found.
func FindPostSidebarFeed(post *models.Post, config *lifecycle.Config, feedConfigs []models.FeedConfig) *models.FeedConfig {
	seriesCfg := parseSeriesConfig(config)

	// Get components config
	components, ok := config.Extra["components"].(models.ComponentsConfig)
	if !ok {
		return nil
	}

	// Check if feed sidebar is enabled
	if components.FeedSidebar.Enabled == nil || !*components.FeedSidebar.Enabled {
		return nil
	}

	// Get configured feed slugs
	feedSlugs := components.FeedSidebar.Feeds
	if len(feedSlugs) == 0 {
		return nil
	}

	// Check if this post belongs to any of the configured feeds
	for _, feedSlug := range feedSlugs {
		// Handle tag-based feeds (tags/xxx)
		if strings.HasPrefix(feedSlug, "tags/") {
			tagName := strings.TrimPrefix(feedSlug, "tags/")
			for _, postTag := range post.Tags {
				if postTag == tagName {
					// Found a matching tag feed
					return GetFeedBySlug(feedSlug, feedConfigs)
				}
			}
		}

		// Handle explicit feed membership via frontmatter
		if post.PrevNextFeed == feedSlug {
			return GetFeedBySlug(feedSlug, feedConfigs)
		}
		if feed, ok := post.Extra["feed"].(string); ok && feed == feedSlug {
			return GetFeedBySlug(feedSlug, feedConfigs)
		}

		// Handle series-based feeds
		if series, ok := post.Extra["series"].(string); ok {
			seriesSlug := getStringFromExtra(post.Extra, "series_slug")
			if seriesSlug == "" {
				seriesSlug = buildSeriesFeedSlug(seriesCfg.SlugPrefix, models.Slugify(series))
			}
			if feedSlug == seriesSlug {
				return GetFeedBySlug(feedSlug, feedConfigs)
			}
		}
	}

	return nil
}

// Ensure SubscriptionFeedsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*SubscriptionFeedsPlugin)(nil)
	_ lifecycle.CollectPlugin  = (*SubscriptionFeedsPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*SubscriptionFeedsPlugin)(nil)
)
