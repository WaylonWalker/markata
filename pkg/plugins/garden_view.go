// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata-go/pkg/buildcache"
	"github.com/WaylonWalker/markata-go/pkg/lifecycle"
	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/templates"
)

// Node type constants for the garden graph.
const (
	gardenNodeTypePost = "post"
	gardenNodeTypeTag  = "tag"
)

// Edge type constants for the garden graph.
const (
	gardenEdgeTypeLink         = "link"
	gardenEdgeTypeTag          = "tag"
	gardenEdgeTypeCooccurrence = "co-occurrence"
)

// GardenNode represents a node in the garden knowledge graph.
type GardenNode struct {
	// ID is the unique identifier for this node (e.g., "post:my-slug" or "tag:go")
	ID string `json:"id"`

	// Type is the node type: "post" or "tag"
	Type string `json:"type"`

	// Label is the display name for this node
	Label string `json:"label"`

	// Href is the URL to the node's page
	Href string `json:"href"`

	// Tags is the list of tags (only for post nodes)
	Tags []string `json:"tags,omitempty"`

	// Date is the publication date (only for post nodes)
	Date string `json:"date,omitempty"`

	// Description is a short description (only for post nodes)
	Description string `json:"description,omitempty"`

	// Count is the number of posts with this tag (only for tag nodes)
	Count int `json:"count,omitempty"`
}

// GardenEdge represents a relationship between two nodes in the garden graph.
type GardenEdge struct {
	// Source is the ID of the source node
	Source string `json:"source"`

	// Target is the ID of the target node
	Target string `json:"target"`

	// Type is the edge type: "link", "tag", or "co-occurrence"
	Type string `json:"type"`

	// Weight is the strength of the relationship (for co-occurrence edges)
	Weight int `json:"weight,omitempty"`
}

// GardenGraph is the complete graph data structure exported as JSON.
type GardenGraph struct {
	// Nodes is the list of all nodes in the graph
	Nodes []GardenNode `json:"nodes"`

	// Edges is the list of all edges in the graph
	Edges []GardenEdge `json:"edges"`
}

// TagCluster represents a tag and its related tags for the garden template.
type TagCluster struct {
	// Name is the tag name
	Name string

	// Count is the number of posts with this tag
	Count int

	// Href is the URL to the tag page
	Href string

	// Related is the list of related tags (by co-occurrence)
	Related []TagRelation
}

// TagRelation represents a related tag with co-occurrence weight.
type TagRelation struct {
	// Name is the related tag name
	Name string

	// Count is the number of posts shared with the tag
	Count int

	// Href is the URL to the related tag page
	Href string
}

// GardenViewPlugin generates a knowledge graph JSON file and optional garden page.
// It runs during the Write stage after link_collector and feeds have finished.
type GardenViewPlugin struct {
	engineMu    sync.RWMutex
	engineCache map[string]*templates.Engine
}

// NewGardenViewPlugin creates a new GardenViewPlugin.
func NewGardenViewPlugin() *GardenViewPlugin {
	return &GardenViewPlugin{
		engineCache: make(map[string]*templates.Engine),
	}
}

// Name returns the unique name of the plugin.
func (p *GardenViewPlugin) Name() string {
	return "garden_view"
}

// Priority returns the plugin's priority for a given stage.
func (p *GardenViewPlugin) Priority(stage lifecycle.Stage) int {
	switch stage {
	case lifecycle.StageWrite:
		// Run late, after link_collector and feeds
		return lifecycle.PriorityLate
	default:
		return lifecycle.PriorityDefault
	}
}

// Write generates the garden graph JSON and optional HTML page.
func (p *GardenViewPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	gardenConfig := p.getGardenConfig(config)
	cache := GetBuildCache(m)

	if !gardenConfig.IsEnabled() {
		return nil
	}

	if cache != nil {
		// Content-only dirty flags do not capture deleted posts until cache
		// cleanup, which runs after this Write stage. The semantic fingerprint
		// is cheap to calculate and includes the current membership, so compare
		// it even when the per-post change flag is clear.
		gardenHash := computeGardenHash(m.Posts(), &gardenConfig, cache)
		if cached := cache.GetGardenHash(); cached != "" && cached == gardenHash {
			return nil
		}
		cache.SetGardenHash(gardenHash)
	}

	posts := m.Posts()
	if len(posts) == 0 {
		log.Printf("[garden_view] No posts found, skipping garden view")
		return nil
	}

	// Filter posts for the graph
	filteredPosts := p.filterPosts(posts, &gardenConfig)
	if len(filteredPosts) == 0 {
		log.Printf("[garden_view] No visible posts after filtering, skipping garden view")
		return nil
	}

	// Build the graph
	graph := p.buildGraph(filteredPosts, &gardenConfig)

	// Apply node limit
	p.applyNodeLimit(&graph, &gardenConfig)

	// Sort for deterministic output
	p.sortGraph(&graph)

	// Create output directory
	outputDir := config.OutputDir
	gardenDir := filepath.Join(outputDir, gardenConfig.GetPath())
	if err := os.MkdirAll(gardenDir, 0o755); err != nil {
		return fmt.Errorf("creating garden directory: %w", err)
	}

	// Export graph.json
	if gardenConfig.IsExportJSON() {
		if err := p.exportGraphJSON(gardenDir, &graph); err != nil {
			return err
		}
	}

	// Render garden page
	if gardenConfig.IsRenderPage() {
		if err := p.renderGardenPage(config, &gardenConfig, &graph); err != nil {
			return err
		}
	}

	log.Printf("[garden_view] Generated /%s/ with %d nodes and %d edges",
		gardenConfig.GetPath(), len(graph.Nodes), len(graph.Edges))

	return nil
}

// filterPosts returns posts that should be included in the garden graph.
func (p *GardenViewPlugin) filterPosts(posts []*models.Post, config *models.GardenConfig) []*models.Post {
	filtered := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if post.Draft || !post.Published || post.Private || post.Skip {
			continue
		}

		// Check if any of the post's tags are excluded
		excluded := false
		for _, tag := range post.Tags {
			if config.IsTagExcluded(tag) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		filtered = append(filtered, post)
	}
	return filtered
}

func computeGardenHash(posts []*models.Post, config *models.GardenConfig, cache *buildcache.Cache) string {
	if config == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(config.GetPath())
	b.WriteByte('\x00')
	b.WriteString(config.Title)
	b.WriteByte('\x00')
	b.WriteString(config.Description)
	b.WriteByte('\x00')
	b.WriteString(config.GetTemplate())
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", config.IsEnabled()))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", config.IsExportJSON()))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", config.IsRenderPage()))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", config.IsIncludeTags()))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%t", config.IsIncludePosts()))
	b.WriteByte('\x00')
	b.WriteString(fmt.Sprintf("%d", config.GetMaxNodes()))
	b.WriteByte('\x00')
	excludeTags := append([]string(nil), config.ExcludeTags...)
	sort.Strings(excludeTags)
	for _, tag := range excludeTags {
		b.WriteString("exclude:")
		b.WriteString(tag)
		b.WriteByte('\x00')
	}
	for _, post := range posts {
		if post.Skip || post.Draft || !post.Published || post.Private {
			continue
		}
		if cache != nil {
			if _, _, gardenHash := cache.GetPostSemanticHashes(post.Path); gardenHash != "" {
				b.WriteString(gardenHash)
				b.WriteByte('\x00')
				continue
			}
		}
		b.WriteString(computePostGardenHash(post))
		b.WriteByte('\x00')
	}
	return buildcache.ContentHash(b.String())
}

// buildGraph constructs the full garden graph from filtered posts.
func (p *GardenViewPlugin) buildGraph(posts []*models.Post, config *models.GardenConfig) GardenGraph {
	graph := GardenGraph{
		Nodes: []GardenNode{},
		Edges: []GardenEdge{},
	}

	// Track tag counts for tag nodes
	tagCounts := make(map[string]int)
	// Track tag co-occurrences for co-occurrence edges
	tagCooccurrence := make(map[string]int) // "tag1:tag2" -> count

	// Build post nodes and collect tag info
	if config.IsIncludePosts() {
		for _, post := range posts {
			node := p.postToNode(post)
			graph.Nodes = append(graph.Nodes, node)

			// Count tags
			for _, tag := range post.Tags {
				if !config.IsTagExcluded(tag) {
					tagCounts[tag]++
				}
			}

			// Count co-occurrences
			visibleTags := p.filterTags(post.Tags, config)
			for i := range visibleTags {
				for j := i + 1; j < len(visibleTags); j++ {
					key := p.cooccurrenceKey(visibleTags[i], visibleTags[j])
					tagCooccurrence[key]++
				}
			}

			// Add post→post edges from internal links
			for _, outlink := range post.Outlinks {
				if outlink.IsInternal && outlink.TargetPost != nil && !outlink.IsSelf {
					// Only include edges to posts that are in our filtered set
					targetSlug := outlink.TargetPost.Slug
					edge := GardenEdge{
						Source: "post:" + post.Slug,
						Target: "post:" + targetSlug,
						Type:   gardenEdgeTypeLink,
					}
					graph.Edges = append(graph.Edges, edge)
				}
			}

			// Add post→tag edges
			if config.IsIncludeTags() {
				for _, tag := range visibleTags {
					edge := GardenEdge{
						Source: "post:" + post.Slug,
						Target: "tag:" + tag,
						Type:   gardenEdgeTypeTag,
					}
					graph.Edges = append(graph.Edges, edge)
				}
			}
		}
	}

	// Build tag nodes
	if config.IsIncludeTags() {
		for tag, count := range tagCounts {
			slug := models.Slugify(tag)
			node := GardenNode{
				ID:    "tag:" + tag,
				Type:  gardenNodeTypeTag,
				Label: tag,
				Href:  "/tags/" + slug + "/",
				Count: count,
			}
			graph.Nodes = append(graph.Nodes, node)
		}

		// Add tag↔tag co-occurrence edges
		for key, count := range tagCooccurrence {
			parts := strings.SplitN(key, "\x00", 2)
			if len(parts) == 2 {
				edge := GardenEdge{
					Source: "tag:" + parts[0],
					Target: "tag:" + parts[1],
					Type:   gardenEdgeTypeCooccurrence,
					Weight: count,
				}
				graph.Edges = append(graph.Edges, edge)
			}
		}
	}

	return graph
}

// filterTags returns tags that are not excluded.
func (p *GardenViewPlugin) filterTags(tags []string, config *models.GardenConfig) []string {
	var filtered []string
	for _, tag := range tags {
		if !config.IsTagExcluded(tag) {
			filtered = append(filtered, tag)
		}
	}
	return filtered
}

// cooccurrenceKey returns a deterministic key for a pair of tags.
// Tags are sorted alphabetically to ensure the key is the same regardless of order.
func (p *GardenViewPlugin) cooccurrenceKey(tag1, tag2 string) string {
	if tag1 > tag2 {
		tag1, tag2 = tag2, tag1
	}
	return tag1 + "\x00" + tag2
}

// postToNode converts a post to a GardenNode.
func (p *GardenViewPlugin) postToNode(post *models.Post) GardenNode {
	node := GardenNode{
		ID:   "post:" + post.Slug,
		Type: gardenNodeTypePost,
		Href: post.Href,
		Tags: post.Tags,
	}

	node.Label = post.PlainTitle()
	if node.Label == "" {
		node.Label = post.Slug
	}

	if post.Date != nil {
		node.Date = post.Date.Format("2006-01-02T15:04:05Z07:00")
	}

	if post.Description != nil {
		node.Description = *post.Description
	}

	return node
}

// applyNodeLimit removes excess tag nodes if the graph exceeds max_nodes.
func (p *GardenViewPlugin) applyNodeLimit(graph *GardenGraph, config *models.GardenConfig) {
	maxNodes := config.GetMaxNodes()
	if len(graph.Nodes) <= maxNodes {
		return
	}

	// Separate post and tag nodes
	var postNodes, tagNodes []GardenNode
	for i := range graph.Nodes {
		if graph.Nodes[i].Type == gardenNodeTypePost {
			postNodes = append(postNodes, graph.Nodes[i])
		} else {
			tagNodes = append(tagNodes, graph.Nodes[i])
		}
	}

	// Keep highest-signal tags first when trimming.
	sort.Slice(tagNodes, func(i, j int) bool {
		return tagNodes[i].Count > tagNodes[j].Count
	})

	// If posts alone exceed the limit, reserve space for tags when available.
	if len(postNodes) >= maxNodes {
		if len(tagNodes) == 0 {
			graph.Nodes = postNodes[:maxNodes]
			p.pruneEdges(graph)
			return
		}

		tagSlots := maxNodes / 5
		if tagSlots < 1 {
			tagSlots = 1
		}
		if tagSlots > len(tagNodes) {
			tagSlots = len(tagNodes)
		}
		if tagSlots >= maxNodes {
			tagSlots = maxNodes - 1
		}

		postSlots := maxNodes - tagSlots
		if postSlots < 0 {
			postSlots = 0
		}

		graph.Nodes = make([]GardenNode, 0, postSlots+tagSlots)
		graph.Nodes = append(graph.Nodes, postNodes[:postSlots]...)
		graph.Nodes = append(graph.Nodes, tagNodes[:tagSlots]...)
		// Remove edges referencing removed nodes
		p.pruneEdges(graph)
		return
	}

	remaining := maxNodes - len(postNodes)
	if remaining > len(tagNodes) {
		remaining = len(tagNodes)
	}
	tagNodes = tagNodes[:remaining]

	graph.Nodes = make([]GardenNode, 0, len(postNodes)+len(tagNodes))
	graph.Nodes = append(graph.Nodes, postNodes...)
	graph.Nodes = append(graph.Nodes, tagNodes...)

	// Remove edges referencing removed nodes
	p.pruneEdges(graph)
}

// pruneEdges removes edges that reference nodes not in the graph.
func (p *GardenViewPlugin) pruneEdges(graph *GardenGraph) {
	nodeSet := make(map[string]bool, len(graph.Nodes))
	for i := range graph.Nodes {
		nodeSet[graph.Nodes[i].ID] = true
	}

	validEdges := make([]GardenEdge, 0, len(graph.Edges))
	for i := range graph.Edges {
		if nodeSet[graph.Edges[i].Source] && nodeSet[graph.Edges[i].Target] {
			validEdges = append(validEdges, graph.Edges[i])
		}
	}
	graph.Edges = validEdges
}

// sortGraph sorts nodes and edges for deterministic output.
func (p *GardenViewPlugin) sortGraph(graph *GardenGraph) {
	sort.Slice(graph.Nodes, func(i, j int) bool {
		return graph.Nodes[i].ID < graph.Nodes[j].ID
	})

	sort.Slice(graph.Edges, func(i, j int) bool {
		if graph.Edges[i].Source != graph.Edges[j].Source {
			return graph.Edges[i].Source < graph.Edges[j].Source
		}
		if graph.Edges[i].Target != graph.Edges[j].Target {
			return graph.Edges[i].Target < graph.Edges[j].Target
		}
		return graph.Edges[i].Type < graph.Edges[j].Type
	})
}

// exportGraphJSON writes the graph data as JSON to the output directory.
func (p *GardenViewPlugin) exportGraphJSON(gardenDir string, graph *GardenGraph) error {
	data, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling garden graph: %w", err)
	}

	outputPath := filepath.Join(gardenDir, "graph.json")
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return fmt.Errorf("writing garden graph.json: %w", err)
	}

	return nil
}

// renderGardenPage renders the garden HTML page.
func (p *GardenViewPlugin) renderGardenPage(config *lifecycle.Config, gardenConfig *models.GardenConfig, graph *GardenGraph) error {
	// Create output directory
	outputDir := config.OutputDir
	gardenDir := filepath.Join(outputDir, gardenConfig.GetPath())
	if err := os.MkdirAll(gardenDir, 0o755); err != nil {
		return fmt.Errorf("creating garden directory: %w", err)
	}

	// Get template engine
	engine, err := p.createTemplateEngine(config)
	if err != nil {
		return err
	}

	templateName := gardenConfig.GetTemplate()
	if !engine.TemplateExists(templateName) {
		log.Printf("[garden_view] Warning: template %q not found, skipping garden page", templateName)
		return nil
	}

	// Build tag clusters for the template
	tagClusters := p.buildTagClusters(graph)

	// Count posts and tags in the graph
	totalPosts := 0
	totalTags := 0
	for i := range graph.Nodes {
		if graph.Nodes[i].Type == gardenNodeTypePost {
			totalPosts++
		} else if graph.Nodes[i].Type == gardenNodeTypeTag {
			totalTags++
		}
	}

	// Build context
	modelsConfig := ToModelsConfig(config)
	title := gardenConfig.Title
	description := gardenConfig.Description
	syntheticPost := &models.Post{
		Slug:        gardenConfig.GetPath(),
		Title:       &title,
		Description: &description,
	}

	ctx := templates.NewContext(syntheticPost, "", modelsConfig)
	ctx.Extra["graph_json"] = "/" + gardenConfig.GetPath() + "/graph.json"
	ctx.Extra["tag_clusters"] = tagClusters
	ctx.Extra["total_posts"] = totalPosts
	ctx.Extra["total_tags"] = totalTags
	ctx.Extra["total_edges"] = len(graph.Edges)

	// Render template
	html, err := engine.Render(templateName, ctx)
	if err != nil {
		return fmt.Errorf("rendering garden template: %w", err)
	}

	// Write output file
	outputPath := filepath.Join(gardenDir, "index.html")
	//nolint:gosec // G306: Output files need 0644 for web serving
	if err := os.WriteFile(outputPath, []byte(html), 0o644); err != nil {
		return fmt.Errorf("writing garden page: %w", err)
	}

	return nil
}

// buildTagClusters creates tag clusters with related tags from co-occurrence data.
func (p *GardenViewPlugin) buildTagClusters(graph *GardenGraph) []TagCluster {
	// Build tag info map
	tagInfo := make(map[string]*TagCluster)
	tagHref := make(map[string]string)
	for i := range graph.Nodes {
		if graph.Nodes[i].Type == gardenNodeTypeTag {
			tagInfo[graph.Nodes[i].ID] = &TagCluster{
				Name:    graph.Nodes[i].Label,
				Count:   graph.Nodes[i].Count,
				Href:    graph.Nodes[i].Href,
				Related: []TagRelation{},
			}
			tagHref[graph.Nodes[i].ID] = graph.Nodes[i].Href
		}
	}

	// Add related tags from co-occurrence edges
	for i := range graph.Edges {
		if graph.Edges[i].Type == gardenEdgeTypeCooccurrence {
			if cluster, ok := tagInfo[graph.Edges[i].Source]; ok {
				// Extract tag name from ID (format: "tag:name")
				targetTag := strings.TrimPrefix(graph.Edges[i].Target, "tag:")
				cluster.Related = append(cluster.Related, TagRelation{
					Name:  targetTag,
					Count: graph.Edges[i].Weight,
					Href:  tagHref[graph.Edges[i].Target],
				})
			}
			if cluster, ok := tagInfo[graph.Edges[i].Target]; ok {
				sourceTag := strings.TrimPrefix(graph.Edges[i].Source, "tag:")
				cluster.Related = append(cluster.Related, TagRelation{
					Name:  sourceTag,
					Count: graph.Edges[i].Weight,
					Href:  tagHref[graph.Edges[i].Source],
				})
			}
		}
	}

	// Sort related tags by weight (desc), then name (asc)
	for _, cluster := range tagInfo {
		sort.Slice(cluster.Related, func(i, j int) bool {
			if cluster.Related[i].Count != cluster.Related[j].Count {
				return cluster.Related[i].Count > cluster.Related[j].Count
			}
			return cluster.Related[i].Name < cluster.Related[j].Name
		})
	}

	// Convert to sorted slice
	clusters := make([]TagCluster, 0, len(tagInfo))
	for _, cluster := range tagInfo {
		clusters = append(clusters, *cluster)
	}

	// Sort clusters by count (descending), then name (ascending)
	sort.Slice(clusters, func(i, j int) bool {
		if clusters[i].Count != clusters[j].Count {
			return clusters[i].Count > clusters[j].Count
		}
		return clusters[i].Name < clusters[j].Name
	})

	return clusters
}

// getOrCreateEngine returns a cached template engine, or creates one if not cached.
func (p *GardenViewPlugin) getOrCreateEngine(templatesDir, themeName string) (*templates.Engine, error) {
	cacheKey := templatesDir + ":" + themeName

	// Fast path: check cache with read lock
	p.engineMu.RLock()
	if engine, ok := p.engineCache[cacheKey]; ok {
		p.engineMu.RUnlock()
		return engine, nil
	}
	p.engineMu.RUnlock()

	// Slow path: create engine with write lock
	p.engineMu.Lock()
	defer p.engineMu.Unlock()

	// Double-check after acquiring write lock
	if engine, ok := p.engineCache[cacheKey]; ok {
		return engine, nil
	}

	engine, err := templates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return nil, err
	}

	p.engineCache[cacheKey] = engine
	return engine, nil
}

// createTemplateEngine creates or retrieves a cached template engine.
func (p *GardenViewPlugin) createTemplateEngine(config *lifecycle.Config) (*templates.Engine, error) {
	templatesDir := PluginNameTemplates
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}

	themeName := getThemeName(config)

	return p.getOrCreateEngine(templatesDir, themeName)
}

// getGardenConfig retrieves garden configuration from the manager config.
func (p *GardenViewPlugin) getGardenConfig(config *lifecycle.Config) models.GardenConfig {
	// Prefer direct access to the full models.Config stored in Extra
	if modelsConfig, ok := config.Extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		return modelsConfig.Garden
	}
	// Fall back to ToModelsConfig reconstruction
	modelsConfig := ToModelsConfig(config)
	if modelsConfig != nil {
		return modelsConfig.Garden
	}
	return models.NewGardenConfig()
}

// Ensure GardenViewPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin         = (*GardenViewPlugin)(nil)
	_ lifecycle.WritePlugin    = (*GardenViewPlugin)(nil)
	_ lifecycle.PriorityPlugin = (*GardenViewPlugin)(nil)
)
