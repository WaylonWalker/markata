package cmd

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/WaylonWalker/markata/pkg/agentskill"
)

// rootSkillRequested makes the bundled, release-matched first-run instructions
// available on stdout without requiring a site or an agent installation.
var rootSkillRequested bool

func init() {
	rootCmd.Flags().BoolVar(&rootSkillRequested, "skill", false, "print the bundled, agent-first setup guide to stdout")
}

func printFirstRunSkill(out io.Writer) error {
	bundle, err := agentskill.SiteSkill()
	if err != nil {
		return fmt.Errorf("open bundled agent skill: %w", err)
	}
	content, err := fs.ReadFile(bundle, "topics/first-setup.md")
	if err != nil {
		return fmt.Errorf("read bundled first-setup skill: %w", err)
	}
	if _, err := out.Write(content); err != nil {
		return fmt.Errorf("print bundled first-setup skill: %w", err)
	}
	return nil
}
