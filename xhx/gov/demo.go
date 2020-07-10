package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	matches, err := filepath.Glob("views/*/*")
	if err != nil {
		fmt.Println("Glob file error. ", err)
		return
	}

	fmt.Printf("%#v\n", matches)
}
