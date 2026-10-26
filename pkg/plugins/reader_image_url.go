package plugins

import (
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var readerGhostSizePattern = regexp.MustCompile(`(/content/images/size/w)([1-9]\d*)(/)`)

// readerImageSrcset preserves the original preview while offering smaller
// raster variants on Ghost's image CDN. Other providers keep their own URLs.
func readerImageSrcset(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != schemeHTTPS || u.Host != "storage.ghost.io" || u.User != nil {
		return ""
	}
	switch strings.ToLower(path.Ext(u.Path)) {
	case extPNG, extJPG, extJPEG, extWebP:
	default:
		return ""
	}
	escaped := u.EscapedPath()
	match := readerGhostSizePattern.FindStringSubmatchIndex(escaped)
	if match == nil {
		return ""
	}
	maximum, err := strconv.Atoi(escaped[match[4]:match[5]])
	if err != nil || maximum <= 320 {
		return ""
	}
	widths := []int{320, 480, 640, 768, 960, 1200, 1600}
	candidates := make([]string, 0, len(widths)+1)
	for _, width := range widths {
		if width >= maximum {
			break
		}
		variant := *u
		variant.RawPath = escaped[:match[4]] + strconv.Itoa(width) + escaped[match[5]:]
		decodedPath, decodeErr := url.PathUnescape(variant.RawPath)
		if decodeErr != nil {
			return ""
		}
		variant.Path = decodedPath
		candidates = append(candidates, variant.String()+" "+strconv.Itoa(width)+"w")
	}
	candidates = append(candidates, raw+" "+strconv.Itoa(maximum)+"w")
	return strings.Join(candidates, ", ")
}

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
