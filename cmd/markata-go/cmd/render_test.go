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
	for _, flags := range [][]string{
		{"--capabilities", "--appearance"},
		{"--capabilities", "--workspace", "some-site"},
	} {
		cmd := newRenderCommand()
		cmd.SetArgs(flags)
		if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
			t.Fatalf("flags %v: expected conflicting mode error, got %v", flags, err)
		}
	}
}

func TestRenderCommand_Stdio(t *testing.T) {
	cmd := newRenderCommand()
	cmd.SetIn(strings.NewReader("{\"id\":\"a\",\"op\":\"capabilities\"}\n{\"id\":\"b\",\"op\":\"render\",\"source\":\"# OK\"}\n"))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--stdio"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected two JSONL responses: %s", out.String())
	}
	var caps map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &caps); err != nil {
		t.Fatal(err)
	}
	if caps["ok"] != true || caps["capabilities"] == nil {
		t.Fatalf("missing capabilities: %+v", caps)
	}
	var response map[string]any
	if err := json.Unmarshal([]byte(lines[1]), &response); err != nil {
		t.Fatal(err)
	}
	if response["ok"] != true || response["result"] == nil {
		t.Fatalf("missing render result: %+v", response)
	}
}

func TestRenderCommand_StdioRejectsConflictingOptions(t *testing.T) {
	for _, flags := range [][]string{
		{"--stdio", "--unsafe"},
		{"--stdio", "--inline-assets"},
	} {
		cmd := newRenderCommand()
		cmd.SetArgs(flags)
		if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "per-request JSON options") {
			t.Fatalf("flags %v: expected stdio option error, got %v", flags, err)
		}
	}
}

func TestRenderCommand_StdioWorkspace(t *testing.T) {
	cmd := newRenderCommand()
	cmd.SetIn(strings.NewReader("{\"id\":\"workspace\",\"op\":\"render\",\"source\":\"# Configured\"}\n"))
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--stdio", "--workspace", t.TempDir()})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	var response struct {
		OK     bool                      `json:"ok"`
		Result *plugins.StandaloneResult `json:"result"`
	}
	if err := json.Unmarshal(out.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if !response.OK || response.Result == nil || !strings.Contains(response.Result.HTML, "Configured") {
		t.Fatalf("workspace stdio render failed: %s", out.String())
	}
}
