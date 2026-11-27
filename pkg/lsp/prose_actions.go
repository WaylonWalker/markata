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

// The ignore commands are available to native editors only. Browser requests
// must never gain filesystem-write access through the LSP bridge.
func (s *Server) handleProseExecuteCommand(_ context.Context, msg *Message) error {
	if s.browser != nil || os.Getenv("MARKATA_PROSE_ENABLED") != "1" {
		return s.sendError(msg.ID, MethodNotFound, "prose commands unavailable")
	}
	var params proseExecuteParams
	if err := json.Unmarshal(msg.Params, &params); err != nil {
		return s.sendError(msg.ID, InvalidParams, "invalid prose command")
	}
	target, err := parseProseIgnoreTarget(params)
	if err != nil {
		return s.sendError(msg.ID, InvalidParams, err.Error())
	}
	content, issue, err := s.findProseIssue(target)
	if err != nil {
		return s.sendError(msg.ID, InvalidParams, err.Error())
	}
	if err := s.applyProseIgnore(params.Command, target, issue, content); err != nil {
		return s.sendError(msg.ID, InternalError, err.Error())
	}
	if err := s.sendResponse(msg.ID, nil); err != nil {
		return err
	}
	return s.publishDiagnostics(target.uri, content)
}

type proseIgnoreTarget struct {
	uri    string
	id     string
	line   int
	column int
}

func parseProseIgnoreTarget(params proseExecuteParams) (proseIgnoreTarget, error) {
	var target proseIgnoreTarget
	if len(params.Arguments) != 4 ||
		(params.Command != "markata.prose.ignoreOnce" && params.Command != "markata.prose.ignoreRule") {
		return target, fmt.Errorf("unsupported prose ignore command")
	}
	if json.Unmarshal(params.Arguments[0], &target.uri) != nil ||
		json.Unmarshal(params.Arguments[1], &target.id) != nil ||
		json.Unmarshal(params.Arguments[2], &target.line) != nil ||
		json.Unmarshal(params.Arguments[3], &target.column) != nil ||
		target.uri == "" || !prose.KnownRule(target.id) {
		return target, fmt.Errorf("invalid prose rule or document")
	}
	return target, nil
}

func (s *Server) findProseIssue(target proseIgnoreTarget) (content string, issue prose.Issue, err error) {
	path := uriToPath(target.uri)
	if s.rootURI != "" {
		root := uriToPath(s.rootURI)
		rel, pathErr := filepath.Rel(root, path)
		if pathErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return "", prose.Issue{}, fmt.Errorf("document is outside the LSP workspace")
		}
	}
	s.docMu.RLock()
	doc := s.documents[target.uri]
	if doc != nil {
		content = doc.Content
	}
	s.docMu.RUnlock()
	if doc == nil {
		return "", prose.Issue{}, fmt.Errorf("document is not open")
	}
	prefs, err := prose.LoadPreferences(path)
	if err != nil {
		return "", prose.Issue{}, err
	}
	for _, candidate := range prose.Check(content, prefs.Options()) {
		if candidate.Rule == target.id && candidate.Line == target.line && candidate.Column == target.column {
			return content, candidate, nil
		}
	}
	return "", prose.Issue{}, fmt.Errorf("stale prose suggestion")
}

func (s *Server) applyProseIgnore(command string, target proseIgnoreTarget, issue prose.Issue, content string) error {
	if command == "markata.prose.ignoreRule" {
		_, err := prose.IgnoreRule(uriToPath(target.uri), target.id)
		return err
	}
	key := proseIssueKey(issue, content)
	s.docMu.Lock()
	defer s.docMu.Unlock()
	if s.proseIgnored == nil {
		s.proseIgnored = make(map[string]map[string]bool)
	}
	if s.proseIgnored[target.uri] == nil {
		s.proseIgnored[target.uri] = make(map[string]bool)
	}
	s.proseIgnored[target.uri][key] = true
	return nil
}
