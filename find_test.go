package sqroot

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFirstN(t *testing.T) {
	hits := FindFirstN(fakeNumber, []int{3, 4}, 3)
	assert.Equal(t, []int{2, 12, 22}, hits)
}

func TestFindAll(t *testing.T) {
	hits := FindAll(fakeNumber.WithSignificant(40), []int{3, 4})
	assert.Equal(t, []int{2, 12, 22, 32}, hits)
}

func TestFind(t *testing.T) {
	pattern := []int{3, 4}
	matches := Find(fakeNumber, pattern)
	pattern[0] = 5
	pattern[1] = 7
	assert.Equal(t, 2, matches())
	assert.Equal(t, 12, matches())
	assert.Equal(t, 22, matches())
}

func TestFindFirstNSingle(t *testing.T) {
	hits := FindFirstN(fakeNumber, []int{1}, 4)
	assert.Equal(t, []int{0, 10, 20, 30}, hits)
}

func TestFindFirst(t *testing.T) {
	assert.Equal(t, 5, FindFirst(fakeNumber, []int{6, 7, 8}))
}

func TestFindFirstNotThere(t *testing.T) {
	assert.Equal(t, -1, FindFirst(Sqrt(100489), []int{5}))
}

func TestFindEmptyPattern(t *testing.T) {
	hits := FindFirstN(fakeNumber, nil, 4)
	assert.Equal(t, []int{0, 1, 2, 3}, hits)
	assert.Equal(t, 0, FindFirst(fakeNumber, nil))
}

func TestFindEmptyPatternIterator(t *testing.T) {
	iter := Find(fakeNumber.WithSignificant(4), nil)
	assert.Equal(t, 0, iter())
	assert.Equal(t, 1, iter())
	assert.Equal(t, 2, iter())
	assert.Equal(t, 3, iter())
	assert.Equal(t, -1, iter())
}

func TestFindFirstNTrickyPattern(t *testing.T) {
	// 12212212122122121221221 ** 2
	radican, ok := new(big.Int).SetString(
		"149138124915706483400311993274596508420730841", 10)
	assert.True(t, ok)
	number := SqrtBigInt(radican)
	hits := FindFirstN(
		number,
		[]int{1, 2, 2, 1, 2, 1, 2, 2, 1, 2, 2, 1},
		3,
	)
	assert.Equal(t, []int{3, 11}, hits)
}

func TestFindLast(t *testing.T) {
	n := fakeNumber.WithSignificant(1000)
	assert.Equal(t, 998, FindLast(n, []int{9, 0}))
	assert.Equal(t, 994, FindLast(n, []int{5, 6}))
	assert.Equal(t, -1, FindLast(n, []int{5, 7}))
	assert.Equal(t, 2, FindLast(Sqrt(1522756), []int{3, 4}))
}

func TestFindLastN(t *testing.T) {
	hits := FindLastN(fakeNumber.WithSignificant(1200), []int{5, 6}, 3)
	assert.Equal(t, []int{1194, 1184, 1174}, hits)
	n := fakeNumber.WithSignificant(1000)
	hits = FindLastN(n, []int{5, 6}, 3)
	assert.Equal(t, []int{994, 984, 974}, hits)
	hits = FindLastN(n.WithStart(975), []int{5, 6}, 3)
	assert.Equal(t, []int{994, 984}, hits)
	hits = FindLastN(n, nil, 4)
	assert.Equal(t, []int{999, 998, 997, 996}, hits)
	hits = FindLastN(n, []int{5, 7, 9}, 3)
	assert.Empty(t, hits)
	hits = FindLastN(n, []int{5, 6}, 0)
	assert.Empty(t, hits)
}

func TestFindZeroNumber(t *testing.T) {
	var n Number
	assert.Equal(t, -1, FindFirst(&n, []int{5}))
	assert.Equal(t, -1, FindFirst(&n, nil))
	assert.Empty(t, FindFirstN(&n, []int{5}, 3))
	assert.Empty(t, FindFirstN(&n, nil, 3))
	assert.Empty(t, FindAll(&n, []int{5}))
	assert.Empty(t, FindAll(&n, nil))
	assert.Equal(t, -1, FindLast(&n, []int{5}))
	assert.Equal(t, -1, FindLast(&n, nil))
	assert.Empty(t, FindLastN(&n, []int{5}, 3))
	assert.Empty(t, FindLastN(&n, nil, 3))
}
