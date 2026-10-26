package servecontrol

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/WaylonWalker/markata/pkg/servefix"
)

// NewWebHandler serves a local view of the same runtime consumed by the TUI.
// Callers must bind the HTTP server to a loopback address for local use.
func NewWebHandler(runtime *Runtime) http.Handler {
	return NewWebHandlerWithSourceRoot(runtime, "")
}

// NewWebHandlerWithSourceRoot adds previewable, stale-safe source fixes when
// a site root is available.
func NewWebHandlerWithSourceRoot(runtime *Runtime, sourceRoot string) http.Handler {
	return NewWebHandlerWithFixHooks(runtime, sourceRoot, nil)
}

// NewWebHandlerWithFixHooks adds local source fixes and optional hooks used by
// Serve to suppress only filesystem events caused by successfully previewed
// replacements. The hook receives root-relative paths and planned after text.
func NewWebHandlerWithFixHooks(runtime *Runtime, sourceRoot string, beforeApply func([]servefix.FilePreview)) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/_markata/assets/pierre-diffs/", pierreDiffAssetsHandler())
	mux.HandleFunc("/_markata/", dashboardHandler(runtime))
	mux.HandleFunc("/_markata/api/state", stateHandler(runtime))
	mux.HandleFunc("/_markata/api/actions", actionsHandler(runtime))
	mux.HandleFunc("/_markata/api/fixes/batch/preview", batchFixPreviewHandler(sourceRoot))
	mux.HandleFunc("/_markata/api/fixes/batch/apply", batchFixApplyHandler(runtime, sourceRoot, beforeApply))
	return mux
}

func dashboardHandler(runtime *Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_markata/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		page := strings.Replace(webHTML, "/*MARKATA_BROWSER_THEME_TOKENS*/", BrowserTokenStylesheet(runtime.Snapshot().Theme), 1)
		if _, err := w.Write([]byte(page)); err != nil {
			return
		}
	}
}

func stateHandler(runtime *Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if err := json.NewEncoder(w).Encode(runtime.Snapshot()); err != nil {
			return
		}
	}
}

func actionsHandler(runtime *Runtime) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", "POST")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}
		if !sameOriginLocalRequest(r, scheme) {
			http.Error(w, "action requires same-origin request", http.StatusForbidden)
			return
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
		var action ActionRequest
		if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
			http.Error(w, "invalid action", 400)
			return
		}
		if err := runtime.Trigger(action); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}
}

func requestScheme(r *http.Request) string {
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func sameOriginLocalRequest(r *http.Request, scheme string) bool {
	if !loopbackAddress(r.RemoteAddr) || !localHost(r.Host) {
		return false
	}
	origin := r.Header.Get("Origin")
	if origin == "" {
		return r.Header.Get("Sec-Fetch-Site") == "" || r.Header.Get("Sec-Fetch-Site") == "same-origin"
	}
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Scheme == scheme && parsed.Host == r.Host
}

func loopbackAddress(address string) bool {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func localHost(address string) bool {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		host = address
	} else if port != "" {
		if _, err := strconv.Atoi(port); err != nil {
			return false
		}
	}
	host = strings.Trim(host, "[]")
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

const webHTML = `<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Serve Control Center · Markata</title>
<style>
/*MARKATA_BROWSER_THEME_TOKENS*/
:root {
  --bg: var(--markata-background);
  --panel: var(--markata-panel);
  --surface: var(--markata-surface);
  --elevated: var(--markata-elevated);
  --text: var(--markata-text-primary);
  --muted: var(--markata-text-secondary);
  --accent: var(--markata-accent);
  --link: var(--markata-link);
  --border: var(--markata-border);
  --focus: var(--markata-focus);
  --success: var(--markata-success);
  --warning: var(--markata-warning);
  --error: var(--markata-error);
  --info: var(--markata-info);
  --code-bg: var(--markata-code-background);
  --code-text: var(--markata-code-text);
  --button-bg: var(--markata-button-background);
  --button-text: var(--markata-button-text);
  --sans: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  --mono: ui-monospace, "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  --page: var(--bg);
  --card: var(--panel);
  --ink: var(--text);
  --line: var(--border);
  --radius: 7px;
  font-family: var(--sans);
  background: var(--bg);
  color: var(--text);
}
* { box-sizing: border-box; }
html { min-height: 100%; background: var(--bg); }
body { min-height: 100vh; margin: 0; background: var(--bg); color: var(--text); font: 14px/1.5 var(--sans); }
button, input { font: inherit; }
button { color: inherit; cursor: pointer; }
a { color: var(--link); }
button:focus-visible, input:focus-visible, a:focus-visible, [tabindex]:focus-visible { outline: 2px solid var(--focus); outline-offset: 3px; }
.app { min-height: 100vh; max-width: 1760px; margin: 0 auto; display: grid; grid-template-columns: 232px minmax(0, 1fr); }
.rail { height: 100vh; min-height: 440px; position: sticky; top: 0; display: flex; flex-direction: column; padding: 22px 14px 16px; background: var(--panel); border-right: 1px solid var(--border); }
.brand { padding: 2px 10px 20px; border-bottom: 1px solid var(--border); color: var(--text); font-size: 18px; font-weight: 700; letter-spacing: -.04em; }
.brand b { display: inline-grid; width: 28px; height: 28px; margin-right: 8px; place-items: center; border: 1px solid var(--border); border-radius: 6px; background: var(--accent); color: var(--button-text); font: 700 15px var(--mono); }
.brand small { display: block; margin: 10px 0 0 36px; color: var(--muted); font: 10px var(--mono); letter-spacing: .1em; text-transform: uppercase; }
.nav { display: grid; gap: 4px; margin-top: 18px; }
.tab { display: flex; justify-content: space-between; gap: 12px; align-items: center; width: 100%; padding: 10px 11px; border: 1px solid transparent; border-radius: 5px; background: transparent; color: var(--muted); text-align: left; font: 12px var(--mono); }
.tab:hover { background: var(--surface); color: var(--text); }
.tab.active { border-color: var(--border); background: var(--surface); color: var(--text); }
.tab span { color: var(--muted); font-variant-numeric: tabular-nums; }
.rail-foot { margin-top: auto; padding: 14px 10px 0; border-top: 1px solid var(--border); color: var(--muted); font: 10px/1.8 var(--mono); letter-spacing: .04em; }
.main { min-width: 0; padding: 28px clamp(18px, 3vw, 46px) 64px; }
.mast { display: flex; justify-content: space-between; gap: 22px; align-items: flex-start; padding: 1px 0 22px; border-bottom: 1px solid var(--border); }
.eyebrow, .kicker, .metric-label, .detail-label { color: var(--muted); font: 10px var(--mono); letter-spacing: .12em; text-transform: uppercase; }
h1, h2, h3, p { margin: 0; }
h1 { margin-top: 7px; font-size: clamp(27px, 3vw, 39px); line-height: 1.08; letter-spacing: -.045em; }
h2 { margin-top: 4px; font-size: 22px; line-height: 1.2; letter-spacing: -.035em; }
h3 { margin: 9px 0; font-size: 20px; line-height: 1.25; letter-spacing: -.03em; overflow-wrap: anywhere; }
.mast p { max-width: 64ch; margin-top: 8px; color: var(--muted); }
.mast-actions { display: flex; flex-wrap: wrap; justify-content: flex-end; gap: 8px; align-items: center; }
.button { min-height: 36px; padding: 7px 11px; border: 1px solid var(--border); border-radius: 5px; background: var(--surface); color: var(--text); font-size: 12px; font-weight: 650; text-decoration: none; }
.button:hover { border-color: var(--accent); background: var(--elevated); text-decoration: none; }
.button.primary { border-color: var(--button-bg); background: var(--button-bg); color: var(--button-text); }
.button.primary:hover { opacity: .88; }
.site-link { color: var(--link); font: 11px var(--mono); }
.metrics { display: grid; grid-template-columns: repeat(5, minmax(0, 1fr)); margin-top: 20px; overflow: hidden; border: 1px solid var(--border); border-radius: var(--radius); background: var(--panel); }
.metric { min-width: 0; padding: 14px 16px; border-right: 1px solid var(--border); }
.metric:last-child { border-right: 0; }
.metric-label { white-space: nowrap; }
.metric-value { margin-top: 6px; overflow: hidden; color: var(--text); font-size: 21px; font-weight: 650; letter-spacing: -.035em; line-height: 1.25; text-overflow: ellipsis; white-space: nowrap; }
.metric-note { overflow: hidden; color: var(--muted); font: 10px var(--mono); text-overflow: ellipsis; white-space: nowrap; }
.metric.warning .metric-value { color: var(--warning); }
.metric.failed .metric-value { color: var(--error); }
.metric.success .metric-value { color: var(--success); }
.content-head { display: flex; justify-content: space-between; align-items: end; gap: 16px; margin: 28px 0 12px; }
.filter { width: min(300px, 100%); min-height: 38px; padding: 8px 10px; border: 1px solid var(--border); border-radius: 5px; background: var(--code-bg); color: var(--text); font: 12px var(--mono); }
.filter::placeholder { color: var(--muted); }
.workspace { display: grid; grid-template-columns: minmax(0, 1fr) minmax(290px, 36%); gap: 14px; align-items: start; }
.list, .detail { min-width: 0; overflow: hidden; border: 1px solid var(--border); border-radius: var(--radius); background: var(--panel); }
.list { min-height: 280px; }
.detail { position: sticky; top: 16px; max-height: calc(100vh - 32px); overflow: auto; }
.item { display: block; width: 100%; padding: 14px 16px; border: 0; border-bottom: 1px solid var(--border); border-radius: 0; background: transparent; color: var(--text); text-align: left; }
.item:last-child { border-bottom: 0; }
.item:hover { background: var(--surface); }
.item.selected { padding-left: 13px; border-left: 3px solid var(--accent); background: var(--surface); }
.item-top { display: flex; gap: 8px; align-items: center; min-width: 0; }
.item-title { overflow: hidden; font-weight: 620; text-overflow: ellipsis; white-space: nowrap; }
.item-meta { margin: 5px 0 0 1px; overflow: hidden; color: var(--muted); font: 10px var(--mono); text-overflow: ellipsis; white-space: nowrap; }
.item-tail { margin-left: auto; color: var(--muted); font: 10px var(--mono); white-space: nowrap; }
.pill { flex: none; padding: 3px 6px; border: 1px solid var(--border); border-radius: 3px; background: var(--surface); color: var(--muted); font: 9px var(--mono); text-transform: uppercase; }
.pill.success, .pill.ready { color: var(--success); }
.pill.running { color: var(--info); }
.pill.warning, .pill.warn { color: var(--warning); }
.pill.failed, .pill.error { color: var(--error); }
.detail-head { padding: 18px 19px; border-bottom: 1px solid var(--border); background: var(--surface); }
.detail-path { color: var(--link); font: 10px var(--mono); overflow-wrap: anywhere; }
.detail-body { padding: 15px 19px 20px; }
.detail-label { margin: 16px 0 6px; }
.detail-label:first-child { margin-top: 0; }
.detail-text { color: var(--text); white-space: pre-wrap; overflow-wrap: anywhere; }
.fix-tools { display: flex; flex-wrap: wrap; align-items: center; gap: 7px; margin: 0 0 12px; }
.fix-tools .summary { margin-right: 8px; color: var(--muted); font-size: 12px; }
.fix-review { width: min(1120px, calc(100vw - 24px)); max-height: min(88vh, 900px); overflow: auto; padding: 18px; }
.fix-review h2 { margin: 0 0 8px; }
.fix-review .fix-file { margin-top: 18px; border-top: 1px solid var(--border); }
.fix-review .fix-file h3 { font: 12px var(--mono); overflow-wrap: anywhere; }
.fix-review .diff-host { min-height: 90px; max-height: 52vh; overflow: auto; border: 1px solid var(--border); border-radius: 5px; }
.fix-review .fix-edits { margin-top: 9px; color: var(--muted); font-size: 12px; }
.fix-review .edit { padding: 12px 0; border-top: 1px solid var(--border); }
.fix-review pre { max-height: 190px; overflow: auto; padding: 10px; background: var(--code-bg); color: var(--code-text); white-space: pre-wrap; overflow-wrap: anywhere; }
.fix-review .actions { display: flex; gap: 8px; justify-content: flex-end; position: sticky; bottom: -18px; padding: 12px 0; background: var(--panel); }
.detail-actions { display: flex; flex-wrap: wrap; gap: 7px; margin-top: 16px; }
.step, .log-line { padding: 7px 0; border-bottom: 1px solid var(--border); font: 10px/1.55 var(--mono); overflow-wrap: anywhere; }
.step { display: flex; gap: 8px; align-items: baseline; }
.step:last-child, .log-line:last-child { border-bottom: 0; }
.time { margin-right: 8px; color: var(--muted); }
.log-line.error { color: var(--error); }
.log-line.warning { color: var(--warning); }
.empty { padding: 42px 18px; color: var(--muted); text-align: center; }
.empty strong { display: block; margin-bottom: 4px; color: var(--text); font-size: 16px; }
.notice { display: none; margin-top: 12px; padding: 10px 12px; border: 1px solid color-mix(in srgb, var(--error) 45%, var(--border)); border-radius: 5px; background: color-mix(in srgb, var(--error) 12%, var(--panel)); color: var(--text); font: 12px var(--mono); }
.notice.show { display: block; }
.mobile-nav { display: none; }
.sr-only { position: absolute; width: 1px; height: 1px; padding: 0; margin: -1px; overflow: hidden; clip: rect(0,0,0,0); white-space: nowrap; border: 0; }
dialog { width: min(440px, calc(100vw - 32px)); padding: 0; border: 1px solid var(--border); border-radius: 8px; background: var(--panel); color: var(--text); box-shadow: 0 24px 80px #0008; }
dialog::backdrop { background: #0009; }
.help-head { display: flex; justify-content: space-between; align-items: center; padding: 15px 18px; border-bottom: 1px solid var(--border); }
.help-list { display: grid; grid-template-columns: auto 1fr; gap: 8px 16px; padding: 16px 18px 20px; color: var(--muted); font-size: 12px; }
kbd { min-width: 26px; padding: 2px 5px; border: 1px solid var(--border); border-radius: 3px; background: var(--surface); color: var(--text); font: 10px var(--mono); text-align: center; }
@media (max-width: 1100px) { .metrics { grid-template-columns: repeat(3, 1fr); } .metric:nth-child(3) { border-right: 0; } .metric:nth-child(n+4) { border-top: 1px solid var(--border); } }
@media (max-width: 800px) { .app { display: block; } .rail { display: none; } .main { padding: 20px 15px 78px; } .mast { flex-wrap: wrap; } .mast-actions { justify-content: flex-start; } .metrics { grid-template-columns: repeat(2, minmax(0, 1fr)); } .metric:nth-child(3) { border-right: 1px solid var(--border); } .metric:nth-child(even) { border-right: 0; } .metric:nth-child(n+3) { border-top: 1px solid var(--border); } .workspace { display: block; } .detail { position: static; max-height: none; margin-top: 12px; } .mobile-nav { position: fixed; z-index: 5; right: 0; bottom: 0; left: 0; display: flex; gap: 2px; padding: 6px max(6px, env(safe-area-inset-left)) calc(6px + env(safe-area-inset-bottom)); border-top: 1px solid var(--border); background: var(--panel); } .mobile-nav .tab { justify-content: center; flex: 1; padding: 10px 2px; font-size: 9px; } .mobile-nav .tab span { display: none; } .content-head { flex-wrap: wrap; } .filter { width: 100%; } }
@media (max-width: 480px) { .metrics { grid-template-columns: repeat(2, minmax(0, 1fr)); } .metric { padding: 11px; } .metric-value { font-size: 18px; } .item { padding: 12px; } .item.selected { padding-left: 9px; } .item-top { flex-wrap: wrap; } .item-tail { margin-left: 0; } }
@media (prefers-reduced-motion: reduce) { *, *::before, *::after { scroll-behavior: auto !important; transition-duration: .01ms !important; animation-duration: .01ms !important; animation-iteration-count: 1 !important; } }
</style></head><body><div class="app"><aside class="rail"><div class="brand"><b>M</b>markata<small>Serve control center</small></div><nav class="nav" id="desktop-nav" aria-label="Sections"></nav><div class="rail-foot">LOCAL SESSION<br><span id="sync-label">Connecting…</span></div></aside><main class="main"><header class="mast"><div><div class="eyebrow">Development / live session</div><h1>Serve Control Center</h1><p>Build activity, page health, and the details behind every warning.</p></div><div class="mast-actions"><a class="button" href="/_markata/subscriptions">Manage subscriptions ↗</a><a class="site-link" id="site-link" hidden target="_blank" rel="noopener">Open site ↗</a><button class="button" id="copy-report-button" type="button" hidden>Copy report</button><button class="button primary" id="build-button" type="button">Trigger build ↗</button></div></header><section class="metrics" id="metrics" aria-label="Current state"></section><div class="notice" id="notice" role="alert"></div><div class="content-head"><div><div class="kicker" id="kicker">Activity</div><h2 id="title">Jobs</h2></div><label><span class="sr-only">Filter current section</span><input class="filter" id="filter" type="search" placeholder="Filter this view…"></label></div><div class="fix-tools" id="fix-tools" hidden></div><div class="workspace"><section class="list" id="list" aria-label="Results" aria-live="polite"></section><aside class="detail" id="detail" aria-label="Selected item" tabindex="-1"></aside></div></main></div><nav class="mobile-nav" id="mobile-nav" aria-label="Sections"></nav><dialog id="keyboard-help" aria-labelledby="keyboard-help-title"><div class="help-head"><h2 id="keyboard-help-title">Keyboard shortcuts</h2><button class="button" id="close-help" type="button">Close</button></div><div class="help-list"><kbd>j</kbd><span>Next item</span><kbd>k</kbd><span>Previous item</span><kbd>Enter</kbd><span>Focus details or activate the focused control</span><kbd>Esc</kbd><span>Clear selection or close help</span><kbd>/</kbd><span>Focus filter</span><kbd>g</kbd><span>First item</span><kbd>G</kbd><span>Last item</span><kbd>r</kbd><span>Rerun selected completed job</span><kbd>w</kbd><span>Warnings</span><kbd>e</kbd><span>Errors</span><kbd>p</kbd><span>Pages</span><kbd>f</kbd><span>Feeds</span><kbd>l</kbd><span>Logs</span><kbd>?</kbd><span>Show this help</span></div></dialog><dialog id="fix-review" class="fix-review" aria-labelledby="fix-review-title"><h2 id="fix-review-title">Review source fixes</h2><div id="fix-review-body"></div><div class="actions"><button class="button" id="fix-cancel" type="button">Cancel</button><button class="button primary" id="fix-approve" type="button">Apply reviewed fixes</button></div></dialog>
<script>
(() => {
'use strict';
const capabilities=Object.freeze({build:true,rerun:true,rebuildPage:true,previewFixes:true,applyFixes:true}),state={snapshot:null,section:'jobs',selected:'',filter:'',severity:'',noSelection:false,batchPreview:null,selectedFixes:new Set(),batchResultMessage:"",batchBaseline:null},tabs=[['jobs','Jobs'],['diagnostics','Problems'],['pages','Pages'],['feeds','Feeds'],['logs','Logs']];
const el=id=>document.getElementById(id),arr=x=>Array.isArray(x)?x:[],esc=x=>String(x==null?'':x).replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c])),when=x=>x&&!isNaN(new Date(x))?new Date(x).toLocaleTimeString():'';
function field(name,value){return value==null||value===''?'':'<div class="detail-label">'+esc(name)+'</div><div class="detail-text">'+esc(value)+'</div>'}
function metric(name,value,note,kind){return '<div class="metric '+esc(kind||'')+'"><div class="metric-label">'+esc(name)+'</div><div class="metric-value">'+esc(value)+'</div><div class="metric-note">'+esc(note)+'</div></div>'}
function key(item,i){if(state.section==='jobs')return 'job:'+item.id;if(state.section==='diagnostics')return 'diagnostic:'+([item.job_id,item.code,item.file||item.page,item.line,item.column,item.message].filter(Boolean).join('|')||i);if(state.section==='pages')return 'page:'+(item.path||item.file||item.slug);if(state.section==='feeds')return 'feed:'+item.name;return 'log:'+([item.time,item.job_id,item.step_id,item.message].filter(Boolean).join('|')||i)}
function routeHash(){const query=new URLSearchParams();if(state.selected)query.set('item',state.selected);if(state.filter)query.set('filter',state.filter);if(state.severity)query.set('severity',state.severity);if(state.noSelection)query.set('none','1');const suffix=query.toString();return '#/'+state.section+(suffix?'?'+suffix:'')}
function writeRoute(push){const hash=routeHash();if(location.hash===hash)return;const prior=history.state&&history.state.markataServeControl?history.state.depth||0:0;const entry={markataServeControl:true,depth:push?prior+1:prior};if(push)history.pushState(entry,'',hash);else history.replaceState(entry,'',hash)}
function readRoute(){const value=location.hash.replace(/^#/,''),parts=value.split('?'),section=parts[0].replace(/^\/+/,''),query=new URLSearchParams(parts.slice(1).join('?'));if(tabs.some(x=>x[0]===section))state.section=section;state.selected=query.get('item')||'';state.filter=query.get('filter')||'';state.severity=query.get('severity')||'';state.noSelection=query.get('none')==='1';el('filter').value=state.filter}
function navigate(section,filter='',selected='',severity=''){state.section=section;state.filter=filter;state.selected=selected;state.severity=severity;state.noSelection=false;el('filter').value=filter;writeRoute(true);render()}
function selectItem(selected){state.selected=selected;state.noSelection=false;writeRoute(true);list()}
function visibleItems(){const query=state.filter.toLowerCase();return resourceItems(state.section).map((item,i)=>({item,k:key(item,i)})).filter(({item})=>(!query||JSON.stringify(item).toLowerCase().includes(query))&&(!state.severity||item.severity===state.severity))}
function diagnosticReport(items){const warnings=items.filter(item=>item.severity==='warning').length,errors=items.filter(item=>item.severity==='error').length,lines=['# Markata Problems Report','',items.length+' problem'+(items.length===1?'':'s')+' ('+errors+' errors, '+warnings+' warnings)'];if(state.severity)lines.push('Severity filter: '+state.severity);if(state.filter)lines.push('Text filter: '+state.filter);if(!items.length){lines.push('','No matching problems.');return lines.join('\n')}items.forEach((item,index)=>{const location=[item.file||item.page,item.line,item.column].filter(value=>value!==undefined&&value!==null&&value!=='').join(':');lines.push('',(index+1)+'. ['+String(item.severity||'problem').toUpperCase()+'] '+(item.code?item.code+' — ':'')+(item.message||'Problem'));if(location)lines.push('   Location: '+location);if(item.explanation)lines.push('   Explanation: '+item.explanation);if(item.suggested_fix)lines.push('   Suggested fix: '+item.suggested_fix);if(item.job_id)lines.push('   Job: '+item.job_id);if(item.step_id)lines.push('   Step: '+item.step_id)});return lines.join('\n')}
async function copyProblemReport(){const items=visibleItems().map(x=>x.item),report=diagnosticReport(items);if(!navigator.clipboard||!navigator.clipboard.writeText){notice('Clipboard access is unavailable in this browser.');return}try{await navigator.clipboard.writeText(report);notice('Copied '+items.length+' problem'+(items.length===1?'':'s')+' to clipboard.')}catch(err){notice('Could not copy problem report: '+(err&&err.message?err.message:'clipboard write failed'))}}
function moveSelection(delta){const items=visibleItems();if(!items.length)return;let index=items.findIndex(x=>x.k===state.selected);index=index<0?(delta>0?-1:items.length):index;index=Math.max(0,Math.min(items.length-1,index+delta));state.selected=items[index].k;state.noSelection=false;writeRoute(true);list();const active=el('list').querySelector('[data-key="'+CSS.escape(state.selected)+'"]');if(active)active.focus()}
function feedPath(item){const name=String(item.path||item.name||'').replace(/^\/+|\/+$/g,'');return name?'/'+name+'/':'/'}
function title(item){return state.section==='jobs'?item.name||item.type||item.id:state.section==='diagnostics'?item.message||item.code:state.section==='pages'?item.path||item.file||item.slug:state.section==='feeds'?item.title||feedPath(item):item.message}
function resourceItems(section){const snapshot=state.snapshot||{};return section==='diagnostics'?arr(snapshot.current_diagnostics||snapshot.diagnostics):arr(snapshot[section])}
function count(section){return resourceItems(section).length}
function updateHTML(target,html){if(target.innerHTML!==html)target.innerHTML=html}
function nav(target){updateHTML(target,tabs.map(([id,name])=>'<button class="tab '+(state.section===id?'active':'')+'" data-section="'+id+'" type="button">'+name+'<span>'+count(id)+'</span></button>').join(''))}
function status(){const s=state.snapshot||{},server=s.server||{},site=s.site||{},jobs=arr(s.jobs),diags=arr(s.diagnostics),warnings=diags.filter(d=>d.severity==='warning').length,errors=diags.filter(d=>d.severity==='error').length,serverState=server.status||(state.snapshot?'ready':'connecting'),siteState=site.status||'waiting';el('metrics').innerHTML=metric('Server',serverState,server.address||'Local session',serverState==='failed'?'failed':'success')+metric('Site',siteState,(site.page_count||0)+' pages'+(site.message?' · '+site.message:''),siteState==='failed'?'failed':'success')+metric('Jobs',jobs.length,jobs.filter(j=>j.state==='running').length+' running')+metric('Warnings',warnings,'Session history',warnings?'warning':'')+metric('Errors',errors,'Session history',errors?'failed':'success');const url=server.address?(String(server.address).startsWith('http')?server.address:'http://'+server.address):'';el('site-link').hidden=!url;if(url)el('site-link').href=url}
function applyTheme(theme){const names={background:'--markata-background',panel:'--markata-panel',surface:'--markata-surface',elevated:'--markata-elevated','text-primary':'--markata-text-primary','text-secondary':'--markata-text-secondary',accent:'--markata-accent',link:'--markata-link',border:'--markata-border',focus:'--markata-focus',success:'--markata-success',warning:'--markata-warning',error:'--markata-error',info:'--markata-info','code-background':'--markata-code-background','code-text':'--markata-code-text','button-background':'--markata-button-background','button-text':'--markata-button-text'};for(const [key,variable] of Object.entries(names)){const color=theme&&theme[key];if(color&&CSS.supports('color',color))document.documentElement.style.setProperty(variable,color)}if(theme&&['dark','light'].includes(theme.mode))document.documentElement.style.colorScheme=theme.mode}
function knownPage(path){return arr((state.snapshot||{}).pages).some(p=>p.path===path)}
function row(item,k){const kind=state.section==='jobs'?item.state:state.section==='diagnostics'?item.severity:state.section==='pages'?(item.state||item.status):state.section==='feeds'?(item.status||'feed'):item.level;const meta=state.section==='jobs'?[item.trigger,item.id].filter(Boolean).join(' · '):state.section==='diagnostics'?[item.code,item.file||item.page,item.line?':'+item.line:''].filter(Boolean).join(' · '):state.section==='pages'?(item.file||item.path||''):state.section==='feeds'?[item.path,arr(item.entries).length+' posts'].filter(Boolean).join(' · '):(item.job_id||'');return '<button class="item '+(state.selected===k?'selected':'')+'" data-key="'+esc(k)+'" type="button"><div class="item-top"><span class="pill '+esc(kind)+'">'+esc(kind||'item')+'</span><span class="item-title">'+esc(title(item))+'</span><span class="item-tail">'+esc(when(item.ended_at||item.time))+'</span></div><div class="item-meta">'+esc(meta)+'</div></button>'}
function logLine(l){return '<div class="log-line '+esc(l.level)+'"><span class="time">'+esc(when(l.time))+'</span>'+esc(l.message)+'</div>'}
function detail(item){if(!item){updateHTML(el('detail'),'<div class="empty"><strong>Nothing selected</strong>Choose an item to inspect it.</div>');return}let head='',body='';if(state.section==='jobs'){head='<span class="pill '+esc(item.state)+'">'+esc(item.state)+'</span><h3>'+esc(title(item))+'</h3><div class="detail-path">'+esc(item.id)+'</div>';body=field('Trigger',item.trigger)+field('Type',item.type)+field('Started',item.started_at)+field('Finished',item.ended_at);if(arr(item.steps).length)body+='<div class="detail-label">Steps</div>'+item.steps.map(x=>'<div class="step"><span class="pill '+esc(x.state)+'">'+esc(x.state)+'</span>'+esc(x.name||x.id)+'</div>').join('');if(arr(item.pages).length)body+='<div class="detail-label">Affected pages</div>'+item.pages.map(path=>'<div class="step">'+(knownPage(path)?'<button class="button" type="button" data-page="'+esc(path)+'">'+esc(path)+' →</button>':esc(path))+'</div>').join('');if(arr(item.diagnostics).length)body+='<div class="detail-label">Diagnostics</div>'+item.diagnostics.map(d=>'<div class="step"><span class="pill '+esc(d.severity)+'">'+esc(d.code||d.severity)+'</span>'+esc(d.message)+'</div>').join('');if(arr(item.logs).length)body+='<div class="detail-label">Recent logs</div>'+item.logs.slice(-20).map(logLine).join('');body+='<div class="detail-actions">'+(capabilities.rerun?'<button class="button" type="button" data-rerun="'+esc(item.id)+'">Rerun job ↗</button>':'')+'<button class="button" type="button" data-job-logs="'+esc(item.id)+'">View logs →</button></div>'}
else if(state.section==='diagnostics'){head='<span class="pill '+esc(item.severity)+'">'+esc(item.severity)+'</span><h3>'+esc(item.message)+'</h3><div class="detail-path">'+esc([item.file||item.page,item.line].filter(Boolean).join(':'))+'</div>';body=field('Code',item.code)+field('Explanation',item.explanation)+field('Suggested fix',item.suggested_fix)+field('Job',item.job_id)+field('Step',item.step_id);body+='<div class="detail-actions">';if(item.job_id)body+='<button class="button" type="button" data-job="'+esc(item.job_id)+'">View job →</button>';if(knownPage(item.page||item.file))body+='<button class="button" type="button" data-page="'+esc(item.page||item.file)+'">View page →</button>';if(capabilities.previewFixes){for(const plan of arr(item.fix_plans)){const id=plan.file+'|'+plan.id;body+='<button class="button" type="button" data-fix-select="'+esc(id)+'">'+(state.selectedFixes.has(id)?'✓ Selected':'Select')+' · '+esc(plan.safety)+' fix</button><button class="button" type="button" data-fix-one="'+esc(id)+'">Review fix</button>'}}body+='</div>'}
else if(state.section==='pages'){const path=item.path||item.file||item.slug||'',diags=arr(item.diagnostics);head='<span class="pill '+esc(item.state||item.status)+'">'+esc(item.state||item.status||'page')+'</span><h3>'+esc(path)+'</h3><div class="detail-path">'+esc(item.file||path)+'</div>';body=field('Last job',item.job_id||item.last_job_id)+'<div class="detail-label">Diagnostics</div>'+(diags.length?diags.map(d=>'<div class="step"><span class="pill '+esc(d.severity)+'">'+esc(d.code||d.severity)+'</span>'+esc(d.message)+'</div>').join(''):'<div class="detail-text">No current diagnostics.</div>');body+='<div class="detail-actions">';if(capabilities.rebuildPage&&/\.(md|markdown)$/i.test(path))body+='<button class="button" type="button" data-rebuild="'+esc(path)+'">Rebuild page ↗</button>';if(item.url){try{const url=new URL(item.url,location.origin);if(url.origin===location.origin)body+='<a class="button" href="'+esc(url.href)+'" target="_blank" rel="noopener">Open page ↗</a>'}catch{}}body+='</div>'}
else if(state.section==='feeds'){head='<span class="pill '+esc(item.status||'feed')+'">'+esc(item.status||'feed')+'</span><h3>'+esc(title(item))+'</h3><div class="detail-path">'+esc(feedPath(item))+'</div>';body=field('Name',item.name)+field('Source/config',item.path)+field('Entries',arr(item.entries).length);if(arr(item.entries).length)body+='<div class="detail-label">Posts</div>'+item.entries.slice(0,250).map(entry=>'<div class="step">'+(knownPage(entry.path)?'<button class="button" type="button" data-page="'+esc(entry.path)+'">'+esc(entry.title||entry.path)+' →</button>':esc(entry.title||entry.path))+'</div>').join('')}
else{head='<span class="pill '+esc(item.level)+'">'+esc(item.level||'log')+'</span><h3>Log entry</h3><div class="detail-path">'+esc(when(item.time))+'</div>';body=field('Message',item.message)+field('Job',item.job_id)+field('Step',item.step_id)}
updateHTML(el('detail'),'<div class="detail-head">'+head+'</div><div class="detail-body">'+body+'</div>')}
function list(){const matches=visibleItems(),items=state.section==='logs'?matches.slice(-200):matches;if(state.section==='logs')el('kicker').textContent='Session output · latest '+items.length+' of '+matches.length+(matches.length>items.length?' · refine the filter to reach older logs':'');if(!items.some(x=>x.k===state.selected)){state.selected=(!state.noSelection&&items.length)?items[0].k:'';writeRoute(false)}updateHTML(el('list'),items.length?items.map(({item,k})=>row(item,k)).join(''):'<div class="empty"><strong>No matching items</strong>Try another filter or wait for the next build.</div>');detail((items.find(x=>x.k===state.selected)||{}).item)}
function fixPlans(){const plans=new Map();for(const diagnostic of resourceItems('diagnostics'))for(const plan of arr(diagnostic.fix_plans))plans.set(plan.file+'|'+plan.id,plan);return plans}
function renderFixTools(){const target=el('fix-tools');target.hidden=state.section!=='diagnostics';if(target.hidden)return;const plans=[...fixPlans().values()],safe=plans.filter(x=>x.safety==='safe'),groups=new Map();for(const plan of safe)groups.set(plan.category,(groups.get(plan.category)||0)+1);updateHTML(target,'<span class="summary">Safe automatic fixes '+safe.length+' · Needs review '+plans.filter(x=>x.safety==='review').length+' · Manual '+resourceItems('diagnostics').filter(x=>!arr(x.fix_plans).length||x.fix_safety==='manual').length+'</span>'+(safe.length?'<button class="button primary" type="button" data-fix-group="all">Review all safe</button>':'')+[...groups].map(([category,count])=>'<button class="button" type="button" data-fix-group="'+esc(category)+'">'+esc(category)+' · '+count+'</button>').join('')+(state.selectedFixes.size?'<button class="button" type="button" data-fix-group="selected">Review selected · '+state.selectedFixes.size+'</button>':''))}
function render(){const focused=document.activeElement,focusedKey=focused?.dataset?.key,focusedSection=focused?.dataset?.section,focusedParent=focused?.parentElement?.id;nav(el('desktop-nav'));nav(el('mobile-nav'));status();el('copy-report-button').hidden=state.section!=='diagnostics';el('title').textContent=tabs.find(x=>x[0]===state.section)[1];el('kicker').textContent=state.section==='diagnostics'?'Problems to resolve':state.section==='pages'?'Built content':state.section==='feeds'?'Generated feeds':state.section==='logs'?'Session output':'Activity';list();renderFixTools();if(focused&&!focused.isConnected){const parent=focusedParent&&el(focusedParent),replacement=focusedKey?el('list').querySelector('[data-key="'+CSS.escape(focusedKey)+'"]'):focusedSection&&parent?parent.querySelector('[data-section="'+CSS.escape(focusedSection)+'"]'):null;if(replacement)replacement.focus({preventScroll:true})}}
function notice(message){el('notice').textContent=message;el('notice').classList.toggle('show',!!message)}
async function poll(){try{const res=await fetch('/_markata/api/state',{cache:'no-store'});if(!res.ok)throw Error('State request failed ('+res.status+')');state.snapshot=await res.json();applyTheme(state.snapshot.theme);el('sync-label').textContent='Live · updated '+when(new Date());if(state.batchBaseline){const jobs=arr(state.snapshot.jobs),latest=jobs[jobs.length-1];if(jobs.length>state.batchBaseline.jobs&&latest&&latest.ended_at){const remaining=new Set(resourceItems('diagnostics').map(d=>[d.code,d.file||d.page,d.line,d.message].join('|'))),resolved=[...state.batchBaseline.keys].filter(key=>!remaining.has(key)).length;state.batchResultMessage+=' '+resolved+' diagnostics resolved after rebuild.';state.batchBaseline=null}}notice(state.batchResultMessage);render()}catch(err){el('sync-label').textContent='Connection interrupted';notice(err.message)}}
async function action(request){try{const res=await fetch('/_markata/api/actions',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(request)});if(!res.ok)throw Error((await res.text()).trim()||'Action failed');notice('');await poll()}catch(err){notice(err.message)}}
let pierreRendererPromise=null,pierreRenderer=null;
function cleanupReviewDiffs(){if(!pierreRenderer)return;el('fix-review-body').querySelectorAll('[data-pierre-file]').forEach(host=>pierreRenderer.cleanupPierreDiff(host))}
function plainReviewDiffs(previews){el('fix-review-body').querySelectorAll('[data-pierre-file]').forEach((host,index)=>{host.innerHTML=arr(previews[index]?.edits).map(edit=>'<pre>− '+esc(edit.before)+'\n+ '+esc(edit.after)+'</pre>').join('')})}
async function reviewFixes(mode){state.batchResultMessage='';const plans=fixPlans(),selected=[...plans].filter(([key,plan])=>mode==='selected'?state.selectedFixes.has(key):mode==='all'?plan.safety==='safe':mode===key||mode===plan.category&&plan.safety==='safe');if(!selected.length){notice('No matching current fixes to review.');return}const byFile=new Map();for(const [,plan] of selected){if(!byFile.has(plan.file))byFile.set(plan.file,[]);byFile.get(plan.file).push(plan.id)}try{const files=[...byFile].map(([path,ids])=>({path,selection:{ids}}));const res=await fetch('/_markata/api/fixes/batch/preview',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({files})});if(!res.ok)throw Error((await res.text()).trim()||'Could not preview fixes');const payload=await res.json(),previews=arr(payload.files);if(!previews.some(file=>arr(file.edits).length)){notice('No current edits match the selection.');return}cleanupReviewDiffs();state.batchPreview=previews;el('fix-review-body').innerHTML='<p>Review '+previews.reduce((n,file)=>n+arr(file.edits).length,0)+' edits across '+previews.length+' files. Files changed since this preview will be skipped.</p>'+previews.map((file,index)=>'<section class="fix-file"><h3>'+esc(file.path)+'</h3><div class="diff-host" data-pierre-file="'+index+'" role="region" aria-label="Changes in '+esc(file.path)+'">Loading diff…</div><div class="fix-edits">'+arr(file.edits).map(edit=>'<div class="edit"><strong>'+esc(edit.category)+'</strong> · '+esc(edit.safety)+' · line '+esc(edit.line)+'<p>'+esc(edit.explanation||edit.message)+'</p></div>').join('')+'</div></section>').join('');el('fix-approve').disabled=true;el('fix-review').showModal();try{pierreRendererPromise??=import('/_markata/assets/pierre-diffs/pierre-review.js');const renderer=await pierreRendererPromise;if(!el('fix-review').open||state.batchPreview!==previews)return;pierreRenderer=renderer;el('fix-review-body').querySelectorAll('[data-pierre-file]').forEach((host,index)=>renderer.renderPierreDiff(host,{path:previews[index].path,before:previews[index].before,after:previews[index].after,diffStyle:matchMedia('(max-width: 800px)').matches?'unified':'split'}))}catch(renderError){if(state.batchPreview===previews&&el('fix-review').open){cleanupReviewDiffs();plainReviewDiffs(previews);notice('Rich diff unavailable; showing text changes.')}}finally{if(state.batchPreview===previews&&el('fix-review').open)el('fix-approve').disabled=false}}catch(err){notice(err.message)}}
async function applyBatch(){const previews=state.batchPreview;if(!previews)return;el('fix-approve').disabled=true;const baseline={jobs:arr(state.snapshot?.jobs).length,keys:new Set(resourceItems('diagnostics').map(d=>[d.code,d.file||d.page,d.line,d.message].join('|')))};try{const files=previews.map(file=>({path:file.path,digest:file.digest,selection:file.selection}));const res=await fetch('/_markata/api/fixes/batch/apply',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({files})});if(!res.ok)throw Error((await res.text()).trim()||'Could not apply fixes');const payload=await res.json(),results=arr(payload.result?.files),applied=arr(payload.result?.applied).length,skipped=results.filter(file=>file.status==='stale'||file.status==='skipped'||file.status==='failed').length;state.batchPreview=null;state.selectedFixes.clear();state.batchBaseline=payload.rebuild_requested?baseline:null;state.batchResultMessage='Applied '+applied+' fixes across '+results.filter(file=>file.status==='applied').length+' files; '+skipped+' skipped. '+(payload.rebuild_error?'Rebuild could not start: '+payload.rebuild_error:payload.rebuild_requested?'One rebuild queued.':'No rebuild needed.');el('fix-review').close();await poll()}catch(err){notice(err.message)}finally{el('fix-approve').disabled=false}}
function moveToEdge(last){const items=visibleItems();if(!items.length)return;state.selected=items[last?items.length-1:0].k;state.noSelection=false;writeRoute(true);list();const active=el('list').querySelector('[data-key="'+CSS.escape(state.selected)+'"]');if(active)active.focus()}
const help=el('keyboard-help');
function editableTarget(target){return !!target.closest('input,textarea,select,[contenteditable="true"],[role="textbox"]')}
window.addEventListener('popstate',()=>{readRoute();render()});window.addEventListener('hashchange',()=>{readRoute();render()});
document.addEventListener('keydown',e=>{if(help.open){if(e.key==='Escape'){help.close();e.preventDefault()}return}if(editableTarget(e.target)||e.altKey||e.ctrlKey||e.metaKey)return;if(e.shiftKey&&e.key!=='?'&&e.key!=='G')return;if(e.key==='?'){help.showModal();el('close-help').focus();e.preventDefault();return}if(e.key==='/'){el('filter').focus();e.preventDefault();return}if(e.key==='j'||e.key==='ArrowDown'){moveSelection(1);e.preventDefault();return}if(e.key==='k'||e.key==='ArrowUp'){moveSelection(-1);e.preventDefault();return}if(e.key==='g'){moveToEdge(false);e.preventDefault();return}if(e.key==='G'){moveToEdge(true);e.preventDefault();return}if(e.key==='Enter'){if(e.target.closest('a,button,summary,input,textarea,select,[role=button],[role=link]'))return;el('detail').focus();e.preventDefault();return}if(e.key==='Escape'){if(el('fix-review').open){el('fix-review').close();state.batchPreview=null;e.preventDefault();return}if(history.state&&history.state.markataServeControl&&(history.state.depth||0)>0){history.back();e.preventDefault();return}state.selected='';state.noSelection=true;writeRoute(false);list();e.preventDefault();return}if(e.key==='r'&&state.section==='jobs'&&capabilities.rerun){const selected=visibleItems().find(item=>item.k===state.selected)?.item;if(selected&&selected.state!=='queued'&&selected.state!=='running'){action({kind:'rerun',job_id:selected.id});e.preventDefault()}return}if(e.key==='w'){navigate('diagnostics','','','warning');return}if(e.key==='e'){navigate('diagnostics','','','error');return}if(e.key==='p'){navigate('pages');return}if(e.key==='f'){navigate('feeds');return}if(e.key==='l'){navigate('logs');return}});
document.addEventListener('click',e=>{const t=e.target.closest('[data-section]');if(t){navigate(t.dataset.section);return}const select=e.target.closest('[data-key]');if(select){selectItem(select.dataset.key);return}const fixGroup=e.target.closest('[data-fix-group]');if(fixGroup){reviewFixes(fixGroup.dataset.fixGroup);return}const fixOne=e.target.closest('[data-fix-one]');if(fixOne){reviewFixes(fixOne.dataset.fixOne);return}const fixSelect=e.target.closest('[data-fix-select]');if(fixSelect){const key=fixSelect.dataset.fixSelect;if(state.selectedFixes.has(key))state.selectedFixes.delete(key);else state.selectedFixes.add(key);render();return}const rerun=e.target.closest('[data-rerun]');if(rerun){action({kind:'rerun',job_id:rerun.dataset.rerun});return}const rebuild=e.target.closest('[data-rebuild]');if(rebuild){action({kind:'rebuild-page',page:rebuild.dataset.rebuild});return}const logs=e.target.closest('[data-job-logs]');if(logs){navigate('logs',logs.dataset.jobLogs);return}const page=e.target.closest('[data-page]');if(page){navigate('pages',page.dataset.page,'page:'+page.dataset.page);return}const job=e.target.closest('[data-job]');if(job){navigate('jobs',job.dataset.job,'job:'+job.dataset.job)}});
el('fix-review').addEventListener('close',()=>{cleanupReviewDiffs();state.batchPreview=null});el('fix-cancel').addEventListener('click',()=>el('fix-review').close());el('fix-approve').addEventListener('click',applyBatch);el('close-help').addEventListener('click',()=>help.close());help.addEventListener('click',e=>{if(e.target===help)help.close()});
el('copy-report-button').addEventListener('click',copyProblemReport);el('build-button').hidden=!capabilities.build;el('build-button').addEventListener('click',()=>action({kind:'build'}));el('filter').addEventListener('input',e=>{state.filter=e.target.value;state.severity='';state.noSelection=false;writeRoute(false);list()});if(!history.state||!history.state.markataServeControl)history.replaceState({markataServeControl:true,depth:0},'',location.href);readRoute();writeRoute(false);render();poll();setInterval(poll,2000);
})();
</script></body></html>`
