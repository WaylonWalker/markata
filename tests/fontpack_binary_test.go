package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/internal/testbinary"
)

func TestInstalledBinaryUsesEmbeddedFontpackOutsideCheckout(t *testing.T) {
	if testing.Short() {
		t.Skip("installed binary integration test")
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repo := filepath.Dir(filepath.Dir(thisFile))
	binary, err := testbinary.Resolve(repo, t.TempDir())
	if err != nil {
		t.Fatalf("resolve standalone binary: %v", err)
	}
	site := t.TempDir()
	if err := os.MkdirAll(filepath.Join(site, "pages"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(site, "markata-go.toml"), []byte("[markata-go]\nfontpack = \"field-notebook\"\noutput_dir = \"public\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(site, "pages", "note.md"), []byte("---\ntitle: Embedded fonts\npublished: true\n---\n\nHello from an installed binary.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(binary, "build", "--fast", "--config", filepath.Join(site, "markata-go.toml"))
	cmd.Dir = site // deliberately outside the Markata source checkout
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("standalone build: %v\n%s", err, output)
	}
	css, err := os.ReadFile(filepath.Join(site, "public", "css", "fonts.css"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(css), "@font-face") {
		t.Fatal("generated font CSS has no @font-face declaration")
	}
	entries, err := os.ReadDir(filepath.Join(site, "public", "assets", "fonts"))
	if err != nil {
		t.Fatal(err)
	}
	fontCount := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".woff2") {
			fontCount++
		}
	}
	if fontCount == 0 {
		t.Fatal("standalone build emitted no WOFF2 assets")
	}
	for _, line := range strings.Split(string(css), "\n") {
		const marker = "url('"
		start := strings.Index(line, marker)
		if start < 0 {
			continue
		}
		start += len(marker)
		end := strings.Index(line[start:], "'")
		if end < 0 {
			t.Fatalf("malformed font URL in CSS line %q", line)
		}
		url := strings.TrimPrefix(line[start:start+end], "/")
		if _, err := os.Stat(filepath.Join(site, "public", url)); err != nil {
			t.Fatalf("CSS URL %q has no emitted file: %v", url, err)
		}
	}
	report := exec.Command(binary, "fonts", "report", "--config", filepath.Join(site, "markata-go.toml"))
	report.Dir = site
	reportOutput, err := report.CombinedOutput()
	if err != nil || !strings.Contains(string(reportOutput), "Files emitted:") {
		t.Fatalf("fonts report did not inspect configured output directory: %v\n%s", err, reportOutput)
	}
	invalidConfig := filepath.Join(site, "markata-go.toml")
	if err := os.WriteFile(invalidConfig, []byte("[markata-go]\nfontpacks_file = \"./does-not-exist.yaml\"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	verifyBuiltin := exec.Command(binary, "fonts", "verify")
	// Run the built-in verification without a site config. On Windows, the
	// child process can otherwise discover the deliberately invalid config
	// above through its working directory before it selects the embedded
	// catalog.
	verifyBuiltin.Dir = t.TempDir()
	if output, err := verifyBuiltin.CombinedOutput(); err != nil {
		t.Fatalf("bare fonts verify depended on invalid custom catalog: %v\n%s", err, output)
	}
	verifyCustom := exec.Command(binary, "fonts", "verify", "field-notebook")
	verifyCustom.Dir = site
	if output, err := verifyCustom.CombinedOutput(); err == nil {
		t.Fatalf("pack-specific verification unexpectedly ignored invalid custom catalog:\n%s", output)
	}
}
