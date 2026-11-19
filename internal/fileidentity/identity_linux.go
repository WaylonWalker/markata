//go:build linux

package fileidentity

import (
	"io/fs"
	"syscall"
)

// Available reports whether the platform supplies verified file identity.
const Available = true

// Read returns the identity of a file on Linux.
func Read(info fs.FileInfo) Identity {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return Identity{}
	}
	return Identity{Device: stat.Dev, Inode: stat.Ino, Size: info.Size(), Mode: uint32(info.Mode()), Modified: info.ModTime().UnixNano(), Changed: stat.Ctim.Nano()}
}
