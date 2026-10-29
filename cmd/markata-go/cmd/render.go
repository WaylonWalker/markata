package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/WaylonWalker/markata/pkg/plugins"

	"github.com/spf13/cobra"
)

func newRenderCommand() *cobra.Command {
	var options plugins.StandaloneOptions
	var jsonOutput bool
	var appearance bool
	cmd := &cobra.Command{
		Use: "render", Short: "Render Markdown from stdin without a site build",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if appearance {
				result, err := plugins.BuildStandaloneAppearance()
				if err != nil {
					return err
				}
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			source, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("read markdown: %w", err)
			}
			result, err := plugins.RenderStandalone(string(source), options)
			if err != nil {
				return err
			}
			if jsonOutput {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			_, err = fmt.Fprint(cmd.OutOrStdout(), result.HTML)
			return err
		},
	}
	cmd.Flags().BoolVar(&appearance, "appearance", false, "Return reusable appearance controls and bundled assets as JSON")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Return HTML, feature signals and asset manifest as JSON")
	cmd.Flags().BoolVar(&options.InlineAssets, "inline-assets", false, "Inline bundled CSS/JS in asset tags including vendored browser libraries")
	cmd.Flags().StringVar(&options.Profile, "profile", "web", "Renderer profile: web, offline-interactive, or offline-static (unsupported)")
	cmd.Flags().BoolVar(&options.Unsafe, "unsafe", false, "Allow raw authored HTML (trusted input only)")
	cmd.Flags().StringVar(&options.AssetBase, "asset-base", "/markata-assets", "URL prefix for serving manifest assets")
	return cmd
}

func init() { rootCmd.AddCommand(newRenderCommand()) }
