package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

func TestBenchmarkJSONTemplateCacheLegacyAndDAG(t *testing.T) {
	for _, executor := range []lifecycle.BuildExecutor{lifecycle.BuildExecutorLegacy, lifecycle.BuildExecutorDAG} {
		t.Run(string(executor), func(t *testing.T) {
			manager := lifecycle.NewManager()
			dir := t.TempDir()
			if err := os.WriteFile(filepath.Join(dir, "post.html"), []byte(`PAGE:{{ body | safe }}`), 0o600); err != nil {
				t.Fatal(err)
			}
			manager.Config().Extra["templates_dir"] = dir
			manager.RegisterPlugins(plugins.NewTemplatesPlugin())
			manager.AddPost(&models.Post{Path: "post.md", InputHash: "hash", Template: "post.html"})
			manager.AddPost(&models.Post{Path: "skip.md", Skip: true})
			var result *BuildResult
			if executor == lifecycle.BuildExecutorDAG {
				var err error
				result, err = runDAGBuildObserved(manager, nil)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				if err := manager.Run(); err != nil {
					t.Fatal(err)
				}
				result = &BuildResult{Executor: manager.BuildExecutor(), Content: manager.ContentDiagnostics()}
			}
			var output bytes.Buffer
			if err := writeBenchmarkJSON(&output, result); err != nil {
				t.Fatal(err)
			}
			var decoded struct {
				Content diagnostics.ContentLedgerSnapshot `json:"content"`
			}
			if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
				t.Fatal(err)
			}
			want := diagnostics.TemplateCacheStats{
				Classified: 1, Skipped: 1, RenderRequired: 1, RenderSucceeded: 1,
				MissReasons: diagnostics.TemplateCacheMissReasons{CacheUnavailable: 1},
			}
			if got := decoded.Content.TemplateCache; got == nil || *got != want {
				t.Fatalf("expanded content schema missing template stats: %+v", got)
			}
		})
	}
}
