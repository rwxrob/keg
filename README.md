# 🌳 Go Bonzai™ Composite Command Tree

*Create a new GitHub project using this template and change this
README.md to match your project. Make all your template changes before
making your first commit.*

[![GoDoc](https://godoc.org/github.com/rwxrob/bonzai-example?status.svg)](https://godoc.org/github.com/rwxrob/bonzai-example)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/rwxrob/bonzai-example/cmd/bonzai-example@latest
```

Composed

```go
package z

import (
	Z "github.com/rwxrob/bonzai/z"
	example "github.com/rwxrob/bonzai-example"
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
complete -C bonzai-example bonzai-example
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.

## Reminders

* Change `bonzai-example` every place to your project name (`git grep
  bonzai-example`)
* Remove anything you don't need
* Change `.github/FUNDING.yaml` to your own information
* Update `.gitignore` to your liking
* Will need to `go get -u` to update dependencies

## Other Examples

* <https://github.com/rwxrob/z> - the one that started it all
