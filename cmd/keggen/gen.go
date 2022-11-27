/*
This tool is to help create the keg binary distributions. It depends only on having the go binary and the GitHub command line tool (gh) installed. Run it from within same directory that contains go.mod.
*/
package main

import (
	"log"
	"os"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/fs"
)

func main() {
	if fs.NotExists(`go.mod`) {
		log.Println(`run from directory containing go.mod`)
		os.Exit(1)
	}
	Z.Exec(`ls`)
}
