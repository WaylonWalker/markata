// Package d2native is an isolated Go 1.27 current-D2 spike, not a Markata plugin.
// It intentionally lives in a nested module to avoid altering production dependencies.
package d2native

import (
	"context"
	"errors"

	"github.com/d2lang/d2/d2graph"
	"github.com/d2lang/d2/d2layouts/d2dagrelayout"
	"github.com/d2lang/d2/d2lib"
	"github.com/d2lang/d2/d2renderers/d2svg"
	"github.com/d2lang/d2/d2themes/d2themescatalog"
	"github.com/d2lang/d2/lib/log"
	"github.com/d2lang/d2/lib/textmeasure"
)

// Engine reuses D2's text measurement ruler across renders.
// This prototype intentionally tests only the built-in Dagre layout.
type Engine struct {
	ruler *textmeasure.Ruler
}

// NewEngine initializes the expensive text ruler once per rendering engine.
func NewEngine() (*Engine, error) {
	ruler, err := textmeasure.NewRuler()
	if err != nil {
		return nil, err
	}
	return &Engine{ruler: ruler}, nil
}

// Render compiles D2 source to self-contained SVG with the Grape Soda theme.
// It does not import files: CompileOptions.FS is intentionally nil.
func (engine *Engine) Render(ctx context.Context, source string) ([]byte, error) {
	if engine == nil || engine.ruler == nil {
		return nil, errors.New("D2 native engine is not initialized")
	}
	resolver := func(name string) (d2graph.LayoutGraph, error) {
		if name != "dagre" {
			return nil, errors.New("D2 spike supports only dagre")
		}
		return d2dagrelayout.DefaultLayout, nil
	}
	renderOpts := &d2svg.RenderOpts{
		ThemeID: &d2themescatalog.GrapeSoda.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: resolver,
		Ruler:          engine.ruler,
	}
	diagram, _, err := d2lib.Compile(log.WithDefault(ctx), source, compileOpts, renderOpts)
	if err != nil {
		return nil, err
	}
	return d2svg.Render(diagram, renderOpts)
}
