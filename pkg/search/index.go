package search

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	bleve "github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search/query"

	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/templates"
)

// Index wraps a bleve index with post-aware operations.
type Index struct {
	idx bleve.Index
	dir string
	mu  sync.RWMutex
}

const analyzerKeyword = "keyword"

// Result represents a single search hit with relevance score.
type Result struct {
	Post  *models.Post
	Doc   Document
	Score float64
}

// Document is the stored bleve document used to answer search results safely.
// It contains only fields that are allowed to be exposed by search.
type Document struct {
	Title       string    `json:"title"`
	Content     string    `json:"content,omitempty"`
	Description string    `json:"description,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Path        string    `json:"path"`
	Slug        string    `json:"slug,omitempty"`
	Href        string    `json:"href,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Published   bool      `json:"published"`
	Private     bool      `json:"private,omitempty"`
	WordCount   int       `json:"word_count,omitempty"`
	MediaURL    string    `json:"media_url,omitempty"`
	MediaType   string    `json:"media_type,omitempty"`
	PosterURL   string    `json:"poster_url,omitempty"`
	VideoMIME   string    `json:"video_mime,omitempty"`
	DocJSON     string    `json:"doc_json,omitempty"`
}

// QueryOptions configures a search query.
type QueryOptions struct {
	Limit     int
	Fuzzy     bool
	Fuzziness int
	Fields    []string // empty = all fields
	DateFrom  *time.Time
	DateTo    *time.Time
	Tags      []string
	Published *bool
}

const (
	indexDir  = "search.bleve"
	hashFile  = "search.hash"
	batchSize = 500
)

// DefaultDir returns the default index directory inside the cache dir.
func DefaultDir(cacheDir string) string {
	return filepath.Join(cacheDir, indexDir)
}

// NamedDir returns an index directory with a process-specific suffix.
// Use this when multiple processes may share the same cache directory
// (e.g., serve and search-server running simultaneously).
func NamedDir(cacheDir, name string) string {
	return filepath.Join(cacheDir, "search-"+name+".bleve")
}

// NamedHashFile returns a hash file path with a process-specific suffix.
func NamedHashFile(cacheDir, name string) string {
	return filepath.Join(cacheDir, "search-"+name+".hash")
}

// Open opens an existing bleve index. Returns an error if no index exists.
func Open(dir string) (*Index, error) {
	idx, err := bleve.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("open search index: %w", err)
	}
	return &Index{idx: idx, dir: dir}, nil
}

// Build creates a new bleve index from posts.
// If dir already exists, it is removed and rebuilt.
func Build(dir string, posts []*models.Post) (*Index, error) {
	if err := os.RemoveAll(dir); err != nil {
		return nil, fmt.Errorf("remove old index: %w", err)
	}

	m := buildMapping()
	idx, err := bleve.New(dir, m)
	if err != nil {
		return nil, fmt.Errorf("create search index: %w", err)
	}

	si := &Index{idx: idx, dir: dir}
	if err := si.indexPosts(posts); err != nil {
		idx.Close()
		return nil, err
	}

	// Index synonym definitions (best-effort — search works without them)
	if synErr := indexSynonyms(idx); synErr != nil {
		fmt.Fprintf(os.Stderr, "warning: synonym indexing failed: %v\n", synErr)
	}

	return si, nil
}

// BuildIfNeeded creates or opens an index, rebuilding only if the content hash changed.
func BuildIfNeeded(cacheDir string, posts []*models.Post) (*Index, error) {
	return BuildIfNeededNamed(cacheDir, "", posts)
}

// BuildIfNeededAt creates or opens an index at an explicit path, rebuilding only if the content hash changed.
func BuildIfNeededAt(dir, hashPath string, posts []*models.Post) (*Index, error) {
	currentHash := ContentHash(posts)
	storedHash := readHashFile(hashPath)

	if currentHash == storedHash {
		si, err := Open(dir)
		if err == nil {
			return si, nil
		}
		// Index corrupted or missing — fall through to rebuild
	}

	si, err := Build(dir, posts)
	if err != nil {
		return nil, err
	}

	if hashPath != "" {
		if mkErr := os.MkdirAll(filepath.Dir(hashPath), 0o755); mkErr == nil {
			_ = os.WriteFile(hashPath, []byte(currentHash), 0o600) //nolint:errcheck // best-effort cache
		}
	}
	return si, nil
}

func hashSearchVisiblePosts(posts []*models.Post) string {
	h := sha256.New()
	paths := make([]string, 0, len(posts))
	docs := make(map[string]Document, len(posts))

	for _, post := range posts {
		if post == nil || post.Skip || post.Draft {
			continue
		}
		paths = append(paths, post.Path)
		docs[post.Path] = toPostDoc(post)
	}

	sort.Strings(paths)
	for _, path := range paths {
		doc := docs[path]
		stored := doc
		stored.DocJSON = ""
		data, err := json.Marshal(stored)
		if err != nil {
			fmt.Fprintf(h, "%s:%v\n", path, stored)
			continue
		}
		fmt.Fprintf(h, "%s:%s\n", path, data)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// BuildIfNeededNamed creates or opens a named index, allowing multiple processes
// to maintain separate indexes in the same cache directory.
func BuildIfNeededNamed(cacheDir, name string, posts []*models.Post) (*Index, error) {
	var dir, hf string
	if name == "" {
		dir = DefaultDir(cacheDir)
		hf = filepath.Join(cacheDir, hashFile)
	} else {
		dir = NamedDir(cacheDir, name)
		hf = NamedHashFile(cacheDir, name)
	}
	return BuildIfNeededAt(dir, hf, posts)
}

// Close closes the index.
func (si *Index) Close() error {
	si.mu.Lock()
	defer si.mu.Unlock()
	return si.idx.Close()
}

// Search performs a ranked full-text search.
func (si *Index) Search(query string, opts QueryOptions, postsByPath map[string]*models.Post) ([]Result, error) {
	si.mu.RLock()
	defer si.mu.RUnlock()

	q := buildQuery(query, opts)

	req := bleve.NewSearchRequest(q)
	if opts.Limit > 0 {
		req.Size = opts.Limit
	} else {
		req.Size = 1000
	}
	req.Fields = []string{"path", "doc_json"}

	searchResult, err := si.idx.Search(req)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	results := make([]Result, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		doc := documentFromFields(hit.Fields)
		path, ok := hit.Fields["path"].(string)
		if !ok {
			path = doc.Path
		}
		var post *models.Post
		if postsByPath != nil {
			post = postsByPath[path]
		}
		if post == nil && doc.Path == "" {
			continue
		}
		results = append(results, Result{Post: post, Doc: doc, Score: hit.Score})
	}
	return results, nil
}

// indexPosts batch-indexes all posts into the bleve index.
// Draft and skipped posts are excluded entirely.
func (si *Index) indexPosts(posts []*models.Post) error {
	batch := si.idx.NewBatch()
	count := 0
	for _, post := range posts {
		if post.Skip || post.Draft {
			continue
		}
		doc := toPostDoc(post)
		if err := batch.Index(post.Path, doc); err != nil {
			return fmt.Errorf("index post %s: %w", post.Path, err)
		}
		count++
		if count%batchSize == 0 {
			if err := si.idx.Batch(batch); err != nil {
				return fmt.Errorf("batch index: %w", err)
			}
			batch = si.idx.NewBatch()
		}
	}
	if batch.Size() > 0 {
		if err := si.idx.Batch(batch); err != nil {
			return fmt.Errorf("batch index: %w", err)
		}
	}
	return nil
}

func toPostDoc(p *models.Post) Document {
	doc := Document{
		Path:      p.Path,
		Slug:      p.Slug,
		Href:      p.Href,
		Published: p.Published,
		Private:   p.Private,
	}
	if !p.Private {
		doc.Tags = p.Tags
	}
	if p.Title != nil && (!p.Private || explicitFrontmatterTitle(p)) {
		doc.Title = p.PlainTitle()
	} else if p.Private {
		doc.Title = ""
	}
	if p.Description != nil && (!p.Private || explicitFrontmatterDescription(p)) {
		doc.Description = *p.Description
	}
	if p.Date != nil {
		doc.Date = *p.Date
	}
	doc.MediaURL, doc.MediaType, doc.PosterURL, doc.VideoMIME = documentMedia(p)
	if p.Extra != nil {
		if wc, ok := p.Extra["word_count"].(int); ok {
			doc.WordCount = wc
		}
	}
	// Private posts expose only metadata in search. Their body content remains hidden
	// until the user visits the page and decrypts it.
	if !p.Private {
		doc.Content = p.Content
	}
	if p.Private {
		sanitizePrivateDocument(&doc, explicitFrontmatterDescription(p))
	}
	stored := doc
	stored.DocJSON = ""
	if data, err := json.Marshal(stored); err == nil {
		doc.DocJSON = string(data)
	}
	return doc
}

func sanitizePrivateDocument(doc *Document, allowDescription bool) {
	if doc == nil || !doc.Private {
		return
	}
	doc.Content = ""
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

func explicitFrontmatterTitle(post *models.Post) bool {
	return post != nil && post.Has("_title_explicit")
}

func explicitFrontmatterDescription(post *models.Post) bool {
	return post != nil && post.Has("_description_explicit")
}

func documentMedia(post *models.Post) (mediaURL, mediaType, posterURL, videoMIME string) {
	if post == nil || post.Extra == nil || post.Private {
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

	return templates.WithSize(mediaURL, 320, 180), "image", "", ""
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

func documentFromFields(fields map[string]interface{}) Document {
	if fields == nil {
		return Document{}
	}
	if raw, ok := fields["doc_json"].(string); ok && strings.TrimSpace(raw) != "" {
		var doc Document
		if err := json.Unmarshal([]byte(raw), &doc); err == nil {
			return doc
		}
	}
	var doc Document
	if path, ok := fields["path"].(string); ok {
		doc.Path = path
	}
	return doc
}

func buildMapping() mapping.IndexMapping {
	im := bleve.NewIndexMapping()
	im.DefaultAnalyzer = "en"

	// Configure synonym source: link the "english" source to the "wordnet" collection
	// using the "en" analyzer for term normalization.
	err := im.AddSynonymSource(synonymSourceName, map[string]interface{}{
		"collection": synonymCollection,
		"analyzer":   synonymAnalyzer,
	})
	if err != nil {
		// Non-fatal: search works without synonyms
		fmt.Fprintf(os.Stderr, "warning: synonym source setup failed: %v\n", err)
	}

	dm := bleve.NewDocumentMapping()

	title := bleve.NewTextFieldMapping()
	title.Analyzer = "en"
	title.Store = false
	title.IncludeTermVectors = true
	title.SynonymSource = synonymSourceName
	dm.AddFieldMappingsAt("title", title)

	content := bleve.NewTextFieldMapping()
	content.Analyzer = "en"
	content.Store = false
	content.IncludeTermVectors = true
	content.SynonymSource = synonymSourceName
	dm.AddFieldMappingsAt("content", content)

	desc := bleve.NewTextFieldMapping()
	desc.Analyzer = "en"
	desc.Store = false
	desc.SynonymSource = synonymSourceName
	dm.AddFieldMappingsAt("description", desc)

	tags := bleve.NewTextFieldMapping()
	tags.Analyzer = analyzerKeyword
	tags.Store = false
	tags.DocValues = true
	dm.AddFieldMappingsAt("tags", tags)

	path := bleve.NewTextFieldMapping()
	path.Analyzer = analyzerKeyword
	path.Store = true
	path.Index = false
	dm.AddFieldMappingsAt("path", path)

	slug := bleve.NewTextFieldMapping()
	slug.Analyzer = analyzerKeyword
	slug.Store = false
	dm.AddFieldMappingsAt("slug", slug)

	date := bleve.NewDateTimeFieldMapping()
	date.Store = false
	date.DocValues = true
	dm.AddFieldMappingsAt("date", date)

	published := mapping.NewBooleanFieldMapping()
	published.Store = false
	published.DocValues = true
	dm.AddFieldMappingsAt("published", published)

	wordCount := bleve.NewNumericFieldMapping()
	wordCount.Store = false
	dm.AddFieldMappingsAt("word_count", wordCount)

	docJSON := bleve.NewTextFieldMapping()
	docJSON.Store = true
	docJSON.Index = false
	dm.AddFieldMappingsAt("doc_json", docJSON)

	im.DefaultMapping = dm
	return im
}

// synonymFields are the fields with SynonymSource configured.
// MatchQuery must target specific fields for synonym expansion to work;
// the default "_all" field has no synonym source.
var synonymFields = []string{"title", "content", "description"}

func buildQuery(queryStr string, opts QueryOptions) query.Query {
	var textQuery query.Query
	if opts.Fuzzy {
		// Build per-field fuzzy queries so synonym expansion triggers
		fuzziness := opts.Fuzziness
		if fuzziness <= 0 {
			fuzziness = 1
		}
		dq := query.NewDisjunctionQuery(nil)
		for _, field := range synonymFields {
			mq := query.NewMatchQuery(queryStr)
			mq.SetFuzziness(fuzziness)
			mq.SetField(field)
			dq.AddQuery(mq)
		}
		textQuery = dq
	} else {
		// Build per-field match queries so synonym expansion triggers
		dq := query.NewDisjunctionQuery(nil)
		for _, field := range synonymFields {
			mq := query.NewMatchQuery(queryStr)
			mq.SetField(field)
			dq.AddQuery(mq)
		}
		textQuery = dq
	}

	// If no additional filters, return the text query directly
	if opts.DateFrom == nil && opts.DateTo == nil && len(opts.Tags) == 0 && opts.Published == nil {
		return textQuery
	}

	// Build boolean query with filters
	bq := query.NewBooleanQuery(nil, nil, nil)
	bq.AddMust(textQuery)

	if opts.DateFrom != nil || opts.DateTo != nil {
		var from, to time.Time
		if opts.DateFrom != nil {
			from = *opts.DateFrom
		}
		if opts.DateTo != nil {
			to = *opts.DateTo
		}
		dq := query.NewDateRangeQuery(from, to)
		dq.SetField("date")
		bq.AddMust(dq)
	}

	for _, tag := range opts.Tags {
		tq := query.NewTermQuery(tag)
		tq.SetField("tags")
		bq.AddMust(tq)
	}

	if opts.Published != nil {
		pq := query.NewBoolFieldQuery(*opts.Published)
		pq.SetField("published")
		bq.AddMust(pq)
	}

	return bq
}

// ContentHash computes a hash of all post paths and content for cache invalidation.
func ContentHash(posts []*models.Post) string {
	return hashSearchVisiblePosts(posts)
}

func readHashFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// PostsByPath builds a lookup map from posts for result resolution.
func PostsByPath(posts []*models.Post) map[string]*models.Post {
	m := make(map[string]*models.Post, len(posts))
	for _, p := range posts {
		if p.Draft || p.Skip {
			continue
		}
		m[p.Path] = p
	}
	return m
}
