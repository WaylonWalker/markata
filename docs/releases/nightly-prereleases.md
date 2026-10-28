# Nightly GitHub and PyPI prereleases

The scheduled `.github/workflows/pypi-go-wheels.yml` release train runs at
02:00 America/Chicago (GitHub's timezone-aware schedule handles DST).
It uses only the default branch. Scheduled Actions can be delayed or skipped by
GitHub; monitoring should alert when the expected nightly is missing.

## Identity and collision safety

From the highest stable `v0.x.y` git tag, increment the patch for the next
nightly. Example from `v0.14.0` on 2026-10-22:

- GitHub tag: `v0.14.1-dev.20261022`
- GitHub archive version: `0.14.1-dev.20261022`
- PyPI wheel version: `0.14.1.dev20261022` (PEP 440)

A date's tag/PyPI version may publish at most once. Both are checked before
building. If one side publishes and the other fails, **do not overwrite or
reuse the version**; investigate, reconcile manually or move to the next
nightly date. Stable release tags and existing PyPI files are never replaced.

## Release gates

- Full Go race suite and GoReleaser validation
- Cross-platform Go archive builds, SBOM generation, checksum signing
- Eight PyPI wheels, metadata checks, Linux clean-venv smoke,
  and macOS/Windows install and CLI smoke
- Immutable GitHub prerelease first; then wheel publishing with PyPI OIDC
- PyPI nightly job receives only `id-token: write`; other jobs do not
  receive PyPI credentials
- GitHub release publish job receives only `contents: write`

The separate manual stable wheel publish path remains gated by the existing
`pypi` environment and `REPLACE_PYTHON_MARKATA_WITH_GO` confirmation.
No `v1.*` tags are produced.

## One-time setup required from repository and PyPI owners

1. Add a GitHub environment named `pypi-nightly`, restrict deployment to
   `main`, and **do not require reviewers** if unattended nightly publication
   is desired. Protect branch and workflow changes via required reviews.
2. In PyPI project `markata` > Publishing > Add trusted publisher, use:
   owner `WaylonWalker`, repo `markata`, workflow
   `pypi-go-wheels.yml`, environment `pypi-nightly`.
   Keep the existing manual `pypi` publisher unchanged.
3. Merge this reviewed PR into `main` only when publication is authorized.
   No branch-local scheduled run publishes; schedule events are default-branch-only.
4. Verify the first `schedule` run and public GitHub + PyPI assets before
   treating unattended publishing as operational.

## Local no-publish test

```bash
python3 -m unittest discover -s scripts -p 'test_nightly_release_version.py' -v
python3 scripts/nightly-release-version.py --date 20261022
```

The existing `scheduled-release.yml` workflow remains a separate non-publishing
dry run while its release train changes are under review.
