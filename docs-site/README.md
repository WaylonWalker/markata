# Markata documentation website and homelab builder

This site's authoring source is `docs/**/*.md`, its custom homepage is `docs-site/home.html`, and the renderer is `scripts/build_docs_site.py`. It uses Python Markdown to preserve the current documentation layout. It is not yet a native Markata CLI site: do not point `markata builder-admin` at the repository root and assume it will publish these pages.

## Local preview

```bash
python -m pip install 'Markdown==3.8.2'
DOCS_BASE_PATH=/ python scripts/build_docs_site.py
python -m http.server --directory docs-site/_site 8000
```

Browse <http://localhost:8000>.

## Homelab hosting

The Argo deployment lives in [WaylonWalker/homelab-argo](https://github.com/WaylonWalker/homelab-argo/tree/main/k8s/markata-docs), with `markata-docs` Argo Application. This follows the builder-and-site split from waylonwalker.com:

- **Public**: https://markata.waylonwalker.com — Nginx serves `/data/site/current` from a persistent Longhorn volume.
- **Protected builder**: https://build.markata.waylonwalker.com — HLab-authenticated admin for Build now, status, releases, and rollback.
- **Background updates**: `docs-site/builder_admin.py` polls the canonical GitHub `main` every 45 seconds and only rebuilds after a source SHA change, or manual trigger. It publishes releases atomically and retains previous versions.
- **Images**: `.github/workflows/docs-site.yml` checks both the static docs site and the live builder on PRs. On `main`, it publishes `ghcr.io/waylonwalker/markata-docs-builder` and the immutable `ghcr.io/waylonwalker/markata-docs` fallback/seed image.
- **No Python runtime on public Nginx**: Python is limited to the builder pod.

The builder image version is pinned explicitly in homelab-argo for controlled upgrades. The docs contents themselves refresh from `main` without a fresh container build or Argo sync. The initial fallback image seeds a new site PVC before the builder's first checkout; it is not the active publishing mechanism afterward.

The builder HTTP handler assumes it is reachable **only through Traefik**, which strips any client-provided identity headers before HLab ForwardAuth. The application also requires an identity header from a configured trusted cluster CIDR and same-origin POSTs. Never expose the Service directly as a public unauthenticated endpoint.

## Debugging

```bash
kubectl -n markata-docs get deploy,pods,pvc,ingress
kubectl -n markata-docs logs deploy/markata-docs-builder -c builder --tail=100
kubectl -n markata-docs rollout status deploy/markata-docs
curl -fsSI https://markata.waylonwalker.com/docs/quickstart/
```

The source clone is ephemeral; the current release, retained releases, and build history are persisted on Longhorn. If the builder is unavailable, the last published site stays online.

Do not apply `kubernetes/markata-go-docs.yaml` alongside the Argo resources. GitHub Pages has been retired for this documentation website.
