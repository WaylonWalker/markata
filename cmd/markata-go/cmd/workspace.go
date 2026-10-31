package cmd

import (
	"encoding/json"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
	"github.com/spf13/cobra"
)

func init() {
	var root, path string
	var pathsOnly bool
	command := &cobra.Command{Use: "workspace", Short: "List or read workspace source as JSON", Args: cobra.NoArgs, RunE: func(cmd *cobra.Command, _ []string) error {
		workspace, err := editorworkspace.New(root)
		if err != nil {
			return err
		}
		if cmd.Flags().Changed("file") {
			source, err := workspace.Read(path)
			if err != nil {
				return err
			}
			return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"source": source, "version": editorworkspace.SourceVersion(source)})
		}
		var files []editorworkspace.File
		if pathsOnly {
			files, err = workspace.Paths()
		} else {
			files, err = workspace.Files()
		}
		if err != nil {
			return err
		}
		return json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{"files": files})
	}}
	command.Flags().BoolVar(&pathsOnly, "paths", false, "list paths without loading document metadata")
	command.Flags().StringVar(&root, "workspace", ".", "workspace root")
	command.Flags().StringVar(&path, "file", "", "read one discovered source file")
	rootCmd.AddCommand(command)
}
