# Markata publisher full source diff

The Go workbench does not request remote scripts. The optional, Pierre-powered
source review is compiled into local ES modules; it loads only after the user
selects **Full diff**, and the highlighter and theme chunks load on demand.

```sh
cd tools/publish-diffs
npm install
npm run build
```

Commit `pkg/workbench/web/pierre-diffs.bundle.js` and the generated
`pkg/workbench/web/pierre-diffs/` chunks alongside dependency version updates.
The lightweight source diff remains available when the optional assets have not
been built.

Pinned: `@pierre/diffs` 1.6.2, `esbuild` 0.25.12. The build uses Pierre's
Highlights backend and does not include Shiki's full language registry.
