package main

import (
	combinations "automata/combinations"
	"testing"
)

const size = 50

func GenerateNums(b *testing.B) {
	n := combinations.New(size)
	for vals := range size + 1 {

		n.Reset(vals)
		for ok := true; ok; {
			ok = n.Next()
		}
	}
}

func BenchmarkNums(b *testing.B) {

	for i := 0; i < b.N; i++ {
		GenerateNums(b)
	}
}
