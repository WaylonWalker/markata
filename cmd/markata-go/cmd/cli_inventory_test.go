package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/spf13/cobra"
)

func TestCLIInventoryRecursiveContract(t *testing.T) {
	root := &cobra.Command{Use: "test-cli", Short: "Fixture CLI"}
	root.PersistentFlags().Bool("global-option", false, "global option")
	group := &cobra.Command{Use: "group", Short: "Fixture group"}
	leaf := &cobra.Command{
		Use:     "leaf <name>",
		Short:   "Fixture leaf",
		Aliases: []string{"other"},
		Hidden:  true,
		Args:    cobra.ExactArgs(1),
		Run:     func(*cobra.Command, []string) {},
	}
	leaf.Flags().String("local-option", "default", "local option")
	group.AddCommand(leaf)
	root.AddCommand(group)
	var out bytes.Buffer
	if err := writeCLIInventory(&out, root); err != nil {
		t.Fatal(err)
	}
	var got []cliInventoryEntry
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[2].Path != "test-cli group leaf" || !got[2].Hidden ||
		!got[2].ArgsChecked || !got[2].Runnable || len(got[2].Aliases) != 1 {
		t.Fatalf("unexpected inventory: %+v", got)
	}
	flags := map[string]cliInventoryFlag{}
	for _, flag := range got[2].Flags {
		if _, exists := flags[flag.Name]; exists {
			t.Fatalf("duplicate flag %s", flag.Name)
		}
		flags[flag.Name] = flag
	}
	if !flags["global-option"].Inherited || flags["local-option"].Inherited ||
		flags["local-option"].Default != "default" || flags["help"].Type != "bool" {
		t.Fatalf("unexpected flag inventory: %+v", flags)
	}
}

func TestRegisteredCLIInventoryIncludesBlogroll(t *testing.T) {
	entries := cliCommandInventory(rootCmd)
	byPath := make(map[string]bool, len(entries))
	for i := range entries {
		entry := &entries[i]
		if byPath[entry.Path] {
			t.Fatalf("duplicate command path: %s", entry.Path)
		}
		byPath[entry.Path] = true
	}
	for _, name := range []string{"add", "list", "edit", "enable", "disable", "remove", "discover", "check", "import", "export", "update"} {
		if !byPath[rootCmd.Name()+" blogroll "+name] {
			t.Errorf("missing registered blogroll command %s", name)
		}
	}
}
