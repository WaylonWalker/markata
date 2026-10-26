package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestGlossaryPlugin_EqualLengthAliasPriority(t *testing.T) {
	title := "Self Host"
	definition := &models.Post{
		Slug: "self-host", Href: "/self-host/", Title: &title,
		Extra: map[string]interface{}{
			"templateKey": "glossary",
			"aliases":     []string{"self-hosted", "self hosted"},
		},
	}

	for range 100 {
		plugin := NewGlossaryPlugin()
		if err := plugin.buildGlossary([]*models.Post{definition}); err != nil {
			t.Fatal(err)
		}
		input := "<p>self-hosted appears before self hosted.</p>"
		got := plugin.linkTerms(input, &models.Post{Slug: "other"})
		if !strings.HasPrefix(got, "<p>self-hosted appears before <a ") || strings.Count(got, `href="/self-host/"`) != 1 {
			t.Fatalf("equal-length aliases changed lexical selection: %s", got)
		}
		if got := plugin.computeTermsHash(); got == buildcache.ContentHash("Self Host\x00/self-host/\x00\x00self-hosted\x01self hosted\x01\x00") {
			t.Fatal("new matching order reused the pre-fix glossary cache identity")
		}
	}
}

func TestGlossaryPlugin_RestoresNestedProtectedContent(t *testing.T) {
	plugin := NewGlossaryPlugin()
	plugin.terms["kirby"] = &GlossaryTerm{Term: "Kirby", Slug: "kirby", Href: "/kirby/"}
	for _, input := range []string{
		`<pre><code><a href="/thought-200/">The One Eyed Fighting Kirby</a></code></pre>`,
		`<code><a href="/thought-200/">The One Eyed Fighting Kirby</a></code>`,
		`<pre><code>Kirby <a href="/thought-200/">Kirby</a></code></pre><code>Kirby</code>`,
	} {
		for range 100 {
			if got := plugin.linkTerms(input, &models.Post{Slug: "other"}); got != input {
				t.Fatalf("protected nested content changed:\ngot: %q\nwant: %q", got, input)
			}
		}
	}
}
