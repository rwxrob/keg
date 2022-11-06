package node_test

import (
	"fmt"

	"github.com/rwxrob/keg/node"
)

func ExampleReadTitle() {
	title, _ := node.ReadTitle(`testdata/sample-node/README.md`)
	fmt.Println(title)
	// Output:
	// This is a title
}

func ExampleReadTitle_no_README() {
	title, _ := node.ReadTitle(`testdata/sample-node`)
	fmt.Println(title)
	// Output:
	// This is a title
}
