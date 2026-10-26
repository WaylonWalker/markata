// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"path"
	"strings"
	"unicode"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// BreadcrumbsPlugin generates breadcrumb navigation for posts based on URL path.
// It supports auto-generation from directory structure, manual override via
// frontmatter, and produces JSON-LD structured data for SEO.
type BreadcrumbsPlugin struct {
	// showHome includes the home link as the first breadcrumb
	showHome bool
	// homeLabel is the label for the home breadcrumb
	homeLabel string
	// separator is the visual separator between breadcrumbs (for display)
	separator string
	// maxDepth limits the breadcrumb depth (0 = unlimited)
	maxDepth int
	// structuredData enables JSON-LD BreadcrumbList generation
	structuredData bool
}

// Breadcrumb represents a single item in the breadcrumb trail.
type Breadcrumb struct {
	// Label is the display text for the breadcrumb
	Label string `json:"label"`
	// URL is the href for the breadcrumb link
	URL string `json:"url"`
	// IsCurrent indicates if this is the current page (last item)
	IsCurrent bool `json:"is_current"`
	// Position is the 1-indexed position in the trail (for JSON-LD)
	Position int `json:"position"`
}

// BreadcrumbConfig holds per-post breadcrumb configuration from frontmatter.
type BreadcrumbConfig struct {
	// Enabled controls whether breadcrumbs are shown (nil = use default)
	Enabled *bool `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	// Items allows manual override of breadcrumb trail
	Items []BreadcrumbItem `json:"items,omitempty" yaml:"items,omitempty"`
	// ShowHome overrides the global show_home setting
	ShowHome *bool `json:"show_home,omitempty" yaml:"show_home,omitempty"`
	// HomeLabel overrides the global home_label
	HomeLabel string `json:"home_label,omitempty" yaml:"home_label,omitempty"`
}

// BreadcrumbItem is a manual breadcrumb entry from frontmatter.
type BreadcrumbItem struct {
	Label string `json:"label" yaml:"label"`
	URL   string `json:"url" yaml:"url"`
}

// BreadcrumbListJSON is the JSON-LD structured data for breadcrumbs.
type BreadcrumbListJSON struct {
	Context         string               `json:"@context"`
	Type            string               `json:"@type"`
	ItemListElement []BreadcrumbListItem `json:"itemListElement"`
}

// BreadcrumbListItem is a single item in JSON-LD BreadcrumbList.
type BreadcrumbListItem struct {
	Type     string `json:"@type"`
	Position int    `json:"position"`
	Name     string `json:"name"`
	Item     string `json:"item,omitempty"`
}

// NewBreadcrumbsPlugin creates a new BreadcrumbsPlugin with default settings.
func NewBreadcrumbsPlugin() *BreadcrumbsPlugin {
	return &BreadcrumbsPlugin{
		showHome:       true,
		homeLabel:      "Home",
		separator:      "/",
		maxDepth:       0, // unlimited
		structuredData: true,
	}
}

// Name returns the unique name of the plugin.
func (p *BreadcrumbsPlugin) Name() string {
	return "breadcrumbs"
}

// Configure reads configuration options for the plugin.
func (p *BreadcrumbsPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra == nil {
		return nil
	}

	// Check for breadcrumbs config in components section
	if components, ok := config.Extra["components"].(map[string]interface{}); ok {
		if bcConfig, ok := components["breadcrumbs"].(map[string]interface{}); ok {
			p.configureFromMap(bcConfig)
		}
	}

	// Also check top-level breadcrumbs config
	if bcConfig, ok := config.Extra["breadcrumbs"].(map[string]interface{}); ok {
		p.configureFromMap(bcConfig)
	}

	return nil
}

// configureFromMap applies configuration from a map.
func (p *BreadcrumbsPlugin) configureFromMap(config map[string]interface{}) {
	if enabled, ok := config["enabled"].(bool); ok && !enabled {
		// Plugin disabled entirely - this will be checked in Transform
		return
	}
	if showHome, ok := config["show_home"].(bool); ok {
		p.showHome = showHome
	}
	if homeLabel, ok := config["home_label"].(string); ok && homeLabel != "" {
		p.homeLabel = homeLabel
	}
	if separator, ok := config["separator"].(string); ok && separator != "" {
		p.separator = separator
	}
	if maxDepth, ok := parseIntFromInterface(config["max_depth"]); ok && maxDepth >= 0 {
		p.maxDepth = maxDepth
	}
	if sd, ok := config["structured_data"].(bool); ok {
		p.structuredData = sd
	}
}

// Transform generates breadcrumbs for each post.
func (p *BreadcrumbsPlugin) Transform(m *lifecycle.Manager) error {
	config := m.Config()

	// Check if plugin is globally disabled
	if p.isDisabled(config) {
		return nil
	}

	siteURL := getSiteURL(config)

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip {
			return false
		}
		postConfig := p.getPostConfig(post)
		if postConfig.Enabled != nil && !*postConfig.Enabled {
			return false
		}
		return true
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		postConfig := p.getPostConfig(post)
		if postConfig.Enabled != nil && !*postConfig.Enabled {
			return nil
		}

		// Generate breadcrumbs
		breadcrumbs := p.generateBreadcrumbs(post, postConfig, siteURL)
		if len(breadcrumbs) == 0 {
			return nil
		}

		// Store breadcrumbs in post.Extra
		post.Set("breadcrumbs", breadcrumbs)
		post.Set("breadcrumb_separator", p.separator)

		// Generate JSON-LD if enabled
		if p.structuredData {
			jsonLD := p.generateJSONLD(breadcrumbs, siteURL)
			post.Set("breadcrumbs_jsonld", jsonLD)
		}

		return nil
	})
}

// isDisabled checks if the plugin is globally disabled.
func (p *BreadcrumbsPlugin) isDisabled(config *lifecycle.Config) bool {
	if config.Extra == nil {
		return false
	}

	// Check components.breadcrumbs.enabled
	if components, ok := config.Extra["components"].(map[string]interface{}); ok {
		if bcConfig, ok := components["breadcrumbs"].(map[string]interface{}); ok {
			if enabled, ok := bcConfig["enabled"].(bool); ok && !enabled {
				return true
			}
		}
	}

	// Check breadcrumbs.enabled
	if bcConfig, ok := config.Extra["breadcrumbs"].(map[string]interface{}); ok {
		if enabled, ok := bcConfig["enabled"].(bool); ok && !enabled {
			return true
		}
	}

	return false
}

// getPostConfig retrieves per-post breadcrumb configuration from frontmatter.
func (p *BreadcrumbsPlugin) getPostConfig(post *models.Post) BreadcrumbConfig {
	config := BreadcrumbConfig{}

	if post.Extra == nil {
		return config
	}

	// Check for breadcrumbs: false to disable
	if bc, ok := post.Extra["breadcrumbs"]; ok {
		switch v := bc.(type) {
		case bool:
			config.Enabled = &v
		case map[string]interface{}:
			p.parsePostConfigMap(v, &config)
		}
	}

	return config
}

// parsePostConfigMap parses breadcrumb configuration from a map.
func (p *BreadcrumbsPlugin) parsePostConfigMap(m map[string]interface{}, config *BreadcrumbConfig) {
	if enabled, ok := m["enabled"].(bool); ok {
		config.Enabled = &enabled
	}
	if showHome, ok := m["show_home"].(bool); ok {
		config.ShowHome = &showHome
	}
	if homeLabel, ok := m["home_label"].(string); ok {
		config.HomeLabel = homeLabel
	}

	// Parse manual items
	if items, ok := m["items"].([]interface{}); ok {
		for _, item := range items {
			if itemMap, ok := item.(map[string]interface{}); ok {
				bi := BreadcrumbItem{}
				if label, ok := itemMap["label"].(string); ok {
					bi.Label = label
				}
				if url, ok := itemMap["url"].(string); ok {
					bi.URL = url
				}
				if bi.Label != "" {
					config.Items = append(config.Items, bi)
				}
			}
		}
	}
}

// generateBreadcrumbs creates the breadcrumb trail for a post.
func (p *BreadcrumbsPlugin) generateBreadcrumbs(post *models.Post, postConfig BreadcrumbConfig, siteURL string) []Breadcrumb {
	// If manual items provided, use those
	if len(postConfig.Items) > 0 {
		return p.buildFromManualItems(postConfig.Items, post, postConfig)
	}

	// Auto-generate from URL path
	return p.buildFromPath(post, postConfig, siteURL)
}

// buildFromManualItems creates breadcrumbs from frontmatter items.
func (p *BreadcrumbsPlugin) buildFromManualItems(items []BreadcrumbItem, post *models.Post, postConfig BreadcrumbConfig) []Breadcrumb {
	breadcrumbs := make([]Breadcrumb, 0, len(items)+2)
	position := 1

	// Determine if we should show home
	showHome := p.showHome
	if postConfig.ShowHome != nil {
		showHome = *postConfig.ShowHome
	}

	// Add home if enabled and not already first item
	if showHome && (len(items) == 0 || items[0].URL != "/") {
		homeLabel := p.homeLabel
		if postConfig.HomeLabel != "" {
			homeLabel = postConfig.HomeLabel
		}
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Label:     homeLabel,
			URL:       "/",
			IsCurrent: false,
			Position:  position,
		})
		position++
	}

	// Add manual items
	for i, item := range items {
		isLast := i == len(items)-1
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Label:     item.Label,
			URL:       item.URL,
			IsCurrent: isLast,
			Position:  position,
		})
		position++
	}

	// If no current page in manual items, add the post
	if len(items) == 0 || items[len(items)-1].URL != post.Href {
		title := p.getPostTitle(post)
		if title != "" {
			breadcrumbs = append(breadcrumbs, Breadcrumb{
				Label:     title,
				URL:       post.Href,
				IsCurrent: true,
				Position:  position,
			})
		}
	}

	return breadcrumbs
}

// buildFromPath auto-generates breadcrumbs from the post's URL path.
func (p *BreadcrumbsPlugin) buildFromPath(post *models.Post, postConfig BreadcrumbConfig, _ string) []Breadcrumb {
	href := post.Href

	// Handle homepage
	if href == "/" || href == "" {
		return nil // No breadcrumbs for homepage
	}

	// Clean the path
	href = strings.TrimPrefix(href, "/")
	href = strings.TrimSuffix(href, "/")

	parts := strings.Split(href, "/")
	if len(parts) == 0 {
		return nil
	}

	breadcrumbs := make([]Breadcrumb, 0, len(parts)+1)
	position := 1

	// Determine if we should show home
	showHome := p.showHome
	if postConfig.ShowHome != nil {
		showHome = *postConfig.ShowHome
	}

	// Add home breadcrumb if enabled
	if showHome {
		homeLabel := p.homeLabel
		if postConfig.HomeLabel != "" {
			homeLabel = postConfig.HomeLabel
		}
		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Label:     homeLabel,
			URL:       "/",
			IsCurrent: false,
			Position:  position,
		})
		position++
	}

	// Build intermediate breadcrumbs from path segments
	currentPath := ""
	for i, part := range parts {
		currentPath = path.Join(currentPath, part)
		isLast := i == len(parts)-1

		// Check max depth
		if p.maxDepth > 0 && position > p.maxDepth {
			break
		}

		// For the last segment, use the post title if available
		label := p.humanizeSegment(part)
		if isLast {
			title := p.getPostTitle(post)
			if title != "" {
				label = title
			}
		}

		breadcrumbs = append(breadcrumbs, Breadcrumb{
			Label:     label,
			URL:       "/" + currentPath + "/",
			IsCurrent: isLast,
			Position:  position,
		})
		position++
	}

	return breadcrumbs
}

// humanizeSegment converts a URL segment to a human-readable label.
func (p *BreadcrumbsPlugin) humanizeSegment(segment string) string {
	// Replace hyphens and underscores with spaces
	label := strings.ReplaceAll(segment, "-", " ")
	label = strings.ReplaceAll(label, "_", " ")

	// Title case
	upperNext := true
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			upperNext = true
			return r
		}
		if upperNext {
			upperNext = false
			return unicode.ToTitle(r)
		}
		return r
	}, label)
}

// getPostTitle returns the post's title for display.
func (p *BreadcrumbsPlugin) getPostTitle(post *models.Post) string {
	if post.PlainTitle() != "" {
		return post.PlainTitle()
	}
	return ""
}

// generateJSONLD creates JSON-LD structured data for breadcrumbs.
func (p *BreadcrumbsPlugin) generateJSONLD(breadcrumbs []Breadcrumb, siteURL string) string {
	if len(breadcrumbs) == 0 {
		return ""
	}

	jsonLD := BreadcrumbListJSON{
		Context:         "https://schema.org",
		Type:            "BreadcrumbList",
		ItemListElement: make([]BreadcrumbListItem, 0, len(breadcrumbs)),
	}

	siteURL = strings.TrimSuffix(siteURL, "/")

	for _, bc := range breadcrumbs {
		item := BreadcrumbListItem{
			Type:     "ListItem",
			Position: bc.Position,
			Name:     bc.Label,
		}

		// Only include item URL for non-current items
		if !bc.IsCurrent {
			item.Item = siteURL + bc.URL
		}

		jsonLD.ItemListElement = append(jsonLD.ItemListElement, item)
	}

	jsonBytes, err := json.MarshalIndent(jsonLD, "", "  ")
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// Priority returns the plugin priority for the given stage.
// Breadcrumbs should run after auto_title so titles are available.
func (p *BreadcrumbsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityDefault // After auto_title (PriorityEarly)
	}
	return lifecycle.PriorityDefault
}

// Ensure BreadcrumbsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*BreadcrumbsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*BreadcrumbsPlugin)(nil)
	_ lifecycle.TransformPlugin = (*BreadcrumbsPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*BreadcrumbsPlugin)(nil)
)
