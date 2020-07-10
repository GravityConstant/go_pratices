package main

import (
	"fmt"

	"xhx/private/b"
)

func main() {
	s := b.Greet("world")
	fmt.Println(s)
	s = b.Hi("world")
	fmt.Println(s)
}
