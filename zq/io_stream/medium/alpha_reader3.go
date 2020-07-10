package main

import (
	"fmt"
	"io"
	"os"
)

type alphaReader struct {
	reader io.Reader
}

func newAlphaReader(r io.Reader) *alphaReader {
	return &alphaReader{reader: r}
}

func alpha(r byte) byte {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

func (a *alphaReader) Read(p []byte) (n int, err error) {
	n, err = a.reader.Read(p)
	if err != nil {
		return n, err
	}

	buf := make([]byte, n)
	for i := 0; i < len(p); i++ {
		if char := alpha(p[i]); char != 0 {
			buf[i] = char
		}
	}
	copy(p, buf)
	return n, nil
}

func main() {
	file, err := os.Open("./alpha_reader2.go")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	p := make([]byte, 4)
	reader := newAlphaReader(file)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}
