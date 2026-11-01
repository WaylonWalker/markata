//go:build js && wasm

package plugins

import (
	"fmt"

	"github.com/WaylonWalker/markata/pkg/models"
)

func availableMermaidModes() []string {
	return []string{mermaidModeClient}
}

func newMermaidRenderer(config models.MermaidConfig, _ *mermaidPaletteColors) (mermaidRenderer, error) {
	if config.Mode != "client" {
		return nil, fmt.Errorf("WASM supports client Mermaid rendering only")
	}
	return &clientRenderer{config: config}, nil
}
