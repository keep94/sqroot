package sqroot_test

import (
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/keep94/consume2"
	"github.com/keep94/sqroot"
	"github.com/stretchr/testify/assert"
)

func TestMantissaReusable(t *testing.T) {
	mantissa, exp := sqroot.SquareRoot(big.NewInt(5), 0)
	assert.Equal(t, 1, exp)
	var answer []int
	mantissa(consume2.Slice(consume2.AppendTo(&answer), 0, 8))
	assert.Equal(t, []int{2, 2, 3, 6, 0, 6, 7, 9}, answer)
	var answer2 []int
	mantissa(consume2.Slice(consume2.AppendTo(&answer2), 0, 8))
	assert.Equal(t, []int{2, 2, 3, 6, 0, 6, 7, 9}, answer2)
}

func Test2(t *testing.T) {
	var answer []int
	radican := big.NewInt(2)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 1, exp)
	mantissa(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 4, 1, 4, 2, 1, 3, 5, 6, 2}, answer)
	assert.Equal(t, big.NewInt(2), radican)
}

func Test3(t *testing.T) {
	var answer []int
	radican := big.NewInt(3)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 1, exp)
	mantissa(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 7, 3, 2, 0, 5, 0, 8, 0, 7}, answer)
	assert.Equal(t, big.NewInt(3), radican)
}

func Test0(t *testing.T) {
	var answer []int
	radican := big.NewInt(0)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 0, exp)
	mantissa(consume2.AppendTo(&answer))
	assert.Equal(t, []int{0}, answer)
	assert.Equal(t, big.NewInt(0), radican)
}

func Test1(t *testing.T) {
	var answer []int
	radican := big.NewInt(1)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 1, exp)
	mantissa(consume2.AppendTo(&answer))
	assert.Equal(t, []int{1}, answer)
	assert.Equal(t, big.NewInt(1), radican)
}

func Test100489(t *testing.T) {
	var answer []int
	radican := big.NewInt(100489)
	mantissa, exp := sqroot.SquareRoot(radican, 0)
	assert.Equal(t, 3, exp)
	mantissa(consume2.AppendTo(&answer))
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
	mantissa(consume2.AppendTo(&answer))
	assert.Equal(t, []int{1, 6}, answer)
	assert.Equal(t, big.NewInt(2560), radican)
}

func Test40(t *testing.T) {
	var answer []int
	radican := big.NewInt(4)
	mantissa, exp := sqroot.SquareRoot(radican, 1)
	assert.Equal(t, 1, exp)
	mantissa(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{6, 3, 2, 4, 5, 5, 5, 3, 2, 0}, answer)
	assert.Equal(t, big.NewInt(4), radican)
}

func Test0026(t *testing.T) {
	var answer []int
	radican := big.NewInt(2600)
	mantissa, exp := sqroot.SquareRoot(radican, -6)
	assert.Equal(t, -1, exp)
	mantissa(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{5, 0, 9, 9, 0, 1, 9, 5, 1, 3}, answer)
	assert.Equal(t, big.NewInt(2600), radican)
}

func Test026(t *testing.T) {
	var answer []int
	radican := big.NewInt(2600)
	mantissa, exp := sqroot.SquareRoot(radican, -5)
	assert.Equal(t, 0, exp)
	mantissa(consume2.Slice(consume2.AppendTo(&answer), 0, 10))
	assert.Equal(t, []int{1, 6, 1, 2, 4, 5, 1, 5, 4, 9}, answer)
	assert.Equal(t, big.NewInt(2600), radican)
}

func ExampleSquareRoot() {
	var mantissaDigits []int

	// Find the square root of 375.2 which is 19.37008002...
	mantissa, exp := sqroot.SquareRoot(big.NewInt(3752), -1)

	mantissa(consume2.Slice(consume2.AppendTo(&mantissaDigits), 0, 10))
	fmt.Println(mantissaDigits)
	fmt.Println(exp)
	// Output:
	// [1 9 3 7 0 0 8 0 0 2]
	// 2
}

func ExamplePrinter() {
	// Find the square root of 2.
	mantissa, exp := sqroot.SquareRoot(big.NewInt(2), 0)
	mantissa(sqroot.NewPrinter(os.Stdout, 10))
	fmt.Printf(" * 10^%d\n", exp)
	// Output:
	// 0.1414213562 * 10^1
}

func ExamplePrinter_format() {
	// Find the square root of 2.
	mantissa, exp := sqroot.SquareRoot(big.NewInt(2), 0)
	fmt.Printf("10^%d *\n", exp)
	mantissa(sqroot.NewPrinter(
		os.Stdout,
		30,
		sqroot.DigitsPerRow(10),
		sqroot.DigitsPerColumn(5),
		sqroot.ShowCount(true)))
	// Output:
	// 10^1 *
	//   0.14142 13562
	// 10  37309 50488
	// 20  01688 72420
}

func TestPrinterNoOptions(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(&builder, 12)
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `0.012345678901`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterColumns(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(&builder, 12, sqroot.DigitsPerColumn(4))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `0.0123 4567 8901`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterColumnsShow(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder, 12, sqroot.DigitsPerColumn(5), sqroot.ShowCount(true))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `0.01234 56789 01`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(&builder, 110, sqroot.DigitsPerRow(10))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `0.0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10Columns(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder, 110, sqroot.DigitsPerRow(10), sqroot.DigitsPerColumn(10))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `0.0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789
  0123456789`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11Columns(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder, 110, sqroot.DigitsPerRow(11), sqroot.DigitsPerColumn(10))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `0.0123456789 0
  1234567890 1
  2345678901 2
  3456789012 3
  4567890123 4
  5678901234 5
  6789012345 6
  7890123456 7
  8901234567 8
  9012345678 9`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10Show(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder, 110, sqroot.DigitsPerRow(10), sqroot.ShowCount(true))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `   0.0123456789
 10  0123456789
 20  0123456789
 30  0123456789
 40  0123456789
 50  0123456789
 60  0123456789
 70  0123456789
 80  0123456789
 90  0123456789
100  0123456789`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows10ColumnsShow(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder,
		110,
		sqroot.DigitsPerRow(10),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `   0.0123456789
 10  0123456789
 20  0123456789
 30  0123456789
 40  0123456789
 50  0123456789
 60  0123456789
 70  0123456789
 80  0123456789
 90  0123456789
100  0123456789`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11ColumnsShow(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder,
		110,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `  0.0123456789 0
11  1234567890 1
22  2345678901 2
33  3456789012 3
44  4567890123 4
55  5678901234 5
66  6789012345 6
77  7890123456 7
88  8901234567 8
99  9012345678 9`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11ColumnsShow109(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder,
		109,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `  0.0123456789 0
11  1234567890 1
22  2345678901 2
33  3456789012 3
44  4567890123 4
55  5678901234 5
66  6789012345 6
77  7890123456 7
88  8901234567 8
99  9012345678`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterRows11ColumnsShow111(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder,
		111,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	expected := `   0.0123456789 0
 11  1234567890 1
 22  2345678901 2
 33  3456789012 3
 44  4567890123 4
 55  5678901234 5
 66  6789012345 6
 77  7890123456 7
 88  8901234567 8
 99  9012345678 9
110  0`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterFewerDigits(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder,
		111,
		sqroot.DigitsPerRow(11),
		sqroot.DigitsPerColumn(10),
		sqroot.ShowCount(true))
	for i := 0; i < 10; i++ {
		p.Consume(i)
	}
	expected := `   0.0123456789`
	assert.Equal(t, expected, builder.String())
}

func TestPrinterNegative(t *testing.T) {
	var builder strings.Builder
	p := sqroot.NewPrinter(
		&builder, -3, sqroot.DigitsPerRow(10), sqroot.ShowCount(true))
	for p.CanConsume() {
		for i := 0; i < 10; i++ {
			p.Consume(i)
		}
	}
	assert.Empty(t, builder.String())
}
