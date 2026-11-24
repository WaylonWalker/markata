# Browser WASM rendering

A js/wasm entrypoint exposes `markataRender(source)` as a JSON-string API using
RenderStandalone schema version 1. Raw HTML is disabled. Invalid inputs and
renderer panics return error envelopes without terminating the runtime.
The renderer accepts at most 200,000 UTF-8 bytes per document.

The static workbench initializes one persistent Go runtime in a dedicated worker.
Edits are debounced; only one render runs at a time and pending edits coalesce.
Revision checks prevent stale previews. Assets are transferred once per worker
and reused between renders. Drafts and editor preferences stay in localStorage.
Storage failures remain visible without preventing editing or export.

The preview uses a script-enabled, opaque-origin sandbox. No Markdown is sent
to a server. Authored external images/links can still access the network.
A versioned service worker caches the complete generated app for offline reload,
only reports readiness after installation, and preserves unrelated caches.
The build emits static files including the matching Go wasm_exec.js runtime;
no CDN or rendering server is required. WASM supports client Mermaid only.

The persistent native StandaloneRenderer has the same output as independent
RenderStandalone calls. Instances support sequential use only. The optional
OmitAssetContent flag links bundled dependencies without reading embedded
asset bytes. Generated scripts retain their contents. WASM embeds only font
metadata and a minimal theme root; native builds retain complete embedded assets.

The build script exports the native appearance catalog and all linked assets.
Font URLs remain relative to the app. The parent embeds preview font data and reads script/style assets through its
controlled origin, because opaque-origin sandbox requests bypass service workers.
The server supports gzip negotiation, loopback binding, and explicit port/output
options. It does not serve a render API or receive document contents.

## Hosted editor preview pilot

The hosted Go editor may offer an opt-in worker preview after its first
authoritative Go preview. The worker must not delay that first preview. It may
run only when the server and worker advertise the same WASM asset fingerprint,
the workspace has no active site configuration, the worker is ready, and a
Go render of the selected Markdown has certified matching HTML and source
lines. Unsupported Markdown, a missing fingerprint, a worker error, or a
fingerprint mismatch keeps the Go renderer active.

The browser must not download or initialize the WASM worker when the browser
does not report device memory greater than 4 GiB or more than four logical
processors, requests reduced data use, reports a 2G or slow-2G connection, or
lacks Worker or WebAssembly support. A WASM response must carry the fingerprint
of the exact uncompressed module; the worker must compare it with the fingerprint
in its URL before starting Go. Hosted worker, runtime, and module responses stay
authenticated, private, and no-store. The first preview remains Go-rendered,
and WASM remains opt-in until hosted compatibility, memory, and edit-to-visible
measurements justify a wider rollout.
When the worker reports readiness or a render result, it MUST include the
current WebAssembly linear-memory size. Local diagnostics MUST retain the
largest value observed in the session without recording source text, paths, or
credentials.

## Reproducible worker benchmarks

The worker benchmark uses the seven renderer integration fixtures without changing their bytes.
Each fixture uses a fresh worker, two excluded warmups, and fifteen measured renders.
The report records fixture hashes, browser version, host, raw timings, median, and p95.
Worker timings exclude startup, message transfer, layout, and paint. They are not HTTP latency measurements.
