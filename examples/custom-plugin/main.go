// A site-owned Markata binary with one deterministic custom transform hook.
// Copy this file into your site's Go module and register your own hooks.
package main

import (
	"fmt"
	"os"

	"github.com/WaylonWalker/markata/cmd/markata-go/cmd"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

type greetingPlugin struct{}

func (*greetingPlugin) Name() string { return "greeting" }

func (*greetingPlugin) Transform(m *lifecycle.Manager) error {
	return m.ProcessPostsConcurrently(func(post *models.Post) error {
		if post.Skip {
			return nil
		}
		post.Set("greeting", "Hello from a site-owned Markata plugin")
		return nil
	})
}

var _ lifecycle.TransformPlugin = (*greetingPlugin)(nil)

func main() {
	plugins.RegisterPluginConstructor("greeting", func() lifecycle.Plugin {
		return &greetingPlugin{}
	})

	// Reuse Markata's existing CLI and build pipeline rather than forking it.
	if err := cmd.Execute(); err != nil {
		if message := cmd.FormatError(err); message != "" {
			_, _ = fmt.Fprintln(os.Stderr, message)
		}
		os.Exit(cmd.ExitCodeForError(err))
	}
}
