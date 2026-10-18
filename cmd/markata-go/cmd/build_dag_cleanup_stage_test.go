package cmd

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/builddag"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

type dagCleanupStageTestPlugin struct {
	name     string
	priority int
	calls    *[]string
	err      error
}

func (p *dagCleanupStageTestPlugin) Name() string { return p.name }

func (p *dagCleanupStageTestPlugin) Priority(lifecycle.Stage) int { return p.priority }

func (p *dagCleanupStageTestPlugin) Cleanup(*lifecycle.Manager) error {
	*p.calls = append(*p.calls, p.name)
	return p.err
}

type dagCleanupCriticalError struct{ err error }

func (e dagCleanupCriticalError) Error() string  { return e.err.Error() }
func (e dagCleanupCriticalError) Unwrap() error  { return e.err }
func (dagCleanupCriticalError) IsCritical() bool { return true }

func TestDAGCleanupPluginStagePreservesOrderAndObserver(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	manager.RegisterPlugins(
		&dagCleanupStageTestPlugin{name: "late", priority: lifecycle.PriorityLate, calls: &calls},
		&dagCleanupStageTestPlugin{name: "early", priority: lifecycle.PriorityEarly, calls: &calls},
	)

	observerEvents := []string{}
	observe := func(stage lifecycle.Stage, starting bool, err error) {
		if err != nil {
			t.Fatalf("observer error for %s: %v", stage, err)
		}
		if stage != lifecycle.StageCleanup {
			t.Fatalf("observer stage = %s, want %s", stage, lifecycle.StageCleanup)
		}
		if starting {
			observerEvents = append(observerEvents, "start")
			return
		}
		observerEvents = append(observerEvents, "finish")
	}

	builder := builddag.NewBuilder()
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageCleanup)}
	addDAGLegacyPluginStage(builder, manager, lifecycle.StageCleanup, provided, observe)
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
	if !manager.HasRun(lifecycle.StageCleanup) {
		t.Fatal("cleanup stage was not marked complete")
	}
	if got, want := result.TaskCount, 4; got != want {
		t.Fatalf("task count = %d, want %d", got, want)
	}
}

func TestDAGCleanupPluginStageContinuesAfterWarning(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	warningErr := errors.New("cleanup warning")
	manager.RegisterPlugins(
		&dagCleanupStageTestPlugin{name: "warn", calls: &calls, err: warningErr},
		&dagCleanupStageTestPlugin{name: "after", calls: &calls},
	)

	builder := builddag.NewBuilder()
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageCleanup)}
	addDAGLegacyPluginStage(builder, manager, lifecycle.StageCleanup, provided, nil)
	graph, err := builder.Compile()
	if err != nil {
		t.Fatalf("Compile() = %v", err)
	}
	executor, err := builddag.NewExecutor(1)
	if err != nil {
		t.Fatalf("NewExecutor() = %v", err)
	}
	if _, err := executor.Execute(context.Background(), graph); err != nil {
		t.Fatalf("Execute() = %v, want warning-only cleanup", err)
	}

	if got, want := calls, []string{"warn", "after"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("plugin calls = %v, want %v", got, want)
	}
	warnings := manager.Warnings()
	if len(warnings) != 1 || warnings[0].Critical || !errors.Is(warnings[0], warningErr) {
		t.Fatalf("warnings = %+v, want one non-critical cleanup warning", warnings)
	}
	if !manager.HasRun(lifecycle.StageCleanup) {
		t.Fatal("warning-only cleanup stage was not marked complete")
	}
}

func TestDAGCleanupPluginStageStopsOnCriticalError(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	artifactErr := errors.New("diagnostics publication failed")
	manager.RegisterPlugins(
		&dagCleanupStageTestPlugin{name: "critical", calls: &calls, err: dagCleanupCriticalError{err: artifactErr}},
		&dagCleanupStageTestPlugin{name: "after", calls: &calls},
	)

	observerEvents := []string{}
	observe := func(stage lifecycle.Stage, starting bool, err error) {
		if stage != lifecycle.StageCleanup {
			t.Fatalf("observer stage = %s, want %s", stage, lifecycle.StageCleanup)
		}
		if starting {
			observerEvents = append(observerEvents, "start")
			return
		}
		if err == nil {
			t.Fatal("cleanup failure observer event had nil error")
		}
		observerEvents = append(observerEvents, "error")
	}

	builder := builddag.NewBuilder()
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageCleanup)}
	addDAGLegacyPluginStage(builder, manager, lifecycle.StageCleanup, provided, observe)
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
		t.Fatal("Execute() = nil, want critical cleanup error")
	}
	var hookErr *lifecycle.HookError
	if !errors.As(err, &hookErr) || !hookErr.Critical || !errors.Is(err, artifactErr) {
		t.Fatalf("Execute() error = %T %v, want critical HookError", err, err)
	}
	if got, want := calls, []string{"critical"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("plugin calls = %v, want %v", got, want)
	}
	if got, want := observerEvents, []string{"start", "error"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("observer events = %v, want %v", got, want)
	}
	if manager.HasRun(lifecycle.StageCleanup) {
		t.Fatal("failed cleanup stage was marked complete")
	}
	if got := manager.Warnings(); len(got) != 0 {
		t.Fatalf("critical cleanup error was recorded as warning: %+v", got)
	}
}
