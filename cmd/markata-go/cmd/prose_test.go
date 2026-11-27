package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/prose"
)

func TestProseCheckCLIHonorsProjectPhrasePreferences(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "draft.md")
	if err := os.WriteFile(file, []byte("Very unique and completely revolutionary."), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := prose.IgnorePhrase(file, "prose/intensifier", "very unique"); err != nil {
		t.Fatal(err)
	}
	oldJSON, oldIgnore := proseJSON, proseIgnore
	t.Cleanup(func() {
		proseJSON, proseIgnore = oldJSON, oldIgnore
		proseCheckCmd.SetOut(nil)
	})
	proseJSON = true
	proseIgnore = nil

	var output bytes.Buffer
	proseCheckCmd.SetOut(&output)
	if err := proseCheckCmd.RunE(proseCheckCmd, []string{file}); err != nil {
		t.Fatal(err)
	}
	var issues []prose.Issue
	if err := json.Unmarshal(output.Bytes(), &issues); err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Rule != "prose/intensifier" ||
		issues[0].Column != len("Very unique and ") {
		t.Fatalf("CLI lost the project phrase ignore: %+v", issues)
	}
	if issues[0].Provider != "markata" || issues[0].ContractVersion != prose.ContractVersion {
		t.Fatalf("CLI lost diagnostic contract metadata: %+v", issues[0])
	}
}
