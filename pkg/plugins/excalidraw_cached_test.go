package plugins

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/visualizations"
)

const testExcalidrawSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="20" height="10"><rect width="20" height="10"/></svg>`

func TestCachedExcalidrawRenderer_UsesMatchingArtifact(t *testing.T) {
	renderer := newTestCachedExcalidrawRenderer(t, excalidrawModeAuto, true)
	artifact := makeTestExcalidrawArtifact(t, renderer.rendererVersion)
	source := withTestExcalidrawArtifact(t, testExcalidrawTextScene, artifact)

	got, err := renderer.RenderSVG(context.Background(), visualizations.Request{Source: source})
	if err != nil {
		t.Fatalf("RenderSVG() error = %v", err)
	}
	if string(got) != testExcalidrawSVG {
		t.Fatalf("RenderSVG() = %q, want cached SVG", got)
	}
}

func TestCachedExcalidrawRenderer_CanonicalFingerprintIgnoresFormatting(t *testing.T) {
	renderer := newTestCachedExcalidrawRenderer(t, excalidrawModeCached, true)
	artifact := makeTestExcalidrawArtifact(t, renderer.rendererVersion)
	formatted := `{
  "files": {},
  "elements": [{ "type": "text", "text": "Flow   overview", "id": "text-1", "isDeleted": false }],
  "version": 2,
  "type": "excalidraw",
  "appState": {}
}`
	source := withTestExcalidrawArtifact(t, formatted, artifact)

	if _, err := renderer.RenderSVG(context.Background(), visualizations.Request{Source: source}); err != nil {
		t.Fatalf("RenderSVG() rejected equivalent formatted scene: %v", err)
	}
}

func TestCachedExcalidrawRenderer_RejectsMissingAndStaleArtifacts(t *testing.T) {
	renderer := newTestCachedExcalidrawRenderer(t, excalidrawModeAuto, true)
	if _, err := renderer.RenderSVG(context.Background(), visualizations.Request{Source: testExcalidrawTextScene}); err == nil || !strings.Contains(err.Error(), "artifact is missing") {
		t.Fatalf("missing artifact error = %v, want actionable missing-artifact error", err)
	}
	artifact := makeTestExcalidrawArtifact(t, renderer.rendererVersion)
	changed := strings.Replace(testExcalidrawTextScene, "Flow   overview", "Updated", 1)
	source := withTestExcalidrawArtifact(t, changed, artifact)
	if _, err := renderer.RenderSVG(context.Background(), visualizations.Request{Source: source}); err == nil || !strings.Contains(err.Error(), "scene fingerprint") {
		t.Fatalf("stale scene error = %v, want a source fingerprint error", err)
	}
}

func TestCachedExcalidrawRenderer_RejectsThemeAndRendererChanges(t *testing.T) {
	darkRenderer := newTestCachedExcalidrawRenderer(t, excalidrawModeAuto, true)
	artifact := makeTestExcalidrawArtifact(t, darkRenderer.rendererVersion)
	source := withTestExcalidrawArtifact(t, testExcalidrawTextScene, artifact)

	lightRenderer := newTestCachedExcalidrawRenderer(t, excalidrawModeAuto, false)
	if _, err := lightRenderer.RenderSVG(context.Background(), visualizations.Request{Source: source}); err == nil || !strings.Contains(err.Error(), "site theme changed") {
		t.Fatalf("theme mismatch error = %v, want a stale theme error", err)
	}

	darkRenderer.rendererVersion += ":changed"
	if _, err := darkRenderer.RenderSVG(context.Background(), visualizations.Request{Source: source}); err == nil || !strings.Contains(err.Error(), "renderer or font assets changed") {
		t.Fatalf("renderer mismatch error = %v, want a stale renderer error", err)
	}
}

func TestDecodeExcalidrawArtifactSVG_EnforcesExpandedSVGLimit(t *testing.T) {
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write(bytes.Repeat([]byte("x"), visualizations.MaxSVGBytes+1)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	artifact := excalidrawSVGArtifact{Encoding: "gzip+base64", SVG: base64.StdEncoding.EncodeToString(compressed.Bytes())}
	if _, err := decodeExcalidrawArtifactSVG(artifact); err == nil || !strings.Contains(err.Error(), "exceeds the") {
		t.Fatalf("decodeExcalidrawArtifactSVG() error = %v, want expanded-size rejection", err)
	}
}

func newTestCachedExcalidrawRenderer(t *testing.T, mode string, darkMode bool) *cachedExcalidrawRenderer {
	t.Helper()
	renderer, err := newCachedExcalidrawRenderer(mode, darkMode)
	if err != nil {
		t.Fatal(err)
	}
	cached, ok := renderer.(*cachedExcalidrawRenderer)
	if !ok {
		t.Fatalf("newCachedExcalidrawRenderer() returned %T", renderer)
	}
	return cached
}

func makeTestExcalidrawArtifact(t *testing.T, rendererVersion string) excalidrawSVGArtifact {
	t.Helper()
	svg := testExcalidrawSVG
	canonical, err := canonicalExcalidrawScene(testExcalidrawTextScene)
	if err != nil {
		t.Fatal(err)
	}
	sourceHash := sha256.Sum256(canonical)
	svgHash := sha256.Sum256([]byte(svg))
	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := io.WriteString(writer, svg); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return excalidrawSVGArtifact{
		Format:          "markata-excalidraw-svg-v1",
		SourceSHA256:    hex.EncodeToString(sourceHash[:]),
		RendererVersion: rendererVersion,
		Theme:           themeModeDark,
		Options:         excalidrawArtifactOptions,
		Encoding:        "gzip+base64",
		SVGSHA256:       hex.EncodeToString(svgHash[:]),
		SVG:             base64.StdEncoding.EncodeToString(compressed.Bytes()),
	}
}

func withTestExcalidrawArtifact(t *testing.T, source string, artifact excalidrawSVGArtifact) string {
	t.Helper()
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(source), &root); err != nil {
		t.Fatal(err)
	}
	artifactJSON, err := json.Marshal(artifact)
	if err != nil {
		t.Fatal(err)
	}
	metadataJSON, err := json.Marshal(map[string]json.RawMessage{"excalidraw_svg": artifactJSON})
	if err != nil {
		t.Fatal(err)
	}
	root["_markata"] = metadataJSON
	encoded, err := json.Marshal(root)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}
