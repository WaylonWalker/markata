package lifecycle

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
)

func TestManagerTemplateCacheSurvivesRunToAndClearsOnReset(t *testing.T) {
	manager := NewManager()
	stats := diagnostics.TemplateCacheStats{Classified: 1, RenderRequired: 1, RenderSucceeded: 1}
	// Advancing an already-started build must not reset its render observations.
	if err := manager.RunTo(StageRender); err != nil {
		t.Fatal(err)
	}
	manager.ContentLedger().SetTemplateCache(&stats)
	if err := manager.RunTo(StageWrite); err != nil {
		t.Fatal(err)
	}
	if got := manager.ContentDiagnostics().TemplateCache; got == nil || *got != stats {
		t.Fatalf("RunTo discarded observations: %+v", got)
	}
	other := NewManager()
	if other.ContentDiagnostics().TemplateCache != nil {
		t.Fatal("another manager inherited observations")
	}
	manager.Reset()
	if manager.ContentDiagnostics().TemplateCache != nil {
		t.Fatal("Reset retained observations")
	}
}
