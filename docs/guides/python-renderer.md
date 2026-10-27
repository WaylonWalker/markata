---
title: "Render Markdown in Python and FastAPI"
description: "Use Markata's renderer with feature signals and bundled CSS and JavaScript"
date: 2026-10-27
published: true
tags:
  - documentation
  - python
  - markdown
---

Install the official Go-based `markata` wheel containing this API (available with
the release that includes this feature):

```bash
uv add markata fastapi uvicorn
```

Render a complete response with conditional CSS and JavaScript inlined:

```python
from fastapi import FastAPI
from fastapi.responses import HTMLResponse
from markata import render_async

app = FastAPI()

@app.get("/", response_class=HTMLResponse)
async def index():
    result = await render_async("# Hello\n\n!!! note\n    Rendered by Markata.", inline_assets=True)
    return HTMLResponse(result.document(title="My app"))
```

`render()` offers the same options for synchronous applications. Both use the
existing Go Markdown renderer shipped inside the Python wheel. Each call launches
one process, so cache repeated content in your application. `render_async()` uses
async subprocess I/O, terminates the process on timeout or cancellation, and does
not block the ASGI event loop. `timeout` defaults to 30 seconds. Failures raise
`markata.RenderError`.

## Fragments and signals

```python
from markata import render

result = render("```go\nfmt.Println(42)\n```", inline_assets=True)
print(result.html)              # Article markup
print(result.features)          # code, admonitions, roads, mermaid, chartjs, ...
print(result.head_html)         # Ordered style tags
print(result.body_end_html)     # Ordered scripts
```

In an existing layout, insert `head_html` into the head, `html` inside an
`<article class="post-content">`, and `body_end_html` after the article. These
strings contain generated HTML; use the appropriate trusted-markup mechanism in
your template engine. `document(title="...", lang="en")` creates that layout for
you and escapes the title and language.

The renderer returns core theme CSS plus only the optional assets the content
requires. Mermaid and Chart.js use Markata's vendored browser libraries; inline
mode includes them too, without CDN requests. Inlining diagrams can make a
response considerably larger. `external_urls` identifies any external browser
dependencies; for the currently supported bundled features it is empty. Authored
images and links still refer to their authored URLs.

Core theme CSS includes global styles. For an application with its own design,
select the assets you want from `result.assets` and supply compatible CSS variables
and `.post-content` styling. This API does not configure a site's theme palette.

## Serve assets separately

The default `inline_assets=False` emits links under `/markata-assets`, or your
`asset_base`. Each asset is a dictionary containing `id`, `kind`, `placement`,
and `content` (or `url` for external assets). Asset IDs are suitable for caching
and deduplication; scripts must keep their manifest order.

A small app can precompute its content and serve the bundled assets directly:

```python
from fastapi import FastAPI, HTTPException
from fastapi.responses import HTMLResponse, Response
from markata import render

app = FastAPI()
page = render("# Example\n\n```python\nprint('hello')\n```", asset_base="/md-assets")
assets = {asset["id"]: asset for asset in page.assets if "content" in asset}

@app.get("/md-assets/{asset_id:path}")
def asset(asset_id: str):
    item = assets.get(asset_id)
    if item is None:
        raise HTTPException(status_code=404)
    media_type = "text/css" if item["kind"] == "style" else "text/javascript"
    return Response(item["content"], media_type=media_type)

@app.get("/", response_class=HTMLResponse)
def index():
    return page.document(title="Example")
```

For multiple pages, combine their bundled assets by ID and retain each page's
ordered generated initialization scripts. Generated chart IDs are local to one
render result: render all Markdown for a document together. Do not deduplicate
generated scripts across unrelated results. If browser content is replaced after
page load, the application must execute its initialization scripts; simply
assigning `innerHTML` does not execute them.

## Scope and trusted content

Supported: Markata's Markdown extensions, syntax highlighting, containers,
admonitions, Mermaid, Chart.js, roads, sortable tables, and opt-in checklists. Rendering does not run a site build;
site-dependent link resolution, embeds, image transforms, and other lifecycle
plugins are outside this API.

Raw authored HTML is disabled by default. Set `unsafe=True` only for trusted
content. Rich Markdown output is **not an HTML sanitizer**; sanitize untrusted
content using your application's policy. Inline styles and scripts also require
a compatible Content Security Policy. Applications with strict CSP should serve
assets themselves and apply their own nonce/hash policy to initialization scripts.

## CLI protocol

The Python API uses the same machine-readable interface available to other languages:

```bash
printf '# Hello\n' | markata render --json --inline-assets
```

Without `--json`, `markata render` emits just the HTML fragment. `--unsafe` enables
raw HTML; `--asset-base /static/markdown` changes linked bundled asset URLs. This
command reads stdin without discovering a site configuration or writing files.
