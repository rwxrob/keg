package scan_test

import (
	"fmt"
	"strings"

	"github.com/rwxrob/keg/kegml/scan"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleTitle_short() {

	s := scanner.New(`# A short title`)
	fmt.Println(scan.Title(s, nil))
	s.Print()

	// Output:
	// true
	// 'e' 14-15 ""

}

func ExampleTitle_parsed_Short() {

	s := scanner.New(`# A short title`)
	title := new(strings.Builder)
	fmt.Println(scan.Title(s, title))
	s.Print()
	fmt.Printf("%q", title)

	// Output:
	// true
	// 'e' 14-15 ""
	// "A short title"

}

func ExampleTitle_long() {

	s := scanner.New(`# A really, really long title that is more than 72 runes long but doesn't get truncated`)
	fmt.Println(scan.Title(s, nil))
	s.Print()

	// Output:
	// false
	// 'g' 74-75 "et truncated"

}

func ExampleTitle_parsed_Long() {

	s := scanner.New(`# A really, really long title that is more than 72 runes long but doesn't get truncated`)
	title := new(strings.Builder)
	fmt.Println(scan.Title(s, title))
	s.Print()
	fmt.Printf("%q", title.String())

	// Output:
	// false
	// 'g' 74-75 "et truncated"
	// ""

}
