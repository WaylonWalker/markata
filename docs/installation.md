---
title: "Install Markata"
description: "Install the Go static site generator with curl, Bash, mise, GitHub Releases, or Go."
date: 2026-10-18
published: true
tags: [documentation, installation]
---

# Install Markata

The preferred executable is **`markata`**. The installer also installs the compatibility command **`markata-go`**.

Markata is a single Go binary; you do not need Python to run it. Choose your installation method.

## Curl (Linux / macOS)

Read [the installer source](https://github.com/WaylonWalker/markata/blob/main/install.sh) before running it, especially on managed machines.

```bash
curl -fsSL https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh | bash
markata version
```

This checks GitHub Releases, detects platform and architecture, and installs into a writable binary directory (usually `~/.local/bin` for non-root users). Add that directory to your `PATH` if necessary.

## Bash (download then inspect)

```bash
curl -fsSLo install-markata.sh https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh
less install-markata.sh
bash install-markata.sh
markata version
```

To choose an install location or a version, pass variables to **the Bash process**:

```bash
MARKATA_INSTALL_DIR="$HOME/.local/bin" MARKATA_VERSION="v0.13.0" bash install-markata.sh
```

The script refuses to overwrite an unrelated `markata` executable, including one belonging to Python Markata.

## mise

Install from GitHub Releases:

```bash
mise use -g github:WaylonWalker/markata
markata version
```

Or omit `-g` to pin it in your project's mise configuration. Check the installed command name: historical release archives may still contain `markata-go` rather than `markata`.

## GitHub download

Get the archive for your operating system and architecture from [GitHub Releases](https://github.com/WaylonWalker/markata/releases). Verify its checksum against the accompanying `checksums.txt`, extract it, and put the executable on your `PATH`.

Historical v0.13.0 release archives are named `markata-go_<version>_<os>_<arch>.tar.gz` (Windows uses `.zip`). New releases may use `markata_...`. Do not guess archive filenames: select one from the actual release assets.

## Go (source build)

Use **Go 1.26 or newer** for the current source module:

```bash
go install github.com/WaylonWalker/markata/cmd/markata@main
markata version
```

Until a release tag contains the new entry point, use `@main` for this canonical command rather than `@latest`. Make sure `$(go env GOPATH)/bin` is in `PATH`.

## PyPI / pip — the original Python project

**This does not install Markata Go.** The `markata` name on PyPI still points to the older Python implementation at [WaylonWalker/markata-py](https://github.com/WaylonWalker/markata-py).

```bash
python -m pip install markata
```

Use that only if you intentionally want **Python Markata**. Go releases and the PyPI migration are separate decisions.

## Verify and start

```bash
markata version
markata --help
mkdir my-site && cd my-site
markata config init
markata new "Hello World"
markata serve
```

Continue with the [quickstart](/docs/quickstart/). If you already use `markata-go`, that command remains supported. Keep your `markata-go.toml` filename and `[markata-go]` configuration section; the Go config file has **not** been renamed. See [Naming migration](/docs/guides/name-migration/).

For more options, check [troubleshooting](/docs/troubleshooting/) and [the release list](https://github.com/WaylonWalker/markata/releases).
