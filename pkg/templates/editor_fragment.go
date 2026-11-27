package templates

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/flosch/pongo2/v6"
)

// editorFragmentLoader is deliberately separate from the unrestricted build
// loader. A preview request may never use a template to read outside the
// selected template directory, including through symlinks.
type editorFragmentLoader struct {
	root *os.Root
}

const maxEditorFragmentBytes = 200000

func validEditorTemplateName(name string) bool {
	return fs.ValidPath(name) && name != "." && !strings.Contains(name, "\\") && filepath.IsLocal(filepath.FromSlash(name))
}

func (l editorFragmentLoader) Abs(base, name string) string {
	if !validEditorTemplateName(name) {
		return ""
	}
	if base != "" {
		candidate := path.Join(path.Dir(base), name)
		if info, err := l.root.Stat(filepath.FromSlash(candidate)); err == nil && info.Mode().IsRegular() {
			return candidate
		}
	}
	return name
}

func (l editorFragmentLoader) Get(name string) (io.Reader, error) {
	if !validEditorTemplateName(name) {
		return nil, fmt.Errorf("invalid editor template name")
	}
	file, err := l.root.Open(filepath.FromSlash(name))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("editor template is not a regular file")
	}
	data, err := io.ReadAll(io.LimitReader(file, maxEditorFragmentBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxEditorFragmentBytes {
		return nil, fmt.Errorf("editor template exceeds size limit")
	}
	return bytes.NewReader(data), nil
}

// RenderEditorFragment renders a site-owned companion using the production
// post projection and filters, but no lifecycle manager, callable build
// helpers, or arbitrary config.Extra (which can contain credentials).
// Draft posts are not entered into the global build template caches.
func RenderEditorFragment(root *os.Root, name string, post *models.Post, cfg *models.Config) (string, error) {
	registerFilters()
	set := pongo2.NewSet("editor-fragment", editorFragmentLoader{root: root})
	tpl, err := set.FromFile(name)
	if err != nil {
		return "", err
	}
	postMap := postToMapUncached(post)
	publicConfig := map[string]interface{}{}
	if cfg != nil {
		publicConfig["title"] = cfg.Title
		publicConfig["url"] = cfg.URL
		publicConfig["description"] = cfg.Description
		publicConfig["author"] = cfg.Author
		publicConfig["theme"] = ThemeToMap(&cfg.Theme)
	}
	ctx := pongo2.Context{"post": postMap, "config": publicConfig}
	addPostContext(&ctx, postMap, post, nil)
	return tpl.Execute(ctx)
}
