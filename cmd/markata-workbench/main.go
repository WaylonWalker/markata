// Command markata-workbench serves the standalone Go Markdown editor.
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WaylonWalker/markata/pkg/workbench"

	"github.com/spf13/cobra"
)

func main() {
	var address string
	command := &cobra.Command{Use: "markata-workbench", Short: "Serve a local Markdown editor with in-process Go rendering", Version: "dev", Args: cobra.NoArgs, SilenceUsage: true, RunE: func(cmd *cobra.Command, _ []string) error {
		fmt.Fprintln(cmd.ErrOrStderr(), "Preparing Markata appearance assets…")
		handler, err := workbench.New()
		if err != nil {
			return err
		}
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return fmt.Errorf("listen on %s: %w", address, err)
		}
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		fmt.Fprintf(cmd.ErrOrStderr(), "Markata Go workbench: http://%s/\nStop with Ctrl+C.\n", listener.Addr())
		return serve(ctx, listener, handler, cmd.ErrOrStderr())
	}}
	command.Flags().StringVar(&address, "listen", "127.0.0.1:8766", "HTTP listen address (trusted local use)")
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

// serve waits for active requests to drain before returning after cancellation.
func serve(ctx context.Context, listener net.Listener, handler http.Handler, output io.Writer) error {
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 16 << 10}
	shutdownDone := make(chan error, 1)
	serveDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := server.Shutdown(shutdown)
			if err != nil {
				if closeErr := server.Close(); closeErr != nil {
					fmt.Fprintf(output, "Close: %v\n", closeErr)
				}
			}
			shutdownDone <- err
		case <-serveDone:
			shutdownDone <- nil
		}
	}()
	err := server.Serve(listener)
	close(serveDone)
	shutdownErr := <-shutdownDone
	if errors.Is(err, http.ErrServerClosed) {
		return shutdownErr
	}
	return err
}
