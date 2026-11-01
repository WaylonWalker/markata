package workbench

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type contentTemplate struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

func (h *Handler) contentTemplates() ([]contentTemplate, error) {
	today := time.Now().Format("2006-01-02")
	builtins := builtinContentTemplates(today)
	templates := append([]contentTemplate(nil), builtins...)
	local, err := loadSiteContentTemplates(h.root)
	if err != nil {
		return nil, err
	}
	templates = append(templates, local...)
	if len(templates) > len(builtins) {
		builtinCount := len(builtins)
		sort.Slice(templates[builtinCount:], func(i, j int) bool {
			return templates[builtinCount+i].Name < templates[builtinCount+j].Name
		})
	}
	return templates, nil
}

func builtinContentTemplates(today string) []contentTemplate {
	return []contentTemplate{
		{ID: "builtin:blog-post", Name: "Blog post", Description: "A long-form article with title, date, and tags.", Source: fmt.Sprintf("---\ntemplateKey: blog-post\ntitle: Untitled\ndate: %s\npublished: false\ntags: []\n---\n\n# Untitled\n\nStart writing here.\n", today)},
		{ID: "builtin:thought", Name: "Thought", Description: "A short link or idea with a little context.", Source: fmt.Sprintf("---\ntemplateKey: thought\ntitle: Untitled thought\ndate: %s\npublished: false\ntags:\n  - thought\n---\n\nAdd a link or thought here.\n", today)},
		{ID: "builtin:note", Name: "Note", Description: "A compact note for the site.", Source: fmt.Sprintf("---\ntemplateKey: note\ntitle: Untitled note\ndate: %s\npublished: false\ntags:\n  - note\n---\n\nWrite a note here.\n", today)},
		{ID: "builtin:photo", Name: "Photo", Description: "A photo post with an image path and caption.", Source: fmt.Sprintf("---\ntemplateKey: photo\ntitle: Untitled photo\ndate: %s\npublished: false\ntags:\n  - photo\nimage: /images/\n---\n\n![Describe the image](/images/)\n\nAdd a caption.\n", today)},
	}
}

func loadSiteContentTemplates(siteRoot string) ([]contentTemplate, error) {
	absoluteSiteRoot, err := filepath.Abs(siteRoot)
	if err != nil {
		return nil, err
	}
	site, err := os.OpenRoot(absoluteSiteRoot)
	if err != nil {
		return nil, err
	}
	defer site.Close()
	templateRoot, err := site.OpenRoot("templates")
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
	name = strings.NewReplacer("-", " ", "_", " ").Replace(name)
	return &contentTemplate{ID: "local:" + filepath.ToSlash(relative), Name: name, Description: "Site template · " + filepath.ToSlash(relative), Source: string(source)}, nil
}
