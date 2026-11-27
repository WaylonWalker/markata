package cmd

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/blogroll"
	"github.com/WaylonWalker/markata/pkg/models"
)

const (
	blogrollImportCategory    = "category"
	blogrollImportDescription = "description"
)

func init() {
	blogrollCmd.AddCommand(newBlogrollImportCmd())
}

func newBlogrollImportCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:     "import <file.opml>",
		Short:   "Import subscriptions from OPML without fetching feeds",
		Example: "  markata blogroll import subscriptions.opml --dry-run",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("open OPML: %w", err)
			}
			defer f.Close()
			feeds, err := blogroll.ReadOPML(f)
			if err != nil {
				return err
			}
			state, err := loadSettingsState()
			if err != nil {
				return fmt.Errorf("load subscription configuration: %w", err)
			}
			added, skipped, err := importBlogrollOPML(state, feeds, dryRun)
			if err != nil {
				return err
			}
			message := "Imported"
			if dryRun {
				message = "Would import"
			}
			_, err = fmt.Fprintf(cmd.OutOrStdout(), "%s %d subscriptions; skipped %d duplicate URLs. Imported feeds are unverified.\n", message, added, skipped)
			return err
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "validate and count changes without writing configuration")
	return cmd
}

// planOPMLImport skips exact duplicate URLs without changing existing feeds.
func planOPMLImport(existing []models.ExternalFeedConfig, incoming []blogroll.OPMLFeed) (add []blogroll.OPMLFeed, skipped int) {
	seen := make(map[string]bool, len(existing)+len(incoming))
	for i := range existing {
		seen[existing[i].URL] = true
	}
	add = make([]blogroll.OPMLFeed, 0, len(incoming))
	skipped = 0
	for i := range incoming {
		feed := &incoming[i]
		if seen[feed.URL] {
			skipped++
			continue
		}
		seen[feed.URL] = true
		add = append(add, *feed)
	}
	return add, skipped
}

// importBlogrollOPML writes the complete import in one atomic transaction,
// then reloads effective config and rolls back if the import was shadowed.
func importBlogrollOPML(state *settingsState, incoming []blogroll.OPMLFeed, dryRun bool) (added, skipped int, err error) {
	add, skipped := planOPMLImport(state.cfg.Blogroll.Feeds, incoming)
	if len(add) == 0 {
		return 0, skipped, nil
	}
	target, err := newFeedTarget(state)
	if err != nil {
		return 0, skipped, err
	}
	backup, err := backupFile(target)
	if err != nil {
		return 0, skipped, err
	}
	content, err := renderOPMLImport(string(backup.data), add, state.cfg.Blogroll.Enabled)
	if err != nil {
		return 0, skipped, err
	}
	if dryRun {
		return len(add), skipped, nil
	}
	mode := backup.mode
	if !backup.existed {
		mode = 0o600
	}
	if err := writeSubscriptionFile(target, []byte(content), mode); err != nil {
		return 0, skipped, err
	}
	updated, err := loadSettingsState()
	if err == nil {
		err = verifyOPMLImport(state.cfg.Blogroll.Feeds, add, updated.cfg.Blogroll)
	}
	if err != nil {
		if rollbackErr := backup.restore(); rollbackErr != nil {
			return 0, skipped, fmt.Errorf("import failed: %w; rollback failed: %w", err, rollbackErr)
		}
		return 0, skipped, fmt.Errorf("import rolled back: %w", err)
	}
	return len(add), skipped, nil
}

func renderOPMLImport(content string, feeds []blogroll.OPMLFeed, enabled bool) (string, error) {
	content = clearEmptySubscriptionList(content)
	if !enabled {
		content = enableSubscriptions(content)
	}
	prefix := detectBlogrollTablePrefix(content)
	var result strings.Builder
	result.WriteString(strings.TrimRight(content, "\r\n"))
	for i := range feeds {
		feed := &feeds[i]
		req := subscriptionRequest{
			URL: feed.URL, Title: feed.Title, Description: feed.Description,
			Category: feed.Category, Type: string(models.ReaderFeedTypeWritten), Active: true,
		}
		fields, err := subscriptionFields(req)
		if err != nil {
			return "", err
		}
		result.WriteString("\n\n[[" + prefix + ".feeds]]\n")
		result.WriteString("url = " + strconv.Quote(feed.URL) + "\n")
		for _, field := range []string{listSortTitle, blogrollImportDescription, blogrollImportCategory, "type", blogrollStatusActive} {
			result.WriteString(field + " = " + fields[field] + "\n")
		}
		if feed.SiteURL != "" {
			result.WriteString("site_url = " + strconv.Quote(feed.SiteURL) + "\n")
		}
	}
	output := result.String()
	if _, err := toml.Decode(output, &map[string]any{}); err != nil {
		return "", fmt.Errorf("generated TOML was invalid: %w", err)
	}
	return output, nil
}

func verifyOPMLImport(before []models.ExternalFeedConfig, added []blogroll.OPMLFeed, after models.BlogrollConfig) error {
	if !after.Enabled || len(after.Feeds) != len(before)+len(added) {
		return errors.New("imported subscriptions did not take effect in composed configuration")
	}
	for i := range before {
		if !reflect.DeepEqual(before[i], after.Feeds[i]) {
			return errors.New("import unexpectedly changed an existing subscription")
		}
	}
	for i := range added {
		want, got := &added[i], &after.Feeds[len(before)+i]
		if got.URL != want.URL || got.Title != strings.TrimSpace(want.Title) ||
			got.Category != strings.TrimSpace(want.Category) || got.Description != strings.TrimSpace(want.Description) ||
			got.SiteURL != want.SiteURL || !got.IsActive() {
			return errors.New("imported subscription metadata did not take effect")
		}
	}
	return nil
}
