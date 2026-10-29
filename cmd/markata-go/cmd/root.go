// Package cmd provides the CLI commands for markata-go.
package cmd

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/WaylonWalker/markata/pkg/suggest"
	"github.com/spf13/cobra"
)

var (
	// cfgFile is the path to the config file specified via --config flag.
	cfgFile string

	// mergeConfigFiles is a list of additional config files to merge with the base config.
	// These are applied in order, with later files taking precedence over earlier ones.
	mergeConfigFiles []string

	// outputDir is the output directory specified via --output flag.
	outputDir string

	// verbose enables verbose output.
	verbose bool

	// quiet suppresses non-essential status output.
	quiet bool

	// noColor disables ANSI color output on all streams.
	noColor bool

	// forceColor forces ANSI color output even for non-TTY streams.
	forceColor bool

	// logFormat controls centralized log formatting.
	logFormat string

	// noInput disables prompts and interactive UI.
	noInput bool

	// cpuProfile is the path to write CPU profile data.
	cpuProfile string

	// memProfile is the path to write memory profile data.
	memProfile string

	// cpuProfileFile holds the open CPU profile file for cleanup.
	cpuProfileFile *os.File
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "markata-go",
	Short: "A plugin-driven static site generator",
	Long: `Markata-go is a static site generator with a powerful feed system.

It processes markdown files with YAML frontmatter and generates a static website
with support for multiple feed formats (RSS, Atom, JSON Feed), automatic archives,
tag pages, and more.

Example usage:
  markata-go build           # Build the site
  markata-go serve           # Build and serve locally with live reload
  markata-go new "My Post"   # Create a new post
  markata-go config show     # Show resolved configuration

Common help:
  markata-go help build      # Explain a subcommand
  markata-go list posts      # Inspect posts from the CLI
  markata-go explain content # Agent workflow for site content

Documentation:
  https://github.com/WaylonWalker/markata/tree/main/docs

Issues:
  https://github.com/WaylonWalker/markata/issues

Profiling:
  markata-go build --cpuprofile cpu.prof   # Write CPU profile
  markata-go build --memprofile mem.prof   # Write memory profile

  Analyze with:
  go tool pprof cpu.prof
  go tool pprof -http=:8080 cpu.prof`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       Version,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 && !strings.HasSuffix(strings.ToLower(args[0]), ".md") && !strings.HasSuffix(strings.ToLower(args[0]), ".markdown") {
			return unknownCommandError(cmd, args[0])
		}
		return cobra.MaximumNArgs(1)(cmd, args)
	},
	RunE: runRootCommand,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		currentCmd = cmd
		if cmd.Flags().Lookup("dag") != nil {
			if err := validateDAGEnvironment(cmd); err != nil {
				return err
			}
		}
		// Build Lab resolves and binds every child to its copied workspace. Do
		// not let the caller's site-directory selection change the parent
		// process's fixture resolution or child command setup.
		// Triage must still work when a requested site directory is broken.
		// It reports the inaccessible path rather than failing before diagnostics.
		if !isBuildLabCommand(cmd) && !isTriageCommand(cmd) {
			if err := activateSiteDir(); err != nil {
				return err
			}
		}
		if noColor && forceColor {
			return fmt.Errorf("cannot use --color and --no-color together")
		}
		configureTUIColor()

		if err := configureCLIThemeFromConfig(); err != nil {
			return err
		}

		// Start CPU profiling if requested
		if cpuProfile != "" {
			f, err := os.Create(cpuProfile)
			if err != nil {
				return fmt.Errorf("failed to create CPU profile: %w", err)
			}
			cpuProfileFile = f
			if err := pprof.StartCPUProfile(f); err != nil {
				cpuProfileFile = nil
				f.Close()
				return fmt.Errorf("failed to start CPU profile: %w", err)
			}
			verbosef("CPU profiling enabled, writing to %s", cpuProfile)
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
		currentCmd = cmd
		// Stop CPU profiling
		stopCPUProfile()

		// Write memory profile if requested
		if memProfile != "" {
			f, err := os.Create(memProfile)
			if err != nil {
				return fmt.Errorf("failed to create memory profile: %w", err)
			}
			defer f.Close()

			// Get the heap profile (most useful for memory analysis)
			if err := pprof.WriteHeapProfile(f); err != nil {
				return fmt.Errorf("failed to write memory profile: %w", err)
			}
			verbosef("Memory profile written to %s", memProfile)
		}
		return nil
	},
}

func runRootCommand(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	if !strings.HasSuffix(strings.ToLower(args[0]), ".md") && !strings.HasSuffix(strings.ToLower(args[0]), ".markdown") {
		return unknownCommandError(cmd, args[0])
	}
	return runBuildCommand(cmd, args)
}

func unknownCommandError(parent *cobra.Command, name string) error {
	families := make([]suggest.Family, 0)
	for _, child := range parent.Commands() {
		if child.IsAvailableCommand() && !child.IsAdditionalHelpTopicCommand() {
			families = append(families, suggest.Family{Name: child.Name(), Aliases: child.Aliases, SuggestFor: child.SuggestFor})
		}
	}
	matches := suggest.Ranked(name, families, 2)
	labels := make(map[string]string, len(families))
	for _, family := range families {
		label := family.Name
		if len(family.Aliases) > 0 {
			label += " (aliases: " + strings.Join(family.Aliases, ", ") + ")"
		}
		labels[family.Name] = label
	}
	for i, match := range matches {
		matches[i] = labels[match]
	}
	return newUsageError(&inputError{
		problem: fmt.Sprintf("unknown command %q", name),
		matches: matches,
		prefix:  parent.CommandPath() + " ",
		next:    "Run '" + parent.CommandPath() + " --help' to see available commands and aliases.",
	})
}

// ExecuteAs invokes the shared CLI under a selected executable name.
// The legacy executable retains its historical spelling and output.
const canonicalCLIName = "markata"

func ExecuteAs(name string) error {
	if name != canonicalCLIName {
		return Execute()
	}

	previousUse, previousLong, previousName := rootCmd.Use, rootCmd.Long, activeCLIName
	rootCmd.Use = canonicalCLIName
	rootCmd.Long = strings.ReplaceAll(rootCmd.Long, "Markata-go is", "Markata is")
	rootCmd.Long = strings.ReplaceAll(rootCmd.Long, "markata-go ", canonicalCLIName+" ")
	activeCLIName = canonicalCLIName
	defer func() {
		rootCmd.Use, rootCmd.Long, activeCLIName = previousUse, previousLong, previousName
	}()
	return Execute()
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	// Cobra skips post-run hooks when a command returns an error. Keep CPU
	// profiles valid for failed Build Lab/product runs as well as successful
	// commands.
	defer stopCPUProfile()
	prepareCommandGroups(rootCmd)
	return rootCmd.Execute()
}

// Cobra otherwise shows help and exits successfully for a mistyped child of a
// non-runnable command group. Give those groups a help action and argument
// validator so unknown children return a useful usage error.
func prepareCommandGroups(parent *cobra.Command) {
	for _, child := range parent.Commands() {
		prepareCommandGroups(child)
	}
	if parent == rootCmd || !parent.HasAvailableSubCommands() || parent.Runnable() {
		return
	}
	parent.Args = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return unknownCommandError(cmd, args[0])
		}
		return nil
	}
	parent.RunE = func(cmd *cobra.Command, _ []string) error { return cmd.Help() }
}

func stopCPUProfile() {
	if cpuProfileFile == nil {
		return
	}
	pprof.StopCPUProfile()
	_ = cpuProfileFile.Close()
	cpuProfileFile = nil
	verbosef("CPU profile written to %s", cpuProfile)
}

func runHelpCommand(_ *cobra.Command, args []string) error {
	if len(args) == 0 {
		return rootCmd.Help()
	}
	target, remaining, err := rootCmd.Find(args)
	if err != nil {
		return newUsageError(err)
	}
	if len(remaining) > 0 {
		return unknownCommandError(target, remaining[0])
	}
	return target.Help()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:   "help [command]",
		Short: "Show help for a command",
		RunE:  runHelpCommand,
	})
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		renderCommandHelp(cmd)
	})
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return newUsageError(clarifyFlagError(cmd, err))
	})

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path (default: auto-discover)")
	rootCmd.PersistentFlags().StringVar(&siteDir, "site-dir", "", "site directory (overrides MARKATA_GO_SITE_DIR and current directory)")
	rootCmd.PersistentFlags().StringSliceVarP(&mergeConfigFiles, "merge-config", "m", nil, "additional config file(s) to merge with base config (can be specified multiple times)")
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "output directory (overrides config)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential status output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&forceColor, "color", false, "force ANSI color output")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable ANSI color on all streams")
	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "auto", "log formatting: auto, plain, or rich")
	rootCmd.PersistentFlags().BoolVar(&noInput, "no-input", false, "disable prompts and interactive UI")

	// Profiling flags
	rootCmd.PersistentFlags().StringVar(&cpuProfile, "cpuprofile", "", "write CPU profile to file")
	rootCmd.PersistentFlags().StringVar(&memProfile, "memprofile", "", "write memory profile to file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Config initialization is handled by the core package when needed
	verbosef("Verbose mode enabled")
}
