package main

import (
	"fmt"
	"io"
	"strings"
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

	j := 0
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			p[j] = char
			j++
		}
	}
	if j == 0 {
		return 0, io.EOF
	}
	return j, nil
}

func main() {
	reader := newAlphaReader(strings.NewReader("Hello! It's 9am, where is the sun?"))
	p := make([]byte, 4)
	for {
		n, err := reader.Read(p)
		if err == io.EOF {
			break
		}
		fmt.Print(string(p[:n]))
	}
	fmt.Println()
}
