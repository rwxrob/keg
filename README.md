# 🌳 KEG Commands

*🚧This project is in active preliminary development. Commits are public
to allow collaboration and exploration of different directions.*

[![GoDoc](https://godoc.org/github.com/rwxrob/keg?status.svg)](https://godoc.org/github.com/rwxrob/keg)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

This `keg` [Bonzai](https://github.com/rwxrob/bonzai) branch contains
all KEG related commands, most of which are exported so they can be
composed individually if preferred.

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/rwxrob/keg/cmd/keg@latest
```

Composed

```go
package z

import (
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/keg"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, keg.Cmd},
}
```

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C keg keg
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.

## Command Line Usage

```
keg help
keg set
ket get

keg conf

keg map - print YAML config for current local keg ids and directories
keg current - print current local keg id
keg set current - sets current local keg id

keg scope - print current default scope
keg set scope SCOPE - set default search scope

keg NAME 

kn - shortcut for `keg node` (monolith only)

-----------

kn edit [NODE] - edit current or specific node (best guess at NODE)
kn create - create a new KEG node

keg nodes - print the KEGNODES file
keg nodes NODEPATTERN - list nodes matching pattern
keg update [KEG] - update KEGNODES (and /index)
keg sync - update current KEG and fetch all FOLLOWS

keg update - updates own KEGNODES, names, links, and follow cache

keg follow add - add an entry to FOLLOWS
keg follow cache - fetch fresh cache of all FOLLOWS
keg avoid add - add an entry to AVOID

keg copy|cp (KEG|NODE ...) KEGTARGET - copy a KEG into into another
keg move|mv (KEG|NODE ...) KEGTARGET - move a KEG or node into another

keg search|query|q|? SEARCH - search content based on preferences

keg [check] orphans - check and list all nodes without links of any kind
keg [check] includes [nodes|files] - check and list nodes with includes
keg [check] links [nodes|files] - check and list all interlinks between nodes
keg [check] urls [FILTER] - check and list all urls pointing to remote content

```

## Configuration

`map` - map of all local keg ids pointing to their directories (like PATH)

## Variables

`current` - current keg from `map`

## Scope Syntax

The method of pattern matching is controlled by the following variables:

* `keg.query.syntax.kegs`
* `keg.query.syntax.nodes`

Default search pattern syntax can be set by the type of content in text
nodes (but most will just set `text.ALL`):

* `keg.query.syntax.text.ALL`
* `keg.query.syntax.text.titles`
* `keg.query.syntax.text.body`
* `keg.query.syntax.text.tags`
* `keg.query.syntax.text.links`
* `keg.query.syntax.text.plain`
* `keg.query.syntax.text.emph`
* `keg.query.syntax.text.strong`
* `keg.query.syntax.text.stremph`
* `keg.query.syntax.text.fenced`
* `keg.query.syntax.text.raw`
* `keg.query.syntax.text.math`
* `keg.query.syntax.text.pre`
* `keg.query.syntax.text.sem`
* `keg.query.syntax.text.tex`
* `keg.query.syntax.text.quote`
* `keg.query.syntax.text.bullet`
* `keg.query.syntax.text.number`
* `keg.query.syntax.text.list`

Data can widely vary depending on the format of the data. The YAML,
JSON, CSV, and TSV data formats are RECOMMENDED by the spec for
all tool implementations. Eventually, data search plugins will be
supported (as autonomous binaries connected by stdin/stdout).

* `keg.query.syntax.data.ALL`
* `keg.query.syntax.data.yaml`
* `keg.query.syntax.data.json`
* `keg.query.syntax.data.csv`
* `keg.query.syntax.data.tsv`

* `keg.query.syntax.figures.ALL`

Examples

```
keg q *:*  - default, searches everything
keg q foo:* - searches everything in the KEG named 'foo'
keg q foo:bar - searches for the keyword 'bar' anywhere in 'foo' KEG
keg q foo.
```

Formal PEGN

```pegn
Scope       <-- Filter+
KegIdent    <-  Glob / Regex
Filter      <-- KegIdent NodeFilter?
NodeFilter  <--
Glob        <-- # find compatible name glob
Regex       <-- # Go compatible regular expression
```

