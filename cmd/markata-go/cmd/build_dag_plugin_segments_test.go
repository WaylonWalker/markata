package cmd

import (
	"context"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/builddag"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

type dagBarrierTransformPlugin struct {
	name string
	run  func(*lifecycle.Manager) error
}

func (p *dagBarrierTransformPlugin) Name() string { return p.name }

func (p *dagBarrierTransformPlugin) Transform(m *lifecycle.Manager) error {
	if p.run == nil {
		return nil
	}
	return p.run(m)
}

func TestDAGPluginExpanderRunsAfterPredecessorCompletes(t *testing.T) {
	manager := lifecycle.NewManager()
	state := "before"
	manager.RegisterPlugin(&dagBarrierTransformPlugin{
		name: "first",
		run: func(*lifecycle.Manager) error {
			state = "after-first"
			return nil
		},
	})
	manager.RegisterPlugin(&dagBarrierTransformPlugin{name: "second"})

	executor, err := builddag.NewExecutor(1)
	if err != nil {
		t.Fatal(err)
	}
	seenState := ""
	expander := func(_ *lifecycle.Manager, _ lifecycle.Stage, task builddag.TaskSpec) []builddag.TaskSpec {
		if strings.HasSuffix(string(task.ID), ".second") {
			seenState = state
		}
		return []builddag.TaskSpec{task}
	}
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageTransform)}
	if _, _, err := executeDAGPluginStageSegments(executor, manager, lifecycle.StageTransform, nil, provided, nil, expander); err != nil {
		t.Fatalf("executeDAGPluginStageSegments() = %v", err)
	}
	if seenState != "after-first" {
		t.Fatalf("second plugin expanded with state %q, want predecessor mutation", seenState)
	}
}

func TestDAGPluginStageSnapshotsPluginListOnce(t *testing.T) {
	manager := lifecycle.NewManager()
	lateRan := false
	late := &dagBarrierTransformPlugin{
		name: "late",
		run: func(*lifecycle.Manager) error {
			lateRan = true
			return nil
		},
	}
	manager.RegisterPlugin(&dagBarrierTransformPlugin{
		name: "registrar",
		run: func(m *lifecycle.Manager) error {
			m.RegisterPlugin(late)
			return nil
		},
	})

	executor, err := builddag.NewExecutor(1)
	if err != nil {
		t.Fatal(err)
	}
	provided := builddag.ArtifactID{Kind: "lifecycle-stage", Key: string(lifecycle.StageTransform)}
	if _, _, err := executeDAGPluginStageSegments(executor, manager, lifecycle.StageTransform, nil, provided, nil, defaultDAGPluginExpander); err != nil {
		t.Fatalf("executeDAGPluginStageSegments() = %v", err)
	}
	if lateRan {
		t.Fatal("plugin registered during Transform joined the already-snapshotted stage")
	}
}

func TestDAGPluginExpansionMustPreserveLegacyBoundary(t *testing.T) {
	required := builddag.ArtifactID{Kind: "legacy-completion", Key: "previous"}
	provided := builddag.ArtifactID{Kind: "legacy-completion", Key: "current"}
	legacyTask := builddag.TaskSpec{
		ID:       "legacy.transform.000.test",
		Requires: []builddag.ArtifactID{required},
		Provides: []builddag.ArtifactID{provided},
		Func:     func(context.Context) error { return nil },
	}

	if err := validateDAGPluginExpansionBoundary(legacyTask, []builddag.TaskSpec{{
		ID:       "replacement",
		Requires: []builddag.ArtifactID{required},
		Provides: []builddag.ArtifactID{provided},
		Func:     func(context.Context) error { return nil },
	}}); err != nil {
		t.Fatalf("valid expansion rejected: %v", err)
	}
	if err := validateDAGPluginExpansionBoundary(legacyTask, []builddag.TaskSpec{{
		ID:       "replacement",
		Requires: []builddag.ArtifactID{required},
		Func:     func(context.Context) error { return nil },
	}}); err == nil {
		t.Fatal("expansion that dropped the legacy completion artifact was accepted")
	}
}
