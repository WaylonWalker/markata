package plugins

import (
	"reflect"
	"sort"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestReferencedMentionHandles(t *testing.T) {
	alice := &mentionEntry{Handle: "alice", SiteURL: "https://alice.example"}
	bob := &mentionEntry{Handle: "bob", SiteURL: "https://bob.example"}
	handles := map[string]*mentionEntry{
		"alice": alice,
		"ally":  alice, // Alias shares the same metadata domain.
		"bob":   bob,
	}

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{name: "plain page has no metadata requests", content: "A plain post with no references"},
		{name: "unknown mention does not fetch", content: "Hello @unknown and ~missing"},
		{name: "known handle", content: "Thanks @Alice!", want: []string{"alice"}},
		{name: "alias and trailing punctuation", content: "Thanks @ally. and @BOB!", want: []string{"ally", "bob"}},
		{name: "chat admonition title", content: "!!! chat \"@alice\"\n    Hello!", want: []string{"alice"}},
		{name: "fenced code is ignored", content: "```md\n@bob\n```\nOutside @alice", want: []string{"alice"}},
		{name: "email address is not a mention", content: "Write to user@bob.example"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMap := referencedMentionHandles([]*models.Post{{Content: tt.content}}, handles)
			got := make([]string, 0, len(gotMap))
			for key := range gotMap {
				got = append(got, key)
			}
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) && !(len(got) == 0 && len(tt.want) == 0) {
				t.Fatalf("referenced handles = %v, want %v", got, tt.want)
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
