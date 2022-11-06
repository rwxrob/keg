package scan

import (
	"io"
	"unicode"

	"github.com/rwxrob/pegn"
)

// Title  <-- '#' SP <(rune){1,70}>
func Title(s pegn.Scanner, w io.Writer) bool {

	// '#'
	if !s.Scan() {
		return false
	}

	// SP
	if !s.Scan() {
		return false
	}

	// <(rune){1,70}>
	buf := make([]rune, 0, 70)
	var i int
	for i = 0; s.Scan() && i < 70; i++ {
		r := s.Rune()
		if r == '\n' {
			break
		}
		if !(unicode.In(r, unicode.PrintRanges...) || r == ' ') {
			return false
		}
		buf = append(buf, r)
	}

	if s.Finished() ||
		(s.Scan() && s.Rune() == '\n' && s.Finished()) ||
		(s.Scan() && s.Rune() == '\n') {
		if w != nil {
			w.Write([]byte(string(buf)))
		}
		return true
	}

	return false

}
