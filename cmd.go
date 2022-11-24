// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package keg

import (
	_ "embed"
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
	Version:   `v0.1.0`,
	Copyright: `Copyright 2022 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Site:      `rwxrob.tv`,
	Source:    `git@github.com:rwxrob/keg.git`,
	Issues:    `github.com/rwxrob/keg/issues`,

	Commands: []*Z.Cmd{
		editCmd, help.Cmd, conf.Cmd, vars.Cmd,
		dexCmd, createCmd, currentCmd, dirCmd, deleteCmd,
		latestCmd, titleCmd, initCmd,
	},

	Shortcuts: Z.ArgMap{
		`set`: {`var`, `set`},
	},

	ConfVars: true,

	Description: `
		The {{aka}} command is for personal and public knowledge
		management as a Knowledge Exchange Graph (sometimes called "personal
		knowledge graph" or "zettelkasten"). Using {{cmd .Name}} you can
		create and share knowledge on the free, decentralized,
		protocol-agnostic, world-wide, Knowledge Exchange Grid.

		Run {{cmd "init"}} inside of a new directory to get started with
		a new keg. After editing the {{pre "keg"}} file you can create your
		first node with {{cmd "create"}}.

		For more about the emerging KEG 2023-01 specification and how to
		create content that complies for knowledge exchange and publication
		(while we work more on linting and validation within the {{cmd .Name}}
		command) have a look at https://github.com/rwxrob/keg-spec

		`,
}

var currentCmd = &Z.Cmd{
	Name:     `current`,
	Summary:  `show the current keg`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{cmd .Name}} command displays the current keg by name, which is
		resolved as follows:

		1. The {{pre "KEG_CURRENT"}} environment variable
		2. The current working directory if {{pre "keg"}} file found
		3. The {{pre "current"}} var setting (see {{cmd "var"}})

		Note that setting the var forces {{cmd .Name}} to always use that
		setting until it is explicitly changed or temporarily overridden
		with {{pre "KEG_CURRENT"}} environment variable.

		It is often useful to have {{pre "current"}} set to the most
		frequently used keg and then change into the working directory of
		another, less updated, keg when needed.

	`,

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
	Name:     `titles`,
	Aliases:  []string{`title`},
	Summary:  `find titles containing keyword`,
	Commands: []*Z.Cmd{help.Cmd},

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
			//fmt.Print(dex.WithTitleText(str).Pretty())
			Z.Page(dex.WithTitleText(str).Pretty())
		} else {
			fmt.Print(dex.WithTitleText(str).AsIncludes())
		}
		return nil
	},
}

var dirCmd = &Z.Cmd{
	Name:     `dir`,
	Aliases:  []string{`d`},
	Summary:  `print path to directory of current keg`,
	Commands: []*Z.Cmd{help.Cmd},

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
	Aliases:  []string{`del`, `rm`},
	Usage:    `(help|INTEGER_NODE_ID|last)`,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		id := args[0]
		if id == "last" {
			if n := Last(keg.Path); n != nil {
				id = n.ID()
			}
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
		return Publish(keg.Path)
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

	// check if current working directory has a keg
	dir, _ = os.Getwd()
	if fs.Exists(filepath.Join(dir, `keg`)) {
		name = filepath.Base(dir)
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

	return nil, fmt.Errorf("no kegs found") // FIXME with better error
}

var dexCmd = &Z.Cmd{
	Name:     `dex`,
	Commands: []*Z.Cmd{help.Cmd, dexUpdateCmd},
	Summary:  `work with indexes`,
}

var dexUpdateCmd = &Z.Cmd{
	Name:     `update`,
	Commands: []*Z.Cmd{help.Cmd},
	Summary:  `update dex/latest.md and dex/nodes.tsv`,
	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller.Caller) // keg dex update
		if err != nil {
			return err
		}
		return MakeDex(keg.Path)
	},
}

var latestCmd = &Z.Cmd{
	Name:     `latest`,
	Aliases:  []string{`last`},
	Summary:  `show last nodes changed (markdown)`,
	UseVars:  true,
	Commands: []*Z.Cmd{help.Cmd, vars.Cmd},
	Shortcuts: Z.ArgMap{
		`default`: {`var`, `get`, `default`},
		`set`:     {`var`, `set`},
	},
	Call: func(x *Z.Cmd, args ...string) error {
		var err error
		n := 1
		if len(args) > 0 {
			n, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}
		} else {
			def, err := x.Get(`default`)
			if err == nil && def != "" {
				n, err = strconv.Atoi(def)
				if err != nil {
					return err
				}
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
			fmt.Print(dex.Pretty())
		} else {
			fmt.Print(dex.AsIncludes())
		}
		return nil
	},
}

//go:embed testdata/samplekeg/keg
var DefaultInfoFile string

//go:embed testdata/samplekeg/0/README.md
var DefaultZeroNode string

var initCmd = &Z.Cmd{
	Name:     `init`,
	Usage:    `(help)`,
	Summary:  `initialize current working dir as new keg`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{cmd .Name}} command creates a {{pre "keg"}} YAML file in the
		current working directory and opens it up for editing. 

		{{cmd .Name}} also creates a **zero node** (/0) typically used for
		linking to planned content from other content nodes. 

		Finally, {{cmd .Name}} creates the {{pre "dex/latest.md"}} and 
		{{pre "dex/nodex.tsv"}} index files and updates the {{pre "keg"}} file
		update field to match the latest update (effectively the same as calling
		{{cmd "dex update"}}).

	`,

	Call: func(_ *Z.Cmd, _ ...string) error {
		if fs.NotExists(`keg`) {
			if err := file.Overwrite(`keg`, DefaultInfoFile); err != nil {
				return err
			}
		}
		if err := file.Overwrite(`0/README.md`, DefaultZeroNode); err != nil {
			return err
		}
		if err := file.Edit(`keg`); err != nil {
			return err
		}
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		if err := MakeDex(dir); err != nil {
			return err
		}
		return Publish(dir)
	},
}

var editCmd = &Z.Cmd{
	Name:     `edit`,
	Aliases:  []string{`e`},
	Usage:    `(help|INTEGER_NODE_ID|last|TITLEWORD)`,
	Summary:  `choose and edit a specific node`,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			return help.Cmd.Call(x, args...)
		}
		if !term.IsInteractive() {
			return titleCmd.Call(x, args...)
		}
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		id := args[0]
		if id == "last" {
			if n := Last(keg.Path); n != nil {
				id = n.ID()
			}
		} else {
			_, err := strconv.Atoi(id)
			if err != nil {
				dex, err := ReadDex(keg.Path)
				if err != nil {
					return err
				}
				key := strings.Join(args, " ")
				hits := dex.WithTitleText(key)
				switch len(hits) {
				case 1:
					id = strconv.Itoa(hits[0].N)
				case 0:
					return fmt.Errorf("no titles match: %v", key)
				default:
					i, _, err := choose.From(hits.PrettyLines())
					if err != nil {
						return err
					}
					if i < 0 {
						return nil
					}
					id = strconv.Itoa(hits[i].N)
				}
			}
		}
		path := filepath.Join(keg.Path, id, `README.md`)
		if !fs.Exists(path) {
			return fmt.Errorf("content node (%s) does not exist in %q", id, keg.Name)
		}
		if err := file.Edit(path); err != nil {
			return err
		}
		if err := MakeDex(keg.Path); err != nil {
			return err
		}
		return Publish(keg.Path)
	},
}

var createCmd = &Z.Cmd{
	Name:     `create`,
	Aliases:  []string{`c`},
	Summary:  `create and edit content node`,
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
			high = 0
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
		return Publish(keg.Path)
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
