package plugins

import (
	"fmt"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestHeadingAnchorsPlugin_Name(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	if p.Name() != "heading_anchors" {
		t.Errorf("expected name 'heading_anchors', got %q", p.Name())
	}
}

func TestHeadingAnchorsPlugin_Priority(t *testing.T) {
	p := NewHeadingAnchorsPlugin()

	// Should have late priority for render stage
	if p.Priority(lifecycle.StageRender) != lifecycle.PriorityLate {
		t.Errorf("expected PriorityLate for render stage, got %d", p.Priority(lifecycle.StageRender))
	}

	// Should have default priority for other stages
	if p.Priority(lifecycle.StageTransform) != lifecycle.PriorityDefault {
		t.Errorf("expected PriorityDefault for transform stage, got %d", p.Priority(lifecycle.StageTransform))
	}
}

func TestHeadingAnchorsPlugin_Configure(t *testing.T) {
	tests := []struct {
		name       string
		config     map[string]interface{}
		wantMin    int
		wantMax    int
		wantPos    string
		wantSymbol string
		wantClass  string
	}{
		{
			name:       "default values",
			config:     nil,
			wantMin:    2,
			wantMax:    4,
			wantPos:    "end",
			wantSymbol: "#",
			wantClass:  "heading-anchor",
		},
		{
			name: "custom values",
			config: map[string]interface{}{
				"heading_anchors": map[string]interface{}{
					"min_level": 1,
					"max_level": 6,
					"position":  "start",
					"symbol":    "link",
					"class":     "anchor-link",
				},
			},
			wantMin:    1,
			wantMax:    6,
			wantPos:    "start",
			wantSymbol: "link",
			wantClass:  "anchor-link",
		},
		{
			name: "partial config",
			config: map[string]interface{}{
				"heading_anchors": map[string]interface{}{
					"min_level": 3,
					"symbol":    ">>",
				},
			},
			wantMin:    3,
			wantMax:    4,
			wantPos:    "end",
			wantSymbol: ">>",
			wantClass:  "heading-anchor",
		},
		{
			name: "invalid position ignored",
			config: map[string]interface{}{
				"heading_anchors": map[string]interface{}{
					"position": "middle",
				},
			},
			wantMin:    2,
			wantMax:    4,
			wantPos:    "end", // Should remain default
			wantSymbol: "#",
			wantClass:  "heading-anchor",
		},
		{
			name: "invalid level range ignored",
			config: map[string]interface{}{
				"heading_anchors": map[string]interface{}{
					"min_level": 0,  // Invalid
					"max_level": 10, // Invalid
				},
			},
			wantMin:    2, // Should remain default
			wantMax:    4, // Should remain default
			wantPos:    "end",
			wantSymbol: "#",
			wantClass:  "heading-anchor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewHeadingAnchorsPlugin()
			m := lifecycle.NewManager()
			if tt.config != nil {
				m.Config().Extra = tt.config
			}

			if err := p.Configure(m); err != nil {
				t.Fatalf("Configure error: %v", err)
			}

			if p.minLevel != tt.wantMin {
				t.Errorf("minLevel = %d, want %d", p.minLevel, tt.wantMin)
			}
			if p.maxLevel != tt.wantMax {
				t.Errorf("maxLevel = %d, want %d", p.maxLevel, tt.wantMax)
			}
			if p.position != tt.wantPos {
				t.Errorf("position = %q, want %q", p.position, tt.wantPos)
			}
			if p.symbol != tt.wantSymbol {
				t.Errorf("symbol = %q, want %q", p.symbol, tt.wantSymbol)
			}
			if p.class != tt.wantClass {
				t.Errorf("class = %q, want %q", p.class, tt.wantClass)
			}
		})
	}
}

func TestHeadingAnchorsPlugin_BasicHeading(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	p.SetLevelRange(1, 6) // Include all levels for testing

	post := &models.Post{
		ArticleHTML: `<h2 id="my-section">My Section</h2>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should have anchor at end (default)
	expected := `<a href="#my-section" class="heading-anchor">#</a>`
	if !strings.Contains(post.ArticleHTML, expected) {
		t.Errorf("expected anchor link in output, got %q", post.ArticleHTML)
	}

	// Should preserve ID
	if !strings.Contains(post.ArticleHTML, `id="my-section"`) {
		t.Errorf("expected id attribute preserved, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_SkipsEmbeddedCardTitles(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	p.SetLevelRange(1, 6)

	card := `<h2 class="card-title p-name"><a class="u-url" href="/post/">Post</a></h2>`
	post := &models.Post{ArticleHTML: card + `<h2>Real Heading</h2>`}

	if err := p.processPost(post); err != nil {
		t.Fatalf("processPost error: %v", err)
	}
	if !strings.Contains(post.ArticleHTML, card) {
		t.Errorf("card title was modified: %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `<a href="#real-heading" class="heading-anchor">`) {
		t.Errorf("real heading should still get an anchor: %q", post.ArticleHTML)
	}
}

func TestWrapHeadingGlyphText_PreservesInlineMarkupAndWearsRegularLinks(t *testing.T) {
	input := `plain <strong>strong</strong> <em>emphasis</em> <code>code</code> <mark>highlight</mark> <a href="https://example.com">linked</a>`
	got := wrapHeadingGlyphText(input)
	for _, tag := range []string{"<strong>", "</strong>", "<em>", "</em>", "<code>", "</code>", "<mark>", "</mark>", `<a href="https://example.com">`, "</a>"} {
		if !strings.Contains(got, tag) {
			t.Fatalf("output lost %s: %s", tag, got)
		}
	}
	if strings.Count(got, `class="heading-wear-glyph"`) != 6 {
		t.Fatalf("expected one glyph wrapper for each text run, got %d: %s", strings.Count(got, `class="heading-wear-glyph"`), got)
	}
}

func TestHeadingWearPlugin_AllLevelsAndInlineMarkup(t *testing.T) {
	p := NewHeadingWearPlugin()
	p.SetEnabled(true)
	post := &models.Post{ArticleHTML: `<h1> one, <mark>two</mark></h1><h2><code>two</code></h2><h3><em>three</em></h3><h4><strong>four</strong></h4><h5><a href="/five">five</a></h5><h6> six </h6>`}
	if err := p.processPost(post); err != nil {
		t.Fatal(err)
	}
	if got := strings.Count(post.ArticleHTML, `class="heading-wear-glyph"`); got != 7 {
		t.Fatalf("wear wrapper count = %d, want 7: %s", got, post.ArticleHTML)
	}
	for _, tag := range []string{"<h1>", "<h2>", "<h3>", "<h4>", "<h5>", "<h6>", `<a href="/five">`} {
		if !strings.Contains(post.ArticleHTML, tag) {
			t.Errorf("output lost %s: %s", tag, post.ArticleHTML)
		}
	}
	if !strings.Contains(post.ArticleHTML, ` six </span>`) {
		t.Errorf("heading whitespace was not preserved: %s", post.ArticleHTML)
	}
}

func TestHeadingWearPlugin_SkipsNoneAndZeroMix(t *testing.T) {
	for _, theme := range []models.ThemeHeadingTextureConfig{
		{Kind: "none", ColorMix: 1},
		{Kind: "screenprint", ColorMix: 0},
	} {
		p := NewHeadingWearPlugin()
		m := lifecycle.NewManager()
		m.Config().Extra = map[string]interface{}{"models_config": &models.Config{Theme: models.ThemeConfig{HeadingTexture: theme}}}
		if err := p.Configure(m); err != nil {
			t.Fatal(err)
		}
		post := &models.Post{ArticleHTML: `<h1>Heading</h1>`}
		m.SetPosts([]*models.Post{post})
		if err := p.Render(m); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(post.ArticleHTML, "heading-wear-glyph") {
			t.Errorf("disabled wear changed output: %s", post.ArticleHTML)
		}
	}
}

func TestHeadingAnchorsPlugin_WearIndependentOfAnchors(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		t.Run(fmt.Sprintf("anchors-enabled-%t", enabled), func(t *testing.T) {
			p := NewHeadingAnchorsPlugin()
			p.SetEnabled(enabled)
			m := lifecycle.NewManager()
			m.Config().Extra = map[string]interface{}{"models_config": &models.Config{Theme: models.ThemeConfig{HeadingTexture: models.ThemeHeadingTextureConfig{Kind: "screenprint", ColorMix: 1}}}}
			post := &models.Post{ArticleHTML: `<h2>Heading</h2>`}
			m.SetPosts([]*models.Post{post})
			if err := p.Configure(m); err != nil {
				t.Fatal(err)
			}
			if err := p.Render(m); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(post.ArticleHTML, "heading-wear-glyph") {
				t.Errorf("wear did not run: %s", post.ArticleHTML)
			}
			hasAnchor := strings.Contains(post.ArticleHTML, "heading-anchor")
			if hasAnchor != enabled {
				t.Errorf("hasAnchor = %v, want %v: %s", hasAnchor, enabled, post.ArticleHTML)
			}
			if enabled && strings.Contains(post.ArticleHTML, `<a href="#heading"><span class="heading-wear-glyph">`) {
				t.Errorf("generated anchor text was worn: %s", post.ArticleHTML)
			}
		})
	}
}

func TestWrapHeadingGlyphText_DoesNotWrapWhitespaceOnlyNodes(t *testing.T) {
	got := wrapHeadingGlyphText("  <strong>text</strong>  ")
	if strings.Contains(got, `class="heading-wear-glyph"> </span>`) {
		t.Fatalf("wrapped whitespace: %s", got)
	}
}

func TestHeadingAnchorsPlugin_AnchorPosition(t *testing.T) {
	tests := []struct {
		name     string
		position string
		want     string
	}{
		{
			name:     "end position",
			position: "end",
			want:     `My Section <a href="#my-section"`,
		},
		{
			name:     "start position",
			position: "start",
			want:     `<a href="#my-section" class="heading-anchor">#</a> My Section`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewHeadingAnchorsPlugin()
			p.SetPosition(tt.position)

			post := &models.Post{
				ArticleHTML: `<h2 id="my-section">My Section</h2>`,
			}

			err := p.processPost(post)
			if err != nil {
				t.Fatalf("processPost error: %v", err)
			}

			if !strings.Contains(post.ArticleHTML, tt.want) {
				t.Errorf("expected %q in output, got %q", tt.want, post.ArticleHTML)
			}
		})
	}
}

func TestHeadingAnchorsPlugin_LevelRange(t *testing.T) {
	tests := []struct {
		name             string
		minLevel         int
		maxLevel         int
		input            string
		shouldHaveAnchor bool
	}{
		{
			name:             "h2 in range 2-4",
			minLevel:         2,
			maxLevel:         4,
			input:            `<h2 id="test">Test</h2>`,
			shouldHaveAnchor: true,
		},
		{
			name:             "h1 outside range 2-4",
			minLevel:         2,
			maxLevel:         4,
			input:            `<h1 id="test">Test</h1>`,
			shouldHaveAnchor: false,
		},
		{
			name:             "h5 outside range 2-4",
			minLevel:         2,
			maxLevel:         4,
			input:            `<h5 id="test">Test</h5>`,
			shouldHaveAnchor: false,
		},
		{
			name:             "h3 in range 2-4",
			minLevel:         2,
			maxLevel:         4,
			input:            `<h3 id="test">Test</h3>`,
			shouldHaveAnchor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewHeadingAnchorsPlugin()
			p.SetLevelRange(tt.minLevel, tt.maxLevel)

			post := &models.Post{ArticleHTML: tt.input}

			err := p.processPost(post)
			if err != nil {
				t.Fatalf("processPost error: %v", err)
			}

			hasAnchor := strings.Contains(post.ArticleHTML, "heading-anchor")
			if hasAnchor != tt.shouldHaveAnchor {
				t.Errorf("hasAnchor = %v, want %v; output: %q", hasAnchor, tt.shouldHaveAnchor, post.ArticleHTML)
			}
		})
	}
}

func TestHeadingAnchorsPlugin_GenerateID(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		wantID string
	}{
		{
			name:   "simple text",
			input:  `<h2>Simple Heading</h2>`,
			wantID: "simple-heading",
		},
		{
			name:   "with special characters",
			input:  `<h2>What's New?</h2>`,
			wantID: "what-s-new",
		},
		{
			name:   "with numbers",
			input:  `<h2>Version 2.0</h2>`,
			wantID: "version-2-0",
		},
		{
			name:   "with nested HTML",
			input:  `<h2><code>func</code> Example</h2>`,
			wantID: "func-example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewHeadingAnchorsPlugin()
			post := &models.Post{ArticleHTML: tt.input}

			err := p.processPost(post)
			if err != nil {
				t.Fatalf("processPost error: %v", err)
			}

			expected := `id="` + tt.wantID + `"`
			if !strings.Contains(post.ArticleHTML, expected) {
				t.Errorf("expected %q in output, got %q", expected, post.ArticleHTML)
			}
		})
	}
}

func TestHeadingAnchorsPlugin_DuplicateIDs(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	post := &models.Post{
		ArticleHTML: `<h2>Section</h2>
<h2>Section</h2>
<h2>Section</h2>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should have unique IDs
	if !strings.Contains(post.ArticleHTML, `id="section"`) {
		t.Error("expected first heading to have id='section'")
	}
	if !strings.Contains(post.ArticleHTML, `id="section-1"`) {
		t.Error("expected second heading to have id='section-1'")
	}
	if !strings.Contains(post.ArticleHTML, `id="section-2"`) {
		t.Error("expected third heading to have id='section-2'")
	}

	// Each anchor should link to its own ID
	if !strings.Contains(post.ArticleHTML, `href="#section"`) {
		t.Error("expected anchor linking to #section")
	}
	if !strings.Contains(post.ArticleHTML, `href="#section-1"`) {
		t.Error("expected anchor linking to #section-1")
	}
	if !strings.Contains(post.ArticleHTML, `href="#section-2"`) {
		t.Error("expected anchor linking to #section-2")
	}
}

func TestHeadingAnchorsPlugin_PreserveExistingID(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	post := &models.Post{
		ArticleHTML: `<h2 id="custom-id">My Heading</h2>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should preserve the existing ID
	if !strings.Contains(post.ArticleHTML, `id="custom-id"`) {
		t.Errorf("expected existing id preserved, got %q", post.ArticleHTML)
	}

	// Anchor should link to existing ID
	if !strings.Contains(post.ArticleHTML, `href="#custom-id"`) {
		t.Errorf("expected anchor to link to #custom-id, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_CustomSymbolAndClass(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	p.SetSymbol("link")
	p.SetClass("my-anchor")

	post := &models.Post{
		ArticleHTML: `<h2 id="test">Test</h2>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="my-anchor"`) {
		t.Errorf("expected custom class, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `>link</a>`) {
		t.Errorf("expected custom symbol, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_SkipPost(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	post := &models.Post{
		ArticleHTML: `<h2 id="test">Test</h2>`,
		Skip:        true,
	}

	originalHTML := post.ArticleHTML
	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if post.ArticleHTML != originalHTML {
		t.Errorf("skipped post should not be modified, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_EmptyArticleHTML(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	post := &models.Post{
		ArticleHTML: "",
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if post.ArticleHTML != "" {
		t.Errorf("empty ArticleHTML should remain empty, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_Disabled(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	p.SetEnabled(false)

	m := lifecycle.NewManager()
	posts := []*models.Post{
		{ArticleHTML: `<h2 id="test">Test</h2>`},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	// Should not modify posts when disabled
	result := m.Posts()
	if strings.Contains(result[0].ArticleHTML, "heading-anchor") {
		t.Error("disabled plugin should not add anchors")
	}
}

func TestHeadingAnchorsPlugin_Render(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	m := lifecycle.NewManager()

	posts := []*models.Post{
		{ArticleHTML: `<h2 id="section-1">Section 1</h2>`},
		{ArticleHTML: `<h2 id="section-2">Section 2</h2>`},
		{ArticleHTML: `<h2 id="section-3">Section 3</h2>`, Skip: true},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	resultPosts := m.Posts()

	// First two should have anchors
	if !strings.Contains(resultPosts[0].ArticleHTML, "heading-anchor") {
		t.Errorf("Post 1 should have anchor: %q", resultPosts[0].ArticleHTML)
	}
	if !strings.Contains(resultPosts[1].ArticleHTML, "heading-anchor") {
		t.Errorf("Post 2 should have anchor: %q", resultPosts[1].ArticleHTML)
	}

	// Skipped post should not have anchor
	if strings.Contains(resultPosts[2].ArticleHTML, "heading-anchor") {
		t.Errorf("Skipped post should not have anchor: %q", resultPosts[2].ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_MultipleHeadings(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	p.SetLevelRange(1, 6)

	post := &models.Post{
		ArticleHTML: `<h1 id="title">Title</h1>
<p>Some content</p>
<h2 id="intro">Introduction</h2>
<p>More content</p>
<h3 id="details">Details</h3>
<p>Even more content</p>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// All headings should have anchors
	if strings.Count(post.ArticleHTML, "heading-anchor") != 3 {
		t.Errorf("expected 3 anchors, got %d in: %q", strings.Count(post.ArticleHTML, "heading-anchor"), post.ArticleHTML)
	}

	// Each should link correctly
	if !strings.Contains(post.ArticleHTML, `href="#title"`) {
		t.Error("expected anchor for #title")
	}
	if !strings.Contains(post.ArticleHTML, `href="#intro"`) {
		t.Error("expected anchor for #intro")
	}
	if !strings.Contains(post.ArticleHTML, `href="#details"`) {
		t.Error("expected anchor for #details")
	}
}

func TestHeadingAnchorsPlugin_HeadingWithAttributes(t *testing.T) {
	p := NewHeadingAnchorsPlugin()

	post := &models.Post{
		ArticleHTML: `<h2 id="test" class="important" data-custom="value">Test Heading</h2>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should preserve other attributes
	if !strings.Contains(post.ArticleHTML, `class="important"`) {
		t.Errorf("expected class attribute preserved, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `data-custom="value"`) {
		t.Errorf("expected data attribute preserved, got %q", post.ArticleHTML)
	}

	// Should add anchor
	if !strings.Contains(post.ArticleHTML, `href="#test"`) {
		t.Errorf("expected anchor link, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_CaseInsensitiveMatching(t *testing.T) {
	p := NewHeadingAnchorsPlugin()
	p.SetLevelRange(1, 6)

	// Test uppercase heading tags
	post := &models.Post{
		ArticleHTML: `<H2 id="test">Test</H2>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should still add anchor
	if !strings.Contains(post.ArticleHTML, "heading-anchor") {
		t.Errorf("expected anchor for uppercase heading, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_SkipsExistingGoldmarkAnchors(t *testing.T) {
	p := NewHeadingAnchorsPlugin()

	post := &models.Post{
		ArticleHTML: `<h2 id="my-heading">My Heading <a class="anchor" href="#my-heading">¶</a></h2>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Should NOT add a second anchor
	if strings.Contains(post.ArticleHTML, `class="heading-anchor"`) {
		t.Errorf("should not add heading-anchor to headings with existing goldmark anchor, got %q", post.ArticleHTML)
	}
}

func TestHeadingAnchorsPlugin_SkipsExistingGoldmarkAnchors_MultipleHeadings(t *testing.T) {
	p := NewHeadingAnchorsPlugin()

	post := &models.Post{
		ArticleHTML: `<h2 id="first">First <a class="anchor" href="#first">¶</a></h2>
<h2 id="second">Second</h2>
<h3 id="third">Third <a class="anchor" href="#third">¶</a></h3>`,
	}

	err := p.processPost(post)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	// Headings with goldmark anchors should be untouched
	if strings.Contains(post.ArticleHTML, `First <a class="anchor" href="#first">¶</a> <a href="#first" class="heading-anchor">`) {
		t.Error("should not add heading-anchor to first heading with existing goldmark anchor")
	}

	// Heading without goldmark anchor should get a heading-anchor
	if !strings.Contains(post.ArticleHTML, `Second <a href="#second" class="heading-anchor">#</a>`) {
		t.Errorf("second heading should have heading-anchor added, got %q", post.ArticleHTML)
	}

	// Third heading (h3) with goldmark anchor should be untouched
	if strings.Contains(post.ArticleHTML, `Third <a class="anchor" href="#third">¶</a> <a href="#third" class="heading-anchor">`) {
		t.Error("should not add heading-anchor to third heading with existing goldmark anchor")
	}
}

// Interface compliance tests

func TestHeadingAnchorsPlugin_Interfaces(_ *testing.T) {
	p := NewHeadingAnchorsPlugin()

	// Verify interface compliance
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.RenderPlugin = p
	var _ lifecycle.PriorityPlugin = p
}
