package cmd

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/builddag"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

type dagWriteStageTestPlugin struct {
	name          string
	priority      int
	calls         *[]string
	err           error
	criticalWrite bool
}

func (p *dagWriteStageTestPlugin) Name() string { return p.name }

func (p *dagWriteStageTestPlugin) Priority(lifecycle.Stage) int { return p.priority }

func (p *dagWriteStageTestPlugin) Write(*lifecycle.Manager) error {
	*p.calls = append(*p.calls, p.name)
	return p.err
}

func (p *dagWriteStageTestPlugin) CriticalStageErrors(stage lifecycle.Stage) bool {
	return p.criticalWrite && stage == lifecycle.StageWrite
}

func TestDAGWritePluginStagePreservesOrderAndObserver(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	manager.RegisterPlugins(
		&dagWriteStageTestPlugin{name: "late", priority: lifecycle.PriorityLate, calls: &calls},
		&dagWriteStageTestPlugin{name: "early", priority: lifecycle.PriorityEarly, calls: &calls},
	)

	observerEvents := []string{}
	observe := func(stage lifecycle.Stage, starting bool, err error) {
		if err != nil {
			t.Fatalf("observer error for %s: %v", stage, err)
		}
		if stage != lifecycle.StageWrite {
			t.Fatalf("observer stage = %s, want %s", stage, lifecycle.StageWrite)
		}
		if starting {
			observerEvents = append(observerEvents, "start")
			return
		}
		observerEvents = append(observerEvents, "finish")
	}

	builder := builddag.NewBuilder()
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageWrite)}
	addDAGLegacyPluginStage(builder, manager, lifecycle.StageWrite, provided, observe)
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
	if !manager.HasRun(lifecycle.StageWrite) {
		t.Fatal("write stage was not marked complete")
	}
	if got, want := result.TaskCount, 4; got != want {
		t.Fatalf("task count = %d, want %d", got, want)
	}
}

func TestDAGWritePluginStageStopsOnPluginDeclaredCriticalError(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	publicationErr := errors.New("publication failed")
	manager.RegisterPlugins(
		&dagWriteStageTestPlugin{name: "critical", calls: &calls, err: publicationErr, criticalWrite: true},
		&dagWriteStageTestPlugin{name: "after", calls: &calls},
	)

	observerEvents := []string{}
	observe := func(stage lifecycle.Stage, starting bool, err error) {
		if stage != lifecycle.StageWrite {
			t.Fatalf("observer stage = %s, want %s", stage, lifecycle.StageWrite)
		}
		if starting {
			observerEvents = append(observerEvents, "start")
			return
		}
		if err == nil {
			t.Fatal("write failure observer event had nil error")
		}
		observerEvents = append(observerEvents, "error")
	}

	builder := builddag.NewBuilder()
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageWrite)}
	addDAGLegacyPluginStage(builder, manager, lifecycle.StageWrite, provided, observe)
	graph, err := builder.Compile()
	if err != nil {
		t.Fatalf("Compile() = %v", err)
	}
	executor, err := builddag.NewExecutor(1)
	if err != nil {
		t.Fatalf("NewExecutor() = %v", err)
	}
	_, err = executor.Execute(context.Background(), graph)
	if err == nil {
		t.Fatal("Execute() = nil, want critical write error")
	}
	var hookErr *lifecycle.HookError
	if !errors.As(err, &hookErr) || !hookErr.Critical || !errors.Is(err, publicationErr) {
		t.Fatalf("Execute() error = %T %v, want plugin-declared critical HookError", err, err)
	}
	if got, want := calls, []string{"critical"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("plugin calls = %v, want %v", got, want)
	}
	if got, want := observerEvents, []string{"start", "error"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("observer events = %v, want %v", got, want)
	}
	if manager.HasRun(lifecycle.StageWrite) {
		t.Fatal("failed write stage was marked complete")
	}
	if got := manager.Warnings(); len(got) != 0 {
		t.Fatalf("critical write error was recorded as warning: %+v", got)
	}
}
