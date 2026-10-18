package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/editorial"
	"github.com/spf13/cobra"
)

var editorialAuthor, editorialRole, editorialSummary, editorialSession, editorialPrimaryAuthor, editorialRecordedBy string
var editorialJSON, editorialDryRun bool

var editorialCmd = &cobra.Command{Use: "editorial", Short: "Create attributed, reviewable post edits", SilenceUsage: true}
var editorialStartCmd = &cobra.Command{Use: "start <post...>", Short: "Capture a clean baseline for an editorial session", Example: "markata-go editorial start posts/article.md --author editor --role copy-editor", Args: cobra.MinimumNArgs(1), RunE: runEditorialStart}
var editorialFinishCmd = &cobra.Command{Use: "finish", Short: "Commit an editorial session and record its links", Example: "markata-go editorial finish --session UUID --summary 'Corrected grammar'", Args: cobra.NoArgs, RunE: runEditorialFinish}
var editorialCommitCmd = &cobra.Command{Use: "commit <post...>", Short: "Commit explicitly owned current post changes", Example: "markata-go editorial commit posts/article.md --author linker --role internal-link-editor --summary 'Added links'", Args: cobra.MinimumNArgs(1), RunE: runEditorialCommit}

func init() {
	editorialCmd.AddCommand(editorialStartCmd, editorialFinishCmd, editorialCommitCmd)
	rootCmd.AddCommand(editorialCmd)
	for _, c := range []*cobra.Command{editorialStartCmd, editorialCommitCmd} {
		c.Flags().StringVar(&editorialAuthor, "author", "", "stable contributor ID")
		c.Flags().StringVar(&editorialRole, "role", "", "copy-editor or internal-link-editor")
		c.Flags().BoolVar(&editorialDryRun, "dry-run", false, "validate and preview without changing files")
		c.Flags().StringVar(&editorialRecordedBy, "recorded-by", "", "ID of the person or agent recording attribution (defaults to author)")
		c.Flags().StringVar(&editorialPrimaryAuthor, "primary-author", "", "primary writer ID when a post has no author list")
		if err := c.MarkFlagRequired("author"); err != nil {
			panic(err)
		}
		if err := c.MarkFlagRequired("role"); err != nil {
			panic(err)
		}
	}
	editorialFinishCmd.Flags().StringVar(&editorialSession, "session", "", "editorial session UUID")
	editorialFinishCmd.Flags().StringVar(&editorialSummary, "summary", "", "human-readable summary of the work")
	editorialFinishCmd.Flags().BoolVar(&editorialDryRun, "dry-run", false, "validate and preview without changing files")
	editorialCommitCmd.Flags().StringVar(&editorialSummary, "summary", "", "human-readable summary of the work")
	for _, c := range []*cobra.Command{editorialStartCmd, editorialFinishCmd, editorialCommitCmd} {
		c.Flags().BoolVar(&editorialJSON, "json", false, "print machine-readable output")
	}
}

func editorialRoot() string {
	if siteDirSelected {
		return activeSiteDir
	}
	if siteDir != "" {
		return siteDir
	}
	return "."
}
func editorialOptions() (editorial.Options, error) {
	primary := editorialPrimaryAuthor
	if cfg, err := config.Load(cfgFile); err != nil {
		return editorial.Options{}, fmt.Errorf("load config: %w", err)
	} else if primary == "" {
		for id := range cfg.Authors.Authors {
			if cfg.Authors.Authors[id].Default {
				primary = id
				break
			}
		}
	}
	return editorial.Options{Root: editorialRoot(), DefaultAuthor: primary, RecordedBy: editorialRecordedBy, DryRun: editorialDryRun}, nil
}
func printEditorial(out io.Writer, v any) error {
	if editorialJSON {
		result := v
		if s, ok := v.(editorial.Session); ok {
			paths := make([]string, 0, len(s.Posts))
			for p := range s.Posts {
				paths = append(paths, p)
			}
			sort.Strings(paths)
			result = struct {
				Session string            `json:"session"`
				Phase   string            `json:"phase"`
				Commits []string          `json:"commits,omitempty"`
				Paths   []string          `json:"paths"`
				Actions map[string]string `json:"actions,omitempty"`
				DryRun  bool              `json:"dry_run"`
			}{s.ID, s.Phase, s.Commits, paths, s.Actions, editorialDryRun}
		}
		b, e := json.MarshalIndent(result, "", "  ")
		if e != nil {
			return e
		}
		fmt.Fprintln(out, string(b))
		return nil
	}
	if editorialDryRun {
		if _, err := fmt.Fprintln(out, "Dry run: no files, commits, or session state changed."); err != nil {
			return err
		}
	}
	switch x := v.(type) {
	case editorial.Session:
		fmt.Fprintf(out, "session %s\n", x.ID)
		if len(x.Commits) > 0 {
			fmt.Fprintf(out, "commits: %s\n", strings.Join(x.Commits, " "))
		}
	default:
		fmt.Fprintln(out, v)
	}
	return nil
}
func runEditorialStart(cmd *cobra.Command, args []string) error {
	o, e := editorialOptions()
	if e != nil {
		return e
	}
	s, e := editorial.Start(o, editorialAuthor, editorialRole, args)
	if e != nil {
		return e
	}
	return printEditorial(cmd.OutOrStdout(), s)
}
func runEditorialFinish(cmd *cobra.Command, _ []string) error {
	if editorialSession == "" {
		return fmt.Errorf("--session is required")
	}
	o, e := editorialOptions()
	if e != nil {
		return e
	}
	s, e := editorial.Finish(o, editorialSession, editorialSummary)
	if e != nil {
		return e
	}
	return printEditorial(cmd.OutOrStdout(), s)
}
func runEditorialCommit(cmd *cobra.Command, args []string) error {
	o, e := editorialOptions()
	if e != nil {
		return e
	}
	s, e := editorial.Commit(o, editorialAuthor, editorialRole, editorialSummary, args)
	if e != nil {
		return e
	}
	return printEditorial(cmd.OutOrStdout(), s)
}
