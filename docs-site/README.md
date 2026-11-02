# Markata documentation site

The public documentation is built from `docs/**/*.md` by Markata itself. The
site-specific configuration is `markata-docs.toml`; it keeps the existing
`/docs/.../` routes, uses Markata's documentation layout and code highlighter,
and sets `docs/index.md` as the homepage.

## Build locally

```bash
go run ./cmd/markata-go --config markata-docs.toml build --clean
go run ./cmd/markata-go --config markata-docs.toml serve
```

The generated site is in `docs-site/_site/`. No separate Python Markdown
renderer or syntax-highlighting package is used.

## Homelab hosting

The Argo CD application combines the Markata Helm chart with the homelab
ingresses. The chart runs the same Markata Go builder used locally, keeps
published releases on the existing `markata-docs-site` volume, and serves the
current release with its public site deployment.

Builder Admin initializes a persistent Git checkout of
`WaylonWalker/markata`. A signed GitHub Push webhook at
`https://build.markata.waylonwalker.com/webhook` fast-forwards `main` and queues
a Markata build. The authenticated Builder Admin interface remains at
`https://build.markata.waylonwalker.com/`.

The webhook HMAC secret is stored in the Kubernetes Secret referenced by
`builderAdmin.webhook` in homelab values. It is not stored in this repository.
Create that secret before rolling out the chart, then configure a GitHub
repository webhook for Push events with the matching secret.

See `homelab-argo/k8s/markata-docs/README.md` for rollout and verification
steps.
