package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

var testPastePNG = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89")

func TestWorkspaceDropperUpload_ForwardsValidatedImage(t *testing.T) {
	t.Setenv("DROPPER_TOKEN", "test-dropper-token")
	var gotFilename, gotContentType, gotAuthorization string
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
		if _, err := io.WriteString(w, `{"url":"https://dropper.example.com/file/pasted.webp"}`); err != nil {
			t.Errorf("write mock Dropper response: %v", err)
		}
	}))
	defer dropper.Close()
	t.Setenv("DROPPER_API_URL", dropper.URL+"/api/upload?format=json")

	h := newHandler(t)
	h.store = nil
	h.siteConfig = &models.Config{ContentTemplates: models.ContentTemplatesConfig{ImageUpload: imageUploadModeDropper}}
	response := postWorkspaceImage(t, h, testPastePNG, "image/png", "http://localhost")
	if response.Code != http.StatusCreated {
		t.Fatalf("upload returned %d: %s", response.Code, response.Body.String())
	}
	var payload map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode upload response: %v", err)
	}
	if payload["url"] != "https://dropper.example.com/file/pasted.webp" || payload["provider"] != imageUploadModeDropper {
		t.Fatalf("unexpected upload response: %#v", payload)
	}
	if strings.Contains(response.Body.String(), "test-dropper-token") {
		t.Fatal("server-side Dropper token leaked in the response")
	}
	if gotAuthorization != "Bearer test-dropper-token" {
		t.Fatalf("Dropper authorization = %q", gotAuthorization)
	}
	if !strings.HasPrefix(gotFilename, "markata-image-") || !strings.HasSuffix(gotFilename, ".png") {
		t.Fatalf("forwarded filename = %q", gotFilename)
	}
	if gotContentType != "image/png" || !bytes.Equal(gotImage, testPastePNG) {
		t.Fatalf("forwarded image content type/data = %q, %x", gotContentType, gotImage)
	}
}

func TestWorkspaceDropperUpload_RejectsUntrustedOrInvalidRequests(t *testing.T) {
	t.Setenv("DROPPER_TOKEN", "test-dropper-token")
	var dropperCalls int
	dropper := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		dropperCalls++
		w.Header().Set("Content-Type", "application/json")
		if _, err := io.WriteString(w, `{"url":"https://dropper.example.com/file/pasted.png"}`); err != nil {
			t.Errorf("write mock Dropper response: %v", err)
		}
	}))
	defer dropper.Close()
	t.Setenv("DROPPER_API_URL", dropper.URL+"/api/upload?format=json")

	h := newHandler(t)
	h.store = nil
	h.siteConfig = &models.Config{ContentTemplates: models.ContentTemplatesConfig{ImageUpload: imageUploadModeDropper}}
	for _, tc := range []struct {
		name        string
		image       []byte
		contentType string
		origin      string
		wantStatus  int
	}{
		{name: "foreign origin", image: testPastePNG, contentType: "image/png", origin: "https://attacker.example", wantStatus: http.StatusForbidden},
		{name: "mismatched image type", image: testPastePNG, contentType: "image/jpeg", origin: "http://localhost", wantStatus: http.StatusUnprocessableEntity},
		{name: "oversized image", image: make([]byte, (15<<20)+1), contentType: "image/png", origin: "http://localhost", wantStatus: http.StatusRequestEntityTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			response := postWorkspaceImage(t, h, tc.image, tc.contentType, tc.origin)
			if response.Code != tc.wantStatus {
				t.Fatalf("upload returned %d, want %d: %s", response.Code, tc.wantStatus, response.Body.String())
			}
		})
	}
	if dropperCalls != 0 {
		t.Fatalf("invalid requests reached Dropper %d times", dropperCalls)
	}
}

func TestWorkspaceDropperUpload_RequiresDropperProvider(t *testing.T) {
	t.Setenv("DROPPER_TOKEN", "test-dropper-token")
	h := newHandler(t)
	h.store = nil
	h.siteConfig = &models.Config{ContentTemplates: models.ContentTemplatesConfig{ImageUpload: imageUploadModeLocal}}
	response := postWorkspaceImage(t, h, testPastePNG, "image/png", "http://localhost")
	if response.Code != http.StatusConflict {
		t.Fatalf("local provider upload returned %d: %s", response.Code, response.Body.String())
	}
}

func TestWorkspaceDropperUpload_UsesCredentialIsolatedUploader(t *testing.T) {
	t.Setenv("DROPPER_TOKEN", "")
	uploader := &fakeDropperUploader{url: "https://dropper.example.com/file/pasted.png"}
	h := newHandler(t)
	h.store = nil
	h.dropperUploader = uploader
	h.siteConfig = &models.Config{ContentTemplates: models.ContentTemplatesConfig{ImageUpload: imageUploadModeAuto}}

	if got := h.imageUploadProvider(); got != imageUploadModeDropper {
		t.Fatalf("auto provider = %q, want %q", got, imageUploadModeDropper)
	}
	response := postWorkspaceImage(t, h, testPastePNG, "image/png", "http://localhost")
	if response.Code != http.StatusCreated {
		t.Fatalf("upload returned %d: %s", response.Code, response.Body.String())
	}
	if uploader.calls != 1 || uploader.filename == "" || uploader.contentType != "image/png" || !bytes.Equal(uploader.data, testPastePNG) {
		t.Fatalf("isolated uploader received calls=%d filename=%q contentType=%q data=%x", uploader.calls, uploader.filename, uploader.contentType, uploader.data)
	}
	if strings.Contains(response.Body.String(), "test-dropper-token") {
		t.Fatal("Dropper token leaked in the workbench response")
	}
}

type fakeDropperUploader struct {
	calls       int
	filename    string
	contentType string
	data        []byte
	url         string
}

func (u *fakeDropperUploader) Upload(_ context.Context, filename, contentType string, data []byte) (string, error) {
	u.calls++
	u.filename = filename
	u.contentType = contentType
	u.data = append([]byte(nil), data...)
	return u.url, nil
}

func postWorkspaceImage(t *testing.T, h *Handler, image []byte, contentType, origin string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="image"; filename="clipboard-image.png"`)
	partHeader.Set("Content-Type", contentType)
	part, err := form.CreatePart(partHeader)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(image); err != nil {
		t.Fatal(err)
	}
	if err := form.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost/workspace/dropper-upload", &body)
	request.RemoteAddr = "127.0.0.1:54321"
	request.Header.Set("Origin", origin)
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Content-Type", form.FormDataContentType())
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
	return response
}

// Multi-photo Dropper cannot guarantee rollback, so fail before any remote
// upload rather than leaving orphaned images when a later step fails.
func TestNewPost_DropperMultiRejectsBeforeUploading(t *testing.T) {
	t.Setenv("DROPPER_TOKEN", "")
	uploader := &fakeDropperUploader{url: "https://dropper.example.com/file/shared.png"}
	h := newHandler(t)
	h.store = nil
	h.dropperUploader = uploader
	h.siteConfig = &models.Config{ContentTemplates: models.ContentTemplatesConfig{ImageUpload: imageUploadModeDropper}}

	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	for name, value := range map[string]string{
		"path":       "posts/shared.md",
		"source":     "---\ntitle: Shared photos\n---\n",
		"image_name": "shared",
	} {
		if err := form.WriteField(name, value); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 2; i++ {
		hdr := make(textproto.MIMEHeader)
		hdr.Set("Content-Disposition", fmt.Sprintf(`form-data; name="image"; filename="shared-%d.png"`, i))
		hdr.Set("Content-Type", "image/png")
		part, err := form.CreatePart(hdr)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write(testPastePNG); err != nil {
			t.Fatal(err)
		}
	}
	if err := form.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequestWithContext(context.Background(), http.MethodPost, "http://localhost/workspace/new-post", &body)
	request.RemoteAddr = "127.0.0.1:54321"
	request.Header.Set("Origin", "http://localhost")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Content-Type", form.FormDataContentType())
	response := httptest.NewRecorder()
	h.ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("multi-photo Dropper upload: HTTP %d: %s", response.Code, response.Body.String())
	}
	if uploader.calls != 0 {
		t.Fatalf("rejected album performed %d remote uploads", uploader.calls)
	}
}
