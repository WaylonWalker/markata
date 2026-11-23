package cmd

import (
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
)

func listenUpFixture(t *testing.T) (net.Listener, int) {
	t.Helper()
	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Skipf("cannot create test loopback listener: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("unexpected listener address: %T", listener.Addr())
	}
	if addr.Port >= 65515 {
		t.Skip("no room above random test port for bounded fallback")
	}
	return listener, addr.Port
}

func TestUpDefaultPortFallbackKeepsOriginalProcess(t *testing.T) {
	listener, busyPort := listenUpFixture(t)
	services, err := planUpServices(upOptions{editorPort: busyPort, editorOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	selected, changes, err := chooseUpPorts(context.Background(), services, upPortPolicy{})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 || changes[0].previous != busyPort || changes[0].selected <= busyPort {
		t.Fatalf("expected bounded fallback from %d: %#v", busyPort, changes)
	}
	next, err := upServicePort(selected[0])
	if err != nil {
		t.Fatal(err)
	}
	if next != changes[0].selected || !strings.Contains(strings.Join(selected[0].args, " "), strconv.Itoa(next)) {
		t.Fatalf("child command and URL disagree: %#v", selected[0])
	}
	// Port discovery cannot kill, reconfigure, or take over the occupant.
	if listener.Addr() == nil {
		t.Fatal("the original listener disappeared")
	}
	if err := checkUpPorts(context.Background(), services); err == nil {
		t.Fatal("original busy port unexpectedly became free")
	}
	// The caller's original service plan must be unchanged.
	original, err := upServicePort(services[0])
	if err != nil || original != busyPort {
		t.Fatalf("selection mutated original service: %#v (%v)", services[0], err)
	}
}

func TestUpOccupiedPreviewFindsNextLoopbackPort(t *testing.T) {
	_, busy := listenUpFixture(t)
	services, err := planUpServices(upOptions{previewPort: busy, builderOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	selected, changes, err := chooseUpPorts(context.Background(), services, upPortPolicy{})
	if err != nil {
		t.Fatal(err)
	}
	if len(changes) != 1 || changes[0].service != upServiceBuilder {
		t.Fatalf("expected builder port fallback: %#v", changes)
	}
	if !strings.Contains(strings.Join(selected[0].args, " "), "--port "+strconv.Itoa(changes[0].selected)) {
		t.Fatalf("builder did not receive selected --port: %#v", selected)
	}
	if !strings.HasPrefix(selected[0].url, "http://127.0.0.1:") {
		t.Fatalf("fallback exposed non-loopback address: %s", selected[0].url)
	}
}

func TestUpExplicitAndStrictPortsNeverFallback(t *testing.T) {
	_, busy := listenUpFixture(t)
	for _, tc := range []upPortPolicy{{editorExplicit: true}, {strict: true}} {
		services, err := planUpServices(upOptions{editorPort: busy, editorOnly: true})
		if err != nil {
			t.Fatal(err)
		}
		selected, changes, err := chooseUpPorts(context.Background(), services, tc)
		if err == nil || !strings.Contains(err.Error(), "--editor-port") ||
			!strings.Contains(err.Error(), "already in use") {
			t.Fatalf("wanted actionable conflict, got selected=%v changes=%v err=%v", selected, changes, err)
		}
	}
}

func TestUpExplicitPreviewPortCannotBeStolenByFallback(t *testing.T) {
	_, busy := listenUpFixture(t)
	if busy+1 > 65535 {
		t.Skip("invalid random test range")
	}
	services, err := planUpServices(upOptions{editorPort: busy, previewPort: busy + 1})
	if err != nil {
		t.Fatal(err)
	}
	selected, changes, err := chooseUpPorts(context.Background(), services, upPortPolicy{previewExplicit: true})
	if err != nil {
		t.Skipf("adjacent fixed preview port unavailable on this host: %v", err)
	}
	if len(changes) != 1 || changes[0].selected == busy+1 {
		t.Fatalf("fallback must not take explicitly reserved preview port: %#v", changes)
	}
	builderPort, err := upServicePort(selected[1])
	if err != nil || builderPort != busy+1 {
		t.Fatalf("explicit builder port was moved: %#v, %v", selected[1], err)
	}
}

func TestUpNoFallbackForUnrelatedListenFailures(t *testing.T) {
	services, err := planUpServices(upOptions{editorPort: 8766, editorOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	services[0].url = "http://0.0.0.0:8766/"
	if _, _, err := chooseUpPorts(context.Background(), services, upPortPolicy{}); err == nil {
		t.Fatal("unsafe non-loopback address must be rejected")
	}
}

func TestUpCancelledPortDiscovery(t *testing.T) {
	services, err := planUpServices(upOptions{editorPort: 8766, editorOnly: true})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := chooseUpPorts(ctx, services, upPortPolicy{}); err == nil {
		t.Fatal("expected cancellation to stop port discovery")
	}
}
