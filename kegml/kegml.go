package kegml

import (
	"unicode"

	"github.com/rwxrob/pegn"
)

var Title = pegn.Scan{

	PEGN: `'#' SP <(rune){1,70}>`,
	Info: `some explanation here`,
	Scan: func(s pegn.Scanner, out *[]rune) bool {

		// '#'
		if !s.Scan() {
			return false
		}

		// SP
		if !s.Scan() {
			return false
		}

		// <(rune){1,70}>
		var i int
		for i = 0; s.Scan() && i < 70; i++ {
			r := s.Rune()
			if r == '\n' {
				break
			}
			if !(unicode.In(r, unicode.PrintRanges...) || r == ' ') {
				return false
			}
			if out != nil {
				*out = append(*out, r)
			}
		}

		if s.Finished() ||
			(s.Rune() == '\n' && s.Scan() && s.Rune() == '\n') ||
			(s.Rune() == '\n' && s.Finished()) {
			return true
		}

		return false
	},
}
