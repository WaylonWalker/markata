package cmd

import (
	"context"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/builddag"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

type dagLoadStageTestPlugin struct {
	name     string
	priority int
	calls    *[]string
}

func (p *dagLoadStageTestPlugin) Name() string { return p.name }

func (p *dagLoadStageTestPlugin) Priority(lifecycle.Stage) int { return p.priority }

func (p *dagLoadStageTestPlugin) Load(*lifecycle.Manager) error {
	*p.calls = append(*p.calls, p.name)
	return nil
}

func TestDAGLoadStagePreservesPluginOrderAndObserver(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	manager.RegisterPlugins(
		&dagLoadStageTestPlugin{name: "late", priority: lifecycle.PriorityLate, calls: &calls},
		&dagLoadStageTestPlugin{name: "early", priority: lifecycle.PriorityEarly, calls: &calls},
	)

	observerEvents := []string{}
	observe := func(stage lifecycle.Stage, starting bool, err error) {
		if err != nil {
			t.Fatalf("observer error for %s: %v", stage, err)
		}
		if stage != lifecycle.StageLoad {
			t.Fatalf("observer stage = %s, want %s", stage, lifecycle.StageLoad)
		}
		if starting {
			observerEvents = append(observerEvents, "start")
			return
		}
		observerEvents = append(observerEvents, "finish")
	}

	builder := builddag.NewBuilder()
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageLoad)}
	addDAGLegacyPluginStage(builder, manager, lifecycle.StageLoad, provided, observe)
	graph, err := builder.Compile()
	if err != nil {
		t.Fatalf("Compile() = %v", err)
	}
	executor, err := builddag.NewExecutor(1)
	if err != nil {
		t.Fatalf("NewExecutor() = %v", err)
	}
	result, err := executor.Execute(context.Background(), graph)
	if err != nil {
		t.Fatalf("Execute() = %v", err)
	}

	if got, want := calls, []string{"early", "late"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("plugin calls = %v, want %v", got, want)
	}
	if got, want := observerEvents, []string{"start", "finish"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("observer events = %v, want %v", got, want)
	}
	if !manager.HasRun(lifecycle.StageLoad) {
		t.Fatal("load stage was not marked complete")
	}
	if got, want := result.TaskCount, 4; got != want {
		t.Fatalf("task count = %d, want %d", got, want)
	}
}
