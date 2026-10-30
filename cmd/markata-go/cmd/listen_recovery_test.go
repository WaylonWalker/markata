package cmd

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestListenWithPortRecovery_NonInteractive(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer occupied.Close()
	command := &cobra.Command{Use: "editor"}
	command.SetIn(strings.NewReader(""))
	command.SetErr(&bytes.Buffer{})
	listener, err := listenWithPortRecovery(command, occupied.Addr().String())
	if listener != nil {
		_ = listener.Close()
		t.Fatal("unexpected listener")
	}
	if err == nil || !strings.Contains(err.Error(), "editor --port ") {
		t.Fatalf("expected rerun instruction, got %v", err)
	}
}
