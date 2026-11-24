package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/WaylonWalker/markata/pkg/diagnostics"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/logging"
	"github.com/WaylonWalker/markata/pkg/palettes"
	"github.com/WaylonWalker/markata/pkg/servecontrol"
	"github.com/WaylonWalker/markata/pkg/servefix"
)

var serveControl struct {
	sync.RWMutex
	runtime        *servecontrol.Runtime
	jobID          string
	tui            bool
	paletteUpdates chan *palettes.Palette
}

func setServeControl(runtime *servecontrol.Runtime, tui bool, paletteUpdates chan *palettes.Palette) {
	serveControl.Lock()
	serveControl.runtime = runtime
	serveControl.jobID = ""
	serveControl.tui = tui
	serveControl.paletteUpdates = paletteUpdates
	serveControl.Unlock()
}

func clearServeControl() {
	serveControl.Lock()
	serveControl.runtime = nil
	serveControl.jobID = ""
	serveControl.tui = false
	serveControl.paletteUpdates = nil
	serveControl.Unlock()
}

func servePaletteForManager(m *lifecycle.Manager) *palettes.Palette {
	if cfg := getModelsConfig(m); cfg != nil {
		if palette, ok := loadLoggerPalette(cfg.Theme, m.Config().Extra); ok {
			return palette
		}
	}
	return nil
}

func publishServePalette(m *lifecycle.Manager) *palettes.Palette {
	palette := servePaletteForManager(m)
	if palette == nil {
		return nil
	}
	colors := servecontrol.BrowserTheme(palette.Resolve, palette.Variant != palettes.VariantLight)
	serveControl.RLock()
	runtime, updates := serveControl.runtime, serveControl.paletteUpdates
	serveControl.RUnlock()
	if runtime != nil {
		runtime.SetTheme(colors)
	}
	if updates != nil {
		select {
		case updates <- palette:
		default:
			select {
			case <-updates:
			default:
			}
			updates <- palette
		}
	}
	return palette
}

func currentServeControl() (*servecontrol.Runtime, string, bool) {
	serveControl.RLock()
	defer serveControl.RUnlock()
	return serveControl.runtime, serveControl.jobID, serveControl.tui
}

func startServeJob(name, trigger string, pages []string) string {
	runtime, _, _ := currentServeControl()
	if runtime == nil {
		return ""
	}
	runtime.SetSite(servecontrol.SiteState{Status: servecontrol.StateRunning, PageCount: len(runtime.Snapshot().Pages)})
	id := runtime.QueueJob(servecontrol.JobSpec{Name: name, Type: "build", Trigger: trigger, Pages: pages})
	runtime.StartJob(id)
	serveControl.Lock()
	serveControl.jobID = id
	serveControl.Unlock()
	return id
}

//nolint:gocyclo // This is the single reconciliation point for lifecycle output and diagnostics.
func finishServeJob(id string, m *lifecycle.Manager, result *BuildResult, buildErr error) {
	runtime, _, _ := currentServeControl()
	if runtime == nil || id == "" {
		return
	}
	var snapshot diagnostics.ContentLedgerSnapshot
	if result != nil {
		snapshot = result.Content
	} else if m != nil {
		snapshot = m.ContentDiagnostics()
	}
	if m != nil {
		feeds := make([]servecontrol.Feed, 0, len(m.Feeds()))
		for _, feed := range m.Feeds() {
			if feed == nil {
				continue
			}
			projected := servecontrol.Feed{Name: feed.Name, Title: feed.Title, Path: feed.Path, Status: servecontrol.StateSuccess, Entries: make([]servecontrol.FeedEntry, 0, len(feed.Posts))}
			for _, post := range feed.Posts {
				if post == nil {
					continue
				}
				entry := servecontrol.FeedEntry{Path: filepath.ToSlash(post.Path), URL: post.Href}
				if post.Title != nil {
					entry.Title = *post.Title
				}
				if post.Date != nil {
					entry.Date = *post.Date
				}
				projected.Entries = append(projected.Entries, entry)
			}
			feeds = append(feeds, projected)
		}
		runtime.SetFeeds(feeds)
	}
	pages := make([]servecontrol.Page, 0, len(snapshot.Entries))
	urls := make(map[string]string)
	if m != nil {
		for _, post := range m.Posts() {
			urls[filepath.ToSlash(post.Path)] = post.Href
		}
	}
	for _, entry := range snapshot.Entries {
		state := servecontrol.StateSuccess
		if !entry.Emitted {
			state = servecontrol.StateWarning
		}
		pages = append(pages, servecontrol.Page{Path: entry.Path, URL: urls[filepath.ToSlash(entry.Path)], Status: state, JobID: id, LastJobID: id})
	}
	affected := make([]string, 0, len(pages))
	for _, page := range pages {
		affected = append(affected, page.Path)
	}
	runtime.SetJobPages(id, affected)
	if len(pages) > 0 || buildErr == nil {
		runtime.SetPages(pages)
	}
	seen := make(map[string]struct{})
	fixPlans := make(map[string]servefix.Plan)
	add := func(issue diagnostics.Issue) {
		key := fmt.Sprintf("%s:%d:%d:%d:%s:%s", issue.File, issue.Range.StartLine, issue.Range.StartCol, issue.Range.EndCol, issue.Code, issue.Message)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		diagnostic := servecontrol.Diagnostic{
			Code: issue.Code, Severity: issue.Severity.String(), Message: issue.Message,
			File: issue.File, Page: issue.File, Line: issue.Range.StartLine + 1,
			Column: issue.Range.StartCol + 1, JobID: id,
			Explanation: serveSuggestedFix(issue.Code), SuggestedFix: serveSuggestedFix(issue.Code),
			FixSafety: string(servefix.SafetyManual),
		}
		if m != nil && issue.File != "" {
			cacheKey := filepath.Clean(issue.File)
			plan, ok := fixPlans[cacheKey]
			if !ok {
				var planErr error
				plan, planErr = servefix.PlanFile(m.Config().ContentDir, issue.File)
				if planErr != nil {
					plan = servefix.Plan{}
				}
				fixPlans[cacheKey] = plan
			}
			for i := range plan.Edits {
				edit := plan.Edits[i]
				if edit.Code != issue.Code || edit.Line != issue.Range.StartLine+1 || edit.Column != issue.Range.StartCol+1 {
					continue
				}
				file := filepath.ToSlash(issue.File)
				if filepath.IsAbs(file) {
					contentRoot, rootErr := filepath.Abs(m.Config().ContentDir)
					if relative, err := filepath.Rel(contentRoot, issue.File); rootErr == nil && err == nil {
						file = filepath.ToSlash(relative)
					}
				}
				diagnostic.FixSafety = string(edit.Safety)
				diagnostic.FixPlans = append(diagnostic.FixPlans, servecontrol.FixPlan{
					ID: edit.ID, Category: edit.Category, Safety: string(edit.Safety), File: file,
					Line: edit.Line, Column: edit.Column, EndLine: edit.EndLine,
					EndColumn: edit.EndColumn, Before: edit.Before, After: edit.After,
					Digest: plan.Digest, Explanation: edit.Explanation,
				})
			}
		}
		runtime.AddDiagnostic(diagnostic)
	}
	for _, entry := range snapshot.Entries {
		for _, issue := range entry.Diagnostics {
			add(issue)
		}
		if !strings.HasSuffix(entry.Path, ".md") && !strings.HasSuffix(entry.Path, ".markdown") {
			continue
		}
		content, err := os.ReadFile(entry.Path)
		if err != nil && m != nil {
			content, err = os.ReadFile(filepath.Join(m.Config().ContentDir, entry.Path))
		}
		if err == nil {
			for _, issue := range diagnostics.Check(entry.Path, string(content), nil) {
				add(issue)
			}
		}
	}
	if result != nil {
		for _, warning := range result.Warnings {
			runtime.AddDiagnostic(servecontrol.Diagnostic{Code: "serve.plugin_warning", Severity: "warning", Message: warning, JobID: id, SuggestedFix: "Inspect the job steps and source files associated with this warning."})
		}
	}
	state := servecontrol.StateSuccess
	if buildErr != nil {
		state = servecontrol.StateFailed
		runtime.AddDiagnostic(servecontrol.Diagnostic{Code: "serve.build_failed", Severity: "error", Message: buildErr.Error(), JobID: id, SuggestedFix: "Inspect the failed step and correct its source or configuration, then trigger a build."})
	} else {
		snapshot := runtime.Snapshot()
		for i := range snapshot.Jobs {
			job := snapshot.Jobs[i]
			if job.ID != id {
				continue
			}
			for i := range job.Diagnostics {
				diagnostic := job.Diagnostics[i]
				if diagnostic.Severity == buildStatusError {
					state = servecontrol.StateFailed
					break
				}
				if diagnostic.Severity == buildStatusWarning && state == servecontrol.StateSuccess {
					state = servecontrol.StateWarning
				}
			}
			break
		}
	}
	runtime.FinishJob(id, state)
	message := ""
	if buildErr != nil {
		message = buildErr.Error()
	}
	runtime.SetSite(servecontrol.SiteState{Status: state, Message: message, PageCount: len(runtime.Snapshot().Pages)})
	serveControl.Lock()
	if serveControl.jobID == id {
		serveControl.jobID = ""
	}
	serveControl.Unlock()
}

func runServeBuild(m *lifecycle.Manager, id string) (*BuildResult, error) {
	runtime, _, _ := currentServeControl()
	if runtime == nil || id == "" {
		return runBuild(m)
	}
	steps := make(map[lifecycle.Stage]string)
	return runBuildObserved(m, func(stage lifecycle.Stage, starting bool, stageErr error) {
		if starting {
			steps[stage] = runtime.StartStep(id, string(stage))
			runtime.AddLog(servecontrol.LogEntry{Level: "info", Message: string(stage) + " started", JobID: id, StepID: steps[stage]})
			return
		}
		state := servecontrol.StateSuccess
		if stageErr != nil {
			state = servecontrol.StateFailed
			runtime.AddDiagnostic(servecontrol.Diagnostic{
				Code: "serve.stage_failed", Severity: "error", Message: stageErr.Error(),
				JobID: id, StepID: steps[stage],
				SuggestedFix: "Inspect this stage's log and correct the reported source or configuration.",
			})
		}
		runtime.FinishStep(id, steps[stage], state)
	})
}

func serveSuggestedFix(code string) string {
	switch code {
	case "missing-alt-text":
		return "Describe the image in its Markdown alt text, for example ![A short description](image.jpg)."
	case "h1-in-content":
		return "Use a level-two heading in the body; the page title supplies the level-one heading."
	case "protocol-less-url":
		return "Add https:// to the link URL."
	case "frontmatter.missing_closing_delimiter", "frontmatter.malformed_closing_delimiter":
		return "Close the frontmatter with a line containing exactly three hyphens (---)."
	case "frontmatter.invalid_date":
		return "Use an ISO date such as 2026-09-26."
	case "frontmatter.duplicate_key", "duplicate-key":
		return "Keep one value for this frontmatter key."
	default:
		return "Open the source location and correct the reported value."
	}
}

func configureServeLogger() {
	runtime, _, _ := currentServeControl()
	if runtime == nil {
		return
	}
	configureServeLoggerWithOptions(currentLogTheme, logging.FormatPlain)
}

func configureServeLoggerWithOptions(theme logging.Theme, format logging.Format) {
	runtime, _, tui := currentServeControl()
	if runtime == nil {
		return
	}
	writer := errWriter()
	if tui {
		writer = io.Discard
	}
	logging.ConfigureStandardLogger(logging.Options{
		Writer: writer, Format: format, ForceColor: forceColor, NoColor: noColor || tui,
		IsTTY: errorOutputIsTerminal(), Theme: theme,
		Observer: func(entry logging.Entry, message string) {
			_, id, _ := currentServeControl()
			level := entry.Level
			if level == "" {
				lower := strings.ToLower(strings.TrimSpace(message))
				switch {
				case strings.HasPrefix(lower, "warning:"), strings.Contains(lower, " warning:"):
					level = "warning"
				case strings.HasPrefix(lower, "error:"), strings.Contains(lower, " error:"):
					level = "error"
				default:
					level = "info"
				}
			}
			runtime.AddLog(servecontrol.LogEntry{Level: level, Message: message, JobID: id})
			if level == "warning" || level == "error" {
				runtime.AddDiagnostic(servecontrol.Diagnostic{Code: "serve.log_" + level, Severity: level, Message: message, JobID: id, SuggestedFix: "Inspect the related job and source files for the cause."})
			}
		},
	})
}
