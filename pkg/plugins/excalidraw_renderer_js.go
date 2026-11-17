//go:build js

package plugins

import (
	"errors"

	"github.com/WaylonWalker/markata/pkg/models"
)

func newExcalidrawRenderer(models.ExcalidrawConfig) (excalidrawSVGRenderer, error) {
	return nil, errors.New("static Excalidraw rendering requires the native Go renderer with Chrome or Chromium")
}

func availableExcalidrawModes() []string { return nil }
