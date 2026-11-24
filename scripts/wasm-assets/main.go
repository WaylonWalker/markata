// Command wasm-assets exports the embedded presentation bundle for the static demo.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/themes"
)

func main() {
	output := flag.String("output", "", "destination directory (required)")
	flag.Parse()
	if *output == "" {
		fmt.Fprintln(os.Stderr, "Use --output DIRECTORY")
		os.Exit(2)
	}
	if err := export(*output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func export(output string) error {
	appearance, err := plugins.BuildStandaloneAppearance()
	if err != nil {
		return err
	}
	names, err := themes.ListStatic()
	if err != nil {
		return err
	}
	assets := map[string][]byte{}
	for _, name := range names {
		if strings.HasPrefix(name, "css/") || strings.HasPrefix(name, "js/") || strings.HasPrefix(name, "assets/vendor/") {
			data, readErr := themes.ReadStatic(name)
			if readErr != nil {
				return readErr
			}
			assets[name] = data
		}
	}
	for name, data := range appearance.Assets {
		assets[name] = data
	}
	// Font URLs resolve inside this app, including under a URL subdirectory.
	assets["css/fonts.css"] = []byte(strings.ReplaceAll(string(assets["css/fonts.css"]), "/assets/fonts/", "../assets/fonts/"))
	sample, err := plugins.RenderStandalone("```go\nvar x = 1\n```", plugins.StandaloneOptions{})
	if err != nil {
		return err
	}
	for _, asset := range sample.Assets {
		if asset.ID == "css/chroma.css" {
			assets[asset.ID] = []byte(asset.Content)
		}
	}
	for name, data := range assets {
		path := filepath.Join(output, "appearance", filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		//nolint:gosec // Generated static assets are public web content.
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
	}
	appearance.Assets = nil
	data, err := json.Marshal(appearance)
	if err != nil {
		return err
	}
	//nolint:gosec // The presentation manifest is public web content.
	return os.WriteFile(filepath.Join(output, "appearance.json"), data, 0o644)
}
