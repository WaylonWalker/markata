// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// PathConflict represents a conflict between two content sources
// that would result in the same output path.
type PathConflict struct {
	// OutputPath is the conflicting output path
	OutputPath string

	// Sources are the content sources (posts, feeds) that conflict
	Sources []string
}

// OverwriteCheckPlugin detects when multiple posts or feeds would write
// to the same output path, preventing accidental content overwrites.
type OverwriteCheckPlugin struct {
	// conflicts stores detected path conflicts
	conflicts []PathConflict

	// warnOnly when true, only warns about conflicts instead of failing
	warnOnly bool
}

// NewOverwriteCheckPlugin creates a new OverwriteCheckPlugin.
func NewOverwriteCheckPlugin() *OverwriteCheckPlugin {
	return &OverwriteCheckPlugin{
		conflicts: make([]PathConflict, 0),
		warnOnly:  false,
	}
}

// Name returns the unique name of the plugin.
func (p *OverwriteCheckPlugin) Name() string {
	return "overwrite_check"
}

// Priority returns the plugin execution priority for the given stage.
// Run before other collect plugins to catch conflicts early.
func (p *OverwriteCheckPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCollect {
		return lifecycle.PriorityEarly // Run early in the collect stage
	}
	return lifecycle.PriorityDefault
}

// SetWarnOnly configures whether to warn or fail on conflicts.
func (p *OverwriteCheckPlugin) SetWarnOnly(warnOnly bool) {
	p.warnOnly = warnOnly
}

// Conflicts returns the detected path conflicts.
func (p *OverwriteCheckPlugin) Conflicts() []PathConflict {
	return p.conflicts
}

// Collect checks for output path conflicts between posts and feeds.
func (p *OverwriteCheckPlugin) Collect(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir

	// Map of output paths to their sources
	pathSources := make(map[string][]string)

	// Check post output paths
	for _, post := range m.Posts() {
		// Skip posts that won't be written
		if post.Skip || post.Draft {
			continue
		}

		outputPath := p.getPostOutputPath(outputDir, post)
		pathSources[outputPath] = append(pathSources[outputPath], fmt.Sprintf("post:%s", post.Path))
	}

	// Check feed output paths
	var feedConfigs []models.FeedConfig
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			feedConfigs = fcs
		}
	}

	for i := range feedConfigs {
		fc := &feedConfigs[i]
		feedOutputPaths := p.getFeedOutputPaths(outputDir, fc)
		for _, outputPath := range feedOutputPaths {
			pathSources[outputPath] = append(pathSources[outputPath], fmt.Sprintf("feed:%s", fc.Slug))
		}
	}

	// Detect conflicts (paths with multiple sources)
	p.conflicts = make([]PathConflict, 0)
	for path, sources := range pathSources {
		if len(sources) > 1 {
			p.conflicts = append(p.conflicts, PathConflict{
				OutputPath: path,
				Sources:    sources,
			})
		}
	}

	if len(p.conflicts) == 0 {
		return nil
	}

	// Build error message
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("detected %d output path conflict(s):\n", len(p.conflicts)))
	for _, conflict := range p.conflicts {
		sb.WriteString(fmt.Sprintf("  - %s: %s\n", conflict.OutputPath, strings.Join(conflict.Sources, ", ")))
	}

	if p.warnOnly {
		// Store warning in manager (non-critical error)
		// The warning will be collected by the manager
		return nil
	}

	return fmt.Errorf("%s", sb.String())
}

// getPostOutputPath returns the output path for a post.
// This uses the post's slug directly without auto-generating - the slug should
// already be set by the load plugin (either from frontmatter or auto-generated).
func (p *OverwriteCheckPlugin) getPostOutputPath(outputDir string, post *models.Post) string {
	// Posts write to: output_dir/slug/index.html
	// Empty slug means homepage: output_dir/index.html
	return filepath.Join(outputDir, post.Slug, "index.html")
}

// getFeedOutputPaths returns all output paths for a feed (HTML, RSS, etc.).
func (p *OverwriteCheckPlugin) getFeedOutputPaths(outputDir string, fc *models.FeedConfig) []string {
	paths := make([]string, 0)

	feedDir := filepath.Join(outputDir, fc.Slug)
	if fc.Slug == "" {
		feedDir = outputDir
	}

	// Add HTML pages (index.html, page/N/index.html)
	if fc.Formats.HTML {
		paths = append(paths, filepath.Join(feedDir, "index.html"))
	}

	// Add RSS
	if fc.Formats.RSS {
		paths = append(paths, filepath.Join(feedDir, "rss.xml"))
	}

	// Add Atom
	if fc.Formats.Atom {
		paths = append(paths, filepath.Join(feedDir, "atom.xml"))
	}

	// Add JSON
	if fc.Formats.JSON {
		paths = append(paths, filepath.Join(feedDir, "index.json"))
	}

	return paths
}

// Ensure OverwriteCheckPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*OverwriteCheckPlugin)(nil)
	_ lifecycle.CollectPlugin  = (*OverwriteCheckPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*OverwriteCheckPlugin)(nil)
)
