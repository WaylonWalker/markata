package workbench

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/html"
)

func newHandler(t testing.TB) *Handler {
	t.Helper()
	h, err := New()
	if err != nil {
		t.Fatal(err)
	}
	return h
}
func request(h *Handler, body string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, "http://localhost/render", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestHandler_RenderAndAssets(t *testing.T) {
	h := newHandler(t)
	source := "# Hello :lucide-house:\n\n!!! note\n    Local preview.\n\n```go\nfmt.Println(42)\n```\n\n```mermaid\ngraph LR\nA-->B\n```\n\n```chartjs\n{\"type\":\"bar\",\"data\":{\"labels\":[\"A\"],\"datasets\":[{\"data\":[1]}]}}\n```"
	body, err := json.Marshal(map[string]any{"markdown": source, "appearance": map[string]string{"data-theme": "light"}})
	if err != nil {
		t.Fatal(err)
	}
	var firstDocument string
	for i := 0; i < 2; i++ {
		w := request(h, string(body))
		if w.Code != 200 {
			t.Fatalf("%d: %s", w.Code, w.Body)
		}
		var response struct {
			Document      string `json:"document"`
			Cached        bool   `json:"cached"`
			InputBytes    int    `json:"input_bytes"`
			DocumentBytes int    `json:"document_bytes"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatal(err)
		}
		if response.Cached != (i == 1) || response.InputBytes != len(source) || response.DocumentBytes != len(response.Document) {
			t.Fatalf("invalid metrics: %+v", response)
		}
		for _, part := range []string{`data-theme="light"`, `data-palette="default-light"`, `class="post h-entry"`, `<svg`, `mermaid`, `chart`} {
			if !strings.Contains(response.Document, part) {
				t.Errorf("missing %s", part)
			}
		}
		if i == 1 && firstDocument != response.Document {
			t.Error("cached render changed document")
		}
		firstDocument = response.Document
		z := html.NewTokenizer(strings.NewReader(response.Document))
		for {
			kind := z.Next()
			if kind == html.ErrorToken {
				break
			}
			if kind != html.StartTagToken {
				continue
			}
			tag := z.Token()
			if tag.Data != "link" && tag.Data != "script" {
				continue
			}
			for _, a := range tag.Attr {
				if a.Key != "src" && a.Key != "href" {
					continue
				}
				if !strings.HasPrefix(a.Val, "/appearance/") {
					t.Fatalf("unexpected asset %s", a.Val)
				}
				asset := httptest.NewRecorder()
				h.ServeHTTP(asset, httptest.NewRequest("GET", a.Val, http.NoBody))
				if asset.Code != 200 || asset.Body.Len() == 0 {
					t.Fatalf("missing asset %s: %d", a.Val, asset.Code)
				}
				if asset.Header().Get("Access-Control-Allow-Origin") != "*" {
					t.Error("sandbox asset lacks CORS")
				}
			}
		}
	}
}

func TestHandler_RequestContract(t *testing.T) {
	h := newHandler(t)
	cases := []struct {
		name, body string
		status     int
	}{
		{"missing", `{}`, 400}, {"null", `{"markdown":null}`, 400}, {"wrong type", `{"markdown":42}`, 400}, {"invalid", `{`, 400}, {"trailing", `{"markdown":"ok"} {}`, 400}, {"unknown", `{"markdown":"ok","unknown":true}`, 400},
		{"bad appearance", `{"markdown":"ok","appearance":{"data-theme":"evil"}}`, 400}, {"unknown appearance", `{"markdown":"ok","appearance":{"onclick":"evil"}}`, 400},
		{"source bytes", `{"markdown":"` + strings.Repeat("é", 100001) + `"}`, 413}, {"body bytes", `{"markdown":"` + strings.Repeat("a", 1<<20) + `"}`, 413},
		{"trailing oversize", `{"markdown":"ok"}` + strings.Repeat(" ", 1<<20), 413},
		{"empty", `{"markdown":""}`, 200}, {"limit", `{"markdown":"` + strings.Repeat("a", maxSource) + `"}`, 200},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			w := request(h, tt.body)
			if w.Code != tt.status {
				t.Fatalf("got %d want %d: %.200s", w.Code, tt.status, w.Body.String())
			}
		})
	}
	for _, test := range []struct {
		method, path, origin, contentType string
		code                              int
	}{
		{"GET", "/render", "", "", 405}, {"POST", "/", "", "", 405}, {"POST", "/render", "https://other.test", "application/json", 403}, {"POST", "/render", "null", "application/json", 403}, {"POST", "/render", "", "text/plain", 415},
		{"GET", "/appearance/missing.css", "", "", 404}, {"GET", "/editor/../handler.go", "", "", 404}, {"GET", "/missing", "", "", 404}, {"GET", "/", "", "", 200}, {"HEAD", "/", "", "", 200}, {"GET", "/editor/app.js", "", "", 200}, {"GET", "/editor/codemirror/lib/codemirror.js", "", "", 200},
	} {
		t.Run(fmt.Sprintf("%s %s %s", test.method, test.path, test.origin), func(t *testing.T) {
			r := httptest.NewRequest(test.method, "http://localhost"+test.path, strings.NewReader(`{"markdown":"hi"}`))
			r.Header.Set("Origin", test.origin)
			r.Header.Set("Content-Type", test.contentType)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != test.code {
				t.Fatalf("got %d want %d", w.Code, test.code)
			}
		})
	}
}

func TestHandler_CancellationAndConcurrency(t *testing.T) {
	h := newHandler(t)
	h.gate <- struct{}{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := httptest.NewRequest("POST", "http://localhost/render", strings.NewReader(`{"markdown":"hello"}`)).WithContext(ctx)
	r.Header.Set("Content-Type", "application/json")
	done := make(chan struct{})
	go func() { h.ServeHTTP(httptest.NewRecorder(), r); close(done) }()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("canceled request blocked")
	}
	<-h.gate
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 3; j++ {
				w := request(h, fmt.Sprintf(`{"markdown":"# Document %d %d"}`, i, j))
				if w.Code != 200 {
					t.Errorf("render: %d", w.Code)
				}
				h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/appearance/css/main.css", http.NoBody))
			}
		}(i)
	}
	wg.Wait()
}

func TestHandler_ConditionalAsset(t *testing.T) {
	h := newHandler(t)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/editor/app.js", http.NoBody))
	r := httptest.NewRequest("GET", "/editor/app.js", http.NoBody)
	r.Header.Set("If-None-Match", w.Header().Get("ETag"))
	cached := httptest.NewRecorder()
	h.ServeHTTP(cached, r)
	if cached.Code != 304 {
		t.Fatalf("got %d", cached.Code)
	}
}

func BenchmarkHandler_Render(b *testing.B) {
	h := newHandler(b)
	source := strings.Repeat("## A heading\n\nSome **Markdown** with a [link](https://example.com).\n\n", 100)
	body, err := json.Marshal(map[string]string{"markdown": source})
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.lastResult = nil
		w := request(h, string(body))
		if w.Code != 200 {
			b.Fatal(w.Code)
		}
	}
}
func BenchmarkHandler_Cached(b *testing.B) {
	h := newHandler(b)
	body := `{"markdown":"# Hello\n\nA repeated draft."}`
	request(h, body)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := request(h, body)
		if !bytes.Contains(w.Body.Bytes(), []byte(`"cached":true`)) {
			b.Fatal("cache miss")
		}
	}
}

// blockedWriter models a client that has stopped reading the response.
type blockedWriter struct {
	*httptest.ResponseRecorder
	entered, release chan struct{}
}

func (w *blockedWriter) Write(data []byte) (int, error) {
	close(w.entered)
	<-w.release
	return w.ResponseRecorder.Write(data)
}
func TestHandler_SlowResponseDoesNotHoldRenderGate(t *testing.T) {
	h := newHandler(t)
	w := &blockedWriter{httptest.NewRecorder(), make(chan struct{}), make(chan struct{})}
	r := httptest.NewRequest("POST", "http://localhost/render", strings.NewReader(`{"markdown":"# Slow reader"}`))
	r.Header.Set("Content-Type", "application/json")
	done := make(chan struct{})
	go func() { h.ServeHTTP(w, r); close(done) }()
	<-w.entered
	defer func() { close(w.release); <-done }()
	select {
	case h.gate <- struct{}{}:
		<-h.gate
	case <-time.After(time.Second):
		t.Fatal("response encoding holds render gate")
	}
}

func TestHandler_FragmentUpdatesOmitSharedCSS(t *testing.T) {
	h := newHandler(t)
	source := "# Update\n\n```go\nfmt.Println(1)\n```\n\n!!! note\n    Conditional styles."
	var first renderResponse
	for _, include := range []bool{true, false, false} {
		body, err := json.Marshal(map[string]any{"markdown": source, "include_document": include})
		if err != nil {
			t.Fatal(err)
		}
		w := request(h, string(body))
		if w.Code != http.StatusOK {
			t.Fatalf("status %d: %s", w.Code, w.Body)
		}
		var result renderResponse
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if result.HTMLBytes != len(result.HTML) || result.InputBytes != len(source) {
			t.Fatal("incorrect fragment byte counts")
		}
		if include && result.Document == nil {
			t.Fatal("initial document missing")
		}
		for _, name := range h.appearance.Styles {
			if strings.Contains(result.FeatureHeadHTML, name) {
				t.Errorf("fragment repeats shared style %s", name)
			}
			if include && strings.Count(*result.Document, styleLink(name)) != 1 {
				t.Errorf("document must include %s exactly once", name)
			}
		}
		if !strings.Contains(result.FeatureHeadHTML, `data-markata-render-head="true"`) || !strings.Contains(result.FeatureHeadHTML, "css/code.css") {
			t.Fatal("missing marked conditional styles")
		}
		if include {
			if result.Document == nil {
				t.Fatal("initial document missing")
			}
			first = result
			continue
		}
		if result.Document != nil || result.DocumentBytes != 0 {
			t.Fatal("fragment includes document shell")
		}
		if !result.Cached || result.HTML != first.HTML || result.FeatureHeadHTML != first.FeatureHeadHTML || result.BodyEndHTML != first.BodyEndHTML {
			t.Fatal("cached fragment differs from initial content")
		}
		if strings.Contains(w.Body.String(), "palette.css") || strings.Contains(w.Body.String(), "fontpack-manifest") {
			t.Fatal("fragment contains shared appearance CSS")
		}
	}
	w := request(h, `{"markdown":"# Plain","include_document":false}`)
	var plain renderResponse
	if err := json.Unmarshal(w.Body.Bytes(), &plain); err != nil {
		t.Fatal(err)
	}
	if plain.FeatureHeadHTML != "" || plain.BodyEndHTML != "" {
		t.Fatal("plain fragment keeps unused feature assets")
	}
	invalid := request(h, `{"markdown":"# Invalid","include_document":"false"}`)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("invalid option: %d", invalid.Code)
	}
}

func TestHandler_RenderProfileAndAssetMetadata(t *testing.T) {
	h := newHandler(t)
	for i, tc := range []struct {
		request string
		status  int
		profile string
		cached  bool
	}{
		{`{"markdown":"# Local","profile":"web"}`, 200, "web", false},
		{`{"markdown":"# Local","profile":"web"}`, 200, "web", true},
		{`{"markdown":"# Local","profile":"offline-interactive"}`, 200, "offline-interactive", false},
		{`{"markdown":"# Local","profile":"offline-interactive"}`, 200, "offline-interactive", true},
		{`{"markdown":"![external](https://example.com/x.png)","profile":"offline-interactive"}`, 422, "", false},
		{`{"markdown":"hi","profile":"offline-static"}`, 422, "", false},
		{`{"markdown":"hi","profile":"not-a-profile"}`, 422, "", false},
	} {
		w := request(h, tc.request)
		if w.Code != tc.status {
			t.Fatalf("case %d: %d want %d: %s", i, w.Code, tc.status, w.Body.String())
		}
		if tc.status != 200 {
			continue
		}
		var result struct {
			Profile  string          `json:"profile"`
			Features map[string]bool `json:"features"`
			Assets   []struct {
				ID        string `json:"id"`
				Kind      string `json:"kind"`
				Placement string `json:"placement"`
			} `json:"assets"`
			Cached bool `json:"cached"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if result.Profile != tc.profile || result.Cached != tc.cached || len(result.Assets) == 0 || result.Features == nil {
			t.Fatalf("case %d: missing or stale profile/manifest: %+v", i, result)
		}
		if result.Assets[0].ID != "css/variables.css" || result.Assets[0].Placement != "head" {
			t.Fatalf("unexpected asset ordering: %+v", result.Assets)
		}
	}
}
