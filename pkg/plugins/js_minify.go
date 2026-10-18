// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/js"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// jsBufPool reuses bytes.Buffer instances across minification calls to reduce allocations.
var jsBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

const jsMinifyPluginName = "js_minify"

// JSMinifyPlugin minifies JavaScript files to reduce file sizes and improve
// Lighthouse performance scores. It runs during the Write stage at
// PriorityLast to ensure all JS files are in their final
// form before minification.
//
// The plugin:
//   - Discovers JS files in the output directory
//   - Minifies each file using tdewolff/minify
//   - Skips excluded files (e.g., already-minified vendor JS)
//   - Reports size reduction statistics
type JSMinifyPlugin struct {
	config   models.JSMinifyConfig
	minifier *minify.M
	exclude  map[string]bool // excluded file patterns for fast lookup
}

// NewJSMinifyPlugin creates a new JSMinifyPlugin.
func NewJSMinifyPlugin() *JSMinifyPlugin {
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)

	return &JSMinifyPlugin{
		minifier: m,
		exclude:  make(map[string]bool),
	}
}

// Name returns the unique name of the plugin.
func (p *JSMinifyPlugin) Name() string {
	return jsMinifyPluginName
}

// Configure reads the JS minify configuration from the manager's config.
func (p *JSMinifyPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Set default config first
	p.config = models.NewJSMinifyConfig()
	p.exclude = make(map[string]bool)

	if config.Extra == nil {
		return nil
	}

	// Try to get js_minify config from Extra
	// It may be a models.JSMinifyConfig or a map[string]interface{} from TOML parsing
	switch v := config.Extra[jsMinifyPluginName].(type) {
	case models.JSMinifyConfig:
		p.config = v
	case map[string]interface{}:
		p.config = p.parseConfigFromMap(v)
	}

	// Build exclude map for fast lookups
	for _, pattern := range p.config.Exclude {
		p.exclude[pattern] = true
	}

	return nil
}

// parseConfigFromMap parses JSMinifyConfig from a raw map (TOML parsing result).
func (p *JSMinifyPlugin) parseConfigFromMap(m map[string]interface{}) models.JSMinifyConfig {
	cfg := models.NewJSMinifyConfig()

	if enabled, ok := m["enabled"].(bool); ok {
		cfg.Enabled = enabled
	}

	// Parse exclude list
	if exclude, ok := m["exclude"].([]interface{}); ok {
		cfg.Exclude = make([]string, 0, len(exclude))
		for _, e := range exclude {
			if s, ok := e.(string); ok {
				cfg.Exclude = append(cfg.Exclude, s)
			}
		}
	}

	return cfg
}

// Write performs JS minification in the output directory.
// Skipped in fast mode (--fast flag) for faster development builds.
func (p *JSMinifyPlugin) Write(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	// Skip minification in fast mode
	if fast, ok := m.Config().Extra["fast_mode"].(bool); ok && fast {
		return nil
	}

	outputDir := m.Config().OutputDir
	jsFiles, err := p.findJSFiles(outputDir)
	if err != nil {
		return fmt.Errorf("finding JS files: %w", err)
	}

	cache, cacheErr := newMinifyCache(m.Config(), p.Name())
	if cacheErr != nil {
		log.Printf("[%s] Warning: initializing asset cache: %v", p.Name(), cacheErr)
	}
	defer cache.close()
	runMinification(p.Name(), jsFiles, p.isExcluded, p.minifyBytes, m.Concurrency(), cache, minifyRecipe(p.Name(), nil))

	return nil
}

// findJSFiles recursively finds all JS files in a directory.
func (p *JSMinifyPlugin) findJSFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Check for JS files (exclude .min.js which are already minified)
		lower := strings.ToLower(path)
		if strings.HasSuffix(lower, ".js") && !strings.HasSuffix(lower, ".min.js") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// isExcluded checks if a file should be excluded from minification.
func (p *JSMinifyPlugin) isExcluded(filePath string) bool {
	if isPagefindAsset(filePath) {
		return true
	}

	filename := filepath.Base(filePath)

	// Always skip already-minified files
	if strings.HasSuffix(strings.ToLower(filename), ".min.js") {
		return true
	}

	return isExcludedByPatterns(filename, p.exclude)
}

func isPagefindAsset(filePath string) bool {
	normalized := filepath.ToSlash(filePath)
	return strings.Contains(normalized, "/_pagefind/")
}

// minifyFile minifies a single JS file in place.
// Returns the original size and minified size.
// Uses a sync.Pool for buffer reuse to reduce allocations under concurrent workloads.
func (p *JSMinifyPlugin) minifyFile(filePath string) (original, minified int64, err error) {
	// Read the original file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("reading file: %w", err)
	}

	original = int64(len(content))
	result, err := p.minifyBytes(content)
	if err != nil {
		return original, 0, err
	}
	if original == 0 {
		return 0, 0, nil
	}
	if err := writeGeneratedFile(filePath, result); err != nil {
		return original, 0, fmt.Errorf("writing minified file: %w", err)
	}
	return original, int64(len(result)), nil
}

func (p *JSMinifyPlugin) minifyBytes(content []byte) ([]byte, error) {
	// Skip empty files
	if len(content) == 0 {
		return []byte{}, nil
	}

	// Minify the content using a pooled buffer
	buf := jsBufPool.Get().(*bytes.Buffer) //nolint:errcheck // type is guaranteed by pool's New func
	buf.Reset()
	defer jsBufPool.Put(buf)

	if err := p.minifier.Minify("application/javascript", buf, bytes.NewReader(content)); err != nil {
		return nil, fmt.Errorf("minifying: %w", err)
	}

	return bytes.Clone(buf.Bytes()), nil
}

// Priority returns the plugin priority for the write stage.
// Should run very late (after all JS-generating plugins) to minify all JS.
func (p *JSMinifyPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		// Run at the very end of Write stage, alongside CSS minification
		return lifecycle.PriorityLast
	}
	return lifecycle.PriorityDefault
}

// Ensure JSMinifyPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*JSMinifyPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*JSMinifyPlugin)(nil)
	_ lifecycle.WritePlugin     = (*JSMinifyPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*JSMinifyPlugin)(nil)
)
