package generator

import (
	"automata/combinations"
	"automata/timecea"
	"automata/types"
	"fmt"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"gonum.org/v1/gonum/stat/combin"
)

var resetOptions = [...]types.ClockReset{
	{Reset1: false, Reset2: false},
	{Reset1: false, Reset2: true},
	{Reset1: true, Reset2: false},
	{Reset1: true, Reset2: true},
}

func GenerateAllTransitions(paretoNum int, statesAmt int) []types.Transition {
	maxConds := max(statesAmt*2, paretoNum*2)
	resultSize := (maxConds + 1) * (maxConds + 1) * 4 * 1 * statesAmt * statesAmt
	result := make([]types.Transition, 0, resultSize)
	color.Yellow("cap=%d, len=%d, resultSize=%d", cap(result), len(result), resultSize)

	for c1 := range types.Time(maxConds + 1) {
		cond1 := c1
		if c1 == 0 {
			cond1 = types.NO_COND
		}
		for c2 := range c1 + 1 {
			cond2 := c2
			if c1 == 0 {
				cond2 = types.NO_COND
			}
			for _, clockreset := range resetOptions {
				var l types.Letter = 'a'
				for stateInit := range types.State(statesAmt) {
					for stateFin := stateInit; stateFin < types.State(statesAmt); stateFin++ {
						result = append(result, types.Transition{
							Input: types.TransitionInput{
								P: stateInit,
								L: l,
							},
							Output: types.TransitionOutput{
								Q: stateFin,
								Cond: types.ClockCondition{
									Cond1: cond1,
									Cond2: cond2},
								Resets: clockreset,
							},
						})
					}
				}

			}
		}
	}

	return result
}

// func PowerSetOld(
// 	transitions *[]types.Transition,
// 	index int,
// 	curr *[]types.Transition,
// 	maxWLen int,
// 	statesAmt int,
// 	paretoNum int,
// ) int {
// 	n := len(*transitions)
// 	if index == n-1 {
// 		color.Cyan("Got to end of recursion, n=%d", n)
// 		tc := timecea.New(statesAmt)
// 		for _, tr := range *transitions {
// 			tc.RegisterTransition(tr)
// 		}
// 		length := tc.TestAutomataForPareto(maxWLen, paretoNum)
// 		return length
// 	}
// 	for i := index + 1; i < n; i++ {
// 		*curr = append(*curr, (*transitions)[i])
// 		// color.Magenta("For, with i=%d, curr=%v", i, curr)

// 		length := PowerSetOld(transitions, i, curr, maxWLen, statesAmt, paretoNum)
// 		if length == -1 {
// 			color.Blue("Case length == -1")
// 			*curr = (*curr)[:len(*curr)-1]
// 		} else {
// 			// color.Red("Return Length")
// 			return length
// 		}
// 		color.Red("Final del loop")

// 	}
// 	return -1
// }

func PowerSetFound(
	curr *[]types.Transition,
	maxWLen int,
	statesAmt int,
	paretoNum int,
) (bool, int) {

	tc := timecea.New(statesAmt)
	for _, tr := range *curr {
		tc.RegisterTransition(tr)
	}
	length := tc.TestAutomataForPareto(maxWLen, paretoNum)

	return (length != -1), length
}

func PowerSet(transitions *[]types.Transition,
	maxWLen int,
	statesAmt int,
	paretoNum int,
) (int, []types.Transition) {
	size := len(*transitions)
	n := combinations.New(size)
	for trAmt := statesAmt; statesAmt <= size; trAmt++ {
		fmt.Print("Testing with combinations of ", trAmt, " transitions\n")
		p := progressbar.Default(int64(combin.Binomial(size, trAmt)))
		n.Reset(trAmt)
		for ok := true; ok; {
			testCase := make([]types.Transition, trAmt)
			for i, index := range n.List {
				testCase[i] = (*transitions)[index]
			}
			if found, length := PowerSetFound(&testCase, maxWLen, statesAmt, paretoNum); found {
				return length, testCase
			}
			ok = n.Next()
			p.Add(1)
		}
	}
	return -1, []types.Transition{}
}

// func PowerSet(
// 	transitions *[]types.Transition,
// 	maxWLen int,
// 	statesAmt int,
// 	paretoNum int,
// ) int {
// 	n := len(*transitions)
// 	bitarray := bitarray(n)
// 	for bitarray.isNotAllTrue() {
// 		testTr := make([]types.Transition, 0, n)
// 		for i := range n {

// 		}
// 	}

// }

// type Counter struct {
// 	Count []int
// 	Max int

// }

// func newCounter(cap int) Counter {
// 	return Counter{Count: make([]int, 0, cap), Max:cap}
// }

// func (c * Counter) AddDigit() {
// 	c.Count = append(c.Count, 0)
// }

// func (c * Counter) IncrementDigit(index int) bool{

// }
// func (c *Counter) Increment() {
// 	if (len(c.Count)==0) {
// 		c.AddDigit()
// 	}
// 	i := len(c.Count)-1
// 	if digit := c.Count[i]; digit == 0 {
// 		panic("asassaas")
// 	}

// }
