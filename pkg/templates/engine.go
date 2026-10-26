// Package templates provides a template engine for markata-go using pongo2 (Jinja2-like syntax).
package templates

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/themes"

	"github.com/flosch/pongo2/v6"
)

// Constants for theme names and embedded file markers
const (
	defaultThemeName   = "default"
	embeddedFilePrefix = "embedded:"
)

// Engine provides template rendering capabilities using pongo2.
type Engine struct {
	set *pongo2.TemplateSet
	dir string
	mu  sync.RWMutex

	// templateCache stores compiled templates for reuse
	templateCache map[string]*pongo2.Template

	// searchPaths is the ordered list of directories to search for templates
	// Resolution order: project templates -> theme templates -> default theme
	searchPaths []string

	// themeName is the current theme name
	themeName string

	// embeddedFS holds the embedded default theme templates as fallback
	embeddedFS fs.FS

	// useEmbedded indicates whether to use embedded templates as fallback
	useEmbedded bool
}

// NewEngine creates a new template engine with the given templates directory.
// The templates directory is used as the base path for template loading.
// If templatesDir is empty, templates can only be rendered from strings.
func NewEngine(templatesDir string) (*Engine, error) {
	return NewEngineWithTheme(templatesDir, "default")
}

// NewEngineWithTheme creates a new template engine with theme support.
func NewEngineWithTheme(templatesDir, themeName string) (*Engine, error) {
	e := &Engine{
		dir:           templatesDir,
		templateCache: make(map[string]*pongo2.Template),
		searchPaths:   make([]string, 0),
		themeName:     themeName,
		embeddedFS:    themes.DefaultTemplates(),
		useEmbedded:   true,
	}

	if e.themeName == "" {
		e.themeName = defaultThemeName
	}

	// Register custom filters
	registerFilters()

	// Build search paths
	e.buildSearchPaths(templatesDir)

	// Create a basic template set (we'll handle loading ourselves)
	e.set = pongo2.NewSet(defaultThemeName, pongo2.MustNewLocalFileSystemLoader(""))

	return e, nil
}

// buildSearchPaths constructs the ordered list of template directories.
func (e *Engine) buildSearchPaths(templatesDir string) {
	e.searchPaths = make([]string, 0)

	// 1. Project templates (highest priority)
	if templatesDir != "" {
		if _, err := os.Stat(templatesDir); err == nil {
			e.searchPaths = append(e.searchPaths, templatesDir)
		}
	}

	// 2. Current theme templates in current directory
	if e.themeName != "" && e.themeName != defaultThemeName {
		themeTemplatesDir := filepath.Join("themes", e.themeName, "templates")
		if _, err := os.Stat(themeTemplatesDir); err == nil {
			e.searchPaths = append(e.searchPaths, themeTemplatesDir)
		}
	}

	// 3. Default theme templates in current directory
	defaultThemeDir := filepath.Join("themes", defaultThemeName, "templates")
	if _, err := os.Stat(defaultThemeDir); err == nil {
		e.searchPaths = append(e.searchPaths, defaultThemeDir)
	}

	// 4. Try themes relative to executable directory
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)

		// Check for themes next to executable (e.g., /usr/local/bin/markata-go -> /usr/local/bin/themes)
		if e.themeName != "" && e.themeName != defaultThemeName {
			themeDir := filepath.Join(exeDir, "themes", e.themeName, "templates")
			if _, err := os.Stat(themeDir); err == nil {
				e.searchPaths = append(e.searchPaths, themeDir)
			}
		}

		defaultExeThemeDir := filepath.Join(exeDir, "themes", defaultThemeName, "templates")
		if _, err := os.Stat(defaultExeThemeDir); err == nil {
			e.searchPaths = append(e.searchPaths, defaultExeThemeDir)
		}

		// Check parent directory (e.g., /usr/local/bin -> /usr/local/share/markata-go/themes)
		parentDir := filepath.Dir(exeDir)
		shareThemeDir := filepath.Join(parentDir, "share", "markata-go", "themes", "default", "templates")
		if _, err := os.Stat(shareThemeDir); err == nil {
			e.searchPaths = append(e.searchPaths, shareThemeDir)
		}
	}

	// If no paths found, use the original templatesDir even if it doesn't exist
	// (embedded templates will be used as final fallback)
	if len(e.searchPaths) == 0 && templatesDir != "" {
		e.searchPaths = append(e.searchPaths, templatesDir)
	}
}

// Render renders a template file with the given context.
// templateName is the path relative to the templates directory.
func (e *Engine) Render(templateName string, ctx Context) (string, error) {
	tpl, err := e.LoadTemplate(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to load template %q: %w", templateName, err)
	}

	result, err := tpl.Execute(ctx.ToPongo2())
	if err != nil {
		return "", fmt.Errorf("failed to execute template %q: %w", templateName, err)
	}

	return result, nil
}

// RenderString renders a template string with the given context.
// This is useful for inline templates in markdown content (jinja_md).
func (e *Engine) RenderString(templateStr string, ctx Context) (string, error) {
	tpl, err := pongo2.FromString(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template string: %w", err)
	}

	result, err := tpl.Execute(ctx.ToPongo2())
	if err != nil {
		return "", fmt.Errorf("failed to execute template string: %w", err)
	}

	return result, nil
}

// LoadTemplate loads and caches a template by name.
// The template is loaded from the search paths in order, with embedded templates as fallback.
func (e *Engine) LoadTemplate(name string) (*pongo2.Template, error) {
	e.mu.RLock()
	if tpl, ok := e.templateCache[name]; ok {
		e.mu.RUnlock()
		return tpl, nil
	}
	e.mu.RUnlock()

	// Find the template in search paths
	templatePath := e.FindTemplate(name)

	var tpl *pongo2.Template
	var err error

	switch {
	case templatePath != "":
		// Found in filesystem - use file-based loading
		absPath, absErr := filepath.Abs(templatePath)
		if absErr != nil {
			absPath = templatePath
		}

		// Create a template set with a multi-directory loader for proper include/extends support
		tplSet := pongo2.NewSet(name, &searchPathLoader{
			searchPaths: e.searchPaths,
			embeddedFS:  e.embeddedFS,
		})

		// Load the template using absolute path
		tpl, err = tplSet.FromFile(absPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %q: %w", name, err)
		}
	case e.useEmbedded && e.embeddedFS != nil:
		// Try embedded templates as fallback
		content, readErr := fs.ReadFile(e.embeddedFS, name)
		if readErr != nil {
			return nil, fmt.Errorf("template %q not found in search paths %v or embedded templates", name, e.searchPaths)
		}

		// Create a template set with embedded loader for include/extends support
		tplSet := pongo2.NewSet(name, &searchPathLoader{
			searchPaths: e.searchPaths,
			embeddedFS:  e.embeddedFS,
		})

		tpl, err = tplSet.FromBytes(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse embedded template %q: %w", name, err)
		}
	default:
		return nil, fmt.Errorf("template %q not found in search paths %v", name, e.searchPaths)
	}

	// Cache the compiled template
	e.mu.Lock()
	e.templateCache[name] = tpl
	e.mu.Unlock()

	return tpl, nil
}

// searchPathLoader implements pongo2.TemplateLoader for multi-directory support.
type searchPathLoader struct {
	searchPaths []string
	embeddedFS  fs.FS
}

func (l *searchPathLoader) Abs(base, name string) string {
	// If name is already absolute, return it
	if filepath.IsAbs(name) {
		return name
	}

	// If name already has embedded prefix, return it
	if strings.HasPrefix(name, embeddedFilePrefix) {
		return name
	}

	// Handle case where base is an embedded file
	if strings.HasPrefix(base, embeddedFilePrefix) {
		if resolved := resolveEmbeddedInclude(base, name, l.searchPaths, l.embeddedFS); resolved != "" {
			return resolved
		}
	}

	// If base is provided and is a regular file, try relative to base directory first
	if base != "" && !strings.HasPrefix(base, embeddedFilePrefix) {
		baseDir := filepath.Dir(base)
		candidate := filepath.Join(baseDir, name)
		if _, err := os.Stat(candidate); err == nil {
			absPath, absErr := filepath.Abs(candidate)
			if absErr == nil {
				return absPath
			}
			return candidate
		}
	}

	// Search through all search paths
	for _, dir := range l.searchPaths {
		candidate := filepath.Join(dir, name)
		if _, err := os.Stat(candidate); err == nil {
			absPath, absErr := filepath.Abs(candidate)
			if absErr == nil {
				return absPath
			}
			return candidate
		}
	}

	// Check embedded filesystem
	if l.embeddedFS != nil {
		if _, err := fs.Stat(l.embeddedFS, name); err == nil {
			// Return a special marker for embedded files
			return embeddedFilePrefix + name
		}
	}

	// Return name with first search path as fallback (will fail gracefully)
	if len(l.searchPaths) > 0 {
		absPath, absErr := filepath.Abs(filepath.Join(l.searchPaths[0], name))
		if absErr == nil {
			return absPath
		}
	}
	return name
}

func resolveEmbeddedInclude(base, name string, searchPaths []string, embeddedFS fs.FS) string {
	embeddedBase := base[len(embeddedFilePrefix):]
	baseDir := filepath.Dir(embeddedBase)
	// Project templates override embedded templates, including nested
	// includes from an embedded parent template.
	projectNames := []string{name}
	if baseDir != "." {
		projectNames = append(projectNames, filepath.Join(baseDir, name))
	}
	if projectPath := findProjectTemplate(searchPaths, projectNames...); projectPath != "" {
		return projectPath
	}

	if embeddedFS == nil {
		return ""
	}
	if _, err := fs.Stat(embeddedFS, name); err == nil {
		return embeddedFilePrefix + name
	}
	if baseDir != "." {
		candidate := filepath.Join(baseDir, name)
		if _, err := fs.Stat(embeddedFS, candidate); err == nil {
			return embeddedFilePrefix + candidate
		}
	}
	return ""
}

func findProjectTemplate(searchPaths []string, names ...string) string {
	for _, dir := range searchPaths {
		for _, name := range names {
			candidate := filepath.Join(dir, name)
			if _, err := os.Stat(candidate); err != nil {
				continue
			}
			absPath, err := filepath.Abs(candidate)
			if err == nil {
				return absPath
			}
			return candidate
		}
	}
	return ""
}

func (l *searchPathLoader) Get(path string) (io.Reader, error) {
	// Check for embedded file marker
	if strings.HasPrefix(path, embeddedFilePrefix) {
		embeddedPath := path[len(embeddedFilePrefix):]
		if l.embeddedFS != nil {
			content, err := fs.ReadFile(l.embeddedFS, embeddedPath)
			if err == nil {
				return bytes.NewReader(content), nil
			}
		}
	}

	// If path is already absolute, just try to open it
	if filepath.IsAbs(path) {
		reader, err := readTemplateFile(path)
		if err == nil {
			return reader, nil
		}
	}

	// Search through all search paths
	for _, dir := range l.searchPaths {
		candidate := filepath.Join(dir, path)
		reader, err := readTemplateFile(candidate)
		if err == nil {
			return reader, nil
		}
	}

	// Try embedded filesystem
	if l.embeddedFS != nil {
		content, err := fs.ReadFile(l.embeddedFS, path)
		if err == nil {
			return bytes.NewReader(content), nil
		}
	}

	// Last resort: try to open path as-is
	reader, err := readTemplateFile(path)
	if err != nil {
		return nil, fmt.Errorf("template %q not found in search paths or embedded", path)
	}
	return reader, nil
}

// readTemplateFile reads and closes a filesystem template before returning it.
// Pongo2 consumes the returned reader but does not close it, so returning an
// in-memory reader prevents file-handle leaks and allows template replacement
// on platforms that lock open files.
func readTemplateFile(path string) (io.Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	content, readErr := io.ReadAll(file)
	closeErr := file.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	return bytes.NewReader(content), nil
}

// TemplateExists checks if a template file exists in any search path or embedded templates.
func (e *Engine) TemplateExists(name string) bool {
	// Check filesystem paths
	for _, dir := range e.searchPaths {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	// Check embedded templates
	if e.useEmbedded && e.embeddedFS != nil {
		if _, err := fs.Stat(e.embeddedFS, name); err == nil {
			return true
		}
	}

	return false
}

// FindTemplate returns the full path to the template, searching through all paths.
// Returns empty string if not found in filesystem (embedded templates don't have paths).
func (e *Engine) FindTemplate(name string) string {
	for _, dir := range e.searchPaths {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

// HasEmbeddedTemplate checks if a template exists in embedded templates.
func (e *Engine) HasEmbeddedTemplate(name string) bool {
	if e.embeddedFS == nil {
		return false
	}
	_, err := fs.Stat(e.embeddedFS, name)
	return err == nil
}

// ClearCache clears the template cache.
// This is useful during development when templates are being modified.
func (e *Engine) ClearCache() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.templateCache = make(map[string]*pongo2.Template)
}

// Dir returns the templates directory path.
func (e *Engine) Dir() string {
	return e.dir
}

// SearchPaths returns the ordered list of template search paths.
func (e *Engine) SearchPaths() []string {
	return e.searchPaths
}

// ThemeName returns the current theme name.
func (e *Engine) ThemeName() string {
	return e.themeName
}

// SetTheme sets the theme and rebuilds search paths.
func (e *Engine) SetTheme(themeName string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.themeName = themeName
	if e.themeName == "" {
		e.themeName = defaultThemeName
	}

	// Clear cache and rebuild search paths
	e.templateCache = make(map[string]*pongo2.Template)
	e.buildSearchPaths(e.dir)
}

// SetDir sets the templates directory and reinitializes the loader.
func (e *Engine) SetDir(dir string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.dir = dir
	e.templateCache = make(map[string]*pongo2.Template)

	// Rebuild search paths
	e.buildSearchPaths(dir)

	return nil
}

// RenderToString renders a template by name with a raw map context.
// This is useful for rendering templates from plugins that don't use the full Context type.
// The ctx map is converted to a pongo2.Context for template execution.
func (e *Engine) RenderToString(templateName string, ctx map[string]interface{}) (string, error) {
	tpl, err := e.LoadTemplate(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to load template %q: %w", templateName, err)
	}

	// Convert map to pongo2.Context
	pctx := pongo2.Context(ctx)

	result, err := tpl.Execute(pctx)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %q: %w", templateName, err)
	}

	return result, nil
}
