package main

import (
	"automata/ones"
	"fmt"
)

// func main() {
// 	paretoNum, statesAmt := 5, 2

// 	transitions := generator.GenerateAllTransitions(paretoNum, statesAmt)
// 	curr := make([]types.Transition, 0, len(transitions)*2)
// 	color.Green("Transitions length: %d\n", len(transitions))
// 	ps := generator.PowerSet(&transitions, -1, &curr, 4, paretoNum, statesAmt)
// 	color.Green("Power Set: %d\n", ps)
// 	fmt.Printf("curr=%v", len(curr))

// }

func main() {
	size := 6
	ones := ones.New(size)
	for n := range size + 1 {
		ones.Reset(n)
		for ok := true; ok; {
			/*

				Run you code

				Here

			*/

			ok = ones.MoveNext()
			fmt.Printf("Ones=%+v\n", ones)
		}
	}
}
