# D2 native Go 1.26 compatibility spike

**Research experiment only.** No new Markata user-facing feature, runtime
dependency, Markdown syntax, or site-builder integration is introduced.

At the time this historical comparison began, Markata targeted Go 1.26. It has since upgraded to **Go 1.27** in [PR #1901](https://github.com/WaylonWalker/markata/pull/1901). D2's official
[v0.8.2](https://github.com/d2lang/d2/blob/v0.8.2/go.mod) and
[v0.9.0](https://github.com/d2lang/d2/blob/v0.9.0/go.mod) modules declare
Go 1.27. The older [v0.7.1](https://github.com/d2lang/d2/blob/v0.7.1/go.mod)
uses the `oss.terrastruct.com/d2` import path and Go 1.24.

This nested Go module pins D2 **v0.7.1** specifically to test whether
native, in-process D2 SVG output is viable under the former Go 1.26 baseline without changing the root module dependency graph.

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

## Initial measurements (GitHub Actions, 2026-10-31)

[Workflow run 41246096778](https://github.com/WaylonWalker/markata/actions/runs/41246096778)
passed both the native Go 1.26 race/smoke tests and the verified CLI comparison.

| Diagram | Native Go v0.7.1 median, warm in-process | D2 CLI v0.9.0 median, fresh process | Native gzip bytes | CLI gzip bytes |
| --- | ---: | ---: | ---: | ---: |
| 3 nodes | 55.77 ms | 33.03 ms | 3,938 | 3,211 |
| 18 nodes | 135.97 ms | 36.89 ms | 5,294 | 4,521 |
| 72 nodes | 461.42 ms | 47.89 ms | 7,766 | 6,818 |

Native medians are from three Go benchmark iterations (with warm, reused
text ruler), while CLI medians are from five sequential subprocess executions
after one warmup. Measurements include different **D2 versions**, so these
numbers do **not** establish that subprocesses are intrinsically faster
than native-Go invocation. Rather, D2 v0.9.0 may have substantial layout,
text/renderer, or allocation optimizations that v0.7.1 lacks.

On these fixtures, the native v0.7.1 benchmark also allocated approximately
82 MiB (3 nodes), 106 MiB (18 nodes), and 198 MiB (72 nodes) *per render*.
Native output sizes were 11,237 / 21,629 / 56,921 uncompressed bytes, while
CLI v0.9.0 output 6,129 / 15,239 / 45,311 bytes.

**Follow-up completed:** [the D2 v0.9.0 native benchmark](../d2-current/) now runs in a separate Go 1.27 module and has passed against the upgraded production Go 1.27 baseline; see [run 41255950545](https://github.com/WaylonWalker/markata/actions/runs/41255950545). No D2 dependency has been added to Markata's root module.

## Explicit constraints

- The benchmark only tests Dagre, the built-in Grape Soda theme and no file
  imports. ELK, TALA, custom fonts, palettes, complex graphs, icons, links,
  nested diagrams, and reproducible offline builds remain unverified.
- A pinned older library version may miss syntax added in newer releases.
  Do not silently imply parity with the current D2 CLI.
- No claim that D2 is safe for untrusted authored diagrams is made here.
- The native D2 v0.9.0 option is favored by measured performance. Implement the production adapter separately under [issue #1896](https://github.com/WaylonWalker/markata/issues/1896), using the SVG foundation merged in [PR #1892](https://github.com/WaylonWalker/markata/pull/1892).
