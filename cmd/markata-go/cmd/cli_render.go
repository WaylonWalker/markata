package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const (
	cliPreferredRuleWidth = 80
	cliNarrowRuleWidth    = 32
)

var cliTerminalWidth = func(writer io.Writer) int {
	file, ok := writer.(*os.File)
	if !ok {
		return cliPreferredRuleWidth
	}
	width, _, err := term.GetSize(int(file.Fd()))
	if err != nil || width <= 0 {
		return cliPreferredRuleWidth
	}
	return width
}

// cliSectionRule returns a terminal-width-aware Unicode section separator.
// It caps rules at 80 columns and uses a compact label on narrow terminals.
func cliSectionRule(label string) string {
	width := min(cliTerminalWidth(outWriter()), cliPreferredRuleWidth)
	if width <= 4 {
		return "─"
	}
	if width < cliNarrowRuleWidth {
		return "─ " + truncateCLILabel(label, width-4) + " ─"
	}

	label = " " + label + " "
	labelWidth := utf8.RuneCountInString(label)
	if labelWidth+4 > width {
		return "─ " + truncateCLILabel(strings.TrimSpace(label), width-4) + " ─"
	}
	left := 3
	right := width - left - labelWidth
	return strings.Repeat("─", left) + label + strings.Repeat("─", right)
}

func truncateCLILabel(label string, maxWidth int) string {
	runes := []rune(label)
	if len(runes) <= maxWidth {
		return label
	}
	if maxWidth <= 1 {
		return "…"
	}
	return string(runes[:maxWidth-1]) + "…"
}

// cliSection writes a themed Unicode section rule for interactive CLI output.
func cliSection(label string) {
	outln(colorizeOutput(cliSectionRule(label), currentLogTheme.Component))
}

func resolveCLITheme(cfg *models.Config, configPaths []string) logging.Theme {
	palette, ok := resolveCLIPalette(cfg, configPaths)
	if !ok {
		return logging.DefaultTheme()
	}
	return logging.ThemeFromPalette(palette)
}

func resolveCLIPalette(cfg *models.Config, configPaths []string) (*palettes.Palette, bool) {
	if cfg == nil {
		return nil, false
	}
	extra := make(map[string]any)
	if len(configPaths) > 0 {
		extra["config_path"] = configPaths[0]
	}
	return loadLoggerPalette(cfg.Theme, extra)
}

func resolveCLIThemeFromConfig() logging.Theme {
	return logging.ThemeFromPalette(resolveCLIPaletteFromConfig())
}

func resolveCLIPaletteFromConfig() *palettes.Palette {
	cfg, _, configPaths, err := loadManagerConfig(cfgFile)
	if err != nil {
		return nil
	}
	palette, _ := resolveCLIPalette(cfg, configPaths)
	return palette
}

func configureCLIThemeFromConfig() error {
	currentCLIPalette = resolveCLIPaletteFromConfig()
	return configureCommandLogger(logging.ThemeFromPalette(currentCLIPalette))
}

//nolint:gocyclo // Handles distinct help sections and narrow-terminal variants.
func renderCommandHelp(cmd *cobra.Command) {
	currentCmd = cmd
	theme := resolveCLIThemeFromConfig()
	writer := cmd.OutOrStdout()

	if cmd.Long != "" {
		renderHelpDescription(writer, cmd.Long, theme)
	} else if cmd.Short != "" {
		_, _ = fmt.Fprintln(writer, colorizeOutput(cmd.Short, theme.Component))
	}
	_, _ = fmt.Fprintln(writer)
	_, _ = fmt.Fprintln(writer, colorizeOutput("Usage:", theme.Component))
	_, _ = fmt.Fprintf(writer, "  %s\n", cmd.UseLine())
	if len(cmd.Aliases) > 0 {
		_, _ = fmt.Fprintln(writer)
		_, _ = fmt.Fprintln(writer, colorizeOutput("Aliases:", theme.Component))
		_, _ = fmt.Fprintf(writer, "  %s, %s\n", cmd.Name(), strings.Join(cmd.Aliases, ", "))
	}

	if cmd.HasAvailableSubCommands() {
		_, _ = fmt.Fprintln(writer)
		_, _ = fmt.Fprintln(writer, colorizeOutput("Commands:", theme.Component))
		width := min(cliTerminalWidth(writer), cliPreferredRuleWidth)
		for _, child := range cmd.Commands() {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			aliases := strings.Join(child.Aliases, ", ")
			name := child.Name()
			if width < 48 {
				label := name
				if aliases != "" {
					label = aliases + " " + name
				}
				_, _ = fmt.Fprintf(writer, "  %s\n", colorizeOutput(label, theme.Success))
				for _, line := range wrapCLIText(child.Short, width-4) {
					_, _ = fmt.Fprintf(writer, "    %s\n", line)
				}
			} else {
				if aliases != "" {
					aliases += "  "
				}
				_, _ = fmt.Fprintf(writer, "  %-8s %-10s %s\n", colorizeOutput(aliases, theme.Component), colorizeOutput(name, theme.Success), child.Short)
			}
		}
	}

	if cmd.HasAvailableLocalFlags() || cmd.HasAvailableInheritedFlags() {
		helpWidth := min(cliTerminalWidth(writer), cliPreferredRuleWidth)
		_, _ = fmt.Fprintln(writer)
		_, _ = fmt.Fprintln(writer, colorizeOutput("Flags:", theme.Component))
		_, _ = fmt.Fprint(writer, cmd.LocalFlags().FlagUsagesWrapped(helpWidth))
		if cmd.HasAvailableInheritedFlags() {
			_, _ = fmt.Fprintln(writer)
			_, _ = fmt.Fprintln(writer, colorizeOutput("Global Flags:", theme.Component))
			_, _ = fmt.Fprint(writer, cmd.InheritedFlags().FlagUsagesWrapped(helpWidth))
		}
	}

	var helpTopics []*cobra.Command
	for _, child := range cmd.Commands() {
		if child.IsAdditionalHelpTopicCommand() {
			helpTopics = append(helpTopics, child)
		}
	}
	if len(helpTopics) > 0 {
		_, _ = fmt.Fprintln(writer)
		_, _ = fmt.Fprintln(writer, colorizeOutput("Additional help topics:", theme.Component))
		for _, topic := range helpTopics {
			if min(cliTerminalWidth(writer), cliPreferredRuleWidth) < 48 {
				_, _ = fmt.Fprintf(writer, "  %s\n", colorizeOutput(topic.Name(), theme.Success))
				for _, line := range wrapCLIText(topic.Short, min(cliTerminalWidth(writer), cliPreferredRuleWidth)-4) {
					_, _ = fmt.Fprintf(writer, "    %s\n", line)
				}
			} else {
				_, _ = fmt.Fprintf(writer, "  %-16s %s\n", colorizeOutput(topic.Name(), theme.Success), topic.Short)
			}
		}
	}

	if cmd.HasAvailableSubCommands() {
		_, _ = fmt.Fprintln(writer)
		footer := "Run '" + cmd.CommandPath() + " <command> --help' for detailed command help."
		if min(cliTerminalWidth(writer), cliPreferredRuleWidth) < 48 {
			for _, line := range wrapCLIText(footer, min(cliTerminalWidth(writer), cliPreferredRuleWidth)) {
				_, _ = fmt.Fprintln(writer, line)
			}
		} else {
			_, _ = fmt.Fprintln(writer, footer)
		}
	}
}

func wrapCLIText(value string, width int) []string {
	if width < 1 {
		return []string{value}
	}
	words := strings.Fields(value)
	var lines []string
	line := ""
	for _, word := range words {
		if line != "" && utf8.RuneCountInString(line)+1+utf8.RuneCountInString(word) > width {
			lines = append(lines, line)
			line = ""
		}
		if line != "" {
			line += " "
		}
		line += word
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

func renderHelpDescription(writer io.Writer, description string, theme logging.Theme) {
	width := min(cliTerminalWidth(writer), cliPreferredRuleWidth)
	for _, line := range strings.Split(description, "\n") {
		if isHelpSectionHeading(line) {
			heading := strings.TrimSuffix(strings.TrimSpace(line), ":")
			_, _ = fmt.Fprintln(writer, colorizeOutput(cliSectionRule(heading), theme.Component))
			continue
		}
		if width < 48 && strings.TrimSpace(line) != "" {
			trimmed := strings.TrimLeft(line, " \t")
			indent := line[:len(line)-len(trimmed)]
			for _, wrapped := range wrapCLIText(trimmed, width-utf8.RuneCountInString(indent)) {
				_, _ = fmt.Fprintln(writer, indent+wrapped)
			}
			continue
		}
		_, _ = fmt.Fprintln(writer, line)
	}
}

func isHelpSectionHeading(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasSuffix(line, ":") || len(line) < 3 || len(line) > 40 {
		return false
	}
	return !strings.Contains(line, "://")
}
