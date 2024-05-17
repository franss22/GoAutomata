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

	THREAD_N := 8

	for trAmt := statesAmt; statesAmt <= size; trAmt++ {
		fmt.Print("Testing with combinations of ", trAmt, " transitions\n")
		iters := int64(combin.Binomial(size, trAmt))
		threadIters := int(iters/int64(THREAD_N) + 1)

		p := progressbar.Default(iters)

		n.Reset(trAmt)
		quit := make(chan bool)
		results := make(chan (bool, int, []types.Transition))//st5ruct de las 3 cosas
		//N veces
		for range THREAD_N {
			go ParallelChecking(n, trAmt, transitions, maxWLen, statesAmt, paretoNum, int(threadIters), quit)
			n.Advance(threadIters)
		}
		//esperar los resultados
		if shouldReturn {
			return Length, testCase
		}
	}
	return -1, []types.Transition{}
}

func ParallelChecking(n combinations.Nums, trAmt int, transitions *[]types.Transition, maxWLen int, statesAmt int, paretoNum int, threadIters int, quit chan bool) (bool, int, []types.Transition) {
	pNums := n.NewPnums(int(threadIters))
	for ok := true; ok; {
		select {
        case <- quit:
            return false, 0, nil
        default:
            testCase := make([]types.Transition, trAmt)
		for i, index := range pNums.Indexes() {
			testCase[i] = (*transitions)[index]
		}
		if found, length := PowerSetFound(&testCase, maxWLen, statesAmt, paretoNum); found {
			return true, length, testCase
		}
		ok = n.Next()
        }
	}
	

	shouldReturn, Length, testCase := CheckCombinationsWithNTransitions(trAmt, &pNums, transitions, maxWLen, statesAmt, paretoNum, quit)

	//mandar shouldReturn, Length, testCase al proceso padre
}

func CheckCombinationsWithNTransitions(trAmt int, n combinations.Combinator, transitions *[]types.Transition, maxWLen int, statesAmt int, paretoNum int) (bool, int, []types.Transition) {
	for ok := true; ok; {
		testCase := make([]types.Transition, trAmt)
		for i, index := range n.Indexes() {
			testCase[i] = (*transitions)[index]
		}
		if found, length := PowerSetFound(&testCase, maxWLen, statesAmt, paretoNum); found {
			return true, length, testCase
		}
		ok = n.Next()
		// p.Add(1)
	}
	return false, 0, nil
}

/*
n Nums [1, 2, 3]
..... iter N
n Nums [8, 9, 10]


N /8
nx8
n0 -> N/8*0 go routine func copy(n)
n avanza N/8
n1 -> N/8*1  go routine func copy(n)
n avanza N/8
n2 -> N/8*2
n1 -> N/8*1
n1 -> N/8*1
n1 -> N/8*1




Sea:
n el iterador de combinaciones, con T cantidad de transiciones y K indices, con L combinaciones posibles
n(i) el iterador en la i-esima combinación
n(0) -> 1^K 0^(T-K)
//n(0)[T=5, K=3] 11100 / [0, 1, 2]


*/
