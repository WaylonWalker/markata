// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// StaticFileConflict represents a conflict between a generated file and a static file.
type StaticFileConflict struct {
	// GeneratedSource is the source of the generated file (e.g., "pages/robots.md")
	GeneratedSource string

	// GeneratedOutput is the output path of the generated file (e.g., "/robots.txt")
	GeneratedOutput string

	// StaticFile is the path to the conflicting static file (e.g., "static/robots.txt")
	StaticFile string

	// OutputPath is the final output path that both would write to
	OutputPath string
}

// StaticFileConflictsPlugin detects when static files would clobber generated content.
// This is a lint rule that warns users when they have both:
// - A generated file (e.g., robots.md → robots.txt)
// - A static file (e.g., static/robots.txt)
//
// The static file always wins (copied last), which can cause unexpected behavior
// like private posts not being added to robots.txt.
type StaticFileConflictsPlugin struct {
	// conflicts stores detected conflicts
	conflicts []StaticFileConflict

	// staticDir is the directory containing static files
	staticDir string

	// enabled controls whether the plugin runs
	enabled bool
}

// NewStaticFileConflictsPlugin creates a new StaticFileConflictsPlugin.
func NewStaticFileConflictsPlugin() *StaticFileConflictsPlugin {
	return &StaticFileConflictsPlugin{
		conflicts: make([]StaticFileConflict, 0),
		staticDir: "static",
		enabled:   true,
	}
}

// Name returns the unique name of the plugin.
func (p *StaticFileConflictsPlugin) Name() string {
	return "static_file_conflicts"
}

// SetStaticDir sets the static directory path.
func (p *StaticFileConflictsPlugin) SetStaticDir(dir string) {
	p.staticDir = dir
}

// SetEnabled enables or disables the plugin.
func (p *StaticFileConflictsPlugin) SetEnabled(enabled bool) {
	p.enabled = enabled
}

// Conflicts returns the detected conflicts.
func (p *StaticFileConflictsPlugin) Conflicts() []StaticFileConflict {
	return p.conflicts
}

// Priority returns the plugin execution priority for the given stage.
// Run during collect stage after posts are processed but before writing.
func (p *StaticFileConflictsPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageCollect {
		return lifecycle.PriorityLate // Run late to catch all generated content
	}
	return lifecycle.PriorityDefault
}

// Collect checks for conflicts between generated and static files.
func (p *StaticFileConflictsPlugin) Collect(m *lifecycle.Manager) error {
	if !p.enabled {
		return nil
	}

	// Check if static directory exists
	if _, err := os.Stat(p.staticDir); os.IsNotExist(err) {
		return nil // No static directory, no conflicts possible
	}

	// Build a map of static files
	staticFiles := make(map[string]string) // normalized output path -> static file path
	err := filepath.Walk(p.staticDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Get relative path from static dir
		relPath, err := filepath.Rel(p.staticDir, path)
		if err != nil {
			return err
		}

		// Normalize to output path (static/foo.txt -> /foo.txt)
		outputPath := "/" + filepath.ToSlash(relPath)
		staticFiles[outputPath] = path

		return nil
	})
	if err != nil {
		return fmt.Errorf("scanning static directory: %w", err)
	}

	// Reset conflicts
	p.conflicts = make([]StaticFileConflict, 0)

	// Check each post for potential conflicts
	for _, post := range m.Posts() {
		// Skip posts that won't be written
		if post.Skip || post.Draft {
			continue
		}

		// Check for various output formats this post might generate
		conflicts := p.checkPostConflicts(post, staticFiles)
		p.conflicts = append(p.conflicts, conflicts...)
	}

	// Check feed configs for conflicts
	feedConflicts := p.checkFeedConflicts(m, staticFiles)
	p.conflicts = append(p.conflicts, feedConflicts...)

	// Report conflicts as warnings
	if len(p.conflicts) > 0 {
		return p.reportConflicts()
	}

	return nil
}

// checkPostConflicts checks if a post's generated files conflict with static files.
func (p *StaticFileConflictsPlugin) checkPostConflicts(post *models.Post, staticFiles map[string]string) []StaticFileConflict {
	conflicts := make([]StaticFileConflict, 0)

	// Get the base filename without extension
	base := filepath.Base(post.Path)
	ext := filepath.Ext(base)
	stem := strings.TrimSuffix(base, ext)

	// Common patterns for root-level generated files
	// robots.md → /robots.txt
	// sitemap.md → /sitemap.xml
	// humans.md → /humans.txt
	// security.md → /security.txt (or /.well-known/security.txt)

	// Check if this is a root-level special file
	dir := filepath.Dir(post.Path)
	isRootLevel := dir == "." || dir == "" || dir == "pages"

	if !isRootLevel {
		return conflicts
	}

	// Map of stem names to their expected output extensions
	specialFiles := map[string][]string{
		"robots":   {".txt"},
		"sitemap":  {".xml"},
		"humans":   {".txt"},
		"security": {".txt"},
		"manifest": {".json", ".webmanifest"},
		"feed":     {".xml", ".json", ".atom"},
		"rss":      {".xml"},
		"atom":     {".xml"},
	}

	stemLower := strings.ToLower(stem)
	if extensions, ok := specialFiles[stemLower]; ok {
		for _, outputExt := range extensions {
			outputPath := "/" + stemLower + outputExt
			if staticPath, exists := staticFiles[outputPath]; exists {
				conflicts = append(conflicts, StaticFileConflict{
					GeneratedSource: post.Path,
					GeneratedOutput: outputPath,
					StaticFile:      staticPath,
					OutputPath:      outputPath,
				})
			}
		}
	}

	// Also check for .txt files that might be generated from .md
	// e.g., changelog.md → /changelog.txt (if txt output is enabled)
	txtOutputPath := "/" + stemLower + ".txt"
	if staticPath, exists := staticFiles[txtOutputPath]; exists {
		// Check if this post has txt template configured
		if post.Templates != nil {
			if _, hasTxt := post.Templates["txt"]; hasTxt {
				conflicts = append(conflicts, StaticFileConflict{
					GeneratedSource: post.Path,
					GeneratedOutput: txtOutputPath,
					StaticFile:      staticPath,
					OutputPath:      txtOutputPath,
				})
			}
		}
	}

	return conflicts
}

// checkFeedConflicts checks if feed-generated files conflict with static files.
func (p *StaticFileConflictsPlugin) checkFeedConflicts(m *lifecycle.Manager, staticFiles map[string]string) []StaticFileConflict {
	conflicts := make([]StaticFileConflict, 0)

	// Check for common feed output paths
	feedPaths := []struct {
		path   string
		source string
	}{
		{"/rss.xml", "feed:rss"},
		{"/atom.xml", "feed:atom"},
		{"/feed.xml", "feed:rss"},
		{"/index.json", "feed:json"},
		{"/feed.json", "feed:json"},
		{"/sitemap.xml", "sitemap"},
	}

	for _, fp := range feedPaths {
		if staticPath, exists := staticFiles[fp.path]; exists {
			conflicts = append(conflicts, StaticFileConflict{
				GeneratedSource: fp.source,
				GeneratedOutput: fp.path,
				StaticFile:      staticPath,
				OutputPath:      fp.path,
			})
		}
	}

	// Check feed configs from cache
	var feedConfigs []models.FeedConfig
	if cached, ok := m.Cache().Get("feed_configs"); ok {
		if fcs, ok := cached.([]models.FeedConfig); ok {
			feedConfigs = fcs
		}
	}

	for i := range feedConfigs {
		fc := &feedConfigs[i]
		feedDir := fc.Slug
		if feedDir == "" {
			feedDir = ""
		}

		// Check each feed format
		if fc.Formats.RSS {
			rssPath := "/" + feedDir
			if feedDir != "" {
				rssPath += "/"
			}
			rssPath += "rss.xml"
			if staticPath, exists := staticFiles[rssPath]; exists {
				conflicts = append(conflicts, StaticFileConflict{
					GeneratedSource: fmt.Sprintf("feed:%s:rss", fc.Slug),
					GeneratedOutput: rssPath,
					StaticFile:      staticPath,
					OutputPath:      rssPath,
				})
			}
		}

		if fc.Formats.Atom {
			atomPath := "/" + feedDir
			if feedDir != "" {
				atomPath += "/"
			}
			atomPath += "atom.xml"
			if staticPath, exists := staticFiles[atomPath]; exists {
				conflicts = append(conflicts, StaticFileConflict{
					GeneratedSource: fmt.Sprintf("feed:%s:atom", fc.Slug),
					GeneratedOutput: atomPath,
					StaticFile:      staticPath,
					OutputPath:      atomPath,
				})
			}
		}

		if fc.Formats.JSON {
			jsonPath := "/" + feedDir
			if feedDir != "" {
				jsonPath += "/"
			}
			jsonPath += "index.json"
			if staticPath, exists := staticFiles[jsonPath]; exists {
				conflicts = append(conflicts, StaticFileConflict{
					GeneratedSource: fmt.Sprintf("feed:%s:json", fc.Slug),
					GeneratedOutput: jsonPath,
					StaticFile:      staticPath,
					OutputPath:      jsonPath,
				})
			}
		}
	}

	return conflicts
}

// reportConflicts returns an error with all detected conflicts formatted as warnings.
func (p *StaticFileConflictsPlugin) reportConflicts() error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("WARNING: %d static file conflict(s) detected\n", len(p.conflicts)))
	sb.WriteString("Static files will override generated content.\n\n")

	for _, c := range p.conflicts {
		sb.WriteString(fmt.Sprintf("  Conflicting %s:\n", c.OutputPath))
		sb.WriteString(fmt.Sprintf("    Generated: %s → %s\n", c.GeneratedSource, c.GeneratedOutput))
		sb.WriteString(fmt.Sprintf("    Static:    %s → %s\n", c.StaticFile, c.OutputPath))
		sb.WriteString("    The static file will override the generated one.\n\n")
	}

	sb.WriteString("To fix:\n")
	sb.WriteString("  - Remove the static file if you want the generated version\n")
	sb.WriteString("  - Remove or rename the source file if you want the static version\n")
	sb.WriteString("  - Disable the static_file_conflicts lint rule in config if intentional\n")

	// Return as a non-critical warning (not an error that stops the build)
	// We use a custom error type that the lifecycle manager can detect as non-critical
	return &StaticFileConflictWarning{Message: sb.String(), Conflicts: p.conflicts}
}

// StaticFileConflictWarning is a warning (not error) about static file conflicts.
type StaticFileConflictWarning struct {
	Message   string
	Conflicts []StaticFileConflict
}

// Error implements the error interface.
func (w *StaticFileConflictWarning) Error() string {
	return w.Message
}

// IsWarning indicates this is a non-critical warning.
func (w *StaticFileConflictWarning) IsWarning() bool {
	return true
}

// Ensure StaticFileConflictsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*StaticFileConflictsPlugin)(nil)
	_ lifecycle.CollectPlugin  = (*StaticFileConflictsPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*StaticFileConflictsPlugin)(nil)
)
