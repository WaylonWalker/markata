package workbench

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

const feedSnapshotPrefix = "/workspace/feed-snapshot/"

// feedSnapshotFile only opens generated files under the selected workspace
// output directory. It will not follow a symlink out of that directory.
func feedSnapshotFile(root string, cfg *models.Config, requested string) (string, error) {
	if cfg == nil {
		return "", os.ErrNotExist
	}
	output := cfg.OutputDir
	if output == "" {
		output = "output"
	}
	if !filepath.IsAbs(output) {
		output = filepath.Join(root, output)
	}
	siteRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}
	outputRoot, err := filepath.EvalSymlinks(output)
	if err != nil {
		return "", err
	}
	underSite, err := filepath.Rel(siteRoot, outputRoot)
	if err != nil || underSite == "." || underSite == ".." || strings.HasPrefix(underSite, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("feed output must stay inside the workspace")
	}
	requested = strings.TrimPrefix(requested, "/")
	if requested == "" || strings.HasSuffix(requested, "/") {
		requested += "index.html"
	}
	if !filepath.IsLocal(filepath.FromSlash(requested)) {
		return "", os.ErrNotExist
	}
	for _, segment := range strings.Split(requested, "/") {
		if segment == "" || strings.HasPrefix(segment, ".") {
			return "", os.ErrNotExist
		}
	}
	candidate := filepath.Join(outputRoot, filepath.FromSlash(requested))
	resolved, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(outputRoot, resolved)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", os.ErrNotExist
	}
	info, err := os.Stat(resolved)
	if err != nil || !info.Mode().IsRegular() {
		return "", os.ErrNotExist
	}
	return resolved, nil
}

// This is the real last-built feed HTML, CSS and media: never a fabricated
// approximation of templates. The browser uses an opaque-origin sandboxed iframe.
func (h *Handler) serveFeedSnapshot(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, feedSnapshotPrefix)
	filename, err := feedSnapshotFile(h.root, h.currentSiteConfig(), name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if strings.HasSuffix(strings.ToLower(filename), ".html") || strings.HasSuffix(strings.ToLower(filename), ".css") {
		file, err := os.Open(filename)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer file.Close()
		source, err := io.ReadAll(io.LimitReader(file, 5<<20+1))
		if err != nil || len(source) > 5<<20 {
			http.Error(w, "rendered site asset is too large", http.StatusRequestEntityTooLarge)
			return
		}
		// Root-absolute asset references must stay within the isolated snapshot
		// mount, rather than resolve to editor routes on the same origin.
		rewritten := strings.NewReplacer(
			`href="/`, `href="`+feedSnapshotPrefix,
			`src="/`, `src="`+feedSnapshotPrefix,
			`poster="/`, `poster="`+feedSnapshotPrefix,
			`srcset="/`, `srcset="`+feedSnapshotPrefix,
			`url("/`, `url("`+feedSnapshotPrefix,
			`url('/`, `url('`+feedSnapshotPrefix,
			`url(/`, `url(`+feedSnapshotPrefix,
		).Replace(string(source))
		serve(w, r, filename, []byte(rewritten))
		return
	}
	http.ServeFile(w, r, filename)
}
