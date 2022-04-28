# 🌳 KEG Commands

[![GoDoc](https://godoc.org/github.com/rwxrob/keg?status.svg)](https://godoc.org/github.com/rwxrob/keg)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

This `keg` [Bonzai][https://github.com/rwxrob/bonzai] branch contains
all KEG related commands, most of which are exported so they can be
composed individually if preferred (for example, `keg.MarkupCmd` for
parsing KEGML in other converter utilities).

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
	example "github.com/rwxrob/keg"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, example.Cmd, example.BazCmd},
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
