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
Use the keymap selector for Default, Vim, Emacs, or Sublime editing.
Code highlighting, icons, admonitions, tables, checklists, Mermaid, and Chart.js use local assets.

The editor stores drafts and preferences in browser storage for the current origin.
Use **Download .md** to save your draft to a file.
A storage error appears beside the render metrics. Editing and downloads remain available.
The server does not store drafts. Drafts on port 8765 remain separate from drafts on other ports.

The preview updates after a 250 ms pause. One request runs at a time, and newer edits replace pending edits.
Old responses cannot replace a newer draft. Errors preserve the last successful preview.
The footer shows render time, round trip time, source size, and preview size.
Repeated source uses a one-entry memory cache. Appearance changes update the preview without a render request.

## Language assistance

The editor connects to the existing Markata LSP through `/lsp` over WebSocket.
The server indexes the directory from `--workspace`. The default directory is the current directory.
The browser LSP accepts same-origin connections on `localhost` or a loopback IP only.

Type `[[` for wiki-link suggestions, including aliases. Type `[label](` for Markdown document paths.
Press **Ctrl+Space** to request suggestions at the cursor. Use the arrow keys and Enter to select a suggestion.
Existing frontmatter, admonition, and mention completion also use the LSP.
Broken wiki links receive an underline. The underline tooltip shows the diagnostic message.

The editor sends unsaved text as a synthetic `.markata/browser-draft.md` document.
This document stays in memory. Markdown path suggestions start at the workspace root because the draft has no source file.
The Files tree opens source files from the selected workspace using Markata glob rules. The Tags view groups those files by loaded frontmatter. Edits stay in browser drafts; download a copy to keep changes outside the browser. See [workspace navigation](./renderer-site-tree/).
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

Use **Open .md** to import a Markdown file. An oversized file leaves the current draft intact.
Use **Download .md** to keep a copy outside browser storage.
The footer reports whether the draft was saved. Storage failures leave editing and downloads available.
See [shared editor features and runtime differences](/docs/guides/renderer-integration-benchmarks/#shared-editor-features-and-runtime-differences) for the comparison.

## Insert Markdown commands

Select **Commands** or press **Ctrl+Shift+P** (**Cmd+Shift+P** on macOS).
Search for special syntax, then press **Enter** to insert it at your selection.
See the [command palette guide](editor-command-palette.md) for the catalog and keyboard controls.

To open the command palette with a complete example loaded, visit `/?example`.
