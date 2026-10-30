package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/WaylonWalker/markata/pkg/workbench"
	"github.com/spf13/cobra"
)

var (
	editorWorkspace = "."
	editorListen    = "127.0.0.1:8766"
)

var editorCmd = &cobra.Command{
	Use:          "editor",
	Short:        "Run the local browser Markdown editor",
	Long:         "Start Markata's Go Markdown editor and preview for a workspace. Browser drafts stay in memory; use Download .md to save a copy.",
	Example:      "  markata editor\n  markata editor --workspace ~/sites/blog\n  markata editor --listen 127.0.0.1:8767",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE:         runEditorCommand,
}

func init() {
	rootCmd.AddCommand(editorCmd)
	editorCmd.Flags().StringVar(&editorWorkspace, "workspace", ".", "Markdown workspace for language assistance")
	editorCmd.Flags().StringVar(&editorListen, "listen", "127.0.0.1:8766", "loopback HTTP address for the editor")
}

func runEditorCommand(cmd *cobra.Command, _ []string) error {
	if err := validateEditorListenAddress(editorListen); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if _, err := fmt.Fprintln(cmd.ErrOrStderr(), "Preparing Markata editor assets…"); err != nil {
		return fmt.Errorf("write editor status: %w", err)
	}
	handler, err := workbench.NewWithWorkspace(ctx, editorWorkspace)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", editorListen)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", editorListen, err)
	}
	if _, err = fmt.Fprintf(cmd.ErrOrStderr(), "Markata editor: http://%s/\nStop with Ctrl+C.\n", listener.Addr()); err != nil {
		_ = listener.Close()
		return fmt.Errorf("write editor address: %w", err)
	}
	return serveEditor(ctx, listener, handler, cmd.ErrOrStderr())
}

func validateEditorListenAddress(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid --listen address %q: %w", address, err)
	}
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("--listen must use localhost or a loopback IP address; the editor serves a trusted workspace")
	}
	return nil
}

func serveEditor(ctx context.Context, listener net.Listener, handler http.Handler, output io.Writer) error {
	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
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
					_, _ = fmt.Fprintf(output, "Close: %v\n", closeErr)
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
