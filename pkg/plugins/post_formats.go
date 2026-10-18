package plugins

import (
	"os"
	"path/filepath"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"gopkg.in/yaml.v3"
)

func resolvePostFormats(post *models.Post, config *lifecycle.Config) models.PostFormatsConfig {
	resolved := getPostFormatsConfig(config)
	if post == nil || post.Extra == nil {
		return resolved
	}

	overrideValue, ok := post.Extra["post_formats"]
	if !ok {
		return resolved
	}

	override, ok := parsePostFormatsOverride(overrideValue)
	if !ok {
		return resolved
	}

	if override.HTML != nil {
		resolved.HTML = override.HTML
	}
	if hasPostFormatKey(overrideValue, "markdown") {
		resolved.Markdown = override.Markdown
	}
	if hasPostFormatKey(overrideValue, "text") {
		resolved.Text = override.Text
	}
	if hasPostFormatKey(overrideValue, "ansi") {
		resolved.ANSI = override.ANSI
	}
	if hasPostFormatKey(overrideValue, "og") {
		resolved.OG = override.OG
	}

	return resolved
}

func parsePostFormatsOverride(value interface{}) (models.PostFormatsConfig, bool) {
	if value == nil {
		return models.PostFormatsConfig{}, false
	}

	switch v := value.(type) {
	case models.PostFormatsConfig:
		return v, true
	case map[string]interface{}:
		return postFormatsConfigFromMap(v), true
	default:
		return models.PostFormatsConfig{}, false
	}
}

func postFormatsConfigFromMap(raw map[string]interface{}) models.PostFormatsConfig {
	encoded, err := yaml.Marshal(raw)
	if err != nil {
		return models.PostFormatsConfig{}
	}

	var parsed struct {
		HTML     *bool `yaml:"html"`
		Markdown *bool `yaml:"markdown"`
		Text     *bool `yaml:"text"`
		ANSI     *bool `yaml:"ansi"`
		OG       *bool `yaml:"og"`
	}
	if err := yaml.Unmarshal(encoded, &parsed); err != nil {
		return models.PostFormatsConfig{}
	}

	return models.PostFormatsConfig{
		HTML:     parsed.HTML,
		Markdown: parsed.Markdown != nil && *parsed.Markdown,
		Text:     parsed.Text != nil && *parsed.Text,
		ANSI:     parsed.ANSI != nil && *parsed.ANSI,
		OG:       parsed.OG != nil && *parsed.OG,
	}
}

func hasPostFormatKey(value interface{}, key string) bool {
	switch v := value.(type) {
	case map[string]interface{}:
		_, ok := v[key]
		return ok
	case models.PostFormatsConfig:
		switch key {
		case "html":
			return v.HTML != nil
		case "markdown", "text", "ansi", "og":
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func applyPostFormatsToConfig(base *models.Config, postFormats models.PostFormatsConfig) *models.Config {
	if base == nil {
		return nil
	}

	resolvedConfig := *base
	resolvedConfig.PostFormats = postFormats
	return &resolvedConfig
}

func removeDisabledPostOutputs(outputDir, slug string, postFormats models.PostFormatsConfig) {
	postDir := filepath.Join(outputDir, slug)
	if !postFormats.IsHTMLEnabled() {
		_ = os.Remove(filepath.Join(postDir, "index.html"))
	}
	if !postFormats.Markdown {
		_ = os.Remove(filepath.Join(outputDir, slug+".md"))
		_ = os.RemoveAll(filepath.Join(postDir, "index.md"))
	}
	if !postFormats.Text {
		_ = os.Remove(filepath.Join(outputDir, slug+".txt"))
		_ = os.RemoveAll(filepath.Join(postDir, "index.txt"))
	}
	if !postFormats.ANSI {
		_ = os.Remove(filepath.Join(outputDir, slug+".ansi"))
		_ = os.RemoveAll(filepath.Join(postDir, "index.ansi"))
	}
	if !postFormats.OG {
		_ = os.RemoveAll(filepath.Join(postDir, "og"))
	}
}

// removePrivatePostOutputs removes alternate formats that may have been
// generated while a post was public. The encrypted HTML page is retained.
func removePrivatePostOutputs(outputDir, slug string) error {
	for _, extension := range []string{"md", "txt", "ansi"} {
		path, err := safeOutputPath(outputDir, slug+"."+extension)
		if err != nil {
			return err
		}
		_ = os.Remove(path)
		redirectPath, err := safeOutputPath(outputDir, filepath.Join(slug, "index."+extension))
		if err != nil {
			return err
		}
		_ = os.RemoveAll(redirectPath)
	}
	ogPath, err := safeOutputPath(outputDir, filepath.Join(slug, "og"))
	if err != nil {
		return err
	}
	_ = os.RemoveAll(ogPath)
	return nil
}
