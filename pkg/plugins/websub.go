package plugins

import (
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

func getWebSubConfig(config *lifecycle.Config) models.WebSubConfig {
	if config != nil && config.Extra != nil {
		if ws, ok := config.Extra["websub"].(models.WebSubConfig); ok {
			return ws
		}
	}
	return models.NewWebSubConfig()
}

func getWebSubHubs(config *lifecycle.Config) []string {
	websub := getWebSubConfig(config)
	if !websub.IsEnabled() {
		return nil
	}

	hubs := make([]string, 0, len(websub.Hubs))
	for _, hub := range websub.Hubs {
		trimmed := strings.TrimSpace(hub)
		if trimmed == "" {
			continue
		}
		hubs = append(hubs, trimmed)
	}
	return hubs
}
