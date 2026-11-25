package plugins

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/WaylonWalker/markata/pkg/visualizations"
)

const excalidrawArtifactOptions = "padding=20;fonts=inline-local-woff2;svg=1"

var errExcalidrawArtifactMissing = errors.New("editor-generated SVG artifact is missing")

type excalidrawSVGArtifact struct {
	Format          string `json:"format"`
	SourceSHA256    string `json:"source_sha256"`
	RendererVersion string `json:"renderer_version"`
	Theme           string `json:"theme"`
	Options         string `json:"options"`
	Encoding        string `json:"encoding"`
	SVGSHA256       string `json:"svg_sha256"`
	SVG             string `json:"svg"`
}

type cachedExcalidrawRenderer struct {
	mode            string
	darkMode        bool
	rendererVersion string
	version         string
}

func newCachedExcalidrawRenderer(mode string, darkMode bool) (excalidrawSVGRenderer, error) {
	rendererVersion, err := visualizations.ExcalidrawRendererVersion()
	if err != nil {
		return nil, err
	}
	return &cachedExcalidrawRenderer{
		mode:            mode,
		darkMode:        darkMode,
		rendererVersion: rendererVersion,
		version:         excalidrawVersionForTheme(rendererVersion, darkMode),
	}, nil
}

func (r *cachedExcalidrawRenderer) ID() string      { return "excalidraw-editor-artifact" }
func (r *cachedExcalidrawRenderer) Version() string { return r.version }
func (r *cachedExcalidrawRenderer) Close() error    { return nil }

func (r *cachedExcalidrawRenderer) RenderSVG(_ context.Context, req visualizations.Request) ([]byte, error) {
	if _, err := validateExcalidrawScene(req.Source); err != nil {
		return nil, err
	}
	artifact, err := readExcalidrawArtifact(req.Source)
	if err != nil {
		if errors.Is(err, errExcalidrawArtifactMissing) {
			return nil, r.staleArtifact(err.Error())
		}
		return nil, err
	}
	canonical, err := canonicalExcalidrawScene(req.Source)
	if err != nil {
		return nil, fmt.Errorf("fingerprint Excalidraw scene: %w", err)
	}
	sceneHash := sha256.Sum256(canonical)
	if artifact.SourceSHA256 != hex.EncodeToString(sceneHash[:]) {
		return nil, r.staleArtifact("scene fingerprint does not match")
	}
	theme := themeModeLight
	if r.darkMode {
		theme = themeModeDark
	}
	if artifact.RendererVersion != r.rendererVersion {
		return nil, r.staleArtifact("bundled Excalidraw renderer or font assets changed")
	}
	if artifact.Theme != theme {
		return nil, r.staleArtifact("site theme changed")
	}
	if artifact.Options != excalidrawArtifactOptions {
		return nil, r.staleArtifact("SVG export options changed")
	}
	if artifact.Format != "markata-excalidraw-svg-v1" || artifact.Encoding != "gzip+base64" {
		return nil, errors.New("unsupported Excalidraw SVG artifact format; reapply the scene in the Go workbench")
	}
	svg, err := decodeExcalidrawArtifactSVG(artifact)
	if err != nil {
		return nil, err
	}
	if _, err := visualizations.NewSVGAsset(svg); err != nil {
		return nil, fmt.Errorf("cached Excalidraw SVG failed static validation: %w", err)
	}
	return svg, nil
}

func (r *cachedExcalidrawRenderer) staleArtifact(reason string) error {
	if r.mode == excalidrawModeCached {
		return fmt.Errorf("cached Excalidraw SVG is stale (%s); reapply the scene in the Go workbench", reason)
	}
	return fmt.Errorf("no matching cached Excalidraw SVG (%s); reapply the scene in the Go workbench or explicitly set [markata-go.excalidraw].mode = \"chromium\"", reason)
}

func readExcalidrawArtifact(source string) (excalidrawSVGArtifact, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(source), &root); err != nil {
		return excalidrawSVGArtifact{}, fmt.Errorf("invalid scene JSON: %w", err)
	}
	var metadata struct {
		Artifact json.RawMessage `json:"excalidraw_svg"`
	}
	if err := json.Unmarshal(root["_markata"], &metadata); err != nil || len(metadata.Artifact) == 0 || string(metadata.Artifact) == "null" {
		return excalidrawSVGArtifact{}, errExcalidrawArtifactMissing
	}
	var artifact excalidrawSVGArtifact
	if err := json.Unmarshal(metadata.Artifact, &artifact); err != nil {
		return excalidrawSVGArtifact{}, fmt.Errorf("invalid Excalidraw SVG artifact: %w", err)
	}
	return artifact, nil
}

func canonicalExcalidrawScene(source string) ([]byte, error) {
	decoder := json.NewDecoder(strings.NewReader(source))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	if decoder.Decode(new(any)) != io.EOF {
		return nil, errors.New("scene JSON contains trailing data")
	}
	root, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("scene JSON must be an object")
	}
	delete(root, "_markata")
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(root); err != nil {
		return nil, err
	}
	canonical := bytes.TrimSuffix(output.Bytes(), []byte("\n"))
	canonical = bytes.ReplaceAll(canonical, []byte(`\u2028`), []byte("\u2028"))
	canonical = bytes.ReplaceAll(canonical, []byte(`\u2029`), []byte("\u2029"))
	return canonical, nil
}

func decodeExcalidrawArtifactSVG(artifact excalidrawSVGArtifact) ([]byte, error) {
	if len(artifact.SVG) > base64.StdEncoding.EncodedLen(visualizations.MaxSVGBytes) {
		return nil, errors.New("cached Excalidraw SVG artifact exceeds its compressed size limit")
	}
	compressed, err := base64.StdEncoding.Strict().DecodeString(artifact.SVG)
	if err != nil {
		return nil, fmt.Errorf("decode cached Excalidraw SVG: %w", err)
	}
	reader := bytes.NewReader(compressed)
	stream, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("decompress cached Excalidraw SVG: %w", err)
	}
	stream.Multistream(false)
	svg, readErr := io.ReadAll(io.LimitReader(stream, visualizations.MaxSVGBytes+1))
	closeErr := stream.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read cached Excalidraw SVG: %w", readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close cached Excalidraw SVG: %w", closeErr)
	}
	if len(svg) > visualizations.MaxSVGBytes {
		return nil, fmt.Errorf("cached Excalidraw SVG exceeds the %d-byte limit", visualizations.MaxSVGBytes)
	}
	if reader.Len() != 0 {
		return nil, errors.New("cached Excalidraw SVG contains trailing compressed data")
	}
	sum := sha256.Sum256(svg)
	if artifact.SVGSHA256 != hex.EncodeToString(sum[:]) {
		return nil, errors.New("cached Excalidraw SVG checksum does not match")
	}
	return svg, nil
}
