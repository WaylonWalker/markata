package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/palettes"
)

// palettePickCmd opens an interactive fuzzy picker for palettes.
var palettePickCmd = &cobra.Command{
	Use:   "pick",
	Short: "Interactively pick a palette and set it in config",
	Long: `Open an interactive fuzzy picker to browse and select a color palette.

Shows color swatches for each palette as you navigate. Type to fuzzy-filter
the palette list. Press Enter to select and set the palette in your config.

Use --no-set to only print the palette name without modifying config.

Example usage:
  markata-go palette pick            # Pick and set in config
  markata-go palette pick --no-set   # Only print the name`,
	RunE: runPalettePickCommand,
}

// palettePickNoSet skips writing the config after picking.
var palettePickNoSet bool

func init() {
	paletteCmd.AddCommand(palettePickCmd)
	palettePickCmd.Flags().BoolVar(&palettePickNoSet, "no-set", false, "Only print the palette name without updating config")
	configurePalettePickerTheme(logging.DefaultTheme())
}

// runPalettePickCommand launches the interactive palette picker.
func runPalettePickCommand(_ *cobra.Command, _ []string) error {
	loader := palettes.NewLoader()
	infos, err := loader.Discover()
	if err != nil {
		return fmt.Errorf("failed to discover palettes: %w", err)
	}

	if len(infos) == 0 {
		return fmt.Errorf("no palettes found")
	}
	configurePalettePickerTheme(currentLogTheme)

	// Pre-load all palettes for preview.
	loaded := make(map[string]*palettes.Palette, len(infos))
	for _, info := range infos {
		p, loadErr := loader.Load(info.Name)
		if loadErr == nil {
			loaded[info.Name] = p
		}
	}

	m := newPickerModel(infos, loaded)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return fmt.Errorf("picker error: %w", err)
	}

	final, ok := result.(pickerModel)
	if !ok || final.canceled {
		return nil
	}

	fmt.Println(final.chosen)

	if !palettePickNoSet {
		configPath := cfgFile
		if configPath == "" {
			var discoverErr error
			configPath, discoverErr = config.Discover()
			if discoverErr != nil {
				return fmt.Errorf("no config file found; create one with 'markata-go init': %w", discoverErr)
			}
		}

		slug := normalizeFileName(final.chosen)
		if err := setPaletteInConfig(configPath, slug); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "Set palette to %q in %s\n", slug, configPath)
	}

	return nil
}

func setPaletteInConfig(configPath, paletteName string) error {
	return config.SetValueInFile(configPath, "theme.palette", paletteName)
}

// --------------------------------------------------------------------
// Bubble Tea model
// --------------------------------------------------------------------

type pickerModel struct {
	// Data
	allInfos []palettes.PaletteInfo
	filtered []palettes.PaletteInfo
	loaded   map[string]*palettes.Palette

	// UI state
	query    string
	cursor   int
	chosen   string
	canceled bool

	// Terminal dimensions
	width  int
	height int
}

func newPickerModel(infos []palettes.PaletteInfo, loaded map[string]*palettes.Palette) pickerModel {
	return pickerModel{
		allInfos: infos,
		filtered: infos,
		loaded:   loaded,
	}
}

func (m pickerModel) Init() tea.Cmd { return nil }

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m pickerModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEscape:
		m.canceled = true
		return m, tea.Quit

	case tea.KeyEnter:
		if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
			m.chosen = m.filtered[m.cursor].Name
		}
		return m, tea.Quit

	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil

	case tea.KeyDown:
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
		return m, nil

	case tea.KeyPgUp:
		m.cursor -= listPageSize(m.height)
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case tea.KeyPgDown:
		m.cursor += listPageSize(m.height)
		if m.cursor >= len(m.filtered) {
			m.cursor = len(m.filtered) - 1
		}
		if m.cursor < 0 {
			m.cursor = 0
		}
		return m, nil

	case tea.KeyBackspace, tea.KeyDelete:
		if m.query != "" {
			m.query = m.query[:len(m.query)-1]
			m.applyFilter()
		}
		return m, nil

	default:
		if msg.Type == tea.KeyRunes {
			m.query += string(msg.Runes)
			m.applyFilter()
		}
		return m, nil
	}
}

func (m *pickerModel) applyFilter() {
	if m.query == "" {
		m.filtered = m.allInfos
		m.cursor = 0
		return
	}

	q := strings.ToLower(m.query)
	var result []palettes.PaletteInfo
	for _, info := range m.allInfos {
		if fuzzyMatch(strings.ToLower(info.Name), q) {
			result = append(result, info)
		}
	}
	m.filtered = result
	if m.cursor >= len(m.filtered) {
		m.cursor = len(m.filtered) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// fuzzyMatch returns true if all characters in pattern appear in str in order.
func fuzzyMatch(str, pattern string) bool {
	pi := 0
	for si := 0; si < len(str) && pi < len(pattern); si++ {
		if str[si] == pattern[pi] {
			pi++
		}
	}
	return pi == len(pattern)
}

// listPageSize returns the number of visible list items for page scrolling.
// Matches the listHeight calculation: termHeight - 2 (help) - 2 (border) - 6 (fixed content).
func listPageSize(termHeight int) int {
	size := termHeight - 10
	if size < 5 {
		return 5
	}
	return size
}

// --------------------------------------------------------------------
// View rendering
// --------------------------------------------------------------------

// Style definitions for the picker.
var (
	pickerTitleStyle         lipgloss.Style
	pickerPromptStyle        lipgloss.Style
	pickerCursorStyle        lipgloss.Style
	pickerItemStyle          lipgloss.Style
	pickerDimStyle           lipgloss.Style
	pickerVariantDarkStyle   lipgloss.Style
	pickerVariantLightStyle  lipgloss.Style
	pickerPreviewHeaderStyle lipgloss.Style
	pickerPreviewLabelStyle  lipgloss.Style
	pickerPreviewValueStyle  lipgloss.Style
	pickerBorderStyle        lipgloss.Style
	pickerHelpStyle          lipgloss.Style
)

func configurePalettePickerTheme(theme logging.Theme) {
	pickerTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Component))
	pickerPromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Component))
	pickerCursorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Success))
	pickerItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Component))
	pickerDimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Timestamp))
	pickerVariantDarkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Timestamp))
	pickerVariantLightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Warning))
	pickerPreviewHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(theme.Component))
	pickerPreviewLabelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Timestamp))
	pickerPreviewValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Component))
	pickerBorderStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color(theme.Timestamp))
	pickerHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Timestamp))
}

func (m pickerModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	// Reserve: 1 help bar + 1 newline before help = 2 lines outside panels
	panelHeight := m.height - 2

	// Layout: left list panel + right preview panel
	listWidth := m.width * 2 / 5
	if listWidth < 25 {
		listWidth = 25
	}
	previewWidth := m.width - listWidth - 4 // account for borders and gap
	if previewWidth < 30 {
		previewWidth = 30
	}

	listPanel := m.renderList(listWidth, panelHeight)
	previewPanel := m.renderPreview(previewWidth, panelHeight)

	// Join side by side
	joined := lipgloss.JoinHorizontal(lipgloss.Top, listPanel, "  ", previewPanel)

	// Help bar at bottom
	help := pickerHelpStyle.Render("  Type to filter | Up/Down to navigate | Enter to select | Esc to cancel")

	return joined + "\n" + help
}

func (m pickerModel) renderList(width, maxHeight int) string {
	var sb strings.Builder

	// Title
	sb.WriteString(pickerTitleStyle.Render("  Palette Picker"))
	sb.WriteString("\n\n")

	// Search prompt
	promptLine := pickerPromptStyle.Render("  > ") + pickerItemStyle.Render(m.query)
	cursorChar := pickerDimStyle.Render("|")
	sb.WriteString(promptLine + cursorChar + "\n\n")

	// Calculate visible range (scrolling)
	// Inside the border we use: title(1) + blank(1) + prompt(1) + blank(1) + items + blank(1) + count(1)
	// The border adds 2 lines (top+bottom), so content height = maxHeight - 2
	// Fixed content lines = 6 (title, blank, prompt, blank, blank-before-count, count)
	listHeight := maxHeight - 2 - 6
	if listHeight < 3 {
		listHeight = 3
	}

	start, end := scrollRange(m.cursor, len(m.filtered), listHeight)

	// Scroll indicator top
	if start > 0 {
		sb.WriteString(pickerDimStyle.Render(fmt.Sprintf("  ... %d more above", start)))
		sb.WriteString("\n")
	}

	// Render items
	for i := start; i < end; i++ {
		info := m.filtered[i]
		name := info.Name

		// Truncate to fit width
		maxName := width - 12
		if maxName < 10 {
			maxName = 10
		}
		if len(name) > maxName {
			name = name[:maxName-3] + "..."
		}

		// Variant badge
		var variantBadge string
		if info.Variant == palettes.VariantDark {
			variantBadge = pickerVariantDarkStyle.Render(" [dark]")
		} else {
			variantBadge = pickerVariantLightStyle.Render(" [light]")
		}

		if i == m.cursor {
			sb.WriteString(pickerCursorStyle.Render("  > "+name) + variantBadge)
		} else {
			sb.WriteString(pickerItemStyle.Render("    "+name) + variantBadge)
		}
		sb.WriteString("\n")
	}

	// Scroll indicator bottom
	if end < len(m.filtered) {
		sb.WriteString(pickerDimStyle.Render(fmt.Sprintf("  ... %d more below", len(m.filtered)-end)))
		sb.WriteString("\n")
	}

	// Count indicator
	sb.WriteString("\n")
	countText := fmt.Sprintf("  %d/%d palettes", len(m.filtered), len(m.allInfos))
	sb.WriteString(pickerDimStyle.Render(countText))

	return pickerBorderStyle.Width(width).Height(maxHeight - 2).Render(sb.String())
}

func (m pickerModel) renderPreview(width, maxHeight int) string {
	if len(m.filtered) == 0 || m.cursor >= len(m.filtered) {
		return pickerBorderStyle.Width(width).Height(maxHeight - 2).Render(
			pickerDimStyle.Render("  No palette selected"))
	}

	info := m.filtered[m.cursor]
	p, ok := m.loaded[info.Name]
	if !ok {
		return pickerBorderStyle.Width(width).Height(maxHeight - 2).Render(
			pickerDimStyle.Render("  Could not load palette"))
	}

	var sb strings.Builder

	// Palette name header
	sb.WriteString(pickerPreviewHeaderStyle.Render("  " + p.Name))
	sb.WriteString("\n")

	// Metadata
	if p.Author != "" {
		sb.WriteString(pickerPreviewLabelStyle.Render("  Author: "))
		sb.WriteString(pickerPreviewValueStyle.Render(p.Author))
		sb.WriteString("\n")
	}
	if p.Description != "" {
		desc := p.Description
		maxDesc := width - 6
		if maxDesc > 0 && len(desc) > maxDesc {
			desc = desc[:maxDesc-3] + "..."
		}
		sb.WriteString(pickerPreviewLabelStyle.Render("  Desc:   "))
		sb.WriteString(pickerPreviewValueStyle.Render(desc))
		sb.WriteString("\n")
	}
	sb.WriteString(pickerPreviewLabelStyle.Render("  Source: "))
	sb.WriteString(pickerPreviewValueStyle.Render(p.Source))
	sb.WriteString("\n\n")

	// Color swatches - raw colors
	sb.WriteString(pickerPreviewLabelStyle.Render("  Colors"))
	sb.WriteString("\n")
	sb.WriteString(renderColorSwatches(p.Colors, width-4))
	sb.WriteString("\n")

	// Semantic preview swatches
	if len(p.Semantic) > 0 {
		sb.WriteString(pickerPreviewLabelStyle.Render("  Semantic"))
		sb.WriteString("\n")
		sb.WriteString(renderSemanticSwatches(p, width-4))
	}

	// Contrast preview block
	sb.WriteString("\n")
	sb.WriteString(renderContrastPreview(p, width-4))

	return pickerBorderStyle.Width(width).Height(maxHeight - 2).Render(sb.String())
}

// renderColorSwatches renders a grid of color swatch blocks with labels.
func renderColorSwatches(colors map[string]string, maxWidth int) string {
	// Sort color names for stable rendering.
	names := make([]string, 0, len(colors))
	for name := range colors {
		names = append(names, name)
	}
	sort.Strings(names)

	var sb strings.Builder
	// Each swatch: 2 chars wide colored block + space = 3 chars minimum,
	// but we want to show name too for key colors. Use a compact grid layout.
	swatchWidth := 4 // "██" occupies 2 cells + 2 padding
	cols := (maxWidth - 2) / swatchWidth
	if cols < 4 {
		cols = 4
	}
	if cols > len(names) {
		cols = len(names)
	}

	// Render swatch grid rows
	for i := 0; i < len(names); i += cols {
		sb.WriteString("  ")
		end := i + cols
		if end > len(names) {
			end = len(names)
		}
		for j := i; j < end; j++ {
			hex := colors[names[j]]
			swatch := lipgloss.NewStyle().
				Foreground(lipgloss.Color(hex)).
				Render("██")
			sb.WriteString(swatch)
			sb.WriteString("  ")
		}
		sb.WriteString("\n")

		// Render names below the swatches
		sb.WriteString("  ")
		for j := i; j < end; j++ {
			name := names[j]
			// Truncate name to fit swatch width
			display := name
			if len(display) > swatchWidth {
				display = display[:swatchWidth-1] + ""
			}
			// Pad to swatch width
			for len(display) < swatchWidth {
				display += " "
			}
			sb.WriteString(pickerDimStyle.Render(display))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderSemanticSwatches renders labeled semantic color swatches.
func renderSemanticSwatches(p *palettes.Palette, maxWidth int) string {
	// Show the most important semantic colors.
	important := []struct {
		key   string
		label string
	}{
		{"text-primary", "text"},
		{"bg-primary", "bg"},
		{"accent", "accent"},
		{"link", "link"},
		{"success", "ok"},
		{"warning", "warn"},
		{"error", "err"},
		{"border", "bord"},
	}

	var sb strings.Builder
	sb.WriteString("  ")

	shown := 0
	swatchWidth := 6
	cols := (maxWidth - 2) / swatchWidth
	if cols < 1 {
		cols = 1
	}

	for _, item := range important {
		hex := p.Resolve(item.key)
		if hex == "" {
			continue
		}
		if shown > 0 && shown%cols == 0 {
			sb.WriteString("\n  ")
		}

		swatch := lipgloss.NewStyle().
			Foreground(lipgloss.Color(hex)).
			Render("██")

		label := item.label
		// Pad label
		for len(label) < swatchWidth-3 {
			label += " "
		}

		sb.WriteString(swatch + " " + pickerDimStyle.Render(label))
		shown++
	}
	sb.WriteString("\n")

	return sb.String()
}

// renderContrastPreview renders a sample text block using the palette's
// primary text/background colors to show real contrast.
func renderContrastPreview(p *palettes.Palette, maxWidth int) string {
	bgHex := p.Resolve("bg-primary")
	fgHex := p.Resolve("text-primary")
	if bgHex == "" || fgHex == "" {
		return ""
	}

	contentWidth := maxWidth - 4
	if contentWidth < 20 {
		contentWidth = 20
	}

	textStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(bgHex)).
		Foreground(lipgloss.Color(fgHex))

	// Build preview lines
	line1 := padRight("  Sample text on bg-primary  ", contentWidth)
	line2 := padRight("  The quick brown fox jumps.  ", contentWidth)

	// Link preview
	linkHex := p.Resolve("link")
	var line3 string
	if linkHex != "" {
		linkStyle := lipgloss.NewStyle().
			Background(lipgloss.Color(bgHex)).
			Foreground(lipgloss.Color(linkHex)).
			Underline(true)
		prefix := textStyle.Render(padRight("  Link: ", 8))
		linkText := linkStyle.Render("https://example.com")
		remaining := contentWidth - 8 - 19
		if remaining < 0 {
			remaining = 0
		}
		pad := textStyle.Render(strings.Repeat(" ", remaining))
		line3 = prefix + linkText + pad
	} else {
		line3 = textStyle.Render(padRight("", contentWidth))
	}

	preview := textStyle.Render(line1) + "\n" +
		textStyle.Render(line2) + "\n" +
		line3

	return "  " + pickerPreviewLabelStyle.Render("Preview") + "\n" +
		"  " + strings.ReplaceAll(preview, "\n", "\n  ") + "\n"
}

// padRight pads a string to the given width with spaces.
func padRight(s string, width int) string {
	if len(s) >= width {
		return s[:width]
	}
	return s + strings.Repeat(" ", width-len(s))
}

// scrollRange calculates the visible start/end indices for a scrollable list.
func scrollRange(cursor, total, height int) (start, end int) {
	if total <= height {
		return 0, total
	}

	// Keep cursor roughly centered
	half := height / 2
	start = cursor - half
	if start < 0 {
		start = 0
	}
	end = start + height
	if end > total {
		end = total
		start = end - height
		if start < 0 {
			start = 0
		}
	}

	return start, end
}
