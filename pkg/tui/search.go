package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/search"
	"github.com/WaylonWalker/markata/pkg/services"
)

// searchResultItem holds a post with its bleve relevance score.
type searchResultItem struct {
	post  *models.Post
	score float64
}

// searchResultsMsg is sent when a search query completes.
type searchResultsMsg struct {
	results []searchResultItem
	query   string
}

// searchModeType tracks which search panel is active.
type searchModeType int

const (
	searchModeSimple   searchModeType = iota // Big search box only
	searchModeAdvanced                       // Search box + filter panel
)

// initSearchView initializes search view state on the model.
func (m *Model) initSearchView() {
	if m.searchInput.Placeholder == "" {
		m.searchInput = textinput.New()
		m.searchInput.Placeholder = "Search posts..."
		m.searchInput.CharLimit = 200
		m.searchInput.Width = 60

		m.searchFilterInput = textinput.New()
		m.searchFilterInput.Placeholder = "e.g. published == True and 'go' in tags"
		m.searchFilterInput.CharLimit = 200
		m.searchFilterInput.Width = 60

		m.searchTable = createSearchTableWithTheme(m.width, m.theme)
	}
}

// createSearchTableWithTheme builds a results table with a SCORE column.
func createSearchTableWithTheme(width int, theme *Theme) table.Model {
	pathWidth := width - 8 - 35 - 12 - 8 - 8 - 18 - 10
	if pathWidth < 10 {
		pathWidth = 10
	}
	columns := []table.Column{
		{Title: "SCORE", Width: 8},
		{Title: "TITLE", Width: 35},
		{Title: "DATE", Width: 12},
		{Title: "WORDS", Width: 8},
		{Title: "READ", Width: 8},
		{Title: "TAGS", Width: 18},
		{Title: "PATH", Width: pathWidth},
	}
	return createTableWithTheme(columns, theme)
}

// handleSearchViewKey processes keys while ViewSearch is active.
func (m Model) handleSearchViewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Escape: leave search → posts
	if key.Matches(msg, keyMap.Escape) {
		return m.handleSearchEscape()
	}

	// Ctrl+F toggles fuzzy matching
	if msg.String() == "ctrl+f" {
		m.searchFuzzy = !m.searchFuzzy
		if m.searchInput.Value() != "" {
			return m, m.executeSearch()
		}
		return m, nil
	}

	// Ctrl+L cycles result limit: 20 → 50 → 100 → 200 → 20
	if msg.String() == "ctrl+l" {
		m.searchLimit = nextSearchLimit(m.searchLimit)
		if m.searchInput.Value() != "" {
			return m, m.executeSearch()
		}
		return m, nil
	}

	// Tab toggles advanced panel
	if msg.String() == "tab" {
		return m.handleSearchTab()
	}

	// Enter on search input → select result / apply filter
	if msg.String() == "enter" {
		return m.handleSearchEnter()
	}

	// Navigate results with up/down/j/k when input is not focused on filter
	if !m.searchFilterFocused {
		switch msg.String() {
		case "down", "ctrl+n":
			m.searchTable.MoveDown(1)
			return m, nil
		case "up", "ctrl+p":
			m.searchTable.MoveUp(1)
			return m, nil
		}
	}

	// Forward text input to the focused field
	if m.searchFilterFocused {
		var cmd tea.Cmd
		m.searchFilterInput, cmd = m.searchFilterInput.Update(msg)
		return m, cmd
	}

	// Update the search input and fire a debounced search
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)

	newQuery := m.searchInput.Value()
	if newQuery != m.searchLastQuery {
		m.searchLastQuery = newQuery
		m.searchDebounceTime = time.Now()
		return m, tea.Batch(cmd, m.searchDebounceCmd())
	}
	return m, cmd
}

func (m Model) handleSearchEscape() (tea.Model, tea.Cmd) {
	if m.searchMode == searchModeAdvanced && m.searchFilterFocused {
		m.searchFilterFocused = false
		m.searchFilterInput.Blur()
		m.searchInput.Focus()
		return m, nil
	}
	m.view = ViewPosts
	m.searchInput.Blur()
	m.searchFilterInput.Blur()
	return m, nil
}

func (m Model) handleSearchTab() (tea.Model, tea.Cmd) {
	switch {
	case m.searchMode == searchModeSimple:
		m.searchMode = searchModeAdvanced
	case !m.searchFilterFocused:
		m.searchFilterFocused = true
		m.searchInput.Blur()
		m.searchFilterInput.Focus()
		return m, textinput.Blink
	default:
		m.searchFilterFocused = false
		m.searchFilterInput.Blur()
		m.searchInput.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

func (m Model) handleSearchEnter() (tea.Model, tea.Cmd) {
	if m.searchFilterFocused {
		m.searchAdvancedFilter = m.searchFilterInput.Value()
		return m, m.executeSearch()
	}
	if len(m.searchResults) > 0 {
		selected := m.searchTable.Cursor()
		if selected >= 0 && selected < len(m.searchResults) {
			m.selectedPost = m.searchResults[selected].post
			m.previousView = ViewSearch
			m.view = ViewPostDetail
			m.initializePostViewport()
			return m, nil
		}
	}
	return m, nil
}

func nextSearchLimit(current int) int {
	switch current {
	case 0, 20:
		return 50
	case 50:
		return 100
	case 100:
		return 200
	default:
		return 20
	}
}

// searchDebounceCmd waits 200ms then fires if the query hasn't changed.
func (m Model) searchDebounceCmd() tea.Cmd {
	snapshot := m.searchDebounceTime
	return tea.Tick(200*time.Millisecond, func(_ time.Time) tea.Msg {
		return searchDebounceMsg{snapshot: snapshot}
	})
}

// searchDebounceMsg carries the timestamp to compare against current debounce.
type searchDebounceMsg struct {
	snapshot time.Time
}

// handleSearchDebounce is called from Update when a debounce tick fires.
func (m Model) handleSearchDebounce(msg searchDebounceMsg) (tea.Model, tea.Cmd) {
	if msg.snapshot != m.searchDebounceTime {
		return m, nil // query changed since this tick was scheduled
	}
	if m.searchInput.Value() == "" {
		m.searchResults = nil
		m.searchTable.SetRows(nil)
		return m, nil
	}
	return m, m.executeSearch()
}

// executeSearch runs a bleve (or fallback substring) search.
func (m Model) executeSearch() tea.Cmd {
	queryStr := m.searchInput.Value()
	filterExpr := m.searchAdvancedFilter
	posts := m.posts
	fuzzy := m.searchFuzzy
	limit := m.searchLimit
	if limit <= 0 {
		limit = 50
	}

	return func() tea.Msg {
		// Apply filter expression if set
		var filteredPosts []*models.Post
		if filterExpr != "" {
			opts := m.app.Posts
			listed, err := opts.List(context.Background(), services.ListOptions{Filter: filterExpr})
			if err != nil {
				filteredPosts = posts
			} else {
				filteredPosts = listed
			}
		} else {
			filteredPosts = posts
		}

		if queryStr == "" {
			return searchResultsMsg{results: nil, query: ""}
		}

		// Try bleve
		results, err := searchPostsBleve(filteredPosts, queryStr, fuzzy, limit)
		if err != nil {
			results = searchPostsSubstring(filteredPosts, queryStr, limit)
		}
		return searchResultsMsg{results: results, query: queryStr}
	}
}

func searchPostsBleve(posts []*models.Post, queryStr string, fuzzy bool, limit int) ([]searchResultItem, error) {
	idx, err := search.BuildIfNeeded(".markata/cache", posts)
	if err != nil {
		return nil, err
	}
	defer idx.Close()

	postsByPath := search.PostsByPath(posts)
	hits, err := idx.Search(queryStr, search.QueryOptions{Limit: limit, Fuzzy: fuzzy}, postsByPath)
	if err != nil {
		return nil, err
	}

	results := make([]searchResultItem, len(hits))
	for i := range hits {
		h := &hits[i]
		results[i] = searchResultItem{post: h.Post, score: h.Score}
	}
	return results, nil
}

func searchPostsSubstring(posts []*models.Post, queryStr string, limit int) []searchResultItem {
	q := strings.ToLower(queryStr)
	var results []searchResultItem
	for _, p := range posts {
		if matchesPostSubstring(p, q) {
			results = append(results, searchResultItem{post: p, score: 1.0})
		}
	}
	if len(results) > limit {
		results = results[:limit]
	}
	return results
}

func matchesPostSubstring(p *models.Post, q string) bool {
	if strings.Contains(strings.ToLower(p.PlainTitle()), q) {
		return true
	}
	if p.Description != nil && strings.Contains(strings.ToLower(*p.Description), q) {
		return true
	}
	if strings.Contains(strings.ToLower(p.Content), q) {
		return true
	}
	for _, tag := range p.Tags {
		if strings.Contains(strings.ToLower(tag), q) {
			return true
		}
	}
	return false
}

// filterVisiblePosts narrows an already loaded list for the TUI slash filter.
// Keeping this local makes the filter predictable and prevents it from
// unexpectedly searching posts that are not in the current drill-down view.
func filterVisiblePosts(posts []*models.Post, query string) []*models.Post {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return posts
	}

	filtered := make([]*models.Post, 0, len(posts))
	for _, post := range posts {
		if matchesPostSubstring(post, query) || strings.Contains(strings.ToLower(post.Path), query) {
			filtered = append(filtered, post)
		}
	}
	return filtered
}

func filterVisibleFeeds(feeds []*lifecycle.Feed, query string) []*lifecycle.Feed {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return feeds
	}

	filtered := make([]*lifecycle.Feed, 0, len(feeds))
	for _, feed := range feeds {
		if strings.Contains(strings.ToLower(feed.Name), query) ||
			strings.Contains(strings.ToLower(feed.Title), query) ||
			strings.Contains(strings.ToLower(feed.Path), query) {
			filtered = append(filtered, feed)
		}
	}
	return filtered
}

// handleSearchResults processes search results arriving from executeSearch.
func (m *Model) handleSearchResults(msg searchResultsMsg) {
	m.searchResults = msg.results
	m.searchTable.SetRows(m.searchResultsToRows())
	if len(msg.results) > 0 {
		m.searchTable.SetCursor(0)
	}
}

func (m Model) searchResultsToRows() []table.Row {
	rows := make([]table.Row, len(m.searchResults))
	for i, r := range m.searchResults {
		row := postToRow(r.post)
		score := fmt.Sprintf("%.2f", r.score)
		rows[i] = append(table.Row{score}, row...)
	}
	return rows
}

// renderSearch renders the Google-style search view.
func (m Model) renderSearch() string {
	theme := m.getTheme()
	var sb strings.Builder

	// Centered title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Colors.Header).
		Align(lipgloss.Center).
		Width(m.width)
	sb.WriteString(titleStyle.Render("🔍 Search"))
	sb.WriteString("\n\n")

	// Big centered search box
	inputWidth := 64
	if m.width < 70 {
		inputWidth = m.width - 6
	}
	m.searchInput.Width = inputWidth

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Colors.Border).
		Padding(0, 1).
		Width(inputWidth + 4).
		Align(lipgloss.Left)

	centeredBox := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center)
	sb.WriteString(centeredBox.Render(boxStyle.Render(m.searchInput.View())))
	sb.WriteString("\n")

	// Hint line
	hintStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Subtle).
		Align(lipgloss.Center).
		Width(m.width)

	// Status indicators for fuzzy and limit
	limit := m.searchLimit
	if limit <= 0 {
		limit = 50
	}
	fuzzyLabel := "off"
	if m.searchFuzzy {
		fuzzyLabel = "on"
	}
	status := fmt.Sprintf("Fuzzy: %s • Limit: %d", fuzzyLabel, limit)
	statusStyle := lipgloss.NewStyle().
		Foreground(theme.Colors.Header).
		Align(lipgloss.Center).
		Width(m.width)
	sb.WriteString(statusStyle.Render(status))
	sb.WriteString("\n")

	hints := "Tab: filters • Ctrl+F: fuzzy • Ctrl+L: limit • Esc: back • ↑/↓: navigate • Enter: open"
	sb.WriteString(hintStyle.Render(hints))
	sb.WriteString("\n\n")

	// Advanced filter panel (if active)
	if m.searchMode == searchModeAdvanced {
		filterLabel := "Filter: "
		if m.searchFilterFocused {
			filterLabel = "Filter (editing): "
		}
		filterBox := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(theme.Colors.Border).
			Padding(0, 1).
			Width(inputWidth + 4)
		filterContent := filterLabel + m.searchFilterInput.View()
		if m.searchAdvancedFilter != "" && !m.searchFilterFocused {
			filterContent += "\n  Active: " + m.searchAdvancedFilter
		}
		sb.WriteString(centeredBox.Render(filterBox.Render(filterContent)))
		sb.WriteString("\n\n")
	}

	// Results
	if m.searchInput.Value() != "" {
		countStyle := lipgloss.NewStyle().
			Foreground(theme.Colors.Header).
			Align(lipgloss.Center).
			Width(m.width)
		sb.WriteString(countStyle.Render(fmt.Sprintf("%d results", len(m.searchResults))))
		sb.WriteString("\n\n")

		if len(m.searchResults) > 0 {
			sb.WriteString(m.searchTable.View())
		}
	}

	return sb.String()
}
