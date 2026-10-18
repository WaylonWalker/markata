package cmd

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"string", "hello", "hello"},
		{"string with spaces", "hello world", "hello world"},
		{"bool true", "true", true},
		{"bool false", "false", false},
		{"bool TRUE", "TRUE", true},
		{"bool FALSE", "FALSE", false},
		{"int", "42", 42},
		{"int negative", "-5", -5},
		{"float", "3.14", 3.14},
		{"json array", `["a", "b", "c"]`, []interface{}{"a", "b", "c"}},
		{"json object", `{"key": "value"}`, map[string]interface{}{"key": "value"}},
		{"string that looks like number", "007", 7}, // Note: leading zeros are parsed as int
		{"url", "https://example.com", "https://example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseValue(tt.input)

			// For slices and maps, compare as JSON strings
			switch expected := tt.expected.(type) {
			case []interface{}:
				gotSlice, ok := got.([]interface{})
				if !ok {
					t.Errorf("parseValue(%q) = %T, want []interface{}", tt.input, got)
					return
				}
				if len(gotSlice) != len(expected) {
					t.Errorf("parseValue(%q) = %v, want %v", tt.input, got, tt.expected)
				}
			case map[string]interface{}:
				gotMap, ok := got.(map[string]interface{})
				if !ok {
					t.Errorf("parseValue(%q) = %T, want map[string]interface{}", tt.input, got)
					return
				}
				if len(gotMap) != len(expected) {
					t.Errorf("parseValue(%q) = %v, want %v", tt.input, got, tt.expected)
				}
			default:
				if got != tt.expected {
					t.Errorf("parseValue(%q) = %v (%T), want %v (%T)", tt.input, got, got, tt.expected, tt.expected)
				}
			}
		})
	}
}

func TestFormatFromPath(t *testing.T) {
	tests := []struct {
		path   string
		expect string
	}{
		{"config.toml", formatTOML},
		{"config.yaml", formatYAML},
		{"config.yml", formatYAML},
		{"config.json", formatJSON},
		{"config.txt", formatTOML}, // default
		{"config.TOML", formatTOML},
		{"config.YAML", formatYAML},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := formatFromPath(tt.path)
			if got != tt.expect {
				t.Errorf("formatFromPath(%q) = %q, want %q", tt.path, got, tt.expect)
			}
		})
	}
}

func TestParseConfigToMap(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		format  string
		wantErr bool
	}{
		{
			name: "toml config",
			data: `[markata-go]
title = "Test"
`,
			format:  formatTOML,
			wantErr: false,
		},
		{
			name: "yaml config",
			data: `markata-go:
  title: Test
`,
			format:  formatYAML,
			wantErr: false,
		},
		{
			name:    "json config",
			data:    `{"markata-go": {"title": "Test"}}`,
			format:  formatJSON,
			wantErr: false,
		},
		{
			name:    "invalid toml",
			data:    `[invalid`,
			format:  formatTOML,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigToMap([]byte(tt.data), tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfigToMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigSetIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name           string
		initialContent string
		format         string
		key            string
		value          string
		checkContains  string
	}{
		{
			name: "set toml string",
			initialContent: `[markata-go]
title = "Old"
`,
			format:        formatTOML,
			key:           "title",
			value:         "New Title",
			checkContains: `title = "New Title"`,
		},
		{
			name: "set toml int",
			initialContent: `[markata-go]
concurrency = 2
`,
			format:        formatTOML,
			key:           "concurrency",
			value:         "8",
			checkContains: "concurrency = 8",
		},
		{
			name: "set yaml string",
			initialContent: `markata-go:
  title: Old
`,
			format:        formatYAML,
			key:           "title",
			value:         "New Title",
			checkContains: "title: New Title",
		},
		{
			name:           "set json bool",
			initialContent: `{"markata-go": {"verbose": false}}`,
			format:         formatJSON,
			key:            "verbose",
			value:          "true",
			checkContains:  `"verbose": true`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			ext := "." + tt.format
			if tt.format == formatYAML {
				ext = ".yaml"
			}
			configPath := filepath.Join(tmpDir, "config"+ext)
			if err := os.WriteFile(configPath, []byte(tt.initialContent), 0o600); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			parsedValue := parseValue(tt.value)
			if err := config.SetValueInFile(configPath, tt.key, parsedValue); err != nil {
				t.Fatalf("Failed to set value: %v", err)
			}

			result, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("Failed to read result: %v", err)
			}
			if !strings.Contains(string(result), tt.checkContains) {
				t.Errorf("Result does not contain %q:\n%s", tt.checkContains, result)
			}
		})
	}
}

func TestFormatValueForDisplay(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		expect string
	}{
		{"nil", nil, "<not set>"},
		{"string", "hello", `"hello"`},
		{"int", 42, "42"},
		{"bool", true, "true"},
		{"slice", []interface{}{"a", "b"}, `["a","b"]`},
		{"map", map[string]interface{}{"key": "value"}, `{"key":"value"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatValueForDisplay(tt.value)
			if got != tt.expect {
				t.Errorf("formatValueForDisplay(%v) = %q, want %q", tt.value, got, tt.expect)
			}
		})
	}
}

func TestConfigToMap(t *testing.T) {
	cfg := config.DefaultConfig()
	m, err := configToMap(cfg)
	if err != nil {
		t.Fatalf("configToMap() error = %v", err)
	}

	// Check some expected keys exist
	if _, ok := m["output_dir"]; !ok {
		t.Error("configToMap() missing output_dir key")
	}
	if _, ok := m["glob"]; !ok {
		t.Error("configToMap() missing glob key")
	}
}

func TestHasPath(t *testing.T) {
	tests := []struct {
		name   string
		m      map[string]interface{}
		path   []string
		expect bool
	}{
		{
			name:   "nil map",
			m:      nil,
			path:   []string{"foo"},
			expect: false,
		},
		{
			name:   "top level key exists",
			m:      map[string]interface{}{"title": "Test"},
			path:   []string{"title"},
			expect: true,
		},
		{
			name: "nested key exists",
			m: map[string]interface{}{
				"glob": map[string]interface{}{"patterns": []interface{}{"posts/**/*.md"}},
			},
			path:   []string{"glob", "patterns"},
			expect: true,
		},
		{
			name:   "key does not exist",
			m:      map[string]interface{}{"title": "Test"},
			path:   []string{"description"},
			expect: false,
		},
		{
			name:   "missing nested key",
			m:      map[string]interface{}{"glob": map[string]interface{}{}},
			path:   []string{"glob", "patterns"},
			expect: false,
		},
		{
			name:   "empty map",
			m:      map[string]interface{}{},
			path:   []string{"title"},
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasPath(tt.m, tt.path)
			if got != tt.expect {
				t.Errorf("hasPath() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestShowDiffConfig(t *testing.T) {
	tests := []struct {
		name           string
		userMap        map[string]interface{}
		configPaths    []string
		expectContains []string
	}{
		{
			name:           "nil user map",
			userMap:        nil,
			configPaths:    nil,
			expectContains: []string{"No user configuration found", "All values are defaults"},
		},
		{
			name:           "empty user map",
			userMap:        map[string]interface{}{},
			configPaths:    nil,
			expectContains: []string{"No user configuration found"},
		},
		{
			name:           "user config with values",
			userMap:        map[string]interface{}{"title": "My Site", "url": "https://example.com"},
			configPaths:    []string{"/path/to/markata-go.toml"},
			expectContains: []string{"User configuration from:", "title: My Site"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := bytes.NewBuffer(nil)
			command := &cobra.Command{Use: "config"}
			command.SetOut(stdout)
			currentCmd = command
			defer func() { currentCmd = nil }()

			err := showDiffConfig(tt.userMap, tt.configPaths)

			if err != nil {
				t.Fatalf("showDiffConfig() error = %v", err)
			}

			output := stdout.String()

			for _, expect := range tt.expectContains {
				if !strings.Contains(output, expect) {
					t.Errorf("showDiffConfig() output does not contain %q\nOutput:\n%s", expect, output)
				}
			}
		})
	}
}

func TestConfigShowAnnotateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	// Create a temporary directory with a config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "markata-go.toml")

	// Write a test config with some user values
	configContent := `[markata-go]
title = "My Test Site"
url = "https://example.com"
output_dir = "dist"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o600); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Change to tmp dir so config discovery works
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}
	defer func() {
		if err := os.Chdir(oldWd); err != nil {
			t.Logf("Warning: failed to restore working directory: %v", err)
		}
	}()

	merged, _, configPaths, err := loadManagerConfig("")
	if err != nil {
		t.Fatalf("loadManagerConfig() error = %v", err)
	}

	configShowDiff = true
	configShowAnnotate = false
	defer func() {
		configShowDiff = false
		configShowAnnotate = false
	}()

	stdout := bytes.NewBuffer(nil)
	command := &cobra.Command{Use: "config"}
	command.SetOut(stdout)
	currentCmd = command
	defer func() { currentCmd = nil }()

	err = runConfigShowWithSources(merged, configPaths, nil)

	if err != nil {
		t.Fatalf("runConfigShowWithSources() error = %v", err)
	}

	output := stdout.String()

	// Should contain user values
	if !strings.Contains(output, "title") {
		t.Errorf("Output should contain 'title', got:\n%s", output)
	}
	if !strings.Contains(output, "My Test Site") {
		t.Errorf("Output should contain 'My Test Site', got:\n%s", output)
	}
}

func TestResolveConfigShowFormat_RejectsConflictingFlags(t *testing.T) {
	command := &cobra.Command{Use: "show"}
	command.Flags().String("format", "yaml", "")
	command.Flags().Bool("json", false, "")
	command.Flags().Bool("toml", false, "")

	if err := command.Flags().Set("json", "true"); err != nil {
		t.Fatalf("set json flag: %v", err)
	}
	if err := command.Flags().Set("toml", "true"); err != nil {
		t.Fatalf("set toml flag: %v", err)
	}

	configFormat = "yaml"
	_, err := resolveConfigShowFormat(command)
	if err == nil {
		t.Fatal("expected conflicting flag error")
	}

	var exitErr *ExitCodeError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitCodeError, got %T", err)
	}
	if exitErr.ExitCode() != exitCodeUsage {
		t.Fatalf("exit code = %d, want %d", exitErr.ExitCode(), exitCodeUsage)
	}
}

func TestConfigToMap_YAMLSafeKeys(t *testing.T) {
	cfg := config.DefaultConfig()
	m, err := configToMap(cfg)
	if err != nil {
		t.Fatalf("configToMap() error = %v", err)
	}

	if _, err := yaml.Marshal(m); err != nil {
		t.Fatalf("yaml.Marshal(configToMap()) error = %v", err)
	}
}

func TestLoadUserConfigMap_MergesConfigPaths(t *testing.T) {
	tmpDir := t.TempDir()
	basePath := filepath.Join(tmpDir, "markata-go.toml")
	overridePath := filepath.Join(tmpDir, "fast.toml")

	base := `[markata-go]
title = "Base"

[markata-go.glob]
patterns = ["posts/**/*.md"]
`
	override := `[markata-go]
url = "https://example.com"

[markata-go.glob]
use_gitignore = false
`

	if err := os.WriteFile(basePath, []byte(base), 0o600); err != nil {
		t.Fatalf("write base config: %v", err)
	}
	if err := os.WriteFile(overridePath, []byte(override), 0o600); err != nil {
		t.Fatalf("write override config: %v", err)
	}

	userMap, err := loadUserConfigMap([]string{basePath, overridePath})
	if err != nil {
		t.Fatalf("loadUserConfigMap() error = %v", err)
	}

	if got, ok := getPathValue(userMap, []string{"title"}); !ok || got != "Base" {
		t.Fatalf("title = %v, %v; want Base, true", got, ok)
	}
	if got, ok := getPathValue(userMap, []string{"url"}); !ok || got != "https://example.com" {
		t.Fatalf("url = %v, %v; want https://example.com, true", got, ok)
	}
	if got, ok := getPathValue(userMap, []string{"glob", "use_gitignore"}); !ok || got != false {
		t.Fatalf("glob.use_gitignore = %v, %v; want false, true", got, ok)
	}
}

func TestExitCodeForError_UsageErrorReturnsTwo(t *testing.T) {
	if got := ExitCodeForError(errors.New("accepts 1 arg(s), received 0")); got != exitCodeUsage {
		t.Fatalf("ExitCodeForError() = %d, want %d", got, exitCodeUsage)
	}
}
