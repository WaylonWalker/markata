# Markata documentation site

Documentation is authored in `docs/**/*.md`, rendered by `scripts/build_docs_site.py`, and served as static HTML. The landing page is `docs-site/home.html`. No JavaScript is needed to read it.

## Local preview

```bash
python -m pip install 'Markdown==3.8.2'
DOCS_BASE_PATH=/ python scripts/build_docs_site.py
python -m http.server --directory docs-site/_site 8000
```

Open <http://localhost:8000/>.

## Production: homelab with Argo CD

**The source of truth is GitOps, not GitHub Pages.** After a docs change is merged to `main`, `.github/workflows/docs-site.yml` verifies the HTML, builds the Nginx image, smoke-tests deep routes, and publishes versioned images to `ghcr.io/waylonwalker/markata-docs`.

The production deployment is defined in [homelab-argo](https://github.com/WaylonWalker/homelab-argo) under `argo-apps/apps/markata-docs.yaml` and `k8s/markata-docs/`. Argo CD manages the deployment, Service, and Traefik ingress. cert-manager issues the HTTPS certificate and ExternalDNS handles the hostname.

Initial hostname: **https://markata.waylonwalker.com**.

Images are versioned as `0.0.<workflow-run-number>` so the existing Renovate configuration in homelab-argo can propose image upgrades. The initial Argo manifest must reference an image tag that has already been published. A merge of Renovate's image update PR advances production; image publishing alone does not replace running pods. Configure auto-merge in Renovate deliberately if unattended promotion is desired.

If the GHCR package is private after its initial creation, make the `markata-docs` package public in GitHub package settings, or configure a pull secret in the deployment. The public documentation should not require an unauthenticated runtime to pull a private image.

The older `kubernetes/markata-go-docs.yaml` is historical and is **not** the active deployment. Don't apply it alongside the Argo app.

To validate local runtime without deploying:

```bash
python -m pip install 'Markdown==3.8.2'
DOCS_BASE_PATH=/ python scripts/build_docs_site.py
docker build -f docs-site/Dockerfile -t markata-docs:local .
docker run --rm -p 8080:8080 markata-docs:local
```

Open <http://localhost:8080/>.

## GitHub Pages

The old Pages deployment has been removed. No `actions/deploy-pages` jobs should run, and enabling GitHub Pages is not necessary.
