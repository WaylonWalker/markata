package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestOEmbedResolver_Resolve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//nolint:errcheck // test helper
		w.Write([]byte(`{"type":"link","version":"1.0","title":"Test Title","provider_name":"Test Provider","thumbnail_url":"https://example.com/image.jpg"}`))
	}))
	defer server.Close()

	provider := oembedProvider{
		Name:           "test",
		Endpoint:       server.URL,
		URLPrefixes:    []string{"https://example.com/"},
		SupportsFormat: false,
	}

	resolver := newOEmbedResolverWithProviders(models.NewEmbedsConfig(), server.Client(), []oembedProvider{provider})

	response, matched, err := resolver.Resolve("https://example.com/post")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if !matched {
		t.Fatal("expected provider match")
	}
	if response == nil {
		t.Fatal("expected response")
	}
	if response.Title != "Test Title" {
		t.Errorf("expected title, got %s", response.Title)
	}
}

func TestOEmbedResolver_DisabledProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		//nolint:errcheck // test helper
		w.Write([]byte(`{"type":"link","version":"1.0","title":"Test Title"}`))
	}))
	defer server.Close()

	provider := oembedProvider{
		Name:           "test",
		Endpoint:       server.URL,
		URLPrefixes:    []string{"https://example.com/"},
		SupportsFormat: false,
	}

	config := models.NewEmbedsConfig()
	config.OEmbedProviders = map[string]models.OEmbedProviderConfig{
		"test": {Enabled: false},
	}

	resolver := newOEmbedResolverWithProviders(config, server.Client(), []oembedProvider{provider})

	response, matched, err := resolver.Resolve("https://example.com/post")
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if !matched {
		t.Fatal("expected provider match")
	}
	if response != nil {
		t.Fatal("expected nil response when provider disabled")
	}
}

func TestOEmbedResolver_BuildEndpointAddsFormat(t *testing.T) {
	provider := oembedProvider{
		Name:           "test",
		Endpoint:       "https://example.com/oembed",
		URLPrefixes:    []string{"https://example.com/"},
		SupportsFormat: true,
	}

	resolver := newOEmbedResolverWithProviders(models.NewEmbedsConfig(), http.DefaultClient, []oembedProvider{provider})

	endpoint, err := resolver.buildEndpoint(&provider, "https://example.com/post")
	if err != nil {
		t.Fatalf("buildEndpoint failed: %v", err)
	}

	if want := "format=json"; !containsString(endpoint, want) {
		t.Errorf("expected format param, got %s", endpoint)
	}

	if want := fmt.Sprintf("url=%s", url.QueryEscape("https://example.com/post")); !containsString(endpoint, want) {
		t.Errorf("expected encoded url param, got %s", endpoint)
	}
}

func TestOEmbedResolver_DiscoverEndpoint(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		//nolint:errcheck // test helper
		w.Write([]byte(`<!doctype html>
<html>
<head>
  <link rel="alternate" type="application/json+oembed" href="https://example.com/oembed?url=https%3A%2F%2Fexample.com%2Fpost">
</head>
<body></body>
</html>`))
	}))
	defer server.Close()

	config := models.NewEmbedsConfig()
	config.OEmbedAutoDiscover = true

	resolver := newOEmbedResolverWithProviders(config, server.Client(), []oembedProvider{})

	endpoint, err := resolver.discoverEndpoint(server.URL)
	if err != nil {
		t.Fatalf("discoverEndpoint failed: %v", err)
	}
	if endpoint == "" {
		t.Fatal("expected endpoint")
	}
}

func TestRenderGistCodeMarkdown(t *testing.T) {
	md := newEmbedMarkdownRenderer(map[string]interface{}{})
	resolver := newOEmbedResolver(models.NewEmbedsConfig(), http.DefaultClient)
	resolver.setMarkdownRenderer(md)

	html, err := renderGistCodeMarkdown(resolver, "go", "package main\n")
	if err != nil {
		t.Fatalf("renderGistCodeMarkdown failed: %v", err)
	}
	if !containsString(html, `class="highlight"`) && !containsString(html, `class="chroma"`) {
		t.Fatalf("expected highlight/chroma wrapper, got: %s", html)
	}
}

func TestExtractGistID(t *testing.T) {
	got, err := extractGistID("https://gist.github.com/user/abcd1234")
	if err != nil {
		t.Fatalf("extractGistID failed: %v", err)
	}
	if got != "abcd1234" {
		t.Fatalf("expected gist id, got %s", got)
	}
}

func TestRenderGistScriptEmbedHTML(t *testing.T) {
	html := renderGistScriptEmbedHTML("https://gist.github.com/user/abcd", "example.go")
	if !containsString(html, `gist.github.com/user/abcd.js?file=example.go`) {
		t.Fatalf("expected gist script embed, got: %s", html)
	}
}
