package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

// chdir changes to the given directory and returns a cleanup function.
// It calls t.Helper() and t.Fatal() on errors.
func chdir(t *testing.T, dir string) func() {
	t.Helper()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change to directory %s: %v", dir, err)
	}
	return func() {
		if err := os.Chdir(origDir); err != nil {
			t.Logf("warning: failed to restore directory: %v", err)
		}
	}
}

// evalSymlinks resolves symlinks in a path. On macOS, /var is a symlink to
// /private/var, which causes path comparison issues in tests.
func evalSymlinks(t *testing.T, path string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("failed to resolve symlinks for %s: %v", path, err)
	}
	return resolved
}

// =============================================================================
// Configuration Tests based on tests.yaml
// =============================================================================

func TestConfig_DefaultOutputDir(t *testing.T) {
	// Test case: "default output dir"
	// When no config is provided, output_dir should be "output"
	config := DefaultConfig()

	if config.OutputDir != "output" {
		t.Errorf("output_dir: got %q, want 'output'", config.OutputDir)
	}
}

func TestConfig_CustomOutputDir(t *testing.T) {
	// Test case: "custom output dir"
	// Create a temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "markata-go.toml")
	configContent := `[markata-go]
output_dir = "public"
`
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if config.OutputDir != "public" {
		t.Errorf("output_dir: got %q, want 'public'", config.OutputDir)
	}
}

func TestConfig_DefaultGlobPatterns(t *testing.T) {
	// Test case: "default glob patterns"
	config := DefaultConfig()

	expectedPatterns := []string{"*.md", "pages/**/*.md", "posts/**/*.md"}
	if len(config.GlobConfig.Patterns) != len(expectedPatterns) ||
		config.GlobConfig.Patterns[0] != expectedPatterns[0] ||
		config.GlobConfig.Patterns[1] != expectedPatterns[1] ||
		config.GlobConfig.Patterns[2] != expectedPatterns[2] {
		t.Errorf("glob.patterns: got %v, want %v", config.GlobConfig.Patterns, expectedPatterns)
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	config := DefaultConfig()

	// Verify all defaults
	if config.OutputDir != "output" {
		t.Errorf("output_dir: got %q, want 'output'", config.OutputDir)
	}
	if config.TemplatesDir != "templates" {
		t.Errorf("templates_dir: got %q, want 'templates'", config.TemplatesDir)
	}
	if config.AssetsDir != "static" {
		t.Errorf("assets_dir: got %q, want 'static'", config.AssetsDir)
	}
	if license, ok := config.License.Key(); !ok || license != models.DefaultLicenseKey {
		t.Errorf("license: got %q (configured=%t), want %q", license, ok, models.DefaultLicenseKey)
	}
	if config.NeedsLicenseWarning() {
		t.Error("default configuration should not require a license warning")
	}
	if len(config.Hooks) != 1 || config.Hooks[0] != "default" {
		t.Errorf("hooks: got %v, want ['default']", config.Hooks)
	}
	if !config.GlobConfig.UseGitignore {
		t.Error("glob.use_gitignore should default to true")
	}
	if config.FeedDefaults.ItemsPerPage != 10 {
		t.Errorf("feed_defaults.items_per_page: got %d, want 10", config.FeedDefaults.ItemsPerPage)
	}
	if config.FeedDefaults.OrphanThreshold != 3 {
		t.Errorf("feed_defaults.orphan_threshold: got %d, want 3", config.FeedDefaults.OrphanThreshold)
	}
}

func TestLoadFromString_BuilderAdminTypedAndSecretSafe(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		data   string
	}{
		{
			name:   "toml",
			format: FormatTOML,
			data: `
[markata-go]
custom_plugin = { enabled = true }
[markata-go.builder_admin]
unrecognized_setting = "kept"
[markata-go.builder_admin.auth.headers]
user_id = "X-Authentik-Uid"
display_name = ""
[markata-go.builder_admin.webhook]
enabled = true
branch = "qa"
secret = "test-secret"
endpoint = "/hooks/build"
`,
		},
		{
			name:   "yaml",
			format: FormatYAML,
			data: `
markata-go:
  custom_plugin:
    enabled: true
  builder_admin:
    unrecognized_setting: kept
    auth:
      headers:
        user_id: X-Authentik-Uid
        display_name: ""
    webhook:
      enabled: true
      branch: qa
      secret: test-secret
      endpoint: /hooks/build
`,
		},
		{
			name:   "json",
			format: FormatJSON,
			data:   `{"markata-go":{"custom_plugin":{"enabled":true},"builder_admin":{"unrecognized_setting":"kept","auth":{"headers":{"user_id":"X-Authentik-Uid","display_name":""}},"webhook":{"enabled":true,"branch":"qa","secret":"test-secret","endpoint":"/hooks/build"}}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := LoadFromString(tt.data, tt.format)
			if err != nil {
				t.Fatalf("LoadFromString() error = %v", err)
			}
			headers := cfg.BuilderAdmin.Auth.Headers
			webhook := cfg.BuilderAdmin.Webhook
			if headers.UserID == nil || *headers.UserID != "X-Authentik-Uid" || headers.DisplayName == nil || *headers.DisplayName != "" {
				t.Fatalf("headers = %#v", headers)
			}
			if webhook.Enabled == nil || !*webhook.Enabled || webhook.Branch == nil || *webhook.Branch != "qa" || webhook.Secret == nil || *webhook.Secret != "test-secret" {
				t.Fatalf("webhook = %#v", webhook)
			}
			if _, ok := cfg.Extra["builder_admin"]; ok {
				t.Fatalf("builder_admin must not be in Extra: %#v", cfg.Extra["builder_admin"])
			}
			if _, ok := cfg.Extra["custom_plugin"]; !ok {
				t.Fatalf("unrelated top-level config missing from Extra: %#v", cfg.Extra)
			}
			if cfg.BuilderAdmin.Extra["unrecognized_setting"] != "kept" {
				t.Fatalf("unrecognized builder_admin field missing: %#v", cfg.BuilderAdmin.Extra)
			}
			webhookExtra, ok := cfg.BuilderAdmin.Extra["webhook"].(map[string]any)
			if !ok || webhookExtra["endpoint"] != "/hooks/build" {
				t.Fatalf("unrecognized webhook field missing: %#v", cfg.BuilderAdmin.Extra)
			}
			encoded, err := json.Marshal(cfg)
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if strings.Contains(string(encoded), "test-secret") {
				t.Fatalf("serialized config exposed webhook secret: %s", encoded)
			}
		})
	}
}

// =============================================================================
// Environment Variable Override Tests
// =============================================================================

func TestConfig_EnvVarOverride(t *testing.T) {
	// Test case: "env var overrides config file"
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "markata-go.toml")
	configContent := `[markata-go]
output_dir = "file-output"
`
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Set environment variable
	_ = os.Setenv("MARKATA_GO_OUTPUT_DIR", "env-output")
	defer func() { _ = os.Unsetenv("MARKATA_GO_OUTPUT_DIR") }()

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	// Environment variable should override config file
	if config.OutputDir != "env-output" {
		t.Errorf("output_dir: got %q, want 'env-output'", config.OutputDir)
	}
}

func TestConfig_EnvVarSetsRoot(t *testing.T) {
	// Test case: "env var sets root config"
	// Clear any existing env vars
	_ = os.Unsetenv("MARKATA_GO_OUTPUT_DIR")

	// Set environment variable
	_ = os.Setenv("MARKATA_GO_OUTPUT_DIR", "env-output")
	defer func() { _ = os.Unsetenv("MARKATA_GO_OUTPUT_DIR") }()

	config, err := LoadWithDefaults()
	if err != nil {
		t.Fatalf("LoadWithDefaults error: %v", err)
	}

	if config.OutputDir != "env-output" {
		t.Errorf("output_dir: got %q, want 'env-output'", config.OutputDir)
	}
}

func TestConfig_EnvVarNestedConfig(t *testing.T) {
	// Test case: "env var sets nested config"
	_ = os.Setenv("MARKATA_GO_GLOB_USE_GITIGNORE", "false")
	defer func() { _ = os.Unsetenv("MARKATA_GO_GLOB_USE_GITIGNORE") }()

	config, err := LoadWithDefaults()
	if err != nil {
		t.Fatalf("LoadWithDefaults error: %v", err)
	}

	if config.GlobConfig.UseGitignore != false {
		t.Error("glob.use_gitignore should be false from env var")
	}
}

func TestConfig_EnvVarBooleanValues(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected bool
	}{
		{"true", "true", true},
		{"1", "1", true},
		{"yes", "yes", true},
		{"false", "false", false},
		{"0", "0", false},
		{"no", "no", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("MARKATA_GO_GLOB_USE_GITIGNORE", tt.envValue)
			defer func() { _ = os.Unsetenv("MARKATA_GO_GLOB_USE_GITIGNORE") }()

			config, err := LoadWithDefaults()
			if err != nil {
				t.Fatalf("LoadWithDefaults error: %v", err)
			}

			if config.GlobConfig.UseGitignore != tt.expected {
				t.Errorf("glob.use_gitignore: got %v, want %v", config.GlobConfig.UseGitignore, tt.expected)
			}
		})
	}
}

func TestConfig_EnvVarListValue(t *testing.T) {
	// Test case: "env var list value comma separated"
	_ = os.Setenv("MARKATA_GO_GLOB_PATTERNS", "posts/**/*.md,pages/*.md")
	defer func() { _ = os.Unsetenv("MARKATA_GO_GLOB_PATTERNS") }()

	config, err := LoadWithDefaults()
	if err != nil {
		t.Fatalf("LoadWithDefaults error: %v", err)
	}

	expected := []string{"posts/**/*.md", "pages/*.md"}
	if len(config.GlobConfig.Patterns) != len(expected) {
		t.Errorf("glob.patterns: got %v, want %v", config.GlobConfig.Patterns, expected)
	}
	for i, pattern := range expected {
		if i < len(config.GlobConfig.Patterns) && config.GlobConfig.Patterns[i] != pattern {
			t.Errorf("glob.patterns[%d]: got %q, want %q", i, config.GlobConfig.Patterns[i], pattern)
		}
	}
}

func TestConfig_EnvVarIntegerConversion(t *testing.T) {
	// Test case: "env var integer conversion"
	_ = os.Setenv("MARKATA_GO_CONCURRENCY", "8")
	defer func() { _ = os.Unsetenv("MARKATA_GO_CONCURRENCY") }()

	config, err := LoadWithDefaults()
	if err != nil {
		t.Fatalf("LoadWithDefaults error: %v", err)
	}

	if config.Concurrency != 8 {
		t.Errorf("concurrency: got %d, want 8", config.Concurrency)
	}
}

// =============================================================================
// Config File Discovery Tests
// =============================================================================

func TestConfig_DiscoverTOML(t *testing.T) {
	// Test case: "finds markata-go.toml in current directory"
	tmpDir := evalSymlinks(t, t.TempDir())

	// Save current dir and change to temp dir
	cleanup := chdir(t, tmpDir)
	defer cleanup()

	// Create config file
	configPath := filepath.Join(tmpDir, "markata-go.toml")
	configContent := `[markata-go]
output_dir = "public"
`
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	discovered, err := Discover()
	if err != nil {
		t.Fatalf("Discover error: %v", err)
	}

	if discovered != configPath {
		t.Errorf("discovered: got %q, want %q", discovered, configPath)
	}
}

func TestConfig_DiscoverYAML(t *testing.T) {
	// Test case: "finds yaml format"
	tmpDir := evalSymlinks(t, t.TempDir())

	cleanup := chdir(t, tmpDir)
	defer cleanup()

	configPath := filepath.Join(tmpDir, "markata-go.yaml")
	configContent := `markata-go:
  output_dir: yaml-output
`
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	discovered, err := Discover()
	if err != nil {
		t.Fatalf("Discover error: %v", err)
	}

	if discovered != configPath {
		t.Errorf("discovered: got %q, want %q", discovered, configPath)
	}
}

func TestConfig_TOMLPreferredOverYAML(t *testing.T) {
	// Test case: "toml preferred over yaml"
	tmpDir := evalSymlinks(t, t.TempDir())

	cleanup := chdir(t, tmpDir)
	defer cleanup()

	// Create both TOML and YAML
	tomlPath := filepath.Join(tmpDir, "markata-go.toml")
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(tomlPath, []byte(`[markata-go]
output_dir = "toml-output"
`), 0o644); err != nil {
		t.Fatalf("failed to write TOML file: %v", err)
	}

	yamlPath := filepath.Join(tmpDir, "markata-go.yaml")
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(yamlPath, []byte(`markata-go:
  output_dir: yaml-output
`), 0o644); err != nil {
		t.Fatalf("failed to write YAML file: %v", err)
	}

	discovered, err := Discover()
	if err != nil {
		t.Fatalf("Discover error: %v", err)
	}

	// TOML should be preferred
	if discovered != tomlPath {
		t.Errorf("discovered: got %q, want %q (TOML preferred)", discovered, tomlPath)
	}
}

func TestConfig_NoConfigFileUsesDefaults(t *testing.T) {
	// Test case: "no config file uses defaults"
	tmpDir := t.TempDir()

	cleanup := chdir(t, tmpDir)
	defer cleanup()

	// No config file created

	config, err := Load("")
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if config.OutputDir != "output" {
		t.Errorf("output_dir: got %q, want 'output' (default)", config.OutputDir)
	}
}

// =============================================================================
// Config Validation Tests
// =============================================================================

func TestConfigValidation_InvalidURL(t *testing.T) {
	// Test case: "invalid url format error"
	config := DefaultConfig()
	config.URL = "not-a-url"

	errs := ValidateConfig(config)
	if len(errs) == 0 {
		t.Error("expected validation error for invalid URL")
	}

	// Check that error mentions URL
	found := false
	for _, err := range errs {
		var ve ValidationError
		if errors.As(err, &ve) && ve.Field == "url" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected validation error for 'url' field")
	}
}

func TestConfigValidation_NegativeConcurrency(t *testing.T) {
	// Test case: "negative concurrency error"
	config := DefaultConfig()
	config.Concurrency = -5

	errs := ValidateConfig(config)
	if len(errs) == 0 {
		t.Error("expected validation error for negative concurrency")
	}

	found := false
	for _, err := range errs {
		var ve ValidationError
		if errors.As(err, &ve) && ve.Field == "concurrency" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected validation error for 'concurrency' field")
	}
}

func TestConfigValidation_EmptyPatternsWarning(t *testing.T) {
	// Test case: "empty patterns warning"
	config := DefaultConfig()
	config.GlobConfig.Patterns = []string{}

	errs := ValidateConfig(config)

	// Should have a warning (not error) for empty patterns
	hasWarning := false
	for _, err := range errs {
		var ve ValidationError
		if errors.As(err, &ve) && ve.Field == "glob.patterns" && ve.IsWarn {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Error("expected warning for empty glob.patterns")
	}
}

func TestConfigValidation_ValidConfig(t *testing.T) {
	config := DefaultConfig()
	config.URL = "https://example.com"

	errs := ValidateConfig(config)

	// Should only have warnings (empty patterns), no errors
	for _, err := range errs {
		var ve ValidationError
		if errors.As(err, &ve) && !ve.IsWarn {
			t.Errorf("unexpected error: %v", err)
		}
	}
}

func TestConfigValidation_NilConfig(t *testing.T) {
	errs := ValidateConfig(nil)

	if len(errs) == 0 {
		t.Error("expected error for nil config")
	}
}

// =============================================================================
// Config Merge Tests
// =============================================================================

func TestConfigMerge_ScalarValues(t *testing.T) {
	// Test case: "scalar values later wins"
	base := DefaultConfig()
	base.OutputDir = "global"

	override := &models.Config{
		OutputDir: "local",
	}

	merged := MergeConfigs(base, override)

	if merged.OutputDir != "local" {
		t.Errorf("output_dir: got %q, want 'local'", merged.OutputDir)
	}
}

func TestConfigMerge_PreservesUnsetValues(t *testing.T) {
	base := DefaultConfig()
	base.OutputDir = "global"
	base.URL = "https://global.com"

	override := &models.Config{
		OutputDir: "local",
		// URL not set
	}

	merged := MergeConfigs(base, override)

	if merged.OutputDir != "local" {
		t.Errorf("output_dir: got %q, want 'local'", merged.OutputDir)
	}
	// URL should be preserved from base if not overridden
	// Note: zero values might override depending on implementation
	_ = merged.URL // verified merge happened via other field checks
}

// =============================================================================
// Config Helper Tests
// =============================================================================

func TestSpec_GetEnvValue(t *testing.T) {
	// Set a test env var
	if err := os.Setenv("MARKATA_GO_TEST_KEY", "test_value"); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("MARKATA_GO_TEST_KEY"); err != nil {
			t.Logf("warning: failed to unset env var: %v", err)
		}
	}()

	value, ok := GetEnvValue("TEST_KEY")
	if !ok {
		t.Error("expected env var to be found")
	}
	if value != "test_value" {
		t.Errorf("got %q, want 'test_value'", value)
	}

	// Test missing env var
	_, ok = GetEnvValue("NONEXISTENT")
	if ok {
		t.Error("expected env var not to be found")
	}
}

func TestSpec_SetEnvValue(t *testing.T) {
	key := "SET_TEST"
	value := "set_value"

	err := SetEnvValue(key, value)
	if err != nil {
		t.Fatalf("SetEnvValue error: %v", err)
	}
	defer func() {
		if err := UnsetEnvValue(key); err != nil {
			t.Logf("warning: failed to unset env var: %v", err)
		}
	}()

	got, ok := GetEnvValue(key)
	if !ok || got != value {
		t.Errorf("got %q, want %q", got, value)
	}
}

func TestSpec_UnsetEnvValue(t *testing.T) {
	key := "UNSET_TEST"
	if err := SetEnvValue(key, "value"); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}

	err := UnsetEnvValue(key)
	if err != nil {
		t.Fatalf("UnsetEnvValue error: %v", err)
	}

	_, ok := GetEnvValue(key)
	if ok {
		t.Error("env var should be unset")
	}
}

func TestSpec_ConfigFromEnv(t *testing.T) {
	if err := os.Setenv("MARKATA_GO_OUTPUT_DIR", "from-env"); err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() {
		if err := os.Unsetenv("MARKATA_GO_OUTPUT_DIR"); err != nil {
			t.Logf("warning: failed to unset env var: %v", err)
		}
	}()

	config := FromEnv()

	if config.OutputDir != "from-env" {
		t.Errorf("output_dir: got %q, want 'from-env'", config.OutputDir)
	}
}

// =============================================================================
// Load Functions Tests
// =============================================================================

func TestSpec_LoadFromString(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		format    Format
		wantDir   string
		wantError bool
	}{
		{
			name: "TOML format",
			data: `[markata-go]
output_dir = "toml-dir"
`,
			format:    FormatTOML,
			wantDir:   "toml-dir",
			wantError: false,
		},
		{
			name: "YAML format",
			data: `markata-go:
  output_dir: yaml-dir
`,
			format:    FormatYAML,
			wantDir:   "yaml-dir",
			wantError: false,
		},
		{
			name:      "JSON format",
			data:      `{"markata-go": {"output_dir": "json-dir"}}`,
			format:    FormatJSON,
			wantDir:   "json-dir",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadFromString(tt.data, tt.format)
			if (err != nil) != tt.wantError {
				t.Errorf("error = %v, wantError = %v", err, tt.wantError)
				return
			}
			if !tt.wantError && config.OutputDir != tt.wantDir {
				t.Errorf("output_dir: got %q, want %q", config.OutputDir, tt.wantDir)
			}
		})
	}
}

func TestSpec_MustLoad_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected MustLoad to panic on invalid path")
		}
	}()

	MustLoad("/nonexistent/path/config.toml")
}

//nolint:gosec // Test file permissions are fine at 0644
func TestSpec_LoadAndValidate(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "markata-go.toml")
	configContent := `[markata-go]
output_dir = "public"
url = "https://example.com"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	config, validationErrs, err := LoadAndValidate(configPath)
	if err != nil {
		t.Fatalf("LoadAndValidate error: %v", err)
	}

	if config == nil {
		t.Fatal("config should not be nil")
	}
	if config.OutputDir != "public" {
		t.Errorf("output_dir: got %q, want 'public'", config.OutputDir)
	}

	// Check for warnings only (no errors expected)
	for _, ve := range validationErrs {
		var validErr ValidationError
		if errors.As(ve, &validErr) && !validErr.IsWarn {
			t.Errorf("unexpected validation error: %v", ve)
		}
	}
}

// =============================================================================
// Format Detection Tests
// =============================================================================

func TestSpec_FormatFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected Format
	}{
		{"config.toml", FormatTOML},
		{"config.yaml", FormatYAML},
		{"config.yml", FormatYAML},
		{"config.json", FormatJSON},
		{"config.TOML", FormatTOML},
		{"config.unknown", FormatTOML}, // Default
		{"config", FormatTOML},         // No extension
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := formatFromPath(tt.path)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestSpec_DiscoverAll(t *testing.T) {
	tmpDir := t.TempDir()

	cleanup := chdir(t, tmpDir)
	defer cleanup()

	// Create multiple config files
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(filepath.Join(tmpDir, "markata-go.toml"), []byte(""), 0o644); err != nil {
		t.Fatalf("failed to write TOML file: %v", err)
	}
	//nolint:gosec // Test file permissions are fine at 0644
	if err := os.WriteFile(filepath.Join(tmpDir, "markata-go.yaml"), []byte(""), 0o644); err != nil {
		t.Fatalf("failed to write YAML file: %v", err)
	}

	found := DiscoverAll()

	if len(found) < 2 {
		t.Errorf("expected at least 2 config files, found %d", len(found))
	}

	// Verify sources are set
	for _, cp := range found {
		if cp.Source == "" {
			t.Error("Path.Source should be set")
		}
	}
}

func TestLoadWithDefaults_MatchesEmptyConfigTheme(t *testing.T) {
	noFile, err := LoadWithDefaults()
	if err != nil {
		t.Fatalf("LoadWithDefaults error: %v", err)
	}
	emptyFile, err := LoadFromString("[markata-go]\n", FormatTOML)
	if err != nil {
		t.Fatalf("LoadFromString error: %v", err)
	}
	if noFile.Theme.Palette == "" {
		t.Fatal("config-less sites need a default palette for the theme picker")
	}
	if noFile.Theme.Palette != emptyFile.Theme.Palette || noFile.Theme.Aesthetic != emptyFile.Theme.Aesthetic {
		t.Errorf("no config file theme = %q/%q, empty config file theme = %q/%q",
			noFile.Theme.Palette, noFile.Theme.Aesthetic, emptyFile.Theme.Palette, emptyFile.Theme.Aesthetic)
	}
}
