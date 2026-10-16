package models

import (
	"testing"
	"time"
)

// =============================================================================
// Feed Tests based on tests.yaml
// =============================================================================

func TestFeed_BasicCreation(t *testing.T) {
	// Test case: "basic feed"
	defaults := NewFeedDefaults()
	feed := NewFeedConfig(defaults)

	feed.Slug = "blog"
	feed.Filter = "True"

	posts := []*Post{
		{Title: strPtr("Post 1"), Slug: "post-1"},
		{Title: strPtr("Post 2"), Slug: "post-2"},
	}
	feed.Posts = posts

	// Verify feed properties
	if feed.Slug != "blog" {
		t.Errorf("slug: got %q, want 'blog'", feed.Slug)
	}
	if len(feed.Posts) != 2 {
		t.Errorf("posts_count: got %d, want 2", len(feed.Posts))
	}
}

func TestFeed_FilteredFeed(t *testing.T) {
	// Test case: "filtered feed"
	// This tests the concept - actual filtering is done by the filter package
	feed := &FeedConfig{
		Slug:   "published",
		Filter: "published == True",
	}

	// Simulate filtered posts
	feed.Posts = []*Post{
		{Title: strPtr("Public"), Published: true},
	}

	if feed.Filter != "published == True" {
		t.Errorf("filter: got %q, want 'published == True'", feed.Filter)
	}
	if len(feed.Posts) != 1 {
		t.Errorf("posts_count: got %d, want 1", len(feed.Posts))
	}
}

func TestFeed_SortedFeed(t *testing.T) {
	// Test case: "sorted feed"
	feed := &FeedConfig{
		Slug:    "recent",
		Filter:  "True",
		Sort:    "date",
		Reverse: true,
	}

	// Verify sort configuration
	if feed.Sort != "date" {
		t.Errorf("sort: got %q, want 'date'", feed.Sort)
	}
	if !feed.Reverse {
		t.Error("reverse: expected true")
	}
}

func TestFeed_Pagination(t *testing.T) {
	// Test case: "paginated feed"
	feed := &FeedConfig{
		Slug:            "archive",
		Filter:          "True",
		ItemsPerPage:    2,
		OrphanThreshold: 2, // Changed threshold so orphan logic doesn't trigger
	}

	// Create 5 posts
	posts := make([]*Post, 5)
	for i := 0; i < 5; i++ {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title, Slug: "post-" + itoa(i+1)}
	}
	feed.Posts = posts

	// Paginate
	feed.Paginate("/archive")

	// With 5 posts, 2 per page, and orphan threshold of 2:
	// Page 1: 2 posts, remaining = 3, 3 >= 2 (threshold), continue
	// Page 2: 2 posts, remaining = 1, 1 < 2 (threshold), merge
	// Result: Page 1 (2 posts), Page 2 (3 posts) = 2 pages
	if len(feed.Pages) != 2 {
		t.Errorf("total_pages: got %d, want 2", len(feed.Pages))
	}

	// First page should have 2 posts
	if len(feed.Pages) > 0 && len(feed.Pages[0].Posts) != 2 {
		t.Errorf("page_1_count: got %d, want 2", len(feed.Pages[0].Posts))
	}

	// Last page should have 3 posts (merged)
	if len(feed.Pages) > 1 && len(feed.Pages[1].Posts) != 3 {
		t.Errorf("page_2_count: got %d, want 3", len(feed.Pages[1].Posts))
	}
}

// Editorial pages are allowed to grow to finish the week; no week's
// photo gallery or featured story should be cut between adjacent pages.
func TestFeed_EditorialPaginationPreservesWeeks(t *testing.T) {
	feed := &FeedConfig{Slug: "archive", Presentation: "editorial", ItemsPerPage: 3, OrphanThreshold: 2}
	// Oct 12–18, Oct 5–11, and Sep 28–Oct 4. Descending date order.
	for _, day := range []int{16, 15, 14, 13, 12, 9, 8, 7, 6, 5, 2, 1} {
		date := time.Date(2026, time.October, day, 12, 0, 0, 0, time.UTC)
		feed.Posts = append(feed.Posts, &Post{Date: &date})
	}
	feed.Paginate("/archive")
	if got := len(feed.Pages); got != 3 {
		t.Fatalf("got %d pages, want 3", got)
	}
	want := []int{5, 5, 2}
	for i, page := range feed.Pages {
		if len(page.Posts) != want[i] {
			t.Errorf("page %d: %d posts, want %d", i+1, len(page.Posts), want[i])
		}
		if page.TotalItems != 12 || page.TotalPages != 3 {
			t.Errorf("page %d metadata: %d total posts, %d pages", i+1, page.TotalItems, page.TotalPages)
		}
	}
	if feed.Pages[0].NextURL != "/archive/page/2/" || feed.Pages[1].PrevURL != "/archive/" || feed.Pages[2].PrevURL != "/archive/page/2/" {
		t.Fatal("editorial page navigation URLs changed")
	}
}

func TestFeed_EditorialPaginationYearBoundaryAndUndated(t *testing.T) {
	feed := &FeedConfig{Presentation: "editorial", ItemsPerPage: 2}
	for _, day := range []time.Time{
		time.Date(2027, 1, 3, 12, 0, 0, 0, time.UTC),
		time.Date(2027, 1, 2, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 31, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 28, 12, 0, 0, 0, time.UTC),
		time.Date(2026, 12, 27, 12, 0, 0, 0, time.UTC),
	} {
		date := day
		feed.Posts = append(feed.Posts, &Post{Date: &date})
	}
	feed.Posts = append(feed.Posts, &Post{}, &Post{})
	feed.Paginate("/archive")
	if len(feed.Pages) != 2 || len(feed.Pages[0].Posts) != 4 || len(feed.Pages[1].Posts) != 3 {
		t.Fatalf("unexpected ISO year and undated boundaries: %+v", feed.Pages)
	}
}

func TestFeed_NonEditorialPaginationKeepsFixedPageSize(t *testing.T) {
	feed := &FeedConfig{Presentation: "classic", ItemsPerPage: 2}
	for i := 0; i < 5; i++ {
		date := time.Date(2026, 10, 12+i, 12, 0, 0, 0, time.UTC)
		feed.Posts = append(feed.Posts, &Post{Date: &date})
	}
	feed.Paginate("/archive")
	if len(feed.Pages) != 3 || len(feed.Pages[0].Posts) != 2 {
		t.Fatalf("classic pagination should be unchanged: %+v", feed.Pages)
	}
}

func TestFeed_OrphanThreshold(t *testing.T) {
	// Test case: "feed with orphan threshold"
	// Last page with items < threshold merges with previous
	feed := &FeedConfig{
		Slug:            "test",
		Filter:          "True",
		ItemsPerPage:    3,
		OrphanThreshold: 2,
	}

	// Create 4 posts
	posts := make([]*Post, 4)
	for i := 0; i < 4; i++ {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title, Slug: "post-" + itoa(i+1)}
	}
	feed.Posts = posts

	// Paginate
	feed.Paginate("/test")

	// With 4 posts, 3 per page, and orphan threshold of 2:
	// Page 1 would have 3, Page 2 would have 1
	// Since 1 < 2 (threshold), they merge into 1 page with 4 posts
	if len(feed.Pages) != 1 {
		t.Errorf("total_pages: got %d, want 1", len(feed.Pages))
	}
	if len(feed.Pages[0].Posts) != 4 {
		t.Errorf("page_1_count: got %d, want 4", len(feed.Pages[0].Posts))
	}
}

func TestFeed_HomePageFeed(t *testing.T) {
	// Test case: "home page feed (empty slug)"
	feed := &FeedConfig{
		Slug:         "",
		Filter:       "published == True",
		ItemsPerPage: 5,
	}

	posts := []*Post{
		{Title: strPtr("Post 1"), Published: true},
	}
	feed.Posts = posts

	// Verify empty slug for home page
	if feed.Slug != "" {
		t.Errorf("slug: got %q, want empty string", feed.Slug)
	}

	// Paginate with empty base URL for home
	feed.Paginate("")

	// First page URL should be root
	if len(feed.Pages) > 0 && feed.Pages[0].Number != 1 {
		t.Errorf("page number: got %d, want 1", feed.Pages[0].Number)
	}
}

// =============================================================================
// FeedConfig Tests
// =============================================================================

func TestFeedConfig_Paginate_RootBaseURLAvoidsProtocolRelativeLinks(t *testing.T) {
	feed := &FeedConfig{
		ItemsPerPage: 2,
		Posts: []*Post{
			{Slug: "one"},
			{Slug: "two"},
			{Slug: "three"},
			{Slug: "four"},
			{Slug: "five"},
		},
	}

	feed.Paginate("/")

	if got, want := len(feed.Pages), 3; got != want {
		t.Fatalf("pages = %d, want %d", got, want)
	}

	wantURLs := []string{"/", "/page/2/", "/page/3/"}
	for i, want := range wantURLs {
		if got := feed.Pages[0].PageURLs[i]; got != want {
			t.Errorf("PageURLs[%d] = %q, want %q", i, got, want)
		}
	}
	if got, want := feed.Pages[0].NextURL, "/page/2/"; got != want {
		t.Errorf("first NextURL = %q, want %q", got, want)
	}
	if got, want := feed.Pages[1].PrevURL, "/"; got != want {
		t.Errorf("second PrevURL = %q, want %q", got, want)
	}
	if got, want := feed.Pages[1].NextURL, "/page/3/"; got != want {
		t.Errorf("second NextURL = %q, want %q", got, want)
	}
	if got, want := feed.Pages[2].PrevURL, "/page/2/"; got != want {
		t.Errorf("third PrevURL = %q, want %q", got, want)
	}
}

func TestFeedConfig_ApplyDefaults(t *testing.T) {
	defaults := FeedDefaults{
		ItemsPerPage:    15,
		OrphanThreshold: 5,
		Formats: FeedFormats{
			HTML: true,
			RSS:  true,
			Atom: true,
		},
	}

	feed := &FeedConfig{
		Slug: "test",
	}

	feed.ApplyDefaults(defaults)

	if feed.ItemsPerPage != 15 {
		t.Errorf("items_per_page: got %d, want 15", feed.ItemsPerPage)
	}
	if feed.OrphanThreshold != 5 {
		t.Errorf("orphan_threshold: got %d, want 5", feed.OrphanThreshold)
	}
}

func TestFeedConfig_PrivateAliasIncludesPrivate(t *testing.T) {
	feed := &FeedConfig{Private: true}
	if !feed.IncludesPrivate() {
		t.Fatal("private alias should opt the feed into private posts")
	}
	feed.ApplyDefaults(FeedDefaults{})
	if !feed.IncludePrivate {
		t.Fatal("ApplyDefaults should normalize the private alias")
	}
}

func TestFeedConfig_DoesNotOverrideExplicitValues(t *testing.T) {
	defaults := FeedDefaults{
		ItemsPerPage:    10,
		OrphanThreshold: 3,
	}

	feed := &FeedConfig{
		Slug:         "test",
		ItemsPerPage: 25,
	}

	feed.ApplyDefaults(defaults)

	// Explicit value should be preserved
	if feed.ItemsPerPage != 25 {
		t.Errorf("items_per_page: got %d, want 25", feed.ItemsPerPage)
	}
	// Default should be applied for unset values
	if feed.OrphanThreshold != 3 {
		t.Errorf("orphan_threshold: got %d, want 3", feed.OrphanThreshold)
	}
}

func TestNewFeedDefaults(t *testing.T) {
	defaults := NewFeedDefaults()

	if defaults.ItemsPerPage != 10 {
		t.Errorf("items_per_page: got %d, want 10", defaults.ItemsPerPage)
	}
	if defaults.OrphanThreshold != 3 {
		t.Errorf("orphan_threshold: got %d, want 3", defaults.OrphanThreshold)
	}
	if !defaults.Formats.HTML {
		t.Error("html format should be enabled by default")
	}
	if !defaults.Formats.RSS {
		t.Error("rss format should be enabled by default")
	}
	if !defaults.Formats.Atom {
		t.Error("atom format should be enabled by default")
	}
	if !defaults.Formats.JSON {
		t.Error("json format should be enabled by default")
	}
	if !defaults.Formats.Sitemap {
		t.Error("sitemap format should be enabled by default")
	}
}

func TestNewFeedConfig(t *testing.T) {
	defaults := NewFeedDefaults()
	feed := NewFeedConfig(defaults)

	if feed.ItemsPerPage != defaults.ItemsPerPage {
		t.Errorf("items_per_page: got %d, want %d", feed.ItemsPerPage, defaults.ItemsPerPage)
	}
	if feed.Posts == nil {
		t.Error("Posts should be initialized")
	}
	if feed.Pages == nil {
		t.Error("Pages should be initialized")
	}
}

// =============================================================================
// Pagination Tests
// =============================================================================

func TestFeedConfig_Paginate_EmptyPosts(t *testing.T) {
	feed := &FeedConfig{
		ItemsPerPage: 10,
		Posts:        []*Post{},
	}

	feed.Paginate("/blog")

	if len(feed.Pages) != 0 {
		t.Errorf("expected 0 pages for empty posts, got %d", len(feed.Pages))
	}
}

func TestFeedConfig_Paginate_SinglePage(t *testing.T) {
	posts := make([]*Post, 5)
	for i := range posts {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title}
	}

	feed := &FeedConfig{
		ItemsPerPage: 10,
		Posts:        posts,
	}

	feed.Paginate("/blog")

	if len(feed.Pages) != 1 {
		t.Errorf("expected 1 page, got %d", len(feed.Pages))
	}
	if feed.Pages[0].HasPrev {
		t.Error("first page should not have prev")
	}
	if feed.Pages[0].HasNext {
		t.Error("single page should not have next")
	}
}

func TestFeedConfig_Paginate_MultiplePages(t *testing.T) {
	posts := make([]*Post, 25)
	for i := range posts {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title}
	}

	feed := &FeedConfig{
		ItemsPerPage:    10,
		OrphanThreshold: 3,
		Posts:           posts,
	}

	feed.Paginate("/blog")

	// 25 posts with 10 per page = 3 pages (10, 10, 5)
	// 5 is above orphan threshold of 3, so no merging
	if len(feed.Pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(feed.Pages))
	}

	// Check page navigation
	if feed.Pages[0].HasPrev {
		t.Error("first page should not have prev")
	}
	if !feed.Pages[0].HasNext {
		t.Error("first page should have next")
	}
	if !feed.Pages[1].HasPrev {
		t.Error("middle page should have prev")
	}
	if !feed.Pages[1].HasNext {
		t.Error("middle page should have next")
	}
	if !feed.Pages[2].HasPrev {
		t.Error("last page should have prev")
	}
	if feed.Pages[2].HasNext {
		t.Error("last page should not have next")
	}
}

func TestFeedConfig_Paginate_URLs(t *testing.T) {
	posts := make([]*Post, 15)
	for i := range posts {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title}
	}

	feed := &FeedConfig{
		ItemsPerPage: 5,
		Posts:        posts,
	}

	feed.Paginate("/blog")

	// Check URLs
	// Page 1: no prev, next is /blog/page/2/
	if feed.Pages[0].PrevURL != "" {
		t.Errorf("page 1 prev URL should be empty, got %q", feed.Pages[0].PrevURL)
	}
	if feed.Pages[0].NextURL != "/blog/page/2/" {
		t.Errorf("page 1 next URL: got %q, want '/blog/page/2/'", feed.Pages[0].NextURL)
	}

	// Page 2: prev is /blog/, next is /blog/page/3/
	if feed.Pages[1].PrevURL != "/blog/" {
		t.Errorf("page 2 prev URL: got %q, want '/blog/'", feed.Pages[1].PrevURL)
	}
	if feed.Pages[1].NextURL != "/blog/page/3/" {
		t.Errorf("page 2 next URL: got %q, want '/blog/page/3/'", feed.Pages[1].NextURL)
	}

	// Page 3: prev is /blog/page/2/, no next
	if feed.Pages[2].PrevURL != "/blog/page/2/" {
		t.Errorf("page 3 prev URL: got %q, want '/blog/page/2/'", feed.Pages[2].PrevURL)
	}
	if feed.Pages[2].NextURL != "" {
		t.Errorf("page 3 next URL should be empty, got %q", feed.Pages[2].NextURL)
	}
}

func TestFeedConfig_Paginate_ZeroItemsPerPage(t *testing.T) {
	posts := make([]*Post, 5)
	for i := range posts {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title}
	}

	feed := &FeedConfig{
		ItemsPerPage: 0, // Should default to 10
		Posts:        posts,
	}

	feed.Paginate("/blog")

	// Should use default of 10, so all 5 posts on 1 page
	if len(feed.Pages) != 1 {
		t.Errorf("expected 1 page with default items_per_page, got %d", len(feed.Pages))
	}
}

func TestFeedConfig_Paginate_PaginationType(t *testing.T) {
	posts := make([]*Post, 15)
	for i := range posts {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title}
	}

	tests := []struct {
		name           string
		paginationType PaginationType
		wantType       PaginationType
	}{
		{"manual explicit", PaginationManual, PaginationManual},
		{"htmx explicit", PaginationHTMX, PaginationHTMX},
		{"htmx-infinite explicit", PaginationHTMXInfinite, PaginationHTMXInfinite},
		{"js explicit", PaginationJS, PaginationJS},
		{"empty defaults to manual", "", PaginationManual},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			feed := &FeedConfig{
				ItemsPerPage:   5,
				PaginationType: tt.paginationType,
				Posts:          posts,
			}

			feed.Paginate("/blog")

			for _, page := range feed.Pages {
				if page.PaginationType != tt.wantType {
					t.Errorf("page %d PaginationType: got %q, want %q", page.Number, page.PaginationType, tt.wantType)
				}
			}
		})
	}
}

func TestFeedConfig_Paginate_Metadata(t *testing.T) {
	posts := make([]*Post, 15)
	for i := range posts {
		title := "Post " + itoa(i+1)
		posts[i] = &Post{Title: &title}
	}

	feed := &FeedConfig{
		ItemsPerPage: 5,
		Posts:        posts,
	}

	feed.Paginate("/blog")

	// Should have 3 pages
	if len(feed.Pages) != 3 {
		t.Fatalf("expected 3 pages, got %d", len(feed.Pages))
	}

	// Check metadata on all pages
	for i, page := range feed.Pages {
		if page.TotalPages != 3 {
			t.Errorf("page %d TotalPages: got %d, want 3", i+1, page.TotalPages)
		}
		if page.TotalItems != 15 {
			t.Errorf("page %d TotalItems: got %d, want 15", i+1, page.TotalItems)
		}
		if page.ItemsPerPage != 5 {
			t.Errorf("page %d ItemsPerPage: got %d, want 5", i+1, page.ItemsPerPage)
		}
		if len(page.PageURLs) != 3 {
			t.Errorf("page %d PageURLs: got %d URLs, want 3", i+1, len(page.PageURLs))
		}
	}

	// Check PageURLs content
	expectedURLs := []string{"/blog/", "/blog/page/2/", "/blog/page/3/"}
	for i, url := range feed.Pages[0].PageURLs {
		if url != expectedURLs[i] {
			t.Errorf("PageURLs[%d]: got %q, want %q", i, url, expectedURLs[i])
		}
	}
}

// =============================================================================
// Feed Formats Tests
// =============================================================================

func TestFeedFormats(t *testing.T) {
	formats := FeedFormats{
		HTML:     true,
		RSS:      true,
		Atom:     false,
		JSON:     true,
		Markdown: false,
		Text:     false,
	}

	if !formats.HTML {
		t.Error("HTML should be enabled")
	}
	if !formats.RSS {
		t.Error("RSS should be enabled")
	}
	if formats.Atom {
		t.Error("Atom should be disabled")
	}
	if !formats.JSON {
		t.Error("JSON should be enabled")
	}
	if formats.Markdown {
		t.Error("Markdown should be disabled")
	}
	if formats.Text {
		t.Error("Text should be disabled")
	}
}

// =============================================================================
// Feed Templates Tests
// =============================================================================

func TestFeedTemplates(t *testing.T) {
	templates := FeedTemplates{
		HTML: "custom-feed.html",
		RSS:  "custom-rss.xml",
		Atom: "custom-atom.xml",
		JSON: "custom-feed.json",
		Card: "custom-card.html",
	}

	if templates.HTML != "custom-feed.html" {
		t.Errorf("HTML template: got %q, want 'custom-feed.html'", templates.HTML)
	}
	if templates.RSS != "custom-rss.xml" {
		t.Errorf("RSS template: got %q, want 'custom-rss.xml'", templates.RSS)
	}
}

// =============================================================================
// Syndication Config Tests
// =============================================================================

func TestSyndicationConfig(t *testing.T) {
	config := SyndicationConfig{
		MaxItems:       20,
		IncludeContent: true,
	}

	if config.MaxItems != 20 {
		t.Errorf("max_items: got %d, want 20", config.MaxItems)
	}
	if !config.IncludeContent {
		t.Error("include_content should be true")
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

func strPtr(s string) *string {
	return &s
}

func TestFeedConfigApplyDefaultsPreservesExplicitFalse(t *testing.T) {
	defaults := NewFeedDefaults()
	for _, tc := range []struct {
		name      string
		overrides map[string]bool
		wantHTML  bool
		wantRSS   bool
		wantSite  bool
	}{
		{"sitemap off", map[string]bool{"sitemap": false}, true, true, false},
		{"all three off", map[string]bool{"html": false, "rss": false, "sitemap": false}, false, false, false},
		{"rss only on", map[string]bool{"html": false, "rss": true, "sitemap": false}, false, true, false},
		{"inherit defaults", nil, true, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			feed := FeedConfig{}
			feed.SetFormatOverrides(tc.overrides)
			feed.ApplyDefaults(defaults)
			if feed.Formats.HTML != tc.wantHTML || feed.Formats.RSS != tc.wantRSS || feed.Formats.Sitemap != tc.wantSite {
				t.Fatalf("formats = %+v, expected html=%v rss=%v sitemap=%v", feed.Formats, tc.wantHTML, tc.wantRSS, tc.wantSite)
			}
		})
	}
}

