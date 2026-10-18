package config

import (
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

// TestApplyEnvOverride_Characterization captures every key accepted by the
// legacy switch. Each assertion compares the complete config so an override
// cannot accidentally change an unrelated field.
func TestApplyEnvOverride_Characterization(t *testing.T) {
	type testCase struct {
		key, value string
		path       []string
		kind       string
	}
	intCases := []testCase{
		{"concurrency", "7", []string{"Concurrency"}, "int"},
		{"feed_defaults_items_per_page", "8", []string{"FeedDefaults", "ItemsPerPage"}, "int"},
		{"feeds_defaults_items_per_page", "9", []string{"FeedDefaults", "ItemsPerPage"}, "int"},
		{"feed_defaults_orphan_threshold", "10", []string{"FeedDefaults", "OrphanThreshold"}, "int"},
		{"feeds_defaults_orphan_threshold", "11", []string{"FeedDefaults", "OrphanThreshold"}, "int"},
		{"feed_defaults_syndication_max_items", "12", []string{"FeedDefaults", "Syndication", "MaxItems"}, "int"},
		{"feeds_defaults_syndication_max_items", "13", []string{"FeedDefaults", "Syndication", "MaxItems"}, "int"},
		{"feeds_defaults_items_per_page", "15", []string{"FeedDefaults", "ItemsPerPage"}, "int"},
		{"feeds_defaults_orphan_threshold", "16", []string{"FeedDefaults", "OrphanThreshold"}, "int"},
		{"encryption_min_password_length", "14", []string{"Encryption", "MinPasswordLength"}, "int"},
	}
	stringCases := []testCase{
		{"output_dir", "out", []string{"OutputDir"}, "string"}, {"url", "https://example.test", []string{"URL"}, "string"},
		{"title", "title", []string{"Title"}, "string"}, {"description", "description", []string{"Description"}, "string"},
		{"author", "author", []string{"Author"}, "string"}, {"language", "en", []string{"Language"}, "string"},
		{"author_url", "https://author.test", []string{"AuthorURL"}, "string"}, {"managing_editor", "editor", []string{"ManagingEditor"}, "string"},
		{"webmaster", "webmaster", []string{"WebMaster"}, "string"}, {"copyright", "copyright", []string{"Copyright"}, "string"},
		{"assets_dir", "assets", []string{"AssetsDir"}, "string"}, {"templates_dir", "templates", []string{"TemplatesDir"}, "string"},
		{"images_path", "images", []string{"Images", "Path"}, "string"}, {"images_template", "image.html", []string{"Images", "Template"}, "string"},
		{"glob_slug_mode", "path", []string{"GlobConfig", "SlugMode"}, "string"}, {"search_endpoint", "/search", []string{"Search", "Endpoint"}, "string"},
		{"search_backend", "bleve", []string{"Search", "Backend"}, "string"}, {"search_bleve_endpoint", ":8080", []string{"Search", "Bleve", "Endpoint"}, "string"},
		{"search_pagefind_cache_dir", "cache", []string{"Search", "Pagefind", "CacheDir"}, "string"}, {"search_pagefind_version", "1.2", []string{"Search", "Pagefind", "Version"}, "string"},
		{"search_pagefind_bundle_dir", "bundle", []string{"Search", "Pagefind", "BundleDir"}, "string"},
		{"encryption_default_key", "key", []string{"Encryption", "DefaultKey"}, "string"}, {"encryption_decryption_hint", "hint", []string{"Encryption", "DecryptionHint"}, "string"},
		{"encryption_min_estimated_crack_time", "centuries", []string{"Encryption", "MinEstimatedCrackTime"}, "string"},
		{"builder_admin_auth_headers_user_id", "x-user", []string{"BuilderAdmin", "Auth", "Headers", "UserID"}, "*string"},
		{"builder_admin_auth_headers_username", "x-name", []string{"BuilderAdmin", "Auth", "Headers", "Username"}, "*string"},
		{"builder_admin_auth_headers_display_name", "x-display", []string{"BuilderAdmin", "Auth", "Headers", "DisplayName"}, "*string"},
		{"builder_admin_auth_headers_email", "x-email", []string{"BuilderAdmin", "Auth", "Headers", "Email"}, "*string"},
		{"builder_admin_auth_headers_groups", "x-groups", []string{"BuilderAdmin", "Auth", "Headers", "Groups"}, "*string"},
		{"builder_admin_auth_headers_roles", "x-roles", []string{"BuilderAdmin", "Auth", "Headers", "Roles"}, "*string"},
		{"builder_admin_auth_headers_scopes", "x-scopes", []string{"BuilderAdmin", "Auth", "Headers", "Scopes"}, "*string"},
		{"builder_admin_webhook_branch", "main", []string{"BuilderAdmin", "Webhook", "Branch"}, "*string"},
		{"builder_admin_webhook_secret", "secret", []string{"BuilderAdmin", "Webhook", "Secret"}, "*string"},
	}
	listCases := []testCase{
		{"hooks", "a, b,, c ", []string{"Hooks"}, "list"}, {"disabled_hooks", "x, y", []string{"DisabledHooks"}, "list"},
		{"glob_patterns", "**/*.md, posts/*.md", []string{"GlobConfig", "Patterns"}, "list"},
		{"markdown_extensions", "tables, footnotes", []string{"MarkdownConfig", "Extensions"}, "list"},
	}
	boolCases := []testCase{
		{"images_enabled", "yes", []string{"Images", "Enabled"}, "*bool"}, {"images_export_json", "1", []string{"Images", "ExportJSON"}, "*bool"},
		{"images_include_unreferenced", "false", []string{"Images", "IncludeUnreferenced"}, "*bool"},
		{"glob_use_gitignore", "TRUE", []string{"GlobConfig", "UseGitignore"}, "bool"},
		{"feed_defaults_formats_html", "true", []string{"FeedDefaults", "Formats", "HTML"}, "bool"},
		{"feed_defaults_formats_rss", "true", []string{"FeedDefaults", "Formats", "RSS"}, "bool"},
		{"feeds_defaults_formats_rss", "no", []string{"FeedDefaults", "Formats", "RSS"}, "bool"},
		{"feeds_defaults_formats_html", "false", []string{"FeedDefaults", "Formats", "HTML"}, "bool"},
		{"feeds_defaults_formats_atom", "no", []string{"FeedDefaults", "Formats", "Atom"}, "bool"},
		{"feeds_defaults_formats_json", "yes", []string{"FeedDefaults", "Formats", "JSON"}, "bool"},
		{"feeds_defaults_formats_markdown", "no", []string{"FeedDefaults", "Formats", "Markdown"}, "bool"},
		{"feeds_defaults_formats_text", "yes", []string{"FeedDefaults", "Formats", "Text"}, "bool"},
		{"feeds_defaults_formats_sitemap", "no", []string{"FeedDefaults", "Formats", "Sitemap"}, "bool"},
		{"feed_defaults_formats_atom", "1", []string{"FeedDefaults", "Formats", "Atom"}, "bool"},
		{"feed_defaults_formats_json", "false", []string{"FeedDefaults", "Formats", "JSON"}, "bool"},
		{"feed_defaults_formats_markdown", "yes", []string{"FeedDefaults", "Formats", "Markdown"}, "bool"},
		{"feed_defaults_formats_text", "0", []string{"FeedDefaults", "Formats", "Text"}, "bool"},
		{"feed_defaults_formats_sitemap", "TRUE", []string{"FeedDefaults", "Formats", "Sitemap"}, "bool"},
		{"feed_defaults_syndication_include_content", "yes", []string{"FeedDefaults", "Syndication", "IncludeContent"}, "bool"},
		{"feed_defaults_syndication_site_archive_disabled", "true", []string{"FeedDefaults", "Syndication", "SiteArchiveDisabled"}, "bool"},
		{"feed_defaults_syndication_feed_archives_disabled", "false", []string{"FeedDefaults", "Syndication", "FeedArchivesDisabled"}, "bool"},
		{"feeds_defaults_syndication_include_content", "no", []string{"FeedDefaults", "Syndication", "IncludeContent"}, "bool"},
		{"feeds_defaults_syndication_site_archive_disabled", "no", []string{"FeedDefaults", "Syndication", "SiteArchiveDisabled"}, "bool"},
		{"feeds_defaults_syndication_feed_archives_disabled", "yes", []string{"FeedDefaults", "Syndication", "FeedArchivesDisabled"}, "bool"},
		{"search_pagefind_auto_install", "true", []string{"Search", "Pagefind", "AutoInstall"}, "*bool"},
		{"search_pagefind_verbose", "no", []string{"Search", "Pagefind", "Verbose"}, "*bool"},
		{"search_enabled", "yes", []string{"Search", "Enabled"}, "*bool"},
		{"search_pagefind_verbose", " TRUE ", []string{"Search", "Pagefind", "Verbose"}, "*bool"},
		{"encryption_enabled", "true", []string{"Encryption", "Enabled"}, "bool"},
		{"encryption_enforce_strength", "no", []string{"Encryption", "EnforceStrength"}, "bool"},
		{"blogroll_enabled", "true", []string{"Blogroll", "Enabled"}, "bool"},
		{"blogroll_refresh_on_build", "false", []string{"Blogroll", "RefreshOnBuild"}, "*bool"},
		{"builder_admin_webhook_enabled", "yes", []string{"BuilderAdmin", "Webhook", "Enabled"}, "*bool"},
	}
	for _, tc := range append(append(append(intCases, stringCases...), listCases...), boolCases...) {
		t.Run(tc.key, func(t *testing.T) {
			got := characterizedConfig()
			before := cloneCharacterizationConfig(got)
			applyEnvOverride(got, tc.key, tc.value)
			want := cloneCharacterizationConfig(before)
			setCharacterizedField(t, reflect.ValueOf(want).Elem(), tc.path, tc.kind, tc.value)
			if !reflect.DeepEqual(*got, *want) {
				t.Fatalf("config after %s=%q differs from complete expected value\ngot:  %#v\nwant: %#v", tc.key, tc.value, *got, *want)
			}
		})
	}

	for _, tc := range []struct{ key, value string }{
		{"content_index_enabled", "yes"}, {"contentindex_enabled", "no"}, {"content_index_output", "index.json"}, {"contentindex_output", "other.json"},
	} {
		t.Run(tc.key, func(t *testing.T) {
			got := &models.Config{Extra: map[string]interface{}{
				"content_index": map[string]interface{}{"keep": "yes"},
				"unrelated":     "preserved",
			}}
			before := cloneCharacterizationConfig(got)
			applyEnvOverride(got, tc.key, tc.value)
			want := before
			want.Extra = map[string]interface{}{
				"content_index": map[string]interface{}{"keep": "yes"},
				"unrelated":     "preserved",
			}
			content := map[string]interface{}{"keep": "yes"}
			if strings.HasSuffix(tc.key, "enabled") {
				content["enabled"] = parseBool(tc.value)
			} else {
				content["output"] = tc.value
			}
			want.Extra["content_index"] = content
			if !reflect.DeepEqual(*got, *want) {
				t.Fatalf("got %#v, want %#v", *got, *want)
			}
		})
	}

	for _, tc := range []struct {
		key, value string
		field      string
	}{{"tailwind_build", "true", "Build"}, {"tailwind_preflight", "no", "Preflight"}} {
		t.Run(tc.key, func(t *testing.T) {
			build, preflight := false, true
			existing := models.TailwindConfig{Input: "input.css", Output: "output.css", Build: &build, Preflight: &preflight}
			got := &models.Config{Extra: map[string]interface{}{"sentinel": "preserved", "tailwind": existing}}
			applyEnvOverride(got, tc.key, tc.value)
			want := cloneCharacterizationConfig(got)
			wantTailwind := existing
			value := parseBool(tc.value)
			if tc.field == "Build" {
				wantTailwind.Build = &value
			} else {
				wantTailwind.Preflight = &value
			}
			want.Extra = map[string]interface{}{"sentinel": "preserved", "tailwind": wantTailwind}
			applyEnvOverride(got, tc.key, tc.value)
			if !reflect.DeepEqual(*got, *want) {
				t.Fatalf("got %#v, want %#v", *got, *want)
			}
		})
	}

	for _, key := range []string{
		"concurrency", "feed_defaults_items_per_page", "feeds_defaults_items_per_page",
		"feed_defaults_orphan_threshold", "feeds_defaults_orphan_threshold",
		"feed_defaults_syndication_max_items", "feeds_defaults_syndication_max_items", "encryption_min_password_length",
	} {
		t.Run(key+"_invalid_integer", func(t *testing.T) {
			got := characterizedConfig()
			got.Concurrency = 31
			got.FeedDefaults.Syndication.MaxItems = 31
			got.Encryption.MinPasswordLength = 31
			before := *got
			applyEnvOverride(got, key, "not-an-integer")
			if !reflect.DeepEqual(*got, before) {
				t.Fatalf("invalid integer mutated config: got %#v, want %#v", *got, before)
			}
		})
	}

	t.Run("case_insensitive_unknown_and_empty_list", func(t *testing.T) {
		got := &models.Config{}
		applyEnvOverride(got, "TiTlE", "mixed")
		applyEnvOverride(got, "unknown_key", "ignored")
		applyEnvOverride(got, "hooks", " ,  , ")
		want := &models.Config{Title: "mixed", Hooks: parseStringList(" ,  , ")}
		if !reflect.DeepEqual(*got, *want) {
			t.Fatalf("got %#v, want %#v", *got, *want)
		}
	})

	t.Run("invalid_boolean_becomes_false", func(t *testing.T) {
		got := characterizedConfig()
		applyEnvOverride(got, "search_enabled", "sometimes")
		want := *characterizedConfig()
		falseValue := false
		want.Search.Enabled = &falseValue
		if !reflect.DeepEqual(*got, want) {
			t.Fatalf("got %#v, want %#v", *got, want)
		}
	})
}

func characterizedConfig() *models.Config {
	falseValue := false
	return &models.Config{
		OutputDir: "existing-output", URL: "https://old.test", Title: "existing-title",
		Hooks: []string{"existing-hook"}, Concurrency: 23,
		GlobConfig:     models.GlobConfig{Patterns: []string{"existing/**/*.md"}, SlugMode: "flat"},
		MarkdownConfig: models.MarkdownConfig{Extensions: []string{"existing-extension"}},
		FeedDefaults:   models.FeedDefaults{ItemsPerPage: 17, Formats: models.FeedFormats{HTML: true, RSS: true}},
		Search:         models.SearchConfig{Endpoint: "/existing", Enabled: &falseValue},
		Images:         models.ImagesConfig{Path: "existing-images"},
		Extra:          map[string]interface{}{"sentinel": map[string]interface{}{"keep": "yes"}},
	}
}

func cloneCharacterizationConfig(config *models.Config) *models.Config {
	cloned := cloneCharacterizationValue(reflect.ValueOf(config))
	configCopy, ok := cloned.Interface().(*models.Config)
	if !ok {
		panic("characterization clone did not return *models.Config")
	}
	return configCopy
}

func cloneCharacterizationValue(value reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}
	switch value.Kind() {
	case reflect.Interface:
		clonedValue := reflect.New(value.Type()).Elem()
		if !value.IsNil() {
			clonedValue.Set(cloneCharacterizationValue(value.Elem()))
		}
		return clonedValue
	case reflect.Pointer:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		clonedValue := reflect.New(value.Type().Elem())
		clonedValue.Elem().Set(cloneCharacterizationValue(value.Elem()))
		return clonedValue
	case reflect.Map:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		clonedValue := reflect.MakeMapWithSize(value.Type(), value.Len())
		iter := value.MapRange()
		for iter.Next() {
			clonedValue.SetMapIndex(cloneCharacterizationValue(iter.Key()), cloneCharacterizationValue(iter.Value()))
		}
		return clonedValue
	case reflect.Slice:
		if value.IsNil() {
			return reflect.Zero(value.Type())
		}
		clonedValue := reflect.MakeSlice(value.Type(), value.Len(), value.Len())
		for i := 0; i < value.Len(); i++ {
			clonedValue.Index(i).Set(cloneCharacterizationValue(value.Index(i)))
		}
		return clonedValue
	case reflect.Struct:
		clonedValue := reflect.New(value.Type()).Elem()
		clonedValue.Set(value)
		for i := 0; i < value.NumField(); i++ {
			if clonedValue.Field(i).CanSet() && value.Field(i).CanInterface() {
				clonedValue.Field(i).Set(cloneCharacterizationValue(value.Field(i)))
			}
		}
		return clonedValue
	default:
		return value
	}
}

func setCharacterizedField(t *testing.T, target reflect.Value, path []string, kind, raw string) {
	t.Helper()
	for _, part := range path[:len(path)-1] {
		target = target.FieldByName(part)
		if target.Kind() == reflect.Pointer {
			target.Set(reflect.New(target.Type().Elem()))
			target = target.Elem()
		}
	}
	target = target.FieldByName(path[len(path)-1])
	switch kind {
	case "string":
		target.SetString(raw)
	case "*string":
		s := raw
		target.Set(reflect.ValueOf(&s))
	case "bool":
		target.SetBool(parseBool(raw))
	case "*bool":
		b := parseBool(raw)
		target.Set(reflect.ValueOf(&b))
	case "int":
		n, err := strconv.Atoi(raw)
		if err != nil {
			t.Fatal(err)
		}
		target.SetInt(int64(n))
	case "list":
		target.Set(reflect.ValueOf(parseStringList(raw)))
	default:
		t.Fatalf("unknown characterization kind %q", kind)
	}
}
