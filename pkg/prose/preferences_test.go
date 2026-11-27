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

func TestProjectDictionaryAddRemoveAndNormalization(t *testing.T) {
	root := t.TempDir()
	post := filepath.Join(root, "draft.md")
	if _, err := AddWord(post, "QuickSync"); err != nil {
		t.Fatal(err)
	}
	if _, err := AddWord(post, "quicksync"); err != nil {
		t.Fatal(err)
	}
	if _, err := AddWord(post, "k3s"); err != nil {
		t.Fatal(err)
	}
	prefs, err := LoadPreferences(post)
	if err != nil {
		t.Fatal(err)
	}
	if len(prefs.AcceptedWords) != 2 ||
		!prefs.Options().AcceptedWords["quicksync"] ||
		!prefs.Options().AcceptedWords["k3s"] {
		t.Fatalf("accepted terms were not normalized and deduplicated: %+v", prefs)
	}
	if _, err := RemoveWord(post, "QUICKSYNC"); err != nil {
		t.Fatal(err)
	}
	prefs, err = LoadPreferences(post)
	if err != nil {
		t.Fatal(err)
	}
	if len(prefs.AcceptedWords) != 1 || prefs.AcceptedWords[0] != "k3s" {
		t.Fatalf("removal lost unrelated accepted terms: %+v", prefs)
	}
}

func TestProjectDictionaryRejectsUnsafeWords(t *testing.T) {
	for _, word := range []string{"", "two words", "../other", "\nword", "name\x00", "über text", " 	"} {
		if _, err := CanonicalWord(word); err == nil {
			t.Fatalf("accepted unsafe dictionary entry %q", word)
		}
	}
	for _, word := range []string{"QuïckSync", "k3s", "C++", "Rustfmt", "can't"} {
		if _, err := CanonicalWord(word); err != nil {
			t.Fatalf("rejected valid technical word %q: %v", word, err)
		}
	}
}

func TestRootDictionaryDoesNotWriteAncestor(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "project")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(parent, PreferencesFile),
		[]byte("{\"accepted_words\":[\"personalword\"]}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	name, err := AddWordAtRoot(root, "QuickSync")
	if err != nil {
		t.Fatal(err)
	}
	if name != filepath.Join(root, PreferencesFile) {
		t.Fatalf("root write escaped workspace: %s", name)
	}
	contents, err := os.ReadFile(filepath.Join(parent, PreferencesFile))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(contents, []byte("{\"accepted_words\":[\"personalword\"]}\n")) {
		t.Fatalf("modified parent dictionary: %s", contents)
	}
}
