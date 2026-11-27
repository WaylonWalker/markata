package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/prose"
)

var proseDictionaryCmd = &cobra.Command{
	Use:   "dictionary",
	Short: "Manage accepted project spelling words",
}

func proseDictionaryDocument(args []string, index int) (string, error) {
	path := "untitled.md"
	if len(args) > index {
		path = args[index]
	}
	return filepath.Abs(path)
}

var proseDictionaryAddCmd = &cobra.Command{
	Use:   "add <word> [markdown-file]",
	Short: "Accept a technical term in the project spelling dictionary",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := proseDictionaryDocument(args, 1)
		if err != nil {
			return err
		}
		name, err := prose.AddWord(path, args[0])
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Accepted %q in %s\n", args[0], name)
		return err
	},
}

var proseDictionaryRemoveCmd = &cobra.Command{
	Use:   "remove <word> [markdown-file]",
	Short: "Revoke an accepted project word",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := proseDictionaryDocument(args, 1)
		if err != nil {
			return err
		}
		name, err := prose.RemoveWord(path, args[0])
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "Removed %q from %s\n", args[0], name)
		return err
	},
}

var proseDictionaryListCmd = &cobra.Command{
	Use:   "list [markdown-file]",
	Short: "List accepted project words",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := proseDictionaryDocument(args, 0)
		if err != nil {
			return err
		}
		prefs, err := prose.LoadPreferences(path)
		if err != nil {
			return err
		}
		words := append([]string(nil), prefs.AcceptedWords...)
		sort.Strings(words)
		for _, word := range words {
			if _, err := fmt.Fprintln(cmd.OutOrStdout(), word); err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	proseDictionaryCmd.AddCommand(proseDictionaryAddCmd, proseDictionaryRemoveCmd, proseDictionaryListCmd)
	proseCmd.AddCommand(proseDictionaryCmd)
}
