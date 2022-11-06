package node

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/keg/kegml/scan"
	"github.com/rwxrob/pegn/scanner"
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

// ReadTitle reads a KEG node title from KEGML file.
func ReadTitle(path string) (string, error) {
	if !strings.HasSuffix(path, `README.md`) {
		path = filepath.Join(path, `README.md`)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	return ParseTitle(f)
}

// ParseTitle parses a KEG node title from a KEGML string, []byte, or
// io.Reader using a pegn.Scanner.
func ParseTitle(a any) (string, error) {
	b := new(strings.Builder)
	s := scanner.New(a)
	if !scan.Title(s, b) {
		return "", fmt.Errorf("failed to parse title") // TODO proper error
	}
	return b.String(), nil
}
