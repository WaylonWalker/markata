package config

import (
	"os"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestApplyEnvOverrides_PagefindSettings(t *testing.T) {
	tests := []struct {
		name         string
		envVars      map[string]string
		expectedFunc func(*models.Config) bool
	}{
		{
			name: "search_endpoint",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_ENDPOINT": "https://go.waylonwalker.com/api/search",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Endpoint == "https://go.waylonwalker.com/api/search"
			},
		},
		{
			name: "search_bleve_endpoint",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_BLEVE_ENDPOINT": "https://go.waylonwalker.com/api/search",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Bleve.Endpoint == "https://go.waylonwalker.com/api/search"
			},
		},
		{
			name: "search_enabled false",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_ENABLED": "false",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Enabled != nil && *c.Search.Enabled == false
			},
		},
		{
			name: "search_pagefind_auto_install true",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_PAGEFIND_AUTO_INSTALL": "true",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Pagefind.AutoInstall != nil && *c.Search.Pagefind.AutoInstall == true
			},
		},
		{
			name: "search_pagefind_auto_install false",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_PAGEFIND_AUTO_INSTALL": "false",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Pagefind.AutoInstall != nil && *c.Search.Pagefind.AutoInstall == false
			},
		},
		{
			name: "search_pagefind_cache_dir",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_PAGEFIND_CACHE_DIR": "/tmp/pagefind-cache",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Pagefind.CacheDir == "/tmp/pagefind-cache"
			},
		},
		{
			name: "search_pagefind_version",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_PAGEFIND_VERSION": "v1.4.0",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Pagefind.Version == "v1.4.0"
			},
		},
		{
			name: "search_pagefind_bundle_dir",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_PAGEFIND_BUNDLE_DIR": "_search",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Pagefind.BundleDir == "_search"
			},
		},
		{
			name: "search_pagefind_verbose true",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_PAGEFIND_VERBOSE": "true",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Pagefind.Verbose != nil && *c.Search.Pagefind.Verbose == true
			},
		},
		{
			name: "multiple pagefind settings",
			envVars: map[string]string{
				"MARKATA_GO_SEARCH_PAGEFIND_AUTO_INSTALL": "true",
				"MARKATA_GO_SEARCH_PAGEFIND_CACHE_DIR":    "/tmp/cache",
				"MARKATA_GO_SEARCH_PAGEFIND_VERSION":      "v1.4.0",
				"MARKATA_GO_SEARCH_PAGEFIND_VERBOSE":      "true",
			},
			expectedFunc: func(c *models.Config) bool {
				return c.Search.Pagefind.AutoInstall != nil && *c.Search.Pagefind.AutoInstall == true &&
					c.Search.Pagefind.CacheDir == "/tmp/cache" &&
					c.Search.Pagefind.Version == "v1.4.0" &&
					c.Search.Pagefind.Verbose != nil && *c.Search.Pagefind.Verbose == true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Clean up environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			// Create a config with default search settings
			config := &models.Config{Search: models.NewSearchConfig()}

			// Apply env overrides
			err := ApplyEnvOverrides(config)
			if err != nil {
				t.Fatalf("ApplyEnvOverrides() error = %v", err)
			}

			// Check expected condition
			if !tt.expectedFunc(config) {
				t.Errorf("ApplyEnvOverrides() did not set expected values for %s", tt.name)
			}
		})
	}
}
