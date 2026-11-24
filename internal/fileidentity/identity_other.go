//go:build !linux

package fileidentity

import "io/fs"

// Available reports whether the platform supplies verified file identity.
const Available = false

// Read returns a zero identity on platforms without platform-specific support.
func Read(_ fs.FileInfo) Identity { return Identity{} }
