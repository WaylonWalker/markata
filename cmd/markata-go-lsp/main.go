// Package main provides the entry point for the markata-go Language Server Protocol server.
//
// The LSP server provides IDE features for markdown files with wikilink support:
//   - Autocomplete for [[wikilinks]]
//   - Diagnostics for broken wikilinks
//   - Hover information showing post title and description
//   - Go to definition for navigating to linked posts
//
// Usage:
//
//	markata-go lsp              # Start LSP server on stdin/stdout
//	markata-go-lsp              # Standalone LSP server
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/WaylonWalker/markata/pkg/lsp"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	versionFlag := flag.Bool("version", false, "print version and exit")
	verboseFlag := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("markata-go-lsp %s (%s)\n", version, commit)
		os.Exit(0)
	}

	// Setup logging
	var logger *log.Logger
	if *verboseFlag {
		logger = log.New(os.Stderr, "[markata-lsp] ", log.LstdFlags|log.Lshortfile)
	} else {
		logger = log.New(os.Stderr, "[markata-lsp] ", log.LstdFlags)
	}

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		logger.Println("Shutting down...")
		cancel()
	}()

	// Create and run the LSP server
	srv := lsp.New(logger)
	if err := srv.Run(ctx, os.Stdin, os.Stdout); err != nil {
		logger.Fatalf("LSP server error: %v", err)
	}
}
