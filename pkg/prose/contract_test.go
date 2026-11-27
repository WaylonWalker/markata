package prose

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSourceContractUnicodeAndCRLF(t *testing.T) {
	source := "Heading\r\n🌱 We very unique.\r\n"
	issues := Check(source, Options{})
	if len(issues) != 1 {
		t.Fatalf("expected one issue: %+v", issues)
	}
	issue := issues[0]
	if issue.ContractVersion != ContractVersion || issue.Provider != "markata" || issue.Severity != "hint" {
		t.Fatalf("missing metadata: %+v", issue)
	}
	if issue.Line != 1 || issue.Column != len("🌱 We ") ||
		issue.StartByte != strings.Index(source, "very unique") ||
		issue.EndByte != issue.StartByte+len("very unique") {
		t.Fatalf("bad byte offsets: %+v", issue)
	}
	text, err := OriginalText(source, issue)
	if err != nil || text != "very unique" {
		t.Fatalf("unexpected source match: %q, %v", text, err)
	}
	span, err := UTF16Range(source, issue)
	if err != nil {
		t.Fatal(err)
	}
	if span.Start.Line != 1 || span.Start.Character != 6 ||
		span.End.Line != 1 || span.End.Character != 17 {
		t.Fatalf("incorrect LSP coordinates: %+v", span)
	}
	if len(issue.Edits) != 1 || issue.Edits[0].Replacement != "unique" ||
		issue.Edits[0].StartByte != issue.StartByte || issue.Edits[0].EndByte != issue.EndByte {
		t.Fatalf("unsafe replacement metadata: %+v", issue)
	}
	jsonResult, err := json.Marshal(issues)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"contract_version", "start_byte", "end_byte", "edits", "provider"} {
		if !strings.Contains(string(jsonResult), field) {
			t.Fatalf("missing %s from JSON contract: %s", field, jsonResult)
		}
	}
}

func TestIssueFactoryRejectsInvalidRanges(t *testing.T) {
	src := "🌱 text"
	for _, test := range []struct{ start, end int }{
		{-1, 1}, {2, 6}, {0, 200}, {6, 2},
	} {
		if _, err := NewIssue(src, "harper", "spelling", "Spelling", test.start, test.end, ""); err == nil {
			t.Fatalf("unsafe range accepted: %+v", test)
		}
	}
	issue, err := NewIssue(src, "harper", "spelling", "Spelling", len("🌱 "), len(src), "words")
	if err != nil {
		t.Fatal(err)
	}
	if issue.Line != 0 || issue.Column != len("🌱 ") || len(issue.Edits) != 1 {
		t.Fatalf("unexpected provider issue: %+v", issue)
	}
}

func TestProtectedMarkdownPreservesSourceSpans(t *testing.T) {
	source := "---\ntitle: very unique\n---\n\n```go\n// very unique\n```\n\n[very unique](https://example.com/very-unique) then very unique.\n"
	issues := Check(source, Options{})
	if len(issues) != 2 {
		t.Fatalf("only visible link text and prose should be linted: %+v", issues)
	}
	for _, issue := range issues {
		phrase, err := OriginalText(source, issue)
		if err != nil || phrase != "very unique" {
			t.Fatalf("range changed protected syntax: %+v, %v", issue, err)
		}
	}
}
