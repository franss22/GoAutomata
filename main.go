package main

import (
	"automata/combinations"
	"automata/generator"
	"automata/timecea"
	"automata/types"
	"fmt"
	"math"
	_ "net/http/pprof"

	"github.com/pkg/profile"
	"github.com/schollz/progressbar/v3"
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

func onesToPretty(list []int) []byte {
	ret := make([]byte, len(list))
	chars := []byte{'-', '#'}
	for i, val := range list {
		ret[i] = chars[val]
	}
	return ret
}

func main() {
	// Start profiling
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	// go func() {
	// 	http.ListenAndServe(":8080", nil)
	// }()

	// fmt.Printf("%s (%v)\n", onesToPretty(nums.ToOnes().List), nums.List)
	// CalcCombinations(30)
	GetMinStatesForEveryPareto(20)

	// TestOneThing()

}

func GetMinStatesForEveryPareto(maxPareto int) {
	for paretonum := 2; paretonum < maxPareto; paretonum++ {
		for statesAmt := 2; ; statesAmt++ {
			fmt.Print("\nTesting pareto ", paretonum, " with ", statesAmt, " states\n")
			transitions := generator.GenerateAllTransitions(paretonum, statesAmt)
			if length, transitions := generator.PowerSet(&transitions, int(math.Round(math.Max(math.Pow(2, float64(statesAmt)), math.Pow(2, float64(paretonum))))), statesAmt, paretonum); length != -1 {

				fmt.Print("Found paretonum ", paretonum, " with ", statesAmt, " states and wlen ", length, "\n")
				fmt.Print("Transitions: \n")
				for _, trans := range transitions {
					fmt.Print(trans)
				}
				break
			}
		}
	}
}

func TestOneThing() {
	c := timecea.New(2)
	c.RegisterTransition(
		types.Transition{
			Input: types.TransitionInput{
				P: 0,
				L: 'a',
			},
			Output: types.TransitionOutput{
				Q: 1,
				Cond: types.ClockCondition{
					Cond1: -1,
					Cond2: -1,
				},
				Resets: types.ClockReset{
					Reset1: false,
					Reset2: false,
				},
			},
		},
	)
	c.RegisterTransition(
		types.Transition{
			Input: types.TransitionInput{
				P: 1,
				L: 'a',
			},
			Output: types.TransitionOutput{
				Q: 1,
				Cond: types.ClockCondition{
					Cond1: -1,
					Cond2: -1,
				},
				Resets: types.ClockReset{
					Reset1: false,
					Reset2: true,
				},
			},
		},
	)
	fmt.Print("Pareto?: ", c.TestAutomataForPareto(4, 2), "\n")
}

func CalcCombinations(size int) {
	nums := combinations.New(size)
	goal := int64(math.Pow(2, float64(size)))
	bar := progressbar.Default(goal)

	for n := range size + 1 {
		// fmt.Printf("Starting with %d values...", n)
		// start := time.Now()
		nums.Reset(n)
		for ok := true; ok; {
			bar.Add(1)

			ok = nums.Next()
		}
		// fmt.Printf("Done in %s\n", time.Since(start))
	}
}

// func main() {

// 	size := 5
// 	ones := ones.New(size)
// 	nums := combinations.New(size - 1)
// 	for n := range size + 1 {
// 		ones.Reset(n)
// 		nums.Reset(n)
// 		for ok := true; ok; {

// 			//ok = ones.MoveNext()
// 			ok = nums.IncrementNext()
// 			// fmt.Printf("Original Ones=%+v\n", ones.List)
// 			fmt.Printf("\t\t      combinations Ones=%+v\n", nums.List)

// 		}
// 	}
// }
