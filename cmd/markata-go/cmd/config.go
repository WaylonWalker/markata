package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Common string constants to avoid goconst warnings.
const (
	formatJSON    = "json"
	formatTOML    = "toml"
	formatYAML    = "yaml"
	extYAML       = ".yaml"
	extYML        = ".yml"
	boolStrTrue   = "true"
	sourceDefault = "default"
)

// configCmd represents the config command group.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
	Long: `Commands for managing markata-go configuration.

Subcommands:
	  show     - Display the resolved configuration
	  get      - Get a specific configuration value
	  set      - Set a configuration value
	  validate - Validate the configuration file
	  init     - Create a new configuration file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return unknownCommandError(cmd, args[0])
		}
		return nil
	},
	RunE: runConfigCommand,
}

// configShowCmd shows the resolved configuration.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display resolved configuration",
	Long: `Display the fully resolved configuration with all defaults applied.

The configuration is merged from:
  1. Default values
  2. Config file (discovered or specified)
  3. Environment variables

Example usage:
  markata-go config show              # Show as YAML
  markata-go config show --json       # Show as JSON
  markata-go config show --toml       # Show as TOML
  markata-go config show              # Show annotated effective YAML
  markata-go config show --no-annotate # Show YAML without source annotations
  markata-go config show --diff       # Show only user-provided values`,
	RunE: runConfigShowCommand,
}

// configGetCmd gets a specific configuration value.
var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get specific config value",
	Long: `Get a specific configuration value by key.

Supports dot notation for nested values.

Example usage:
  markata-go config get output_dir
  markata-go config get glob.patterns
  markata-go config get feed_defaults.items_per_page`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGetCommand,
}

// configValidateCmd validates the configuration.
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration",
	Long: `Validate the configuration file and report any errors or warnings.

Exit codes:
  0 - Configuration is valid (warnings may be present)
  1 - Configuration has errors

Example usage:
  markata-go config validate
  markata-go config validate -c custom-config.toml`,
	RunE: runConfigValidateCommand,
}

// configInitCmd creates a new configuration file.
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create new config file",
	Long: `Create a new configuration file with sensible defaults.

Supports TOML, YAML, and JSON formats (detected from filename).

Example usage:
  markata-go config init                    # Creates markata-go.toml
  markata-go config init markata-go.yaml    # Creates YAML config
  markata-go config init --force            # Overwrite existing file`,
	RunE: runConfigInitCommand,
}

// configSetCmd sets a configuration value.
var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set configuration value",
	Long: `Set a configuration value in the config file.

Supports dot notation for nested values. Values are automatically
type-detected (string, int, bool, array/object via JSON).

Example usage:
  markata-go config set output_dir "dist"
  markata-go config set url "https://example.com"
  markata-go config set concurrency 4
  markata-go config set glob.use_gitignore true
  markata-go config set glob.patterns '["posts/**/*.md", "pages/*.md"]'
  markata-go config set feed_defaults.items_per_page 15

Flags:
  --dry-run   Show what would be changed without writing
  --backup    Create a backup before modifying the config file`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSetCommand,
}

var (
	// configFormat specifies the output format for config show.
	configFormat string

	// configShowAnnotate shows source of each config value.
	configShowAnnotate bool

	// configShowNoAnnotate disables source comments in YAML output.
	configShowNoAnnotate bool

	// configShowSources explicitly enables source comments in YAML output.
	configShowSources bool

	// configShowDiff shows only user-provided values that differ from defaults.
	configShowDiff bool

	// configInitForce overwrites existing config file.
	configInitForce bool

	// configSetDryRun shows what would be changed without writing.
	configSetDryRun bool

	// configSetBackup creates a backup before modifying.
	configSetBackup bool
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configInitCmd)

	configShowCmd.Flags().StringVar(&configFormat, "format", "yaml", "output format (yaml, json, toml)")
	configShowCmd.Flags().Bool("json", false, "output as JSON (shorthand for --format=json)")
	configShowCmd.Flags().Bool("toml", false, "output as TOML (shorthand for --format=toml)")
	configShowCmd.Flags().BoolVar(&configShowAnnotate, "annotate", true, "show source comments in YAML output")
	configShowCmd.Flags().BoolVar(&configShowNoAnnotate, "no-annotate", false, "show YAML without source comments")
	configShowCmd.Flags().BoolVar(&configShowSources, "sources", false, "show source comments in YAML output")
	configShowCmd.Flags().BoolVar(&configShowDiff, "diff", false, "show only user-provided values that differ from defaults")

	configInitCmd.Flags().BoolVar(&configInitForce, "force", false, "overwrite existing file")

	configSetCmd.Flags().BoolVar(&configSetDryRun, "dry-run", false, "show what would be changed without writing")
	configSetCmd.Flags().BoolVar(&configSetBackup, "backup", false, "create backup before modifying")
}

func runConfigCommand(cmd *cobra.Command, _ []string) error {
	currentCmd = cmd
	return runConfigShowCommand(configShowCmd, nil)
}

func runConfigShowCommand(cmd *cobra.Command, _ []string) error {
	format, err := resolveConfigShowFormat(cmd)
	if err != nil {
		return err
	}

	cfg, _, configPaths, err := loadManagerConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	cliSources := make(map[string]string)
	if outputDir != "" {
		cfg.OutputDir = outputDir
		cliSources["output_dir"] = "CLI --output — pass a different flag value"
	}

	// Handle --annotate and --diff modes
	if ((configShowAnnotate || configShowSources) && !configShowNoAnnotate && format == formatYAML) || configShowDiff {
		return runConfigShowWithSources(cfg, configPaths, cliSources)
	}

	displayMap, err := configToMap(cfg)
	if err != nil {
		return fmt.Errorf("failed to convert config to map: %w", err)
	}

	// Output in requested format
	switch format {
	case formatJSON:
		data, err := json.MarshalIndent(displayMap, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		outln(string(data))

	case formatTOML:
		if err := toml.NewEncoder(outWriter()).Encode(displayMap); err != nil {
			return fmt.Errorf("failed to marshal TOML: %w", err)
		}

	case formatYAML:
		data, err := yaml.Marshal(displayMap)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML: %w", err)
		}
		outText(string(data))

	default:
		return newUsageError(fmt.Errorf("invalid format %q (expected yaml, json, or toml)", format))
	}

	return nil
}

func resolveConfigShowFormat(cmd *cobra.Command) (string, error) {
	format := strings.ToLower(strings.TrimSpace(configFormat))
	if format == "" {
		format = formatYAML
	}

	jsonFlag, err := cmd.Flags().GetBool(formatJSON)
	if err != nil {
		return "", fmt.Errorf("failed to get json flag: %w", err)
	}
	tomlFlag, err := cmd.Flags().GetBool(formatTOML)
	if err != nil {
		return "", fmt.Errorf("failed to get toml flag: %w", err)
	}

	if cmd.Flags().Changed(formatJSON) && cmd.Flags().Changed(formatTOML) && jsonFlag && tomlFlag {
		return "", newUsageError(fmt.Errorf("cannot use --json and --toml together"))
	}

	if cmd.Flags().Changed("format") {
		if jsonFlag && format != formatJSON {
			return "", newUsageError(fmt.Errorf("cannot use --format=%s with --json", format))
		}
		if tomlFlag && format != formatTOML {
			return "", newUsageError(fmt.Errorf("cannot use --format=%s with --toml", format))
		}
	}

	if jsonFlag {
		format = formatJSON
	}
	if tomlFlag {
		format = formatTOML
	}

	switch format {
	case formatJSON, formatTOML, formatYAML:
		return format, nil
	default:
		return "", newUsageError(fmt.Errorf("invalid format %q (expected yaml, json, or toml)", format))
	}
}

// runConfigShowWithSources handles --annotate and --diff flags.
// It compares user config with defaults to show value sources.
func runConfigShowWithSources(merged *models.Config, configPaths []string, cliSources map[string]string) error {
	mergedMap, err := configToMap(merged)
	if err != nil {
		return fmt.Errorf("failed to convert config to map: %w", err)
	}

	sourcePaths, err := discoverConfigSourcePaths(configPaths)
	if err != nil {
		return fmt.Errorf("discover config sources: %w", err)
	}

	userMap, err := loadUserConfigMap(sourcePaths)
	if err != nil {
		return err
	}
	redactBuilderAdminWebhookSecret(userMap)

	if configShowDiff {
		// Show only values that differ from defaults
		return showDiffConfig(userMap, sourcePaths)
	}
	fileSources, err := collectConfigFileSources(sourcePaths)
	if err != nil {
		return err
	}

	// Show annotated config
	return showAnnotatedConfig(mergedMap, userMap, sourcePaths, fileSources, collectEnvironmentSources(), cliSources)
}

// redactBuilderAdminWebhookSecret removes the credential from raw user config
// before config show --diff or --annotate writes it to the terminal.
func redactBuilderAdminWebhookSecret(config map[string]interface{}) {
	markata, ok := config["markata-go"].(map[string]interface{})
	if !ok {
		return
	}
	builderAdmin, ok := markata["builder_admin"].(map[string]interface{})
	if !ok {
		return
	}
	webhook, ok := builderAdmin["webhook"].(map[string]interface{})
	if !ok {
		return
	}
	delete(webhook, "secret")
}

// configToMap converts a Config struct to a map[string]interface{}.
func configToMap(cfg *models.Config) (map[string]interface{}, error) {
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// showDiffConfig prints only user-provided values that differ from defaults.
func showDiffConfig(userMap map[string]interface{}, configPaths []string) error {
	if len(userMap) == 0 {
		outln("# No user configuration found")
		outln("# All values are defaults")
		return nil
	}

	if len(configPaths) > 0 {
		outln("# User configuration from:")
		for _, path := range configPaths {
			outlnf("#   - %s", path)
		}
	}
	outln("# Values below differ from defaults:")
	outln()

	// Output the user config as YAML with comments
	data, err := yaml.Marshal(userMap)
	if err != nil {
		return fmt.Errorf("failed to marshal diff: %w", err)
	}
	outText(string(data))

	return nil
}

// showAnnotatedConfig prints the merged config with source annotations.
//
//nolint:gocyclo // YAML indentation and source annotations require line-aware handling.
func showAnnotatedConfig(merged, user map[string]interface{}, configPaths []string, fileSources map[string]configFileSource, envSources, cliSources map[string]string) error {
	userSource := userSourceLabel(configPaths)
	var output strings.Builder

	// Print header
	fmt.Fprintln(&output, "# Configuration with source annotations")
	if len(configPaths) > 0 {
		fmt.Fprintln(&output, "# User config files:")
		for _, path := range configPaths {
			fmt.Fprintf(&output, "#   - %s\n", path)
		}
	} else {
		fmt.Fprintln(&output, "# User config: (none found, using defaults)")
	}
	fmt.Fprintln(&output)

	// Get YAML representation of merged config
	data, err := yaml.Marshal(merged)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	type yamlPathFrame struct {
		indent int
		key    string
	}

	var stack []yamlPathFrame

	// Split into lines and annotate each
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line == "" {
			fmt.Fprintln(&output)
			continue
		}

		// Parse the line to extract key path
		trimmed := strings.TrimLeft(line, " ")
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "#") {
			// Array item or comment - just print
			fmt.Fprintln(&output, line)
			continue
		}

		// Try to extract key from "key: value" format
		colonIdx := strings.Index(trimmed, ":")
		if colonIdx == -1 {
			fmt.Fprintln(&output, line)
			continue
		}

		key := trimmed[:colonIdx]
		indent := len(line) - len(trimmed)

		for len(stack) > 0 && indent <= stack[len(stack)-1].indent {
			stack = stack[:len(stack)-1]
		}

		path := make([]string, 0, len(stack)+1)
		for _, frame := range stack {
			path = append(path, frame.key)
		}
		path = append(path, key)

		// Determine source based on the effective file assignment. Preserve the
		// coarse user label for container values that have no leaf assignment.
		pathKey := strings.Join(path, ".")
		source := sourceComment(path, fileSources)
		if envSource, ok := envSources[pathKey]; ok {
			source = envSource
		}
		if cliSource, ok := cliSources[pathKey]; ok {
			source = cliSource
		}
		if strings.HasPrefix(source, "default —") && hasPath(user, path) {
			source = userSource
		}

		// Calculate padding for alignment
		padding := 40 - len(line)
		if padding < 2 {
			padding = 2
		}

		// Print with annotation
		// Only annotate leaf values (lines with values after the colon)
		valueAfterColon := strings.TrimSpace(trimmed[colonIdx+1:])
		if valueAfterColon != "" && !strings.HasPrefix(valueAfterColon, "|") && !strings.HasPrefix(valueAfterColon, ">") {
			fmt.Fprintf(&output, "%s%s# %s\n", line, strings.Repeat(" ", padding), source)
		} else {
			if hasPath(user, path) {
				fmt.Fprintf(&output, "%s%s# %s\n", line, strings.Repeat(" ", padding), source)
			} else {
				fmt.Fprintln(&output, line)
			}
			stack = append(stack, yamlPathFrame{indent: indent, key: key})
		}
	}

	_, chromaStyle := resolveLSPSetupRendering()
	outText(highlightLSPSetupSnippet(output.String(), "yaml", chromaStyle))

	return nil
}

func loadUserConfigMap(configPaths []string) (map[string]interface{}, error) {
	var userMap map[string]interface{}

	for _, configPath := range configPaths {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
		}

		wrapper, err := parseConfigToMap(data, formatFromPath(configPath))
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
		}

		inner, ok := wrapper["markata-go"].(map[string]interface{})
		if !ok {
			continue
		}

		userMap = mergeConfigMaps(userMap, inner)
	}

	return userMap, nil
}

func mergeConfigMaps(base, override map[string]interface{}) map[string]interface{} {
	if base == nil {
		base = make(map[string]interface{})
	}

	for key, value := range override {
		overrideMap, ok := value.(map[string]interface{})
		if !ok {
			base[key] = value
			continue
		}

		baseMap, ok := base[key].(map[string]interface{})
		if !ok {
			baseMap = make(map[string]interface{})
		}
		base[key] = mergeConfigMaps(baseMap, overrideMap)
	}

	return base
}

func userSourceLabel(configPaths []string) string {
	if len(configPaths) == 0 {
		return "user"
	}
	if len(configPaths) == 1 {
		return "user: " + filepath.Base(configPaths[0])
	}
	return fmt.Sprintf("user: merged config (%d files)", len(configPaths))
}

func hasPath(m map[string]interface{}, path []string) bool {
	_, ok := getPathValue(m, path)
	return ok
}

func getPathValue(m map[string]interface{}, path []string) (interface{}, bool) {
	if m == nil || len(path) == 0 {
		return nil, false
	}

	var current interface{} = m
	for _, key := range path {
		currentMap, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = currentMap[key]
		if !ok {
			return nil, false
		}
	}

	return current, true
}

func runConfigGetCommand(_ *cobra.Command, args []string) error {
	key := args[0]

	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.Discover()
		if err != nil {
			return fmt.Errorf("no config file found (use -c to specify one or run 'config init' first): %w", err)
		}
	}

	value, err := config.GetValueFromFile(configPath, key)
	if err != nil {
		if errors.Is(err, config.ErrKeyNotFound) {
			return newUsageError(configKeyError(key))
		}
		return fmt.Errorf("cannot read %q from %s: %w", key, configPath, err)
	}

	// Print value
	switch v := value.(type) {
	case string:
		outln(v)
	case []string:
		for _, s := range v {
			outln(s)
		}
	default:
		data, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			outlnf("%v", v)
		} else {
			outln(string(data))
		}
	}

	return nil
}

func runConfigValidateCommand(_ *cobra.Command, _ []string) error {
	// Load and validate configuration
	cfg, _, _, err := loadManagerConfig(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	validationErrs := config.ValidateConfig(cfg)
	actualErrors, warnings := config.SplitErrorsAndWarnings(validationErrs)

	configPath := cfgFile
	if configPath == "" {
		if len(mergeConfigFiles) > 0 {
			configPath = fmt.Sprintf("(merged config: %s)", strings.Join(mergeConfigFiles, ", "))
		} else if discovered, discoverErr := config.Discover(); discoverErr == nil {
			configPath = discovered
		}
	}
	if configPath == "" && len(mergeConfigFiles) > 0 {
		configPath = fmt.Sprintf("(merged config: %s)", strings.Join(mergeConfigFiles, ", "))
	}
	if configPath == "" {
		configPath = "(defaults)"
	}

	// Print warnings
	if len(warnings) > 0 {
		errln(colorize("Warnings:", currentLogTheme.Warning, colorEnabledFor(errorOutputIsTerminal())))
		for _, w := range warnings {
			errlnf("  - %v", w)
		}
		errln()
	}

	// Print errors
	if len(actualErrors) > 0 {
		errln(colorize("Errors:", currentLogTheme.Error, colorEnabledFor(errorOutputIsTerminal())))
		for _, e := range actualErrors {
			errlnf("  - %v", e)
		}
		errln("Next: Correct the listed settings and rerun 'markata-go config validate'.")
		return newExitCodeError(1, nil)
	}

	// Print success
	outlnf("Configuration is valid: %s", configPath)

	if verbose {
		errln("\nConfiguration summary:")
		errlnf("  Output directory: %s", cfg.OutputDir)
		errlnf("  Site URL: %s", cfg.URL)
		errlnf("  Site title: %s", cfg.Title)
		errlnf("  Glob patterns: %v", cfg.GlobConfig.Patterns)
		errlnf("  Feeds defined: %d", len(cfg.Feeds))
	}

	return nil
}

func runConfigInitCommand(_ *cobra.Command, args []string) error {
	// Determine output filename
	filename := defaultConfigFilename
	if len(args) > 0 {
		filename = args[0]
	}

	// Check if file exists
	if _, err := os.Stat(filename); err == nil && !configInitForce {
		return fmt.Errorf("file already exists: %s (use --force to overwrite)", filename)
	}

	// Determine format from filename
	var content string
	switch {
	case strings.HasSuffix(filename, extYAML) || strings.HasSuffix(filename, extYML):
		content = defaultConfigYAML
	case strings.HasSuffix(filename, ".json"):
		content = defaultConfigJSON
	default:
		content = defaultConfigTOML
	}

	// Write file (0o644 is appropriate for config files that should be world-readable)
	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil { //nolint:gosec // config files should be readable
		return fmt.Errorf("failed to write config file: %w", err)
	}

	outlnf("Created: %s", filename)
	return nil
}

func runConfigSetCommand(_ *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Discover config file
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.Discover()
		if err != nil {
			return fmt.Errorf("no config file found (use -c to specify one or run 'config init' first): %w", err)
		}
	}

	// Parse the value with automatic type detection
	parsedValue := parseValue(value)

	oldValue, err := config.GetValueFromFile(configPath, key)
	if err != nil {
		oldValue = nil
		if errors.Is(err, config.ErrKeyNotFound) {
			candidate := configKeyError(key)
			var input *inputError
			if errors.As(candidate, &input) && len(input.matches) > 0 && safeConfigSetSuggestion(key, input.matches[0]) {
				return newUsageError(candidate)
			}
		}
	}

	// If dry-run, just show what would change
	if configSetDryRun {
		outlnf("Would update %s in %s:", key, configPath)
		outlnf("  Old: %v", formatValueForDisplay(oldValue))
		outlnf("  New: %v", formatValueForDisplay(parsedValue))
		return nil
	}

	// Create backup if requested
	if configSetBackup {
		backupPath := configPath + ".backup." + time.Now().Format("20060102-150405")
		if err := copyFile(configPath, backupPath); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		errlnf("Created backup: %s", backupPath)
	}

	if err := config.SetValueInFile(configPath, key, parsedValue); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	outlnf("Updated %s in %s", key, configPath)
	return nil
}

func safeConfigSetSuggestion(key, candidate string) bool {
	key = strings.TrimPrefix(key, "markata-go.")
	candidate = strings.TrimPrefix(candidate, "markata-go.")
	keyParent, _, keyNested := strings.Cut(key, ".")
	candidateParent, _, candidateNested := strings.Cut(candidate, ".")
	if keyNested || candidateNested {
		return keyNested && candidateNested && keyParent == candidateParent
	}
	return strings.Contains(key, "_") && strings.Contains(candidate, "_")
}

func configSettingsForSuggestions() []string {
	fields := config.Settings(config.DefaultConfig())
	keys := make([]string, 0, len(fields))
	for i := range fields {
		keys = append(keys, fields[i].Key)
	}
	return keys
}

// formatFromPath determines the config format from a file path.
func formatFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".toml":
		return formatTOML
	case extYAML, extYML:
		return formatYAML
	case ".json":
		return formatJSON
	default:
		return formatTOML
	}
}

// parseConfigToMap parses config data into a generic map[string]interface{}.
// The config is wrapped in a "markata-go" key.
func parseConfigToMap(data []byte, format string) (map[string]interface{}, error) {
	var wrapper map[string]interface{}

	switch format {
	case formatTOML:
		if err := toml.Unmarshal(data, &wrapper); err != nil {
			return nil, err
		}
	case formatYAML:
		if err := yaml.Unmarshal(data, &wrapper); err != nil {
			return nil, err
		}
	case formatJSON:
		if err := json.Unmarshal(data, &wrapper); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	// If the config doesn't have the markata-go wrapper, add it
	if _, ok := wrapper["markata-go"]; !ok {
		wrapper = map[string]interface{}{"markata-go": wrapper}
	}

	return wrapper, nil
}

// parseValue parses a string value with automatic type detection.
func parseValue(value string) interface{} {
	// Check for boolean
	lower := strings.ToLower(value)
	if lower == boolStrTrue {
		return true
	}
	if lower == "false" {
		return false
	}

	// Check for integer
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	// Check for float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		// Only treat as float if it has a decimal point
		if strings.Contains(value, ".") {
			return floatVal
		}
	}

	// Check for JSON array or object
	if strings.HasPrefix(value, "[") || strings.HasPrefix(value, "{") {
		var jsonVal interface{}
		if err := json.Unmarshal([]byte(value), &jsonVal); err == nil {
			return jsonVal
		}
	}

	// Default to string
	return value
}

// formatValueForDisplay formats a value for human-readable display.
func formatValueForDisplay(v interface{}) string {
	if v == nil {
		return "<not set>"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case []interface{}:
		data, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(data)
	case map[string]interface{}:
		data, err := json.Marshal(val)
		if err != nil {
			return fmt.Sprintf("%v", val)
		}
		return string(data)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// Default configuration templates
const defaultConfigTOML = `# Markata-go configuration file

# Site metadata
title = "My Site"
url = "https://example.com"
description = "A site built with markata-go"
author = "Your Name"

# Output settings
output_dir = "output"
templates_dir = "templates"
assets_dir = "static"

# File discovery
[glob]
patterns = ["pages/**/*.md", "posts/**/*.md"]
use_gitignore = true

# Post output formats
[post_formats]
html = true
markdown = true
text = true
og = true

# Feed defaults
[feed_defaults]
items_per_page = 10
orphan_threshold = 3

[feed_defaults.formats]
html = true
simple_html = true
rss = true
atom = true
json = true
sitemap = true

# Default feed
[[feeds]]
slug = "archive"
title = "Archive"
description = "All posts"
filter = "published == true"
sort = "date"
reverse = true
`

const defaultConfigYAML = `# Markata-go configuration file

# Site metadata
title: My Site
url: https://example.com
description: A site built with markata-go
author: Your Name

# Output settings
output_dir: output
templates_dir: templates
assets_dir: static

# File discovery
glob:
  patterns:
    - "pages/**/*.md"
    - "posts/**/*.md"
  use_gitignore: true

# Post output formats
post_formats:
  html: true
  markdown: true
  text: true
  og: true

# Feed defaults
feed_defaults:
  items_per_page: 10
  orphan_threshold: 3
  formats:
    html: true
    simple_html: true
    rss: true
    atom: true
    json: true
    sitemap: true

# Default feed
feeds:
  - slug: archive
    title: Archive
    description: All posts
    filter: "published == true"
    sort: date
    reverse: true
`

const defaultConfigJSON = `{
  "title": "My Site",
  "url": "https://example.com",
  "description": "A site built with markata-go",
  "author": "Your Name",
  "output_dir": "output",
  "templates_dir": "templates",
  "assets_dir": "static",
  "glob": {
    "patterns": ["pages/**/*.md", "posts/**/*.md"],
    "use_gitignore": true
  },
  "post_formats": {
    "html": true,
    "markdown": true,
    "text": true,
    "og": true
  },
  "feed_defaults": {
    "items_per_page": 10,
    "orphan_threshold": 3,
    "formats": {
      "html": true,
      "simple_html": true,
      "rss": true,
      "atom": true,
      "json": true,
      "sitemap": true
    }
  },
  "feeds": [
    {
      "slug": "archive",
      "title": "Archive",
      "description": "All posts",
      "filter": "published == true",
      "sort": "date",
      "reverse": true
    }
  ]
}
`
