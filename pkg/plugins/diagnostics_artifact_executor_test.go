package plugins

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

func TestDiagnosticsArtifactPersistsBuildExecutor(t *testing.T) {
	manager := lifecycle.NewManager()
	config := manager.Config()
	config.OutputDir = t.TempDir()
	config.ContentDir = t.TempDir()
	if err := manager.SetBuildExecutor(lifecycle.BuildExecutorDAG); err != nil {
		t.Fatalf("SetBuildExecutor() = %v", err)
	}

	for _, stage := range []lifecycle.Stage{
		lifecycle.StageConfigure,
		lifecycle.StageValidate,
		lifecycle.StageGlob,
		lifecycle.StageLoad,
		lifecycle.StageTransform,
		lifecycle.StageRender,
		lifecycle.StageCollect,
		lifecycle.StageWrite,
		lifecycle.StageCleanup,
	} {
		if err := manager.MarkStageComplete(stage); err != nil {
			t.Fatalf("MarkStageComplete(%s) = %v", stage, err)
		}
	}

	plugin := NewDiagnosticsArtifactPlugin()
	if err := plugin.Cleanup(manager); err != nil {
		t.Fatalf("Cleanup() = %v", err)
	}
	data, err := os.ReadFile(filepath.Join(config.OutputDir, diagnostics.DefaultArtifactPath))
	if err != nil {
		t.Fatalf("read diagnostics artifact: %v", err)
	}
	artifact, err := diagnostics.ParseArtifact(data)
	if err != nil {
		t.Fatalf("ParseArtifact() = %v", err)
	}
	if artifact.Executor != diagnostics.ArtifactExecutorDAG {
		t.Fatalf("executor = %q, want %q", artifact.Executor, diagnostics.ArtifactExecutorDAG)
	}
}
