package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/WaylonWalker/markata/pkg/prose"
	"github.com/spf13/cobra"
)

var proseJSON bool
var proseIgnore []string

var proseCmd = &cobra.Command{
	Use: "prose",
	Short: "Check writing style without changing your text",
}

var proseCheckCmd = &cobra.Command{
	Use: "check <markdown-file>",
	Short: "Show non-blocking prose suggestions in a Markdown file",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(args[0])
		if err != nil { return fmt.Errorf("read %s: %w", args[0], err) }
		ignored := make(map[string]bool, len(proseIgnore))
		for _, rule := range proseIgnore { ignored[strings.TrimSpace(rule)] = true }
		issues := prose.Check(string(content), prose.Options{IgnoredRules: ignored})
		if proseJSON {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(issues)
		}
		for _, issue := range issues {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s:%d:%d: %s: %s\n",
				args[0], issue.Line+1, issue.Column+1, issue.Rule, issue.Message); err != nil {
				return err
			}
		}
		if len(issues) == 0 { _, err = fmt.Fprintln(cmd.OutOrStdout(), "No prose suggestions."); return err }
		return nil // Editorial observations never block publishing or CI by default.
	},
}

func init() {
	proseCheckCmd.Flags().BoolVar(&proseJSON, "json", false, "emit stable JSON diagnostics for tools and agents")
	proseCheckCmd.Flags().StringSliceVar(&proseIgnore, "ignore-rule", nil, "suppress a rule ID (repeatable)")
	proseCmd.AddCommand(proseCheckCmd)
	rootCmd.AddCommand(proseCmd)
}
