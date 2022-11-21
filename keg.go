package keg

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	_fs "github.com/rwxrob/fs"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/keg/kegml"
	"github.com/rwxrob/keg/model"
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

// ScanDex takes the target path to a keg root directory returns a
// model.Dex object.
func ScanDex(kegdir string) (*model.Dex, error) {
	var dex model.Dex
	dirs, _, _ := NodePaths(kegdir)
	sort.Slice(dirs, func(i, j int) bool {
		_, iinfo := _fs.LatestChange(dirs[i].Path)
		_, jinfo := _fs.LatestChange(dirs[j].Path)
		return iinfo.ModTime().After(jinfo.ModTime())
	})
	for _, d := range dirs {
		_, i := _fs.LatestChange(d.Path)
		title, _ := kegml.ReadTitle(d.Path)
		// TODO change this to DexEntry.MD
		id, err := strconv.Atoi(d.Info.Name())
		if err != nil {
			continue
		}
		entry := model.DexEntry{U: i.ModTime(), T: title, N: id}
		dex = append(dex, entry)
	}
	return &dex, nil
}

// MakeDex calls ScanDex and writes (or overwrites) the output to the
// reserved dex node file within the kegdir passed. File-level
// locking is attempted using the go-internal/lockedfile (used by Go
// itself).
func MakeDex(kegdir string) error {
	dex, err := ScanDex(kegdir)
	if err != nil {
		return err
	}
	jsonpath := filepath.Join(kegdir, `dex`, `nodes.json`)
	if err := file.Overwrite(jsonpath, dex.String()); err != nil {
		return err
	}
	mdpath := filepath.Join(kegdir, `dex`, `nodes.md`)
	if err := file.Overwrite(mdpath, dex.MD()); err != nil {
		return err
	}
	// TODO UpdateUpdated
	return nil
}

// MkTempNode creates a text node directory containing a README.md
// file within a directory created with os.MkdirTemp and returns a full
// path to the README.md file itself. Directory names
// are always prefixed with keg-node.
func MkTempNode() (string, error) {
	dir, err := os.MkdirTemp("", `keg-node-*`)
	if err != nil {
		return "", err
	}
	readme := path.Join(dir, `README.md`)
	err = file.Touch(readme)
	if err != nil {
		return "", err
	}
	return readme, nil
}

// ImportNode moves the nodedir into the KEG directory for the kegid giving
// it the nodeid name. Import will fail if the given nodeid already
// existing the the target KEG.
func ImportNode(from, to, nodeid string) error {
	to = path.Join(to, nodeid)
	if _fs.Exists(to) {
		return _fs.ErrorExists{to}
	}
	return os.Rename(from, to)
}

// UpdateUpdated sets the updated YAML field in the keg info file.
func UpdateUpdated(kegpath string) error {
	kegfile := filepath.Join(kegpath, `keg`)
	updated := UpdatedString(kegpath)
	return file.ReplaceAllString(
		kegfile, `(^|\n)updated:.+(\n|$)`, `${1}updated: `+updated+`${2}`,
	)
}

// Updated parses the most recent change time in the dex/node.md file
// (the first line) and returns the time stamp it contains as
// a time.Time. If a time stamp could not be determined returns time.
func Updated(kegpath string) (*time.Time, error) {
	kegfile := filepath.Join(kegpath, `dex`, `nodes.md`)
	str, err := file.FindString(kegfile, model.IsoDateExpStr)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(model.IsoDateFmt, str)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// UpdatedString returns Updated time as a string or an empty string if
// there is a error.
func UpdatedString(kegpath string) string {
	u, err := Updated(kegpath)
	if err != nil {
		log.Println(err)
		return ""
	}
	return (*u).Format(model.IsoDateFmt)
}
