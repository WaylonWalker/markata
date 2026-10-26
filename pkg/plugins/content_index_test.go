package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/contentindex"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestContentIndexPlugin_WritesSafeMetadataAndFeeds(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "content-index.json")
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{ContentDir: dir, OutputDir: dir, Extra: map[string]interface{}{"content_index": map[string]interface{}{"enabled": true, "output": output}, "markata_version": "test"}})
	public := models.NewPost("posts/public.md")
	public.Slug, public.Href, public.Published, public.Tags = "public", "/public/", true, []string{"z", "a"}
	public.Extra["image"] = "https://images.example.test/public.webp"
	public.Extra["video"] = "https://videos.example.test/public.mp4"
	public.Extra["bio"] = "Author biography"
	public.Extra["thumbnail"] = "https://images.example.test/thumb.webp"
	public.Extra["cover_image"] = "https://images.example.test/cover.webp"
	public.Extra["og_image"] = "https://images.example.test/og.webp"
	public.Extra["category"] = "notes"
	public.Extra["categories"] = []interface{}{"writing", "personal"}
	public.Author = contentIndexStringPtr("waylon")
	public.Authors = []string{"waylon", "guest"}
	authorAvatar, authorBio := "https://images.example.test/author.webp", "Author biography"
	public.AuthorObjects = []models.Author{{ID: "waylon", Avatar: &authorAvatar, Bio: &authorBio}}
	unpublished := models.NewPost("pages/unpublished.md")
	unpublished.Slug, unpublished.Href = "unpublished", "/unpublished/"
	private := models.NewPost("posts/private.md")
	privateTitle := "Private title"
	private.Private, private.Published = true, true
	private.Slug, private.Href = "private", "/private/"
	private.Title = &privateTitle
	private.Description = contentIndexStringPtr("body-derived private description")
	private.Content = "private body that must not be serialized"
	private.ArticleHTML = "<p>private body that must not be serialized</p>"
	private.SecretKey = "private-key"
	private.Set("_title_explicit", true)
	private.Set("image", "https://images.example.test/private.webp")
	private.Set("video", "https://videos.example.test/private.mp4")
	private.Set("avatar", "https://images.example.test/private-avatar.webp")
	private.Set("bio", "private biography")
	private.Set("thumbnail", "https://images.example.test/private-thumb.webp")
	private.Set("cover", "https://images.example.test/private-cover.webp")
	private.Set("og_image", "https://images.example.test/private-og.webp")
	draft := models.NewPost("posts/draft.md")
	draft.Draft = true
	m.SetPosts([]*models.Post{private, draft, public, unpublished})
	m.SetFeeds([]*lifecycle.Feed{{Name: "blog", IncludePrivate: true, Posts: []*models.Post{public, private}}, {Name: "archive", Posts: []*models.Post{public, private}}, {Name: "draft", Posts: []*models.Post{unpublished}}})
	if err := NewContentIndexPlugin().Write(m); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	index, err := contentindex.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if index.Scope != contentindex.PublicMetadataScope {
		t.Fatalf("scope = %q, want %q", index.Scope, contentindex.PublicMetadataScope)
	}
	if index.DocumentCount != 3 || index.Documents[0].Path != unpublished.Path || index.Documents[1].Path != private.Path || index.Documents[2].Path != public.Path {
		t.Fatalf("document filtering failed: %#v", index)
	}
	if got := index.Documents[2].Feeds; len(got) != 2 || got[0] != "archive" || got[1] != "blog" {
		t.Fatalf("feeds = %v", got)
	}
	if index.Documents[0].Published || len(index.Documents[0].Feeds) != 1 || index.Documents[0].Feeds[0] != "draft" {
		t.Fatalf("unpublished direct page or draft feed was conflated: %#v", index.Documents[0])
	}
	privateDocument := index.Documents[1]
	if !privateDocument.Private || privateDocument.Title == nil || *privateDocument.Title != privateTitle || privateDocument.Description != nil || len(privateDocument.Feeds) != 1 || privateDocument.Feeds[0] != "blog" {
		t.Fatalf("private metadata = %#v", privateDocument)
	}
	if privateDocument.Image != nil || privateDocument.Video != nil || privateDocument.Bio != nil || privateDocument.Thumbnail != nil || privateDocument.Cover != nil || privateDocument.OGImage != nil || privateDocument.Avatar == nil || *privateDocument.Avatar != "https://images.example.test/private-avatar.webp" {
		t.Fatalf("private sensitive metadata = %#v", privateDocument)
	}
	if strings.Contains(string(data), "private body that must not be serialized") || strings.Contains(string(data), "private-key") {
		t.Fatalf("private content or key leaked into index: %s", data)
	}
	document := index.Documents[2]
	if document.Image == nil || *document.Image != "https://images.example.test/public.webp" || document.Video == nil || document.Avatar == nil || *document.Avatar != authorAvatar || document.Bio == nil || *document.Bio != authorBio || document.Thumbnail == nil || document.Cover == nil || *document.Cover != "https://images.example.test/cover.webp" || document.OGImage == nil {
		t.Fatalf("media metadata = %#v", document)
	}
	if document.Author == nil || *document.Author != "waylon" || len(document.Authors) != 2 || document.Category == nil || *document.Category != "notes" || len(document.Categories) != 2 {
		t.Fatalf("author/category metadata = %#v", document)
	}
}

func TestContentIndexPlugin_ExplicitV1RemainsPublicOnly(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "content-index-v1.json")
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{ContentDir: dir, OutputDir: dir, Extra: map[string]interface{}{
		"content_index": map[string]interface{}{"enabled": true, "output": output, "schema_version": 1},
	}})
	public := models.NewPost("public.md")
	public.Slug, public.Href, public.Published = "public", "/public/", true
	private := models.NewPost("private.md")
	private.Slug, private.Href, private.Published, private.Private = "private", "/private/", true, true
	m.SetPosts([]*models.Post{private, public})

	if err := NewContentIndexPlugin().Write(m); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	index, err := contentindex.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	if index.SchemaVersion != 1 || index.Scope != contentindex.PublicScope || len(index.Documents) != 1 || index.Documents[0].Path != public.Path {
		t.Fatalf("explicit v1 output changed compatibility behavior: %#v", index)
	}
}

func TestDocumentFromPost_PrivateExplicitMarkersRequireTrue(t *testing.T) {
	markerValues := []struct {
		name  string
		value interface{}
		set   bool
	}{
		{name: "missing"},
		{name: "false", value: false, set: true},
		{name: "nil", value: nil, set: true},
		{name: "string", value: "true", set: true},
		{name: "integer", value: 1, set: true},
		{name: "true", value: true, set: true},
	}
	for _, marker := range markerValues {
		t.Run(marker.name, func(t *testing.T) {
			title := "Private title"
			description := "Private description"
			post := models.NewPost("private.md")
			post.Private = true
			post.Title = &title
			post.Description = &description
			post.Extra = make(map[string]interface{})
			if marker.set {
				post.Set("_title_explicit", marker.value)
				post.Set("_description_explicit", marker.value)
			}

			document := documentFromPost(post, nil)
			if marker.name == "true" {
				if document.Title == nil || document.Description == nil {
					t.Fatalf("true markers did not preserve explicit metadata: %#v", document)
				}
				return
			}
			if document.Title != nil || document.TitleText != nil || document.Description != nil {
				t.Fatalf("non-boolean-true markers exposed metadata: %#v", document)
			}
		})
	}
}

func contentIndexStringPtr(value string) *string { return &value }

func TestContentIndexPlugin_DisabledByDefault(t *testing.T) {
	dir := t.TempDir()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{OutputDir: dir, Extra: map[string]interface{}{}})
	if err := NewContentIndexPlugin().Write(m); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "content-index.json")); !os.IsNotExist(err) {
		t.Fatalf("index was written while disabled: %v", err)
	}
}

func TestContentIndexPlugin_RejectsOutputTraversal(t *testing.T) {
	dir := t.TempDir()
	m := lifecycle.NewManager()
	m.SetConfig(&lifecycle.Config{OutputDir: dir, Extra: map[string]interface{}{"content_index": map[string]interface{}{"enabled": true, "output": "../escape.json"}}})
	if err := NewContentIndexPlugin().Write(m); err == nil {
		t.Fatal("expected output traversal error")
	}
}

func TestContentIndexPlugin_RejectsInvalidConfiguration(t *testing.T) {
	for name, config := range map[string]interface{}{
		"invalid table":      true,
		"invalid enabled":    map[string]interface{}{"enabled": "true"},
		"invalid output":     map[string]interface{}{"enabled": true, "output": 42},
		"fractional version": map[string]interface{}{"enabled": true, "schema_version": 1.5},
	} {
		t.Run(name, func(t *testing.T) {
			m := lifecycle.NewManager()
			m.SetConfig(&lifecycle.Config{Extra: map[string]interface{}{"content_index": config}})
			if err := NewContentIndexPlugin().Write(m); err == nil {
				t.Fatal("expected configuration error")
			}
		})
	}
}
