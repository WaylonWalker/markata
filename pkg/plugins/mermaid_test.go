package plugins

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

type countingMermaidRenderer struct {
	calls int
	svg   string
	err   error
}

func (r *countingMermaidRenderer) render(_ string) (string, error) {
	r.calls++
	if r.err != nil {
		return "", r.err
	}
	return r.svg, nil
}

func (r *countingMermaidRenderer) close() error { return nil }

func TestMermaidPlugin_Name(t *testing.T) {
	p := NewMermaidPlugin()
	if p.Name() != "mermaid" {
		t.Errorf("expected name 'mermaid', got %q", p.Name())
	}
}

func TestMermaidPlugin_DefaultConfig(t *testing.T) {
	p := NewMermaidPlugin()
	config := p.Config()

	if !config.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if config.CDNURL != "/assets/vendor/mermaid/mermaid.min.js" {
		t.Errorf("unexpected default CDN URL: %q", config.CDNURL)
	}
	if config.Theme != "default" {
		t.Errorf("expected default theme 'default', got %q", config.Theme)
	}
}

func TestMermaidPlugin_Configure(t *testing.T) {
	tests := []struct {
		name        string
		extra       map[string]interface{}
		wantEnabled bool
		wantCDNURL  string
		wantTheme   string
	}{
		{
			name:        "no config",
			extra:       nil,
			wantEnabled: true,
			wantCDNURL:  "/assets/vendor/mermaid/mermaid.min.js",
			wantTheme:   "default",
		},
		{
			name:        "empty extra",
			extra:       map[string]interface{}{},
			wantEnabled: true,
			wantCDNURL:  "/assets/vendor/mermaid/mermaid.min.js",
			wantTheme:   "default",
		},
		{
			name: "custom config",
			extra: map[string]interface{}{
				"mermaid": map[string]interface{}{
					"enabled": false,
					"cdn_url": "https://example.com/mermaid.js",
					"theme":   "dark",
				},
			},
			wantEnabled: false,
			wantCDNURL:  "https://example.com/mermaid.js",
			wantTheme:   "dark",
		},
		{
			name: "partial config - theme only",
			extra: map[string]interface{}{
				"mermaid": map[string]interface{}{
					"theme": "forest",
				},
			},
			wantEnabled: true,
			wantCDNURL:  "/assets/vendor/mermaid/mermaid.min.js",
			wantTheme:   "forest",
		},
		{
			name: "disabled",
			extra: map[string]interface{}{
				"mermaid": map[string]interface{}{
					"enabled": false,
				},
			},
			wantEnabled: false,
			wantCDNURL:  "/assets/vendor/mermaid/mermaid.min.js",
			wantTheme:   "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewMermaidPlugin()
			m := lifecycle.NewManager()
			m.Config().Extra = tt.extra

			if err := p.Configure(m); err != nil {
				t.Fatalf("Configure returned error: %v", err)
			}

			config := p.Config()
			if config.Enabled != tt.wantEnabled {
				t.Errorf("Enabled = %v, want %v", config.Enabled, tt.wantEnabled)
			}
			if config.CDNURL != tt.wantCDNURL {
				t.Errorf("CDNURL = %q, want %q", config.CDNURL, tt.wantCDNURL)
			}
			if config.Theme != tt.wantTheme {
				t.Errorf("Theme = %q, want %q", config.Theme, tt.wantTheme)
			}
		})
	}
}

func TestMermaidPlugin_Priority(t *testing.T) {
	p := NewMermaidPlugin()

	// Should run late in the render stage (after render_markdown)
	renderPriority := p.Priority(lifecycle.StageRender)
	if renderPriority != lifecycle.PriorityLate {
		t.Errorf("render stage priority = %d, want %d (PriorityLate)", renderPriority, lifecycle.PriorityLate)
	}

	// Other stages should use default priority
	transformPriority := p.Priority(lifecycle.StageTransform)
	if transformPriority != lifecycle.PriorityDefault {
		t.Errorf("transform stage priority = %d, want %d (PriorityDefault)", transformPriority, lifecycle.PriorityDefault)
	}
}

func TestMermaidPlugin_BasicFlowchart(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A[Start] --&gt; B{Decision}
    B --&gt;|Yes| C[Action 1]
    B --&gt;|No| D[Action 2]</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should have mermaid pre block
	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected <pre class=\"mermaid\"> in output, got %q", post.ArticleHTML)
	}

	// Should NOT have the old code block
	if strings.Contains(post.ArticleHTML, `class="language-mermaid"`) {
		t.Errorf("should not contain language-mermaid class after processing, got %q", post.ArticleHTML)
	}

	// Should have the diagram code (HTML entities decoded)
	if !strings.Contains(post.ArticleHTML, "graph TD") {
		t.Errorf("expected diagram code in output, got %q", post.ArticleHTML)
	}

	// Should have decoded HTML entities
	if !strings.Contains(post.ArticleHTML, "-->") {
		t.Errorf("expected decoded arrows (-->) in output, got %q", post.ArticleHTML)
	}

	// Should have injected script
	if !strings.Contains(post.ArticleHTML, "<script type=\"module\">") {
		t.Errorf("expected mermaid script in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "mermaid.initialize") {
		t.Errorf("expected mermaid.initialize in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_SequenceDiagram(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">sequenceDiagram
    Alice-&gt;&gt;Bob: Hello Bob
    Bob--&gt;&gt;Alice: Hi Alice</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "sequenceDiagram") {
		t.Errorf("expected sequenceDiagram in output, got %q", post.ArticleHTML)
	}
	// HTML entities should be decoded
	if !strings.Contains(post.ArticleHTML, "->>") {
		t.Errorf("expected decoded arrows (->>) in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_ClassDiagram(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">classDiagram
    class Animal {
        +name string
        +makeSound()
    }</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "classDiagram") {
		t.Errorf("expected classDiagram in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_PieChart(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">pie title Pets
    "Dogs" : 386
    "Cats" : 85
    "Birds" : 15</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "pie title Pets") {
		t.Errorf("expected pie chart in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_GanttChart(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">gantt
    title A Gantt Diagram
    section Section
    Task 1 :a1, 2024-01-01, 30d</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "gantt") {
		t.Errorf("expected gantt chart in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_ERDiagram(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "erDiagram") {
		t.Errorf("expected erDiagram in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_StateDiagram(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">stateDiagram-v2
    [*] --&gt; Still
    Still --&gt; Moving
    Moving --&gt; Still</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "stateDiagram-v2") {
		t.Errorf("expected stateDiagram in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_GitGraph(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">gitGraph
    commit
    branch develop
    commit
    checkout main
    merge develop</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "gitGraph") {
		t.Errorf("expected gitGraph in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_Mindmap(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">mindmap
  root((mindmap))
    Origins
    Research</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "mindmap") {
		t.Errorf("expected mindmap in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_Timeline(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">timeline
    title History of Social Media
    2002 : LinkedIn
    2004 : Facebook</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "timeline") {
		t.Errorf("expected timeline in output, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_MultipleDiagrams(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<p>First diagram:</p>
<pre><code class="language-mermaid">graph LR
    A --&gt; B</code></pre>
<p>Second diagram:</p>
<pre><code class="language-mermaid">pie title Test
    "A" : 50
    "B" : 50</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Count mermaid pre blocks
	count := strings.Count(post.ArticleHTML, `<pre class="mermaid">`)
	if count != 2 {
		t.Errorf("expected 2 mermaid pre blocks, got %d in %q", count, post.ArticleHTML)
	}

	// Script should only be injected once
	scriptCount := strings.Count(post.ArticleHTML, "<script type=\"module\">")
	if scriptCount != 1 {
		t.Errorf("expected 1 mermaid script injection, got %d", scriptCount)
	}
}

func TestMermaidPlugin_NoMermaid(t *testing.T) {
	p := NewMermaidPlugin()
	originalHTML := `<p>Hello world</p>
<pre><code class="language-python">print("hello")</code></pre>`
	post := &models.Post{
		ArticleHTML: originalHTML,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should not inject script when no mermaid blocks
	if strings.Contains(post.ArticleHTML, "<script type=\"module\">") {
		t.Errorf("should not inject script when no mermaid blocks, got %q", post.ArticleHTML)
	}

	// Original content should be preserved
	if !strings.Contains(post.ArticleHTML, `class="language-python"`) {
		t.Errorf("original content should be preserved, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_SkipPost(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A --> B</code></pre>`,
		Skip: true,
	}

	originalHTML := post.ArticleHTML
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// HTML should remain unchanged for skipped posts
	if post.ArticleHTML != originalHTML {
		t.Errorf("skipped post HTML should not change, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_EmptyContent(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: "",
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if post.ArticleHTML != "" {
		t.Errorf("empty content should remain empty, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_DisabledPlugin(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled: false,
		CDNURL:  "https://example.com/mermaid.js",
		Theme:   "dark",
	})

	m := lifecycle.NewManager()
	posts := []*models.Post{
		{
			ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A --> B</code></pre>`,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// HTML should remain unchanged when plugin is disabled
	if strings.Contains(posts[0].ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("disabled plugin should not process mermaid, got %q", posts[0].ArticleHTML)
	}
}

func TestMermaidPlugin_CustomTheme(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled:         true,
		Mode:            "client",
		CDNURL:          "/assets/vendor/mermaid/mermaid.min.js",
		Theme:           "dark",
		UseCSSVariables: false,
	})

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A --&gt; B</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `theme: 'dark'`) {
		t.Errorf("expected dark theme in script, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_CustomCDN(t *testing.T) {
	p := NewMermaidPlugin()
	customCDN := "https://example.com/custom/mermaid.js"
	p.SetConfig(models.MermaidConfig{
		Enabled:         true,
		Mode:            "client",
		CDNURL:          customCDN,
		Theme:           "default",
		UseCSSVariables: false,
	})

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A --&gt; B</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, customCDN) {
		t.Errorf("expected custom CDN URL in script, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_Render(t *testing.T) {
	p := NewMermaidPlugin()
	m := lifecycle.NewManager()

	posts := []*models.Post{
		{
			ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A --&gt; B</code></pre>`,
		},
		{
			ArticleHTML: `<p>No mermaid here</p>`,
		},
		{
			ArticleHTML: `<pre><code class="language-mermaid">pie title Test
    "A" : 50</code></pre>`,
			Skip: true,
		},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	resultPosts := m.Posts()

	// First post should be processed
	if !strings.Contains(resultPosts[0].ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("Post 1 should have mermaid pre block: %q", resultPosts[0].ArticleHTML)
	}

	// Second post should not have script (no mermaid)
	if strings.Contains(resultPosts[1].ArticleHTML, "<script type=\"module\">") {
		t.Errorf("Post 2 should not have mermaid script: %q", resultPosts[1].ArticleHTML)
	}

	// Third post should be unchanged (skipped)
	if strings.Contains(resultPosts[2].ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("Skipped post should not be processed: %q", resultPosts[2].ArticleHTML)
	}
}

func TestMermaidPlugin_HTMLEntitiesDecoded(t *testing.T) {
	p := NewMermaidPlugin()
	// Test various HTML entities that goldmark might produce
	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A[&quot;Start&quot;] --&gt; B{&lt;Decision&gt;}
    B --&gt;|&amp;Yes&amp;| C</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// All HTML entities should be decoded
	if strings.Contains(post.ArticleHTML, "&quot;") {
		t.Errorf("&quot; should be decoded, got %q", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "&gt;") {
		t.Errorf("&gt; should be decoded, got %q", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "&lt;") {
		t.Errorf("&lt; should be decoded, got %q", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "&amp;") {
		t.Errorf("&amp; should be decoded, got %q", post.ArticleHTML)
	}

	// Check decoded values are present
	if !strings.Contains(post.ArticleHTML, `"Start"`) {
		t.Errorf("expected decoded quotes, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "-->") {
		t.Errorf("expected decoded arrows, got %q", post.ArticleHTML)
	}
}

func TestMermaidPlugin_PreserveSurroundingContent(t *testing.T) {
	p := NewMermaidPlugin()
	post := &models.Post{
		ArticleHTML: `<h1>Title</h1>
<p>Some text before</p>
<pre><code class="language-mermaid">graph TD
    A --&gt; B</code></pre>
<p>Some text after</p>
<pre><code class="language-python">print("keep me")</code></pre>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Surrounding content should be preserved
	if !strings.Contains(post.ArticleHTML, "<h1>Title</h1>") {
		t.Errorf("h1 should be preserved, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<p>Some text before</p>") {
		t.Errorf("text before should be preserved, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<p>Some text after</p>") {
		t.Errorf("text after should be preserved, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `class="language-python"`) {
		t.Errorf("python code block should be preserved, got %q", post.ArticleHTML)
	}
}

// Interface compliance tests

func TestMermaidPlugin_Interfaces(_ *testing.T) {
	p := NewMermaidPlugin()

	// Verify interface compliance
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.RenderPlugin = p
	var _ lifecycle.PriorityPlugin = p
}

func TestMermaidPlugin_RenderDiagram_UsesCache(t *testing.T) {
	cache := buildcache.New(t.TempDir())

	renderer := &countingMermaidRenderer{svg: "<svg>cached</svg>"}
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled:         true,
		Mode:            mermaidModeChromium,
		Theme:           "default",
		UseCSSVariables: true,
		Lightbox:        true,
		ChromiumConfig: &models.ChromiumRendererConfig{
			Timeout:       30,
			MaxConcurrent: 4,
			NoSandbox:     true,
		},
	})
	p.renderer = renderer
	p.cache = cache

	post := &models.Post{Path: "pages/diagram.md"}
	first, err := p.renderDiagram(post, "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("first renderDiagram: %v", err)
	}
	second, err := p.renderDiagram(post, "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("second renderDiagram: %v", err)
	}

	if first != "<svg>cached</svg>" || second != "<svg>cached</svg>" {
		t.Fatalf("unexpected SVG outputs: %q %q", first, second)
	}
	if renderer.calls != 1 {
		t.Fatalf("renderer calls = %d, want 1", renderer.calls)
	}
}

func TestMermaidPlugin_RenderDiagram_CacheKeyIncludesTheme(t *testing.T) {
	cache := buildcache.New(t.TempDir())

	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{Enabled: true, Mode: mermaidModeChromium, Theme: "default"})
	p.renderKey = "chromium:/usr/bin/chromium:Chromium 1:mermaidjs-v10"
	p.renderer = &countingMermaidRenderer{svg: "<svg>default</svg>"}
	p.cache = cache

	post := &models.Post{Path: "pages/diagram.md"}
	if _, err := p.renderDiagram(post, "graph TD\nA-->B"); err != nil {
		t.Fatalf("render default: %v", err)
	}

	renderer := &countingMermaidRenderer{svg: "<svg>dark</svg>"}
	p.SetConfig(models.MermaidConfig{Enabled: true, Mode: mermaidModeChromium, Theme: "dark"})
	p.renderKey = "chromium:/usr/bin/chromium:Chromium 1:mermaidjs-v10"
	p.renderer = renderer
	p.cache = cache

	got, err := p.renderDiagram(post, "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("render dark: %v", err)
	}
	if got != "<svg>dark</svg>" {
		t.Fatalf("got %q, want dark SVG", got)
	}
	if renderer.calls != 1 {
		t.Fatalf("renderer calls = %d, want 1", renderer.calls)
	}
}

func TestMermaidPlugin_RenderDiagram_PropagatesRenderErrors(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{Enabled: true, Mode: mermaidModeChromium})
	p.renderKey = "chromium:/usr/bin/chromium:Chromium 1:mermaidjs-v10"
	p.renderer = &countingMermaidRenderer{err: errors.New("render failed")}

	_, err := p.renderDiagram(&models.Post{Path: "pages/diagram.md"}, "graph TD\nA-->B")
	if err == nil {
		t.Fatal("expected render error")
	}
}

func TestMermaidPlugin_RenderDiagram_CacheKeyIncludesRendererFingerprint(t *testing.T) {
	cache := buildcache.New(t.TempDir())

	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{Enabled: true, Mode: mermaidModeChromium, Theme: "default"})
	p.renderKey = "chromium:/usr/bin/chromium:Chromium 1:mermaidjs-v10"
	p.renderer = &countingMermaidRenderer{svg: "<svg>old</svg>"}
	p.cache = cache

	post := &models.Post{Path: "pages/diagram.md"}
	if _, err := p.renderDiagram(post, "graph TD\nA-->B"); err != nil {
		t.Fatalf("render old: %v", err)
	}

	renderer := &countingMermaidRenderer{svg: "<svg>new</svg>"}
	p.renderKey = "chromium:/usr/bin/chromium:Chromium 2:mermaidjs-v10"
	p.renderer = renderer

	got, err := p.renderDiagram(post, "graph TD\nA-->B")
	if err != nil {
		t.Fatalf("render new: %v", err)
	}
	if got != "<svg>new</svg>" {
		t.Fatalf("got %q, want new SVG", got)
	}
	if renderer.calls != 1 {
		t.Fatalf("renderer calls = %d, want 1", renderer.calls)
	}
}

func TestMermaidPlugin_DiagramRenderHashIncludesCLIConfigFileContents(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "mermaid.json")
	if err := os.WriteFile(configPath, []byte(`{"theme":"default"}`), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled: true,
		Mode:    mermaidModeCLI,
		Theme:   "default",
		CLIConfig: &models.CLIRendererConfig{
			ExtraArgs: "--configFile " + configPath,
		},
	})
	p.renderKey = "cli:/usr/bin/mmdc:1.0.0"

	first := p.diagramRenderHash("graph TD\nA-->B")
	if err := os.WriteFile(configPath, []byte(`{"theme":"dark"}`), 0o600); err != nil {
		t.Fatalf("rewrite config file: %v", err)
	}
	second := p.diagramRenderHash("graph TD\nA-->B")

	if first == second {
		t.Fatal("diagram render hash did not change after CLI config file contents changed")
	}
}

func TestMermaidPlugin_DiagramRenderHashIncludesCLIEqualsFileArg(t *testing.T) {
	dir := t.TempDir()
	cssPath := filepath.Join(dir, "mermaid.css")
	if err := os.WriteFile(cssPath, []byte("svg { color: black; }"), 0o600); err != nil {
		t.Fatalf("write css file: %v", err)
	}

	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled: true,
		Mode:    mermaidModeCLI,
		Theme:   "default",
		CLIConfig: &models.CLIRendererConfig{
			ExtraArgs: "--cssFile=" + cssPath,
		},
	})
	p.renderKey = "cli:/usr/bin/mmdc:1.0.0"

	first := p.diagramRenderHash("graph TD\nA-->B")
	if err := os.WriteFile(cssPath, []byte("svg { color: white; }"), 0o600); err != nil {
		t.Fatalf("rewrite css file: %v", err)
	}
	second := p.diagramRenderHash("graph TD\nA-->B")

	if first == second {
		t.Fatal("diagram render hash did not change after CLI css file contents changed")
	}
}

// TestMermaidPlugin_CLIMode_Fallback tests that CLI mode falls back to client rendering
// when rendering fails (e.g., mmdc not found). This verifies graceful degradation.
func TestMermaidPlugin_CLIMode_Fallback(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled: true,
		Mode:    "cli",
		Theme:   "default",
		CDNURL:  "/assets/vendor/mermaid/mermaid.min.js",
		CLIConfig: &models.CLIRendererConfig{
			// Use a non-existent path to force fallback
			MMDCPath:  "/nonexistent/mmdc",
			ExtraArgs: "",
		},
	})

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A --&gt; B</code></pre>`,
	}

	// processPost should not error even if mmdc is missing
	// It should fall back to client rendering (mermaid pre block)
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should fall back to client rendering (mermaid pre block)
	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected fallback to mermaid pre block, got %q", post.ArticleHTML)
	}
}

// TestMermaidPlugin_ChromiumMode_Fallback tests that Chromium mode falls back to client rendering
// when rendering fails (e.g., browser not found). This verifies graceful degradation.
func TestMermaidPlugin_ChromiumMode_Fallback(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled: true,
		Mode:    "chromium",
		Theme:   "default",
		CDNURL:  "/assets/vendor/mermaid/mermaid.min.js",
		ChromiumConfig: &models.ChromiumRendererConfig{
			// Use a non-existent path to force fallback
			BrowserPath:   "/nonexistent/chrome",
			Timeout:       30,
			MaxConcurrent: 4,
		},
	})

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">pie title Test
    "A" : 50</code></pre>`,
	}

	// processPost should not error even if browser is missing
	// It should fall back to client rendering (mermaid pre block)
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should fall back to client rendering (mermaid pre block)
	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected fallback to mermaid pre block, got %q", post.ArticleHTML)
	}
}

// TestMermaidPlugin_InvalidMode_DefaultsToClient tests that an invalid mode
// defaults to client mode behavior when no explicit mode is set.
func TestMermaidPlugin_InvalidMode_DefaultsToClient(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled:         true,
		Mode:            "invalid", // Invalid mode
		Theme:           "dark",
		UseCSSVariables: false,
		CDNURL:          "/assets/vendor/mermaid/mermaid.min.js",
	})

	post := &models.Post{
		ArticleHTML: `<pre><code class="language-mermaid">graph TD
    A --&gt; B</code></pre>`,
	}

	// processPost should treat invalid mode as an unknown mode and not inject script
	// (since it only handles "client", "cli", and "chromium")
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// With an invalid mode, the code block should be converted but no script injected
	// (line 236 checks p.config.Mode == "client" specifically)
	if !strings.Contains(post.ArticleHTML, `<pre class="mermaid">`) {
		t.Errorf("expected mermaid pre block, got %q", post.ArticleHTML)
	}
}

// TestMermaidPlugin_MultipleModesInPost tests that a single post with multiple
// mermaid diagrams processes all of them.
func TestMermaidPlugin_MultipleModesInPost(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled:         true,
		Mode:            "client",
		Theme:           "dark",
		UseCSSVariables: false,
		CDNURL:          "/assets/vendor/mermaid/mermaid.min.js",
	})

	post := &models.Post{
		ArticleHTML: `
<p>First diagram:</p>
<pre><code class="language-mermaid">graph TD
    A --&gt; B</code></pre>
<p>Second diagram:</p>
<pre><code class="language-mermaid">pie title Test
    "A" : 50</code></pre>
`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Count how many mermaid pre blocks are in the result
	count := strings.Count(post.ArticleHTML, `<pre class="mermaid">`)
	if count != 2 {
		t.Errorf("expected 2 mermaid pre blocks, got %d", count)
	}

	// Script should be injected once
	if !strings.Contains(post.ArticleHTML, `theme: 'dark'`) {
		t.Errorf("expected dark theme in script, got %q", post.ArticleHTML)
	}
}

// TestMermaidPlugin_CLIConfig tests that CLI configuration options are properly loaded.
func TestMermaidPlugin_CLIConfig(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled: true,
		Mode:    "cli",
		CLIConfig: &models.CLIRendererConfig{
			MMDCPath:  "/usr/local/bin/mmdc",
			ExtraArgs: "--theme dark --width 1024",
		},
	})

	cfg := p.Config()
	if cfg.CLIConfig.MMDCPath != "/usr/local/bin/mmdc" {
		t.Errorf("expected MMDCPath /usr/local/bin/mmdc, got %q", cfg.CLIConfig.MMDCPath)
	}
	if cfg.CLIConfig.ExtraArgs != "--theme dark --width 1024" {
		t.Errorf("expected ExtraArgs, got %q", cfg.CLIConfig.ExtraArgs)
	}
}

// TestMermaidPlugin_ChromiumConfig tests that Chromium configuration options are properly loaded.
func TestMermaidPlugin_ChromiumConfig(t *testing.T) {
	p := NewMermaidPlugin()
	p.SetConfig(models.MermaidConfig{
		Enabled: true,
		Mode:    "chromium",
		ChromiumConfig: &models.ChromiumRendererConfig{
			BrowserPath:   "/usr/bin/chromium",
			Timeout:       60,
			MaxConcurrent: 8,
		},
	})

	cfg := p.Config()
	if cfg.ChromiumConfig.BrowserPath != "/usr/bin/chromium" {
		t.Errorf("expected BrowserPath /usr/bin/chromium, got %q", cfg.ChromiumConfig.BrowserPath)
	}
	if cfg.ChromiumConfig.Timeout != 60 {
		t.Errorf("expected Timeout 60, got %d", cfg.ChromiumConfig.Timeout)
	}
	if cfg.ChromiumConfig.MaxConcurrent != 8 {
		t.Errorf("expected MaxConcurrent 8, got %d", cfg.ChromiumConfig.MaxConcurrent)
	}
}

func TestMermaidCSSVariablesUseThemeSurfaceAndBorder(t *testing.T) {
	script := NewMermaidPlugin().cssVariablesScript()
	for _, expected := range []string{
		"primaryColor: surface",
		"nodeBkg: surface",
		"edgeLabelBackground: surface",
		"primaryTextColor: foreground",
		"nodeTextColor: foreground",
		"primaryBorderColor: border",
		"nodeBorder: border",
		"lineColor: border",
		"clusterBorder: border",
		"clusterBkg: background",
		"new WeakMap()",
		"removeAttribute('data-processed')",
		"diagramSources.set(block, block.textContent || '')",
		"window.addEventListener('theme-change'",
		"window.addEventListener('palette-change'",
		"window.addEventListener('color-mode-change'",
	} {
		if !strings.Contains(script, expected) {
			t.Errorf("Mermaid theme script missing %q", expected)
		}
	}
	if strings.Contains(script, "prefers-color-scheme: dark") {
		t.Error("Mermaid palette should not override selected site theme with OS preference")
	}
	if strings.Contains(script, "* { cursor: pointer; }") {
		t.Error("Mermaid should not make every diagram element appear interactive")
	}
	if strings.Contains(script, "nodeBkg: css('--color-code-bg'") {
		t.Error("Mermaid node fill should not use code background")
	}
}
