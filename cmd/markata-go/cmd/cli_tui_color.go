package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/WaylonWalker/markata/pkg/logging"
)

// configureTUIColor applies the shared CLI color policy to Lipgloss-based
// renderers, including Bubble Tea and Huh forms.
func configureTUIColor() {
	configureTUIColorFor(outputIsTerminal())
}

func configureTUIColorFor(isTTY bool) {
	format := logging.Format(logFormat)
	if !logging.ColorEnabled(isTTY, format, forceColor, noColor) {
		lipgloss.SetColorProfile(termenv.Ascii)
		return
	}

	forceFromEnv := os.Getenv("FORCE_COLOR") != "" && os.Getenv("FORCE_COLOR") != "0"
	if forceColor || forceFromEnv {
		lipgloss.SetColorProfile(termenv.TrueColor)
		return
	}

	lipgloss.SetColorProfile(termenv.EnvColorProfile())
}
