package servetui

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/servecontrol"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func key(value string) tea.KeyMsg {
	if value == "enter" {
		return tea.KeyMsg{Type: tea.KeyEnter}
	}
	if value == "esc" {
		return tea.KeyMsg{Type: tea.KeyEscape}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}

func update(t *testing.T, model Model, msg tea.Msg) Model {
	t.Helper()
	next, _ := model.Update(msg)
	result, ok := next.(Model)
	if !ok {
		t.Fatalf("Update returned %T", next)
	}
	return result
}

func sampleSnapshot() servecontrol.Snapshot {
	diagnostic := servecontrol.Diagnostic{Code: "MARKATA-W014", Severity: "warning", Message: "Unknown frontmatter key", Page: "posts/foo.md", File: "posts/foo.md", Line: 4, JobID: "job-1", SuggestedFix: "published: true"}
	return servecontrol.Snapshot{
		Server:      servecontrol.ServerState{Status: servecontrol.StateSuccess, Address: "localhost:8000"},
		Site:        servecontrol.SiteState{Status: servecontrol.StateWarning, PageCount: 2},
		Jobs:        []servecontrol.Job{{ID: "job-1", Name: "initial build", State: servecontrol.StateWarning, Steps: []servecontrol.Step{{ID: "step-1", Name: "render", State: servecontrol.StateSuccess}}, Diagnostics: []servecontrol.Diagnostic{diagnostic}, Pages: []string{"posts/foo.md"}}},
		Diagnostics: []servecontrol.Diagnostic{diagnostic},
		Pages:       []servecontrol.Page{{Path: "posts/foo.md", Status: servecontrol.StateWarning, JobID: "job-1", LastJobID: "job-1", Diagnostics: []servecontrol.Diagnostic{diagnostic}}, {Path: "posts/clean.md", Status: servecontrol.StateSuccess, LastJobID: "job-1"}},
		Logs:        []servecontrol.LogEntry{{Time: time.Date(2026, 9, 26, 12, 0, 0, 0, time.UTC), Level: "info", Message: "first log", JobID: "job-1"}},
	}
}

func TestAdminAffordanceUsesMountedURL(t *testing.T) {
	snapshot := sampleSnapshot()
	without := NewModel(snapshot, nil)
	if strings.Contains(without.View(), "Admin") {
		t.Fatal("disabled Admin should not be advertised")
	}
	snapshot.Server.AdminURL = "http://localhost:8000/_markata/"
	with := NewModel(snapshot, nil)
	if !strings.Contains(with.View(), "Admin (a)") {
		t.Fatal("mounted Admin should be discoverable")
	}
	_, command := with.Update(key("a"))
	if command == nil {
		t.Fatal("Admin shortcut should open the mounted URL")
	}
	_, command = without.Update(key("a"))
	if command != nil {
		t.Fatal("disabled Admin shortcut should do nothing")
	}
}

func TestNavigationAndDiagnosticsRemainDiscoverable(t *testing.T) {
	m := NewModel(sampleSnapshot(), nil)
	m = update(t, m, key("enter"))
	if m.view != screenJob || !strings.Contains(m.View(), "render") {
		t.Fatalf("job detail missing: %q", m.View())
	}
	m = update(t, m, key("w"))
	if m.view != screenWarnings || !strings.Contains(m.View(), "MARKATA-W014") {
		t.Fatalf("warning inbox missing: %q", m.View())
	}
	m = update(t, m, key("enter"))
	if m.view != screenPage || !strings.Contains(m.View(), "published: true") {
		t.Fatalf("page diagnostics missing: %q", m.View())
	}
	m = update(t, m, key("esc"))
	if m.view != screenWarnings {
		t.Fatalf("page back should return to warnings, got %s", m.view)
	}
	m = update(t, m, key("esc"))
	if m.view != screenJob {
		t.Fatalf("warnings back should return to job, got %s", m.view)
	}
	m = update(t, m, key("w"))
	m = update(t, m, key("enter"))
	m = update(t, m, key("b"))
	if m.view != screenJob {
		t.Fatalf("diagnostic page should open its job, got %s", m.view)
	}
	m = update(t, m, key("esc"))
	if m.view != screenPage {
		t.Fatalf("job back should return to page, got %s", m.view)
	}
	m = update(t, m, key("esc"))
	if m.view != screenWarnings {
		t.Fatalf("page back should return to warnings, got %s", m.view)
	}
	snapshot := sampleSnapshot()
	for range 200 {
		snapshot.Logs = append(snapshot.Logs, servecontrol.LogEntry{Level: "info", Message: "later log"})
	}
	m = update(t, m, snapshotMsg{snapshot})
	if !strings.Contains(m.View(), "MARKATA-W014") {
		t.Fatal("diagnostic disappeared after later logs")
	}
}

func TestLogScrollPreservesPositionAndCountsNewLines(t *testing.T) {
	snapshot := sampleSnapshot()
	for range 50 {
		snapshot.Logs = append(snapshot.Logs, servecontrol.LogEntry{Message: "old log"})
	}
	m := NewModel(snapshot, nil)
	m = update(t, m, key("l"))
	if !m.follow || m.offset == 0 {
		t.Fatal("logs should open at live tail")
	}
	m = update(t, m, key("k"))
	offset := m.offset
	if m.follow {
		t.Fatal("scrolling up should pause live tail")
	}
	snapshot.Logs = append(snapshot.Logs, servecontrol.LogEntry{Message: "new log"})
	m = update(t, m, snapshotMsg{snapshot})
	if m.offset != offset || m.newLines != 1 || !strings.Contains(m.View(), "1 new lines") {
		t.Fatalf("scroll anchor/new line count: offset=%d new=%d", m.offset, m.newLines)
	}
	m = update(t, m, key("G"))
	if !m.follow || m.newLines != 0 || m.offset != m.maxOffset() {
		t.Fatal("G should return to live tail")
	}
}

func TestSearchAndPageInventory(t *testing.T) {
	m := NewModel(sampleSnapshot(), nil)
	m = update(t, m, key("p"))
	if len(m.pages()) != 2 {
		t.Fatalf("want two pages including clean page, got %d", len(m.pages()))
	}
	m = update(t, m, key("/"))
	m = update(t, m, key("clean"))
	m = update(t, m, key("enter"))
	if len(m.pages()) != 1 || !strings.Contains(m.View(), "posts/clean.md") || strings.Contains(m.View(), "posts/foo.md") {
		t.Fatalf("page search failed: %q", m.View())
	}
}

func TestNarrowResizeBoundsOutput(t *testing.T) {
	snapshot := sampleSnapshot()
	snapshot.Jobs[0].Name = "日本語の長いビルド名"
	m := NewModel(snapshot, nil)
	m = update(t, m, tea.WindowSizeMsg{Width: 36, Height: 12})
	for _, line := range strings.Split(m.View(), "\n") {
		if lipgloss.Width(line) > 36 {
			t.Fatalf("line exceeds terminal width: %q", line)
		}
	}
	if got := len(strings.Split(m.View(), "\n")); got > 12 {
		t.Fatalf("view exceeds terminal height: %d", got)
	}
	m = update(t, m, tea.WindowSizeMsg{Width: 100, Height: 30})
	if m.width != 100 || m.height != 30 {
		t.Fatal("resize did not update dimensions")
	}
}

func TestRerunActionUsesSelectedJob(t *testing.T) {
	var got servecontrol.ActionRequest
	m := NewModel(sampleSnapshot(), func(request servecontrol.ActionRequest) error { got = request; return nil })
	_, command := m.Update(key("r"))
	if command == nil {
		t.Fatal("expected rerun command")
	}
	command()
	if got.Kind != "rerun" || got.JobID != "job-1" {
		t.Fatalf("unexpected action: %+v", got)
	}
}

func TestPageUsesCurrentDiagnosticsWhileInboxRetainsHistory(t *testing.T) {
	snapshot := sampleSnapshot()
	snapshot.Pages[0].Status = servecontrol.StateSuccess
	snapshot.Pages[0].Diagnostics = nil
	m := NewModel(snapshot, nil)
	m = update(t, m, key("p"))
	if !strings.Contains(m.View(), "posts/foo.md · success · 0 diagnostics") {
		t.Fatalf("page inventory counted a historical warning: %q", m.View())
	}
	m = update(t, m, key("j"))
	m = update(t, m, key("enter"))
	if strings.Contains(m.View(), "MARKATA-W014") {
		t.Fatal("current page shows a resolved historical warning")
	}
	m = update(t, m, key("w"))
	if !strings.Contains(m.View(), "MARKATA-W014") {
		t.Fatal("session warning history was lost")
	}
}

func TestRebuildSelectedPageAction(t *testing.T) {
	var got servecontrol.ActionRequest
	m := NewModel(sampleSnapshot(), func(request servecontrol.ActionRequest) error { got = request; return nil })
	m = update(t, m, key("p"))
	_, command := m.Update(key("R"))
	if command == nil {
		t.Fatal("expected rebuild command")
	}
	command()
	if got.Kind != "rebuild-page" || got.Page != "posts/clean.md" {
		t.Fatalf("unexpected page action: %+v", got)
	}
}

func TestThemeHelpAndRunningAnimation(t *testing.T) {
	snapshot := sampleSnapshot()
	snapshot.Site.Status = servecontrol.StateRunning
	snapshot.Jobs[0].State = servecontrol.StateRunning
	palette := palettes.NewPalette("test-site", palettes.VariantDark)
	palette.Colors["accent"] = "#123456"
	palette.Colors["warning"] = "#abcdef"
	palette.Colors["error"] = "#fedcba"
	m := NewModelWithPalette(snapshot, nil, palette)
	if m.theme.accent != lipgloss.Color("#123456") || m.theme.warning != lipgloss.Color("#abcdef") {
		t.Fatalf("site palette not applied: %+v", m.theme)
	}
	if m.Init() == nil {
		t.Fatal("running work should schedule animation")
	}
	before := m.View()
	m = update(t, m, tickMsg{})
	if before == m.View() || !strings.Contains(m.View(), pulseFrames[1]) {
		t.Fatal("running indicator did not animate")
	}
	m = update(t, m, tea.WindowSizeMsg{Width: 80, Height: 35})
	m = update(t, m, key("?"))
	view := m.View()
	if !strings.Contains(view, "Navigation") || !strings.Contains(view, "R rebuild selected Markdown page") {
		t.Fatalf("help lacks controls: %q", view)
	}
	if !strings.Contains(view, strings.Repeat(" ", 35)) || centerLine("Navigation", 80) == "Navigation" {
		t.Fatal("help heading is not centered")
	}
	second := palettes.NewPalette("other-site", palettes.VariantLight)
	second.Colors["accent"] = "#654321"
	m = update(t, m, paletteMsg{second})
	if m.theme.accent != lipgloss.Color("#654321") {
		t.Fatal("palette update was not applied")
	}
}

func TestFeedsSearchAndOpenExistingPageDetail(t *testing.T) {
	snapshot := sampleSnapshot()
	snapshot.Feeds = []servecontrol.Feed{
		{Name: "archive", Title: "Archive", Path: "/archive/", Entries: []servecontrol.FeedEntry{{Path: "posts/other.md", Title: "Other"}}},
		{Name: "thoughts", Title: "Thoughts", Path: "/thoughts/", Entries: []servecontrol.FeedEntry{{Path: "posts/cli.md", Title: "CLI serve controls"}}},
	}
	snapshot.Pages = append(snapshot.Pages, servecontrol.Page{Path: "posts/cli.md", URL: "/cli/", Status: servecontrol.StateSuccess})
	m := NewModel(snapshot, nil)
	m = update(t, m, key("f"))
	if m.view != screenFeeds || !strings.Contains(m.View(), "/thoughts/  1 posts") {
		t.Fatalf("feed inventory missing: %q", m.View())
	}
	m = update(t, m, key("/"))
	for _, r := range "thoughts" {
		m = update(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	m = update(t, m, key("enter"))
	m = update(t, m, key("enter"))
	if m.view != screenFeed || !strings.Contains(m.View(), "Thoughts") {
		t.Fatalf("feed detail did not open: %q", m.View())
	}
	m = update(t, m, key("/"))
	for _, r := range "cli" {
		m = update(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	if !strings.Contains(m.View(), "CLI serve controls") {
		t.Fatalf("feed entry search did not match: %q", m.View())
	}
	m = update(t, m, key("enter"))
	m = update(t, m, key("enter"))
	if m.view != screenPage || m.page != "posts/cli.md" {
		t.Fatalf("feed entry did not open existing page detail: view=%s page=%s", m.view, m.page)
	}
	m = update(t, m, key("esc"))
	if m.view != screenFeed {
		t.Fatalf("page back should return to feed, got %s", m.view)
	}
	m = update(t, m, key("esc"))
	if m.view != screenFeeds {
		t.Fatalf("feed back should return to feeds, got %s", m.view)
	}
}

func TestJobFooterAdvertisesLogsKey(t *testing.T) {
	m := NewModel(sampleSnapshot(), nil)
	m = update(t, m, key("enter"))
	if !strings.Contains(m.footer(m.width), "l logs") {
		t.Fatalf("job footer should advertise l for logs: %q", m.footer(m.width))
	}
	m = update(t, m, key("l"))
	if m.view != screenLogs || m.logJobID != "job-1" {
		t.Fatalf("l should open this job's logs: view=%s job=%s", m.view, m.logJobID)
	}
}

func TestSourceOpeningUsesConfiguredContentRoot(t *testing.T) {
	root := t.TempDir()
	m := NewModel(sampleSnapshot(), nil)
	m.view = screenPage
	m.page = "posts/foo.md"
	m.sourceRoot = root
	path, _ := m.selectedSource()
	if want := filepath.Join(root, "posts", "foo.md"); path != want {
		t.Fatalf("source path = %q, want %q", path, want)
	}
}

func TestPageJobTraversalAndGeneratedPageAction(t *testing.T) {
	snapshot := sampleSnapshot()
	snapshot.Pages = append(snapshot.Pages, servecontrol.Page{Path: "feeds/index.html", Status: servecontrol.StateSuccess})
	m := NewModel(snapshot, nil)
	m = update(t, m, key("enter"))
	m = update(t, m, key("p"))
	if m.view != screenPage || m.page != "posts/foo.md" {
		t.Fatalf("job to page traversal failed: %s %s", m.view, m.page)
	}
	m = update(t, m, key("b"))
	if m.view != screenJob || m.jobID != "job-1" {
		t.Fatal("page to last job traversal failed")
	}
	m = update(t, m, key("p"))
	m = update(t, m, key("p"))
	if m.view != screenPages {
		t.Fatal("expected page inventory")
	}
	_, command := m.Update(key("R"))
	if command != nil {
		t.Fatal("generated HTML page should not offer source rebuild")
	}
}

func TestSiteFailureRemainsVisible(t *testing.T) {
	snapshot := sampleSnapshot()
	snapshot.Site = servecontrol.SiteState{Status: servecontrol.StateFailed, Message: "posts/foo.md:17: invalid date", PageCount: 2}
	m := NewModel(snapshot, nil)
	if !strings.Contains(m.View(), "Build failed: posts/foo.md:17: invalid date") {
		t.Fatalf("failed current site not pinned: %q", m.View())
	}
}
