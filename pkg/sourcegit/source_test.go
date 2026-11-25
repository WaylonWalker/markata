package sourcegit

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/fileidentity"
)

func TestReadSourceState(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "post.md"), "one")
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Content Index Test")
	git(t, dir, "add", "post.md")
	git(t, dir, "commit", "-m", "initial")

	state, err := Read(context.Background(), dir)
	if err != nil || state.Commit == "" || state.Dirty == nil || *state.Dirty {
		t.Fatalf("clean state = %#v, err = %v", state, err)
	}

	writeFile(t, filepath.Join(dir, "post.md"), "changed")
	dirtyState, err := Read(context.Background(), dir)
	if err != nil || dirtyState.Dirty == nil || !*dirtyState.Dirty || state.Equal(dirtyState) {
		t.Fatalf("tracked dirty state = %#v, err = %v", state, err)
	}
	writeFile(t, filepath.Join(dir, "post.md"), "changed again")
	changedDirtyState, err := Read(context.Background(), dir)
	if err != nil || dirtyState.Equal(changedDirtyState) {
		t.Fatalf("changed dirty state was not detected: %#v -> %#v, err = %v", dirtyState, changedDirtyState, err)
	}
	state = dirtyState

	git(t, dir, "add", "post.md")
	state, err = Read(context.Background(), dir)
	if err != nil || state.Dirty == nil || !*state.Dirty {
		t.Fatalf("staged dirty state = %#v, err = %v", state, err)
	}

	git(t, dir, "reset", "--hard", "HEAD")
	if err := os.Remove(filepath.Join(dir, "post.md")); err != nil {
		t.Fatal(err)
	}
	state, err = Read(context.Background(), dir)
	if err != nil || state.Dirty == nil || !*state.Dirty {
		t.Fatalf("deleted dirty state = %#v, err = %v", state, err)
	}
	git(t, dir, "restore", "post.md")
	writeFile(t, filepath.Join(dir, "new.md"), "untracked source")
	state, err = Read(context.Background(), dir)
	if err != nil || state.Dirty == nil || !*state.Dirty {
		t.Fatalf("untracked dirty state = %#v, err = %v", state, err)
	}
}

func TestReadSourceStateDetectsTrackedModeChanges(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Git executable mode bits are not portable on Windows")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "post.md")
	writeFile(t, path, "one")
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Content Index Test")
	git(t, dir, "add", "post.md")
	git(t, dir, "commit", "-m", "initial")
	clean, err := Read(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}
	modeChanged, err := Read(context.Background(), dir)
	if err != nil {
		t.Fatal(err)
	}
	if clean.Equal(modeChanged) {
		t.Fatal("tracked mode change was not detected")
	}
}

func TestReadMetadataSourceStateDetectsSameSizeSameMtimeRewrite(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("metadata fingerprints require Linux file identity")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "post.md")
	writeFile(t, path, "one!")
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Content Index Test")
	git(t, dir, "add", "post.md")
	git(t, dir, "commit", "-m", "initial")

	before, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(10 * time.Millisecond)
	writeFile(t, path, "two!")
	if err := os.Chtimes(path, info.ModTime(), info.ModTime()); err != nil {
		t.Fatal(err)
	}
	after, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if before.Equal(after) {
		t.Fatal("same-size, same-mtime source rewrite did not change metadata fingerprint")
	}
}

func TestReadFilteredMetadataTracksUntrackedIgnoredAndRenamedPaths(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("metadata fingerprints require Linux file identity")
	}
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".gitignore"), "ignored.md\n")
	writeFile(t, filepath.Join(dir, "post.md"), "tracked post")
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Content Index Test")
	git(t, dir, "add", ".gitignore", "post.md")
	git(t, dir, "commit", "-m", "initial")

	clean, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "new.md"), "untracked")
	untracked, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil || untracked.Dirty == nil || !*untracked.Dirty || clean.Equal(untracked) {
		t.Fatalf("untracked metadata state = %#v, err = %v", untracked, err)
	}
	if err := os.Remove(filepath.Join(dir, "new.md")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "ignored.md"), "ignored")
	ignored, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil || ignored.Dirty == nil || !*ignored.Dirty || clean.Equal(ignored) {
		t.Fatalf("ignored Markdown metadata state = %#v, err = %v", ignored, err)
	}
	if err := os.Remove(filepath.Join(dir, "ignored.md")); err != nil {
		t.Fatal(err)
	}
	git(t, dir, "mv", "post.md", "renamed.md")
	renamed, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil || renamed.Dirty == nil || !*renamed.Dirty || clean.Equal(renamed) {
		t.Fatalf("renamed metadata state = %#v, err = %v", renamed, err)
	}
}

func TestReadMetadataRejectsGitIndexFlagsThatHideChanges(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("metadata fingerprints require Linux file identity")
	}
	for _, flag := range []string{"--assume-unchanged", "--skip-worktree"} {
		t.Run(strings.TrimPrefix(flag, "--"), func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "post.md")
			writeFile(t, path, "tracked post")
			git(t, dir, "init")
			git(t, dir, "config", "user.email", "test@example.invalid")
			git(t, dir, "config", "user.name", "Content Index Test")
			git(t, dir, "add", "post.md")
			git(t, dir, "commit", "-m", "initial")
			git(t, dir, "update-index", flag, "post.md")
			writeFile(t, path, "changed post")
			if _, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{}); err == nil {
				t.Fatal("metadata snapshot accepted an index flag that can hide tracked changes")
			}
		})
	}
}

func TestReadMetadataRejectsGitStatSettingsThatCanHideChanges(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("metadata fingerprints require Linux file identity")
	}
	for _, setting := range []struct {
		name  string
		key   string
		value string
	}{
		{name: "trustctime_false", key: "core.trustctime", value: "false"},
		{name: "checkstat_minimal", key: "core.checkStat", value: "minimal"},
		{name: "fsmonitor_enabled", key: "core.fsmonitor", value: "true"},
	} {
		t.Run(setting.name, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, filepath.Join(dir, "post.md"), "tracked post")
			git(t, dir, "init")
			git(t, dir, "config", "user.email", "test@example.invalid")
			git(t, dir, "config", "user.name", "Content Index Test")
			git(t, dir, "add", "post.md")
			git(t, dir, "commit", "-m", "initial")
			git(t, dir, "config", setting.key, setting.value)
			if _, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{}); err == nil {
				t.Fatal("metadata snapshot accepted a Git stat setting that can hide changed files")
			}
		})
	}
}

func TestReadFilteredMetadataRejectsGitOnlyIgnoredPaths(t *testing.T) {
	if !fileidentity.Available {
		t.Skip("metadata fingerprints require Linux file identity")
	}
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "post.md"), "tracked post")
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Content Index Test")
	git(t, dir, "add", "post.md")
	git(t, dir, "commit", "-m", "initial")
	writeFile(t, filepath.Join(dir, ".git", "info", "exclude"), "hidden/\nhidden.md\n")
	if err := os.MkdirAll(filepath.Join(dir, "hidden"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "hidden", "post.md"), "hidden input")
	writeFile(t, filepath.Join(dir, "hidden.md"), "hidden input")
	_, err := ReadMetadataWithOptions(context.Background(), dir, ReadOptions{
		IgnoredPathFilter: func(string) (bool, error) { return false, nil },
	})
	if err == nil {
		t.Fatal("Git-only excludes hid build paths from a filtered metadata snapshot")
	}
}

func TestReadSourceStateWithoutGit(t *testing.T) {
	if _, err := Read(context.Background(), t.TempDir()); err == nil {
		t.Fatal("expected unavailable Git state error")
	}
}

func TestReadSourceStateDetectsIgnoredMarkdownChanges(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".gitignore"), "ignored.md\n")
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Content Index Test")
	git(t, dir, "add", ".gitignore")
	git(t, dir, "commit", "-m", "initial")
	writeFile(t, filepath.Join(dir, "ignored.md"), "one")
	first, err := Read(context.Background(), dir)
	if err != nil || first.Dirty == nil || !*first.Dirty {
		t.Fatalf("ignored source was not marked dirty: %#v, %v", first, err)
	}
	writeFile(t, filepath.Join(dir, "ignored.md"), "two")
	second, err := Read(context.Background(), dir)
	if err != nil || first.Equal(second) {
		t.Fatalf("ignored source change was not detected: %#v -> %#v, %v", first, second, err)
	}
}

func TestReadSourceStateIgnoredDirectoryDependsOnInputPolicy(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".gitignore"), "ignored/\n")
	git(t, dir, "init")
	git(t, dir, "config", "user.email", "test@example.invalid")
	git(t, dir, "config", "user.name", "Content Index Test")
	git(t, dir, "add", ".gitignore")
	git(t, dir, "commit", "-m", "initial")
	if err := os.MkdirAll(filepath.Join(dir, "ignored"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "ignored", "post.md"), "one")

	filtered, err := ReadWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if filtered.Dirty == nil || *filtered.Dirty {
		t.Fatalf("ignored generated directory marked filtered source dirty: %#v", filtered)
	}
	writeFile(t, filepath.Join(dir, "ignored", "post.md"), "two")
	filteredAfter, err := ReadWithOptions(context.Background(), dir, ReadOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if !filtered.Equal(filteredAfter) {
		t.Fatal("ignored directory changed a filtered source state")
	}

	inputs, err := ReadWithOptions(context.Background(), dir, ReadOptions{IncludeIgnoredContent: true})
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "ignored", "post.md"), "three")
	inputsAfter, err := ReadWithOptions(context.Background(), dir, ReadOptions{IncludeIgnoredContent: true})
	if err != nil {
		t.Fatal(err)
	}
	if inputs.Equal(inputsAfter) {
		t.Fatal("ignored source input change was not detected")
	}
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	if output, err := Command(context.Background(), dir, args...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
