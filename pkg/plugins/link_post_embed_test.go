package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLinkPostDestination(t *testing.T) {
	tests := []struct {
		name     string
		template string
		fields   map[string]interface{}
		private  bool
		want     string
	}{
		{name: "link", template: "link", fields: map[string]interface{}{"link": "https://example.com/post"}, want: "https://example.com/post"},
		{name: "url alias", template: "bookmark", fields: map[string]interface{}{"url": "https://example.com/bookmark"}, want: "https://example.com/bookmark"},
		{name: "link wins", template: "stars", fields: map[string]interface{}{"link": "https://example.com/link", "url": "https://example.com/url"}, want: "https://example.com/link"},
		{name: "invalid link falls back", template: "til", fields: map[string]interface{}{"link": "javascript:alert(1)", "url": "https://example.com/url"}, want: "https://example.com/url"},
		{name: "normal article", template: "blog-post", fields: map[string]interface{}{"link": "https://example.com/article"}},
		{name: "relative source", template: "link", fields: map[string]interface{}{"link": "/local"}},
		{name: "unsafe scheme", template: "link", fields: map[string]interface{}{"url": "javascript:alert(1)"}},
		{name: "private link", template: "link", fields: map[string]interface{}{"url": "https://example.com/secret"}, private: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			post := &models.Post{Template: tt.template, Extra: tt.fields, Private: tt.private}
			if got := linkPostDestination(post); got != tt.want {
				t.Fatalf("linkPostDestination() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLinkPostEmbedCacheHashTracksFrontmatterDestination(t *testing.T) {
	post := &models.Post{Template: "link", Content: "Authored note", Extra: map[string]interface{}{"url": "https://example.com/first"}}
	first := linkPostEmbedCacheHash(post, "config")
	post.Extra["url"] = "https://example.com/second"
	second := linkPostEmbedCacheHash(post, "config")
	if first == second {
		t.Fatal("changing the destination must invalidate the embed transform cache")
	}
	post.Template = "post"
	article := linkPostEmbedCacheHash(post, "config")
	if article == second {
		t.Fatal("changing from link post to article must invalidate the embed transform cache")
	}
}

func TestLinkPostDestinationEmbedsWithoutOverwritingAuthoredText(t *testing.T) {
	p := NewEmbedsPlugin()
	p.config.FetchExternal = false
	p.config.OEmbedEnabled = false

	post := &models.Post{
		Template: "link",
		Content:  "My commentary stays here.",
		Extra:    map[string]interface{}{"link": "https://example.com/interesting-post"},
	}
	out := p.embedLinkPostDestination(post, post.Content)
	if !strings.Contains(out, `href="https://example.com/interesting-post"`) {
		t.Fatalf("missing source link in embed: %s", out)
	}
	if !strings.Contains(out, "My commentary stays here.") {
		t.Fatalf("missing commentary: %s", out)
	}
	if strings.Index(out, "My commentary") < strings.Index(out, "embed-card") {
		t.Fatalf("embed should precede commentary: %s", out)
	}

	post.Content = ""
	if out := p.embedLinkPostDestination(post, ""); !strings.Contains(out, "embed-card") {
		t.Fatalf("link post without body must still generate preview: %s", out)
	}
}

func TestLinkPostEmbedAvoidsDuplicateAuthoredEmbeds(t *testing.T) {
	const source = "https://example.com/video"
	p := NewEmbedsPlugin()
	p.config.FetchExternal = false
	p.config.OEmbedEnabled = false
	for _, authored := range []string{
		"![embed](" + source + ")",
		"[!embed](" + source + "|card)",
		"![[" + source + "]]",
		source,
		"<" + source + ">",
	} {
		post := &models.Post{Template: "link", Content: authored, Extra: map[string]interface{}{"link": source}}
		if got := p.embedLinkPostDestination(post, authored); got != authored {
			t.Errorf("authored embed was duplicated for %q: %s", authored, got)
		}
	}
}

func TestLinkPostYouTubeDestinationUsesExistingLiteEmbed(t *testing.T) {
	p := NewEmbedsPlugin()
	p.config.FetchExternal = false
	p.config.OEmbedEnabled = false

	post := &models.Post{
		Template: "link",
		Extra:    map[string]interface{}{"url": "https://www.youtube.com/watch?v=R_uf5OfMGio&t=3327s"},
	}
	got := p.embedLinkPostDestination(post, "")
	if !containsLiteYouTubeEmbed(got) {
		t.Fatalf("YouTube link posts should reuse Lite YouTube embedding: %s", got)
	}
	if !strings.Contains(got, "R_uf5OfMGio") {
		t.Fatalf("YouTube video ID missing: %s", got)
	}
}
