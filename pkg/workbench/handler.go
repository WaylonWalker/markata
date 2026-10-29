// Package workbench provides a server backed Markdown editor that can be embedded
// in Go applications. The handler owns no listeners and stores no drafts.
package workbench

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/plugins"
)

//go:embed web
var web embed.FS

const maxSource = 200000

var manifestPattern = regexp.MustCompile(`--(palette|aesthetic|fontpack)-manifest:\s*'((?:\\.|[^'])*)'`)

// Handler serves the editor, render API and bundled appearance assets.
// Mount it at the root of an origin. Authored HTML is enabled for trusted local
// editing; isolate its preview using an opaque-origin iframe sandbox.
type Handler struct {
	appearance *plugins.StandaloneAppearance
	choices    map[string][]string
	assets     map[string][]byte
	index      []byte
	gate       chan struct{}
	lastSource string
	lastResult *plugins.StandaloneResult
}

// New builds immutable appearance assets once. Reuse the handler for all requests.
func New() (*Handler, error) {
	a, err := plugins.BuildStandaloneAppearance()
	if err != nil {
		return nil, err
	}
	h := &Handler{appearance: a, assets: a.Assets, gate: make(chan struct{}, 1), choices: map[string][]string{"data-theme": {"light", "dark"}, "data-text-size": {"small", "medium", "large", "x-large"}}}
	h.assets["css/fonts.css"] = bytes.ReplaceAll(h.assets["css/fonts.css"], []byte("/assets/fonts/"), []byte("/appearance/assets/fonts/"))
	for _, name := range []string{"palette", "aesthetic", "fonts"} {
		match := manifestPattern.FindSubmatch(h.assets["css/"+name+".css"])
		if len(match) != 3 {
			return nil, fmt.Errorf("missing %s manifest", name)
		}
		var entries []struct {
			Name string `json:"name"`
		}
		raw := strings.ReplaceAll(strings.ReplaceAll(string(match[2]), `\'`, `'`), `\\`, `\`)
		if err = json.Unmarshal([]byte(raw), &entries); err != nil {
			return nil, fmt.Errorf("decode %s manifest: %w", name, err)
		}
		for _, entry := range entries {
			h.choices["data-"+string(match[1])] = append(h.choices["data-"+string(match[1])], entry.Name)
		}
	}
	page, err := web.ReadFile("web/index.html")
	if err != nil {
		return nil, err
	}
	// Editor CSS precedes theme overrides. Site article layout belongs in the iframe.
	editorHead := `<link rel="stylesheet" href="/editor/codemirror/lib/codemirror.css">`
	for _, name := range []string{"lib/codemirror", "mode/xml/xml", "mode/markdown/markdown", "mode/gfm/gfm", "mode/javascript/javascript", "mode/python/python", "mode/css/css", "mode/htmlmixed/htmlmixed", "addon/mode/overlay", "addon/edit/matchbrackets", "keymap/vim", "keymap/emacs", "keymap/sublime"} {
		editorHead += `<script src="/editor/codemirror/` + name + `.js"></script>`
	}
	controls := ""
	for _, name := range a.Styles {
		if name != "css/main.css" && name != "css/components.css" {
			controls += styleLink(name)
		}
	}
	scripts := ""
	for _, name := range a.Scripts {
		scripts += `<script src="/appearance/` + name + `"></script>`
	}
	h.index = []byte(strings.NewReplacer("<!--EDITOR_HEAD-->", editorHead, "<!--APPEARANCE_HEAD-->", controls, "<!--PICKER-->", a.PickerHTML, "<!--APPEARANCE_SCRIPTS-->", scripts).Replace(string(page)))
	return h, nil
}

func styleLink(name string) string { return `<link rel="stylesheet" href="/appearance/` + name + `">` }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	switch {
	case r.URL.Path == "/render":
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.render(w, r)
	case r.Method != http.MethodGet && r.Method != http.MethodHead:
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	case r.URL.Path == "/":
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'self'")
		serve(w, r, "index.html", h.index)
	case strings.HasPrefix(r.URL.Path, "/appearance/"):
		// The only mutation occurs under the render gate. Copy the slice header while
		// holding it, then let concurrent asset requests proceed independently.
		select {
		case h.gate <- struct{}{}:
		case <-r.Context().Done():
			return
		}
		data, ok := h.assets[strings.TrimPrefix(r.URL.Path, "/appearance/")]
		<-h.gate
		if !ok {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		serve(w, r, r.URL.Path, data)
	case strings.HasPrefix(r.URL.Path, "/editor/"):
		name := strings.TrimPrefix(r.URL.Path, "/editor/")
		data, err := web.ReadFile("web/" + name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		serve(w, r, name, data)
	default:
		http.NotFound(w, r)
	}
}

func serve(w http.ResponseWriter, r *http.Request, name string, data []byte) {
	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(name)))
	w.Header().Set("ETag", fmt.Sprintf(`"%x"`, sha256.Sum256(data)))
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(data))
}

func fail(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		return
	}
}

func (h *Handler) render(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	if origin := r.Header.Get("Origin"); origin != "" {
		u, err := url.Parse(origin)
		if err != nil || u.Host != r.Host || (u.Scheme != "http" && u.Scheme != "https") {
			fail(w, 403, "cross-origin rendering is disabled")
			return
		}
	}
	if r.Header.Get("Sec-Fetch-Site") == "cross-site" {
		fail(w, 403, "cross-origin rendering is disabled")
		return
	}
	kind, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || kind != "application/json" {
		fail(w, 415, "send Content-Type: application/json")
		return
	}
	req := h.decodeRequest(w, r)
	if req == nil {
		return
	}
	select {
	case h.gate <- struct{}{}:
	case <-r.Context().Done():
		return
	}
	started := time.Now()
	document, cached, err := h.prepareRender(r, req)
	if err != nil {
		if r.Context().Err() == nil {
			fail(w, 422, err.Error())
		}
		return
	}
	ms := float64(time.Since(started).Microseconds()) / 1000
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server-Timing", fmt.Sprintf("markata;dur=%.3f", ms))
	if err := json.NewEncoder(w).Encode(map[string]any{"document": document, "render_ms": ms, "input_bytes": len(*req.Markdown), "document_bytes": len(document), "cached": cached}); err != nil {
		return
	}
}

func (h *Handler) prepareRender(r *http.Request, req *renderRequest) (document string, cacheHit bool, renderErr error) {
	defer func() { <-h.gate }()
	if err := r.Context().Err(); err != nil {
		return "", false, err
	}
	result := h.lastResult
	cached := result != nil && h.lastSource == *req.Markdown
	if !cached {
		var err error
		result, err = plugins.RenderStandalone(*req.Markdown, plugins.StandaloneOptions{Unsafe: true, AssetBase: "/appearance"})
		if err != nil {
			return "", false, err
		}
		for _, asset := range result.Assets {
			if !strings.HasPrefix(asset.ID, "generated-") && asset.Content != "" {
				h.assets[asset.ID] = []byte(asset.Content)
			}
		}
		h.lastSource, h.lastResult = *req.Markdown, result
	}
	return h.document(result, req.Appearance), cached, nil
}

type renderRequest struct {
	Markdown   *string           `json:"markdown"`
	Appearance map[string]string `json:"appearance"`
}

func (h *Handler) decodeRequest(w http.ResponseWriter, r *http.Request) *renderRequest {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req renderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err == nil {
		var extra any
		if trailing := decoder.Decode(&extra); !errors.Is(trailing, io.EOF) {
			if trailing != nil {
				err = trailing
			} else {
				err = errors.New("expected one JSON object")
			}
		}
	}
	if err != nil {
		var large *http.MaxBytesError
		if errors.As(err, &large) {
			fail(w, 413, "request exceeds 1 MiB")
		} else {
			fail(w, 400, "expected one JSON object with markdown and optional appearance")
		}
		return nil
	}
	if req.Markdown == nil {
		fail(w, 400, "markdown must be a string")
		return nil
	}
	if len(*req.Markdown) > maxSource {
		fail(w, 413, "document exceeds 200000 UTF-8 bytes")
		return nil
	}
	for key, value := range req.Appearance {
		valid := false
		for _, choice := range h.choices[key] {
			if choice == value {
				valid = true
				break
			}
		}
		if !valid {
			fail(w, 400, "unknown appearance choice: "+key)
			return nil
		}
	}
	return &req
}

func (h *Handler) document(result *plugins.StandaloneResult, attributes map[string]string) string {
	defaults := map[string]string{"data-theme": "dark", "data-palette": "default-dark", "data-aesthetic": "minimal", "data-fontpack": "system", "data-text-size": "large"}
	for key, value := range attributes {
		defaults[key] = value
	}
	if _, ok := attributes["data-palette"]; !ok && defaults["data-theme"] == "light" {
		defaults["data-palette"] = "default-light"
	}
	attrs := ""
	for _, key := range []string{"data-theme", "data-palette", "data-aesthetic", "data-fontpack", "data-text-size"} {
		attrs += " " + key + `="` + html.EscapeString(defaults[key]) + `"`
	}
	styles := ""
	for _, name := range h.appearance.Styles {
		styles += styleLink(name)
	}
	choices, err := json.Marshal(h.choices)
	if err != nil {
		panic("appearance choices are not serializable: " + err.Error())
	}
	bridge := `<script>(()=>{const allowed=` + string(choices) + `;addEventListener('message',e=>{if(e.source!==parent||e.data?.type!=='appearance')return;for(const [key,values] of Object.entries(allowed)){const v=e.data.attributes?.[key];if(values.includes(v))document.documentElement.setAttribute(key,v)}dispatchEvent(new Event('theme-change'))})})();</script>`
	bg, fg := "#131721", "#e5e7eb"
	if defaults["data-theme"] == "light" {
		bg, fg = "#fff", "#1f2937"
	}
	initial := `<style>html{color-scheme:` + defaults["data-theme"] + `;background:` + bg + `}body{background:var(--color-background,` + bg + `);color:var(--color-text,` + fg + `)}</style>`
	return `<!doctype html><html lang="en"` + attrs + `><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Markdown preview</title>` + initial + result.HeadHTML + styles + `</head><body><main class="main-content"><div class="content-wrapper"><div class="post-shell"><article class="post h-entry"><div class="post-content e-content">` + result.HTML + `</div></article></div></div></main>` + result.BodyEndHTML + bridge + `</body></html>`
}
