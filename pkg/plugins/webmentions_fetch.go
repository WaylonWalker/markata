// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/buildstats"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
)

var webmentionsFetchLog = logging.Component("webmentions_fetch").Phase("transform")

// Compile-time interface verification.
var (
	_ lifecycle.Plugin          = (*WebmentionsFetchPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*WebmentionsFetchPlugin)(nil)
	_ lifecycle.TransformPlugin = (*WebmentionsFetchPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*WebmentionsFetchPlugin)(nil)
)

// WebmentionIOResponse represents the response from webmention.io API.
type WebmentionIOResponse struct {
	Type     string               `json:"type"`
	Name     string               `json:"name"`
	Children []ReceivedWebMention `json:"children"`
}

// WebmentionsFetchPlugin fetches incoming webmentions from webmention.io.
type WebmentionsFetchPlugin struct {
	config        models.WebMentionsConfig
	siteURL       string
	httpClient    *http.Client
	mentions      []ReceivedWebMention
	mentionsByURL map[string][]ReceivedWebMention
}

// NewWebmentionsFetchPlugin creates a new WebmentionsFetchPlugin.
func NewWebmentionsFetchPlugin() *WebmentionsFetchPlugin {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	buildstats.InstrumentHTTPClient(client)

	return &WebmentionsFetchPlugin{
		config:        models.NewWebMentionsConfig(),
		httpClient:    client,
		mentions:      make([]ReceivedWebMention, 0),
		mentionsByURL: make(map[string][]ReceivedWebMention),
	}
}

// Name returns the unique name of the plugin.
func (p *WebmentionsFetchPlugin) Name() string {
	return "webmentions_fetch"
}

// Configure reads configuration options for the plugin.
func (p *WebmentionsFetchPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Get site URL
	if siteURL, ok := config.Extra["url"].(string); ok {
		p.siteURL = siteURL
	}

	// Get webmentions config from Extra as a map
	if wm, ok := config.Extra["webmentions"]; ok {
		// Handle as map[string]interface{} from TOML parsing
		if wmMap, ok := wm.(map[string]interface{}); ok {
			// Extract fields manually
			if enabled, ok := wmMap["enabled"].(bool); ok {
				p.config.Enabled = enabled
			}
			if outgoing, ok := wmMap["outgoing"].(bool); ok {
				p.config.Outgoing = outgoing
			}
			if cacheDir, ok := wmMap["cache_dir"].(string); ok {
				p.config.CacheDir = cacheDir
			}
			if timeout, ok := wmMap["timeout"].(string); ok {
				p.config.Timeout = timeout
			}
			if token, ok := wmMap["webmention_io_token"].(string); ok {
				p.config.WebmentionIOToken = token
			}
			if userAgent, ok := wmMap["user_agent"].(string); ok {
				p.config.UserAgent = userAgent
			}
			if concurrentRequests, ok := parseIntFromInterface(wmMap["concurrent_requests"]); ok {
				p.config.ConcurrentRequests = concurrentRequests
			}
		}
	}

	// Set timeout if configured
	if p.config.Timeout != "" {
		if timeout, err := time.ParseDuration(p.config.Timeout); err == nil {
			p.httpClient.Timeout = timeout
		}
	}

	return nil
}

// Priority returns the execution priority for this plugin.
// For Transform stage, run before other plugins that might cache post maps
// (like jinja_md which runs at PriorityEarly).
func (p *WebmentionsFetchPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		// Run before jinja_md (PriorityEarly = -100) so webmentions are attached
		// before post maps are cached.
		return lifecycle.PriorityFirst // -200
	}
	return lifecycle.PriorityDefault
}

// Transform loads cached webmentions and attaches them to posts.
func (p *WebmentionsFetchPlugin) Transform(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	// Load cached mentions
	if err := p.loadMentionsCache(); err != nil {
		// Not an error if cache doesn't exist yet
		return nil
	}

	if len(p.mentions) == 0 {
		// No mentions cached
		return nil
	}

	posts := m.Posts()

	// Incremental build optimization: skip attaching when no posts will rebuild
	if cache := GetBuildCache(m); cache != nil {
		if !postsNeedWebmentionAttachments(posts, cache) {
			return nil
		}
	}

	// Build URL mapping for faster lookups
	p.mentionsByURL = p.GroupMentionsByURL()

	// Attach mentions to each post
	attachCount := 0
	for i := range posts {
		post := posts[i]

		// Build full URL for this post
		postURL := p.siteURL + post.Href

		// Normalize the URL (remove double slashes after protocol)
		postURL = strings.Replace(postURL, "://", "PROTOCOL_SEP", 1)
		postURL = strings.ReplaceAll(postURL, "//", "/")
		postURL = strings.Replace(postURL, "PROTOCOL_SEP", "://", 1)

		// Try various URL formats
		urlVariants := []string{
			postURL,
			postURL + "/",
			strings.TrimSuffix(postURL, "/"),
		}

		// Add legacy URL prefixes for migration compatibility
		// Common case: posts used to be at /blog/slug but now at /slug
		legacyPrefixes := []string{"/blog"}
		for _, prefix := range legacyPrefixes {
			legacyURL := p.siteURL + prefix + post.Href
			// Normalize
			legacyURL = strings.Replace(legacyURL, "://", "PROTOCOL_SEP", 1)
			legacyURL = strings.ReplaceAll(legacyURL, "//", "/")
			legacyURL = strings.Replace(legacyURL, "PROTOCOL_SEP", "://", 1)
			urlVariants = append(urlVariants,
				legacyURL,
				legacyURL+"/",
				strings.TrimSuffix(legacyURL, "/"),
			)
		}

		// Also try with double slashes after domain (common in webmention.io data)
		// e.g., https://waylonwalker.com//2025-07-04-notes
		if strings.Contains(postURL, "://") {
			parts := strings.SplitN(postURL, "://", 2)
			if len(parts) == 2 && strings.Contains(parts[1], "/") {
				// Split domain from path
				domainPath := strings.SplitN(parts[1], "/", 2)
				if len(domainPath) == 2 {
					// Create URL with double slash after domain: https://domain//path
					doubleSlashURL := parts[0] + "://" + domainPath[0] + "//" + domainPath[1]
					urlVariants = append(urlVariants,
						doubleSlashURL,
						doubleSlashURL+"/",
						strings.TrimSuffix(doubleSlashURL, "/"),
					)
				}
			}
		}

		// Collect all mentions for this post
		mentionMap := make(map[string]ReceivedWebMention)
		for _, url := range urlVariants {
			for i := range p.mentionsByURL[url] {
				mention := &p.mentionsByURL[url][i]
				// Use source URL as key to avoid duplicates
				mentionMap[mention.Source] = *mention
			}
		}

		if len(mentionMap) > 0 {
			uniqueMentions := make([]ReceivedWebMention, 0, len(mentionMap))
			for key := range mentionMap {
				uniqueMentions = append(uniqueMentions, mentionMap[key])
			}

			// Add to post's Extra data
			if post.Extra == nil {
				post.Extra = make(map[string]interface{})
			}
			post.Extra["webmentions"] = uniqueMentions
			attachCount++
		}
	}

	if attachCount > 0 {
		webmentionsFetchLog.Printf("Attached webmentions to %d posts", attachCount)
	}

	return nil
}

func postsNeedWebmentionAttachments(posts []*models.Post, cache *buildcache.Cache) bool {
	if cache == nil {
		return true
	}
	for _, post := range posts {
		if post.Skip || post.InputHash == "" {
			continue
		}
		if cache.ShouldRebuild(post.Path, post.InputHash, post.Template) {
			return true
		}
	}
	return false
}

// loadMentionsCache loads mentions from the cache file.
func (p *WebmentionsFetchPlugin) loadMentionsCache() error {
	if p.config.CacheDir == "" {
		return nil
	}

	cacheFile := filepath.Join(p.config.CacheDir, "received_mentions.json")

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Cache doesn't exist yet
		}
		return err
	}

	if err := json.Unmarshal(data, &p.mentions); err != nil {
		return err
	}
	for i := range p.mentions {
		normalizeReceivedWebMention(&p.mentions[i])
	}
	return nil
}

// FetchMentions fetches all webmentions from webmention.io for the configured domain.
func (p *WebmentionsFetchPlugin) FetchMentions() error {
	// Get token from config or environment variable
	token := p.config.WebmentionIOToken
	if token == "" || strings.HasPrefix(token, "${") {
		// Try environment variable
		token = os.Getenv("WEBMENTION_IO_TOKEN")
	}

	if token == "" {
		return fmt.Errorf("webmention_io_token not configured - set it in config or WEBMENTION_IO_TOKEN environment variable")
	}

	if p.siteURL == "" {
		return fmt.Errorf("site URL not configured")
	}

	// Extract domain from site URL
	domain := extractDomain(p.siteURL)

	// Build API URL
	apiURL := fmt.Sprintf("https://webmention.io/api/mentions.jf2?token=%s&domain=%s&per-page=1000",
		token, domain)

	// Fetch mentions
	ctx, cancel := context.WithTimeout(context.Background(), p.httpClient.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("User-Agent", p.config.UserAgent)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fetch mentions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) //nolint:errcheck // Best effort read for error message
		return fmt.Errorf("webmention.io API error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse response
	var wmResp WebmentionIOResponse
	if err := json.NewDecoder(resp.Body).Decode(&wmResp); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	p.mentions = wmResp.Children
	for i := range p.mentions {
		normalizeReceivedWebMention(&p.mentions[i])
	}

	// Save to cache
	if err := p.saveMentionsCache(); err != nil {
		return fmt.Errorf("save cache: %w", err)
	}

	return nil
}

// saveMentionsCache saves the fetched mentions to a JSON file.
func (p *WebmentionsFetchPlugin) saveMentionsCache() error {
	if p.config.CacheDir == "" {
		return nil
	}

	// Create cache directory
	if err := os.MkdirAll(p.config.CacheDir, 0o755); err != nil {
		return err
	}

	cacheFile := filepath.Join(p.config.CacheDir, "received_mentions.json")
	for i := range p.mentions {
		normalizeReceivedWebMention(&p.mentions[i])
	}

	data, err := json.MarshalIndent(p.mentions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0o600)
}

// GetMentions returns all fetched mentions.
func (p *WebmentionsFetchPlugin) GetMentions() []ReceivedWebMention {
	return p.mentions
}

// GetMentionsForURL returns mentions for a specific URL.
func (p *WebmentionsFetchPlugin) GetMentionsForURL(targetURL string) []ReceivedWebMention {
	var results []ReceivedWebMention
	for i := range p.mentions {
		if p.mentions[i].Target == targetURL {
			results = append(results, p.mentions[i])
		}
	}
	return results
}

// GroupMentionsByURL groups mentions by their target URL.
func (p *WebmentionsFetchPlugin) GroupMentionsByURL() map[string][]ReceivedWebMention {
	groups := make(map[string][]ReceivedWebMention)
	for i := range p.mentions {
		groups[p.mentions[i].Target] = append(groups[p.mentions[i].Target], p.mentions[i])
	}
	return groups
}

// extractDomain extracts the domain from a URL.
func extractDomain(siteURL string) string {
	// Remove protocol
	domain := siteURL
	if idx := strings.Index(domain, "://"); idx != -1 {
		domain = domain[idx+3:]
	}
	// Remove path
	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}
	// Remove port
	if idx := strings.Index(domain, ":"); idx != -1 {
		domain = domain[:idx]
	}
	return domain
}
