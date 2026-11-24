// Command shortcuts-modal-render renders the actual Markata Pongo2 template
// for browser acceptance; the browser test must not strip template syntax.
package main

import (
	"fmt"
	"os"

	"github.com/WaylonWalker/markata/pkg/templates"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: shortcuts-modal-render <template-directory>")
		os.Exit(2)
	}
	engine, err := templates.NewEngine(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	config := map[string]interface{}{
		"theme": map[string]interface{}{
			"switcher": map[string]interface{}{
				"mode_toggle": true,
				"enabled":     true,
			},
		},
		"header": map[string]interface{}{
			"show_theme_toggle": true,
		},
		"shortcuts": map[string]interface{}{},
	}
	html, err := engine.RenderToString("partials/shortcuts-modal.html", map[string]interface{}{"config": config})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(html)
}
