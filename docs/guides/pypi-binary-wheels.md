---
title: "PyPI installation of Go Markata using Simon Willison's go-to-wheel"
description: "Decision (October 18, 2026): pip install markata should install the Go"
date: 2026-10-31
published: true
tags:
  - documentation
---


**Decision (October 18, 2026): `pip install markata` should install the Go
implementation.** This is an intentional handoff of the existing PyPI
distribution name, not a rename of the Python import-compatible library.
Markata remains on the **0.x** release line. Existing Python releases remain
immutable on PyPI and in [WaylonWalker/markata-py](https://github.com/WaylonWalker/markata-py).

## Packaging tool

We use Simon Willison's [go-to-wheel](https://github.com/simonw/go-to-wheel)
(version `0.2`) to compile our `./cmd/markata` Go CLI into eight native wheels:

- Linux amd64/arm64, both glibc and musl (four wheels).
- macOS Intel and Apple Silicon (two wheels).
- Windows amd64 and arm64 (two wheels).

`go-to-wheel` requires the `go.mod` directory to also be the CLI package.
Markata's CLI is under `cmd/markata`. Our
[`scripts/build-pypi-wheels.sh`](../../scripts/build-pypi-wheels.sh)
temporarily copies its existing `main.go` into the module root solely for
compilation, then removes it, preserving the canonical CLI implementation and
normal Go module. It sets version and commit metadata in the built binary.

When installed, `markata` is a minimal Python entry point which executes the
bundled Go binary, **not** the historical Python Markata CLI or library. There
is no network download at first run and no Go compiler required on end-user
machines.

## Build and test, without publishing

Run from the Go repository root with Go 1.27 and Python 3.10+:

```bash
python -m pip install 'go-to-wheel==0.2'
bash scripts/build-pypi-wheels.sh 0.14.0.dev0
MARKATA_WHEEL_VERSION=0.14.0.dev0 python scripts/check-pypi-wheels.py dist/pypi-wheels

python -m venv .venv-markata
. .venv-markata/bin/activate
python -m pip install ./dist/pypi-wheels/markata-0.14.0.dev0-py3-none-manylinux_2_17_x86_64.whl
markata --help
markata version
```

The [wheel CI workflow](../../.github/workflows/pypi-go-wheels.yml)
automatically repeats the build and verifies tags, embedded OS executables,
metadata, installation in a clean venv, and `python -m markata`.

## Automated PyPI publication

The Python-to-Go handoff of the existing `markata` PyPI project was
completed with Go v0.14.0. The one-time migration confirmation is **retired**.
The historical Python releases remain immutable and available by version pin.

For new stable 0.x releases:

1. A maintainer/agent publishes a tested `v0.x.y` tag using the
   [Release workflow](../../.github/workflows/release.yml).
2. Once that workflow **succeeds**, the
   [wheel workflow](../../.github/workflows/pypi-go-wheels.yml) starts via
   `workflow_run`. This also works for GitHub releases generated using
   `GITHUB_TOKEN`, which do not reliably trigger other release events.
3. The wheel workflow validates the completed run, published GitHub release,
   immutable tag commit and existing PyPI state. It then builds eight wheels,
   runs clean-install and platform smoke checks, and requires editor parity
   for newly tagged versions.
4. PyPA Trusted Publishing uses short-lived OIDC credentials in the protected
   `pypi` environment. A final job verifies all eight wheel filenames
   appear on PyPI. **There is no second manual confirmation.**

To recover a missed upload, choose **Actions → Go Markata PyPI wheels →
Run workflow**, select branch `main`, and supply the *existing GitHub
release* version, e.g. `0.15.0`. The workflow refuses mismatched/unpublished
tags and does not republish a complete version already present on PyPI.
A partial existing PyPI release is treated as an error rather than overwritten.
Pull requests only perform the build and tests, never the upload.

The original `v0.15.0` tag predates the newer editor browser-parity gate
and `pkg/renderprotocol` package. Its historical release verification
is preserved for a safe backfill without changing or moving that tag.
All later releases must satisfy their full modern checks.

The registered publisher remains `WaylonWalker/markata`,
`pypi-go-wheels.yml`, environment `pypi`. Keep the historical
`markata-py` repository's publishing disabled; it must not race the Go
package on the `markata` project name.

## Compatibility consequences

- Users needing the historical Python library should pin the last Python
  release in a separate environment: `python -m pip install markata==0.10.0`.
  This is **not** the Go executable and is not a maintained 0.x Go release.
- `import markata` from a new Go wheel is only a tiny binary-launching
  wrapper; it is not API compatible with old Python Markata plugins.
- Existing Go binaries, `markata-go` entry points, Kubernetes Helm settings,
  and Go GitHub releases are unchanged by this PyPI handoff.
- On unsupported Python/platform/architecture combinations, pip may select an
  *older Python release* rather than the Go wheel. The eight-wheel matrix must
  be checked before making blanket cross-platform promises.
- No stable Go **1.0** release or repository version tag is required.

The initial name handoff was separately authorized; subsequent versions are published automatically from verified GitHub releases.

## Current installation status

The PyPI handoff has occurred: current `markata` releases install the Go CLI after a version break from the Python implementation. Historical Python wheels remain immutable and can be pinned. This handoff does not imply Python API compatibility. Version pins preserve historical installers; the automated process above documents the current release path.
