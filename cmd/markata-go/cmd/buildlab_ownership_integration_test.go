package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/buildlab"
)

// These are strict cold/warm equivalence contracts for known incremental
// ownership and invalidation defects. Keep them isolated from the shared
// fixture so source mutations cannot influence each other.
func TestBuildLabPublishedOutputOwnershipParity(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := t.TempDir()
	writeBuildLabFixtureFile(t, fixture, "markata-go.toml", `[markata-go]
title = "HTML ownership regression"
license = false
output_dir = "output"

[markata-go.glob]
patterns = ["content/**/*.md"]
use_gitignore = false

[markata-go.search]
enabled = false

[markata-go.tailwind]
include = false
build = false
`)
	writeBuildLabFixtureFile(t, fixture, "content/deck.md", "---\ntitle: Generated deck\nslug: deck\npublished: true\n---\n\nGenerated page body.\n")
	writeBuildLabFixtureFile(t, fixture, "static/deck/index.html", "<!doctype html><title>Static deck</title>")
	runBuildLabOwnershipScenario(t, fixture, "static-generated-html-ownership", []buildlab.Operation{
		{Type: buildlab.OpClearCache},
		{Type: buildlab.OpClearOutput},
		{Type: buildlab.OpBuild},
		{Type: buildlab.OpBuild},
	})
}

func TestBuildLabNoAliasSlugRenameParity(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := filepath.Join(moduleRoot(t), "cmd", "markata-go", "cmd", "testdata", "buildlab-site")
	runBuildLabOwnershipScenario(t, fixture, "rename-target-without-alias", []buildlab.Operation{
		{Type: buildlab.OpClearCache},
		{Type: buildlab.OpClearOutput},
		{Type: buildlab.OpBuild},
		{Type: buildlab.OpBuild},
		{Type: buildlab.OpReplaceExact, Path: "content/target.md", Old: "title: Target post", New: "title: Target post\nslug: changed-target"},
		{Type: buildlab.OpBuild},
	})
}

func TestBuildLabGardenDeletionParity(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := t.TempDir()
	writeBuildLabFixtureFile(t, fixture, "markata-go.toml", `[markata-go]
title = "Garden deletion regression"
license = false
output_dir = "output"

[markata-go.glob]
patterns = ["content/**/*.md"]
use_gitignore = false

[markata-go.garden]
enabled = true

[markata-go.search]
enabled = false

[markata-go.tailwind]
include = false
build = false
`)
	writeBuildLabFixtureFile(t, fixture, "content/probe.md", "---\ntitle: Probe\npublished: true\n---\n\nLink to [[neighbor]].\n")
	writeBuildLabFixtureFile(t, fixture, "content/neighbor.md", "---\ntitle: Neighbor\npublished: true\n---\n\nNeighbor.\n")
	runBuildLabOwnershipScenario(t, fixture, "garden-first-deletion", []buildlab.Operation{
		{Type: buildlab.OpClearCache},
		{Type: buildlab.OpClearOutput},
		{Type: buildlab.OpBuild},
		{Type: buildlab.OpDelete, Path: "content/probe.md"},
		{Type: buildlab.OpBuild},
	})
}

func writeBuildLabFixtureFile(t *testing.T, root, name, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func runBuildLabOwnershipScenario(t *testing.T, fixture, name string, operations []buildlab.Operation) {
	t.Helper()
	result, err := buildlab.RunScenario(context.Background(), buildlab.ScenarioRunConfig{
		Fixture: fixture,
		Scenario: buildlab.Scenario{
			ID:         name,
			Version:    "1",
			Operations: operations,
		},
		Baseline: buildlab.BuildCommand{
			Binary:    buildTestBinary(t),
			Args:      []string{"build", "-c", "markata-go.toml"},
			OutputDir: "output",
			Timeout:   5 * time.Minute,
			Env:       []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"},
		},
		Candidate: buildlab.BuildCommand{
			Binary:    buildTestBinary(t),
			Args:      []string{"build", "-c", "markata-go.toml"},
			OutputDir: "output",
			Timeout:   5 * time.Minute,
			Env:       []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"},
		},
		Classes:          map[string]buildlab.OutputClass{".markata/diagnostics.json": buildlab.ClassVolatile, ".well-known/time": buildlab.ClassVolatile},
		CheckDeterminism: true,
		GOMAXPROCS:       1,
	})
	if err != nil {
		t.Fatalf("Build Lab %s failed: %v; diagnostics=%+v", name, err, result.Diagnostics)
	}
	if result.Verdict != buildLabPassVerdict {
		for i := range result.Checkpoints {
			checkpoint := &result.Checkpoints[i]
			t.Logf("%s operation %d (%s): incremental=%v applicable=%v determinism=%v diff=%+v diagnostics=%+v",
				name, checkpoint.OperationIndex, checkpoint.Operation, checkpoint.Correctness.IncrementalEqual,
				checkpoint.Correctness.IncrementalApplicable, checkpoint.Correctness.DeterministicEqual,
				checkpoint.Correctness.IncrementalDiff, checkpoint.Diagnostics)
		}
		t.Fatalf("Build Lab %s verdict=%q, diagnostics=%+v", name, result.Verdict, result.Diagnostics)
	}
	for i := range result.Checkpoints {
		checkpoint := &result.Checkpoints[i]
		if checkpoint.Correctness.IncrementalApplicable && !checkpoint.Correctness.IncrementalEqual {
			t.Fatalf("%s operation %d differs from clean output: %+v", name, checkpoint.OperationIndex, checkpoint.Correctness)
		}
	}
}
