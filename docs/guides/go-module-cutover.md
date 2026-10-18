# Go module import path cutover

The canonical Go source module for **Markata 1.0** is
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
- **Pinned pre-1.0 consumers:** existing historical tags for
  `github.com/WaylonWalker/markata-go` have not been deleted.
  Avoid depending on both the historical and new module paths in the
  same dependency graph; they describe the same source with different
  module identities and may result in duplicate packages.
- **Python users:** the Python package import `markata` and the
  existing `markata` distribution on PyPI are unchanged. A
  Go CLI install must not silently overwrite an existing Python
  `markata` executable.

## Canonical source install

```bash
go install github.com/WaylonWalker/markata/cmd/markata@main

# Legacy Go executable remains available from the canonical module:
go install github.com/WaylonWalker/markata/cmd/markata-go@main
```

For reproducible installations, pin a verified release tag rather than
`@main`. Do not tag 1.0 until canonical Go imports, builder images,
release archives, installer URLs, Argo CD source references and representative
site builds are validated.

## Change control

The module path migration changed only tracked Go files, `go.mod`,
and GoReleaser's linker symbol package paths. The deterministic rewrite
script lives at `scripts/migrate-go-module-path.py`; the dedicated
PR workflow fails if any old module import remains.

These are **not** reasons to rename Kubernetes PVCs, releases, secrets,
Helm chart names, historical Docker tags, old Go binaries, or Python
PyPI distributions.
