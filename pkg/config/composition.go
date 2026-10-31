package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"
)

// loadResolvedConfig loads a config file, resolves recursive includes, then
// materializes the merged raw config into a typed models.Config.
func loadResolvedConfig(configPath string) (*models.Config, error) {
	rawWrapper, err := loadResolvedRawConfig(configPath)
	if err != nil {
		return nil, err
	}

	warnings := normalizeRenderingTheme(rawWrapper)
	defaultRaw, err := rawWrapperFromDefaultConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to encode default config: %w", err)
	}

	rawWrapper = mergeRawMaps(nil, defaultRaw, rawWrapper)
	normalizeRenderingTheme(rawWrapper)

	config, err := configFromRawWrapper(rawWrapper)
	if err != nil {
		return nil, err
	}

	if len(warnings) > 0 {
		config.Extra["theme_migration_warnings"] = warnings
	}
	return config, nil
}

// normalizeRenderingTheme returns conflicts in the supplied values. Collect
// warnings before merging defaults; later passes include generated values.
func normalizeRenderingTheme(wrapper map[string]any) []string {
	markata, ok := wrapper["markata-go"].(map[string]any)
	if !ok {
		return nil
	}
	theme, warnings := renderingcontract.NormalizeTheme(markata)
	markata["theme"] = theme
	wrapper["markata-go"] = markata
	return warnings
}

// loadResolvedRawConfig loads a config file, recursively resolves include
// entries relative to the declaring file, and merges the results.
func loadResolvedRawConfig(configPath string) (map[string]any, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path %s: %w", configPath, err)
	}

	loader := &rawConfigLoader{loaded: make(map[string]bool)}
	return loader.load(absPath, nil)
}

// DiscoverIncludedConfigPaths returns the root config file plus any recursively
// included config files in the order they are resolved.
func DiscoverIncludedConfigPaths(configPath string) ([]string, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path %s: %w", configPath, err)
	}

	collector := &rawConfigPathCollector{seen: make(map[string]bool)}
	return collector.collect(absPath, nil)
}

type rawConfigLoader struct {
	loaded map[string]bool
}

type rawConfigPathCollector struct {
	seen map[string]bool
}

func (l *rawConfigLoader) load(configPath string, stack []string) (map[string]any, error) {
	configPath = filepath.Clean(configPath)

	if containsPath(stack, configPath) {
		cycle := append(append([]string{}, stack...), configPath)
		return nil, fmt.Errorf("config include cycle detected: %s", strings.Join(cycle, " -> "))
	}

	if l.loaded[configPath] {
		return map[string]any{}, nil
	}

	rawWrapper, err := loadRawConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	includes, err := extractIncludePatterns(rawWrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to parse includes in %s: %w", configPath, err)
	}

	merged := cloneMap(rawWrapper)
	childStack := append(append([]string{}, stack...), configPath)
	baseDir := filepath.Dir(configPath)

	for _, pattern := range includes {
		matches, err := resolveIncludePattern(baseDir, pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve include %q in %s: %w", pattern, configPath, err)
		}

		for _, match := range matches {
			child, err := l.load(match, childStack)
			if err != nil {
				return nil, err
			}
			merged = mergeRawMaps(nil, merged, child)
		}
	}

	l.loaded[configPath] = true
	return merged, nil
}

func (c *rawConfigPathCollector) collect(configPath string, stack []string) ([]string, error) {
	configPath = filepath.Clean(configPath)

	if containsPath(stack, configPath) {
		cycle := append(append([]string{}, stack...), configPath)
		return nil, fmt.Errorf("config include cycle detected: %s", strings.Join(cycle, " -> "))
	}

	if c.seen[configPath] {
		return nil, nil
	}

	rawWrapper, err := loadRawConfigFile(configPath)
	if err != nil {
		return nil, err
	}

	c.seen[configPath] = true
	paths := []string{configPath}
	childStack := append(append([]string{}, stack...), configPath)
	baseDir := filepath.Dir(configPath)

	includes, err := extractIncludePatterns(rawWrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to parse includes in %s: %w", configPath, err)
	}

	for _, pattern := range includes {
		matches, err := resolveIncludePattern(baseDir, pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve include %q in %s: %w", pattern, configPath, err)
		}

		for _, match := range matches {
			childPaths, err := c.collect(match, childStack)
			if err != nil {
				return nil, err
			}
			paths = append(paths, childPaths...)
		}
	}

	return paths, nil
}

func loadRawConfigFile(configPath string) (map[string]any, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	rawWrapper, err := loadRawConfigData(data, formatFromPath(configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	return rawWrapper, nil
}

func loadRawConfigData(data []byte, format Format) (map[string]any, error) {
	var rawWrapper map[string]any

	switch format {
	case FormatTOML:
		if err := toml.Unmarshal(data, &rawWrapper); err != nil {
			return nil, err
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &rawWrapper); err != nil {
			return nil, err
		}
	case FormatJSON:
		if err := json.Unmarshal(data, &rawWrapper); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported config format: %s", format)
	}

	normalized, ok := normalizeValue(rawWrapper).(map[string]any)
	if !ok {
		return map[string]any{}, nil
	}
	if err := likelyConfigTypo(normalized); err != nil {
		return nil, err
	}

	return normalized, nil
}

func extractIncludePatterns(rawWrapper map[string]any) ([]string, error) {
	markataGo, ok := rawWrapper["markata-go"].(map[string]any)
	if !ok {
		return nil, nil
	}

	includeValue, ok := markataGo["include"]
	if !ok {
		return nil, nil
	}

	return stringSliceFromValue(includeValue)
}

func resolveIncludePattern(baseDir, pattern string) ([]string, error) {
	resolvedPattern := pattern
	if !filepath.IsAbs(resolvedPattern) {
		resolvedPattern = filepath.Join(baseDir, resolvedPattern)
	}

	resolvedPattern = filepath.Clean(resolvedPattern)
	hasGlob := strings.ContainsAny(pattern, "*?[")

	if !hasGlob {
		info, err := os.Stat(resolvedPattern)
		if err != nil {
			return nil, fmt.Errorf("include target not found: %s", resolvedPattern)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("include target must be a file or glob, got directory: %s", resolvedPattern)
		}
		return []string{resolvedPattern}, nil
	}

	matches, err := filepath.Glob(resolvedPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid include glob %q: %w", pattern, err)
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("include pattern matched no files: %s", resolvedPattern)
	}

	files := make([]string, 0, len(matches))
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			return nil, fmt.Errorf("failed to stat include target %s: %w", match, err)
		}
		if info.IsDir() {
			continue
		}
		files = append(files, filepath.Clean(match))
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("include pattern matched no files: %s", resolvedPattern)
	}

	sort.Strings(files)
	return files, nil
}

func configFromRawWrapper(rawWrapper map[string]any) (*models.Config, error) {
	data, err := json.Marshal(rawWrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resolved config: %w", err)
	}

	var wrapper struct {
		MarkataGo jsonConfig `json:"markata-go"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to decode resolved config: %w", err)
	}
	if err := wrapper.MarkataGo.Head.validate(); err != nil {
		return nil, fmt.Errorf("failed to decode resolved config: %w", err)
	}

	config := wrapper.MarkataGo.toConfig()
	populateExtra(config, rawWrapper)
	return config, nil
}

func rawWrapperFromConfig(config *models.Config) (map[string]any, error) {
	data, err := json.Marshal(map[string]any{"markata-go": config})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}

	var rawWrapper map[string]any
	if err := json.Unmarshal(data, &rawWrapper); err != nil {
		return nil, fmt.Errorf("failed to decode default config: %w", err)
	}

	return rawWrapper, nil
}

// rawWrapperFromDefaultConfig keeps unset formats on built-in feeds absent
// from the raw layer. JSON serialization otherwise turns their zero values
// into explicit false values before user config is merged.
func rawWrapperFromDefaultConfig() (map[string]any, error) {
	defaults := DefaultConfig()
	raw, err := rawWrapperFromConfig(defaults)
	if err != nil {
		return nil, err
	}
	root, ok := raw["markata-go"].(map[string]any)
	if !ok {
		return raw, nil
	}
	feeds, ok := root["feeds"].([]any)
	if !ok {
		return raw, nil
	}
	for i := range defaults.Feeds {
		if i >= len(feeds) || defaults.Feeds[i].Formats.HasAnyEnabled() {
			continue
		}
		if rawFeed, ok := feeds[i].(map[string]any); ok {
			delete(rawFeed, "formats")
		}
	}
	return raw, nil
}

func populateExtra(config *models.Config, rawWrapper map[string]any) {
	markataGoRaw, ok := rawWrapper["markata-go"].(map[string]any)
	if !ok {
		return
	}

	knownKeys := map[string]bool{
		"output_dir": true, "url": true, "title": true, "description": true,
		"author": true, "language": true, "author_url": true, "managing_editor": true,
		"webmaster": true, "copyright": true, "license": true, "assets_dir": true,
		"templates_dir": true, "templates": true, "nav": true, "footer": true,
		"hooks": true, "disabled_hooks": true, "glob": true, "markdown": true,
		"feeds": true, "feed_defaults": true, "concurrency": true, "theme": true,
		"post_formats": true, "well_known": true, "seo": true, "indieauth": true,
		"webmention": true, "components": true, "layout": true, "sidebar": true,
		"toc": true, "header": true, "blogroll": true, "mentions": true,
		"template_presets": true, "default_templates": true, "head": true,
		"content_templates": true, "footer_layout": true, "search": true,
		"plugins": true, "thoughts": true, "wikilinks": true, "tags": true,
		"tag_aggregator": true, "websub": true, "shortcuts": true, "view_transitions": true,
		"encryption": true, "authors": true, "garden": true, "feeds_page": true,
		"assets": true, "resource_hints": true, "error_pages": true, "theme_calendar": true,
		"builder_admin": true, "include": true, "fontpack": true, "fontpacks_file": true,
		"images": true,
	}

	if config.Extra == nil {
		config.Extra = make(map[string]any)
	}

	for key, value := range markataGoRaw {
		if !knownKeys[key] {
			config.Extra[key] = value
		}
	}

	populateBuilderAdminExtra(config, markataGoRaw)
}

func mergeRawMaps(path []string, base, override map[string]any) map[string]any {
	if base == nil {
		return cloneMap(override)
	}
	if override == nil {
		return cloneMap(base)
	}

	result := cloneMap(base)
	for key, overrideValue := range override {
		currentPath := append(append([]string{}, path...), key)
		if baseValue, ok := result[key]; ok {
			result[key] = mergeRawValue(currentPath, baseValue, overrideValue)
			continue
		}
		result[key] = cloneValue(overrideValue)
	}

	return result
}

func mergeRawValue(path []string, base, override any) any {
	baseMap, baseIsMap := base.(map[string]any)
	overrideMap, overrideIsMap := override.(map[string]any)
	if baseIsMap && overrideIsMap {
		return mergeRawMaps(path, baseMap, overrideMap)
	}

	baseSlice, baseIsSlice := base.([]any)
	overrideSlice, overrideIsSlice := override.([]any)
	if baseIsSlice && overrideIsSlice {
		if isFeedsPath(path) {
			return mergeFeedSlices(baseSlice, overrideSlice)
		}
		return cloneSlice(overrideSlice)
	}

	return cloneValue(override)
}

func mergeFeedSlices(base, override []any) []any {
	result := cloneSlice(base)
	indexes := make(map[string]int)

	for i, value := range result {
		if feedMap, ok := value.(map[string]any); ok {
			if slug, ok := feedMap["slug"].(string); ok && slug != "" {
				indexes[slug] = i
			}
		}
	}

	for _, value := range override {
		feedMap, ok := value.(map[string]any)
		if !ok {
			result = append(result, cloneValue(value))
			continue
		}

		slugValue := feedMap["slug"]
		slug, ok := slugValue.(string)
		if !ok || slug == "" {
			result = append(result, cloneValue(feedMap))
			continue
		}

		if index, ok := indexes[slug]; ok {
			baseMap, ok := result[index].(map[string]any)
			if !ok {
				result[index] = cloneValue(feedMap)
				continue
			}
			result[index] = mergeRawMaps([]string{"markata-go", "feeds", slug}, baseMap, feedMap)
			continue
		}

		indexes[slug] = len(result)
		result = append(result, cloneValue(feedMap))
	}

	return result
}

func isFeedsPath(path []string) bool {
	return len(path) == 2 && path[0] == "markata-go" && path[1] == "feeds"
}

func stringSliceFromValue(value any) ([]string, error) {
	switch typed := value.(type) {
	case string:
		return []string{typed}, nil
	case []string:
		return append([]string{}, typed...), nil
	case []any:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("include entries must be strings")
			}
			result = append(result, text)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("include must be a string or list of strings")
	}
}

func normalizeValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		normalized := make(map[string]any, len(typed))
		for key, item := range typed {
			normalized[key] = normalizeValue(item)
		}
		return normalized
	case map[any]any:
		normalized := make(map[string]any, len(typed))
		for key, item := range typed {
			normalized[fmt.Sprint(key)] = normalizeValue(item)
		}
		return normalized
	case []any:
		normalized := make([]any, 0, len(typed))
		for _, item := range typed {
			normalized = append(normalized, normalizeValue(item))
		}
		return normalized
	default:
		reflected := reflect.ValueOf(value)
		switch reflected.Kind() {
		case reflect.Map:
			normalized := make(map[string]any, reflected.Len())
			for _, key := range reflected.MapKeys() {
				normalized[fmt.Sprint(key.Interface())] = normalizeValue(reflected.MapIndex(key).Interface())
			}
			return normalized
		case reflect.Slice, reflect.Array:
			normalized := make([]any, 0, reflected.Len())
			for i := 0; i < reflected.Len(); i++ {
				normalized = append(normalized, normalizeValue(reflected.Index(i).Interface()))
			}
			return normalized
		default:
			return value
		}
	}
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = cloneValue(value)
	}
	return cloned
}

func cloneSlice(input []any) []any {
	if input == nil {
		return nil
	}
	cloned := make([]any, 0, len(input))
	for _, value := range input {
		cloned = append(cloned, cloneValue(value))
	}
	return cloned
}

func cloneValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneMap(typed)
	case []any:
		return cloneSlice(typed)
	default:
		return typed
	}
}

func containsPath(paths []string, target string) bool {
	for _, path := range paths {
		if path == target {
			return true
		}
	}
	return false
}
