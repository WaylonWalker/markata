package buildlab

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

const ResultSchemaVersion = 1

const (
	verdictFail       = "fail"
	verdictPass       = "pass"
	buildLabWindowsOS = "windows"
)

// BuildCommand describes one real Markata-Go process.  Binary should normally
// be an absolute path to a built executable so the baseline and candidate
// identity can be recorded without relying on process state.
type BuildCommand struct {
	Binary    string
	Args      []string
	OutputDir string
	Timeout   time.Duration
	Env       []string
}

// Workspace is an isolated copy of a fixture and its process-local state.
// IsolationConfig redirects built-in Markata-Go cache settings into this
// workspace and is appended after fixture-provided merge configs.
type Workspace struct {
	Root            string
	SiteDir         string
	HomeDir         string
	XDGCacheDir     string
	XDGConfigDir    string
	XDGDataDir      string
	AppDataDir      string
	LocalAppDataDir string
	TempDir         string
	MarkataCache    string
	IsolationConfig string
}

// NewWorkspace copies source inputs into a private site and creates the
// directories used by home, XDG cache, temporary files, and Markata cache
// state. Known generated build state is not copied.
func NewWorkspace(fixture, parent string) (Workspace, error) {
	return newWorkspace(fixture, parent, buildLabFixtureExclusions())
}

func newWorkspace(fixture, parent string, excluded []string) (Workspace, error) {
	if fixture == "" {
		return Workspace{}, fmt.Errorf("fixture path is empty")
	}
	fixtureInfo, err := os.Lstat(fixture)
	if err != nil {
		return Workspace{}, fmt.Errorf("inspect fixture: %w", err)
	}
	if fixtureInfo.Mode()&os.ModeSymlink != 0 {
		return Workspace{}, fmt.Errorf("fixture is a symlink: %q", fixture)
	}
	if !fixtureInfo.IsDir() {
		return Workspace{}, fmt.Errorf("fixture is not a directory: %q", fixture)
	}
	if parent == "" {
		parent = os.TempDir()
	}
	root, err := os.MkdirTemp(parent, "markata-buildlab-")
	if err != nil {
		return Workspace{}, fmt.Errorf("create workspace: %w", err)
	}
	fail := func(e error) (Workspace, error) {
		_ = os.RemoveAll(root)
		return Workspace{}, e
	}
	w := Workspace{
		Root:            root,
		SiteDir:         filepath.Join(root, "site"),
		HomeDir:         filepath.Join(root, "home"),
		XDGCacheDir:     filepath.Join(root, "xdg-cache"),
		XDGConfigDir:    filepath.Join(root, "xdg-config"),
		XDGDataDir:      filepath.Join(root, "xdg-data"),
		AppDataDir:      filepath.Join(root, "appdata"),
		LocalAppDataDir: filepath.Join(root, "local-appdata"),
		TempDir:         filepath.Join(root, "tmp"),
		MarkataCache:    filepath.Join(root, "markata-cache"),
		IsolationConfig: filepath.Join(root, "buildlab-isolation.json"),
	}
	if err := validateFixtureCopyPaths(fixture, w.SiteDir); err != nil {
		return fail(fmt.Errorf("validate workspace copy paths: %w", err))
	}
	if err := copyDir(fixture, w.SiteDir, excluded); err != nil {
		return fail(fmt.Errorf("copy fixture: %w", err))
	}
	for _, dir := range []string{w.HomeDir, w.XDGCacheDir, w.XDGConfigDir, w.XDGDataDir, w.AppDataDir, w.LocalAppDataDir, w.TempDir, w.MarkataCache} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fail(fmt.Errorf("create workspace state %s: %w", dir, err))
		}
	}
	if err := writeIsolationConfig(w.IsolationConfig, w.MarkataCache); err != nil {
		return fail(fmt.Errorf("write workspace isolation config: %w", err))
	}
	return w, nil
}

func buildLabFixtureExclusions(commands ...BuildCommand) []string {
	excluded := []string{".markata", ".markata.cache", ".markata-cache", "cache", "output"}
	excluded = append(excluded, []string{".markata-css_minify-cache", ".markata-fontpack-cache", ".markata-fonts.json", ".markata-js_minify-cache"}...)
	for _, command := range commands {
		output := command.OutputDir
		if output == "" {
			output = outputArg(command.Args)
		}
		if output == "" || filepath.IsAbs(output) {
			continue
		}
		clean := filepath.Clean(filepath.FromSlash(output))
		if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
			continue
		}
		excluded = append(excluded, clean)
	}
	return excluded
}

func writeIsolationConfig(path, cacheDir string) error {
	pluginCache := func(name string) string {
		return filepath.Join(cacheDir, name)
	}
	data, err := json.Marshal(map[string]any{
		"markata-go": map[string]any{
			"cache_cleanup_async": false,
			"cache_dir":           cacheDir,
			"assets": map[string]any{
				"cache_dir": pluginCache("assets"),
			},
			"blogroll": map[string]any{
				"cache_dir": pluginCache("blogroll"),
			},
			"embeds": map[string]any{
				"cache_dir": pluginCache("embeds"),
			},
			"image_optimization": map[string]any{
				"cache_dir": pluginCache("image-optimization"),
			},
			"mentions": map[string]any{
				"cache_dir": pluginCache("mentions"),
			},
			"search": map[string]any{
				"pagefind": map[string]any{
					"cache_dir": pluginCache("pagefind"),
				},
			},
			"tailwind": map[string]any{
				"cache_dir": pluginCache("tailwind"),
			},
			"webmentions": map[string]any{
				"cache_dir": pluginCache("webmentions"),
			},
		},
	})
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

// Remove deletes all files created for the workspace.
func (w Workspace) Remove() error {
	if w.Root == "" {
		return nil
	}
	return os.RemoveAll(w.Root)
}

// Environment returns a child environment containing only the reviewed ambient
// variables and workspace-controlled paths. Explicit variables are validated
// before they are added.
func (w Workspace) Environment(extra []string, gomaxprocs int) ([]string, error) {
	if err := ValidateEnvironment(extra); err != nil {
		return nil, err
	}
	if err := w.validateEnvironmentPaths(); err != nil {
		return nil, err
	}
	values := make(map[string]string)
	for _, key := range inheritedBuildEnvironmentKeys() {
		if value, ok := os.LookupEnv(key); ok {
			values[key] = value
		}
	}
	for _, entry := range extra {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
		}
	}
	configDir := workspaceDir(w.XDGConfigDir, w.Root, "xdg-config")
	dataDir := workspaceDir(w.XDGDataDir, w.Root, "xdg-data")
	values["HOME"] = w.HomeDir
	values["XDG_CACHE_HOME"] = w.XDGCacheDir
	values["XDG_CONFIG_HOME"] = configDir
	values["XDG_DATA_HOME"] = dataDir
	values["TMPDIR"] = w.TempDir
	values["TMP"] = w.TempDir
	values["TEMP"] = w.TempDir
	values["TZ"] = "UTC"
	values["LANG"] = "C.UTF-8"
	values["LC_ALL"] = "C.UTF-8"
	values["SOURCE_DATE_EPOCH"] = "0"
	values[runtimeenv.EnvDisableDotEnv] = "1"
	if runtime.GOOS == buildLabWindowsOS {
		values["USERPROFILE"] = w.HomeDir
		values["APPDATA"] = workspaceDir(w.AppDataDir, w.Root, "appdata")
		values["LOCALAPPDATA"] = workspaceDir(w.LocalAppDataDir, w.Root, "local-appdata")
	}
	if gomaxprocs > 0 {
		values["GOMAXPROCS"] = strconv.Itoa(gomaxprocs)
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	env := make([]string, 0, len(keys))
	for _, key := range keys {
		env = append(env, key+"="+values[key])
	}
	return env, nil
}

func (w Workspace) validateEnvironmentPaths() error {
	if w.Root == "" {
		return fmt.Errorf("workspace root is empty")
	}
	paths := map[string]string{
		"HOME":             w.HomeDir,
		"XDG_CACHE_HOME":   w.XDGCacheDir,
		"XDG_CONFIG_HOME":  workspaceDir(w.XDGConfigDir, w.Root, "xdg-config"),
		"XDG_DATA_HOME":    workspaceDir(w.XDGDataDir, w.Root, "xdg-data"),
		"TMPDIR":           w.TempDir,
		"site directory":   w.SiteDir,
		"markata cache":    w.MarkataCache,
		"isolation config": w.IsolationConfig,
	}
	if runtime.GOOS == buildLabWindowsOS {
		paths["APPDATA"] = workspaceDir(w.AppDataDir, w.Root, "appdata")
		paths["LOCALAPPDATA"] = workspaceDir(w.LocalAppDataDir, w.Root, "local-appdata")
	}
	for name, path := range paths {
		if path == "" {
			// Optional fields are not used by every manually constructed
			// Workspace, but any value that is supplied must be confined.
			if name == "site directory" || name == "markata cache" || name == "isolation config" {
				continue
			}
			return fmt.Errorf("workspace path for %s is empty", name)
		}
		if err := validateWorkspacePath(w.Root, path); err != nil {
			return fmt.Errorf("workspace path for %s: %w", name, err)
		}
	}
	return nil
}

func validateWorkspacePath(root, target string) error {
	root, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve workspace root: %w", err)
	}
	target, err = filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolve workspace path: %w", err)
	}
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path %q is outside workspace root %q", target, root)
	}

	rootInfo, err := os.Lstat(root)
	if err == nil {
		if rootInfo.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("workspace root is a symlink: %q", root)
		}
		if !rootInfo.IsDir() {
			return fmt.Errorf("workspace root is not a directory: %q", root)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("inspect workspace root: %w", err)
	}

	current := root
	for _, component := range strings.Split(relative, string(filepath.Separator)) {
		if component == "" || component == "." {
			continue
		}
		current = filepath.Join(current, component)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			return fmt.Errorf("inspect %q: %w", current, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path enters symlink %q", current)
		}
		if !info.IsDir() && current != target {
			return fmt.Errorf("path component is not a directory: %q", current)
		}
	}
	return nil
}

func inheritedBuildEnvironmentKeys() []string {
	keys := []string{"PATH"}
	if runtime.GOOS == buildLabWindowsOS {
		keys = append(keys, "SystemRoot", "WINDIR", "COMSPEC", "PATHEXT")
	}
	return keys
}

func workspaceDir(configured, root, name string) string {
	if configured != "" {
		return configured
	}
	return filepath.Join(root, name)
}

var buildLabAllowedEnvironmentKeys = map[string]struct{}{
	"PATH":                          {},
	runtimeenv.EnvEncryptionEnabled: {},
	runtimeenv.EnvOffline:           {},
}

// ValidateEnvironment rejects malformed or unapproved entries before a caller
// forwards them to an isolated build process.
func ValidateEnvironment(entries []string) error {
	seen := make(map[string]struct{}, len(entries))
	for i, entry := range entries {
		key, _, ok := strings.Cut(entry, "=")
		trimmedKey := strings.TrimSpace(key)
		if trimmedKey != key {
			return fmt.Errorf("environment entry %d has invalid key %q", i, trimmedKey)
		}
		key = trimmedKey
		if !ok || key == "" {
			return fmt.Errorf("invalid environment entry %d; use KEY=value", i)
		}
		if !validEnvironmentKey(key) {
			return fmt.Errorf("environment entry %d has invalid key %q", i, key)
		}
		if _, ok := buildLabAllowedEnvironmentKeys[key]; !ok {
			return fmt.Errorf("environment key %q is not allowed by Build Lab", key)
		}
		if _, ok := seen[key]; ok {
			return fmt.Errorf("environment key %q is specified more than once", key)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func validEnvironmentKey(key string) bool {
	for i := 0; i < len(key); i++ {
		ch := key[i]
		if i == 0 {
			if !(ch == '_' || ch >= 'A' && ch <= 'Z') {
				return false
			}
			continue
		}
		if !(ch == '_' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9') {
			return false
		}
	}
	return key != ""
}

func (c BuildCommand) run(ctx context.Context, w Workspace, clean bool, gomaxprocs int) RunResult {
	if c.Binary == "" {
		return RunResult{Err: fmt.Errorf("build command binary is empty")}
	}
	args, err := c.args(w, clean)
	if err != nil {
		return RunResult{Err: err}
	}
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Minute
	}
	environment, err := w.Environment(c.Env, gomaxprocs)
	if err != nil {
		return RunResult{Err: fmt.Errorf("build environment: %w", err)}
	}
	return Run(ctx, RunConfig{
		Command: c.Binary,
		Args:    args,
		CWD:     w.SiteDir,
		Env:     environment,
		Timeout: timeout,
	})
}

func (c BuildCommand) args(w Workspace, clean bool) (args []string, err error) {
	output, err := c.outputPath(w)
	if err != nil {
		return nil, err
	}
	args = append([]string(nil), c.Args...)
	// The isolation config must follow all fixture and user merge configs so
	// its private cache directory cannot be overridden by the fixture.
	if w.IsolationConfig != "" {
		args = append(args, "--merge-config", w.IsolationConfig)
	}
	// Always override the command's output path with a workspace-local path.
	// This prevents a copied fixture containing an absolute output_dir from
	// causing baseline and candidate processes to write into the same tree.
	// Bind site selection to the copied workspace. This is defense in depth
	// against both caller-provided flags and ambient site-directory settings.
	args = append(args, "--output", output, "--site-dir", w.SiteDir)
	if clean && !hasArg(args, "--clean") {
		args = append(args, "--clean")
	}
	return args, nil
}

func (c BuildCommand) outputPath(w Workspace) (string, error) {
	output := c.OutputDir
	if output == "" {
		output = outputArg(c.Args)
	}
	if output == "" {
		output = "output"
	}
	if filepath.IsAbs(output) {
		clean := filepath.Clean(output)
		if err := validateOutputPath(w.SiteDir, clean); err != nil {
			return "", fmt.Errorf("build output directory %q escapes workspace", output)
		}
		return clean, nil
	}
	clean := filepath.Clean(filepath.Join(w.SiteDir, output))
	if err := validateOutputPath(w.SiteDir, clean); err != nil {
		return "", fmt.Errorf("build output directory %q escapes workspace", output)
	}
	return clean, nil
}

func validateOutputPath(root, output string) error {
	rel, err := filepath.Rel(root, output)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("output path escapes workspace")
	}
	part := root
	for _, component := range strings.Split(rel, string(filepath.Separator)) {
		if component == "." || component == "" {
			continue
		}
		part = filepath.Join(part, component)
		info, err := os.Lstat(part)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output path enters symlink %q", part)
		}
	}
	return nil
}

func hasArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}

type BuildIdentity struct {
	Binary       string `json:"binary"`
	BinarySHA256 string `json:"binary_sha256,omitempty"`
	GitSHA       string `json:"git_sha,omitempty"`
}

type FixtureIdentity struct {
	Path   string `json:"path"`
	Digest string `json:"digest"`
}

type EnvironmentIdentity struct {
	Timezone             string            `json:"timezone"`
	Locale               string            `json:"locale"`
	SourceDateEpoch      string            `json:"source_date_epoch"`
	GOMAXPROCS           int               `json:"gomaxprocs"`
	GOOS                 string            `json:"goos"`
	GOARCH               string            `json:"goarch"`
	GoVersion            string            `json:"go_version"`
	ExternalToolVersions map[string]string `json:"external_tool_versions,omitempty"`
}

func defaultEnvironment(gomaxprocs int, tools map[string]string) EnvironmentIdentity {
	if gomaxprocs <= 0 {
		gomaxprocs = runtime.GOMAXPROCS(0)
	}
	copyTools := make(map[string]string, len(tools))
	for key, value := range tools {
		copyTools[key] = value
	}
	return EnvironmentIdentity{
		Timezone: "UTC", Locale: "C.UTF-8", SourceDateEpoch: "0",
		GOMAXPROCS: gomaxprocs, GOOS: runtime.GOOS, GOARCH: runtime.GOARCH,
		GoVersion: runtime.Version(), ExternalToolVersions: copyTools,
	}
}

type RunObservation struct {
	ExitCode        int          `json:"exit_code"`
	Successful      bool         `json:"successful"`
	TimedOut        bool         `json:"timed_out"`
	StdoutTruncated bool         `json:"stdout_truncated,omitempty"`
	StderrTruncated bool         `json:"stderr_truncated,omitempty"`
	DurationSeconds float64      `json:"duration_seconds"`
	ManifestDigest  string       `json:"manifest_digest,omitempty"`
	StdoutSHA256    string       `json:"stdout_sha256,omitempty"`
	StderrSHA256    string       `json:"stderr_sha256,omitempty"`
	FailureClass    FailureClass `json:"failure_class,omitempty"`
	Error           string       `json:"error,omitempty"`
}

type FailureClass string

// FailureClass distinguishes a failed configured build from a failure in the
// measuring instrument. Empty means that the observation succeeded.
const (
	FailureProduct FailureClass = "product"
	FailureHarness FailureClass = "harness"
)

// Diagnostic explains a failed build or comparison without requiring callers
// to reconstruct the failure from process hashes and manifest records.
type Diagnostic struct {
	Class   FailureClass  `json:"class"`
	Scope   string        `json:"scope"`
	Message string        `json:"message"`
	Diff    *ManifestDiff `json:"diff,omitempty"`
}

type CorrectnessResult struct {
	DifferentialEqual     bool         `json:"differential_equal"`
	IncrementalEqual      bool         `json:"incremental_equal"`
	IncrementalApplicable bool         `json:"incremental_applicable"`
	DeterministicEqual    bool         `json:"deterministic_equal"`
	DifferentialDiff      ManifestDiff `json:"differential_diff,omitempty"`
	IncrementalDiff       ManifestDiff `json:"incremental_diff,omitempty"`
	DeterminismDiff       ManifestDiff `json:"determinism_diff,omitempty"`
}

type Checkpoint struct {
	OperationIndex       int               `json:"operation_index"`
	Operation            string            `json:"operation"`
	Baseline             RunObservation    `json:"baseline"`
	CandidateClean       RunObservation    `json:"candidate_clean"`
	CandidateDeterminism *RunObservation   `json:"candidate_determinism,omitempty"`
	CandidateIncremental RunObservation    `json:"candidate_incremental"`
	Correctness          CorrectnessResult `json:"correctness"`
	Diagnostics          []Diagnostic      `json:"diagnostics,omitempty"`
}

// Result is stable, machine-readable evidence from a scenario run.  Timing
// and process hashes are observations; correctness fields are the gate.
type Result struct {
	SchemaVersion int                 `json:"schema_version"`
	Baseline      BuildIdentity       `json:"baseline"`
	Candidate     BuildIdentity       `json:"candidate"`
	Fixture       FixtureIdentity     `json:"fixture"`
	Scenario      Scenario            `json:"scenario"`
	Environment   EnvironmentIdentity `json:"environment"`
	Diagnostics   []Diagnostic        `json:"diagnostics,omitempty"`
	Checkpoints   []Checkpoint        `json:"checkpoints"`
	Verdict       string              `json:"verdict"`
	FailureClass  FailureClass        `json:"failure_class,omitempty"`
}

// CanonicalJSON returns deterministic JSON for a result with checkpoint order
// preserved from the declarative scenario.
func (r Result) CanonicalJSON() ([]byte, error) {
	if r.SchemaVersion == 0 {
		r.SchemaVersion = ResultSchemaVersion
	}
	return json.Marshal(r)
}

func (r Result) Digest() (string, error) {
	b, err := r.CanonicalJSON()
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:]), nil
}

// ScenarioRunConfig describes a baseline/candidate comparison.
type ScenarioRunConfig struct {
	Fixture          string
	Scenario         Scenario
	Baseline         BuildCommand
	Candidate        BuildCommand
	ParentDir        string
	Classes          map[string]OutputClass
	Comparators      map[string]func(FileRecord, FileRecord) bool
	CheckDeterminism bool
	GOMAXPROCS       int
	ExternalTools    map[string]string
}

// RunScenario executes each build operation and compares all three required
// paths.  A failed command is represented in Result and makes the verdict
// fail; callers can persist the result even when a candidate is broken.
//
//nolint:gocyclo // Scenario validation and checkpoint orchestration are one public operation.
func RunScenario(ctx context.Context, cfg ScenarioRunConfig) (result Result, runErr error) {
	result = Result{
		SchemaVersion: ResultSchemaVersion,
		Baseline:      identityForCommand(cfg.Baseline),
		Candidate:     identityForCommand(cfg.Candidate),
		Scenario:      cfg.Scenario,
		Environment:   defaultEnvironment(cfg.GOMAXPROCS, cfg.ExternalTools),
		Checkpoints:   []Checkpoint{},
		Verdict:       verdictFail,
	}
	if cfg.Fixture == "" {
		return failResult(result, "fixture", errors.New("fixture path is empty"))
	}
	if err := cfg.Scenario.Validate(); err != nil {
		return failResult(result, "scenario", fmt.Errorf("scenario: %w", err))
	}
	excluded := buildLabFixtureExclusions(cfg.Baseline, cfg.Candidate)
	result.Fixture.Path = cfg.Fixture
	fixtureManifest, err := BuildManifest(cfg.Fixture, nil, excluded...)
	if err != nil {
		return failResult(result, "fixture-manifest", fmt.Errorf("fixture manifest: %w", err))
	}
	fixtureDigest, err := fixtureManifest.Digest()
	if err != nil {
		return failResult(result, "fixture-digest", fmt.Errorf("fixture digest: %w", err))
	}
	result.Fixture.Digest = fixtureDigest
	source, err := newWorkspace(cfg.Fixture, cfg.ParentDir, excluded)
	if err != nil {
		return failResult(result, "workspace", err)
	}
	defer func() {
		if err := cleanupWorkspace(source); err != nil {
			result.Verdict = verdictFail
			result.Diagnostics = append(result.Diagnostics, Diagnostic{
				Class: FailureHarness, Scope: "workspace-cleanup", Message: err.Error(),
			})
			result.FailureClass = FailureHarness
			runErr = errors.Join(runErr, err)
		}
	}()

	builtOnce := false
	pendingMutation := false
	pendingMutationIndex := -1
	for operationIndex := range cfg.Scenario.Operations {
		operation := cfg.Scenario.Operations[operationIndex]
		if operation.Type != OpBuild {
			if err := applyWorkspaceOperation(source, operation, cfg.Candidate); err != nil {
				return failResult(result, "scenario-operation", &ScenarioError{Operation: operationIndex, Err: err})
			}
			if isSemanticMutation(operation.Type) {
				if pendingMutation {
					return failResult(result, "scenario-operation", &ScenarioError{Operation: operationIndex, Err: fmt.Errorf("mutation must be followed by a build")})
				}
				pendingMutation = true
				pendingMutationIndex = operationIndex
			}
			continue
		}

		checkpoint, checkpointErr := runScenarioCheckpoint(ctx, cfg, source, operationIndex, operation, builtOnce, result.Environment.GOMAXPROCS)
		result.Checkpoints = append(result.Checkpoints, checkpoint)
		if checkpointErr != nil {
			return failResult(result, "checkpoint", &ScenarioError{Operation: operationIndex, Err: checkpointErr})
		}
		builtOnce = true
		pendingMutation = false
		pendingMutationIndex = -1
	}
	if pendingMutation {
		return failResult(result, "scenario-operation", &ScenarioError{Operation: pendingMutationIndex, Err: fmt.Errorf("scenario ends with a mutation without a following build")})
	}
	if len(result.Checkpoints) == 0 {
		return failResult(result, "scenario", errors.New("scenario contains no build operation"))
	}
	result.Verdict = verdictPass
	for i := range result.Checkpoints {
		checkpoint := &result.Checkpoints[i]
		checkpoint.Diagnostics = checkpointDiagnostics(*checkpoint)
		if !checkpoint.Baseline.Successful || !checkpoint.CandidateClean.Successful ||
			!checkpoint.CandidateIncremental.Successful ||
			!checkpoint.Correctness.DifferentialEqual || !checkpoint.Correctness.DeterministicEqual ||
			(checkpoint.Correctness.IncrementalApplicable &&
				!checkpoint.Correctness.IncrementalEqual) {
			result.Verdict = verdictFail
			break
		}
	}
	result.FailureClass = resultFailureClass(result)
	return result, nil
}

func failResult(result Result, scope string, err error) (Result, error) {
	result.Verdict = verdictFail
	result.FailureClass = FailureHarness
	result.Diagnostics = append(result.Diagnostics, Diagnostic{
		Class: FailureHarness, Scope: scope, Message: err.Error(),
	})
	return result, err
}

// runScenarioCheckpoint compares the candidate's incremental output with
// clean baseline and candidate builds for one build operation. Derived clean
// workspaces are released after each manifest is captured so large generated
// output trees do not accumulate across the checkpoint.
func runScenarioCheckpoint(
	ctx context.Context,
	cfg ScenarioRunConfig,
	source Workspace,
	operationIndex int,
	operation Operation,
	incrementalApplicable bool,
	gomaxprocs int,
) (checkpoint Checkpoint, runErr error) {
	var incremental RunResult
	if incrementalApplicable {
		incremental = cfg.Candidate.run(ctx, source, false, gomaxprocs)
	} else {
		incremental = cfg.Candidate.run(ctx, source, true, gomaxprocs)
	}
	var incrementalManifest Manifest
	var incrementalManifestErr error
	if incremental.Successful() {
		incrementalManifest, incrementalManifestErr = manifestForRun(source, cfg.Candidate, cfg.Classes)
	}
	incrementalObservation := observe(incremental, incrementalManifest, incrementalManifestErr)
	checkpoint = Checkpoint{
		OperationIndex:       operationIndex,
		Operation:            string(operation.Type),
		CandidateIncremental: incrementalObservation,
		Correctness: CorrectnessResult{
			IncrementalApplicable: incrementalApplicable,
			DeterministicEqual:    !cfg.CheckDeterminism,
		},
	}

	excluded := buildLabFixtureExclusions(cfg.Baseline, cfg.Candidate)
	baselineObservation, baselineManifest, baselineErr := runCleanWorkspace(ctx, source.SiteDir, cfg.ParentDir, excluded, cfg.Baseline, cfg.Classes, gomaxprocs)
	candidateCleanObservation, candidateCleanManifest, candidateCleanErr := runCleanWorkspace(ctx, source.SiteDir, cfg.ParentDir, excluded, cfg.Candidate, cfg.Classes, gomaxprocs)

	correctness := checkpoint.Correctness
	if baselineObservation.Successful && candidateCleanObservation.Successful {
		correctness.DifferentialDiff = CompareManifests(baselineManifest, candidateCleanManifest, cfg.Comparators)
		correctness.DifferentialEqual = correctness.DifferentialDiff.Equal()
	}
	if incrementalApplicable && incremental.Successful() && candidateCleanObservation.Successful {
		correctness.IncrementalDiff = CompareManifests(candidateCleanManifest, incrementalManifest, cfg.Comparators)
		correctness.IncrementalEqual = correctness.IncrementalDiff.Equal()
	}
	correctness.DeterministicEqual = !cfg.CheckDeterminism
	var determinismObservation *RunObservation
	var determinismErr error
	if cfg.CheckDeterminism {
		observation, determinismManifest, err := runCleanWorkspace(ctx, source.SiteDir, cfg.ParentDir, excluded, cfg.Candidate, cfg.Classes, gomaxprocs)
		determinismObservation = &observation
		determinismErr = err
		if candidateCleanObservation.Successful && observation.Successful {
			correctness.DeterminismDiff = CompareManifests(candidateCleanManifest, determinismManifest, cfg.Comparators)
			correctness.DeterministicEqual = correctness.DeterminismDiff.Equal()
		}
	}
	checkpoint = Checkpoint{
		OperationIndex: operationIndex, Operation: string(operation.Type),
		Baseline: baselineObservation, CandidateClean: candidateCleanObservation,
		CandidateDeterminism: determinismObservation, CandidateIncremental: incrementalObservation,
		Correctness: correctness,
	}
	checkpoint.Diagnostics = checkpointDiagnostics(checkpoint)
	return checkpoint, errors.Join(baselineErr, candidateCleanErr, determinismErr)
}

func observationFromError(err error) RunObservation {
	return RunObservation{FailureClass: FailureHarness, Error: errorString(err)}
}

func runCleanWorkspace(
	ctx context.Context,
	fixture, parent string,
	excluded []string,
	command BuildCommand,
	classes map[string]OutputClass,
	gomaxprocs int,
) (observation RunObservation, manifest Manifest, runErr error) {
	workspace, err := newWorkspace(fixture, parent, excluded)
	if err != nil {
		return observationFromError(err), manifest, err
	}
	defer func() {
		runErr = errors.Join(runErr, cleanupWorkspace(workspace))
	}()
	run := command.run(ctx, workspace, true, gomaxprocs)
	var manifestErr error
	if run.Successful() {
		manifest, manifestErr = manifestForRun(workspace, command, classes)
	}
	observation = observe(run, manifest, manifestErr)
	return observation, manifest, nil
}

func isSemanticMutation(operation OperationType) bool {
	switch operation {
	case OpWriteFile, OpReplaceExact, OpDelete, OpRename, OpCopy, OpSetConfig, OpTouch:
		return true
	default:
		return false
	}
}

func cleanupWorkspace(w Workspace) error {
	if err := w.Remove(); err != nil {
		return fmt.Errorf("remove Build Lab workspace %q: %w", w.Root, err)
	}
	return nil
}

func applyWorkspaceOperation(w Workspace, operation Operation, outputCommand BuildCommand) error {
	if operation.Type == OpClearOutput {
		output, err := outputCommand.outputPath(w)
		if err != nil {
			return err
		}
		if err := os.RemoveAll(output); err != nil {
			return fmt.Errorf("clear workspace output: %w", err)
		}
		return nil
	}
	if err := ApplyOperation(w.SiteDir, operation); err != nil {
		return err
	}
	if (operation.Type != OpCleanCache && operation.Type != opLegacyClearCache) || w.MarkataCache == "" {
		return nil
	}
	if err := os.RemoveAll(w.MarkataCache); err != nil {
		return fmt.Errorf("clear workspace cache: %w", err)
	}
	if err := os.MkdirAll(w.MarkataCache, 0o755); err != nil {
		return fmt.Errorf("recreate workspace cache: %w", err)
	}
	return nil
}

func identityForCommand(command BuildCommand) BuildIdentity {
	identity := BuildIdentity{Binary: command.Binary}
	if hash, err := hashFile(command.Binary); err == nil {
		identity.BinarySHA256 = hash
	}
	return identity
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func manifestForRun(w Workspace, command BuildCommand, classes map[string]OutputClass) (Manifest, error) {
	output, err := command.outputPath(w)
	if err != nil {
		return Manifest{Records: []FileRecord{}}, err
	}
	if info, statErr := os.Stat(output); statErr != nil {
		return Manifest{Records: []FileRecord{}}, fmt.Errorf("output directory %q: %w", output, statErr)
	} else if !info.IsDir() {
		return Manifest{Records: []FileRecord{}}, fmt.Errorf("output path %q is not a directory", output)
	}
	return BuildManifest(output, classes)
}

func outputArg(args []string) string {
	output := ""
	for i, arg := range args {
		if (arg == "--output" || arg == "-o") && i+1 < len(args) {
			output = args[i+1]
		}
		if strings.HasPrefix(arg, "--output=") {
			output = strings.TrimPrefix(arg, "--output=")
		}
	}
	return output
}

func observe(run RunResult, manifest Manifest, manifestErr error) RunObservation {
	digest := ""
	if value, err := manifest.Digest(); err == nil {
		digest = value
	}
	commandErr := run.Err
	if manifestErr != nil {
		if run.Err == nil {
			run.Err = manifestErr
		} else {
			run.Err = errors.Join(run.Err, manifestErr)
		}
	}
	failureClass := FailureClass("")
	// A command that starts and returns a non-zero status is the product under
	// measurement. Start, containment, timeout, truncation, and manifest
	// collection errors belong to the harness.
	switch {
	case run.StdoutTruncated || run.StderrTruncated || run.TimedOut:
		failureClass = FailureHarness
	case commandErr != nil && run.ExitCode != 0:
		failureClass = FailureProduct
	case commandErr != nil || manifestErr != nil:
		failureClass = FailureHarness
	}
	return RunObservation{
		ExitCode: run.ExitCode, Successful: run.Successful(), TimedOut: run.TimedOut,
		StdoutTruncated: run.StdoutTruncated, StderrTruncated: run.StderrTruncated,
		DurationSeconds: run.Duration.Seconds(), ManifestDigest: digest,
		StdoutSHA256: hashBytes(run.Stdout), StderrSHA256: hashBytes(run.Stderr),
		FailureClass: failureClass, Error: errorString(run.Err),
	}
}

func checkpointDiagnostics(checkpoint Checkpoint) []Diagnostic {
	diagnostics := make([]Diagnostic, 0, 4)
	observations := []struct {
		scope string
		value RunObservation
	}{
		{scope: "baseline", value: checkpoint.Baseline},
		{scope: "candidate-clean", value: checkpoint.CandidateClean},
		{scope: "candidate-incremental", value: checkpoint.CandidateIncremental},
	}
	if checkpoint.CandidateDeterminism != nil {
		observations = append(observations, struct {
			scope string
			value RunObservation
		}{scope: "candidate-determinism", value: *checkpoint.CandidateDeterminism})
	}
	for _, observation := range observations {
		if observation.value.FailureClass == "" {
			continue
		}
		message := observation.scope + " build failed"
		if observation.value.Error != "" {
			message += ": " + observation.value.Error
		}
		diagnostics = append(diagnostics, Diagnostic{
			Class:   observation.value.FailureClass,
			Scope:   observation.scope,
			Message: message,
		})
	}
	if !checkpoint.Correctness.DifferentialEqual && checkpoint.Baseline.Successful && checkpoint.CandidateClean.Successful {
		diff := checkpoint.Correctness.DifferentialDiff
		diagnostics = append(diagnostics, Diagnostic{
			Class:   FailureProduct,
			Scope:   "clean-comparison",
			Message: "clean baseline and candidate output differ (" + diffSummary(diff) + ")",
			Diff:    &diff,
		})
	}
	if checkpoint.Correctness.IncrementalApplicable && !checkpoint.Correctness.IncrementalEqual && checkpoint.CandidateClean.Successful && checkpoint.CandidateIncremental.Successful {
		diff := checkpoint.Correctness.IncrementalDiff
		diagnostics = append(diagnostics, Diagnostic{
			Class:   FailureProduct,
			Scope:   "incremental-comparison",
			Message: "candidate incremental and clean output differ (" + diffSummary(diff) + ")",
			Diff:    &diff,
		})
	}
	if !checkpoint.Correctness.DeterministicEqual && checkpoint.CandidateClean.Successful && checkpoint.CandidateDeterminism != nil && checkpoint.CandidateDeterminism.Successful {
		diff := checkpoint.Correctness.DeterminismDiff
		diagnostics = append(diagnostics, Diagnostic{
			Class:   FailureProduct,
			Scope:   "determinism-comparison",
			Message: "candidate clean rebuilds differ (" + diffSummary(diff) + ")",
			Diff:    &diff,
		})
	}
	return diagnostics
}

func diffSummary(diff ManifestDiff) string {
	return fmt.Sprintf("missing=%d, extra=%d, changed=%d", len(diff.Missing), len(diff.Extra), len(diff.Changed))
}

func resultFailureClass(result Result) FailureClass {
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Class == FailureHarness {
			return FailureHarness
		}
	}
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Class == FailureProduct {
			return FailureProduct
		}
	}
	for checkpointIndex := range result.Checkpoints {
		checkpoint := &result.Checkpoints[checkpointIndex]
		for _, diagnostic := range checkpoint.Diagnostics {
			if diagnostic.Class == FailureHarness {
				return FailureHarness
			}
		}
	}
	for checkpointIndex := range result.Checkpoints {
		checkpoint := &result.Checkpoints[checkpointIndex]
		for _, diagnostic := range checkpoint.Diagnostics {
			if diagnostic.Class == FailureProduct {
				return FailureProduct
			}
		}
	}
	return ""
}

func hashBytes(value []byte) string {
	h := sha256.Sum256(value)
	return hex.EncodeToString(h[:])
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// WriteResult writes canonical result JSON and creates its parent directory.
func WriteResult(path string, result Result) error {
	data, err := result.CanonicalJSON()
	if err != nil {
		return err
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}
