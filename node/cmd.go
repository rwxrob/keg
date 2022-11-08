package node

import (
	"path"

	"github.com/gosimple/slug"
	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/fs"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/help"
	"github.com/rwxrob/keg/node/parse"
	"github.com/rwxrob/keg/node/read"
	"github.com/rwxrob/vars"
)

var Cmd = &Z.Cmd{
	Name:    `node`,
	Aliases: []string{`n`},
	Summary: `work with a single KEG node`,
	Commands: []*Z.Cmd{
		help.Cmd, vars.Cmd, conf.Cmd,
		createCmd, parse.Cmd,
	},
	Shortcuts: Z.ArgMap{},
}

var createCmd = &Z.Cmd{
	Name:     `create`,
	Aliases:  []string{`c`},
	Summary:  `create node in current KEG`,
	Commands: []*Z.Cmd{help.Cmd},

	Call: func(x *Z.Cmd, args ...string) error {

		curkeg, err := x.Caller.Caller.Get(`current`)
		if err != nil {
			return err
		}

		if curkeg == "" {
			return Z.MissingVar{x.Caller.Caller.Path(`current`)}
		}

		curdir, err := x.Caller.Caller.C(`map.` + curkeg)
		if curdir == "" {
			return Z.MissingConf{x.Caller.Caller.Path(`map.` + curkeg)}
		}

		curdir = fs.Tilde2Home(curdir)

		readme, err := MkTemp()
		if err != nil {
			return err
		}

		if err := file.Edit(readme); err != nil {
			return err
		}

		title, err := read.Title(readme)
		if err != nil {
			return err
		}

		id := slug.Make(title)
		if err := Import(path.Dir(readme), curdir, id); err != nil {
			return err
		}
		// TODO update the index, but only for the new node
		return nil
	},
}
