package main

import (
	"fmt"
	"os"
)

func main() {
	proverbs := []string{
		"Channels orchestrate muxtexes serialize\n",
		"Cgo is not go\n",
		"Errors are values\n",
		"Don't panic\n",
	}

	file, err := os.Create("./proverbs.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, p := range proverbs {
		n, err := file.Write([]byte(p))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n != len(p) {
			fmt.Println("failed to write data")
			os.Exit(1)
		}
	}
	fmt.Println("file write done")
}
