//go:build js && wasm

package fontcatalog

import (
	"embed"
	"io/fs"
)

// Browser rendering links font files from the exported static bundle. The
// catalog remains available without embedding a second copy of every WOFF2.
//
//go:embed markata-fontpacks.yaml markata-fonts.lock.yaml */manifest.yaml
var Files embed.FS

// FS returns font metadata; browser hosts serve font bytes separately.
func FS() fs.FS { return Files }
