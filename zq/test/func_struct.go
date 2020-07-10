package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"time"
	// "github.com/jinzhu/now"
)

type Adder interface {
	add(params ...int) int
}

type AdderFunc func(params ...int) int

func (f AdderFunc) add(params ...int) int {
	params = append(params, 100)
	return f(params...)
}

var IntAdder Adder = AdderFunc(intAdd)
var IntAdder2 Adder = AdderFunc(intAdd2)

func intAdd(params ...int) int {
	var sum int

	for _, val := range params {
		sum += val
	}

	return sum
}

func intAdd2(params ...int) int {
	var sum int

	for i := range params {
		sum += i
	}

	return sum
}

func main() {
	arr := []int{1, 2, 3, 4, 5}
	fmt.Println(IntAdder.add(arr...))
	fmt.Println(IntAdder2.add(arr...))

	cur := time.Now()

	curWeekday := cur.Weekday()
	fmt.Println(curWeekday, curWeekday+1 > 6)

}
