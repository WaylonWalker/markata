package searchapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/search"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// Handler serves the bleve search API.
// The index is built lazily on first search and reused across requests.
// It is only rebuilt when the search-visible post data changes.
type Handler struct {
	posts       []*models.Post
	postsByPath map[string]*models.Post
	cacheDir    string
	config      Config
	mu          sync.RWMutex
	idx         *search.Index
	postsHash   string
	searchSem   chan struct{}
}

// Config controls search API behavior.
type Config struct {
	DefaultLimit int
	MaxLimit     int
	DefaultFuzzy bool
	CORSOrigins  []string
	IndexName    string // process-specific index name to avoid conflicts
	IndexDir     string
	HashPath     string
	ReadOnly     bool
	Rebuild      bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultFuzzy: false,
		CORSOrigins:  []string{"*"},
	}
}

// SearchResponse is the JSON response envelope.
type SearchResponse struct {
	Query   string         `json:"query"`
	Total   int            `json:"total"`
	Fuzzy   bool           `json:"fuzzy"`
	Limit   int            `json:"limit"`
	Results []SearchResult `json:"results"`
}

// SearchResult is a single search hit in the API response.
type SearchResult struct {
	Title       string   `json:"title"`
	Path        string   `json:"path"`
	Slug        string   `json:"slug"`
	Href        string   `json:"href"`
	MediaURL    string   `json:"media_url,omitempty"`
	MediaType   string   `json:"media_type,omitempty"`
	PosterURL   string   `json:"poster_url,omitempty"`
	VideoMIME   string   `json:"video_mime,omitempty"`
	Description string   `json:"description,omitempty"`
	Date        string   `json:"date,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Score       float64  `json:"score"`
	WordCount   int      `json:"word_count,omitempty"`
	ReadTime    string   `json:"read_time,omitempty"`
	Private     bool     `json:"private,omitempty"`
}

// NewHandler creates a new search API handler.
// The bleve index is built lazily on the first search request,
// not at construction time, to avoid blocking server startup.
func NewHandler(posts []*models.Post, cacheDir string, cfg Config) *Handler {
	// Limit concurrent searches to prevent resource exhaustion.
	// Default to 4 concurrent searches; large sites (3000+ posts) with
	// bleve can use significant memory per query.
	maxConcurrent := 4
	return &Handler{
		posts:       posts,
		postsByPath: search.PostsByPath(posts),
		cacheDir:    cacheDir,
		config:      cfg,
		postsHash:   postsFingerprint(posts),
		searchSem:   make(chan struct{}, maxConcurrent),
	}
}

// NewReadOnlyHandler creates a handler that serves from an existing bleve index
// without loading or rebuilding site content.
func NewReadOnlyHandler(indexDir string, cfg Config) *Handler {
	cfg.ReadOnly = true
	cfg.IndexDir = indexDir
	maxConcurrent := 4
	return &Handler{
		cacheDir:  filepath.Dir(indexDir),
		config:    cfg,
		searchSem: make(chan struct{}, maxConcurrent),
	}
}

// UpdatePosts replaces the post list (e.g., after a rebuild).
// Only rebuilds the index if the post set actually changed.
func (h *Handler) UpdatePosts(posts []*models.Post) {
	newHash := postsFingerprint(posts)

	h.mu.Lock()
	defer h.mu.Unlock()

	if newHash == h.postsHash {
		return // nothing changed
	}

	h.posts = posts
	h.postsByPath = search.PostsByPath(posts)
	h.postsHash = newHash

	// Close old index so it rebuilds on next search
	if h.idx != nil {
		h.idx.Close()
		h.idx = nil
	}
}

// Warm ensures the current index is opened or built before the handler starts serving queries.
func (h *Handler) Warm() error {
	_, _, err := h.getOrBuildIndex()
	return err
}

// Close releases the underlying search index.
func (h *Handler) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.idx != nil {
		h.idx.Close()
		h.idx = nil
	}
}

// ServeHTTP handles GET requests for search.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS
	origin := r.Header.Get("Origin")
	if origin != "" && h.corsAllowed(origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	queryStr := r.URL.Query().Get("q")
	if queryStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing required parameter: q"})
		return
	}

	opts, limit, fuzzy := h.parseQueryOptions(r)

	// Acquire search semaphore to limit concurrent searches
	select {
	case h.searchSem <- struct{}{}:
		defer func() { <-h.searchSem }()
	default:
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "too many concurrent searches, try again"})
		return
	}

	idx, postsByPath, err := h.getOrBuildIndex()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "search index unavailable"})
		return
	}

	results, err := idx.Search(queryStr, opts, postsByPath)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "search failed"})
		return
	}

	writeJSON(w, http.StatusOK, buildResponse(queryStr, results, fuzzy, limit))
}

// parseQueryOptions extracts search parameters from the HTTP request.
func (h *Handler) parseQueryOptions(r *http.Request) (search.QueryOptions, int, bool) {
	q := r.URL.Query()

	fuzzy := h.config.DefaultFuzzy
	if f := q.Get("fuzzy"); f != "" {
		fuzzy = f == "true" || f == "1"
	}

	limit := h.config.DefaultLimit
	if l := q.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > h.config.MaxLimit {
		limit = h.config.MaxLimit
	}

	opts := search.QueryOptions{
		Limit: limit,
		Fuzzy: fuzzy,
	}

	if tags := q.Get("tags"); tags != "" {
		opts.Tags = parseTags(tags)
	}
	if from := q.Get("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			opts.DateFrom = &t
		}
	}
	if to := q.Get("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			opts.DateTo = &t
		}
	}

	published := true
	opts.Published = &published

	return opts, limit, fuzzy
}

func parseTags(raw string) []string {
	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		tag := strings.TrimSpace(part)
		if tag == "" {
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

func buildResponse(queryStr string, results []search.Result, fuzzy bool, limit int) SearchResponse {
	resp := SearchResponse{
		Query:   queryStr,
		Total:   len(results),
		Fuzzy:   fuzzy,
		Limit:   limit,
		Results: make([]SearchResult, len(results)),
	}

	for i := range results {
		hit := &results[i]
		doc := hit.Doc
		if hit.Post != nil {
			doc = search.Document{}
			if !hit.Post.Private || hit.Post.Has("_title_explicit") {
				doc.Title = hit.Post.PlainTitle()
			}
			doc.Path = hit.Post.Path
			doc.Slug = hit.Post.Slug
			doc.Href = hit.Post.Href
			if !hit.Post.Private {
				doc.Tags = append([]string(nil), hit.Post.Tags...)
			}
			doc.Private = hit.Post.Private
			if hit.Post.Description != nil && (!hit.Post.Private || explicitFrontmatterDescription(hit.Post)) {
				doc.Description = *hit.Post.Description
			}
			if hit.Post.Date != nil {
				doc.Date = *hit.Post.Date
			}
			if hit.Post.Extra != nil {
				if wc, ok := hit.Post.Extra["word_count"].(int); ok {
					doc.WordCount = wc
				}
			}
			doc.MediaURL, doc.MediaType, doc.PosterURL, doc.VideoMIME = searchMedia(hit.Post)
			if hit.Post.Private {
				sanitizePrivateDocument(&doc, explicitFrontmatterDescription(hit.Post))
			}
		}
		if doc.Private {
			sanitizePrivateDocument(&doc, doc.Description != "")
		}

		sr := SearchResult{
			Title:   doc.Title,
			Path:    doc.Path,
			Slug:    doc.Slug,
			Href:    doc.Href,
			Tags:    append([]string(nil), doc.Tags...),
			Score:   hit.Score,
			Private: doc.Private,
		}
		sr.Description = doc.Description
		sr.MediaURL = doc.MediaURL
		sr.MediaType = doc.MediaType
		sr.PosterURL = doc.PosterURL
		sr.VideoMIME = doc.VideoMIME
		if !doc.Date.IsZero() {
			sr.Date = doc.Date.Format(time.RFC3339)
		}
		if doc.WordCount > 0 {
			sr.WordCount = doc.WordCount
			sr.ReadTime = readTime(doc.WordCount)
		}
		resp.Results[i] = sr
	}
	return resp
}

func sanitizePrivateDocument(doc *search.Document, allowDescription bool) {
	if doc == nil || !doc.Private {
		return
	}
	if !allowDescription {
		doc.Description = ""
	}
	doc.Tags = nil
	doc.WordCount = 0
	doc.MediaURL = ""
	doc.MediaType = ""
	doc.PosterURL = ""
	doc.VideoMIME = ""
}

func searchMedia(post *models.Post) (mediaURL, mediaType, posterURL, videoMIME string) {
	if post == nil || post.Extra == nil {
		return "", "", "", ""
	}
	if post.Private {
		return "", "", "", ""
	}

	imageURL := firstExtraString(post.Extra, "image", "cover", "cover_image", "og_image")
	videoURL := firstExtraString(post.Extra, "video")
	mediaURL = imageURL
	if mediaURL == "" {
		mediaURL = videoURL
	}
	if mediaURL == "" {
		return "", "", "", ""
	}

	if templates.IsVideoURL(mediaURL) {
		mediaType = "video"
		videoMIME = templates.VideoMIMEType(mediaURL)
		posterURL = templates.PosterURLFromMap(post.Extra, mediaURL)
		mediaURL = templates.WithSize(mediaURL, 320, 180)
		if posterURL != "" {
			posterURL = templates.WithSize(posterURL, 320, 180)
		}
		return mediaURL, mediaType, posterURL, videoMIME
	}

	mediaType = "image"
	mediaURL = templates.WithSize(mediaURL, 320, 180)
	return mediaURL, mediaType, "", ""
}

func firstExtraString(extra map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		value, ok := extra[key].(string)
		if ok && strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func explicitFrontmatterDescription(post *models.Post) bool {
	if post == nil || post.Description == nil {
		return false
	}
	if post.Extra == nil {
		return false
	}
	value, ok := post.Extra["description"]
	if !ok {
		return false
	}
	text, ok := value.(string)
	return ok && strings.TrimSpace(text) != ""
}

// getOrBuildIndex returns the current index, building it if needed.
// The index is reused across requests and only rebuilt when posts change.
func (h *Handler) getOrBuildIndex() (*search.Index, map[string]*models.Post, error) {
	// Fast path: index already exists
	h.mu.RLock()
	if h.idx != nil {
		idx := h.idx
		postsByPath := h.postsByPath
		h.mu.RUnlock()
		return idx, postsByPath, nil
	}
	h.mu.RUnlock()

	// Slow path: build index
	h.mu.Lock()
	defer h.mu.Unlock()

	// Double-check after acquiring write lock
	if h.idx != nil {
		return h.idx, h.postsByPath, nil
	}

	idx, err := h.openOrBuildIndex()
	if err != nil {
		return nil, nil, err
	}
	h.idx = idx
	return h.idx, h.postsByPath, nil
}

func (h *Handler) openOrBuildIndex() (*search.Index, error) {
	indexDir := h.config.IndexDir
	hashPath := h.config.HashPath
	if indexDir == "" {
		if h.config.IndexName == "" {
			indexDir = search.DefaultDir(h.cacheDir)
			hashPath = filepath.Join(h.cacheDir, "search.hash")
		} else {
			indexDir = search.NamedDir(h.cacheDir, h.config.IndexName)
			hashPath = search.NamedHashFile(h.cacheDir, h.config.IndexName)
		}
	}
	if hashPath == "" && !h.config.ReadOnly {
		hashPath = filepath.Join(filepath.Dir(indexDir), filepath.Base(indexDir)+".hash")
	}

	if h.config.ReadOnly {
		return search.Open(indexDir)
	}

	if h.config.Rebuild {
		return search.Build(indexDir, h.posts)
	}

	return search.BuildIfNeededAt(indexDir, hashPath, h.posts)
}

func (h *Handler) corsAllowed(origin string) bool {
	for _, allowed := range h.config.CORSOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

func readTime(words int) string {
	minutes := words / 200
	if minutes < 1 {
		return "1 min"
	}
	return fmt.Sprintf("%d min", minutes)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	//nolint:errcheck // HTTP response writer errors are handled by the framework
	enc.Encode(v)
}

// postsFingerprint creates a cheap hash of the post set for change detection.
// It matches the search-visible content used to build the bleve index.
func postsFingerprint(posts []*models.Post) string {
	return search.ContentHash(posts)
}
