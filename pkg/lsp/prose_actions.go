package lsp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/prose"
)

// Prose code actions intentionally use standard LSP edits and commands so
// Neovim needs no Markata-specific UI extension.
type proseCodeActionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

type proseCodeAction struct {
	Title       string              `json:"title"`
	Kind        string              `json:"kind"`
	IsPreferred bool                `json:"isPreferred,omitempty"`
	Edit        *proseWorkspaceEdit `json:"edit,omitempty"`
	Command     *proseCommand       `json:"command,omitempty"`
}

type proseWorkspaceEdit struct {
	DocumentChanges []proseDocumentEdit `json:"documentChanges"`
}

type proseDocumentEdit struct {
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`
	Edits        []TextEdit                      `json:"edits"`
}

type proseCommand struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments"`
}

type proseExecuteParams struct {
	Command   string            `json:"command"`
	Arguments []json.RawMessage `json:"arguments"`
}

func (s *Server) handleProseCodeAction(_ context.Context, msg *Message) error {
	actions := make([]proseCodeAction, 0)
	if os.Getenv("MARKATA_PROSE_ENABLED") != "1" {
		return s.sendResponse(msg.ID, actions)
	}
	var params proseCodeActionParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		return s.sendError(msg.ID, InvalidParams, "invalid code action params")
	}
	s.docMu.RLock()
	doc := s.documents[params.TextDocument.URI]
	var content string
	var version int
	if doc != nil {
		content, version = doc.Content, doc.Version
	}
	s.docMu.RUnlock()
	if doc == nil {
		return s.sendResponse(msg.ID, actions)
	}
	prefs, err := prose.LoadPreferences(uriToPath(params.TextDocument.URI))
	if err != nil {
		s.logger.Printf("Cannot load prose preferences: %v", err)
		return s.sendResponse(msg.ID, actions)
	}
	for _, issue := range prose.Check(content, prefs.Options()) {
		issueRange := proseLSPRange(content, issue)
		if !proseRangeIntersects(params.Range, issueRange) || s.proseIssueIgnored(params.TextDocument.URI, issue, content) {
			continue
		}
		if issue.Suggestion != "" {
			actions = append(actions, proseCodeAction{
				Title:       "Replace with " + issue.Suggestion,
				Kind:        "quickfix",
				IsPreferred: false,
				Edit: &proseWorkspaceEdit{DocumentChanges: []proseDocumentEdit{{
					TextDocument: VersionedTextDocumentIdentifier{
						TextDocumentIdentifier: TextDocumentIdentifier{URI: params.TextDocument.URI},
						Version:                version,
					},
					Edits: []TextEdit{{Range: issueRange, NewText: issue.Suggestion}},
				}}},
			})
		}
		// Browser sessions deliberately cannot mutate server-side workspace
		// preferences. Their editor UI will own ignore persistence.
		if s.browser != nil {
			continue
		}
		args := []interface{}{params.TextDocument.URI, issue.Rule, issue.Line, issue.Column}
		actions = append(actions, proseCodeAction{
			Title: "Ignore this occurrence",
			Kind:  "quickfix",
			Command: &proseCommand{
				Title:     "Ignore this occurrence",
				Command:   "markata.prose.ignoreOnce",
				Arguments: args,
			},
		})
		actions = append(actions, proseCodeAction{
			Title: "Always ignore " + issue.Rule,
			Kind:  "quickfix",
			Command: &proseCommand{
				Title:     "Always ignore " + issue.Rule,
				Command:   "markata.prose.ignoreRule",
				Arguments: args,
			},
		})
	}
	return s.sendResponse(msg.ID, actions)
}

func proseLSPRange(content string, issue prose.Issue) Range {
	lines := strings.Split(content, "\n")
	start := Position{Line: issue.Line, Character: issue.Column}
	end := Position{Line: issue.EndLine, Character: issue.EndColumn}
	if start.Line >= 0 && start.Line < len(lines) {
		start.Character = utf16Column(lines[start.Line], start.Character)
	}
	if end.Line >= 0 && end.Line < len(lines) {
		end.Character = utf16Column(lines[end.Line], end.Character)
	}
	return Range{Start: start, End: end}
}

func proseComparePosition(a, b Position) int {
	switch {
	case a.Line < b.Line:
		return -1
	case a.Line > b.Line:
		return 1
	case a.Character < b.Character:
		return -1
	case a.Character > b.Character:
		return 1
	default:
		return 0
	}
}

func proseRangeIntersects(a, b Range) bool {
	return proseComparePosition(a.Start, b.End) <= 0 && proseComparePosition(b.Start, a.End) <= 0
}

func proseIssueKey(issue prose.Issue, content string) string {
	lines := strings.Split(content, "\n")
	if issue.Line < 0 || issue.Line >= len(lines) || issue.EndLine != issue.Line {
		return ""
	}
	line := lines[issue.Line]
	if issue.Column < 0 || issue.EndColumn > len(line) || issue.Column >= issue.EndColumn {
		return ""
	}
	return fmt.Sprintf("%s:%d:%d:%s", issue.Rule, issue.Line, issue.Column, line[issue.Column:issue.EndColumn])
}

func (s *Server) proseIssueIgnored(uri string, issue prose.Issue, content string) bool {
	key := proseIssueKey(issue, content)
	s.docMu.RLock()
	defer s.docMu.RUnlock()
	return key != "" && s.proseIgnored[uri][key]
}

// The ignore commands are available to local LSP sessions, not the browser
// bridge, which must not provide arbitrary workspace write capabilities.
func (s *Server) handleProseExecuteCommand(_ context.Context, msg *Message) error {
	if s.browser != nil || os.Getenv("MARKATA_PROSE_ENABLED") != "1" {
		return s.sendError(msg.ID, MethodNotFound, "prose commands unavailable")
	}
	var params proseExecuteParams
	if err := json.Unmarshal(msg.Params, &params); err != nil || len(params.Arguments) != 4 {
		return s.sendError(msg.ID, InvalidParams, "invalid prose command arguments")
	}
	if params.Command != "markata.prose.ignoreOnce" && params.Command != "markata.prose.ignoreRule" {
		return s.sendError(msg.ID, MethodNotFound, "unknown prose command")
	}
	var uri, id string
	var line, column int
	if json.Unmarshal(params.Arguments[0], &uri) != nil ||
		json.Unmarshal(params.Arguments[1], &id) != nil ||
		json.Unmarshal(params.Arguments[2], &line) != nil ||
		json.Unmarshal(params.Arguments[3], &column) != nil ||
		uri == "" || !prose.KnownRule(id) {
		return s.sendError(msg.ID, InvalidParams, "invalid prose rule or document")
	}
	path := uriToPath(uri)
	if s.rootURI != "" {
		root := uriToPath(s.rootURI)
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return s.sendError(msg.ID, InvalidParams, "document is outside the LSP workspace")
		}
	}
	s.docMu.RLock()
	doc := s.documents[uri]
	var content string
	if doc != nil {
		content = doc.Content
	}
	s.docMu.RUnlock()
	if doc == nil {
		return s.sendError(msg.ID, InvalidParams, "document is not open")
	}
	prefs, err := prose.LoadPreferences(path)
	if err != nil {
		return s.sendError(msg.ID, InvalidParams, "invalid prose preferences")
	}
	var selected *prose.Issue
	for _, issue := range prose.Check(content, prefs.Options()) {
		if issue.Rule == id && issue.Line == line && issue.Column == column {
			copyIssue := issue
			selected = &copyIssue
			break
		}
	}
	if selected == nil {
		return s.sendError(msg.ID, InvalidParams, "stale prose suggestion")
	}
	if params.Command == "markata.prose.ignoreRule" {
		if _, err := prose.IgnoreRule(path, id); err != nil {
			return s.sendError(msg.ID, InternalError, "unable to save prose preferences")
		}
	} else {
		key := proseIssueKey(*selected, content)
		s.docMu.Lock()
		if s.proseIgnored[uri] == nil {
			s.proseIgnored[uri] = make(map[string]bool)
		}
		s.proseIgnored[uri][key] = true
		s.docMu.Unlock()
	}
	if err := s.sendResponse(msg.ID, nil); err != nil {
		return err
	}
	return s.publishDiagnostics(uri, content)
}
