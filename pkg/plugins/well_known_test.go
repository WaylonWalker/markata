package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

type linksDomainFixture struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
	Links  []struct {
		SourceURL string `json:"sourceUrl"`
		TargetURL string `json:"targetUrl"`
	} `json:"links"`
}

type internalTargetFixture struct {
	TargetURL string `json:"targetUrl"`
	Count     int    `json:"count"`
	Links     []struct {
		SourceURL string `json:"sourceUrl"`
		TargetURL string `json:"targetUrl"`
	} `json:"links"`
}

func TestWellKnownPlugin_Write_DefaultEntries(t *testing.T) {
	outputDir := t.TempDir()
	wellKnownConfig := models.NewWellKnownConfig()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"url":         "https://example.com",
			"title":       "Example Site",
			"description": "Example description",
			"author":      "Jane Doe",
			"well_known":  wellKnownConfig,
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	plugin := NewWellKnownPlugin()
	plugin.now = func() time.Time {
		return time.Date(2026, time.February, 4, 12, 34, 56, 0, time.UTC)
	}

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	expected := []string{
		".well-known/host-meta",
		".well-known/host-meta.json",
		".well-known/webfinger",
		".well-known/nodeinfo",
		".well-known/time",
		".well-known/links",
		".well-known/internal-links",
		"nodeinfo/2.0",
		"external-links/index.html",
		"internal-links/index.html",
	}

	for _, rel := range expected {
		path := filepath.Join(outputDir, filepath.FromSlash(rel))
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}

	content, err := os.ReadFile(filepath.Join(outputDir, filepath.FromSlash(".well-known/time")))
	if err != nil {
		t.Fatalf("reading time file: %v", err)
	}
	if string(content) != "2026-02-04T12:34:56Z\n" {
		t.Fatalf("time content = %q, want %q", string(content), "2026-02-04T12:34:56Z\n")
	}
}

func TestWellKnownPlugin_Write_ExternalLinksPage(t *testing.T) {
	outputDir := t.TempDir()
	enabled := true
	wellKnownConfig := models.WellKnownConfig{
		Enabled:      &enabled,
		AutoGenerate: []string{"links"},
	}

	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"url":        "https://example.com",
			"title":      "Example Site",
			"well_known": wellKnownConfig,
		},
	}

	posts := []*models.Post{
		{
			Href: "/alpha/",
			Outlinks: []*models.Link{
				{SourceURL: "https://example.com/alpha/", TargetURL: "https://go.dev/doc", IsInternal: false},
				{SourceURL: "https://example.com/alpha/", TargetURL: "https://example.com/about/", IsInternal: true},
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)
	m.SetPosts(posts)

	plugin := NewWellKnownPlugin()
	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	pagePath := filepath.Join(outputDir, filepath.FromSlash("external-links/index.html"))
	content, err := os.ReadFile(pagePath)
	if err != nil {
		t.Fatalf("reading external links page: %v", err)
	}

	htmlBody := string(content)
	if !strings.Contains(htmlBody, "External Links") {
		t.Fatalf("expected title on external links page")
	}
	if !strings.Contains(htmlBody, "https://go.dev/doc") {
		t.Fatalf("expected outbound link on external links page")
	}
	if !strings.Contains(htmlBody, "/.well-known/links") {
		t.Fatalf("expected well-known links reference on external links page")
	}
	if !strings.Contains(htmlBody, "external-links-bar-fill") {
		t.Fatalf("expected bar graph rows on external links page")
	}

	internalPagePath := filepath.Join(outputDir, filepath.FromSlash("internal-links/index.html"))
	internalContent, err := os.ReadFile(internalPagePath)
	if err != nil {
		t.Fatalf("reading internal links page: %v", err)
	}
	internalHTML := string(internalContent)
	if !strings.Contains(internalHTML, "Internal Links") {
		t.Fatalf("expected title on internal links page")
	}
	if !strings.Contains(internalHTML, "/.well-known/internal-links") {
		t.Fatalf("expected well-known internal links reference on internal links page")
	}
}

func TestWellKnownPlugin_Write_OptionalEntriesOnly(t *testing.T) {
	outputDir := t.TempDir()
	enabled := true
	wellKnownConfig := models.WellKnownConfig{
		Enabled:         &enabled,
		AutoGenerate:    []string{},
		SSHFingerprint:  "SHA256:abcdef",
		KeybaseUsername: "alice",
	}
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"url":        "https://example.com",
			"well_known": wellKnownConfig,
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	plugin := NewWellKnownPlugin()
	plugin.now = func() time.Time { return time.Date(2026, time.February, 4, 0, 0, 0, 0, time.UTC) }

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	sshfpPath := filepath.Join(outputDir, filepath.FromSlash(".well-known/sshfp"))
	if _, err := os.Stat(sshfpPath); err != nil {
		t.Fatalf("expected %s to exist: %v", sshfpPath, err)
	}

	keybasePath := filepath.Join(outputDir, filepath.FromSlash(".well-known/keybase.txt"))
	if _, err := os.Stat(keybasePath); err != nil {
		t.Fatalf("expected %s to exist: %v", keybasePath, err)
	}

	if _, err := os.Stat(filepath.Join(outputDir, filepath.FromSlash(".well-known/host-meta"))); err == nil {
		t.Fatalf("did not expect host-meta to be generated when auto_generate is empty")
	}
}

func TestWellKnownPlugin_Write_WithAuthorImage(t *testing.T) {
	outputDir := t.TempDir()
	wellKnownConfig := models.NewWellKnownConfig()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"url":        "https://example.com",
			"title":      "Example Site",
			"author":     "Jane Doe",
			"well_known": wellKnownConfig,
			"seo": map[string]interface{}{
				"author_image": "/images/avatar.png",
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	plugin := NewWellKnownPlugin()
	plugin.now = func() time.Time { return time.Date(2026, time.February, 4, 12, 0, 0, 0, time.UTC) }

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check that avatar endpoint was created
	avatarPath := filepath.Join(outputDir, filepath.FromSlash(".well-known/avatar"))
	if _, err := os.Stat(avatarPath); err != nil {
		t.Fatalf("expected %s to exist: %v", avatarPath, err)
	}

	// Check avatar content contains redirect to image
	content, err := os.ReadFile(avatarPath)
	if err != nil {
		t.Fatalf("reading avatar file: %v", err)
	}
	if !strings.Contains(string(content), "https://example.com/images/avatar.png") {
		t.Errorf("avatar content should contain image URL, got: %s", string(content))
	}

	// Check webfinger contains avatar link
	webfingerPath := filepath.Join(outputDir, filepath.FromSlash(".well-known/webfinger"))
	webfingerContent, err := os.ReadFile(webfingerPath)
	if err != nil {
		t.Fatalf("reading webfinger file: %v", err)
	}
	if !strings.Contains(string(webfingerContent), "http://webfinger.net/rel/avatar") {
		t.Errorf("webfinger should contain avatar rel, got: %s", string(webfingerContent))
	}
	if !strings.Contains(string(webfingerContent), "https://example.com/images/avatar.png") {
		t.Errorf("webfinger should contain avatar URL, got: %s", string(webfingerContent))
	}
}

func TestWellKnownPlugin_Write_WithoutAuthorImage(t *testing.T) {
	outputDir := t.TempDir()
	wellKnownConfig := models.NewWellKnownConfig()
	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"url":        "https://example.com",
			"title":      "Example Site",
			"author":     "Jane Doe",
			"well_known": wellKnownConfig,
			// No SEO config with author_image
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)

	plugin := NewWellKnownPlugin()
	plugin.now = func() time.Time { return time.Date(2026, time.February, 4, 12, 0, 0, 0, time.UTC) }

	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Check that avatar endpoint was NOT created (no author image)
	avatarPath := filepath.Join(outputDir, filepath.FromSlash(".well-known/avatar"))
	if _, err := os.Stat(avatarPath); err == nil {
		t.Fatalf("did not expect %s to exist when no author image is configured", avatarPath)
	}

	// Check webfinger does NOT contain avatar link
	webfingerPath := filepath.Join(outputDir, filepath.FromSlash(".well-known/webfinger"))
	webfingerContent, err := os.ReadFile(webfingerPath)
	if err != nil {
		t.Fatalf("reading webfinger file: %v", err)
	}
	if strings.Contains(string(webfingerContent), "http://webfinger.net/rel/avatar") {
		t.Errorf("webfinger should NOT contain avatar rel when no image configured, got: %s", string(webfingerContent))
	}
}

func TestGetAuthorImageURL(t *testing.T) {
	tests := []struct {
		name     string
		config   *lifecycle.Config
		siteURL  string
		expected string
	}{
		{
			name:     "nil config",
			config:   nil,
			siteURL:  "https://example.com",
			expected: "",
		},
		{
			name: "relative author image",
			config: &lifecycle.Config{
				Extra: map[string]interface{}{
					"seo": map[string]interface{}{
						"author_image": "/images/avatar.png",
					},
				},
			},
			siteURL:  "https://example.com",
			expected: "https://example.com/images/avatar.png",
		},
		{
			name: "absolute author image",
			config: &lifecycle.Config{
				Extra: map[string]interface{}{
					"seo": map[string]interface{}{
						"author_image": "https://cdn.example.com/avatar.png",
					},
				},
			},
			siteURL:  "https://example.com",
			expected: "https://cdn.example.com/avatar.png",
		},
		{
			name: "fallback to default_image",
			config: &lifecycle.Config{
				Extra: map[string]interface{}{
					"seo": map[string]interface{}{
						"default_image": "/images/default.png",
					},
				},
			},
			siteURL:  "https://example.com",
			expected: "https://example.com/images/default.png",
		},
		{
			name: "author_image takes precedence over default_image",
			config: &lifecycle.Config{
				Extra: map[string]interface{}{
					"seo": map[string]interface{}{
						"author_image":  "/images/author.png",
						"default_image": "/images/default.png",
					},
				},
			},
			siteURL:  "https://example.com",
			expected: "https://example.com/images/author.png",
		},
		{
			name: "no seo config",
			config: &lifecycle.Config{
				Extra: map[string]interface{}{},
			},
			siteURL:  "https://example.com",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAuthorImageURL(tt.config, tt.siteURL)
			if result != tt.expected {
				t.Errorf("getAuthorImageURL() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestWellKnownPlugin_Write_LinksResource(t *testing.T) {
	outputDir := t.TempDir()
	enabled := true
	wellKnownConfig := models.WellKnownConfig{
		Enabled:      &enabled,
		AutoGenerate: []string{"links"},
	}

	config := &lifecycle.Config{
		OutputDir: outputDir,
		Extra: map[string]interface{}{
			"url":        "https://example.com",
			"well_known": wellKnownConfig,
		},
	}

	posts := []*models.Post{
		{
			Href: "/alpha/",
			Outlinks: []*models.Link{
				{SourceURL: "https://example.com/alpha/", TargetURL: "https://go.dev/doc", IsInternal: false},
				{SourceURL: "https://example.com/alpha/", TargetURL: "https://go.dev/doc", IsInternal: false}, // duplicate
				{SourceURL: "https://example.com/alpha/", TargetURL: "https://pkg.go.dev/", IsInternal: false},
				{SourceURL: "https://example.com/alpha/", TargetURL: "https://example.com/internal/", IsInternal: true},
			},
		},
		{
			Href: "/beta/",
			Outlinks: []*models.Link{
				{SourceURL: "https://example.com/beta/", TargetURL: "https://go.dev/ref/spec", IsInternal: false},
			},
		},
	}

	m := lifecycle.NewManager()
	m.SetConfig(config)
	m.SetPosts(posts)

	plugin := NewWellKnownPlugin()
	if err := plugin.Write(m); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	linksPath := filepath.Join(outputDir, filepath.FromSlash(".well-known/links"))
	content, err := os.ReadFile(linksPath)
	if err != nil {
		t.Fatalf("reading links file: %v", err)
	}

	var domains []linksDomainFixture
	if err := json.Unmarshal(content, &domains); err != nil {
		t.Fatalf("unmarshal links JSON: %v", err)
	}

	if len(domains) != 2 {
		t.Fatalf("expected 2 domains, got %d", len(domains))
	}

	if domains[0].Domain != "go.dev" || domains[0].Count != 2 {
		t.Fatalf("expected go.dev with count 2, got %#v", domains[0])
	}
	if domains[1].Domain != "pkg.go.dev" || domains[1].Count != 1 {
		t.Fatalf("expected pkg.go.dev with count 1, got %#v", domains[1])
	}

	internalPath := filepath.Join(outputDir, filepath.FromSlash(".well-known/internal-links"))
	internalContent, err := os.ReadFile(internalPath)
	if err != nil {
		t.Fatalf("reading internal links file: %v", err)
	}

	var targets []internalTargetFixture
	if err := json.Unmarshal(internalContent, &targets); err != nil {
		t.Fatalf("unmarshal internal links JSON: %v", err)
	}

	if len(targets) != 1 {
		t.Fatalf("expected 1 internal target, got %d", len(targets))
	}
	if targets[0].TargetURL != "/internal" && targets[0].TargetURL != "/internal/" {
		t.Fatalf("expected internal target /internal, got %q", targets[0].TargetURL)
	}
	if targets[0].Count != 1 {
		t.Fatalf("expected internal target count 1, got %d", targets[0].Count)
	}
}
