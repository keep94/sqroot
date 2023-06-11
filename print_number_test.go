package sqroot

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	// fakeNumber = 0.12345678901234567890...
	fakeNumber = &Number{spec: funcNumberSpec(
		func() func() int {
			i := 0
			return func() int {
				i++
				return i % 10
			}
		})}
)

func TestPrintZeroDigits(t *testing.T) {
	assert.Equal(t, "", Sprint(fakeNumber, UpTo(0)))
	assert.Equal(t, "", Sprint(fakeNumber, UpTo(-1)))
}

func TestPrintNoOptions(t *testing.T) {
	actual := Sprint(fakeNumber, UpTo(12))
	expected := `0.12345 67890 12`
	assert.Equal(t, expected, actual)
}

func TestPrintLessThanOneRow(t *testing.T) {
	actual := Sprint(
		fakeNumber, UpTo(12), DigitsPerRow(12), DigitsPerColumn(0))
	expected := `0.123456789012`
	assert.Equal(t, expected, actual)
}

func TestPrintColumns(t *testing.T) {
	actual := Sprint(
		fakeNumber,
		UpTo(12),
		DigitsPerColumn(4),
		DigitsPerRow(0),
		ShowCount(false))
	expected := `0.1234 5678 9012`
	assert.Equal(t, expected, actual)
}

func TestPrintColumnsShow(t *testing.T) {
	actual := Sprint(fakeNumber, UpTo(12), DigitsPerColumn(5), DigitsPerRow(0))
	expected := `0.12345 67890 12`
	assert.Equal(t, expected, actual)
}

func TestPrinterRows10(t *testing.T) {
	actual := Sprint(
		fakeNumber,
		UpTo(110),
		DigitsPerRow(10),
		DigitsPerColumn(0),
		ShowCount(false))
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
	actual := Sprint(
		fakeNumber,
		UpTo(110),
		DigitsPerRow(10),
		DigitsPerColumn(10),
		ShowCount(false))
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
	actual := Sprint(
		fakeNumber,
		UpTo(110),
		DigitsPerRow(11),
		DigitsPerColumn(10),
		ShowCount(false))
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
	actual := Sprint(
		fakeNumber, UpTo(110), DigitsPerRow(10), DigitsPerColumn(0))
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
	actual := Sprint(
		fakeNumber, UpTo(110), DigitsPerRow(10), DigitsPerColumn(10))
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
	actual := Sprint(
		fakeNumber, UpTo(110), DigitsPerRow(11), DigitsPerColumn(10))
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
	actual := Sprint(
		fakeNumber, UpTo(109), DigitsPerRow(11), DigitsPerColumn(10))
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
	actual := Sprint(
		fakeNumber,
		UpTo(111),
		DigitsPerRow(11),
		DigitsPerColumn(10),
		ShowCount(true))
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

func TestPrinterWithBetween(t *testing.T) {
	actual := Sprint(
		fakeNumber,
		Between(50, 70),
		DigitsPerRow(11),
		DigitsPerColumn(10))
	expected := `44  ......1234 5
55  6789012345 6
66  7890`
	assert.Equal(t, expected, actual)
}

func TestPrinterWithPositions(t *testing.T) {
	var pb PositionsBuilder
	actual := Sprint(
		fakeNumber,
		pb.Add(45).Add(48).AddRange(50, 52).Build(),
		DigitsPerRow(11),
		DigitsPerColumn(10),
		MissingDigit('-'))
	expected := `44  -6--9-12`
	assert.Equal(t, expected, actual)
}

func TestPrinterFewerDigits(t *testing.T) {
	actual := Sprint(
		fakeNumber.WithSignificant(9),
		UpTo(111),
		DigitsPerRow(11),
		DigitsPerColumn(10))
	expected := `   0.123456789`
	assert.Equal(t, expected, actual)
}

func TestPrinterNegative(t *testing.T) {
	actual := Sprint(fakeNumber, UpTo(-3), DigitsPerRow(10))
	assert.Equal(t, "", actual)
}

func TestPrinterCountBytes(t *testing.T) {
	w := &maxBytesWriter{maxBytes: 10000}

	// Prints 20 rows. Each row 10 columns 6 chars per column + (3+2) chars
	// for the margin. 65*20-1=1299 bytes because last line doesn't get a
	// line feed char.
	n, err := Fprint(w, fakeNumber, UpTo(1000))
	assert.Equal(t, 1299, n)
	assert.NoError(t, err)
}

func TestErrorAtAllStages(t *testing.T) {

	// Force an error at each point of the printing
	for i := 0; i < 1299; i++ {
		w := &maxBytesWriter{maxBytes: i}
		n, err := Fprint(w, fakeNumber, UpTo(1000))
		assert.Equal(t, i, n)
		assert.Error(t, err)
	}
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

type funcNumberSpec func() func() int

func (f funcNumberSpec) IteratorAt(index int) func() int {
	return fastForward(f(), index)
}

func (f funcNumberSpec) At(index int) int {
	return simpleAt(f(), index)
}

func (f funcNumberSpec) IsMemoize() bool { return false }

func (f funcNumberSpec) FirstN(n int) []int {
	panic("FirstN not supported")
}
