package plugins

import (
	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
)

// SinglePageRootPlugin publishes every loaded post at the site root.
// It is used by single-file builds (markata-go build post.md) so the one
// source renders to output/index.html instead of output/<slug>/index.html.
type SinglePageRootPlugin struct{}

// NewSinglePageRootPlugin creates a new SinglePageRootPlugin.
func NewSinglePageRootPlugin() *SinglePageRootPlugin {
	return &SinglePageRootPlugin{}
}

// Name returns the unique name of the plugin.
func (p *SinglePageRootPlugin) Name() string {
	return "single_page_root"
}

// Priority runs after other Load plugins so the root slug wins over any
// frontmatter or path-derived slug before Transform plugins read hrefs.
func (p *SinglePageRootPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageLoad {
		return lifecycle.PriorityLast
	}
	return lifecycle.PriorityDefault
}

// Load rewrites loaded posts to the homepage slug and href.
func (p *SinglePageRootPlugin) Load(m *lifecycle.Manager) error {
	for _, post := range m.Posts() {
		if post == nil {
			continue
		}
		post.Slug = ""
		post.Set("_slug_explicit", true)
		post.GenerateHref()
	}
	return nil
}

// singlePageExcludedPlugins lists default plugins that generate site-level
// pages, feeds, or indexes that do not belong in a single-file build.
var singlePageExcludedPlugins = map[string]bool{
	"build_cache":             true,
	"blogroll":                true,
	"series":                  true,
	"subscription_feeds":      true,
	"feeds":                   true,
	"auto_feeds":              true,
	"prevnext":                true,
	"webmentions_leaderboard": true,
	"publish_feeds":           true,
	"well_known":              true,
	"images":                  true,
	"random_post":             true,
	"redirects":               true,
	"error_pages":             true,
	"tags_listing":            true,
	"feeds_listing":           true,
	"garden_view":             true,
	"sitemap":                 true,
	"llms_txt":                true,
	"content_index":           true,
	"pagefind":                true,
}

// SinglePagePlugins returns the default plugin set trimmed for rendering one
// Markdown file to output/index.html. Rendering, theme, and asset plugins are
// kept; feeds, listings, sitemaps, search indexes, and other site-wide pages
// are dropped.
func SinglePagePlugins() []lifecycle.Plugin {
	defaults := DefaultPlugins()
	result := make([]lifecycle.Plugin, 0, len(defaults)+1)
	for _, p := range defaults {
		if singlePageExcludedPlugins[p.Name()] {
			continue
		}
		result = append(result, p)
	}
	return append(result, NewSinglePageRootPlugin())
}

var (
	_ lifecycle.Plugin         = (*SinglePageRootPlugin)(nil)
	_ lifecycle.LoadPlugin     = (*SinglePageRootPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*SinglePageRootPlugin)(nil)
)
