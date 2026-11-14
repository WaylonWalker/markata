package workbench

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type contentTemplate struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Directory         string   `json:"directory"`
	Aliases           []string `json:"aliases,omitempty"`
	NumberPrefix      string   `json:"number_prefix,omitempty"`
	NumberedByDefault bool     `json:"numbered_by_default,omitempty"`
	NextNumber        int      `json:"next_number,omitempty"`
	Source            string   `json:"source"`
}

type siteStarterMetadata struct {
	name              string
	description       string
	directory         string
	aliases           []string
	numberPrefix      string
	numberedByDefault bool
	cleanFrontmatter  []byte
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
	for i := range templates {
		if templates[i].NumberPrefix == "" {
			continue
		}
		next, err := nextContentNumber(h.root, templates[i].Directory, templates[i].NumberPrefix)
		if err != nil {
			return nil, fmt.Errorf("find next number for starter %q: %w", templates[i].Name, err)
		}
		templates[i].NextNumber = next
	}
	if len(templates) > len(builtins) {
		builtinCount := len(builtins)
		sort.Slice(templates[builtinCount:], func(i, j int) bool {
			return strings.ToLower(templates[builtinCount+i].Name) < strings.ToLower(templates[builtinCount+j].Name)
		})
	}
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
		{ID: "builtin:post", Name: "Post", Description: "A standard blog post.", Directory: "pages/post", Source: "---\ntemplate: post\n---\n\nWrite your content here...\n"},
		{ID: "builtin:page", Name: "Page", Description: "A static site page.", Directory: "pages", Source: "---\ntemplate: page\n---\n\nWrite your page content here...\n"},
		{ID: "builtin:docs", Name: "Docs", Description: "A documentation page.", Directory: "docs", Source: "---\ntemplate: docs\n---\n\nWrite your documentation here...\n"},
		{ID: "builtin:article", Name: "Article", Description: "A long-form article.", Directory: "pages/article", Aliases: []string{"blog-post", "essay", "tutorial"}, Source: "---\ntemplate: article\n---\n\nWrite your article here...\n"},
		{ID: "builtin:blog-post", Name: "Blog post", Description: "A long-form article with title, date, and tags.", Directory: "pages/blog", Source: fmt.Sprintf("---\ntemplateKey: blog-post\ntitle: \"{{title}}\"\ndate: %s\npublished: false\ntags: []\n---\n\n# {{title}}\n\nStart writing here.\n", today)},
		{ID: "builtin:note", Name: "Note", Description: "A quick note.", Directory: "pages/note", Aliases: []string{"thought", "status", "tweet"}, Source: "---\ntemplate: note\n---\n\nA quick note...\n"},
		{ID: "builtin:ping", Name: "Ping", Description: "A short ping post.", Directory: "pages/ping", NumberPrefix: "ping", NumberedByDefault: true, Source: "---\ntemplate: ping\ntemplateKey: ping\n---\n\nA quick ping...\n"},
		{ID: "builtin:thought", Name: "Thought", Description: "A short link or idea with a little context.", Directory: "pages/thoughts", Source: fmt.Sprintf("---\ntemplateKey: thought\ntitle: \"{{title}}\"\ndate: %s\npublished: false\ntags:\n  - thought\n---\n\nAdd a link or thought here.\n", today)},
		{ID: "builtin:photo", Name: "Photo", Description: "A photo post with an image path and caption.", Directory: "pages/photo", Aliases: []string{"shots", "image", "gallery"}, Source: "---\ntemplate: photo\nimage: \"\"\n---\n\nPhoto caption...\n"},
		{ID: "builtin:shot", Name: "Shot", Description: "A shot post with an image.", Directory: "pages/photo", Source: "---\ntemplate: shot\ntemplateKey: shot\nimage: \"\"\n---\n\nAdd a caption.\n"},
		{ID: "builtin:video", Name: "Video", Description: "A video post.", Directory: "pages/video", Aliases: []string{"clip", "cast", "stream"}, Source: "---\ntemplate: video\nvideo: \"\"\nimage: \"\"\nduration: \"\"\n---\n\nVideo description...\n"},
		{ID: "builtin:link", Name: "Link", Description: "Save a web link with a title and source URL.", Directory: "pages/link", Aliases: []string{"bookmark", "til", "stars"}, Source: "---\ntemplate: link\nurl: \"\"\nimage: \"\"\n---\n\nWhy I'm sharing this link...\n"},
		{ID: "builtin:quote", Name: "Quote", Description: "A quote with attribution.", Directory: "pages/quote", Aliases: []string{"quotation"}, Source: "---\ntemplate: quote\nquote: \"\"\nsource: \"\"\n---\n\nAdditional commentary on this quote...\n"},
		{ID: "builtin:guide", Name: "Guide", Description: "A step-by-step guide.", Directory: "pages/guide", Aliases: []string{"series", "step", "chapter"}, Source: "---\ntemplate: guide\n---\n\n## Introduction\n\nWrite your guide here...\n"},
		{ID: "builtin:inline", Name: "Inline", Description: "A short, feed-first post.", Directory: "pages/inline", Aliases: []string{"micro"}, Source: "---\ntemplate: inline\n---\n\nInline content...\n"},
		{ID: "builtin:gratitude", Name: "Gratitude", Description: "A short gratitude post.", Directory: "pages/inline", NumberPrefix: "gratitude", NumberedByDefault: true, Source: "---\ntemplate: gratitude\ntemplateKey: gratitude\n---\n\nI am grateful for...\n"},
		{ID: "builtin:contact", Name: "Contact", Description: "A contact or profile page.", Directory: "pages/contact", Aliases: []string{"character", "person"}, Source: "---\ntemplate: contact\nhandle: \"\"\naliases: []\navatar: \"\"\nurl: \"\"\n---\n\nBio or contact details...\n"},
		{ID: "builtin:author", Name: "Author", Description: "An author profile.", Directory: "pages/author", Source: "---\ntemplate: author\nname: \"\"\nbio: \"\"\nrole: \"\"\navatar: \"\"\nurl: \"\"\nemail: \"\"\nsocial: {}\nactive: true\nguest: false\n---\n\nExtended author bio...\n"},
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
	if (strings.HasPrefix(trimmed, "---\n") || strings.HasPrefix(trimmed, "---\r\n")) && bytes.Equal(body, source) {
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
	return &contentTemplate{
		ID: "local:" + filepath.ToSlash(relative), Name: metadata.name,
		Description: metadata.description, Directory: metadata.directory, Aliases: metadata.aliases,
		NumberPrefix: metadata.numberPrefix, NumberedByDefault: metadata.numberedByDefault,
		Source: rendered.String(),
	}, nil
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
	numberPrefix, numberedByDefault, err := siteStarterNumbering(relative, starter)
	if err != nil {
		return siteStarterMetadata{}, err
	}
	delete(metadata, "starter")
	delete(metadata, "_directory")
	cleanFrontmatter, err := yaml.Marshal(metadata)
	if err != nil {
		return siteStarterMetadata{}, err
	}
	return siteStarterMetadata{
		name: name, description: description, directory: outputDirectory, aliases: aliases,
		numberPrefix: numberPrefix, numberedByDefault: numberedByDefault, cleanFrontmatter: cleanFrontmatter,
	}, nil
}

func siteStarterNumbering(relative string, starter map[string]any) (string, bool, error) {
	prefix := ""
	if raw, exists := starter["number_prefix"]; exists {
		value, ok := raw.(string)
		if !ok || !validNumberPrefix(value) {
			return "", false, fmt.Errorf("starter %s starter.number_prefix must be a lowercase slug", filepath.ToSlash(relative))
		}
		prefix = value
	}
	numberedByDefault := false
	if raw, exists := starter["numbered_by_default"]; exists {
		value, ok := raw.(bool)
		if !ok {
			return "", false, fmt.Errorf("starter %s starter.numbered_by_default must be a boolean", filepath.ToSlash(relative))
		}
		numberedByDefault = value
	}
	if numberedByDefault && prefix == "" {
		return "", false, fmt.Errorf("starter %s must set starter.number_prefix when numbering is enabled", filepath.ToSlash(relative))
	}
	return prefix, numberedByDefault, nil
}

var numberPrefixPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func validNumberPrefix(prefix string) bool {
	return numberPrefixPattern.MatchString(prefix)
}

func nextContentNumber(siteRoot, directory, prefix string) (int, error) {
	if !safeStarterOutputDirectory(directory) || !validNumberPrefix(prefix) {
		return 0, fmt.Errorf("invalid numbering directory or prefix")
	}
	root, err := filepath.Abs(siteRoot)
	if err != nil {
		return 0, err
	}
	folder := filepath.Join(root, filepath.FromSlash(directory))
	resolved, err := filepath.EvalSymlinks(folder)
	if os.IsNotExist(err) {
		return 1, nil
	}
	if err != nil {
		return 0, err
	}
	relative, err := filepath.Rel(root, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return 0, fmt.Errorf("numbering directory escapes workspace")
	}
	entries, err := os.ReadDir(resolved)
	if err != nil {
		return 0, err
	}
	pattern := regexp.MustCompile(`^` + regexp.QuoteMeta(prefix) + `-(\d+)\.(?:md|markdown|mdown|mkd)$`)
	maximum := 0
	for _, entry := range entries {
		match := pattern.FindStringSubmatch(entry.Name())
		if len(match) != 2 || entry.IsDir() {
			continue
		}
		number, parseErr := strconv.Atoi(match[1])
		if parseErr == nil && number > maximum {
			maximum = number
		}
	}
	return maximum + 1, nil
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
