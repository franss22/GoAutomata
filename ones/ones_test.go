package ones

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func ListEquals(t *testing.T, value []int, expected []int) {
// 	if len(value) != len(expected) {
// 		t.Errorf("Diff arrays:\n\tgot=%v\n\texp=%v\nDiff length %d vs %d", value, expected, len(value), len(expected))

// 	}

// 	for i := range len(value) {
// 		if value[i] != expected[i] {
// 			t.Errorf("Diff arrays:\n\tgot=%v\n\texp=%v\nDiff values %d vs %d at index %d", value, expected, value[i], expected[i], i)
// 		}
// 	}
// }

func assertOnes(t *testing.T, ones Ones, list []int, onesAmt int, rightmostOnes int) {
	assert.Equal(t, list, ones.List)
	assert.Equal(t, len(list), ones.Size)
	assert.Equal(t, onesAmt, ones.OneAmount)
	assert.Equal(t, rightmostOnes, ones.RightmostOnes)
}

func TestConstructor(t *testing.T) {
	ones := New(5)
	assertOnes(t, ones, []int{0, 0, 0, 0, 0}, 0, 0)
}

func TestReset(t *testing.T) {
	ones := New(5)
	ones.Reset(4)
	assertOnes(t, ones, []int{1, 1, 1, 1, 0}, 4, 0)

}
func TestResetFull(t *testing.T) {
	ones := New(5)
	ones.Reset(5)
	assertOnes(t, ones, []int{1, 1, 1, 1, 1}, 5, 0)

}
func TestMoveNextOnce(t *testing.T) {
	ones := New(5)
	ones.Reset(1)
	ok := ones.Next()
	assert.Equal(t, true, ok)

	assertOnes(t, ones, []int{0, 1, 0, 0, 0}, 1, 0)
}

func TestMoveNextTwice(t *testing.T) {
	ones := New(5)
	ones.Reset(1)
	ok := ones.Next()
	assert.Equal(t, true, ok)
	ok = ones.Next()
	assert.Equal(t, true, ok)

	assertOnes(t, ones, []int{0, 0, 1, 0, 0}, 1, 0)
}

func TestLastOneIndex(t *testing.T) {
	ones := New(5)
	ones.Reset(2)
	i := ones.LastOneIndex()

	assert.Equal(t, 1, i)
	assertOnes(t, ones, []int{1, 1, 0, 0, 0}, 2, 0)
}

func TestLastOneIndexWithoutSetRightmostOnes(t *testing.T) {
	ones := New(5)
	ones.Reset(4)
	ones.List = []int{1, 1, 0, 1, 1}
	i := ones.LastOneIndex()

	assert.Equal(t, 4, i)
	assertOnes(t, ones, []int{1, 1, 0, 1, 1}, 4, 0)
}

func TestLastOneIndexWithSetRightmostOnes(t *testing.T) {
	ones := New(5)
	ones.Reset(4)
	ones.List = []int{1, 1, 0, 1, 1}
	ones.RightmostOnes = 2
	i := ones.LastOneIndex()

	assert.Equal(t, 1, i)
	assertOnes(t, ones, []int{1, 1, 0, 1, 1}, 4, 2)
}

func TestCountRightmostOnes(t *testing.T) {
	ones := New(5)
	ones.Reset(4)
	ones.List = []int{1, 1, 0, 1, 1}
	ones.CountRightmostOnes()
	assertOnes(t, ones, []int{1, 1, 0, 1, 1}, 4, 2)
}
