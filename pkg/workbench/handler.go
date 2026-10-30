// Package workbench provides a server backed Markdown editor that can be embedded
// in Go applications. The handler owns no listeners and stores no drafts.
package workbench

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
	"github.com/WaylonWalker/markata/pkg/lsp"
	"github.com/WaylonWalker/markata/pkg/models"
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
	siteConfig     *models.Config
	workspace      *editorworkspace.Workspace
	language       http.Handler
	languageConfig lsp.BrowserConfig
	appearance     *plugins.StandaloneAppearance
	choices        map[string][]string
	assets         map[string][]byte
	index          []byte
	bridge         string
	gate           chan struct{}
	lastSource     string
	lastResult     *plugins.StandaloneResult
}

// New builds immutable appearance assets once. Reuse the handler for all requests.
func New() (*Handler, error) { return NewWithWorkspace(context.Background(), ".") }

// NewWithWorkspace enables LSP for a server-selected workspace. Cancel ctx on shutdown.
func NewWithWorkspace(ctx context.Context, root string) (*Handler, error) {
	workspace, err := editorworkspace.New(root)
	if err != nil {
		return nil, err
	}
	config, err := lsp.BrowserWorkspace(root)
	if err != nil {
		return nil, err
	}
	language, err := lsp.NewWebSocketHandler(ctx, root, log.Default())
	if err != nil {
		return nil, err
	}
	siteConfig, err := workspace.Config()
	if err != nil {
		return nil, err
	}
	a, err := plugins.BuildStandaloneAppearanceForSite(siteConfig, root)
	if err != nil {
		return nil, err
	}
	h := &Handler{siteConfig: siteConfig, workspace: workspace, language: language, languageConfig: config, appearance: a, assets: a.Assets, gate: make(chan struct{}, 1), choices: map[string][]string{"data-theme": {"light", "dark"}, "data-text-size": {"small", "medium", "large", "x-large"}}}
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
	choices, err := json.Marshal(h.choices)
	if err != nil {
		return nil, fmt.Errorf("encode appearance choices: %w", err)
	}
	bridge, err := web.ReadFile("web/preview.js")
	if err != nil {
		return nil, err
	}
	h.bridge = "<script>" + strings.ReplaceAll(string(bridge), "__MARKATA_CHOICES__", string(choices)) + "</script>"
	page, err := web.ReadFile("web/index.html")
	if err != nil {
		return nil, err
	}
	// Editor CSS precedes theme overrides. Site article layout belongs in the iframe.
	editorHead := `<link rel="stylesheet" href="/editor/codemirror/lib/codemirror.css"><link rel="stylesheet" href="/editor/codemirror/addon/hint/show-hint.css"><link rel="stylesheet" href="/editor/lsp.css"><script src="/editor/lsp.js"></script>`
	for _, name := range []string{"lib/codemirror", "mode/xml/xml", "mode/markdown/markdown", "mode/gfm/gfm", "mode/javascript/javascript", "mode/python/python", "mode/css/css", "mode/htmlmixed/htmlmixed", "addon/mode/overlay", "addon/edit/matchbrackets", "addon/hint/show-hint", "keymap/vim", "keymap/emacs", "keymap/sublime"} {
		editorHead += `<script src="/editor/codemirror/` + name + `.js"></script>`
	}
	controls := ""
	for _, name := range a.Styles {
		if name != "css/main.css" && name != "css/components.css" {
			controls += h.versionedStyleLink(name)
		}
	}
	scripts := ""
	for _, name := range a.Scripts {
		scripts += `<script src="/appearance/` + name + `"></script>`
	}
	h.index = []byte(strings.NewReplacer("<!--EDITOR_HEAD-->", editorHead, "<!--APPEARANCE_HEAD-->", controls, "<!--PICKER-->", a.PickerHTML, "<!--APPEARANCE_SCRIPTS-->", scripts).Replace(string(page)))
	if siteConfig != nil {
		attrs := `<html lang="en"`
		for _, key := range []string{"data-theme", "data-palette", "data-aesthetic", "data-fontpack", "data-text-size"} {
			attrs += " " + key + `="` + html.EscapeString(a.Defaults[key]) + `"`
		}
		h.index = regexp.MustCompile(`<html[^>]*>`).ReplaceAll(h.index, []byte(attrs+">"))
	}
	h.versionEditorAssets()
	return h, nil
}

// Version editor references so restarting on the same port cannot reuse stale code.
var editorAssetReference = regexp.MustCompile(`(?:src|href)="/(editor|appearance)/([^"?]+)"`)

func (h *Handler) versionEditorAssets() {
	h.index = editorAssetReference.ReplaceAllFunc(h.index, func(reference []byte) []byte {
		match := editorAssetReference.FindSubmatch(reference)
		var data []byte
		if string(match[1]) == "editor" {
			var err error
			data, err = web.ReadFile("web/" + string(match[2]))
			if err != nil {
				return reference
			}
		} else {
			data = h.assets[string(match[2])]
		}
		hash := sha256.Sum256(data)
		return []byte(strings.TrimSuffix(string(reference), `"`) + fmt.Sprintf("?v=%x", hash[:8]) + `"`)
	})
}

func (h *Handler) versionedStyleLink(name string) string {
	hash := sha256.Sum256(h.assets[name])
	return `<link rel="stylesheet" href="/appearance/` + name + fmt.Sprintf("?v=%x", hash[:8]) + `">`
}

func styleLink(name string) string { return `<link rel="stylesheet" href="/appearance/` + name + `">` }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	switch {
	case (r.URL.Path == "/workspace" || r.URL.Path == "/workspace/file") && r.Method == http.MethodGet:
		peer, _, peerErr := net.SplitHostPort(r.RemoteAddr)
		if peerErr != nil || net.ParseIP(peer) == nil || !net.ParseIP(peer).IsLoopback() {
			http.Error(w, "local workspace access required", http.StatusForbidden)
			return
		}
		host, hostErr := url.Parse("//" + r.Host)
		if hostErr != nil || host.User != nil || host.Path != "" || (host.Hostname() != "localhost" && (net.ParseIP(host.Hostname()) == nil || !net.ParseIP(host.Hostname()).IsLoopback())) {
			http.Error(w, "local workspace access required", http.StatusForbidden)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/workspace/file" {
			source, err := h.workspace.Read(r.URL.Query().Get("path"))
			if err != nil {
				fail(w, 404, err.Error())
				return
			}
			if err := json.NewEncoder(w).Encode(map[string]string{"source": source}); err != nil {
				return
			}
			return
		}
		files, err := h.workspace.Files()
		if err != nil {
			fail(w, 422, err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"files": files, "rootUri": h.languageConfig.RootURI}); err != nil {
			return
		}
	case r.URL.Path == "/lsp":
		h.language.ServeHTTP(w, r)
	case r.URL.Path == "/lsp/config" && r.Method == http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		if err := json.NewEncoder(w).Encode(h.languageConfig); err != nil {
			return
		}
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
		w.Header().Set("Cache-Control", "no-cache")
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
		w.Header().Set("Cache-Control", "no-cache")
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
	response, err := h.prepareRender(r, req)
	if err != nil {
		if r.Context().Err() == nil {
			fail(w, 422, err.Error())
		}
		return
	}
	ms := float64(time.Since(started).Microseconds()) / 1000
	response.RenderMS = ms
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server-Timing", fmt.Sprintf("markata;dur=%.3f", ms))
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *Handler) prepareRender(r *http.Request, req *renderRequest) (*renderResponse, error) {
	defer func() { <-h.gate }()
	if err := r.Context().Err(); err != nil {
		return nil, err
	}
	result := h.lastResult
	cached := result != nil && h.lastSource == *req.Markdown
	if !cached {
		var err error
		result, err = plugins.RenderStandalone(*req.Markdown, plugins.StandaloneOptions{Config: h.siteConfig, Unsafe: true, AssetBase: "/appearance"})
		if err != nil {
			return nil, err
		}
		for _, asset := range result.Assets {
			if !strings.HasPrefix(asset.ID, "generated-") && asset.Content != "" {
				h.assets[asset.ID] = []byte(asset.Content)
			}
		}
		h.lastSource, h.lastResult = *req.Markdown, result
	}
	response := &renderResponse{HTML: result.HTML, FeatureHeadHTML: h.featureHead(result), BodyEndHTML: result.BodyEndHTML, InputBytes: len(*req.Markdown), HTMLBytes: len(result.HTML), Cached: cached}
	if req.IncludeDocument == nil || *req.IncludeDocument {
		document := h.document(result, req.Appearance)
		response.Document, response.DocumentBytes = &document, len(document)
	}
	return response, nil
}

type renderRequest struct {
	Markdown        *string           `json:"markdown"`
	Appearance      map[string]string `json:"appearance"`
	IncludeDocument *bool             `json:"include_document"`
}

// renderResponse matches the Python reference's full-document/fragment contract.
type renderResponse struct {
	Document        *string `json:"document"`
	HTML            string  `json:"html"`
	FeatureHeadHTML string  `json:"feature_head_html"`
	BodyEndHTML     string  `json:"body_end_html"`
	RenderMS        float64 `json:"render_ms"`
	InputBytes      int     `json:"input_bytes"`
	HTMLBytes       int     `json:"html_bytes"`
	DocumentBytes   int     `json:"document_bytes"`
	Cached          bool    `json:"cached"`
}

// featureHead omits styles already owned by the persistent appearance shell.
// Mark feature styles so subsequent fragments can remove unused styles as well.
func (h *Handler) featureHead(result *plugins.StandaloneResult) string {
	head := result.HeadHTML
	for _, name := range h.appearance.Styles {
		head = strings.ReplaceAll(head, styleLink(name)+"\n", "")
	}
	return strings.NewReplacer("<link ", `<link data-markata-render-head="true" `, "<style>", `<style data-markata-render-head="true">`).Replace(head)
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
	for key, value := range h.appearance.Defaults {
		defaults[key] = value
	}
	for key, value := range attributes {
		defaults[key] = value
	}
	if _, ok := attributes["data-palette"]; !ok && defaults["data-theme"] == "light" {
		if h.siteConfig == nil {
			defaults["data-palette"] = "default-light"
		}
	}
	attrs := ""
	for _, key := range []string{"data-theme", "data-palette", "data-aesthetic", "data-fontpack", "data-text-size"} {
		attrs += " " + key + `="` + html.EscapeString(defaults[key]) + `"`
	}
	styles := ""
	for _, name := range h.appearance.Styles {
		styles += h.versionedStyleLink(name)
	}
	bg, fg := "#131721", "#e5e7eb"
	if defaults["data-theme"] == "light" {
		bg, fg = "#fff", "#1f2937"
	}
	initial := `<style>html{color-scheme:` + defaults["data-theme"] + `;background:` + bg + `}body{background:var(--color-background,` + bg + `);color:var(--color-text,` + fg + `)}</style>`
	return `<!doctype html><html lang="en"` + attrs + `><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Markdown preview</title>` + initial + styles + h.featureHead(result) + `</head><body><main class="main-content"><div class="content-wrapper"><div class="post-shell"><article class="post h-entry"><div class="post-content e-content">` + result.HTML + `</div></article></div></div></main>` + strings.ReplaceAll(result.BodyEndHTML, "<script", `<script data-markata-generated="true"`) + h.bridge + `</body></html>`
}
