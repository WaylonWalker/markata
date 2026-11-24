//go:build js && wasm

package themes

import "embed"

// WASM hosts link presentation assets from the generated static bundle.
// Keep a small filesystem root for the shared theme API.
//
//go:embed default/static/css/variables.css
var defaultTheme embed.FS
