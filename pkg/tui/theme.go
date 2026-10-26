package tui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/WaylonWalker/markata/pkg/palettes"
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

// DefaultColors returns fallback colors when no palette is configured.
// These match the original hardcoded ANSI codes for backward compatibility.
func DefaultColors() *Colors {
	return &Colors{
		Header:        lipgloss.Color("99"),  // Purple
		Subtle:        lipgloss.Color("241"), // Gray
		Selected:      lipgloss.Color("212"), // Pink
		SelectedBg:    lipgloss.Color("57"),  // Purple bg
		SelectedText:  lipgloss.Color("229"), // Light yellow
		Border:        lipgloss.Color("240"), // Gray
		BorderFocus:   lipgloss.Color("99"),  // Purple
		TableHeader:   lipgloss.Color("99"),  // Purple
		TableCell:     lipgloss.Color("252"), // Light gray
		TableSelected: lipgloss.Color("229"), // Light yellow
		FeedPrivate:   lipgloss.Color("203"), // Red/pink
		FeedHidden:    lipgloss.Color("214"), // Orange
		FeedAutomated: lipgloss.Color("81"),  // Cyan
		MenuFg:        lipgloss.Color("252"), // Light gray
		MenuBg:        lipgloss.Color("236"), // Dark gray
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

// colorsFromPalette extracts TUI colors from a loaded palette.
func colorsFromPalette(p *palettes.Palette) *Colors {
	colors := DefaultColors()

	// Map palette semantic colors to TUI colors
	// Use accent for header styling
	if hex := p.Resolve("accent"); hex != "" {
		colors.Header = lipgloss.Color(hex)
		colors.BorderFocus = lipgloss.Color(hex)
		colors.TableHeader = lipgloss.Color(hex)
	}

	// Use text-muted for subtle styling
	if hex := p.Resolve("text-muted"); hex != "" {
		colors.Subtle = lipgloss.Color(hex)
	}

	// Use link color for selected items
	if hex := p.Resolve("link"); hex != "" {
		colors.Selected = lipgloss.Color(hex)
		colors.TableSelected = lipgloss.Color(hex)
	}
	if hex := p.Resolve("error"); hex != "" {
		colors.FeedPrivate = lipgloss.Color(hex)
	}
	if hex := p.Resolve("warning"); hex != "" {
		colors.FeedHidden = lipgloss.Color(hex)
	}
	if hex := p.Resolve("info"); hex != "" {
		colors.FeedAutomated = lipgloss.Color(hex)
	}

	// Use accent-hover for selection background
	if hex := p.Resolve("accent-hover"); hex != "" {
		colors.SelectedBg = lipgloss.Color(hex)
	} else if hex := p.Resolve("bg-elevated"); hex != "" {
		colors.SelectedBg = lipgloss.Color(hex)
	}

	// Use text-primary for selected text
	if hex := p.Resolve("text-primary"); hex != "" {
		colors.SelectedText = lipgloss.Color(hex)
	}

	// Use border colors
	if hex := p.Resolve("border"); hex != "" {
		colors.Border = lipgloss.Color(hex)
	}

	// Use text-primary for table cells
	if hex := p.Resolve("text-primary"); hex != "" {
		colors.TableCell = lipgloss.Color(hex)
	}

	// Menu colors from nav or surface colors
	if hex := p.Resolve("text-primary"); hex != "" {
		colors.MenuFg = lipgloss.Color(hex)
	}
	if hex := p.Resolve("bg-surface"); hex != "" {
		colors.MenuBg = lipgloss.Color(hex)
	}

	return colors
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
