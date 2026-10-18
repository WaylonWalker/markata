# Shipping the Go Markata CLI through Python wheels

The Go implementation can be distributed through PyPI using platform-specific
Python wheel files. **This is a build/packaging integration, not a PyPI handoff.**
It does not change our **0.x** release policy.

## What tool we use

[deadnews/go-wheel-action](https://github.com/deadnews/go-wheel-action) packages
the Go command at `./cmd/markata` directly. It uses `CGO_ENABLED=0` and
cross-compiles native executables into eight wheels: glibc/musl Linux (amd64,
arm64), macOS (amd64, arm64), and Windows (amd64, arm64).

It creates a small Python entry point named `markata` that launches the bundled
Go binary. Installing the wheel does **not** require users to install Go, clone
this repository, or download binaries on the first command invocation.
The upstream action is pinned to an exact Git commit in our workflow.

CI: [Go CLI PyPI Wheel Smoke](../../.github/workflows/pypi-go-wheels.yml).

## What is live now

- `pip install markata` continues to install the **original Python project**
  maintained in [WaylonWalker/markata-py](https://github.com/WaylonWalker/markata-py).
- The original Python project still owns the PyPI distribution name `markata`,
  Python import `markata` and its historical releases.
- The Go repository produces no-publish **test wheels** with the eventual
  canonical package/entry-point identity `markata`. Using that identity in CI
  **does not** upload or overwrite anything on PyPI.
- The wheel test uses a synthetic, never-published version
  `0.0.0.dev0`, verifies all eight archive formats, and installs/runs the
  Linux wheel in a clean virtual environment.

To try the generated wheel after downloading the CI artifact:

```bash
python -m venv .venv-markata-go
. .venv-markata-go/bin/activate
python -m pip install ./markata-0.0.0.dev0-py3-none-manylinux_2_17_x86_64.whl
markata --help
markata version
```

Install and smoke-test the wheel **inside an isolated environment**: the Go
command and historical Python command both use the executable name `markata`.

## Deciding whether to publish

This project has intentionally **no PyPI upload job**, no PyPI credentials,
no publishing `id-token: write` permission, and no tag-triggered wheel release.
GoReleaser GitHub binary releases are still separate.

Before publishing a Go wheel, explicitly choose a distribution identity:

1. **Safe coexistence:** publish under a different PyPI distribution name
   (for example `markata-go-bin`) while leaving the Python `markata`
   distribution untouched. This requires adjusting the wheel project and
   console entry-point identity or using a packager that supports separate
   names; do not assume renaming the wheel file is valid.
2. **Handoff of `markata` on PyPI:** an intentional major compatibility
   decision. Requires communicating that `pip install markata` changes
   implementations, protecting existing Python users and imports, and an
   explicitly approved release plan. A GitHub repository rename is **not**
   authorization to do this.

Publishing should use a separate, manually approved workflow with a pinned
version and PyPI trusted-publishing setup. Never silently publish a Go binary
under the historical Python distribution just because a tag was pushed.

The existing [compatibility checklist](https://github.com/WaylonWalker/markata/issues/1682)
and [name migration roadmap](https://github.com/WaylonWalker/markata/issues/1628)
continue to track the unresolved operational checks.
