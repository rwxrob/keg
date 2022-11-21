package keg

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/rwxrob/json"
)

const IsoDateFmt = `2006-01-02T15:04:05Z`
const IsoDateFmtMD = `2006-01-02 15:04:05Z`
const IsoDateExpStr = `\d\d\d\d-\d\d-\d\dT\d\d:\d\d:\d\dZ`
const IsoDateExpStrMD = `\d\d\d\d-\d\d-\d\d \d\d:\d\d:\d\dZ`

// DexEntry represents a single line in the dex/nodes.md or
// dex/nodes.json file. All three fields are always required.
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

// String fulfills the fmt.Stringer interface as JSON. Any error returns
// a "null" string.
func (e DexEntry) String() string {
	byt, err := e.MarshalJSON()
	if err != nil {
		return "null"
	}
	return string(byt)
}

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
		"* %v [%v](/%v)\n",
		e.U.Format(IsoDateFmtMD),
		e.T, e.N,
	)
}

// Dex is a collection of DexEntry structs. This allows mapping methods
// for its serialization to either the dex/nodes.md or dex/nodes.json or
// other output formats. Marshaling a Dex is unique in that each
// DexEntry is on its own line and there is never any HTML escaping of
// any JSON value (unlike standard JSON marshaling).
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
func (e Dex) String() string {
	byt, err := e.MarshalJSON()
	if err != nil {
		return "null"
	}
	return string(byt)
}

// MD renders the entire Dex as a Markdown list suitable for the
// dex/nodes.md file.
func (e Dex) MD() string {
	var str string
	for _, entry := range e {
		str += entry.MD()
	}
	return str
}
