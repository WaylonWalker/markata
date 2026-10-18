package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestBuildPostCopyPayloads(t *testing.T) {
	post := &models.Post{
		Slug:        "hello-world",
		Href:        "/hello-world/",
		Title:       stringPtr("Hello World"),
		Description: stringPtr("A friendly introduction"),
		Content:     "## Intro\n\n```python\nprint('hello')\n```\n",
		ArticleHTML: "<h2 id=\"intro\">Intro <a href=\"#intro\" class=\"heading-anchor\">#</a></h2><pre><code class=\"language-python\">print('hello')\n</code></pre>",
	}

	payloads := buildPostCopyPayloads(post, nil, "https://example.com")
	if payloads.URL != "https://example.com/hello-world/" {
		t.Fatalf("unexpected url payload: %q", payloads.URL)
	}
	if payloads.MarkdownURL != "https://example.com/hello-world.md" {
		t.Fatalf("unexpected markdown url payload: %q", payloads.MarkdownURL)
	}
	if payloads.TextURL != "https://example.com/hello-world.txt" {
		t.Fatalf("unexpected text url payload: %q", payloads.TextURL)
	}
	if payloads.ANSICurl != "" {
		t.Fatalf("unexpected ansi curl payload: %q", payloads.ANSICurl)
	}
	if payloads.Markdown != "" || payloads.Text != "" {
		t.Fatalf("published formats should not duplicate full post content in HTML, got markdown=%q text=%q", payloads.Markdown, payloads.Text)
	}
}

func TestBuildPostCopyPayloads_HidesDisabledFormatRoutes(t *testing.T) {
	htmlEnabled := true
	post := &models.Post{
		Slug:        "hello-world",
		Href:        "/hello-world/",
		Title:       stringPtr("Hello World"),
		Content:     "hello",
		ArticleHTML: "<p>hello</p>",
	}
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: false,
				Text:     false,
				ANSI:     false,
			},
		},
	}

	payloads := buildPostCopyPayloads(post, config, "https://example.com")
	if payloads.MarkdownURL != "" {
		t.Fatalf("expected markdown route to be hidden, got %q", payloads.MarkdownURL)
	}
	if payloads.TextURL != "" {
		t.Fatalf("expected text route to be hidden, got %q", payloads.TextURL)
	}
	if payloads.ANSICurl != "" {
		t.Fatalf("expected ansi route to be hidden, got %q", payloads.ANSICurl)
	}
	if !strings.Contains(payloads.Markdown, "# Hello World") || !strings.Contains(payloads.Text, "Source: https://example.com/hello-world/") {
		t.Fatalf("disabled formats need inline fallback, markdown=%q text=%q", payloads.Markdown, payloads.Text)
	}
}

func TestBuildPostCopyPayloads_ShowsOnlyEnabledRoutes(t *testing.T) {
	htmlEnabled := true
	post := &models.Post{
		Slug:        "hello-world",
		Href:        "/hello-world/",
		Title:       stringPtr("Hello World"),
		Content:     "hello",
		ArticleHTML: "<p>hello</p>",
	}
	config := &lifecycle.Config{
		Extra: map[string]interface{}{
			"post_formats": models.PostFormatsConfig{
				HTML:     &htmlEnabled,
				Markdown: true,
				Text:     false,
				ANSI:     true,
			},
		},
	}

	payloads := buildPostCopyPayloads(post, config, "https://example.com")
	if payloads.MarkdownURL != "https://example.com/hello-world.md" {
		t.Fatalf("unexpected markdown route: %q", payloads.MarkdownURL)
	}
	if payloads.TextURL != "" {
		t.Fatalf("expected text route to be hidden, got %q", payloads.TextURL)
	}
	if payloads.ANSICurl != "curl https://example.com/hello-world.ansi" {
		t.Fatalf("unexpected ansi route: %q", payloads.ANSICurl)
	}
}

func stringPtr(value string) *string {
	return &value
}
