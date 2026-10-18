package templates

import (
	"strings"
	"testing"
)

// A feed card must render the article's Markdown output, not the raw Markdown
// source or the plain-text description intended for metadata.
func TestNoteCardPreservesRenderedMarkdown(t *testing.T) {
	for _, tree := range []string{"../../templates", "../themes/default/templates"} {
		t.Run(tree, func(t *testing.T) {
			engine, err := NewEngine(tree)
			if err != nil {
				t.Fatal(err)
			}
			ctx := map[string]interface{}{
				"post": map[string]interface{}{
					"href":         "/ping-84/",
					"title":        "If an agent can't verify",
					"content":      "Verify ==nothing else matters==, **bold**, *italic*, [read](https://example.com).",
					"description":  "Plain metadata description, not the visible Markdown excerpt.",
					"article_html": "<p>Verify <mark>nothing else matters</mark>, <strong>bold</strong>, <em>italic</em>, <a href=\"https://example.com\">read</a>.</p>",
				},
			}

			got, err := engine.RenderToString("partials/cards/note-card.html", ctx)
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{
				"<mark>nothing else matters</mark>",
				"<strong>bold</strong>",
				"<em>italic</em>",
				"<a href=\"https://example.com\">read</a>",
			} {
				if !strings.Contains(got, want) {
					t.Errorf("note card missing formatted excerpt %q: %s", want, got)
				}
			}
			for _, unwanted := range []string{
				"==nothing else matters==",
				"**bold**",
				"*italic*",
				"[read](https://example.com)",
				"Plain metadata description",
				"&lt;mark&gt;",
			} {
				if strings.Contains(got, unwanted) {
					t.Errorf("note card leaked raw/escaped Markdown %q: %s", unwanted, got)
				}
			}
		})
	}
}

func TestNoteCardEscapesPlainDescriptionFallback(t *testing.T) {
	for _, tree := range []string{"../../templates", "../themes/default/templates"} {
		t.Run(tree, func(t *testing.T) {
			engine, err := NewEngine(tree)
			if err != nil {
				t.Fatal(err)
			}
			got, err := engine.RenderToString("partials/cards/note-card.html", map[string]interface{}{
				"post": map[string]interface{}{
					"href":        "/ping-fallback/",
					"title":       "No rendered body",
					"description": "Escaped <script>alert(1)</script> & still readable",
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(got, "<script>") ||
				!strings.Contains(got, "&lt;script&gt;alert(1)&lt;/script&gt;") ||
				!strings.Contains(got, "still readable") {
				t.Fatalf("plain-text fallback not properly escaped: %s", got)
			}
		})
	}
}
