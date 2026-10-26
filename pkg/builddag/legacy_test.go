package builddag

import (
	"context"
	"reflect"
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

type legacyTaskTestPlugin struct {
	name     string
	priority int
	calls    *[]string
}

func (p *legacyTaskTestPlugin) Name() string { return p.name }

func (p *legacyTaskTestPlugin) Priority(lifecycle.Stage) int { return p.priority }

func (p *legacyTaskTestPlugin) Transform(*lifecycle.Manager) error {
	*p.calls = append(*p.calls, p.name)
	return nil
}

type legacyTaskNameOnlyPlugin struct{ name string }

func (p *legacyTaskNameOnlyPlugin) Name() string { return p.name }

func TestLegacyTasksPreservePriorityOrderAndDependencies(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	manager.RegisterPlugins(
		&legacyTaskTestPlugin{name: "default-a", calls: &calls},
		&legacyTaskTestPlugin{name: "early", priority: lifecycle.PriorityEarly, calls: &calls},
		&legacyTaskNameOnlyPlugin{name: "not-transform"},
		&legacyTaskTestPlugin{name: "default-b", calls: &calls},
	)

	start := ArtifactID{Kind: "stage", Key: "load"}
	tasks := LegacyTasks(manager, lifecycle.StageTransform, []ArtifactID{start})
	if len(tasks) != 3 {
		t.Fatalf("LegacyTasks() returned %d tasks, want 3", len(tasks))
	}

	wantIDs := []TaskID{
		"legacy.transform.000.early",
		"legacy.transform.001.default-a",
		"legacy.transform.002.default-b",
	}
	gotIDs := make([]TaskID, 0, len(tasks))
	for i := range tasks {
		gotIDs = append(gotIDs, tasks[i].ID)
		if !tasks[i].Exclusive {
			t.Fatalf("task %q is not exclusive", tasks[i].ID)
		}
		if tasks[i].ParallelSafe {
			t.Fatalf("task %q unexpectedly marked parallel-safe", tasks[i].ID)
		}
	}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Fatalf("task IDs = %v, want %v", gotIDs, wantIDs)
	}
	if !reflect.DeepEqual(tasks[0].Requires, []ArtifactID{start}) {
		t.Fatalf("first task requires = %v, want %v", tasks[0].Requires, []ArtifactID{start})
	}
	for i := 1; i < len(tasks); i++ {
		if !reflect.DeepEqual(tasks[i].Requires, tasks[i-1].Provides) {
			t.Fatalf("task %q requires = %v, want previous provides %v", tasks[i].ID, tasks[i].Requires, tasks[i-1].Provides)
		}
	}
}

func TestLegacyTasksExecuteThroughSerialGraph(t *testing.T) {
	manager := lifecycle.NewManager()
	calls := []string{}
	manager.RegisterPlugins(
		&legacyTaskTestPlugin{name: "second", priority: lifecycle.PriorityDefault, calls: &calls},
		&legacyTaskTestPlugin{name: "first", priority: lifecycle.PriorityEarly, calls: &calls},
	)

	start := ArtifactID{Kind: "stage", Key: "load"}
	builder := NewBuilder()
	builder.AddExternal(start)
	for _, task := range LegacyTasks(manager, lifecycle.StageTransform, []ArtifactID{start}) {
		builder.AddTask(task)
	}
	graph, err := builder.Compile()
	if err != nil {
		t.Fatalf("Compile() = %v", err)
	}
	executor, err := NewExecutor(1)
	if err != nil {
		t.Fatalf("NewExecutor() = %v", err)
	}
	if _, err := executor.Execute(context.Background(), graph); err != nil {
		t.Fatalf("Execute() = %v", err)
	}

	want := []string{"first", "second"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("plugin calls = %v, want %v", calls, want)
	}
}

func TestLegacyTasksReturnsNilForNilManager(t *testing.T) {
	if tasks := LegacyTasks(nil, lifecycle.StageTransform, nil); tasks != nil {
		t.Fatalf("LegacyTasks(nil) = %v, want nil", tasks)
	}
}
