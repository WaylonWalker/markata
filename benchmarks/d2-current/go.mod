module github.com/WaylonWalker/markata/benchmarks/d2-current

go 1.27.0

// Experimental sidecar-only module; never part of Markata's production go.mod.
require (
	github.com/WaylonWalker/markata v0.0.0
	github.com/d2lang/d2 v0.9.0
)

// Test against the actual Markata SVG contract in this checkout.
// This replace is local to the nested spike module, not the production graph.
replace github.com/WaylonWalker/markata => ../..
