package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
)

const (
	dropperUploaderPath    = "/v1/upload"
	dropperUploaderMaxBody = 15 << 20
	dropperUploaderMaxResp = 1 << 20
)

// DropperUploader provides image uploads without giving the workbench process
// access to Dropper credentials.
type DropperUploader interface {
	Upload(context.Context, string, string, []byte) (string, error)
}

type unixDropperUploader struct {
	socket string
	client *http.Client
}

// NewUnixDropperUploader connects the workbench to a credential-bearing
// Dropper uploader over a Unix domain socket.
func NewUnixDropperUploader(socket string) DropperUploader {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", socket)
		},
	}
	return &unixDropperUploader{
		socket: socket,
		client: &http.Client{
			Transport: transport,
			Timeout:   35 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (u *unixDropperUploader) Upload(ctx context.Context, filename, contentType string, data []byte) (string, error) {
	if strings.TrimSpace(u.socket) == "" {
		return "", errors.New("dropper uploader socket is not configured")
	}
	if len(data) == 0 || len(data) > dropperUploaderMaxBody {
		return "", fmt.Errorf("dropper image size must be between 1 byte and %d MiB", dropperUploaderMaxBody>>20)
	}
	if err := validateDropperUploadFilename(filename); err != nil {
		return "", err
	}
	if _, err := editorworkspace.ImageExtension(contentType, data); err != nil {
		return "", err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://dropper-uploader"+dropperUploaderPath, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("create local Dropper upload request: %w", err)
	}
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("X-Markata-Filename", filename)
	response, err := u.client.Do(request)
	if err != nil {
		return "", fmt.Errorf("contact local Dropper uploader: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, dropperUploaderMaxResp+1))
	if err != nil {
		return "", fmt.Errorf("read local Dropper uploader response: %w", err)
	}
	if len(body) > dropperUploaderMaxResp {
		return "", errors.New("local Dropper uploader response exceeded its limit")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var failure struct {
			Error string `json:"error"`
		}
		message := ""
		if err := json.Unmarshal(body, &failure); err == nil {
			message = strings.TrimSpace(failure.Error)
		}
		if message == "" {
			message = http.StatusText(response.StatusCode)
		}
		return "", fmt.Errorf("local Dropper uploader rejected request (%d): %s", response.StatusCode, message)
	}
	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("decode local Dropper uploader response: %w", err)
	}
	imageURL, err := validatedDropperImageURL(result.URL)
	if err != nil {
		return "", fmt.Errorf("local Dropper uploader returned an invalid image URL")
	}
	return imageURL, nil
}

// DropperUploaderServerConfig fixes the credential-bearing uploader to one
// Dropper upload endpoint and token.
type DropperUploaderServerConfig struct {
	SocketPath string
	Endpoint   string
	Token      string
}

type dropperUploaderService struct {
	endpoint string
	token    string
}

func newDropperUploaderService(config DropperUploaderServerConfig) (*dropperUploaderService, error) {
	if strings.TrimSpace(config.Token) == "" {
		return nil, errors.New("dropper uploader token is required")
	}
	endpoint := strings.TrimSpace(config.Endpoint)
	if endpoint == "" {
		endpoint = "https://dropper.waylonwalker.com/api/upload?format=json"
	}
	if err := validateDropperUploadEndpoint(endpoint); err != nil {
		return nil, err
	}
	return &dropperUploaderService{endpoint: endpoint, token: strings.TrimSpace(config.Token)}, nil
}

// ServeDropperUploader serves a fixed image-upload capability on a Unix
// socket. The Dropper token never enters the workbench container or browser.
func ServeDropperUploader(ctx context.Context, config DropperUploaderServerConfig) error {
	service, err := newDropperUploaderService(config)
	if err != nil {
		return err
	}
	if strings.TrimSpace(config.SocketPath) == "" {
		return errors.New("dropper uploader socket path is required")
	}
	socket, err := filepath.Abs(config.SocketPath)
	if err != nil {
		return fmt.Errorf("resolve Dropper uploader socket: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(socket), 0o770); err != nil {
		return fmt.Errorf("create Dropper uploader socket directory: %w", err)
	}
	if info, err := os.Lstat(socket); err == nil {
		if info.Mode()&os.ModeSocket == 0 {
			return fmt.Errorf("refusing to replace non-socket path %s", socket)
		}
		if err := os.Remove(socket); err != nil {
			return fmt.Errorf("remove stale Dropper uploader socket: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("inspect Dropper uploader socket: %w", err)
	}
	listener, err := (&net.ListenConfig{}).Listen(ctx, "unix", socket)
	if err != nil {
		return fmt.Errorf("listen on Dropper uploader socket: %w", err)
	}
	defer func() {
		_ = listener.Close()
		_ = os.Remove(socket)
	}()
	if err := os.Chmod(socket, 0o660); err != nil {
		return fmt.Errorf("set Dropper uploader socket permissions: %w", err)
	}
	server := &http.Server{
		Handler:           service.handler(),
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    8192,
	}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Dropper uploader shutdown: %v", err)
		}
	}()
	err = server.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) || errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}

func (s *dropperUploaderService) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(dropperUploaderPath, s.handleUpload)
	return mux
}

func (s *dropperUploaderService) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		dropperUploaderError(w, http.StatusMethodNotAllowed, "POST is required")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, dropperUploaderMaxBody+1)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			dropperUploaderError(w, http.StatusRequestEntityTooLarge, "image exceeds 15 MiB")
			return
		}
		dropperUploaderError(w, http.StatusBadRequest, "cannot read image")
		return
	}
	if len(data) == 0 || len(data) > dropperUploaderMaxBody {
		dropperUploaderError(w, http.StatusRequestEntityTooLarge, "image must be between 1 byte and 15 MiB")
		return
	}
	filename := r.Header.Get("X-Markata-Filename")
	if err := validateDropperUploadFilename(filename); err != nil {
		dropperUploaderError(w, http.StatusBadRequest, err.Error())
		return
	}
	contentType := strings.TrimSpace(r.Header.Get("Content-Type"))
	if _, err := editorworkspace.ImageExtension(contentType, data); err != nil {
		dropperUploaderError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	imageURL, err := uploadImageToDropper(r.Context(), s.endpoint, s.token, filename, contentType, data)
	if err != nil {
		dropperUploaderError(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"url": imageURL}); err != nil {
		return
	}
}

func uploadImageToDropper(ctx context.Context, endpoint, token, filename, contentType string, data []byte) (string, error) {
	if err := validateDropperUploadEndpoint(endpoint); err != nil {
		return "", err
	}
	payload, formContentType, err := dropperUploadPayload(filename, contentType, data)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, &payload)
	if err != nil {
		return "", fmt.Errorf("create Dropper request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Content-Type", formContentType)
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if len(via) > 0 && !strings.EqualFold(request.URL.Host, via[0].URL.Host) {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("dropper upload request failed: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("read Dropper response: %w", err)
	}
	return dropperUploadResponseURL(response.StatusCode, response.Status, responseBody)
}

func validateDropperUploadFilename(filename string) error {
	if filename == "" || filename != filepath.Base(filename) || strings.ContainsAny(filename, "\\/\r\n\x00") {
		return errors.New("dropper image filename must be a plain file name")
	}
	return nil
}

func validatedDropperImageURL(raw string) (string, error) {
	imageURL, err := url.ParseRequestURI(raw)
	if err != nil || imageURL.Host == "" || imageURL.User != nil || imageURL.Scheme != secureScheme && !(imageURL.Scheme == insecureScheme && isLoopbackHost(imageURL.Hostname())) {
		return "", errors.New("dropper returned an invalid image URL")
	}
	return imageURL.String(), nil
}

func dropperUploaderError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("encode Dropper uploader error response: %v", err)
	}
}
