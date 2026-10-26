// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

const (
	tailwindIncludeCSS      = "css"
	tailwindIncludeJS       = "js"
	tailwindDefaultInputCSS = "@tailwind base;\n@tailwind components;\n@tailwind utilities;\n"
)

type tailwindInstaller interface {
	Install() (string, error)
}

var (
	tailwindLookPath     = exec.LookPath
	newTailwindInstaller = func(config TailwindInstallerConfig) tailwindInstaller {
		return NewTailwindInstallerWithConfig(config)
	}
)

var tailwindLog = logging.Component("tailwind").Phase("cleanup")

// TailwindPlugin runs the Tailwind standalone CLI and wires inclusion into the head.
type TailwindPlugin struct {
	config    models.TailwindConfig
	assetURLs map[string]string
}

// NewTailwindPlugin creates a new TailwindPlugin with default settings.
func NewTailwindPlugin() *TailwindPlugin {
	return &TailwindPlugin{
		config: models.NewTailwindConfig(),
	}
}

// Name returns the unique name of the plugin.
func (p *TailwindPlugin) Name() string {
	return "tailwind"
}

// Priority ensures Tailwind runs early in Configure so output CSS exists before asset copying.
func (p *TailwindPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageConfigure {
		return lifecycle.PriorityEarly
	}
	if stage == lifecycle.StageCleanup {
		return lifecycle.PriorityEarly
	}
	return lifecycle.PriorityDefault
}

// Configure reads Tailwind config, runs the CLI, and injects CSS/JS includes.
func (p *TailwindPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config == nil {
		return nil
	}

	if config.Extra == nil {
		config.Extra = make(map[string]interface{})
	}

	if config.Extra["models_config"] == nil {
		config.Extra["models_config"] = p.buildModelsConfig(config)
	}

	p.assetURLs = assetURLsFromConfig(config)

	p.config = p.parseTailwindConfig(config.Extra)

	includeMode := normalizeTailwindInclude(p.config.IncludeMode())

	if includeMode != "" {
		if err := p.injectIncludes(config, includeMode); err != nil {
			return err
		}
	}

	if modelsConfig, ok := config.Extra["models_config"].(*models.Config); ok && modelsConfig != nil {
		config.Extra["theme"] = modelsConfig.Theme
		config.Extra["head"] = modelsConfig.Head
		templates.ClearConfigMapCache()
		if p.config.IsVerbose() {
			tailwindLog.Printf("theme in extra: %#v", config.Extra["theme"])
			tailwindLog.Printf("modelsConfig theme: %#v", modelsConfig.Theme)
		}
	}

	return nil
}

// Cleanup rebuilds Tailwind after HTML has been written so content globs that
// point at generated output (for example output/**/*.html) can produce the
// final utility set before minification and CSS purge run.
func (p *TailwindPlugin) Cleanup(m *lifecycle.Manager) error {
	config := m.Config()
	if config == nil {
		return nil
	}

	includeMode := normalizeTailwindInclude(p.config.IncludeMode())
	if includeMode != tailwindIncludeCSS || !p.config.IsBuildEnabled() {
		return nil
	}

	plan, err := p.prepareTailwindBuildPlan(m)
	if err != nil {
		return err
	}
	defer plan.cleanup()

	if plan.shouldBuild {
		if err := p.runTailwindBuild(config, plan.contentPaths); err != nil {
			return err
		}
	}

	if cache := GetBuildCache(m); cache != nil && plan.manifestHash != "" {
		cache.SetTailwindManifestHash(plan.manifestHash)
	}

	return p.syncBuiltCSSToOutput(config)
}

func (p *TailwindPlugin) buildModelsConfig(config *lifecycle.Config) *models.Config {
	if config == nil {
		return nil
	}
	return &models.Config{
		OutputDir: config.OutputDir,
		URL:       getStringFromExtra(config.Extra, "url"),
		Title:     getStringFromExtra(config.Extra, "title"),
		Author:    getStringFromExtra(config.Extra, "author"),
		AssetsDir: getStringFromExtra(config.Extra, "assets_dir"),
		Theme:     extractThemeConfig(config.Extra),
		Head:      extractHeadConfig(config.Extra),
	}
}

func extractThemeConfig(extra map[string]interface{}) models.ThemeConfig {
	if extra == nil {
		return models.ThemeConfig{}
	}
	if theme, ok := extra["theme"].(models.ThemeConfig); ok {
		return theme
	}
	if theme, ok := extra["theme"].(map[string]interface{}); ok {
		result := models.ThemeConfig{}
		if name, ok := theme["name"].(string); ok && name != "" {
			result.Name = name
		}
		if palette, ok := theme["palette"].(string); ok && palette != "" {
			result.Palette = palette
		}
		if paletteLight, ok := theme["palette_light"].(string); ok && paletteLight != "" {
			result.PaletteLight = paletteLight
		}
		if paletteDark, ok := theme["palette_dark"].(string); ok && paletteDark != "" {
			result.PaletteDark = paletteDark
		}
		if customCSS, ok := theme["custom_css"].(string); ok {
			result.CustomCSS = customCSS
		}
		if result.CustomCSS == "" {
			if customCSS, ok := theme["custom_css"].([]byte); ok {
				result.CustomCSS = string(customCSS)
			}
		}
		return result
	}
	return models.ThemeConfig{}
}

func extractHeadConfig(extra map[string]interface{}) models.HeadConfig {
	if extra == nil {
		return models.HeadConfig{}
	}
	if head, ok := extra["head"].(models.HeadConfig); ok {
		return head
	}
	return models.HeadConfig{}
}

func (p *TailwindPlugin) parseTailwindConfig(extra map[string]interface{}) models.TailwindConfig {
	result := models.NewTailwindConfig()
	if extra == nil {
		return result
	}

	if cfg, ok := extra["tailwind"].(models.TailwindConfig); ok {
		return cfg
	}

	raw, ok := extra["tailwind"].(map[string]interface{})
	if !ok {
		return result
	}

	p.parseTailwindConfigInclude(raw, &result)
	p.parseTailwindConfigStrings(raw, &result)
	p.parseTailwindConfigBools(raw, &result)

	if rawExtra := raw["extra_args"]; rawExtra != nil {
		result.ExtraArgs = tailwindConfigStringSlice(rawExtra)
	}

	return result
}

func (p *TailwindPlugin) parseTailwindConfigInclude(raw map[string]interface{}, result *models.TailwindConfig) {
	if include, ok := raw["include"]; ok {
		if v := tailwindConfigString(include); v != "" {
			result.Include = &v
		} else if b, ok := include.(bool); ok {
			if b {
				value := tailwindIncludeCSS
				result.Include = &value
			} else {
				value := BoolFalse
				result.Include = &value
			}
		}
	}
}

func (p *TailwindPlugin) parseTailwindConfigStrings(raw map[string]interface{}, result *models.TailwindConfig) {
	if v := tailwindConfigString(raw["input"]); v != "" {
		result.Input = v
	}
	if v := tailwindConfigString(raw["output"]); v != "" {
		result.Output = v
	}
	if v := tailwindConfigString(raw["config_file"]); v != "" {
		result.ConfigFile = v
	}
	if v := tailwindConfigString(raw["config"]); v != "" {
		result.ConfigFile = v
	}
	if v := tailwindConfigString(raw["version"]); v != "" {
		if normalized, err := tailwindNormalizeVersion(v); err == nil {
			result.Version = normalized
		}
	}
	if v := tailwindConfigString(raw["cache_dir"]); v != "" {
		result.CacheDir = v
	}
	if v := tailwindConfigString(raw["binary"]); v != "" {
		result.Binary = v
	}
}

func (p *TailwindPlugin) parseTailwindConfigBools(raw map[string]interface{}, result *models.TailwindConfig) {
	if rawBuild, ok := raw["build"]; ok {
		value := tailwindConfigBoolean(rawBuild, true)
		result.Build = &value
	}
	if rawPreflight, ok := raw["preflight"]; ok {
		value := tailwindConfigBoolean(rawPreflight, false)
		result.Preflight = &value
	}
	if rawMinify, ok := raw["minify"]; ok {
		value := tailwindConfigBoolean(rawMinify, true)
		result.Minify = &value
	}
	if rawAuto, ok := raw["auto_install"]; ok {
		value := tailwindConfigBoolean(rawAuto, true)
		result.AutoInstall = &value
	}
	if rawVerbose, ok := raw["verbose"]; ok {
		value := tailwindConfigBoolean(rawVerbose, false)
		result.Verbose = &value
	}
}

func (p *TailwindPlugin) runTailwindBuild(config *lifecycle.Config, contentPaths []string) error {
	inputPath, cleanupInput, err := p.resolveBuildInput(config)
	if err != nil {
		return err
	}
	defer cleanupInput()

	configPath, cleanupConfig, err := p.resolveBuildConfigFile(config, contentPaths)
	if err != nil {
		return err
	}
	defer cleanupConfig()

	outputPath := p.resolveAssetPath(config, p.config.Output)

	if inputPath == "" || outputPath == "" {
		return fmt.Errorf("tailwind: input and output paths must be set")
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return fmt.Errorf("tailwind: creating output directory: %w", err)
	}

	cliPath, err := p.findOrInstallTailwind()
	if err != nil {
		return err
	}
	if cliPath == "" {
		return nil
	}
	args := []string{"-i", inputPath, "-o", outputPath}
	if configPath != "" {
		args = append(args, "--config", configPath)
	}
	if p.config.IsMinifyEnabled() {
		args = append(args, "--minify")
	}
	if len(p.config.ExtraArgs) > 0 {
		args = append(args, p.config.ExtraArgs...)
	}

	cmd := exec.Command(cliPath, args...)
	cmd.Dir = "."
	cmd.Env = os.Environ()
	start := time.Now()
	tailwindLog.Printf("Running subprocess: %s %s", cliPath, strings.Join(args, " "))
	var stdoutLog, stderrLog *subprocessLineWriter
	var stderrTail *subprocessTailWriter
	cmd.Stdout = io.Discard
	if p.config.IsVerbose() {
		stdoutLog = newSubprocessLineWriter(tailwindLog)
		stderrLog = newSubprocessLineWriter(tailwindLog.Level("error"))
		cmd.Stdout, cmd.Stderr = stdoutLog, stderrLog
	} else {
		stderrTail = newSubprocessTailWriter(64 * 1024)
		cmd.Stderr = stderrTail
	}
	if err := cmd.Run(); err != nil {
		if stdoutLog != nil {
			stdoutLog.Flush()
			stderrLog.Flush()
		} else if output := stderrTail.String(); output != "" {
			logSubprocessOutput(tailwindLog.Level("error"), output)
		}
		return fmt.Errorf("tailwind build failed: %w", err)
	}
	if stdoutLog != nil {
		stdoutLog.Flush()
		stderrLog.Flush()
	}
	tailwindLog.Printf("Subprocess completed in %v", time.Since(start))

	return nil
}

func (p *TailwindPlugin) findOrInstallTailwind() (string, error) {
	if p.config.Binary != "" {
		if _, err := os.Stat(p.config.Binary); err == nil {
			return p.config.Binary, nil
		}
	}

	if path, err := tailwindLookPath(getTailwindBinaryName()); err == nil {
		return path, nil
	}

	if !p.config.IsAutoInstallEnabled() {
		tailwindLog.Warnf("tailwindcss not found in PATH, skipping build")
		tailwindLog.Printf("Install it or set [markata-go.tailwind].auto_install = true")
		return "", nil
	}

	version := p.config.Version
	if normalized, err := tailwindNormalizeVersion(version); err == nil {
		version = normalized
	}

	installer := newTailwindInstaller(TailwindInstallerConfig{
		Version:  version,
		CacheDir: p.config.CacheDir,
		Verbose:  p.config.Verbose,
	})

	if p.config.IsVerbose() {
		tailwindLog.Printf("Using managed Tailwind CLI %s", version)
	}

	installedPath, err := installer.Install()
	if err != nil {
		return "", fmt.Errorf("tailwind auto-install failed: %w", err)
	}

	return installedPath, nil
}

func (p *TailwindPlugin) injectIncludes(config *lifecycle.Config, includeMode string) error {
	modelsConfig, ok := config.Extra["models_config"].(*models.Config)
	if !ok || modelsConfig == nil {
		return nil
	}

	if includeMode == tailwindIncludeCSS {
		if strings.TrimSpace(modelsConfig.Theme.CustomCSS) == "" {
			includePath := p.includeAssetPath(config, p.config.Output)
			if includePath != "" {
				modelsConfig.Theme.CustomCSS = includePath
			}
		}
		return nil
	}

	if includeMode == tailwindIncludeJS {
		scriptSrc := resolveAssetURL(p.assetURLs, "tailwindcss-js", "https://cdn.tailwindcss.com")
		if !headHasScript(modelsConfig.Head.Script, scriptSrc) {
			modelsConfig.Head.Script = append(modelsConfig.Head.Script, models.ScriptTag{Src: scriptSrc})
		}
	}

	return nil
}

func (p *TailwindPlugin) resolveAssetPath(config *lifecycle.Config, path string) string {
	if path == "" {
		return ""
	}

	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}

	if filepath.IsAbs(trimmed) {
		return trimmed
	}

	assetsDir := ""
	if config.Extra != nil {
		if v, ok := config.Extra["assets_dir"].(string); ok && v != "" {
			assetsDir = v
		}
	}
	if assetsDir == "" {
		assetsDir = StaticDir
	}

	if strings.HasPrefix(trimmed, assetsDir+string(filepath.Separator)) || strings.HasPrefix(trimmed, assetsDir+"/") {
		return trimmed
	}

	return filepath.Join(assetsDir, trimmed)
}

func (p *TailwindPlugin) relativeAssetPath(config *lifecycle.Config, path string) string {
	if path == "" {
		return ""
	}

	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return ""
	}

	assetsDir := ""
	if config.Extra != nil {
		if v, ok := config.Extra["assets_dir"].(string); ok && v != "" {
			assetsDir = v
		}
	}
	if assetsDir == "" {
		assetsDir = StaticDir
	}
	absAssetsDir := assetsDir
	if !filepath.IsAbs(absAssetsDir) {
		if cwd, err := os.Getwd(); err == nil {
			absAssetsDir = filepath.Join(cwd, assetsDir)
		}
	}

	if filepath.IsAbs(trimmed) {
		rel, err := filepath.Rel(absAssetsDir, trimmed)
		if err != nil || strings.HasPrefix(rel, "..") {
			return trimmed
		}
		return filepath.ToSlash(rel)
	}

	if strings.HasPrefix(trimmed, assetsDir+string(filepath.Separator)) || strings.HasPrefix(trimmed, assetsDir+"/") {
		rel := strings.TrimPrefix(trimmed, assetsDir)
		rel = strings.TrimPrefix(rel, string(filepath.Separator))
		rel = strings.TrimPrefix(rel, "/")
		return filepath.ToSlash(rel)
	}

	return filepath.ToSlash(trimmed)
}

func (p *TailwindPlugin) includeAssetPath(config *lifecycle.Config, path string) string {
	relPath := p.relativeAssetPath(config, path)
	if relPath == "" {
		return ""
	}
	if isAbsoluteOrRootedPath(relPath) {
		return ""
	}
	return relPath
}

func (p *TailwindPlugin) syncBuiltCSSToOutput(config *lifecycle.Config) error {
	outputPath := p.resolveAssetPath(config, p.config.Output)
	if outputPath == "" {
		return nil
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("tailwind: reading built css: %w", err)
	}

	relPath := p.includeAssetPath(config, p.config.Output)
	if relPath == "" {
		return nil
	}

	destPath := filepath.Join(config.OutputDir, filepath.FromSlash(relPath))
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return fmt.Errorf("tailwind: creating output css directory: %w", err)
	}
	//nolint:gosec // Generated CSS must be readable by the web server.
	if err := os.WriteFile(destPath, data, 0o644); err != nil {
		return fmt.Errorf("tailwind: syncing built css to output: %w", err)
	}
	return p.writeHashedTailwindCopy(destPath, data)
}

func (p *TailwindPlugin) writeHashedTailwindCopy(destPath string, data []byte) error {
	if len(data) == 0 {
		return nil
	}
	ext := filepath.Ext(destPath)
	if ext == "" {
		return nil
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))[:8]
	hashedPath := strings.TrimSuffix(destPath, ext) + "." + hash + ext
	if hashedPath == destPath {
		return nil
	}
	//nolint:gosec // Generated CSS must be readable by the web server.
	if err := os.WriteFile(hashedPath, data, 0o644); err != nil {
		return fmt.Errorf("tailwind: writing hashed css copy: %w", err)
	}
	relPath := filepath.ToSlash(filepath.Base(destPath))
	if dir := filepath.ToSlash(filepath.Dir(destPath)); dir != "." {
		relPath = filepath.ToSlash(filepath.Join(filepath.Base(dir), filepath.Base(destPath)))
	}
	templates.SetAssetHashes(map[string]string{relPath: hash})
	return nil
}

func (p *TailwindPlugin) resolveBuildInput(config *lifecycle.Config) (inputPath string, cleanup func(), err error) {
	inputPath = p.resolveAssetPath(config, p.config.Input)
	if inputPath != "" {
		if _, err := os.Stat(inputPath); err == nil {
			return inputPath, func() {}, nil
		} else if err != nil && !os.IsNotExist(err) {
			return "", nil, fmt.Errorf("tailwind: checking input file: %w", err)
		}
	}

	tmpFile, err := os.CreateTemp("", "markata-tailwind-*.css")
	if err != nil {
		return "", nil, fmt.Errorf("tailwind: creating default input css: %w", err)
	}

	cleanup = func() {
		_ = os.Remove(tmpFile.Name())
	}

	if _, err := tmpFile.WriteString(tailwindDefaultInputCSS); err != nil {
		cleanup()
		_ = tmpFile.Close()
		return "", nil, fmt.Errorf("tailwind: writing default input css: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("tailwind: closing default input css: %w", err)
	}

	if p.config.IsVerbose() {
		tailwindLog.Printf("Input CSS missing, using generated default input")
	}

	return tmpFile.Name(), cleanup, nil
}

func (p *TailwindPlugin) resolveBuildConfigFile(_ *lifecycle.Config, contentPaths []string) (configPath string, cleanup func(), err error) {
	if p.config.ConfigFile != "" || len(p.config.ExtraArgs) > 0 {
		return p.config.ConfigFile, func() {}, nil
	}

	if len(contentPaths) == 0 {
		return "", func() {}, nil
	}

	tmpFile, err := os.CreateTemp("", "markata-tailwind-*.config.js")
	if err != nil {
		return "", nil, fmt.Errorf("tailwind: creating default config file: %w", err)
	}

	cleanup = func() {
		_ = os.Remove(tmpFile.Name())
	}

	content := make([]string, 0, len(contentPaths))
	for _, pattern := range contentPaths {
		content = append(content, fmt.Sprintf("  %q", filepath.ToSlash(pattern)))
	}

	configJS := "module.exports = {\ncontent: [\n" + strings.Join(content, ",\n") + "\n],\ncorePlugins: {\n  preflight: " + boolToJS(p.config.IsPreflightEnabled()) + "\n}\n}\n"
	if _, err := tmpFile.WriteString(configJS); err != nil {
		cleanup()
		_ = tmpFile.Close()
		return "", nil, fmt.Errorf("tailwind: writing default config file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("tailwind: closing default config file: %w", err)
	}

	if p.config.IsVerbose() {
		tailwindLog.Printf("Generated default config with %d content patterns", len(contentPaths))
	}

	return tmpFile.Name(), cleanup, nil
}

func (p *TailwindPlugin) usesGeneratedContentConfig() bool {
	return p.config.ConfigFile == "" && len(p.config.ExtraArgs) == 0
}

func (p *TailwindPlugin) tailwindOutputExists(config *lifecycle.Config) bool {
	outputPath := p.resolveAssetPath(config, p.config.Output)
	if outputPath == "" {
		return false
	}
	_, err := os.Stat(outputPath)
	return err == nil
}

func (p *TailwindPlugin) computeTailwindManifestHash(config *lifecycle.Config, manifest string) string {
	var builder strings.Builder
	builder.WriteString(manifest)
	builder.WriteString("\n--tailwind-version:")
	builder.WriteString(strings.TrimSpace(p.config.Version))
	builder.WriteString("\n--tailwind-include:")
	builder.WriteString(strings.TrimSpace(p.config.IncludeMode()))
	builder.WriteString("\n--tailwind-output:")
	builder.WriteString(strings.TrimSpace(p.config.Output))
	builder.WriteString("\n--tailwind-input-hash:")
	builder.WriteString(p.tailwindInputHash(config))
	builder.WriteString("\n--tailwind-minify:")
	if p.config.IsMinifyEnabled() {
		builder.WriteString("true")
	} else {
		builder.WriteString("false")
	}
	builder.WriteString("\n--tailwind-preflight:")
	if p.config.IsPreflightEnabled() {
		builder.WriteString("true")
	} else {
		builder.WriteString("false")
	}
	if p.config.ConfigFile != "" {
		builder.WriteString("\n--tailwind-config-hash:")
		builder.WriteString(hashFileContents(p.config.ConfigFile))
	}
	if len(p.config.ExtraArgs) > 0 {
		builder.WriteString("\n--tailwind-extra-args:")
		builder.WriteString(strings.Join(p.config.ExtraArgs, "\x00"))
	}
	return buildcache.ContentHash(builder.String())
}

func (p *TailwindPlugin) tailwindInputHash(config *lifecycle.Config) string {
	inputPath := p.resolveAssetPath(config, p.config.Input)
	if inputPath == "" {
		return buildcache.ContentHash(tailwindDefaultInputCSS)
	}
	if data, err := os.ReadFile(inputPath); err == nil {
		return buildcache.ContentHash(string(data))
	}
	return buildcache.ContentHash(tailwindDefaultInputCSS)
}

func normalizeTailwindInclude(value string) string {
	return tailwindNormalizeInclude(value)
}

func isAbsoluteOrRootedPath(path string) bool {
	if filepath.IsAbs(path) {
		return true
	}
	if len(path) >= 3 && (path[1] == ':' && (path[2] == '\\' || path[2] == '/')) {
		return true
	}
	if path != "" && (path[0] == '/' || path[0] == '\\') {
		return true
	}
	return false
}

func headHasScript(scripts []models.ScriptTag, src string) bool {
	if src == "" {
		return false
	}
	for _, script := range scripts {
		if strings.TrimSpace(script.Src) == src {
			return true
		}
	}
	return false
}

func boolToJS(value bool) string {
	if value {
		return BoolTrue
	}
	return BoolFalse
}

func tailwindNormalizeInclude(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "", "false", "off", "0", "none":
		return ""
	case "css", "js":
		return value
	default:
		return ""
	}
}

// Ensure TailwindPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*TailwindPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*TailwindPlugin)(nil)
	_ lifecycle.CleanupPlugin   = (*TailwindPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*TailwindPlugin)(nil)
)

// No additional helpers.
