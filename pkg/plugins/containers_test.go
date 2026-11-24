// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/anchor"
)

func renderContainerMarkdown(t *testing.T, input string) string {
	t.Helper()

	md := goldmark.New(goldmark.WithExtensions(&ContainerExtension{}))
	var output bytes.Buffer
	if err := md.Convert([]byte(input), &output); err != nil {
		t.Fatalf("Convert() error = %v", err)
	}
	return output.String()
}

func TestContainerExtension_SectionHead(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"heading and intro", "## Heading\n\nIntro.\n\nBody.", []string{
			`<div class="md-section-head">`, "<h2>Heading</h2>", "<p>Intro.</p>",
			"</div>", `<div class="md-section-body">`, "<p>Body.</p>", "</div>", "</section>",
		}},
		{"heading only", "## Heading\n\n- Body", []string{
			`<div class="md-section-head">`, "<h2>Heading</h2>", "</div>",
			`<div class="md-section-body">`, "<ul>", "<li>Body</li>",
		}},
		{"no heading", "Body.\n\n## Later heading", []string{
			`<div class="md-section-body">`, "<p>Body.</p>", "<h2>Later heading</h2>",
		}},
		{"other heading", "### Heading\n\nBody.", []string{
			`<div class="md-section-body">`, "<h3>Heading</h3>", "<p>Body.</p>",
		}},
		{"empty", "", []string{`<div class="md-section-body">`, "</div>", "</section>"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderContainerMarkdown(t, "::: md-section {#era}\n\n"+tt.input+"\n\n:::")
			assertContainerHTMLInOrder(t, got, append([]string{`<section id="era" class="md-section">`}, tt.want...)...)
			if !strings.HasPrefix(tt.input, "## ") && strings.Contains(got, `class="md-section-head"`) {
				t.Fatalf("unexpected head: %s", got)
			}
		})
	}
}

func TestContainerExtension_SectionKickerPreservesHeadingID(t *testing.T) {
	for _, position := range []anchor.Position{anchor.Before, anchor.After} {
		md := goldmark.New(
			goldmark.WithExtensions(&ContainerExtension{}, &anchor.Extender{Position: position}),
			goldmark.WithParserOptions(parser.WithAutoHeadingID(),
				parser.WithASTTransformers(util.Prioritized(&HeadingIDTransformer{}, -90))),
		)
		var output bytes.Buffer
		if err := md.Convert([]byte("::: md-section\n\n## 2026 · Agent workflows and **access**\n\nIntro.\n\n:::"), &output); err != nil {
			t.Fatal(err)
		}
		got := output.String()
		assertContainerHTMLInOrder(t, got,
			`<h2 id="2026-agent-workflows-and-access">`,
			`<span class="md-section-kicker">2026</span>`,
			`Agent workflows and <strong>access</strong>`,
		)
		if !strings.Contains(got, `href="#2026-agent-workflows-and-access"`) {
			t.Fatalf("permalink lost original ID: %s", got)
		}
	}
}

func TestContainerExtension_SectionKickerRequiresShortTokenAndSeparator(t *testing.T) {
	for _, heading := range []string{"2026 Agent workflows", "2026·Agent workflows", "A long heading · Title", "12345678901234567 · Title", "2026 · "} {
		got := renderContainerMarkdown(t, "::: md-section\n\n## "+heading+"\n\n:::")
		if strings.Contains(got, `class="md-section-kicker"`) {
			t.Errorf("unexpected split for %q: %s", heading, got)
		}
	}
}

func TestContainerExtension_SectionKickerWithDefaultExtensions(t *testing.T) {
	md := NewRenderMarkdownPlugin().md
	for _, heading := range []string{"2024 · Delivery", "2023 · Foundations", "2026 · **Access**"} {
		var output bytes.Buffer
		if err := md.Convert([]byte("::: md-section\n\n## "+heading+"\n\n:::"), &output); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(output.String(), `class="md-section-kicker"`) {
			t.Errorf("missing kicker for %q: %s", heading, output.String())
		}
	}
}

func TestContainerExtension_GridVariants(t *testing.T) {
	for _, variant := range []string{"rail", "pending"} {
		got := renderContainerMarkdown(t, "::: md-grid "+variant+" {data-columns=\"3\"}\n\n:::: md-card\n\n### Card\n\n::::\n\n:::")
		assertContainerHTMLInOrder(t, got,
			`<div class="md-grid `+variant+`" data-columns="3">`,
			`<div class="md-card">`, "<h3>Card</h3>", "</div>", "</div>",
		)
	}
}

func assertContainerHTMLInOrder(t *testing.T, got string, wantInOrder ...string) {
	t.Helper()

	cursor := 0
	for _, want := range wantInOrder {
		idx := strings.Index(got[cursor:], want)
		if idx < 0 {
			t.Fatalf("rendered HTML missing %q after byte %d\nGot:\n%s", want, cursor, got)
		}
		cursor += idx + len(want)
	}
}

func TestContainerExtension_NestedContainersUseMatchingDelimiterDepth(t *testing.T) {
	got := renderContainerMarkdown(t, `::: outer
outer before

:::: inner
inner body
::::

outer after
:::`)

	assertContainerHTMLInOrder(t, got,
		`<div class="outer">`,
		`<p>outer before</p>`,
		`<div class="inner">`,
		`<p>inner body</p>`,
		`</div>`,
		`<p>outer after</p>`,
		`</div>`,
	)
}

func TestContainerCloseMarkerMustBeColonOnly(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{name: "too short", line: "::", want: false},
		{name: "trailing whitespace", line: "::: ", want: false},
		{name: "three colons exact", line: ":::", want: true},
		{name: "four colons exact", line: "::::", want: true},
		{name: "named container", line: "::: card", want: false},
		{name: "attribute container", line: `::: {.card}`, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isContainerClose(tt.line); got != tt.want {
				t.Fatalf("isContainerClose(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestContainerExtension_RichMarkdownLayoutPrimitives(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-layout

:::: md-box {data-align="center"}
Cluster

::::

:::: md-columns {data-gap="sm"}

::::: md-card
Website

:::::

::::: md-card
Builder

:::::

::::

:::: md-connector {data-direction="both"}

::::

:::: md-stack

::::: md-divider

:::::

::::

:::`)
	assertContainerHTMLInOrder(t, got,
		`<div class="md-layout">`,
		`<div class="md-box" data-align="center">`,
		`<p>Cluster</p>`,
		`</div>`,
		`<div class="md-columns" data-gap="sm">`,
		`<div class="md-card">`,
		`<p>Website</p>`,
		`</div>`,
		`<div class="md-card">`,
		`<p>Builder</p>`,
		`</div>`,
		`</div>`,
		`<div class="md-connector" data-direction="both">`,
		`</div>`,
		`<div class="md-stack">`,
		`<div class="md-divider">`,
		`</div>`,
		`</div>`,
		`</div>`,
	)
}

func TestSourceDisclosureRendersNativeDetailsWithEscapedLabel(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-source {data-label="Show & source"}

~~~markdown
# Hello **Markdown**
~~~

:::`)
	assertContainerHTMLInOrder(t, got,
		`<details class="md-source" data-markata-ui="source-preview"><summary>Show &amp; source</summary>`,
		`<pre><code class="language-markdown">`,
		`# Hello **Markdown**`,
		`</code></pre>`,
		`</details>`,
	)
	if strings.Contains(got, `<div class="md-source"`) {
		t.Fatalf("md-source must use native details, got: %s", got)
	}
}

func TestSourceDisclosureDefaultsToReadableLabel(t *testing.T) {
	got := renderContainerMarkdown(t, `::: md-source

A short example.

:::`)
	assertContainerHTMLInOrder(t, got,
		`<details class="md-source" data-markata-ui="source-preview"><summary>Show Markdown source</summary>`,
		`<p>A short example.</p>`,
		`</details>`,
	)
}
