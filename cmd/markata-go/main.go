// Package main provides the entry point for the markata-go CLI.
package main

import (
	"fmt"
	"os"

	"github.com/WaylonWalker/markata/cmd/markata-go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		exitCode := cmd.ExitCodeForError(err)
		if message := cmd.FormatError(err); message != "" {
			fmt.Fprintln(os.Stderr, message)
		}
		os.Exit(exitCode)
	}
}
