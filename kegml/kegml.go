package kegml

import (
	"github.com/rwxrob/pegn"
	"github.com/rwxrob/pegn/ast"
)

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
