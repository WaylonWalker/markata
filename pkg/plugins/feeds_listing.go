// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

const homeFeedTitle = "Home"

const (
	feedVariantPage    = "page"
	feedVariantExport  = "export"
	feedVariantFeed    = DefaultFeedPath
	feedVariantArchive = defaultArchivePrefix
)

// FeedVariantLink describes one publicly accessible output for a feed.
type FeedVariantLink struct {
	Label string
	Href  string
	Kind  string
}

type FeedListingSection struct {
	ID          string
	Title       string
	Description string
	TotalCount  int
	MoreHref    string
	MoreLabel   string
	Pagination  *FeedListingPage
	Feeds       []FeedListingInfo
}

type FeedListingPage struct {
	Number       int
	Feeds        []FeedListingInfo
	HasPrev      bool
	HasNext      bool
	PrevURL      string
	NextURL      string
	TotalPages   int
	TotalItems   int
	ItemsPerPage int
	PageURLs     []string
}

type sparklineWindow struct {
	Start time.Time
	End   time.Time
}

type SparklinePoint struct {
	X     float64
	Y     float64
	HitX  float64
	Month string
	Value int
}

// FeedListingInfo contains the data rendered on the /feeds page.
type FeedListingInfo struct {
	Title             string
	Slug              string
	Description       string
	Href              string
	PostCount         int
	LatestPostDate    string
	SubscribeCount    int
	ArchiveCount      int
	DisplayVariants   []FeedVariantLink
	PrimaryVariants   []FeedVariantLink
	ArchiveVariants   []FeedVariantLink
	UtilityVariants   []FeedVariantLink
	LatestPostTime    time.Time
	SparklinePoints   string
	SparklineData     []SparklinePoint
	SparklineValues   string
	SparklineTitle    string
	SparklineSummary  string
	SparklineStart    string
	SparklineStartISO string
	SparklineEnd      string
	GeneratedBySite   bool
}

// FeedsListingPlugin generates a feeds listing page at /feeds.
type FeedsListingPlugin struct {
	engineMu    sync.RWMutex
	engineCache map[string]*templates.Engine
}

// NewFeedsListingPlugin creates a new FeedsListingPlugin.
func NewFeedsListingPlugin() *FeedsListingPlugin {
	return &FeedsListingPlugin{engineCache: make(map[string]*templates.Engine)}
}

// Name returns the unique name of the plugin.
func (p *FeedsListingPlugin) Name() string {
	return "feeds_listing"
}

// Priority returns the plugin's priority for a given stage.
func (p *FeedsListingPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Write generates the feeds listing page.
func (p *FeedsListingPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	feedsPage := getFeedsPageConfig(config)
	if !feedsPage.IsEnabled() {
		return nil
	}

	feedConfigs := getCachedFeedConfigs(m)
	if len(feedConfigs) == 0 {
		return nil
	}
	feedDefaults := getFeedDefaults(config)
	cache := GetBuildCache(m)
	listingHash := computeFeedsListingHash(feedConfigs, config, &feedsPage, feedDefaults, cache)
	if cache != nil {
		if cached := cache.GetFeedsListingHash(); cached != "" && cached == listingHash {
			if feedsListingOutputsExist(config, &feedsPage, feedConfigs, feedDefaults) {
				return nil
			}
		}
	}
	sparklineRange := computeSparklineWindow(m.Posts(), false)

	sections, generatedFeedPages := p.collectFeedSections(feedConfigs, config, &feedsPage, feedDefaults, sparklineRange)
	if len(sections) == 0 {
		return nil
	}

	if err := p.renderFeedsPage(config, &feedsPage, sections, feedsPage.SlugPrefix, feedsPage.Title, feedsPage.Description, nil, nil); err != nil {
		return err
	}

	if len(generatedFeedPages) > 0 {
		generatedTitle := "Generated Feeds"
		generatedDescription := "Automatically updated feeds for broader site sections, archives, and collections."
		for i := range generatedFeedPages {
			pageSlug := filepath.ToSlash(filepath.Join(feedsPage.SlugPrefix, "generated"))
			if generatedFeedPages[i].Number > 1 {
				pageSlug = filepath.ToSlash(filepath.Join(pageSlug, "page", fmt.Sprintf("%d", generatedFeedPages[i].Number)))
			}
			section := FeedListingSection{
				ID:          "generated-feeds",
				Title:       generatedTitle,
				Description: generatedDescription,
				TotalCount:  generatedFeedPages[i].TotalItems,
				Feeds:       generatedFeedPages[i].Feeds,
			}
			if err := p.renderFeedsPage(
				config,
				&feedsPage,
				[]FeedListingSection{section},
				pageSlug,
				generatedTitle,
				generatedDescription,
				nil,
				&generatedFeedPages[i],
			); err != nil {
				return err
			}
		}
	}

	if cache != nil {
		cache.SetFeedsListingHash(listingHash)
	}

	return nil
}

func computeFeedsListingHash(
	feedConfigs []models.FeedConfig,
	config *lifecycle.Config,
	feedsPage *models.FeedsPageConfig,
	feedDefaults models.FeedDefaults,
	cache *buildcache.Cache,
) string {
	h := sha256.New()
	writeString := func(value string) {
		fmt.Fprintf(h, "%d:", len(value))
		h.Write([]byte(value))
		h.Write([]byte{0})
	}
	writeBool := func(value bool) {
		fmt.Fprintf(h, "%t", value)
		h.Write([]byte{0})
	}
	writeInt := func(value int) {
		fmt.Fprintf(h, "%d", value)
		h.Write([]byte{0})
	}

	writeBool(feedsPage != nil && feedsPage.IsEnabled())
	if feedsPage != nil {
		writeString(feedsPage.Title)
		writeString(feedsPage.Description)
		writeString(feedsPage.Template)
		writeString(feedsPage.SlugPrefix)
		writeString(feedsPage.Robots)
		for _, slug := range feedsPage.ShowPrivateFeeds {
			writeString(slug)
		}
	}

	writeInt(feedDefaults.ItemsPerPage)
	writeInt(feedDefaults.OrphanThreshold)

	configuredSlugs := configuredFeedSlugs(config)
	configuredOrder := make([]string, 0, len(configuredSlugs))
	for slug := range configuredSlugs {
		configuredOrder = append(configuredOrder, slug)
	}
	sort.Slice(configuredOrder, func(i, j int) bool {
		return configuredSlugs[configuredOrder[i]] < configuredSlugs[configuredOrder[j]]
	})
	for _, slug := range configuredOrder {
		writeString(slug)
	}

	for i := range feedConfigs {
		fc := &feedConfigs[i]
		writeString(fc.Slug)
		if cache != nil {
			if feedHash := cache.GetFeedHash(fc.Slug); feedHash != "" {
				writeString(feedHash)
				continue
			}
		}
		writeString(fc.Title)
		writeString(fc.Description)
		writeBool(fc.IncludesPrivate())
		writeBool(fc.Formats.HTML)
		writeBool(fc.Formats.SimpleHTML)
		writeBool(fc.Formats.RSS)
		writeBool(fc.Formats.Atom)
		writeBool(fc.Formats.JSON)
		writeBool(fc.Formats.Markdown)
		writeBool(fc.Formats.Text)
		writeBool(fc.Formats.Sitemap)
		for _, post := range fc.Posts {
			if post == nil {
				continue
			}
			writeString(post.Slug)
			if post.Date != nil {
				writeString(post.Date.UTC().Format(time.RFC3339Nano))
			}
			writeBool(post.Published)
			writeBool(post.Draft)
			writeBool(post.Private)
			writeBool(post.Skip)
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}

func feedsListingOutputsExist(config *lifecycle.Config, feedsPage *models.FeedsPageConfig, feedConfigs []models.FeedConfig, feedDefaults models.FeedDefaults) bool {
	for _, outputPath := range expectedFeedsListingOutputPaths(config, feedsPage, feedConfigs, feedDefaults) {
		if _, err := os.Stat(outputPath); err != nil {
			return false
		}
	}
	return true
}

func expectedFeedsListingOutputPaths(config *lifecycle.Config, feedsPage *models.FeedsPageConfig, feedConfigs []models.FeedConfig, feedDefaults models.FeedDefaults) []string {
	if config == nil || feedsPage == nil {
		return nil
	}

	paths := []string{filepath.Join(config.OutputDir, filepath.FromSlash(feedsPage.SlugPrefix), "index.html")}
	generatedCount := 0
	configuredSlugs := configuredFeedSlugs(config)
	for i := range feedConfigs {
		fc := &feedConfigs[i]
		if fc.IncludesPrivate() && !feedsPageAllowsPrivateFeed(feedsPage, fc.Slug) {
			continue
		}
		if _, isConfigured := configuredSlugs[fc.Slug]; !isConfigured {
			generatedCount++
		}
	}

	generatedPages := feedListingPageCount(generatedCount, feedDefaults)
	for pageNum := 1; pageNum <= generatedPages; pageNum++ {
		if pageNum == 1 {
			paths = append(paths, filepath.Join(config.OutputDir, filepath.FromSlash(feedsPage.SlugPrefix), "generated", "index.html"))
			continue
		}
		paths = append(paths, filepath.Join(config.OutputDir, filepath.FromSlash(feedsPage.SlugPrefix), "generated", "page", fmt.Sprintf("%d", pageNum), "index.html"))
	}

	return paths
}

func feedListingPageCount(total int, defaults models.FeedDefaults) int {
	if total <= 0 {
		return 0
	}
	itemsPerPage := defaults.ItemsPerPage
	if itemsPerPage <= 0 {
		itemsPerPage = 10
	}
	orphanThreshold := defaults.OrphanThreshold
	if orphanThreshold <= 0 {
		orphanThreshold = 3
	}
	pages := 0
	for i := 0; i < total; i += itemsPerPage {
		end := i + itemsPerPage
		if end > total {
			end = total
		}
		remaining := total - end
		if remaining > 0 && remaining < orphanThreshold {
			end = total
		}
		pages++
		if end >= total {
			break
		}
	}
	return pages
}

func (p *FeedsListingPlugin) collectFeedSections(
	feedConfigs []models.FeedConfig,
	config *lifecycle.Config,
	feedsPage *models.FeedsPageConfig,
	feedDefaults models.FeedDefaults,
	sparklineRange sparklineWindow,
) ([]FeedListingSection, []FeedListingPage) {
	syndication := getSyndicationConfig(config)
	userDefined := make([]FeedListingInfo, 0, len(feedConfigs))
	generated := make([]FeedListingInfo, 0, len(feedConfigs))
	configuredSlugs := configuredFeedSlugs(config)

	for i := range feedConfigs {
		fc := &feedConfigs[i]
		if fc.IncludesPrivate() && !feedsPageAllowsPrivateFeed(feedsPage, fc.Slug) {
			continue
		}

		postCount, latestDate, latestTime := feedStats(fc.Posts, fc.IncludesPrivate())
		display, primary, archive, utility := splitFeedVariants(fc, syndication)
		_, isConfigured := configuredSlugs[fc.Slug]
		sparklineData := buildFeedSparklineData(fc.Posts, sparklineRange, fc.IncludesPrivate())
		info := FeedListingInfo{
			Title:             feedDisplayTitle(fc),
			Slug:              fc.Slug,
			Description:       fc.Description,
			Href:              feedHTMLHref(fc),
			PostCount:         postCount,
			LatestPostDate:    latestDate,
			LatestPostTime:    latestTime,
			SubscribeCount:    feedSubscribeCount(fc, syndication, postCount),
			ArchiveCount:      postCount,
			DisplayVariants:   display,
			PrimaryVariants:   primary,
			ArchiveVariants:   archive,
			UtilityVariants:   utility,
			SparklinePoints:   buildFeedSparkline(fc.Posts, sparklineRange, fc.IncludesPrivate()),
			SparklineData:     sparklineData,
			SparklineValues:   encodeFeedSparklineValues(sparklineData),
			SparklineTitle:    buildFeedSparklineTitle(fc.Posts, sparklineRange, fc.IncludesPrivate()),
			SparklineSummary:  buildFeedSparklineSummary(fc.Posts, sparklineRange, fc.IncludesPrivate()),
			SparklineStart:    buildFeedSparklineStart(sparklineRange),
			SparklineStartISO: buildFeedSparklineStartISO(sparklineRange),
			SparklineEnd:      buildFeedSparklineEnd(sparklineRange),
			GeneratedBySite:   !isConfigured,
		}

		if isConfigured {
			userDefined = append(userDefined, info)
		} else {
			generated = append(generated, info)
		}
	}

	sort.SliceStable(userDefined, func(i, j int) bool {
		leftOrder, leftOK := configuredSlugs[userDefined[i].Slug]
		rightOrder, rightOK := configuredSlugs[userDefined[j].Slug]
		if leftOK && rightOK {
			return leftOrder < rightOrder
		}
		if leftOK != rightOK {
			return leftOK
		}
		return userDefined[i].Title < userDefined[j].Title
	})
	sort.SliceStable(generated, func(i, j int) bool {
		if generated[i].PostCount != generated[j].PostCount {
			return generated[i].PostCount > generated[j].PostCount
		}
		if !generated[i].LatestPostTime.Equal(generated[j].LatestPostTime) {
			return generated[i].LatestPostTime.After(generated[j].LatestPostTime)
		}
		return generated[i].Title < generated[j].Title
	})

	sections := make([]FeedListingSection, 0, 2)
	generatedPages := make([]FeedListingPage, 0)
	if len(userDefined) > 0 {
		sections = append(sections, FeedListingSection{
			ID:          "configured-feeds",
			Title:       "Curated Feeds",
			Description: "Hand-picked collections grouped around the main themes of the site.",
			TotalCount:  len(userDefined),
			Feeds:       userDefined,
		})
	}
	if len(generated) > 0 {
		generatedPages = paginateFeedListings(
			generated,
			feedDefaults,
			"/"+filepath.ToSlash(filepath.Join(feedsPage.SlugPrefix, "generated"))+"/",
			"/"+filepath.ToSlash(filepath.Join(feedsPage.SlugPrefix, "generated")),
		)
		previewPages := paginateFeedListings(
			generated,
			feedDefaults,
			"/"+feedsPage.SlugPrefix+"/",
			"/"+filepath.ToSlash(filepath.Join(feedsPage.SlugPrefix, "generated")),
		)
		preview := generated
		truncated := false
		if len(previewPages) > 0 {
			preview = previewPages[0].Feeds
			truncated = previewPages[0].TotalPages > 1
		}

		section := FeedListingSection{
			ID:          "generated-feeds",
			Title:       "Generated Feeds",
			Description: "Automatically updated feeds for broader site sections, archives, and collections.",
			TotalCount:  len(generated),
			Feeds:       preview,
			Pagination:  nil,
		}
		if truncated {
			section.Pagination = &previewPages[0]
		}
		sections = append(sections, section)
	}

	return sections, generatedPages
}

func feedsPageAllowsPrivateFeed(feedsPage *models.FeedsPageConfig, slug string) bool {
	if feedsPage == nil {
		return false
	}
	normalized := strings.Trim(slug, "/")
	for _, allowed := range feedsPage.ShowPrivateFeeds {
		if strings.Trim(allowed, "/") == normalized {
			return true
		}
	}
	return false
}

func feedStats(posts []*models.Post, includePrivate bool) (count int, latestDate string, latestTime time.Time) {
	var latest time.Time
	for _, post := range posts {
		if post == nil || post.Skip || post.Draft || !post.Published || (post.Private && !includePrivate) {
			continue
		}
		count++
		if post.Date != nil && post.Date.After(latest) {
			latest = *post.Date
		}
	}
	if latest.IsZero() {
		return count, "", time.Time{}
	}
	return count, latest.Format("2006-01-02"), latest
}

func feedDisplayTitle(fc *models.FeedConfig) string {
	if fc.Title != "" {
		return fc.Title
	}
	if fc.Slug == "" {
		return homeFeedTitle
	}
	return fc.Slug
}

func feedHTMLHref(fc *models.FeedConfig) string {
	if fc.Slug == "" {
		return "/"
	}
	return "/" + fc.Slug + "/"
}

func feedVariantLinks(fc *models.FeedConfig, syndication models.SyndicationConfig) []FeedVariantLink {
	variants := make([]FeedVariantLink, 0, 10)
	baseHref := feedHTMLHref(fc)

	if fc.Formats.HTML {
		variants = append(variants, FeedVariantLink{Label: "HTML", Href: baseHref, Kind: "page"})
	}
	if fc.Formats.SimpleHTML {
		variants = append(variants, FeedVariantLink{Label: "Simple", Href: pathJoinURL(baseHref, "simple/"), Kind: "page"})
	}
	if fc.Formats.RSS {
		variants = append(variants, FeedVariantLink{Label: "RSS", Href: pathJoinURL(baseHref, "rss.xml"), Kind: "feed"})
	}
	if fc.Formats.Atom {
		variants = append(variants, FeedVariantLink{Label: "Atom", Href: pathJoinURL(baseHref, "atom.xml"), Kind: "feed"})
	}
	if fc.Formats.JSON {
		variants = append(variants, FeedVariantLink{Label: "JSON", Href: pathJoinURL(baseHref, "feed.json"), Kind: "feed"})
	}
	if fc.Formats.Markdown && fc.Slug != "" {
		variants = append(variants, FeedVariantLink{Label: "Markdown", Href: "/" + fc.Slug + ".md", Kind: "export"})
	}
	if fc.Formats.Text && fc.Slug != "" {
		variants = append(variants, FeedVariantLink{Label: "Text", Href: "/" + fc.Slug + ".txt", Kind: "export"})
	}
	if fc.Formats.Sitemap {
		variants = append(variants, FeedVariantLink{Label: "Sitemap", Href: pathJoinURL(baseHref, "sitemap.xml"), Kind: "meta"})
	}
	if shouldGenerateFeedArchive(fc, syndication) {
		if fc.Formats.RSS {
			variants = append(variants, FeedVariantLink{Label: "Archive RSS", Href: feedArchiveURL(fc.Slug, "rss.xml"), Kind: "archive"})
		}
		if fc.Formats.Atom {
			variants = append(variants, FeedVariantLink{Label: "Archive Atom", Href: feedArchiveURL(fc.Slug, "atom.xml"), Kind: "archive"})
		}
		if fc.Formats.JSON {
			variants = append(variants, FeedVariantLink{Label: "Archive JSON", Href: feedArchiveURL(fc.Slug, "feed.json"), Kind: "archive"})
		}
	}

	return variants
}

func splitFeedVariants(fc *models.FeedConfig, syndication models.SyndicationConfig) (display, primary, archive, utility []FeedVariantLink) {
	variants := feedVariantLinks(fc, syndication)
	for _, variant := range variants {
		switch variant.Kind {
		case feedVariantPage, feedVariantExport:
			display = append(display, variant)
		case feedVariantFeed:
			primary = append(primary, variant)
		case feedVariantArchive:
			archive = append(archive, variant)
		default:
			utility = append(utility, variant)
		}
	}
	return display, primary, archive, utility
}

func feedSubscribeCount(fc *models.FeedConfig, syndication models.SyndicationConfig, totalPosts int) int {
	if fc == nil {
		return 0
	}
	if isArchiveFeed(fc) {
		return totalPosts
	}
	if syndication.MaxItems <= 0 || totalPosts < syndication.MaxItems {
		return totalPosts
	}
	return syndication.MaxItems
}

func configuredFeedSlugs(config *lifecycle.Config) map[string]int {
	slugs := make(map[string]int)
	if config == nil || config.Extra == nil {
		return slugs
	}
	feeds, ok := config.Extra["feeds"].([]models.FeedConfig)
	if !ok {
		return slugs
	}
	for i := range feeds {
		slugs[feeds[i].Slug] = i
	}
	return slugs
}

func buildFeedSparkline(posts []*models.Post, window sparklineWindow, includePrivate bool) string {
	buckets, months := monthlyPostBuckets(posts, window, includePrivate)
	if len(buckets) < 2 {
		return ""
	}
	return sparklinePolylinePoints(buckets, months)
}

func buildFeedSparklineData(posts []*models.Post, window sparklineWindow, includePrivate bool) []SparklinePoint {
	buckets, months := monthlyPostBuckets(posts, window, includePrivate)
	if len(buckets) < 2 {
		return nil
	}
	maxValue := 0
	for _, count := range buckets {
		if count > maxValue {
			maxValue = count
		}
	}
	if maxValue == 0 {
		return nil
	}

	const width = 96
	const height = 28
	step := float64(width) / float64(len(buckets)-1)
	points := make([]SparklinePoint, 0, len(buckets))
	for i, count := range buckets {
		x := float64(i) * step
		y := float64(height)
		if maxValue > 0 {
			y = float64(height) - ((float64(count) / float64(maxValue)) * float64(height-4))
		}
		points = append(points, SparklinePoint{
			X:     x,
			Y:     y,
			HitX:  x - 2.6,
			Month: months[i].Format("Jan 2006"),
			Value: count,
		})
	}
	return points
}

func sparklinePolylinePoints(buckets []int, months []time.Time) string {
	if len(buckets) == 0 || len(buckets) != len(months) {
		return ""
	}

	maxValue := 0
	for _, count := range buckets {
		if count > maxValue {
			maxValue = count
		}
	}
	if maxValue == 0 {
		return ""
	}

	const width = 96
	const height = 28
	step := float64(width) / float64(len(buckets)-1)
	points := make([]string, 0, len(buckets))
	for i, count := range buckets {
		x := float64(i) * step
		y := float64(height)
		if maxValue > 0 {
			y = float64(height) - ((float64(count) / float64(maxValue)) * float64(height-4))
		}
		points = append(points, fmt.Sprintf("%.1f,%.1f", x, y))
	}
	return strings.Join(points, " ")
}

func buildFeedSparklineTitle(posts []*models.Post, window sparklineWindow, includePrivate bool) string {
	buckets, _ := monthlyPostBuckets(posts, window, includePrivate)
	if len(buckets) == 0 {
		return ""
	}
	return fmt.Sprintf(
		"Posts published per month from %s to %s",
		window.Start.Format("2006-01"),
		window.End.Format("2006-01"),
	)
}

func buildFeedSparklineSummary(posts []*models.Post, window sparklineWindow, includePrivate bool) string {
	data := buildFeedSparklineData(posts, window, includePrivate)
	if len(data) == 0 {
		return ""
	}
	last := data[len(data)-1]
	return fmt.Sprintf("%s | %d posts", last.Month, last.Value)
}

func buildFeedSparklineStart(window sparklineWindow) string {
	if window.Start.IsZero() {
		return ""
	}
	return window.Start.Format("Jan 2006")
}

func buildFeedSparklineStartISO(window sparklineWindow) string {
	if window.Start.IsZero() {
		return ""
	}
	return window.Start.Format("2006-01")
}

func encodeFeedSparklineValues(data []SparklinePoint) string {
	if len(data) == 0 {
		return ""
	}
	values := make([]string, len(data))
	for i, point := range data {
		values[i] = fmt.Sprint(point.Value)
	}
	return strings.Join(values, ",")
}

func buildFeedSparklineEnd(window sparklineWindow) string {
	if window.End.IsZero() {
		return ""
	}
	return window.End.Format("Jan 2006")
}

func monthlyPostBuckets(posts []*models.Post, window sparklineWindow, includePrivate bool) ([]int, []time.Time) {
	if window.Start.IsZero() || window.End.IsZero() || window.End.Before(window.Start) {
		return nil, nil
	}
	counts := map[string]int{}
	for _, post := range posts {
		if post == nil || post.Skip || post.Draft || !post.Published || post.Date == nil || (post.Private && !includePrivate) {
			continue
		}
		month := firstOfMonth(post.Date.UTC())
		if month.Before(window.Start) || month.After(window.End) {
			continue
		}
		monthKey := month.Format("2006-01")
		counts[monthKey]++
	}
	buckets := make([]int, 0, monthsBetweenInclusive(window.Start, window.End))
	months := make([]time.Time, 0, monthsBetweenInclusive(window.Start, window.End))
	for month := window.Start; !month.After(window.End); month = month.AddDate(0, 1, 0) {
		buckets = append(buckets, counts[month.Format("2006-01")])
		months = append(months, month)
	}
	return buckets, months
}

func computeSparklineWindow(posts []*models.Post, includePrivate bool) sparklineWindow {
	dates := make([]time.Time, 0, len(posts))
	for _, post := range posts {
		if post == nil || post.Skip || post.Draft || !post.Published || post.Date == nil || (post.Private && !includePrivate) {
			continue
		}
		date := post.Date.UTC()
		if !isSanePublishDate(date) {
			continue
		}
		dates = append(dates, firstOfMonth(date))
	}
	if len(dates) == 0 {
		return sparklineWindow{}
	}
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})
	start := dates[0]
	end := dates[len(dates)-1]
	if len(dates) >= 50 {
		start = dates[len(dates)/100]
		endIndex := len(dates) - 1 - (len(dates) / 100)
		if endIndex >= 0 && endIndex < len(dates) {
			end = dates[endIndex]
		}
		if end.Before(start) {
			start = dates[0]
			end = dates[len(dates)-1]
		}
	}
	return sparklineWindow{Start: start, End: end}
}

func isSanePublishDate(date time.Time) bool {
	if date.IsZero() {
		return false
	}
	if date.Before(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)) {
		return false
	}
	if date.After(time.Now().AddDate(1, 0, 0)) {
		return false
	}
	return true
}

func firstOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func monthsBetweenInclusive(start, end time.Time) int {
	if end.Before(start) {
		return 0
	}
	return (end.Year()-start.Year())*12 + int(end.Month()-start.Month()) + 1
}

func paginateFeedListings(feeds []FeedListingInfo, defaults models.FeedDefaults, firstPageURL, baseURL string) []FeedListingPage {
	if len(feeds) == 0 {
		return nil
	}

	itemsPerPage := defaults.ItemsPerPage
	if itemsPerPage <= 0 {
		itemsPerPage = 10
	}
	orphanThreshold := defaults.OrphanThreshold
	if orphanThreshold <= 0 {
		orphanThreshold = 3
	}

	pages := make([]FeedListingPage, 0)
	for i := 0; i < len(feeds); i += itemsPerPage {
		end := i + itemsPerPage
		if end > len(feeds) {
			end = len(feeds)
		}
		remaining := len(feeds) - end
		if remaining > 0 && remaining < orphanThreshold {
			end = len(feeds)
		}
		pageNum := len(pages) + 1
		pages = append(pages, FeedListingPage{
			Number:  pageNum,
			Feeds:   feeds[i:end],
			HasPrev: pageNum > 1,
		})
		if end >= len(feeds) {
			break
		}
	}

	totalPages := len(pages)
	pageURLs := make([]string, totalPages)
	for i := 0; i < totalPages; i++ {
		if i == 0 {
			pageURLs[i] = firstPageURL
		} else {
			pageURLs[i] = baseURL + "/page/" + fmt.Sprintf("%d", i+1) + "/"
		}
	}
	for i := range pages {
		pages[i].HasNext = i < totalPages-1
		pages[i].TotalPages = totalPages
		pages[i].TotalItems = len(feeds)
		pages[i].ItemsPerPage = itemsPerPage
		pages[i].PageURLs = pageURLs
		if pages[i].HasPrev {
			if i == 1 {
				pages[i].PrevURL = firstPageURL
			} else {
				pages[i].PrevURL = baseURL + "/page/" + fmt.Sprintf("%d", i) + "/"
			}
		}
		if pages[i].HasNext {
			pages[i].NextURL = baseURL + "/page/" + fmt.Sprintf("%d", i+2) + "/"
		}
	}

	return pages
}

func pathJoinURL(base, suffix string) string {
	if base == "/" {
		return "/" + suffix
	}
	return base + suffix
}

func (p *FeedsListingPlugin) renderFeedsPage(
	config *lifecycle.Config,
	feedsPage *models.FeedsPageConfig,
	sections []FeedListingSection,
	pageSlug string,
	title string,
	description string,
	pageLinks []FeedVariantLink,
	pagination *FeedListingPage,
) error {
	feedsDir := filepath.Join(config.OutputDir, filepath.FromSlash(pageSlug))
	if err := os.MkdirAll(feedsDir, 0o755); err != nil {
		return fmt.Errorf("creating feeds directory: %w", err)
	}

	engine, err := p.createTemplateEngine(config)
	if err != nil {
		return err
	}

	if !engine.TemplateExists(feedsPage.Template) {
		log.Printf("[feeds_listing] Warning: template %q not found, skipping feeds listing page", feedsPage.Template)
		return nil
	}

	modelsConfig := ToModelsConfig(config)
	pageRobots := feedsPage.Robots
	if modelsConfig != nil && modelsConfig.FeedsPage.Robots != "" {
		pageRobots = modelsConfig.FeedsPage.Robots
	}
	syntheticPost := &models.Post{
		Slug:        pageSlug,
		Title:       &title,
		Description: &description,
	}

	ctx := templates.NewContext(syntheticPost, "", modelsConfig)
	ctx.Extra["feeds_page"] = map[string]interface{}{
		"robots": pageRobots,
	}
	ctx.Extra["feeds_page_robots"] = pageRobots
	ctx.Extra["feed_sections"] = sections
	ctx.Extra["page_links"] = pageLinks
	ctx.Extra["pagination"] = pagination
	totalFeeds := 0
	for _, section := range sections {
		totalFeeds += section.TotalCount
	}
	ctx.Extra["total_feeds"] = totalFeeds

	html, err := engine.Render(feedsPage.Template, ctx)
	if err != nil {
		return fmt.Errorf("rendering feeds template: %w", err)
	}
	if pageRobots != "" && !strings.Contains(html, `name="robots"`) {
		html = strings.Replace(html, "</head>", `<meta name="robots" content="`+pageRobots+`">`+"\n</head>", 1)
	}

	outputPath := filepath.Join(feedsDir, "index.html")
	if err := os.WriteFile(outputPath, []byte(html), 0o644); err != nil { //nolint:gosec // web output must be world-readable
		return fmt.Errorf("writing feeds listing page: %w", err)
	}

	log.Printf("[feeds_listing] Generated /%s/ with %d feeds", pageSlug, totalFeeds)
	return nil
}

func (p *FeedsListingPlugin) createTemplateEngine(config *lifecycle.Config) (*templates.Engine, error) {
	templatesDir := PluginNameTemplates
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}

	themeName := getThemeName(config)
	cacheKey := templatesDir + ":" + themeName

	p.engineMu.RLock()
	if engine, ok := p.engineCache[cacheKey]; ok {
		p.engineMu.RUnlock()
		return engine, nil
	}
	p.engineMu.RUnlock()

	p.engineMu.Lock()
	defer p.engineMu.Unlock()
	if engine, ok := p.engineCache[cacheKey]; ok {
		return engine, nil
	}

	engine, err := templates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return nil, err
	}
	p.engineCache[cacheKey] = engine
	return engine, nil
}

var (
	_ lifecycle.Plugin         = (*FeedsListingPlugin)(nil)
	_ lifecycle.WritePlugin    = (*FeedsListingPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*FeedsListingPlugin)(nil)
)
