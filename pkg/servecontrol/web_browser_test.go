package servecontrol

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"

	"github.com/WaylonWalker/markata/pkg/servefix"
)

func TestBrowserKeyboardNavigationAndHistory(t *testing.T) {
	browserPath := ""
	for _, candidate := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if path, err := exec.LookPath(candidate); err == nil {
			browserPath = path
			break
		}
	}
	if browserPath == "" {
		t.Skip("Chromium is unavailable")
	}

	runtime := NewRuntime()
	firstJob := runtime.QueueJob(JobSpec{Name: "First build", Type: "build"})
	runtime.StartJob(firstJob)
	runtime.FinishJob(firstJob, StateSuccess)
	secondJob := runtime.QueueJob(JobSpec{Name: "Second build", Type: "build"})
	runtime.StartJob(secondJob)
	runtime.FinishJob(secondJob, StateFailed)
	runtime.SetPages([]Page{{Path: "posts/first.md"}, {Path: "posts/second.md"}})
	actions := make(chan ActionRequest, 1)
	runtime.SetActionHandler(func(request ActionRequest) error {
		actions <- request
		return nil
	})
	server := httptest.NewServer(NewWebHandler(runtime))
	defer server.Close()

	allocator, cancelAllocator := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WSURLReadTimeout(60*time.Second), chromedp.ExecPath(browserPath), chromedp.Headless, chromedp.NoFirstRun,
		chromedp.Flag("no-sandbox", true), chromedp.Flag("disable-dev-shm-usage", true),
	)...)
	defer cancelAllocator()
	browser, cancelBrowser := chromedp.NewContext(allocator)
	defer cancelBrowser()
	browser, cancelTimeout := context.WithTimeout(browser, 120*time.Second)
	defer cancelTimeout()

	var hash, section, filter string
	if err := chromedp.Do(browser,
		chromedp.EmulateViewport(1280, 900),
		chromedp.Navigate(server.URL+"/_markata/"),
		chromedp.WaitVisible(`#list [data-key]`),
		chromedp.KeyEvent("p"),
		evaluateInto(`location.hash`, &hash),
	); err != nil {
		t.Fatal(err)
	}
	if hash != "#/pages?item=page%3Aposts%2Ffirst.md" {
		t.Fatalf("p shortcut hash = %q", hash)
	}
	if err := chromedp.Do(browser,
		chromedp.KeyEvent("j"),
		evaluateInto(`location.hash`, &hash),
	); err != nil {
		t.Fatal(err)
	}
	if hash != "#/pages?item=page%3Aposts%2Fsecond.md" {
		t.Fatalf("j selection hash = %q", hash)
	}
	if err := chromedp.Do(browser,
		chromedp.KeyEvent(kb.Escape),
		chromedp.Poll[chromedp.Void](`location.hash === "#/pages?item=page%3Aposts%2Ffirst.md"`, chromedp.WithPollingTimeout(5*time.Second)),
		evaluateInto(`location.hash`, &hash),
		evaluateInto(`document.querySelector('#title').textContent`, &section),
	); err != nil {
		t.Fatal(err)
	}
	if hash != "#/pages?item=page%3Aposts%2Ffirst.md" || section != "Pages" {
		t.Fatalf("Escape did not return to the prior selection: hash=%q title=%q", hash, section)
	}
	if err := chromedp.Do(browser,
		evaluateVoid(`history.forward()`),
		chromedp.Poll[chromedp.Void](`location.hash === "#/pages?item=page%3Aposts%2Fsecond.md"`, chromedp.WithPollingTimeout(5*time.Second)),
		evaluateInto(`location.hash`, &hash),
	); err != nil {
		t.Fatal(err)
	}
	if hash != "#/pages?item=page%3Aposts%2Fsecond.md" {
		t.Fatalf("browser Forward did not restore selected item: %q", hash)
	}
	if err := chromedp.Do(browser,
		chromedp.KeyEvent(kb.Escape),
		chromedp.Poll[chromedp.Void](`location.hash === "#/pages?item=page%3Aposts%2Ffirst.md"`, chromedp.WithPollingTimeout(5*time.Second)),
		evaluateInto(`location.hash`, &hash),
	); err != nil {
		t.Fatal(err)
	}
	if hash != "#/pages?item=page%3Aposts%2Ffirst.md" {
		t.Fatalf("second Escape did not return to prior route: %q", hash)
	}
	if err := chromedp.Do(browser,
		chromedp.Focus(`#filter`),
		chromedp.KeyEvent("p"),
		evaluateInto(`document.querySelector('#filter').value`, &filter),
		evaluateInto(`document.querySelector('#title').textContent`, &section),
	); err != nil {
		t.Fatal(err)
	}
	if filter != "p" || section != "Pages" {
		t.Fatalf("typing in filter was intercepted: filter=%q section=%q", filter, section)
	}
	if err := chromedp.Do(browser,
		chromedp.Focus(`#desktop-nav [data-section="feeds"]`),
		chromedp.KeyEvent(kb.Enter),
		evaluateInto(`document.querySelector('#title').textContent`, &section),
	); err != nil {
		t.Fatal(err)
	}
	if section != "Feeds" {
		t.Fatalf("Enter did not activate focused nav button: title=%q", section)
	}
	if err := chromedp.Do(browser,
		chromedp.Click(`#desktop-nav [data-section="jobs"]`),
		chromedp.KeyEvent("j"),
		chromedp.KeyEvent("r"),
	); err != nil {
		t.Fatal(err)
	}
	select {
	case action := <-actions:
		if action.Kind != "rerun" || action.JobID != secondJob {
			t.Fatalf("r shortcut action = %+v, want rerun for %s", action, secondJob)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("r shortcut did not rerun the selected completed job")
	}
	if err := chromedp.Do(browser, evaluateInto(`location.hash`, &hash)); err != nil {
		t.Fatal(err)
	}
	if hash != "#/jobs?item=job%3Ajob-2" {
		t.Fatalf("section navigation after browser interaction = %q", hash)
	}
	var focused string
	if err := chromedp.Do(browser,
		evaluateVoid(`document.querySelector('#list [data-key="job:job-2"]').focus()`),
		chromedp.Sleep(2300*time.Millisecond),
		evaluateInto(`document.activeElement?.dataset?.key || ''`, &focused),
	); err != nil {
		t.Fatal(err)
	}
	if focused != "job:job-2" {
		t.Fatalf("poll replaced focused row: focus = %q", focused)
	}
	runtime.QueueJob(JobSpec{Name: "Third build", Type: "build"})
	if err := chromedp.Do(browser,
		chromedp.Sleep(2300*time.Millisecond),
		evaluateInto(`document.activeElement?.dataset?.key || ''`, &focused),
	); err != nil {
		t.Fatal(err)
	}
	if focused != "job:job-2" {
		t.Fatalf("state update lost focused row: focus = %q", focused)
	}
}

func TestBrowserGroupedSafeFixReviewAndApply(t *testing.T) {
	browserPath := ""
	for _, candidate := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if path, err := exec.LookPath(candidate); err == nil {
			browserPath = path
			break
		}
	}
	if browserPath == "" {
		t.Skip("Chromium is unavailable")
	}
	root := t.TempDir()
	runtime := NewRuntime()
	for i, name := range []string{"first.md", "second.md"} {
		content := "---\ntitle: Post\npublished: true\n---\n\n[link](//example.com)\n"
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		plan, err := servefix.PlanFile(root, name)
		if err != nil || len(plan.Edits) == 0 {
			t.Fatalf("plan %s: %v, edits=%d", name, err, len(plan.Edits))
		}
		edit := plan.Edits[0]
		runtime.AddDiagnostic(Diagnostic{Code: edit.Code, Severity: "warning", Message: edit.Message, File: name, Page: name, Line: edit.Line, JobID: fmt.Sprintf("job-%d", i), FixSafety: string(edit.Safety), FixPlans: []FixPlan{{ID: edit.ID, File: name, Category: edit.Category, Safety: string(edit.Safety), Digest: plan.Digest, Before: edit.Before, After: edit.After}}})
	}
	actions := make(chan ActionRequest, 2)
	runtime.SetActionHandler(func(action ActionRequest) error {
		actions <- action
		if action.Kind == "build" {
			job := runtime.QueueJob(JobSpec{Name: "Fix rebuild", Type: "build"})
			runtime.StartJob(job)
			runtime.SetPages(nil)
			runtime.FinishJob(job, StateSuccess)
		}
		return nil
	})
	server := httptest.NewServer(NewWebHandlerWithSourceRoot(runtime, root))
	defer server.Close()
	allocator, cancelAllocator := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:], chromedp.WSURLReadTimeout(60*time.Second), chromedp.ExecPath(browserPath), chromedp.Headless, chromedp.NoFirstRun, chromedp.Flag("no-sandbox", true), chromedp.Flag("disable-dev-shm-usage", true))...)
	defer cancelAllocator()
	browser, cancelBrowser := chromedp.NewContext(allocator)
	defer cancelBrowser()
	browser, cancelTimeout := context.WithTimeout(browser, 120*time.Second)
	defer cancelTimeout()
	var reviewOpen bool
	var richDiffRendered bool
	if err := chromedp.Do(browser,
		chromedp.Navigate(server.URL+"/_markata/#/diagnostics"),
		chromedp.WaitVisible(`#fix-tools [data-fix-group="all"]`),
		chromedp.Click(`#fix-tools [data-fix-group="all"]`),
		chromedp.WaitVisible(`#fix-review[open]`),
		chromedp.WaitVisible(`#fix-approve:not([disabled])`),
		evaluateInto(`document.querySelector('#fix-review-body').textContent.includes('2 edits across 2 files')`, &reviewOpen),
		pollInto(`[...document.querySelectorAll('.diff-host')].length === 2 && [...document.querySelectorAll('.diff-host')].every(host => host.shadowRoot?.querySelector('[data-diffs-header]') && host.shadowRoot.textContent.includes('https://example.com'))`, &richDiffRendered, chromedp.WithPollingTimeout(10*time.Second)),
	); err != nil {
		t.Fatal(err)
	}
	if !reviewOpen {
		t.Fatal("grouped preview did not show both edits")
	}
	if !richDiffRendered {
		t.Fatal("Pierre did not render both frozen source diffs")
	}
	if err := chromedp.Do(browser, chromedp.Click(`#fix-approve`)); err != nil {
		t.Fatal(err)
	}
	select {
	case action := <-actions:
		if action.Kind != "build" {
			t.Fatalf("batch action = %+v, want one build", action)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("grouped apply did not request a rebuild")
	}
	select {
	case action := <-actions:
		t.Fatalf("grouped apply requested another rebuild: %+v", action)
	case <-time.After(200 * time.Millisecond):
	}
	for _, name := range []string{"first.md", "second.md"} {
		content, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || !strings.Contains(string(content), "https://example.com") {
			t.Fatalf("%s was not repaired: %v, %q", name, err, content)
		}
	}
	var noticeText string
	if err := chromedp.Do(browser,
		chromedp.Sleep(2300*time.Millisecond),
		evaluateInto(`document.querySelector('#notice').textContent`, &noticeText),
	); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(noticeText, "Applied 2 fixes") || !strings.Contains(noticeText, "2 diagnostics resolved") {
		t.Fatalf("batch outcome did not persist after live poll: %q", noticeText)
	}
}
