package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/prose"
)

var proseJSON bool
var proseIgnore []string

var proseCmd = &cobra.Command{
	Use:   "prose",
	Short: "Check writing style without changing your text",
}

var proseCheckCmd = &cobra.Command{
	Use:   "check <markdown-file>",
	Short: "Show non-blocking prose suggestions in a Markdown file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("read %s: %w", args[0], err)
		}
		preferences, err := prose.LoadPreferences(args[0])
		if err != nil {
			return err
		}
		options := preferences.Options()
		for _, rule := range proseIgnore {
			options.IgnoredRules[strings.TrimSpace(rule)] = true
		}
		issues := prose.Check(string(content), options)
		if proseJSON {
			return json.NewEncoder(cmd.OutOrStdout()).Encode(issues)
		}
		for _, issue := range issues {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s:%d:%d: %s: %s\n",
				args[0], issue.Line+1, issue.Column+1, issue.Rule, issue.Message); err != nil {
				return err
			}
		}
		if len(issues) == 0 {
			_, err = fmt.Fprintln(cmd.OutOrStdout(), "No prose suggestions.")
			return err
		}
		return nil
	},
}

var proseIgnoreCmd = &cobra.Command{
	Use:   "ignore <rule-id> [markdown-file]",
	Short: "Always ignore a style rule in this project",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "untitled.md"
		if len(args) > 1 {
			path = args[1]
		}
		absolute, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		name, err := prose.IgnoreRule(absolute, args[0])
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Ignored %s in %s\n", args[0], name)
		return err
	},
}

var proseIgnorePhraseCmd = &cobra.Command{
	Use:   "ignore-phrase <rule-id> <phrase> [markdown-file]",
	Short: "Always ignore one phrase without disabling the style rule",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := "untitled.md"
		if len(args) > 2 {
			path = args[2]
		}
		absolute, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		name, err := prose.IgnorePhrase(absolute, args[0], args[1])
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Ignored phrase %q (%s) in %s\n", args[1], args[0], name)
		return err
	},
}

func init() {
	proseCheckCmd.Flags().BoolVar(&proseJSON, "json", false, "emit stable JSON diagnostics for tools and agents")
	proseCheckCmd.Flags().StringSliceVar(&proseIgnore, "ignore-rule", nil, "suppress a rule ID (repeatable)")
	proseCmd.AddCommand(proseCheckCmd)
	proseCmd.AddCommand(proseIgnoreCmd)
	proseCmd.AddCommand(proseIgnorePhraseCmd)
	rootCmd.AddCommand(proseCmd)
}
