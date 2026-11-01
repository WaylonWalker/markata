# Markata publisher full source diff

The Go workbench runs without npm and does not request remote scripts. The
optional, Pierre-powered source review is built as a self-contained asset.

```sh
cd tools/publish-diffs
npm install
npm run build
```

Commit the generated `pkg/workbench/web/pierre-diffs.bundle.js` alongside
dependency version updates. The runtime loads this asset **only** after the
user selects **Full diff**. The lightweight source diff remains available
even when the optional asset hasn't been built.

Pinned: `@pierre/diffs` 1.6.2, `esbuild` 0.25.12. The bundle should be
reviewed for size, licenses and offline behavior before a production merge.
