package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/internal/testbinary"
	"github.com/WaylonWalker/markata/pkg/buildlab"
)

// This test is deliberately observational. Current main may report product
// failures for linked-target invalidation; those failures are evidence, not a
// failure of this harness. A harness failure remains a test failure.
func TestBuildLab_LinkedAndFixtureMutationsCharacterizeProduct(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := filepath.Join(moduleRoot(t), "cmd", "markata-go", "cmd", "testdata", "buildlab-site")
	binary := buildTestBinary(t)
	// The focused Build Lab integration tests below retain deterministic-replay
	// coverage. Avoid an extra clean build at each of this scenario's eight
	// mutation checkpoints; each checkpoint still compares baseline, clean
	// candidate, and incremental candidate output.
	result, runErr := buildlab.RunScenario(context.Background(), buildlab.ScenarioRunConfig{
		Fixture: fixture,
		Scenario: buildlab.Scenario{ID: "cli-buildlab-fixture-mutations", Version: "1", Operations: []buildlab.Operation{
			{Type: buildlab.OpClearCache}, {Type: buildlab.OpClearOutput}, {Type: buildlab.OpBuild}, {Type: buildlab.OpBuild},
			{Type: buildlab.OpReplaceExact, Path: "content/target.md", Old: "Original target content.", New: "Updated target content."}, {Type: buildlab.OpBuild},
			{Type: buildlab.OpWriteFile, Path: "content/future-target.md", Content: fixtureFutureTarget}, {Type: buildlab.OpBuild},
			{Type: buildlab.OpDelete, Path: "content/target.md"}, {Type: buildlab.OpBuild},
			{Type: buildlab.OpWriteFile, Path: "content/target.md", Content: fixtureRecreatedTarget}, {Type: buildlab.OpBuild},
			{Type: buildlab.OpRename, Path: "content/target.md", Dest: "content/renamed.md"}, {Type: buildlab.OpBuild},
			{Type: buildlab.OpSetConfig, Path: "markata-go.toml", Key: "title", Value: "Changed by scenario"}, {Type: buildlab.OpBuild},
		}},
		Baseline:   buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Candidate:  buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Classes:    map[string]buildlab.OutputClass{".markata/diagnostics.json": buildlab.ClassVolatile, ".well-known/time": buildlab.ClassVolatile},
		GOMAXPROCS: 1,
	})
	for checkpointIndex := range result.Checkpoints {
		checkpoint := &result.Checkpoints[checkpointIndex]
		observations := []buildlab.RunObservation{checkpoint.Baseline, checkpoint.CandidateClean, checkpoint.CandidateIncremental}
		if checkpoint.CandidateDeterminism != nil {
			observations = append(observations, *checkpoint.CandidateDeterminism)
		}
		for _, observation := range observations {
			if observation.FailureClass == buildlab.FailureHarness {
				t.Fatalf("harness failure at operation %d: %+v (run error: %v)", checkpoint.OperationIndex, observation, runErr)
			}
			if observation.FailureClass == buildlab.FailureProduct {
				t.Logf("observed product failure at operation %d: error=%q stderr_sha256=%s", checkpoint.OperationIndex, observation.Error, observation.StderrSHA256)
			}
		}
	}
	if runErr != nil && !hasBuildLabProductFailure(result) {
		t.Fatalf("Build Lab run error was not represented as a product failure: %v; failure_class=%q", runErr, result.FailureClass)
	}
	if result.Verdict == buildLabPassVerdict {
		t.Logf("current main passed linked-target characterization")
	} else {
		t.Logf("current main product result: verdict=%s, failure_class=%s, checkpoints=%d", result.Verdict, result.FailureClass, len(result.Checkpoints))
		for checkpointIndex := range result.Checkpoints {
			checkpoint := &result.Checkpoints[checkpointIndex]
			for _, diagnostic := range checkpoint.Diagnostics {
				t.Logf("operation %d: %s: %s", checkpoint.OperationIndex, diagnostic.Scope, diagnostic.Message)
			}
		}
	}
	if len(result.Checkpoints) != 8 {
		t.Fatalf("checkpoints = %d, want 8", len(result.Checkpoints))
	}
}

func TestBuildLab_DeletePostPrunesGeneratedOutput(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := filepath.Join(moduleRoot(t), "cmd", "markata-go", "cmd", "testdata", "buildlab-site")
	binary := buildTestBinary(t)
	result, runErr := buildlab.RunScenario(context.Background(), buildlab.ScenarioRunConfig{
		Fixture: fixture,
		Scenario: buildlab.Scenario{ID: "cli-delete-post", Version: "1", Operations: []buildlab.Operation{
			{Type: buildlab.OpClearCache},
			{Type: buildlab.OpClearOutput},
			{Type: buildlab.OpBuild},
			{Type: buildlab.OpBuild},
			{Type: buildlab.OpDelete, Path: "content/target.md"},
			{Type: buildlab.OpBuild},
		}},
		Baseline:         buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Candidate:        buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Classes:          map[string]buildlab.OutputClass{".markata/diagnostics.json": buildlab.ClassVolatile, ".well-known/time": buildlab.ClassVolatile},
		CheckDeterminism: true,
		GOMAXPROCS:       1,
	})
	if runErr != nil {
		t.Fatalf("Build Lab run error = %v", runErr)
	}
	if result.Verdict != buildLabPassVerdict {
		t.Fatalf("Build Lab verdict = %s, want %s; diagnostics = %+v", result.Verdict, buildLabPassVerdict, result.Diagnostics)
	}
	if len(result.Checkpoints) != 3 {
		t.Fatalf("checkpoints = %d, want 3", len(result.Checkpoints))
	}
}

func TestBuildLab_LinkedTargetEditMatchesCleanPageAndFeeds(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := filepath.Join(moduleRoot(t), "cmd", "markata-go", "cmd", "testdata", "buildlab-site")
	binary := buildTestBinary(t)
	result, runErr := buildlab.RunScenario(context.Background(), buildlab.ScenarioRunConfig{
		Fixture: fixture,
		Scenario: buildlab.Scenario{ID: "cli-linked-target-edit", Version: "1", Operations: []buildlab.Operation{
			{Type: buildlab.OpClearCache},
			{Type: buildlab.OpClearOutput},
			{Type: buildlab.OpBuild},
			{Type: buildlab.OpReplaceExact, Path: "content/target.md", Old: "Original target content.", New: "Updated target content."},
			{Type: buildlab.OpBuild},
		}},
		Baseline:         buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Candidate:        buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Classes:          map[string]buildlab.OutputClass{".markata/diagnostics.json": buildlab.ClassVolatile, ".well-known/time": buildlab.ClassVolatile},
		CheckDeterminism: true,
		GOMAXPROCS:       1,
	})
	if runErr != nil {
		t.Fatalf("Build Lab run error = %v", runErr)
	}
	if result.Verdict != buildLabPassVerdict {
		t.Fatalf("Build Lab verdict = %s, want %s; diagnostics = %+v", result.Verdict, buildLabPassVerdict, result.Diagnostics)
	}
	if len(result.Checkpoints) != 2 {
		t.Fatalf("checkpoints = %d, want 2", len(result.Checkpoints))
	}
	checkpoint := result.Checkpoints[1]
	if !checkpoint.Correctness.IncrementalApplicable || !checkpoint.Correctness.IncrementalEqual {
		t.Fatalf("linked-target incremental comparison = %+v", checkpoint.Correctness)
	}

	assertLinkedTargetIncrementalScope(t, fixture, binary)
}

func assertLinkedTargetIncrementalScope(t *testing.T, fixture, binary string) {
	t.Helper()
	workspace, err := buildlab.NewWorkspace(fixture, t.TempDir())
	if err != nil {
		t.Fatalf("create Build Lab workspace: %v", err)
	}
	t.Cleanup(func() {
		if err := workspace.Remove(); err != nil {
			t.Errorf("remove Build Lab workspace: %v", err)
		}
	})
	environment, err := workspace.Environment([]string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}, 1)
	if err != nil {
		t.Fatalf("create Build Lab environment: %v", err)
	}
	outputDir := filepath.Join(workspace.SiteDir, "output")
	runBuild := func() []byte {
		t.Helper()
		//nolint:gosec // The test controls the freshly built binary and isolated workspace paths.
		command := exec.Command(
			binary,
			"build", "-c", "markata-go.toml",
			"--merge-config", workspace.IsolationConfig,
			"--output", outputDir,
			"--site-dir", workspace.SiteDir,
		)
		command.Dir = workspace.SiteDir
		command.Env = environment
		output, err := command.CombinedOutput()
		if err != nil {
			t.Fatalf("run fixture build: %v\n%s", err, output)
		}
		return output
	}

	runBuild()
	targetPath := filepath.Join(workspace.SiteDir, "content", "target.md")
	target, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("read target fixture: %v", err)
	}
	oldContent := []byte("Original target content.")
	newContent := []byte("Updated target content.")
	if bytes.Count(target, oldContent) != 1 {
		t.Fatalf("target fixture occurrence count = %d, want 1", bytes.Count(target, oldContent))
	}
	if err := os.WriteFile(targetPath, bytes.Replace(target, oldContent, newContent, 1), 0o600); err != nil {
		t.Fatalf("update target fixture: %v", err)
	}

	incrementalLog := runBuild()
	if !bytes.Contains(incrementalLog, []byte("Invalidated 1 changed posts and 1 dependent posts")) {
		t.Fatalf("incremental build did not report dependency-scoped invalidation:\n%s", incrementalLog)
	}
	if !bytes.Contains(incrementalLog, []byte("Phase 1a classify: 1 cacheable, 2 need render")) {
		t.Fatalf("incremental build did not keep the unrelated post cacheable:\n%s", incrementalLog)
	}
	if !bytes.Contains(incrementalLog, []byte("Incremental build: 1 skipped, 2 rebuilt")) {
		t.Fatalf("incremental build did not skip only the unrelated post:\n%s", incrementalLog)
	}

	for _, relativePath := range []string{"source/index.html", "archive/atom.xml", "archive/feed.json"} {
		generated, err := os.ReadFile(filepath.Join(outputDir, relativePath))
		if err != nil {
			t.Fatalf("read generated %s: %v", relativePath, err)
		}
		if !bytes.Contains(generated, newContent) || bytes.Contains(generated, oldContent) {
			t.Fatalf("generated %s did not contain only updated embedded content", relativePath)
		}
	}
}

func TestBuildLab_ConfigChangeRebuildsFeedsListing(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := filepath.Join(moduleRoot(t), "cmd", "markata-go", "cmd", "testdata", "buildlab-site")
	binary := buildTestBinary(t)
	result, runErr := buildlab.RunScenario(context.Background(), buildlab.ScenarioRunConfig{
		Fixture: fixture,
		Scenario: buildlab.Scenario{ID: "cli-config-change", Version: "1", Operations: []buildlab.Operation{
			{Type: buildlab.OpClearCache},
			{Type: buildlab.OpClearOutput},
			{Type: buildlab.OpBuild},
			{Type: buildlab.OpSetConfig, Path: "markata-go.toml", Key: "title", Value: "Changed by scenario"},
			{Type: buildlab.OpBuild},
		}},
		Baseline:         buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Candidate:        buildlab.BuildCommand{Binary: binary, Args: []string{"build", "-c", "markata-go.toml"}, OutputDir: "output", Timeout: 5 * time.Minute, Env: []string{"MARKATA_GO_ENCRYPTION_ENABLED=false"}},
		Classes:          map[string]buildlab.OutputClass{".markata/diagnostics.json": buildlab.ClassVolatile, ".well-known/time": buildlab.ClassVolatile},
		CheckDeterminism: true,
		GOMAXPROCS:       1,
	})
	if runErr != nil {
		t.Fatalf("Build Lab run error = %v", runErr)
	}
	if result.Verdict != buildLabPassVerdict {
		for _, checkpoint := range result.Checkpoints {
			t.Logf("Build Lab operation %d (%s): correctness = %+v; diagnostics = %+v", checkpoint.OperationIndex, checkpoint.Operation, checkpoint.Correctness, checkpoint.Diagnostics)
		}
		t.Fatalf("Build Lab verdict = %s, want %s; diagnostics = %+v", result.Verdict, buildLabPassVerdict, result.Diagnostics)
	}
	if len(result.Checkpoints) != 2 {
		t.Fatalf("checkpoints = %d, want 2", len(result.Checkpoints))
	}
}

func TestBuildLabCLI_EmitsMachineReadableResultAndUsesOrdinaryBuild(t *testing.T) {
	requireLinuxBuildLab(t)
	fixture := t.TempDir()
	if err := os.WriteFile(filepath.Join(fixture, "markata-go.toml"), []byte("[markata-go]\noutput_dir = \"output\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(t.TempDir(), "build.sh")
	const scriptText = `#!/bin/sh
out=""
while [ "$#" -gt 0 ]; do
  case "$1" in
    --dag*) exit 91 ;;
    --output|-o) out="$2"; shift 2 ;;
    --output=*) out="${1#--output=}"; shift ;;
    *) shift ;;
  esac
done
mkdir -p "$out"
`
	if err := os.WriteFile(script, []byte(scriptText), 0o700); err != nil { //nolint:gosec // The test fixture must be executable as a build command.
		t.Fatal(err)
	}
	scenario := filepath.Join(t.TempDir(), "scenario.json")
	if err := os.WriteFile(scenario, []byte(`{"id":"cli-contract","version":"1","seed":42,"operations":[{"type":"build"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	resultPath := filepath.Join(t.TempDir(), "result.json")
	command := exec.Command(buildTestBinary(t), "buildlab", "run", "--fixture", fixture, "--baseline", script, "--candidate", script, "--scenario", scenario, "--check-determinism=false", "--result", resultPath) //nolint:gosec // Test-controlled fixture and executable paths are intentional.
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err != nil {
		t.Fatalf("CLI failed: %v\nstderr=%s\nstdout=%s", err, stderr.String(), stdout.String())
	}
	var got, fileResult buildlab.Result
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &got); err != nil {
		t.Fatalf("stdout is not JSON: %v", err)
	}
	data, err := os.ReadFile(resultPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &fileResult); err != nil {
		t.Fatal(err)
	}
	if got.Verdict != buildLabPassVerdict || got.Scenario.Seed != 42 || fileResult.Verdict != got.Verdict {
		t.Fatalf("stdout/file result = %+v / %+v", got, fileResult)
	}
}

func TestBuildLabCLI_EmitsSetupFailureAndFinalizesCPUProfile(t *testing.T) {
	requireLinuxBuildLab(t)
	binary := buildTestBinary(t)
	profile := filepath.Join(t.TempDir(), "buildlab.prof")
	missingFixture := filepath.Join(t.TempDir(), "missing-fixture")
	command := exec.Command(binary, "--cpuprofile", profile, "buildlab", "run", "--fixture", missingFixture)
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	if err := command.Run(); err == nil {
		t.Fatal("Build Lab setup failure unexpectedly succeeded")
	}
	var result buildlab.Result
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &result); err != nil {
		t.Fatalf("setup failure stdout is not JSON: %v\nstdout=%s\nstderr=%s", err, stdout.String(), stderr.String())
	}
	if result.Verdict != "fail" || result.FailureClass != buildlab.FailureHarness || len(result.Diagnostics) != 1 {
		t.Fatalf("setup failure result = %+v", result)
	}
	if !bytes.Contains(stderr.Bytes(), []byte("Build Lab harness (setup)")) {
		t.Fatalf("setup diagnostics missing from stderr: %s", stderr.String())
	}
	info, err := os.Stat(profile)
	if err != nil {
		t.Fatalf("CPU profile was not written: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("CPU profile is empty after failed Build Lab command")
	}
}

const fixtureFutureTarget = "---\ntitle: Future target\ndate: 2026-08-03\npublished: true\n---\n\n# Future target\n\nCreated after the initial build.\n"
const fixtureRecreatedTarget = "---\ntitle: Target post\ndate: 2026-08-01\npublished: true\n---\n\n# Target post\n\nRecreated target content.\n"

func hasBuildLabProductFailure(result buildlab.Result) bool {
	for checkpointIndex := range result.Checkpoints {
		checkpoint := &result.Checkpoints[checkpointIndex]
		observations := []buildlab.RunObservation{checkpoint.Baseline, checkpoint.CandidateClean, checkpoint.CandidateIncremental}
		if checkpoint.CandidateDeterminism != nil {
			observations = append(observations, *checkpoint.CandidateDeterminism)
		}
		for _, observation := range observations {
			if observation.FailureClass == buildlab.FailureProduct {
				return true
			}
		}
	}
	return false
}

func buildTestBinary(t *testing.T) string {
	t.Helper()
	binary, err := testbinary.Resolve(moduleRoot(t), t.TempDir())
	if err != nil {
		t.Fatalf("resolve markata-go test binary: %v", err)
	}
	return binary
}
func requireLinuxBuildLab(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skip("Build Lab process groups are implemented on Linux")
	}
}
func moduleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not locate go.mod")
		}
		dir = parent
	}
}
