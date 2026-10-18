package diagnostics

import (
	"strings"
	"testing"
)

func TestUnknownQuoteAttributions(t *testing.T) {
	resolver := &mockResolver{handles: map[string]bool{"known": true, "known.dev": true}}
	source := "---\ntitle: Test\n---\n~missing\n@unknown\n~known\n@known.dev.\ncontact@example.com\n~~not-a-mention~~\n```\n~ignored\n```\n"
	var got []Issue
	for _, issue := range Check("ping.md", source, resolver) {
		if issue.Code == "unknown-mention" {
			got = append(got, issue)
		}
	}
	if len(got) != 2 {
		t.Fatalf("want exactly two unknown mentions, got %v", got)
	}
	if !strings.Contains(got[0].Message, "~missing") || !strings.Contains(got[1].Message, "@unknown") {
		t.Errorf("missing marker in warnings: %+v", got)
	}
	for _, issue := range got {
		if !strings.Contains(issue.Message, "blogroll") {
			t.Errorf("must suggest blogroll fix: %s", issue.Message)
		}
	}
}
