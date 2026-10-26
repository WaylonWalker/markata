package plugins

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"

	"github.com/WaylonWalker/markata/pkg/contentindex"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/sourcegit"
)

// ContentIndexPlugin writes the opt-in Markata Content Index after feed
// membership has been resolved.
type ContentIndexPlugin struct {
	initialSource    sourcegit.State
	initialSourceErr error
	capturedSource   bool
}

func NewContentIndexPlugin() *ContentIndexPlugin { return &ContentIndexPlugin{} }
func (p *ContentIndexPlugin) Name() string       { return "content_index" }
func (p *ContentIndexPlugin) Priority(stage lifecycle.Stage) int {
	if stage == lifecycle.StageWrite {
		return lifecycle.PriorityLast
	}
	return lifecycle.PriorityDefault
}

// Configure captures source state before content discovery starts. The final
// write compares this snapshot with a second snapshot before publishing.
func (p *ContentIndexPlugin) Configure(m *lifecycle.Manager) error {
	cfg, ok, err := contentIndexConfig(m.Config().Extra)
	if err != nil {
		return err
	}
	if !ok || !cfg.Enabled {
		return nil
	}
	p.initialSource, p.initialSourceErr = sourcegit.ReadWithOptions(context.Background(), m.Config().ContentDir, sourceSnapshotOptions(m.Config()))
	p.capturedSource = true
	return nil
}

//nolint:gocyclo // The writer coordinates the complete build-time artifact pipeline.
func (p *ContentIndexPlugin) Write(m *lifecycle.Manager) (err error) {
	defer func() {
		if err != nil {
			err = &contentIndexWriteError{err: err}
		}
	}()
	cfg, ok, configErr := contentIndexConfig(m.Config().Extra)
	if configErr != nil {
		return configErr
	}
	if !ok || !cfg.Enabled {
		return nil
	}
	schemaVersion := cfg.SchemaVersion
	if schemaVersion == 0 {
		schemaVersion = contentindex.CurrentVersion
	}
	if schemaVersion != 1 && schemaVersion != contentindex.CurrentVersion {
		return fmt.Errorf("content_index.schema_version %d is unsupported; latest is %d", schemaVersion, contentindex.CurrentVersion)
	}
	stateBefore, beforeErr := p.initialSource, p.initialSourceErr
	if !p.capturedSource {
		stateBefore, beforeErr = sourcegit.ReadWithOptions(context.Background(), m.Config().ContentDir, sourceSnapshotOptions(m.Config()))
	}

	posts := m.Posts()
	posts = slices.DeleteFunc(posts, func(post *models.Post) bool { return post == nil })
	sort.SliceStable(posts, func(i, j int) bool { return posts[i].Path < posts[j].Path })
	eligible := make(map[string]bool, len(posts))
	for _, post := range posts {
		if post != nil && !post.Draft && !post.Skip && (schemaVersion == contentindex.CurrentVersion || !post.Private) {
			eligible[post.Path] = true
		}
	}
	feeds := make(map[string][]string)
	for _, feed := range m.Feeds() {
		if feed == nil {
			continue
		}
		for _, post := range feed.Posts {
			if post != nil && eligible[post.Path] && (!post.Private || feed.IncludePrivate) {
				feeds[post.Path] = append(feeds[post.Path], feed.Name)
			}
		}
	}
	for path := range feeds {
		sort.Strings(feeds[path])
	}

	docs := make([]contentindex.Document, 0, len(eligible))
	for _, post := range posts {
		if post == nil || !eligible[post.Path] {
			continue
		}
		docs = append(docs, documentFromPost(post, feeds[post.Path]))
	}
	version := "dev"
	if v, ok := m.Config().Extra["markata_version"].(string); ok && v != "" {
		version = v
	}
	scope := contentindex.PublicMetadataScope
	if schemaVersion == 1 {
		scope = contentindex.PublicScope
	}
	index := contentindex.Index{Schema: contentindex.Schema, SchemaVersion: schemaVersion, Scope: scope, Generator: contentindex.Generator{Name: contentindex.GeneratorName, Version: version}, DocumentCount: len(docs), Documents: docs}
	stateAfter, afterErr := sourcegit.ReadWithOptions(context.Background(), m.Config().ContentDir, sourceSnapshotOptions(m.Config()))
	if beforeErr == nil && afterErr == nil {
		if !stateBefore.Equal(stateAfter) {
			return fmt.Errorf("source Git state changed during Content Index build")
		}
		index.Source.Commit = stateAfter.Commit
		index.Source.Dirty = stateAfter.Dirty
	}
	data, err := contentindex.Marshal(index)
	if err != nil {
		return fmt.Errorf("marshal content index: %w", err)
	}
	destination, err := contentIndexDestination(m.Config().OutputDir, cfg.Output)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return fmt.Errorf("create content index directory: %w", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(destination), ".content-index-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary content index: %w", err)
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0o644); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("set content index permissions: %w", err)
	}
	if _, err := temporary.Write(data); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("write temporary content index: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary content index: %w", err)
	}
	if err := replaceContentIndex(temporaryName, destination); err != nil {
		return fmt.Errorf("write content index: %w", err)
	}
	return nil
}

func sourceSnapshotOptions(config *lifecycle.Config) sourcegit.ReadOptions {
	// Git-ignore filtering defaults to enabled. The loader copies the effective
	// value into Extra for plugins, so only an explicit false enables the more
	// expensive ignored-input fingerprint.
	useGitignore := true
	if value, ok := config.Extra["use_gitignore"].(bool); ok {
		useGitignore = value
	}
	return sourcegit.ReadOptions{IncludeIgnoredContent: !useGitignore}
}

func contentIndexDestination(outputDir, configured string) (string, error) {
	destination := configured
	if destination == "" {
		destination = "content-index.json"
	}
	if filepath.IsAbs(destination) {
		return destination, nil
	}
	clean := filepath.Clean(destination)
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("content_index.output must stay within output_dir: %q", configured)
	}
	return filepath.Join(outputDir, clean), nil
}

func replaceContentIndex(source, destination string) error {
	if runtime.GOOS == tailwindOSWindows {
		if err := os.Remove(destination); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return os.Rename(source, destination)
}

type contentIndexWriteError struct{ err error }

func (e *contentIndexWriteError) Error() string    { return e.err.Error() }
func (e *contentIndexWriteError) Unwrap() error    { return e.err }
func (e *contentIndexWriteError) IsCritical() bool { return true }

func documentFromPost(post *models.Post, feedNames []string) contentindex.Document {
	document := contentindex.Document{Path: post.Path, Slug: post.Slug, Href: post.Href, Title: post.Title, TitleText: titleText(post), Date: post.Date, Modified: post.Modified, Published: post.Published, Draft: post.Draft, Private: post.Private, Template: post.Template, Tags: sortedStrings(post.Tags), Description: post.Description, Feeds: append([]string(nil), feedNames...), Aliases: aliases(post), Image: extraString(post, "image"), Video: extraString(post, "video"), Avatar: resolvedAvatar(post), Bio: resolvedBio(post), Thumbnail: extraString(post, "thumbnail"), Cover: firstExtraString(post, "cover", "cover_image"), OGImage: firstExtraString(post, "og_image", "social_image"), Author: post.Author, Authors: append([]string(nil), post.GetAuthors()...), Category: extraString(post, "category"), Categories: extraStrings(post, "categories")}
	if !post.Private {
		return document
	}

	// Private records are useful for metadata consumers, but they must not
	// inherit values derived from the private body or sensitive media URLs.
	if !postExtraBool(post, "_title_explicit") {
		document.Title = nil
		document.TitleText = nil
	}
	if !postExtraBool(post, "_description_explicit") {
		document.Description = nil
	}
	document.Image = nil
	document.Video = nil
	document.Bio = nil
	document.Thumbnail = nil
	document.Cover = nil
	document.OGImage = nil
	if extraString(post, "avatar") == nil {
		document.Avatar = nil
	}
	return document
}

func firstExtraString(post *models.Post, keys ...string) *string {
	for _, key := range keys {
		if value := extraString(post, key); value != nil {
			return value
		}
	}
	return nil
}

func resolvedAvatar(post *models.Post) *string {
	if value := extraString(post, "avatar"); value != nil {
		return value
	}
	for index := range post.AuthorObjects {
		author := &post.AuthorObjects[index]
		if author.Avatar != nil && *author.Avatar != "" {
			return author.Avatar
		}
	}
	return nil
}

func resolvedBio(post *models.Post) *string {
	if value := extraString(post, "bio"); value != nil {
		return value
	}
	for index := range post.AuthorObjects {
		author := &post.AuthorObjects[index]
		if author.Bio != nil && *author.Bio != "" {
			return author.Bio
		}
	}
	return nil
}

func extraString(post *models.Post, key string) *string {
	value, ok := post.Extra[key].(string)
	if !ok || value == "" {
		return nil
	}
	return &value
}

func extraStrings(post *models.Post, key string) []string {
	value := post.Extra[key]
	var result []string
	switch values := value.(type) {
	case []string:
		result = append(result, values...)
	case []interface{}:
		for _, item := range values {
			if value, ok := item.(string); ok {
				result = append(result, value)
			}
		}
	case string:
		result = append(result, values)
	}
	return sortedStrings(result)
}

func titleText(post *models.Post) *string {
	if post.Title == nil && !post.TitleTextDerived && post.TitleText == "" {
		return nil
	}
	v := post.PlainTitle()
	return &v
}

func aliases(post *models.Post) []string {
	values, ok := post.Extra["aliases"]
	if !ok {
		return nil
	}
	result := make([]string, 0)
	switch v := values.(type) {
	case []string:
		result = append(result, v...)
	case []interface{}:
		for _, item := range v {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
	}
	return sortedStrings(result)
}

func sortedStrings(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}

type contentIndexOptions struct {
	Enabled       bool
	Output        string
	SchemaVersion int
}

func contentIndexConfig(extra map[string]interface{}) (contentIndexOptions, bool, error) {
	value, ok := extra["content_index"]
	if !ok {
		return contentIndexOptions{}, false, nil
	}
	result := contentIndexOptions{Enabled: true}
	switch cfg := value.(type) {
	case map[string]interface{}:
		if enabled, exists := cfg["enabled"]; exists {
			value, valid := enabled.(bool)
			if !valid {
				return contentIndexOptions{}, false, fmt.Errorf("content_index.enabled must be boolean")
			}
			result.Enabled = value
		}
		if output, exists := cfg["output"]; exists {
			value, valid := output.(string)
			if !valid {
				return contentIndexOptions{}, false, fmt.Errorf("content_index.output must be a string")
			}
			result.Output = value
		}
		if value, exists := cfg["schema_version"]; exists {
			version, valid := parseIntFromInterface(value)
			if !valid {
				return contentIndexOptions{}, false, fmt.Errorf("content_index.schema_version must be an integer")
			}
			result.SchemaVersion = version
		}
	case contentIndexOptions:
		result = cfg
	default:
		return contentIndexOptions{}, false, fmt.Errorf("content_index must be a table")
	}
	return result, true, nil
}
