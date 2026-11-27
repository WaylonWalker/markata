package workbench

import (
	"errors"
	"fmt"
	"html"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/templates"
)

const templatePreviewWarning = `<aside class="workbench-frontmatter-warning" role="status">Site template preview unavailable; showing the basic draft preview.</aside>`

// renderSiteFrontmatterPreview keeps generated template markup separate from
// source-backed Markdown blocks. The existing opaque preview iframe owns
// appearance and incremental updates.
func (h *Handler) renderSiteFrontmatterPreview(source, markdown, sourcePath string) string {
	metadata, _, err := plugins.ParseFrontmatter(source)
	if err != nil {
		return renderFrontmatterPreview(source, markdown)
	}
	// The client path is only a template-selection hint, never a file read.
	if sourcePath == "" || !fs.ValidPath(sourcePath) || !filepath.IsLocal(filepath.FromSlash(sourcePath)) || strings.Contains(sourcePath, "\\") {
		sourcePath = "draft.md"
	}
	post, resolved, err := plugins.ResolveDraftTemplate(metadata, sourcePath, h.siteConfig)
	if err != nil {
		return templatePreviewWarning + renderFrontmatterMetadataPreview(metadata, markdown)
	}
	kind := frontmatterTemplateName(metadata)
	if !isPhotoTemplate(kind) && !isLinkTemplate(kind) {
		kind = normalizedPreviewTemplate(resolved)
	}
	fallback := renderFrontmatterMetadataForTemplate(metadata, markdown, kind)
	if h.siteConfig == nil {
		return fallback
	}
	if !fs.ValidPath(resolved) || strings.Contains(resolved, "\\") || !filepath.IsLocal(filepath.FromSlash(resolved)) {
		return templatePreviewWarning + fallback
	}
	return h.renderCompanionPreview(post, resolved, markdown, fallback)
}

func (h *Handler) openPreviewTemplateRoot() (*os.Root, error) {
	// All custom files, including includes, must stay within the selected
	// workspace's configured template directory. Do not use the build loader,
	// which intentionally also permits absolute paths and external themes.
	dir := h.siteConfig.TemplatesDir
	if dir == "" {
		dir = "templates"
	}
	if filepath.IsAbs(dir) {
		relative, err := filepath.Rel(h.root, dir)
		if err != nil {
			return nil, err
		}
		dir = relative
	}
	if !filepath.IsLocal(dir) {
		return nil, fmt.Errorf("unsafe editor template directory")
	}
	root, err := os.OpenRoot(h.root)
	if err != nil {
		return nil, err
	}
	defer root.Close()
	return root.OpenRoot(dir)
}

func (h *Handler) renderCompanionPreview(post *models.Post, resolved, markdown, fallback string) string {
	templateRoot, err := h.openPreviewTemplateRoot()
	if errors.Is(err, fs.ErrNotExist) {
		return fallback
	}
	if err != nil {
		return templatePreviewWarning + fallback
	}
	defer templateRoot.Close()

	// Companion fragments replace the basic metadata prefix. Markdown remains
	// rendered exactly once by the standalone renderer, preserving scroll sync.
	fragment := path.Join("editor", resolved)
	if _, err = templateRoot.Stat(filepath.FromSlash(fragment)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return templatePreviewWarning + fallback
		}
		if _, err = templateRoot.Stat(filepath.FromSlash(resolved)); err == nil {
			return templatePreviewWarning + fallback
		}
		return fallback
	}
	post.Content = markdown
	post.TitleHTML = html.EscapeString(post.PlainTitle())
	markup, err := templates.RenderEditorFragment(templateRoot, fragment, post, h.siteConfig)
	if err != nil {
		return templatePreviewWarning + fallback
	}
	return `<div class="workbench-template-preview">` + markup + `</div>`
}
