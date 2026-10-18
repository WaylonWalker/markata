// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"

	"github.com/PuerkitoBio/goquery"
	"github.com/andybalholm/cascadia"
)

// LinkAvatarsConfig holds configuration for the link_avatars plugin.
type LinkAvatarsConfig struct {
	// Enabled controls whether the plugin is active.
	// Default: false
	Enabled bool

	// Mode controls how avatars are applied: "js", "local", or "hosted".
	// Default: "js"
	Mode string

	// Selector is the CSS selector for links to enhance.
	// Default: "a[href^='http']"
	Selector string

	// Service is the avatar service provider: "duckduckgo", "google", "custom"
	// Default: "duckduckgo"
	Service string

	// Template is a custom URL template (only used when Service = "custom").
	// Supports placeholders: {origin}, {host}
	Template string

	// IgnoreDomains is a list of domains to skip.
	IgnoreDomains []string

	// IgnoreOrigins is a list of full origins to skip (includes protocol).
	IgnoreOrigins []string

	// IgnoreSelectors is a list of CSS selectors to exclude.
	IgnoreSelectors []string

	// IgnoreClasses is a list of CSS classes to exclude.
	IgnoreClasses []string

	// IgnoreIDs is a list of element IDs to exclude.
	IgnoreIDs []string

	// Size is the avatar icon size in pixels.
	// Default: 16
	Size int

	// Position is where to place the avatar: "before" or "after" link text.
	// Default: "before"
	Position string

	// HostedBaseURL is the base URL for hosted mode assets.
	// Used when Mode = "hosted".
	HostedBaseURL string
}

const (
	linkAvatarModeJS     = "js"
	linkAvatarModeLocal  = "local"
	linkAvatarModeHosted = "hosted"

	linkAvatarIconExtICO = ".ico"
	linkAvatarMissExt    = ".missing"
)

// defaultLinkAvatarsConfig returns the default configuration.
func defaultLinkAvatarsConfig() LinkAvatarsConfig {
	return LinkAvatarsConfig{
		Enabled:         true,
		Mode:            linkAvatarModeJS,
		Selector:        "a[href^='http']",
		Service:         "duckduckgo",
		Template:        "",
		IgnoreDomains:   []string{},
		IgnoreOrigins:   []string{},
		IgnoreSelectors: []string{},
		IgnoreClasses:   []string{"no-avatar"},
		IgnoreIDs:       []string{},
		Size:            16,
		Position:        "before",
		HostedBaseURL:   "",
	}
}

// LinkAvatarsPlugin adds favicon/avatar icons next to external links.
// It generates client-side JavaScript and CSS assets that enhance links
// at runtime in the browser.
type LinkAvatarsPlugin struct {
	config      LinkAvatarsConfig
	siteOrigin  string
	siteURLPath string
	client      *http.Client
	cssHash     string
	jsHash      string
}

// NewLinkAvatarsPlugin creates a new LinkAvatarsPlugin.
func NewLinkAvatarsPlugin() *LinkAvatarsPlugin {
	client := &http.Client{Timeout: 10 * time.Second}
	buildstats.InstrumentHTTPClient(client)

	return &LinkAvatarsPlugin{
		config: defaultLinkAvatarsConfig(),
		client: client,
	}
}

// Name returns the unique name of the plugin.
func (p *LinkAvatarsPlugin) Name() string {
	return "link_avatars"
}

// Configure loads plugin configuration from the manager.
func (p *LinkAvatarsPlugin) Configure(m *lifecycle.Manager) error {
	p.config = parseLinkAvatarsConfig(m.Config())
	if err := p.validateConfig(); err != nil {
		return err
	}

	p.siteOrigin = getSiteOrigin(m.Config())
	p.siteURLPath = getSiteURLPath(m.Config())

	if p.config.Enabled {
		p.computeHashes()
		p.injectHeadTags(m.Config())
	}

	return nil
}

// computeHashes computes content hashes for CSS and JS assets.
func (p *LinkAvatarsPlugin) computeHashes() {
	cssContent := p.generateCSS()
	p.cssHash = fmt.Sprintf("%x", sha256.Sum256([]byte(cssContent)))[:8]

	if p.config.Mode == linkAvatarModeJS {
		jsContent := p.generateJavaScript()
		p.jsHash = fmt.Sprintf("%x", sha256.Sum256([]byte(jsContent)))[:8]
	}
}

// Write generates the JavaScript and CSS assets.
func (p *LinkAvatarsPlugin) Write(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	cfg := m.Config()
	outputDir := cfg.OutputDir
	if outputDir == "" {
		outputDir = defaultOutputDir
	}

	cssDir := filepath.Join(outputDir, "css")
	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		return fmt.Errorf("creating link_avatars css directory: %w", err)
	}

	cssContent := p.generateCSS()
	cssHash := fmt.Sprintf("%x", sha256.Sum256([]byte(cssContent)))[:8]
	p.cssHash = cssHash

	cssPath := filepath.Join(cssDir, "link-avatars.css")
	if err := os.WriteFile(cssPath, []byte(cssContent), 0o644); err != nil { //nolint:gosec // static CSS needs world-readable permissions
		return fmt.Errorf("writing link-avatars.css: %w", err)
	}

	m.SetAssetHash("css/link-avatars.css", cssHash)
	templates.SetAssetHashes(map[string]string{"css/link-avatars.css": cssHash})

	if p.config.Mode == linkAvatarModeJS {
		jsContent := p.generateJavaScript()
		jsHash := fmt.Sprintf("%x", sha256.Sum256([]byte(jsContent)))[:8]
		p.jsHash = jsHash

		jsDir := filepath.Join(outputDir, "js")
		if err := os.MkdirAll(jsDir, 0o755); err != nil {
			return fmt.Errorf("creating link_avatars js directory: %w", err)
		}

		jsPath := filepath.Join(jsDir, "link-avatars.js")
		if err := os.WriteFile(jsPath, []byte(jsContent), 0o644); err != nil { //nolint:gosec // static JS needs world-readable permissions
			return fmt.Errorf("writing link-avatars.js: %w", err)
		}

		m.SetAssetHash("js/link-avatars.js", jsHash)
		templates.SetAssetHashes(map[string]string{"js/link-avatars.js": jsHash})
	}

	return nil
}

// Render injects build-time avatars for local/hosted modes.
func (p *LinkAvatarsPlugin) Render(m *lifecycle.Manager) error {
	if !p.config.Enabled || p.config.Mode == linkAvatarModeJS {
		return nil
	}

	outputDir := resolveOutputDir(m.Config())
	assetsDir := filepath.Join(outputDir, "assets", "markata", "link-avatars")
	if err := os.MkdirAll(assetsDir, 0o755); err != nil {
		return fmt.Errorf("creating link_avatars icon directory: %w", err)
	}

	publicBase, err := p.iconBaseURL()
	if err != nil {
		return err
	}

	cache := GetBuildCache(m)

	// Thread-safe icon cache for concurrent access
	var iconCache sync.Map

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.ArticleHTML == "" {
			return false
		}
		return true
	})

	if lifecycle.IsServeIncremental(m) {
		if affected := lifecycle.GetServeAffectedPaths(m); len(affected) > 0 {
			filtered := posts[:0]
			for _, post := range posts {
				if affected[post.Path] {
					filtered = append(filtered, post)
				}
			}
			posts = filtered
		}
	}

	// Phase 1: Restore cached results for unchanged posts
	var needProcessing []*models.Post
	if cache != nil {
		for _, post := range posts {
			articleHash := buildcache.ContentHash(post.ArticleHTML)
			if cached, ok := cache.GetCachedLinkAvatarsHTML(post.Path, articleHash); ok {
				post.ArticleHTML = cached
			} else {
				needProcessing = append(needProcessing, post)
			}
		}
	} else {
		needProcessing = posts
	}

	if len(needProcessing) == 0 {
		return nil
	}

	// Phase 2: Process posts that need updating, concurrently
	return m.ProcessPostsSliceConcurrently(needProcessing, func(post *models.Post) error {
		articleHash := buildcache.ContentHash(post.ArticleHTML)
		updated, processErr := p.processHTMLConcurrent(post.ArticleHTML, publicBase, assetsDir, &iconCache)
		if processErr != nil {
			return fmt.Errorf("link_avatars render %q: %w", post.Path, processErr)
		}
		post.ArticleHTML = updated
		if cache != nil {
			cache.CacheLinkAvatarsHTML(post.Path, articleHash, updated)
		}
		return nil
	})
}

// generateJavaScript generates the client-side JavaScript.
func (p *LinkAvatarsPlugin) generateJavaScript() string {
	// Marshal config values for JavaScript - this won't fail with basic Go types
	configJSON, err := json.Marshal(map[string]interface{}{
		"selector":        p.config.Selector,
		"service":         p.config.Service,
		"template":        p.config.Template,
		"ignoreDomains":   p.config.IgnoreDomains,
		"ignoreOrigins":   p.config.IgnoreOrigins,
		"ignoreSelectors": p.config.IgnoreSelectors,
		"ignoreClasses":   p.config.IgnoreClasses,
		"ignoreIds":       p.config.IgnoreIDs,
		"size":            p.config.Size,
		"position":        p.config.Position,
	})
	if err != nil {
		// Fallback to empty config on error (should never happen with basic types)
		configJSON = []byte("{}")
	}

	duckduckgoTemplate := "https://icons.duckduckgo.com/ip3/{host}" + linkAvatarIconExtICO

	return `/**
 * Link Avatars - markata-go
 * Adds favicon icons next to external links
 */
(function() {
  'use strict';

  var config = ` + string(configJSON) + `;

  // Service URL templates
  var serviceTemplates = {
    'duckduckgo': '` + duckduckgoTemplate + `',
    'google': 'https://www.google.com/s2/favicons?domain={host}&sz=' + config.size
  };

  function getFaviconURL(href) {
    try {
      var url = new URL(href);
      var host = url.hostname;
      var origin = url.origin;
      var template;

      if (config.service === 'custom' && config.template) {
        template = config.template;
      } else {
        template = serviceTemplates[config.service] || serviceTemplates['duckduckgo'];
      }

      return template
        .replace('{host}', host)
        .replace('{origin}', encodeURIComponent(origin));
    } catch (e) {
      return null;
    }
  }

  function shouldIgnoreLink(link) {
    var href = link.href;
    var url;

    try {
      url = new URL(href);
    } catch (e) {
      return true; // Invalid URL
    }

    // Skip same-origin links
    if (url.origin === window.location.origin) {
      return true;
    }

    // Check ignore domains
    for (var i = 0; i < config.ignoreDomains.length; i++) {
      if (url.hostname === config.ignoreDomains[i] ||
          url.hostname.endsWith('.' + config.ignoreDomains[i])) {
        return true;
      }
    }

    // Check ignore origins
    for (var j = 0; j < config.ignoreOrigins.length; j++) {
      if (url.origin === config.ignoreOrigins[j]) {
        return true;
      }
    }

    // Check ignore classes
    for (var k = 0; k < config.ignoreClasses.length; k++) {
      if (link.classList.contains(config.ignoreClasses[k])) {
        return true;
      }
    }

    // Check ignore IDs (link is inside an element with ignored ID)
    for (var l = 0; l < config.ignoreIds.length; l++) {
      if (link.closest('#' + config.ignoreIds[l])) {
        return true;
      }
    }

    // Skip links that wrap images
    if (link.querySelector('img, picture')) {
      return true;
    }

    // Check ignore selectors
    for (var m = 0; m < config.ignoreSelectors.length; m++) {
      try {
        if (link.matches(config.ignoreSelectors[m])) {
          return true;
        }
      } catch (e) {
        // Invalid selector, skip
      }
    }

    return false;
  }

  function processLink(link) {
    if (link.classList.contains('has-avatar')) {
      return; // Already processed
    }

    if (shouldIgnoreLink(link)) {
      return;
    }

    var faviconURL = getFaviconURL(link.href);
    if (!faviconURL) {
      return;
    }

    link.setAttribute('data-favicon', faviconURL);
    link.style.setProperty('--favicon-url', 'url("' + faviconURL + '")');
    link.classList.add('has-avatar');
    link.classList.add('has-avatar-' + config.position);
  }

  function processLinks() {
    var links = document.querySelectorAll(config.selector);
    links.forEach(processLink);
  }

  // Use Intersection Observer for lazy loading if available
  function observeLinks() {
    if (!('IntersectionObserver' in window)) {
      processLinks();
      return;
    }

    var observer = new IntersectionObserver(function(entries) {
      entries.forEach(function(entry) {
        if (entry.isIntersecting) {
          processLink(entry.target);
          observer.unobserve(entry.target);
        }
      });
    }, { rootMargin: '100px' });

    var links = document.querySelectorAll(config.selector);
    links.forEach(function(link) {
      if (!shouldIgnoreLink(link)) {
        observer.observe(link);
      }
    });
  }

  // Initialize
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', observeLinks);
  } else {
    observeLinks();
  }
})();
`
}

// generateCSS generates the CSS styles for link avatars.
func (p *LinkAvatarsPlugin) generateCSS() string {
	size := p.config.Size
	if size <= 0 {
		size = 16
	}

	return fmt.Sprintf(`/**
 * Link Avatars - markata-go
 * Styles for favicon icons next to external links
 */

/* Common avatar styles */
a.has-avatar {
  position: relative;
}

a.has-avatar::before,
a.has-avatar::after {
  content: '';
  display: inline-block;
  width: %dpx;
  height: %dpx;
  background-image: var(--favicon-url);
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
  vertical-align: middle;
  opacity: 0;
  transition: opacity 0.2s ease;
}

/* Show avatar when image loads */
a.has-avatar::before,
a.has-avatar::after {
  opacity: 0.8;
}

/* Avatar before link text */
a.has-avatar-before::before {
  margin-right: 0.35em;
}

a.has-avatar-before::after {
  display: none;
}

/* Avatar after link text */
a.has-avatar-after::after {
  margin-left: 0.35em;
}

a.has-avatar-after::before {
  display: none;
}

/* Hover effect */
a.has-avatar:hover::before,
a.has-avatar:hover::after {
  opacity: 1;
}

/* Hide avatar if image fails to load (fallback) */
a.has-avatar[data-favicon-error]::before,
a.has-avatar[data-favicon-error]::after {
  display: none;
}
`, size, size)
}

// injectHeadTags adds the CSS and JS references to the config head.
func (p *LinkAvatarsPlugin) injectHeadTags(cfg *lifecycle.Config) {
	if cfg == nil || cfg.Extra == nil {
		return
	}

	modelsConfig, ok := cfg.Extra["models_config"].(*models.Config)
	if !ok {
		return
	}

	cssPath := buildHashedURL("css/link-avatars.css", p.cssHash)

	modelsConfig.Head.Link = append(modelsConfig.Head.Link, models.LinkTag{
		Rel:  "stylesheet",
		Href: cssPath,
	})

	if p.config.Mode == linkAvatarModeJS {
		jsPath := buildHashedURL("js/link-avatars.js", p.jsHash)

		modelsConfig.Head.Script = append(modelsConfig.Head.Script, models.ScriptTag{
			Src: jsPath,
		})
	}

	cfg.Extra["head"] = modelsConfig.Head
}

func buildHashedURL(assetPath, hash string) string {
	if assetPath == "" {
		return ""
	}
	if hash != "" {
		ext := path.Ext(assetPath)
		base := strings.TrimSuffix(assetPath, ext)
		assetPath = base + "." + hash + ext
	}
	if !strings.HasPrefix(assetPath, "/") {
		assetPath = "/" + assetPath
	}
	return assetPath
}

// parseLinkAvatarsConfig parses the configuration from the manager config.
func parseLinkAvatarsConfig(cfg *lifecycle.Config) LinkAvatarsConfig {
	result := defaultLinkAvatarsConfig()
	if cfg == nil || cfg.Extra == nil {
		return result
	}

	raw := getLinkAvatarsConfigRaw(cfg.Extra)
	if raw == nil {
		return result
	}

	// Allow directly providing a typed config
	if typed, ok := raw.(LinkAvatarsConfig); ok {
		return typed
	}

	m := coerceToMapAny(raw)
	if m == nil {
		return result
	}

	applyLinkAvatarsConfigMap(&result, m)
	return result
}

// getLinkAvatarsConfigRaw retrieves the raw config from Extra.
func getLinkAvatarsConfigRaw(extra map[string]interface{}) any {
	if extra == nil {
		return nil
	}
	if v, ok := extra["link_avatars"]; ok {
		return v
	}
	// Back-compat: allow a nested "markata-go" map
	if markataGo, ok := extra["markata-go"].(map[string]any); ok {
		return markataGo["link_avatars"]
	}
	if markataGo, ok := extra["markata-go"].(map[string]interface{}); ok {
		return markataGo["link_avatars"]
	}
	return nil
}

// applyLinkAvatarsConfigMap applies map values to the config.
func applyLinkAvatarsConfigMap(dst *LinkAvatarsConfig, m map[string]any) {
	if dst == nil || m == nil {
		return
	}

	applyLinkAvatarsBasicFields(dst, m)
	applyLinkAvatarsIgnoreFields(dst, m)
	applyLinkAvatarsSizeAndPosition(dst, m)
}

// applyLinkAvatarsBasicFields applies basic string/bool fields.
func applyLinkAvatarsBasicFields(dst *LinkAvatarsConfig, m map[string]any) {
	if v, ok := m["enabled"].(bool); ok {
		dst.Enabled = v
	}
	if v, ok := m["mode"].(string); ok && strings.TrimSpace(v) != "" {
		dst.Mode = strings.ToLower(strings.TrimSpace(v))
	}
	if v, ok := m["selector"].(string); ok && strings.TrimSpace(v) != "" {
		dst.Selector = v
	}
	if v, ok := m["service"].(string); ok && strings.TrimSpace(v) != "" {
		dst.Service = strings.ToLower(strings.TrimSpace(v))
	}
	if v, ok := m["template"].(string); ok {
		dst.Template = v
	}
	if v, ok := m["hosted_base_url"].(string); ok {
		dst.HostedBaseURL = strings.TrimSpace(v)
	}
}

// applyLinkAvatarsIgnoreFields applies ignore list fields.
func applyLinkAvatarsIgnoreFields(dst *LinkAvatarsConfig, m map[string]any) {
	if v, ok := m["ignore_domains"]; ok {
		dst.IgnoreDomains = coerceStringSlice(v)
	}
	if v, ok := m["ignore_origins"]; ok {
		dst.IgnoreOrigins = coerceStringSlice(v)
	}
	if v, ok := m["ignore_selectors"]; ok {
		dst.IgnoreSelectors = coerceStringSlice(v)
	}
	if v, ok := m["ignore_classes"]; ok {
		dst.IgnoreClasses = coerceStringSlice(v)
	}
	if v, ok := m["ignore_ids"]; ok {
		dst.IgnoreIDs = coerceStringSlice(v)
	}
}

// applyLinkAvatarsSizeAndPosition applies size and position fields.
func applyLinkAvatarsSizeAndPosition(dst *LinkAvatarsConfig, m map[string]any) {
	if v, ok := parseIntFromInterface(m["size"]); ok && v > 0 {
		dst.Size = v
	}

	if v, ok := m["position"].(string); ok {
		pos := strings.ToLower(strings.TrimSpace(v))
		if pos == "before" || pos == "after" {
			dst.Position = pos
		}
	}
}

func (p *LinkAvatarsPlugin) validateConfig() error {
	if !p.config.Enabled {
		return nil
	}
	if p.config.Mode == "" {
		p.config.Mode = linkAvatarModeJS
	}

	switch p.config.Mode {
	case linkAvatarModeJS, linkAvatarModeLocal, linkAvatarModeHosted:
	default:
		return fmt.Errorf("link_avatars mode must be \"js\", \"local\", or \"hosted\"")
	}

	if p.config.Mode == linkAvatarModeHosted && strings.TrimSpace(p.config.HostedBaseURL) == "" {
		return fmt.Errorf("link_avatars hosted_base_url is required when mode = \"hosted\"")
	}

	return nil
}

func resolveOutputDir(cfg *lifecycle.Config) string {
	if cfg == nil || cfg.OutputDir == "" {
		return defaultOutputDir
	}
	return cfg.OutputDir
}

func getSiteOrigin(cfg *lifecycle.Config) string {
	if cfg == nil || cfg.Extra == nil {
		return ""
	}

	if modelsConfig, ok := cfg.Extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		return normalizeOrigin(modelsConfig.URL)
	}

	if urlValue, ok := cfg.Extra["url"].(string); ok {
		return normalizeOrigin(urlValue)
	}

	return ""
}

func normalizeOrigin(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ""
	}
	return parsed.Scheme + "://" + parsed.Host
}

func getSiteURLPath(cfg *lifecycle.Config) string {
	if cfg == nil || cfg.Extra == nil {
		return ""
	}

	if modelsConfig, ok := cfg.Extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		return extractURLPath(modelsConfig.URL)
	}

	if urlValue, ok := cfg.Extra["url"].(string); ok {
		return extractURLPath(urlValue)
	}

	return ""
}

func extractURLPath(raw string) string {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Path == "" {
		return ""
	}
	return strings.TrimSuffix(parsed.Path, "/")
}

func (p *LinkAvatarsPlugin) iconBaseURL() (string, error) {
	switch p.config.Mode {
	case linkAvatarModeLocal:
		return "/assets/markata/link-avatars", nil
	case linkAvatarModeHosted:
		base := strings.TrimRight(p.config.HostedBaseURL, "/")
		if base == "" {
			return "", fmt.Errorf("link_avatars hosted_base_url is required when mode = \"hosted\"")
		}
		return base, nil
	default:
		return "", fmt.Errorf("link_avatars mode must be \"local\" or \"hosted\" for build-time injection")
	}
}

func (p *LinkAvatarsPlugin) processHTMLConcurrent(htmlContent, publicBase, assetsDir string, iconCache *sync.Map) (string, error) {
	wrapped := `<div id="__link-avatars-root">` + htmlContent + `</div>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wrapped))
	if err != nil {
		return htmlContent, err
	}

	root := doc.Find("#__link-avatars-root")
	if root.Length() == 0 {
		return htmlContent, nil
	}

	selectors := parseIgnoreSelectors(p.config.IgnoreSelectors)
	positionClass := "has-avatar-" + p.config.Position

	root.Find(p.config.Selector).Each(func(_ int, link *goquery.Selection) {
		if link.HasClass("has-avatar") {
			return
		}

		href, ok := link.Attr("href")
		if !ok {
			return
		}

		parsed, parseErr := url.Parse(strings.TrimSpace(href))
		if parseErr != nil || !parsed.IsAbs() {
			return
		}
		if parsed.Scheme != schemeHTTP && parsed.Scheme != schemeHTTPS {
			return
		}

		origin := parsed.Scheme + "://" + parsed.Host
		host := parsed.Hostname()
		if host == "" {
			return
		}

		if shouldIgnoreLink(link, host, origin, p.siteOrigin, selectors, p.config) {
			return
		}

		var iconURL string
		if cached, loaded := iconCache.Load(host); loaded {
			if s, ok := cached.(string); ok {
				iconURL = s
			}
		} else {
			fileName, fetchErr := p.ensureIconForHost(assetsDir, host, origin)
			if fetchErr != nil {
				return
			}
			iconURL = strings.TrimRight(publicBase, "/") + "/" + fileName
			iconCache.Store(host, iconURL)
		}

		style, _ := link.Attr("style")
		style = updateStyleAttribute(style, "--favicon-url", fmt.Sprintf("url('%s')", iconURL))
		link.SetAttr("style", style)
		link.SetAttr("data-favicon", iconURL)
		link.AddClass("has-avatar")
		link.AddClass(positionClass)
	})

	updated, err := root.Html()
	if err != nil {
		return htmlContent, err
	}

	return updated, nil
}

func parseIgnoreSelectors(selectors []string) []cascadia.Sel {
	if len(selectors) == 0 {
		return nil
	}

	parsed := make([]cascadia.Sel, 0, len(selectors))
	for _, raw := range selectors {
		value := strings.TrimSpace(raw)
		if value == "" {
			continue
		}
		sel, err := cascadia.Parse(value)
		if err != nil {
			continue
		}
		parsed = append(parsed, sel)
	}

	return parsed
}

func shouldIgnoreLink(link *goquery.Selection, host, origin, siteOrigin string, selectors []cascadia.Sel, cfg LinkAvatarsConfig) bool {
	if siteOrigin != "" && origin == siteOrigin {
		return true
	}

	if link.Find("img, picture").Length() > 0 {
		return true
	}

	lowerHost := strings.ToLower(host)
	for _, domain := range cfg.IgnoreDomains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}
		if lowerHost == domain || strings.HasSuffix(lowerHost, "."+domain) {
			return true
		}
	}

	for _, ignoredOrigin := range cfg.IgnoreOrigins {
		ignoredOrigin = strings.TrimSpace(ignoredOrigin)
		if ignoredOrigin == "" {
			continue
		}
		if origin == ignoredOrigin {
			return true
		}
	}

	for _, className := range cfg.IgnoreClasses {
		className = strings.TrimSpace(className)
		if className == "" {
			continue
		}
		if link.HasClass(className) {
			return true
		}
	}

	if hasIgnoredID(link, cfg.IgnoreIDs) {
		return true
	}

	if matchesIgnoreSelector(link, selectors) {
		return true
	}

	return false
}

func hasIgnoredID(link *goquery.Selection, ids []string) bool {
	if len(ids) == 0 {
		return false
	}
	for _, id := range ids {
		value := strings.TrimSpace(id)
		if value == "" {
			continue
		}
		selector := "#" + value
		if link.Is(selector) {
			return true
		}
		if link.ParentsFiltered(selector).Length() > 0 {
			return true
		}
	}
	return false
}

func matchesIgnoreSelector(link *goquery.Selection, selectors []cascadia.Sel) bool {
	if len(selectors) == 0 {
		return false
	}
	for _, sel := range selectors {
		for _, node := range link.Nodes {
			if sel.Match(node) {
				return true
			}
		}
	}
	return false
}

func (p *LinkAvatarsPlugin) ensureIconForHost(assetsDir, host, origin string) (string, error) {
	safeHost := sanitizeHost(host)
	if safeHost == "" {
		return "", fmt.Errorf("invalid host")
	}

	if cached, ok := findCachedIcon(assetsDir, safeHost); ok {
		return cached, nil
	}
	if hasCachedIconMiss(assetsDir, safeHost) {
		return "", fmt.Errorf("favicon request previously failed")
	}

	faviconURL, err := p.faviconURL(host, origin)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, faviconURL, http.NoBody)
	if err != nil {
		return "", err
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
			cacheIconMiss(assetsDir, safeHost)
		}
		return "", fmt.Errorf("favicon request failed: %s", resp.Status)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return "", err
	}

	ext := iconExtension(resp.Header.Get("Content-Type"), faviconURL)
	fileName := safeHost + ext
	outputPath := filepath.Join(assetsDir, fileName)
	if err := os.WriteFile(outputPath, data, 0o644); err != nil { //nolint:gosec // static assets need world-readable permissions
		return "", err
	}

	return fileName, nil
}

func (p *LinkAvatarsPlugin) faviconURL(host, origin string) (string, error) {
	service := p.config.Service
	if service == "custom" && strings.TrimSpace(p.config.Template) != "" {
		return replaceTemplatePlaceholders(p.config.Template, host, origin), nil
	}

	switch service {
	case "google":
		return fmt.Sprintf("https://www.google.com/s2/favicons?domain=%s&sz=%d", host, p.config.Size), nil
	case "duckduckgo", "":
		return fmt.Sprintf("https://icons.duckduckgo.com/ip3/%s%s", host, linkAvatarIconExtICO), nil
	default:
		return "", fmt.Errorf("unknown link_avatars service %q", service)
	}
}

func replaceTemplatePlaceholders(template, host, origin string) string {
	encodedOrigin := url.QueryEscape(origin)
	return strings.NewReplacer(
		"{host}", host,
		"{origin}", encodedOrigin,
	).Replace(template)
}

func findCachedIcon(assetsDir, safeHost string) (string, bool) {
	for _, ext := range []string{linkAvatarIconExtICO, ".png", ".jpg", ".jpeg", ".svg"} {
		fileName := safeHost + ext
		if _, err := os.Stat(filepath.Join(assetsDir, fileName)); err == nil {
			return fileName, true
		}
	}
	return "", false
}

func hasCachedIconMiss(assetsDir, safeHost string) bool {
	_, err := os.Stat(filepath.Join(assetsDir, safeHost+linkAvatarMissExt))
	return err == nil
}

func cacheIconMiss(assetsDir, safeHost string) {
	//nolint:errcheck,gosec // negative cache is best-effort and static output is intentionally readable
	_ = os.WriteFile(filepath.Join(assetsDir, safeHost+linkAvatarMissExt), []byte("missing\n"), 0o644)
}

func iconExtension(contentType, faviconURL string) string {
	if contentType != "" {
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err == nil {
			switch mediaType {
			case "image/png":
				return ".png"
			case "image/jpeg":
				return ".jpg"
			case "image/svg+xml":
				return ".svg"
			case "image/x-icon", "image/vnd.microsoft.icon":
				return linkAvatarIconExtICO
			}
		}
	}

	parsed, err := url.Parse(faviconURL)
	if err == nil {
		ext := strings.ToLower(path.Ext(parsed.Path))
		switch ext {
		case ".png", ".jpg", ".jpeg", ".svg", linkAvatarIconExtICO:
			return ext
		}
	}

	return linkAvatarIconExtICO
}

func sanitizeHost(host string) string {
	value := strings.ToLower(strings.TrimSpace(host))
	if value == "" {
		return ""
	}
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '.' || r == '-':
			return r
		default:
			return '-'
		}
	}, value)
}

func updateStyleAttribute(style, varName, value string) string {
	style = strings.TrimSpace(style)
	if strings.Contains(style, varName) {
		parts := strings.Split(style, ";")
		filtered := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" {
				continue
			}
			if strings.HasPrefix(trimmed, varName) {
				continue
			}
			filtered = append(filtered, trimmed)
		}
		style = strings.Join(filtered, ";")
	}

	if style != "" && !strings.HasSuffix(style, ";") {
		style += ";"
	}

	return style + fmt.Sprintf("%s: %s;", varName, value)
}

// SetConfig sets the link avatars configuration directly (for testing).
func (p *LinkAvatarsPlugin) SetConfig(config LinkAvatarsConfig) {
	p.config = config
}

// Config returns the current configuration (for testing).
func (p *LinkAvatarsPlugin) Config() LinkAvatarsConfig {
	return p.config
}

// Ensure LinkAvatarsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*LinkAvatarsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*LinkAvatarsPlugin)(nil)
	_ lifecycle.RenderPlugin    = (*LinkAvatarsPlugin)(nil)
	_ lifecycle.WritePlugin     = (*LinkAvatarsPlugin)(nil)
)
