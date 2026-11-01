package lsp

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/imageindex"
)

func loadManagedImages(root string) []imageindex.Image {
	data, err := os.ReadFile(filepath.Join(root, "output", "images.json"))
	if err != nil {
		return nil
	}
	index, err := imageindex.Parse(data)
	if err != nil {
		return nil
	}
	images := make([]imageindex.Image, 0, len(index.Images))
	for i := range index.Images {
		image := &index.Images[i]
		parsed, err := url.Parse(image.Src)
		if err != nil || parsed.IsAbs() || parsed.Host != "" || strings.Contains(image.Src, "${") || !strings.HasPrefix(image.MIMEType, "image/") {
			continue
		}
		images = append(images, *image)
	}
	return images
}
