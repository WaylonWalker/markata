package editorial

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func safetyWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func safetyRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func safetyGit(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v: %s", args, err, out)
	}
	return string(out)
}

func safetyRepo(t *testing.T, post string) string {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	root := t.TempDir()
	safetyGit(t, root, "init", "-q")
	safetyGit(t, root, "config", "user.name", "Editorial test")
	safetyGit(t, root, "config", "user.email", "editorial@example.invalid")
	safetyGit(t, root, "remote", "add", "origin", "https://github.com/example/site.git")
	safetyWrite(t, filepath.Join(root, "posts", "a post's.md"), post)
	safetyWrite(t, filepath.Join(root, "other.txt"), "original\n")
	safetyGit(t, root, "add", ".")
	safetyGit(t, root, "commit", "-qm", "baseline")
	return root
}

func safetyRecord(t *testing.T, text string) map[string]interface{} {
	t.Helper()
	_, node, err := split(text)
	if err != nil {
		t.Fatal(err)
	}
	var history []map[string]interface{}
	if err := mapping(&node, "editorial_history").Decode(&history); err != nil {
		t.Fatal(err)
	}
	return history[len(history)-1]
}

func TestEditorialSafety_IsolatedIndexAndEncodedPaths(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthors:\n  - writer\n---\nBody\n")
	post := "posts/a post's.md"
	safetyWrite(t, filepath.Join(root, "other.txt"), "staged user change\n")
	safetyGit(t, root, "add", "other.txt")
	staged := safetyGit(t, root, "show", ":other.txt")
	safetyWrite(t, filepath.Join(root, "other.txt"), "unstaged user change\n")
	safetyWrite(t, filepath.Join(root, post), "---\ntitle: Post\nauthors:\n  - writer\n---\nChanged by the editor\n")
	session, err := Commit(Options{Root: root}, "human", "copy-editor", "Corrected wording", []string{post})
	if err != nil {
		t.Fatal(err)
	}
	if len(session.Commits) != 2 {
		t.Fatalf("commits = %v", session.Commits)
	}
	for _, sha := range session.Commits {
		files := safetyGit(t, root, "diff-tree", "--no-commit-id", "--name-only", "-z", "-r", sha)
		if files != post+"\x00" {
			t.Fatalf("unexpected committed paths: %q", files)
		}
	}
	if safetyGit(t, root, "show", ":other.txt") != staged || safetyRead(t, filepath.Join(root, "other.txt")) != "unstaged user change\n" || safetyGit(t, root, "show", "HEAD:other.txt") != "original\n" {
		t.Fatal("unrelated staged or unstaged work was altered")
	}
	record := safetyRecord(t, safetyRead(t, filepath.Join(root, post)))
	fileURL, fileOK := record["file_url"].(string)
	if !fileOK {
		t.Fatalf("missing file URL: %#v", record)
	}
	if record["commit"] != session.Commits[0] || record["action"] != "edited" || !strings.Contains(fileURL, "a%20post%27s.md") {
		t.Fatalf("wrong attribution: %#v", record)
	}
}

func TestEditorialSafety_DryRunLeavesFilesHistoryAndSessionUntouched(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	post := "posts/a post's.md"
	path := filepath.Join(root, post)
	s, err := Start(Options{Root: root}, "human", "copy-editor", []string{post})
	if err != nil {
		t.Fatal(err)
	}
	safetyWrite(t, path, "---\ntitle: Post\nauthor: writer\n---\nMy current edits\n")
	before, head := safetyRead(t, path), safetyGit(t, root, "rev-parse", "HEAD")
	manifest := filepath.Join(sessionDir(root), s.ID+".json")
	state := safetyRead(t, manifest)
	if _, err := Finish(Options{Root: root, DryRun: true}, s.ID, "Preview edits"); err != nil {
		t.Fatal(err)
	}
	if before != safetyRead(t, path) || head != safetyGit(t, root, "rev-parse", "HEAD") || state != safetyRead(t, manifest) {
		t.Fatal("dry-run finish modified repository or session")
	}
	entries, err := os.ReadDir(sessionDir(root))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Commit(Options{Root: root, DryRun: true}, "human", "copy-editor", "Preview existing edits", []string{post}); err != nil {
		t.Fatal(err)
	}
	afterEntries, err := os.ReadDir(sessionDir(root))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != len(afterEntries) || before != safetyRead(t, path) || head != safetyGit(t, root, "rev-parse", "HEAD") {
		t.Fatal("dry-run commit modified repository")
	}
}

func TestEditorialSafety_StagedRejectionBeforeRewrites(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	post := "posts/a post's.md"
	path := filepath.Join(root, post)
	safetyWrite(t, path, "---\ntitle: Post\nauthor: writer\n---\nStaged edits\n")
	safetyGit(t, root, "add", post)
	before := safetyRead(t, path)
	if _, err := Commit(Options{Root: root}, "human", "copy-editor", "Own edits", []string{post}); err == nil {
		t.Fatal("staged selected post accepted")
	}
	if safetyRead(t, path) != before {
		t.Fatal("rejected staged edit was modified")
	}
}

func TestEditorialSafety_PrimaryAuthorPersistsAndUnknownHistorySurvives(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\neditorial_history:\n  - schema_version: 99\n    future_field: keep me\n---\nBody\n")
	post := "posts/a post's.md"
	s, err := Start(Options{Root: root, DefaultAuthor: "writer"}, "human", "copy-editor", []string{post})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Finish(Options{Root: root}, s.ID, "Reviewed"); err != nil {
		t.Fatal(err)
	}
	text := safetyRead(t, filepath.Join(root, post))
	if !strings.Contains(text, "future_field: keep me") || !strings.Contains(text, "- writer") {
		t.Fatalf("writer or older history lost:\n%s", text)
	}
	if safetyRecord(t, text)["action"] != "reviewed" {
		t.Fatal("unchanged post falsely marked edited")
	}
}

func TestEditorialSafety_NestedSiteDir(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	safetyWrite(t, filepath.Join(root, "nested", "content", "post.md"), "---\ntitle: Nested\nauthor: writer\n---\nBody\n")
	safetyGit(t, root, "add", "nested")
	safetyGit(t, root, "commit", "-qm", "nested post")
	s, err := Start(Options{Root: filepath.Join(root, "nested")}, "human", "copy-editor", []string{"content/post.md"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Posts["nested/content/post.md"]; !ok {
		t.Fatalf("wrong repository-relative path: %#v", s.Posts)
	}
}

func TestEditorialSafety_RetryDetectsModifiedPendingText(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX shell hook fixture")
	}
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	post := "posts/a post's.md"
	path := filepath.Join(root, post)
	hook := filepath.Join(root, ".git", "hooks", "commit-msg")
	safetyWrite(t, hook, "#!/bin/sh\ncase \"$(cat \"$1\")\" in *'editorial: record'*) exit 1;; esac\n")
	if err := os.Chmod(hook, 0o700); err != nil {
		t.Fatal(err)
	}
	s, err := Start(Options{Root: root}, "human", "copy-editor", []string{post})
	if err != nil {
		t.Fatal(err)
	}
	safetyWrite(t, path, "---\ntitle: Post\nauthor: writer\n---\nEdited body\n")
	if _, err := Finish(Options{Root: root}, s.ID, "Edited"); err == nil {
		t.Fatal("failing attribution hook unexpectedly passed")
	}
	pending, err := Load(root, s.ID)
	if err != nil {
		t.Fatal(err)
	}
	head := safetyGit(t, root, "rev-parse", "HEAD")
	safetyWrite(t, path, safetyRead(t, path)+"Another writer's change\n")
	changed := safetyRead(t, path)
	if err := os.Remove(hook); err != nil {
		t.Fatal(err)
	}
	if _, err := Finish(Options{Root: root}, s.ID, "Retry"); err == nil {
		t.Fatal("retry committed intervening changes")
	}
	if safetyRead(t, path) != changed || safetyGit(t, root, "rev-parse", "HEAD") != head {
		t.Fatal("retry altered intervening changes")
	}
	safetyWrite(t, path, pending.Pending[post])
	done, err := Finish(Options{Root: root}, s.ID, "Retry")
	if err != nil || len(done.Commits) != 2 {
		t.Fatalf("restored pending retry failed: commits=%v err=%v", done.Commits, err)
	}
}

func TestEditorialSafety_DryRunRejectsMalformedCurrentPost(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	post := "posts/a post's.md"
	path := filepath.Join(root, post)
	safetyWrite(t, path, "---\ntitle: [broken\n---\nOwned body\n")
	before := safetyRead(t, path)
	if _, err := Commit(Options{Root: root, DryRun: true}, "human", "copy-editor", "Preview", []string{post}); err == nil {
		t.Fatal("dry-run accepted invalid working frontmatter")
	}
	if safetyRead(t, path) != before {
		t.Fatal("failed dry-run changed post")
	}
}

func TestEditorialSafety_AgentRecordsHumanEdits(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	post := "posts/a post's.md"
	safetyWrite(t, filepath.Join(root, post), "---\ntitle: Post\nauthor: writer\n---\nEdited by a human\n")
	session, err := Commit(Options{Root: root, RecordedBy: "agent-recorder"}, "human-editor", "copy-editor", "Human copy edits", []string{post})
	if err != nil {
		t.Fatal(err)
	}
	record := safetyRecord(t, safetyRead(t, filepath.Join(root, post)))
	if record["author"] != "human-editor" || record["recorded_by"] != "agent-recorder" {
		t.Fatalf("human/agent attribution conflated: %#v", record)
	}
	message := safetyGit(t, root, "show", "--format=%B", "--no-patch", session.Commits[1])
	if !strings.Contains(message, "Editorial-Author: agent-recorder") {
		t.Fatalf("metadata commit did not credit recorder: %s", message)
	}
}

func TestEditorialSafety_DirectoryAliasAndSymlinkedPost(t *testing.T) {
	root := safetyRepo(t, "---\nauthor: writer\n---\nOriginal.\n")
	alias := filepath.Join(t.TempDir(), "site")
	if err := os.Symlink(root, alias); err != nil {
		t.Skipf("directory symlinks unavailable: %v", err)
	}
	if _, err := Start(Options{Root: alias}, "editor", "copy-editor", []string{"posts/a post's.md"}); err != nil {
		t.Fatal(err)
	}
	linked := filepath.Join(root, "posts", "linked.md")
	if err := os.Symlink(filepath.Join(root, "posts", "a post's.md"), linked); err != nil {
		t.Fatal(err)
	}
	if _, err := Start(Options{Root: alias}, "editor", "copy-editor", []string{"posts/linked.md"}); err == nil {
		t.Fatal("symlinked post was accepted")
	}
}

func TestEditorialSafety_ValidateAllPostsBeforeAnyRewrite(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	first := "posts/a post's.md"
	second := "posts/zz-invalid.md"
	safetyWrite(t, filepath.Join(root, second), "---\ntitle: Second\nauthor: writer\n---\nBody\n")
	safetyGit(t, root, "add", second)
	safetyGit(t, root, "commit", "-qm", "add second post")
	s, err := Start(Options{Root: root}, "editor", "copy-editor", []string{first, second})
	if err != nil {
		t.Fatal(err)
	}
	modified := "---\ntitle: Post\nauthor: writer\n---\nHuman edit\n"
	safetyWrite(t, filepath.Join(root, first), modified)
	safetyWrite(t, filepath.Join(root, second), "---\ntitle: [broken\n---\nBody\n")
	head := safetyGit(t, root, "rev-parse", "HEAD")
	if _, err := Finish(Options{Root: root}, s.ID, "Edit both"); err == nil {
		t.Fatal("accepted invalid second post")
	}
	if got := safetyRead(t, filepath.Join(root, first)); got != modified {
		t.Fatalf("first post was rewritten before the second post validated:\n%s", got)
	}
	if head != safetyGit(t, root, "rev-parse", "HEAD") {
		t.Fatal("validation failure committed some selected files")
	}
}

func TestEditorialSafety_ResumeAfterContentCommitWithoutSavedPhase(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	post := "posts/a post's.md"
	started, err := Start(Options{Root: root}, "editor", "copy-editor", []string{post})
	if err != nil {
		t.Fatal(err)
	}
	safetyWrite(t, filepath.Join(root, post), "---\ntitle: Post\nauthor: writer\n---\nUpdated text\n")
	done, err := Finish(Options{Root: root}, started.ID, "Edited copy")
	if err != nil {
		t.Fatal(err)
	}
	if len(done.Commits) != 2 {
		t.Fatalf("expected two commits, got %v", done.Commits)
	}
	// Recreate the crash window: the content commit exists, but the manifest
	// still reports "started". The retry must adopt it, not create another.
	safetyGit(t, root, "reset", "--hard", done.Commits[0])
	if err := save(root, started); err != nil {
		t.Fatal(err)
	}
	recovered, err := Finish(Options{Root: root}, started.ID, "Edited copy")
	if err != nil {
		t.Fatalf("recover content commit: %v", err)
	}
	if recovered.Phase != phaseDone || len(recovered.Commits) != 2 || recovered.Commits[0] != done.Commits[0] {
		t.Fatalf("incorrect recovered session: %#v", recovered)
	}
}

func TestEditorialSafety_ResumeAfterMetadataCommitWithoutSavedPhase(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	post := "posts/a post's.md"
	safetyWrite(t, filepath.Join(root, post), "---\ntitle: Post\nauthor: writer\n---\nUpdated text\n")
	done, err := Commit(Options{Root: root}, "editor", "copy-editor", "Updated copy", []string{post})
	if err != nil {
		t.Fatal(err)
	}
	// Metadata committed, but its final phase transition did not persist.
	pending := done
	pending.Phase = phaseMetadataPending
	pending.Commits = append([]string(nil), done.Commits[:1]...)
	if err := save(root, pending); err != nil {
		t.Fatal(err)
	}
	sha := strings.TrimSpace(safetyGit(t, root, "rev-parse", "HEAD"))
	recovered, err := Finish(Options{Root: root}, pending.ID, "Updated copy")
	if err != nil {
		t.Fatalf("recover metadata commit: %v", err)
	}
	if recovered.Phase != phaseDone || len(recovered.Commits) != 2 ||
		recovered.Commits[1] != sha || strings.TrimSpace(safetyGit(t, root, "rev-parse", "HEAD")) != sha {
		t.Fatalf("retry created another metadata commit: %#v", recovered)
	}
}

func TestEditorialSafety_ResumePartiallyWrittenMetadata(t *testing.T) {
	root := safetyRepo(t, "---\ntitle: Post\nauthor: writer\n---\nBody\n")
	first := "posts/a post's.md"
	second := "posts/b.md"
	safetyWrite(t, filepath.Join(root, second), "---\ntitle: Second\nauthor: writer\n---\nBody\n")
	safetyGit(t, root, "add", second)
	safetyGit(t, root, "commit", "-qm", "add second post")
	s, err := Start(Options{Root: root}, "editor", "copy-editor", []string{first, second})
	if err != nil {
		t.Fatal(err)
	}
	safetyWrite(t, filepath.Join(root, first), "---\ntitle: Post\nauthor: writer\n---\nUpdated text\n")
	safetyWrite(t, filepath.Join(root, second), "---\ntitle: Second\nauthor: writer\n---\nUpdated text\n")
	done, err := Finish(Options{Root: root}, s.ID, "Edited both")
	if err != nil {
		t.Fatal(err)
	}
	safetyGit(t, root, "reset", "--hard", done.Commits[0])
	pending := done
	pending.Phase = phaseMetadataPending
	pending.Commits = append([]string(nil), done.Commits[:1]...)
	if err := save(root, pending); err != nil {
		t.Fatal(err)
	}
	safetyWrite(t, filepath.Join(root, first), pending.Pending[first])
	recovered, err := Finish(Options{Root: root}, pending.ID, "Edited both")
	if err != nil {
		t.Fatalf("complete partial metadata write: %v", err)
	}
	if recovered.Phase != phaseDone || len(recovered.Commits) != 2 ||
		safetyRead(t, filepath.Join(root, second)) != pending.Pending[second] {
		t.Fatalf("incomplete attribution retry: %#v", recovered)
	}
}
