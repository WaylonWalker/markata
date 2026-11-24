package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
)

type configFileSource struct {
	path string
	line int
}

func discoverConfigSourcePaths(configPaths []string) ([]string, error) {
	var paths []string
	seen := make(map[string]bool)
	for _, configPath := range configPaths {
		included, err := config.DiscoverIncludedConfigPaths(configPath)
		if err != nil {
			return nil, err
		}
		for _, path := range included {
			if !seen[path] {
				seen[path] = true
				paths = append(paths, path)
			}
		}
	}
	return paths, nil
}

func collectConfigFileSources(paths []string) (map[string]configFileSource, error) {
	sources := make(map[string]configFileSource)
	for _, path := range paths {
		fileSources, err := scanConfigFileSources(path)
		if err != nil {
			return nil, err
		}
		for key, source := range fileSources {
			sources[key] = source
		}
	}
	return sources, nil
}

func scanConfigFileSources(path string) (map[string]configFileSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config source %s: %w", path, err)
	}
	sources := make(map[string]configFileSource)
	var section []string
	for index, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(strings.SplitN(line, "#", 2)[0])
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			section = strings.Split(strings.Trim(trimmed, "[]"), ".")
			if len(section) > 0 && section[0] == "markata-go" {
				section = section[1:]
			}
			continue
		}
		key, ok := configSourceKey(trimmed)
		if !ok {
			continue
		}
		pathParts := append(append([]string{}, section...), key)
		sources[strings.Join(pathParts, ".")] = configFileSource{path: filepath.Clean(path), line: index + 1}
	}
	return sources, nil
}

func configSourceKey(line string) (string, bool) {
	if index := strings.Index(line, "="); index > 0 {
		return strings.Trim(strings.TrimSpace(line[:index]), "\"'"), true
	}
	if index := strings.Index(line, ":"); index > 0 {
		return strings.Trim(strings.TrimSpace(line[:index]), "\"'"), true
	}
	return "", false
}

func sourceComment(path []string, sources map[string]configFileSource) string {
	key := strings.Join(path, ".")
	if source, ok := sources[key]; ok {
		return fmt.Sprintf("file %s:%d:1 — edit it or run: markata-go config set %s <value>", source.path, source.line, key)
	}
	return fmt.Sprintf("default — run: markata-go config set %s <value>", key)
}

func collectEnvironmentSources() map[string]string {
	sources := make(map[string]string)
	for _, entry := range os.Environ() {
		name, _, _ := strings.Cut(entry, "=")
		if !strings.HasPrefix(name, "MARKATA_GO_") {
			continue
		}
		path := strings.ToLower(strings.ReplaceAll(strings.TrimPrefix(name, "MARKATA_GO_"), "_", "."))
		for old, newValue := range map[string]string{"output.dir": "output_dir", "assets.dir": "assets_dir", "templates.dir": "templates_dir", "author.url": "author_url", "use.gitignore": "use_gitignore"} {
			path = strings.ReplaceAll(path, old, newValue)
		}
		sources[path] = "environment " + name + " — export a new value"
	}
	return sources
}
