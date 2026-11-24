package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	stdhtml "html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
)

// OEmbedResponse represents the JSON oEmbed response fields we care about.
// See: https://oembed.com
type OEmbedResponse struct {
	Type            string `json:"type"`
	Version         string `json:"version"`
	Title           string `json:"title"`
	URL             string `json:"url"`
	AuthorName      string `json:"author_name"`
	AuthorURL       string `json:"author_url"`
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	ThumbnailURL    string `json:"thumbnail_url"`
	ThumbnailWidth  int    `json:"thumbnail_width"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	HTML            string `json:"html"`
	Width           string `json:"width"`
	Height          string `json:"height"`
	CacheAge        int    `json:"cache_age"`

	// Extra holds provider-specific metadata for rendering.
	// This field is not part of the oEmbed spec and is not serialized.
	Extra map[string]string `json:"-"`
}

// oembedProviderJSON represents a provider from providers.json
type oembedProviderJSON struct {
	ProviderName string `json:"provider_name"`
	ProviderURL  string `json:"provider_url"`
	Endpoints    []struct {
		Schemes   []string `json:"schemes"`
		URL       string   `json:"url"`
		Discovery bool     `json:"discovery"`
		Formats   []string `json:"formats"`
	} `json:"endpoints"`
}

// oembedProvider describes a single oEmbed provider.
type oembedProvider struct {
	Name             string
	Endpoint         string
	URLPrefixes      []string
	RequiresAuth     bool
	SupportsFormat   bool
	SupportsDiscover bool
	CustomFetch      func(resolver *oembedResolver, rawURL string) (*OEmbedResponse, error)
	// IsCustom indicates this is a hardcoded provider that should be tried before providers.json
	IsCustom bool
}

// oembedResolver resolves oEmbed data for URLs.
type oembedResolver struct {
	client          *http.Client
	config          models.EmbedsConfig
	providers       []oembedProvider
	jsonProviders   []oembedProvider
	jsonProvidersMu sync.RWMutex
	jsonFetchedAt   time.Time
	markdown        goldmark.Markdown
	htmlFetcher     func(rawURL string) (string, error)
}

const jsonProvidersCacheDuration = 24 * time.Hour

func newOEmbedResolver(config models.EmbedsConfig, client *http.Client) *oembedResolver {
	return &oembedResolver{
		client:    client,
		config:    config,
		providers: defaultOEmbedProviders(),
	}
}

func newOEmbedResolverWithProviders(config models.EmbedsConfig, client *http.Client, providers []oembedProvider) *oembedResolver {
	return &oembedResolver{
		client:    client,
		config:    config,
		providers: providers,
	}
}

func (r *oembedResolver) updateConfig(config models.EmbedsConfig) {
	r.config = config
	r.markdown = nil
}

func (r *oembedResolver) setMarkdownRenderer(md goldmark.Markdown) {
	r.markdown = md
}

func (r *oembedResolver) setHTMLFetcher(fetcher func(rawURL string) (string, error)) {
	r.htmlFetcher = fetcher
}

func newEmbedMarkdownRenderer(extra map[string]interface{}) goldmark.Markdown {
	chromaTheme, lineNumbers, enabled := resolveEmbedHighlightConfig(extra)
	if !enabled {
		chromaTheme = disabledHighlightTheme
		lineNumbers = false
	}
	if chromaTheme == "" {
		chromaTheme = palettes.DefaultChromaThemeDark
	}

	formatOptions := []chromahtml.Option{
		chromahtml.WithClasses(true),
		chromahtml.WithAllClasses(true),
	}
	if lineNumbers {
		formatOptions = append(formatOptions, chromahtml.WithLineNumbers(true))
	}
	options := []highlighting.Option{
		highlighting.WithStyle(chromaTheme),
		highlighting.WithFormatOptions(formatOptions...),
	}

	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.Linkify,
			extension.TaskList,
			highlighting.NewHighlighting(options...),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
}

func resolveEmbedHighlightConfig(extra map[string]interface{}) (chromaTheme string, lineNumbers, enabled bool) {
	enabled = true

	if highlight, ok := highlightConfigFromExtra(extra); ok {
		enabled = highlight.IsEnabled()
		chromaTheme = highlight.Theme
		lineNumbers = highlight.LineNumbers
	}

	if chromaTheme == "" {
		if themeExtra, ok := extra["theme"].(map[string]interface{}); ok {
			if palette, ok := themeExtra["palette"].(string); ok && palette != "" {
				chromaTheme = palettes.ChromaTheme(palette)
			}
		}
	}

	return chromaTheme, lineNumbers, enabled
}

func renderGistCodeMarkdown(resolver *oembedResolver, language, content string) (string, error) {
	if language == "" {
		language = formatText
	}
	md := resolver.markdown
	if md == nil {
		md = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Table,
				extension.Strikethrough,
				extension.Linkify,
				extension.TaskList,
				highlighting.NewHighlighting(),
			),
		)
	}

	markdown := fmt.Sprintf("```%s\n%s\n```\n", language, content)
	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderGistMarkdownEmbedHTML(gistURL, filename, language, codeHTML string) string {
	var sb strings.Builder
	sb.WriteString(`<div class="embed-gist">`)
	sb.WriteString("\n")
	sb.WriteString(`  <div class="embed-gist-header">`)
	sb.WriteString("\n")
	sb.WriteString(`    <a href="`)
	sb.WriteString(stdhtml.EscapeString(gistURL))
	sb.WriteString(`" target="_blank" rel="noopener noreferrer">`)
	sb.WriteString(stdhtml.EscapeString(filename))
	sb.WriteString(`</a>`)
	sb.WriteString("\n")
	if language != "" {
		sb.WriteString(`    <span class="embed-gist-language">`)
		sb.WriteString(stdhtml.EscapeString(strings.ToLower(language)))
		sb.WriteString(`</span>`)
		sb.WriteString("\n")
	}
	sb.WriteString(`  </div>`)
	sb.WriteString("\n")
	sb.WriteString(codeHTML)
	sb.WriteString("\n")
	sb.WriteString(`</div>`)
	sb.WriteString("\n")

	return sb.String()
}

// fetchProvidersJSON fetches and parses the providers.json file.
// This is called lazily on first resolve attempt.
func (r *oembedResolver) fetchProvidersJSON() error {
	providersURL := r.config.OEmbedProvidersURL
	if providersURL == "" {
		return nil // Disabled by user
	}

	// Check if we already have cached providers
	r.jsonProvidersMu.RLock()
	if len(r.jsonProviders) > 0 && time.Since(r.jsonFetchedAt) < jsonProvidersCacheDuration {
		r.jsonProvidersMu.RUnlock()
		return nil
	}
	r.jsonProvidersMu.RUnlock()

	// Fetch providers.json
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, providersURL, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create request for providers.json: %w", err)
	}
	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch providers.json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("providers.json returned status %s", resp.Status)
	}

	var jsonProviders []oembedProviderJSON
	if err := json.NewDecoder(resp.Body).Decode(&jsonProviders); err != nil {
		return fmt.Errorf("failed to parse providers.json: %w", err)
	}

	// Convert to internal format
	providers := make([]oembedProvider, 0, len(jsonProviders))
	for _, jp := range jsonProviders {
		for _, ep := range jp.Endpoints {
			if ep.URL == "" {
				continue
			}
			provider := oembedProvider{
				Name:             jp.ProviderName,
				URLPrefixes:      ep.Schemes,
				Endpoint:         ep.URL,
				SupportsDiscover: ep.Discovery,
			}
			// Check if format=json is required
			for _, f := range ep.Formats {
				if f == "json" {
					provider.SupportsFormat = true
					break
				}
			}
			providers = append(providers, provider)
		}
	}

	r.jsonProvidersMu.Lock()
	r.jsonProviders = providers
	r.jsonFetchedAt = time.Now()
	r.jsonProvidersMu.Unlock()

	return nil
}

func (r *oembedResolver) Resolve(rawURL string) (*OEmbedResponse, bool, error) {
	// First, try custom (hardcoded) providers - these have priority
	provider := r.matchCustomProvider(rawURL)
	if provider != nil {
		return r.resolveWithProvider(provider, rawURL)
	}

	// Try providers.json if available
	if r.config.OEmbedProvidersURL != "" {
		if err := r.fetchProvidersJSON(); err == nil {
			provider = r.matchJSONProvider(rawURL)
			if provider != nil {
				return r.resolveWithProvider(provider, rawURL)
			}
		}
	}

	// Fall back to trying non-custom hardcoded providers (for backwards compatibility)
	provider = r.matchProvider(rawURL)
	if provider != nil {
		return r.resolveWithProvider(provider, rawURL)
	}

	// Fall back to auto-discovery if enabled
	if !r.config.OEmbedAutoDiscover {
		return nil, false, nil
	}

	endpoint, err := r.discoverEndpoint(rawURL)
	if err != nil {
		if errors.Is(err, errOEmbedDiscoveryDisabled) {
			return nil, false, nil
		}
		return nil, false, err
	}

	payload, err := r.fetchOEmbedResponse(endpoint)
	if err != nil {
		return nil, false, err
	}

	return payload, true, nil
}

// resolveWithProvider handles the actual resolution given a matched provider
func (r *oembedResolver) resolveWithProvider(provider *oembedProvider, rawURL string) (*OEmbedResponse, bool, error) {
	if !r.isProviderEnabled(provider.Name) {
		return nil, true, nil
	}

	if provider.SupportsDiscover {
		endpoint, err := r.discoverEndpoint(rawURL)
		if err != nil {
			if errors.Is(err, errOEmbedDiscoveryDisabled) {
				return nil, false, nil
			}
			return nil, false, err
		}

		payload, err := r.fetchOEmbedResponse(endpoint)
		if err != nil {
			return nil, false, err
		}

		return payload, true, nil
	}

	if provider.CustomFetch != nil {
		payload, err := provider.CustomFetch(r, rawURL)
		if err != nil {
			return nil, true, err
		}
		return payload, true, nil
	}

	endpoint, err := r.buildEndpoint(provider, rawURL)
	if err != nil {
		return nil, true, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, true, err
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "application/json")

	payload, err := r.fetchOEmbedResponse(endpoint)
	if err != nil {
		return nil, true, err
	}

	return payload, true, nil
}

// matchCustomProvider only matches hardcoded providers (IsCustom=true)
func (r *oembedResolver) matchCustomProvider(rawURL string) *oembedProvider {
	for _, provider := range r.providers {
		if !provider.IsCustom {
			continue
		}
		for _, prefix := range provider.URLPrefixes {
			if strings.HasPrefix(rawURL, prefix) {
				return &provider
			}
		}
	}
	return nil
}

// matchJSONProvider matches providers loaded from providers.json
func (r *oembedResolver) matchJSONProvider(rawURL string) *oembedProvider {
	r.jsonProvidersMu.RLock()
	defer r.jsonProvidersMu.RUnlock()

	for _, provider := range r.jsonProviders {
		for _, prefix := range provider.URLPrefixes {
			// Support wildcards in schemes
			if strings.Contains(prefix, "*") {
				pattern := strings.ReplaceAll(prefix, "*", ".*")
				matched, err := regexp.MatchString(pattern, rawURL)
				if err == nil && matched {
					return &provider
				}
			}
			if strings.HasPrefix(rawURL, prefix) {
				return &provider
			}
		}
	}
	return nil
}

func (r *oembedResolver) matchProvider(rawURL string) *oembedProvider {
	for _, provider := range r.providers {
		for _, prefix := range provider.URLPrefixes {
			if strings.HasPrefix(rawURL, prefix) {
				return &provider
			}
		}
	}

	if r.config.OEmbedAutoDiscover {
		for _, provider := range r.providers {
			if provider.SupportsDiscover {
				return &provider
			}
		}
	}

	return nil
}

func (r *oembedResolver) isProviderEnabled(name string) bool {
	if r.config.OEmbedProviders == nil {
		return true
	}

	providerConfig, ok := r.config.OEmbedProviders[strings.ToLower(name)]
	if !ok {
		return true
	}

	return providerConfig.Enabled
}

func (r *oembedResolver) buildEndpoint(provider *oembedProvider, rawURL string) (string, error) {
	encodedURL := url.QueryEscape(rawURL)
	if provider.SupportsFormat {
		return fmt.Sprintf("%s?url=%s&format=json", provider.Endpoint, encodedURL), nil
	}

	parsed, err := url.Parse(provider.Endpoint)
	if err != nil {
		return "", err
	}

	query := parsed.Query()
	query.Set("url", rawURL)
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}

func (r *oembedResolver) fetchOEmbedResponse(endpoint string) (*OEmbedResponse, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oembed request failed: %s", resp.Status)
	}

	var payload OEmbedResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func (r *oembedResolver) discoverEndpoint(rawURL string) (string, error) {
	if !r.config.OEmbedAutoDiscover {
		return "", errOEmbedDiscoveryDisabled
	}

	htmlContent, err := r.fetchPageHTML(rawURL)
	if err != nil {
		return "", err
	}

	endpoint := findOEmbedLink(htmlContent)
	if endpoint == "" {
		return "", fmt.Errorf("oembed discovery: endpoint not found")
	}

	if strings.HasPrefix(endpoint, "//") {
		endpoint = "https:" + endpoint
	}

	if strings.HasPrefix(endpoint, "/") {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return "", err
		}
		endpoint = parsed.ResolveReference(&url.URL{Path: endpoint}).String()
	}

	return endpoint, nil
}

func (r *oembedResolver) fetchPageHTML(rawURL string) (string, error) {
	if r.htmlFetcher != nil {
		return r.htmlFetcher(rawURL)
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, http.NoBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("oembed discovery failed: %s", resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func findOEmbedLink(html string) string {
	patterns := []*regexp.Regexp{
		oembedLinkRelFirstRe,
		oembedLinkTypeFirstRe,
		oembedLinkHrefMidRe,
	}

	for _, pattern := range patterns {
		match := pattern.FindStringSubmatch(html)
		if len(match) > 1 {
			return match[1]
		}
	}

	return ""
}

// Pre-compiled regexes for oEmbed link discovery.
var (
	oembedLinkRelFirstRe  = regexp.MustCompile(`<link[^>]*rel=["']alternate["'][^>]*type=["']application/json\+oembed["'][^>]*href=["']([^"']+)["']`)
	oembedLinkTypeFirstRe = regexp.MustCompile(`<link[^>]*type=["']application/json\+oembed["'][^>]*rel=["']alternate["'][^>]*href=["']([^"']+)["']`)
	oembedLinkHrefMidRe   = regexp.MustCompile(`<link[^>]*rel=["']alternate["'][^>]*href=["']([^"']+)["'][^>]*type=["']application/json\+oembed["']`)
)

var errOEmbedDiscoveryDisabled = errors.New("oembed discovery disabled")

func defaultOEmbedProviders() []oembedProvider {
	return []oembedProvider{
		{
			Name:           "youtube",
			Endpoint:       "https://www.youtube.com/oembed",
			URLPrefixes:    []string{"https://www.youtube.com/", "https://youtu.be/"},
			RequiresAuth:   false,
			SupportsFormat: true,
			IsCustom:       true,
		},
		{
			Name:           "vimeo",
			Endpoint:       "https://vimeo.com/api/oembed.json",
			URLPrefixes:    []string{"https://vimeo.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "tiktok",
			Endpoint:       "https://www.tiktok.com/oembed",
			URLPrefixes:    []string{"https://www.tiktok.com/", "https://tiktok.com/", "https://vm.tiktok.com/", "https://vt.tiktok.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			CustomFetch:    fetchTikTokEmbed,
			IsCustom:       true,
		},
		{
			Name:           "flickr",
			Endpoint:       "https://www.flickr.com/services/oembed/",
			URLPrefixes:    []string{"https://www.flickr.com/", "https://flickr.com/"},
			RequiresAuth:   false,
			SupportsFormat: true,
			IsCustom:       true,
		},
		{
			Name:           "spotify",
			Endpoint:       "https://open.spotify.com/oembed",
			URLPrefixes:    []string{"https://open.spotify.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "soundcloud",
			Endpoint:       "https://soundcloud.com/oembed",
			URLPrefixes:    []string{"https://soundcloud.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "codepen",
			Endpoint:       "https://codepen.io/api/oembed",
			URLPrefixes:    []string{"https://codepen.io/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "codesandbox",
			Endpoint:       "https://codesandbox.io/oembed",
			URLPrefixes:    []string{"https://codesandbox.io/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "jsfiddle",
			Endpoint:       "https://jsfiddle.net/services/oembed/",
			URLPrefixes:    []string{"https://jsfiddle.net/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "observable",
			Endpoint:       "https://api.observablehq.com/oembed",
			URLPrefixes:    []string{"https://observablehq.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "github",
			Endpoint:       "https://github.com/services/oembed",
			URLPrefixes:    []string{"https://gist.github.com/", "https://github.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			CustomFetch:    fetchGitHubGistEmbed,
			IsCustom:       true,
		},
		{
			Name:           "slideshare",
			Endpoint:       "https://www.slideshare.net/api/oembed/2",
			URLPrefixes:    []string{"https://www.slideshare.net/", "https://slideshare.net/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "prezi",
			Endpoint:       "https://prezi.com/services/oembed/",
			URLPrefixes:    []string{"https://prezi.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "speakerdeck",
			Endpoint:       "https://speakerdeck.com/oembed.json",
			URLPrefixes:    []string{"https://speakerdeck.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "issuu",
			Endpoint:       "https://issuu.com/oembed",
			URLPrefixes:    []string{"https://issuu.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "datawrapper",
			Endpoint:       "https://api.datawrapper.de/v3/oembed/",
			URLPrefixes:    []string{"https://datawrapper.de/", "https://www.datawrapper.de/", "https://app.datawrapper.de/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "flourish",
			Endpoint:       "https://app.flourish.studio/api/v1/oembed",
			URLPrefixes:    []string{"https://public.flourish.studio/", "https://app.flourish.studio/", "https://flourish.studio/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "infogram",
			Endpoint:       "https://infogram.com/oembed",
			URLPrefixes:    []string{"https://infogram.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "reddit",
			Endpoint:       "https://www.reddit.com/oembed",
			URLPrefixes:    []string{"https://www.reddit.com/", "https://reddit.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			CustomFetch:    fetchRedditImage,
			IsCustom:       true,
		},
		{
			Name:           "dailymotion",
			Endpoint:       "https://www.dailymotion.com/services/oembed",
			URLPrefixes:    []string{"https://www.dailymotion.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "wistia",
			Endpoint:       "https://fast.wistia.com/oembed.json",
			URLPrefixes:    []string{"https://wistia.com/", "https://fast.wistia.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			IsCustom:       true,
		},
		{
			Name:           "giphy",
			Endpoint:       "https://giphy.com/services/oembed",
			URLPrefixes:    []string{"https://giphy.com/"},
			RequiresAuth:   false,
			SupportsFormat: false,
			CustomFetch:    fetchGiphyEmbed,
			IsCustom:       true,
		},
		{
			Name:             "oembed",
			Endpoint:         "",
			URLPrefixes:      []string{},
			RequiresAuth:     false,
			SupportsFormat:   false,
			SupportsDiscover: true,
		},
	}
}

type redditPostData struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	PostHint string `json:"post_hint"`
	Preview  struct {
		Images []struct {
			Source struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"source"`
		} `json:"images"`
	} `json:"preview"`
	Thumbnail json.RawMessage `json:"thumbnail"`
}

type redditAPIResponse struct {
	Data struct {
		Children []struct {
			Data redditPostData `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func fetchRedditImage(resolver *oembedResolver, rawURL string) (*OEmbedResponse, error) {
	jsonURL := rawURL + "/.json"

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, jsonURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("reddit fetch: %w", err)
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "application/json")

	resp, err := resolver.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("reddit fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reddit fetch failed: %s", resp.Status)
	}

	var apiResp []redditAPIResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("reddit parse: %w", err)
	}

	if len(apiResp) == 0 || len(apiResp[0].Data.Children) == 0 {
		return nil, fmt.Errorf("reddit: no post data found")
	}

	post := apiResp[0].Data.Children[0].Data

	var thumbnailURL string
	var thumbnailWidth, thumbnailHeight int

	switch {
	case len(post.Preview.Images) > 0:
		thumbnailURL = post.Preview.Images[0].Source.URL
		thumbnailWidth = post.Preview.Images[0].Source.Width
		thumbnailHeight = post.Preview.Images[0].Source.Height
	case post.PostHint == "image":
		thumbnailURL = post.URL
	case len(post.Thumbnail) > 0 && !strings.HasPrefix(string(post.Thumbnail), `"`):
		var thumb struct {
			Source struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"source"`
		}
		if err := json.Unmarshal(post.Thumbnail, &thumb); err == nil {
			thumbnailURL = thumb.Source.URL
			thumbnailWidth = thumb.Source.Width
			thumbnailHeight = thumb.Source.Height
		}
	}

	thumbnailURL = strings.ReplaceAll(thumbnailURL, "&amp;", "&")

	return &OEmbedResponse{
		Type:    "rich",
		Version: "1.0",
		Title:   post.Title,
		Extra: map[string]string{
			"image_alt": post.Title,
		},
		ThumbnailURL:    thumbnailURL,
		ThumbnailWidth:  thumbnailWidth,
		ThumbnailHeight: thumbnailHeight,
		ProviderName:    "Reddit",
		ProviderURL:     "https://reddit.com",
	}, nil
}

func fetchTikTokEmbed(resolver *oembedResolver, rawURL string) (*OEmbedResponse, error) {
	endpoint := "https://www.tiktok.com/oembed"

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("tiktok: parse endpoint: %w", err)
	}

	query := parsed.Query()
	query.Set("url", rawURL)
	parsed.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, parsed.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("tiktok: create request: %w", err)
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "application/json")

	resp, err := resolver.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tiktok: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tiktok: status %s", resp.Status)
	}

	var payload OEmbedResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("tiktok: decode: %w", err)
	}

	return &payload, nil
}

// fetchGitHubGistEmbed fetches embed HTML for GitHub Gists.
// GitHub doesn't have a working oEmbed endpoint, so we fetch the gist page
// and extract the embed script from it.
func fetchGitHubGistEmbed(resolver *oembedResolver, rawURL string) (*OEmbedResponse, error) {
	// Ensure we're getting a gist URL
	if !strings.Contains(rawURL, "gist.github.com") {
		return nil, fmt.Errorf("not a gist URL")
	}

	gistID, err := extractGistID(rawURL)
	if err != nil {
		return nil, fmt.Errorf("gist: %w", err)
	}

	apiURL := fmt.Sprintf("https://api.github.com/gists/%s", gistID)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, apiURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("gist: create request: %w", err)
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := resolver.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gist: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gist: status %s", resp.Status)
	}

	var gistData struct {
		Files map[string]struct {
			Filename  string `json:"filename"`
			Content   string `json:"content"`
			Type      string `json:"type"`
			Language  string `json:"language"`
			Size      int    `json:"size"`
			RawURL    string `json:"raw_url"`
			Truncated bool   `json:"truncated"`
		} `json:"files"`
		Description string `json:"description"`
		HTMLURL     string `json:"html_url"`
		Owner       struct {
			Login string `json:"login"`
		} `json:"owner"`
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&gistData); err != nil {
		return nil, fmt.Errorf("gist: decode: %w", err)
	}

	// Pick the first file by name for deterministic output
	var firstFile struct {
		Filename  string
		Content   string
		Language  string
		RawURL    string
		Size      int
		Truncated bool
	}
	if len(gistData.Files) > 0 {
		filenames := make([]string, 0, len(gistData.Files))
		for name := range gistData.Files {
			filenames = append(filenames, name)
		}
		sort.Strings(filenames)
		file := gistData.Files[filenames[0]]
		firstFile = struct {
			Filename  string
			Content   string
			Language  string
			RawURL    string
			Size      int
			Truncated bool
		}{
			Filename:  file.Filename,
			Content:   file.Content,
			Language:  file.Language,
			RawURL:    file.RawURL,
			Size:      file.Size,
			Truncated: file.Truncated,
		}
	}

	description := gistData.Description
	if description == "" {
		description = firstFile.Filename
	}

	// Generate embed HTML
	content := strings.TrimRight(firstFile.Content, "\n")
	if firstFile.RawURL != "" {
		fetched, err := fetchRemoteContent(resolver.client, firstFile.RawURL)
		if err == nil && fetched != "" {
			content = strings.TrimRight(fetched, "\n")
		}
	}
	codeHTML, err := renderGistCodeMarkdown(resolver, firstFile.Language, content)
	if err != nil || codeHTML == "" {
		embedHTML := fmt.Sprintf(`<script src="https://gist.github.com/%s.js?file=%s"></script>`,
			strings.TrimSuffix(strings.TrimPrefix(rawURL, "https://gist.github.com/"), ".json"),
			firstFile.Filename)
		return &OEmbedResponse{
			Type:         "rich",
			Version:      "1.0",
			Title:        description,
			HTML:         embedHTML,
			ProviderName: "GitHub",
			ProviderURL:  "https://github.com",
		}, nil
	}

	gistURL := strings.TrimSuffix(rawURL, ".json")
	if !strings.HasPrefix(gistURL, "https://gist.github.com/") {
		gistURL = rawURL
	}
	embedHTML := renderGistMarkdownEmbedHTML(gistURL, firstFile.Filename, firstFile.Language, codeHTML)

	return &OEmbedResponse{
		Type:         "rich",
		Version:      "1.0",
		Title:        description,
		HTML:         embedHTML,
		ProviderName: "GitHub",
		ProviderURL:  "https://github.com",
		Extra: map[string]string{
			"needs_code_css": BoolTrue,
		},
	}, nil
}

func fetchRemoteContent(client *http.Client, rawURL string) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, rawURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("fetch content: %w", err)
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch content: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("fetch content: %w", err)
	}

	return string(body), nil
}

func extractGistID(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}

	path := strings.TrimSuffix(parsed.Path, "/")
	parts := strings.Split(path, "/")
	for i := len(parts) - 1; i >= 0; i-- {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}
		part = strings.TrimSuffix(part, ".json")
		if part != "" {
			return part, nil
		}
	}

	return "", fmt.Errorf("missing gist id")
}

func renderGistScriptEmbedHTML(gistURL, filename string) string {
	if filename == "" {
		filename = "gist"
	}

	return fmt.Sprintf(`<script src="%s.js?file=%s"></script>`,
		strings.TrimSuffix(gistURL, ".json"),
		url.QueryEscape(filename))
}

var (
	giphyGifsRe  = regexp.MustCompile(`giphy\.com/gifs/[\w-]+-(\w+)`)
	giphyMediaRe = regexp.MustCompile(`media\.giphy\.com/media/(\w+)/`)
)

// fetchGiphyEmbed fetches GIPHY embed data using the oEmbed API.
func fetchGiphyEmbed(resolver *oembedResolver, rawURL string) (*OEmbedResponse, error) {
	endpoint := "https://giphy.com/services/oembed"
	parsed, err := url.Parse(endpoint)
	if err != nil {
		if fallback := buildGiphyFallback(rawURL); fallback != nil {
			return fallback, nil
		}
		return nil, fmt.Errorf("giphy: parse endpoint: %w", err)
	}

	query := parsed.Query()
	query.Set("url", rawURL)
	parsed.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, parsed.String(), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("giphy: create request: %w", err)
	}

	req.Header.Set("User-Agent", "markata-go/1.0 (+https://github.com/WaylonWalker/markata)")
	req.Header.Set("Accept", "application/json")

	resp, err := resolver.client.Do(req)
	if err != nil {
		if fallback := buildGiphyFallback(rawURL); fallback != nil {
			return fallback, nil
		}
		return nil, fmt.Errorf("giphy: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if fallback := buildGiphyFallback(rawURL); fallback != nil {
			return fallback, nil
		}
		return nil, fmt.Errorf("giphy: status %s", resp.Status)
	}

	var payload OEmbedResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("giphy: decode: %w", err)
	}

	// Ensure we have a thumbnail URL for image embeds
	if payload.ThumbnailURL == "" {
		if fallback := buildGiphyFallback(rawURL); fallback != nil {
			payload.ThumbnailURL = fallback.ThumbnailURL
			if payload.URL == "" {
				payload.URL = fallback.URL
			}
			if payload.Title == "" {
				payload.Title = fallback.Title
			}
			if payload.ProviderName == "" {
				payload.ProviderName = fallback.ProviderName
				payload.ProviderURL = fallback.ProviderURL
			}
		} else {
			payload.ThumbnailURL = payload.URL
		}
	}
	if payload.Extra == nil {
		payload.Extra = make(map[string]string)
	}
	if payload.Title != "" {
		payload.Extra["image_alt"] = payload.Title
	}

	return &payload, nil
}

func buildGiphyFallback(rawURL string) *OEmbedResponse {
	gifID := extractGiphyIDFromURL(rawURL)
	if gifID == "" {
		return nil
	}

	imageURL := fmt.Sprintf("https://media.giphy.com/media/%s/giphy.gif", gifID)
	return &OEmbedResponse{
		Type:         "photo",
		Version:      "1.0",
		Title:        "Giphy GIF",
		URL:          imageURL,
		ThumbnailURL: imageURL,
		ProviderName: "Giphy",
		ProviderURL:  "https://giphy.com",
		Extra: map[string]string{
			"image_alt": "Giphy GIF",
		},
	}
}

func extractGiphyIDFromURL(rawURL string) string {
	if matches := giphyGifsRe.FindStringSubmatch(rawURL); len(matches) > 1 {
		return matches[1]
	}
	if matches := giphyMediaRe.FindStringSubmatch(rawURL); len(matches) > 1 {
		return matches[1]
	}
	return ""
}
