# 🌳 KEG Commands

[![GoDoc](https://godoc.org/github.com/rwxrob/keg?status.svg)](https://godoc.org/github.com/rwxrob/keg)
[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

This `keg` [Bonzai](https://github.com/rwxrob/bonzai) branch contains all KEG related commands, most of which are exported so they can be composed individually if preferred.

## Install

You can just download from the [releases page](https://github.com/rwxrob/keg/releases).

```
curl -L https://github.com/rwxrob/keg/releases/latest/download/keg-linux-amd64 -o ~/.local/bin/keg
curl -L https://github.com/rwxrob/keg/releases/latest/download/keg-darwin-amd64 -o ~/.local/bin/keg
curl -L https://github.com/rwxrob/keg/releases/latest/download/keg-darwin-arm64 -o ~/.local/bin/keg
curl -L https://github.com/rwxrob/keg/releases/latest/download/keg-windows-amd64 -o ~/.local/bin/keg

```

Or with `go`:

```
go install github.com/rwxrob/keg/cmd/keg@latest
```

You might want to create a small script to encapsulate `KEG_CURRENT` rather than changing into the directory all the time. Note that aliases and functions do not reliably work from within `vim`, only executables (which is 80% of the reason to use `keg` in the first place).

```shell
#!/bin/bash
KEG_CURRENT=zet keg "$@"
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
```

## Configuration

`map` - map of all local keg ids pointing to their directories (like PATH)

## Variables

`current` - current keg from `map`

## Build and Release Instructions

Building workflow uses the [`good`](https://github.com/rwxrob/good) Go helper tool (often composited into bonzai personal command trees (`z go`):

```
cd cmd/keg
good build
gh release create
gh release upload TAGVER build/*
```
