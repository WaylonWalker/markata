package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/spf13/cobra"
)

var currentCmd *cobra.Command

func activeCmd() *cobra.Command {
	return currentCmd
}

func outWriter() io.Writer {
	if cmd := activeCmd(); cmd != nil {
		return cmd.OutOrStdout()
	}
	return os.Stdout
}

func errWriter() io.Writer {
	if cmd := activeCmd(); cmd != nil {
		return cmd.ErrOrStderr()
	}
	return os.Stderr
}

func inReader() io.Reader {
	if cmd := activeCmd(); cmd != nil {
		return cmd.InOrStdin()
	}
	return os.Stdin
}

func outln(args ...any) {
	_, _ = fmt.Fprintln(outWriter(), args...)
}

func outlnf(format string, args ...any) {
	_, _ = fmt.Fprintf(outWriter(), format+"\n", args...)
}

func outText(text string) {
	_, _ = fmt.Fprint(outWriter(), text)
}

func out(format string, args ...any) {
	_, _ = fmt.Fprintf(outWriter(), format, args...)
}

func errln(args ...any) {
	_, _ = fmt.Fprintln(errWriter(), args...)
}

func errlnf(format string, args ...any) {
	_, _ = fmt.Fprintf(errWriter(), format+"\n", args...)
}

func errf(format string, args ...any) {
	_, _ = fmt.Fprintf(errWriter(), format, args...)
}

func cliLogger() logging.Logger {
	if cmd := activeCmd(); cmd != nil {
		return logging.Component(cmd.Name())
	}
	return logging.Component("cli")
}

func infof(format string, args ...any) {
	if quiet {
		return
	}
	cliLogger().Infof(format, args...)
}

func verbosef(format string, args ...any) {
	if !verbose || quiet {
		return
	}
	cliLogger().Phase("cli").Level("debug").Printf(format, args...)
}

func warnf(format string, args ...any) {
	if quiet {
		return
	}
	cliLogger().Warnf(format, args...)
}
