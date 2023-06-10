package sqroot

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/keep94/consume2"
	"github.com/stretchr/testify/assert"
)

func TestNumberReusable(t *testing.T) {
	radican := big.NewInt(5)
	n := SqrtBigInt(radican)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.22360 679", n.Sprint(UpTo(8)))
	assert.Equal(t, big.NewInt(5), radican)
	radican.SetInt64(7)
	assert.Equal(t, "0.22360 679", n.Sprint(UpTo(8)))
	assert.Equal(t, big.NewInt(7), radican)
}

func Test2(t *testing.T) {
	n := Sqrt(2)
	assert.False(t, n.IsZero())
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.14142 13562", n.Sprint(UpTo(10)))
}

func Test3(t *testing.T) {
	n := Sqrt(3)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.17320 50807", n.Sprint(UpTo(10)))
}

func Test0(t *testing.T) {
	n := Sqrt(0)
	assert.Zero(t, *zeroNumber)
	assert.Same(t, zeroNumber, n)
	iter := n.Iterator()
	assert.Equal(t, -1, iter())
}

func Test1(t *testing.T) {
	n := Sqrt(1)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.1", n.Sprint(UpTo(10)))
}

func Test100489(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, 3, n.Exponent())
	assert.Equal(t, "0.317", n.Sprint(UpTo(10)))
}

func Test100489Iterator(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, 3, n.Exponent())
	iter := n.Iterator()
	assert.Equal(t, 3, iter())
	assert.Equal(t, 1, iter())
	assert.Equal(t, 7, iter())
	assert.Equal(t, -1, iter())
	assert.Equal(t, -1, iter())
	iter = n.Iterator()
	assert.Equal(t, 3, iter())
	assert.Equal(t, 1, iter())
	assert.Equal(t, 7, iter())
	assert.Equal(t, -1, iter())
	assert.Equal(t, -1, iter())
}

func TestIteratorPersistence(t *testing.T) {
	n := Sqrt(7)
	iter := n.Iterator()
	n = Sqrt(11)
	assert.Equal(t, 2, iter())
	assert.Equal(t, 6, iter())
	assert.Equal(t, 4, iter())
	assert.Equal(t, 5, iter())
}

func TestIteratorAt(t *testing.T) {
	n := Sqrt(100489)
	iter := n.IteratorAt(3)
	assert.Equal(t, -1, iter())
	iter = n.IteratorAt(2)
	assert.Equal(t, 7, iter())
	assert.Equal(t, -1, iter())
	iter = n.IteratorAt(0)
	assert.Equal(t, 3, iter())
	assert.Equal(t, 1, iter())
	assert.Equal(t, 7, iter())
	assert.Equal(t, -1, iter())
	assert.Panics(t, func() { n.IteratorAt(-1) })
}

func TestNegative(t *testing.T) {
	assert.Panics(t, func() { Sqrt(-1) })
}

func Test256(t *testing.T) {
	n := Sqrt(256)
	assert.Equal(t, 2, n.Exponent())
	assert.Equal(t, "0.16", n.Sprint(UpTo(10)))
}

func Test40(t *testing.T) {
	n := Sqrt(40)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.63245 55320", n.Sprint(UpTo(10)))
}

func Test0026(t *testing.T) {
	n := SqrtRat(2600, 1000000)
	assert.Equal(t, -1, n.Exponent())
	assert.Equal(t, "0.50990 19513", n.Sprint(UpTo(10)))
}

func Test026(t *testing.T) {
	n := SqrtRat(26, 1000)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.16124 51549", n.Sprint(UpTo(10)))
}

func Test2401Over400(t *testing.T) {
	n := SqrtRat(2401, 4)
	assert.Equal(t, 2, n.Exponent())
	assert.Equal(t, "0.245", n.Sprint(UpTo(10)))
}

func Test3Over7(t *testing.T) {
	n := SqrtRat(3, 7)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.65465 36707 0797", n.Sprint(UpTo(14)))
}

func Test3Over70000Reusable(t *testing.T) {
	radican := big.NewRat(3, 70000)
	n := SqrtBigRat(radican)
	assert.Equal(t, -2, n.Exponent())
	assert.Equal(t, "0.65465 36707 0797", n.Sprint(UpTo(14)))
	assert.Equal(t, big.NewRat(3, 70000), radican)
	radican.Num().SetInt64(5)
	radican.Denom().SetInt64(80000)
	assert.Equal(t, "0.65465 36707 0797", n.Sprint(UpTo(14)))
	assert.Equal(t, big.NewInt(5), radican.Num())
	assert.Equal(t, big.NewInt(80000), radican.Denom())
}

func TestSquareRootFixed(t *testing.T) {
	number := Sqrt(10)
	actual := fmt.Sprintf("%f", number)
	assert.Equal(t, "3.162277", actual)
}

func TestSquareRootString(t *testing.T) {
	number := Sqrt(10)
	assert.Equal(t, "3.162277660168379", number.String())
}

func TestCubeRoot2(t *testing.T) {
	assert.Equal(t, "1.25992104989487", fmt.Sprintf("%.15g", CubeRoot(2)))
}

func TestCubeRoot2Big(t *testing.T) {
	n := CubeRootBigInt(big.NewInt(2))
	assert.Equal(t, "1.25992104989487", fmt.Sprintf("%.15g", n))
}

func TestCubeRoot35223040952(t *testing.T) {
	n := CubeRoot(35223040952)
	assert.Equal(t, "3278", n.String())
	assert.Equal(t, 4, n.Exponent())
	iter := n.Iterator()
	assert.Equal(t, 3, iter())
	assert.Equal(t, 2, iter())
	assert.Equal(t, 7, iter())
	assert.Equal(t, 8, iter())
	assert.Equal(t, -1, iter())
	assert.Equal(t, -1, iter())
}

func TestCubeRootRat(t *testing.T) {
	n := CubeRootRat(35223040952, 8000)
	assert.Equal(t, "163.9", n.String())
}

func TestCubeRootBigRat(t *testing.T) {
	n := CubeRootBigRat(big.NewRat(35223040952, 8000))
	assert.Equal(t, "163.9", n.String())
}

func TestCubeRootSmallRat(t *testing.T) {
	n := CubeRootRat(2, 73952)
	assert.Equal(t, -1, n.Exponent())
	assert.Equal(t, "0.030016498129266", fmt.Sprintf("%.14g", n))
}

func TestNegDenom(t *testing.T) {
	radican := big.NewRat(1, 700)
	radican.Denom().SetInt64(-500)
	radican.Num().SetInt64(3)
	assert.Panics(t, func() { SqrtBigRat(radican) })
}

func TestWithSignificant(t *testing.T) {
	// Resolves to 6 significant digits
	n := Sqrt(2).WithSignificant(9).WithSignificant(6).WithSignificant(10)
	assert.Equal(t, "1.41421", n.String())
}

func TestWithSignificantPanics(t *testing.T) {
	var n Number
	assert.Panics(t, func() { n.WithSignificant(-1) })
}

func TestWithSignificantToZero(t *testing.T) {
	assert.Zero(t, *zeroNumber)
	assert.Same(t, zeroNumber, Sqrt(2).WithSignificant(0))
}

func TestAt(t *testing.T) {
	n := Sqrt(2)
	assert.Equal(t, 5, n.At(15))
	assert.Equal(t, 7, n.At(25))
	assert.Equal(t, -1, n.At(-1))
}

func TestAtFinite(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, 3, n.At(0))
	assert.Equal(t, 7, n.At(2))
	assert.Equal(t, -1, n.At(3))
}

func TestZeroNumber(t *testing.T) {
	var n Number
	assert.Equal(t, -1, n.At(0))
	assert.Zero(t, n.Exponent())
	assert.True(t, n.IsZero())
	assert.Zero(t, AllDigits(&n))
	assert.True(t, n.IsMemoize())
	assert.Same(t, &n, n.WithMemoize())
	assert.Same(t, &n, n.WithSignificant(5))
	assert.Equal(t, -1, n.Iterator()())
	assert.Equal(t, "0", n.String())
	s := n.WithSignificant(2000000000).WithStart(1900000000)
	assert.Zero(t, AllDigits(s))
}

func TestSameNumber(t *testing.T) {
	n := Sqrt(6)
	sixDigits := n.WithSignificant(6)
	assert.Same(t, sixDigits, sixDigits.WithSignificant(6))
	memoized := sixDigits.WithMemoize()
	assert.Same(t, memoized, memoized.WithMemoize())
	sevenDigits := memoized.WithSignificant(7)
	assert.Same(t, sevenDigits, sevenDigits.WithSignificant(8))
	assert.Same(t, sevenDigits, sevenDigits.WithMemoize())
}

func TestNumberWithStart(t *testing.T) {
	n := Sqrt(19)
	pattern := getSequentialDigits(n, 500, 2)
	assert.Less(t, FindFirst(n, pattern), 500)
	expected := findFirstNAfter(n, 500, pattern, 3)
	assert.Equal(t, 500, expected[0])
	actual := FindFirstN(n.WithStart(500), pattern, 3)
	assert.Equal(t, expected, actual)

	firstTwoResults := FindFirstN(
		n.WithSignificant(expected[2]+1).WithStart(500), pattern, 3)
	assert.Equal(t, expected[:2], firstTwoResults)
}

func TestNumberWithStartEmpty(t *testing.T) {
	n := Sqrt(19)
	s := n.WithSignificant(10).WithStart(300000)
	assert.Zero(t, AllDigits(s))
	s = n.WithSignificant(10).WithStart(10)
	assert.Zero(t, AllDigits(s))
}

func TestNumberWithStartZeroOrNegative(t *testing.T) {
	n := Sqrt(19)
	assert.Same(t, n, n.WithStart(0))
	assert.Same(t, n, n.WithStart(-1))
}

func TestFiniteLengthNumberWithStart(t *testing.T) {
	n := Sqrt(100489)
	s := n.WithStart(1)
	assert.Equal(t, []int{1}, FindAll(s, []int{1, 7}))
	assert.Equal(t, []int{1}, FindAll(s, []int{1, 7}))
	assert.Empty(t, FindAll(n.WithStart(2), []int{1, 7}))
	assert.Empty(t, FindAll(n.WithStart(300000), []int{1, 7}))
}

func TestNumberWithStartAndMemoize(t *testing.T) {
	n := Sqrt(23)
	pattern := getSequentialDigits(n, 500, 2)
	assert.Less(t, FindFirst(n, pattern), 500)
	expected := findFirstNAfter(n, 500, pattern, 3)
	assert.Equal(t, 500, expected[0])
	s := n.WithMemoize().WithStart(500)
	assert.Equal(t, expected, FindFirstN(s, pattern, 3))
	assert.Equal(t, expected, FindFirstN(s, pattern, 3))
}

func TestNumberGetDigits(t *testing.T) {
	n := Sqrt(2)
	var pb PositionsBuilder
	for i := 0; i < 10000; i += 2 {
		pb.Add(i)
	}
	p := pb.Build()
	assert.Equal(
		t, GetDigits(n, p).Sprint(), GetDigits(n.WithMemoize(), p).Sprint())
}

func findFirstNAfter(n *Number, start int, pattern []int, count int) []int {
	pipeline := consume2.PFilter(func(x int) bool { return x >= start })
	pipeline = consume2.Join(pipeline, consume2.PSlice[int](0, count))
	var result []int
	consume2.FromIntGenerator(Find(n, pattern), pipeline.AppendTo(&result))
	return result
}

func getSequentialDigits(n *Number, start, length int) []int {
	var pb PositionsBuilder
	digits := GetDigits(n, pb.AddRange(start, start+length).Build())
	result := make([]int, 0, length)
	for i := start; i < start+length; i++ {
		result = append(result, digits.At(i))
	}
	return result
}
