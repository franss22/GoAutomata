package generator

import (
	"automata/combinations"
	"automata/timecea"
	"automata/types"
	"context"
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
	maxConds := 0 //max(statesAmt*2, paretoNum*2)
	resultSize := (maxConds + 1) * (maxConds + 1) * 4 * 1 * statesAmt * statesAmt
	result := make([]types.Transition, 0, resultSize)
	// color.Yellow("cap=%d, len=%d, resultSize=%d", cap(result), len(result), resultSize)

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

type CheckResult struct {
	ShouldReturn bool
	Length       int
	TestCase     []types.Transition
	Iter         int
}

type ParetoParams struct {
	trAmt       int
	transitions *[]types.Transition
	maxWLen     int
	statesAmt   int
	paretoNum   int
}

const THREAD_N = 16
const MAX_TRANSITIONS = 10

func PowerSet(transitions *[]types.Transition,
	maxWLen int,
	statesAmt int,
	paretoNum int,
	parallel bool,
) (int, []types.Transition) {
	size := len(*transitions)
	n := combinations.New(size)

	for trAmt := statesAmt; statesAmt <= min(size, MAX_TRANSITIONS); trAmt++ {

		pp := ParetoParams{
			trAmt:       trAmt,
			transitions: transitions,
			maxWLen:     maxWLen,
			statesAmt:   statesAmt,
			paretoNum:   paretoNum,
		}
		iters := int64(combin.Binomial(size, trAmt))

		n.Reset(trAmt)
		fmt.Print(color.BlueString("\nTesting with combinations of %d transitions (%d combinations per thread)\n", trAmt, int(iters/int64(THREAD_N)+1)))

		// fmt.Printf("iters: %d, t_iters: %d, t_iters sum: %d\n", iters, int(iters/int64(THREAD_N)+1), int(iters/int64(THREAD_N)+1)*THREAD_N)
		var shouldReturn bool
		var length int
		var testCase []types.Transition
		if parallel {
			shouldReturn, length, testCase = ParallelCheckCombinationsWithNTransitions(THREAD_N, n, int(iters/int64(THREAD_N)+1), pp)
		} else {
			shouldReturn, length, testCase = CheckCombinationsWithNTransitions(pp, &n, progressbar.Default(iters))
		}

		if shouldReturn {
			return length, testCase
		}
	}
	return -1, []types.Transition{}
}

func ParallelCheckCombinationsWithNTransitions(ROUTINE_N int, n combinations.Nums, routIters int, pp ParetoParams) (bool, int, []types.Transition) {
	resultsChannel := make(chan CheckResult, ROUTINE_N)
	ctx, ctxCancel := context.WithCancel(context.Background())

	fmt.Print(color.CyanString("Starting [%d] routines...", ROUTINE_N))
	go ParallelStart(ctx, ROUTINE_N, n, routIters, pp, resultsChannel)

	fmt.Printf("Waiting for (%d) routines...\n ", ROUTINE_N)

	for i := range ROUTINE_N {
		fmt.Printf("(%d),", i)
		res := <-resultsChannel

		if res.ShouldReturn {
			fmt.Print(color.GreenString("Found match in combination %d\n", res.Iter))
			ctxCancel()
			return true, res.Length, res.TestCase
		}

	}
	ctxCancel()
	fmt.Print(color.MagentaString("Nothing Found\n"))
	return false, -1, nil
}

func ParallelStart(ctx context.Context, ROUTINE_N int, n combinations.Nums, routIters int, pp ParetoParams, resultsChannel chan CheckResult) {
	for i := range ROUTINE_N {
		select {
		case <-ctx.Done():
			return
		default:
			pnums := n.NewPnums(routIters)
			go ParallelChecking(pp, pnums, int(routIters), resultsChannel, ctx)
			fmt.Print(color.CyanString(" [%d],", i))
			n.Advance(routIters)
		}
	}
}

func ParallelChecking(pp ParetoParams, pNums combinations.PNums, threadIters int, results chan CheckResult, ctx context.Context) {
	for ok := true; ok; {
		select {
		case <-ctx.Done():
			return
		default:
			testCase := make([]types.Transition, pp.trAmt)
			for i, index := range pNums.Indexes() {
				testCase[i] = (*pp.transitions)[index]
			}
			if found, length := PowerSetFound(&testCase, pp.maxWLen, pp.statesAmt, pp.paretoNum); found {
				// fmt.Printf("FOUND IN ITER %d\n", pNums.CurrentIteration())
				results <- CheckResult{true, length, testCase, pNums.Nums.ItersDone}
				return
			}
			ok = pNums.Next()
		}
	}
	results <- CheckResult{false, -1, nil, -1}

	//mandar shouldReturn, Length, testCase al proceso padre
}

func CheckCombinationsWithNTransitions(pp ParetoParams, n combinations.Combinator, p *progressbar.ProgressBar) (bool, int, []types.Transition) {
	for ok := true; ok; {
		testCase := make([]types.Transition, pp.trAmt)
		for i, index := range n.Indexes() {
			testCase[i] = (*pp.transitions)[index]
		}
		if found, length := PowerSetFound(&testCase, pp.maxWLen, pp.statesAmt, pp.paretoNum); found {
			fmt.Printf("FOUND IN ITER %d\n", n.CurrentIteration())
			return true, length, testCase
		}
		ok = n.Next()
		p.Add(1)
	}
	return false, -1, nil
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
