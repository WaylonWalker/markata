// Package cmd provides the CLI commands for markata-go.
package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
)

// Version information set via ldflags at build time.
// These are overwritten by goreleaser during releases.
var (
	// Version is the semantic version (e.g., "0.1.0")
	Version = "dev"

	// Commit is the git commit SHA
	Commit = "none"

	// Date is the build date in RFC3339 format
	Date = "unknown"
)

// activeCLIName controls human-facing and User-Agent branding for each entry point.
var activeCLIName = "markata-go"

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version, commit hash, build date, and Go runtime version.`,
	Run: func(cmd *cobra.Command, _ []string) {
		short, err := cmd.Flags().GetBool("short")
		if err != nil {
			errlnf("Error getting flag: %v", err)
			return
		}
		currentCmd = cmd
		if short {
			outln(GetVersion())
			return
		}
		outln(GetVersionInfo())
	},
}

func init() {
	versionCmd.Flags().BoolP("short", "s", false, "Print only the version number")
	rootCmd.AddCommand(versionCmd)
}

// GetVersion returns just the version string.
func GetVersion() string {
	return Version
}

// GetVersionInfo returns full version information as a formatted string.
func GetVersionInfo() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s %s\n", activeCLIName, Version))
	sb.WriteString(fmt.Sprintf("  commit:    %s\n", Commit))
	sb.WriteString(fmt.Sprintf("  built:     %s\n", Date))
	sb.WriteString(fmt.Sprintf("  go:        %s\n", runtime.Version()))
	sb.WriteString(fmt.Sprintf("  os/arch:   %s/%s", runtime.GOOS, runtime.GOARCH))

	// Include VCS info if available from debug.BuildInfo (for dev builds)
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				switch setting.Key {
				case "vcs.revision":
					if Commit == "none" && len(setting.Value) >= 7 {
						sb.WriteString(fmt.Sprintf("\n  vcs.rev:   %s", setting.Value[:7]))
					}
				case "vcs.modified":
					if setting.Value == "true" {
						sb.WriteString(" (dirty)")
					}
				}
			}
		}
	}

	return sb.String()
}

// GetShortVersionInfo returns a one-line version string suitable for User-Agent etc.
func GetShortVersionInfo() string {
	return fmt.Sprintf("%s/%s (%s; %s/%s)", activeCLIName, Version, Commit[:minInt(7, len(Commit))], runtime.GOOS, runtime.GOARCH)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
