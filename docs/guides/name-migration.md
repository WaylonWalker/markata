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
legacy release where needed, installing both command names. If a different application (including Python Markata) already owns the `markata` command at the destination, the installer stops rather than overwriting it. Move the conflicting executable or choose a separate install directory.

```bash
curl -sSL https://raw.githubusercontent.com/WaylonWalker/markata/main/install.sh | bash
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

`MARKATA_GO_OUTPUT_DIR`, `MARKATA_GO_SITE_DIR` and other existing
`MARKATA_GO_*` overrides continue working. Configuration overrides applied by
the Go config loader also accept the canonical `MARKATA_*` aliases, for example
`MARKATA_OUTPUT_DIR`. The canonical name takes precedence if both are set;
conflicting values record `name_migration_warnings` on the loaded config
(without exposing the values). Other runtime-specific settings may still use
legacy names; migrate those only after their individual support is verified.

Renaming the config file to `markata.toml` or the TOML section to `[markata]`
is **not yet supported** by the Go CLI, and may be confused with an original
Python Markata configuration. Keep existing TOML filenames/sections for now.

For the installer, `MARKATA_INSTALL_DIR` and `MARKATA_VERSION` are preferred;
the legacy `MARKATA_GO_INSTALL_DIR` and `MARKATA_GO_VERSION` still work.

## GitHub repository cutover completed

The repository swap was completed on October 18, 2026. The Go implementation
is now [WaylonWalker/markata](https://github.com/WaylonWalker/markata)
(repository ID `1140034778`); the original Python implementation is
[WaylonWalker/markata-py](https://github.com/WaylonWalker/markata-py)
(repository ID `333449359`). The old Go GitHub URL may redirect to the
canonical Go repository, but use the new URL explicitly in new automation.
Old Python bookmarks targeting `WaylonWalker/markata` now land on Go.

Read the [repository cutover runbook](repository-cutover.md) for the owner
rename checks and post-cutover infrastructure audit. The GitHub rename alone
**does not** migrate the Go module import path, Python PyPI ownership,
legacy configuration names, container registry permissions, or live sites.

## Package and production migration

The Go module is still declared as `github.com/WaylonWalker/markata-go`
until the separate module/import migration passes tests. During that
transition, pinned old Go imports may continue referencing older releases;
avoid mixing both module paths in a single module's dependency graph.
For the first Go 1.0 release, new imports should use
`github.com/WaylonWalker/markata`.

Existing `markata-go.toml`, `[markata-go]`, `MARKATA_GO_*`,
`markata-go` binaries and legacy GHCR builder references are intentionally
preserved. Do not rename Kubernetes PVCs, Helm release names or secrets merely
to match branding. Audit Argo CD chart repo URLs and permissions separately,
and validate a release candidate before touching production images.
