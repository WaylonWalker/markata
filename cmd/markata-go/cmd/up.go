package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const (
	upDefaultEditorPort  = 8766
	upDefaultPreviewPort = 8000
	upShutdownTimeout    = 8 * time.Second
)

var (
	upEditorPort  = upDefaultEditorPort
	upPreviewPort = upDefaultPreviewPort
	upEditorOnly  bool
	upBuilderOnly bool
)

// up starts the existing Go editor and site builder in independent child
// processes. This keeps their global CLI state and failure domains isolated.
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the editor and live site builder together",
	Long:  "Start the Markdown editor immediately and the live-reloading site builder alongside it. Both listen on loopback; publishing remains an explicit action. Ctrl+C stops both processes. Background service management will be added separately.",
	Args:  cobra.NoArgs,
	RunE:  runUpCommand,
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.Flags().IntVar(&upEditorPort, "editor-port", upDefaultEditorPort, "local editor port")
	upCmd.Flags().IntVar(&upPreviewPort, "preview-port", upDefaultPreviewPort, "local site preview and builder port")
	upCmd.Flags().BoolVar(&upEditorOnly, "editor-only", false, "start only the editor")
	upCmd.Flags().BoolVar(&upBuilderOnly, "builder-only", false, "start only the builder and live preview")
}

type upOptions struct {
	editorPort  int
	previewPort int
	editorOnly  bool
	builderOnly bool
}

type upService struct {
	name string
	args []string
	url  string
}

type upExit struct {
	name string
	err  error
}

func planUpServices(opts upOptions) ([]upService, error) {
	if opts.editorOnly && opts.builderOnly {
		return nil, fmt.Errorf("--editor-only and --builder-only cannot both be set")
	}
	checkPort := func(name string, port int) error {
		if port < 1 || port > 65535 {
			return fmt.Errorf("%s port %d is invalid (expected 1-65535)", name, port)
		}
		return nil
	}
	if !opts.builderOnly {
		if err := checkPort("editor", opts.editorPort); err != nil {
			return nil, err
		}
	}
	if !opts.editorOnly {
		if err := checkPort("preview", opts.previewPort); err != nil {
			return nil, err
		}
	}
	if !opts.builderOnly && !opts.editorOnly && opts.editorPort == opts.previewPort {
		return nil, fmt.Errorf("editor and preview need different ports")
	}

	services := make([]upService, 0, 2)
	if !opts.builderOnly {
		services = append(services, upService{
			name: "editor",
			args: []string{"editor", "--listen", net.JoinHostPort("127.0.0.1", strconv.Itoa(opts.editorPort))},
			url:  "http://" + net.JoinHostPort("127.0.0.1", strconv.Itoa(opts.editorPort)) + "/",
		})
	}
	if !opts.editorOnly {
		services = append(services, upService{
			name: "builder",
			args: []string{"serve", "--no-tui", "--host", "127.0.0.1", "--port", strconv.Itoa(opts.previewPort)},
			url:  "http://" + net.JoinHostPort("127.0.0.1", strconv.Itoa(opts.previewPort)) + "/",
		})
	}
	return services, nil
}

func upGlobalArgs(workspace string) ([]string, error) {
	args := []string{"--no-input"}
	absoluteConfig := func(file string) (string, error) {
		if filepath.IsAbs(file) {
			return file, nil
		}
		return filepath.Abs(filepath.Join(workspace, file))
	}
	if cfgFile != "" {
		path, err := absoluteConfig(cfgFile)
		if err != nil {
			return nil, err
		}
		args = append(args, "--config", path)
	}
	for _, file := range mergeConfigFiles {
		path, err := absoluteConfig(file)
		if err != nil {
			return nil, err
		}
		args = append(args, "--merge-config", path)
	}
	if outputDir != "" {
		args = append(args, "--output", outputDir)
	}
	if noColor {
		args = append(args, "--no-color")
	}
	return args, nil
}

func checkUpPorts(services []upService) error {
	for _, service := range services {
		host, port, err := net.SplitHostPort(service.url[len("http://") : len(service.url)-1])
		if err != nil {
			return err
		}
		listener, err := net.Listen("tcp", net.JoinHostPort(host, port))
		if err != nil {
			return fmt.Errorf("%s cannot listen on %s:%s: %w", service.name, host, port, err)
		}
		if err := listener.Close(); err != nil {
			return err
		}
	}
	return nil
}

func runUpCommand(cmd *cobra.Command, _ []string) error {
	opts := upOptions{upEditorPort, upPreviewPort, upEditorOnly, upBuilderOnly}
	services, err := planUpServices(opts)
	if err != nil {
		return err
	}
	workspace, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("resolve workspace: %w", err)
	}
	globalArgs, err := upGlobalArgs(workspace)
	if err != nil {
		return fmt.Errorf("resolve child options: %w", err)
	}
	// Prevent an accidental connection to an unrelated local listener. Child
	// commands still validate their own binding if another process races us.
	if err := checkUpPorts(services); err != nil {
		return err
	}
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate Markata executable: %w", err)
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Markata workspace: %s\n", workspace)
	for _, service := range services {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s starting: %s\n", service.name, service.url)
	}
	_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Ctrl+C stops this workspace. Save does not publish.")
	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	return superviseUp(ctx, executable, workspace, globalArgs, services, cmd.OutOrStdout(), cmd.ErrOrStderr())
}

func superviseUp(ctx context.Context, executable, workspace string, globalArgs []string, services []upService, stdout, stderr io.Writer) error {
	exits := make(chan upExit, len(services))
	running := make(map[string]*exec.Cmd, len(services))
	for _, service := range services {
		child := exec.Command(executable, append(append([]string(nil), globalArgs...), service.args...)...)
		child.Dir = workspace
		child.Stdout = stdout
		child.Stderr = stderr
		if err := child.Start(); err != nil {
			if service.name != "builder" || len(running) == 0 {
				stopUpServices(running, exits)
				return fmt.Errorf("start %s: %w", service.name, err)
			}
			_, _ = fmt.Fprintf(stderr, "Builder failed to start: %v (editor remains available)\n", err)
			continue
		}
		running[service.name] = child
		go func(name string, process *exec.Cmd) {
			exits <- upExit{name: name, err: process.Wait()}
		}(service.name, child)
	}
	for len(running) > 0 {
		select {
		case <-ctx.Done():
			stopUpServices(running, exits)
			return nil
		case result := <-exits:
			if _, ok := running[result.name]; !ok {
				continue
			}
			delete(running, result.name)
			if result.name == "builder" && len(running) > 0 {
				_, _ = fmt.Fprintf(stderr, "Builder stopped (%v). Editor is still running; inspect the builder output and restart Markata to retry.\n", result.err)
				continue
			}
			stopUpServices(running, exits)
			if result.err != nil {
				return fmt.Errorf("%s exited: %w", result.name, result.err)
			}
			return fmt.Errorf("%s exited; workspace stopped", result.name)
		}
	}
	return fmt.Errorf("no Markata services started")
}

// Stop only child processes started by this invocation. Allow the existing
// editor and builder signal handlers to finish cleanly before forcing exit.
func stopUpServices(running map[string]*exec.Cmd, exits <-chan upExit) {
	for _, child := range running {
		if err := child.Process.Signal(os.Interrupt); err != nil {
			_ = child.Process.Kill()
		}
	}
	timer := time.NewTimer(upShutdownTimeout)
	defer timer.Stop()
	for len(running) > 0 {
		select {
		case result := <-exits:
			delete(running, result.name)
		case <-timer.C:
			for _, child := range running {
				_ = child.Process.Kill()
			}
			for len(running) > 0 {
				result := <-exits
				delete(running, result.name)
			}
		}
	}
}
