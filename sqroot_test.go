package sqroot_test

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/keep94/consume2"
	"github.com/keep94/sqroot"
	"github.com/stretchr/testify/assert"
)

var (
	// fakeMantissa = 0.12345678901234567890...
	fakeMantissa sqroot.Mantissa = func(consumer consume2.Consumer[int]) {
		for consumer.CanConsume() {
			for i := 1; i <= 10; i++ {
				consumer.Consume(i % 10)
			}
		}
	}

	// fakeMantissaFiniteDigits = 0.123456789
	fakeMantissaFiniteDigits sqroot.Mantissa = func(
		consumer consume2.Consumer[int]) {
		for i := 1; i < 10; i++ {
			consumer.Consume(i)
		}
	}
)

func TestMantissaReusable(t *testing.T) {
	mantissa, exp := sqroot.SquareRoot(big.NewInt(5), 0)
	assert.Equal(t, 1, exp)
	var answer []int
	mantissa.Send(consume2.Slice(consume2.AppendTo(&answer), 0, 8))
	assert.Equal(t, []int{2, 2, 3, 6, 0, 6, 7, 9}, answer)
	var answer2 []int
	mantissa.Send(consume2.Slice(consume2.AppendTo(&answer2), 0, 8))
	assert.Equal(t, []int{2, 2, 3, 6, 0, 6, 7, 9}, answer2)
}

func Test2(t *testing.T) {
	var answer []int
	radican := big.NewInt(2)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 1, exp)
	mantissa.Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 4, 1, 4, 2, 1, 3, 5, 6, 2}, answer)
	assert.Equal(t, big.NewInt(2), radican)
}

func Test3(t *testing.T) {
	var answer []int
	radican := big.NewInt(3)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 1, exp)
	mantissa.Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 7, 3, 2, 0, 5, 0, 8, 0, 7}, answer)
	assert.Equal(t, big.NewInt(3), radican)
}

func Test0(t *testing.T) {
	var answer []int
	radican := big.NewInt(0)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 0, exp)
	assert.Nil(t, mantissa)
	mantissa.Send(consume2.AppendTo(&answer))
	assert.Empty(t, answer)
	assert.Equal(t, big.NewInt(0), radican)
}

func Test1(t *testing.T) {
	var answer []int
	radican := big.NewInt(1)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 1, exp)
	mantissa.Send(consume2.AppendTo(&answer))
	assert.Equal(t, []int{1}, answer)
	assert.Equal(t, big.NewInt(1), radican)
}

func Test100489(t *testing.T) {
	var answer []int
	radican := big.NewInt(100489)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 3, exp)
	mantissa.Send(consume2.AppendTo(&answer))
	assert.Equal(t, []int{3, 1, 7}, answer)
	assert.Equal(t, big.NewInt(100489), radican)
}

func TestNegative(t *testing.T) {
	assert.Panics(t, func() { sqroot.SquareRoot(big.NewInt(-1), 0) })
}

func Test256(t *testing.T) {
	var answer []int
	radican := big.NewInt(2560)
	mantissa, exp := sqroot.SquareRoot(radican, -1)
	assert.Equal(t, 2, exp)
	mantissa.Send(consume2.AppendTo(&answer))
	assert.Equal(t, []int{1, 6}, answer)
	assert.Equal(t, big.NewInt(2560), radican)
}

func Test40(t *testing.T) {
	var answer []int
	radican := big.NewInt(4)
	mantissa, exp := sqroot.SquareRoot(radican, 1)
	assert.Equal(t, 1, exp)
	mantissa.Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{6, 3, 2, 4, 5, 5, 5, 3, 2, 0}, answer)
	assert.Equal(t, big.NewInt(4), radican)
}

func Test0026(t *testing.T) {
	var answer []int
	radican := big.NewInt(2600)
	mantissa, exp := sqroot.SquareRoot(radican, -6)
	assert.Equal(t, -1, exp)
	mantissa.Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{5, 0, 9, 9, 0, 1, 9, 5, 1, 3}, answer)
	assert.Equal(t, big.NewInt(2600), radican)
}

func Test026(t *testing.T) {
	var answer []int
	radican := big.NewInt(2600)
	mantissa, exp := sqroot.SquareRoot(radican, -5)
	assert.Equal(t, 0, exp)
	mantissa.Send(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 6, 1, 2, 4, 5, 1, 5, 4, 9}, answer)
	assert.Equal(t, big.NewInt(2600), radican)
}

func ExampleSquareRoot() {
	var mantissaDigits []int

	// Find the square root of 375.2 which is 19.37008002...
	mantissa, exp := sqroot.SquareRoot(big.NewInt(3752), -1)

	mantissa.Send(consume2.Slice(consume2.AppendTo(&mantissaDigits), 0, 10))
	fmt.Println(mantissaDigits)
	fmt.Println(exp)
	// Output:
	// [1 9 3 7 0 0 8 0 0 2]
	// 2
}

func ExampleMantissa_Format() {

	// Find the square root of 5.
	mantissa, exp := sqroot.SquareRoot(big.NewInt(5), 0)

	fmt.Printf("%.50f * 10^%d\n", mantissa, exp)
	// Output:
	// 0.22360679774997896964091736687312762354406183596115 * 10^1
}

func ExampleMantissa_Print() {

	// Find the square root of 2.
	mantissa, exp := sqroot.SquareRoot(big.NewInt(2), 0)

	// Print first 5 digits
	mantissa.Print(5)

	fmt.Printf(" * 10^%d\n", exp)
	// Output:
	// 0.14142 * 10^1
}

func ExampleMantissa_Print_format() {

	// Find the square root of 2.
	mantissa, exp := sqroot.SquareRoot(big.NewInt(2), 0)

	fmt.Printf("10^%d *\n", exp)
	mantissa.Print(
		1000,
		sqroot.DigitsPerRow(50),
		sqroot.DigitsPerColumn(5),
		sqroot.ShowCount(true))
	// Output:
	// 10^1 *
	//    0.14142 13562 37309 50488 01688 72420 96980 78569 67187 53769
	//  50  48073 17667 97379 90732 47846 21070 38850 38753 43276 41572
	// 100  73501 38462 30912 29702 49248 36055 85073 72126 44121 49709
	// 150  99358 31413 22266 59275 05592 75579 99505 01152 78206 05714
	// 200  70109 55997 16059 70274 53459 68620 14728 51741 86408 89198
	// 250  60955 23292 30484 30871 43214 50839 76260 36279 95251 40798
	// 300  96872 53396 54633 18088 29640 62061 52583 52395 05474 57502
	// 350  87759 96172 98355 75220 33753 18570 11354 37460 34084 98847
	// 400  16038 68999 70699 00481 50305 44027 79031 64542 47823 06849
	// 450  29369 18621 58057 84631 11596 66871 30130 15618 56898 72372
	// 500  35288 50926 48612 49497 71542 18334 20428 56860 60146 82472
	// 550  07714 35854 87415 56570 69677 65372 02264 85447 01585 88016
	// 600  20758 47492 26572 26002 08558 44665 21458 39889 39443 70926
	// 650  59180 03113 88246 46815 70826 30100 59485 87040 03186 48034
	// 700  21948 97278 29064 10450 72636 88131 37398 55256 11732 20402
	// 750  45091 22770 02269 41127 57362 72804 95738 10896 75040 18369
	// 800  86836 84507 25799 36472 90607 62996 94138 04756 54823 72899
	// 850  71803 26802 47442 06292 69124 85905 21810 04459 84215 05911
	// 900  20249 44134 17285 31478 10580 36033 71077 30918 28693 14710
	// 950  17111 16839 16581 72688 94197 58716 58215 21282 29518 48847
}

func TestPrintNoOptions(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(&builder, 12)
	expected := `0.123456789012`
	assert.Equal(t, expected, builder.String())
}

func TestPrintColumns(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(&builder, 12, sqroot.DigitsPerColumn(4))
	expected := `0.1234 5678 9012`
	assert.Equal(t, expected, builder.String())
}

func TestPrintColumnsShow(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder, 12, sqroot.DigitsPerColumn(5), sqroot.ShowCount(true))
	expected := `0.12345 67890 12`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(&builder, 110, sqroot.DigitsPerRow(10))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10Columns(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder, 110, sqroot.DigitsPerRow(10), sqroot.DigitsPerColumn(10))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11Columns(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder, 110, sqroot.DigitsPerRow(11), sqroot.DigitsPerColumn(10))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10Show(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder, 110, sqroot.DigitsPerRow(10), sqroot.ShowCount(true))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10ColumnsShow(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder,
		110,
		sqroot.DigitsPerRow(10),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11ColumnsShow(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder,
		110,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11ColumnsShow109(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder,
		109,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11ColumnsShow111(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder,
		111,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
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
	assert.Equal(t, expected, builder.String())
}

func TestPrinterFewerDigits(t *testing.T) {
	var builder strings.Builder
	fakeMantissaFiniteDigits.Fprint(
		&builder,
		111,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
	expected := `   0.123456789`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterNegative(t *testing.T) {
	var builder strings.Builder
	fakeMantissa.Fprint(
		&builder,
		-3,
		sqroot.DigitsPerRow(10),
		sqroot.ShowCount(true))
	assert.Empty(t, builder.String())
}

func TestPrinterCountBytes(t *testing.T) {
	w := &maxBytesWriter{maxBytes: 10000}

	// Prints 20 rows. Each row 10 columns 6 chars per column + (3+2) chars
	// for the margin. 65*20-1=1299 bytes because last line doesn't get a
	// line feed char.
	n, err := fakeMantissa.Fprint(
		w,
		1000,
		sqroot.DigitsPerRow(50),
		sqroot.DigitsPerColumn(5),
		sqroot.ShowCount(true))
	assert.Equal(t, 1299, n)
	assert.NoError(t, err)
}

func TestErrorAtAllStages(t *testing.T) {

	// Force an error at each point of the printing
	for i := 0; i < 1299; i++ {
		w := &maxBytesWriter{maxBytes: i}
		n, err := fakeMantissa.Fprint(
			w,
			1000,
			sqroot.DigitsPerRow(50),
			sqroot.DigitsPerColumn(5),
			sqroot.ShowCount(true))
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
}

func TestFormat(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%.14f", fakeMantissa)
	assert.Equal(t, "0.12345678901234", builder.String())
}

func TestFormatNoPrecision(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%f", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", builder.String())
}

func TestFormatNotInfinite(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%.14f", fakeMantissaFiniteDigits)
	assert.Equal(t, "0.12345678900000", builder.String())
}

func TestFormatNotInfiniteNoPrecision(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%f", fakeMantissaFiniteDigits)
	assert.Equal(t, "0.1234567890000000", builder.String())
}

func TestFormatWidth(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%20f", fakeMantissa)
	assert.Equal(t, "  0.1234567890123456", builder.String())
}

func TestFormatShortWidth(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%10f", fakeMantissa)
	assert.Equal(t, "0.1234567890123456", builder.String())
}

func TestFormatWidthLeftJustify(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%-20f", fakeMantissa)
	assert.Equal(t, "0.1234567890123456  ", builder.String())
}

func TestFormatWidthAndPrecision(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%-20.13f", fakeMantissa)
	assert.Equal(t, "0.1234567890123     ", builder.String())
}

func TestFormatWidthAndPrecisionNotInfinite(t *testing.T) {
	var builder strings.Builder
	n, err := fmt.Fprintf(&builder, "%-20.13f", fakeMantissaFiniteDigits)
	assert.Equal(t, "0.1234567890000     ", builder.String())
	assert.Equal(t, 20, n)
	assert.NoError(t, err)
}

func TestFormatPrecisionSetToZero(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%.0f", fakeMantissa)
	assert.Equal(t, "0", builder.String())
}

func TestFormatWidthAndPrecisionNotInfiniteError(t *testing.T) {
	for i := 0; i < 20; i++ {
		w := &maxBytesWriter{maxBytes: i}
		n, err := fmt.Fprintf(w, "%-20.13f", fakeMantissaFiniteDigits)
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
}

func TestPrintZero(t *testing.T) {
	var mantissa sqroot.Mantissa
	var builder strings.Builder
	mantissa.Fprint(&builder, 45)
	assert.Equal(t, "0", builder.String())
}

func TestFormatZero(t *testing.T) {
	var mantissa sqroot.Mantissa
	var builder strings.Builder
	fmt.Fprintf(&builder, "%f", mantissa)
	assert.Equal(t, "0.0000000000000000", builder.String())
}

func TestFormatZeroPrecision(t *testing.T) {
	var mantissa sqroot.Mantissa
	var builder strings.Builder
	fmt.Fprintf(&builder, "%.5f", mantissa)
	assert.Equal(t, "0.00000", builder.String())
}

func TestFormatZeroPrecisionSmallWidth(t *testing.T) {
	var mantissa sqroot.Mantissa
	var builder strings.Builder
	fmt.Fprintf(&builder, "%5.5f", mantissa)
	assert.Equal(t, "0.00000", builder.String())
}

func TestFormatZeroPrecisionWidth(t *testing.T) {
	var mantissa sqroot.Mantissa
	var builder strings.Builder
	fmt.Fprintf(&builder, "%10.5f", mantissa)
	assert.Equal(t, "   0.00000", builder.String())
}

func TestFormatZeroPrecisionWidthLeft(t *testing.T) {
	var mantissa sqroot.Mantissa
	var builder strings.Builder
	fmt.Fprintf(&builder, "%-10.5f", mantissa)
	assert.Equal(t, "0.00000   ", builder.String())
}

func TestFormatZeroWidth(t *testing.T) {
	var mantissa sqroot.Mantissa
	var builder strings.Builder
	fmt.Fprintf(&builder, "%4.0f", mantissa)
	assert.Equal(t, "   0", builder.String())
}

func TestFormatBadVerb(t *testing.T) {
	var builder strings.Builder
	fmt.Fprintf(&builder, "%g", fakeMantissa)
	assert.Equal(t, "%!g(mantissa)", builder.String())
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
