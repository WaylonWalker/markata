package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/charmbracelet/huh"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
)

// Template name constants used for conditional wizard groups.
const tmplPost = "post"
const privacyChoicePublic = string(encryptPostActionPublic)
const privacyChoicePrivate = "private"

const (
	tmplPhoto   = "photo"
	tmplVideo   = "video"
	tmplLink    = "link"
	tmplQuote   = "quote"
	tmplContact = "contact"
	tmplAuthor  = "author"

	dirCustomSentinel = "__custom__" // sentinel value for custom directory selection
)

// newWizardState holds all state gathered during the new post huh wizard.
type newWizardState struct {
	Title            string
	Template         string
	Directory        string
	DirChoice        string // intermediate: selected directory or __custom__
	CustomDir        string // intermediate: custom directory input
	Tags             []string
	CustomTag        string
	Private          bool   // legacy state; new wizard tracks an explicit privacy choice
	PrivateChoice    string // inherit, private, or public
	Authors          []string
	UseDefaultAuthor bool

	// Template-specific fields
	ImageURL string // photo, video (thumbnail), link (preview)
	VideoURL string // video
	Duration string // video
	LinkURL  string // link
	Quote    string // quote
	Source   string // quote
	Handle   string // contact
	Name     string // author
	Bio      string // author
	Email    string // author
	Role     string // author
}

// siteContext holds discovered site metadata used to populate wizard options.
type siteContext struct {
	Config             *models.Config
	Templates          map[string]ContentTemplate
	ExistingTags       []string
	ExistingDirs       []string
	HasMultipleAuthors bool
	AuthorOptions      []authorOption
	DefaultAuthorID    string
}

// authorOption represents an author choice in the wizard.
type authorOption struct {
	ID   string
	Name string
}

// loadSiteContext loads config, discovers tags, directories, and authors.
func loadSiteContext(templates map[string]ContentTemplate) *siteContext {
	ctx := &siteContext{
		Templates: templates,
	}

	// Load full config
	cfg, err := config.Load("")
	if err != nil {
		cfg, _ = config.LoadWithDefaults() //nolint:errcheck // best-effort fallback to defaults
	}
	ctx.Config = cfg

	// Discover existing tags from content
	ctx.ExistingTags = discoverExistingTags(cfg)

	// Discover existing directories
	ctx.ExistingDirs = discoverExistingDirs(cfg)

	// Discover author configuration
	if cfg != nil && len(cfg.Authors.Authors) > 1 {
		ctx.HasMultipleAuthors = true
		ctx.AuthorOptions = make([]authorOption, 0, len(cfg.Authors.Authors))

		for id := range cfg.Authors.Authors {
			author := cfg.Authors.Authors[id]
			if !author.Active && !author.Default {
				continue
			}
			ctx.AuthorOptions = append(ctx.AuthorOptions, authorOption{
				ID:   id,
				Name: author.Name,
			})
			if author.Default {
				ctx.DefaultAuthorID = id
			}
		}

		// Sort authors by name for consistent display
		sort.Slice(ctx.AuthorOptions, func(i, j int) bool {
			return ctx.AuthorOptions[i].Name < ctx.AuthorOptions[j].Name
		})

		// If no default was set, use first active author
		if ctx.DefaultAuthorID == "" && len(ctx.AuthorOptions) > 0 {
			ctx.DefaultAuthorID = ctx.AuthorOptions[0].ID
		}
	}

	return ctx
}

// discoverExistingTags scans content files to collect all unique tags.
func discoverExistingTags(cfg *models.Config) []string {
	if cfg == nil {
		return nil
	}

	patterns := cfg.GlobConfig.Patterns
	if len(patterns) == 0 {
		patterns = []string{"**/*.md"}
	}

	tagSet := make(map[string]bool)

	for _, pattern := range patterns {
		matches, err := doublestar.FilepathGlob(pattern)
		if err != nil {
			continue
		}
		for _, path := range matches {
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			metadata, _, err := plugins.ParseFrontmatter(string(content))
			if err != nil {
				continue
			}
			for _, tag := range plugins.GetStringSlice(metadata, "tags") {
				tagSet[tag] = true
			}
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

// discoverExistingDirs finds all directories that contain markdown files.
func discoverExistingDirs(cfg *models.Config) []string {
	if cfg == nil {
		return nil
	}

	patterns := cfg.GlobConfig.Patterns
	if len(patterns) == 0 {
		patterns = []string{"**/*.md"}
	}

	dirSet := make(map[string]bool)

	for _, pattern := range patterns {
		matches, err := doublestar.FilepathGlob(pattern)
		if err != nil {
			continue
		}
		for _, path := range matches {
			dir := filepath.Dir(path)
			if dir != "." {
				dirSet[dir] = true
			}
		}
	}

	dirs := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	return dirs
}

// runHuhNewWizard runs the huh-based interactive wizard for creating new content.
// All groups are in a single form so Shift+Tab navigates backwards.
func runHuhNewWizard(templates map[string]ContentTemplate) (*newWizardState, error) {
	ctx := loadSiteContext(templates)

	// Use the palette resolved for the active command configuration.
	theme := createHuhThemeForPalette(currentCLIPalette)

	state := &newWizardState{
		Template:      tmplPost,
		PrivateChoice: "inherit",
	}
	if _, resolvedName, ok := resolveTemplateSelection(newTemplate, templates); ok {
		state.Template = resolvedName
	}

	// Build all groups for a single form with back navigation.

	// --- Group 1: Template and Title ---
	templateOptions := buildTemplateOptions(ctx.Templates)
	titleGroup := huh.NewGroup(
		huh.NewNote().
			Title("Create New Content").
			Description("Let's set up your new content file.\nUse Shift+Tab to go back to a previous step."),
		huh.NewSelect[string]().
			Title("Template").
			Description("Choose the content type (/ to filter)").
			Options(templateOptions...).
			Value(&state.Template),
		huh.NewInput().
			Title("Title").
			Description("The title of your new content").
			Value(&state.Title).
			Placeholder("My New Post").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("title is required")
				}
				return nil
			}),
	)

	// --- Group 2: Directory selection (dynamic based on template) ---
	dirGroup := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Directory").
			Description("Where to create the file (/ to filter)").
			OptionsFunc(func() []huh.Option[string] {
				selectedTemplate := ctx.Templates[state.Template]
				return buildDirectoryOptions(selectedTemplate, ctx.ExistingDirs)
			}, &state.Template).
			Value(&state.DirChoice),
	)

	// --- Group 2b: Custom directory (shown only if __custom__ selected) ---
	customDirGroup := huh.NewGroup(
		huh.NewInput().
			Title("Custom Directory").
			Description("Enter the directory path").
			Value(&state.CustomDir).
			Placeholder("pages/custom"),
	).WithHideFunc(func() bool {
		return state.DirChoice != dirCustomSentinel
	})

	// --- Template-specific groups ---

	// Photo: image URL
	photoGroup := huh.NewGroup(
		huh.NewInput().
			Title("Image URL").
			Description("Path or URL to the image").
			Value(&state.ImageURL).
			Placeholder("/images/photo.jpg"),
	).WithHideFunc(func() bool {
		return state.Template != tmplPhoto
	})

	// Video: video URL, thumbnail, duration
	videoGroup := huh.NewGroup(
		huh.NewInput().
			Title("Video URL").
			Description("Path or URL to the video").
			Value(&state.VideoURL).
			Placeholder("https://youtube.com/watch?v=..."),
		huh.NewInput().
			Title("Thumbnail Image").
			Description("Path or URL to the video thumbnail (optional)").
			Value(&state.ImageURL).
			Placeholder("/images/thumb.jpg"),
		huh.NewInput().
			Title("Duration").
			Description("Video duration (optional)").
			Value(&state.Duration).
			Placeholder("5:30"),
	).WithHideFunc(func() bool {
		return state.Template != tmplVideo
	})

	// Link: URL and preview image
	linkGroup := huh.NewGroup(
		huh.NewInput().
			Title("Link URL").
			Description("The URL you are sharing").
			Value(&state.LinkURL).
			Placeholder("https://example.com/article"),
		huh.NewInput().
			Title("Preview Image").
			Description("Image for the link preview (optional)").
			Value(&state.ImageURL).
			Placeholder("/images/preview.jpg"),
	).WithHideFunc(func() bool {
		return state.Template != tmplLink
	})

	// Quote: quote text and source
	quoteGroup := huh.NewGroup(
		huh.NewInput().
			Title("Quote").
			Description("The quote text").
			Value(&state.Quote).
			Placeholder("To be or not to be..."),
		huh.NewInput().
			Title("Source").
			Description("Where the quote is from (optional)").
			Value(&state.Source).
			Placeholder("Shakespeare, Hamlet"),
	).WithHideFunc(func() bool {
		return state.Template != tmplQuote
	})

	// Contact: handle
	contactGroup := huh.NewGroup(
		huh.NewInput().
			Title("Handle").
			Description("Short handle for @mentions (e.g. alice)").
			Value(&state.Handle).
			Placeholder("alice"),
	).WithHideFunc(func() bool {
		return state.Template != tmplContact
	})

	// Author: name, bio, email, role
	authorGroup := huh.NewGroup(
		huh.NewInput().
			Title("Author Name").
			Description("Full name of the author").
			Value(&state.Name).
			Placeholder("Jane Doe"),
		huh.NewInput().
			Title("Bio").
			Description("Short bio (optional)").
			Value(&state.Bio).
			Placeholder("Software engineer and writer"),
		huh.NewInput().
			Title("Email").
			Description("Contact email (optional)").
			Value(&state.Email).
			Placeholder("jane@example.com"),
		huh.NewInput().
			Title("Role").
			Description("Role or title (optional)").
			Value(&state.Role).
			Placeholder("author"),
	).WithHideFunc(func() bool {
		return state.Template != tmplAuthor
	})

	// --- Group 3: Tags ---
	var tagFields []huh.Field

	if len(ctx.ExistingTags) > 0 {
		tagOptions := make([]huh.Option[string], 0, len(ctx.ExistingTags))
		for _, tag := range ctx.ExistingTags {
			tagOptions = append(tagOptions, huh.NewOption(tag, tag))
		}
		tagFields = append(tagFields,
			huh.NewMultiSelect[string]().
				Title("Tags").
				Description("Select from existing tags (/ to filter, space to select)").
				Options(tagOptions...).
				Filterable(true).
				Value(&state.Tags).
				Height(min(len(tagOptions)+1, 12)),
		)
	}

	tagFields = append(tagFields,
		huh.NewInput().
			Title("Additional Tags").
			Description("Enter additional tags, comma-separated (or leave blank)").
			Value(&state.CustomTag).
			Placeholder("new-tag, another-tag"),
	)

	tagsGroup := huh.NewGroup(tagFields...)

	// --- Group 4: Privacy ---
	// Inheriting is distinct from explicitly selecting public (false).
	privateGroup := huh.NewGroup(
		huh.NewSelect[string]().
			Title("Privacy").
			Description("Keep the starter setting, or explicitly override it").
			Options(
				huh.NewOption("Keep starter default", "inherit"),
				huh.NewOption("Private", privacyChoicePrivate),
				huh.NewOption("Public", privacyChoicePublic),
			).
			Value(&state.PrivateChoice),
	)

	// --- Group 5: Authors (conditional - only for multi-author sites) ---
	// Default author confirm
	defaultAuthorName := ctx.DefaultAuthorID
	for _, a := range ctx.AuthorOptions {
		if a.ID == ctx.DefaultAuthorID {
			defaultAuthorName = a.Name
			break
		}
	}

	defaultAuthorGroup := huh.NewGroup(
		huh.NewConfirm().
			Title(fmt.Sprintf("Use default author (%s)?", defaultAuthorName)).
			Description("You can add multiple authors if needed").
			Value(&state.UseDefaultAuthor).
			Affirmative("Yes").
			Negative("Choose authors"),
	).WithHideFunc(func() bool {
		return !ctx.HasMultipleAuthors
	})

	// Author multi-select (shown only if user chose to pick authors)
	authorOpts := buildAuthorOptions(ctx)
	authorSelectGroup := huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("Authors").
			Description("Select authors for this post (/ to filter, space to select)").
			Options(authorOpts...).
			Filterable(true).
			Value(&state.Authors),
	).WithHideFunc(func() bool {
		return !ctx.HasMultipleAuthors || state.UseDefaultAuthor
	})

	// --- Group 6: Summary and confirm ---
	var confirmed bool
	summaryGroup := huh.NewGroup(
		huh.NewNote().
			Title("Summary").
			DescriptionFunc(func() string {
				return buildSummary(state, ctx)
			}, state),
		huh.NewConfirm().
			Title("Create this file?").
			Value(&confirmed).
			Affirmative("Create").
			Negative("Cancel"),
	)

	// --- Single form: all groups, Shift+Tab goes back ---
	form := huh.NewForm(
		titleGroup,
		dirGroup,
		customDirGroup,
		photoGroup,
		videoGroup,
		linkGroup,
		quoteGroup,
		contactGroup,
		authorGroup,
		tagsGroup,
		privateGroup,
		defaultAuthorGroup,
		authorSelectGroup,
		summaryGroup,
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		return nil, fmt.Errorf("wizard canceled: %w", err)
	}

	if !confirmed {
		return nil, fmt.Errorf("canceled by user")
	}

	// Resolve directory
	if state.DirChoice == dirCustomSentinel {
		if state.CustomDir != "" {
			state.Directory = state.CustomDir
		} else {
			state.Directory = ctx.Templates[state.Template].Directory
		}
	} else {
		state.Directory = state.DirChoice
	}

	// Merge custom tags
	if state.CustomTag != "" {
		for _, tag := range parseTags(state.CustomTag) {
			if !containsString(state.Tags, tag) {
				state.Tags = append(state.Tags, tag)
			}
		}
	}

	// Resolve default author
	if ctx.HasMultipleAuthors && state.UseDefaultAuthor {
		state.Authors = []string{ctx.DefaultAuthorID}
	}

	return state, nil
}

// buildSummary generates the summary text for the confirmation step.
func buildSummary(state *newWizardState, ctx *siteContext) string {
	slug := generateSlug(state.Title)
	filename := slug + ".md"

	dir := state.DirChoice
	if dir == dirCustomSentinel && state.CustomDir != "" {
		dir = state.CustomDir
	} else if dir == dirCustomSentinel {
		dir = ctx.Templates[state.Template].Directory
	}
	fullPath := filepath.Join(dir, filename)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Title:     %s\n", state.Title))
	sb.WriteString(fmt.Sprintf("Template:  %s\n", state.Template))
	sb.WriteString(fmt.Sprintf("File:      %s\n", fullPath))
	selectedTemplate := ctx.Templates[state.Template]
	effectiveDraft := starterDraftDefault(selectedTemplate)
	sb.WriteString(fmt.Sprintf("Draft:     %t (effective)\n", effectiveDraft))
	sb.WriteString(fmt.Sprintf("Published: %t (effective)\n", !effectiveDraft))
	effectivePrivate := false
	if inheritedPrivacy, ok := selectedTemplate.Frontmatter["private"].(bool); ok {
		effectivePrivate = inheritedPrivacy
	}
	switch state.PrivateChoice {
	case privacyChoicePrivate:
		effectivePrivate = true
	case privacyChoicePublic:
		effectivePrivate = false
	}
	sb.WriteString(fmt.Sprintf("Private:   %t (effective)\n", effectivePrivate))

	// Template-specific fields
	sb.WriteString(buildTemplateSummary(state))

	// Merge selected + custom tags for display
	allTags := append([]string{}, state.Tags...)
	if state.CustomTag != "" {
		for _, tag := range parseTags(state.CustomTag) {
			if !containsString(allTags, tag) {
				allTags = append(allTags, tag)
			}
		}
	}
	if len(allTags) > 0 {
		sb.WriteString(fmt.Sprintf("Tags:      %s\n", strings.Join(allTags, ", ")))
	}
	if ctx.HasMultipleAuthors {
		if state.UseDefaultAuthor {
			sb.WriteString(fmt.Sprintf("Author:    %s (default)\n", ctx.DefaultAuthorID))
		} else if len(state.Authors) > 0 {
			sb.WriteString(fmt.Sprintf("Authors:   %s\n", strings.Join(state.Authors, ", ")))
		}
	}

	return sb.String()
}

// buildTemplateOptions creates huh options from available templates.
func buildTemplateOptions(templates map[string]ContentTemplate) []huh.Option[string] {
	// Sort template names for consistent display
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	options := make([]huh.Option[string], 0, len(names))
	for _, name := range names {
		t := templates[name]
		label := formatTemplateEntry(name, t)
		options = append(options, huh.NewOption(label, name))
	}
	return options
}

// buildDirectoryOptions creates huh options for directory selection.
func buildDirectoryOptions(template ContentTemplate, existingDirs []string) []huh.Option[string] {
	seen := make(map[string]bool)
	var options []huh.Option[string]

	// First option: template's default directory
	defaultDir := template.Directory
	label := fmt.Sprintf("%s (default)", defaultDir)
	options = append(options, huh.NewOption(label, defaultDir))
	seen[defaultDir] = true

	// Add existing directories from the site
	for _, dir := range existingDirs {
		if !seen[dir] {
			options = append(options, huh.NewOption(dir, dir))
			seen[dir] = true
		}
	}

	// Custom option
	options = append(options, huh.NewOption("Custom...", dirCustomSentinel))

	return options
}

// buildAuthorOptions creates huh options for author multi-select.
func buildAuthorOptions(ctx *siteContext) []huh.Option[string] {
	opts := make([]huh.Option[string], 0, len(ctx.AuthorOptions))
	for _, a := range ctx.AuthorOptions {
		label := a.Name
		if a.ID == ctx.DefaultAuthorID {
			label = fmt.Sprintf("%s (default)", a.Name)
		}
		opts = append(opts, huh.NewOption(label, a.ID))
	}
	return opts
}

// buildTemplateSummary returns summary lines for template-specific fields.
func buildTemplateSummary(state *newWizardState) string {
	var sb strings.Builder

	switch state.Template {
	case tmplPhoto:
		if state.ImageURL != "" {
			sb.WriteString(fmt.Sprintf("Image:     %s\n", state.ImageURL))
		}
	case tmplVideo:
		if state.VideoURL != "" {
			sb.WriteString(fmt.Sprintf("Video:     %s\n", state.VideoURL))
		}
		if state.ImageURL != "" {
			sb.WriteString(fmt.Sprintf("Thumbnail: %s\n", state.ImageURL))
		}
		if state.Duration != "" {
			sb.WriteString(fmt.Sprintf("Duration:  %s\n", state.Duration))
		}
	case tmplLink:
		if state.LinkURL != "" {
			sb.WriteString(fmt.Sprintf("URL:       %s\n", state.LinkURL))
		}
	case tmplQuote:
		if state.Quote != "" {
			sb.WriteString(fmt.Sprintf("Quote:     %s\n", state.Quote))
		}
		if state.Source != "" {
			sb.WriteString(fmt.Sprintf("Source:    %s\n", state.Source))
		}
	case tmplContact:
		if state.Handle != "" {
			sb.WriteString(fmt.Sprintf("Handle:    @%s\n", state.Handle))
		}
	case tmplAuthor:
		if state.Name != "" {
			sb.WriteString(fmt.Sprintf("Name:      %s\n", state.Name))
		}
		if state.Role != "" {
			sb.WriteString(fmt.Sprintf("Role:      %s\n", state.Role))
		}
	}

	return sb.String()
}

// applyTemplateFields sets template-specific frontmatter fields from wizard state.
func applyTemplateFields(state *newWizardState, fm map[string]interface{}) {
	switch state.Template {
	case tmplPhoto:
		if state.ImageURL != "" {
			fm["image"] = state.ImageURL
		}
	case tmplVideo:
		if state.VideoURL != "" {
			fm["video"] = state.VideoURL
		}
		if state.ImageURL != "" {
			fm["image"] = state.ImageURL
		}
		if state.Duration != "" {
			fm["duration"] = state.Duration
		}
	case tmplLink:
		if state.LinkURL != "" {
			fm["url"] = state.LinkURL
		}
		if state.ImageURL != "" {
			fm["image"] = state.ImageURL
		}
	case tmplQuote:
		if state.Quote != "" {
			fm["quote"] = state.Quote
		}
		if state.Source != "" {
			fm["source"] = state.Source
		}
	case tmplContact:
		if state.Handle != "" {
			fm["handle"] = state.Handle
		}
	case tmplAuthor:
		if state.Name != "" {
			fm["name"] = state.Name
		}
		if state.Bio != "" {
			fm["bio"] = state.Bio
		}
		if state.Email != "" {
			fm["email"] = state.Email
		}
		if state.Role != "" {
			fm["role"] = state.Role
		}
	}
}

// containsString checks if a string slice contains a value.
func containsString(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}

// applyNewWizardState generates and writes the content file from wizard state.
func applyNewWizardState(state *newWizardState, templates map[string]ContentTemplate) error {
	template, _, ok := resolveTemplateSelection(state.Template, templates)
	if !ok {
		return fmt.Errorf("unknown template %q", state.Template)
	}
	if !safeStarterDirectory(template.Directory) {
		return fmt.Errorf("starter %q has an unsafe output directory %q", template.Name, template.Directory)
	}

	// Override directory from wizard state
	outputDir := state.Directory
	if outputDir == "" {
		outputDir = template.Directory
	}
	if !safeStarterDirectory(outputDir) {
		return fmt.Errorf("unsafe output directory %q", outputDir)
	}

	slug := generateSlug(state.Title)

	// Build enhanced template that includes wizard-gathered fields
	enhancedTemplate := template
	if enhancedTemplate.Frontmatter == nil {
		enhancedTemplate.Frontmatter = make(map[string]interface{})
	}

	// Inherited values stay in the starter; only explicit choices are written.
	switch state.PrivateChoice {
	case privacyChoicePrivate:
		enhancedTemplate.Frontmatter["private"] = true
	case privacyChoicePublic:
		enhancedTemplate.Frontmatter["private"] = false
	}

	// Set authors if provided
	if len(state.Authors) > 0 {
		enhancedTemplate.Frontmatter["authors"] = state.Authors
	}

	// Apply template-specific fields from wizard
	applyTemplateFields(state, enhancedTemplate.Frontmatter)

	// An unmodified tag selector is not an explicit empty-tags override.
	// Preserve the starter's tags unless the user supplied replacements.
	tags := state.Tags
	if len(tags) == 0 {
		tags = nil
	}

	// The wizard has no publication control. Preserve starter/site defaults.
	return writeContentFile(state.Title, slug, outputDir, false, false, tags, enhancedTemplate)
}
