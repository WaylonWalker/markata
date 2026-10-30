package cmd

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// listenWithPortRecovery retains the offered socket while the user decides.
func listenWithPortRecovery(cmd *cobra.Command, address string) (net.Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err == nil || !isAddressInUseError(err) {
		return listener, err
	}
	host, portText, splitErr := net.SplitHostPort(address)
	if splitErr != nil {
		return nil, err
	}
	port, parseErr := strconv.Atoi(portText)
	if parseErr != nil {
		return nil, err
	}
	for candidate := port + 1; candidate <= min(port+20, 65535); candidate++ {
		offered, offerErr := net.Listen("tcp", net.JoinHostPort(host, strconv.Itoa(candidate)))
		if offerErr != nil {
			continue
		}
		retry := fmt.Sprintf("%s --port %d", cmd.CommandPath(), candidate)
		if noInput || !streamIsTerminal(cmd.InOrStdin()) || !streamIsTerminal(cmd.ErrOrStderr()) {
			_ = offered.Close()
			return nil, fmt.Errorf("port %d is already in use; rerun with %s", port, retry)
		}
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Port %d is already in use.\nRerun with: %s\nUse port %d now? Press Enter to continue, or q then Enter or Ctrl+C to quit: ", port, retry, candidate)
		reader := bufio.NewReader(cmd.InOrStdin())
		for {
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			type result struct {
				answer string
				err    error
			}
			input := make(chan result, 1)
			go func() { answer, err := reader.ReadString('\n'); input <- result{answer, err} }()
			var answer string
			var readErr error
			select {
			case value := <-input:
				answer, readErr = value.answer, value.err
			case <-ctx.Done():
				readErr = ctx.Err()
			}
			stop()
			if readErr != nil {
				_ = offered.Close()
				return nil, fmt.Errorf("port selection canceled: %w", readErr)
			}
			answer = strings.TrimSpace(answer)
			if answer == "" {
				return offered, nil
			}
			if strings.EqualFold(answer, "q") {
				_ = offered.Close()
				return nil, fmt.Errorf("port selection canceled")
			}
			_, _ = fmt.Fprint(cmd.ErrOrStderr(), "Press Enter to continue, or q then Enter to quit: ")
		}
	}
	return nil, err
}
