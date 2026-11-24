package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestLoadFromStringProjectsSupportedTypedFields(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		data   string
	}{
		{
			name:   "toml",
			format: FormatTOML,
			data: `[markata-go]
[markata-go.head]
text = "<meta name=\"fixture\">"
[[markata-go.head.meta]]
name = "robots"
content = "noindex"

[markata-go.theme_calendar]
enabled = false
[[markata-go.theme_calendar.rules]]
name = "Winter"
start_date = "12-01"
end_date = "02-28"
palette = "winter-frost"

[markata-go.error_pages]
enable_404 = false
custom_404_template = "not-found.html"
max_suggestions = 9

[markata-go.resource_hints]
enabled = false
auto_detect = false
exclude_domains = ["excluded.example"]
[[markata-go.resource_hints.domains]]
domain = "cdn.example"
hint_types = ["preconnect"]

[markata-go.markdown.highlight]
enabled = false
theme = "monokai"
line_numbers = true

[markata-go.template_presets.article]
html = "article.html"
txt = "article.txt"
[markata-go.default_templates]
html = "default.html"
`,
		},
		{
			name:   "yaml",
			format: FormatYAML,
			data: `markata-go:
  head:
    text: '<meta name="fixture">'
    meta:
      - name: robots
        content: noindex
  theme_calendar:
    enabled: false
    rules:
      - name: Winter
        start_date: 12-01
        end_date: 02-28
        palette: winter-frost
  error_pages:
    enable_404: false
    custom_404_template: not-found.html
    max_suggestions: 9
  resource_hints:
    enabled: false
    auto_detect: false
    exclude_domains: [excluded.example]
    domains:
      - domain: cdn.example
        hint_types: [preconnect]
  markdown:
    highlight:
      enabled: false
      theme: monokai
      line_numbers: true
  template_presets:
    article:
      html: article.html
      txt: article.txt
  default_templates:
    html: default.html
`,
		},
		{
			name:   "json",
			format: FormatJSON,
			data:   `{"markata-go":{"head":{"text":"<meta name=\"fixture\">","meta":[{"name":"robots","content":"noindex"}]},"theme_calendar":{"enabled":false,"rules":[{"name":"Winter","start_date":"12-01","end_date":"02-28","palette":"winter-frost"}]},"error_pages":{"enable_404":false,"custom_404_template":"not-found.html","max_suggestions":9},"resource_hints":{"enabled":false,"auto_detect":false,"exclude_domains":["excluded.example"],"domains":[{"domain":"cdn.example","hint_types":["preconnect"]}]},"markdown":{"highlight":{"enabled":false,"theme":"monokai","line_numbers":true}},"template_presets":{"article":{"html":"article.html","txt":"article.txt"}},"default_templates":{"html":"default.html"}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadFromString(tt.data, tt.format)
			if err != nil {
				t.Fatalf("LoadFromString() error = %v", err)
			}

			if config.Head.Text != "<meta name=\"fixture\">" || len(config.Head.Meta) != 1 {
				t.Fatalf("Head = %+v, want typed head config", config.Head)
			}
			if config.TemplatePresets["article"].Text != "article.txt" {
				t.Errorf("TemplatePresets = %+v, want article txt template", config.TemplatePresets)
			}
			if config.DefaultTemplates["html"] != "default.html" {
				t.Errorf("DefaultTemplates = %+v, want html template", config.DefaultTemplates)
			}
			if config.ThemeCalendar.IsEnabled() || len(config.ThemeCalendar.Rules) != 1 {
				t.Errorf("ThemeCalendar = %+v, want disabled calendar with one rule", config.ThemeCalendar)
			}
			if config.ThemeCalendar.Rules[0].Palette != "winter-frost" {
				t.Errorf("ThemeCalendar rule = %+v, want winter-frost", config.ThemeCalendar.Rules[0])
			}
			if config.ErrorPages.Is404Enabled() || config.ErrorPages.Custom404Template != "not-found.html" || config.ErrorPages.MaxSuggestions != 9 {
				t.Errorf("ErrorPages = %+v, want explicit values including false", config.ErrorPages)
			}
			if config.ResourceHints.IsEnabled() || config.ResourceHints.IsAutoDetectEnabled() || len(config.ResourceHints.Domains) != 1 {
				t.Errorf("ResourceHints = %+v, want explicit disabled values", config.ResourceHints)
			}
			if config.MarkdownConfig.Highlight.IsEnabled() || config.MarkdownConfig.Highlight.Theme != "monokai" || !config.MarkdownConfig.Highlight.LineNumbers {
				t.Errorf("Markdown highlight = %+v, want explicit values including false", config.MarkdownConfig.Highlight)
			}
		})
	}
}

func TestLoadFromStringNormalizesHeadTextBlocks(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		data   string
	}{
		{
			name:   "toml",
			format: FormatTOML,
			data: `[markata-go.head]
[[markata-go.head.link]]
rel = "webmention"
href = "https://example.test/webmention"
[[markata-go.head.text]]
value = "<style>one</style>"
[[markata-go.head.text]]
value = "<script>two</script>"
`,
		},
		{
			name:   "yaml",
			format: FormatYAML,
			data: `markata-go:
  head:
    link:
      - rel: webmention
        href: https://example.test/webmention
    text:
      - value: <style>one</style>
      - value: <script>two</script>
`,
		},
		{
			name:   "json",
			format: FormatJSON,
			data:   `{"markata-go":{"head":{"link":[{"rel":"webmention","href":"https://example.test/webmention"}],"text":[{"value":"<style>one</style>"},{"value":"<script>two</script>"}]}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadFromString(tt.data, tt.format)
			if err != nil {
				t.Fatalf("LoadFromString() error = %v", err)
			}
			if config.Head.Text != "<style>one</style>\n<script>two</script>" {
				t.Errorf("Head.Text = %q, want normalized text blocks", config.Head.Text)
			}
			if len(config.Head.Link) != 1 || config.Head.Link[0].Rel != "webmention" {
				t.Errorf("Head.Link = %+v, want legacy head link preserved", config.Head.Link)
			}
		})
	}
}

func TestLoadFromStringRejectsInvalidHeadText(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		data   string
	}{
		{
			name:   "toml",
			format: FormatTOML,
			data:   "[markata-go.head]\ntext = 42\n",
		},
		{
			name:   "yaml",
			format: FormatYAML,
			data:   "markata-go:\n  head:\n    text: 42\n",
		},
		{
			name:   "json",
			format: FormatJSON,
			data:   `{"markata-go":{"head":{"text":42}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := LoadFromString(tt.data, tt.format); err == nil {
				t.Fatal("LoadFromString() error = nil, want invalid head.text error")
			}
		})
	}
}

func TestLoadProjectsAffectedTypedFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "markata-go.toml")
	data := `[markata-go]
[markata-go.head]
text = "loaded head"
[markata-go.theme_calendar]
enabled = false
[markata-go.error_pages]
enable_404 = false
[markata-go.resource_hints]
enabled = false
[markata-go.markdown.highlight]
enabled = false
theme = "monokai"
[markata-go.template_presets.article]
html = "article.html"
[markata-go.default_templates]
html = "default.html"
`
	//nolint:gosec // Test files use temporary paths and fixture permissions.
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}

	config, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config.Head.Text != "loaded head" || config.TemplatePresets["article"].HTML != "article.html" || config.DefaultTemplates["html"] != "default.html" {
		t.Errorf("loaded template/head projection = head=%+v presets=%+v defaults=%+v", config.Head, config.TemplatePresets, config.DefaultTemplates)
	}
	if config.ThemeCalendar.IsEnabled() || config.ErrorPages.Is404Enabled() || config.ResourceHints.IsEnabled() || config.MarkdownConfig.Highlight.IsEnabled() {
		t.Errorf("explicit false values were lost: calendar=%+v errors=%+v hints=%+v highlight=%+v", config.ThemeCalendar, config.ErrorPages, config.ResourceHints, config.MarkdownConfig.Highlight)
	}
}

func TestParseFunctionsProjectAffectedTypedFields(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte) (*models.Config, error)
		data  string
	}{
		{
			name:  "toml",
			parse: ParseTOML,
			data: `[markata-go]
[markata-go.head]
text = "direct parser"
[markata-go.theme_calendar]
enabled = false
[markata-go.error_pages]
enable_404 = false
[markata-go.resource_hints]
enabled = false
[markata-go.markdown.highlight]
enabled = false
theme = "monokai"
[markata-go.template_presets.article]
html = "article.html"
[markata-go.default_templates]
html = "default.html"
`,
		},
		{
			name:  "yaml",
			parse: ParseYAML,
			data: `markata-go:
  head: {text: direct parser}
  theme_calendar: {enabled: false}
  error_pages: {enable_404: false}
  resource_hints: {enabled: false}
  markdown:
    highlight: {enabled: false, theme: monokai}
  template_presets:
    article: {html: article.html}
  default_templates: {html: default.html}
`,
		},
		{
			name:  "json",
			parse: ParseJSON,
			data:  `{"markata-go":{"head":{"text":"direct parser"},"theme_calendar":{"enabled":false},"error_pages":{"enable_404":false},"resource_hints":{"enabled":false},"markdown":{"highlight":{"enabled":false,"theme":"monokai"}},"template_presets":{"article":{"html":"article.html"}},"default_templates":{"html":"default.html"}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tt.parse([]byte(tt.data))
			if err != nil {
				t.Fatalf("parse() error = %v", err)
			}
			if config.Head.Text != "direct parser" || config.TemplatePresets["article"].HTML != "article.html" || config.DefaultTemplates["html"] != "default.html" {
				t.Errorf("head/templates were not projected: head=%+v presets=%+v defaults=%+v", config.Head, config.TemplatePresets, config.DefaultTemplates)
			}
			if config.ThemeCalendar.IsEnabled() || config.ErrorPages.Is404Enabled() || config.ResourceHints.IsEnabled() || config.MarkdownConfig.Highlight.IsEnabled() || config.MarkdownConfig.Highlight.Theme != "monokai" {
				t.Errorf("explicit false values were lost: calendar=%+v errors=%+v hints=%+v highlight=%+v", config.ThemeCalendar, config.ErrorPages, config.ResourceHints, config.MarkdownConfig.Highlight)
			}
		})
	}
}

func TestLoadFromStringPreservesExplicitZeroAndEmptyValues(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		data   string
	}{
		{
			name:   "toml",
			format: FormatTOML,
			data: `[markata-go]
[markata-go.error_pages]
custom_404_template = ""
max_suggestions = 0
[markata-go.template_presets.empty]
html = ""
[markata-go.default_templates]
html = ""
`,
		},
		{
			name:   "yaml",
			format: FormatYAML,
			data: `markata-go:
  error_pages:
    custom_404_template: ""
    max_suggestions: 0
  template_presets:
    empty:
      html: ""
  default_templates:
    html: ""
`,
		},
		{
			name:   "json",
			format: FormatJSON,
			data:   `{"markata-go":{"error_pages":{"custom_404_template":"","max_suggestions":0},"template_presets":{"empty":{"html":""}},"default_templates":{"html":""}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadFromString(tt.data, tt.format)
			if err != nil {
				t.Fatalf("LoadFromString() error = %v", err)
			}
			if config.ErrorPages.Custom404Template != "" || config.ErrorPages.MaxSuggestions != 0 {
				t.Errorf("ErrorPages = %+v, want explicit empty template and zero suggestions", config.ErrorPages)
			}
			preset, ok := config.TemplatePresets["empty"]
			if !ok || preset.HTML != "" {
				t.Errorf("TemplatePresets[empty] = %+v, present=%v, want explicit empty html template", preset, ok)
			}
			defaultTemplate, ok := config.DefaultTemplates["html"]
			if !ok || defaultTemplate != "" {
				t.Errorf("DefaultTemplates[html] = %q, present=%v, want explicit empty template", defaultTemplate, ok)
			}
		})
	}
}

func TestDefaultConfig_AffectedDefaults(t *testing.T) {
	config := DefaultConfig()

	if !config.MarkdownConfig.Highlight.IsEnabled() {
		t.Error("MarkdownConfig.Highlight should be enabled by default")
	}
	if config.ThemeCalendar.IsEnabled() {
		t.Error("ThemeCalendar should be disabled by default")
	}
	if !config.ErrorPages.Is404Enabled() {
		t.Error("ErrorPages 404 should be enabled by default")
	}
	if !config.ResourceHints.IsEnabled() || !config.ResourceHints.IsAutoDetectEnabled() {
		t.Error("ResourceHints should be enabled with auto-detection by default")
	}
}

func TestLoadWithMergeProjectsAffectedFields(t *testing.T) {
	dir := t.TempDir()
	basePath := filepath.Join(dir, "base.toml")
	overridePath := filepath.Join(dir, "override.yaml")

	base := `[markata-go]
[markata-go.head]
text = "base head"
[markata-go.markdown.highlight]
theme = "base-theme"
[markata-go.theme_calendar]
enabled = true
[[markata-go.theme_calendar.rules]]
name = "Base"
start_date = "01-01"
end_date = "01-31"
[markata-go.error_pages]
custom_404_template = "base-404.html"
[markata-go.resource_hints]
exclude_domains = ["base.example"]
[markata-go.template_presets.article]
html = "base.html"
[markata-go.default_templates]
html = "base-default.html"
`
	override := `markata-go:
  markdown:
    highlight:
      enabled: false
      line_numbers: true
  theme_calendar:
    enabled: false
  error_pages:
    enable_404: false
  resource_hints:
    enabled: false
  template_presets:
    article:
      txt: override.txt
  default_templates:
    txt: override.txt
`
	//nolint:gosec // Test files use temporary paths and fixture permissions.
	if err := os.WriteFile(basePath, []byte(base), 0o644); err != nil {
		t.Fatal(err)
	}
	//nolint:gosec // Test files use temporary paths and fixture permissions.
	if err := os.WriteFile(overridePath, []byte(override), 0o644); err != nil {
		t.Fatal(err)
	}

	config, err := LoadWithMerge(basePath, overridePath)
	if err != nil {
		t.Fatalf("LoadWithMerge() error = %v", err)
	}

	if config.Head.Text != "base head" {
		t.Errorf("Head.Text = %q, want base head", config.Head.Text)
	}
	if config.MarkdownConfig.Highlight.IsEnabled() || config.MarkdownConfig.Highlight.Theme != "base-theme" || !config.MarkdownConfig.Highlight.LineNumbers {
		t.Errorf("Markdown highlight = %+v, want deep-merged override", config.MarkdownConfig.Highlight)
	}
	if config.ThemeCalendar.IsEnabled() || len(config.ThemeCalendar.Rules) != 1 {
		t.Errorf("ThemeCalendar = %+v, want disabled base rule", config.ThemeCalendar)
	}
	if config.ErrorPages.Is404Enabled() || config.ErrorPages.Custom404Template != "base-404.html" {
		t.Errorf("ErrorPages = %+v, want deep-merged override", config.ErrorPages)
	}
	if config.ResourceHints.IsEnabled() || len(config.ResourceHints.ExcludeDomains) != 1 || config.ResourceHints.ExcludeDomains[0] != "base.example" {
		t.Errorf("ResourceHints = %+v, want deep-merged override", config.ResourceHints)
	}
	if preset := config.TemplatePresets["article"]; preset.HTML != "base.html" || preset.Text != "override.txt" {
		t.Errorf("TemplatePresets[article] = %+v, want fields from both files", preset)
	}
	if config.DefaultTemplates["html"] != "base-default.html" || config.DefaultTemplates["txt"] != "override.txt" {
		t.Errorf("DefaultTemplates = %+v, want fields from both files", config.DefaultTemplates)
	}
}

func TestLoadWithMergePreservesExplicitZeroAndEmptyValues(t *testing.T) {
	dir := t.TempDir()
	basePath := filepath.Join(dir, "base.toml")
	overridePath := filepath.Join(dir, "override.yaml")

	base := `[markata-go]
[markata-go.head]
text = "base head"
[[markata-go.head.meta]]
name = "robots"
content = "noindex"
[markata-go.error_pages]
custom_404_template = "base-404.html"
max_suggestions = 7
[markata-go.markdown.highlight]
theme = "base-theme"
line_numbers = true
[markata-go.template_presets.article]
html = "base.html"
[markata-go.default_templates]
html = "base-default.html"
`
	override := `markata-go:
  head:
    text: ""
  error_pages:
    custom_404_template: ""
    max_suggestions: 0
  markdown:
    highlight:
      theme: ""
  template_presets:
    article:
      html: ""
  default_templates:
    html: ""
`
	//nolint:gosec // Test files use temporary paths and fixture permissions.
	if err := os.WriteFile(basePath, []byte(base), 0o644); err != nil {
		t.Fatal(err)
	}
	//nolint:gosec // Test files use temporary paths and fixture permissions.
	if err := os.WriteFile(overridePath, []byte(override), 0o644); err != nil {
		t.Fatal(err)
	}

	config, err := LoadWithMergeOptions(LoadOptions{DisableDotEnv: true, DisableEnvOverrides: true}, basePath, overridePath)
	if err != nil {
		t.Fatalf("LoadWithMergeOptions() error = %v", err)
	}
	if config.Head.Text != "" || len(config.Head.Meta) != 1 || config.Head.Meta[0].Content != "noindex" {
		t.Errorf("Head = %+v, want empty text with base meta preserved", config.Head)
	}
	if config.ErrorPages.Custom404Template != "" || config.ErrorPages.MaxSuggestions != 0 {
		t.Errorf("ErrorPages = %+v, want explicit empty and zero overrides", config.ErrorPages)
	}
	if config.MarkdownConfig.Highlight.Theme != "" || !config.MarkdownConfig.Highlight.LineNumbers {
		t.Errorf("Highlight = %+v, want empty theme with base line numbers preserved", config.MarkdownConfig.Highlight)
	}
	if preset, ok := config.TemplatePresets["article"]; !ok || preset.HTML != "" {
		t.Errorf("TemplatePresets[article] = %+v, present=%v, want explicit empty html override", preset, ok)
	}
	if defaultTemplate, ok := config.DefaultTemplates["html"]; !ok || defaultTemplate != "" {
		t.Errorf("DefaultTemplates[html] = %q, present=%v, want explicit empty override", defaultTemplate, ok)
	}
}
