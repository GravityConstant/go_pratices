package main

import (
	"math/rand"
	"testing"
	"time"
)

var numbers []int

func init() {
	rand.Seed(time.Now().UnixNano())
	numbers = generateList(1e3)
}

func generateList(totalNumbers int) []int {
	numbers := make([]int, totalNumbers)
	for i := 0; i < totalNumbers; i++ {
		numbers[i] = rand.Intn(totalNumbers)
	}
	return numbers
}

func Loop(nums []int, step int) {
	l := len(nums)
	for i := 0; i < step; i++ {
		for j := i; j < l; j += step {
			nums[j] = 4
		}
	}
}

func BenchmarkLoopStepOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Loop(numbers, 1)
	}
}

func BenchmarkLoopStepSixteen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Loop(numbers, 16)
	}
}
