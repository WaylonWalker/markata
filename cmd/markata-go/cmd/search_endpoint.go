package cmd

import (
	"net/url"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
)

const defaultSearchEndpoint = "/api/search"

func configuredSearchEndpoints(cfg *models.Config) (clientEndpoint, handlerPath string) {
	clientEndpoint = defaultSearchEndpoint
	if cfg != nil {
		clientEndpoint = cfg.Search.Bleve.EndpointOrDefault(cfg.Search.SearchEndpoint())
	}
	return clientEndpoint, normalizeSearchEndpointPath(clientEndpoint)
}

func normalizeSearchEndpointPath(endpoint string) string {
	trimmed := strings.TrimSpace(endpoint)
	if trimmed == "" {
		return defaultSearchEndpoint
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return ensureLeadingSlash(trimmed)
	}

	if parsed.Path == "" {
		return "/"
	}

	return ensureLeadingSlash(parsed.Path)
}

func ensureLeadingSlash(value string) string {
	if strings.HasPrefix(value, "/") {
		return value
	}
	return "/" + value
}
