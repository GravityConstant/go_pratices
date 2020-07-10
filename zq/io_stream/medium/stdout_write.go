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

	for _, p := range proverbs {
		n, err := os.Stdout.Write([]byte(p))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if n != len(p) {
			fmt.Println("fail write to stdout")
			os.Exit(1)
		}
	}
}
