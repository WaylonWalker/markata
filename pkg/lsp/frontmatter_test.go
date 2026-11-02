package lsp

import (
	"testing"
)

func TestFindFrontmatterBoundaries(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantStart int
		wantEnd   int
	}{
		{
			name:      "valid frontmatter",
			content:   "---\ntitle: Test\n---\n\nContent",
			wantStart: 0,
			wantEnd:   2,
		},
		{
			name:      "no frontmatter",
			content:   "# Title\n\nContent",
			wantStart: -1,
			wantEnd:   -1,
		},
		{
			name:      "unclosed frontmatter",
			content:   "---\ntitle: Test\n",
			wantStart: 0,
			wantEnd:   -1,
		},
		{
			name:      "multiline frontmatter",
			content:   "---\ntitle: Test\ndescription: A description\npublished: true\n---\n\nContent",
			wantStart: 0,
			wantEnd:   4,
		},
		{
			name:      "empty file",
			content:   "",
			wantStart: -1,
			wantEnd:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := splitLines(tt.content)
			gotStart, gotEnd := findFrontmatterBoundaries(lines)
			if gotStart != tt.wantStart {
				t.Errorf("startLine = %d, want %d", gotStart, tt.wantStart)
			}
			if gotEnd != tt.wantEnd {
				t.Errorf("endLine = %d, want %d", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestCollectExistingFields(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantKeys []string
	}{
		{
			name:     "basic fields",
			content:  "---\ntitle: Test\ndescription: A test\n---",
			wantKeys: []string{"title", "description"},
		},
		{
			name:     "with list",
			content:  "---\ntitle: Test\ntags:\n  - tag1\n  - tag2\n---",
			wantKeys: []string{"title", "tags"},
		},
		{
			name:     "empty frontmatter",
			content:  "---\n---",
			wantKeys: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := splitLines(tt.content)
			startLine, endLine := findFrontmatterBoundaries(lines)
			got := collectExistingFields(lines, startLine, endLine)

			if len(got) != len(tt.wantKeys) {
				t.Errorf("got %d fields, want %d", len(got), len(tt.wantKeys))
			}
			for _, key := range tt.wantKeys {
				if !got[key] {
					t.Errorf("missing expected field: %s", key)
				}
			}
		})
	}
}

func TestGetFrontmatterContext(t *testing.T) {
	tests := []struct {
		name             string
		content          string
		line             int
		col              int
		wantInFM         bool
		wantIsFieldName  bool
		wantIsFieldValue bool
		wantCurrentField string
		wantPrefix       string
	}{
		{
			name:            "on field name",
			content:         "---\ntit\n---",
			line:            1,
			col:             3,
			wantInFM:        true,
			wantIsFieldName: true,
			wantPrefix:      "tit",
		},
		{
			name:             "on field value",
			content:          "---\ntitle: Test\n---",
			line:             1,
			col:              14, // Position at end of "Test"
			wantInFM:         true,
			wantIsFieldValue: true,
			wantCurrentField: "title",
			wantPrefix:       "Test",
		},
		{
			name:             "on empty field value",
			content:          "---\ntitle: \n---",
			line:             1,
			col:              7,
			wantInFM:         true,
			wantIsFieldValue: true,
			wantCurrentField: "title",
			wantPrefix:       "",
		},
		{
			name:            "empty line in frontmatter",
			content:         "---\ntitle: Test\n\n---",
			line:            2,
			col:             0,
			wantInFM:        true,
			wantIsFieldName: true,
			wantPrefix:      "",
		},
		{
			name:     "outside frontmatter - before",
			content:  "---\ntitle: Test\n---\n\nContent",
			line:     0,
			col:      0,
			wantInFM: false,
		},
		{
			name:     "outside frontmatter - after",
			content:  "---\ntitle: Test\n---\n\nContent",
			line:     4,
			col:      3,
			wantInFM: false,
		},
		{
			name:     "on closing delimiter",
			content:  "---\ntitle: Test\n---\n\nContent",
			line:     2,
			col:      0,
			wantInFM: false,
		},
		{
			name:             "published field value",
			content:          "---\npublished: tr\n---",
			line:             1,
			col:              14,
			wantInFM:         true,
			wantIsFieldValue: true,
			wantCurrentField: "published",
			wantPrefix:       "tr",
		},
		{
			name:             "in list item",
			content:          "---\ntags:\n  - goals\n---",
			line:             2,
			col:              9,
			wantInFM:         true,
			wantIsFieldValue: true,
			wantCurrentField: "tags",
			wantPrefix:       "goals",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := getFrontmatterContext(tt.content, tt.line, tt.col)

			if ctx.InFrontmatter != tt.wantInFM {
				t.Errorf("InFrontmatter = %v, want %v", ctx.InFrontmatter, tt.wantInFM)
			}
			if ctx.IsFieldName != tt.wantIsFieldName {
				t.Errorf("IsFieldName = %v, want %v", ctx.IsFieldName, tt.wantIsFieldName)
			}
			if ctx.IsFieldValue != tt.wantIsFieldValue {
				t.Errorf("IsFieldValue = %v, want %v", ctx.IsFieldValue, tt.wantIsFieldValue)
			}
			if ctx.CurrentField != tt.wantCurrentField {
				t.Errorf("CurrentField = %q, want %q", ctx.CurrentField, tt.wantCurrentField)
			}
			if ctx.Prefix != tt.wantPrefix {
				t.Errorf("Prefix = %q, want %q", ctx.Prefix, tt.wantPrefix)
			}
			if tt.name == "in list item" && ctx.StartCol != 4 {
				t.Errorf("StartCol = %d, want 4", ctx.StartCol)
			}
		})
	}
}

func TestGetFieldNameCompletions(t *testing.T) {
	ctx := &FrontmatterContext{
		InFrontmatter:  true,
		IsFieldName:    true,
		Prefix:         "",
		StartCol:       0,
		ExistingFields: map[string]bool{},
	}

	params := CompletionParams{
		Position: Position{Line: 1, Character: 0},
	}

	items := getFieldNameCompletions(ctx, params)

	// Should return all frontmatter fields
	if len(items) == 0 {
		t.Error("expected completion items, got none")
	}

	// Check that required fields appear first (sorted)
	foundTitle := false
	foundDate := false
	for _, item := range items {
		if item.Label == "title" {
			foundTitle = true
		}
		if item.Label == "date" {
			foundDate = true
		}
	}

	if !foundTitle {
		t.Error("expected 'title' field in completions")
	}
	if !foundDate {
		t.Error("expected 'date' field in completions")
	}
}

func TestGetFieldNameCompletions_WithPrefix(t *testing.T) {
	ctx := &FrontmatterContext{
		InFrontmatter:  true,
		IsFieldName:    true,
		Prefix:         "pub",
		StartCol:       0,
		ExistingFields: map[string]bool{},
	}

	params := CompletionParams{
		Position: Position{Line: 1, Character: 3},
	}

	items := getFieldNameCompletions(ctx, params)

	// Should filter to fields starting with "pub"
	if len(items) != 1 {
		t.Errorf("expected 1 completion item for prefix 'pub', got %d", len(items))
	}
	if len(items) > 0 && items[0].Label != "published" {
		t.Errorf("expected 'published' field, got %q", items[0].Label)
	}
}

func TestGetFieldNameCompletions_ExcludesExisting(t *testing.T) {
	ctx := &FrontmatterContext{
		InFrontmatter: true,
		IsFieldName:   true,
		Prefix:        "",
		StartCol:      0,
		ExistingFields: map[string]bool{
			"title": true,
			"date":  true,
		},
	}

	params := CompletionParams{
		Position: Position{Line: 1, Character: 0},
	}

	items := getFieldNameCompletions(ctx, params)

	// Should not include title or date
	for _, item := range items {
		if item.Label == "title" {
			t.Error("should not suggest 'title' when it already exists")
		}
		if item.Label == "date" {
			t.Error("should not suggest 'date' when it already exists")
		}
	}
}

func TestGetFieldValueCompletions(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		prefix     string
		wantCount  int
		wantValues []string
	}{
		{
			name:       "published field",
			field:      "published",
			prefix:     "",
			wantCount:  2,
			wantValues: []string{"true", "false"},
		},
		{
			name:       "draft field with prefix",
			field:      "draft",
			prefix:     "t",
			wantCount:  1,
			wantValues: []string{"true"},
		},
		{
			name:       "title field - no predefined values",
			field:      "title",
			prefix:     "",
			wantCount:  0,
			wantValues: []string{},
		},
		{
			name:       "unknown field",
			field:      "unknown_field",
			prefix:     "",
			wantCount:  0,
			wantValues: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &FrontmatterContext{
				InFrontmatter:  true,
				IsFieldValue:   true,
				CurrentField:   tt.field,
				Prefix:         tt.prefix,
				StartCol:       7,
				ExistingFields: map[string]bool{},
			}

			params := CompletionParams{
				Position: Position{Line: 1, Character: 7 + len(tt.prefix)},
			}

			items := getFieldValueCompletions(ctx, params)

			if len(items) != tt.wantCount {
				t.Errorf("got %d items, want %d", len(items), tt.wantCount)
			}

			for i, wantValue := range tt.wantValues {
				if i < len(items) && items[i].Label != wantValue {
					t.Errorf("item %d: got %q, want %q", i, items[i].Label, wantValue)
				}
			}
		})
	}
}

func TestFormatFieldDocumentation(t *testing.T) {
	field := &FrontmatterField{
		Name:         "published",
		Type:         "boolean",
		Description:  "Whether the post is published",
		Required:     false,
		Values:       []string{"true", "false"},
		DefaultValue: "false",
	}

	doc := formatFieldDocumentation(field)

	// Check that documentation contains expected parts
	if doc == "" {
		t.Error("expected non-empty documentation")
	}

	expectedParts := []string{
		"**published**",
		"Whether the post is published",
		"*Type: boolean*",
		"*Allowed values: true, false*",
		"*Default: false*",
	}

	for _, part := range expectedParts {
		if !containsFrontmatter(doc, part) {
			t.Errorf("documentation missing: %q", part)
		}
	}
}

func TestFrontmatterFieldsDefinitions(t *testing.T) {
	// Verify that required fields have proper snippets
	requiredFields := []string{"title", "date"}

	for _, name := range requiredFields {
		found := false
		for _, field := range frontmatterFields {
			if field.Name == name {
				found = true
				if !field.Required {
					t.Errorf("field %q should be marked as required", name)
				}
				if field.Snippet == "" {
					t.Errorf("field %q missing snippet", name)
				}
				break
			}
		}
		if !found {
			t.Errorf("required field %q not found in frontmatterFields", name)
		}
	}

	// Verify boolean fields have values
	boolFields := []string{"published", "draft", "skip", "toc", "sidebar"}
	for _, name := range boolFields {
		for _, field := range frontmatterFields {
			if field.Name == name {
				if field.Type != "boolean" {
					t.Errorf("field %q should have type 'boolean', got %q", name, field.Type)
				}
				if len(field.Values) != 2 {
					t.Errorf("field %q should have 2 values (true/false), got %d", name, len(field.Values))
				}
				break
			}
		}
	}
}

// Helper function to split content into lines
func splitLines(content string) []string {
	if content == "" {
		return []string{}
	}
	result := []string{}
	start := 0
	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			result = append(result, content[start:i])
			start = i + 1
		}
	}
	if start <= len(content) {
		result = append(result, content[start:])
	}
	return result
}

// Helper function to check if string contains substring
// Note: uses strings.Contains - this is a wrapper for readability in tests
func containsFrontmatter(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && containsHelperFrontmatter(s, substr))
}

func containsHelperFrontmatter(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
