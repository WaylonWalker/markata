package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata-go/pkg/models"
)

const shortsRoutesFile = ".markata-shorts-routes.json"

// shortsRouteLedger records the exact routes owned by the previous successful
// publication, including feeds that may no longer appear in configuration.
type shortsRouteLedger struct {
	Version int               `json:"version"`
	Routes  map[string]string `json:"routes"`
}

func shortsPublishedOwner(outputDir, route string) string {
	if validateShortsPath(route) != nil {
		return ""
	}
	marker := filepath.Join(outputDir, strings.Trim(route, "/"), "data", "index.json")
	data, err := os.ReadFile(marker)
	if err != nil {
		return ""
	}
	var index struct {
		FeedSlug string `json:"feed_slug"`
	}
	if json.Unmarshal(data, &index) != nil {
		return ""
	}
	return index.FeedSlug
}

// reconcileShortsRoutes runs only after all current feed writes complete.
// Compare the current route set with the previous successful build's owners.
// Do not scan the whole site on every warm build, or remove routes whose
// marker no longer identifies the original feed.
func reconcileShortsRoutes(outputDir string, feeds []models.FeedConfig, posts []*models.Post) error {
	current := make(map[string]string)
	for i := range feeds {
		feed := &feeds[i]
		if feed.HasView(models.FeedViewShorts) {
			current[feed.ShortsURL()] = feed.Slug
		}
	}
	ledgerPath := filepath.Join(outputDir, shortsRoutesFile)
	previous := shortsRouteLedger{Version: 1, Routes: make(map[string]string)}
	data, err := os.ReadFile(ledgerPath)
	if err == nil {
		if err := json.Unmarshal(data, &previous); err != nil || previous.Version != 1 {
			return fmt.Errorf("invalid Shorts route ledger %s", ledgerPath)
		}
	} else if !os.IsNotExist(err) {
		return err
	} else {
		// Upgrade existing sites that generated the standard /shorts/ route
		// before the ownership ledger existed.
		if slug := shortsPublishedOwner(outputDir, "/shorts/"); slug != "" {
			previous.Routes["/shorts/"] = slug
		}
	}
	reserved := make(map[string]bool)
	for i := range feeds {
		feed := &feeds[i]
		if feed.Formats.HTML {
			route := "/" + strings.Trim(feed.Slug, "/") + "/"
			if feed.Slug == "" {
				route = "/"
			}
			reserved[route] = true
		}
	}
	for _, post := range posts {
		if post != nil && !post.Skip && !post.Draft && post.Published {
			reserved["/"+strings.Trim(post.Slug, "/")+"/"] = true
		}
	}
	for route, owner := range previous.Routes {
		if _, stillPublished := current[route]; stillPublished {
			continue
		}
		if shortsPublishedOwner(outputDir, route) != owner {
			continue
		}
		if reserved[route] {
			// A current post or feed owns index.html. Only remove the stale
			// Shorts metadata; never delete another publisher's HTML.
			if err := os.RemoveAll(filepath.Join(outputDir, strings.Trim(route, "/"), "data")); err != nil {
				return err
			}
			continue
		}
		if err := cleanupDisabledShorts(outputDir, route, owner); err != nil {
			return err
		}
	}
	if len(current) == 0 && len(previous.Routes) == 0 {
		return nil
	}
	encoded, err := json.Marshal(shortsRouteLedger{Version: 1, Routes: current})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(outputDir, ".shorts-routes-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(append(encoded, '\n')); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), ledgerPath)
}
