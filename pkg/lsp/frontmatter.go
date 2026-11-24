package lsp

import (
	"fmt"
	"sort"
	"strings"
)

// FrontmatterField represents a frontmatter field definition for auto-completion.
type FrontmatterField struct {
	// Name is the field name (e.g., "title", "date")
	Name string

	// Type is the field type for documentation (e.g., "string", "boolean", "date")
	Type string

	// Description is a human-readable description of the field
	Description string

	// Required indicates if the field is required
	Required bool

	// Values contains allowed values for enum-like fields (e.g., ["true", "false"])
	Values []string

	// DefaultValue is the default value if any
	DefaultValue string

	// Snippet is the completion snippet (uses $1, $2 for placeholders)
	Snippet string
}

// frontmatterFields contains all known frontmatter fields for markata-go posts.
var frontmatterFields = []FrontmatterField{
	{
		Name:        "title",
		Type:        "string",
		Description: "The post title displayed in the browser and feeds",
		Required:    true,
		Snippet:     "title: ${1:My Post Title}",
	},
	{
		Name:        "date",
		Type:        "date",
		Description: "Publication date in YYYY-MM-DD format",
		Required:    true,
		Snippet:     "date: ${1:2024-01-01}",
	},
	{
		Name:        "published",
		Type:        "boolean",
		Description: "Whether the post is published (visible on the site)",
		Required:    false,
		Values:      []string{"true", "false"},
		Snippet:     "published: ${1|true,false|}",
	},
	{
		Name:         "draft",
		Type:         "boolean",
		Description:  "Whether the post is a draft (not published)",
		Required:     false,
		Values:       []string{"true", "false"},
		DefaultValue: "false",
		Snippet:      "draft: ${1|true,false|}",
	},
	{
		Name:        "description",
		Type:        "string",
		Description: "Short description for SEO and feed summaries",
		Required:    false,
		Snippet:     "description: ${1:A brief description of the post}",
	},
	{
		Name:        "slug",
		Type:        "string",
		Description: "URL-safe identifier (auto-generated from filename if not set)",
		Required:    false,
		Snippet:     "slug: ${1:my-post-slug}",
	},
	{
		Name:        "tags",
		Type:        "list",
		Description: "List of tags for categorization",
		Required:    false,
		Snippet:     "tags:\n  - ${1:tag1}\n  - ${2:tag2}",
	},
	{
		Name:        "template",
		Type:        "string",
		Description: "Template file to use for rendering (default: post.html)",
		Required:    false,
		Snippet:     "template: ${1:post.html}",
	},
	{
		Name:        "templateKey",
		Type:        "string",
		Description: "Built-in or site template key for this post",
		Required:    false,
		Snippet:     "templateKey: ${1:blog-post}",
	},
	{
		Name:         "skip",
		Type:         "boolean",
		Description:  "Skip this post during processing",
		Required:     false,
		Values:       []string{"true", "false"},
		DefaultValue: "false",
		Snippet:      "skip: ${1|true,false|}",
	},
	{
		Name:        "prevnext_feed",
		Type:        "string",
		Description: "Feed/series slug for prev/next navigation",
		Required:    false,
		Snippet:     "prevnext_feed: ${1:series-name}",
	},
	{
		Name:        "image",
		Type:        "string",
		Description: "Featured image URL for Open Graph and social sharing",
		Required:    false,
		Snippet:     "image: ${1:/images/featured.jpg}",
	},
	{
		Name:        "author",
		Type:        "string",
		Description: "Post author name",
		Required:    false,
		Snippet:     "author: ${1:Author Name}",
	},
	{
		Name:        "canonical_url",
		Type:        "string",
		Description: "Canonical URL if this post is republished from another source",
		Required:    false,
		Snippet:     "canonical_url: ${1:https://example.com/original-post}",
	},
	{
		Name:        "layout",
		Type:        "string",
		Description: "Layout to use for this post (overrides default)",
		Required:    false,
		Snippet:     "layout: ${1|default,wide,full|}",
	},
	{
		Name:        "toc",
		Type:        "boolean",
		Description: "Enable table of contents for this post",
		Required:    false,
		Values:      []string{"true", "false"},
		Snippet:     "toc: ${1|true,false|}",
	},
	{
		Name:        "sidebar",
		Type:        "boolean",
		Description: "Enable sidebar for this post",
		Required:    false,
		Values:      []string{"true", "false"},
		Snippet:     "sidebar: ${1|true,false|}",
	},
}

// FrontmatterContext contains information about the cursor position within frontmatter.
type FrontmatterContext struct {
	// InFrontmatter indicates if the cursor is within the frontmatter section
	InFrontmatter bool

	// IsFieldName indicates if the cursor is in a position for a field name
	IsFieldName bool

	// IsFieldValue indicates if the cursor is in a position for a field value
	IsFieldValue bool

	// CurrentField is the field name if we're editing a value
	CurrentField string

	// Prefix is the text before the cursor (for filtering)
	Prefix string

	// StartCol is the column where the prefix starts
	StartCol int

	// ExistingFields contains field names already present in the frontmatter
	ExistingFields map[string]bool
}

// getFrontmatterContext analyzes the document to determine frontmatter context.
// Returns information about whether the cursor is in frontmatter and what kind of completion is needed.
func getFrontmatterContext(content string, line, col int) *FrontmatterContext {
	lines := strings.Split(content, "\n")

	// Find frontmatter boundaries
	startLine, endLine := findFrontmatterBoundaries(lines)
	if startLine == -1 || endLine == -1 {
		return &FrontmatterContext{InFrontmatter: false}
	}

	// Check if cursor is within frontmatter
	if line <= startLine || line >= endLine {
		return &FrontmatterContext{InFrontmatter: false}
	}

	// Collect existing fields
	existingFields := collectExistingFields(lines, startLine, endLine)

	// Get the current line content
	if line >= len(lines) {
		return &FrontmatterContext{InFrontmatter: true, ExistingFields: existingFields}
	}

	currentLine := lines[line]
	if col > len(currentLine) {
		col = len(currentLine)
	}

	// Analyze what we're completing
	return analyzeLineContext(lines, line, currentLine, col, existingFields)
}

// findFrontmatterBoundaries finds the start and end lines of frontmatter.
// Returns (-1, -1) if no valid frontmatter is found.
func findFrontmatterBoundaries(lines []string) (startLine, endLine int) {
	startLine = -1
	endLine = -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "---" {
			if startLine == -1 {
				startLine = i
			} else {
				endLine = i
				break
			}
		}
	}

	return startLine, endLine
}

// collectExistingFields collects field names already present in the frontmatter.
func collectExistingFields(lines []string, startLine, endLine int) map[string]bool {
	existingFields := make(map[string]bool)

	for i := startLine + 1; i < endLine && i < len(lines); i++ {
		line := lines[i]
		// Match field name at the start of line (not indented)
		if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") {
			if idx := strings.Index(line, ":"); idx > 0 {
				fieldName := strings.TrimSpace(line[:idx])
				if fieldName != "" {
					existingFields[fieldName] = true
				}
			}
		}
	}

	return existingFields
}

// analyzeLineContext determines what kind of completion is needed on the current line.
func analyzeLineContext(lines []string, lineNumber int, line string, col int, existingFields map[string]bool) *FrontmatterContext {
	ctx := &FrontmatterContext{
		InFrontmatter:  true,
		ExistingFields: existingFields,
	}

	textBeforeCursor := ""
	if col <= len(line) {
		textBeforeCursor = line[:col]
	}

	// Check if we're in a list item (indented with - )
	trimmedLine := strings.TrimLeft(line, " \t")
	if strings.HasPrefix(trimmedLine, "- ") {
		// Keep the parent key for completions in YAML list values (for example,
		// tags: followed by - go). Only the current item is replaced.
		for previous := lineNumber - 1; previous >= 0; previous-- {
			parent := lines[previous]
			if strings.TrimSpace(parent) == "" {
				continue
			}
			if len(parent) == len(strings.TrimLeft(parent, " \t")) {
				if colon := strings.IndexByte(parent, ':'); colon > 0 {
					ctx.CurrentField = strings.TrimSpace(parent[:colon])
					ctx.IsFieldValue = true
					itemStart := len(line) - len(trimmedLine) + 2
					for itemStart < len(line) && (line[itemStart] == ' ' || line[itemStart] == '\t') {
						itemStart++
					}
					ctx.StartCol = itemStart
					if col > itemStart {
						ctx.Prefix = strings.TrimSpace(line[itemStart:col])
					}
				}
				break
			}
			// The nearest indented non-list value ends the parent list.
			if !strings.HasPrefix(strings.TrimLeft(parent, " \t"), "-") {
				break
			}
		}
		return ctx
	}

	// Check if this line has a colon
	colonIdx := strings.Index(textBeforeCursor, ":")
	if colonIdx == -1 {
		// No colon yet - completing a field name
		ctx.IsFieldName = true
		ctx.Prefix = strings.TrimSpace(textBeforeCursor)
		// Start column is after any leading whitespace
		ctx.StartCol = len(line) - len(strings.TrimLeft(line, " \t"))
		return ctx
	}

	// We have a colon - check if we're before or after it
	if col <= colonIdx {
		// Cursor is before the colon - completing field name
		ctx.IsFieldName = true
		ctx.Prefix = strings.TrimSpace(textBeforeCursor)
		ctx.StartCol = len(line) - len(strings.TrimLeft(line, " \t"))
		return ctx
	}

	// Cursor is after the colon - completing field value
	ctx.IsFieldValue = true
	fieldName := strings.TrimSpace(line[:colonIdx])
	ctx.CurrentField = fieldName

	// Get the value part after the colon
	valueStart := colonIdx + 1
	if valueStart < col {
		ctx.Prefix = strings.TrimSpace(line[valueStart:col])
		ctx.StartCol = valueStart
		// Skip leading space after colon
		if valueStart < len(line) && line[valueStart] == ' ' {
			ctx.StartCol = valueStart + 1
		}
	}

	return ctx
}

// getFrontmatterCompletions returns completion items for frontmatter.
func getFrontmatterCompletions(ctx *FrontmatterContext, params CompletionParams) []CompletionItem {
	var items []CompletionItem

	if ctx.IsFieldName {
		items = getFieldNameCompletions(ctx, params)
	} else if ctx.IsFieldValue {
		items = getFieldValueCompletions(ctx, params)
	}

	return items
}

// getFieldNameCompletions returns completions for field names.
func getFieldNameCompletions(ctx *FrontmatterContext, params CompletionParams) []CompletionItem {
	items := make([]CompletionItem, 0, len(frontmatterFields))
	prefix := strings.ToLower(ctx.Prefix)

	// Sort fields: required first, then by name
	sortedFields := make([]FrontmatterField, len(frontmatterFields))
	copy(sortedFields, frontmatterFields)
	sort.Slice(sortedFields, func(i, j int) bool {
		if sortedFields[i].Required != sortedFields[j].Required {
			return sortedFields[i].Required
		}
		return sortedFields[i].Name < sortedFields[j].Name
	})

	for i, field := range sortedFields {
		// Skip fields that already exist
		if ctx.ExistingFields[field.Name] {
			continue
		}

		// Filter by prefix
		if prefix != "" && !strings.HasPrefix(strings.ToLower(field.Name), prefix) {
			continue
		}

		// Build documentation
		doc := formatFieldDocumentation(&field)

		// Determine detail text
		detail := field.Type
		if field.Required {
			detail += " (required)"
		}

		item := CompletionItem{
			Label:  field.Name,
			Kind:   CompletionItemKindProperty,
			Detail: detail,
			Documentation: &MarkupContent{
				Kind:  "markdown",
				Value: doc,
			},
			InsertText:       field.Snippet,
			InsertTextFormat: InsertTextFormatSnippet,
			FilterText:       field.Name,
			SortText:         fmt.Sprintf("%d%s", boolToInt(!field.Required), field.Name),
		}

		// Use TextEdit to replace prefix if we have one
		if ctx.Prefix != "" {
			item.TextEdit = &TextEdit{
				Range: Range{
					Start: Position{Line: params.Position.Line, Character: ctx.StartCol},
					End:   Position{Line: params.Position.Line, Character: params.Position.Character},
				},
				NewText: field.Snippet,
			}
		}

		items = append(items, item)

		// Limit to reasonable number
		if len(items) >= 20 {
			break
		}

		_ = i // avoid unused variable warning
	}

	return items
}

// getFieldValueCompletions returns completions for field values.
func getFieldValueCompletions(ctx *FrontmatterContext, params CompletionParams) []CompletionItem {
	var items []CompletionItem

	// Find the field definition
	var field *FrontmatterField
	for i := range frontmatterFields {
		if frontmatterFields[i].Name == ctx.CurrentField {
			field = &frontmatterFields[i]
			break
		}
	}

	if field == nil {
		return items
	}

	// If the field has predefined values, suggest them
	if len(field.Values) > 0 {
		prefix := strings.ToLower(ctx.Prefix)
		for i, value := range field.Values {
			if prefix != "" && !strings.HasPrefix(strings.ToLower(value), prefix) {
				continue
			}

			item := CompletionItem{
				Label:            value,
				Kind:             CompletionItemKindValue,
				Detail:           fmt.Sprintf("Value for %s", field.Name),
				InsertText:       value,
				InsertTextFormat: InsertTextFormatPlainText,
				SortText:         fmt.Sprintf("%02d", i),
			}

			// Use TextEdit to replace prefix
			if ctx.Prefix != "" {
				item.TextEdit = &TextEdit{
					Range: Range{
						Start: Position{Line: params.Position.Line, Character: ctx.StartCol},
						End:   Position{Line: params.Position.Line, Character: params.Position.Character},
					},
					NewText: value,
				}
			}

			items = append(items, item)
		}
	}

	return items
}

// formatFieldDocumentation formats field info for display in completion documentation.
func formatFieldDocumentation(field *FrontmatterField) string {
	var sb strings.Builder

	sb.WriteString("**")
	sb.WriteString(field.Name)
	sb.WriteString("**")
	if field.Required {
		sb.WriteString(" (required)")
	}
	sb.WriteString("\n\n")

	sb.WriteString(field.Description)
	sb.WriteString("\n\n")

	sb.WriteString("*Type: ")
	sb.WriteString(field.Type)
	sb.WriteString("*")

	if len(field.Values) > 0 {
		sb.WriteString("\n\n*Allowed values: ")
		sb.WriteString(strings.Join(field.Values, ", "))
		sb.WriteString("*")
	}

	if field.DefaultValue != "" {
		sb.WriteString("\n\n*Default: ")
		sb.WriteString(field.DefaultValue)
		sb.WriteString("*")
	}

	return sb.String()
}

// boolToInt converts a boolean to int (0 for false, 1 for true).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
