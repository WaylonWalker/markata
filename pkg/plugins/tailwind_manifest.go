package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
)

var tailwindClassAttrPattern = regexp.MustCompile(`(?i)\bclass\s*=\s*(?:"([^"]*)"|'([^']*)')`)

type tailwindBuildPlan struct {
	contentPaths []string
	manifestHash string
	shouldBuild  bool
	cleanup      func()
}

func (p *TailwindPlugin) prepareTailwindBuildPlan(m *lifecycle.Manager) (*tailwindBuildPlan, error) {
	plan := &tailwindBuildPlan{cleanup: func() {}}
	config := m.Config()
	if config == nil {
		return plan, nil
	}

	outputExists := p.tailwindOutputExists(config)
	fastMode := lifecycle.IsServeFastMode(m)
	if !fastMode {
		if extra := config.Extra; extra != nil {
			if fast, ok := extra["fast_mode"].(bool); ok && fast {
				fastMode = true
			}
		}
	}

	if !p.usesGeneratedContentConfig() {
		plan.shouldBuild = !fastMode || !outputExists
		return plan, nil
	}

	manifest := p.buildTailwindManifest(m)

	manifestHash := p.computeTailwindManifestHash(config, manifest)
	plan.manifestHash = manifestHash

	cache := GetBuildCache(m)
	if fastMode && outputExists {
		plan.shouldBuild = false
		return plan, nil
	}
	if outputExists && cache != nil && cache.GetTailwindManifestHash() == manifestHash {
		plan.shouldBuild = false
		return plan, nil
	}

	manifestPath, cleanup, err := writeTailwindManifest(manifest)
	if err != nil {
		return nil, err
	}
	plan.contentPaths = p.generatedTailwindContentPaths(config, manifestPath)
	plan.shouldBuild = true
	plan.cleanup = cleanup
	return plan, nil
}

func (p *TailwindPlugin) buildTailwindManifest(m *lifecycle.Manager) string {
	posts := m.Posts()
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Path < posts[j].Path
	})

	cache := GetBuildCache(m)
	manifestTokens := make(map[string]struct{}, 512)

	for _, post := range posts {
		if post == nil || post.Skip {
			continue
		}

		if strings.TrimSpace(post.HTML) == "" {
			if tokens, ok := getStoredTailwindTokens(cache, post.Path); ok {
				for _, token := range strings.Fields(tokens) {
					manifestTokens[token] = struct{}{}
				}
			}
			continue
		}

		htmlHash := buildcache.ContentHash(post.HTML)
		tokens, ok := getCachedTailwindTokens(cache, post.Path, htmlHash)
		if !ok {
			tokens = extractTailwindTokens(post.HTML)
			if cache != nil {
				cache.CacheTailwindTokens(post.Path, htmlHash, tokens)
			}
		}

		for _, token := range strings.Fields(tokens) {
			manifestTokens[token] = struct{}{}
		}
	}

	tokens := make([]string, 0, len(manifestTokens))
	for token := range manifestTokens {
		tokens = append(tokens, token)
	}
	sort.Strings(tokens)
	return strings.Join(tokens, "\n")
}

func (p *TailwindPlugin) generatedTailwindContentPaths(config *lifecycle.Config, manifestPath string) []string {
	paths := make([]string, 0, 4)
	seen := make(map[string]struct{}, 4)
	add := func(value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		value = filepath.ToSlash(value)
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		paths = append(paths, value)
	}

	add(manifestPath)
	if assetsPattern := tailwindAssetsJavaScriptPattern(config); assetsPattern != "" {
		add(assetsPattern)
	}
	for _, pattern := range tailwindTemplatePatterns(config) {
		add(pattern)
	}

	return paths
}

func tailwindAssetsJavaScriptPattern(config *lifecycle.Config) string {
	assetsDir := ""
	if config != nil && config.Extra != nil {
		if value, ok := config.Extra["assets_dir"].(string); ok {
			assetsDir = strings.TrimSpace(value)
		}
	}
	if assetsDir == "" {
		assetsDir = "static"
	}
	if !pathExists(assetsDir) {
		return ""
	}
	return filepath.Join(assetsDir, "**", "*.js")
}

func tailwindTemplatePatterns(config *lifecycle.Config) []string {
	if config == nil || config.Extra == nil {
		return nil
	}
	value, ok := config.Extra["templates_dir"].(string)
	if !ok || strings.TrimSpace(value) == "" || !pathExists(value) {
		return nil
	}
	return []string{
		filepath.Join(value, "**", "*.html"),
		filepath.Join(value, "**", "*.js"),
		filepath.Join(value, "**", "*.md"),
	}
}

func writeTailwindManifest(content string) (path string, cleanup func(), err error) {
	tmpFile, err := os.CreateTemp("", "markata-tailwind-manifest-*.txt")
	if err != nil {
		return "", nil, fmt.Errorf("tailwind: creating token manifest: %w", err)
	}

	cleanup = func() {
		_ = os.Remove(tmpFile.Name())
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		cleanup()
		_ = tmpFile.Close()
		return "", nil, fmt.Errorf("tailwind: writing token manifest: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("tailwind: closing token manifest: %w", err)
	}

	return tmpFile.Name(), cleanup, nil
}

func getCachedTailwindTokens(cache *buildcache.Cache, path, htmlHash string) (string, bool) {
	if cache == nil {
		return "", false
	}
	return cache.GetCachedTailwindTokens(path, htmlHash)
}

func getStoredTailwindTokens(cache *buildcache.Cache, path string) (string, bool) {
	if cache == nil {
		return "", false
	}
	return cache.GetStoredTailwindTokens(path)
}

func extractTailwindTokens(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	tokens := make(map[string]struct{}, 32)
	for _, matches := range tailwindClassAttrPattern.FindAllStringSubmatch(content, -1) {
		classValue := matches[1]
		if classValue == "" {
			classValue = matches[2]
		}
		for _, token := range strings.Fields(classValue) {
			token = strings.TrimSpace(token)
			if token == "" {
				continue
			}
			tokens[token] = struct{}{}
		}
	}

	if len(tokens) == 0 {
		return ""
	}

	values := make([]string, 0, len(tokens))
	for token := range tokens {
		values = append(values, token)
	}
	sort.Strings(values)
	return strings.Join(values, " ")
}

func hashFileContents(path string) string {
	if strings.TrimSpace(path) == "" {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return buildcache.ContentHash(string(data))
}

func pathExists(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}
