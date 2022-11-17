package keg

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rogpeppe/go-internal/lockedfile"
	_fs "github.com/rwxrob/fs"
	"github.com/rwxrob/keg/node/read"
)

// NodePaths returns a list of node directory paths contained in the
// keg root directory path passed. Paths returns are fully qualified and
// cleaned. Only directories with valid integer node IDs will be
// recognized. Empty slice is returned if kegroot doesn't point to
// directory containing node directories with integer names.
//
// The lowest and highest integer names are returned as well to help
// determine what to name a new directory.
//
// File and directories that do not have an integer name will be
// ignored.
var NodePaths = _fs.IntDirs

// MakeDex takes the target path to a keg root directory and creates or
// replaces the dex keg index file there. Each line of a dex file
// contains three things:
//
//     1. Second last changed in UTC in ISO8601 (RFC3339) (just digits)
//     2. Unique integer identifier
//     3. Current title (always first line of README.md)
//
// Note that the second of last change is based on *any* file within the
// node directory changing, not just the README.md or meta files.

func Dex(kegdir string) (string, error) {
	var dex string
	dirs, _, _ := NodePaths(kegdir)
	sort.Slice(dirs, func(i, j int) bool {
		_, iinfo := _fs.LatestChange(dirs[i].Path)
		_, jinfo := _fs.LatestChange(dirs[j].Path)
		return iinfo.ModTime().After(jinfo.ModTime())
	})
	for _, d := range dirs {
		_, i := _fs.LatestChange(d.Path)
		title, _ := read.Title(d.Path)
		dex += fmt.Sprintf("%s %s %s\n",
			_fs.IsosecModTime(i), d.Info.Name(), title,
		)
	}
	return dex, nil
}

// MakeDex calls Dex and writes (or overwrites) the output to the 'dex'
// file within the kegdir passed. File-level locking is attempted using
// the go-internal/lockedfile (used by Go itself).
func MakeDex(kegdir string) error {
	dex, err := Dex(kegdir)
	if err != nil {
		return err
	}
	return lockedfile.Write(
		filepath.Join(kegdir, `dex`),
		strings.NewReader(dex),
		fs.FileMode(0644),
	)
}
