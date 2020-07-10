package main

import (
	"fmt"
	"reflect"
)

func main() {
	makeFunctionDemo()
	makeFunctionPratice()
}

func makeFunctionDemo() {
	// swap is the implementation passed to MakeFunc.
	// It must work in terms of reflect.Values so that it is possible
	// to write code without knowing beforehand what the types
	// will be.
	swap := func(in []reflect.Value) []reflect.Value {
		return []reflect.Value{in[1], in[0]}
	}

	// makeSwap expects fptr to be a pointer to a nil function.
	// It sets that pointer to a new function created with MakeFunc.
	// When the function is invoked, reflect turns the arguments
	// into Values, calls swap, and then turns swap's result slice
	// into the values returned by the new function.
	makeSwap := func(fptr interface{}) {
		// fptr is a pointer to a function.
		// Obtain the function value itself (likely nil) as a reflect.Value
		// so that we can query its type and then set the value.
		fn := reflect.ValueOf(fptr).Elem()

		// Make a function of the right type.
		v := reflect.MakeFunc(fn.Type(), swap)

		// Assign it to the value fn represents.
		fn.Set(v)
	}

	// Make and call a swap function for ints.
	var intSwap func(int, int) (int, int)
	makeSwap(&intSwap)
	fmt.Println(intSwap(0, 1))

	// Make and call a swap function for float64s.
	var floatSwap func(float64, float64) (float64, float64)
	makeSwap(&floatSwap)
	fmt.Println(floatSwap(2.72, 3.14))
}

func makeFunctionPratice() {
	add := func(args []reflect.Value) []reflect.Value {
		var resi int64
		for _, arg := range args {
			switch arg.Kind() {
			case reflect.Int, reflect.Int64:
				resi += arg.Int()
			}
		}
		return []reflect.Value{reflect.ValueOf(int(resi))}
	}
	makeAdd := func(fptr interface{}) {
		fn := reflect.ValueOf(fptr).Elem()
		v := reflect.MakeFunc(fn.Type(), add)
		fn.Set(v)
	}
	var intAdd func(i, j, m, n int) int
	makeAdd(&intAdd)
	fmt.Println(intAdd(1, 2, 3, 4))
}
