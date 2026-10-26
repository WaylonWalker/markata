package plugins

import (
	"encoding/json"
	"html"
	"net/url"
	"strings"

	"github.com/WaylonWalker/markata/pkg/buildcache"
	"github.com/WaylonWalker/markata/pkg/models"
)

// encodeLocalPreviews prepares one immutable lookup shared by concurrent renders.
func encodeLocalPreviews(previews map[string]map[string]interface{}) map[string]string {
	encoded := make(map[string]string, len(previews))
	for path, preview := range previews {
		data, err := json.Marshal(preview)
		if err == nil {
			encoded[path] = string(data)
		}
	}
	return encoded
}

func localPreviewPath(href, basePath, siteURL string) string {
	ref, err := url.Parse(strings.TrimSpace(href))
	if err != nil || ref.Path == "" {
		return ""
	}
	if ref.IsAbs() || ref.Host != "" {
		if ref.Scheme != schemeHTTP && ref.Scheme != schemeHTTPS {
			return ""
		}
		siteHost := ""
		if site, err := url.Parse(siteURL); err == nil && site != nil {
			siteHost = site.Hostname()
		}
		host := strings.ToLower(ref.Hostname())
		if host != "localhost" && host != "127.0.0.1" && host != "::1" && (siteHost == "" || !strings.EqualFold(host, siteHost)) {
			return ""
		}
	}
	if strings.HasPrefix(ref.Path, "/") {
		return cleanNavPath(ref.Path)
	}
	base := &url.URL{Path: basePath}
	return cleanNavPath(base.ResolveReference(ref).Path)
}

func localPreviewForTag(tag string, previews map[string]string) string {
	return previews[cleanNavPath("/tags/"+models.Slugify(tag))]
}

func localPreviewHash(post *models.Post, previews map[string]string, siteURL string) string {
	var values strings.Builder
	for _, tag := range externalAnchorOpenRegex.FindAllString(post.ArticleHTML, -1) {
		attrs := parseAnchorAttrs(tag)
		if _, skip := attrs["data-no-preview"]; skip {
			continue
		}
		if path := localPreviewPath(attrs["href"], post.Href, siteURL); path != "" {
			values.WriteString(previews[path])
			values.WriteByte(0)
		}
	}
	for _, tag := range post.Tags {
		values.WriteString(localPreviewForTag(tag, previews))
		values.WriteByte(0)
	}
	if values.Len() == 0 {
		return ""
	}
	return buildcache.ContentHash(values.String())
}

func annotateLocalPreviewLinks(article, basePath, siteURL string, previews map[string]string) string {
	if !strings.Contains(article, "<a") {
		return article
	}
	return externalAnchorOpenRegex.ReplaceAllStringFunc(article, func(tag string) string {
		attrs := parseAnchorAttrs(tag)
		for _, name := range []string{"data-no-preview", "data-hover-card", "data-link-preview", "data-local-preview"} {
			if _, skip := attrs[name]; skip {
				return tag
			}
		}
		path := localPreviewPath(attrs["href"], basePath, siteURL)
		preview := previews[path]
		if preview == "" {
			return tag
		}
		end := len(tag) - 1
		if strings.HasSuffix(tag, "/>") {
			end--
		}
		return tag[:end] + ` data-local-preview="` + html.EscapeString(preview) + `"` + tag[end:]
	})
}
