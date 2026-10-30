package config

import (
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

// validateRendererDefaults checks integration-only settings, keeping full site
// validation independent from transport-specific render policy.
func validateRendererDefaults(config *models.Config) []error {
	var errs []error
	// The profile is an integration preference, not the site-build lifecycle.
		// Reject misspellings instead of recording a setting no consumer can honor.
		switch config.MarkdownConfig.RenderProfile {
		case "", "web", "offline-interactive", "offline-static":
			// Recognized. Offline-static is reserved and standalone renders reject it.
		default:
			errs = append(errs, ValidationError{
				Field: "markdown.render_profile",
				Message: "must be one of: web, offline-interactive, offline-static",
			})
		}
	// Standalone integration asset defaults are independent of site-build assets.
		switch config.MarkdownConfig.RenderAssetMode {
		case "", "linked", "inline":
		default:
			errs = append(errs, ValidationError{Field: "markdown.render_asset_mode",
				Message: "must be one of: linked, inline"})
		}
		assetBase := config.MarkdownConfig.RenderAssetBase
		if assetBase != "" && (!strings.HasPrefix(assetBase, "/") ||
			strings.HasPrefix(assetBase, "//") ||
			strings.ContainsAny(assetBase, "?#\\\r\n") ||
			strings.Contains(assetBase, "/../")) {
			errs = append(errs, ValidationError{Field: "markdown.render_asset_base",
				Message: "must be a local absolute URL path without a query, fragment, or parent traversal"})
		}
	return errs
}
