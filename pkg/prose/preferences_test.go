package prose

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectPreferencesFollowMarkdownFile(t *testing.T) {
	root := t.TempDir()
	post := filepath.Join(root, "posts", "sample.md")
	if err := os.MkdirAll(filepath.Dir(post), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(post, []byte("A rich tapestry."), 0o644); err != nil {
		t.Fatal(err)
	}
	preferencesPath := filepath.Join(root, PreferencesFile)
	if err := os.WriteFile(preferencesPath, []byte("{\"ignored_rules\":[\"prose/cliche\"]}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := LoadPreferences(post)
	if err != nil {
		t.Fatal(err)
	}
	if len(Check("A rich tapestry.", p.Options())) != 0 {
		t.Fatal("shared preferences did not suppress the style rule")
	}
	name, err := IgnoreRule(post, "prose/wordiness")
	if err != nil {
		t.Fatal(err)
	}
	if name != preferencesPath {
		t.Fatalf("wrote outside the discovered project: %s", name)
	}
	if _, err := IgnoreRule(post, "prose/wordiness"); err != nil {
		t.Fatal(err)
	}
	p, err = LoadPreferences(post)
	if err != nil {
		t.Fatal(err)
	}
	if len(p.IgnoredRules) != 2 {
		t.Fatalf("expected deduplicated project preferences: %+v", p)
	}
}

func TestMalformedPreferencesReturnError(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, PreferencesFile), []byte("{oops"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPreferences(filepath.Join(root, "test.md")); err == nil {
		t.Fatal("expected parse error for malformed preferences")
	}
}
