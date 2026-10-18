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
	"github.com/tdewolff/minify/v2/css"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// cssBufPool reuses bytes.Buffer instances across CSS minification calls to reduce allocations.
var cssBufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

const cssMinifyPluginName = "css_minify"

// CSSMinifyPlugin minifies CSS files to reduce file sizes and improve
// Lighthouse performance scores. It runs during the Write stage, after
// all other CSS-generating plugins (palette_css, chroma_css, css_bundle)
// have completed.
//
// The plugin:
// - Discovers CSS files in the output directory
// - Minifies each file using tdewolff/minify
// - Preserves specified comments (e.g., copyright notices)
// - Skips excluded files
// - Reports size reduction statistics
type CSSMinifyPlugin struct {
	config   models.CSSMinifyConfig
	minifier *minify.M
	exclude  map[string]bool // excluded file patterns for fast lookup
}

// NewCSSMinifyPlugin creates a new CSSMinifyPlugin.
func NewCSSMinifyPlugin() *CSSMinifyPlugin {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)

	return &CSSMinifyPlugin{
		minifier: m,
		exclude:  make(map[string]bool),
	}
}

// Name returns the unique name of the plugin.
func (p *CSSMinifyPlugin) Name() string {
	return cssMinifyPluginName
}

// Configure reads the CSS minify configuration from the manager's config.
func (p *CSSMinifyPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()

	// Set default config first
	p.config = models.NewCSSMinifyConfig()
	p.exclude = make(map[string]bool)

	if config.Extra == nil {
		return nil
	}

	// Try to get css_minify config from Extra
	// It may be a models.CSSMinifyConfig or a map[string]interface{} from TOML parsing
	switch v := config.Extra[cssMinifyPluginName].(type) {
	case models.CSSMinifyConfig:
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

// parseConfigFromMap parses CSSMinifyConfig from a raw map (TOML parsing result).
func (p *CSSMinifyPlugin) parseConfigFromMap(m map[string]interface{}) models.CSSMinifyConfig {
	cfg := models.NewCSSMinifyConfig()

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

	// Parse preserve_comments list
	if preserveComments, ok := m["preserve_comments"].([]interface{}); ok {
		cfg.PreserveComments = make([]string, 0, len(preserveComments))
		for _, c := range preserveComments {
			if s, ok := c.(string); ok {
				cfg.PreserveComments = append(cfg.PreserveComments, s)
			}
		}
	}

	return cfg
}

// Write performs CSS minification in the output directory.
// Skipped in fast mode (--fast flag) for faster development builds.
func (p *CSSMinifyPlugin) Write(m *lifecycle.Manager) error {
	if !p.config.Enabled {
		return nil
	}

	// Skip minification in fast mode
	if fast, ok := m.Config().Extra["fast_mode"].(bool); ok && fast {
		return nil
	}

	outputDir := m.Config().OutputDir
	cssFiles, err := p.findCSSFiles(outputDir)
	if err != nil {
		return fmt.Errorf("finding CSS files: %w", err)
	}

	cache, cacheErr := newMinifyCache(m.Config(), p.Name())
	if cacheErr != nil {
		log.Printf("[%s] Warning: initializing asset cache: %v", p.Name(), cacheErr)
	}
	defer cache.close()
	runMinification(p.Name(), cssFiles, p.isExcluded, p.minifyBytes, m.Concurrency(), cache, minifyRecipe(p.Name(), p.config.PreserveComments))

	return nil
}

// findCSSFiles recursively finds all CSS files in a directory.
func (p *CSSMinifyPlugin) findCSSFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// Check for CSS files
		if strings.HasSuffix(strings.ToLower(path), ".css") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// isExcluded checks if a file should be excluded from minification.
func (p *CSSMinifyPlugin) isExcluded(filePath string) bool {
	if isPagefindAsset(filePath) {
		return true
	}

	filename := filepath.Base(filePath)
	// The generated font stylesheet lives at css/fonts.css. Minifying it
	// again can rewrite unicode-range declarations on incremental builds.
	// Do not exclude unrelated stylesheets just because they share its name.
	if filename == "fonts.css" && filepath.Base(filepath.Dir(filePath)) == "css" {
		return true
	}
	return isExcludedByPatterns(filename, p.exclude)
}

// minifyFile minifies a single CSS file in place.
// Returns the original size and minified size.
// Uses a sync.Pool for buffer reuse to reduce allocations under concurrent workloads.
func (p *CSSMinifyPlugin) minifyFile(filePath string) (original, minified int64, err error) {
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
	if err := writeGeneratedFile(filePath, result); err != nil {
		return original, 0, fmt.Errorf("writing minified file: %w", err)
	}
	return original, int64(len(result)), nil
}

func (p *CSSMinifyPlugin) minifyBytes(content []byte) ([]byte, error) {
	// Preserve specified comments
	preservedComments := p.extractPreservedComments(string(content))

	// Minify the content using a pooled buffer
	buf := cssBufPool.Get().(*bytes.Buffer) //nolint:errcheck // type is guaranteed by pool's New func
	buf.Reset()
	defer cssBufPool.Put(buf)

	if err := p.minifier.Minify("text/css", buf, bytes.NewReader(content)); err != nil {
		return nil, fmt.Errorf("minifying: %w", err)
	}

	// Build result: prepend preserved comments + minified content
	// Use a single pooled buffer for the result to avoid a second allocation
	result := cssBufPool.Get().(*bytes.Buffer) //nolint:errcheck // type is guaranteed by pool's New func
	result.Reset()
	defer cssBufPool.Put(result)

	for _, comment := range preservedComments {
		result.WriteString(comment)
		result.WriteByte('\n')
	}
	result.Write(buf.Bytes())

	return bytes.Clone(result.Bytes()), nil
}

// extractPreservedComments extracts comments that match the preserve patterns.
func (p *CSSMinifyPlugin) extractPreservedComments(content string) []string {
	if len(p.config.PreserveComments) == 0 {
		return nil
	}

	var preserved []string

	// Find all CSS comments
	i := 0
	for i < len(content) {
		// Look for comment start
		start := strings.Index(content[i:], "/*")
		if start == -1 {
			break
		}
		start += i

		// Look for comment end
		end := strings.Index(content[start+2:], "*/")
		if end == -1 {
			break
		}
		end += start + 2 + 2 // Include the */

		comment := content[start:end]

		// Check if this comment should be preserved
		for _, pattern := range p.config.PreserveComments {
			if strings.Contains(comment, pattern) {
				preserved = append(preserved, comment)
				break
			}
		}

		i = end
	}

	return preserved
}

// Priority returns the plugin priority for the write stage.
// Should run late (after palette_css, chroma_css, css_bundle) to minify all CSS.
func (p *CSSMinifyPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		// Run very late in Write stage, after CSS generators and bundler
		return lifecycle.PriorityLast
	}
	return lifecycle.PriorityDefault
}

// Ensure CSSMinifyPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*CSSMinifyPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*CSSMinifyPlugin)(nil)
	_ lifecycle.WritePlugin     = (*CSSMinifyPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*CSSMinifyPlugin)(nil)
)
