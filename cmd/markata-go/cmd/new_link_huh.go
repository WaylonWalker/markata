package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func runHuhNewLinkDefaults(cmd *cobra.Command, state *newWizardState, ctx *siteContext, theme *huh.Theme) (bool, error) {
	urlGroup := huh.NewGroup(
		huh.NewInput().
			Title("Link URL").
			Description("Fetch defaults from the page you are sharing").
			Value(&state.LinkURL).
			Placeholder("https://example.com/article").
			Validate(validateNewLinkURL),
	)
	if err := huh.NewForm(urlGroup).WithTheme(theme).Run(); err != nil {
		return false, fmt.Errorf("wizard canceled: %w", err)
	}
	metadata, err := resolveNewLinkMetadata(cmd, state.LinkURL)
	if err != nil {
		return false, err
	}
	initializeNewLinkWizard(cmd, state, ctx.Templates[state.Template], metadata)
	if !metadata.TrustedTitle {
		titleGroup := huh.NewGroup(
			huh.NewInput().
				Title("Confirm Title").
				Description("No trusted Open Graph title. Edit or accept this suggestion.").
				Value(&state.Title).
				Validate(func(title string) error {
					if strings.TrimSpace(title) == "" {
						return fmt.Errorf("title is required")
					}
					return nil
				}),
		)
		if err := huh.NewForm(titleGroup).WithTheme(theme).Run(); err != nil {
			return false, fmt.Errorf("wizard canceled: %w", err)
		}
	}
	action := newLinkActionCreate
	actionGroup := huh.NewGroup(
		huh.NewNote().Title("Link Post").Description(buildSummary(state, ctx)),
		huh.NewSelect[string]().
			Title("Next").
			Options(
				huh.NewOption("Create", newLinkActionCreate),
				huh.NewOption("Edit details", newLinkActionEdit),
				huh.NewOption("Cancel", "cancel"),
			).
			Value(&action),
	)
	if err := huh.NewForm(actionGroup).WithTheme(theme).Run(); err != nil {
		return false, fmt.Errorf("wizard canceled: %w", err)
	}
	if action == "cancel" {
		return false, fmt.Errorf("canceled by user")
	}
	return action == newLinkActionCreate, nil
}

func initializeNewLinkWizard(cmd *cobra.Command, state *newWizardState, template ContentTemplate, metadata newLinkMetadata) {
	state.Title = metadata.Title
	state.Description = metadata.Description
	state.ImageURL = metadata.Image
	if state.Description == "" {
		if description, ok := template.Frontmatter["description"].(string); ok {
			state.Description = description
		}
	}
	if state.ImageURL == "" {
		if image, ok := template.Frontmatter["image"].(string); ok {
			state.ImageURL = image
		}
	}
	state.Directory = template.Directory
	state.DirChoice = template.Directory
	state.UseDefaultAuthor = true
	if cmd.Flags().Changed("dir") && newDir != "" {
		state.Directory = newDir
		state.DirChoice = newDir
	}
	if cmd.Flags().Changed("tags") {
		state.Tags = parseTags(newTags)
	}
}
