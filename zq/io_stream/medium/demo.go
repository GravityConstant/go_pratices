package main

import (
	"fmt"
)

func main() {
	buf := make([]byte, 4)
	buf[0] = 'a'
	buf[2] = 'b'
	fmt.Println(string(buf))
}
