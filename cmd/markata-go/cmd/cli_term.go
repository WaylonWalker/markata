package cmd

import (
	"os"

	"github.com/WaylonWalker/markata/pkg/logging"
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
	return logging.ColorEnabled(isTTY, logging.Format(logFormat), forceColor, noColor)
}
