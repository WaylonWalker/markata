// Package editorworkspace discovers and loads a trusted local editor workspace.
package editorworkspace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

// File describes a discovered source document.
type File struct {
	Path  string   `json:"path"`
	Title *string  `json:"title,omitempty"`
	Tags  []string `json:"tags,omitempty"`
	Image string   `json:"image,omitempty"`
}

// Workspace confines discovery and reads to one directory.
type Workspace struct{ root string }

// New validates the workspace root without changing the process directory.
func New(root string) (*Workspace, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	absolute, err = filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(absolute)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("workspace must be a directory")
	}
	return &Workspace{root: absolute}, nil
}

// Config loads the selected site's configuration, or returns nil if none exists.
func (w *Workspace) Config() (*models.Config, error) {
	if !runtimeenv.DotEnvDisabled() {
		if err := config.LoadDotEnvFile(filepath.Join(w.root, ".env")); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("load workspace environment: %w", err)
		}
	}
	for _, name := range []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"} {
		filename := filepath.Join(w.root, name)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			continue
		}
		return config.Load(filename)
	}
	return nil, nil
}

func (w *Workspace) manager() (*lifecycle.Manager, error) {
	cfg, err := w.Config()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = models.NewConfig()
	}
	m := lifecycle.NewManager()
	lc := lifecycle.NewConfig()
	lc.ContentDir = w.root
	lc.GlobPatterns = cfg.GlobConfig.Patterns
	lc.Extra["models_config"] = cfg
	lc.Extra["use_gitignore"] = cfg.GlobConfig.UseGitignore
	m.SetConfig(lc)
	glob := plugins.NewGlobPlugin()
	if err := glob.Configure(m); err != nil {
		return nil, err
	}
	if err := glob.Glob(m); err != nil {
		return nil, err
	}
	safe := []string{}
	for _, path := range m.Files() {
		if _, err := w.resolve(path); err == nil {
			safe = append(safe, path)
		}
	}
	m.SetFiles(safe)
	return m, nil
}

func (w *Workspace) resolve(path string) (string, error) {
	if filepath.IsAbs(path) || !filepath.IsLocal(filepath.FromSlash(path)) {
		return "", fmt.Errorf("invalid workspace path")
	}
	resolved, err := filepath.EvalSymlinks(filepath.Join(w.root, filepath.FromSlash(path)))
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(w.root, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes workspace")
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("not a regular source file")
	}
	return resolved, nil
}

// ReadImage returns a workspace static image for the editor preview.
func (w *Workspace) ReadImage(name string) ([]byte, error) {
	if filepath.IsAbs(name) || !filepath.IsLocal(filepath.FromSlash(name)) {
		return nil, fmt.Errorf("invalid image path")
	}
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif", ".avif":
	default:
		return nil, fmt.Errorf("unsupported image type")
	}
	staticRoot := filepath.Join(w.root, "static")
	resolved, err := w.resolve(filepath.ToSlash(filepath.Join("static", filepath.FromSlash(name))))
	if err != nil {
		return nil, err
	}
	relative, err := filepath.Rel(staticRoot, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("image escapes static directory")
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return nil, err
	}
	if info.Size() > 15<<20 {
		return nil, fmt.Errorf("image exceeds 15 MiB")
	}
	return os.ReadFile(resolved)
}

// Paths discovers source paths without reading or parsing document bodies.
func (w *Workspace) Paths() ([]File, error) {
	m, err := w.manager()
	if err != nil {
		return nil, err
	}
	files := make([]File, 0, len(m.Files()))
	for _, path := range m.Files() {
		files = append(files, File{Path: filepath.ToSlash(path)})
	}
	return files, nil
}

// Files runs Markata's glob and load plugins without rendering or writing a site.
func (w *Workspace) Files() ([]File, error) {
	m, err := w.manager()
	if err != nil {
		return nil, err
	}
	loader := plugins.NewLoadPlugin()
	if err := loader.Configure(m); err != nil {
		return nil, err
	}
	if err := loader.Load(m); err != nil {
		return nil, err
	}
	files := make([]File, 0, len(m.Posts()))
	for _, post := range m.Posts() {
		image := ""
		if value, ok := post.Extra["image"].(string); ok {
			image = value
		}
		files = append(files, File{Path: filepath.ToSlash(post.Path), Title: post.Title, Tags: post.Tags, Image: image})
	}
	return files, nil
}

// Read returns bounded, original source, including frontmatter, for a discovered file.
func (w *Workspace) Read(path string) (string, error) {
	m, err := w.manager()
	if err != nil {
		return "", err
	}
	found := false
	for _, candidate := range m.Files() {
		if filepath.ToSlash(candidate) == path {
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("source file is not in the workspace content globs")
	}
	return w.readSource(path)
}

func (w *Workspace) readSource(path string) (string, error) {
	resolved, err := w.resolve(path)
	if err != nil {
		return "", err
	}
	file, err := os.Open(resolved)
	if err != nil {
		return "", err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, 200001))
	if err != nil {
		return "", err
	}
	if len(data) > 200000 {
		return "", fmt.Errorf("source exceeds 200,000 bytes")
	}
	return string(data), nil
}
