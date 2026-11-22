package cmd

import (
	"bytes"
	"context"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type synchronizedUpTestOutput struct {
	mu sync.Mutex
	buf bytes.Buffer
}

func (w *synchronizedUpTestOutput) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *synchronizedUpTestOutput) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}

func TestPlanUpServices(t *testing.T) {
	t.Run("both start loopback only", func(t *testing.T) {
		services, err := planUpServices(upOptions{editorPort: 8766, previewPort: 8000})
		if err != nil {
			t.Fatal(err)
		}
		if len(services) != 2 || services[0].name != upServiceEditor || services[1].name != upServiceBuilder {
			t.Fatalf("unexpected service order: %#v", services)
		}
		if !strings.Contains(strings.Join(services[0].args, " "), "127.0.0.1:8766") ||
			!strings.Contains(strings.Join(services[1].args, " "), "--host 127.0.0.1") {
			t.Fatalf("service configuration does not bind to loopback: %#v", services)
		}
	})
	t.Run("editor-only", func(t *testing.T) {
		services, err := planUpServices(upOptions{editorPort: 8766, editorOnly: true})
		if err != nil || len(services) != 1 || services[0].name != upServiceEditor {
			t.Fatalf("editor-only = %#v, %v", services, err)
		}
	})
	t.Run("builder-only", func(t *testing.T) {
		services, err := planUpServices(upOptions{previewPort: 8000, builderOnly: true})
		if err != nil || len(services) != 1 || services[0].name != upServiceBuilder {
			t.Fatalf("builder-only = %#v, %v", services, err)
		}
	})
	for _, tc := range []upOptions{
		{editorPort: 8000, previewPort: 8000},
		{editorPort: 0, previewPort: 8000},
		{editorPort: 8766, previewPort: 65536},
		{editorPort: 8766, previewPort: 8000, editorOnly: true, builderOnly: true},
	} {
		if _, err := planUpServices(tc); err == nil {
			t.Fatalf("expected invalid option rejection: %#v", tc)
		}
	}
}

func TestUpPortConflict(t *testing.T) {
	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("loopback listener unavailable: %v", err)
	}
	defer listener.Close()
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener address %T", listener.Addr())
	}
	port := strconv.Itoa(tcpAddr.Port)
	if err := checkUpPorts(context.Background(), []upService{{name: upServiceEditor, url: "http://127.0.0.1:" + port + "/"}}); err == nil {
		t.Fatal("expected occupied port to be rejected before starting children")
	}
}

func TestUpSupervisorBuilderFailureKeepsEditor(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-based subprocess fixture requires POSIX")
	}
	dir := t.TempDir()
	fixture := filepath.Join(dir, "markata-test-child")
	script := "#!/bin/sh\ncase \"$*\" in\n *serve*) exit 23 ;;\n *editor*) exec sleep 5 ;;\nesac\nexit 42\n"
	if err := os.WriteFile(fixture, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	//nolint:gosec // Test subprocess fixtures must be executable.
	if err := os.Chmod(fixture, 0o700); err != nil {
		t.Fatal(err)
	}
	services, err := planUpServices(upOptions{editorPort: 8766, previewPort: 8000})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var out, errs synchronizedUpTestOutput
	done := make(chan error, 1)
	go func() {
		done <- superviseUp(ctx, fixture, dir, []string{"--no-input"}, services, &out, &errs)
	}()
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.NewTimer(4 * time.Second)
	defer deadline.Stop()
	for {
		select {
		case err := <-done:
			t.Fatalf("workspace stopped before builder failure could be handled: %v (stderr: %s)", err, errs.String())
		case <-deadline.C:
			t.Fatalf("builder failure not reported: %q", errs.String())
		case <-ticker.C:
			if !strings.Contains(errs.String(), "Editor is still running") {
				continue
			}
			// If the supervisor is still alive after processing the builder exit,
			// the editor remains available independently of the failed builder.
			select {
			case err := <-done:
				t.Fatalf("workspace unexpectedly exited after builder stopped: %v", err)
			default:
			}
			cancel()
			if err := <-done; err != nil {
				t.Fatalf("shutdown failed: %v", err)
			}
			return
		}
	}
}

func TestUpSupervisorEditorExitIsFatal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-based subprocess fixture requires POSIX")
	}
	dir := t.TempDir()
	fixture := filepath.Join(dir, "markata-test-child")
	script := "#!/bin/sh\ncase \"$*\" in\n *editor*) exit 19 ;;\n *serve*) exec sleep 5 ;;\nesac\n"
	if err := os.WriteFile(fixture, []byte(script), 0o600); err != nil {
		t.Fatal(err)
	}
	//nolint:gosec // Test subprocess fixtures must be executable.
	if err := os.Chmod(fixture, 0o700); err != nil {
		t.Fatal(err)
	}
	services, err := planUpServices(upOptions{editorPort: 8766, previewPort: 8000})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var out, errs synchronizedUpTestOutput
	err = superviseUp(ctx, fixture, dir, []string{"--no-input"}, services, &out, &errs)
	if err == nil || !strings.Contains(err.Error(), "editor exited") {
		t.Fatalf("expected editor exit to stop the workspace, got %v", err)
	}
}
