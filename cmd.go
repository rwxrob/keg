// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package keg

import (
	"embed"
	"log"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/emb"
	"github.com/rwxrob/fs"
	"github.com/rwxrob/help"
	"github.com/rwxrob/keg/node"
	"github.com/rwxrob/term"
	"github.com/rwxrob/vars"
)

//go:embed files
var files embed.FS

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
	emb.FS = files
	emb.Top = "files"
}

var Cmd = &Z.Cmd{
	Name:      `keg`,
	Summary:   `manage knowledge exchange graphs (KEG)`,
	Version:   `v0.0.0`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Site:      `rwxrob.tv`,
	Source:    `git@github.com:rwxrob/keg.git`,
	Issues:    `github.com/rwxrob/keg/issues`,

	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, vars.Cmd, emb.Cmd, node.Cmd, dirCmd,
	},

	Shortcuts: Z.ArgMap{
		`current`: {`var`, `get`, `current`},
		`set`:     {`var`, `set`},
		`map`:     {`conf`, `query`, `.keg.map`},
		`pegn`:    {`emb`, `cat`, `kegml.pegn`},
	},

	ConfVars: true,

	Description: `
		The **{{.Name}}** command composes together the branches and
		commands used to search, read, create, and share knowledge on the
		free, decentralized, protocol-agnostic, world-wide, Knowledge
		Exchange Grid, a modern replacement for the very broken WorldWideWeb
		(see keg.pub for more).`,

	Call: func(x *Z.Cmd, args ...string) error {
		log.Print("made")

		return nil
	},
}

var dirCmd = &Z.Cmd{
	Name: `dir`,
	Call: func(x *Z.Cmd, args ...string) error {
		curkeg, err := x.Caller.Get(`current`)
		if err != nil {
			return err
		}
		curdir, err := x.Caller.C(`map.` + curkeg)
		if err != nil {
			return err
		}
		term.Print(fs.Tilde2Home(curdir))
		return nil
	},
}
