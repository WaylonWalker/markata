package cmd

import (
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
)

var currentLogTheme = logging.DefaultTheme()

func configureCommandLogger(theme logging.Theme) error {
	format, err := logging.ParseFormat(logFormat)
	if err != nil {
		return err
	}
	currentLogTheme = theme

	logging.ConfigureStandardLogger(logging.Options{
		Writer:     errWriter(),
		Format:     format,
		ForceColor: forceColor,
		NoColor:    noColor,
		IsTTY:      errorOutputIsTerminal(),
		Theme:      theme,
	})

	return nil
}

func configureLoggerForManager(m *lifecycle.Manager) {
	theme := loggerThemeForManager(m)
	if err := configureCommandLogger(theme); err != nil {
		errlnf("Warning: failed to configure themed logging: %v", err)
	}
}

// configureServeLoggerForManager applies the site's normal logging theme and
// immediately installs Serve's terminal ownership and runtime observation
// policy. Keep this as the only manager-specific logger setup used by Serve,
// including rebuild managers.
func configureServeLoggerForManager(m *lifecycle.Manager) {
	theme := loggerThemeForManager(m)
	currentLogTheme = theme
	format, err := logging.ParseFormat(logFormat)
	if err != nil {
		format = logging.FormatAuto
	}
	configureServeLoggerWithOptions(theme, format)
	if err != nil {
		logging.Component("serve").Warnf("failed to configure themed logging: %v", err)
	}
}

func loggerThemeForManager(m *lifecycle.Manager) logging.Theme {
	theme := logging.DefaultTheme()
	if resolved, ok := resolveLoggerTheme(m); ok {
		theme = resolved
	}
	return theme
}

func resolveLoggerTheme(m *lifecycle.Manager) (logging.Theme, bool) {
	if m == nil || m.Config() == nil || m.Config().Extra == nil {
		return logging.Theme{}, false
	}

	modelsConfig, ok := m.Config().Extra["models_config"].(*models.Config)
	if !ok || modelsConfig == nil {
		return logging.Theme{}, false
	}

	palette, ok := loadLoggerPalette(modelsConfig.Theme, m.Config().Extra)
	if !ok {
		return logging.Theme{}, false
	}

	return logging.ThemeFromPalette(palette), true
}

func loadLoggerPalette(theme models.ThemeConfig, extra map[string]any) (*palettes.Palette, bool) {
	loader := palettes.NewLoader()
	if configPath, ok := extra["config_path"].(string); ok && configPath != "" {
		loader.AddPath(filepath.Join(filepath.Dir(configPath), "palettes"))
	}

	if theme.Palette == "generated" {
		if theme.SeedColor == "" {
			return nil, false
		}
		variant := palettes.VariantDark
		if strings.EqualFold(theme.FallbackMode, "light") {
			variant = palettes.VariantLight
		}
		generated, err := palettes.GenerateTriadicPalette(theme.SeedColor, variant)
		if err != nil {
			return nil, false
		}
		return generated, true
	}

	light, dark := palettes.GetEffectivePalettes(theme.Palette, theme.PaletteLight, theme.PaletteDark)
	selected := dark
	if strings.EqualFold(theme.FallbackMode, "light") {
		selected = light
	}
	if selected == "" {
		selected = theme.Palette
	}
	if selected == "" {
		return nil, false
	}

	palette, err := loader.Load(selected)
	if err != nil {
		return nil, false
	}
	return palette, true
}
