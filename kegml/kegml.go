package kegml

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/ast"
	"github.com/rwxrob/pegn/scanner"
)

//go:embed kegml.pegn
var PEGN string

const (
	Untyped int = iota
	Title
)

// ------------------------------- Title ------------------------------

func ScanTitle(s pegn.Scanner, buf *[]rune) bool {
	m := s.Mark()
	if !s.Scan() || s.Rune() != '#' {
		return s.Revert(m, Title)
	}
	if !s.Scan() || s.Rune() != ' ' {
		return s.Revert(m, Title)
	}
	var count int
	for s.Scan() {
		if count >= 70 {
			return s.Revert(m, Title)
		}
		r := s.Rune()
		if r == '\n' {
			if count > 0 {
				return true
			} else {
				return s.Revert(m, Title)
			}
		}
		if buf != nil {
			*buf = append(*buf, r)
		}
		count++
	}
	return true
}

func ParseTitle(s pegn.Scanner) *ast.Node {
	buf := make([]rune, 0, 70)
	if !ScanTitle(s, &buf) {
		return nil
	}
	return &ast.Node{T: Title, V: string(buf)}
}

var Scanner pegn.Scanner

func init() {
	Scanner = scanner.New()
	Scanner.SetErrFmtFunc(
		func(e error) string {
			return fmt.Sprintf("custom %q\n", e)
		})
}

// ReadTitle reads a KEG node title from KEGML file.
func ReadTitle(path string) (string, error) {
	if !strings.HasSuffix(path, `README.md`) {
		path = filepath.Join(path, `README.md`)
	}
	if err := Scanner.Open(path); err != nil {
		return "", err
	}
	nd := ParseTitle(Scanner)
	if nd == nil {
		return "", Scanner
	}
	return nd.V, nil
}
