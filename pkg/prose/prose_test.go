package prose

import "testing"

func TestCheckProseAndProtectedMarkdown(t *testing.T) {
	source := "---\ntitle: delve into\n---\nHere we delve into the API.\n\n```go\n// delve into\n```\n\nUse `delve into` or [delve into](https://example.com/delve-into).\n"
	issues := Check(source, Options{})
	if len(issues) != 2 {
		t.Fatalf("expected prose in the paragraph and link label, got %+v", issues)
	}
	if issues[0].Rule != "prose/cliche" || issues[0].Line != 3 || issues[0].Column != 8 {
		t.Fatalf("unexpected issue: %+v", issues[0])
	}
}

func TestIgnoreRule(t *testing.T) {
	source := "It is important to note that this is very unique."
	all := Check(source, Options{})
	if len(all) != 2 {
		t.Fatalf("expected two issues, got %+v", all)
	}
	filtered := Check(source, Options{IgnoredRules: map[string]bool{"prose/wordiness": true}})
	if len(filtered) != 1 || filtered[0].Rule != "prose/intensifier" {
		t.Fatalf("unexpected filtered issues: %+v", filtered)
	}
}

func TestUTF8ByteOffset(t *testing.T) {
	source := "☀️ We delve into testing."
	issues := Check(source, Options{})
	if len(issues) != 1 || issues[0].Column != len("☀️ We ") {
		t.Fatalf("expected byte column %d, got %+v", len("☀️ We "), issues)
	}
}

func TestNoFalseAIAuthorshipClaim(t *testing.T) {
	for _, issue := range Check("A rich tapestry.", Options{}) {
		if issue.Rule == "" || issue.Message == "" {
			t.Fatalf("missing details: %+v", issue)
		}
	}
}
