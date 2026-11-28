package prose

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode/utf8"
)

func harperCLIJSONFixture(t *testing.T, source, matched string, suggestions []string) []byte {
	t.Helper()
	start := strings.Index(source, matched)
	if start < 0 {
		t.Fatal("missing fixture match")
	}
	lint := harperCLILint{
		Rule:        "Spelling",
		Kind:        "Spelling",
		Message:     "Check this spelling.",
		MatchedText: matched,
		Suggestions: suggestions,
	}
	lint.Span.CharStart = utf8.RuneCountInString(source[:start])
	lint.Span.CharEnd = lint.Span.CharStart + utf8.RuneCountInString(matched)
	data, err := json.Marshal([]harperCLIFileResult{{
		File: "draft.md", Lints: []harperCLILint{lint},
	}})
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestHarperJSONConvertsUnicodeCharacterSpansToBytes(t *testing.T) {
	source := "🌱 I have a apple."
	data := harperCLIJSONFixture(t, source, "a apple", []string{"Replace with: “an apple”"})
	issues, err := ParseHarperCLIJSON(source, data, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected one grammar issue, got %+v", issues)
	}
	issue := issues[0]
	expectedStart := strings.Index(source, "a apple")
	if issue.StartByte != expectedStart || issue.EndByte != expectedStart+len("a apple") ||
		issue.Provider != "harper" || issue.Rule != "harper/Spelling" {
		t.Fatalf("Harper offsets or rule metadata were changed: %+v", issue)
	}
	if issue.Suggestion != "an apple" || len(issue.Edits) != 1 ||
		issue.Edits[0].Replacement != "an apple" {
		t.Fatalf("expected safe replacement, got %+v", issue)
	}
	span, err := UTF16Range(source, issue)
	if err != nil {
		t.Fatal(err)
	}
	if span.Start.Character != utf8.RuneCountInString(source[:expectedStart])+1 {
		t.Fatalf("invalid UTF-16 column after emoji: %+v", span)
	}
}

func TestHarperJSONFiltersSharedPreferences(t *testing.T) {
	source := "🌱 I have a apple."
	data := harperCLIJSONFixture(t, source, "a apple", nil)
	ignored := Options{IgnoredPhrases: map[string]map[string]bool{
		"harper/Spelling": {"a apple": true},
	}}
	issues, err := ParseHarperCLIJSON(source, data, ignored)
	if err != nil || len(issues) != 0 {
		t.Fatalf("phrase should be ignored, got %+v, %v", issues, err)
	}
}

func TestHarperJSONNeverMistakesHumanActionForReplacement(t *testing.T) {
	source := "Fix this."
	for _, suggestion := range []string{"Insert “,”", "Remove error", "Replace with: incomplete"} {
		data := harperCLIJSONFixture(t, source, "Fix", []string{suggestion})
		issues, err := ParseHarperCLIJSON(source, data, Options{})
		if err != nil || len(issues) != 1 {
			t.Fatalf("missing harmless observation: %+v, %v", issues, err)
		}
		if len(issues[0].Edits) != 0 || issues[0].Suggestion != "" {
			t.Fatalf("misparsed human-readable action %q: %+v", suggestion, issues[0])
		}
	}
}

func TestHarperJSONRejectsMismatchedOrCorruptSourceRanges(t *testing.T) {
	source := "🌱 Some words."
	data := harperCLIJSONFixture(t, source, "Some", nil)
	var files []harperCLIFileResult
	if err := json.Unmarshal(data, &files); err != nil {
		t.Fatal(err)
	}
	files[0].Lints[0].MatchedText = "other"
	broken, marshalErr := json.Marshal(files)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	if _, err := ParseHarperCLIJSON(source, broken, Options{}); err == nil {
		t.Fatal("accepted a suggestion pointing to unrelated source text")
	}
	files[0].Lints[0].MatchedText = "Some"
	files[0].Lints[0].Span.CharEnd = 900
	broken, marshalErr = json.Marshal(files)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	if _, err := ParseHarperCLIJSON(source, broken, Options{}); err == nil {
		t.Fatal("accepted out-of-bounds character offsets")
	}
	if _, err := ParseHarperCLIJSON(source, []byte("not json"), Options{}); err == nil {
		t.Fatal("accepted malformed JSON")
	}
	if _, err := ParseHarperCLIJSON(string([]byte{0xff}), data, Options{}); err == nil {
		t.Fatal("accepted invalid UTF-8 source")
	}
}

func TestHarperCLIBinaryMissingDoesNotBreakDefaultChecker(t *testing.T) {
	if len(Check("This is very unique.", Options{})) != 1 {
		t.Fatal("default checker must work independently of Harper")
	}
	_, err := CheckWithHarperCLI(context.Background(), "This is an example.", HarperCLIOptions{
		Binary: "markata-harper-test-command-does-not-exist",
	})
	if err == nil {
		t.Fatal("explicitly requested missing Harper binary should report an error")
	}
}

func TestHarperCLIExpectedFindingExitStatusIsNotFatal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is validated on Linux and macOS")
	}
	source := "🌱 This is very unique."
	jsonResult := harperCLIJSONFixture(t, source, "very unique", []string{"Replace with: “unique”"})
	binary := filepath.Join(t.TempDir(), "harper-cli-fixture")
	// Harper intentionally exits with code 1 when it reports findings.
	script := "#!/bin/sh\nprintf '%s' '" + string(jsonResult) + "'\nexit 1\n"
	if err := os.WriteFile(binary, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	// #nosec G302 -- a private test fixture requires executable permissions.
	if err := os.Chmod(binary, 0o700); err != nil {
		t.Fatal(err)
	}
	issues, err := CheckWithHarperCLI(context.Background(), source, HarperCLIOptions{Binary: binary})
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Provider != "harper" || issues[0].Suggestion != "unique" {
		t.Fatalf("expected Harper finding despite exit 1, got %+v", issues)
	}
}

func TestHarperCLINonzeroExitWithoutLintsIsError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell fixture is validated on Linux and macOS")
	}
	binary := filepath.Join(t.TempDir(), "harper-cli-fixture")
	script := "#!/bin/sh\nprintf '%s' '[{\"file\":\"draft.md\",\"lints\":[]}]'\nexit 1\n"
	if err := os.WriteFile(binary, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	// #nosec G302 -- a private test fixture requires executable permissions.
	if err := os.Chmod(binary, 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := CheckWithHarperCLI(context.Background(), "This is fine.", HarperCLIOptions{Binary: binary}); err == nil {
		t.Fatal("failure exit with no lints cannot be accepted as a successful check")
	}
}
