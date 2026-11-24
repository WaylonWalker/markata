// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
)

// AutoFeedsConfig configures automatic feed generation.
type AutoFeedsConfig struct {
	// Tags configures automatic tag feeds
	Tags AutoFeedTypeConfig `json:"tags" yaml:"tags" toml:"tags"`

	// Categories configures automatic category feeds
	Categories AutoFeedTypeConfig `json:"categories" yaml:"categories" toml:"categories"`

	// Archives configures automatic date archive feeds
	Archives AutoArchiveConfig `json:"archives" yaml:"archives" toml:"archives"`
}

// AutoFeedTypeConfig configures a type of auto-generated feed (tags, categories).
type AutoFeedTypeConfig struct {
	// Enabled enables generation of this feed type
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// SlugPrefix is the URL prefix for feeds (e.g., "tags" -> /tags/python/)
	SlugPrefix string `json:"slug_prefix" yaml:"slug_prefix" toml:"slug_prefix"`

	// Formats specifies which output formats to generate
	Formats models.FeedFormats `json:"formats" yaml:"formats" toml:"formats"`

	// Robots controls the robots meta tag for generated HTML feed pages.
	Robots string `json:"robots,omitempty" yaml:"robots,omitempty" toml:"robots,omitempty"`
}

// AutoArchiveConfig configures automatic date archive feeds.
type AutoArchiveConfig struct {
	// Enabled enables generation of archive feeds
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// SlugPrefix is the URL prefix for archives (e.g., "archive" -> /archive/2024/)
	SlugPrefix string `json:"slug_prefix" yaml:"slug_prefix" toml:"slug_prefix"`

	// YearlyFeeds enables year-based archive feeds
	YearlyFeeds bool `json:"yearly_feeds" yaml:"yearly_feeds" toml:"yearly_feeds"`

	// MonthlyFeeds enables month-based archive feeds
	MonthlyFeeds bool `json:"monthly_feeds" yaml:"monthly_feeds" toml:"monthly_feeds"`

	// Formats specifies which output formats to generate
	Formats models.FeedFormats `json:"formats" yaml:"formats" toml:"formats"`

	// Robots controls the robots meta tag for generated HTML archive pages.
	Robots string `json:"robots,omitempty" yaml:"robots,omitempty" toml:"robots,omitempty"`
}

// Default slug prefix constants for auto-generated feeds.
const (
	autoFeedsPluginName     = "auto_feeds"
	defaultTagsPrefix       = "tags"
	defaultCategoriesPrefix = "categories"
	defaultArchivePrefix    = "archive"
)

// AutoFeedsPlugin automatically generates feeds for tags, categories, and date archives.
// It implements ConfigurePlugin to pre-register synthetic posts for wikilink resolution,
// and CollectPlugin to generate the actual feed content.
type AutoFeedsPlugin struct{}

type autoStringGroup struct {
	SlugPart string
	Display  string
	Variants []string
}

type autoTagGroup struct {
	autoStringGroup
}

// NewAutoFeedsPlugin creates a new AutoFeedsPlugin.
func NewAutoFeedsPlugin() *AutoFeedsPlugin {
	return &AutoFeedsPlugin{}
}

// Name returns the unique name of the plugin.
func (p *AutoFeedsPlugin) Name() string {
	return autoFeedsPluginName
}

// Priority returns the plugin's priority for a given stage.
func (p *AutoFeedsPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageLoad:
		// Run late in Load to ensure all posts are loaded
		// This allows us to scan for tags/categories before Transform stage
		return lifecycle.PriorityLate
	case lifecycle.StageCollect:
		// Run early in Collect to generate feeds before other plugins need them
		return lifecycle.PriorityDefault
	default:
		return lifecycle.PriorityDefault
	}
}

// Load pre-registers synthetic posts for auto-generated feeds
// so they can be resolved by wikilinks during the Transform stage.
// This runs after posts are loaded but before Transform stage.
func (p *AutoFeedsPlugin) Load(m *lifecycle.Manager) error {
	posts := m.Posts()
	config := m.Config()

	autoConfig := getAutoFeedsConfig(config)

	// Pre-register tag feed synthetic posts
	if autoConfig.Tags.Enabled {
		p.registerTagSyntheticPosts(m, posts, autoConfig.Tags)
	}

	// Pre-register category feed synthetic posts
	if autoConfig.Categories.Enabled {
		p.registerCategorySyntheticPosts(m, posts, autoConfig.Categories)
	}

	// Pre-register archive feed synthetic posts
	if autoConfig.Archives.Enabled {
		p.registerArchiveSyntheticPosts(m, posts, autoConfig.Archives)
	}

	return nil
}

// autoFeedsStrPtr returns a pointer to the given string.
func autoFeedsStrPtr(s string) *string { return &s }

func autoFeedSlugPrefix(prefix, fallback string) string {
	if prefix == "" {
		return fallback
	}
	return prefix
}

func autoFeedAliases(variants []string, slugPart string) []interface{} {
	aliases := make([]interface{}, 0, len(variants)+1)
	for _, variant := range variants {
		aliases = append(aliases, variant)
	}
	return append(aliases, slugPart)
}

func addSyntheticGroupedPost(m *lifecycle.Manager, prefix, slugPart, title, description string, variants []string) {
	slug := prefix + "/" + slugPart
	m.AddPost(&models.Post{
		Slug:        slug,
		Title:       autoFeedsStrPtr(title),
		Description: autoFeedsStrPtr(description),
		Href:        "/" + slug + "/",
		Published:   true,
		Skip:        true,
		Extra: map[string]interface{}{
			"aliases": autoFeedAliases(variants, slugPart),
		},
	})
}

func buildAutoFeedConfigs(groups []autoStringGroup, prefix, titleFormat, descriptionFormat string, filterFn func([]string) string, formats models.FeedFormats, robots string, changedSlugs map[string]bool) []models.FeedConfig {
	feeds := make([]models.FeedConfig, 0, len(groups))
	for _, group := range groups {
		if changedSlugs != nil && !changedSlugs[group.SlugPart] {
			continue
		}
		slug := prefix + "/" + group.SlugPart
		feeds = append(feeds, models.FeedConfig{
			Slug:        slug,
			Title:       fmt.Sprintf(titleFormat, group.Display),
			Description: fmt.Sprintf(descriptionFormat, group.Display),
			Filter:      filterFn(group.Variants),
			Sort:        "date",
			Reverse:     true,
			Formats:     formats,
			Robots:      robots,
		})
	}
	return feeds
}

// registerTagSyntheticPosts creates synthetic posts for tag feeds.
func (p *AutoFeedsPlugin) registerTagSyntheticPosts(m *lifecycle.Manager, posts []*models.Post, config AutoFeedTypeConfig) {
	prefix := autoFeedSlugPrefix(config.SlugPrefix, defaultTagsPrefix)
	groups := collectAutoTagGroups(posts)

	for _, group := range groups {
		addSyntheticGroupedPost(
			m,
			prefix,
			group.SlugPart,
			fmt.Sprintf("Posts tagged: %s", group.Display),
			fmt.Sprintf("All posts with the tag %q", group.Display),
			group.Variants,
		)
	}
}

// registerCategorySyntheticPosts creates synthetic posts for category feeds.
func (p *AutoFeedsPlugin) registerCategorySyntheticPosts(m *lifecycle.Manager, posts []*models.Post, config AutoFeedTypeConfig) {
	prefix := autoFeedSlugPrefix(config.SlugPrefix, defaultCategoriesPrefix)
	groups := collectAutoCategoryGroups(posts)

	for _, group := range groups {
		addSyntheticGroupedPost(
			m,
			prefix,
			group.SlugPart,
			fmt.Sprintf("Category: %s", group.Display),
			fmt.Sprintf("All posts in the %q category", group.Display),
			group.Variants,
		)
	}
}

// registerArchiveSyntheticPosts creates synthetic posts for archive feeds.
func (p *AutoFeedsPlugin) registerArchiveSyntheticPosts(m *lifecycle.Manager, posts []*models.Post, config AutoArchiveConfig) {
	prefix := config.SlugPrefix
	if prefix == "" {
		prefix = defaultArchivePrefix
	}

	// Collect all unique year/month combinations
	yearsMap := make(map[int]bool)
	yearMonthsMap := make(map[string]bool)

	for _, post := range posts {
		if post.Date != nil {
			year := post.Date.Year()
			month := post.Date.Month()
			yearsMap[year] = true
			yearMonthsMap[fmt.Sprintf("%04d/%02d", year, month)] = true
		}
	}

	// Register synthetic posts for yearly archives
	if config.YearlyFeeds {
		p.registerYearlyArchivePosts(m, yearsMap, prefix)
	}

	// Register synthetic posts for monthly archives
	if config.MonthlyFeeds {
		p.registerMonthlyArchivePosts(m, yearMonthsMap, prefix)
	}
}

// registerYearlyArchivePosts creates synthetic posts for yearly archive feeds.
func (p *AutoFeedsPlugin) registerYearlyArchivePosts(m *lifecycle.Manager, yearsMap map[int]bool, prefix string) {
	years := make([]int, 0, len(yearsMap))
	for year := range yearsMap {
		years = append(years, year)
	}
	sort.Ints(years)

	for _, year := range years {
		slug := fmt.Sprintf("%s/%04d", prefix, year)
		syntheticPost := &models.Post{
			Slug:        slug,
			Title:       autoFeedsStrPtr(fmt.Sprintf("Archive: %d", year)),
			Description: autoFeedsStrPtr(fmt.Sprintf("All posts from %d", year)),
			Href:        "/" + slug + "/",
			Published:   true,
			Skip:        true,
		}
		m.AddPost(syntheticPost)
	}
}

// registerMonthlyArchivePosts creates synthetic posts for monthly archive feeds.
func (p *AutoFeedsPlugin) registerMonthlyArchivePosts(m *lifecycle.Manager, yearMonthsMap map[string]bool, prefix string) {
	yearMonths := make([]string, 0, len(yearMonthsMap))
	for ym := range yearMonthsMap {
		yearMonths = append(yearMonths, ym)
	}
	sort.Strings(yearMonths)

	for _, ym := range yearMonths {
		var year, month int
		//nolint:errcheck // best-effort parsing
		fmt.Sscanf(ym, "%d/%d", &year, &month)

		slug := fmt.Sprintf("%s/%s", prefix, ym)
		monthName := time.Month(month).String()
		syntheticPost := &models.Post{
			Slug:        slug,
			Title:       autoFeedsStrPtr(fmt.Sprintf("Archive: %s %d", monthName, year)),
			Description: autoFeedsStrPtr(fmt.Sprintf("All posts from %s %d", monthName, year)),
			Href:        "/" + slug + "/",
			Published:   true,
			Skip:        true,
		}
		m.AddPost(syntheticPost)
	}
}

// Collect generates automatic feeds for tags, categories, and date archives.
func (p *AutoFeedsPlugin) Collect(m *lifecycle.Manager) error {
	var generation, filtering, recording, preparation time.Duration
	var feedCount, matchedCount, selectedCount, observationCount int
	defer func() {
		logging.Component(p.Name()).Phase("collect").Level("debug").Printf(
			"generation_ns=%d filtering_sorting_ns=%d selection_recording_ns=%d pagination_preparation_ns=%d feeds=%d matched=%d selected=%d observations=%d",
			generation.Nanoseconds(), filtering.Nanoseconds(), recording.Nanoseconds(), preparation.Nanoseconds(),
			feedCount, matchedCount, selectedCount, observationCount)
	}()
	start := time.Now()
	posts := m.Posts()
	config := m.Config()

	autoConfig := getAutoFeedsConfig(config)

	// Collect only auto-generated feed configs
	var autoFeedConfigs []models.FeedConfig

	privateTagSlugs := getPrivateTagSlugs(config)

	// Auto-feed config generation is cheap, so always rebuild the full set.
	// Partial regeneration drops unchanged auto feeds from the current build.
	if autoConfig.Tags.Enabled {
		tagFeeds := p.generateTagFeeds(posts, autoConfig.Tags, privateTagSlugs)
		autoFeedConfigs = append(autoFeedConfigs, tagFeeds...)
	}

	if autoConfig.Categories.Enabled {
		categoryFeeds := p.generateCategoryFeeds(posts, autoConfig.Categories)
		autoFeedConfigs = append(autoFeedConfigs, categoryFeeds...)
	}

	if autoConfig.Archives.Enabled {
		archiveFeeds := p.generateArchiveFeeds(posts, autoConfig.Archives)
		autoFeedConfigs = append(autoFeedConfigs, archiveFeeds...)
	}
	generation = time.Since(start)
	feedCount = len(autoFeedConfigs)

	// If no auto-feeds were generated, nothing to do
	if len(autoFeedConfigs) == 0 {
		return nil
	}

	start = time.Now()
	filterCache := newFeedFilterCache(posts)
	selection := newFeedSelectionRecorder(m.ContentLedger(), posts)
	// Process auto-generated feeds
	feedDefaults := getFeedDefaults(config)

	// Get existing feeds from FeedsPlugin
	feeds := m.Feeds()

	// Get existing feed configs from cache to append auto-generated ones
	allFeedConfigs := make([]models.FeedConfig, 0, len(autoFeedConfigs))
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			allFeedConfigs = append(allFeedConfigs, fcs...)
		}
	}
	preparation += time.Since(start)

	for i := range autoFeedConfigs {
		start = time.Now()
		fc := &autoFeedConfigs[i]

		// Apply defaults
		fc.ApplyDefaults(feedDefaults)
		preparation += time.Since(start)
		start = time.Now()

		// Filter posts for this feed
		filteredPosts, err := filterCache.FilterPosts(fc.Filter, fc.IncludesPrivate())
		if err != nil {
			filtering += time.Since(start)
			return fmt.Errorf("auto feed %q: %w", fc.Slug, err)
		}
		filteredPosts = cloneFeedPosts(filteredPosts)

		// Sort posts by date, newest first
		sortPosts(filteredPosts, "date", true)
		filtering += time.Since(start)
		start = time.Now()
		matchedPosts := filteredPosts
		filteredPosts = applyFeedLimitOffset(filteredPosts, fc)

		// Store posts in feed config
		fc.Posts = filteredPosts
		preparation += time.Since(start)
		start = time.Now()
		selection.record(fc.Slug, selection.sources, matchedPosts, filteredPosts, fc.Filter, fc.IncludesPrivate(), fc.Offset, fc.Limit)
		recording += time.Since(start)
		matchedCount += len(matchedPosts)
		selectedCount += len(filteredPosts)
		observationCount += len(selection.sources)
		start = time.Now()

		// Get base URL for pagination
		baseURL := "/" + fc.Slug

		// Paginate results
		fc.Paginate(baseURL)

		// Create lifecycle.Feed
		feed := &lifecycle.Feed{
			Name:           fc.Slug,
			Title:          fc.Title,
			Posts:          filteredPosts,
			Path:           fc.Slug,
			IncludePrivate: fc.IncludesPrivate(),
			Automated:      true,
		}

		feeds = append(feeds, feed)
		allFeedConfigs = append(allFeedConfigs, *fc)
		preparation += time.Since(start)
	}

	start = time.Now()
	m.SetFeeds(feeds)

	// Update cache with all feed configs (original + auto-generated)
	m.Cache().Set("feed_configs", allFeedConfigs)
	preparation += time.Since(start)

	return nil
}

// generateTagFeeds creates feed configurations for each unique tag.
func (p *AutoFeedsPlugin) generateTagFeeds(posts []*models.Post, config AutoFeedTypeConfig, privateTagSlugs map[string]bool) []models.FeedConfig {
	groups := collectAutoTagGroups(posts)
	prefix := autoFeedSlugPrefix(config.SlugPrefix, defaultTagsPrefix)
	stringGroups := make([]autoStringGroup, 0, len(groups))
	for _, group := range groups {
		stringGroups = append(stringGroups, group.autoStringGroup)
	}
	feeds := buildAutoFeedConfigs(stringGroups, prefix, "Posts tagged: %s", "All posts with the tag %q", buildTagFilterExpression, config.Formats, config.Robots, nil)
	for i := range feeds {
		slug := strings.TrimPrefix(feeds[i].Slug, prefix+"/")
		if privateTagSlugs[slug] {
			feeds[i].IncludePrivate = true
		}
	}
	return feeds
}

func (p *AutoFeedsPlugin) generateTagFeedsForChanged(posts []*models.Post, config AutoFeedTypeConfig, changedSet, privateTagSlugs map[string]bool) []models.FeedConfig {
	if len(changedSet) == 0 {
		return nil
	}
	changedSlugs := collectChangedTagSlugs(posts, changedSet)
	if len(changedSlugs) == 0 {
		return nil
	}
	groups := collectAutoTagGroups(posts)
	prefix := autoFeedSlugPrefix(config.SlugPrefix, defaultTagsPrefix)
	stringGroups := make([]autoStringGroup, 0, len(groups))
	for _, group := range groups {
		stringGroups = append(stringGroups, group.autoStringGroup)
	}
	feeds := buildAutoFeedConfigs(stringGroups, prefix, "Posts tagged: %s", "All posts with the tag %q", buildTagFilterExpression, config.Formats, config.Robots, changedSlugs)
	for i := range feeds {
		slug := strings.TrimPrefix(feeds[i].Slug, prefix+"/")
		if privateTagSlugs[slug] {
			feeds[i].IncludePrivate = true
		}
	}
	return feeds
}

// generateCategoryFeeds creates feed configurations for each unique category.
func (p *AutoFeedsPlugin) generateCategoryFeeds(posts []*models.Post, config AutoFeedTypeConfig) []models.FeedConfig {
	groups := collectAutoCategoryGroups(posts)
	prefix := autoFeedSlugPrefix(config.SlugPrefix, defaultCategoriesPrefix)
	return buildAutoFeedConfigs(groups, prefix, "Category: %s", "All posts in the %q category", buildCategoryFilterExpression, config.Formats, config.Robots, nil)
}

// generateArchiveFeeds creates feed configurations for year and month archives.
func (p *AutoFeedsPlugin) generateArchiveFeeds(posts []*models.Post, config AutoArchiveConfig) []models.FeedConfig {
	// Collect all unique year/month combinations
	yearMonths := make(map[string]bool)
	years := make(map[int]bool)

	for _, post := range posts {
		if post.Date != nil {
			year := post.Date.Year()
			month := post.Date.Month()

			years[year] = true
			yearMonths[fmt.Sprintf("%04d/%02d", year, month)] = true
		}
	}

	var feeds []models.FeedConfig
	prefix := config.SlugPrefix
	if prefix == "" {
		prefix = defaultArchivePrefix
	}

	// Create yearly feeds
	if config.YearlyFeeds {
		var yearList []int
		for year := range years {
			yearList = append(yearList, year)
		}
		sort.Sort(sort.Reverse(sort.IntSlice(yearList)))

		for _, year := range yearList {
			slug := fmt.Sprintf("%s/%04d", prefix, year)
			startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
			endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

			feeds = append(feeds, models.FeedConfig{
				Slug:        slug,
				Title:       fmt.Sprintf("Archive: %d", year),
				Description: fmt.Sprintf("All posts from %d", year),
				Filter:      fmt.Sprintf("date >= %q and date < %q", startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)),
				Sort:        "date",
				Reverse:     true,
				Formats:     config.Formats,
			})
		}
	}

	// Create monthly feeds
	if config.MonthlyFeeds {
		var ymList []string
		for ym := range yearMonths {
			ymList = append(ymList, ym)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(ymList)))

		for _, ym := range ymList {
			var year, month int
			//nolint:errcheck // best-effort parsing, invalid format will result in zero values
			fmt.Sscanf(ym, "%d/%d", &year, &month)

			slug := fmt.Sprintf("%s/%s", prefix, ym)
			monthName := time.Month(month).String()
			startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endDate := startDate.AddDate(0, 1, 0)

			feeds = append(feeds, models.FeedConfig{
				Slug:        slug,
				Title:       fmt.Sprintf("Archive: %s %d", monthName, year),
				Description: fmt.Sprintf("All posts from %s %d", monthName, year),
				Filter:      fmt.Sprintf("date >= %q and date < %q", startDate.Format(time.RFC3339), endDate.Format(time.RFC3339)),
				Sort:        "date",
				Reverse:     true,
				Formats:     config.Formats,
			})
		}
	}

	return feeds
}

// getAutoFeedsConfig retrieves auto feeds configuration from the manager config.
func getAutoFeedsConfig(config *lifecycle.Config) AutoFeedsConfig {
	defaultConfig := AutoFeedsConfig{
		Tags: AutoFeedTypeConfig{
			Enabled:    true, // Tag feeds enabled by default
			SlugPrefix: defaultTagsPrefix,
			Formats: models.FeedFormats{
				HTML:     true,
				RSS:      true,
				Atom:     true,
				Markdown: true,
				Text:     true,
			},
		},
		Categories: AutoFeedTypeConfig{
			Enabled:    false,
			SlugPrefix: defaultCategoriesPrefix,
			Formats: models.FeedFormats{
				HTML:     true,
				RSS:      true,
				Atom:     true,
				Markdown: true,
				Text:     true,
			},
		},
		Archives: AutoArchiveConfig{
			Enabled:      false,
			SlugPrefix:   defaultArchivePrefix,
			YearlyFeeds:  true,
			MonthlyFeeds: false,
			Formats: models.FeedFormats{
				HTML:     true,
				RSS:      false,
				Atom:     false,
				Markdown: true,
				Text:     true,
			},
		},
	}

	if config.Extra == nil {
		return defaultConfig
	}

	if autoFeeds, ok := config.Extra[autoFeedsPluginName]; ok {
		if ac, ok := autoFeeds.(AutoFeedsConfig); ok {
			return ac
		}
		if raw, ok := autoFeeds.(map[string]any); ok {
			ac := defaultConfig
			if tagsRaw, ok := raw["tags"].(map[string]any); ok {
				applyAutoFeedTypeOverrides(&ac.Tags, tagsRaw)
			}
			if categoriesRaw, ok := raw["categories"].(map[string]any); ok {
				applyAutoFeedTypeOverrides(&ac.Categories, categoriesRaw)
			}
			if archivesRaw, ok := raw["archives"].(map[string]any); ok {
				applyAutoArchiveOverrides(&ac.Archives, archivesRaw)
			}
			return ac
		}
	}

	return defaultConfig
}

func applyAutoFeedTypeOverrides(cfg *AutoFeedTypeConfig, raw map[string]any) {
	if v, ok := raw["enabled"].(bool); ok {
		cfg.Enabled = v
	}
	if v, ok := raw["slug_prefix"].(string); ok && v != "" {
		cfg.SlugPrefix = v
	}
	if v, ok := raw["robots"].(string); ok {
		cfg.Robots = v
	}
	if formatsRaw, ok := raw["formats"].(map[string]any); ok {
		applyFeedFormatOverrides(&cfg.Formats, formatsRaw)
	}
}

func applyAutoArchiveOverrides(cfg *AutoArchiveConfig, raw map[string]any) {
	if v, ok := raw["enabled"].(bool); ok {
		cfg.Enabled = v
	}
	if v, ok := raw["slug_prefix"].(string); ok && v != "" {
		cfg.SlugPrefix = v
	}
	if v, ok := raw["yearly_feeds"].(bool); ok {
		cfg.YearlyFeeds = v
	}
	if v, ok := raw["monthly_feeds"].(bool); ok {
		cfg.MonthlyFeeds = v
	}
	if v, ok := raw["robots"].(string); ok {
		cfg.Robots = v
	}
	if formatsRaw, ok := raw["formats"].(map[string]any); ok {
		applyFeedFormatOverrides(&cfg.Formats, formatsRaw)
	}
}

func applyFeedFormatOverrides(cfg *models.FeedFormats, raw map[string]any) {
	if v, ok := raw["html"].(bool); ok {
		cfg.HTML = v
	}
	if v, ok := raw["simple_html"].(bool); ok {
		cfg.SimpleHTML = v
	}
	if v, ok := raw["rss"].(bool); ok {
		cfg.RSS = v
	}
	if v, ok := raw["atom"].(bool); ok {
		cfg.Atom = v
	}
	if v, ok := raw["json"].(bool); ok {
		cfg.JSON = v
	}
	if v, ok := raw["markdown"].(bool); ok {
		cfg.Markdown = v
	}
	if v, ok := raw["text"].(bool); ok {
		cfg.Text = v
	}
	if v, ok := raw["sitemap"].(bool); ok {
		cfg.Sitemap = v
	}
}

// slugify converts a string to a URL-safe slug.
// This is a convenience wrapper around models.Slugify for internal use.
func slugify(s string) string {
	return models.Slugify(s)
}

func collectAutoTagGroups(posts []*models.Post) []autoTagGroup {
	groups := make(map[string]*autoTagGroup)
	for _, post := range posts {
		for _, tag := range post.Tags {
			slug := slugify(tag)
			if slug == "" {
				continue
			}
			group, ok := groups[slug]
			if !ok {
				group = &autoTagGroup{autoStringGroup: autoStringGroup{SlugPart: slug, Display: tag}}
				groups[slug] = group
			}
			group.Display = pickPreferredAutoLabel(group.Display, tag)
			group.Variants = appendUniqueSorted(group.Variants, tag)
		}
	}

	result := make([]autoTagGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, *group)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SlugPart < result[j].SlugPart
	})
	return result
}

func collectAutoCategoryGroups(posts []*models.Post) []autoStringGroup {
	groups := make(map[string]*autoStringGroup)
	for _, post := range posts {
		cat, ok := post.Extra["category"].(string)
		if !ok || cat == "" {
			continue
		}
		slug := slugify(cat)
		if slug == "" {
			continue
		}
		group, ok := groups[slug]
		if !ok {
			group = &autoStringGroup{SlugPart: slug, Display: cat}
			groups[slug] = group
		}
		group.Display = pickPreferredAutoLabel(group.Display, cat)
		group.Variants = appendUniqueSorted(group.Variants, cat)
	}

	result := make([]autoStringGroup, 0, len(groups))
	for _, group := range groups {
		result = append(result, *group)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SlugPart < result[j].SlugPart
	})
	return result
}

func collectChangedTagSlugs(posts []*models.Post, changedSet map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for _, post := range posts {
		if !changedSet[post.Slug] {
			continue
		}
		for _, tag := range post.Tags {
			slug := slugify(tag)
			if slug != "" {
				result[slug] = true
			}
		}
	}
	return result
}

func pickPreferredAutoLabel(current, candidate string) string {
	if current == "" {
		return candidate
	}
	currentScore, currentKey := autoLabelPreference(current)
	candidateScore, candidateKey := autoLabelPreference(candidate)
	if candidateScore < currentScore {
		return candidate
	}
	if candidateScore > currentScore {
		return current
	}
	if candidateKey < currentKey {
		return candidate
	}
	return current
}

func autoLabelPreference(value string) (score int, key string) {
	trimmed := strings.TrimSpace(value)
	score = 1
	if strings.Contains(trimmed, " ") {
		score = 0
	} else if strings.Contains(trimmed, "-") {
		score = 2
	}
	key = strings.ToLower(trimmed)
	return score, key
}

func appendUniqueSorted(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	values = append(values, value)
	sort.Strings(values)
	return values
}

// getPrivateTagSlugs extracts the set of private tag slugs from the encryption config.
// These tag feeds should include private posts so encrypted content renders on the page.
func getPrivateTagSlugs(config *lifecycle.Config) map[string]bool {
	slugs := make(map[string]bool)
	if modelsConfig, ok := getModelsConfig(config); ok {
		for tag := range modelsConfig.Encryption.PrivateTags {
			slug := slugify(tag)
			if slug != "" {
				slugs[slug] = true
			}
		}
	}
	return slugs
}

func buildTagFilterExpression(variants []string) string {
	parts := make([]string, 0, len(variants))
	for _, variant := range variants {
		parts = append(parts, fmt.Sprintf("%q in tags", variant))
	}
	return strings.Join(parts, " or ")
}

func buildCategoryFilterExpression(variants []string) string {
	parts := make([]string, 0, len(variants))
	for _, variant := range variants {
		parts = append(parts, fmt.Sprintf("category == %q", variant))
	}
	return strings.Join(parts, " or ")
}

// Ensure AutoFeedsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*AutoFeedsPlugin)(nil)
	_ lifecycle.LoadPlugin     = (*AutoFeedsPlugin)(nil)
	_ lifecycle.CollectPlugin  = (*AutoFeedsPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*AutoFeedsPlugin)(nil)
)
