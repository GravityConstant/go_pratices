package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Data []BindPhone

type BindPhone struct {
	Id    string
	phone string
}

func init() {

}

func (d Data) Random() {
	rand.Seed(time.Now().Unix())
	strsLen := len(d)

	ns := make([]int, strsLen)
	for i := 0; i < strsLen; i++ {
		ns[i] = -1
	}
	// fmt.Println("init:", ns, strsLen)

	res := make(Data, strsLen)

	for i := 0; i < strsLen; i++ {
		n := rand.Intn(strsLen)
		n = GetPosition(n, strsLen, ns, true)
		// fmt.Println("real n: ", n)
		ns[i] = n
		res[i] = d[n]
	}
	copy(d, res)
}

func main() {
	s := "abc and "
	fmt.Println(strings.TrimSuffix(s, "and "))
}

// 还是要用结构体写，不然麻烦
func Random(strs []string) []string {
	rand.Seed(time.Now().Unix())
	strsLen := len(strs)

	ns := make([]int, strsLen)
	for i := 0; i < strsLen; i++ {
		ns[i] = -1
	}
	// fmt.Println("init:", ns, strsLen)

	res := make([]string, strsLen)

	for i := 0; i < strsLen; i++ {
		n := rand.Intn(strsLen)
		n = GetPosition(n, strsLen, ns, true)
		// fmt.Println("real n: ", n)
		ns[i] = n
		res[i] = strs[n]
	}
	return res
}

func GetPosition(e, l int, ns []int, up bool) (v int) {
	fmt.Println("in position:", e, ns)
	for _, n := range ns {
		if n == e {
			if e == 0 {
				e = GetPosition(e+1, l, ns, true)
			} else if e == l-1 {
				e = GetPosition(e-1, l, ns, false)
			} else {
				if up {
					e = GetPosition(e+1, l, ns, true)
				} else {
					e = GetPosition(e-1, l, ns, false)
				}
			}
		}
	}
	return e
}
