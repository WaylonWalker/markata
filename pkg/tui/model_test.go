package tui

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/services"
)

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input int
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{12, "12"},
		{123, "123"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatNumber(tt.input)
			if got != tt.want {
				t.Errorf("formatNumber(%d) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", "", 0},
		{"single word", "hello", 1},
		{"multiple words", "hello world", 2},
		{"with newlines", "hello\nworld\nfoo", 3},
		{"with extra spaces", "  hello   world  ", 2},
		{"markdown content", "# Title\n\nThis is a paragraph with several words.", 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countWords(tt.input)
			if got != tt.want {
				t.Errorf("countWords(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetContentPreview(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		maxChars int
		maxLines int
		maxWidth int
		wantLen  int // Check that output length is reasonable
	}{
		{"empty", "", 500, 12, 80, 7}, // "(empty)"
		{"short", "hello", 500, 12, 80, 5},
		{"multiline", "line1\nline2\nline3", 500, 12, 80, 17},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getContentPreview(tt.content, tt.maxChars, tt.maxLines, tt.maxWidth)
			if len(got) < tt.wantLen {
				t.Errorf("getContentPreview() length = %d, want at least %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestHandleEnter_PostsView(t *testing.T) {
	// Create a model with some posts
	m := Model{
		view:   ViewPosts,
		cursor: 0,
		posts: []*models.Post{
			{Path: "test.md", Content: "Test content"},
		},
	}

	// Simulate pressing Enter
	newModel, _ := m.handleEnter()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	if newM.view != ViewPostDetail {
		t.Errorf("view = %q, want %q", newM.view, ViewPostDetail)
	}

	if newM.selectedPost == nil {
		t.Error("selectedPost is nil, expected a post")
	}

	if newM.previousView != ViewPosts {
		t.Errorf("previousView = %q, want %q", newM.previousView, ViewPosts)
	}
}

func TestHandleEscape_DetailView(t *testing.T) {
	title := "Test Post"
	m := Model{
		view:         ViewPostDetail,
		previousView: ViewPosts,
		selectedPost: &models.Post{Title: &title},
	}

	newModel, _ := m.handleEscape()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	if newM.view != ViewPosts {
		t.Errorf("view = %q, want %q", newM.view, ViewPosts)
	}

	if newM.selectedPost != nil {
		t.Error("selectedPost should be nil after escape")
	}
}

func TestHandleDetailViewKey_Quit(t *testing.T) {
	m := Model{
		view: ViewPostDetail,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := m.handleDetailViewKey(msg)

	// Check that the command is tea.Quit
	if cmd == nil {
		t.Error("expected quit command, got nil")
	}
}

func TestHandleDetailViewKey_Escape(t *testing.T) {
	title := "Test"
	m := Model{
		view:         ViewPostDetail,
		previousView: ViewPosts,
		selectedPost: &models.Post{Title: &title},
	}

	msg := tea.KeyMsg{Type: tea.KeyEscape}
	newModel, _ := m.handleDetailViewKey(msg)
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	if newM.view != ViewPosts {
		t.Errorf("view = %q, want %q", newM.view, ViewPosts)
	}
}

func TestRenderPostDetail_NilPost(t *testing.T) {
	m := Model{
		view:         ViewPostDetail,
		selectedPost: nil,
		width:        80,
		height:       24,
	}

	result := m.renderPostDetail()
	if result != "No post selected." {
		t.Errorf("expected 'No post selected.', got %q", result)
	}
}

func TestRenderPostDetail_WithPost(t *testing.T) {
	title := "Test Title"
	desc := "Test description"
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	m := Model{
		view:   ViewPosts,
		cursor: 0,
		posts: []*models.Post{
			{
				Title:       &title,
				Path:        "posts/test.md",
				Date:        &date,
				Published:   true,
				Tags:        []string{"go", "test"},
				Description: &desc,
				Content:     "# Test\n\nThis is test content.",
			},
		},
		width:  80,
		height: 24,
		theme:  DefaultTheme(),
	}

	// First, handle Enter to initialize the viewport and switch to detail view
	newModel, cmd := m.handleEnter()
	var ok bool
	m, ok = newModel.(Model)
	if !ok {
		t.Fatal("handleEnter returned unexpected type")
	}
	_ = cmd // Command not needed in this test

	result := m.renderPostDetail()

	// Check that key elements are present
	if !contains(result, "Test Title") {
		t.Error("expected title in output")
	}
	if !contains(result, "posts/test.md") {
		t.Error("expected path in output")
	}
	if !contains(result, "2024-01-15") {
		t.Error("expected date in output")
	}
	if !contains(result, "true") {
		t.Error("expected published status in output")
	}
	if !contains(result, "go, test") {
		t.Error("expected tags in output")
	}
	if !contains(result, "[e]dit") {
		t.Error("expected edit keybinding in output")
	}
	if !contains(result, "[Esc] back") {
		t.Error("expected escape keybinding in output")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Tests for drill-down navigation (Issue #250)

func TestHandleEnter_TagsView_DrillDown(t *testing.T) {
	// Create a model with tags
	m := Model{
		view:   ViewTags,
		cursor: 0,
		tags: []services.TagInfo{
			{Name: "go", Count: 5, Slug: "go"},
			{Name: "python", Count: 3, Slug: "python"},
		},
	}

	// Simulate pressing Enter on the first tag
	newModel, cmd := m.handleEnter()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// View should change to posts
	if newM.view != ViewPosts {
		t.Errorf("view = %q, want %q", newM.view, ViewPosts)
	}

	// Active filter should be set
	if newM.activeFilter == nil {
		t.Fatal("activeFilter should not be nil")
	}
	if newM.activeFilter.Type != "tag" {
		t.Errorf("activeFilter.Type = %q, want %q", newM.activeFilter.Type, "tag")
	}
	if newM.activeFilter.Name != "go" {
		t.Errorf("activeFilter.Name = %q, want %q", newM.activeFilter.Name, "go")
	}

	// Cursor should be reset
	if newM.cursor != 0 {
		t.Errorf("cursor = %d, want 0", newM.cursor)
	}

	// A command should be returned to load posts
	if cmd == nil {
		t.Error("expected command to load posts, got nil")
	}
}

func TestHandleEnter_FeedsView_DrillDown(t *testing.T) {
	// Create a model with feeds
	theme := DefaultTheme()
	feedsTable := createFeedsTableWithTheme(80, theme)
	feedsTable.SetCursor(1) // Select the second feed

	feeds := []*lifecycle.Feed{
		{Name: "main", Path: "feed/main.xml"},
		{Name: "blog", Path: "feed/blog.xml"},
	}

	m := Model{
		view:       ViewFeeds,
		cursor:     1,
		feedsTable: feedsTable,
		feeds:      feeds,
	}
	m.feedsTable.SetRows(m.feedsToRows())

	// Simulate pressing Enter on the second feed
	newModel, cmd := m.handleEnter()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// View should change to posts
	if newM.view != ViewPosts {
		t.Errorf("view = %q, want %q", newM.view, ViewPosts)
	}

	// Active filter should be set
	if newM.activeFilter == nil {
		t.Fatal("activeFilter should not be nil")
	}
	if newM.activeFilter.Type != "feed" {
		t.Errorf("activeFilter.Type = %q, want %q", newM.activeFilter.Type, "feed")
	}
	if newM.activeFilter.Name != "blog" {
		t.Errorf("activeFilter.Name = %q, want %q", newM.activeFilter.Name, "blog")
	}

	// Cursor should be reset
	if newM.cursor != 0 {
		t.Errorf("cursor = %d, want 0", newM.cursor)
	}

	// A command should be returned to load posts
	if cmd == nil {
		t.Error("expected command to load posts, got nil")
	}
}

func TestHandleEscape_ClearActiveFilter(t *testing.T) {
	// Create a model with an active filter
	m := Model{
		view: ViewPosts,
		activeFilter: &FilterContext{
			Type: "tag",
			Name: "go",
		},
		cursor: 5,
	}

	// Simulate pressing Escape
	newModel, cmd := m.handleEscape()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// Active filter should be cleared
	if newM.activeFilter != nil {
		t.Error("activeFilter should be nil after escape")
	}

	// Cursor should be reset
	if newM.cursor != 0 {
		t.Errorf("cursor = %d, want 0", newM.cursor)
	}

	// A command should be returned to reload posts
	if cmd == nil {
		t.Error("expected command to load posts, got nil")
	}
}

func TestHandleEscape_NoFilterNoOp(t *testing.T) {
	// Create a model without an active filter in posts view
	m := Model{
		view:         ViewPosts,
		activeFilter: nil,
		cursor:       5,
	}

	// Simulate pressing Escape
	newModel, cmd := m.handleEscape()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// Cursor should not change
	if newM.cursor != 5 {
		t.Errorf("cursor = %d, want 5", newM.cursor)
	}

	// No command should be returned
	if cmd != nil {
		t.Error("expected nil command when no filter active")
	}
}

func TestFilterContext(t *testing.T) {
	tests := []struct {
		name     string
		filter   *FilterContext
		wantType string
		wantName string
	}{
		{
			name:     "tag filter",
			filter:   &FilterContext{Type: "tag", Name: "golang"},
			wantType: "tag",
			wantName: "golang",
		},
		{
			name:     "feed filter",
			filter:   &FilterContext{Type: "feed", Name: "main"},
			wantType: "feed",
			wantName: "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.filter.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", tt.filter.Type, tt.wantType)
			}
			if tt.filter.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tt.filter.Name, tt.wantName)
			}
		})
	}
}

func TestRenderLayout_WithActiveFilter(t *testing.T) {
	m := Model{
		view:  ViewPosts,
		width: 80,
		activeFilter: &FilterContext{
			Type: "tag",
			Name: "go",
		},
		sortBy:    "date",
		sortOrder: services.SortDesc,
		theme:     DefaultTheme(),
	}

	result := m.renderLayout("test content")

	// Should show the filter indicator
	if !contains(result, "tag: go") {
		t.Error("expected filter indicator in header")
	}

	// Should show Esc hint for clearing filter
	if !contains(result, "Esc:clear filter") {
		t.Error("expected 'Esc:clear filter' in status bar")
	}
}

func TestRenderLayout_WithoutActiveFilter(t *testing.T) {
	m := Model{
		view:         ViewPosts,
		width:        80,
		activeFilter: nil,
		sortBy:       "date",
		sortOrder:    services.SortDesc,
		theme:        DefaultTheme(),
	}

	result := m.renderLayout("test content")

	// Should not show filter indicator
	if contains(result, "tag:") || contains(result, "feed:") {
		t.Error("should not show filter indicator when no active filter")
	}

	// Should not show Esc hint for clearing filter
	if contains(result, "Esc:clear filter") {
		t.Error("should not show 'Esc:clear filter' when no active filter")
	}
}

func TestHandleEnter_EmptyTags(t *testing.T) {
	// Create a model with no tags
	m := Model{
		view:   ViewTags,
		cursor: 0,
		tags:   []services.TagInfo{},
	}

	// Simulate pressing Enter
	newModel, cmd := m.handleEnter()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// View should not change when no tags available
	if newM.view != ViewTags {
		t.Errorf("view should remain %q when no tags, got %q", ViewTags, newM.view)
	}

	// No command should be returned
	if cmd != nil {
		t.Error("expected nil command when no tags")
	}
}

func TestHandleEnter_EmptyFeeds(t *testing.T) {
	// Create a model with no feeds
	theme := DefaultTheme()
	feedsTable := createFeedsTableWithTheme(80, theme)

	m := Model{
		view:       ViewFeeds,
		cursor:     0,
		feedsTable: feedsTable,
		feeds:      []*lifecycle.Feed{},
	}

	// Simulate pressing Enter
	newModel, cmd := m.handleEnter()
	newM, ok := newModel.(Model)
	if !ok {
		t.Fatal("expected Model type")
	}

	// View should not change when no feeds available
	if newM.view != ViewFeeds {
		t.Errorf("view should remain %q when no feeds, got %q", ViewFeeds, newM.view)
	}

	// No command should be returned
	if cmd != nil {
		t.Error("expected nil command when no feeds")
	}
}

func TestFilterVisiblePosts_SearchesOnlyCurrentList(t *testing.T) {
	title := "Go tooling"
	description := "A guide to terminal workflows"
	posts := []*models.Post{
		{Path: "posts/go.md", Title: &title},
		{Path: "notes/terminal.md", Description: &description},
		{Path: "posts/rust.md"},
	}

	for _, tt := range []struct {
		query string
		want  int
	}{
		{"go", 1},
		{"terminal", 1},
		{"rust", 1},
		{"missing", 0},
	} {
		t.Run(tt.query, func(t *testing.T) {
			got := filterVisiblePosts(posts, tt.query)
			if len(got) != tt.want {
				t.Fatalf("filterVisiblePosts(%q) returned %d posts, want %d", tt.query, len(got), tt.want)
			}
		})
	}
}

func TestFeedToRow_LabelsSpecialFeeds(t *testing.T) {
	model := Model{theme: DefaultTheme()}
	row := model.feedToRow(&lifecycle.Feed{
		Name:           "secret",
		IncludePrivate: true,
		Hidden:         true,
		Automated:      true,
	})

	if !contains(row[0], "secret [private,hidden,auto]") {
		t.Fatalf("feed title = %q, missing status labels", row[0])
	}
}

func TestRenderPostsTable_PreservesHeaders(t *testing.T) {
	theme := DefaultTheme()
	title := "Example"
	m := Model{
		theme:      theme,
		sortBy:     "date",
		postsTable: createPostsTableWithTheme(120, theme, "date"),
		posts:      []*models.Post{{Title: &title}},
	}
	m.postsTable.SetRows(m.postsToRows())

	view := m.renderPostsTable()
	for _, header := range []string{"TITLE", "DATE", "WORDS", "READ", "TAGS", "PATH"} {
		if !contains(view, header) {
			t.Errorf("posts table is missing %q header: %s", header, view)
		}
	}
	if contains(view, "GAGS") {
		t.Error("posts table must not replace the TAGS header with its shortcut")
	}
}

func TestFilterVisibleFeeds_MatchesNameTitleAndPath(t *testing.T) {
	feeds := []*lifecycle.Feed{
		{Name: "daily", Title: "Daily Notes", Path: "feeds/daily"},
		{Name: "tags/go", Title: "Posts tagged: Go", Path: "feeds/tags/go"},
		{Name: "archive", Title: "Archive", Path: "feeds/archive"},
	}

	for _, tt := range []struct {
		query string
		want  int
	}{
		{"daily", 1},
		{"posts tagged", 1},
		{"tags/go", 1},
		{"missing", 0},
	} {
		t.Run(tt.query, func(t *testing.T) {
			if got := filterVisibleFeeds(feeds, tt.query); len(got) != tt.want {
				t.Fatalf("filterVisibleFeeds(%q) returned %d feeds, want %d", tt.query, len(got), tt.want)
			}
		})
	}
}
