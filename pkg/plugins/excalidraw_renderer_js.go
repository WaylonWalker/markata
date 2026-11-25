//go:build js

package plugins

import (
	"errors"
	"fmt"

	"github.com/WaylonWalker/markata/pkg/models"
)

func newExcalidrawRenderer(config models.ExcalidrawConfig, darkMode bool) (excalidrawSVGRenderer, error) {
	switch config.Mode {
	case excalidrawModeAuto, excalidrawModeCached:
		return newCachedExcalidrawRenderer(config.Mode, darkMode)
	case excalidrawModeChromium:
		return nil, errors.New("excalidraw mode chromium requires the native Go renderer with Chrome or Chromium")
	default:
		return nil, fmt.Errorf("unsupported Excalidraw mode %q (supported: auto, cached, chromium)", config.Mode)
	}
}

func availableExcalidrawModes() []string { return []string{excalidrawModeAuto, excalidrawModeCached} }
