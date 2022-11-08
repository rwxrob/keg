package read_test

import (
	"fmt"

	"github.com/rwxrob/keg/node/read"
)

func ExampleTitle() {
	title, _ := read.Title(`testdata/sample-node/README.md`)
	fmt.Println(title)
	// Output:
	// This is a title
}

func ExampleTitle_no_README() {
	title, _ := read.Title(`testdata/sample-node`)
	fmt.Println(title)
	// Output:
	// This is a title
}
