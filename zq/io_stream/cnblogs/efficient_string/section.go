package main

import (
	"bytes"
	"fmt"
	"unsafe"
)

const RECORD_LEN = 23

type User struct {
	name [20]byte
	age  int16
	sex  bool
}

func main() {
	man := User{age: 21, sex: true}
	name := man.name[:]
	copy(name, "Jack")
	fmt.Println(man)

	rw := &bytes.Buffer{}
	pBytes := (*[RECORD_LEN]byte)(unsafe.Pointer(&man))
	rw.Write(pBytes[:])
	fmt.Printf("% x\n", rw.Bytes())

	pUser := (*User)(unsafe.Pointer(&rw.Bytes()[0]))
	orgin := *pUser //copy!
	fmt.Println(orgin)

	fmt.Println("Hello, ", string(name))
}
