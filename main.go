package main

import (
	"automata/generator"
	"fmt"
	"math"
	_ "net/http/pprof"

	"github.com/pkg/profile"
)

func main() {
	// Start profiling
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	GetMinStatesForEveryPareto(20)

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
					fmt.Println(trans)
				}
				break
			}
		}
	}
}
