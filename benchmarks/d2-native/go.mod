module github.com/WaylonWalker/markata/benchmarks/d2-native

go 1.26.0

// D2 v0.8+ requires Go 1.27; test the last compatible module line
// without changing the Go requirement or dependencies of Markata itself.
require oss.terrastruct.com/d2 v0.7.1
