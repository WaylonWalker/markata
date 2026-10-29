package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestServe_DrainsActiveRequest(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	entered, release := make(chan struct{}), make(chan struct{})
	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(entered)
		<-release
		w.WriteHeader(http.StatusNoContent)
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- serve(ctx, listener, handler, io.Discard) }()
	clientDone := make(chan error, 1)
	go func() {
		req, requestErr := http.NewRequestWithContext(context.Background(), http.MethodGet, "http://"+listener.Addr().String(), http.NoBody)
		if requestErr != nil {
			clientDone <- requestErr
			return
		}
		client := &http.Client{Timeout: 5 * time.Second}
		response, requestErr := client.Do(req)
		if requestErr == nil {
			requestErr = response.Body.Close()
		}
		clientDone <- requestErr
	}()
	<-entered
	cancel()
	select {
	case err := <-done:
		close(release)
		t.Fatalf("server exited before active request completed: %v", err)
	case <-time.After(100 * time.Millisecond):
	}
	close(release)
	if err := <-clientDone; err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not finish shutdown")
	}
}
