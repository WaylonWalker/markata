package plugins

import (
	"fmt"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// JinjaMdPlugin processes Jinja2 templates in markdown content.
// This allows using template syntax within markdown files before they are
// converted to HTML. It operates during the transform stage.
//
// Posts must have `jinja: true` in their frontmatter to be processed.
//
// Available template variables:
//   - post: The current post object
//   - config: The site configuration
//   - posts: All posts (via core.Posts())
//   - core: The lifecycle manager (for filter/map operations)
type JinjaMdPlugin struct {
	engine *templates.Engine
}

// NewJinjaMdPlugin creates a new jinja_md plugin.
func NewJinjaMdPlugin() *JinjaMdPlugin {
	return &JinjaMdPlugin{}
}

// Name returns the plugin name.
func (p *JinjaMdPlugin) Name() string {
	return "jinja_md"
}

// Configure initializes the template engine.
// If the templates plugin has already initialized an engine, reuse it.
func (p *JinjaMdPlugin) Configure(m *lifecycle.Manager) error {
	// Try to get existing engine from cache (set by templates plugin)
	if cached, ok := m.Cache().Get("templates.engine"); ok {
		if engine, ok := cached.(*templates.Engine); ok {
			p.engine = engine
			return nil
		}
	}

	// Create our own engine (without templates directory, we only need string rendering)
	engine, err := templates.NewEngine("")
	if err != nil {
		return fmt.Errorf("failed to initialize template engine: %w", err)
	}
	p.engine = engine

	return nil
}

// Transform processes jinja templates in markdown content.
// Only posts with `jinja: true` or `jinja_md: true` in their frontmatter are processed.
func (p *JinjaMdPlugin) Transform(m *lifecycle.Manager) error {
	if p.engine == nil {
		return fmt.Errorf("template engine not initialized")
	}

	// Get config for template context
	config := m.Config()

	// Get all posts for access in templates
	allPosts := m.Posts()

	// Collect private paths for robots.txt and similar templates
	privatePaths := collectPrivatePathsForJinja(allPosts)

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.Content == "" {
			return false
		}
		return isJinjaEnabled(post)
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		// Create template context
		ctx := templates.NewContext(post, "", ToModelsConfig(config))
		ctx = ctx.WithCore(m)
		ctx = ctx.WithPosts(allPosts)

		// Add all_posts alias for convenience (same as posts)
		ctx.Set("all_posts", templates.PostsToMaps(allPosts))

		// Add helper functions as extra context
		ctx.Set("filter", createFilterFunc(m))
		ctx.Set("map", createMapFunc(m))
		ctx.Set("feed_posts", createFeedPostsFunc(m))
		ctx.Set("render_feed", createRenderFeedFunc(m))
		ctx.Set("render_slashes", createRenderSlashesFunc(m))
		ctx.Set("include_post", createIncludePostFunc(m))

		// Add private_paths for robots.txt generation
		ctx.Set("private_paths", privatePaths)

		// Render the content as a template
		rendered, err := p.engine.RenderString(post.Content, ctx)
		if err != nil {
			return fmt.Errorf("failed to render jinja in post %q: %w", post.Path, err)
		}

		// Update the post content with rendered result
		post.Content = rendered
		return nil
	})
}

// collectPrivatePathsForJinja returns a list of paths (hrefs) for all private posts.
// These paths are used in robots.txt templates to add Disallow directives.
// Includes all format variants (.txt, .ansi, .md, .og) and excludes the robots post itself.
func collectPrivatePathsForJinja(posts []*models.Post) []string {
	var paths []string
	for _, post := range posts {
		if post.Private && !post.Draft && !post.Skip {
			// Skip the robots post itself to avoid self-reference
			if post.Slug == "robots" {
				continue
			}
			// Add base href (e.g., /slug/)
			paths = append(paths, post.Href)
			// Add format variants
			// For regular posts: /slug.txt, /slug.ansi, /slug.md, /slug.og/
			if post.Slug != "" {
				paths = append(paths,
					"/"+post.Slug+".txt",
					"/"+post.Slug+".ansi",
					"/"+post.Slug+".md",
					"/"+post.Slug+".og/",
				)
			}
		}
	}
	return paths
}

// Priority returns the plugin priority for the given stage.
// JinjaMd should run early in the transform stage, before other transforms.
func (p *JinjaMdPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityEarly // Run early to process templates before other transforms
	}
	return lifecycle.PriorityDefault
}

// isJinjaEnabled checks if jinja processing is enabled for a post.
// Returns true if the post has `jinja: true` or `jinja_md: true` in frontmatter or Extra.
func isJinjaEnabled(post *models.Post) bool {
	if post.Extra == nil {
		return false
	}

	// Check for "jinja" or "jinja_md" key (both are valid)
	for _, key := range []string{"jinja", "jinja_md"} {
		if val, ok := post.Extra[key]; ok {
			switch v := val.(type) {
			case bool:
				if v {
					return true
				}
			case string:
				if v == BoolTrue || v == "yes" || v == "1" {
					return true
				}
			}
		}
	}

	return false
}

// filterFuncWrapper wraps the Manager.Filter method for use in templates.
type filterFuncWrapper struct {
	manager *lifecycle.Manager
}

// Call implements a callable for filter operations.
func (f *filterFuncWrapper) Call(expr string) ([]*models.Post, error) {
	return f.manager.Filter(expr)
}

// createFilterFunc creates a filter function for templates.
// Usage in templates: {% for post in filter("published==true") %}
func createFilterFunc(m *lifecycle.Manager) *filterFuncWrapper {
	return &filterFuncWrapper{manager: m}
}

// mapFuncWrapper wraps the Manager.Map method for use in templates.
type mapFuncWrapper struct {
	manager *lifecycle.Manager
}

// Call implements a callable for map operations.
func (mf *mapFuncWrapper) Call(field, filter, sort string, reverse bool) ([]interface{}, error) {
	return mf.manager.Map(field, filter, sort, reverse)
}

// createMapFunc creates a map function for templates.
// Usage in templates: {{ map("title", "published==true", "date", true) }}
func createMapFunc(m *lifecycle.Manager) *mapFuncWrapper {
	return &mapFuncWrapper{manager: m}
}

// Engine returns the template engine for use by other plugins.
func (p *JinjaMdPlugin) Engine() *templates.Engine {
	return p.engine
}
