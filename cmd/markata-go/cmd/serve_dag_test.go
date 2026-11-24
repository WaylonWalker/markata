package cmd

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestServeRegistersDAGFlag(t *testing.T) {
	flag := serveCmd.Flags().Lookup("dag")
	if flag == nil {
		t.Fatal("serve --dag flag is not registered")
	}
	if flag.DefValue != envValueDisabled {
		t.Fatalf("serve --dag default = %q, want false", flag.DefValue)
	}
}

func TestDAGBuildPreservesServeStageObserver(t *testing.T) {
	previous := buildDAG
	buildDAG = true
	t.Cleanup(func() { buildDAG = previous })
	t.Setenv(dagBuildEnv, envValueDisabled)

	manager := lifecycle.NewManager()
	got := make([]string, 0, len(dagLifecycleStages)*2)
	_, err := runBuildObserved(manager, func(stage lifecycle.Stage, starting bool, stageErr error) {
		if stageErr != nil {
			t.Fatalf("observer stage %s error = %v", stage, stageErr)
		}
		phase := "finish"
		if starting {
			phase = "start"
		}
		got = append(got, fmt.Sprintf("%s:%s", stage, phase))
	})
	if err != nil {
		t.Fatalf("runBuildObserved() with DAG = %v", err)
	}

	want := make([]string, 0, len(dagLifecycleStages)*2)
	for _, stage := range dagLifecycleStages {
		want = append(want, fmt.Sprintf("%s:start", stage), fmt.Sprintf("%s:finish", stage))
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("observer events = %v, want %v", got, want)
	}
}
