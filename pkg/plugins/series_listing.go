// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// SeriesListingInfo contains the data rendered for one series on the series index.
type SeriesListingInfo struct {
	Name        string
	Slug        string
	Title       string
	Description string
	Href        string
	Count       int
}

// Write generates a top-level series listing page for the series discovered
// during the build. SeriesPlugin already participates in Collect; defining the
// Write hook here keeps series discovery and series navigation in one plugin.
func (p *SeriesPlugin) Write(m *lifecycle.Manager) error {
	config := m.Config()
	seriesCfg := parseSeriesConfig(config)
	listingSlug := strings.Trim(seriesCfg.SlugPrefix, "/")
	if listingSlug == "" {
		log.Printf("[series] warning: cannot generate series listing with an empty slug_prefix")
		return nil
	}

	seriesList := p.collectSeriesListing(m.Posts(), seriesCfg)
	if len(seriesList) == 0 {
		return nil
	}

	return p.renderSeriesListing(config, listingSlug, seriesList)
}

func (p *SeriesPlugin) collectSeriesListing(posts []*models.Post, cfg seriesConfig) []SeriesListingInfo {
	groups := p.groupPostsBySeries(posts, cfg)
	seriesList := make([]SeriesListingInfo, 0, len(groups))

	for _, group := range groups {
		published := filterSeriesOutputPosts(group.posts)
		if len(published) == 0 {
			continue
		}

		feedSlug := buildSeriesFeedSlug(cfg.SlugPrefix, group.slug)
		info := SeriesListingInfo{
			Name:  group.name,
			Slug:  group.slug,
			Title: seriesDisplayTitle(group.name, group.cfg),
			Href:  "/" + strings.Trim(feedSlug, "/") + "/",
			Count: len(published),
		}
		if group.cfg != nil {
			info.Description = group.cfg.Description
		}
		seriesList = append(seriesList, info)
	}

	sort.Slice(seriesList, func(i, j int) bool {
		left := strings.ToLower(seriesList[i].Title)
		right := strings.ToLower(seriesList[j].Title)
		if left == right {
			return seriesList[i].Slug < seriesList[j].Slug
		}
		return left < right
	})

	return seriesList
}

func (p *SeriesPlugin) renderSeriesListing(config *lifecycle.Config, listingSlug string, seriesList []SeriesListingInfo) error {
	seriesDir := filepath.Join(config.OutputDir, filepath.FromSlash(listingSlug))
	if err := os.MkdirAll(seriesDir, 0o755); err != nil {
		return fmt.Errorf("creating series listing directory: %w", err)
	}

	templatesDir := PluginNameTemplates
	if extra, ok := config.Extra["templates_dir"].(string); ok && extra != "" {
		templatesDir = extra
	}
	engine, err := templates.NewEngineWithTheme(templatesDir, getThemeName(config))
	if err != nil {
		return fmt.Errorf("creating series listing template engine: %w", err)
	}
	if !engine.TemplateExists("series.html") {
		log.Printf("[series] warning: template %q not found, skipping series listing page", "series.html")
		return nil
	}

	title := "Series"
	description := "Browse all series."
	syntheticPost := &models.Post{
		Slug:        listingSlug,
		Title:       &title,
		Description: &description,
	}
	ctx := templates.NewContext(syntheticPost, "", ToModelsConfig(config))
	ctx.Extra["series_list"] = seriesList
	ctx.Extra["total_series"] = len(seriesList)

	html, err := engine.Render("series.html", ctx)
	if err != nil {
		return fmt.Errorf("rendering series listing template: %w", err)
	}

	outputPath := filepath.Join(seriesDir, "index.html")
	//nolint:gosec // G306: generated site output must be readable by the web server.
	if err := os.WriteFile(outputPath, []byte(html), 0o644); err != nil {
		return fmt.Errorf("writing series listing page: %w", err)
	}

	log.Printf("[series] Generated /%s/ with %d series", listingSlug, len(seriesList))
	return nil
}

var _ lifecycle.WritePlugin = (*SeriesPlugin)(nil)
