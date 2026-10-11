package themes_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WaylonWalker/markata-go/pkg/models"
	"github.com/WaylonWalker/markata-go/pkg/templates"
	"github.com/WaylonWalker/markata-go/pkg/themes"

	"github.com/chromedp/chromedp"
)

func TestEditorialHistory_LazySafeComparisonBrowser(t *testing.T) {
	browserPath := ""
	for _, candidate := range []string{"chromium", "chromium-browser", "google-chrome", "google-chrome-stable"} {
		if found, err := exec.LookPath(candidate); err == nil {
			browserPath = found
			break
		}
	}
	if browserPath == "" {
		t.Skip("Chromium is unavailable")
	}
	sha := strings.Repeat("b", 40)
	role := "copy editor"
	writer := "author"
	post := &models.Post{AuthorObjects: []models.Author{
		{ID: "writer", Name: "The Writer", Role: &writer},
		{ID: "editor", Name: "Main Editor", Role: &role},
		{ID: "other", Name: "Another Editor", Role: &role},
	}, Extra: map[string]interface{}{"editorial_history": []interface{}{map[string]interface{}{
		"schema_version": 1, "author": "editor", "role": "copy-editor", "action": "edited",
		"commit": sha, "base_commit": strings.Repeat("a", 40), "path": "pages/post.md",
		"file_url": "https://github.com/example/site/blob/" + sha + "/pages/post.md",
		"summary":  "Corrected wording.",
	}}}}
	engine, err := templates.NewEngine("")
	if err != nil {
		t.Fatal(err)
	}
	ctx := templates.NewContext(post, "", &models.Config{})
	byline, err := engine.Render("components/post_byline.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	history, err := engine.Render("components/editorial_history.html", ctx)
	if err != nil {
		t.Fatal(err)
	}
	fixture := `<!doctype html><html><head><meta name="viewport" content="width=device-width"></head><body>` + byline + history + `<script>
window.revisionRequests = 0;
window.failRevision = false;
window.fetch = async function(url) {
 window.revisionRequests++;
 if (window.failRevision) return new Response('', {status: 404});
 var before = url.includes('aaaaaaaa');
 var tail = Array(100).fill('unchanged context ').join('');
 return new Response('---\\ntitle: Revision\\n---\\n' + (before ? 'A tpyο.\\n<script>window.bad = true;<\\/script>' : 'A typo.\\n<img src=x onerror="window.bad = true">') + tail);
};
</script><script src="/js/editorial-history.js"></script></body></html>`
	var viewerRequests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/js/") {
			if r.URL.Path == "/js/editorial-diff.js" {
				viewerRequests.Add(1)
			}
			asset, readErr := themes.ReadStatic(strings.TrimPrefix(r.URL.Path, "/"))
			if readErr != nil {
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-Type", "text/javascript")
			if _, writeErr := w.Write(asset); writeErr != nil {
				t.Errorf("write asset: %v", writeErr)
			}
			return
		}
		if _, writeErr := fmt.Fprint(w, fixture); writeErr != nil {
			t.Errorf("write fixture: %v", writeErr)
		}
	}))
	defer server.Close()
	allocator, cancelAllocator := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath(browserPath), chromedp.Flag("no-sandbox", true), chromedp.Flag("disable-dev-shm-usage", true),
	)...)
	defer cancelAllocator()
	browser, cancelBrowser := chromedp.NewContext(allocator)
	defer cancelBrowser()
	browser, cancelTimeout := context.WithTimeout(browser, 45*time.Second)
	defer cancelTimeout()
	var requests int
	var safe bool
	if err := chromedp.Run(browser,
		chromedp.Navigate(server.URL),
		chromedp.WaitReady(`[data-editorial-show]:not([hidden])`),
		chromedp.Evaluate(`window.revisionRequests`, &requests),
	); err != nil {
		t.Fatal(err)
	}
	if requests != 0 || viewerRequests.Load() != 0 {
		t.Fatal("revision data or viewer loaded before reader action")
	}
	if err := chromedp.Run(browser,
		chromedp.Evaluate(`document.querySelector('.editorial-history > summary').focus()`, nil),
		chromedp.KeyEvent("\r"),
		chromedp.WaitVisible(`[data-editorial-show]`),
		chromedp.Evaluate(`window.failRevision = true`, nil),
		chromedp.Click(`[data-editorial-show]`),
		chromedp.Poll(`document.querySelector('[data-editorial-show]').textContent === 'Retry changes'`, nil),
		chromedp.Evaluate(`window.failRevision = false`, nil),
		chromedp.Click(`[data-editorial-show]`),
		chromedp.WaitVisible(`.editorial-history ins`),
		chromedp.Evaluate(`!window.bad && document.querySelectorAll('[data-editorial-result] script, [data-editorial-result] img').length === 0 && document.querySelector('[data-editorial-show]').getAttribute('aria-expanded') === 'true' && document.querySelectorAll('.editorial-history__omitted').length > 0 && [...document.querySelectorAll('.editorial-history__omitted')].every(x => !x.open)`, &safe),
	); err != nil {
		t.Fatal(err)
	}
	if !safe || viewerRequests.Load() != 1 {
		t.Fatalf("unsafe render or duplicate viewer loads: safe=%v loads=%d", safe, viewerRequests.Load())
	}
	if err := chromedp.Run(browser,
		chromedp.Click(`[data-editorial-show]`),
		chromedp.Evaluate(`document.querySelector('.editorial-history__comparison').hidden`, &safe),
		chromedp.Evaluate(`window.dispatchEvent(new Event('view-transition-complete')); window.dispatchEvent(new Event('view-transition-complete'))`, nil),
		chromedp.Click(`[data-editorial-show]`),
		chromedp.Evaluate(`!document.querySelector('.editorial-history__comparison').hidden`, &safe),
	); err != nil {
		t.Fatal(err)
	}
	if !safe || viewerRequests.Load() != 1 {
		t.Fatal("toggle/view-transition initialization duplicated work")
	}
}
