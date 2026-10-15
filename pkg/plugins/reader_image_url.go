package plugins

import (
	"net/url"
	"strings"
)

// normalizeReaderImageURL resolves feed-relative media and OG image references
// against the article URL, not the reader's own domain. Discard unsupported
// schemes so an unusable image becomes a normal preview fallback.
func normalizeReaderImageURL(rawImageURL, articleURL string) string {
	rawImageURL = strings.TrimSpace(rawImageURL)
	if rawImageURL == "" {
		return ""
	}
	image, err := url.Parse(rawImageURL)
	if err != nil {
		return ""
	}
	if image.Scheme != "" && image.Scheme != schemeHTTP && image.Scheme != schemeHTTPS {
		return ""
	}
	if !image.IsAbs() {
		article, err := url.Parse(articleURL)
		if err != nil || (article.Scheme != schemeHTTP && article.Scheme != schemeHTTPS) || article.Hostname() == "" {
			return ""
		}
		image = article.ResolveReference(image)
	}
	if (image.Scheme != schemeHTTP && image.Scheme != schemeHTTPS) || image.Hostname() == "" {
		return ""
	}
	return normalizeExternalURL(image.String())
}
