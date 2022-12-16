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
	"text/template"

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

// -------------------------------- get -------------------------------

func get(x *Z.Cmd, it string) (keg *Local, id string, entry *DexEntry, err error) {

	keg, err = current(x.Caller)
	if err != nil {
		return
	}

	switch it {

	case "same":
		if entry = LastChanged(keg.Path); entry != nil {
			id = entry.ID()
		}

	case "last":
		if entry = Last(keg.Path); entry != nil {
			id = entry.ID()
		}

	default:

		var dex *Dex
		dex, err = ReadDex(keg.Path)
		if err != nil {
			return
		}

		var idn int
		idn, err = strconv.Atoi(it)

		if err == nil {

			entry = dex.Lookup(idn)
			if entry == nil {
				err = fmt.Errorf(_NodeNotFound, idn)
				return
			}
			id = entry.ID()

		} else {

			var pre string
			pre, err = x.Caller.Get(`regxpre`)
			if err != nil {
				return
			}
			if pre == "" {
				pre = `(?i)`
			}

			var re *regexp.Regexp
			re, err = regexp.Compile(pre + it)
			if err != nil {
				return
			}

			entry = dex.ChooseWithTitleTextExp(re)
			if entry == nil {
				err = fmt.Errorf(_ChooseTitleFail)
				return
			}

			id = entry.ID()
		}
	}
	return
}

// ------------------------------ current -----------------------------

// has to stay here because needs vars package from x
func current(x *Z.Cmd) (*Local, error) {
	var name, dir string

	// if we have an env it beats config settings
	name = os.Getenv(`KEG_CURRENT`)
	if name != "" {

		switch name[0] {

		case os.PathSeparator:
			local := &Local{Path: name}
			if strings.HasSuffix(name, string(os.PathSeparator)+`docs`) {
				local.Name = filepath.Base(filepath.Dir(name))
				return local, nil
			}
			local.Name = filepath.Base(name)
			return local, nil

		case '~':
			local := &Local{Path: name}
			dir = fs.Tilde2Home(dir)
			if strings.HasSuffix(name, string(os.PathSeparator)+`docs`) {
				local.Name = filepath.Base(filepath.Dir(name))
				return local, nil
			}
			local.Name = filepath.Base(name)
			return local, nil

		default:
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

	}

	// check vars and conf
	name, _ = x.Get(`current`)
	if name != "" {

		if name[0] == os.PathSeparator || name[0] == '~' {
			os.Setenv(`KEG_CURRENT`, name)
			return current(x)
		}

		dir, _ = x.C(`map.` + name)
		if !(dir == "" || dir == "null") {
			dir = fs.Tilde2Home(dir)
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

	return nil, fmt.Errorf(_NoKegsFound)
}

// ------------------------------- Cmds -------------------------------

var Cmd = &Z.Cmd{
	Name:        `keg`,
	Aliases:     []string{`kn`},
	Version:     `v0.9.0`,
	UseVars:     true,
	Copyright:   `Copyright 2022 Robert S Muhlestein`,
	License:     `Apache-2.0`,
	Site:        `rwxrob.tv`,
	Source:      `git@github.com:rwxrob/keg.git`,
	Issues:      `https://github.com/rwxrob/keg/issues`,
	ConfVars:    true,
	Summary:     help.S(_keg),
	Description: help.D(_keg),

	Commands: []*Z.Cmd{
		editCmd, help.Cmd, conf.Cmd, vars.Cmd,
		indexCmd, createCmd, currentCmd, directoryCmd, deleteCmd,
		lastCmd, changesCmd, titlesCmd, initCmd, randomCmd,
		importCmd, grepCmd, viewCmd, columnsCmd, linkCmd, tagCmd,
	},

	Shortcuts: Z.ArgMap{
		`set`:    {`var`, `set`},
		`unset`:  {`var`, `unset`},
		`get`:    {`var`, `get`},
		`sample`: {`create`, `sample`},
	},
}

var currentCmd = &Z.Cmd{
	Name:        `current`,
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     help.S(_current),
	Description: help.D(_current),

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
	Name:        `titles`,
	Aliases:     []string{`title`},
	Usage:       `(help|REGEXP)`,
	UseVars:     true,
	Summary:     help.S(_titles),
	Description: help.D(_titles),
	Commands:    []*Z.Cmd{help.Cmd, vars.Cmd},

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

var directoryCmd = &Z.Cmd{
	Name:        `directory`,
	Aliases:     []string{`d`, `dir`},
	Usage:       `[help|REGEXP]`,
	MaxArgs:     1,
	Summary:     help.S(_directory),
	Description: help.D(_directory),
	Commands:    []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		if len(args) > 0 {
			dex, _ := ReadDex(keg.Path)
			re, err := regexp.Compile(args[0])
			if err != nil {
				return err
			}
			choice := dex.ChooseWithTitleTextExp(re)
			term.Print(filepath.Join(keg.Path, strconv.Itoa(choice.N)))
			return nil
		}

		term.Print(keg.Path)

		return nil
	},
}

var deleteCmd = &Z.Cmd{
	Name:        `delete`,
	Usage:       `(help|INTEGER_NODE_ID|last|same)`,
	Aliases:     []string{`del`, `rm`},
	Summary:     help.S(_delete),
	Description: help.D(_delete),
	MinArgs:     1,
	Commands:    []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, id, entry, err := get(x, args[0])
		if err != nil {
			return err
		}

		dir := filepath.Join(keg.Path, id)

		log.Println("❌", dir)

		if err := os.RemoveAll(dir); err != nil {
			return err
		}

		if err := DexRemove(keg.Path, entry); err != nil {
			return err
		}
		return Publish(keg.Path)

	},
}

var indexCmd = &Z.Cmd{
	Name:        `index`,
	Aliases:     []string{`dex`},
	Commands:    []*Z.Cmd{help.Cmd, dexUpdateCmd},
	Summary:     help.S(_index),
	Description: help.D(_index),
}

var dexUpdateCmd = &Z.Cmd{
	Name:        `update`,
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     help.S(_index_update),
	Description: help.D(_index_update),
	Call: func(x *Z.Cmd, args ...string) error {
		keg, err := current(x.Caller.Caller) // keg dex update
		if err != nil {
			return err
		}
		return MakeDex(keg.Path)
	},
}

var lastCmd = &Z.Cmd{
	Name:        `last`,
	Usage:       `[help|dir|id|title|time]`,
	Params:      []string{`dir`, `id`, `title`, `time`},
	MaxArgs:     1,
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     help.S(_last),
	Description: help.D(_last),

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
	Name:        `changes`,
	Aliases:     []string{`changed`},
	Usage:       `[help|COUNT|default|set default COUNT]`,
	UseVars:     true,
	Summary:     help.S(_changes),
	Description: help.D(_changes),
	Commands:    []*Z.Cmd{help.Cmd, vars.Cmd},

	Dynamic: template.FuncMap{
		`changesdef`: func() int { return ChangesDefault },
	},

	Shortcuts: Z.ArgMap{
		`default`: {`var`, `get`, `default`},
		`set`:     {`var`, `set`},
	},

	Call: func(x *Z.Cmd, args ...string) error {
		var err error
		var n int

		if len(args) > 0 {
			n, _ = strconv.Atoi(args[0])
		}

		if n <= 0 {
			def, err := x.Get(`default`)
			if err == nil && def != "" {
				n, err = strconv.Atoi(def)
				if err != nil {
					return err
				}
			}
		}

		if n <= 0 {
			n = ChangesDefault
		}

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		path := filepath.Join(keg.Path, `dex/changes.md`)
		if !fs.Exists(path) {
			return fmt.Errorf(_FileNotFound, `dex/changes.md`)
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

var initCmd = &Z.Cmd{
	Name:        `init`,
	Usage:       `[help]`,
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     help.S(_init),
	Description: help.D(_init),

	Call: func(_ *Z.Cmd, _ ...string) error {

		if fs.NotExists(`keg`) {
			if err := file.Overwrite(`keg`, _kegyaml); err != nil {
				return err
			}
		}

		if fs.NotExists(`0/README.md`) {
			if err := file.Overwrite(`0/README.md`, _zero_node); err != nil {
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
	Name:        `edit`,
	Aliases:     []string{`e`},
	Params:      []string{`last`, `same`},
	Usage:       `(help|ID|last|same|REGEX)`,
	Summary:     help.S(_edit),
	Description: help.D(_edit),
	Commands:    []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		if len(args) == 0 {
			return help.Cmd.Call(x, args...)
		}

		if !term.IsInteractive() {
			return titlesCmd.Call(x, args...)
		}

		keg, id, entry, err := get(x, args[0])
		if err != nil {
			return err
		}

		path := filepath.Join(keg.Path, id, `README.md`)
		if !fs.Exists(path) {
			return fmt.Errorf(_NodeNotFound, id)
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
		if atime.After(btime) {
			return Publish(keg.Path)
		}
		return nil

	},
}

var createCmd = &Z.Cmd{
	Name:        `create`,
	Aliases:     []string{`c`},
	Params:      []string{`sample`},
	MaxArgs:     1,
	Summary:     help.S(_create),
	Description: help.D(_create),
	Commands:    []*Z.Cmd{help.Cmd},

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
	Name:        `random`,
	Aliases:     []string{`rand`},
	Usage:       `[help|title|id|dir|edit]`,
	Params:      []string{`title`, `id`, `dir`, `edit`},
	MaxArgs:     1,
	Summary:     help.S(_random),
	Description: help.D(_random),
	Commands:    []*Z.Cmd{help.Cmd},

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
	Name:        `import`,
	Usage:       `[help|(DIR|NODEDIR)...]`,
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     help.S(_import),
	Description: help.D(_import),

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			d := dir.Abs()
			if d == "" {
				return fmt.Errorf(_AbsPathFail)
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

// columns first looks for term.WinSize.Col to have been set. If not
// found, the columns variable (from vars) is checked and used if found.
// Finally, the package global DefColumns will be used.
func columns(x *Z.Cmd) int {

	col := int(term.WinSize.Col) // only > 0 for interactive terminals
	if col > 0 {
		return col
	}

	colstr, err := x.Caller.Get(`columns`)
	if err == nil && colstr != "" {
		col, err = strconv.Atoi(colstr)
		if err == nil {
			return col
		}
	}

	return DefColumns

}

var columnsCmd = &Z.Cmd{
	Name:        `columns`,
	Usage:       `(help|col|cols)`,
	Summary:     help.S(_columns),
	Description: help.D(_columns),
	MaxArgs:     1,
	Commands:    []*Z.Cmd{help.Cmd},
	Dynamic:     template.FuncMap{`columns`: func() int { return DefColumns }},

	Call: func(x *Z.Cmd, args ...string) error {
		term.Print(columns(x))
		return nil
	},
}

type grepChoice struct {
	hit grep.Result
	str string
}

func (c grepChoice) String() string { return c.str }

var grepCmd = &Z.Cmd{
	Name:        `grep`,
	Usage:       `(help|REGEXP)`,
	MinArgs:     1,
	Commands:    []*Z.Cmd{help.Cmd},
	Summary:     help.S(_grep),
	Description: help.D(_grep),

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

		col := columns(x) - 14
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
	Name:        `view`,
	Usage:       `(help|ID|REGEXP)`,
	Summary:     help.S(_view),
	Description: help.D(_view),
	Params:      []string{`last`, `same`},
	MinArgs:     1,
	Commands:    []*Z.Cmd{help.Cmd},

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
					return fmt.Errorf(_ChooseTitleFail)
				}

				id = strconv.Itoa(choice.N)
			}
		}

		path := filepath.Join(keg.Path, id, `README.md`)

		if !fs.Exists(path) {
			return fmt.Errorf(_NodeNotFound, id)
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

var lastfmtExp = regexp.MustCompile(`(?:^|\n)linkfmt:\s*(.+)(?:\n|$)`)

var linkCmd = &Z.Cmd{
	Name:        `link`,
	Aliases:     []string{`url`},
	Summary:     help.S(_link),
	Description: help.D(_link),
	Commands:    []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, id, _, err := get(x, args[0])
		if err != nil {
			return err
		}

		kegfile := filepath.Join(keg.Path, `keg`)

		buf, err := os.ReadFile(kegfile)
		if err != nil {
			return err
		}

		f := lastfmtExp.FindStringSubmatch(string(buf))
		if f == nil {
			return fmt.Errorf(_NotInKegFile, `lastfmt`)
		}
		url := f[1]

		i := strings.Index(url, `{{id}}`)
		if i < 0 {
			return fmt.Errorf(_StringHasNo, `{{id}}`)
		}
		term.Print(url[:i] + id + url[i+6:])

		return nil
	},
}

var tagCmd = &Z.Cmd{
	Name:        `tag`,
	Aliases:     []string{`tags`},
	Params:      []string{`edit`},
	Usage:       `[help|edit|all|TAGS (NODEID|same|last|REGEXP)]`,
	Summary:     help.S(_tag),
	Description: help.D(_tag),
	Commands:    []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			args = append(args, `all`)
		}

		if len(args) == 1 {

			if args[0] == `edit` {
				return file.Edit(filepath.Join(keg.Path, `dex`, `tags`))
			}

			if args[0] == `list` {
				tags := Tags(keg.Path)
				if len(tags) > 0 {
					term.Print(tags)
				}
				return nil
			}

			if args[0] == `all` {
				str, err := os.ReadFile(filepath.Join(keg.Path, `dex`, `tags`))
				if err != nil {
					return err
				}
				fmt.Print(string(str))
				return nil
			}

			str, err := GrepTags(keg.Path, args[0])
			if err != nil {
				return err
			}
			fmt.Print(str)
			return nil
		}

		keg, id, _, err := get(x, args[1])

		return Tag(keg.Path, id, args[0])
	},
}
