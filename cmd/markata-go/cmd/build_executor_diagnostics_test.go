package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestBenchmarkJSONIncludesExecutor(t *testing.T) {
	previous := buildDAG
	t.Cleanup(func() { buildDAG = previous })
	t.Setenv(dagBuildEnv, envValueDisabled)

	for _, test := range []struct {
		name string
		dag  bool
		want lifecycle.BuildExecutor
	}{
		{name: "legacy", want: lifecycle.BuildExecutorLegacy},
		{name: "dag", dag: true, want: lifecycle.BuildExecutorDAG},
	} {
		t.Run(test.name, func(t *testing.T) {
			buildDAG = !test.dag
			var output bytes.Buffer
			if err := writeBenchmarkJSON(&output, &BuildResult{Executor: test.want}); err != nil {
				t.Fatalf("writeBenchmarkJSON() = %v", err)
			}
			var decoded struct {
				Executor lifecycle.BuildExecutor `json:"executor"`
			}
			if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
				t.Fatalf("benchmark JSON = %v", err)
			}
			if decoded.Executor != test.want {
				t.Fatalf("executor = %q, want %q", decoded.Executor, test.want)
			}
		})
	}
}

func TestBenchmarkJSONRejectsInvalidExecutor(t *testing.T) {
	var output bytes.Buffer
	if err := writeBenchmarkJSON(&output, &BuildResult{Executor: "invalid"}); err == nil {
		t.Fatal("invalid executor was serialized as a successful benchmark")
	}
}

func TestRunDAGBuildRecordsExecutorOnManager(t *testing.T) {
	manager := lifecycle.NewManager()
	if _, err := runDAGBuildObserved(manager, nil); err != nil {
		t.Fatalf("runDAGBuildObserved() = %v", err)
	}
	if got := manager.BuildExecutor(); got != lifecycle.BuildExecutorDAG {
		t.Fatalf("BuildExecutor() = %q, want %q", got, lifecycle.BuildExecutorDAG)
	}
}

func TestLegacyBuildResetsDAGManagerIdentity(t *testing.T) {
	previous, previousCommand := buildDAG, currentCmd
	buildDAG, currentCmd = false, nil
	t.Cleanup(func() { buildDAG, currentCmd = previous, previousCommand })
	t.Setenv(dagBuildEnv, "false")
	manager := lifecycle.NewManager()
	if _, err := runDAGBuildObserved(manager, nil); err != nil {
		t.Fatal(err)
	}
	result, err := runBuildObserved(manager, nil)
	if err != nil {
		t.Fatal(err)
	}
	if result.Executor != lifecycle.BuildExecutorLegacy || manager.BuildExecutor() != lifecycle.BuildExecutorLegacy {
		t.Fatalf("legacy build retained DAG identity: result=%q manager=%q", result.Executor, manager.BuildExecutor())
	}
}
