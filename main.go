package main

import (
	"automata/generator"
	"math"
	_ "net/http/pprof"

	"github.com/fatih/color"
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
			color.Red("\n=============================\nTesting pareto %d with %d states \n=============================\n", paretonum, statesAmt)
			transitions := generator.GenerateAllTransitions(paretonum, statesAmt)
			if length, transitions := generator.PowerSet(&transitions, int(math.Round(math.Max(math.Pow(2, float64(statesAmt)), math.Pow(2, float64(paretonum))))), statesAmt, paretonum); length != -1 {

				color.Yellow("Found paretonum %d with %d states and wlen %d\nTransitions: \n", paretonum, statesAmt, length)
				for _, trans := range transitions {
					color.Green("\t%v\n", trans)
				}
				break
			}
		}
	}
}
