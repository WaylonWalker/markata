package visualizations

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Publish atomically writes a content-addressed SVG under the site output
// directory. The destination path is derived from the bytes, not user input.
// Existing identical assets are reused; mismatches are treated as corruption.
func Publish(outputDir string, asset Asset) (string, error) {
	canonical, err := NewSVGAsset(asset.SVG)
	if err != nil {
		return "", err
	}
	if asset.Path != canonical.Path || asset.SHA256 != canonical.SHA256 ||
		asset.Width != canonical.Width || asset.Height != canonical.Height {
		return "", errors.New("SVG asset identity/dimensions do not match bytes")
	}
	path := filepath.Join(outputDir, filepath.FromSlash(canonical.Path))
	if existing, err := os.ReadFile(path); err == nil {
		if bytes.Equal(existing, asset.SVG) {
			return canonical.Path, nil
		}
		return "", fmt.Errorf("SVG path %s contains different content", canonical.Path)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".svg-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(asset.SVG); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Chmod(0o644); err != nil {
		file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return "", err
	}
	return canonical.Path, nil
}

// FigureHTML emits a fully static figure with intrinsic dimensions so the
// browser reserves space before the external SVG is fetched. No rendering JS
// is needed. External site CSS can make the image responsive (max-width:100%).
func FigureHTML(asset Asset, alt, caption string, aboveFold bool) (string, error) {
	return FigureHTMLWithBasePath(asset, alt, caption, aboveFold, "")
}

// FigureHTMLWithBasePath emits the same figure with an optional site URL path
// prefix, such as "/blog" for a site hosted below the domain root.
func FigureHTMLWithBasePath(asset Asset, alt, caption string, aboveFold bool, basePath string) (string, error) {
	if strings.TrimSpace(alt) == "" {
		return "", errors.New("SVG figure requires descriptive alternative text")
	}
	basePath = strings.TrimRight(strings.TrimSpace(basePath), "/")
	if basePath != "" {
		decoded, err := url.PathUnescape(basePath)
		if err != nil || !strings.HasPrefix(basePath, "/") || strings.HasPrefix(basePath, "//") ||
			strings.ContainsAny(decoded, "?#\\") || path.Clean(decoded) != decoded {
			return "", errors.New("SVG figure base path must be a normalized local URL path")
		}
	}
	canonical, err := NewSVGAsset(asset.SVG)
	if err != nil {
		return "", err
	}
	if asset.Path != canonical.Path || asset.SHA256 != canonical.SHA256 || asset.Width != canonical.Width || asset.Height != canonical.Height {
		return "", errors.New("SVG figure asset identity/dimensions do not match bytes")
	}
	loading := "lazy"
	if aboveFold {
		loading = "eager"
	}
	markup := `<figure class="markata-visualization"><img src="` + html.EscapeString(basePath+"/"+asset.Path) +
		`" width="` + strconv.Itoa(asset.Width) + `" height="` + strconv.Itoa(asset.Height) +
		`" loading="` + loading + `" decoding="async" alt="` + html.EscapeString(alt) + `">`
	if caption != "" {
		markup += "<figcaption>" + html.EscapeString(caption) + "</figcaption>"
	}
	return markup + "</figure>", nil
}
