# Markata documentation site

The public documentation is built from `docs/**/*.md` by Markata itself. The
site-specific configuration is `markata-docs.toml`; it keeps the existing
`/docs/.../` routes and uses Markata's bundled default templates, theme, and
code highlighter. `docs/index.md` is the homepage and uses the bundled landing
layout. The `/docs/` path prefix keeps the docs layout and generated sidebar.
The docs site does not load the repository's main-site template overrides.

## Build locally

```bash
go run ./cmd/markata-go --config markata-docs.toml build --clean
go run ./cmd/markata-go --config markata-docs.toml serve
```

The generated site is in `docs-site/_site/`. No separate Python Markdown
renderer or syntax-highlighting package is used.

## Browser and visual checks

CI installs the pinned Playwright, pixelmatch, and pngjs versions, builds the
site, then runs `node --test scripts/docs-site.browser.test.cjs`. The suite
checks 320, 390, 768, 1024, and 1440 pixel viewports in light and dark themes.
It compares homepage and Quickstart screenshots with the reviewed PNGs in
`scripts/visual-baselines/docs-site/`; failures upload screenshots and pixel
diffs as workflow artifacts. The comparison ignores anti-aliasing, allows a
small height delta, uses a per-pixel color threshold, and caps changed pixels
at 2% of the full page. Larger changes fail and upload screenshots and pixel
diffs for review.

To run the browser check locally, install the same pinned packages and
Chromium, build the site, then run the test command above. Review screenshot
diffs before changing a baseline. Set `UPDATE_DOCS_VISUAL_BASELINES=1` only
after review to write replacement baselines.

## Homelab hosting

The Argo CD application combines the Markata Helm chart with the homelab
ingresses. The chart runs the same Markata Go builder used locally, keeps
published releases on the existing `markata-docs-site` volume, and serves the
current release with its public site deployment.

Builder Admin initializes a persistent Git checkout of
`WaylonWalker/markata`. A signed GitHub Push webhook at
`https://markata.waylonwalker.com/webhook` fast-forwards `main` and queues
a Markata build. The authenticated Builder Admin interface remains at
`https://build.markata.waylonwalker.com/`.

The webhook HMAC secret is stored in the Kubernetes Secret referenced by
`builderAdmin.webhook` in homelab values. It is not stored in this repository.
Create that secret before rolling out the chart, then configure a GitHub
repository webhook for Push events with the matching secret.

See `homelab-argo/k8s/markata-docs/README.md` for rollout and verification
steps.
