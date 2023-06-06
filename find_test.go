package sqroot

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFirstN(t *testing.T) {
	number := Sqrt(5)
	hits := FindFirstN(number, []int{9, 7}, 3)
	assert.Equal(t, []int{7, 12, 59}, hits)
}

func TestFindAll(t *testing.T) {
	number := Sqrt(5).WithSignificant(100)
	hits := FindAll(number, []int{9, 7})
	assert.Equal(t, []int{7, 12, 59}, hits)
}

func TestFind(t *testing.T) {
	number := Sqrt(5)
	pattern := []int{9, 7}
	matches := Find(number, pattern)
	pattern[0] = 2
	pattern[1] = 3
	assert.Equal(t, 7, matches())
	assert.Equal(t, 12, matches())
	assert.Equal(t, 59, matches())
}

func TestFindFirstNSingle(t *testing.T) {
	number := Sqrt(11)
	hits := FindFirstN(number, []int{3}, 4)
	assert.Equal(t, []int{0, 1, 10, 13}, hits)
}

func TestFindFirst(t *testing.T) {
	number := Sqrt(2)
	assert.Equal(t, 1, FindFirst(number, []int{4, 1, 4}))
}

func TestFindFirstNotThere(t *testing.T) {
	number := Sqrt(100489)
	assert.Equal(t, -1, FindFirst(number, []int{5}))
}

func TestFindEmptyPattern(t *testing.T) {
	number := Sqrt(2)
	hits := FindFirstN(number, nil, 4)
	assert.Equal(t, []int{0, 1, 2, 3}, hits)
	assert.Equal(t, 0, FindFirst(number, nil))
}

func TestFindEmptyPatternIterator(t *testing.T) {
	number := Sqrt(2).WithSignificant(4)
	iter := Find(number, nil)
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
	number := Sqrt(5).WithSignificant(1000)
	assert.Equal(t, 936, FindLast(number, []int{9, 7}))
	assert.Equal(t, -1, FindLast(number, []int{0, 1, 2, 3, 4}))
}

func TestFindLastN(t *testing.T) {
	number := Sqrt(5).WithSignificant(1000)
	hits := FindLastN(number, []int{9, 7}, 3)
	assert.Equal(t, []int{936, 718, 600}, hits)
	hits = FindLastN(number.WithStart(601), []int{9, 7}, 3)
	assert.Equal(t, []int{936, 718}, hits)
	hits = FindLastN(number.WithStart(1001), []int{9, 7}, 3)
	assert.Empty(t, hits)
	hits = FindLastN(Sqrt(5).WithSignificant(1300), []int{9, 7}, 3)
	assert.Equal(t, []int{1276, 1221, 936}, hits)
	hits = FindLastN(Sqrt(5).WithSignificant(0), []int{9, 7}, 3)
	assert.Empty(t, hits)
	hits = FindLastN(number, nil, 4)
	assert.Equal(t, []int{999, 998, 997, 996}, hits)
	hits = FindLastN(number, []int{1, 2, 3}, 3)
	assert.Equal(t, []int{815, 579}, hits)
	hits = FindLastN(number, []int{1, 2, 3}, 0)
	assert.Empty(t, hits)
	short := Sqrt(1522756)
	assert.Equal(t, 2, FindLast(short, []int{3, 4}))
}

func TestFindLastNMemoize(t *testing.T) {
	n := Sqrt(5).WithMemoize()
	hits := FindLastN(n.WithSignificant(1300), []int{9, 7}, 3)
	assert.Equal(t, []int{1276, 1221, 936}, hits)
	hits = FindLastN(n.WithSignificant(0), []int{9, 7}, 3)
	assert.Empty(t, hits)
	n1000 := n.WithSignificant(1000)
	hits = FindLastN(n1000, []int{9, 7}, 3)
	assert.Equal(t, []int{936, 718, 600}, hits)
	hits = FindLastN(n1000.WithStart(601), []int{9, 7}, 3)
	assert.Equal(t, []int{936, 718}, hits)
	hits = FindLastN(n1000.WithStart(1001), []int{9, 7}, 3)
	assert.Empty(t, hits)
	hits = FindLastN(n1000, nil, 4)
	assert.Equal(t, []int{999, 998, 997, 996}, hits)
	hits = FindLastN(n1000, []int{1, 2, 3}, 3)
	assert.Equal(t, []int{815, 579}, hits)
	hits = FindLastN(n1000, []int{1, 2, 3}, 0)
	assert.Empty(t, hits)
	short := Sqrt(1522756).WithMemoize()
	assert.Equal(t, 2, FindLast(short, []int{3, 4}))
}

func TestFindLastNDigits(t *testing.T) {
	str := "v1:01201[10]010101[20]120101"
	var d Digits
	d.UnmarshalText([]byte(str))
	pattern := []int{0, 1, 0, 1}
	hits := FindLastN(d, pattern, 4)
	assert.Equal(t, []int{22, 12, 10}, hits)
	hits = FindLastN(d, pattern, 3)
	assert.Equal(t, []int{22, 12, 10}, hits)
	hits = FindLastN(d, pattern, 2)
	assert.Equal(t, []int{22, 12}, hits)
	hits = FindLastN(d, pattern, 0)
	assert.Empty(t, hits)
	hits = FindLastN(d, nil, 8)
	assert.Equal(t, []int{25, 24, 23, 22, 21, 20, 15, 14}, hits)
}

func TestFindLastNDigits2(t *testing.T) {
	digits2 := AllDigits(Sqrt(5).WithSignificant(1300))
	digits := digits2.WithEnd(1000)
	hits := FindLastN(digits, []int{9, 7}, 3)
	assert.Equal(t, []int{936, 718, 600}, hits)
	hits = FindLastN(digits2, []int{9, 7}, 3)
	assert.Equal(t, []int{1276, 1221, 936}, hits)
	hits = FindLastN(digits, nil, 4)
	assert.Equal(t, []int{999, 998, 997, 996}, hits)
	hits = FindLastN(digits, []int{1, 2, 3}, 3)
	assert.Equal(t, []int{815, 579}, hits)
	hits = FindLastN(digits, []int{1, 2, 3}, 0)
	assert.Empty(t, hits)
	assert.Equal(t, 936, FindLast(digits, []int{9, 7}))
	assert.Equal(t, -1, FindLast(digits, []int{0, 1, 2, 3, 4}))
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

func TestFindZeroDigits(t *testing.T) {
	var d Digits
	assert.Equal(t, -1, FindFirst(d, []int{5}))
	assert.Equal(t, -1, FindFirst(d, nil))
	assert.Empty(t, FindFirstN(d, []int{5}, 3))
	assert.Empty(t, FindFirstN(d, nil, 3))
	assert.Empty(t, FindAll(d, []int{5}))
	assert.Empty(t, FindAll(d, nil))
	assert.Equal(t, -1, FindLast(d, []int{5}))
	assert.Equal(t, -1, FindLast(d, nil))
	assert.Empty(t, FindLastN(d, []int{5}, 3))
	assert.Empty(t, FindLastN(d, nil, 3))
}
