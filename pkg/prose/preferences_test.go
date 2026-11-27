package prose

import (
	"bytes"
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
	if err := os.WriteFile(post, []byte("A rich tapestry."), 0o600); err != nil {
		t.Fatal(err)
	}
	preferencesPath := filepath.Join(root, PreferencesFile)
	if err := os.WriteFile(preferencesPath, []byte("{\"ignored_rules\":[\"prose/cliche\"]}\n"), 0o600); err != nil {
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
	if err := os.WriteFile(filepath.Join(root, PreferencesFile), []byte("{oops"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadPreferences(filepath.Join(root, "test.md")); err == nil {
		t.Fatal("expected parse error for malformed preferences")
	}
}

func TestIgnoreRuleAtRootDoesNotWriteAncestor(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "site")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	ancestor := filepath.Join(parent, PreferencesFile)
	original := []byte("{\"ignored_rules\":[\"prose/wordiness\"]}\n")
	if err := os.WriteFile(ancestor, original, 0o600); err != nil {
		t.Fatal(err)
	}
	name, err := IgnoreRuleAtRoot(root, "prose/cliche")
	if err != nil {
		t.Fatal(err)
	}
	if name != filepath.Join(root, PreferencesFile) {
		t.Fatalf("browser ignore escaped root: %q", name)
	}
	contents, err := os.ReadFile(ancestor)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(contents, original) {
		t.Fatalf("ancestor preference was changed: %s", contents)
	}
	prefs, err := LoadPreferences(filepath.Join(root, "post.md"))
	if err != nil || !prefs.Options().IgnoredRules["prose/cliche"] {
		t.Fatalf("root preference missing: %+v, err=%v", prefs, err)
	}
}
