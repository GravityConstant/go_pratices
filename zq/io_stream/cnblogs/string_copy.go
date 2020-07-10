package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	sr := strings.NewReader("Hello")
	sb := strings.Builder{}
	io.Copy(&sb, sr)
	io.WriteString(&sb, ", 世界")
	fmt.Println(sr, &sb)
}
