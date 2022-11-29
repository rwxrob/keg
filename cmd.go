// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package keg

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/dir"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"github.com/rwxrob/vars"
)

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
}

// has to stay here because needs vars package from x
func current(x *Z.Cmd) (*Local, error) {
	var name, dir string

	// if we have an env it beats config settings
	name = os.Getenv(`KEG_CURRENT`)
	if name != "" {
		dir, _ = x.C(`map.` + name)
		if !(dir == "" || dir == "null") {

			dir = fs.Tilde2Home(dir)
			if fs.NotExists(dir) {
				return nil, fs.ErrNotExist{dir}
			}
			docsdir := filepath.Join(dir, `docs`)
			if fs.Exists(docsdir) {
				dir = docsdir
			}
			return &Local{Path: dir, Name: name}, nil
		}
	}

	// check if current working directory has a keg
	dir, _ = os.Getwd()
	if fs.Exists(filepath.Join(dir, `keg`)) {
		name = filepath.Base(dir)
		if name == `docs` {
			name = filepath.Base(filepath.Dir(dir))
		}
		return &Local{Path: dir, Name: name}, nil
	}

	// check if current working directory has a docs/keg
	dir, _ = os.Getwd()
	if fs.Exists(filepath.Join(dir, `docs`, `keg`)) {
		name = filepath.Base(dir)
		dir = filepath.Join(dir, `docs`)
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

var Cmd = &Z.Cmd{
	Name:      `keg`,
	Aliases:   []string{`kn`},
	Summary:   `create and manage knowledge exchange graphs`,
	Version:   `v0.6.0`,
	Copyright: `Copyright 2022 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Site:      `rwxrob.tv`,
	Source:    `git@github.com:rwxrob/keg.git`,
	Issues:    `github.com/rwxrob/keg/issues`,

	Commands: []*Z.Cmd{
		editCmd, help.Cmd, conf.Cmd, vars.Cmd,
		dexCmd, createCmd, currentCmd, dirCmd, deleteCmd,
		lastCmd, changesCmd, titleCmd, initCmd, randomCmd,
		importCmd,
	},

	Shortcuts: Z.ArgMap{
		`set`:    {`var`, `set`},
		`sample`: {`create`, `sample`},
	},

	ConfVars: true,

	Description: `
		The {{aka}} command is for personal and public knowledge management
		as a Knowledge Exchange Graph (sometimes called "personal knowledge
		graph" or "zettelkasten"). Using {{cmd .Name}} you can create,
		update, search, and organize everything that passes through your
		brain that you may want to recall later, for whatever reason: school,
		training, team knowledge, or publishing a paper, article, blog, or
		book.
		
		Getting Started
		
		1. Create a directory and change into it
		2. Run the {{cmd "init"}} command
		3. Update the YAML file it opens
		4. Exit your editor
		5. List contents of directory to see what was created
		6. Run the {{cmd "create sample"}} command to create your first node
		7. Read and understand the sample
		8. Exit your editor
		9. Check your index with {{cmd "changes"}} or {{cmd "titles"}}
		10. Repeat 6-9 creating several nodes (optionally omitting {{cmd "sample"}})
		11. Search titles with the {{cmd "titles"}} command
		12. Edit node with keywords with {{cmd "edit WORD"}} command
		13. Notice that {{cmd "edit"}} is the default (ex: {{cmd .Name}} WORD)
		
		Learning KEG Markup Language
		
		Use the {{cmd "create sample"}} command to automatically create
		a new content node sample that explains everything about the KEG
		Markup Language (KEGML). You can delete it later after reading it.
		Or, you can use it instead of just {{cmd "create"}} (which gives
		you a blank) to help you remember how to write KEGML until you get
		proficient enough not to have to look it up every time.
		
		For more about the emerging KEG 2023-01 specification and how to
		create content that complies for knowledge exchange and publication
		(while we work more on linting and validation within the {{cmd
		.Name}} command) have a look at https://github.com/rwxrob/keg-spec
		
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
		2. The {{pre "docs"}} directory in current working if found
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

		var dex *Dex
		str := strings.Join(args, " ")
		dex, err = ReadDex(keg.Path)
		if err != nil {
			return err
		}

		if term.IsInteractive() {
			Z.Page(dex.WithTitleText(str).Pretty())
			return nil
		}

		fmt.Print(dex.WithTitleText(str).AsIncludes())
		return nil
	},
}

var dirCmd = &Z.Cmd{
	Name:     `dir`,
	Aliases:  []string{`d`},
	MaxArgs:  1,
	Summary:  `print path to directory of current keg or node`,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		if len(args) > 0 {
			dex, _ := ReadDex(keg.Path)
			choice := dex.ChooseWithTitleText(strings.Join(args, " "))
			term.Print(filepath.Join(keg.Path, strconv.Itoa(choice.N)))
			return nil
		}

		term.Print(keg.Path)

		return nil
	},
}

var deleteCmd = &Z.Cmd{
	Name:     `delete`,
	Summary:  `delete node from current keg`,
	MinArgs:  1,
	Aliases:  []string{`del`, `rm`},
	Usage:    `(help|INTEGER_NODE_ID|last|same)`,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		id := args[0]
		if id == "same" {
			if n := LastChanged(keg.Path); n != nil {
				id = n.ID()
			}
		}

		if id == "last" {
			if n := Last(keg.Path); n != nil {
				id = n.ID()
			}
		}

		if _, err = strconv.Atoi(id); err != nil {
			return x.UsageError()
		}

		dir := filepath.Join(keg.Path, id)
		log.Println("deleting", dir)

		if err := os.RemoveAll(dir); err != nil {
			return err
		}

		if err := MakeDex(keg.Path); err != nil {
			return err
		}

		return Publish(keg.Path)
	},
}

var dexCmd = &Z.Cmd{
	Name:     `dex`,
	Commands: []*Z.Cmd{help.Cmd, dexUpdateCmd},
	Summary:  `work with indexes`,
}

var dexUpdateCmd = &Z.Cmd{
	Name:     `update`,
	Commands: []*Z.Cmd{help.Cmd},
	Summary:  `update dex/changes.md and dex/nodes.tsv`,
	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller.Caller) // keg dex update
		if err != nil {
			return err
		}

		return MakeDex(keg.Path)
	},
}

var lastCmd = &Z.Cmd{
	Name:     `last`,
	Usage:    `[help|dir|id|title|time]`,
	Params:   []string{`dir`, `id`, `title`, `time`},
	MaxArgs:  1,
	Summary:  `show last created node`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} command shows information about the last content node
		that was created, which is assumed to be the one with the highest
		integer identifier within the current keg directory. By default the
		colorized form is displayed to interactive terminals and a KEGML
		include link when non-interactive (assuming !! from vim, for example).

		* {{pre "dir"}} shows only the full directory path
		* {{pre "id"}} shows only the node ID
		* {{pre "title"}} shows only the title
		* {{pre "time"}} shows only the time of last change

		Note that this is different than the latest {{cmd "changes"}} command.

	`,

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		last := Last(keg.Path)

		if len(args) == 0 {
			if term.IsInteractive() {
				fmt.Print(last.Pretty())
			} else {
				fmt.Print(last.MD())
			}
			return nil
		}

		switch args[0] {
		case `dir`:
			term.Print(filepath.Join(keg.Path, last.ID()))
		case `time`:
			term.Print(last.U.Format(IsoDateFmt))
		case `title`:
			term.Print(last.T)
		case `id`:
			term.Print(last.ID())
		}

		return nil
	},
}

var ChangesDefault = 5

var changesCmd = &Z.Cmd{
	Name:     `changes`,
	Aliases:  []string{`changed`},
	Usage:    `[help|COUNT|default|set default COUNT]`,
	Summary:  `show most recent n nodes changed`,
	UseVars:  true,
	Commands: []*Z.Cmd{help.Cmd, vars.Cmd},

	Shortcuts: Z.ArgMap{
		`default`: {`var`, `get`, `default`},
		`set`:     {`var`, `set`},
	},

	Call: func(x *Z.Cmd, args ...string) error {
		var err error
		n := ChangesDefault

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

		path := filepath.Join(keg.Path, `dex/changes.md`)
		if !fs.Exists(path) {
			return fmt.Errorf("dex/changes.md file does not exist")
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
			return nil
		}

		fmt.Print(dex.AsIncludes())
		return nil
	},
}

//go:embed testdata/samplekeg/keg
var DefaultInfoFile string

//go:embed testdata/samplekeg/0/README.md
var DefaultZeroNode string

var initCmd = &Z.Cmd{
	Name:     `init`,
	Usage:    `[help]`,
	Summary:  `initialize current working dir as new keg`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} command creates a {{pre "keg"}} YAML file in the
		current working directory and opens it up for editing.

		{{aka}} also creates a **zero node** (/0) typically used for
		linking to planned content from other content nodes.

		Finally, {{aka}} creates the {{pre "dex/changes.md"}} and
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

		if fs.NotExists(`0/README.md`) {
			if err := file.Overwrite(`0/README.md`, DefaultZeroNode); err != nil {
				return err
			}
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
	Params:   []string{`last`, `same`},
	Usage:    `(help|ID|last|same|TITLEWORD)`,
	Summary:  `choose and edit a specific node (default)`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} command opens a content node README.md file for editing.
		It is the default command when no other arguments match other
		commands. Nodes can be identified by integer ID, TITLEWORD contained
		in the title, or the special {{pre "last"}} (last created) or {{pre
		"same"}} (last updated) parameters.

		For TITLEWORD if more than one match is found the user is prompted
		to choose between them. Otherwise, the match is opened in the
		EDITOR. See rwxrob/fs.file.Edit for more about how editor is resolved.

	`,

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

		switch id {

		case "same":
			if n := LastChanged(keg.Path); n != nil {
				id = n.ID()
			}

		case "last":
			if n := Last(keg.Path); n != nil {
				id = n.ID()
			}

		default:
			_, err := strconv.Atoi(id)

			if err != nil {

				dex, err := ReadDex(keg.Path)
				if err != nil {
					return err
				}

				key := strings.Join(args, " ")
				choice := dex.ChooseWithTitleText(key)
				if choice == nil {
					return fmt.Errorf("unable to choose a title")
				}

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

		if file.IsEmpty(path) {
			if err = os.RemoveAll(filepath.Dir(path)); err != nil {
				return err
			}
		}

		// FIXME: shouldn't make the entire dex every time
		if err := MakeDex(keg.Path); err != nil {
			return err
		}

		return Publish(keg.Path)
	},
}

var createCmd = &Z.Cmd{
	Name:     `create`,
	Aliases:  []string{`c`},
	Params:   []string{`sample`},
	Summary:  `create and edit content node`,
	MaxArgs:  1,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		entry, err := MakeNode(keg.Path)
		if err != nil {
			return err
		}

		if len(args) > 0 && args[0] == `sample` {
			if err := WriteSample(keg.Path, entry); err != nil {
				return err
			}
		}

		if err := Edit(keg.Path, entry.N); err != nil {
			return err
		}

		path := filepath.Join(keg.Path, entry.ID(), `README.md`)

		if file.IsEmpty(path) {
			if err = os.RemoveAll(filepath.Dir(path)); err != nil {
				return err
			}
			return nil
		}

		if err := DexUpdate(keg.Path, entry); err != nil {
			return err
		}

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

var randomCmd = &Z.Cmd{
	Name:     `random`,
	Aliases:  []string{`rand`},
	Usage:    `[help|title|id|dir|edit]`,
	Params:   []string{`title`, `id`, `dir`, `edit`},
	MaxArgs:  1,
	Summary:  `return random node, gamify content editing`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} command randomizes the selection of a single node and
		returns the title, id, or directory; or opens the editor on a random
		node.

		One of the core tenets of the Zettelkasten approach is regularly and
		randomly reviewing the knowledge that is stored in it to bring it to
		the forefront of your mind so that it can inspire new ideas. Looking
		at a random content node is one way to accomplish this and break
		writers block by giving you something random to focus on to get you
		started.

    Defaults to {{pre "edit"}} if no argument given.
	`,

	Call: func(x *Z.Cmd, args ...string) error {
		if len(args) == 0 {
			args = append(args, `edit`)
		}
		keg, err := current(x.Caller)
		if err != nil {
			return err
		}
		dex, err := ReadDex(keg.Path)
		r := dex.Random()
		switch args[0] {
		case `id`:
			term.Print(r.N)
		case `title`:
			term.Print(r.T)
		case `edit`:
			return editCmd.Call(x, strconv.Itoa(r.N))
		case `dir`:
			term.Print(filepath.Join(strconv.Itoa(r.N)))
		}
		return nil
	},
}

var importCmd = &Z.Cmd{
	Name:     `import`,
	Usage:    `[help|(DIR|NODEDIR)...]`,
	Summary:  `import nodes into current keg`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} command imports a specific NODEDIR or all the apparent
		node directories within DIR into the current node. If no argument is
		passed, imports the current working directory into the current keg.
		If any of the arguments end in an integer they are assumed to be
		node directories. Arguments without a base integer are assumed to be
		directories containing node directories with integer identifiers.

		This command is useful when indirectly migrating nodes from one keg
		into another by way of an intermediary directory (like {{pre "tmp"}})

		Currently, there is no resolution of links within any of the
		imported nodes. Each node to be imported should be checked
		individually to ensure that any dependencies are met or adjusted.

	`,

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			d := dir.Abs()
			if d == "" {
				return fmt.Errorf("unable to determine absolute path to current directory")
			}
			args = append(args, d)
		}

		if err := Import(keg.Path, args...); err != nil {
			return err
		}

		if err := MakeDex(keg.Path); err != nil {
			return err
		}

		return Publish(keg.Path)

	},
}
