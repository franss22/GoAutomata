package ones

import (
	"github.com/fatih/color"
)

type Ones struct {
	//List of ones and 0s
	List []int
	//Size of the list (to avoid calling len() each time)
	Size int
	//Amount of 1s in the list
	OneAmount int
	//Amount of contiguous ones stacked at the end of the list
	RightmostOnes int
	//Index of the currently moving 1
	MovingIndex int
}

func New(Size int) Ones {

	return Ones{List: make([]int, Size), Size: Size}
}

func (o *Ones) LastOneIndex() int {
	last_bit := o.Size - 1 - o.RightmostOnes
	for last_bit >= 0 && (o.List)[last_bit] == 0 {
		last_bit--
	}
	return last_bit
}

// 0010011

func (o *Ones) Reset(onesAmt int) {
	o.OneAmount = onesAmt
	o.RightmostOnes = 0
	o.List = make([]int, o.Size)
	for i := range onesAmt {
		o.List[i] = 1
	}
	o.MovingIndex = onesAmt - 1
}

func (o *Ones) CanMove(index int) bool {
	if index >= o.Size-1 {
		return false
	}
	return o.List[index+1] == 0
}

func (o *Ones) MoveOneRightOnce(index int) {

	if o.List[index] == 0 {
		panic("Tried to move a 1 from an index with value 0")
	}
	if o.List[index+1] == 1 {
		panic("Tried to move a 1 to an index with value 1")
	}

	o.List[index] = 0
	o.List[index+1] = 1
	// color.Cyan("Moved 1 (index %d) to the right: %v", index, o.List)

	// return true
}

// Checks if there any slots with 0 value to the right of index
// RightmostOnes must have been previsously counted
func (o *Ones) CheckRightsideGaps(index int) bool {
	return index < o.Size-1-o.RightmostOnes
}

// Reagroupa todos los rightmostOnes a la derecha de index, y rellena el resto con 0s
func (o *Ones) RegroupLeft(index int) {
	//Mueve todos los rightmostOnes a la derecha de index
	for range o.RightmostOnes {
		o.List[index] = 1
		index++
	}
	//Si se rellena con 0s, no hay rightmostOnes
	if index < o.Size {
		o.RightmostOnes = 0
	}
	//rellena hasta el final con 0s
	for index < o.Size {
		o.List[index] = 0
		index++
	}
}

func (o *Ones) CountRightmostOnes() {
	last_bit := o.Size - 1 - o.RightmostOnes
	for last_bit >= 0 && o.List[last_bit] == 1 {
		last_bit--
		o.RightmostOnes++
	}
}

// Moves the noxt one in o.List.
// Returns true while it can still move, false when no more ones can be moved
func (o *Ones) Next() bool {
	//fmt.print(color.YellowString("%d -> ", o.List))

	if o.OneAmount == 0 || o.OneAmount == o.Size {
		return false
	}
	if o.CanMove(o.MovingIndex) { //A: Mover I a la derecha
		//fmt.print(color.GreenString("A (%d): ", o.MovingIndex))
		o.MoveOneRightOnce(o.MovingIndex)
		o.MovingIndex++
		return true
	} else { //B: I llegó al final, hay que mover el siguiente 1
		//fmt.print(color.GreenString(" B: "))

		o.CountRightmostOnes()
		if o.RightmostOnes == o.OneAmount { // E: Todos los 1s están al final de la lista
			//fmt.print(color.GreenString("E: "))
			return false
		}
		o.MovingIndex = o.LastOneIndex()

		if !o.CanMove(o.MovingIndex) { //...huh?
			color.Red("huh?=%v", o)
			return false
		}

		o.MoveOneRightOnce(o.MovingIndex)
		o.MovingIndex++
		if o.CheckRightsideGaps(o.MovingIndex) { //C

			//fmt.print(color.GreenString("C: "))
			o.RegroupLeft(o.MovingIndex + 1)
			o.MovingIndex = o.LastOneIndex()
			return true

		} else { //D
			//fmt.print(color.GreenString("D: "))
			o.RightmostOnes++
			return true
		}

	}

}

/*

A - Avanzar el último 1 hasta llegar al final
B - Buscar el siguiente 1
C - Si ese 1 puede avanzar, avanzar y mover los rightmostones a su derecha. GOTO A
D - Si no puede, GOTO A
E - Si todos los 1s están al final, return false
111000 A
  ^
110100 A
   ^
110010 A
    ^
110001 A -> B
     ^<
101100 C -> A
  ^^
101010 A
    ^
101001 A -> B
     ^<
100110 C -> A
   ^^
100101 A -> B
     ^<
100011 C1 -> A -> B
	^<
*/
