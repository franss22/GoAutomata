package main

import (
	"automata/generator"
	"fmt"
	"math"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/profile"
)

func main() {
	// Start profiling
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	start := time.Now()

	repeats := int64(1)

	for i := range repeats {
		fmt.Printf("%d, ", i)
		tempOut := os.Stdout
		// os.Stdout = nil
		GetMinStatesForEveryPareto(5, true)
		os.Stdout = tempOut

	}

	color.Red("\n Finished in %v, avg=%vms", time.Since(start), time.Since(start).Milliseconds()/repeats)

}

const MAX_STATES = 2

func GetMinStatesForEveryPareto(maxPareto int, parallel bool) {
	for paretonum := 2; paretonum < maxPareto; paretonum++ {
		for statesAmt := 2; MAX_STATES == -1 || statesAmt <= MAX_STATES; statesAmt++ {
			fmt.Print(color.RedString("\n=============================\nTesting pareto %d with %d states \n=============================\n", paretonum, statesAmt))
			transitions := generator.GenerateAllTransitions(paretonum, statesAmt)
			maxWlen := int(math.Round(math.Max(math.Pow(2, float64(statesAmt)), math.Pow(2, float64(paretonum)))))
			if length, transitions := generator.PowerSet(&transitions, maxWlen, statesAmt, paretonum, parallel); length != -1 {

				fmt.Print(color.YellowString("Found paretonum %d with %d states and wlen %d\nTransitions: \n", paretonum, statesAmt, length))
				for _, trans := range transitions {
					fmt.Print(color.GreenString("\t%v\n", trans))
				}
				break
			}
		}
	}
}
