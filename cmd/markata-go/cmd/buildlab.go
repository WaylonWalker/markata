package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildlab"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/spf13/cobra"
)

var (
	buildLabFixture          string
	buildLabBaseline         string
	buildLabCandidate        string
	buildLabBuildConfig      string
	buildLabScenario         string
	buildLabResult           string
	buildLabSeed             int64
	buildLabTimeout          time.Duration
	buildLabGOMAXPROCS       int
	buildLabCheckDeterminism bool
	buildLabVolatile         string
	buildLabFast             bool
	buildLabExternalTools    string
	buildLabEnvironment      string
)

const buildLabPassVerdict = "pass"

var buildLabCmd = &cobra.Command{Use: "buildlab", Short: "Run isolated baseline and incremental build comparisons"}
var buildLabRunCmd = &cobra.Command{Use: "run", Short: "Run a Build Lab scenario", Args: cobra.NoArgs, RunE: runBuildLabCommand}

func init() {
	rootCmd.AddCommand(buildLabCmd)
	buildLabCmd.AddCommand(buildLabRunCmd)
	flags := buildLabRunCmd.Flags()
	flags.StringVar(&buildLabFixture, "fixture", "", "site fixture directory")
	flags.StringVar(&buildLabBaseline, "baseline", "", "baseline markata-go binary (defaults to this binary)")
	flags.StringVar(&buildLabCandidate, "candidate", "", "candidate markata-go binary (defaults to this binary)")
	flags.StringVar(&buildLabBuildConfig, "build-config", "markata-go.toml", "config path relative to the fixture")
	flags.StringVar(&buildLabScenario, "scenario", "", "scenario JSON file; defaults to a clean build")
	flags.StringVar(&buildLabResult, "result", "", "write structured result JSON to this path")
	flags.Int64Var(&buildLabSeed, "seed", 1, "scenario seed")
	flags.DurationVar(&buildLabTimeout, "timeout", 10*time.Minute, "per-build timeout")
	flags.IntVar(&buildLabGOMAXPROCS, "gomaxprocs", 1, "GOMAXPROCS used by isolated builds")
	flags.BoolVar(&buildLabCheckDeterminism, "check-determinism", true, "compare two clean candidate builds")
	flags.StringVar(&buildLabVolatile, "volatile", ".well-known/time", "comma-separated output paths whose bytes may vary")
	flags.BoolVar(&buildLabFast, "fast", false, "add --fast to baseline and candidate builds")
	flags.StringVar(&buildLabExternalTools, "tool-version", "", "comma-separated tool=version metadata to record")
	flags.StringVar(&buildLabEnvironment, "env", "", "comma-separated PATH, MARKATA_GO_ENCRYPTION_ENABLED, or MARKATA_GO_OFFLINE entries for each build")
}

func isBuildLabCommand(command *cobra.Command) bool {
	for current := command; current != nil; current = current.Parent() {
		if current == buildLabCmd {
			return true
		}
	}
	return false
}

//nolint:gocyclo // CLI validation, configuration resolution, and result emission are one command boundary.
func runBuildLabCommand(cmd *cobra.Command, _ []string) error {
	setupFailure := func(err error) error {
		result := buildLabSetupFailureResult(err)
		emitErr := emitBuildLabResult(outWriter(), buildLabResult, result)
		writeBuildLabDiagnostics(errWriter(), result)
		if emitErr != nil {
			return errors.Join(err, emitErr)
		}
		return err
	}
	if buildLabFixture == "" {
		return setupFailure(fmt.Errorf("--fixture is required"))
	}
	fixture, err := absoluteBuildLabFixture(buildLabFixture)
	if err != nil {
		return setupFailure(err)
	}
	configPathArg := buildLabBuildConfig
	if !cmd.Flags().Changed("build-config") && configPathArg == defaultConfigFilename {
		if _, statErr := os.Stat(filepath.Join(fixture, configPathArg)); os.IsNotExist(statErr) {
			configPathArg = ""
			for _, candidate := range []string{"markata-go.yaml", "markata-go.yml", "markata-go.json"} {
				if _, candidateErr := os.Stat(filepath.Join(fixture, candidate)); candidateErr == nil {
					configPathArg = candidate
					break
				}
			}
		}
	}
	configPath, mergePaths, effectiveOutput, err := resolveBuildLabConfig(fixture, configPathArg, mergeConfigFiles)
	if err != nil {
		return setupFailure(err)
	}
	baseline, err := absoluteCommandPath(buildLabBaseline)
	if err != nil {
		return setupFailure(fmt.Errorf("baseline binary: %w", err))
	}
	candidate, err := absoluteCommandPath(buildLabCandidate)
	if err != nil {
		return setupFailure(fmt.Errorf("candidate binary: %w", err))
	}
	scenario, err := loadBuildLabScenario(buildLabScenario)
	if err != nil {
		return setupFailure(err)
	}
	if buildLabScenario == "" || cmd.Flags().Changed("seed") {
		scenario.Seed = buildLabSeed
	}
	baselineArgs := []string{"build"}
	if buildLabFast {
		baselineArgs = append(baselineArgs, "--fast")
	}
	if configPath != "" {
		baselineArgs = append(baselineArgs, "-c", configPath)
	}
	for _, merge := range mergePaths {
		if merge != "" {
			baselineArgs = append(baselineArgs, "-m", merge)
		}
	}
	if outputDir != "" {
		effectiveOutput, err = buildLabRelativePath(fixture, outputDir)
		if err != nil {
			return setupFailure(fmt.Errorf("output directory: %w", err))
		}
		baselineArgs = append(baselineArgs, "-o", effectiveOutput)
	}
	candidateArgs := append([]string(nil), baselineArgs...)
	classes := parseBuildLabVolatile(buildLabVolatile)
	externalTools, err := parseBuildLabToolVersions(buildLabExternalTools)
	if err != nil {
		return setupFailure(err)
	}
	environment, err := parseBuildLabEnvironment(buildLabEnvironment)
	if err != nil {
		return setupFailure(err)
	}
	result, runErr := buildlab.RunScenario(context.Background(), buildlab.ScenarioRunConfig{
		Fixture: fixture, Scenario: scenario,
		Baseline:  buildlab.BuildCommand{Binary: baseline, Args: baselineArgs, OutputDir: effectiveOutput, Timeout: buildLabTimeout, Env: environment},
		Candidate: buildlab.BuildCommand{Binary: candidate, Args: candidateArgs, OutputDir: effectiveOutput, Timeout: buildLabTimeout, Env: environment},
		Classes:   classes, CheckDeterminism: buildLabCheckDeterminism, GOMAXPROCS: buildLabGOMAXPROCS, ExternalTools: externalTools,
	})
	if err := emitBuildLabResult(outWriter(), buildLabResult, result); err != nil {
		return err
	}
	writeBuildLabDiagnostics(errWriter(), result)
	if runErr != nil {
		return runErr
	}
	if result.Verdict != buildLabPassVerdict {
		return fmt.Errorf("build lab verdict: %s", result.Verdict)
	}
	return nil
}

func buildLabSetupFailureResult(err error) buildlab.Result {
	return buildlab.Result{
		SchemaVersion: buildlab.ResultSchemaVersion,
		Checkpoints:   []buildlab.Checkpoint{},
		Verdict:       "fail",
		FailureClass:  buildlab.FailureHarness,
		Diagnostics: []buildlab.Diagnostic{{
			Class: buildlab.FailureHarness, Scope: "setup", Message: err.Error(),
		}},
	}
}

func parseBuildLabVolatile(value string) map[string]buildlab.OutputClass {
	classes := make(map[string]buildlab.OutputClass)
	for _, path := range strings.Split(value, ",") {
		path = filepath.ToSlash(filepath.Clean(strings.TrimSpace(path)))
		if path != "" && path != "." {
			classes[path] = buildlab.ClassVolatile
		}
	}
	return classes
}

// emitBuildLabResult always writes stdout before trying the optional result file.
func emitBuildLabResult(output io.Writer, resultPath string, result buildlab.Result) error {
	data, err := result.CanonicalJSON()
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(output, string(data)); err != nil {
		return fmt.Errorf("write Build Lab result to stdout: %w", err)
	}
	if resultPath != "" {
		if err := buildlab.WriteResult(resultPath, result); err != nil {
			return fmt.Errorf("write Build Lab result: %w", err)
		}
	}
	return nil
}

func writeBuildLabDiagnostics(output io.Writer, result buildlab.Result) {
	for _, diagnostic := range result.Diagnostics {
		_, _ = fmt.Fprintf(output, "Build Lab %s (%s): %s\n", diagnostic.Class, diagnostic.Scope, diagnostic.Message)
	}
	for checkpointIndex := range result.Checkpoints {
		checkpoint := &result.Checkpoints[checkpointIndex]
		for _, diagnostic := range checkpoint.Diagnostics {
			_, _ = fmt.Fprintf(output, "Build Lab %s at checkpoint %d (%s): %s\n", diagnostic.Class, checkpointIndex, diagnostic.Scope, diagnostic.Message)
		}
	}
}

func absoluteBuildLabFixture(value string) (string, error) {
	path, err := filepath.Abs(value)
	if err != nil {
		return "", fmt.Errorf("fixture path: %w", err)
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("fixture path: %w", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("fixture path: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("fixture path %q is not a directory", value)
	}
	return filepath.Clean(path), nil
}

func buildLabRelativePath(root, value string) (string, error) {
	root, err := canonicalBuildLabPath(root)
	if err != nil {
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("path is empty")
	}
	path := value
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	path, err = canonicalBuildLabPath(path)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes fixture", value)
	}
	if err := verifyBuildLabPath(root, path); err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

func verifyBuildLabPath(root, path string) error {
	resolvedRoot, err := canonicalBuildLabPath(root)
	if err != nil {
		return err
	}
	root = resolvedRoot
	for current := path; ; current = filepath.Dir(current) {
		resolved, err := filepath.EvalSymlinks(current)
		if err == nil {
			if !buildLabPathWithin(root, resolved) {
				return fmt.Errorf("path %q resolves outside fixture", path)
			}
			return nil
		}
		if !os.IsNotExist(err) {
			return err
		}
		if info, statErr := os.Lstat(current); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path %q contains unresolved symlink %q", path, current)
		} else if statErr != nil && !os.IsNotExist(statErr) {
			return statErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			return nil
		}
	}
}

// canonicalBuildLabPath resolves the existing portion of path before
// appending any missing components. This keeps containment checks correct
// when a platform exposes the same directory through multiple path spellings,
// such as /var and /private/var on macOS or short and long paths on Windows.
func canonicalBuildLabPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	missing := make([]string, 0)
	current := abs
	for {
		resolved, evalErr := filepath.EvalSymlinks(current)
		if evalErr == nil {
			for i := len(missing) - 1; i >= 0; i-- {
				resolved = filepath.Join(resolved, missing[i])
			}
			return filepath.Clean(resolved), nil
		}
		if !os.IsNotExist(evalErr) {
			return "", evalErr
		}
		if info, statErr := os.Lstat(current); statErr == nil && info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("path contains unresolved symlink %q", current)
		} else if statErr != nil && !os.IsNotExist(statErr) {
			return "", statErr
		}
		parent := filepath.Dir(current)
		if parent == current {
			return abs, nil
		}
		missing = append(missing, filepath.Base(current))
		current = parent
	}
}

func buildLabPathWithin(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func resolveBuildLabConfig(fixture, configPath string, mergePaths []string) (resolvedConfig string, resolvedMerges []string, effectiveOutput string, err error) {
	relConfig := ""
	if configPath != "" {
		relConfig, err = buildLabRelativePath(fixture, configPath)
		if err != nil {
			return "", nil, "", fmt.Errorf("build config: %w", err)
		}
	}
	relMerges := make([]string, 0, len(mergePaths))
	absMerges := make([]string, 0, len(mergePaths))
	for _, path := range mergePaths {
		if path == "" {
			continue
		}
		rel, e := buildLabRelativePath(fixture, path)
		if e != nil {
			return "", nil, "", fmt.Errorf("merge config %q: %w", path, e)
		}
		relMerges = append(relMerges, rel)
		absMerges = append(absMerges, filepath.Join(fixture, filepath.FromSlash(rel)))
	}
	basePath := ""
	if relConfig != "" {
		basePath = filepath.Join(fixture, filepath.FromSlash(relConfig))
	}
	cfg, err := loadBuildLabConfig(basePath, absMerges...)
	if err != nil {
		return "", nil, "", fmt.Errorf("resolve Build Lab output directory: %w", err)
	}
	effectiveOutput = cfg.OutputDir
	if effectiveOutput == "" {
		effectiveOutput = defaultOutputDir
	}
	if filepath.IsAbs(effectiveOutput) {
		effectiveOutput, err = buildLabRelativePath(fixture, effectiveOutput)
	} else {
		baseDir := fixture
		if basePath != "" {
			baseDir = filepath.Dir(basePath)
		}
		effectiveOutput, err = buildLabRelativePath(fixture, filepath.Join(baseDir, effectiveOutput))
	}
	if err != nil {
		return "", nil, "", fmt.Errorf("configured output directory: %w", err)
	}
	return relConfig, relMerges, effectiveOutput, nil
}

func loadBuildLabConfig(basePath string, mergePaths ...string) (*models.Config, error) {
	return config.LoadWithMergeOptions(config.LoadOptions{DisableDotEnv: true, DisableEnvOverrides: true}, basePath, mergePaths...)
}

func parseBuildLabToolVersions(value string) (map[string]string, error) {
	tools := make(map[string]string)
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		key, version, ok := strings.Cut(item, "=")
		key, version = strings.TrimSpace(key), strings.TrimSpace(version)
		if !ok || key == "" || version == "" {
			return nil, fmt.Errorf("invalid --tool-version value %q; use name=version", item)
		}
		tools[key] = version
	}
	return tools, nil
}
func parseBuildLabEnvironment(value string) ([]string, error) {
	entries := make([]string, 0)
	for _, item := range strings.Split(value, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		key, envValue, ok := strings.Cut(item, "=")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid --env entry %d; use KEY=value", len(entries))
		}
		entries = append(entries, strings.TrimSpace(key)+"="+envValue)
	}
	if err := buildlab.ValidateEnvironment(entries); err != nil {
		return nil, err
	}
	return entries, nil
}
func absoluteCommandPath(value string) (string, error) {
	if value == "" {
		value, err := exec.LookPath(os.Args[0])
		if err != nil {
			return "", err
		}
		return filepath.Abs(value)
	}
	if !strings.ContainsRune(value, filepath.Separator) {
		if path, err := exec.LookPath(value); err == nil {
			return filepath.Abs(path)
		}
	}
	path, err := filepath.Abs(value)
	if err != nil {
		return "", err
	}
	return path, nil
}
func loadBuildLabScenario(path string) (buildlab.Scenario, error) {
	if path == "" {
		return buildlab.Scenario{ID: "cli-smoke", Version: "1", Operations: []buildlab.Operation{{Type: buildlab.OpClearCache}, {Type: buildlab.OpClearOutput}, {Type: buildlab.OpBuild}, {Type: buildlab.OpBuild}}}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return buildlab.Scenario{}, fmt.Errorf("read scenario: %w", err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return buildlab.Scenario{}, fmt.Errorf("parse scenario: %w", err)
	}
	for _, field := range []string{"id", "version", "seed", "operations"} {
		value, ok := fields[field]
		if !ok || strings.TrimSpace(string(value)) == "" || strings.TrimSpace(string(value)) == "null" {
			return buildlab.Scenario{}, fmt.Errorf("scenario field %q is required", field)
		}
	}
	var scenario buildlab.Scenario
	if err := json.Unmarshal(data, &scenario); err != nil {
		return buildlab.Scenario{}, fmt.Errorf("parse scenario: %w", err)
	}
	if err := scenario.Validate(); err != nil {
		return buildlab.Scenario{}, fmt.Errorf("validate scenario: %w", err)
	}
	return scenario, nil
}
