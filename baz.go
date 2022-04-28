// Copyright 2022 bonzai-example Authors
// SPDX-License-Identifier: Apache-2.0

package example

import (
	"log"

	Z "github.com/rwxrob/bonzai/z"
)

// exported leaf
var BazCmd = &Z.Cmd{
	Name: `baz`,
	Call: func(caller *Z.Cmd, none ...string) error {
		log.Print("Baz, suncreen song")
		return nil
	},
}
