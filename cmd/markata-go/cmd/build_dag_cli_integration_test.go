package cmd

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildlab"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestDAGCLISelection(t *testing.T) {
	binary := buildTestBinary(t)
	for _, test := range []struct {
		name     string
		env      string
		flags    []string
		executor lifecycle.BuildExecutor
		wantErr  bool
		dryRun   bool
	}{
		{name: "environment opt-in", env: "true", executor: lifecycle.BuildExecutorDAG},
		{name: "explicit opt-out", env: "true", flags: []string{"--dag=false"}, executor: lifecycle.BuildExecutorLegacy},
		{name: "invalid environment", env: "invalid", wantErr: true},
		{name: "invalid overridden false", env: "invalid", flags: []string{"--dag=false"}, executor: lifecycle.BuildExecutorLegacy},
		{name: "invalid overridden true", env: "invalid", flags: []string{"--dag"}, executor: lifecycle.BuildExecutorDAG},
		{name: "dry-run opt-in", env: "true", flags: []string{"--dag", "--dry-run"}, dryRun: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			workspace, err := buildlab.NewWorkspace(filepath.Join(moduleRoot(t), "cmd", "markata-go", "cmd", "testdata", "buildlab-site"), t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			environment, err := workspace.Environment([]string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}, 2)
			if err != nil {
				t.Fatal(err)
			}
			environment = append(environment, dagBuildEnv+"="+test.env)
			benchmark := filepath.Join(workspace.Root, "benchmark.json")
			args := append([]string{"build", "--no-color", "--no-input", "--merge-config", workspace.IsolationConfig, "--benchmark-json=" + benchmark}, test.flags...)
			ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, binary, args...)
			command.Dir = workspace.SiteDir
			command.Env = environment
			output, err := command.CombinedOutput()
			if test.wantErr {
				if err == nil || !strings.Contains(string(output), dagBuildEnv+" must be a boolean") {
					t.Fatalf("expected invalid environment error, got %v:\n%s", err, output)
				}
				return
			}
			if err != nil {
				t.Fatalf("build failed: %v:\n%s", err, output)
			}
			wantLabel := !test.dryRun && test.executor == lifecycle.BuildExecutorDAG
			if strings.Contains(string(output), "serial DAG (experimental)") != wantLabel {
				t.Fatalf("unexpected executor label:\n%s", output)
			}
			if test.dryRun {
				if _, err := os.Stat(filepath.Join(workspace.SiteDir, "output", ".markata", "diagnostics.json")); !os.IsNotExist(err) {
					t.Fatalf("dry-run produced diagnostics or failed to inspect them: %v", err)
				}
				return
			}
			data, err := os.ReadFile(benchmark)
			if err != nil {
				t.Fatal(err)
			}
			var result struct {
				Executor lifecycle.BuildExecutor `json:"executor"`
			}
			if err := json.Unmarshal(data, &result); err != nil {
				t.Fatal(err)
			}
			if result.Executor != test.executor {
				t.Fatalf("actual executor = %q, want %q", result.Executor, test.executor)
			}
		})
	}
}
