package main

import (
	"bytes"
	"io"
	"os"
)

func main() {
	proverbs := new(bytes.Buffer)
	proverbs.Write([]byte("Channels orchestrate muxtexes serialize\n"))
	proverbs.Write([]byte("Cgo is not go\n"))
	proverbs.Write([]byte("Errors are values\n"))
	proverbs.Write([]byte("Don't panic\n"))

	r, w := io.Pipe()
	defer r.Close()

	go func() {
		defer w.Close()
		io.Copy(w, proverbs)
	}()

	io.Copy(os.Stdout, r)
}
