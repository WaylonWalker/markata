package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestPrintFirstRunSkillIsStableAndSelfContained(t *testing.T) {
	var first, second bytes.Buffer
	if err := printFirstRunSkill(&first); err != nil {
		t.Fatal(err)
	}
	if err := printFirstRunSkill(&second); err != nil {
		t.Fatal(err)
	}
	if first.String() != second.String() {
		t.Fatal("agent-first setup guide changed between reads")
	}
	for _, required := range []string{
		"Markata — agent-first site setup",
		"markata config init",
		"markata build",
		"markata agent install",
		"template",
		"plugin",
		"deployment",
		"explicitly approves",
	} {
		if !strings.Contains(first.String(), required) {
			t.Errorf("bundled guide missing %q", required)
		}
	}
	if first.Len() == 0 || first.Bytes()[first.Len()-1] != '\n' {
		t.Fatal("guide must be nonempty newline-terminated text")
	}
}

type failingSkillWriter struct{}

func (failingSkillWriter) Write([]byte) (int, error) {
	return 0, errors.New("broken output")
}

func TestPrintFirstRunSkillReturnsWriteError(t *testing.T) {
	if err := printFirstRunSkill(failingSkillWriter{}); err == nil {
		t.Fatal("expected stdout writer failure to propagate")
	}
}

func TestRootSkillNeedsNoWorkspace(t *testing.T) {
	old := rootSkillRequested
	rootSkillRequested = true
	t.Cleanup(func() { rootSkillRequested = old })

	var output bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)
	if err := runRootCommand(cmd, nil); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(output.String(), "# Markata") {
		t.Fatalf("unexpected guide opening: %q", output.String()[:min(output.Len(), 70)])
	}
	if err := runRootCommand(cmd, []string{"post.md"}); err == nil {
		t.Fatal("--skill must refuse file arguments")
	}
}
