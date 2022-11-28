# Q: What is a closure?

A closure is a way to *enclose* a reference to a variable or function within a function. Anything that returns the value of a reference that is not otherwise available is an *enclsure*

```go
package main

import "fmt"

var myage = 54

// Yes, this qualifies.

func Age1() int { return myage }

// This is more like what you would see in most examples.
// The returned function is a closure since it encloses the private, local
// myage variable.

func Age2() func() int {
	myage := 20
	someage := func() int { return myage }
	return someage
}

func main() {
	fmt.Println(Age1())
	fmt.Println(Age2()())
}
```
