package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/blogroll"
)

func init() {
	blogrollCmd.AddCommand(newBlogrollExportCmd())
}

func newBlogrollExportCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "export [file.opml]",
		Short:   "Export subscriptions as OPML",
		Example: "  markata blogroll export subscriptions.opml\n  markata blogroll export > subscriptions.opml",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := loadSettingsState()
			if err != nil {
				return fmt.Errorf("load subscription configuration: %w", err)
			}
			feeds := make([]blogroll.OPMLFeed, 0, len(state.cfg.Blogroll.Feeds))
			for i := range state.cfg.Blogroll.Feeds {
				feed := &state.cfg.Blogroll.Feeds[i]
				feeds = append(feeds, blogroll.OPMLFeed{
					URL:         feed.URL,
					Title:       feed.Title,
					SiteURL:     feed.SiteURL,
					Description: feed.Description,
					Category:    feed.Category,
				})
			}
			var buf bytes.Buffer
			if err := blogroll.WriteOPML(&buf, "Markata subscriptions", feeds); err != nil {
				return err
			}
			if len(args) == 0 || args[0] == "-" {
				_, err = io.Copy(cmd.OutOrStdout(), &buf)
				return err
			}
			return writeNewOPMLFile(args[0], &buf)
		},
	}
}

// Never overwrite an existing file, including a symlink to another file.
func writeNewOPMLFile(path string, content io.Reader) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create OPML export: %w", err)
	}
	_, writeErr := io.Copy(f, content)
	closeErr := f.Close()
	if writeErr != nil || closeErr != nil {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("write OPML export failed: %w; remove incomplete export: %w", writeErr, err)
		}
		if writeErr != nil {
			return fmt.Errorf("write OPML export: %w", writeErr)
		}
		return fmt.Errorf("close OPML export: %w", closeErr)
	}
	return nil
}
