# Markata documentation website

This is a small static documentation website, using Markdown from `docs/` as its source of truth. The landing page is `docs-site/home.html`; shared CSS and a tiny progressive-enhancement installer tab script are in `docs-site/`.

## Preview locally

```bash
python -m pip install "Markdown==3.8.2"
python scripts/build_docs_site.py
python -m http.server --directory docs-site/_site 8000
```

The default build path is `/markata/` for the GitHub Pages project site. For local preview at the site root:

```bash
DOCS_BASE_PATH=/ python scripts/build_docs_site.py
python -m http.server --directory docs-site/_site 8000
```

Then open <http://localhost:8000/>.

## Publishing

`.github/workflows/docs-site.yml` builds and checks on pull requests. On pushes to `main`, it uploads and deploys to GitHub Pages, **once the repo's Pages source is configured to GitHub Actions** under Settings → Pages.

The expected default URL is `https://waylonwalker.github.io/markata/` when Pages is enabled. Do not claim the URL is live until a successful Pages deploy. Custom domains require separate Pages, DNS, and `DOCS_BASE_PATH` configuration. This workflow does not change live homelab ingress or Argo CD.

No client-side JavaScript is required to read docs or install commands: the installer tabs are enhanced only after the page loads, and all methods remain visible without JS.
