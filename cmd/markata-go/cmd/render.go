package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/renderprotocol"

	"github.com/spf13/cobra"
)

func newRenderCommand() *cobra.Command {
	var options plugins.StandaloneOptions
	var jsonOutput bool
	var appearance bool
	var workspaceRoot string
	var capabilities bool
	var stdio bool
	cmd := &cobra.Command{
		Use: "render", Short: "Render Markdown from stdin without a site build",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if stdio {
				if stdioHasConflictingFlags(cmd, appearance, capabilities, jsonOutput) {
					return fmt.Errorf("--stdio takes per-request JSON options and cannot be combined with other render mode flags")
				}
				if cmd.Flags().Changed("workspace") {
					workspace, err := editorworkspace.New(workspaceRoot)
					if err != nil {
						return err
					}
					options.Config, err = workspace.Config()
					if err != nil {
						return err
					}
				}
				return renderprotocol.ServeWithOptions(cmd.InOrStdin(), cmd.OutOrStdout(), options)
			}
			if capabilities {
				if appearance || cmd.Flags().Changed("workspace") {
					return fmt.Errorf("--capabilities cannot be combined with --appearance or --workspace")
				}
				return json.NewEncoder(cmd.OutOrStdout()).Encode(plugins.GetStandaloneCapabilities())
			}
			if cmd.Flags().Changed("workspace") {
				workspace, err := editorworkspace.New(workspaceRoot)
				if err != nil {
					return err
				}
				options.Config, err = workspace.Config()
				if err != nil {
					return err
				}
			}
			if appearance {
				result, err := plugins.BuildStandaloneAppearanceForSite(options.Config, workspaceRoot)
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
	cmd.Flags().StringVar(&workspaceRoot, "workspace", "", "apply site configuration from this workspace")
	cmd.Flags().BoolVar(&stdio, "stdio", false, "Keep rendering JSON Lines requests from stdin until EOF or shutdown")
	cmd.Flags().BoolVar(&capabilities, "capabilities", false, "Describe supported standalone render features and modes as JSON")
	cmd.Flags().BoolVar(&appearance, "appearance", false, "Return reusable appearance controls and bundled assets as JSON")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Return HTML, feature signals and asset manifest as JSON")
	cmd.Flags().BoolVar(&options.InlineAssets, "inline-assets", false, "Inline bundled CSS/JS in asset tags including vendored browser libraries")
	cmd.Flags().BoolVar(&options.Unsafe, "unsafe", false, "Allow raw authored HTML (trusted input only)")
	cmd.Flags().StringVar(&options.AssetBase, "asset-base", "/markata-assets", "URL prefix for serving manifest assets")
	cmd.Flags().StringVar(&options.MermaidMode, "mermaid-mode", "client", "Mermaid rendering: client (default) or go (native SVG)")
	return cmd
}

func stdioHasConflictingFlags(cmd *cobra.Command, appearance, capabilities, jsonOutput bool) bool {
	return appearance || capabilities || jsonOutput ||
		cmd.Flags().Changed("inline-assets") ||
		cmd.Flags().Changed("asset-base") ||
		cmd.Flags().Changed("unsafe") ||
		cmd.Flags().Changed("mermaid-mode")
}

func init() { rootCmd.AddCommand(newRenderCommand()) }
