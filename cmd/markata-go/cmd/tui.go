package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/listcache"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/services"
	"github.com/WaylonWalker/markata/pkg/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive terminal UI",
	Long: `Launch an interactive terminal user interface for browsing and managing
your markata-go site content.

The TUI provides:
  - Post list with filtering and sorting
  - Tag and feed browsing
  - Quick editing via $EDITOR
  - Build triggering

Navigation:
  j/k or ↑/↓  Move selection
  Enter       View post details
  /           Filter posts
  :           Command mode
  q           Quit`,
	RunE: runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(_ *cobra.Command, _ []string) error {
	// Create the manager using the existing helper
	manager, err := createManager(cfgFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create services app
	app := services.NewApp(manager)
	configHash, err := configFilesHash(cfgFile, mergeConfigFiles)
	if err != nil {
		return err
	}
	listcache.SetOptions(manager, listcache.Options{
		CacheDir:   listcache.DefaultCacheDir,
		ConfigHash: configHash,
	})

	// Load theme from configuration
	var palette *palettes.Palette
	if cfg := getModelsConfig(manager); cfg != nil {
		palette, _ = loadLoggerPalette(cfg.Theme, manager.Config().Extra)
	}
	colors := tui.ColorsFromPalette(palette)
	theme := tui.NewTheme(colors)

	// Create and run TUI with theme
	model := tui.NewModelWithTheme(app, theme)
	p := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
