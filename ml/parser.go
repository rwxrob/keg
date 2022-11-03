package mark

import (
	"log"
	"unicode"

	"github.com/rwxrob/scan"
	"github.com/rwxrob/structs/tree"
	"github.com/rwxrob/to"
)

const (
	Unknown = iota + 1
	Tag
)

var DefaultParser = &Parser{}

// Parser parses KEG Mark (KEGML) and will populate an AST node tree (T)
// if defined. Otherwise, simply parses and validates. Parser embeds
// a scan.R (S) which should usually be used directly for testing
// and debugging purposes only. However, occasionally, it may be
// beneficial to advance the scanner over non-Mark data outside the
// scope of this parser and then continue parsing as if the skipped
// section was not a part of the original document data.
type Parser struct {
	// TODO add multiple errors to Parser rather than just first
	S      scan.R
	T      *tree.E[string]
	Errors []error

	cur *tree.Node[string]
}

// add creates a new node in the node tree (if defined) containing the
// string from the specified position to the current position.
func (p *Parser) add(typ, beg int) {
	if p.T == nil {
		return
	}
	p.cur.Add(typ, string(p.S.B[beg:p.S.P]))
	return
}

// Parse is a convenience function that passes it's input to
// DefaultParser and returns a pointer to it containing the parsed node
// tree (T) and whether the parse was successful.  For more
// control over the parsing process create a Parser instead. For just
// syntax validation consider Check instead.
func Parse(in any) (*Parser, bool) {
	// TODO provide the types to the tree
	DefaultParser.T = tree.New[string]()
	return DefaultParser, DefaultParser.Parse(in)
}

// Check calls Parse from the DefaultParser with an empty node tree to
// provide a syntax check of the input returning any errors.
func Check(in any) []error {
	DefaultParser.T = nil
	DefaultParser.Parse(in)
	return DefaultParser.Errors
}

// Parse parses an io.Reader, string, []byte or []rune returning true if
// parsing was successful. If not, false is returned and Errors is
// populated with one or more errors during parsing. If the Parser's
// node tree is not nil it will be populated with nodes based on the
// capture rules specified in mark.pegn. Document is the root syntax
// rule and is the first parsed.
func (p *Parser) Parse(in any) bool {
	p.S.B = to.Bytes(in)
	return p.Document()
}

// please keep the lowest leaves at the top

func (p *Parser) EOB() bool {
	if p.S.Peek(EOB) {
		p.S.P += 2
		return true
	}
	return false
}

func (p *Parser) ULetter() bool {
	if !p.S.Scan() {
		return false
	}
	if unicode.IsLetter(p.S.R) {
		return true
	}
	p.S.P = p.S.LP
	return false
}

func (p *Parser) UPrint() bool {
	if !p.S.Scan() {
		return false
	}
	if unicode.IsPrint(p.S.R) {
		return true
	}
	p.S.P = p.S.LP
	return false
}

func (p *Parser) Emoji() bool {
	if !p.S.Scan() {
		return false
	}
	// FIXME: create and check a unicode.RangeTable for emojis
	return unicode.IsGraphic(p.S.R) && !unicode.IsSpace(p.S.R)
}

func (p *Parser) Text() bool {
	var n int
	for p.UPrint() {
		n++
	}
	if n > 0 {
		return true
	}
	return false
}

func (p *Parser) Tag() bool {
	p.S.Scan()
	if p.S.R != '#' {
		p.S.P = p.S.LP
		return false
	}
	beg := p.S.P
	var n int
	for p.ULetter() {
		n++
		if n > 20 {
			p.S.Error(TagTooLong{})
		}
	}
	if n > 0 {
		p.add(Tag, beg)
		return true
	}
	return false
}

func (p *Parser) Mono() bool {
	beg := p.S.P
	if !p.S.Scan() {
		return false
	}
	if p.S.R != BKTICK {
		p.S.P = beg
		p.S.Error(Expected{BKTICK})
		return false
	}
	p.S.Log()
	if !p.Text() {
		p.S.P = beg
		return false
	}
	p.S.Log()
	if p.S.R != BKTICK {
		p.S.P = beg
		p.S.Error(Expected{BKTICK})
		return false
	}
	return true
}

func (p *Parser) Document() bool {
	log.Println("would parse document, not yet implemented")
	return true
}
