package cmd

import (
	"strings"
	"testing"
)

func TestPatchSubscriptionBlockPreservesUnrelatedFields(t *testing.T) {
	original := "[[blogroll.feeds]]\nurl = \"https://example.test/feed\"\n" +
		"title = \"Old # title\" # keep this comment\n" +
		"tags = [\"go\", \"rss\"]\nmax_entries = 5\n"
	got, err := patchSubscriptionBlock(original, map[string]string{"title": `"New title"`}, []string{"title"})
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Replace(original, `"Old # title"`, `"New title"`, 1)
	if got != want {
		t.Fatalf("unexpected edit:\n%s\nwant:\n%s", got, want)
	}
}

func TestPatchSubscriptionBlockRejectsUnsafeMultiline(t *testing.T) {
	original := "[[blogroll.feeds]]\nurl = \"https://example.test/feed\"\ntags = [\n  \"one\",\n  \"two\",\n]\n"
	if _, err := patchSubscriptionBlock(original, map[string]string{"tags": `["new"]`}, []string{"tags"}); err == nil {
		t.Fatal("multiline field must not be partially rewritten")
	}
}
