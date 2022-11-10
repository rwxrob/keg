package kegml

import (
	"unicode"

	"github.com/rwxrob/pegn"
)

var EndBlock = pegn.Parser{
	Ident: `EndBlock`,
	PEGN:  `!. / LF !. / LF{2}`,
}

var Start = pegn.Parser{
	Ident: `Start`,
	PEGN:  `!<.`,
}

var End = pegn.Parser{
	Ident: `End`,
	PEGN:  `!.`,
}

var Title = pegn.Parser{

	Type: 1, Ident: `Title`,
	PEGN: `Start '#' SP < rune{1,70} > EndBlock`,

	Info: `A title MUST NOT have anything before it (first line) and MUST be identified with a hashtag (#) followed by a single space. The title text itself must be 1-70 runes (UNICODE printable code points or simple spaces). Titles are blocks and therefore end with an EndBlock (end of data or two line returns).`,

	Scan: func(s pegn.Scanner, out *[]rune) bool {

		// TODO check that at first position in the scanner

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

	Parse: func(s pegn.Scanner) *Node {
		buf := make([]rune, 0, 70)
		s.Scan()

		return nil
	},
}
