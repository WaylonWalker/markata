package plugins

import (
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestGenerateAtom_ArchiveMetadataAndFallbacks(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"url":           "https://example.com",
		"title":         "Example Site",
		"description":   "Site description",
		"language":      "en",
		"author_url":    "https://example.com/about/",
		"models_config": &models.Config{},
	}

	feed := &lifecycle.Feed{
		Title:       "Blog Archive",
		Description: "Blog posts",
		Path:        "blog/archive",
		Posts: []*models.Post{{
			Slug: "one",
			Href: "/one/",
		}},
	}

	atom, err := GenerateAtom(feed, config)
	if err != nil {
		t.Fatalf("GenerateAtom() error = %v", err)
	}

	checks := []string{
		`xml:lang="en"`,
		`xmlns:fh="http://purl.org/syndication/history/1.0"`,
		`<subtitle>Blog posts</subtitle>`,
		`<name>Example Site</name>`,
		`<uri>https://example.com/about/</uri>`,
		`<fh:complete></fh:complete>`,
		`href="https://example.com/blog/atom.xml" rel="current" type="application/atom+xml"`,
		`<updated>1970-01-01T00:00:00Z</updated>`,
	}
	for _, want := range checks {
		if !strings.Contains(atom, want) {
			t.Fatalf("expected atom feed to contain %q\n%s", want, atom)
		}
	}
}

func TestGenerateRSS_UsesFeedSpecificMetadata(t *testing.T) {
	date := time.Date(2024, 2, 2, 12, 0, 0, 0, time.UTC)
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"url":             "https://example.com",
		"title":           "Example Site",
		"description":     "Site description",
		"language":        "en-us",
		"managing_editor": "editor@example.com (Editor)",
		"webmaster":       "webmaster@example.com (Webmaster)",
		"copyright":       "Copyright 2026 Example",
		"models_config": &models.Config{
			Authors: models.AuthorsConfig{Authors: map[string]models.Author{
				"waylon": {ID: "waylon", Name: "Waylon", Email: testStringPtr("waylon@example.com")},
			}},
		},
	}

	feed := &lifecycle.Feed{
		Title:       "Blog",
		Description: "Blog posts",
		Path:        "blog",
		Posts: []*models.Post{{
			Slug:      "one",
			Href:      "/one/",
			Title:     testStringPtr("One"),
			Date:      &date,
			Tags:      []string{"go", "feeds"},
			Authors:   []string{"waylon"},
			Published: true,
		}},
	}

	rss, err := GenerateRSS(feed, config)
	if err != nil {
		t.Fatalf("GenerateRSS() error = %v", err)
	}

	checks := []string{
		`<link>https://example.com/blog/</link>`,
		`<description>Blog posts</description>`,
		`<language>en-us</language>`,
		`<managingEditor>editor@example.com (Editor)</managingEditor>`,
		`<webMaster>webmaster@example.com (Webmaster)</webMaster>`,
		`<copyright>Copyright 2026 Example</copyright>`,
		`<generator>markata-go</generator>`,
		`<docs>https://www.rssboard.org/rss-specification</docs>`,
		`<author>waylon@example.com</author>`,
		`<category>go</category>`,
		`<category>feeds</category>`,
	}
	for _, want := range checks {
		if !strings.Contains(rss, want) {
			t.Fatalf("expected rss feed to contain %q\n%s", want, rss)
		}
	}
}

func TestGenerateJSONFeed_UsesFeedMetadata(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"url":         "https://example.com",
		"title":       "Example Site",
		"description": "Site description",
		"author":      "Waylon",
		"author_url":  "https://example.com/about/",
		"language":    "en",
		"models_config": &models.Config{
			SEO: models.SEOConfig{LogoURL: "https://example.com/logo.png"},
		},
	}

	feed := &lifecycle.Feed{
		Title:       "Blog Archive",
		Description: "Blog posts",
		Path:        "blog/archive",
		Posts: []*models.Post{{
			Slug:      "one",
			Href:      "/one/",
			Published: true,
		}},
	}

	jsonFeed, err := GenerateJSONFeed(feed, config)
	if err != nil {
		t.Fatalf("GenerateJSONFeed() error = %v", err)
	}

	checks := []string{
		`"home_page_url": "https://example.com/blog/"`,
		`"description": "Blog posts"`,
		`"language": "en"`,
		`"icon": "https://example.com/logo.png"`,
		`"name": "Waylon"`,
		`"url": "https://example.com/about/"`,
	}
	for _, want := range checks {
		if !strings.Contains(jsonFeed, want) {
			t.Fatalf("expected json feed to contain %q\n%s", want, jsonFeed)
		}
	}
}

func TestGenerateJSONFeed_IncludePrivateUsesEncryptedHTML(t *testing.T) {
	config := lifecycle.NewConfig()
	config.Extra = map[string]interface{}{
		"url":   "https://example.com",
		"title": "Example Site",
		"models_config": &models.Config{Authors: models.AuthorsConfig{Authors: map[string]models.Author{
			"waylon": {
				ID:     "waylon",
				Name:   "Waylon",
				Email:  testStringPtr("private-author@example.com"),
				URL:    testStringPtr("https://private.example.com/profile"),
				Avatar: testStringPtr("https://private.example.com/avatar.webp"),
				Bio:    testStringPtr("private author bio"),
			},
		}}},
	}

	date := time.Date(2024, 2, 2, 12, 0, 0, 0, time.UTC)
	feed := &lifecycle.Feed{
		Title:          "Blog Archive",
		Path:           "blog/archive",
		IncludePrivate: true,
		Posts: []*models.Post{{
			Slug:        "one",
			Href:        "/one/",
			Title:       testStringPtr("One"),
			Published:   true,
			Private:     true,
			Date:        &date,
			Authors:     []string{"waylon"},
			Description: testStringPtr("secret summary"),
			Content:     "secret body",
			ArticleHTML: `<div class="encrypted-content" data-encrypted="ciphertext" data-key-name="private-key">locked</div>`,
			Extra:       map[string]interface{}{"avatar": "https://example.com/explicit-avatar.webp", "image": "https://example.com/private-image.webp", "secret_value": "private secret"},
		}},
	}

	jsonFeed, err := GenerateJSONFeed(feed, config)
	if err != nil {
		t.Fatalf("GenerateJSONFeed() error = %v", err)
	}

	if !strings.Contains(jsonFeed, `"content_html": "\u003cdiv class=\"encrypted-content\" data-encrypted=\"ciphertext\"\u003elocked\u003c/div\u003e"`) {
		t.Fatalf("expected json feed to include encrypted HTML\n%s", jsonFeed)
	}
	if strings.Contains(jsonFeed, "data-key-name") || strings.Contains(jsonFeed, "private-key") || strings.Contains(jsonFeed, "secret summary") || strings.Contains(jsonFeed, "secret body") || strings.Contains(jsonFeed, "private-image.webp") || strings.Contains(jsonFeed, "private secret") || strings.Contains(jsonFeed, "private.example.com") || strings.Contains(jsonFeed, `"content_text"`) {
		t.Fatalf("json feed should not expose private plaintext\n%s", jsonFeed)
	}
	if !strings.Contains(jsonFeed, "explicit-avatar.webp") {
		t.Fatalf("json feed should retain explicitly authored avatar metadata\n%s", jsonFeed)
	}

	rss, err := GenerateRSS(feed, config)
	if err != nil {
		t.Fatalf("GenerateRSS() error = %v", err)
	}
	atom, err := GenerateAtom(feed, config)
	if err != nil {
		t.Fatalf("GenerateAtom() error = %v", err)
	}
	for name, output := range map[string]string{"rss": rss, "atom": atom} {
		if !strings.Contains(output, "encrypted-content") || !strings.Contains(output, "ciphertext") {
			t.Errorf("%s feed should retain encrypted content\n%s", name, output)
		}
		if strings.Contains(output, "data-key-name") || strings.Contains(output, "private-key") || strings.Contains(output, "secret summary") || strings.Contains(output, "secret body") || strings.Contains(output, "private-image.webp") || strings.Contains(output, "private secret") || strings.Contains(output, "private.example.com") {
			t.Errorf("%s feed should not expose private plaintext or key names\n%s", name, output)
		}
	}
}

func testStringPtr(s string) *string {
	return &s
}
