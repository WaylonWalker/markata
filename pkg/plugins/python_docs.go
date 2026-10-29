// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	markatatemplates "github.com/WaylonWalker/markata/pkg/templates"
	"github.com/bmatcuk/doublestar/v4"
)

const pythonDocsPluginName = "python_docs"

// PythonDocsConfig holds configuration for the python_docs plugin.
type PythonDocsConfig struct {
	Enabled            bool     `json:"enabled" yaml:"enabled" toml:"enabled"`
	Patterns           []string `json:"patterns" yaml:"patterns" toml:"patterns"`
	Directories        []string `json:"directories" yaml:"directories" toml:"directories"`
	Exclude            []string `json:"exclude" yaml:"exclude" toml:"exclude"`
	UseGitignore       bool     `json:"use_gitignore" yaml:"use_gitignore" toml:"use_gitignore"`
	SlugPrefix         string   `json:"slug_prefix" yaml:"slug_prefix" toml:"slug_prefix"`
	Template           string   `json:"template" yaml:"template" toml:"template"`
	Published          bool     `json:"published" yaml:"published" toml:"published"`
	IncludePrivate     bool     `json:"include_private" yaml:"include_private" toml:"include_private"`
	IncludeSource      bool     `json:"include_source" yaml:"include_source" toml:"include_source"`
	IncludeModuleCode  bool     `json:"include_module_code" yaml:"include_module_code" toml:"include_module_code"`
	SymbolTemplate     string   `json:"symbol_template" yaml:"symbol_template" toml:"symbol_template"`
	SymbolTemplateName string   `json:"symbol_template_name" yaml:"symbol_template_name" toml:"symbol_template_name"`
	Tags               []string `json:"tags" yaml:"tags" toml:"tags"`
	Interpreter        string   `json:"interpreter" yaml:"interpreter" toml:"interpreter"`
}

func defaultPythonDocsConfig() PythonDocsConfig {
	return PythonDocsConfig{
		Enabled:           false,
		Patterns:          []string{"**/*.py"},
		Exclude:           []string{"**/.venv/**", "**/venv/**", "**/__pycache__/**", "**/site-packages/**", "**/node_modules/**"},
		UseGitignore:      true,
		SlugPrefix:        "api",
		Published:         false,
		IncludePrivate:    false,
		IncludeSource:     true,
		IncludeModuleCode: false,
		Tags:              []string{"python", "docs"},
	}
}

// PythonDocsPlugin generates docs posts from Python source files.
type PythonDocsPlugin struct {
	config            PythonDocsConfig
	interpreter       string
	gitignorePatterns []string
	strictHookOptIn   bool
}

// NewPythonDocsPlugin creates a new PythonDocsPlugin.
func NewPythonDocsPlugin() *PythonDocsPlugin {
	return &PythonDocsPlugin{config: defaultPythonDocsConfig()}
}

// Name returns the unique plugin name.
func (p *PythonDocsPlugin) Name() string {
	return pythonDocsPluginName
}

// Configure reads plugin configuration and resolves the Python interpreter when enabled.
func (p *PythonDocsPlugin) Configure(m *lifecycle.Manager) error {
	p.config = defaultPythonDocsConfig()
	config := m.Config()
	p.strictHookOptIn = pythonDocsExplicitlyEnabled(config)
	if config.Extra != nil {
		if raw, ok := config.Extra[pythonDocsPluginName].(map[string]interface{}); ok {
			p.parseConfig(raw)
		}
	}

	if !p.config.Enabled || !p.strictHookOptIn {
		return nil
	}

	baseDir := config.ContentDir
	if baseDir == "" {
		baseDir = "."
	}

	if p.config.UseGitignore {
		if err := p.loadGitignore(baseDir); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("%s: load gitignore: %w", pythonDocsPluginName, err)
		}
	}

	interpreter, err := p.resolveInterpreter()
	if err != nil {
		return fmt.Errorf("%s: %w", pythonDocsPluginName, err)
	}
	p.interpreter = interpreter

	return nil
}

// Load discovers Python files and appends generated docs posts.
func (p *PythonDocsPlugin) Load(m *lifecycle.Manager) error {
	if !p.config.Enabled || !p.strictHookOptIn {
		return nil
	}

	baseDir := m.Config().ContentDir
	if baseDir == "" {
		baseDir = "."
	}
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return fmt.Errorf("%s: resolve content dir: %w", pythonDocsPluginName, err)
	}

	files, modTimes, err := p.scanPythonFiles(absBaseDir)
	if err != nil {
		return fmt.Errorf("%s: scan files: %w", pythonDocsPluginName, err)
	}
	if len(files) == 0 {
		return nil
	}

	modules, err := p.extractModules(absBaseDir, files)
	if err != nil {
		return fmt.Errorf("%s: extract module docs: %w", pythonDocsPluginName, err)
	}

	renderer, err := newPythonSymbolRenderer(m.Config(), p.config)
	if err != nil {
		return fmt.Errorf("%s: initialize symbol renderer: %w", pythonDocsPluginName, err)
	}

	moduleIndex, symbolIndex := buildPythonDocIndexes(p.config.SlugPrefix, modules)
	for i := range modules {
		module := modules[i]
		post, err := p.makePost(module, modTimes[module.SourcePath], moduleIndex, symbolIndex, renderer)
		if err != nil {
			return fmt.Errorf("%s: build post for %s: %w", pythonDocsPluginName, module.ModuleName, err)
		}
		m.AddPost(post)
	}

	return nil
}

func (p *PythonDocsPlugin) parseConfig(cfg map[string]interface{}) {
	if enabled, ok := cfg["enabled"].(bool); ok {
		p.config.Enabled = enabled
	}
	if patterns := pythonDocsStringSlice(cfg["patterns"]); len(patterns) > 0 {
		p.config.Patterns = patterns
	}
	if directories := pythonDocsStringSlice(cfg["directories"]); len(directories) > 0 {
		p.config.Directories = directories
	}
	if directories := pythonDocsStringSlice(cfg["content_directories"]); len(directories) > 0 {
		p.config.Directories = directories
	}
	if exclude := pythonDocsStringSlice(cfg["exclude"]); len(exclude) > 0 {
		p.config.Exclude = exclude
	}
	if useGitignore, ok := cfg["use_gitignore"].(bool); ok {
		p.config.UseGitignore = useGitignore
	}
	if slugPrefix, ok := cfg["slug_prefix"].(string); ok && slugPrefix != "" {
		p.config.SlugPrefix = slugPrefix
	}
	if template, ok := cfg["template"].(string); ok {
		p.config.Template = template
	}
	if published, ok := cfg["published"].(bool); ok {
		p.config.Published = published
	}
	if includePrivate, ok := cfg["include_private"].(bool); ok {
		p.config.IncludePrivate = includePrivate
	}
	if includeSource, ok := cfg["include_source"].(bool); ok {
		p.config.IncludeSource = includeSource
	}
	if includeModuleCode, ok := cfg["include_module_code"].(bool); ok {
		p.config.IncludeModuleCode = includeModuleCode
	}
	if symbolTemplate, ok := cfg["symbol_template"].(string); ok {
		p.config.SymbolTemplate = symbolTemplate
	}
	if symbolTemplateName, ok := cfg["symbol_template_name"].(string); ok {
		p.config.SymbolTemplateName = symbolTemplateName
	}
	if tags := pythonDocsStringSlice(cfg["tags"]); len(tags) > 0 {
		p.config.Tags = tags
	}
	if interpreter, ok := cfg["interpreter"].(string); ok && interpreter != "" {
		p.config.Interpreter = interpreter
	}
}

func (p *PythonDocsPlugin) resolveInterpreter() (string, error) {
	if p.config.Interpreter != "" {
		if _, err := exec.LookPath(p.config.Interpreter); err != nil {
			return "", fmt.Errorf("configured interpreter %q not found", p.config.Interpreter)
		}
		return p.config.Interpreter, nil
	}

	for _, candidate := range []string{"python3", "python"} {
		if path, err := exec.LookPath(candidate); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("python interpreter not found; set [markata-go.python_docs].interpreter or install python3")
}

func (p *PythonDocsPlugin) scanPythonFiles(absBaseDir string) (files []string, modTimes map[string]time.Time, err error) {
	patterns := append([]string{}, p.config.Patterns...)
	for _, dir := range p.config.Directories {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		patterns = append(patterns, filepath.ToSlash(filepath.Join(dir, "**", "*.py")))
	}

	fileSet := make(map[string]struct{})
	modTimes = make(map[string]time.Time)

	for _, pattern := range patterns {
		fullPattern := pattern
		if !filepath.IsAbs(pattern) {
			fullPattern = filepath.Join(absBaseDir, pattern)
		}

		matches, err := doublestar.FilepathGlob(fullPattern)
		if err != nil {
			return nil, nil, fmt.Errorf("glob python files %q: %w", pattern, err)
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil || info.IsDir() {
				continue
			}

			relPath, err := filepath.Rel(absBaseDir, match)
			if err != nil {
				continue
			}
			relPath = filepath.ToSlash(relPath)

			if p.isIgnored(relPath) || p.isExcluded(relPath) {
				continue
			}

			fileSet[relPath] = struct{}{}
			modTimes[relPath] = info.ModTime().UTC()
		}
	}

	files = make([]string, 0, len(fileSet))
	for file := range fileSet {
		files = append(files, file)
	}
	sort.Strings(files)

	return files, modTimes, nil
}

func (p *PythonDocsPlugin) loadGitignore(baseDir string) error {
	gitignorePath := filepath.Join(baseDir, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return err
	}
	defer file.Close()

	p.gitignorePatterns = p.gitignorePatterns[:0]
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		p.gitignorePatterns = append(p.gitignorePatterns, line)
	}

	return scanner.Err()
}

func (p *PythonDocsPlugin) isIgnored(path string) bool {
	if !p.config.UseGitignore || len(p.gitignorePatterns) == 0 {
		return false
	}

	normalizedPath := filepath.ToSlash(path)
	for _, pattern := range p.gitignorePatterns {
		if strings.HasPrefix(pattern, "!") {
			continue
		}
		normalizedPattern := strings.TrimSuffix(filepath.ToSlash(pattern), "/")

		matched, err := doublestar.Match(normalizedPattern, normalizedPath)
		if err == nil && matched {
			return true
		}

		if strings.HasPrefix(normalizedPath, normalizedPattern+"/") {
			return true
		}

		filename := filepath.Base(normalizedPath)
		matched, err = doublestar.Match(normalizedPattern, filename)
		if err == nil && matched {
			return true
		}

		if !strings.HasPrefix(normalizedPattern, "**/") && !strings.HasPrefix(normalizedPattern, "/") {
			matched, err = doublestar.Match("**/"+normalizedPattern, normalizedPath)
			if err == nil && matched {
				return true
			}
		}
	}

	return false
}

func (p *PythonDocsPlugin) isExcluded(path string) bool {
	normalizedPath := filepath.ToSlash(path)
	for _, pattern := range p.config.Exclude {
		matched, err := doublestar.Match(filepath.ToSlash(pattern), normalizedPath)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func (p *PythonDocsPlugin) extractModules(baseDir string, files []string) ([]pythonModuleDoc, error) {
	input, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(p.interpreter, "-c", pythonDocsExtractorScript, baseDir) // #nosec G204 -- interpreter is resolved from user config and LookPath
	cmd.Stdin = bytes.NewReader(input)
	cmd.Env = append(os.Environ(), "MARKATA_GO_PYTHON_DOCS_INCLUDE_PRIVATE="+pythonDocsBoolEnv(p.config.IncludePrivate))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("python extraction failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}

	var modules []pythonModuleDoc
	if err := json.Unmarshal(stdout.Bytes(), &modules); err != nil {
		return nil, fmt.Errorf("decode extractor output: %w", err)
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].ModuleName < modules[j].ModuleName
	})

	return modules, nil
}

func buildPythonDocIndexes(slugPrefix string, modules []pythonModuleDoc) (moduleIndex, symbolIndex map[string]string) {
	moduleIndex = make(map[string]string, len(modules))
	symbolIndex = make(map[string]string)

	for i := range modules {
		module := &modules[i]
		href := "/" + buildPythonDocSlug(slugPrefix, module.ModuleName) + "/"
		moduleIndex[module.ModuleName] = href
		symbolIndex[module.ModuleName] = href

		for classIdx := range module.Classes {
			classDoc := &module.Classes[classIdx]
			classAnchor := pythonDocAnchor(classDoc.Name)
			symbolIndex[module.ModuleName+"."+classDoc.Name] = href + "#" + classAnchor
			for methodIdx := range classDoc.Methods {
				method := &classDoc.Methods[methodIdx]
				symbolIndex[module.ModuleName+"."+classDoc.Name+"."+method.Name] = href + "#" + pythonDocAnchor(classDoc.Name+"-"+method.Name)
			}
		}
		for functionIdx := range module.Functions {
			function := &module.Functions[functionIdx]
			symbolIndex[module.ModuleName+"."+function.Name] = href + "#" + pythonDocAnchor(function.Name)
		}
	}

	return moduleIndex, symbolIndex
}

func (p *PythonDocsPlugin) makePost(module pythonModuleDoc, modTime time.Time, moduleIndex, symbolIndex map[string]string, renderer *pythonSymbolRenderer) (*models.Post, error) {
	title := module.ModuleName
	post := models.NewPost(module.SourcePath)
	post.Title = &title
	content, err := renderPythonModuleMarkdown(module, p.config, moduleIndex, symbolIndex, renderer)
	if err != nil {
		return nil, err
	}
	post.Content = content
	post.Template = p.config.Template
	post.Published = p.config.Published
	post.Tags = append([]string{}, p.config.Tags...)
	post.Slug = buildPythonDocSlug(p.config.SlugPrefix, module.ModuleName)
	post.GenerateHref()
	post.Date = &modTime
	post.Modified = &modTime

	if description := firstParagraph(module.Docstring); description != "" {
		post.Description = &description
	}

	post.Set(pythonDocsPluginName, true)
	post.Set("python_module", module.ModuleName)
	post.Set("source_path", module.SourcePath)

	return post, nil
}

func buildPythonDocSlug(prefix, moduleName string) string {
	modulePath := strings.ReplaceAll(moduleName, ".", "/")
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return strings.Trim(modulePath, "/")
	}
	if modulePath == "" {
		return prefix
	}
	return prefix + "/" + modulePath
}

func renderPythonModuleMarkdown(module pythonModuleDoc, cfg PythonDocsConfig, moduleIndex, symbolIndex map[string]string, renderer *pythonSymbolRenderer) (string, error) {
	var b strings.Builder

	b.WriteString("# `" + module.ModuleName + "`\n\n")
	b.WriteString("Source: `" + module.SourcePath + "`\n\n")

	if module.Docstring != "" {
		b.WriteString(linkPythonDocstring(module.Docstring, moduleIndex, symbolIndex))
		b.WriteString("\n\n")
	}

	if len(module.Imports) > 0 {
		b.WriteString("## Imports\n\n")
		for _, imp := range module.Imports {
			b.WriteString("- ")
			b.WriteString(renderPythonImport(imp, moduleIndex, symbolIndex))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(module.Classes) > 0 || len(module.Functions) > 0 {
		b.WriteString("## API\n\n")
		for i := range module.Classes {
			classDoc := module.Classes[i]
			b.WriteString("- [Class `" + classDoc.Name + "`](#" + pythonDocAnchor(classDoc.Name) + ")\n")
		}
		for i := range module.Functions {
			function := module.Functions[i]
			b.WriteString("- [Function `" + function.Name + "`](#" + pythonDocAnchor(function.Name) + ")\n")
		}
		b.WriteString("\n")
	}

	for i := range module.Classes {
		classDoc := module.Classes[i]
		if err := renderPythonDocItem(&b, module, classDoc, 2, cfg, moduleIndex, symbolIndex, classDoc.Name, renderer); err != nil {
			return "", err
		}
	}
	for i := range module.Functions {
		function := module.Functions[i]
		if err := renderPythonDocItem(&b, module, function, 2, cfg, moduleIndex, symbolIndex, function.Name, renderer); err != nil {
			return "", err
		}
	}

	if cfg.IncludeSource && cfg.IncludeModuleCode && module.Source != "" {
		b.WriteString("## Module Source\n\n")
		b.WriteString("```python\n")
		b.WriteString(strings.TrimRight(module.Source, "\n"))
		b.WriteString("\n```\n")
	}

	return strings.TrimSpace(b.String()) + "\n", nil
}

func renderPythonDocItem(b *strings.Builder, module pythonModuleDoc, item pythonDocItem, level int, cfg PythonDocsConfig, moduleIndex, symbolIndex map[string]string, anchorName string, renderer *pythonSymbolRenderer) error {
	headingPrefix := strings.Repeat("#", level)
	anchor := pythonDocAnchor(anchorName)
	rendered, err := renderer.Render(module, item, level, headingPrefix, anchor, cfg, moduleIndex, symbolIndex)
	if err != nil {
		return err
	}
	b.WriteString(rendered)

	for i := range item.Methods {
		method := item.Methods[i]
		if err := renderPythonDocItem(b, module, method, level+1, cfg, moduleIndex, symbolIndex, item.Name+"-"+method.Name, renderer); err != nil {
			return err
		}
	}

	return nil
}

func renderPythonImport(imp pythonImport, moduleIndex, symbolIndex map[string]string) string {
	moduleRef := linkPythonReference(imp.Module, moduleIndex, symbolIndex)
	if imp.Kind == "import" {
		parts := make([]string, 0, len(imp.Names))
		for _, name := range imp.Names {
			part := moduleRef
			if name.AsName != "" {
				part += " as `" + name.AsName + "`"
			}
			parts = append(parts, part)
		}
		if len(parts) == 0 {
			return "import " + moduleRef
		}
		return "import " + strings.Join(parts, ", ")
	}

	imported := make([]string, 0, len(imp.Names))
	for _, name := range imp.Names {
		qualified := name.Name
		if imp.Module != "" {
			qualified = imp.Module + "." + name.Name
		}
		ref := linkPythonReference(qualified, moduleIndex, symbolIndex)
		if ref == "`"+qualified+"`" || ref == "`"+name.Name+"`" {
			ref = "`" + name.Name + "`"
		}
		if name.AsName != "" {
			ref += " as `" + name.AsName + "`"
		}
		imported = append(imported, ref)
	}
	if len(imported) == 0 {
		return "from " + moduleRef + " import *"
	}
	return "from " + moduleRef + " import " + strings.Join(imported, ", ")
}

var (
	pythonInlineRefRegex  = regexp.MustCompile("`([A-Za-z_][A-Za-z0-9_.]*)`")
	pythonSphinxRefRegex  = regexp.MustCompile(`:(?:mod|func|class|meth):` + "`" + `~?([^` + "`" + `]+)` + "`")
	pythonWhitespaceRegex = regexp.MustCompile(`\n{3,}`)
)

func linkPythonDocstring(text string, moduleIndex, symbolIndex map[string]string) string {
	text = pythonSphinxRefRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := pythonSphinxRefRegex.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		return linkPythonReference(parts[1], moduleIndex, symbolIndex)
	})

	text = pythonInlineRefRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := pythonInlineRefRegex.FindStringSubmatch(match)
		if len(parts) != 2 {
			return match
		}
		return linkPythonReference(parts[1], moduleIndex, symbolIndex)
	})

	return strings.TrimSpace(pythonWhitespaceRegex.ReplaceAllString(text, "\n\n"))
}

func linkPythonReference(ref string, moduleIndex, symbolIndex map[string]string) string {
	ref = strings.TrimSpace(strings.TrimPrefix(ref, "~"))
	if ref == "" {
		return "``"
	}
	if href, ok := symbolIndex[ref]; ok {
		return "[`" + ref + "`]" + "(" + href + ")"
	}
	if href, ok := moduleIndex[ref]; ok {
		return "[`" + ref + "`]" + "(" + href + ")"
	}
	if idx := strings.LastIndex(ref, "."); idx != -1 {
		if href, ok := symbolIndex[ref]; ok {
			return "[`" + ref + "`]" + "(" + href + ")"
		}
		if href, ok := moduleIndex[ref[:idx]]; ok {
			return "[`" + ref + "`]" + "(" + href + ")"
		}
	}
	return "`" + ref + "`"
}

var pythonReferenceTokenRegex = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_.]*`)

func extractPythonReferences(text string) []string {
	matches := pythonReferenceTokenRegex.FindAllString(text, -1)
	seen := make(map[string]struct{}, len(matches))
	refs := make([]string, 0, len(matches))
	for _, match := range matches {
		if !strings.Contains(match, ".") {
			continue
		}
		if _, ok := seen[match]; ok {
			continue
		}
		seen[match] = struct{}{}
		refs = append(refs, match)
	}
	return refs
}

func pythonDocAnchor(name string) string {
	return "symbol-" + models.Slugify(strings.ReplaceAll(name, ".", "-"))
}

func pythonKindLabel(kind string) string {
	switch kind {
	case htmlAttrClass:
		return "Class"
	case "function":
		return "Function"
	case "method":
		return "Method"
	default:
		return "Symbol"
	}
}

func firstParagraph(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	if idx := strings.Index(text, "\n\n"); idx != -1 {
		return strings.TrimSpace(text[:idx])
	}
	if idx := strings.Index(text, "\n"); idx != -1 {
		return strings.TrimSpace(text[:idx])
	}
	return text
}

func pythonDocsStringSlice(value interface{}) []string {
	switch v := value.(type) {
	case []string:
		return append([]string{}, v...)
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
				result = append(result, s)
			}
		}
		return result
	default:
		return nil
	}
}

func pythonDocsBoolEnv(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func pythonDocsExplicitlyEnabled(cfg *lifecycle.Config) bool {
	if cfg == nil {
		return false
	}
	if modelsConfig, ok := cfg.Extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		for _, hook := range modelsConfig.Hooks {
			if hook == pythonDocsPluginName {
				return true
			}
		}
		return false
	}
	if hooks := pythonDocsStringSlice(cfg.Extra["hooks"]); len(hooks) > 0 {
		for _, hook := range hooks {
			if hook == pythonDocsPluginName {
				return true
			}
		}
		return false
	}
	return false
}

func pythonDocsModelsConfig(cfg *lifecycle.Config) *models.Config {
	if cfg == nil || cfg.Extra == nil {
		return nil
	}
	modelsConfig, ok := cfg.Extra["models_config"].(*models.Config)
	if !ok {
		return nil
	}
	return modelsConfig
}

type pythonSymbolRenderer struct {
	config       PythonDocsConfig
	engine       *markatatemplates.Engine
	modelsConfig *models.Config
}

func newPythonSymbolRenderer(cfg *lifecycle.Config, pluginCfg PythonDocsConfig) (*pythonSymbolRenderer, error) {
	renderer := &pythonSymbolRenderer{config: pluginCfg, modelsConfig: pythonDocsModelsConfig(cfg)}
	if pluginCfg.SymbolTemplate == "" && pluginCfg.SymbolTemplateName == "" {
		return renderer, nil
	}

	templatesDir := ""
	themeName := templateTypeDefault
	if renderer.modelsConfig != nil {
		templatesDir = renderer.modelsConfig.TemplatesDir
		if renderer.modelsConfig.Theme.Name != "" {
			themeName = renderer.modelsConfig.Theme.Name
		}
	}

	engine, err := markatatemplates.NewEngineWithTheme(templatesDir, themeName)
	if err != nil {
		return nil, err
	}
	renderer.engine = engine
	return renderer, nil
}

func (r *pythonSymbolRenderer) Render(module pythonModuleDoc, item pythonDocItem, level int, headingPrefix, anchor string, cfg PythonDocsConfig, moduleIndex, symbolIndex map[string]string) (string, error) {
	if r == nil {
		return renderDefaultPythonSymbol(item, headingPrefix, anchor, cfg, moduleIndex, symbolIndex), nil
	}
	if r.engine != nil && (r.config.SymbolTemplate != "" || r.config.SymbolTemplateName != "") {
		custom, err := r.renderCustom(module, item, level, headingPrefix, anchor, cfg, moduleIndex, symbolIndex)
		if err != nil {
			return "", err
		}
		if custom != "" {
			return custom, nil
		}
	}
	return renderDefaultPythonSymbol(item, headingPrefix, anchor, cfg, moduleIndex, symbolIndex), nil
}

func (r *pythonSymbolRenderer) renderCustom(module pythonModuleDoc, item pythonDocItem, level int, headingPrefix, anchor string, cfg PythonDocsConfig, moduleIndex, symbolIndex map[string]string) (string, error) {
	if r.engine == nil {
		return "", nil
	}

	ctx := markatatemplates.NewContext(nil, "", r.modelsConfig)
	ctx.Set("module", pythonModuleToMap(module))
	ctx.Set("item", pythonDocItemToMap(item))
	ctx.Set("level", level)
	ctx.Set("heading_prefix", headingPrefix)
	ctx.Set("anchor", anchor)
	ctx.Set("docstring", linkPythonDocstring(item.Docstring, moduleIndex, symbolIndex))
	ctx.Set("extends", pythonLinksForBases(item.Bases, moduleIndex, symbolIndex))
	ctx.Set("references", pythonLinksForText(item.Signature+"\n"+item.Docstring, moduleIndex, symbolIndex))
	ctx.Set("source_enabled", cfg.IncludeSource)

	if r.config.SymbolTemplate != "" {
		return r.engine.RenderString(r.config.SymbolTemplate, ctx)
	}
	if r.config.SymbolTemplateName != "" {
		return r.engine.Render(r.config.SymbolTemplateName, ctx)
	}
	return "", nil
}

func renderDefaultPythonSymbol(item pythonDocItem, headingPrefix, anchor string, cfg PythonDocsConfig, moduleIndex, symbolIndex map[string]string) string {
	var b strings.Builder
	b.WriteString("<a id=\"" + anchor + "\"></a>\n")
	b.WriteString(headingPrefix + " " + pythonKindLabel(item.Kind) + " `" + item.Name + "`\n\n")

	if item.Signature != "" {
		b.WriteString("```python\n")
		b.WriteString(strings.TrimRight(item.Signature, "\n"))
		b.WriteString("\n```\n\n")
	}

	if len(item.Bases) > 0 {
		b.WriteString("Extends: ")
		b.WriteString(strings.Join(pythonLinksForBases(item.Bases, moduleIndex, symbolIndex), ", "))
		b.WriteString("\n\n")
	}

	related := pythonLinksForText(item.Signature+"\n"+item.Docstring, moduleIndex, symbolIndex)
	if len(related) > 0 {
		b.WriteString("Related: ")
		b.WriteString(strings.Join(related, ", "))
		b.WriteString("\n\n")
	}

	if item.Docstring != "" {
		b.WriteString(linkPythonDocstring(item.Docstring, moduleIndex, symbolIndex))
		b.WriteString("\n\n")
	}

	if cfg.IncludeSource && item.Source != "" {
		b.WriteString("<details>\n<summary>Source</summary>\n\n")
		b.WriteString("```python\n")
		b.WriteString(strings.TrimRight(item.Source, "\n"))
		b.WriteString("\n```\n\n</details>\n\n")
	}

	return b.String()
}

func pythonLinksForBases(bases []string, moduleIndex, symbolIndex map[string]string) []string {
	links := make([]string, 0, len(bases))
	for _, base := range bases {
		links = append(links, pythonLinksForText(base, moduleIndex, symbolIndex)...)
	}
	if len(links) == 0 {
		for _, base := range bases {
			links = append(links, linkPythonReference(base, moduleIndex, symbolIndex))
		}
	}
	return dedupeStrings(links)
}

func pythonLinksForText(text string, moduleIndex, symbolIndex map[string]string) []string {
	refs := extractPythonReferences(text)
	links := make([]string, 0, len(refs))
	for _, ref := range refs {
		candidates := []string{ref}
		if idx := strings.LastIndex(ref, "."); idx != -1 {
			candidates = append(candidates, ref[:idx])
		}
		for _, candidate := range candidates {
			linked := linkPythonReference(candidate, moduleIndex, symbolIndex)
			if linked != "`"+candidate+"`" {
				links = append(links, linked)
			}
		}
	}
	return dedupeStrings(links)
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func pythonModuleToMap(module pythonModuleDoc) map[string]interface{} {
	return map[string]interface{}{
		"source_path": module.SourcePath,
		"module_name": module.ModuleName,
		"docstring":   module.Docstring,
		"source":      module.Source,
	}
}

func pythonDocItemToMap(item pythonDocItem) map[string]interface{} {
	methods := make([]map[string]interface{}, 0, len(item.Methods))
	for i := range item.Methods {
		method := item.Methods[i]
		methods = append(methods, pythonDocItemToMap(method))
	}
	return map[string]interface{}{
		"name":      item.Name,
		"kind":      item.Kind,
		"signature": item.Signature,
		"docstring": item.Docstring,
		"source":    item.Source,
		"bases":     append([]string{}, item.Bases...),
		"methods":   methods,
	}
}

type pythonModuleDoc struct {
	SourcePath string          `json:"source_path"`
	ModuleName string          `json:"module_name"`
	Docstring  string          `json:"docstring"`
	Source     string          `json:"source"`
	Imports    []pythonImport  `json:"imports"`
	Classes    []pythonDocItem `json:"classes"`
	Functions  []pythonDocItem `json:"functions"`
}

type pythonImport struct {
	Kind   string             `json:"kind"`
	Module string             `json:"module"`
	Names  []pythonImportName `json:"names"`
}

type pythonImportName struct {
	Name   string `json:"name"`
	AsName string `json:"asname,omitempty"`
}

type pythonDocItem struct {
	Name      string          `json:"name"`
	Kind      string          `json:"kind"`
	Signature string          `json:"signature"`
	Docstring string          `json:"docstring"`
	Source    string          `json:"source"`
	Bases     []string        `json:"bases,omitempty"`
	Methods   []pythonDocItem `json:"methods,omitempty"`
}

const pythonDocsExtractorScript = `
import ast
import json
import os
import sys

base_dir = sys.argv[1]
files = json.load(sys.stdin)


def module_name_from_path(rel_path):
    rel_path = rel_path.replace("\\", "/")
    if rel_path.endswith(".py"):
        rel_path = rel_path[:-3]
    parts = [part for part in rel_path.split("/") if part]
    if parts and parts[-1] == "__init__":
        parts = parts[:-1]
    return ".".join(parts)


def docstring_expr(node):
    body = getattr(node, "body", None) or []
    if not body:
        return None
    first = body[0]
    if isinstance(first, ast.Expr) and isinstance(getattr(first, "value", None), ast.Constant) and isinstance(first.value.value, str):
        return first
    return None


def source_slice(lines, start, end):
    if not start or not end or end < start:
        return ""
    return "\n".join(lines[start - 1:end]).rstrip() + "\n"


def header_slice(node, lines):
    start = node.lineno
    if getattr(node, "decorator_list", None):
        start = min([decorator.lineno for decorator in node.decorator_list] + [node.lineno])
    body = getattr(node, "body", None) or []
    if body:
        end = body[0].lineno - 1
    else:
        end = getattr(node, "end_lineno", node.lineno)
    return source_slice(lines, start, end).rstrip()


def implementation_slice(node, lines):
    body = getattr(node, "body", None) or []
    if not body:
        return ""
    doc = docstring_expr(node)
    if doc is not None:
        start = getattr(doc, "end_lineno", doc.lineno) + 1
    else:
        start = body[0].lineno
    end = getattr(node, "end_lineno", start)
    if start > end:
        return ""
    return source_slice(lines, start, end)


def module_body_slice(tree, lines):
    body = getattr(tree, "body", None) or []
    if not body:
        return ""
    doc = docstring_expr(tree)
    if doc is not None:
        start = getattr(doc, "end_lineno", doc.lineno) + 1
    else:
        start = body[0].lineno
    if start > len(lines):
        return ""
    return "\n".join(lines[start - 1:]).rstrip() + "\n"


def is_private(name):
    return name.startswith("_") and not (name.startswith("__") and name.endswith("__"))


def resolve_relative_module(current_module, module, level):
    if not level:
        return module or ""
    parts = current_module.split(".") if current_module else []
    drop = max(level - 1, 0)
    base = parts[:-drop] if drop <= len(parts) and drop > 0 else parts[:]
    if module:
        base.append(module)
    return ".".join([part for part in base if part])


def import_entry(node, current_module):
    if isinstance(node, ast.Import):
        names = []
        for alias in node.names:
            names.append({"name": alias.name, "asname": alias.asname or ""})
        module = names[0]["name"] if len(names) == 1 else ""
        return {"kind": "import", "module": module, "names": names}
    if isinstance(node, ast.ImportFrom):
        names = []
        for alias in node.names:
            names.append({"name": alias.name, "asname": alias.asname or ""})
        return {
            "kind": "from",
            "module": resolve_relative_module(current_module, node.module or "", node.level),
            "names": names,
        }
    return None


def item_doc(node, kind, lines, include_private):
    if not include_private and is_private(node.name):
        return None
    doc = ast.get_docstring(node, clean=True) or ""
    item = {
        "name": node.name,
        "kind": kind,
        "signature": header_slice(node, lines),
        "docstring": doc,
        "source": implementation_slice(node, lines),
    }
    if isinstance(node, ast.ClassDef):
        bases = []
        for base in node.bases:
            try:
                bases.append(ast.unparse(base))
            except Exception:
                pass
        if bases:
            item["bases"] = bases
        methods = []
        for child in node.body:
            if isinstance(child, (ast.FunctionDef, ast.AsyncFunctionDef)):
                method = item_doc(child, "method", lines, include_private)
                if method is not None:
                    methods.append(method)
        if methods:
            item["methods"] = methods
    return item


documents = []
for rel_path in files:
    full_path = os.path.join(base_dir, rel_path)
    with open(full_path, "r", encoding="utf-8") as handle:
        source = handle.read()
    tree = ast.parse(source, filename=rel_path)
    lines = source.splitlines()
    module_name = module_name_from_path(rel_path)
    module_doc = ast.get_docstring(tree, clean=True) or ""
    include_private = bool(os.environ.get("MARKATA_GO_PYTHON_DOCS_INCLUDE_PRIVATE") == "1")

    imports = []
    classes = []
    functions = []
    for node in tree.body:
        if isinstance(node, (ast.Import, ast.ImportFrom)):
            entry = import_entry(node, module_name)
            if entry is not None:
                imports.append(entry)
        elif isinstance(node, ast.ClassDef):
            doc = item_doc(node, "class", lines, include_private)
            if doc is not None:
                classes.append(doc)
        elif isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)):
            doc = item_doc(node, "function", lines, include_private)
            if doc is not None:
                functions.append(doc)

    documents.append(
        {
            "source_path": rel_path.replace("\\", "/"),
            "module_name": module_name,
            "docstring": module_doc,
            "source": module_body_slice(tree, lines),
            "imports": imports,
            "classes": classes,
            "functions": functions,
        }
    )

json.dump(documents, sys.stdout)
`

// Ensure PythonDocsPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*PythonDocsPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*PythonDocsPlugin)(nil)
	_ lifecycle.LoadPlugin      = (*PythonDocsPlugin)(nil)
)
