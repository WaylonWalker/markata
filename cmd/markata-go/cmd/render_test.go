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

func TestRenderCommand_Capabilities(t *testing.T) {
	cmd := newRenderCommand()
	cmd.SetIn(strings.NewReader("this is not rendered"))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--capabilities"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var capabilities plugins.StandaloneCapabilities
	if err := json.Unmarshal(out.Bytes(), &capabilities); err != nil {
		t.Fatal(err)
	}
	if capabilities.SchemaVersion != plugins.StandaloneSchemaVersion ||
		capabilities.Renderer != "markata-standalone" || capabilities.LoadsSiteConfig ||
		capabilities.ExternalURLAudit {
		t.Fatalf("unexpected capabilities: %+v", capabilities)
	}
	if strings.Contains(out.String(), "this is not rendered") {
		t.Fatal("capability discovery should not read or render stdin")
	}
}

func TestRenderCommand_CapabilitiesConflict(t *testing.T) {
	cmd := newRenderCommand()
	cmd.SetArgs([]string{"--capabilities", "--appearance"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("expected conflicting mode error, got %v", err)
	}
}
