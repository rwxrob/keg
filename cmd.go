// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package keg

import (
	"fmt"
	"log"
	"path"
	"strconv"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"github.com/rwxrob/vars"
)

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
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
		help.Cmd, conf.Cmd, vars.Cmd, nodeCmd, dirCmd, dexCmd,
	},

	Shortcuts: Z.ArgMap{
		`current`: {`var`, `get`, `current`},
		`set`:     {`var`, `set`},
		`map`:     {`conf`, `query`, `.keg.map`},
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

var dexCmd = &Z.Cmd{
	Name:     `dex`,
	Commands: []*Z.Cmd{help.Cmd},
	Params:   []string{`update`},
	Call: func(x *Z.Cmd, args ...string) error {
		curkeg, err := x.Caller.Get(`current`)
		if err != nil {
			return err
		}
		curdir, err := x.Caller.C(`map.` + curkeg)
		if err != nil {
			return err
		}
		curdir = fs.Tilde2Home(curdir)
		if len(args) > 0 && args[0] == `update` {
			return MakeDex(curdir)
		}
		return nil
	},
}

// --------------------------------------------------------------------

var nodeCmd = &Z.Cmd{
	Name:    `node`,
	Aliases: []string{`n`},
	Summary: `work with a single KEG node`,
	Commands: []*Z.Cmd{
		help.Cmd, vars.Cmd, conf.Cmd,
		nodeCreateCmd,
	},
	Shortcuts: Z.ArgMap{},
}

var nodeCreateCmd = &Z.Cmd{
	Name:     `create`,
	Aliases:  []string{`c`},
	Summary:  `create KEG node`,
	Usage:    `[help|KEG]`,
	MaxArgs:  1,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {
		var curkeg string
		var err error

		if len(args) > 0 {
			curkeg = args[0]
		} else {
			curkeg, err = x.Caller.Caller.Get(`current`)
			if err != nil {
				return err
			}
		}

		if curkeg == "" {
			return Z.MissingVar{x.Caller.Caller.Path(`current`)}
		}

		curdir, err := x.Caller.Caller.C(`map.` + curkeg)
		if curdir == "" {
			return Z.MissingConf{x.Caller.Caller.Path(`map.` + curkeg)}
		}

		curdir = fs.Tilde2Home(curdir)

		readme, err := MkTempNode()
		if err != nil {
			return err
		}

		if err := file.Edit(readme); err != nil {
			return err
		}

		_, _, high := NodePaths(curdir)
		if high < 0 {
			return fmt.Errorf(`Can't determine last id`)
		}
		high++
		if err := ImportNode(path.Dir(readme), curdir, strconv.Itoa(high)); err != nil {
			return err
		}
		MakeDex(curdir)
		return nil
	},
}

// ----------------------------- node ast -----------------------------

/*
var nodeParseCmd = &Z.Cmd{
	Name:    `parse`,
	Summary: `parse/print semantic node content`,
	Usage:   `[TYPE [FILTER|FILE|DIR]]`,
	Commands: []*Z.Cmd{help.Cmd, conf.Cmd, vars.Cmd,
		yamlCmd, jsonCmd, xmlCmd,
	},
	ConfVars: true,
	VarDefs:  Z.VarVals{`nl`: `KEGNL`},
	Shortcuts: Z.ArgMap{
		`get`:  {`var`, `get`},
		`pegn`: {`emb`, `cat`, `kegml.pegn`},
	},
	Params: []string{
		`title`, `heading`, `block`, `include`, `incfile`, `incnode`,
		`bulleted`, `numbered`, `figure`, `fenced`, `tex`, `quote`, `raw`,
		`ref`, `refs`, `link`, `linkfile`, `linknode`, `tags`, `tag`, `div`,
		`para`, `bullet`, `number`, `span`, `inflect`, `bold`, `verbatim`, `math`,
		`deleted`, `squoted`, `dquoted`, `quoted`, `bracketed`, `parens`,
		`braced`, `angled`, `url`, `longdash`, `shortdash`, `plain`, `ellipsis`,
		`word`,
	},
	Description: `
		The {{ cmd .Name }} command parses and prints different (semantic)
		parts of the KEG node. Matches are printed one per line with any
		line returns replaced with {{ pre KEGNL }} (which can be changed
		with the {{ cmd "set nl" }} command.

		The first parameter indicates the type of parsed content wanted from
		the KEGML file. Type names come from the supported KEGML PEGN
		specification available for reference from the {{ cmd "pegn" }}
		command.

		The second argument indicates the node (or nodes) to parse by KEG
		node identifier or scope filter. See the {{ cmd "keg" }} command
		help for more information about KEG.

		The second argument may also simply be a file system path to a file
		or directory containing a README.md file.

		If the second argument is omitted, the current node is assumed
		{{ pre "set current" }}. If no current node it set, the parent caller's
		{{ pre "current" }} value is used (if it exists). If even then no
		current node can be resolved, the README.md file within the current
		working directory is assumed.

	`,
}
*/
