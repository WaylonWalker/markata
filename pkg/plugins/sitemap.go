// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// SitemapPlugin generates the root sitemap index and the primary pages sitemap.
type SitemapPlugin struct{}

// NewSitemapPlugin creates a new SitemapPlugin.
func NewSitemapPlugin() *SitemapPlugin { return &SitemapPlugin{} }

// Name returns the unique name of the plugin.
func (p *SitemapPlugin) Name() string { return "sitemap" }

// Priority returns the plugin priority for the given stage.
func (p *SitemapPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Write generates and writes sitemap files.
func (p *SitemapPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	siteURL := getSiteURL(config)
	if siteURL == "" {
		siteURL = DefaultSiteURL
	}

	pagesSitemap := p.buildPagesSitemap(m, siteURL)
	pagesOutput, err := xml.MarshalIndent(pagesSitemap, "", "    ")
	if err != nil {
		return fmt.Errorf("marshaling pages sitemap: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "sitemap-pages.xml"), []byte(xml.Header+string(pagesOutput)), 0o644); err != nil { //nolint:gosec // sitemap files must be world-readable
		return fmt.Errorf("writing pages sitemap: %w", err)
	}

	index := p.buildSitemapIndex(m, siteURL)
	indexOutput, err := xml.MarshalIndent(index, "", "    ")
	if err != nil {
		return fmt.Errorf("marshaling sitemap index: %w", err)
	}
	if err := os.WriteFile(filepath.Join(outputDir, "sitemap.xml"), []byte(xml.Header+string(indexOutput)), 0o644); err != nil { //nolint:gosec // sitemap files must be world-readable
		return fmt.Errorf("writing sitemap index: %w", err)
	}

	return nil
}

func (p *SitemapPlugin) buildPagesSitemap(m *lifecycle.Manager, siteURL string) *URLSet {
	sitemap := &URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]SitemapURL, 0),
	}

	posts := m.Posts()
	var latestDate time.Time
	for _, post := range posts {
		if post != nil && post.Published && !post.Draft && !post.Skip && !post.Private && !sitemapPostNoindex(post) && post.Date != nil && post.Date.After(latestDate) {
			latestDate = *post.Date
		}
	}

	homeURL := SitemapURL{Loc: siteURL + "/", ChangeFreq: "daily", Priority: "1.0"}
	if !latestDate.IsZero() {
		homeURL.LastMod = latestDate.Format("2006-01-02")
	}
	sitemap.URLs = append(sitemap.URLs, homeURL)

	for _, post := range posts {
		if post == nil || !post.Published || post.Draft || post.Skip || post.Private || sitemapPostNoindex(post) {
			continue
		}
		url := SitemapURL{Loc: siteURL + post.Href, ChangeFreq: "weekly", Priority: "0.8"}
		if post.Date != nil {
			url.LastMod = post.Date.Format("2006-01-02")
		}
		sitemap.URLs = append(sitemap.URLs, url)
	}

	sitemap.URLs = append(sitemap.URLs, p.generatedIndexPages(m, siteURL)...)
	sitemap.URLs = dedupeSitemapURLs(sitemap.URLs)

	return sitemap
}

// buildSitemap is retained for unit tests that verify the public URL set content.
func (p *SitemapPlugin) buildSitemap(m *lifecycle.Manager, siteURL string) *URLSet {
	sitemap := &URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]SitemapURL, 0),
	}

	posts := m.Posts()
	var latestDate time.Time
	for _, post := range posts {
		if post.Published && !post.Draft && !post.Skip && !post.Private && post.Date != nil && post.Date.After(latestDate) {
			latestDate = *post.Date
		}
	}

	homeURL := SitemapURL{Loc: siteURL + "/", ChangeFreq: "daily", Priority: "1.0"}
	if !latestDate.IsZero() {
		homeURL.LastMod = latestDate.Format("2006-01-02")
	}
	sitemap.URLs = append(sitemap.URLs, homeURL)

	for _, post := range posts {
		if !post.Published || post.Draft || post.Skip || post.Private {
			continue
		}
		url := SitemapURL{Loc: siteURL + post.Href, ChangeFreq: "weekly", Priority: "0.8"}
		if post.Date != nil {
			url.LastMod = post.Date.Format("2006-01-02")
		}
		sitemap.URLs = append(sitemap.URLs, url)
	}

	return sitemap
}

func (p *SitemapPlugin) generatedIndexPages(m *lifecycle.Manager, siteURL string) []SitemapURL {
	feedConfigs := getCachedFeedConfigs(m)
	urls := make([]SitemapURL, 0, len(feedConfigs)+2)

	for i := range feedConfigs {
		fc := &feedConfigs[i]
		if !eligibleFeedLandingForSitemap(fc) {
			continue
		}

		var latest time.Time
		for _, post := range fc.Posts {
			if post == nil || post.Private || post.Skip || post.Draft || !post.Published || post.Date == nil {
				continue
			}
			if post.Date.After(latest) {
				latest = *post.Date
			}
		}

		url := SitemapURL{Loc: siteURL + "/" + fc.Slug + "/", ChangeFreq: "weekly", Priority: "0.6"}
		if !latest.IsZero() {
			url.LastMod = latest.Format("2006-01-02")
		}
		urls = append(urls, url)
	}

	if tagsConfig, ok := m.Config().Extra["tags"].(models.TagsConfig); ok && tagsConfig.IsEnabled() {
		urls = append(urls, SitemapURL{Loc: siteURL + "/" + tagsConfig.SlugPrefix + "/", ChangeFreq: "weekly", Priority: "0.5"})
	}
	feedsPage := getFeedsPageConfig(m.Config())
	if feedsPage.IsEnabled() && !sitemapNoindex(feedsPage.Robots) {
		urls = append(urls, SitemapURL{Loc: siteURL + "/" + feedsPage.SlugPrefix + "/", ChangeFreq: "weekly", Priority: "0.5"})
		if hasGeneratedFeeds(feedConfigs, m.Config()) {
			urls = append(urls, SitemapURL{Loc: siteURL + "/" + feedsPage.SlugPrefix + "/generated/", ChangeFreq: "weekly", Priority: "0.4"})
		}
	}

	return urls
}

func eligibleFeedLandingForSitemap(feed *models.FeedConfig) bool {
	return feed != nil && feed.Slug != "" && feed.Formats.HTML && feed.Formats.Sitemap &&
		!feed.IncludesPrivate() && !sitemapNoindex(feed.Robots)
}

// sitemapNoindex accepts a comma- or whitespace-separated robots directive.
// Do not confuse "noindex" with "index" or "noimageindex".
func sitemapNoindex(value string) bool {
	for _, token := range strings.FieldsFunc(strings.ToLower(value), func(r rune) bool {
		return r == ',' || r == ';' || r == ' ' || r == '\t' || r == '\n'
	}) {
		if token == "noindex" {
			return true
		}
	}
	return false
}

func sitemapPostNoindex(post *models.Post) bool {
	if post == nil || post.Extra == nil {
		return false
	}
	robots, ok := post.Extra["robots"].(string)
	return ok && sitemapNoindex(robots)
}

func dedupeSitemapURLs(urls []SitemapURL) []SitemapURL {
	seen := make(map[string]struct{}, len(urls))
	unique := urls[:0]
	for _, item := range urls {
		if _, present := seen[item.Loc]; present {
			continue
		}
		seen[item.Loc] = struct{}{}
		unique = append(unique, item)
	}
	return unique
}

func hasGeneratedFeeds(feedConfigs []models.FeedConfig, config *lifecycle.Config) bool {
	configured := configuredFeedSlugs(config)
	count := 0
	for i := range feedConfigs {
		fc := &feedConfigs[i]
		if fc.IncludesPrivate() {
			continue
		}
		if _, ok := configured[fc.Slug]; !ok {
			count++
		}
	}
	defaults := getFeedDefaults(config)
	pageSize := defaults.ItemsPerPage
	if pageSize <= 0 {
		pageSize = 10
	}
	return count > pageSize
}

func (p *SitemapPlugin) buildSitemapIndex(m *lifecycle.Manager, siteURL string) *SitemapIndex {
	// Use the actual public pages sitemap, including generated landing pages.
	// Unpublished, private, skipped, and noindex posts cannot affect its lastmod.
	pagesLatest := sitemapURLsLastMod(p.buildPagesSitemap(m, siteURL).URLs)
	pagesLoc := siteURL + "/sitemap-pages.xml"
	entries := []SitemapIndexEntry{{Loc: pagesLoc, LastMod: pagesLatest}}
	seen := map[string]bool{pagesLoc: true}
	feedConfigs := getCachedFeedConfigs(m)

	for i := range feedConfigs {
		fc := &feedConfigs[i]
		if !fc.Formats.Sitemap || fc.Slug == "" || fc.IncludesPrivate() || sitemapNoindex(fc.Robots) {
			continue
		}
		loc := siteURL + "/" + fc.Slug + "/sitemap.xml"
		if seen[loc] {
			continue
		}
		seen[loc] = true
		entries = append(entries, SitemapIndexEntry{
			Loc:     loc,
			LastMod: sitemapLastMod(fc.Posts),
		})
	}

	return &SitemapIndex{
		XMLNS:    "http://www.sitemaps.org/schemas/sitemap/0.9",
		Sitemaps: entries,
	}
}

// URLSet represents a sitemap urlset document.
type URLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	XMLNS   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapURL represents a single URL entry in a sitemap.
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// SitemapIndex represents a sitemap index document.
type SitemapIndex struct {
	XMLName  xml.Name            `xml:"sitemapindex"`
	XMLNS    string              `xml:"xmlns,attr"`
	Sitemaps []SitemapIndexEntry `xml:"sitemap"`
}

// SitemapIndexEntry represents a sitemap reference inside a sitemap index.
type SitemapIndexEntry struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

// sitemapURLsLastMod derives index metadata from URLs actually emitted in
// sitemap-pages.xml, rather than from posts omitted by sitemap eligibility.
func sitemapURLsLastMod(urls []SitemapURL) string {
	latest := ""
	for _, url := range urls {
		// Sitemap dates use fixed-width YYYY-MM-DD, so lexical order is date order.
		if url.LastMod > latest {
			latest = url.LastMod
		}
	}
	return latest
}

// sitemapLastMod mirrors publishSitemap's post eligibility. An empty feed
// sitemap has no defensible modification date, even if it contains scheduled
// or private source posts.
func sitemapLastMod(posts []*models.Post) string {
	var latest time.Time
	for _, post := range posts {
		if post == nil || !post.Published || post.Draft || post.Skip || post.Private || post.Date == nil {
			continue
		}
		if post.Date.After(latest) {
			latest = *post.Date
		}
	}
	if latest.IsZero() {
		return ""
	}
	return latest.Format("2006-01-02")
}

var (
	_ lifecycle.Plugin         = (*SitemapPlugin)(nil)
	_ lifecycle.WritePlugin    = (*SitemapPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*SitemapPlugin)(nil)
)
