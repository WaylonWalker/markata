package cmd

import (
	"fmt"

	"github.com/WaylonWalker/markata/pkg/palettes"
)

const ansiReset = "\033[0m"

func colorizeOutput(text, color string) string {
	return colorize(text, color, colorEnabledOnOutput())
}

func colorize(text, color string, enabled bool) string {
	if !enabled {
		return text
	}
	if rgb := ansiTrueColor(color); rgb != "" {
		return rgb + text + ansiReset
	}
	return color + text + ansiReset
}

func ansiTrueColor(hex string) string {
	if hex == "" || hex[0] != '#' {
		return ""
	}
	parsed, err := palettes.ParseHexColor(hex)
	if err != nil {
		return ""
	}
	r, g, b, _ := parsed.RGBA()
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)
}
