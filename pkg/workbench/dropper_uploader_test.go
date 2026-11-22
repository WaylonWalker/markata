package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func newDropperTestSocketPath(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp(os.TempDir(), "m-")
	if err != nil {
		t.Fatalf("create short socket test directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Errorf("remove socket test directory: %v", err)
		}
	})
	return filepath.Join(dir, "u.sock")
}

func TestDropperUploaderService_ValidatesAndForwardsImage(t *testing.T) {
	var gotAuthorization, gotFilename, gotContentType string
	var gotImage []byte
	dropper := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuthorization = r.Header.Get("Authorization")
		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			t.Errorf("parse forwarded content type: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		reader := multipart.NewReader(r.Body, params["boundary"])
		part, err := reader.NextPart()
		if err != nil {
			t.Errorf("read forwarded image: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		gotFilename = part.FileName()
		gotContentType = part.Header.Get("Content-Type")
		gotImage, err = io.ReadAll(part)
		if err != nil {
			t.Errorf("read forwarded image data: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"url":"https://dropper.example.com/file/uploaded.png"}`); err != nil {
			t.Errorf("write mock Dropper response: %v", err)
		}
	}))
	defer dropper.Close()

	service, err := newDropperUploaderService(DropperUploaderServerConfig{
		Endpoint: dropper.URL + "/api/upload?format=json",
		Token:    "test-dropper-token",
	})
	if err != nil {
		t.Fatalf("create uploader service: %v", err)
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, dropperUploaderPath, strings.NewReader(string(testPastePNG)))
	request.Header.Set("Content-Type", "image/png")
	request.Header.Set("X-Markata-Filename", "pasted.png")
	service.handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("upload returned %d: %s", response.Code, response.Body.String())
	}
	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode upload response: %v", err)
	}
	if payload["url"] != "https://dropper.example.com/file/uploaded.png" {
		t.Fatalf("upload URL = %q", payload["url"])
	}
	if gotAuthorization != "Bearer test-dropper-token" {
		t.Fatalf("Dropper authorization = %q", gotAuthorization)
	}
	if gotFilename != "pasted.png" || gotContentType != "image/png" || !bytes.Equal(gotImage, testPastePNG) {
		t.Fatalf("forwarded file = %q, %q, %x", gotFilename, gotContentType, gotImage)
	}
	if strings.Contains(response.Body.String(), "test-dropper-token") {
		t.Fatal("Dropper token leaked from uploader response")
	}
}

func TestUnixDropperUploader_UsesSocketAndValidatesReturnedURL(t *testing.T) {
	dropper := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"url":"javascript:alert(1)"}`); err != nil {
			t.Errorf("write mock Dropper response: %v", err)
		}
	}))
	defer dropper.Close()
	service, err := newDropperUploaderService(DropperUploaderServerConfig{
		Endpoint: dropper.URL + "/api/upload?format=json",
		Token:    "test-dropper-token",
	})
	if err != nil {
		t.Fatalf("create uploader service: %v", err)
	}
	socket := newDropperTestSocketPath(t)
	listener, err := (&net.ListenConfig{}).Listen(context.Background(), "unix", socket)
	if err != nil {
		t.Fatalf("listen on test socket: %v", err)
	}
	if err := os.Chmod(socket, 0o660); err != nil {
		t.Fatalf("chmod test socket: %v", err)
	}
	server := &http.Server{Handler: service.handler(), ReadHeaderTimeout: 5 * time.Second}
	defer server.Close()
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("serve test Dropper uploader: %v", err)
		}
	}()

	uploader := NewUnixDropperUploader(socket)
	if _, err := uploader.Upload(context.Background(), "pasted.png", "image/png", testPastePNG); err == nil {
		t.Fatal("uploader accepted an invalid response URL")
	}
}

func TestDropperUploaderService_RejectsInvalidRequests(t *testing.T) {
	dropperCalls := 0
	dropper := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		dropperCalls++
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"url":"https://dropper.example.com/file/uploaded.png"}`); err != nil {
			t.Errorf("write mock Dropper response: %v", err)
		}
	}))
	defer dropper.Close()
	service, err := newDropperUploaderService(DropperUploaderServerConfig{Endpoint: dropper.URL, Token: "test-token"})
	if err != nil {
		t.Fatalf("create uploader service: %v", err)
	}
	for _, tc := range []struct {
		name        string
		method      string
		path        string
		contentType string
		filename    string
		body        []byte
		wantStatus  int
	}{
		{name: "wrong method", method: http.MethodGet, path: dropperUploaderPath, wantStatus: http.StatusMethodNotAllowed},
		{name: "wrong path", method: http.MethodPost, path: "/v1/admin", body: testPastePNG, wantStatus: http.StatusNotFound},
		{name: "path traversal filename", method: http.MethodPost, path: dropperUploaderPath, contentType: "image/png", filename: "../pasted.png", body: testPastePNG, wantStatus: http.StatusBadRequest},
		{name: "mismatched image type", method: http.MethodPost, path: dropperUploaderPath, contentType: "image/jpeg", filename: "pasted.jpg", body: testPastePNG, wantStatus: http.StatusUnprocessableEntity},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequestWithContext(context.Background(), tc.method, tc.path, strings.NewReader(string(tc.body)))
			if tc.contentType != "" {
				request.Header.Set("Content-Type", tc.contentType)
			}
			if tc.filename != "" {
				request.Header.Set("X-Markata-Filename", tc.filename)
			}
			service.handler().ServeHTTP(response, request)
			if response.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.wantStatus, response.Body.String())
			}
		})
	}
	if dropperCalls != 0 {
		t.Fatalf("invalid requests reached Dropper %d times", dropperCalls)
	}
}

func TestNewDropperUploaderService_RequiresCredentialAndSafeEndpoint(t *testing.T) {
	if _, err := newDropperUploaderService(DropperUploaderServerConfig{Endpoint: "https://dropper.example.com/api/upload"}); err == nil {
		t.Fatal("uploader accepted a missing token")
	}
	if _, err := newDropperUploaderService(DropperUploaderServerConfig{Endpoint: "http://example.com/api/upload", Token: "test-token"}); err == nil {
		t.Fatal("uploader accepted a non-loopback HTTP endpoint")
	}
}

func TestDropperUploader_DoesNotForwardUpstreamBodyOrCredential(t *testing.T) {
	const secret = "test-dropper-token"
	dropper := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := io.WriteString(w, "bad credential: "+secret); err != nil {
			t.Errorf("write mock Dropper error: %v", err)
		}
	}))
	defer dropper.Close()
	service, err := newDropperUploaderService(DropperUploaderServerConfig{
		Endpoint: dropper.URL + "/api/upload?format=json",
		Token:    secret,
	})
	if err != nil {
		t.Fatalf("create uploader service: %v", err)
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, dropperUploaderPath, strings.NewReader(string(testPastePNG)))
	request.Header.Set("Content-Type", "image/png")
	request.Header.Set("X-Markata-Filename", "pasted.png")
	service.handler().ServeHTTP(response, request)
	if response.Code != http.StatusBadGateway {
		t.Fatalf("upload returned %d, want %d: %s", response.Code, http.StatusBadGateway, response.Body.String())
	}
	if strings.Contains(response.Body.String(), secret) || strings.Contains(response.Body.String(), "bad credential") {
		t.Fatalf("uploader leaked the upstream response or credential: %s", response.Body.String())
	}
}

func TestDropperUploaderSocketPermissions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	socket := newDropperTestSocketPath(t)
	started := make(chan error, 1)
	go func() {
		started <- ServeDropperUploader(ctx, DropperUploaderServerConfig{
			SocketPath: socket,
			Endpoint:   "https://dropper.example.com/api/upload?format=json",
			Token:      "test-token",
		})
	}()
	for i := 0; i < 100; i++ {
		if info, err := os.Stat(socket); err == nil && info.Mode()&os.ModeSocket != 0 {
			if info.Mode().Perm() != 0o660 {
				t.Fatalf("socket permissions = %o, want 660", info.Mode().Perm())
			}
			cancel()
			if err := <-started; err != nil {
				t.Fatalf("uploader server shutdown: %v", err)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	cancel()
	select {
	case err := <-started:
		t.Fatalf("uploader server did not create socket: %v", err)
	default:
		t.Fatal("uploader server did not create socket")
	}
}
