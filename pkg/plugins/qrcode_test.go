package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestQRCodePlugin_Name(t *testing.T) {
	p := NewQRCodePlugin()
	if got := p.Name(); got != "qrcode" {
		t.Errorf("Name() = %q, want %q", got, "qrcode")
	}
}

func TestQRCodePlugin_ProcessPost_Basic(t *testing.T) {
	p := NewQRCodePlugin()
	tmpDir := t.TempDir()
	p.SetOutputDir(tmpDir)
	p.SetBaseURL("https://example.com")

	// Create QR code output directory
	qrDir := filepath.Join(tmpDir, p.Config().OutputDir)
	if err := os.MkdirAll(qrDir, 0o755); err != nil {
		t.Fatal(err)
	}

	post := &models.Post{
		Slug: "test-post",
		Href: "/test-post/",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Check QR code file was created
	qrPath := filepath.Join(qrDir, "test-post.svg")
	if _, err := os.Stat(qrPath); os.IsNotExist(err) {
		t.Error("Expected QR code file to be created")
	}

	// Check qrcode_url was set in post
	qrURL, ok := post.Extra["qrcode_url"].(string)
	if !ok {
		t.Error("Expected qrcode_url to be set in post.Extra")
	}
	if qrURL != "/qrcodes/test-post.svg" {
		t.Errorf("qrcode_url = %q, want %q", qrURL, "/qrcodes/test-post.svg")
	}
}

func TestQRCodePlugin_ProcessPost_PNG(t *testing.T) {
	p := NewQRCodePlugin()
	p.SetConfig(models.QRCodeConfig{
		Enabled:         true,
		Format:          "png",
		Size:            100,
		OutputDir:       "qrcodes",
		ErrorCorrection: "M",
		Foreground:      "#000000",
		Background:      "#ffffff",
	})

	tmpDir := t.TempDir()
	p.SetOutputDir(tmpDir)
	p.SetBaseURL("https://example.com")

	// Create QR code output directory
	qrDir := filepath.Join(tmpDir, "qrcodes")
	if err := os.MkdirAll(qrDir, 0o755); err != nil {
		t.Fatal(err)
	}

	post := &models.Post{
		Slug: "png-test",
		Href: "/png-test/",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Check PNG file was created
	qrPath := filepath.Join(qrDir, "png-test.png")
	if _, err := os.Stat(qrPath); os.IsNotExist(err) {
		t.Error("Expected PNG QR code file to be created")
	}
}

func TestQRCodePlugin_ProcessPost_SkipPost(t *testing.T) {
	p := NewQRCodePlugin()
	tmpDir := t.TempDir()
	p.SetOutputDir(tmpDir)

	post := &models.Post{
		Skip: true,
		Slug: "skip-test",
		Href: "/skip-test/",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Check QR code file was NOT created
	qrPath := filepath.Join(tmpDir, "qrcodes", "skip-test.svg")
	if _, err := os.Stat(qrPath); !os.IsNotExist(err) {
		t.Error("QR code should not be created for skipped posts")
	}
}

func TestQRCodePlugin_ProcessPost_NoHref(t *testing.T) {
	p := NewQRCodePlugin()
	tmpDir := t.TempDir()
	p.SetOutputDir(tmpDir)

	// Post with no href or slug
	post := &models.Post{
		Path: "test.md",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// No qrcode_url should be set
	if post.Extra != nil {
		if _, ok := post.Extra["qrcode_url"]; ok {
			t.Error("qrcode_url should not be set for posts without href")
		}
	}
}

func TestQRCodePlugin_ProcessPost_SlugOnly(t *testing.T) {
	p := NewQRCodePlugin()
	tmpDir := t.TempDir()
	p.SetOutputDir(tmpDir)
	p.SetBaseURL("https://example.com")

	// Create QR code output directory
	qrDir := filepath.Join(tmpDir, "qrcodes")
	if err := os.MkdirAll(qrDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Post with slug but no href (should derive href from slug)
	post := &models.Post{
		Slug: "slug-only",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Check QR code was created
	qrPath := filepath.Join(qrDir, "slug-only.svg")
	if _, err := os.Stat(qrPath); os.IsNotExist(err) {
		t.Error("Expected QR code file to be created for slug-only post")
	}
}

func TestQRCodePlugin_BuildAbsoluteURL(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		href    string
		slug    string
		want    string
	}{
		{
			name:    "with base URL and href",
			baseURL: "https://example.com",
			href:    "/my-post/",
			want:    "https://example.com/my-post/",
		},
		{
			name:    "base URL with trailing slash",
			baseURL: "https://example.com/",
			href:    "/my-post/",
			want:    "https://example.com/my-post/",
		},
		{
			name:    "no base URL",
			baseURL: "",
			href:    "/my-post/",
			want:    "/my-post/",
		},
		{
			name:    "slug fallback",
			baseURL: "https://example.com",
			href:    "",
			slug:    "fallback-post",
			want:    "https://example.com/fallback-post/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewQRCodePlugin()
			p.SetBaseURL(tt.baseURL)

			post := &models.Post{
				Href: tt.href,
				Slug: tt.slug,
			}

			got := p.buildAbsoluteURL(post)
			if got != tt.want {
				t.Errorf("buildAbsoluteURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestQRCodePlugin_BuildFilename(t *testing.T) {
	tests := []struct {
		name   string
		slug   string
		path   string
		format string
		want   string
	}{
		{
			name:   "with slug",
			slug:   "my-post",
			format: "svg",
			want:   "my-post.svg",
		},
		{
			name:   "with slug PNG",
			slug:   "my-post",
			format: "png",
			want:   "my-post.png",
		},
		{
			name:   "fallback to path",
			slug:   "",
			path:   "content/my-file.md",
			format: "svg",
			want:   "my-file.svg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewQRCodePlugin()
			p.config.Format = tt.format

			post := &models.Post{
				Slug: tt.slug,
				Path: tt.path,
			}

			got := p.buildFilename(post)
			if got != tt.want {
				t.Errorf("buildFilename() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestQRCodePlugin_ErrorCorrectionLevels(t *testing.T) {
	tests := []struct {
		level string
	}{
		{"L"},
		{"M"},
		{"Q"},
		{"H"},
		{"invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(_ *testing.T) {
			p := NewQRCodePlugin()
			p.config.ErrorCorrection = tt.level

			// Just verify it doesn't panic
			_ = p.getErrorCorrectionLevel()
		})
	}
}

func TestQRCodePlugin_SVGContent(t *testing.T) {
	p := NewQRCodePlugin()
	p.config.Foreground = "#123456"
	p.config.Background = "#abcdef"

	tmpDir := t.TempDir()
	qrPath := filepath.Join(tmpDir, "test.svg")

	err := p.writeSVG("https://example.com", p.getErrorCorrectionLevel(), qrPath)
	if err != nil {
		t.Errorf("writeSVG() error = %v", err)
	}

	// Read and verify SVG content
	content, err := os.ReadFile(qrPath)
	if err != nil {
		t.Fatal(err)
	}

	svg := string(content)

	// Check SVG structure
	if !strings.Contains(svg, `<?xml version="1.0"`) {
		t.Error("Expected XML declaration")
	}
	if !strings.Contains(svg, `<svg xmlns="http://www.w3.org/2000/svg"`) {
		t.Error("Expected SVG element")
	}
	if !strings.Contains(svg, `fill="#abcdef"`) {
		t.Error("Expected background color")
	}
	if !strings.Contains(svg, `fill="#123456"`) {
		t.Error("Expected foreground color")
	}
}

func TestQRCodePlugin_Config(t *testing.T) {
	p := NewQRCodePlugin()

	// Default config
	cfg := p.Config()
	if !cfg.Enabled {
		t.Error("Expected Enabled to be true by default")
	}
	if cfg.Format != "svg" {
		t.Errorf("Expected Format to be 'svg', got %q", cfg.Format)
	}
	if cfg.Size != 200 {
		t.Errorf("Expected Size to be 200, got %d", cfg.Size)
	}
	if cfg.OutputDir != "qrcodes" {
		t.Errorf("Expected OutputDir to be 'qrcodes', got %q", cfg.OutputDir)
	}
	if cfg.ErrorCorrection != "M" {
		t.Errorf("Expected ErrorCorrection to be 'M', got %q", cfg.ErrorCorrection)
	}

	// Set custom config
	customCfg := models.QRCodeConfig{
		Enabled:         false,
		Format:          "png",
		Size:            300,
		OutputDir:       "custom-qr",
		ErrorCorrection: "H",
		Foreground:      "#ff0000",
		Background:      "#00ff00",
	}
	p.SetConfig(customCfg)

	cfg = p.Config()
	if cfg.Enabled {
		t.Error("Expected Enabled to be false")
	}
	if cfg.Format != "png" {
		t.Errorf("Expected Format to be 'png', got %q", cfg.Format)
	}
	if cfg.Size != 300 {
		t.Errorf("Expected Size to be 300, got %d", cfg.Size)
	}
}

func TestQRCodePlugin_CustomOutputDir(t *testing.T) {
	p := NewQRCodePlugin()
	p.config.OutputDir = "custom-qrcodes"

	tmpDir := t.TempDir()
	p.SetOutputDir(tmpDir)
	p.SetBaseURL("https://example.com")

	// Create QR code output directory
	qrDir := filepath.Join(tmpDir, "custom-qrcodes")
	if err := os.MkdirAll(qrDir, 0o755); err != nil {
		t.Fatal(err)
	}

	post := &models.Post{
		Slug: "custom-dir-test",
		Href: "/custom-dir-test/",
	}

	err := p.processPost(post)
	if err != nil {
		t.Errorf("processPost() error = %v", err)
	}

	// Check QR code file was created in custom directory
	qrPath := filepath.Join(qrDir, "custom-dir-test.svg")
	if _, err := os.Stat(qrPath); os.IsNotExist(err) {
		t.Error("Expected QR code file in custom directory")
	}

	// Check qrcode_url uses custom directory
	qrURL, ok := post.Extra["qrcode_url"].(string)
	if !ok {
		t.Fatal("qrcode_url not set")
	}
	if qrURL != "/custom-qrcodes/custom-dir-test.svg" {
		t.Errorf("qrcode_url = %q, want %q", qrURL, "/custom-qrcodes/custom-dir-test.svg")
	}
}
