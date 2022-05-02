package mark_test

import (
	"fmt"

	"github.com/rwxrob/keg/mark"
)

func ExampleParser_EOB() {
	p := new(mark.Parser)
	p.S.B = []byte("\n\n")
	fmt.Println(p.EOB())
	// Output:
	// true
}

func ExampleParser_UPrint() {
	p := new(mark.Parser)
	p.S.B = []byte(`foo`)
	fmt.Println(p.UPrint())
	// Output:
	// true
}

func ExampleParser_UPrint_false() {
	p := new(mark.Parser)
	p.S.B = []byte("\x00")
	fmt.Println(p.UPrint())
	// Output:
	// false
}

func ExampleParser_Text() {
	p := new(mark.Parser)
	p.S.B = []byte(`foo`)
	fmt.Println(p.Text())
	// Output:
	// true
}

func ExampleParser_Text_false() {
	p := new(mark.Parser)
	p.S.B = []byte(``)
	fmt.Println(p.Text())
	// Output:
	// false
}

func ExampleParser_Emoji() {
	p := new(mark.Parser)
	p.S.B = []byte(`💢`)
	p.Emoji()
	p.S.Print()
	// Output:
	// 4 '💢' ""
}

func ExampleParser_ULetter() {
	p := new(mark.Parser)
	p.S.B = []byte(`A`)
	fmt.Println(p.ULetter())
	// Output:
	// true
}

func ExampleParser_ULetter_false() {
	p := new(mark.Parser)
	p.S.B = []byte(` `)
	fmt.Println(p.ULetter())
	// Output:
	// false
}

func ExampleParser_Tag() {
	p := new(mark.Parser)
	p.S.B = []byte(`#foo`)
	fmt.Println(p.Tag())
	// Output:
	// true
}

func ExampleParser_Mono() {
	p := new(mark.Parser)
	p.S.B = []byte("`some`")
	fmt.Println(p.Mono())
	// Output:
	// true
}

func ExampleParser_Mono_false() {
	p := new(mark.Parser)
	p.S.B = []byte("``")
	fmt.Println(p.Mono())
	// Output:
	// false
}
