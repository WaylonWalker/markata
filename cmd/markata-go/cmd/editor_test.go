package cmd

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestValidateEditorListenAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{name: "IPv4 loopback", address: "127.0.0.1:8766"},
		{name: "IPv6 loopback", address: "[::1]:8766"},
		{name: "localhost", address: "localhost:8766"},
		{name: "all interfaces", address: "0.0.0.0:8766", wantErr: true},
		{name: "non-loopback IP", address: "192.0.2.1:8766", wantErr: true},
		{name: "hostname", address: "editor.example:8766", wantErr: true},
		{name: "missing port", address: "127.0.0.1", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEditorListenAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateEditorListenAddress(%q) error = %v, wantErr %t", tt.address, err, tt.wantErr)
			}
		})
	}
}

func TestEditorCommand_ServesAndShutsDown(t *testing.T) {
	previousWorkspace, previousListen := editorWorkspace, editorListen
	t.Cleanup(func() {
		editorWorkspace, editorListen = previousWorkspace, previousListen
	})
	editorWorkspace = t.TempDir()
	editorListen = "127.0.0.1:0"

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	cli := &cobra.Command{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cli.SetContext(ctx)
	cli.SetErr(writer)
	done := make(chan error, 1)
	go func() { done <- runEditorCommand(cli, nil) }()

	lines := bufio.NewReader(reader)
	readLine := func() string {
		t.Helper()
		type lineResult struct {
			value string
			err   error
		}
		line := make(chan lineResult, 1)
		go func() {
			value, err := lines.ReadString('\n')
			line <- lineResult{value: value, err: err}
		}()
		select {
		case result := <-line:
			if result.err != nil {
				t.Fatalf("read editor startup output: %v", result.err)
			}
			return strings.TrimSpace(result.value)
		case <-time.After(10 * time.Second):
			t.Fatal("editor did not print startup output")
			return ""
		}
	}
	if got := readLine(); got != "Preparing Markata editor assets…" {
		t.Fatalf("startup status = %q", got)
	}
	urlLine := readLine()
	urlText, ok := strings.CutPrefix(urlLine, "Markata editor: ")
	if !ok {
		t.Fatalf("startup URL line = %q", urlLine)
	}
	if got := readLine(); got != "Stop with Ctrl+C." {
		t.Fatalf("stop hint = %q", got)
	}
	editorURL, err := url.Parse(urlText)
	if err != nil || editorURL.Host == "" {
		t.Fatalf("invalid editor URL %q: %v", urlText, err)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, editorURL.String(), http.NoBody)
	if err != nil {
		t.Fatalf("create editor request: %v", err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("GET editor URL: %v", err)
	}
	_ = response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET editor status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	cancel()
	select {
	case err = <-done:
		if err != nil {
			t.Fatalf("editor shutdown error: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("editor did not stop after context cancellation")
	}
}

func TestEditorCommand_RejectsNonLoopbackBeforeStarting(t *testing.T) {
	previousListen := editorListen
	editorListen = "0.0.0.0:8766"
	t.Cleanup(func() { editorListen = previousListen })
	if err := runEditorCommand(&cobra.Command{}, nil); err == nil || !strings.Contains(err.Error(), "loopback") {
		t.Fatalf("runEditorCommand() error = %v, want loopback restriction", err)
	}
}
