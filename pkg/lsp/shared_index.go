package lsp

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/WaylonWalker/markata/pkg/models"
	"gopkg.in/yaml.v3"
)

// IndexProvider returns an immutable workspace snapshot and its revision.
// The returned index and its PostInfo values must not be mutated after publication.
type IndexProvider func() (*Index, uint64)

// LoadedPostInfo converts an already loaded post without reading its source again.
func LoadedPostInfo(root string, post *models.Post) *PostInfo {
	metadata := make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(post.RawFrontmatter), &metadata); err != nil {
		metadata = make(map[string]interface{})
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	if aliases := post.Get("aliases"); aliases != nil {
		metadata["aliases"] = aliases
	}
	path := post.Path
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	title := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if post.Title != nil {
		title = *post.Title
	}
	description := extractExcerpt(post.Content, 200)
	if post.Description != nil {
		description = *post.Description
	}
	return &PostInfo{URI: pathToURI(path), Path: path, Slug: post.Slug, Title: title, Description: description, Metadata: metadata, Aliases: extractAliases(metadata), Wikilinks: findWikilinks(post.Content)}
}

// NewWorkspaceIndex publishes parsed workspace information for browser sessions.
func NewWorkspaceIndex(logger *log.Logger, root string, posts []*PostInfo) *Index {
	index := NewIndex(logger)
	for _, post := range posts {
		index.posts[post.Slug] = post
		index.uriToSlug[post.URI] = post.Slug
	}
	for _, post := range posts {
		for _, alias := range post.Aliases {
			key := strings.ToLower(alias)
			if _, exists := index.posts[key]; !exists {
				index.posts[key] = post
			}
		}
	}
	index.indexBlogrollMentions(root)
	index.indexFromPostsMentions(root)
	return index
}

// Clone copies index maps while sharing immutable parsed workspace values.
// Mutations through Update and Remove affect only the new index.
func (idx *Index) Clone() *Index {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	idx.mentionMu.RLock()
	defer idx.mentionMu.RUnlock()
	cloned := NewIndex(idx.logger)
	for key, value := range idx.posts {
		cloned.posts[key] = value
	}
	for key, value := range idx.uriToSlug {
		cloned.uriToSlug[key] = value
	}
	for key, value := range idx.mentions {
		cloned.mentions[key] = value
	}
	return cloned
}

func (s *Server) refreshSharedIndex() {
	if s.sharedIndex == nil {
		return
	}
	index, revision := s.sharedIndex()
	if index == nil || (s.sharedLoaded && revision == s.sharedRevision) {
		return
	}
	next := index.Clone()
	s.docMu.RLock()
	for uri, document := range s.documents {
		if err := next.Update(uri, document.Content); err != nil {
			s.logger.Printf("index draft: %v", err)
		}
	}
	s.docMu.RUnlock()
	s.index = next
	s.sharedRevision = revision
	s.sharedLoaded = true
}
