# Go render workbench

The standalone Go workbench provides an embeddable `http.Handler` and a loopback
HTTP demo. Rendering calls `plugins.RenderStandalone` directly in process;
there is no Python subprocess, WASM runtime, site build, or rendering network call.

`GET /` serves the editor. `POST /render` accepts JSON containing `markdown` and
optional appearance attributes. Source is limited to 200,000 UTF-8 bytes; the
whole request is limited to 1 MiB. Invalid JSON and appearance options return
400; oversized requests return 413. Render failures return 422. Only one render
runs at a time per handler; queue admission observes request cancellation.

Responses contain `document`, `html`, `feature_head_html`, `body_end_html`,
`render_ms`, `input_bytes`, `html_bytes`, `document_bytes`, and `cached`.
The optional `include_document` boolean defaults to true for API compatibility.
With false, `document` is null and `document_bytes` is zero. Fragment responses
contain no shared appearance CSS or document shell.

The editor loads its sandboxed preview document once, then sends fragment updates
to the same frame. Shared styles and loaded script dependencies persist across
edits. Conditional styles are reconciled without duplicates, old Chart.js
instances are destroyed, and generated initialization scripts run in order.
The initial document includes each shared stylesheet once. Feature styles are
marked for replacement and removal. The bridge accepts updates only from its
parent; it does not grant same-origin access or replay DOMContentLoaded.
Asset URLs are normalized against the document base before deduplication.
Sortable tables and other features initialize new elements after updates.
Authored HTML and styles persist in fragments; only generated scripts rerun.
Reloading the editor applies changes to authored scripts.
The preview uses the shared standalone appearance assets and ordered renderer
assets. The handler permits raw authored HTML for local trusted editing. The
browser isolates preview scripts in an opaque-origin iframe sandbox. All assets
are served locally. Authored images and links may access external networks.

The browser debounces edits, runs at most one request at a time, coalesces pending
edits, and rejects stale revisions. Draft and keymap persistence are best effort;
storage failure is visible and does not prevent editing or downloading Markdown.
Appearance updates do not re-render Markdown. The server does not store drafts.

The demo binds to 127.0.0.1 by default. It rejects cross-origin render requests
and non-JSON content types. It has HTTP timeouts and graceful shutdown.
The handler is reusable by a future builder; it does not depend on CLI globals.
