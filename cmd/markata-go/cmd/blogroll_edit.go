package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	blogrollCmd.AddCommand(newBlogrollEditCmd())
}

func newBlogrollEditCmd() *cobra.Command {
	var options subscriptionRequest
	var dryRun bool
	cmd := &cobra.Command{
		Use:     "edit <handle-or-url>",
		Short:   "Edit selected fields of an external subscription",
		Example: "  markata blogroll edit engineering --category Technology --dry-run",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := loadSettingsState()
			if err != nil {
				return fmt.Errorf("load subscription configuration: %w", err)
			}
			feed, err := selectBlogrollFeed(state.cfg.Blogroll.Feeds, args[0])
			if err != nil {
				return err
			}
			req := subscriptionRequest{
				Action:      subscriptionActionUpdate,
				URL:         feed.URL,
				Title:       feed.Title,
				Description: feed.Description,
				Category:    feed.Category,
				Type:        string(feed.GetReaderType()),
				Tags:        strings.Join(feed.Tags, ","),
				Handle:      feed.Handle,
				Active:      feed.IsActive(),
			}
			for _, field := range []struct {
				name  string
				dest  *string
				value string
			}{
				{"title", &req.Title, options.Title},
				{"description", &req.Description, options.Description},
				{"category", &req.Category, options.Category},
				{"type", &req.Type, options.Type},
				{"tags", &req.Tags, options.Tags},
				{"handle", &req.Handle, options.Handle},
			} {
				if cmd.Flags().Changed(field.name) {
					*field.dest = strings.TrimSpace(field.value)
					req.selectedFields = append(req.selectedFields, field.name)
				}
			}
			if len(req.selectedFields) == 0 {
				return newUsageError(&Diagnostic{
					Problem:    "no subscription fields selected",
					Fix:        "Specify a field flag such as --title or --category; see blogroll edit --help.",
					StatusCode: exitCodeUsage,
				})
			}
			if err := changeSubscriptionWithOptions(state, req, dryRun); err != nil {
				return fmt.Errorf("edit subscription: %w", err)
			}
			message := "Edited"
			if dryRun {
				message = "Would edit"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %s: %s\n", message, feed.URL, strings.Join(req.selectedFields, ", "))
			return err
		},
	}
	cmd.Flags().StringVar(&options.Title, "title", "", "set the subscription title")
	cmd.Flags().StringVar(&options.Description, "description", "", "set the description")
	cmd.Flags().StringVar(&options.Category, "category", "", "set the category")
	cmd.Flags().StringVar(&options.Type, "type", "", "set the reader type: written, video, podcast")
	cmd.Flags().StringVar(&options.Tags, "tags", "", "replace comma-separated tags; empty clears tags")
	cmd.Flags().StringVar(&options.Handle, "handle", "", "set the subscription handle")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "validate the selected changes without writing configuration")
	return cmd
}
