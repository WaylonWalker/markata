package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestRenderMarkdownPlugin_Name(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	if p.Name() != "render_markdown" {
		t.Errorf("expected name 'render_markdown', got %q", p.Name())
	}
}

func TestRenderMarkdownPlugin_Configure(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	m := lifecycle.NewManager()
	if err := p.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}
}

func TestRenderMarkdownPlugin_BasicParagraph(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "Hello world"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	expected := "<p>Hello world</p>"
	if !strings.Contains(post.ArticleHTML, expected) {
		t.Errorf("expected %q in output, got %q", expected, post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_WrapsHeadingHighlights(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "## A ==wrapped highlight==\n"}
	if err := p.renderPost(post); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(post.ArticleHTML, `<span class="heading-highlight"><mark style="background-color:var(--color-highlight)!important;color:var(--color-highlight-text)!important">wrapped highlight</mark></span>`) {
		t.Fatalf("heading highlight wrapper missing: %s", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Headings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"h1", "# Heading 1", "<h1"},
		{"h2", "## Heading 2", "<h2"},
		{"h3", "### Heading 3", "<h3"},
		{"h4", "#### Heading 4", "<h4"},
		{"h5", "##### Heading 5", "<h5"},
		{"h6", "###### Heading 6", "<h6"},
	}

	p := NewRenderMarkdownPlugin()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Content: tt.input}
			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}
			if !strings.Contains(post.ArticleHTML, tt.expected) {
				t.Errorf("expected %q in output, got %q", tt.expected, post.ArticleHTML)
			}
		})
	}
}

func TestRenderMarkdownPlugin_ExpressiveH1AndH2(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "# **Strong** _emphasis_ ==mark== `code` [link](/test) ~~delete~~\n\n## ==**nested**=="}
	if err := p.renderPost(post); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"<strong>Strong</strong> <em>emphasis</em> <span class=\"heading-highlight\"><mark style=\"background-color:var(--color-highlight)!important;color:var(--color-highlight-text)!important\">mark</mark></span> <code>code</code> <a href=\"/test\">link</a> <del>delete</del>",
		"<span class=\"heading-highlight\"><mark style=\"background-color:var(--color-highlight)!important;color:var(--color-highlight-text)!important\"><strong>nested</strong></mark></span>",
	} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("rendered headings %q do not contain %q", post.ArticleHTML, want)
		}
	}
}

func TestRenderMarkdownPlugin_Emphasis(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "*italic* and **bold**"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<em>italic</em>") {
		t.Errorf("expected <em>italic</em> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<strong>bold</strong>") {
		t.Errorf("expected <strong>bold</strong> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_UnorderedList(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "- item 1\n- item 2\n- item 3"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<ul>") {
		t.Errorf("expected <ul> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<li>") {
		t.Errorf("expected <li> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_OrderedList(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "1. first\n2. second\n3. third"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<ol>") {
		t.Errorf("expected <ol> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<li>") {
		t.Errorf("expected <li> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_CodeBlock(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "```python\nprint('hello')\n```"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<pre") {
		t.Errorf("expected <pre> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<code") {
		t.Errorf("expected <code> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_InlineCode(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "Use `code` inline"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<code>code</code>") {
		t.Errorf("expected <code>code</code> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Links(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "[link text](https://example.com)"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<a href="https://example.com"`) {
		t.Errorf("expected link in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "link text") {
		t.Errorf("expected link text in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Images(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "![alt text](image.png)"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<img src="image.png"`) {
		t.Errorf("expected img tag in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `alt="alt text"`) {
		t.Errorf("expected alt text in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Blockquote(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> This is a quote"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<blockquote>") {
		t.Errorf("expected <blockquote> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_FigureCaptionFromFollowingBlockquote(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "![Alt text](image.png)\n> This is the caption."}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<figure>") {
		t.Fatalf("expected <figure> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<figcaption><p>This is the caption.</p></figcaption>") {
		t.Fatalf("expected blockquote caption to merge into figure, got %q", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "<blockquote>") {
		t.Fatalf("expected blockquote to be consumed as figcaption, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_StandaloneBlockquoteUnaffectedByFigureMerge(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "Paragraph before\n\n> This stays a blockquote."}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<blockquote>") {
		t.Fatalf("expected standalone blockquote to remain, got %q", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "<figcaption>") {
		t.Fatalf("expected no figcaption for standalone blockquote, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_ComparisonContainerBlockquoteRemainsOutsideFigure(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "::: wa-comparison\n![Before](before.webp)\n![After](after.webp)\n:::\n\n> comparison caption"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<div class="wa-comparison">`) {
		t.Fatalf("expected wa-comparison container, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<blockquote>") {
		t.Fatalf("expected blockquote to remain outside comparison figure, got %q", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, "<figcaption><p>comparison caption</p></figcaption>") {
		t.Fatalf("expected comparison blockquote not to merge into inner figure, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_SamplePostBothBlockquoteCaptionsMerge(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "::: wa-comparison\n![Before](before.webp)\n![](after.webp)\n:::\n\n![Before](before.webp)\n> black and white\n\n![](after.webp)\n> original"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if strings.Count(post.ArticleHTML, "<figcaption><p>black and white</p></figcaption>") != 1 {
		t.Fatalf("expected first blockquote caption to merge once, got %q", post.ArticleHTML)
	}
	if strings.Count(post.ArticleHTML, "<figcaption><p>original</p></figcaption>") != 1 {
		t.Fatalf("expected second blockquote caption to merge once, got %q", post.ArticleHTML)
	}
	if strings.Contains(post.ArticleHTML, `<div class="wa-comparison"><figure><img src="before.webp" alt="Before" /><img src="after.webp" alt="" /><figcaption>`) {
		t.Fatalf("expected comparison inner figure not to absorb following caption, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteAttributionPlainText(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\nsomeone"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>") {
		t.Fatalf("expected <footer> in output, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 someone</footer>") {
		t.Fatalf("expected footer with em dash and name, got %q", post.ArticleHTML)
	}

	// Attribution should be inside the blockquote
	if !strings.Contains(post.ArticleHTML, "</p>\n  <footer>") {
		t.Fatalf("expected footer inside blockquote after quote body, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteAttributionMention(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\n@waylonwalker"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 @waylonwalker</footer>") {
		t.Fatalf("expected footer with @mention, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteAttributionWikilink(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\n[[Wiki Link]]"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 [[Wiki Link]]</footer>") {
		t.Fatalf("expected footer with wikilink, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteAttributionLink(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\n[Author Name](/author)"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 <a href=\"/author\">Author Name</a></footer>") {
		t.Fatalf("expected footer with link, got %q", post.ArticleHTML)
	}

	// Single person link should NOT have cite wrapping
	if strings.Contains(post.ArticleHTML, "<cite>") {
		t.Fatalf("expected no <cite> for a single person link, got %q", post.ArticleHTML)
	}

	// No cite attribute on blockquote for person-only links
	if strings.Contains(post.ArticleHTML, `<blockquote cite="`) {
		t.Fatalf("expected no cite attr for person-only attribution, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteAttributionDash(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\n\u2014 Author Name"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	// User-provided dash should be kept; no additional dash added
	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 Author Name</footer>") {
		t.Fatalf("expected footer with single dash, got %q", post.ArticleHTML)
	}

	if strings.Count(post.ArticleHTML, "\u2014") != 1 {
		t.Fatalf("expected exactly one em dash (no double dash), got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteAttributionWithBlankLine(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\n\nsomeone"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if strings.Contains(post.ArticleHTML, "<footer>") {
		t.Fatalf("expected no footer for blockquote with blank line, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<blockquote>") {
		t.Fatalf("expected blockquote to remain, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<p>someone</p>") {
		t.Fatalf("expected someone to be a separate paragraph, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteAttributionMultiParagraph(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> Paragraph 1\n>\n> Paragraph 2\n\u2014 Author"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 Author</footer>") {
		t.Fatalf("expected footer after multi-paragraph quote body, got %q", post.ArticleHTML)
	}

	// Footer should be inside the blockquote, after the last </p>
	if !strings.Contains(post.ArticleHTML, "</p>\n  <footer>") {
		t.Fatalf("expected footer inside blockquote after quote body, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_BlockquoteLazyContinuationNotAttribution(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\nthis is a long continuation that should not be treated as attribution text"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if strings.Contains(post.ArticleHTML, "<footer>") {
		t.Fatalf("expected no footer for lazy continuation, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<blockquote>") {
		t.Fatalf("expected blockquote to remain, got %q", post.ArticleHTML)
	}
}

// Attribution + source link: @mention [via](url) → blockquote[cite] + footer with cite
func TestRenderMarkdownPlugin_BlockquoteAttributionMentionWithSource(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\n@jane [via](https://example.com/post)"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<blockquote cite="https://example.com/post"`) {
		t.Fatalf("expected blockquote with cite attr, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>") {
		t.Fatalf("expected footer, got %q", post.ArticleHTML)
	}

	// Source link should be wrapped in <cite>
	if !strings.Contains(post.ArticleHTML, "<cite><a href=\"https://example.com/post\"") {
		t.Fatalf("expected source link in <cite>, got %q", post.ArticleHTML)
	}
}

// Attribution + source: plain name, [Source Title](url) → blockquote[cite] + cite wrapping
func TestRenderMarkdownPlugin_BlockquoteAttributionPlainNameWithSource(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> Good systems make the easy path the correct path.\nJane Doe, [Designing Better Defaults](https://example.com/post)"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<blockquote cite="https://example.com/post"`) {
		t.Fatalf("expected blockquote with cite attr, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<cite><a href=\"https://example.com/post\"") {
		t.Fatalf("expected source link in <cite>, got %q", post.ArticleHTML)
	}
}

// Source-only attribution: [via](url) → blockquote[cite] + cite wrapping, no person
func TestRenderMarkdownPlugin_BlockquoteAttributionSourceOnly(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> Good systems make the easy path the correct path.\n[via](https://example.com/post)"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<blockquote cite="https://example.com/post"`) {
		t.Fatalf("expected blockquote with cite attr, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<cite><a href=\"https://example.com/post\"") {
		t.Fatalf("expected source link in <cite>, got %q", post.ArticleHTML)
	}
}

// Multi-paragraph quote with attribution inside blockquote
func TestRenderMarkdownPlugin_BlockquoteAttributionInsideBlockquote(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> This is the first paragraph.\n>\n> This is the second paragraph.\nJane Doe"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 Jane Doe</footer>") {
		t.Fatalf("expected footer inside blockquote, got %q", post.ArticleHTML)
	}

	// Footer should be inside blockquote
	if !strings.Contains(post.ArticleHTML, "</p>\n  <footer>") {
		t.Fatalf("expected footer after last paragraph, got %q", post.ArticleHTML)
	}
}

// Attribution not triggered when text is inside quote (has > prefix)
func TestRenderMarkdownPlugin_BlockquoteAttributionNotForInsideQuote(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> This is the quote.\n> Jane Doe"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if strings.Contains(post.ArticleHTML, "<footer>") {
		t.Fatalf("expected no footer for Jane Doe inside quote, got %q", post.ArticleHTML)
	}

	// Jane Doe should be inside the blockquote paragraph via lazy continuation
	if !strings.Contains(post.ArticleHTML, "<blockquote>") {
		t.Fatalf("expected blockquote, got %q", post.ArticleHTML)
	}
}

// Attribution: [Name](/people) + [Source](url) → first link person, second link source
func TestRenderMarkdownPlugin_BlockquoteAttributionLinkPersonWithSource(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> quoted text\n[Jane Doe](/people/jane), [Designing Better Defaults](https://example.com/post)"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<blockquote cite="https://example.com/post"`) {
		t.Fatalf("expected blockquote with cite attr, got %q", post.ArticleHTML)
	}

	// Person link should NOT be inside <cite>, source link SHOULD be
	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 <a href=\"/people/jane\">Jane Doe</a>, <cite><a href=\"https://example.com/post\"") {
		t.Fatalf("expected person link outside cite, source link inside cite, got %q", post.ArticleHTML)
	}
}

// Attribution not triggered after blank line; text becomes a normal paragraph
func TestRenderMarkdownPlugin_BlockquoteAttributionNormalAfterParagraph(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> This is the quote.\n\nThis is a normal paragraph after the quote."}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if strings.Contains(post.ArticleHTML, "<footer>") {
		t.Fatalf("expected no footer, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<p>This is a normal paragraph after the quote.</p>") {
		t.Fatalf("expected normal paragraph after blockquote, got %q", post.ArticleHTML)
	}
}

// Attribution followed by blank line then normal paragraph
func TestRenderMarkdownPlugin_BlockquoteAttributionWithFollowingParagraph(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> This is the quote.\nJane Doe\n\nThis is a normal paragraph after the quote."}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 Jane Doe</footer>") {
		t.Fatalf("expected footer with attribution, got %q", post.ArticleHTML)
	}

	if !strings.Contains(post.ArticleHTML, "<p>This is a normal paragraph after the quote.</p>") {
		t.Fatalf("expected following paragraph to remain, got %q", post.ArticleHTML)
	}
}

// Organization attribution
func TestRenderMarkdownPlugin_BlockquoteAttributionOrganization(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "> Move fast without breaking trust.\nThe Markata Project"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<footer>\u2014 The Markata Project</footer>") {
		t.Fatalf("expected footer with organization name, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Table(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "| Header 1 | Header 2 |\n|----------|----------|\n| Cell 1   | Cell 2   |"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<table>") {
		t.Errorf("expected <table> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<th>") {
		t.Errorf("expected <th> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<td>") {
		t.Errorf("expected <td> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Strikethrough(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "~~strikethrough~~"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<del>strikethrough</del>") {
		t.Errorf("expected <del>strikethrough</del> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_TaskList(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "- [ ] unchecked\n- [x] checked"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `type="checkbox"`) {
		t.Errorf("expected checkbox input in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Autolinks(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "Visit https://example.com for more"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<a href="https://example.com"`) {
		t.Errorf("expected autolinked URL in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_SkipPost(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{
		Content: "# Should not render",
		Skip:    true,
	}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if post.ArticleHTML != "" {
		t.Errorf("expected empty ArticleHTML for skipped post, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_EmptyContent(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: ""}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if post.ArticleHTML != "" {
		t.Errorf("expected empty ArticleHTML for empty content, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_RawHTML(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "<div class='custom'>content</div>"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `<div class='custom'>content</div>`) {
		t.Errorf("expected raw HTML in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_AutoHeadingIDs(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "# My Heading"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `id="`) {
		t.Errorf("expected heading ID in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Render(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	m := lifecycle.NewManager()

	posts := []*models.Post{
		{Content: "# Post 1"},
		{Content: "# Post 2"},
		{Content: "# Post 3", Skip: true},
	}
	m.SetPosts(posts)

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	resultPosts := m.Posts()

	if !strings.Contains(resultPosts[0].ArticleHTML, "<h1") {
		t.Errorf("Post 1 not rendered correctly: %q", resultPosts[0].ArticleHTML)
	}
	if !strings.Contains(resultPosts[1].ArticleHTML, "<h1") {
		t.Errorf("Post 2 not rendered correctly: %q", resultPosts[1].ArticleHTML)
	}
	if resultPosts[2].ArticleHTML != "" {
		t.Errorf("Skipped post should have empty ArticleHTML: %q", resultPosts[2].ArticleHTML)
	}
}

// Admonition tests

func TestAdmonition_BasicNote(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "!!! note \"Important\"\n    This is important information"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, `class="admonition note"`) {
		t.Errorf("expected admonition note class in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `class="admonition-title"`) {
		t.Errorf("expected admonition-title class in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "Important") {
		t.Errorf("expected title 'Important' in output, got %q", post.ArticleHTML)
	}
}

func TestAdmonition_Types(t *testing.T) {
	types := []string{"note", "warning", "tip", "important", "danger", "caution"}

	p := NewRenderMarkdownPlugin()

	for _, adType := range types {
		t.Run(adType, func(t *testing.T) {
			post := &models.Post{Content: "!!! " + adType + "\n    Content"}

			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}

			expected := `class="admonition ` + adType + `"`
			if !strings.Contains(post.ArticleHTML, expected) {
				t.Errorf("expected %q in output, got %q", expected, post.ArticleHTML)
			}
		})
	}
}

func TestRenderMarkdownPlugin_AdmonitionDefaultTitle(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "!!! warning\n    Warning content"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	// Should use capitalized type as default title
	if !strings.Contains(post.ArticleHTML, "Warning") {
		t.Errorf("expected default title 'Warning' in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_AdmonitionCustomType(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "!!! source \"Primary Source\"\n    Content"}

	if err := p.renderPost(post); err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	for _, want := range []string{
		"<div class=\"admonition source\">",
		"<p class=\"admonition-title\">Primary Source</p>",
		"<p>Content</p>",
	} {
		if !strings.Contains(post.ArticleHTML, want) {
			t.Errorf("expected %q in output, got %q", want, post.ArticleHTML)
		}
	}
}

// Interface compliance tests

func TestRenderMarkdownPlugin_Interfaces(_ *testing.T) {
	p := NewRenderMarkdownPlugin()

	// Verify interface compliance
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.RenderPlugin = p
}

// Highlight configuration tests

func TestRenderMarkdownPlugin_ConfigureWithPalette(t *testing.T) {
	tests := []struct {
		name        string
		extra       map[string]interface{}
		wantContain string // string that should be in rendered code
	}{
		{
			name: "catppuccin-mocha palette",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "catppuccin-mocha",
				},
			},
			wantContain: `class="chroma"`, // Uses CSS classes, not inline styles
		},
		{
			name: "explicit theme override",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "catppuccin-mocha",
				},
				"markdown": map[string]interface{}{
					"highlight": map[string]interface{}{
						"theme": "dracula",
					},
				},
			},
			wantContain: `class="kd"`, // keyword declaration class
		},
		{
			name:        "no config uses default",
			extra:       map[string]interface{}{},
			wantContain: `class="chroma"`, // Uses CSS classes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRenderMarkdownPlugin()
			m := lifecycle.NewManager()
			m.Config().Extra = tt.extra

			err := p.Configure(m)
			if err != nil {
				t.Fatalf("Configure error: %v", err)
			}

			// Verify rendering works with the configured theme
			post := &models.Post{Content: "```go\nfunc main() {}\n```"}
			err = p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}

			if !strings.Contains(post.ArticleHTML, tt.wantContain) {
				t.Errorf("expected %q in output, got %q", tt.wantContain, post.ArticleHTML)
			}
		})
	}
}

func TestRenderMarkdownPlugin_ResolveHighlightConfig(t *testing.T) {
	tests := []struct {
		name      string
		extra     map[string]interface{}
		wantTheme string
		wantLN    bool
	}{
		{
			name:      "empty config - default dark theme",
			extra:     map[string]interface{}{},
			wantTheme: "github-dark", // Default for dark variant
			wantLN:    false,
		},
		{
			name: "explicit theme",
			extra: map[string]interface{}{
				"markdown": map[string]interface{}{
					"highlight": map[string]interface{}{
						"theme": "monokai",
					},
				},
			},
			wantTheme: "monokai",
			wantLN:    false,
		},
		{
			name: "line numbers enabled",
			extra: map[string]interface{}{
				"markdown": map[string]interface{}{
					"highlight": map[string]interface{}{
						"theme":        "dracula",
						"line_numbers": true,
					},
				},
			},
			wantTheme: "dracula",
			wantLN:    true,
		},
		{
			name: "palette-derived theme - catppuccin-mocha",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "catppuccin-mocha",
				},
			},
			wantTheme: "catppuccin-mocha",
			wantLN:    false,
		},
		{
			name: "palette-derived theme - gruvbox-light",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "gruvbox-light",
				},
			},
			wantTheme: "gruvbox-light",
			wantLN:    false,
		},
		{
			name: "explicit theme overrides palette",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "catppuccin-mocha",
				},
				"markdown": map[string]interface{}{
					"highlight": map[string]interface{}{
						"theme": "nord",
					},
				},
			},
			wantTheme: "nord",
			wantLN:    false,
		},
		{
			name: "unknown palette uses variant default - light",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "custom-light-theme",
				},
			},
			wantTheme: "github", // Light variant default
			wantLN:    false,
		},
		{
			name: "unknown palette uses variant default - dark",
			extra: map[string]interface{}{
				"theme": map[string]interface{}{
					"palette": "custom-dark-theme",
				},
			},
			wantTheme: "github-dark", // Dark variant default
			wantLN:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRenderMarkdownPlugin()
			gotTheme, gotLN := p.resolveHighlightConfig(tt.extra)

			if gotTheme != tt.wantTheme {
				t.Errorf("theme = %q, want %q", gotTheme, tt.wantTheme)
			}
			if gotLN != tt.wantLN {
				t.Errorf("lineNumbers = %v, want %v", gotLN, tt.wantLN)
			}
		})
	}
}

func TestRenderMarkdownPlugin_AttributeSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "image with class",
			input:    "![alt text](image.webp){.more-cinematic}",
			expected: `class="more-cinematic"`,
		},
		{
			name:     "image with multiple classes",
			input:    "![photo](photo.jpg){.shadow .bordered}",
			expected: `class="shadow bordered"`,
		},
		{
			name:     "image with id",
			input:    "![hero](hero.png){#hero-image}",
			expected: `id="hero-image"`,
		},
		{
			name:     "image with class and id",
			input:    "![banner](banner.jpg){#main-banner .full-width}",
			expected: `id="main-banner"`,
		},
		{
			name:     "heading with class",
			input:    "## Section Title {.highlighted}",
			expected: `class="highlighted"`,
		},
		{
			name:     "heading with custom id",
			input:    "## Installation {#install}",
			expected: `id="install"`,
		},
		{
			name:     "link with class",
			input:    "[Click here](https://example.com){.external}",
			expected: `class="external"`,
		},
		{
			name:     "link with id",
			input:    "[Main link](https://example.com){#main-link}",
			expected: `id="main-link"`,
		},
		{
			name:     "paragraph with class",
			input:    "A paragraph\n{.lead}",
			expected: `class="lead"`,
		},
		{
			name:     "paragraph with class no space",
			input:    "A paragraph{.lead}",
			expected: `class="lead"`,
		},
		{
			name:     "blockquote with class",
			input:    "> Quoted text\n{.callout}",
			expected: `class="callout"`,
		},
		{
			name:     "blockquote with class line",
			input:    "> Quoted text\n> {.callout}",
			expected: `class="callout"`,
		},
		{
			name:     "list item with class",
			input:    "- Item text\n", // list item class syntax is not supported
			expected: `<li>`,
		},
		{
			name:     "list item with inline class",
			input:    "- hello {.item}",
			expected: `class="item"`,
		},
		{
			name:     "list item with inline class no space",
			input:    "- hello{.item}",
			expected: `class="item"`,
		},
		{
			name:     "list item with attribute line",
			input:    "- hello\n  {.item}",
			expected: `class="item"`,
		},
		{
			name:     "emphasis with class",
			input:    "*note*{.accent}",
			expected: `class="accent"`,
		},
		{
			name:     "inline code with class",
			input:    "`code`{.mono}",
			expected: `class="mono"`,
		},
		{
			name:     "strikethrough with class",
			input:    "~~old~~{.deprecated}",
			expected: `class="deprecated"`,
		},
		{
			name:     "table cell with class",
			input:    "| Col |\n|---|\n| Cell {.cell} |",
			expected: `class="cell"`,
		},
	}

	p := NewRenderMarkdownPlugin()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Content: tt.input}
			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}
			if !strings.Contains(post.ArticleHTML, tt.expected) {
				t.Errorf("expected %q in output, got %q", tt.expected, post.ArticleHTML)
			}
		})
	}
}

func TestRenderMarkdownPlugin_AttributeSyntax_ImageClassInOutput(t *testing.T) {
	// Specific test for the issue #404 use case
	p := NewRenderMarkdownPlugin()
	post := &models.Post{Content: "![alt text](image.webp){.more-cinematic}"}

	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	// Verify the img tag has the class attribute
	if !strings.Contains(post.ArticleHTML, `<img`) {
		t.Errorf("expected <img> tag in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, `class="more-cinematic"`) {
		t.Errorf("expected class=\"more-cinematic\" in output, got %q", post.ArticleHTML)
	}
	// Verify the attribute syntax {.more-cinematic} is NOT in the output
	if strings.Contains(post.ArticleHTML, "{.more-cinematic}") {
		t.Errorf("attribute syntax should be removed from output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_GetPaletteVariant(t *testing.T) {
	p := NewRenderMarkdownPlugin()

	tests := []struct {
		palette string
		want    string
	}{
		// Light variants
		{"default-light", "light"},
		{"gruvbox-light", "light"},
		{"catppuccin-latte", "light"},
		{"rose-pine-dawn", "light"},
		{"tokyo-night-day", "light"},
		{"kanagawa-lotus", "light"},

		// Dark variants
		{"default-dark", "dark"},
		{"gruvbox-dark", "dark"},
		{"catppuccin-mocha", "dark"},
		{"rose-pine", "dark"},
		{"tokyo-night", "dark"},
		{"dracula", "dark"},
		{"matte-black", "dark"},
	}

	for _, tt := range tests {
		t.Run(tt.palette, func(t *testing.T) {
			got := p.getPaletteVariant(tt.palette)
			if string(got) != tt.want {
				t.Errorf("getPaletteVariant(%q) = %q, want %q", tt.palette, got, tt.want)
			}
		})
	}
}

func TestRenderMarkdownPlugin_DetectCSSRequirements(t *testing.T) {
	tests := []struct {
		name               string
		content            string
		wantAdmonitionsCSS bool
		wantCodeCSS        bool
	}{
		{
			name:               "plain text - no CSS needed",
			content:            "Hello world, this is plain text.",
			wantAdmonitionsCSS: false,
			wantCodeCSS:        false,
		},
		{
			name:               "fenced code block - needs code CSS",
			content:            "```go\nfunc main() {}\n```",
			wantAdmonitionsCSS: false,
			wantCodeCSS:        true,
		},
		{
			name:               "inline code - no special code CSS needed",
			content:            "Use `fmt.Println()` to print.",
			wantAdmonitionsCSS: false,
			wantCodeCSS:        false, // inline code doesn't need syntax highlighting CSS
		},
		{
			name:               "headings and links - no special CSS",
			content:            "# Heading\n\n[Link](https://example.com)",
			wantAdmonitionsCSS: false,
			wantCodeCSS:        false,
		},
		{
			name:               "code block without language - still needs code CSS",
			content:            "```\nplain code block\n```",
			wantAdmonitionsCSS: false,
			wantCodeCSS:        true,
		},
		{
			name:               "admonition - needs admonitions CSS only",
			content:            "!!! note \"Note\"\n    This is a note admonition.",
			wantAdmonitionsCSS: true,
			wantCodeCSS:        false, // indented body is a paragraph, not a code block
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRenderMarkdownPlugin()
			post := &models.Post{Content: tt.content}

			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}

			// Check admonitions CSS flag
			gotAdmonitions := false
			if post.Extra != nil {
				if v, ok := post.Extra["needs_admonitions_css"].(bool); ok {
					gotAdmonitions = v
				}
			}
			if gotAdmonitions != tt.wantAdmonitionsCSS {
				t.Errorf("needs_admonitions_css = %v, want %v", gotAdmonitions, tt.wantAdmonitionsCSS)
			}

			// Check code CSS flag
			gotCode := false
			if post.Extra != nil {
				if v, ok := post.Extra["needs_code_css"].(bool); ok {
					gotCode = v
				}
			}
			if gotCode != tt.wantCodeCSS {
				t.Errorf("needs_code_css = %v, want %v", gotCode, tt.wantCodeCSS)
			}
		})
	}
}

func TestRenderMarkdownPlugin_Typographer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "smart double quotes",
			input:    `"hello world"`,
			expected: "&ldquo;", // left double quotation mark
		},
		{
			name:     "smart single quotes/apostrophe",
			input:    `It's a nice day`,
			expected: "&rsquo;", // right single quotation mark (apostrophe)
		},
		{
			name:     "en dash",
			input:    `9--5`,
			expected: "&ndash;",
		},
		{
			name:     "em dash",
			input:    `hello---world`,
			expected: "&mdash;",
		},
		{
			name:     "ellipsis",
			input:    `wait...`,
			expected: "&hellip;",
		},
	}

	p := NewRenderMarkdownPlugin()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Content: tt.input}
			err := p.renderPost(post)
			if err != nil {
				t.Fatalf("renderPost error: %v", err)
			}
			if !strings.Contains(post.ArticleHTML, tt.expected) {
				t.Errorf("expected %q in output, got %q", tt.expected, post.ArticleHTML)
			}
		})
	}
}

func TestRenderMarkdownPlugin_DefinitionList(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	content := `Term 1
:   Definition 1

Term 2
:   Definition 2a
:   Definition 2b`

	post := &models.Post{Content: content}
	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "<dl>") {
		t.Errorf("expected <dl> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<dt>") {
		t.Errorf("expected <dt> in output, got %q", post.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "<dd>") {
		t.Errorf("expected <dd> in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_Footnote(t *testing.T) {
	p := NewRenderMarkdownPlugin()
	content := `Here's a sentence with a footnote.[^1]

[^1]: This is the footnote content.`

	post := &models.Post{Content: content}
	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}

	if !strings.Contains(post.ArticleHTML, "footnote") {
		t.Errorf("expected 'footnote' class in output, got %q", post.ArticleHTML)
	}
}

func TestRenderMarkdownPlugin_ExtensionConfig_Disabled(t *testing.T) {
	// Test that extensions can be disabled via config
	config := MarkdownExtensionConfig{
		TypographerEnabled:    false,
		DefinitionListEnabled: false,
		FootnoteEnabled:       false,
	}

	p := &RenderMarkdownPlugin{
		md: createMarkdownRenderer("github", false, config),
	}

	// Test typographer is disabled - straight quotes should remain
	post := &models.Post{Content: `"hello"`}
	err := p.renderPost(post)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}
	// With typographer disabled, quotes should remain as straight quotes
	// Check that typographic quotes are NOT present
	if strings.Contains(post.ArticleHTML, "&ldquo;") || strings.Contains(post.ArticleHTML, "&rdquo;") {
		t.Error("Typographer should be disabled but smart quotes found in output")
	}

	// Test definition list is disabled - should render as paragraphs
	post2 := &models.Post{Content: "Term\n:   Definition"}
	err = p.renderPost(post2)
	if err != nil {
		t.Fatalf("renderPost error: %v", err)
	}
	// Without definition list extension, should not have <dl>
	if strings.Contains(post2.ArticleHTML, "<dl>") {
		t.Error("Definition list should be disabled but <dl> found in output")
	}
}

func TestRenderMarkdownPlugin_DetectRoads(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"road container", "::: road timeline\n\n- [Kickoff](#a) Jan\n:::", true},
		{"road container with attributes", "::: road {data-road=\"git\"}\n\n- [Init](#a) main\n:::", true},
		{"raw div with data-road", "<div data-road=\"winding\">\n\n- [Start](#a)\n\n</div>", true},
		{"similar class name", "<div class=\"roadmap\">plan</div>", false},
		{"plain list", "- [Kickoff](#a) Jan", false},
		{"word road in text", "The road goes ever on.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewRenderMarkdownPlugin()
			post := &models.Post{Content: tt.content}
			if err := p.renderPost(post); err != nil {
				t.Fatalf("renderPost error: %v", err)
			}
			got := post.Extra["needs_roads"] == true
			if got != tt.want {
				t.Errorf("needs_roads = %v, want %v\nhtml: %s", got, tt.want, post.ArticleHTML)
			}
		})
	}
}
