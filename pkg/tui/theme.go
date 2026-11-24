package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/terminaltheme"
)

// Colors holds resolved colors for TUI styling.
// These are mapped from palette semantic colors to TUI-specific uses.
type Colors struct {
	// Header/title colors
	Header lipgloss.Color

	// Subtle/muted text colors
	Subtle lipgloss.Color

	// Selected/highlighted item colors
	Selected     lipgloss.Color
	SelectedBg   lipgloss.Color
	SelectedText lipgloss.Color

	// Border colors
	Border      lipgloss.Color
	BorderFocus lipgloss.Color

	// Table colors
	TableHeader   lipgloss.Color
	TableCell     lipgloss.Color
	TableSelected lipgloss.Color

	// Feed status colors
	FeedPrivate   lipgloss.Color
	FeedHidden    lipgloss.Color
	FeedAutomated lipgloss.Color

	// Menu colors
	MenuFg lipgloss.Color
	MenuBg lipgloss.Color
}

// DefaultColors returns the shared fallback colors when no palette is configured.
func DefaultColors() *Colors {
	return colorsFromTerminalTheme(terminaltheme.Default())
}

func colorsFromTerminalTheme(theme terminaltheme.Theme) *Colors {
	return &Colors{
		Header:        lipgloss.Color(theme.Accent),
		Subtle:        lipgloss.Color(theme.TextMuted),
		Selected:      lipgloss.Color(theme.Link),
		SelectedBg:    lipgloss.Color(theme.Elevated),
		SelectedText:  lipgloss.Color(theme.TextPrimary),
		Border:        lipgloss.Color(theme.Border),
		BorderFocus:   lipgloss.Color(theme.BorderFocus),
		TableHeader:   lipgloss.Color(theme.Accent),
		TableCell:     lipgloss.Color(theme.TextPrimary),
		TableSelected: lipgloss.Color(theme.Link),
		FeedPrivate:   lipgloss.Color(theme.Error),
		FeedHidden:    lipgloss.Color(theme.Warning),
		FeedAutomated: lipgloss.Color(theme.Info),
		MenuFg:        lipgloss.Color(theme.TextPrimary),
		MenuBg:        lipgloss.Color(theme.Background),
	}
}

// LoadColors loads TUI colors from a palette.
// If the palette cannot be loaded, returns default colors.
func LoadColors(paletteName string) *Colors {
	if paletteName == "" {
		return DefaultColors()
	}

	loader := palettes.NewLoader()
	palette, err := loader.Load(paletteName)
	if err != nil {
		return DefaultColors()
	}

	return colorsFromPalette(palette)
}

// ColorsFromPalette maps a resolved palette to the shared TUI color roles.
func ColorsFromPalette(palette *palettes.Palette) *Colors {
	return colorsFromPalette(palette)
}

// colorsFromPalette extracts TUI colors from a loaded palette.
func colorsFromPalette(p *palettes.Palette) *Colors {
	return colorsFromTerminalTheme(terminaltheme.FromPalette(p))
}

// Theme holds the Lipgloss styles for the TUI.
// Styles are initialized from Colors.
type Theme struct {
	Colors *Colors

	// Pre-built styles
	HeaderStyle       lipgloss.Style
	SubtleStyle       lipgloss.Style
	SelectedStyle     lipgloss.Style
	DetailLabelStyle  lipgloss.Style
	DetailBoxStyle    lipgloss.Style
	DetailStatusStyle lipgloss.Style
	SortMenuStyle     lipgloss.Style
}

// NewTheme creates a new Theme from colors.
func NewTheme(colors *Colors) *Theme {
	if colors == nil {
		colors = DefaultColors()
	}

	return &Theme{
		Colors: colors,
		HeaderStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Header),
		SubtleStyle: lipgloss.NewStyle().
			Foreground(colors.Subtle),
		SelectedStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Selected),
		DetailLabelStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Header).
			Width(12),
		DetailBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.BorderFocus).
			Padding(1, 0),
		DetailStatusStyle: lipgloss.NewStyle().
			Foreground(colors.Subtle).
			Padding(0, 1),
		SortMenuStyle: lipgloss.NewStyle().
			Foreground(colors.MenuFg).
			Background(colors.MenuBg),
	}
}

// DefaultTheme returns a theme with default colors.
func DefaultTheme() *Theme {
	return NewTheme(DefaultColors())
}

// GetPaletteNameFromConfig extracts the palette name from config.Extra.
func GetPaletteNameFromConfig(extra map[string]interface{}) string {
	if extra == nil {
		return ""
	}

	theme, ok := extra["theme"].(map[string]interface{})
	if !ok {
		return ""
	}

	palette, ok := theme["palette"].(string)
	if !ok {
		return ""
	}

	return palette
}
