# Naming and identity transition (Go 1.0)

The Go implementation becomes **Markata**, with `markata` as its primary command.
The previous Go identity `markata-go` remains a compatibility alias. The
historical Python project will be renamed **markata-py** at a coordinated cutover.

## Stage 1: Additive CLI and release assets

- Both `cmd/markata` and `cmd/markata-go` build against the same command package.
- `markata` prints the canonical name in help and version output; the old
  `markata-go` binary retains its previous name.
- Release artifacts include `markata_<version>_<os>_<arch>` and
  `markata-go_<version>_<os>_<arch>`, keeping published older archives valid.
- The official installer tries the canonical asset first, falls back to the
  legacy archive, and installs both executable names.
- `MARKATA_INSTALL_DIR` / `MARKATA_VERSION` take precedence over legacy
  `MARKATA_GO_INSTALL_DIR` / `MARKATA_GO_VERSION`. Existing old variables
  remain supported.
- Existing `markata-go.toml`, `[markata-go]`, `MARKATA_GO_*`, module path,
  Docker tags, Helm settings and application state are **unchanged** in stage 1.

## Stage 2: Coordinated project cutover

The GitHub repository `WaylonWalker/markata` is currently occupied by the
original Python implementation. Rename it to `markata-py` before moving the
Go repository to `markata`; preserve tags, releases, history and canonical
package pointers. Treat Go module path / imports, CLI file discovery, config
namespace, environment variable precedence, Docker images, Helm charts,
plugin skill bundles, GitHub actions, site infrastructure, and package registry
ownership as one audited change set, not a search/replace.

The canonical config name `markata.toml` MUST NOT silently parse old Python
config with Go semantics. A later cutover must implement an explicit format
discriminator or compatibility warning and test precedence where both names
and both environment prefixes are present.

PyPI `markata` cannot simply be overwritten with a Go distribution: existing
Python users may depend on `markata<1`. Prepare a versioned transition and
document a `markata-py` installation strategy before any registry release.

## Release gate

Do not publish `1.0.0` until:
- old and new CLIs produce the same builds on representative sites;
- both installer paths and version selection are tested;
- Go modules and imported consumers build, across the published platforms;
- Docker / Helm production deployments and CI references have been audited;
- migration of Python naming and PyPI ownership is verified;
- RC release and migration guide are checked end-to-end.
