package keg

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

var nextCmd = &Z.Cmd{
	Name:        `next`,
	UseVars:     true,
	Summary:     help.S(_next),
	Description: help.D(_next),
	Commands:    []*Z.Cmd{help.Cmd, vars.Cmd},
	VarDefs:     Z.VarVals{`order`: `oldest`},

	Shortcuts: Z.ArgMap{
		`set`:   {`var`, `set`},
		`unset`: {`var`, `unset`},
		`get`:   {`var`, `get`},
	},

	Call: func(x *Z.Cmd, args ...string) error {

		keg, err := current(x.Caller)
		if err != nil {
			return err
		}

		cur, _ := x.Get(`node`)
		var entry *DexEntry

		order, err := x.Get(`order`)
		if err != nil {
			return err
		}

		if cur == `` {
			switch order {
			case `oldest`:
				entry = First(keg.Path)
			case `newest`:
				entry = Last(keg.Path)
			case `changed`:
				entry = LastChanged(keg.Path)
			}
			if entry == nil {
				return fmt.Errorf(_CantGetNextNode, nil)
			}
			cur = entry.ID()
			if err := x.Set(`node`, cur); err != nil {
				return err
			}
		} else {
			_, entry = Lookup(keg.Path, cur)
			if entry == nil {
				return fmt.Errorf(_CantGetNextNode, nil)
			}
			cur = entry.ID()
		}
		/*
		   // TODO
		   		var nxt *DexEntry
		   		switch order {
		   		case `oldest`:
		   			nxt = NextOldest(keg.Path, entry)
		   		case `newest`:
		   			nxt = NextNewest(keg.Path, entry)
		   		case `changed`:
		   			nxt = NextChanged(keg.Path, entry)
		   		}

		   		return Edit(keg.Path, nxt.N)
		*/
		return nil
	},
}
