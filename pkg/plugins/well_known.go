package plugins

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

const (
	wellKnownDir           = ".well-known"
	wellKnownGeneratorName = "markata-go"
	wellKnownGeneratorVer  = "unknown"
	wellKnownNodeInfoRel   = "http://nodeinfo.diaspora.software/ns/schema/2.0"
	wellKnownWebfingerRel  = "http://webfinger.net/rel/profile-page"
	wellKnownAvatarRel     = "http://webfinger.net/rel/avatar"
	wellKnownHostMetaRel   = "lrdd"
	wellKnownHostMetaType  = "application/jrd+json"
	wellKnownWebfingerType = "text/html"
)

type wellKnownData struct {
	SiteURL          string
	SiteTitle        string
	SiteDescription  string
	SiteHost         string
	Author           string
	AuthorImageURL   string
	BuildTime        string
	WebfingerSubject string
	WebfingerURL     string
	WebfingerAliases []string
	NodeInfoURL      string
	GeneratorName    string
	GeneratorVersion string
	SSHFingerprint   string
	KeybaseUsername  string
	Links            []wellKnownLinksDomain
	InternalLinks    []wellKnownInternalTarget
}

type wellKnownLink struct {
	SourceURL string `json:"sourceUrl"`
	TargetURL string `json:"targetUrl"`
}

type wellKnownLinksDomain struct {
	Domain string          `json:"domain"`
	Count  int             `json:"count"`
	Links  []wellKnownLink `json:"links"`
}

type wellKnownInternalTarget struct {
	TargetURL string          `json:"targetUrl"`
	Count     int             `json:"count"`
	Links     []wellKnownLink `json:"links"`
}

func (d wellKnownData) toMap() map[string]interface{} {
	return map[string]interface{}{
		"site_url":          d.SiteURL,
		"site_title":        d.SiteTitle,
		"site_description":  d.SiteDescription,
		"site_host":         d.SiteHost,
		"author":            d.Author,
		"author_image_url":  d.AuthorImageURL,
		"build_time":        d.BuildTime,
		"webfinger_subject": d.WebfingerSubject,
		"webfinger_url":     d.WebfingerURL,
		"webfinger_aliases": d.WebfingerAliases,
		"nodeinfo_url":      d.NodeInfoURL,
		"generator":         d.GeneratorName,
		"generator_version": d.GeneratorVersion,
		"ssh_fingerprint":   d.SSHFingerprint,
		"keybase_username":  d.KeybaseUsername,
	}
}

type wellKnownEntry struct {
	name     string
	path     string
	template string
	fallback func(wellKnownData) string
}

// WellKnownPlugin generates .well-known endpoints from site metadata.
type WellKnownPlugin struct {
	now func() time.Time
}

// NewWellKnownPlugin creates a new WellKnownPlugin.
func NewWellKnownPlugin() *WellKnownPlugin {
	return &WellKnownPlugin{now: time.Now}
}

// Name returns the unique name of the plugin.
func (p *WellKnownPlugin) Name() string {
	return "well_known"
}

// Write generates .well-known files during the write stage.
func (p *WellKnownPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	wellKnownConfig := getWellKnownConfig(config)
	if !wellKnownConfig.IsEnabled() {
		return nil
	}

	if err := os.MkdirAll(config.OutputDir, 0o755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	var engine *templates.Engine
	if cached, ok := m.Cache().Get("templates.engine"); ok && cached != nil {
		if e, ok := cached.(*templates.Engine); ok {
			engine = e
		}
	}

	data := buildWellKnownData(config, wellKnownConfig, p.now())
	data.Links = buildWellKnownLinks(m.Posts())
	data.InternalLinks = buildWellKnownInternalLinks(m.Posts())
	entries := resolveWellKnownEntries(wellKnownConfig)

	// Add avatar endpoint if author image is available
	if data.AuthorImageURL != "" {
		entries = append(entries, wellKnownEntries["avatar"])
	}

	for _, entry := range entries {
		content, err := p.renderEntry(entry, engine, config, m, data)
		if err != nil {
			return err
		}
		// Skip empty content (e.g., avatar with no image)
		if content == "" || content == "\n" {
			continue
		}
		if err := writeWellKnownFile(config.OutputDir, entry.path, content); err != nil {
			return err
		}
	}

	return nil
}

func getWellKnownConfig(config *lifecycle.Config) models.WellKnownConfig {
	if config != nil && config.Extra != nil {
		if wk, ok := config.Extra["well_known"].(models.WellKnownConfig); ok {
			return wk
		}
	}
	return models.NewWellKnownConfig()
}

func buildWellKnownData(config *lifecycle.Config, wellKnownConfig models.WellKnownConfig, buildTime time.Time) wellKnownData {
	siteURL := strings.TrimSuffix(getSiteURL(config), "/")
	if siteURL == "" {
		siteURL = DefaultSiteURL
	}

	siteTitle := getSiteTitle(config)
	siteDescription := getSiteDescription(config)
	siteHost := extractHost(siteURL)
	author := getStringFromExtra(config.Extra, "author")
	authorImageURL := getAuthorImageURL(config, siteURL)
	webfingerURL := siteURL + "/" + wellKnownDir + "/webfinger"
	nodeInfoURL := siteURL + "/nodeinfo/2.0"

	webfingerSubject := siteURL
	if author != "" && siteHost != "" {
		webfingerSubject = fmt.Sprintf("acct:%s@%s", models.Slugify(author), siteHost)
	}

	return wellKnownData{
		SiteURL:          siteURL,
		SiteTitle:        siteTitle,
		SiteDescription:  siteDescription,
		SiteHost:         siteHost,
		Author:           author,
		AuthorImageURL:   authorImageURL,
		BuildTime:        buildTime.UTC().Format(time.RFC3339),
		WebfingerSubject: webfingerSubject,
		WebfingerURL:     webfingerURL,
		WebfingerAliases: []string{siteURL},
		NodeInfoURL:      nodeInfoURL,
		GeneratorName:    wellKnownGeneratorName,
		GeneratorVersion: wellKnownGeneratorVer,
		SSHFingerprint:   wellKnownConfig.SSHFingerprint,
		KeybaseUsername:  wellKnownConfig.KeybaseUsername,
	}
}

func extractHost(siteURL string) string {
	parsed, err := url.Parse(siteURL)
	if err != nil {
		return ""
	}
	host := parsed.Hostname()
	if host == "" {
		return strings.TrimSpace(siteURL)
	}
	return host
}

// getAuthorImageURL returns the absolute URL of the author's avatar image.
// It checks seo.author_image first, falling back to seo.default_image.
// Relative paths are resolved against the site URL.
func getAuthorImageURL(config *lifecycle.Config, siteURL string) string {
	if config == nil || config.Extra == nil {
		return ""
	}

	// Try to get SEO config
	seoVal, ok := config.Extra["seo"]
	if !ok {
		return ""
	}

	seoMap, ok := seoVal.(map[string]interface{})
	if !ok {
		// Try as SEOConfig struct
		if seo, ok := seoVal.(models.SEOConfig); ok {
			imageURL := seo.AuthorImage
			if imageURL == "" {
				imageURL = seo.DefaultImage
			}
			return resolveImageURL(imageURL, siteURL)
		}
		return ""
	}

	// Get author_image, fallback to default_image
	imageURL := ""
	if authorImage, ok := seoMap["author_image"].(string); ok && authorImage != "" {
		imageURL = authorImage
	} else if defaultImage, ok := seoMap["default_image"].(string); ok && defaultImage != "" {
		imageURL = defaultImage
	}

	return resolveImageURL(imageURL, siteURL)
}

// resolveImageURL converts a potentially relative image URL to an absolute URL.
func resolveImageURL(imageURL, siteURL string) string {
	if imageURL == "" {
		return ""
	}

	// Already absolute
	if strings.HasPrefix(imageURL, "http://") || strings.HasPrefix(imageURL, "https://") {
		return imageURL
	}

	// Relative path - resolve against site URL
	if strings.HasPrefix(imageURL, "/") {
		return strings.TrimSuffix(siteURL, "/") + imageURL
	}

	// Relative without leading slash
	return strings.TrimSuffix(siteURL, "/") + "/" + imageURL
}

func resolveWellKnownEntries(wellKnownConfig models.WellKnownConfig) []wellKnownEntry {
	entries := make([]wellKnownEntry, 0)
	seen := make(map[string]bool)

	addEntry := func(entry wellKnownEntry) {
		if entry.name == "" || seen[entry.name] {
			return
		}
		entries = append(entries, entry)
		seen[entry.name] = true
	}

	for _, name := range wellKnownConfig.AutoGenerateList() {
		switch name {
		case "nodeinfo":
			addEntry(wellKnownEntries["nodeinfo"])
			addEntry(wellKnownEntries["nodeinfo-2.0"])
		case "links":
			addEntry(wellKnownEntries["links"])
			addEntry(wellKnownEntries["internal-links"])
			addEntry(wellKnownEntries["external-links"])
			addEntry(wellKnownEntries["internal-links-page"])
		default:
			if entry, ok := wellKnownEntries[name]; ok {
				addEntry(entry)
			}
		}
	}

	if wellKnownConfig.SSHFingerprint != "" {
		addEntry(wellKnownEntries["sshfp"])
	}
	if wellKnownConfig.KeybaseUsername != "" {
		addEntry(wellKnownEntries["keybase"])
	}

	return entries
}

func (p *WellKnownPlugin) renderEntry(entry wellKnownEntry, engine *templates.Engine, config *lifecycle.Config, m *lifecycle.Manager, data wellKnownData) (string, error) {
	if engine != nil && engine.TemplateExists(entry.template) {
		ctx := templates.NewContext(nil, "", ToModelsConfig(config))
		ctx.Extra["well_known"] = data.toMap()
		ctx.Extra["well_known_links"] = data.Links
		ctx.Extra["well_known_internal_links"] = data.InternalLinks
		ctx = ctx.WithCore(m)
		result, err := engine.Render(entry.template, ctx)
		if err == nil {
			return normalizeWellKnownContent(result), nil
		}
	}

	if entry.fallback == nil {
		return "", fmt.Errorf("no template or fallback for %s", entry.name)
	}

	return normalizeWellKnownContent(entry.fallback(data)), nil
}

func buildWellKnownLinks(posts []*models.Post) []wellKnownLinksDomain {
	byDomain := make(map[string][]wellKnownLink)
	seen := make(map[string]struct{})

	for _, post := range posts {
		if post == nil {
			continue
		}
		for _, link := range post.Outlinks {
			if link == nil || link.IsInternal {
				continue
			}

			targetURL := strings.TrimSpace(link.TargetURL)
			if targetURL == "" {
				continue
			}
			targetParsed, err := url.Parse(targetURL)
			if err != nil {
				continue
			}
			domain := strings.ToLower(strings.TrimSpace(targetParsed.Hostname()))
			if domain == "" {
				continue
			}

			sourceURL := strings.TrimSpace(link.SourceURL)
			if sourceURL == "" && post.Href != "" {
				sourceURL = post.Href
			}

			dedupeKey := domain + "|" + sourceURL + "|" + targetURL
			if _, ok := seen[dedupeKey]; ok {
				continue
			}
			seen[dedupeKey] = struct{}{}

			byDomain[domain] = append(byDomain[domain], wellKnownLink{
				SourceURL: sourceURL,
				TargetURL: targetURL,
			})
		}
	}

	result := make([]wellKnownLinksDomain, 0, len(byDomain))
	for domain, links := range byDomain {
		sort.Slice(links, func(i, j int) bool {
			if links[i].SourceURL == links[j].SourceURL {
				return links[i].TargetURL < links[j].TargetURL
			}
			return links[i].SourceURL < links[j].SourceURL
		})
		result = append(result, wellKnownLinksDomain{
			Domain: domain,
			Count:  len(links),
			Links:  links,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Domain < result[j].Domain
		}
		return result[i].Count > result[j].Count
	})

	return result
}

func buildWellKnownInternalLinks(posts []*models.Post) []wellKnownInternalTarget {
	byTarget := make(map[string][]wellKnownLink)
	seen := make(map[string]struct{})

	for _, post := range posts {
		if post == nil {
			continue
		}
		for _, link := range post.Outlinks {
			if link == nil || !link.IsInternal {
				continue
			}

			targetURL := strings.TrimSpace(link.TargetURL)
			if targetURL == "" {
				continue
			}

			normalizedTarget := normalizeInternalTarget(targetURL)
			if normalizedTarget == "" {
				continue
			}

			sourceURL := strings.TrimSpace(link.SourceURL)
			if sourceURL == "" && post.Href != "" {
				sourceURL = post.Href
			}

			dedupeKey := normalizedTarget + "|" + sourceURL + "|" + targetURL
			if _, ok := seen[dedupeKey]; ok {
				continue
			}
			seen[dedupeKey] = struct{}{}

			byTarget[normalizedTarget] = append(byTarget[normalizedTarget], wellKnownLink{
				SourceURL: sourceURL,
				TargetURL: targetURL,
			})
		}
	}

	result := make([]wellKnownInternalTarget, 0, len(byTarget))
	for targetURL, links := range byTarget {
		sort.Slice(links, func(i, j int) bool {
			if links[i].SourceURL == links[j].SourceURL {
				return links[i].TargetURL < links[j].TargetURL
			}
			return links[i].SourceURL < links[j].SourceURL
		})
		result = append(result, wellKnownInternalTarget{
			TargetURL: targetURL,
			Count:     len(links),
			Links:     links,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].TargetURL < result[j].TargetURL
		}
		return result[i].Count > result[j].Count
	})

	return result
}

func normalizeInternalTarget(targetURL string) string {
	trimmed := strings.TrimSpace(targetURL)
	if trimmed == "" {
		return ""
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return trimmed
	}

	path := parsed.Path
	if path == "" {
		path = trimmed
	}
	if strings.HasSuffix(path, "/") && path != "/" {
		path = strings.TrimSuffix(path, "/")
	}
	if parsed.RawQuery != "" {
		path = path + "?" + parsed.RawQuery
	}

	return path
}

func marshalWellKnownLinks(links []wellKnownLinksDomain) string {
	if links == nil {
		links = []wellKnownLinksDomain{}
	}
	b, err := json.Marshal(links)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func marshalWellKnownInternalLinks(links []wellKnownInternalTarget) string {
	if links == nil {
		links = []wellKnownInternalTarget{}
	}
	b, err := json.Marshal(links)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func renderExternalLinksHTML(data wellKnownData) string {
	title := "External Links"
	if data.SiteTitle != "" {
		title = data.SiteTitle + " External Links"
	}

	totalLinks := 0
	maxCount := 0
	for _, domain := range data.Links {
		totalLinks += domain.Count
		if domain.Count > maxCount {
			maxCount = domain.Count
		}
	}

	var builder strings.Builder
	builder.WriteString("<!doctype html>\n")
	builder.WriteString("<html lang=\"en\">\n<head>\n")
	builder.WriteString("  <meta charset=\"utf-8\">\n")
	builder.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	builder.WriteString("  <title>" + html.EscapeString(title) + "</title>\n")
	builder.WriteString("  <style>\n")
	builder.WriteString("    :root{color-scheme:light dark;--accent:#2563eb;--bar:#93c5fd;--bar-bg:#e5e7eb44;}body{font-family:ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,sans-serif;max-width:70rem;margin:2rem auto;padding:0 1rem;line-height:1.5;}h1{margin-bottom:.25rem;}p{margin-top:0;color:#666;}section{border:1px solid #ccc3;border-radius:.6rem;padding:1rem;margin:1rem 0;}ul{padding-left:1.2rem;margin:.5rem 0 0;}li+li{margin-top:.35rem;}small{color:#777;}code{font-size:.9em;}a{color:inherit;}\n")
	builder.WriteString("    .external-links-chart{margin:1.25rem 0 2rem;}\n")
	builder.WriteString("    .external-links-chart ol{list-style:decimal;margin:0;padding-left:1.5rem;}\n")
	builder.WriteString("    .external-links-chart li{margin:.35rem 0;}\n")
	builder.WriteString("    .external-links-row{display:grid;grid-template-columns:minmax(12rem,24rem) minmax(8rem,1fr) auto;align-items:center;gap:.75rem;}\n")
	builder.WriteString("    .external-links-domain{display:flex;align-items:center;gap:.5rem;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;}\n")
	builder.WriteString("    .external-links-domain img{width:16px;height:16px;border-radius:3px;flex:none;}\n")
	builder.WriteString("    .external-links-bar{height:.65rem;background:var(--bar-bg);border-radius:999px;overflow:hidden;}\n")
	builder.WriteString("    .external-links-bar-fill{display:block;height:100%;background:linear-gradient(90deg,var(--bar),var(--accent));border-radius:999px;}\n")
	builder.WriteString("    .external-links-count{font-variant-numeric:tabular-nums;min-width:2.5rem;text-align:right;}\n")
	builder.WriteString("    .external-links-links{margin-top:.75rem;}\n")
	builder.WriteString("    @media (max-width:720px){.external-links-row{grid-template-columns:1fr auto;}.external-links-bar{grid-column:1 / -1;}}\n")
	builder.WriteString("  </style>\n")
	builder.WriteString("</head>\n<body>\n")
	builder.WriteString("  <header>\n")
	builder.WriteString("    <h1>External Links</h1>\n")
	builder.WriteString("    <p>")
	builder.WriteString(fmt.Sprintf("%d links across %d domains", totalLinks, len(data.Links)))
	if data.SiteURL != "" {
		builder.WriteString(" from <a href=\"")
		builder.WriteString(html.EscapeString(withTrailingSlash(data.SiteURL)))
		builder.WriteString("\">")
		builder.WriteString(html.EscapeString(withTrailingSlash(data.SiteURL)))
		builder.WriteString("</a>")
	}
	builder.WriteString(". JSON source: <a href=\"/.well-known/links\"><code>/.well-known/links</code></a>.</p>\n")
	builder.WriteString("  </header>\n")
	builder.WriteString("  <section class=\"external-links-chart\">\n")
	builder.WriteString("    <h2>Domains</h2>\n")
	builder.WriteString("    <ol>\n")
	for _, domain := range data.Links {
		barWidth := 0
		if maxCount > 0 {
			barWidth = (domain.Count * 100) / maxCount
		}
		if barWidth < 1 && domain.Count > 0 {
			barWidth = 1
		}
		faviconURL := "https://www.google.com/s2/favicons?domain=" + url.QueryEscape(domain.Domain) + "&sz=32"

		builder.WriteString("      <li><div class=\"external-links-row\">\n")
		builder.WriteString("        <a class=\"external-links-domain\" href=\"#domain-")
		builder.WriteString(html.EscapeString(domain.Domain))
		builder.WriteString("\">\n")
		builder.WriteString("          <img src=\"")
		builder.WriteString(html.EscapeString(faviconURL))
		builder.WriteString("\" alt=\"\">\n")
		builder.WriteString("          <span>")
		builder.WriteString(html.EscapeString(domain.Domain))
		builder.WriteString("</span>\n")
		builder.WriteString("        </a>\n")
		builder.WriteString("        <span class=\"external-links-bar\" aria-hidden=\"true\"><span class=\"external-links-bar-fill\" style=\"width:")
		builder.WriteString(fmt.Sprintf("%d", barWidth))
		builder.WriteString("%\"></span></span>\n")
		builder.WriteString("        <span class=\"external-links-count\">")
		builder.WriteString(fmt.Sprintf("%d", domain.Count))
		builder.WriteString("</span>\n")
		builder.WriteString("      </div></li>\n")
	}
	builder.WriteString("    </ol>\n")
	builder.WriteString("  </section>\n")

	for _, domain := range data.Links {
		builder.WriteString("  <section>\n")
		builder.WriteString("    <h2 id=\"domain-" + html.EscapeString(domain.Domain) + "\">" + html.EscapeString(domain.Domain) + " <small>(" + fmt.Sprintf("%d", domain.Count) + ")</small></h2>\n")
		builder.WriteString("    <ul class=\"external-links-links\">\n")
		for _, link := range domain.Links {
			source := link.SourceURL
			if source == "" {
				source = "(unknown source)"
			}
			builder.WriteString("      <li><a href=\"")
			builder.WriteString(html.EscapeString(link.TargetURL))
			builder.WriteString("\">")
			builder.WriteString(html.EscapeString(link.TargetURL))
			builder.WriteString("</a> <small>from ")
			if link.SourceURL != "" {
				builder.WriteString("<a href=\"")
				builder.WriteString(html.EscapeString(link.SourceURL))
				builder.WriteString("\">")
				builder.WriteString(html.EscapeString(source))
				builder.WriteString("</a>")
			} else {
				builder.WriteString(html.EscapeString(source))
			}
			builder.WriteString("</small></li>\n")
		}
		builder.WriteString("    </ul>\n")
		builder.WriteString("  </section>\n")
	}

	builder.WriteString("</body>\n</html>\n")

	return builder.String()
}

func renderInternalLinksHTML(data wellKnownData) string {
	title := "Internal Links"
	if data.SiteTitle != "" {
		title = data.SiteTitle + " Internal Links"
	}

	totalLinks := 0
	maxCount := 0
	for _, target := range data.InternalLinks {
		totalLinks += target.Count
		if target.Count > maxCount {
			maxCount = target.Count
		}
	}

	var builder strings.Builder
	builder.WriteString("<!doctype html>\n")
	builder.WriteString("<html lang=\"en\">\n<head>\n")
	builder.WriteString("  <meta charset=\"utf-8\">\n")
	builder.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	builder.WriteString("  <title>" + html.EscapeString(title) + "</title>\n")
	builder.WriteString("  <style>\n")
	builder.WriteString("    :root{color-scheme:light dark;--accent:#059669;--bar:#6ee7b7;--bar-bg:#e5e7eb44;}body{font-family:ui-sans-serif,system-ui,-apple-system,Segoe UI,Roboto,sans-serif;max-width:70rem;margin:2rem auto;padding:0 1rem;line-height:1.5;}h1{margin-bottom:.25rem;}p{margin-top:0;color:#666;}section{border:1px solid #ccc3;border-radius:.6rem;padding:1rem;margin:1rem 0;}ul{padding-left:1.2rem;margin:.5rem 0 0;}li+li{margin-top:.35rem;}small{color:#777;}code{font-size:.9em;}a{color:inherit;}\n")
	builder.WriteString("    .internal-links-chart{margin:1.25rem 0 2rem;}\n")
	builder.WriteString("    .internal-links-chart ol{list-style:decimal;margin:0;padding-left:1.5rem;}\n")
	builder.WriteString("    .internal-links-chart li{margin:.35rem 0;}\n")
	builder.WriteString("    .internal-links-row{display:grid;grid-template-columns:minmax(16rem,26rem) minmax(8rem,1fr) auto;align-items:center;gap:.75rem;}\n")
	builder.WriteString("    .internal-links-target{white-space:nowrap;overflow:hidden;text-overflow:ellipsis;}\n")
	builder.WriteString("    .internal-links-bar{height:.65rem;background:var(--bar-bg);border-radius:999px;overflow:hidden;}\n")
	builder.WriteString("    .internal-links-bar-fill{display:block;height:100%;background:linear-gradient(90deg,var(--bar),var(--accent));border-radius:999px;}\n")
	builder.WriteString("    .internal-links-count{font-variant-numeric:tabular-nums;min-width:2.5rem;text-align:right;}\n")
	builder.WriteString("    .internal-links-links{margin-top:.75rem;}\n")
	builder.WriteString("    @media (max-width:720px){.internal-links-row{grid-template-columns:1fr auto;}.internal-links-bar{grid-column:1 / -1;}}\n")
	builder.WriteString("  </style>\n")
	builder.WriteString("</head>\n<body>\n")
	builder.WriteString("  <header>\n")
	builder.WriteString("    <h1>Internal Links</h1>\n")
	builder.WriteString("    <p>")
	builder.WriteString(fmt.Sprintf("%d links across %d destinations", totalLinks, len(data.InternalLinks)))
	if data.SiteURL != "" {
		builder.WriteString(" on <a href=\"")
		builder.WriteString(html.EscapeString(withTrailingSlash(data.SiteURL)))
		builder.WriteString("\">")
		builder.WriteString(html.EscapeString(withTrailingSlash(data.SiteURL)))
		builder.WriteString("</a>")
	}
	builder.WriteString(". JSON source: <a href=\"/.well-known/internal-links\"><code>/.well-known/internal-links</code></a>.</p>\n")
	builder.WriteString("  </header>\n")
	builder.WriteString("  <section class=\"internal-links-chart\">\n")
	builder.WriteString("    <h2>Destinations</h2>\n")
	builder.WriteString("    <ol>\n")
	for idx, target := range data.InternalLinks {
		barWidth := 0
		if maxCount > 0 {
			barWidth = (target.Count * 100) / maxCount
		}
		if barWidth < 1 && target.Count > 0 {
			barWidth = 1
		}

		builder.WriteString("      <li><div class=\"internal-links-row\">\n")
		builder.WriteString("        <a class=\"internal-links-target\" href=\"#target-")
		builder.WriteString(fmt.Sprintf("%d", idx+1))
		builder.WriteString("\">")
		builder.WriteString(html.EscapeString(target.TargetURL))
		builder.WriteString("</a>\n")
		builder.WriteString("        <span class=\"internal-links-bar\" aria-hidden=\"true\"><span class=\"internal-links-bar-fill\" style=\"width:")
		builder.WriteString(fmt.Sprintf("%d", barWidth))
		builder.WriteString("%\"></span></span>\n")
		builder.WriteString("        <span class=\"internal-links-count\">")
		builder.WriteString(fmt.Sprintf("%d", target.Count))
		builder.WriteString("</span>\n")
		builder.WriteString("      </div></li>\n")
	}
	builder.WriteString("    </ol>\n")
	builder.WriteString("  </section>\n")

	for idx, target := range data.InternalLinks {
		builder.WriteString("  <section>\n")
		builder.WriteString("    <h2 id=\"target-")
		builder.WriteString(fmt.Sprintf("%d", idx+1))
		builder.WriteString("\">")
		builder.WriteString(html.EscapeString(target.TargetURL))
		builder.WriteString(" <small>(")
		builder.WriteString(fmt.Sprintf("%d", target.Count))
		builder.WriteString(")</small></h2>\n")
		builder.WriteString("    <ul class=\"internal-links-links\">\n")
		for _, link := range target.Links {
			source := link.SourceURL
			if source == "" {
				source = "(unknown source)"
			}
			builder.WriteString("      <li>")
			if link.SourceURL != "" {
				builder.WriteString("<a href=\"")
				builder.WriteString(html.EscapeString(link.SourceURL))
				builder.WriteString("\">")
				builder.WriteString(html.EscapeString(source))
				builder.WriteString("</a>")
			} else {
				builder.WriteString(html.EscapeString(source))
			}
			builder.WriteString(" <small>to <a href=\"")
			builder.WriteString(html.EscapeString(link.TargetURL))
			builder.WriteString("\">")
			builder.WriteString(html.EscapeString(link.TargetURL))
			builder.WriteString("</a></small></li>\n")
		}
		builder.WriteString("    </ul>\n")
		builder.WriteString("  </section>\n")
	}

	builder.WriteString("</body>\n</html>\n")

	return builder.String()
}

func writeWellKnownFile(outputDir, relativePath, content string) error {
	fullPath := filepath.Join(outputDir, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("creating directory for %s: %w", fullPath, err)
	}
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", fullPath, err)
	}
	return nil
}

func normalizeWellKnownContent(content string) string {
	trimmed := strings.TrimRight(content, "\n")
	return trimmed + "\n"
}

func withTrailingSlash(value string) string {
	if value == "" {
		return value
	}
	if strings.HasSuffix(value, "/") {
		return value
	}
	return value + "/"
}

var wellKnownEntries = map[string]wellKnownEntry{
	"host-meta": {
		name:     "host-meta",
		path:     wellKnownDir + "/host-meta",
		template: "well-known/host-meta.xml",
		fallback: func(data wellKnownData) string {
			return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">
  <Subject>%s</Subject>
  <Link rel="%s" type="%s" template="%s?resource={uri}" />
</XRD>`, withTrailingSlash(data.SiteURL), wellKnownHostMetaRel, wellKnownHostMetaType, data.WebfingerURL)
		},
	},
	"host-meta.json": {
		name:     "host-meta.json",
		path:     wellKnownDir + "/host-meta.json",
		template: "well-known/host-meta.json",
		fallback: func(data wellKnownData) string {
			return fmt.Sprintf(`{
  "subject": "%s",
  "links": [
    {
      "rel": "%s",
      "type": "%s",
      "template": "%s?resource={uri}"
    }
  ]
}`, withTrailingSlash(data.SiteURL), wellKnownHostMetaRel, wellKnownHostMetaType, data.WebfingerURL)
		},
	},
	"webfinger": {
		name:     "webfinger",
		path:     wellKnownDir + "/webfinger",
		template: "well-known/webfinger.json",
		fallback: func(data wellKnownData) string {
			// Build links array - always include profile-page
			links := fmt.Sprintf(`{
      "rel": "%s",
      "type": "%s",
      "href": "%s"
    }`, wellKnownWebfingerRel, wellKnownWebfingerType, withTrailingSlash(data.SiteURL))

			// Add avatar link if author image is available
			if data.AuthorImageURL != "" {
				links += fmt.Sprintf(`,
    {
      "rel": "%s",
      "href": "%s"
    }`, wellKnownAvatarRel, data.AuthorImageURL)
			}

			return fmt.Sprintf(`{
  "subject": "%s",
  "aliases": [
    "%s"
  ],
  "links": [
    %s
  ]
}`, data.WebfingerSubject, withTrailingSlash(data.SiteURL), links)
		},
	},
	"nodeinfo": {
		name:     "nodeinfo",
		path:     wellKnownDir + "/nodeinfo",
		template: "well-known/nodeinfo.json",
		fallback: func(data wellKnownData) string {
			return fmt.Sprintf(`{
  "links": [
    {
      "rel": "%s",
      "href": "%s"
    }
  ]
}`, wellKnownNodeInfoRel, data.NodeInfoURL)
		},
	},
	"nodeinfo-2.0": {
		name:     "nodeinfo-2.0",
		path:     "nodeinfo/2.0",
		template: "well-known/nodeinfo-2.0.json",
		fallback: func(data wellKnownData) string {
			return fmt.Sprintf(`{
  "version": "2.0",
  "software": {
    "name": "%s",
    "version": "%s"
  },
  "protocols": [],
  "services": {
    "inbound": [],
    "outbound": []
  },
  "openRegistrations": false,
  "usage": {
    "users": {
      "total": 0
    }
  },
  "metadata": {
    "site": {
      "name": "%s",
      "description": "%s",
      "url": "%s"
    }
  }
}`, data.GeneratorName, data.GeneratorVersion, data.SiteTitle, data.SiteDescription, withTrailingSlash(data.SiteURL))
		},
	},
	"time": {
		name:     "time",
		path:     wellKnownDir + "/time",
		template: "well-known/time.txt",
		fallback: func(data wellKnownData) string {
			return data.BuildTime
		},
	},
	"links": {
		name:     "links",
		path:     wellKnownDir + "/links",
		template: "well-known/links.json",
		fallback: func(data wellKnownData) string {
			return marshalWellKnownLinks(data.Links)
		},
	},
	"internal-links": {
		name:     "internal-links",
		path:     wellKnownDir + "/internal-links",
		template: "well-known/internal-links.json",
		fallback: func(data wellKnownData) string {
			return marshalWellKnownInternalLinks(data.InternalLinks)
		},
	},
	"external-links": {
		name:     "external-links",
		path:     "external-links/index.html",
		template: "external-links.html",
		fallback: renderExternalLinksHTML,
	},
	"internal-links-page": {
		name:     "internal-links-page",
		path:     "internal-links/index.html",
		template: "internal-links.html",
		fallback: renderInternalLinksHTML,
	},
	"sshfp": {
		name:     "sshfp",
		path:     wellKnownDir + "/sshfp",
		template: "well-known/sshfp.txt",
		fallback: func(data wellKnownData) string {
			return data.SSHFingerprint
		},
	},
	"keybase": {
		name:     "keybase",
		path:     wellKnownDir + "/keybase.txt",
		template: "well-known/keybase.txt",
		fallback: func(data wellKnownData) string {
			return fmt.Sprintf("keybase: %s", data.KeybaseUsername)
		},
	},
	"avatar": {
		name:     "avatar",
		path:     wellKnownDir + "/avatar",
		template: "well-known/avatar",
		fallback: func(data wellKnownData) string {
			// Return an HTML redirect to the author image
			// This works better than raw image bytes for static hosting
			if data.AuthorImageURL == "" {
				return ""
			}
			return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta http-equiv="refresh" content="0;url=%s">
<link rel="canonical" href="%s">
</head>
<body>
<a href="%s">Avatar</a>
</body>
</html>`, data.AuthorImageURL, data.AuthorImageURL, data.AuthorImageURL)
		},
	},
}

// Ensure WellKnownPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin      = (*WellKnownPlugin)(nil)
	_ lifecycle.WritePlugin = (*WellKnownPlugin)(nil)
)
