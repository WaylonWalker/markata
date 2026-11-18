package visualizations

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
	"sync"
)

const (
	// ExcalidrawPackageVersion identifies the bundled upstream editor package.
	ExcalidrawPackageVersion = "0.18.1"
	// ExcalidrawSVGExportVersion identifies Markata's SVG export contract.
	ExcalidrawSVGExportVersion = "excalidraw-" + ExcalidrawPackageVersion + "-static-svg-v2"
)

//go:embed excalidraw
var excalidrawFiles embed.FS

// ExcalidrawAssetFS returns the locally bundled browser assets used by the
// editor and the static SVG renderer.
func ExcalidrawAssetFS() (fs.FS, error) {
	assets, err := fs.Sub(excalidrawFiles, "excalidraw")
	if err != nil {
		return nil, fmt.Errorf("open embedded Excalidraw assets: %w", err)
	}
	return assets, nil
}

// ReadExcalidrawAsset reads one locally bundled Excalidraw asset.
func ReadExcalidrawAsset(name string) ([]byte, error) {
	assets, err := ExcalidrawAssetFS()
	if err != nil {
		return nil, err
	}
	return fs.ReadFile(assets, name)
}

// Cache the embedded asset fingerprint: this can be requested for each
// workbench preview, but its input files do not change during a process.
var excalidrawRendererVersionOnce sync.Once
var excalidrawRendererVersionValue string
var excalidrawRendererVersionErr error

// ExcalidrawRendererVersion fingerprints the pinned SVG export contract,
// bundled font bytes, and the exported JS/CSS entrypoints. The JS entrypoint
// references content-hashed chunks, so rebuilt dependency changes also change
// this fingerprint.
func ExcalidrawRendererVersion() (string, error) {
	excalidrawRendererVersionOnce.Do(func() {
		excalidrawRendererVersionValue, excalidrawRendererVersionErr = computeExcalidrawRendererVersion()
	})
	return excalidrawRendererVersionValue, excalidrawRendererVersionErr
}

func computeExcalidrawRendererVersion() (string, error) {
	assets, err := ExcalidrawAssetFS()
	if err != nil {
		return "", err
	}
	manifestData, err := fs.ReadFile(assets, "excalidraw-fonts/manifest.json")
	if err != nil {
		return "", fmt.Errorf("read local Excalidraw font manifest: %w", err)
	}
	var manifest struct {
		RendererVersion string   `json:"renderer_version"`
		Fonts           []string `json:"fonts"`
	}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return "", fmt.Errorf("parse local Excalidraw font manifest: %w", err)
	}
	if manifest.RendererVersion != ExcalidrawSVGExportVersion || len(manifest.Fonts) == 0 {
		return "", fmt.Errorf("local Excalidraw font manifest has unsupported version or no fonts")
	}
	fonts := append([]string(nil), manifest.Fonts...)
	sort.Strings(fonts)
	hash := sha256.New()
	previous := ""
	for _, font := range fonts {
		if !fs.ValidPath(font) || path.Clean(font) != font || strings.Contains(font, "\\") || (previous != "" && font == previous) {
			return "", fmt.Errorf("invalid or duplicate path in Excalidraw font manifest: %q", font)
		}
		previous = font
		name := path.Join("excalidraw-fonts", font)
		data, err := fs.ReadFile(assets, name)
		if err != nil {
			return "", fmt.Errorf("read bundled Excalidraw font %q: %w", font, err)
		}
		_, _ = hash.Write([]byte(font))
		_, _ = hash.Write([]byte{0})
		_, _ = hash.Write(data)
	}
	for _, name := range []string{"excalidraw-editor.js", "excalidraw-editor.css"} {
		data, err := fs.ReadFile(assets, name)
		if err != nil {
			return "", fmt.Errorf("read bundled Excalidraw export asset %q: %w", name, err)
		}
		_, _ = hash.Write([]byte(name))
		_, _ = hash.Write([]byte{0})
		_, _ = hash.Write(data)
	}
	return manifest.RendererVersion + ":" + hex.EncodeToString(hash.Sum(nil)), nil
}
