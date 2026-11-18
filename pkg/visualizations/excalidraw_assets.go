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

// ExcalidrawRendererVersion identifies the SVG exporter and every bundled
// Excalidraw asset, including its JavaScript and fonts.
func ExcalidrawRendererVersion() (string, error) {
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
	return manifest.RendererVersion + ":" + hex.EncodeToString(hash.Sum(nil)), nil
}
