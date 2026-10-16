package config

import (
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"

	"github.com/BurntSushi/toml"
)

func TestParseTOML(t *testing.T) {
	data := []byte(`
[markata-go]
output_dir = "public"
url = "https://example.com"
title = "Test Site"

[markata-go.components.share]
enabled = false
title = "Share this"
position = "bottom"
platforms = ["twitter", "copy"]

[markata-go.components.share.custom]
mastodon = { name = "Mastodon", icon = "mastodon.svg", url = "https://mastodon.social/share?text={{title}}&url={{url}}" }

[markata-go.glob]
patterns = ["**/*.md"]
use_gitignore = true

[[markata-go.feeds]]
slug = "blog"
title = "Blog"
filter = "published == True"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.OutputDir != "public" {
		t.Errorf("OutputDir = %q, want %q", config.OutputDir, "public")
	}
	if config.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", config.URL, "https://example.com")
	}
	if len(config.Feeds) != 1 {
		t.Errorf("len(Feeds) = %d, want 1", len(config.Feeds))
	}
	if config.Components.Share.Enabled == nil || *config.Components.Share.Enabled != false {
		t.Errorf("Components.Share.Enabled = %v, want false", config.Components.Share.Enabled)
	}
	if config.Components.Share.Title != "Share this" {
		t.Errorf("Components.Share.Title = %q, want %q", config.Components.Share.Title, "Share this")
	}
	if config.Components.Share.Position != "bottom" {
		t.Errorf("Components.Share.Position = %q, want %q", config.Components.Share.Position, "bottom")
	}
	if len(config.Components.Share.Platforms) != 2 || config.Components.Share.Platforms[0] != "twitter" || config.Components.Share.Platforms[1] != "copy" {
		t.Errorf("Components.Share.Platforms = %v, want [twitter copy]", config.Components.Share.Platforms)
	}
	if config.Components.Share.Custom["mastodon"].Name != "Mastodon" {
		t.Errorf("Components.Share.Custom[mastodon].Name = %q, want %q", config.Components.Share.Custom["mastodon"].Name, "Mastodon")
	}
}

func TestParseTOML_WellKnown(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"

[markata-go.well_known]
enabled = false
auto_generate = ["time"]
ssh_fingerprint = "SHA256:abc"
keybase_username = "alice"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.WellKnown.Enabled == nil || *config.WellKnown.Enabled != false {
		t.Fatalf("WellKnown.Enabled = %v, want false", config.WellKnown.Enabled)
	}
	if len(config.WellKnown.AutoGenerate) != 1 || config.WellKnown.AutoGenerate[0] != "time" {
		t.Fatalf("WellKnown.AutoGenerate = %v, want [time]", config.WellKnown.AutoGenerate)
	}
	if config.WellKnown.SSHFingerprint != "SHA256:abc" {
		t.Fatalf("WellKnown.SSHFingerprint = %q, want %q", config.WellKnown.SSHFingerprint, "SHA256:abc")
	}
	if config.WellKnown.KeybaseUsername != "alice" {
		t.Fatalf("WellKnown.KeybaseUsername = %q, want %q", config.WellKnown.KeybaseUsername, "alice")
	}
}

func TestParseTOML_WebSub(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"

[markata-go.websub]
enabled = true
hubs = ["https://hub.example.com/"]
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.WebSub.Enabled == nil || *config.WebSub.Enabled != true {
		t.Fatalf("WebSub.Enabled = %v, want true", config.WebSub.Enabled)
	}
	if len(config.WebSub.Hubs) != 1 || config.WebSub.Hubs[0] != "https://hub.example.com/" {
		t.Fatalf("WebSub.Hubs = %v, want hub list", config.WebSub.Hubs)
	}
}

func TestParseTOML_ThemeSwitcherModeToggle(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"

[markata-go.header]
show_theme_toggle = false

[markata-go.theme.switcher]
mode_toggle = true
enabled = false
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.Theme.Switcher.ModeToggle == nil || *config.Theme.Switcher.ModeToggle != true {
		t.Fatalf("Theme.Switcher.ModeToggle = %v, want true", config.Theme.Switcher.ModeToggle)
	}
	if config.Theme.Switcher.Enabled == nil || *config.Theme.Switcher.Enabled != false {
		t.Fatalf("Theme.Switcher.Enabled = %v, want false", config.Theme.Switcher.Enabled)
	}
	if config.Header.ShowThemeToggle == nil || *config.Header.ShowThemeToggle != false {
		t.Fatalf("Header.ShowThemeToggle = %v, want false", config.Header.ShowThemeToggle)
	}
}

func TestParseThemeTextSize(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte) (*models.Config, error)
		data  []byte
	}{
		{
			name:  "toml",
			parse: ParseTOML,
			data: []byte(`
[markata-go.theme]
text_size = "medium"
show_text_size_control = false
`),
		},
		{
			name:  "yaml",
			parse: ParseYAML,
			data:  []byte("markata-go:\n  theme:\n    text_size: medium\n    show_text_size_control: false\n"),
		},
		{
			name:  "json",
			parse: ParseJSON,
			data:  []byte(`{"markata-go":{"theme":{"text_size":"medium","show_text_size_control":false}}}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tt.parse(tt.data)
			if err != nil {
				t.Fatalf("parse error = %v", err)
			}
			if config.Theme.TextSize != models.TextSizeMedium {
				t.Fatalf("TextSize = %q, want %q", config.Theme.TextSize, models.TextSizeMedium)
			}
			if config.Theme.EffectiveTextSize() != models.TextSizeMedium {
				t.Errorf("EffectiveTextSize() = %q, want %q", config.Theme.EffectiveTextSize(), models.TextSizeMedium)
			}
			if config.Theme.ShowTextSizeControl == nil || *config.Theme.ShowTextSizeControl {
				t.Errorf("ShowTextSizeControl = %v, want false", config.Theme.ShowTextSizeControl)
			}
			if config.Theme.IsTextSizeControlEnabled() {
				t.Error("IsTextSizeControlEnabled() = true, want false")
			}
		})
	}
}

func TestValidateThemeTextSize(t *testing.T) {
	config := models.NewConfig()
	config.Theme.TextSize = "gigantic"

	if errs := ValidateConfig(config); len(errs) == 0 {
		t.Fatal("ValidateConfig() returned no error for invalid theme.text_size")
	}
}

func TestParseTOML_PostConnectionsComponent(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"

[markata-go.components.post_connections]
enabled = true
display = ["list", "graph"]
graph_min_links = 3
graph_max_links = 0
list_min_links = 1
list_max_links = 5
inlinks_limit = 4
outlinks_limit = 6
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.Components.PostConnections.Enabled == nil || *config.Components.PostConnections.Enabled != true {
		t.Fatalf("PostConnections.Enabled = %v, want true", config.Components.PostConnections.Enabled)
	}
	if len(config.Components.PostConnections.Display) != 2 {
		t.Fatalf("PostConnections.Display len = %d, want 2", len(config.Components.PostConnections.Display))
	}
	if config.Components.PostConnections.ListMaxLinks != 5 {
		t.Fatalf("PostConnections.ListMaxLinks = %d, want 5", config.Components.PostConnections.ListMaxLinks)
	}
	if config.Components.PostConnections.InlinksLimit != 4 {
		t.Fatalf("PostConnections.InlinksLimit = %d, want 4", config.Components.PostConnections.InlinksLimit)
	}
	if config.Components.PostConnections.OutlinksLimit != 6 {
		t.Fatalf("PostConnections.OutlinksLimit = %d, want 6", config.Components.PostConnections.OutlinksLimit)
	}
}

func TestParseFeedSidebarMaxPostsAcrossFormats(t *testing.T) {
	t.Run("toml", func(t *testing.T) {
		config, err := ParseTOML([]byte(`
[markata-go.components.feed_sidebar]
enabled = true
max_posts = 51
`))
		if err != nil {
			t.Fatalf("ParseTOML() error = %v", err)
		}
		if config.Components.FeedSidebar.MaxPosts != 51 {
			t.Fatalf("FeedSidebar.MaxPosts = %d, want 51", config.Components.FeedSidebar.MaxPosts)
		}
	})

	t.Run("yaml", func(t *testing.T) {
		config, err := ParseYAML([]byte(`
markata-go:
  components:
    feed_sidebar:
      enabled: true
      max_posts: 51
`))
		if err != nil {
			t.Fatalf("ParseYAML() error = %v", err)
		}
		if config.Components.FeedSidebar.MaxPosts != 51 {
			t.Fatalf("FeedSidebar.MaxPosts = %d, want 51", config.Components.FeedSidebar.MaxPosts)
		}
	})

	t.Run("json", func(t *testing.T) {
		config, err := ParseJSON([]byte(`{
  "markata-go": {
    "components": {
      "feed_sidebar": {
        "enabled": true,
        "max_posts": 51
      }
    }
  }
}`))
		if err != nil {
			t.Fatalf("ParseJSON() error = %v", err)
		}
		if config.Components.FeedSidebar.MaxPosts != 51 {
			t.Fatalf("FeedSidebar.MaxPosts = %d, want 51", config.Components.FeedSidebar.MaxPosts)
		}
	})
}

func TestParseTOML_InvalidSyntax(t *testing.T) {
	data := []byte(`invalid toml {{{{ syntax`)

	_, err := ParseTOML(data)
	if err == nil {
		t.Error("ParseTOML() should return error for invalid TOML")
	}
}

func TestParseTOML_EmptySection(t *testing.T) {
	data := []byte(`
[markata-go]
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	// Should return empty config
	if config.OutputDir != "" {
		t.Errorf("OutputDir = %q, want empty", config.OutputDir)
	}
}

func TestParseFeedPrivacyOptions_AllFormats(t *testing.T) {
	tests := []struct {
		name  string
		parse func([]byte) (*models.Config, error)
		data  string
	}{
		{
			name:  "toml",
			parse: ParseTOML,
			data: `[markata-go]

[[markata-go.feeds]]
slug = "include"
include_private = true

[[markata-go.feeds]]
slug = "alias"
private = true
`,
		},
		{
			name:  "yaml",
			parse: ParseYAML,
			data: `markata-go:
  feeds:
    - slug: include
      include_private: true
    - slug: alias
      private: true
`,
		},
		{
			name:  "json",
			parse: ParseJSON,
			data:  `{"markata-go":{"feeds":[{"slug":"include","include_private":true},{"slug":"alias","private":true}]}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tt.parse([]byte(tt.data))
			if err != nil {
				t.Fatalf("parse() error = %v", err)
			}
			assertFeedPrivacyOptions(t, config)
		})
	}
}

func assertFeedPrivacyOptions(t *testing.T, config *models.Config) {
	t.Helper()
	if len(config.Feeds) != 2 {
		t.Fatalf("len(Feeds) = %d, want 2", len(config.Feeds))
	}
	if !config.Feeds[0].IncludePrivate || config.Feeds[0].Private {
		t.Errorf("include feed privacy = (include_private=%v, private=%v), want (true, false)", config.Feeds[0].IncludePrivate, config.Feeds[0].Private)
	}
	if config.Feeds[1].IncludePrivate || !config.Feeds[1].Private {
		t.Errorf("alias feed privacy = (include_private=%v, private=%v), want (false, true)", config.Feeds[1].IncludePrivate, config.Feeds[1].Private)
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte(`
markata-go:
  output_dir: public
  url: https://example.com
  title: Test Site
  components:
    share:
      enabled: false
      title: Share this
      position: bottom
      platforms:
        - twitter
        - copy
      custom:
        mastodon:
          name: Mastodon
          icon: mastodon.svg
          url: "https://mastodon.social/share?text={{title}}&url={{url}}"
  glob:
    patterns:
      - "**/*.md"
    use_gitignore: true
  feeds:
    - slug: blog
      title: Blog
      filter: "published == True"
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if config.OutputDir != "public" {
		t.Errorf("OutputDir = %q, want %q", config.OutputDir, "public")
	}
	if config.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", config.URL, "https://example.com")
	}
	if len(config.Feeds) != 1 {
		t.Errorf("len(Feeds) = %d, want 1", len(config.Feeds))
	}
	if config.Components.Share.Enabled == nil || *config.Components.Share.Enabled != false {
		t.Errorf("Components.Share.Enabled = %v, want false", config.Components.Share.Enabled)
	}
	if config.Components.Share.Title != "Share this" {
		t.Errorf("Components.Share.Title = %q, want %q", config.Components.Share.Title, "Share this")
	}
	if config.Components.Share.Position != "bottom" {
		t.Errorf("Components.Share.Position = %q, want %q", config.Components.Share.Position, "bottom")
	}
	if len(config.Components.Share.Platforms) != 2 || config.Components.Share.Platforms[0] != "twitter" || config.Components.Share.Platforms[1] != "copy" {
		t.Errorf("Components.Share.Platforms = %v, want [twitter copy]", config.Components.Share.Platforms)
	}
	if config.Components.Share.Custom["mastodon"].Icon != "mastodon.svg" {
		t.Errorf("Components.Share.Custom[mastodon].Icon = %q, want %q", config.Components.Share.Custom["mastodon"].Icon, "mastodon.svg")
	}
}

func TestParseYAML_InvalidSyntax(t *testing.T) {
	data := []byte(`
markata-go:
  invalid:
    - unclosed list
   bad indent
`)

	_, err := ParseYAML(data)
	if err == nil {
		t.Error("ParseYAML() should return error for invalid YAML")
	}
}

func TestParseYAML_EmptySection(t *testing.T) {
	data := []byte(`
markata-go:
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if config.OutputDir != "" {
		t.Errorf("OutputDir = %q, want empty", config.OutputDir)
	}
}

func TestParseJSON(t *testing.T) {
	data := []byte(`{
  "markata-go": {
    "output_dir": "public",
    "url": "https://example.com",
    "title": "Test Site",
    "components": {
      "share": {
        "enabled": false,
        "title": "Share this",
        "position": "bottom",
        "platforms": ["twitter", "copy"],
        "custom": {
          "mastodon": {
            "name": "Mastodon",
            "icon": "mastodon.svg",
            "url": "https://mastodon.social/share?text={{title}}&url={{url}}"
          }
        }
      }
    },
    "glob": {
      "patterns": ["**/*.md"],
      "use_gitignore": true
    },
    "feeds": [
      {
        "slug": "blog",
        "title": "Blog",
        "filter": "published == True"
      }
    ]
  }
}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if config.OutputDir != "public" {
		t.Errorf("OutputDir = %q, want %q", config.OutputDir, "public")
	}
	if config.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", config.URL, "https://example.com")
	}
	if len(config.Feeds) != 1 {
		t.Errorf("len(Feeds) = %d, want 1", len(config.Feeds))
	}
	if config.Components.Share.Enabled == nil || *config.Components.Share.Enabled != false {
		t.Errorf("Components.Share.Enabled = %v, want false", config.Components.Share.Enabled)
	}
	if config.Components.Share.Title != "Share this" {
		t.Errorf("Components.Share.Title = %q, want %q", config.Components.Share.Title, "Share this")
	}
	if config.Components.Share.Position != "bottom" {
		t.Errorf("Components.Share.Position = %q, want %q", config.Components.Share.Position, "bottom")
	}
	if len(config.Components.Share.Platforms) != 2 || config.Components.Share.Platforms[0] != "twitter" || config.Components.Share.Platforms[1] != "copy" {
		t.Errorf("Components.Share.Platforms = %v, want [twitter copy]", config.Components.Share.Platforms)
	}
	if config.Components.Share.Custom["mastodon"].URL == "" {
		t.Error("Components.Share.Custom[mastodon].URL should be set")
	}
}

func TestParseJSON_InvalidSyntax(t *testing.T) {
	data := []byte(`{invalid json}`)

	_, err := ParseJSON(data)
	if err == nil {
		t.Error("ParseJSON() should return error for invalid JSON")
	}
}

func TestParseJSON_EmptySection(t *testing.T) {
	data := []byte(`{"markata-go": {}}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if config.OutputDir != "" {
		t.Errorf("OutputDir = %q, want empty", config.OutputDir)
	}
}

func TestParseTOML_SearchConfig(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"

[markata-go.search]
backend = "bleve"
endpoint = "http://localhost:3001/api/search"

[markata-go.search.bleve]
fuzzy = true
limit = 20
max_limit = 100
cors_origins = ["*"]
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.Search.Backend != "bleve" {
		t.Fatalf("Search.Backend = %q, want %q", config.Search.Backend, "bleve")
	}
	if config.Search.Endpoint != "http://localhost:3001/api/search" {
		t.Fatalf("Search.Endpoint = %q, want %q", config.Search.Endpoint, "http://localhost:3001/api/search")
	}
	if config.Search.Bleve.Fuzzy == nil || !*config.Search.Bleve.Fuzzy {
		t.Fatalf("Search.Bleve.Fuzzy = %v, want true", config.Search.Bleve.Fuzzy)
	}
	if config.Search.Bleve.Limit != 20 {
		t.Fatalf("Search.Bleve.Limit = %d, want 20", config.Search.Bleve.Limit)
	}
	if config.Search.Bleve.MaxLimit != 100 {
		t.Fatalf("Search.Bleve.MaxLimit = %d, want 100", config.Search.Bleve.MaxLimit)
	}
	if len(config.Search.Bleve.CORSOrigins) != 1 || config.Search.Bleve.CORSOrigins[0] != "*" {
		t.Fatalf("Search.Bleve.CORSOrigins = %v, want [*]", config.Search.Bleve.CORSOrigins)
	}
}

func TestParseYAML_SearchConfig(t *testing.T) {
	data := []byte(`
markata-go:
  title: Test Site
  search:
    backend: bleve
    endpoint: http://localhost:3001/api/search
    bleve:
      fuzzy: true
      limit: 20
      max_limit: 100
      cors_origins:
        - "*"
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if config.Search.Backend != "bleve" {
		t.Fatalf("Search.Backend = %q, want %q", config.Search.Backend, "bleve")
	}
	if config.Search.Endpoint != "http://localhost:3001/api/search" {
		t.Fatalf("Search.Endpoint = %q, want %q", config.Search.Endpoint, "http://localhost:3001/api/search")
	}
	if config.Search.Bleve.Fuzzy == nil || !*config.Search.Bleve.Fuzzy {
		t.Fatalf("Search.Bleve.Fuzzy = %v, want true", config.Search.Bleve.Fuzzy)
	}
	if config.Search.Bleve.Limit != 20 {
		t.Fatalf("Search.Bleve.Limit = %d, want 20", config.Search.Bleve.Limit)
	}
}

func TestParseJSON_SearchConfig(t *testing.T) {
	data := []byte(`{
  "markata-go": {
    "title": "Test Site",
    "search": {
      "backend": "bleve",
      "endpoint": "http://localhost:3001/api/search",
      "bleve": {
        "fuzzy": true,
        "limit": 20,
        "max_limit": 100,
        "cors_origins": ["*"]
      }
    }
  }
}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if config.Search.Backend != "bleve" {
		t.Fatalf("Search.Backend = %q, want %q", config.Search.Backend, "bleve")
	}
	if config.Search.Endpoint != "http://localhost:3001/api/search" {
		t.Fatalf("Search.Endpoint = %q, want %q", config.Search.Endpoint, "http://localhost:3001/api/search")
	}
	if config.Search.Bleve.Fuzzy == nil || !*config.Search.Bleve.Fuzzy {
		t.Fatalf("Search.Bleve.Fuzzy = %v, want true", config.Search.Bleve.Fuzzy)
	}
	if config.Search.Bleve.Limit != 20 {
		t.Fatalf("Search.Bleve.Limit = %d, want 20", config.Search.Bleve.Limit)
	}
}

func TestParseTOML_UseGitignoreFalse(t *testing.T) {
	data := []byte(`
[markata-go]
[markata-go.glob]
use_gitignore = false
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.GlobConfig.UseGitignore != false {
		t.Error("UseGitignore should be false")
	}
}

func TestParseTOML_GlobSlugMode(t *testing.T) {
	data := []byte(`
[markata-go]
[markata-go.glob]
slug_mode = "path"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.GlobConfig.SlugMode != "path" {
		t.Errorf("SlugMode = %q, want %q", config.GlobConfig.SlugMode, "path")
	}
}

func TestParseTOML_GlobSlugRules(t *testing.T) {
	data := []byte(`
[markata-go]
[markata-go.glob]
slug_mode = "flat"

[[markata-go.glob.slug_rules]]
prefix = "posts/blog"
mode = "flat"

[[markata-go.glob.slug_rules]]
prefix = "posts/notes"
mode = "path"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.GlobConfig.SlugRules) != 2 {
		t.Fatalf("len(SlugRules) = %d, want 2", len(config.GlobConfig.SlugRules))
	}
	if config.GlobConfig.SlugRules[1].Prefix != "posts/notes" {
		t.Errorf("Prefix = %q, want %q", config.GlobConfig.SlugRules[1].Prefix, "posts/notes")
	}
}

func TestParseTOML_GlobSlugRulesDecoder(t *testing.T) {
	data := []byte(`
[markata-go]

[markata-go.glob]
slug_mode = "flat"

[[markata-go.glob.slug_rules]]
prefix = "posts/blog"
mode = "flat"

[[markata-go.glob.slug_rules]]
prefix = "posts/notes"
mode = "path"
`)

	var wrapper struct {
		MarkataGo tomlConfig `toml:"markata-go"`
	}
	if err := toml.Unmarshal(data, &wrapper); err != nil {
		t.Fatalf("toml.Unmarshal() error = %v", err)
	}

	if len(wrapper.MarkataGo.Glob.SlugRules) != 2 {
		t.Fatalf("decoder len(SlugRules) = %d, want 2", len(wrapper.MarkataGo.Glob.SlugRules))
	}
}

func TestParseTOML_FeedFormats(t *testing.T) {
	data := []byte(`
[markata-go]
[[markata-go.feeds]]
slug = "blog"
[markata-go.feeds.formats]
html = true
rss = true
atom = false
json = false
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Feeds) != 1 {
		t.Fatalf("len(Feeds) = %d, want 1", len(config.Feeds))
	}

	formats := config.Feeds[0].Formats
	if !formats.HTML {
		t.Error("HTML should be true")
	}
	if !formats.RSS {
		t.Error("RSS should be true")
	}
	if formats.Atom {
		t.Error("Atom should be false")
	}
	if formats.JSON {
		t.Error("JSON should be false")
	}
}

func TestParseTOML_FeedDefaults(t *testing.T) {
	data := []byte(`
[markata-go]
[markata-go.feed_defaults]
items_per_page = 15
orphan_threshold = 4
[markata-go.feed_defaults.formats]
html = true
rss = true
[markata-go.feed_defaults.syndication]
max_items = 30
include_content = true
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.FeedDefaults.ItemsPerPage != 15 {
		t.Errorf("ItemsPerPage = %d, want 15", config.FeedDefaults.ItemsPerPage)
	}
	if config.FeedDefaults.OrphanThreshold != 4 {
		t.Errorf("OrphanThreshold = %d, want 4", config.FeedDefaults.OrphanThreshold)
	}
	if !config.FeedDefaults.Formats.HTML {
		t.Error("Formats.HTML should be true")
	}
	if config.FeedDefaults.Syndication.MaxItems != 30 {
		t.Errorf("Syndication.MaxItems = %d, want 30", config.FeedDefaults.Syndication.MaxItems)
	}
}

func TestParseTOML_FeedLimitOffset(t *testing.T) {
	data := []byte(`
[markata-go]
[[markata-go.feeds]]
slug = "blog"
limit = 3
offset = 2
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Feeds) != 1 {
		t.Fatalf("len(Feeds) = %d, want 1", len(config.Feeds))
	}

	feed := config.Feeds[0]
	if feed.Limit != 3 {
		t.Errorf("Limit = %d, want 3", feed.Limit)
	}
	if feed.Offset != 2 {
		t.Errorf("Offset = %d, want 2", feed.Offset)
	}
}

func TestParseYAML_FeedLimitOffset(t *testing.T) {
	data := []byte(`
markata-go:
  feeds:
    - slug: blog
      limit: 4
      offset: 1
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if len(config.Feeds) != 1 {
		t.Fatalf("len(Feeds) = %d, want 1", len(config.Feeds))
	}

	feed := config.Feeds[0]
	if feed.Limit != 4 {
		t.Errorf("Limit = %d, want 4", feed.Limit)
	}
	if feed.Offset != 1 {
		t.Errorf("Offset = %d, want 1", feed.Offset)
	}
}

func TestParseJSON_FeedLimitOffset(t *testing.T) {
	data := []byte(`
{
  "markata-go": {
    "feeds": [
      {
        "slug": "blog",
        "limit": 2,
        "offset": 3
      }
    ]
  }
}
`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if len(config.Feeds) != 1 {
		t.Fatalf("len(Feeds) = %d, want 1", len(config.Feeds))
	}

	feed := config.Feeds[0]
	if feed.Limit != 2 {
		t.Errorf("Limit = %d, want 2", feed.Limit)
	}
	if feed.Offset != 3 {
		t.Errorf("Offset = %d, want 3", feed.Offset)
	}
}

func TestParseTOML_MarkdownConfig(t *testing.T) {
	data := []byte(`
[markata-go]
[markata-go.markdown]
extensions = ["tables", "footnotes", "syntax-highlighting"]
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	expected := []string{"tables", "footnotes", "syntax-highlighting"}
	if len(config.MarkdownConfig.Extensions) != len(expected) {
		t.Fatalf("len(Extensions) = %d, want %d", len(config.MarkdownConfig.Extensions), len(expected))
	}
	for i, ext := range expected {
		if config.MarkdownConfig.Extensions[i] != ext {
			t.Errorf("Extensions[%d] = %q, want %q", i, config.MarkdownConfig.Extensions[i], ext)
		}
	}
}

func TestParseTOML_Hooks(t *testing.T) {
	data := []byte(`
[markata-go]
hooks = ["markdown", "template", "sitemap"]
disabled_hooks = ["seo", "analytics"]
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Hooks) != 3 {
		t.Errorf("len(Hooks) = %d, want 3", len(config.Hooks))
	}
	if len(config.DisabledHooks) != 2 {
		t.Errorf("len(DisabledHooks) = %d, want 2", len(config.DisabledHooks))
	}
}

func TestParseTOML_MultipleFeeds(t *testing.T) {
	data := []byte(`
[markata-go]
[[markata-go.feeds]]
slug = "blog"
title = "Blog"

[[markata-go.feeds]]
slug = "projects"
title = "Projects"

[[markata-go.feeds]]
slug = "notes"
title = "Notes"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Feeds) != 3 {
		t.Fatalf("len(Feeds) = %d, want 3", len(config.Feeds))
	}

	expectedSlugs := []string{"blog", "projects", "notes"}
	for i, slug := range expectedSlugs {
		if config.Feeds[i].Slug != slug {
			t.Errorf("Feeds[%d].Slug = %q, want %q", i, config.Feeds[i].Slug, slug)
		}
	}
}

func TestParseTOML_FeedTemplates(t *testing.T) {
	data := []byte(`
[markata-go]
[[markata-go.feeds]]
slug = "blog"
[markata-go.feeds.templates]
html = "custom-feed.html"
rss = "custom-rss.xml"
atom = "custom-atom.xml"
card = "custom-card.html"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Feeds) != 1 {
		t.Fatalf("len(Feeds) = %d, want 1", len(config.Feeds))
	}

	templates := config.Feeds[0].Templates
	if templates.HTML != "custom-feed.html" {
		t.Errorf("Templates.HTML = %q, want %q", templates.HTML, "custom-feed.html")
	}
	if templates.RSS != "custom-rss.xml" {
		t.Errorf("Templates.RSS = %q, want %q", templates.RSS, "custom-rss.xml")
	}
	if templates.Atom != "custom-atom.xml" {
		t.Errorf("Templates.Atom = %q, want %q", templates.Atom, "custom-atom.xml")
	}
	if templates.Card != "custom-card.html" {
		t.Errorf("Templates.Card = %q, want %q", templates.Card, "custom-card.html")
	}
}

func TestParseYAML_NestedConfig(t *testing.T) {
	data := []byte(`
markata-go:
  output_dir: public
  glob:
    patterns:
      - "posts/**/*.md"
      - "pages/*.md"
    use_gitignore: true
  feed_defaults:
    items_per_page: 20
    formats:
      html: true
      rss: true
    syndication:
      max_items: 50
      include_content: true
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if len(config.GlobConfig.Patterns) != 2 {
		t.Errorf("len(Patterns) = %d, want 2", len(config.GlobConfig.Patterns))
	}
	if config.FeedDefaults.ItemsPerPage != 20 {
		t.Errorf("ItemsPerPage = %d, want 20", config.FeedDefaults.ItemsPerPage)
	}
	if config.FeedDefaults.Syndication.MaxItems != 50 {
		t.Errorf("Syndication.MaxItems = %d, want 50", config.FeedDefaults.Syndication.MaxItems)
	}
}

func TestParseJSON_NestedConfig(t *testing.T) {
	data := []byte(`{
  "markata-go": {
    "output_dir": "public",
    "glob": {
      "patterns": ["posts/**/*.md", "pages/*.md"],
	      "use_gitignore": true,
	      "slug_mode": "path"
    },
    "feed_defaults": {
      "items_per_page": 20,
      "formats": {
        "html": true,
        "rss": true
      },
      "syndication": {
        "max_items": 50,
        "include_content": true
      }
    }
  }
}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if len(config.GlobConfig.Patterns) != 2 {
		t.Errorf("len(Patterns) = %d, want 2", len(config.GlobConfig.Patterns))
	}
	if config.GlobConfig.SlugMode != "path" {
		t.Errorf("SlugMode = %q, want %q", config.GlobConfig.SlugMode, "path")
	}
	if config.FeedDefaults.ItemsPerPage != 20 {
		t.Errorf("ItemsPerPage = %d, want 20", config.FeedDefaults.ItemsPerPage)
	}
	if config.FeedDefaults.Syndication.MaxItems != 50 {
		t.Errorf("Syndication.MaxItems = %d, want 50", config.FeedDefaults.Syndication.MaxItems)
	}
}

func TestParseTOML_AllFields(t *testing.T) {
	data := []byte(`
[markata-go]
output_dir = "public"
url = "https://example.com"
title = "My Site"
description = "A great site"
author = "John Doe"
assets_dir = "assets"
templates_dir = "themes/default"
hooks = ["markdown", "template"]
disabled_hooks = ["seo"]
concurrency = 4

[markata-go.glob]
patterns = ["**/*.md"]
use_gitignore = true

[markata-go.markdown]
extensions = ["tables"]

[[markata-go.feeds]]
slug = "blog"
title = "Blog"
description = "All posts"
filter = "published == True"
sort = "date"
reverse = true
items_per_page = 10
orphan_threshold = 3

[markata-go.feeds.formats]
html = true
rss = true

[markata-go.feed_defaults]
items_per_page = 15
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	// Verify all top-level fields
	tests := []struct {
		field string
		got   interface{}
		want  interface{}
	}{
		{"OutputDir", config.OutputDir, "public"},
		{"URL", config.URL, "https://example.com"},
		{"Title", config.Title, "My Site"},
		{"Description", config.Description, "A great site"},
		{"Author", config.Author, "John Doe"},
		{"AssetsDir", config.AssetsDir, "assets"},
		{"TemplatesDir", config.TemplatesDir, "themes/default"},
		{"Concurrency", config.Concurrency, 4},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.field, tt.got, tt.want)
		}
	}
}

func TestParseAssetsConfigAcrossFormats(t *testing.T) {
	tests := []struct {
		name  string
		data  []byte
		parse func([]byte) (*models.Config, error)
	}{
		{
			name: "TOML",
			data: []byte(`[markata-go]
[markata-go.assets]
mode = "auto"
cache_dir = "tmp/assets"
output_dir = "vendor"
verify_integrity = false
`),
			parse: ParseTOML,
		},
		{
			name: "YAML",
			data: []byte(`markata-go:
  assets:
    mode: auto
    cache_dir: tmp/assets
    output_dir: vendor
    verify_integrity: false
`),
			parse: ParseYAML,
		},
		{
			name:  "JSON",
			data:  []byte(`{"markata-go":{"assets":{"mode":"auto","cache_dir":"tmp/assets","output_dir":"vendor","verify_integrity":false}}}`),
			parse: ParseJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := tt.parse(tt.data)
			if err != nil {
				t.Fatalf("parse assets config: %v", err)
			}
			if config.Assets.Mode != "auto" || config.Assets.CacheDir != "tmp/assets" || config.Assets.OutputDir != "vendor" {
				t.Fatalf("assets config = %+v", config.Assets)
			}
			if config.Assets.VerifyIntegrity == nil || *config.Assets.VerifyIntegrity {
				t.Fatalf("assets verify_integrity = %v, want false", config.Assets.VerifyIntegrity)
			}
		})
	}
}

func TestParseTOML_BlogrollExternalFeedFields(t *testing.T) {
	maxEntries := 5
	primary := true

	data := []byte(`
[markata-go]

[markata-go.blogroll]
enabled = true

[[markata-go.blogroll.feeds]]
url = "https://example.com/feed.xml"
title = "Example Blog"
handle = "exampleblog"
aliases = ["example", "ex"]
max_entries = 5
primary = true
primary_person = "mainauthor"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Blogroll.Feeds) != 1 {
		t.Fatalf("len(Blogroll.Feeds) = %d, want 1", len(config.Blogroll.Feeds))
	}

	feed := config.Blogroll.Feeds[0]

	// Verify all external feed fields are parsed correctly
	tests := []struct {
		field string
		got   interface{}
		want  interface{}
	}{
		{"URL", feed.URL, "https://example.com/feed.xml"},
		{"Title", feed.Title, "Example Blog"},
		{"Handle", feed.Handle, "exampleblog"},
		{"PrimaryPerson", feed.PrimaryPerson, "mainauthor"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %v, want %v", tt.field, tt.got, tt.want)
		}
	}

	// Check slice fields
	if len(feed.Aliases) != 2 {
		t.Errorf("len(Aliases) = %d, want 2", len(feed.Aliases))
	} else if feed.Aliases[0] != "example" || feed.Aliases[1] != "ex" {
		t.Errorf("Aliases = %v, want [example ex]", feed.Aliases)
	}

	// Check pointer fields
	if feed.MaxEntries == nil || *feed.MaxEntries != maxEntries {
		t.Errorf("MaxEntries = %v, want %d", feed.MaxEntries, maxEntries)
	}
	if feed.Primary == nil || *feed.Primary != primary {
		t.Errorf("Primary = %v, want %v", feed.Primary, primary)
	}
}

func TestParseYAML_BlogrollExternalFeedFields(t *testing.T) {
	maxEntries := 10
	primary := false

	data := []byte(`
markata-go:
  blogroll:
    enabled: true
    feeds:
      - url: "https://blog.example.org/rss"
        title: "Another Blog"
        handle: "anotherblog"
        aliases:
          - another
          - blog
        max_entries: 10
        primary: false
        primary_person: "someauthor"
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if len(config.Blogroll.Feeds) != 1 {
		t.Fatalf("len(Blogroll.Feeds) = %d, want 1", len(config.Blogroll.Feeds))
	}

	feed := config.Blogroll.Feeds[0]

	if feed.URL != "https://blog.example.org/rss" {
		t.Errorf("URL = %q, want %q", feed.URL, "https://blog.example.org/rss")
	}
	if feed.Handle != "anotherblog" {
		t.Errorf("Handle = %q, want %q", feed.Handle, "anotherblog")
	}
	if feed.PrimaryPerson != "someauthor" {
		t.Errorf("PrimaryPerson = %q, want %q", feed.PrimaryPerson, "someauthor")
	}
	if len(feed.Aliases) != 2 {
		t.Errorf("len(Aliases) = %d, want 2", len(feed.Aliases))
	}
	if feed.MaxEntries == nil || *feed.MaxEntries != maxEntries {
		t.Errorf("MaxEntries = %v, want %d", feed.MaxEntries, maxEntries)
	}
	if feed.Primary == nil || *feed.Primary != primary {
		t.Errorf("Primary = %v, want %v", feed.Primary, primary)
	}
}

func TestParseJSON_BlogrollExternalFeedFields(t *testing.T) {
	maxEntries := 3

	data := []byte(`{
  "markata-go": {
    "blogroll": {
      "enabled": true,
      "feeds": [
        {
          "url": "https://json.example.com/feed",
          "title": "JSON Blog",
          "handle": "jsonblog",
          "aliases": ["json", "jb"],
          "max_entries": 3,
          "primary": true,
          "primary_person": "jsonauthor"
        }
      ]
    }
  }
}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if len(config.Blogroll.Feeds) != 1 {
		t.Fatalf("len(Blogroll.Feeds) = %d, want 1", len(config.Blogroll.Feeds))
	}

	feed := config.Blogroll.Feeds[0]

	if feed.URL != "https://json.example.com/feed" {
		t.Errorf("URL = %q, want %q", feed.URL, "https://json.example.com/feed")
	}
	if feed.Handle != "jsonblog" {
		t.Errorf("Handle = %q, want %q", feed.Handle, "jsonblog")
	}
	if feed.PrimaryPerson != "jsonauthor" {
		t.Errorf("PrimaryPerson = %q, want %q", feed.PrimaryPerson, "jsonauthor")
	}
	if len(feed.Aliases) != 2 || feed.Aliases[0] != "json" || feed.Aliases[1] != "jb" {
		t.Errorf("Aliases = %v, want [json jb]", feed.Aliases)
	}
	if feed.MaxEntries == nil || *feed.MaxEntries != maxEntries {
		t.Errorf("MaxEntries = %v, want %d", feed.MaxEntries, maxEntries)
	}
	if feed.Primary == nil || !*feed.Primary {
		t.Errorf("Primary = %v, want true", feed.Primary)
	}
}

// TestParseTOML_BlogrollFallbackImageService tests parsing of fallback_image_service from TOML
func TestParseTOML_BlogrollFallbackImageService(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"

[markata-go.blogroll]
enabled = true
fallback_image_service = "https://shots.waylonwalker.com/shot/?url={url}&height=160&width=240"

[[markata-go.blogroll.feeds]]
url = "https://simonwillison.net/atom/everything/"
title = "Simon Willison"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if !config.Blogroll.Enabled {
		t.Error("Blogroll.Enabled should be true")
	}

	expectedURL := "https://shots.waylonwalker.com/shot/?url={url}&height=160&width=240"
	if config.Blogroll.FallbackImageService != expectedURL {
		t.Errorf("Blogroll.FallbackImageService = %q, want %q",
			config.Blogroll.FallbackImageService, expectedURL)
	}

	if len(config.Blogroll.Feeds) != 1 {
		t.Errorf("len(Blogroll.Feeds) = %d, want 1", len(config.Blogroll.Feeds))
	}
}

// TestParseYAML_BlogrollFallbackImageService tests parsing of fallback_image_service from YAML
func TestParseYAML_BlogrollFallbackImageService(t *testing.T) {
	data := []byte(`
markata-go:
  title: "Test Site"
  blogroll:
    enabled: true
    fallback_image_service: "https://shots.waylonwalker.com/shot/?url={url}&height=160&width=240"
    feeds:
      - url: "https://simonwillison.net/atom/everything/"
        title: "Simon Willison"
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if !config.Blogroll.Enabled {
		t.Error("Blogroll.Enabled should be true")
	}

	expectedURL := "https://shots.waylonwalker.com/shot/?url={url}&height=160&width=240"
	if config.Blogroll.FallbackImageService != expectedURL {
		t.Errorf("Blogroll.FallbackImageService = %q, want %q",
			config.Blogroll.FallbackImageService, expectedURL)
	}

	if len(config.Blogroll.Feeds) != 1 {
		t.Errorf("len(Blogroll.Feeds) = %d, want 1", len(config.Blogroll.Feeds))
	}
}

// TestParseJSON_BlogrollFallbackImageService tests parsing of fallback_image_service from JSON
func TestParseJSON_BlogrollFallbackImageService(t *testing.T) {
	data := []byte(`{
  "markata-go": {
    "title": "Test Site",
    "blogroll": {
      "enabled": true,
      "fallback_image_service": "https://shots.waylonwalker.com/shot/?url={url}&height=160&width=240",
      "feeds": [
        {
          "url": "https://simonwillison.net/atom/everything/",
          "title": "Simon Willison"
        }
      ]
    }
  }
}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if !config.Blogroll.Enabled {
		t.Error("Blogroll.Enabled should be true")
	}

	expectedURL := "https://shots.waylonwalker.com/shot/?url={url}&height=160&width=240"
	if config.Blogroll.FallbackImageService != expectedURL {
		t.Errorf("Blogroll.FallbackImageService = %q, want %q",
			config.Blogroll.FallbackImageService, expectedURL)
	}

	if len(config.Blogroll.Feeds) != 1 {
		t.Errorf("len(Blogroll.Feeds) = %d, want 1", len(config.Blogroll.Feeds))
	}
}

// TestParseTOML_Shortcuts tests parsing of shortcuts from TOML config
func TestParseTOML_Shortcuts(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"

[markata-go.shortcuts.navigation]
"g t" = "/tags/"
"g a" = "/about/"
"g h" = "/"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if config.Shortcuts.Navigation == nil {
		t.Fatal("Shortcuts.Navigation should not be nil")
	}

	if len(config.Shortcuts.Navigation) != 3 {
		t.Errorf("len(Shortcuts.Navigation) = %d, want 3", len(config.Shortcuts.Navigation))
	}

	tests := []struct {
		key  string
		want string
	}{
		{"g t", "/tags/"},
		{"g a", "/about/"},
		{"g h", "/"},
	}

	for _, tt := range tests {
		if got := config.Shortcuts.Navigation[tt.key]; got != tt.want {
			t.Errorf("Shortcuts.Navigation[%q] = %q, want %q", tt.key, got, tt.want)
		}
	}
}

// TestParseTOML_Shortcuts_Empty tests parsing of empty shortcuts from TOML config
func TestParseTOML_Shortcuts_Empty(t *testing.T) {
	data := []byte(`
[markata-go]
title = "Test Site"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	// Empty shortcuts should still work (navigation may be nil or empty)
	if len(config.Shortcuts.Navigation) != 0 {
		t.Errorf("Shortcuts.Navigation should be nil or empty, got %v", config.Shortcuts.Navigation)
	}
}

// TestParseYAML_Shortcuts tests parsing of shortcuts from YAML config
func TestParseYAML_Shortcuts(t *testing.T) {
	data := []byte(`
markata-go:
  title: Test Site
  shortcuts:
    navigation:
      "g t": "/tags/"
      "g a": "/about/"
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if config.Shortcuts.Navigation == nil {
		t.Fatal("Shortcuts.Navigation should not be nil")
	}

	if len(config.Shortcuts.Navigation) != 2 {
		t.Errorf("len(Shortcuts.Navigation) = %d, want 2", len(config.Shortcuts.Navigation))
	}

	if got := config.Shortcuts.Navigation["g t"]; got != "/tags/" {
		t.Errorf("Shortcuts.Navigation[\"g t\"] = %q, want %q", got, "/tags/")
	}

	if got := config.Shortcuts.Navigation["g a"]; got != "/about/" {
		t.Errorf("Shortcuts.Navigation[\"g a\"] = %q, want %q", got, "/about/")
	}
}

// TestParseJSON_Shortcuts tests parsing of shortcuts from JSON config
func TestParseJSON_Shortcuts(t *testing.T) {
	data := []byte(`{
  "markata-go": {
    "title": "Test Site",
    "shortcuts": {
      "navigation": {
        "g t": "/tags/",
        "g a": "/about/"
      }
    }
  }
}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if config.Shortcuts.Navigation == nil {
		t.Fatal("Shortcuts.Navigation should not be nil")
	}

	if len(config.Shortcuts.Navigation) != 2 {
		t.Errorf("len(Shortcuts.Navigation) = %d, want 2", len(config.Shortcuts.Navigation))
	}

	if got := config.Shortcuts.Navigation["g t"]; got != "/tags/" {
		t.Errorf("Shortcuts.Navigation[\"g t\"] = %q, want %q", got, "/tags/")
	}

	if got := config.Shortcuts.Navigation["g a"]; got != "/about/" {
		t.Errorf("Shortcuts.Navigation[\"g a\"] = %q, want %q", got, "/about/")
	}
}

// TestParseTOML_MentionsFromPostsDefault tests that default from_posts source
// is applied when mentions section exists but no from_posts are configured.
func TestParseTOML_MentionsFromPostsDefault(t *testing.T) {
	data := []byte(`
[markata-go]
output_dir = "public"

[markata-go.mentions]
css_class = "mention"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Mentions.FromPosts) == 0 {
		t.Fatal("Mentions.FromPosts should have default entries when not configured")
	}

	if config.Mentions.FromPosts[0].Filter != "template == 'contact'" {
		t.Errorf("default FromPosts[0].Filter = %q, want %q",
			config.Mentions.FromPosts[0].Filter, "template == 'contact'")
	}

	if config.Mentions.FromPosts[0].HandleField != "handle" {
		t.Errorf("default FromPosts[0].HandleField = %q, want %q",
			config.Mentions.FromPosts[0].HandleField, "handle")
	}
}

// TestParseTOML_MentionsFromPostsExplicit tests that explicit from_posts
// overrides the defaults and includes AvatarField.
func TestParseTOML_MentionsFromPostsExplicit(t *testing.T) {
	data := []byte(`
[markata-go]
output_dir = "public"

[markata-go.mentions]
css_class = "mention"

[[markata-go.mentions.from_posts]]
filter = "'team' in tags"
handle_field = "github_handle"
aliases_field = "nicknames"
avatar_field = "photo"
`)

	config, err := ParseTOML(data)
	if err != nil {
		t.Fatalf("ParseTOML() error = %v", err)
	}

	if len(config.Mentions.FromPosts) != 1 {
		t.Fatalf("len(Mentions.FromPosts) = %d, want 1", len(config.Mentions.FromPosts))
	}

	src := config.Mentions.FromPosts[0]
	if src.Filter != "'team' in tags" {
		t.Errorf("FromPosts[0].Filter = %q, want %q", src.Filter, "'team' in tags")
	}
	if src.HandleField != "github_handle" {
		t.Errorf("FromPosts[0].HandleField = %q, want %q", src.HandleField, "github_handle")
	}
	if src.AliasesField != "nicknames" {
		t.Errorf("FromPosts[0].AliasesField = %q, want %q", src.AliasesField, "nicknames")
	}
	if src.AvatarField != "photo" {
		t.Errorf("FromPosts[0].AvatarField = %q, want %q", src.AvatarField, "photo")
	}
}

// TestParseYAML_MentionsFromPostsDefault tests default from_posts in YAML config.
func TestParseYAML_MentionsFromPostsDefault(t *testing.T) {
	data := []byte(`
markata-go:
  output_dir: public
  mentions:
    css_class: mention
`)

	config, err := ParseYAML(data)
	if err != nil {
		t.Fatalf("ParseYAML() error = %v", err)
	}

	if len(config.Mentions.FromPosts) == 0 {
		t.Fatal("Mentions.FromPosts should have default entries when not configured")
	}

	if config.Mentions.FromPosts[0].Filter != "template == 'contact'" {
		t.Errorf("default FromPosts[0].Filter = %q, want %q",
			config.Mentions.FromPosts[0].Filter, "template == 'contact'")
	}
}

// TestParseJSON_MentionsFromPostsDefault tests default from_posts in JSON config.
func TestParseJSON_MentionsFromPostsDefault(t *testing.T) {
	data := []byte(`{
  "markata-go": {
    "output_dir": "public",
    "mentions": {
      "css_class": "mention"
    }
  }
}`)

	config, err := ParseJSON(data)
	if err != nil {
		t.Fatalf("ParseJSON() error = %v", err)
	}

	if len(config.Mentions.FromPosts) == 0 {
		t.Fatal("Mentions.FromPosts should have default entries when not configured")
	}

	if config.Mentions.FromPosts[0].Filter != "template == 'contact'" {
		t.Errorf("default FromPosts[0].Filter = %q, want %q",
			config.Mentions.FromPosts[0].Filter, "template == 'contact'")
	}
}

func TestViewTransitionsNativeIsOptInAcrossFormats(t *testing.T) {
	config, err := ParseTOML([]byte(`
[markata-go.view_transitions]
enabled = true
native = true
`))
	if err != nil {
		t.Fatal(err)
	}
	if config.ViewTransitions.Native == nil || !*config.ViewTransitions.Native {
		t.Fatalf("native navigation setting was not parsed: %+v", config.ViewTransitions)
	}

	enabled := true
	for name, got := range map[string]models.ViewTransitionsConfig{
		"toml": (&tomlViewTransitionsConfig{Native: &enabled}).toViewTransitionsConfig(),
		"yaml": (&yamlViewTransitionsConfig{Native: &enabled}).toViewTransitionsConfig(),
		"json": (&jsonViewTransitionsConfig{Native: &enabled}).toViewTransitionsConfig(),
	} {
		if got.Native == nil || !*got.Native {
			t.Fatalf("%s native navigation missing: %+v", name, got)
		}
	}

	defaults := models.NewViewTransitionsConfig()
	if defaults.Native != nil && *defaults.Native {
		t.Fatalf("native navigation must be opt-in: %+v", defaults)
	}

	disabled := false
	merged := mergeViewTransitionsConfig(defaults, models.ViewTransitionsConfig{Native: &enabled})
	if merged.Native == nil || !*merged.Native {
		t.Fatalf("native flag lost during merge: %+v", merged)
	}
	merged = mergeViewTransitionsConfig(merged, models.ViewTransitionsConfig{Native: &disabled})
	if merged.Native == nil || *merged.Native {
		t.Fatalf("explicit native=false ignored during merge: %+v", merged)
	}
}

func TestParsedFeedFormatFalseAcrossFormats(t *testing.T) {
	testCases := []struct {
		name  string
		data  string
		parse func([]byte) (*models.Config, error)
	}{
		{
			name:  "toml",
			data:  "[markata-go]\n[[markata-go.feeds]]\nslug = \"draft\"\n[markata-go.feeds.formats]\nsitemap = false\n",
			parse: ParseTOML,
		},
		{
			name:  "yaml",
			data:  "markata-go:\n  feeds:\n    - slug: draft\n      formats:\n        sitemap: false\n",
			parse: ParseYAML,
		},
		{
			name:  "json",
			data:  `{"markata-go":{"feeds":[{"slug":"draft","formats":{"sitemap":false}}]}}`,
			parse: ParseJSON,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := tc.parse([]byte(tc.data))
			if err != nil {
				t.Fatal(err)
			}
			if len(cfg.Feeds) != 1 {
				t.Fatalf("expected 1 feed, got %d", len(cfg.Feeds))
			}
			feed := cfg.Feeds[0]
			feed.ApplyDefaults(models.NewFeedDefaults())
			if feed.Formats.Sitemap || !feed.Formats.HTML || !feed.Formats.RSS {
				t.Fatalf("explicit false lost or absent defaults not inherited: %+v", feed.Formats)
			}
		})
	}
}

