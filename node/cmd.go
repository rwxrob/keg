package node

import (
	"log"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/fs/file"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

var Cmd = &Z.Cmd{
	Name:    `node`,
	Summary: `work with a single KEG node`,
	Commands: []*Z.Cmd{
		help.Cmd, vars.Cmd, conf.Cmd,
		createCmd,
	},
	Shortcuts: Z.ArgMap{},
}

var createCmd = &Z.Cmd{
	Name:     `create`,
	Aliases:  []string{`add`, `new`},
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

		readme, err := MkTemp()
		if err != nil {
			return err
		}

		if err := file.Edit(readme); err != nil {
			return err
		}

		title, err := ReadTitle(readme)
		if err != nil {
			return err
		}

		log.Print("title", title)

		return nil

		/*
			log.Println("TODO slugify the title gosimple/slug")
			log.Println("TODO call Import(path,current), fail if no current")
			return nil
		*/
	},
}
