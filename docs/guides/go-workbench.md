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

Run the editor from the Markata repository:

```bash
go run ./cmd/markata-workbench --listen 127.0.0.1:8767
```

Open <http://127.0.0.1:8767/>. Edit Markdown on the left and view the preview on the right.
Press Ctrl+C to stop the server. Without `--listen`, the server uses port 8766.

The Go server calls Markata's standalone renderer directly. It starts no Python process and needs no WASM runtime.
It uses the same Markdown renderer as the [Python API](./python-renderer/) and [FastAPI editor](./fastapi-editor/).
It does not build a site or load site configuration.

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

The JSON response contains `document`, `render_ms`, `input_bytes`, `document_bytes`, and `cached`.
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

This editor supports trusted local authoring. Authored HTML and scripts are enabled inside an opaque-origin iframe sandbox.
Preview scripts cannot read the editor's DOM or storage. Authored links and images can access external networks.
Keep the default loopback binding for local use. Add authentication and suitable resource controls before remote deployment.

All built-in assets come from the server. Offline reload is not supported.
For site builds and site-specific plugins, use the normal [CLI](../reference/cli/).
