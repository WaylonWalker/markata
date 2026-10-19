package cmd

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func inputIsTerminal() bool {
	return streamIsTerminal(inReader())
}

func outputIsTerminal() bool {
	return streamIsTerminal(outWriter())
}

func errorOutputIsTerminal() bool {
	return streamIsTerminal(errWriter())
}

func streamIsTerminal(stream any) bool {
	return fileLikeTerminal(stream)
}

func fileLikeTerminal(stream any) bool {
	file, ok := stream.(*os.File)
	if !ok {
		return false
	}
	stat, err := file.Stat()
	if err != nil {
		return false
	}
	// Character devices include /dev/null. Check that the descriptor is an actual terminal.
	return (stat.Mode()&os.ModeCharDevice) != 0 && term.IsTerminal(int(file.Fd()))
}

func colorEnabledOnOutput() bool {
	return colorEnabledFor(outputIsTerminal())
}

func colorEnabledFor(isTTY bool) bool {
	if noColor || logFormat == "plain" {
		return false
	}
	if forceColor {
		return true
	}
	if !isTTY {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return strings.TrimSpace(os.Getenv("NO_COLOR")) == ""
}
