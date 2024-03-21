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
	assert.Equal(t, []int{3, 1, 7}, exhaust(n.Iterator(), 0))
	assert.Equal(t, []int{3, 1, 7}, exhaust(n.Iterator(), 0))
}

func TestIteratorPersistence(t *testing.T) {
	n := Sqrt(7)
	iter := n.Iterator()
	n = Sqrt(11)
	assert.Equal(t, []int{2, 6, 4, 5}, exhaust(iter, 4))
}

func TestReverse(t *testing.T) {
	// n = 2.2360679
	n := Sqrt(5).WithSignificant(8)
	iter := n.Reverse()
	assert.Equal(t, []int{9, 7, 6, 0, 6, 3, 2, 2}, exhaust(iter, 0))
}

func TestIteratorAt(t *testing.T) {
	n := Sqrt(100489)
	assert.Empty(t, exhaust(n.WithStart(3).Iterator(), 0))
	assert.Equal(t, []int{7}, exhaust(n.WithStart(2).Iterator(), 0))
	assert.Equal(t, []int{3, 1, 7}, exhaust(n.WithStart(0).Iterator(), 0))
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
	assert.Equal(t, []int{3, 2, 7, 8}, exhaust(n.Iterator(), 0))
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

func TestExact(t *testing.T) {
	n := fakeNumber().withExponent(1)
	assert.Equal(t, "1.2345678901234567890", n.WithSignificant(20).Exact())
}

func TestExactShort(t *testing.T) {
	n, _ := NewNumberForTesting([]int{5, 0, 0, 1}, nil, 3)
	assert.Equal(t, "500.1", n.WithSignificant(20).Exact())
	assert.Equal(t, "500", n.WithSignificant(3).Exact())
	assert.Equal(t, "0.50e+03", n.WithSignificant(2).Exact())
	assert.Equal(t, "0.5e+03", n.WithSignificant(1).Exact())
	assert.Equal(t, "0", n.WithSignificant(0).Exact())
	smallN := n.withExponent(-3)
	assert.Equal(t, "0.0005001", smallN.WithSignificant(4).Exact())
	assert.Equal(t, "0.00050", smallN.WithSignificant(2).Exact())
}

func TestExactZero(t *testing.T) {
	var n FiniteNumber
	assert.Equal(t, "0", n.Exact())
}

func TestNewNumberForTesting(t *testing.T) {
	n, err := NewNumberForTesting([]int{1, 0, 2}, []int{0, 0, 3, 4}, 2)
	assert.Equal(t, "10.20034003400340", n.String())
	assert.NoError(t, err)
}

func TestNewNumberForTestingNoExp(t *testing.T) {
	n, err := NewNumberForTesting([]int{1, 0, 2}, []int{0, 0, 3, 4}, 0)
	assert.Equal(t, "0.1020034003400340", n.String())
	assert.NoError(t, err)
}

func TestNewNumberForTestingNegExp(t *testing.T) {
	n, err := NewNumberForTesting([]int{1, 0, 2}, []int{0, 0, 3, 4}, -2)
	assert.Equal(t, "0.001020034003400340", n.String())
	assert.NoError(t, err)
}

func TestNewNumberForTestingNoFixed(t *testing.T) {
	n, err := NewNumberForTesting(nil, []int{1, 0, 3, 4}, 0)
	assert.Equal(t, "0.1034103410341034", n.String())
	assert.NoError(t, err)
}

func TestNewNumberForTestingNoRepeat(t *testing.T) {
	n, err := NewNumberForTesting([]int{1, 0, 2}, nil, 0)
	assert.Equal(t, "0.102", n.String())
	assert.NoError(t, err)
}

func TestNewNumberForTestingRepeatZeros(t *testing.T) {
	n, err := NewNumberForTesting([]int{1, 0, 2}, []int{0}, -2)
	assert.Equal(t, "0.001020000000000000", n.String())
	assert.NoError(t, err)
}

func TestNewNumberForTestingZero(t *testing.T) {
	n, err := NewNumberForTesting(nil, nil, 5)
	assert.True(t, n.IsZero())
	assert.NoError(t, err)
}

func TestNewNumberForTestingLeadingZero(t *testing.T) {
	_, err := NewNumberForTesting(nil, []int{0, 3}, 5)
	assert.Error(t, err)
}

func TestNewNumberForTestingIllegalDigits(t *testing.T) {
	_, err := NewNumberForTesting([]int{10}, nil, 5)
	assert.Error(t, err)
	_, err = NewNumberForTesting(nil, []int{-1}, 5)
	assert.Error(t, err)
}

func TestNewNumber(t *testing.T) {
	// n = 0.12112111211112....
	n := NewNumber(&testgenerator{first: 1, second: 2})
	assert.Equal(t, "0.1211211121111211", n.String())
}

func TestNewNumberIllegal(t *testing.T) {
	n := NewNumber(&testgenerator{first: 5, second: 10})
	assert.Equal(t, "0.5", n.String())
}

func TestNewNumberZero(t *testing.T) {
	n := NewNumber(&testgenerator{first: 10, second: 5, exp: 3})
	assert.True(t, n.IsZero())
}

func TestNewNumberZero2(t *testing.T) {
	n := NewNumber(&testgenerator{first: -1, second: -1, exp: 3})
	assert.True(t, n.IsZero())
}

func TestNewNumberZeroLeadingZero(t *testing.T) {
	n := NewNumber(&testgenerator{first: 0, second: 5, exp: 3})
	assert.True(t, n.IsZero())
}

func TestNewNumberFromBigRat(t *testing.T) {
	var r big.Rat
	r.SetString("2/7")
	assert.Equal(t, "0.285714", fmt.Sprintf("%.6g", NewNumberFromBigRat(&r)))
	r.SetString("47")
	assert.Equal(t, "47", NewNumberFromBigRat(&r).String())
}

func TestNewNumberFromBigRatZero(t *testing.T) {
	var r big.Rat
	assert.True(t, NewNumberFromBigRat(&r).IsZero())
}

func TestNewNumberFromBigRatPanic(t *testing.T) {
	var r big.Rat
	r.SetString("-3/5")
	assert.Panics(t, func() { NewNumberFromBigRat(&r) })
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
	assert.Equal(t, "1.41421", n.Exact())
}

func TestWithSignificantPanics(t *testing.T) {
	n := Sqrt(2)
	assert.Panics(t, func() { n.WithSignificant(-1) })
}

func TestWithSignificantToZero(t *testing.T) {
	assert.Zero(t, *zeroNumber)
	assert.Same(t, zeroNumber, Sqrt(2).WithSignificant(0))
}

func TestZeroNumber(t *testing.T) {
	var n FiniteNumber
	assertEmpty(t, &n)
	assert.Equal(t, -1, n.At(0))
	assert.Zero(t, n.Exponent())
	assert.True(t, n.IsZero())
	assert.Equal(t, "0", n.String())
	assert.Same(t, &n, n.WithSignificant(5))
	assertEmpty(t, n.WithEnd(17))
	assertEmpty(t, n.FiniteWithStart(5))
}

func TestSameNumber(t *testing.T) {
	n := Sqrt(6)
	sixDigits := n.WithSignificant(6)
	assert.Same(t, sixDigits, sixDigits.WithSignificant(6))
	assert.Same(t, sixDigits, sixDigits.WithSignificant(7))
}

func TestNumberWithStartEmpty(t *testing.T) {
	n := Sqrt(19)
	assertEmpty(t, n.WithSignificant(10).FiniteWithStart(300000))
	assertEmpty(t, n.WithSignificant(10).FiniteWithStart(10))
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

func TestNumberSubSequence(t *testing.T) {
	n := fakeNumber()
	assertStartsAt(t, n, 0)
	assertRange(t, n.WithStart(62).WithEnd(404), 62, 404)
}

func TestNumberSubSequenceWithEnd(t *testing.T) {
	n := fakeNumber().WithEnd(357)
	assertRange(t, n, 0, 357)
	assertRange(t, n.WithStart(62).WithEnd(404), 62, 357)
	assertRange(t, n.WithStart(100).WithEnd(150), 100, 150)
	assertEmpty(t, n.WithStart(357).WithEnd(400))
}

func TestNumberSubSequenceWithStart(t *testing.T) {
	seq := fakeNumber().WithStart(423)
	assertStartsAt(t, seq, 423)
	assertRange(t, seq.WithStart(357).WithEnd(504), 423, 504)
	assertRange(t, seq.WithStart(424).WithEnd(425), 424, 425)
}

func TestNumberSubSequenceWithStartAndEnd(t *testing.T) {
	n := fakeNumber().WithEnd(541)
	seq := n.FiniteWithStart(423)
	assertRange(t, seq, 423, 541)
	assertRange(t, seq.WithStart(357).WithEnd(600), 423, 541)
	assertEmpty(t, seq.WithStart(357).WithEnd(358))
	assertRange(t, seq.WithStart(424).WithEnd(425), 424, 425)
	assertEmpty(t, n.FiniteWithStart(541))
	assertEmpty(t, n.FiniteWithStart(542))
}

func TestNumberSubSequenceSame(t *testing.T) {
	n := fakeNumber()
	assert.Same(t, n, n.WithStart(0))
	assert.Same(t, n, n.WithStart(-1))
	endSeq := n.WithEnd(457)
	assert.Same(t, endSeq, endSeq.WithEnd(458))
	startEndSeq := endSeq.WithStart(303)
	assert.Same(t, startEndSeq, startEndSeq.WithStart(-2))
	assert.Same(t, startEndSeq, startEndSeq.WithEnd(458))
	assertEmpty(t, startEndSeq.WithEnd(-3))
}

func TestNumberInfSequenceSame(t *testing.T) {
	n := Sqrt(11)
	s := n.WithStart(10)
	assert.Same(t, s, s.WithStart(10))
	assert.Same(t, s, s.WithStart(9))
}

func TestTypeAssertions(t *testing.T) {
	n := Sqrt(6)
	_, ok := n.(*FiniteNumber)
	assert.False(t, ok)
}

func TestTypeAssertionsWithPositiveStart(t *testing.T) {
	s := Sqrt(6).WithStart(2).WithStart(3).WithStart(1)
	_, ok := s.(FiniteSequence)
	assert.False(t, ok)
}

func TestTypeAssertionsWithEnd(t *testing.T) {
	s := Sqrt(6).WithEnd(1000).WithStart(1)
	_, ok := s.(FiniteSequence)
	assert.True(t, ok)
}

func TestTypeAssertionsWithStartAndEnd(t *testing.T) {
	s := Sqrt(6).WithStart(5).WithEnd(1000).WithStart(10)
	_, ok := s.(FiniteSequence)
	assert.True(t, ok)
}

func TestTypeAssertionsWithSignificant(t *testing.T) {
	s := Sqrt(6).WithSignificant(1000).WithStart(0)
	_, ok := s.(*FiniteNumber)
	assert.True(t, ok)
}

func assertStartsAt(t *testing.T, s Sequence, start int) {
	t.Helper()
	iter := s.Iterator()
	d, ok := iter()
	assert.True(t, ok)
	assert.Equal(t, start, d.Position)
	assert.Equal(t, (start+1)%10, d.Value)
}

func assertRange(t *testing.T, s FiniteSequence, start, end int) {
	t.Helper()
	assertForwardRange(t, s, start, end)
	assertReverseRange(t, s, start, end)
}

func assertForwardRange(t *testing.T, s Sequence, start, end int) bool {
	t.Helper()
	iter := s.Iterator()
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

func assertReverseRange(t *testing.T, s FiniteSequence, start, end int) bool {
	t.Helper()
	iter := s.Reverse()
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

func assertEmpty(t *testing.T, s FiniteSequence) {
	t.Helper()
	assertRange(t, s, 0, 0)
}

func exhaust(iter func() (Digit, bool), max int) []int {
	var result []int
	for digit, ok := iter(); ok; digit, ok = iter() {
		result = append(result, digit.Value)
		if len(result) == max {
			return result
		}
	}
	return result
}

type testgenerator struct {
	first  int
	second int
	exp    int
}

func (g *testgenerator) Generate() (func() int, int) {
	onesLeft := 1
	onesCount := 1
	digits := func() int {
		if onesLeft == 0 {
			onesCount++
			onesLeft = onesCount
			return g.second
		}
		onesLeft--
		return g.first
	}
	return digits, g.exp
}
