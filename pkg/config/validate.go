package config

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/fontpacks"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/renderingcontract"
	"github.com/WaylonWalker/markata/pkg/suggest"
)

// ValidationError represents a configuration validation error.
// This type is kept for backward compatibility.
type ValidationError struct {
	Field   string
	Message string
	IsWarn  bool // If true, this is a warning rather than an error
}

func (e ValidationError) Error() string {
	prefix := "error"
	if e.IsWarn {
		prefix = "warning"
	}
	return fmt.Sprintf("config %s: %s: %s", prefix, e.Field, e.Message)
}

// ValidateConfig validates a configuration and returns any errors or warnings.
// Errors are returned first, followed by warnings.
// The returned slice will be nil if no issues are found.
func ValidateConfig(config *models.Config) []error {
	if config == nil {
		return []error{errors.New("config is nil")}
	}

	var errs []error
	errs = append(errs, validateRenderingTheme(config)...)

	// Validate URL format if provided
	if config.URL != "" {
		if err := validateURL(config.URL); err != nil {
			errs = append(errs, ValidationError{
				Field:   "url",
				Message: err.Error(),
			})
		}
	}

	// Validate concurrency
	if config.Concurrency < 0 {
		errs = append(errs, ValidationError{
			Field:   "concurrency",
			Message: "must be >= 0 (0 means auto-detect)",
		})
	}

	errs = append(errs, validateImageLibraryConfig(config)...)

	// Warn on empty glob patterns
	if len(config.GlobConfig.Patterns) == 0 {
		errs = append(errs, ValidationError{
			Field:   "glob.patterns",
			Message: "no glob patterns specified, no files will be processed",
			IsWarn:  true,
		})
	}

	if config.GlobConfig.SlugMode != "" && !models.IsValidSlugMode(config.GlobConfig.SlugMode) {
		errs = append(errs, ValidationError{
			Field:   "glob.slug_mode",
			Message: `must be one of: "flat", "path"`,
		})
	}
	for i, rule := range config.GlobConfig.SlugRules {
		if strings.TrimSpace(rule.Prefix) == "" {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("glob.slug_rules[%d].prefix", i),
				Message: "must not be empty",
			})
		}
		if !models.IsValidSlugMode(rule.Mode) {
			errs = append(errs, ValidationError{
				Field:   fmt.Sprintf("glob.slug_rules[%d].mode", i),
				Message: `must be one of: "flat", "path"`,
			})
		}
	}

	// Validate license configuration
	if !config.License.IsDisabled() {
		if key, ok := config.License.Key(); ok {
			if _, valid := models.GetLicenseOption(key); !valid {
				supported := strings.Join(models.LicenseKeys(), ", ")
				errs = append(errs, ValidationError{
					Field:   "license",
					Message: fmt.Sprintf("unsupported license key %q (supported: %s)", key, supported),
				})
			}
		} else if !config.License.HasValue() {
			supported := strings.Join(models.LicenseKeys(), ", ")
			errs = append(errs, ValidationError{
				Field:   "license",
				Message: fmt.Sprintf("license not configured (supported: %s). Set license = %q or license = false to skip the footer.", supported, models.DefaultLicenseKey),
				IsWarn:  true,
			})
		}
	}

	// Validate feed configurations
	// Validate raw view declarations before applying defaults so an explicit
	// empty list remains distinguishable from an omitted list.
	for i := range config.Feeds {
		errs = append(errs, validateFeedViews(fmt.Sprintf("feeds[%d].views", i), config.Feeds[i].Views)...)
		feedWithDefaults := config.Feeds[i]
		feedWithDefaults.ApplyDefaults(config.FeedDefaults)
		feedErrs := validateFeedConfig(i, &feedWithDefaults)
		errs = append(errs, feedErrs...)
	}

	// Validate feed defaults
	errs = append(errs, validateFeedViews("feed_defaults.views", config.FeedDefaults.Views)...)
	if config.FeedDefaults.ItemsPerPage < 0 {
		errs = append(errs, ValidationError{
			Field:   "feed_defaults.items_per_page",
			Message: "must be >= 0",
		})
	}
	if config.FeedDefaults.OrphanThreshold < 0 {
		errs = append(errs, ValidationError{
			Field:   "feed_defaults.orphan_threshold",
			Message: "must be >= 0",
		})
	}

	// Warn if Tailwind CSS is included but css_purge is disabled
	if tailwindIncludeMode(config) == CSS && !isCSSPurgeEnabled(config) {
		errs = append(errs, ValidationError{
			Field:   "tailwind.include",
			Message: "Tailwind CSS output is large; enable css_purge to reduce CSS size",
			IsWarn:  true,
		})
	}

	// Sort errors first, then warnings
	sortErrors(errs)

	return errs
}

func validateImageLibraryConfig(config *models.Config) []error {
	var errs []error
	if config.Images.PageSize < 0 {
		errs = append(errs, ValidationError{
			Field:   "images.page_size",
			Message: "must be a nonnegative integer (0 disables pagination)",
		})
	}
	if path := strings.TrimSpace(config.Images.Path); path != "" && imagePathEscapesRoot(path) {
		errs = append(errs, ValidationError{
			Field:   "images.path",
			Message: "must be a relative path within output_dir",
		})
	}
	if template := strings.TrimSpace(config.Images.Template); template != "" && imagePathEscapesRoot(template) {
		errs = append(errs, ValidationError{
			Field:   "images.template",
			Message: "must be a relative template path",
		})
	}
	return errs
}

func imagePathEscapesRoot(path string) bool {
	clean := filepath.ToSlash(filepath.Clean(path))
	return filepath.IsAbs(path) || strings.HasPrefix(path, "/") || clean == ".." || strings.HasPrefix(clean, "../")
}

//nolint:gocyclo // Each contract dimension produces an independent validation diagnostic.
func validateRenderingTheme(config *models.Config) []error {
	c, err := renderingcontract.Load()
	if err != nil {
		return []error{ValidationError{Field: "theme", Message: "cannot load rendering contract: " + err.Error()}}
	}
	var errs []error
	valid := func(group, value, field string) {
		if value == "" {
			return
		}
		for _, candidate := range c.Enums[group] {
			if value == candidate {
				return
			}
		}
		message := fmt.Sprintf("unsupported value %q; choose one of: %s", value, strings.Join(c.Enums[group], ", "))
		if matches := suggest.Closest(value, c.Enums[group], 1); len(matches) > 0 {
			message += fmt.Sprintf("; did you mean %q?", matches[0])
		}
		errs = append(errs, ValidationError{Field: field, Message: message})
	}
	if palette := config.Theme.Palette; palette != "" && !knownPalette(&c, palette) {
		message := fmt.Sprintf("unknown palette %q; run 'markata-go palette list' to see available palettes", palette)
		candidates := make([]string, 0, len(c.Palettes))
		for _, item := range c.Palettes {
			candidates = append(candidates, item.ID)
		}
		if matches := suggest.Closest(palette, candidates, 1); len(matches) > 0 {
			message += fmt.Sprintf("; did you mean %q?", matches[0])
		}
		errs = append(errs, ValidationError{Field: "theme.palette", Message: message})
	}
	valid("aesthetics", config.Theme.Aesthetic, "theme.aesthetic")
	if config.Theme.TextSize != "" && !isValidTextSize(config.Theme.TextSize) {
		errs = append(errs, ValidationError{Field: "theme.text_size", Message: `must be one of: "small", "medium", "large", "x-large"; using "large"`, IsWarn: true})
	}
	if !fontpackInCatalog(config) {
		// A legacy top-level fontpack is copied into theme.fontpack during
		// normalization; report it under the key the user actually wrote.
		field := "theme.fontpack"
		if config.Fontpack != "" && config.Fontpack == config.Theme.Fontpack {
			field = "fontpack"
		}
		valid("fontpacks", config.Theme.Fontpack, field)
	}
	valid("textures", config.Theme.Texture.Kind, "theme.texture.kind")
	valid("scopes", config.Theme.Texture.Scope, "theme.texture.scope")
	valid("heading_textures", config.Theme.HeadingTexture.Kind, "theme.heading_texture.kind")
	valid("motifs", config.Theme.Motif.Kind, "theme.motif.kind")
	valid("motif_layers", config.Theme.Motif.Layer, "theme.motif.layer")
	valid("motif_colors", config.Theme.Motif.Color, "theme.motif.color")
	for field, value := range map[string]float64{"theme.texture.color_mix": config.Theme.Texture.ColorMix, "theme.heading_texture.color_mix": config.Theme.HeadingTexture.ColorMix, "theme.motif.color_mix": config.Theme.Motif.ColorMix, "theme.motif.row_offset": config.Theme.Motif.RowOffset, "theme.motif.wobble": config.Theme.Motif.Wobble, "theme.motif.scatter": config.Theme.Motif.Scatter} {
		if value < 0 || value > 1 {
			errs = append(errs, ValidationError{Field: field, Message: "must be between 0 and 1"})
		}
	}
	for field, value := range map[string]float64{"theme.texture.scale": config.Theme.Texture.Scale, "theme.heading_texture.scale": config.Theme.HeadingTexture.Scale} {
		if value != 0 && (value < .25 || value > 3) {
			errs = append(errs, ValidationError{Field: field, Message: "must be between 0.25 and 3"})
		}
	}
	if warnings, ok := config.Extra["theme_migration_warnings"].([]string); ok {
		for _, warning := range warnings {
			errs = append(errs, ValidationError{Field: "theme", Message: warning, IsWarn: true})
		}
	}
	return errs
}

func isValidTextSize(value string) bool {
	switch value {
	case models.TextSizeSmall, models.TextSizeMedium, models.TextSizeLarge, models.TextSizeXLarge:
		return true
	default:
		return false
	}
}

// ValidateConfigWithPositions validates a configuration with file position tracking.
// It returns a ConfigErrors collection with detailed error information including
// file paths, line numbers, and fix suggestions.
func ValidateConfigWithPositions(config *models.Config, tracker *PositionTracker) *ConfigErrors {
	configErrors := &ConfigErrors{}

	if config == nil {
		configErrors.Add(&ConfigError{
			Field:   "config",
			Message: "configuration is nil",
		})
		return configErrors
	}

	// Validate URL format if provided
	if config.URL != "" {
		validateURLWithDetails(config.URL, tracker, configErrors)
	}

	// Validate concurrency
	if config.Concurrency < 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"concurrency",
			fmt.Sprintf("%d", config.Concurrency),
			"must be >= 0 (0 means auto-detect)",
			GetFixSuggestion("negative_value", "concurrency", ""),
			false,
		))
	}

	// Warn on empty glob patterns
	if len(config.GlobConfig.Patterns) == 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"glob.patterns",
			"[]",
			"no glob patterns specified, no files will be processed",
			GetFixSuggestion("empty_patterns", "", ""),
			true,
		))
	}

	if config.GlobConfig.SlugMode != "" && !models.IsValidSlugMode(config.GlobConfig.SlugMode) {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"glob.slug_mode",
			config.GlobConfig.SlugMode,
			`must be one of: "flat", "path"`,
			"Use \"flat\" for filename-based slugs or \"path\" for directory-based slugs.",
			false,
		))
	}
	for i, rule := range config.GlobConfig.SlugRules {
		if strings.TrimSpace(rule.Prefix) == "" {
			configErrors.Add(NewConfigErrorWithFix(
				tracker,
				fmt.Sprintf("glob.slug_rules[%d].prefix", i),
				rule.Prefix,
				"must not be empty",
				"Set a content path prefix like \"posts/blog\" or remove the rule.",
				false,
			))
		}
		if !models.IsValidSlugMode(rule.Mode) {
			configErrors.Add(NewConfigErrorWithFix(
				tracker,
				fmt.Sprintf("glob.slug_rules[%d].mode", i),
				rule.Mode,
				`must be one of: "flat", "path"`,
				"Use \"flat\" or \"path\".",
				false,
			))
		}
	}

	// Validate feed configurations
	for i := range config.Feeds {
		validateFeedViewsWithPositions(fmt.Sprintf("feeds[%d].views", i), config.Feeds[i].Views, tracker, configErrors)
		feedWithDefaults := config.Feeds[i]
		feedWithDefaults.ApplyDefaults(config.FeedDefaults)
		validateFeedConfigWithPositions(i, &feedWithDefaults, tracker, configErrors)
	}

	// Validate feed defaults
	validateFeedViewsWithPositions("feed_defaults.views", config.FeedDefaults.Views, tracker, configErrors)
	if config.FeedDefaults.ItemsPerPage < 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"feed_defaults.items_per_page",
			fmt.Sprintf("%d", config.FeedDefaults.ItemsPerPage),
			"must be >= 0",
			GetFixSuggestion("negative_value", "items_per_page", ""),
			false,
		))
	}
	if config.FeedDefaults.OrphanThreshold < 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"feed_defaults.orphan_threshold",
			fmt.Sprintf("%d", config.FeedDefaults.OrphanThreshold),
			"must be >= 0",
			GetFixSuggestion("negative_value", "orphan_threshold", ""),
			false,
		))
	}

	// Warn if Tailwind CSS is included but css_purge is disabled
	if tailwindIncludeMode(config) == CSS && !isCSSPurgeEnabled(config) {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"tailwind.include",
			CSS,
			"Tailwind CSS output is large; enable css_purge to reduce CSS size",
			"Enable [markata-go.css_purge].enabled = true",
			true,
		))
	}

	return configErrors
}

// validateURLWithDetails validates a URL and adds detailed errors to configErrors.
func validateURLWithDetails(rawURL string, tracker *PositionTracker, configErrors *ConfigErrors) {
	u, err := url.Parse(rawURL)
	if err != nil {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"url",
			rawURL,
			fmt.Sprintf("invalid URL format: %v", err),
			GetFixSuggestion("url_no_scheme", "url", rawURL),
			false,
		))
		return
	}

	if u.Scheme == "" {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"url",
			rawURL,
			"URL must include a scheme (e.g., https://)",
			GetFixSuggestion("url_no_scheme", "url", rawURL),
			false,
		))
		return
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"url",
			rawURL,
			fmt.Sprintf("URL scheme must be http or https, got %q", u.Scheme),
			GetFixSuggestion("url_invalid_scheme", "url", rawURL),
			false,
		))
		return
	}

	if u.Host == "" {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			"url",
			rawURL,
			"URL must include a host",
			GetFixSuggestion("url_no_host", "url", rawURL),
			false,
		))
	}
}

// validateFeedConfigWithPositions validates a feed configuration with position tracking.
func validateFeedConfigWithPositions(index int, feed *models.FeedConfig, tracker *PositionTracker, configErrors *ConfigErrors) {
	prefix := fmt.Sprintf("feeds[%d]", index)

	// Validate items_per_page
	if feed.ItemsPerPage < 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			prefix+".items_per_page",
			fmt.Sprintf("%d", feed.ItemsPerPage),
			"must be >= 0",
			GetFixSuggestion("negative_value", "items_per_page", ""),
			false,
		))
	}

	// Validate limit
	if feed.Limit < 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			prefix+".limit",
			fmt.Sprintf("%d", feed.Limit),
			"must be >= 0",
			GetFixSuggestion("negative_value", "limit", ""),
			false,
		))
	}

	// Validate offset
	if feed.Offset < 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			prefix+".offset",
			fmt.Sprintf("%d", feed.Offset),
			"must be >= 0",
			GetFixSuggestion("negative_value", "offset", ""),
			false,
		))
	}

	// Validate orphan_threshold
	if feed.OrphanThreshold < 0 {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			prefix+".orphan_threshold",
			fmt.Sprintf("%d", feed.OrphanThreshold),
			"must be >= 0",
			GetFixSuggestion("negative_value", "orphan_threshold", ""),
			false,
		))
	}

	// Warn if no output formats are enabled
	if !hasAnyFormat(feed.Formats) {
		configErrors.Add(NewConfigErrorWithFix(
			tracker,
			prefix+".formats",
			"{}",
			"no output formats enabled, feed will not produce any output",
			GetFixSuggestion("no_formats", "", ""),
			true,
		))
	}
}

// validateURL validates a URL string.
func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if u.Scheme == "" {
		return errors.New("URL must include a scheme (e.g., https://)")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL scheme must be http or https, got %q", u.Scheme)
	}

	if u.Host == "" {
		return errors.New("URL must include a host")
	}

	return nil
}

// validateFeedConfig validates a single feed configuration.
func validateFeedConfig(index int, feed *models.FeedConfig) []error {
	var errs []error
	prefix := fmt.Sprintf("feeds[%d]", index)

	// Note: Empty slug is allowed - it represents the home page feed (index.html)
	// No validation needed for slug field

	// Validate items_per_page
	if feed.ItemsPerPage < 0 {
		errs = append(errs, ValidationError{
			Field:   prefix + ".items_per_page",
			Message: "must be >= 0",
		})
	}

	// Validate limit
	if feed.Limit < 0 {
		errs = append(errs, ValidationError{
			Field:   prefix + ".limit",
			Message: "must be >= 0",
		})
	}

	// Validate offset
	if feed.Offset < 0 {
		errs = append(errs, ValidationError{
			Field:   prefix + ".offset",
			Message: "must be >= 0",
		})
	}

	// Validate orphan_threshold
	if feed.OrphanThreshold < 0 {
		errs = append(errs, ValidationError{
			Field:   prefix + ".orphan_threshold",
			Message: "must be >= 0",
		})
	}

	// Warn if no output formats are enabled
	if !hasAnyFormat(feed.Formats) {
		errs = append(errs, ValidationError{
			Field:   prefix + ".formats",
			Message: "no output formats enabled, feed will not produce any output",
			IsWarn:  true,
		})
	}

	return errs
}

// hasAnyFormat returns true if any feed format is enabled.
func hasAnyFormat(formats models.FeedFormats) bool {
	return formats.HTML || formats.SimpleHTML || formats.RSS || formats.Atom ||
		formats.JSON || formats.Markdown || formats.Text || formats.Sitemap
}

// sortErrors sorts validation errors so that errors come before warnings.
func sortErrors(errs []error) {
	if len(errs) <= 1 {
		return
	}

	// Simple bubble sort - stable and fine for small slices
	for i := 0; i < len(errs)-1; i++ {
		for j := 0; j < len(errs)-i-1; j++ {
			// If current is warning and next is error, swap
			currWarn := isWarning(errs[j])
			nextWarn := isWarning(errs[j+1])
			if currWarn && !nextWarn {
				errs[j], errs[j+1] = errs[j+1], errs[j]
			}
		}
	}
}

// isWarning returns true if the error is a ValidationError with IsWarn=true
// or a ConfigError with IsWarn=true.
func isWarning(err error) bool {
	var ve ValidationError
	if errors.As(err, &ve) {
		return ve.IsWarn
	}
	var ce *ConfigError
	if errors.As(err, &ce) {
		return ce.IsWarn
	}
	return false
}

// HasErrors returns true if the error slice contains any actual errors (not warnings).
func HasErrors(errs []error) bool {
	for _, err := range errs {
		if !isWarning(err) {
			return true
		}
	}
	return false
}

// HasWarnings returns true if the error slice contains any warnings.
func HasWarnings(errs []error) bool {
	for _, err := range errs {
		if isWarning(err) {
			return true
		}
	}
	return false
}

// SplitErrorsAndWarnings separates errors and warnings into separate slices.
func SplitErrorsAndWarnings(errs []error) (errsOut, warnsOut []error) {
	for _, err := range errs {
		if isWarning(err) {
			warnsOut = append(warnsOut, err)
		} else {
			errsOut = append(errsOut, err)
		}
	}
	return
}

// ValidateAndWarn validates config and logs warnings but only returns actual errors.
// This is useful when you want to proceed with warnings but stop on errors.
func ValidateAndWarn(config *models.Config, warnFunc func(error)) []error {
	allErrs := ValidateConfig(config)
	actualErrors, warnings := SplitErrorsAndWarnings(allErrs)

	for _, w := range warnings {
		if warnFunc != nil {
			warnFunc(w)
		}
	}

	return actualErrors
}

func tailwindIncludeMode(config *models.Config) string {
	if config == nil {
		return ""
	}

	if config.Extra == nil {
		return ""
	}

	raw, ok := config.Extra["tailwind"]
	if !ok {
		return ""
	}

	switch value := raw.(type) {
	case models.TailwindConfig:
		return normalizeTailwindInclude(value.IncludeMode())
	case *models.TailwindConfig:
		return normalizeTailwindInclude(value.IncludeMode())
	case map[string]interface{}:
		if include, ok := value["include"]; ok {
			if s, ok := include.(string); ok {
				return normalizeTailwindInclude(s)
			}
			if b, ok := include.(bool); ok {
				if !b {
					return ""
				}
				return CSS
			}
		}
		return ""
	default:
		return ""
	}
}

func isCSSPurgeEnabled(config *models.Config) bool {
	if config == nil {
		return false
	}

	if config.Extra == nil {
		return false
	}

	raw, ok := config.Extra["css_purge"]
	if !ok {
		return false
	}

	switch value := raw.(type) {
	case models.CSSPurgeConfig:
		return value.Enabled
	case *models.CSSPurgeConfig:
		return value.Enabled
	case map[string]interface{}:
		if enabled, ok := value["enabled"].(bool); ok {
			return enabled
		}
		return false
	default:
		return false
	}
}

func normalizeTailwindInclude(value string) string {
	return configNormalizeTailwindInclude(value)
}

func configNormalizeTailwindInclude(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	switch value {
	case "", "false", "off", "0", "none":
		return ""
	case CSS, "js":
		return value
	default:
		return ""
	}
}

// FormatConfigError formats a ConfigError for CLI display with colors and context.
func FormatConfigError(err *ConfigError) string {
	var sb strings.Builder

	// Header with location
	switch {
	case err.File != "" && err.Line > 0 && err.Column > 0:
		sb.WriteString(fmt.Sprintf("Error: configuration error in %s:%d:%d\n", err.File, err.Line, err.Column))
	case err.File != "" && err.Line > 0:
		sb.WriteString(fmt.Sprintf("Error: configuration error in %s:%d\n", err.File, err.Line))
	case err.File != "":
		sb.WriteString(fmt.Sprintf("Error: configuration error in %s\n", err.File))
	default:
		sb.WriteString("Error: configuration error\n")
	}

	// Context with arrow pointing to problematic line
	if len(err.Context) > 0 {
		sb.WriteString("\n")
		for _, line := range err.Context {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}

	// Error details
	sb.WriteString("\n")
	sb.WriteString(err.Field)
	sb.WriteString(": ")
	sb.WriteString(err.Message)
	sb.WriteString("\n")

	// Fix suggestion
	if err.Fix != "" {
		sb.WriteString("\nFix: ")
		sb.WriteString(err.Fix)
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatConfigErrors formats multiple ConfigErrors for CLI display.
func FormatConfigErrors(configErrors *ConfigErrors) string {
	if configErrors == nil || len(configErrors.Errors) == 0 {
		return ""
	}

	if len(configErrors.Errors) == 1 {
		return FormatConfigError(configErrors.Errors[0])
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Error: configuration validation failed with %d issues\n", len(configErrors.Errors)))

	for i, err := range configErrors.Errors {
		if i > 0 {
			sb.WriteString("\n")
			sb.WriteString(strings.Repeat("-", 60))
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
		sb.WriteString(FormatConfigError(err))
	}

	return sb.String()
}

// fontpackInCatalog reports whether theme.fontpack names a pack (or alias) the
// fontpack plugin can resolve.
func fontpackInCatalog(config *models.Config) bool {
	return KnownFontpack(config.Theme.Fontpack, config.FontpacksFile)
}

// KnownFontpack reports whether name resolves to a built-in font pack or
// alias. Any name is accepted when fontpacksFile selects a custom catalog.
func KnownFontpack(name, fontpacksFile string) bool {
	if name == "" || strings.TrimSpace(fontpacksFile) != "" {
		return true
	}
	source, err := fontpacks.BuiltinSource()
	if err != nil {
		return false
	}
	_, _, err = source.Catalog.ResolvePack(name)
	return err == nil
}

// KnownPalette reports whether name is a rendering-contract palette, a
// contract alias, or a site-local or user palette the loader can find.
func KnownPalette(name string) bool {
	c, err := renderingcontract.Load()
	if err != nil {
		return false
	}
	return knownPalette(&c, name)
}

func knownPalette(c *renderingcontract.Contract, name string) bool {
	if alias, ok := c.Aliases[name]; ok {
		name = alias
	}
	for _, item := range c.Palettes {
		if item.ID == name {
			return true
		}
	}
	// Site-local and user palettes (./palettes, ~/.config/markata-go/palettes)
	// are valid even though they are not part of the rendering contract.
	_, err := palettes.NewLoader().Load(name)
	return err == nil
}
