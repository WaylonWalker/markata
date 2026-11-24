package config

import (
	"bufio"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

const (
	envPrefix          = "MARKATA_GO_" // Historical prefix, retained for old sites.
	canonicalEnvPrefix = "MARKATA_"
)

// Common string constants used in environment variable processing.
const (
	envKeyURL         = "url"
	envKeyConcurrency = "concurrency"
	envKeyJSON        = "json"
)

// ApplyEnvOverrides applies MARKATA_* configuration overrides with
// MARKATA_GO_* as a backward-compatible alias. When both prefixes set a key,
// the canonical value wins and a value-free migration warning is recorded.
// Nested keys use underscores: MARKATA_FEEDS_DEFAULTS_ITEMS_PER_PAGE.
// Boolean values: "true", "1", "yes" -> true; "false", "0", "no" -> false.
// List values: comma-separated strings.
func ApplyEnvOverrides(config *models.Config) error {
	legacy := make(map[string]string)
	canonical := make(map[string]string)
	for _, entry := range os.Environ() {
		name, value, ok := strings.Cut(entry, "=")
		if !ok {
			continue
		}
		switch {
		case strings.HasPrefix(name, envPrefix):
			legacy[strings.TrimPrefix(name, envPrefix)] = value
		case strings.HasPrefix(name, canonicalEnvPrefix):
			canonical[strings.TrimPrefix(name, canonicalEnvPrefix)] = value
		}
	}

	// Copy legacy values first, then canonical values. Never depend on the
	// (unspecified) iteration order of os.Environ or the maps for precedence.
	overrides := make(map[string]string, len(legacy)+len(canonical))
	for key, value := range legacy {
		overrides[key] = value
	}
	var warnings []string
	for key, value := range canonical {
		if old, exists := legacy[key]; exists && old != value {
			// Do not log values: configuration can contain secrets.
			warnings = append(warnings, canonicalEnvPrefix+key+" overrides "+envPrefix+key)
		}
		overrides[key] = value
	}
	if len(warnings) > 0 {
		sort.Strings(warnings)
		config.Extra = ensureExtra(config.Extra)
		config.Extra["name_migration_warnings"] = warnings
	}

	// Apply simple overrides. Every recognized key uses existing parsing.
	for key, value := range overrides {
		applyEnvOverride(config, key, value)
	}
	return nil
}

// applyEnvOverride applies a single environment variable override. The related
// setters keep each config area together, which makes supported keys easier to
// inspect without changing the permissive behavior of unknown keys.
func applyEnvOverride(config *models.Config, key, value string) {
	key = strings.ToLower(key)
	if applySiteEnvOverride(config, key, value) ||
		applyInputEnvOverride(config, key, value) ||
		applyFeedEnvOverride(config, key, value) ||
		applySearchEnvOverride(config, key, value) ||
		applyEncryptionEnvOverride(config, key, value) ||
		applyBlogrollEnvOverride(config, key, value) {
		return
	}
	applyBuilderAdminEnvOverride(config, key, value)
}

func applySiteEnvOverride(config *models.Config, key, value string) bool {
	switch key {
	case "output_dir":
		config.OutputDir = value
	case envKeyURL:
		config.URL = value
	case "title":
		config.Title = value
	case "description":
		config.Description = value
	case "author":
		config.Author = value
	case "language":
		config.Language = value
	case "author_url":
		config.AuthorURL = value
	case "managing_editor":
		config.ManagingEditor = value
	case "webmaster":
		config.WebMaster = value
	case "copyright":
		config.Copyright = value
	case "assets_dir":
		config.AssetsDir = value
	case "templates_dir":
		config.TemplatesDir = value
	case envKeyConcurrency:
		if n, err := strconv.Atoi(value); err == nil {
			config.Concurrency = n
		}
	case "hooks":
		config.Hooks = parseStringList(value)
	case "disabled_hooks":
		config.DisabledHooks = parseStringList(value)
	default:
		return false
	}
	return true
}

func applyInputEnvOverride(config *models.Config, key, value string) bool {
	switch key {
	case "images_enabled":
		config.Images.Enabled = boolPointer(value)
	case "images_path":
		config.Images.Path = value
	case "images_template":
		config.Images.Template = value
	case "images_export_json":
		config.Images.ExportJSON = boolPointer(value)
	case "images_include_unreferenced":
		config.Images.IncludeUnreferenced = boolPointer(value)
	case "glob_patterns":
		config.GlobConfig.Patterns = parseStringList(value)
	case "glob_use_gitignore":
		config.GlobConfig.UseGitignore = parseBool(value)
	case "glob_slug_mode":
		config.GlobConfig.SlugMode = value
	case "markdown_extensions":
		config.MarkdownConfig.Extensions = parseStringList(value)
	default:
		return false
	}
	return true
}

func applyFeedEnvOverride(config *models.Config, key, value string) bool {
	switch key {
	case "feed_defaults_items_per_page", "feeds_defaults_items_per_page":
		if n, err := strconv.Atoi(value); err == nil {
			config.FeedDefaults.ItemsPerPage = n
		}
	case "feed_defaults_orphan_threshold", "feeds_defaults_orphan_threshold":
		if n, err := strconv.Atoi(value); err == nil {
			config.FeedDefaults.OrphanThreshold = n
		}
	case "feed_defaults_formats_html", "feeds_defaults_formats_html":
		config.FeedDefaults.Formats.HTML = parseBool(value)
	case "feed_defaults_formats_rss", "feeds_defaults_formats_rss":
		config.FeedDefaults.Formats.RSS = parseBool(value)
	case "feed_defaults_formats_atom", "feeds_defaults_formats_atom":
		config.FeedDefaults.Formats.Atom = parseBool(value)
	case "feed_defaults_formats_json", "feeds_defaults_formats_json":
		config.FeedDefaults.Formats.JSON = parseBool(value)
	case "feed_defaults_formats_markdown", "feeds_defaults_formats_markdown":
		config.FeedDefaults.Formats.Markdown = parseBool(value)
	case "feed_defaults_formats_text", "feeds_defaults_formats_text":
		config.FeedDefaults.Formats.Text = parseBool(value)
	case "feed_defaults_formats_sitemap", "feeds_defaults_formats_sitemap":
		config.FeedDefaults.Formats.Sitemap = parseBool(value)
	case "feed_defaults_syndication_max_items", "feeds_defaults_syndication_max_items":
		if n, err := strconv.Atoi(value); err == nil {
			config.FeedDefaults.Syndication.MaxItems = n
		}
	case "feed_defaults_syndication_include_content", "feeds_defaults_syndication_include_content":
		config.FeedDefaults.Syndication.IncludeContent = parseBool(value)
	case "feed_defaults_syndication_site_archive_disabled", "feeds_defaults_syndication_site_archive_disabled":
		config.FeedDefaults.Syndication.SiteArchiveDisabled = parseBool(value)
	case "feed_defaults_syndication_feed_archives_disabled", "feeds_defaults_syndication_feed_archives_disabled":
		config.FeedDefaults.Syndication.FeedArchivesDisabled = parseBool(value)
	default:
		return false
	}
	return true
}

func applySearchEnvOverride(config *models.Config, key, value string) bool {
	switch key {
	case "search_endpoint":
		config.Search.Endpoint = value
	case "search_backend":
		config.Search.Backend = value
	case "search_bleve_endpoint":
		config.Search.Bleve.Endpoint = value
	case "search_pagefind_auto_install":
		config.Search.Pagefind.AutoInstall = boolPointer(value)
	case "search_pagefind_cache_dir":
		config.Search.Pagefind.CacheDir = value
	case "search_pagefind_version":
		config.Search.Pagefind.Version = value
	case "search_pagefind_bundle_dir":
		config.Search.Pagefind.BundleDir = value
	case "search_pagefind_verbose":
		config.Search.Pagefind.Verbose = boolPointer(value)
	case "search_enabled":
		config.Search.Enabled = boolPointer(value)
	case "content_index_enabled", "contentindex_enabled", "content_index_output", "contentindex_output":
		applyContentIndexEnvOverride(config, key, value)
	case "tailwind_build":
		config.Extra = ensureExtra(config.Extra)
		tw, ok := config.Extra["tailwind"].(models.TailwindConfig)
		if !ok {
			tw = models.NewTailwindConfig()
		}
		tw.Build = boolPointer(value)
		config.Extra["tailwind"] = tw
	case "tailwind_preflight":
		config.Extra = ensureExtra(config.Extra)
		tw, ok := config.Extra["tailwind"].(models.TailwindConfig)
		if !ok {
			tw = models.NewTailwindConfig()
		}
		tw.Preflight = boolPointer(value)
		config.Extra["tailwind"] = tw
	default:
		return false
	}
	return true
}

func applyContentIndexEnvOverride(config *models.Config, key, value string) {
	config.Extra = ensureExtra(config.Extra)
	contentIndex, contentIndexOK := config.Extra["content_index"].(map[string]interface{})
	if !contentIndexOK {
		contentIndex = nil
	}
	if contentIndex == nil {
		contentIndex = make(map[string]interface{})
	}
	if strings.HasSuffix(key, "enabled") {
		contentIndex["enabled"] = parseBool(value)
	} else {
		contentIndex["output"] = value
	}
	config.Extra["content_index"] = contentIndex
}

func ensureExtra(extra map[string]interface{}) map[string]interface{} {
	if extra == nil {
		return make(map[string]interface{})
	}
	return extra
}

func applyEncryptionEnvOverride(config *models.Config, key, value string) bool {
	switch key {
	case "encryption_enabled":
		config.Encryption.Enabled = parseBool(value)
	case "encryption_default_key":
		config.Encryption.DefaultKey = value
	case "encryption_decryption_hint":
		config.Encryption.DecryptionHint = value
	case "encryption_enforce_strength":
		config.Encryption.EnforceStrength = parseBool(value)
	case "encryption_min_estimated_crack_time":
		config.Encryption.MinEstimatedCrackTime = value
	case "encryption_min_password_length":
		if n, err := strconv.Atoi(value); err == nil {
			config.Encryption.MinPasswordLength = n
		}
	default:
		return false
	}
	return true
}

func applyBlogrollEnvOverride(config *models.Config, key, value string) bool {
	switch key {
	case "blogroll_enabled":
		config.Blogroll.Enabled = parseBool(value)
	case "blogroll_refresh_on_build":
		config.Blogroll.RefreshOnBuild = boolPointer(value)
	default:
		return false
	}
	return true
}

func applyBuilderAdminEnvOverride(config *models.Config, key, value string) {
	valuePointer := &value
	switch key {
	case "builder_admin_auth_headers_user_id":
		config.BuilderAdmin.Auth.Headers.UserID = valuePointer
	case "builder_admin_auth_headers_username":
		config.BuilderAdmin.Auth.Headers.Username = valuePointer
	case "builder_admin_auth_headers_display_name":
		config.BuilderAdmin.Auth.Headers.DisplayName = valuePointer
	case "builder_admin_auth_headers_email":
		config.BuilderAdmin.Auth.Headers.Email = valuePointer
	case "builder_admin_auth_headers_groups":
		config.BuilderAdmin.Auth.Headers.Groups = valuePointer
	case "builder_admin_auth_headers_roles":
		config.BuilderAdmin.Auth.Headers.Roles = valuePointer
	case "builder_admin_auth_headers_scopes":
		config.BuilderAdmin.Auth.Headers.Scopes = valuePointer
	case "builder_admin_auth_headers_auth_method":
		config.BuilderAdmin.Auth.Headers.AuthMethod = valuePointer
	case "builder_admin_webhook_enabled":
		config.BuilderAdmin.Webhook.Enabled = boolPointer(value)
	case "builder_admin_webhook_branch":
		config.BuilderAdmin.Webhook.Branch = valuePointer
	case "builder_admin_webhook_secret":
		config.BuilderAdmin.Webhook.Secret = valuePointer
	}
}

func boolPointer(value string) *bool {
	parsed := parseBool(value)
	return &parsed
}

// parseBool parses a string into a boolean.
// "true", "1", "yes" -> true
// "false", "0", "no" -> false
// All comparisons are case-insensitive.
func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	}
	return false
}

// parseStringList parses a comma-separated string into a slice.
func parseStringList(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// GetEnvValue prefers a canonical MARKATA_* variable over its legacy
// MARKATA_GO_* alias, retaining the old lookup as a fallback.
func GetEnvValue(key string) (string, bool) {
	key = strings.ToUpper(key)
	if value, ok := os.LookupEnv(canonicalEnvPrefix + key); ok {
		return value, true
	}
	return os.LookupEnv(envPrefix + key)
}

// SetEnvValue sets an environment variable with the MARKATA_GO_ prefix.
// This is primarily useful for testing.
func SetEnvValue(key, value string) error {
	return os.Setenv(envPrefix+strings.ToUpper(key), value)
}

// UnsetEnvValue unsets an environment variable with the MARKATA_GO_ prefix.
// This is primarily useful for testing.
func UnsetEnvValue(key string) error {
	return os.Unsetenv(envPrefix + strings.ToUpper(key))
}

// FromEnv creates a Config entirely from environment variables.
// This is useful when no config file is available.
func FromEnv() *models.Config {
	config := DefaultConfig()
	_ = ApplyEnvOverrides(config) //nolint:errcheck // Best-effort env override
	return config
}

// StructToEnvKeys returns a map of environment variable keys for a struct.
// This is useful for documentation and debugging.
func StructToEnvKeys(prefix string, v interface{}) map[string]string {
	result := make(map[string]string)
	structToEnvKeysRecursive(prefix, reflect.TypeOf(v), result)
	return result
}

func structToEnvKeysRecursive(prefix string, t reflect.Type, result map[string]string) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Get the field name for the environment variable
		name := field.Name
		if tag := field.Tag.Get("json"); tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] != "" && parts[0] != "-" {
				name = parts[0]
			}
		}

		envKey := prefix + strings.ToUpper(name)

		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		switch fieldType.Kind() {
		case reflect.Struct:
			structToEnvKeysRecursive(envKey+"_", fieldType, result)
		case reflect.Slice:
			result[envKey] = "comma-separated list"
		case reflect.Bool:
			result[envKey] = "true/false/1/0/yes/no"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result[envKey] = "integer"
		case reflect.String:
			result[envKey] = "string"
		default:
			// Other types (uint, float, complex, etc.) are not currently supported for env vars
		}
	}
}

// LoadDotEnv loads environment variables from a .env file in the current directory.
// Lines starting with '#' are treated as comments. Empty lines are skipped.
// Values can optionally be quoted with single or double quotes.
// Variables are set into the process environment so they are available to
// os.Environ() and os.Getenv() calls (including encryption key lookups).
//
// This function is safe to call even if no .env file exists - it will silently
// return nil in that case.
func LoadDotEnv() error {
	return LoadDotEnvFile(".env")
}

// LoadDotEnvFile loads environment variables from the specified file path.
// Returns nil if the file does not exist.
func LoadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No .env file is fine
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Strip surrounding quotes (single or double)
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		// Only set if not already set in the environment
		// (real env vars take precedence over .env)
		if _, exists := os.LookupEnv(key); !exists {
			if err := os.Setenv(key, value); err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}
