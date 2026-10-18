// Package servetui presents a servecontrol session in a terminal.
// The runtime owns build state; this package only keeps navigation state.
package servetui

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/servecontrol"
	"github.com/WaylonWalker/markata/pkg/serveopen"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen string

const (
	screenJobs     screen = "jobs"
	screenJob      screen = "job"
	screenLogs     screen = "logs"
	screenWarnings screen = "warnings"
	screenErrors   screen = "errors"
	screenPages    screen = "pages"
	screenPage     screen = "page"
	screenFeeds    screen = "feeds"
	screenFeed     screen = "feed"
	screenHelp     screen = "help"
)

const (
	keyEsc   = "esc"
	keyEnter = "enter"
)

type snapshotMsg struct{ snapshot servecontrol.Snapshot }
type actionResultMsg struct{ err error }
type paletteMsg struct{ palette *palettes.Palette }
type tickMsg struct{}

var pulseFrames = [...]string{"·", "•", "●", "•", "·"}

type theme struct {
	accent, text, muted, success, warning, failure, focus lipgloss.Color
}

func themeFromPalette(palette *palettes.Palette) theme {
	// These defaults match the bundled default site palette's semantic roles.
	t := theme{accent: "#8462af", text: "#e5e5e5", muted: "#777777", success: "#3c9d73", warning: "#cb903c", failure: "#ce5364", focus: "#49335f"}
	if palette == nil {
		return t
	}
	resolve := func(names []string, fallback lipgloss.Color) lipgloss.Color {
		for _, name := range names {
			if color := palette.Resolve(name); color != "" {
				return lipgloss.Color(color)
			}
		}
		return fallback
	}
	t.accent = resolve([]string{"accent", "link"}, t.accent)
	t.text = resolve([]string{"text-primary", "text"}, t.text)
	t.muted = resolve([]string{"text-muted", "text-secondary"}, t.muted)
	t.success = resolve([]string{"success", "info"}, t.success)
	t.warning = resolve([]string{"warning"}, t.warning)
	t.failure = resolve([]string{"error"}, t.failure)
	t.focus = resolve([]string{"bg-elevated", "accent-hover"}, t.focus)
	return t
}

// Model renders an immutable runtime snapshot. It never reconstructs job or
// diagnostic state from log text.
type Model struct {
	snapshot   servecontrol.Snapshot
	action     func(servecontrol.ActionRequest) error
	view       screen
	navStack   []screen
	sourceRoot string
	jobID      string
	logJobID   string
	page       string
	feedName   string
	selected   int
	width      int
	height     int
	offset     int
	query      string
	search     bool
	follow     bool
	newLines   int
	lastLogs   int
	message    string
	theme      theme
	frame      int
}

// NewModel creates a testable terminal model from a runtime snapshot.
func NewModel(snapshot servecontrol.Snapshot, action func(servecontrol.ActionRequest) error) Model {
	return NewModelWithPalette(snapshot, action, nil)
}

// NewModelWithPalette maps site colors into terminal status and focus colors.
func NewModelWithPalette(snapshot servecontrol.Snapshot, action func(servecontrol.ActionRequest) error, palette *palettes.Palette) Model {
	return Model{snapshot: snapshot, action: action, view: screenJobs, width: 80, height: 24, follow: true, lastLogs: len(snapshot.Logs), theme: themeFromPalette(palette)}
}

func (m Model) Init() tea.Cmd {
	if m.hasRunningWork() {
		return nextTick()
	}
	return nil
}

func nextTick() tea.Cmd {
	return tea.Tick(350*time.Millisecond, func(time.Time) tea.Msg { return tickMsg{} })
}

func (m Model) hasRunningWork() bool {
	if m.snapshot.Site.Status == servecontrol.StateRunning {
		return true
	}
	for i := range m.snapshot.Jobs {
		job := m.snapshot.Jobs[i]
		if job.State == servecontrol.StateRunning {
			return true
		}
	}
	return false
}

// Update handles terminal input and new snapshots without touching the runtime.
//
//nolint:gocyclo // Centralized input routing keeps navigation behavior consistent across resources.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = max(1, msg.Width), max(1, msg.Height)
		m.clampOffset()
	case snapshotMsg:
		wasRunning := m.hasRunningWork()
		added := len(msg.snapshot.Logs) - m.lastLogs
		if added > 0 && !m.follow {
			m.newLines += added
		}
		m.snapshot = msg.snapshot
		m.lastLogs = len(msg.snapshot.Logs)
		if m.follow && m.view == screenLogs {
			m.offset = m.maxOffset()
		} else {
			m.clampOffset()
		}
		if !wasRunning && m.hasRunningWork() {
			return m, nextTick()
		}
	case paletteMsg:
		m.theme = themeFromPalette(msg.palette)
	case tickMsg:
		if m.hasRunningWork() {
			m.frame = (m.frame + 1) % len(pulseFrames)
			return m, nextTick()
		}
	case actionResultMsg:
		if msg.err != nil {
			m.message = msg.err.Error()
		} else {
			m.message = "Action queued"
		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.scroll(-3)
		case tea.MouseButtonWheelDown:
			m.scroll(3)
		default:
		}
	case tea.KeyMsg:
		if m.search {
			switch msg.String() {
			case keyEsc:
				m.search = false
			case keyEnter:
				m.search = false
			case "backspace":
				runes := []rune(m.query)
				if len(runes) > 0 {
					m.query = string(runes[:len(runes)-1])
				}
				m.selected, m.offset = 0, 0
			default:
				if len(msg.Runes) > 0 {
					m.query += string(msg.Runes)
					m.selected, m.offset = 0, 0
				}
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "a":
			if m.snapshot.Server.AdminURL != "" {
				cmd, err := serveopen.BrowserCommand(m.snapshot.Server.AdminURL)
				if err != nil {
					m.message = err.Error()
					return m, nil
				}
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return actionResultMsg{err: err} })
			}
		case "?":
			m.changeView(screenHelp)
		case keyEsc:
			switch {
			case len(m.navStack) > 0:
				m.view = m.navStack[len(m.navStack)-1]
				m.navStack = m.navStack[:len(m.navStack)-1]
			case m.query != "":
				m.query = ""
			default:
				m.view = screenJobs
			}
			m.offset, m.selected = 0, 0
		case "/":
			m.search = true
		case "j", "down":
			m.move(1)
		case "k", "up":
			m.move(-1)
		case "pgdown", "ctrl+d":
			m.scroll(max(1, m.bodyHeight()/2))
		case "pgup", "ctrl+u":
			m.scroll(-max(1, m.bodyHeight()/2))
		case "g", "home":
			m.selected, m.offset, m.follow = 0, 0, false
		case "G", "end":
			m.selected = max(0, m.itemCount()-1)
			m.offset = m.maxOffset()
			m.follow, m.newLines = true, 0
		case keyEnter:
			m.enter()
		case "l":
			if m.view == screenJob {
				m.logJobID = m.jobID
			} else {
				m.logJobID = ""
			}
			m.changeView(screenLogs)
			m.follow, m.newLines = true, 0
			m.offset = m.maxOffset()
		case "w":
			m.changeView(screenWarnings)
		case "e":
			m.changeView(screenErrors)
		case "p":
			//nolint:gocritic // The action is shared by both diagnostic inboxes.
			if m.view == screenWarnings || m.view == screenErrors {
				if d, ok := m.selectedDiagnostic(); ok && d.Page != "" {
					m.page = d.Page
					m.changeView(screenPage)
				}
			} else if m.view == screenJob {
				for i := range m.snapshot.Jobs {
					job := m.snapshot.Jobs[i]
					if job.ID == m.jobID && len(job.Pages) > 0 {
						m.page = job.Pages[0]
						m.changeView(screenPage)
						break
					}
				}
			} else {
				m.changeView(screenPages)
			}
		case "f":
			m.changeView(screenFeeds)
		case "b":
			if m.view == screenWarnings || m.view == screenErrors || m.view == screenPage {
				jobID := ""
				if d, ok := m.selectedDiagnostic(); ok {
					jobID = d.JobID
				}
				if jobID == "" && m.view == screenPage {
					for _, page := range m.snapshot.Pages {
						if page.Path == m.page {
							jobID = page.LastJobID
							break
						}
					}
				}
				if jobID != "" {
					m.jobID = jobID
					m.changeView(screenJob)
				}
			}
		case "t":
			return m, m.trigger(servecontrol.ActionRequest{Kind: "build"})
		case "r":
			if m.view == screenJob && m.jobID != "" {
				return m, m.trigger(servecontrol.ActionRequest{Kind: "rerun", JobID: m.jobID})
			}
			if m.view == screenJobs {
				if job, ok := m.selectedJob(); ok {
					return m, m.trigger(servecontrol.ActionRequest{Kind: "rerun", JobID: job.ID})
				}
			}
		case "R":
			page := ""
			if m.view == screenPage {
				page = m.page
			}
			if m.view == screenPages {
				pages := m.pages()
				if m.selected >= 0 && m.selected < len(pages) {
					page = pages[m.selected]
				}
			}
			if m.rebuildablePage(page) {
				return m, m.trigger(servecontrol.ActionRequest{Kind: "rebuild-page", Page: page})
			}
		case "o":
			path, line := m.selectedSource()
			if path != "" {
				cmd, err := serveopen.EditorCommand(path, line)
				if err != nil {
					m.message = err.Error()
					return m, nil
				}
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return actionResultMsg{err: err} })
			}
		case "v":
			if page, ok := m.selectedPage(); ok && page.URL != "" && m.snapshot.Server.Address != "" {
				base := "http://" + m.snapshot.Server.Address
				preview, err := url.JoinPath(base, page.URL)
				if err != nil {
					m.message = err.Error()
					return m, nil
				}
				cmd, err := serveopen.BrowserCommand(preview)
				if err != nil {
					m.message = err.Error()
					return m, nil
				}
				return m, tea.ExecProcess(cmd, func(err error) tea.Msg { return actionResultMsg{err: err} })
			}
		}
	}
	return m, nil
}

func (m Model) trigger(req servecontrol.ActionRequest) tea.Cmd {
	if m.action == nil {
		return nil
	}
	return func() tea.Msg { return actionResultMsg{err: m.action(req)} }
}

func (m *Model) changeView(view screen) {
	if view == m.view {
		return
	}
	m.navStack = append(m.navStack, m.view)
	m.view, m.selected, m.offset = view, 0, 0
}

func (m *Model) enter() {
	switch m.view {
	case screenFeeds:
		feeds := m.filteredFeeds()
		if m.selected < len(feeds) {
			m.feedName = feeds[m.selected].Name
			m.query = ""
			m.changeView(screenFeed)
		}
	case screenFeed:
		feed, ok := m.selectedFeed()
		if ok {
			entries := m.filteredFeedEntries(feed)
			if m.selected < len(entries) {
				m.page = entries[m.selected].Path
				m.query = ""
				m.changeView(screenPage)
			}
		}
	case screenJobs:
		if job, ok := m.selectedJob(); ok {
			m.jobID = job.ID
			m.changeView(screenJob)
		}
	case screenPages:
		pages := m.pages()
		if m.selected < len(pages) {
			m.page = pages[m.selected]
			m.changeView(screenPage)
		}
	case screenWarnings, screenErrors:
		if d, ok := m.selectedDiagnostic(); ok && d.Page != "" {
			m.page = d.Page
			m.changeView(screenPage)
		}
	default:
	}
}

func (m Model) filteredJobs() []servecontrol.Job {
	jobs := make([]servecontrol.Job, 0, len(m.snapshot.Jobs))
	for i := range m.snapshot.Jobs {
		job := m.snapshot.Jobs[i]
		if m.match(job.ID + " " + job.Name + " " + job.Type + " " + fmt.Sprint(job.State)) {
			jobs = append(jobs, job)
		}
	}
	return jobs
}

func (m Model) filteredFeeds() []servecontrol.Feed {
	feeds := make([]servecontrol.Feed, 0, len(m.snapshot.Feeds))
	for _, feed := range m.snapshot.Feeds {
		if m.match(feed.Name + " " + feed.Title + " " + feed.Path) {
			feeds = append(feeds, feed)
		}
	}
	sort.SliceStable(feeds, func(i, j int) bool { return feedPathOrName(feeds[i]) < feedPathOrName(feeds[j]) })
	return feeds
}

func feedPathOrName(feed servecontrol.Feed) string {
	name := feed.Path
	if name == "" {
		name = feed.Name
	}
	name = strings.Trim(name, "/")
	if name == "" {
		return "/"
	}
	return "/" + name + "/"
}

func (m Model) selectedFeed() (servecontrol.Feed, bool) {
	for _, feed := range m.snapshot.Feeds {
		if feed.Name == m.feedName {
			return feed, true
		}
	}
	return servecontrol.Feed{}, false
}

func (m Model) filteredFeedEntries(feed servecontrol.Feed) []servecontrol.FeedEntry {
	entries := make([]servecontrol.FeedEntry, 0, len(feed.Entries))
	for _, entry := range feed.Entries {
		if m.match(entry.Path + " " + entry.Title + " " + entry.URL) {
			entries = append(entries, entry)
		}
	}
	return entries
}

func (m Model) filteredDiagnostics() []servecontrol.Diagnostic {
	source := m.snapshot.Diagnostics
	if m.view == screenPage {
		source = nil
		for _, page := range m.snapshot.Pages {
			if page.Path == m.page {
				source = page.Diagnostics
				break
			}
		}
	}
	diags := make([]servecontrol.Diagnostic, 0, len(source))
	for i := range source {
		d := source[i]
		if m.view == screenWarnings && !strings.EqualFold(d.Severity, "warning") {
			continue
		}
		if m.view == screenErrors && !strings.EqualFold(d.Severity, "error") {
			continue
		}
		if m.match(d.Code + " " + d.Message + " " + d.Page + " " + d.File) {
			diags = append(diags, d)
		}
	}
	return diags
}

func (m Model) selectedJob() (servecontrol.Job, bool) {
	jobs := m.filteredJobs()
	if m.selected >= 0 && m.selected < len(jobs) {
		return jobs[m.selected], true
	}
	return servecontrol.Job{}, false
}

func (m Model) selectedDiagnostic() (servecontrol.Diagnostic, bool) {
	diags := m.filteredDiagnostics()
	if m.selected >= 0 && m.selected < len(diags) {
		return diags[m.selected], true
	}
	return servecontrol.Diagnostic{}, false
}

func (m Model) selectedSource() (path string, line int) {
	if m.view == screenWarnings || m.view == screenErrors || m.view == screenPage {
		if diagnostic, ok := m.selectedDiagnostic(); ok {
			path := diagnostic.File
			if path == "" {
				path = diagnostic.Page
			}
			return m.resolveSource(path), diagnostic.Line
		}
	}
	if page, ok := m.selectedPage(); ok {
		return m.resolveSource(page.Path), 0
	}
	return "", 0
}

func (m Model) resolveSource(path string) string {
	if path == "" || m.sourceRoot == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(m.sourceRoot, path)
}

func (m Model) selectedPage() (servecontrol.Page, bool) {
	path := m.page
	if m.view == screenPages {
		pages := m.pages()
		if m.selected >= 0 && m.selected < len(pages) {
			path = pages[m.selected]
		}
	}
	for _, page := range m.snapshot.Pages {
		if page.Path == path {
			return page, true
		}
	}
	return servecontrol.Page{}, false
}

func (m Model) pages() []string {
	seen := map[string]bool{}
	for _, page := range m.snapshot.Pages {
		if page.Path != "" {
			seen[page.Path] = true
		}
	}
	pages := make([]string, 0, len(seen))
	for page := range seen {
		pages = append(pages, page)
	}
	sort.Strings(pages)
	filtered := pages[:0]
	for _, page := range pages {
		if m.match(page) {
			filtered = append(filtered, page)
		}
	}
	return filtered
}

func (m Model) rebuildablePage(path string) bool {
	if path == "" || !(strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown")) {
		return false
	}
	for _, page := range m.snapshot.Pages {
		if page.Path == path {
			return true
		}
	}
	return false
}

func (m Model) match(value string) bool {
	return m.query == "" || strings.Contains(strings.ToLower(value), strings.ToLower(m.query))
}

func (m Model) itemCount() int {
	switch m.view {
	case screenJobs:
		return len(m.filteredJobs())
	case screenFeeds:
		return len(m.filteredFeeds())
	case screenFeed:
		if feed, ok := m.selectedFeed(); ok {
			return len(m.filteredFeedEntries(feed))
		}
		return 0
	case screenWarnings, screenErrors, screenPage:
		return len(m.filteredDiagnostics())
	case screenPages:
		return len(m.pages())
	default:
		return len(m.bodyLines())
	}
}

func (m *Model) move(delta int) {
	if m.view == screenLogs || m.view == screenJob || m.view == screenHelp {
		m.scroll(delta)
		return
	}
	m.selected = min(max(0, m.selected+delta), max(0, m.itemCount()-1))
	if m.selected < m.offset {
		m.offset = m.selected
	}
	if m.selected >= m.offset+m.bodyHeight() {
		m.offset = m.selected - m.bodyHeight() + 1
	}
}

func (m *Model) scroll(delta int) {
	m.offset = min(max(0, m.offset+delta), m.maxOffset())
	if m.view == screenLogs {
		m.follow = m.offset == m.maxOffset()
		if m.follow {
			m.newLines = 0
		}
	}
}

func (m *Model) clampOffset()   { m.offset = min(m.offset, m.maxOffset()) }
func (m Model) maxOffset() int  { return max(0, len(m.bodyLines())-m.bodyHeight()) }
func (m Model) bodyHeight() int { return max(1, m.height-7) }

//nolint:gocyclo // Rendering stays grouped by resource to keep each view's context together.
func (m Model) bodyLines() []string {
	lines := []string{}
	switch m.view {
	case screenJobs:
		if m.snapshot.Site.Status == servecontrol.StateFailed && m.snapshot.Site.Message != "" {
			lines = append(lines, "✗ Build failed: "+m.snapshot.Site.Message, "")
		} else if m.snapshot.Server.Status == servecontrol.StateFailed && m.snapshot.Server.Message != "" {
			lines = append(lines, "✗ "+m.snapshot.Server.Message, "")
		}
		jobs := m.filteredJobs()
		for i := range jobs {
			job := jobs[i]
			mark := "  "
			if i == m.selected {
				mark = "> "
			}
			lines = append(lines, fmt.Sprintf("%s%s %-12s %s", mark, m.symbol(job.State), job.State, job.Name))
		}
	case screenJob:
		for i := range m.snapshot.Jobs {
			job := m.snapshot.Jobs[i]
			if job.ID != m.jobID {
				continue
			}
			lines = append(lines, "Job: "+job.Name, "ID: "+job.ID, "State: "+fmt.Sprint(job.State), "Trigger: "+job.Trigger, "Duration: "+jobDuration(job), "")
			if len(job.Pages) > 0 {
				lines = append(lines, fmt.Sprintf("Affected pages: %d", len(job.Pages)))
				for _, path := range job.Pages {
					lines = append(lines, "  "+path)
				}
				lines = append(lines, "")
			}
			for _, step := range job.Steps {
				lines = append(lines, fmt.Sprintf("  %s %-12s %s", m.symbol(step.State), step.State, step.Name))
			}
			lines = append(lines, "", fmt.Sprintf("%d logs · %d diagnostics", len(job.Logs), len(job.Diagnostics)))
			for i := range job.Diagnostics {
				d := job.Diagnostics[i]
				lines = append(lines, fmt.Sprintf("%s %s %s", d.Severity, d.Code, d.Message))
			}
			for _, log := range job.Logs {
				lines = append(lines, formatLog(log))
			}
			break
		}
	case screenLogs:
		logs := m.snapshot.Logs
		if m.logJobID != "" {
			for i := range m.snapshot.Jobs {
				job := m.snapshot.Jobs[i]
				if job.ID == m.logJobID {
					logs = job.Logs
					break
				}
			}
		}
		for _, log := range logs {
			if m.match(log.Message) {
				lines = append(lines, formatLog(log))
			}
		}
	case screenWarnings, screenErrors, screenPage:
		if m.view == screenPage {
			lines = append(lines, "Page: "+m.page, "")
			for _, page := range m.snapshot.Pages {
				if page.Path == m.page {
					lines = append(lines, "Status: "+string(page.Status), "Last job: "+page.LastJobID)
					if page.URL != "" {
						lines = append(lines, "Preview: "+page.URL)
					}
					lines = append(lines, "")
					break
				}
			}
		}
		diagnostics := m.filteredDiagnostics()
		for i := range diagnostics {
			d := diagnostics[i]
			mark := "  "
			if i == m.selected {
				mark = "> "
			}
			location := d.File
			if d.Line > 0 {
				location = fmt.Sprintf("%s:%d", location, d.Line)
			}
			lines = append(lines, fmt.Sprintf("%s%s %s %s", mark, d.Code, d.Message, location))
			if i == m.selected {
				if d.Explanation != "" {
					lines = append(lines, "    "+d.Explanation)
				}
				if d.SuggestedFix != "" {
					lines = append(lines, "    Fix: "+d.SuggestedFix)
				}
			}
		}
	case screenPages:
		for i, page := range m.pages() {
			mark := "  "
			if i == m.selected {
				mark = "> "
			}
			count := 0
			status := ""
			for _, record := range m.snapshot.Pages {
				if record.Path == page {
					status = string(record.Status)
					count = len(record.Diagnostics)
					break
				}
			}
			lines = append(lines, fmt.Sprintf("%s%s · %s · %d diagnostics", mark, page, status, count))
		}
	case screenFeeds:
		for i, feed := range m.filteredFeeds() {
			mark := "  "
			if i == m.selected {
				mark = "> "
			}
			lines = append(lines, fmt.Sprintf("%s%s  %d posts", mark, feedPathOrName(feed), len(feed.Entries)))
		}
	case screenFeed:
		if feed, ok := m.selectedFeed(); ok {
			label := feed.Title
			if label == "" {
				label = feedPathOrName(feed)
			}
			lines = append(lines, "Feed: "+label, feedPathOrName(feed), fmt.Sprintf("%d posts", len(feed.Entries)))
			if feed.Status != "" {
				lines = append(lines, "Status: "+string(feed.Status))
			}
			if feed.Path != "" {
				lines = append(lines, "Source/config: "+feed.Path)
			}
			lines = append(lines, "")
			for i, entry := range m.filteredFeedEntries(feed) {
				mark := "  "
				if i == m.selected {
					mark = "> "
				}
				date := "          "
				if !entry.Date.IsZero() {
					date = entry.Date.Format("2006-01-02") + " "
				}
				title := entry.Title
				if title == "" {
					title = entry.Path
				}
				lines = append(lines, fmt.Sprintf("%s%s %s", mark, date, title))
			}
		}
	case screenHelp:
		lines = []string{"Navigation", "j/k, arrows  move or scroll", "Enter        inspect selected resource", "Esc          back", "/            search current view", "g/G          top/bottom", "PgUp/PgDn    scroll", "", "Views", "w warnings  e errors  p pages  f feeds  l logs", "b job from diagnostic or page", "", "Actions", "t trigger build  r rerun selected job", "R rebuild selected Markdown page", "q quit"}
		if m.snapshot.Server.AdminURL != "" {
			lines = append(lines, "a open Admin in browser")
		}
	}
	if len(lines) == 0 {
		return []string{"No items to show."}
	}
	return lines
}

func stateSymbol(state string) string {
	switch strings.ToLower(state) {
	case "success":
		return "✓"
	case "failed", "error":
		return "✗"
	case "warning":
		return "⚠"
	case "running":
		return "●"
	default:
		return "○"
	}
}

func (m Model) symbol(state servecontrol.State) string {
	if state == servecontrol.StateRunning {
		return pulseFrames[m.frame%len(pulseFrames)]
	}
	return stateSymbol(string(state))
}

func formatLog(log servecontrol.LogEntry) string {
	stamp := "--:--:--"
	if !log.Time.IsZero() {
		stamp = log.Time.Format(time.TimeOnly)
	}
	return fmt.Sprintf("%s %-5s %s", stamp, log.Level, log.Message)
}

func jobDuration(job servecontrol.Job) string {
	if job.StartedAt.IsZero() {
		return "not started"
	}
	end := job.EndedAt
	if end.IsZero() {
		end = time.Now()
	}
	return end.Sub(job.StartedAt).Round(10 * time.Millisecond).String()
}

// View renders from snapshot and navigation state. Rendering is bounded by the
// terminal height, even when the runtime holds many pages or log entries.
func (m Model) View() string {
	width := max(1, m.width)
	jobs := len(m.snapshot.Jobs)
	warnings, errors := 0, 0
	for i := range m.snapshot.Diagnostics {
		d := m.snapshot.Diagnostics[i]
		switch strings.ToLower(d.Severity) {
		case "warning":
			warnings++
		case "error":
			errors++
		}
	}
	header := fmt.Sprintf("markata-go serve · %s", strings.ToUpper(string(m.view)))
	if m.snapshot.Server.AdminURL != "" {
		header += " · ⚙ Admin (a)"
	}
	if m.snapshot.Server.Address != "" {
		header += " · " + m.snapshot.Server.Address
	}
	summary := fmt.Sprintf("Jobs %d  Warnings %d  Errors %d", jobs, warnings, errors)
	if m.snapshot.Server.Status != "" {
		summary = fmt.Sprintf("Server %s  ", m.snapshot.Server.Status) + summary
	}
	if m.snapshot.Site.Status != "" {
		summary = fmt.Sprintf("Site %s %s (%d pages)  ", m.symbol(m.snapshot.Site.Status), m.snapshot.Site.Status, m.snapshot.Site.PageCount) + summary
	}
	if width >= 80 {
		summary += fmt.Sprintf("  Pages %d", len(m.pages()))
	}
	if m.view == screenLogs && m.newLines > 0 {
		summary += fmt.Sprintf("  · %d new lines", m.newLines)
	}
	lines := m.bodyLines()
	start := min(m.offset, max(0, len(lines)-1))
	end := min(len(lines), start+m.bodyHeight())
	body := make([]string, 0, m.bodyHeight())
	for _, line := range lines[start:end] {
		plain := truncate(line, width)
		if m.view == screenHelp {
			plain = centerLine(plain, width)
		}
		body = append(body, m.styleLine(plain))
	}
	for len(body) < m.bodyHeight() {
		body = append(body, "")
	}
	footer := m.footer(width)
	if m.view == screenLogs && m.newLines > 0 {
		footer = fmt.Sprintf("%d new lines · G live tail · %s", m.newLines, footer)
	}
	if m.search {
		footer = "Search: " + m.query + "▏  Enter accept · Esc cancel"
	} else if m.query != "" {
		footer = "Filter: " + m.query + " · / edit · Esc clear"
	}
	if m.message != "" {
		footer = m.message + " · " + footer
	}
	border := strings.Repeat("─", width)
	return strings.Join([]string{
		lipgloss.NewStyle().Bold(true).Foreground(m.theme.accent).Render(truncate(header, width)),
		m.styleSummary(truncate(summary, width)), lipgloss.NewStyle().Foreground(m.theme.muted).Render(border),
		strings.Join(body, "\n"), border,
		lipgloss.NewStyle().Foreground(m.theme.muted).Render(truncate(footer, width)),
	}, "\n")
}

func (m Model) footer(width int) string {
	var footer string
	switch m.view {
	case screenJobs:
		footer = "j/k move  enter job  t build  r rerun  w warnings  e errors  p pages  f feeds  l logs  / search"
	case screenJob:
		footer = "j/k scroll  p page  l logs  r rerun  Esc back"
	case screenPages:
		footer = "j/k move  enter inspect  o source  v preview  R rebuild  / search"
	case screenPage:
		footer = "j/k move  o source  v preview  b job  R rebuild  Esc back"
	case screenFeeds:
		footer = "j/k move  enter feed  / search  Esc back"
	case screenFeed:
		footer = "j/k move  enter page  / search entries  Esc back"
	case screenWarnings, screenErrors:
		footer = "j/k move  enter page  o source  b job  / search  Esc back"
	case screenLogs:
		footer = "j/k scroll  G follow  g oldest  / search  PgUp/PgDn"
	default:
		footer = "j/k move  enter inspect  Esc back"
	}
	if width < 80 {
		if m.snapshot.Server.AdminURL != "" {
			return "j/k move  enter inspect  a admin  / search  ? help  q quit"
		}
		return "j/k move  enter inspect  / search  ? help  q quit"
	}
	if m.snapshot.Server.AdminURL != "" {
		footer += "  a admin"
	}
	return footer + "  ? help  q quit"
}

func (m Model) styleSummary(line string) string {
	color := m.theme.success
	//nolint:gocritic // Failure takes precedence over warning and running status.
	if m.snapshot.Site.Status == servecontrol.StateFailed || m.snapshot.Server.Status == servecontrol.StateFailed {
		color = m.theme.failure
	} else if m.snapshot.Site.Status == servecontrol.StateWarning {
		color = m.theme.warning
	} else if m.snapshot.Site.Status == servecontrol.StateRunning {
		color = m.theme.accent
	}
	return lipgloss.NewStyle().Foreground(color).Render(line)
}

func (m Model) styleLine(line string) string {
	if line == "" {
		return line
	}
	style := lipgloss.NewStyle()
	switch {
	case strings.HasPrefix(line, "> "):
		style = style.Bold(true).Foreground(m.theme.accent).Background(m.theme.focus)
	case m.view == screenHelp && (line == centerLine("Navigation", m.width) || line == centerLine("Views", m.width) || line == centerLine("Actions", m.width)):
		style = style.Bold(true).Foreground(m.theme.accent)
	case strings.Contains(line, "✗") || m.view == screenErrors || strings.Contains(line, " error "):
		style = style.Foreground(m.theme.failure)
	case strings.Contains(line, "⚠") || m.view == screenWarnings || strings.Contains(line, " warning "):
		style = style.Foreground(m.theme.warning)
	case strings.Contains(line, "✓") || strings.Contains(line, " success "):
		style = style.Foreground(m.theme.success)
	case strings.Contains(line, "●") || strings.Contains(line, "•") || strings.Contains(line, "·"):
		style = style.Foreground(m.theme.accent)
	default:
		style = style.Foreground(m.theme.text)
	}
	return style.Render(line)
}

func centerLine(line string, width int) string {
	if width < 40 || line == "" {
		return line
	}
	padding := max(0, (width-lipgloss.Width(line))/2)
	return strings.Repeat(" ", padding) + line
}

func truncate(s string, width int) string {
	if lipgloss.Width(s) <= width {
		return s
	}
	if width <= 1 {
		return "…"
	}
	var out strings.Builder
	cells := 0
	for _, r := range s {
		runeWidth := lipgloss.Width(string(r))
		if cells+runeWidth > width-1 {
			break
		}
		out.WriteRune(r)
		cells += runeWidth
	}
	return out.String() + "…"
}

// Run starts the terminal client and stops its runtime subscription on exit.
func Run(ctx context.Context, runtime interface {
	Snapshot() servecontrol.Snapshot
	Subscribe() (<-chan servecontrol.Snapshot, func())
	Trigger(servecontrol.ActionRequest) error
}, palette *palettes.Palette, paletteUpdates ...<-chan *palettes.Palette) error {
	return RunWithSourceRoot(ctx, runtime, palette, "", paletteUpdates...)
}

// RunWithSourceRoot starts the terminal client with source paths resolved from
// the configured content directory.
func RunWithSourceRoot(ctx context.Context, runtime interface {
	Snapshot() servecontrol.Snapshot
	Subscribe() (<-chan servecontrol.Snapshot, func())
	Trigger(servecontrol.ActionRequest) error
}, palette *palettes.Palette, sourceRoot string, paletteUpdates ...<-chan *palettes.Palette) error {
	updates, unsubscribe := runtime.Subscribe()
	defer unsubscribe()
	model := NewModelWithPalette(runtime.Snapshot(), runtime.Trigger, palette)
	model.sourceRoot = sourceRoot
	program := tea.NewProgram(model, tea.WithContext(ctx), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if len(paletteUpdates) > 0 && paletteUpdates[0] != nil {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case palette, ok := <-paletteUpdates[0]:
					if !ok {
						return
					}
					program.Send(paletteMsg{palette: palette})
				}
			}
		}()
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case snapshot, ok := <-updates:
				if !ok {
					return
				}
				program.Send(snapshotMsg{snapshot: snapshot})
			}
		}
	}()
	_, err := program.Run()
	return err
}
