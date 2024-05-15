package main

import (
	"automata/combinations"
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
	CalcCombinations(30)

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
