package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var contentTemplateAliases = map[string][]string{
	"article": {"blog-post", "essay", "tutorial"},
	"note":    {"ping", "thought", "status", "tweet"},
	"photo":   {"shot", "shots", "image", "gallery"},
	"video":   {"clip", "cast", "stream"},
	"link":    {"bookmark", "til", "stars"},
	"quote":   {"quotation"},
	"guide":   {"series", "step", "chapter"},
	"inline":  {"gratitude", "micro"},
	"contact": {"character", "person"},
}

var (
	// newDir is the directory for new posts (overrides template placement).
	newDir string

	// newDraft creates the post as a draft.
	newDraft bool

	// newTags is a comma-separated list of tags.
	newTags string

	// newTemplate specifies the content template to use.
	newTemplate string

	// newList lists available templates.
	newList bool

	// newPlain uses plain text prompts instead of the huh TUI wizard.
	newPlain bool
)

// ContentTemplate represents a content template with its configuration.
type ContentTemplate struct {
	Name            string
	Description     string
	Aliases         []string
	ValidationError string
	Directory       string
	Frontmatter     map[string]interface{}
	Body            string
	Source          string // "builtin", "config", or "file"
}

// ToMarkdown converts a ContentTemplate to a markdown file format with YAML frontmatter.
// The _directory field is included to preserve the placement information.
func (ct *ContentTemplate) ToMarkdown() string {
	// Build frontmatter with _directory
	fm := make(map[string]interface{})
	fm["_directory"] = ct.Directory
	for k, v := range ct.Frontmatter {
		fm[k] = v
	}

	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		// Fallback to minimal frontmatter
		return fmt.Sprintf("---\n_directory: %s\n---\n\n%s\n", ct.Directory, ct.Body)
	}

	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(fmBytes)
	sb.WriteString("---\n\n")
	sb.WriteString(ct.Body)
	sb.WriteString("\n")

	return sb.String()
}

// BuiltinTemplates returns the default built-in content templates.
// This is exported for use by the init command for vending templates.
func BuiltinTemplates() map[string]ContentTemplate {
	return builtinTemplates()
}

// builtinTemplates returns the default built-in content templates.
func builtinTemplates() map[string]ContentTemplate {
	return map[string]ContentTemplate{
		"post": {
			Name:      "post",
			Directory: "pages/post",
			Frontmatter: map[string]interface{}{
				"template": "post",
			},
			Body:   "Write your content here...",
			Source: "builtin",
		},
		"page": {
			Name:      "page",
			Directory: "pages",
			Frontmatter: map[string]interface{}{
				"template": "page",
			},
			Body:   "Write your page content here...",
			Source: "builtin",
		},
		"docs": {
			Name:      "docs",
			Directory: "docs",
			Frontmatter: map[string]interface{}{
				"template": "docs",
			},
			Body:   "Write your documentation here...",
			Source: "builtin",
		},
		"article": {
			Name:      "article",
			Directory: "pages/article",
			Frontmatter: map[string]interface{}{
				"template": "article",
			},
			Body:   "Write your article here...",
			Source: "builtin",
		},
		"note": {
			Name:      "note",
			Directory: "pages/note",
			Frontmatter: map[string]interface{}{
				"template": "note",
			},
			Body:   "A quick note...",
			Source: "builtin",
		},
		"photo": {
			Name:      "photo",
			Directory: "pages/photo",
			Frontmatter: map[string]interface{}{
				"template": "photo",
				"image":    "",
			},
			Body:   "Photo caption...",
			Source: "builtin",
		},
		"video": {
			Name:      "video",
			Directory: "pages/video",
			Frontmatter: map[string]interface{}{
				"template": "video",
				"video":    "",
				"image":    "",
				"duration": "",
			},
			Body:   "Video description...",
			Source: "builtin",
		},
		"link": {
			Name:      "link",
			Directory: "pages/link",
			Frontmatter: map[string]interface{}{
				"template": "link",
				"url":      "",
				"image":    "",
			},
			Body:   "Why I'm sharing this link...",
			Source: "builtin",
		},
		"quote": {
			Name:      "quote",
			Directory: "pages/quote",
			Frontmatter: map[string]interface{}{
				"template": "quote",
				"quote":    "",
				"source":   "",
			},
			Body:   "Additional commentary on this quote...",
			Source: "builtin",
		},
		"guide": {
			Name:      "guide",
			Directory: "pages/guide",
			Frontmatter: map[string]interface{}{
				"template": "guide",
			},
			Body:   "## Introduction\n\nWrite your guide here...",
			Source: "builtin",
		},
		"inline": {
			Name:      "inline",
			Directory: "pages/inline",
			Frontmatter: map[string]interface{}{
				"template": "inline",
			},
			Body:   "Inline content...",
			Source: "builtin",
		},
		"contact": {
			Name:      "contact",
			Directory: "pages/contact",
			Frontmatter: map[string]interface{}{
				"template": "contact",
				"handle":   "",
				"aliases":  []string{},
				"avatar":   "",
				"url":      "",
			},
			Body:   "Bio or contact details...",
			Source: "builtin",
		},
		"author": {
			Name:      "author",
			Directory: "pages/author",
			Frontmatter: map[string]interface{}{
				"template": "author",
				"name":     "",
				"bio":      "",
				"role":     "",
				"avatar":   "",
				"url":      "",
				"email":    "",
				"social":   map[string]string{},
				"active":   true,
				"guest":    false,
			},
			Body:   "Extended author bio...",
			Source: "builtin",
		},
	}
}

// newCmd represents the new command.
var newCmd = &cobra.Command{
	Use:   "new [title]",
	Short: "Create a new content file",
	Long: `Create a new markdown content file with frontmatter template.

The command generates a new markdown file with:
  - Title set from the argument (or prompted if not provided)
  - Slug generated from the title
  - Current date
  - Draft status (configurable)
  - Tags (optional)
  - Template-specific frontmatter and placement

Template System:
  Templates control the default frontmatter and output directory for new content.
  Built-in templates: post, page, docs, article, note, photo, video, link, quote, guide, inline, contact, author

  Custom templates can be defined:
  1. In markata-go.toml under [content_templates]
  2. As markdown files in the content-templates/ directory

Example usage:
  markata-go new "My First Post"                    # Create posts/my-first-post.md (default: post)
  markata-go new "About" --template page            # Create pages/about.md
  markata-go new "Getting Started" --template docs  # Create docs/getting-started.md
  markata-go new "Hello World" --dir blog           # Override directory: blog/hello-world.md
  markata-go new --list                             # List available templates
  markata-go new                                    # Interactive mode`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNewCommand,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringVar(&newDir, "dir", "", "directory for new content (overrides template placement)")
	newCmd.Flags().BoolVar(&newDraft, "draft", false, "create as draft (default: false, posts are published by default)")
	newCmd.Flags().StringVar(&newTags, "tags", "", "comma-separated list of tags")
	newCmd.Flags().StringVarP(&newTemplate, "template", "t", "post", "content template to use")
	newCmd.Flags().BoolVarP(&newList, "list", "l", false, "list available templates")
	newCmd.Flags().BoolVar(&newPlain, "plain", false, "use plain text prompts (for non-TTY environments)")
}

// loadTemplates discovers and loads all available content templates.
func loadTemplates() map[string]ContentTemplate {
	templates := builtinTemplates()

	// Load templates from config if available
	cfg := loadConfigSafe()
	if cfg != nil {
		// Apply placement overrides from config
		for name, dir := range cfg.ContentTemplates.Placement {
			if t, exists := templates[name]; exists {
				t.Directory = dir
				templates[name] = t
			}
		}

		// Add/override templates from config
		for _, ct := range cfg.ContentTemplates.Templates {
			templates[ct.Name] = ContentTemplate{
				Name:        ct.Name,
				Directory:   ct.Directory,
				Frontmatter: ct.Frontmatter,
				Body:        ct.Body,
				Source:      "config",
			}
		}

		// Load templates from content-templates directory
		templatesDir := cfg.ContentTemplates.Directory
		if templatesDir == "" {
			templatesDir = "content-templates"
		}
		loadTemplatesFromDir(templatesDir, templates)
	} else {
		// No config, try default directory
		loadTemplatesFromDir("content-templates", templates)
	}

	return templates
}

func templateAliases(name string, template ContentTemplate) []string {
	aliases := append([]string(nil), contentTemplateAliases[name]...)
	aliases = append(aliases, template.Aliases...)
	return aliases
}

func formatTemplateEntry(name string, template ContentTemplate) string {
	aliases := templateAliases(name, template)
	label := name
	if len(aliases) > 0 {
		label = fmt.Sprintf("%s (aka: %s)", name, strings.Join(aliases, ", "))
	}
	entry := fmt.Sprintf("%s -> %s/ (%s)", label, template.Directory, template.Source)
	if template.Description != "" {
		entry += " — " + template.Description
	}
	if template.ValidationError != "" {
		entry += " [invalid: " + template.ValidationError + "]"
	}
	return entry
}

func resolveTemplateSelection(name string, templates map[string]ContentTemplate) (ContentTemplate, string, bool) {
	// Canonical names always win over aliases, even when an alias is ambiguous.
	if template, exists := templates[name]; exists {
		return template, name, true
	}
	canonical := ""
	for candidate, item := range templates {
		for _, alias := range templateAliases(candidate, item) {
			if alias != name {
				continue
			}
			if canonical != "" && canonical != candidate {
				return ContentTemplate{}, "", false
			}
			canonical = candidate
		}
	}
	if canonical == "" {
		return ContentTemplate{}, "", false
	}
	template := templates[canonical]
	frontmatter := make(map[string]interface{}, len(template.Frontmatter))
	for key, value := range template.Frontmatter {
		frontmatter[key] = value
	}
	// Retain the existing alias-specific page template behavior.
	frontmatter["template"] = name
	frontmatter["templateKey"] = name
	template.Frontmatter = frontmatter
	return template, canonical, true
}

// loadTemplatesFromDir loads templates from markdown files in a directory.
func loadTemplatesFromDir(dir string, templates map[string]ContentTemplate) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Directory doesn't exist or can't be read - that's fine
		return
	}

	sources := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		path := filepath.Join(dir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		template := parseTemplateFile(name, string(content))
		template.Source = "file"
		if previous, duplicate := sources[template.Name]; duplicate {
			first := templates[template.Name]
			first.ValidationError = fmt.Sprintf("duplicate canonical starter %q in %s and %s", template.Name, previous, path)
			templates[template.Name] = first
			continue
		}
		sources[template.Name] = path
		templates[template.Name] = template
	}
}

// parseTemplateFile parses a markdown template file with frontmatter.
func parseTemplateFile(name, content string) ContentTemplate {
	template := ContentTemplate{
		Name:        name,
		Directory:   name, // Default directory is same as template name
		Frontmatter: make(map[string]interface{}),
		Body:        "",
	}

	// Check for frontmatter
	if !strings.HasPrefix(content, "---") {
		template.Body = strings.TrimSpace(content)
		return template
	}

	// Find end of frontmatter
	parts := strings.SplitN(content[3:], "---", 2)
	if len(parts) < 2 {
		template.ValidationError = "frontmatter is missing its closing delimiter"
		template.Body = strings.TrimSpace(content)
		return template
	}

	// Parse frontmatter YAML
	frontmatterYAML := strings.TrimSpace(parts[0])
	if err := yaml.Unmarshal([]byte(frontmatterYAML), &template.Frontmatter); err != nil {
		template.ValidationError = fmt.Sprintf("invalid starter frontmatter: %v", err)
	} else {
		// Extract directory from frontmatter if present
		if dir, ok := template.Frontmatter["_directory"].(string); ok {
			template.Directory = dir
			delete(template.Frontmatter, "_directory")
		}
		if rawStarter, exists := template.Frontmatter["starter"]; exists {
			starter, ok := rawStarter.(map[string]interface{})
			if !ok {
				template.ValidationError = "starter metadata must be a mapping"
				delete(template.Frontmatter, "starter")
			} else {
				if description, ok := starter["description"].(string); ok {
					template.Description = description
				} else if _, exists := starter["description"]; exists {
					template.ValidationError = "starter.description must be a string"
				}
				if canonicalName, ok := starter["name"].(string); ok && strings.TrimSpace(canonicalName) != "" {
					defaultDirectory := template.Name
					template.Name = strings.TrimSpace(canonicalName)
					if template.Directory == defaultDirectory {
						template.Directory = template.Name
					}
				} else if _, exists := starter["name"]; exists {
					template.ValidationError = "starter.name must be a non-empty string"
				}
				if rawAliases, exists := starter["aliases"]; exists {
					aliases, ok := rawAliases.([]interface{})
					if !ok {
						template.ValidationError = "starter.aliases must be a list of strings"
					}
					for _, alias := range aliases {
						value, ok := alias.(string)
						if !ok || strings.TrimSpace(value) == "" {
							template.ValidationError = "starter.aliases must contain non-empty strings"
							continue
						}
						template.Aliases = append(template.Aliases, strings.TrimSpace(value))
					}
				}
				delete(template.Frontmatter, "starter")
			}
		}
	}

	// Body is everything after frontmatter, but preserve the template markers
	template.Body = strings.TrimSpace(parts[1])

	return template
}

// loadConfigSafe attempts to load the config without errors.
func loadConfigSafe() *configWrapper {
	var cfg *configWrapper
	if cfgFile != "" {
		parsed, err := parseConfigFile(cfgFile)
		if err == nil {
			cfg = parsed
		}
	}

	if cfg == nil && cfgFile == "" {
		configPaths := []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"}
		for _, path := range configPaths {
			parsed, err := parseConfigFile(path)
			if err == nil {
				cfg = parsed
				break
			}
		}
	}

	for _, path := range mergeConfigFiles {
		parsed, err := parseConfigFile(path)
		if err != nil {
			continue
		}
		if cfg == nil {
			cfg = &configWrapper{}
		}
		mergeTemplateConfig(cfg, parsed)
	}
	return cfg
}

func mergeTemplateConfig(base, override *configWrapper) {
	if override.ContentTemplates.Directory != "" {
		base.ContentTemplates.Directory = override.ContentTemplates.Directory
	}
	if len(override.ContentTemplates.Placement) > 0 {
		if base.ContentTemplates.Placement == nil {
			base.ContentTemplates.Placement = make(map[string]string)
		}
		for name, directory := range override.ContentTemplates.Placement {
			base.ContentTemplates.Placement[name] = directory
		}
	}
	if override.ContentTemplates.Templates != nil {
		base.ContentTemplates.Templates = override.ContentTemplates.Templates
	}
}

// configWrapper wraps the content templates config for safe loading.
type configWrapper struct {
	ContentTemplates struct {
		Directory string            `yaml:"directory" toml:"directory" json:"directory"`
		Placement map[string]string `yaml:"placement" toml:"placement" json:"placement"`
		Templates []struct {
			Name        string                 `yaml:"name" toml:"name" json:"name"`
			Directory   string                 `yaml:"directory" toml:"directory" json:"directory"`
			Frontmatter map[string]interface{} `yaml:"frontmatter" toml:"frontmatter" json:"frontmatter"`
			Body        string                 `yaml:"body" toml:"body" json:"body"`
		} `yaml:"templates" toml:"templates" json:"templates"`
	} `yaml:"content_templates" toml:"content_templates" json:"content_templates"`
}

// parseConfigFile parses a config file to extract content templates config.
func parseConfigFile(path string) (*configWrapper, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg configWrapper
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(content, &cfg)
	case ".toml":
		err = toml.Unmarshal(content, &cfg)
	case ".json":
		err = json.Unmarshal(content, &cfg)
	default:
		return nil, fmt.Errorf("unsupported config format: %s", ext)
	}

	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// interactiveInput holds the collected input from interactive mode.
type interactiveInput struct {
	title    string
	template ContentTemplate
	tags     []string
}

// runInteractiveMode prompts the user for input when no title is provided.
func runInteractiveMode(cmd *cobra.Command, templates map[string]ContentTemplate) (*interactiveInput, error) {
	reader := bufio.NewReader(inReader())

	errln()

	// Get title
	title := promptNew(reader, "Title", "")
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	template, _, ok := resolveTemplateSelection(newTemplate, templates)
	if !ok {
		return nil, fmt.Errorf("unknown template %q", newTemplate)
	}

	// Get template (only prompt if not explicitly set via flag)
	if !cmd.Flags().Changed("template") {
		templateNames := make([]string, 0, len(templates))
		for name := range templates {
			templateNames = append(templateNames, name)
		}
		sort.Strings(templateNames)
		errln("Available templates:")
		for _, name := range templateNames {
			errlnf("  %s", formatTemplateEntry(name, templates[name]))
		}
		templateInput := promptNew(reader, "Template", "post")
		if t, _, ok := resolveTemplateSelection(templateInput, templates); ok {
			newTemplate = templateInput
			template = t
		}
	}

	// Get directory (only prompt if not explicitly set via flag)
	if !cmd.Flags().Changed("dir") {
		dirDefault := template.Directory
		newDir = promptNew(reader, "Directory", dirDefault)
	}

	// Get tags (only prompt if not explicitly set via flag)
	var tags []string
	if !cmd.Flags().Changed("tags") {
		tagsInput := promptNew(reader, "Tags (comma-separated)", "")
		if tagsInput != "" {
			tags = parseTags(tagsInput)
		}
	} else if newTags != "" {
		tags = parseTags(newTags)
	}

	// Get draft status (only prompt if not explicitly set via flag)
	if !cmd.Flags().Changed("draft") {
		newDraft = promptYesNoNew(reader, "Create as draft?", false)
	}

	errln()

	return &interactiveInput{
		title:    title,
		template: template,
		tags:     tags,
	}, nil
}

// writeContentFile creates the content file with the given parameters.
func writeContentFile(title, slug, outputDir string, draft, draftSet bool, tags []string, template ContentTemplate) error {
	if template.ValidationError != "" {
		return fmt.Errorf("invalid starter %q: %s", template.Name, template.ValidationError)
	}
	filename := slug + ".md"
	fullPath := filepath.Join(outputDir, filename)

	// Check if file already exists
	if _, err := os.Stat(fullPath); err == nil {
		return fmt.Errorf("file already exists: %s", fullPath)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate content
	now := time.Now()
	content := generateTemplatedContent(title, slug, now, draft, draftSet, tags, template)

	// Write file (0o644 is appropriate for content files that should be world-readable)
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil { //nolint:gosec // content files should be readable
		return fmt.Errorf("failed to write file: %w", err)
	}

	outlnf("Created: %s", fullPath)
	if verbose {
		errlnf("  Template: %s", newTemplate)
		errlnf("  Title: %s", title)
		errlnf("  Slug: %s", slug)
		errlnf("  Date: %s", now.Format("2006-01-02"))
		errlnf("  Draft: %t", draft)
		if len(tags) > 0 {
			errlnf("  Tags: %s", strings.Join(tags, ", "))
		}
	}

	return nil
}

func runNewCommand(cmd *cobra.Command, args []string) error {
	currentCmd = cmd
	// Handle --list flag
	if newList {
		return listTemplates()
	}

	// Load templates
	templates := loadTemplates()

	// Validate template exists
	template, _, exists := resolveTemplateSelection(newTemplate, templates)
	if !exists {
		availableNames := make([]string, 0, len(templates))
		for name := range templates {
			availableNames = append(availableNames, name)
		}
		sort.Strings(availableNames)
		return fmt.Errorf("unknown template %q; available templates: %s", newTemplate, strings.Join(availableNames, ", "))
	}
	var title string
	var tags []string

	// If no title provided, run interactive mode
	if len(args) == 0 {
		if noInput {
			return fmt.Errorf("title is required when --no-input is set; run 'markata-go new \"My Post\"' or pass a title argument")
		}
		// Determine if we should use the huh TUI wizard
		usePlain := newPlain || !inputIsTerminal() || !outputIsTerminal()

		if usePlain {
			// Use legacy plain text interactive mode
			input, err := runInteractiveMode(cmd, templates)
			if err != nil {
				return err
			}
			title = input.title
			template = input.template
			tags = input.tags
		} else {
			// Use huh TUI wizard
			state, err := runHuhNewWizard(templates)
			if err != nil {
				return err
			}
			return applyNewWizardState(state, templates)
		}
	} else {
		title = args[0]
		if newTags != "" {
			tags = parseTags(newTags)
		} else if cmd.Flags().Changed("tags") {
			tags = []string{}
		}
	}

	// Determine output directory
	outputDir := template.Directory
	if cmd.Flags().Changed("dir") || newDir != "" {
		outputDir = newDir
	}
	if outputDir == "" {
		outputDir = template.Directory
	}
	if !cmd.Flags().Changed("dir") && !safeStarterDirectory(outputDir) {
		return fmt.Errorf("starter %q has an unsafe output directory %q", template.Name, outputDir)
	}
	if siteDirSelected && !filepath.IsAbs(outputDir) {
		outputDir = filepath.Join(activeSiteDir, outputDir)
	}

	// Generate slug from title
	slug := generateSlug(title)

	draftSet := cmd.Flags().Changed("draft")
	if len(args) == 0 && !draftSet {
		// An interactive answer is an override only if it changes the starter default.
		draftSet = newDraft != starterDraftDefault(template)
	}
	return writeContentFile(title, slug, outputDir, newDraft, draftSet, tags, template)
}

func starterDraftDefault(template ContentTemplate) bool {
	if draft, ok := template.Frontmatter["draft"].(bool); ok {
		return draft
	}
	if published, ok := template.Frontmatter["published"].(bool); ok {
		return !published
	}
	return false
}

func safeStarterDirectory(directory string) bool {
	if directory == "" || filepath.IsAbs(directory) {
		return false
	}
	clean := filepath.Clean(directory)
	return clean != ".." && !strings.HasPrefix(clean, ".."+string(filepath.Separator))
}

// listTemplates prints available templates.
func listTemplates() error {
	templates := loadTemplates()

	if len(templates) == 0 {
		outln("No templates available.")
		return nil
	}

	// Sort template names
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	outln("Available content templates:")
	outln()
	for _, name := range names {
		t := templates[name]
		outlnf("  %s", formatTemplateEntry(name, t))
	}
	outln()
	outln("Use --template <name> or -t <name> to select a template.")
	return nil
}

// generateTemplatedContent creates the markdown content with template-specific frontmatter.
func generateTemplatedContent(title, slug string, date time.Time, draft, draftSet bool, tags []string, template ContentTemplate) string {
	// Only the selected starter's frontmatter and explicit user overrides are
	// serialized. A value inferred from site/built-in defaults is not a field.
	fm := make(map[string]interface{}, len(template.Frontmatter))
	for key, value := range template.Frontmatter {
		fm[key] = expandStarterValue(value, title, slug, date)
	}
	if _, exists := fm["title"]; exists {
		fm["title"] = title
	}
	if draftSet {
		fm["published"] = !draft
		fm["draft"] = draft
	}
	if tags != nil {
		fm["tags"] = tags
	}

	body := expandStarterText(template.Body, title, slug, date)
	if body == "" {
		body = "Write your content here..."
	}
	if len(fm) == 0 {
		return body + "\n"
	}
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		// Preserve the existing fallback on invalid frontmatter values.
		return generatePostContentWithTags(title, slug, date, draft, tags)
	}
	var sb strings.Builder
	sb.WriteString("---\n")
	sb.Write(fmBytes)
	sb.WriteString("---\n\n")
	sb.WriteString(body)
	sb.WriteString("\n")
	return sb.String()
}

func expandStarterValue(value interface{}, title, slug string, date time.Time) interface{} {
	switch typed := value.(type) {
	case string:
		return expandStarterText(typed, title, slug, date)
	case []interface{}:
		out := make([]interface{}, len(typed))
		for i, item := range typed {
			out[i] = expandStarterValue(item, title, slug, date)
		}
		return out
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			out[key] = expandStarterValue(item, title, slug, date)
		}
		return out
	default:
		return value
	}
}

func expandStarterText(value, title, slug string, date time.Time) string {
	replacer := strings.NewReplacer("{{title}}", title, "{{slug}}", slug, "{{date}}", date.Format("2006-01-02"), "{{datetime}}", date.UTC().Format(time.RFC3339))
	return replacer.Replace(value)
}

// generateSlug creates a URL-safe slug from a title.
func generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove non-alphanumeric characters (except hyphens)
	reg := regexp.MustCompile(`[^a-z0-9\-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Collapse multiple hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// generatePostContent creates the markdown content with frontmatter.
func generatePostContent(title, slug string, date time.Time, draft bool) string {
	return generatePostContentWithTags(title, slug, date, draft, nil)
}

// generatePostContentWithTags creates the markdown content with frontmatter and optional tags.
func generatePostContentWithTags(title, slug string, date time.Time, draft bool, tags []string) string {
	published := !draft

	// Format tags as YAML array
	tagsYAML := "[]"
	if len(tags) > 0 {
		var quotedTags []string
		for _, tag := range tags {
			quotedTags = append(quotedTags, fmt.Sprintf("%q", tag))
		}
		tagsYAML = "[" + strings.Join(quotedTags, ", ") + "]"
	}

	return fmt.Sprintf(`---
title: "%s"
slug: "%s"
date: %s
published: %t
draft: %t
tags: %s
description: ""
---

Write your content here...
`, title, slug, date.Format("2006-01-02"), published, draft, tagsYAML)
}

// promptNew displays a question and returns the user's response or a default value.
func promptNew(reader *bufio.Reader, question, defaultVal string) string {
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

// promptYesNoNew displays a yes/no question and returns the boolean result.
func promptYesNoNew(reader *bufio.Reader, question string, defaultYes bool) bool {
	defaultStr := "y/N"
	if defaultYes {
		defaultStr = "Y/n"
	}
	errf("%s (%s): ", question, defaultStr)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultYes
	}
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return defaultYes
	}
	return input == "y" || input == "yes"
}

// parseTags splits a comma-separated tag string into a slice of trimmed tags.
func parseTags(tagsStr string) []string {
	var tags []string
	for _, tag := range strings.Split(tagsStr, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}
