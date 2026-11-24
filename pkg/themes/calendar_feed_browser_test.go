package themes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/chromedp/chromedp"
)

func TestCalendarYearPagination_Browser(t *testing.T) {
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
	js, err := ReadStatic("js/calendar-feed.js")
	if err != nil {
		t.Fatal(err)
	}
	fixture := `<!doctype html><html><body>
<div data-calendar-feed data-calendar-url-state data-calendar-default="list">
<div data-calendar-switch hidden><button data-calendar-mode="list">List</button><button data-calendar-mode="calendar">Calendar</button></div>
<div data-calendar-years hidden></div><ul data-calendar-list data-calendar-primary>
<li data-calendar-post data-date="2026-01-01"><a href="/new/">New</a></li>
<li data-calendar-post data-date="2024-02-29"><a href="/leap/">Leap</a></li>
<li data-calendar-post data-date="2023-01-01"><a href="/older/">Older</a></li>
</ul></div><script src="/calendar.js"></script></body></html>`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/calendar.js" {
			w.Header().Set("Content-Type", "text/javascript")
			if _, writeErr := w.Write(js); writeErr != nil {
				t.Errorf("write script: %v", writeErr)
			}
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if _, writeErr := fmt.Fprint(w, fixture); writeErr != nil {
			t.Errorf("write fixture: %v", writeErr)
		}
	}))
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
	var months int
	var year, query string
	if err := chromedp.Do(browser,
		chromedp.Navigate(server.URL),
		chromedp.WaitVisible(`[data-calendar-switch]`),
		evaluateInto(`document.querySelectorAll('.calendar-month').length`, &months),
	); err != nil {
		t.Fatal(err)
	}
	if months != 0 {
		t.Fatalf("list view created %d calendar months, want zero", months)
	}
	if err := chromedp.Do(browser,
		chromedp.Click(`[data-calendar-mode="calendar"]`),
		chromedp.WaitVisible(`#calendar-year-2026`),
		evaluateInto(`document.querySelectorAll('.calendar-month').length`, &months),
		chromedp.Click(`.calendar-year-older`),
		chromedp.WaitVisible(`#calendar-year-2024`),
		evaluateInto(`location.search`, &query),
	); err != nil {
		t.Fatal(err)
	}
	if months != 12 || !strings.Contains(query, "year=2024") {
		t.Fatalf("annual paging: months=%d query=%q, want 12 months and populated 2024", months, query)
	}
	if err := chromedp.Do(browser,
		evaluateVoid(`history.back()`),
		chromedp.WaitVisible(`#calendar-year-2026`),
		evaluateInto(`document.querySelector('.calendar-year-select').value`, &year),
	); err != nil {
		t.Fatal(err)
	}
	if year != "2026" {
		t.Fatalf("Back restored year %q, want 2026", year)
	}
	if err := chromedp.Do(browser,
		evaluateVoid(`history.forward()`),
		chromedp.WaitVisible(`#calendar-year-2024`),
		evaluateInto(`document.querySelector('.calendar-year-select').value`, &year),
	); err != nil {
		t.Fatal(err)
	}
	if year != "2024" {
		t.Fatalf("Forward restored year %q, want 2024", year)
	}
	if err := chromedp.Do(browser,
		chromedp.Navigate(server.URL+"?view=calendar&year=9999&keep=test"),
		chromedp.WaitVisible(`#calendar-year-2026`),
		chromedp.Click(`.calendar-year-older`),
		chromedp.WaitVisible(`#calendar-year-2024`),
		evaluateInto(`location.search`, &query),
		evaluateInto(`document.querySelectorAll('.calendar-month').length`, &months),
	); err != nil {
		t.Fatal(err)
	}
	if months != 12 || !strings.Contains(query, "keep=test") || !strings.Contains(query, "year=2024") {
		t.Fatalf("invalid year fallback/preserved URL: months=%d query=%q", months, query)
	}
}
