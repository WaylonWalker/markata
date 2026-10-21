package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

// The sample page contains both "Clippy" and "Clippy No Simpy" glossary
// entries. A later short-term pass must not corrupt the long-term anchor's
// attributes or insert a nested anchor into its visible text.
func TestGlossaryPlugin_OverlappingTermsDoNotCorruptNewLinks(t *testing.T) {
	shortTitle := "Clippy"
	longTitle := "Clippy No Simpy"
	longDescription := `"Clippy No Simpy" warns about defending bad business practices.`
	p := NewGlossaryPlugin()
	if err := p.buildGlossary([]*models.Post{
		{
			Title: &shortTitle, Slug: "clippy", Href: "/glossary/clippy/",
			Extra: map[string]interface{}{"templateKey": "glossary"},
		},
		{
			Title: &longTitle, Description: &longDescription,
			Slug: "clippy-no-simpy", Href: "/glossary/clippy-no-simpy/",
			Extra: map[string]interface{}{"templateKey": "glossary"},
		},
	}); err != nil {
		t.Fatal(err)
	}

	input := `<p>clippy no simpy and clippy; <a href="/other">clippy</a>.</p>`
	want := `<p><a href="/glossary/clippy-no-simpy/" class="glossary-term" title="&#34;Clippy No Simpy&#34; warns about defending bad business practices." data-hover-title="Clippy No Simpy">clippy no simpy</a> and <a href="/glossary/clippy/" class="glossary-term">clippy</a>; <a href="/other">clippy</a>.</p>`
	for range 10 {
		got := p.linkTerms(input, &models.Post{Slug: "sample"})
		if got != want {
			t.Fatalf("overlapping glossary terms corrupted the generated HTML:\ngot:  %s\nwant: %s", got, want)
		}
		if strings.Contains(got, "\x00") {
			t.Fatalf("internal placeholders leaked into the page: %q", got)
		}
	}
}

func TestGlossaryPlugin_OverlappingTermsWithUnlimitedLinks(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.MaxLinksPerTerm = 0
	p.terms = map[string]*GlossaryTerm{
		"markdown": {Term: "Markdown", Slug: "markdown", Href: "/markdown/"},
		"mark":     {Term: "Mark", Slug: "mark", Href: "/mark/"},
	}
	got := p.linkTerms("<p>Markdown Markdown and Mark.</p>", &models.Post{Slug: "sample"})
	if strings.Count(got, `href="/markdown/"`) != 2 ||
		strings.Count(got, `href="/mark/"`) != 1 ||
		strings.Count(got, "<a ") != 3 {
		t.Fatalf("generated links were modified or counts changed: %s", got)
	}
}
