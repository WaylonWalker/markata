package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/plugins"
)

func TestRenderCommand_JSON(t *testing.T) {
	cmd := newRenderCommand()
	cmd.SetIn(strings.NewReader("# Test\n"))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--json", "--inline-assets"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var result plugins.StandaloneResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.SchemaVersion != 1 || !strings.Contains(result.HTML, "Test") || !strings.Contains(result.HeadHTML, "<style>") {
		t.Fatal("incomplete result")
	}
}
