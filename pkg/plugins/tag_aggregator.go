package plugins

import (
	"log"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

// TagAggregatorPlugin normalizes tags using synonyms and expands them with hierarchical relationships.
//
// Features:
// - Synonym normalization: Replace variant tag names with canonical tags (e.g., "k8s" -> "kubernetes")
// - Hierarchical expansion: Automatically add parent/related tags (e.g., "pandas" adds "data" and "python")
// - Recursive expansion: Additional tags are applied recursively to build complete tag hierarchies
//
// This plugin runs in the Load stage with priority 50 (after posts are loaded but before AutoFeedsPlugin
// at PriorityLate=100) so that expanded tags are visible to auto-generated tag feeds.
type TagAggregatorPlugin struct {
	manager *lifecycle.Manager
}

var (
	_ lifecycle.Plugin          = (*TagAggregatorPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*TagAggregatorPlugin)(nil)
	_ lifecycle.LoadPlugin      = (*TagAggregatorPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*TagAggregatorPlugin)(nil)
)

// NewTagAggregatorPlugin creates a new tag aggregator plugin.
func NewTagAggregatorPlugin() *TagAggregatorPlugin {
	return &TagAggregatorPlugin{}
}

// Name returns the plugin name.
func (p *TagAggregatorPlugin) Name() string {
	return "tag_aggregator"
}

// Priority returns the plugin's priority for a given stage.
// We run in Load stage with priority 50 - after posts are loaded (default=0)
// but before AutoFeedsPlugin (PriorityLate=100) creates tag feeds.
func (p *TagAggregatorPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageLoad {
		return 50 // Between default (0) and late (100)
	}
	return lifecycle.PriorityDefault
}

// Configure stores the manager reference.
func (p *TagAggregatorPlugin) Configure(m *lifecycle.Manager) error {
	p.manager = m
	return nil
}

// Load normalizes and expands tags for all posts.
// This runs in the Load stage so expanded tags are visible to auto-generated feeds.
func (p *TagAggregatorPlugin) Load(m *lifecycle.Manager) error {
	// Get models.Config from lifecycle.Config
	modelsConfig, ok := getModelsConfig(m.Config())
	if !ok {
		// Gracefully skip if models.Config is not available
		// This can happen in tests or minimal configurations
		return nil
	}

	cfg := &modelsConfig.TagAggregator

	// Skip if disabled
	if !cfg.IsEnabled() {
		return nil
	}

	// Skip if no synonyms or additional tags configured
	if len(cfg.Synonyms) == 0 && len(cfg.Additional) == 0 {
		return nil
	}

	posts := m.Posts()
	synonymCount := 0
	addedCount := 0

	for _, post := range posts {
		if len(post.Tags) == 0 {
			continue
		}

		originalTags := post.Tags

		// Step 1: Normalize tags using synonyms
		normalizedTags := normalizeTags(post.Tags, cfg.Synonyms)
		if len(normalizedTags) != len(originalTags) || !tagsEqual(normalizedTags, originalTags) {
			synonymCount++
		}

		// Step 2: Recursively expand with additional tags
		expandedTags := expandTags(normalizedTags, cfg.Additional)
		addedTagCount := len(expandedTags) - len(normalizedTags)
		if addedTagCount > 0 {
			addedCount += addedTagCount
		}

		// Update post tags (sorted)
		post.Tags = sortTags(expandedTags)
	}

	if synonymCount > 0 || addedCount > 0 {
		log.Printf("[tag_aggregator] Processed tags: %d posts with synonym normalization, %d tags added via expansion",
			synonymCount, addedCount)
	}

	return nil
}

// normalizeTags replaces synonym tags with their canonical versions.
func normalizeTags(tags []string, synonyms map[string][]string) []string {
	normalized := make([]string, 0, len(tags))
	seen := make(map[string]bool)

	for _, tag := range tags {
		canonicalTag := tag

		// Check if this tag is a synonym
		for canonical, variants := range synonyms {
			for _, variant := range variants {
				if strings.EqualFold(tag, variant) {
					canonicalTag = canonical
					break
				}
			}
			if canonicalTag != tag {
				break
			}
		}

		// Add canonical tag if not already seen (de-duplicate)
		if !seen[canonicalTag] {
			normalized = append(normalized, canonicalTag)
			seen[canonicalTag] = true
		}
	}

	return normalized
}

// expandTags recursively adds additional tags based on the configured relationships.
func expandTags(tags []string, additional map[string][]string) []string {
	result := make(map[string]bool)

	// Add all initial tags
	for _, tag := range tags {
		result[tag] = true
	}

	// Process tags to expand (use a queue to avoid infinite loops)
	queue := make([]string, len(tags))
	copy(queue, tags)
	processed := make(map[string]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if processed[current] {
			continue
		}
		processed[current] = true

		// Add additional tags for this tag
		if additionalTags, ok := additional[current]; ok {
			for _, additionalTag := range additionalTags {
				if !result[additionalTag] {
					result[additionalTag] = true
					queue = append(queue, additionalTag)
				}
			}
		}
	}

	// Convert map to slice
	expanded := make([]string, 0, len(result))
	for tag := range result {
		expanded = append(expanded, tag)
	}

	return expanded
}

// sortTags returns a sorted copy of tags.
func sortTags(tags []string) []string {
	sorted := make([]string, len(tags))
	copy(sorted, tags)
	sort.Strings(sorted)
	return sorted
}

// tagsEqual checks if two tag slices contain the same tags (order-independent).
func tagsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]bool)
	for _, tag := range a {
		aMap[tag] = true
	}

	for _, tag := range b {
		if !aMap[tag] {
			return false
		}
	}

	return true
}
