// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package keg

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/choose"
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
	Aliases:   []string{`kn`},
	Summary:   `manage knowledge exchange graphs (KEG)`,
	Version:   `v0.0.0`,
	Copyright: `Copyright 2021 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Site:      `rwxrob.tv`,
	Source:    `git@github.com:rwxrob/keg.git`,
	Issues:    `github.com/rwxrob/keg/issues`,

	Commands: []*Z.Cmd{
		help.Cmd, conf.Cmd, vars.Cmd,
		dexCmd, createCmd, currentCmd, dirCmd, deleteCmd, editCmd,
		latestCmd, titleCmd,
	},

	Shortcuts: Z.ArgMap{
		`set`: {`var`, `set`},
	},

	ConfVars: true,

	Description: `
		The {{cmd .Name}} command is for personal and public knowledge
		management as a Knowledge Exchange Graph (sometimes called "personal
		knowledge graph" or "zettelkasten"). Using {{cmd .Name}} you can
		create and share knowledge on the free, decentralized,
		protocol-agnostic, world-wide, Knowledge Exchange Grid.`,
}

var currentCmd = &Z.Cmd{
	Name: `current`,
	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		term.Print(keg.Name)
		return nil
	},
}

var titleCmd = &Z.Cmd{
	Name:    `titles`,
	Aliases: []string{`title`},
	Summary: `find titles containing keyword`,
	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, "")
		}
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		str := strings.Join(args, " ")
		var dex *Dex
		dex, err = ReadDex(keg.Path)
		if err != nil {
			return err
		}
		if term.IsInteractive() {
			fmt.Print(dex.WithTitleText(str))
		} else {
			fmt.Print(dex.WithTitleText(str).AsIncludes())
		}
		return nil
	},
}

var dirCmd = &Z.Cmd{
	Name:    `dir`,
	Aliases: []string{`d`},
	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		term.Print(keg.Path)
		return nil
	},
}

var deleteCmd = &Z.Cmd{
	Name:     `delete`,
	Summary:  `delete node by ID from current keg`,
	Aliases:  []string{`del`},
	Usage:    `(help|INTEGER_NODE_ID|last)`,
	MinArgs:  1,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		id := args[0]
		if id == "last" {
			id = Last(keg.Path)
		}
		_, err = strconv.Atoi(id)
		if err != nil {
			return x.UsageError()
		}
		dir := filepath.Join(keg.Path, id)
		log.Println("deleting", dir)
		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
		err = MakeDex(keg.Path)
		if err != nil {
			return err
		}
		return nil
	},
}

func current(x *Z.Cmd) (*Local, error) {
	var name, dir string

	// if we have an env it beats config settings
	name = os.Getenv(`KEG_CURRENT`)
	dir, _ = x.C(`map.` + name)
	if !(dir == "" || dir == "null") {
		dir = fs.Tilde2Home(dir)
		return &Local{Path: dir, Name: name}, nil
	}

	// check vars and conf
	name, _ = x.Get(`current`)
	if name != "" {
		dir, _ = x.C(`map.` + name)
		if !(dir == "" || dir == "null") {
			dir = fs.Tilde2Home(dir)
			return &Local{Path: dir, Name: name}, nil
		}
	}

	// map entry that matches executable
	dir, _ = x.C(`map.` + Z.ExeName)
	if !(dir == "" || dir == "null") {
		dir = fs.Tilde2Home(dir)
		return &Local{Path: dir, Name: Z.ExeName}, nil
	}

	// check if current working directory has a keg
	dir, _ = os.Getwd()
	if fs.Exists(filepath.Join(dir, `keg`)) {
		name = filepath.Base(dir)
		return &Local{Path: dir, Name: name}, nil
	}

	return nil, fmt.Errorf("no kegs found") // FIXME with better error
}

var dexCmd = &Z.Cmd{
	Name:     `dex`,
	Commands: []*Z.Cmd{help.Cmd},
	Params:   []string{`update`},
	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		if len(args) > 0 && args[0] == `update` {
			return MakeDex(keg.Path)
		}
		return nil
	},
}

var latestCmd = &Z.Cmd{
	Name:     `latest`,
	Aliases:  []string{`last`},
	Summary:  `show last nodes changed (markdown)`,
	MaxArgs:  1,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, args ...string) error {
		var err error
		n := 1
		if len(args) > 0 {
			n, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		}
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		path := filepath.Join(keg.Path, `dex/latest.md`)
		if !fs.Exists(path) {
			return fmt.Errorf("dex/latest.md file does not exist")
		}
		lines, err := file.Head(path, n)
		if err != nil {
			return err
		}
		dex, err := ParseDex(strings.Join(lines, "\n"))
		if err != nil {
			return nil
		}
		if term.IsInteractive() {
			fmt.Print(dex)
		} else {
			fmt.Print(dex.AsIncludes())
		}
		return nil
	},
}

var editCmd = &Z.Cmd{
	Name:     `edit`,
	Aliases:  []string{`e`},
	Usage:    `(help|INTEGER_NODE_ID|last|TITLEWORD)`,
	Summary:  `edit a specific node README.md file`,
	Commands: []*Z.Cmd{help.Cmd},
	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		id := args[0]
		if id == "last" {
			id = Last(keg.Path)
		} else {

			_, err := strconv.Atoi(id)
			if err != nil {
				dex, err := ReadDex(keg.Path)
				if err != nil {
					return err
				}
				choice, err := choose.Choices[DexEntry](dex.WithTitleText(strings.Join(args, " "))).Choose()
				id = strconv.Itoa(choice.N)
			}
		}
		path := filepath.Join(keg.Path, id, `README.md`)
		if !fs.Exists(path) {
			return fmt.Errorf("content node (%s) does not exist in %q", id, keg.Name)
		}
		if err := file.Edit(path); err != nil {
			return err
		}
		return MakeDex(keg.Path)
	},
}

var createCmd = &Z.Cmd{
	Name:     `create`,
	Aliases:  []string{`c`},
	Summary:  `create KEG content node`,
	MaxArgs:  1,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		readme, err := MkTempNode()
		if err != nil {
			return err
		}
		if err := file.Edit(readme); err != nil {
			return err
		}
		_, _, high := NodePaths(keg.Path)
		if high < 0 {
			return fmt.Errorf(`Can't determine last id: %v`, keg.Path)
		}
		high++
		if err := ImportNode(path.Dir(readme), keg.Path, strconv.Itoa(high)); err != nil {
			return err
		}
		if err := MakeDex(keg.Path); err != nil {
			return err
		}
		hd, _ := file.Head(filepath.Join(keg.Path, `dex`, `latest.md`), 1)
		fmt.Println(hd[0])
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
