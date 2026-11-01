package workbench

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type contentTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Directory   string   `json:"directory"`
	Aliases     []string `json:"aliases,omitempty"`
	Source      string   `json:"source"`
}

type siteStarterMetadata struct {
	name             string
	description      string
	directory        string
	aliases          []string
	cleanFrontmatter []byte
}

func (h *Handler) contentTemplates() ([]contentTemplate, error) {
	today := time.Now().Format("2006-01-02")
	builtins := builtinContentTemplates(today)
	templates := append([]contentTemplate(nil), builtins...)
	directory := "content-templates"
	if h.siteConfig != nil {
		configured := h.siteConfig.ContentTemplates
		if configured.Directory != "" {
			directory = configured.Directory
		}
		for i := range templates {
			if placement, ok := configured.Placement[contentTemplateKey(templates[i])]; ok {
				if !safeStarterOutputDirectory(placement) {
					return nil, fmt.Errorf("unsafe content starter placement %q", placement)
				}
				templates[i].Directory = placement
			}
		}
		for _, configTemplate := range configured.Templates {
			if strings.TrimSpace(configTemplate.Name) == "" {
				return nil, fmt.Errorf("configured content starter needs a name")
			}
			outputDir := configTemplate.Directory
			if outputDir == "" {
				outputDir = configTemplate.Name
			}
			if !safeStarterOutputDirectory(outputDir) {
				return nil, fmt.Errorf("configured starter %q has unsafe directory %q", configTemplate.Name, outputDir)
			}
			var source strings.Builder
			if len(configTemplate.Frontmatter) > 0 {
				header, err := yaml.Marshal(configTemplate.Frontmatter)
				if err != nil {
					return nil, fmt.Errorf("configured starter %q: %w", configTemplate.Name, err)
				}
				source.WriteString("---\n")
				source.Write(header)
				source.WriteString("---\n\n")
			}
			source.WriteString(configTemplate.Body)
			templates = upsertContentTemplate(templates, contentTemplate{
				ID: "config:" + configTemplate.Name, Name: configTemplate.Name,
				Description: "Site-configured starter", Directory: outputDir, Source: source.String(),
			})
		}
	}
	local, err := loadSiteContentTemplates(h.root, directory)
	if err != nil {
		return nil, err
	}
	// Site-authored Markdown is the highest-priority starter definition.
	// Duplicate names among Markdown files are errors, not silent overrides.
	seenLocal := make(map[string]string, len(local))
	for _, siteTemplate := range local {
		key := strings.ToLower(strings.TrimSpace(siteTemplate.Name))
		if previous, exists := seenLocal[key]; exists {
			return nil, fmt.Errorf("duplicate content starter %q in %s and %s", siteTemplate.Name, previous, siteTemplate.ID)
		}
		seenLocal[key] = siteTemplate.ID
		templates = upsertContentTemplate(templates, siteTemplate)
	}
	if err := validateContentTemplateNames(templates); err != nil {
		return nil, err
	}
	sort.SliceStable(templates, func(i, j int) bool {
		return strings.ToLower(templates[i].Name) < strings.ToLower(templates[j].Name)
	})
	return templates, nil
}

func contentTemplateKey(template contentTemplate) string {
	if strings.HasPrefix(template.ID, "builtin:") {
		return strings.TrimPrefix(template.ID, "builtin:")
	}
	return template.Name
}

func upsertContentTemplate(templates []contentTemplate, selected contentTemplate) []contentTemplate {
	for i := range templates {
		if strings.EqualFold(contentTemplateKey(templates[i]), selected.Name) {
			templates[i] = selected
			return templates
		}
	}
	return append(templates, selected)
}

func validateContentTemplateNames(templates []contentTemplate) error {
	owners := make(map[string]string, len(templates))
	for _, template := range templates {
		canonicalName := template.Name
		if strings.HasPrefix(template.ID, "builtin:") {
			canonicalName = strings.TrimPrefix(template.ID, "builtin:")
		}
		canonical := strings.ToLower(strings.TrimSpace(canonicalName))
		if owner, exists := owners[canonical]; exists && owner != template.ID {
			return fmt.Errorf("content starter name %q conflicts with another starter", template.Name)
		}
		owners[canonical] = template.ID
	}
	return nil
}

func builtinContentTemplates(today string) []contentTemplate {
	return []contentTemplate{
		{ID: "builtin:blog-post", Name: "Blog post", Description: "A long-form article with title, date, and tags.", Directory: "pages/blog", Source: fmt.Sprintf("---\ntemplateKey: blog-post\ntitle: \"{{title}}\"\ndate: %s\npublished: false\ntags: []\n---\n\n# {{title}}\n\nStart writing here.\n", today)},
		{ID: "builtin:thought", Name: "Thought", Description: "A short link or idea with a little context.", Directory: "pages/thoughts", Source: fmt.Sprintf("---\ntemplateKey: thought\ntitle: \"{{title}}\"\ndate: %s\npublished: false\ntags:\n  - thought\n---\n\nAdd a link or thought here.\n", today)},
		{ID: "builtin:note", Name: "Note", Description: "A compact note for the site.", Directory: "pages/notes", Source: fmt.Sprintf("---\ntemplateKey: note\ntitle: \"{{title}}\"\ndate: %s\npublished: false\ntags:\n  - note\n---\n\nWrite a note here.\n", today)},
		{ID: "builtin:photo", Name: "Photo", Description: "A photo post with an image path and caption.", Directory: "pages/photos", Source: fmt.Sprintf("---\ntemplateKey: photo\ntitle: \"{{title}}\"\ndate: %s\npublished: false\ntags:\n  - photo\nimage: /images/\n---\n\n![Describe the image](/images/)\n\nAdd a caption.\n", today)},
	}
}

func loadSiteContentTemplates(siteRoot, directory string) ([]contentTemplate, error) {
	absoluteSiteRoot, err := filepath.Abs(siteRoot)
	if err != nil {
		return nil, err
	}
	site, err := os.OpenRoot(absoluteSiteRoot)
	if err != nil {
		return nil, err
	}
	defer site.Close()
	templateRoot, err := site.OpenRoot(filepath.Clean(directory))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer templateRoot.Close()
	var templates []contentTemplate
	err = fs.WalkDir(templateRoot.FS(), ".", func(path string, entry fs.DirEntry, walkErr error) error {
		return appendSiteContentTemplate(templateRoot, path, entry, walkErr, &templates)
	})
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func appendSiteContentTemplate(templateRoot *os.Root, path string, entry fs.DirEntry, walkErr error, templates *[]contentTemplate) error {
	if walkErr != nil {
		return walkErr
	}
	if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
		return nil
	}
	info, err := entry.Info()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() || info.Size() > maxSource {
		return nil
	}
	template, err := readSiteContentTemplate(templateRoot, path)
	if err != nil {
		return err
	}
	if template != nil {
		*templates = append(*templates, *template)
	}
	return nil
}

func readSiteContentTemplate(templateRoot *os.Root, relative string) (*contentTemplate, error) {
	localPath := filepath.FromSlash(relative)
	file, err := templateRoot.Open(localPath)
	if err != nil {
		return nil, err
	}
	openedInfo, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if !openedInfo.Mode().IsRegular() || openedInfo.Size() > maxSource {
		_ = file.Close()
		return nil, nil
	}
	source, readErr := io.ReadAll(io.LimitReader(file, maxSource+1))
	closeErr := file.Close()
	if readErr != nil {
		return nil, readErr
	}
	if closeErr != nil {
		return nil, closeErr
	}
	if len(source) > maxSource {
		return nil, nil
	}
	name := strings.TrimSuffix(filepath.Base(localPath), filepath.Ext(localPath))
	return siteContentTemplateFromSource(relative, name, source)
}

func siteContentTemplateFromSource(relative, name string, source []byte) (*contentTemplate, error) {
	frontmatter, body := splitStarterSource(string(source))
	trimmed := strings.TrimPrefix(string(source), "\ufeff")
	if (strings.HasPrefix(trimmed, "---\n") || strings.HasPrefix(trimmed, "---\r\n")) && string(body) == string(source) {
		return nil, fmt.Errorf("starter %s is missing its closing frontmatter delimiter", filepath.ToSlash(relative))
	}
	if len(frontmatter) == 0 {
		name = strings.NewReplacer("-", " ", "_", " ").Replace(name)
		return &contentTemplate{ID: "local:" + filepath.ToSlash(relative), Name: name, Description: "Site template · " + filepath.ToSlash(relative), Directory: name, Source: string(source)}, nil
	}
	metadata, err := parseSiteStarterMetadata(relative, name, frontmatter)
	if err != nil {
		return nil, err
	}
	var rendered bytes.Buffer
	rendered.WriteString("---\n")
	rendered.Write(metadata.cleanFrontmatter)
	rendered.WriteString("---\n")
	rendered.Write(body)
	return &contentTemplate{ID: "local:" + filepath.ToSlash(relative), Name: metadata.name, Description: metadata.description, Directory: metadata.directory, Aliases: metadata.aliases, Source: rendered.String()}, nil
}

func parseSiteStarterMetadata(relative, name string, frontmatter []byte) (siteStarterMetadata, error) {
	var metadata map[string]any
	if err := yaml.Unmarshal(frontmatter, &metadata); err != nil {
		return siteStarterMetadata{}, fmt.Errorf("parse starter %s: %w", filepath.ToSlash(relative), err)
	}
	if metadata == nil {
		return siteStarterMetadata{}, fmt.Errorf("parse starter %s: frontmatter must be a mapping", filepath.ToSlash(relative))
	}
	starter, err := starterMetadataMap(relative, metadata)
	if err != nil {
		return siteStarterMetadata{}, err
	}
	if rawName, exists := starter["name"]; exists {
		value, ok := rawName.(string)
		if !ok || strings.TrimSpace(value) == "" {
			return siteStarterMetadata{}, fmt.Errorf("starter %s starter.name must be a non-empty string", filepath.ToSlash(relative))
		}
		name = strings.TrimSpace(value)
	}
	description, err := siteStarterDescription(relative, starter)
	if err != nil {
		return siteStarterMetadata{}, err
	}
	outputDirectory, err := siteStarterDirectory(relative, name, metadata)
	if err != nil {
		return siteStarterMetadata{}, err
	}
	aliases, err := siteStarterAliases(relative, starter)
	if err != nil {
		return siteStarterMetadata{}, err
	}
	delete(metadata, "starter")
	delete(metadata, "_directory")
	cleanFrontmatter, err := yaml.Marshal(metadata)
	if err != nil {
		return siteStarterMetadata{}, err
	}
	return siteStarterMetadata{name: name, description: description, directory: outputDirectory, aliases: aliases, cleanFrontmatter: cleanFrontmatter}, nil
}

func starterMetadataMap(relative string, metadata map[string]any) (map[string]any, error) {
	if rawStarter, exists := metadata["starter"]; exists {
		starter, ok := rawStarter.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("starter %s metadata must be a mapping", filepath.ToSlash(relative))
		}
		return starter, nil
	}
	return map[string]any{}, nil
}

func siteStarterDescription(relative string, starter map[string]any) (string, error) {
	description := "Site template · " + filepath.ToSlash(relative)
	if raw, exists := starter["description"]; exists {
		value, ok := raw.(string)
		if !ok {
			return "", fmt.Errorf("starter %s starter.description must be a string", filepath.ToSlash(relative))
		}
		description = value
	}
	return description, nil
}

func siteStarterDirectory(relative, name string, metadata map[string]any) (string, error) {
	directory := name
	if raw, exists := metadata["_directory"]; exists {
		value, ok := raw.(string)
		if !ok || strings.TrimSpace(value) == "" {
			return "", fmt.Errorf("starter %s has invalid _directory", filepath.ToSlash(relative))
		}
		directory = value
	}
	if !safeStarterOutputDirectory(directory) {
		return "", fmt.Errorf("starter %s has unsafe _directory %q", filepath.ToSlash(relative), directory)
	}
	return directory, nil
}

func siteStarterAliases(relative string, starter map[string]any) ([]string, error) {
	aliases := []string{}
	if raw, exists := starter["aliases"]; exists {
		values, ok := raw.([]any)
		if !ok {
			return nil, fmt.Errorf("starter %s aliases must be a list of strings", filepath.ToSlash(relative))
		}
		for _, item := range values {
			alias, ok := item.(string)
			if !ok || strings.TrimSpace(alias) == "" {
				return nil, fmt.Errorf("starter %s aliases must be a list of non-empty strings", filepath.ToSlash(relative))
			}
			aliases = append(aliases, strings.TrimSpace(alias))
		}
	}
	return aliases, nil
}

func safeStarterOutputDirectory(directory string) bool {
	if directory == "" || filepath.IsAbs(directory) {
		return false
	}
	clean := filepath.Clean(directory)
	return clean != ".." && !strings.HasPrefix(clean, ".."+string(filepath.Separator))
}

func splitStarterSource(source string) (frontmatter, body []byte) {
	trimmed := strings.TrimPrefix(source, "\ufeff")
	if !strings.HasPrefix(trimmed, "---\n") && !strings.HasPrefix(trimmed, "---\r\n") {
		return nil, []byte(source)
	}
	lines := strings.SplitAfter(trimmed, "\n")
	var fm bytes.Buffer
	for index := 1; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		if line == "---" || line == "..." {
			return fm.Bytes(), []byte(strings.Join(lines[index+1:], ""))
		}
		fm.WriteString(lines[index])
	}
	return nil, []byte(source)
}
