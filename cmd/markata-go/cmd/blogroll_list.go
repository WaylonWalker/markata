package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	blogrollCmd.AddCommand(newBlogrollListCmd())
}

func newBlogrollListCmd() *cobra.Command {
	var asJSON bool
	var filter string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List external subscriptions and their configuration sources",
		Long:  "List configured blogroll subscriptions without making network requests or modifying configuration.",
		Example: "  markata blogroll list\n" +
			"  markata blogroll list --filter engineering\n" +
			"  markata blogroll list --json",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			state, err := loadSettingsState()
			if err != nil {
				return fmt.Errorf("load subscription configuration: %w", err)
			}
			return renderBlogrollList(cmd.OutOrStdout(), describeSubscriptions(state), filter, asJSON)
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "write machine-readable subscription information")
	cmd.Flags().StringVar(&filter, "filter", "", "filter by title, URL, handle, category or tag")
	return cmd
}

func renderBlogrollList(w io.Writer, response subscriptionResponse, filter string, asJSON bool) error {
	filter = strings.ToLower(strings.TrimSpace(filter))
	feeds := make([]subscriptionView, 0, len(response.Feeds))
	for i := range response.Feeds {
		feed := &response.Feeds[i]
		search := strings.Join([]string{feed.Title, feed.URL, feed.Handle, feed.Category, strings.Join(feed.Tags, " ")}, " ")
		if filter == "" || strings.Contains(strings.ToLower(search), filter) {
			feeds = append(feeds, *feed)
		}
	}
	response.Feeds = feeds
	if asJSON {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(response)
	}
	if !response.Enabled {
		if _, err := fmt.Fprintln(w, "Blogroll is disabled in configuration."); err != nil {
			return err
		}
	}
	if len(feeds) == 0 {
		_, err := fmt.Fprintln(w, "No matching subscriptions. Run 'markata blogroll add' to add a feed.")
		return err
	}
	for i := range feeds {
		feed := &feeds[i]
		title := feed.Title
		if title == "" {
			title = feed.URL
		}
		status := "active"
		if !feed.IsActive() {
			status = "paused"
		}
		if _, err := fmt.Fprintf(w, "%s [%s] (%s)\n  URL: %s\n  Handle: %s\n  Category: %s\n  Source: %s\n", title, status, feed.GetReaderType(), feed.URL, feed.Handle, feed.Category, feed.Source); err != nil {
			return err
		}
	}
	return nil
}
