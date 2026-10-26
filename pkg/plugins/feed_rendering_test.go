package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestSplitFeedRenderablePosts_PartitionsVisibleAndRenderablePosts(t *testing.T) {
	t.Parallel()

	title := "Visible"
	posts := []*models.Post{
		{Slug: "visible", Title: &title, Published: true, Content: "body", ArticleHTML: "<p>body</p>"},
		{Slug: "title-only", Title: &title, Published: true},
		{Slug: "draft", Title: &title, Draft: true, Published: false, Content: "draft"},
		{Slug: "private", Title: &title, Private: true, Content: "private", ArticleHTML: `<div class="encrypted-content" data-encrypted="ciphertext">locked</div>`},
	}

	pagePosts, outputPosts := splitFeedRenderablePosts(posts, false)
	if got := len(pagePosts); got != 2 {
		t.Fatalf("pagePosts len = %d, want 2", got)
	}
	if got := len(outputPosts); got != 1 {
		t.Fatalf("outputPosts len = %d, want 1", got)
	}
	if pagePosts[0].Slug != "visible" || pagePosts[1].Slug != "title-only" {
		t.Fatalf("unexpected pagePosts order: %#v", []string{pagePosts[0].Slug, pagePosts[1].Slug})
	}
	if outputPosts[0].Slug != "visible" {
		t.Fatalf("unexpected outputPosts: %q", outputPosts[0].Slug)
	}

	pagePosts, outputPosts = splitFeedRenderablePosts(posts, true)
	if got := len(pagePosts); got != 3 {
		t.Fatalf("pagePosts len with private = %d, want 3", got)
	}
	if got := len(outputPosts); got != 2 {
		t.Fatalf("outputPosts len with private = %d, want 2", got)
	}
}

func TestSafeFeedPost_ScrubsPrivateContentAndSecrets(t *testing.T) {
	title := "Private title"
	description := "Safe description"
	post := &models.Post{
		Path:         "private.md",
		Slug:         "private",
		Href:         "/private/",
		Title:        &title,
		Description:  &description,
		Private:      true,
		Content:      "private body",
		ArticleHTML:  `<div class="encrypted-content" data-encrypted="ciphertext" data-key-name="private-key">locked</div>`,
		SecretKey:    "private-key",
		InputHash:    "private-input-hash",
		PrevNextFeed: "private-feed",
		Tags:         []string{"private"},
		Extra: map[string]interface{}{
			"_title_explicit":       true,
			"_description_explicit": true,
			"has_encrypted_content": true,
			"encryption_key_name":   "private-key",
			"image":                 "https://example.test/private.webp",
			"secret_value":          "private secret",
			"category":              "notes",
		},
	}

	safe := safeFeedPost(post)
	if safe == post {
		t.Fatal("private feed post should be copied")
	}
	if safe.Content != "" || safe.HTML != "" || safe.SecretKey != "" || safe.InputHash != "" || safe.PrevNextFeed != "" {
		t.Fatalf("private content or key was retained: %#v", safe)
	}
	if safe.ArticleHTML == "" || !strings.Contains(safe.ArticleHTML, `data-encrypted="ciphertext"`) {
		t.Fatalf("encrypted article HTML was not retained: %q", safe.ArticleHTML)
	}
	if strings.Contains(safe.ArticleHTML, "data-key-name") || strings.Contains(safe.ArticleHTML, "private-key") {
		t.Fatalf("feed wrapper exposed the encryption key name: %q", safe.ArticleHTML)
	}
	if post.ArticleHTML == "" || !strings.Contains(post.ArticleHTML, `data-key-name="private-key"`) {
		t.Fatal("safe projection mutated the source post")
	}
	if safe.Title == nil || *safe.Title != title || safe.Description == nil || *safe.Description != description {
		t.Fatalf("explicit metadata was not retained: %#v", safe)
	}
	if _, ok := safe.Extra["secret_value"]; ok {
		t.Fatal("arbitrary private Extra value was retained")
	}
	if _, ok := safe.Extra["encryption_key_name"]; ok {
		t.Fatal("encryption key name was retained")
	}
	if _, ok := safe.Extra["image"]; ok {
		t.Fatal("private media was retained")
	}
	if safe.Extra["category"] != "notes" {
		t.Fatal("safe authored metadata was not retained")
	}
}

func TestSafeFeedPost_DropsUnencryptedPrivateBody(t *testing.T) {
	post := &models.Post{Private: true, Content: "private body", ArticleHTML: "<p>private body</p>"}
	safe := safeFeedPost(post)
	if safe.Content != "" || safe.ArticleHTML != "" {
		t.Fatalf("unencrypted private content was retained: %#v", safe)
	}
}

func TestSafeFeedPost_RemovesNestedKeyNames(t *testing.T) {
	post := &models.Post{
		Private: true,
		ArticleHTML: `<div data-key-name="family&quot;archive" class="encrypted-content" data-encrypted="ciphertext">
  <div data-key-name="nested-key">locked</div>
</div>`,
		Extra: map[string]interface{}{"has_encrypted_content": true},
	}

	safe := safeFeedPost(post)
	if safe.ArticleHTML == "" {
		t.Fatal("encrypted wrapper was dropped")
	}
	if strings.Contains(strings.ToLower(safe.ArticleHTML), "data-key-name") || strings.Contains(safe.ArticleHTML, "family") || strings.Contains(safe.ArticleHTML, "nested-key") {
		t.Fatalf("key name remained in safe wrapper: %q", safe.ArticleHTML)
	}
	if !strings.Contains(safe.ArticleHTML, `data-encrypted="ciphertext"`) {
		t.Fatalf("ciphertext was not retained: %q", safe.ArticleHTML)
	}
	if !strings.Contains(post.ArticleHTML, "family&quot;archive") || !strings.Contains(post.ArticleHTML, "nested-key") {
		t.Fatal("safe projection mutated the source wrapper")
	}
}

func TestSafeFeedPost_DropsAmbiguousEncryptedWrapper(t *testing.T) {
	posts := []string{
		`<span class="encrypted-content" data-encrypted="ciphertext">locked</span>`,
		`<div class="encrypted-content" data-encrypted="ciphertext">one</div><div class="encrypted-content" data-encrypted="ciphertext">two</div>`,
		`<div class="encrypted-content" data-encrypted="">locked</div>`,
	}
	for _, articleHTML := range posts {
		safe := safeFeedPost(&models.Post{
			Private:     true,
			ArticleHTML: articleHTML,
			Extra:       map[string]interface{}{"has_encrypted_content": true},
		})
		if safe.ArticleHTML != "" {
			t.Errorf("ambiguous wrapper %q was retained as %q", articleHTML, safe.ArticleHTML)
		}
	}
}
