---
title: "Markata naming migration"
description: "Moving from markata-go to Markata while preserving existing Go and Python installations"
date: 2026-10-16
published: true
tags:
  - documentation
  - migration
---

# Markata naming migration

The Go static site generator previously known as **markata-go** is becoming
**Markata** for the 1.0 release line. The existing Python implementation will
be called **markata-py** after the repositories and packages are coordinated.

## What works now

New development builds provide the preferred `markata` executable. The legacy
`markata-go` command remains supported, including existing tagged downloads.
The official install script attempts a canonical archive and falls back to a
legacy release where needed, installing both command names.

```bash
curl -sSL https://raw.githubusercontent.com/WaylonWalker/markata-go/main/install.sh | bash
markata --help
markata build

# Existing automation still works:
markata-go build
```

If you install from Go source rather than binary releases, the new command is
available from `main` until the first tagged version with the new entry point:

```bash
go install github.com/WaylonWalker/markata-go/cmd/markata@main
```

## Configuration stays as-is for now

There is no immediate need to rename configuration:

```toml
# markata-go.toml remains correct during the compatibility window
[markata-go]
output_dir = "output"
```

`MARKATA_GO_OUTPUT_DIR`, `MARKATA_GO_SITE_DIR` and all other existing
`MARKATA_GO_*` overrides continue working. Renaming the file to
`markata.toml` or the TOML section to `[markata]` is **not yet supported**
by the Go CLI, and may be confused with a Python Markata configuration.

For the installer, `MARKATA_INSTALL_DIR` and `MARKATA_VERSION` are preferred;
the legacy `MARKATA_GO_INSTALL_DIR` and `MARKATA_GO_VERSION` still work.

## Future 1.0 cutover

The Go GitHub repository is still
[WaylonWalker/markata-go](https://github.com/WaylonWalker/markata-go),
because [WaylonWalker/markata](https://github.com/WaylonWalker/markata) is
the historical Python repository. Their names will swap only after Python
moves to `markata-py`, with explicit link and dependency migration.

Do not change Go import paths, container references, Helm releases or Python
dependencies just because the CLI is now named `markata`. Versioned upgrade
guidance for those surfaces will follow the repository and package cutover.

For now: use `markata` interactively, retain old config and deployments, and
keep `markata-go` anywhere exact compatibility is essential.
