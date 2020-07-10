package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, err := os.Open("./proverbs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if _, err := io.Copy(os.Stdout, file); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}