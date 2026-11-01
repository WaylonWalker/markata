package workbench

import (
	"fmt"
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
	templates := []contentTemplate{
		{ID: "builtin:blog-post", Name: "Blog post", Description: "A long-form article with title, date, and tags.", Source: fmt.Sprintf("---\ntemplateKey: blog-post\ntitle: Untitled\ndate: %s\npublished: false\ntags: []\n---\n\n# Untitled\n\nStart writing here.\n", today)},
		{ID: "builtin:thought", Name: "Thought", Description: "A short link or idea with a little context.", Source: fmt.Sprintf("---\ntemplateKey: thought\ntitle: Untitled thought\ndate: %s\npublished: false\ntags:\n  - thought\n---\n\nAdd a link or thought here.\n", today)},
		{ID: "builtin:note", Name: "Note", Description: "A compact note for the site.", Source: fmt.Sprintf("---\ntemplateKey: note\ntitle: Untitled note\ndate: %s\npublished: false\ntags:\n  - note\n---\n\nWrite a note here.\n", today)},
		{ID: "builtin:photo", Name: "Photo", Description: "A photo post with an image path and caption.", Source: fmt.Sprintf("---\ntemplateKey: photo\ntitle: Untitled photo\ndate: %s\npublished: false\ntags:\n  - photo\nimage: /images/\n---\n\n![Describe the image](/images/)\n\nAdd a caption.\n", today)},
	}
	root := filepath.Join(h.root, "templates")
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return templates, nil
	} else if err != nil {
		return nil, err
	}
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	relativeRoot, err := filepath.Rel(h.root, resolvedRoot)
	if err != nil || relativeRoot == ".." || strings.HasPrefix(relativeRoot, ".."+string(filepath.Separator)) {
		return templates, nil
	}
	root = resolvedRoot
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
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
		if info.Size() > maxSource {
			return nil
		}
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		name := strings.TrimSuffix(filepath.Base(relative), filepath.Ext(relative))
		name = strings.NewReplacer("-", " ", "_", " ").Replace(name)
		templates = append(templates, contentTemplate{ID: "local:" + filepath.ToSlash(relative), Name: name, Description: "Site template · " + filepath.ToSlash(relative), Source: string(source)})
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(templates) > 4 {
		sort.Slice(templates[4:], func(i, j int) bool { return templates[4+i].Name < templates[4+j].Name })
	}
	return templates, nil
}
