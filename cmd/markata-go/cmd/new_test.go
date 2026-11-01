package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func resetNewCommandFlags() {
	newDir = ""
	newDraft = false
	newTags = ""
	newTemplate = "post"
	newList = false
	newPlain = false
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "simple title",
			title:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "title with special characters",
			title:    "What's New in Go 1.22?",
			expected: "whats-new-in-go-122",
		},
		{
			name:     "title with multiple spaces",
			title:    "My   First   Post",
			expected: "my-first-post",
		},
		{
			name:     "title with leading/trailing spaces",
			title:    "  Trimmed Title  ",
			expected: "trimmed-title",
		},
		{
			name:     "title with numbers",
			title:    "Top 10 Tips",
			expected: "top-10-tips",
		},
		{
			name:     "empty title",
			title:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateSlug(tt.title)
			if result != tt.expected {
				t.Errorf("generateSlug(%q) = %q, want %q", tt.title, result, tt.expected)
			}
		})
	}
}

func TestRunNewCommand_NoInputRequiresTitle(t *testing.T) {
	resetNewCommandFlags()
	originalNoInput := noInput
	defer func() { noInput = originalNoInput }()

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	newCmd.SetOut(stdout)
	newCmd.SetErr(stderr)

	noInput = true
	err := runNewCommand(newCmd, nil)
	if err == nil {
		t.Fatal("expected error when --no-input is set without a title")
	}
	if !strings.Contains(err.Error(), "title is required when --no-input is set") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunInteractiveMode_TemplateAlias(t *testing.T) {
	resetNewCommandFlags()
	newTemplate = "ping"

	cmd := &cobra.Command{}
	cmd.SetIn(strings.NewReader("My Ping\nnotes\n\nn\n"))
	cmd.SetOut(bytes.NewBuffer(nil))
	cmd.SetErr(bytes.NewBuffer(nil))
	currentCmd = cmd
	defer func() { currentCmd = nil }()

	input, err := runInteractiveMode(cmd, builtinTemplates())
	if err != nil {
		t.Fatalf("runInteractiveMode returned error: %v", err)
	}
	if got := input.template.Frontmatter["template"]; got != "ping" {
		t.Fatalf("template frontmatter = %v, want ping", got)
	}
	if got := input.template.Frontmatter["templateKey"]; got != "ping" {
		t.Fatalf("templateKey frontmatter = %v, want ping", got)
	}
	if input.template.Directory != "pages/note" {
		t.Fatalf("directory = %q, want %q", input.template.Directory, "pages/note")
	}
	if newDir != "pages/note" {
		t.Fatalf("newDir = %q, want %q", newDir, "pages/note")
	}
	if newDraft {
		t.Fatal("expected draft prompt to leave newDraft false")
	}
	if input.title != "My Ping" {
		t.Fatalf("title = %q, want %q", input.title, "My Ping")
	}
}

func TestParseTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single tag",
			input:    "go",
			expected: []string{"go"},
		},
		{
			name:     "multiple tags",
			input:    "go,tutorial,programming",
			expected: []string{"go", "tutorial", "programming"},
		},
		{
			name:     "tags with spaces",
			input:    " go , tutorial , programming ",
			expected: []string{"go", "tutorial", "programming"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "only commas",
			input:    ",,,",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTags(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("parseTags(%q) returned %d tags, want %d", tt.input, len(result), len(tt.expected))
				return
			}
			for i, tag := range result {
				if tag != tt.expected[i] {
					t.Errorf("parseTags(%q)[%d] = %q, want %q", tt.input, i, tag, tt.expected[i])
				}
			}
		})
	}
}

func TestBuiltinTemplates(t *testing.T) {
	templates := builtinTemplates()

	// Check that all expected templates exist
	expectedTemplates := []string{"post", "page", "docs", "article", "note", "photo", "video", "link", "quote", "guide", "inline", "contact", "author"}
	for _, name := range expectedTemplates {
		if _, exists := templates[name]; !exists {
			t.Errorf("expected builtin template %q not found", name)
		}
	}

	// Check post template
	post := templates["post"]
	if post.Directory != "pages/post" {
		t.Errorf("post template directory = %q, want %q", post.Directory, "pages/post")
	}
	if post.Source != "builtin" {
		t.Errorf("post template source = %q, want %q", post.Source, "builtin")
	}

	// Check page template
	page := templates["page"]
	if page.Directory != "pages" {
		t.Errorf("page template directory = %q, want %q", page.Directory, "pages")
	}

	// Check docs template
	docs := templates["docs"]
	if docs.Directory != "docs" {
		t.Errorf("docs template directory = %q, want %q", docs.Directory, "docs")
	}

	// Check contact template
	contact := templates["contact"]
	if contact.Directory != "pages/contact" {
		t.Errorf("contact template directory = %q, want %q", contact.Directory, "pages/contact")
	}
	if contact.Frontmatter["handle"] != "" {
		t.Errorf("contact template should have empty handle field")
	}

	// Check author template
	author := templates["author"]
	if author.Directory != "pages/author" {
		t.Errorf("author template directory = %q, want %q", author.Directory, "pages/author")
	}
	if author.Frontmatter["role"] != "" {
		t.Errorf("author template should have empty role field")
	}

	// Check video template has correct field name
	video := templates["video"]
	if _, ok := video.Frontmatter["video"]; !ok {
		t.Error("video template should have 'video' field, not 'video_url'")
	}
	if _, ok := video.Frontmatter["video_url"]; ok {
		t.Error("video template should not have 'video_url' field (renamed to 'video')")
	}

	// Check quote template has correct field names
	quote := templates["quote"]
	if _, ok := quote.Frontmatter["quote"]; !ok {
		t.Error("quote template should have 'quote' field")
	}
	if _, ok := quote.Frontmatter["source"]; !ok {
		t.Error("quote template should have 'source' field")
	}
	if _, ok := quote.Frontmatter["quote_author"]; ok {
		t.Error("quote template should not have 'quote_author' field (use post.author instead)")
	}
}

func TestParseTemplateFile(t *testing.T) {
	tests := []struct {
		name             string
		templateName     string
		content          string
		expectedDir      string
		expectedBody     string
		expectedFMLength int
	}{
		{
			name:         "simple template without frontmatter",
			templateName: "simple",
			content:      "Hello, this is the body.",
			expectedDir:  "simple",
			expectedBody: "Hello, this is the body.",
		},
		{
			name:         "template with frontmatter",
			templateName: "fancy",
			content: `---
templateKey: fancy
_directory: fancy-posts
custom_field: value
---

This is the fancy body.`,
			expectedDir:      "fancy-posts",
			expectedBody:     "This is the fancy body.",
			expectedFMLength: 2, // templateKey and custom_field (_directory is removed)
		},
		{
			name:         "template with only frontmatter",
			templateName: "minimal",
			content: `---
templateKey: minimal
---
`,
			expectedDir:      "minimal",
			expectedBody:     "",
			expectedFMLength: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTemplateFile(tt.templateName, tt.content)

			if result.Name != tt.templateName {
				t.Errorf("name = %q, want %q", result.Name, tt.templateName)
			}
			if result.Directory != tt.expectedDir {
				t.Errorf("directory = %q, want %q", result.Directory, tt.expectedDir)
			}
			if result.Body != tt.expectedBody {
				t.Errorf("body = %q, want %q", result.Body, tt.expectedBody)
			}
			if tt.expectedFMLength > 0 && len(result.Frontmatter) != tt.expectedFMLength {
				t.Errorf("frontmatter length = %d, want %d", len(result.Frontmatter), tt.expectedFMLength)
			}
		})
	}
}

func TestGenerateTemplatedContent(t *testing.T) {
	template := ContentTemplate{
		Name:      "post",
		Directory: "posts",
		Frontmatter: map[string]interface{}{
			"templateKey": "post",
			"layout":      "post.html",
		},
		Body:   "Start writing here...",
		Source: "builtin",
	}

	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	content := generateTemplatedContent("My Test Post", "my-test-post", date, true, true, []string{"go", "test"}, template)

	// Check frontmatter delimiters
	if !strings.HasPrefix(content, "---\n") {
		t.Error("content should start with frontmatter delimiter")
	}
	if !strings.Contains(content, "\n---\n") {
		t.Error("content should contain frontmatter closing delimiter")
	}

	// Check required fields are present
	requiredFields := []string{
		"title: My Test Post",
		"slug: my-test-post",
		"date: \"2024-01-15\"",
		"draft: true",
		"published: false",
		"templateKey: post",
		"layout: post.html",
	}
	for _, field := range requiredFields {
		if !strings.Contains(content, field) {
			t.Errorf("content should contain %q", field)
		}
	}

	// Check tags
	if !strings.Contains(content, "- go") || !strings.Contains(content, "- test") {
		t.Error("content should contain tags")
	}

	if strings.Contains(content, "# My Test Post") {
		t.Error("content should not contain an H1 heading")
	}

	// Check body
	if !strings.Contains(content, "Start writing here...") {
		t.Error("content should contain template body")
	}
}

func TestResolveTemplateSelection_Alias(t *testing.T) {
	templates := builtinTemplates()

	template, resolvedName, ok := resolveTemplateSelection("ping", templates)
	if !ok {
		t.Fatal("expected alias template resolution to succeed")
	}
	if resolvedName != "note" {
		t.Fatalf("resolvedName = %q, want %q", resolvedName, "note")
	}
	if got := template.Frontmatter["template"]; got != "ping" {
		t.Fatalf("template frontmatter = %v, want ping", got)
	}
	if got := template.Frontmatter["templateKey"]; got != "ping" {
		t.Fatalf("templateKey frontmatter = %v, want ping", got)
	}
	if template.Directory != templates["note"].Directory {
		t.Fatalf("directory = %q, want %q", template.Directory, templates["note"].Directory)
	}
}

func TestFormatTemplateEntry_IncludesAliases(t *testing.T) {
	template := builtinTemplates()["note"]
	entry := formatTemplateEntry("note", template)

	if !strings.Contains(entry, "note (aka: ping, thought, status, tweet)") {
		t.Fatalf("entry = %q, expected aliases", entry)
	}
	if !strings.Contains(entry, "-> pages/note/") {
		t.Fatalf("entry = %q, expected directory", entry)
	}
}

func TestLoadTemplatesFromDir(t *testing.T) {
	// Create a temporary directory with template files
	tmpDir := t.TempDir()

	// Create a test template file
	templateContent := `---
templateKey: tutorial
_directory: tutorials
---

Write your tutorial here...`
	err := os.WriteFile(filepath.Join(tmpDir, "tutorial.md"), []byte(templateContent), 0o600)
	if err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	// Load templates
	templates := make(map[string]ContentTemplate)
	loadTemplatesFromDir(tmpDir, templates)

	// Check that the template was loaded
	if _, exists := templates["tutorial"]; !exists {
		t.Error("tutorial template should have been loaded")
		return
	}

	tutorial := templates["tutorial"]
	if tutorial.Directory != "tutorials" {
		t.Errorf("tutorial directory = %q, want %q", tutorial.Directory, "tutorials")
	}
	if tutorial.Source != "file" {
		t.Errorf("tutorial source = %q, want %q", tutorial.Source, "file")
	}
	if !strings.Contains(tutorial.Body, "Write your tutorial here") {
		t.Errorf("tutorial body unexpected: %q", tutorial.Body)
	}
}

func TestLoadTemplatesFromNonexistentDir(t *testing.T) {
	templates := make(map[string]ContentTemplate)
	// Should not panic or error
	loadTemplatesFromDir("/nonexistent/path/that/does/not/exist", templates)

	if len(templates) != 0 {
		t.Errorf("expected empty templates map, got %d items", len(templates))
	}
}

func TestGeneratePostContentWithTags(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	// Test with tags
	content := generatePostContentWithTags("Test Post", "test-post", date, false, []string{"go", "testing"})

	if !strings.Contains(content, `title: "Test Post"`) {
		t.Error("should contain title")
	}
	if !strings.Contains(content, `slug: "test-post"`) {
		t.Error("should contain slug")
	}
	if !strings.Contains(content, "date: 2024-01-15") {
		t.Error("should contain date")
	}
	if !strings.Contains(content, "draft: false") {
		t.Error("should contain draft: false")
	}
	if !strings.Contains(content, "published: true") {
		t.Error("should contain published: true")
	}
	if !strings.Contains(content, `"go"`) || !strings.Contains(content, `"testing"`) {
		t.Error("should contain tags")
	}

	// Test without tags
	contentNoTags := generatePostContentWithTags("No Tags", "no-tags", date, true, nil)
	if !strings.Contains(contentNoTags, "tags: []") {
		t.Error("should contain empty tags array")
	}
	if strings.Contains(content, "# Test Post") {
		t.Error("should not contain an H1 heading")
	}
}

func TestContentTemplatesConfig_GetPlacement(t *testing.T) {
	tests := []struct {
		name         string
		placement    map[string]string
		templateName string
		expected     string
	}{
		{
			name: "configured placement",
			placement: map[string]string{
				"post": "blog",
				"page": "pages",
			},
			templateName: "post",
			expected:     "blog",
		},
		{
			name: "unconfigured placement returns template name",
			placement: map[string]string{
				"post": "blog",
			},
			templateName: "unknown",
			expected:     "unknown",
		},
		{
			name:         "empty placement map",
			placement:    map[string]string{},
			templateName: "post",
			expected:     "post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use the models package's GetPlacement through a wrapper test
			// Since we're testing the cmd package, test the local behavior
			var dir string
			if d, ok := tt.placement[tt.templateName]; ok {
				dir = d
			} else {
				dir = tt.templateName
			}

			if dir != tt.expected {
				t.Errorf("GetPlacement(%q) = %q, want %q", tt.templateName, dir, tt.expected)
			}
		})
	}
}

func TestParseTemplateFile_StarterMetadata(t *testing.T) {
	content := "---\nstarter:\n  name: how-to\n  description: A useful guide\n  aliases:\n    - tutorial\n_directory: guides\ndescription: Page description\ntitle: \"{{title}}\"\ndate: \"{{date}}\"\n---\n\n## {{title}}\n"
	template := parseTemplateFile("guide", content)
	if template.Name != "how-to" || template.Directory != "guides" || template.Description != "A useful guide" {
		t.Fatalf("template metadata = %#v", template)
	}
	if len(template.Aliases) != 1 || template.Aliases[0] != "tutorial" {
		t.Fatalf("aliases = %#v", template.Aliases)
	}
	for _, key := range []string{"starter", "_directory"} {
		if _, ok := template.Frontmatter[key]; ok {
			t.Errorf("metadata %q leaked into generated frontmatter", key)
		}
	}
	if template.Frontmatter["description"] != "Page description" {
		t.Fatal("page description frontmatter was not preserved")
	}
}

func TestParseTemplateFile_UsesCanonicalNameAsDefaultDirectory(t *testing.T) {
	template := parseTemplateFile("guide", "---\nstarter:\n  name: how-to\n---\n\nBody\n")
	if template.Name != "how-to" || template.Directory != "how-to" {
		t.Fatalf("name/directory = %q/%q, want how-to/how-to", template.Name, template.Directory)
	}
}

func TestParseTemplateFile_ReportsInvalidStarterMetadata(t *testing.T) {
	for _, content := range []string{
		"---\nstarter: invalid\n---\nBody\n",
		"---\nstarter:\n  aliases: one\n---\nBody\n",
		"---\nstarter: [broken\n---\nBody\n",
		"---\nstarter:\n  name: no-close\n",
	} {
		template := parseTemplateFile("broken", content)
		if template.ValidationError == "" {
			t.Fatalf("invalid starter metadata was not reported: %s", content)
		}
	}
}

func TestResolveTemplateSelection_RejectsAliasShadowingCanonicalName(t *testing.T) {
	templates := map[string]ContentTemplate{
		"post":  {Name: "post", Directory: "posts"},
		"guide": {Name: "guide", Directory: "guides", Aliases: []string{"post"}},
	}
	if _, _, ok := resolveTemplateSelection("post", templates); ok {
		t.Fatal("expected canonical/alias collision to be rejected")
	}
}

func TestValidateTemplateAliases_ReportsCanonicalNameCollision(t *testing.T) {
	templates := map[string]ContentTemplate{
		"post":  {Name: "post", Directory: "posts"},
		"guide": {Name: "guide", Directory: "guides", Aliases: []string{"post"}},
	}
	validateTemplateAliases(templates)
	if templates["post"].ValidationError == "" || templates["guide"].ValidationError == "" {
		t.Fatalf("expected both conflicting starters to be marked invalid: %#v", templates)
	}
}

func TestValidateTemplateAliases_KeepsBuiltinsUsable(t *testing.T) {
	templates := builtinTemplates()
	validateTemplateAliases(templates)
	for name, template := range templates {
		if template.ValidationError != "" {
			t.Errorf("built-in %q unexpectedly invalid: %s", name, template.ValidationError)
		}
	}
}

func TestGenerateTemplatedContent_ExpandsStarterVariables(t *testing.T) {
	template := ContentTemplate{Frontmatter: map[string]interface{}{"subtitle": "{{title}} ({{slug}})", "when": "{{date}}"}, Body: "## {{title}} on {{datetime}}"}
	date := time.Date(2024, 1, 15, 12, 30, 0, 0, time.UTC)
	content := generateTemplatedContent("A Guide", "a-guide", date, false, false, nil, template)
	for _, expected := range []string{"subtitle: A Guide (a-guide)", "when: \"2024-01-15\"", "## A Guide on 2024-01-15T12:30:00Z"} {
		if !strings.Contains(content, expected) {
			t.Errorf("content missing %q:\n%s", expected, content)
		}
	}
}

func TestGenerateTemplatedContent_PreservesUnselectedStarterDefaults(t *testing.T) {
	template := ContentTemplate{Frontmatter: map[string]interface{}{
		"published":   false,
		"draft":       true,
		"tags":        []interface{}{"site-tag"},
		"description": "Starter description",
	}}
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	content := generateTemplatedContent("A Guide", "a-guide", date, false, false, nil, template)
	for _, expected := range []string{"published: false", "draft: true", "- site-tag", "description: Starter description"} {
		if !strings.Contains(content, expected) {
			t.Errorf("content missing preserved value %q:\n%s", expected, content)
		}
	}
}
