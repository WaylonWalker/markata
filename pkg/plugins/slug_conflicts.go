// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// SlugConflict represents a conflict between content sources that resolve
// to the same slug.
type SlugConflict struct {
	// Slug is the conflicting slug (empty string means homepage)
	Slug string

	// Sources describes the conflicting content sources
	Sources []string

	// ConflictType indicates the type of conflict
	ConflictType string // "post-post", "post-feed"
}

// SlugConflictsPlugin detects slug conflicts between posts and feeds.
// It runs during the Collect stage after posts are loaded and feeds are configured.
//
// Conflicts detected:
//   - Multiple posts resolving to the same slug
//   - Post slug matching a feed slug
type SlugConflictsPlugin struct {
	// conflicts stores detected slug conflicts
	conflicts []SlugConflict

	// enabled controls whether the plugin runs
	enabled bool
}

// NewSlugConflictsPlugin creates a new SlugConflictsPlugin.
func NewSlugConflictsPlugin() *SlugConflictsPlugin {
	return &SlugConflictsPlugin{
		conflicts: make([]SlugConflict, 0),
		enabled:   true,
	}
}

// Name returns the unique name of the plugin.
func (p *SlugConflictsPlugin) Name() string {
	return "slug_conflicts"
}

// SetEnabled enables or disables the plugin.
func (p *SlugConflictsPlugin) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// Conflicts returns the detected slug conflicts.
func (p *SlugConflictsPlugin) Conflicts() []SlugConflict {
	return p.conflicts
}

// Priority returns the plugin execution priority for the given stage.
// Run very early in the collect stage to catch conflicts before other processing.
func (p *SlugConflictsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCollect {
		return lifecycle.PriorityFirst // Run first to catch conflicts early
	}
	return lifecycle.PriorityDefault
}

// Collect checks for slug conflicts between posts and feeds.
func (p *SlugConflictsPlugin) Collect(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}

	// Reset conflicts
	p.conflicts = make([]SlugConflict, 0)

	// Collect post and feed slugs
	postSlugs := p.collectPostSlugs(m)
	feedSlugs := p.collectFeedSlugs(m)

	// Detect conflicts
	p.detectPostPostConflicts(postSlugs)
	p.detectPostFeedConflicts(postSlugs, feedSlugs)

	// Sort conflicts for deterministic output
	sort.Slice(p.conflicts, func(i, j int) bool {
		return p.conflicts[i].Slug < p.conflicts[j].Slug
	})
	if m.ContentLedger() != nil {
		for _, conflict := range p.conflicts {
			for _, source := range conflict.Sources {
				path, ok := strings.CutPrefix(source, "post:")
				if ok {
					m.ContentLedger().AddReason(path, diagnostics.ReasonContentDuplicateSlug)
				}
			}
		}
	}

	if len(p.conflicts) == 0 {
		return nil
	}

	return p.formatError()
}

// collectPostSlugs gathers slugs from all publishable posts.
func (p *SlugConflictsPlugin) collectPostSlugs(m *lifecycle.Manager) map[string][]string {
	postSlugs := make(map[string][]string) // slug -> list of post paths
	for _, post := range m.Posts() {
		// Skip posts that won't be written
		if post.Skip || post.Draft {
			continue
		}
		normalizedSlug := strings.ToLower(post.Slug)
		postSlugs[normalizedSlug] = append(postSlugs[normalizedSlug], post.Path)
	}
	return postSlugs
}

// collectFeedSlugs gathers slugs from feeds that generate HTML output.
func (p *SlugConflictsPlugin) collectFeedSlugs(m *lifecycle.Manager) map[string]string {
	feedSlugs := make(map[string]string) // slug -> feed identifier
	feedConfigs := p.getFeedConfigs(m)
	feedDefaults := p.getFeedDefaults(m)

	for i := range feedConfigs {
		fc := feedConfigs[i]
		// Apply defaults to check if HTML format will be enabled
		// (FeedsPlugin applies these later, but we need to check now)
		if !fc.Formats.HasAnyEnabled() {
			fc.Formats = feedDefaults.Formats
		}
		// Only check feeds that generate HTML (which would conflict with post index.html)
		if fc.Formats.HTML {
			feedSlugs[strings.ToLower(fc.Slug)] = fmt.Sprintf("feed:%s", fc.Slug)
		}
	}
	return feedSlugs
}

// getFeedConfigs retrieves feed configurations from cache or config.
func (p *SlugConflictsPlugin) getFeedConfigs(m *lifecycle.Manager) []models.FeedConfig {
	// Try cache first
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			return fcs
		}
	}

	// Fall back to config.Extra
	if m.Config().Extra != nil {
		if feeds, ok := m.Config().Extra["feeds"]; ok {
			if fcs, ok := feeds.([]models.FeedConfig); ok {
				return fcs
			}
		}
	}
	return nil
}

// getFeedDefaults retrieves feed defaults from config.
func (p *SlugConflictsPlugin) getFeedDefaults(m *lifecycle.Manager) models.FeedDefaults {
	feedDefaults := models.NewFeedDefaults()
	if m.Config().Extra != nil {
		if defaults, ok := m.Config().Extra["feed_defaults"]; ok {
			if fd, ok := defaults.(models.FeedDefaults); ok {
				return fd
			}
		}
	}
	return feedDefaults
}

// detectPostPostConflicts finds multiple posts with the same slug.
func (p *SlugConflictsPlugin) detectPostPostConflicts(postSlugs map[string][]string) {
	for slug, posts := range postSlugs {
		if len(posts) > 1 {
			p.conflicts = append(p.conflicts, SlugConflict{
				Slug:         slug,
				Sources:      formatPostSources(posts),
				ConflictType: "post-post",
			})
		}
	}
}

// detectPostFeedConflicts finds posts whose slugs match feed slugs.
func (p *SlugConflictsPlugin) detectPostFeedConflicts(postSlugs map[string][]string, feedSlugs map[string]string) {
	for slug, posts := range postSlugs {
		if feedID, exists := feedSlugs[slug]; exists {
			sources := append(formatPostSources(posts), feedID)
			p.conflicts = append(p.conflicts, SlugConflict{
				Slug:         slug,
				Sources:      sources,
				ConflictType: "post-feed",
			})
		}
	}
}

// formatPostSources formats post paths as source identifiers.
func formatPostSources(paths []string) []string {
	sources := make([]string, len(paths))
	for i, path := range paths {
		sources[i] = fmt.Sprintf("post:%s", path)
	}
	return sources
}

// formatError builds a detailed error message for all conflicts.
func (p *SlugConflictsPlugin) formatError() error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("ERROR: %d slug conflict(s) detected\n\n", len(p.conflicts)))

	for _, c := range p.conflicts {
		slugDisplay := c.Slug
		if slugDisplay == "" {
			slugDisplay = "(homepage - empty slug)"
		}

		sb.WriteString(fmt.Sprintf("  Conflicting slug: %s\n", slugDisplay))
		sb.WriteString(fmt.Sprintf("    Type: %s\n", c.ConflictType))
		sb.WriteString("    Sources:\n")
		for _, src := range c.Sources {
			sb.WriteString(fmt.Sprintf("      - %s\n", src))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("To fix:\n")
	sb.WriteString("  - Change the slug in one of the conflicting posts' frontmatter\n")
	sb.WriteString("  - Remove one of the conflicting content sources\n")
	sb.WriteString("  - For post-feed conflicts, use a different slug for the feed\n")

	return &SlugConflictError{Message: sb.String(), Conflicts: p.conflicts}
}

// SlugConflictError is a critical error about slug conflicts.
// It implements lifecycle.CriticalError to ensure the build fails.
type SlugConflictError struct {
	Message   string
	Conflicts []SlugConflict
}

// Error implements the error interface.
func (e *SlugConflictError) Error() string {
	return e.Message
}

// IsCritical marks this as a critical error that should halt the build.
func (e *SlugConflictError) IsCritical() bool {
	return true
}

// Ensure SlugConflictsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*SlugConflictsPlugin)(nil)
	_ lifecycle.CollectPlugin  = (*SlugConflictsPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*SlugConflictsPlugin)(nil)
)
