package main

import (
	combinations "automata/combinations"
	"os"
	"syscall"
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

// func BenchmarkNums(b *testing.B) {

// 	for i := 0; i < b.N; i++ {
// 		GenerateNums(b)
// 	}
// }

func BenchmarkParetoCheckingParallel(b *testing.B) {
	defer func(stdout *os.File) {
		os.Stdout = stdout
	}(os.Stdout)
	os.Stdout = os.NewFile(uintptr(syscall.Stdin), os.DevNull)

	paretoNum := 4

	for i := 0; i < b.N; i++ {
		GetMinStatesForEveryPareto(paretoNum, true)
	}

}

func BenchmarkParetoCheckingSingle(b *testing.B) {
	defer func(stdout *os.File) {
		os.Stdout = stdout
	}(os.Stdout)
	os.Stdout = os.NewFile(uintptr(syscall.Stdin), os.DevNull)

	paretoNum := 4

	for i := 0; i < b.N; i++ {
		GetMinStatesForEveryPareto(paretoNum, false)
	}

}
