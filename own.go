// Copyright 2022 bonzai-example Authors
// SPDX-License-Identifier: Apache-2.0

package example

import (
	"log"

	Z "github.com/rwxrob/bonzai/z"
)

// private leaf
var ownCmd = &Z.Cmd{
	Name: `own`,
	Call: func(caller *Z.Cmd, none ...string) error {
		log.Print("I'm in my own file.")
		return nil
	},
}
