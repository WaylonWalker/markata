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

func TestRenderCommand_Profiles(t *testing.T) {
	t.Run("offline interactive", func(t *testing.T) {
		cmd := newRenderCommand()
		cmd.SetIn(strings.NewReader("# Hello"))
		var output bytes.Buffer
		cmd.SetOut(&output)
		cmd.SetArgs([]string{"--json", "--profile", "offline-interactive"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
		var result plugins.StandaloneResult
		if err := json.Unmarshal(output.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if result.RenderProfile != plugins.RenderProfileOfflineInteractive || !strings.Contains(result.HeadHTML, "<style>") {
			t.Fatalf("offline rendering policy was not applied: %+v", result)
		}
	})
	t.Run("static must not silently use JavaScript", func(t *testing.T) {
		cmd := newRenderCommand()
		cmd.SetIn(strings.NewReader("# Hello"))
		cmd.SetArgs([]string{"--profile", "offline-static"})
		if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "not supported") {
			t.Fatalf("expected unsupported static profile, got %v", err)
		}
	})
	t.Run("site profile requires explicit config", func(t *testing.T) {
		cmd := newRenderCommand()
		cmd.SetArgs([]string{"--site-render-profile"})
		if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "requires --config") {
			t.Fatalf("missing explicit config should fail, got %v", err)
		}
	})
}
