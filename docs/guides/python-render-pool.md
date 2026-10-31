---
title: "Warm renderer pool for FastAPI"
description: "Keep Markata renderer workers warm for concurrent FastAPI requests"
date: 2026-10-31
published: true
tags: [documentation, python, fastapi, performance]
---

The Python standalone render API normally starts the bundled Go executable for
each call. Applications that render repeatedly can instead own a fixed pool of
persistent renderer processes using `RendererPool`. The pool is started and
stopped by the application's lifespan; **no separate daemon or CLI setup is
required**.

## FastAPI lifespan integration

```python
from contextlib import asynccontextmanager

from fastapi import FastAPI, Request
from pydantic import BaseModel
from markata import RendererPool

class Preview(BaseModel):
    markdown: str

@asynccontextmanager
async def lifespan(app: FastAPI):
    async with RendererPool(workers=2, max_pending=64) as pool:
        app.state.renderer_pool = pool
        yield

app = FastAPI(lifespan=lifespan)

@app.post("/preview")
async def preview(body: Preview, request: Request):
    result = await request.app.state.renderer_pool.render(body.markdown)
    return {"html": result.html, "features": result.features,
            "assets": result.assets}
```

The pool starts all workers at startup, verifies the persistent protocol via
a capability request, and maintains a queue of idle renderers. Each worker is
one bundled Go process that handles sequential JSON-lines requests. Requests
can run concurrently across workers. The Go protocol caches configured parser
instances so subsequent requests avoid parser reinitialization too.

Set `inline_assets=True` to embed styles and scripts in the returned result
or serve `result.assets` separately, as with `markata.render_async`. Rendering
options `unsafe`, `asset_base`, and `inline_assets` can be supplied on individual
pool requests.

## Lifecycle and limits

- `workers=2` starts two subprocesses per running application process. With
  `uvicorn --workers 4`, that is **eight** Go subprocesses. Size the pool
  accordingly; it does not share processes across Uvicorn workers.
- `max_pending=64` limits the number of waiting requests in addition to the
  currently active worker slots. When saturated, callers receive `RenderError`
  rather than growing an unbounded backlog.
- `render(..., timeout=30)` covers queue wait and rendering. A timed-out or
  cancelled in-flight request kills and replaces that worker so a late result
  cannot contaminate a later request.
- Unexpected child exit or broken protocol also replaces the affected process.
  A failed replacement makes the pool report an error rather than hang forever.
- Exiting the lifespan terminates and reaps its child processes. Do not share
  the pool across Python event loops or fork a process after starting it.
- Requests retain the standalone renderer's 200,000-byte source limit. The
  wire protocol has a 1 MiB request frame limit.

The pool requires a bundled Go executable supporting `markata render --stdio`.
The original `render()` and `render_async()` remain backward-compatible one-shot
APIs and do **not** implicitly start a pool.

The initial stdio protocol does not load site configuration, so `workspace` is
not a pool option yet. This is not an HTML sanitizer; keep `unsafe=False` for
untrusted authored HTML and apply your application's own safety policy.

## Performance verification

Measure median and p95 of warm pool calls versus `render_async()` on the same
machine, with at least 100 requests for each fixture. Also measure throughput
at concurrency 1, 2, and 4, resident memory per pool worker, initial startup
latency, and large Chart.js response transfer. The existing standalone Go
benchmarks are not end-to-end pool timings. Do not assume the same latency
until the Python-to-Go pipe and JSON costs are measured.
