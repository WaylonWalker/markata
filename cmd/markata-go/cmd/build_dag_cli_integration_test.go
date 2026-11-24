package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildlab"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

// dagCLIFastConfig turns off the plugins that dominate a fixture build but
// play no part in executor selection: the Web Awesome vendor download (a
// network request) and CSS/JS minification. Builds then take a fraction of a
// second locally instead of seconds, which keeps slow, contended CI runners
// well inside the build timeout.
const dagCLIFastConfig = `[markata-go.webawesome]
enabled = false

[markata-go.css_minify]
enabled = false

[markata-go.js_minify]
enabled = false
`

// dagCLIBuildTimeout bounds one CLI build. It uses the time left before the
// test binary's deadline, so `go test -timeout` governs slow runners, and
// falls back to a fixed limit when no deadline is set.
func dagCLIBuildTimeout(t *testing.T) time.Duration {
	t.Helper()
	const fallback = 5 * time.Minute
	deadline, ok := t.Deadline()
	if !ok {
		return fallback
	}
	// Leave a margin so a timed-out build is reported before the test binary
	// panics on its own deadline.
	remaining := time.Until(deadline) - 30*time.Second
	if remaining < 30*time.Second {
		return 30 * time.Second
	}
	return min(remaining, fallback)
}

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
			fastConfig := filepath.Join(workspace.Root, "dag-cli-fast.toml")
			if err := os.WriteFile(fastConfig, []byte(dagCLIFastConfig), 0o600); err != nil {
				t.Fatal(err)
			}
			benchmark := filepath.Join(workspace.Root, "benchmark.json")
			args := append([]string{
				"build", "--no-color", "--no-input",
				"--merge-config", fastConfig,
				"--merge-config", workspace.IsolationConfig,
				"--benchmark-json=" + benchmark,
			}, test.flags...)
			timeout := dagCLIBuildTimeout(t)
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			command := exec.CommandContext(ctx, binary, args...)
			command.Dir = workspace.SiteDir
			command.Env = environment
			output, err := command.CombinedOutput()
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				t.Fatalf("build timed out after %s (killed, not a build failure):\n%s", timeout, output)
			}
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
