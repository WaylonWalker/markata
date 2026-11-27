package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/blogroll"
)

func init() {
	blogrollCmd.AddCommand(newBlogrollDiscoverCmd(), newBlogrollCheckCmd())
}

func newBlogrollDiscoverCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:     "discover <website-or-feed-url>",
		Short:   "Find and validate an advertised feed",
		Example: "  markata blogroll discover https://example.com --json",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 20*time.Second)
			defer cancel()
			result, err := discoverSubscription(ctx, strings.TrimSpace(args[0]))
			if err != nil {
				return fmt.Errorf("discover feed: %w", err)
			}
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(result)
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), result.URL)
			return err
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "write discovered feed information as JSON")
	return cmd
}

func newBlogrollCheckCmd() *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:     "check <handle-or-feed-url>",
		Short:   "Check that a public feed is reachable and has a valid document",
		Example: "  markata blogroll check engineering --json",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := strings.TrimSpace(args[0])
			if !strings.HasPrefix(address, "https://") && !strings.HasPrefix(address, "http://") {
				state, err := loadSettingsState()
				if err != nil {
					return fmt.Errorf("load subscription configuration: %w", err)
				}
				feed, err := selectBlogrollFeed(state.cfg.Blogroll.Feeds, address)
				if err != nil {
					return err
				}
				address = feed.URL
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), 20*time.Second)
			defer cancel()
			client := safeSubscriptionClient()
			defer client.CloseIdleConnections()
			body, finalURL, err := fetchSubscriptionURL(ctx, client, address)
			if err != nil {
				return fmt.Errorf("check feed: %w", err)
			}
			format, err := blogroll.ValidateFeedDocument(body)
			if err != nil {
				return fmt.Errorf("check feed: %w", err)
			}
			if asJSON {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(struct {
					URL    string `json:"url"`
					Format string `json:"format"`
					Valid  bool   `json:"valid"`
				}{finalURL, format, true})
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "Valid %s feed: %s\n", format, finalURL)
			return err
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "write successful feed check as JSON")
	return cmd
}
