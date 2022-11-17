package keg_test

import (
	"fmt"
	"path/filepath"

	"github.com/rwxrob/keg"
)

func ExampleNodePaths() {

	dirs, low, high := keg.NodePaths("testdata/samplekeg")

	fmt.Println(low)
	fmt.Println(high)

	for _, d := range dirs {
		fmt.Printf("%v ", filepath.Base(d.Path))
	}

	// Output:
	// 0
	// 12
	// 0 1 10 11 12 2 3 4 5 6 7 8 9

}

/*
func ExampleDex() {
	fmt.Println(keg.Dex("testdata/samplekeg"))

	// Output:
	// ignored
}
*/

/*
func ExampleMakeDex() {
	fmt.Println(keg.MakeDex("testdata/samplekeg"))

	// Output:
	// ignored
}
*/
