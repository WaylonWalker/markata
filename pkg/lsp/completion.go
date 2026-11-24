package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// CompletionParams contains the parameters for textDocument/completion.
type CompletionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
	Context      *CompletionContext     `json:"context,omitempty"`
}

// CompletionContext contains additional information about the context.
type CompletionContext struct {
	TriggerKind      int     `json:"triggerKind"`
	TriggerCharacter *string `json:"triggerCharacter,omitempty"`
}

// CompletionList represents a collection of completion items.
type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

// CompletionItem represents a completion suggestion.
type CompletionItem struct {
	Label            string         `json:"label"`
	Kind             int            `json:"kind,omitempty"`
	Detail           string         `json:"detail,omitempty"`
	Documentation    *MarkupContent `json:"documentation,omitempty"`
	InsertText       string         `json:"insertText,omitempty"`
	InsertTextFormat int            `json:"insertTextFormat,omitempty"`
	TextEdit         *TextEdit      `json:"textEdit,omitempty"`
	FilterText       string         `json:"filterText,omitempty"`
	SortText         string         `json:"sortText,omitempty"`
	Data             interface{}    `json:"data,omitempty"`
	AdditionalEdits  []TextEdit     `json:"additionalTextEdits,omitempty"`
	CommitCharacters []string       `json:"commitCharacters,omitempty"`
}

// CompletionItemKind defines types of completion items.
const (
	CompletionItemKindText          = 1
	CompletionItemKindMethod        = 2
	CompletionItemKindFunction      = 3
	CompletionItemKindConstructor   = 4
	CompletionItemKindField         = 5
	CompletionItemKindVariable      = 6
	CompletionItemKindClass         = 7
	CompletionItemKindInterface     = 8
	CompletionItemKindModule        = 9
	CompletionItemKindProperty      = 10
	CompletionItemKindUnit          = 11
	CompletionItemKindValue         = 12
	CompletionItemKindEnum          = 13
	CompletionItemKindKeyword       = 14
	CompletionItemKindSnippet       = 15
	CompletionItemKindColor         = 16
	CompletionItemKindFile          = 17
	CompletionItemKindReference     = 18
	CompletionItemKindFolder        = 19
	CompletionItemKindEnumMember    = 20
	CompletionItemKindConstant      = 21
	CompletionItemKindStruct        = 22
	CompletionItemKindEvent         = 23
	CompletionItemKindOperator      = 24
	CompletionItemKindTypeParameter = 25
)

// InsertTextFormat constants.
const (
	InsertTextFormatPlainText = 1
	InsertTextFormatSnippet   = 2
)

// handleCompletion handles textDocument/completion requests.
func (s *Server) handleCompletion(_ context.Context, msg *Message) error {
	var params CompletionParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		return s.sendError(msg.ID, InvalidParams, "invalid completion params")
	}

	// Get the document
	s.docMu.RLock()
	doc, ok := s.documents[params.TextDocument.URI]
	s.docMu.RUnlock()

	if !ok {
		return s.sendCompletion(msg.ID, params, &CompletionList{Items: []CompletionItem{}})
	}

	// Get the line at the cursor position
	lines := strings.Split(doc.Content, "\n")
	if params.Position.Line < 0 || params.Position.Character < 0 || params.Position.Line >= len(lines) {
		return s.sendCompletion(msg.ID, params, &CompletionList{Items: []CompletionItem{}})
	}

	line := lines[params.Position.Line]
	col := byteColumn(line, params.Position.Character)
	params.Position.Character = col
	if col > len(line) {
		col = len(line)
	}
	if strings.HasSuffix(params.TextDocument.URI, ".toml") || strings.HasSuffix(params.TextDocument.URI, ".yaml") || strings.HasSuffix(params.TextDocument.URI, ".yml") || strings.HasSuffix(params.TextDocument.URI, ".json") {
		if items := getThemeConfigCompletions(doc.Content, params.Position.Line, line[:col]); items != nil {
			return s.sendCompletion(msg.ID, params, &CompletionList{Items: items})
		}
	}

	// Check if we're inside frontmatter
	frontmatterCtx := getFrontmatterContext(doc.Content, params.Position.Line, col)
	if frontmatterCtx.InFrontmatter && (frontmatterCtx.IsFieldName || frontmatterCtx.IsFieldValue) {
		items := getFrontmatterCompletions(frontmatterCtx, params)
		if frontmatterCtx.IsFieldValue {
			items = append(items, s.workspaceFrontmatterValues(frontmatterCtx, params, line)...)
		}
		return s.sendCompletion(msg.ID, params, &CompletionList{
			IsIncomplete: false,
			Items:        items,
		})
	}

	// Check if we're in an admonition context (after !!!, ???, or ???+)
	if adCtx, inAdmonition := getAdmonitionContext(line, col); inAdmonition {
		items := getAdmonitionCompletions(adCtx, params)
		return s.sendCompletion(msg.ID, params, &CompletionList{Items: items})
	}

	if handled, err := s.handleMarkdownCompletion(msg, params, lines, line, col); handled {
		return err
	}

	// Check if we're inside a mention (@handle)
	if prefix, startCol, inMention := getMentionContext(line, col); inMention {
		return s.handleMentionCompletion(msg, params, prefix, startCol, col)
	}

	if prefix, start, ok := getDocumentLinkContext(line, col); ok {
		return s.handleDocumentLinkCompletion(msg, params, prefix, start, col)
	}

	// Check if we're inside a wikilink
	prefix, startCol, inWikilink := getWikilinkContext(line, col)
	if !inWikilink {
		return s.sendCompletion(msg.ID, params, &CompletionList{Items: []CompletionItem{}})
	}

	return s.handleWikilinkCompletion(msg, params, prefix, startCol, col)
}

func (s *Server) workspaceFrontmatterValues(ctx *FrontmatterContext, params CompletionParams, line string) []CompletionItem {
	if s.index == nil {
		return nil
	}
	values := make(map[string]struct{})
	for _, post := range s.index.AllPosts() {
		if ctx.CurrentField == "tags" {
			for _, tag := range extractTagsFromMetadata(post.Metadata["tags"]) {
				cleanTag := strings.TrimSpace(tag)
				if cleanTag != "" && !strings.EqualFold(cleanTag, "(No visible text)") {
					values[tag] = struct{}{}
				}
			}
			continue
		}
		if value, ok := post.Metadata[ctx.CurrentField]; ok {
			if text, ok := value.(string); ok && text != "" {
				values[text] = struct{}{}
			}
		}
	}
	start := ctx.StartCol
	if start == 0 {
		if colon := strings.IndexByte(line, ':'); colon >= 0 {
			start = colon + 1
			for start < len(line) && (line[start] == ' ' || line[start] == '\t') {
				start++
			}
		}
	}
	items := make([]CompletionItem, 0, len(values))
	for value := range values {
		if strings.TrimSpace(value) == "" {
			continue
		}
		if ctx.Prefix != "" && !strings.HasPrefix(strings.ToLower(value), strings.ToLower(ctx.Prefix)) {
			continue
		}
		items = append(items, CompletionItem{Label: value, Kind: CompletionItemKindValue,
			Detail: "Value used in this site", InsertText: value, InsertTextFormat: InsertTextFormatPlainText,
			TextEdit: &TextEdit{Range: Range{Start: Position{Line: params.Position.Line, Character: start}, End: Position{Line: params.Position.Line, Character: params.Position.Character}}, NewText: value}})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Label < items[j].Label })
	return items
}

// handleWikilinkCompletion handles completion for [[wikilinks]].
func (s *Server) handleWikilinkCompletion(msg *Message, params CompletionParams, prefix string, startCol, col int) error {
	// Get all posts
	posts := s.index.AllPosts()
	if s.browser != nil {
		filtered := make([]*PostInfo, 0, len(posts))
		for _, post := range posts {
			if post.URI != s.browser.DocumentURI {
				filtered = append(filtered, post)
			}
		}
		posts = filtered
	}

	// Build completion entries - include both slugs and aliases
	entries := buildWikilinkCompletionEntries(posts, prefix)

	// Sort: slugs first, then by title/alias name
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].sortKey < entries[j].sortKey
	})

	// Build completion items
	items := make([]CompletionItem, 0, len(entries))
	for i, entry := range entries {
		item := buildWikilinkCompletionItem(entry, i, params, prefix, startCol, col)
		items = append(items, item)
	}

	result := &CompletionList{
		IsIncomplete: false,
		Items:        items,
	}

	return s.sendCompletion(msg.ID, params, result)
}

// wikilinkCompletionEntry represents a potential wikilink completion.
type wikilinkCompletionEntry struct {
	label   string
	post    *PostInfo
	isAlias bool
	sortKey string
}

// buildWikilinkCompletionEntries builds completion entries for all posts including aliases.
func buildWikilinkCompletionEntries(posts []*PostInfo, prefix string) []wikilinkCompletionEntry {
	// Pre-calculate capacity: one entry per post (for slug) plus aliases
	totalAliases := 0
	for _, post := range posts {
		totalAliases += len(post.Aliases)
	}
	entries := make([]wikilinkCompletionEntry, 0, len(posts)+totalAliases)

	for _, post := range posts {
		// Add the slug as a completion entry
		entries = append(entries, wikilinkCompletionEntry{
			label:   post.Slug,
			post:    post,
			isAlias: false,
			sortKey: "0" + strings.ToLower(post.Title), // Slugs sort first
		})

		// Add each alias as a separate completion entry
		for _, alias := range post.Aliases {
			entries = append(entries, wikilinkCompletionEntry{
				label:   alias,
				post:    post,
				isAlias: true,
				sortKey: "1" + strings.ToLower(alias), // Aliases sort after slugs
			})
		}
	}

	// Filter by prefix if provided
	if prefix != "" {
		prefixLower := strings.ToLower(prefix)
		filtered := make([]wikilinkCompletionEntry, 0, len(entries))
		for _, entry := range entries {
			if strings.HasPrefix(strings.ToLower(entry.label), prefixLower) ||
				strings.Contains(strings.ToLower(entry.post.Title), prefixLower) {
				filtered = append(filtered, entry)
			}
		}
		return filtered
	}

	return entries
}

// buildWikilinkCompletionItem creates a CompletionItem from an entry.
func buildWikilinkCompletionItem(entry wikilinkCompletionEntry, index int, params CompletionParams, prefix string, startCol, col int) CompletionItem {
	// Create completion item
	detail := entry.post.Title
	if entry.isAlias {
		detail = fmt.Sprintf("(alias for: %s)", entry.post.Slug)
	}

	// Build filter text - include slug, title, and aliases for better matching
	filterParts := []string{entry.label, entry.post.Title}
	if !entry.isAlias {
		filterParts = append(filterParts, entry.post.Aliases...)
	}

	item := CompletionItem{
		Label:  entry.label,
		Kind:   CompletionItemKindReference,
		Detail: detail,
		Documentation: &MarkupContent{
			Kind:  "markdown",
			Value: formatPostDocumentation(entry.post),
		},
		InsertText:       entry.label,
		InsertTextFormat: InsertTextFormatPlainText,
		FilterText:       strings.Join(filterParts, " "),
		SortText:         fmt.Sprintf("%05d", index), // Preserve sort order
	}

	// If we have a prefix, use TextEdit to replace it
	if prefix != "" {
		item.TextEdit = &TextEdit{
			Range: Range{
				Start: Position{Line: params.Position.Line, Character: startCol},
				End:   Position{Line: params.Position.Line, Character: col},
			},
			NewText: entry.label,
		}
	}

	return item
}

// handleMentionCompletion handles completion for @mentions.
func (s *Server) handleMentionCompletion(msg *Message, params CompletionParams, prefix string, startCol, col int) error {
	// Get matching mentions
	var mentions []*MentionInfo
	if prefix == "" {
		mentions = s.index.AllMentions()
	} else {
		mentions = s.index.SearchMentions(prefix)
	}

	// Sort by handle for consistent ordering
	sort.Slice(mentions, func(i, j int) bool {
		return mentions[i].Handle < mentions[j].Handle
	})

	// Build completion items
	items := make([]CompletionItem, 0, len(mentions))
	for i, mention := range mentions {
		// Create completion item
		item := CompletionItem{
			Label:  "@" + mention.Handle,
			Kind:   CompletionItemKindReference,
			Detail: mention.Title,
			Documentation: &MarkupContent{
				Kind:  "markdown",
				Value: formatMentionDocumentation(mention),
			},
			InsertText:       mention.Handle,
			InsertTextFormat: InsertTextFormatPlainText,
			FilterText:       mention.Handle + " " + mention.Title + " " + strings.Join(mention.Aliases, " "),
			SortText:         fmt.Sprintf("%05d", i), // Preserve sort order
		}

		// Use TextEdit to replace the prefix (everything after @)
		item.TextEdit = &TextEdit{
			Range: Range{
				Start: Position{Line: params.Position.Line, Character: startCol},
				End:   Position{Line: params.Position.Line, Character: col},
			},
			NewText: mention.Handle,
		}

		items = append(items, item)
	}

	result := &CompletionList{
		IsIncomplete: false,
		Items:        items,
	}

	return s.sendCompletion(msg.ID, params, result)
}

// mentionHandleRegex matches the handle part after @.
var mentionHandleRegex = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_.-]*$`)

// getMentionContext checks if the cursor is inside a mention and returns the prefix.
// Returns (prefix, startColumn, isInMention).
func getMentionContext(line string, col int) (prefix string, startCol int, inMention bool) {
	if col > len(line) {
		col = len(line)
	}

	textBeforeCursor := line[:col]

	// Find the last @ that could start a mention
	atIdx := strings.LastIndex(textBeforeCursor, "@")
	if atIdx == -1 {
		return "", 0, false
	}

	// Check that @ is preceded by a valid boundary (start of line, space, or non-word char)
	// and not by another @ or word character
	if atIdx > 0 {
		prevChar := textBeforeCursor[atIdx-1]
		// If preceded by @ or word character, not a valid mention
		if prevChar == '@' || (prevChar >= 'a' && prevChar <= 'z') ||
			(prevChar >= 'A' && prevChar <= 'Z') || (prevChar >= '0' && prevChar <= '9') || prevChar == '_' {
			return "", 0, false
		}
	}

	// Get the text after @
	afterAt := textBeforeCursor[atIdx+1:]

	// Check if there's a valid handle pattern after @
	// Handle must start with a letter
	match := mentionHandleRegex.FindStringIndex(afterAt)
	if len(match) < 2 || match[0] != 0 {
		// No valid handle starting right after @
		return "", 0, false
	}

	// Return the handle prefix
	prefix = afterAt[match[0]:match[1]]
	startCol = atIdx + 1 // Position right after @

	return prefix, startCol, true
}

// wikilinkStartRegex matches [[ at the beginning of a wikilink.
var wikilinkStartRegex = regexp.MustCompile(`\[\[([^\]|]*)$`)

// getWikilinkContext checks if the cursor is inside a wikilink and returns the prefix.
// Returns (prefix, startColumn, isInWikilink).
func getWikilinkContext(line string, col int) (prefix string, startCol int, inWikilink bool) {
	if col > len(line) {
		col = len(line)
	}

	// Look for [[ before the cursor
	textBeforeCursor := line[:col]

	// Check if we're inside a wikilink (after [[ but not after ]])
	match := wikilinkStartRegex.FindStringSubmatchIndex(textBeforeCursor)
	if match == nil {
		return "", 0, false
	}

	// match[0:2] is the full match position
	// match[2:4] is the captured slug prefix position
	startCol = match[2] // Start of the slug prefix
	prefix = textBeforeCursor[startCol:]

	// Check if we're in the display text part (after |)
	if strings.Contains(prefix, "|") {
		return "", 0, false
	}

	return prefix, startCol, true
}

// formatPostDocumentation formats post info for display in completion documentation.
func formatPostDocumentation(post *PostInfo) string {
	var sb strings.Builder

	sb.WriteString("**")
	sb.WriteString(post.Title)
	sb.WriteString("**\n\n")

	if post.Description != "" {
		sb.WriteString(post.Description)
		sb.WriteString("\n\n")
	}

	sb.WriteString("*Path: ")
	sb.WriteString(post.Path)
	sb.WriteString("*")

	return sb.String()
}

// formatMentionDocumentation formats mention info for display in completion documentation.
func formatMentionDocumentation(mention *MentionInfo) string {
	var sb strings.Builder

	sb.WriteString("**@")
	sb.WriteString(mention.Handle)
	sb.WriteString("**")

	if mention.Title != "" {
		sb.WriteString(" - ")
		sb.WriteString(mention.Title)
	}
	sb.WriteString("\n\n")

	if mention.Description != "" {
		sb.WriteString(mention.Description)
		sb.WriteString("\n\n")
	}

	if mention.SiteURL != "" {
		sb.WriteString("*Site: ")
		sb.WriteString(mention.SiteURL)
		sb.WriteString("*\n")
	}

	if len(mention.Aliases) > 0 {
		sb.WriteString("*Aliases: @")
		sb.WriteString(strings.Join(mention.Aliases, ", @"))
		sb.WriteString("*")
	}

	return sb.String()
}

// mentionAtPositionRegex matches @handle patterns for position detection.
var mentionAtPositionRegex = regexp.MustCompile(`@([a-zA-Z][a-zA-Z0-9_.-]*)`)

// getMentionAtPosition returns the mention handle at the given cursor position.
// Returns (handle without @, range, found).
func getMentionAtPosition(line string, col, lineNum int) (string, *Range) {
	// Find all mentions on this line
	matches := mentionAtPositionRegex.FindAllStringSubmatchIndex(line, -1)

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		// match[0:2] is the full match position (including @)
		// match[2:4] is the handle group position (without @)
		start := match[0]
		end := match[1]

		// Check if cursor is within this mention
		if col >= start && col <= end {
			// Validate that @ is at a valid boundary
			if start > 0 {
				prevChar := line[start-1]
				// If preceded by word character, not a valid mention
				if (prevChar >= 'a' && prevChar <= 'z') ||
					(prevChar >= 'A' && prevChar <= 'Z') ||
					(prevChar >= '0' && prevChar <= '9') || prevChar == '_' || prevChar == '@' {
					continue
				}
			}

			handle := line[match[2]:match[3]]
			mentionRange := &Range{
				Start: Position{Line: lineNum, Character: start},
				End:   Position{Line: lineNum, Character: end},
			}
			return handle, mentionRange
		}
	}

	return "", nil
}

// sendCompletion converts byte-based parser ranges to LSP UTF-16 ranges.
func (s *Server) sendCompletion(id json.RawMessage, params CompletionParams, result *CompletionList) error {
	// Browser adapters request again as the prefix changes. Keep large workspaces responsive.
	if s.browser != nil && len(result.Items) > 200 {
		result.Items = result.Items[:200]
		result.IsIncomplete = true
	}
	s.docMu.RLock()
	doc := s.documents[params.TextDocument.URI]
	s.docMu.RUnlock()
	if doc != nil {
		lines := strings.Split(doc.Content, "\n")
		for i := range result.Items {
			edit := result.Items[i].TextEdit
			if edit == nil {
				continue
			}
			if edit.Range.Start.Line >= 0 && edit.Range.Start.Line < len(lines) {
				edit.Range.Start.Character = utf16Column(lines[edit.Range.Start.Line], edit.Range.Start.Character)
			}
			if edit.Range.End.Line >= 0 && edit.Range.End.Line < len(lines) {
				edit.Range.End.Character = utf16Column(lines[edit.Range.End.Line], edit.Range.End.Character)
			}
		}
	}
	return s.sendResponse(id, result)
}

var documentLinkStart = regexp.MustCompile(`!?\[[^\]\n]*\]\(([^)\s]*)$`)

func getDocumentLinkContext(line string, col int) (prefix string, startCol int, inLink bool) {
	match := documentLinkStart.FindStringSubmatchIndex(line[:col])
	if match == nil {
		return "", 0, false
	}
	prefix = line[match[2]:col]
	// Remote URLs and fragment-only links are not document path completions.
	if strings.ContainsAny(prefix, ":#") {
		return "", 0, false
	}
	return prefix, match[2], true
}

func (s *Server) handleDocumentLinkCompletion(msg *Message, params CompletionParams, prefix string, start, col int) error {
	line := s.documentLine(params.TextDocument.URI, params.Position.Line)
	if isManagedImageContext(line, col) {
		return s.handleImageLinkCompletion(msg, params, prefix, start, col)
	}
	if prefix == "" && isEmbedImageContext(line, col) {
		return s.sendCompletion(msg.ID, params, &CompletionList{Items: []CompletionItem{}})
	}
	items := make([]CompletionItem, 0)
	source := uriToPath(params.TextDocument.URI)
	// The browser draft has no source file. Its document paths start at the workspace root.
	directory := filepath.Dir(source)
	if s.browser != nil {
		directory = uriToPath(s.browser.RootURI)
	}
	for _, post := range s.index.AllPosts() {
		if post.URI == params.TextDocument.URI {
			continue
		}
		relative, err := filepath.Rel(directory, post.Path)
		if err != nil {
			continue
		}
		target := (&url.URL{Path: filepath.ToSlash(relative)}).EscapedPath()
		if strings.HasPrefix(prefix, "/") && s.rootURI != "" {
			relative, err = filepath.Rel(uriToPath(s.rootURI), post.Path)
			if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				continue
			}
			target = "/" + (&url.URL{Path: filepath.ToSlash(relative)}).EscapedPath()
		} else if strings.HasPrefix(prefix, "./") {
			target = "./" + target
		}
		if prefix != "" && !strings.HasPrefix(strings.ToLower(target), strings.ToLower(prefix)) && !strings.Contains(strings.ToLower(post.Title), strings.ToLower(prefix)) {
			continue
		}
		items = append(items, CompletionItem{Label: target, Kind: CompletionItemKindFile, Detail: post.Title,
			Documentation: &MarkupContent{Kind: "markdown", Value: formatPostDocumentation(post)},
			TextEdit:      &TextEdit{Range: Range{Start: Position{Line: params.Position.Line, Character: start}, End: Position{Line: params.Position.Line, Character: col}}, NewText: target}})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Label < items[j].Label })
	return s.sendCompletion(msg.ID, params, &CompletionList{Items: items})
}

func (s *Server) documentLine(uri string, lineNumber int) string {
	s.docMu.RLock()
	doc := s.documents[uri]
	s.docMu.RUnlock()
	if doc == nil {
		return ""
	}
	lines := strings.Split(doc.Content, "\n")
	if lineNumber < 0 || lineNumber >= len(lines) {
		return ""
	}
	return lines[lineNumber]
}

func isManagedImageContext(line string, col int) bool {
	match := documentLinkStart.FindStringSubmatchIndex(line[:min(col, len(line))])
	if len(match) < 2 || !strings.HasPrefix(line[match[0]:match[1]], "!") {
		return false
	}
	full := line[match[0]:match[1]]
	linkClose := strings.Index(full, "](")
	open := strings.IndexByte(full, '[')
	return linkClose >= 0 && open >= 0 && !strings.EqualFold(strings.TrimSpace(full[open+1:linkClose]), "embed")
}

func isEmbedImageContext(line string, col int) bool {
	match := documentLinkStart.FindStringSubmatchIndex(line[:min(col, len(line))])
	if len(match) < 2 || !strings.HasPrefix(line[match[0]:match[1]], "!") {
		return false
	}
	full := line[match[0]:match[1]]
	linkClose := strings.Index(full, "](")
	open := strings.IndexByte(full, '[')
	return linkClose >= 0 && open >= 0 && strings.EqualFold(strings.TrimSpace(full[open+1:linkClose]), "embed")
}

func (s *Server) handleImageLinkCompletion(msg *Message, params CompletionParams, prefix string, start, col int) error {
	items := make([]CompletionItem, 0)
	for i := range s.images {
		image := &s.images[i]
		if prefix != "" && !strings.HasPrefix(strings.ToLower(image.Src), strings.ToLower(prefix)) {
			continue
		}
		detail := image.MIMEType
		if image.Width > 0 && image.Height > 0 {
			detail = fmt.Sprintf("%s · %d × %d", detail, image.Width, image.Height)
		}
		items = append(items, CompletionItem{Label: image.Src, Kind: CompletionItemKindFile,
			Detail: detail, Documentation: &MarkupContent{Kind: "markdown", Value: "**Alt text:** " + image.Alt},
			TextEdit: &TextEdit{Range: Range{Start: Position{Line: params.Position.Line, Character: start}, End: Position{Line: params.Position.Line, Character: col}}, NewText: image.Src}})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Label < items[j].Label })
	return s.sendCompletion(msg.ID, params, &CompletionList{Items: items})
}
