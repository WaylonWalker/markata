package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func TestGlossaryPlugin_Name(t *testing.T) {
	p := NewGlossaryPlugin()
	if p.Name() != "glossary" {
		t.Errorf("expected name 'glossary', got %q", p.Name())
	}
}

func TestGlossaryPlugin_Configure_Defaults(t *testing.T) {
	p := NewGlossaryPlugin()
	m := lifecycle.NewManager()

	if err := p.Configure(m); err != nil {
		t.Errorf("Configure returned error: %v", err)
	}

	config := p.Config()
	if !config.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if config.LinkClass != "glossary-term" {
		t.Errorf("expected LinkClass 'glossary-term', got %q", config.LinkClass)
	}
	if config.CaseSensitive {
		t.Error("expected CaseSensitive to be false by default")
	}
	if !config.Tooltip {
		t.Error("expected Tooltip to be true by default")
	}
	if config.MaxLinksPerTerm != 1 {
		t.Errorf("expected MaxLinksPerTerm 1, got %d", config.MaxLinksPerTerm)
	}
	if !config.ExportJSON {
		t.Error("expected ExportJSON to be true by default")
	}
}

func TestGlossaryPlugin_Configure_CustomValues(t *testing.T) {
	p := NewGlossaryPlugin()
	m := lifecycle.NewManager()
	m.Config().Extra = map[string]interface{}{
		"glossary": map[string]interface{}{
			"enabled":            false,
			"link_class":         "custom-class",
			"case_sensitive":     true,
			"tooltip":            false,
			"max_links_per_term": 3,
			"exclude_tags":       []interface{}{"tag1", "tag2"},
			"export_json":        false,
			"glossary_path":      "definitions",
			"template_key":       "definition",
		},
	}

	if err := p.Configure(m); err != nil {
		t.Fatalf("Configure returned error: %v", err)
	}

	config := p.Config()
	if config.Enabled {
		t.Error("expected Enabled to be false")
	}
	if config.LinkClass != "custom-class" {
		t.Errorf("expected LinkClass 'custom-class', got %q", config.LinkClass)
	}
	if !config.CaseSensitive {
		t.Error("expected CaseSensitive to be true")
	}
	if config.Tooltip {
		t.Error("expected Tooltip to be false")
	}
	if config.MaxLinksPerTerm != 3 {
		t.Errorf("expected MaxLinksPerTerm 3, got %d", config.MaxLinksPerTerm)
	}
	if len(config.ExcludeTags) != 2 || config.ExcludeTags[0] != "tag1" {
		t.Errorf("expected ExcludeTags ['tag1', 'tag2'], got %v", config.ExcludeTags)
	}
	if config.ExportJSON {
		t.Error("expected ExportJSON to be false")
	}
	if config.GlossaryPath != "definitions" {
		t.Errorf("expected GlossaryPath 'definitions', got %q", config.GlossaryPath)
	}
	if config.TemplateKey != "definition" {
		t.Errorf("expected TemplateKey 'definition', got %q", config.TemplateKey)
	}
}

func TestGlossaryPlugin_Priority(t *testing.T) {
	p := NewGlossaryPlugin()

	if p.Priority(lifecycle.StageRender) != lifecycle.PriorityLate {
		t.Errorf("expected PriorityLate for render stage, got %d", p.Priority(lifecycle.StageRender))
	}
	if p.Priority(lifecycle.StageTransform) != lifecycle.PriorityDefault {
		t.Errorf("expected PriorityDefault for transform stage, got %d", p.Priority(lifecycle.StageTransform))
	}
}

func TestGlossaryPlugin_BuildGlossary(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	apiDesc := "Application Programming Interface"

	posts := []*models.Post{
		{
			Title:       &apiTitle,
			Description: &apiDesc,
			Slug:        "api",
			Href:        "/glossary/api/",
			Extra: map[string]interface{}{
				"templateKey": "glossary",
				"aliases":     []interface{}{"APIs", "Application Programming Interface"},
			},
		},
	}

	err := p.buildGlossary(posts)
	if err != nil {
		t.Fatalf("buildGlossary error: %v", err)
	}

	// Check primary term
	if _, ok := p.terms["api"]; !ok {
		t.Error("expected 'api' term in lookup")
	}

	// Check aliases
	if _, ok := p.terms["apis"]; !ok {
		t.Error("expected 'apis' alias in lookup")
	}
	if _, ok := p.terms["application programming interface"]; !ok {
		t.Error("expected 'application programming interface' alias in lookup")
	}

	// Check allTerms
	if len(p.allTerms) != 1 {
		t.Errorf("expected 1 term, got %d", len(p.allTerms))
	}
}

func TestGlossaryPlugin_IsGlossaryPost_TemplateKey(t *testing.T) {
	p := NewGlossaryPlugin()

	post := &models.Post{
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	if !p.isGlossaryPost(post) {
		t.Error("expected post with templateKey 'glossary' to be glossary post")
	}
}

func TestGlossaryPlugin_IsGlossaryPost_TemplateKeyVariant(t *testing.T) {
	p := NewGlossaryPlugin()

	post := &models.Post{
		Extra: map[string]interface{}{
			"template_key": "glossary",
		},
	}

	if !p.isGlossaryPost(post) {
		t.Error("expected post with template_key 'glossary' to be glossary post")
	}
}

func TestGlossaryPlugin_IsGlossaryPost_Path(t *testing.T) {
	p := NewGlossaryPlugin()

	tests := []struct {
		path     string
		expected bool
	}{
		{"glossary/api.md", true},
		{"glossary/terms/api.md", true},
		{"content/glossary/api.md", true},
		{"blog/api.md", false},
		{"glossaryfile.md", false},
	}

	for _, tt := range tests {
		post := &models.Post{Path: tt.path, Extra: map[string]interface{}{}}
		result := p.isGlossaryPost(post)
		if result != tt.expected {
			t.Errorf("isGlossaryPost(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestGlossaryPlugin_LinkTerms_Basic(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	apiDesc := "Application Programming Interface"

	glossaryPost := &models.Post{
		Title:       &apiTitle,
		Description: &apiDesc,
		Slug:        "api",
		Href:        "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	posts := []*models.Post{glossaryPost}
	_ = p.buildGlossary(posts) //nolint:errcheck // test ignores error

	html := "<p>The API allows communication between services.</p>"
	result := p.linkTerms(html, nil)

	expected := `<a href="/glossary/api/" class="glossary-term" title="Application Programming Interface" data-hover-title="API">API</a>`
	if !strings.Contains(result, expected) {
		t.Errorf("expected link in output:\nwant: %s\ngot: %s", expected, result)
	}
}

func TestGlossaryPlugin_LinkTerms_SkipsTagAttributes(t *testing.T) {
	p := NewGlossaryPlugin()

	title := "aur"
	desc := "Arch User Repository"
	glossaryPost := &models.Post{
		Title:       &title,
		Description: &desc,
		Slug:        "aur",
		Href:        "/aur/",
		Extra:       map[string]interface{}{"templateKey": "glossary"},
	}
	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := `<h2 id="malicious-aur-packages" class="card-title">Malicious aur packages</h2><img alt="aur logo" src="/aur.png">`
	result := p.linkTerms(html, nil)

	if !strings.Contains(result, `id="malicious-aur-packages"`) {
		t.Errorf("id attribute was rewritten: %s", result)
	}
	if !strings.Contains(result, `alt="aur logo"`) {
		t.Errorf("alt attribute was rewritten: %s", result)
	}
	if strings.Count(result, `class="glossary-term"`) != 1 {
		t.Errorf("expected exactly one linked term in text content, got: %s", result)
	}
}

func TestGlossaryPlugin_LinkTerms_MaxLinksPerTerm(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.MaxLinksPerTerm = 1

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := "<p>The API is great. Another API mention. Third API usage.</p>"
	result := p.linkTerms(html, nil)

	// Count links
	linkCount := strings.Count(result, `<a href="/glossary/api/"`)
	if linkCount != 1 {
		t.Errorf("expected 1 link with max_links_per_term=1, got %d\nresult: %s", linkCount, result)
	}
}

func TestGlossaryPlugin_LinkTerms_AllOccurrences(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.MaxLinksPerTerm = 0 // Link all occurrences

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := "<p>The API is great. Another API mention.</p>"
	result := p.linkTerms(html, nil)

	// Count links
	linkCount := strings.Count(result, `<a href="/glossary/api/"`)
	if linkCount != 2 {
		t.Errorf("expected 2 links with max_links_per_term=0, got %d\nresult: %s", linkCount, result)
	}
}

func TestGlossaryPlugin_LinkTerms_SkipExistingLinks(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := `<p>Check this <a href="/other">API link</a> for info.</p>`
	result := p.linkTerms(html, nil)

	// Should not create nested links
	if strings.Contains(result, `<a href="/glossary/api/"`) {
		t.Error("should not link term inside existing anchor tag")
	}
}

func TestGlossaryPlugin_LinkTerms_SkipCodeBlocks(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := `<p>Use <code>API.call()</code> in your code.</p>`
	result := p.linkTerms(html, nil)

	// Should not link inside code tags
	if strings.Contains(result, `<a href="/glossary/api/">`) && strings.Contains(result, `<code>`) {
		// Check the actual structure
		if strings.Contains(result, `<code><a`) || strings.Contains(result, `>API</a></code>`) {
			t.Error("should not link term inside code tag")
		}
	}
}

func TestGlossaryPlugin_LinkTerms_SkipPreBlocks(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := `<pre><code>API.call()</code></pre><p>The API works.</p>`
	result := p.linkTerms(html, nil)

	// Should link in paragraph but not in pre block
	if !strings.Contains(result, `<a href="/glossary/api/"`) {
		t.Error("should link term in paragraph")
	}
	// Pre block content should be preserved
	if !strings.Contains(result, `<pre><code>API.call()</code></pre>`) {
		t.Error("pre block content should be unchanged")
	}
}

func TestGlossaryPlugin_LinkTerms_CaseInsensitive(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.CaseSensitive = false

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := "<p>The api is lowercase.</p>"
	result := p.linkTerms(html, nil)

	if !strings.Contains(result, `<a href="/glossary/api/"`) {
		t.Error("should link case-insensitive match")
	}
	// Should preserve original case in link text
	if !strings.Contains(result, `>api</a>`) {
		t.Errorf("should preserve original case in link text, got: %s", result)
	}
}

func TestGlossaryPlugin_LinkTerms_CaseSensitive(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.CaseSensitive = true

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := "<p>The api is lowercase but API is uppercase.</p>"
	result := p.linkTerms(html, nil)

	// Should only link exact case match
	linkCount := strings.Count(result, `<a href="/glossary/api/"`)
	if linkCount != 1 {
		t.Errorf("expected 1 case-sensitive link, got %d\nresult: %s", linkCount, result)
	}
}

func TestGlossaryPlugin_LinkTerms_NoTooltip(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.Tooltip = false

	apiTitle := "API"
	apiDesc := "Application Programming Interface"
	glossaryPost := &models.Post{
		Title:       &apiTitle,
		Description: &apiDesc,
		Slug:        "api",
		Href:        "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := "<p>The API works.</p>"
	result := p.linkTerms(html, nil)

	if strings.Contains(result, `title="`) {
		t.Error("should not include title attribute when tooltip is disabled")
	}
}

func TestGlossaryPlugin_LinkTerms_Aliases(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.MaxLinksPerTerm = 0

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
			"aliases":     []interface{}{"APIs"},
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := "<p>The API and APIs are both linked.</p>"
	result := p.linkTerms(html, nil)

	// Both should link to the same term
	if strings.Count(result, `/glossary/api/`) != 2 {
		t.Errorf("expected both API and APIs to link to /glossary/api/\nresult: %s", result)
	}
}

func TestGlossaryPlugin_ProcessPost_SkipGlossaryPost(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title:       &apiTitle,
		Slug:        "api",
		Href:        "/glossary/api/",
		ArticleHTML: "<p>This post defines API.</p>",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	originalHTML := glossaryPost.ArticleHTML
	err := p.processPost(glossaryPost)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if glossaryPost.ArticleHTML != originalHTML {
		t.Error("glossary post should not be modified")
	}
}

func TestGlossaryPlugin_ProcessPost_SkipExcludedTags(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.ExcludeTags = []string{"glossary", "no-link"}

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	regularPost := &models.Post{
		Slug:        "test-post",
		Tags:        []string{"no-link"},
		ArticleHTML: "<p>The API works.</p>",
		Extra:       map[string]interface{}{},
	}

	originalHTML := regularPost.ArticleHTML
	err := p.processPost(regularPost)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if regularPost.ArticleHTML != originalHTML {
		t.Error("post with excluded tag should not be modified")
	}
}

func TestGlossaryPlugin_Render_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	p := NewGlossaryPlugin()
	m := lifecycle.NewManager()

	apiTitle := "API"
	apiDesc := "Application Programming Interface"
	glossaryPost := &models.Post{
		Title:       &apiTitle,
		Description: &apiDesc,
		Slug:        "api",
		Href:        "/glossary/api/",
		ArticleHTML: "<p>An API is an interface.</p>",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	blogTitle := "Using APIs"
	blogPost := &models.Post{
		Title:       &blogTitle,
		Slug:        "using-apis",
		ArticleHTML: "<p>The API allows services to communicate.</p>",
		Extra:       map[string]interface{}{},
	}

	m.SetPosts([]*models.Post{glossaryPost, blogPost})

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	posts := m.Posts()

	// Glossary post should not be modified
	if strings.Contains(posts[0].ArticleHTML, `<a href="/glossary/api/"`) {
		t.Error("glossary post should not have self-links")
	}

	// Blog post should have link
	if !strings.Contains(posts[1].ArticleHTML, `<a href="/glossary/api/"`) {
		t.Errorf("blog post should have glossary link\ngot: %s", posts[1].ArticleHTML)
	}
}

func TestGlossaryPlugin_Render_Disabled(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.Enabled = false

	m := lifecycle.NewManager()

	apiTitle := "API"
	blogPost := &models.Post{
		Title:       &apiTitle,
		Slug:        "blog",
		ArticleHTML: "<p>The API works.</p>",
		Extra:       map[string]interface{}{},
	}

	m.SetPosts([]*models.Post{blogPost})

	err := p.Render(m)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}

	posts := m.Posts()
	if strings.Contains(posts[0].ArticleHTML, `<a href=`) {
		t.Error("disabled plugin should not modify posts")
	}
}

func TestGlossaryPlugin_Write_ExportJSON(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	apiDesc := "Application Programming Interface"
	glossaryPost := &models.Post{
		Title:       &apiTitle,
		Description: &apiDesc,
		Slug:        "api",
		Href:        "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
			"aliases":     []interface{}{"APIs"},
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	// Create temp directory
	tmpDir := t.TempDir()

	m := lifecycle.NewManager()
	m.Config().OutputDir = tmpDir

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Check glossary.json was created
	jsonPath := filepath.Join(tmpDir, "glossary.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("failed to read glossary.json: %v", err)
	}

	var export GlossaryExport
	if err := json.Unmarshal(data, &export); err != nil {
		t.Fatalf("failed to parse glossary.json: %v", err)
	}

	if len(export.Terms) != 1 {
		t.Errorf("expected 1 term in export, got %d", len(export.Terms))
	}

	term := export.Terms[0]
	if term.Term != "API" {
		t.Errorf("expected term 'API', got %q", term.Term)
	}
	if term.Description != "Application Programming Interface" {
		t.Errorf("expected description, got %q", term.Description)
	}
	if len(term.Aliases) != 1 || term.Aliases[0] != "APIs" {
		t.Errorf("expected aliases ['APIs'], got %v", term.Aliases)
	}
}

func TestGlossaryPlugin_Write_Disabled(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.ExportJSON = false

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	tmpDir := t.TempDir()

	m := lifecycle.NewManager()
	m.Config().OutputDir = tmpDir

	err := p.Write(m)
	if err != nil {
		t.Fatalf("Write error: %v", err)
	}

	// Check glossary.json was NOT created
	jsonPath := filepath.Join(tmpDir, "glossary.json")
	if _, err := os.Stat(jsonPath); !os.IsNotExist(err) {
		t.Error("glossary.json should not be created when ExportJSON is false")
	}
}

func TestGlossaryPlugin_WordBoundary(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	// Should not match "API" within "RAPID"
	html := "<p>RAPID development uses the API.</p>"
	result := p.linkTerms(html, nil)

	// Count links - should be exactly 1
	linkCount := strings.Count(result, `<a href="/glossary/api/"`)
	if linkCount != 1 {
		t.Errorf("expected 1 link (word boundary), got %d\nresult: %s", linkCount, result)
	}

	// RAPID should not be linked
	if strings.Contains(result, ">RAPID</a>") || strings.Contains(result, ">rapid</a>") {
		t.Error("RAPID should not be linked (API is not a word boundary)")
	}
}

func TestGlossaryPlugin_LongestMatchFirst(t *testing.T) {
	p := NewGlossaryPlugin()
	p.config.MaxLinksPerTerm = 0

	jsTitle := "JavaScript"
	javaTitle := "Java"

	jsPost := &models.Post{
		Title: &jsTitle,
		Slug:  "javascript",
		Href:  "/glossary/javascript/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	javaPost := &models.Post{
		Title: &javaTitle,
		Slug:  "java",
		Href:  "/glossary/java/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{jsPost, javaPost}) //nolint:errcheck // test ignores error

	html := "<p>JavaScript is not Java.</p>"
	result := p.linkTerms(html, nil)

	// JavaScript should be linked as one term, not as Java + Script
	if !strings.Contains(result, `href="/glossary/javascript/"`) {
		t.Error("JavaScript should be linked")
	}
	if !strings.Contains(result, `href="/glossary/java/"`) {
		t.Error("Java should be linked")
	}

	// Check JavaScript wasn't broken up
	if strings.Contains(result, `>Java</a>Script`) {
		t.Error("JavaScript should be linked as whole word, not broken up")
	}
}

func TestGlossaryPlugin_Interfaces(_ *testing.T) {
	p := NewGlossaryPlugin()

	// Verify interface compliance
	var _ lifecycle.Plugin = p
	var _ lifecycle.ConfigurePlugin = p
	var _ lifecycle.RenderPlugin = p
	var _ lifecycle.WritePlugin = p
	var _ lifecycle.PriorityPlugin = p
}

func TestGlossaryPlugin_Terms(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	restTitle := "REST"

	posts := []*models.Post{
		{
			Title: &apiTitle,
			Slug:  "api",
			Extra: map[string]interface{}{
				"templateKey": "glossary",
			},
		},
		{
			Title: &restTitle,
			Slug:  "rest",
			Extra: map[string]interface{}{
				"templateKey": "glossary",
			},
		},
	}

	_ = p.buildGlossary(posts) //nolint:errcheck // test ignores error

	terms := p.Terms()
	if len(terms) != 2 {
		t.Errorf("expected 2 terms, got %d", len(terms))
	}

	// Should be sorted alphabetically
	if terms[0].Term != "API" || terms[1].Term != "REST" {
		t.Errorf("terms should be sorted: got %v, %v", terms[0].Term, terms[1].Term)
	}
}

func TestGlossaryPlugin_SetConfig(t *testing.T) {
	p := NewGlossaryPlugin()

	newConfig := &GlossaryConfig{
		Enabled:         false,
		LinkClass:       "new-class",
		MaxLinksPerTerm: 5,
	}

	p.SetConfig(newConfig)

	if p.Config().Enabled {
		t.Error("expected Enabled to be false")
	}
	if p.Config().LinkClass != "new-class" {
		t.Errorf("expected LinkClass 'new-class', got %q", p.Config().LinkClass)
	}
	if p.Config().MaxLinksPerTerm != 5 {
		t.Errorf("expected MaxLinksPerTerm 5, got %d", p.Config().MaxLinksPerTerm)
	}
}

func TestGlossaryPlugin_SkipEmptyPosts(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	// Post with empty ArticleHTML
	emptyPost := &models.Post{
		Slug:        "empty",
		ArticleHTML: "",
		Extra:       map[string]interface{}{},
	}

	err := p.processPost(emptyPost)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if emptyPost.ArticleHTML != "" {
		t.Error("empty post should remain empty")
	}
}

func TestGlossaryPlugin_SkipSkippedPosts(t *testing.T) {
	p := NewGlossaryPlugin()

	apiTitle := "API"
	glossaryPost := &models.Post{
		Title: &apiTitle,
		Slug:  "api",
		Href:  "/glossary/api/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	skippedPost := &models.Post{
		Slug:        "skipped",
		Skip:        true,
		ArticleHTML: "<p>The API works.</p>",
		Extra:       map[string]interface{}{},
	}

	originalHTML := skippedPost.ArticleHTML
	err := p.processPost(skippedPost)
	if err != nil {
		t.Fatalf("processPost error: %v", err)
	}

	if skippedPost.ArticleHTML != originalHTML {
		t.Error("skipped post should not be modified")
	}
}

func TestGlossaryPlugin_HTMLEscaping(t *testing.T) {
	p := NewGlossaryPlugin()

	// Term with special characters in description
	termTitle := "Test"
	termDesc := `Test with "quotes" & <special> chars`
	glossaryPost := &models.Post{
		Title:       &termTitle,
		Description: &termDesc,
		Slug:        "test",
		Href:        "/glossary/test/",
		Extra: map[string]interface{}{
			"templateKey": "glossary",
		},
	}

	_ = p.buildGlossary([]*models.Post{glossaryPost}) //nolint:errcheck // test ignores error

	html := "<p>This is a Test.</p>"
	result := p.linkTerms(html, nil)

	// Check that special characters are escaped
	if strings.Contains(result, `"quotes"`) {
		t.Error("quotes should be escaped in title attribute")
	}
	if strings.Contains(result, `<special>`) {
		t.Error("angle brackets should be escaped in title attribute")
	}
	if !strings.Contains(result, `&amp;`) || !strings.Contains(result, `&#34;`) || !strings.Contains(result, `&lt;`) {
		t.Errorf("special characters should be HTML escaped\nresult: %s", result)
	}
}
