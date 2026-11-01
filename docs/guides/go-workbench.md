---
title: "Live Markdown editor with Go"
description: "Run or embed a Go server with a live Markata Markdown editor"
date: 2026-10-28
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

Use the appearance picker to select a palette, color mode, style, font, or text size.
Open **Editor settings** to choose line wrapping, keymap, and layout. Drag the
divider between the source and preview panes to resize them. Choose **Editor
only** or **Preview only** to focus on one pane; **Editor and preview** restores
the split.
Scrolling either pane keeps the corresponding Markdown visible in the other.
When you write or move the cursor, the preview follows the editor. Scroll sync
is on by default; turn off **Sync editor and preview scrolling** in **Editor
settings** to scroll each pane independently. The choice is saved in this browser.
Choose Default, Vim, Emacs, or Sublime keymaps in **Editor settings**.
In Vim mode, press `/` to search forward, `?` to search backward, then `n` or
`N` to move between matches.
Code highlighting, icons, admonitions, tables, checklists, Mermaid, and Chart.js use local assets.

Press **Ctrl+P** (**Cmd+P** on macOS) to find a site file by its path, title,
slug, or tag. The finder shows the current and recently opened files when its
search is empty.
Select a result or press **Enter** to open it in the editor.
Choose **Feeds** or **Tags** in the workspace selector to browse posts by feed
or tag. Feed names are shown as roots with folders relative to the workspace.

Edits to selected site files save back to the workspace after a short pause.
Unsaved scratch documents stay in browser storage. The browser also keeps a
recovery copy of workspace edits. Use **Ctrl+Shift+P** and search for **Download
current file** to export a copy, or **Open Markdown file** to import a file as
an unsaved draft.
A storage error appears beside the render metrics. Editing and downloads remain available.
The server writes only selected workspace files. Scratch drafts remain local to
the browser origin, so drafts on one port remain separate from drafts on another.

The preview updates after a 250 ms pause. One request runs at a time, and newer edits replace pending edits.
Old responses cannot replace a newer draft. Errors preserve the last successful preview.
The footer shows render time, round trip time, source size, and preview size.
Repeated source uses a one-entry memory cache. Appearance changes update the preview without a render request.
Workspace post embeds in `![[post-slug]]` or `![embed](path/to/post.md)` format
render as cards linked to the matching site post. Fenced code examples remain
literal source.

## Language assistance

The editor connects to the existing Markata LSP through `/lsp` over WebSocket.
The server indexes the directory from `--workspace`. The default directory is the current directory.
The browser LSP accepts same-origin connections on `localhost` or a loopback IP only.

Type `[[` for wiki-link suggestions, including aliases. Type `[label](` for Markdown document paths.
Press **Ctrl+Space** to request suggestions at the cursor. Use the arrow keys and Enter to select a suggestion.
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
The footer reports whether the draft was saved. Storage failures leave editing and downloads available.
See [shared editor features and runtime differences](/docs/guides/renderer-integration-benchmarks/#shared-editor-features-and-runtime-differences) for the comparison.

## Command palette and Markdown commands

Select **Commands** or press **Ctrl+Shift+P** (**Cmd+Shift+P** on macOS).
Search for special syntax, then press **Enter** to insert it at your selection.
See the [command palette guide](editor-command-palette.md) for the catalog and keyboard controls.

To open the command palette with a complete example loaded, visit `/?example`.
