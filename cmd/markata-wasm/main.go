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
	renderer := plugins.NewStandaloneRenderer(plugins.StandaloneOptions{AssetBase: "./appearance", OmitAssetContent: true})
	render := js.FuncOf(func(_ js.Value, args []js.Value) (response any) {
		defer func() {
			if value := recover(); value != nil {
				response = errorJSON(fmt.Sprintf("render failed: %v", value))
			}
		}()
		if len(args) != 1 || args[0].Type() != js.TypeString {
			return errorJSON("Pass one Markdown string")
		}
		source := args[0].String()
		if len(source) > 200000 || !utf8.ValidString(source) {
			return errorJSON("Markdown must contain at most 200,000 UTF-8 bytes")
		}
		result, err := renderer.Render(source)
		if err != nil {
			return errorJSON(err.Error())
		}
		// Bundled dependencies are served once as static files. Only generated code
		// travels with each render, never the multi-megabyte Mermaid bundle.
		for i := range result.Assets {
			asset := &result.Assets[i]
			if asset.ID != "css/chroma.css" && len(asset.ID) >= 10 && asset.ID[:10] == "generated-" {
				continue
			}
			asset.Content = ""
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
