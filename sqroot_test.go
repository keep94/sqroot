package sqroot

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

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
	assert.Zero(t, n)
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

func TestPrintP(t *testing.T) {
	p := new(PositionsBuilder).AddRange(5, 8).AddRange(45, 50).Build()
	actual := GetDigits(fakeMantissa, p).Sprint(
		DigitsPerRow(11), DigitsPerColumn(10))
	expected := `  0......678.. .
44  .67890`
	assert.Equal(t, expected, actual)
}

func TestPrintP2(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(0, 2).AddRange(3, 5).AddRange(8, 11).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(11), DigitsPerColumn(10))
	expected := `0.12.45...90 1`
	assert.Equal(t, expected, actual)
}

func TestPrintPGaps(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(22, 44).AddRange(66, 77).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(11), DigitsPerColumn(10))
	expected := `22  3456789012 3
33  4567890123 4
66  7890123456 7`
	assert.Equal(t, expected, actual)
}

func TestPrintPGaps2(t *testing.T) {
	var pb PositionsBuilder
	pb.AddRange(0, 10).AddRange(11, 21).AddRange(33, 43).AddRange(66, 76)
	p := pb.Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(11), DigitsPerColumn(10))
	expected := `  0.1234567890 .
11  2345678901 .
33  4567890123 .
66  7890123456`
	assert.Equal(t, expected, actual)
}

func TestPrintPGaps3(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(21, 33).AddRange(65, 77).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(
		DigitsPerRow(11), DigitsPerColumn(10), MissingDigit('-'))
	expected := `11  ---------- 2
22  3456789012 3
55  ---------- 6
66  7890123456 7`
	assert.Equal(t, expected, actual)
}

func TestPrintPNoShowCount(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(21, 33).AddRange(65, 77).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(
		DigitsPerRow(11), DigitsPerColumn(10), ShowCount(false))
	expected := `0........... .
  .......... 2
  3456789012 3
  .......... .
  .......... .
  .......... 6
  7890123456 7`
	assert.Equal(t, expected, actual)
}

func TestPrintPNarrow(t *testing.T) {
	var pb PositionsBuilder
	p := pb.Add(3).Add(5).Add(8).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint(DigitsPerRow(1))
	expected := `3  4
5  6
8  9`
	assert.Equal(t, expected, actual)
}

func TestPrintPDefaults(t *testing.T) {
	var pb PositionsBuilder
	p := pb.AddRange(0, 75).Build()
	digits := GetDigits(fakeMantissa, p)
	actual := digits.Sprint()
	expected := `  0.12345 67890 12345 67890 12345 67890 12345 67890 12345 67890
50  12345 67890 12345 67890 12345`
	assert.Equal(t, expected, actual)
}

func TestPrintPTooShort(t *testing.T) {
	n := Sqrt(100489)
	p := new(PositionsBuilder).AddRange(3, 5).Build()
	digits := GetDigits(n.Mantissa(), p)
	assert.Zero(t, digits)
	assert.Empty(t, digits.Sprint())
}

func TestPrintPZero(t *testing.T) {
	var m Mantissa
	p := new(PositionsBuilder).AddRange(3, 5).Build()
	digits := GetDigits(m, p)
	assert.Zero(t, digits)
	assert.Empty(t, digits.Sprint())
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
	hits := FindFirstN(number.Mantissa(), []int{9, 7}, 3)
	assert.Equal(t, []int{7, 12, 59}, hits)
}

func TestFindAll(t *testing.T) {
	number := Sqrt(5).WithSignificant(100)
	hits := FindAll(number.Mantissa(), []int{9, 7})
	assert.Equal(t, []int{7, 12, 59}, hits)
}

func TestFind(t *testing.T) {
	number := Sqrt(5)
	pattern := []int{9, 7}
	matches := Find(number.Mantissa(), pattern)
	pattern[0] = 2
	pattern[1] = 3
	assert.Equal(t, 7, matches())
	assert.Equal(t, 12, matches())
	assert.Equal(t, 59, matches())
}

func TestFindFirstNSingle(t *testing.T) {
	number := Sqrt(11)
	hits := FindFirstN(number.Mantissa(), []int{3}, 4)
	assert.Equal(t, []int{0, 1, 10, 13}, hits)
}

func TestFindFirst(t *testing.T) {
	number := Sqrt(2)
	assert.Equal(t, 1, FindFirst(number.Mantissa(), []int{4, 1, 4}))
}

func TestFindFirstNotThere(t *testing.T) {
	number := Sqrt(100489)
	assert.Equal(t, -1, FindFirst(number.Mantissa(), []int{5}))
}

func TestFindEmptyPattern(t *testing.T) {
	number := Sqrt(2)
	hits := FindFirstN(number.Mantissa(), nil, 4)
	assert.Equal(t, []int{0, 1, 2, 3}, hits)
	assert.Equal(t, 0, FindFirst(number.Mantissa(), nil))
}

func TestFindEmptyPatternIterator(t *testing.T) {
	number := Sqrt(2).WithSignificant(4)
	iter := Find(number.Mantissa(), nil)
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
		number.Mantissa(),
		[]int{1, 2, 2, 1, 2, 1, 2, 2, 1, 2, 2, 1},
		3,
	)
	assert.Equal(t, []int{3, 11}, hits)
}

func TestFindLast(t *testing.T) {
	number := Sqrt(5).WithSignificant(1000)
	assert.Equal(t, 936, FindLast(number.Mantissa(), []int{9, 7}))
	assert.Equal(t, -1, FindLast(number.Mantissa(), []int{0, 1, 2, 3, 4}))
}

func TestFindLastN(t *testing.T) {
	number := Sqrt(5).WithSignificant(1000)
	hits := FindLastN(number.Mantissa(), []int{9, 7}, 3)
	assert.Equal(t, []int{936, 718, 600}, hits)
	hits = FindLastN(Sqrt(5).WithSignificant(1300).Mantissa(), []int{9, 7}, 3)
	assert.Equal(t, []int{1276, 1221, 936}, hits)
	hits = FindLastN(number.Mantissa(), nil, 4)
	assert.Equal(t, []int{999, 998, 997, 996}, hits)
	hits = FindLastN(number.Mantissa(), []int{1, 2, 3}, 3)
	assert.Equal(t, []int{815, 579}, hits)
	hits = FindLastN(number.Mantissa(), []int{1, 2, 3}, 0)
	assert.Empty(t, hits)
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
	digits2 := Sqrt(5).WithSignificant(1300).Mantissa().Digits()
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

func TestFindZeroMantissa(t *testing.T) {
	var m Mantissa
	assert.Equal(t, -1, FindFirst(m, []int{5}))
	assert.Equal(t, -1, FindFirst(m, nil))
	assert.Empty(t, FindFirstN(m, []int{5}, 3))
	assert.Empty(t, FindFirstN(m, nil, 3))
	assert.Empty(t, FindAll(m, []int{5}))
	assert.Empty(t, FindAll(m, nil))
	assert.Equal(t, -1, FindLast(m, []int{5}))
	assert.Equal(t, -1, FindLast(m, nil))
	assert.Empty(t, FindLastN(m, []int{5}, 3))
	assert.Empty(t, FindLastN(m, nil, 3))
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
	assert.Zero(t, n.WithSignificant(5))
}

func TestWithSignificantToZero(t *testing.T) {
	assert.Zero(t, Sqrt(2).WithSignificant(0))
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

func TestDigits(t *testing.T) {
	m := Sqrt(2).Mantissa()
	positions := new(PositionsBuilder).Add(25).Add(15).Add(50).Build()
	digits := GetDigits(m, positions)
	assert.Equal(t, 5, digits.At(15))
	assert.Equal(t, 7, digits.At(25))
	assert.Equal(t, 4, digits.At(50))
	assert.Equal(t, -1, digits.At(26))
	checkFullIter(t, digits.Items(), 15, 5, 25, 7, 50, 4)
	checkFullIter(t, digits.WithStart(25).Items(), 25, 7, 50, 4)
	assert.Equal(t, 15, digits.Min())
	assert.Equal(t, 50, digits.Max())
	checkFullIter(t, digits.ReverseItems(), 50, 4, 25, 7, 15, 5)
	checkFullIter(t, digits.WithEnd(50).ReverseItems(), 25, 7, 15, 5)
}

func TestGetDigitsFromDigits(t *testing.T) {
	m := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.AddRange(0, 100).AddRange(200, 300).AddRange(400, 500)
	digits := GetDigits(m, pb.Build())
	pb.AddRange(100, 200).AddRange(300, 400).AddRange(500, 600)
	assert.Zero(t, GetDigits(digits, pb.Build()))
}

func TestGetDigitsFromDigits2(t *testing.T) {
	m := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.AddRange(100, 200).AddRange(300, 400).AddRange(500, 600)
	digits := GetDigits(m, pb.Build())

	// force GetDigits to do a full scan rather than picking
	pb.AddRange(0, 101).AddRange(200, 301).Add(500).AddRange(1000, 2000000000)
	digits = GetDigits(digits, pb.Build())
	expected := GetDigits(m, pb.Add(100).Add(300).Add(500).Build())
	assert.Equal(t, expected.Sprint(), digits.Sprint())
}

func TestGetDigitsFromDigitsPick(t *testing.T) {
	m := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.AddRange(100, 200).AddRange(300, 400).AddRange(500, 600)
	digits := GetDigits(m, pb.Build())
	pb.Add(153).Add(200)
	digits = GetDigits(digits, pb.Build())
	expected := GetDigits(m, pb.Add(153).Build())
	assert.Equal(t, expected.Sprint(), digits.Sprint())
}

func TestDigitsFinite(t *testing.T) {
	m := Sqrt(2).WithSignificant(50).Mantissa()
	positions := new(PositionsBuilder).Add(25).Add(15).Add(50).Build()
	digits := GetDigits(m, positions)
	assert.Equal(t, 5, digits.At(15))
	assert.Equal(t, 7, digits.At(25))
	assert.Equal(t, -1, digits.At(50))
	checkFullIter(t, digits.Items(), 15, 5, 25, 7)
	checkFullIter(t, digits.ReverseItems(), 25, 7, 15, 5)
}

func TestDigitsNone(t *testing.T) {
	m := Sqrt(2).Mantissa()
	var positions Positions
	assert.Zero(t, GetDigits(m, positions))
}

func TestDigitsNoneZeroMantissa(t *testing.T) {
	var m Mantissa
	var p Positions
	var pb PositionsBuilder
	assert.Zero(t, GetDigits(m, p))
	assert.Zero(t, GetDigits(m, pb.AddRange(0, 100).Build()))
}

func TestDigitsNoneFromDigits(t *testing.T) {
	digits := Sqrt(2).WithSignificant(1000).Mantissa().Digits()
	assert.Equal(t, 1000, digits.Len())
	var zeroPosits Positions
	zeroDigits := GetDigits(digits, zeroPosits)
	assert.Zero(t, zeroDigits)
	assert.Zero(t, GetDigits(zeroDigits, zeroPosits))
	var pb PositionsBuilder
	assert.Zero(t, GetDigits(zeroDigits, pb.AddRange(0, 100).Build()))
}

func TestDigitsZero(t *testing.T) {
	var digits Digits
	iter := digits.Items()
	_, ok := iter()
	assert.False(t, ok)
	iter = digits.ReverseItems()
	_, ok = iter()
	assert.False(t, ok)
	assert.Equal(t, -1, digits.At(0))
	assert.Empty(t, digits.Sprint())
	assert.Equal(t, -1, digits.Min())
	assert.Equal(t, -1, digits.Max())
	assert.Equal(t, 0, digits.Len())
	assert.Empty(t, FindAll(digits, nil))
	assert.Empty(t, FindAll(digits, []int{5}))
	assert.Equal(t, 0, digits.WithStart(5).Len())
	assert.Equal(t, 0, digits.WithEnd(5).Len())
}

func TestDigitsBinary(t *testing.T) {
	mantissa := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.AddRange(1000, 2000).AddRange(5000, 5999).AddRange(10000, 10999).Add(11000)
	digits := GetDigits(mantissa, pb.Build())
	arr, err := digits.MarshalBinary()
	assert.Equal(t, 1509, len(arr))
	assert.NoError(t, err)
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinary2(t *testing.T) {
	mantissa := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.Add(50).Add(25).Add(15).Add(0).AddRange(100, 102)
	digits := GetDigits(mantissa, pb.Build())
	arr, err := digits.MarshalBinary()
	assert.NoError(t, err)
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinary3(t *testing.T) {
	mantissa := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.AddRange(0, 4).Add(26)
	digits := GetDigits(mantissa, pb.Build())
	arr, err := digits.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, 6, len(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinaryZero(t *testing.T) {
	var digits Digits
	arr, err := digits.MarshalBinary()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalBinary(arr))
	assert.Zero(t, copy)
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsBinaryEmpty(t *testing.T) {
	var digits Digits
	assert.Error(t, digits.UnmarshalBinary(nil))
}

func TestDigitsBinaryBadVersion(t *testing.T) {
	var digits Digits
	assert.Error(t, digits.UnmarshalBinary([]byte{51}))
}

func TestDigitsBinaryUnmarshalErrors(t *testing.T) {
	// A very large position gap results in a negative position delta.
	data := []byte{0xbb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x44}
	var digits Digits
	assert.Error(t, digits.UnmarshalBinary(data))
	data = []byte{0xbb, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x6d}
	assert.Error(t, digits.UnmarshalBinary(data))
	data = []byte{0xbb, 0xff}
	assert.Error(t, digits.UnmarshalBinary(data))
}

func TestDigitsText(t *testing.T) {
	mantissa := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.AddRange(0, 4).Add(26)
	digits := GetDigits(mantissa, pb.Build())
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:1414[26]2", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsText2(t *testing.T) {
	mantissa := Sqrt(2).Mantissa()
	var pb PositionsBuilder
	pb.AddRange(3, 6).AddRange(10, 12)
	digits := GetDigits(mantissa, pb.Build())
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:[3]421[10]37", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsText3(t *testing.T) {
	digits := Sqrt(2).WithSignificant(6).Mantissa().Digits()
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:141421", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsTextZero(t *testing.T) {
	var digits Digits
	arr, err := digits.MarshalText()
	assert.NoError(t, err)
	assert.Equal(t, "v1:", string(arr))
	var copy Digits
	assert.NoError(t, copy.UnmarshalText(arr))
	assert.Zero(t, copy)
	assert.Equal(t, digits.Sprint(), copy.Sprint())
}

func TestDigitsTextUnmarshalErrors(t *testing.T) {
	text := []byte("v1:12345[")
	var digits Digits
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v1:12345[67]")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v1:12abc")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v1:12345[6a]")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("v2:")
	assert.Error(t, digits.UnmarshalText(text))
	text = []byte("")
	assert.Error(t, digits.UnmarshalText(text))
}

func TestDigitsWithStartAndEnd(t *testing.T) {
	digits := Sqrt(2).WithSignificant(1000).Mantissa().Digits()
	assert.NotEqual(t, -1, digits.At(700))
	assert.NotEqual(t, -1, digits.At(399))
	digits = digits.WithStart(200).WithEnd(900)
	digits = digits.WithStart(400).WithEnd(700).WithStart(300).WithEnd(800)
	assert.Equal(t, 300, digits.Len())
	assert.Equal(t, 400, digits.Min())
	assert.Equal(t, 699, digits.Max())
	assert.Equal(t, 700, digits.limit())
	assert.Equal(t, -1, digits.At(700))
	assert.Equal(t, -1, digits.At(399))
	assert.NotEqual(t, -1, digits.At(400))
	assert.NotEqual(t, -1, digits.At(699))
	assert.Len(t, digits.Sprint(), 389)
}

func TestDigitsWithStartZero(t *testing.T) {
	digits := Sqrt(2).WithSignificant(1000).Mantissa().Digits()
	digits = digits.WithStart(1000)
	assert.Equal(t, 0, digits.Len())
	assert.Equal(t, -1, digits.Min())
	assert.Equal(t, -1, digits.Max())
	assert.Equal(t, -1, digits.At(500))
}

func TestDigitsWithEndZero(t *testing.T) {
	digits := Sqrt(2).WithSignificant(1000).Mantissa().Digits()
	digits = digits.WithEnd(0)
	assert.Equal(t, 0, digits.Len())
	assert.Equal(t, -1, digits.Min())
	assert.Equal(t, -1, digits.Max())
	assert.Equal(t, -1, digits.At(500))
}

func TestDigitsBuilder(t *testing.T) {
	var builder digitsBuilder
	assert.NoError(t, builder.AddDigit(0, 3))
	assert.NoError(t, builder.AddDigit(1, 1))
	assert.NoError(t, builder.AddDigit(2, 4))
	assert.NoError(t, builder.AddDigit(3, 1))
	assert.NoError(t, builder.AddDigit(6, 4))
	assert.NoError(t, builder.AddDigit(7, 1))
	assert.NoError(t, builder.AddDigit(8, 4))
	digits := builder.Build()
	assert.NoError(t, builder.AddDigit(0, 4))
	assert.NoError(t, builder.AddDigit(1, 6))
	assert.NoError(t, builder.AddDigit(2, 6))
	assert.NoError(t, builder.AddDigit(3, 9))
	assert.NoError(t, builder.AddDigit(4, 2))
	assert.NoError(t, builder.AddDigit(5, 0))
	feigenbaum := builder.Build()
	actual := digits.Sprint(
		DigitsPerRow(0), DigitsPerColumn(0), ShowCount(false))
	assert.Equal(t, "0.3141..414", actual)
	assert.Equal(t, 7, digits.Len())
	matches := FindAll(digits, []int{1, 4})
	assert.Equal(t, []int{1, 7}, matches)
	assert.Equal(t, 6, feigenbaum.Len())
}

func TestDigitBuilderErrors(t *testing.T) {
	var builder digitsBuilder
	assert.Error(t, builder.AddDigit(-1, 3))
	assert.Error(t, builder.AddDigit(0, -1))
	assert.Error(t, builder.AddDigit(0, 10))
	assert.NoError(t, builder.AddDigit(0, 2))
	assert.Error(t, builder.AddDigit(0, 1))
	digits := builder.Build()
	assert.Equal(t, 1, digits.Len())
}

func TestPositionsBuilder(t *testing.T) {
	var pb PositionsBuilder
	pb.AddRange(0, 2).Add(4).Add(10).AddRange(-1, 3)
	pb.AddRange(15, 17)
	pb.AddRange(-3, -1)
	pb.AddRange(13, 15)
	pb.AddRange(17, 19)
	pb.Add(1)
	pb.AddRange(20, 25)
	pb.AddRange(21, 26)
	pb.AddRange(22, 23)
	assert.True(t, pb.unsorted)
	p := pb.Build()
	assert.False(t, pb.unsorted)
	assert.Len(t, pb.ranges, 0)
	expected := []positionRange{
		{Start: 0, End: 3},
		{Start: 4, End: 5},
		{Start: 10, End: 11},
		{Start: 13, End: 19},
		{Start: 20, End: 26},
	}
	assert.Equal(t, expected, p.ranges)
	assert.Equal(t, 26, p.limit())
}

func TestPositionsBuilderSorted(t *testing.T) {
	var pb PositionsBuilder
	pb.AddRange(0, 3).AddRange(1, 4).Add(2)
	pb.AddRange(4, 6).AddRange(10, 15).AddRange(6, 6).AddRange(7, 5)
	pb.AddRange(12, 17)
	for i := 100; i < 200; i++ {
		pb.Add(i)
	}
	assert.False(t, pb.unsorted)
	assert.Len(t, pb.ranges, 3)
	p := pb.Build()
	assert.False(t, pb.unsorted)
	assert.Len(t, pb.ranges, 0)
	expected := []positionRange{
		{Start: 0, End: 6},
		{Start: 10, End: 17},
		{Start: 100, End: 200},
	}
	assert.Equal(t, expected, p.ranges)
	assert.Equal(t, 200, p.limit())
}

func TestPositionsBuilderNegative(t *testing.T) {
	var pb PositionsBuilder
	pb.Add(-1)
	assert.Zero(t, pb.Build())
}

func TestPositionsBuilderZero(t *testing.T) {
	var pb PositionsBuilder
	assert.Zero(t, pb.Build())
}

func TestDigitLookup(t *testing.T) {
	n := Sqrt(11).WithSignificant(10000)
	pattern := []int{4, 5, 7}
	finds := FindAll(n.Mantissa(), pattern)
	assert.Len(t, finds, 7)
	var pb PositionsBuilder
	for _, find := range finds {
		pb.AddRange(find, find+3)
	}
	digits := GetDigits(n.Mantissa(), pb.Build())
	iter := digits.Items()
	count := 0
	for digit, ok := iter(); ok; digit, ok = iter() {
		assert.Equal(t, pattern[count%3], digit.Value)
		assert.Equal(t, finds[count/3]+count%3, digit.Position)
		count++
	}
	assert.Equal(t, 7*3, count)
}

func TestFindWithDigits(t *testing.T) {
	n := Sqrt(5)
	// '000' in Sqrt(5) at 424 569 3663 4506 4601 6113 7173 9110 9114
	var pb PositionsBuilder
	pb.AddRange(7173, 7176).AddRange(9110, 9117).AddRange(4500, 4600)
	digits := GetDigits(n.Mantissa(), pb.Build())
	finds := FindAll(digits, []int{0, 0, 0})
	assert.Equal(t, []int{4506, 7173, 9110, 9114}, finds)
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

func checkFullIter(t *testing.T, iter func() (Digit, bool), positionDigits ...int) {
	t.Helper()
	var actual []int
	for digit, ok := iter(); ok; digit, ok = iter() {
		actual = append(actual, digit.Position, digit.Value)
	}
	assert.Equal(t, positionDigits, actual)
	_, ok := iter()
	assert.False(t, ok)
}
