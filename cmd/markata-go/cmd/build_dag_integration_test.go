package cmd

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildlab"
)

func TestBuildDAGMatchesLegacyInBuildLab(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := filepath.Join(moduleRoot(t), "cmd", "markata-go", "cmd", "testdata", "buildlab-site")
	binary := buildTestBinary(t)
	result, runErr := buildlab.RunScenario(context.Background(), buildlab.ScenarioRunConfig{
		Fixture: fixture,
		Scenario: buildlab.Scenario{ID: "cli-serial-dag-equivalence", Version: "1", Operations: []buildlab.Operation{
			{Type: buildlab.OpClearCache},
			{Type: buildlab.OpClearOutput},
			{Type: buildlab.OpBuild},
			{Type: buildlab.OpReplaceExact, Path: "content/target.md", Old: "Original target content.", New: "Updated target content."},
			{Type: buildlab.OpBuild},
		}},
		Baseline: buildlab.BuildCommand{
			Binary:    binary,
			Args:      []string{"build", "-c", "markata-go.toml"},
			OutputDir: "output",
			Timeout:   5 * time.Minute,
			Env:       []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"},
		},
		Candidate: buildlab.BuildCommand{
			Binary:    binary,
			Args:      []string{"build", "--dag", "-c", "markata-go.toml"},
			OutputDir: "output",
			Timeout:   5 * time.Minute,
			Env:       []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"},
		},
		Classes: map[string]buildlab.OutputClass{
			".markata/diagnostics.json": buildlab.ClassVolatile,
			".well-known/time":          buildlab.ClassVolatile,
		},
		CheckDeterminism: true,
		GOMAXPROCS:       1,
	})
	if runErr != nil {
		t.Fatalf("Build Lab run error = %v", runErr)
	}
	if result.Verdict != buildLabPassVerdict {
		t.Fatalf("Build Lab verdict = %s, want %s; diagnostics = %+v", result.Verdict, buildLabPassVerdict, result.Diagnostics)
	}
	if len(result.Checkpoints) != 2 {
		t.Fatalf("checkpoints = %d, want 2", len(result.Checkpoints))
	}
	for index := range result.Checkpoints {
		correctness := result.Checkpoints[index].Correctness
		if !correctness.DifferentialEqual {
			t.Fatalf("checkpoint %d legacy-vs-DAG mismatch: %+v", index, correctness.DifferentialDiff)
		}
		if correctness.IncrementalApplicable && !correctness.IncrementalEqual {
			t.Fatalf("checkpoint %d DAG incremental mismatch: %+v", index, correctness)
		}
		if !correctness.DeterministicEqual {
			t.Fatalf("checkpoint %d DAG determinism mismatch: %+v", index, correctness.DeterminismDiff)
		}
	}
}
