// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// seriesKey is the frontmatter and config key for series.
const seriesKey = "series"

func buildSeriesFeedSlug(prefix, seriesSlug string) string {
	trimmed := strings.Trim(prefix, "/")
	if trimmed == "" {
		return seriesSlug
	}
	return trimmed + "/" + seriesSlug
}

// SeriesPlugin scans posts for `series` frontmatter and auto-generates
// series feed configs. It runs early in the Collect stage so that the
// feeds plugin can process the generated configs.
//
// # Processing Steps
//
//  1. Scan all posts for `series` frontmatter
//  2. Group posts by series name
//  3. Sort posts within each series (by series_order or date ascending)
//  4. Set guide navigation (Prev/Next) on each post
//  5. Set PrevNextContext with position info
//  6. Inject series FeedConfigs into config.Extra["feeds"]
type SeriesPlugin struct{}

// NewSeriesPlugin creates a new SeriesPlugin.
func NewSeriesPlugin() *SeriesPlugin {
	return &SeriesPlugin{}
}

// Name returns the unique name of the plugin.
func (p *SeriesPlugin) Name() string {
	return seriesKey
}

// Priority returns the plugin's priority for a given stage.
// Series runs early in Collect so feed configs are ready for the feeds plugin.
func (p *SeriesPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCollect {
		return lifecycle.PriorityEarly
	}
	return lifecycle.PriorityDefault
}

// seriesGroup holds posts grouped by series name.
type seriesGroup struct {
	name  string          // raw series name from frontmatter
	slug  string          // slugified series name
	posts []*models.Post  // posts in this series
	cfg   *seriesOverride // per-series config override (may be nil)
}

// seriesConfig holds the parsed series configuration from config.Extra.
type seriesConfig struct {
	SlugPrefix  string                     // URL prefix (default: "series")
	AutoSidebar bool                       // auto-enable sidebar (default: true)
	Defaults    seriesDefaults             // defaults for all series feeds
	Overrides   map[string]*seriesOverride // per-series overrides keyed by series name
}

// seriesDefaults holds default settings for series feeds.
type seriesDefaults struct {
	ItemsPerPage int                 // default: 0 (no pagination)
	Sidebar      bool                // default: true
	Formats      *models.FeedFormats // if set, overrides feed defaults
}

// seriesOverride holds per-series configuration.
type seriesOverride struct {
	Title        string
	Description  string
	ItemsPerPage *int
	Formats      *models.FeedFormats
}

// Collect scans posts for series frontmatter and injects series FeedConfigs.
func (p *SeriesPlugin) Collect(m *lifecycle.Manager) error {
	posts := m.Posts()
	config := m.Config()

	// Parse series config
	seriesCfg := parseSeriesConfig(config)

	// Group posts by series
	groups := p.groupPostsBySeries(posts, seriesCfg)

	if len(groups) == 0 {
		return nil
	}

	// Get existing feed configs
	feedConfigs := getFeedConfigs(config)

	for _, group := range groups {
		if len(group.posts) == 0 {
			continue
		}

		// Sort posts within the series
		sortSeriesPosts(group, true)

		// Build feed slug
		feedSlug := buildSeriesFeedSlug(seriesCfg.SlugPrefix, group.slug)

		// Build ordered list of published posts for output/navigation
		publishedPosts := filterSeriesOutputPosts(group.posts)
		setSeriesMetadata(group.posts, feedSlug, len(publishedPosts))

		if len(publishedPosts) == 0 {
			log.Printf("[series] warning: no published posts in series %q", group.name)
			continue
		}

		// Set guide navigation (Prev/Next) on published posts
		setGuideNavigation(publishedPosts, feedSlug)

		// Set PrevNextContext with position info
		seriesTitle := p.seriesTitle(group)
		p.setPrevNextContext(publishedPosts, feedSlug, seriesTitle)

		// Build the FeedConfig
		fc := p.buildFeedConfig(group, feedSlug, seriesCfg, publishedPosts)

		feedConfigs = append(feedConfigs, fc)
	}

	// Store updated feed configs back
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}
	config.Extra["feeds"] = feedConfigs

	return nil
}

// groupPostsBySeries scans posts and groups them by series name.
func (p *SeriesPlugin) groupPostsBySeries(posts []*models.Post, cfg seriesConfig) []*seriesGroup {
	groupMap := make(map[string]*seriesGroup)
	var groupOrder []string // preserve discovery order for determinism

	for _, post := range posts {
		seriesName := getStringFromExtra(post.Extra, seriesKey)
		if seriesName == "" {
			continue
		}

		slug := slugify(seriesName)

		group, ok := groupMap[slug]
		if !ok {
			group = &seriesGroup{
				name: seriesName,
				slug: slug,
			}
			group.cfg = resolveSeriesOverride(cfg, seriesName, slug)
			groupMap[slug] = group
			groupOrder = append(groupOrder, slug)
		}

		group.posts = append(group.posts, post)
	}

	// Build result in discovery order
	groups := make([]*seriesGroup, 0, len(groupOrder))
	for _, slug := range groupOrder {
		groups = append(groups, groupMap[slug])
	}

	return groups
}

func resolveSeriesOverride(cfg seriesConfig, seriesName, seriesSlug string) *seriesOverride {
	if override, exists := cfg.Overrides[seriesSlug]; exists {
		return override
	}
	if override, exists := cfg.Overrides[seriesName]; exists {
		return override
	}
	return nil
}

// sortSeriesPosts sorts posts within a series according to ordering rules:
//  1. If any post has series_order, sort by series_order ascending
//  2. Posts without series_order are placed after ordered posts, sorted by date
//  3. If no post has series_order, sort by date ascending
//  4. Ties broken by file path
func sortSeriesPosts(group *seriesGroup, logDuplicates bool) {
	// Check if any post has series_order
	hasExplicitOrder := false
	for _, post := range group.posts {
		if _, ok := getSeriesOrder(post); ok {
			hasExplicitOrder = true
			break
		}
	}

	if hasExplicitOrder && logDuplicates {
		// Check for duplicate series_order values
		orderSeen := make(map[int]string) // order -> first post path
		for _, post := range group.posts {
			if order, ok := getSeriesOrder(post); ok {
				if prevPath, exists := orderSeen[order]; exists {
					log.Printf("[series] warning: duplicate series_order %d in series %q: %q and %q",
						order, group.name, prevPath, post.Path)
				} else {
					orderSeen[order] = post.Path
				}
			}
		}

		// Sort: posts with order first (ascending), then posts without order (by date ascending)
		sort.SliceStable(group.posts, func(i, j int) bool {
			orderI, hasI := getSeriesOrder(group.posts[i])
			orderJ, hasJ := getSeriesOrder(group.posts[j])

			// Both have order: compare orders
			if hasI && hasJ {
				if orderI != orderJ {
					return orderI < orderJ
				}
				// Tie-break by date then path
				return tieBreakByDateThenPath(group.posts[i], group.posts[j])
			}

			// Only i has order: i comes first
			if hasI && !hasJ {
				return true
			}

			// Only j has order: j comes first
			if !hasI && hasJ {
				return false
			}

			// Neither has order: sort by date ascending then path
			return tieBreakByDateThenPath(group.posts[i], group.posts[j])
		})
	} else {
		// No explicit order: sort by date ascending (oldest first for sequential reading)
		sort.SliceStable(group.posts, func(i, j int) bool {
			return tieBreakByDateThenPath(group.posts[i], group.posts[j])
		})
	}
}

// tieBreakByDateThenPath compares posts by date ascending, then path ascending.
func tieBreakByDateThenPath(a, b *models.Post) bool {
	dateA := a.Date
	dateB := b.Date

	switch {
	case dateA != nil && dateB != nil:
		if !dateA.Equal(*dateB) {
			return dateA.Before(*dateB)
		}
	case dateA != nil:
		return true // post with date comes first
	case dateB != nil:
		return false
	}

	return a.Path < b.Path
}

func filterSeriesOutputPosts(posts []*models.Post) []*models.Post {
	if len(posts) == 0 {
		return nil
	}
	result := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if !post.Published || post.Draft || post.Skip || post.Private {
			continue
		}
		result = append(result, post)
	}
	return result
}

func setSeriesMetadata(posts []*models.Post, feedSlug string, total int) {
	for _, post := range posts {
		if post.Extra == nil {
			post.Extra = make(map[string]interface{})
		}
		post.Extra["series_slug"] = feedSlug
		post.Extra["series_total"] = total
	}
}

// getSeriesOrder extracts the series_order from a post's Extra map.
func getSeriesOrder(post *models.Post) (int, bool) {
	if post.Extra == nil {
		return 0, false
	}
	v, ok := post.Extra["series_order"]
	if !ok {
		return 0, false
	}
	return parseIntFromInterface(v)
}

// setPrevNextContext sets PrevNextContext on each post in a series group.
func (p *SeriesPlugin) setPrevNextContext(posts []*models.Post, feedSlug, title string) {
	total := len(posts)

	for i, post := range posts {
		post.PrevNextContext = &models.PrevNextContext{
			FeedSlug:  feedSlug,
			FeedTitle: title,
			Position:  i + 1,
			Total:     total,
			Prev:      post.Prev,
			Next:      post.Next,
		}
	}
}

// seriesTitle returns the display title for a series.
// Uses override title if set, otherwise derives from the series name.
func (p *SeriesPlugin) seriesTitle(group *seriesGroup) string {
	return seriesDisplayTitle(group.name, group.cfg)
}

func seriesDisplayTitle(seriesName string, override *seriesOverride) string {
	if override != nil && override.Title != "" {
		return override.Title
	}

	// Derive title from series name: replace hyphens with spaces and title-case
	title := strings.ReplaceAll(seriesName, "-", " ")
	return toTitleCase(title)
}

// buildFeedConfig creates a FeedConfig for a series group.
func (p *SeriesPlugin) buildFeedConfig(group *seriesGroup, feedSlug string, cfg seriesConfig, posts []*models.Post) models.FeedConfig {
	title := p.seriesTitle(group)

	fc := models.FeedConfig{
		Slug:    feedSlug,
		Title:   title,
		Type:    models.FeedTypeSeries,
		Filter:  fmt.Sprintf("series_slug == %q and published == true and draft == false and skip == false", feedSlug),
		Sort:    "series_order",
		Reverse: false, // ascending order
		Sidebar: func() *bool { v := cfg.Defaults.Sidebar; return &v }(),
		Posts:   posts,
	}

	// Apply description from override
	if group.cfg != nil && group.cfg.Description != "" {
		fc.Description = group.cfg.Description
	}

	// Apply items_per_page: series default is 0 (no pagination)
	fc.ItemsPerPage = cfg.Defaults.ItemsPerPage
	if group.cfg != nil && group.cfg.ItemsPerPage != nil {
		fc.ItemsPerPage = *group.cfg.ItemsPerPage
	}

	// Apply formats
	if group.cfg != nil && group.cfg.Formats != nil {
		fc.Formats = *group.cfg.Formats
	} else if cfg.Defaults.Formats != nil {
		fc.Formats = *cfg.Defaults.Formats
	}
	// If no formats set, ApplyDefaults in the feeds plugin will apply feed defaults

	return fc
}

// parseSeriesConfig parses series configuration from config.Extra.
func parseSeriesConfig(config *lifecycle.Config) seriesConfig {
	cfg := seriesConfig{
		SlugPrefix:  seriesKey,
		AutoSidebar: true,
		Defaults: seriesDefaults{
			ItemsPerPage: 0,    // no pagination by default for series
			Sidebar:      true, // always show sidebar
		},
		Overrides: make(map[string]*seriesOverride),
	}

	if config.Extra == nil {
		return cfg
	}

	seriesRaw, ok := config.Extra[seriesKey]
	if !ok {
		return cfg
	}

	seriesMap, ok := seriesRaw.(map[string]interface{})
	if !ok {
		return cfg
	}

	// Parse slug_prefix
	if v, ok := seriesMap["slug_prefix"].(string); ok && v != "" {
		cfg.SlugPrefix = v
	}

	// Parse auto_sidebar
	if v, ok := seriesMap["auto_sidebar"].(bool); ok {
		cfg.AutoSidebar = v
	}

	// Parse defaults
	parseSeriesDefaults(seriesMap, &cfg.Defaults)

	// Parse overrides
	parseSeriesOverrides(seriesMap, cfg.Overrides)

	return cfg
}

// parseSeriesDefaults parses the defaults section from the series config map.
func parseSeriesDefaults(seriesMap map[string]interface{}, defaults *seriesDefaults) {
	defaultsRaw, ok := seriesMap["defaults"]
	if !ok {
		return
	}

	defaultsMap, ok := defaultsRaw.(map[string]interface{})
	if !ok {
		return
	}

	if v, ok := defaultsMap["items_per_page"]; ok {
		if n, ok := parseIntFromInterface(v); ok {
			defaults.ItemsPerPage = n
		}
	}
	if v, ok := defaultsMap["sidebar"].(bool); ok {
		defaults.Sidebar = v
	}
	if formatsRaw, ok := defaultsMap["formats"]; ok {
		if formats := parseFeedFormatsFromMap(formatsRaw); formats != nil {
			defaults.Formats = formats
		}
	}
}

// parseSeriesOverrides parses the overrides section from the series config map.
func parseSeriesOverrides(seriesMap map[string]interface{}, overrides map[string]*seriesOverride) {
	overridesRaw, ok := seriesMap["overrides"]
	if !ok {
		return
	}

	overridesMap, ok := overridesRaw.(map[string]interface{})
	if !ok {
		return
	}

	for name, overrideRaw := range overridesMap {
		overrideMap, ok := overrideRaw.(map[string]interface{})
		if !ok {
			continue
		}

		override := &seriesOverride{}
		if v, ok := overrideMap["title"].(string); ok {
			override.Title = v
		}
		if v, ok := overrideMap["description"].(string); ok {
			override.Description = v
		}
		if v, ok := overrideMap["items_per_page"]; ok {
			if n, ok := parseIntFromInterface(v); ok {
				override.ItemsPerPage = &n
			}
		}
		if formatsRaw, ok := overrideMap["formats"]; ok {
			override.Formats = parseFeedFormatsFromMap(formatsRaw)
		}
		overrides[name] = override
	}
}

// parseFeedFormatsFromMap parses a FeedFormats from a generic map.
func parseFeedFormatsFromMap(raw interface{}) *models.FeedFormats {
	m, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}

	formats := &models.FeedFormats{}
	if v, ok := m["html"].(bool); ok {
		formats.HTML = v
	}
	if v, ok := m["simple_html"].(bool); ok {
		formats.SimpleHTML = v
	}
	if v, ok := m["rss"].(bool); ok {
		formats.RSS = v
	}
	if v, ok := m["atom"].(bool); ok {
		formats.Atom = v
	}
	if v, ok := m["json"].(bool); ok {
		formats.JSON = v
	}
	if v, ok := m["markdown"].(bool); ok {
		formats.Markdown = v
	}
	if v, ok := m["text"].(bool); ok {
		formats.Text = v
	}
	if v, ok := m["sitemap"].(bool); ok {
		formats.Sitemap = v
	}

	return formats
}

// Ensure SeriesPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*SeriesPlugin)(nil)
	_ lifecycle.CollectPlugin  = (*SeriesPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*SeriesPlugin)(nil)
)
