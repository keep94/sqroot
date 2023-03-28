package sqroot

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMantissaReusable(t *testing.T) {
	radican := big.NewInt(5)
	n := SqrtBigInt(radican)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.22360679", n.Mantissa().Sprint(8))
	assert.Equal(t, big.NewInt(5), radican)
	radican.SetInt64(7)
	assert.Equal(t, "0.22360679", n.Mantissa().Sprint(8))
	assert.Equal(t, big.NewInt(7), radican)
}

func Test2(t *testing.T) {
	n := Sqrt(2)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.1414213562", n.Mantissa().Sprint(10))
}

func Test3(t *testing.T) {
	n := Sqrt(3)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.1732050807", n.Mantissa().Sprint(10))
}

func Test0(t *testing.T) {
	n := Sqrt(0)
	assert.Zero(t, *zeroNumber)
	assert.Same(t, zeroNumber, n)
	iter := n.Mantissa().Iterator()
	assert.Equal(t, -1, iter())
}

func Test1(t *testing.T) {
	n := Sqrt(1)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.1", n.Mantissa().Sprint(10))
}

func Test100489(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, 3, n.Exponent())
	assert.Equal(t, "0.317", n.Mantissa().Sprint(10))
}

func Test100489Iterator(t *testing.T) {
	n := Sqrt(100489)
	assert.Equal(t, 3, n.Exponent())
	iter := n.Mantissa().Iterator()
	assert.Equal(t, 3, iter())
	assert.Equal(t, 1, iter())
	assert.Equal(t, 7, iter())
	assert.Equal(t, -1, iter())
	assert.Equal(t, -1, iter())
	iter = n.Mantissa().Iterator()
	assert.Equal(t, 3, iter())
	assert.Equal(t, 1, iter())
	assert.Equal(t, 7, iter())
	assert.Equal(t, -1, iter())
	assert.Equal(t, -1, iter())
}

func TestIteratorPersistence(t *testing.T) {
	n := Sqrt(7)
	m := n.Mantissa()
	iter := m.Iterator()
	m = Sqrt(11).Mantissa()
	assert.Equal(t, 2, iter())
	assert.Equal(t, 6, iter())
	assert.Equal(t, 4, iter())
	assert.Equal(t, 5, iter())
}

func TestNegative(t *testing.T) {
	assert.Panics(t, func() { Sqrt(-1) })
}

func Test256(t *testing.T) {
	n := Sqrt(256)
	assert.Equal(t, 2, n.Exponent())
	assert.Equal(t, "0.16", n.Mantissa().Sprint(10))
}

func Test40(t *testing.T) {
	n := Sqrt(40)
	assert.Equal(t, 1, n.Exponent())
	assert.Equal(t, "0.6324555320", n.Mantissa().Sprint(10))
}

func Test0026(t *testing.T) {
	n := SqrtRat(2600, 1000000)
	assert.Equal(t, -1, n.Exponent())
	assert.Equal(t, "0.5099019513", n.Mantissa().Sprint(10))
}

func Test026(t *testing.T) {
	n := SqrtRat(26, 1000)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.1612451549", n.Mantissa().Sprint(10))
}

func Test2401Over400(t *testing.T) {
	n := SqrtRat(2401, 4)
	assert.Equal(t, 2, n.Exponent())
	assert.Equal(t, "0.245", n.Mantissa().Sprint(10))
}

func Test3Over7(t *testing.T) {
	n := SqrtRat(3, 7)
	assert.Equal(t, 0, n.Exponent())
	assert.Equal(t, "0.65465367070797", n.Mantissa().Sprint(14))
}

func Test3Over70000Reusable(t *testing.T) {
	radican := big.NewRat(3, 70000)
	n := SqrtBigRat(radican)
	assert.Equal(t, -2, n.Exponent())
	assert.Equal(t, "0.65465367070797", n.Mantissa().Sprint(14))
	assert.Equal(t, big.NewRat(3, 70000), radican)
	radican.Num().SetInt64(5)
	radican.Denom().SetInt64(80000)
	assert.Equal(t, "0.65465367070797", n.Mantissa().Sprint(14))
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
	m := Sqrt(2).Mantissa()
	assert.Equal(t, 5, m.At(15))
	assert.Equal(t, 7, m.At(25))
	assert.Equal(t, -1, m.At(-1))
}

func TestAtFinite(t *testing.T) {
	m := Sqrt(100489).Mantissa()
	assert.Equal(t, 3, m.At(0))
	assert.Equal(t, 7, m.At(2))
	assert.Equal(t, -1, m.At(3))
}

func TestZeroMantissa(t *testing.T) {
	var m Mantissa
	assert.Equal(t, -1, m.At(0))
	assert.Zero(t, m.Digits())
	assert.True(t, m.Memoize())
	assert.Same(t, &m, m.WithMemoize())
	assert.Same(t, &m, m.WithSignificant(5))
	assert.Equal(t, -1, m.Iterator()())
	assert.Equal(t, "0", m.String())
}

func TestZeroNumber(t *testing.T) {
	var n Number
	assert.Same(t, &n, n.WithSignificant(5))
	assert.Same(t, &n, n.WithMemoize())
	assert.Zero(t, *zeroMantissa)
	assert.Same(t, zeroMantissa, n.Mantissa())
	assert.Zero(t, n.Exponent())
	assert.Equal(t, "0", n.String())
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
