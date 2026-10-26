package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigKeyTypoSuggestsNestedKey(t *testing.T) {
	_, err := LoadFromString("[markata-go.theme]\npalete = 'ayu-dark'\n", FormatTOML)
	if err == nil || !strings.Contains(err.Error(), "theme.palette") {
		t.Fatalf("error = %v; want theme.palette suggestion", err)
	}
}

func TestGetValueFromFileReadsExistingKey(t *testing.T) {
	path := filepath.Join(t.TempDir(), "markata-go.toml")
	if err := os.WriteFile(path, []byte("[markata-go]\noutput_dir = 'public'\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"output_dir", "markata-go.output_dir"} {
		value, err := GetValueFromFile(path, key)
		if err != nil || value != "public" {
			t.Errorf("GetValueFromFile(%q) = %v, %v; want public", key, value, err)
		}
	}
	_, err := GetValueFromFile(path, "missing_key")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("missing key error = %v; want ErrKeyNotFound", err)
	}
}

func TestConfigKeyTypoSuggestsTopLevelKey(t *testing.T) {
	_, err := LoadFromString("[markata-go]\noutput_dr = 'public'\n", FormatTOML)
	if err == nil || !strings.Contains(err.Error(), "output_dir") {
		t.Fatalf("error = %v; want output_dir suggestion", err)
	}
}

func TestConfigUnknownPluginSectionRemainsValid(t *testing.T) {
	for _, section := range []string{"stats", "webmentions"} {
		_, err := LoadFromString("[markata-go."+section+"]\ncount = 3\n", FormatTOML)
		if err != nil {
			t.Fatalf("custom plugin section %q: %v", section, err)
		}
	}
}

func TestMisplacedFeedFormatKeyAcrossFormats(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		format Format
	}{
		{
			name: "toml",
			data: "[[markata-go.feeds]]\nslug = 'published'\n[markata-go.feeds.formats]\nsitemap = true\n[[markata-go.feeds]]\nslug = 'draft'\nsitemap = false\n",
			format: FormatTOML,
		},
		{
			name: "yaml",
			data: "markata-go:\n  feeds:\n    - slug: published\n      formats:\n        sitemap: true\n    - slug: draft\n      sitemap: false\n",
			format: FormatYAML,
		},
		{
			name: "json",
			data: `{"markata-go":{"feeds":[{"slug":"published","formats":{"sitemap":true}},{"slug":"draft","sitemap":false}]}}`,
			format: FormatJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := LoadFromString(tt.data, tt.format)
			if err == nil {
				t.Fatal("expected misplaced feed format diagnostic")
			}
			for _, want := range []string{"feeds[1]", "draft", "formats.sitemap"} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error = %v; missing %q", err, want)
				}
			}
		})
	}
}

func TestFeedFormatKeysAndExtensionFields(t *testing.T) {
	for _, key := range []string{"html", "simple_html", "rss", "atom", "json", "markdown", "text", "sitemap"} {
		t.Run(key, func(t *testing.T) {
			_, err := LoadFromString("[[markata-go.feeds]]\nslug = 'draft'\n"+key+" = false\n", FormatTOML)
			if err == nil || !strings.Contains(err.Error(), "formats."+key) {
				t.Fatalf("error = %v; want formats.%s", err, key)
			}
		})
	}
	_, err := LoadFromString("[[markata-go.feeds]]\nslug = 'draft'\nplugin_custom = true\n[markata-go.feeds.formats]\nsitemap = false\n", FormatTOML)
	if err != nil {
		t.Fatalf("valid nested formats and extension field: %v", err)
	}
}

func TestIncludedFeedFormatDiagnosticNamesSource(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "markata-go.toml")
	child := filepath.Join(dir, "feeds.toml")
	if err := os.WriteFile(root, []byte("[markata-go]\ninclude = ['feeds.toml']\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(child, []byte("[[markata-go.feeds]]\nslug = 'draft'\nsitemap = false\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := loadResolvedConfig(root)
	if err == nil || !strings.Contains(err.Error(), child) || !strings.Contains(err.Error(), "formats.sitemap") {
		t.Fatalf("error = %v; want included file path and nested format suggestion", err)
	}
}
