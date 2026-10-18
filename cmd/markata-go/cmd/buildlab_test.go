package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/buildlab"
)

func TestBuildLabCLI_OnlyExposesReviewedDAGOptIn(t *testing.T) {
	if buildLabRunCmd.Flags().Lookup("candidate-dag") != nil || buildLabRunCmd.Flags().Lookup("candidate-dag-random-ready") != nil {
		t.Fatal("Build Lab exposes scheduler-specific flags")
	}

	dagFlag := buildCmd.Flags().Lookup("dag")
	if dagFlag == nil {
		t.Fatal("ordinary build does not expose the reviewed --dag opt-in")
	}
	if dagFlag.DefValue != envValueDisabled {
		t.Fatalf("build --dag default = %q, want false", dagFlag.DefValue)
	}
	if buildCmd.Flags().Lookup("dag-seed") != nil || buildCmd.Flags().Lookup("dag-random-ready") != nil {
		t.Fatal("ordinary build exposes legacy scheduler tuning flags")
	}
}

func TestIsBuildLabCommandSkipsSiteActivationOnlyForBuildLabTree(t *testing.T) {
	if !isBuildLabCommand(buildLabRunCmd) {
		t.Fatal("Build Lab command was not recognized")
	}
	if isBuildLabCommand(buildCmd) {
		t.Fatal("ordinary build was recognized as Build Lab")
	}
}

func TestParseBuildLabVolatileNormalizesPlatformSeparators(t *testing.T) {
	path := filepath.Join("nested", "time.txt")
	got := parseBuildLabVolatile(path)
	if got[filepath.ToSlash(path)] != buildlab.ClassVolatile || len(got) != 1 {
		t.Fatalf("volatile paths = %v", got)
	}
}

func TestBuildLabRelativePathAcceptsAliasedAbsolutePath(t *testing.T) {
	actualParent := t.TempDir()
	fixture := filepath.Join(actualParent, "fixture")
	if err := os.MkdirAll(fixture, 0o755); err != nil {
		t.Fatal(err)
	}
	aliasParent := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(actualParent, aliasParent); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	aliasedFixture := filepath.Join(aliasParent, "fixture")
	configPath := filepath.Join(fixture, "markata-go.toml")
	if err := os.WriteFile(configPath, []byte("[markata-go]\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := buildLabRelativePath(aliasedFixture, filepath.Join(aliasParent, "fixture", "markata-go.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if got != "markata-go.toml" {
		t.Fatalf("relative path = %q, want markata-go.toml", got)
	}
	got, err = buildLabRelativePath(aliasedFixture, filepath.Join(aliasParent, "fixture", "output"))
	if err != nil {
		t.Fatal(err)
	}
	if got != "output" {
		t.Fatalf("relative missing path = %q, want output", got)
	}
}

func TestBuildLabRelativePathRejectsDanglingSymlink(t *testing.T) {
	fixture := t.TempDir()
	dangling := filepath.Join(fixture, "link")
	if err := os.Symlink(filepath.Join(t.TempDir(), "missing"), dangling); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	_, err := buildLabRelativePath(fixture, filepath.Join(dangling, "output"))
	if err == nil || !strings.Contains(err.Error(), "unresolved symlink") {
		t.Fatalf("dangling symlink path was accepted: %v", err)
	}
}

func TestResolveBuildLabConfigUsesFixtureRelativeOutputAndRejectsEscape(t *testing.T) {
	fixture := t.TempDir()
	base := filepath.Join(fixture, "markata-go.toml")
	merge := filepath.Join(fixture, "fast.toml")
	if err := os.WriteFile(base, []byte("[markata-go]\noutput_dir = \"public\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(merge, []byte("[markata-go]\noutput_dir = \"dist\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	configPath, merges, output, err := resolveBuildLabConfig(fixture, base, []string{merge})
	if err != nil || configPath != "markata-go.toml" || len(merges) != 1 || merges[0] != "fast.toml" || output != "dist" {
		t.Fatalf("resolved config = %q, merges=%v, output=%q, err=%v", configPath, merges, output, err)
	}
	if err := os.WriteFile(base, []byte("[markata-go]\noutput_dir = \"../outside\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, _, err := resolveBuildLabConfig(fixture, base, nil); err == nil || !strings.Contains(err.Error(), "escapes fixture") {
		t.Fatalf("outside output error = %v", err)
	}
}

func TestResolveBuildLabConfigKeepsNestedOutputBoundToFixture(t *testing.T) {
	fixture := t.TempDir()
	configDir := filepath.Join(fixture, "config")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "markata-go.toml")
	if err := os.WriteFile(configPath, []byte("[markata-go]\noutput_dir = \"public\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	resolvedConfig, merges, output, err := resolveBuildLabConfig(fixture, configPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	if resolvedConfig != "config/markata-go.toml" || len(merges) != 0 || output != "config/public" {
		t.Fatalf("resolved nested config = %q, merges=%v, output=%q", resolvedConfig, merges, output)
	}
}

func TestParseBuildLabEnvironmentAllowsOnlyReviewedKeys(t *testing.T) {
	for _, value := range []string{"PATH=/tools", "MARKATA_GO_ENCRYPTION_ENABLED=false", "MARKATA_GO_OFFLINE=true"} {
		if _, err := parseBuildLabEnvironment(value); err != nil {
			t.Fatalf("allowed entry %q rejected: %v", value, err)
		}
	}
	for _, value := range []string{"SAFE=value", "API_TOKEN=value", "DATABASE_URL=super-secret", "HOME=/unsafe", "BROKEN"} {
		if _, err := parseBuildLabEnvironment(value); err == nil || strings.Contains(err.Error(), "super-secret") {
			t.Fatalf("entry %q was accepted or leaked: %v", value, err)
		}
	}
}

func TestLoadBuildLabScenarioValidatesRequiredFieldsAndOperations(t *testing.T) {
	tests := []string{
		`{"version":"1","seed":1,"operations":[]}`,
		`{"id":"scenario","seed":1,"operations":[]}`,
		`{"id":"scenario","version":"1","operations":[]}`,
		`{"id":"scenario","version":"2","seed":1,"operations":[]}`,
		`{"id":"scenario","version":"1","seed":1,"operations":[{"type":"unknown"}]}`,
	}
	for _, data := range tests {
		path := filepath.Join(t.TempDir(), "scenario.json")
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := loadBuildLabScenario(path); err == nil {
			t.Fatalf("invalid scenario %s was accepted", data)
		}
	}
	path := filepath.Join(t.TempDir(), "scenario.json")
	if err := os.WriteFile(path, []byte(`{"id":"declared","version":"1","seed":42,"operations":[{"type":"build"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	scenario, err := loadBuildLabScenario(path)
	if err != nil || scenario.Seed != 42 {
		t.Fatalf("scenario = %+v, err=%v", scenario, err)
	}
}

func TestEmitBuildLabResultPrintsBeforeResultFileFailure(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(blocker, []byte("blocker"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout bytes.Buffer
	err := emitBuildLabResult(&stdout, filepath.Join(blocker, "result.json"), buildlab.Result{SchemaVersion: buildlab.ResultSchemaVersion, Verdict: "fail"})
	if err == nil {
		t.Fatal("result-file failure was not reported")
	}
	var result buildlab.Result
	if err := json.Unmarshal(bytes.TrimSpace(stdout.Bytes()), &result); err != nil || result.Verdict != "fail" {
		t.Fatalf("stdout result = %q, err=%v", stdout.String(), err)
	}
}

func TestWriteBuildLabDiagnosticsUsesHumanReadableStderrFormat(t *testing.T) {
	var output bytes.Buffer
	writeBuildLabDiagnostics(&output, buildlab.Result{
		Diagnostics: []buildlab.Diagnostic{{
			Class: buildlab.FailureHarness, Scope: "scenario", Message: "scenario operation 1: path escapes fixture",
		}},
		Checkpoints: []buildlab.Checkpoint{{
			Diagnostics: []buildlab.Diagnostic{{
				Class: buildlab.FailureProduct, Scope: "incremental-comparison", Message: "candidate output differs",
			}},
		}},
	})
	if got, want := output.String(), "Build Lab harness (scenario): scenario operation 1: path escapes fixture\nBuild Lab product at checkpoint 0 (incremental-comparison): candidate output differs\n"; got != want {
		t.Fatalf("diagnostics = %q, want %q", got, want)
	}
}

func TestBuildLabSetupFailureIsVersionedHarnessResult(t *testing.T) {
	result := buildLabSetupFailureResult(errors.New("fixture path is not a directory"))
	if result.SchemaVersion != buildlab.ResultSchemaVersion || result.Verdict != "fail" || result.FailureClass != buildlab.FailureHarness {
		t.Fatalf("setup failure result = %+v", result)
	}
	if len(result.Diagnostics) != 1 || result.Diagnostics[0].Scope != "setup" {
		t.Fatalf("setup failure diagnostics = %+v", result.Diagnostics)
	}
	data, err := result.CanonicalJSON()
	if err != nil || !bytes.Contains(data, []byte(`"schema_version":1`)) {
		t.Fatalf("setup failure JSON = %s, err=%v", data, err)
	}
}
