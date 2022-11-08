package read

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rwxrob/keg/node/parse"
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
	return parse.Title(f)
}
