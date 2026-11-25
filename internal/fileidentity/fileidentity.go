// Package fileidentity supplies identities for reusing verified local file data.
package fileidentity

import "time"

// Identity captures the inode and the metadata that changes when its bytes change.
type Identity struct {
	Device   uint64
	Inode    uint64
	Size     int64
	Mode     uint32
	Modified int64
	Changed  int64
}

// StableBefore checks that both file timestamps precede the observation by more
// than the resolution of a timestamp on a coarse filesystem.
func (id Identity) StableBefore(observed time.Time) bool {
	threshold := observed.Add(-time.Second).UnixNano()
	return id.Inode != 0 && id.Modified < threshold && id.Changed < threshold
}
