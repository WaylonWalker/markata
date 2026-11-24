---
title: "Markata GitHub repository cutover — owner runbook"
description: "| GitHub repository | Before | After | Stable repository ID |"
date: 2026-10-18
published: true
tags:
  - documentation
---


> **Historical runbook:** The GitHub repository swap was completed on
> October 18, 2026. Do **not** rerun the rename command. Go is at
> `WaylonWalker/markata` and Python at `WaylonWalker/markata-py`; the Go
> project remains on **0.x** without a release tag planned by this migration.

> **No release tag, PyPI overwrite, or production deployment** is implied
> by a repository rename. Source-path, package and deployment migrations
> are separate gates.

## Identity and sequencing

| GitHub repository | Before | After | Stable repository ID |
| --- | --- | --- | --- |
| Historical Python | `WaylonWalker/markata` | `WaylonWalker/markata-py` | `333449359` |
| New Go implementation | `WaylonWalker/markata-go` | `WaylonWalker/markata` | `1140034778` |

**Order is mandatory.** The Python repository must vacate `markata` before the
Go repository can claim it. When Go takes over `WaylonWalker/markata`, existing
Python links to that name will start opening the **Go repository**, not Python.
Do not rely on GitHub redirects to keep Python bookmarks pointing to Python.

The GitHub connector currently does not expose repository rename permissions.
An owner can use repository Settings → General → Repository name, or the guarded
script below with an authenticated `gh` CLI that has repository administration
permissions.

## Before moving either repository

- Confirm [Python PR #175](https://github.com/WaylonWalker/markata-py/pull/175)
  (identity notice) and [#176](https://github.com/WaylonWalker/markata-py/pull/176)
  (disable automatic PyPI publishing, stabilize Python source URLs) are merged.
- Confirm Go [#1629](https://github.com/WaylonWalker/markata/pull/1629)
  (dual CLI), [#1671](https://github.com/WaylonWalker/markata/pull/1671)
  (canonical/legacy config environment aliases),
  [#1672](https://github.com/WaylonWalker/markata/pull/1672)
  (legacy/canonical GHCR images), and
  [#1673](https://github.com/WaylonWalker/markata/pull/1673)
  (GoReleaser release target) are merged.
- Verify GitHub Actions, branch rulesets, collaborators, Cloudflare Pages
  integrations, GitHub Apps and registry publication rights for both repos.
  Confirm the current Go production builder image is still retrievable.
- **Do not** publish another Python package to PyPI as an incidental result of
  a code/metadata edit. PyPI `markata` is still the historic Python package.
  `markata-py` is not automatically a new PyPI distribution.
- Do not change site configs called `markata-go.toml`, old Go module paths,
  sealed secrets, existing PVC names, Kubernetes Service identities, or Helm
  release names simply to match branding.

Run from a clone of the **Go** repository:

```bash
gh auth status
bash scripts/markata-repo-cutover.sh preflight
```

The script checks immutable GitHub repository IDs, merged preparation PRs and
that `markata-py` is unoccupied. It refuses ambiguous API/auth failures.

## Perform the owner-level rename

Use GitHub Settings manually in the exact order above, **or** run:

```bash
bash scripts/markata-repo-cutover.sh apply --execute
```

This requires an interactive explicit confirmation and checks repository IDs
both before and after each rename. **There is no automatic rollback.** If the
first rename succeeds and the second fails, stop and diagnose the second
failure; Python will already be at `markata-py`.

Once done, re-clone the Go repository or update its checkout:

```bash
git remote set-url origin https://github.com/WaylonWalker/markata.git
git remote -v
bash scripts/markata-repo-cutover.sh postcheck
```

For Python checkouts, point `origin` to
`https://github.com/WaylonWalker/markata-py.git`.

## Post-swap compatibility work — while remaining on 0.x

1. **Go source/module:** migrate `go.mod`, import paths, generated code,
   build/linker flags, GoReleaser import targets, CI, embedded paths, and
   code generators to `github.com/WaylonWalker/markata` together. Test
   `go test ./...`, representative external Go consumers, cross-platform
   installs and `go mod tidy`. Never silently rewrite historical Python
   `markata.toml` configuration or convert its semantics to Go.
2. **Installer/docs:** update canonical install scripts, release links,
   README badges, bundled agent skill and documentation to the *new* Go
   repository. Keep both `markata` and `markata-go` executables and
   old release archives. A raw.githubusercontent.com URL pointing at an old
   repository path must not be assumed to redirect after the move.
3. **Images/Argo:** verify the Docker job pushes identical tagged content to
   `ghcr.io/waylonwalker/markata`,
   `ghcr.io/waylonwalker/markata-go`,
   `ghcr.io/waylonwalker/markata-builder`, and
   `ghcr.io/waylonwalker/markata-go-builder` after the repo rename.
   Confirm GHCR Actions permissions to create/push the new package names.
   Pin images by digest when practical; do not delete old image tags.
4. **Home infrastructure:** audit `WaylonWalker/homelab-argo` Argo
   sources and their repo access credentials, especially:

   - `argo-apps/apps/waylonwalker-com-prod.yaml` (currently pinned to a commit)
   - `argo-apps/apps/go-waylonwalker-com.yaml`
   - `argo-apps/apps/wyattbubbylee-com-prod.yaml`
   - `argo-apps/apps/rhiannonwalker-com-prod.yaml`

   Their Git repository URLs were updated by [homelab-argo#111](https://github.com/WaylonWalker/homelab-argo/pull/111).
   Verify actual Argo CD repository access and application health rather than
   treating a PR merge as successful production reconciliation. Also audit
   `k8s/recipes-waylonwalker-com/` and `k8s/rhiannonwalker/source-seed-job.yaml`
   for historical builder images. Separate cosmetic renaming from functional
   Kubernetes resource changes.
5. **Python archival:** verify Python source/issue/docs URLs and Cloudflare
   pages deployments point to `markata-py`. Preserve old Python tags,
   import namespace and current installations. Plan PyPI `markata-py`
   distribution/alias separately rather than overwriting `markata` with
   an untested Go wrapper.
6. **Ongoing 0.x validation:** check old/new CLI parity on representative
   sites, installer compatibility, archive formats, Go consumers,
   Docker/Helm production smoke, registry permissions and GitHub source paths.
   Snapshot packaging can be tested without publishing a release. Any
   versioned tag requires a separate explicit decision; **no 1.0 release or
   release candidate is planned as part of the rename**.

Tracker: https://github.com/WaylonWalker/markata/issues/1628
