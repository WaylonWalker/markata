package lsp

import (
	"os"
	"strings"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/prose"
)

// indexResolver adapts the Index to the diagnostics.Resolver interface.
type indexResolver struct {
	index *Index
}

func (r *indexResolver) ResolveSlug(slug string) bool {
	return r.index.GetBySlug(slug) != nil
}

func (r *indexResolver) ResolveHandle(handle string) bool {
	return r.index.GetByHandle(handle) != nil
}

// publishDiagnostics publishes diagnostics for a document.
func (s *Server) publishDiagnostics(uri, content string) error {
	diagnosticsList := s.computeDiagnostics(uri, content)

	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnosticsList,
	}

	s.docMu.RLock()
	if doc := s.documents[uri]; doc != nil {
		version := doc.Version
		params.Version = &version
	}
	s.docMu.RUnlock()
	return s.sendNotification("textDocument/publishDiagnostics", params)
}

// computeDiagnostics computes diagnostics for a document.
func (s *Server) computeDiagnostics(uri, content string) []Diagnostic {
	// Convert URI to file path for diagnostics
	filePath := uriToPath(uri)

	// Create resolver adapter for the index
	resolver := &indexResolver{index: s.index}

	// Use shared diagnostics package
	issues := diagnostics.Check(filePath, content, resolver)

	diagnosticsList := make([]Diagnostic, 0, len(issues))
	lines := strings.Split(content, "\n")

	for _, issue := range issues {
		diag := Diagnostic{
			Range: Range{
				Start: Position{Line: issue.Range.StartLine, Character: issue.Range.StartCol},
				End:   Position{Line: issue.Range.EndLine, Character: issue.Range.EndCol},
			},
			Severity: convertSeverity(issue.Severity),
			Source:   "markata-go",
			Message:  issue.Message,
			Code:     issue.Code,
		}
		if diag.Range.Start.Line >= 0 && diag.Range.Start.Line < len(lines) {
			diag.Range.Start.Character = utf16Column(lines[diag.Range.Start.Line], diag.Range.Start.Character)
		}
		if diag.Range.End.Line >= 0 && diag.Range.End.Line < len(lines) {
			diag.Range.End.Character = utf16Column(lines[diag.Range.End.Line], diag.Range.End.Character)
		}
		diagnosticsList = append(diagnosticsList, diag)
	}

	// Editorial prose observations use the same source positions as CLI
	// diagnostics. They are informational rather than build-blocking errors.
	if os.Getenv("MARKATA_PROSE_ENABLED") == "1" {
		preferences, err := prose.LoadPreferences(filePath)
		if err != nil {
			s.logger.Printf("Cannot load prose preferences: %v", err)
		}
		for _, issue := range prose.Check(content, preferences.Options()) {
			start, end := Position{Line: issue.Line, Character: issue.Column},
				Position{Line: issue.EndLine, Character: issue.EndColumn}
			if start.Line >= 0 && start.Line < len(lines) {
				start.Character = utf16Column(lines[start.Line], start.Character)
			}
			if end.Line >= 0 && end.Line < len(lines) {
				end.Character = utf16Column(lines[end.Line], end.Character)
			}
			diagnosticsList = append(diagnosticsList, Diagnostic{
				Range: Range{Start: start, End: end},
				Severity: DiagnosticSeverityHint,
				Source:   "markata-prose",
				Code:     issue.Rule,
				Message:  issue.Message,
			})
		}
	}
	return diagnosticsList
}

// convertSeverity converts diagnostics.Severity to LSP severity.
func convertSeverity(s diagnostics.Severity) int {
	switch s {
	case diagnostics.SeverityError:
		return DiagnosticSeverityError
	case diagnostics.SeverityWarning:
		return DiagnosticSeverityWarning
	case diagnostics.SeverityInfo:
		return DiagnosticSeverityInformation
	default:
		return DiagnosticSeverityWarning
	}
}

// clearDiagnostics clears diagnostics for a document.
func (s *Server) clearDiagnostics(uri string) error {
	params := PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []Diagnostic{},
	}

	return s.sendNotification("textDocument/publishDiagnostics", params)
}
