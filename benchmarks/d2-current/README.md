# Current D2 native Go 1.27 comparison

This isolated nested module pins D2 v0.9.0, which requires Go 1.27, and
duplicates the [Go 1.26 / D2 v0.7.1 benchmark](../d2-native/) fixtures.

It is **not** imported by Markata or compiled into the default CLI. It exists
to separate the effects of D2 engine version and process spawning.

```sh
cd benchmarks/d2-current
go mod tidy
go test -race -v -run TestNativeD2SVG ./...
go test -run '^$' -bench BenchmarkNativeD2 -benchmem -benchtime=200ms -count=3 ./...
```

The path-filtered GitHub Action compares:

1. D2 v0.7.1 in-process with Go 1.26.
2. D2 v0.9.0 in-process with Go 1.27.
3. D2 v0.9.0 subprocess CLI (including process startup).

Compare equal version and equal fixture before drawing conclusions about
native versus subprocess execution. This is still a research experiment:
Markata's production module has upgraded to Go 1.27 in [PR #1901](https://github.com/WaylonWalker/markata/pull/1901), so a separate worker is no longer required solely for the Go version. The production adapter and its trust/network/asset-policy integration remain unimplemented.

Tracked in [issue #1896](https://github.com/WaylonWalker/markata/issues/1896).

## First same-version results (2026-10-31)

[Go 1.27 benchmark workflow run 41246400810](https://github.com/WaylonWalker/markata/actions/runs/41246400810) passed the in-process D2 v0.9.0 SVG tests, race detector and benchmarks.

| Fixture | Native D2 v0.9.0 median | D2 v0.9.0 CLI median including process launch |
| --- | ---: | ---: |
| 3 nodes | 0.782 ms | 33.03 ms |
| 18 nodes | 2.765 ms | 36.89 ms |
| 72 nodes | 9.315 ms | 47.89 ms |

Native medians are three Go benchmark samples, reusing a text ruler; CLI
medians are five subprocess runs after one warmup. Native SVG files for these
fixtures were 6,133, 15,197 and 45,315 bytes uncompressed and contained no
embedded font data URIs in these test cases.

**Conclusion:** the earlier v0.7.1 native results were dominated by version
differences, not by the intrinsic overhead of in-process Go rendering.
D2 v0.9.0 native is substantially faster than the subprocess CLI on the
same fixtures. Markata has since upgraded to Go 1.27. Strongly prefer evaluating
native v0.9.0 for production, but first check binary size, layouts,
themes/fonts, assets, offline policy and actual builder integration.

## SVG contract integration check

The `TestD2CurrentSVGCanPublishAsMarkataAsset` test imports
`github.com/WaylonWalker/markata/pkg/visualizations` through the isolated
nested module's local `replace => ../..`. It runs real D2 v0.9.0 diagrams
through the shared SVG validator, immutable publisher, and figure-markup
helper. The workflow test selector must explicitly run this integration test.

This experiment cannot establish universal syntax/theme compatibility;
external links, images, HTML labels, and custom font resources still need
separate fixtures.

## Retest against upgraded Markata main

[Workflow 41255950545](https://github.com/WaylonWalker/markata/actions/runs/41255950545) passed on a branch containing the merged Go 1.27 Markata baseline. Its native D2 v0.9.0 benchmark yielded medians of **1.050 / 3.692 / 13.382 ms** for 3 / 18 / 72 nodes (three samples each). Its Markata SVG policy/publisher test explicitly executed and passed all five fixture groups: simple, medium, large, labels and shapes. The root Markata binary still does **not** use D2; this is research evidence for the follow-up production adapter.
