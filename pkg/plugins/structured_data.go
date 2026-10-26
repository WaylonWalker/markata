// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"encoding/json"
	"net/url"
	"strings"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

// StructuredDataPlugin generates JSON-LD Schema.org markup, OpenGraph meta tags,
// and Twitter Cards for SEO and social media optimization.
type StructuredDataPlugin struct{}

// NewStructuredDataPlugin creates a new StructuredDataPlugin.
func NewStructuredDataPlugin() *StructuredDataPlugin {
	return &StructuredDataPlugin{}
}

// Name returns the unique name of the plugin.
func (p *StructuredDataPlugin) Name() string {
	return "structured_data"
}

// Priority returns the plugin priority for the given stage.
// Runs in mid-transform, after description plugin has run.
func (p *StructuredDataPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageTransform {
		return lifecycle.PriorityDefault // 500
	}
	return lifecycle.PriorityDefault
}

// Transform generates structured data for each post.
func (p *StructuredDataPlugin) Transform(m *lifecycle.Manager) error {
	config := m.Config()

	// Check if structured data is enabled
	seoConfig := getSEOConfig(config)
	if !seoConfig.StructuredData.IsEnabled() {
		return nil
	}

	posts := m.FilterPosts(func(post *models.Post) bool {
		if post.Skip || post.Draft {
			return false
		}
		return post.PlainTitle() != ""
	})

	return m.ProcessPostsSliceConcurrently(posts, func(post *models.Post) error {
		return p.generateStructuredData(post, config, &seoConfig)
	})
}

// generateStructuredData creates all structured data for a post.
func (p *StructuredDataPlugin) generateStructuredData(post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) error {
	sd := models.NewStructuredData()

	// Generate JSON-LD
	jsonLD, err := p.generateJSONLD(post, config, seoConfig)
	if err == nil && jsonLD != "" {
		sd.JSONLD = jsonLD
	}

	// Generate OpenGraph tags
	p.generateOpenGraph(sd, post, config, seoConfig)

	// Generate Twitter Card tags
	p.generateTwitterCard(sd, post, config, seoConfig)

	// Store in post.Extra
	post.Set("structured_data", sd)

	return nil
}

// generateJSONLD creates JSON-LD Schema.org markup for a post.
func (p *StructuredDataPlugin) generateJSONLD(post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) (string, error) {
	siteURL := getSiteURL(config)
	postURL := siteURL + post.Href

	// Create BlogPosting schema
	bp := models.NewBlogPosting(post.PlainTitle(), postURL)

	// Add description
	if post.Description != nil {
		bp.Description = *post.Description
	}

	// Add dates
	if post.Date != nil {
		bp.DatePublished = post.Date.Format("2006-01-02T15:04:05Z07:00")

		// Use Modified field if set, otherwise fall back to Date
		if post.Modified != nil {
			bp.DateModified = post.Modified.Format("2006-01-02T15:04:05Z07:00")
		} else {
			bp.DateModified = bp.DatePublished
		}
	}

	// Add image
	imageURL := p.getPostImage(post, seoConfig)
	if imageURL != "" {
		bp.Image = p.makeAbsoluteURL(imageURL, siteURL)
	}

	// Add keywords from tags
	if len(post.Tags) > 0 {
		bp.Keywords = post.Tags
	}

	// Add author
	bp.Author = p.getAuthor(post, config, seoConfig)

	// Add publisher
	bp.Publisher = p.getPublisher(config, seoConfig)

	// Marshal to JSON
	jsonBytes, err := json.MarshalIndent(bp, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// generateOpenGraph creates OpenGraph meta tags.
func (p *StructuredDataPlugin) generateOpenGraph(sd *models.StructuredData, post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) {
	siteURL := getSiteURL(config)
	siteTitle := getSiteTitle(config)
	siteDescription := getSiteDescription(config)
	postURL := siteURL + post.Href

	// Required tags
	sd.AddOpenGraph("og:title", post.PlainTitle())
	sd.AddOpenGraph("og:url", postURL)
	sd.AddOpenGraph("og:site_name", siteTitle)

	// Type - article for posts with dates, website otherwise
	if post.Date != nil {
		sd.AddOpenGraph("og:type", "article")
	} else {
		sd.AddOpenGraph("og:type", "website")
	}

	// Description
	if post.Description != nil {
		sd.AddOpenGraph("og:description", *post.Description)
	} else if siteDescription != "" {
		sd.AddOpenGraph("og:description", siteDescription)
	}

	// Image - use og_image_service if available, otherwise fall back to post image
	imageURL := p.getOGImageURL(post, config, seoConfig)
	if imageURL != "" {
		sd.AddOpenGraph("og:image", imageURL)
		sd.AddOpenGraph("og:image:width", "1200")
		sd.AddOpenGraph("og:image:height", "600")
	}

	// Locale
	sd.AddOpenGraph("og:locale", "en_US")

	// Article-specific tags
	if post.Date != nil {
		sd.AddOpenGraph("article:published_time", post.Date.Format("2006-01-02T15:04:05Z07:00"))

		// Modified time - use Modified field if set
		if post.Modified != nil {
			sd.AddOpenGraph("article:modified_time", post.Modified.Format("2006-01-02T15:04:05Z07:00"))
		}

		// Author URL
		author := p.getAuthor(post, config, seoConfig)
		if author != nil && author.URL != "" {
			sd.AddOpenGraph("article:author", author.URL)
		}

		// Tags
		for _, tag := range post.Tags {
			sd.AddOpenGraph("article:tag", tag)
		}
	}
}

// generateTwitterCard creates Twitter Card meta tags.
func (p *StructuredDataPlugin) generateTwitterCard(sd *models.StructuredData, post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) {
	// Card type - summary_large_image if we have an image
	// Use og_image_service if available, otherwise fall back to post image
	imageURL := p.getTwitterImageURL(post, config, seoConfig)
	if imageURL != "" {
		sd.AddTwitter("twitter:card", "summary_large_image")
		sd.AddTwitter("twitter:image", imageURL)
	} else {
		sd.AddTwitter("twitter:card", "summary")
	}

	// Site handle
	if seoConfig.TwitterHandle != "" {
		sd.AddTwitter("twitter:site", "@"+seoConfig.TwitterHandle)
	}

	// Creator handle (use first author's twitter if available, otherwise site handle)
	creatorHandle := p.getFirstAuthorTwitterHandle(post, config, seoConfig)
	if creatorHandle != "" {
		sd.AddTwitter("twitter:creator", "@"+creatorHandle)
	}

	// Title
	sd.AddTwitter("twitter:title", post.PlainTitle())

	// Description (truncated to 200 chars for Twitter)
	if post.Description != nil {
		desc := *post.Description
		if len(desc) > 200 {
			desc = desc[:197] + "..."
		}
		sd.AddTwitter("twitter:description", desc)
	}
}

// getPostImage returns the image URL for a post.
// Checks frontmatter image, social_image, then falls back to default.
func (p *StructuredDataPlugin) getPostImage(post *models.Post, seoConfig *models.SEOConfig) string {
	// Check for social_image override first
	if socialImage, ok := post.Extra["social_image"]; ok {
		if imgStr, ok := socialImage.(string); ok && imgStr != "" {
			return imgStr
		}
	}

	// Check for image in frontmatter
	if image, ok := post.Extra["image"]; ok {
		if imgStr, ok := image.(string); ok && imgStr != "" {
			return imgStr
		}
	}

	// Fall back to default image
	return seoConfig.DefaultImage
}

// getOGImageURL returns the OG image URL for a post.
// If og_image_service is configured, generates a screenshot service URL.
// Otherwise falls back to getPostImage().
func (p *StructuredDataPlugin) getOGImageURL(post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) string {
	siteURL := getSiteURL(config)

	// If og_image_service is configured, use the screenshot service
	if seoConfig.OGImageService != "" {
		// Build the URL for the OG card page: {site_url}{post_href}og/
		ogCardURL := siteURL + post.Href + "og/"

		// Build the screenshot service URL with OG image dimensions (1200x600)
		// URL-encode the ogCardURL to handle special characters safely
		// Format: {og_image_service}?url={og_card_url}&height=600&width=1200&scaled_width=1200&scaled_height=600&format=jpg
		return seoConfig.OGImageService + "?url=" + url.QueryEscape(ogCardURL) + "&height=600&width=1200&scaled_width=1200&scaled_height=600&format=jpg"
	}

	// Fall back to post image
	imageURL := p.getPostImage(post, seoConfig)
	if imageURL != "" {
		return p.makeAbsoluteURL(imageURL, siteURL)
	}

	return ""
}

// getTwitterImageURL returns the Twitter card image URL for a post.
// If og_image_service is configured, generates a screenshot service URL with Twitter dimensions.
// Otherwise falls back to getPostImage().
func (p *StructuredDataPlugin) getTwitterImageURL(post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) string {
	siteURL := getSiteURL(config)

	// If og_image_service is configured, use the screenshot service with Twitter dimensions
	if seoConfig.OGImageService != "" {
		// Build the URL for the OG card page: {site_url}{post_href}og/
		ogCardURL := siteURL + post.Href + "og/"

		// Build the screenshot service URL with Twitter image dimensions (1280x640)
		// URL-encode the ogCardURL to handle special characters safely
		// Format: {og_image_service}?url={og_card_url}&height=640&width=1280&scaled_width=1280&scaled_height=640&format=jpg
		return seoConfig.OGImageService + "?url=" + url.QueryEscape(ogCardURL) + "&height=640&width=1280&scaled_width=1280&scaled_height=640&format=jpg"
	}

	// Fall back to post image
	imageURL := p.getPostImage(post, seoConfig)
	if imageURL != "" {
		return p.makeAbsoluteURL(imageURL, siteURL)
	}

	return ""
}

// getAuthor returns author SchemaAgent for a post.
func (p *StructuredDataPlugin) getAuthor(post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) *models.SchemaAgent {
	// Check for author in frontmatter
	var authorName string
	if author, ok := post.Extra["author"]; ok {
		if authorStr, ok := author.(string); ok && authorStr != "" {
			authorName = authorStr
		}
	}

	// If we have a custom author name, create a basic agent
	if authorName != "" {
		return models.NewSchemaAgent("Person", authorName)
	}

	// Use default author from config
	if seoConfig.StructuredData.DefaultAuthor != nil {
		da := seoConfig.StructuredData.DefaultAuthor
		return models.NewSchemaAgent(da.Type, da.Name).WithURL(da.URL)
	}

	// Fall back to site author
	siteAuthor := getSiteAuthor(config)
	if siteAuthor != "" {
		return models.NewSchemaAgent("Person", siteAuthor)
	}

	return nil
}

// getPublisher returns the publisher SchemaAgent for the site.
func (p *StructuredDataPlugin) getPublisher(config *lifecycle.Config, seoConfig *models.SEOConfig) *models.SchemaAgent {
	siteURL := getSiteURL(config)

	// Use publisher from config
	if seoConfig.StructuredData.Publisher != nil {
		pub := seoConfig.StructuredData.Publisher
		agent := models.NewSchemaAgent(pub.Type, pub.Name).WithURL(pub.URL)
		if pub.Logo != "" {
			agent.WithLogo(p.makeAbsoluteURL(pub.Logo, siteURL))
		} else if seoConfig.LogoURL != "" {
			agent.WithLogo(p.makeAbsoluteURL(seoConfig.LogoURL, siteURL))
		}
		return agent
	}

	// Fall back to site title as Organization
	siteTitle := getSiteTitle(config)
	if siteTitle != "" {
		agent := models.NewSchemaAgent("Organization", siteTitle).WithURL(siteURL)
		if seoConfig.LogoURL != "" {
			agent.WithLogo(p.makeAbsoluteURL(seoConfig.LogoURL, siteURL))
		}
		return agent
	}

	return nil
}

// getFirstAuthorTwitterHandle returns the Twitter handle for the first author of a post.
// Checks post frontmatter first, then author config, then falls back to site handle.
func (p *StructuredDataPlugin) getFirstAuthorTwitterHandle(post *models.Post, config *lifecycle.Config, seoConfig *models.SEOConfig) string {
	// Check for author's twitter handle in post frontmatter
	if twitterHandle, ok := post.Extra["twitter"]; ok {
		if handleStr, ok := twitterHandle.(string); ok && handleStr != "" {
			// Remove @ if present
			return strings.TrimPrefix(handleStr, "@")
		}
	}

	// Check first author's social links if multi-author is configured
	authorIDs := post.GetAuthors()
	if len(authorIDs) > 0 {
		// Get models.Config from lifecycle.Config.Extra
		if modelsConfig, ok := config.Extra["models_config"].(*models.Config); ok {
			if modelsConfig.Authors.Authors != nil {
				if author, ok := modelsConfig.Authors.Authors[authorIDs[0]]; ok {
					if twitter, ok := author.Social["twitter"]; ok && twitter != "" {
						return strings.TrimPrefix(twitter, "@")
					}
				}
			}
		}
	}

	// Fall back to site handle
	return seoConfig.TwitterHandle
}

// makeAbsoluteURL converts a relative URL to an absolute URL.
func (p *StructuredDataPlugin) makeAbsoluteURL(rawURL, siteURL string) string {
	if rawURL == "" {
		return ""
	}

	// Already absolute
	if strings.HasPrefix(rawURL, "http://") || strings.HasPrefix(rawURL, "https://") {
		return rawURL
	}

	// Protocol-relative
	if strings.HasPrefix(rawURL, "//") {
		return "https:" + rawURL
	}

	// Relative URL - prepend site URL
	siteURL = strings.TrimSuffix(siteURL, "/")
	if !strings.HasPrefix(rawURL, "/") {
		rawURL = "/" + rawURL
	}
	return siteURL + rawURL
}

// getSEOConfig retrieves the SEOConfig from lifecycle.Config.Extra.
func getSEOConfig(config *lifecycle.Config) models.SEOConfig {
	if config.Extra != nil {
		if seo, ok := config.Extra["seo"].(models.SEOConfig); ok {
			return seo
		}
	}
	return models.NewSEOConfig()
}

// Ensure StructuredDataPlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*StructuredDataPlugin)(nil)
	_ lifecycle.TransformPlugin = (*StructuredDataPlugin)(nil)
	_ lifecycle.PriorityPlugin  = (*StructuredDataPlugin)(nil)
)
