package timecea

import (
	"automata/types"
	"container/list"
)

// transitionmap podría ser una lista 2d de listas de outputs (transition input es una tupla (int, byte))
// ocuparía harta mas memoria, y crear un nuevo TimeCEA sería lento
// posiblemente sería bueno usar un sparse array si es que es rápido
// como statesAmt es muy chiquito (por ahora siempre 2) hacer 255 arrays de tamaño statesAmt no ocuparía tanta memoria
// type [][][]types.TransitionOutput
type TransitionMap map[types.TransitionInput][]types.TransitionOutput

type ClockArray []*list.List

func (cm ClockArray) AddClock(q types.State, clock types.Clock) bool {
	clockSet := cm[q]
	if clockSet == nil {
		cm[q] = list.New()
		clockSet = cm[q]
	}

	for e := clockSet.Front(); e != nil; {
		otherclock := e.Value.(types.Clock)
		if otherclock.LessOrEqualThan(clock) {
			return false
		}
		if clock.LessOrEqualThan(otherclock) {
			n := e.Next()
			clockSet.Remove(e)
			e = n
		} else {
			e = e.Next()
		}
	}
	clockSet.PushBack(clock)
	return true
}

type TimeCEA struct {
	paretoFound bool
	StatesAmt   int
	ClocksAmt   int
	Q0          types.State
	Transitions TransitionMap
	Clocks      ClockArray
}

func New(statesAmt int) TimeCEA {
	newSet := list.New()
	newSet.PushBack(types.Clock{Clock1: 0, Clock2: 0})
	clocks := make(ClockArray, statesAmt)
	clocks[0] = newSet
	return TimeCEA{
		paretoFound: false,
		StatesAmt:   statesAmt,
		ClocksAmt:   2,
		Q0:          0,
		Transitions: TransitionMap{},
		Clocks:      clocks,
	}
}

func (tc *TimeCEA) RegisterTransition(tr types.Transition) {
	outputs, ok := tc.Transitions[tr.Input]
	if !ok {
		tc.Transitions[tr.Input] = append(make([]types.TransitionOutput, 0, tc.StatesAmt), tr.Output)
	} else {

		tc.Transitions[tr.Input] = append(outputs, tr.Output)
	}
}

// O(n^3)
func (tc *TimeCEA) ReceiveWord(w types.Word) bool {
	// newClocks := make(ClockArray, tc.StatesAmt)
	flag := false
	for initialState, clocks := range tc.Clocks {
		if clocks == nil {
			continue
		}
		for e := clocks.Front(); e != nil; e = e.Next() {
			clock := e.Value.(types.Clock)
			flag = tc.ExecuteTransition(types.State(initialState), clock, w, tc.Clocks) || flag
		}
	}
	return flag
	// tc.Clocks = newClocks
}

// O(n)
func (tc *TimeCEA) ExecuteTransition(currentState types.State, clock types.Clock, word types.Word, clocks ClockArray) bool {
	flag := false
	if outputs, ok := tc.Transitions[types.TransitionInput{P: currentState, L: word.Lettr}]; ok {
		for _, output := range outputs {
			newClock := clock.ProcessWordNewClock(word)
			if newClock.CheckCondition(output.Cond) {
				newClock.ResetClock(output.Resets)
				flag = clocks.AddClock(output.Q, newClock) || flag
			}
		}

	}
	clocks.AddClock(tc.Q0, types.Clock{Clock1: 0, Clock2: 0})
	return flag
}

// O(N^4)
func (tc *TimeCEA) TestAutomataForPareto(maxWlen int, paretoNum int) int {
	for length := range maxWlen + 1 {
		for _, clock := range tc.Clocks {
			if clock == nil {
				continue
			}
			if clock.Len() >= paretoNum {
				return length
			}
		}
		if !tc.ReceiveWord(types.Word{Lettr: 'a', TimeDelta: 1}) {
			return -1
		}

	}
	return -1
}
