package main

import (
	combinations "automata/combinations"
	"automata/ones"
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
func GenerateOnes(b *testing.B) {
	o := ones.New(size)
	for vals := range size + 1 {

		o.Reset(vals)
		for ok := true; ok; {
			ok = o.Next()
		}
	}
}

func BenchmarkNums(b *testing.B) {

	for i := 0; i < b.N; i++ {
		GenerateNums(b)
	}
}

func BenchmarkOnes(b *testing.B) {

	for i := 0; i < b.N; i++ {
		GenerateOnes(b)
	}
}
