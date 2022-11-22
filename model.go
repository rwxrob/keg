package keg

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rwxrob/json"
)

const IsoDateFmt = `2006-01-02 15:04:05Z`
const IsoDateExpStr = `\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\dZ`

// Local contains a name to full path mapping for kegs stored locally.
type Local struct {
	Name string
	Path string
}

// DexEntry represents a single line in an index (usually the latest.md
// or nodes.tsv file). All three fields are always required.
type DexEntry struct {
	U time.Time // updated
	T string    // title
	N int       // node id
}

// MarshalJSON produces JSON text that contains one DexEntry per line
// that has not been HTML escaped (unlike the default) and that uses
// a consistent DateTime format. Note that the (broken) encoding/json
// encoder is not used at all.
func (e *DexEntry) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	buf.WriteRune('{')
	buf.WriteString(`"U":"` + e.U.Format(IsoDateFmt) + `",`)
	buf.WriteString(`"N":` + strconv.Itoa(e.N) + `,`)
	buf.WriteString(`"T":"` + json.Escape(e.T) + `"`)
	buf.WriteRune('}')
	return buf.Bytes(), nil
}

func (e DexEntry) TSV() string {
	return fmt.Sprintf("%v\t%v\t%v", e.N, e.U.Format(IsoDateFmt), e.T)
}

// String fulfills the fmt.Stringer interface as a markdown link.
func (e DexEntry) String() string { return e.TSV() }

// MD returns the entry as a single Markdown list item for inclusion in
// the dex/nodex.md file:
//
//     1. Second last changed in UTC in ISO8601 (RFC3339)
//     2. Current title (always first line of README.md)
//     2. Unique node integer identifier
//
// Note that the second of last change is based on *any* file within the
// node directory changing, not just the README.md or meta files.
func (e DexEntry) MD() string {
	return fmt.Sprintf(
		"* %v [%v](/%v)",
		e.U.Format(IsoDateFmt),
		e.T, e.N,
	)
}

// Asinclude returns a KEGML include link list item without the time
// suitable for creating include blocks in node files.
func (e DexEntry) AsInclude() string {
	return fmt.Sprintf("* [%v](/%v)", e.T, e.N)
}

// Dex is a collection of DexEntry structs. This allows mapping methods
// for its serialization to different output formats.
type Dex []DexEntry

// MarshalJSON produces JSON text that contains one DexEntry per line
// that has not been HTML escaped (unlike the default).
func (d *Dex) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 0))
	buf.WriteString("[")
	for _, entry := range *d {
		byt, _ := entry.MarshalJSON()
		buf.Write(byt)
		buf.WriteString(",\n")
	}
	byt := buf.Bytes()
	byt[len(byt)-2] = ']'
	return byt, nil
}

// String fulfills the fmt.Stringer interface as JSON. Any error returns
// a "null" string.
func (e Dex) String() string { return e.TSV() }

// MD renders the entire Dex as a Markdown list suitable for the
// standard dex/latest.md file.
func (e Dex) MD() string {
	var str string
	for _, entry := range e {
		str += entry.MD() + "\n"
	}
	return str
}

// AsIncludes renders the entire Dex as a KEGML include list (markdown
// bulleted list) and cab be useful from within editing sessions to
// include from the current keg without leaving the terminal editor.
func (e Dex) AsIncludes() string {
	var str string
	for _, entry := range e {
		str += entry.AsInclude() + "\n"
	}
	return str
}

// TSV renders the entire Dex as a loadable tab-separated values file.
func (e Dex) TSV() string {
	var str string
	for _, entry := range e {
		str += entry.TSV() + "\n"
	}
	return str
}

// ByID orders the Dex from lowest to highest node ID integer.
func (e Dex) ByID() Dex {
	sort.Slice(e, func(i, j int) bool {
		return e[i].N < e[j].N
	})
	return e
}

// WithTitleText filters all nodes with titles that do not contain the text
// substring in the title.
func (e Dex) WithTitleText(keyword string) Dex {
	dex := Dex{}
	for _, d := range e {
		if strings.Index(d.T, keyword) >= 0 {
			dex = append(dex, d)
		}
	}
	return dex
}
