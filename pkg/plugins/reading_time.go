// Package plugins provides lifecycle plugins for markata-go.
package plugins

import (
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
)

const readingTimePluginName = "reading_time"

// ReadingTimePlugin calculates the word count and estimated reading time
// for each post during the transform stage.
type ReadingTimePlugin struct {
	// wordsPerMinute is the average reading speed (default: 200)
	wordsPerMinute int
}

// NewReadingTimePlugin creates a new ReadingTimePlugin with default settings.
func NewReadingTimePlugin() *ReadingTimePlugin {
	return &ReadingTimePlugin{
		wordsPerMinute: defaultWordsPerMinute,
	}
}

// Name returns the unique name of the plugin.
func (p *ReadingTimePlugin) Name() string {
	return readingTimePluginName
}

// Configure reads configuration options for the plugin.
func (p *ReadingTimePlugin) Configure(m *lifecycle.Manager) error {
	config := m.Config()
	if config.Extra != nil {
		if wpm, ok := parseIntFromInterface(config.Extra["words_per_minute"]); ok && wpm > 0 {
			p.wordsPerMinute = wpm
		}
	}
	return nil
}

// Transform calculates word count and reading time for each eligible post.
// The legacy lifecycle keeps its existing concurrent slice processing; the
// post-local operation is separated so the DAG can schedule explicit item
// tasks without duplicating reading-time semantics.
func (p *ReadingTimePlugin) Transform(m *lifecycle.Manager) error {
	return m.ProcessPostsSliceConcurrently(p.PostsForTransform(m), p.TransformPost)
}

// PostsForTransform returns the posts reading_time would process for the
// current manager state. Keeping selection here ensures the legacy hook and
// scheduler-owned item tasks share skip/empty and Serve-incremental behavior.
func (p *ReadingTimePlugin) PostsForTransform(m *lifecycle.Manager) []*models.Post {
	if m == nil {
		return nil
	}

	posts := m.FilterPosts(readingTimePostEligible)
	if !lifecycle.IsServeIncremental(m) {
		return posts
	}

	affected := lifecycle.GetServeAffectedPaths(m)
	if len(affected) == 0 {
		return posts
	}
	filtered := posts[:0]
	for _, post := range posts {
		if affected[post.Path] {
			filtered = append(filtered, post)
		}
	}
	return filtered
}

// TransformPost calculates and stores reading-time fields for one post. It is
// deliberately post-local: callers own selection/order while this method owns
// only the supplied post's reading-time fields.
func (p *ReadingTimePlugin) TransformPost(post *models.Post) error {
	if !readingTimePostEligible(post) {
		return nil
	}

	metrics := calculateReadingTimeMetrics(post.Content, p.wordsPerMinute, false)
	post.Set("word_count", metrics.WordCount)
	post.Set("reading_time", metrics.ReadingTime)
	post.Set("reading_time_text", metrics.ReadingTimeText)
	return nil
}

func readingTimePostEligible(post *models.Post) bool {
	return post != nil && !post.Skip && post.Content != ""
}

// countWords counts the number of words in markdown content.
// It excludes code blocks, URLs, and other non-prose elements.
func (p *ReadingTimePlugin) countWords(content string) int {
	return countReadingWords(content, false)
}

// calculateReadingTime estimates reading time in minutes based on word count.
// Returns at least 1 minute for any non-empty content.
func (p *ReadingTimePlugin) calculateReadingTime(wordCount int) int {
	return calculateReadingMinutes(wordCount, p.wordsPerMinute)
}

// formatReadingTime creates a human-readable reading time string.
func (p *ReadingTimePlugin) formatReadingTime(minutes int) string {
	return formatReadingTimeText(minutes)
}

// SetWordsPerMinute sets the average reading speed.
func (p *ReadingTimePlugin) SetWordsPerMinute(wpm int) {
	if wpm > 0 {
		p.wordsPerMinute = wpm
	}
}

// ReadingTimeResult holds the calculated reading metrics for a post.
type ReadingTimeResult struct {
	// WordCount is the number of words in the post
	WordCount int `json:"word_count"`

	// ReadingTime is the estimated reading time in minutes
	ReadingTime int `json:"reading_time"`

	// ReadingTimeText is a formatted reading time string
	ReadingTimeText string `json:"reading_time_text"`
}

// GetReadingTime extracts reading time data from a post's Extra map.
// Returns nil if reading time hasn't been calculated.
func GetReadingTime(post *models.Post) *ReadingTimeResult {
	if post.Extra == nil {
		return nil
	}

	wordCount, hasWC := post.Extra["word_count"].(int)
	readingTime, hasRT := post.Extra["reading_time"].(int)
	readingTimeText, hasRTT := post.Extra["reading_time_text"].(string)

	if !hasWC || !hasRT {
		return nil
	}

	result := &ReadingTimeResult{
		WordCount:   wordCount,
		ReadingTime: readingTime,
	}

	if hasRTT {
		result.ReadingTimeText = readingTimeText
	}

	return result
}

// Ensure ReadingTimePlugin implements the required interfaces.
var (
	_ lifecycle.Plugin          = (*ReadingTimePlugin)(nil)
	_ lifecycle.ConfigurePlugin = (*ReadingTimePlugin)(nil)
	_ lifecycle.TransformPlugin = (*ReadingTimePlugin)(nil)
)
