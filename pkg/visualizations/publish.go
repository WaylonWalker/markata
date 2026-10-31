package visualizations

import (
 "bytes"
 "errors"
 "fmt"
 "html"
 "os"
 "path/filepath"
 "strings"
)

// Publish atomically writes a fingerprinted SVG under the output directory.
// The site-relative path is derived from content, never from authored input.
func Publish(outputDir string, asset Asset) (string, error) {
 canonical, err := NewSVGAsset(asset.SVG)
 if err != nil { return "", err }
 if asset.Path != canonical.Path || asset.SHA256 != canonical.SHA256 {
  return "", errors.New("SVG asset identity does not match bytes")
 }
 path := filepath.Join(outputDir, filepath.FromSlash(canonical.Path))
 if existing, err := os.ReadFile(path); err == nil {
  if bytes.Equal(existing, asset.SVG) { return canonical.Path, nil }
  return "", fmt.Errorf("SVG path %s contains different content", canonical.Path)
 } else if !os.IsNotExist(err) { return "", err }
 if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return "", err }
 file, err := os.CreateTemp(filepath.Dir(path), ".svg-*")
 if err != nil { return "", err }
 defer os.Remove(file.Name())
 if _, err := file.Write(asset.SVG); err != nil {
  file.Close()
  return "", err
 }
 if err := file.Chmod(0o644); err != nil {
  file.Close()
  return "", err
 }
 if err := file.Close(); err != nil { return "", err }
 if err := os.Rename(file.Name(), path); err != nil { return "", err }
 return canonical.Path, nil
}

// FigureHTML emits static HTML with intrinsic dimensions to prevent layout
// shift. The site theme can set image max-width:100% for responsiveness.
// aboveFold should be true for a diagram critical to initial page rendering.
func FigureHTML(asset Asset, alt, caption string, aboveFold bool) (string, error) {
 if strings.TrimSpace(alt) == "" {
  return "", errors.New("SVG figure requires descriptive alternative text")
 }
 canonical, err := NewSVGAsset(asset.SVG)
 if err != nil { return "", err }
 if asset.Path != canonical.Path || asset.SHA256 != canonical.SHA256 ||
  asset.Width != canonical.Width || asset.Height != canonical.Height {
   return "", errors.New("SVG figure asset identity/dimensions do not match bytes")
 }
 loading := "lazy"
 if aboveFold { loading = "eager" }
 markup := fmt.Sprintf(`<figure class="markata-visualization"><img src="/%s" width="%d" height="%d" loading="%s" decoding="async" alt="%s">`,
  html.EscapeString(asset.Path), asset.Width, asset.Height, loading, html.EscapeString(alt))
 if caption != "" {
  markup += "<figcaption>" + html.EscapeString(caption) + "</figcaption>"
 }
 return markup + "</figure>", nil
}
