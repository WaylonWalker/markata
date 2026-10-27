package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/themes"
)

// initCmd represents the init command.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new markata-go project",
	Long: `Initialize a new markata-go project with interactive setup.

This command provides a rich TUI (Text User Interface) wizard for
project initialization. The wizard guides you through:

  - Setting site information (title, description, author, URL)
  - Configuring optional features (themes, SEO, post formats, feeds)
  - Vending built-in assets for customization
  - Creating your first post

The wizard automatically falls back to plain text prompts when stdin
is not a TTY. Use --plain to force plain text mode.

If a configuration file already exists, you can add new features or
update site information interactively.

Example usage:
  markata-go init           # Interactive project setup (TUI)
  markata-go init --plain   # Use plain text prompts
  markata-go init --force   # Overwrite existing files`,
	RunE: runInitCommand,
}

var (
	// initForce overwrites existing files.
	initForce bool
	// initPlain uses plain text prompts instead of the huh TUI wizard.
	initPlain bool
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initForce, "force", false, "overwrite existing files")
	initCmd.Flags().BoolVar(&initPlain, "plain", false, "use plain text prompts (for non-TTY environments)")
}

// prompt displays a question and returns the user's response or a default value.
func prompt(reader *bufio.Reader, question, defaultVal string) string {
	if defaultVal != "" {
		errf("%s [%s]: ", question, defaultVal)
	} else {
		errf("%s: ", question)
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultVal
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

// Feature name constants for the init wizard.
const (
	featureTheme        = "theme"
	featureSEO          = "seo"
	featurePostFormats  = "post_formats"
	featureAdvancedFeed = "advanced_feeds"
)

// defaultConfigFilename is the default name for the configuration file.
const defaultConfigFilename = "markata-go.toml"

// featureInfo describes an available feature for the wizard.
type featureInfo struct {
	Name        string
	Description string
	Configured  bool
}

// detectConfiguredFeatures checks which features are already configured.
func detectConfiguredFeatures(cfg *models.Config) map[string]bool {
	features := make(map[string]bool)
	features[featureTheme] = cfg.Theme.Palette != "" && cfg.Theme.Palette != "default-light"
	features[featureSEO] = cfg.SEO.TwitterHandle != "" || cfg.SEO.DefaultImage != ""
	features[featurePostFormats] = cfg.PostFormats.Markdown || cfg.PostFormats.Text || cfg.PostFormats.OG || !cfg.PostFormats.IsHTMLEnabled()
	features[featureAdvancedFeed] = cfg.FeedDefaults.Formats.Atom || cfg.FeedDefaults.Formats.JSON || cfg.FeedDefaults.Formats.Sitemap
	return features
}

// getAvailableFeatures returns the list of features that can be added.
func getAvailableFeatures(configured map[string]bool) []featureInfo {
	allFeatures := []featureInfo{
		{
			Name:        featureTheme,
			Description: "Theme/Palette system (color schemes)",
			Configured:  configured[featureTheme],
		},
		{
			Name:        featureSEO,
			Description: "SEO metadata (Twitter, Open Graph)",
			Configured:  configured[featureSEO],
		},
		{
			Name:        featurePostFormats,
			Description: "Post output formats (HTML, markdown, text, OG cards)",
			Configured:  configured[featurePostFormats],
		},
		{
			Name:        featureAdvancedFeed,
			Description: "Feed formats (HTML, RSS, Atom, JSON, sitemap)",
			Configured:  configured[featureAdvancedFeed],
		},
	}

	return allFeatures
}

// promptMenuChoice displays a numbered menu and returns the selected option.
func promptMenuChoice(reader *bufio.Reader, question string, options []string) int {
	errln()
	errln(question)
	for i, opt := range options {
		errlnf("  %d) %s", i+1, opt)
	}
	errf("\nEnter choice [1]: ")

	input, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(input) == "" {
		return 0
	}

	var choice int
	_, err = fmt.Sscanf(strings.TrimSpace(input), "%d", &choice)
	if err != nil || choice < 1 || choice > len(options) {
		return 0
	}
	return choice - 1
}

// promptMenuChoiceDefault displays a numbered menu and returns the selected option.
// Uses defaultIndex when input is empty or invalid.
func promptMenuChoiceDefault(reader *bufio.Reader, question string, options []string, defaultIndex int) int {
	if defaultIndex < 0 || defaultIndex >= len(options) {
		defaultIndex = 0
	}

	errln()
	errln(question)
	for i, opt := range options {
		errlnf("  %d) %s", i+1, opt)
	}
	errf("\nEnter choice [%d]: ", defaultIndex+1)

	input, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(input) == "" {
		return defaultIndex
	}

	var choice int
	_, err = fmt.Sscanf(strings.TrimSpace(input), "%d", &choice)
	if err != nil || choice < 1 || choice > len(options) {
		return defaultIndex
	}
	return choice - 1
}

// promptRadioBool displays an enabled/disabled choice and returns the selected value.
func promptRadioBool(reader *bufio.Reader, question string, defaultEnabled bool) bool {
	options := []string{"Enabled", "Disabled"}
	defaultIndex := 1
	if defaultEnabled {
		defaultIndex = 0
	}
	choice := promptMenuChoiceDefault(reader, question, options, defaultIndex)
	return choice == 0
}

// promptFeatureSelection displays checkboxes for feature selection.
func promptFeatureSelection(reader *bufio.Reader, features []featureInfo) []string {
	errln()
	errln("Select features to add (enter numbers separated by spaces):")
	errln()

	availableIdx := []int{}
	for i, f := range features {
		status := "[ ]"
		if f.Configured {
			status = "[x]"
			errlnf("  %d) %s %s (already configured)", i+1, status, f.Description)
		} else {
			errlnf("  %d) %s %s", i+1, status, f.Description)
			availableIdx = append(availableIdx, i)
		}
	}

	if len(availableIdx) == 0 {
		errln("\n  All features are already configured!")
		return nil
	}

	errf("\nEnter numbers (e.g., 1 3): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil
	}

	selected := []string{}
	parts := strings.Fields(strings.TrimSpace(input))
	for _, p := range parts {
		var idx int
		if _, err := fmt.Sscanf(p, "%d", &idx); err == nil {
			if idx >= 1 && idx <= len(features) && !features[idx-1].Configured {
				selected = append(selected, features[idx-1].Name)
			}
		}
	}

	return selected
}

// backupConfig creates a backup of the existing config file.
func backupConfig(path string) error {
	backupPath := path + ".backup"
	// If backup already exists, add timestamp
	if _, err := os.Stat(backupPath); err == nil {
		backupPath = path + "." + time.Now().Format("20060102-150405") + ".backup"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config for backup: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0o644); err != nil { //nolint:gosec // backup files should be readable
		return fmt.Errorf("failed to write backup: %w", err)
	}

	outlnf("  Backed up config to %s", backupPath)
	return nil
}

// addFeatureTheme prompts for theme/palette configuration.
func addFeatureTheme(reader *bufio.Reader, cfg *models.Config) error {
	outln()
	cliSection("Theme/Palette Configuration")

	// List some available palettes
	loader := palettes.NewLoader()
	availablePalettes, err := loader.Discover()
	if err != nil {
		outln("(Could not discover palettes, using default)")
	} else {
		outln("Available palettes:")
		// Show a sample of palettes
		shown := 0
		for _, p := range availablePalettes {
			if shown < 10 {
				outlnf("  - %s (%s)", p.Name, p.Variant)
				shown++
			}
		}
		if len(availablePalettes) > 10 {
			outlnf("  ... and %d more (run 'markata-go palette list' to see all)", len(availablePalettes)-10)
		}
	}

	palette := prompt(reader, "\nPalette name", "default-light")
	cfg.Theme.Palette = palette

	return nil
}

// addFeatureSEO prompts for SEO configuration.
func addFeatureSEO(reader *bufio.Reader, cfg *models.Config) error {
	outln()
	cliSection("SEO Configuration")

	handle := prompt(reader, "Twitter/X handle (without @)", cfg.SEO.TwitterHandle)
	if handle != "" {
		// Remove @ if user accidentally includes it
		handle = strings.TrimPrefix(handle, "@")
		cfg.SEO.TwitterHandle = handle
	}

	defaultImage := prompt(reader, "Default Open Graph image URL", cfg.SEO.DefaultImage)
	cfg.SEO.DefaultImage = defaultImage

	logoURL := prompt(reader, "Site logo URL (for Schema.org)", cfg.SEO.LogoURL)
	cfg.SEO.LogoURL = logoURL

	return nil
}

// addFeaturePostFormats prompts for post format configuration.
func addFeaturePostFormats(reader *bufio.Reader, cfg *models.Config) error {
	outln()
	cliSection("Post Output Formats")

	htmlEnabled := promptRadioBool(reader, "HTML output", cfg.PostFormats.IsHTMLEnabled())
	cfg.PostFormats.HTML = &htmlEnabled

	cfg.PostFormats.Markdown = promptRadioBool(reader, "Markdown source output (generates /slug.md)", cfg.PostFormats.Markdown)
	cfg.PostFormats.Text = promptRadioBool(reader, "Text output (generates /slug.txt)", cfg.PostFormats.Text)
	cfg.PostFormats.OG = promptRadioBool(reader, "OG card output (generates /slug/og/index.html for social images)", cfg.PostFormats.OG)

	return nil
}

// addFeatureAdvancedFeeds prompts for advanced feed configuration.
func addFeatureAdvancedFeeds(reader *bufio.Reader, cfg *models.Config) error {
	outln()
	cliSection("Advanced Feed Formats")

	cfg.FeedDefaults.Formats.HTML = promptRadioBool(reader, "HTML feed output", cfg.FeedDefaults.Formats.HTML)
	cfg.FeedDefaults.Formats.RSS = promptRadioBool(reader, "RSS feed output", cfg.FeedDefaults.Formats.RSS)
	cfg.FeedDefaults.Formats.Atom = promptRadioBool(reader, "Atom feed output", cfg.FeedDefaults.Formats.Atom)
	cfg.FeedDefaults.Formats.JSON = promptRadioBool(reader, "JSON Feed output", cfg.FeedDefaults.Formats.JSON)
	cfg.FeedDefaults.Formats.Sitemap = promptRadioBool(reader, "Sitemap output", cfg.FeedDefaults.Formats.Sitemap)

	return nil
}

// addFeature adds a specific feature to the configuration.
func addFeature(reader *bufio.Reader, feature string, cfg *models.Config) error {
	switch feature {
	case featureTheme:
		return addFeatureTheme(reader, cfg)
	case featureSEO:
		return addFeatureSEO(reader, cfg)
	case featurePostFormats:
		return addFeaturePostFormats(reader, cfg)
	case featureAdvancedFeed:
		return addFeatureAdvancedFeeds(reader, cfg)
	}
	return nil
}

// displayCurrentConfig shows the current configuration.
func displayCurrentConfig(cfg *models.Config) {
	outln()
	cliSection("Current Configuration")
	outln()
	outln("Site Information:")
	outlnf("  Title:       %s", cfg.Title)
	outlnf("  URL:         %s", cfg.URL)
	outlnf("  Description: %s", cfg.Description)
	outlnf("  Author:      %s", cfg.Author)
	outln()
	outln("Directories:")
	outlnf("  Output:      %s", cfg.OutputDir)
	outlnf("  Templates:   %s", cfg.TemplatesDir)
	outlnf("  Assets:      %s", cfg.AssetsDir)
	outln()
	outln("Theme:")
	outlnf("  Palette:     %s", cfg.Theme.Palette)
	outln()
	outln("SEO:")
	outlnf("  Twitter:     %s", valueOrNone(cfg.SEO.TwitterHandle))
	outlnf("  Default Img: %s", valueOrNone(cfg.SEO.DefaultImage))
	outln()
	outln("Post Formats:")
	outlnf("  HTML:        %v", cfg.PostFormats.IsHTMLEnabled())
	outlnf("  Markdown:    %v", cfg.PostFormats.Markdown)
	outlnf("  Text:        %v", cfg.PostFormats.Text)
	outlnf("  OG Cards:    %v", cfg.PostFormats.OG)
	outln()
	outln("Feed Formats (defaults):")
	outlnf("  HTML:        %v", cfg.FeedDefaults.Formats.HTML)
	outlnf("  RSS:         %v", cfg.FeedDefaults.Formats.RSS)
	outlnf("  Atom:        %v", cfg.FeedDefaults.Formats.Atom)
	outlnf("  JSON:        %v", cfg.FeedDefaults.Formats.JSON)
	outlnf("  Sitemap:     %v", cfg.FeedDefaults.Formats.Sitemap)
	outln()
}

func valueOrNone(s string) string {
	if s == "" {
		return "(not set)"
	}
	return s
}

// =============================================================================
// Vending Functions - Export built-in assets for local customization
// =============================================================================

// venderContentTemplates exports built-in content templates to content-templates/.
// Each template is written as a markdown file with YAML frontmatter.
// A README.md is also generated explaining how to use and customize templates.
func venderContentTemplates(force bool) error {
	dir := "content-templates"

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	templates := BuiltinTemplates()

	// Sort template names for consistent output
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	written := 0
	skipped := 0
	for _, name := range names {
		t := templates[name]
		filename := filepath.Join(dir, name+".md")

		// Check if file exists
		if _, err := os.Stat(filename); err == nil && !force {
			outlnf("  ! Skipped %s (exists, use --force to overwrite)", filename)
			skipped++
			continue
		}

		content := t.ToMarkdown()
		if err := os.WriteFile(filename, []byte(content), 0o644); err != nil { //nolint:gosec // template files should be readable
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
		outlnf("  Created %s", filename)
		written++
	}

	// Write README
	readmePath := filepath.Join(dir, "README.md")
	if _, err := os.Stat(readmePath); err != nil || force {
		readme := generateContentTemplatesREADME(templates)
		if err := os.WriteFile(readmePath, []byte(readme), 0o644); err != nil { //nolint:gosec // readme should be readable
			return fmt.Errorf("failed to write README: %w", err)
		}
		outlnf("  Created %s", readmePath)
	}

	outlnf("\n  Vendored %d content templates (%d skipped)", written, skipped)
	return nil
}

// venderPalettes exports built-in palettes to palettes/.
// Each palette is written as a TOML file.
// A README.md is also generated explaining how to customize palettes.
func venderPalettes(force bool) error {
	dir := "palettes"

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Get all built-in palette files
	paletteFiles, err := palettes.ListBuiltinFiles()
	if err != nil {
		return fmt.Errorf("failed to list built-in palettes: %w", err)
	}

	written := 0
	skipped := 0
	for filename, data := range paletteFiles {
		fullPath := filepath.Join(dir, filename)

		// Check if file exists
		if _, err := os.Stat(fullPath); err == nil && !force {
			outlnf("  ! Skipped %s (exists, use --force to overwrite)", fullPath)
			skipped++
			continue
		}

		if err := os.WriteFile(fullPath, data, 0o644); err != nil { //nolint:gosec // palette files should be readable
			return fmt.Errorf("failed to write %s: %w", fullPath, err)
		}
		outlnf("  Created %s", fullPath)
		written++
	}

	// Write README
	readmePath := filepath.Join(dir, "README.md")
	if _, err := os.Stat(readmePath); err != nil || force {
		readme := generatePalettesREADME()
		if err := os.WriteFile(readmePath, []byte(readme), 0o644); err != nil { //nolint:gosec // readme should be readable
			return fmt.Errorf("failed to write README: %w", err)
		}
		outlnf("  Created %s", readmePath)
	}

	outlnf("\n  Vendored %d palettes (%d skipped)", written, skipped)
	return nil
}

// venderLayouts exports built-in HTML templates to templates/.
// This includes base templates, partials, and components.
// A README.md is also generated explaining template customization.
func venderLayouts(force bool) error {
	dir := defaultTemplatesDir

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Get all template files
	templateFiles, err := themes.ListTemplates()
	if err != nil {
		return fmt.Errorf("failed to list templates: %w", err)
	}

	written := 0
	skipped := 0
	for _, relPath := range templateFiles {
		fullPath := filepath.Join(dir, relPath)

		// Create subdirectories if needed
		subdir := filepath.Dir(fullPath)
		if err := os.MkdirAll(subdir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", subdir, err)
		}

		// Check if file exists
		if _, err := os.Stat(fullPath); err == nil && !force {
			outlnf("  ! Skipped %s (exists, use --force to overwrite)", fullPath)
			skipped++
			continue
		}

		// Read from embedded FS
		data, err := themes.ReadTemplate(relPath)
		if err != nil {
			return fmt.Errorf("failed to read template %s: %w", relPath, err)
		}

		if err := os.WriteFile(fullPath, data, 0o644); err != nil { //nolint:gosec // template files should be readable
			return fmt.Errorf("failed to write %s: %w", fullPath, err)
		}
		outlnf("  Created %s", fullPath)
		written++
	}

	// Write README
	readmePath := filepath.Join(dir, "README.md")
	if _, err := os.Stat(readmePath); err != nil || force {
		readme := generateLayoutsREADME()
		if err := os.WriteFile(readmePath, []byte(readme), 0o644); err != nil { //nolint:gosec // readme should be readable
			return fmt.Errorf("failed to write README: %w", err)
		}
		outlnf("  Created %s", readmePath)
	}

	outlnf("\n  Vendored %d layout templates (%d skipped)", written, skipped)
	return nil
}

// generateContentTemplatesREADME creates documentation for content templates.
func generateContentTemplatesREADME(templates map[string]ContentTemplate) string {
	var sb strings.Builder

	sb.WriteString(`# Content Templates

This directory contains content templates that control how new content files are created
when using ` + "`markata-go new`" + `.

## How Templates Work

Each ` + "`.md`" + ` file in this directory defines a content template:
- The filename (without .md) becomes the template name
- The frontmatter defines default fields for new content
- The body provides the starting content

## Special Frontmatter Fields

| Field | Purpose |
|-------|---------|
| ` + "`_directory`" + ` | Default output directory for new content using this template |
| ` + "`template`" + ` | HTML template to use for rendering (e.g., "post", "docs") |

## Available Templates

`)

	// Sort template names
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		t := templates[name]
		sb.WriteString(fmt.Sprintf("### %s\n\n", name))
		sb.WriteString(fmt.Sprintf("- **Directory:** `%s/`\n", t.Directory))
		if tpl, ok := t.Frontmatter["template"].(string); ok {
			sb.WriteString(fmt.Sprintf("- **HTML Template:** `%s`\n", tpl))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(`## Creating Custom Templates

1. Create a new ` + "`.md`" + ` file in this directory
2. Add frontmatter with your default fields
3. Include ` + "`_directory`" + ` to set where new content goes
4. Add any body content as a starting point

Example custom template (` + "`recipe.md`" + `):

` + "```yaml" + `
---
_directory: recipes
template: recipe
servings: 4
prep_time: ""
cook_time: ""
---

## Ingredients

-

## Instructions

1.
` + "```" + `

## Usage

` + "```bash" + `
markata-go new "My Recipe" --template recipe
` + "```" + `

This creates ` + "`recipes/my-recipe.md`" + ` with the template's frontmatter and body.

## Obsidian Compatibility

These templates work well with Obsidian's template plugin. Simply point Obsidian
to this directory as your templates folder.
`)

	return sb.String()
}

// generatePalettesREADME creates documentation for palettes.
func generatePalettesREADME() string {
	return `# Color Palettes

This directory contains color palettes that control your site's appearance.
Local palettes override built-in palettes with the same name.

## Palette Structure

Each ` + "`.toml`" + ` file defines a complete color palette:

` + "```toml" + `
[palette]
name = "My Custom Theme"
variant = "dark"  # or "light"
author = "Your Name"
description = "A custom color scheme"

[palette.colors]
# Base colors (raw hex values)
background = "#1e1e2e"
foreground = "#cdd6f4"
primary = "#89b4fa"
# ... more colors

[palette.semantic]
# Semantic mappings (reference colors by name)
text = "foreground"
link = "primary"
# ... more mappings

[palette.components]
# Component-specific colors (reference colors or semantic names)
header_bg = "background"
nav_link = "link"
# ... more components
` + "```" + `

## Color Resolution

Colors are resolved in this order:
1. Direct hex values (e.g., ` + "`#89b4fa`" + `)
2. Raw color names from ` + "`[palette.colors]`" + `
3. Semantic names from ` + "`[palette.semantic]`" + `
4. Component names from ` + "`[palette.components]`" + `

## Creating a Custom Palette

1. Copy an existing palette as a starting point
2. Modify the colors to your liking
3. Update the name, author, and description
4. Save with a new filename

## Using Your Palette

In ` + "`markata-go.toml`" + `:

` + "```toml" + `
[markata-go.theme]
palette = "my-custom-theme"  # Matches your palette's name field
` + "```" + `

## Variants

Each palette should specify a variant:
- ` + "`light`" + ` - For light backgrounds with dark text
- ` + "`dark`" + ` - For dark backgrounds with light text

This helps with automatic contrast and accessibility checks.

## Required Colors

At minimum, palettes should define:
- ` + "`background`" + ` - Main page background
- ` + "`foreground`" + ` - Primary text color
- ` + "`primary`" + ` - Accent color for links and highlights

## Resources

- [Catppuccin](https://catppuccin.com/) - Popular palette collection
- [Coolors](https://coolors.co/) - Color scheme generator
- [Lospec](https://lospec.com/palette-list) - Pixel art palettes
`
}

// generateLayoutsREADME creates documentation for HTML templates.
func generateLayoutsREADME() string {
	return `# HTML Layout Templates

This directory contains Jinja2/Pongo2 templates that control how your content is rendered.
Local templates override built-in templates with the same name.

## Directory Structure

` + "```" + `
templates/
├── base.html           # Base template with HTML structure
├── post.html           # Individual post template
├── feed.html           # Feed/listing template
├── components/         # Reusable components
│   ├── header.html
│   ├── footer.html
│   └── nav.html
└── partials/           # Partial templates
    └── cards/          # Card templates for feeds
        ├── article-card.html
        ├── note-card.html
        └── ...
` + "```" + `

## Template Inheritance

Templates use Jinja2-style inheritance:

` + "```html" + `
{% extends "base.html" %}

{% block content %}
  <article>{{ post.article_html|safe }}</article>
{% endblock %}
` + "```" + `

## Available Variables

### In all templates:
| Variable | Type | Description |
|----------|------|-------------|
| ` + "`config`" + ` | Config | Site configuration |
| ` + "`config.title`" + ` | string | Site title |
| ` + "`config.url`" + ` | string | Site URL |
| ` + "`config.description`" + ` | string | Site description |
| ` + "`config.author`" + ` | string | Site author |

### In post templates (post.html):
| Variable | Type | Description |
|----------|------|-------------|
| ` + "`post`" + ` | Post | Current post object |
| ` + "`post.title`" + ` | string | Post title |
| ` + "`post.slug`" + ` | string | URL slug |
| ` + "`post.href`" + ` | string | Full path (e.g., /my-post/) |
| ` + "`post.date`" + ` | time | Publication date |
| ` + "`post.published`" + ` | bool | Whether post is published |
| ` + "`post.tags`" + ` | []string | Post tags |
| ` + "`post.description`" + ` | string | Post description |
| ` + "`post.content`" + ` | string | Raw markdown content |
| ` + "`post.article_html`" + ` | string | Rendered HTML content |
| ` + "`post.template`" + ` | string | Template name |
| ` + "`post.prev`" + ` | Post | Previous post in sequence |
| ` + "`post.next`" + ` | Post | Next post in sequence |

### In feed templates (feed.html):
| Variable | Type | Description |
|----------|------|-------------|
| ` + "`feed`" + ` | Feed | Current feed object |
| ` + "`feed.title`" + ` | string | Feed title |
| ` + "`feed.slug`" + ` | string | Feed slug |
| ` + "`feed.posts`" + ` | []Post | Posts in this feed |
| ` + "`feed.description`" + ` | string | Feed description |

## Filters

| Filter | Usage | Description |
|--------|-------|-------------|
| ` + "`safe`" + ` | ` + "`{{ html|safe }}`" + ` | Mark string as safe HTML |
| ` + "`date`" + ` | ` + "`{{ post.date|date:\"Jan 2, 2006\" }}`" + ` | Format date |
| ` + "`default`" + ` | ` + "`{{ value|default:\"fallback\" }}`" + ` | Provide default |
| ` + "`length`" + ` | ` + "`{{ items|length }}`" + ` | Get length |
| ` + "`join`" + ` | ` + "`{{ tags|join:\", \" }}`" + ` | Join array |

## Customization Tips

1. **Start small** - Override just the template you need
2. **Use partials** - Break large templates into reusable pieces
3. **Check variables** - Use ` + "`{% if post.description %}...{% endif %}`" + `
4. **Test locally** - Run ` + "`markata-go serve`" + ` to preview changes

## Card Templates

Cards control how posts appear in feeds. The system selects cards based on
the post's ` + "`template`" + ` field:

| Post Template | Card Template |
|--------------|---------------|
| article | partials/cards/article-card.html |
| note | partials/cards/note-card.html |
| photo | partials/cards/photo-card.html |
| video | partials/cards/video-card.html |
| link | partials/cards/link-card.html |
| quote | partials/cards/quote-card.html |
| guide | partials/cards/guide-card.html |
| (default) | partials/cards/default-card.html |

## CSS Variables

Templates use CSS custom properties from your palette:

` + "```css" + `
:root {
  --color-background: {{ palette.background }};
  --color-foreground: {{ palette.foreground }};
  --color-primary: {{ palette.primary }};
  /* ... */
}
` + "```" + `

Access these in your CSS: ` + "`color: var(--color-primary);`" + `
`
}

// promptVendSelection displays vending options and returns selected items.
func promptVendSelection(reader *bufio.Reader) []string {
	vendOptions := []struct {
		name        string
		description string
	}{
		{"templates", "Content templates (for 'markata-go new')"},
		{"palettes", "Color palettes (TOML files)"},
		{"layouts", "HTML layout templates (Jinja2/Pongo2)"},
	}

	outln()
	outln("Select assets to vend (enter numbers separated by spaces):")
	outln()
	for i, opt := range vendOptions {
		outlnf("  %d) [ ] %s", i+1, opt.description)
	}
	outln()
	out("Enter numbers (e.g., 1 2 3 for all): ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return nil
	}

	var selected []string
	parts := strings.Fields(strings.TrimSpace(input))
	for _, p := range parts {
		var idx int
		if _, err := fmt.Sscanf(p, "%d", &idx); err == nil {
			if idx >= 1 && idx <= len(vendOptions) {
				selected = append(selected, vendOptions[idx-1].name)
			}
		}
	}

	return selected
}

// runVendAssets runs the vending process for selected asset types.
func runVendAssets(force bool, assets []string) error {
	for _, asset := range assets {
		outln()
		switch asset {
		case "templates":
			outln("Vending content templates...")
			if err := venderContentTemplates(force); err != nil {
				return fmt.Errorf("failed to vend content templates: %w", err)
			}
		case "palettes":
			outln("Vending palettes...")
			if err := venderPalettes(force); err != nil {
				return fmt.Errorf("failed to vend palettes: %w", err)
			}
		case "layouts":
			outln("Vending layout templates...")
			if err := venderLayouts(force); err != nil {
				return fmt.Errorf("failed to vend layouts: %w", err)
			}
		}
	}
	return nil
}

// updateSiteInfo prompts for updated site information.
func updateSiteInfo(reader *bufio.Reader, cfg *models.Config) error {
	outln()
	outln("Update Site Information")
	outln("-----------------------")

	cfg.Title = prompt(reader, "Site title", cfg.Title)
	cfg.Description = prompt(reader, "Description", cfg.Description)
	cfg.Author = prompt(reader, "Author", cfg.Author)
	cfg.URL = prompt(reader, "URL", cfg.URL)

	return nil
}

// writeConfigTOML writes the configuration to a TOML file.
func writeConfigTOML(path string, cfg *models.Config) error {
	// Build TOML content manually for better formatting
	var sb strings.Builder

	sb.WriteString("# Markata-go configuration file\n\n")

	sb.WriteString("[markata-go]\n")
	sb.WriteString("# Site metadata\n")
	sb.WriteString(fmt.Sprintf("title = %q\n", cfg.Title))
	sb.WriteString(fmt.Sprintf("url = %q\n", cfg.URL))
	sb.WriteString(fmt.Sprintf("description = %q\n", cfg.Description))
	sb.WriteString(fmt.Sprintf("author = %q\n", cfg.Author))
	sb.WriteString("\n")

	sb.WriteString("# Output settings\n")
	sb.WriteString(fmt.Sprintf("output_dir = %q\n", cfg.OutputDir))
	sb.WriteString(fmt.Sprintf("templates_dir = %q\n", cfg.TemplatesDir))
	sb.WriteString(fmt.Sprintf("assets_dir = %q\n", cfg.AssetsDir))
	sb.WriteString("\n")

	// Glob config
	sb.WriteString("# File discovery\n")
	sb.WriteString("[markata-go.glob]\n")
	if len(cfg.GlobConfig.Patterns) > 0 {
		sb.WriteString("patterns = [")
		for i, p := range cfg.GlobConfig.Patterns {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q", p))
		}
		sb.WriteString("]\n")
	} else {
		sb.WriteString("patterns = [\"**/*.md\"]\n")
	}
	sb.WriteString(fmt.Sprintf("use_gitignore = %v\n", cfg.GlobConfig.UseGitignore))
	sb.WriteString("\n")

	// Theme config
	if cfg.Theme.Palette != "" {
		sb.WriteString("# Theme configuration\n")
		sb.WriteString("[markata-go.theme]\n")
		if cfg.Theme.Name != "" {
			sb.WriteString(fmt.Sprintf("name = %q\n", cfg.Theme.Name))
		}
		sb.WriteString(fmt.Sprintf("palette = %q\n", cfg.Theme.Palette))
		sb.WriteString("\n")
	}

	// SEO config
	if cfg.SEO.TwitterHandle != "" || cfg.SEO.DefaultImage != "" || cfg.SEO.LogoURL != "" {
		sb.WriteString("# SEO configuration\n")
		sb.WriteString("[markata-go.seo]\n")
		if cfg.SEO.TwitterHandle != "" {
			sb.WriteString(fmt.Sprintf("twitter_handle = %q\n", cfg.SEO.TwitterHandle))
		}
		if cfg.SEO.DefaultImage != "" {
			sb.WriteString(fmt.Sprintf("default_image = %q\n", cfg.SEO.DefaultImage))
		}
		if cfg.SEO.LogoURL != "" {
			sb.WriteString(fmt.Sprintf("logo_url = %q\n", cfg.SEO.LogoURL))
		}
		sb.WriteString("\n")
	}

	// Post formats config
	sb.WriteString("# Post output formats\n")
	sb.WriteString("[markata-go.post_formats]\n")
	sb.WriteString(fmt.Sprintf("html = %v\n", cfg.PostFormats.IsHTMLEnabled()))
	sb.WriteString(fmt.Sprintf("markdown = %v\n", cfg.PostFormats.Markdown))
	sb.WriteString(fmt.Sprintf("text = %v\n", cfg.PostFormats.Text))
	sb.WriteString(fmt.Sprintf("og = %v\n", cfg.PostFormats.OG))
	sb.WriteString("\n")

	// Feed defaults
	sb.WriteString("# Feed defaults\n")
	sb.WriteString("[markata-go.feed_defaults]\n")
	sb.WriteString(fmt.Sprintf("items_per_page = %d\n", cfg.FeedDefaults.ItemsPerPage))
	sb.WriteString(fmt.Sprintf("orphan_threshold = %d\n", cfg.FeedDefaults.OrphanThreshold))
	sb.WriteString("\n")

	sb.WriteString("[markata-go.feed_defaults.formats]\n")
	sb.WriteString(fmt.Sprintf("html = %v\n", cfg.FeedDefaults.Formats.HTML))
	sb.WriteString(fmt.Sprintf("rss = %v\n", cfg.FeedDefaults.Formats.RSS))
	sb.WriteString(fmt.Sprintf("atom = %v\n", cfg.FeedDefaults.Formats.Atom))
	sb.WriteString(fmt.Sprintf("json = %v\n", cfg.FeedDefaults.Formats.JSON))
	sb.WriteString(fmt.Sprintf("sitemap = %v\n", cfg.FeedDefaults.Formats.Sitemap))
	sb.WriteString("\n")

	// Custom feeds (commented example)
	sb.WriteString("# Define custom feeds\n")
	sb.WriteString("# [[markata-go.feeds]]\n")
	sb.WriteString("# slug = \"blog\"\n")
	sb.WriteString("# title = \"Blog Posts\"\n")
	sb.WriteString("# filter = \"published == true\"\n")
	sb.WriteString("# sort = \"date\"\n")
	sb.WriteString("# reverse = true\n")

	return os.WriteFile(path, []byte(sb.String()), 0o644) //nolint:gosec // config files should be readable
}

// runExistingProjectWizard runs the wizard for an existing project.
func runExistingProjectWizard(reader *bufio.Reader, configPath string) error {
	outln()
	outln("Found existing markata-go.toml")

	// Load existing config
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load existing config: %w", err)
	}

	menuOptions := []string{
		"Add new features",
		"Update site information",
		"Vend built-in assets (templates, palettes, layouts)",
		"View current configuration",
		"Exit",
	}

	for {
		choice := promptMenuChoice(reader, "What would you like to do?", menuOptions)

		switch choice {
		case 0: // Add new features
			configured := detectConfiguredFeatures(cfg)
			features := getAvailableFeatures(configured)
			selected := promptFeatureSelection(reader, features)

			if len(selected) == 0 {
				outln("\nNo features selected.")
				continue
			}

			// Add each selected feature
			for _, feature := range selected {
				if err := addFeature(reader, feature, cfg); err != nil {
					return fmt.Errorf("failed to configure %s: %w", feature, err)
				}
			}

			// Backup and write config
			if err := backupConfig(configPath); err != nil {
				warnf("%v", err)
			}

			if err := writeConfigTOML(configPath, cfg); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			outln()
			outln("  Updated markata-go.toml with new features")
			outln()
			outln("Run 'markata-go build' to apply changes!")
			return nil

		case 1: // Update site information
			if err := updateSiteInfo(reader, cfg); err != nil {
				return err
			}

			// Backup and write config
			if err := backupConfig(configPath); err != nil {
				warnf("%v", err)
			}

			if err := writeConfigTOML(configPath, cfg); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			outln()
			outln("  Updated markata-go.toml")
			outln()
			outln("Run 'markata-go build' to apply changes!")
			return nil

		case 2: // Vend built-in assets
			selected := promptVendSelection(reader)
			if len(selected) == 0 {
				outln("\nNo assets selected.")
				continue
			}

			if err := runVendAssets(initForce, selected); err != nil {
				return err
			}

			outln()
			outln("Assets vendored successfully!")
			outln("Local files will now override built-in defaults.")
			return nil

		case 3: // View current configuration
			displayCurrentConfig(cfg)
			// Continue the loop to show menu again

		case 4: // Exit
			outln("\nExiting without changes.")
			return nil
		}
	}
}

func runInitCommand(cmd *cobra.Command, _ []string) error {
	currentCmd = cmd
	if noInput {
		return fmt.Errorf("interactive setup is disabled by --no-input; rerun without it or pass a config file path to 'markata-go config init'")
	}
	// Determine if we should use the huh TUI wizard
	// Use plain mode if: --plain flag is set, or stdin is not a TTY
	usePlain := initPlain || !inputIsTerminal() || !outputIsTerminal()

	configPath := defaultConfigFilename

	// Check for existing config file
	if _, err := os.Stat(configPath); err == nil {
		if initForce {
			// Force mode - proceed with new project setup
			infof("--force specified, overwriting existing configuration")
		} else {
			// Run the wizard for existing projects
			if usePlain {
				reader := bufio.NewReader(inReader())
				return runExistingProjectWizard(reader, configPath)
			}
			// Try to load existing config to get palette for theming
			theme := createHuhThemeForPalette(getPaletteFromConfig(configPath))
			return runHuhExistingProjectWizard(configPath, theme, initForce)
		}
	}

	// New project wizard
	if usePlain {
		return runPlainNewProjectWizard()
	}

	// Use huh TUI wizard
	theme := createHuhTheme("")
	state, err := runHuhNewProjectWizard(theme)
	if err != nil {
		return err
	}

	return applyWizardState(state, initForce)
}

// getPaletteFromConfig attempts to extract the palette name from an existing config file.
func getPaletteFromConfig(configPath string) *palettes.Palette {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil
	}
	palette, _ := loadLoggerPalette(cfg.Theme, map[string]any{"config_path": configPath})
	return palette
}

// runPlainNewProjectWizard runs the plain text wizard for new projects.
func runPlainNewProjectWizard() error {
	reader := bufio.NewReader(inReader())

	outln()
	outln("Welcome to markata-go!")
	outln()

	// Gather site information
	title := prompt(reader, "Site title", "My Site")
	description := prompt(reader, "Description", "A site built with markata-go")
	author := prompt(reader, "Author", "")
	url := prompt(reader, "URL", "https://example.com")

	outln()
	outln("Creating project structure...")

	// Create directories
	dirs := []string{"posts", "static"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		outlnf("  Created %s/", dir)
	}

	// Create config with defaults and user input
	cfg := models.NewConfig()
	cfg.Title = title
	cfg.Description = description
	cfg.Author = author
	cfg.URL = url
	cfg.GlobConfig.UseGitignore = true

	// Ask about optional features for new projects
	outln()
	if promptRadioBool(reader, "Configure additional features", false) {
		configured := detectConfiguredFeatures(cfg)
		features := getAvailableFeatures(configured)
		selected := promptFeatureSelection(reader, features)

		for _, feature := range selected {
			if err := addFeature(reader, feature, cfg); err != nil {
				return fmt.Errorf("failed to configure %s: %w", feature, err)
			}
		}
	}

	// Write config file
	configPath := defaultConfigFilename
	if err := writeConfigTOML(configPath, cfg); err != nil {
		return fmt.Errorf("failed to write %s: %w", defaultConfigFilename, err)
	}
	outlnf("  Created %s", defaultConfigFilename)

	outln()

	// Offer to vend built-in assets
	if promptRadioBool(reader, "Vend built-in assets for customization? (Obsidian-compatible)", false) {
		selected := promptVendSelection(reader)
		if len(selected) > 0 {
			if err := runVendAssets(initForce, selected); err != nil {
				return err
			}
			outln()
			outln("Assets vendored! Local files will override built-in defaults.")
		}
	}

	outln()

	// Offer to create first post
	if promptRadioBool(reader, "Create your first post", true) {
		postTitle := prompt(reader, "Post title", "Hello World")

		slug := generateSlug(postTitle)
		filename := slug + ".md"
		fullPath := filepath.Join("posts", filename)

		// Check if file already exists
		if _, err := os.Stat(fullPath); err == nil && !initForce {
			outlnf("  ! Post already exists: %s (skipped)", fullPath)
		} else {
			now := time.Now()
			content := generatePostContent(title, slug, now, false)
			if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil { //nolint:gosec // content files should be readable
				return fmt.Errorf("failed to write post: %w", err)
			}
			outlnf("  Created %s", fullPath)
		}
	}

	outln()
	outln("Done! Run 'markata-go serve' to start.")
	outln()

	return nil
}

// Ensure toml package is used (for potential future direct encoding)
var _ = toml.Unmarshal
