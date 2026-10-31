# D2 native Go 1.26 compatibility spike

**Research experiment only.** No new Markata user-facing feature, runtime
dependency, Markdown syntax, or site-builder integration is introduced.

The production Markata module currently targets Go 1.26. D2's official
[v0.8.2](https://github.com/d2lang/d2/blob/v0.8.2/go.mod) and
[v0.9.0](https://github.com/d2lang/d2/blob/v0.9.0/go.mod) modules declare
Go 1.27. The older [v0.7.1](https://github.com/d2lang/d2/blob/v0.7.1/go.mod)
uses the `oss.terrastruct.com/d2` import path and Go 1.24.

This nested Go module pins D2 **v0.7.1** specifically to test whether
native, in-process D2 SVG output is a viable renderer without changing
Markata's Go toolchain or increasing its default binary size.

## Run

Requires Go 1.26 and network access once for the external module.

```sh
cd benchmarks/d2-native
go mod tidy
go test -race -v -run TestNativeD2SVG ./...
go test -run '^$' -bench BenchmarkNativeD2 -benchmem -benchtime=200ms -count=3 ./...
```

A path-filtered workflow runs these commands on pull requests changing this
directory. The workflow also downloads the official **D2 v0.9.0** Linux CLI
archive, verifies it against its release SHA256SUMS, and runs the same graphs
through the CLI in five sequential subprocess invocations. That provides a
measured backend process-startup baseline. The warm native benchmarks
compile+render small, medium and large linear graphs with one shared D2 text
ruler. They print median-ready Go benchmark
samples, raw SVG bytes, gzip size and allocation statistics. The smoke test
logs a coarse embedded-font signal.

These results are **not** published-page performance measurements. This
first comparison captures native warm in-process rendering against CLI
per-process execution, but it does not yet measure full site-build latency,
peak memory, custom-font/theme compatibility or matching visual output. A
production decision must also compare binary size, fonts/theme compatibility,
PNG/SVG fidelity and whether the resulting SVG meets
`pkg/visualizations.NewSVGAsset` restrictions.

## Explicit constraints

- The benchmark only tests Dagre, the built-in Grape Soda theme and no file
  imports. ELK, TALA, custom fonts, palettes, complex graphs, icons, links,
  nested diagrams, and reproducible offline builds remain unverified.
- A pinned older library version may miss syntax added in newer releases.
  Do not silently imply parity with the current D2 CLI.
- No claim that D2 is safe for untrusted authored diagrams is made here.
- Once evidence exists, choose the production adapter in [issue #1896](https://github.com/WaylonWalker/markata/issues/1896)
  and implement on top of [foundation PR #1892](https://github.com/WaylonWalker/markata/pull/1892).
