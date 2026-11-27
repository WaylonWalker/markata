package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

const (
	blogrollListCommand = "list"
	blogrollStatusActive = "active"
)

func init() {
	blogrollCmd.AddCommand(newBlogrollListCmd())
}

func newBlogrollListCmd() *cobra.Command {
	var asJSON bool
	var filter string
	cmd := &cobra.Command{
		Use:   blogrollListCommand,
		Short: "List external subscriptions and their configuration sources",
		Long:  `List configured blogroll subscriptions without making network requests or modifying configuration.`,
		Example: `  markata blogroll list
  markata blogroll list --filter engineering
  markata blogroll list --json`,
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
		return encoder.Encode(struct {
			Feeds   []subscriptionView `json:"feeds"`
			Enabled bool               `json:"enabled"`
			Target  string             `json:"target,omitempty"`
			Error   string             `json:"error,omitempty"`
		}{feeds, response.Enabled, response.Target, response.Error})
	}
	if !response.Enabled {
		if _, err := fmt.Fprintln(w, "Blogroll is disabled in configuration."); err != nil {
			return err
		}
	}
	if len(feeds) == 0 {
		message := "No subscriptions configured. Run 'markata blogroll add --help' to get started."
		if filter != "" {
			message = "No matching subscriptions. Try a different --filter or list without a filter."
		}
		_, err := fmt.Fprintln(w, message)
		return err
	}
	for i := range feeds {
		feed := &feeds[i]
		title := feed.Title
		if title == "" {
			title = feed.URL
		}
		status := blogrollStatusActive
		if !feed.IsActive() {
			status = "paused"
		}
		if _, err := fmt.Fprintf(w, "%s [%s] (%s)\n  URL: %s\n  Handle: %s\n  Category: %s\n  Source: %s\n", title, status, feed.GetReaderType(), feed.URL, feed.Handle, feed.Category, feed.Source); err != nil {
			return err
		}
	}
	return nil
}
