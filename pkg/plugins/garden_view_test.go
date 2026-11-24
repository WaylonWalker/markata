package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func strPtrG(s string) *string { return &s }

func timePtrG(t time.Time) *time.Time { return &t }

func boolPtrG(b bool) *bool { return &b }

// newTestGardenPlugin creates a GardenViewPlugin for testing.
func newTestGardenPlugin() *GardenViewPlugin {
	return NewGardenViewPlugin()
}

// newTestGardenConfig creates a GardenConfig with defaults for testing.
func newTestGardenConfig() models.GardenConfig {
	return models.NewGardenConfig()
}

// newTestPosts creates a set of test posts for garden view testing.
func newTestPosts() []*models.Post {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	postB := &models.Post{
		Path:      "post-b.md",
		Slug:      "post-b",
		Href:      "/post-b/",
		Title:     strPtrG("Post B"),
		Tags:      []string{"go", "programming"},
		Date:      timePtrG(date),
		Published: true,
	}
	postA := &models.Post{
		Path:      "post-a.md",
		Slug:      "post-a",
		Href:      "/post-a/",
		Title:     strPtrG("Post A"),
		Tags:      []string{"go", "tutorial"},
		Date:      timePtrG(date.AddDate(0, 0, 1)),
		Published: true,
		Outlinks: []*models.Link{
			{
				IsInternal: true,
				TargetPost: postB,
				IsSelf:     false,
			},
		},
	}
	postC := &models.Post{
		Path:        "post-c.md",
		Slug:        "post-c",
		Href:        "/post-c/",
		Title:       strPtrG("Post C"),
		Tags:        []string{"python", "tutorial"},
		Date:        timePtrG(date.AddDate(0, 0, 2)),
		Published:   true,
		Description: strPtrG("A Python tutorial"),
	}
	return []*models.Post{postA, postB, postC}
}

func TestGardenViewPlugin_Name(t *testing.T) {
	p := newTestGardenPlugin()
	if got := p.Name(); got != "garden_view" {
		t.Errorf("Name() = %q, want %q", got, "garden_view")
	}
}

func TestGardenViewPlugin_Priority(t *testing.T) {
	p := newTestGardenPlugin()
	if got := p.Priority(lifecycle.StageWrite); got != lifecycle.PriorityLate {
		t.Errorf("Priority(StageWrite) = %d, want %d", got, lifecycle.PriorityLate)
	}
	if got := p.Priority(lifecycle.StageRender); got != lifecycle.PriorityDefault {
		t.Errorf("Priority(StageRender) = %d, want %d", got, lifecycle.PriorityDefault)
	}
}

func TestGardenViewPlugin_ImplementsInterfaces(_ *testing.T) {
	var _ lifecycle.Plugin = (*GardenViewPlugin)(nil)
	var _ lifecycle.WritePlugin = (*GardenViewPlugin)(nil)
	var _ lifecycle.PriorityPlugin = (*GardenViewPlugin)(nil)
}

func TestGardenViewPlugin_FilterPosts(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		posts     []*models.Post
		config    models.GardenConfig
		wantLen   int
		wantSlugs []string
	}{
		{
			name:    "visible posts pass through",
			posts:   newTestPosts(),
			config:  config,
			wantLen: 3,
		},
		{
			name: "draft posts filtered out",
			posts: []*models.Post{
				{Slug: "visible", Published: true, Date: timePtrG(date)},
				{Slug: "draft", Draft: true, Published: true, Date: timePtrG(date)},
			},
			config:    config,
			wantLen:   1,
			wantSlugs: []string{"visible"},
		},
		{
			name: "unpublished posts filtered out",
			posts: []*models.Post{
				{Slug: "visible", Published: true, Date: timePtrG(date)},
				{Slug: "unpub", Published: false, Date: timePtrG(date)},
			},
			config:    config,
			wantLen:   1,
			wantSlugs: []string{"visible"},
		},
		{
			name: "private posts filtered out",
			posts: []*models.Post{
				{Slug: "visible", Published: true, Date: timePtrG(date)},
				{Slug: "private", Published: true, Private: true, Date: timePtrG(date)},
			},
			config:    config,
			wantLen:   1,
			wantSlugs: []string{"visible"},
		},
		{
			name: "skip posts filtered out",
			posts: []*models.Post{
				{Slug: "visible", Published: true, Date: timePtrG(date)},
				{Slug: "skipped", Published: true, Skip: true, Date: timePtrG(date)},
			},
			config:    config,
			wantLen:   1,
			wantSlugs: []string{"visible"},
		},
		{
			name: "excluded tags filter posts",
			posts: []*models.Post{
				{Slug: "visible", Published: true, Tags: []string{"go"}, Date: timePtrG(date)},
				{Slug: "excluded", Published: true, Tags: []string{"secret"}, Date: timePtrG(date)},
			},
			config: func() models.GardenConfig {
				c := newTestGardenConfig()
				c.ExcludeTags = []string{"secret"}
				return c
			}(),
			wantLen:   1,
			wantSlugs: []string{"visible"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.filterPosts(tt.posts, &tt.config)
			if len(got) != tt.wantLen {
				t.Errorf("filterPosts() returned %d posts, want %d", len(got), tt.wantLen)
			}
			if tt.wantSlugs != nil {
				for i, slug := range tt.wantSlugs {
					if i < len(got) && got[i].Slug != slug {
						t.Errorf("filterPosts()[%d].Slug = %q, want %q", i, got[i].Slug, slug)
					}
				}
			}
		})
	}
}

func TestGardenViewPlugin_BuildGraph_PostNodes(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// Count post nodes
	postNodes := 0
	for _, node := range graph.Nodes {
		if node.Type == "post" {
			postNodes++
		}
	}
	if postNodes != 3 {
		t.Errorf("expected 3 post nodes, got %d", postNodes)
	}
}

func TestGardenViewPlugin_BuildGraph_TagNodes(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// Count tag nodes
	tagNodeMap := make(map[string]GardenNode)
	for _, node := range graph.Nodes {
		if node.Type == "tag" {
			tagNodeMap[node.Label] = node
		}
	}

	// We expect: go (2 posts), programming (1), tutorial (2), python (1)
	expectedTags := map[string]int{
		"go":          2,
		"programming": 1,
		"tutorial":    2,
		"python":      1,
	}
	if len(tagNodeMap) != len(expectedTags) {
		t.Errorf("expected %d tag nodes, got %d", len(expectedTags), len(tagNodeMap))
	}
	for tag, wantCount := range expectedTags {
		node, ok := tagNodeMap[tag]
		if !ok {
			t.Errorf("missing tag node %q", tag)
			continue
		}
		if node.Count != wantCount {
			t.Errorf("tag %q count = %d, want %d", tag, node.Count, wantCount)
		}
	}
}

func TestGardenViewPlugin_BuildGraph_LinkEdges(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// Check for post-to-post link edge (post-a -> post-b)
	found := false
	for _, edge := range graph.Edges {
		if edge.Source == "post:post-a" && edge.Target == "post:post-b" && edge.Type == "link" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected link edge from post:post-a to post:post-b")
	}
}

func TestGardenViewPlugin_BuildGraph_TagEdges(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// Check for post→tag edges
	postTagEdges := 0
	for _, edge := range graph.Edges {
		if edge.Type == "tag" {
			postTagEdges++
		}
	}
	// post-a has 2 tags, post-b has 2 tags, post-c has 2 tags = 6 post→tag edges
	if postTagEdges != 6 {
		t.Errorf("expected 6 post→tag edges, got %d", postTagEdges)
	}
}

func TestGardenViewPlugin_BuildGraph_CooccurrenceEdges(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// Check co-occurrence edges exist
	coEdges := make(map[string]int)
	for _, edge := range graph.Edges {
		if edge.Type == "co-occurrence" {
			key := edge.Source + " -> " + edge.Target
			coEdges[key] = edge.Weight
		}
	}

	// Expected co-occurrences:
	// go + programming (post-b) -> weight 1
	// go + tutorial (post-a) -> weight 1
	// python + tutorial (post-c) -> weight 1
	if len(coEdges) != 3 {
		t.Errorf("expected 3 co-occurrence edges, got %d: %v", len(coEdges), coEdges)
	}
}

func TestGardenViewPlugin_BuildGraph_NoTags(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	config.IncludeTags = boolPtrG(false)
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// Should have only post nodes, no tag nodes
	for _, node := range graph.Nodes {
		if node.Type == "tag" {
			t.Error("expected no tag nodes when include_tags is false")
			break
		}
	}

	// Should have no tag or co-occurrence edges
	for _, edge := range graph.Edges {
		if edge.Type == "tag" || edge.Type == "co-occurrence" {
			t.Errorf("expected no tag/co-occurrence edges when include_tags is false, got type=%q", edge.Type)
			break
		}
	}
}

func TestGardenViewPlugin_BuildGraph_NoPosts(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	config.IncludePosts = boolPtrG(false)
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// Should have no post nodes and no edges
	if len(graph.Nodes) != 0 {
		t.Errorf("expected 0 nodes when include_posts is false, got %d", len(graph.Nodes))
	}
	if len(graph.Edges) != 0 {
		t.Errorf("expected 0 edges when include_posts is false, got %d", len(graph.Edges))
	}
}

func TestGardenViewPlugin_BuildGraph_ExcludeTags(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	config.ExcludeTags = []string{"tutorial"}
	posts := newTestPosts()

	graph := p.buildGraph(posts, &config)

	// "tutorial" tag should not appear as a node
	for _, node := range graph.Nodes {
		if node.Type == "tag" && node.Label == "tutorial" {
			t.Error("expected 'tutorial' tag to be excluded from nodes")
		}
	}

	// No edges should reference tag:tutorial
	for _, edge := range graph.Edges {
		if edge.Source == "tag:tutorial" || edge.Target == "tag:tutorial" {
			t.Errorf("expected no edges referencing tag:tutorial, got %v", edge)
		}
	}
}

func TestGardenViewPlugin_PostToNode(t *testing.T) {
	p := newTestGardenPlugin()
	date := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	post := &models.Post{
		Slug:        "test-post",
		Href:        "/test-post/",
		Title:       strPtrG("Test Post"),
		Tags:        []string{"go"},
		Date:        timePtrG(date),
		Description: strPtrG("A test"),
	}

	node := p.postToNode(post)

	if node.ID != "post:test-post" {
		t.Errorf("ID = %q, want %q", node.ID, "post:test-post")
	}
	if node.Type != "post" {
		t.Errorf("Type = %q, want %q", node.Type, "post")
	}
	if node.Label != "Test Post" {
		t.Errorf("Label = %q, want %q", node.Label, "Test Post")
	}
	if node.Href != "/test-post/" {
		t.Errorf("Href = %q, want %q", node.Href, "/test-post/")
	}
	if node.Description != "A test" {
		t.Errorf("Description = %q, want %q", node.Description, "A test")
	}
	if node.Date != "2024-06-15T00:00:00Z" {
		t.Errorf("Date = %q, want %q", node.Date, "2024-06-15T00:00:00Z")
	}
}

func TestGardenViewPlugin_PostToNode_NoTitle(t *testing.T) {
	p := newTestGardenPlugin()
	post := &models.Post{
		Slug: "no-title",
		Href: "/no-title/",
	}

	node := p.postToNode(post)
	if node.Label != "no-title" {
		t.Errorf("Label = %q, want slug fallback %q", node.Label, "no-title")
	}
}

func TestGardenViewPlugin_ApplyNodeLimit(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	config.MaxNodes = 5

	graph := GardenGraph{
		Nodes: []GardenNode{
			{ID: "post:a", Type: "post"},
			{ID: "post:b", Type: "post"},
			{ID: "post:c", Type: "post"},
			{ID: "tag:popular", Type: "tag", Count: 10},
			{ID: "tag:medium", Type: "tag", Count: 5},
			{ID: "tag:rare", Type: "tag", Count: 1},
		},
		Edges: []GardenEdge{
			{Source: "post:a", Target: "tag:popular", Type: "tag"},
			{Source: "post:a", Target: "tag:rare", Type: "tag"},
			{Source: "post:b", Target: "tag:medium", Type: "tag"},
			{Source: "tag:popular", Target: "tag:medium", Type: "co-occurrence"},
			{Source: "tag:popular", Target: "tag:rare", Type: "co-occurrence"},
		},
	}

	p.applyNodeLimit(&graph, &config)

	// With limit 5: 3 posts + top 2 tags (popular, medium)
	if len(graph.Nodes) != 5 {
		t.Errorf("expected 5 nodes after limit, got %d", len(graph.Nodes))
	}

	// "rare" tag should be removed
	for _, node := range graph.Nodes {
		if node.ID == "tag:rare" {
			t.Error("expected tag:rare to be removed by node limit")
		}
	}

	// Edges referencing tag:rare should be pruned
	for _, edge := range graph.Edges {
		if edge.Source == "tag:rare" || edge.Target == "tag:rare" {
			t.Error("expected edges referencing tag:rare to be pruned")
		}
	}
}

func TestGardenViewPlugin_ApplyNodeLimit_UnderLimit(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	config.MaxNodes = 100

	graph := GardenGraph{
		Nodes: []GardenNode{
			{ID: "post:a", Type: "post"},
			{ID: "tag:go", Type: "tag", Count: 5},
		},
		Edges: []GardenEdge{
			{Source: "post:a", Target: "tag:go", Type: "tag"},
		},
	}

	p.applyNodeLimit(&graph, &config)

	// Under limit, nothing should change
	if len(graph.Nodes) != 2 {
		t.Errorf("expected 2 nodes (under limit), got %d", len(graph.Nodes))
	}
	if len(graph.Edges) != 1 {
		t.Errorf("expected 1 edge (under limit), got %d", len(graph.Edges))
	}
}

func TestGardenViewPlugin_ApplyNodeLimit_ReservesTagsWhenPostsExceedLimit(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	config.MaxNodes = 5

	graph := GardenGraph{
		Nodes: []GardenNode{
			{ID: "post:a", Type: "post"},
			{ID: "post:b", Type: "post"},
			{ID: "post:c", Type: "post"},
			{ID: "post:d", Type: "post"},
			{ID: "post:e", Type: "post"},
			{ID: "post:f", Type: "post"},
			{ID: "tag:top", Type: "tag", Count: 9},
			{ID: "tag:other", Type: "tag", Count: 3},
		},
		Edges: []GardenEdge{
			{Source: "post:a", Target: "tag:top", Type: "tag"},
			{Source: "post:b", Target: "tag:other", Type: "tag"},
		},
	}

	p.applyNodeLimit(&graph, &config)

	if len(graph.Nodes) != 5 {
		t.Fatalf("expected 5 nodes after limit, got %d", len(graph.Nodes))
	}

	tagCount := 0
	for _, node := range graph.Nodes {
		if node.Type == "tag" {
			tagCount++
		}
	}
	if tagCount == 0 {
		t.Fatalf("expected at least one tag node to be preserved")
	}
}

func TestGardenViewPlugin_SortGraph_Deterministic(t *testing.T) {
	p := newTestGardenPlugin()

	graph := GardenGraph{
		Nodes: []GardenNode{
			{ID: "tag:zebra", Type: "tag"},
			{ID: "post:middle", Type: "post"},
			{ID: "post:alpha", Type: "post"},
			{ID: "tag:apple", Type: "tag"},
		},
		Edges: []GardenEdge{
			{Source: "tag:zebra", Target: "tag:apple", Type: "co-occurrence"},
			{Source: "post:alpha", Target: "tag:apple", Type: "tag"},
			{Source: "post:alpha", Target: "post:middle", Type: "link"},
		},
	}

	p.sortGraph(&graph)

	// Nodes sorted by ID
	expectedNodeIDs := []string{"post:alpha", "post:middle", "tag:apple", "tag:zebra"}
	for i, want := range expectedNodeIDs {
		if graph.Nodes[i].ID != want {
			t.Errorf("Nodes[%d].ID = %q, want %q", i, graph.Nodes[i].ID, want)
		}
	}

	// Edges sorted by (source, target, type)
	if graph.Edges[0].Source != "post:alpha" {
		t.Errorf("Edges[0].Source = %q, want %q", graph.Edges[0].Source, "post:alpha")
	}
}

func TestGardenViewPlugin_CooccurrenceKey_Deterministic(t *testing.T) {
	p := newTestGardenPlugin()

	key1 := p.cooccurrenceKey("go", "python")
	key2 := p.cooccurrenceKey("python", "go")

	if key1 != key2 {
		t.Errorf("cooccurrenceKey is not symmetric: %q != %q", key1, key2)
	}
}

func TestGardenViewPlugin_ExportGraphJSON(t *testing.T) {
	p := newTestGardenPlugin()
	dir := t.TempDir()

	graph := GardenGraph{
		Nodes: []GardenNode{
			{ID: "post:test", Type: "post", Label: "Test Post", Href: "/test/"},
			{ID: "tag:go", Type: "tag", Label: "go", Href: "/tags/go/", Count: 1},
		},
		Edges: []GardenEdge{
			{Source: "post:test", Target: "tag:go", Type: "tag"},
		},
	}

	err := p.exportGraphJSON(dir, &graph)
	if err != nil {
		t.Fatalf("exportGraphJSON() error = %v", err)
	}

	// Verify file exists
	outputPath := filepath.Join(dir, "graph.json")
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("reading graph.json: %v", err)
	}

	// Verify it's valid JSON
	var parsed GardenGraph
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("graph.json is not valid JSON: %v", err)
	}

	// Verify content
	if len(parsed.Nodes) != 2 {
		t.Errorf("expected 2 nodes in JSON, got %d", len(parsed.Nodes))
	}
	if len(parsed.Edges) != 1 {
		t.Errorf("expected 1 edge in JSON, got %d", len(parsed.Edges))
	}
}

func TestGardenViewPlugin_BuildTagClusters(t *testing.T) {
	p := newTestGardenPlugin()

	graph := GardenGraph{
		Nodes: []GardenNode{
			{ID: "tag:go", Type: "tag", Label: "go", Href: "/tags/go/", Count: 5},
			{ID: "tag:python", Type: "tag", Label: "python", Href: "/tags/python/", Count: 3},
			{ID: "tag:tutorial", Type: "tag", Label: "tutorial", Href: "/tags/tutorial/", Count: 2},
		},
		Edges: []GardenEdge{
			{Source: "tag:go", Target: "tag:tutorial", Type: "co-occurrence", Weight: 2},
			{Source: "tag:python", Target: "tag:tutorial", Type: "co-occurrence", Weight: 1},
		},
	}

	clusters := p.buildTagClusters(&graph)

	if len(clusters) != 3 {
		t.Fatalf("expected 3 clusters, got %d", len(clusters))
	}

	// Clusters should be sorted by count descending
	if clusters[0].Name != "go" {
		t.Errorf("clusters[0].Name = %q, want %q", clusters[0].Name, "go")
	}
	if clusters[0].Count != 5 {
		t.Errorf("clusters[0].Count = %d, want %d", clusters[0].Count, 5)
	}

	// "go" should have "tutorial" as related
	if len(clusters[0].Related) != 1 || clusters[0].Related[0].Name != "tutorial" {
		t.Errorf("clusters[0].Related = %v, want [tutorial]", clusters[0].Related)
	}
	if clusters[0].Related[0].Count != 2 {
		t.Errorf("clusters[0].Related[0].Count = %d, want %d", clusters[0].Related[0].Count, 2)
	}

	// "tutorial" should have both "go" and "python" as related
	var tutorialCluster *TagCluster
	for i := range clusters {
		if clusters[i].Name == "tutorial" {
			tutorialCluster = &clusters[i]
			break
		}
	}
	if tutorialCluster == nil {
		t.Fatal("missing tutorial cluster")
	}
	if len(tutorialCluster.Related) != 2 {
		t.Errorf("tutorial.Related = %v, want 2 related tags", tutorialCluster.Related)
	}
	if tutorialCluster.Related[0].Count < tutorialCluster.Related[1].Count {
		t.Errorf("tutorial.Related not sorted by count desc: %v", tutorialCluster.Related)
	}
}

func TestGardenViewPlugin_Write_Disabled(t *testing.T) {
	p := newTestGardenPlugin()
	m := lifecycle.NewManager()
	dir := t.TempDir()

	config := lifecycle.NewConfig()
	config.OutputDir = dir
	gardenConfig := models.NewGardenConfig()
	gardenConfig.Enabled = boolPtrG(false)
	config.Extra["garden"] = gardenConfig
	config.Extra["models_config"] = &models.Config{Garden: gardenConfig, OutputDir: dir}
	m.SetConfig(config)

	date := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	m.SetPosts([]*models.Post{
		{Slug: "test", Published: true, Tags: []string{"go"}, Date: timePtrG(date)},
	})

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// No files should be created
	gardenDir := filepath.Join(dir, "garden")
	if _, err := os.Stat(gardenDir); err == nil {
		t.Error("expected garden directory to not be created when disabled")
	}
}

func TestGardenViewPlugin_Write_NoPosts(t *testing.T) {
	p := newTestGardenPlugin()
	m := lifecycle.NewManager()
	dir := t.TempDir()

	config := lifecycle.NewConfig()
	config.OutputDir = dir
	gardenConfig := models.NewGardenConfig()
	config.Extra["garden"] = gardenConfig
	config.Extra["models_config"] = &models.Config{Garden: gardenConfig, OutputDir: dir}
	m.SetConfig(config)
	m.SetPosts([]*models.Post{})

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}
}

func TestGardenViewPlugin_Write_ExportJSON(t *testing.T) {
	p := newTestGardenPlugin()
	m := lifecycle.NewManager()
	dir := t.TempDir()

	gardenConfig := models.NewGardenConfig()
	gardenConfig.RenderPage = boolPtrG(false) // Only test JSON export

	config := lifecycle.NewConfig()
	config.OutputDir = dir
	config.Extra["garden"] = gardenConfig
	config.Extra["models_config"] = &models.Config{Garden: gardenConfig, OutputDir: dir}
	m.SetConfig(config)

	m.SetPosts(newTestPosts())

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check graph.json was created
	graphPath := filepath.Join(dir, "garden", "graph.json")
	data, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("expected graph.json to be created: %v", err)
	}

	// Verify it's valid JSON with expected structure
	var graph GardenGraph
	if err := json.Unmarshal(data, &graph); err != nil {
		t.Fatalf("graph.json is not valid JSON: %v", err)
	}

	// Should have post and tag nodes
	postCount, tagCount := 0, 0
	for _, node := range graph.Nodes {
		switch node.Type {
		case "post":
			postCount++
		case "tag":
			tagCount++
		}
	}
	if postCount != 3 {
		t.Errorf("expected 3 post nodes, got %d", postCount)
	}
	if tagCount != 4 {
		t.Errorf("expected 4 tag nodes (go, programming, tutorial, python), got %d", tagCount)
	}

	// Verify deterministic ordering: nodes sorted by ID
	for i := 1; i < len(graph.Nodes); i++ {
		if graph.Nodes[i].ID < graph.Nodes[i-1].ID {
			t.Errorf("nodes not sorted: %q comes after %q", graph.Nodes[i].ID, graph.Nodes[i-1].ID)
		}
	}
}

func TestGardenViewPlugin_FilterTags(t *testing.T) {
	p := newTestGardenPlugin()
	config := newTestGardenConfig()
	config.ExcludeTags = []string{"secret", "internal"}

	tags := []string{"go", "secret", "tutorial", "internal", "python"}
	filtered := p.filterTags(tags, &config)

	expected := []string{"go", "tutorial", "python"}
	if len(filtered) != len(expected) {
		t.Fatalf("filterTags() returned %d tags, want %d", len(filtered), len(expected))
	}
	for i, want := range expected {
		if filtered[i] != want {
			t.Errorf("filterTags()[%d] = %q, want %q", i, filtered[i], want)
		}
	}
}

func TestGardenConfig_Defaults(t *testing.T) {
	config := models.NewGardenConfig()

	if !config.IsEnabled() {
		t.Error("expected IsEnabled() == true by default")
	}
	if !config.IsExportJSON() {
		t.Error("expected IsExportJSON() == true by default")
	}
	if !config.IsRenderPage() {
		t.Error("expected IsRenderPage() == true by default")
	}
	if !config.IsIncludeTags() {
		t.Error("expected IsIncludeTags() == true by default")
	}
	if !config.IsIncludePosts() {
		t.Error("expected IsIncludePosts() == true by default")
	}
	if config.GetPath() != "garden" {
		t.Errorf("GetPath() = %q, want %q", config.GetPath(), "garden")
	}
	if config.GetTemplate() != "garden.html" {
		t.Errorf("GetTemplate() = %q, want %q", config.GetTemplate(), "garden.html")
	}
	if config.GetMaxNodes() != 2000 {
		t.Errorf("GetMaxNodes() = %d, want %d", config.GetMaxNodes(), 2000)
	}
}

func TestGardenConfig_NilDefaults(t *testing.T) {
	config := models.GardenConfig{}

	// All nil pointers should default to true
	if !config.IsEnabled() {
		t.Error("nil Enabled should default to true")
	}
	if !config.IsExportJSON() {
		t.Error("nil ExportJSON should default to true")
	}
	if !config.IsRenderPage() {
		t.Error("nil RenderPage should default to true")
	}
	if !config.IsIncludeTags() {
		t.Error("nil IncludeTags should default to true")
	}
	if !config.IsIncludePosts() {
		t.Error("nil IncludePosts should default to true")
	}

	// Zero-value strings should return defaults
	if config.GetPath() != "garden" {
		t.Errorf("empty Path should default to %q", "garden")
	}
	if config.GetTemplate() != "garden.html" {
		t.Errorf("empty Template should default to %q", "garden.html")
	}
	if config.GetMaxNodes() != 2000 {
		t.Errorf("zero MaxNodes should default to %d", 2000)
	}
}

func TestGardenConfig_IsTagExcluded(t *testing.T) {
	config := models.GardenConfig{
		ExcludeTags: []string{"secret", "internal"},
	}

	if !config.IsTagExcluded("secret") {
		t.Error("expected 'secret' to be excluded")
	}
	if !config.IsTagExcluded("internal") {
		t.Error("expected 'internal' to be excluded")
	}
	if config.IsTagExcluded("go") {
		t.Error("'go' should not be excluded")
	}
}

func TestComputeGardenHash_UsesCachedSemanticHashes(t *testing.T) {
	posts := newTestPosts()
	cache := buildcache.New("")
	for _, post := range posts {
		cache.UpdatePostSemanticHashes(post.Path, "feed-"+post.Slug, computePostTagIndexHash(post), computePostGardenHash(post))
	}

	config := newTestGardenConfig()
	want := computeGardenHash(posts, &config, nil)
	got := computeGardenHash(posts, &config, cache)
	if got != want {
		t.Fatalf("computeGardenHash with cache = %q, want %q", got, want)
	}
}

func TestComputeGardenHash_ChangesWhenConfigChanges(t *testing.T) {
	posts := newTestPosts()
	base := newTestGardenConfig()
	baseHash := computeGardenHash(posts, &base, nil)

	modified := newTestGardenConfig()
	modified.Title = "Custom Garden"
	modified.ExcludeTags = []string{"go"}
	modifiedHash := computeGardenHash(posts, &modified, nil)

	if baseHash == modifiedHash {
		t.Fatal("computeGardenHash did not change after config update")
	}
}
