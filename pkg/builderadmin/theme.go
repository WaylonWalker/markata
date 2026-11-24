package builderadmin

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/servecontrol"
)

// uiTheme contains the semantic palette values used by the builder-admin UI.
// It has complete defaults so a missing or invalid site palette never prevents
// the build service from starting.
type uiTheme struct {
	Background string
	Panel      string
	Surface    string
	Elevated   string
	Text       string
	Muted      string
	Accent     string
	Link       string
	Border     string
	Focus      string
	Success    string
	Warning    string
	Error      string
	Info       string
	CodeBG     string
	CodeText   string
	ButtonBG   string
	ButtonText string
	IsDark     bool
}

func defaultUITheme() uiTheme {
	return uiTheme{
		Background: "#09090b", Panel: "#18181b", Surface: "#27272a", Elevated: "#3f3f46",
		Text: "#f4f4f5", Muted: "#a1a1aa", Accent: "#fafafa", Link: "#fafafa", Border: "#3f3f46",
		Focus: "#93c5fd", Success: "#22c55e", Warning: "#f59e0b", Error: "#ef4444", Info: "#3b82f6",
		CodeBG: "#111113", CodeText: "#f4f4f5", ButtonBG: "#fafafa", ButtonText: "#18181b", IsDark: true,
	}
}

func loadUITheme(cfg Config) uiTheme {
	theme := defaultUITheme()
	configPath := resolveThemeConfigPath(cfg)
	if configPath == "" {
		return theme
	}

	siteConfig, err := config.Load(configPath)
	if err != nil {
		return theme
	}
	loader := palettes.NewLoader()
	loader.AddPath(filepath.Join(cfg.SourceDir, "palettes"))
	if siteConfig.Theme.Palette == "generated" && siteConfig.Theme.SeedColor != "" {
		variant := palettes.VariantDark
		if strings.EqualFold(siteConfig.Theme.FallbackMode, "light") {
			variant = palettes.VariantLight
		}
		if generated, err := palettes.GenerateTriadicPalette(siteConfig.Theme.SeedColor, variant); err == nil {
			loader.AddPalette("generated-"+string(variant), generated)
		}
	}
	palette, ok := loadThemePalette(loader, siteConfig.Theme)
	if !ok {
		return theme
	}

	colors := servecontrol.BrowserTheme(palette.Resolve, palette.Variant != palettes.VariantLight)
	theme.IsDark = colors["mode"] == "dark"
	theme.Background = colors["background"]
	theme.Panel = colors["panel"]
	theme.Surface = colors["surface"]
	theme.Elevated = colors["elevated"]
	theme.Text = colors["text-primary"]
	theme.Muted = colors["text-secondary"]
	theme.Accent = colors["accent"]
	theme.Link = colors["link"]
	theme.Border = colors["border"]
	theme.Focus = colors["focus"]
	theme.Success = colors["success"]
	theme.Warning = colors["warning"]
	theme.Error = colors["error"]
	theme.Info = colors["info"]
	theme.CodeBG = colors["code-background"]
	theme.CodeText = colors["code-text"]
	theme.ButtonBG = colors["button-background"]
	theme.ButtonText = colors["button-text"]
	return theme
}

func resolveThemeConfigPath(cfg Config) string {
	if cfg.ConfigPath == "" {
		return sourceConfigPath(cfg.SourceDir)
	}
	if filepath.IsAbs(cfg.ConfigPath) {
		return cfg.ConfigPath
	}
	return filepath.Join(cfg.SourceDir, cfg.ConfigPath)
}

func sourceConfigPath(sourceDir string) string {
	for _, name := range []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"} {
		path := filepath.Join(sourceDir, name)
		if fileExists(path) {
			return path
		}
	}
	return ""
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func selectedPaletteName(theme models.ThemeConfig) string {
	light, dark := palettes.GetEffectivePalettes(theme.Palette, theme.PaletteLight, theme.PaletteDark)
	if strings.EqualFold(theme.FallbackMode, "light") {
		return light
	}
	return dark
}

func loadThemePalette(loader *palettes.Loader, theme models.ThemeConfig) (*palettes.Palette, bool) {
	mode := "dark"
	if strings.EqualFold(theme.FallbackMode, "light") {
		mode = "light"
	}
	light, dark := palettes.GetEffectivePalettes(theme.Palette, theme.PaletteLight, theme.PaletteDark)
	candidates := []string{selectedPaletteName(theme)}
	if mode == "light" {
		candidates = append(candidates, light)
	} else {
		candidates = append(candidates, dark)
	}
	if !strings.HasSuffix(theme.Palette, "-light") && !strings.HasSuffix(theme.Palette, "-dark") {
		candidates = append(candidates, theme.Palette+"-"+mode)
	}
	candidates = append(candidates, theme.Palette)
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		palette, err := loader.Load(candidate)
		if err == nil {
			return palette, true
		}
	}
	return nil, false
}

func browserThemeMap(theme uiTheme) map[string]string {
	mode := "dark"
	if !theme.IsDark {
		mode = "light"
	}
	return map[string]string{
		"mode": mode, "background": theme.Background, "panel": theme.Panel,
		"surface": theme.Surface, "elevated": theme.Elevated,
		"text-primary": theme.Text, "text-secondary": theme.Muted,
		"accent": theme.Accent, "link": theme.Link, "border": theme.Border,
		"focus": theme.Focus, "success": theme.Success, "warning": theme.Warning,
		"error": theme.Error, "info": theme.Info, "code-background": theme.CodeBG,
		"code-text": theme.CodeText, "button-background": theme.ButtonBG,
		"button-text": theme.ButtonText,
	}
}
