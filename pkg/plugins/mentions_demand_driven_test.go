package plugins

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestMentionsReferencedHandlesDemandDriven(t *testing.T) {
	p := NewMentionsPlugin()
	alice := &mentionEntry{Handle: "alice", SiteURL: "https://alice.example"}
	bob := &mentionEntry{Handle: "bob", SiteURL: "https://bob.example"}
	domain := &mentionEntry{Handle: "domain.example", SiteURL: "https://domain.example"}
	handles := map[string]*mentionEntry{
		"alice": alice, "alias": alice, "bob": bob, "domain.example": domain,
	}
	tests := []struct {
		name  string
		posts []*models.Post
		want  []string
	}{
		{"no mentions", []*models.Post{{Content: "Just a normal article", Published: true}}, nil},
		{"unknown mention", []*models.Post{{Content: "Hello @missing", Published: true}}, nil},
		{"fenced code only", []*models.Post{{Content: "```text\n@alice\n```", Published: true}}, nil},
		{"nonchat header", []*models.Post{{Content: "!!! note @alice\n    nothing else", Published: true}}, nil},
		{"domain with sentence punctuation", []*models.Post{{Content: "Thanks @domain.example.", Published: true}}, []string{"domain.example"}},
		{"known with unused contact", []*models.Post{{Content: "Hello @alice", Published: true}}, []string{"alice"}},
		{"alias and punctuation", []*models.Post{{Content: "Thanks @ALIAS.", Published: true}}, []string{"alias"}},
		{"chat title", []*models.Post{{Content: "!!! chat \"@alice\"\n    Hello", Published: true}}, []string{"alice"}},
		{"attribution", []*models.Post{{Content: "Quoted from ~bob", Published: true}}, []string{"bob"}},
		{"skip other post", []*models.Post{{Content: "@alice", Skip: true}, {Content: "@bob", Published: true}}, []string{"bob"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := lifecycle.NewManager()
			m.SetPosts(tc.posts)
			got := p.referencedHandles(m, handles)
			if len(got) != len(tc.want) {
				t.Fatalf("got %d handles, want %d: %#v", len(got), len(tc.want), got)
			}
			for _, handle := range tc.want {
				if got[handle] != handles[handle] {
					t.Errorf("missing referenced handle %s", handle)
				}
			}
		})
	}
}

func TestMentionsTransformPlainPageWithoutBlogrollMetadata(t *testing.T) {
	p := NewMentionsPlugin()
	m := lifecycle.NewManager()
	active := true
	m.Config().Extra = map[string]interface{}{
		"blogroll": models.BlogrollConfig{
			Enabled: true,
			Feeds: []models.ExternalFeedConfig{{
				URL: "https://unrelated.invalid/feed.xml", Handle: "unrelated",
				SiteURL: "https://unrelated.invalid", Active: &active,
			}},
		},
	}
	const plain = "A post without mentions."
	m.AddPost(&models.Post{Path: "plain.md", Content: plain})
	if err := p.Transform(m); err != nil {
		t.Fatalf("Transform: %v", err)
	}
	if got := m.Posts()[0].Content; got != plain {
		t.Fatalf("plain content changed: %q", got)
	}
}
