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
