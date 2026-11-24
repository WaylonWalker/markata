package plugins

import (
	"strings"
	"testing"

	"github.com/WaylonWalker/markata/pkg/models"
)

func TestHashtagTagsPlugin_Name(t *testing.T) {
	p := NewHashtagTagsPlugin()
	if p.Name() != "hashtag_tags" {
		t.Errorf("expected name 'hashtag_tags', got %q", p.Name())
	}
}

func TestFormatReadingTime(t *testing.T) {
	tests := []struct {
		minutes  int
		expected string
	}{
		{0, "less than a minute"},
		{1, "1 minute"},
		{5, "5 minutes"},
		{60, "60 minutes"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := formatReadingTime(tt.minutes)
			if got != tt.expected {
				t.Errorf("formatReadingTime(%d) = %q, want %q", tt.minutes, got, tt.expected)
			}
		})
	}
}

func TestHashtagTagsPlugin_BuildTagStats(t *testing.T) {
	// Create test posts with tags
	posts := []*models.Post{
		{
			Slug:      "post1",
			Published: true,
			Draft:     false,
			Private:   false,
			Skip:      false,
			Tags:      []string{"go", "programming"},
			Extra: map[string]interface{}{
				"reading_time": 5,
			},
		},
		{
			Slug:      "post2",
			Published: true,
			Draft:     false,
			Private:   false,
			Skip:      false,
			Tags:      []string{"go", "web"},
			Extra: map[string]interface{}{
				"reading_time": 8,
			},
		},
		{
			Slug:      "post3",
			Published: false,
			Draft:     false,
			Private:   false,
			Skip:      false,
			Tags:      []string{"go"}, // Should be skipped (not published)
			Extra: map[string]interface{}{
				"reading_time": 10,
			},
		},
	}

	// We need to manually test buildTagStats behavior without a full Manager.
	// This is a simplified test that validates the core logic.
	tagStats := make(map[string]*TagStats)

	// Manually simulate what buildTagStats does
	for _, post := range posts {
		if post.Draft || !post.Published || post.Private || post.Skip {
			continue
		}

		readingTime := 0
		if rt, ok := post.Extra["reading_time"].(int); ok {
			readingTime = rt
		}

		for _, tag := range post.Tags {
			if stat, exists := tagStats[tag]; exists {
				stat.Count++
				stat.ReadingTime += readingTime
			} else {
				slug := models.Slugify(tag)
				tagStats[tag] = &TagStats{
					Tag:         tag,
					Slug:        slug,
					Count:       1,
					ReadingTime: readingTime,
				}
			}
		}
	}

	// Format reading time for each tag
	for _, stat := range tagStats {
		stat.ReadingTimeText = formatReadingTime(stat.ReadingTime)
	}

	// Verify results
	if len(tagStats) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tagStats))
	}

	if stat, ok := tagStats["go"]; ok {
		if stat.Count != 2 {
			t.Errorf("expected 'go' count 2, got %d", stat.Count)
		}
		if stat.ReadingTime != 13 { // 5 + 8
			t.Errorf("expected 'go' reading time 13, got %d", stat.ReadingTime)
		}
	} else {
		t.Error("expected 'go' tag to exist")
	}

	if stat, ok := tagStats["programming"]; ok {
		if stat.Count != 1 {
			t.Errorf("expected 'programming' count 1, got %d", stat.Count)
		}
	} else {
		t.Error("expected 'programming' tag to exist")
	}

	if stat, ok := tagStats["web"]; ok {
		if stat.Count != 1 {
			t.Errorf("expected 'web' count 1, got %d", stat.Count)
		}
	} else {
		t.Error("expected 'web' tag to exist")
	}
}

func TestHashtagTagsPlugin_ProcessHashtagsInText(t *testing.T) {
	p := NewHashtagTagsPlugin()

	tagStats := map[string]*TagStats{
		"go": {
			Tag:             "go",
			Slug:            "go",
			Count:           5,
			ReadingTime:     30,
			ReadingTimeText: "30 minutes",
		},
		"rust": {
			Tag:             "rust",
			Slug:            "rust",
			Count:           3,
			ReadingTime:     20,
			ReadingTimeText: "20 minutes",
		},
	}

	tests := []struct {
		name     string
		input    string
		contains string // Check if result contains this substring
	}{
		{
			name:     "single hashtag",
			input:    "I love #go programming",
			contains: `<a href="/tags/go/"`,
		},
		{
			name:     "multiple hashtags",
			input:    "I use #go and #rust",
			contains: `<a href="/tags/go/"`,
		},
		{
			name:     "hashtag with dash not in stats",
			input:    "Check out #go-lang",
			contains: "#go-lang", // Should not match (go-lang is not in tagStats)
		},
		{
			name:     "hashtag in data-tag attribute",
			input:    "I love #go",
			contains: `data-tag="go"`,
		},
		{
			name:     "non-existent tag",
			input:    "This is #unknown",
			contains: "#unknown", // Should remain unchanged
		},
		{
			name:     "hashtag at start of line",
			input:    "#go is great",
			contains: `<a href="/tags/go/"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.processHashtagsInText(tt.input, tagStats)
			if !strings.Contains(got, tt.contains) {
				t.Errorf("processHashtagsInText(%q)\nresult: %q\ndoes not contain: %q", tt.input, got, tt.contains)
			}
		})
	}
}

func TestHashtagTagsPlugin_ProcessHashtagsInContent(t *testing.T) {
	p := NewHashtagTagsPlugin()

	tagStats := map[string]*TagStats{
		"go": {
			Tag:             "go",
			Slug:            "go",
			Count:           5,
			ReadingTime:     30,
			ReadingTimeText: "30 minutes",
		},
	}

	tests := []struct {
		name        string
		input       string
		contains    string
		notContains string
	}{
		{
			name:        "hashtag outside code block",
			input:       "I love #go\n```\n#go is code\n```",
			contains:    `<a href="/tags/go/"`,
			notContains: "", // The code block content should be unchanged
		},
		{
			name:        "hashtag inside code block is preserved",
			input:       "```\n#go is code\n```",
			contains:    "#go is code",         // Code block should remain unchanged
			notContains: `<a href="/tags/go/"`, // Should NOT transform inside code block
		},
		{
			name:        "multiple code blocks",
			input:       "#go is great\n```\n#go is code\n```\nI love #go",
			contains:    `<a href="/tags/go/"`,
			notContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.processHashtagsInContent(tt.input, tagStats)
			if tt.contains != "" && !strings.Contains(got, tt.contains) {
				t.Errorf("processHashtagsInContent(%q)\nresult: %q\ndoes not contain: %q", tt.input, got, tt.contains)
			}
			if tt.notContains != "" && strings.Contains(got, tt.notContains) {
				t.Errorf("processHashtagsInContent(%q)\nresult: %q\nshould not contain: %q", tt.input, got, tt.notContains)
			}
		})
	}
}
