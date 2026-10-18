// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// AutoTitlePlugin auto-generates human-readable titles for posts that don't have one.
// It uses a comprehensive fallback strategy to ensure every post has a title:
//  1. Frontmatter title (explicit, highest priority)
//  2. First H1 heading from markdown content
//  3. Filename-based title (with date prefix stripping)
//  4. Directory name (for index.md files)
//  5. Generated fallback with timestamp
type AutoTitlePlugin struct{}

// dateRegex matches common date prefixes in filenames: YYYY-MM-DD with optional separator
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[-_]?`)

// h1Regex matches Markdown H1 headings at the start of a line
var h1Regex = regexp.MustCompile(`^#\s+(.+)$`)

// NewAutoTitlePlugin creates a new AutoTitlePlugin.
func NewAutoTitlePlugin() *AutoTitlePlugin {
	return &AutoTitlePlugin{}
}

// Name returns the unique name of the plugin.
func (p *AutoTitlePlugin) Name() string {
	return "auto_title"
}

// Priority returns the plugin priority for the given stage.
// Auto title should run very early in transform to have title available for other plugins.
func (p *AutoTitlePlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityFirst
	}
	return lifecycle.PriorityDefault
}

// Transform generates titles for posts that don't have one.
// Uses a comprehensive fallback strategy to ensure no post has a nil title.
func (p *AutoTitlePlugin) Transform(m *lifecycle.Manager) error {
	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip {
			return false
		}
		if post.Title != nil && *post.Title != "" {
			return false
		}
		return true
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		title := p.inferTitle(post)
		post.Title = &title
		return nil
	})
}

// inferTitle attempts to generate a title using multiple fallback strategies.
// Priority order:
//  1. First H1 heading from content
//  2. Date-stripped filename
//  3. Directory name (for index.md files)
//  4. Plain filename-based title
//  5. Generated fallback with path identifier
func (p *AutoTitlePlugin) inferTitle(post *models.Post) string {
	// Strategy 1: Extract from first H1 heading in content
	if title := p.extractFromContent(post.Content); title != "" {
		return title
	}

	// Strategy 2: Check for index files and use directory name
	base := filepath.Base(post.Path)
	if strings.EqualFold(base, "index.md") || strings.EqualFold(base, "index.markdown") {
		if title := p.extractFromDirectory(post.Path); title != "" {
			return title
		}
	}

	// Strategy 3: Generate from filename (with date stripping)
	if title := p.generateTitle(post.Path); title != "" {
		return title
	}

	// Strategy 4: Final fallback - never return empty
	return p.generateFallback(post.Path)
}

// extractFromContent extracts a title from the first H1 heading in markdown content.
func (p *AutoTitlePlugin) extractFromContent(content string) string {
	if content == "" {
		return ""
	}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if matches := h1Regex.FindStringSubmatch(trimmed); len(matches) > 1 {
			// Clean up the title - remove any trailing markup
			title := strings.TrimSpace(matches[1])
			// Remove trailing # characters (alternate H1 syntax: # Title #)
			title = strings.TrimRight(title, "#")
			title = strings.TrimSpace(title)
			if title != "" {
				return title
			}
		}
	}
	return ""
}

// extractFromDirectory extracts a title from the parent directory name.
// Used for index.md files where the directory name is more meaningful.
func (p *AutoTitlePlugin) extractFromDirectory(path string) string {
	dir := filepath.Dir(path)
	if dir == "" || dir == "." || dir == "/" {
		return ""
	}

	// Get the immediate parent directory name
	dirName := filepath.Base(dir)
	if dirName == "" || dirName == "." || dirName == "/" {
		return ""
	}

	// Convert directory name to title case
	title := strings.ReplaceAll(dirName, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")
	return toTitleCase(title)
}

// stripDatePrefix removes a date prefix (YYYY-MM-DD) from a filename.
// Returns the cleaned filename and whether a date was found.
func (p *AutoTitlePlugin) stripDatePrefix(filename string) (string, bool) {
	if match := dateRegex.FindString(filename); match != "" {
		return strings.TrimLeft(filename[len(match):], "-_"), true
	}
	return filename, false
}

// generateFallback creates a fallback title when all other strategies fail.
// Uses a combination of path info and timestamp to ensure uniqueness.
func (p *AutoTitlePlugin) generateFallback(path string) string {
	if path != "" {
		// Try to create something from the path
		base := filepath.Base(path)
		stem := strings.TrimSuffix(base, filepath.Ext(base))
		if stem != "" && stem != "." {
			return "Untitled " + stem
		}
	}
	// Last resort - use timestamp
	return "Untitled " + time.Now().Format("2006-01-02-150405")
}

// generateTitle creates a human-readable title from a file path.
// Includes date prefix stripping for files like "2024-01-15-my-post.md".
func (p *AutoTitlePlugin) generateTitle(path string) string {
	// Extract filename without extension
	base := filepath.Base(path)
	stem := strings.TrimSuffix(base, filepath.Ext(base))

	if stem == "" {
		return ""
	}

	// Strip date prefix if present
	stem, _ = p.stripDatePrefix(stem)
	if stem == "" {
		return ""
	}

	// Replace hyphens and underscores with spaces
	title := strings.ReplaceAll(stem, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Apply title case
	title = toTitleCase(title)

	return title
}

// toTitleCase converts a string to title case.
// Each word starts with an uppercase letter, rest are lowercase.
func toTitleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if word == "" {
			continue
		}
		// Capitalize first rune, lowercase the rest
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// Ensure AutoTitlePlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*AutoTitlePlugin)(nil)
	_ lifecycle.TransformPlugin = (*AutoTitlePlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*AutoTitlePlugin)(nil)
)
