// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// GlossaryConfig holds configuration for the glossary plugin.
type GlossaryConfig struct {
	// Enabled controls whether the plugin is active
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// LinkClass is the CSS class for glossary links (default: "glossary-term")
	LinkClass string `json:"link_class" yaml:"link_class" toml:"link_class"`

	// CaseSensitive controls whether term matching is case-sensitive (default: false)
	CaseSensitive bool `json:"case_sensitive" yaml:"case_sensitive" toml:"case_sensitive"`

	// Tooltip controls whether to add hover-card data (a title attribute with the
	// description and data-hover-title with the term name) (default: true)
	Tooltip bool `json:"tooltip" yaml:"tooltip" toml:"tooltip"`

	// MaxLinksPerTerm limits how many times each term is linked (0 = all, default: 1)
	MaxLinksPerTerm int `json:"max_links_per_term" yaml:"max_links_per_term" toml:"max_links_per_term"`

	// ExcludeTags lists tags that should not have glossary terms linked
	ExcludeTags []string `json:"exclude_tags" yaml:"exclude_tags" toml:"exclude_tags"`

	// ExportJSON controls whether to export glossary.json (default: true)
	ExportJSON bool `json:"export_json" yaml:"export_json" toml:"export_json"`

	// GlossaryPath is the path prefix for glossary posts (default: "glossary")
	GlossaryPath string `json:"glossary_path" yaml:"glossary_path" toml:"glossary_path"`

	// TemplateKey identifies glossary posts by templateKey frontmatter
	TemplateKey string `json:"template_key" yaml:"template_key" toml:"template_key"`
}

// NewGlossaryConfig creates a GlossaryConfig with default values.
func NewGlossaryConfig() *GlossaryConfig {
	return &GlossaryConfig{
		Enabled:         true,
		LinkClass:       "glossary-term",
		CaseSensitive:   false,
		Tooltip:         true,
		MaxLinksPerTerm: 1,
		ExcludeTags:     []string{"glossary"},
		ExportJSON:      true,
		GlossaryPath:    "glossary",
		TemplateKey:     "glossary",
	}
}

// GlossaryTerm represents a single glossary term with its definition.
type GlossaryTerm struct {
	// Term is the primary term name
	Term string `json:"term"`

	// Slug is the URL-safe identifier
	Slug string `json:"slug"`

	// Description is the term description
	Description string `json:"description"`

	// Aliases are alternative names for the term
	Aliases []string `json:"aliases,omitempty"`

	// Href is the URL path to the term's page
	Href string `json:"href"`

	// post is the source post (not exported)
	post *models.Post
}

// GlossaryExport represents the JSON export format.
type GlossaryExport struct {
	Terms []*GlossaryTerm `json:"terms"`
}

// GlossaryPlugin automatically links glossary terms in post content.
type GlossaryPlugin struct {
	config *GlossaryConfig

	// terms maps lowercase term/alias -> GlossaryTerm for lookup
	terms map[string]*GlossaryTerm

	// allTerms holds all unique glossary terms for export
	allTerms []*GlossaryTerm

	// mu protects terms map during concurrent access
	mu sync.RWMutex

	// termRegexCache caches compiled regexes for terms to avoid repeated compilation.
	// Key format: "caseSensitive:termKey"
	termRegexCache map[string]*regexp.Regexp
	regexCacheMu   sync.RWMutex
}

// NewGlossaryPlugin creates a new GlossaryPlugin with default configuration.
func NewGlossaryPlugin() *GlossaryPlugin {
	return &GlossaryPlugin{
		config:         NewGlossaryConfig(),
		terms:          make(map[string]*GlossaryTerm),
		allTerms:       make([]*GlossaryTerm, 0),
		termRegexCache: make(map[string]*regexp.Regexp),
	}
}

// Name returns the unique name of the plugin.
func (p *GlossaryPlugin) Name() string {
	return "glossary"
}

// Configure reads configuration options for the plugin.
func (p *GlossaryPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Read glossary config from Extra["glossary"]
	glossaryConfig, ok := config.Extra["glossary"].(map[string]interface{})
	if !ok {
		return nil
	}

	if enabled, ok := glossaryConfig["enabled"].(bool); ok {
		p.config.Enabled = enabled
	}
	if linkClass, ok := glossaryConfig["link_class"].(string); ok {
		p.config.LinkClass = linkClass
	}
	if caseSensitive, ok := glossaryConfig["case_sensitive"].(bool); ok {
		p.config.CaseSensitive = caseSensitive
	}
	if tooltip, ok := glossaryConfig["tooltip"].(bool); ok {
		p.config.Tooltip = tooltip
	}
	if maxLinks, ok := parseIntFromInterface(glossaryConfig["max_links_per_term"]); ok {
		p.config.MaxLinksPerTerm = maxLinks
	}
	if exportJSON, ok := glossaryConfig["export_json"].(bool); ok {
		p.config.ExportJSON = exportJSON
	}
	if glossaryPath, ok := glossaryConfig["glossary_path"].(string); ok {
		p.config.GlossaryPath = glossaryPath
	}
	if templateKey, ok := glossaryConfig["template_key"].(string); ok {
		p.config.TemplateKey = templateKey
	}

	// Parse exclude_tags
	if excludeTags, ok := glossaryConfig["exclude_tags"].([]interface{}); ok {
		p.config.ExcludeTags = make([]string, 0, len(excludeTags))
		for _, tag := range excludeTags {
			if tagStr, ok := tag.(string); ok {
				p.config.ExcludeTags = append(p.config.ExcludeTags, tagStr)
			}
		}
	}
	if excludeTags, ok := glossaryConfig["exclude_tags"].([]string); ok {
		p.config.ExcludeTags = excludeTags
	}

	return nil
}

// Priority returns the plugin priority for the given stage.
// Glossary should run late in render (post_render) to process article_html.
func (p *GlossaryPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Render processes glossary terms and links them in post content.
func (p *GlossaryPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	posts := m.Posts()

	// Build glossary term lookup
	if err := p.buildGlossary(posts); err != nil {
		return fmt.Errorf("building glossary: %w", err)
	}

	// No terms found, nothing to do
	if len(p.terms) == 0 {
		return nil
	}

	// Compute a stable hash of all glossary terms so we can detect when terms change.
	// This is combined with each post's ArticleHTML hash to form the cache key.
	termsHash := p.computeTermsHash()

	cache := GetBuildCache(m)

	postsToProcess := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip {
			return false
		}
		if post.ArticleHTML == "" {
			return false
		}
		if p.isGlossaryPost(post) {
			return false
		}
		if p.hasExcludedTag(post) {
			return false
		}
		return true
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

	// Phase 1: Restore cached results for unchanged posts
	var needProcessing []*models.Post
	if cache != nil {
		for _, post := range postsToProcess {
			combinedHash := buildcache.ContentHash(post.ArticleHTML + termsHash)
			if cached, ok := cache.GetCachedGlossaryHTML(post.Path, combinedHash); ok {
				post.ArticleHTML = cached
			} else {
				needProcessing = append(needProcessing, post)
			}
		}
	} else {
		needProcessing = postsToProcess
	}

	if len(needProcessing) == 0 {
		return nil
	}

	// Phase 2: Process posts that need updating, concurrently
	return m.ProcessPostsSliceConcurrently(needProcessing, func(post *models.Post) error {
		combinedHash := buildcache.ContentHash(post.ArticleHTML + termsHash)
		post.ArticleHTML = p.linkTerms(post.ArticleHTML, post)
		if cache != nil {
			cache.CacheGlossaryHTML(post.Path, combinedHash, post.ArticleHTML)
		}
		return nil
	})
}

// computeTermsHash builds a deterministic hash of all glossary terms and their
// configuration so that cache entries are invalidated when glossary content changes.
func (p *GlossaryPlugin) computeTermsHash() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// allTerms is already sorted by term name from buildGlossary
	var b strings.Builder
	b.WriteString("glossary-deterministic-protection-v3\x00")
	for _, term := range p.allTerms {
		b.WriteString(term.Term)
		b.WriteByte('\x00')
		b.WriteString(term.Href)
		b.WriteByte('\x00')
		b.WriteString(term.Description)
		b.WriteByte('\x00')
		for _, alias := range term.Aliases {
			b.WriteString(alias)
			b.WriteByte('\x01')
		}
		b.WriteByte('\x00')
	}
	return buildcache.ContentHash(b.String())
}

// Write exports the glossary JSON file if configured.
func (p *GlossaryPlugin) Write(m *lifecycle.Manager) error {
	if !p.config.Enabled || !p.config.ExportJSON {
		return nil
	}

	if len(p.allTerms) == 0 {
		return nil
	}

	config := m.Config()
	outputPath := filepath.Join(config.OutputDir, "glossary.json")

	// Ensure output directory exists
	if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	export := &GlossaryExport{Terms: p.allTerms}
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling glossary JSON: %w", err)
	}

	//nolint:gosec // G306: glossary.json is a public file, 0644 is appropriate
	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return fmt.Errorf("writing glossary.json: %w", err)
	}

	return nil
}

// buildGlossary scans posts for glossary terms and builds the lookup map.
func (p *GlossaryPlugin) buildGlossary(posts []*models.Post) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.terms = make(map[string]*GlossaryTerm)
	p.allTerms = make([]*GlossaryTerm, 0)

	for _, post := range posts {
		if post.Skip {
			continue
		}

		if !p.isGlossaryPost(post) {
			continue
		}

		term := p.extractGlossaryTerm(post)
		if term == nil {
			continue
		}

		p.allTerms = append(p.allTerms, term)

		// Add primary term to lookup
		key := p.normalizeKey(term.Term)
		p.terms[key] = term

		// Add aliases to lookup
		for _, alias := range term.Aliases {
			aliasKey := p.normalizeKey(alias)
			p.terms[aliasKey] = term
		}
	}

	// Sort allTerms by term name for consistent export
	sort.Slice(p.allTerms, func(i, j int) bool {
		return strings.ToLower(p.allTerms[i].Term) < strings.ToLower(p.allTerms[j].Term)
	})

	return nil
}

// isGlossaryPost checks if a post is a glossary definition post.
func (p *GlossaryPlugin) isGlossaryPost(post *models.Post) bool {
	// Check templateKey in Extra
	if templateKey, ok := post.Extra["templateKey"].(string); ok {
		if templateKey == p.config.TemplateKey {
			return true
		}
	}

	// Check template_key variant
	if templateKey, ok := post.Extra["template_key"].(string); ok {
		if templateKey == p.config.TemplateKey {
			return true
		}
	}

	// Check if post is in glossary path
	if p.config.GlossaryPath != "" {
		if strings.HasPrefix(post.Path, p.config.GlossaryPath+"/") ||
			strings.HasPrefix(post.Path, p.config.GlossaryPath+"\\") ||
			strings.Contains(post.Path, "/"+p.config.GlossaryPath+"/") ||
			strings.Contains(post.Path, "\\"+p.config.GlossaryPath+"\\") {
			return true
		}
	}

	return false
}

// extractGlossaryTerm creates a GlossaryTerm from a glossary post.
func (p *GlossaryPlugin) extractGlossaryTerm(post *models.Post) *GlossaryTerm {
	// Get the term name (title)
	termName := ""
	if post.PlainTitle() != "" {
		termName = post.PlainTitle()
	} else {
		// Use slug as fallback
		termName = post.Slug
	}

	if termName == "" {
		return nil
	}

	// Get description
	description := ""
	if post.Description != nil {
		description = *post.Description
	}

	// Get aliases from Extra
	var aliases []string
	if aliasesRaw, ok := post.Extra["aliases"].([]interface{}); ok {
		for _, a := range aliasesRaw {
			if alias, ok := a.(string); ok {
				aliases = append(aliases, alias)
			}
		}
	}
	if aliasesStr, ok := post.Extra["aliases"].([]string); ok {
		aliases = aliasesStr
	}

	// Build href
	href := post.Href
	if href == "" {
		href = "/" + post.Slug + "/"
	}

	return &GlossaryTerm{
		Term:        termName,
		Slug:        post.Slug,
		Description: description,
		Aliases:     aliases,
		Href:        href,
		post:        post,
	}
}

// processPost links glossary terms in a single post's article_html.
//
//nolint:unparam // error return matches PostProcessor signature; kept for test compatibility
func (p *GlossaryPlugin) processPost(post *models.Post) error {
	if post.Skip {
		return nil
	}

	// Skip if no article HTML
	if post.ArticleHTML == "" {
		return nil
	}

	// Skip glossary posts themselves
	if p.isGlossaryPost(post) {
		return nil
	}

	// Skip posts with excluded tags
	if p.hasExcludedTag(post) {
		return nil
	}

	// Process the article HTML
	post.ArticleHTML = p.linkTerms(post.ArticleHTML, post)

	return nil
}

// hasExcludedTag checks if a post has any of the excluded tags.
func (p *GlossaryPlugin) hasExcludedTag(post *models.Post) bool {
	for _, postTag := range post.Tags {
		for _, excludeTag := range p.config.ExcludeTags {
			if strings.EqualFold(postTag, excludeTag) {
				return true
			}
		}
	}
	return false
}

// normalizeKey normalizes a term for lookup based on case sensitivity setting.
func (p *GlossaryPlugin) normalizeKey(term string) string {
	if p.config.CaseSensitive {
		return term
	}
	return strings.ToLower(term)
}

// Regex patterns for protected content in glossary processing
var (
	// glossaryAnchorTagRegex matches content inside <a>...</a> tags (including nested tags)
	glossaryAnchorTagRegex = regexp.MustCompile(`(?is)<a\s[^>]*>.*?</a>`)

	// glossaryCodeTagRegex matches content inside <code>...</code> tags
	glossaryCodeTagRegex = regexp.MustCompile(`(?is)<code[^>]*>.*?</code>`)

	// glossaryPreTagRegex matches content inside <pre>...</pre> tags
	glossaryPreTagRegex = regexp.MustCompile(`(?is)<pre[^>]*>.*?</pre>`)

	// glossaryTagRegex matches any remaining HTML tag so terms inside
	// attributes (id, href, alt, title, ...) are never rewritten.
	glossaryTagRegex = regexp.MustCompile(`<[^>]+>`)
)

// placeholder is used to mark protected content
const placeholder = "\x00GLOSSARY_PROTECTED_%d\x00"

// linkTerms replaces glossary terms with linked versions in HTML content.
func (p *GlossaryPlugin) linkTerms(htmlContent string, currentPost *models.Post) string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Track how many times each term has been linked
	linkCounts := make(map[string]int)

	// Protect content that should not be modified
	protectedSegments := make(map[string]string)
	protectedIdx := 0

	// Protect <a> tags
	htmlContent = glossaryAnchorTagRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		key := fmt.Sprintf(placeholder, protectedIdx)
		protectedSegments[key] = match
		protectedIdx++
		return key
	})

	// Protect <pre> tags (before <code> since pre may contain code)
	htmlContent = glossaryPreTagRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		key := fmt.Sprintf(placeholder, protectedIdx)
		protectedSegments[key] = match
		protectedIdx++
		return key
	})

	// Protect <code> tags
	htmlContent = glossaryCodeTagRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		key := fmt.Sprintf(placeholder, protectedIdx)
		protectedSegments[key] = match
		protectedIdx++
		return key
	})

	// Protect every other tag so attribute values are left untouched
	htmlContent = glossaryTagRegex.ReplaceAllStringFunc(htmlContent, func(match string) string {
		key := fmt.Sprintf(placeholder, protectedIdx)
		protectedSegments[key] = match
		protectedIdx++
		return key
	})

	// Build list of terms sorted by length (longest first) to match longer terms first
	termList := make([]string, 0, len(p.terms))
	for term := range p.terms {
		termList = append(termList, term)
	}
	sort.Slice(termList, func(i, j int) bool {
		if len(termList[i]) == len(termList[j]) {
			return termList[i] < termList[j]
		}
		return len(termList[i]) > len(termList[j])
	})

	// Process each term
	for _, termKey := range termList {
		glossaryTerm := p.terms[termKey]

		// Don't link a term to itself (glossary post linking to its own page)
		if glossaryTerm.post == currentPost {
			continue
		}

		// Check max links per term
		termID := glossaryTerm.Slug
		if p.config.MaxLinksPerTerm > 0 && linkCounts[termID] >= p.config.MaxLinksPerTerm {
			continue
		}

		// Get cached regex pattern for this term (word boundary matching)
		pattern := p.getTermRegex(termKey, p.config.CaseSensitive)

		// Find and replace matches
		htmlContent = pattern.ReplaceAllStringFunc(htmlContent, func(match string) string {
			// Check if we've hit the limit
			if p.config.MaxLinksPerTerm > 0 && linkCounts[termID] >= p.config.MaxLinksPerTerm {
				return match
			}

			// Protect generated links as well as existing HTML. Later, shorter
			// glossary terms must never match inside href, title, or link text.
			key := fmt.Sprintf(placeholder, protectedIdx)
			protectedSegments[key] = p.buildLink(glossaryTerm, match)
			protectedIdx++
			linkCounts[termID]++
			return key
		})
	}

	// Restore protected segments
	// Later protected segments can contain earlier placeholders, so restore
	// outer segments before their contents.
	for i := protectedIdx - 1; i >= 0; i-- {
		key := fmt.Sprintf(placeholder, i)
		htmlContent = strings.Replace(htmlContent, key, protectedSegments[key], 1)
	}

	return htmlContent
}

// buildLink creates an HTML anchor tag for a glossary term.
func (p *GlossaryPlugin) buildLink(term *GlossaryTerm, matchedText string) string {
	var attrs strings.Builder

	//nolint:gocritic // sprintfQuotedString: %q produces Go escaping, but we need HTML entity escaping here
	_, _ = attrs.WriteString(fmt.Sprintf(`href="%s"`, html.EscapeString(term.Href)))

	if p.config.LinkClass != "" {
		//nolint:gocritic // sprintfQuotedString: %q produces Go escaping, but we need HTML entity escaping here
		_, _ = attrs.WriteString(fmt.Sprintf(` class="%s"`, html.EscapeString(p.config.LinkClass)))
	}

	if p.config.Tooltip && term.Description != "" {
		//nolint:gocritic // sprintfQuotedString: %q produces Go escaping, but we need HTML entity escaping here
		_, _ = attrs.WriteString(fmt.Sprintf(` title="%s"`, html.EscapeString(term.Description)))
		if term.Term != "" {
			//nolint:gocritic // sprintfQuotedString: %q produces Go escaping, but we need HTML entity escaping here
			_, _ = attrs.WriteString(fmt.Sprintf(` data-hover-title="%s"`, html.EscapeString(term.Term)))
		}
	}

	return fmt.Sprintf(`<a %s>%s</a>`, attrs.String(), html.EscapeString(matchedText))
}

// Config returns the current glossary configuration.
func (p *GlossaryPlugin) Config() *GlossaryConfig {
	return p.config
}

// SetConfig sets the glossary configuration.
func (p *GlossaryPlugin) SetConfig(config *GlossaryConfig) {
	p.config = config
}

// Terms returns a copy of all glossary terms.
func (p *GlossaryPlugin) Terms() []*GlossaryTerm {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]*GlossaryTerm, len(p.allTerms))
	copy(result, p.allTerms)
	return result
}

// getTermRegex returns a cached regex for matching the given term.
// The caseSensitive flag determines whether the regex should be case-insensitive.
func (p *GlossaryPlugin) getTermRegex(termKey string, caseSensitive bool) *regexp.Regexp {
	// Build cache key
	cacheKey := "i:" + termKey // case-insensitive prefix
	if caseSensitive {
		cacheKey = "s:" + termKey // case-sensitive prefix
	}

	// Check cache first
	p.regexCacheMu.RLock()
	pattern, ok := p.termRegexCache[cacheKey]
	p.regexCacheMu.RUnlock()

	if ok {
		return pattern
	}

	// Compile and cache
	escapedTerm := regexp.QuoteMeta(termKey)
	var patternStr string
	if caseSensitive {
		patternStr = `\b(` + escapedTerm + `)\b`
	} else {
		patternStr = `(?i)\b(` + escapedTerm + `)\b`
	}
	pattern = regexp.MustCompile(patternStr)

	p.regexCacheMu.Lock()
	p.termRegexCache[cacheKey] = pattern
	p.regexCacheMu.Unlock()

	return pattern
}

// Ensure GlossaryPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*GlossaryPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*GlossaryPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*GlossaryPlugin)(nil)
	_ lifecycle.WritePlugin     = (*GlossaryPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*GlossaryPlugin)(nil)
)
