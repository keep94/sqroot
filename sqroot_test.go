package sqroot

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberReusable(t *testing.T) {
	radican := big.NewInt(5)
	n := SqrtBigInt(radican)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.22360 679", Sprint(n, UpTo(8)))
	assert.Equal(t, big.NewInt(5), radican)
	radican.SetInt64(7)
	assert.Equal(t, "0.22360 679", Sprint(n, UpTo(8)))
	assert.Equal(t, big.NewInt(7), radican)
}

func Test2(t *testing.T) {
	n := Sqrt(2)
	assert.False(t, n.IsZero())
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.14142 13562", Sprint(n, UpTo(10)))
}

func Test3(t *testing.T) {
	n := Sqrt(3)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.17320 50807", Sprint(n, UpTo(10)))
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
	assert.Equal(t, "0.1", Sprint(n, UpTo(10)))
}

func Test100489(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, 3, n.Exponent())
	assert.Equal(t, "0.317", Sprint(n, UpTo(10)))
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
	assert.Equal(t, "0.16", Sprint(n, UpTo(10)))
}

func Test40(t *testing.T) {
	n := Sqrt(40)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.63245 55320", Sprint(n, UpTo(10)))
}

func Test0026(t *testing.T) {
	n := SqrtRat(2600, 1000000)
	assert.Equal(t, -1, n.Exponent())
	assert.Equal(t, "0.50990 19513", Sprint(n, UpTo(10)))
}

func Test026(t *testing.T) {
	n := SqrtRat(26, 1000)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.16124 51549", Sprint(n, UpTo(10)))
}

func Test2401Over400(t *testing.T) {
	n := SqrtRat(2401, 4)
	assert.Equal(t, 2, n.Exponent())
	assert.Equal(t, "0.245", Sprint(n, UpTo(10)))
}

func Test3Over7(t *testing.T) {
	n := SqrtRat(3, 7)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.65465 36707 0797", Sprint(n, UpTo(14)))
}

func Test3Over70000Reusable(t *testing.T) {
	radican := big.NewRat(3, 70000)
	n := SqrtBigRat(radican)
	assert.Equal(t, -2, n.Exponent())
	assert.Equal(t, "0.65465 36707 0797", Sprint(n, UpTo(14)))
	assert.Equal(t, big.NewRat(3, 70000), radican)
	radican.Num().SetInt64(5)
	radican.Denom().SetInt64(80000)
	assert.Equal(t, "0.65465 36707 0797", Sprint(n, UpTo(14)))
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

func TestNumberAt(t *testing.T) {
	n := fakeNumber
	assert.Equal(t, -1, n.At(-1))
	assert.Equal(t, 3, n.At(322))
	assert.Equal(t, 1, n.At(0))
	assert.Equal(t, 2, n.At(1))
	assert.Equal(t, 3, n.At(102))
	assert.Equal(t, 0, n.At(399))
}

func TestNumberAtFiniteLength(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, -1, n.At(-1))
	assert.Equal(t, 7, n.At(2))
	assert.Equal(t, 3, n.At(0))
	assert.Equal(t, -1, n.At(3))
}

func TestNumberAtMemoize(t *testing.T) {
	n := fakeNumber.WithMemoize()
	assert.Equal(t, -1, n.At(-1))
	assert.Equal(t, 3, n.At(322))
	assert.Equal(t, 1, n.At(0))
	assert.Equal(t, 2, n.At(1))
	assert.Equal(t, 3, n.At(102))
	assert.Equal(t, 0, n.At(399))
}

func TestNumberAtSig(t *testing.T) {
	n := fakeNumber.WithSignificant(357)
	assert.Equal(t, -1, n.At(-1))
	assert.Equal(t, 3, n.At(322))
	assert.Equal(t, 1, n.At(0))
	assert.Equal(t, 4, n.At(303))
	assert.Equal(t, 7, n.At(356))
	assert.Equal(t, -1, n.At(357))
}

func TestNumberAtSigMemoize(t *testing.T) {
	n := fakeNumber.WithSignificant(357).WithMemoize()
	assert.Equal(t, -1, n.At(-1))
	assert.Equal(t, 3, n.At(322))
	assert.Equal(t, 1, n.At(0))
	assert.Equal(t, 4, n.At(303))
	assert.Equal(t, 7, n.At(356))
	assert.Equal(t, -1, n.At(357))
}

func TestNumberAtMemoizeSig(t *testing.T) {
	n := fakeNumber.WithMemoize().WithSignificant(357)
	assert.Equal(t, -1, n.At(-1))
	assert.Equal(t, 3, n.At(322))
	assert.Equal(t, 1, n.At(0))
	assert.Equal(t, 4, n.At(303))
	assert.Equal(t, 7, n.At(356))
	assert.Equal(t, -1, n.At(357))
}

func TestNumberInterfaces(t *testing.T) {
	n := fakeNumber
	assertStartsAt(t, n, 0)
	assert.False(t, hasReverse(n))
	assert.False(t, hasSubRange(n))
}

func TestNumberInterfacesSig(t *testing.T) {
	n := fakeNumber.WithSignificant(357)
	assertRange(t, n, 0, 357)
	assert.False(t, hasReverse(n))
	assert.False(t, hasSubRange(n))
}

func TestNumberInterfacesMemoize(t *testing.T) {
	n := fakeNumber.WithMemoize()
	assertStartsAt(t, n, 0)
	assert.True(t, hasReverse(n))
	assert.True(t, hasSubRange(n))
	assertRange(t, subRange(n, 62, 404), 62, 404)
	assertEmpty(t, subRange(n, 62, 62))
}

func TestNumberInterfacesSigMemoize(t *testing.T) {
	n := fakeNumber.WithSignificant(357).WithMemoize()
	assertRange(t, n, 0, 357)
	assert.True(t, hasReverse(n))
	assert.True(t, hasSubRange(n))
	assertRange(t, subRange(n, 62, 404), 62, 357)
	assertEmpty(t, subRange(n, 62, 62))
	assertEmpty(t, subRange(n, 357, 400))
}

func TestNumberInterfacesMemoizeSig(t *testing.T) {
	n := fakeNumber.WithMemoize().WithSignificant(357)
	assertRange(t, n, 0, 357)
	assert.True(t, hasReverse(n))
	assert.True(t, hasSubRange(n))
	assertRange(t, subRange(n, 62, 404), 62, 357)
	assertEmpty(t, subRange(n, 62, 62))
	assertEmpty(t, subRange(n, 357, 400))
}

func TestWithStart(t *testing.T) {
	n := fakeNumber
	seq := n.WithStart(423)
	assertStartsAt(t, seq, 423)
	assert.False(t, hasReverse(seq))
	assert.False(t, hasSubRange(seq))
}

func TestWithStartSig(t *testing.T) {
	n := fakeNumber.WithSignificant(541)
	seq := n.WithStart(423)
	assertRange(t, seq, 423, 541)
	assert.False(t, hasReverse(seq))
	assert.False(t, hasSubRange(seq))
	assertEmpty(t, n.WithStart(541))
	assertEmpty(t, n.WithStart(542))
}

func TestWithStartMemoize(t *testing.T) {
	n := fakeNumber.WithMemoize()
	seq := n.WithStart(423)
	assertStartsAt(t, seq, 423)
	assert.True(t, hasReverse(seq))
	assert.True(t, hasSubRange(seq))
	assertRange(t, subRange(seq, 357, 504), 423, 504)
}

func TestWithStartSigMemoize(t *testing.T) {
	n := fakeNumber.WithSignificant(541).WithMemoize()
	seq := n.WithStart(423)
	assertRange(t, seq, 423, 541)
	assert.True(t, hasReverse(seq))
	assert.True(t, hasSubRange(seq))
	assertRange(t, subRange(seq, 357, 600), 423, 541)
	assertEmpty(t, n.WithStart(541))
	assertEmpty(t, n.WithStart(542))
}

func TestWithStartMemoizeSig(t *testing.T) {
	n := fakeNumber.WithMemoize().WithSignificant(541)
	seq := n.WithStart(423)
	assertRange(t, seq, 423, 541)
	assert.True(t, hasReverse(seq))
	assert.True(t, hasSubRange(seq))
	assertRange(t, subRange(seq, 357, 600), 423, 541)
	assertEmpty(t, n.WithStart(541))
	assertEmpty(t, n.WithStart(542))
}

func assertStartsAt(t *testing.T, s Sequence, start int) {
	t.Helper()
	iter := s.digitIter()
	d, ok := iter()
	assert.True(t, ok)
	assert.Equal(t, start, d.Position)
	assert.Equal(t, (start+1)%10, d.Value)
}

func assertRange(t *testing.T, s Sequence, start, end int) {
	t.Helper()
	if assertForwardRange(t, s, start, end) && hasReverse(s) {
		assertReverseRange(t, s, start, end)
	}
}

func assertForwardRange(t *testing.T, s Sequence, start, end int) bool {
	t.Helper()
	iter := s.digitIter()
	for i := start; i < end; i++ {
		d, ok := iter()
		if !assert.True(t, ok) {
			return false
		}
		if !assert.Equal(t, i, d.Position) {
			return false
		}
		if !assert.Equal(t, (i+1)%10, d.Value) {
			return false
		}
	}
	_, ok := iter()
	return assert.False(t, ok)
}

func assertReverseRange(t *testing.T, s Sequence, start, end int) bool {
	t.Helper()
	r := s.(reverseSequence)
	iter := r.reverseDigitIter()
	for i := end - 1; i >= start; i-- {
		d, ok := iter()
		if !assert.True(t, ok) {
			return false
		}
		if !assert.Equal(t, i, d.Position) {
			return false
		}
		if !assert.Equal(t, (i+1)%10, d.Value) {
			return false
		}
	}
	_, ok := iter()
	return assert.False(t, ok)
}

func assertEmpty(t *testing.T, s Sequence) {
	t.Helper()
	assertRange(t, s, 0, 0)
}

func hasReverse(s Sequence) bool {
	r, ok := s.(reverseSequence)
	return ok && r.canReverse()
}

func hasSubRange(s Sequence) bool {
	sr, ok := s.(subRangeSequence)
	return ok && sr.canSubRange()
}

func subRange(s Sequence, start, end int) Sequence {
	sr := s.(subRangeSequence)
	return sr.subRange(start, end)
}
