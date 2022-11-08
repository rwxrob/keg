package node

import (
	"os"
	"path"

	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/file"
)

// MkTemp creates a text node directory containing a README.md
// file within a directory created with os.MkdirTemp and returns a full
// path to the README.md file itself. Directory names
// are always prefixed with keg-node.
func MkTemp() (string, error) {
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

// Import moves the nodedir into the KEG directory for the kegid giving
// it the nodeid name. Import will fail if the given nodeid already
// existing the the target KEG.
func Import(from, to, nodeid string) error {
	to = path.Join(to, nodeid)
	if fs.Exists(to) {
		return fs.ErrorExists{to}
	}
	return os.Rename(from, to)
}
