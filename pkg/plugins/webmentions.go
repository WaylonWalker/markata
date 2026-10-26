// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// WebMention represents a webmention to be sent.
type WebMention struct {
	// Source is the URL of the page containing the link (your page)
	Source string `json:"source"`

	// Target is the URL being linked to (external page)
	Target string `json:"target"`

	// Endpoint is the discovered webmention endpoint
	Endpoint string `json:"endpoint,omitempty"`

	// Sent indicates whether the webmention was successfully sent
	Sent bool `json:"sent"`

	// SentAt is when the webmention was sent
	SentAt *time.Time `json:"sent_at,omitempty"`

	// Error contains any error message from sending
	Error string `json:"error,omitempty"`

	// StatusCode is the HTTP status code from the endpoint
	StatusCode int `json:"status_code,omitempty"`
}

// WebMentionsPlugin sends outgoing webmentions for external links in posts.
// It runs in the Collect stage after posts have been rendered and links collected.
type WebMentionsPlugin struct {
	// config holds the plugin configuration
	config models.WebMentionsConfig

	// siteURL is the base URL of the site
	siteURL string

	// httpClient is the HTTP client for sending webmentions
	httpClient *http.Client

	// mentions is the list of webmentions discovered
	mentions []*WebMention

	// mu protects concurrent access to mentions
	mu sync.Mutex

	// sentCache tracks already-sent webmentions to avoid duplicates
	sentCache map[string]bool
}

// NewWebMentionsPlugin creates a new WebMentionsPlugin.
func NewWebMentionsPlugin() *WebMentionsPlugin {
	return &WebMentionsPlugin{
		config:    models.NewWebMentionsConfig(),
		mentions:  make([]*WebMention, 0),
		sentCache: make(map[string]bool),
	}
}

// Name returns the unique name of the plugin.
func (p *WebMentionsPlugin) Name() string {
	return "webmentions"
}

// Configure reads configuration options for the plugin.
func (p *WebMentionsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Get webmentions config
	if wm, ok := config.Extra["webmentions"]; ok {
		if wmConfig, ok := wm.(models.WebMentionsConfig); ok {
			p.config = wmConfig
		}
	}

	// Get site URL
	if siteURL, ok := config.Extra["url"].(string); ok {
		p.siteURL = siteURL
	}

	// Set up HTTP client with configured timeout
	timeout, err := time.ParseDuration(p.config.Timeout)
	if err != nil {
		timeout = 30 * time.Second
	}

	p.httpClient = &http.Client{
		Timeout: timeout,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}
	buildstats.InstrumentHTTPClient(p.httpClient)

	// Load sent cache from disk
	if p.config.CacheDir != "" {
		p.loadSentCache()
	}

	return nil
}

// Priority returns the plugin's priority for a given stage.
func (p *WebMentionsPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageCollect:
		// Run after link_collector to have access to outlinks
		return lifecycle.PriorityLate + 50
	default:
		return lifecycle.PriorityDefault
	}
}

// Collect discovers and sends webmentions for external links.
func (p *WebMentionsPlugin) Collect(m *lifecycle.Manager) error {
	if !p.config.Enabled || !p.config.Outgoing {
		return nil
	}

	if p.siteURL == "" {
		// Skip if no site URL configured - can't determine source URLs
		return nil
	}

	posts := m.Posts()
	if len(posts) == 0 {
		return nil
	}

	// Collect all external links from posts
	var allMentions []*WebMention
	for _, post := range posts {
		if post.Skip || len(post.Outlinks) == 0 {
			continue
		}

		// Build source URL for this post
		sourceURL := strings.TrimSuffix(p.siteURL, "/") + post.Href

		// Process external links
		for _, link := range post.Outlinks {
			if link.IsInternal {
				continue
			}

			// Create webmention
			mention := &WebMention{
				Source: sourceURL,
				Target: link.TargetURL,
			}

			allMentions = append(allMentions, mention)
		}
	}

	if len(allMentions) == 0 {
		return nil
	}

	// Process webmentions concurrently
	concurrency := p.config.ConcurrentRequests
	if concurrency <= 0 {
		concurrency = 5
	}

	semaphore := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, mention := range allMentions {
		// Skip if already sent (cached)
		cacheKey := p.cacheKey(mention.Source, mention.Target)
		if p.sentCache[cacheKey] {
			mention.Sent = true
			continue
		}

		wg.Add(1)
		go func(m *WebMention) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			p.processWebMention(m)
		}(mention)
	}

	wg.Wait()

	// Store mentions
	p.mu.Lock()
	p.mentions = allMentions
	p.mu.Unlock()

	// Save sent cache
	if p.config.CacheDir != "" {
		p.saveSentCache()
	}

	// Store in manager cache for potential reporting
	m.Cache().Set("webmentions", allMentions)

	return nil
}

// processWebMention discovers the endpoint and sends the webmention.
func (p *WebMentionsPlugin) processWebMention(mention *WebMention) {
	// Discover webmention endpoint
	endpoint, err := p.discoverEndpoint(mention.Target)
	if err != nil {
		mention.Error = fmt.Sprintf("endpoint discovery failed: %v", err)
		return
	}

	if endpoint == "" {
		// No webmention endpoint found - not an error, just not supported
		return
	}

	mention.Endpoint = endpoint

	// Send the webmention
	if err := p.sendWebMention(mention); err != nil {
		mention.Error = err.Error()
		return
	}

	// Mark as sent and cache
	mention.Sent = true
	now := time.Now()
	mention.SentAt = &now

	cacheKey := p.cacheKey(mention.Source, mention.Target)
	p.mu.Lock()
	p.sentCache[cacheKey] = true
	p.mu.Unlock()
}

// discoverEndpoint discovers the webmention endpoint for a target URL.
// It checks both HTTP Link headers and HTML <link> elements.
func (p *WebMentionsPlugin) discoverEndpoint(targetURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.httpClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", p.config.UserAgent)
	req.Header.Set("Accept", "text/html, application/xhtml+xml, */*")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch target: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	// Check HTTP Link header first
	if endpoint := p.extractEndpointFromHeader(resp.Header, targetURL); endpoint != "" {
		return endpoint, nil
	}

	// Read body for HTML parsing (limit to 512KB to avoid huge pages)
	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	// Check HTML for <link rel="webmention">
	if endpoint := p.extractEndpointFromHTML(string(body), targetURL); endpoint != "" {
		return endpoint, nil
	}

	return "", nil
}

// linkHeaderRegex matches Link headers with rel="webmention"
var linkHeaderRegex = regexp.MustCompile(`<([^>]+)>;\s*rel=["']?webmention["']?`)

// extractEndpointFromHeader extracts webmention endpoint from HTTP Link header.
func (p *WebMentionsPlugin) extractEndpointFromHeader(headers http.Header, baseURL string) string {
	linkHeaders := headers.Values("Link")
	for _, header := range linkHeaders {
		// Parse each Link header value
		for _, part := range strings.Split(header, ",") {
			part = strings.TrimSpace(part)

			// Check if this Link has rel="webmention"
			if strings.Contains(strings.ToLower(part), "rel=\"webmention\"") ||
				strings.Contains(strings.ToLower(part), "rel='webmention'") ||
				strings.Contains(strings.ToLower(part), "rel=webmention") {
				// Extract the URL from <...>
				match := linkHeaderRegex.FindStringSubmatch(part)
				if len(match) >= 2 {
					return p.resolveURL(baseURL, match[1])
				}
			}
		}
	}
	return ""
}

// webmentionLinkRegex matches <link rel="webmention" href="..."> in HTML
var webmentionLinkRegex = regexp.MustCompile(`(?i)<link[^>]+rel=["']?webmention["']?[^>]+href=["']([^"']+)["']`)
var webmentionLinkRegex2 = regexp.MustCompile(`(?i)<link[^>]+href=["']([^"']+)["'][^>]+rel=["']?webmention["']?`)

// extractEndpointFromHTML extracts webmention endpoint from HTML content.
func (p *WebMentionsPlugin) extractEndpointFromHTML(html, baseURL string) string {
	// Try both attribute orderings
	if match := webmentionLinkRegex.FindStringSubmatch(html); len(match) >= 2 {
		return p.resolveURL(baseURL, match[1])
	}

	if match := webmentionLinkRegex2.FindStringSubmatch(html); len(match) >= 2 {
		return p.resolveURL(baseURL, match[1])
	}

	// Also check for <a rel="webmention"> which some sites use
	aTagRegex := regexp.MustCompile(`(?i)<a[^>]+rel=["']?webmention["']?[^>]+href=["']([^"']+)["']`)
	if match := aTagRegex.FindStringSubmatch(html); len(match) >= 2 {
		return p.resolveURL(baseURL, match[1])
	}

	aTagRegex2 := regexp.MustCompile(`(?i)<a[^>]+href=["']([^"']+)["'][^>]+rel=["']?webmention["']?`)
	if match := aTagRegex2.FindStringSubmatch(html); len(match) >= 2 {
		return p.resolveURL(baseURL, match[1])
	}

	return ""
}

// resolveURL resolves a potentially relative URL against a base URL.
func (p *WebMentionsPlugin) resolveURL(baseURL, href string) string {
	if href == "" {
		return ""
	}

	// Already absolute
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return href
	}

	ref, err := url.Parse(href)
	if err != nil {
		return href
	}

	return base.ResolveReference(ref).String()
}

// sendWebMention sends a webmention to the discovered endpoint.
func (p *WebMentionsPlugin) sendWebMention(mention *WebMention) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.httpClient.Timeout)
	defer cancel()

	// Prepare form data
	data := url.Values{}
	data.Set("source", mention.Source)
	data.Set("target", mention.Target)

	req, err := http.NewRequestWithContext(ctx, "POST", mention.Endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", p.config.UserAgent)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	defer resp.Body.Close()

	mention.StatusCode = resp.StatusCode

	// Accept any 2xx status as success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	// Read error body for debugging
	//nolint:errcheck // Error body read is best-effort for debugging
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
}

// cacheKey generates a unique key for a source/target pair.
func (p *WebMentionsPlugin) cacheKey(source, target string) string {
	hash := sha256.Sum256([]byte(source + "|" + target))
	return hex.EncodeToString(hash[:16])
}

// loadSentCache loads the sent cache from disk.
func (p *WebMentionsPlugin) loadSentCache() {
	cacheFile := filepath.Join(p.config.CacheDir, "webmentions_sent.json")

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return // Cache doesn't exist yet
	}

	var cache map[string]bool
	if err := json.Unmarshal(data, &cache); err != nil {
		return
	}

	p.mu.Lock()
	p.sentCache = cache
	p.mu.Unlock()
}

// saveSentCache saves the sent cache to disk.
func (p *WebMentionsPlugin) saveSentCache() {
	if err := os.MkdirAll(p.config.CacheDir, 0o755); err != nil {
		return
	}

	cacheFile := filepath.Join(p.config.CacheDir, "webmentions_sent.json")

	p.mu.Lock()
	data, err := json.MarshalIndent(p.sentCache, "", "  ")
	p.mu.Unlock()

	if err != nil {
		return
	}

	//nolint:errcheck // Cache writes are best-effort
	os.WriteFile(cacheFile, data, 0o600)
}

// Mentions returns the discovered webmentions.
func (p *WebMentionsPlugin) Mentions() []*WebMention {
	p.mu.Lock()
	defer p.mu.Unlock()
	result := make([]*WebMention, len(p.mentions))
	copy(result, p.mentions)
	return result
}

// Ensure WebMentionsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*WebMentionsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*WebMentionsPlugin)(nil)
	_ lifecycle.CollectPlugin   = (*WebMentionsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*WebMentionsPlugin)(nil)
)
