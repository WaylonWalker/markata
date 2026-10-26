package templates

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/htmltotext"
	"github.com/WaylonWalker/markata/pkg/models"

	"github.com/flosch/pongo2/v6"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var registerOnce sync.Once

// HumanDateFormat is the canonical visible HTML date format.
const HumanDateFormat = "Jan 2, 2006"

// assetHashRegistry stores asset path -> hash mappings for cache busting.
// This is set by SetAssetHashes() before rendering templates.
var assetHashRegistry map[string]string
var assetHashMu sync.RWMutex

// SetAssetHashes merges new hashes into the global asset hash registry.
// This should be called before rendering templates that use theme_asset_hashed.
// Multiple plugins can call this to register their generated assets.
func SetAssetHashes(hashes map[string]string) {
	assetHashMu.Lock()
	defer assetHashMu.Unlock()

	// Initialize the registry if needed
	if assetHashRegistry == nil {
		assetHashRegistry = make(map[string]string)
	}

	// Merge new hashes (later calls can override earlier ones)
	for path, hash := range hashes {
		assetHashRegistry[path] = hash
	}
}

// GetAssetHash retrieves the hash for a given asset path.
func GetAssetHash(path string) (string, bool) {
	assetHashMu.RLock()
	defer assetHashMu.RUnlock()
	if assetHashRegistry == nil {
		return "", false
	}
	hash, ok := assetHashRegistry[path]
	return hash, ok
}

// Pre-compiled regex patterns for performance - compiled once at package init
var (
	// Used by filterStripTags and stripHTMLTagsHelper
	htmlTagRe = regexp.MustCompile(`<[^>]*>`)

	// Used by filterStripTags, cleanExcerptHTML, and stripHTMLTagsHelper
	whitespaceRe = regexp.MustCompile(`\s+`)

	// Used by filterLinebreaks
	paragraphSplitRe = regexp.MustCompile(`\n\n+`)

	// Used by removeAdmonitions
	admonitionDivRe     = regexp.MustCompile(`(?s)<div class="admonition[^"]*">.*?</div>`)
	admonitionDetailsRe = regexp.MustCompile(`(?s)<details class="admonition[^"]*">.*?</details>`)

	// Used by collectParagraphs and filterExcerpt
	excerptParagraphRe = regexp.MustCompile(`(?s)<p[^>]*>(.*?)</p>`)

	// Used by cleanExcerptHTML - pre-compiled patterns for each block tag
	blockTagOpenRe  = compileBlockTagPatterns(true)
	blockTagCloseRe = compileBlockTagPatterns(false)
)

// compileBlockTagPatterns creates regex patterns for block tags.
// If openTag is true, creates opening tag patterns; otherwise closing tag patterns.
func compileBlockTagPatterns(openTag bool) map[string]*regexp.Regexp {
	blockTags := []string{"div", "span", "section", "article", "header", "footer", "nav", "aside", "pre"}
	patterns := make(map[string]*regexp.Regexp, len(blockTags))

	for _, tag := range blockTags {
		if openTag {
			patterns[tag] = regexp.MustCompile(`(?i)<` + tag + `[^>]*>`)
		} else {
			patterns[tag] = regexp.MustCompile(`(?i)</` + tag + `>`)
		}
	}
	return patterns
}

// registerFilters registers all custom template filters with pongo2.
// This is called once when the first Engine is created.
// The pongo2 registration functions return errors only for duplicate registrations,
// which won't occur due to sync.Once protection.
//
//nolint:errcheck // pongo2 filter registration errors are only for duplicates, protected by sync.Once
func registerFilters() {
	registerOnce.Do(func() {
		// Date formatting filters
		pongo2.RegisterFilter("rss_date", filterRSSDate)
		pongo2.RegisterFilter("atom_date", filterAtomDate)
		pongo2.RegisterFilter("date_format", filterDateFormat)
		pongo2.RegisterFilter("human_date", filterHumanDate)
		// Override the built-in date filter to handle *time.Time and string parsing
		pongo2.ReplaceFilter("date", filterDate)

		// String manipulation filters
		pongo2.RegisterFilter("slugify", filterSlugify)
		pongo2.RegisterFilter("truncate", filterTruncate)
		pongo2.RegisterFilter("truncatewords", filterTruncateWords)
		pongo2.RegisterFilter("endswith", filterEndsWith)
		pongo2.RegisterFilter("startswith", filterStartsWith)
		pongo2.RegisterFilter("split", filterSplit)
		pongo2.RegisterFilter("replace", filterReplace)

		// Default/fallback filter
		pongo2.RegisterFilter("default_if_none", filterDefaultIfNone)

		// Collection filters
		pongo2.RegisterFilter("length", filterLength)
		pongo2.RegisterFilter("first", filterFirst)
		pongo2.RegisterFilter("last", filterLast)
		pongo2.RegisterFilter("join", filterJoin)
		pongo2.RegisterFilter("reverse", filterReverse)
		pongo2.RegisterFilter("sort", filterSort)
		pongo2.RegisterFilter("selectattr", filterSelectAttr)
		pongo2.RegisterFilter("rejectattr", filterRejectAttr)
		pongo2.RegisterFilter("getitem", filterGetItem)

		// HTML/text filters
		pongo2.ReplaceFilter("striptags", filterStripTags)
		pongo2.RegisterFilter("plaintext", filterPlaintext)
		pongo2.RegisterFilter("summary", filterSummary)
		pongo2.RegisterFilter("pin_preview", filterPinPreview)
		pongo2.RegisterFilter("linebreaks", filterLinebreaks)
		pongo2.RegisterFilter("linebreaksbr", filterLinebreaksBR)

		// URL filters
		pongo2.RegisterFilter("urlencode", filterURLEncode)
		pongo2.RegisterFilter("absolute_url", filterAbsoluteURL)
		pongo2.RegisterFilter("domain", filterDomain)
		pongo2.RegisterFilter("external_url", filterExternalURL)

		// Theme/asset filters (per THEMES.md spec)
		pongo2.RegisterFilter("theme_asset", filterThemeAsset)
		pongo2.RegisterFilter("theme_asset_hashed", filterThemeAssetHashed)
		pongo2.RegisterFilter("asset_url", filterAssetURL)

		// ISO date format filter (per THEMES.md spec)
		pongo2.RegisterFilter("isoformat", filterISOFormat)

		// Font configuration filter
		pongo2.RegisterFilter("google_fonts_url", filterGoogleFontsURL)

		// String repeat filter for text output
		pongo2.RegisterFilter("repeat", filterRepeat)

		// Reading time filter
		pongo2.RegisterFilter("reading_time", filterReadingTime)

		// Excerpt filter
		pongo2.RegisterFilter("excerpt", filterExcerpt)
		pongo2.RegisterFilter("slides_reveal", filterSlidesReveal)

		// Type conversion filter
		pongo2.RegisterFilter("string", filterString)

		// Contribution data filter for Cal-Heatmap
		pongo2.RegisterFilter("contribution_data", filterContributionData)

		// Media detection filters
		pongo2.RegisterFilter("is_video", filterIsVideo)
		pongo2.RegisterFilter("media_url", filterMediaURL)
		pongo2.RegisterFilter("with_size", filterWithSize)
		pongo2.RegisterFilter("responsive_srcset", filterResponsiveSrcset)
		pongo2.RegisterFilter("editorial_card_image", filterEditorialCardImage)
		pongo2.RegisterFilter("original_media_url", filterOriginalMediaURL)
		pongo2.RegisterFilter("video_mime", filterVideoMIME)
		pongo2.RegisterFilter("poster_url", filterPosterURL)
	})
}

// filterRSSDate formats a date for RSS feeds.
// Format: "Mon, 02 Jan 2006 15:04:05 -0700"
func filterRSSDate(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, err := toTime(in)
	if err != nil {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(t.Format(time.RFC1123Z)), nil
}

// filterAtomDate formats a date for Atom feeds.
// Format: RFC3339 (e.g., "2006-01-02T15:04:05Z07:00")
func filterAtomDate(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, err := toTime(in)
	if err != nil {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(t.Format(time.RFC3339)), nil
}

// filterDate is a replacement for pongo2's built-in date filter that handles
// *time.Time pointers and string parsing in addition to time.Time values.
// Uses Go's time formatting (e.g., "2006-01-02", "January 2, 2006").
func filterDate(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, err := toTime(in)
	if err != nil {
		return pongo2.AsValue(""), nil
	}

	format := param.String()
	if format == "" {
		format = "2006-01-02"
	}

	return pongo2.AsValue(t.Format(format)), nil
}

// filterDateFormat formats a date using a custom format string.
// Uses Go's time formatting (e.g., "2006-01-02", "January 2, 2006").
func filterDateFormat(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, err := toTime(in)
	if err != nil {
		return pongo2.AsValue(""), nil
	}

	format := param.String()
	if format == "" {
		format = "2006-01-02"
	}

	return pongo2.AsValue(t.Format(format)), nil
}

// FormatHumanDate formats a time value for visible HTML text.
func FormatHumanDate(t time.Time) string {
	return t.Format(HumanDateFormat)
}

func filterHumanDate(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, err := toTime(in)
	if err != nil {
		return pongo2.AsValue(""), nil
	}

	return pongo2.AsValue(FormatHumanDate(t)), nil
}

// filterSlugify converts a string to a URL-safe slug.
// Converts to lowercase, replaces non-alphanumeric chars with hyphens,
// collapses multiple hyphens, and trims leading/trailing hyphens.
func filterSlugify(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	return pongo2.AsValue(models.Slugify(s)), nil
}

// filterTruncate truncates a string to a specified length with an ellipsis.
// Usage: {{ text|truncate:100 }} or {{ text|truncate:"50" }}
func filterTruncate(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	length := param.Integer()

	if length <= 0 {
		length = 80 // default
	}

	if len(s) <= length {
		return in, nil
	}

	// Truncate and add ellipsis
	truncated := s[:length]
	// Try to break at a word boundary
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > length/2 {
		truncated = truncated[:lastSpace]
	}

	return pongo2.AsValue(truncated + "..."), nil
}

// filterTruncateWords truncates a string to a specified number of words.
// Usage: {{ text|truncatewords:20 }}
func filterTruncateWords(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	wordCount := param.Integer()

	if wordCount <= 0 {
		wordCount = 20 // default
	}

	words := strings.Fields(s)
	if len(words) <= wordCount {
		return in, nil
	}

	return pongo2.AsValue(strings.Join(words[:wordCount], " ") + "..."), nil
}

// filterEndsWith checks if a string ends with a given suffix.
// Usage: {{ filename|endswith:".mp4" }}
func filterEndsWith(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	suffix := param.String()
	return pongo2.AsValue(strings.HasSuffix(s, suffix)), nil
}

// filterStartsWith checks if a string starts with a given prefix.
// Usage: {{ filename|startswith:"http" }}
func filterStartsWith(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	prefix := param.String()
	return pongo2.AsValue(strings.HasPrefix(s, prefix)), nil
}

// filterSplit splits a string by a delimiter and returns a slice.
// Usage: {{ "a.b.c"|split:"." }} -> ["a", "b", "c"]
func filterSplit(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	delim := param.String()
	if delim == "" {
		delim = " " // default to space
	}
	return pongo2.AsValue(strings.Split(s, delim)), nil
}

// filterReplace replaces occurrences of a substring with another.
// Usage: {{ url|replace:"https://," }} - replaces "https://" with ""
// Format: "old,new" where old is replaced with new
func filterReplace(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	paramStr := param.String()

	// Parse "old,new" format
	parts := strings.SplitN(paramStr, ",", 2)
	if len(parts) != 2 {
		// If no comma, treat param as string to remove (replace with "")
		return pongo2.AsValue(strings.ReplaceAll(s, paramStr, "")), nil
	}

	old := parts[0]
	replacement := parts[1]
	return pongo2.AsValue(strings.ReplaceAll(s, old, replacement)), nil
}

// filterDefaultIfNone returns a default value if the input is nil or empty.
// Usage: {{ value|default_if_none:"fallback" }}
func filterDefaultIfNone(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsNil() || (in.String() == "" && !in.IsBool()) {
		return param, nil
	}
	return in, nil
}

// filterLength returns the length of a string, slice, or map.
func filterLength(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.Len()), nil
}

// filterFirst returns the first element of a slice.
func filterFirst(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.Len() == 0 {
		return pongo2.AsValue(nil), nil
	}
	return in.Index(0), nil
}

// filterLast returns the last element of a slice.
func filterLast(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	length := in.Len()
	if length == 0 {
		return pongo2.AsValue(nil), nil
	}
	return in.Index(length - 1), nil
}

// filterGetItem gets an item from a map by key.
// Usage: {{ map|getitem:key }} or {{ map|getitem:"literal_key" }}
// Returns nil if the key doesn't exist.
func filterGetItem(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// Get the key as a string
	key := param.String()
	if key == "" {
		return pongo2.AsValue(nil), nil
	}

	// Try to get the underlying value as a map
	iface := in.Interface()
	if iface == nil {
		return pongo2.AsValue(nil), nil
	}

	// Handle map[string]string
	if m, ok := iface.(map[string]string); ok {
		if val, exists := m[key]; exists {
			return pongo2.AsValue(val), nil
		}
		return pongo2.AsValue(nil), nil
	}

	// Handle map[string]interface{}
	if m, ok := iface.(map[string]interface{}); ok {
		if val, exists := m[key]; exists {
			return pongo2.AsValue(val), nil
		}
		return pongo2.AsValue(nil), nil
	}

	// Handle map[string]any (same as interface{})
	if m, ok := iface.(map[string]any); ok {
		if val, exists := m[key]; exists {
			return pongo2.AsValue(val), nil
		}
		return pongo2.AsValue(nil), nil
	}

	return pongo2.AsValue(nil), nil
}

// filterJoin joins slice elements with a separator.
// Usage: {{ list|join:", " }}
func filterJoin(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.CanSlice() {
		return pongo2.AsValue(in.String()), nil
	}

	separator := param.String()
	if separator == "" {
		separator = ", "
	}

	var items []string
	for i := 0; i < in.Len(); i++ {
		items = append(items, in.Index(i).String())
	}

	return pongo2.AsValue(strings.Join(items, separator)), nil
}

// filterReverse reverses a slice or string.
func filterReverse(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// For strings, reverse the characters
	if !in.CanSlice() {
		s := in.String()
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return pongo2.AsValue(string(runes)), nil
	}

	// For slices, create a reversed copy
	length := in.Len()
	reversed := make([]interface{}, length)
	for i := 0; i < length; i++ {
		reversed[length-1-i] = in.Index(i).Interface()
	}
	return pongo2.AsValue(reversed), nil
}

// filterSort sorts a slice of comparable values.
func filterSort(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.CanSlice() || in.Len() == 0 {
		return in, nil
	}

	// Convert to string slice and sort
	var items []string
	for i := 0; i < in.Len(); i++ {
		items = append(items, in.Index(i).String())
	}

	// Simple string sort
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i] > items[j] {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	// Convert back to interface slice
	result := make([]interface{}, len(items))
	for i, s := range items {
		result[i] = s
	}

	return pongo2.AsValue(result), nil
}

// filterSelectAttr filters a slice of maps/structs to only include items
// where the specified attribute equals the given value.
// Usage: {{ items|selectattr:"key:value" }}
// Example: {{ webmentions|selectattr:"WMProperty:like-of" }}
func filterSelectAttr(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.CanSlice() {
		return in, nil
	}

	// Parse param as "key:value"
	paramStr := param.String()
	parts := strings.SplitN(paramStr, ":", 2)
	if len(parts) != 2 {
		return in, nil // Invalid param format, return original
	}
	key := parts[0]
	value := parts[1]

	var result []interface{}
	for i := 0; i < in.Len(); i++ {
		item := in.Index(i)
		// Try to get the attribute
		attr := getAttr(item, key)
		if attr != nil && attr.String() == value {
			result = append(result, item.Interface())
		}
	}

	return pongo2.AsValue(result), nil
}

// filterRejectAttr filters a slice of maps/structs to exclude items
// where the specified attribute equals the given value.
// Usage: {{ items|rejectattr:"key:value" }}
// Example: {{ webmentions|rejectattr:"WMProperty:like-of" }}
func filterRejectAttr(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.CanSlice() {
		return in, nil
	}

	// Parse param as "key:value"
	paramStr := param.String()
	parts := strings.SplitN(paramStr, ":", 2)
	if len(parts) != 2 {
		return in, nil // Invalid param format, return original
	}
	key := parts[0]
	value := parts[1]

	var result []interface{}
	for i := 0; i < in.Len(); i++ {
		item := in.Index(i)
		// Try to get the attribute
		attr := getAttr(item, key)
		if attr == nil || attr.String() != value {
			result = append(result, item.Interface())
		}
	}

	return pongo2.AsValue(result), nil
}

// getAttr gets an attribute from a pongo2.Value (works with maps and structs).
func getAttr(v *pongo2.Value, key string) *pongo2.Value {
	// Get the underlying interface
	iface := v.Interface()
	if iface == nil {
		return nil
	}

	rv := reflect.ValueOf(iface)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		// Try string key first
		mv := rv.MapIndex(reflect.ValueOf(key))
		if mv.IsValid() {
			return pongo2.AsValue(mv.Interface())
		}
	case reflect.Struct:
		fv := rv.FieldByName(key)
		if fv.IsValid() {
			return pongo2.AsValue(fv.Interface())
		}
	default:
		// Other types don't support field/key access
		return nil
	}

	return nil
}

// filterStripTags removes HTML tags from a string and cleans up entities.
func filterStripTags(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()

	// Remove HTML tags
	s = htmlTagRe.ReplaceAllString(s, "")

	// Decode common HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&apos;", "'")

	// Collapse multiple whitespace into single space
	s = whitespaceRe.ReplaceAllString(s, " ")

	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)

	return pongo2.AsValue(s), nil
}

// filterPlaintext converts HTML content to plain text with footnote-style link
// references. It decodes HTML entities, strips tags while preserving block
// structure, and formats links as "[text] [N]" with a references section.
// The output is marked as safe to prevent pongo2 from re-escaping entities.
//
// Usage: {{ post.content | plaintext }}
func filterPlaintext(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	if s == "" {
		return pongo2.AsSafeValue(""), nil
	}
	return pongo2.AsSafeValue(htmltotext.Convert(s)), nil
}

var (
	// Drop rendered heading self-links while leaving normal link text intact.
	summaryDropBlockRe  = regexp.MustCompile(`(?is)<(pre|code|script|style|figure|table|svg)\b[^>]*>.*?</\s*(pre|code|script|style|figure|table|svg)\s*>`)
	summaryDropAnchorRe = regexp.MustCompile(`(?is)<a\b[^>]*\bclass\s*=\s*(?:"(?:[^"]*\s)?heading-anchor(?:\s[^"]*)?"|'(?:[^']*\s)?heading-anchor(?:\s[^']*)?')[^>]*>.*?</a\s*>`)
	summaryWhitespaceRe = regexp.MustCompile(`\s+`)
	summaryPunctGapRe   = regexp.MustCompile(` +([.,;:!?)\]])`)
)

// filterSummary turns rendered HTML into a short, single-line, escaped
// summary suitable for cards and previews. Unlike plaintext it drops code,
// tables, and figures rather than flattening them, collapses all whitespace
// (so the result never breaks a surrounding HTML block when embedded in
// Markdown), and truncates to the optional character limit (default 300)
// with an ellipsis.
func filterSummary(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	if s == "" {
		return pongo2.AsSafeValue(""), nil
	}

	limit := 300
	if param != nil && !param.IsNil() {
		if n := param.Integer(); n > 0 {
			limit = n
		}
	}

	s = summaryDropBlockRe.ReplaceAllString(s, " ")
	s = summaryDropAnchorRe.ReplaceAllString(s, " ")
	s = htmlTagRe.ReplaceAllString(s, " ")
	s = html.UnescapeString(s)
	s = strings.TrimSpace(summaryWhitespaceRe.ReplaceAllString(s, " "))
	s = summaryPunctGapRe.ReplaceAllString(s, "$1")

	if runes := []rune(s); len(runes) > limit {
		s = strings.TrimRight(string(runes[:limit]), " ") + "…"
	}

	return pongo2.AsSafeValue(html.EscapeString(s)), nil
}

// filterLinebreaks converts newlines to <p> and <br> tags.
func filterLinebreaks(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	// Split by double newlines for paragraphs
	paragraphs := paragraphSplitRe.Split(s, -1)
	var result []string
	for _, p := range paragraphs {
		p = strings.TrimSpace(p)
		if p != "" {
			// Convert single newlines to <br>
			p = strings.ReplaceAll(p, "\n", "<br>")
			result = append(result, "<p>"+p+"</p>")
		}
	}
	return pongo2.AsValue(strings.Join(result, "\n")), nil
}

// filterLinebreaksBR converts newlines to <br> tags.
func filterLinebreaksBR(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	return pongo2.AsValue(strings.ReplaceAll(s, "\n", "<br>")), nil
}

// filterURLEncode URL-encodes a string.
func filterURLEncode(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	// Simple URL encoding for common characters
	s = strings.ReplaceAll(s, " ", "%20")
	s = strings.ReplaceAll(s, "&", "%26")
	s = strings.ReplaceAll(s, "=", "%3D")
	s = strings.ReplaceAll(s, "?", "%3F")
	s = strings.ReplaceAll(s, "#", "%23")
	return pongo2.AsValue(s), nil
}

// filterAbsoluteURL converts a relative URL to an absolute URL.
// Requires the site URL to be passed as the parameter.
// Usage: {{ post.href|absolute_url:config.url }}
func filterAbsoluteURL(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	path := in.String()
	baseURL := param.String()

	if baseURL == "" {
		return in, nil
	}

	// Remove trailing slash from base URL
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return pongo2.AsValue(baseURL + path), nil
}

// filterDomain extracts the hostname from a URL string.
// Usage: {{ "https://htmx.org/examples/foo/" | domain }} => "htmx.org"
func filterDomain(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	raw := in.String()
	if raw == "" {
		return pongo2.AsValue(""), nil
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" {
		return in, nil
	}
	return pongo2.AsValue(u.Hostname()), nil
}

// filterExternalURL returns only absolute HTTP(S) URLs suitable for hrefs.
// HTML escaping does not prevent unsafe schemes such as javascript:.
func filterExternalURL(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	raw := strings.TrimSpace(in.String())
	if raw == "" || strings.ContainsAny(raw, "\r\n\t<>") {
		return pongo2.AsValue(""), nil
	}
	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != editorialHTTPS) || u.Host == "" || u.User != nil {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(raw), nil
}

// filterThemeAsset returns a URL path for theme static assets.
// Usage: {{ 'css/main.css' | theme_asset }}
// Returns: /css/main.css (theme assets are copied to root of output)
func filterThemeAsset(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	path := in.String()
	if path == "" {
		return pongo2.AsValue(""), nil
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return pongo2.AsValue(path), nil
}

// filterAssetURL returns a URL path for project static assets.
// Usage: {{ 'images/logo.png' | asset_url }}
// Returns: /images/logo.png (project assets are at root of output)
func filterAssetURL(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	path := in.String()
	if path == "" {
		return pongo2.AsValue(""), nil
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return pongo2.AsValue(path), nil
}

// filterThemeAssetHashed returns a cache-busted URL for theme assets.
// Looks up content hash from global registry and appends to filename.
// Usage: {{ 'css/main.css' | theme_asset_hashed }}
// Returns: /css/main.abc12345.css (if hash registered), /css/main.css (fallback)
func filterThemeAssetHashed(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	path := in.String()
	if path == "" {
		return pongo2.AsValue(""), nil
	}

	// Get hash from global registry
	if hash, ok := GetAssetHash(path); ok && hash != "" {
		// Insert hash before the file extension
		// "css/main.css" → "css/main.abc12345.css"
		ext := filepath.Ext(path)
		base := strings.TrimSuffix(path, ext)
		path = base + "." + hash + ext
	}

	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return pongo2.AsValue(path), nil
}

// filterISOFormat formats a date in ISO 8601 format.
// Usage: {{ post.date | isoformat }}
// Returns: 2006-01-02T15:04:05Z07:00
func filterISOFormat(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	t, err := toTime(in)
	if err != nil {
		return pongo2.AsValue(""), nil
	}
	return pongo2.AsValue(t.Format(time.RFC3339)), nil
}

// filterGoogleFontsURL generates a Google Fonts CSS URL from a FontConfig.
// Usage: {{ config.theme.font | google_fonts_url }}
// Returns: https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&family=Fira+Code:wght@400;500;600;700&display=swap
func filterGoogleFontsURL(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsNil() {
		return pongo2.AsValue(""), nil
	}

	// Try to get google_fonts from the FontConfig struct
	var googleFonts []string

	// Check if it has a GoogleFonts field (FontConfig struct)
	iface := in.Interface()
	if fc, ok := iface.(interface{ GetGoogleFontsURL() string }); ok {
		return pongo2.AsValue(fc.GetGoogleFontsURL()), nil
	}

	// Try map access for template contexts where struct is converted to map
	if m, ok := iface.(map[string]interface{}); ok {
		if gf, exists := m["google_fonts"]; exists {
			switch fonts := gf.(type) {
			case []string:
				googleFonts = fonts
			case []interface{}:
				googleFonts = make([]string, 0, len(fonts))
				for _, f := range fonts {
					if s, ok := f.(string); ok {
						googleFonts = append(googleFonts, s)
					}
				}
			}
		}
	}

	if len(googleFonts) == 0 {
		return pongo2.AsValue(""), nil
	}

	// Build Google Fonts URL
	families := make([]string, 0, len(googleFonts))
	for _, font := range googleFonts {
		encoded := strings.ReplaceAll(font, " ", "+")
		families = append(families, "family="+encoded+":wght@400;500;600;700")
	}

	url := "https://fonts.googleapis.com/css2?" + strings.Join(families, "&") + "&display=swap"
	return pongo2.AsValue(url), nil
}

// filterRepeat repeats a string N times.
// The input is the count (length), and the parameter is the string to repeat.
// Usage: {{ post.title|length|repeat:"=" }}
func filterRepeat(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	count := in.Integer()
	char := param.String()

	if count <= 0 || char == "" {
		return pongo2.AsValue(""), nil
	}

	return pongo2.AsValue(strings.Repeat(char, count)), nil
}

// toTime attempts to convert a pongo2 value to a time.Time.
func toTime(in *pongo2.Value) (time.Time, error) {
	if in.IsNil() {
		return time.Time{}, fmt.Errorf("nil value")
	}

	// Check if it's already a time.Time
	if t, ok := in.Interface().(time.Time); ok {
		return t, nil
	}

	// Check if it's a *time.Time
	if t, ok := in.Interface().(*time.Time); ok {
		if t != nil {
			return *t, nil
		}
		return time.Time{}, fmt.Errorf("nil time pointer")
	}

	// Try to parse as string
	s := in.String()
	if s == "" {
		return time.Time{}, fmt.Errorf("empty string")
	}

	// Try common formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"January 2, 2006",
		"Jan 2, 2006",
		"02 Jan 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", s)
}

// filterReadingTime calculates estimated reading time for content.
// Assumes ~200 words per minute reading speed.
// Returns a string like "5 min read" or "< 1 min read"
func filterReadingTime(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	content := in.String()
	if content == "" {
		return pongo2.AsValue("< 1 min read"), nil
	}

	// Count words (simple: split by whitespace)
	words := strings.Fields(content)
	wordCount := len(words)

	// Assume 200 words per minute
	const wordsPerMinute = 200
	minutes := (wordCount + wordsPerMinute - 1) / wordsPerMinute // Round up

	if minutes < 1 {
		return pongo2.AsValue("< 1 min read"), nil
	}

	if minutes == 1 {
		return pongo2.AsValue("1 min read"), nil
	}

	return pongo2.AsValue(fmt.Sprintf("%d min read", minutes)), nil
}

// excerptConfig holds configuration for excerpt extraction.
type excerptConfig struct {
	maxParagraphs int
	maxChars      int
}

// defaultExcerptConfig returns the default excerpt configuration.
func defaultExcerptConfig() excerptConfig {
	return excerptConfig{
		maxParagraphs: 3,
		maxChars:      1500,
	}
}

// parseExcerptParams parses excerpt filter parameters.
// Supports: "paragraphs=N", "chars=N", or "paragraphs=N,chars=M"
func parseExcerptParams(param *pongo2.Value) excerptConfig {
	cfg := defaultExcerptConfig()

	if param == nil || param.String() == "" {
		return cfg
	}

	parts := strings.Split(param.String(), ",")
	for _, part := range parts {
		kv := strings.Split(strings.TrimSpace(part), "=")
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		switch key {
		case "paragraphs":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.maxParagraphs = n
			}
		case "chars":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.maxChars = n
			}
		}
	}

	return cfg
}

// removeAdmonitions strips admonition blocks from HTML.
// Admonitions are supplementary content that shouldn't appear in excerpts.
func removeAdmonitions(html string) string {
	// Remove standard admonition divs
	html = admonitionDivRe.ReplaceAllString(html, "")

	// Remove collapsible admonitions (details elements)
	html = admonitionDetailsRe.ReplaceAllString(html, "")

	return html
}

// truncateAtWordBoundary truncates text at a word boundary.
// Returns the truncated text with "..." appended if truncation occurred.
func truncateAtWordBoundary(text string, maxLen int, addEllipsis bool) string {
	if len(text) <= maxLen {
		return text
	}

	truncated := text[:maxLen]
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > maxLen/2 {
		truncated = truncated[:lastSpace]
	}

	if addEllipsis {
		return truncated + "..."
	}
	return truncated
}

// collectParagraphs extracts paragraphs from HTML up to the configured limits.
func collectParagraphs(html string, cfg excerptConfig) []string {
	matches := excerptParagraphRe.FindAllStringSubmatch(html, -1)

	if len(matches) == 0 {
		return nil
	}

	// Pre-allocate with expected capacity
	capacity := cfg.maxParagraphs
	if len(matches) < capacity {
		capacity = len(matches)
	}
	paragraphs := make([]string, 0, capacity)
	totalChars := 0

	for _, match := range matches {
		if len(paragraphs) >= cfg.maxParagraphs {
			break
		}

		text := cleanExcerptHTML(match[1])
		text = strings.TrimSpace(text)

		if text == "" {
			continue
		}

		pLen := len(text)

		// Check if adding this paragraph would exceed maxChars
		if totalChars+pLen > cfg.maxChars && len(paragraphs) > 0 {
			remaining := cfg.maxChars - totalChars
			if remaining > 50 {
				truncated := truncateAtWordBoundary(text, remaining, true)
				paragraphs = append(paragraphs, "<p>"+truncated+"</p>")
			}
			break
		}

		paragraphs = append(paragraphs, "<p>"+text+"</p>")
		totalChars += pLen

		if totalChars >= cfg.maxChars {
			break
		}
	}

	return paragraphs
}

// filterExcerpt extracts an excerpt from HTML content.
// Extracts the first N paragraphs or M characters (whichever is shorter).
// Usage: {{ post.article_html|excerpt }} - uses defaults (3 paragraphs or 1500 chars)
// Usage: {{ post.article_html|excerpt:"paragraphs=3" }}
// Usage: {{ post.article_html|excerpt:"chars=500" }}
// Usage: {{ post.article_html|excerpt:"paragraphs=2,chars=800" }}
func filterExcerpt(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	html := in.String()
	if html == "" {
		return pongo2.AsValue(""), nil
	}

	cfg := parseExcerptParams(param)
	html = removeAdmonitions(html)

	// Try to extract paragraphs
	paragraphs := collectParagraphs(html, cfg)

	// Fallback: no paragraphs found, strip all tags and truncate
	if len(paragraphs) == 0 {
		stripped := stripHTMLTagsHelper(html)
		if stripped == "" {
			return pongo2.AsValue(""), nil
		}
		result := truncateAtWordBoundary(stripped, cfg.maxChars, len(stripped) > cfg.maxChars)
		return pongo2.AsValue(result), nil
	}

	result := strings.Join(paragraphs, "\n")

	// Add ellipsis if we truncated (more content exists)
	allMatches := excerptParagraphRe.FindAllStringSubmatch(html, -1)
	if len(allMatches) > len(paragraphs) && !strings.HasSuffix(result, "...") {
		result += "\n<p>...</p>"
	}

	return pongo2.AsSafeValue(result), nil
}

// cleanExcerptHTML preserves inline formatting tags (code, strong, em, a, etc.)
// while removing block elements and cleaning up content for excerpts
func cleanExcerptHTML(s string) string {
	// Remove block-level tags but keep their content
	blockTags := []string{"div", "span", "section", "article", "header", "footer", "nav", "aside", "pre"}
	for _, tag := range blockTags {
		// Remove opening tags
		s = blockTagOpenRe[tag].ReplaceAllString(s, "")
		// Remove closing tags
		s = blockTagCloseRe[tag].ReplaceAllString(s, "")
	}

	// Collapse multiple whitespace
	s = whitespaceRe.ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}

// stripHTMLTagsHelper removes all HTML tags from a string (helper for filterExcerpt)
func stripHTMLTagsHelper(s string) string {
	s = htmlTagRe.ReplaceAllString(s, "")

	// Decode common HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&apos;", "'")

	// Collapse multiple whitespace
	s = whitespaceRe.ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}

type slideDeck struct {
	Horizontal []slideGroup
}

type slideGroup struct {
	Vertical []string
}

// filterSlidesReveal converts rendered article HTML into reveal.js sections.
// Horizontal slides start at h2 or hr boundaries. H3 starts a vertical child slide.
func filterSlidesReveal(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	htmlContent := in.String()
	if strings.TrimSpace(htmlContent) == "" {
		return pongo2.AsSafeValue(`<section><div class="slide-content"></div></section>`), nil
	}

	deck, err := buildRevealSlideDeck(htmlContent)
	if err != nil {
		return pongo2.AsSafeValue(`<section><div class="slide-content">` + htmlContent + `</div></section>`), nil
	}

	return pongo2.AsSafeValue(renderRevealSlideDeck(deck)), nil
}

func buildRevealSlideDeck(htmlContent string) (slideDeck, error) {
	context := &html.Node{Type: html.ElementNode, Data: "div", DataAtom: atom.Div}
	fragment, err := html.ParseFragment(strings.NewReader(htmlContent), context)
	if err != nil {
		return slideDeck{}, err
	}

	deck := slideDeck{}
	currentHorizontal := slideGroup{}
	currentVertical := strings.Builder{}
	hasVerticalContent := false
	hasHorizontal := false

	flushVertical := func() {
		content := strings.TrimSpace(currentVertical.String())
		if content == "" {
			currentVertical.Reset()
			hasVerticalContent = false
			return
		}
		currentHorizontal.Vertical = append(currentHorizontal.Vertical, content)
		currentVertical.Reset()
		hasVerticalContent = false
	}

	flushHorizontal := func() {
		flushVertical()
		if len(currentHorizontal.Vertical) == 0 {
			return
		}
		deck.Horizontal = append(deck.Horizontal, currentHorizontal)
		currentHorizontal = slideGroup{}
		hasHorizontal = false
	}

	for _, node := range fragment {
		if isHorizontalRuleSeparator(node) {
			flushHorizontal()
			continue
		}

		if isHorizontalHeadingSeparator(node) {
			if hasVerticalContent {
				flushHorizontal()
			}
		}

		if isVerticalSeparator(node) {
			if hasVerticalContent {
				flushVertical()
			}
			hasHorizontal = true
		}

		currentVertical.WriteString(renderNode(node))
		hasVerticalContent = true
		if !hasHorizontal {
			hasHorizontal = true
		}
	}

	flushHorizontal()

	if len(deck.Horizontal) == 0 {
		deck.Horizontal = append(deck.Horizontal, slideGroup{Vertical: []string{htmlContent}})
	}

	return deck, nil
}

func isHorizontalHeadingSeparator(node *html.Node) bool {
	if node == nil || node.Type != html.ElementNode {
		return false
	}
	return node.Data == "h2"
}

func isHorizontalRuleSeparator(node *html.Node) bool {
	if node == nil || node.Type != html.ElementNode {
		return false
	}
	return node.Data == "hr"
}

func isVerticalSeparator(node *html.Node) bool {
	if node == nil || node.Type != html.ElementNode {
		return false
	}
	return node.Data == "h3"
}

func renderRevealSlideDeck(deck slideDeck) string {
	var out strings.Builder
	for _, horizontal := range deck.Horizontal {
		if len(horizontal.Vertical) <= 1 {
			out.WriteString("<section>")
			if len(horizontal.Vertical) == 1 {
				out.WriteString(`<div class="slide-content">`)
				out.WriteString(horizontal.Vertical[0])
				out.WriteString(`</div>`)
			}
			out.WriteString("</section>")
			continue
		}

		out.WriteString("<section>")
		for _, vertical := range horizontal.Vertical {
			out.WriteString("<section>")
			out.WriteString(`<div class="slide-content">`)
			out.WriteString(vertical)
			out.WriteString(`</div></section>`)
		}
		out.WriteString("</section>")
	}
	return out.String()
}

func renderNode(node *html.Node) string {
	if node == nil {
		return ""
	}
	var out strings.Builder
	if err := html.Render(&out, node); err != nil {
		return ""
	}
	return out.String()
}

// filterString converts a value to a string representation.
// Usage: {{ number|string }} or {{ post.excerpt_paragraphs|string }}
func filterString(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in.String()), nil
}

// filterIsVideo returns true if the input string has a video file extension.
// Usage: {% if post.image|is_video %}...{% endif %}
func filterIsVideo(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(IsVideoURL(in.String())), nil
}

// filterMediaURL resolves a media URL from multiple fields.
// Returns the first non-empty value from the input and parameter, normalizing trusted media hosts to https.
// Usage: {{ post.image|media_url:post.video }}
func filterMediaURL(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// Try input first (primary field)
	if !in.IsNil() && in.String() != "" {
		return pongo2.AsValue(normalizeTrustedMediaURL(in.String())), nil
	}
	// Fall back to parameter (secondary field)
	if param != nil && !param.IsNil() && param.String() != "" {
		return pongo2.AsValue(normalizeTrustedMediaURL(param.String())), nil
	}
	return pongo2.AsValue(""), nil
}

// filterWithSize decorates a trusted URL with width/height query params.
// Usage: {{ media_src | with_size:"1200,675" }}  -- width + height
//
//	{{ media_src | with_size:"1200" }}       -- width-only (CDN preserves aspect ratio)
func filterWithSize(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsNil() {
		return in, nil
	}
	if param == nil || param.IsNil() {
		return in, nil
	}
	parts := strings.SplitN(param.String(), ",", 2)
	width, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || width <= 0 {
		return in, nil
	}
	height := 0
	if len(parts) == 2 {
		h, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err == nil && h > 0 {
			height = h
		}
	}
	return pongo2.AsValue(WithSize(in.String(), width, height)), nil
}

// filterEditorialCardImage avoids downloading full-resolution YouTube posters
// for tiny archive cards. The transformation is strictly host/path-scoped;
// every non-YouTube image URL retains its existing CDN responsive behavior.
func filterEditorialCardImage(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsNil() {
		return in, nil
	}
	thumbnail := strings.TrimSpace(in.String())
	u, err := url.Parse(thumbnail)
	if err != nil || u.Scheme != "https" {
		return in, nil
	}
	switch strings.ToLower(u.Hostname()) {
	case "i.ytimg.com", "img.youtube.com":
		// Do not rewrite arbitrary paths or query strings on other hosts.
		if strings.HasPrefix(u.Path, "/vi/") && strings.HasSuffix(u.Path, "/maxresdefault.jpg") {
			size := "mqdefault.jpg" // 320px: secondary link previews.
			if param != nil && param.IsTrue() {
				size = "sddefault.jpg" // 640px: featured lead previews.
			}
			u.Path = strings.TrimSuffix(u.Path, "maxresdefault.jpg") + size
			return pongo2.AsValue(u.String()), nil
		}
	}
	return in, nil
}

// filterResponsiveSrcset creates width-descriptor candidates for trusted
// remote media services from a comma-separated list of widths.
func filterResponsiveSrcset(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsNil() || param == nil || param.IsNil() {
		return pongo2.AsValue(""), nil
	}

	widths := make([]int, 0)
	for _, item := range strings.Split(param.String(), ",") {
		width, err := strconv.Atoi(strings.TrimSpace(item))
		if err == nil && width > 0 {
			widths = append(widths, width)
		}
	}
	return pongo2.AsValue(ResponsiveSrcset(in.String(), widths)), nil
}

func filterOriginalMediaURL(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsNil() {
		return in, nil
	}
	return pongo2.AsValue(OriginalMediaURL(in.String())), nil
}

// filterVideoMIME infers the video MIME type.
func filterVideoMIME(in, _ *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(VideoMIMEType(in.String())), nil
}

// filterPosterURL resolves a poster candidate map or derived fallback.
// Param is the video URL used for poster derivation.
func filterPosterURL(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if data := valueToStringMap(in); data != nil {
		mediaURL := ""
		if param != nil && !param.IsNil() {
			mediaURL = param.String()
		}
		return pongo2.AsValue(PosterURLFromMap(data, mediaURL)), nil
	}
	candidate := strings.TrimSpace(in.String())
	if candidate != "" {
		return pongo2.AsValue(normalizeTrustedMediaURL(candidate)), nil
	}
	return pongo2.AsValue(""), nil
}

func valueToStringMap(v *pongo2.Value) map[string]interface{} {
	if v == nil || v.IsNil() {
		return nil
	}
	if data, ok := v.Interface().(map[string]interface{}); ok {
		return data
	}
	if data, ok := v.Interface().(map[string]string); ok {
		result := make(map[string]interface{}, len(data))
		for k, val := range data {
			result[k] = val
		}
		return result
	}
	if post, ok := v.Interface().(*models.Post); ok {
		return GetPostMap(post)
	}
	return nil
}

// filterContributionData generates Cal-Heatmap compatible JSON data from posts.
// Takes a slice of posts and a year parameter, returns JSON array of {date, value} objects.
// Usage: {{ posts|contribution_data:2024 }}
// Returns: [{"date": "2024-01-01", "value": 2}, {"date": "2024-01-15", "value": 1}, ...]
func filterContributionData(in, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	year := param.Integer()
	if year == 0 {
		year = time.Now().Year()
	}

	// Count posts per day for the specified year
	postsByDate := make(map[string]int)

	// Iterate through the input (should be a slice of posts)
	if !in.CanSlice() {
		return pongo2.AsValue("[]"), nil
	}

	for i := 0; i < in.Len(); i++ {
		item := in.Index(i)
		if item.IsNil() {
			continue
		}

		// Try to get the date from the post
		var postDate time.Time
		var found bool

		// Try different ways to access the date
		iface := item.Interface()
		if post, ok := iface.(*models.Post); ok && post != nil && post.Date != nil {
			postDate = *post.Date
			found = true
		} else if post, ok := iface.(models.Post); ok && post.Date != nil {
			postDate = *post.Date
			found = true
		} else {
			// Try to get date field from map or struct (try both lowercase and capitalized)
			dateVal := getAttr(item, "date")
			if dateVal == nil {
				dateVal = getAttr(item, "Date")
			}
			if dateVal != nil {
				if t, err := toTime(dateVal); err == nil {
					postDate = t
					found = true
				}
			}
		}

		if !found {
			continue
		}

		// Only include posts from the specified year
		if postDate.Year() != year {
			continue
		}

		// Format date as YYYY-MM-DD
		dateStr := postDate.Format("2006-01-02")
		postsByDate[dateStr]++
	}

	// Convert to Cal-Heatmap data format
	type dataPoint struct {
		Date  string `json:"date"`
		Value int    `json:"value"`
	}

	data := make([]dataPoint, 0, len(postsByDate))
	for date, count := range postsByDate {
		data = append(data, dataPoint{Date: date, Value: count})
	}

	// Sort by date for consistent output
	for i := 0; i < len(data)-1; i++ {
		for j := i + 1; j < len(data); j++ {
			if data[i].Date > data[j].Date {
				data[i], data[j] = data[j], data[i]
			}
		}
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return pongo2.AsSafeValue("[]"), nil
	}

	// Return as safe value to prevent HTML escaping of JSON
	return pongo2.AsSafeValue(string(jsonBytes)), nil
}
