// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"net/url"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"golang.org/x/net/html"
)

// LinkCollectorPlugin collects all hyperlinks from posts and tracks inlinks
// (pages linking TO a post) and outlinks (pages a post links TO).
// It runs in the render stage after render_markdown.
type LinkCollectorPlugin struct {
	// includeFeeds includes feed pages in inlinks when true
	includeFeeds bool

	// includeIndex includes index page in inlinks when true
	includeIndex bool

	// siteURL is the base URL of the site
	siteURL string

	// siteDomain is the domain extracted from siteURL
	siteDomain string
}

// NewLinkCollectorPlugin creates a new LinkCollectorPlugin with default settings.
func NewLinkCollectorPlugin() *LinkCollectorPlugin {
	return &LinkCollectorPlugin{
		includeFeeds: false,
		includeIndex: false,
	}
}

// Name returns the unique name of the plugin.
func (p *LinkCollectorPlugin) Name() string {
	return "link_collector"
}

// Configure reads configuration options for the plugin.
func (p *LinkCollectorPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		// Check for link_collector config section
		if linkCollectorConfig, ok := config.Extra["link_collector"].(map[string]interface{}); ok {
			if includeFeeds, ok := linkCollectorConfig["include_feeds"].(bool); ok {
				p.includeFeeds = includeFeeds
			}
			if includeIndex, ok := linkCollectorConfig["include_index"].(bool); ok {
				p.includeIndex = includeIndex
			}
		}

		// Get site URL from config
		if siteURL, ok := config.Extra["url"].(string); ok {
			p.siteURL = siteURL
			if parsed, err := url.Parse(siteURL); err == nil {
				p.siteDomain = parsed.Host
			}
		}
	}
	return nil
}

// Priority returns the plugin's priority for the render stage.
// Returns a late priority to ensure it runs after render_markdown.
func (p *LinkCollectorPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Render collects all links from posts and populates inlinks/outlinks.
func (p *LinkCollectorPlugin) Render(m *lifecycle.Manager) error {
	posts := m.Posts()
	if len(posts) == 0 {
		return nil
	}

	cache := GetBuildCache(m)
	useCache := cache != nil

	// Use the shared PostIndex from the lifecycle manager
	idx := m.PostIndex()

	// Collect all links from all posts
	var allLinks []*models.Link
	var linksMu sync.Mutex

	postsToProcess := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && post.ArticleHTML != ""
	})

	if lifecycle.IsServeIncremental(m) {
		if affected := lifecycle.GetServeAffectedPaths(m); len(affected) > 0 {
			filtered := postsToProcess[:0]
			for _, post := range postsToProcess {
				if affected[post.Path] {
					filtered = append(filtered, post)
				}
			}
			postsToProcess = filtered
		}
	}

	// Process each post to extract hrefs and create Link objects
	err := m.ProcessPostsSliceConcurrently(postsToProcess, func(post *models.Post) error {
		// Build base URL for this post
		baseURL := p.buildBaseURL(post)

		var hrefs []string
		var hrefTextMap map[string]string
		if useCache {
			articleHash := buildcache.ContentHash(post.ArticleHTML)
			if cached := cache.GetCachedLinkHrefs(post.Path, articleHash); cached != nil {
				hrefs = cached
				// When restoring from cache we only have hrefs; build text map from HTML
				hrefTextMap = extractHrefTextMap(post.ArticleHTML)
			} else {
				hrefs, hrefTextMap = extractHrefsAndText(post.ArticleHTML)
				cache.CacheLinkHrefs(post.Path, articleHash, hrefs)
			}
		} else {
			hrefs, hrefTextMap = extractHrefsAndText(post.ArticleHTML)
		}
		post.Hrefs = hrefs

		// Create Link objects for each href
		var postLinks []*models.Link
		for _, href := range hrefs {
			link := p.createLink(post, baseURL, href, idx, hrefTextMap)
			if link != nil {
				postLinks = append(postLinks, link)
			}
		}

		// Add links to global collection (thread-safe)
		linksMu.Lock()
		allLinks = append(allLinks, postLinks...)
		linksMu.Unlock()

		return nil
	})
	if err != nil {
		return err
	}

	// Store all links in cache for potential use by other plugins
	m.Cache().Set("links", allLinks)

	// Assign inlinks and outlinks to each post
	p.assignLinksToPost(posts, allLinks)

	return nil
}

// buildBaseURL constructs the absolute URL for a post.
func (p *LinkCollectorPlugin) buildBaseURL(post *models.Post) string {
	if p.siteURL == "" {
		return post.Href
	}
	return strings.TrimSuffix(p.siteURL, "/") + post.Href
}

type anchorExtractionState struct {
	hrefs       []string
	textMap     map[string]string
	seen        map[string]bool
	currentHref string
	currentText strings.Builder
	inAnchor    bool
}

func newAnchorExtractionState() *anchorExtractionState {
	return &anchorExtractionState{
		hrefs:   make([]string, 0),
		textMap: make(map[string]string),
		seen:    make(map[string]bool),
	}
}

func (s *anchorExtractionState) flushAnchor() {
	href := strings.TrimSpace(s.currentHref)
	if href == "" || href == "#" || s.seen[href] {
		return
	}
	s.seen[href] = true
	s.hrefs = append(s.hrefs, href)
	s.textMap[href] = strings.Join(strings.Fields(s.currentText.String()), " ")
}

func (s *anchorExtractionState) startAnchor(token html.Token) {
	// Malformed nested anchors should not discard the outer href.
	if s.inAnchor {
		s.flushAnchor()
	}
	s.currentHref = anchorHref(token)
	s.currentText.Reset()
	s.inAnchor = true
}

func (s *anchorExtractionState) finishAnchor() {
	if s.inAnchor {
		s.flushAnchor()
	}
	s.currentHref = ""
	s.currentText.Reset()
	s.inAnchor = false
}

func (s *anchorExtractionState) addText(text string) {
	if s.inAnchor {
		s.currentText.WriteString(text)
	}
}

func (s *anchorExtractionState) addSeparator(tag string) {
	if s.inAnchor && linkTextSeparatorTag(tag) && s.currentText.Len() > 0 {
		s.currentText.WriteByte(' ')
	}
}

func anchorHref(token html.Token) string {
	for _, attr := range token.Attr {
		if strings.EqualFold(attr.Key, "href") {
			return attr.Val
		}
	}
	return ""
}

// extractHrefsAndText extracts all href values and their associated visible
// link text from HTML. It uses the HTML tokenizer so anchors with nested markup
// (for example card divs, spans, or images) are collected correctly.
func extractHrefsAndText(htmlContent string) (hrefs []string, textMap map[string]string) {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlContent))
	state := newAnchorExtractionState()

	for {
		switch tokenizer.Next() {
		case html.ErrorToken:
			state.finishAnchor()
			return state.hrefs, state.textMap
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				state.startAnchor(token)
				continue
			}
			state.addSeparator(token.Data)
		case html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				state.startAnchor(token)
				state.finishAnchor()
				continue
			}
			state.addSeparator(token.Data)
		case html.TextToken:
			state.addText(tokenizer.Token().Data)
		case html.EndTagToken:
			token := tokenizer.Token()
			if token.Data == "a" {
				state.finishAnchor()
				continue
			}
			state.addSeparator(token.Data)
		case html.CommentToken, html.DoctypeToken:
			continue
		}
	}
}

func linkTextSeparatorTag(tag string) bool {
	switch tag {
	case "br", "div", "p", "li", "dt", "dd", "h1", "h2", "h3", "h4", "h5", "h6", "blockquote":
		return true
	default:
		return false
	}
}

// extractHrefTextMap builds a map from href to link text by scanning HTML.
// Used when hrefs are restored from cache but link text is not cached.
func extractHrefTextMap(htmlContent string) map[string]string {
	_, textMap := extractHrefsAndText(htmlContent)
	return textMap
}

// extractHrefs extracts all href values from HTML content.
func extractHrefs(htmlContent string) []string {
	hrefs, _ := extractHrefsAndText(htmlContent)
	return hrefs
}

// createLink creates a Link object from an href found in a post.
func (p *LinkCollectorPlugin) createLink(sourcePost *models.Post, baseURL, href string, idx *lifecycle.PostIndex, hrefTextMap map[string]string) *models.Link {
	// Resolve the href against the base URL
	targetURL := resolveURL(baseURL, href)
	if targetURL == "" {
		return nil
	}

	// Parse the target URL to extract domain
	parsed, err := url.Parse(targetURL)
	if err != nil {
		return nil
	}

	targetDomain := parsed.Host

	// Determine if internal link
	isInternal := p.isInternalLink(targetDomain)

	// Look up target post if internal
	var targetPost *models.Post
	if isInternal {
		targetPost = p.resolveTargetPost(parsed.Path, idx)
	}

	// Determine if self-link
	isSelf := targetPost != nil && targetPost.Slug == sourcePost.Slug

	// Look up link text from pre-extracted map (avoids re-scanning HTML per href)
	sourceText := hrefTextMap[href]

	link := &models.Link{
		SourceURL:    baseURL,
		SourcePost:   sourcePost,
		TargetPost:   targetPost,
		RawTarget:    href,
		TargetURL:    targetURL,
		TargetDomain: targetDomain,
		IsInternal:   isInternal,
		IsSelf:       isSelf,
		SourceText:   sourceText,
	}

	// Set target text if we have a target post
	if targetPost != nil && targetPost.Title != nil {
		link.TargetText = *targetPost.Title
	}

	return link
}

// resolveURL resolves a potentially relative URL against a base URL.
func resolveURL(baseURL, href string) string {
	// Handle absolute URLs
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	// Handle protocol-relative URLs
	if strings.HasPrefix(href, "//") {
		return "https:" + href
	}

	// Handle anchor-only links
	if strings.HasPrefix(href, "#") {
		return baseURL + href
	}

	// Handle mailto, tel, javascript, etc.
	if strings.Contains(href, ":") && !strings.HasPrefix(href, "/") {
		return href
	}

	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	// Parse href
	ref, err := url.Parse(href)
	if err != nil {
		return href
	}

	// Resolve against base
	resolved := base.ResolveReference(ref)
	return resolved.String()
}

// isInternalLink checks if the given domain matches the site domain.
func (p *LinkCollectorPlugin) isInternalLink(targetDomain string) bool {
	if p.siteDomain == "" {
		// If no site URL configured, treat root-relative links as internal
		return targetDomain == ""
	}
	return strings.EqualFold(targetDomain, p.siteDomain)
}

// resolveTargetPost looks up a post by its path/slug using the shared PostIndex.
func (p *LinkCollectorPlugin) resolveTargetPost(path string, idx *lifecycle.PostIndex) *models.Post {
	// Normalize the path
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimPrefix(path, "/")

	// Try direct slug lookup (BySlug is already lowercase)
	if post, ok := idx.BySlug[strings.ToLower(path)]; ok {
		return post
	}

	// Try href lookup with trailing slash
	hrefPath := "/" + path + "/"
	if post, ok := idx.ByHref[hrefPath]; ok {
		return post
	}

	// Try without trailing slash
	hrefPathNoSlash := "/" + path
	if post, ok := idx.ByHref[hrefPathNoSlash]; ok {
		return post
	}

	return nil
}

// assignLinksToPost populates inlinks and outlinks for each post.
func (p *LinkCollectorPlugin) assignLinksToPost(posts []*models.Post, allLinks []*models.Link) {
	// Build inlinks map: target post slug -> list of links
	inlinksMap := make(map[string][]*models.Link)
	// Build outlinks map: source post slug -> list of links
	outlinksMap := make(map[string][]*models.Link)

	for _, link := range allLinks {
		// Skip self-links
		if link.IsSelf {
			continue
		}

		// Add to outlinks for source post
		if link.SourcePost != nil {
			outlinksMap[link.SourcePost.Slug] = append(outlinksMap[link.SourcePost.Slug], link)
		}

		// Add to inlinks for target post (if internal and resolved)
		if link.TargetPost != nil {
			// Check if we should include this inlink based on config
			if !p.shouldIncludeInlink(link) {
				continue
			}
			inlinksMap[link.TargetPost.Slug] = append(inlinksMap[link.TargetPost.Slug], link)
		}
	}

	// Assign to posts
	for _, post := range posts {
		// Deduplicate inlinks by source URL
		post.Inlinks = deduplicateLinksBySource(inlinksMap[post.Slug])

		// Deduplicate outlinks by target URL
		post.Outlinks = deduplicateLinksByTarget(outlinksMap[post.Slug])
	}
}

// shouldIncludeInlink checks if an inlink should be included based on configuration.
func (p *LinkCollectorPlugin) shouldIncludeInlink(link *models.Link) bool {
	if link.SourcePost == nil {
		return false
	}

	slug := link.SourcePost.Slug

	// Check index exclusion
	if !p.includeIndex && slug == "index" {
		return false
	}

	// Check feed exclusion (simple heuristic: check if post has no content)
	// A more robust check would involve checking against known feed slugs
	if !p.includeFeeds {
		// Check if the source post appears to be a feed page
		if isFeedPost(link.SourcePost) {
			return false
		}
	}

	return true
}

// isFeedPost checks if a post appears to be a feed/index page.
// Feed pages typically have a "feed" template or are in the feeds list.
func isFeedPost(post *models.Post) bool {
	if post.Template == "feed.html" || post.Template == "archive.html" {
		return true
	}
	// Check for feed marker in extra
	if isFeed, ok := post.Extra["is_feed"].(bool); ok && isFeed {
		return true
	}
	return false
}

// deduplicateLinksBySource removes duplicate links with the same source URL.
func deduplicateLinksBySource(links []*models.Link) []*models.Link {
	if len(links) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	result := make([]*models.Link, 0, len(links))

	for _, link := range links {
		if !seen[link.SourceURL] {
			seen[link.SourceURL] = true
			result = append(result, link)
		}
	}

	return result
}

// deduplicateLinksByTarget removes duplicate links with the same target URL.
func deduplicateLinksByTarget(links []*models.Link) []*models.Link {
	if len(links) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	result := make([]*models.Link, 0, len(links))

	for _, link := range links {
		if !seen[link.TargetURL] {
			seen[link.TargetURL] = true
			result = append(result, link)
		}
	}

	return result
}

// SetIncludeFeeds enables or disables including feed pages in inlinks.
func (p *LinkCollectorPlugin) SetIncludeFeeds(include bool) {
	p.includeFeeds = include
}

// SetIncludeIndex enables or disables including the index page in inlinks.
func (p *LinkCollectorPlugin) SetIncludeIndex(include bool) {
	p.includeIndex = include
}

// SetSiteURL sets the site URL for determining internal vs external links.
func (p *LinkCollectorPlugin) SetSiteURL(siteURL string) {
	p.siteURL = siteURL
	if parsed, err := url.Parse(siteURL); err == nil {
		p.siteDomain = parsed.Host
	}
}

// Ensure LinkCollectorPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*LinkCollectorPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*LinkCollectorPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*LinkCollectorPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*LinkCollectorPlugin)(nil)
)
