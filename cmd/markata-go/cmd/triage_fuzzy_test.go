package cmd

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestTriageMarkdownCandidatesFindsPostsWithoutReadingFiles(t *testing.T) {
	root := t.TempDir()
	write := func(path string) {
		t.Helper()
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("private body not included in triage"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{
		"posts/alpha.md", "posts/nested/beta.MARKDOWN", "pages/about.md",
		"docs/tutorial.md", "content/new-post.md", "src/posts/extra.md",
		"posts/.hidden.md", "posts/.drafts/private.md",
		"posts/node_modules/dependency.md", "posts/build/generated.md",
		"static/asset.md", "other/custom.md",
	} {
		write(path)
	}
	// Symlinks are never traversed or offered as Markdown candidates.
	if err := os.Symlink(filepath.Join(root, "posts", "alpha.md"), filepath.Join(root, "posts", "linked.md")); err != nil {
		t.Logf("symlink unsupported here: %v", err)
	}
	want := []string{
		"content/new-post.md", "docs/tutorial.md", "pages/about.md",
		"posts/alpha.md", "posts/nested/beta.MARKDOWN", "src/posts/extra.md",
	}
	got := triageMarkdownCandidates(root)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected post choices: got %q, want %q", got, want)
	}
	if paths := triageMarkdownCandidates(filepath.Join(root, "missing")); len(paths) != 0 {
		t.Fatalf("missing site directory should return no choices: %v", paths)
	}
}

func TestTriagePostPickerAllowsManualEntry(t *testing.T) {
	resetTriageIntakeFlags(t)
	beforeSiteDir := siteDir
	t.Cleanup(func() { siteDir = beforeSiteDir })
	siteDir = t.TempDir()
	post := filepath.Join(siteDir, "posts", "known.md")
	if err := os.MkdirAll(filepath.Dir(post), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(post, []byte("# Known"), 0o600); err != nil {
		t.Fatal(err)
	}
	triageAction = "review"
	var prompt strings.Builder
	cmd := &cobra.Command{Use: "post"}
	cmd.SetIn(strings.NewReader("custom/my-post.md\n"))
	cmd.SetErr(&prompt)
	picked := false
	req, _, err := resolveTriageIntakeWithPicker(cmd, true, func(q triageQuestion, options []string) (string, error) {
		if q.key != "post" || len(options) < 2 ||
			options[0] != "posts/known.md" || options[len(options)-1] != triageManualPostOption {
			t.Fatalf("bad post choices for %s: %v", q.key, options)
		}
		picked = true
		return triageManualPostOption, nil
	})
	if err != nil || !picked || req.Post != "custom/my-post.md" {
		t.Fatalf("manual path failed: %#v; %v; picker called=%t", req, err, picked)
	}
	if !strings.Contains(prompt.String(), "Which Markdown post") {
		t.Fatalf("post prompt absent: %q", prompt.String())
	}
}

func TestTriagePostPickerFallsBackToManualWhenNoFilesFound(t *testing.T) {
	question := triageQuestions["post"]
	var shown strings.Builder
	answer, err := askTriageQuestion(bufio.NewReader(strings.NewReader("elsewhere/note.md\n")), &shown,
		question, func(_ triageQuestion, _ []string) (string, error) {
			t.Fatal("no picker expected for empty directory")
			return "", nil
		}, t.TempDir())
	if err != nil || answer != "elsewhere/note.md" {
		t.Fatalf("manual fallback: %q, %v", answer, err)
	}
}

func TestTriageFuzzyPickerCancellationDoesNotContinue(t *testing.T) {
	resetTriageIntakeFlags(t)
	var prompt strings.Builder
	cmd := &cobra.Command{Use: "post"}
	cmd.SetErr(&prompt)
	picked := 0
	_, _, err := resolveTriageIntakeWithPicker(cmd, true, func(_ triageQuestion, _ []string) (string, error) {
		picked++
		return "", errTriageIntakeCanceled
	})
	if !errors.Is(err, errTriageIntakeCanceled) || picked != 1 {
		t.Fatalf("cancellation should stop at first choice: picks=%d err=%v", picked, err)
	}
}

func TestTriageFuzzyPickerRejectsUnknownChoice(t *testing.T) {
	resetTriageIntakeFlags(t)
	cmd := &cobra.Command{Use: "post"}
	_, _, err := resolveTriageIntakeWithPicker(cmd, true, func(_ triageQuestion, _ []string) (string, error) {
		return "invented", nil
	})
	if err == nil || !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("unknown fuzzy choice should be rejected: %v", err)
	}
}

func TestTriagePostFlagSkipsFileDiscoveryAndPicker(t *testing.T) {
	resetTriageIntakeFlags(t)
	triagePost = "custom/article.md"
	triageAction = "review"
	cmd := &cobra.Command{Use: "post"}
	req, _, err := resolveTriageIntakeWithPicker(cmd, true, func(q triageQuestion, _ []string) (string, error) {
		t.Fatalf("no picker expected when all flags specified; got %s", q.key)
		return "", nil
	})
	if err != nil || req.Post != "custom/article.md" {
		t.Fatalf("flag-supplied path should be used as-is: %#v %v", req, err)
	}
}
