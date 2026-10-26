package lsp

import (
	"github.com/WaylonWalker/markata/pkg/diagnostics"
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
		diagnosticsList = append(diagnosticsList, diag)
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
