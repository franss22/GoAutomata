package timecea

import (
	"automata/types"
)

type TransitionMap map[types.TransitionInput][]types.TransitionOutput

type ClockMap map[types.State][]types.Clock

func (cm ClockMap) AddClock(q types.State, clock types.Clock) {
	cm[q] = append((cm)[q], clock)
}

type TimeCEA struct {
	StatesAmt   int
	ClocksAmt   int
	Q0          types.State
	Transitions TransitionMap
	Clocks      ClockMap
}

func New(statesAmt int) TimeCEA {
	return TimeCEA{StatesAmt: statesAmt, ClocksAmt: 2, Q0: 0, Transitions: TransitionMap{}, Clocks: ClockMap{0: []types.Clock{{Clock1: 0, Clock2: 0}}}}
}

func (tc *TimeCEA) RegisterTransition(tr types.Transition) {
	tc.Transitions[tr.Input] = append(tc.Transitions[tr.Input], tr.Output)
}

// O(n^3)
func (tc *TimeCEA) ReceiveWord(w types.Word) {
	newClocks := ClockMap{}

	for initialState, clocks := range tc.Clocks {
		for _, clock := range clocks {
			tc.ExecuteTransition(initialState, clock, w, newClocks)
		}
	}
	for state, clocks := range newClocks {
		// O(n^2)
		newClocks[state] = types.GetPareto(clocks)
	}
	tc.Clocks = newClocks
}

// O(n)
func (tc *TimeCEA) ExecuteTransition(currentState types.State, clock types.Clock, word types.Word, clocks ClockMap) {
	if outputs, ok := tc.Transitions[types.TransitionInput{P: currentState, L: word.Lettr}]; ok {
		for _, output := range outputs {
			newClock := clock.ProcessWordNewClock(word)
			if newClock.CheckCondition(output.Cond) {
				newClock.ResetClock(output.Resets)
				clocks.AddClock(output.Q, newClock)
			}
		}
	}
}

// O(N^4)
func (tc *TimeCEA) TestAutomataForPareto(maxWlen int, paretoNum int) int {
	for length := range maxWlen + 1 {
		for _, clock := range tc.Clocks {
			if len(clock) >= paretoNum {
				return length
			}
		}
		tc.ReceiveWord(types.Word{Lettr: 'a', TimeDelta: 1})
	}
	return -1
}
