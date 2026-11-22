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
	"testing"
	"time"
)

func TestPlanUpServices(t *testing.T) {
	t.Run("both start loopback only", func(t *testing.T) {
		services, err := planUpServices(upOptions{editorPort: 8766, previewPort: 8000})
		if err != nil {
			t.Fatal(err)
		}
		if len(services) != 2 || services[0].name != "editor" || services[1].name != "builder" {
			t.Fatalf("unexpected service order: %#v", services)
		}
		if !strings.Contains(strings.Join(services[0].args, " "), "127.0.0.1:8766") ||
			!strings.Contains(strings.Join(services[1].args, " "), "--host 127.0.0.1") {
			t.Fatalf("service configuration does not bind to loopback: %#v", services)
		}
	})
	t.Run("editor-only", func(t *testing.T) {
		services, err := planUpServices(upOptions{editorPort: 8766, editorOnly: true})
		if err != nil || len(services) != 1 || services[0].name != "editor" {
			t.Fatalf("editor-only = %#v, %v", services, err)
		}
	})
	t.Run("builder-only", func(t *testing.T) {
		services, err := planUpServices(upOptions{previewPort: 8000, builderOnly: true})
		if err != nil || len(services) != 1 || services[0].name != "builder" {
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
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("loopback listener unavailable: %v", err)
	}
	defer listener.Close()
	port := strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	if err := checkUpPorts([]upService{{name: "editor", url: "http://127.0.0.1:" + port + "/"}}); err == nil {
		t.Fatal("expected occupied port to be rejected before starting children")
	}
}

func TestUpSupervisorBuilderFailureKeepsEditor(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-based subprocess fixture requires POSIX")
	}
	dir := t.TempDir()
	fixture := filepath.Join(dir, "markata-test-child")
	// The builder exits immediately; the editor must remain running until
	// context cancellation, without turning the builder failure into success.
	script := "#!/bin/sh\ncase \"$*\" in\n *serve*) exit 23 ;;\n *editor*) exec sleep 5 ;;\nesac\nexit 42\n"
	if err := os.WriteFile(fixture, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	services, err := planUpServices(upOptions{editorPort: 8766, previewPort: 8000})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()
	var out, errs bytes.Buffer
	start := time.Now()
	if err := superviseUp(ctx, fixture, dir, []string{"--no-input"}, services, &out, &errs); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) < 200*time.Millisecond {
		t.Fatal("supervisor exited as soon as the builder failed")
	}
	if !strings.Contains(errs.String(), "Editor is still running") {
		t.Fatalf("builder failure not reported: %q", errs.String())
	}
}

func TestUpSupervisorEditorExitIsFatal(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-based subprocess fixture requires POSIX")
	}
	dir := t.TempDir()
	fixture := filepath.Join(dir, "markata-test-child")
	script := "#!/bin/sh\ncase \"$*\" in\n *editor*) exit 19 ;;\n *serve*) exec sleep 5 ;;\nesac\n"
	if err := os.WriteFile(fixture, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	services, err := planUpServices(upOptions{editorPort: 8766, previewPort: 8000})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var out, errs bytes.Buffer
	err = superviseUp(ctx, fixture, dir, []string{"--no-input"}, services, &out, &errs)
	if err == nil || !strings.Contains(err.Error(), "editor exited") {
		t.Fatalf("expected editor exit to stop the workspace, got %v", err)
	}
}
