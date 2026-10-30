// Package themes provides embedded theme files for markata-go.
//
// This package embeds the default theme templates and static assets
// so they are always available regardless of where the binary is run.
package themes

import (
	"io/fs"
	"path"
	"strings"
)

// DefaultTemplates returns a filesystem containing the default theme templates.
func DefaultTemplates() fs.FS {
	sub, err := fs.Sub(defaultTheme, "default/templates")
	if err != nil {
		// This should never happen with embedded files
		return nil
	}
	return sub
}

// DefaultStatic returns a filesystem containing the default theme static files.
func DefaultStatic() fs.FS {
	sub, err := fs.Sub(defaultTheme, "default/static")
	if err != nil {
		// This should never happen with embedded files
		return nil
	}
	return sub
}

// DefaultTheme returns the full embedded filesystem for the default theme.
func DefaultTheme() fs.FS {
	sub, err := fs.Sub(defaultTheme, "default")
	if err != nil {
		return nil
	}
	return sub
}

// ReadTemplate reads a template file from the embedded default theme.
func ReadTemplate(name string) ([]byte, error) {
	return defaultTheme.ReadFile(path.Join("default", "templates", name))
}

// ReadStatic reads a static file from the embedded default theme.
func ReadStatic(name string) ([]byte, error) {
	return defaultTheme.ReadFile(path.Join("default", "static", name))
}

// ListTemplates returns all template files in the default theme.
func ListTemplates() ([]string, error) {
	var templates []string
	const prefix = "default/templates/"
	err := fs.WalkDir(defaultTheme, "default/templates", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// Strip the prefix to get relative path (embedded FS always uses forward slashes)
			if strings.HasPrefix(p, prefix) {
				templates = append(templates, p[len(prefix):])
			}
		}
		return nil
	})
	return templates, err
}

// ListStatic returns all static files in the default theme.
func ListStatic() ([]string, error) {
	var files []string
	const prefix = "default/static/"
	err := fs.WalkDir(defaultTheme, "default/static", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// Strip the prefix to get relative path (embedded FS always uses forward slashes)
			if strings.HasPrefix(p, prefix) {
				files = append(files, p[len(prefix):])
			}
		}
		return nil
	})
	return files, err
}
