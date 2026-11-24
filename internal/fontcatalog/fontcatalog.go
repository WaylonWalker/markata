//go:build !js

// Package fontcatalog owns the built-in font catalog and its immutable assets.
// Keeping the embed here prevents runtime code from depending on the process
// working directory or on a source checkout being present.
package fontcatalog

import (
	"embed"
	"io/fs"
)

// Files contains the built-in catalog, lockfile, manifests, licenses, and
// WOFF2 files shipped in the markata-go executable.
//
//go:embed markata-fontpacks.yaml markata-fonts.lock.yaml */*
var Files embed.FS

// FS returns the embedded built-in font asset filesystem.
func FS() fs.FS { return Files }
