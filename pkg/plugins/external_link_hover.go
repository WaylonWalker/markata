// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"html"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// ExternalLinkHoverConfig configures the external_link_hover plugin.
type ExternalLinkHoverConfig struct {
	// Enabled controls whether external links get hover preview data (default: false).
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled"`

	// Fetch resolves metadata for URLs with no cached embeds metadata (default: false).
	// When false the plugin makes no network requests.
	Fetch bool `json:"fetch" yaml:"fetch" toml:"fetch"`

	// IncludeImage adds data-link-image when metadata has an image (default: true).
	IncludeImage bool `json:"include_image" yaml:"include_image" toml:"include_image"`

	// FaviconService is a URL template for a site icon. {host} and {origin} are
	// replaced. Empty disables icons (default: "").
	FaviconService string `json:"favicon_service" yaml:"favicon_service" toml:"favicon_service"`

	// IgnoreDomains lists hosts (and their subdomains) to skip.
	IgnoreDomains []string `json:"ignore_domains" yaml:"ignore_domains" toml:"ignore_domains"`
}

// NewExternalLinkHoverConfig returns the default configuration.
func NewExternalLinkHoverConfig() ExternalLinkHoverConfig {
	return ExternalLinkHoverConfig{IncludeImage: true}
}

// ExternalLinkHoverPlugin adds data-link-* attributes to external links so the
// theme's hover cards can preview them. Metadata comes from the embeds cache,
// blogroll feed config, and the link's own title attribute.
type ExternalLinkHoverPlugin struct {
	config     ExternalLinkHoverConfig
	siteHost   string
	embeds     *EmbedsPlugin
	blogroll   map[string]models.ExternalFeedConfig
	mu         sync.Mutex
	resolved   map[string]*externalLinkResolution
	fetchTitle string
}

type externalLinkResolution struct {
	once     sync.Once
	metadata *OGMetadata
}

// NewExternalLinkHoverPlugin creates the plugin with default settings.
func NewExternalLinkHoverPlugin() *ExternalLinkHoverPlugin {
	return &ExternalLinkHoverPlugin{config: NewExternalLinkHoverConfig()}
}

// Name returns the unique name of the plugin.
func (p *ExternalLinkHoverPlugin) Name() string {
	return "external_link_hover"
}

// Priority runs the plugin after encryption (50) and, by registration order
// within PriorityLate, after glossary and before templates render ArticleHTML.
func (p *ExternalLinkHoverPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageRender {
		return lifecycle.PriorityLate
	}
	return lifecycle.PriorityDefault
}

// Configure reads [markata-go.external_link_hover] and prepares metadata sources.
func (p *ExternalLinkHoverPlugin) Configure(m *lifecycle.Manager) error {
	cfg := m.Config()
	p.config = parseExternalLinkHoverConfig(cfg)
	p.resolved = map[string]*externalLinkResolution{}
	if !p.config.Enabled {
		return nil
	}

	if origin := getSiteOrigin(cfg); origin != "" {
		if u, err := url.Parse(origin); err == nil {
			p.siteHost = strings.ToLower(u.Hostname())
		}
	}

	p.embeds = NewEmbedsPlugin()
	if err := p.embeds.Configure(m); err != nil {
		return err
	}
	p.fetchTitle = p.embeds.config.FallbackTitle

	p.blogroll = map[string]models.ExternalFeedConfig{}
	feeds := getBlogrollConfig(cfg).Feeds
	for i := range feeds {
		feed := feeds[i]
		for _, raw := range []string{feed.SiteURL, feed.URL} {
			host := externalLinkHost(raw)
			if host == "" {
				continue
			}
			if _, exists := p.blogroll[host]; !exists {
				p.blogroll[host] = feed
			}
		}
	}
	return nil
}

func parseExternalLinkHoverConfig(cfg *lifecycle.Config) ExternalLinkHoverConfig {
	result := NewExternalLinkHoverConfig()
	if cfg == nil || cfg.Extra == nil {
		return result
	}
	raw, ok := cfg.Extra["external_link_hover"].(map[string]interface{})
	if !ok {
		return result
	}
	if v, ok := raw["enabled"].(bool); ok {
		result.Enabled = v
	}
	if v, ok := raw["fetch"].(bool); ok {
		result.Fetch = v
	}
	if v, ok := raw["include_image"].(bool); ok {
		result.IncludeImage = v
	}
	if v, ok := raw["favicon_service"].(string); ok {
		result.FaviconService = strings.TrimSpace(v)
	}
	for _, d := range parseStringSlice(raw["ignore_domains"]) {
		if d = strings.ToLower(strings.TrimSpace(d)); d != "" {
			result.IgnoreDomains = append(result.IgnoreDomains, d)
		}
	}
	return result
}

// Render annotates external links in every eligible post.
func (p *ExternalLinkHoverPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}
	posts := m.FilterPosts(func(post *models.Post) bool {
		return !post.Skip && !post.Private && strings.Contains(post.ArticleHTML, "http")
	})
	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		post.ArticleHTML = p.annotateHTML(post.ArticleHTML)
		return nil
	})
}

var (
	externalAnchorOpenRegex = regexp.MustCompile(`<a\s[^>]*>`)
	externalAnchorAttrRegex = regexp.MustCompile(`([a-zA-Z_:][-a-zA-Z0-9_:.]*)(?:\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'>]+)))?`)
)

var externalLinkSkipClasses = []string{"wikilink", "mention", "glossary-term", "heading-anchor"}

var externalLinkSkipAttrs = []string{
	"data-link-preview", "data-hover-card", "data-hover-title", "data-hover-description", "data-no-preview",
}

func (p *ExternalLinkHoverPlugin) annotateHTML(content string) string {
	if !strings.Contains(content, "<a") {
		return content
	}
	return externalAnchorOpenRegex.ReplaceAllStringFunc(content, p.annotateTag)
}

// parseAnchorAttrs returns the lower-cased attribute names and unescaped values
// of an opening <a ...> tag.
func parseAnchorAttrs(tag string) map[string]string {
	inner := strings.TrimSuffix(strings.TrimPrefix(tag, "<a"), ">")
	inner = strings.TrimSuffix(inner, "/")
	attrs := map[string]string{}
	for _, m := range externalAnchorAttrRegex.FindAllStringSubmatch(inner, -1) {
		name := strings.ToLower(m[1])
		if _, seen := attrs[name]; seen {
			continue
		}
		attrs[name] = html.UnescapeString(m[2] + m[3] + m[4])
	}
	return attrs
}

func (p *ExternalLinkHoverPlugin) annotateTag(tag string) string {
	attrs := parseAnchorAttrs(tag)
	href := strings.TrimSpace(attrs["href"])
	if href == "" {
		return tag
	}
	for _, name := range externalLinkSkipAttrs {
		if _, ok := attrs[name]; ok {
			return tag
		}
	}
	for _, class := range strings.Fields(attrs["class"]) {
		for _, skip := range externalLinkSkipClasses {
			if class == skip {
				return tag
			}
		}
		if strings.Contains(class, "embed-card") || strings.Contains(class, "link-card") {
			return tag
		}
	}

	parsed, err := url.Parse(href)
	if err != nil || (parsed.Scheme != schemeHTTP && parsed.Scheme != schemeHTTPS) || parsed.Host == "" {
		return tag
	}
	host := strings.ToLower(parsed.Hostname())
	if host == p.siteHost || p.ignored(host) {
		return tag
	}

	data := p.linkData(href, parsed, host, attrs["title"])
	var b strings.Builder
	b.WriteString(` data-link-preview data-link-host="`)
	b.WriteString(html.EscapeString(strings.TrimPrefix(host, "www.")))
	b.WriteByte('"')
	for _, kv := range data {
		if kv[1] == "" {
			continue
		}
		b.WriteString(` data-link-`)
		b.WriteString(kv[0])
		b.WriteString(`="`)
		b.WriteString(html.EscapeString(kv[1]))
		b.WriteByte('"')
	}

	end := len(tag) - 1
	if strings.HasSuffix(tag, "/>") {
		end = len(tag) - 2
	}
	return tag[:end] + b.String() + tag[end:]
}

func (p *ExternalLinkHoverPlugin) ignored(host string) bool {
	for _, d := range p.config.IgnoreDomains {
		if host == d || strings.HasSuffix(host, "."+d) {
			return true
		}
	}
	return false
}

// linkData returns ordered [name, value] pairs for the data-link-* attributes.
func (p *ExternalLinkHoverPlugin) linkData(href string, parsed *url.URL, host, titleAttr string) [][2]string {
	var title, desc, site, image string
	if meta := p.metadata(href); meta != nil {
		if meta.Title != p.fetchTitle {
			title = strings.TrimSpace(meta.Title)
		}
		desc = strings.TrimSpace(meta.Description)
		site = strings.TrimSpace(firstNonEmptyString(meta.SiteName, meta.ProviderName))
		image = strings.TrimSpace(meta.Image)
	}
	if feed, ok := p.blogrollFeed(host); ok {
		site = firstNonEmptyString(site, feed.Title)
		desc = firstNonEmptyString(desc, feed.Description)
		image = firstNonEmptyString(image, feed.ImageURL)
	}
	desc = firstNonEmptyString(desc, strings.TrimSpace(titleAttr))
	if !p.config.IncludeImage {
		image = ""
	}
	if image != "" {
		if resolved, err := parsed.Parse(image); err == nil {
			image = resolved.String()
		}
	}

	return [][2]string{
		{"title", truncatePreviewText(title, 140)},
		{"description", truncatePreviewText(desc, 280)},
		{"site", site},
		{"image", image},
		{"icon", p.faviconURL(parsed)},
	}
}

func (p *ExternalLinkHoverPlugin) blogrollFeed(host string) (models.ExternalFeedConfig, bool) {
	feed, ok := p.blogroll[strings.TrimPrefix(host, "www.")]
	return feed, ok
}

func (p *ExternalLinkHoverPlugin) faviconURL(parsed *url.URL) string {
	if p.config.FaviconService == "" {
		return ""
	}
	r := strings.NewReplacer(
		"{host}", url.PathEscape(parsed.Hostname()),
		"{origin}", url.QueryEscape(parsed.Scheme+"://"+parsed.Host),
	)
	return r.Replace(p.config.FaviconService)
}

// metadata returns cached (or, with fetch enabled, fetched) embeds metadata,
// resolving each URL at most once per build.
func (p *ExternalLinkHoverPlugin) metadata(rawURL string) *OGMetadata {
	if p.embeds == nil {
		return nil
	}
	p.mu.Lock()
	res, ok := p.resolved[rawURL]
	if !ok {
		res = &externalLinkResolution{}
		p.resolved[rawURL] = res
	}
	p.mu.Unlock()

	res.once.Do(func() {
		resolved := normalizeHackerNewsURL(rawURL)
		res.metadata = p.embeds.fetchCachedMetadata(resolved, "oembed")
		if res.metadata == nil {
			res.metadata = p.embeds.fetchCachedMetadata(resolved, "og")
		}
		if res.metadata == nil && p.config.Fetch {
			res.metadata = p.embeds.fetchExternalMetadata(resolved)
		}
	})
	return res.metadata
}

func firstNonEmptyString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// externalLinkHost returns the lower-cased host of an absolute URL, without "www.".
func externalLinkHost(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" {
		return ""
	}
	return strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
}

// SetConfig sets the plugin configuration directly (useful in tests).
func (p *ExternalLinkHoverPlugin) SetConfig(config ExternalLinkHoverConfig) {
	p.config = config
	if p.resolved == nil {
		p.resolved = map[string]*externalLinkResolution{}
	}
}

// Ensure ExternalLinkHoverPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*ExternalLinkHoverPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ExternalLinkHoverPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*ExternalLinkHoverPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*ExternalLinkHoverPlugin)(nil)
)
