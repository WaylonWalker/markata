package plugins

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"
	"github.com/WaylonWalker/markata/pkg/templates"
)

var paletteCSSLog = logging.Component("palette_css")

const (
	themeModeDark  = "dark"
	themeModeLight = "light"
)

// PaletteCSSPlugin generates CSS variables from the configured color palette.
// It runs during the Write stage and creates/overwrites css/palette.css
// with the palette's CSS custom properties. It runs after static_assets
// to overwrite the default palette.css with palette-specific values.
//
// The plugin supports intelligent light/dark palette mapping:
// - palette = "everforest" will auto-detect everforest-light and everforest-dark
// - palette_light and palette_dark can override the auto-detected variants
//
// The plugin maps palette colors to the theme's expected CSS variable names,
// preserving fonts, spacing, and other non-color variables from the default theme.
type PaletteCSSPlugin struct{}

// NewPaletteCSSPlugin creates a new PaletteCSSPlugin.
func NewPaletteCSSPlugin() *PaletteCSSPlugin {
	return &PaletteCSSPlugin{}
}

// Name returns the unique name of the plugin.
func (p *PaletteCSSPlugin) Name() string {
	return "palette_css"
}

// Configure generates CSS from the configured palette and registers its hash
// so templates can use the correct hashed filename. This runs in Configure stage
// before templates are rendered in Render stage.
func (p *PaletteCSSPlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	delete(config.Extra, "palette_core_css")
	delete(config.Extra, "palette_core_enabled")

	// Get palette configuration from config.Extra["theme"].
	paletteName, paletteLight, paletteDark, seedColor := p.getPaletteConfig(config.Extra)
	paletteName = canonicalPaletteName(paletteName)
	paletteLight = canonicalPaletteName(paletteLight)
	paletteDark = canonicalPaletteName(paletteDark)
	fallbackMode := p.getThemeFallbackMode(config.Extra)
	userVariables := p.getThemeVariables(config.Extra)
	if paletteName == "" {
		return nil
	}

	// Load palettes
	loader := palettes.NewLoader()

	if paletteName == "generated" {
		if seedColor == "" {
			paletteCSSLog.Phase("configure").Warnf("palette is 'generated' but no seed_color provided. Using fallback.")
		} else {
			lightP, err := palettes.GenerateTriadicPalette(seedColor, palettes.VariantLight)
			if err == nil {
				loader.AddPalette("generated-light", lightP)
			} else {
				paletteCSSLog.Phase("configure").Errorf("generating light palette from seed %q: %v", seedColor, err)
			}
			darkP, err := palettes.GenerateTriadicPalette(seedColor, palettes.VariantDark)
			if err == nil {
				loader.AddPalette("generated-dark", darkP)
			} else {
				paletteCSSLog.Phase("configure").Errorf("generating dark palette from seed %q: %v", seedColor, err)
			}
		}
	}

	// Publish the effective light/dark palette names so templates can pin
	// data-palette to the palette matching the resolved color mode before
	// palette.css is available to scripts.
	effectiveLight, effectiveDark := palettes.GetEffectivePalettes(paletteName, paletteLight, paletteDark)
	config.Extra["palette_light_effective"] = normalizePaletteName(effectiveLight)
	config.Extra["palette_dark_effective"] = normalizePaletteName(effectiveDark)

	// Check if theme switcher is enabled
	switcherEnabled := p.isSwitcherEnabled(config.Extra)

	var css string
	if switcherEnabled {
		css = p.generateMultiPaletteCSS(loader, config.Extra, paletteName, paletteLight, paletteDark, userVariables, fallbackMode)
	} else {
		css = p.generateSinglePaletteCSS(loader, paletteName, paletteLight, paletteDark, userVariables, fallbackMode)
	}
	css += p.generateContractPaletteCSS(config)
	if core, ok := config.Extra["palette_core_css"].(string); ok && switcherEnabled {
		core += p.generateContractPaletteCSS(config, normalizePaletteName(effectiveLight), normalizePaletteName(effectiveDark))
		config.Extra["palette_core_css"] = core
		coreHash := fmt.Sprintf("%x", sha256.Sum256([]byte(core)))[:8]
		m.SetAssetHash("css/palette-core.css", coreHash)
		templates.SetAssetHashes(map[string]string{"css/palette-core.css": coreHash})
		config.Extra["palette_core_enabled"] = true
	}

	// Compute hash of generated CSS
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(css)))[:8]

	// Register hash so templates can use it
	m.SetAssetHash("css/palette.css", hash)
	templates.SetAssetHashes(map[string]string{"css/palette.css": hash})

	paletteCSSLog.Phase("configure").Printf("Registered hash %s for palette.css", hash)

	return nil
}

// Write generates CSS from the configured palette and writes it to the output directory.
func (p *PaletteCSSPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	outputDir := config.OutputDir
	if hashedAssetReadyInFastMode(m, "css/palette.css") &&
		(config.Extra["palette_core_enabled"] != true || hashedAssetReadyInFastMode(m, "css/palette-core.css")) {
		return nil
	}

	// Get palette configuration from config.Extra["theme"].
	paletteName, paletteLight, paletteDark, seedColor := p.getPaletteConfig(config.Extra)
	paletteName = canonicalPaletteName(paletteName)
	paletteLight = canonicalPaletteName(paletteLight)
	paletteDark = canonicalPaletteName(paletteDark)
	fallbackMode := p.getThemeFallbackMode(config.Extra)
	userVariables := p.getThemeVariables(config.Extra)
	if paletteName == "" {
		// No palette configured, skip
		paletteCSSLog.Phase("write").Printf("No palette configured, skipping CSS generation")
		return nil
	}

	paletteCSSLog.Phase("write").Printf("Generating CSS for palette: %s (light: %s, dark: %s)", paletteName, paletteLight, paletteDark)

	// Publish the effective light/dark palette names so templates can pin
	// data-palette to the palette matching the resolved color mode before
	// palette.css is available to scripts.
	effectiveLight, effectiveDark := palettes.GetEffectivePalettes(paletteName, paletteLight, paletteDark)
	config.Extra["palette_light_effective"] = normalizePaletteName(effectiveLight)
	config.Extra["palette_dark_effective"] = normalizePaletteName(effectiveDark)

	// Check if theme switcher is enabled
	switcherEnabled := p.isSwitcherEnabled(config.Extra)

	// Load palettes
	loader := palettes.NewLoader()

	if paletteName == "generated" {
		if seedColor == "" {
			paletteCSSLog.Phase("write").Warnf("palette is 'generated' but no seed_color provided. Using fallback.")
		} else {
			lightP, err := palettes.GenerateTriadicPalette(seedColor, palettes.VariantLight)
			if err == nil {
				loader.AddPalette("generated-light", lightP)
			} else {
				paletteCSSLog.Phase("write").Errorf("generating light palette from seed %q: %v", seedColor, err)
			}
			darkP, err := palettes.GenerateTriadicPalette(seedColor, palettes.VariantDark)
			if err == nil {
				loader.AddPalette("generated-dark", darkP)
			} else {
				paletteCSSLog.Phase("write").Errorf("generating dark palette from seed %q: %v", seedColor, err)
			}
		}
	}

	var css string
	if switcherEnabled {
		// Generate CSS for all palettes when switcher is enabled
		css = p.generateMultiPaletteCSS(loader, config.Extra, paletteName, paletteLight, paletteDark, userVariables, fallbackMode)
	} else {
		// Generate CSS for just the configured light/dark pair
		css = p.generateSinglePaletteCSS(loader, paletteName, paletteLight, paletteDark, userVariables, fallbackMode)
	}
	css += p.generateContractPaletteCSS(config)
	if core, ok := config.Extra["palette_core_css"].(string); ok && switcherEnabled {
		core += p.generateContractPaletteCSS(config, normalizePaletteName(effectiveLight), normalizePaletteName(effectiveDark))
		if err := writePaletteCoreCSS(outputDir, core, m.GetAssetHash("css/palette-core.css")); err != nil {
			return err
		}
	}

	// Write to output directory
	cssDir := filepath.Join(outputDir, "css")
	cssPath := filepath.Join(cssDir, "palette.css")
	if existing, err := os.ReadFile(cssPath); err == nil {
		if bytes.Equal(existing, []byte(css)) {
			paletteCSSLog.Phase("write").Printf("CSS unchanged, skipping write")
			return nil
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading existing palette CSS: %w", err)
	}

	if err := os.MkdirAll(cssDir, 0o755); err != nil {
		return fmt.Errorf("creating css directory: %w", err)
	}
	//nolint:gosec // G306: palette.css is a public CSS file, 0644 is appropriate
	if err := os.WriteFile(cssPath, []byte(css), 0o644); err != nil {
		return fmt.Errorf("writing palette CSS: %w", err)
	}

	if hash := m.GetAssetHash("css/palette.css"); hash != "" {
		base := strings.TrimSuffix(filepath.Base(cssPath), filepath.Ext(cssPath))
		hashedPath := filepath.Join(cssDir, fmt.Sprintf("%s.%s.css", base, hash))
		//nolint:gosec // G306: hashed palette.css is a public CSS file, 0644 is appropriate
		if err := os.WriteFile(hashedPath, []byte(css), 0o644); err != nil {
			return fmt.Errorf("writing hashed palette CSS: %w", err)
		}
	}

	paletteCSSLog.Phase("write").Printf("Wrote %d bytes to %s", len(css), cssPath)

	return nil
}

// canonicalPaletteName resolves contract aliases before the legacy palette
// loader is called. The loader remains a compatibility projection, but it must
// receive the same identity selected by the canonical contract.
func canonicalPaletteName(name string) string {
	if name == "" {
		return name
	}
	contract, err := renderingcontract.Load()
	if err != nil {
		return name
	}
	seen := map[string]bool{}
	for {
		if seen[name] {
			return name
		}
		seen[name] = true
		alias, ok := contract.Aliases[name]
		if !ok || strings.HasPrefix(alias, "theme.") {
			return name
		}
		name = alias
	}
}

//nolint:gocyclo // Palette projection has separate fallback and variant branches.
func (p *PaletteCSSPlugin) generateContractPaletteCSS(config *lifecycle.Config, paletteIDs ...string) string {
	modelsConfig, ok := config.Extra["models_config"].(*models.Config)
	if !ok {
		return ""
	}
	contract, err := renderingcontract.Load()
	if err != nil {
		return ""
	}
	resolve := func(id string) string {
		for i := 0; i < len(contract.Aliases); i++ {
			alias, ok := contract.Aliases[id]
			if !ok || strings.HasPrefix(alias, "theme.") {
				return id
			}
			id = alias
		}
		return id
	}
	selected := resolve(modelsConfig.Theme.Palette)
	find := func(id string) *renderingcontract.Palette {
		for index := range contract.Palettes {
			if contract.Palettes[index].ID == id {
				return &contract.Palettes[index]
			}
		}
		return nil
	}
	active := find(selected)
	if active == nil {
		return ""
	}
	var light, dark *renderingcontract.Palette
	for index := range contract.Palettes {
		palette := &contract.Palettes[index]
		if palette.Family != active.Family {
			continue
		}
		if palette.Variant == "light" {
			light = palette
		}
		if palette.Variant == "dark" {
			dark = palette
		}
	}
	write := func(builder *strings.Builder, palette *renderingcontract.Palette) {
		if palette == nil {
			return
		}
		projected := renderingcontract.FinalRenderPalette(*palette)
		background, surface := projected["background"], projected["surface"]
		ink, accent := projected["text"], projected["link"]
		roles := map[string]string{"background": background, "surface": surface, "ink": ink, "accent": accent}
		for _, role := range []string{"background", "surface", "ink", "accent"} {
			fmt.Fprintf(builder, "  --%s: %s;\n", map[string]string{"background": "bg", "surface": "panel", "ink": "ink", "accent": "accent"}[role], roles[role])
		}
		// These are the variables consumed by the active Markata templates.
		// The contract's 3:1 "border" role (ink) is kept as --color-border-strong
		// for control boundaries; decorative borders use a low-contrast mix so
		// cards, dividers, and surfaces read as soft edges instead of hard lines.
		fmt.Fprintf(builder, "  --color-background: %s;\n  --color-surface: %s;\n  --color-text: %s;\n  --color-primary: %s;\n", background, surface, ink, accent)
		fmt.Fprintf(builder, "  --color-border-strong: %s;\n  --color-border: color-mix(in srgb, %s 16%%, %s);\n", ink, ink, background)
	}
	var builder strings.Builder
	builder.WriteString("\n@layer tokens {\n")
	aesthetic := modelsConfig.Theme.Aesthetic
	if aesthetic == "" {
		aesthetic = contract.Defaults.Aesthetic
	}
	if values, ok := contract.Aesthetics[aesthetic]; ok {
		fmt.Fprintf(&builder, ":root { --aesthetic: %q;\n", aesthetic)
		keys := make([]string, 0, len(values))
		for key := range values {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			fmt.Fprintf(&builder, "  --theme-%s: %v;\n", key, values[key])
		}
		builder.WriteString("}\n")
	}
	if dark != nil {
		// Mode fallbacks use zero specificity so any [data-palette] block,
		// including visitor theme picks outside the contract, overrides them.
		builder.WriteString(":where(:root:not([data-theme=\"light\"]), [data-theme=\"dark\"]) {\n")
		write(&builder, dark)
		builder.WriteString("}\n")
	}
	if light != nil {
		builder.WriteString(":where([data-theme=\"light\"]) {\n")
		write(&builder, light)
		builder.WriteString("}\n")
	}
	for _, palette := range contract.Palettes {
		if len(paletteIDs) > 0 {
			include := false
			for _, id := range paletteIDs {
				include = include || palette.ID == resolve(id)
			}
			if !include {
				continue
			}
		}
		fmt.Fprintf(&builder, "[data-palette=%q], :root[data-palette=%q] {\n", palette.ID, palette.ID)
		write(&builder, &palette)
		builder.WriteString("}\n")
	}
	builder.WriteString("}\n")
	return builder.String()
}

// isSwitcherEnabled checks if the theme switcher is enabled in config.
func (p *PaletteCSSPlugin) isSwitcherEnabled(extra map[string]interface{}) bool {
	return themeSwitcherEnabled(extra)
}

// themeSwitcherEnabled resolves theme.switcher.enabled from plugin Extra.
// The theme picker is on unless a config explicitly turns it off.
func themeSwitcherEnabled(extra map[string]interface{}) bool {
	defaults := models.NewThemeSwitcherConfig()
	if extra == nil {
		return defaults.IsEnabled()
	}
	if themeConfig, ok := extra["theme"].(models.ThemeConfig); ok {
		return themeConfig.Switcher.IsEnabled()
	}
	// Check if it's a map (raw TOML)
	if theme, ok := extra["theme"].(map[string]interface{}); ok {
		if switcher, ok := theme["switcher"].(map[string]interface{}); ok {
			if enabled, ok := switcher["enabled"].(bool); ok {
				return enabled
			}
		}
	}
	if configured, ok := extra["models_config"].(*models.Config); ok && configured != nil {
		return configured.Theme.Switcher.IsEnabled()
	}
	return defaults.IsEnabled()
}

// getSwitcherConfig gets the theme switcher configuration from Extra.
func (p *PaletteCSSPlugin) getSwitcherConfig(extra map[string]interface{}) models.ThemeSwitcherConfig {
	if extra == nil {
		return models.NewThemeSwitcherConfig()
	}
	if themeConfig, ok := extra["theme"].(models.ThemeConfig); ok {
		return themeConfig.Switcher
	}
	// Return default if not found
	return models.NewThemeSwitcherConfig()
}

// generateSinglePaletteCSS generates CSS for a single light/dark palette pair.
func (p *PaletteCSSPlugin) generateSinglePaletteCSS(loader *palettes.Loader, paletteName, paletteLight, paletteDark string, overrides map[string]string, fallbackMode string) string {
	// Get effective light and dark palette names
	lightName, darkName := palettes.GetEffectivePalettes(paletteName, paletteLight, paletteDark)

	var lightPalette, darkPalette *palettes.Palette
	var err error

	// Load light palette
	if lightName != "" {
		lightPalette, err = loader.Load(lightName)
		if err != nil {
			// Fall back to base palette, ignore error as we'll handle nil palette later
			//nolint:errcheck // fallback load failure is handled by nil check below
			lightPalette, _ = loader.Load(paletteName)
		}
	}

	// Load dark palette
	if darkName != "" {
		darkPalette, err = loader.Load(darkName)
		if err != nil {
			// Fall back to base palette, ignore error as we'll handle nil palette later
			//nolint:errcheck // fallback load failure is handled by nil check below
			darkPalette, _ = loader.Load(paletteName)
		}
	}

	// If we have no palettes loaded, try to load the base palette
	if lightPalette == nil && darkPalette == nil {
		palette, err := loader.Load(paletteName)
		if err == nil {
			lightPalette = palette
			darkPalette = palette
		}
	}

	// Generate theme-compatible CSS with both variants
	return p.generateThemeCSSWithVariants(lightPalette, darkPalette, lightName, darkName, overrides, fallbackMode)
}

// PaletteManifestEntry represents one pickable theme in the switcher manifest.
//
// Each entry is a concrete palette for a single color mode. Counterpart names
// the palette the picker switches to when the visitor flips dark/light mode.
type PaletteManifestEntry struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Variant     string `json:"variant"`
	BaseName    string `json:"baseName"`
	Counterpart string `json:"counterpart,omitempty"`
	Derived     bool   `json:"derived,omitempty"`

	// loadName is the name passed to the palette loader (not serialized).
	loadName string
}

// generateMultiPaletteCSS generates CSS for all available palettes when switcher is enabled.
func (p *PaletteCSSPlugin) generateMultiPaletteCSS(loader *palettes.Loader, extra map[string]interface{}, paletteName, paletteLight, paletteDark string, overrides map[string]string, fallbackMode string) string {
	var buf bytes.Buffer

	buf.WriteString("@layer reset, tokens, base, components, utilities, overrides;\n\n")
	buf.WriteString("@layer tokens {\n")

	// Get all available palettes
	allPalettes, err := loader.Discover()
	if err != nil {
		// Fall back to single palette CSS on error
		return p.generateSinglePaletteCSS(loader, paletteName, paletteLight, paletteDark, overrides, fallbackMode)
	}

	// Filter palettes based on switcher config
	switcherConfig := p.getSwitcherConfig(extra)
	filteredPalettes := p.filterPalettes(allPalettes, switcherConfig)

	// Get effective light and dark palette names for the default
	lightName, darkName := palettes.GetEffectivePalettes(paletteName, paletteLight, paletteDark)

	// Generate one CSS block per pickable palette. Blocks use a bare
	// [data-palette] attribute selector so the same rules can scope any
	// element, which lets the theme picker render true-color preview cards.
	candidates := p.generatePaletteManifest(filteredPalettes, normalizePaletteName(lightName), normalizePaletteName(darkName))
	var blocks bytes.Buffer
	var coreBlocks bytes.Buffer
	manifest := make([]PaletteManifestEntry, 0, len(candidates))
	coreManifest := make([]PaletteManifestEntry, 0, 2)
	for _, entry := range candidates {
		palette, err := loader.Load(entry.loadName)
		if err != nil && entry.loadName != entry.Name {
			palette, err = loader.Load(entry.Name)
		}
		if err != nil {
			continue
		}
		blocks.WriteString(fmt.Sprintf("/* Palette: %s (%s) */\n", entry.Name, entry.Variant))
		fmt.Fprintf(&blocks, "[data-palette=%q], :root[data-palette=%q] {\n", entry.Name, entry.Name)
		p.writePaletteVariablesIndented(&blocks, palette, "  ")
		blocks.WriteString("}\n\n")
		if entry.Name == normalizePaletteName(lightName) || entry.Name == normalizePaletteName(darkName) {
			coreManifest = append(coreManifest, entry)
			fmt.Fprintf(&coreBlocks, "[data-palette=%q], :root[data-palette=%q] {\n", entry.Name, entry.Name)
			p.writePaletteVariablesIndented(&coreBlocks, palette, "  ")
			coreBlocks.WriteString("}\n\n")
		}
		manifest = append(manifest, entry)
	}
	manifest = linkCounterparts(manifest)
	coreManifest = linkCounterparts(coreManifest)
	if len(coreManifest) == 0 {
		coreManifest = manifest
	}
	if extra != nil {
		extra["palette_seasonal_json"] = seasonalScheduleJSON(manifest, time.Now())
	}

	// Header comment
	buf.WriteString("/* CSS Custom Properties - Generated by markata-go */\n")
	buf.WriteString("/* Multi-palette theme switcher enabled */\n")
	buf.WriteString(fmt.Sprintf("/* Default: Light=%s, Dark=%s */\n\n", lightName, darkName))

	// Generate palette manifest for JavaScript
	escapedManifest := paletteManifestCSS(manifest)
	coreManifestCSS := paletteManifestCSS(coreManifest)

	// Write non-color variables and palette manifest in :root
	buf.WriteString("/* Global configuration and non-color variables */\n")
	buf.WriteString(":root {\n")
	buf.WriteString(fmt.Sprintf("  --palette-light: %q;\n", lightName))
	buf.WriteString(fmt.Sprintf("  --palette-dark: %q;\n", darkName))
	buf.WriteString(fmt.Sprintf("  --palette-manifest: '%s';\n", escapedManifest))
	buf.WriteString("  --palette-switcher-enabled: 1;\n")

	buf.WriteString("}\n\n")
	buf.Write(blocks.Bytes())

	if normalizeThemeFallbackMode(fallbackMode) == themeModeLight {
		// Light fallback (default)
		if lightName != "" {
			lightPalette, err := loader.Load(lightName)
			if err == nil {
				buf.WriteString(fmt.Sprintf("/* Default light mode - %s */\n", lightName))
				buf.WriteString(":root:not([data-theme=\"dark\"]):not([data-palette]),\n")
				buf.WriteString("[data-theme=\"light\"]:not([data-palette]) {\n")
				p.writePaletteVariables(&buf, lightPalette)
				buf.WriteString("}\n\n")
			}
		}

		if darkName != "" {
			darkPalette, err := loader.Load(darkName)
			if err == nil {
				buf.WriteString(fmt.Sprintf("/* Dark mode - %s */\n", darkName))
				buf.WriteString("[data-theme=\"dark\"]:not([data-palette]) {\n")
				p.writePaletteVariables(&buf, darkPalette)
				buf.WriteString("}\n\n")
			}
		}
	} else {
		// Dark fallback (default)
		if darkName != "" {
			darkPalette, err := loader.Load(darkName)
			if err == nil {
				buf.WriteString(fmt.Sprintf("/* Default dark mode - %s */\n", darkName))
				buf.WriteString(":root:not([data-theme=\"light\"]):not([data-palette]),\n")
				buf.WriteString("[data-theme=\"dark\"]:not([data-palette]) {\n")
				p.writePaletteVariables(&buf, darkPalette)
				buf.WriteString("}\n\n")
			}
		}

		if lightName != "" {
			lightPalette, err := loader.Load(lightName)
			if err == nil {
				buf.WriteString(fmt.Sprintf("/* Light mode - %s */\n", lightName))
				buf.WriteString("[data-theme=\"light\"]:not([data-palette]) {\n")
				p.writePaletteVariables(&buf, lightPalette)
				buf.WriteString("}\n\n")
			}
		}
	}

	buf.WriteString("}\n\n")
	p.writeThemeOverrides(&buf, overrides)
	setPaletteCoreCSS(extra, buf.String(), blocks.String(), coreBlocks.String(), escapedManifest, coreManifestCSS)

	return buf.String()
}

// paletteManifestCSS escapes catalog metadata for a single-quoted CSS value.
func paletteManifestCSS(manifest []PaletteManifestEntry) string {
	encoded, err := json.Marshal(manifest)
	if err != nil {
		return "[]"
	}
	return strings.ReplaceAll(string(encoded), "'", "\\'")
}

// setPaletteCoreCSS substitutes the configured pair for the full pickable set.
func setPaletteCoreCSS(extra map[string]interface{}, fullCSS, paletteBlocks, coreBlocks, manifest, coreManifest string) {
	if extra != nil && paletteBlocks != "" {
		core := strings.Replace(fullCSS, paletteBlocks, coreBlocks, 1)
		extra["palette_core_css"] = strings.Replace(core, "--palette-manifest: '"+manifest+"';", "--palette-manifest: '"+coreManifest+"';", 1)
	}
}

// writePaletteCoreCSS keeps the legacy full stylesheet intact while publishing
// a content-addressed entry point for the default light/dark pair.
func writePaletteCoreCSS(outputDir, css, hash string) error {
	if hash == "" {
		hash = fmt.Sprintf("%x", sha256.Sum256([]byte(css)))[:8]
	}
	dir := filepath.Join(outputDir, "css")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating palette core directory: %w", err)
	}
	for _, name := range []string{"palette-core.css", "palette-core." + hash + ".css"} {
		//nolint:gosec // Public stylesheet, readable by the web server.
		if err := os.WriteFile(filepath.Join(dir, name), []byte(css), 0o644); err != nil {
			return fmt.Errorf("writing palette core CSS: %w", err)
		}
	}
	return nil
}

// filterPalettes filters palettes based on switcher configuration.
func (p *PaletteCSSPlugin) filterPalettes(allPalettes []palettes.PaletteInfo, switcherConfig models.ThemeSwitcherConfig) []palettes.PaletteInfo {
	if switcherConfig.IsIncludeAll() {
		// Include all, then exclude specified
		excludeSet := make(map[string]bool)
		for _, name := range switcherConfig.Exclude {
			excludeSet[strings.ToLower(name)] = true
			excludeSet[normalizePaletteName(name)] = true
		}

		var result []palettes.PaletteInfo
		for _, info := range allPalettes {
			lowerName := strings.ToLower(info.Name)
			normalized := normalizePaletteName(info.Name)
			if !excludeSet[lowerName] && !excludeSet[normalized] {
				result = append(result, info)
			}
		}
		return result
	}

	// Include only specified palettes
	includeSet := make(map[string]bool)
	for _, name := range switcherConfig.Include {
		includeSet[strings.ToLower(name)] = true
		includeSet[normalizePaletteName(name)] = true
	}

	var result []palettes.PaletteInfo
	for _, info := range allPalettes {
		lowerName := strings.ToLower(info.Name)
		normalized := normalizePaletteName(info.Name)
		if includeSet[lowerName] || includeSet[normalized] {
			result = append(result, info)
		}
	}
	return result
}

// generatePaletteManifest builds the list of pickable themes for the switcher.
//
// Every explicit (non-derived) light or dark palette becomes an entry for its
// own mode, paired with the counterpart resolved by palettes.GetEffectivePalettes.
// A derived counterpart is added only when the pair has no explicit palette for
// the opposite mode, so each mode lists every family exactly once without
// duplicate or derived-of-derived entries. The configured default light and
// dark palettes are always included.
//
//nolint:gocyclo // Pairing explicit and derived palettes requires mode-specific fallback rules.
func (p *PaletteCSSPlugin) generatePaletteManifest(paletteInfos []palettes.PaletteInfo, defaultLight, defaultDark string) []PaletteManifestEntry {
	explicit := make(map[string]palettes.PaletteInfo)
	derived := make(map[string]palettes.PaletteInfo)
	for _, info := range paletteInfos {
		if info.Variant != palettes.VariantLight && info.Variant != palettes.VariantDark {
			continue
		}
		id := normalizePaletteName(info.Name)
		if info.Derived {
			derived[id] = info
		} else {
			explicit[id] = info
		}
	}

	ids := make([]string, 0, len(explicit))
	for id := range explicit {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	entries := make(map[string]PaletteManifestEntry)
	add := func(entry PaletteManifestEntry) {
		if entry.Name == "" {
			return
		}
		if _, exists := entries[entry.Name]; exists {
			return
		}
		entries[entry.Name] = entry
	}

	for _, id := range ids {
		info := explicit[id]
		variant := string(info.Variant)
		light, dark := palettes.GetEffectivePalettes(id, "", "")
		other := normalizePaletteName(dark)
		if info.Variant == palettes.VariantDark {
			other = normalizePaletteName(light)
		}
		if other == id {
			other = ""
		}
		label := paletteFamilyLabel(info.Name)
		add(PaletteManifestEntry{
			Name:        id,
			DisplayName: label,
			Variant:     variant,
			BaseName:    getBaseName(info.Name),
			Counterpart: other,
			loadName:    info.Name,
		})
		if other == "" {
			continue
		}
		if _, isExplicit := explicit[other]; isExplicit {
			continue
		}
		if dInfo, ok := derived[other]; ok {
			add(PaletteManifestEntry{
				Name:        other,
				DisplayName: label,
				Variant:     string(dInfo.Variant),
				BaseName:    getBaseName(info.Name),
				Counterpart: id,
				Derived:     true,
				loadName:    dInfo.Name,
			})
		}
	}

	// The site default must always be selectable, even when filtered out.
	for _, def := range []struct{ name, variant, other string }{
		{defaultLight, themeModeLight, defaultDark},
		{defaultDark, themeModeDark, defaultLight},
	} {
		if def.name == "" {
			continue
		}
		if _, ok := entries[def.name]; ok {
			continue
		}
		other := def.other
		if other == def.name {
			other = ""
		}
		add(PaletteManifestEntry{
			Name:        def.name,
			DisplayName: paletteFamilyLabel(strings.ReplaceAll(def.name, "-", " ")),
			Variant:     def.variant,
			BaseName:    getBaseName(def.name),
			Counterpart: other,
			loadName:    def.name,
		})
	}

	manifest := make([]PaletteManifestEntry, 0, len(entries))
	for _, entry := range entries {
		manifest = append(manifest, entry)
	}
	sort.Slice(manifest, func(i, j int) bool {
		li, lj := strings.ToLower(manifest[i].DisplayName), strings.ToLower(manifest[j].DisplayName)
		if li != lj {
			return li < lj
		}
		if manifest[i].Variant != manifest[j].Variant {
			// Light variants come before dark
			return manifest[i].Variant > manifest[j].Variant
		}
		return manifest[i].Name < manifest[j].Name
	})

	return manifest
}

// linkCounterparts clears counterpart references to palettes that are not
// part of the manifest, so the browser never switches to an unstyled name, and
// fills missing counterparts with the opposite-mode entry sharing the same
// family label when exactly one exists.
func linkCounterparts(manifest []PaletteManifestEntry) []PaletteManifestEntry {
	names := make(map[string]bool, len(manifest))
	byLabel := make(map[string][]int)
	for i, entry := range manifest {
		names[entry.Name] = true
		key := strings.ToLower(entry.DisplayName) + "|" + entry.Variant
		byLabel[key] = append(byLabel[key], i)
	}
	for i := range manifest {
		if manifest[i].Counterpart != "" && !names[manifest[i].Counterpart] {
			manifest[i].Counterpart = ""
		}
		if manifest[i].Counterpart != "" {
			continue
		}
		opposite := themeModeDark
		if manifest[i].Variant == themeModeDark {
			opposite = themeModeLight
		}
		if matches := byLabel[strings.ToLower(manifest[i].DisplayName)+"|"+opposite]; len(matches) == 1 {
			manifest[i].Counterpart = manifest[matches[0]].Name
		}
	}
	return manifest
}

// paletteFamilyLabel returns a mode-neutral display label for a palette name
// by trimming a trailing "Light"/"Dark" word ("Ayu Dark" -> "Ayu").
func paletteFamilyLabel(name string) string {
	label := strings.TrimSpace(name)
	lower := strings.ToLower(label)
	for _, suffix := range []string{" light", " dark", "-light", "-dark"} {
		if strings.HasSuffix(lower, suffix) && len(label) > len(suffix) {
			label = strings.TrimSpace(label[:len(label)-len(suffix)])
			break
		}
	}
	if label == "" {
		return name
	}
	words := strings.Fields(label)
	for i, word := range words {
		if word != "" && word[0] >= 'a' && word[0] <= 'z' {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// normalizePaletteName converts a palette name to a CSS-safe format.
func normalizePaletteName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "_", "-")
	return name
}

// getBaseName extracts the base name from a palette name (removes -light, -dark suffix).
func getBaseName(name string) string {
	lower := strings.ToLower(name)
	lower = strings.TrimSuffix(lower, "-light")
	lower = strings.TrimSuffix(lower, "-dark")
	return lower
}

// generateThemeCSSWithVariants generates CSS with both light and dark palette variants.
// The CSS uses data-theme attribute for switching between variants.
// It preserves all non-color variables (typography, spacing, layout, gradients) from the default theme.
func (p *PaletteCSSPlugin) generateThemeCSSWithVariants(lightPalette, darkPalette *palettes.Palette, lightName, darkName string, overrides map[string]string, fallbackMode string) string {
	var buf bytes.Buffer

	buf.WriteString("@layer reset, tokens, base, components, utilities, overrides;\n\n")
	buf.WriteString("@layer tokens {\n")

	// Header comment
	buf.WriteString("/* CSS Custom Properties - Generated by markata-go */\n")
	buf.WriteString(fmt.Sprintf("/* Light: %s, Dark: %s */\n\n", lightName, darkName))

	// Generate palette info CSS custom properties for JavaScript
	buf.WriteString("/* Palette configuration for JavaScript theme toggle */\n")
	buf.WriteString(":root {\n")
	// CSS custom property values need the quotes as part of the value for JS to read
	buf.WriteString(fmt.Sprintf("  --palette-light: %q;\n", lightName))
	buf.WriteString(fmt.Sprintf("  --palette-dark: %q;\n", darkName))

	// Write non-color variables (typography, spacing, layout, gradients)
	// These are written once in :root and don't change between light/dark modes

	buf.WriteString("}\n\n")

	if normalizeThemeFallbackMode(fallbackMode) == themeModeLight {
		if lightPalette != nil {
			buf.WriteString(fmt.Sprintf("/* Default light mode - %s */\n", lightPalette.Name))
			buf.WriteString(":root:not([data-theme=\"dark\"]),\n")
			buf.WriteString("[data-theme=\"light\"] {\n")
			p.writePaletteVariables(&buf, lightPalette)
			buf.WriteString("}\n\n")
		}

		if darkPalette != nil {
			buf.WriteString(fmt.Sprintf("/* Dark mode - %s */\n", darkPalette.Name))
			buf.WriteString("[data-theme=\"dark\"] {\n")
			p.writePaletteVariables(&buf, darkPalette)
			buf.WriteString("}\n\n")
		}
	} else {
		// Dark fallback (default)
		if darkPalette != nil {
			buf.WriteString(fmt.Sprintf("/* Dark mode - %s */\n", darkPalette.Name))
			buf.WriteString(":root:not([data-theme=\"light\"]),\n")
			buf.WriteString("[data-theme=\"dark\"] {\n")
			p.writePaletteVariables(&buf, darkPalette)
			buf.WriteString("}\n\n")
		}

		if lightPalette != nil {
			buf.WriteString(fmt.Sprintf("/* Light mode - %s */\n", lightPalette.Name))
			buf.WriteString("[data-theme=\"light\"] {\n")
			p.writePaletteVariables(&buf, lightPalette)
			buf.WriteString("}\n\n")
		}
	}

	buf.WriteString("}\n\n")
	p.writeThemeOverrides(&buf, overrides)

	return buf.String()
}

func (p *PaletteCSSPlugin) getThemeFallbackMode(extra map[string]interface{}) string {
	if extra == nil {
		return themeModeDark
	}

	if modelsConfig, ok := extra["models_config"].(*models.Config); ok {
		return normalizeThemeFallbackMode(modelsConfig.Theme.FallbackMode)
	}

	if themeConfig, ok := extra["theme"].(models.ThemeConfig); ok {
		return normalizeThemeFallbackMode(themeConfig.FallbackMode)
	}

	if theme, ok := extra["theme"].(map[string]interface{}); ok {
		if mode, ok := theme["fallback_mode"].(string); ok {
			return normalizeThemeFallbackMode(mode)
		}
	}

	return themeModeDark
}

func normalizeThemeFallbackMode(mode string) string {
	if strings.EqualFold(strings.TrimSpace(mode), themeModeLight) {
		return themeModeLight
	}
	return themeModeDark
}

// writePaletteVariables writes CSS custom properties for a palette.
func (p *PaletteCSSPlugin) writePaletteVariables(buf *bytes.Buffer, palette *palettes.Palette) {
	p.writePaletteVariablesIndented(buf, palette, "  ")
}

// resolveWithContrast resolves a color from the palette and adjusts it to meet WCAG contrast ratio against a background color.
func (p *PaletteCSSPlugin) resolveWithContrast(palette *palettes.Palette, fgKey string, minRatio float64) string {
	return p.resolveWithContrastAgainst(palette, fgKey, "bg-primary", minRatio)
}

func (p *PaletteCSSPlugin) resolveWithContrastAgainst(palette *palettes.Palette, fgKey, bgKey string, minRatio float64) string {
	fgHex := palette.Resolve(fgKey)
	if fgHex == "" {
		return ""
	}

	bgHex := palette.Resolve(bgKey)
	if bgHex == "" {
		return fgHex
	}

	fgColor, errFg := palettes.ParseHexColor(fgHex)
	bgColor, errBg := palettes.ParseHexColor(bgHex)
	if errFg != nil || errBg != nil {
		return fgHex
	}

	adjusted, _ := fgColor.AdjustForContrast(bgColor, minRatio)
	return adjusted.Hex()
}

// writePaletteVariablesIndented writes CSS custom properties with custom indentation.
//
//nolint:gocyclo // complexity is acceptable for a CSS generation function with many rules
func (p *PaletteCSSPlugin) writePaletteVariablesIndented(buf *bytes.Buffer, palette *palettes.Palette, indent string) {
	// Primary/accent colors
	if accent := palette.Resolve("accent"); accent != "" {
		fmt.Fprintf(buf, "%s--color-primary: %s;\n", indent, accent)
	}
	if accentHover := palette.Resolve("accent-hover"); accentHover != "" {
		fmt.Fprintf(buf, "%s--color-primary-light: %s;\n", indent, accentHover)
		fmt.Fprintf(buf, "%s--color-primary-dark: %s;\n", indent, accentHover)
	}

	fmt.Fprintf(buf, "\n%s/* Semantic colors */\n", indent)

	// Text colors
	if textPrimary := p.resolveWithContrast(palette, "text-primary", 4.5); textPrimary != "" {
		fmt.Fprintf(buf, "%s--color-text: %s;\n", indent, textPrimary)
	}
	if textSecondary := p.resolveWithContrast(palette, "text-secondary", 4.5); textSecondary != "" {
		fmt.Fprintf(buf, "%s--color-text-secondary: %s;\n", indent, textSecondary)
	}
	if textMuted := p.resolveWithContrast(palette, "text-muted", 4.5); textMuted != "" {
		fmt.Fprintf(buf, "%s--color-text-muted: %s;\n", indent, textMuted)
	}

	// Background colors
	if bgPrimary := palette.Resolve("bg-primary"); bgPrimary != "" {
		fmt.Fprintf(buf, "%s--color-background: %s;\n", indent, bgPrimary)
	}
	if bgSurface := palette.Resolve("bg-surface"); bgSurface != "" {
		fmt.Fprintf(buf, "%s--color-surface: %s;\n", indent, bgSurface)
	}

	// Border
	if border := palette.Resolve("border"); border != "" {
		fmt.Fprintf(buf, "%s--color-border: %s;\n", indent, border)
	}

	fmt.Fprintf(buf, "\n%s/* Status colors */\n", indent)

	// Status colors
	if success := p.resolveWithContrast(palette, "success", 3.0); success != "" {
		fmt.Fprintf(buf, "%s--color-success: %s;\n", indent, success)
	}
	if warning := p.resolveWithContrast(palette, "warning", 3.0); warning != "" {
		fmt.Fprintf(buf, "%s--color-warning: %s;\n", indent, warning)
	}
	if errorColor := p.resolveWithContrast(palette, "error", 3.0); errorColor != "" {
		fmt.Fprintf(buf, "%s--color-error: %s;\n", indent, errorColor)
	}
	if info := p.resolveWithContrast(palette, "info", 3.0); info != "" {
		fmt.Fprintf(buf, "%s--color-info: %s;\n", indent, info)
	}

	// Add link colors if available
	if link := p.resolveWithContrast(palette, "link", 4.5); link != "" {
		fmt.Fprintf(buf, "\n%s/* Link colors */\n", indent)
		fmt.Fprintf(buf, "%s--color-link: %s;\n", indent, link)
		if linkHover := p.resolveWithContrast(palette, "link-hover", 4.5); linkHover != "" {
			fmt.Fprintf(buf, "%s--color-link-hover: %s;\n", indent, linkHover)
		}
		if linkVisited := p.resolveWithContrast(palette, "link-visited", 4.5); linkVisited != "" {
			fmt.Fprintf(buf, "%s--color-link-visited: %s;\n", indent, linkVisited)
		}
	}

	// Add code colors if available
	if codeBg := palette.Resolve("code-bg"); codeBg != "" {
		fmt.Fprintf(buf, "\n%s/* Code colors */\n", indent)
		fmt.Fprintf(buf, "%s--color-code-bg: %s;\n", indent, codeBg)
		if codeText := p.resolveWithContrastAgainst(palette, "code-text", "code-bg", 4.5); codeText != "" {
			fmt.Fprintf(buf, "%s--color-code-text: %s;\n", indent, codeText)
		}
		for _, role := range []string{
			"comment",
			"keyword",
			"string",
			"number",
			"function",
			"type",
			"operator",
		} {
			componentName := "code-" + role
			if value := p.resolveWithContrastAgainst(palette, componentName, "code-bg", 4.5); value != "" {
				fmt.Fprintf(buf, "%s--color-%s: %s;\n", indent, componentName, value)
			}
		}
	}

	// Add admonition colors if available
	admonitionTypes := []string{
		"note", "info", "tip", "hint", "success",
		"warn", "warning", "caution", "important",
		"danger", "error", "bug",
		"example", "quote", "abstract",
		"chat", "chat-reply",
	}

	hasAdmonitions := false
	for _, adType := range admonitionTypes {
		if palette.Resolve("admonition-"+adType+"-bg") != "" ||
			palette.Resolve("admonition-"+adType+"-border") != "" {
			hasAdmonitions = true
			break
		}
	}

	if hasAdmonitions {
		fmt.Fprintf(buf, "\n%s/* Admonition colors */\n", indent)
		for _, adType := range admonitionTypes {
			if bg := palette.Resolve("admonition-" + adType + "-bg"); bg != "" {
				fmt.Fprintf(buf, "%s--admonition-%s-bg: %s;\n", indent, adType, bg)
			}
			if border := palette.Resolve("admonition-" + adType + "-border"); border != "" {
				fmt.Fprintf(buf, "%s--admonition-%s-border: %s;\n", indent, adType, border)
			}
		}
	}

	// Add mark/highlight colors
	// These are used for ==highlighted text== rendered as <mark> elements
	p.writeMarkColors(buf, palette, indent)
}

// writeMarkColors generates CSS variables for mark/highlight elements.
// If mark-bg/mark-text are not defined in the palette, computes them from the warning color.
//
//nolint:gocyclo // Mark colors require explicit fallback and contrast branches.
func (p *PaletteCSSPlugin) writeMarkColors(buf *bytes.Buffer, palette *palettes.Palette, indent string) {
	fmt.Fprintf(buf, "\n%s/* Mark/highlight colors */\n", indent)
	// Highlight semantics are palette roles, not Web Awesome warning tokens.
	// Keep the legacy mark names as compatibility aliases for existing themes.
	highlightBg := p.resolveWithContrast(palette, "accent", 4.5)
	highlightText := palette.Resolve("accent-ink")
	if highlightText == "" {
		highlightText = palette.Resolve("background")
	}
	if bgColor, err := palettes.ParseHexColor(highlightBg); err == nil {
		if textColor, err := palettes.ParseHexColor(highlightText); err == nil {
			if adjusted, ok := textColor.AdjustForContrast(bgColor, 4.5); ok {
				highlightText = adjusted.Hex()
			}
		}
	}
	if highlightBg != "" {
		fmt.Fprintf(buf, "%s--color-highlight: %s;\n", indent, highlightBg)
	}
	if highlightText != "" {
		fmt.Fprintf(buf, "%s--color-highlight-text: %s;\n", indent, highlightText)
	}
	// Keep heading contours in sync with the *effective* highlight tokens,
	// including author overrides applied after palette.css loads. The old
	// mark-presentation script did this by reading computed styles and writing
	// inline !important variables; CSS aliases avoid that forced reflow.
	if highlightBg != "" {
		fmt.Fprintf(buf, "%s--heading-mark-bg: var(--color-highlight);\n", indent)
	}
	if highlightText != "" {
		fmt.Fprintf(buf, "%s--heading-mark-fg: var(--color-highlight-text);\n", indent)
	}

	// Try explicit mark colors first
	markBg := palette.Resolve("mark-bg")
	markText := palette.Resolve("mark-text")

	// If not defined, compute from warning color
	if markBg == "" {
		warningHex := p.resolveWithContrast(palette, "warning", 3.0)
		if warningHex != "" {
			warningColor, err := palettes.ParseHexColor(warningHex)
			if err == nil {
				// Light palettes: lighten warning for subtle background
				// Dark palettes: darken warning for subtle background
				if palette.Variant == palettes.VariantLight {
					markBg = warningColor.Lighten(0.75).Hex()
				} else {
					markBg = warningColor.Darken(0.6).Hex()
				}
			}
		}
	}

	// Write mark-bg if we have it
	if markBg != "" {
		fmt.Fprintf(buf, "%s--color-mark-bg: %s;\n", indent, markBg)
	}

	// Compute mark-text with contrast checking
	if markText == "" && markBg != "" {
		bgColor, err := palettes.ParseHexColor(markBg)
		if err == nil {
			// Get the text color from the palette
			textHex := p.resolveWithContrast(palette, "text-primary", 4.5)
			if textHex == "" {
				// Fallback: use black for light bg, white for dark bg
				if bgColor.RelativeLuminance() > 0.5 {
					textHex = "#1a1a1a"
				} else {
					textHex = "#f5f5f5"
				}
			}

			textColor, err := palettes.ParseHexColor(textHex)
			if err == nil {
				// Ensure WCAG AA contrast (4.5:1 for normal text)
				adjustedText, ok := textColor.AdjustForContrast(bgColor, 4.5)
				if ok {
					markText = adjustedText.Hex()
				} else {
					// Couldn't adjust, use high contrast fallback
					if bgColor.RelativeLuminance() > 0.5 {
						markText = "#000000"
					} else {
						markText = "#ffffff"
					}
				}
			}
		}
	}

	// Write mark-text if we have it
	if markText != "" {
		fmt.Fprintf(buf, "%s--color-mark-text: %s;\n", indent, markText)
	}

	// Add selection colors (for user text selection)
	// These default to mark colors but can be overridden
	p.writeSelectionColors(buf, palette, indent)
}

// writeSelectionColors generates CSS variables for user text selection.
// Defaults to mark colors for consistency, but allows explicit override via palette.
func (p *PaletteCSSPlugin) writeSelectionColors(buf *bytes.Buffer, palette *palettes.Palette, indent string) {
	// Check for explicit selection colors
	selectionBg := palette.Resolve("selection-bg")
	selectionText := palette.Resolve("selection-text")

	// Only write if explicitly set (CSS fallback handles the mark color default)
	if selectionBg != "" || selectionText != "" {
		fmt.Fprintf(buf, "\n%s/* Selection colors */\n", indent)
		if selectionBg != "" {
			fmt.Fprintf(buf, "%s--color-selection-bg: %s;\n", indent, selectionBg)
		}
		if selectionText != "" {
			fmt.Fprintf(buf, "%s--color-selection-text: %s;\n", indent, selectionText)
		}
	}
}

// getPaletteConfig extracts palette configuration from config.Extra.
// Returns the base palette name and optional light/dark overrides.
func (p *PaletteCSSPlugin) getPaletteConfig(extra map[string]interface{}) (palette, paletteLight, paletteDark, seedColor string) {
	if extra == nil {
		return "", "", "", ""
	}

	if modelsConfig, ok := extra["models_config"].(*models.Config); ok {
		if modelsConfig.Theme.Palette != "" || modelsConfig.Theme.PaletteLight != "" || modelsConfig.Theme.PaletteDark != "" {
			return modelsConfig.Theme.Palette, modelsConfig.Theme.PaletteLight, modelsConfig.Theme.PaletteDark, modelsConfig.Theme.SeedColor
		}
	}

	// Check if theme is a models.ThemeConfig (from core.go)
	if themeConfig, ok := extra["theme"].(models.ThemeConfig); ok {
		return themeConfig.Palette, themeConfig.PaletteLight, themeConfig.PaletteDark, themeConfig.SeedColor
	}

	// Check if theme is a map[string]interface{} (from benchmark.go or raw TOML)
	theme, ok := extra["theme"].(map[string]interface{})
	if !ok {
		return "", "", "", ""
	}

	// Get base palette
	if pal, ok := theme["palette"].(string); ok && pal != "" {
		palette = pal
	}

	// Get optional light override
	if pal, ok := theme["palette_light"].(string); ok && pal != "" {
		paletteLight = pal
	}

	// Get optional dark override
	if pal, ok := theme["palette_dark"].(string); ok && pal != "" {
		paletteDark = pal
	}

	// Get optional seed color
	if seed, ok := theme["seed_color"].(string); ok && seed != "" {
		seedColor = seed
	}

	return palette, paletteLight, paletteDark, seedColor
}

func (p *PaletteCSSPlugin) getThemeVariables(extra map[string]interface{}) map[string]string {
	if extra == nil {
		return nil
	}

	if modelsConfig, ok := extra["models_config"].(*models.Config); ok {
		if len(modelsConfig.Theme.Variables) > 0 {
			return normalizeThemeVariables(modelsConfig.Theme.Variables)
		}
	}

	if themeConfig, ok := extra["theme"].(models.ThemeConfig); ok {
		return normalizeThemeVariables(themeConfig.Variables)
	}

	theme, ok := extra["theme"].(map[string]interface{})
	if !ok {
		return nil
	}

	if vars, ok := theme["variables"].(map[string]string); ok {
		return normalizeThemeVariables(vars)
	}

	if vars, ok := theme["variables"].(map[string]interface{}); ok {
		converted := make(map[string]string, len(vars))
		for key, value := range vars {
			str, ok := value.(string)
			if !ok {
				continue
			}
			converted[key] = str
		}
		return normalizeThemeVariables(converted)
	}

	return nil
}

func normalizeThemeVariables(vars map[string]string) map[string]string {
	if len(vars) == 0 {
		return nil
	}

	normalized := make(map[string]string, len(vars))
	for key, value := range vars {
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "" || value == "" {
			continue
		}
		if !strings.HasPrefix(key, "--") {
			continue
		}
		normalized[key] = value
	}

	if len(normalized) == 0 {
		return nil
	}

	return normalized
}

func (p *PaletteCSSPlugin) writeThemeOverrides(buf *bytes.Buffer, overrides map[string]string) {
	if len(overrides) == 0 {
		return
	}

	keys := make([]string, 0, len(overrides))
	for key := range overrides {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	buf.WriteString("@layer overrides {\n")
	buf.WriteString("  /* Theme variable overrides */\n")
	buf.WriteString("  :root {\n")
	for _, key := range keys {
		fmt.Fprintf(buf, "    %s: %s;\n", key, overrides[key])
	}
	buf.WriteString("  }\n")
	buf.WriteString("}\n")
}

// Priority returns the plugin priority for the write stage.
// Should run after static_assets so it can overwrite the default palette.css
func (p *PaletteCSSPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityDefault // After static_assets (PriorityEarly)
	}
	return lifecycle.PriorityDefault
}

// Ensure PaletteCSSPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*PaletteCSSPlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*PaletteCSSPlugin)(nil)
	_ lifecycle.WritePlugin     = (*PaletteCSSPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*PaletteCSSPlugin)(nil)
)
