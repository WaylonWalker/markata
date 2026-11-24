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
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/WaylonWalker/markata/pkg/workbench"
	"github.com/spf13/cobra"
)

var (
	editorWorkspace             = "."
	editorPort                  int
	editorListen                = "127.0.0.1:8766"
	editorPublicOrigin          string
	editorTrustedProxyCIDRs     []string
	editorRequiredScope         string
	editorGitPublisherSocket    string
	editorDropperUploaderSocket string
)

var editorCmd = &cobra.Command{
	Use:          "editor",
	Short:        "Run the browser Markdown editor",
	Long:         "Start Markata's Go Markdown editor and preview for a workspace. It binds to loopback unless hosted Hlab ForwardAuth settings are supplied. Browser drafts stay in memory; use Download .md to save a copy.",
	Example:      "  markata editor\n  markata editor --workspace ~/sites/blog\n  markata editor --listen 127.0.0.1:8767\n  markata editor --listen 0.0.0.0:8766 --public-origin https://edit.example.com --trusted-proxy-cidr 10.42.0.0/24 --required-scope markata-editor:admin",
	Args:         cobra.NoArgs,
	SilenceUsage: true,
	RunE:         runEditorCommand,
}

func init() {
	rootCmd.AddCommand(editorCmd)
	editorCmd.Flags().IntVarP(&editorPort, "port", "p", 8766, "port to serve on")
	editorCmd.Flags().StringVar(&editorWorkspace, "workspace", ".", "Markdown workspace for language assistance")
	editorCmd.Flags().StringVar(&editorListen, "listen", "127.0.0.1:8766", "HTTP listener address (non-loopback requires hosted Hlab authentication)")
	editorCmd.Flags().StringVar(&editorPublicOrigin, "public-origin", "", "canonical HTTPS origin for hosted Hlab ForwardAuth")
	editorCmd.Flags().StringSliceVar(&editorTrustedProxyCIDRs, "trusted-proxy-cidr", nil, "trusted reverse-proxy source CIDR (repeatable)")
	editorCmd.Flags().StringVar(&editorRequiredScope, "required-scope", "", "app-specific Hlab writer scope required for every request")
	editorCmd.Flags().StringVar(&editorGitPublisherSocket, "git-publisher-socket", "", "delegate hosted Git remote checks and pushes to a local Unix-socket publisher")
	editorCmd.Flags().StringVar(&editorDropperUploaderSocket, "dropper-uploader-socket", "", "delegate hosted Dropper image uploads to a local Unix-socket uploader")
}

func runEditorCommand(cmd *cobra.Command, _ []string) error {
	if cmd.Flags().Changed("port") {
		host, _, err := net.SplitHostPort(editorListen)
		if err != nil {
			return err
		}
		editorListen = net.JoinHostPort(host, strconv.Itoa(editorPort))
	}
	hosted := editorPublicOrigin != "" || len(editorTrustedProxyCIDRs) > 0 || editorRequiredScope != ""
	var hostedConfig workbench.HostedAuthConfig
	if hosted {
		hostedConfig = workbench.HostedAuthConfig{
			PublicOrigin:      editorPublicOrigin,
			TrustedProxyCIDRs: append([]string(nil), editorTrustedProxyCIDRs...),
			RequiredScope:     editorRequiredScope,
		}
		if err := workbench.ValidateHostedAuthConfig(hostedConfig); err != nil {
			return err
		}
	}
	if err := validateEditorListenAddressWithHosted(editorListen, hosted); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if _, err := fmt.Fprintln(cmd.ErrOrStderr(), colorize("Preparing Markata editor assets…", currentLogTheme.Component, colorEnabledFor(streamIsTerminal(cmd.ErrOrStderr())))); err != nil {
		return fmt.Errorf("write editor status: %w", err)
	}
	var publisher workbench.GitPublisher
	if editorGitPublisherSocket != "" {
		publisher = workbench.NewUnixGitPublisher(editorGitPublisherSocket)
	}
	var uploader workbench.DropperUploader
	if editorDropperUploaderSocket != "" {
		uploader = workbench.NewUnixDropperUploader(editorDropperUploaderSocket)
	}
	handler, err := workbench.NewWithWorkspaceAndGitPublisherAndDropperUploader(ctx, editorWorkspace, publisher, uploader)
	if err != nil {
		return err
	}
	if hosted {
		if err := handler.EnableHostedAuth(hostedConfig); err != nil {
			_ = handler.Close()
			return err
		}
	}
	listener, err := listenWithPortRecovery(cmd, editorListen)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", editorListen, err)
	}
	if _, err = fmt.Fprintf(cmd.ErrOrStderr(), "Markata editor: %s\nStop with Ctrl+C.\n", colorize(fmt.Sprintf("http://%s/", listener.Addr()), currentLogTheme.Success, colorEnabledFor(streamIsTerminal(cmd.ErrOrStderr())))); err != nil {
		_ = listener.Close()
		return fmt.Errorf("write editor address: %w", err)
	}
	var serving http.Handler = handler
	if !hosted {
		serving = withUpIdentity(serving, upServiceEditor, editorWorkspace)
	}
	return serveEditor(ctx, listener, serving, cmd.ErrOrStderr())
}

func validateEditorListenAddress(address string) error {
	return validateEditorListenAddressWithHosted(address, false)
}

func validateEditorListenAddressWithHosted(address string, hosted bool) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid --listen address %q: %w", address, err)
	}
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return nil
	}
	if hosted && ip != nil {
		return nil
	}
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
