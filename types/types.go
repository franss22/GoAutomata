package types

import (
	"sort"
)

const NO_COND Time = -1

type Letter byte
type State int
type Time int

type Word struct {
	Lettr     Letter
	TimeDelta Time
}

type ClockCondition struct {
	// Condicion de Clock 1
	Cond1 Time
	// Condicion de Clock 2
	Cond2 Time
}

type ClockReset struct {
	// Reset de Clock 1
	Reset1 bool
	// Reset de Clock 2
	Reset2 bool
}

type Clock struct {
	Clock1 Time
	Clock2 Time
}

type ClockSet map[Clock]struct{}

func (c *Clock) LessOrEqualThan(c2 Clock) bool {
	return c.Clock1 <= c2.Clock1 && c.Clock2 <= c2.Clock2
}

func (c *Clock) CheckCondition(conds ClockCondition) bool {

	check1 := conds.Cond1 == NO_COND || (c.Clock1 <= conds.Cond1)
	check2 := conds.Cond2 == NO_COND || (c.Clock2 <= conds.Cond2)
	return check1 && check2

}

func (c *Clock) ResetClock(reset ClockReset) {
	if reset.Reset1 {
		c.Clock1 = 0
	}
	if reset.Reset2 {
		c.Clock2 = 0
	}
}

func (c *Clock) ProcessWordNewClock(w Word) Clock {
	return Clock{c.Clock1 + w.TimeDelta, c.Clock2 + w.TimeDelta}
}

// O(n)
// (2,1) (1,2) (3,3) devuelve (2,1) (1,2)
func (clock *Clock) IsMinimal(clocks []Clock) bool {
	for _, otherclock := range clocks {
		if otherclock == *clock {
			continue
		}
		if otherclock.LessOrEqualThan(*clock) {
			return false
		}
	}
	return true
}

func (clockset ClockSet) Sort() []Clock {
	result := []Clock{}
	for clock := range clockset {
		result = append(result, clock)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].LessOrEqualThan(result[j])
	})

	return result
}

// O(n^2) (podría ser O(nlogn))
func GetPareto(cs ClockSet) ClockSet {
	ret := ClockSet{}
	sortedClocks := cs.Sort()

	// fmt.Print(sortedClocks, "\n")
	for clock := range cs {
		if clock.IsMinimal(sortedClocks) {
			ret[clock] = struct{}{}
		}
	}
	return ret
}

type Transition struct {
	Input  TransitionInput
	Output TransitionOutput
}

type TransitionInput struct {
	// Estado inicial
	P State
	// Letra leida
	L Letter
}

type TransitionOutput struct {
	//Estado Final
	Q State
	// Condiciones de los Clocks
	Cond ClockCondition
	// Resets de los Clocks
	Resets ClockReset
}
