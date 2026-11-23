package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	setupLegacyJSON    = "markata-go.json"
	setupAlternateTOML = "markata.toml"
	setupAlternateYAML = "markata.yaml"
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
	configNames := []string{defaultConfigFilename, "markata-go.yaml", "markata-go.yml", setupLegacyJSON}
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
		for _, name := range []string{setupAlternateTOML, setupAlternateYAML, "markata.yml", "markata.json"} {
			candidate := filepath.Join(abs, name)
			if _, err := os.Stat(candidate); err == nil {
				return fmt.Errorf("found %s, which Markata does not automatically discover; rename it to %s (or pass --config explicitly)", candidate, defaultConfigFilename)
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("inspect configuration %s: %w", candidate, err)
			}
		}
		target := filepath.Join(abs, defaultConfigFilename)
		if err := writeSetupConfig(target); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Created starter configuration: %s\n", target)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "Found existing configuration: %s\n", filepath.Join(abs, existing))
	}
	fmt.Fprintf(cmd.OutOrStdout(), "\nWorkspace ready: %s\n\nNext:\n  cd %q\n  markata editor\n\nBuild and preview your site:\n  markata serve\n", abs, abs)
	return nil
}

// writeSetupConfig stages a complete default configuration in the target
// directory and links it into place atomically, without replacing any file.
// Unlike writing directly to the final name, failed writes and interrupts
// cannot leave a partially initialized workspace behind. The hard link is
// created only after the complete staged file has been flushed and closed.
func writeSetupConfig(target string) error {
	staged, err := os.CreateTemp(filepath.Dir(target), ".markata-setup-*.tmp")
	if err != nil {
		return fmt.Errorf("stage configuration: %w", err)
	}
	defer func() { _ = os.Remove(staged.Name()) }()

	if _, err := staged.WriteString(defaultConfigTOML); err != nil {
		_ = staged.Close()
		return fmt.Errorf("write staged configuration: %w", err)
	}
	if err := staged.Chmod(0o644); err != nil {
		_ = staged.Close()
		return fmt.Errorf("set configuration permissions: %w", err)
	}
	if err := staged.Sync(); err != nil {
		_ = staged.Close()
		return fmt.Errorf("flush configuration: %w", err)
	}
	if err := staged.Close(); err != nil {
		return fmt.Errorf("close configuration: %w", err)
	}
	if err := os.Link(staged.Name(), target); err != nil {
		return fmt.Errorf("create configuration without overwriting existing data: %w", err)
	}
	return nil
}
