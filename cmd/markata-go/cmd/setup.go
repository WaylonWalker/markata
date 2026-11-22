package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// setupCmd provides a safe, non-interactive first-run path without changing
// configuration belonging to an existing site.
var setupCmd = &cobra.Command{
	Use:   "setup [workspace]",
	Short: "Prepare a Markdown workspace for the Markata editor",
	Long:  "Find an existing Markata configuration or create a starter configuration without overwriting files. The editor and builder remain independent until explicitly started.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSetupCommand,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetupCommand(cmd *cobra.Command, args []string) error {
	workspace := "."
	if len(args) > 0 {
		workspace = args[0]
	}
	abs, err := filepath.Abs(workspace)
	if err != nil {
		return fmt.Errorf("resolve workspace: %w", err)
	}
	if info, err := os.Stat(abs); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("workspace is not a directory: %s", abs)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(abs, 0o755); err != nil {
			return fmt.Errorf("create workspace: %w", err)
		}
	} else {
		return fmt.Errorf("inspect workspace: %w", err)
	}

	// Respect recognized legacy config filenames; never overwrite existing sites.
	configNames := []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"}
	existing := ""
	for _, name := range configNames {
		candidate := filepath.Join(abs, name)
		if _, err := os.Stat(candidate); err == nil {
			existing = name
			break
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect configuration %s: %w", candidate, err)
		}
	}
	if existing == "" {
		// Alternate filenames are not currently auto-discovered by Markata.
		// Do not silently create a second configuration beside one of them.
		for _, name := range []string{"markata.toml", "markata.yaml", "markata.yml", "markata.json"} {
			candidate := filepath.Join(abs, name)
			if _, err := os.Stat(candidate); err == nil {
				return fmt.Errorf("found %s, which Markata does not automatically discover; rename it to markata-go.toml (or pass --config explicitly)", candidate)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("inspect configuration %s: %w", candidate, err)
			}
		}
		target := filepath.Join(abs, defaultConfigFilename)
		f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			return fmt.Errorf("create configuration (will not overwrite): %w", err)
		}
		if _, err := f.WriteString(defaultConfigTOML); err != nil {
			_ = f.Close()
			return fmt.Errorf("write configuration: %w", err)
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("close configuration: %w", err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created starter configuration: %s\n", target)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Found existing configuration: %s\n", filepath.Join(abs, existing))
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\nWorkspace ready: %s\n\nNext:\n  cd %q\n  markata editor\n\nBuild and preview your site:\n  markata serve\n", abs, abs)
	return nil
}
