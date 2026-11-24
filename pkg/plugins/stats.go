// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// statsZeroMin is the default text for zero duration.
const (
	statsZeroMin    = "0 min"
	statsPluginName = "stats"
)

// StatsPlugin calculates comprehensive content statistics for posts and feeds.
// It provides word count, reading time, character count, and code block metrics
// for individual posts, and aggregates these statistics at the feed level.
type StatsPlugin struct {
	// wordsPerMinute is the average reading speed (default: 200)
	wordsPerMinute int
	// includeCodeInCount includes code block content in word count
	includeCodeInCount bool
	// trackCodeBlocks enables counting of code block lines
	trackCodeBlocks bool
}

// PostStats contains calculated statistics for a single post.
type PostStats struct {
	// WordCount is the number of words in the post
	WordCount int `json:"word_count"`
	// CharCount is the number of characters (excluding whitespace)
	CharCount int `json:"char_count"`
	// ReadingTime is the estimated reading time in minutes
	ReadingTime int `json:"reading_time"`
	// ReadingTimeText is a formatted reading time string
	ReadingTimeText string `json:"reading_time_text"`
	// CodeLines is the number of lines of code in code blocks
	CodeLines int `json:"code_lines"`
	// CodeBlocks is the number of code blocks in the post
	CodeBlocks int `json:"code_blocks"`
}

// FeedStats contains aggregated statistics for a feed.
type FeedStats struct {
	// PostCount is the total number of posts in the feed
	PostCount int `json:"post_count"`
	// TotalWords is the sum of word counts across all posts
	TotalWords int `json:"total_words"`
	// TotalChars is the sum of character counts across all posts
	TotalChars int `json:"total_chars"`
	// TotalReadingTime is the sum of reading times in minutes
	TotalReadingTime int `json:"total_reading_time"`
	// TotalReadingTimeText is a formatted total reading time
	TotalReadingTimeText string `json:"total_reading_time_text"`
	// AverageWords is the average word count per post
	AverageWords int `json:"average_words"`
	// AverageReadingTime is the average reading time per post
	AverageReadingTime int `json:"average_reading_time"`
	// AverageReadingTimeText is formatted average reading time
	AverageReadingTimeText string `json:"average_reading_time_text"`
	// TotalCodeLines is the sum of code lines across all posts
	TotalCodeLines int `json:"total_code_lines"`
	// TotalCodeBlocks is the sum of code blocks across all posts
	TotalCodeBlocks int `json:"total_code_blocks"`
	// PostsByYear maps year to post count for this feed
	PostsByYear map[int]int `json:"posts_by_year"`
	// WordsByYear maps year to total word count for this feed
	WordsByYear map[int]int `json:"words_by_year"`
	// PostsByTag maps tag name to post count for this feed
	PostsByTag map[string]int `json:"posts_by_tag"`
}

// SiteStats contains global statistics across all posts.
type SiteStats struct {
	// TotalPosts is the total number of posts on the site
	TotalPosts int `json:"total_posts"`
	// TotalWords is the sum of word counts across all posts
	TotalWords int `json:"total_words"`
	// TotalChars is the sum of character counts
	TotalChars int `json:"total_chars"`
	// TotalReadingTime is the total reading time in minutes
	TotalReadingTime int `json:"total_reading_time"`
	// TotalReadingTimeText is formatted total reading time
	TotalReadingTimeText string `json:"total_reading_time_text"`
	// AverageWords is the average word count per post
	AverageWords int `json:"average_words"`
	// AverageReadingTime is the average reading time per post
	AverageReadingTime int `json:"average_reading_time"`
	// AverageReadingTimeText is formatted average reading time
	AverageReadingTimeText string `json:"average_reading_time_text"`
	// TotalCodeLines is the total lines of code across the site
	TotalCodeLines int `json:"total_code_lines"`
	// TotalCodeBlocks is the total number of code blocks
	TotalCodeBlocks int `json:"total_code_blocks"`
	// PostsByYear maps year to post count
	PostsByYear map[int]int `json:"posts_by_year"`
	// WordsByYear maps year to total word count
	WordsByYear map[int]int `json:"words_by_year"`
	// PostsByTag maps tag name to post count
	PostsByTag map[string]int `json:"posts_by_tag"`
}

// NewStatsPlugin creates a new StatsPlugin with default settings.
func NewStatsPlugin() *StatsPlugin {
	return &StatsPlugin{
		wordsPerMinute:     defaultWordsPerMinute,
		includeCodeInCount: false,
		trackCodeBlocks:    true,
	}
}

// Name returns the unique name of the plugin.
func (p *StatsPlugin) Name() string {
	return statsPluginName
}

// Configure reads configuration options for the plugin.
func (p *StatsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Check for stats plugin config
	if statsConfig, ok := config.Extra[statsPluginName].(map[string]interface{}); ok {
		if wpm, ok := parseIntFromInterface(statsConfig["words_per_minute"]); ok && wpm > 0 {
			p.wordsPerMinute = wpm
		}
		if includeCode, ok := statsConfig["include_code_in_count"].(bool); ok {
			p.includeCodeInCount = includeCode
		}
		if trackCode, ok := statsConfig["track_code_blocks"].(bool); ok {
			p.trackCodeBlocks = trackCode
		}
	}

	// Also check top-level words_per_minute for backwards compatibility
	if wpm, ok := parseIntFromInterface(config.Extra["words_per_minute"]); ok && wpm > 0 {
		p.wordsPerMinute = wpm
	}

	return nil
}

// Transform calculates statistics for each post.
func (p *StatsPlugin) Transform(m *lifecycle.Manager) error {
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.Content != ""
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		stats := p.calculatePostStats(post.Content)

		// Stats calculates and stores its own values. ReadingTimePlugin may run
		// later and overwrite the public reading-time fields without changing the
		// Stats-specific PostStats value.
		post.Set("word_count", stats.WordCount)
		post.Set("char_count", stats.CharCount)
		post.Set("reading_time", stats.ReadingTime)
		post.Set("reading_time_text", stats.ReadingTimeText)
		post.Set("code_lines", stats.CodeLines)
		post.Set("code_blocks", stats.CodeBlocks)

		// Also store as a stats object
		post.Set(statsPluginName, stats)

		return nil
	})
}

// Collect aggregates statistics at the feed level.
func (p *StatsPlugin) Collect(m *lifecycle.Manager) error {
	feeds := m.Feeds()
	siteStats := &SiteStats{
		PostsByYear: make(map[int]int),
		WordsByYear: make(map[int]int),
		PostsByTag:  make(map[string]int),
	}

	// Track posts already counted for site stats to avoid double-counting
	countedPaths := make(map[string]bool)

	for _, feed := range feeds {
		feedStats := p.calculateFeedStatsFromLifecycle(feed)

		// Store in feed's Extra or runtime data
		// Note: FeedConfig doesn't have Extra, so we store in cache
		cacheKey := fmt.Sprintf("feed_stats.%s", feed.Name)
		m.Cache().Set(cacheKey, feedStats)

		// Aggregate to site stats (avoiding double counting)
		for _, post := range feed.Posts {
			if countedPaths[post.Path] {
				continue
			}
			countedPaths[post.Path] = true
			stats := p.getPostStats(post)
			siteStats.TotalPosts++
			siteStats.TotalWords += stats.WordCount
			siteStats.TotalChars += stats.CharCount
			siteStats.TotalReadingTime += stats.ReadingTime
			siteStats.TotalCodeLines += stats.CodeLines
			siteStats.TotalCodeBlocks += stats.CodeBlocks

			// Aggregate by year
			if post.Date != nil {
				year := post.Date.Year()
				siteStats.PostsByYear[year]++
				siteStats.WordsByYear[year] += stats.WordCount
			}

			// Aggregate by tag
			for _, tag := range post.Tags {
				siteStats.PostsByTag[tag]++
			}
		}
	}

	// Calculate averages for site
	if siteStats.TotalPosts > 0 {
		siteStats.AverageWords = siteStats.TotalWords / siteStats.TotalPosts
		siteStats.AverageReadingTime = siteStats.TotalReadingTime / siteStats.TotalPosts
	}
	siteStats.TotalReadingTimeText = p.formatDuration(siteStats.TotalReadingTime)
	siteStats.AverageReadingTimeText = p.formatReadingTime(siteStats.AverageReadingTime)

	// Store site stats in cache
	m.Cache().Set("site_stats", siteStats)

	// Also store in config extra for template access
	config := m.Config()
	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}
	config.Extra["site_stats"] = siteStats

	// Store stats helper object for template function access
	statsHelper := NewStatsHelper(m)
	m.Cache().Set("stats_helper", statsHelper)
	config.Extra[statsPluginName] = statsHelper

	return nil
}

// calculatePostStats computes all statistics for a post's content.
// This function is optimized for performance:
// - Uses the shared normalized metrics for words, characters, and duration
// - Uses fast string scanning for Stats-specific code metrics
// - Minimizes string allocations where possible
func (p *StatsPlugin) calculatePostStats(content string) *PostStats {
	metrics := calculateReadingTimeMetrics(content, p.wordsPerMinute, p.includeCodeInCount)
	stats := &PostStats{
		WordCount:       metrics.WordCount,
		CharCount:       metrics.CharCount,
		ReadingTime:     metrics.ReadingTime,
		ReadingTimeText: metrics.ReadingTimeText,
	}

	if !p.trackCodeBlocks {
		return stats
	}

	// Extract code blocks using fast string scanning instead of regex.
	codeBlocks := extractCodeBlocks(content)
	stats.CodeBlocks = len(codeBlocks)

	for _, block := range codeBlocks {
		// Count lines without creating intermediate slice
		lines := strings.Count(block, "\n") + 1
		// Subtract 2 for opening and closing fences
		codeLineCount := lines - 2
		if codeLineCount < 0 {
			codeLineCount = 0
		}
		stats.CodeLines += codeLineCount
	}

	return stats
}

// extractCodeBlocks finds all fenced code blocks (``` or ~~~). This is much
// faster than regex and is used for Stats-specific code metrics.
//
//nolint:gocyclo // Complexity is acceptable for this performance-critical string parsing
func extractCodeBlocks(s string) (blocks []string) {
	blocks = make([]string, 0, 4) // Pre-allocate for typical number of code blocks

	i := 0
	for i < len(s) {
		// Look for fence start (``` or ~~~)
		if i+2 < len(s) && (s[i] == '`' && s[i+1] == '`' && s[i+2] == '`' ||
			s[i] == '~' && s[i+1] == '~' && s[i+2] == '~') {
			fenceChar := s[i]
			fenceStart := i
			closed := false

			// Find end of opening fence line
			i += 3
			for i < len(s) && s[i] != '\n' {
				i++
			}
			if i < len(s) {
				i++ // Skip the newline
			}

			// Find closing fence
			blockStart := i
			for i < len(s) {
				// Look for closing fence at start of line
				if (i == blockStart || s[i-1] == '\n') &&
					i+2 < len(s) &&
					s[i] == fenceChar && s[i+1] == fenceChar && s[i+2] == fenceChar {
					// Found closing fence - capture up to end of fence (not trailing newline)
					fenceEnd := i + 3
					// Skip any additional fence chars or text on same line
					for fenceEnd < len(s) && s[fenceEnd] != '\n' {
						fenceEnd++
					}
					// Do NOT include the trailing newline in the block
					// This matches the old regex behavior

					blocks = append(blocks, s[fenceStart:fenceEnd])
					closed = true

					// Skip past the newline for next iteration
					if fenceEnd < len(s) && s[fenceEnd] == '\n' {
						fenceEnd++
					}
					i = fenceEnd
					break
				}
				i++
			}
			if !closed {
				// Unclosed fence contains no complete code block.
				break
			}
		} else {
			i++
		}
	}

	return blocks
}

// calculateFeedStats aggregates stats for a feed's posts (models.FeedConfig).
func (p *StatsPlugin) calculateFeedStats(feed *models.FeedConfig) *FeedStats {
	return p.aggregateFeedStats(feed.Posts)
}

// calculateFeedStatsFromLifecycle aggregates stats for a lifecycle.Feed.
func (p *StatsPlugin) calculateFeedStatsFromLifecycle(feed *lifecycle.Feed) *FeedStats {
	return p.aggregateFeedStats(feed.Posts)
}

// aggregateFeedStats calculates feed statistics from a slice of posts.
func (p *StatsPlugin) aggregateFeedStats(posts []*models.Post) *FeedStats {
	stats := &FeedStats{
		PostCount:   len(posts),
		PostsByYear: make(map[int]int),
		WordsByYear: make(map[int]int),
		PostsByTag:  make(map[string]int),
	}

	for _, post := range posts {
		postStats := p.getPostStats(post)
		stats.TotalWords += postStats.WordCount
		stats.TotalChars += postStats.CharCount
		stats.TotalReadingTime += postStats.ReadingTime
		stats.TotalCodeLines += postStats.CodeLines
		stats.TotalCodeBlocks += postStats.CodeBlocks

		// Aggregate by year
		if post.Date != nil {
			year := post.Date.Year()
			stats.PostsByYear[year]++
			stats.WordsByYear[year] += postStats.WordCount
		}

		// Aggregate by tag
		for _, tag := range post.Tags {
			stats.PostsByTag[tag]++
		}
	}

	// Calculate averages
	if stats.PostCount > 0 {
		stats.AverageWords = stats.TotalWords / stats.PostCount
		stats.AverageReadingTime = stats.TotalReadingTime / stats.PostCount
	}

	// Format text values
	stats.TotalReadingTimeText = p.formatDuration(stats.TotalReadingTime)
	stats.AverageReadingTimeText = p.formatReadingTime(stats.AverageReadingTime)

	return stats
}

// getPostStats retrieves stats from a post's Extra map.
func (p *StatsPlugin) getPostStats(post *models.Post) *PostStats {
	stats := &PostStats{}

	if post == nil || post.Extra == nil {
		return stats
	}

	// The structured Stats result is authoritative for Stats aggregates. The
	// public reading-time fields may be overwritten later by ReadingTimePlugin.
	if postStats, ok := post.Extra[statsPluginName].(*PostStats); ok && postStats != nil {
		return postStats
	}

	if wc, ok := post.Extra["word_count"].(int); ok {
		stats.WordCount = wc
	}
	if cc, ok := post.Extra["char_count"].(int); ok {
		stats.CharCount = cc
	}
	if rt, ok := post.Extra["reading_time"].(int); ok {
		stats.ReadingTime = rt
	}
	if cl, ok := post.Extra["code_lines"].(int); ok {
		stats.CodeLines = cl
	}
	if cb, ok := post.Extra["code_blocks"].(int); ok {
		stats.CodeBlocks = cb
	}

	return stats
}

// formatReadingTime creates a human-readable reading time string.
func (p *StatsPlugin) formatReadingTime(minutes int) string {
	return formatReadingTimeText(minutes)
}

// formatDuration formats a duration in minutes as hours and minutes.
func (p *StatsPlugin) formatDuration(minutes int) string {
	if minutes == 0 {
		return statsZeroMin
	}
	if minutes < 60 {
		return fmt.Sprintf("%d min", minutes)
	}

	hours := minutes / 60
	mins := minutes % 60

	if mins == 0 {
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}

	if hours == 1 {
		return fmt.Sprintf("1 hour %d min", mins)
	}
	return fmt.Sprintf("%d hours %d min", hours, mins)
}

// Priority returns the plugin priority for the given stage.
// Stats runs early in transform (after content is loaded) and late during
// collection (after feeds are populated).
func (p *StatsPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageTransform:
		return lifecycle.PriorityEarly // Calculate post stats early
	case lifecycle.StageCollect:
		return lifecycle.PriorityLate // Aggregate after feeds are built
	default:
		return lifecycle.PriorityDefault
	}
}

// GetFeedStats retrieves feed statistics from the cache.
func GetFeedStats(m *lifecycle.Manager, feedSlug string) *FeedStats {
	cacheKey := fmt.Sprintf("feed_stats.%s", feedSlug)
	if stats, ok := m.Cache().Get(cacheKey); ok {
		if feedStats, ok := stats.(*FeedStats); ok {
			return feedStats
		}
	}
	return nil
}

// GetSiteStats retrieves site-wide statistics from the cache.
func GetSiteStats(m *lifecycle.Manager) *SiteStats {
	if stats, ok := m.Cache().Get("site_stats"); ok {
		if siteStats, ok := stats.(*SiteStats); ok {
			return siteStats
		}
	}
	return nil
}

// Ensure StatsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*StatsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*StatsPlugin)(nil)
	_ lifecycle.TransformPlugin = (*StatsPlugin)(nil)
	_ lifecycle.CollectPlugin   = (*StatsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*StatsPlugin)(nil)
)

// StatsHelper provides template-friendly access to site statistics.
// It exposes methods and properties that can be used in Jinja2-style templates
// for building analytics dashboards and "year in review" style posts.
//
// Template usage examples:
//   - {{ stats.total_posts }} - Total number of posts
//   - {{ stats.total_words }} - Total word count
//   - {{ stats.posts_by_year }} - Map of year to post count
//   - {{ stats.kpi("total_posts") }} - Get a specific KPI value
//   - {{ stats.for_feed("blog").total_posts }} - Feed-specific stats
type StatsHelper struct {
	manager   *lifecycle.Manager
	siteStats *SiteStats
}

// NewStatsHelper creates a new stats helper for template access.
func NewStatsHelper(m *lifecycle.Manager) *StatsHelper {
	return &StatsHelper{
		manager:   m,
		siteStats: GetSiteStats(m),
	}
}

// TotalPosts returns the total number of posts.
func (h *StatsHelper) TotalPosts() int {
	if h.siteStats == nil {
		return 0
	}
	return h.siteStats.TotalPosts
}

// TotalWords returns the total word count.
func (h *StatsHelper) TotalWords() int {
	if h.siteStats == nil {
		return 0
	}
	return h.siteStats.TotalWords
}

// TotalReadingTime returns the total reading time in minutes.
func (h *StatsHelper) TotalReadingTime() int {
	if h.siteStats == nil {
		return 0
	}
	return h.siteStats.TotalReadingTime
}

// TotalReadingTimeText returns the formatted total reading time.
func (h *StatsHelper) TotalReadingTimeText() string {
	if h.siteStats == nil {
		return statsZeroMin
	}
	return h.siteStats.TotalReadingTimeText
}

// AverageWords returns the average word count per post.
func (h *StatsHelper) AverageWords() int {
	if h.siteStats == nil {
		return 0
	}
	return h.siteStats.AverageWords
}

// AverageReadingTime returns the average reading time per post.
func (h *StatsHelper) AverageReadingTime() int {
	if h.siteStats == nil {
		return 0
	}
	return h.siteStats.AverageReadingTime
}

// TotalCodeLines returns the total lines of code.
func (h *StatsHelper) TotalCodeLines() int {
	if h.siteStats == nil {
		return 0
	}
	return h.siteStats.TotalCodeLines
}

// TotalCodeBlocks returns the total number of code blocks.
func (h *StatsHelper) TotalCodeBlocks() int {
	if h.siteStats == nil {
		return 0
	}
	return h.siteStats.TotalCodeBlocks
}

// PostsByYear returns a map of year to post count.
func (h *StatsHelper) PostsByYear() map[int]int {
	if h.siteStats == nil {
		return make(map[int]int)
	}
	return h.siteStats.PostsByYear
}

// WordsByYear returns a map of year to total word count.
func (h *StatsHelper) WordsByYear() map[int]int {
	if h.siteStats == nil {
		return make(map[int]int)
	}
	return h.siteStats.WordsByYear
}

// PostsByTag returns a map of tag name to post count.
func (h *StatsHelper) PostsByTag() map[string]int {
	if h.siteStats == nil {
		return make(map[string]int)
	}
	return h.siteStats.PostsByTag
}

// Kpi returns a specific KPI value by name.
// Supported KPIs: total_posts, total_words, total_reading_time, average_words,
// average_reading_time, total_code_lines, total_code_blocks
func (h *StatsHelper) Kpi(name string) interface{} {
	if h.siteStats == nil {
		return 0
	}
	switch name {
	case "total_posts":
		return h.siteStats.TotalPosts
	case "total_words":
		return h.siteStats.TotalWords
	case "total_reading_time":
		return h.siteStats.TotalReadingTime
	case "total_reading_time_text":
		return h.siteStats.TotalReadingTimeText
	case "average_words":
		return h.siteStats.AverageWords
	case "average_reading_time":
		return h.siteStats.AverageReadingTime
	case "average_reading_time_text":
		return h.siteStats.AverageReadingTimeText
	case "total_code_lines":
		return h.siteStats.TotalCodeLines
	case "total_code_blocks":
		return h.siteStats.TotalCodeBlocks
	default:
		return nil
	}
}

// ForFeed returns a FeedStatsHelper for feed-specific statistics.
func (h *StatsHelper) ForFeed(feedName string) *FeedStatsHelper {
	feedStats := GetFeedStats(h.manager, feedName)
	return &FeedStatsHelper{feedStats: feedStats}
}

// FeedStatsHelper provides template-friendly access to feed-specific statistics.
type FeedStatsHelper struct {
	feedStats *FeedStats
}

// PostCount returns the number of posts in the feed.
func (h *FeedStatsHelper) PostCount() int {
	if h.feedStats == nil {
		return 0
	}
	return h.feedStats.PostCount
}

// TotalWords returns the total word count for the feed.
func (h *FeedStatsHelper) TotalWords() int {
	if h.feedStats == nil {
		return 0
	}
	return h.feedStats.TotalWords
}

// TotalReadingTime returns the total reading time for the feed.
func (h *FeedStatsHelper) TotalReadingTime() int {
	if h.feedStats == nil {
		return 0
	}
	return h.feedStats.TotalReadingTime
}

// TotalReadingTimeText returns the formatted total reading time.
func (h *FeedStatsHelper) TotalReadingTimeText() string {
	if h.feedStats == nil {
		return statsZeroMin
	}
	return h.feedStats.TotalReadingTimeText
}

// AverageWords returns the average word count per post.
func (h *FeedStatsHelper) AverageWords() int {
	if h.feedStats == nil {
		return 0
	}
	return h.feedStats.AverageWords
}

// AverageReadingTime returns the average reading time per post.
func (h *FeedStatsHelper) AverageReadingTime() int {
	if h.feedStats == nil {
		return 0
	}
	return h.feedStats.AverageReadingTime
}

// TotalCodeLines returns the total lines of code in the feed.
func (h *FeedStatsHelper) TotalCodeLines() int {
	if h.feedStats == nil {
		return 0
	}
	return h.feedStats.TotalCodeLines
}

// TotalCodeBlocks returns the total number of code blocks in the feed.
func (h *FeedStatsHelper) TotalCodeBlocks() int {
	if h.feedStats == nil {
		return 0
	}
	return h.feedStats.TotalCodeBlocks
}

// PostsByYear returns a map of year to post count for this feed.
func (h *FeedStatsHelper) PostsByYear() map[int]int {
	if h.feedStats == nil {
		return make(map[int]int)
	}
	return h.feedStats.PostsByYear
}

// WordsByYear returns a map of year to word count for this feed.
func (h *FeedStatsHelper) WordsByYear() map[int]int {
	if h.feedStats == nil {
		return make(map[int]int)
	}
	return h.feedStats.WordsByYear
}

// PostsByTag returns a map of tag to post count for this feed.
func (h *FeedStatsHelper) PostsByTag() map[string]int {
	if h.feedStats == nil {
		return make(map[string]int)
	}
	return h.feedStats.PostsByTag
}

// Kpi returns a specific KPI value by name for this feed.
func (h *FeedStatsHelper) Kpi(name string) interface{} {
	if h.feedStats == nil {
		return 0
	}
	switch name {
	case "post_count", "total_posts":
		return h.feedStats.PostCount
	case "total_words":
		return h.feedStats.TotalWords
	case "total_reading_time":
		return h.feedStats.TotalReadingTime
	case "total_reading_time_text":
		return h.feedStats.TotalReadingTimeText
	case "average_words":
		return h.feedStats.AverageWords
	case "average_reading_time":
		return h.feedStats.AverageReadingTime
	case "average_reading_time_text":
		return h.feedStats.AverageReadingTimeText
	case "total_code_lines":
		return h.feedStats.TotalCodeLines
	case "total_code_blocks":
		return h.feedStats.TotalCodeBlocks
	default:
		return nil
	}
}

// GetStatsHelper retrieves the stats helper from the cache.
func GetStatsHelper(m *lifecycle.Manager) *StatsHelper {
	if helper, ok := m.Cache().Get("stats_helper"); ok {
		if statsHelper, ok := helper.(*StatsHelper); ok {
			return statsHelper
		}
	}
	return nil
}
