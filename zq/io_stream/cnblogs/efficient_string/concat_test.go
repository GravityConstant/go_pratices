package efficient_string

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

var global string

// nextString is an iterator we use to represent a process
// that returns strings that we want to concatenate in order.
func nextString() func() string {
	n := 0
	// closure captures variable n
	return func() string {
		n += 1
		return strconv.Itoa(n)
	}
}

func BenchmarkBufferString(b *testing.B) {
	var ns string
	var buffer bytes.Buffer
	for i := 0; i < b.N; i++ {
		next := nextString()
		for u := 0; u < 100; u++ {
			buffer.WriteString(next())
		}
	}
	ns = buffer.String()
	global = ns
}

func BenchmarkPlus(b *testing.B) {
	var ns string
	for i := 0; i < b.N; i++ {
		next := nextString()
		for u := 0; u < 100; u++ {
			ns += next()
		}
	}
	global = ns
}

func BenchmarkStringJoin(b *testing.B) {
	var ns string
	s := []string{}
	for i := 0; i < b.N; i++ {
		next := nextString()
		for u := 0; u < 100; u++ {
			s = append(s, next())
		}
	}
	ns = strings.Join(s, "")
	global = ns
}
