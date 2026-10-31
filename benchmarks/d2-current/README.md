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
packaging a separate Go 1.27-built native sidecar would require an explicit
release/build design and trust/network/asset-policy integration.

Tracked in [issue #1896](https://github.com/WaylonWalker/markata/issues/1896).
