// Package main provides the canonical Markata CLI entry point.
package main

import (
	"fmt"
	"os"

	"github.com/WaylonWalker/markata/cmd/markata-go/cmd"
)

func main() {
	if err := cmd.ExecuteAs("markata"); err != nil {
		exitCode := cmd.ExitCodeForError(err)
		if message := cmd.FormatError(err); message != "" {
			fmt.Fprintln(os.Stderr, message)
		}
		os.Exit(exitCode)
	}
}
