// Copyright 2022 Robert Muhlestein.
// SPDX-License-Identifier: Apache-2.0

package keg

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/glamour"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/choose"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/dir"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/grep"
	"github.com/rwxrob/help"
	"github.com/rwxrob/term"
	"github.com/rwxrob/to"
	"github.com/rwxrob/vars"
)

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
}

var DefColumns = 100

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
	Version:   `v0.8.0`,
	UseVars:   true,
	Copyright: `Copyright 2022 Robert S Muhlestein`,
	License:   `Apache-2.0`,
	Site:      `rwxrob.tv`,
	Source:    `git@github.com:rwxrob/keg.git`,
	Issues:    `github.com/rwxrob/keg/issues`,

	Commands: []*Z.Cmd{
		editCmd, help.Cmd, conf.Cmd, vars.Cmd,
		dexCmd, createCmd, currentCmd, dirCmd, deleteCmd,
		lastCmd, changesCmd, titlesCmd, initCmd, randomCmd,
		importCmd, grepCmd, viewCmd,
	},

	Shortcuts: Z.ArgMap{
		`set`:    {`var`, `set`},
		`get`:    {`var`, `get`},
		`sample`: {`create`, `sample`},
	},

	ConfVars: true,

	Description: `
		The {{aka}} command is for personal and public knowledge management
		as a Knowledge Exchange Graph (sometimes called "personal knowledge
		graph" or "zettelkasten"). Using {{aka}} you can create,
		update, search, and organize everything that passes through your
		brain that you may want to recall later, for whatever reason: school,
		training, team knowledge, or publishing a paper, article, blog, or
		book.

	`,

	Other: []Z.Section{
		{`Getting Started`, `
		
		The steps to create your first KEG directory are below, but first
		a little about the structure of this directory, which does not
		necessarily require this {{aka}} command to create and maintain.
		
		A KEG directory (aka "keg") is just a directory containing a
		{{pre "keg"}} YAML file and a number of directories that have a
		{{pre "README.md"}} file called **content nodes**. A node directory must
		have an incrementing integer name as would be used in a database
		table. A keg usually also has a {{pre "dex"}} directory containing
		at least two files:
		
		1. Latest changes {{pre "dex/changes.md"}}
		2. All nodes by ID {{pre "dex/nodes.tsv"}}
		
		The {{aka}} command keeps these files up to date.
		
		A special **zero node** is used by convention as a target for links
		to nodes that have yet to be created.
		
		Okay, here are the specific steps to get started by creating your
		first keg directory. If you plan on using Git or GitHub hold off on
		doing anything with git for now.
		
		1. Create a directory and change into it
		2. Run the {{aka}} {{cmd "init"}} command
		3. Update the YAML file it opens
		4. Exit your editor
		5. List contents of directory to see what was created
		6. Run the {{aka}} {{cmd "create sample"}} command to create your first node
		7. Read and understand the sample
		8. Exit your editor
		9. Check your index with {{aka}} {{cmd "changes"}} or {{aka}} {{cmd "titles"}}
		10. Repeat 6-9 creating several nodes (optionally omitting {{cmd "sample"}})
		11. Search titles with the {{aka}} {{cmd "titles"}} command
		12. Edit node with title keywords with {{aka}} {{cmd "edit WORD"}} command
		13. Edit node with grep regexp matches with {{aka}} {{cmd "grep WORD"}} command
		14. Notice that {{aka}} {{cmd "edit"}} is the default (ex: {{aka}} WORD)
		
		`,
		}, {`Git and GitHub`, `

		It's important when using Git that either the remote git repo has
		been fully created (so that {{cmd "git pull"}} will work) or that
		{{cmd "git"}} has not been run at all (no {{pre ".git"}} directory).
		Otherwise, {{aka}} will attempt to pull and fail. These instructions
		assume the reader understands {{cmd "git"}} and the {{cmd "gh"}}
		commands.
		
		Here are the steps to follow when Git and GitHub are wanted. They
		are essentially the same as {{pre "Getting Started"}} but include
		creating a GitHub repo with the {{cmd "gh"}} command afterward.
		
		1. Create a directory and change into it
		2. Run the {{aka}} {{cmd "init"}} command
		3. Update the YAML file it opens
		4. Exit your editor
		5. Create and push as Git repo with {{cmd "gh repo create"}}
		6. Continue with steps 9+ from Getting Started
		
		Alternatively, one can simply create a GitHub repo from the web site
		and {{cmd "git clone"}} it down to the local machine and then run
		{{aka}} {{cmd "init"}} from within it.

		`}, {`Learning KEG Markup Language`, `
		
		Use the {{aka}} {{cmd "create sample"}} command to automatically create
		a new content node sample that introduces the KEG Markup Language
		(KEGML). You can delete it later after reading it. Or, you can use
		it instead of just {{aka}} {{cmd "create"}} (which gives you a blank) to
		help you remember how to write KEGML until you get proficient enough
		not to have to look it up every time.
		
		For more about the emerging KEG 2023-01 specification and how to
		create content that complies for knowledge exchange and publication
		(while we work more on linting and validation within the
		{{aka}} command) have a look at https://github.com/rwxrob/keg-spec
		
		`},
	},
}

var currentCmd = &Z.Cmd{
	Name:     `current`,
	Summary:  `show the current keg`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} command displays the current keg by name, which is
		resolved as follows:
		
		1. The {{pre "KEG_CURRENT"}} environment variable
		2. The current working directory if {{pre "keg"}} file found
		2. The {{pre "docs"}} directory in current working if found
		3. The {{pre "current"}} var setting (see {{cmd "var"}})
		
		Note that setting the var forces {{aka}} to always use that
		setting until it is explicitly changed or temporarily overridden
		with {{pre "KEG_CURRENT"}} environment variable.
		
		      keg()
		      {
		        KEG_CURRENT=zet keg "$@"
		      }
		
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

var titlesCmd = &Z.Cmd{
	Name:    `titles`,
	Aliases: []string{`title`},
	Usage:   `(help|REGEXP)`,
	Summary: `find titles containing regular expression`,
	UseVars: true,

	Description: `
		The {{aka}} command returns a paged list of all titles matching the
		given regular expression. By default all regular expressions
		({{pre "REGEXP"}}) are made case insensitive by adding the prefix
		{{pre "(?i)"}} which can be explicitly overridden with {{pre "(?-i)"}} for
		one search or changed as the default by assigning it to the
		{{pre "regxpre"}} variable:

		      keg set regxpre '(?-i)'

		Note that if set, {{pre "regxpre"}} applies to *all* searches, which
		includes the {{cmd "edit"}} and {{cmd "grep"}} commands.

	`,

	Commands: []*Z.Cmd{help.Cmd, vars.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		if len(args) == 0 {
			args = append(args, "")
		}

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		var dex *Dex
		dex, err = ReadDex(keg.Path)
		if err != nil {
			return err
		}

		pre, err := x.Caller.Get(`regxpre`)
		if err != nil {
			return err
		}
		if pre == "" {
			pre = `(?i)`
		}

		re, err := regexp.Compile(pre + args[0])
		if err != nil {
			return err
		}

		if term.IsInteractive() {
			Z.Page(dex.WithTitleTextExp(re).Pretty())
			return nil
		}

		fmt.Print(dex.WithTitleTextExp(re).AsIncludes())
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
	Commands: []*Z.Cmd{help.Cmd},

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
		{{pre "dex/nodes.tsv"}} index files and updates the {{pre "keg"}} file
		{{pre "updated"}} field to match the latest update (effectively the
		same as calling {{cmd "dex update"}}).

		Also see the *Getting Started* docs in the main {{cmd "help"}} command.

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
		in the title, or the special {{pre "last"}} (last created) or 
		{{pre "same"}} (last updated) parameters.

		For TITLEWORD if more than one match is found the user is prompted
		to choose between them. Otherwise, the match is opened in the
		EDITOR. See rwxrob/fs.file.Edit for more about how editor is resolved.

	`,

	Call: func(x *Z.Cmd, args ...string) error {

		if len(args) == 0 {
			return help.Cmd.Call(x, args...)
		}

		if !term.IsInteractive() {
			return titlesCmd.Call(x, args...)
		}

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		id := args[0]
		var entry *DexEntry

		switch id {

		case "same":
			if entry = LastChanged(keg.Path); entry != nil {
				id = entry.ID()
			}

		case "last":
			if entry = Last(keg.Path); entry != nil {
				id = entry.ID()
			}

		default:

			dex, err := ReadDex(keg.Path)
			if err != nil {
				return err
			}

			idn, err := strconv.Atoi(id)

			if err == nil {
				entry = dex.Lookup(idn)
			} else {

				pre, err := x.Caller.Get(`regxpre`)
				if err != nil {
					return err
				}
				if pre == "" {
					pre = `(?i)`
				}

				re, err := regexp.Compile(pre + args[0])
				if err != nil {
					return err
				}

				entry = dex.ChooseWithTitleTextExp(re)
				if entry == nil {
					return fmt.Errorf("unable to choose a title")
				}

				id = strconv.Itoa(entry.N)
			}
		}

		path := filepath.Join(keg.Path, id, `README.md`)

		if !fs.Exists(path) {
			return fmt.Errorf("content node (%s) does not exist in %q", id, keg.Name)
		}

		btime := fs.ModTime(path)

		if err := file.Edit(path); err != nil {
			return err
		}

		if file.IsEmpty(path) {
			if err = os.RemoveAll(filepath.Dir(path)); err != nil {
				return err
			}
			if err := DexRemove(keg.Path, entry); err != nil {
				return err
			}
			return Publish(keg.Path)
		} else {
			if err := DexUpdate(keg.Path, entry); err != nil {
				return err
			}
		}

		atime := fs.ModTime(path)
		if !atime.After(btime) {
			return nil
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

type grepChoice struct {
	hit grep.Result
	str string
}

func (c grepChoice) String() string { return c.str }

var grepCmd = &Z.Cmd{
	Name:     `grep`,
	Usage:    `(help|REGEXP)`,
	MinArgs:  1,
	Summary:  `grep regular expression out of all nodes`,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} performs a simple regular expression grep of all node
		README.md files. (Does not depend on host {{pre "grep"}} command. By
		default, all regular expressions are case-sensitive (unlike default
		{{cmd "title"}} or {{cmd "edit"}} commands).  This can be
		explicitly overridden for a given search by adding {{pre "(?i)"}} or
		for all searches by setting the global {{pre "regxpre"}} variable to
		the same:
		
		      keg grep set regxpre '(?i)'
		
		Note that if set, {{pre "regxpre"}} applies to *all* searches, which
		includes the {{cmd "titles"}} and {{cmd "edit"}} commands.
		
		`,

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		dirs, _, _ := fs.IntDirs(keg.Path)
		dpaths := []string{}
		for _, d := range dirs {
			dpaths = append(dpaths, filepath.Join(d.Path, `README.md`))
		}

		// figure out columns (yes it's complicated)
		col := int(term.WinSize.Col) // only > 0 for interactive terminals
		if col <= 0 {

			colstr, err := x.Caller.Get(`columns`)
			if err != nil {
				return err
			}

			if colstr != "" {
				col, err = strconv.Atoi(colstr)
				if err != nil {
					return err
				}
			}

			if col <= 0 {
				col = DefColumns
			}
		}

		col -= 14
		results, err := grep.This(args[0], col, dpaths...)
		if err != nil {
			return err
		}

		if term.IsInteractive() {

			var choices []grepChoice
			for _, hit := range results.Hits {
				id := filepath.Base(filepath.Dir(hit.File))
				match := to.CrunchSpaceVisible(hit.Text[hit.TextBeg:hit.TextEnd])
				before := to.CrunchSpaceVisible(hit.Text[0:hit.TextBeg])
				after := to.CrunchSpaceVisible(hit.Text[hit.TextEnd:])
				width := len(match) + len(before) + len(after)
				if width > col {
					chop := (width - col) / 2
					lafter := len(after)
					lbefore := len(before)
					switch {
					case lbefore > chop && lafter > chop:
						after = after[:len(after)-chop]
						before = before[chop:]
					case lbefore > chop && lafter < chop:
						before = before[chop-(chop-lafter):]
					case lafter > chop && lbefore < chop:
						after = after[:len(after)-(chop-lbefore)]
					}
				}
				out := before + term.Red + match + term.X + after
				choices = append(choices, grepChoice{
					hit: hit,
					str: fmt.Sprintf("%v%6v%v %v", term.Green, id, term.X, out),
				})
			}
			i, c, err := choose.From(choices)
			if err != nil {
				return err
			}
			if i > 0 {
				id := filepath.Base(filepath.Dir(c.hit.File))
				return editCmd.Call(x, id)
			}
			return nil
		}

		dex, err := ReadDex(keg.Path)
		if err != nil {
			return err
		}
		var lastid int
		for _, hit := range results.Hits {
			id, err := strconv.Atoi(filepath.Base(filepath.Dir(hit.File)))
			if err != nil {
				return err
			}
			if id == lastid {
				continue
			}
			lastid = id
			fmt.Println(dex.Lookup(id).AsInclude())
		}
		return nil
	},
}

//go:embed testdata/keg-dark.json
var dark []byte

//go:embed testdata/keg-notty.json
var notty []byte

var viewCmd = &Z.Cmd{
	Name:     `view`,
	Summary:  `view a specific node`,
	Usage:    `(help|ID|REGEXP)`,
	Params:   []string{`last`, `same`},
	MinArgs:  1,
	Commands: []*Z.Cmd{help.Cmd},

	Description: `
		The {{aka}} command renders a specific node for viewing in the
		terminal suitable for being cutting and pasting into other text
		documents and description fields. The argument passed may be an
		integer ID or a regular expression to be matched in the title text
		(as with {{cmd "edit"}} and {{cmd "title"}} commands. When matting
		a REGEXP case insensitive matching is assumed (prefix {{pre "(?i)"}}
		is added. (See {{cmd "grep"}} for how this default an be changed.)

		The {{aka}} command uses the https://github.com/charmbracelet/glamour
		package for rendering markdown directly to the terminal and
		therefore can be customized by setting the GLAMOUR_STYLE environment
		variable for those who wish. Since the popular GitHub command line
		utility uses this as well the same customization can be applied to
		both {{cmd "keg"}} and {{cmd "gh"}}.  By default, a variation on the
		{{pre "dark"}} style is used with line wrapping and margins disabled
		(for better cutting and pasting). To get a full copy of the style
		JSON used see the {{cmd "style"}} command.

		If the output is not to a terminal then the {{pre "notty"}} Glamour
		theme is used automatically.
		
	`,

	Call: func(x *Z.Cmd, args ...string) error {

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

				pre, err := x.Caller.Get(`regxpre`)
				if err != nil {
					return err
				}
				if pre == "" {
					pre = `(?i)`
				}

				re, err := regexp.Compile(pre + args[0])
				if err != nil {
					return err
				}

				choice := dex.ChooseWithTitleTextExp(re)
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

		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var r *glamour.TermRenderer
		if !term.IsInteractive() {
			r, err = glamour.NewTermRenderer(
				glamour.WithWordWrap(-1),
				glamour.WithStylesFromJSONBytes(notty),
			)
			if err != nil {
				return err
			}
			out, err := r.Render(string(buf))
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		}

		glamenv := os.Getenv(`GLAMOUR_STYLE`)
		if glamenv != "" {
			r, err = glamour.NewTermRenderer(
				glamour.WithEnvironmentConfig(),
				glamour.WithWordWrap(-1),
			)
			if err != nil {
				return err
			}
		} else {
			r, err = glamour.NewTermRenderer(
				glamour.WithStylesFromJSONBytes(dark),
				glamour.WithWordWrap(-1),
			)
			if err != nil {
				return err
			}
		}

		out, err := r.Render(string(buf))
		if err != nil {
			return err
		}
		Z.Page(out)

		return nil
	},
}
