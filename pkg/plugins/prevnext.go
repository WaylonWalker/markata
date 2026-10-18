package plugins

import (
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// PrevNextPlugin calculates previous/next post links for navigation.
// It runs in the collect stage after feeds have been created.
type PrevNextPlugin struct{}

// NewPrevNextPlugin creates a new PrevNextPlugin.
func NewPrevNextPlugin() *PrevNextPlugin {
	return &PrevNextPlugin{}
}

// Name returns the unique name of the plugin.
func (p *PrevNextPlugin) Name() string {
	return "prevnext"
}

// Priority ensures this plugin runs after the feeds plugin.
func (p *PrevNextPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCollect {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Collect calculates prev/next links for all posts based on the configured strategy.
func (p *PrevNextPlugin) Collect(m *lifecycle.Manager) error {
	config := getPrevNextConfig(m.Config())
	seriesCfg := parseSeriesConfig(m.Config())

	// Check if enabled
	if !config.Enabled {
		return nil
	}

	// Apply defaults
	config.ApplyDefaults()

	// Build feed lookup: post slug -> list of feeds containing the post
	feeds := m.Feeds()
	postToFeeds := buildPostToFeedsMap(feeds)

	// Process each post
	posts := m.Posts()
	for _, post := range posts {
		p.processPost(post, config, seriesCfg, feeds, postToFeeds)
	}

	return nil
}

// processPost calculates prev/next for a single post.
func (p *PrevNextPlugin) processPost(
	post *models.Post,
	config models.PrevNextConfig,
	seriesCfg seriesConfig,
	feeds []*lifecycle.Feed,
	postToFeeds map[string][]*lifecycle.Feed,
) {
	if post.PrevNextFeed != "" || post.PrevNextContext != nil {
		return
	}

	// Determine which feed to use for navigation
	feed := p.resolveFeed(post, config, seriesCfg, feeds, postToFeeds)
	if feed == nil {
		// Post is not in any feed, no prev/next navigation
		return
	}

	// Find the post's position in the feed
	position := -1
	for i, feedPost := range feed.Posts {
		if feedPost.Slug == post.Slug {
			position = i
			break
		}
	}

	if position == -1 {
		// Post not found in feed (shouldn't happen)
		return
	}

	// Set prev/next, skipping private posts to avoid linking to encrypted content
	var prev, next *models.Post
	for i := position - 1; i >= 0; i-- {
		if !feed.Posts[i].Private {
			prev = feed.Posts[i]
			break
		}
	}
	for i := position + 1; i < len(feed.Posts); i++ {
		if !feed.Posts[i].Private {
			next = feed.Posts[i]
			break
		}
	}

	// Update post fields
	post.Prev = prev
	post.Next = next
	post.PrevNextFeed = feed.Name
	post.PrevNextContext = &models.PrevNextContext{
		FeedSlug:  feed.Name,
		FeedTitle: feed.Title,
		Position:  position + 1, // 1-indexed
		Total:     len(feed.Posts),
		Prev:      prev,
		Next:      next,
	}
}

// resolveFeed determines which feed to use for a post based on the strategy.
func (p *PrevNextPlugin) resolveFeed(
	post *models.Post,
	config models.PrevNextConfig,
	seriesCfg seriesConfig,
	feeds []*lifecycle.Feed,
	postToFeeds map[string][]*lifecycle.Feed,
) *lifecycle.Feed {
	switch config.Strategy {
	case models.StrategyExplicitFeed:
		return findFeedBySlug(feeds, config.DefaultFeed)

	case models.StrategySeries:
		// Check for series in frontmatter
		seriesSlug := getStringFromExtra(post.Extra, "series_slug")
		if seriesSlug == "" {
			if series := getStringFromExtra(post.Extra, "series"); series != "" {
				seriesSlug = buildSeriesFeedSlug(seriesCfg.SlugPrefix, slugify(series))
			}
		}
		if seriesSlug != "" {
			if feed := findFeedBySlug(feeds, seriesSlug); feed != nil {
				return feed
			}
		}
		// Fall back to first_feed
		return p.getFirstFeed(post, postToFeeds)

	case models.StrategyFrontmatter:
		// Check for prevnext_feed in frontmatter
		if feedSlug := getStringFromExtra(post.Extra, "prevnext_feed"); feedSlug != "" {
			if feed := findFeedBySlug(feeds, feedSlug); feed != nil {
				return feed
			}
		}
		// Fall back to first_feed
		return p.getFirstFeed(post, postToFeeds)

	case models.StrategyFirstFeed:
		return p.getFirstFeed(post, postToFeeds)
	default:
		return p.getFirstFeed(post, postToFeeds)
	}
}

// getFirstFeed returns the first feed that contains the post.
func (p *PrevNextPlugin) getFirstFeed(
	post *models.Post,
	postToFeeds map[string][]*lifecycle.Feed,
) *lifecycle.Feed {
	feeds, ok := postToFeeds[post.Slug]
	if !ok || len(feeds) == 0 {
		return nil
	}
	return feeds[0]
}

// buildPostToFeedsMap creates a mapping from post slugs to the feeds they appear in.
func buildPostToFeedsMap(feeds []*lifecycle.Feed) map[string][]*lifecycle.Feed {
	result := make(map[string][]*lifecycle.Feed)

	for _, feed := range feeds {
		for _, post := range feed.Posts {
			result[post.Slug] = append(result[post.Slug], feed)
		}
	}

	return result
}

// findFeedBySlug finds a feed by its slug/name.
func findFeedBySlug(feeds []*lifecycle.Feed, slug string) *lifecycle.Feed {
	for _, feed := range feeds {
		if feed.Name == slug {
			return feed
		}
	}
	return nil
}

// getPrevNextConfig retrieves the prevnext configuration from the manager config.
func getPrevNextConfig(config *lifecycle.Config) models.PrevNextConfig {
	if config.Extra == nil {
		return models.NewPrevNextConfig()
	}

	if prevnextCfg, ok := config.Extra["prevnext"]; ok {
		if cfg, ok := prevnextCfg.(models.PrevNextConfig); ok {
			return cfg
		}
		// Try to parse from map
		if cfgMap, ok := prevnextCfg.(map[string]interface{}); ok {
			return parsePrevNextConfigFromMap(cfgMap)
		}
	}

	return models.NewPrevNextConfig()
}

// parsePrevNextConfigFromMap parses a PrevNextConfig from a generic map.
func parsePrevNextConfigFromMap(m map[string]interface{}) models.PrevNextConfig {
	config := models.NewPrevNextConfig()

	if enabled, ok := m["enabled"].(bool); ok {
		config.Enabled = enabled
	}

	if strategy, ok := m["strategy"].(string); ok {
		config.Strategy = models.PrevNextStrategy(strategy)
	}

	if defaultFeed, ok := m["default_feed"].(string); ok {
		config.DefaultFeed = defaultFeed
	}

	return config
}

// Ensure PrevNextPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*PrevNextPlugin)(nil)
	_ lifecycle.CollectPlugin  = (*PrevNextPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*PrevNextPlugin)(nil)
)
