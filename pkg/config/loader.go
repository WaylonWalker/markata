package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

// Supported config file names in discovery order
var configFileNames = []string{
	"markata-go.toml",
	"markata-go.yaml",
	"markata-go.yml",
	"markata-go.json",
}

// Format represents a configuration file format.
type Format string

const (
	FormatTOML Format = "toml"
	FormatYAML Format = "yaml"
	FormatJSON Format = "json"
)

// ErrConfigNotFound is returned when no config file is found.
var ErrConfigNotFound = errors.New("no configuration file found")

// Load loads configuration from the specified file path.
// If configPath is empty, it will attempt to discover a config file.
// Environment variable overrides are applied after loading the file.
// A .env file in the current directory is loaded first (if present)
// so that encryption keys and other settings can be stored there.
func Load(configPath string) (*models.Config, error) {
	// Load .env file before anything else so env vars are available
	// for both config parsing and encryption key lookups.
	// Errors loading .env are non-fatal (file may not exist).
	if !runtimeenv.DotEnvDisabled() {
		_ = LoadDotEnv() //nolint:errcheck // .env loading is best-effort
	}

	var config *models.Config
	var err error

	if configPath == "" {
		// Try to discover a config file
		configPath, err = Discover()
		if err != nil {
			if errors.Is(err, ErrConfigNotFound) {
				// No config file found, use defaults with env overrides
				return LoadWithDefaults()
			}
			return nil, err
		}
	}

	config, err = loadResolvedConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Apply environment variable overrides
	if err := ApplyEnvOverrides(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	return config, nil
}

// Discover searches for a configuration file in the standard locations.
// It returns the path to the first config file found, or ErrConfigNotFound.
//
// Discovery order:
//  1. ./markata-go.toml
//  2. ./markata-go.yaml (or .yml)
//  3. ./markata-go.json
//  4. ~/.config/markata-go/config.toml
func Discover() (string, error) {
	// Check current directory first
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	for _, name := range configFileNames {
		path := filepath.Join(cwd, name)
		if fileExists(path) {
			return path, nil
		}
	}

	// Check user config directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		userConfigPath := filepath.Join(homeDir, ".config", "markata-go", "config.toml")
		if fileExists(userConfigPath) {
			return userConfigPath, nil
		}
	}

	return "", ErrConfigNotFound
}

// LoadWithDefaults returns a configuration with default values and
// environment variable overrides applied.
func LoadWithDefaults() (*models.Config, error) {
	// Run defaults through the same theme normalization as a config file so a
	// config-less site (e.g. `markata-go build page.md`) gets the contract
	// default theme, including a palette for the theme picker.
	emptyRaw := map[string]any{"markata-go": map[string]any{}}
	warnings := normalizeRenderingTheme(emptyRaw)
	defaultRaw, err := rawWrapperFromDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to encode default config: %w", err)
	}
	mergedRaw := mergeRawMaps(nil, defaultRaw, emptyRaw)
	normalizeRenderingTheme(mergedRaw)
	config, err := configFromRawWrapper(mergedRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode default config: %w", err)
	}
	if len(warnings) > 0 {
		config.Extra["theme_migration_warnings"] = warnings
	}

	if err := ApplyEnvOverrides(config); err != nil {
		return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
	}

	return config, nil
}

// LoadWithMerge loads a base configuration and merges it with one or more
// override configurations. Override configs are applied in order, with later
// configs taking precedence over earlier ones.
//
// This is useful for environment-specific overrides, e.g.:
//   - markata-go.toml (base config)
//   - markata-go.local.toml (local overrides, gitignored)
//   - fast-markata-go.toml (fast build overrides)
//
// Example usage:
//
//	// Load base + local overrides
//	config, err := LoadWithMerge("markata-go.toml", "markata-go.local.toml")
//
//	// Load with multiple overrides
//	config, err := LoadWithMerge("markata-go.toml", "fast-markata-go.toml", "debug.toml")
func LoadWithMerge(basePath string, overridePaths ...string) (*models.Config, error) {
	return LoadWithMergeOptions(LoadOptions{}, basePath, overridePaths...)
}

// LoadOptions controls optional configuration-loading behavior.
type LoadOptions struct {
	// DisableDotEnv prevents loading a .env file before parsing configuration.
	DisableDotEnv bool
	// DisableEnvOverrides prevents ambient MARKATA_GO_* variables from changing
	// the explicitly supplied configuration.
	DisableEnvOverrides bool
	// Overlay is a raw config wrapper (keyed by "markata-go") merged after
	// every file and before environment overrides. markata-go serve uses it
	// to preview unsaved settings; see SettingsOverlay.
	Overlay map[string]any
	// Remove lists key paths (below the markata-go wrapper) deleted from the
	// merged files and overlay before defaults apply, so each falls back to
	// its default. markata-go serve uses it to preview "reset to default".
	Remove [][]string
}

// LoadWithMergeOptions loads and merges configuration with explicit options.
// Isolated callers can use these options to avoid process-local configuration
// sources while ordinary callers keep the existing defaults.
func LoadWithMergeOptions(options LoadOptions, basePath string, overridePaths ...string) (*models.Config, error) {
	// Load .env file first
	if !options.DisableDotEnv && !runtimeenv.DotEnvDisabled() {
		_ = LoadDotEnv() //nolint:errcheck // .env loading is best-effort
	}

	// Load base config (or use empty config if no base path provided)
	var err error
	var mergedRaw map[string]any
	if basePath != "" {
		mergedRaw, err = loadResolvedRawConfig(basePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load base config %s: %w", basePath, err)
		}
	} else {
		mergedRaw = map[string]any{}
	}

	// Apply each override config in order
	for _, overridePath := range overridePaths {
		if overridePath == "" {
			continue
		}

		overrideRaw, err := loadResolvedRawConfig(overridePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load override config %s: %w", overridePath, err)
		}

		mergedRaw = mergeRawMaps(nil, mergedRaw, overrideRaw)
	}

	if len(options.Overlay) > 0 {
		mergedRaw = mergeRawMaps(nil, mergedRaw, options.Overlay)
	}
	for _, path := range options.Remove {
		removeRawSetting(mergedRaw, path)
	}

	warnings := normalizeRenderingTheme(mergedRaw)
	defaultRaw, err := rawWrapperFromDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to encode default config: %w", err)
	}

	mergedRaw = mergeRawMaps(nil, defaultRaw, mergedRaw)
	normalizeRenderingTheme(mergedRaw)

	baseConfig, err := configFromRawWrapper(mergedRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to decode merged config: %w", err)
	}

	// Apply environment variable overrides last (highest precedence), unless
	// the caller needs a configuration identity independent of ambient state.
	if !options.DisableEnvOverrides {
		if err := ApplyEnvOverrides(baseConfig); err != nil {
			return nil, fmt.Errorf("failed to apply environment overrides: %w", err)
		}
	}

	if len(warnings) > 0 {
		baseConfig.Extra["theme_migration_warnings"] = warnings
	}
	return baseConfig, nil
}

// LoadSingleConfig loads a single config file without merging defaults or env overrides.
// This is useful for loading partial configs that will be merged with others.
func LoadSingleConfig(configPath string) (*models.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	format := formatFromPath(configPath)

	var config *models.Config
	switch format {
	case FormatTOML:
		_, err = ParseTOML(data)
	case FormatYAML:
		_, err = ParseYAML(data)
	case FormatJSON:
		_, err = ParseJSON(data)
	default:
		return nil, fmt.Errorf("unsupported config format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	resolvedRaw, rawErr := loadRawConfigData(data, format)
	if rawErr != nil {
		return nil, rawErr
	}
	warnings := normalizeRenderingTheme(resolvedRaw)
	config, err = configFromRawWrapper(resolvedRaw)
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		config.Extra["theme_migration_warnings"] = warnings
	}
	return config, nil
}

// LoadFromString parses configuration from a string with the specified format.
func LoadFromString(data string, format Format) (*models.Config, error) {
	resolvedRaw, err := loadRawConfigData([]byte(data), format)
	if err != nil {
		return nil, err
	}

	switch format {
	case FormatTOML:
		_, err = ParseTOML([]byte(data))
	case FormatYAML:
		_, err = ParseYAML([]byte(data))
	case FormatJSON:
		_, err = ParseJSON([]byte(data))
	default:
		return nil, fmt.Errorf("unsupported config format: %s", format)
	}

	if err != nil {
		return nil, err
	}

	warnings := normalizeRenderingTheme(resolvedRaw)
	defaultRaw, err := rawWrapperFromDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to encode default config: %w", err)
	}

	mergedRaw := mergeRawMaps(nil, defaultRaw, resolvedRaw)
	normalizeRenderingTheme(mergedRaw)
	config, err := configFromRawWrapper(mergedRaw)
	if err != nil {
		return nil, err
	}

	if len(warnings) > 0 {
		config.Extra["theme_migration_warnings"] = warnings
	}
	return config, nil
}

// MustLoad loads configuration and panics on error.
// This is useful for initialization code where config loading failure
// should be fatal.
func MustLoad(configPath string) *models.Config {
	config, err := Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return config
}

// LoadAndValidate loads and validates configuration.
// Returns the config and any validation errors/warnings.
func LoadAndValidate(configPath string) (*models.Config, []error, error) {
	config, err := Load(configPath)
	if err != nil {
		return nil, nil, err
	}

	validationErrs := ValidateConfig(config)
	return config, validationErrs, nil
}

// LoadAndValidateWithPositions loads configuration and validates it with
// position tracking for enhanced error messages.
// Returns the config, detailed errors with file positions, and any loading error.
func LoadAndValidateWithPositions(configPath string) (*models.Config, *ConfigErrors, error) {
	var actualPath string
	var err error

	if configPath == "" {
		// Try to discover a config file
		actualPath, err = Discover()
		if err != nil {
			if errors.Is(err, ErrConfigNotFound) {
				// No config file found, use defaults with env overrides
				config, loadErr := LoadWithDefaults()
				if loadErr != nil {
					return nil, nil, loadErr
				}
				// Validate without position tracking since there's no file
				configErrors := ValidateConfigWithPositions(config, nil)
				return config, configErrors, nil
			}
			return nil, nil, err
		}
	} else {
		actualPath = configPath
	}

	// Read the root file for position tracking.
	data, err := os.ReadFile(actualPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read config file %s: %w", actualPath, err)
	}

	tracker := NewPositionTracker(data, actualPath)

	config, err := Load(actualPath)
	if err != nil {
		configErrors := &ConfigErrors{}
		configErrors.Add(&ConfigError{
			File:    actualPath,
			Message: fmt.Sprintf("failed to parse configuration: %v", err),
			Field:   "syntax",
		})
		return nil, configErrors, err
	}

	// Validate with position tracking
	configErrors := ValidateConfigWithPositions(config, tracker)

	return config, configErrors, nil
}

// formatFromPath determines the config format from a file path.
func formatFromPath(path string) Format {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".toml":
		return FormatTOML
	case ".yaml", ".yml":
		return FormatYAML
	case ".json":
		return FormatJSON
	default:
		return FormatTOML // Default to TOML
	}
}

// fileExists returns true if the file exists and is not a directory.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// Path holds the result of config file discovery with additional metadata.
type Path struct {
	Path   string // Full path to the config file
	Format Format // Format of the config file
	Source string // Where it was found: "cli", "cwd", "user"
}

// DiscoverAll finds all config files in the standard locations.
// This is useful for debugging or showing available configs.
func DiscoverAll() []Path {
	var found []Path

	cwd, err := os.Getwd()
	if err == nil {
		for _, name := range configFileNames {
			path := filepath.Join(cwd, name)
			if fileExists(path) {
				found = append(found, Path{
					Path:   path,
					Format: formatFromPath(path),
					Source: "cwd",
				})
			}
		}
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		userConfigPath := filepath.Join(homeDir, ".config", "markata-go", "config.toml")
		if fileExists(userConfigPath) {
			found = append(found, Path{
				Path:   userConfigPath,
				Format: FormatTOML,
				Source: "user",
			})
		}
	}

	return found
}

// removeRawSetting deletes markata-go.<path...> from raw, a merged config
// wrapper whose maps the caller owns.
func removeRawSetting(raw map[string]any, path []string) {
	if len(path) == 0 {
		return
	}
	node, ok := raw["markata-go"].(map[string]any)
	if !ok {
		return
	}
	for _, key := range path[:len(path)-1] {
		if node, ok = node[key].(map[string]any); !ok {
			return
		}
	}
	delete(node, path[len(path)-1])
}
