package cmd

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

type dagLoadRegistrarPlugin struct {
	transform lifecycle.Plugin
}

func (p *dagLoadRegistrarPlugin) Name() string { return "dag_load_registrar" }

func (p *dagLoadRegistrarPlugin) Load(m *lifecycle.Manager) error {
	m.RegisterPlugin(p.transform)
	return nil
}

type dagLateTransformPlugin struct {
	ran *bool
}

func (p *dagLateTransformPlugin) Name() string { return "dag_late_transform" }

func (p *dagLateTransformPlugin) Transform(_ *lifecycle.Manager) error {
	*p.ran = true
	return nil
}

func TestDAGCompilesTransformAfterLoadMaterializesState(t *testing.T) {
	manager := lifecycle.NewManager()
	lateTransformRan := false
	manager.RegisterPlugin(&dagLoadRegistrarPlugin{
		transform: &dagLateTransformPlugin{ran: &lateTransformRan},
	})

	if _, err := runDAGBuildObserved(manager, nil); err != nil {
		t.Fatalf("runDAGBuildObserved() = %v", err)
	}
	if !lateTransformRan {
		t.Fatal("transform plugin registered during Load was absent from the later Transform graph")
	}
}

func TestDAGStagedCompilationPreservesObserverBoundaries(t *testing.T) {
	manager := lifecycle.NewManager()
	starts := make(map[lifecycle.Stage]int)
	finishes := make(map[lifecycle.Stage]int)

	_, err := runDAGBuildObserved(manager, func(stage lifecycle.Stage, starting bool, stageErr error) {
		if stageErr != nil {
			t.Fatalf("observer stage %s error = %v", stage, stageErr)
		}
		if starting {
			starts[stage]++
			return
		}
		finishes[stage]++
	})
	if err != nil {
		t.Fatalf("runDAGBuildObserved() = %v", err)
	}

	for _, stage := range dagLifecycleStages {
		if starts[stage] != 1 || finishes[stage] != 1 {
			t.Fatalf("stage %s observer counts = start:%d finish:%d, want 1/1", stage, starts[stage], finishes[stage])
		}
	}
}
