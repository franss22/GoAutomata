package combinations

import "automata/ones"

type Nums struct {
	List    []int
	ValsAmt int
	MaxVal  int
}

func New(size int) Nums {
	return Nums{MaxVal: size - 1}
}

func (n *Nums) Reset(valsAmt int) {
	if valsAmt > n.MaxVal+1 {
		panic("Cannot have more values than maxval+1")
	}
	n.ValsAmt = valsAmt
	l := make([]int, valsAmt)
	for i := range valsAmt {
		l[i] = i
	}
	n.List = l
}

func (n *Nums) Next() bool {
	if n.ValsAmt > n.MaxVal || n.ValsAmt == 0 {
		return false
	}
	index := n.ValsAmt - 1
	max := n.MaxVal

	return n.incrementVal(index, max)

}

func (n *Nums) incrementVal(index int, max int) bool {
	val := n.List[index]
	if val == max {
		if index == 0 {
			return false
		}
		if !n.incrementVal(index-1, max-1) {
			return false
		}
		prevVal := n.List[index-1]
		n.List[index] = prevVal + 1
		return true

	} else {
		n.List[index]++
		return true
	}

}

func (n *Nums) ToOnes() ones.Ones {
	size := n.MaxVal + 1
	ones := ones.New(size)
	// ones.List = make([]int, size)
	// color.Red("%v", ones.List)
	ones.OneAmount = n.ValsAmt
	for _, i := range n.List {
		ones.List[i] = 1
	}
	ones.CountRightmostOnes()
	return ones
}
