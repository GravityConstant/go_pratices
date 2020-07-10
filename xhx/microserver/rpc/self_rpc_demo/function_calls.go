package main

func TopLevel(x int) {}

type Pointer struct{}

func (*Pointer) M(int) {}

type Value struct{}

func (Value) M(int) {}

type Interface interface{ M(int) }

var literal = func(x int) {}

func main() {

	// direct call of top-level func

	TopLevel(1)

	// direct call of method with value receiver (two spellings, but same)

	var v Value

	v.M(1)

	Value.M(v, 1)

	// direct call of method with pointer receiver (two spellings, but same)

	var p Pointer

	(&p).M(1)

	(*Pointer).M(&p, 1)

	// indirect call of func value (×4)

	f1 := TopLevel

	f1(1)

	f2 := Value.M

	f2(v, 1)

	f3 := (*Pointer).M

	f3(&p, 1)

	f4 := literal

	f4(1)

	// indirect call of method on interface (×3)

	var i Interface

	i = v

	i.M(1)

	i = &v

	i.M(1)

	i = &p

	i.M(1)

	Interface.M(i, 1)

	Interface.M(v, 1)

	Interface.M(&p, 1)

}
