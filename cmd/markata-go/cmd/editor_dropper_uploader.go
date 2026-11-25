package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/WaylonWalker/markata/pkg/workbench"

	"github.com/spf13/cobra"
)

var (
	dropperUploaderSocketPath string
	dropperUploaderEndpoint   string
)

var editorDropperUploaderCmd = &cobra.Command{
	Use:          "editor-dropper-uploader",
	Short:        "Run the hosted editor's local Dropper upload broker",
	Long:         "Internal sidecar process for hosted editor deployments. It keeps the Dropper token in this process and exposes only validated image uploads over a Unix socket.",
	Example:      "  markata-go editor-dropper-uploader --socket /run/markata-dropper/uploader.sock --endpoint https://dropper.example.com/api/upload?format=json",
	Args:         cobra.NoArgs,
	Hidden:       true,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return workbench.ServeDropperUploader(ctx, workbench.DropperUploaderServerConfig{
			SocketPath: dropperUploaderSocketPath,
			Endpoint:   dropperUploaderEndpoint,
			Token:      os.Getenv("DROPPER_TOKEN"),
		})
	},
}

func init() {
	rootCmd.AddCommand(editorDropperUploaderCmd)
	editorDropperUploaderCmd.Flags().StringVar(&dropperUploaderSocketPath, "socket", "", "Unix socket path for local editor upload requests")
	editorDropperUploaderCmd.Flags().StringVar(&dropperUploaderEndpoint, "endpoint", "", "fixed Dropper JSON upload endpoint")
}
