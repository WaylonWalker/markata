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
`source_lines`, `source_anchors`, `render_ms`, `input_bytes`, `html_bytes`,
`document_bytes`, and `cached`. `source_lines` maps top-level rendered blocks;
`source_anchors` maps nested container starts and ends to editor lines.
The optional `include_document` boolean defaults to true for API compatibility.
With false, `document` is null and `document_bytes` is zero. Fragment responses
contain no shared appearance CSS or document shell.

The editor loads its sandboxed preview document once, then sends fragment updates
to the same frame. Shared styles and loaded script dependencies persist across
edits. Conditional styles are reconciled without duplicates, old Chart.js
instances are destroyed, and generated initialization scripts run in order.
The initial preview MUST report content-ready after the article and bridge are
available, without waiting for deferred feature initialization. It MUST accept
latest-wins fragment updates after blocking preview scripts are parsed, without
waiting for deferred module initialization or the iframe load event. Diagnostics
MUST distinguish initial frame content readiness, script readiness, and the
iframe load event from later iframe patch timing.
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
`GET /favicon.ico` serves the bundled Markata favicon so browsers do not report
a missing icon for the workbench origin.

The Go workbench serves its Excalidraw editor and fonts from the embedded local
asset set. The upstream editor requests font files below an additional `fonts/`
URL segment; `/editor/excalidraw-fonts/fonts/<path>` MUST resolve to the same
bundled asset as `/editor/excalidraw-fonts/<path>`. Split editor JavaScript
chunks MUST also be served below `/editor/excalidraw-editor-chunks/`. These
paths MUST remain contained in the embedded asset filesystem. This keeps text
editing and SVG export offline.

The Go workbench MUST keep YAML frontmatter out of the rendered Markdown while
using valid frontmatter to preview the post title and description. Link and
bookmark templates MUST preview their `link` or `url`, title, description,
image, site/provider, media type, and author metadata as an external link card
when those fields are available. Photo and shot templates MUST preview their
frontmatter image above the body. Absolute HTTP(S) image URLs MAY load directly;
root-relative `/images/` paths MUST resolve only to regular files under the
workspace `static/` directory. Preview rendering MUST NOT fetch page metadata
or modify source. Invalid YAML MUST leave the Markdown body preview available
and show that frontmatter metadata could not be rendered.

Pasting a PNG, JPEG, GIF, or WebP image into the Markdown editor MUST upload it
through the workbench's Dropper integration and insert an image Markdown link at
the selection. The editor MUST replace the selection immediately with a unique,
accessible upload status placeholder. On success, it MUST replace only that
placeholder with the image Markdown. On failure, it MUST restore the selected
text when the placeholder remains. If the author changes documents while the
upload is pending, completion MUST NOT insert Markdown into the newly active
document; the originating placeholder MUST be resolved if that document is
opened again. If the author removes the placeholder, the editor MUST preserve
the current source and report the uploaded Markdown so it can be recovered.

`POST /workspace/dropper-upload` MUST accept one multipart `image` file only
for same-origin loopback workspace requests. It MUST enforce the existing
15 MiB limit and verify that the declared MIME type matches PNG, JPEG, GIF, or
WebP image data before contacting Dropper. Dropper credentials MUST remain
server-side. Uploads MUST be disabled when the effective
`content_templates.image_upload` provider is `local`. A successful response
MUST return the validated Dropper image URL. It MUST NOT write a workspace file.


Workspace reads require a loopback client and a localhost or loopback Host.
`PUT` and `POST /workspace/file` also require an exact same-origin `Origin`;
when `Sec-Fetch-Site` is present, it must be `same-origin`. Missing, malformed,
or foreign origins are rejected with HTTP 403 before the workspace is changed.

Workspace metadata indexing MUST parse large uncached file sets with the
loader's bounded worker pool. Small file sets MAY load sequentially to avoid
worker overhead. Index construction remains asynchronous to editor startup,
and one workspace store MUST reuse parsed entries for unchanged files for the
life of its server process.

The Git Changes view MUST request a fresh working-tree inventory from
`GET /workspace/changes` when the editor attaches, so its badge is ready before
the view opens. The browser MUST retain the last inventory in memory, render it
immediately when opening the view, and refresh it in the background. Successful
workspace saves, window focus, manual refresh, and the open-view polling interval
MUST trigger revalidation. The server MUST report current Git status rather than
serving a long-lived cached snapshot. The active-post Review action MUST be
located in the Changes panel header.

The Changes view MUST provide a review queue that distinguishes Markdown posts
from other files, supports filtering by path, and exposes each entry's Git
status, staged/unstaged state, and workspace scope. Selecting a supported text
file MUST show its saved workspace source compared with the last committed Git
baseline in a readable diff. The UI MUST label Git HEAD as the last commit and
MUST NOT imply that it is deployed or live. Deleted files, conflicts, binary or
oversized files, and files outside the selected site MUST keep a clear status
and explanation when a source diff is unavailable.

For an eligible Markdown post, the review queue MUST provide a clear path to
the existing rendered, side-by-side, source, and full-diff review modes. The
quick source view MUST remain readable without the optional full-diff renderer.

The review queue's **Keep in workspace** action MUST leave the file, Git index,
and history unchanged. It only marks the entry reviewed for the current queue
session and MUST be described as such. It MUST NOT stage, commit, push, publish,
or hide the change from the Git inventory. **Review and reset** MUST open the
existing single-document review for an eligible workspace Markdown file. Reset
MUST remain unavailable for staged changes, deleted files, conflicts, and
out-of-site or non-Markdown files; the UI MUST explain why. The existing
confirmation, exact revision checks, and stale-draft refusal MUST remain in
force. Reset MUST NOT affect unrelated paths or Git history.

The review queue MUST keep navigation usable with keyboard and screen readers.
On narrow screens, the file list and selected diff MUST stack without hiding the
selected path, status, or decision actions. Closing the queue MUST preserve all
workspace changes. Progress is session-only and MUST reset when the user
requests a refresh or Git status changes. Background revalidation with the same
inventory MUST preserve the current session's decisions.
Workspace navigation MUST use stable tabs for Recent, Files, Feeds, Tags, and
Changes. Recent files MUST occupy their own panel rather than appearing inside
Files. The tablist MUST support arrow-key navigation and expose the selected tab
and panel accessibly. Switching tabs MUST NOT resize or move the sidebar shell.
Unavailable Feeds and Tags tabs MUST remain disabled until their data is ready,
or when the site has no matching entries.

The Preview speed setting MUST support Auto, Fast, Normal, and Slow, and MUST
persist in browser storage. Normal is the default and waits 50 ms after the
latest edit, with a 250 ms maximum delay during continuous typing. Slow waits
300 ms and starts a preview within 1,000 ms during continuous typing. Fast starts
on each editor change without a debounce. At most one render request may be
active at a time, and the client MUST keep at most one pending revision. Newer
edits MUST replace pending work, and stale revisions MUST NOT be mounted in the
preview. The iframe MUST keep at most one active update and one replaceable latest
update, and MUST discard stale queued revisions. Auto MUST use Normal by default.
It MUST select Slow when the browser reports data saver, a 2G connection, four
or fewer logical CPU cores, 4 GiB or less device memory, or a discharging battery
at or below 20%, when those signals are available. Auto MUST NOT select Fast.
Performance Diagnostics MUST show the latest preview as a four-stage timeline:
waiting (including debounce), server queue plus Go rendering, HTTP/response
overhead, and preview display. HTTP/response overhead is the browser-measured
request duration minus server queue and render timings; it includes server
response encoding and browser JSON decoding, and MUST NOT be described as pure
network time. The display stage MUST include iframe patch timing when available,
including HTML parsing, DOM reconciliation, feature setup, scroll-map collection,
scroll-map ordering, and caret-marker synchronization. Nested scroll anchors MUST
map to their top-level preview block in linear time. The copied local diagnostics
report MUST include the latest stage values and iframe details without source
content or paths.
Content patches preserve unchanged preview blocks. When the editor has focus,
the actual caret line remains the synchronization target through the completed
patch, including edits that leave the caret on the same source line. Its preview
marker is refreshed against the new source map. Draft and keymap persistence are
best effort; storage failure is visible and does not prevent editing or
downloading Markdown. Workbench-only source-line attributes on Markdown
containers provide additional scroll anchors without changing the editor text
or its line count. For nested content, the preview marker selects the narrowest
container span containing the caret. Source ordered mapping handles side-by-side
content; preview scrolling uses visual order and preserves the previously active
source branch when several blocks share one pixel row. Appearance updates do not
re-render Markdown. The server
does not store drafts.

Automatic LSP completion waits for a short typing pause instead of requesting a
completion after every character. Ctrl+Space requests completion immediately.

Workspace autosave compares each write against the exact source version the
browser last read or successfully wrote. A stale version MUST return HTTP 409
without changing the file, unless the current file already contains the exact
requested bytes; that idempotent retry MUST return success with the current
version and MUST NOT write the file again. After a 409, the browser MUST retain
its recovery draft and stop automatic writes for that file until the user
resolves the conflict. **Compare copies…** shows the browser draft and the
current disk copy. Keeping the browser draft requires an explicit confirmation
and writes against the reviewed disk version; using the disk copy replaces the
editor buffer. Either choice updates the expected version before autosave
resumes.

The demo binds to 127.0.0.1 by default. It rejects cross-origin render requests
and non-JSON content types. It has HTTP timeouts and graceful shutdown.
The handler is reusable by a future builder; it does not depend on CLI globals.
On desktop, the editor shell fills the viewport without document-level scrolling;
the workspace panes own their scroll areas. Narrow mobile layouts may scroll the
stacked editor and preview vertically.
In split layout, the preview marks the rendered block for the active source
cursor and shows a compact labeled caret. When source text matches rendered
text, the preview caret SHOULD map the editor's character offset to a rendered
text range. For transformed or replaced content, it MUST fall back to the
nearest source-mapped position. In Vim mode, the presence caret MUST reflect
normal, insert, replace, and visual cursor shapes. Presence markers use
participant IDs, labels, and colors so additional editor cursors can be added
without changing the marker rendering model. Clicking rendered text in the
preview MUST move the source cursor to the matching source position when a
source match exists; links MUST retain their normal navigation behavior.
When the local cursor follows trailing spaces that Markdown collapses, the
preview SHOULD show subtle whitespace markers across the collapsed spaces and
place the caret after the spaces before the source cursor. This cue MUST NOT
change rendered Markdown or apply to non-trailing spaces.


## Integration benchmark

The benchmark reuses the seven Python-reference fixture files without changing
source bytes. It records fixture and binary hashes, raw samples, environment,
and source revisions. Initial and fragment cache misses are separate from
repeated cached fragments. A separate unmeasured render invalidates the one-entry
cache before each miss sample; every measured response validates its cache flag.
Each mode excludes two warm-up calls and measures fifteen calls by default.
HTTP timing excludes JSON encoding/decoding, binary startup, and asset requests.
Reports identify linked Go assets and inline Python assets as different transport
policies. Comparisons against the stored Python report require matching fixture
hashes and identify the measurements as separate runs on a shared workstation.
