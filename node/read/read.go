package read

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/keg/kegml"
	"github.com/rwxrob/pegn/scanner"
)

// Title reads a KEG node title from KEGML file.
func Title(path string) (string, error) {
	if !strings.HasSuffix(path, `README.md`) {
		path = filepath.Join(path, `README.md`)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	s := scanner.New(f)
	nd := kegml.ParseTitle(s)
	if nd == nil {
		return "", s
	}
	return nd.V, nil
}
