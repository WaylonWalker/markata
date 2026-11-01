package lsp

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/config"
	"github.com/WaylonWalker/markata/pkg/imageindex"
)

func loadManagedImages(root string) []imageindex.Image {
	outputDir := filepath.Join(root, "output")
	for _, name := range []string{"markata-go.toml", "markata-go.yaml", "markata-go.yml", "markata-go.json"} {
		cfg, err := config.LoadSingleConfig(filepath.Join(root, name))
		if err != nil {
			continue
		}
		if cfg.OutputDir != "" {
			outputDir = cfg.OutputDir
			if !filepath.IsAbs(outputDir) {
				outputDir = filepath.Join(root, outputDir)
			}
		}
		break
	}
	data, err := os.ReadFile(filepath.Join(outputDir, "images.json"))
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
