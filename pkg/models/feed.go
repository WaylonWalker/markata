package models

import "strings"

// PaginationType represents the type of pagination to use.
type PaginationType string

const (
	// PaginationManual uses traditional page links with full page reloads.
	PaginationManual PaginationType = "manual"

	// PaginationHTMX uses HTMX for seamless AJAX-based page loading.
	PaginationHTMX PaginationType = "htmx"

	// PaginationHTMXInfinite uses HTMX for infinite scroll pagination.
	PaginationHTMXInfinite PaginationType = "htmx-infinite"

	// PaginationJS uses client-side JavaScript for pagination.
	PaginationJS PaginationType = "js"
)

// FeedType represents the type of feed.
type FeedType string

const (
	// FeedTypeBlog is a standard blog feed sorted by date.
	FeedTypeBlog FeedType = "blog"

	// FeedTypeSeries is a series of related posts.
	FeedTypeSeries FeedType = "series"

	// FeedTypeGuide is a guide series with ordered content and prerequisites.
	FeedTypeGuide FeedType = "guide"
)

const (
	// FeedViewDefault is the feed's normal card/grid presentation.
	FeedViewDefault = "default"
	// FeedViewSimple is the compact simple-list presentation.
	FeedViewSimple = "simple"
	// FeedViewCalendar is the publishing-calendar presentation.
	FeedViewCalendar = "calendar"
	// FeedViewShorts is the immersive full-screen media presentation.
	FeedViewShorts = "shorts"
)

// DefaultFeedViews returns the presentation views enabled for feeds by default.
func DefaultFeedViews() []string {
	return []string{FeedViewDefault, FeedViewSimple, FeedViewCalendar}
}

// IsKnownFeedView reports whether view names a built-in feed presentation.
func IsKnownFeedView(view string) bool {
	switch view {
	case FeedViewDefault, FeedViewSimple, FeedViewCalendar, FeedViewShorts:
		return true
	default:
		return false
	}
}

// FeedConfig represents a feed configuration.
type FeedConfig struct {
	// Slug is the URL-safe identifier for the feed
	Slug string `json:"slug" yaml:"slug" toml:"slug"`

	// Title is the feed title
	Title string `json:"title" yaml:"title" toml:"title"`

	// Description is the feed description
	Description string `json:"description" yaml:"description" toml:"description"`

	// Type is the feed type (blog, series, guide)
	Type FeedType `json:"type" yaml:"type" toml:"type"`

	// Presentation selects the HTML feed layout (classic, auto, editorial, gallery, reading, timeline, resources).
	Presentation string `json:"presentation,omitempty" yaml:"presentation,omitempty" toml:"presentation,omitempty"`

	// Filter is the filter expression for selecting posts
	Filter string `json:"filter" yaml:"filter" toml:"filter"`

	// Sort is the field to sort posts by
	Sort string `json:"sort" yaml:"sort" toml:"sort"`

	// Reverse indicates if the sort order should be reversed
	Reverse bool `json:"reverse" yaml:"reverse" toml:"reverse"`

	// Views controls which HTML presentations are exposed for this feed.
	// Nil/unset inherits FeedDefaults.Views; an explicit empty list is invalid.
	Views []string `json:"views,omitempty" yaml:"views,omitempty" toml:"views,omitempty"`

	// ShortsPath optionally moves the immersive view to /shorts/ (or another route).
	ShortsPath string `json:"shorts_path,omitempty" yaml:"shorts_path,omitempty" toml:"shorts_path,omitempty"`

	// ItemsPerPage is the number of items per page (default: 10)
	ItemsPerPage int `json:"items_per_page" yaml:"items_per_page" toml:"items_per_page"`

	// OrphanThreshold is the minimum number of items for a separate page (default: 3)
	OrphanThreshold int `json:"orphan_threshold" yaml:"orphan_threshold" toml:"orphan_threshold"`

	// PaginationType specifies the pagination strategy (manual, htmx, js)
	PaginationType PaginationType `json:"pagination_type" yaml:"pagination_type" toml:"pagination_type"`

	// Limit is a hard cap on the number of posts in the feed (0 = unlimited)
	Limit int `json:"limit" yaml:"limit" toml:"limit"`

	// Offset skips the first N posts after filtering/sorting (0 = none)
	Offset int `json:"offset" yaml:"offset" toml:"offset"`

	// Formats specifies which output formats to generate
	Formats FeedFormats `json:"formats" yaml:"formats" toml:"formats"`

	// FormatOverrides preserves presence as well as values from parsed config.
	// An explicit false must not be mistaken for an omitted format.
	// Runtime-only: serializers still expose the original formats object.
	formatOverrides map[string]bool

	// Templates specifies custom templates for each format
	Templates FeedTemplates `json:"templates" yaml:"templates" toml:"templates"`

	// Sidebar controls whether this feed appears in the sidebar feed picker.
	// Default (nil/unset): included. Set to false to exclude from the picker.
	Sidebar *bool `json:"sidebar,omitempty" yaml:"sidebar,omitempty" toml:"sidebar,omitempty"`

	// SidebarTitle overrides the feed title in sidebar navigation
	SidebarTitle string `json:"sidebar_title" yaml:"sidebar_title" toml:"sidebar_title"`

	// SidebarOrder controls the position in multi-feed sidebars (lower = first, default: 0)
	SidebarOrder int `json:"sidebar_order" yaml:"sidebar_order" toml:"sidebar_order"`

	// SidebarGroupBy groups posts by a frontmatter field in the sidebar (e.g., "category")
	SidebarGroupBy string `json:"sidebar_group_by" yaml:"sidebar_group_by" toml:"sidebar_group_by"`

	// Primary marks this feed as part of the primary sidebar rotation on post pages.
	Primary bool `json:"primary,omitempty" yaml:"primary,omitempty" toml:"primary,omitempty"`

	// IncludePrivate allows including private posts in this feed (default: false)
	IncludePrivate bool `json:"include_private,omitempty" yaml:"include_private,omitempty" toml:"include_private,omitempty"`

	// Private marks this feed as private. It is retained as a compatibility
	// alias for site configurations that use `private = true`; private feeds
	// also include private posts.
	Private bool `json:"private,omitempty" yaml:"private,omitempty" toml:"private,omitempty"`

	// ArchiveDisabled disables the generated archive syndication endpoints for this feed.
	ArchiveDisabled bool `json:"archive_disabled,omitempty" yaml:"archive_disabled,omitempty" toml:"archive_disabled,omitempty"`

	// Robots controls the robots meta tag for HTML feed pages.
	// Example: "noindex,follow"
	Robots string `json:"robots,omitempty" yaml:"robots,omitempty" toml:"robots,omitempty"`

	// Posts holds the filtered posts at runtime (not serialized)
	Posts []*Post `json:"-" yaml:"-" toml:"-"`

	// Pages holds the paginated results at runtime (not serialized)
	Pages []FeedPage `json:"-" yaml:"-" toml:"-"`
}

// IncludesPrivate reports whether this feed explicitly opts into private posts.
// Private is retained as a compatibility alias for IncludePrivate.
func (f FeedConfig) IncludesPrivate() bool {
	return f.IncludePrivate || f.Private
}

// HasView reports whether this feed exposes a presentation view. Unresolved
// feed configs preserve the built-in three-view default for template callers.
func (f FeedConfig) HasView(view string) bool {
	views := f.Views
	if views == nil {
		views = DefaultFeedViews()
	}
	for _, candidate := range views {
		if candidate == view {
			return true
		}
	}
	return false
}

// ShortsURL returns the canonical path of this feed's full-screen viewer.
func (f FeedConfig) ShortsURL() string {
	if f.ShortsPath != "" {
		return f.ShortsPath
	}
	if strings.Trim(f.Slug, "/") == "" {
		return "/shorts/"
	}
	return "/" + strings.Trim(f.Slug, "/") + "/shorts/"
}

// GetSidebarTitle returns the effective title for sidebar navigation.
// Returns SidebarTitle if set, otherwise returns Title.
func (f *FeedConfig) GetSidebarTitle() string {
	if f.SidebarTitle != "" {
		return f.SidebarTitle
	}
	return f.Title
}

// FeedFormats specifies which output formats to generate for a feed.
type FeedFormats struct {
	// HTML generates an HTML page
	HTML bool `json:"html" yaml:"html" toml:"html"`

	// SimpleHTML generates a compact, dense list view of posts
	SimpleHTML bool `json:"simple_html" yaml:"simple_html" toml:"simple_html"`

	// RSS generates an RSS feed
	RSS bool `json:"rss" yaml:"rss" toml:"rss"`

	// Atom generates an Atom feed
	Atom bool `json:"atom" yaml:"atom" toml:"atom"`

	// JSON generates a JSON feed
	JSON bool `json:"json" yaml:"json" toml:"json"`

	// Markdown generates a Markdown file
	Markdown bool `json:"markdown" yaml:"markdown" toml:"markdown"`

	// Text generates a plain text file
	Text bool `json:"text" yaml:"text" toml:"text"`

	// Sitemap generates a sitemap XML file
	Sitemap bool `json:"sitemap" yaml:"sitemap" toml:"sitemap"`
}

// HasAnyEnabled returns true if any output format is enabled.
func (f FeedFormats) HasAnyEnabled() bool {
	return f.HTML || f.SimpleHTML || f.RSS || f.Atom || f.JSON || f.Markdown || f.Text || f.Sitemap
}

// FeedTemplates specifies custom templates for feed formats.
type FeedTemplates struct {
	// HTML is the template for HTML output
	HTML string `json:"html" yaml:"html" toml:"html"`

	// SimpleHTML is the template for simple HTML list output
	SimpleHTML string `json:"simple_html" yaml:"simple_html" toml:"simple_html"`

	// RSS is the template for RSS output
	RSS string `json:"rss" yaml:"rss" toml:"rss"`

	// Atom is the template for Atom output
	Atom string `json:"atom" yaml:"atom" toml:"atom"`

	// JSON is the template for JSON output
	JSON string `json:"json" yaml:"json" toml:"json"`

	// Card is the template for individual post cards
	Card string `json:"card" yaml:"card" toml:"card"`

	// Sitemap is the template for sitemap output
	Sitemap string `json:"sitemap" yaml:"sitemap" toml:"sitemap"`
}

// FeedPage represents a single page of paginated feed results.
type FeedPage struct {
	// Number is the page number (1-indexed)
	Number int `json:"number" yaml:"number" toml:"number"`

	// Posts is the list of posts on this page
	Posts []*Post `json:"posts" yaml:"posts" toml:"posts"`

	// HasPrev indicates if there is a previous page
	HasPrev bool `json:"has_prev" yaml:"has_prev" toml:"has_prev"`

	// HasNext indicates if there is a next page
	HasNext bool `json:"has_next" yaml:"has_next" toml:"has_next"`

	// PrevURL is the URL of the previous page
	PrevURL string `json:"prev_url" yaml:"prev_url" toml:"prev_url"`

	// NextURL is the URL of the next page
	NextURL string `json:"next_url" yaml:"next_url" toml:"next_url"`

	// TotalPages is the total number of pages in the feed
	TotalPages int `json:"total_pages" yaml:"total_pages" toml:"total_pages"`

	// TotalItems is the total number of posts in the feed
	TotalItems int `json:"total_items" yaml:"total_items" toml:"total_items"`

	// ItemsPerPage is the number of posts per page
	ItemsPerPage int `json:"items_per_page" yaml:"items_per_page" toml:"items_per_page"`

	// PageURLs contains URLs for all pages (for numbered pagination)
	PageURLs []string `json:"page_urls" yaml:"page_urls" toml:"page_urls"`

	// PaginationType is the pagination strategy used
	PaginationType PaginationType `json:"pagination_type" yaml:"pagination_type" toml:"pagination_type"`
}

// FeedDefaults provides default values that feeds inherit.
type FeedDefaults struct {
	// Presentation is the default HTML feed layout.
	Presentation string `json:"presentation,omitempty" yaml:"presentation,omitempty" toml:"presentation,omitempty"`
	// ItemsPerPage is the default number of items per page
	ItemsPerPage int `json:"items_per_page" yaml:"items_per_page" toml:"items_per_page"`

	// OrphanThreshold is the default minimum number of items for a separate page
	OrphanThreshold int `json:"orphan_threshold" yaml:"orphan_threshold" toml:"orphan_threshold"`

	// PaginationType is the default pagination strategy
	PaginationType PaginationType `json:"pagination_type" yaml:"pagination_type" toml:"pagination_type"`

	// Views specifies the default HTML presentations exposed by feeds.
	Views []string `json:"views,omitempty" yaml:"views,omitempty" toml:"views,omitempty"`

	// Formats specifies the default output formats
	Formats FeedFormats `json:"formats" yaml:"formats" toml:"formats"`

	// Templates specifies the default templates
	Templates FeedTemplates `json:"templates" yaml:"templates" toml:"templates"`

	// Syndication configures syndication feed behavior
	Syndication SyndicationConfig `json:"syndication" yaml:"syndication" toml:"syndication"`
}

// SyndicationConfig configures syndication feed behavior.
type SyndicationConfig struct {
	// MaxItems is the maximum number of items in syndication feeds
	MaxItems int `json:"max_items" yaml:"max_items" toml:"max_items"`

	// IncludeContent determines if full content is included in feeds
	IncludeContent bool `json:"include_content" yaml:"include_content" toml:"include_content"`

	// SiteArchiveDisabled disables the built-in site archive feeds.
	SiteArchiveDisabled bool `json:"site_archive_disabled,omitempty" yaml:"site_archive_disabled,omitempty" toml:"site_archive_disabled,omitempty"`

	// FeedArchivesDisabled disables generated archive syndication endpoints for non-root feeds.
	FeedArchivesDisabled bool `json:"feed_archives_disabled,omitempty" yaml:"feed_archives_disabled,omitempty" toml:"feed_archives_disabled,omitempty"`
}

// NewFeedDefaults creates FeedDefaults with sensible default values.
func NewFeedDefaults() FeedDefaults {
	return FeedDefaults{
		Presentation:    "classic",
		ItemsPerPage:    10,
		OrphanThreshold: 3,
		PaginationType:  PaginationManual,
		Views:           DefaultFeedViews(),
		Formats: FeedFormats{
			HTML:       true,
			SimpleHTML: true,
			RSS:        true,
			Atom:       true,
			JSON:       true,
			Sitemap:    true,
		},
		Templates: FeedTemplates{
			HTML:       "feed.html",
			SimpleHTML: "simple-feed.html",
			RSS:        "feed.xml",
			Atom:       "atom.xml",
			JSON:       "feed.json",
			Card:       "card.html",
			Sitemap:    "sitemap.xml",
		},
		Syndication: SyndicationConfig{
			MaxItems:             20,
			IncludeContent:       true,
			SiteArchiveDisabled:  false,
			FeedArchivesDisabled: false,
		},
	}
}

// NewFeedConfig creates a new FeedConfig with default values from FeedDefaults.
func NewFeedConfig(defaults FeedDefaults) *FeedConfig {
	views := defaults.Views
	if views == nil {
		views = DefaultFeedViews()
	}
	return &FeedConfig{
		Presentation:    defaults.Presentation,
		ItemsPerPage:    defaults.ItemsPerPage,
		OrphanThreshold: defaults.OrphanThreshold,
		PaginationType:  defaults.PaginationType,
		Views:           append([]string(nil), views...),
		Formats:         defaults.Formats,
		Templates:       defaults.Templates,
		Posts:           []*Post{},
		Pages:           []FeedPage{},
	}
}

// ApplyDefaults applies default values from FeedDefaults to a FeedConfig
// for any fields that are not explicitly set.
func (f *FeedConfig) ApplyDefaults(defaults FeedDefaults) {
	if f.Presentation == "" {
		f.Presentation = defaults.Presentation
	}
	if f.Private {
		f.IncludePrivate = true
	}
	if f.ItemsPerPage == 0 {
		f.ItemsPerPage = defaults.ItemsPerPage
	}
	if f.OrphanThreshold == 0 {
		f.OrphanThreshold = defaults.OrphanThreshold
	}
	if f.PaginationType == "" {
		f.PaginationType = defaults.PaginationType
	}
	if f.Views == nil {
		views := defaults.Views
		if views == nil {
			views = DefaultFeedViews()
		}
		f.Views = append([]string(nil), views...)
	}

	f.applyFormatDefaults(defaults.Formats)

	// Apply template defaults for any empty template paths
	if f.Templates.HTML == "" {
		f.Templates.HTML = defaults.Templates.HTML
	}
	if f.Templates.SimpleHTML == "" {
		f.Templates.SimpleHTML = defaults.Templates.SimpleHTML
	}
	if f.Templates.RSS == "" {
		f.Templates.RSS = defaults.Templates.RSS
	}
	if f.Templates.Atom == "" {
		f.Templates.Atom = defaults.Templates.Atom
	}
	if f.Templates.JSON == "" {
		f.Templates.JSON = defaults.Templates.JSON
	}
	if f.Templates.Card == "" {
		f.Templates.Card = defaults.Templates.Card
	}
	if f.Templates.Sitemap == "" {
		f.Templates.Sitemap = defaults.Templates.Sitemap
	}
}

// SetFormatOverrides retains explicit parsed values without exposing parser
// bookkeeping as a configuration field or changing config documentation.
func (f *FeedConfig) SetFormatOverrides(values map[string]bool) {
	f.formatOverrides = values
}

func (f *FeedConfig) applyFormatDefaults(defaults FeedFormats) {
	if len(f.formatOverrides) == 0 {
		if !f.Formats.HasAnyEnabled() {
			f.Formats = defaults
		}
		return
	}
	formats := defaults
	for name, enabled := range f.formatOverrides {
		switch name {
		case "html":
			formats.HTML = enabled
		case "simple_html":
			formats.SimpleHTML = enabled
		case "rss":
			formats.RSS = enabled
		case "atom":
			formats.Atom = enabled
		case "json":
			formats.JSON = enabled
		case "markdown":
			formats.Markdown = enabled
		case "text":
			formats.Text = enabled
		case "sitemap":
			formats.Sitemap = enabled
		}
	}
	f.Formats = formats
}

// Paginate divides the Posts slice into pages based on ItemsPerPage and OrphanThreshold.
func (f *FeedConfig) Paginate(baseURL string) {
	baseURL = strings.TrimSuffix(baseURL, "/")
	if len(f.Posts) == 0 {
		f.Pages = []FeedPage{}
		return
	}

	itemsPerPage := f.ItemsPerPage
	if itemsPerPage <= 0 {
		itemsPerPage = 10
	}

	// Determine pagination type (default to manual)
	paginationType := f.PaginationType
	if paginationType == "" {
		paginationType = PaginationManual
	}

	totalPosts := len(f.Posts)
	var pages []FeedPage

	for i := 0; i < totalPosts; {
		end := i + itemsPerPage
		if end > totalPosts {
			end = totalPosts
		}

		// Editorial feeds present entire calendar weeks as a single section.
		// Expand the page at its boundary instead of splitting the final week.
		// Other layouts keep exact item-count pagination.
		if f.Presentation == "editorial" && paginationType != PaginationJS {
			for end < totalPosts && sameFeedCalendarWeek(f.Posts[end-1], f.Posts[end]) {
				end++
			}
		}

		// Check orphan threshold: if remaining items are below threshold,
		// add them to the current page instead of creating a new page
		remaining := totalPosts - end
		if remaining > 0 && remaining < f.OrphanThreshold {
			end = totalPosts
		}

		pageNum := len(pages) + 1
		page := FeedPage{
			Number:  pageNum,
			Posts:   f.Posts[i:end],
			HasPrev: pageNum > 1,
		}

		pages = append(pages, page)

		// Advance by the actual boundary, which may include extra posts
		// from the trailing editorial week or orphan merge.
		i = end
	}

	totalPages := len(pages)

	// Generate page URLs for numbered navigation
	pageURLs := make([]string, totalPages)
	for i := 0; i < totalPages; i++ {
		if i == 0 {
			pageURLs[i] = baseURL + "/"
		} else {
			pageURLs[i] = baseURL + "/page/" + itoa(i+1) + "/"
		}
	}

	// Set HasNext, URLs, and metadata for each page
	for i := range pages {
		pages[i].HasNext = i < totalPages-1
		pages[i].TotalPages = totalPages
		pages[i].TotalItems = totalPosts
		pages[i].ItemsPerPage = itemsPerPage
		pages[i].PageURLs = pageURLs
		pages[i].PaginationType = paginationType

		if pages[i].HasPrev {
			if i == 1 {
				pages[i].PrevURL = baseURL + "/"
			} else {
				pages[i].PrevURL = baseURL + "/page/" + itoa(i) + "/"
			}
		}

		if pages[i].HasNext {
			pages[i].NextURL = baseURL + "/page/" + itoa(i+2) + "/"
		}
	}

	f.Pages = pages
}

// sameFeedCalendarWeek uses ISO week boundaries, including year transitions.
// Posts without dates stay together in the editorial "Undated" section.
func sameFeedCalendarWeek(left, right *Post) bool {
	if left == nil || right == nil {
		return false
	}
	if left.Date == nil || right.Date == nil {
		return left.Date == nil && right.Date == nil
	}
	leftYear, leftWeek := left.Date.ISOWeek()
	rightYear, rightWeek := right.Date.ISOWeek()
	return leftYear == rightYear && leftWeek == rightWeek
}

// itoa converts an integer to a string without importing strconv
func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var digits []byte
	negative := n < 0
	if negative {
		n = -n
	}

	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}

	if negative {
		digits = append([]byte{'-'}, digits...)
	}

	return string(digits)
}
