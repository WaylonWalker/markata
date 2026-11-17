package visualizations

import (
	"embed"
	"fmt"
	"io/fs"
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
