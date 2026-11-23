// Package editorworkspace discovers and loads a trusted local editor workspace.
package editorworkspace

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/runtimeenv"
)

// File describes a discovered source document.
type File struct {
	Path   string   `json:"path"`
	Slug   string   `json:"slug,omitempty"`
	Title  *string  `json:"title,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Image  string   `json:"image,omitempty"`
	Date   string   `json:"date,omitempty"`
	Status string   `json:"status,omitempty"`
}

// fileDate formats a post's publication date as YYYY-MM-DD.
func fileDate(post *models.Post) string {
	if post.Date == nil {
		return ""
	}
	return post.Date.Format("2006-01-02")
}

// fileStatus reports whether a post is a draft or published, or empty when neither.
func fileStatus(post *models.Post) string {
	switch {
	case post.Draft:
		return "draft"
	case post.Published:
		return "published"
	default:
		return ""
	}
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
	m, _, err := w.managerWithIgnoreRules()
	return m, err
}

func (w *Workspace) managerWithIgnoreRules() (*lifecycle.Manager, func(string) bool, error) {
	cfg, err := w.Config()
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}
	if err := glob.Glob(m); err != nil {
		return nil, nil, err
	}
	safe := []string{}
	for _, path := range m.Files() {
		if _, err := w.resolve(path); err == nil {
			safe = append(safe, path)
		}
	}
	m.SetFiles(safe)
	return m, glob.IsIgnoredPath, nil
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
	const maxImageBytes = 15 << 20
	if info.Size() > maxImageBytes {
		return nil, fmt.Errorf("image exceeds 15 MiB")
	}
	file, err := os.Open(resolved)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if !openedInfo.Mode().IsRegular() {
		return nil, fmt.Errorf("not a regular image file")
	}
	data, err := io.ReadAll(io.LimitReader(file, maxImageBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxImageBytes {
		return nil, fmt.Errorf("image exceeds 15 MiB")
	}
	return data, nil
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
		files = append(files, File{Path: filepath.ToSlash(post.Path), Slug: post.Slug, Title: post.Title, Tags: post.Tags, Image: image, Date: fileDate(post), Status: fileStatus(post)})
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

// Write atomically replaces a discovered workspace source file.
func (w *Workspace) Write(path, source string) error {
	if len(source) > 200000 {
		return fmt.Errorf("source exceeds 200,000 bytes")
	}
	m, err := w.manager()
	if err != nil {
		return err
	}
	for _, candidate := range m.Files() {
		if filepath.ToSlash(candidate) == path {
			return w.writeSource(path, source)
		}
	}
	return fmt.Errorf("source file is not in the workspace content globs")
}

// Delete removes only an existing regular Markdown source whose bytes still
// match the reviewed version. Symlinked paths are refused, including symlinked
// parent directories. Callers must separately verify Git's untracked status.
func (w *Workspace) Delete(path, expectedVersion string) error {
	local := filepath.FromSlash(path)
	if !filepath.IsLocal(local) || !isMarkdownSource(local) {
		return fmt.Errorf("invalid workspace Markdown path")
	}
	target := filepath.Join(w.root, local)
	resolved, err := w.resolve(path)
	if err != nil {
		return err
	}
	if resolved != target {
		return fmt.Errorf("cannot delete source through symlinked path")
	}
	info, err := os.Lstat(target)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("cannot delete non-regular workspace source")
	}
	source, err := w.readSource(path)
	if err != nil {
		return err
	}
	if SourceVersion(source) != expectedVersion {
		return fmt.Errorf("workspace source changed since review")
	}
	// Local workbench CAS protects its own requests; external filesystem
	// processes are not covered by an atomic cross-process transaction.
	return os.Remove(target)
}

// Create atomically adds a Markdown file inside an existing workspace directory.
func (w *Workspace) Create(path, source string) error {
	if len(source) > 200000 {
		return fmt.Errorf("source exceeds 200,000 bytes")
	}
	local := filepath.FromSlash(path)
	if !filepath.IsLocal(local) || !isMarkdownSource(local) {
		return fmt.Errorf("new workspace path must be a local Markdown file")
	}
	local = filepath.Clean(local)
	parent, err := w.ensureDirectory(filepath.Dir(local))
	if err != nil {
		return err
	}
	target := filepath.Join(parent, filepath.Base(local))
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("workspace file already exists")
	} else if !os.IsNotExist(err) {
		return err
	}
	return writeNewSource(parent, target, source)
}

// CheckCreate verifies a new Markdown path is available without writing its source.
func (w *Workspace) CheckCreate(path string) error {
	local := filepath.FromSlash(path)
	if !filepath.IsLocal(local) || !isMarkdownSource(local) {
		return fmt.Errorf("new workspace path must be a local Markdown file")
	}
	local = filepath.Clean(local)
	parent, err := w.ensureDirectory(filepath.Dir(local))
	if err != nil {
		return err
	}
	return requireNewFile(filepath.Join(parent, filepath.Base(local)))
}

// UploadedImage is a validated photo supplied alongside a new Markdown post.
// Filename is intentionally not used for paths: the server generates safe slugs.
type UploadedImage struct {
	ContentType string
	Data        []byte
}

// SourceWithUploadedImages adds a cover image and an explicit Markdown gallery.
// Only server-generated image URLs should be passed to this function.
func SourceWithUploadedImages(source string, urls []string) string {
	if len(urls) == 0 {
		return source
	}
	result := sourceWithImage(source, urls[0])
	for i, imageURL := range urls {
		result += fmt.Sprintf("\n\n![Shared photo %d](%s)", i+1, imageURL)
	}
	return result + "\n"
}

// CreateWithImages creates a Markdown file and 2-8 local photos as one
// rollback-safe transaction. Existing files are never overwritten.
func (w *Workspace) CreateWithImages(path, source, imageName string, images []UploadedImage) ([]string, error) {
	if len(images) < 2 || len(images) > 8 {
		return nil, fmt.Errorf("provide 2 to 8 images")
	}
	if len(source) > 200000 || !validUploadName(imageName) {
		return nil, fmt.Errorf("invalid Markdown source or image slug")
	}
	local := filepath.Clean(filepath.FromSlash(path))
	if !filepath.IsLocal(local) || !isMarkdownSource(local) {
		return nil, fmt.Errorf("new workspace path must be a local Markdown file")
	}
	sourceParent, err := w.ensureDirectory(filepath.Dir(local))
	if err != nil {
		return nil, err
	}
	imageParent, err := w.ensureDirectory(filepath.Join("static", "images"))
	if err != nil {
		return nil, err
	}
	sourceTarget := filepath.Join(sourceParent, filepath.Base(local))
	if err := requireNewFile(sourceTarget); err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(images))
	targets := make([]string, 0, len(images))
	temps := make([]string, 0, len(images))
	defer func() {
		for _, tmp := range temps {
			_ = os.Remove(tmp)
		}
	}()
	total := 0
	for i, photo := range images {
		total += len(photo.Data)
		if len(photo.Data) == 0 || len(photo.Data) > 15<<20 || total > 48<<20 {
			return nil, fmt.Errorf("photos must be nonempty, at most 15 MiB each and 48 MiB total")
		}
		extension, err := ImageExtension(photo.ContentType, photo.Data)
		if err != nil {
			return nil, err
		}
		filename := fmt.Sprintf("%s-%d.%s", imageName, i+1, extension)
		target := filepath.Join(imageParent, filename)
		if err := requireNewFile(target); err != nil {
			return nil, err
		}
		tmp, err := writeTempBytes(imageParent, ".markata-image-*", photo.Data)
		if err != nil {
			return nil, err
		}
		temps = append(temps, tmp)
		targets = append(targets, target)
		urls = append(urls, "/images/"+filename)
	}
	finalSource := SourceWithUploadedImages(source, urls)
	if len(finalSource) > 200000 {
		return nil, fmt.Errorf("source exceeds 200,000 bytes")
	}
	sourceTemp, err := writeTempBytes(sourceParent, ".markata-editor-*", []byte(finalSource))
	if err != nil {
		return nil, err
	}
	temps = append(temps, sourceTemp)
	linked := make([]string, 0, len(images))
	completed := false
	defer func() {
		if !completed {
			for _, target := range linked {
				_ = os.Remove(target)
			}
		}
	}()
	for i, target := range targets {
		if err := os.Link(temps[i], target); err != nil {
			return nil, fmt.Errorf("create workspace image: %w", err)
		}
		linked = append(linked, target)
	}
	if err := os.Link(sourceTemp, sourceTarget); err != nil {
		return nil, fmt.Errorf("create workspace source: %w", err)
	}
	completed = true
	return urls, nil
}

// CreateWithImage atomically creates a Markdown file and a local static image.
// Existing source or image paths are never replaced.
func (w *Workspace) CreateWithImage(path, source, imageName, contentType string, image []byte) (string, error) {
	if len(source) > 200000 {
		return "", fmt.Errorf("source exceeds 200,000 bytes")
	}
	if len(image) == 0 || len(image) > 15<<20 {
		return "", fmt.Errorf("image must be between 1 byte and 15 MiB")
	}
	local := filepath.FromSlash(path)
	if !filepath.IsLocal(local) || !isMarkdownSource(local) {
		return "", fmt.Errorf("new workspace path must be a local Markdown file")
	}
	local = filepath.Clean(local)
	if !validUploadName(imageName) {
		return "", fmt.Errorf("image name must be a lowercase slug")
	}
	extension, err := ImageExtension(contentType, image)
	if err != nil {
		return "", err
	}
	sourceParent, err := w.ensureDirectory(filepath.Dir(local))
	if err != nil {
		return "", err
	}
	imageParent, err := w.ensureDirectory(filepath.Join("static", "images"))
	if err != nil {
		return "", err
	}
	sourceTarget := filepath.Join(sourceParent, filepath.Base(local))
	imageFilename := imageName + "." + extension
	imageTarget := filepath.Join(imageParent, imageFilename)
	if err := requireNewFile(sourceTarget); err != nil {
		return "", err
	}
	if err := requireNewFile(imageTarget); err != nil {
		return "", fmt.Errorf("image file already exists: %s", imageFilename)
	}
	imageURL := "/images/" + imageFilename
	source = sourceWithImage(source, imageURL)
	if len(source) > 200000 {
		return "", fmt.Errorf("source exceeds 200,000 bytes")
	}
	imageTemp, err := writeTempBytes(imageParent, ".markata-image-*", image)
	if err != nil {
		return "", err
	}
	defer os.Remove(imageTemp)
	sourceTemp, err := writeTempBytes(sourceParent, ".markata-editor-*", []byte(source))
	if err != nil {
		return "", err
	}
	defer os.Remove(sourceTemp)
	if err := os.Link(imageTemp, imageTarget); err != nil {
		return "", fmt.Errorf("create workspace image: %w", err)
	}
	if err := os.Link(sourceTemp, sourceTarget); err != nil {
		removeErr := os.Remove(imageTarget)
		if removeErr != nil {
			return "", errors.Join(
				fmt.Errorf("create workspace source: %w", err),
				fmt.Errorf("remove new image: %w", removeErr),
			)
		}
		return "", fmt.Errorf("create workspace source: %w", err)
	}
	return imageURL, nil
}

// SetContentTemplateDefaults saves a starter's site-wide folder and filename mode.
func (w *Workspace) SetContentTemplateDefaults(name, directory, filenameMode string) error {
	if !validUploadName(name) || !filepath.IsLocal(filepath.FromSlash(directory)) || filenameMode != "title" && filenameMode != "numbered" {
		return fmt.Errorf("invalid content template defaults")
	}
	cleanDirectory := filepath.Clean(filepath.FromSlash(directory))
	if cleanDirectory == ".." || strings.HasPrefix(cleanDirectory, ".."+string(filepath.Separator)) {
		return fmt.Errorf("content template directory must stay inside the workspace")
	}
	configName, err := w.contentTemplateConfigName()
	if err != nil {
		return err
	}
	return w.writeContentTemplateDefaults(configName, name, filepath.ToSlash(cleanDirectory), filenameMode)
}

func (w *Workspace) contentTemplateConfigName() (string, error) {
	for _, candidate := range []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"} {
		if _, err := os.Lstat(filepath.Join(w.root, candidate)); err == nil {
			return candidate, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("inspect site configuration: %w", err)
		}
	}
	return "markata-go.toml", nil
}

func (w *Workspace) writeContentTemplateDefaults(configName, name, directory, filenameMode string) error {
	target := filepath.Join(w.root, configName)
	existing, err := os.Lstat(target)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("inspect site configuration: %w", err)
	}
	exists := err == nil
	mode := os.FileMode(0o644)
	var source []byte
	if exists {
		if existing.Mode()&os.ModeSymlink != 0 || !existing.Mode().IsRegular() {
			return fmt.Errorf("site configuration must be a regular file")
		}
		mode = existing.Mode().Perm()
		source, err = os.ReadFile(target)
		if err != nil {
			return fmt.Errorf("read site configuration: %w", err)
		}
	}
	return writeContentTemplateConfig(w.root, configName, target, name, directory, filenameMode, source, mode, exists)
}

func writeContentTemplateConfig(root, configName, target, name, directory, filenameMode string, source []byte, mode os.FileMode, exists bool) error {
	temporary, err := os.CreateTemp(root, ".markata-settings-*"+filepath.Ext(configName))
	if err != nil {
		return fmt.Errorf("create site settings temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	err = temporary.Chmod(mode)
	if err == nil {
		_, err = temporary.Write(source)
	}
	if err == nil {
		err = temporary.Sync()
	}
	if closeErr := temporary.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("write site settings temporary file: %w", err)
	}
	placementKey := "content_templates.placement." + name
	modeKey := "content_templates.filename_mode." + name
	if err := config.SetValueInFileDirect(temporaryPath, placementKey, directory); err != nil {
		return fmt.Errorf("update content template placement: %w", err)
	}
	if err := config.SetValueInFileDirect(temporaryPath, modeKey, filenameMode); err != nil {
		return fmt.Errorf("update content template filename mode: %w", err)
	}
	if exists {
		if err := os.Rename(temporaryPath, target); err != nil {
			return fmt.Errorf("save site settings: %w", err)
		}
		return nil
	}
	if err := os.Link(temporaryPath, target); err != nil {
		return fmt.Errorf("create site settings: %w", err)
	}
	return nil
}

func (w *Workspace) ensureDirectory(relative string) (string, error) {
	local := filepath.Clean(filepath.FromSlash(relative))
	if local == "." {
		return w.root, nil
	}
	if !filepath.IsLocal(local) {
		return "", fmt.Errorf("workspace directory must be local")
	}
	current := w.root
	for _, segment := range strings.Split(local, string(filepath.Separator)) {
		if segment == "" || segment == "." {
			continue
		}
		current = filepath.Join(current, segment)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			if err := os.Mkdir(current, 0o755); err != nil && !os.IsExist(err) {
				return "", fmt.Errorf("create workspace directory: %w", err)
			}
			info, err = os.Lstat(current)
		}
		if err != nil {
			return "", fmt.Errorf("inspect workspace directory: %w", err)
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
			return "", fmt.Errorf("workspace directory contains a symlink or non-directory")
		}
	}
	resolved, err := filepath.EvalSymlinks(current)
	if err != nil {
		return "", err
	}
	relativeToRoot, err := filepath.Rel(w.root, resolved)
	if err != nil || relativeToRoot == ".." || strings.HasPrefix(relativeToRoot, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("workspace directory escapes the workspace")
	}
	return resolved, nil
}

func requireNewFile(path string) error {
	if _, err := os.Lstat(path); err == nil {
		return fmt.Errorf("workspace file already exists")
	} else if !os.IsNotExist(err) {
		return err
	}
	return nil
}

func validUploadName(name string) bool {
	if name == "" || len(name) > 100 {
		return false
	}
	for _, r := range name {
		if !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') {
			return false
		}
	}
	return name[0] != '-' && name[len(name)-1] != '-'
}

// ImageExtension validates an uploaded image's declared type and returns its safe extension.
func ImageExtension(contentType string, data []byte) (string, error) {
	detected := http.DetectContentType(data[:min(len(data), 512)])
	declared := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if declared != detected {
		return "", fmt.Errorf("image content type does not match its data")
	}
	switch {
	case detected == "image/png" && bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")):
		return "png", nil
	case detected == "image/jpeg" && bytes.HasPrefix(data, []byte{0xff, 0xd8, 0xff}):
		return "jpg", nil
	case detected == "image/gif" && (bytes.HasPrefix(data, []byte("GIF87a")) || bytes.HasPrefix(data, []byte("GIF89a"))):
		return "gif", nil
	case detected == "image/webp" && len(data) >= 12 && string(data[8:12]) == "WEBP":
		return "webp", nil
	default:
		return "", fmt.Errorf("unsupported image format; use PNG, JPEG, GIF, or WebP")
	}
}

func sourceWithImage(source, imageURL string) string {
	quoted := fmt.Sprintf("image: %q", imageURL)
	match := sourceFrontmatterPattern.FindStringSubmatchIndex(source)
	if match == nil {
		return "---\n" + quoted + "\n---\n\n" + source
	}
	frontmatterStart, frontmatterEnd := match[4], match[5]
	header := source[frontmatterStart:frontmatterEnd]
	lines := strings.Split(header, "\n")
	found := false
	for i, line := range lines {
		lineEnding := ""
		if strings.HasSuffix(line, "\r") {
			lineEnding = "\r"
			line = strings.TrimSuffix(line, "\r")
		}
		if strings.HasPrefix(line, "image:") {
			lines[i] = quoted + lineEnding
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, quoted)
	}
	return source[:frontmatterStart] + strings.Join(lines, "\n") + source[frontmatterEnd:]
}

var sourceFrontmatterPattern = regexp.MustCompile(`(?s)^(---\r?\n)(.*?)(\r?\n---(?:\r?\n|$))`)

func writeTempBytes(directory, pattern string, data []byte) (string, error) {
	tmp, err := os.CreateTemp(directory, pattern)
	if err != nil {
		return "", fmt.Errorf("create workspace temporary file: %w", err)
	}
	tmpPath := tmp.Name()
	if chmodErr := tmp.Chmod(0o644); chmodErr == nil {
		_, err = tmp.Write(data)
	} else {
		err = chmodErr
	}
	if err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		_ = os.Remove(tmpPath)
		return "", fmt.Errorf("write workspace temporary file: %w", err)
	}
	return tmpPath, nil
}

func isMarkdownSource(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown", ".mdown", ".mkd":
		return true
	default:
		return false
	}
}

func writeNewSource(parent, target, source string) error {
	tmp, err := os.CreateTemp(parent, ".markata-editor-*")
	if err != nil {
		return fmt.Errorf("create workspace temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err = tmp.Chmod(0o644); err == nil {
		_, err = tmp.WriteString(source)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("write workspace temp file: %w", err)
	}
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("workspace file already exists")
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Link(tmpName, target); err != nil {
		return fmt.Errorf("create workspace source: %w", err)
	}
	if err := os.Remove(tmpName); err != nil {
		return fmt.Errorf("remove workspace temp file: %w", err)
	}
	return nil
}

func (w *Workspace) writeSource(path, source string) error {
	resolved, err := w.resolve(path)
	if err != nil {
		return err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(resolved), ".markata-editor-*")
	if err != nil {
		return fmt.Errorf("create workspace temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err = tmp.Chmod(info.Mode().Perm()); err == nil {
		_, err = tmp.WriteString(source)
	}
	if err == nil {
		err = tmp.Sync()
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("write workspace temp file: %w", err)
	}
	if err := os.Rename(tmpName, resolved); err != nil {
		return fmt.Errorf("replace workspace source: %w", err)
	}
	return nil
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
