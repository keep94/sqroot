package sqroot

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/keep94/consume2"
	"github.com/stretchr/testify/assert"
)

var (
	// fakeMantissa = 0.12345678901234567890...
	fakeMantissa = Mantissa{spec: funcMantissaSpec(
		func() func() int {
			i := 0
			return func() int {
				i++
				return i % 10
			}
		})}

	// fakeMantissaFiniteDigits = 0.123456789
	fakeMantissaFiniteDigits = Mantissa{spec: funcMantissaSpec(
		func() func() int {
			i := 0
			return func() int {
				if i == 9 {
					return -1
				}
				i++
				return i
			}
		})}

	// fakeMantissaShort = 0.123
	fakeMantissaShort = Mantissa{spec: funcMantissaSpec(
		func() func() int {
			i := 0
			return func() int {
				if i == 3 {
					return -1
				}
				i++
				return i
			}
		})}
)

func TestMantissaReusable(t *testing.T) {
	radican := big.NewInt(5)
	n := SqrtBigInt(radican)
	assert.Equal(t, 1, n.Exponent())
	var answer []int
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 8))
	assert.Equal(t, []int{2, 2, 3, 6, 0, 6, 7, 9}, answer)
	assert.Equal(t, big.NewInt(5), radican)
	radican.SetInt64(7)
	var answer2 []int
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer2), 0, 8))
	assert.Equal(t, []int{2, 2, 3, 6, 0, 6, 7, 9}, answer2)
	assert.Equal(t, big.NewInt(7), radican)
}

func Test2(t *testing.T) {
	var answer []int
	n := Sqrt(2)
	assert.Equal(t, 1, n.Exponent())
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 4, 1, 4, 2, 1, 3, 5, 6, 2}, answer)
}

func Test3(t *testing.T) {
	var answer []int
	n := Sqrt(3)
	assert.Equal(t, 1, n.Exponent())
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 7, 3, 2, 0, 5, 0, 8, 0, 7}, answer)
}

func Test0(t *testing.T) {
	var answer []int
	n := Sqrt(0)
	assert.Zero(t, n)
	n.Mantissa().Send(consume2.AppendTo(&answer))
	assert.Empty(t, answer)
}

func Test1(t *testing.T) {
	var answer []int
	n := Sqrt(1)
	assert.Equal(t, 1, n.Exponent())
	n.Mantissa().Send(consume2.AppendTo(&answer))
	assert.Equal(t, []int{1}, answer)
}

func Test100489(t *testing.T) {
	var answer []int
	n := Sqrt(100489)
	assert.Equal(t, 3, n.Exponent())
	n.Mantissa().Send(consume2.AppendTo(&answer))
	assert.Equal(t, []int{3, 1, 7}, answer)
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
	var answer []int
	n := Sqrt(256)
	assert.Equal(t, 2, n.Exponent())
	n.Mantissa().Send(consume2.AppendTo(&answer))
	assert.Equal(t, []int{1, 6}, answer)
}

func Test40(t *testing.T) {
	var answer []int
	n := Sqrt(40)
	assert.Equal(t, 1, n.Exponent())
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{6, 3, 2, 4, 5, 5, 5, 3, 2, 0}, answer)
}

func Test0026(t *testing.T) {
	var answer []int
	n := SqrtRat(2600, 1000000)
	assert.Equal(t, -1, n.Exponent())
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{5, 0, 9, 9, 0, 1, 9, 5, 1, 3}, answer)
}

func Test026(t *testing.T) {
	var answer []int
	n := SqrtRat(26, 1000)
	assert.Equal(t, 0, n.Exponent())
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 6, 1, 2, 4, 5, 1, 5, 4, 9}, answer)
}

func Test2401Over400(t *testing.T) {
	var answer []int
	n := SqrtRat(2401, 4)
	assert.Equal(t, 2, n.Exponent())
	n.Mantissa().Send(consume2.AppendTo(&answer))
	assert.Equal(t, []int{2, 4, 5}, answer)
}

func Test3Over7(t *testing.T) {
	var answer []int
	n := SqrtRat(3, 7)
	assert.Equal(t, 0, n.Exponent())
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 14))
	assert.Equal(t, []int{6, 5, 4, 6, 5, 3, 6, 7, 0, 7, 0, 7, 9, 7}, answer)
}

func Test3Over70000Reusable(t *testing.T) {
	var answer []int
	radican := big.NewRat(3, 70000)
	n := SqrtBigRat(radican)
	assert.Equal(t, -2, n.Exponent())
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer), 0, 14))
	assert.Equal(t, []int{6, 5, 4, 6, 5, 3, 6, 7, 0, 7, 0, 7, 9, 7}, answer)
	assert.Equal(t, big.NewRat(3, 70000), radican)
	radican.Num().SetInt64(5)
	radican.Denom().SetInt64(80000)
	var answer2 []int
	n.Mantissa().Send(consume2.Slice(consume2.AppendTo(&answer2), 0, 14))
	assert.Equal(t, []int{6, 5, 4, 6, 5, 3, 6, 7, 0, 7, 0, 7, 9, 7}, answer2)
	assert.Equal(t, big.NewInt(5), radican.Num())
	assert.Equal(t, big.NewInt(80000), radican.Denom())
}

func TestNegDenom(t *testing.T) {
	radican := big.NewRat(1, 700)
	radican.Denom().SetInt64(-500)
	radican.Num().SetInt64(3)
	assert.Panics(t, func() { SqrtBigRat(radican) })
}

func TestPrintZeroDigits(t *testing.T) {
	assert.Equal(t, "0", fakeMantissa.Sprint(0))
	assert.Equal(t, "0", fakeMantissa.Sprint(-1))
}

func TestPrintNoOptions(t *testing.T) {
	actual := fakeMantissa.Sprint(12)
	expected := `0.123456789012`
	assert.Equal(t, expected, actual)
}

func TestPrintLessThanOneRow(t *testing.T) {
	actual := fakeMantissa.Sprint(12, ShowCount(true), DigitsPerRow(12))
	expected := `0.123456789012`
	assert.Equal(t, expected, actual)
}

func TestPrintColumns(t *testing.T) {
	actual := fakeMantissa.Sprint(12, DigitsPerColumn(4))
	expected := `0.1234 5678 9012`
	assert.Equal(t, expected, actual)
}

func TestPrintColumnsShow(t *testing.T) {
	actual := fakeMantissa.Sprint(
		12, DigitsPerColumn(5), ShowCount(true))
	expected := `0.12345 67890 12`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows10(t *testing.T) {
	actual := fakeMantissa.Sprint(110, DigitsPerRow(10))
	expected := `0.1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows10Columns(t *testing.T) {
	actual := fakeMantissa.Sprint(
		110, DigitsPerRow(10), DigitsPerColumn(10))
	expected := `0.1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890
  1234567890`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows11Columns(t *testing.T) {
	actual := fakeMantissa.Sprint(
		110, DigitsPerRow(11), DigitsPerColumn(10))
	expected := `0.1234567890 1
  2345678901 2
  3456789012 3
  4567890123 4
  5678901234 5
  6789012345 6
  7890123456 7
  8901234567 8
  9012345678 9
  0123456789 0`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows10Show(t *testing.T) {
	actual := fakeMantissa.Sprint(
		110, DigitsPerRow(10), ShowCount(true))
	expected := `   0.1234567890
 10  1234567890
 20  1234567890
 30  1234567890
 40  1234567890
 50  1234567890
 60  1234567890
 70  1234567890
 80  1234567890
 90  1234567890
100  1234567890`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows10ColumnsShow(t *testing.T) {
	actual := fakeMantissa.Sprint(
		110, DigitsPerRow(10), DigitsPerColumn(10), ShowCount(true))
	expected := `   0.1234567890
 10  1234567890
 20  1234567890
 30  1234567890
 40  1234567890
 50  1234567890
 60  1234567890
 70  1234567890
 80  1234567890
 90  1234567890
100  1234567890`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows11ColumnsShow(t *testing.T) {
	actual := fakeMantissa.Sprint(
		110, DigitsPerRow(11), DigitsPerColumn(10), ShowCount(true))
	expected := `  0.1234567890 1
11  2345678901 2
22  3456789012 3
33  4567890123 4
44  5678901234 5
55  6789012345 6
66  7890123456 7
77  8901234567 8
88  9012345678 9
99  0123456789 0`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows11ColumnsShow109(t *testing.T) {
	actual := fakeMantissa.Sprint(
		109, DigitsPerRow(11), DigitsPerColumn(10), ShowCount(true))
	expected := `  0.1234567890 1
11  2345678901 2
22  3456789012 3
33  4567890123 4
44  5678901234 5
55  6789012345 6
66  7890123456 7
77  8901234567 8
88  9012345678 9
99  0123456789`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows11ColumnsShow111(t *testing.T) {
	actual := fakeMantissa.Sprint(
		111, DigitsPerRow(11), DigitsPerColumn(10), ShowCount(true))
	expected := `   0.1234567890 1
 11  2345678901 2
 22  3456789012 3
 33  4567890123 4
 44  5678901234 5
 55  6789012345 6
 66  7890123456 7
 77  8901234567 8
 88  9012345678 9
 99  0123456789 0
110  1`
	assert.Equal(t, expected, actual)
}

func TestPrinterFewerDigits(t *testing.T) {
	actual := fakeMantissaFiniteDigits.Sprint(
		111, DigitsPerRow(11), DigitsPerColumn(10), ShowCount(true))
	expected := `   0.123456789`
	assert.Equal(t, expected, actual)
}

func TestPrinterNegative(t *testing.T) {
	actual := fakeMantissa.Sprint(
		-3, DigitsPerRow(10), ShowCount(true))
	assert.Equal(t, "0", actual)
}

func TestPrinterCountBytes(t *testing.T) {
	w := &maxBytesWriter{maxBytes: 10000}

	// Prints 20 rows. Each row 10 columns 6 chars per column + (3+2) chars
	// for the margin. 65*20-1=1299 bytes because last line doesn't get a
	// line feed char.
	n, err := fakeMantissa.Fprint(
		w, 1000, DigitsPerRow(50), DigitsPerColumn(5), ShowCount(true))
	assert.Equal(t, 1299, n)
	assert.NoError(t, err)
}

func TestErrorAtAllStages(t *testing.T) {

	// Force an error at each point of the printing
	for i := 0; i < 1299; i++ {
		w := &maxBytesWriter{maxBytes: i}
		n, err := fakeMantissa.Fprint(
			w, 1000, DigitsPerRow(50), DigitsPerColumn(5), ShowCount(true))
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
}

func TestFormat(t *testing.T) {
	actual := fmt.Sprintf("%.14f", fakeMantissa)
	assert.Equal(t, "0.12345678901234", actual)
}

func TestFormatNoPrecision(t *testing.T) {
	actual := fmt.Sprintf("%f", fakeMantissa)
	assert.Equal(t, "0.123456", actual)
}

func TestFormatNoPrecisionCapital(t *testing.T) {
	actual := fmt.Sprintf("%F", fakeMantissa)
	assert.Equal(t, "0.123456", actual)
}

func TestFormatNotInfinite(t *testing.T) {
	actual := fmt.Sprintf("%.14f", fakeMantissaFiniteDigits)
	assert.Equal(t, "0.12345678900000", actual)
}

func TestFormatNotInfiniteNoPrecision(t *testing.T) {
	actual := fmt.Sprintf("%f", fakeMantissaShort)
	assert.Equal(t, "0.123000", actual)
}

func TestFormatWidth(t *testing.T) {
	actual := fmt.Sprintf("%10f", fakeMantissa)
	assert.Equal(t, "  0.123456", actual)
}

func TestFormatShortWidth(t *testing.T) {
	actual := fmt.Sprintf("%5f", fakeMantissa)
	assert.Equal(t, "0.123456", actual)
}

func TestFormatWidthLeftJustify(t *testing.T) {
	actual := fmt.Sprintf("%-10f", fakeMantissa)
	assert.Equal(t, "0.123456  ", actual)
}

func TestFormatWidthAndPrecision(t *testing.T) {
	actual := fmt.Sprintf("%-20.13f", fakeMantissa)
	assert.Equal(t, "0.1234567890123     ", actual)
}

func TestFormatPrecisionSetToZero(t *testing.T) {
	actual := fmt.Sprintf("%.0f", fakeMantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatWidthAndPrecisionNotInfinite(t *testing.T) {
	var builder strings.Builder
	n, err := fmt.Fprintf(&builder, "%-20.13f", fakeMantissaFiniteDigits)
	assert.Equal(t, "0.1234567890000     ", builder.String())
	assert.Equal(t, 20, n)
	assert.NoError(t, err)
}

func TestFormatWidthAndPrecisionNotInfiniteError(t *testing.T) {
	for i := 0; i < 20; i++ {
		w := &maxBytesWriter{maxBytes: i}
		n, err := fmt.Fprintf(w, "%-20.13f", fakeMantissaFiniteDigits)
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
}

func TestFormatNoPrecisionG(t *testing.T) {
	actual := fmt.Sprintf("%g", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestFormatNoPrecisionCapitalG(t *testing.T) {
	actual := fmt.Sprintf("%G", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestFormatNotInfiniteG14(t *testing.T) {
	actual := fmt.Sprintf("%.14g", fakeMantissaFiniteDigits)
	assert.Equal(t, "0.123456789", actual)
}

func TestFormatNotInfiniteG(t *testing.T) {
	actual := fmt.Sprintf("%g", fakeMantissaFiniteDigits)
	assert.Equal(t, "0.123456789", actual)
}

func TestFormatG7(t *testing.T) {
	actual := fmt.Sprintf("%.7g", fakeMantissa)
	assert.Equal(t, "0.1234567", actual)
}

func TestFormatG0(t *testing.T) {
	actual := fmt.Sprintf("%.0g", fakeMantissa)
	assert.Equal(t, "0.1", actual)
}

func TestFormatV(t *testing.T) {
	actual := fmt.Sprintf("%v", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestFormatE(t *testing.T) {
	actual := fmt.Sprintf("%e", fakeMantissa)
	assert.Equal(t, "0.123456e+00", actual)
}

func TestPrintZero(t *testing.T) {
	var mantissa Mantissa
	actual := mantissa.Sprint(45)
	assert.Equal(t, "0", actual)
}

func TestFormatZero(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%f", mantissa)
	assert.Equal(t, "0.000000", actual)
}

func TestFormatZeroPrecision(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.5f", mantissa)
	assert.Equal(t, "0.00000", actual)
}

func TestFormatZeroWidth(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%4.0f", mantissa)
	assert.Equal(t, "   0", actual)
}

func TestFormatZeroG(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%G", mantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatZeroPrecisionG(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.5G", mantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatZeroZeroPrecisionG(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.0G", mantissa)
	assert.Equal(t, "0", actual)
}

func TestFormatZeroV(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%5v", mantissa)
	assert.Equal(t, "    0", actual)
}

func TestFormatZeroE(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprintf("%.5E", mantissa)
	assert.Equal(t, "0.00000E+00", actual)
}

func TestFormatBadVerb(t *testing.T) {
	actual := fmt.Sprintf("%h", fakeMantissa)
	assert.Equal(t, "%!h(mantissa=0.1234567890123456)", actual)
}

func TestPrint(t *testing.T) {
	actual := fmt.Sprint(fakeMantissa)
	assert.Equal(t, "0.1234567890123456", actual)
}

func TestPrintNil(t *testing.T) {
	var mantissa Mantissa
	actual := fmt.Sprint(mantissa)
	assert.Equal(t, "0", actual)
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

func TestNumberZeroValueString(t *testing.T) {
	var number Number
	assert.Equal(t, "0", number.String())
}

func TestNumberFPositiveExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: 5}
	actual := fmt.Sprintf("%f", number)
	assert.Equal(t, "12345.678901", actual)
	actual = fmt.Sprintf("%.1f", number)
	assert.Equal(t, "12345.6", actual)
	actual = fmt.Sprintf("%.0f", number)
	assert.Equal(t, "12345", actual)
}

func TestNumberFPositiveExponentFiniteDigits(t *testing.T) {
	number := Number{mantissa: fakeMantissaFiniteDigits, exponent: 5}
	actual := fmt.Sprintf("%F", number)
	assert.Equal(t, "12345.678900", actual)
}

func TestNumberFNegExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: -5}
	actual := fmt.Sprintf("%f", number)
	assert.Equal(t, "0.000001", actual)
	actual = fmt.Sprintf("%.10f", number)
	assert.Equal(t, "0.0000012345", actual)
	actual = fmt.Sprintf("%.5f", number)
	assert.Equal(t, "0.00000", actual)
	actual = fmt.Sprintf("%.1f", number)
	assert.Equal(t, "0.0", actual)
	actual = fmt.Sprintf("%.0f", number)
	assert.Equal(t, "0", actual)
}

func TestNumberFZero(t *testing.T) {
	var number Number
	actual := fmt.Sprintf("%f", number)
	assert.Equal(t, "0.000000", actual)
	actual = fmt.Sprintf("%.3f", number)
	assert.Equal(t, "0.000", actual)
	actual = fmt.Sprintf("%.1f", number)
	assert.Equal(t, "0.0", actual)
	actual = fmt.Sprintf("%.0f", number)
	assert.Equal(t, "0", actual)
}

func TestNumberGPositiveExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: 5}
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "12345.67890123456", actual)
	actual = fmt.Sprintf("%.8g", number)
	assert.Equal(t, "12345.678", actual)
	actual = fmt.Sprintf("%.5g", number)
	assert.Equal(t, "12345", actual)
	actual = fmt.Sprintf("%.4g", number)
	assert.Equal(t, "0.1234e+05", actual)
	actual = fmt.Sprintf("%.0G", number)
	assert.Equal(t, "0.1E+05", actual)
}

func TestNumberGPositiveExponentShort(t *testing.T) {
	number := Number{mantissa: fakeMantissaShort, exponent: 5}
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "12300", actual)
	actual = fmt.Sprintf("%.5g", number)
	assert.Equal(t, "12300", actual)
	actual = fmt.Sprintf("%.4g", number)
	assert.Equal(t, "0.123e+05", actual)
}

func TestNumberGPositiveExponentFiniteDigits(t *testing.T) {
	number := Number{mantissa: fakeMantissaFiniteDigits, exponent: 5}
	actual := fmt.Sprintf("%G", number)
	assert.Equal(t, "12345.6789", actual)
}

func TestNumberGNegExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: -3}
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "0.0001234567890123456", actual)
	actual = fmt.Sprintf("%.8g", number)
	assert.Equal(t, "0.00012345678", actual)
	actual = fmt.Sprintf("%.0g", number)
	assert.Equal(t, "0.0001", actual)
}

func TestNumberGZero(t *testing.T) {
	var number Number
	actual := fmt.Sprintf("%G", number)
	assert.Equal(t, "0", actual)
	actual = fmt.Sprintf("%.0g", number)
	assert.Equal(t, "0", actual)
}

func TestNumberGLargePosExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: 7}
	actual := fmt.Sprintf("%G", number)
	assert.Equal(t, "0.1234567890123456E+07", actual)
	actual = fmt.Sprintf("%.8g", number)
	assert.Equal(t, "0.12345678e+07", actual)
	actual = fmt.Sprintf("%.0g", number)
	assert.Equal(t, "0.1e+07", actual)
	number = Number{mantissa: fakeMantissa, exponent: 6}
	actual = fmt.Sprintf("%.6g", number)
	assert.Equal(t, "123456", actual)
	number = Number{mantissa: fakeMantissa, exponent: 10}
	actual = fmt.Sprintf("%.10g", number)
	assert.Equal(t, "0.1234567890e+10", actual)
}

func TestNumberGLargePosExponentFiniteDigits(t *testing.T) {
	number := Number{
		mantissa: fakeMantissaFiniteDigits, exponent: 7}
	actual := fmt.Sprintf("%g", number)
	assert.Equal(t, "0.123456789e+07", actual)
}

func TestNumberGLargeNegExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: -4}
	actual := fmt.Sprintf("%G", number)
	assert.Equal(t, "0.1234567890123456E-04", actual)
}

func TestNumberEPositiveExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: 5}
	actual := fmt.Sprintf("%e", number)
	assert.Equal(t, "0.123456e+05", actual)
	actual = fmt.Sprintf("%.1E", number)
	assert.Equal(t, "0.1E+05", actual)
	actual = fmt.Sprintf("%.0e", number)
	assert.Equal(t, "0e+05", actual)
}

func TestNumberEPositiveExponentFiniteDigits(t *testing.T) {
	number := Number{mantissa: fakeMantissaFiniteDigits, exponent: 5}
	actual := fmt.Sprintf("%.14e", number)
	assert.Equal(t, "0.12345678900000e+05", actual)
}

func TestNumberENegExponent(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: -5}
	actual := fmt.Sprintf("%e", number)
	assert.Equal(t, "0.123456e-05", actual)
	actual = fmt.Sprintf("%.1E", number)
	assert.Equal(t, "0.1E-05", actual)
	actual = fmt.Sprintf("%.0e", number)
	assert.Equal(t, "0e-05", actual)
}

func TestNumberEZero(t *testing.T) {
	var number Number
	actual := fmt.Sprintf("%E", number)
	assert.Equal(t, "0.000000E+00", actual)
	actual = fmt.Sprintf("%.1e", number)
	assert.Equal(t, "0.0e+00", actual)
	actual = fmt.Sprintf("%.0e", number)
	assert.Equal(t, "0e+00", actual)
}

func TestNumberWidth(t *testing.T) {
	number := Number{mantissa: fakeMantissa, exponent: 5}
	actual := fmt.Sprintf("%20v", number)
	assert.Equal(t, "   12345.67890123456", actual)
	actual = fmt.Sprintf("%16v", number)
	assert.Equal(t, "12345.67890123456", actual)
	actual = fmt.Sprintf("%-20v", number)
	assert.Equal(t, "12345.67890123456   ", actual)
	actual = fmt.Sprintf("%-16v", number)
	assert.Equal(t, "12345.67890123456", actual)
	actual = fmt.Sprintf("%6.5v", number)
	assert.Equal(t, " 12345", actual)
}

func TestNumberString(t *testing.T) {
	number := Number{mantissa: fakeMantissaFiniteDigits, exponent: 6}
	assert.Equal(t, "123456.789", number.String())
	number = Number{mantissa: fakeMantissa, exponent: 6}
	assert.Equal(t, "123456.7890123456", number.String())
	number = Number{mantissa: fakeMantissa, exponent: 7}
	assert.Equal(t, "0.1234567890123456e+07", number.String())
	number = Number{mantissa: fakeMantissa, exponent: 11}
	assert.Equal(t, "0.1234567890123456e+11", number.String())
	number = Number{mantissa: fakeMantissa, exponent: -3}
	assert.Equal(t, "0.0001234567890123456", number.String())
	number = Number{mantissa: fakeMantissa, exponent: -4}
	assert.Equal(t, "0.1234567890123456e-04", number.String())
	number = Number{}
	assert.Equal(t, "0", number.String())
}

func TestNumberBadVerb(t *testing.T) {
	number := Number{mantissa: fakeMantissaFiniteDigits, exponent: 5}
	actual := fmt.Sprintf("%h", number)
	assert.Equal(t, "%!h(number=12345.6789)", actual)
}

func TestFindFirstN(t *testing.T) {
	number := Sqrt(5)
	hits := number.Mantissa().FindFirstN([]int{9, 7}, 3)
	assert.Equal(t, []int{7, 12, 59}, hits)
}

func TestFindAllSlice(t *testing.T) {
	number := Sqrt(5).WithSignificant(100)
	hits := number.Mantissa().FindAllSlice([]int{9, 7})
	assert.Equal(t, []int{7, 12, 59}, hits)
}

func TestFind(t *testing.T) {
	number := Sqrt(5)
	pattern := []int{9, 7}
	matches := number.Mantissa().Find(pattern)
	pattern[0] = 2
	pattern[1] = 3
	assert.Equal(t, 7, matches())
	assert.Equal(t, 12, matches())
	assert.Equal(t, 59, matches())
}

func TestFindFirstNSingle(t *testing.T) {
	number := Sqrt(11)
	hits := number.Mantissa().FindFirstN([]int{3}, 4)
	assert.Equal(t, []int{0, 1, 10, 13}, hits)
}

func TestFindFirst(t *testing.T) {
	number := Sqrt(2)
	assert.Equal(t, 1, number.Mantissa().FindFirst([]int{4, 1, 4}))
}

func TestFindFirstNotThere(t *testing.T) {
	number := Sqrt(100489)
	assert.Equal(t, -1, number.Mantissa().FindFirst([]int{5}))
}

func TestFindEmptyPattern(t *testing.T) {
	number := Sqrt(2)
	hits := number.Mantissa().FindFirstN(nil, 4)
	assert.Equal(t, []int{0, 1, 2, 3}, hits)
	assert.Equal(t, 0, number.Mantissa().FindFirst(nil))
}

func TestFindFirstNTrickyPattern(t *testing.T) {
	// 12212212122122121221221 ** 2
	radican, ok := new(big.Int).SetString(
		"149138124915706483400311993274596508420730841", 10)
	assert.True(t, ok)
	number := SqrtBigInt(radican)
	hits := number.Mantissa().FindFirstN(
		[]int{1, 2, 2, 1, 2, 1, 2, 2, 1, 2, 2, 1}, 3)
	assert.Equal(t, []int{3, 11}, hits)
}

func TestFindAll(t *testing.T) {
	var indexes []int

	// Find all 000's in sqrt(2). Stop when we find 4 or more zeros in a row
	consumer := consume2.TakeWhile(
		consume2.AppendTo(&indexes),
		func(index int) bool {
			length := len(indexes)
			return !(length >= 2 && (indexes[length-1] == indexes[length-2]+1) && (index > indexes[length-1]+1))
		},
	)
	number := Sqrt(2)
	number.Mantissa().FindAll([]int{0, 0, 0}, consumer)
	assert.Equal(t, []int{
		1879,
		8189,
		8937,
		9025,
		9221,
		10762,
		10801,
		11460,
		12378,
		12413,
		12561,
		12655,
		12656,
	}, indexes)
}

func TestFindZeroMantissa(t *testing.T) {
	var m Mantissa
	assert.Equal(t, -1, m.FindFirst([]int{5}))
	assert.Equal(t, -1, m.FindFirst(nil))
	assert.Empty(t, m.FindFirstN([]int{5}, 3))
	assert.Empty(t, m.FindFirstN(nil, 3))
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

func TestWithSignificantZero(t *testing.T) {
	var n Number
	n2 := n.WithSignificant(5)
	assert.Zero(t, n2)
}

func TestWithSignificantToZero(t *testing.T) {
	n := Sqrt(2).WithSignificant(0)
	assert.Zero(t, n)
}

type maxBytesWriter struct {
	maxBytes     int
	bytesWritten int
}

func (m *maxBytesWriter) Write(p []byte) (n int, err error) {
	length := len(p)
	if length <= m.maxBytes-m.bytesWritten {
		m.bytesWritten += length
		return length, nil
	}
	diff := m.maxBytes - m.bytesWritten
	m.bytesWritten += diff
	return diff, errors.New("Ran out of space")
}

type funcMantissaSpec func() func() int

func (f funcMantissaSpec) Iterator() func() int {
	return f()
}
