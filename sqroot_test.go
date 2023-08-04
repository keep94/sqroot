package sqroot

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNumberNoSideEffects(t *testing.T) {
	radican := big.NewInt(5)
	n := SqrtBigInt(radican)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "2.2360679", fmt.Sprintf("%.8g", n))
	assert.Equal(t, big.NewInt(5), radican)
}

func TestNumberNoSideEffects2(t *testing.T) {
	radican := big.NewInt(5)
	n := SqrtBigInt(radican)
	radican.SetInt64(7)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "2.2360679", fmt.Sprintf("%.8g", n))
	assert.Equal(t, big.NewInt(7), radican)
}

func Test2(t *testing.T) {
	n := Sqrt(2)
	assert.False(t, n.IsZero())
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "1.414213562", fmt.Sprintf("%.10g", n))
}

func Test3(t *testing.T) {
	n := Sqrt(3)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "1.732050807", fmt.Sprintf("%.10g", n))
}

func Test0(t *testing.T) {
	n := Sqrt(0)
	assert.Zero(t, *zeroNumber)
	assert.Same(t, zeroNumber, n)
}

func Test1(t *testing.T) {
	n := Sqrt(1)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "1", fmt.Sprintf("%.10g", n))
}

func Test100489(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, 3, n.Exponent())
	assert.Equal(t, "317", fmt.Sprintf("%.10g", n))
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
	assert.Equal(t, "16", fmt.Sprintf("%.10g", n))
}

func Test40(t *testing.T) {
	n := Sqrt(40)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "6.324555320", fmt.Sprintf("%.10g", n))
}

func Test0026(t *testing.T) {
	n := SqrtRat(2600, 1000000)
	assert.Equal(t, -1, n.Exponent())
	assert.Equal(t, "0.05099019513", fmt.Sprintf("%.10g", n))
}

func Test026(t *testing.T) {
	n := SqrtRat(26, 1000)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.1612451549", fmt.Sprintf("%.10g", n))
}

func Test2401Over400(t *testing.T) {
	n := SqrtRat(2401, 4)
	assert.Equal(t, 2, n.Exponent())
	assert.Equal(t, "24.5", fmt.Sprintf("%.10g", n))
}

func Test3Over7(t *testing.T) {
	n := SqrtRat(3, 7)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.65465367070797", fmt.Sprintf("%.14g", n))
}

func Test3Over70000NoSideEffects(t *testing.T) {
	radican := big.NewRat(3, 70000)
	n := SqrtBigRat(radican)
	assert.Equal(t, -2, n.Exponent())
	assert.Equal(t, "0.0065465367070797", fmt.Sprintf("%.14g", n))
	assert.Equal(t, big.NewRat(3, 70000), radican)
}

func Test3Over70000NoSideEffects2(t *testing.T) {
	radican := big.NewRat(3, 70000)
	n := SqrtBigRat(radican)
	radican.Num().SetInt64(17)
	radican.Denom().SetInt64(80000)
	assert.Equal(t, -2, n.Exponent())
	assert.Equal(t, "0.0065465367070797", fmt.Sprintf("%.14g", n))
	assert.Equal(t, big.NewInt(17), radican.Num())
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
	assertEmpty(t, &n)
	assert.Equal(t, -1, n.At(0))
	assert.Zero(t, n.Exponent())
	assert.True(t, n.IsZero())
	assert.Equal(t, -1, n.Iterator()())
	assert.Equal(t, -1, n.IteratorAt(5)())
	assert.Equal(t, "0", n.String())
	assert.Same(t, &n, n.WithSignificant(5))
	assert.Same(t, &n, n.WithStart(1900000000))
}

func TestSameNumber(t *testing.T) {
	n := Sqrt(6)
	sixDigits := n.WithSignificant(6)
	assert.Same(t, sixDigits, sixDigits.WithSignificant(6))
	assert.Same(t, sixDigits, sixDigits.WithSignificant(7))
}

func TestNumberWithStartEmpty(t *testing.T) {
	n := Sqrt(19)
	assertEmpty(t, n.WithSignificant(10).WithStart(300000))
	assertEmpty(t, n.WithSignificant(10).WithStart(10))
}

func TestNumberWithStartZeroOrNegative(t *testing.T) {
	n := Sqrt(19)
	assert.Same(t, n, n.WithStart(0))
	assert.Same(t, n, n.WithStart(-1))
}

func TestNumberAt(t *testing.T) {
	n := fakeNumber()
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

func TestNumberAtSig(t *testing.T) {
	n := fakeNumber().WithSignificant(357)
	assert.Equal(t, -1, n.At(-1))
	assert.Equal(t, 3, n.At(322))
	assert.Equal(t, 1, n.At(0))
	assert.Equal(t, 4, n.At(303))
	assert.Equal(t, 7, n.At(356))
	assert.Equal(t, -1, n.At(357))
	assert.Equal(t, -1, n.At(2000000000))
}

func TestNumberInterfaces(t *testing.T) {
	n := fakeNumber()
	assertStartsAt(t, n, 0)
	assertRange(t, n.subRange(62, 404), 62, 404)
}

func TestNumberInterfacesSig(t *testing.T) {
	n := fakeNumber().WithSignificant(357)
	assertRange(t, n, 0, 357)
	assertRange(t, n.subRange(62, 404), 62, 357)
	assertRange(t, n.subRange(100, 150), 100, 150)
	assertEmpty(t, n.subRange(357, 400))
}

func TestWithStart(t *testing.T) {
	n := fakeNumber()
	seq := n.WithStart(423)
	assertStartsAt(t, seq, 423)
	assertRange(t, seq.subRange(357, 504), 423, 504)
	assertRange(t, seq.subRange(424, 425), 424, 425)
}

func TestWithStartSig(t *testing.T) {
	n := fakeNumber().WithSignificant(541)
	seq := n.WithStart(423)
	assertRange(t, seq, 423, 541)
	assertRange(t, seq.subRange(357, 600), 423, 541)
	assertEmpty(t, seq.subRange(357, 358))
	assertRange(t, seq.subRange(424, 425), 424, 425)
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
	assertForwardRange(t, s, start, end)
	assertReverseRange(t, s, start, end)
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
	iter := s.reverseDigitIter()
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
