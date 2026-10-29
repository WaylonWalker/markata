package tui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/services"
)

// View represents different TUI views
type View string

const (
	ViewPosts      View = "posts"
	ViewTags       View = "tags"
	ViewFeeds      View = "feeds"
	ViewHelp       View = "help"
	ViewPostDetail View = "post_detail"
	ViewConfig     View = "config"
	ViewSearch     View = "search"
)

// Mode represents the input mode
type Mode string

const (
	ModeNormal  Mode = "normal"
	ModeFilter  Mode = "filter"
	ModeCommand Mode = "command"
)

// FilterContext tracks the active filter for drill-down navigation
type FilterContext struct {
	Type string // "tag" or "feed"
	Name string // The tag name or feed name
}

// sortHotkeyMap maps capital letter keys to sort fields (k9s-inspired)
var sortHotkeyMap = map[string]string{
	"T": "title",
	"D": "date",
	"W": "words",
	"P": "path",
	"R": "reading_time",
	"G": "tags",
}

// footerButton represents a clickable button in the footer
type footerButton struct {
	label  string
	key    string
	startX int
	endX   int
	action func(*Model) (tea.Model, tea.Cmd)
}

// Model is the main Bubble Tea model
type Model struct {
	app               *services.App
	posts             []*models.Post
	visiblePosts      []*models.Post    // Unfiltered posts in the current list; used by / filtering.
	visibleFeeds      []*lifecycle.Feed // Unfiltered feeds in the current list; used by / filtering.
	tags              []services.TagInfo
	feeds             []*lifecycle.Feed
	postsTable        table.Model
	tagsTable         table.Model
	feedsTable        table.Model
	cursor            int
	view              View
	previousView      View // Track previous view for returning from detail
	mode              Mode
	filter            string
	filterBeforeInput string
	filterInput       textinput.Model
	cmdInput          textinput.Model
	width             int
	height            int
	err               error
	selectedPost      *models.Post   // The post being viewed in detail
	postViewport      viewport.Model // Viewport for scrolling post detail content
	helpViewport      viewport.Model // Viewport for scrolling help content

	// Sort state
	sortBy       string             // "date", "title", "words", "path"
	sortOrder    services.SortOrder // SortAsc or SortDesc
	showSortMenu bool
	sortMenuIdx  int

	// Theme styling
	theme *Theme

	// Drill-down filter state
	activeFilter *FilterContext // Active tag/feed filter for drill-down navigation

	// Footer button tracking for mouse clicks
	footerButtons []footerButton
	mouseX        int // Current mouse X position
	mouseY        int // Current mouse Y position

	// Help search state
	helpSearchInput  textinput.Model // Search input for help view
	helpSearchMode   bool            // Whether search is active in help view
	helpSearchQuery  string          // Current search query
	helpContentLines []string        // All help content lines
	helpMatchedLines []int           // Indices of lines that match search
	helpCurrentMatch int             // Current match index for n/N navigation

	// Refresh state
	lastRefresh time.Time // Track last refresh time
	refreshing  bool      // Indicate refresh in progress
	loading     bool      // Initial site data is loading in the background

	// Search view state
	searchInput          textinput.Model    // Main search input
	searchFilterInput    textinput.Model    // Advanced filter input
	searchResults        []searchResultItem // Current search results
	searchTable          table.Model        // Results table
	searchMode           searchModeType     // Simple or advanced
	searchFilterFocused  bool               // Whether filter input is focused
	searchAdvancedFilter string             // Active filter expression
	searchLastQuery      string             // Last query (for debounce)
	searchDebounceTime   time.Time          // Timestamp for debounce
	searchFuzzy          bool               // Fuzzy matching enabled
	searchLimit          int                // Max results (0 = default 50)

	// Config view state
	configSections []configSection // Expanded config data
	configCursor   int             // Current cursor position in config view
	configFilter   string          // Search filter for config keys
	configExpanded map[string]bool // Track which sections are expanded
	configViewport viewport.Model  // Viewport for scrolling config content
}

// Messages
type postsLoadedMsg struct {
	posts []*models.Post
}

type tagsLoadedMsg struct {
	tags []services.TagInfo
}

type feedsLoadedMsg struct {
	feeds []*lifecycle.Feed
}

type errMsg struct {
	err error
}

type editorFinishedMsg struct {
	err error
}

type refreshStartedMsg struct{}

type refreshCompletedMsg struct {
	posts []*models.Post
	tags  []services.TagInfo
	feeds []*lifecycle.Feed
}

type initialDataLoadedMsg struct {
	posts []*models.Post
	tags  []services.TagInfo
	feeds []*lifecycle.Feed
}

// sortOption represents a sort field option
type sortOption struct {
	label string
	value string
}

// sortOptions available for sorting posts
var sortOptions = []sortOption{
	{"Date", "date"},
	{"Title", "title"},
	{"Word Count", "words"},
	{"Path", "path"},
	{"Reading Time", "reading_time"},
	{"Tags", "tags"},
}

// configSection represents a section in the config view
type configSection struct {
	name     string       // Section name (e.g., "Site Metadata")
	key      string       // Key identifier for expansion tracking
	items    []configItem // Items in this section
	expanded bool         // Whether section is expanded
}

// configItem represents a single config key-value pair
type configItem struct {
	key   string // Config key name
	value string // Formatted value
	level int    // Indentation level (0 = top-level)
}

// NewModel creates a new TUI model with default theme.
func NewModel(app *services.App) Model {
	return NewModelWithTheme(app, nil)
}

// NewModelWithTheme creates a new TUI model with the specified theme.
// If theme is nil, the default theme is used.
func NewModelWithTheme(app *services.App, theme *Theme) Model {
	if theme == nil {
		theme = DefaultTheme()
	}

	filterInput := textinput.New()
	filterInput.Placeholder = "Filter visible posts..."
	filterInput.CharLimit = 100

	cmdInput := textinput.New()
	cmdInput.Placeholder = "Command..."
	cmdInput.CharLimit = 100

	helpSearchInput := textinput.New()
	helpSearchInput.Placeholder = "Search help..."
	helpSearchInput.CharLimit = 100

	// Initialize posts table with columns and theme
	postsTable := createPostsTableWithTheme(80, theme, "date") // Default width, updated on resize

	// Initialize tags table with columns and theme
	tagsTable := createTagsTableWithTheme(80, theme) // Default width, will be updated on resize

	// Initialize feeds table with columns and theme
	feedsTable := createFeedsTableWithTheme(80, theme) // Default width, will be updated on resize

	m := Model{
		app:             app,
		view:            ViewPosts,
		mode:            ModeNormal,
		filterInput:     filterInput,
		cmdInput:        cmdInput,
		helpSearchInput: helpSearchInput,
		postsTable:      postsTable,
		tagsTable:       tagsTable,
		feedsTable:      feedsTable,
		sortBy:          "date",
		sortOrder:       services.SortDesc,
		theme:           theme,
		configExpanded:  make(map[string]bool),
		loading:         true,
	}

	return m
}

// createTableWithTheme creates and configures a table with theme colors and the given columns.
func createTableWithTheme(columns []table.Column, theme *Theme) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10), // Will be updated on resize
	)

	// Apply theme-aware styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.Colors.Border).
		BorderBottom(true).
		Bold(true).
		Foreground(theme.Colors.TableHeader)
	s.Selected = s.Selected.
		Foreground(theme.Colors.SelectedText).
		Background(theme.Colors.SelectedBg).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(theme.Colors.TableCell)
	t.SetStyles(s)

	return t
}

// createTagsTableWithTheme creates and configures the tags table with theme colors.
func createTagsTableWithTheme(width int, theme *Theme) table.Model {
	// Column widths: TAG(30) + COUNT(10) + WORDS(10) + READ(10) + SLUG(remaining)
	// Account for padding/borders (approximately 10 chars)
	slugWidth := width - 30 - 10 - 10 - 10 - 10
	if slugWidth < 15 {
		slugWidth = 15
	}

	columns := []table.Column{
		{Title: "TAG", Width: 30},
		{Title: "COUNT", Width: 10},
		{Title: "WORDS", Width: 10},
		{Title: "READ", Width: 10},
		{Title: "SLUG", Width: slugWidth},
	}

	return createTableWithTheme(columns, theme)
}

// createFeedsTableWithTheme creates and configures the feeds table with theme colors.
func createFeedsTableWithTheme(width int, theme *Theme) table.Model {
	// Column widths: NAME(32) + POSTS(8) + WORDS(10) + TOT TIME(10) + AVG TIME(10) + OUTPUT(remaining)
	// Account for padding/borders (approximately 12 chars)
	outputWidth := width - 32 - 8 - 10 - 10 - 10 - 12
	if outputWidth < 15 {
		outputWidth = 15
	}

	columns := []table.Column{
		{Title: "NAME", Width: 20},
		{Title: "POSTS", Width: 8},
		{Title: "WORDS", Width: 10},
		{Title: "TOT TIME", Width: 10},
		{Title: "AVG TIME", Width: 10},
		{Title: "OUTPUT", Width: outputWidth},
	}

	return createTableWithTheme(columns, theme)
}

// createPostsTableWithTheme creates and configures the posts table with theme colors.
func createPostsTableWithTheme(width int, theme *Theme, sortBy string) table.Model {
	return createTableWithTheme(postsTableColumns(width, theme, sortBy), theme)
}

func postsTableColumns(width int, _ *Theme, _ string) []table.Column {
	// Column widths: TITLE(35) + DATE(12) + WORDS(8) + READ(8) + TAGS(18) + PATH(remaining)
	// Account for padding/borders (approximately 10 chars)
	pathWidth := width - 35 - 12 - 8 - 8 - 18 - 10
	if pathWidth < 10 {
		pathWidth = 10
	}

	return []table.Column{
		{Title: "TITLE", Width: 35},
		{Title: "DATE", Width: 12},
		{Title: "WORDS", Width: 8},
		{Title: "READ", Width: 8},
		{Title: "TAGS", Width: 18},
		{Title: "PATH", Width: pathWidth},
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.loadInitialData()
}

// handleWindowResize handles terminal window resize events.
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	// Update table dimensions with theme
	m.postsTable = createPostsTableWithTheme(msg.Width, m.theme, m.sortBy)
	m.postsTable.SetHeight(msg.Height - 10) // Leave room for header/footer
	m.tagsTable = createTagsTableWithTheme(msg.Width, m.theme)
	m.tagsTable.SetHeight(msg.Height - 10)
	m.feedsTable = createFeedsTableWithTheme(msg.Width, m.theme)
	m.feedsTable.SetHeight(msg.Height - 10)
	if m.searchInput.Placeholder != "" {
		m.searchTable = createSearchTableWithTheme(msg.Width, m.theme)
		m.searchTable.SetHeight(msg.Height - 14)
		if len(m.searchResults) > 0 {
			m.searchTable.SetRows(m.searchResultsToRows())
		}
	}

	// Repopulate table if we have posts
	if len(m.posts) > 0 {
		m.postsTable.SetRows(m.postsToRows())
	}
	// Repopulate tags table if we have tags
	if len(m.tags) > 0 {
		m.tagsTable.SetRows(m.tagsToRows())
	}
	// Repopulate feeds table if we have feeds
	if len(m.feeds) > 0 {
		m.feedsTable.SetRows(m.feedsToRows())
	}

	// Update viewport dimensions based on current view
	m.updateViewportDimensions(msg.Width, msg.Height)

	return m, nil
}

// updateViewportDimensions updates the dimensions of the active viewport.
func (m *Model) updateViewportDimensions(w, h int) {
	width := w
	if width < 40 {
		width = 80
	}
	if width > 100 {
		width = 100
	}
	viewportHeight := h - 8
	if viewportHeight < 10 {
		viewportHeight = 10
	}

	if m.view == ViewPostDetail && m.selectedPost != nil {
		m.postViewport.Width = width - 4
		m.postViewport.Height = viewportHeight
	}
	if m.view == ViewHelp {
		m.helpViewport.Width = width - 4
		m.helpViewport.Height = viewportHeight
	}
	if m.view == ViewConfig {
		m.configViewport.Width = width - 4
		m.configViewport.Height = viewportHeight
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg)

	case tea.KeyMsg:
		if m.loading && msg.String() != "q" && msg.String() != "ctrl+c" {
			return m, nil
		}
		return m.handleKey(msg)

	case tea.MouseMsg:
		return m.handleMouse(msg)

	case postsLoadedMsg:
		m.visiblePosts = msg.posts
		m.posts = filterVisiblePosts(msg.posts, m.filter)
		m.postsTable.SetRows(m.postsToRows())
		return m, nil

	case initialDataLoadedMsg:
		m.loading = false
		m.visiblePosts = msg.posts
		m.posts = filterVisiblePosts(msg.posts, m.filter)
		m.tags = msg.tags
		m.visibleFeeds = msg.feeds
		m.feeds = filterVisibleFeeds(msg.feeds, m.filter)
		m.postsTable.SetRows(m.postsToRows())
		m.tagsTable.SetRows(m.tagsToRows())
		m.feedsTable.SetRows(m.feedsToRows())
		return m, nil

	case tagsLoadedMsg:
		m.tags = msg.tags
		m.tagsTable.SetRows(m.tagsToRows())
		return m, nil

	case feedsLoadedMsg:
		m.visibleFeeds = msg.feeds
		m.feeds = filterVisibleFeeds(msg.feeds, m.filter)
		m.feedsTable.SetRows(m.feedsToRows())
		return m, nil

	case errMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case editorFinishedMsg:
		// Reload posts in case content changed
		if msg.err != nil {
			m.err = msg.err
		}
		// Trigger refresh after editing
		return m, m.refreshData()

	case refreshStartedMsg:
		m.refreshing = true
		return m, nil

	case refreshCompletedMsg:
		m.refreshing = false
		m.lastRefresh = time.Now()
		m.visiblePosts = msg.posts
		m.posts = filterVisiblePosts(msg.posts, m.filter)
		m.tags = msg.tags
		m.visibleFeeds = msg.feeds
		m.feeds = filterVisibleFeeds(msg.feeds, m.filter)
		m.postsTable.SetRows(m.postsToRows())
		m.tagsTable.SetRows(m.tagsToRows())
		m.feedsTable.SetRows(m.feedsToRows())
		return m, nil

	case searchResultsMsg:
		m.handleSearchResults(msg)
		return m, nil

	case searchDebounceMsg:
		return m.handleSearchDebounce(msg)
	}

	return m, nil
}

// postsToRows converts posts to table rows
func (m Model) postsToRows() []table.Row {
	rows := make([]table.Row, len(m.posts))
	for i, p := range m.posts {
		rows[i] = postToRow(p)
	}
	return rows
}

// tagsToRows converts tags to table rows
func (m Model) tagsToRows() []table.Row {
	rows := make([]table.Row, len(m.tags))
	for i, t := range m.tags {
		rows[i] = m.tagToRow(t)
	}
	return rows
}

// feedsToRows converts feeds to table rows
func (m Model) feedsToRows() []table.Row {
	rows := make([]table.Row, len(m.feeds))
	for i, f := range m.feeds {
		rows[i] = m.feedToRow(f)
	}
	return rows
}

// tagToRow converts a single tag to a table row with statistics
func (m Model) tagToRow(t services.TagInfo) table.Row {
	// Tag name (truncate to 28 chars to leave room for selection indicator)
	name := t.Name
	if len(name) > 28 {
		name = name[:25] + "..."
	}

	// Count
	count := fmt.Sprintf("%d", t.Count)

	// Calculate statistics for this tag
	totalWords := 0
	totalReadingTime := 0

	// Get posts for this tag to calculate stats
	for _, post := range m.posts {
		hasTag := false
		for _, tag := range post.Tags {
			if tag == t.Name {
				hasTag = true
				break
			}
		}
		if hasTag {
			if wc, ok := post.Extra["word_count"].(int); ok {
				totalWords += wc
			}
			if rt, ok := post.Extra["reading_time"].(int); ok {
				totalReadingTime += rt
			}
		}
	}

	// Format statistics
	wordsStr := formatWordCount(totalWords)
	readTimeStr := formatReadingTime(totalReadingTime)

	// Slug
	slug := t.Slug
	if slug == "" {
		slug = t.Name
	}

	return table.Row{name, count, wordsStr, readTimeStr, slug}
}

// feedToRow converts a single feed to a table row
func (m Model) feedToRow(f *lifecycle.Feed) table.Row {
	// Name comes first so the feed remains recognizable even when status labels
	// are present. The table clips the complete cell to its column width.
	name := f.Name
	labels := make([]string, 0, 3)
	if f.IncludePrivate {
		labels = append(labels, "private")
	}
	if f.Hidden {
		labels = append(labels, "hidden")
	}
	if f.Automated {
		labels = append(labels, "auto")
	}
	if len(name) > 18 {
		name = name[:15] + "..."
	}
	if len(labels) > 0 {
		name += " [" + strings.Join(labels, ",") + "]"
	}
	name = m.styleFeedTitle(f, name)

	// Posts count
	postsCount := fmt.Sprintf("%d", len(f.Posts))

	// Calculate feed statistics from posts
	totalWords, totalReadingTime := calculateFeedStats(f.Posts)
	avgReadingTime := 0
	if len(f.Posts) > 0 {
		avgReadingTime = totalReadingTime / len(f.Posts)
	}

	// Format statistics
	wordsStr := formatWordCount(totalWords)
	totTimeStr := formatReadingTime(totalReadingTime)
	avgTimeStr := formatReadingTime(avgReadingTime)

	// Output path
	output := f.Path

	return table.Row{name, postsCount, wordsStr, totTimeStr, avgTimeStr, output}
}

// styleFeedTitle gives each non-standard feed state a distinct visual signal.
// The textual labels remain present so the distinction survives monochrome
// terminals and color-disabled output.
func (m Model) styleFeedTitle(feed *lifecycle.Feed, title string) string {
	style := lipgloss.NewStyle()
	switch {
	case feed.IncludePrivate:
		style = style.Foreground(m.getTheme().Colors.FeedPrivate)
	case feed.Hidden:
		style = style.Foreground(m.getTheme().Colors.FeedHidden)
	case feed.Automated:
		style = style.Foreground(m.getTheme().Colors.FeedAutomated)
	default:
		return title
	}
	return style.Bold(true).Render(title)
}

// postToRow converts a single post to a table row
func postToRow(p *models.Post) table.Row {
	// Title (truncate to 33 chars to leave room for selection indicator)
	title := "(untitled)"
	if p.PlainTitle() != "" {
		title = p.PlainTitle()
	}
	if len(title) > 33 {
		title = title[:30] + "..."
	}

	// Date (YYYY-MM-DD format)
	date := ""
	if p.Date != nil {
		date = p.Date.Format("2006-01-02")
	}

	// Word count (from Extra field, populated by StatsPlugin)
	words := ""
	if wc, ok := p.Extra["word_count"].(int); ok {
		words = formatWordCount(wc)
	}

	// Reading time (from Extra field, populated by StatsPlugin)
	readTime := ""
	if rt, ok := p.Extra["reading_time"].(int); ok {
		readTime = formatReadingTime(rt)
	}

	// Tags (truncate to 16 chars)
	tags := strings.Join(p.Tags, ", ")
	if len(tags) > 16 {
		tags = tags[:13] + "..."
	}

	// Path (will be truncated by table column width)
	path := p.Path

	return table.Row{title, date, words, readTime, tags, path}
}

// formatWordCount formats a word count in a human-readable format (e.g., "1.5k")
func formatWordCount(count int) string {
	if count < 1000 {
		return fmt.Sprintf("%d", count)
	}
	if count < 10000 {
		// Show one decimal place for 1k-9.9k
		return fmt.Sprintf("%.1fk", float64(count)/1000)
	}
	// Round to nearest k for 10k+
	return fmt.Sprintf("%dk", count/1000)
}

// formatReadingTime formats reading time in a compact format (e.g., "2 min")
func formatReadingTime(minutes int) string {
	if minutes == 0 {
		return "<1 min"
	}
	if minutes == 1 {
		return "1 min"
	}
	if minutes < 60 {
		return fmt.Sprintf("%d min", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle sort menu keys when visible
	if m.showSortMenu {
		return m.handleSortMenuKey(msg)
	}

	// Handle mode-specific input
	switch m.mode {
	case ModeFilter:
		return m.handleFilterMode(msg)
	case ModeCommand:
		return m.handleCommandMode(msg)
	case ModeNormal:
		// Fall through to normal mode handling below
	}

	// Handle detail view keys separately
	if m.view == ViewPostDetail {
		return m.handleDetailViewKey(msg)
	}

	// Handle help view keys separately (for scrolling and search)
	if m.view == ViewHelp {
		return m.handleHelpViewKey(msg)
	}

	// Handle config view keys separately (for navigation and scrolling)
	if m.view == ViewConfig {
		return m.handleConfigViewKey(msg)
	}

	// Handle search view keys separately
	if m.view == ViewSearch {
		return m.handleSearchViewKey(msg)
	}

	// Normal mode key handling
	return m.handleNormalModeKey(msg)
}

// handleMouse handles mouse events
func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	// Track mouse position for hover effects
	m.mouseX = msg.X
	m.mouseY = msg.Y

	// Only handle mouse in normal mode
	if m.mode != ModeNormal {
		return m, nil
	}

	// Don't handle mouse in sort menu, filter, or command mode
	if m.showSortMenu {
		return m, nil
	}

	// Handle wheel scrolling
	if msg.Action == tea.MouseActionPress {
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			// Navigate up (same as 'k' key)
			return m.handleMouseNavigation(true)

		case tea.MouseButtonWheelDown:
			// Navigate down (same as 'j' key)
			return m.handleMouseNavigation(false)

		case tea.MouseButtonLeft:
			// Check if click is on footer button
			if cmd := m.handleFooterClick(msg.X, msg.Y); cmd != nil {
				return m, cmd
			}

			// Click to select (same as Enter key)
			if m.view == ViewPostDetail {
				// In detail view, clicking does nothing (could implement back navigation)
				return m, nil
			}
			return m.handleEnter()

		default:
			// Ignore other mouse buttons
			return m, nil
		}
	}

	return m, nil
}

// handleFooterClick checks if a click occurred on a footer button and triggers the action
func (m *Model) handleFooterClick(x, y int) tea.Cmd {
	// Calculate footer row (it's at the bottom of the screen)
	// Footer is the last line: height - 1 (0-indexed)
	footerY := m.height - 1

	// Check if click is on footer row
	if y != footerY {
		return nil
	}

	// Check each button to see if click is within its bounds
	for _, btn := range m.footerButtons {
		if x >= btn.startX && x <= btn.endX {
			// Click is on this button, trigger its action
			newModel, cmd := btn.action(m)
			if model, ok := newModel.(Model); ok {
				*m = model
			}
			return cmd
		}
	}

	return nil
}

// handleMouseNavigation handles mouse wheel scrolling
func (m Model) handleMouseNavigation(up bool) (tea.Model, tea.Cmd) {
	// Create a simulated key message for up/down
	var keyMsg tea.KeyMsg
	if up {
		keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	} else {
		keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	}

	// Let the table handle navigation when in posts view
	if m.view == ViewPosts {
		var cmd tea.Cmd
		m.postsTable, cmd = m.postsTable.Update(keyMsg)
		m.cursor = m.postsTable.Cursor()
		return m, cmd
	}
	// Let the table handle navigation when in tags view
	if m.view == ViewTags {
		var cmd tea.Cmd
		m.tagsTable, cmd = m.tagsTable.Update(keyMsg)
		m.cursor = m.tagsTable.Cursor()
		return m, cmd
	}
	// Let the table handle navigation when in feeds view
	if m.view == ViewFeeds {
		var cmd tea.Cmd
		m.feedsTable, cmd = m.feedsTable.Update(keyMsg)
		m.cursor = m.feedsTable.Cursor()
		return m, cmd
	}
	return m, nil
}

// handleNormalModeKey handles key input in normal mode
func (m Model) handleNormalModeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, keyMap.Quit):
		return m.handleQuitKey()

	case key.Matches(msg, keyMap.Up), key.Matches(msg, keyMap.Down):
		return m.handleNavigation(msg)

	case key.Matches(msg, keyMap.Filter):
		if m.view != ViewPosts && m.view != ViewFeeds {
			return m, nil
		}
		m.mode = ModeFilter
		m.filterBeforeInput = m.filter
		m.filterInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, keyMap.Command):
		m.mode = ModeCommand
		m.cmdInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, keyMap.Help):
		m.view = ViewHelp
		m.initializeHelpViewport()
		return m, nil

	case key.Matches(msg, keyMap.Posts):
		return m.handlePostsKey()

	case key.Matches(msg, keyMap.Tags):
		m.view = ViewTags
		m.cursor = 0
		m.tagsTable.SetCursor(0)
		return m, m.loadTags()

	case key.Matches(msg, keyMap.Feeds):
		m.view = ViewFeeds
		m.filter = ""
		m.filterInput.SetValue("")
		m.cursor = 0
		m.feedsTable.SetCursor(0)
		return m, m.loadFeeds()

	case key.Matches(msg, keyMap.Config):
		return m.handleConfigKey()

	case key.Matches(msg, keyMap.Search):
		m.initSearchView()
		m.view = ViewSearch
		m.searchInput.Focus()
		return m, textinput.Blink

	case key.Matches(msg, keyMap.Enter):
		return m.handleEnter()

	case key.Matches(msg, keyMap.Escape):
		return m.handleEscape()

	case key.Matches(msg, keyMap.Edit):
		return m.handleEditKey()

	case key.Matches(msg, keyMap.Sort):
		return m.handleSortKey()

	case key.Matches(msg, keyMap.Refresh):
		return m.handleRefreshKey()

	default:
		// Handle capital letter hotkeys for sorting (k9s-inspired)
		if field, ok := sortHotkeyMap[msg.String()]; ok {
			return m.handleSortHotkey(field)
		}
	}

	return m, nil
}

// handleQuitKey handles the quit key, clearing active filter first if present
func (m Model) handleQuitKey() (tea.Model, tea.Cmd) {
	// If there's an active filter in posts view, clear it first instead of quitting
	if m.view == ViewPosts && m.activeFilter != nil {
		m.activeFilter = nil
		m.cursor = 0
		m.postsTable.SetCursor(0)
		return m, m.loadPosts()
	}
	return m, tea.Quit
}

// handlePostsKey handles navigation to posts view
func (m Model) handlePostsKey() (tea.Model, tea.Cmd) {
	m.view = ViewPosts
	m.cursor = 0
	m.postsTable.SetCursor(0)
	m.activeFilter = nil // Clear any active filter when explicitly navigating to posts
	m.filter = ""
	m.filterInput.SetValue("")
	return m, m.loadPosts()
}

// handleEditKey handles the edit key for editing posts
func (m Model) handleEditKey() (tea.Model, tea.Cmd) {
	if m.view == ViewPosts {
		return m, m.openInEditor()
	}
	return m, nil
}

// handleConfigKey handles navigation to config view
func (m Model) handleConfigKey() (tea.Model, tea.Cmd) {
	m.view = ViewConfig
	m.configCursor = 0
	m.configFilter = ""
	m.buildConfigSections()
	m.initializeConfigViewport()
	return m, nil
}

func (m Model) handleSortKey() (tea.Model, tea.Cmd) {
	if m.view == ViewPosts {
		m.showSortMenu = true
		// Set sortMenuIdx to current sort field
		for i, opt := range sortOptions {
			if opt.value == m.sortBy {
				m.sortMenuIdx = i
				break
			}
		}
	}
	return m, nil
}

// handleSortHotkey handles capital letter hotkeys for direct column sorting (k9s-inspired)
func (m Model) handleSortHotkey(field string) (tea.Model, tea.Cmd) {
	// Only allow sorting in posts view in normal mode
	if m.view != ViewPosts || m.mode != ModeNormal {
		return m, nil
	}

	// If already sorting by this field, toggle the order
	if m.sortBy == field {
		if m.sortOrder == services.SortAsc {
			m.sortOrder = services.SortDesc
		} else {
			m.sortOrder = services.SortAsc
		}
	} else {
		// New field: set to descending by default
		m.sortBy = field
		m.sortOrder = services.SortDesc
	}
	m.postsTable.SetColumns(postsTableColumns(m.width, m.theme, m.sortBy))

	// Reset cursor and reload posts
	m.cursor = 0
	m.postsTable.SetCursor(0)
	return m, m.loadPosts()
}

// handleRefreshKey handles the refresh key for manually refreshing data
func (m Model) handleRefreshKey() (tea.Model, tea.Cmd) {
	if m.refreshing {
		return m, nil // Already refreshing
	}
	return m, m.refreshData()
}

func (m Model) handleNavigation(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Let the table handle navigation when in posts view
	if m.view == ViewPosts {
		var cmd tea.Cmd
		m.postsTable, cmd = m.postsTable.Update(msg)
		m.cursor = m.postsTable.Cursor()
		return m, cmd
	}
	// Let the table handle navigation when in tags view
	if m.view == ViewTags {
		var cmd tea.Cmd
		m.tagsTable, cmd = m.tagsTable.Update(msg)
		m.cursor = m.tagsTable.Cursor()
		return m, cmd
	}
	// Let the table handle navigation when in feeds view
	if m.view == ViewFeeds {
		var cmd tea.Cmd
		m.feedsTable, cmd = m.feedsTable.Update(msg)
		m.cursor = m.feedsTable.Cursor()
		return m, cmd
	}
	return m, nil
}

func (m Model) handleFilterMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.view == ViewPosts && m.visiblePosts == nil {
		m.visiblePosts = m.posts
	}
	if m.view == ViewFeeds && m.visibleFeeds == nil {
		m.visibleFeeds = m.feeds
	}
	switch msg.Type {
	case tea.KeyEscape:
		m.mode = ModeNormal
		m.filterInput.Blur()
		m.filter = m.filterBeforeInput
		m.filterInput.SetValue(m.filter)
		m.applyLocalFilter()
		return m, nil

	case tea.KeyEnter:
		m.filter = m.filterInput.Value()
		m.mode = ModeNormal
		m.filterInput.Blur()
		m.cursor = 0
		m.postsTable.SetCursor(0)
		m.applyLocalFilter()
		return m, nil

	default:
		// Handle other keys through the text input
	}

	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)
	// Slash is a quick local filter, not the structured expression language.
	m.filter = m.filterInput.Value()
	m.applyLocalFilter()
	return m, cmd
}

func (m *Model) applyLocalFilter() {
	if m.view == ViewFeeds {
		m.feeds = filterVisibleFeeds(m.visibleFeeds, m.filter)
		m.feedsTable.SetCursor(0)
		m.feedsTable.SetRows(m.feedsToRows())
		m.cursor = m.feedsTable.Cursor()
		return
	}
	m.posts = filterVisiblePosts(m.visiblePosts, m.filter)
	m.cursor = 0
	m.postsTable.SetCursor(0)
	m.postsTable.SetRows(m.postsToRows())
}

func (m Model) handleCommandMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		m.mode = ModeNormal
		m.cmdInput.Blur()
		m.cmdInput.SetValue("")
		return m, nil

	case tea.KeyEnter:
		cmd := m.cmdInput.Value()
		m.cmdInput.SetValue("")
		m.mode = ModeNormal
		m.cmdInput.Blur()
		return m.executeCommand(cmd)

	default:
		// Handle other keys through the text input
	}

	var cmd tea.Cmd
	m.cmdInput, cmd = m.cmdInput.Update(msg)
	return m, cmd
}

func (m Model) executeCommand(cmd string) (tea.Model, tea.Cmd) {
	cmd = strings.TrimSpace(cmd)
	switch cmd {
	case "posts", "p":
		m.view = ViewPosts
		m.cursor = 0
		m.postsTable.SetCursor(0)
		m.activeFilter = nil // Clear any active filter when explicitly navigating to posts
		return m, m.loadPosts()
	case "tags", "t":
		m.view = ViewTags
		m.cursor = 0
		m.tagsTable.SetCursor(0)
		return m, m.loadTags()
	case "feeds", "f":
		m.view = ViewFeeds
		m.cursor = 0
		m.feedsTable.SetCursor(0)
		return m, m.loadFeeds()
	case "config", "c":
		return m.handleConfigKey()
	case "q", "quit":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.view {
	case ViewPosts:
		return m.handleEnterPostsList()
	case ViewPostDetail:
		// Already in detail view, Enter does nothing
	case ViewHelp:
		// Return to previous view
		m.view = ViewPosts
	case ViewTags:
		return m.handleEnterTagsList()
	case ViewFeeds:
		return m.handleEnterFeedsList()
	case ViewConfig:
		// Toggle section expansion (handled in handleConfigViewKey)
		return m.toggleConfigSection()
	case ViewSearch:
		// Handled in handleSearchViewKey
	}
	return m, nil
}

// handleEnterPostsList handles the Enter key when viewing the posts list.
func (m Model) handleEnterPostsList() (tea.Model, tea.Cmd) {
	if len(m.posts) == 0 || m.cursor >= len(m.posts) {
		return m, nil
	}

	m.selectedPost = m.posts[m.cursor]
	m.previousView = m.view
	m.view = ViewPostDetail

	// Initialize viewport with post detail content
	m.initializePostViewport()
	return m, nil
}

// handleEnterTagsList handles the Enter key when viewing the tags list.
func (m Model) handleEnterTagsList() (tea.Model, tea.Cmd) {
	if len(m.tags) == 0 || m.cursor >= len(m.tags) {
		return m, nil
	}

	selectedTag := m.tags[m.cursor]
	m.activeFilter = &FilterContext{
		Type: "tag",
		Name: selectedTag.Name,
	}
	m.filter = ""
	m.filterInput.SetValue("")
	m.view = ViewPosts
	m.cursor = 0
	m.postsTable.SetCursor(0)
	return m, m.loadPostsForTag(selectedTag.Name)
}

// handleEnterFeedsList handles the Enter key when viewing the feeds list.
func (m Model) handleEnterFeedsList() (tea.Model, tea.Cmd) {
	if len(m.feeds) == 0 || m.cursor >= len(m.feeds) {
		return m, nil
	}

	selectedFeed := m.feeds[m.cursor]
	m.activeFilter = &FilterContext{
		Type: "feed",
		Name: selectedFeed.Name,
	}
	m.filter = ""
	m.filterInput.SetValue("")
	m.view = ViewPosts
	m.cursor = 0
	m.postsTable.SetCursor(0)
	return m, m.loadPostsForFeed(selectedFeed.Name)
}

// initializePostViewport sets up the viewport for viewing post details.
func (m *Model) initializePostViewport() {
	p := m.selectedPost
	if p == nil {
		return
	}

	// Calculate available width and height for viewport
	width := m.calculateViewportWidth()
	viewportHeight := m.calculateViewportHeight()

	// Build viewport content
	viewportContent := m.buildPostViewportContent(p, width)

	// Initialize viewport
	m.postViewport = viewport.New(width-4, viewportHeight)
	m.postViewport.SetContent(viewportContent)
	m.postViewport.YPosition = 0
}

// calculateViewportWidth returns the appropriate width for the viewport.
func (m Model) calculateViewportWidth() int {
	width := m.width
	if width < 40 {
		width = 80
	}
	if width > 100 {
		width = 100
	}
	return width
}

// calculateViewportHeight returns the appropriate height for the viewport.
func (m Model) calculateViewportHeight() int {
	viewportHeight := m.height - 8
	if viewportHeight < 10 {
		viewportHeight = 10
	}
	return viewportHeight
}

// buildPostViewportContent constructs the content string for the post detail viewport.
func (m Model) buildPostViewportContent(p *models.Post, width int) string {
	var metadata strings.Builder
	theme := m.getTheme()

	// Build metadata section
	m.appendPostMetadata(&metadata, p, theme, width)

	// Add separator
	contentWidth := width - 4
	separator := strings.Repeat("─", contentWidth)

	// Render markdown content
	renderedContent := m.renderPostContent(p, contentWidth)

	return metadata.String() + "\n  " + separator + "\n\n  " + theme.DetailLabelStyle.Render("Content:") + "\n\n" + renderedContent
}

// appendPostMetadata appends post metadata to the provided string builder.
func (m Model) appendPostMetadata(metadata *strings.Builder, p *models.Post, theme *Theme, width int) {
	title := "(untitled)"
	if p.PlainTitle() != "" {
		title = p.PlainTitle()
	}
	fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Title:"), title)
	fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Path:"), p.Path)

	dateStr := "(not set)"
	if p.Date != nil {
		dateStr = p.Date.Format("2006-01-02")
	}
	fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Date:"), dateStr)

	publishedStr := "false"
	if p.Published {
		publishedStr = "true"
	}
	fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Published:"), publishedStr)

	tagsStr := "(none)"
	if len(p.Tags) > 0 {
		tagsStr = strings.Join(p.Tags, ", ")
	}
	fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Tags:"), tagsStr)

	wordCount := 0
	if wc, ok := p.Extra["word_count"].(int); ok {
		wordCount = wc
	} else {
		wordCount = countWords(p.Content)
	}
	fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Words:"), formatNumber(wordCount))

	if rt, ok := p.Extra["reading_time"].(int); ok {
		fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Read Time:"), formatReadingTime(rt))
	}

	if cc, ok := p.Extra["char_count"].(int); ok {
		fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Chars:"), formatNumber(cc))
	}

	if p.Description != nil && *p.Description != "" {
		desc := *p.Description
		contentWidth := width - 4
		if len(desc) > contentWidth-15 {
			desc = desc[:contentWidth-18] + "..."
		}
		fmt.Fprintf(metadata, "  %s  %s\n", theme.DetailLabelStyle.Render("Description:"), desc)
	}
}

// renderPostContent renders the post's markdown content using glamour.
func (m Model) renderPostContent(p *models.Post, contentWidth int) string {
	if p.Content == "" {
		return "  (empty)"
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(contentWidth-4),
	)
	if err != nil {
		return "  " + strings.ReplaceAll(p.Content, "\n", "\n  ")
	}

	rendered, err := renderer.Render(p.Content)
	if err != nil {
		return "  " + strings.ReplaceAll(p.Content, "\n", "\n  ")
	}

	return "  " + strings.ReplaceAll(strings.TrimRight(rendered, "\n"), "\n", "\n  ")
}

func (m Model) handleEscape() (tea.Model, tea.Cmd) {
	switch m.view {
	case ViewPostDetail:
		m.view = m.previousView
		if m.view == "" {
			m.view = ViewPosts
		}
		m.selectedPost = nil
	case ViewHelp:
		m.view = ViewPosts
	case ViewPosts:
		// If there's an active filter, clear it and reload all posts
		if m.activeFilter != nil {
			m.activeFilter = nil
			m.cursor = 0
			m.postsTable.SetCursor(0)
			return m, m.loadPosts()
		}
		if m.filter != "" {
			m.filter = ""
			m.filterInput.SetValue("")
			m.posts = m.visiblePosts
			m.cursor = 0
			m.postsTable.SetCursor(0)
			m.postsTable.SetRows(m.postsToRows())
		}
	case ViewTags:
		// Escape does nothing in the tag list view.
	case ViewFeeds:
		if m.filter != "" {
			m.filter = ""
			m.filterInput.SetValue("")
			m.applyLocalFilter()
		}
	case ViewConfig:
		// Return to posts view
		m.view = ViewPosts
	case ViewSearch:
		// Handled in handleSearchViewKey
	}
	return m, nil
}

func (m Model) handleDetailViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keyMap.Quit):
		return m, tea.Quit

	case key.Matches(msg, keyMap.Escape):
		m.view = m.previousView
		if m.view == "" {
			m.view = ViewPosts
		}
		m.selectedPost = nil
		return m, nil

	case key.Matches(msg, keyMap.Edit):
		// Edit functionality - show coming soon message for now
		// This will be implemented in issue #221
		m.err = fmt.Errorf("edit feature coming soon (see issue #221)")
		return m, nil

	case key.Matches(msg, keyMap.Up), key.Matches(msg, keyMap.Down):
		// Handle viewport scrolling
		m.postViewport, cmd = m.postViewport.Update(msg)
		return m, cmd
	}

	// Also handle page up/down, mouse wheel, etc through viewport
	m.postViewport, cmd = m.postViewport.Update(msg)
	return m, cmd
}

// handleHelpViewKey handles key input when viewing the help screen
func (m Model) handleHelpViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Handle search mode keys
	if m.helpSearchMode {
		switch msg.Type {
		case tea.KeyEscape:
			// Clear search and exit search mode
			m.helpSearchMode = false
			m.helpSearchQuery = ""
			m.helpSearchInput.SetValue("")
			m.helpSearchInput.Blur()
			m.helpMatchedLines = nil
			m.helpCurrentMatch = 0
			m.updateHelpViewportContent()
			return m, nil

		case tea.KeyEnter:
			// Just exit search mode, keep the filter active
			m.helpSearchMode = false
			m.helpSearchInput.Blur()
			return m, nil

		default:
			// Update search input
			m.helpSearchInput, cmd = m.helpSearchInput.Update(msg)
			newQuery := m.helpSearchInput.Value()
			if newQuery != m.helpSearchQuery {
				m.helpSearchQuery = newQuery
				m.filterHelpContent()
				m.updateHelpViewportContent()
			}
			return m, cmd
		}
	}

	// Normal mode in help view
	switch {
	case key.Matches(msg, keyMap.Quit):
		return m, tea.Quit

	case key.Matches(msg, keyMap.Escape):
		// Clear search if active, otherwise go back
		if m.helpSearchQuery != "" {
			m.helpSearchQuery = ""
			m.helpSearchInput.SetValue("")
			m.helpMatchedLines = nil
			m.helpCurrentMatch = 0
			m.updateHelpViewportContent()
			return m, nil
		}
		m.view = ViewPosts
		return m, nil

	case msg.String() == "/":
		// Activate search mode
		m.helpSearchMode = true
		m.helpSearchInput.Focus()
		return m, textinput.Blink

	case msg.String() == "n":
		// Go to next match
		if len(m.helpMatchedLines) > 0 {
			m.helpCurrentMatch++
			if m.helpCurrentMatch >= len(m.helpMatchedLines) {
				m.helpCurrentMatch = 0
			}
			m.scrollToHelpMatch()
		}
		return m, nil

	case msg.String() == "N":
		// Go to previous match
		if len(m.helpMatchedLines) > 0 {
			m.helpCurrentMatch--
			if m.helpCurrentMatch < 0 {
				m.helpCurrentMatch = len(m.helpMatchedLines) - 1
			}
			m.scrollToHelpMatch()
		}
		return m, nil

	case key.Matches(msg, keyMap.Up), key.Matches(msg, keyMap.Down):
		// Handle viewport scrolling
		m.helpViewport, cmd = m.helpViewport.Update(msg)
		return m, cmd
	}

	// Handle other viewport controls
	m.helpViewport, cmd = m.helpViewport.Update(msg)
	return m, cmd
}

// initializeHelpViewport sets up the viewport for viewing help content
func (m *Model) initializeHelpViewport() {
	// Calculate available width and height for viewport
	width := m.width
	if width > 120 {
		width = 120
	}
	if width < 40 {
		width = 80
	}
	viewportHeight := m.height - 8
	if viewportHeight < 10 {
		viewportHeight = 10
	}

	// Build help content
	m.populateHelpContentLines()

	// Initialize viewport
	m.helpViewport = viewport.New(width-4, viewportHeight)
	m.updateHelpViewportContent()
	m.helpViewport.YPosition = 0
}

// populateHelpContentLines creates the help content lines
func (m *Model) populateHelpContentLines() {
	helpText := `Help

Navigation:
  j / ↓      Move down
  k / ↑      Move up
  Enter      Select / view details
  Esc        Cancel / go back / clear filter

Views:
  p          Posts view
  t          Tags view
  f          Feeds view
  S          Search view (full-text search with ranked results)

Drill-Down Navigation:
  Enter      In tags view: show posts with selected tag
             In feeds view: show posts in selected feed
  Esc        Clear active filter, return to all posts

Actions:
  e          Edit selected post in $EDITOR
  s          Sort menu (Date, Title, Word Count, Path)

Modes:
  /          Filter visible posts as you type
  :          Command mode

Filter:
   Press / in the posts view and type to narrow the posts already visible.
   Matches title, path, description, tags, and content. Esc cancels; Enter keeps.

Sort Menu:
  j/k/↑/↓    Navigate sort options
  a          Set ascending order
  d          Set descending order
  Enter      Apply sort
  Esc        Cancel

Commands:
  :posts     Show posts
  :tags      Show tags
  :feeds     Show feeds
  :quit      Exit

Help Search:
  /          Activate search mode
  Esc        Clear search / return
  n          Next match
  N          Previous match

Search View:
  S          Open search from any view
  Tab        Toggle advanced filters panel
  ↑/↓        Navigate results
  Enter      Open selected post / apply filter
  Esc        Return to previous view

Press Esc to return.`

	m.helpContentLines = strings.Split(helpText, "\n")
}

// filterHelpContent filters help content based on search query
func (m *Model) filterHelpContent() {
	m.helpMatchedLines = nil
	m.helpCurrentMatch = 0

	if m.helpSearchQuery == "" {
		return
	}

	query := strings.ToLower(m.helpSearchQuery)
	for i, line := range m.helpContentLines {
		if strings.Contains(strings.ToLower(line), query) {
			m.helpMatchedLines = append(m.helpMatchedLines, i)
		}
	}
}

// updateHelpViewportContent updates the viewport content with current filtering/highlighting
func (m *Model) updateHelpViewportContent() {
	var content strings.Builder

	if m.helpSearchQuery == "" {
		// No search, show all content
		content.WriteString(strings.Join(m.helpContentLines, "\n"))
	} else {
		// Highlight matches
		query := strings.ToLower(m.helpSearchQuery)
		for i, line := range m.helpContentLines {
			if strings.Contains(strings.ToLower(line), query) {
				// Highlight this line
				highlighted := m.highlightSearchMatch(line, query)
				content.WriteString(highlighted)
			} else {
				// Show dimmed
				content.WriteString(m.theme.SubtleStyle.Render(line))
			}
			if i < len(m.helpContentLines)-1 {
				content.WriteString("\n")
			}
		}
	}

	m.helpViewport.SetContent(content.String())
}

// highlightSearchMatch highlights the search query in a line
func (m Model) highlightSearchMatch(line, query string) string {
	lowerLine := strings.ToLower(line)
	lowerQuery := strings.ToLower(query)

	// Find all occurrences
	var result strings.Builder
	lastIndex := 0

	for {
		index := strings.Index(lowerLine[lastIndex:], lowerQuery)
		if index == -1 {
			// No more matches
			result.WriteString(line[lastIndex:])
			break
		}

		actualIndex := lastIndex + index

		// Add text before match
		result.WriteString(line[lastIndex:actualIndex])

		// Add highlighted match
		matchText := line[actualIndex : actualIndex+len(query)]
		highlightStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(m.theme.Colors.SelectedText).
			Background(m.theme.Colors.SelectedBg)
		result.WriteString(highlightStyle.Render(matchText))

		lastIndex = actualIndex + len(query)
	}

	return result.String()
}

// scrollToHelpMatch scrolls the viewport to show the current match
func (m *Model) scrollToHelpMatch() {
	if len(m.helpMatchedLines) == 0 {
		return
	}

	lineNumber := m.helpMatchedLines[m.helpCurrentMatch]

	// Calculate the target Y offset
	// We want to position the matched line in the middle of the viewport
	targetY := lineNumber - (m.helpViewport.Height / 2)
	if targetY < 0 {
		targetY = 0
	}

	maxY := len(m.helpContentLines) - m.helpViewport.Height
	if maxY < 0 {
		maxY = 0
	}
	if targetY > maxY {
		targetY = maxY
	}

	m.helpViewport.YOffset = targetY
}

func (m Model) handleSortMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.sortMenuIdx < len(sortOptions)-1 {
			m.sortMenuIdx++
		}
	case "k", "up":
		if m.sortMenuIdx > 0 {
			m.sortMenuIdx--
		}
	case "a":
		m.sortOrder = services.SortAsc
	case "d":
		m.sortOrder = services.SortDesc
	case "enter":
		m.sortBy = sortOptions[m.sortMenuIdx].value
		m.showSortMenu = false
		m.postsTable.SetColumns(postsTableColumns(m.width, m.theme, m.sortBy))
		m.cursor = 0
		m.postsTable.SetCursor(0)
		return m, m.loadPosts()
	case "esc", "q":
		m.showSortMenu = false
	}
	return m, nil
}

// Commands
func (m Model) loadPosts() tea.Cmd {
	return func() tea.Msg {
		opts := services.ListOptions{
			SortBy:    m.sortBy,
			SortOrder: m.sortOrder,
			// m.filter is the local slash filter and is applied after loading.
			Filter: "",
		}
		posts, err := m.app.Posts.List(context.Background(), opts)
		if err != nil {
			return errMsg{err}
		}
		return postsLoadedMsg{posts}
	}
}

// loadInitialData performs the expensive site load after Bubble Tea has
// started rendering. This lets the TUI show its loading screen immediately
// instead of blocking terminal initialization on disk I/O and decryption.
func (m Model) loadInitialData() tea.Cmd {
	return func() tea.Msg {
		if err := m.app.Build.LoadForTUI(context.Background()); err != nil {
			return errMsg{err}
		}

		posts, err := m.app.Posts.List(context.Background(), services.ListOptions{
			SortBy:    m.sortBy,
			SortOrder: m.sortOrder,
		})
		if err != nil {
			return errMsg{err}
		}
		tags, err := m.app.Tags.List(context.Background())
		if err != nil {
			return errMsg{err}
		}
		feeds, err := m.app.Feeds.List(context.Background())
		if err != nil {
			return errMsg{err}
		}
		return initialDataLoadedMsg{posts: posts, tags: tags, feeds: feeds}
	}
}

func (m Model) loadTags() tea.Cmd {
	return func() tea.Msg {
		tags, err := m.app.Tags.List(context.Background())
		if err != nil {
			return errMsg{err}
		}
		return tagsLoadedMsg{tags}
	}
}

func (m Model) loadFeeds() tea.Cmd {
	return func() tea.Msg {
		feeds, err := m.app.Feeds.List(context.Background())
		if err != nil {
			return errMsg{err}
		}
		return feedsLoadedMsg{feeds}
	}
}

// refreshData reloads all data (posts, tags, and feeds)
func (m Model) refreshData() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			return refreshStartedMsg{}
		},
		func() tea.Msg {
			if err := m.app.Build.LoadForTUI(context.Background()); err != nil {
				return errMsg{err}
			}

			// Load all data in parallel
			opts := services.ListOptions{
				SortBy:    m.sortBy,
				SortOrder: m.sortOrder,
			}

			// If there's an active filter, load the appropriate data
			var posts []*models.Post
			var err error

			if m.activeFilter != nil {
				switch m.activeFilter.Type {
				case "tag":
					posts, err = m.app.Tags.GetPosts(context.Background(), m.activeFilter.Name, opts)
				case "feed":
					posts, err = m.app.Feeds.GetPosts(context.Background(), m.activeFilter.Name, opts)
				default:
					opts.Filter = m.filter
					posts, err = m.app.Posts.List(context.Background(), opts)
				}
			} else {
				opts.Filter = m.filter
				posts, err = m.app.Posts.List(context.Background(), opts)
			}

			if err != nil {
				return errMsg{err}
			}

			tags, err := m.app.Tags.List(context.Background())
			if err != nil {
				return errMsg{err}
			}

			feeds, err := m.app.Feeds.List(context.Background())
			if err != nil {
				return errMsg{err}
			}

			return refreshCompletedMsg{
				posts: posts,
				tags:  tags,
				feeds: feeds,
			}
		},
	)
}

// loadPostsForTag loads posts filtered by a specific tag
func (m Model) loadPostsForTag(tag string) tea.Cmd {
	return func() tea.Msg {
		opts := services.ListOptions{
			SortBy:    m.sortBy,
			SortOrder: m.sortOrder,
		}
		posts, err := m.app.Tags.GetPosts(context.Background(), tag, opts)
		if err != nil {
			return errMsg{err}
		}
		return postsLoadedMsg{posts}
	}
}

// loadPostsForFeed loads posts filtered by a specific feed
func (m Model) loadPostsForFeed(feedName string) tea.Cmd {
	return func() tea.Msg {
		opts := services.ListOptions{
			SortBy:    m.sortBy,
			SortOrder: m.sortOrder,
		}
		posts, err := m.app.Feeds.GetPosts(context.Background(), feedName, opts)
		if err != nil {
			return errMsg{err}
		}
		return postsLoadedMsg{posts}
	}
}

// getSelectedPost returns the currently selected post, or nil if none selected
func (m Model) getSelectedPost() *models.Post {
	if m.view != ViewPosts || len(m.posts) == 0 {
		return nil
	}
	if m.cursor < 0 || m.cursor >= len(m.posts) {
		return nil
	}
	return m.posts[m.cursor]
}

// getEditor returns the editor command to use based on environment variables
func getEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	// Check if vim exists
	if _, err := exec.LookPath("vim"); err == nil {
		return "vim"
	}
	return "nano"
}

// openInEditor opens the selected post in the user's editor
func (m Model) openInEditor() tea.Cmd {
	post := m.getSelectedPost()
	if post == nil {
		return nil
	}

	editor := getEditor()
	c := exec.Command(editor, post.Path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

// getTheme returns the model's theme, or a default theme if nil.
func (m Model) getTheme() *Theme {
	if m.theme == nil {
		return DefaultTheme()
	}
	return m.theme
}

// View renders the UI
func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress q to quit.", m.err)
	}

	var content string
	switch m.view {
	case ViewPosts:
		content = m.renderPosts()
	case ViewTags:
		content = m.renderTags()
	case ViewFeeds:
		content = m.renderFeeds()
	case ViewHelp:
		return m.renderHelp()
	case ViewPostDetail:
		return m.renderPostDetail()
	case ViewConfig:
		return m.renderConfig()
	case ViewSearch:
		content = m.renderSearch()
	}

	rendered := m.renderLayout(content)

	// Overlay sort menu if visible
	if m.showSortMenu {
		rendered = m.renderWithSortMenu(rendered)
	}

	return rendered
}

func (m Model) renderLayout(content string) string {
	// Header
	header := m.theme.HeaderStyle.Render("markata-go")
	header += " " + m.theme.SubtleStyle.Render(fmt.Sprintf("[%s]", m.view))

	// Show active filter in header if present
	if m.activeFilter != nil && m.view == ViewPosts {
		filterLabel := fmt.Sprintf(" → %s: %s", m.activeFilter.Type, m.activeFilter.Name)
		activeFilterStyle := lipgloss.NewStyle().Bold(true).
			Foreground(m.theme.Colors.SelectedText).
			Background(m.theme.Colors.SelectedBg)
		header += " " + activeFilterStyle.Render(filterLabel)
	}

	// Status bar with clickable buttons
	var statusBar string
	if m.loading {
		statusBar = m.theme.SubtleStyle.Render("Loading site data...  q:quit")
	} else {
		switch m.mode {
		case ModeFilter:
			statusBar = "Filter: " + m.filterInput.View()
		case ModeCommand:
			statusBar = ":" + m.cmdInput.View()
		default:
			// Build sort indicator
			sortArrow := "↓"
			if m.sortOrder == services.SortAsc {
				sortArrow = "↑"
			}
			sortIndicator := fmt.Sprintf("[%s%s]", sortArrow, m.sortBy)
			statusBar = m.renderFooter(sortIndicator)
		}
	}

	return fmt.Sprintf("%s\n\n%s\n\n%s", header, content, statusBar)
}

// renderFooter builds the footer with clickable buttons and tracks their positions
func (m *Model) renderFooter(sortIndicator string) string {
	// Reset footer buttons
	m.footerButtons = []footerButton{}

	// Track current position in footer
	currentX := 0
	var footerParts []string

	// Helper function to add a clickable button
	addButton := func(label, key string, action func(*Model) (tea.Model, tea.Cmd)) {
		buttonText := key + ":" + label
		startX := currentX
		endX := currentX + len(buttonText) - 1

		// Create hover style - underline and brighten when mouse is over
		isHover := m.mouseY == m.height-1 && m.mouseX >= startX && m.mouseX <= endX
		var styledText string
		if isHover {
			styledText = lipgloss.NewStyle().
				Foreground(m.theme.Colors.Header).
				Underline(true).
				Render(buttonText)
		} else {
			styledText = m.theme.SubtleStyle.Render(buttonText)
		}

		footerParts = append(footerParts, styledText)
		m.footerButtons = append(m.footerButtons, footerButton{
			label:  label,
			key:    key,
			startX: startX,
			endX:   endX,
			action: action,
		})

		// Update position (add length + 2 spaces for separator)
		currentX = endX + 3
	}

	// Add sort indicator (non-clickable)
	footerParts = append(footerParts, m.theme.SubtleStyle.Render(sortIndicator))
	currentX += len(sortIndicator) + 2

	// Determine which buttons to show based on context
	if m.activeFilter != nil && m.view == ViewPosts {
		// When filter is active, show different set of buttons
		addButton("help", "?", func(model *Model) (tea.Model, tea.Cmd) {
			model.view = ViewHelp
			return *model, nil
		})

		addButton("refresh", "r", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleRefreshKey()
		})

		addButton("sort", "s", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleSortKey()
		})

		addButton("edit", "e", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleEditKey()
		})

		addButton("clear filter", "Esc", func(model *Model) (tea.Model, tea.Cmd) {
			model.activeFilter = nil
			model.cursor = 0
			model.postsTable.SetCursor(0)
			return *model, model.loadPosts()
		})

		addButton("quit", "q", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleQuitKey()
		})
	} else {
		// Normal footer buttons
		addButton("help", "?", func(model *Model) (tea.Model, tea.Cmd) {
			model.view = ViewHelp
			return *model, nil
		})

		addButton("refresh", "r", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleRefreshKey()
		})

		addButton("sort", "s", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleSortKey()
		})

		addButton("edit", "e", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleEditKey()
		})

		addButton("feeds", "f", func(model *Model) (tea.Model, tea.Cmd) {
			model.view = ViewFeeds
			model.cursor = 0
			model.feedsTable.SetCursor(0)
			return *model, model.loadFeeds()
		})

		addButton("filter", "/", func(model *Model) (tea.Model, tea.Cmd) {
			model.mode = ModeFilter
			model.filterInput.Focus()
			return *model, textinput.Blink
		})

		addButton("search", "S", func(model *Model) (tea.Model, tea.Cmd) {
			model.initSearchView()
			model.view = ViewSearch
			model.searchInput.Focus()
			return *model, textinput.Blink
		})

		addButton("cmd", ":", func(model *Model) (tea.Model, tea.Cmd) {
			model.mode = ModeCommand
			model.cmdInput.Focus()
			return *model, textinput.Blink
		})

		addButton("quit", "q", func(model *Model) (tea.Model, tea.Cmd) {
			return model.handleQuitKey()
		})
	}

	// Add refresh status indicator
	refreshStatus := ""
	if m.refreshing {
		refreshStatus = " [Refreshing...]"
	} else if !m.lastRefresh.IsZero() {
		elapsed := time.Since(m.lastRefresh)
		switch {
		case elapsed < time.Minute:
			refreshStatus = fmt.Sprintf(" [%ds ago]", int(elapsed.Seconds()))
		case elapsed < time.Hour:
			refreshStatus = fmt.Sprintf(" [%dm ago]", int(elapsed.Minutes()))
		default:
			refreshStatus = fmt.Sprintf(" [%s]", m.lastRefresh.Format("15:04"))
		}
	}
	if refreshStatus != "" {
		footerParts = append(footerParts, m.theme.SubtleStyle.Render(refreshStatus))
	}

	return strings.Join(footerParts, "  ")
}

func (m Model) renderPosts() string {
	if m.loading {
		return "Loading site data..."
	}
	if len(m.posts) == 0 {
		return "No posts found."
	}

	var sb strings.Builder

	// Render the table with header showing count
	header := fmt.Sprintf("Posts (%d)", len(m.posts))
	sb.WriteString(m.theme.HeaderStyle.Render(header))
	sb.WriteString("\n\n")
	sb.WriteString(m.renderPostsTable())

	return sb.String()
}

// renderPostsTable colors each Shift-letter shortcut after Bubble Table
// has laid out its plain headers. Coloring column titles before the table
// renders them makes its width calculation count ANSI escape bytes as text,
// which can hide or misalign headers.
func (m Model) renderPostsTable() string {
	view := m.postsTable.View()
	lines := strings.Split(view, "\n")
	shortcuts := []struct {
		title    string
		shortcut string
		field    string
	}{
		{title: "TITLE", shortcut: "T", field: "title"},
		{title: "DATE", shortcut: "D", field: "date"},
		{title: "WORDS", shortcut: "W", field: "words"},
		{title: "READ", shortcut: "R", field: "reading_time"},
		{title: "TAGS", shortcut: "G", field: "tags"},
		{title: "PATH", shortcut: "P", field: "path"},
	}

	for i, line := range lines {
		if !strings.Contains(line, "TITLE") {
			continue
		}
		for _, key := range shortcuts {
			searchTitle := key.title
			if key.title == "TAGS" {
				// G is the shortcut for tags, but it is not part of the word
				// TAGS. Reuse two padded cell spaces for a visible G prefix.
				searchTitle += "  "
			}
			index := strings.Index(line, searchTitle)
			if index < 0 {
				continue
			}
			color := m.theme.Colors.FeedAutomated
			if key.field == m.sortBy {
				color = m.theme.Colors.Selected
			}
			style := lipgloss.NewStyle().Bold(true).Foreground(color)
			colored := style.Render(key.shortcut)
			if key.title == "TAGS" {
				line = line[:index] + colored + " " + key.title + line[index+len(searchTitle):]
			} else {
				line = line[:index] + colored + line[index+1:]
			}
		}
		lines[i] = line
		break
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderTags() string {
	if len(m.tags) == 0 {
		return "No tags found."
	}

	var sb strings.Builder

	// Render the table with header showing count
	header := fmt.Sprintf("Tags (%d)", len(m.tags))
	sb.WriteString(m.theme.HeaderStyle.Render(header))
	sb.WriteString("\n\n")
	sb.WriteString(m.tagsTable.View())

	return sb.String()
}

func (m Model) renderFeeds() string {
	if len(m.feeds) == 0 {
		return "No feeds found."
	}

	var sb strings.Builder

	// Render the table with header showing count
	header := fmt.Sprintf("Feeds (%d)", len(m.feeds))
	sb.WriteString(m.theme.HeaderStyle.Render(header))
	sb.WriteString("\n")
	sb.WriteString(m.theme.SubtleStyle.Render("[private] private posts  [hidden] sidebar-hidden  [auto] generated"))
	sb.WriteString("\n\n")
	sb.WriteString(m.feedsTable.View())

	return sb.String()
}

// calculateFeedStats calculates total words and reading time for a feed's posts
func calculateFeedStats(posts []*models.Post) (totalWords, totalReadingTime int) {
	for _, post := range posts {
		if wc, ok := post.Extra["word_count"].(int); ok {
			totalWords += wc
		}
		if rt, ok := post.Extra["reading_time"].(int); ok {
			totalReadingTime += rt
		}
	}
	return totalWords, totalReadingTime
}

func (m Model) renderHelp() string {
	theme := m.getTheme()

	// Calculate available width
	width := m.width
	if width < 40 {
		width = 80 // Default minimum
	}
	if width > 100 {
		width = 100 // Max width for readability
	}

	var sb strings.Builder

	// Header
	header := theme.HeaderStyle.Render("markata-go")
	header += " " + theme.SubtleStyle.Render("[help]")
	sb.WriteString(header)
	sb.WriteString("\n\n")

	// Search input area
	if m.helpSearchMode {
		searchPrompt := "Search: " + m.helpSearchInput.View()
		sb.WriteString(searchPrompt)
		sb.WriteString("\n\n")
	} else if m.helpSearchQuery != "" {
		// Show search status
		matchCount := len(m.helpMatchedLines)
		matchInfo := ""
		if matchCount > 0 {
			matchInfo = fmt.Sprintf("Search: %q - %d matches (match %d/%d) [n/N: next/prev, Esc: clear]",
				m.helpSearchQuery, matchCount, m.helpCurrentMatch+1, matchCount)
		} else {
			matchInfo = fmt.Sprintf("Search: %q - no matches [Esc: clear]", m.helpSearchQuery)
		}
		sb.WriteString(theme.SubtleStyle.Render(matchInfo))
		sb.WriteString("\n\n")
	}

	// Create help box with viewport content
	helpBox := theme.DetailBoxStyle.
		Width(width).
		Render(m.helpViewport.View())
	sb.WriteString(helpBox)
	sb.WriteString("\n")

	// Footer with controls
	var footer string
	switch {
	case m.helpSearchMode:
		footer = theme.SubtleStyle.Render("Esc: cancel search  Enter: apply filter")
	case m.helpSearchQuery != "":
		footer = theme.SubtleStyle.Render(fmt.Sprintf("n/N: next/prev match  Esc: clear search  q: quit  %.0f%%", m.helpViewport.ScrollPercent()*100))
	default:
		footer = theme.SubtleStyle.Render(fmt.Sprintf("/: search  ↑/↓: scroll  Esc: return  q: quit  %.0f%%", m.helpViewport.ScrollPercent()*100))
	}
	sb.WriteString(footer)

	return sb.String()
}

func (m Model) renderPostDetail() string {
	if m.selectedPost == nil {
		return "No post selected."
	}

	theme := m.getTheme()

	// Calculate available width
	width := m.width
	if width < 40 {
		width = 80 // Default minimum
	}
	if width > 100 {
		width = 100 // Max width for readability
	}

	// Status bar with scroll percentage
	statusBar := theme.DetailStatusStyle.
		Width(width).
		Render(fmt.Sprintf("  [↑/↓] scroll  [e]dit  [Esc] back  [q]uit  %.0f%%", m.postViewport.ScrollPercent()*100))

	// Header
	header := theme.HeaderStyle.Render("markata-go")
	header += " " + theme.SubtleStyle.Render("[post_detail]")

	// Create detail box with viewport content
	detailBox := theme.DetailBoxStyle.
		Width(width).
		Render(m.postViewport.View())

	return header + "\n\n" + detailBox + "\n" + statusBar
}

// countWords counts the number of words in a string
func countWords(s string) int {
	if s == "" {
		return 0
	}
	return len(strings.Fields(s))
}

// formatNumber formats a number with comma separators
func formatNumber(n int) string {
	if n < 1000 {
		return fmt.Sprintf("%d", n)
	}

	// Format with commas
	s := fmt.Sprintf("%d", n)
	var result strings.Builder
	length := len(s)

	for i, c := range s {
		if i > 0 && (length-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}

	return result.String()
}

// getContentPreview returns a truncated preview of the content
func getContentPreview(content string, maxChars, maxLines, maxWidth int) string {
	if content == "" {
		return "(empty)"
	}

	// Split into lines
	lines := strings.Split(content, "\n")

	var result strings.Builder
	charCount := 0
	lineCount := 0

	for _, line := range lines {
		if lineCount >= maxLines || charCount >= maxChars {
			break
		}

		// Truncate long lines
		if utf8.RuneCountInString(line) > maxWidth {
			line = string([]rune(line)[:maxWidth-3]) + "..."
		}

		result.WriteString(line)
		result.WriteString("\n")

		charCount += len(line)
		lineCount++
	}

	output := strings.TrimRight(result.String(), "\n")

	// Add ellipsis if content was truncated
	if charCount >= maxChars || lineCount >= maxLines {
		output += "\n..."
	}

	return output
}

// renderWithSortMenu overlays the sort menu on top of the existing content
func (m Model) renderWithSortMenu(content string) string {
	menu := m.renderSortMenu()

	// Split content into lines
	contentLines := strings.Split(content, "\n")
	menuLines := strings.Split(menu, "\n")

	// Calculate position to center the menu
	menuWidth := 23 // Width of the menu box
	menuHeight := len(menuLines)

	startX := (m.width - menuWidth) / 2
	if startX < 0 {
		startX = 0
	}
	startY := (m.height - menuHeight) / 2
	if startY < 2 {
		startY = 2
	}

	// Overlay menu on content
	for i, menuLine := range menuLines {
		contentY := startY + i
		if contentY < len(contentLines) {
			line := contentLines[contentY]
			// Ensure line is long enough
			for len(line) < startX {
				line += " "
			}
			// Insert menu line
			runes := []rune(line)
			menuRunes := []rune(menuLine)
			if startX < len(runes) {
				newLine := string(runes[:startX]) + string(menuRunes)
				if startX+len(menuRunes) < len(runes) {
					newLine += string(runes[startX+len(menuRunes):])
				}
				contentLines[contentY] = newLine
			} else {
				contentLines[contentY] = line + menuLine
			}
		}
	}

	return strings.Join(contentLines, "\n")
}

// renderSortMenu renders the sort menu box
func (m Model) renderSortMenu() string {
	var sb strings.Builder

	// Build menu content
	sb.WriteString("┌─ Sort By ─────────┐\n")

	for i, opt := range sortOptions {
		prefix := "  "
		if i == m.sortMenuIdx {
			prefix = "> "
		}
		// Pad label to fixed width
		label := opt.label
		for len(label) < 16 {
			label += " "
		}
		sb.WriteString(fmt.Sprintf("│ %s%s │\n", prefix, label))
	}

	sb.WriteString("├───────────────────┤\n")

	// Show current order with highlight
	ascStyle := ""
	descStyle := ""
	if m.sortOrder == services.SortAsc {
		ascStyle = "*"
	} else {
		descStyle = "*"
	}
	sb.WriteString(fmt.Sprintf("│ [a]sc%s  [d]esc%s   │\n", ascStyle, descStyle))
	sb.WriteString("│ [Enter] apply     │\n")
	sb.WriteString("└───────────────────┘")

	return m.theme.SortMenuStyle.Render(sb.String())
}

// handleConfigViewKey handles key events in the config view.
func (m Model) handleConfigViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, keyMap.Quit):
		return m, tea.Quit

	case key.Matches(msg, keyMap.Escape):
		m.view = ViewPosts
		return m, nil

	case key.Matches(msg, keyMap.Up):
		if m.configCursor > 0 {
			m.configCursor--
		}
		m.syncConfigViewportToCursor()
		return m, nil

	case key.Matches(msg, keyMap.Down):
		maxCursor := m.getMaxConfigCursor()
		if m.configCursor < maxCursor {
			m.configCursor++
		}
		m.syncConfigViewportToCursor()
		return m, nil

	case key.Matches(msg, keyMap.Enter):
		// Toggle section expansion
		return m.toggleConfigSection()

	case key.Matches(msg, keyMap.Filter):
		// TODO: Implement filter mode for config in future phase
		return m, nil
	}

	// Also handle page up/down, mouse wheel, etc through viewport
	m.configViewport, cmd = m.configViewport.Update(msg)
	return m, cmd
}

// toggleConfigSection toggles expansion of the current section
func (m Model) toggleConfigSection() (tea.Model, tea.Cmd) {
	// Count total items to find which section we're in
	currentLine := 0
	cursorPos := m.configCursor

	for i := range m.configSections {
		section := &m.configSections[i]
		if cursorPos >= currentLine && cursorPos < currentLine+1 {
			// Toggle this section
			section.expanded = !section.expanded
			m.configExpanded[section.key] = section.expanded
			// Rebuild and refresh viewport
			m.buildConfigSections()
			m.refreshConfigViewport()
			return m, nil
		}
		currentLine++ // Section header
		if section.expanded {
			currentLine += len(section.items)
		}
	}
	return m, nil
}

// getMaxConfigCursor returns the maximum valid cursor position in the config view
func (m Model) getMaxConfigCursor() int {
	total := 0
	for _, section := range m.configSections {
		total++ // Section header
		if section.expanded {
			total += len(section.items)
		}
	}
	return total - 1
}

// syncConfigViewportToCursor adjusts viewport scroll to keep cursor visible
func (m *Model) syncConfigViewportToCursor() {
	viewportHeight := m.configViewport.Height
	if m.configCursor < m.configViewport.YOffset {
		m.configViewport.YOffset = m.configCursor
	} else if m.configCursor >= m.configViewport.YOffset+viewportHeight {
		m.configViewport.YOffset = m.configCursor - viewportHeight + 1
	}
}

// buildConfigSections builds the config sections from the current configuration
func (m *Model) buildConfigSections() {
	cfg := m.app.Manager.Config()
	extra := cfg.Extra

	m.configSections = []configSection{}

	// Site Metadata section
	siteItems := []configItem{
		{key: "url", value: getStringFromExtra(extra, "url"), level: 0},
		{key: "title", value: getStringFromExtra(extra, "title"), level: 0},
		{key: "description", value: getStringFromExtra(extra, "description"), level: 0},
		{key: "author", value: getStringFromExtra(extra, "author"), level: 0},
	}
	m.configSections = append(m.configSections, configSection{
		name:     "Site Metadata",
		key:      "site",
		items:    siteItems,
		expanded: m.configExpanded["site"],
	})

	// Directories section
	dirItems := []configItem{
		{key: "output_dir", value: cfg.OutputDir, level: 0},
		{key: "content_dir", value: cfg.ContentDir, level: 0},
		{key: "assets_dir", value: getStringFromExtra(extra, "assets_dir"), level: 0},
		{key: "templates_dir", value: getStringFromExtra(extra, "templates_dir"), level: 0},
	}
	m.configSections = append(m.configSections, configSection{
		name:     "Directories",
		key:      "dirs",
		items:    dirItems,
		expanded: m.configExpanded["dirs"],
	})

	// Theme section
	themeItems := []configItem{}
	if themeMap, ok := extra["theme"].(map[string]interface{}); ok {
		themeItems = append(themeItems,
			configItem{key: "name", value: getStringFromMap(themeMap, "name"), level: 0},
			configItem{key: "palette", value: getStringFromMap(themeMap, "palette"), level: 0},
			configItem{key: "palette_light", value: getStringFromMap(themeMap, "palette_light"), level: 0},
			configItem{key: "palette_dark", value: getStringFromMap(themeMap, "palette_dark"), level: 0},
		)
	} else {
		themeItems = append(themeItems,
			configItem{key: "name", value: "default", level: 0},
			configItem{key: "palette", value: "default-light", level: 0},
		)
	}
	m.configSections = append(m.configSections, configSection{
		name:     "Theme",
		key:      "theme",
		items:    themeItems,
		expanded: m.configExpanded["theme"],
	})

	// Build Options section
	buildItems := []configItem{
		{key: "concurrency", value: fmt.Sprintf("%d", getIntFromExtra(extra, "concurrency")), level: 0},
		{key: "glob_patterns", value: strings.Join(cfg.GlobPatterns, ", "), level: 0},
	}
	m.configSections = append(m.configSections, configSection{
		name:     "Build Options",
		key:      "build",
		items:    buildItems,
		expanded: m.configExpanded["build"],
	})

	// Feeds section
	feedItems := []configItem{}
	if feedsRaw, ok := extra["feeds"].([]interface{}); ok {
		feedItems = append(feedItems, configItem{key: "count", value: fmt.Sprintf("%d feeds configured", len(feedsRaw)), level: 0})
		for i, feedRaw := range feedsRaw {
			if i >= 5 {
				feedItems = append(feedItems, configItem{key: "...", value: fmt.Sprintf("and %d more", len(feedsRaw)-5), level: 1})
				break
			}
			if feedMap, ok := feedRaw.(map[string]interface{}); ok {
				name := getStringFromMap(feedMap, "name")
				if name == "" {
					name = fmt.Sprintf("feed_%d", i+1)
				}
				filter := getStringFromMap(feedMap, "filter")
				if filter == "" {
					filter = "(no filter)"
				}
				feedItems = append(feedItems, configItem{key: name, value: filter, level: 1})
			}
		}
	} else {
		feedItems = append(feedItems, configItem{key: "count", value: "0 feeds configured", level: 0})
	}
	m.configSections = append(m.configSections, configSection{
		name:     "Feeds",
		key:      "feeds",
		items:    feedItems,
		expanded: m.configExpanded["feeds"],
	})

	// Layout section
	layoutItems := []configItem{}
	if layoutMap, ok := extra["layout"].(map[string]interface{}); ok {
		layoutItems = append(layoutItems,
			configItem{key: "type", value: getStringFromMap(layoutMap, "type"), level: 0},
			configItem{key: "max_width", value: getStringFromMap(layoutMap, "max_width"), level: 0},
		)
	} else {
		layoutItems = append(layoutItems, configItem{key: "type", value: "default", level: 0})
	}
	m.configSections = append(m.configSections, configSection{
		name:     "Layout",
		key:      "layout",
		items:    layoutItems,
		expanded: m.configExpanded["layout"],
	})

	// Blogroll section
	blogrollItems := []configItem{}
	if blogrollMap, ok := extra["blogroll"].(map[string]interface{}); ok {
		if enabled, ok := blogrollMap["enabled"].(bool); ok {
			blogrollItems = append(blogrollItems, configItem{key: "enabled", value: fmt.Sprintf("%v", enabled), level: 0})
		}
		if feeds, ok := blogrollMap["feeds"].([]interface{}); ok {
			blogrollItems = append(blogrollItems, configItem{key: "feeds_count", value: fmt.Sprintf("%d feeds", len(feeds)), level: 0})
		}
	} else {
		blogrollItems = append(blogrollItems, configItem{key: "enabled", value: "false", level: 0})
	}
	m.configSections = append(m.configSections, configSection{
		name:     "Blogroll",
		key:      "blogroll",
		items:    blogrollItems,
		expanded: m.configExpanded["blogroll"],
	})
}

// initializeConfigViewport sets up the viewport for the config view
func (m *Model) initializeConfigViewport() {
	width := m.calculateViewportWidth()
	viewportHeight := m.calculateViewportHeight()

	content := m.buildConfigContent()

	m.configViewport = viewport.New(width-4, viewportHeight)
	m.configViewport.SetContent(content)
	m.configViewport.YPosition = 0
}

// refreshConfigViewport updates the viewport content after changes
func (m *Model) refreshConfigViewport() {
	content := m.buildConfigContent()
	m.configViewport.SetContent(content)
}

// buildConfigContent builds the content string for the config viewport
func (m Model) buildConfigContent() string {
	var sb strings.Builder
	theme := m.getTheme()
	currentLine := 0

	for _, section := range m.configSections {
		// Section header
		expandIcon := "▶"
		if section.expanded {
			expandIcon = "▼"
		}
		header := fmt.Sprintf("%s %s", expandIcon, section.name)

		// Highlight if cursor is on this line
		if currentLine == m.configCursor {
			header = theme.SelectedStyle.Render(header)
		} else {
			header = theme.HeaderStyle.Render(header)
		}
		sb.WriteString(header)
		sb.WriteString("\n")
		currentLine++

		// Section items (if expanded)
		if section.expanded {
			for _, item := range section.items {
				indent := strings.Repeat("  ", item.level+1)
				value := item.value
				if value == "" {
					value = "(not set)"
				}
				line := fmt.Sprintf("%s%s: %s", indent, theme.DetailLabelStyle.Render(item.key), value)

				// Highlight if cursor is on this line
				if currentLine == m.configCursor {
					line = theme.SelectedStyle.Render(line)
				}
				sb.WriteString(line)
				sb.WriteString("\n")
				currentLine++
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// renderConfig renders the config view
func (m Model) renderConfig() string {
	theme := m.getTheme()

	// Calculate available width
	width := m.width
	if width < 40 {
		width = 80
	}
	if width > 100 {
		width = 100
	}

	// Status bar with scroll percentage
	statusBar := theme.DetailStatusStyle.
		Width(width).
		Render(fmt.Sprintf("  [↑/↓] scroll  [Enter] expand/collapse  [Esc] back  [q]uit  %.0f%%", m.configViewport.ScrollPercent()*100))

	// Header
	header := theme.HeaderStyle.Render("markata-go")
	header += " " + theme.SubtleStyle.Render("[config]")

	// Create config box with viewport content
	configBox := theme.DetailBoxStyle.
		Width(width).
		Render(m.configViewport.View())

	return header + "\n\n" + configBox + "\n" + statusBar
}

// Helper functions for extracting config values

func getStringFromExtra(extra map[string]interface{}, name string) string {
	if extra == nil {
		return ""
	}
	if v, ok := extra[name].(string); ok {
		return v
	}
	return ""
}

func getStringFromMap(data map[string]interface{}, name string) string {
	if data == nil {
		return ""
	}
	if v, ok := data[name].(string); ok {
		return v
	}
	return ""
}

func getIntFromExtra(extra map[string]interface{}, name string) int {
	if extra == nil {
		return 0
	}
	switch v := extra[name].(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	}
	return 0
}

// CI trigger
