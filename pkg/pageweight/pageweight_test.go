package pageweight

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestAnalyzePostAndActionableFindings(t *testing.T) {
	limits := Limits{PostHTMLBytes: 60, OtherHTMLBytes: 10000, EmbeddedJSONBytes: 10, InlineJSBytes: 10, SidebarJSONBytes: 10}
	html := []byte(`<article class="post h-entry"></article><script type="application/json" id="feed-sidebar-data">{"posts":[1,2,3,4,5]}</script><script src="/js/app.js"></script><script>var startup = "extra work";</script><img src="x"><link rel="stylesheet" href="a.css">`)
	page := Analyze("posts/a/index.html", html, limits)
	if page.Kind != "post" || page.ExternalScripts != 1 || page.StylesheetLinks != 1 || page.ImageElements != 1 {
		t.Fatalf("unexpected page counters: %+v", page)
	}
	if page.SidebarJSONBytes == 0 || page.EmbeddedJSONBytes == 0 || page.InlineJSBytes == 0 {
		t.Fatalf("missing inline measurements: %+v", page)
	}
	want := []string{"HTML_OVER_BUDGET", "EMBEDDED_JSON_OVER_BUDGET", "INLINE_JS_OVER_BUDGET", "SIDEBAR_JSON_OVER_BUDGET"}
	codes := make([]string, 0, len(page.Findings))
	for _, f := range page.Findings {
		if f.Advice == "" || f.BudgetBytes <= 0 || f.Bytes <= f.BudgetBytes {
			t.Fatalf("invalid finding: %+v", f)
		}
		codes = append(codes, f.Code)
	}
	if !reflect.DeepEqual(codes, want) {
		t.Fatalf("finding codes = %v, want %v", codes, want)
	}
}

func TestScanDeterministicAndReportNotScanned(t *testing.T) {
	root := t.TempDir()
	for name, body := range map[string]string{
		"b/index.html":          "<html><body>short</body></html>",
		"a/index.html":          "<html><body>" + strings.Repeat("x", 150) + "</body></html>",
		".markata/ignored.html": strings.Repeat("x", 100000),
		"binary.txt":            "not HTML",
	} {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	limits := Limits{OtherHTMLBytes: 100}
	first, err := Scan(root, limits)
	if err != nil {
		t.Fatal(err)
	}
	if first.Summary.Pages != 2 || first.Summary.WarnedPages != 1 || first.Pages[0].Path != "a/index.html" {
		t.Fatalf("unexpected report: %+v", first)
	}
	out := filepath.Join(root, ".markata", "page-weight.json")
	if err := WriteJSON(out, first); err != nil {
		t.Fatal(err)
	}
	second, err := Scan(root, limits)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("reports are not deterministic")
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	var decoded Report
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Measurement != "uncompressed_generated_html_source" || decoded.SchemaVersion != SchemaVersion {
		t.Fatalf("unexpected serialized report: %+v", decoded)
	}
}

func TestNoWarningsOnSmallPage(t *testing.T) {
	p := Analyze("index.html", []byte("<h1>Hi</h1>"), DefaultLimits())
	if len(p.Findings) != 0 {
		t.Fatalf("unexpected findings: %+v", p.Findings)
	}
}
