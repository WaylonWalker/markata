package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/blogroll"
	"github.com/WaylonWalker/markata/pkg/models"
)

func init() {
	blogrollCmd.AddCommand(newBlogrollActiveCmd("enable", true), newBlogrollActiveCmd("disable", false), newBlogrollRemoveCmd())
}

func newBlogrollActiveCmd(name string, active bool) *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:     name + " <handle-or-url>",
		Short:   name + " an external subscription",
		Example: "  markata blogroll " + name + " engineering --dry-run",
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
			req := subscriptionRequest{Action: subscriptionActionActive, URL: feed.URL, Active: active}
			if err := changeSubscriptionWithOptions(state, req, dryRun); err != nil {
				return fmt.Errorf("change subscription: %w", err)
			}
			message := "Updated"
			if dryRun {
				message = "Would update"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %s: active = %t\n", message, feed.URL, active)
			return err
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "validate the proposed change without writing configuration")
	return cmd
}

// selectBlogrollFeed requires an exact URL or an unambiguous case-insensitive
// handle. Never mutate a feed on a substring match.
func selectBlogrollFeed(feeds []models.ExternalFeedConfig, selector string) (models.ExternalFeedConfig, error) {
	return blogroll.SelectSubscription(feeds, selector)
}

func newBlogrollRemoveCmd() *cobra.Command {
	var dryRun, yes bool
	cmd := &cobra.Command{
		Use:     "remove <handle-or-url>",
		Short:   "Remove an external subscription",
		Example: "  markata blogroll remove engineering --dry-run",
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
			if !dryRun && !yes {
				if err := confirmBlogrollRemoval(feed.URL); err != nil {
					return err
				}
			}
			req := subscriptionRequest{Action: subscriptionActionRemove, URL: feed.URL}
			if err := changeSubscriptionWithOptions(state, req, dryRun); err != nil {
				return fmt.Errorf("remove subscription: %w", err)
			}
			message := "Removed"
			if dryRun {
				message = "Would remove"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", message, feed.URL)
			return err
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "validate removal without writing configuration")
	cmd.Flags().BoolVar(&yes, "yes", false, "confirm removal without prompting")
	return cmd
}

func confirmBlogrollRemoval(address string) error {
	if noInput || !inputIsTerminal() || !outputIsTerminal() {
		return newUsageError(&Diagnostic{
			Problem:    "subscription removal requires confirmation",
			Fix:        "Use --yes to confirm removal or --dry-run to preview.",
			StatusCode: exitCodeUsage,
		})
	}
	confirmed := false
	form := huh.NewForm(huh.NewGroup(huh.NewConfirm().
		Title("Remove subscription?").
		Description(address).
		Value(&confirmed))).WithTheme(createHuhThemeForPalette(currentCLIPalette))
	if err := form.Run(); err != nil {
		return err
	}
	if !confirmed {
		return fmt.Errorf("subscription removal canceled; no changes written")
	}
	return nil
}
