package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/prose"
)

func TestProseDictionaryCLIAddListRemove(t *testing.T) {
	file := filepath.Join(t.TempDir(), "draft.md")
	var buf bytes.Buffer
	proseDictionaryAddCmd.SetOut(&buf)
	t.Cleanup(func() {
		proseDictionaryAddCmd.SetOut(nil)
		proseDictionaryListCmd.SetOut(nil)
		proseDictionaryRemoveCmd.SetOut(nil)
	})
	if err := proseDictionaryAddCmd.RunE(proseDictionaryAddCmd, []string{"QuickSync", file}); err != nil {
		t.Fatal(err)
	}
	prefs, err := prose.LoadPreferences(file)
	if err != nil || !prefs.Options().AcceptedWords["quicksync"] {
		t.Fatalf("CLI did not persist word: %+v, %v", prefs, err)
	}
	buf.Reset()
	proseDictionaryListCmd.SetOut(&buf)
	if err := proseDictionaryListCmd.RunE(proseDictionaryListCmd, []string{file}); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(buf.String()) != "quicksync" {
		t.Fatalf("unexpected dictionary list: %q", buf.String())
	}
	buf.Reset()
	proseDictionaryRemoveCmd.SetOut(&buf)
	if err := proseDictionaryRemoveCmd.RunE(proseDictionaryRemoveCmd, []string{"QUICKSYNC", file}); err != nil {
		t.Fatal(err)
	}
	prefs, err = prose.LoadPreferences(file)
	if err != nil || len(prefs.AcceptedWords) != 0 {
		t.Fatalf("CLI did not revoke accepted word: %+v, %v", prefs, err)
	}
}
