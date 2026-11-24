// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/filter"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// FeedsPlugin processes configured feeds during the collect stage.
// It filters posts, sorts them, and paginates the results.
type FeedsPlugin struct{}

// NewFeedsPlugin creates a new FeedsPlugin.
func NewFeedsPlugin() *FeedsPlugin {
	return &FeedsPlugin{}
}

// Name returns the unique name of the plugin.
func (p *FeedsPlugin) Name() string {
	return "feeds"
}

// Collect processes each FeedConfig and creates feeds with filtered, sorted, and paginated posts.
func (p *FeedsPlugin) Collect(m *lifecycle.Manager) error {
	posts := m.Posts()
	config := m.Config()
	useIncremental := false
	if config.Extra != nil {
		if enabled, ok := config.Extra["feeds_incremental"].(bool); ok && enabled {
			useIncremental = true
		}
	}
	filterCache := newFeedFilterCache(posts)
	selection := newFeedSelectionRecorder(m.ContentLedger(), posts)

	// Get feed configs from manager's extra config
	feedConfigs := getFeedConfigs(config)
	feedDefaults := getFeedDefaults(config)

	feeds := make([]*lifecycle.Feed, 0, len(feedConfigs))

	for i := range feedConfigs {
		fc := &feedConfigs[i]
		// Apply defaults before the incremental decision so the diagnostic
		// window matches the same effective feed configuration used below.
		fc.ApplyDefaults(feedDefaults)
		if useIncremental && !hasFeedLimitOffset(fc) && p.shouldSkipFeedCollect(fc, m) {
			selection.record(fc.Slug, selection.sources, fc.Posts, fc.Posts, fc.Filter, fc.IncludesPrivate(), 0, 0)
			continue
		}

		usePresetPosts := fc.Type == models.FeedTypeSeries && len(fc.Posts) > 0

		var filteredPosts []*models.Post
		if usePresetPosts {
			// Series feeds may provide an explicit post list. Apply the same
			// privacy and renderability boundary as filter-based feeds before
			// limits and pagination are calculated.
			filteredPosts = filterFeedPagePosts(fc.Posts, fc.IncludesPrivate())
		} else {
			// Filter posts
			var err error
			filteredPosts, err = filterCache.FilterPosts(fc.Filter, fc.IncludesPrivate())
			if err != nil {
				return fmt.Errorf("feed %q: %w", fc.Slug, err)
			}
			filteredPosts = cloneFeedPosts(filteredPosts)
		}
		if !usePresetPosts {
			// Determine sort field and direction based on feed type
			sortField := fc.Sort
			reverse := fc.Reverse

			// For guide-type feeds, default to sorting by guide_order (ascending)
			if fc.Type == models.FeedTypeGuide {
				if sortField == "" {
					sortField = "guide_order"
					reverse = false // Guides should be in ascending order by default
				}
			} else {
				// Default to date sorting for non-guide feeds
				if sortField == "" {
					sortField = "date"
					reverse = true // Newest first by default
				}
			}

			sortPosts(filteredPosts, sortField, reverse)

			// For guide-type feeds, set up prev/next navigation on each post
			if fc.Type == models.FeedTypeGuide || fc.Type == models.FeedTypeSeries {
				setGuideNavigation(filteredPosts, fc.Slug)
			}
		}

		consideredSources := selection.sources
		if usePresetPosts {
			consideredSources = prepareFeedSources(fc.Posts)
		}
		matchedPosts := filteredPosts

		// Store posts in feed config
		filteredPosts = applyFeedLimitOffset(filteredPosts, fc)
		fc.Posts = filteredPosts
		selection.record(fc.Slug, consideredSources, matchedPosts, filteredPosts, fc.Filter, fc.IncludesPrivate(), fc.Offset, fc.Limit)

		// Get base URL for pagination
		baseURL := "/" + fc.Slug
		if fc.Slug == "" {
			baseURL = ""
		}

		// Paginate results
		fc.Paginate(baseURL)

		// Create lifecycle.Feed for each page
		feed := &lifecycle.Feed{
			Name:           fc.Slug,
			Title:          fc.Title,
			Posts:          filteredPosts,
			Path:           fc.Slug,
			IncludePrivate: fc.IncludesPrivate(),
			Hidden:         fc.Sidebar != nil && !*fc.Sidebar,
		}

		feeds = append(feeds, feed)
	}

	m.SetFeeds(feeds)

	// Store feed configs back in cache for publish_feeds to access
	m.Cache().Set("feed_configs", feedConfigs)

	return nil
}

func hasFeedLimitOffset(fc *models.FeedConfig) bool {
	return fc != nil && (fc.Offset > 0 || fc.Limit > 0)
}

func recordFeedSelection(m *lifecycle.Manager, feedName string, selected []*models.Post, filterExpr string, includePrivate bool) {
	if m == nil {
		return
	}
	recordFeedSelectionForPosts(m, feedName, m.Posts(), selected, selected, filterExpr, includePrivate, 0, 0)
}

func recordFeedSelectionForPosts(
	m *lifecycle.Manager,
	feedName string,
	considered, matched, selected []*models.Post,
	filterExpr string,
	includePrivate bool,
	offset, limit int,
) {
	if m == nil || m.ContentLedger() == nil {
		return
	}
	recorder := newFeedSelectionRecorder(m.ContentLedger(), considered)
	recorder.record(feedName, recorder.sources, matched, selected, filterExpr, includePrivate, offset, limit)
}

// Prepared per occurrence, not per pointer or path. These flags are valid only
// during the synchronous Collect hook; every invocation refreshes them.
type feedSelectionSource struct {
	path  string
	flags uint8
}

const (
	feedUnpublished uint8 = 1 << iota
	feedSkipped
	feedDraft
	feedPrivate
)

func prepareFeedSources(posts []*models.Post) []feedSelectionSource {
	sources := make([]feedSelectionSource, 0, len(posts))
	for _, post := range posts {
		if post == nil || post.Path == "" {
			continue
		}
		var flags uint8
		if !post.Published {
			flags |= feedUnpublished
		}
		if post.Skip {
			flags |= feedSkipped
		}
		if post.Draft {
			flags |= feedDraft
		}
		if post.Private {
			flags |= feedPrivate
		}
		sources = append(sources, feedSelectionSource{path: post.Path, flags: flags})
	}
	return sources
}

type feedSelectionRecorder struct {
	ledger        *diagnostics.ContentLedger
	sources       []feedSelectionSource
	matchedPaths  map[string]int
	selectedPaths map[string]struct{}
	observations  []diagnostics.ContentFeedObservation
	reasons       []string
}

func newFeedSelectionRecorder(ledger *diagnostics.ContentLedger, posts []*models.Post) *feedSelectionRecorder {
	return &feedSelectionRecorder{
		ledger: ledger, sources: prepareFeedSources(posts),
		matchedPaths: make(map[string]int), selectedPaths: make(map[string]struct{}),
	}
}

func (r *feedSelectionRecorder) record(feedName string, sources []feedSelectionSource, matched, selected []*models.Post, filterExpr string, includePrivate bool, offset, limit int) {
	if r.ledger == nil {
		return
	}
	clear(r.matchedPaths)
	clear(r.selectedPaths)
	for index, post := range matched {
		if post == nil || post.Path == "" {
			continue
		}
		if _, exists := r.matchedPaths[post.Path]; !exists {
			r.matchedPaths[post.Path] = index
		}
	}

	for _, post := range selected {
		if post != nil && post.Path != "" {
			r.selectedPaths[post.Path] = struct{}{}
		}
	}

	// At most four eligibility reasons and one window/fallback reason per
	// occurrence. Pre-sizing keeps each row in the same reusable arena.
	if cap(r.reasons) < 5*len(sources) {
		r.reasons = make([]string, 0, 5*len(sources))
	}
	clear(r.observations)
	if cap(r.observations) < len(sources) {
		r.observations = make([]diagnostics.ContentFeedObservation, 0, len(sources))
	}
	r.observations = r.observations[:0]
	r.reasons = r.reasons[:0]
	for _, source := range sources {
		start := len(r.reasons)
		flags := source.flags
		if includePrivate {
			flags &^= feedPrivate
		}
		if flags != 0 {
			for i, reason := range [...]string{
				diagnostics.ReasonContentPublishedFalse, diagnostics.ReasonContentSkip,
				diagnostics.ReasonContentDraft, diagnostics.ReasonContentPrivate,
			} {
				if flags&(1<<i) != 0 {
					r.reasons = append(r.reasons, reason)
				}
			}
		}
		index, matchedInFeed := r.matchedPaths[source.path]
		included := false
		if matchedInFeed {
			_, included = r.selectedPaths[source.path]
		}
		reasons := r.reasons[start:]
		switch {
		case !matchedInFeed:
			reasons = feedUnmatchedReasons(reasons, filterExpr)
			included = false
		case !included:
			reasons = feedWindowExcludedReasons(reasons, index, len(matched), offset, limit)
		default:
			included = len(reasons) == 0
		}
		r.reasons = r.reasons[:start+len(reasons)]
		r.observations = append(r.observations, diagnostics.ContentFeedObservation{
			Path: source.path, Included: included, Reasons: reasons[:len(reasons):len(reasons)],
		})
	}
	r.ledger.RecordFeedBatch(feedName, r.observations)
}

func feedUnmatchedReasons(reasons []string, filterExpr string) []string {
	if filterExpr != "" && len(reasons) == 0 {
		reasons = append(reasons, diagnostics.ReasonContentFiltered)
	}
	if len(reasons) == 0 {
		reasons = append(reasons, diagnostics.ReasonContentNoOutput)
	}
	return reasons
}

func feedWindowExcludedReasons(reasons []string, index, matchedCount, offset, limit int) []string {
	if reason := feedWindowReason(index, matchedCount, offset, limit); reason != "" {
		reasons = append(reasons, reason)
	}
	if len(reasons) == 0 {
		reasons = append(reasons, diagnostics.ReasonContentNoOutput)
	}
	return reasons
}

func feedEligibilityReasons(post *models.Post, includePrivate bool) []string {
	reasons := make([]string, 0, 4)
	if !post.Published {
		reasons = append(reasons, diagnostics.ReasonContentPublishedFalse)
	}
	if post.Skip {
		reasons = append(reasons, diagnostics.ReasonContentSkip)
	}
	if post.Draft {
		reasons = append(reasons, diagnostics.ReasonContentDraft)
	}
	if post.Private && !includePrivate {
		reasons = append(reasons, diagnostics.ReasonContentPrivate)
	}
	return reasons
}

func feedWindowReason(index, matchedCount, offset, limit int) string {
	if offset < 0 {
		offset = 0
	}
	if limit < 0 {
		limit = 0
	}
	if index < offset || offset >= matchedCount {
		return diagnostics.ReasonFeedOffset
	}
	if limit > 0 && index >= offset+limit {
		return diagnostics.ReasonFeedLimit
	}
	return ""
}

func applyFeedLimitOffset(posts []*models.Post, fc *models.FeedConfig) []*models.Post {
	if len(posts) == 0 {
		return posts
	}

	offset := fc.Offset
	if offset < 0 {
		offset = 0
	}

	limit := fc.Limit
	if limit < 0 {
		limit = 0
	}

	if offset == 0 && limit == 0 {
		return posts
	}

	if offset >= len(posts) {
		return posts[:0]
	}

	limited := posts
	if offset > 0 {
		limited = limited[offset:]
	}

	if limit == 0 {
		return limited
	}

	if limit >= len(limited) {
		return limited
	}

	return limited[:limit]
}

func (p *FeedsPlugin) shouldSkipFeedCollect(fc *models.FeedConfig, m *lifecycle.Manager) bool {
	cache := GetBuildCache(m)
	if cache == nil {
		return false
	}
	if len(fc.Posts) == 0 {
		// Feed membership is runtime state and is not serialized in the site
		// configuration. A fresh incremental manager therefore cannot safely
		// skip collection when Posts is empty: doing so would replace a cached
		// feed with an empty one and would make feed diagnostics incorrect.
		return false
	}

	for _, post := range fc.Posts {
		currentHash := computePostFeedItemHash(post)
		cachedHash, _, _ := cache.GetPostSemanticHashes(post.Path)
		if cachedHash != currentHash {
			return false
		}
	}

	return true
}

type feedFilterKey struct {
	expr           string
	includePrivate bool
}

type feedFilterCache struct {
	parsed      map[string]*filter.Filter
	filtered    map[feedFilterKey][]*models.Post
	publicPosts []*models.Post
	allPosts    []*models.Post
}

func newFeedFilterCache(posts []*models.Post) *feedFilterCache {
	publicPosts := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if !post.Private {
			publicPosts = append(publicPosts, post)
		}
	}

	return &feedFilterCache{
		parsed:      make(map[string]*filter.Filter),
		filtered:    make(map[feedFilterKey][]*models.Post),
		publicPosts: publicPosts,
		allPosts:    posts,
	}
}

func (c *feedFilterCache) FilterPosts(filterExpr string, includePrivate bool) ([]*models.Post, error) {
	key := feedFilterKey{expr: filterExpr, includePrivate: includePrivate}
	if cached, ok := c.filtered[key]; ok {
		return cached, nil
	}

	candidatePosts := c.publicPosts
	if includePrivate {
		candidatePosts = c.allPosts
	}

	if filterExpr == "" {
		c.filtered[key] = candidatePosts
		return candidatePosts, nil
	}

	f, ok := c.parsed[filterExpr]
	if !ok {
		var err error
		f, err = filter.Parse(filterExpr)
		if err != nil {
			return nil, fmt.Errorf("invalid filter expression: %w", err)
		}
		c.parsed[filterExpr] = f
	}

	result := f.MatchAll(candidatePosts)
	c.filtered[key] = result
	return result, nil
}

// getFeedConfigs retrieves feed configurations from the manager config.
func getFeedConfigs(config *lifecycle.Config) []models.FeedConfig {
	if config.Extra == nil {
		return nil
	}

	if feeds, ok := config.Extra["feeds"]; ok {
		if feedConfigs, ok := feeds.([]models.FeedConfig); ok {
			return feedConfigs
		}
	}

	return nil
}

// getFeedDefaults retrieves feed defaults from the manager config.
func getFeedDefaults(config *lifecycle.Config) models.FeedDefaults {
	if config.Extra == nil {
		return models.NewFeedDefaults()
	}

	if defaults, ok := config.Extra["feed_defaults"]; ok {
		if feedDefaults, ok := defaults.(models.FeedDefaults); ok {
			return feedDefaults
		}
	}

	return models.NewFeedDefaults()
}

// filterPosts applies a filter expression to posts.
// If includePrivate is false, private posts are excluded before filtering.
func filterPosts(posts []*models.Post, filterExpr string, includePrivate bool) ([]*models.Post, error) {
	// First, filter out private posts if includePrivate is false
	var candidatePosts []*models.Post
	if includePrivate {
		candidatePosts = posts
	} else {
		for _, post := range posts {
			if !post.Private {
				candidatePosts = append(candidatePosts, post)
			}
		}
	}

	if filterExpr == "" {
		// Return a copy of all candidate posts
		result := make([]*models.Post, len(candidatePosts))
		copy(result, candidatePosts)
		return result, nil
	}

	f, err := filter.Parse(filterExpr)
	if err != nil {
		return nil, fmt.Errorf("invalid filter expression: %w", err)
	}

	return f.MatchAll(candidatePosts), nil
}

func cloneFeedPosts(posts []*models.Post) []*models.Post {
	if len(posts) == 0 {
		return nil
	}
	result := make([]*models.Post, len(posts))
	copy(result, posts)
	return result
}

// sortPosts sorts posts by the specified field.
// When primary sort values are equal, uses Path as a secondary sort key
// to ensure deterministic ordering across builds.
func sortPosts(posts []*models.Post, field string, reverse bool) {
	sort.SliceStable(posts, func(i, j int) bool {
		vi := getFieldValue(posts[i], field)
		vj := getFieldValue(posts[j], field)

		cmp := compareFieldValues(vi, vj)

		// Use path as tie-breaker for deterministic ordering
		if cmp == 0 {
			cmp = strings.Compare(posts[i].Path, posts[j].Path)
		}

		if reverse {
			return cmp > 0
		}
		return cmp < 0
	})
}

// getFieldValue retrieves a field value from a post.
func getFieldValue(post *models.Post, field string) interface{} {
	// Check Extra fields first
	if post.Extra != nil {
		if v, ok := post.Extra[field]; ok {
			return v
		}
	}

	// Use reflection for struct fields
	v := reflect.ValueOf(post).Elem()
	t := v.Type()

	// Try to find field by name (case-insensitive)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if strings.EqualFold(f.Name, field) {
			fv := v.Field(i)
			if fv.Kind() == reflect.Ptr {
				if fv.IsNil() {
					return nil
				}
				return fv.Elem().Interface()
			}
			return fv.Interface()
		}
	}

	return nil
}

// compareFieldValues compares two field values for sorting.
func compareFieldValues(a, b interface{}) int {
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Compare time.Time
	if ta, ok := a.(time.Time); ok {
		if tb, ok := b.(time.Time); ok {
			if ta.Before(tb) {
				return -1
			}
			if ta.After(tb) {
				return 1
			}
			return 0
		}
	}

	// Compare strings
	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			return strings.Compare(sa, sb)
		}
	}

	// Compare integers (for guide_order)
	if ia, ok := toInt(a); ok {
		if ib, ok := toInt(b); ok {
			if ia < ib {
				return -1
			}
			if ia > ib {
				return 1
			}
			return 0
		}
	}

	// Compare as formatted strings
	return strings.Compare(fmt.Sprintf("%v", a), fmt.Sprintf("%v", b))
}

// toInt attempts to convert an interface{} to an int.
// It handles int, int64, float64, and other common numeric types.
func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case int32:
		return int(n), true
	case float64:
		return int(n), true
	case float32:
		return int(n), true
	default:
		return 0, false
	}
}

// setGuideNavigation sets up Prev/Next navigation for posts in a guide series.
// Posts should already be sorted in the desired order before calling this function.
func setGuideNavigation(posts []*models.Post, feedSlug string) {
	for i, post := range posts {
		// Set feed reference for navigation context
		post.PrevNextFeed = feedSlug

		// Set Prev pointer
		if i > 0 {
			post.Prev = posts[i-1]
		} else {
			post.Prev = nil
		}

		// Set Next pointer
		if i < len(posts)-1 {
			post.Next = posts[i+1]
		} else {
			post.Next = nil
		}
	}
}

// Ensure FeedsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin        = (*FeedsPlugin)(nil)
	_ lifecycle.CollectPlugin = (*FeedsPlugin)(nil)
)
