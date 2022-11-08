package kegml_test

import (
	"fmt"

	"github.com/rwxrob/keg/kegml"
	"github.com/rwxrob/pegn/scanner"
)

func ExampleTitle_short() {

	s := scanner.New(`# A short title`)
	fmt.Println(kegml.Title.Scan(s, nil))
	s.Print()

	// Output:
	// true
	// 'e' 14-15 ""

}

func ExampleTitle_parsed_Short() {

	s := scanner.New(`# A short title`)
	title := make([]rune, 0, 100)
	fmt.Println(kegml.Title.Scan(s, &title))
	s.Print()
	fmt.Printf("%q", string(title))

	// Output:
	// true
	// 'e' 14-15 ""
	// "A short title"

}

func ExampleTitle_long() {

	s := scanner.New(`# A really, really long title that is more than 72 runes long but doesn't get truncated`)
	fmt.Println(kegml.Title.Scan(s, nil))
	s.Print()

	// Output:
	// false
	// 't' 72-73 " get truncated"

}

func ExampleTitle_parsed_Long() {

	s := scanner.New(`# A really, really long title that is more than 72 runes long but doesn't get truncated`)
	title := make([]rune, 0, 70)
	fmt.Println(kegml.Title.Scan(s, &title))
	s.Print()
	fmt.Printf("%q", string(title))

	// Output:
	// false
	// 't' 72-73 " get truncated"
	// "A really, really long title that is more than 72 runes long but doesn'"

}
