//go:build js && wasm

// Command markata-wasm exposes the standalone renderer inside a browser worker.
package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"
	"unicode/utf8"

	"github.com/WaylonWalker/markata/pkg/plugins"
)

func main() {
	renderers := map[string]*plugins.StandaloneRenderer{
		"web:false": plugins.NewStandaloneRenderer(plugins.StandaloneOptions{
			Profile: plugins.RenderProfileWeb, AssetBase: "./appearance", OmitAssetContent: true,
		}),
	}
	render := js.FuncOf(func(_ js.Value, args []js.Value) (response any) {
		defer func() {
			if value := recover(); value != nil {
				response = errorJSON(fmt.Sprintf("render failed: %v", value))
			}
		}()
		if len(args) < 1 || len(args) > 3 || args[0].Type() != js.TypeString {
			return errorJSON("Pass a Markdown string, optional profile, and inline-assets flag")
		}
		profile := plugins.RenderProfileWeb
		if len(args) >= 2 {
			if args[1].Type() != js.TypeString {
				return errorJSON("Render profile must be a string")
			}
			profile = args[1].String()
		}
		inline := false
		if len(args) >= 3 {
			if args[2].Type() != js.TypeBoolean {
				return errorJSON("Inline-assets must be a boolean")
			}
			inline = args[2].Bool()
		}
		source := args[0].String()
		if len(source) > 200000 || !utf8.ValidString(source) {
			return errorJSON("Markdown must contain at most 200,000 UTF-8 bytes")
		}
		// WASM intentionally does not duplicate its large CSS/JS bundle inside
		// the module. The static host and service worker provide those bytes,
		// and the parent preview hydrates them before inserting into its iframe.
		options := plugins.StandaloneOptions{
			Profile: profile, AssetBase: "./appearance", OmitAssetContent: true,
			OfflineAssetHost: profile == plugins.RenderProfileOfflineInteractive,
		}
		key := fmt.Sprintf("%s:%t", profile, inline)
		renderer, ok := renderers[key]
		if !ok {
			renderer = plugins.NewStandaloneRenderer(options)
			renderers[key] = renderer
		}
		result, err := renderer.Render(source)
		if err != nil {
			return errorJSON(err.Error())
		}
		// Bundled dependencies are served once as static files. Only generated code
		// travels with each render, never the multi-megabyte Mermaid bundle.
		if options.OmitAssetContent {
			for i := range result.Assets {
				asset := &result.Assets[i]
				if asset.ID != "css/chroma.css" && len(asset.ID) >= 10 && asset.ID[:10] == "generated-" {
					continue
				}
				asset.Content = ""
			}
		}
		data, err := json.Marshal(result)
		if err != nil {
			return errorJSON(err.Error())
		}
		return string(data)
	})
	js.Global().Set("markataRender", render)
	select {}
}

func errorJSON(message string) string {
	data, _ := json.Marshal(map[string]string{"error": message})
	return string(data)
}
