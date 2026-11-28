package cmd

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// cliInventory is machine-readable command metadata, not a claim that a
// command's runtime behavior, documentation or error paths have passed audit.
type cliInventoryEntry struct {
	Path         string             `json:"path"`
	Usage        string             `json:"usage"`
	Summary      string             `json:"summary"`
	Aliases      []string           `json:"aliases,omitempty"`
	Hidden       bool               `json:"hidden"`
	Deprecated   string             `json:"deprecated,omitempty"`
	Runnable     bool               `json:"runnable"`
	ArgsChecked  bool               `json:"args_validator_registered"`
	Examples     string             `json:"examples,omitempty"`
	Flags        []cliInventoryFlag `json:"flags"`
}

type cliInventoryFlag struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Default     string `json:"default"`
	Inherited   bool   `json:"inherited"`
	Hidden      bool   `json:"hidden"`
	Deprecated  string `json:"deprecated,omitempty"`
}

var cliInventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "Print the recursive CLI command and flag inventory as JSON",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return writeCLIInventory(cmd.OutOrStdout(), cmd.Root())
	},
}

func init() {
	group := &cobra.Command{
		Use:   "cli",
		Short: "Inspect the installed CLI",
	}
	group.AddCommand(cliInventoryCmd)
	rootCmd.AddCommand(group)
}

// cliCommandInventory traverses the actual registered commands, including
// hidden commands. Flag defaults are declaration defaults, not live values.
func cliCommandInventory(root *cobra.Command) []cliInventoryEntry {
	result := make([]cliInventoryEntry, 0)
	var walk func(*cobra.Command)
	walk = func(cmd *cobra.Command) {
		cmd.InitDefaultHelpFlag()
		entry := cliInventoryEntry{
			Path:        cmd.CommandPath(),
			Usage:       cmd.UseLine(),
			Summary:     cmd.Short,
			Aliases:     cmd.Aliases,
			Hidden:      cmd.Hidden,
			Deprecated:  cmd.Deprecated,
			Runnable:    cmd.Runnable(),
			ArgsChecked: cmd.Args != nil,
			Examples:    cmd.Example,
			Flags:       make([]cliInventoryFlag, 0),
		}
		collect := func(flags *pflag.FlagSet, inherited bool) {
			flags.VisitAll(func(flag *pflag.Flag) {
				entry.Flags = append(entry.Flags, cliInventoryFlag{
					Name:        flag.Name,
					Shorthand:   flag.Shorthand,
					Type:        flag.Value.Type(),
					Description: flag.Usage,
					Default:     flag.DefValue,
					Inherited:   inherited,
					Hidden:      flag.Hidden,
					Deprecated:  flag.Deprecated,
				})
			})
		}
		collect(cmd.NonInheritedFlags(), false)
		collect(cmd.InheritedFlags(), true)
		result = append(result, entry)
		for _, child := range cmd.Commands() {
			walk(child)
		}
	}
	walk(root)
	return result
}

func writeCLIInventory(out io.Writer, root *cobra.Command) error {
	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cliCommandInventory(root))
}
