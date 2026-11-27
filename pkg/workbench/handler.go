// Package workbench provides a server backed Markdown editor that can be embedded
// in Go applications. The handler owns no listeners and stores no drafts.
package workbench

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/WaylonWalker/markata/pkg/editorworkspace"
	"github.com/WaylonWalker/markata/pkg/lifecycle"
	"github.com/WaylonWalker/markata/pkg/lsp"
	"github.com/WaylonWalker/markata/pkg/models"
	"github.com/WaylonWalker/markata/pkg/plugins"
	"github.com/WaylonWalker/markata/pkg/prose"
	"github.com/WaylonWalker/markata/pkg/visualizations"
)

const (
	imageUploadModeLocal   = "local"
	imageUploadModeDropper = "dropper"
	imageUploadModeAuto    = "auto"
)

const (
	workbenchTextSizeAttribute  = "data-text-size"
	workbenchThemeAttribute     = "data-theme"
	workbenchAestheticAttribute = "data-aesthetic"
	workbenchPaletteAttribute   = "data-palette"
	workbenchFontpackAttribute  = "data-fontpack"
	workbenchLightMode          = "light"
	workbenchDarkMode           = "dark"
)

//go:embed web
var web embed.FS

const maxSource = 200000

var manifestPattern = regexp.MustCompile(`--(palette|aesthetic|fontpack)-manifest:\s*'((?:\\.|[^'])*)'`)
var stylesheetLinkPattern = regexp.MustCompile(`(?i)<link\b[^>]*\brel\s*=\s*(?:"stylesheet"|'stylesheet')[^>]*>`)
var stylesheetHrefPattern = regexp.MustCompile(`(?i)\bhref\s*=\s*"([^"]+)"|\bhref\s*=\s*'([^']+)'`)

func wasmBundleVersion() string {
	wasm, err := web.ReadFile("web/wasm/markata.wasm")
	if err != nil {
		return ""
	}
	hash := sha256.Sum256(wasm)
	return fmt.Sprintf("%x", hash[:])
}

func wasmPreviewFingerprint(wasmVersion string, siteConfig *models.Config) string {
	if siteConfig != nil {
		return ""
	}
	return wasmVersion
}

type workspaceSource interface {
	Config() (*models.Config, error)
	Paths() ([]editorworkspace.File, error)
	Files() ([]editorworkspace.File, error)
	Read(string) (string, error)
	Write(string, string) error
	Create(string, string) error
	CreateWithImage(string, string, string, string, []byte) (string, error)
	ReadImage(string) ([]byte, error)
}

type workspaceMultiImageCreator interface {
	CreateWithImages(string, string, string, []editorworkspace.UploadedImage) ([]string, error)
}

type contentTemplateDefaultsWriter interface {
	SetContentTemplateDefaults(string, string, string) error
}

type workspaceCreateChecker interface {
	CheckCreate(string) error
}

// Handler serves the editor, render API and bundled appearance assets.
// Mount it at the root of an origin. Authored HTML is enabled for trusted local
// editing; isolate its preview using an opaque-origin iframe sandbox.
type Handler struct {
	store              *editorworkspace.Store
	root               string
	configRevision     uint64
	renderer           *plugins.StandaloneRenderer
	siteConfig         *models.Config
	workspace          workspaceSource
	language           http.Handler
	languageConfig     lsp.BrowserConfig
	appearance         *plugins.StandaloneAppearance
	choices            map[string][]string
	assets             map[string][]byte
	index              []byte
	bridge             string
	gate               chan struct{}
	workspaceWriteMu   sync.Mutex
	lastSource         string
	lastResult         *plugins.StandaloneResult
	wasmVersion        string
	previewFingerprint string
	hostedAuth         *hostedAuth
	gitPublisher       GitPublisher
	dropperUploader    DropperUploader
}

// New builds immutable appearance assets once. Reuse the handler for all requests.
func New() (*Handler, error) { return NewWithWorkspace(context.Background(), ".") }

// NewWithWorkspace enables LSP for a server-selected workspace. Cancel ctx on shutdown.
func NewWithWorkspace(ctx context.Context, root string) (*Handler, error) {
	return NewWithWorkspaceAndGitPublisher(ctx, root, nil)
}

// NewWithWorkspaceAndGitPublisher enables LSP for a server-selected workspace
// and delegates hosted remote Git operations to publisher. Cancel ctx on
// shutdown. A nil publisher preserves direct local Git behavior.
func NewWithWorkspaceAndGitPublisher(ctx context.Context, root string, publisher GitPublisher) (*Handler, error) {
	return NewWithWorkspaceAndGitPublisherAndDropperUploader(ctx, root, publisher, nil)
}

// NewWithWorkspaceAndGitPublisherAndDropperUploader enables LSP for a
// server-selected workspace, delegates hosted remote Git operations to
// publisher, and delegates hosted Dropper uploads to uploader. Nil services
// preserve the local editor's direct behavior. Cancel ctx on shutdown.
func NewWithWorkspaceAndGitPublisherAndDropperUploader(ctx context.Context, root string, publisher GitPublisher, uploader DropperUploader) (*Handler, error) {
	started := time.Now()
	workspace, err := editorworkspace.Open(ctx, root)
	if err != nil {
		return nil, err
	}
	log.Printf("editor startup: workspace discovery ready in %s", time.Since(started).Round(time.Millisecond))
	lspStarted := time.Now()
	config, err := lsp.BrowserWorkspace(root)
	if err != nil {
		_ = workspace.Close()
		return nil, err
	}
	// LSP setup can initialize a large site index; this is optional for
	// writing. Bring the HTTP editor online first and warm LSP asynchronously.
	language := &deferredLanguage{}
	siteConfig, revision := workspace.Configuration()
	appearanceStarted := time.Now()
	handler, err := buildPresentation(root, workspace, config, language, siteConfig)
	if err != nil {
		_ = workspace.Close()
		return nil, err
	}
	handler.store = workspace
	handler.configRevision = revision
	handler.gitPublisher = publisher
	handler.dropperUploader = uploader
	go func() {
		service, err := lsp.NewWebSocketHandlerWithIndexAndOriginCheck(ctx, root, log.Default(), workspace.Snapshot, handler.validLSPOrigin)
		if err != nil {
			if ctx.Err() == nil {
				log.Printf("editor startup: language service unavailable: %v", err)
			}
			return
		}
		language.ready(service)
		log.Printf("editor startup: LSP ready asynchronously in %s", time.Since(lspStarted).Round(time.Millisecond))
	}()
	log.Printf("editor startup: appearance and browser assets ready in %s (total %s)", time.Since(appearanceStarted).Round(time.Millisecond), time.Since(started).Round(time.Millisecond))
	return handler, nil
}

// Close releases resources owned by a workbench created without a cancelable context.
func (h *Handler) Close() error {
	if h.store != nil {
		return h.store.Close()
	}
	return nil
}

func buildPresentation(root string, workspace workspaceSource, config lsp.BrowserConfig, language http.Handler, siteConfig *models.Config) (*Handler, error) {
	a, err := plugins.BuildStandaloneAppearanceForSite(siteConfig, root)
	if err != nil {
		return nil, err
	}
	excalidrawVersion, err := visualizations.ExcalidrawRendererVersion()
	if err != nil {
		return nil, fmt.Errorf("fingerprint editor Excalidraw assets: %w", err)
	}
	excalidrawTheme := workbenchDarkMode
	if siteConfig != nil && siteConfig.Theme.FallbackMode == "light" {
		excalidrawTheme = "light"
	}
	wasmVersion := wasmBundleVersion()
	previewFingerprint := wasmPreviewFingerprint(wasmVersion, siteConfig)
	h := &Handler{root: root, renderer: plugins.NewStandaloneRenderer(plugins.StandaloneOptions{Config: siteConfig, Unsafe: true, AssetBase: "/appearance", MermaidMode: "client"}), siteConfig: siteConfig, workspace: workspace, language: language, languageConfig: config, appearance: a, assets: a.Assets, gate: make(chan struct{}, 1), wasmVersion: wasmVersion, previewFingerprint: previewFingerprint, choices: map[string][]string{workbenchThemeAttribute: {workbenchLightMode, workbenchDarkMode}, workbenchTextSizeAttribute: {"small", "medium", "large", "x-large"}}}
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
	// Only Docker (or another packaging step) supplies the version-matched
	// WASM bytes. Bare Go builds must remain fully functional without them.
	// Editor CSS precedes theme overrides. Site article layout belongs in the iframe.
	editorHead := fmt.Sprintf(`<script>window.MARKATA_EXCALIDRAW_RENDERER_VERSION=%q;window.MARKATA_EXCALIDRAW_EXPORT_THEME=%q;window.MARKATA_WASM_VERSION=%q;window.MARKATA_PREVIEW_FINGERPRINT=%q;</script><link rel="stylesheet" href="/editor/codemirror/lib/codemirror.css"><link rel="stylesheet" href="/editor/codemirror/addon/hint/show-hint.css"><link rel="stylesheet" href="/editor/codemirror/addon/dialog/dialog.css"><link rel="stylesheet" href="/editor/lsp.css"><script src="/editor/lsp.js"></script>`, excalidrawVersion, excalidrawTheme, wasmVersion, previewFingerprint)
	for _, name := range []string{"lib/codemirror", "mode/xml/xml", "mode/markdown/markdown", "mode/gfm/gfm", "mode/javascript/javascript", "mode/python/python", "mode/css/css", "mode/htmlmixed/htmlmixed", "addon/mode/overlay", "addon/edit/matchbrackets", "addon/hint/show-hint", "addon/search/searchcursor", "addon/dialog/dialog", "keymap/vim", "keymap/emacs", "keymap/sublime"} {
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
		for _, key := range []string{workbenchThemeAttribute, workbenchPaletteAttribute, workbenchAestheticAttribute, workbenchFontpackAttribute, workbenchTextSizeAttribute} {
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
	if h.hostedAuth != nil {
		w.Header().Set("Cache-Control", "private, no-store")
		w = &privateNoStoreResponseWriter{ResponseWriter: w}
		if !h.authorizeHostedRequest(r) {
			http.Error(w, "authenticated editor access required", http.StatusUnauthorized)
			return
		}
	}
	switch {
	case isWorkspaceRoute(r.URL.Path, r.Method):
		h.serveWorkspace(w, r)
	case r.URL.Path == "/lsp":
		h.language.ServeHTTP(w, r)
	case r.URL.Path == "/preview-fonts.css":
		h.servePreviewFonts(w, r)
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
		h.serveIndex(w, r)
	case r.URL.Path == "/manifest.webmanifest":
		h.handlePWAManifest(w, r)
	case r.URL.Path == "/sw.js":
		h.handlePWAServiceWorker(w, r)
	case r.URL.Path == "/pwa-icon.svg":
		h.handlePWAIcon(w, r)
	case r.URL.Path == "/pwa-icon-192.png" || r.URL.Path == rasterLargeIconPath:
		h.handlePWARasterIcon(w, r)
	case r.URL.Path == "/favicon.ico":
		data, err := web.ReadFile("web/favicon.ico")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		serve(w, r, "favicon.ico", data)
	case strings.HasPrefix(r.URL.Path, "/images/"):
		// Inline gallery Markdown uses /images/... while the cover image is
		// rendered through /workspace/image. Both must use the same contained
		// workspace reader, never a raw filesystem path.
		if h.workspace == nil {
			http.NotFound(w, r)
			return
		}
		name := strings.TrimPrefix(r.URL.Path, "/")
		data, err := h.workspace.ReadImage(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		serve(w, r, name, data)
	case strings.HasPrefix(r.URL.Path, "/appearance/"):
		w.Header().Set("Cache-Control", "no-cache")
		// The only mutation occurs under the render gate. Copy the slice header while
		// holding it, then let concurrent asset requests proceed independently.
		select {
		case h.gate <- struct{}{}:
		case <-r.Context().Done():
			return
		}
		if err := h.refreshPresentation(); err != nil {
			<-h.gate
			fail(w, 422, err.Error())
			return
		}
		data, ok := h.assets[strings.TrimPrefix(r.URL.Path, "/appearance/")]
		<-h.gate
		if !ok {
			http.NotFound(w, r)
			return
		}
		if h.hostedAuth == nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		serve(w, r, r.URL.Path, data)
	case strings.HasPrefix(r.URL.Path, "/editor/"):
		w.Header().Set("Cache-Control", "no-cache")
		name := strings.TrimPrefix(r.URL.Path, "/editor/")
		data, err := web.ReadFile("web/" + name)
		if err != nil && (name == "excalidraw-editor.js" || name == "excalidraw-editor.css" || strings.HasPrefix(name, "excalidraw-fonts/") || strings.HasPrefix(name, "excalidraw-editor-chunks/")) {
			assetName := name
			if strings.HasPrefix(assetName, "excalidraw-fonts/fonts/") {
				assetName = "excalidraw-fonts/" + strings.TrimPrefix(assetName, "excalidraw-fonts/fonts/")
			}
			data, err = visualizations.ReadExcalidrawAsset(assetName)
		}
		if err != nil {
			http.NotFound(w, r)
			return
		}
		if name == "wasm/markata.wasm" {
			if h.wasmVersion != "" {
				w.Header().Set("X-Markata-Wasm-Version", h.wasmVersion)
			}
			// The authenticated worker is served with the same authorization
			// boundary as /editor/. Content-encoding saves the large cold
			// download without compressing at runtime or caching private data.
			w.Header().Set("Vary", "Accept-Encoding")
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				if compressed, err := web.ReadFile("web/wasm/markata.wasm.gz"); err == nil {
					w.Header().Set("Content-Encoding", "gzip")
					serve(w, r, name, compressed)
					return
				}
			}
		}
		serve(w, r, name, data)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	if !h.acquire(r) {
		return
	}
	if err := h.refreshPresentation(); err != nil {
		<-h.gate
		fail(w, 422, err.Error())
		return
	}
	index := h.index
	<-h.gate
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Security-Policy", "frame-ancestors 'self'")
	serve(w, r, "index.html", index)
}

func serve(w http.ResponseWriter, r *http.Request, name string, data []byte) {
	if w.Header().Get("Cache-Control") == "" {
		w.Header().Set("Cache-Control", "public, max-age=3600, must-revalidate")
	}
	contentType := mime.TypeByExtension(path.Ext(name))
	if strings.EqualFold(path.Ext(name), ".woff2") {
		contentType = "font/woff2"
	} else if strings.EqualFold(path.Ext(name), ".wasm") {
		contentType = "application/wasm"
	}
	w.Header().Set("Content-Type", contentType)
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
	requestStarted := time.Now()
	w.Header().Set("Cache-Control", "no-store")
	if h.hostedAuth != nil {
		if !h.sameOriginWorkspaceMutation(r) {
			fail(w, http.StatusForbidden, "same-origin rendering is required")
			return
		}
	} else if origin := r.Header.Get("Origin"); origin != "" {
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
	if !h.acquire(r) {
		return
	}
	if err := h.refreshPresentation(); err != nil {
		<-h.gate
		fail(w, 422, err.Error())
		return
	}
	req := h.decodeRequest(w, r)
	if req == nil {
		<-h.gate
		return
	}
	started := time.Now()
	queueMS := float64(started.Sub(requestStarted).Microseconds()) / 1000
	response, err := h.prepareRender(r, req)
	if err != nil {
		if r.Context().Err() == nil {
			fail(w, 422, err.Error())
		}
		return
	}
	ms := float64(time.Since(started).Microseconds()) / 1000
	response.RenderMS = ms
	response.QueueMS = queueMS
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
	var index *lsp.Index
	var revision uint64
	if snapshot, ok := h.workspace.(interface{ Snapshot() (*lsp.Index, uint64) }); ok {
		index, revision = snapshot.Snapshot()
	}
	markdown := renderWorkbenchEmbeds(*req.Markdown, index, h.root)
	markdown = renderWorkbenchWikilinks(markdown, index)
	markdown, sourceAnchors := markSourceContainers(markdown)
	frontmatterPreview := renderFrontmatterPreview(frontmatterSource(req), markdown)
	// Index revisions change while the first preview is open. Plain Markdown
	// has no workspace references, so an index update does not invalidate it.
	if !strings.Contains(*req.Markdown, "[[") && !strings.Contains(*req.Markdown, "![embed](") {
		revision = 0
	}
	cacheKey := fmt.Sprintf("%d\x00%s", revision, markdown)
	result := h.lastResult
	cached := result != nil && h.lastSource == cacheKey
	if !cached {
		var err error
		result, err = h.renderer.Render(markdown)
		if err != nil {
			return nil, err
		}
		for _, asset := range result.Assets {
			if !strings.HasPrefix(asset.ID, "generated-") && asset.Content != "" {
				h.assets[asset.ID] = []byte(asset.Content)
			}
		}
		h.lastSource, h.lastResult = cacheKey, result
	}
	previewResult := *result
	previewResult.HTML = frontmatterPreview + result.HTML
	fontpack := req.Appearance[workbenchFontpackAttribute]
	if fontpack == "" {
		fontpack = h.appearance.Defaults[workbenchFontpackAttribute]
	}
	response := &renderResponse{SchemaVersion: 1, RendererFingerprint: h.previewFingerprint, HTML: previewResult.HTML, SourceLines: h.renderer.SourceLines(markdown), SourceAnchors: sourceAnchors, FeatureHeadHTML: h.featureHead(result, fontpack), BodyEndHTML: result.BodyEndHTML, InputBytes: len(*req.Markdown), HTMLBytes: len(previewResult.HTML), Cached: cached}
	if req.IncludeDocument == nil || *req.IncludeDocument {
		document := h.document(&previewResult, req.Appearance)
		response.Document, response.DocumentBytes = &document, len(document)
	}
	return response, nil
}

func frontmatterSource(req *renderRequest) string {
	if req == nil || req.FrontmatterSource == nil {
		return ""
	}
	return *req.FrontmatterSource
}

type renderRequest struct {
	Markdown          *string           `json:"markdown"`
	FrontmatterSource *string           `json:"frontmatter_source,omitempty"`
	Appearance        map[string]string `json:"appearance"`
	IncludeDocument   *bool             `json:"include_document"`
}

// renderResponse matches the Python reference's full-document/fragment contract.
type renderResponse struct {
	SchemaVersion       int     `json:"schema_version"`
	RendererFingerprint string  `json:"renderer_fingerprint"`
	Document            *string `json:"document"`
	HTML                string  `json:"html"`
	SourceLines         []int   `json:"source_lines"`
	SourceAnchors       []int   `json:"source_anchors"`
	FeatureHeadHTML     string  `json:"feature_head_html"`
	BodyEndHTML         string  `json:"body_end_html"`
	RenderMS            float64 `json:"render_ms"`
	QueueMS             float64 `json:"queue_ms"`
	InputBytes          int     `json:"input_bytes"`
	HTMLBytes           int     `json:"html_bytes"`
	DocumentBytes       int     `json:"document_bytes"`
	Cached              bool    `json:"cached"`
}

// featureHead omits styles already owned by the persistent appearance shell.
// Mark feature styles so subsequent fragments can remove unused styles as well.
func (h *Handler) featureHead(result *plugins.StandaloneResult, fontpack string) string {
	head := result.HeadHTML
	for _, name := range h.appearance.Styles {
		head = strings.ReplaceAll(head, styleLink(name)+"\n", "")
	}
	head = h.inlinePreviewStylesheetMarkup(head, fontpack)
	return strings.NewReplacer("<link ", `<link data-markata-render-head="true" `, "<style>", `<style data-markata-render-head="true">`).Replace(head)
}

func (h *Handler) inlinePreviewStylesheetMarkup(markup, fontpack string) string {
	return stylesheetLinkPattern.ReplaceAllStringFunc(markup, func(tag string) string {
		match := stylesheetHrefPattern.FindStringSubmatch(tag)
		if len(match) != 3 {
			return tag
		}
		rawURL := match[1]
		if rawURL == "" {
			rawURL = match[2]
		}
		parsed, err := url.Parse(rawURL)
		if err != nil || parsed.IsAbs() || parsed.Host != "" || !strings.HasPrefix(parsed.Path, "/appearance/") {
			return tag
		}
		name := strings.TrimPrefix(parsed.Path, "/appearance/")
		data, ok := h.assets[name]
		if !ok {
			return tag
		}
		css := h.inlinePreviewCSS(name, data, fontpack)
		dataURL := "data:text/css;base64," + base64.StdEncoding.EncodeToString(css)
		quote := `"`
		if strings.Contains(match[0], "'") {
			quote = `'`
		}
		return strings.Replace(tag, match[0], "href="+quote+dataURL+quote, 1)
	})
}

func (h *Handler) inlinePreviewCSS(name string, css []byte, fontpack string) []byte {
	if name != "css/fonts.css" || h.appearance == nil {
		return css
	}
	fontSources := make(map[string]bool)
	for _, source := range h.appearance.FontPackSources[fontpack] {
		fontSources[source] = true
	}
	if len(fontSources) == 0 {
		return css
	}
	result := append([]byte(nil), css...)
	for assetURL, source := range h.appearance.FontAssetSources {
		if !fontSources[source] {
			continue
		}
		asset, ok := h.assets[strings.TrimPrefix(assetURL, "/")]
		if !ok {
			continue
		}
		previewURL := "/appearance" + assetURL
		dataURL := "data:font/woff2;base64," + base64.StdEncoding.EncodeToString(asset)
		result = bytes.ReplaceAll(result, []byte(previewURL), []byte(dataURL))
	}
	return result
}

func (h *Handler) previewFontCSS(fontpack string) ([]byte, bool) {
	if h.appearance == nil {
		return nil, false
	}
	if fontpack == "" {
		fontpack = h.appearance.Defaults[workbenchFontpackAttribute]
	}
	if _, ok := h.appearance.FontPackSources[fontpack]; !ok {
		return nil, false
	}
	css, ok := h.assets["css/fonts.css"]
	if !ok {
		return nil, false
	}
	return h.inlinePreviewCSS("css/fonts.css", css, fontpack), true
}

func (h *Handler) servePreviewFonts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	css, ok := h.previewFontCSS(r.URL.Query().Get("fontpack"))
	if !ok {
		http.Error(w, "unknown font pack", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "private, no-store")
	if r.Method == http.MethodHead {
		return
	}
	//nolint:gosec // G705: this response is text/css from local appearance assets, not HTML.
	if _, err := w.Write(css); err != nil {
		return
	}
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
	if req.FrontmatterSource != nil && len(*req.FrontmatterSource) > maxSource {
		fail(w, 413, "frontmatter exceeds 200000 UTF-8 bytes")
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
	defaults := map[string]string{workbenchThemeAttribute: workbenchDarkMode, workbenchPaletteAttribute: "default-dark", workbenchAestheticAttribute: "minimal", workbenchFontpackAttribute: "system", workbenchTextSizeAttribute: "large"}
	for key, value := range h.appearance.Defaults {
		defaults[key] = value
	}
	for key, value := range attributes {
		defaults[key] = value
	}
	if _, ok := attributes[workbenchPaletteAttribute]; !ok && defaults[workbenchThemeAttribute] == workbenchLightMode {
		if h.siteConfig == nil {
			defaults[workbenchPaletteAttribute] = "default-light"
		}
	}
	attrs := ""
	for _, key := range []string{workbenchThemeAttribute, workbenchPaletteAttribute, workbenchAestheticAttribute, workbenchFontpackAttribute, workbenchTextSizeAttribute} {
		attrs += " " + key + `="` + html.EscapeString(defaults[key]) + `"`
	}
	styles := ""
	for _, name := range h.appearance.Styles {
		data, ok := h.assets[name]
		if !ok {
			styles += h.versionedStyleLink(name)
			continue
		}
		css := h.inlinePreviewCSS(name, data, defaults[workbenchFontpackAttribute])
		if name == "css/fonts.css" {
			encoded, err := json.Marshal(string(css))
			if err == nil {
				styles += `<script>(function(){const style=document.createElement('style');style.id="markata-preview-fonts";style.textContent=` + string(encoded) + `;document.head.append(style)})()</script>`
			}
			continue
		}
		styles += dataStylesheetLink(css)
	}
	bg, fg := "#131721", "#e5e7eb"
	if defaults[workbenchThemeAttribute] == workbenchLightMode {
		bg, fg = "#fff", "#1f2937"
	}
	initial := `<style>html{color-scheme:` + defaults[workbenchThemeAttribute] + `;background:` + bg + `}body{background:var(--color-background,` + bg + `);color:var(--color-text,` + fg + `)}</style>`
	return `<!doctype html><html lang="en"` + attrs + `><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Markdown preview</title>` + initial + styles + h.featureHead(result, defaults[workbenchFontpackAttribute]) + `</head><body><main class="main-content"><div class="content-wrapper"><div class="post-shell"><article class="post h-entry"><div class="post-content e-content">` + result.HTML + `</div></article></div></div></main><script>window.markataEditorPreview=true</script>` + h.bridge + strings.ReplaceAll(result.BodyEndHTML, "<script", `<script data-markata-generated="true"`) + `<script>parent.postMessage({type:'preview-features-ready'},'*')</script></body></html>`
}

func dataStylesheetLink(css []byte) string {
	return `<link rel="stylesheet" href="data:text/css;base64,` + base64.StdEncoding.EncodeToString(css) + `">`
}

// serveProseIgnoreRule is the only persistence bridge for web prose ignores.
// It never accepts a filesystem path or arbitrary configuration mutation.
func (h *Handler) serveProseIgnoreRule(w http.ResponseWriter, r *http.Request) {
	if !h.sameOriginWorkspaceMutation(r) {
		fail(w, http.StatusForbidden, "same-origin workspace writes are required")
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		fail(w, http.StatusUnsupportedMediaType, "send Content-Type: application/json")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	var req struct {
		Rule   string `json:"rule"`
		Phrase string `json:"phrase,omitempty"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		fail(w, http.StatusBadRequest, "invalid prose ignore request")
		return
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		fail(w, http.StatusBadRequest, "expected one prose ignore request")
		return
	}
	if !prose.KnownRule(req.Rule) || (req.Phrase != "" && !prose.ValidPhrase(req.Rule, req.Phrase)) {
		fail(w, http.StatusBadRequest, "unknown prose rule or invalid phrase")
		return
	}
	h.workspaceWriteMu.Lock()
	defer h.workspaceWriteMu.Unlock()
	// Site root is server-owned; clients cannot pick paths outside it.
	var err error
	if req.Phrase == "" {
		_, err = prose.IgnoreRuleAtRoot(h.root, req.Rule)
	} else {
		_, err = prose.IgnorePhraseAtRoot(h.root, req.Rule, req.Phrase)
	}
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "could not save prose preference")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if _, err := io.WriteString(w, "{\"ok\":true}"); err != nil {
		return
	}
}

func isWorkspaceRoute(routePath, method string) bool {
	switch {
	case routePath == "/workspace" || routePath == "/workspace/templates" || routePath == "/workspace/content-number" || routePath == "/workspace/image" || routePath == "/workspace/index-status" || routePath == "/workspace/changes":
		return method == http.MethodGet
	case routePath == "/workspace/new-post" || routePath == "/workspace/dropper-upload":
		return method == http.MethodPost
	case routePath == "/workspace/link-metadata" || routePath == "/workspace/prose/ignore-rule":
		return method == http.MethodPost
	case routePath == "/workspace/template-settings":
		return method == http.MethodPost
	case routePath == "/workspace/file":
		return method == http.MethodGet || method == http.MethodPut || method == http.MethodPost
	case routePath == "/workspace/publish/commit" || routePath == "/workspace/publish/push" || routePath == "/workspace/publish/discard":
		return method == http.MethodPost
	case routePath == "/workspace/publish/deployment" || routePath == "/workspace/publish/review" || routePath == "/workspace/publish/readiness":
		return method == http.MethodGet
	case routePath == frontmatterInspectPath || routePath == frontmatterEditPath:
		return method == http.MethodPost
	default:
		return false
	}
}

func (h *Handler) imageUploadProvider() string {
	mode := imageUploadModeAuto
	if config := h.currentSiteConfig(); config != nil && config.ContentTemplates.ImageUpload != "" {
		mode = config.ContentTemplates.ImageUpload
	}
	switch mode {
	case imageUploadModeLocal:
		return imageUploadModeLocal
	case imageUploadModeDropper:
		return imageUploadModeDropper
	case imageUploadModeAuto:
		if h.dropperUploader != nil || strings.TrimSpace(os.Getenv("DROPPER_TOKEN")) != "" {
			return imageUploadModeDropper
		}
		return imageUploadModeLocal
	default:
		return "invalid"
	}
}

func (h *Handler) dropperUploadEndpoint() string {
	if configured := strings.TrimSpace(os.Getenv("DROPPER_API_URL")); configured != "" {
		return configured
	}
	if config := h.currentSiteConfig(); config != nil && strings.TrimSpace(config.ContentTemplates.DropperUploadURL) != "" {
		return strings.TrimSpace(config.ContentTemplates.DropperUploadURL)
	}
	return "https://dropper.waylonwalker.com/api/upload?format=json"
}

func (h *Handler) currentSiteConfig() *models.Config {
	if h.store != nil {
		if config, _ := h.store.Configuration(); config != nil {
			return config
		}
	}
	return h.siteConfig
}

func (h *Handler) uploadImageToDropper(ctx context.Context, filename, contentType string, data []byte) (string, error) {
	if h.dropperUploader != nil {
		return h.dropperUploader.Upload(ctx, filename, contentType, data)
	}
	token := strings.TrimSpace(os.Getenv("DROPPER_TOKEN"))
	if token == "" {
		return "", fmt.Errorf("dropper uploads require the server-side DROPPER_TOKEN")
	}
	return uploadImageToDropper(ctx, h.dropperUploadEndpoint(), token, filename, contentType, data)
}

func validateDropperUploadEndpoint(endpoint string) error {
	parsedEndpoint, err := url.ParseRequestURI(endpoint)
	if err != nil || parsedEndpoint.Host == "" || parsedEndpoint.User != nil || parsedEndpoint.Scheme != "https" && !(parsedEndpoint.Scheme == "http" && isLoopbackHost(parsedEndpoint.Hostname())) {
		return fmt.Errorf("dropper upload URL must use HTTPS or a loopback HTTP URL")
	}
	return nil
}

func dropperUploadPayload(filename, contentType string, data []byte) (bytes.Buffer, string, error) {
	var payload bytes.Buffer
	form := multipart.NewWriter(&payload)
	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename=%q`, filename))
	partHeader.Set("Content-Type", contentType)
	part, err := form.CreatePart(partHeader)
	if err != nil {
		return payload, "", fmt.Errorf("prepare dropper upload: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return payload, "", fmt.Errorf("prepare dropper image: %w", err)
	}
	if err := form.Close(); err != nil {
		return payload, "", fmt.Errorf("finish dropper upload: %w", err)
	}
	return payload, form.FormDataContentType(), nil
}

func dropperUploadResponseURL(statusCode int, status string, responseBody []byte) (string, error) {
	if statusCode < 200 || statusCode >= 300 {
		return "", fmt.Errorf("dropper upload failed (%s)", status)
	}
	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", fmt.Errorf("dropper returned invalid JSON: %w", err)
	}
	return validatedDropperImageURL(result.URL)
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func (h *Handler) serveWorkspace(w http.ResponseWriter, r *http.Request) {
	if h.hostedAuth == nil {
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
	}
	if r.URL.Path == "/workspace/image" {
		name := r.URL.Query().Get("path")
		data, err := h.workspace.ReadImage(name)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		serve(w, r, name, data)
		return
	}
	if r.URL.Path == "/workspace/dropper-upload" {
		h.serveWorkspaceDropperUpload(w, r)
		return
	}
	if r.URL.Path == "/workspace/prose/ignore-rule" {
		h.serveProseIgnoreRule(w, r)
		return
	}
	if r.URL.Path == "/workspace/new-post" {
		if !h.sameOriginWorkspaceMutation(r) {
			fail(w, http.StatusForbidden, "same-origin workspace writes are required")
			return
		}
		const maxImageBytes = 15 << 20
		r.Body = http.MaxBytesReader(w, r.Body, maxSource+(48<<20)+(1<<20))
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			fail(w, http.StatusBadRequest, "cannot read new post upload")
			return
		}
		if r.MultipartForm != nil {
			defer func() {
				if err := r.MultipartForm.RemoveAll(); err != nil {
					log.Printf("cleanup multipart image upload: %v", err)
				}
			}()
		}
		source := r.FormValue("source")
		if len(source) > maxSource {
			fail(w, http.StatusRequestEntityTooLarge, "document exceeds 200,000 bytes")
			return
		}
		headers := r.MultipartForm.File["image"]
		if len(headers) == 0 || len(headers) > 8 || len(r.MultipartForm.File) != 1 {
			fail(w, http.StatusBadRequest, "provide 1 to 8 photos under the image field")
			return
		}
		uploads := make([]editorworkspace.UploadedImage, 0, len(headers))
		total := 0
		for _, header := range headers {
			image, err := header.Open()
			if err != nil {
				fail(w, http.StatusBadRequest, "cannot open image")
				return
			}
			data, readErr := io.ReadAll(io.LimitReader(image, maxImageBytes+1))
			closeErr := image.Close()
			if readErr != nil || closeErr != nil {
				fail(w, http.StatusBadRequest, "cannot read image")
				return
			}
			total += len(data)
			if len(data) == 0 || len(data) > maxImageBytes || total > 48<<20 {
				fail(w, http.StatusRequestEntityTooLarge, "photos must be 15 MiB each and 48 MiB total")
				return
			}
			if _, err := editorworkspace.ImageExtension(header.Header.Get("Content-Type"), data); err != nil {
				fail(w, http.StatusUnprocessableEntity, err.Error())
				return
			}
			uploads = append(uploads, editorworkspace.UploadedImage{
				ContentType: header.Header.Get("Content-Type"), Data: data,
			})
		}
		h.workspaceWriteMu.Lock()
		defer h.workspaceWriteMu.Unlock()
		postPath := r.FormValue("path")
		imageName := r.FormValue("image_name")
		provider := h.imageUploadProvider()
		includeGallery := r.FormValue("include_gallery") == pwaGalleryFlagTrue
		var images []string
		var createErr error
		storedSource := source
		switch provider {
		case imageUploadModeDropper:
			// Dropper has no transactional delete API. Do not leak previous
			// uploads when a later photo fails or the post is rejected.
			if len(uploads) > 1 {
				fail(w, http.StatusConflict, "multi-photo Dropper uploads require cleanup support; use local image storage")
				return
			}
			checker, ok := h.workspace.(workspaceCreateChecker)
			if !ok {
				fail(w, http.StatusNotImplemented, "workspace cannot validate a new post path")
				return
			}
			if !validNumberPrefix(imageName) || len(imageName) > 100 {
				fail(w, http.StatusBadRequest, "image name must be a lowercase slug")
				return
			}
			if checkErr := checker.CheckCreate(postPath); checkErr != nil {
				fail(w, http.StatusUnprocessableEntity, checkErr.Error())
				return
			}
			for i, photo := range uploads {
				extension, extensionErr := editorworkspace.ImageExtension(photo.ContentType, photo.Data)
				if extensionErr != nil {
					createErr = extensionErr
					break
				}
				name := imageName
				if len(uploads) > 1 {
					name = fmt.Sprintf("%s-%d", imageName, i+1)
				}
				imageURL, err := h.uploadImageToDropper(r.Context(), name+"."+extension, photo.ContentType, photo.Data)
				if err != nil {
					createErr = err
					break
				}
				images = append(images, imageURL)
			}
			if createErr == nil {
				if len(images) > 1 || includeGallery {
					storedSource = editorworkspace.SourceWithUploadedImages(source, images)
				} else {
					storedSource = sourceWithImage(source, images[0])
				}
				createErr = h.workspace.Create(postPath, storedSource)
			}
		case imageUploadModeLocal:
			if len(uploads) == 1 {
				var uploadedURL string
				singleSource := source
				if includeGallery {
					extension, extensionErr := editorworkspace.ImageExtension(uploads[0].ContentType, uploads[0].Data)
					if extensionErr != nil {
						fail(w, http.StatusUnprocessableEntity, extensionErr.Error())
						return
					}
					localURL := "/images/" + imageName + "." + extension
					singleSource = editorworkspace.SourceWithUploadedImages(source, []string{localURL})
				}
				uploadedURL, createErr = h.workspace.CreateWithImage(postPath, singleSource, imageName, uploads[0].ContentType, uploads[0].Data)
				if createErr == nil {
					images = []string{uploadedURL}
					storedSource = sourceWithImage(singleSource, uploadedURL)
				}
			} else {
				creator, ok := h.workspace.(workspaceMultiImageCreator)
				if !ok {
					fail(w, http.StatusNotImplemented, "workspace cannot create a multi-image post")
					return
				}
				images, createErr = creator.CreateWithImages(postPath, source, imageName, uploads)
				if createErr == nil {
					storedSource = editorworkspace.SourceWithUploadedImages(source, images)
				}
			}
		default:
			createErr = fmt.Errorf("unsupported image upload provider %q", provider)
		}
		if createErr != nil {
			fail(w, http.StatusUnprocessableEntity, createErr.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(map[string]any{
			workspaceVersionKey: editorworkspace.SourceVersion(storedSource),
			"image":             images[0],
			pwaSharePhotosField: images,
			"provider":          provider,
		}); err != nil {
			return
		}
		return
	}
	if r.URL.Path == "/workspace/link-metadata" {
		h.serveLinkMetadata(w, r)
		return
	}
	if r.URL.Path == "/workspace/template-settings" {
		if !h.sameOriginWorkspaceMutation(r) {
			fail(w, http.StatusForbidden, "same-origin workspace writes are required")
			return
		}
		writer, ok := h.workspace.(contentTemplateDefaultsWriter)
		if !ok {
			fail(w, http.StatusNotImplemented, "site content defaults cannot be saved by this workspace")
			return
		}
		var request struct {
			ID           string `json:"id"`
			Directory    string `json:"directory"`
			FilenameMode string `json:"filename_mode"`
		}
		decoder := json.NewDecoder(io.LimitReader(r.Body, 32<<10))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&request); err != nil {
			fail(w, http.StatusBadRequest, "invalid site content defaults")
			return
		}
		if !h.acquire(r) {
			return
		}
		refreshErr := h.refreshPresentation()
		templates, err := h.contentTemplates()
		<-h.gate
		if refreshErr != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot refresh site configuration: "+refreshErr.Error())
			return
		}
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot load content templates: "+err.Error())
			return
		}
		var key string
		for i := range templates {
			if templates[i].ID == request.ID {
				key = contentTemplateKey(templates[i])
				break
			}
		}
		if key == "" || !validNumberPrefix(key) || !safeStarterOutputDirectory(request.Directory) || request.FilenameMode != contentFilenameModeTitle && request.FilenameMode != contentFilenameModeNumbered {
			fail(w, http.StatusBadRequest, "invalid content type, folder, or filename mode")
			return
		}
		h.workspaceWriteMu.Lock()
		err = writer.SetContentTemplateDefaults(key, request.Directory, request.FilenameMode)
		h.workspaceWriteMu.Unlock()
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot save site content defaults: "+err.Error())
			return
		}
		if !h.acquire(r) {
			return
		}
		err = h.refreshPresentation()
		<-h.gate
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot reload site content defaults: "+err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]bool{"saved": true}); err != nil {
			return
		}
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Path == "/workspace/publish/commit" {
		h.servePublishCommit(w, r)
		return
	}
	if r.URL.Path == "/workspace/publish/discard" {
		h.servePublishDiscard(w, r)
		return
	}
	if r.URL.Path == "/workspace/publish/push" {
		h.servePublishPush(w, r)
		return
	}
	if r.URL.Path == "/workspace/publish/deployment" {
		h.servePublishDeployment(w, r)
		return
	}
	if r.URL.Path == "/workspace/frontmatter/inspect" || r.URL.Path == "/workspace/frontmatter/edit" {
		h.workspaceWriteMu.Lock()
		defer h.workspaceWriteMu.Unlock()
		h.serveFrontmatterEditor(w, r)
		return
	}
	if r.URL.Path == "/workspace/publish/review" {
		h.servePublishReview(w, r)
		return
	}
	if r.URL.Path == "/workspace/publish/readiness" {
		h.servePublishReadiness(w, r)
		return
	}
	if r.URL.Path == "/workspace/changes" {
		h.serveWorkspaceChanges(w, r)
		return
	}
	if r.URL.Path == "/workspace/index-status" {
		if h.store == nil {
			fail(w, http.StatusServiceUnavailable, "site index telemetry is unavailable")
			return
		}
		if err := json.NewEncoder(w).Encode(h.store.IndexStatus()); err != nil {
			return
		}
		return
	}
	if r.URL.Path == "/workspace/templates" {
		if !h.acquire(r) {
			return
		}
		refreshErr := h.refreshPresentation()
		templates, err := h.contentTemplates()
		<-h.gate
		if refreshErr != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot refresh site configuration: "+refreshErr.Error())
			return
		}
		if err != nil {
			fail(w, 422, "cannot load content templates: "+err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]any{"templates": templates, "image_upload": h.imageUploadProvider()}); err != nil {
			return
		}
		return
	}
	if r.URL.Path == "/workspace/content-number" {
		directory := r.URL.Query().Get("directory")
		prefix := r.URL.Query().Get("prefix")
		if !safeStarterOutputDirectory(directory) || !validNumberPrefix(prefix) {
			fail(w, http.StatusBadRequest, "invalid content directory or number prefix")
			return
		}
		next, err := nextContentNumber(h.root, directory, prefix)
		if err != nil {
			fail(w, http.StatusUnprocessableEntity, "cannot find next content number: "+err.Error())
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]int{"next_number": next}); err != nil {
			return
		}
		return
	}
	if r.URL.Path == "/workspace/file" {
		if r.Method == http.MethodPut || r.Method == http.MethodPost {
			if !h.sameOriginWorkspaceMutation(r) {
				fail(w, http.StatusForbidden, "same-origin workspace writes are required")
				return
			}
			// Serialize local editor writes with version-checked Discard restores.
			h.workspaceWriteMu.Lock()
			defer h.workspaceWriteMu.Unlock()
			body, err := io.ReadAll(io.LimitReader(r.Body, maxSource+1))
			if err != nil {
				fail(w, 400, "cannot read source")
				return
			}
			if len(body) > maxSource {
				fail(w, 413, "document exceeds 200,000 bytes")
				return
			}
			sourcePath := r.URL.Query().Get("path")
			if r.Method == http.MethodPut {
				// Browser autosave is compare-and-swap against the exact source bytes
				// it last loaded or successfully wrote. Reject stale saves rather
				// than overwriting a newer file from another editor process.
				expected := r.Header.Get("X-Markata-Expected-Version")
				if !gitHashPattern.MatchString(expected) || len(expected) != 64 {
					fail(w, http.StatusPreconditionRequired, "expected source version required; reload before saving")
					return
				}
				current, readErr := h.workspace.Read(sourcePath)
				if readErr != nil {
					fail(w, http.StatusUnprocessableEntity, "cannot read workspace source before saving")
					return
				}
				currentVersion := editorworkspace.SourceVersion(current)
				if currentVersion != expected {
					// A duplicate save may arrive after another tab already wrote the
					// exact same bytes. Acknowledge that idempotently without touching
					// the file again; a different current source remains a conflict.
					if current != string(body) {
						fail(w, http.StatusConflict, "file changed on disk; browser recovery copy kept; compare copies to resolve conflict")
						return
					}
					w.Header().Set("Content-Type", "application/json")
					if err := json.NewEncoder(w).Encode(map[string]string{workspaceVersionKey: currentVersion}); err != nil {
						return
					}
					return
				}
			}
			if r.Method == http.MethodPost {
				err = h.workspace.Create(sourcePath, string(body))
			} else {
				err = h.workspace.Write(sourcePath, string(body))
			}
			if err != nil {
				fail(w, 422, err.Error())
				return
			}
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
			}
			if err := json.NewEncoder(w).Encode(map[string]string{"version": editorworkspace.SourceVersion(string(body))}); err != nil {
				return
			}
			return
		}
		source, err := h.workspace.Read(r.URL.Query().Get("path"))
		if err != nil {
			fail(w, 404, err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(map[string]string{"source": source, "version": editorworkspace.SourceVersion(source)}); err != nil {
			return
		}
		return
	}
	var files []editorworkspace.File
	var err error
	listStarted := time.Now()
	if r.URL.Query().Get("metadata") == "true" {
		files, err = h.workspace.Files()
	} else {
		files, err = h.workspace.Paths()
	}
	if err != nil {
		fail(w, 422, err.Error())
		return
	}
	// A path/metadata request may spend additional time assembling configured
	// feeds after reading the cached site index. Report both costs separately.
	indexMS := float64(time.Since(listStarted).Microseconds()) / 1000
	response := map[string]any{"files": files, "rootUri": h.languageConfig.RootURI}
	if r.URL.Query().Get("metadata") == "true" {
		feedsStarted := time.Now()
		feeds := h.workspaceFeeds()
		feedsDuration := time.Since(feedsStarted)
		response["feeds"] = feeds
		w.Header().Set("Server-Timing", fmt.Sprintf("site-index;dur=%.1f, site-feeds;dur=%.1f", indexMS, float64(feedsDuration.Microseconds())/1000))
		if feedsDuration >= 200*time.Millisecond {
			log.Printf("editor site feeds: %d feeds assembled in %s", len(feeds), feedsDuration.Round(time.Millisecond))
		}
	} else {
		w.Header().Set("Server-Timing", fmt.Sprintf("site-paths;dur=%.1f", indexMS))
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *Handler) serveWorkspaceDropperUpload(w http.ResponseWriter, r *http.Request) {
	if !h.sameOriginWorkspaceMutation(r) {
		fail(w, http.StatusForbidden, "same-origin workspace writes are required")
		return
	}
	if h.imageUploadProvider() != imageUploadModeDropper {
		fail(w, http.StatusConflict, "Dropper image uploads are not enabled for this site")
		return
	}
	if h.dropperUploader == nil && strings.TrimSpace(os.Getenv("DROPPER_TOKEN")) == "" {
		fail(w, http.StatusServiceUnavailable, "Dropper uploads are not configured on this server")
		return
	}

	const maxImageBytes = 15 << 20
	r.Body = http.MaxBytesReader(w, r.Body, maxImageBytes+(1<<20))
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			fail(w, http.StatusRequestEntityTooLarge, "image exceeds 15 MiB")
			return
		}
		fail(w, http.StatusBadRequest, "cannot read image upload")
		return
	}
	if r.MultipartForm != nil {
		defer func() {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				log.Printf("cleanup multipart Dropper upload: %v", err)
			}
		}()
	}
	files := r.MultipartForm.File["image"]
	if len(r.MultipartForm.File) != 1 || len(files) != 1 {
		fail(w, http.StatusBadRequest, "provide exactly one image file")
		return
	}
	image, header, err := r.FormFile("image")
	if err != nil {
		fail(w, http.StatusBadRequest, "image upload is missing")
		return
	}
	defer image.Close()
	data, err := io.ReadAll(io.LimitReader(image, maxImageBytes+1))
	if err != nil {
		fail(w, http.StatusBadRequest, "cannot read image upload")
		return
	}
	if len(data) > maxImageBytes {
		fail(w, http.StatusRequestEntityTooLarge, "image exceeds 15 MiB")
		return
	}
	extension, err := editorworkspace.ImageExtension(header.Header.Get("Content-Type"), data)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	filename := fmt.Sprintf("markata-image-%d.%s", time.Now().UTC().UnixNano(), extension)
	imageURL, err := h.uploadImageToDropper(r.Context(), filename, header.Header.Get("Content-Type"), data)
	if err != nil {
		fail(w, http.StatusBadGateway, err.Error())
		return
	}
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"url": imageURL, "provider": imageUploadModeDropper}); err != nil {
		return
	}
}

type linkMetadataRequest struct {
	URL string `json:"url"`
}

type linkMetadataResponse struct {
	Title        string `json:"title,omitempty"`
	Description  string `json:"description,omitempty"`
	Image        string `json:"image,omitempty"`
	SiteName     string `json:"site_name,omitempty"`
	Type         string `json:"type,omitempty"`
	ProviderName string `json:"provider_name,omitempty"`
	AuthorName   string `json:"author_name,omitempty"`
	AuthorURL    string `json:"author_url,omitempty"`
	Source       string `json:"source,omitempty"`
	Available    bool   `json:"available"`
}

func (h *Handler) serveLinkMetadata(w http.ResponseWriter, r *http.Request) {
	if !h.sameOriginWorkspaceMutation(r) {
		fail(w, http.StatusForbidden, "same-origin workspace requests are required")
		return
	}
	request, err := decodeLinkMetadataRequest(w, r)
	if err != nil {
		fail(w, http.StatusBadRequest, err.Error())
		return
	}
	destination := strings.TrimSpace(request.URL)
	if !validLinkMetadataURL(destination) {
		fail(w, http.StatusBadRequest, "enter an absolute HTTP(S) URL without credentials")
		return
	}
	if !h.acquire(r) {
		return
	}
	refreshErr := h.refreshPresentation()
	siteConfig := h.siteConfig
	<-h.gate
	if refreshErr != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot refresh site configuration: "+refreshErr.Error())
		return
	}
	resolver, err := h.newEmbedMetadataResolver(siteConfig)
	if err != nil {
		fail(w, http.StatusUnprocessableEntity, "cannot configure embed metadata: "+err.Error())
		return
	}
	resolver.SetMetadataHTTPClient(newPublicMetadataHTTPClient(r.Context(), resolver.Config().Timeout))
	metadata, err := resolver.ResolveExternalMetadata(destination)
	if err != nil {
		fail(w, http.StatusBadRequest, err.Error())
		return
	}
	title := strings.TrimSpace(metadata.Title)
	if title == strings.TrimSpace(resolver.Config().FallbackTitle) {
		title = ""
	}
	response := linkMetadataResponse{
		Title:        title,
		Description:  metadata.Description,
		Image:        absoluteMetadataURL(metadata.Image, destination),
		SiteName:     metadata.SiteName,
		Type:         metadata.Type,
		ProviderName: metadata.ProviderName,
		AuthorName:   metadata.AuthorName,
		AuthorURL:    absoluteMetadataURL(metadata.AuthorURL, destination),
		Source:       metadata.Source,
	}
	response.Available = linkMetadataIsAvailable(response)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func decodeLinkMetadataRequest(w http.ResponseWriter, r *http.Request) (linkMetadataRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request linkMetadataRequest
	if err := decoder.Decode(&request); err != nil {
		return linkMetadataRequest{}, errors.New("expected a JSON object with a link URL")
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return linkMetadataRequest{}, errors.New("expected one link metadata request")
	}
	return request, nil
}

func validLinkMetadataURL(destination string) bool {
	if len(destination) > 8192 {
		return false
	}
	parsed, err := url.Parse(destination)
	if err != nil || parsed.Hostname() == "" || parsed.User != nil {
		return false
	}
	return publicMetadataURL(parsed)
}

func linkMetadataIsAvailable(response linkMetadataResponse) bool {
	return response.Title != "" || response.Description != "" || response.Image != "" || response.SiteName != "" || response.Type != "" || response.ProviderName != "" || response.AuthorName != "" || response.AuthorURL != ""
}

func absoluteMetadataURL(rawURL, baseURL string) string {
	if strings.TrimSpace(rawURL) == "" {
		return ""
	}
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return ""
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	resolved := base.ResolveReference(parsed)
	if (resolved.Scheme != "http" && resolved.Scheme != "https") || resolved.Hostname() == "" || resolved.User != nil {
		return ""
	}
	return resolved.String()
}

func (h *Handler) newEmbedMetadataResolver(siteConfig *models.Config) (*plugins.EmbedsPlugin, error) {
	resolver := plugins.NewEmbedsPlugin()
	manager := lifecycle.NewManager()
	config := lifecycle.NewConfig()
	config.ContentDir = h.root
	if siteConfig != nil {
		config.Extra = siteConfig.Extra
	}
	manager.SetConfig(config)
	if err := resolver.Configure(manager); err != nil {
		return nil, fmt.Errorf("configure embeds plugin: %w", err)
	}
	embedsConfig := resolver.Config()
	if embedsConfig.CacheDir != "" && !filepath.IsAbs(embedsConfig.CacheDir) {
		root, err := filepath.Abs(h.root)
		if err != nil {
			return nil, fmt.Errorf("resolve workspace path: %w", err)
		}
		embedsConfig.CacheDir = filepath.Join(root, embedsConfig.CacheDir)
		resolver.SetConfig(embedsConfig)
	}
	return resolver, nil
}

func sameOriginWorkspaceMutation(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	u, err := url.Parse(origin)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if r.URL.Scheme != "" {
		scheme = r.URL.Scheme
	}
	if err != nil || u.User != nil || u.Host == "" || u.Host != r.Host || u.Path != "" || u.RawPath != "" || u.RawQuery != "" || u.Fragment != "" || u.Scheme != scheme {
		return false
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" && site != "same-origin" {
		return false
	}
	return true
}

func (h *Handler) sameOriginWorkspaceMutation(r *http.Request) bool {
	if h.hostedAuth == nil {
		return sameOriginWorkspaceMutation(r)
	}
	origins := r.Header.Values("Origin")
	if len(origins) != 1 || origins[0] != h.hostedAuth.origin {
		return false
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" && site != "same-origin" {
		return false
	}
	return true
}

func (h *Handler) validLSPOrigin(r *http.Request) bool {
	if h.hostedAuth == nil {
		return lsp.ValidBrowserOrigin(r)
	}
	return h.authorizeHostedRequest(r) && h.sameOriginWorkspaceMutation(r)
}

func (h *Handler) acquire(r *http.Request) bool {
	select {
	case h.gate <- struct{}{}:
		return true
	case <-r.Context().Done():
		return false
	}
}

func (h *Handler) refreshPresentation() error {
	if h.store == nil {
		return nil
	}
	cfg, revision := h.store.Configuration()
	if revision == h.configRevision {
		return nil
	}
	next, err := buildPresentation(h.root, h.workspace, h.languageConfig, h.language, cfg)
	if err != nil {
		return err
	}
	h.siteConfig = next.siteConfig
	h.appearance = next.appearance
	h.choices = next.choices
	h.assets = next.assets
	h.index = next.index
	h.bridge = next.bridge
	h.renderer = next.renderer
	h.wasmVersion = next.wasmVersion
	h.previewFingerprint = next.previewFingerprint
	h.lastResult = nil
	h.lastSource = ""
	h.configRevision = revision
	return nil
}
