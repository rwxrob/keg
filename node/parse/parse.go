package parse

import (
	"github.com/rwxrob/keg/kegml"
	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/scanner"
)

// Title parses a KEG node title from a KEGML string, []byte, or
// io.Reader using a pegn.Scanner.
func Title(a any) (string, error) {
	b := make([]rune, 0, 70)
	s := scanner.New(a)
	if !kegml.Title.Scan(s, &b) {
		return "", pegn.ScanError{s, kegml.Title}
	}
	return string(b), nil
}
