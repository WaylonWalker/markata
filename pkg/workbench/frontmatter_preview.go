package workbench

import (
	"html"
	"net/url"
	"path"
	"strings"

	"github.com/WaylonWalker/markata/pkg/plugins"
)

func renderFrontmatterPreview(frontmatterSource, markdown string) string {
	if strings.TrimSpace(frontmatterSource) == "" {
		return ""
	}
	metadata, _, err := plugins.ParseFrontmatter(frontmatterSource)
	if err != nil {
		return `<aside class="workbench-frontmatter-warning" role="status">Frontmatter preview is unavailable until the YAML is valid.</aside>`
	}
	if len(metadata) == 0 {
		return ""
	}
	return renderFrontmatterMetadataPreview(metadata, markdown)
}

func renderFrontmatterMetadataPreview(metadata map[string]interface{}, markdown string) string {
	templateName := strings.ToLower(strings.TrimSpace(plugins.GetString(metadata, "template")))
	if templateName == "" {
		templateName = strings.ToLower(strings.TrimSpace(plugins.GetString(metadata, "templateKey")))
	}
	title := strings.TrimSpace(plugins.GetString(metadata, "title"))
	description := strings.TrimSpace(plugins.GetString(metadata, "description"))
	image := strings.TrimSpace(plugins.GetString(metadata, frontmatterImageKey))
	linkURL := strings.TrimSpace(plugins.GetString(metadata, "link"))
	if linkURL == "" {
		linkURL = strings.TrimSpace(plugins.GetString(metadata, "url"))
	}

	var preview strings.Builder
	preview.WriteString(renderFrontmatterHeader(markdown, title, description, isLinkTemplate(templateName), isPhotoTemplate(templateName)))

	switch templateName {
	case "link", contentTemplateAliasBookmark, contentTemplateAliasStars, contentTemplateAliasTIL:
		if card := renderFrontmatterLinkCard(metadata, linkURL); card != "" {
			preview.WriteString(card)
		}
	case "photo", "shot", contentTemplateAliasShots, frontmatterImageKey:
		preview.WriteString(renderFrontmatterPhotoImage(image, title, description, markdown))
	}

	return preview.String()
}

func renderFrontmatterHeader(markdown, title, description string, linkTemplate, photoTemplate bool) string {
	showTitle := title != "" && !linkTemplate && !bodyStartsWithH1(markdown)
	showDescription := description != "" && !linkTemplate && !photoTemplate
	if !showTitle && !showDescription {
		return ""
	}

	var header strings.Builder
	header.WriteString(`<header class="workbench-frontmatter-header">`)
	if showTitle {
		header.WriteString(`<h1 class="p-name">`)
		header.WriteString(html.EscapeString(title))
		header.WriteString(`</h1>`)
	}
	if showDescription {
		header.WriteString(`<p class="p-summary">`)
		header.WriteString(html.EscapeString(description))
		header.WriteString(`</p>`)
	}
	header.WriteString(`</header>`)
	return header.String()
}

func renderFrontmatterPhotoImage(image, title, description, markdown string) string {
	if strings.Contains(markdown, "![Post image]") {
		return ""
	}
	imageURL := workbenchPreviewImageURL(image)
	if imageURL == "" {
		return ""
	}
	alt := title
	if alt == "" {
		alt = description
	}

	var figure strings.Builder
	figure.WriteString(`<figure class="photo-figure workbench-frontmatter-image"><img class="u-photo" src="`)
	figure.WriteString(html.EscapeString(imageURL))
	figure.WriteString(`" alt="`)
	figure.WriteString(html.EscapeString(alt))
	figure.WriteString(`" loading="lazy" decoding="async">`)
	if description != "" {
		figure.WriteString(`<figcaption class="p-summary">`)
		figure.WriteString(html.EscapeString(description))
		figure.WriteString(`</figcaption>`)
	}
	figure.WriteString(`</figure>`)
	return figure.String()
}

func isLinkTemplate(templateName string) bool {
	switch templateName {
	case "link", contentTemplateAliasBookmark, contentTemplateAliasStars, contentTemplateAliasTIL:
		return true
	default:
		return false
	}
}

func isPhotoTemplate(templateName string) bool {
	switch templateName {
	case "photo", "shot", contentTemplateAliasShots, frontmatterImageKey:
		return true
	default:
		return false
	}
}

func renderFrontmatterLinkCard(metadata map[string]interface{}, rawURL string) string {
	destination, parsedURL, ok := safePreviewHTTPURL(rawURL)
	if !ok {
		return ""
	}
	title := strings.TrimSpace(plugins.GetString(metadata, "title"))
	if title == "" {
		title = parsedURL.Hostname()
	}
	description := strings.TrimSpace(plugins.GetString(metadata, "description"))
	siteName := strings.TrimSpace(plugins.GetString(metadata, "site_name"))
	providerName := strings.TrimSpace(plugins.GetString(metadata, "provider_name"))
	mediaType := strings.TrimSpace(plugins.GetString(metadata, "type"))
	authorName := strings.TrimSpace(plugins.GetString(metadata, "author_name"))
	authorURL := strings.TrimSpace(plugins.GetString(metadata, "author_url"))
	imageURL := workbenchPreviewImageURL(plugins.GetString(metadata, frontmatterImageKey))

	var card strings.Builder
	card.WriteString(`<aside class="embed-card embed-card-external workbench-link-preview" data-markata-embed="true"><a class="embed-card-link" href="`)
	card.WriteString(html.EscapeString(destination))
	card.WriteString(`" target="_blank" rel="noopener noreferrer">`)
	if imageURL != "" {
		card.WriteString(`<div class="embed-card-image"><img src="`)
		card.WriteString(html.EscapeString(imageURL))
		card.WriteString(`" alt="" loading="lazy" decoding="async"></div>`)
	}
	card.WriteString(`<div class="embed-card-content"><div class="embed-card-title">`)
	card.WriteString(html.EscapeString(title))
	card.WriteString(`</div>`)
	if description != "" {
		card.WriteString(`<div class="embed-card-description">`)
		card.WriteString(html.EscapeString(description))
		card.WriteString(`</div>`)
	}
	if siteName == "" {
		siteName = providerName
	}
	card.WriteString(`<div class="embed-card-meta">`)
	if siteName != "" {
		card.WriteString(html.EscapeString(siteName))
		card.WriteString(` &middot; `)
	}
	card.WriteString(html.EscapeString(parsedURL.Hostname()))
	if providerName != "" && providerName != siteName {
		card.WriteString(` &middot; `)
		card.WriteString(html.EscapeString(providerName))
	}
	if mediaType != "" {
		card.WriteString(` &middot; `)
		card.WriteString(html.EscapeString(mediaType))
	}
	card.WriteString(`</div>`)
	card.WriteString(`</div></a>`)
	if authorName != "" {
		card.WriteString(`<div class="embed-card-source">`)
		if _, _, valid := safePreviewHTTPURL(authorURL); valid {
			card.WriteString(`<a href="`)
			card.WriteString(html.EscapeString(authorURL))
			card.WriteString(`" target="_blank" rel="noopener noreferrer">`)
			card.WriteString(html.EscapeString(authorName))
			card.WriteString(`</a>`)
		} else {
			card.WriteString(html.EscapeString(authorName))
		}
		card.WriteString(`</div>`)
	}
	card.WriteString(`</aside>`)
	return card.String()
}

func workbenchPreviewImageURL(rawImage string) string {
	rawImage = strings.TrimSpace(rawImage)
	if rawImage == "" {
		return ""
	}
	if absolute, _, ok := safePreviewHTTPURL(rawImage); ok {
		return absolute
	}
	parsed, err := url.Parse(rawImage)
	if err != nil || parsed.Scheme != "" || parsed.Host != "" || parsed.RawQuery != "" || parsed.Fragment != "" || !strings.HasPrefix(parsed.Path, "/images/") {
		return ""
	}
	staticPath := path.Clean(strings.TrimPrefix(parsed.Path, "/"))
	if staticPath == "images" || !strings.HasPrefix(staticPath, "images/") || strings.HasPrefix(staticPath, "../") {
		return ""
	}
	return "/workspace/image?path=" + url.QueryEscape(staticPath)
}

func safePreviewHTTPURL(rawURL string) (string, *url.URL, bool) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Hostname() == "" || parsed.User != nil || (parsed.Scheme != insecureScheme && parsed.Scheme != secureScheme) {
		return "", nil, false
	}
	canonicalURL := parsed.String()
	return canonicalURL, parsed, true
}

// Preserve an authored leading H1 as the page heading. A frontmatter title
// may intentionally differ from that heading; rendering both creates two H1s.
func bodyStartsWithH1(markdown string) bool {
	for _, line := range strings.Split(markdown, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return strings.HasPrefix(line, "# ")
	}
	return false
}
