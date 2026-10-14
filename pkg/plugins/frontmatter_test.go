package plugins

import (
	"reflect"
	"testing"
)

// =============================================================================
// Frontmatter Parsing Tests based on tests.yaml
// =============================================================================

func TestParseFrontmatter_BasicYAML(t *testing.T) {
	// Test case: "basic yaml frontmatter"
	content := `---
title: Hello World
date: 2024-01-15
---
Content here`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify title
	if metadata["title"] != "Hello World" {
		t.Errorf("title: got %q, want 'Hello World'", metadata["title"])
	}

	// Verify date (can be string or time.Time depending on YAML parser)
	if metadata["date"] == nil {
		t.Error("date should be set")
	}

	// Verify body
	if body != "Content here" {
		t.Errorf("body: got %q, want 'Content here'", body)
	}
}

func TestParseFrontmatter_ScalarTags(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []string
	}{
		{name: "single tag", raw: "tags: go", want: []string{"go"}},
		{name: "quoted string containing comma", raw: "tags: 'go, rust'", want: []string{"go, rust"}},
		{name: "empty tag", raw: "tags: ''", want: nil},
		{name: "invalid numeric tag", raw: "tags: 42", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := "---\n" + tt.raw + "\n---\nBody"
			metadata, body, err := ParseFrontmatter(content)
			if err != nil {
				t.Fatalf("ParseFrontmatter: %v", err)
			}
			if body != "Body" {
				t.Errorf("body = %q, want Body", body)
			}
			if got := GetStringSlice(metadata, "tags"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tags = %v, want %v", got, tt.want)
			}
			withRaw, _, raw, err := ParseFrontmatterWithRaw(content)
			if err != nil {
				t.Fatalf("ParseFrontmatterWithRaw: %v", err)
			}
			if raw != tt.raw {
				t.Errorf("raw = %q, want %q", raw, tt.raw)
			}
			if got := GetStringSlice(withRaw, "tags"); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tags with raw = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFrontmatter_TagsList(t *testing.T) {
	// Test case: "frontmatter with tags list"
	content := `---
title: Tagged Post
tags: [python, tutorial, beginner]
---
Post content`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify title
	if metadata["title"] != "Tagged Post" {
		t.Errorf("title: got %q, want 'Tagged Post'", metadata["title"])
	}

	// Verify tags
	tags := GetStringSlice(metadata, "tags")
	expectedTags := []string{"python", "tutorial", "beginner"}
	if !reflect.DeepEqual(tags, expectedTags) {
		t.Errorf("tags: got %v, want %v", tags, expectedTags)
	}

	// Verify body
	if body != "Post content" {
		t.Errorf("body: got %q, want 'Post content'", body)
	}
}

func TestParseFrontmatter_NestedObject(t *testing.T) {
	// Test case: "frontmatter with nested object"
	content := `---
title: Custom Template
template:
  name: special.html
  layout: wide
---
Content`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify title
	if metadata["title"] != "Custom Template" {
		t.Errorf("title: got %q, want 'Custom Template'", metadata["title"])
	}

	// Verify nested template object
	templateData, ok := metadata["template"].(map[string]interface{})
	if !ok {
		t.Fatalf("template should be a map, got %T", metadata["template"])
	}
	if templateData["name"] != "special.html" {
		t.Errorf("template.name: got %q, want 'special.html'", templateData["name"])
	}
	if templateData["layout"] != "wide" {
		t.Errorf("template.layout: got %q, want 'wide'", templateData["layout"])
	}

	// Verify body
	if body != "Content" {
		t.Errorf("body: got %q, want 'Content'", body)
	}
}

func TestParseFrontmatter_Missing(t *testing.T) {
	// Test case: "missing frontmatter"
	content := `# Just Markdown
No frontmatter here`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify empty metadata
	if len(metadata) != 0 {
		t.Errorf("metadata should be empty, got %d items", len(metadata))
	}

	// Verify body is full content
	expected := "# Just Markdown\nNo frontmatter here"
	if body != expected {
		t.Errorf("body: got %q, want %q", body, expected)
	}
}

func TestParseFrontmatter_Empty(t *testing.T) {
	// Test case: "empty frontmatter"
	content := `---
---
Content only`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify empty metadata
	if len(metadata) != 0 {
		t.Errorf("metadata should be empty, got %d items", len(metadata))
	}

	// Verify body
	if body != "Content only" {
		t.Errorf("body: got %q, want 'Content only'", body)
	}
}

func TestParseFrontmatter_NormalizeAliases(t *testing.T) {
	content := `---
aliases:
  - alpha
alias: beta
handles:
  - gamma
handle: delta
---
Content`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if body != "Content" {
		t.Errorf("body: got %q, want 'Content'", body)
	}

	aliases := GetStringSlice(metadata, "aliases")
	expected := []string{"alpha", "beta", "gamma", "delta"}
	if !reflect.DeepEqual(aliases, expected) {
		t.Errorf("aliases: got %v, want %v", aliases, expected)
	}
}

func TestParseFrontmatter_BooleanValues(t *testing.T) {
	// Test case: "frontmatter boolean values"
	content := `---
published: true
draft: false
featured: yes
---
Post`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify booleans
	if GetBool(metadata, "published", false) != true {
		t.Errorf("published: expected true")
	}
	if GetBool(metadata, "draft", true) != false {
		t.Errorf("draft: expected false")
	}
	// YAML "yes" is parsed as boolean true
	if GetBool(metadata, "featured", false) != true {
		t.Errorf("featured: expected true (yes should be parsed as true)")
	}

	// Verify body
	if body != "Post" {
		t.Errorf("body: got %q, want 'Post'", body)
	}
}

// =============================================================================
// ExtractFrontmatter Tests
// =============================================================================

func TestExtractFrontmatter_Basic(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantFM   string
		wantBody string
		wantErr  bool
	}{
		{
			name:     "basic frontmatter",
			content:  "---\ntitle: Test\n---\nBody",
			wantFM:   "title: Test",
			wantBody: "Body",
			wantErr:  false,
		},
		{
			name:     "no frontmatter",
			content:  "Just content",
			wantFM:   "",
			wantBody: "Just content",
			wantErr:  false,
		},
		{
			name:     "empty frontmatter",
			content:  "---\n---\nBody",
			wantFM:   "",
			wantBody: "Body",
			wantErr:  false,
		},
		{
			name:     "unclosed frontmatter",
			content:  "---\ntitle: Test\nNo closing",
			wantFM:   "",
			wantBody: "",
			wantErr:  true,
		},
		{
			name:     "not frontmatter start",
			content:  "---something\ncontent",
			wantFM:   "",
			wantBody: "---something\ncontent",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := ExtractFrontmatter(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if fm != tt.wantFM {
					t.Errorf("frontmatter = %q, want %q", fm, tt.wantFM)
				}
				if body != tt.wantBody {
					t.Errorf("body = %q, want %q", body, tt.wantBody)
				}
			}
		})
	}
}

func TestExtractFrontmatter_CRLFLineEndings(t *testing.T) {
	// Test Windows-style line endings
	content := "---\r\ntitle: Test\r\n---\r\nBody"
	fm, body, err := ExtractFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fm != "title: Test" {
		t.Errorf("frontmatter = %q, want 'title: Test'", fm)
	}
	if body != "Body" {
		t.Errorf("body = %q, want 'Body'", body)
	}
}

func TestExtractFrontmatter_StrictDelimiterAndBOM(t *testing.T) {
	frontmatter, body, err := ExtractFrontmatter("\ufeff---\ntitle: Test\n---\nBody")
	if err != nil {
		t.Fatalf("BOM frontmatter returned error: %v", err)
	}
	if frontmatter != "title: Test" || body != "Body" {
		t.Fatalf("BOM result = (%q, %q), want (%q, %q)", frontmatter, body, "title: Test", "Body")
	}

	if _, _, err := ExtractFrontmatter("---\ntitle: Test\n--- \nBody"); err == nil {
		t.Fatal("malformed closing delimiter returned nil error")
	}
}

// =============================================================================
// Helper Function Tests
// =============================================================================

func TestFrontmatter_GetString(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		key      string
		expected string
	}{
		{
			name:     "existing string",
			metadata: map[string]interface{}{"title": "Hello"},
			key:      "title",
			expected: "Hello",
		},
		{
			name:     "missing key",
			metadata: map[string]interface{}{"title": "Hello"},
			key:      "description",
			expected: "",
		},
		{
			name:     "wrong type",
			metadata: map[string]interface{}{"count": 42},
			key:      "count",
			expected: "",
		},
		{
			name:     "nil map",
			metadata: nil,
			key:      "title",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetString(tt.metadata, tt.key)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFrontmatter_GetBool(t *testing.T) {
	tests := []struct {
		name       string
		metadata   map[string]interface{}
		key        string
		defaultVal bool
		expected   bool
	}{
		{
			name:       "true value",
			metadata:   map[string]interface{}{"published": true},
			key:        "published",
			defaultVal: false,
			expected:   true,
		},
		{
			name:       "false value",
			metadata:   map[string]interface{}{"draft": false},
			key:        "draft",
			defaultVal: true,
			expected:   false,
		},
		{
			name:       "string 'true'",
			metadata:   map[string]interface{}{"enabled": "true"},
			key:        "enabled",
			defaultVal: false,
			expected:   true,
		},
		{
			name:       "string 'yes'",
			metadata:   map[string]interface{}{"enabled": "yes"},
			key:        "enabled",
			defaultVal: false,
			expected:   true,
		},
		{
			name:       "string 'false'",
			metadata:   map[string]interface{}{"enabled": "false"},
			key:        "enabled",
			defaultVal: true,
			expected:   false,
		},
		{
			name:       "string 'no'",
			metadata:   map[string]interface{}{"enabled": "no"},
			key:        "enabled",
			defaultVal: true,
			expected:   false,
		},
		{
			name:       "missing key returns default",
			metadata:   map[string]interface{}{},
			key:        "missing",
			defaultVal: true,
			expected:   true,
		},
		{
			name:       "wrong type returns default",
			metadata:   map[string]interface{}{"count": 42},
			key:        "count",
			defaultVal: false,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBool(tt.metadata, tt.key, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFrontmatter_GetStringSlice(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]interface{}
		key      string
		expected []string
	}{
		{
			name:     "[]string type",
			metadata: map[string]interface{}{"tags": []string{"a", "b", "c"}},
			key:      "tags",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "[]interface{} type",
			metadata: map[string]interface{}{"tags": []interface{}{"a", "b", "c"}},
			key:      "tags",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "mixed interface slice",
			metadata: map[string]interface{}{"tags": []interface{}{"a", 42, "c"}},
			key:      "tags",
			expected: []string{"a", "c"}, // Non-strings are skipped
		},
		{
			name:     "missing key",
			metadata: map[string]interface{}{},
			key:      "tags",
			expected: nil,
		},
		{
			name:     "scalar shorthand",
			metadata: map[string]interface{}{"tags": "not-a-slice"},
			key:      "tags",
			expected: []string{"not-a-slice"},
		},
		{
			name:     "empty scalar",
			metadata: map[string]interface{}{"tags": ""},
			key:      "tags",
			expected: nil,
		},
		{
			name:     "wrong type",
			metadata: map[string]interface{}{"tags": 42},
			key:      "tags",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetStringSlice(tt.metadata, tt.key)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Edge Cases and Error Handling
// =============================================================================

func TestParseFrontmatter_InvalidYAML(t *testing.T) {
	// Test case: "invalid frontmatter yaml"
	content := `---
title: "unclosed quote
---
Content`

	_, _, err := ParseFrontmatter(content)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseFrontmatter_MultipleDocuments(t *testing.T) {
	// Content with what might look like multiple YAML documents
	content := `---
title: First
---
Content with --- in it
More content`

	metadata, body, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata["title"] != "First" {
		t.Errorf("title: got %q, want 'First'", metadata["title"])
	}

	// Body should include the --- in content
	if body != "Content with --- in it\nMore content" {
		t.Errorf("body: got %q", body)
	}
}

func TestParseFrontmatter_ListValues(t *testing.T) {
	content := `---
tags:
  - python
  - go
  - rust
---
Content`

	metadata, _, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tags := GetStringSlice(metadata, "tags")
	expected := []string{"python", "go", "rust"}
	if !reflect.DeepEqual(tags, expected) {
		t.Errorf("tags: got %v, want %v", tags, expected)
	}
}

func TestParseFrontmatter_NestedLists(t *testing.T) {
	content := `---
authors:
  - name: Alice
    email: alice@example.com
  - name: Bob
    email: bob@example.com
---
Content`

	metadata, _, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	authors, ok := metadata["authors"].([]interface{})
	if !ok {
		t.Fatalf("authors should be a slice, got %T", metadata["authors"])
	}
	if len(authors) != 2 {
		t.Errorf("expected 2 authors, got %d", len(authors))
	}
}

func TestParseFrontmatter_NumericValues(t *testing.T) {
	content := `---
count: 42
rating: 4.5
---
Content`

	metadata, _, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata["count"] != 42 {
		t.Errorf("count: got %v (%T), want 42", metadata["count"], metadata["count"])
	}
	if metadata["rating"] != 4.5 {
		t.Errorf("rating: got %v (%T), want 4.5", metadata["rating"], metadata["rating"])
	}
}

func TestParseFrontmatter_NullValues(t *testing.T) {
	content := `---
title: Test
description: null
---
Content`

	metadata, _, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata["title"] != "Test" {
		t.Errorf("title: got %q, want 'Test'", metadata["title"])
	}
	if metadata["description"] != nil {
		t.Errorf("description: got %v, want nil", metadata["description"])
	}
}

func TestParseFrontmatter_MultilineStrings(t *testing.T) {
	content := `---
description: |
  This is a multiline
  description that spans
  multiple lines.
---
Content`

	metadata, _, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	desc := GetString(metadata, "description")
	if desc == "" {
		t.Error("description should not be empty")
	}
	// Multiline string should contain newlines
	if len(desc) < 20 {
		t.Errorf("description seems truncated: %q", desc)
	}
}

func TestParseFrontmatter_DateFormats(t *testing.T) {
	content := `---
date1: 2024-01-15
date2: "2024-01-15"
date3: 2024-01-15T10:30:00Z
---
Content`

	metadata, _, err := ParseFrontmatter(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// YAML parser may interpret dates in different ways
	// Just verify they're present and not nil
	if metadata["date1"] == nil {
		t.Error("date1 should be set")
	}
	if metadata["date2"] == nil {
		t.Error("date2 should be set")
	}
	if metadata["date3"] == nil {
		t.Error("date3 should be set")
	}
}
