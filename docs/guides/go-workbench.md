---
title: "Live Markdown editor with Go"
description: "Run or embed a Go server with a live Markata Markdown editor"
date: 2026-11-20
published: true
tags:
  - documentation
  - go
  - markdown
---

Run the shipped Go editor command:

```bash
markata editor --workspace ./docs
```

Open the local URL printed by the command. Edit Markdown on the left and view the preview on the right.
Press Ctrl+C to stop the server. The default address is `127.0.0.1:8766`.
For a checkout build, run `go run ./cmd/markata-go editor --workspace ./docs`.

The Go server calls Markata's standalone renderer directly. It starts no Python process and needs no WASM runtime.
It uses the same Markdown renderer as the [Python API](./python-renderer/) and [FastAPI editor](./fastapi-editor/).
The preview does not build a site. The language service indexes Markdown and reads workspace configuration for existing mention support.

### Run behind Hlab Auth

The editor binds to loopback by default. Keep that setting for local work.
Hosted mode requires an HTTPS public origin, the reverse proxy's source CIDR,
and an app-specific Hlab writer scope:

```bash
markata editor \
  --workspace /data/source \
  --listen 0.0.0.0:8766 \
  --public-origin https://edit.example.com \
  --trusted-proxy-cidr 10.42.0.0/24 \
  --required-scope markata-editor:admin
```

Use the exact proxy source ranges enforced by the cluster NetworkPolicy. The
server checks the TCP peer address, public Host, stable `X-Hlab-User-Id`, and
the exact scope on every request, including assets, workspace APIs, and LSP
WebSockets. It does not use `X-Forwarded-For` or `X-Hlab-Auth-Method`.
Traefik must strip client-supplied `X-Hlab-*` headers before ForwardAuth and
forward only the identity headers returned by Hlab Auth. Restrict the editor
service to Traefik in the NetworkPolicy.

Register a dedicated Hlab app for the editor. Require its administrative scope
at the editor so a viewer scope cannot save files, upload images, or publish.
Do not reuse the broader Build Admin scope. The proxy route must use Hlab
ForwardAuth and return the `__Host-hlab-app-session` cookie for browser login
handoff. Hosted mutations also require the exact configured public Origin, and
hosted responses are marked private and no-store.

The rendered preview keeps its opaque-origin script sandbox. The workbench
delivers its appearance styles and selected font files inside that preview, so
they do not need a separate Hlab login handoff.

Do not expose the listener directly to the internet. Non-loopback binding is
rejected unless all hosted authentication settings are valid.

The Markata Helm chart can run the editor as a sidecar beside Builder Admin. It
shares the site's source checkout, but uses a separate authenticated ingress:

```yaml
builderAdmin:
  editor:
    enabled: true
    requiredScope: markata-editor:admin
    ingress:
      host: edit.example.com
      ingressClassName: traefik
      tls:
        secretName: edit-example-com-tls
    gitAuth:
      repository: https://github.com/owner/site
      branch: main
      existingSecretName: markata-editor-git-publish
      privateKeyKey: ssh-privatekey
      knownHosts: |
        github.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMqqnkVzrm0SdG6UOoqKLsabgH5C9okWi0dh2l9GKJl
```

The chart reuses the Builder Admin ForwardAuth service URL and trusted proxy
CIDRs. It creates a separate middleware, service, ingress, and NetworkPolicy
rule for the editor port. Create the TLS secret with the cluster's certificate
controller. Register a dedicated Hlab app for `edit.example.com`; do not add
this host to the Build Admin app or grant its operator scope to the editor.
The GitHub deploy key must be write-enabled for only the site's repository;
store its private key in a namespace-scoped SealedSecret and use the verified
GitHub host key shown above for strict SSH host checking. An init container
copies the Secret into an ephemeral volume with owner-only file permissions;
the Git publisher sidecar is the only long-running container that mounts that
volume. It uses the configured repository and branch as an allowlist, and
accepts requests only over a local Unix socket. The editor has no private-key
mount, SSH command, or agent socket; it continues to require the authenticated
review, commit, and push flow.

Hosted deployments can enable `builderAdmin.editor.dropperUploader` and point
it at a namespace-scoped Secret containing the Dropper token. The chart runs a
separate uploader sidecar that exposes only validated image uploads over a Unix
socket. Only that sidecar receives the token; the editor container and browser
do not. The broker accepts PNG, JPEG, GIF, and WebP images up to 15 MiB and
forwards them to its fixed Dropper upload endpoint. Configure the chart values:

```yaml
builderAdmin:
  editor:
    dropperUploader:
      enabled: true
      existingSecretName: markata-editor-dropper
      tokenKey: token
      endpoint: https://dropper.waylonwalker.com/api/upload?format=json
```

Use the [Excalidraw canvas](excalidraw-editor.md) to create and edit drawings
inside a Markdown page.

### Paste an image into Markdown

When the site's effective `content_templates.image_upload` setting uses
Dropper, paste a PNG, JPEG, GIF, or WebP image into the Markdown editor. The
workbench places an **Uploading image…** marker at the cursor while it uploads,
then replaces that marker with an image link such as `![image](https://…)`.
The workbench uses a server-side Dropper credential; the browser never receives
it. In hosted deployments, the credential stays in the uploader sidecar and is
not available to the editor process. The upload limit is 15 MiB.

Paste while the Markdown editor has focus. If the browser reports image data but
does not make the image bytes available, the workbench shows an error in the
editor toolbar. Try copying the image again or drop the image file into the
editor.

With the Vim keymap, press `p` or `P` in normal or visual mode to upload an image
from the system clipboard. If the clipboard has no image, those keys keep their
usual Vim register-paste behavior. `Ctrl+V` remains Vim's blockwise selection
command.

If an upload fails, the workbench restores the selected text when the marker is
still present. If you remove the marker while an upload is running, the editor
does not overwrite your edits and reports the Markdown returned by Dropper.
Uploads are unavailable when `content_templates.image_upload` resolves to
`local`. See the [content starter guide](content-starters.md) for provider
configuration.

The preview shows its rendered article and accepts updates before deferred
feature initialization finishes. Mermaid and other preview features can finish
initializing afterward. Diagnostics report script readiness and full iframe
load separately from later preview patch timings.

### Understand startup and site index timing

The editor startup has distinct milestones. The CLI prints elapsed time for
workspace discovery/watcher setup, LSP preparation, and appearance/browser
assets. Once the HTTP listener is ready, the browser still loads its source
file list and lazily requests frontmatter plus the site link index; these are
**not** part of the initial server-ready message.

Filesystem watchers prune directories matched by the workspace's `.gitignore`
rules. This keeps generated output and cache trees out of editor startup while
preserving watches for discovered content and the workspace root.

Large uncached metadata loads use a bounded worker pool. The store then reuses
unchanged metadata while that server process stays open. A new server process
builds its own index, so large sites still incur a one-time indexing cost per
process; this work runs in the background and does not block writing or preview.

Open **Diagnostics** in the footer to see **Loading site index…** and its phases:
waiting for metadata, reading frontmatter, constructing the link index, and
ready/error. Expand **Site index ready** afterward to see file counts and
the measured index load time (frontmatter parsing and link construction
separately), configured feed assembly and metadata HTTP timing, browser file
list readiness, and time to the first usable preview. **Copy diagnostics** captures a local JSON snapshot of those
counts and durations, not document contents.

For a local command-line check:

```bash
curl -s http://127.0.0.1:8766/workspace/index-status
curl -sD - -o /dev/null 'http://127.0.0.1:8766/workspace?metadata=true'
```

The first command is lightweight and does not initiate indexing. The second
explicitly loads metadata and returns a `Server-Timing: site-index;dur=…`
header (and a `site-feeds;dur=…` value for configured feed assembly). This is local diagnostic telemetry, not a distributed tracing system;
there is no automatic upload of timing data. If it exposes a slow stage on
a large site, the next step is to profile that stage rather than guess at a
percentage-complete bar. Related: #1987.

Use the appearance picker to select a palette, color mode, style, font, or text size.
Open [Editor settings](./editor-workspace/) to choose line wrapping, keymap, and layout. Drag the
divider between the source and preview panes to resize them. Choose **Editor
only** or **Preview only** to focus on one pane; **Editor and preview** restores
the split.
Scrolling either pane keeps the corresponding Markdown visible in the other.
When you write or move the cursor, the preview follows the editor. Scroll sync
is on by default; turn off **Sync editor and preview scrolling** in **Editor
settings** to scroll each pane independently. The choice is saved in this browser.
At the end of a document, both panes keep extra scroll space so the last lines
can move above the bottom edge while you write.
In split view, a subtle marker shows the corresponding preview block, with a
thin caret and **You** flag at the active location. When rendered text matches
the source line, the preview caret follows the character cursor. For content
that a renderer transforms or replaces, it falls back to the nearest mapped
source position. If the cursor follows trailing spaces that Markdown hides,
small dots show the spaces and the preview caret advances past them. Click
rendered text to move the source caret to the matching text when possible;
links keep their normal navigation behavior. The source
pane keeps its normal caret on the actual text
cursor. With the Vim keymap, the preview caret changes shape for normal,
insert, replace, and visual modes. Preview updates preserve
unchanged blocks, and layout changes refresh the source-to-preview mapping,
including after images and fonts load. While you edit, the preview follows the
caret after the matching content update completes, including when you keep
typing on the same line. Scroll changes caused by a render or a delayed
image/font layout are not mistaken for manual scrolling. The preview marker is
subtle and marks the rendered block for the current caret line. Nested Markdown
container boundaries add source anchors to improve scroll alignment in long
sections; these markers exist only in the preview render and do not change your
source. Within nested content, the marker follows the narrowest container that
contains the caret line. Generated frontmatter previews, including link cards,
do not shift the marker away from the matching Markdown block.
For the Go site-backed workbench, **Editor settings → Keyboard Shortcuts** opens
application-wide key bindings. Choose Writer, VS Code, Vim, or Emacs application
presets independently of the editor's Default/Vim/Emacs/Sublime text keymaps.
The command palette also exposes **Editor keymap: Default**, **Vim**, **Emacs**,
and **Sublime**. Selecting one uses the existing editor keymap preference;
it does not change the independent Writer/VS Code/Vim/Emacs application
shortcut preset. This lets writers change editing motions without navigating
to Settings.
The command palette displays active shortcuts; the registry also supports
multi-key chords, search, per-preset overrides, and resets. The preview pane
includes a **Back to editor** control for keyboard navigation.
See [Keyboard shortcuts](./editor-keybindings/).

Choose Default, Vim, Emacs, or Sublime keymaps in **Editor settings**.
In Vim mode, press `/` to search forward, `?` to search backward, then `n` or
`N` to move between matches.
Code highlighting, icons, admonitions, tables, checklists, Mermaid, and Chart.js use local assets.

Press **Ctrl+P** (**Cmd+P** on macOS) to find a site file by its path, title,
slug, or tag. The finder shows the current and recently opened files when its
search is empty.
Select a result or press **Enter** to open it in the editor.
The sidebar has tabs for **Recent**, **Files**, **Feeds**, **Tags**, and
**Changes**. **Recent** lists the last eight workspace documents opened on this
device. It is saved per workspace in browser storage, works without network
access, and excludes unsaved scratch drafts. Select a recent entry to reopen it;
files no longer in the workspace are omitted. **Feeds** and **Tags** browse
posts by feed or tag. Feed names are shown as roots with folders relative to the
workspace.

Edits to selected site files save back to the workspace after a short pause.
Unsaved scratch documents stay in browser storage. The browser also keeps a
recovery copy of workspace edits. Use **Ctrl+Shift+P** and search for **Download
current file** to export a copy, or **Open Markdown file** to import a file as
an unsaved draft.
A storage error appears beside the render metrics. Editing and downloads remain available.
The server writes only selected workspace files. Scratch drafts remain local to
the browser origin, so drafts on one port remain separate from drafts on another.
Workspace saves and new-post creation require a same-origin browser request.
The server rejects missing or foreign origins before it changes a site file.

Choose **Preview speed** in **Editor settings → Preview**:

- **Auto** uses Normal by default. It selects Slow when the browser reports data saver,
  a 2G connection, four or fewer logical CPU cores, 4 GiB or less device memory, or a
  discharging battery at or below 20%, when those signals are available.
- **Fast** requests a preview on every editor change. If a render is in progress,
  newer edits replace pending work and stale revisions are skipped.
- **Normal** waits 50 ms after the latest edit and starts an update within 250 ms
  during continuous typing. This is the default.
- **Slow** waits 300 ms after the latest edit and starts an update within 1,000 ms
  during continuous typing.

The setting stays in this browser. Auto never selects Fast. Slow reduces how often
the editor sends work; it does not change the size of each render request.
The first preview can appear while the site index is still loading. Workspace links and embeds
resolve when indexing finishes, using the current editor text.
New edits can update the first preview as soon as its page script is ready, even if an image or
font in that page is still loading.
Edits made during indexing update the changed post without restarting metadata loading for
every other post.
Old responses cannot replace a newer draft. Errors preserve the last successful preview.
Rendered Mermaid diagrams stay in place when you edit other text, so a large post does not redraw
every diagram after each keystroke. Diagrams below the visible page render as you scroll near them.
Changing a diagram's source redraws that diagram when it is near view.
Mermaid and Chart.js fences skip generic code syntax detection before their own renderers process
them, which keeps long posts with many diagrams responsive while typing.
The footer shows Go render time, total edit-to-preview time, source size, and preview size.
Open **Diagnostics** to see a visual breakdown of the latest preview: waiting
(including the edit debounce), server queue and Go rendering, HTTP/response
overhead, and display in the preview frame. The display details include HTML
parsing, DOM updates, and feature setup when available. The HTTP/response stage
is the full request time after measured server queue and rendering, so it also
includes response encoding and browser JSON decoding. The display detail also
shows scroll-map collection and ordering plus caret-marker synchronization.
**Copy full report** adds these stage timings to the local diagnostics snapshot
without including Markdown content or file paths.
Repeated source uses a one-entry memory cache. Appearance changes update the preview without a render request.
Workspace post embeds in `![[post-slug]]` or `![embed](path/to/post.md)` format
render as cards linked to the matching site post. Fenced code examples remain
literal source.

The preview uses valid frontmatter without displaying the YAML itself. It shows
the title and description, renders Link and bookmark metadata as a link card,
and places Photo and Shot images above the body. Remote HTTP(S) images preview
directly; `/images/` paths load from the site's `static/` directory. If the
frontmatter YAML is invalid, the Markdown body still previews and the frame shows
that it could not use the page metadata.

## Language assistance

The editor connects to the existing Markata LSP through `/lsp` over WebSocket.
The server indexes the directory from `--workspace`. The default directory is the current directory.
The browser LSP accepts same-origin connections on `localhost` or a loopback IP only.

Type `[[` for wiki-link suggestions, including aliases. Type `[label](` for Markdown document paths.
Automatic suggestions wait for a short typing pause. Press **Ctrl+Space** to request suggestions immediately at the cursor. Use the arrow keys and Enter to select a suggestion.
Existing frontmatter, admonition, and mention completion also use the LSP.
Broken wiki links receive an underline. The underline tooltip shows the diagnostic message.
The workbench preview resolves valid `[[wikilinks]]` against the site's workspace index.
Click a resolved wikilink in the preview to open that post in the editor.
Select the error and warning counts in the bottom status bar to inspect diagnostics and
jump to their source locations.

The editor sends unsaved text as a synthetic `.markata/browser-draft.md` document.
This document stays in memory. Markdown path suggestions start at the workspace root because the draft has no source file.
The Files tree opens source files from the selected workspace using Markata glob rules. The Feeds view groups files by configured feed filters, and the Tags view groups them by loaded frontmatter. The selector shows file, feed, and tag counts. See [workspace navigation](./renderer-site-tree/).
Each WebSocket connection owns its own draft and index. Reload the page to index external file changes.

## Edit page properties

Select a site file, then choose **Properties** to edit its frontmatter in the
full-height page-details view. The live preview stays visible beside it. Choose
**Back to Markdown** to return to the source, or **Edit YAML** to work directly
with the frontmatter block. Press **Escape** to return to the article. The
command palette also has actions to open page details, add a property, and
return to Markdown.

The page summary shows the title, URL slug or source path, and publication
state. Common fields appear first. Choose **Add property**, search for a known
field, and select it.
Open **Custom property** to add a site-specific YAML field. Edit tags as chips:
remove one tag at a time or enter a tag and choose **Add**. Advanced fields
remain available under **Show advanced properties**.

Property edits update the Markdown source and use the normal workspace save
and preview behavior. Press **Ctrl+Z** or **Cmd+Z** to undo an edit. Use
**Remove** to delete a property from the page.

The footer reports the language connection state. Editing and preview continue if the connection fails.
Reload the page to reconnect. The browser adapter is shared with the [FastAPI example](./fastapi-editor/).

For an embedded application, use `workbench.NewWithWorkspace(ctx, root)`.
Cancel `ctx` during shutdown to close active language connections.
`workbench.New()` uses the current directory and the lifetime of the application.

### WebSocket protocol

`GET /lsp/config` returns `rootUri` and `documentUri` for the selected workspace.
Connect to `/lsp` with the page origin. Each WebSocket text message contains one JSON-RPC object.
WebSocket messages omit the stdio `Content-Length` headers.
Send `initialize`, `initialized`, and `textDocument/didOpen` before language requests.
The session accepts the synthetic draft URI from the configuration response.
The server ignores browser workspace-change notifications and overrides the client workspace during initialization.

Browser completion responses contain at most 200 suggestions. More specific text requests a new list.
Messages have a 2 MiB limit. A disconnect releases the session.
The endpoint rejects missing or foreign origins. It provides no authentication or shared workspace management.
Keep this example on a trusted local address.

## Embed the handler

Mount the handler at the root of a dedicated origin:

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/WaylonWalker/markata/pkg/workbench"
)

func main() {
    handler, err := workbench.New()
    if err != nil {
        log.Fatal(err)
    }
    server := &http.Server{
        Addr: "127.0.0.1:8767", Handler: handler,
        ReadHeaderTimeout: 5 * time.Second,
    }
    log.Fatal(server.ListenAndServe())
}
```

Create one handler at startup. It prepares the shared appearance bundle once.
The handler owns no listeners, signal handlers, or application configuration.
A future builder can embed it without the demo command.

## Render API

```bash
curl http://127.0.0.1:8767/render \
  -H 'Content-Type: application/json' \
  -d '{"markdown":"# Hello Go","appearance":{"data-theme":"light"}}'
```

The first editor request receives a complete preview document. Later requests use `"include_document": false` and receive fragments.
The iframe keeps its document, shared styles, and loaded libraries between edits.
New features add their styles once; removed features remove their conditional styles.
Chart instances are destroyed before replacement, and feature scripts initialize the new content.

The JSON response contains `document`, `html`, `feature_head_html`, `body_end_html`, and `cached`.
Metrics include `render_ms`, `input_bytes`, `html_bytes`, and `document_bytes`.
With `include_document` set to false, `document` is null and `document_bytes` is zero.
Omit this option to retain the original full-document response.
The `Server-Timing` response header also reports render time.
`document` contains a complete HTML preview with links to this server's assets.
Keep the server available while you view the document.

Optional appearance keys are `data-theme`, `data-palette`, `data-aesthetic`, `data-fontpack`, and `data-text-size`.
Values must occur in the bundled appearance catalogs.
The API rejects unknown fields, malformed JSON, and invalid appearance choices with HTTP 400.
It returns HTTP 413 for source above 200,000 UTF-8 bytes or a request above 1 MiB.
Requests need `application/json`. Cross-origin render requests receive HTTP 403.

The handler runs one render at a time. Canceled requests leave the queue before rendering.
An active render completes before the next render starts.
The browser stops waiting after 30 seconds. The demo server has read and write timeouts.

This editor supports trusted local authoring. Authored HTML and styles remain available inside an opaque-origin iframe sandbox.
Only renderer-generated scripts run during fragment updates. Reload the page to apply changes to authored scripts.
Preview scripts cannot read the editor's DOM or storage. Authored links and images can access external networks.
Keep the default loopback binding for local use. Add authentication and suitable resource controls before remote deployment.

All built-in assets come from the server. Offline reload is not supported.
For site builds and site-specific plugins, use the normal [CLI](../reference/cli/).

## Check persistent preview updates

Start the Go workbench, then run its browser regression tests:

```bash
MARKATA_GO_WORKBENCH_URL=http://127.0.0.1:8767 \
MARKATA_TEST_CHROMIUM=/usr/bin/chromium \
uv run --with playwright python python/test_go_editor_browser.py -v
```

Set `MARKATA_TEST_CHROMIUM` to your Chromium executable.
The tests check retained stylesheets, fragment responses, conditional styles, charts, and diagrams.

The publishing browser test uses its own temporary Git site and bare remote. It
edits page properties, reviews the rendered and source diffs, commits an edited
message, and pushes only to that temporary remote:

```bash
go build -o /tmp/markata-go-editor ./cmd/markata-go
MARKATA_GO_EDITOR_BINARY=/tmp/markata-go-editor \
uv run --with playwright python -m unittest python.test_go_editor_publish_flow -v
```

The pull request fast gate runs this test when Python browser-test files change.
It installs Chromium and uses the same disposable Git fixture.


## Benchmark the server

The [integration benchmark guide](/docs/guides/renderer-integration-benchmarks/#go-server-workbench) compares seven existing fixtures with the Python reference.
It separates initial renders, uncached edits, and cached edits. Raw samples and reproduction commands are included.

## Shared editor behavior

Use the command palette to import or export a Markdown file. An oversized import leaves the current document intact.
Search for **New Post** in the command palette to open a fuzzy picker for
built-in starters and `.md` files in the site's `templates/` directory. Choose
a starter to open its details form. The title and destination folder start with
template defaults. The slug field stays empty when using the suggestion; its
ghost text and the destination path update as the title changes. Enter a custom
slug only when you want to override the title-derived filename. Creating the
post opens the new file in the editor, where later edits autosave to the site.
The built-in **Link** starter accepts a source URL and can fill in the title,
description, and preview image from the same configured oEmbed and Open Graph
providers used by [external embeds](./embeds/). It also imports available site,
provider, media type, and author metadata. Enter a URL and pause briefly for the
lookup. The form keeps every field editable and preserves values you enter.
Metadata lookup failure does not remove the URL; enter a title manually to
continue. The saved post keeps the destination URL as entered and starts as an
unpublished draft.

The footer reports whether the draft was saved. Storage failures leave editing and downloads available.
See [shared editor features and runtime differences](/docs/guides/renderer-integration-benchmarks/#shared-editor-features-and-runtime-differences) for the comparison.

## Command palette and Markdown commands

Select **Commands** or press **Ctrl+Shift+P** (**Cmd+Shift+P** on macOS).
Search for special syntax, then press **Enter** to insert it at your selection.
See the [command palette guide](editor-command-palette.md) for the catalog and keyboard controls.

To open the command palette with a complete example loaded, visit `/?example`.


## Review changes before publishing (experimental)

Select the **Changes** tab in the workspace sidebar to view the Git inventory.
It enumerates the entire Git working tree, not only posts
opened in the editor. Editable posts are listed first, followed by other files.
Use **Find a file** or the **All**, **Posts**, and **Other** filters to narrow
the list. Select **Review queue** to open the review surface. It shows each
file's Git status and compares supported text files with the last committed
version. The review labels Git HEAD as **Last commit**; it does not mean that
version is live. Deleted files, conflicts, binary or oversized files, and files
outside the selected site show why a source diff is unavailable.

**Keep in workspace** leaves the saved file and Git status unchanged. It marks
the file as reviewed only for the current queue session. **Review and reset**
opens the full review for an unstaged Markdown post. **Open in editor** opens
an editable Markdown post without changing its worktree status. In the reset
review, select **Discard changes…** and confirm to restore the exact last
committed source.
Select **Full review…** in the diff header to choose Rendered diff, Side by
side, Source diff, or Full diff for a Markdown post.
The workbench disables reset for staged files, deleted files, conflicts,
non-Markdown files, and files outside the selected site. The server also
rejects stale source versions or a changed Git baseline. Other files can be
compared when they are supported text files, but this flow does not reset them.

The panel refreshes after a successful commit, push, or discard and
periodically while open. This is an inspection and triage view; multi-file
selection and **Publish selected** are not available. A local commit is not a
remote publication, and a push does not prove the site is live. If Git status
fails, the panel reports an error rather than pretending there are no changes.
The change count is loaded in the background when the editor starts and refreshed
when a site file saves or the editor regains focus. Opening the view shows its
last inventory immediately while a fresh Git status check runs. **Review** for
the active post is in the Changes panel header.

Use **Publish post** in the header or **Ctrl/Cmd+S** (Emacs **Ctrl-X Ctrl-S**,
Vim `:w`) to publish *only the active site Markdown post*. This explicitly
marks the post as public, saves it to disk, verifies the remote branch equals
the local pre-publish commit, then creates an isolated one-file commit and
pushes it. Any unrelated local commits block the quick path rather than
being pushed unexpectedly. The status says **Published to Git**, not Live;
deployment verification is separate. If the network push fails, the in-memory
session can retry the pending commit without creating a duplicate. If the
editor is reloaded in that state, inspect Git and remote status before
retrying; automatic recovery across browser restarts is not yet implemented.

Use the separate **Review changes** control or **Ctrl/Cmd+Shift+Enter** when
you want a detailed comparison of the active post. The optional advanced
review dialog still offers separate **Commit selected page** and
**Push to origin** controls; these are not required for the primary quick
publish path. For a Git-tracked file, **Discard changes…** requires a
confirmation and restores exactly the last committed source in the local workbench.
The server refuses stale source versions, changed Git HEAD, or staged edits to
that file. New **untracked, non-ignored** Markdown posts can also be discarded. After
explicit confirmation, the server checks Git status, HEAD, and the exact reviewed
source version before deleting just that local file. The browser clears only
that file's recovery draft and reloads workspace discovery. Staged, ignored,
or tracked files cannot be deleted through this control. Discard does not create a Git commit or deploy anything.
A different process changing the filesystem simultaneously is not covered by a
cross-process transaction; use a disposable test post while this is experimental.

### Local-first responsiveness (experimental)

The editor now focuses and accepts locally recoverable draft text without
waiting for the initial site-file list or index. The first Go preview is
deferred until the selected site file has finished attaching; a slow preview
is not a prerequisite for writing. Language assistance starts asynchronously
after the server's editor handler is ready, and the browser retries its
WebSocket when the service is still warming up or temporarily unavailable.
The LSP's connection state never determines whether a document is editable.

The startup diagnostics show **Editor interactive** separately from
**Browser file list**, **Metadata index**, and **First preview**. Use
**Copy diagnostics** alongside server startup logs to distinguish delays
from site discovery, language assistance, rendering, and browser updates.
The live preview status now reports Go render time separately from
**queue time**, which includes waiting for the render gate and preparing the
request before the renderer starts. Large queue times suggest contention;
large render times suggest a slow Markdown/feature renderer.
A large workspace's initial path discovery and appearance construction
still run before the HTTP listener becomes available; these are remaining
cold-start profiling targets, not claimed fixed by this pass.

Pierre Full diff uses compact hunk separators with virtual spacer buffers
disabled inside its nested review dialog. The supplied browser bundle is
updated along with its vendored build input. Run the browser regression
against a long document containing a tiny edit to verify the visual result.

### Autosave and external edits

The browser keeps a local recovery draft and autosaves to the loopback Go editor
workspace. Each successful `GET /workspace/file` returns a 64-character source
version (SHA-256 of the exact file bytes). The browser sends this as
`X-Markata-Expected-Version` with each `PUT /workspace/file`.

The Go editor rejects a missing version with HTTP 428 and a stale version with
HTTP 409 **without writing the file**, except when the current file already
contains the exact submitted bytes. In that case, the server acknowledges the
duplicate save and returns the current version without writing again. For a
real conflict, the browser keeps the recovery copy and pauses automatic writes
for that file, so it does not repeat the same stale request after every edit.
Select **Compare copies…** to review the browser draft and the current file on
disk. You can keep the browser draft after confirming the overwrite, or replace
it with the disk copy. Autosave resumes against the version you reviewed. On a
phone the two panels stack vertically. If the server is unavailable, the
browser draft can still be downloaded.
An older recovery record without a version token also
fails closed. An acknowledged save updates the browser's recovery version so a
normal reload does not invent a conflict. The current local Go handler
serializes its own workspace writes and Discard restores, but unrelated
filesystem writers can still race with the final write. This is not a
cross-process or hosted synchronization transaction.

The dialog first waits until the exact text in the editor has been saved to
the workspace. It then compares that saved source with the file from the last
Git commit. **Rendered diff** opens first. **Side by side** shows both rendered
versions and keeps matching sections in view while you scroll; mobile can switch
between versions by tapping or swiping. **Source diff** shows unified Markdown,
and **Full diff** opens the detailed source viewer. The **Commit message** tab
shows the message for the local commit and lets you edit it. Metadata chips
identify draft/published flag changes, slug changes, and conflicting publication
flags.

The review can create a local Git commit for the selected page, using the
message shown in the **Commit message** tab. After that commit, **Push to
origin** is a separate action. Pushing does not prove that a deployment
succeeded; when Builder Admin verification is configured, **Check live status**
can confirm the deployed source revision. “Last committed” is not necessarily
the live site revision. A scratch file must first be created or opened as a
workspace Markdown post. If the workspace has no Git history, the review
reports this rather than assuming a baseline.

The publish contract and planned deployment safety gates are defined in
[Editor publishing](../../spec/spec/EDITOR_PUBLISHING.md).


### Verify a remote publication with Builder Admin (opt-in)

Publishing has deliberately separate steps: review the exact Markdown and
frontmatter, commit only that selected page, then explicitly push to `origin`.
A successful push **does not** prove the site is live.

For a site deployed through Markata Builder Admin, the local editor can
optionally verify the current release's recorded source Git commit:

```sh
export MARKATA_PUBLISH_BUILDER_URL="https://builder.example.com"
export MARKATA_PUBLISH_BUILDER_TOKEN="<your authorized bearer token>"
markata editor --workspace ./site
```

Use an API token authorized for Builder Admin's `/api/v1/state`. The editor
stores neither the token nor builder URL in the browser; the local backend
makes a short, authenticated, server-side GET request. Remote builder URLs
must use HTTPS (loopback HTTP is permitted for development). Redirects are
not followed, so bearer tokens cannot be forwarded to a different host.

After pushing, select **Check live status**. This checks
`revision.live_source_commit`, with a fallback to the current release's
source revision when available. It reports **Live · Verified** only if that
revision matches the exact pushed SHA. Unknown revision, unavailable API,
unconfigured auth, failed build, and different live SHA never count as success.

This verification identifies the release's **source Git revision**, not
its byte-for-byte rendered output, CDN cache rollout, or availability at
every public edge. Site URL health checks and other deployment providers
need dedicated follow-up adapters.


### Credential-bearing Dropper uploader platform support

A hosted Markata editor can delegate photo uploads to a **separate local
Dropper uploader process**, so the editor and browser never receive the
upstream upload token. That capability currently relies on a Unix-domain
socket with POSIX owner/group permissions (mode `0660`). It is supported
only on operating systems where those permission bits are enforceable.

On Windows, the credential-bearing Unix-socket uploader **fails closed** with
an unsupported-platform error. Windows may report `0666` for Unix sockets
even after `chmod 0660`; a successful chmod call is not evidence of a
restricted Windows ACL. Use a secure Linux/macOS uploader backend or wait for
a Windows ACL-based transport to be implemented and tested. This restriction
does **not** disable the standard local editor or Markdown site builder.
