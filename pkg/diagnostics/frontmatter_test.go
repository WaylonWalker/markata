package diagnostics

import "testing"

func TestInspectFrontmatter_ValidBlock(t *testing.T) {
	inspection := InspectFrontmatter("---\ntitle: Test\n---\nBody")

	if !inspection.HasFrontmatter {
		t.Fatal("InspectFrontmatter() did not find frontmatter")
	}
	if inspection.Frontmatter != "title: Test" {
		t.Fatalf("Frontmatter = %q, want %q", inspection.Frontmatter, "title: Test")
	}
	if inspection.Body != "Body" {
		t.Fatalf("Body = %q, want %q", inspection.Body, "Body")
	}
	if inspection.OpeningLine != 1 || inspection.ClosingLine != 3 {
		t.Fatalf("delimiter lines = %d, %d; want 1, 3", inspection.OpeningLine, inspection.ClosingLine)
	}
}

func TestAnalyzeFrontmatter_DelimiterProblems(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantCode  string
		wantRange Range
	}{
		{
			name:     "suspicious opening delimiter",
			content:  "----\ntitle: Test\n---\nBody",
			wantCode: ReasonFrontmatterSuspiciousDelimiter,
			wantRange: Range{
				StartLine: 0,
				EndLine:   0,
				EndCol:    4,
			},
		},
		{
			name:     "leading whitespace",
			content:  " ---\ntitle: Test\n---\nBody",
			wantCode: ReasonFrontmatterLeadingWhitespace,
			wantRange: Range{
				StartLine: 0,
				EndLine:   0,
				EndCol:    4,
			},
		},
		{
			name:     "malformed closing delimiter",
			content:  "---\ntitle: Test\n--- \nBody",
			wantCode: ReasonFrontmatterMalformedClosing,
			wantRange: Range{
				StartLine: 2,
				EndLine:   2,
				EndCol:    4,
			},
		},
		{
			name:     "missing closing delimiter",
			content:  "---\ntitle: Test\nBody",
			wantCode: ReasonFrontmatterMissingClosing,
			wantRange: Range{
				StartLine: 0,
				EndLine:   0,
				EndCol:    3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeFrontmatter("post.md", tt.content)
			issue := issueWithCode(analysis.Issues, tt.wantCode)
			if issue == nil {
				t.Fatalf("AnalyzeFrontmatter() did not return %q: %#v", tt.wantCode, analysis.Issues)
			}
			if issue.Range != tt.wantRange {
				t.Errorf("issue range = %+v, want %+v", issue.Range, tt.wantRange)
			}
		})
	}
}

func TestInspectFrontmatter_TrailingOpeningWhitespaceIsNotLeading(t *testing.T) {
	inspection := InspectFrontmatter("--- \ntitle: Test\n---\nBody")
	if inspection.LeadingWhitespace {
		t.Fatal("trailing opening whitespace was reported as leading whitespace")
	}
	if inspection.HasFrontmatter {
		t.Fatal("non-exact opening delimiter was treated as frontmatter")
	}
}

func TestAnalyzeFrontmatter_InvalidYAML(t *testing.T) {
	analysis := AnalyzeFrontmatter("post.md", "---\ntitle: [\n---\nBody")

	issue := issueWithCode(analysis.Issues, ReasonFrontmatterParseError)
	if issue == nil {
		t.Fatalf("AnalyzeFrontmatter() did not report invalid YAML: %#v", analysis.Issues)
	}
	if analysis.Valid {
		t.Error("AnalyzeFrontmatter() marked invalid YAML as valid")
	}
	if issue.Range.StartLine != 1 {
		t.Errorf("parse issue line = %d, want 1", issue.Range.StartLine)
	}
}

func TestAnalyzeFrontmatter_MetadataTypes(t *testing.T) {
	analysis := AnalyzeFrontmatter("post.md", "---\ntitle: [Test]\npublished: 1\ntags: [ok, 2]\n---\nBody")

	for _, code := range []string{ReasonFrontmatterInvalidType} {
		if issueWithCode(analysis.Issues, code) == nil {
			t.Fatalf("AnalyzeFrontmatter() did not report %q: %#v", code, analysis.Issues)
		}
	}
	if countIssuesWithCode(analysis.Issues, ReasonFrontmatterInvalidType) != 3 {
		t.Errorf("invalid type issue count = %d, want 3", countIssuesWithCode(analysis.Issues, ReasonFrontmatterInvalidType))
	}
	issue := issueWithCode(analysis.Issues, ReasonFrontmatterInvalidType)
	if issue.Range.StartLine != 1 {
		t.Errorf("invalid type issue line = %d, want 1", issue.Range.StartLine)
	}
}

func TestAnalyzeFrontmatter_ScalarTagsWarning(t *testing.T) {
	tests := []struct {
		name, value, wantMessage string
		wantIssues               int
	}{
		{
			name: "scalar shorthand", value: "go", wantIssues: 1,
			wantMessage: `frontmatter field "tags" uses scalar shorthand; use a list of strings`,
		},
		{
			name: "numeric tags", value: "42", wantIssues: 1,
			wantMessage: `frontmatter field "tags" must be a list of strings`,
		},
		{name: "list", value: "[go, rust]", wantIssues: 0},
		{name: "empty", value: "null", wantIssues: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := AnalyzeFrontmatter("post.md", "---\ntags: "+tt.value+"\n---\nBody")
			if !analysis.Valid {
				t.Fatal("valid YAML frontmatter rejected")
			}
			if got := countIssuesWithCode(analysis.Issues, ReasonFrontmatterInvalidType); got != tt.wantIssues {
				t.Fatalf("invalid type issues = %d, want %d: %#v", got, tt.wantIssues, analysis.Issues)
			}
			if tt.wantIssues == 0 {
				return
			}
			issue := issueWithCode(analysis.Issues, ReasonFrontmatterInvalidType)
			if issue.Message != tt.wantMessage || issue.Severity != SeverityWarning || issue.Range.StartLine != 1 {
				t.Errorf("issue = %+v, want message %q, warning on second line", *issue, tt.wantMessage)
			}
		})
	}
}

func TestAnalyzeFrontmatter_InvalidDate(t *testing.T) {
	analysis := AnalyzeFrontmatter("post.md", "---\ndate: 2020-1-15T00:00:00\n---\nBody")

	issue := issueWithCode(analysis.Issues, "invalid-date")
	if issue == nil {
		t.Fatalf("AnalyzeFrontmatter() did not report invalid date: %#v", analysis.Issues)
	}
	if ReasonForDiagnosticCode(issue.Code) != ReasonFrontmatterInvalidDate {
		t.Errorf("ReasonForDiagnosticCode(%q) = %q, want %q", issue.Code, ReasonForDiagnosticCode(issue.Code), ReasonFrontmatterInvalidDate)
	}
}

func issueWithCode(issues []Issue, code string) *Issue {
	for index := range issues {
		if issues[index].Code == code {
			return &issues[index]
		}
	}
	return nil
}

func countIssuesWithCode(issues []Issue, code string) int {
	count := 0
	for _, issue := range issues {
		if issue.Code == code {
			count++
		}
	}
	return count
}
