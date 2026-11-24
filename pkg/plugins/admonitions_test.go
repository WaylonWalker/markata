package plugins

import (
	"bytes"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// =============================================================================
// Admonition AST Node Tests
// =============================================================================

func TestAdmonitionNode_NewAdmonition(t *testing.T) {
	ad := NewAdmonition("note", "Important", false, false, "")

	if ad.AdmonitionType != "note" {
		t.Errorf("expected type 'note', got %q", ad.AdmonitionType)
	}
	if ad.AdmonitionTitle != "Important" {
		t.Errorf("expected title 'Important', got %q", ad.AdmonitionTitle)
	}
	if ad.Collapsible {
		t.Error("expected Collapsible to be false")
	}
	if ad.DefaultOpen {
		t.Error("expected DefaultOpen to be false")
	}
}

func TestAdmonitionNode_NewAdmonitionCollapsible(t *testing.T) {
	ad := NewAdmonition("note", "Title", true, true, "")

	if !ad.Collapsible {
		t.Error("expected Collapsible to be true")
	}
	if !ad.DefaultOpen {
		t.Error("expected DefaultOpen to be true")
	}
}

func TestAdmonitionNode_NewAdmonitionAside(t *testing.T) {
	ad := NewAdmonition("aside", "Definition", false, false, "left")

	if ad.AdmonitionType != "aside" {
		t.Errorf("expected type 'aside', got %q", ad.AdmonitionType)
	}
	if ad.Position != "left" {
		t.Errorf("expected Position 'left', got %q", ad.Position)
	}
}

func TestAdmonitionNode_Kind(t *testing.T) {
	ad := NewAdmonition("note", "Test", false, false, "")

	if ad.Kind() != KindAdmonition {
		t.Errorf("expected KindAdmonition, got %v", ad.Kind())
	}
}

// =============================================================================
// AdmonitionParser Tests
// =============================================================================

func TestAdmonitionParser_Trigger(t *testing.T) {
	p := NewAdmonitionParser()
	triggers := p.Trigger()

	// Should trigger on both ! and ?
	if len(triggers) != 2 {
		t.Errorf("expected 2 triggers, got %d", len(triggers))
	}
	hasExclaim := false
	hasQuestion := false
	for _, tr := range triggers {
		if tr == '!' {
			hasExclaim = true
		}
		if tr == '?' {
			hasQuestion = true
		}
	}
	if !hasExclaim {
		t.Error("expected '!' trigger")
	}
	if !hasQuestion {
		t.Error("expected '?' trigger")
	}
}

func TestAdmonitionParser_CanInterruptParagraph(t *testing.T) {
	p := NewAdmonitionParser()
	if !p.CanInterruptParagraph() {
		t.Error("expected CanInterruptParagraph to return true")
	}
}

func TestAdmonitionParser_CanAcceptIndentedLine(t *testing.T) {
	p := NewAdmonitionParser()
	if p.CanAcceptIndentedLine() {
		t.Error("expected CanAcceptIndentedLine to return false")
	}
}

// =============================================================================
// Admonition Rendering Tests
// =============================================================================

func renderAdmonitionMarkdown(input string) string {
	md := goldmark.New(
		goldmark.WithExtensions(&AdmonitionExtension{}),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		return ""
	}
	return buf.String()
}

func TestAdmonitionRender_NoteWithTitle(t *testing.T) {
	// Test case from tests.yaml: "admonition note"
	input := `!!! note "Important"
    This is important information.`

	output := renderAdmonitionMarkdown(input)

	// Check for admonition structure
	if !strings.Contains(output, `class="admonition note"`) {
		t.Errorf("expected admonition note class in output, got %q", output)
	}
	if !strings.Contains(output, `class="admonition-title"`) {
		t.Errorf("expected admonition-title class in output, got %q", output)
	}
	if !strings.Contains(output, "Important") {
		t.Errorf("expected title 'Important' in output, got %q", output)
	}
}

func TestAdmonitionRender_WarningType(t *testing.T) {
	// Test case from tests.yaml: "admonition warning"
	input := `!!! warning
    Be careful here.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, `class="admonition warning"`) {
		t.Errorf("expected admonition warning class in output, got %q", output)
	}
}

func TestAdmonitionRender_DefaultTitleFromType(t *testing.T) {
	// When no title is specified, use capitalized type
	input := `!!! warning
    Content here.`

	output := renderAdmonitionMarkdown(input)

	// Should have "Warning" as default title
	if !strings.Contains(output, "Warning") {
		t.Errorf("expected default title 'Warning' in output, got %q", output)
	}
}

func TestAdmonitionRender_AllSupportedTypes(t *testing.T) {
	// Test all supported admonition types
	types := []string{"note", "info", "tip", "hint", "success", "warning", "caution", "important", "danger", "error", "bug", "example", "quote", "abstract"}

	for _, adType := range types {
		t.Run(adType, func(t *testing.T) {
			input := "!!! " + adType + "\n    Content here."
			output := renderAdmonitionMarkdown(input)

			expected := `class="admonition ` + adType + `"`
			if !strings.Contains(output, expected) {
				t.Errorf("expected %q in output, got %q", expected, output)
			}
		})
	}
}

func TestAdmonitionRender_CustomTypes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "custom type uses capitalized default title",
			input: "!!! scripture\n    Verse text.",
			want: []string{
				"<div class=\"admonition scripture\">",
				"<p class=\"admonition-title\">Scripture</p>",
				"<p>Verse text.</p>",
				"</div>",
			},
		},
		{
			name:  "hyphenated custom type",
			input: "!!! background-thought\n    A passing thought.",
			want: []string{
				"<div class=\"admonition background-thought\">",
				"<p class=\"admonition-title\">Background-thought</p>",
				"<p>A passing thought.</p>",
			},
		},
		{
			name:  "custom type with quoted title",
			input: "!!! source \"Primary Source\"\n    Source text.",
			want: []string{
				"<div class=\"admonition source\">",
				"<p class=\"admonition-title\">Primary Source</p>",
				"<p>Source text.</p>",
			},
		},
		{
			name:  "custom type with unquoted title",
			input: "!!! scripture Ancient Text\n    Passage.",
			want: []string{
				"<p class=\"admonition-title\">Ancient Text</p>",
				"<p>Passage.</p>",
			},
		},
		{
			name:  "collapsed custom type",
			input: "??? scripture\n    Hidden verse.",
			want: []string{
				"<details class=\"admonition scripture\">",
				"<summary class=\"admonition-title\">Scripture</summary>",
				"<p>Hidden verse.</p>",
				"</details>",
			},
		},
		{
			name:  "expanded custom type with title",
			input: "???+ source \"Further reading\"\n    Visible source.",
			want: []string{
				"<details class=\"admonition source\" open>",
				"<summary class=\"admonition-title\">Further reading</summary>",
				"<p>Visible source.</p>",
				"</details>",
			},
		},
		{
			name:  "custom type is case insensitive",
			input: "!!! ScRiPtUrE\n    Verse.",
			want: []string{
				"<div class=\"admonition scripture\">",
				"<p class=\"admonition-title\">Scripture</p>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderAdmonitionMarkdown(tt.input)
			for _, want := range tt.want {
				if !strings.Contains(output, want) {
					t.Errorf("expected %q in output, got %q", want, output)
				}
			}
		})
	}
}

func TestAdmonitionRender_CustomTitleWithQuotes(t *testing.T) {
	input := `!!! note "Custom Title Here"
    Some content.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "Custom Title Here") {
		t.Errorf("expected custom title in output, got %q", output)
	}
}

// =============================================================================
// Collapsible Admonition Tests
// =============================================================================

func TestAdmonitionRender_CollapsibleClosed(t *testing.T) {
	input := `??? note "Collapsed by default"
    Hidden content.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "<details") {
		t.Errorf("expected <details> element for collapsible, got %q", output)
	}
	if !strings.Contains(output, `class="admonition note"`) {
		t.Errorf("expected admonition class, got %q", output)
	}
	if !strings.Contains(output, "<summary") {
		t.Errorf("expected <summary> element, got %q", output)
	}
	// Should NOT have open attribute
	if strings.Contains(output, "open") {
		t.Errorf("collapsed admonition should not have open attribute, got %q", output)
	}
}

func TestAdmonitionRender_CollapsibleOpen(t *testing.T) {
	input := `???+ note "Expanded by default"
    Visible content.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "<details") {
		t.Errorf("expected <details> element for collapsible, got %q", output)
	}
	if !strings.Contains(output, " open") {
		t.Errorf("expanded admonition should have open attribute, got %q", output)
	}
	if !strings.Contains(output, "</details>") {
		t.Errorf("expected closing </details>, got %q", output)
	}
}

func TestAdmonitionRender_CollapsibleAllTypes(t *testing.T) {
	types := []string{"note", "warning", "tip", "danger"}

	for _, adType := range types {
		t.Run("collapsed_"+adType, func(t *testing.T) {
			input := "??? " + adType + "\n    Content."
			output := renderAdmonitionMarkdown(input)

			if !strings.Contains(output, "<details") {
				t.Errorf("expected <details> for ??? %s", adType)
			}
			if !strings.Contains(output, `class="admonition `+adType) {
				t.Errorf("expected admonition class for %s", adType)
			}
		})

		t.Run("expanded_"+adType, func(t *testing.T) {
			input := "???+ " + adType + "\n    Content."
			output := renderAdmonitionMarkdown(input)

			if !strings.Contains(output, "<details") {
				t.Errorf("expected <details> for ???+ %s", adType)
			}
			if !strings.Contains(output, " open") {
				t.Errorf("expected open attribute for ???+ %s", adType)
			}
		})
	}
}

// =============================================================================
// Aside Admonition Tests
// =============================================================================

func TestAdmonitionRender_AsideBasic(t *testing.T) {
	input := `!!! aside
    This is a marginal note.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "<aside") {
		t.Errorf("expected <aside> element, got %q", output)
	}
	// Default position is right
	if !strings.Contains(output, "aside-right") {
		t.Errorf("expected aside-right class (default), got %q", output)
	}
	if !strings.Contains(output, "</aside>") {
		t.Errorf("expected closing </aside>, got %q", output)
	}
}

func TestAdmonitionRender_AsideWithTitle(t *testing.T) {
	input := `!!! aside "Definition"
    A static site generator converts source files into static HTML.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "<aside") {
		t.Errorf("expected <aside> element, got %q", output)
	}
	if !strings.Contains(output, "Definition") {
		t.Errorf("expected title 'Definition', got %q", output)
	}
}

func TestAdmonitionRender_AsideNoDefaultTitle(t *testing.T) {
	// Aside has no default title per spec
	input := `!!! aside
    Content without title.`

	output := renderAdmonitionMarkdown(input)

	// Should NOT have admonition-title with "Aside" text
	if strings.Contains(output, ">Aside<") {
		t.Errorf("aside should not have default title 'Aside', got %q", output)
	}
}

func TestAdmonitionRender_AsideLeft(t *testing.T) {
	input := `!!! aside left
    Positioned on the left.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "aside-left") {
		t.Errorf("expected aside-left class, got %q", output)
	}
}

func TestAdmonitionRender_AsideRight(t *testing.T) {
	input := `!!! aside right
    Positioned on the right.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "aside-right") {
		t.Errorf("expected aside-right class, got %q", output)
	}
}

func TestAdmonitionRender_AsideInline(t *testing.T) {
	// 'inline' is an alias for 'left' (Material for MkDocs compat)
	input := `!!! aside inline
    Floats to the left.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "aside-left") {
		t.Errorf("expected aside-left class (inline alias), got %q", output)
	}
}

func TestAdmonitionRender_AsideInlineEnd(t *testing.T) {
	// 'inline end' is an alias for 'right' (Material for MkDocs compat)
	input := `!!! aside inline end
    Floats to the right.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "aside-right") {
		t.Errorf("expected aside-right class (inline end alias), got %q", output)
	}
}

func TestAdmonitionRender_AsideInlineWithTitle(t *testing.T) {
	input := `!!! aside inline "Side Note"
    This is a side note.`

	output := renderAdmonitionMarkdown(input)

	if !strings.Contains(output, "aside-left") {
		t.Errorf("expected aside-left class, got %q", output)
	}
	if !strings.Contains(output, "Side Note") {
		t.Errorf("expected title 'Side Note', got %q", output)
	}
}

// =============================================================================
// Regex Tests
// =============================================================================

func TestAdmonitionRegex_Matching(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantMarker string
		wantType   string
		wantMod    string
		wantTitle  string
		wantMatch  bool
	}{
		{
			name:       "basic note",
			input:      "!!! note",
			wantMarker: "!!!",
			wantType:   "note",
			wantMod:    "",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "note with title",
			input:      `!!! note "My Title"`,
			wantMarker: "!!!",
			wantType:   "note",
			wantMod:    "",
			wantTitle:  "My Title",
			wantMatch:  true,
		},
		{
			name:       "collapsible closed",
			input:      "??? note",
			wantMarker: "???",
			wantType:   "note",
			wantMod:    "",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "collapsible open",
			input:      "???+ note",
			wantMarker: "???+",
			wantType:   "note",
			wantMod:    "",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "collapsible with title",
			input:      `??? warning "Click to expand"`,
			wantMarker: "???",
			wantType:   "warning",
			wantMod:    "",
			wantTitle:  "Click to expand",
			wantMatch:  true,
		},
		{
			name:       "aside basic",
			input:      "!!! aside",
			wantMarker: "!!!",
			wantType:   "aside",
			wantMod:    "",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "aside left",
			input:      "!!! aside left",
			wantMarker: "!!!",
			wantType:   "aside",
			wantMod:    "left",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "aside right",
			input:      "!!! aside right",
			wantMarker: "!!!",
			wantType:   "aside",
			wantMod:    "right",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "aside inline",
			input:      "!!! aside inline",
			wantMarker: "!!!",
			wantType:   "aside",
			wantMod:    "inline",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "aside inline end",
			input:      "!!! aside inline end",
			wantMarker: "!!!",
			wantType:   "aside",
			wantMod:    "inline end",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "aside left with title",
			input:      `!!! aside left "Note"`,
			wantMarker: "!!!",
			wantType:   "aside",
			wantMod:    "left",
			wantTitle:  "Note",
			wantMatch:  true,
		},
		{
			name:       "warning type",
			input:      "!!! warning",
			wantMarker: "!!!",
			wantType:   "warning",
			wantMod:    "",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:       "tip with spaces before title",
			input:      `!!! tip   "Helpful"`,
			wantMarker: "!!!",
			wantType:   "tip",
			wantMod:    "",
			wantTitle:  "Helpful",
			wantMatch:  true,
		},
		{
			name:       "empty title",
			input:      `!!! note ""`,
			wantMarker: "!!!",
			wantType:   "note",
			wantMod:    "",
			wantTitle:  "",
			wantMatch:  true,
		},
		{
			name:      "not enough exclamation marks",
			input:     "!! note",
			wantMatch: false,
		},
		{
			name:      "missing space after type",
			input:     `!!!note"Title"`,
			wantMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := admonitionRegex.FindStringSubmatch(tt.input)
			gotMatch := matches != nil

			if gotMatch != tt.wantMatch {
				t.Errorf("regex match = %v, want %v for input %q", gotMatch, tt.wantMatch, tt.input)
				return
			}

			if !tt.wantMatch {
				return
			}

			if len(matches) < 3 {
				t.Fatalf("expected at least 3 groups in match, got %d: %v", len(matches), matches)
			}

			if matches[1] != tt.wantMarker {
				t.Errorf("marker = %q, want %q", matches[1], tt.wantMarker)
			}

			if matches[2] != tt.wantType {
				t.Errorf("type = %q, want %q", matches[2], tt.wantType)
			}

			if len(matches) >= 4 && matches[3] != tt.wantMod {
				t.Errorf("modifier = %q, want %q", matches[3], tt.wantMod)
			}

			if len(matches) >= 5 && matches[4] != tt.wantTitle {
				t.Errorf("title = %q, want %q", matches[4], tt.wantTitle)
			}
		})
	}
}

func TestAdmonitionRender_StructureWithDiv(t *testing.T) {
	input := `!!! note "Test Title"
    Test content.`

	output := renderAdmonitionMarkdown(input)

	// Should have proper HTML structure
	if !strings.Contains(output, "<div") {
		t.Errorf("expected <div> element, got %q", output)
	}
	if !strings.Contains(output, "</div>") {
		t.Errorf("expected closing </div> tag, got %q", output)
	}
	if !strings.Contains(output, "<p") {
		t.Errorf("expected <p> element for title, got %q", output)
	}
}

func TestAdmonitionParser_Open(t *testing.T) {
	p := NewAdmonitionParser()

	tests := []struct {
		name     string
		line     string
		wantOpen bool
	}{
		{"valid note", "!!! note\n", true},
		{"valid note with title", `!!! note "Title"` + "\n", true},
		{"valid warning", "!!! warning\n", true},
		{"collapsible closed", "??? note\n", true},
		{"collapsible open", "???+ note\n", true},
		{"aside", "!!! aside\n", true},
		{"aside left", "!!! aside left\n", true},
		{"aside right", "!!! aside right\n", true},
		{"aside inline", "!!! aside inline\n", true},
		{"custom type", "!!! invalid\n", true},
		{"invalid marker", "!! invalid\n", false},
		{"not admonition", "Not an admonition\n", false},
		{"only exclamation", "!!!\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := text.NewReader([]byte(tt.line))
			node, _ := p.Open(nil, reader, parser.NewContext())
			gotOpen := node != nil

			if gotOpen != tt.wantOpen {
				t.Errorf("Open() returned node = %v, want %v", gotOpen, tt.wantOpen)
			}

			if gotOpen {
				ad, ok := node.(*Admonition)
				if !ok {
					t.Errorf("expected *Admonition node, got %T", node)
				} else if ad.AdmonitionType == "" {
					t.Error("expected AdmonitionType to be set")
				}
			}
		})
	}
}

func TestAdmonitionExtension_Extend(t *testing.T) {
	ext := &AdmonitionExtension{}
	md := goldmark.New()

	// Should not panic
	ext.Extend(md)

	// Verify it can render admonitions after extension
	input := `!!! note
    Content`

	var buf bytes.Buffer
	err := md.Convert([]byte(input), &buf)
	if err != nil {
		t.Fatalf("conversion error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "admonition") {
		t.Errorf("expected admonition after extension, got %q", output)
	}
}

func TestAdmonitionRenderer_RegisterFuncs(_ *testing.T) {
	r := NewAdmonitionRenderer()

	// Just verify it doesn't panic
	// The actual registration is tested through the full rendering tests
	_ = r
}

func TestAdmonitionNode_Dump(_ *testing.T) {
	ad := NewAdmonition("note", "Test Title", false, false, "")

	// Should not panic
	ad.Dump([]byte("source"), 0)
}

// =============================================================================
// Integration with RenderMarkdownPlugin
// =============================================================================

func TestAdmonitionIntegration_AllTypes(t *testing.T) {
	p := NewRenderMarkdownPlugin()

	tests := []struct {
		name      string
		content   string
		wantClass string
		wantTitle string
	}{
		{
			name:      "note with title",
			content:   "!!! note \"Important\"\n    This is important.",
			wantClass: `class="admonition note"`,
			wantTitle: "Important",
		},
		{
			name:      "warning default title",
			content:   "!!! warning\n    Be careful.",
			wantClass: `class="admonition warning"`,
			wantTitle: "Warning",
		},
		{
			name:      "tip type",
			content:   "!!! tip\n    Here's a tip.",
			wantClass: `class="admonition tip"`,
			wantTitle: "Tip",
		},
		{
			name:      "danger type",
			content:   "!!! danger\n    Danger zone!",
			wantClass: `class="admonition danger"`,
			wantTitle: "Danger",
		},
		{
			name:      "important type",
			content:   "!!! important\n    This is important.",
			wantClass: `class="admonition important"`,
			wantTitle: "Important",
		},
		{
			name:      "caution type",
			content:   "!!! caution\n    Use caution.",
			wantClass: `class="admonition caution"`,
			wantTitle: "Caution",
		},
		{
			name:      "info type",
			content:   "!!! info\n    Information.",
			wantClass: `class="admonition info"`,
			wantTitle: "Info",
		},
		{
			name:      "success type",
			content:   "!!! success\n    Success!",
			wantClass: `class="admonition success"`,
			wantTitle: "Success",
		},
		{
			name:      "bug type",
			content:   "!!! bug\n    Bug report.",
			wantClass: `class="admonition bug"`,
			wantTitle: "Bug",
		},
		{
			name:      "example type",
			content:   "!!! example\n    Example content.",
			wantClass: `class="admonition example"`,
			wantTitle: "Example",
		},
		{
			name:      "quote type",
			content:   "!!! quote\n    A quote.",
			wantClass: `class="admonition quote"`,
			wantTitle: "Quote",
		},
		{
			name:      "abstract type",
			content:   "!!! abstract\n    Abstract text.",
			wantClass: `class="admonition abstract"`,
			wantTitle: "Abstract",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Content: tt.content}

			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}

			if !strings.Contains(post.ArticleHTML, tt.wantClass) {
				t.Errorf("expected %q in output, got %q", tt.wantClass, post.ArticleHTML)
			}
			if !strings.Contains(post.ArticleHTML, tt.wantTitle) {
				t.Errorf("expected title %q in output, got %q", tt.wantTitle, post.ArticleHTML)
			}
		})
	}
}

func TestAdmonitionIntegration_Collapsible(t *testing.T) {
	p := NewRenderMarkdownPlugin()

	tests := []struct {
		name       string
		content    string
		wantOpen   bool
		wantDetail bool
	}{
		{
			name:       "collapsed",
			content:    "??? note \"Hidden\"\n    Hidden content.",
			wantOpen:   false,
			wantDetail: true,
		},
		{
			name:       "expanded",
			content:    "???+ warning \"Visible\"\n    Visible content.",
			wantOpen:   true,
			wantDetail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Content: tt.content}

			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}

			if tt.wantDetail && !strings.Contains(post.ArticleHTML, "<details") {
				t.Errorf("expected <details> element, got %q", post.ArticleHTML)
			}
			if tt.wantOpen && !strings.Contains(post.ArticleHTML, " open") {
				t.Errorf("expected open attribute, got %q", post.ArticleHTML)
			}
			if !tt.wantOpen && strings.Contains(post.ArticleHTML, " open") {
				t.Errorf("did not expect open attribute, got %q", post.ArticleHTML)
			}
		})
	}
}

func TestAdmonitionIntegration_Aside(t *testing.T) {
	p := NewRenderMarkdownPlugin()

	tests := []struct {
		name      string
		content   string
		wantClass string
		wantAside bool
	}{
		{
			name:      "basic aside (default right)",
			content:   "!!! aside\n    Marginal note.",
			wantClass: "aside-right",
			wantAside: true,
		},
		{
			name:      "aside left",
			content:   "!!! aside left\n    Left note.",
			wantClass: "aside-left",
			wantAside: true,
		},
		{
			name:      "aside right",
			content:   "!!! aside right\n    Right note.",
			wantClass: "aside-right",
			wantAside: true,
		},
		{
			name:      "aside inline (alias for left)",
			content:   "!!! aside inline\n    Left note.",
			wantClass: "aside-left",
			wantAside: true,
		},
		{
			name:      "aside inline end (alias for right)",
			content:   "!!! aside inline end\n    Right note.",
			wantClass: "aside-right",
			wantAside: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Content: tt.content}

			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}

			if tt.wantAside && !strings.Contains(post.ArticleHTML, "<aside") {
				t.Errorf("expected <aside> element, got %q", post.ArticleHTML)
			}
			if !strings.Contains(post.ArticleHTML, tt.wantClass) {
				t.Errorf("expected %q in output, got %q", tt.wantClass, post.ArticleHTML)
			}
		})
	}
}

func TestAdmonitionIntegration_CustomTypeRendered(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "!!! scripture\n    Content here"}

	if err := p.renderPost(post); err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	for _, want := range []string{
		"<div class=\"admonition scripture\">",
		"<p class=\"admonition-title\">Scripture</p>",
		"<p>Content here</p>",
	} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("expected %q in output, got %q", want, post.ArticleHTML)
		}
	}
}

// =============================================================================
// Content Leak Regression Tests (issue #745)
//
// Admonitions with multiple indented paragraphs must close properly when
// followed by unindented content. Previously, the blank-line look-ahead in
// Continue() advanced the reader past blank lines and returned Close, causing
// goldmark's paragraph parser to merge unindented text into the admonition.
// =============================================================================

func TestAdmonitionRender_ContentLeakRegression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantInside  []string // content that should appear before the closing tag
		wantOutside []string // content that should appear after the closing tag
		closeTag    string   // the expected closing tag
	}{
		{
			name:        "single paragraph then outside",
			input:       "!!! note\n    Inside.\n\nOutside.\n",
			wantInside:  []string{"Inside."},
			wantOutside: []string{"Outside."},
			closeTag:    "</div>",
		},
		{
			name:        "multi paragraph then outside",
			input:       "!!! note\n    Para 1.\n\n    Para 2.\n\nOutside.\n",
			wantInside:  []string{"Para 1.", "Para 2."},
			wantOutside: []string{"Outside."},
			closeTag:    "</div>",
		},
		{
			name:        "three paragraphs then outside",
			input:       "!!! note\n    Para 1.\n\n    Para 2.\n\n    Para 3.\n\nOutside.\n",
			wantInside:  []string{"Para 1.", "Para 2.", "Para 3."},
			wantOutside: []string{"Outside."},
			closeTag:    "</div>",
		},
		{
			name:        "multi paragraph then multiple outside",
			input:       "!!! note\n    In 1.\n\n    In 2.\n\nOut 1.\n\nOut 2.\n",
			wantInside:  []string{"In 1.", "In 2."},
			wantOutside: []string{"Out 1.", "Out 2."},
			closeTag:    "</div>",
		},
		{
			name:        "multi paragraph then heading",
			input:       "!!! note\n    In 1.\n\n    In 2.\n\n## Heading\n\nAfter.\n",
			wantInside:  []string{"In 1.", "In 2."},
			wantOutside: []string{"Heading", "After."},
			closeTag:    "</div>",
		},
		{
			name:        "chat multi paragraph then outside",
			input:       "!!! chat @user\n    Msg 1.\n\n    Msg 2.\n\nOutside.\n",
			wantInside:  []string{"Msg 1.", "Msg 2."},
			wantOutside: []string{"Outside."},
			closeTag:    "</div>",
		},
		{
			name:        "chat-reply multi paragraph then outside",
			input:       "!!! chat-reply @other\n    Reply 1.\n\n    Reply 2.\n\nOutside.\n",
			wantInside:  []string{"Reply 1.", "Reply 2."},
			wantOutside: []string{"Outside."},
			closeTag:    "</div>",
		},
		{
			name:        "collapsible multi paragraph then outside",
			input:       "??? note\n    In 1.\n\n    In 2.\n\nOutside.\n",
			wantInside:  []string{"In 1.", "In 2."},
			wantOutside: []string{"Outside."},
			closeTag:    "</details>",
		},
		{
			name:        "aside multi paragraph then outside",
			input:       "!!! aside\n    In 1.\n\n    In 2.\n\nOutside.\n",
			wantInside:  []string{"In 1.", "In 2."},
			wantOutside: []string{"Outside."},
			closeTag:    "</aside>",
		},
		{
			name:        "double blank line then outside",
			input:       "!!! note\n    In 1.\n\n    In 2.\n\n\nOutside.\n",
			wantInside:  []string{"In 1.", "In 2."},
			wantOutside: []string{"Outside."},
			closeTag:    "</div>",
		},
		{
			name:        "blank line after header",
			input:       "!!! note\n\n    Para 1.\n\n    Para 2.\n\nOutside.\n",
			wantInside:  []string{"Para 1.", "Para 2."},
			wantOutside: []string{"Outside."},
			closeTag:    "</div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := renderAdmonitionMarkdown(tt.input)

			closeIdx := strings.Index(output, tt.closeTag)
			if closeIdx < 0 {
				t.Fatalf("no %s found in output:\n%s", tt.closeTag, output)
			}

			insidePart := output[:closeIdx]
			outsidePart := output[closeIdx+len(tt.closeTag):]

			for _, want := range tt.wantInside {
				if !strings.Contains(insidePart, want) {
					t.Errorf("expected %q INSIDE admonition, not found.\nInside:\n%s\nFull:\n%s", want, insidePart, output)
				}
			}

			for _, want := range tt.wantOutside {
				if !strings.Contains(outsidePart, want) {
					t.Errorf("expected %q OUTSIDE admonition (after %s), not found.\nOutside:\n%s\nFull:\n%s", want, tt.closeTag, outsidePart, output)
				}
			}
		})
	}
}

func TestAdmonitionIntegration_ContentLeakRegression(t *testing.T) {
	p := NewRenderMarkdownPlugin()

	tests := []struct {
		name        string
		content     string
		wantInside  []string
		wantOutside []string
	}{
		{
			name:        "multi paragraph then outside",
			content:     "!!! note\n    Para 1.\n\n    Para 2.\n\nOutside.\n",
			wantInside:  []string{"Para 1.", "Para 2."},
			wantOutside: []string{"Outside."},
		},
		{
			name:        "chat conversation then outside",
			content:     "!!! chat @alice\n    Hello!\n\n    How are you?\n\nThis is not part of the chat.\n",
			wantInside:  []string{"Hello!", "How are you?"},
			wantOutside: []string{"This is not part of the chat."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Content: tt.content}

			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}

			closeIdx := strings.Index(post.ArticleHTML, "</div>")
			if closeIdx < 0 {
				t.Fatalf("no </div> found in output:\n%s", post.ArticleHTML)
			}

			insidePart := post.ArticleHTML[:closeIdx]
			outsidePart := post.ArticleHTML[closeIdx+len("</div>"):]

			for _, want := range tt.wantInside {
				if !strings.Contains(insidePart, want) {
					t.Errorf("expected %q INSIDE admonition, not found.\nOutput:\n%s", want, post.ArticleHTML)
				}
			}

			for _, want := range tt.wantOutside {
				if !strings.Contains(outsidePart, want) {
					t.Errorf("expected %q OUTSIDE admonition, not found.\nOutput:\n%s", want, post.ArticleHTML)
				}
			}
		})
	}
}

// TestAdmonitionParser_FirstBodyLineIsNotCodeBlock guards against the first
// indented body line being parsed as an indented code block. The Open method
// must not advance onto the next line, otherwise goldmark tries child parsers
// before Continue strips the 4-space indent.
func TestAdmonitionParser_FirstBodyLineIsNotCodeBlock(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
		notWant string
	}{
		{
			name:    "paragraph",
			content: "!!! note \"T\"\n    First line.\n    Second line.\n",
			want:    "<p>First line.\nSecond line.</p>",
			notWant: "<pre>",
		},
		{
			name:    "bullet list",
			content: "!!! note\n    * one\n    * two\n",
			want:    "<ul>\n<li>one</li>\n<li>two</li>\n</ul>",
			notWant: "<pre>",
		},
		{
			name:    "heading",
			content: "!!! note\n    # Heading\n\nafter",
			want:    "<h1 id=\"heading\">Heading</h1>\n</div>\n<p>after</p>",
			notWant: "<pre>",
		},
		{
			name:    "blank line before body",
			content: "!!! note\n\n    para\n",
			want:    "<p>para</p>",
			notWant: "<pre>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRenderMarkdownPlugin()
			post := &models.Post{Content: tt.content}
			if err := p.renderPost(post); err != nil {
				t.Fatalf("renderPost() error = %v", err)
			}
			if !strings.Contains(post.ArticleHTML, tt.want) {
				t.Errorf("output missing %q:\n%s", tt.want, post.ArticleHTML)
			}
			if strings.Contains(post.ArticleHTML, tt.notWant) {
				t.Errorf("output contains %q:\n%s", tt.notWant, post.ArticleHTML)
			}
		})
	}
}
