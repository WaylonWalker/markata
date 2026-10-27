# Go module import path cutover

The canonical Go source module for **Markata** (still on the 0.x line) is
`github.com/WaylonWalker/markata`, following the October 18, 2026
GitHub repository swap. It is independent of the original Python module,
which now lives at [WaylonWalker/markata-py](https://github.com/WaylonWalker/markata-py).

## Existing users

- **Using the CLI or builder containers:** you can keep the
  `markata-go` binary alias, `markata-go.toml`, `[markata-go]`,
  `MARKATA_GO_*` and legacy GHCR names for now. See the
  [naming migration guide](name-migration.md).
- **Importing Go packages:** replace imports starting with
  `github.com/WaylonWalker/markata-go/` with
  `github.com/WaylonWalker/markata/`, and run `go mod tidy`.
  The Go module's own declaration has changed as well.
- **Pinned historical Go consumers:** existing historical tags for
  `github.com/WaylonWalker/markata-go` have not been deleted.
  Avoid depending on both the historical and new module paths in the
  same dependency graph; they describe the same source with different
  module identities and may result in duplicate packages.
- **Python users:** historical pinned Python releases remain available on PyPI, but
  current PyPI releases install the Go CLI after a version break. A
  Go CLI install must not silently overwrite an existing Python
  `markata` executable.

## Canonical source install

```bash
go install github.com/WaylonWalker/markata/cmd/markata@main

# Legacy Go executable remains available from the canonical module:
go install github.com/WaylonWalker/markata/cmd/markata-go@main
```

For reproducible installations, pin an already published, verified release
tag rather than `@main` when that tag supports your command and module path.
The name migration does **not** authorize any new release, prerelease or tag.
Validate Go imports, installer compatibility, builder images, Argo CD sources
and representative site builds independently while remaining on 0.x.

## Change control

The module path migration changed only tracked Go files, `go.mod`,
and GoReleaser's linker symbol package paths. The deterministic rewrite
script lives at `scripts/migrate-go-module-path.py`; the dedicated
PR workflow fails if any old module import remains.

These are **not** reasons to rename Kubernetes PVCs, releases, secrets,
Helm chart names, historical Docker tags, old Go binaries, or Python
PyPI distributions.
