package templates

import (
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/models"
)

var (
	trustedMediaMu      sync.RWMutex
	trustedMediaDomains = buildTrustedMediaDomainSet(models.DefaultTrustedMediaDomains)
	videoMimeTypes      = map[string]string{
		".avi":  "video/x-msvideo",
		".mkv":  "video/x-matroska",
		".mp4":  "video/mp4",
		".m4v":  "video/x-m4v",
		".webm": "video/webm",
		".mov":  "video/quicktime",
		".ogv":  "video/ogg",
		".ogg":  "video/ogg",
	}
	posterAliasOrder = []string{
		"poster_image",
		"poster",
		"video_poster",
		"video_thumbnail",
		"thumbnail",
		"thumb",
	}
)

// SetTrustedMediaDomains replaces the trusted domains allowlist for media helpers.
func SetTrustedMediaDomains(domains []string) {
	trustedMediaMu.Lock()
	defer trustedMediaMu.Unlock()
	trustedMediaDomains = map[string]struct{}{}
	if len(domains) == 0 {
		domains = models.DefaultTrustedMediaDomains
	}
	for _, domain := range domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}
		trustedMediaDomains[domain] = struct{}{}
	}
}

func buildTrustedMediaDomainSet(domains []string) map[string]struct{} {
	set := make(map[string]struct{}, len(domains))
	for _, domain := range domains {
		domain = strings.ToLower(strings.TrimSpace(domain))
		if domain == "" {
			continue
		}
		set[domain] = struct{}{}
	}
	return set
}

// WithSize appends or overwrites the w/h query parameters for trusted URLs.
// When height is 0 only the width param is set, letting the CDN preserve
// the original aspect ratio (width-only sizing).
func WithSize(raw string, width, height int) string {
	if width <= 0 {
		return raw
	}
	raw = normalizeTrustedMediaURL(raw)
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if !isTrustedURL(u) {
		return raw
	}
	qs := u.Query()
	qs.Set("w", strconv.Itoa(width))
	if height > 0 {
		qs.Set("h", strconv.Itoa(height))
	} else {
		qs.Del("h")
	}
	u.RawQuery = qs.Encode()
	return u.String()
}

// OriginalMediaURL removes resizing parameters from trusted media URLs while
// preserving access tokens and other query parameters. Unknown URLs are left
// unchanged.
func OriginalMediaURL(raw string) string {
	raw = normalizeTrustedMediaURL(raw)
	u, err := url.Parse(raw)
	if err != nil || !isTrustedURL(u) {
		return raw
	}
	query := u.Query()
	for _, key := range []string{"w", "h", "width", "height"} {
		query.Del(key)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

// ResponsiveSrcset returns width-descriptor candidates for trusted remote
// media URLs. Local paths are handled by image_optimization, which generates
// actual files instead of relying on a remote resizing service.
func ResponsiveSrcset(raw string, widths []int) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != editorialHTTPS) || !IsTrustedMediaURL(raw) {
		return ""
	}

	unique := make(map[int]struct{}, len(widths))
	ordered := make([]int, 0, len(widths))
	for _, width := range widths {
		if width <= 0 {
			continue
		}
		if _, exists := unique[width]; exists {
			continue
		}
		unique[width] = struct{}{}
		ordered = append(ordered, width)
	}
	if len(ordered) < 2 {
		return ""
	}
	sort.Ints(ordered)

	candidates := make([]string, 0, len(ordered))
	for _, width := range ordered {
		candidates = append(candidates, WithSize(raw, width, 0)+" "+strconv.Itoa(width)+"w")
	}
	return strings.Join(candidates, ", ")
}

// MediaDimensionsFromURL extracts known width/height query parameters from a media URL.
// It recognizes both w/h and width/height parameter names.
func MediaDimensionsFromURL(raw string) (width, height int, ok bool) {
	u, err := url.Parse(raw)
	if err != nil {
		return 0, 0, false
	}

	qs := u.Query()
	if width = positiveIntFromString(firstNonEmpty(qs.Get("w"), qs.Get("width"))); width > 0 {
		ok = true
	}
	if height = positiveIntFromString(firstNonEmpty(qs.Get("h"), qs.Get("height"))); height > 0 {
		ok = true
	}
	return width, height, ok
}

// IsVideoURL reports whether a URL ends with a known video extension.
func IsVideoURL(raw string) bool {
	ext := extensionFromURL(raw)
	_, ok := videoMimeTypes[ext]
	return ok
}

// VideoMIMEType infers the video MIME type for a URL.
func VideoMIMEType(raw string) string {
	ext := extensionFromURL(raw)
	return videoMimeTypes[ext]
}

// IsTrustedMediaURL reports whether the host/url is trusted for helper modifications.
func IsTrustedMediaURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return isTrustedURL(u)
}

// PosterURLFromMap resolves the first poster alias value or derives a .webp poster for trusted video URLs.
func PosterURLFromMap(data map[string]interface{}, mediaURL string) string {
	if data == nil {
		return ""
	}
	for _, alias := range posterAliasOrder {
		if val := stringFromMap(data, alias); val != "" {
			return normalizeTrustedMediaURL(val)
		}
	}
	if mediaURL == "" || !IsTrustedMediaURL(mediaURL) {
		return ""
	}
	return normalizeTrustedMediaURL(derivePosterFromVideo(normalizeTrustedMediaURL(mediaURL)))
}

func isTrustedURL(u *url.URL) bool {
	if u == nil {
		return false
	}
	if u.Host == "" {
		return true
	}
	host := strings.ToLower(u.Hostname())
	trustedMediaMu.RLock()
	defer trustedMediaMu.RUnlock()
	_, ok := trustedMediaDomains[host]
	return ok
}

func extensionFromURL(raw string) string {
	u, err := url.Parse(raw)
	// A Windows drive letter is parsed as a URL scheme, leaving the path
	// empty. Fall back to filepath.Ext for local paths in that form.
	if err != nil {
		return strings.ToLower(filepath.Ext(raw))
	}
	if u.Path == "" {
		if !isWindowsLocalPath(raw) {
			return ""
		}
		return strings.ToLower(filepath.Ext(raw))
	}
	return strings.ToLower(filepath.Ext(u.Path))
}

func isWindowsLocalPath(raw string) bool {
	if filepath.VolumeName(raw) != "" {
		return true
	}
	return len(raw) >= 2 && raw[1] == ':' &&
		((raw[0] >= 'a' && raw[0] <= 'z') || (raw[0] >= 'A' && raw[0] <= 'Z'))
}

func derivePosterFromVideo(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	ext := filepath.Ext(u.Path)
	if ext == "" {
		return ""
	}
	u.Path = strings.TrimSuffix(u.Path, ext) + ".webp"
	// Video URLs may carry resize hints for the source (e.g. ?width=400).
	// A poster is a different media resource; don't propagate those stale
	// dimensions into its URL alongside the poster's own w/h dimensions.
	// Preserve unrelated query parameters such as signed access tokens.
	query := u.Query()
	for _, key := range []string{"w", "h", "width", "height"} {
		query.Del(key)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

func normalizeTrustedMediaURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	if !isTrustedURL(u) {
		return raw
	}
	if u.Host != "" && u.Scheme != editorialHTTPS {
		u.Scheme = "https"
		return u.String()
	}
	return raw
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func positiveIntFromString(raw string) int {
	v, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil || v <= 0 {
		return 0
	}
	return v
}

func stringFromMap(data map[string]interface{}, key string) string {
	if v, ok := data[key]; ok {
		if s, ok := v.(string); ok {
			s = strings.TrimSpace(s)
			if s != "" {
				return s
			}
		}
	}
	return ""
}
