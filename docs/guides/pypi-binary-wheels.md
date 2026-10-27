# PyPI installation of Go Markata using Simon Willison's go-to-wheel

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

Run from the Go repository root with Go 1.26 and Python 3.10+:

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

## Publishing the name handoff

As of this documentation change, **the public PyPI project is still Python
Markata until a Go wheel upload succeeds and is confirmed on PyPI**.
The existing Python release `markata==0.10.0` remains available as a pin.
The newest prior Python prerelease is `0.11.0.dev11`. PyPI does not permit
replacing those immutable files with Go binaries.

To publish Go Markata (first proposed PyPI Go version: **`0.14.0`**), the
owner must do the following once:

1. Sign into the existing [PyPI `markata` project](https://pypi.org/manage/project/markata/),
   open **Publishing**, and add a **GitHub Actions Trusted Publisher**:
   - owner: `WaylonWalker`
   - repository: `markata`
   - workflow filename: `pypi-go-wheels.yml`
   - environment: `pypi`
2. In GitHub `WaylonWalker/markata`, create or verify a **pypi** environment
   with protection (reviewers and restrictions to `main`) so an unreviewed
   workflow cannot publish packages.
3. Once CI passes and the trusted publisher is configured, open **Actions →
   Go Markata PyPI wheels → Run workflow**, on branch `main`, set:
   - `version`: `0.14.0` (or another approved, *unpublished* 0.x version)
   - `publish`: true
   - `confirm`: `REPLACE_PYTHON_MARKATA_WITH_GO`
   - Approve the `pypi` environment gate.
4. Check the resulting PyPI release and **test a new** `python -m venv` with
   `python -m pip install markata` and `markata version`.

The publish job is a separate, minimal-permission job using PyPA
`pypa/gh-action-pypi-publish` and short-lived OIDC tokens. Pull requests,
manual dry-runs, and non-main branches never publish. An explicit confirmation
string and GitHub environment gate prevent accidental uploads.

**Freeze the Python repo's old manual PyPI release workflow** as part of the
handoff so a future Python release does not retake the latest PyPI version.
It currently needs owner attention because GitHub rejects branch creation in
`WaylonWalker/markata-py`. Do not reuse its existing `pypi_password` secret.

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

This release is an explicitly authorized new PyPI packaging path, distinct from
the earlier GitHub repository rename. The version is chosen on dispatch.

## Current installation status

The PyPI handoff has occurred: current `markata` releases install the Go CLI after a version break from the Python implementation. Historical Python wheels remain immutable and can be pinned. This handoff does not imply Python API compatibility. The release preparation steps above document the transition process, not the present installer identity.
