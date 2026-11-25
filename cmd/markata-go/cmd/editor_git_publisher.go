package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/WaylonWalker/markata/pkg/workbench"

	"github.com/spf13/cobra"
)

var (
	gitPublisherSourceDir  string
	gitPublisherSocketPath string
	gitPublisherOriginURL  string
	gitPublisherBranch     string
	gitPublisherPrivateKey string
	gitPublisherKnownHosts string
)

var editorGitPublisherCmd = &cobra.Command{
	Use:          "editor-git-publisher",
	Short:        "Run the hosted editor's local Git credential broker",
	Long:         "Internal sidecar process for hosted editor deployments. It serves only on a Unix socket and accepts remote verification and constrained single-Markdown pushes for one configured repository and branch.",
	Example:      "  markata-go editor-git-publisher --source-dir /data/source --socket /run/markata-publisher/publisher.sock --origin-url https://github.com/owner/repo --branch refs/heads/main --private-key /run/markata-git/id_ed25519 --known-hosts /etc/ssh/known_hosts",
	Args:         cobra.NoArgs,
	Hidden:       true,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		return workbench.ServeGitPublisher(ctx, workbench.GitPublisherServerConfig{
			SourceDir:  gitPublisherSourceDir,
			SocketPath: gitPublisherSocketPath,
			OriginURL:  gitPublisherOriginURL,
			Branch:     gitPublisherBranch,
			PrivateKey: gitPublisherPrivateKey,
			KnownHosts: gitPublisherKnownHosts,
		})
	},
}

func init() {
	rootCmd.AddCommand(editorGitPublisherCmd)
	editorGitPublisherCmd.Flags().StringVar(&gitPublisherSourceDir, "source-dir", "", "fixed Git checkout to publish")
	editorGitPublisherCmd.Flags().StringVar(&gitPublisherSocketPath, "socket", "", "Unix socket path for local editor requests")
	editorGitPublisherCmd.Flags().StringVar(&gitPublisherOriginURL, "origin-url", "", "allowlisted HTTPS GitHub repository URL")
	editorGitPublisherCmd.Flags().StringVar(&gitPublisherBranch, "branch", "", "allowlisted fully qualified branch, such as refs/heads/main")
	editorGitPublisherCmd.Flags().StringVar(&gitPublisherPrivateKey, "private-key", "", "private key file mounted only in this process")
	editorGitPublisherCmd.Flags().StringVar(&gitPublisherKnownHosts, "known-hosts", "", "verified SSH known_hosts file")
}
