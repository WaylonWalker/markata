package lsp

import (
	"testing"
)

func TestContainerCompletion_MatchesCommandPaletteFamilies(t *testing.T) {
	params := CompletionParams{Position: Position{Line: 0}}
	cases := []struct{ line, want string }{
		{"::::: md-i", "md-icon"},
		{"::::: md-g", "md-grid"},
		{"::: wa-c", "wa-carousel"},
		{"::: road e", "editorial"},
		{":::: md-grid ", "rail"},
		{":::: md-grid r", "rail"},
	}
	for _, tt := range cases {
		t.Run(tt.line, func(t *testing.T) {
			items, ok := containerCompletions(tt.line, len(tt.line), params)
			if !ok {
				t.Fatal("container context not detected")
			}
			for _, item := range items {
				if item.Label == tt.want {
					return
				}
			}
			t.Fatalf("missing %q from %+v", tt.want, items)
		})
	}
}

func TestContainerAttributeCompletion_SuggestsNamesAndValues(t *testing.T) {
	params := CompletionParams{Position: Position{Line: 0}}
	cases := []struct{ line, want string }{
		{"::::: md-icon {data-s", "data-size"},
		{"::: wa-details {sum", "summary"},
		{"::: wa-carousel {navi", "navigation"},
		{"::: wa-badge {var", "variant"},
		{"::: md-card {.", ".class"},
		{"::: md-card {#", "#id"},
		{"::::: md-grid {data-columns=\"", "3"},
		{"::::: md-icon {data-size=\"", "lg"},
		{"::::: md-icon {data-size=\"l", "lg"},
		{"::: wa-badge {variant=\"b", "brand"},
		{"::: wa-carousel {navigation=\"t", "true"},
		{"::: wa-carousel {navigation=\"true\" data-s", "data-size"},
	}
	for _, tt := range cases {
		t.Run(tt.line, func(t *testing.T) {
			items, ok := attributeCompletions(tt.line, len(tt.line), params)
			if !ok {
				t.Fatal("attribute context not detected")
			}
			for _, item := range items {
				if item.Label == tt.want {
					return
				}
			}
			t.Fatalf("missing %q from %+v", tt.want, items)
		})
	}
}

func TestMermaidCompletion_SuggestsDiagramKindsDirectionsAndStyles(t *testing.T) {
	cases := []struct {
		name   string
		lines  []string
		line   int
		cursor int
		want   string
	}{
		{"diagram kind", []string{"```mermaid", "g"}, 1, 1, "graph"},
		{"flowchart starter", []string{"```mermaid", ""}, 1, 0, "flowchart TD"},
		{"direction", []string{"```mermaid", "flowchart T"}, 1, 11, "TD"},
		{"class style", []string{"```mermaid", "classDef note c"}, 1, 15, "color"},
		{"node shape", []string{"```mermaid", "flowchart TD", "  A("}, 2, 4, ")"},
		{"outside code", []string{"g"}, 0, 1, ""},
		{"after Mermaid fence", []string{"```mermaid", "flowchart TD", "```", "g"}, 3, 1, ""},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			params := CompletionParams{Position: Position{Line: tt.line}}
			items, ok := mermaidCompletions(tt.lines, tt.line, tt.lines[tt.line], tt.cursor, params)
			if tt.want == "" {
				if ok {
					t.Fatal("recognized non-Mermaid context")
				}
				return
			}
			if !ok {
				t.Fatal("Mermaid context not detected")
			}
			for _, item := range items {
				if item.Label == tt.want {
					return
				}
			}
			t.Fatalf("missing %q from %+v", tt.want, items)
		})
	}
}
