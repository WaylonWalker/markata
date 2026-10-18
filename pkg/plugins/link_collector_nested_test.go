package plugins

import (
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestExtractHrefsAndTextNestedHTML(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		wantHrefs     []string
		wantTextByURL map[string]string
	}{
		{
			name:          "nested span",
			html:          `<a href="/page"><span>Nested text</span></a>`,
			wantHrefs:     []string{"/page"},
			wantTextByURL: map[string]string{"/page": "Nested text"},
		},
		{
			name:          "image and text",
			html:          `<a href="/image"><img alt="icon" src="icon.png"> Open image</a>`,
			wantHrefs:     []string{"/image"},
			wantTextByURL: map[string]string{"/image": "Open image"},
		},
		{
			name:          "card divs",
			html:          `<a href="/card"><div class="title">Card title</div><div class="body">Card body</div></a>`,
			wantHrefs:     []string{"/card"},
			wantTextByURL: map[string]string{"/card": "Card title Card body"},
		},
		{
			name:          "image only still collects href",
			html:          `<a href="/photo"><img alt="Photo" src="photo.jpg"></a>`,
			wantHrefs:     []string{"/photo"},
			wantTextByURL: map[string]string{"/photo": ""},
		},
		{
			name:          "duplicate keeps first link text",
			html:          `<a href="/same"><span>First</span></a><a href="/same">Second</a>`,
			wantHrefs:     []string{"/same"},
			wantTextByURL: map[string]string{"/same": "First"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHrefs, gotTextByURL := extractHrefsAndText(tt.html)
			if !reflect.DeepEqual(gotHrefs, tt.wantHrefs) {
				t.Fatalf("extractHrefsAndText() hrefs = %v, want %v", gotHrefs, tt.wantHrefs)
			}
			if !reflect.DeepEqual(gotTextByURL, tt.wantTextByURL) {
				t.Fatalf("extractHrefsAndText() text = %v, want %v", gotTextByURL, tt.wantTextByURL)
			}
		})
	}
}

func TestLinkCollectorPluginRenderNestedAnchor(t *testing.T) {
	p := NewLinkCollectorPlugin()
	p.SetSiteURL("https://example.com")
	m := lifecycle.NewManager()

	sourceTitle := "Source"
	targetTitle := "Target"
	posts := []*models.Post{
		{
			Path:        "source.md",
			Slug:        "source",
			Href:        "/source/",
			Title:       &sourceTitle,
			ArticleHTML: `<a href="/target/"><div class="card"><span>Target card</span></div></a>`,
		},
		{
			Path:        "target.md",
			Slug:        "target",
			Href:        "/target/",
			Title:       &targetTitle,
			ArticleHTML: `<p>Target body</p>`,
		},
	}
	m.SetPosts(posts)

	if err := p.Render(m); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	source := m.Posts()[0]
	target := m.Posts()[1]
	if !reflect.DeepEqual(source.Hrefs, []string{"/target/"}) {
		t.Fatalf("source.Hrefs = %v, want [/target/]", source.Hrefs)
	}
	if len(source.Outlinks) != 1 {
		t.Fatalf("source.Outlinks = %d, want 1", len(source.Outlinks))
	}
	if source.Outlinks[0].TargetPost != target {
		t.Fatalf("outlink target = %#v, want target post", source.Outlinks[0].TargetPost)
	}
	if source.Outlinks[0].SourceText != "Target card" {
		t.Fatalf("outlink SourceText = %q, want %q", source.Outlinks[0].SourceText, "Target card")
	}
	if len(target.Inlinks) != 1 {
		t.Fatalf("target.Inlinks = %d, want 1", len(target.Inlinks))
	}
	if target.Inlinks[0].SourcePost != source {
		t.Fatalf("inlink source = %#v, want source post", target.Inlinks[0].SourcePost)
	}
}
