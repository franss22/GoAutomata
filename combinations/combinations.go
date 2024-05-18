package combinations

type Nums struct {
	List      []int
	ValsAmt   int
	MaxVal    int
	ItersDone int
}

type Combinator interface {
	Next() bool
	Indexes() []int
	CurrentIteration() int
}

func (n *Nums) Indexes() []int {
	return n.List
}

func (n *Nums) CurrentIteration() int {
	return n.ItersDone
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
	n.ItersDone = 0
}

func (n *Nums) Next() bool {
	if n.ValsAmt > n.MaxVal || n.ValsAmt == 0 {
		return false
	}
	index := n.ValsAmt - 1
	max := n.MaxVal
	n.ItersDone++

	return n.incrementVal(index, max)

}
func (n *Nums) Advance(iters int) *Nums {
	for range iters {
		n.Next()
	}
	return n
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

type PNums struct {
	Nums             Nums
	Iters            int
	ItersDoneAtStart int
}

func (pn *PNums) Next() bool {
	pn.Iters--
	return pn.Nums.Next() && pn.Iters > 0
}

func (pn *PNums) Indexes() []int {
	return pn.Nums.List
}

func (pn *PNums) CurrentIteration() int {
	return pn.Nums.ItersDone
}

func (n *Nums) NewPnums(maxIters int) PNums {
	newNums := New(n.MaxVal + 1)
	newNums.List = make([]int, len(n.List))
	newNums.ValsAmt = n.ValsAmt
	copy(newNums.List, n.List)
	return PNums{Nums: newNums, Iters: maxIters, ItersDoneAtStart: n.ItersDone}
}
