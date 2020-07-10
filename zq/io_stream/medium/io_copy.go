package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	proverbs := new(bytes.Buffer)
	proverbs.Write([]byte("Channels orchestrate muxtexes serialize\n"))
	proverbs.Write([]byte("Cgo is not go\n"))
	proverbs.Write([]byte("Errors are values\n"))
	proverbs.Write([]byte("Don't panic\n"))

	file, err := os.Create("./proverbs")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := io.Copy(file, proverbs); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("file created")
}
