package plugins

import (
	"io"
	"net/http"
	"strings"
	"sync/atomic"
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

type countingMentionTransport func(*http.Request) (*http.Response, error)

func (fn countingMentionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

// Unlike selector-only tests, this proves Transform's real metadata-fetch path
// sends zero requests for inert references and one request per used domain.
// All requests are intercepted; CI never contacts an external network.
func TestMentionsTransform_MetadataHTTPCounts(t *testing.T) {
	tests := []struct {
		name, content string
		want          int64
	}{
		{"plain post", "No relevant contacts", 0},
		{"unknown mention", "Hello @unknown", 0},
		{"only fenced code", "```md\n@alice\n```", 0},
		{"nonchat header", "!!! note @alice\n    Other text", 0},
		{"known handle", "Hello @alice", 1},
		{"alias", "Hello @ally", 1},
		{"repeat domain", "Hello @alice and @ally again", 1},
		{"chat title", "!!! chat \"@alice\"\n    Hello", 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int64
			p := NewMentionsPlugin()
			p.metadataHTTPClient = &http.Client{
				Transport: countingMentionTransport(func(req *http.Request) (*http.Response, error) {
					calls.Add(1)
					if req.URL.Scheme != "https" || req.URL.Host != "alice.example" {
						t.Errorf("unexpected request URL: %s", req.URL)
					}
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader("<html><head><title>Alice</title></head></html>")),
						Request:    req,
					}, nil
				}),
			}
			config := models.NewMentionsConfig()
			config.CacheDir = t.TempDir()
			active := true
			m := lifecycle.NewManager()
			m.Config().Extra = map[string]interface{}{
				"mentions": config,
				"blogroll": models.BlogrollConfig{
					Enabled: true,
					Feeds: []models.ExternalFeedConfig{{
						URL: "https://alice.example/feed.xml",
						SiteURL: "https://alice.example",
						Handle: "alice",
						Aliases: []string{"ally"},
						Active: &active,
					}},
				},
			}
			m.AddPost(&models.Post{Path: "example.md", Content: tc.content})
			if err := p.Transform(m); err != nil {
				t.Fatal(err)
			}
			if got := calls.Load(); got != tc.want {
				t.Fatalf("metadata HTTP requests: got %d, want %d for %q", got, tc.want, tc.content)
			}
		})
	}
}
